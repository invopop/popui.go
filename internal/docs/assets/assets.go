// Package assets holds the embedded assets for the UI.
package assets

import (
	"embed"

	popui "github.com/invopop/popui.go"
)

//go:embed scripts/* prism-popui.css

// Content holds the embedded assets.
var Content embed.FS

// Versioned provides the versioned path for the given path assuming the file
// exists in the Content.
func Versioned(file ...string) string {
	return popui.Versioned(Content, file...)
}
