// Package schema provides the JSON Schema for jsonui panels, generated
// from the jsonui component registry and the popui props structs.
package schema

import _ "embed"

//go:generate go run ./internal/gen

// ID is the canonical location of the panel schema; panels reference
// it with "$schema" so editors can offer completion and validation. It names
// the major version, where each release publishes its schema, so it stays the
// same across minor versions; change it when the major version changes.
const ID = "https://assets.invopop.com/popui/v0/panel.json"

// Panel is the generated JSON Schema for a jsonui panel.
//
//go:embed panel.json
var Panel []byte
