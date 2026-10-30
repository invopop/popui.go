package jsonui

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/a-h/templ"
)

// Lists and conditions are rendered as Alpine templates, so they follow the
// data as it changes, e.g. when an action stores a list fetched by the host
// with "into". Nodes repeated with "each" or shown with "show" must render a
// single element, as Alpine templates only keep their first one.

const (
	optionType = "option"
	selectType = "select"
	textType   = "text"
	// xText is the Alpine directive that sets an element's text.
	xText = "x-text"
	// jsonNull is null in JSON, as a prop set to it is unset.
	jsonNull = "null"
	// propDisabled is the prop that disables a control, kept live by its
	// Alpine binding, bindDisabled.
	propDisabled = "disabled"
	bindDisabled = ":disabled"
	// loopItem and loopIndex are the names an "each" node and its children
	// use for the current item and its position.
	loopItem  = "item"
	loopIndex = "index"
)

// Option are the props of the "option" node, a choice in a select. Options
// rendered from a list with "each" keep the select's model selected as the
// list changes:
//
//	{
//	  "is": "select",
//	  "model": "config.connection_id",
//	  "children": [
//	    {"is": "option", "props": {"value": "", "label": "Default"}},
//	    {"is": "option", "each": "connections", "bind": {"value": "item.id", "label": "item.name"}}
//	  ]
//	}
type Option struct {
	Value    string
	Label    string
	Disabled bool
	Selected bool
	// Attributes are added to the option element.
	Attributes templ.Attributes
}

// optionComponent renders an option. Its value is rendered even when empty,
// as an option without one would use its label.
func optionComponent(p ...Option) templ.Component {
	var o Option
	if len(p) > 0 {
		o = p[0]
	}
	return option(o)
}

// liveProps lists the props that "bind" keeps in sync with the data, by
// component, and the Alpine directive that does it. Other props are only
// resolved at render time.
var liveProps = map[string]map[string]string{
	optionType: {"value": ":value", "label": xText, propDisabled: bindDisabled},
	"input":    controlLiveProps,
	"textarea": controlLiveProps,
	selectType: controlLiveProps,
	"checkbox": controlLiveProps,
	"button":   controlLiveProps,
}

// controlLiveProps are the live props of form controls and buttons.
var controlLiveProps = map[string]string{propDisabled: bindDisabled}

// selectState is what a single select's options need to stay in step with
// its model.
type selectState struct {
	// model is the Alpine expression of the select's model.
	model string
	// value is the model's value at render time.
	value any
	// statics are the values of the options not rendered from a list.
	statics []string
	// loops are the lists options are rendered from, and the expression of
	// each option's value in the list.
	loops []selectLoop
}

type selectLoop struct {
	list, value string
}

// livePath checks a data path read by an Alpine binding and returns its
// expression, and whether it's relative to the current "each" item. Other
// paths must be in the data in strict mode, so they can't reach the page's
// globals.
func (b *docBuilder) livePath(path string) (string, bool, error) {
	if !safePath(path) {
		return "", false, fmt.Errorf("%q is not a data path", path)
	}
	if len(b.loops) > 0 {
		root, _, _ := strings.Cut(path, ".")
		if root == loopItem || root == loopIndex {
			return alpinePath(path), true, nil
		}
	}
	if b.data == nil {
		return "", false, errors.New("document data must be set")
	}
	if !b.opts.Trusted {
		if _, err := lookup(b.data, path); err != nil {
			return "", false, err
		}
	}
	return alpinePath(path), false, nil
}

// each renders the node once for every item of the list at its "each" path.
func (b *docBuilder) each(n *Node) (templ.Component, error) {
	list, _, err := b.livePath(n.Each)
	if err != nil {
		return nil, fmt.Errorf("%s each: %w", n.Component, err)
	}
	if v, err := lookup(b.data, n.Each); err == nil && v != nil {
		if _, ok := v.([]any); !ok {
			return nil, fmt.Errorf("%s each: %q is not a list", n.Component, n.Each)
		}
	}
	item := *n
	item.Each = ""
	b.loops = append(b.loops, list)
	c, err := b.element(&item)
	b.loops = b.loops[:len(b.loops)-1]
	if err != nil {
		return nil, err
	}
	return withChildren(loopTemplate(fmt.Sprintf("(%s, %s) in %s", loopItem, loopIndex, list)), c), nil
}

// show renders the node only while its "show" condition holds.
func (b *docBuilder) show(n *Node) (templ.Component, error) {
	cond, err := b.condition(n.Show)
	if err != nil {
		return nil, fmt.Errorf("%s show: %w", n.Component, err)
	}
	inner := *n
	inner.Show = nil
	c, err := b.element(&inner)
	if err != nil {
		return nil, err
	}
	return withChildren(ifTemplate(cond), c), nil
}

// element builds a node for a template, which keeps only its first element,
// so plain text is given one.
func (b *docBuilder) element(n *Node) (templ.Component, error) {
	if n.Component == textType && len(n.Bind) == 0 && n.Each == "" && len(n.Show) == 0 {
		return plainText(n.Text), nil
	}
	return b.build(n)
}

