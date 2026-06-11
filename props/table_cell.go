package props

import "github.com/a-h/templ"

// TableCell renders the content of a single table cell through a consistent
// style. Use it inside a <td> (or <th>) within popui.Table. An empty Value
// renders a muted dash placeholder so missing data reads differently from a
// real value. For non-text content (buttons, tags, links) put that directly
// in the cell instead of using TableCell.
type TableCell struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Value      string // text content; empty renders a muted dash
	Mono       bool   // render Value in the monospace face (ids, hashes)
}
