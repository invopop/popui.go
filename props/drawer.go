package props

import "github.com/a-h/templ"

// Drawer renders a floating, fixed-position side panel that overlays
// the right edge of the viewport. Mounted once near the App root and
// toggled open/closed via the `popui-drawer-open` / `popui-drawer-close`
// window events (event.detail must equal the drawer's ID). Stays
// mounted across HTMX content swaps so the open/close state survives
// — the caller fills the drawer's content slot via hx-target on a
// child container and dispatches the open event after the swap.
//
// Mirrors console-ui's JobDetailPanel pattern: fixed right-0 top-0
// bottom-0 z-50, non-blocking (no backdrop, rest of the app stays
// interactive), `fly`-style transition in from the right, closes on
// Escape. No grid reflow on open/close — the panel overlays the main
// content area instead of pushing it.
type Drawer struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Width sets the drawer's fixed pixel width. Zero falls back to
	// 400 — the JobDetailPanel default. Anything between 320 and
	// 720 reads comfortably without overlapping the main content.
	Width int
}

// DrawerHeader renders a sticky header bar at the top of a popui.Drawer:
// an X close button on the left, a single-line truncated title in the
// middle, and an optional action slot on the right (children of the
// component). Lives inside a popui.Drawer's children block so it sits
// above the scrollable content area.
//
// The close button dispatches `popui-drawer-close` (event.detail =
// DrawerID) by default. Pass CloseAttributes to add extra behavior on
// click — e.g. clearing a URL query param via history.replaceState.
// CloseAttributes are appended after the default attributes, so any
// extra `@click` REPLACES the default rather than merging — give the
// override responsibility for re-dispatching the close event.
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

	// CloseAttributes override the close button's default attributes
	// (a single `@click` that dispatches the close event). Pass your
	// own `@click` if you need extra logic (e.g. URL rewrite) — make
	// sure to still dispatch the close event yourself.
	CloseAttributes templ.Attributes
}
