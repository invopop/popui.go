package jsonui

import "embed"

// Assets holds jsonui.js, the browser runtime that executes the actions
// declared in a document's "on" handlers. Serve it alongside Alpine.js.
//
//go:embed assets/*
var Assets embed.FS
