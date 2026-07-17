package props

import "github.com/a-h/templ"

// ButtonCopy Templ component props.
//
// Deprecated: use Button with the Copy, CopyPrefixLength and
// CopySuffixLength fields instead.
type ButtonCopy struct {
	ID           string
	Class        string
	Attributes   templ.Attributes
	Variant      string // Forwarded to Button.Variant
	Size         string // Forwarded to Button.Size
	Value        string // The text value to display and copy
	PrefixLength int    // Number of characters to show at the beginning
	SuffixLength int    // Number of characters to show at the end
}
