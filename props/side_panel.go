package props

import "github.com/a-h/templ"

// SidePanel position constants for the Position prop. Side panels anchor to
// one edge of the viewport; pick the side that's furthest from the
// app's primary navigation.
const (
	SidePanelPositionRight string = "right"
	SidePanelPositionLeft  string = "left"
)

// SidePanel renders a floating, fixed-position side panel that overlays
// one edge of the viewport. Non-blocking by default — no backdrop, the
// rest of the app stays interactive — and stays mounted across HTMX
// content swaps so the open/close state survives. Consumers typically
// host a swap target inside the panel's children block, fill it via
// hx-target, and dispatch the open event after the swap.
//
// Toggled via the `popui-sidepanel-open` / `popui-sidepanel-close` window
// events whose `event.detail` must equal the panel's ID. The
// controller scopes those by ID so several side panels can coexist on a
// single page without trampling each other; it also re-broadcasts the
// matching event whenever the open state flips, so external listeners
// react to every close path (X click, Escape, programmatic).
type SidePanel struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Position picks the viewport edge the side panel anchors to.
	// One of SidePanelPositionRight (default) or SidePanelPositionLeft.
	// Empty / unknown values fall back to right.
	Position string

	// Width sets the panel's fixed pixel width. Zero falls back to
	// 400 — a comfortable detail-panel default. Anything between
	// 320 and 720 reads well without crowding the main content.
	Width int
}

// SidePanelHeader renders a sticky header bar at the top of a popui.SidePanel:
// an X close button on the leading edge, a single-line truncated title
// in the middle, and an optional action slot on the trailing edge
// (children of the component). Lives inside a popui.SidePanel's children
// block so it sits above the scrollable content area.
//
// The close button dispatches `popui-sidepanel-close` (event.detail =
// SidePanelID) by default. Pass CloseAttributes to override that click
// behavior — e.g. to also clear a URL query param via
// history.replaceState — but make sure the override still dispatches
// the close event so the panel's Alpine state syncs.
type SidePanelHeader struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// SidePanelID is the parent SidePanel's ID — used by the default close
	// button to dispatch `popui-sidepanel-close` with the matching
	// `event.detail` so the right panel toggles closed.
	SidePanelID string

	// Title is the h1 heading rendered in the centre cell. Truncates
	// with `truncate` if it overflows the available width.
	Title string

	// CloseAttributes replace the close button's default attributes
	// (a single `@click` that dispatches the close event). When set,
	// the override is fully responsible for closing the panel — keep
	// the `popui-sidepanel-close` dispatch in your handler.
	CloseAttributes templ.Attributes
}
