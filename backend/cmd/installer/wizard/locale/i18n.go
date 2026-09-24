package locale

import (
	"fmt"
	"strings"
	"sync"
)

// Language is a BCP 47 tag of a supported installer UI language.
type Language string

const (
	LanguageEnglish           Language = "en"
	LanguageChineseSimplified Language = "zh-CN"

	// DefaultLanguage is the source language of locale.go and the fallback.
	DefaultLanguage = LanguageEnglish
)

// SupportedLanguages lists UI languages in switching order (Ctrl+L cycles through it).
// Adding a language: append it here, give it a native name below, extend
// normalizeTable if needed, and add catalog files registering every text var.
var SupportedLanguages = []Language{LanguageEnglish, LanguageChineseSimplified}

// nativeNames are shown in the language's own script so users can find theirs
// whatever the current UI language is; never translate them.
var nativeNames = map[Language]string{
	LanguageEnglish:           "English",
	LanguageChineseSimplified: "简体中文",
}

// NativeName returns the language name written in that language.
func (l Language) NativeName() string {
	if name, ok := nativeNames[l]; ok {
		return name
	}
	return string(l)
}

var (
	mu        sync.Mutex
	current   = DefaultLanguage
	catalogs  = map[Language]map[*string]string{}
	source    map[*string]string
	snapshotO sync.Once
)

// registerCatalog adds translations for lang; called from init() of catalog_*.go
// files. Keys are addresses of the text vars in locale.go, so a renamed or removed
// var is a compile error in every catalog.
func registerCatalog(lang Language, entries map[*string]string) {
	catalog, ok := catalogs[lang]
	if !ok {
		catalog = make(map[*string]string, len(entries))
		catalogs[lang] = catalog
	}
	for ptr, text := range entries {
		if _, dup := catalog[ptr]; dup {
			panic(fmt.Sprintf("locale: duplicate %s translation for %q", lang, *ptr))
		}
		catalog[ptr] = text
	}
}

// Current returns the active UI language.
func Current() Language {
	mu.Lock()
	defer mu.Unlock()
	return current
}

// SetLanguage switches every text var to lang. Screens that cached strings at
// construction must be rebuilt afterwards (the wizard recreates its registry).
// Must not race with rendering: call it before the TUI starts or from the
// bubbletea update loop.
func SetLanguage(lang Language) error {
	if !IsSupported(lang) {
		return fmt.Errorf("unsupported language %q (supported: %s)", lang, SupportedLanguageTags())
	}

	mu.Lock()
	defer mu.Unlock()

	// capture the English source once, before the first switch mutates anything
	snapshotO.Do(func() {
		source = make(map[*string]string)
		for _, catalog := range catalogs {
			for ptr := range catalog {
				source[ptr] = *ptr
			}
		}
	})

	for ptr, text := range source {
		*ptr = text
	}
	for ptr, text := range catalogs[lang] {
		*ptr = text
	}
	current = lang

	return nil
}

// NextLanguage returns the language after the current one in SupportedLanguages.
func NextLanguage() Language {
	cur := Current()
	for i, lang := range SupportedLanguages {
		if lang == cur {
			return SupportedLanguages[(i+1)%len(SupportedLanguages)]
		}
	}
	return DefaultLanguage
}

// IsSupported reports whether lang is one of SupportedLanguages.
func IsSupported(lang Language) bool {
	for _, supported := range SupportedLanguages {
		if supported == lang {
			return true
		}
	}
	return false
}

// normalizeTable maps a lowercased primary subtag to a supported language when
// the full tag has no exact match. Every Chinese variant (zh_TW, zh-Hant…) maps
// to Simplified Chinese, the only Chinese translation shipped.
var normalizeTable = map[string]Language{
	"en": LanguageEnglish,
	"zh": LanguageChineseSimplified,
}

// Normalize maps a locale tag such as "zh_CN.UTF-8", "en-GB" or "zh-Hans" onto a
// supported language.
func Normalize(tag string) (Language, bool) {
	tag = strings.TrimSpace(tag)
	if i := strings.IndexAny(tag, ".@"); i >= 0 {
		tag = tag[:i]
	}
	tag = strings.ToLower(strings.ReplaceAll(tag, "_", "-"))
	if tag == "" {
		return "", false
	}

	for _, lang := range SupportedLanguages {
		if strings.ToLower(string(lang)) == tag {
			return lang, true
		}
	}

	primary, _, _ := strings.Cut(tag, "-")
	lang, ok := normalizeTable[primary]
	return lang, ok
}

// DetectSystemLanguage derives the UI language from the environment following
// POSIX/GNU gettext precedence (LANGUAGE list unless the locale is C/POSIX, then
// LC_ALL > LC_MESSAGES > LANG), then the OS user preference (macOS, Windows),
// falling back to English.
func DetectSystemLanguage(getenv func(string) string) Language {
	if lang, ok := detectFromEnv(getenv); ok {
		return lang
	}
	for _, tag := range osPreferredLanguages() {
		if lang, ok := Normalize(tag); ok {
			return lang
		}
	}
	return DefaultLanguage
}

// detectFromEnv reports ok=false only when no locale variable is set at all, so
// an explicit but unsupported locale (e.g. de_DE) yields English instead of
// consulting the OS.
func detectFromEnv(getenv func(string) string) (Language, bool) {
	locale := ""
	for _, name := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := strings.TrimSpace(getenv(name)); value != "" {
			locale = value
			break
		}
	}

	base, _, _ := strings.Cut(locale, ".")
	if base == "C" || base == "POSIX" {
		return LanguageEnglish, true
	}

	for _, tag := range strings.Split(getenv("LANGUAGE"), ":") {
		if lang, ok := Normalize(tag); ok {
			return lang, true
		}
	}

	if locale == "" {
		return "", false
	}
	if lang, ok := Normalize(locale); ok {
		return lang, true
	}
	return DefaultLanguage, true
}

// SupportedLanguageTags lists the accepted language tags, e.g. "en, zh-CN".
func SupportedLanguageTags() string {
	tags := make([]string, len(SupportedLanguages))
	for i, lang := range SupportedLanguages {
		tags[i] = string(lang)
	}
	return strings.Join(tags, ", ")
}
