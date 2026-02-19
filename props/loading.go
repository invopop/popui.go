package props

import "github.com/a-h/templ"

// Loading variants
const (
	LoadingSkeleton string = "skeleton"
	LoadingSpinner  string = "spinner"
)

// Loading sizes
const (
	LoadingSizeXS     string = "xs"
	LoadingSizeSmall  string = "sm"
	LoadingSizeMedium string = "md"
	LoadingSizeLarge  string = "lg"
	LoadingSizeXL     string = "xl"
)

// Loading provides props for the loading component
type Loading struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Variant defines the loading style: "skeleton" (default) or "spinner".
	Variant string

	// Size defines the loading indicator height, see the LoadingSize constants.
	// Default is "md". Only applies to skeleton variant.
	Size string

	// Width defines the skeleton width as a percentage (0-100).
	// Default is 100 (full width). Only applies to skeleton variant.
	Width int

	// Rounded applies full border-radius for circular/pill shapes.
	// Default is false (uses standard border-radius). Only applies to skeleton variant.
	Rounded bool

	// Label is optional text to display alongside the spinner.
	// Only applies to spinner variant.
	Label string
}
