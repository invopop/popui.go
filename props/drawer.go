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
}
