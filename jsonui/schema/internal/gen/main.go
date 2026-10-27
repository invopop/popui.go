// Command gen generates document.json from the jsonui registry, using the props structs'
// Go doc comments as descriptions. Run with
// `go generate ./jsonui/schema`.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	"github.com/invopop/jsonschema"
	"github.com/invopop/popui.go/jsonui"
	"github.com/invopop/popui.go/jsonui/schema"
)

const outPath = "jsonui/schema/document.json"

// JSON Schema type names.
const (
	typeObject = "object"
	typeString = "string"
	typeArray  = "array"
)

func main() {
	if err := generate(); err != nil {
		panic(err)
	}
}

// unsupported marks props that can't be expressed in JSON (templ.Component
// values such as icons) so they can be dropped from the output.
var unsupported = &jsonschema.Schema{Description: "unsupported"}

func generate() error {
	// Comment paths are resolved relative to the module root; go generate
	// runs in the schema package directory.
	if err := os.Chdir("../.."); err != nil {
		return err
	}

	componentType := reflect.TypeFor[templ.Component]()
	r := &jsonschema.Reflector{
		KeyNamer:                   lowerCamel,
		RequiredFromJSONSchemaTags: true,
		Mapper: func(t reflect.Type) *jsonschema.Schema {
			if t == componentType {
				return unsupported
			}
			return nil
		},
	}
	if err := r.AddGoComments("github.com/invopop/popui.go", "./props"); err != nil {
		return fmt.Errorf("reading comments: %w", err)
	}
	if err := r.AddGoComments("github.com/invopop/popui.go", "./jsonui"); err != nil {
		return fmt.Errorf("reading comments: %w", err)
	}
	nodeComment := func(field string) string {
		return r.CommentMap["github.com/invopop/popui.go/jsonui.Node."+field]
	}

	defs := jsonschema.Definitions{}
	merge := func(s *jsonschema.Schema) {
		for k, v := range s.Definitions {
			defs[k] = v
		}
	}

	doc := r.Reflect(jsonui.Document{})
	merge(doc)
	// Data is a json.RawMessage, which reflects as an untyped value.
	docDef := defs["Document"]
	if p, ok := docDef.Properties.Get("data"); ok {
		p.Type = typeObject
	}
	docDef.Properties.Set("template", &jsonschema.Schema{
		Type:        typeArray,
		Description: r.CommentMap["github.com/invopop/popui.go/jsonui.Document.Template"],
		Items:       &jsonschema.Schema{Ref: "#/$defs/Node"},
	})
	docDef.Required = []string{"template"}
	if p, ok := docDef.Properties.Get("handlers"); ok {
		p.AdditionalProperties = &jsonschema.Schema{Ref: "#/$defs/HandlerDefinition"}
	}
	delete(defs, "Node") // replaced by the per-type union below

	children := &jsonschema.Schema{
		Type:        typeArray,
		Description: nodeComment("Children"),
		Items:       &jsonschema.Schema{Ref: "#/$defs/Node"},
	}

	text := node("text")
	text.Description = "Plain text, optionally bound to an Alpine expression."
	text.Properties.Set("text", &jsonschema.Schema{Type: typeString, Description: nodeComment("Text")})
	text.Properties.Set("bind", &jsonschema.Schema{Type: typeString, Description: "Alpine expression shown with x-text. A plain data path such as \"customer.name\" is also resolved from object data at render time."})
	defs["node.text"] = text
	union := []*jsonschema.Schema{{Ref: "#/$defs/node.text"}}

	for _, c := range jsonui.Components() {
		ps := r.ReflectFromType(c.Props)
		merge(ps)
		n := node(c.Type)
		n.Description = fmt.Sprintf("PopUI %s component.", c.Props.Name())
		n.Properties.Set("props", &jsonschema.Schema{Ref: "#/$defs/" + c.Props.Name()})
		if !c.Leaf {
			n.Properties.Set("children", children)
			n.Properties.Set("bind", bindSchema(defs[c.Props.Name()]))
		}
		if c.Model {
			n.Properties.Set("model", &jsonschema.Schema{Type: typeString, Description: nodeComment("Model")})
		}
		n.Properties.Set("on", &jsonschema.Schema{
			Type:                 typeObject,
			Description:          nodeComment("On"),
			PropertyNames:        &jsonschema.Schema{Pattern: `^[a-z][a-z0-9-]*(\.[a-z0-9-]+)*$`},
			AdditionalProperties: &jsonschema.Schema{Ref: "#/$defs/Handler"},
		})
		defs["node."+c.Type] = n
		union = append(union, &jsonschema.Schema{Ref: "#/$defs/node." + c.Type})
	}
	addActions(defs)
	defs["Node"] = &jsonschema.Schema{
		Description: "A component in the template, identified by its type.",
		OneOf:       union,
	}

	for _, d := range defs {
		dropUnsupported(d)
	}
	for name, fields := range enums {
		for field, values := range fields {
			p, ok := defs[name].Properties.Get(field)
			if !ok {
				return fmt.Errorf("enum for unknown property %s.%s", name, field)
			}
			for _, v := range values {
				p.Enum = append(p.Enum, v)
			}
		}
	}

	out := &jsonschema.Schema{
		Version:     jsonschema.Version,
		ID:          schema.ID,
		Ref:         "#/$defs/Document",
		Title:       "PopUI JSON document",
		Definitions: defs,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s (%d definitions)\n", outPath, len(defs))
	return nil
}

// enums lists the accepted values of string props, keyed by props type and
// property, mirroring the constants and comments in the props package. These
// could move onto the props types as JSONSchemaExtend methods.
//
//nolint:goconst // prop names read more clearly inline than as constants
var enums = map[string]map[string][]string{
	"Button": {
		"variant": {"", "primary", "secondary", "danger", "transparent", "inverse"},
		"size":    {"", "sm", "lg", "icon", "icon-lg", "icon-sm", "icon-xs"},
		"type":    {"", "button", "submit", "reset"},
	},
	"ButtonGroup": {"align": {"", "left", "center", "right"}},
	"StatusBadge": {"status": {"success", "failed", "warning", "running"}},
	"TagStatus":   {"status": {"grey", "green", "yellow", "red", "orange", "blue", "purple", "olive", "teal", "crimson", "steelBlue", "empty"}},
	"Avatar":      {"color": {"", "accent", "success", "warning", "critical", "info", "document-xml", "document-pdf", "document-png"}, "size": {"", "lg"}},
	"Input":       {"size": {"", "sm", "lg"}},
	"Checkbox":    {"variant": {"", "switch"}},
	"Fieldset":    {"variant": {"", "card"}},
}

// param is an extra property accepted by a built-in action.
type param struct {
	name     string
	required bool
	schema   *jsonschema.Schema
}

// commonParams are handled by assets/jsonui.js for every action.
var commonParams = []param{
	{"into", false, &jsonschema.Schema{Type: typeString, Description: "Data path to store the action's result in."}},
	{"busy", false, &jsonschema.Schema{Type: typeString, Description: "Data path set to true while the action runs."}},
}

// builtinActions are implemented by assets/jsonui.js.
var builtinActions = []struct {
	name, desc string
	params     []param
}{
	{"alert", "Show a browser alert.", []param{
		{"message", true, &jsonschema.Schema{Type: typeString, Description: "Text to show; {{path}} placeholders are replaced with data values."}},
	}},
	{"set", "Set a data path to a value.", []param{
		{"path", true, &jsonschema.Schema{Type: typeString, Description: "Data path to change, e.g. \"customer.name\"."}},
		{"value", false, &jsonschema.Schema{Description: "New value; strings may use {{path}} placeholders."}},
	}},
	{"reset", "Restore the document data to its initial values.", nil},
	{"dispatch", "Dispatch a bubbling DOM CustomEvent for the host page to handle.", []param{
		{"event", true, &jsonschema.Schema{Type: typeString, Description: "Event name, e.g. \"customer:save\"."}},
		{"detail", false, &jsonschema.Schema{Description: "Event detail; defaults to a copy of the document data."}},
	}},
}

// actionRef refers to the Action definition.
const actionRef = "#/$defs/Action"

// addActions defines Action (one built-in or host-registered action) and
// Handler (an action or a list of actions).
func addActions(defs jsonschema.Definitions) {
	var options []*jsonschema.Schema
	var names []any
	for _, a := range builtinActions {
		s := &jsonschema.Schema{
			Type:                 typeObject,
			Description:          a.desc,
			Properties:           jsonschema.NewProperties(),
			Required:             []string{"action"},
			AdditionalProperties: jsonschema.FalseSchema,
		}
		s.Properties.Set("action", &jsonschema.Schema{Const: a.name})
		for _, p := range append(a.params, commonParams...) {
			s.Properties.Set(p.name, p.schema)
			if p.required {
				s.Required = append(s.Required, p.name)
			}
		}
		defs["action."+a.name] = s
		options = append(options, &jsonschema.Schema{Ref: "#/$defs/action." + a.name})
		names = append(names, a.name)
	}
	custom := &jsonschema.Schema{
		Type:        typeObject,
		Description: "An action registered by the host page with popui.actions.register(name, fn). Other properties are passed to it as params.",
		Properties:  jsonschema.NewProperties(),
		Required:    []string{"action"},
	}
	custom.Properties.Set("action", &jsonschema.Schema{Type: typeString, Not: &jsonschema.Schema{Enum: names}})
	options = append(options, custom)
	defs["Action"] = &jsonschema.Schema{Description: "An action to run in response to an event.", AnyOf: options}
	// Handlers defined on the document hold actions only; "on" values can
	// also name one of them.
	defs["HandlerDefinition"] = &jsonschema.Schema{
		AnyOf: []*jsonschema.Schema{
			{Ref: actionRef},
			{Type: typeArray, Items: &jsonschema.Schema{Ref: actionRef}},
		},
	}
	actionOrName := &jsonschema.Schema{AnyOf: []*jsonschema.Schema{
		{Ref: actionRef},
		{Type: typeString, Description: "The name of one of the document's handlers."},
	}}
	defs["Handler"] = &jsonschema.Schema{
		AnyOf: append(actionOrName.AnyOf, &jsonschema.Schema{Type: typeArray, Items: actionOrName}),
	}
}

// bindSchema describes a component's "bind" map: keys are the component's
// own prop names, values are data paths.
func bindSchema(props *jsonschema.Schema) *jsonschema.Schema {
	var names []any
	for pair := props.Properties.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value != unsupported {
			names = append(names, pair.Key)
		}
	}
	return &jsonschema.Schema{
		Type:                 typeObject,
		Description:          "Props to take from the document data, as prop name → data path, e.g. {\"value\": \"customer.name\"}. Requires object data; resolved at render time.",
		PropertyNames:        &jsonschema.Schema{Enum: names},
		AdditionalProperties: &jsonschema.Schema{Type: typeString},
	}
}

// node returns the base schema for a node of the given type.
func node(typ string) *jsonschema.Schema {
	s := &jsonschema.Schema{
		Type:                 typeObject,
		Properties:           jsonschema.NewProperties(),
		Required:             []string{"type"},
		AdditionalProperties: jsonschema.FalseSchema,
	}
	s.Properties.Set("type", &jsonschema.Schema{Const: typ})
	return s
}

func dropUnsupported(s *jsonschema.Schema) {
	if s.Properties == nil {
		return
	}
	for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value == unsupported || (pair.Value.Items != nil && pair.Value.Items == unsupported) {
			s.Properties.Delete(pair.Key)
		}
	}
}

// lowerCamel converts Go field names to the camelCase used in documents:
// "Label" → "label", "ID" → "id", "HTMLClass" → "htmlClass". Decoding is
// case-insensitive, so these map back onto the same fields.
func lowerCamel(s string) string {
	r := []rune(s)
	n := 0
	for n < len(r) && unicode.IsUpper(r[n]) {
		n++
	}
	if n > 1 && n < len(r) {
		n-- // keep the start of the next word capitalised
	}
	return strings.ToLower(string(r[:n])) + string(r[n:])
}
