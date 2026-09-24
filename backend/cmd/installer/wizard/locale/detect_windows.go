//go:build windows

package locale

import "golang.org/x/sys/windows"

// osPreferredLanguages returns the user's preferred Windows UI languages.
func osPreferredLanguages() []string {
	tags, err := windows.GetUserPreferredUILanguages(windows.MUI_LANGUAGE_NAME)
	if err != nil {
		return nil
	}
	return tags
}
