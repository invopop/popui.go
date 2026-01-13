package props

import "github.com/a-h/templ"

const (
	// FormVariantCard represents a card-style form.
	FormVariantCard string = "card"
)

// Form properties to configure a form element with standard HTML form attributes.
type Form struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Action specifies where to send the form data when submitted.
	Action string
	// Method specifies the HTTP method to use when submitting the form (GET, POST, etc.).
	Method string
	// Enctype specifies how the form data should be encoded when submitting (e.g., "multipart/form-data").
	Enctype string
	// Target specifies where to display the response after submitting the form (_self, _blank, etc.).
	Target string
	// Autocomplete specifies whether the form should have autocomplete on or off.
	Autocomplete string
	// Novalidate when true, the form will not be validated when submitted.
	Novalidate bool
	// AcceptCharset specifies the character encodings used for form submission.
	AcceptCharset string
	// Name sets the name of the form.
	Name string
	// Rel specifies the relationship between the current document and the linked document.
	Rel string
}
