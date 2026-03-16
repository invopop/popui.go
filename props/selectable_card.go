package props

import "github.com/a-h/templ"

// SelectableCard defines the properties for a selectable card component.
// Can render as a radio-based card (with Name for grouping) or a standalone
// clickable card. Supports title, description, and trailing children content.
type SelectableCard struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Label is the primary text shown in the card.
	Label string
	// Description is secondary text shown below the label.
	Description string
	// Value is the value attribute for the radio input.
	Value string
	// Name groups radio-based selectable cards together.
	Name string
	// Checked marks this card as initially selected.
	Checked bool
	// Disabled prevents interaction with the card.
	Disabled bool
}
