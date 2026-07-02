package props

import "github.com/a-h/templ"

const (
	// FieldsetVariantCard represents a card-style fieldset.
	FieldsetVariantCard string = "card"
)

// Fieldset properties to configure a simple fieldset.
type Fieldset struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Legend     string
	Variant    string
}

// FieldsetCard Templ component props. FieldsetCard groups a set of form fields
// inside a tinted card body, with an optional title and a secondary description
// rendered above it. Use it for settings sections and grouped form blocks.
type FieldsetCard struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Title renders a bold heading above the card body.
	Title string
	// Description renders next to the Title, separated by a middot, in a muted color.
	Description string
}
