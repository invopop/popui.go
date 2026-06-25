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
	// table: the body scrolls sideways and the first column is automatically
	// frozen — pinned to the left with a full-height divider that stays put
	// during the scroll, so the identity column never detaches from the rows.
	ScrollHorizontal bool

	// StickyHeader pins the <thead> row to the top while the body scrolls.
	// The RootClass must make the wrapper a vertical scroll container with a
	// bounded height (e.g. "flex-1 min-h-0 overflow-y-auto").
	StickyHeader bool

	// Resizable lets the user drag the right edge of each header cell to
	// resize that column. The drag handles and behavior are provided by
	// popui (CSS + popui.js); no consumer script is needed.
	Resizable bool
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
