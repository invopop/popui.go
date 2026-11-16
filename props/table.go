package props

import "github.com/a-h/templ"

// Table Templ component props.
//
// Full-width layout, row-hover tint, and column-title dividers are baked into
// every Table — they're part of the standard data-table styling and need no
// configuration.
type Table struct {
	ID         string
	Class      string
	RootClass  string
	Attributes templ.Attributes
	Variant    string // "card" adds outer border

	// ScrollHorizontal turns the table into a wide, horizontally-scrolling
	// table: the body scrolls sideways. Pair it with StickyColumn to keep the
	// first column pinned during the scroll.
	ScrollHorizontal bool

	// StickyColumn pins the first column to the left — a full-height divider
	// stays put during a horizontal scroll so the identity column never
	// detaches from the rows. Only has a visible effect alongside
	// ScrollHorizontal (there's nothing to scroll past otherwise).
	StickyColumn bool

	// StickyHeader pins the <thead> row to the top while the body scrolls.
	// The RootClass must make the wrapper a vertical scroll container with a
	// bounded height (e.g. "flex-1 min-h-0 overflow-y-auto").
	StickyHeader bool

	// Resizable lets the user drag the right edge of each header cell to
	// resize that column. The drag handles and behavior are provided by
	// popui (CSS + popui.js); no consumer script is needed. With Columns,
	// each drag ends by dispatching a bubbling popui-table-column-resize
	// event from the table, with detail {key, width}.
	Resizable bool

	// Columns describes the table's columns. When set, Table renders the
	// <colgroup> and the <thead> from them, so the children are just the
	// rows (a <tbody>), and the table takes fixed column widths: it is as
	// wide as its columns or its container, whichever is wider, the last
	// column takes any space left, it scrolls sideways when it doesn't
	// fit, and cells truncate rather than wrap. ScrollHorizontal and
	// StickyColumn are implied (use Frozen instead) and ignored.
	Columns []TableColumn

	// Freezable adds Freeze and Unfreeze items to the header menu of each
	// column. Freezing a column freezes those before it too, and unfreezing
	// it unfreezes those after it. Each change dispatches a bubbling
	// popui-table-columns-freeze event from the table, with detail
	// {frozen: [keys]}. Only applies with Columns.
	Freezable bool
}

// TableColumn alignments.
const (
	TableColumnAlignLeft  = "left"
	TableColumnAlignRight = "right"
)

// TableColumn sort directions.
const (
	TableColumnSortAscending  = "asc"
	TableColumnSortDescending = "desc"
)

// TableColumn describes one column of a Table: its header, its width and
// how it behaves. Clicking or right-clicking a header opens its menu, when
// it has items: sorting, freezing (with the table's Freezable) and any
// MenuItems.
type TableColumn struct {
	// Key identifies the column in events. Defaults to its position,
	// counting from 1.
	Key string
	// Title is the header's text.
	Title string
	// Header replaces Title with any content, e.g. a select-all Checkbox.
	Header templ.Component
	// Width in pixels. Defaults to 150.
	Width int
	// MinWidth in pixels is the narrowest a user can resize the column to.
	// Defaults to 60, or Width if that's narrower.
	MinWidth int
	// Align is "left" (default) or "right", for the header and cells
	// alike, e.g. "right" for amounts.
	Align string `jsonschema:"enum=,enum=left,enum=right"`
	// Frozen keeps the column in place when the table scrolls sideways.
	// Only leading columns can be frozen, so this freezes every column
	// before it too.
	Frozen bool

	// Sort is the direction the table is sorted by this column, "asc" or
	// "desc", or empty when it isn't. It marks the header with an arrow and
	// the matching menu item with a tick.
	Sort string `jsonschema:"enum=,enum=asc,enum=desc"`
	// SortAscending and SortDescending are the attributes of the "Sort
	// ascending" and "Sort descending" header menu items, typically htmx
	// attributes (hx-get, hx-target, …) that reload the table sorted.
	// Columns without them aren't sortable.
	SortAscending  templ.Attributes
	SortDescending templ.Attributes
	// MenuItems adds items to the end of the column's header menu.
	MenuItems templ.Component
	// Attributes are added to the column's header cell.
	Attributes templ.Attributes
}

// TablePaginationElements defines custom attributes for pagination interactive elements
type TablePaginationElements struct {
	First      templ.Attributes
	Prev       templ.Attributes
	Next       templ.Attributes
	Last       templ.Attributes
	Page       templ.Attributes
	Select     templ.Attributes
	TotalPages templ.Attributes
	TotalItems templ.Attributes
}

// TablePagination Templ component props
type TablePagination struct {
	ID                 string
	Class              string
	Attributes         templ.Attributes
	CurrentPage        int
	TotalPages         int
	TotalItems         int
	RowsPerPage        int
	RowsPerPageOptions []int
	ShowRowsPerPage    bool
	ItemsLabel         string
	FirstPageURL       templ.SafeURL
	PrevPageURL        templ.SafeURL
	NextPageURL        templ.SafeURL
	LastPageURL        templ.SafeURL
	Elements           TablePaginationElements
}
