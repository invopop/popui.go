package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

const (
	// CheckboxVariantSwitch represents a switch-style checkbox.
	CheckboxVariantSwitch string = "switch"
)

// Checkbox defines the properties for checkbox inputs.
type Checkbox struct {
	ID         string
	Class      string
	Name       string
	Value      string
	Checked    bool
	Autofocus  bool
	Disabled   bool
	Attributes templ.Attributes

	// Variant specifies the style of the checkbox, e.g.,
	// `CheckboxVariantSwitch` for toggle switch style.
	Variant string

	// Label is the text label associated with the checkbox.
	Label string
	// Description is additional descriptive text for the checkbox that will
	// be shown below the label.
	Description string
}

// GenerateID returns a new Checkbox instance with either the existing ID
// or a new randomly generated one. This is useful for ensuring that
// checkboxes have unique IDs when none are provided and is designed to be used
// inline:
//
//	@popui.Checkbox(props.Checkbox{Name: "agree", Label: "I agree"}.GenerateID())
func (c Checkbox) GenerateID() Checkbox {
	if c.ID != "" {
		return c
	}
	// generate a short random identifier
	c.ID = fmt.Sprintf("%s%06d", c.Name, rand.Intn(100000))
	return c
}

// OptionGroup properties for grouping form elements like checkboxes and radio buttons.
type OptionGroup struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
}
