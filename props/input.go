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
	Size string `jsonschema:"enum=,enum=sm,enum=lg"`

	// Icon embeds the content inside the input field on the left side.
	Icon templ.Component

	// Prefix adds a text just before the input field inside a @Description.
	Prefix string

	// Local shows date, time, datetime-local, week, and month inputs in the
	// browser's own control, formatted for its locale. By default they're
	// text fields holding the ISO 8601 value as written, such as
	// "2026-09-20", with the system picker behind a calendar button;
	// datetime-local inputs then hold an RFC 3339 timestamp with a zone, such
	// as "2026-09-20T12:30:00Z". Other values are submitted the same either
	// way.
	Local bool

	Autofocus bool
	Readonly  bool
	Required  bool
	Disabled  bool
	Error     Error
}
