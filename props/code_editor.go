package props

import (
	"strconv"

	"github.com/a-h/templ"
)

// CodeEditor Templ component props
type CodeEditor struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Name        string
	Placeholder string
	Value       string
	Label       string
	Disabled    bool
	Readonly    bool
	Required    bool
	Autofocus   bool
	Rows        int

	// MaxRows sets the maximum number of visible rows before scrolling.
	// Defaults to 20 if not set.
	MaxRows int

	Error Error
}

// GetRows returns the Rows prop as a string with a default if not present
func (c CodeEditor) GetRows() string {
	if c.Rows == 0 {
		return "4"
	}
	return strconv.Itoa(c.Rows)
}

// GetMaxRows returns the MaxRows prop as a string with a default if not present
func (c CodeEditor) GetMaxRows() string {
	if c.MaxRows == 0 {
		return "20"
	}
	return strconv.Itoa(c.MaxRows)
}
