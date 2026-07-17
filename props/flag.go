package props

import "github.com/a-h/templ"

// Flag Templ component props.
//
// Deprecated: use Image with https://assets.invopop.com/flags/<cc>.svg instead.
type Flag struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Country    string // ISO 3166-1 alpha-2 country code (e.g., "ES", "US", "GB")
}
