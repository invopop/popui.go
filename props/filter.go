package props

import "github.com/a-h/templ"

// Filter renders a search-bar–style row at the top of a data view: a
// "+ Filter" menu plus one editable chip per active filter. Multi-filter
// UX — picking a field from the menu appends an editable chip (it does not
// replace active filters); several filters apply at once (AND), laid out
// left-to-right in add order. Backed by the `filterRow` Alpine controller
// registered in popui.js (active list, add/remove, auto-open on add,
// cleared values on remove).
//
// The whole row is one <form> that fires hx-get on submit. Caller wires
// the HTMX target/select/swap so the swap can be scoped to the data
// region (keeps the filter row alive — its open option panel survives
// submissions).
type Filter struct {
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

	// Select is the hx-select selector — extracts a matching region from
	// the response (so the server can return a full page and only the
	// region is swapped). Defaults to Target when empty, which is the
	// usual case (select and swap the same region).
	Select string

	// Swap is the hx-swap mode (innerHTML, outerHTML, etc.). Defaults to
	// "outerHTML" when empty — replace the target region element whole.
	Swap string

	// PageSize is carried as a hidden `size` input so the page-size
	// preference survives filter changes. Pass 0 to omit.
	PageSize int

	Inputs []FilterInput
}

// Filter editor types for FilterInput.Type. The type picks the value
// editor rendered in the chip and the operator label shown next to it.
const (
	// FilterTypeText renders a free-text <input>. Operator: "matches".
	// This is the default when Type is empty.
	FilterTypeText string = "text"
	// FilterTypeSelect renders a single-choice option list. Same inline
	// format as FilterTypeMultiple (an option's Color shows a TagStatus
	// dot); picking a value replaces the previous one. Operator: "matches".
	FilterTypeSelect string = "select"
	// FilterTypeMultiple renders a multi-choice option list — one chip, many
	// values. Operator: "is any of".
	FilterTypeMultiple string = "multiple"
	// FilterTypeCalendar renders a dual-month date-range calendar with a
	// preset rail (this/last week, month, quarter + custom). The selected
	// range is submitted as a single "YYYY-MM-DD..YYYY-MM-DD" value under
	// Name. Operator: "is between".
	FilterTypeCalendar string = "calendar"
)

// FilterInput describes one filterable field rendered as a chip. Type
// selects the value editor; see the FilterType* constants.
type FilterInput struct {
	Name  string
	Label string

	// Type selects the value editor and operator label. One of the
	// FilterType* constants; defaults to FilterTypeText when empty.
	Type string

	// PluralLabel overrides the auto-pluralised label shown in the chip's
	// summary when several options are selected ("3 statuses"). Leave empty
	// to use the built-in English pluralizer; set it for irregular plurals
	// (e.g. Label "Person" → PluralLabel "people"). Only used by
	// FilterTypeMultiple.
	PluralLabel string

	// Icon is rendered both in the "+ Filter" menu and on the active chip.
	Icon templ.Component

	// Values are the currently-applied values (typically pulled from the
	// URL query string by the caller). A non-empty Values means the chip
	// starts active on render.
	Values []string

	// Options are the choices for FilterTypeSelect and FilterTypeMultiple.
	Options []FilterOption

	// Presets selects the calendar preset rail for FilterTypeCalendar — which
	// range shortcuts appear and in what order. Empty uses
	// DefaultCalendarPresets(). Ignored for non-calendar types.
	Presets []CalendarPreset
}

// FilterOption is one row in a select/multiple FilterInput's value list.
// Color is a TagStatus.Status value (green, orange, blue, etc.); a
// non-empty Color shows the matching dot beside the option.
type FilterOption struct {
	Value string
	Label string
	Color string
}

// MaxStackedDots caps how many overlapping colour dots are rendered in a
// multi-selection summary; selected options beyond the cap add to the
// count label but draw no dot.
const MaxStackedDots = 3
