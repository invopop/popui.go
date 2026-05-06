// Package popui provides the Templ components and support functions required to use
// POPUI in a Go application.
package popui

import (
	"crypto/md5"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"
)

// We use the tailwindcss CLI directly here for performance and to avoid
// needing to install the NodeJS toolchain.
//go:generate tailwindcss -i ./styles.css -o ./assets/popui.css --minify

const (
	// AssetPath determines where assets should be served from.
	AssetPath = "/_popui"
)

// Assets provides access to the CSS generated style
//
//go:embed assets/*
var Assets embed.FS

var versionCache = map[string]string{}

// Versioned will find the file inside the provided filesystem and
// add a version hash as a query parameter so that when the asset is loaded
// by the browser it'll always use the latest version.
//
// For example, to load a JS file inside a Templ component:
//
//	<script src={ popui.Versioned(assets.Content, "scripts", "app.js") }></script>
//
// Where `assets.Content` is the source of the file and `"scripts", "app.js"`
// identify the file's location. The above example might produce:
//
//	<script src="scripts/app.js?v=a1b2c3d4"></script>
//
// A simple version cache is used, and will only be renewed upon reloading
// the application. Paths must be unique for this to work correctly.
func Versioned(content fs.FS, file ...string) string {
	p := path.Join(file...)
	if v, ok := versionCache[p]; ok {
		return v
	}
	f, err := content.Open(p)
	if err != nil {
		return p
	}
	defer f.Close() //nolint:errcheck

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return p
	}

	v := fmt.Sprintf("%x", h.Sum(nil))[0:8]
	vp := p + "?v=" + v
	versionCache[p] = vp
	return vp
}