// condition converts a "show" condition to an Alpine expression. It is
// either a data path, true when its value is, optionally negated with "!",
// e.g. "!config.sandbox", or an object of data paths to the list of values
// each may have, all of which must match, e.g.
// {"workflow.schema": ["bill/invoice"], "config.mode": ["a", "b"]}.
// Conditions are data rather than expressions, so documents can't run
// script through them.
func (b *docBuilder) condition(raw json.RawMessage) (string, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	switch v := v.(type) {
	case string:
		path, negate := strings.CutPrefix(v, "!")
		expr, _, err := b.livePath(path)
		if err != nil {
			return "", err
		}
		if negate {
			return "!" + expr, nil
		}
		return expr, nil
	case map[string]any:
		if len(v) == 0 {
			return "", errors.New("expected at least one data path")
		}
		paths := make([]string, 0, len(v))
		for p := range v {
			paths = append(paths, p)
		}
		sort.Strings(paths) // stable output
		parts := make([]string, len(paths))
		for i, p := range paths {
			expr, _, err := b.livePath(p)
			if err != nil {
				return "", err
			}
			values, ok := v[p].([]any)
			if !ok || len(values) == 0 {
				return "", fmt.Errorf("%s: expected a list of values", p)
			}
			for _, val := range values {
				switch val.(type) {
				case string, float64, bool, nil:
				default:
					return "", fmt.Errorf("%s: values must be strings, numbers, booleans or null", p)
				}
			}
			lits, err := json.Marshal(values)
			if err != nil {
				return "", err
			}
			parts[i] = fmt.Sprintf("%s.includes(%s)", lits, expr)
		}
		return strings.Join(parts, " && "), nil
	default:
		return "", errors.New("expected a data path or an object of data paths to values")
	}
}

// selectOption adds what keeps an option of a single select in step with the
// select's model: it's selected while its value is the model's, whether the
// option, the model or the other options change. value is the Alpine
// expression of a live value, or the value itself.
func (b *docBuilder) selectOption(props map[string]json.RawMessage, attrs map[string]any, value string, live bool) error {
	s := b.sel
	if s == nil {
		return errors.New("option: must be a child of a select without options or multiple")
	}
	if s.model == "" {
		return nil
	}
	expr := value
	if !live {
		s.statics = append(s.statics, value)
		lit, err := json.Marshal(value)
		if err != nil {
			return err
		}
		expr = string(lit)
		if props["selected"], err = json.Marshal(value == display(s.value)); err != nil {
			return err
		}
	} else if len(b.loops) > 0 {
		s.loops = append(s.loops, selectLoop{list: b.loops[len(b.loops)-1], value: value})
	}
	attrs[":selected"] = fmt.Sprintf("String(%s) === String(%s)", expr, s.model)
	return nil
}

// unavailableOption is shown when the select's model has a value none of its
// options has, e.g. a saved choice since removed from the list it came from.
// The value is kept rather than cleared, so nothing changes until the user
// chooses another. It's also shown while lists are still empty, such as while
// they load, so the saved choice is still the one shown.
func (s *selectState) unavailableOption() templ.Component {
	if len(s.loops) == 0 {
		return nil
	}
	statics, _ := json.Marshal(append([]string{""}, s.statics...)) //nolint:errcheck // strings always marshal
	conds := []string{
		fmt.Sprintf("%s != null", s.model),
		fmt.Sprintf("!%s.includes(String(%s))", statics, s.model),
	}
	lengths := make([]string, len(s.loops))
	for i, l := range s.loops {
		conds = append(conds, fmt.Sprintf("!(%s || []).some((%s, %s) => String(%s) === String(%s))", l.list, loopItem, loopIndex, l.value, s.model))
		lengths[i] = fmt.Sprintf("(%s || []).length", l.list)
	}
	label := fmt.Sprintf("String(%s) + (%s ? ' (unavailable)' : '')", s.model, strings.Join(lengths, " + "))
	return withChildren(ifTemplate(strings.Join(conds, " && ")), unavailable(s.model, label))
}

// liveBind adds the Alpine bindings for the props a component keeps in sync
// with the data. Props bound to the "each" item can only be live ones, as
// they can't be resolved at render time.
func (b *docBuilder) liveBind(typ, prop, path string, attrs map[string]any) (bool, error) {
	expr, inLoop, err := b.livePath(path)
	if err != nil {
		return false, err
	}
	directive, ok := liveProps[typ][strings.ToLower(prop)]
	if !ok {
		if inLoop {
			return false, fmt.Errorf("can't be bound to the %s; %s", loopItem, livePropNames(typ))
		}
		return false, nil
	}
	attrs[directive] = expr
	return inLoop, nil
}

// livePropNames says which of a component's props can be bound to the item.
func livePropNames(typ string) string {
	names := make([]string, 0, len(liveProps[typ]))
	for k := range liveProps[typ] {
		names = append(names, k)
	}
	if len(names) == 0 {
		return "no props can"
	}
	slices.Sort(names)
	return "only " + strings.Join(names, ", ") + " can"
}
