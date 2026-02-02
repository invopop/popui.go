package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Popover Templ component props
type Popover struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// ButtonLabel will be the text shown on the popover button, the
	// default if not provided is "···" (three mid-dots)
	ButtonLabel string
	// ButtonVariant sets the button style, see the button component
	// for available variants.
	ButtonVariant string
	// RightAlign determines if the popover menu is aligned to the right
	// of the button, default is left aligned. This will also change the
	// location of the popover.
	RightAlign bool
}

// PopoverItem menu item component props
type PopoverItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// GenerateID generates a unique ID for the Popover if none is provided.
func (po Popover) GenerateID() Popover {
	if po.ID != "" {
		return po
	}
	// generate a short random identifier
	po.ID = fmt.Sprintf("popover-%06d", rand.Intn(100000))
	return po
}
