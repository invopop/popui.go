package jsonui

import (
	"fmt"
	"regexp"
	"slices"
)

// KeyPattern is what a node's "key" looks like: letters, digits, dashes and
// underscores, starting with a letter.
const KeyPattern = `^[A-Za-z][A-Za-z0-9_-]*$`

var keyFormat = regexp.MustCompile(KeyPattern)

// Patch describes changes to a panel, kept apart from it, so that a host
// can let its users adapt panels they didn't write: an app's panel in a
// workspace, say. It refers to nodes by their "key" rather than their
// place, so it applies to later versions of the panel as long as the keys
// are still there, and it can only take away or rearrange what the panel
// has, so a patched panel is as safe as the panel itself.
//
//	{"hide": ["payment-terms"], "order": ["totals", "customer", "supplier"]}
//
// EXPERIMENTAL: this is a first, minimal patch. Adding nodes isn't
// supported yet.
type Patch struct {
	// Hide lists the keys of the nodes to leave out, with their children.
	Hide []string `json:"hide,omitempty"`
	// Order lists keys in the order their nodes should appear. Nodes are
	// only moved among those listed with the same parent, which trade
	// places: the first listed takes the place of whichever came first,
	// and so on, so nodes that aren't listed stay where they are. Keys
	// with different parents can share the list.
	Order []string `json:"order,omitempty"`
}

// Empty reports whether the patch changes nothing.
func (p *Patch) Empty() bool {
	return p == nil || (len(p.Hide) == 0 && len(p.Order) == 0)
}

// Apply returns a copy of the panel with the patch applied, leaving the
// panel as it was. Keys the panel doesn't have are skipped and returned, in
// the order the patch names them, so the host can tell that part of a
// patch no longer applies.
func (d *Panel) Apply(patch *Patch) (*Panel, []string, error) {
	keys, err := d.keys()
	if err != nil {
		return nil, nil, err
	}
	out := *d
	if patch.Empty() {
		out.Components = copyNodes(d.Components, nil)
		return &out, nil, nil
	}
	var missing []string
	for _, k := range slices.Concat(patch.Hide, patch.Order) {
		if !keys[k] && !slices.Contains(missing, k) {
			missing = append(missing, k)
		}
	}
	out.Components = copyNodes(d.Components, patch.Hide)
	reorder(out.Components, patch.Order)
	return &out, missing, nil
}

// Keys returns the keys of the panel's nodes, in the order they appear,
// which are what a Patch can refer to.
func (d *Panel) Keys() ([]string, error) {
	if _, err := d.keys(); err != nil {
		return nil, err
	}
	var list []string
	walk(d.Components, func(n *Node) {
		if n.Key != "" {
			list = append(list, n.Key)
		}
	})
	return list, nil
}

// keys returns the set of keys in the panel, or an error if one is
// malformed or used twice.
func (d *Panel) keys() (map[string]bool, error) {
	keys := make(map[string]bool)
	var err error
	walk(d.Components, func(n *Node) {
		switch {
		case n.Key == "" || err != nil:
		case !keyFormat.MatchString(n.Key):
			err = fmt.Errorf("%s: key %q must be letters, digits, dashes or underscores, starting with a letter", n.Component, n.Key)
		case keys[n.Key]:
			err = fmt.Errorf("%s: key %q is already used in the panel", n.Component, n.Key)
		default:
			keys[n.Key] = true
		}
	})
	return keys, err
}

// walk calls fn for every node, parents before their children.
func walk(nodes []Node, fn func(*Node)) {
	for i := range nodes {
		fn(&nodes[i])
		walk(nodes[i].Children, fn)
	}
}

// copyNodes copies the nodes and their children, without those whose key
// is hidden. What nodes hold besides their children is shared, as nothing
// changes it.
func copyNodes(nodes []Node, hide []string) []Node {
	if nodes == nil {
		return nil
	}
	out := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		if n.Key != "" && slices.Contains(hide, n.Key) {
			continue
		}
		n.Children = copyNodes(n.Children, hide)
		out = append(out, n)
	}
	return out
}

// reorder arranges, among each list of siblings, the nodes whose keys are
// in the order, in the places they hold between them.
func reorder(nodes []Node, order []string) {
	var places []int
	for i := range nodes {
		reorder(nodes[i].Children, order)
		if nodes[i].Key != "" && slices.Contains(order, nodes[i].Key) {
			places = append(places, i)
		}
	}
	if len(places) < 2 {
		return
	}
	moved := make([]Node, len(places))
	for i, at := range places {
		moved[i] = nodes[at]
	}
	slices.SortStableFunc(moved, func(a, b Node) int {
		return slices.Index(order, a.Key) - slices.Index(order, b.Key)
	})
	for i, at := range places {
		nodes[at] = moved[i]
	}
}
