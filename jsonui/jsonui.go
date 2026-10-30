// Package jsonui renders PopUI components from a JSON description of the
// component tree. It is EXPERIMENTAL: the same code is compiled natively for
// server-side rendering and to WebAssembly (see ./wasm) so a browser can turn
// JSON straight into PopUI HTML without a Go server round-trip.
//
// A document is an object with a list of components and optional Alpine.js
// data, which is also the source for bound values:
//
//	{
//	  "$schema": "https://…/jsonui/schema/document.json",
//	  "data": {"customer": {"name": "Acme S.L."}},
//	  "components": [
//	    {
//	      "is": "card",
//	      "children": [
//	        {"is": "card_header", "bind": {"title": "customer.name"}},
//	        {"is": "text", "bind": "customer.name"}
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
	"maps"
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
	// Components is the list of nodes to render.
	Components []Node `json:"components"`
	// Handlers are named actions, or lists of actions, that "on" handlers
	// can refer to by name so a handler used in several places is defined
	// once, e.g. {"update": {"action": "step.update", "summary": "…"}} and
	// then {"on": {"change": "update"}}. Handlers can't refer to each other.
	Handlers map[string]json.RawMessage `json:"handlers,omitempty"`
	// Init is run once when the document starts, like an "on" handler: an
	// action, the name of one of the document's handlers, or a list of
	// either, e.g. "update" so the host receives the initial data.
	Init json.RawMessage `json:"init,omitempty"`
	// Schemas map data paths, such as "config", to the JSON Schema their
	// value must match. The data is validated as it changes, and errors are
	// shown with the form control whose "model" is their path, or after the
	// document for paths without one. Schemas also give controls their
	// label, from "title", and selects without options their options, from
	// "enum" or the "const" and "title" of each "oneOf". Schemas can only
	// refer to themselves with "$ref". A schema's "errorMessage" replaces the
	// messages of its errors: one message, or one for each keyword, such as
	// {"format": "Enter an email address"}.
	Schemas map[string]json.RawMessage `json:"schemas,omitempty"`
}

// Node is a single element in the JSON component tree.
type Node struct {
	// Component is the registered component name, e.g. "button" or
	// "card_header", or "text" for plain text. In JSON it is "is", after
	// HTML's global "is" attribute, which names the custom element to use;
	// "type" is left to the props, where inputs and buttons use it.
	Component string `json:"is"`
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
	// to data paths, e.g. {"value": "customer.name"}, resolved at render
	// time. Some props, such as "disabled", also follow the data as it
	// changes, and only those can be bound to an "each" item.
	Bind json.RawMessage `json:"bind,omitempty"`
	// Each repeats the node for every item of the list at a data path, e.g.
	// "connections". The node and its children refer to the current item as
	// "item" and its position as "index", e.g. {"value": "item.id"} in
	// "bind". The node must render a single element.
	Each string `json:"each,omitempty"`
	// Show renders the node only while a condition on the data holds: a data
	// path that is true when its value is, optionally negated with "!", e.g.
	// "config.sandbox", or an object of data paths to the values each may
	// have, all of which must match, e.g. {"workflow.schema":
	// ["bill/invoice"]}. The node must render a single element.
	Show json.RawMessage `json:"show,omitempty"`
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
	selectType: "options",
	"checkbox": "checked",
}

// Component describes a registered component.
type Component struct {
	// Name is the component name used in JSON.
	Name string
	// Props is the props struct the node's "props" object decodes into.
	Props reflect.Type
	// Model is true for form controls that accept a "model" binding.
	Model bool
	// Children is true for nodes that render their "children".
	Children bool
	// Bind is true for nodes whose props can be set with "bind".
	Bind bool
	// Events lists the events the node accepts in "on", or nil for any DOM
	// event.
	Events []string
}

// builder converts raw props into a component, children are attached later.
type builder func(raw json.RawMessage) (templ.Component, error)

type entry struct {
	props reflect.Type
	build builder
	// children is true when the component renders its children.
	children bool
}

var registry = map[string]entry{}

// register adds a PopUI component whose signature follows the standard
// `templ Name(p ...props.X)` convention.
func register[P any](name string, fn func(...P) templ.Component) {
	var zero P
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
		children: rendersChildren(fn(zero)),
	}
}

// childMarker is rendered as the children of a component to find out
// whether it renders them.
const childMarker = "\x00jsonui-children\x00"

