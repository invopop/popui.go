package apps

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestRoutes renders every page of the prototype through the real router so
// a template that panics or a handler that 500s fails the build, and checks
// a few landmarks that the access port must keep.
func TestRoutes(t *testing.T) {
	e := echo.New()
	Register(e)

	get := func(t *testing.T, path string) string {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: status %d\n%s", path, rec.Code, rec.Body.String())
		}
		return rec.Body.String()
	}
	send := func(t *testing.T, method, path string, form url.Values) string {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		req.Header.Set("HX-Request", "true")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s %s: status %d\n%s", method, path, rec.Code, rec.Body.String())
		}
		return rec.Body.String()
	}
	want := func(t *testing.T, body string, parts ...string) {
		t.Helper()
		for _, p := range parts {
			if !strings.Contains(body, p) {
				t.Errorf("missing %q", p)
			}
		}
	}

	t.Run("index", func(t *testing.T) {
		body := get(t, "/examples/apps")
		want(t, body, "Manage applications", "Spain Government", "gov-es", "Create app")
		if strings.Contains(body, "Acme Sandbox") {
			t.Error("deactivated app listed under active apps")
		}
		body = get(t, "/examples/apps?deactivated=1")
		want(t, body, "Acme Sandbox")
	})

	t.Run("new and create", func(t *testing.T) {
		body := get(t, "/examples/apps/new")
		want(t, body, "Create app", `id="app-form"`, `form="app-form"`)
		body = send(t, http.MethodPost, "/examples/apps", url.Values{
			"id": {"demo"}, "name": {"Demo"}, "categories": {"doc", "sync"}, "client_secret": {"cs_once"},
		})
		// Creation lands on the Integration tab with the one-time secrets.
		want(t, body, "Application created.", "OAuth client secret", "cs_once", `form="app-integration-form"`)
		body = get(t, "/examples/apps/gov-es/integration")
		if strings.Contains(body, `name="client_secret"`) {
			t.Error("secret shown again on a later visit")
		}
	})

	t.Run("edit and update", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es")
		want(t, body, "Spain Government", "VERI*FACTU", `name="categories"`)
		if strings.Contains(body, `name="transport"`) || strings.Contains(body, "OAuth client ID") || strings.Contains(body, `name="scopes"`) {
			t.Error("delivery and OAuth belong on the Integration tab, not the App tab")
		}
		body = send(t, http.MethodPut, "/examples/apps/gov-es", url.Values{
			"rev": {"14-8f2c1a"}, "name": {"Spain Government (edited)"},
		})
		want(t, body, "Application updated.", "Spain Government (edited)")
	})

	t.Run("listing", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es/listing")
		want(t, body, "Spain Government", `maxlength="1000"`, "Register invoices with the AEAT in real time", `name="feature_order"`, "listingFeatures([true,true,true])", "No key features yet", `name="countries"`, `name="privacy_url"`, `name="support_email"`)
		if strings.Contains(body, "Works with") || strings.Contains(body, "Response time") || strings.Contains(body, "Pricing model") {
			t.Error("listing must not show Works with or Response time")
		}
		body = send(t, http.MethodPut, "/examples/apps/gov-es/listing", url.Values{
			"about": {"Short about"}, "publisher": {"Invopop"},
			"feature_titles[]":       {"First", "", "Third"},
			"feature_descriptions[]": {"d1", "", "d3"},
			"feature_image_urls[]":   {"", "", ""},
			"feature_order":          {"feature-3,feature-1,feature-2"},
		})
		want(t, body, "Listing saved.", ">Short about<")
		// The deck order and the empty slot are applied: Third, First, then an empty slot.
		if i, j := strings.Index(body, `value="Third"`), strings.Index(body, `value="First"`); i < 0 || j < 0 || i > j {
			t.Errorf("feature order not applied (Third at %d, First at %d)", i, j)
		}
	})

	t.Run("availability", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es/availability")
		want(t, body, "Danger zone", "Disable application", "Allowed org IDs", "Visibility")
		body = get(t, "/examples/apps/gov-es/availability?sudo=0")
		if strings.Contains(body, "Allowed org IDs") {
			t.Error("sudo section rendered for a non-sudo viewer")
		}
		body = send(t, http.MethodPut, "/examples/apps/pdf/availability", url.Values{"disabled": {"true"}})
		want(t, body, "Enable application", "Application disabled.")
		body = send(t, http.MethodPut, "/examples/apps/pdf/availability", url.Values{})
		want(t, body, "Disable application", "Application enabled.")
		body = send(t, http.MethodPut, "/examples/apps/gov-es/sudo", url.Values{
			"visibility": {"public"}, "allow_org_ids[]": {"0192b2a5-2a1f-7e1a-9c7d-6b1b2c3d4e5f"},
		})
		want(t, body, "Saved", `value="public" selected`)
	})

	t.Run("actions", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es/actions")
		want(t, body, "Spain Government", "gov-es.verifactu.register", "VERI*FACTU", "Seat · monthly", "Releases seat", "Default channel", "All-schema actions", "Invoice actions", `name="order"`, "New action")
		req := httptest.NewRequest(http.MethodPost, "/examples/apps/gov-es/actions/groups/verifactu/bill-invoice/reorder", strings.NewReader("order=gov-es.verifactu.cancel,gov-es.verifactu.register"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("reorder: status %d", rec.Code)
		}
		if got := sampleActions["gov-es"].List[0].Id; got != "gov-es.verifactu.cancel" {
			t.Errorf("reorder not applied, first action is %s", got)
		}
		body = send(t, http.MethodPost, "/examples/apps/gov-es/actions/gov-es.sii.register/disable", nil)
		want(t, body, "Action disabled.")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/actions/gov-es.sii.register/enable", nil)
		want(t, body, "Action enabled.")
		body = get(t, "/examples/apps/gov-es/actions/new")
		want(t, body, "Action ID", "actionId(&#39;gov-es&#39;)", "Other…")
		body = get(t, "/examples/apps/gov-es/actions/gov-es.verifactu.register")
		want(t, body, "Register with VERI*FACTU", "accepted-with-errors", "es-verifactu-v1", "Worker count",
			`name="config_schema"`, "QUEUED", "Task (request)", "retried with back-off")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/actions", url.Values{
			"channel": {"sii"}, "action_verb": {"send"}, "qualifier": {"batch"}, "name": {"Send SII batch"},
			"results[0].status": {"1"}, "results[0].code": {"ok"},
		})
		want(t, body, "Saved", "gov-es.sii.send.batch", "Send SII batch", "&#34;code&#34;:&#34;ok&#34;")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/actions", url.Values{
			"channel": {"nope"}, "action_verb": {"send"}, "name": {"x"},
		})
		want(t, body, "unknown channel")
	})

	t.Run("workflows", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es/workflows")
		want(t, body, "Spain Government", "VERI*FACTU", "SII issue invoice", "Event workflows", "Default channel", "Facturae generate invoice", `name="order"`)
		body = get(t, "/examples/apps/gov-es/workflows/new?channel=sii&schema=org/party")
		want(t, body, "New workflow", `value="sii" selected`, `value="org/party" selected`)
		body = get(t, "/examples/apps/gov-es/workflows/sii-register")
		want(t, body, "SII register party", "gov-es.sii.wait.upload")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/workflows", url.Values{
			"name": {"SII correct invoice"}, "channel": {"sii"}, "schema": {"bill/invoice"},
			"content": {`{"steps":[{"name":"Correct","provider":"gov-es.sii.correct"}]}`},
		})
		want(t, body, "Workflow created.", `value="sii-correct-invoice"`)
		body = send(t, http.MethodPost, "/examples/apps/gov-es/workflows", url.Values{
			"name": {"Bad"}, "channel": {"sii"}, "schema": {"bill/invoice"}, "content": {"{nope"},
		})
		want(t, body, "not valid workflow JSON")
		req := httptest.NewRequest(http.MethodPost, "/examples/apps/gov-es/workflows/groups/sii/bill-invoice/reorder", strings.NewReader("order=sii-received,sii-issued"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("reorder: status %d", rec.Code)
		}
		if got := sampleWorkflows["gov-es"].Group("sii/bill-invoice").Workflows[0].ID; got != "sii-received" {
			t.Errorf("reorder not applied, first workflow is %s", got)
		}
		body = send(t, http.MethodPost, "/examples/apps/gov-es/workflows/sii-correct-invoice/delete", nil)
		want(t, body, "Workflow removed.")
	})

	t.Run("documents", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es/documents")
		want(t, body, "Spain Government", "SII", "Invoice documents", "Party documents", "Default channel", "B2B Standard Invoice", "Facturae B2G Invoice", `name="order"`, `hx-disinherit="*"`)
		if strings.Contains(body, "New category") {
			t.Error("documents are grouped by channel, there are no categories")
		}
		body = get(t, "/examples/apps/gov-es/documents/new?channel=sii&schema=org/party")
		want(t, body, "New document", `value="sii" selected`, `value="org/party" selected`)
		body = get(t, "/examples/apps/gov-es/documents/sii-b2c")
		want(t, body, "B2C Simplified Invoice", "es-sii-v1")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/documents", url.Values{
			"title": {"B2B Rounding Example"}, "channel": {"sii"}, "schema": {"bill/invoice"}, "content": {"{}"},
		})
		want(t, body, "Document created.", `value="b2b-rounding-example"`)
		req := httptest.NewRequest(http.MethodPost, "/examples/apps/gov-es/documents/groups/sii/bill-invoice/reorder", strings.NewReader("order=sii-b2b,sii-b2c"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("reorder: status %d", rec.Code)
		}
		if got := sampleDocuments["gov-es"].Group("sii/bill-invoice").Documents[0].ID; got != "sii-b2b" {
			t.Errorf("reorder not applied, first doc is %s", got)
		}
		body = send(t, http.MethodPost, "/examples/apps/gov-es/documents/b2b-rounding-example/delete", nil)
		want(t, body, "Document removed.")
	})

	t.Run("integration", func(t *testing.T) {
		body := get(t, "/examples/apps/gov-es/integration")
		want(t, body, "Spain Government", "How to connect", "tls://connect.ngs.global", "gw.gov-es.task",
			`name="transport"`, "OAuth client ID", "oauth/token", `name="scopes"`, "enrolled", "Reset secret", "How the app gets a token",
			`name="settings_schema"`, `name="allowed_origins[]"`,
			"Production worker", "Revoked", "Issue credentials")
		body = get(t, "/examples/apps/slack/integration")
		want(t, body, "HTTP webhook", "Invopop-Signature", "Rotate secret")
		if strings.Contains(body, "Issue credentials") {
			t.Error("HTTP apps do not issue NATS credentials")
		}
		body = send(t, http.MethodPut, "/examples/apps/gov-es/integration", url.Values{
			"transport": {"http"}, "api_url": {""},
		})
		want(t, body, "API URL is required")
		body = send(t, http.MethodPut, "/examples/apps/gov-es/integration", url.Values{
			"transport": {"http"}, "api_url": {"https://gov-es.invopop.com/tasks"}, "client_id": {"ci_gov_es_4d1f9b7e"},
		})
		want(t, body, "Integration settings saved.", `value="http" selected`)
		body = send(t, http.MethodPost, "/examples/apps/slack/integration/rotate-secret", nil)
		want(t, body, "Signing secret rotated", "whsec_")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/integration/reset-client-secret", nil)
		want(t, body, "OAuth client secret reset", "cs_", `name="client_secret"`)
		body = send(t, http.MethodPost, "/examples/apps/gov-es/integration/credentials", url.Values{"label": {"CI runner"}})
		want(t, body, "Credentials issued", "BEGIN NATS USER JWT", "credential-result")
		body = send(t, http.MethodPost, "/examples/apps/gov-es/integration/credentials/cred_01/revoke", nil)
		want(t, body, "Credential revoked.")
	})

	t.Run("htmx request renders only the body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/examples/apps", nil)
		req.Header.Set("HX-Request", "true")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if strings.Contains(rec.Body.String(), "<html") {
			t.Error("full document returned for an HTMX request")
		}
	})
}
