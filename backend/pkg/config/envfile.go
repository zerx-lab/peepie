package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// envAssignRe matches an active (uncommented) assignment line and captures the
// variable name; godotenv also accepts an optional `export ` prefix.
var envAssignRe = regexp.MustCompile(`^\s*(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=`)

// ReadEnvFileSettings parses the .env file at path and returns the catalogued
// runtime settings it assigns, grouped by category and normalized (an empty
// assignment resolves to the variable's default, as it does at boot).
// Variables outside the settings catalogue are ignored.
func ReadEnvFileSettings(path string) (map[string]map[string]string, error) {
	values, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}

	out := make(map[string]map[string]string)
	for name, value := range values {
		def, ok := settingsByEnv[name]
		if !ok {
			continue
		}
		if out[def.category] == nil {
			out[def.category] = make(map[string]string)
		}
		out[def.category][def.key] = NormalizeSetting(def.category, def.key, value)
	}

	return out, nil
}

// WriteEnvFileSettings assigns each NAME=value of values (keyed by environment
// variable name) in the existing .env file at path: active assignments are
// rewritten in place, commented-out or absent variables are appended. Every
// other line (comments, ordering, unrelated variables) is preserved.
//
// The file is truncated and rewritten in place rather than replaced via
// rename, so a single-file Docker bind mount of the host .env keeps pointing
// at the same inode the host sees.
func WriteEnvFileSettings(path string, values map[string]string) error {
	for name, value := range values {
		if strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("value of %s must not contain line breaks", name)
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("'%s' is not a regular file", path)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(raw), "\n")
	written := make(map[string]bool, len(values))
	for i, line := range lines {
		m := envAssignRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		value, ok := values[m[1]]
		if !ok {
			continue
		}
		lines[i] = m[1] + "=" + quoteEnvValue(value)
		written[m[1]] = true
	}

	missing := make([]string, 0, len(values))
	for name := range values {
		if !written[name] {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)

	if len(missing) > 0 {
		// Append before the trailing newline (the empty last element), and
		// make sure the file ends with one afterwards.
		if n := len(lines); n > 0 && lines[n-1] == "" {
			lines = lines[:n-1]
		}
		for _, name := range missing {
			lines = append(lines, name+"="+quoteEnvValue(values[name]))
		}
		lines = append(lines, "")
	}

	updated := strings.Join(lines, "\n")
	if updated == string(raw) {
		return nil
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(updated); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// quoteEnvValue leaves plain values bare (the format the installer TUI writes
// and reads) and single-quotes anything godotenv or Docker Compose would
// otherwise reinterpret: whitespace, inline comments, quotes, `$` expansion.
func quoteEnvValue(value string) string {
	if !strings.ContainsAny(value, " \t#'\"\\$`") {
		return value
	}
	if !strings.Contains(value, "'") {
		return "'" + value + "'"
	}
	escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`, `$`, `\$`).Replace(value)
	return `"` + escaped + `"`
}

// SettingsEnvFileStatus returns the absolute path of the settings .env file and
// whether it exists and can be written, i.e. whether Web UI changes are being
// mirrored into it.
func (c *Config) SettingsEnvFileStatus() (string, bool) {
	if c.SettingsEnvFile == "" {
		return "", false
	}

	path, err := filepath.Abs(c.SettingsEnvFile)
	if err != nil {
		path = c.SettingsEnvFile
	}

	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return path, false
	}
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return path, false
	}
	_ = f.Close()

	return path, true
}

// settingsMu serializes ApplySettings against ReloadOverrides. Without it a
// periodic reload that read the database/file just before a Web UI write could
// publish its stale snapshot right after the write, reverting the change until
// the next tick.
var settingsMu sync.Mutex

// ApplySettings publishes Web UI changes of one category: it normalizes values
// (in place), mirrors them into the settings .env file, runs persist (the
// database write) and finally applies them to Overrides. A missing settings
// file is skipped (the deployment does not expose one); any other file error
// aborts before anything is persisted, since the file would otherwise revert
// the change on the next reload.
func (c *Config) ApplySettings(category string, values map[string]string, persist func() error) error {
	settingsMu.Lock()
	defer settingsMu.Unlock()

	byEnv := make(map[string]string, len(values))
	for key, value := range values {
		name, ok := SettingEnvName(category, key)
		if !ok {
			return fmt.Errorf("setting %s/%s has no environment variable mapping", category, key)
		}
		values[key] = NormalizeSetting(category, key, value)
		byEnv[name] = values[key]
	}

	if c.SettingsEnvFile != "" && len(byEnv) > 0 {
		err := WriteEnvFileSettings(c.SettingsEnvFile, byEnv)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to update settings file %s: %w", c.SettingsEnvFile, err)
		}
	}

	if persist != nil {
		if err := persist(); err != nil {
			return err
		}
	}

	for key, value := range values {
		c.Overrides.Set(category, key, value)
	}
	return nil
}

// ReloadOverrides replaces the live overrides with the rows returned by
// loadRows (the persisted system_settings table) layered under the settings
// .env file. The file wins because it is what the installer TUI and operators
// edit directly, and every Web UI change is written to both. A missing file is
// not an error; a loadRows error or any other file read failure leaves the
// current overrides untouched and is returned, so a file caught mid-write never
// flips live settings back to stale values.
func (c *Config) ReloadOverrides(loadRows func() (map[string]map[string]string, error)) error {
	settingsMu.Lock()
	defer settingsMu.Unlock()

	dbRows, err := loadRows()
	if err != nil {
		return err
	}

	merged := make(map[string]map[string]string, len(dbRows))
	put := func(category, key, value string) {
		if merged[category] == nil {
			merged[category] = make(map[string]string)
		}
		merged[category][key] = NormalizeSetting(category, key, value)
	}

	for category, kv := range dbRows {
		for key, value := range kv {
			put(category, key, value)
		}
	}

	if c.SettingsEnvFile != "" {
		fileRows, err := ReadEnvFileSettings(c.SettingsEnvFile)
		switch {
		case err == nil:
			for category, kv := range fileRows {
				for key, value := range kv {
					put(category, key, value)
				}
			}
		case errors.Is(err, fs.ErrNotExist):
		default:
			return fmt.Errorf("failed to read settings file %s: %w", c.SettingsEnvFile, err)
		}
	}

	c.Overrides.Load(merged)
	return nil
}
