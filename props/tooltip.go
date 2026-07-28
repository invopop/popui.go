package props

import "github.com/a-h/templ"

const (
	// TooltipPositionTop shows the card above the trigger, centered on it. This
	// is the default.
	TooltipPositionTop string = "top"
	// TooltipPositionTopStart shows the card above the trigger, with their
	// leading edges aligned. Use it when a centered card would overflow to the
	// left, e.g. next to a label at the start of a row.
	TooltipPositionTopStart string = "top-start"
	// TooltipPositionTopEnd shows the card above the trigger, with their
	// trailing edges aligned. Use it when a centered card would overflow to the
	// right.
	TooltipPositionTopEnd string = "top-end"
	// TooltipPositionBottom shows the card below the trigger, centered on it.
	TooltipPositionBottom string = "bottom"
	// TooltipPositionBottomStart shows the card below the trigger, with their
	// leading edges aligned.
	TooltipPositionBottomStart string = "bottom-start"
	// TooltipPositionBottomEnd shows the card below the trigger, with their
	// trailing edges aligned.
	TooltipPositionBottomEnd string = "bottom-end"
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
	// Title is the card heading. Optional, but at least one of Title,
	// Description or Image must be set for the card to render.
	Title string
	// Description renders below the title in a muted color, parsed as Markdown:
	// emphasis, code, links and bullet or numbered lists all work, and a single
	// newline is a line break, so multi-line text reads as written. Raw HTML is
	// escaped rather than passed through. Optional, but at least one of Title,
	// Description or Image must be set for the card to render.
	Description string
	// Image is the URL of an optional illustration shown above the title.
	// The component applies the rounded corners and border to it.
	Image string
	// Position places the card relative to the trigger: one of the
	// TooltipPosition constants. Defaults to top.
	//
	// The card has no collision detection, so a wide one centered on a trigger
	// near the edge of its container can overflow. The -start and -end variants
	// align the card with the trigger's leading or trailing edge instead of
	// centering it, which keeps it inside the container.
	Position string
}

// Empty reports whether the tooltip has no content to show.
func (t Tooltip) Empty() bool {
	return t.Title == "" && t.Description == "" && t.Image == ""
}
