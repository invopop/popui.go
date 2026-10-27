package jsonui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/a-h/templ"
)

// BuiltinActions are the actions implemented by assets/jsonui.js.
var BuiltinActions = []string{"alert", "set", "reset", "dispatch"}

// Options control how documents are checked when they are built.
//
// By default documents are strict: the only way they can make the page do
// anything is through the named actions in their "on" handlers, so a host
// can embed documents from another source and decide exactly what they may
// trigger. Strict mode rejects:
//
//   - attributes in props that Alpine or the browser would run as script:
//     x-*, @*, :* and on* (use "model" and "on" instead);
//   - text "bind" values other than plain data paths such as "customer.name";
//   - props that components use as Alpine variable names (such as the tabs
//     "model") unless they are plain identifiers;
//   - href, src, action, formaction and url values with a scheme other than
//     http, https, mailto or tel (e.g. javascript:);
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

// urlKeys are the props and attribute names whose values are URLs.
var urlKeys = []string{"href", "src", "action", "formaction", "url"}

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
			if s, ok := val.(string); ok && slices.Contains(urlKeys, key) && !safeURL(s) {
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
			strings.HasPrefix(name, ":") || strings.HasPrefix(name, "on") {
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
