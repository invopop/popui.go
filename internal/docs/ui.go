// Package docs provides user interface components and utilities for PopUI.
package docs

import (
	"sort"
	"strings"

	"github.com/a-h/templ"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/popui.go/internal/docs/components"
)

// Page defines a single documentation page that can be rendered.
type Page struct {
	Title    string
	Desc     string
	Path     string
	Template templ.Component
}

// Group defines a group of related pages.
type Group struct {
	Title string
	Path  string
	Pages []*Page
}

// DocsIndex is the index of all documentation pages.
var groups = []*Group{
	{
		Title: "Foundations",
		Path:  "foundations",
		Pages: []*Page{
			{
				Title:    "Icons",
				Desc:     "Browse the full set of icons available in PopUI. Click any icon to copy its usage code.",
				Path:     "icons",
				Template: components.Icons(),
			},
			{
				Title:    "Tokens",
				Desc:     "Tailwind CSS design tokens including colors, spacing, shadows, and typography.",
				Path:     "tokens",
				Template: components.Tokens(),
			},
		},
	},
	{
		Title: "Components",
		Path:  "components",
		Pages: []*Page{
			{
				Title: "App",
				Desc: here.Doc(`
					The main application layout component with header, nav, aside, main,
					and footer sections. Use this as the base for all popui applications,
					and extend as needed.
				`),
				Path:     "app",
				Template: components.App(),
			},
			{
				Title:    "Card",
				Desc:     "Flexible container components for displaying content in a bordered box.",
				Path:     "card",
				Template: components.Card(),
			},
			{
				Title:    "PageState",
				Desc:     "Layout for displaying page states with illustrations and call-to-action.",
				Path:     "page-state",
				Template: components.PageState(),
			},
			{
				Title:    "CardDeck",
				Desc:     "Tinted container that stacks regular Cards into one visual unit.",
				Path:     "card-deck",
				Template: components.CardDeck(),
			},
			{
				Title:    "Separator",
				Desc:     "Horizontal divider with a dashed line for visual separation.",
				Path:     "separator",
				Template: components.Separator(),
			},
			{
				Title:    "Accordion",
				Desc:     "Vertically stacked interactive sections to organize content.",
				Path:     "accordion",
				Template: components.Accordion(),
			},
			{
				Title:    "Avatar",
				Desc:     "Display user profile images or initials in circular containers.",
				Path:     "avatar",
				Template: components.Avatar(),
			},
			{
				Title:    "Breadcrumbs",
				Desc:     "Navigation breadcrumbs to show the current page location.",
				Path:     "breadcrumb",
				Template: components.Breadcrumbs(),
			},
			{
				Title:    "Button",
				Desc:     "Trigger actions and events with customizable button components.",
				Path:     "button",
				Template: components.Button(),
			},
			{
				Title:    "Calendar",
				Desc:     "Dual-month, range-selection calendar with a preset rail and month navigation. Backed by the rangeCalendar Alpine controller.",
				Path:     "calendar",
				Template: components.Calendar(),
			},
			{
				Title:    "Checkbox",
				Desc:     "Checkbox inputs with optional toggle switch variant.",
				Path:     "checkbox",
				Template: components.Checkbox(),
			},
			{
				Title:    "Menu",
				Desc:     "A dropdown menu of actions that opens from a trigger button.",
				Path:     "menu",
				Template: components.Menu(),
			},
			{
				Title:    "DescriptionList",
				Desc:     "Semantic HTML definition list for displaying term-description pairs.",
				Path:     "description-list",
				Template: components.DescriptionList(),
			},
			{
				Title:    "Fieldset",
				Desc:     "Groups form fields together with proper spacing and optional legend.",
				Path:     "fieldset",
				Template: components.Fieldset(),
			},
			{
				Title:    "File",
				Desc:     "File components for selecting, uploading and displaying files. Use InputFile for basic file selection, FileUpload for avatar/image uploads with preview, and FileDownload for displaying stored files.",
				Path:     "file",
				Template: components.File(),
			},
			{
				Title:    "Filter",
				Desc:     "Search-bar style filter for data views: a \"+ Filter\" menu plus one editable chip per active filter, with multi-field filtering, colored option lists, and keyboard navigation.",
				Path:     "filter",
				Template: components.Filter(),
			},
			{
				Title:    "Form",
				Desc:     "Form element with proper spacing and standard HTML form attributes for handling submissions.",
				Path:     "form",
				Template: components.Form(),
			},
			{
				Title:    "Image",
				Desc:     "Displays images with rounded corners and proper object fit.",
				Path:     "image",
				Template: components.Image(),
			},
			{
				Title:    "Input",
				Desc:     "Text input field for capturing user data with various types and validation.",
				Path:     "input",
				Template: components.Input(),
			},
			{
				Title:    "Label",
				Desc:     "Form label element with optional hint tooltip.",
				Path:     "label",
				Template: components.Label(),
			},
			{
				Title:    "Notification",
				Desc:     "Feedback messages with different severity types and icons.",
				Path:     "notification",
				Template: components.Notification(),
			},
			{
				Title:    "Popover",
				Desc:     "A centered modal dialog using the HTML popover API with backdrop overlay.",
				Path:     "popover",
				Template: components.Popover(),
			},
			{
				Title:    "Radio",
				Desc:     "Radio button inputs for selecting a single option from a group.",
				Path:     "radio",
				Template: components.Radio(),
			},
			{
				Title:    "Select",
				Desc:     "Dropdown selection control for choosing one or more options from a list.",
				Path:     "select",
				Template: components.Select(),
			},
			{
				Title:    "SidePanel",
				Desc:     "Floating side panel that overlays one edge of the viewport.",
				Path:     "side-panel",
				Template: components.SidePanel(),
			},
			{
				Title:    "SignaturePad",
				Desc:     "A signature capture dialog with typed-name preview and freehand drawing tabs.",
				Path:     "signature-pad",
				Template: components.SignaturePad(),
			},
			{
				Title:    "Slider",
				Desc:     "A range slider input for selecting numeric values.",
				Path:     "slider",
				Template: components.Slider(),
			},
			{
				Title:    "Table",
				Desc:     "Display data in a structured table format with automatic styling for headers, cells, and borders.",
				Path:     "table",
				Template: components.Table(),
			},
			{
				Title:    "Tabs",
				Desc:     "Interactive tab navigation component with default and pill variants.",
				Path:     "tabs",
				Template: components.Tabs(),
			},
			{
				Title:    "TagStatus",
				Desc:     "Status indicators with optional dots and different color variants.",
				Path:     "tag-status",
				Template: components.TagStatus(),
			},
			{
				Title:    "StatusBadge",
				Desc:     "Icon-led outcome pill with an optional label: success, failed, warning or running.",
				Path:     "status-badge",
				Template: components.StatusBadge(),
			},
			{
				Title:    "Textarea",
				Desc:     "Multi-line text input field for capturing longer user input with support for labels and validation.",
				Path:     "textarea",
				Template: components.Textarea(),
			},
			{
				Title:    "Toast",
				Desc:     "Dark floating notification with a type icon, message, optional description, and optional action button.",
				Path:     "toast",
				Template: components.Toast(),
			},
			{
				Title:    "Tooltip",
				Desc:     "Dark floating card with a title, description and optional illustration, revealed on hover.",
				Path:     "tooltip",
				Template: components.Tooltip(),
			},
			{
				Title:    "Typography",
				Desc:     "Text components for headings, paragraphs, descriptions, and formatted content.",
				Path:     "typography",
				Template: components.Typography(),
			},
		},
	},
}

// Pages within each group are presented alphabetically regardless of the
// order they are added in above; the groups themselves keep their authored
// order. The Get Started guide is pinned above the groups by the sidebar
// template.
func init() {
	for _, g := range groups {
		sort.Slice(g.Pages, func(i, j int) bool {
			return strings.ToLower(g.Pages[i].Title) < strings.ToLower(g.Pages[j].Title)
		})
	}
}
