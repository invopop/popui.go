package props

import "github.com/a-h/templ"

// App defines the properties used to configure an App.
type App struct {
	Title       string
	Description string

	// Accent color for the application.
	AccentColor string

	// Head allows for additional headers to be added to the App's Head
	// component.
	Head templ.Component

	// Data adds the x-data attribute to the application body wrapper's
	// contents. When used with HTMX, this data will also be replaced.
	Data string

	// HTMX when true loads the htmx javascript library
	HTMX bool

	// Axios when true loads the axios javascript library
	Axios bool

	// Auth when true enables authentication token handling in included scripts.
	Auth bool

	// SkipCSS when true will not include the default CSS in the head. This is useful
	// for applications that provide their own CSS on top of those used by PopUI.
	SkipCSS bool

	// Importmap is a list of module specifier mappings for an ES import map.
	// When non-empty, rendered before module and regular scripts.
	Importmap []Import

	// Modules is a list of ES module scripts to include.
	// Rendered as <script type="module"> after the importmap.
	Modules []Module

	// Scripts is a list of additional script paths to include as an alternative
	// to defining in the head property.
	Scripts []Script
	// Stylesheets is a list of additional stylesheet links to include as an
	// alternative to defining in the head property.
	Stylesheets []Link
}

// Main for configuring the App main content component.
type Main struct {
	ID    string
	Class string

	// Cloak when true adds the x-cloak attribute to the main element.
	Cloak bool

	// Data adds the x-data attribute to the main element.
	Data string

	// Center when true will ensure that the main content is centered in
	// the available space with a maximum width.
	Center bool

	Attributes templ.Attributes
}

// Article for configuring the App article component.
type Article struct {
	ID    string
	Class string

	// Icon is a URL for an icon to show before the title.
	Icon string
	// Title is a simple string title of this article.
	Title string
	// Subtitle is a simple string subtitle of this article.
	Subtitle string

	// FullWidth when true makes the article take the full width of the
	// main content area while maintaining padding.
	FullWidth bool

	// HideTitleSeparator when true hides the separator line between the title,
	// subtitle, and content. The line will only be present if the Title is set.
	HideTitleSeparator bool

	Attributes templ.Attributes
}

// Section defines a section within an article that has its own title, subtitle, and content.
// Sections are used to break up content within an article into smaller, more manageable pieces.
type Section struct {
	ID    string
	Class string

	// Icon is a URL for an icon to show before the title.
	Icon string
	// Title is a simple string title of this section.
	Title string
	// Subtitle is a simple string subtitle of this section.
	Subtitle string

	Attributes templ.Attributes
}

// Block for configuring the App block component.
type Block struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// Header for configuring the App header component and set of
// breadcrumbs if required.
type Header struct {
	ID    string
	Class string

	// Title allows for a custom title component to be provided on the left
	// side of the header before the Breadcrumbs. This can be used for logos,
	// texts, and a closer button.
	Title templ.Component

	// Data adds the x-data attribute to the header element.
	Data string

	// Attributes adds additional attributes to the header element.
	Attributes templ.Attributes

	// Breadcrumbs is the set of breadcrumbs to show in the header.
	Breadcrumbs []Breadcrumb
}

// Footer for configuring the App footer component.
type Footer struct {
	ID    string
	Class string
	// Data adds the x-data attribute to the footer element.
	Data       string
	Attributes templ.Attributes
}

// Nav for configuring the App nav component.
type Nav struct {
	ID    string
	Class string
	// Data adds the x-data attribute to the nav element.
	Data       string
	Attributes templ.Attributes
}

// Aside for configuring the App aside component.
type Aside struct {
	ID    string
	Class string
	// Data adds the x-data attribute to the aside element.
	Data       string
	Attributes templ.Attributes
}
