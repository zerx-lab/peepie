package locale

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// identifierConstPrefixes are the only consts allowed in locale.go: internal
// identifiers that must never be translated. Any other string belongs in a
// `var` block so SetLanguage can translate it.
var identifierConstPrefixes = []string{"EmbedderProviderID"}

// sourceTexts parses locale.go and returns every text var with its English value.
func sourceTexts(t *testing.T) map[string]string {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), "locale.go", nil, 0)
	if err != nil {
		t.Fatalf("parse locale.go: %v", err)
	}

	texts := map[string]string{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || (gen.Tok != token.VAR && gen.Tok != token.CONST) {
			continue
		}
		for _, spec := range gen.Specs {
			vs := spec.(*ast.ValueSpec)
			for i, name := range vs.Names {
				if gen.Tok == token.CONST {
					if !slices.ContainsFunc(identifierConstPrefixes, func(p string) bool {
						return strings.HasPrefix(name.Name, p)
					}) {
						t.Errorf("%s is a const: user-visible text must be declared in a var block to be translatable", name.Name)
					}
					continue
				}
				if i >= len(vs.Values) {
					t.Fatalf("%s has no value", name.Name)
				}
				texts[name.Name] = stringValue(t, vs.Values[i])
			}
		}
	}

	return texts
}

// catalogTexts parses catalog_*.go and returns language → var name → translation.
func catalogTexts(t *testing.T) map[string]map[string]string {
	t.Helper()

	files, err := filepath.Glob("catalog_*.go")
	if err != nil {
		t.Fatal(err)
	}

	result := map[string]map[string]string{}
	for _, path := range files {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			if fn, ok := call.Fun.(*ast.Ident); !ok || fn.Name != "registerCatalog" || len(call.Args) != 2 {
				return true
			}

			langIdent, ok := call.Args[0].(*ast.Ident)
			if !ok {
				t.Fatalf("%s: registerCatalog language must be a Language constant", path)
			}
			lang := string(languageByIdent(t, langIdent.Name))
			if result[lang] == nil {
				result[lang] = map[string]string{}
			}

			entries := call.Args[1].(*ast.CompositeLit)
			for _, elt := range entries.Elts {
				kv := elt.(*ast.KeyValueExpr)
				name := kv.Key.(*ast.UnaryExpr).X.(*ast.Ident).Name
				if _, dup := result[lang][name]; dup {
					t.Errorf("%s: duplicate %s translation of %s", path, lang, name)
				}
				result[lang][name] = stringValue(t, kv.Value)
			}
			return false
		})
	}

	return result
}

func languageByIdent(t *testing.T, ident string) Language {
	t.Helper()
	switch ident {
	case "LanguageEnglish":
		return LanguageEnglish
	case "LanguageChineseSimplified":
		return LanguageChineseSimplified
	}
	t.Fatalf("unknown language constant %s", ident)
	return ""
}

// stringValue evaluates string literals and their `+` concatenations.
func stringValue(t *testing.T, expr ast.Expr) string {
	t.Helper()
	switch e := expr.(type) {
	case *ast.BasicLit:
		value, err := strconv.Unquote(e.Value)
		if err != nil || e.Kind != token.STRING {
			t.Fatalf("not a string literal: %s", e.Value)
		}
		return value
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			return stringValue(t, e.X) + stringValue(t, e.Y)
		}
	case *ast.ParenExpr:
		return stringValue(t, e.X)
	}
	t.Fatalf("unsupported expression %T: locale texts must be string literals", expr)
	return ""
}

var printfVerb = regexp.MustCompile(`%(?:\[\d+\])?[-+#0]*(?:\d+|\*)?(?:\.(?:\d+|\*))?[a-zA-Z%]`)

func TestEveryTextIsTranslated(t *testing.T) {
	source := sourceTexts(t)
	catalogs := catalogTexts(t)

	for _, lang := range SupportedLanguages {
		if lang == DefaultLanguage {
			continue
		}
		catalog := catalogs[string(lang)]

		for name, english := range source {
			translated, ok := catalog[name]
			if !ok {
				t.Errorf("%s: missing translation for %s = %q", lang, name, english)
				continue
			}
			if strings.TrimSpace(translated) == "" && strings.TrimSpace(english) != "" {
				t.Errorf("%s: empty translation for %s", lang, name)
			}
			if got, want := printfVerb.FindAllString(translated, -1), printfVerb.FindAllString(english, -1); !slices.Equal(got, want) {
				t.Errorf("%s: %s format verbs %v, source has %v", lang, name, got, want)
			}
		}
	}
}

