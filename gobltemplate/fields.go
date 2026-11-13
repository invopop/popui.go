package gobltemplate

import (
	"strings"
	"sync"

	"github.com/invopop/icons"
	"github.com/invopop/popui.go/goblfields"
	"github.com/invopop/popui.go/props"
)

// AnyGroup is the name of the group with the variables that most schemas
// carry.
const AnyGroup = "Any"

// Format is the FieldPicker format for the tree that Fields returns.
const Format = "{{.%s}}"

// Fields returns the FieldPicker tree for templates about the given schemas:
//
//   - An Any group with the variables that more than half of the schemas
//     carry. Each one is described by the schemas that carry it, or "All
//     schemas".
//   - A group for each schema, with all of its variables, and then its
//     document under doc. The group's Description is the schema's ID.
//
// The groups are headings, so a pick inserts the same text from any group.
// The document leaves out the fields that a template cannot print with
// Format, such as the fields inside arrays.
func Fields(schemas ...Schema) []props.Field {
	carriers := map[string][]string{}
	for _, s := range schemas {
		for _, v := range Variables(s.ID) {
			carriers[v.Key] = append(carriers[v.Key], s.Label)
		}
	}

	var shared []props.Field
	for _, v := range variables {
		labels := carriers[v.Key]
		if len(labels) <= len(schemas)/2 {
			continue
		}
		field := variableField(v)
		if len(labels) == len(schemas) {
			field.Description = "All schemas"
		} else {
			field.Description = strings.Join(labels, ", ")
		}
		shared = append(shared, field)
	}

	fields := make([]props.Field, 0, len(schemas)+1)
	group := goblfields.Group(AnyGroup, shared)
	group.Icon = icons.Stack()
	fields = append(fields, group)

	for _, s := range schemas {
		var own []props.Field
		for _, v := range Variables(s.ID) {
			own = append(own, variableField(v))
		}
		if doc := documentFields(s.ID); doc != nil {
			nested := goblfields.Nest(docKey, doc)
			nested.Divider = true
			own = append(own, nested)
		}
		group := goblfields.Group(s.Label, own)
		group.Description = s.ID
		if s.Icon != nil {
			group.Icon = s.Icon()
		}
		fields = append(fields, group)
	}
	return fields
}

func variableField(v Variable) props.Field {
	return props.Field{
		Name:          v.Key,
		Path:          v.Key,
		Description:   v.Description,
		Type:          "string",
		AlwaysPresent: v.always,
	}
}

var (
	documentFieldsMu    sync.Mutex
	documentFieldsCache = map[string][]props.Field{}
)

// documentFields builds a schema's document tree once. A schema that the
// app's gobl release does not have gets no tree.
func documentFields(schema string) []props.Field {
	documentFieldsMu.Lock()
	defer documentFieldsMu.Unlock()
	if fields, ok := documentFieldsCache[schema]; ok {
		return fields
	}
	fields, err := goblfields.Build(schema, goblfields.Options{Descriptions: true, GoTemplate: true})
	if err != nil {
		fields = nil
	}
	documentFieldsCache[schema] = fields
	return fields
}
