package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// MenuItemVariant constants
const (
	// MenuItemVariantDanger colors the item red for destructive actions
	// (delete, sign out, revoke).
	MenuItemVariantDanger string = "danger"
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
	// TriggerClass adds classes to the unstyled <button> wrapping a custom
	// Trigger — e.g. "w-full" so a row-shaped trigger stretches. Ignored
	// without Trigger.
	TriggerClass string
	// RootClass adds classes to the element holding the trigger and panel,
	// which is inline-block by default — set "block w-full" (with a
	// TriggerClass of "w-full") to let the menu fill its container.
	RootClass string
	// RightAlign determines if the context menu is aligned to the right
	// of the button, default is left aligned. This will also change the
	// location of the context menu.
	RightAlign bool
	// DropUp opens the panel above the trigger instead of below it, for
	// triggers that sit near the bottom of the viewport (a sidebar footer,
	// a table's last row).
	DropUp bool
}

// MenuItem menu item component props
type MenuItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Variant styles the item; see the MenuItemVariant constants.
	Variant string
	// Selected marks the item as the current choice with a trailing tick,
	// for menus that act as pickers.
	Selected bool
}

// MenuLabel is a non-interactive group heading within a Menu.
type MenuLabel struct {
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
