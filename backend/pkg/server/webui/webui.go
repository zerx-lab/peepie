// Package webui holds the frontend build embedded into the backend binary.
//
// `task build` copies frontend/dist into ./dist before compiling. Without that
// step only the placeholder .gitkeep is embedded, Embedded reports false, and
// the server falls back to serving STATIC_DIR from disk (Docker image layout).
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// Embedded returns the embedded frontend build rooted at its index.html, and
// whether a build was embedded at all.
func Embedded() (fs.FS, bool) {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, false
	}
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		return nil, false
	}

	return sub, true
}
