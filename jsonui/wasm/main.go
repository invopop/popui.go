//go:build js && wasm

// Command wasm exposes the jsonui renderer to the browser:
//
//	popuiRender(json, { actions: ["alert", "ui.notify"], trusted: false })
//	  → { html } or { error }
//
// Options are optional and default to strict mode with the built-in actions.
//
//	GOOS=js GOARCH=wasm go build -o popui.wasm ./jsonui/wasm
package main

import (
	"bytes"
	"context"
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
	types := jsonui.Types()
	list := make([]any, len(types))
	for i, t := range types {
		list[i] = t
	}
	js.Global().Set("popuiTypes", js.ValueOf(list))
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
