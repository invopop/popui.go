// Package jsonui renders PopUI components from a JSON description of the
// component tree. It is EXPERIMENTAL: the same code is compiled natively for
// server-side rendering and to WebAssembly (see ./wasm) so a browser can turn
// JSON straight into PopUI HTML without a Go server round-trip.
//
// A document is an object with a template of nodes and optional Alpine.js
// data, which is also the source for bound values:
//
//	{
//	  "$schema": "https://…/jsonui/schema/document.json",
//	  "data": {"customer": {"name": "Acme S.L."}},
//	  "template": [
//	    {
//	      "type": "card",
//	      "children": [
//	        {"type": "card_header", "bind": {"title": "customer.name"}},
//	        {"type": "text", "bind": "customer.name"}
//	      ]
//	    }
//	  ]
//	}
//
// "props" is decoded straight into the matching struct from the props
// package (field names are matched case-insensitively, unknown fields are
// rejected). Props that hold templ.Component values (icons, etc.) are not
// supported yet. The JSON Schema in ./schema is generated from the same
// registry and props structs.
package jsonui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	popui "github.com/invopop/popui.go"
)

// Document is the root of a JSON UI description.
type Document struct {
	// Schema identifies the JSON Schema of the document so editors can offer
	// completion and validation. It is ignored when rendering.
	Schema string `json:"$schema,omitempty"`
	// Data is a JSON object with the document's state. It becomes the
	// Alpine.js `x-data` of the document, is the source for node "bind" and
	// "model" values, and is what actions in "on" handlers read and change.
	// Behaviour is never embedded here: actions are implemented in JS (see
	// assets/jsonui.js) and only referenced by name.
	Data json.RawMessage `json:"data,omitempty"`
	// Template is the list of nodes to render.
	Template []Node `json:"template"`
	// Handlers are named actions, or lists of actions, that "on" handlers
	// can refer to by name so a handler used in several places is defined
	// once, e.g. {"update": {"action": "step.update", "summary": "…"}} and
	// then {"on": {"change": "update"}}. Handlers can't refer to each other.
	Handlers map[string]json.RawMessage `json:"handlers,omitempty"`
}

// Node is a single element in the JSON component tree.
type Node struct {
	// Type is the registered component name, e.g. "button" or "card_header".
	Type string `json:"type"`
	// Props is decoded into the component's props struct.
	Props json.RawMessage `json:"props,omitempty"`
	// Children are rendered as the component's `{ children... }`.
	Children []Node `json:"children,omitempty"`
	// Text is only used by the "text" node type.
	Text string `json:"text,omitempty"`
	// Bind takes values from the document data. On "text" nodes it is an
	// Alpine expression shown with x-text; a plain path such as
	// "customer.name" is also resolved from object data at render time so the
	// text is correct before Alpine starts. On components it maps prop names
	// to data paths, e.g. {"value": "customer.name"}, resolved at render time.
	Bind json.RawMessage `json:"bind,omitempty"`
	// Model two-way binds a form control (input, textarea, select, checkbox)
	// to a data path, e.g. "customer.email". The initial value is also
	// rendered from the data.
	Model string `json:"model,omitempty"`
	// On maps DOM events, with optional Alpine modifiers such as
	// "submit.prevent", to an action, the name of one of the document's
	// handlers, or a list of either to run in order, e.g.
	// {"click": {"action": "alert", "message": "Hi {{customer.name}}"}} or
	// {"change": "update"}.
	On map[string]json.RawMessage `json:"on,omitempty"`
}

// valueProp is the prop holding a text control's value.
const valueProp = "value"

// nodeEvents lists the "on" events accepted by nodes that don't take every
// DOM event.
var nodeEvents = map[string][]string{}

// modelProps lists the components that accept "model" and the prop holding
// their rendered value.
var modelProps = map[string]string{
	"input":    valueProp,
	"textarea": valueProp,
	"select":   "options",
	"checkbox": "checked",
}

// Component describes a registered node type.
type Component struct {
	// Type is the node type name used in JSON.
	Type string
	// Props is the props struct the node's "props" object decodes into.
	Props reflect.Type
	// Model is true for form controls that accept a "model" binding.
	Model bool
	// Leaf is true for nodes that take no "children" or "bind".
	Leaf bool
	// Events lists the events the node accepts in "on", or nil for any DOM
	// event.
	Events []string
}

// builder converts raw props into a component, children are attached later.
type builder func(raw json.RawMessage) (templ.Component, error)

type entry struct {
	props reflect.Type
	build builder
}

var registry = map[string]entry{}

// register adds a PopUI component whose signature follows the standard
// `templ Name(p ...props.X)` convention.
func register[P any](name string, fn func(...P) templ.Component) {
	registry[name] = entry{
		props: reflect.TypeFor[P](),
		build: func(raw json.RawMessage) (templ.Component, error) {
			var p P
			if len(raw) > 0 {
				if err := decodeStrict(raw, &p); err != nil {
					return nil, fmt.Errorf("%s props: %w", name, err)
				}
			}
			return fn(p), nil
		},
	}
}

