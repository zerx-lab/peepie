# Installer TUI internationalization (i18n)

The installer (`backend/cmd/installer`) ships in **English (`en`, source)** and **Simplified Chinese (`zh-CN`)**. All user-visible text lives in `backend/cmd/installer/wizard/locale`; tests fail when a string is not translated.

## How the language is chosen

| Priority | Source |
|---|---|
| 1 | `-l <tag>` flag for this run (`en`, `zh-CN`, `zh`, `zh_CN.UTF-8` …), not saved |
| 2 | Choice saved with **Ctrl+L** in `<user config dir>/pentagi/installer.json` (`~/.config/pentagi/installer.json` on Linux, `~/Library/Application Support/pentagi/installer.json` on macOS) |
| 3 | System locale: `LANGUAGE` list (ignored when the locale is `C`/`POSIX`), then `LC_ALL` > `LC_MESSAGES` > `LANG`. A set but unsupported locale means English |
| 4 | Only when no locale variable is set: macOS `AppleLanguages`, Windows preferred UI languages |
| 5 | English |

Any `zh*` tag (including `zh_TW`, `zh-Hant`) maps to `zh-CN`.

**Ctrl+L** switches to the next language on the Welcome screen and the Main Menu; the footer shows the target language in its own script (`Ctrl+L: Switch to 简体中文` / `Ctrl+L: 切换到 English`). It is hidden and ignored while an operation runs (install, apply changes, password reset …) because switching recreates every screen. Command-line usage, startup messages and fatal errors use the same catalog.

## How it works

- `locale.go` declares every text as a package-level **`var`** holding the English source. Only internal identifiers stay `const` (allow-listed in `locale_test.go`).
- `catalog_zh_cn_*.go` files register translations keyed by the var's address:

  ```go
  func init() {
      registerCatalog(LanguageChineseSimplified, map[*string]string{
          &UIStatistics: "统计",
      })
  }
  ```

  A renamed or removed var is a compile error in every catalog.
- `SetLanguage` restores the English snapshot and applies the target catalog; `App.switchLanguage` then rebuilds the screen registry so texts cached in constructors (list options, form fields) are rebuilt.
- `i18n.go`: languages, detection, `Normalize`, `NextLanguage`; `preference.go`: saved choice; `detect_*.go`: OS fallbacks.

## Rules for new code

1. Never put user-visible text in `models/`, `controller/`, `main.go` or elsewhere — add a `var` to the matching block in `locale.go` and reference `locale.X`.
2. Add the translation to the `catalog_zh_cn_*.go` file for that area in the same change; `go test ./cmd/installer/wizard/locale/` (`TestEveryTextIsTranslated`) fails otherwise. An identical translation is fine for texts that are the same in both languages (model names, URLs).
3. Format strings: translations must keep exactly the same printf verbs in the same order (the test checks it). Use `errors.New(locale.X)` rather than `fmt.Errorf(locale.X)` for texts without verbs (`go vet` rejects non-constant format strings).
4. Read `locale.X` when building or rendering a screen, never in package-level variables (they would freeze the language active at init). `TestEveryScreenFitsTheWindowInEveryLanguage` (registry package) renders every screen in every language at 80/120/200 columns and fails when a screen ignores the language or overflows the window.
5. Keep one sentence per var instead of concatenating fragments; word order differs between languages.
6. Mind width: CJK characters are two cells wide (lipgloss accounts for it); keep key hints and labels short.
7. Adding a language: extend `SupportedLanguages`, `nativeNames`, `normalizeTable` (if needed) and `languageByIdent` in the test, then add a full set of catalog files.

## Not translated

Environment variable names and values, commands, Docker/compose output and terminal output, logs (`logger.*`), the EULA (`EULA.md`, legal text), product/provider/model names other than the user-facing Peepie name, URLs, keyboard key names. Keep upstream `pentagi` technical identifiers and external service names unchanged.

Chinese style and terminology follow the glossary in [`frontend/docs/i18n.md`](../../frontend/docs/i18n.md#chinese-style-guide-and-glossary) so the web UI and the installer use the same words.
