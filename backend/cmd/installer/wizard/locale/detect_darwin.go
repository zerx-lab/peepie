//go:build darwin

package locale

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// osPreferredLanguages reads the macOS UI language list (System Settings →
// Language & Region). Used only when no locale environment variable is set,
// which happens e.g. in terminals launched without "Set locale environment
// variables on startup".
func osPreferredLanguages() []string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "defaults", "read", "-g", "AppleLanguages").Output()
	if err != nil {
		return nil
	}

	return parseAppleLanguages(string(out))
}

// parseAppleLanguages parses the plist array printed by `defaults read`:
//
//	(
//	    "zh-Hans-CN",
//	    en
//	)
func parseAppleLanguages(out string) []string {
	var tags []string
	for _, line := range strings.Split(out, "\n") {
		tag := strings.Trim(strings.TrimSpace(line), `(),"`)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
