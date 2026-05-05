package props

import (
	"encoding/json"

	"github.com/a-h/templ"
)

// Head Templ component properties.
type Head struct {
	Title       string
	Description string

	// AlpineJS when true includes Alpine.js in the head.
	AlpineJS bool

	// Auth when true enables authentication token handling in included scripts.
	Auth bool

	// HTMX when true loads the htmx library including special authentication
	// token handling.
	HTMX bool
	// Axios when true loads the axios javascript library alongside an interceptor
	// to add authentication tokens to requests automatically.
	Axios bool

	// SkipCSS when true will not include the default CSS in the head. This is useful
	// when you want to provide your own CSS or use the components without the default
	// styling.
	SkipCSS bool

	// Importmap is a list of module specifier mappings. When non-empty,
	// rendered as <script type="importmap">{"imports":{...}}</script>
	// before any module scripts. Order is preserved.
	Importmap []Import

	// Modules is a list of ES module scripts to include.
	// Rendered as <script type="module" src="..."></script> after
	// the importmap (if any) and before regular Scripts.
	Modules []Module

	// Scripts is a list of additional script paths to include.
	Scripts []Script
	// Stylesheets is a list of additional stylesheet links to include.
	Stylesheets []Link
}

// Body Templ component props.
type Body struct {
	ID    string
	Class string

	// Data adds the x-data attribute to the body element.
	Data string

	Attributes templ.Attributes
}

// PopupLayout Templ component props
type PopupLayout struct {
	Title string
}

// Import defines a single module specifier mapping for an import map.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script/type/importmap
type Import struct {
	// Name is the bare module specifier (e.g. "codemirror", "@codemirror/state").
	Name string
	// URL is the resolved URL for the module (e.g. "https://esm.sh/codemirror@6.0.1").
	URL string
}

// importmapData is the JSON structure for an ES module import map.
type importmapData struct {
	Imports [][2]string `json:"-"`
}

// MarshalJSON produces ordered JSON: {"imports":{"name":"url",...}}
func (d importmapData) MarshalJSON() ([]byte, error) {
	buf := []byte(`{"imports":{`)
	for i, pair := range d.Imports {
		if i > 0 {
			buf = append(buf, ',')
		}
		key, err := json.Marshal(pair[0])
		if err != nil {
			return nil, err
		}
		val, err := json.Marshal(pair[1])
		if err != nil {
			return nil, err
		}
		buf = append(buf, key...)
		buf = append(buf, ':')
		buf = append(buf, val...)
	}
	buf = append(buf, "}}"...)
	return buf, nil
}

// ImportmapScript returns a templ.Component that renders the Import slice
// as a <script type="importmap"> tag with ordered JSON content.
func ImportmapScript(imports []Import) templ.Component {
	pairs := make([][2]string, len(imports))
	for i, im := range imports {
		pairs[i] = [2]string{im.Name, im.URL}
	}
	return templ.JSONScript("", importmapData{Imports: pairs}).WithType("importmap")
}

// Module defines an ES module script to include in the head section.
// These are rendered as <script type="module" src="..."></script>.
type Module struct {
	Src string
}

// Script defines a script to include in the head section.
type Script struct {
	Src   string
	Defer bool
	Async bool
}

// Link defines a link to include in the head section.
type Link struct {
	Href string
	Rel  string // default rel will be "stylesheet"
}

// RelOrStylesheet returns the rel or "stylesheet" value if not set.
func (l Link) RelOrStylesheet() string {
	if l.Rel == "" {
		return "stylesheet"
	}
	return l.Rel
}
