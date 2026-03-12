package props

import "github.com/a-h/templ"

// WizardHeader Templ component props
type WizardHeader struct {
	Class      string
	Attributes templ.Attributes
}

// WizardContent Templ component props
type WizardContent struct {
	Class          string
	CenterVertical bool
	Attributes     templ.Attributes
}

// WizardFooter Templ component props
type WizardFooter struct {
	Class      string
	Attributes templ.Attributes
}
