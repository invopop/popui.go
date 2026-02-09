package props

import "github.com/a-h/templ"

// Input sizes if not using default.
const (
	InputSizeSmall string = "sm"
	InputSizeLarge string = "lg"
)

// Input defines the properties that can be used with input fields.
type Input struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Type        string
	Placeholder string
	Value       string
	Name        string
	Label       string

	// Set the height of the input field, see the InputSize constants.
	Size string

	// Icon embeds the content inside the input field on the left side.
	Icon templ.Component

	// Prefix adds a text just before the input field inside a @Description.
	Prefix string

	Autofocus bool
	Readonly  bool
	Required  bool
	Disabled  bool
	Error     Error
}