func TestSetLanguageRoundTrip(t *testing.T) {
	t.Cleanup(func() { _ = SetLanguage(DefaultLanguage) })

	if err := SetLanguage(LanguageChineseSimplified); err != nil {
		t.Fatal(err)
	}
	if Current() != LanguageChineseSimplified {
		t.Fatalf("current = %s", Current())
	}
	for ptr, text := range catalogs[LanguageChineseSimplified] {
		if *ptr != text {
			t.Fatalf("var not switched: got %q want %q", *ptr, text)
		}
	}

	if err := SetLanguage(LanguageEnglish); err != nil {
		t.Fatal(err)
	}
	if UILoading != "Loading..." {
		t.Fatalf("English not restored: UILoading = %q", UILoading)
	}

	if err := SetLanguage("fr"); err == nil {
		t.Fatal("expected error for unsupported language")
	}
}

func TestNextLanguageCycles(t *testing.T) {
	t.Cleanup(func() { _ = SetLanguage(DefaultLanguage) })

	seen := map[Language]bool{}
	for range SupportedLanguages {
		seen[Current()] = true
		if err := SetLanguage(NextLanguage()); err != nil {
			t.Fatal(err)
		}
	}
	if len(seen) != len(SupportedLanguages) || Current() != DefaultLanguage {
		t.Fatalf("cycle visited %v, ended at %s", seen, Current())
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]Language{
		"en":          LanguageEnglish,
		"en_US.UTF-8": LanguageEnglish,
		"zh_CN.UTF-8": LanguageChineseSimplified,
		"zh-Hans-CN":  LanguageChineseSimplified,
		"zh_TW":       LanguageChineseSimplified,
		"ZH-cn":       LanguageChineseSimplified,
	}
	for tag, want := range cases {
		if got, ok := Normalize(tag); !ok || got != want {
			t.Errorf("Normalize(%q) = %q, %v; want %q", tag, got, ok, want)
		}
	}
	for _, tag := range []string{"", "de_DE", "C"} {
		if got, ok := Normalize(tag); ok {
			t.Errorf("Normalize(%q) = %q, want no match", tag, got)
		}
	}
}

func TestDetectFromEnv(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want Language
		ok   bool
	}{
		{"LANG chinese", map[string]string{"LANG": "zh_CN.UTF-8"}, LanguageChineseSimplified, true},
		{"LC_ALL overrides LANG", map[string]string{"LC_ALL": "en_US.UTF-8", "LANG": "zh_CN.UTF-8"}, LanguageEnglish, true},
		{"LC_MESSAGES overrides LANG", map[string]string{"LC_MESSAGES": "zh_CN.UTF-8", "LANG": "en_US.UTF-8"}, LanguageChineseSimplified, true},
		{"LANGUAGE list first supported", map[string]string{"LANGUAGE": "de:zh_CN:en", "LANG": "en_US.UTF-8"}, LanguageChineseSimplified, true},
		{"C locale ignores LANGUAGE", map[string]string{"LANGUAGE": "zh_CN", "LC_ALL": "C"}, LanguageEnglish, true},
		{"unsupported locale is English", map[string]string{"LANG": "de_DE.UTF-8"}, LanguageEnglish, true},
		{"nothing set", map[string]string{}, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := detectFromEnv(func(key string) string { return tc.env[key] })
			if got != tc.want || ok != tc.ok {
				t.Fatalf("got %q, %v; want %q, %v", got, ok, tc.want, tc.ok)
			}
		})
	}
}

func TestPreferenceRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "installer.json")

	if _, ok := LoadPreference(path); ok {
		t.Fatal("missing file must not yield a preference")
	}
	if err := SavePreference(path, LanguageChineseSimplified); err != nil {
		t.Fatal(err)
	}
	if got, ok := LoadPreference(path); !ok || got != LanguageChineseSimplified {
		t.Fatalf("got %q, %v", got, ok)
	}

	if err := os.WriteFile(path, []byte(`{"language":"xx"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, ok := LoadPreference(path); ok {
		t.Fatal("unsupported stored language must be ignored")
	}
	if err := SavePreference(path, "xx"); err == nil {
		t.Fatal("saving an unsupported language must fail")
	}
}
