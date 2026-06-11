package props

import "github.com/a-h/templ"

// DetailList is the container for a console-ui-style detail view — a
// vertical list of label-on-left / value-on-right rows, as used in side
// panels and drawers. Wrap popui.DetailItem / popui.DetailRow children in
// it. App-agnostic: the styling ships with popui so consumers don't copy a
// <style> block.
type DetailList struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Bordered draws a bottom divider, turning the list into a labelled
	// section that's visually separated from what follows.
	Bordered bool
	// Indent insets the list horizontally so a Bordered divider lines up
	// with the label/value gutter instead of spanning the full width.
	Indent bool
}

// DetailItem is a single text row inside a popui.DetailList: a muted,
// fixed-width Label on the left and a Value on the right. An empty Value
// renders as an em-dash so column alignment stays stable for absent fields.
// For non-text values (tags, links, copy buttons) use popui.DetailRow.
type DetailItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	Value      string
}

// DetailRow is a detail row whose value cell is filled by arbitrary
// children — a popui.Tag, a link, a popui.ButtonCopy, etc. — instead of a
// plain string. Same label column and hover-tinted value cell as
// DetailItem.
type DetailRow struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Copyable turns the value cell into a reveal-on-hover affordance row:
	// children carrying the `popui-detail-action` class (and any
	// popui.ButtonCopy icons) stay invisible until the row is hovered, and
	// popui.ButtonCopy chrome is stripped so the value reads as plain text.
	// Use for id / copy-to-clipboard rows.
	Copyable bool
}
