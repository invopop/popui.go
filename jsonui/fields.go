package jsonui

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"

	"github.com/invopop/popui.go/classes"
)

// Form controls with a "model" are fields: they show the error for their
// model's path, which jsonui.js finds by validating the data against the
// document's schemas, or which actions report, such as a host saving the
// data. Errors from schemas are only shown once the field has been changed
// or left, or after the "validate" action, so a new form isn't covered in
// them.

// errorIDName is the name of the Alpine x-id scope of a field, and errorID
// the expression of its error's id, unique to the field on the page.
const (
	errorIDName = "jsonui-error"
	errorID     = "$id('" + errorIDName + "')"
)

// fieldAttrs adds the bindings that mark a control as invalid, and described
// by its error, while its path has an error shown.
func (b *docBuilder) fieldAttrs(path string, attrs map[string]any) {
	b.fields++
	shown := shownError(path)
	if _, ok := attrs[":class"]; !ok { // only trusted documents set their own
		attrs[":class"] = shown + " ? '" + classes.FormFieldError() + "' : ''"
	}
	attrs[":aria-invalid"] = "!!" + shown
	_, own := attrs["aria-describedby"]
	if _, bound := attrs[":aria-describedby"]; !own && !bound {
		attrs[":aria-describedby"] = shown + " ? " + errorID + " : null"
	}
}

// schemaProps fills in the props a control's schema provides that the node
// hasn't set: its label from "title", its description from "description",
// and a select's options from its choices.
func (b *docBuilder) schemaProps(n *Node, props map[string]json.RawMessage) error {
	s := b.validator.schemaAt(n.Model)
	if s == nil {
		return nil
	}
	typ := registry[n.Component].props
	for prop, key := range map[string]string{"Label": "title", "Description": "description"} {
		text, ok := s[key].(string)
		if !ok || text == "" || !hasField(typ, prop) || hasKey(props, prop) {
			continue
		}
		var err error
		if props[strings.ToLower(prop)], err = json.Marshal(text); err != nil {
			return err
		}
	}
	if n.Component != selectType || len(n.Children) > 0 || hasKey(props, "options") || multiple(props) {
		return nil
	}
	opts, ok := schemaOptions(s)
	if !ok {
		return nil
	}
	// The select must have an option for its value, or it would show the
	// first one while the data has another.
	value, _ := lookup(b.data, n.Model) //nolint:errcheck // model checks it
	current := display(value)
	if !slices.ContainsFunc(opts, func(o Option) bool { return o.Value == current }) {
		opts = append([]Option{{Value: current, Label: current}}, opts...)
	}
	list := make([]map[string]string, len(opts))
	for i, o := range opts {
		list[i] = map[string]string{"value": o.Value, "label": o.Label}
	}
	var err error
	props["options"], err = json.Marshal(list)
	return err
}

// hasField reports whether a props struct has a field, which JSON matches
// case-insensitively.
func hasField(t reflect.Type, name string) bool {
	if t == nil || t.Kind() != reflect.Struct {
		return false
	}
	_, ok := t.FieldByName(name)
	return ok
}

// multiple reports whether a select's props make it a multiple select.
func multiple(props map[string]json.RawMessage) bool {
	for k, v := range props {
		if strings.EqualFold(k, "multiple") {
			var on bool
			_ = json.Unmarshal(v, &on) //nolint:errcheck // decoding the props reports it
			return on
		}
	}
	return false
}

// hasKey reports whether props set a prop, matching names as decoding does.
func hasKey(props map[string]json.RawMessage, name string) bool {
	for k, v := range props {
		if strings.EqualFold(k, name) && string(v) != jsonNull {
			return true
		}
	}
	return false
}
