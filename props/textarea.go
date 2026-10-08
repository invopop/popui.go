package props

import (
	"strconv"

	"github.com/a-h/templ"
)

// Textarea Templ component props
type Textarea struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Name        string
	Placeholder string
	Value       string
	Label       string
	Disabled    bool

	// Tooltip shows a question mark icon after the label that reveals a
	// Tooltip card on hover. Only used when Label is set.
	Tooltip   Tooltip
	Readonly  bool
	Required  bool
	Autofocus bool
	Rows      int

	// Monospace indicates whether to use a monospace font for the textarea
	// by adding the appropriate classes.
	Monospace bool

	// VariableFormat turns Contenteditable into a template editor. It is the
	// syntax variables are written in, %s standing for the name: "{{.%s}}"
	// for a Go template. Each variable in the text is shown as a chip, and
	// the text with the variables written out — the template itself — is
	// what Value holds, what is submitted under Name, and what an Alpine
	// x-model on the component binds. FieldPicker inserts into such an
	// editor as chips. Textarea ignores it.
	VariableFormat string
	// ViewToggle adds a Rich / Plain switch above the editor: Rich shows
	// the chips, Plain shows the template as text, in a monospaced
	// textarea. Only with VariableFormat.
	ViewToggle bool
	// View is the view the editor starts in: ContenteditableViewRich (the
	// default) or ContenteditableViewPlain. Only with VariableFormat.
	View string

	Error Error
}

// The views a Contenteditable with a VariableFormat can show.
const (
	ContenteditableViewRich  = "rich"
	ContenteditableViewPlain = "plain"
)

// GetRows returns the Rows prop as a string with a default if not present
func (t Textarea) GetRows() string {
	if t.Rows == 0 {
		return "4"
	}

	return strconv.Itoa(t.Rows)
}
