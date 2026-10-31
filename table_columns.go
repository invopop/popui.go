package popui

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

const (
	tableColumnDefaultWidth    = 150
	tableColumnDefaultMinWidth = 60
)

// tableColumn is a props.TableColumn with its defaults applied.
type tableColumn struct {
	props.TableColumn
	Index int // counting from 1, as nth-child does
}

func (c tableColumn) right() bool {
	return c.Align == props.TableColumnAlignRight
}

func (c tableColumn) sortable() bool {
	return c.SortAscending != nil || c.SortDescending != nil
}

// tableColumnsOf applies the defaults to the columns.
func tableColumnsOf(cols []props.TableColumn) []tableColumn {
	out := make([]tableColumn, len(cols))
	for i, c := range cols {
		if c.Key == "" {
			c.Key = strconv.Itoa(i + 1)
		}
		if c.Width <= 0 {
			c.Width = tableColumnDefaultWidth
		}
		if c.MinWidth <= 0 {
			c.MinWidth = min(tableColumnDefaultMinWidth, c.Width)
		}
		out[i] = tableColumn{TableColumn: c, Index: i + 1}
	}
	return out
}

// tableFrozenCount is the number of leading columns that are frozen: up to
// the last one marked Frozen, as freezing a column freezes those before it.
func tableFrozenCount(cols []tableColumn) int {
	n := 0
	for _, c := range cols {
		if c.Frozen {
			n = c.Index
		}
	}
	return n
}

// tableColumnsClass names the stylesheet tableColumnsCSS writes for tables
// shaped like this one. It depends only on what the stylesheet does, so
// tables of the same shape share it, and rendering stays deterministic.
func tableColumnsClass(cols []tableColumn) string {
	h := fnv.New32a()
	fmt.Fprintf(h, "%d", len(cols))
	for _, c := range cols {
		if c.right() {
			fmt.Fprintf(h, ",%d", c.Index)
		}
	}
	return fmt.Sprintf("popui-table-cols-%08x", h.Sum32())
}

// tableColumnsCSS is the stylesheet for a table with Columns. Everything
// that changes as the user works with the table is kept out of it, on the
// table element, so popui.js only ever updates that:
//
//   - The widths are the custom properties --popui-col-N, which size the
//     <col> elements and the table, and offset the frozen columns.
//   - The frozen columns are listed by position in data-frozen ("1 2"),
//     and the last of them in data-frozen-last.
func tableColumnsCSS(cols []tableColumn) string {
	c := "." + tableColumnsClass(cols)
	n := len(cols)
	var b strings.Builder
	sum := make([]string, 0, n)
	for _, col := range cols {
		sum = append(sum, tableColumnVar(col.Index))
	}
	// As wide as the columns, or the container when that's wider.
	fmt.Fprintf(&b, "%s{width:max(100%%,calc(%s))}", c, strings.Join(sum, " + "))
	for _, col := range cols {
		k := col.Index
		// The last <col> has no width, so it takes any space left.
		if k < n {
			fmt.Fprintf(&b, "%s>colgroup>:nth-child(%d){width:%s}", c, k, tableColumnVar(k))
		}
		if col.right() {
			fmt.Fprintf(&b, "%s>tbody>tr>:nth-child(%d){text-align:right}", c, k)
			fmt.Fprintf(&b, "%s>tbody>tr>:nth-child(%d)>div{justify-content:flex-end}", c, k)
		}
		left := "0px"
		if k > 1 {
			left = "calc(" + strings.Join(sum[:k-1], " + ") + ")"
		}
		frozen := fmt.Sprintf(`%s[data-frozen~="%d"]>*>tr>:nth-child(%d)`, c, k, k)
		fmt.Fprintf(&b, "%s{position:sticky;left:%s;z-index:5}", frozen, left)
		// Above the sticky header's cells, which scroll under it.
		fmt.Fprintf(&b, `%s[data-frozen~="%d"]>thead>tr>:nth-child(%d){z-index:15}`, c, k, k)
		fmt.Fprintf(&b, `%s[data-frozen-last="%d"]>*>tr>:nth-child(%d){border-right:1px solid var(--color-neutral-20)}`, c, k, k)
		// The border separates the next column, so it needs no divider.
		if k < n {
			fmt.Fprintf(&b, `%s[data-frozen-last="%d"]>thead>tr>:nth-child(%d)::before{display:none}`, c, k, k+1)
		}
	}
	return b.String()
}

func tableColumnVar(k int) string {
	return fmt.Sprintf("var(--popui-col-%d)", k)
}

// tableColumnsAttributes are the attributes of a table with Columns: its
// widths and frozen columns, followed by the page's own attributes, with
// any style of theirs kept after the widths.
func tableColumnsAttributes(cols []tableColumn, attrs templ.Attributes) templ.Attributes {
	var style strings.Builder
	for _, c := range cols {
		fmt.Fprintf(&style, "--popui-col-%d:%dpx;", c.Index, c.Width)
	}
	out := templ.Attributes{"data-popui-columns": true}
	for k, v := range attrs {
		if k == "style" {
			if s, ok := v.(string); ok {
				style.WriteString(s)
				continue
			}
		}
		out[k] = v
	}
	out["style"] = style.String()
	if n := tableFrozenCount(cols); n > 0 {
		out["data-frozen"] = tableFrozenList(n)
		out["data-frozen-last"] = strconv.Itoa(n)
	}
	return out
}

func tableFrozenList(n int) string {
	s := make([]string, n)
	for i := range s {
		s[i] = strconv.Itoa(i + 1)
	}
	return strings.Join(s, " ")
}

// tableColumnAttributes are the attributes of a column's header cell.
func tableColumnAttributes(c tableColumn, frozen bool, resizable bool) templ.Attributes {
	out := templ.Attributes{
		"scope":    "col",
		"data-key": c.Key,
	}
	if resizable {
		out["data-min-width"] = strconv.Itoa(c.MinWidth)
	}
	if frozen {
		out["data-frozen"] = true
	}
	switch c.Sort {
	case "asc":
		out["aria-sort"] = "ascending"
	case "desc":
		out["aria-sort"] = "descending"
	}
	for k, v := range c.Attributes {
		out[k] = v
	}
	return out
}

// tableColumnHasMenu reports whether the column's header opens a menu.
func tableColumnHasMenu(c tableColumn, freezable bool) bool {
	return freezable || c.sortable() || c.MenuItems != nil
}

// tableColumnLabel is the accessible name of a column, for its resize
// handle.
func tableColumnLabel(c tableColumn) string {
	if c.Title != "" {
		return c.Title
	}
	return c.Key
}
