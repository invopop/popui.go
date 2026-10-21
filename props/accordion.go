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
}

// AccordionGroup provides props for a container that visually joins several
// accordions into one bordered panel with dividers between items.
type AccordionGroup struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
