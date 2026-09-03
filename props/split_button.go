package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// SplitButton defines the properties for the SplitButton component: a main
// action button with an attached dropdown toggle exposing related actions.
type SplitButton struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Label is the text of the main action segment.
	Label string

	// Content replaces Label with any component as the main action's
	// content — e.g. an icon followed by a `<span>` label. Label is
	// ignored when set.
	Content templ.Component

	// Button configures the main action segment. Variant and Size also
	// apply to the dropdown toggle, and Disabled disables both segments.
	// The form fields (Type, Name, Value, Form) make the main action a
	// form button; Href makes it an anchor.
	Button Button

	// ToggleLabel is the accessible name (aria-label) of the dropdown
	// toggle. Default is "More actions".
	ToggleLabel string

	// RightAlign aligns the dropdown panel with the right edge of the
	// control instead of the left.
	RightAlign bool

	// DropUp opens the panel above the control instead of below it, for
	// controls that sit near the bottom of the viewport.
	DropUp bool

	// MenuClass adds classes to the dropdown panel.
	MenuClass string
}

// GenerateID generates a unique ID for the SplitButton if none is provided.
func (sb SplitButton) GenerateID() SplitButton {
	if sb.ID != "" {
		return sb
	}
	// generate a short random identifier
	sb.ID = fmt.Sprintf("split-button-%06d", rand.Intn(100000))
	return sb
}
