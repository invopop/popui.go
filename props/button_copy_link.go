package props

import "github.com/a-h/templ"

// ButtonCopyLink renders an id-style value with two affordances: a
// copy-to-clipboard button (the duplicate icon) and an external-link button
// (opens URL in a new tab). Use it for cells that reference a record both
// worth copying and viewable elsewhere — e.g. a silo entry id that links out
// to the Invopop console.
type ButtonCopyLink struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Variant / Size are forwarded to the copy button (popui.ButtonCopy):
	// Variant "outline" (default) or "ghost"/"transparent"; Size "sm"/"md"/"lg".
	Variant string
	Size    string

	// Value is displayed (truncated to Prefix/Suffix) and copied to the
	// clipboard by the copy button.
	Value string

	// URL is opened in a new tab by the external-link button. Caller-trusted
	// (built by the consumer), so it's a templ.SafeURL. Empty hides the link
	// button, leaving just the copy button.
	URL templ.SafeURL

	PrefixLength int
	SuffixLength int
}
