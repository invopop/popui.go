package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Field is one entry in a FieldPicker's tree. Trees are usually built from a
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
	// Description is the schema's long description, shown after the path and
	// as the row's tooltip. It truncates before the path does, and is never
	// searched: filtering matches the path alone. Optional — goblfields
	// leaves it out unless asked for, and a row without one shows just its
	// path.
	Description string
	// Type is the JSON Schema type shown at the end of the row: string,
	// object, array, integer, number, boolean.
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
	// Group marks a heading over the fields beneath it rather than a field
	// of the data: it opens like an object and its Name is part of the path
	// the user reads and filters by, but it is left out of what its
	// descendants emit. It lets one tree hold several — the variables of
	// each document schema a message might be about, say — without the
	// grouping changing the values picked. goblfields.Group builds one.
	Group bool
	// Children are the fields one level down, for a field that can be drilled
	// into.
	Children []Field
}

// FieldPicker Templ component props.
//
// The component is a panel that picks one field of a tree — filtering the
// whole tree by path, or browsing it a level at a time — hung off whichever
// trigger suits: its own button by default, or anything placed in its
// children slot (a Button, an Input showing the value, …). A supplied trigger
// is rendered as given, so give it the picker's ID when there is a Label, so
// the label points at it, and the popup attributes the default button
// carries — aria-haspopup="listbox", :aria-expanded="open" and
// :aria-controls="$id('field-picker') + '-listbox'" — which resolve in the
// picker's Alpine scope. A pick is wrapped
// in Format, inserted at the caret of the Target element when there is one,
// submitted under Name via a hidden input, and announced with a
// `field-select` event.
type FieldPicker struct {
	ID string
	// Class is applied to the default trigger button.
	Class string
	// Attributes are applied to the component root, the element carrying
	// x-modelable, so an Alpine x-model there two-way binds the value.
	Attributes templ.Attributes
	// Name is the hidden input's name, for form submission.
	Name string
	// Label renders above the trigger.
	Label string
	// Tooltip shows a question mark icon after the label that reveals a
	// Tooltip card on hover. Only used when Label is set.
	Tooltip Tooltip
	// TriggerLabel is the default trigger button's text. Defaults to
	// "Insert field". Ignored when the children slot supplies the trigger.
	TriggerLabel string
	// Value is the initial value. The panel opens where it lives rather
	// than at the root.
	Value string
	// Root labels the first breadcrumb in the panel, naming what the tree
	// describes: "bill/invoice", "envelope".
	Root string
	// Fields is the tree to pick from.
	Fields []Field
	// Format wraps a picked path, %s standing for the path: "{{.doc.%s}}"
	// turns supplier.name into {{.doc.supplier.name}}. Defaults to "%s". A
	// field with a Value of its own is emitted as is.
	Format string
	// Target is a CSS selector for an input, textarea or contenteditable
	// element. When set, a pick is inserted there at the caret, replacing
	// any selection, with the caret left after it. When the target is a
	// template editor (a Contenteditable with a VariableFormat), typing the
	// format's opening characters there — {{ — opens the picker at the
	// caret with the filter focused, and the pick replaces them.
	Target string
	// ScalarsOnly makes objects and arrays browse-only: they can be opened
	// but not picked, and one with nothing underneath is left out.
	// For a template, where a picked object would print as a Go map dump.
	// Off by default, since a picker may be choosing a field for its own
	// sake rather than for its value.
	ScalarsOnly bool
	Disabled    bool
	Required    bool
	Error       Error
}

// GenerateID returns a new FieldPicker with either the existing ID or a
// newly generated one, so the label links to the trigger. It is designed to be
// used inline:
//
//	@popui.FieldPicker(props.FieldPicker{Name: "field"}.GenerateID())
func (f FieldPicker) GenerateID() FieldPicker {
	if f.ID != "" {
		return f
	}
	f.ID = fmt.Sprintf("field-picker-%06d", rand.Intn(100000))
	return f
}