func init() {
	// Layout & typography
	register("title", popui.Title)
	register("title_group", popui.TitleGroup)
	register("subtitle", popui.Subtitle)
	register("description", popui.Description)
	register("p", popui.P)
	register("separator", popui.Separator)

	// Cards & details
	register("card", popui.Card)
	register("card_header", popui.CardHeader)
	register("card_content", popui.CardContent)
	register("card_progress_bar", popui.CardProgressBar)
	register("description_list", popui.DescriptionList)
	register("description_list_item", popui.DescriptionListItem)
	register("dt", popui.DT)
	register("dd", popui.DD)

	// Status & identity
	register("status_badge", popui.StatusBadge)
	register("tag_status", popui.TagStatus)
	register("avatar", popui.Avatar)

	// Actions
	register("button", popui.Button)
	register("button_group", popui.ButtonGroup)

	// Forms
	register("form", popui.Form)
	register("fieldset", popui.Fieldset)
	register("input", popui.Input)
	register("textarea", popui.Textarea)
	register("select", popui.Select)
	register("checkbox", popui.Checkbox)
	register("option_group", popui.OptionGroup)

	// Navigation
	register("tabs", popui.Tabs)
	register("tab", popui.Tab)
}

// Components returns the registered component types sorted by name. The
// built-in "text" node is not included as it has no props.
func Components() []Component {
	out := make([]Component, 0, len(registry))
	for k, e := range registry {
		_, model := modelProps[k]
		// Entries without a props builder are custom leaf nodes (upload_file).
		out = append(out, Component{Type: k, Props: e.props, Model: model, Leaf: e.build == nil, Events: nodeEvents[k]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Type < out[j].Type })
	return out
}

// Types returns the sorted list of supported node types.
func Types() []string {
	out := []string{"text"}
	for _, c := range Components() {
		out = append(out, c.Type)
	}
	sort.Strings(out)
	return out
}

// Component builds the document's template. When the document has data,
// the template is wrapped in an Alpine.js component root.
func (d *Document) Component(opts Options) (templ.Component, error) {
	if d.Template == nil {
		return nil, errors.New("document: missing template")
	}
	b := &docBuilder{opts: opts}
	xdata := ""
	if len(d.Data) > 0 {
		var err error
		if xdata, b.data, err = parseData(d.Data); err != nil {
			return nil, err
		}
	}
	// Resolve every handler up front so mistakes are reported even when a
	// handler isn't used.
	b.handlers = make(map[string][]map[string]any, len(d.Handlers))
	for name, raw := range d.Handlers {
		list, err := b.resolve(raw, false)
		if err != nil {
			return nil, fmt.Errorf("handler %q: %w", name, err)
		}
		b.handlers[name] = list
	}
	body, err := b.buildAll(d.Template)
	if err != nil {
		return nil, err
	}
	if xdata == "" {
		return body, nil
	}
	return withChildren(documentRoot(xdata), body), nil
}

// parseData returns the compact x-data JSON for the document along with
// its decoded value for resolving binds and models.
func parseData(raw json.RawMessage) (string, map[string]any, error) {
	var v any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return "", nil, fmt.Errorf("document data: %w", err)
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return "", nil, errors.New("document data: must be an object")
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, raw); err != nil {
		return "", nil, fmt.Errorf("document data: %w", err)
	}
	return buf.String(), obj, nil
}

// docBuilder holds the state shared while building a document's nodes.
type docBuilder struct {
	opts Options
	// data is the document data when it is a JSON object, or nil.
	data map[string]any
	// handlers are the document's named handlers, resolved to their actions.
	handlers map[string][]map[string]any
}

func (b *docBuilder) buildAll(nodes []Node) (templ.Component, error) {
	comps := make([]templ.Component, len(nodes))
	for i := range nodes {
		c, err := b.build(&nodes[i])
		if err != nil {
			return nil, err
		}
		comps[i] = c
	}
	return templ.Join(comps...), nil
}

func (b *docBuilder) build(n *Node) (templ.Component, error) {
	if n.Type == "text" {
		return b.text(n)
	}
	if !b.opts.Trusted {
		if err := checkProps(n.Type, n.Props); err != nil {
			return nil, err
		}
	}
	if n.Type == uploadFileType {
		return b.uploadFile(n)
	}
	e, ok := registry[n.Type]
	if !ok {
		return nil, fmt.Errorf("unknown component type %q", n.Type)
	}
	raw, err := b.props(n)
	if err != nil {
		return nil, err
	}
	c, err := e.build(raw)
	if err != nil {
		return nil, err
	}
	if len(n.Children) == 0 {
		return c, nil
	}
	children, err := b.buildAll(n.Children)
	if err != nil {
		return nil, err
	}
	return withChildren(c, children), nil
}

