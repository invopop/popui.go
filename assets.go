// Package popui provides the Templ components and support functions required to use
// POPUI in a Go application.
package popui

import (
	"crypto/md5"
	"embed"
	"fmt"
	"io"
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

// Versioned provides the versioned path for the given path assuming the file
// exists in the Assets embedded filestyem.
func Versioned(file ...string) string {
	p := path.Join(file...)
	if v, ok := versionCache[p]; ok {
		return p + "?v=" + v
	}
	f, err := Assets.Open(p)
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

	versionCache[p] = v

	return p + "?v=" + v
}
