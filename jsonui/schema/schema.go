// Package schema provides the JSON Schema for jsonui documents, generated
// from the jsonui component registry and the popui props structs.
package schema

import _ "embed"

//go:generate go run ./internal/gen

// ID is the canonical location of the document schema; documents reference
// it with "$schema" so editors can offer completion and validation.
const ID = "https://raw.githubusercontent.com/invopop/popui.go/main/jsonui/schema/document.json"

// Document is the generated JSON Schema for a jsonui document.
//
//go:embed document.json
var Document []byte
