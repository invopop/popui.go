package gobltemplate

import (
	"encoding/json"
	"strings"
	"text/template"
	"text/template/parse"
)

// docKey is the key of the document in the data.
const docKey = "doc"

// Job holds the values that come from the job running the template, not
// from the document.
type Job struct {
	// ID is printed by job_id.
	ID string
	// SiloEntryID is printed by silo_entry_id.
	SiloEntryID string
	// Link is printed by meta_link.
	Link string
	// Faults are printed by faults, one per line.
	Faults []Fault
}

// Fault is a problem that a step of the job reported.
type Fault struct {
	Code    string
	Message string
}

// Data returns what a template about a document executes against:
//
//   - Every variable of the schema, as a string. A value the document does
//     not have is "".
//   - The document under doc, as its JSON decodes.
//
// The schema is a short path or a full id. The document is anything that
// encodes to the document's JSON: a gobl type, a map or a json.RawMessage.
// A nil document leaves doc out.
func Data(schema string, doc any, job Job) (map[string]any, error) {
	d, err := decode(doc)
	if err != nil {
		return nil, err
	}
	vars := Variables(schema)
	data := make(map[string]any, len(vars)+1)
	for _, v := range vars {
		data[v.Key] = v.value(d, job)
	}
	if d != nil {
		data[docKey] = map[string]any(d)
	}
	return data, nil
}

func decode(doc any) (document, error) {
	if doc == nil {
		return nil, nil
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	var d document
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// Render executes a template against data, as Data returns it. A value that
// the template reads and data does not have prints as nothing, not as
// "<no value>". Render adds those values to data. The name shows in errors.
func Render(name, text string, data map[string]any) (string, error) {
	t, err := template.New(name).Parse(text)
	if err != nil {
		return "", err
	}
	if data == nil {
		data = map[string]any{}
	}
	ensurePaths(t, data)
	var b strings.Builder
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// ensurePaths gives every field chain that the template reads, such as
// .doc.customer.name, a value in data. It sets a missing value to "", and
// adds the missing objects on the way. It does not change a value that is
// there, and it stops at a value that is not an object.
func ensurePaths(t *template.Template, data map[string]any) {
	if t.Tree == nil || t.Root == nil {
		return
	}
	walkFields(t.Root, func(idents []string) { ensurePath(data, idents) })
}

func ensurePath(data map[string]any, idents []string) {
	cur := data
	for i, name := range idents {
		next, ok := cur[name]
		if i == len(idents)-1 {
			if !ok {
				cur[name] = ""
			}
			return
		}
		if !ok {
			m := map[string]any{}
			cur[name] = m
			cur = m
			continue
		}
		m, ok := next.(map[string]any)
		if !ok {
			return
		}
		cur = m
	}
}

// walkFields calls fn with the identifiers of every field chain in the tree,
// also inside if, with and range blocks.
func walkFields(node parse.Node, fn func([]string)) {
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, child := range n.Nodes {
			walkFields(child, fn)
		}
	case *parse.ActionNode:
		walkPipe(n.Pipe, fn)
	case *parse.IfNode:
		walkBranch(&n.BranchNode, fn)
	case *parse.WithNode:
		walkBranch(&n.BranchNode, fn)
	case *parse.RangeNode:
		walkBranch(&n.BranchNode, fn)
	}
}

func walkBranch(b *parse.BranchNode, fn func([]string)) {
	walkPipe(b.Pipe, fn)
	walkFields(b.List, fn)
	walkFields(b.ElseList, fn)
}

func walkPipe(pipe *parse.PipeNode, fn func([]string)) {
	if pipe == nil {
		return
	}
	for _, cmd := range pipe.Cmds {
		for _, arg := range cmd.Args {
			if field, ok := arg.(*parse.FieldNode); ok {
				fn(field.Ident)
			}
		}
	}
}
