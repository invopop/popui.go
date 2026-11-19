package apps

import (
	"encoding/json"
	"strings"
)

// WorkflowStep is one step of an example workflow: the action (provider) it
// runs and the label shown in the Console.
type WorkflowStep struct {
	Name     string         `json:"name"`
	Provider string         `json:"provider"`
	Config   map[string]any `json:"config,omitempty"`
}

// Workflow is an example workflow published in the app's documentation (the
// "Workflows" tab of the docs page). Like Document, there is no access model
// for it yet; see HANDOFF.md.
type Workflow struct {
	ID          string
	Name        string
	Description string
	// Channel is the app channel the workflow belongs to, empty for default.
	Channel string
	// Schema is the GOBL schema the workflow runs on.
	Schema string
	Steps  []WorkflowStep
	Rescue []WorkflowStep
}

// Content renders the workflow as the JSON readers copy into the Console.
func (w *Workflow) Content() string {
	out, _ := json.MarshalIndent(struct {
		Name        string         `json:"name"`
		Description string         `json:"description,omitempty"`
		Schema      string         `json:"schema"`
		Steps       []WorkflowStep `json:"steps"`
		Rescue      []WorkflowStep `json:"rescue,omitempty"`
	}{w.Name, w.Description, w.Schema, w.Steps, w.Rescue}, "", "  ")
	return string(out)
}

// WorkflowGroup is one deck on the Workflows tab: every workflow of one
// channel that runs on one schema, in display order.
type WorkflowGroup struct {
	Channel   string
	Schema    string
	Workflows []*Workflow
}

// ID identifies the group in routes ("<channel>/<schema>" with the default
// channel spelled "default").
func (g *WorkflowGroup) ID() string {
	ch := g.Channel
	if ch == "" {
		ch = "default"
	}
	return ch + "/" + strings.ReplaceAll(g.Schema, "/", "-")
}

// Reorder applies the order of workflow IDs the CardDeck submitted; IDs it
// does not mention keep their relative order at the end.
func (g *WorkflowGroup) Reorder(ids []string) {
	seen := map[string]bool{}
	var out []*Workflow
	for _, id := range ids {
		for _, w := range g.Workflows {
			if w.ID == id && !seen[id] {
				out = append(out, w)
				seen[id] = true
			}
		}
	}
	for _, w := range g.Workflows {
		if !seen[w.ID] {
			out = append(out, w)
		}
	}
	g.Workflows = out
}

// WorkflowCollection is the ordered set of groups for one app.
type WorkflowCollection struct {
	Groups []*WorkflowGroup
}

// Count returns the number of workflows across all groups.
func (wc *WorkflowCollection) Count() int {
	n := 0
	for _, g := range wc.Groups {
		n += len(g.Workflows)
	}
	return n
}

// Group finds a group by its route ID.
func (wc *WorkflowCollection) Group(id string) *WorkflowGroup {
	for _, g := range wc.Groups {
		if g.ID() == id {
			return g
		}
	}
	return nil
}

// ForChannel lists the groups of one channel in display order.
func (wc *WorkflowCollection) ForChannel(channel string) []*WorkflowGroup {
	var out []*WorkflowGroup
	for _, g := range wc.Groups {
		if g.Channel == channel {
			out = append(out, g)
		}
	}
	return out
}

// Find locates a workflow and its group.
func (wc *WorkflowCollection) Find(id string) (*WorkflowGroup, *Workflow) {
	for _, g := range wc.Groups {
		for _, w := range g.Workflows {
			if w.ID == id {
				return g, w
			}
		}
	}
	return nil, nil
}

// Add places a workflow in the group for its channel and schema, creating
// the group when it is the first of its kind.
func (wc *WorkflowCollection) Add(w *Workflow) {
	for _, g := range wc.Groups {
		if g.Channel == w.Channel && g.Schema == w.Schema {
			g.Workflows = append(g.Workflows, w)
			return
		}
	}
	wc.Groups = append(wc.Groups, &WorkflowGroup{Channel: w.Channel, Schema: w.Schema, Workflows: []*Workflow{w}})
}

// Remove drops a workflow, and its group once empty.
func (wc *WorkflowCollection) Remove(id string) {
	for gi, g := range wc.Groups {
		for i, w := range g.Workflows {
			if w.ID != id {
				continue
			}
			g.Workflows = append(g.Workflows[:i], g.Workflows[i+1:]...)
			if len(g.Workflows) == 0 {
				wc.Groups = append(wc.Groups[:gi], wc.Groups[gi+1:]...)
			}
			return
		}
	}
}

