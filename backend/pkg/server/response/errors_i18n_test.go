package response

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
)

// The web UI translates REST errors by code (frontend/src/i18n/locales/<lng>/apiErrors.json,
// see frontend/docs/i18n.md). Adding or changing an error here requires updating those
// catalogs: English mirrors the message verbatim, every other language needs the code.
const apiErrorsCatalogDir = "../../../../frontend/src/i18n/locales"

func declaredHttpErrors(t *testing.T) map[string]string {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), "errors.go", nil, 0)
	if err != nil {
		t.Fatalf("parse errors.go: %v", err)
	}

	errs := map[string]string{}
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if fn, ok := call.Fun.(*ast.Ident); !ok || fn.Name != "NewHttpError" || len(call.Args) != 3 {
			return true
		}

		code, codeErr := strconv.Unquote(call.Args[1].(*ast.BasicLit).Value)
		msg, msgErr := strconv.Unquote(call.Args[2].(*ast.BasicLit).Value)
		if codeErr != nil || msgErr != nil {
			t.Fatalf("NewHttpError arguments must be string literals")
		}
		errs[code] = msg
		return false
	})

	return errs
}

func loadCatalog(t *testing.T, language string) map[string]string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(apiErrorsCatalogDir, language, "apiErrors.json"))
	if err != nil {
		t.Fatalf("read %s apiErrors catalog: %v", language, err)
	}

	catalog := map[string]string{}
	if err := json.Unmarshal(data, &catalog); err != nil {
		t.Fatalf("parse %s apiErrors catalog: %v", language, err)
	}

	return catalog
}

func TestHttpErrorsHaveWebTranslations(t *testing.T) {
	declared := declaredHttpErrors(t)
	if len(declared) == 0 {
		t.Fatal("no NewHttpError declarations found")
	}

	if english := loadCatalog(t, "en"); !maps.Equal(english, declared) {
		t.Errorf("en/apiErrors.json must mirror errors.go exactly\n got: %v\nwant: %v", english, declared)
	}

	entries, err := os.ReadDir(apiErrorsCatalogDir)
	if err != nil {
		t.Fatal(err)
	}
	codes := slices.Sorted(maps.Keys(declared))
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "en" {
			continue
		}
		catalog := loadCatalog(t, entry.Name())
		if got := slices.Sorted(maps.Keys(catalog)); !slices.Equal(got, codes) {
			t.Errorf("%s/apiErrors.json codes differ from errors.go\n got: %v\nwant: %v", entry.Name(), got, codes)
		}
		for code, text := range catalog {
			if text == "" {
				t.Errorf("%s/apiErrors.json: empty translation for %s", entry.Name(), code)
			}
		}
	}
}
