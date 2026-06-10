package props

import "github.com/a-h/templ"

// Drawer position constants for the Position prop. Drawers anchor to
// one edge of the viewport; pick the side that's furthest from the
// app's primary navigation.
const (
	DrawerPositionRight string = "right"
	DrawerPositionLeft  string = "left"
)

// Drawer renders a floating, fixed-position side panel that overlays
// one edge of the viewport. Non-blocking by default — no backdrop, the
// rest of the app stays interactive — and stays mounted across HTMX
// content swaps so the open/close state survives. Consumers typically
// host a swap target inside the drawer's children block, fill it via
// hx-target, and dispatch the open event after the swap.
//
// Toggled via the `popui-drawer-open` / `popui-drawer-close` window
// events whose `event.detail` must equal the drawer's ID. The
// controller scopes those by ID so several drawers can coexist on a
// single page without trampling each other; it also re-broadcasts the
// matching event whenever the open state flips, so external listeners
// react to every close path (X click, Escape, programmatic).
type Drawer struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Position picks the viewport edge the drawer anchors to.
	// One of DrawerPositionRight (default) or DrawerPositionLeft.
	// Empty / unknown values fall back to right.
	Position string

	// Width sets the drawer's fixed pixel width. Zero falls back to
	// 400 — a comfortable detail-panel default. Anything between
	// 320 and 720 reads well without crowding the main content.
	Width int
}

// DrawerHeader renders a sticky header bar at the top of a popui.Drawer:
// an X close button on the leading edge, a single-line truncated title
// in the middle, and an optional action slot on the trailing edge
// (children of the component). Lives inside a popui.Drawer's children
// block so it sits above the scrollable content area.
//
// The close button dispatches `popui-drawer-close` (event.detail =
// DrawerID) by default. Pass CloseAttributes to override that click
// behavior — e.g. to also clear a URL query param via
// history.replaceState — but make sure the override still dispatches
// the close event so the drawer's Alpine state syncs.
type DrawerHeader struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// DrawerID is the parent Drawer's ID — used by the default close
	// button to dispatch `popui-drawer-close` with the matching
	// `event.detail` so the right drawer toggles closed.
	DrawerID string

	// Title is the h1 heading rendered in the centre cell. Truncates
	// with `truncate` if it overflows the available width.
	Title string

	// CloseAttributes replace the close button's default attributes
	// (a single `@click` that dispatches the close event). When set,
	// the override is fully responsible for closing the drawer — keep
	// the `popui-drawer-close` dispatch in your handler.
	CloseAttributes templ.Attributes
}
