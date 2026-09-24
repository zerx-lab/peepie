package wizard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pentagi/cmd/installer/checker"
	"pentagi/cmd/installer/files"
	"pentagi/cmd/installer/state"
	"pentagi/cmd/installer/wizard/locale"
	"pentagi/cmd/installer/wizard/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newTestApp(t *testing.T) (*App, string) {
	t.Helper()
	t.Cleanup(func() { _ = locale.SetLanguage(locale.LanguageEnglish) })
	if err := locale.SetLanguage(locale.LanguageEnglish); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	appState, err := state.NewState(envPath)
	if err != nil {
		t.Fatal(err)
	}

	prefPath := filepath.Join(dir, "config", "installer.json")
	app := NewApp(appState, checker.CheckResult{}, files.NewFiles(), prefPath)
	app.Update(tea.WindowSizeMsg{Width: 200, Height: 50})

	return app, prefPath
}

func pressCtrlL(app *App) {
	app.Update(tea.KeyMsg{Type: tea.KeyCtrlL})
}

func TestCtrlLSwitchesLanguageRebuildsScreensAndPersists(t *testing.T) {
	app, prefPath := newTestApp(t)

	if view := app.View(); !strings.Contains(view, "Ctrl+L: Switch to 简体中文") {
		t.Fatalf("footer does not offer the language switch:\n%s", view)
	}
	registryBefore := app.registry

	pressCtrlL(app)

	if got := locale.Current(); got != locale.LanguageChineseSimplified {
		t.Fatalf("language = %s, want zh-CN", got)
	}
	if app.registry == registryBefore {
		t.Fatal("screens were not rebuilt after switching language")
	}
	if saved, ok := locale.LoadPreference(prefPath); !ok || saved != locale.LanguageChineseSimplified {
		t.Fatalf("saved preference = %q, %v", saved, ok)
	}
	view := app.View()
	if !strings.Contains(view, "Ctrl+L: 切换到 English") || strings.Contains(view, "Ctrl+Q: Exit") {
		t.Fatalf("view is not rendered in Chinese:\n%s", view)
	}

	// header, content and wrapped footer hints must stay inside a narrow terminal
	app.Update(tea.WindowSizeMsg{Width: 60, Height: 24})
	for i, line := range strings.Split(app.View(), "\n") {
		if width := lipgloss.Width(line); width > 60 {
			t.Fatalf("line %d is %d cells wide in a 60-column terminal: %q", i, width, line)
		}
	}

	pressCtrlL(app)

	if got := locale.Current(); got != locale.LanguageEnglish {
		t.Fatalf("second Ctrl+L: language = %s, want en", got)
	}
}

func TestCtrlLIgnoredOnFormScreens(t *testing.T) {
	app, prefPath := newTestApp(t)

	app.Update(models.NavigationMsg{Target: models.LLMProvidersScreen})
	if strings.Contains(app.View(), "Ctrl+L") {
		t.Fatal("language switch must not be offered outside the welcome screen and main menu")
	}

	pressCtrlL(app)

	if got := locale.Current(); got != locale.LanguageEnglish {
		t.Fatalf("language changed on a form screen: %s", got)
	}
	if _, ok := locale.LoadPreference(prefPath); ok {
		t.Fatal("preference saved although switching was not allowed")
	}
}
