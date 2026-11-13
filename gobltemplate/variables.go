package gobltemplate

import (
	"strconv"
	"strings"
)

// Variable is a value that a template prints by its key, such as
// {{.full_code}}. Saved templates use the keys, so a key is never renamed or
// removed.
type Variable struct {
	Key         string `json:"key"`
	Description string `json:"description"`

	// schemas carry the variable. Nil means every document, of any schema.
	schemas []Schema
	// always marks a value that every job has.
	always bool
	value  func(d document, job Job) string
}

// variables is the table of every variable, in the order the picker lists
// them.
var variables = []Variable{
	{Key: "type", Description: "Type", schemas: billSchemas, value: text("type")},
	{Key: "series", Description: "Series", schemas: billSchemas, value: text("series")},
	{Key: "code", Description: "Code", schemas: billSchemas, value: text("code")},
	{Key: "full_code", Description: "Full code (series-code)", schemas: billSchemas, value: func(d document, _ Job) string {
		return fullCode(d)
	}},
	{Key: "issue_date", Description: "Issue date", schemas: billSchemas, value: text("issue_date")},
	{Key: "supplier_name", Description: "Supplier name", schemas: billSchemas, value: text("supplier", "name")},
	{Key: "supplier_email", Description: "Supplier email", schemas: billSchemas, value: func(d document, _ Job) string {
		return firstEmail(d.object("supplier"))
	}},
	{Key: "customer_name", Description: "Customer name", schemas: billSchemas, value: text("customer", "name")},
	{Key: "customer_email", Description: "Customer email", schemas: billSchemas, value: func(d document, _ Job) string {
		return firstEmail(d.object("customer"))
	}},

	{Key: "total", Description: "Total", schemas: []Schema{Payment}, value: text("total")},
	{Key: "currency", Description: "Currency", schemas: []Schema{Payment, Item}, value: text("currency")},
	{Key: "preceding_codes", Description: "Preceding documents", schemas: []Schema{Invoice, Payment, Order, Delivery}, value: func(d document, _ Job) string {
		var codes []string
		for _, ref := range d.list("preceding") {
			codes = append(codes, fullCode(ref))
		}
		return joinList(codes)
	}},
	{Key: "due_dates", Description: "Due dates", schemas: []Schema{Invoice}, value: func(d document, _ Job) string {
		var dates []string
		for _, due := range d.list("payment", "terms", "due_dates") {
			dates = append(dates, due.text("date"))
		}
		return joinList(dates)
	}},

	{Key: "line_codes", Description: "Line documents", schemas: []Schema{Status}, value: func(d document, _ Job) string {
		var codes []string
		for _, line := range d.list("lines") {
			codes = append(codes, fullCode(line.object("doc")))
		}
		return joinList(codes)
	}},
	{Key: "line_statuses", Description: "Line statuses", schemas: []Schema{Status}, value: func(d document, _ Job) string {
		var keys []string
		for _, line := range d.list("lines") {
			keys = append(keys, line.text("key"))
		}
		return joinList(keys)
	}},
	{Key: "line_summary", Description: "Line summary", schemas: []Schema{Status}, value: func(d document, _ Job) string {
		var summary []string
		for _, line := range d.list("lines") {
			summary = append(summary, joinNonEmpty(": ", fullCode(line.object("doc")), line.text("key")))
		}
		return joinList(summary)
	}},
	{Key: "line_reasons", Description: "Line reasons", schemas: []Schema{Status}, value: func(d document, _ Job) string {
		var reasons []string
		for _, line := range d.list("lines") {
			for _, reason := range line.list("reasons") {
				if desc := reason.text("description"); desc != "" {
					reasons = append(reasons, desc)
				} else {
					reasons = append(reasons, reason.text("key"))
				}
			}
		}
		return joinList(reasons)
	}},

	{Key: "name", Description: "Name", schemas: []Schema{Party, Item}, value: text("name")},
	{Key: "tax_id", Description: "Tax ID", schemas: []Schema{Party}, value: func(d document, _ Job) string {
		id := d.object("tax_id")
		return id.text("country") + id.text("code")
	}},
	{Key: "email", Description: "Email", schemas: []Schema{Party}, value: func(d document, _ Job) string {
		return firstEmail(d)
	}},
	{Key: "alias", Description: "Alias", schemas: []Schema{Party}, value: text("alias")},
	{Key: "ref", Description: "Reference", schemas: []Schema{Item}, value: text("ref")},
	{Key: "description", Description: "Description", schemas: []Schema{Item}, value: text("description")},
	{Key: "price", Description: "Price", schemas: []Schema{Item}, value: text("price")},

	{Key: "job_id", Description: "Job ID", always: true, value: func(_ document, j Job) string { return j.ID }},
	{Key: "silo_entry_id", Description: "Silo entry ID", always: true, value: func(_ document, j Job) string { return j.SiloEntryID }},
	{Key: "faults", Description: "Faults", value: func(_ document, j Job) string { return faults(j.Faults) }},
	{Key: "meta_link", Description: "Link", value: func(_ document, j Job) string { return j.Link }},
}

// Variables returns the variables that a template about the schema can
// print, given as a short path or a full id. A schema with no variables of
// its own gets the job's variables only.
func Variables(schema string) []Variable {
	var out []Variable
	for _, v := range variables {
		if v.carriedBy(schema) {
			out = append(out, v)
		}
	}
	return out
}

func (v Variable) carriedBy(schema string) bool {
	if v.schemas == nil {
		return true
	}
	for _, s := range v.schemas {
		if s.matches(schema) {
			return true
		}
	}
	return false
}

// document is a GOBL document as its JSON decodes: objects are maps and
// arrays are slices.
type document map[string]any

// at returns the value at a path of object keys, or nil.
func (d document) at(path ...string) any {
	var cur any = map[string]any(d)
	for _, key := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[key]
	}
	return cur
}

// text returns the scalar at a path as a template prints it, or "".
func (d document) text(path ...string) string {
	switch v := d.at(path...).(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	}
	return ""
}

// object returns the object at a path, or nil.
func (d document) object(path ...string) document {
	m, _ := d.at(path...).(map[string]any)
	return m
}

// list returns the objects of the array at a path.
func (d document) list(path ...string) []document {
	items, _ := d.at(path...).([]any)
	out := make([]document, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func text(path ...string) func(document, Job) string {
	return func(d document, _ Job) string { return d.text(path...) }
}

// fullCode joins a document's series and code as gobl's cbc.Code.Join does.
func fullCode(d document) string {
	return joinNonEmpty("-", d.text("series"), d.text("code"))
}

// firstEmail returns a party's first email address.
func firstEmail(party document) string {
	for _, email := range party.list("emails") {
		if addr := email.text("addr"); addr != "" {
			return addr
		}
	}
	return ""
}

func faults(list []Fault) string {
	var b strings.Builder
	for _, f := range list {
		if f.Code != "" {
			b.WriteString(f.Code)
			b.WriteString(": ")
		}
		b.WriteString(f.Message)
		b.WriteString("\n")
	}
	return b.String()
}

// joinList joins the non-empty values of a list, in order.
func joinList(values []string) string {
	return joinNonEmpty(", ", values...)
}

func joinNonEmpty(sep string, values ...string) string {
	kept := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			kept = append(kept, v)
		}
	}
	return strings.Join(kept, sep)
}
