package props

import "github.com/a-h/templ"

const (
	// ToastTypeSuccess shows a green success icon.
	ToastTypeSuccess string = "success"
	// ToastTypeError shows a red error icon.
	ToastTypeError string = "error"
	// ToastTypeInfo shows a neutral info icon.
	ToastTypeInfo string = "info"
)

const (
	// ToastPositionTopLeft fixes the toast to the top left of the screen.
	ToastPositionTopLeft string = "top-left"
	// ToastPositionTopCenter fixes the toast to the top center of the screen.
	ToastPositionTopCenter string = "top-center"
	// ToastPositionTopRight fixes the toast to the top right of the screen.
	ToastPositionTopRight string = "top-right"
	// ToastPositionBottomLeft fixes the toast to the bottom left of the screen.
	ToastPositionBottomLeft string = "bottom-left"
	// ToastPositionBottomCenter fixes the toast to the bottom center of the screen.
	ToastPositionBottomCenter string = "bottom-center"
	// ToastPositionBottomRight fixes the toast to the bottom right of the screen.
	// This is the default.
	ToastPositionBottomRight string = "bottom-right"
)

// Toast Templ component props. Toast is a dark floating notification with a
// type icon, a message, an optional secondary description, and an optional
// action button on the right.
//
// Toasts are hidden until shown via popui.showToast(id) or a button with a
// data-toast-trigger="<toast-id>" attribute, and hide again automatically
// after Duration milliseconds. Only one toast is visible at a time.
type Toast struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Type selects the leading icon: "success", "error" or "info".
	Type string
	// Message is the main toast text.
	Message string
	// Description renders below the message in a muted color.
	Description string
	// Position fixes the toast to a screen corner or edge. One of the
	// ToastPosition constants. Defaults to bottom-right.
	Position string
	// Duration is how long the toast stays visible, in milliseconds.
	// Defaults to 3000 (three seconds).
	Duration int
	// ActionLabel renders an action button on the right side of the toast.
	ActionLabel string
	// ActionHref renders the action as a link instead of a button.
	ActionHref templ.SafeURL
	// ActionAttributes are additional HTML attributes for the action element,
	// e.g. click handlers.
	ActionAttributes templ.Attributes
}
