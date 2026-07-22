package props

import "github.com/a-h/templ"

// DrawerContextItem renders one row inside a drawer-style
// option list. Mirrors @invopop/popui Svelte DrawerContextItem.svelte:
// optional colour-tagged prefix (TagStatus dot-only), label, and a
// right-side affordance — either a reactive checkbox (multi-mode) or a
// tick icon (single-mode, shown only when this row is selected).
//
// Selection state is REACTIVE: the parent owns an Alpine `values` array
// (or single `value`) and passes IsSelectedExpr — an Alpine expression
// evaluating to bool — so the checkbox / tick toggles without a server
// round-trip. OnClickExpr is the Alpine handler invoked on row click.
type DrawerContextItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	Value string
	Label string

	// Color is a TagStatus.Status value (green, orange, blue, etc.). When
	// non-empty, a TagStatus dot-only pill is rendered to the left of the
	// label, giving the row the same colour vocabulary as the matching
	// TagStatus rendered elsewhere on the page.
	Color string

	// Multiple switches the right-side affordance from a tick (single
	// mode, shown only when selected) to a checkbox (multi mode, always
	// rendered, reactive to IsSelectedExpr).
	Multiple bool

	// IsSelectedExpr is an Alpine expression evaluating to bool — used to
	// drive the right-side checkbox/tick visibility. Example:
	//   "values.includes('production')"
	IsSelectedExpr string

	// OnClickExpr is the Alpine handler run when the row is clicked.
	// Example: "toggle('production')"
	OnClickExpr string

	// Tabindex, when non-empty, sets the row button's tabindex attribute.
	// Pass "-1" to take the row out of the Tab order — used by virtual-cursor
	// lists (e.g. Filter's inline option list) where a parent element owns
	// keyboard focus and the highlight is driven by an index, not DOM focus,
	// so individual rows must not steal focus. Empty → no tabindex (the row is
	// a normal tab stop).
	Tabindex string

	// HighlightExpr is an optional Alpine expression evaluating to bool. When
	// non-empty, it drives a keyboard-highlight class on the row (the same
	// `bg-background-default-secondary` used for hover) so a parent controller
	// can paint the row under an arrow-key cursor. Empty → no binding, so
	// existing callers are unaffected. Example:
	//   "activeIndex === 2"
	HighlightExpr string
}
