//go:build !darwin && !windows

package locale

// osPreferredLanguages has no source beyond the locale environment variables on
// Linux and other Unix systems.
func osPreferredLanguages() []string {
	return nil
}
