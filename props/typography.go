package props

import "github.com/a-h/templ"

// P defines the props for the P component.
type P struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Size of the text, see the TextSize constants. Default is "base".
	Size string `jsonschema:"enum=,enum=sm,enum=base,enum=lg,enum=xl"`
	// Tone is the colour of the text, see the TextTone constants. Default
	// is "default".
	Tone string `jsonschema:"enum=,enum=default,enum=secondary,enum=tertiary"`
	// Numeric shows digits in equal widths so figures line up, for amounts.
	Numeric bool
}

// Title defines the props for the Title component.
type Title struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Size of the title: "lg" for a page or card, or "base" for a section
	// within one, in a medium weight. Default is "lg".
	Size string `jsonschema:"enum=,enum=base,enum=lg"`
}

// TitleGroup defines the props for the TitleGroup component: a title with
// an optional secondary description, stacked by default or arranged on one
// baseline with Horizontal.
type TitleGroup struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Title       string
	Description string
	// Horizontal places the description next to the title on the same
	// baseline, separated by a middot, instead of stacking it below.
	Horizontal bool
}

// Subtitle defines the props for the Subtitle component.
type Subtitle struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// Description defines the props for the Description component.
type Description struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Size of the text, "sm" or "base". Default is "base".
	Size string `jsonschema:"enum=,enum=sm,enum=base"`
}

// Info defines the props for the Info component.
type Info struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// Warning defines the props for the Warning component.
type Warning struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// Mark defines the props for the Mark component.
type Mark struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// Prose defines the props for the Prose component.
type Prose struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
