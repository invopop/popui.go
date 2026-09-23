package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Field is one entry in a FieldSelector's tree. Trees are usually built from a
// GOBL schema with the goblfields package, but any nested data shape can be
// described by hand — a consumer can prepend its own entries (a job id, a
// rendered fault list) to a schema-derived tree so both are picked the same
// way.
type Field struct {
	// Name is the last segment of the path, without the array marker.
	Name string
	// Path is the full dotted path from the root, with [] marking each array
	// crossed on the way: lines[].item.name. It is what the picker emits
	// unless Value overrides it.
	Path string
	// Value overrides what selecting the field writes into the hidden input
	// and reports on the select event. Use it when the path the user reads is
	// not the string the consumer needs — a Go template variable, say.
	Value string
	// Title is the human label from the schema, shown after the path.
	Title string
	// Description is the schema's long description, shown as the row's
	// tooltip. Optional: goblfields leaves it out unless asked for.
	Description string
	// Type is the label shown at the end of the row: string, object,
	// string[], object[].
	Type string
	// Array marks a field that is itself an array. Its children describe one
	// item, so their paths carry the [] marker of this field.
	Array bool
	// Required marks a field its immediate parent requires.
	Required bool
	// AlwaysPresent marks a field required at every step from the root, and
	// so present on every document of the schema. Inside an array it reads
	// per item: lines[].item is marked when every line that exists carries an
	// item. The picker shows these with a green dot.
	AlwaysPresent bool
	// Children are the fields one level down, for a field that can be drilled
	// into.
	Children []Field
}

// FieldSelector Templ component props.
//
// The component renders a field-shaped trigger and a dropdown that either
// filters the whole tree or browses it one level at a time. Fields are
// rendered from the tree given here — the component makes no requests of its
// own — and the chosen path is submitted under Name via a hidden input.
type FieldSelector struct {
	ID    string
	Class string
	// Attributes are applied to the component root, the element carrying
	// x-modelable, so an Alpine x-model there two-way binds the selected
	// value.
	Attributes templ.Attributes
	// Name is the hidden input's name, for form submission.
	Name  string
	Label string
	// Tooltip shows a question mark icon after the label that reveals a
	// Tooltip card on hover. Only used when Label is set.
	Tooltip Tooltip
	// Placeholder is shown in the trigger when nothing is selected.
	Placeholder string
	// Value is the initially selected path. The dropdown opens where it
	// lives rather than at the root.
	Value string
	// Root labels the first breadcrumb in the dropdown, naming what the tree
	// describes: "bill/invoice", "envelope".
	Root string
	// Fields is the tree to pick from.
	Fields   []Field
	Disabled bool
	Required bool
	Error    Error
}

// GenerateID returns a new FieldSelector with either the existing ID or a
// newly generated one, so the label links to the trigger. It is designed to be
// used inline:
//
//	@popui.FieldSelector(props.FieldSelector{Name: "field"}.GenerateID())
func (f FieldSelector) GenerateID() FieldSelector {
	if f.ID != "" {
		return f
	}
	f.ID = fmt.Sprintf("field-selector-%06d", rand.Intn(100000))
	return f
}
