package jsonui

import (
	"encoding/json"
	"strings"
	"testing"
)

// pathTo is the path of the email config's recipient.
const pathTo = "config.to"

func mustValidator(t *testing.T, schemas string) *Validator {
	t.Helper()
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(schemas), &raw); err != nil {
		t.Fatal(err)
	}
	v, err := NewValidator(raw)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func validateJSON(t *testing.T, v *Validator, data string) map[string]string {
	t.Helper()
	errs, err := v.ValidateJSON([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range errs {
		out[e.Path] = e.Message
	}
	return out
}

func assertErrors(t *testing.T, got, want map[string]string) {
	t.Helper()
	for path, msg := range want {
		if got[path] != msg {
			t.Errorf("%s: got %q, want %q", path, got[path], msg)
		}
	}
	for path, msg := range got {
		if _, ok := want[path]; !ok {
			t.Errorf("unexpected error at %s: %q", path, msg)
		}
	}
}

const emailSchema = `{"config": {
	"type": "object",
	"required": ["connection_id", "to"],
	"additionalProperties": false,
	"properties": {
		"connection_id": {"type": "string", "minLength": 1, "title": "Connection"},
		"to": {"type": "string", "format": "email"},
		"mode": {"oneOf": [{"const": "html", "title": "HTML"}, {"const": "text", "title": "Plain text"}]},
		"priority": {"enum": ["low", "high"]},
		"retries": {"type": "integer", "minimum": 0, "maximum": 5},
		"cc": {"type": "array", "items": {"$ref": "#/$defs/address"}}
	},
	"$defs": {"address": {"type": "string", "format": "email"}}
}}`

func TestValidate(t *testing.T) {
	v := mustValidator(t, emailSchema)

	assertErrors(t, validateJSON(t, v, `{"config": {"connection_id": "c1", "to": "a@example.com"}}`), nil)

	// Each error is at the path of the value it's about.
	assertErrors(t, validateJSON(t, v, `{"config": {
		"to": "nope", "mode": "rtf", "priority": "", "retries": 9,
		"cc": ["b@example.com", "x"], "extra": true
	}}`), map[string]string{
		"config.connection_id": msgRequired,
		pathTo:                 msgEmail,
		"config.mode":          "Must be one of html, text",
		"config.priority":      "Required",
		"config.retries":       "Must be at most 5",
		"config.cc.1":          msgEmail,
		"config.extra":         "Not allowed",
	})

	assertErrors(t, validateJSON(t, v, `{"config": {"connection_id": "", "to": "a@example.com", "retries": 1.5}}`), map[string]string{
		"config.connection_id": msgRequired,
		"config.retries":       "Must be an integer",
	})

	// Missing data is validated as null.
	assertErrors(t, validateJSON(t, v, `{}`), map[string]string{"config": msgRequired})
}

func TestValidateErrorMessages(t *testing.T) {
	v := mustValidator(t, `{"config": {
		"type": "object",
		"required": ["to", "from", "subject"],
		"errorMessage": {"required": {"from": "Choose who it's from"}},
		"properties": {
			"to": {"type": "string", "format": "email", "errorMessage": "Enter the recipient's email address"},
			"from": {"type": "string"},
			"subject": {"type": "string", "maxLength": 5, "errorMessage": {"maxLength": "Keep it short", "required": "Add a subject"}}
		}
	}}`)
	assertErrors(t, validateJSON(t, v, `{"config": {}}`), map[string]string{
		pathTo:           "Enter the recipient's email address",
		"config.from":    "Choose who it's from",
		"config.subject": "Add a subject",
	})
	assertErrors(t, validateJSON(t, v, `{"config": {"to": "x", "from": "a", "subject": "too long"}}`), map[string]string{
		pathTo:           "Enter the recipient's email address",
		"config.subject": "Keep it short",
	})
}

func TestNewValidatorErrors(t *testing.T) {
	for schemas, want := range map[string]string{
		`{"config": {"$ref": "https://example.com/schema.json"}}`: "may only refer to themselves",
		`{"config": {"type": "nope"}}`:                            "type",
		`{"__proto__.x": {}}`:                                     errNotDataPath,
		`{"a b": {}}`:                                             errNotDataPath,
	} {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(schemas), &raw); err != nil {
			t.Fatal(err)
		}
		_, err := NewValidator(raw)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: got %v, want %q", schemas, err, want)
		}
	}
	// Documents report them too.
	_, err := render(t, Options{}, `{"data": {}, "schemas": {"config": {"$ref": "file:///etc/passwd"}}, "components": []}`)
	if err == nil || !strings.Contains(err.Error(), "schemas config") {
		t.Errorf("got %v", err)
	}
}

// TestFields checks that form controls show the error for their model, and
// take what they don't set from their schema.
func TestFields(t *testing.T) {
	html := mustRender(t, `{
		"data": {"config": {"connection_id": "", "to": "", "mode": "text", "priority": ""}, "items": [{"v": ""}]},
		"schemas": `+emailSchema+`,
		"components": [
			{"is": "input", "model": "config.connection_id"},
			{"is": "input", "model": "config.to", "props": {"label": "Recipient"}},
			{"is": "select", "model": "config.mode"},
			{"is": "select", "model": "config.priority"},
			{"is": "input", "each": "items", "model": "item.v"}
		]
	}`)
	assertContains(t, html,
		`x-jsonui-errors="{&#34;config&#34;:{&#34;type&#34;:&#34;object&#34;`,
		`<div data-jsonui-field="config.to">`,
		`:aria-invalid="!!$errors.shown(&#39;config.to&#39;)"`,
		`:class="$errors.shown(&#39;config.to&#39;) ? &#39;!text-foreground-critical`,
		`x-show="$errors.shown(&#39;config.to&#39;)" x-text="$errors.shown(&#39;config.to&#39;)"`,
		`<template x-for="e in $errors.unplaced()" :key="e.path">`,
		// The label comes from the title unless the node sets one.
		`>Connection</label>`,
		`>Recipient</label>`,
		// Options come from the choices, with the current value.
		`<option value="html">HTML</option>`,
		`<option value="text" selected>Plain text</option>`,
		`<option value="" selected></option>`,
		`<option value="low">low</option>`,
	)
	// Fields in a list have no path of their own.
	if strings.Contains(html, `data-jsonui-field="item.v"`) {
		t.Errorf("list item is a field in:\n%s", html)
	}
	// Selects with an option for their value don't get another.
	if strings.Contains(html, `<option value="text" selected>text</option>`) {
		t.Errorf("extra option in:\n%s", html)
	}
}

// TestFieldsWithoutSchemas checks that controls show errors actions report
// without schemas, and that documents without controls are unchanged.
func TestFieldsWithoutSchemas(t *testing.T) {
	html := mustRender(t, `{"data": {"name": ""}, "components": [{"is": "input", "model": "name"}]}`)
	assertContains(t, html, `x-jsonui-errors=""`, `<div data-jsonui-field="name">`)

	html = mustRender(t, `{"data": {"name": "x"}, "components": [{"is": "text", "bind": "name"}]}`)
	if strings.Contains(html, "jsonui-errors") || strings.Contains(html, "$errors") {
		t.Errorf("errors tracked without fields in:\n%s", html)
	}
}

func TestValidateAction(t *testing.T) {
	html := mustRender(t, `{
		"data": {"config": {}},
		"schemas": {"config": {"type": "object"}},
		"components": [{"is": "button", "on": {"click": [{"action": "validate"}, {"action": "dispatch", "event": "save"}]}}]
	}`)
	assertContains(t, html, `{&#34;action&#34;:&#34;validate&#34;}`)
}
