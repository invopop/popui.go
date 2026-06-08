package props

import "github.com/a-h/templ"

// FilterRow renders a search-bar–style row at the top of a data view: a
// "+ Filter" menu plus one editable chip per active filter. Single-filter
// UX — picking a field from the menu replaces whatever chip was active
// before. Backed by the `filterRow` Alpine controller registered in
// popui.js (active list, add/remove, auto-open on add, cleared values on
// swap).
//
// The whole row is one <form> that fires hx-get on submit. Caller wires
// the HTMX target/select/swap so the swap can be scoped to the data
// region (keeps the filter row alive — its popovers stay open across
// submissions).
type FilterRow struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// BaseURL is the form's hx-get target (typically the path the page
	// lives at). Filter values become query params on this URL.
	BaseURL string

	// Target is the hx-target selector — the region to swap when the
	// form submits. If empty, no hx-target is emitted (HTMX defaults to
	// document.body, which is rarely what you want; pass an explicit
	// id selector like "#dashboard-content").
	Target string

	// Select is the hx-select selector — extracts a matching region
	// from the response. Set this to the same selector as Target for
	// in-place swaps where the server returns the full page.
	Select string

	// Swap is the hx-swap mode (innerHTML, outerHTML, etc.). If empty,
	// HTMX's default applies.
	Swap string

	// PageSize is carried as a hidden `size` input so the page-size
	// preference survives filter changes. Pass 0 to omit.
	PageSize int

	Inputs []FilterInput
}

// FilterInput describes one filterable field rendered as a chip. The
// chip's value editor is chosen by the input's shape:
//   - Options with a Color → popui.DropdownSelect with TagStatus dots.
//   - Options without Color → plain <select>.
//   - No Options → free-text <input>.
type FilterInput struct {
	Name  string
	Label string

	// Icon is rendered both in the "+ Filter" menu and on the active chip.
	Icon templ.Component

	// Values are the currently-applied values (typically pulled from the
	// URL query string by the caller). A non-empty Values means the chip
	// starts active on render.
	Values []string

	Options []FilterOption

	// Multi turns the DropdownSelect into a multi-select (one chip,
	// many values). Only honoured for the DropdownSelect editor path
	// (options with Color). Also flips the operator label from
	// "matches" to "is any of".
	Multi bool
}

// FilterOption is one row in a FilterInput's value dropdown. Color is a
// TagStatus.Status value (green, orange, blue, etc.) — non-empty Color
// switches the editor from a plain <select> to popui.DropdownSelect with
// the matching coloured dot.
type FilterOption struct {
	Value string
	Label string
	Color string
}
