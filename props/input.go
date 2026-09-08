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

	// Tooltip shows a question mark icon after the label that reveals a
	// Tooltip card on hover. Only used when Label is set.
	Tooltip Tooltip

	// Set the height of the input field, see the InputSize constants.
	Size string

	// Icon embeds the content inside the input field on the left side.
	Icon templ.Component

	// Prefix adds a text just before the input field inside a @Description.
	Prefix string

	// ISO turns a datetime-local input into a text field holding an RFC 3339
	// timestamp, such as "2026-09-20T12:30:00Z", so the value can be read,
	// copied, and pasted directly. The calendar button opens a native
	// datetime-local picker that fills the field with the chosen instant in
	// the browser's local offset (popui.js). The field is validated against
	// an RFC 3339 pattern on submit; an empty field submits an empty string.
	// Only used when Type is "datetime-local".
	ISO bool

	Autofocus bool
	Readonly  bool
	Required  bool
	Disabled  bool
	Error     Error
}
