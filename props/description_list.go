package props

import "github.com/a-h/templ"

// DescriptionList defines the props for the DescriptionList component.
type DescriptionList struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// DescriptionListAction is one hover-revealed action button on a
// DescriptionListItem's value. Provide an Icon plus a behavior: Copy
// (clipboard), URL (external link), or a plain button driven by Attributes.
type DescriptionListAction struct {
	// Icon is the glyph rendered inside the button — any popui icon component,
	// e.g. icons.Duplicate() or icons.ExternalLink().
	Icon templ.Component
	// Copy, when set, wires the button to copy this text to the clipboard.
	Copy string
	// URL, when set, renders the button as an anchor that opens the URL in a new
	// tab (target=_blank, rel=noopener noreferrer).
	URL templ.SafeURL
	// Attributes wire extra behavior/extras on the button (aria-label, onclick,
	// …) and override the defaults a Copy / URL action sets.
	Attributes templ.Attributes
}

// DescriptionListItem defines the props for the DescriptionListItem component.
type DescriptionListItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	Value      string
	// Inline lays the item out as label-left / value-right instead of the
	// default stacked (label above value).
	Inline bool
	// Mono renders the value in the monospace face (ids, hashes).
	Mono bool
	// PrefixLength / SuffixLength truncate the DISPLAYED value to
	// "prefix…suffix" (or "prefix…" when SuffixLength is 0); the full value is
	// what Actions like copy act on. Both 0 shows it untruncated.
	PrefixLength int
	SuffixLength int
	// Actions are hover-revealed buttons rendered next to the value (copy,
	// external-link, …), in order of placement.
	Actions []DescriptionListAction
}

// DT defines the props for the DT component (definition term).
type DT struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// DD defines the props for the DD component (definition description).
type DD struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
