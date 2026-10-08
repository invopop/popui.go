package props

import "github.com/a-h/templ"

// Accordion provides props for the details element wrapper.
//
// When Title is set the accordion renders its own trigger (title, optional
// description and chevron) and treats children as the collapsible content,
// mirroring Mintlify's <Accordion title="..."> usage. Leave Title empty to
// compose the trigger and content manually with AccordionTrigger and
// AccordionContent.
type Accordion struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Open        bool
	Title       string
	Description string
}

// AccordionGroup provides props for a container that visually joins several
// accordions into one bordered panel with dividers between items.
type AccordionGroup struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// AccordionTrigger provides props for the clickable summary element
type AccordionTrigger struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Description string
}

// AccordionContent provides props for the collapsible content area
type AccordionContent struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
