package props

import "github.com/a-h/templ"

// Notification defines the props for the Notification component.
type Notification struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Type selects the color scheme and leading icon: "info", "warning",
	// "error", "success" or "neutral".
	Type string
	// Text is the main notification text, rendered across the full width of
	// the notification.
	Text string
	// Description renders below the text in a muted color.
	Description string
	// Actions renders one or more buttons on the right of the description
	// row, aligned to the right edge of the notification.
	Actions templ.Component
}
