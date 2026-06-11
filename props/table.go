package props

import "github.com/a-h/templ"

// Table Templ component props
type Table struct {
	ID               string
	Class            string
	RootClass        string
	Attributes       templ.Attributes
	Variant          string // "card" adds outer border
	ScrollHorizontal bool   // Enable horizontal scrolling for wide tables

	// StickyHeader pins the <thead> row to the top while the body scrolls.
	// The RootClass must make the wrapper a vertical scroll container with a
	// bounded height (e.g. "flex-1 min-h-0 overflow-y-auto").
	StickyHeader bool

	// FreezeFirstColumn pins the first column (header + every body row) to
	// the left edge while the table scrolls horizontally. Its right edge
	// carries a full-height divider that stays in place during the scroll,
	// so the identity column never detaches from the rows. Requires
	// ScrollHorizontal.
	FreezeFirstColumn bool

	// ColumnDividers draws a label-height vertical divider before each
	// header column, separating the column titles. When FreezeFirstColumn
	// is also set, the divider before the 2nd column is suppressed because
	// the frozen column's full-height right border already separates 1|2.
	ColumnDividers bool

	// FullWidth stretches the table to fill its container width even when
	// the columns don't need it, while still growing past the container
	// when they do. Pairs with ScrollHorizontal for wide tables.
	FullWidth bool
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
