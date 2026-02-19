package props

import (
	"github.com/a-h/templ"
	"github.com/google/uuid"
)

// ContextMenu provides props for the context menu component.
type ContextMenu struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// RightAlign aligns the menu to the right of the button.
	RightAlign bool

	// ButtonVariant defines the style of the trigger button.
	ButtonVariant string

	// ButtonLabel sets a text label on the trigger button instead of the default dots icon.
	ButtonLabel string
}

// GenerateID ensures the context menu has a unique ID and returns itself.
func (c ContextMenu) GenerateID() ContextMenu {
	if c.ID == "" {
		c.ID = "ctx-" + uuid.New().String()[:8]
	}
	return c
}

// ContextMenuItem provides props for individual items within a context menu.
type ContextMenuItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
