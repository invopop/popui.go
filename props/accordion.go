package props

import "github.com/a-h/templ"

// Accordion provides props for a collapsible details element. Title renders
// in the trigger and children are the collapsible content.
type Accordion struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Title      string
	Open       bool
	// Variant is "default", or "panel": a panel of a page's side column,
	// with a compact header divided from its content, a plain icon, and a
	// shadow that sets it apart from a white work area.
	Variant string `jsonschema:"enum=,enum=default,enum=panel"`

	// Icon goes before the title, as an action's icon in a step.
	Icon templ.Component
	// Actions go after the title, before the chevron: a duration, a status
	// or a button. Clicking a button or link there doesn't open or close
	// the accordion.
	Actions templ.Component
}

// AccordionGroup provides props for a container that visually joins several
// accordions into one bordered panel with dividers between items.
type AccordionGroup struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
