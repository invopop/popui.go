package props

import "github.com/a-h/templ"

// Stat sizes.
const (
	StatSizeMedium = "md"
	StatSizeLarge  = "lg"
)

// Stats defines the props for the Stats component.
type Stats struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Align is "end" to keep the stats at the end of the group, with their
	// labels and figures aligned to it, also on the rows they wrap to,
	// while a large stat leads at the start, as in a document's totals.
	// Default is "start".
	Align string `jsonschema:"enum=,enum=start,enum=end"`
}

// Stat defines the props for the Stat component.
type Stat struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label says what the figure is, shown above it.
	Label string
	// Value is the figure. Without one, the stat's children are shown as
	// the value.
	Value string
	// Size is "lg" for the headline figure of a group, which leads its
	// row, or "md" for the others. Default is "md".
	Size string `jsonschema:"enum=,enum=md,enum=lg"`
}
