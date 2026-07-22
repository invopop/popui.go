package props

import "github.com/a-h/templ"

const (
	// TooltipPositionTop shows the card above the trigger. This is the default.
	TooltipPositionTop string = "top"
	// TooltipPositionBottom shows the card below the trigger.
	TooltipPositionBottom string = "bottom"
	// TooltipPositionLeft shows the card to the left of the trigger.
	TooltipPositionLeft string = "left"
	// TooltipPositionRight shows the card to the right of the trigger.
	TooltipPositionRight string = "right"
)

// Tooltip Templ component props. Tooltip wraps a trigger element (the
// children) and reveals a dark floating card with a title, a description,
// and an optional illustration on hover or keyboard focus. The card keeps
// its dark color scheme in both light and dark modes.
type Tooltip struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Title is the card heading. Required.
	Title string
	// Description renders below the title in a muted color. Required.
	Description string
	// Image is the URL of an optional illustration shown above the title.
	// The component applies the rounded corners and border to it.
	Image string
	// Position places the card relative to the trigger: one of the
	// TooltipPosition constants. Defaults to top.
	Position string
}

// Empty reports whether the tooltip has no content to show.
func (t Tooltip) Empty() bool {
	return t.Title == "" && t.Description == "" && t.Image == ""
}