// rendersChildren reports whether c renders its children, by rendering it
// with its default props and a marker child, so that documents can't supply
// children a component would silently drop.
func rendersChildren(c templ.Component) bool {
	marker := templ.Raw(childMarker)
	var buf bytes.Buffer
	if err := c.Render(templ.WithChildren(context.Background(), marker), &buf); err != nil {
		return false
	}
	return bytes.Contains(buf.Bytes(), []byte(childMarker))
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
	register(selectType, popui.Select)
	register("checkbox", popui.Checkbox)
	register("option_group", popui.OptionGroup)
	register(optionType, optionComponent)

	// Navigation
	register("tabs", popui.Tabs)
	register("tab", popui.Tab)
}

// Components returns the registered components sorted by name. The
// built-in "text" node is not included as it has no props.
func Components() []Component {
	out := make([]Component, 0, len(registry))
	for k, e := range registry {
		_, model := modelProps[k]
		// Entries without a props builder are custom nodes (upload_file),
		// which build their props themselves.
		out = append(out, Component{Name: k, Props: e.props, Model: model, Children: e.children, Bind: e.build != nil, Events: nodeEvents[k]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Names returns the sorted names of the supported components, including
// "text".
func Names() []string {
	out := []string{textType}
	for _, c := range Components() {
		out = append(out, c.Name)
	}
	sort.Strings(out)
	return out
}

// Component builds the document's components. When the document has data,
// they are wrapped in an Alpine.js component root.
func (d *Document) Component(opts Options) (templ.Component, error) {
	if d.Components == nil {
		return nil, errors.New("document: missing components")
	}
	b := &docBuilder{opts: opts}
	var schemas string
	if len(d.Schemas) > 0 {
		var err error
		if b.validator, err = NewValidator(d.Schemas); err != nil {
			return nil, err
		}
		if schemas, err = compactSchemas(d.Schemas); err != nil {
			return nil, err
		}
	}
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
	body, err := b.buildAll(d.Components)
	if err != nil {
		return nil, err
	}
	xinit := ""
	if len(d.Init) > 0 {
		list, err := b.actionList(d.Init)
		if err != nil {
			return nil, fmt.Errorf("init: %w", err)
		}
		// There's no event when the document starts.
		xinit = "$actions(" + list + ")"
		if xdata == "" {
			xdata = "{}" // init runs on the document root
		}
	}
	errors := b.fields > 0 || b.validator != nil
	if errors {
		body = templ.Join(body, unplacedErrors())
		if xdata == "" {
			xdata = "{}" // errors are tracked on the document root
		}
	}
	if xdata == "" {
		return body, nil
	}
	return withChildren(documentRoot(xdata, xinit, errors, schemas), body), nil
}

// compactSchemas returns the schemas as the compact JSON jsonui.js
// validates the data with.
func compactSchemas(schemas map[string]json.RawMessage) (string, error) {
	out := make(map[string]json.RawMessage, len(schemas))
	for k, raw := range schemas {
		var buf bytes.Buffer
		if err := json.Compact(&buf, raw); err != nil {
			return "", fmt.Errorf("schemas %s: %w", k, err)
		}
		out[k] = buf.Bytes()
	}
	data, err := json.Marshal(out)
	return string(data), err
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
	// loops are the list expressions of the "each" nodes being built.
	loops []string
	// sel is the state of the single select whose options are being built.
	sel *selectState
	// validator checks the data against the document's schemas, if any.
	validator *Validator
	// fields counts the form controls that show errors for their model.
	fields int
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
	if n.Each != "" {
		return b.each(n)
	}
	if len(n.Show) > 0 {
		return b.show(n)
	}
	if n.Component == textType {
		return b.text(n)
	}
	if !b.opts.Trusted {
		if err := checkProps(n.Component, n.Props); err != nil {
			return nil, err
		}
	}
	if n.Component == uploadFileType {
		return b.uploadFile(n)
	}
	e, ok := registry[n.Component]
	if !ok {
		return nil, fmt.Errorf("unknown component %q", n.Component)
	}
	raw, err := b.props(n)
	if err != nil {
		return nil, err
	}
	c, err := e.build(raw)
	if err != nil {
		return nil, err
	}
	if len(n.Children) > 0 {
		if !e.children {
			return nil, fmt.Errorf("%s: children are not supported", n.Component)
		}
		if c, err = b.withChildNodes(n, c); err != nil {
			return nil, err
		}
	}
	if path := b.fieldPath(n); path != "" {
		c = withChildren(field(path, shownError(path)), c)
	}
	return c, nil
}

// fieldPath is the data path a form control shows errors for: its model,
// unless that's in an "each" item, which has no path of its own.
func (b *docBuilder) fieldPath(n *Node) string {
	if _, ok := modelProps[n.Component]; !ok || n.Model == "" || !safePath(n.Model) {
		return ""
	}
	if len(b.loops) > 0 {
		if root, _, _ := strings.Cut(n.Model, "."); root == loopItem || root == loopIndex {
			return ""
		}
	}
	return n.Model
}

// shownError is the Alpine expression of the error shown for a path, or
// "" when there's none. Paths are data paths, so can be quoted as is.
func shownError(path string) string {
	return "$errors.shown('" + path + "')"
}

// withChildNodes builds the node's children into c. A select's children are
// its options, which follow its model.
func (b *docBuilder) withChildNodes(n *Node, c templ.Component) (templ.Component, error) {
	outer := b.sel
	defer func() { b.sel = outer }()
	b.sel = nil
	if n.Component == selectType {
		var err error
		if b.sel, err = b.selectState(n); err != nil {
			return nil, err
		}
	}
	children, err := b.buildAll(n.Children)
	if err != nil {
		return nil, err
	}
	if b.sel != nil {
		if u := b.sel.unavailableOption(); u != nil {
			children = templ.Join(u, children)
		}
	}
	return withChildren(c, children), nil
}

// selectState prepares the state of a select with options as children,
// which only the single select renders.
func (b *docBuilder) selectState(n *Node) (*selectState, error) {
	var p struct {
		Multiple bool
		Options  json.RawMessage
	}
	if len(n.Props) > 0 {
		if err := json.Unmarshal(n.Props, &p); err != nil {
			return nil, fmt.Errorf("%s props: %w", n.Component, err)
		}
	}
	if p.Multiple || (len(p.Options) > 0 && string(p.Options) != jsonNull) {
		return nil, fmt.Errorf("%s: children are only supported by single selects without options", n.Component)
	}
	s := new(selectState)
	if n.Model != "" {
		expr, inLoop, err := b.livePath(n.Model)
		if err != nil {
			return nil, fmt.Errorf("%s model: %w", n.Component, err)
		}
		s.model = expr
		if !inLoop {
			s.value, _ = lookup(b.data, n.Model) //nolint:errcheck // checked by livePath in strict mode
		}
	}
	return s, nil
}

func (b *docBuilder) text(n *Node) (templ.Component, error) {
	if len(n.Bind) == 0 {
		return templ.Raw(templ.EscapeString(n.Text)), nil
	}
	var expr string
	if err := json.Unmarshal(n.Bind, &expr); err != nil {
		return nil, errors.New("text bind: must be a string expression")
	}
	if len(b.loops) > 0 {
		if root, _, _ := strings.Cut(expr, "."); root == loopItem || root == loopIndex {
			path, _, err := b.livePath(expr)
			if err != nil {
				return nil, fmt.Errorf("text bind: %w", err)
			}
			return boundText(n.Text, path), nil
		}
	}
	text := n.Text
	if !b.opts.Trusted {
		// Alpine resolves names it can't find in the data from the page's
		// globals, so x-text="document.cookie" would read the host page. In
		// strict mode the path must therefore be in the document's data.
		if !safePath(expr) {
			return nil, fmt.Errorf("text bind: %q must be a data path in strict mode", expr)
		}
		if b.data == nil {
			return nil, errors.New("text bind: document data must be set")
		}
		v, err := lookup(b.data, expr)
		if err != nil {
			return nil, fmt.Errorf("text bind: %w", err)
		}
		text = display(v)
	} else if b.data != nil && dataPath.MatchString(expr) {
		if v, err := lookup(b.data, expr); err == nil {
			text = display(v)
		}
	}
	return boundText(text, alpinePath(expr)), nil
}

// props returns the node's props with any bind, model and on values from
// the node applied.
func (b *docBuilder) props(n *Node) (json.RawMessage, error) {
	if len(n.Bind) == 0 && n.Model == "" && len(n.On) == 0 && n.Component != optionType {
		return n.Props, nil
	}
	var props map[string]json.RawMessage
	if len(n.Props) > 0 {
		if err := json.Unmarshal(n.Props, &props); err != nil {
			return nil, fmt.Errorf("%s props: %w", n.Component, err)
		}
	}
	if props == nil { // missing or null
		props = map[string]json.RawMessage{}
	}
	live := map[string]any{}
	if len(n.Bind) > 0 {
		if err := b.bindProps(n.Component, props, live, n.Bind); err != nil {
			return nil, err
		}
		if !b.opts.Trusted {
			// Bound values come from the data, so check them like props.
			raw, err := json.Marshal(props)
			if err != nil {
				return nil, err
			}
			if err := checkProps(n.Component, raw); err != nil {
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
				return nil, fmt.Errorf("%s attributes: %w", n.Component, err)
			}
		}
	}
	if attrs == nil { // missing or null
		attrs = map[string]any{}
	}
	maps.Copy(attrs, live)
	if path := b.fieldPath(n); path != "" {
		if err := b.schemaProps(n, props); err != nil {
			return nil, err
		}
		b.fieldAttrs(path, attrs)
	}
	if n.Component == optionType {
		if err := b.optionSelection(props, attrs); err != nil {
			return nil, err
		}
	}
	if n.Model != "" {
		if err := b.model(n.Component, props, attrs, n.Model); err != nil {
			return nil, err
		}
	}
	for event, raw := range n.On {
		handler, err := b.actions(event, raw)
		if err != nil {
			return nil, fmt.Errorf("%s on: %w", n.Component, err)
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

// optionSelection keeps an option selected while its value is its select's
// model, whether the value is bound or set in its props.
func (b *docBuilder) optionSelection(props map[string]json.RawMessage, attrs map[string]any) error {
	if expr, ok := attrs[":value"].(string); ok {
		return b.selectOption(props, attrs, expr, true)
	}
	var value string
	for k, raw := range props {
		if strings.EqualFold(k, "value") {
			if err := json.Unmarshal(raw, &value); err != nil {
				return fmt.Errorf("option value: %w", err)
			}
		}
	}
	return b.selectOption(props, attrs, value, false)
}

// bindProps overlays the props with values taken from the document data,
// and adds the bindings of the props that follow the data to live.
func (b *docBuilder) bindProps(typ string, props map[string]json.RawMessage, live map[string]any, bind json.RawMessage) error {
	var paths map[string]string
	if err := decodeStrict(bind, &paths); err != nil {
		return fmt.Errorf("%s bind: must map prop names to data paths: %w", typ, err)
	}
	if b.data == nil && len(b.loops) == 0 {
		return fmt.Errorf("%s bind: document data must be set", typ)
	}
	for prop, path := range paths {
		inLoop, err := b.liveBind(typ, prop, path, live)
		if err != nil {
			return fmt.Errorf("%s bind %s: %w", typ, prop, err)
		}
		if inLoop {
			continue // only known once Alpine renders the item
		}
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
	if expr, inLoop, err := b.livePath(path); err == nil && inLoop {
		// The item's value is only known once Alpine renders it.
		attrs["x-model"] = expr
		return nil
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
	attrs[directive] = alpinePath(path)
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
	list, err := b.actionList(raw)
	if err != nil {
		return "", fmt.Errorf("%s: %w", event, err)
	}
	return "$actions(" + list + ", $event)", nil
}

// actionList resolves a handler to the JSON list of actions it runs.
func (b *docBuilder) actionList(raw json.RawMessage) (string, error) {
	list, err := b.resolve(raw, true)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(list)
	if err != nil {
		return "", err
	}
	return string(out), nil
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
	list := []map[string]any{} // an empty list renders as [], not null
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
			if slices.Contains(BuiltinActions, name) {
				for _, param := range RequiredParams[name] {
					if v, _ := item[param].(string); v == "" {
						return nil, fmt.Errorf("action %q: %s is required", name, param)
					}
				}
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

// alpinePath writes a data path as an Alpine expression. Segments that
// aren't identifiers, such as array indexes, use brackets, so
// "invoice.lines.0.total" becomes "invoice.lines[0].total". Anything other
// than a data path, such as an expression in a trusted document, is
// returned unchanged.
func alpinePath(path string) string {
	if !dataPath.MatchString(path) {
		return path
	}
	var out strings.Builder
	for i, key := range strings.Split(path, ".") {
		switch {
		case identifier.MatchString(key):
			if i > 0 {
				out.WriteByte('.')
			}
			out.WriteString(key)
		case strings.Trim(key, "0123456789") == "":
			out.WriteString("[" + key + "]")
		default:
			out.WriteString("['" + key + "']")
		}
	}
	return out.String()
}

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
	if name == actionSet {
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

// lookup resolves a dotted path such as "invoice.lines.0.total" in data. A
// list's "length" is its number of items.
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
			if key == "length" {
				// A list's length, as JavaScript reads it, e.g. for a
				// condition that the list is empty.
				cur = json.Number(strconv.Itoa(len(v)))
				continue
			}
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
