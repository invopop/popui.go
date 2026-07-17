package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Menu Templ component props
type Menu struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// ButtonLabel will be the text shown on the context menu button, the
	// default if not provided is "···" (three mid-dots)
	ButtonLabel string
	// ButtonVariant sets the button style, see the button component
	// for available variants.
	ButtonVariant string
	// Trigger replaces the default button with any component (an Avatar, an
	// icon, a styled row, …) as the menu opener. It is wrapped in an
	// unstyled <button> that carries the popover wiring, so pass
	// non-interactive content. ButtonLabel and ButtonVariant are ignored
	// when set.
	Trigger templ.Component
	// RightAlign determines if the context menu is aligned to the right
	// of the button, default is left aligned. This will also change the
	// location of the context menu.
	RightAlign bool
}

// MenuItem menu item component props
type MenuItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// GenerateID generates a unique ID for the Menu if none is provided.
func (cm Menu) GenerateID() Menu {
	if cm.ID != "" {
		return cm
	}
	// generate a short random identifier
	cm.ID = fmt.Sprintf("context-menu-%06d", rand.Intn(100000))
	return cm
}
