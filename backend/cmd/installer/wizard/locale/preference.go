package locale

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type preferenceFile struct {
	Language Language `json:"language"`
}

// PreferencePath is where an explicit language choice (Ctrl+L) is remembered:
// <user config dir>/pentagi/installer.json.
func PreferencePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve user config dir: %w", err)
	}
	return filepath.Join(dir, "pentagi", "installer.json"), nil
}

// LoadPreference returns the saved language; ok is false when there is no
// readable, supported choice (missing file, corrupt JSON, removed language).
func LoadPreference(path string) (Language, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	var pref preferenceFile
	if err := json.Unmarshal(data, &pref); err != nil || !IsSupported(pref.Language) {
		return "", false
	}

	return pref.Language, true
}

// SavePreference stores lang atomically so an interrupted write never leaves a
// truncated file behind.
func SavePreference(path string, lang Language) error {
	if !IsSupported(lang) {
		return fmt.Errorf("unsupported language %q", lang)
	}

	data, err := json.Marshal(preferenceFile{Language: lang})
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("cannot create %s: %w", filepath.Dir(path), err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".installer-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}