// schemaLabel names a deck after the schema it holds and what kind of thing
// is in it: schemaLabel("bill/invoice", "workflows") is "Invoice workflows".
func schemaLabel(schema, kind string) string {
	nouns := map[string]string{
		"bill/invoice":  "Invoice",
		"org/party":     "Party",
		"bill/status":   "Event",
		"bill/payment":  "Payment",
		"bill/order":    "Order",
		"bill/delivery": "Delivery",
	}
	if noun, ok := nouns[schema]; ok {
		return noun + " " + kind
	}
	return schema
}

func state(s string) WorkflowStep {
	return WorkflowStep{Name: "Set state", Provider: "silo.state", Config: map[string]any{"state": s}}
}

func step(name, provider string) WorkflowStep {
	return WorkflowStep{Name: name, Provider: provider}
}

// sampleWorkflows mirrors the "Workflows" tab of the Spain app page on
// docs.invopop.com (apps/spain.mdx): per channel, invoice workflows then
// party workflows, with NO VERI*FACTU's event workflows last.
var sampleWorkflows = map[string]*WorkflowCollection{
	"gov-es": {Groups: []*WorkflowGroup{
		{Channel: "ticketbai", Schema: "bill/invoice", Workflows: []*Workflow{
			{ID: "tbai-issue", Name: "TicketBAI issue invoice", Description: "Sign the invoice and register it with TicketBAI", Channel: "ticketbai", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Add sequential code", "sequence.enumerate"), step("Sign envelope", "silo.close"), step("Send invoice to TicketBAI", "gov-es.ticketbai.register"), state("registered"), step("Generate PDF", "pdf.generate")},
				Rescue: []WorkflowStep{state("error")}},
			{ID: "tbai-cancel", Name: "TicketBAI cancel invoice", Description: "Cancel a registered invoice in TicketBAI", Channel: "ticketbai", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Cancel invoice in TicketBAI", "gov-es.ticketbai.cancel"), state("void")},
				Rescue: []WorkflowStep{state("error")}},
		}},
		{Channel: "ticketbai", Schema: "org/party", Workflows: []*Workflow{
			{ID: "tbai-register", Name: "TicketBAI register party", Description: "Register a supplier with TicketBAI", Channel: "ticketbai", Schema: "org/party",
				Steps:  []WorkflowStep{step("Register supplier with TicketBAI", "gov-es.ticketbai.register"), state("processing"), step("Wait for TicketBAI supplier agreement upload", "gov-es.ticketbai.wait.upload"), state("registered"), step("Wait for TicketBAI supplier approval", "gov-es.ticketbai.wait.approval")},
				Rescue: []WorkflowStep{state("rejected"), step("Unregister supplier with TicketBAI", "gov-es.ticketbai.unregister")}},
			{ID: "tbai-unregister", Name: "TicketBAI unregister party", Description: "Release the supplier's TicketBAI registration", Channel: "ticketbai", Schema: "org/party",
				Steps: []WorkflowStep{step("Unregister supplier with TicketBAI", "gov-es.ticketbai.unregister"), state("empty")}},
		}},
		{Channel: "", Schema: "bill/invoice", Workflows: []*Workflow{
			{ID: "facturae-generate", Name: "Facturae generate invoice", Description: "Generate and sign a Facturae XML for a public administration", Channel: "", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Add sequential code", "sequence.enumerate"), step("Sign envelope", "silo.close"), step("Generate Facturae XML", "gov-es.facturae.generate"), state("sent")},
				Rescue: []WorkflowStep{state("error")}},
		}},
		{Channel: "sii", Schema: "bill/invoice", Workflows: []*Workflow{
			{ID: "sii-issued", Name: "SII issue invoice", Description: "Record an issued invoice with SII", Channel: "sii", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Add sequential code", "sequence.enumerate"), step("Sign envelope", "silo.close"), step("Record issued invoice with SII", "gov-es.sii.register"), state("registered"), step("Generate PDF", "pdf.generate")},
				Rescue: []WorkflowStep{state("error")}},
			{ID: "sii-received", Name: "SII record received invoice", Description: "Record a received invoice with SII", Channel: "sii", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Record received invoice with SII", "gov-es.sii.register.received"), state("registered")},
				Rescue: []WorkflowStep{state("error")}},
		}},
		{Channel: "sii", Schema: "org/party", Workflows: []*Workflow{
			{ID: "sii-register", Name: "SII register party", Description: "Register a supplier with SII", Channel: "sii", Schema: "org/party",
				Steps:  []WorkflowStep{step("Register supplier with SII", "gov-es.sii.register"), state("processing"), step("Wait for SII supplier agreement upload", "gov-es.sii.wait.upload"), state("registered"), step("Wait for SII supplier approval", "gov-es.sii.wait.approval")},
				Rescue: []WorkflowStep{state("rejected"), step("Unregister supplier from SII", "gov-es.sii.unregister")}},
			{ID: "sii-unregister", Name: "SII unregister party", Description: "Release the supplier's SII registration", Channel: "sii", Schema: "org/party",
				Steps: []WorkflowStep{step("Unregister supplier from SII", "gov-es.sii.unregister"), state("empty")}},
		}},
		{Channel: "verifactu", Schema: "bill/invoice", Workflows: []*Workflow{
			{ID: "nvf-issue", Name: "NO VERI*FACTU issue invoice", Description: "Generate and store a NO VERI*FACTU invoice record", Channel: "verifactu", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Add sequential code", "sequence.enumerate"), step("Sign envelope", "silo.close"), step("Generate invoice for NO VERI*FACTU", "gov-es.verifactu.generate"), step("Record for NO VERI*FACTU", "gov-es.verifactu.record"), state("registered"), step("Generate PDF", "pdf.generate")},
				Rescue: []WorkflowStep{state("error")}},
			{ID: "nvf-cancel", Name: "NO VERI*FACTU cancel invoice", Description: "Generate a cancellation record for NO VERI*FACTU", Channel: "verifactu", Schema: "bill/invoice",
				Steps:  []WorkflowStep{state("processing"), step("Generate cancellation for NO VERI*FACTU", "gov-es.verifactu.cancel"), step("Record for NO VERI*FACTU", "gov-es.verifactu.record"), state("void")},
				Rescue: []WorkflowStep{state("error")}},
		}},
		{Channel: "verifactu", Schema: "org/party", Workflows: []*Workflow{
			{ID: "nvf-register", Name: "NO VERI*FACTU register supplier", Description: "Register a supplier with NO VERI*FACTU", Channel: "verifactu", Schema: "org/party",
				Steps:  []WorkflowStep{step("Register supplier with NO VERI*FACTU", "gov-es.verifactu.register"), state("registered")},
				Rescue: []WorkflowStep{state("rejected")}},
			{ID: "nvf-unregister", Name: "NO VERI*FACTU unregister supplier", Description: "Release the supplier's NO VERI*FACTU seat", Channel: "verifactu", Schema: "org/party",
				Steps: []WorkflowStep{step("Unregister supplier from NO VERI*FACTU", "gov-es.verifactu.unregister"), state("empty")}},
		}},
		{Channel: "verifactu", Schema: "bill/status", Workflows: []*Workflow{
			{ID: "nvf-event", Name: "NO VERI*FACTU process event", Description: "Generate and store a NO VERI*FACTU event record", Channel: "verifactu", Schema: "bill/status",
				Steps:  []WorkflowStep{state("processing"), step("Add sequential code", "sequence.enumerate"), step("Sign envelope", "silo.close"), step("Generate event for NO VERI*FACTU", "gov-es.verifactu.event"), step("Record for NO VERI*FACTU", "gov-es.verifactu.record"), state("registered")},
				Rescue: []WorkflowStep{state("error")}},
			{ID: "nvf-summary", Name: "NO VERI*FACTU summary event", Description: "Build and record the periodic summary event", Channel: "verifactu", Schema: "bill/status",
				Steps:  []WorkflowStep{state("processing"), step("Build status for NO VERI*FACTU", "gov-es.verifactu.status"), step("Record for NO VERI*FACTU", "gov-es.verifactu.record"), state("registered")},
				Rescue: []WorkflowStep{state("error")}},
		}},
	}},
	"pdf": {Groups: []*WorkflowGroup{
		{Channel: "", Schema: "bill/invoice", Workflows: []*Workflow{
			{ID: "pdf-invoice", Name: "Generate invoice PDF", Description: "Sign the invoice and attach a PDF", Channel: "", Schema: "bill/invoice",
				Steps: []WorkflowStep{step("Sign envelope", "silo.close"), step("Generate PDF", "pdf.generate"), state("sent")}},
		}},
	}},
}
