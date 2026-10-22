// Package apps is a living, backend-free prototype of the Invopop Access
// admin "Applications" area (apps, actions, availability and credentials),
// built with PopUI. It mirrors the templates in
// invopop/access/internal/interfaces/web/components/admin so that the UI can
// be designed here and later ported back with the real backend wired in.
// See HANDOFF.md in this directory for the porting guide.
//
// Everything in this file stands in for the access domain models and the
// static vocabularies (categories, scopes, schemas…) those models expose.
// The field names match the access models so the templates read the same.
package apps

import (
	"strings"
	"time"
)

// Org stands in for models.Org.
type Org struct {
	ID   string
	Name string
}

// Channel stands in for models.Channel: a named destination within an app
// that actions attach to (e.g. "verifactu" inside "gov-es").
type Channel struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// Application visibility values, as in models.AppVisibility*.
const (
	AppVisibilityPrivate = "private"
	AppVisibilityPublic  = "public"
	AppVisibilityGated   = "gated"
)

// Gate features, as in models.FeatureCountry / models.FeatureIntegration.
const (
	FeatureCountry     = "country"
	FeatureIntegration = "integration"
)

// Transport values: how the gateway delivers tasks to the app.
const (
	TransportNATS = "nats" // NATS request/reply on gw.<app>.task (today's only mode)
	TransportHTTP = "http" // signed POST to the app's API URL, result in the response
)

// NATS connection details shared by every app; shown next to credentials.
const (
	natsURL = "tls://connect.ngs.global"
)

// Application stands in for models.Application.
type Application struct {
	ID               string
	Rev              string
	OrgID            string
	Name             string
	Description      string
	Categories       []string
	Visibility       string
	GateFeature      string
	Channels         []Channel
	ClientID         string
	ClientSecret     string
	ClientSecretHash string
	RedirectURL      string
	Scopes           []string
	Tags             []string
	IconURL          string
	LogoURL          string
	ConfigURL        string
	LaunchURL        string
	InfoURL          string
	APIURL           string
	// Transport, SigningSecretHash, SettingsSchema and AllowedOrigins are
	// proposed additions for third-party apps (see HANDOFF.md).
	Transport         string
	SigningSecret     string // never persisted; shown once
	SigningSecretHash string
	SettingsSchema    string   // JSON Schema for enrollment data
	AllowedOrigins    []string // browser origins allowed to call the API
	AllowOrgIDs       []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DisabledAt        *time.Time
}

// Disabled reports whether the app has been deactivated.
func (a *Application) Disabled() bool { return a.DisabledAt != nil }

// HasChannel mirrors models.Application.HasChannel: the empty (default)
// channel is always valid.
func (a *Application) HasChannel(key string) bool {
	if key == "" {
		return true
	}
	for _, c := range a.Channels {
		if c.Key == key {
			return true
		}
	}
	return false
}

// ChannelForActionID mirrors models.Application.ChannelForActionID: action
// IDs follow "<app_id>[.<channel>].<action>[.<qualifier>]", so the segment
// after the app id is the channel when it matches a declared one.
func (a *Application) ChannelForActionID(id string) string {
	rest, ok := strings.CutPrefix(id, a.ID+".")
	if !ok {
		return ""
	}
	seg, _, _ := strings.Cut(rest, ".")
	if seg != "" && a.HasChannel(seg) {
		return seg
	}
	return ""
}

// ChannelLabel mirrors models.Application.ChannelLabel: the human-friendly
// label for one of the app's channels, or the key itself if undeclared.
func (a *Application) ChannelLabel(key string) string {
	for _, c := range a.Channels {
		if c.Key == key {
			return c.Label
		}
	}
	return key
}

// ApplicationCollection stands in for models.ApplicationCollection.
type ApplicationCollection struct {
	List []*Application
}

// Get finds an app by ID.
func (ac *ApplicationCollection) Get(id string) *Application {
	for _, a := range ac.List {
		if a.ID == id {
			return a
		}
	}
	return nil
}

// Enabled lists the apps that are not deactivated.
func (ac *ApplicationCollection) Enabled() []*Application {
	var out []*Application
	for _, a := range ac.List {
		if !a.Disabled() {
			out = append(out, a)
		}
	}
	return out
}

