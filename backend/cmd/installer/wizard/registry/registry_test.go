package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pentagi/cmd/installer/checker"
	"pentagi/cmd/installer/files"
	"pentagi/cmd/installer/processor"
	"pentagi/cmd/installer/state"
	"pentagi/cmd/installer/wizard/controller"
	"pentagi/cmd/installer/wizard/locale"
	"pentagi/cmd/installer/wizard/styles"
	"pentagi/cmd/installer/wizard/window"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TestEveryScreenFitsTheWindowInEveryLanguage renders all screens in every
// supported language at common terminal sizes. Translations are longer or
// double-width (CJK), so a screen that lays out text by hand could overflow or
// panic only in one language; a screen whose text is cached outside locale
// would render the same in every language.
func TestEveryScreenFitsTheWindowInEveryLanguage(t *testing.T) {
	t.Cleanup(func() { _ = locale.SetLanguage(locale.DefaultLanguage) })

	envPath := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(envPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	appState, err := state.NewState(envPath)
	if err != nil {
		t.Fatal(err)
	}
	checkResult := checker.CheckResult{}
	f := files.NewFiles()

	// views in the source language, to prove every screen actually switches language
	sourceViews := map[string]string{}

	for _, lang := range locale.SupportedLanguages {
		if err := locale.SetLanguage(lang); err != nil {
			t.Fatal(err)
		}

		for _, size := range [][2]int{{80, 24}, {120, 40}, {200, 50}} {
			width, height := size[0], size[1]

			// fresh screens per language, as App.switchLanguage does
			w := window.New()
			w.SetWindowSize(width, height)
			c := controller.NewController(appState, f, checkResult)
			p := processor.NewProcessorModel(appState, c.GetChecker(), f)
			r := NewRegistry(c, styles.New(), w, f, p).(*registry)

			contentWidth, contentHeight := w.GetContentSize()
			resize := tea.WindowSizeMsg{Width: contentWidth, Height: contentHeight}

			for id, screen := range r.screens {
				screen.Init()
				screen.Update(resize)
				view := screen.View()

				if strings.TrimSpace(view) == "" {
					t.Errorf("%s %dx%d %s: empty view", lang, width, height, id)
				}
				key := fmt.Sprintf("%dx%d %s", width, height, id)
				if lang == locale.DefaultLanguage {
					sourceViews[key] = view
				} else if view == sourceViews[key] {
					t.Errorf("%s %s: view is identical to %s, the screen ignores the language", lang, key, locale.DefaultLanguage)
				}
				for i, line := range strings.Split(view, "\n") {
					if lineWidth := lipgloss.Width(line); lineWidth > contentWidth {
						t.Errorf("%s %dx%d %s: line %d is %d cells wide, content area is %d: %q",
							lang, width, height, id, i, lineWidth, contentWidth, line)
						break
					}
				}
			}
		}
	}
}
