package props

import "github.com/a-h/templ"

// SidebarVariant constants
const (
	SidebarVariantDark string = "dark"
)

// Sidebar Templ component props
type Sidebar struct {
	Title      string
	Class      string
	Variant    string
	Attributes templ.Attributes
}

// SidebarSection Templ component props
type SidebarSection struct {
	Title      string
	Class      string
	Attributes templ.Attributes
}

// SidebarItem defines the property for a single sidebar navigation item.
type SidebarItem struct {
	ID         string
	Class      string
	Selected   bool
	Attributes templ.Attributes
	Href       templ.SafeURL
}

// SidebarGroup nests items under a named header row — an organization and
// its areas, a project and its pages. Children render indented behind a
// tree line.
type SidebarGroup struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label is the header row's text.
	Label string
	// Header replaces Label with arbitrary content (an Avatar and a name,
	// say) when set.
	Header templ.Component
	// HeaderAttributes attach behaviour to the header row — an href, an
	// hx-get, an @click. When present the row renders as a link with a
	// trailing chevron; otherwise it is a plain heading.
	HeaderAttributes templ.Attributes
	// HeaderMenu turns the header row into a Menu trigger opening this
	// component (MenuItems) — a switcher between sibling groups, say. Takes
	// precedence over HeaderAttributes.
	HeaderMenu templ.Component
}

// SidebarUser is the person signed in, pinned to the sidebar's footer:
// avatar, name and email, opening a menu of account actions (the children,
// as MenuItems).
type SidebarUser struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Name is the user's display name.
	Name string
	// Email renders muted under the name.
	Email string
	// AvatarURL is the user's picture; without it the avatar falls back to
	// the name's initial.
	AvatarURL string
}