// Disabled lists the deactivated apps.
func (ac *ApplicationCollection) Disabled() []*Application {
	var out []*Application
	for _, a := range ac.List {
		if a.Disabled() {
			out = append(out, a)
		}
	}
	return out
}

// ActionResult stands in for gateway.ActionResult.
type ActionResult struct {
	Status      int32  `json:"status,omitempty"` // 0 NA, 1 OK, 2 KO
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

// Action stands in for gateway.Action (the gateway protocol's action
// definition). Field names follow the protobuf-generated Go struct.
type Action struct {
	Id                  string
	Rev                 string
	ApplicationId       string
	Name                string
	Category            string
	FeatureId           string
	Pops                int32
	Seat                bool
	RenewalPeriod       string
	Unroll              bool
	Description         string
	IconUrl             string
	ConfigUrl           string
	InfoUrl             string
	ConfigSchema        string // JSON Schema for the step's config (proposed)
	Schemas             []string
	Countries           []string
	TaxRegimes          []string
	Addons              []string
	ExcludeSiloEntry    bool
	ConvertIntoCurrency string
	SharedMeta          bool
	RequireSignature    bool
	RemoveIncludedTax   bool
	OptionalSiloEntry   bool
	Results             []*ActionResult
	WorkerCount         int32
	DisabledAt          string
	CreatedAt           string
	UpdatedAt           string
}

// ActionCollection stands in for gateway.ActionCollection.
type ActionCollection struct {
	List []*Action
}

// Get finds an action by ID.
func (ac *ActionCollection) Get(id string) *Action {
	for _, a := range ac.List {
		if a.Id == id {
			return a
		}
	}
	return nil
}

// ActionGroup is one deck on the Actions tab: the actions of one channel that
// accept the same set of schemas, in display order.
type ActionGroup struct {
	Channel string
	Schemas []string // empty means every schema
	Actions []*Action
}

// SchemaKey identifies the schema set in routes: "all" or the schemas joined
// with "+" and their slashes turned into dashes.
func (g *ActionGroup) SchemaKey() string {
	if len(g.Schemas) == 0 {
		return "all"
	}
	return strings.ReplaceAll(strings.Join(g.Schemas, "+"), "/", "-")
}

// Label names the deck after its schema set.
func (g *ActionGroup) Label() string {
	switch len(g.Schemas) {
	case 0:
		return "All-schema actions"
	case 1:
		return schemaLabel(g.Schemas[0], "actions")
	}
	return "Mixed-schema actions"
}

// Groups splits the actions of an app into decks: one per channel and schema
// set, in order of first appearance, so the Actions tab can render a section
// per channel like the Workflows tab.
func (ac *ActionCollection) Groups(app *Application) []*ActionGroup {
	var out []*ActionGroup
	for _, act := range ac.List {
		ch := app.ChannelForActionID(act.Id)
		key := strings.Join(act.Schemas, "+")
		var g *ActionGroup
		for _, cand := range out {
			if cand.Channel == ch && strings.Join(cand.Schemas, "+") == key {
				g = cand
				break
			}
		}
		if g == nil {
			g = &ActionGroup{Channel: ch, Schemas: act.Schemas}
			out = append(out, g)
		}
		g.Actions = append(g.Actions, act)
	}
	return out
}

// GroupsForChannel filters Groups to one channel.
func (ac *ActionCollection) GroupsForChannel(app *Application, channel string) []*ActionGroup {
	var out []*ActionGroup
	for _, g := range ac.Groups(app) {
		if g.Channel == channel {
			out = append(out, g)
		}
	}
	return out
}

// Reorder rearranges the named actions into the given order, keeping the
// slots they occupy in the overall list so other decks are untouched.
func (ac *ActionCollection) Reorder(ids []string) {
	want := map[string]bool{}
	for _, id := range ids {
		want[id] = true
	}
	var slots []int
	byID := map[string]*Action{}
	for i, act := range ac.List {
		if want[act.Id] {
			slots = append(slots, i)
			byID[act.Id] = act
		}
	}
	if len(slots) != len(ids) {
		return
	}
	for n, id := range ids {
		ac.List[slots[n]] = byID[id]
	}
}

// Credential stands in for models.NATSCredential.
type Credential struct {
	ID            string
	Label         string
	UserPublicKey string
	Creds         string
	IssuedBy      string
	CreatedAt     time.Time
	RevokedAt     *time.Time
	RevokedBy     string
}

// Revoked reports whether the credential has been revoked.
func (c *Credential) Revoked() bool { return c.RevokedAt != nil }

// CredentialCollection stands in for models.NATSCredentialCollection.
type CredentialCollection struct {
	List []*Credential
}

// Len returns the number of credentials.
func (cc *CredentialCollection) Len() int { return len(cc.List) }

// Definition is a key with a human name and description, standing in for
// gobl's cbc.Definition as used by models.ApplicationCategories().
type Definition struct {
	Key  string
	Name string
	Desc string
}

// categories mirrors models.ApplicationCategories().
var categories = []*Definition{
	{Key: "doc", Name: "Document", Desc: "General purpose document handling."},
	{Key: "gov", Name: "Government", Desc: "Interactions with government systems or compliance."},
	{Key: "notify", Name: "Notification", Desc: "Sending notifications via email, SMS, or other channels."},
	{Key: "sync", Name: "Synchronization", Desc: "Data synchronization between systems."},
	{Key: "net", Name: "Network", Desc: "Sending and receiving data over networks."},
	{Key: "format", Name: "Format", Desc: "Data formatting and conversion."},
	{Key: "storage", Name: "Storage", Desc: "Data storage and management."},
	{Key: "payment", Name: "Payment", Desc: "Payment processing and management."},
}

// ApplicationCategories mirrors models.ApplicationCategories().
func ApplicationCategories() []*Definition { return categories }

// Scope stands in for pkg/scope: a key plus the description surfaced by
// scope.Key.Description().
type Scope struct {
	Key         string
	Description string
}

// scopes mirrors scope.Options() with scope.Key.Description().
var scopes = []Scope{
	{"auth", "Authenticate users without accessing other endpoints."},
	{"refresh", "Issue new tokens with broader access from a refresh token."},
	{"enrolled", "Manage a company's enrollment data, including silo metadata."},
	{"admin", "Administrative access across all scopes."},
	{"admin:transform", "Manage workflows, create jobs, and read results."},
	{"admin:silo", "Access and create silo entries, including attachments."},
	{"read", "Read-only access across all scopes."},
	{"read:transform", "Read-only access to workflows and jobs."},
	{"read:silo", "Read data from the silo."},
}

// ScopeOptions mirrors scope.Options().
func ScopeOptions() []Scope { return scopes }

// actionSchemas mirrors the admin package's actionSchemas.
var actionSchemas = []string{
	"bill/invoice",
	"bill/payment",
	"bill/order",
	"bill/delivery",
	"bill/status",
	"org/party",
}

// seatRenewalPeriods mirrors the admin package's seatRenewalPeriods
// (models.Interval values).
var seatRenewalPeriods = []struct {
	Value string
	Label string
}{
	{"monthly", "Monthly"},
	{"quarterly", "Quarterly"},
	{"yearly", "Yearly"},
}

// actionVerbs mirrors the admin package's controlled vocabulary for the
// action segment of an action ID.
var actionVerbs = []string{
	"register", "unregister", "send", "fetch", "sync", "cancel", "correct",
	"record", "sign", "update", "validate", "convert", "embed", "generate",
	"import", "notify", "archive", "wait", "enumerate",
}

// addons stands in for gobl's tax.AllAddonDefs(): key plus display name.
var addons = []*Definition{
	{Key: "es-verifactu-v1", Name: "Spain VERI*FACTU V1"},
	{Key: "es-tbai-v1", Name: "Spain TicketBAI"},
	{Key: "es-facturae-v3", Name: "Spain FacturaE"},
	{Key: "it-sdi-v1", Name: "Italy SDI FatturaPA v1.x"},
	{Key: "mx-cfdi-v4", Name: "Mexican SAT CFDI v4.X"},
	{Key: "br-nfse-v1", Name: "Brazil NFS-e 1.X"},
	{Key: "co-dian-v2", Name: "Colombia DIAN UBL 2.X"},
	{Key: "de-xrechnung-v3", Name: "German XRechnung 3.X"},
	{Key: "de-zugferd-v2", Name: "German ZUGFeRD 2.X"},
	{Key: "eu-en16931-v2017", Name: "EN 16931-1:2017 (EU)"},
	{Key: "fr-facturx-v1", Name: "France Factur-X v1"},
	{Key: "gr-mydata-v1", Name: "Greece MyData v1.x"},
	{Key: "pt-saft-v1", Name: "Portugal SAF-T"},
	{Key: "pl-ksef-v2", Name: "Poland KSeF v2"},
}

// currencies is a short stand-in for gobl's currency.Definitions().
var currencies = []*Definition{
	{Key: "EUR", Name: "Euro"},
	{Key: "USD", Name: "US Dollar"},
	{Key: "GBP", Name: "British Pound"},
	{Key: "MXN", Name: "Mexican Peso"},
	{Key: "BRL", Name: "Brazilian Real"},
	{Key: "COP", Name: "Colombian Peso"},
	{Key: "PLN", Name: "Polish Zloty"},
}

// ----------------------------------------------------------------------------
// Fixtures. The store is in-memory and resets whenever the server restarts.
// Only the availability toggle mutates it; form saves re-render what was
// posted without persisting anything.
// ----------------------------------------------------------------------------

var sampleOrg = &Org{
	ID:   "0192b2a5-2a1f-7e1a-9c7d-6b1b2c3d4e5f",
	Name: "Invopop",
}

func ts(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func tsp(s string) *time.Time {
	t := ts(s)
	return &t
}

var sampleApps = &ApplicationCollection{List: []*Application{
	{
		ID:          "gov-es",
		Rev:         "14-8f2c1a",
		OrgID:       sampleOrg.ID,
		Name:        "Spain Government",
		Description: "Register invoices with the Spanish tax agency through VERI*FACTU, SII and TicketBAI.",
		Categories:  []string{"gov"},
		Visibility:  AppVisibilityGated,
		GateFeature: FeatureCountry,
		Channels: []Channel{
			{Key: "verifactu", Label: "VERI*FACTU"},
			{Key: "sii", Label: "SII"},
			{Key: "ticketbai", Label: "TicketBAI"},
		},
		ClientID:         "ci_gov_es_4d1f9b7e",
		ClientSecretHash: "sha256:…",
		RedirectURL:      "https://gov-es.invopop.com/oauth/callback",
		Scopes:           []string{"enrolled", "read:silo", "admin:silo"},
		Tags:             []string{"ES"},
		ConfigURL:        "https://gov-es.invopop.com/config",
		LaunchURL:        "https://gov-es.invopop.com/launch",
		InfoURL:          "https://gov-es.invopop.com/info.md",
		APIURL:           "https://gov-es.invopop.com/api",
		Transport:        TransportNATS,
		SettingsSchema: `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Spain settings",
  "type": "object",
  "properties": {
    "environment": { "type": "string", "enum": ["production", "testing"], "default": "testing" },
    "certificate_id": { "type": "string", "title": "Certificate", "description": "ID of the uploaded signing certificate" }
  },
  "required": ["environment"]
}`,
		AllowedOrigins: []string{"https://gov-es.invopop.com"},
		CreatedAt:      ts("2023-04-12T09:14:00Z"),
		UpdatedAt:      ts("2026-09-28T16:02:00Z"),
	},
	{
		ID:               "pdf",
		Rev:              "9-c0ffee",
		OrgID:            sampleOrg.ID,
		Name:             "PDF Generator",
		Description:      "Render GOBL documents as branded PDF files with configurable templates and languages.",
		Categories:       []string{"doc", "format"},
		Visibility:       AppVisibilityPublic,
		ClientID:         "ci_pdf_7a3e52c1",
		ClientSecretHash: "sha256:…",
		RedirectURL:      "https://pdf.invopop.com/oauth/callback",
		Scopes:           []string{"enrolled", "read:silo"},
		ConfigURL:        "https://pdf.invopop.com/config",
		InfoURL:          "https://pdf.invopop.com/info.md",
		CreatedAt:        ts("2022-11-02T11:30:00Z"),
		UpdatedAt:        ts("2026-08-14T08:45:00Z"),
	},
	{
		ID:                "slack",
		Rev:               "3-1b2c3d",
		OrgID:             sampleOrg.ID,
		Name:              "Slack",
		Description:       "Post a message to a Slack channel whenever a workflow step completes or fails.",
		Categories:        []string{"notify"},
		Visibility:        AppVisibilityPublic,
		ClientID:          "ci_slack_9e8d7c6b",
		ClientSecretHash:  "sha256:…",
		RedirectURL:       "https://slack.invopop.com/oauth/callback",
		Scopes:            []string{"enrolled", "read"},
		Tags:              []string{"BETA"},
		ConfigURL:         "https://slack.invopop.com/config",
		APIURL:            "https://slack.invopop.com/api/tasks",
		Transport:         TransportHTTP,
		SigningSecretHash: "sha256:…",
		AllowedOrigins:    []string{"https://slack.invopop.com"},
		CreatedAt:         ts("2025-02-20T15:00:00Z"),
		UpdatedAt:         ts("2026-06-01T10:12:00Z"),
	},
	{
		ID:               "stripe",
		Rev:              "21-aa11bb",
		OrgID:            sampleOrg.ID,
		Name:             "Stripe",
		Description:      "Keep invoices and payments in sync with Stripe customers, subscriptions and charges.",
		Categories:       []string{"sync", "payment"},
		Visibility:       AppVisibilityGated,
		GateFeature:      FeatureIntegration,
		ClientID:         "ci_stripe_2f4e6a8c",
		ClientSecretHash: "sha256:…",
		RedirectURL:      "https://stripe.invopop.com/oauth/callback",
		Scopes:           []string{"enrolled", "admin:silo", "admin:transform"},
		ConfigURL:        "https://stripe.invopop.com/config",
		LaunchURL:        "https://stripe.invopop.com/launch",
		InfoURL:          "https://stripe.invopop.com/info.md",
		AllowOrgIDs:      []string{"0192b2a5-2a1f-7e1a-9c7d-6b1b2c3d4e5f"},
		CreatedAt:        ts("2024-07-09T13:21:00Z"),
		UpdatedAt:        ts("2026-09-30T09:00:00Z"),
	},
	{
		ID:               "acme-sandbox",
		Rev:              "2-deadbe",
		OrgID:            sampleOrg.ID,
		Name:             "Acme Sandbox",
		Description:      "Internal test application used while developing the enrollment flow. Superseded by the console sandbox.",
		Categories:       []string{"doc"},
		Visibility:       AppVisibilityPrivate,
		ClientID:         "ci_acme_0000aaaa",
		ClientSecretHash: "sha256:…",
		Scopes:           []string{"auth"},
		CreatedAt:        ts("2023-09-01T08:00:00Z"),
		UpdatedAt:        ts("2024-01-15T17:40:00Z"),
		DisabledAt:       tsp("2024-01-15T17:40:00Z"),
	},
}}

var sampleActions = map[string]*ActionCollection{
	"gov-es": {List: []*Action{
		{
			Id: "gov-es.verifactu.register", Rev: "6-a1", ApplicationId: "gov-es",
			Name: "Register with VERI*FACTU", Category: "gov", Pops: 1, Seat: true, RenewalPeriod: "monthly",
			Description: "Send the invoice to the AEAT VERI*FACTU service and record the response.",
			IconUrl:     "",
			ConfigUrl:   "https://gov-es.invopop.com/config/verifactu",
			InfoUrl:     "https://gov-es.invopop.com/info/verifactu.md",
			Schemas:     []string{"bill/invoice"}, Countries: []string{"ES"},
			Addons: []string{"es-verifactu-v1"}, ExcludeSiloEntry: true,
			ConfigSchema: `{
  "type": "object",
  "properties": {
    "representative": { "type": "boolean", "title": "Send as representative", "default": false }
  }
}`,
			Results: []*ActionResult{
				{Status: 1, Code: "accepted", Description: "Accepted by the AEAT", Default: true},
				{Status: 1, Code: "accepted-with-errors", Description: "Accepted with warnings"},
				{Status: 2, Code: "rejected", Description: "Rejected by the AEAT", Default: true},
			},
			WorkerCount: 4, CreatedAt: "2025-01-10T10:00:00Z", UpdatedAt: "2026-09-12T12:30:00Z",
		},
		{
			Id: "gov-es.verifactu.cancel", Rev: "3-b2", ApplicationId: "gov-es",
			Name: "Cancel VERI*FACTU registration", Category: "gov", Pops: 1,
			Description: "Send a cancellation record for a previously registered invoice.",
			IconUrl:     "",
			Schemas:     []string{"bill/invoice"}, Countries: []string{"ES"},
			Addons: []string{"es-verifactu-v1"}, ExcludeSiloEntry: true,
			Results: []*ActionResult{
				{Status: 1, Code: "cancelled", Description: "Cancellation accepted", Default: true},
				{Status: 2, Code: "rejected", Description: "Cancellation rejected", Default: true},
			},
			WorkerCount: 2, CreatedAt: "2025-01-10T10:05:00Z", UpdatedAt: "2026-05-02T09:10:00Z",
		},
		{
			Id: "gov-es.verifactu.unregister", Rev: "1-c3", ApplicationId: "gov-es",
			Name: "Release VERI*FACTU seat", Category: "gov", Seat: true, Unroll: true,
			Description: "Release the billing seat allocated when the workspace registered with VERI*FACTU.",
			IconUrl:     "",
			Schemas:     []string{"org/party"}, Countries: []string{"ES"}, ExcludeSiloEntry: true,
			WorkerCount: 1, CreatedAt: "2025-03-18T14:00:00Z", UpdatedAt: "2025-03-18T14:00:00Z",
		},
		{
			Id: "gov-es.sii.register", Rev: "11-d4", ApplicationId: "gov-es",
			Name: "Register with SII", Category: "gov", Pops: 1, Seat: true, RenewalPeriod: "monthly",
			Description: "Submit the invoice to the SII immediate supply of information system.",
			IconUrl:     "",
			Schemas:     []string{"bill/invoice"}, Countries: []string{"ES"}, ExcludeSiloEntry: true,
			Results: []*ActionResult{
				{Status: 1, Code: "accepted", Description: "Accepted", Default: true},
				{Status: 2, Code: "rejected", Description: "Rejected", Default: true},
			},
			WorkerCount: 4, CreatedAt: "2023-06-01T08:00:00Z", UpdatedAt: "2026-02-11T11:11:00Z",
		},
		{
			Id: "gov-es.ticketbai.register", Rev: "8-e5", ApplicationId: "gov-es",
			Name: "Register with TicketBAI", Category: "gov", Pops: 1, Seat: true, RenewalPeriod: "monthly",
			Description: "Sign and send the invoice to the Basque TicketBAI service.",
			IconUrl:     "",
			Schemas:     []string{"bill/invoice"}, Countries: []string{"ES"},
			Addons: []string{"es-tbai-v1"}, ExcludeSiloEntry: true, RequireSignature: true,
			Results: []*ActionResult{
				{Status: 1, Code: "accepted", Description: "Accepted", Default: true},
				{Status: 2, Code: "rejected", Description: "Rejected", Default: true},
			},
			WorkerCount: 2, CreatedAt: "2023-06-01T08:30:00Z", UpdatedAt: "2025-12-01T10:00:00Z",
		},
		{
			Id: "gov-es.fetch.legacy", Rev: "2-f6", ApplicationId: "gov-es",
			Name: "Fetch (legacy)", Category: "gov",
			Description: "Deprecated polling action kept for workflows created before 2024.",
			IconUrl:     "",
			WorkerCount: 1, DisabledAt: "2024-11-30T00:00:00Z",
			CreatedAt: "2023-02-01T08:00:00Z", UpdatedAt: "2024-11-30T00:00:00Z",
		},
	}},
	"pdf": {List: []*Action{
		{
			Id: "pdf.generate", Rev: "5-11", ApplicationId: "pdf",
			Name: "Generate PDF", Category: "doc", Pops: 1,
			Description:       "Render the document as a PDF and attach it to the silo entry.",
			IconUrl:           "",
			ConfigUrl:         "https://pdf.invopop.com/config",
			Schemas:           []string{"bill/invoice", "bill/order", "bill/delivery"},
			OptionalSiloEntry: false, SharedMeta: true,
			WorkerCount: 8, CreatedAt: "2022-11-02T12:00:00Z", UpdatedAt: "2026-08-14T08:45:00Z",
		},
	}},
	"slack": {List: []*Action{
		{
			Id: "slack.notify", Rev: "2-22", ApplicationId: "slack",
			Name: "Send Slack message", Category: "notify",
			Description:      "Post a message to the configured Slack channel.",
			IconUrl:          "",
			ExcludeSiloEntry: true, OptionalSiloEntry: true,
			WorkerCount: 2, CreatedAt: "2025-02-20T15:10:00Z", UpdatedAt: "2026-06-01T10:12:00Z",
		},
	}},
	"stripe": {List: []*Action{
		{
			Id: "stripe.sync", Rev: "7-33", ApplicationId: "stripe",
			Name: "Sync with Stripe", Category: "sync", Pops: 1,
			Description: "Create or update the matching Stripe invoice.",
			IconUrl:     "",
			Schemas:     []string{"bill/invoice"}, ConvertIntoCurrency: "USD",
			WorkerCount: 4, CreatedAt: "2024-07-09T13:30:00Z", UpdatedAt: "2026-09-30T09:00:00Z",
		},
		{
			Id: "stripe.fetch.payments", Rev: "4-44", ApplicationId: "stripe",
			Name: "Fetch payments", Category: "sync",
			Description: "Pull payment events for the invoice from Stripe.",
			IconUrl:     "",
			Schemas:     []string{"bill/invoice", "bill/payment"},
			WorkerCount: 2, CreatedAt: "2024-09-01T09:00:00Z", UpdatedAt: "2026-01-20T16:00:00Z",
		},
	}},
	"acme-sandbox": {},
}

var sampleCredentials = map[string]*CredentialCollection{
	"gov-es": {List: []*Credential{
		{
			ID: "cred_01", Label: "Production worker", UserPublicKey: "UAK7JX3QZ2M4N6P8R0T2V4X6Z8B0D2F4H6J8L0N2P4R6T8V0X2Z4B6D8",
			IssuedBy: "sam@invopop.com", CreatedAt: ts("2026-03-04T10:20:00Z"),
		},
		{
			ID: "cred_02", Label: "Staging worker", UserPublicKey: "UBQ2W4E6R8T0Y2U4I6O8P0A2S4D6F8G0H2J4K6L8Z0X2C4V6B8N0M2Q4",
			IssuedBy: "mark@invopop.com", CreatedAt: ts("2026-05-19T08:05:00Z"),
		},
		{
			ID: "cred_03", Label: "Laptop (Sam)", UserPublicKey: "UCZ9X7C5V3B1N9M7L5K3J1H9G7F5D3S1A9P7O5I3U1Y9T7R5E3W1Q9",
			IssuedBy: "sam@invopop.com", CreatedAt: ts("2025-11-11T11:11:00Z"),
			RevokedAt: tsp("2026-01-08T09:30:00Z"), RevokedBy: "sam@invopop.com",
		},
	}},
	"stripe": {List: []*Credential{
		{
			ID: "cred_11", Label: "Production worker", UserPublicKey: "UDM3N5B7V9C1X3Z5L7K9J1H3G5F7D9S1A3Q5W7E9R1T3Y5U7I9O1P3",
			IssuedBy: "sam@invopop.com", CreatedAt: ts("2026-02-02T12:00:00Z"),
		},
	}},
}

// issuedCredential is the fake credential returned when "Issue credentials"
// is submitted, so the one-time result modal can be designed.
func issuedCredential(app *Application, label string) *Credential {
	return &Credential{
		ID:            "cred_new",
		Label:         label,
		UserPublicKey: "UEN3W1T9H2A7N5K6S3F0R7P2R0T0T2Y1P3N6G5X8C4V2B9M1Q7Z5",
		IssuedBy:      "mark@invopop.com",
		CreatedAt:     time.Now(),
		Creds: strings.Join([]string{
			"-----BEGIN NATS USER JWT-----",
			"eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiJQUk9UT1RZUEUiLCJpYXQiOjE3NTk3NDA4MDAsImlzcyI6IkFCQyIsIm5hbWUiOiIiICsgIiIsInN1YiI6IlVFTjNXMVQ5SDJBN041SzZTM0YwUjdQMiIsIm5hdHMiOnsicHViIjp7fSwic3ViIjp7fSwic3VicyI6LTEsImRhdGEiOi0xLCJwYXlsb2FkIjotMSwidHlwZSI6InVzZXIiLCJ2ZXJzaW9uIjoyfX0.PROTOTYPE_SIGNATURE_" + app.ID,
			"------END NATS USER JWT------",
			"",
			"************************* IMPORTANT *************************",
			"NKEY Seed printed below can be used to sign and prove identity.",
			"NKEYs are sensitive and should be treated as secrets.",
			"",
			"-----BEGIN USER NKEY SEED-----",
			"SUAPROTOTYPE4N6K2S8F0R7P2R0T0T2Y1P3N6G5X8C4V2B9M1Q7Z5L3J1H9G7F5D3S1A",
			"------END USER NKEY SEED------",
			"",
			"*************************************************************",
		}, "\n"),
	}
}
