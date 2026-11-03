//go:build js && wasm

// Command wasm exposes the jsonui renderer to the browser:
//
//	popuiRender(json, { actions: ["alert", "ui.notify"], trusted: false })
//	  → { html } or { error }
//
// Options are optional and default to strict mode with the built-in actions.
//
//	popuiValidate(schemas, data)
//	  → { errors: [{ path, message }] } or { error }
//
// validates panel data, as JSON, against a panel's schemas, as the
// JSON object of data paths to schemas jsonui.js finds on panels.
//
//	GOOS=js GOARCH=wasm go build -o popui.wasm ./jsonui/wasm
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"syscall/js"

	"github.com/invopop/popui.go/jsonui"
)

func main() {
	js.Global().Set("popuiRender", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) < 1 || args[0].Type() != js.TypeString {
			return map[string]any{"error": "popuiRender expects a JSON string"}
		}
		var opts jsonui.Options
		if len(args) > 1 && args[1].Type() == js.TypeObject {
			opts = options(args[1])
		}
		var buf bytes.Buffer
		if err := opts.Render(context.Background(), []byte(args[0].String()), &buf); err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"html": buf.String()}
	}))
	js.Global().Set("popuiValidate", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) < 2 || args[0].Type() != js.TypeString || args[1].Type() != js.TypeString {
			return map[string]any{"error": "popuiValidate expects the schemas and data as JSON strings"}
		}
		v, err := validators.get(args[0].String())
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		errs, err := v.ValidateJSON([]byte(args[1].String()))
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		list := make([]any, len(errs))
		for i, e := range errs {
			list[i] = map[string]any{"path": e.Path, "message": e.Message}
		}
		return map[string]any{"errors": list}
	}))
	names := jsonui.Names()
	list := make([]any, len(names))
	for i, name := range names {
		list[i] = name
	}
	js.Global().Set("popuiComponents", js.ValueOf(list))
	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("popui:ready"))

	select {} // keep the Go runtime alive for callbacks
}

// options reads jsonui.Options from a JS object.
func options(v js.Value) jsonui.Options {
	var o jsonui.Options
	if t := v.Get("trusted"); t.Type() == js.TypeBoolean {
		o.Trusted = t.Bool()
	}
	if a := v.Get("actions"); a.Type() == js.TypeObject && a.Get("length").Type() == js.TypeNumber {
		o.Actions = make([]string, a.Length())
		for i := range o.Actions {
			o.Actions[i] = a.Index(i).String()
		}
	}
	return o
}

// validators keeps the validators of the panels' schemas, as the same
// ones validate the data every time it changes.
var validators = validatorCache{byJSON: map[string]*jsonui.Validator{}}

type validatorCache struct {
	byJSON map[string]*jsonui.Validator
}

// maxValidators bounds the cache; panels rarely change schemas.
const maxValidators = 16

func (c *validatorCache) get(schemas string) (*jsonui.Validator, error) {
	if v, ok := c.byJSON[schemas]; ok {
		return v, nil
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(schemas), &raw); err != nil {
		return nil, err
	}
	v, err := jsonui.NewValidator(raw)
	if err != nil {
		return nil, err
	}
	if len(c.byJSON) >= maxValidators {
		clear(c.byJSON)
	}
	c.byJSON[schemas] = v
	return v, nil
}
