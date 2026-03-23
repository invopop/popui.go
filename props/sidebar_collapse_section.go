package props

import "github.com/a-h/templ"

// SidebarCollapseSection Templ component props
type SidebarCollapseSection struct {
	Title       string
	Class       string
	Collapsable bool
	Attributes  templ.Attributes
}