func (b *docBuilder) text(n *Node) (templ.Component, error) {
	if len(n.Bind) == 0 {
		return templ.Raw(templ.EscapeString(n.Text)), nil
	}
	var expr string
	if err := json.Unmarshal(n.Bind, &expr); err != nil {
		return nil, errors.New("text bind: must be a string expression")
	}
	isPath := dataPath.MatchString(expr)
	if !isPath && !b.opts.Trusted {
		return nil, fmt.Errorf("text bind: %q must be a data path in strict mode", expr)
	}
	text := n.Text
	if b.data != nil && isPath {
		if v, err := lookup(b.data, expr); err == nil {
			text = display(v)
		}
	}
	return boundText(text, expr), nil
}

// props returns the node's props with any bind, model and on values from
// the node applied.
func (b *docBuilder) props(n *Node) (json.RawMessage, error) {
	if len(n.Bind) == 0 && n.Model == "" && len(n.On) == 0 {
		return n.Props, nil
	}
	var props map[string]json.RawMessage
	if len(n.Props) > 0 {
		if err := json.Unmarshal(n.Props, &props); err != nil {
			return nil, fmt.Errorf("%s props: %w", n.Type, err)
		}
	}
	if props == nil { // missing or null
		props = map[string]json.RawMessage{}
	}
	if len(n.Bind) > 0 {
		if err := b.bindProps(n.Type, props, n.Bind); err != nil {
			return nil, err
		}
		if !b.opts.Trusted {
			// Bound values come from the data, so check them like props.
			raw, err := json.Marshal(props)
			if err != nil {
				return nil, err
			}
			if err := checkProps(n.Type, raw); err != nil {
				return nil, err
			}
		}
	}
	var attrs map[string]any
	attrsKey := "attributes"
	for k, v := range props {
		if strings.EqualFold(k, attrsKey) {
			attrsKey = k
			if err := json.Unmarshal(v, &attrs); err != nil {
				return nil, fmt.Errorf("%s attributes: %w", n.Type, err)
			}
		}
	}
	if attrs == nil { // missing or null
		attrs = map[string]any{}
	}
	if n.Model != "" {
		if err := b.model(n.Type, props, attrs, n.Model); err != nil {
			return nil, err
		}
	}
	for event, raw := range n.On {
		handler, err := b.actions(event, raw)
		if err != nil {
			return nil, fmt.Errorf("%s on: %w", n.Type, err)
		}
		attrs["x-on:"+event] = handler
	}
	if len(attrs) > 0 {
		var err error
		if props[attrsKey], err = json.Marshal(attrs); err != nil {
			return nil, err
		}
	}
	return json.Marshal(props)
}

// bindProps overlays the props with values taken from the document data.
func (b *docBuilder) bindProps(typ string, props map[string]json.RawMessage, bind json.RawMessage) error {
	var paths map[string]string
	if err := decodeStrict(bind, &paths); err != nil {
		return fmt.Errorf("%s bind: must map prop names to data paths: %w", typ, err)
	}
	if b.data == nil {
		return fmt.Errorf("%s bind: document data must be set", typ)
	}
	for prop, path := range paths {
		v, err := lookup(b.data, path)
		if err != nil {
			return fmt.Errorf("%s bind %s: %w", typ, prop, err)
		}
		if props[prop], err = json.Marshal(v); err != nil {
			return err
		}
	}
	return nil
}

// model adds the x-model binding for a form control and renders its initial
// value from the data.
func (b *docBuilder) model(typ string, props map[string]json.RawMessage, attrs map[string]any, path string) error {
	prop, ok := modelProps[typ]
	if !ok {
		return fmt.Errorf("%s: model is only supported on form controls", typ)
	}
	if b.data == nil {
		return fmt.Errorf("%s model: document data must be set", typ)
	}
	if !safePath(path) {
		return fmt.Errorf("%s model: %q is not a data path", typ, path)
	}
	v, err := lookup(b.data, path)
	if err != nil {
		return fmt.Errorf("%s model: %w", typ, err)
	}
	directive := "x-model"
	if _, isNum := v.(json.Number); isNum {
		directive = "x-model.number" // keep numbers as numbers in the data
	}
	attrs[directive] = path
	switch prop {
	case valueProp:
		props[prop], err = json.Marshal(display(v))
	case "checked":
		on, _ := v.(bool)
		props[prop], err = json.Marshal(on)
	case "options":
		err = selectOptions(props, v)
	}
	return err
}

