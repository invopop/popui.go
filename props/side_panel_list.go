package props

import "github.com/a-h/templ"

// SidePanelList is the container for a console-ui-style detail view — a
// vertical list of label-on-left / value-on-right rows, as used in side
// panels. Wrap popui.SidePanelItem / popui.SidePanelRow children in
// it. App-agnostic: the styling ships with popui so consumers don't copy a
// <style> block.
type SidePanelList struct {
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

// SidePanelItem is a single text row inside a popui.SidePanelList: a muted,
// fixed-width Label on the left and a Value on the right. An empty Value
// renders as an em-dash so column alignment stays stable for absent fields.
// For non-text values (tags, links, copy buttons) use popui.SidePanelRow.
type SidePanelItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	Value      string
}

// SidePanelHeading is a section title for a detail view (e.g. "Details",
// "Inboxes"). Renders the heading with the standard padding/typography so
// apps don't hand-roll the class string. Use Title, or pass children.
type SidePanelHeading struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Title      string
}

// SidePanelEmpty is the muted placeholder shown when a detail section / tab has
// no data ("No inboxes", "No reporting state", …). Use Text, or children.
type SidePanelEmpty struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Text       string
}

// SidePanelCopyRow is a copyable detail row: the Value reads as plain text and,
// on hover, the row tints and a copy button (and an optional external-link
// button) fade in — the clipboard wiring is built in, so apps just pass a
// label and value. It's the batteries-included form of popui.SidePanelRow with
// Copyable set.
type SidePanelCopyRow struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Value is shown (optionally truncated to Prefix…Suffix) and is what the
	// copy button writes to the clipboard in full.
	Value string
	// Stacked places the label above the value (see SidePanelRow.Stacked).
	Stacked bool
	// Mono renders the value in the monospace face (ids, hashes).
	Mono bool
	// PrefixLength / SuffixLength truncate the DISPLAYED value to
	// "prefix…suffix" (or "prefix…" when SuffixLength is 0); the full Value is
	// still copied. Both 0 shows the value untruncated.
	PrefixLength int
	SuffixLength int
	// URL, when set, adds a hover-revealed external-link button that opens it
	// in a new tab.
	URL templ.SafeURL
}

// SidePanelActionRow is the general form of SidePanelCopyRow: a label + value
// row with a single hover-revealed action button whose icon and behavior the
// caller supplies. Use it instead of hand-rolling a SidePanelRow with a
// popui.Button child — it gives copy, external-link, download, … rows the same
// hover reveal and sizing (24px, or 16px when Stacked) without duplication.
type SidePanelActionRow struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Value is shown (optionally truncated to Prefix…Suffix); the full value is
	// what behaviors like copy act on.
	Value string
	// Stacked places the label above the value (see SidePanelRow.Stacked) and
	// switches the action button to the compact 16px size.
	Stacked bool
	// Mono renders the value in the monospace face (ids, hashes).
	Mono bool
	// PrefixLength / SuffixLength truncate the DISPLAYED value to
	// "prefix…suffix" (or "prefix…" when SuffixLength is 0). Both 0 shows it
	// untruncated.
	PrefixLength int
	SuffixLength int
	// Icon is the glyph rendered inside the hover-revealed action button — any
	// popui icon component, e.g. icons.Duplicate() or icons.ExternalLink(). Nil
	// renders the row with no action button (value only).
	Icon templ.Component
	// URL, when set, turns the action button into an anchor that opens the URL
	// in a new tab (target=_blank, rel=noopener noreferrer) — e.g. an external
	// link. Leave empty for a plain button driven by ButtonAttributes.
	URL templ.SafeURL
	// ButtonAttributes wire the action button's behavior and extras — an
	// onclick for a non-link action, plus aria-label etc. Applied whether the
	// button renders as a <button> or (with URL set) an <a>.
	ButtonAttributes templ.Attributes
}

// SidePanelRow is a detail row whose value cell is filled by arbitrary
// children — a popui.Tag, a link, a popui.ButtonCopy, etc. — instead of a
// plain string. Same label column and hover-tinted value cell as
// SidePanelItem.
type SidePanelRow struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Stacked places the label above the value (a vertical row) instead of
	// the default label-left / value-right layout.
	Stacked bool
	// Copyable turns the value cell into a reveal-on-hover affordance row:
	// children carrying the `popui-detail-action` class (and any
	// popui.ButtonCopy icons) stay invisible until the row is hovered, and
	// popui.ButtonCopy chrome is stripped so the value reads as plain text.
	// Use for id / copy-to-clipboard rows.
	Copyable bool
}
