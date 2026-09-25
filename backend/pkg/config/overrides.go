package config

import (
	"strconv"
	"sync"
)

// Overrides is a thread-safe, in-memory store of runtime configuration values
// that were changed through the Web UI (Settings) and persisted in the
// system_settings table, keyed by (category, key). Every *Config carries one
// live Overrides instance (see NewConfig), so any code holding a *Config
// pointer automatically observes updates applied via Set/Load without needing
// to re-read the process environment or restart.
//
// Only fields explicitly wired to read through Overrides become hot-
// reloadable; the underlying struct field on Config (populated once from the
// environment at boot) remains the fallback default for every category/key
// that has no row in system_settings yet.
type Overrides struct {
	mu   sync.RWMutex
	data map[string]map[string]string // category -> key -> value
}

// NewOverrides creates an empty overrides store.
func NewOverrides() *Overrides {
	return &Overrides{data: make(map[string]map[string]string)}
}

// Load replaces the entire snapshot atomically. Used at boot and by the
// periodic refresh so that a multi-replica deployment converges on the
// latest values written by any instance within the refresh interval.
func (o *Overrides) Load(rows map[string]map[string]string) {
	snapshot := make(map[string]map[string]string, len(rows))
	for category, kv := range rows {
		m := make(map[string]string, len(kv))
		for k, v := range kv {
			m[k] = v
		}
		snapshot[category] = m
	}

	o.mu.Lock()
	o.data = snapshot
	o.mu.Unlock()
}

// Set applies a single (category, key) -> value change immediately, without
// waiting for the next periodic refresh. Callers persisting the change to the
// database should call this right after a successful write.
func (o *Overrides) Set(category, key, value string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.data == nil {
		o.data = make(map[string]map[string]string)
	}
	if o.data[category] == nil {
		o.data[category] = make(map[string]string)
	}
	o.data[category][key] = value
}

// Delete removes a single override, reverting reads of that key back to the
// environment-sourced struct field default.
func (o *Overrides) Delete(category, key string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.data[category] != nil {
		delete(o.data[category], key)
	}
}

// Get returns the raw stored value for (category, key), if any. Safe to call
// on a nil *Overrides (e.g. a Config built directly in tests without going
// through NewConfig): behaves as an always-empty store.
func (o *Overrides) Get(category, key string) (string, bool) {
	if o == nil {
		return "", false
	}

	o.mu.RLock()
	defer o.mu.RUnlock()

	kv, ok := o.data[category]
	if !ok {
		return "", false
	}
	v, ok := kv[key]
	return v, ok
}

// GetString returns the override value for (category, key), or fallback when
// no override has been set.
func (o *Overrides) GetString(category, key, fallback string) string {
	if v, ok := o.Get(category, key); ok {
		return v
	}
	return fallback
}

// GetBool returns the override value parsed as bool, or fallback.
func (o *Overrides) GetBool(category, key string, fallback bool) bool {
	if v, ok := o.Get(category, key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

// GetInt returns the override value parsed as int, or fallback.
func (o *Overrides) GetInt(category, key string, fallback int) int {
	if v, ok := o.Get(category, key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// Snapshot returns a copy of every key/value currently overridden for a
// category, e.g. for building a GraphQL response that must show live values.
// Like Get, it is safe to call on a nil *Overrides.
func (o *Overrides) Snapshot(category string) map[string]string {
	if o == nil {
		return map[string]string{}
	}

	o.mu.RLock()
	defer o.mu.RUnlock()

	kv := o.data[category]
	out := make(map[string]string, len(kv))
	for k, v := range kv {
		out[k] = v
	}
	return out
}

// Categories returns a copy of the whole snapshot, grouped by category.
func (o *Overrides) Categories() map[string]map[string]string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	out := make(map[string]map[string]string, len(o.data))
	for category, kv := range o.data {
		m := make(map[string]string, len(kv))
		for k, v := range kv {
			m[k] = v
		}
		out[category] = m
	}
	return out
}