// selectOptions marks the select options with the given value as selected,
// or with any of the values when it is a list, as for a multiple select.
func selectOptions(props map[string]json.RawMessage, value any) error {
	var values []string
	if list, ok := value.([]any); ok {
		for _, v := range list {
			values = append(values, display(v))
		}
	} else {
		values = []string{display(value)}
	}
	for k, raw := range props {
		if !strings.EqualFold(k, "options") {
			continue
		}
		var opts []map[string]any
		if err := json.Unmarshal(raw, &opts); err != nil {
			return fmt.Errorf("select options: %w", err)
		}
		for _, o := range opts {
			for ok := range o {
				if strings.EqualFold(ok, "selected") {
					delete(o, ok)
				}
			}
			o["selected"] = slices.Contains(values, fmt.Sprint(o["value"]))
		}
		var err error
		props[k], err = json.Marshal(opts)
		return err
	}
	return nil
}

// actions converts an "on" handler to the Alpine expression that runs its
// actions through the jsonui.js runtime.
func (b *docBuilder) actions(event string, raw json.RawMessage) (string, error) {
	if !eventName.MatchString(event) {
		return "", fmt.Errorf("invalid event name %q", event)
	}
	list, err := b.resolve(raw, true)
	if err != nil {
		return "", fmt.Errorf("%s: %w", event, err)
	}
	out, err := json.Marshal(list)
	if err != nil {
		return "", err
	}
	return "$actions(" + string(out) + ", $event)", nil
}

// resolve converts a handler (an action, a handler name when refs is true,
// or a list of either) to the list of actions it runs, checking that each
// action is allowed.
func (b *docBuilder) resolve(raw json.RawMessage, refs bool) ([]map[string]any, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	items, ok := v.([]any)
	if !ok {
		items = []any{v}
	}
	var list []map[string]any
	for _, item := range items {
		switch item := item.(type) {
		case string:
			if !refs {
				return nil, fmt.Errorf("handlers can't refer to other handlers (%q)", item)
			}
			h, ok := b.handlers[item]
			if !ok {
				return nil, fmt.Errorf("unknown handler %q", item)
			}
			list = append(list, h...)
		case map[string]any:
			name, _ := item["action"].(string)
			if name == "" {
				return nil, errors.New("every action needs an \"action\" name")
			}
			if !b.opts.allowsAction(name) {
				return nil, fmt.Errorf("action %q is not allowed", name)
			}
			if err := checkActionPaths(name, item); err != nil {
				return nil, err
			}
			list = append(list, item)
		default:
			return nil, errors.New("expected an action, a handler name or a list of them")
		}
	}
	return list, nil
}

// dataPath matches the plain dotted paths that can be resolved in Go, as
// opposed to arbitrary Alpine expressions.
var dataPath = regexp.MustCompile(`^[A-Za-z_$][\w$]*(\.[\w$]+)*$`)

// unsafeKeys are path segments that reach JavaScript object prototypes.
var unsafeKeys = []string{"__proto__", "prototype", "constructor"}

// safePath reports whether path is a plain data path that doesn't reach an
// object prototype, so writing to it can't pollute the host page.
func safePath(path string) bool {
	if !dataPath.MatchString(path) {
		return false
	}
	for _, key := range strings.Split(path, ".") {
		if slices.Contains(unsafeKeys, key) {
			return false
		}
	}
	return true
}

// checkActionPaths checks the action params that jsonui.js writes to as
// data paths: into and busy, and the path of the built-in set action.
func checkActionPaths(name string, action map[string]any) error {
	keys := []string{"into", "busy"}
	if name == "set" {
		keys = append(keys, "path")
	}
	for _, key := range keys {
		v, ok := action[key]
		if !ok {
			continue
		}
		if path, _ := v.(string); !safePath(path) {
			raw, _ := json.Marshal(v)
			return fmt.Errorf("action %q: %s %s is not a data path", name, key, raw)
		}
	}
	return nil
}

// lookup resolves a dotted path such as "invoice.lines.0.total" in data.
func lookup(data map[string]any, path string) (any, error) {
	var cur any = data
	for _, key := range strings.Split(path, ".") {
		switch v := cur.(type) {
		case map[string]any:
			next, ok := v[key]
			if !ok {
				return nil, fmt.Errorf("data has no %q", path)
			}
			cur = next
		case []any:
			i, err := strconv.Atoi(key)
			if err != nil || i < 0 || i >= len(v) {
				return nil, fmt.Errorf("data has no %q", path)
			}
			cur = v[i]
		default:
			return nil, fmt.Errorf("data has no %q", path)
		}
	}
	return cur, nil
}

// display formats a data value as text.
func display(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		out, _ := json.Marshal(v)
		return string(out)
	}
}

func withChildren(c, children templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, children), w)
	})
}

// decodeStrict rejects unknown fields so typos in props surface as errors
// instead of being silently ignored, and anything after the JSON value.
func decodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("unexpected data after the JSON value")
	}
	return nil
}
