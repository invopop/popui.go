package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// ContextMenu Templ component props
type ContextMenu struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// ButtonLabel will be the text shown on the context menu button, the
	// default if not provided is "···" (three mid-dots)
	ButtonLabel string
	// ButtonVariant sets the button style, see the button component
	// for available variants.
	ButtonVariant string
	// RightAlign determines if the context menu is aligned to the right
	// of the button, default is left aligned. This will also change the
	// location of the context menu.
	RightAlign bool
}

// ContextMenuItem menu item component props
type ContextMenuItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// GenerateID generates a unique ID for the ContextMenu if none is provided.
func (cm ContextMenu) GenerateID() ContextMenu {
	if cm.ID != "" {
		return cm
	}
	// generate a short random identifier
	cm.ID = fmt.Sprintf("context-menu-%06d", rand.Intn(100000))
	return cm
}
