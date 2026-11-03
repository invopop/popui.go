package props

import "github.com/a-h/templ"

// Details layouts.
const (
	DetailsLayoutInline  = "inline"
	DetailsLayoutStacked = "stacked"
)

// Details defines the props for the Details component.
type Details struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Layout is "inline" for details whose labels are beside their values,
	// close together, as in a card or side panel, or "stacked" for details
	// whose labels are above their values, unless they set Inline. Default
	// is "inline".
	Layout string `jsonschema:"enum=,enum=inline,enum=stacked"`
}

// DetailAction is one hover-revealed action button on a Detail's value. Provide an Icon plus a behavior: Copy
// (clipboard), URL (external link), or a plain button driven by Attributes.
type DetailAction struct {
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

// Detail defines the props for the Detail component.
type Detail struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Value is the text shown for the label. Without one, the detail's
	// children are shown as the value, such as a TagStatus.
	Value string
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
	Actions []DetailAction
}

// DT defines the props for the DT component (definition term).
//
// Deprecated: Detail has a Label, and takes a component as its value.
type DT struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// DD defines the props for the DD component (definition description).
//
// Deprecated: Detail takes a component as its value.
type DD struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// DescriptionList defines the props for the DescriptionList component.
//
// Deprecated: use Details.
type DescriptionList = Details

// DescriptionListItem defines the props for the DescriptionListItem
// component.
//
// Deprecated: use Detail.
type DescriptionListItem = Detail

// DescriptionListAction is one hover-revealed action button on a
// DescriptionListItem's value.
//
// Deprecated: use DetailAction.
type DescriptionListAction = DetailAction
