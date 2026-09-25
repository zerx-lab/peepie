package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type EnvVar struct {
	Name      string // variable name
	Value     string // variable value
	IsChanged bool   // was the value changed manually
	IsComment bool   // is this line a comment (not saved, updated on value change)
	Default   string // default value from config struct (not saved, used for display)
	Line      int    // line number in file (-1 if not present, e.g. for new vars)
}

func (e *EnvVar) IsDefault() bool {
	return e.Value == e.Default || (e.Value == "" && e.Default != "")
}

func (e *EnvVar) IsPresent() bool {
	return e.Line != -1
}

type EnvFile interface {
	Del(name string)
	Set(name, value string)
	Get(name string) (EnvVar, bool)
	GetAll() map[string]EnvVar
	SetAll(vars map[string]EnvVar)
	Save(path string) error
	Clone() EnvFile
}

type envFile struct {
	vars map[string]*EnvVar
	perm os.FileMode
	raw  string
	mx   *sync.Mutex
}

func (e *envFile) Del(name string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	delete(e.vars, name)
}

func (e *envFile) Set(name, value string) {
	e.mx.Lock()
	defer e.mx.Unlock()

	name, value = trim(name), trim(value)

	if envVar, ok := e.vars[name]; !ok {
		e.vars[name] = &EnvVar{
			Name:      name,
			Value:     value,
			IsChanged: true,
			Line:      -1,
		}
	} else {
		if envVar.Value != value {
			envVar.IsChanged = true
			envVar.Value = value
		}
	}
}

func (e *envFile) Get(name string) (EnvVar, bool) {
	e.mx.Lock()
	defer e.mx.Unlock()

	if envVar, ok := e.vars[name]; !ok {
		return EnvVar{
			Name: name,
			Line: -1,
		}, false
	} else {
		return *envVar, true
	}
}

func (e *envFile) GetAll() map[string]EnvVar {
	e.mx.Lock()
	defer e.mx.Unlock()

	result := make(map[string]EnvVar, len(e.vars))
	for name, envVar := range e.vars {
		result[name] = *envVar
	}

	return result
}

func (e *envFile) SetAll(vars map[string]EnvVar) {
	e.mx.Lock()
	defer e.mx.Unlock()

	for name := range vars {
		envVar := vars[name]
		e.vars[name] = &envVar
	}
}

func (e *envFile) Save(path string) error {
	e.mx.Lock()
	defer e.mx.Unlock()

	// The backend (Web UI -> Settings -> System) writes runtime settings into
	// this same file while the installer may be open; rebase onto the current
	// disk content so saving here never reverts those edits.
	if disk, err := os.ReadFile(path); err == nil && string(disk) != e.raw {
		e.rebase(string(disk))
	}

	// check if there are any changes to the file to avoid unnecessary writes
	curRaw := e.raw
	e.patchRaw()
	isChanged := e.raw != curRaw
	for _, envVar := range e.vars {
		if envVar.IsChanged {
			isChanged = true
			break
		}
	}
	if !isChanged {
		return nil
	}

	backupDir := filepath.Join(filepath.Dir(path), ".bak")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return fmt.Errorf("'%s' is a directory", path)
	} else if err == nil {
		// Copy rather than rename: the file is rewritten in place below so a
		// Docker single-file bind mount of it (the pentagi container mounts
		// .env to mirror Web UI settings) keeps pointing at the live file.
		current, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read current file state: %w", err)
		}
		curTimeStr := time.Unix(time.Now().Unix(), 0).Format("20060102150405")
		backupPath := filepath.Join(backupDir, fmt.Sprintf("%s.%s", filepath.Base(path), curTimeStr))
		if err := os.WriteFile(backupPath, current, info.Mode().Perm()); err != nil {
			return fmt.Errorf("failed to create backup file: %w", err)
		}
	}

	if err := os.WriteFile(path, []byte(e.raw), e.perm); err != nil {
		return fmt.Errorf("failed to write new file state: %w", err)
	}

	for _, envVar := range e.vars {
		envVar.IsChanged = false
	}

	return nil
}

// rebase adopts disk as the new base content: variables with pending changes
// keep their new value (re-anchored to their line on disk), every other
// variable takes its value from disk, and variables only present on disk are
// kept instead of being dropped as deleted.
func (e *envFile) rebase(disk string) {
	fresh := loadVars(disk)
	for name, envVar := range e.vars {
		diskVar, onDisk := fresh[name]
		switch {
		case envVar.IsChanged && onDisk:
			envVar.Line = diskVar.Line
		case envVar.IsChanged:
			envVar.Line = -1
		case onDisk:
			envVar.Value, envVar.Line, envVar.IsComment = diskVar.Value, diskVar.Line, diskVar.IsComment
			envVar.IsChanged = diskVar.IsChanged
		default:
			envVar.Value, envVar.Line = "", -1
		}
	}
	for name, diskVar := range fresh {
		if _, ok := e.vars[name]; !ok {
			e.vars[name] = diskVar
		}
	}
	e.raw = disk
}

func (e *envFile) Clone() EnvFile {
	e.mx.Lock()
	defer e.mx.Unlock()

	clone := envFile{
		vars: make(map[string]*EnvVar, len(e.vars)),
		perm: e.perm,
		raw:  e.raw,
		mx:   &sync.Mutex{},
	}

	for name, envVar := range e.vars {
		v := *envVar
		clone.vars[name] = &v
	}

	return &clone
}

func (e *envFile) patchRaw() {
	lines := strings.Split(e.raw, "\n")
	hasLastEmpty := len(lines) > 0 && trim(lines[len(lines)-1]) == ""
	for ldx := len(lines) - 1; ldx >= 0 && trim(lines[ldx]) == ""; ldx-- {
		lines = lines[:ldx]
	}

	// First pass: mark lines for deletion and update existing variables
	var linesToDelete []int
	for ldx, line := range lines {
		line = trim(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		varName := trim(parts[0])

		// Check if this variable still exists
		if envVar, exists := e.vars[varName]; exists {
			if envVar.IsChanged && !envVar.IsComment {
				lines[ldx] = fmt.Sprintf("%s=%s", envVar.Name, envVar.Value)
				envVar.Line = ldx
			}
		} else {
			// Mark line for deletion
			linesToDelete = append(linesToDelete, ldx)
		}
	}

	// Remove lines in reverse order to maintain indices
	for i := len(linesToDelete) - 1; i >= 0; i-- {
		lineIdx := linesToDelete[i]
		lines = append(lines[:lineIdx], lines[lineIdx+1:]...)

		// Update line numbers for remaining variables
		for _, envVar := range e.vars {
			if envVar.Line > lineIdx {
				envVar.Line--
			}
		}
	}

	// Second pass: add new variables
	for _, envVar := range e.vars {
		if !envVar.IsChanged || envVar.IsComment {
			continue
		}

		line := fmt.Sprintf("%s=%s", envVar.Name, envVar.Value)
		if !envVar.IsPresent() || envVar.Line >= len(lines) {
			lines = append(lines, line)
			envVar.Line = len(lines) - 1
		} else {
			lines[envVar.Line] = line
		}
	}

	if hasLastEmpty {
		lines = append(lines, "")
	}
	e.raw = strings.Join(lines, "\n")
}

func trim(value string) string {
	return strings.Trim(value, "\n\r\t ")
}
