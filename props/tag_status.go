package props

import "github.com/a-h/templ"

// TagStatus Templ component props
type TagStatus struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Status is the colour of the tag.
	Status string `jsonschema:"enum=grey,enum=green,enum=yellow,enum=red,enum=orange,enum=blue,enum=purple,enum=olive,enum=teal,enum=crimson,enum=steelBlue,enum=empty"`
	Dot    bool
}
