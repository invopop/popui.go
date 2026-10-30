package jsonui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/a-h/templ"
)

// Names of the built-in actions.
const (
	actionAlert    = "alert"
	actionSet      = "set"
	actionReset    = "reset"
	actionDispatch = "dispatch"
	actionValidate = "validate"
)

// BuiltinActions are the actions implemented by assets/jsonui.js.
var BuiltinActions = []string{actionAlert, actionSet, actionReset, actionDispatch, actionValidate}

// RequiredParams lists the string params each built-in action needs.
// Actions registered by the host have their params passed through unchecked.
var RequiredParams = map[string][]string{
	actionAlert:    {"message"},
	actionSet:      {"path"},
	actionDispatch: {"event"},
}

// Options control how documents are checked when they are built.
//
// By default documents are strict: the only way they can make the page do
// anything is through the named actions in their "on" handlers, so a host
// can embed documents from another source and decide exactly what they may
// trigger. Strict mode rejects:
//
//   - attributes in props that Alpine, htmx or the browser would run as
//     script: x-*, @*, :*, on*, hx-* and data-hx-* (use "model" and "on"
//     instead), including props set with "bind";
//   - text "bind" values other than plain data paths such as "customer.name";
//   - data paths in "model", and in the path, into and busy action params,
//     that reach __proto__, prototype or constructor;
//   - props that components use as Alpine variable names (such as the tabs
//     "model") unless they are plain identifiers;
//   - URLs with a scheme other than http, https, mailto or tel (e.g.
//     javascript:), in the href, src, action and formaction attributes and
//     in URL props (see URLFormat);
//   - actions not listed in Actions.
type Options struct {
	// Actions lists the action names documents may use in "on" handlers.
	// Nil allows BuiltinActions; an empty slice allows none. Hosts that
	// register their own actions with popui.actions.register list them here
	// too, e.g. append(jsonui.BuiltinActions, "ui.notify", "modal.close").
	Actions []string
	// Trusted turns off the strict checks, allowing Alpine directives and
	// event attributes in props, Alpine expressions in text binds and any
	// action name. Only use it for documents from a fully trusted source.
	Trusted bool
}

// Build parses a strict JSON document into a renderable component using the
// default Options.
func Build(data []byte) (templ.Component, error) {
	return Options{}.Build(data)
}

// Render parses a strict JSON document using the default Options and writes
// the resulting HTML to w.
func Render(ctx context.Context, data []byte, w io.Writer) error {
	return Options{}.Render(ctx, data, w)
}

// Build parses a JSON document into a renderable component.
func (o Options) Build(data []byte) (templ.Component, error) {
	var doc Document
	if err := decodeStrict(data, &doc); err != nil {
		return nil, err
	}
	return doc.Component(o)
}

// Render parses the JSON document and writes the resulting HTML to w.
func (o Options) Render(ctx context.Context, data []byte, w io.Writer) error {
	c, err := o.Build(data)
	if err != nil {
		return err
	}
	return c.Render(ctx, w)
}

func (o Options) allowsAction(name string) bool {
	if o.Trusted {
		return true
	}
	if o.Actions == nil {
		return slices.Contains(BuiltinActions, name)
	}
	return slices.Contains(o.Actions, name)
}

// eventName matches DOM event names with optional Alpine modifiers, e.g.
// "click", "submit.prevent" or "keydown.enter".
var eventName = regexp.MustCompile(`^[a-z][a-z0-9-]*(\.[a-z0-9-]+)*$`)

// URLFormat is the JSON Schema format of URL props. String props in the
// props package that hold URLs are tagged `jsonschema:"format=uri-reference"`
// and templ.SafeURL props are URLs by type; strict mode checks both.
const URLFormat = "uri-reference"

// urlAttributes are the HTML attributes whose values are URLs, in lower case.
var urlAttributes = []string{"href", "src", "action", "formaction"}

// urlKeys returns the lower case names of the URL attributes and of every
// URL prop in the registered components' props, including nested structs
// such as Input.Tooltip.
var urlKeys = sync.OnceValue(func() []string {
	keys := slices.Clone(urlAttributes)
	seen := map[reflect.Type]bool{}
	var walk func(t reflect.Type)
	walk = func(t reflect.Type) {
		for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice || t.Kind() == reflect.Array || t.Kind() == reflect.Map {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct || seen[t] {
			return
		}
		seen[t] = true
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Type == safeURLType || slices.Contains(strings.Split(f.Tag.Get("jsonschema"), ","), "format="+URLFormat) {
				if name := strings.ToLower(f.Name); !slices.Contains(keys, name) {
					keys = append(keys, name)
				}
			}
			walk(f.Type)
		}
	}
	for _, e := range registry {
		walk(e.props)
	}
	return keys
})

var safeURLType = reflect.TypeFor[templ.SafeURL]()

// checkProps applies the strict checks to a node's own props, before the
// renderer adds the x-model and x-on attributes generated from "model" and
// "on".
func checkProps(typ string, raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return fmt.Errorf("%s props: %w", typ, err)
	}
	if err := checkIdentifiers(typ, v); err != nil {
		return err
	}
	return checkValue(typ+" props", v)
}

// identifierProps lists props that components paste into Alpine expressions
// as variable names, which must therefore be plain identifiers.
var identifierProps = map[string][]string{
	"tabs": {"model"},
	"tab":  {"model"},
}

var identifier = regexp.MustCompile(`^[A-Za-z_$][\w$]*$`)

func checkIdentifiers(typ string, v any) error {
	props, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	for _, name := range identifierProps[typ] {
		for k, val := range props {
			if !strings.EqualFold(k, name) {
				continue
			}
			if s, _ := val.(string); s != "" && !identifier.MatchString(s) {
				return fmt.Errorf("%s props.%s: %q must be a plain identifier", typ, k, s)
			}
		}
	}
	return nil
}

func checkValue(at string, v any) error {
	switch v := v.(type) {
	case map[string]any:
		for k, val := range v {
			key := strings.ToLower(k)
			if key == "attributes" {
				if err := checkAttributes(at+"."+k, val); err != nil {
					return err
				}
			}
			if s, ok := val.(string); ok && slices.Contains(urlKeys(), key) && !safeURL(s) {
				return fmt.Errorf("%s.%s: unsafe URL %q", at, k, s)
			}
			if err := checkValue(at+"."+k, val); err != nil {
				return err
			}
		}
	case []any:
		for i, val := range v {
			if err := checkValue(fmt.Sprintf("%s[%d]", at, i), val); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkAttributes(at string, v any) error {
	attrs, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	for k := range attrs {
		name := strings.ToLower(k)
		if strings.HasPrefix(name, "x-") || strings.HasPrefix(name, "@") ||
			strings.HasPrefix(name, ":") || strings.HasPrefix(name, "on") ||
			strings.HasPrefix(name, "hx-") || strings.HasPrefix(name, "data-hx-") {
			return fmt.Errorf("%s: attribute %q is not allowed in strict mode; use \"model\" or \"on\"", at, k)
		}
	}
	return nil
}

// safeURL reports whether a URL can't run script when followed.
func safeURL(s string) bool {
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil {
		return false
	}
	switch strings.ToLower(u.Scheme) {
	case "", "http", "https", "mailto", "tel":
		return true
	}
	return false
}
