package apps

import (
	"errors"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go"
	"github.com/invopop/popui.go/flash"
	"github.com/invopop/popui.go/htmx"
	"github.com/labstack/echo/v4"
)

// Prefix is where the prototype is mounted by `popui serve`.
const Prefix = "/examples"

// Register mounts the prototype routes. The paths mirror the access admin
// controller (adminAppsController.serve) minus the org segment, so every
// template URL helper maps 1:1 to a real admin route:
//
//	GET  /apps                      index (?deactivated=1 for the other tab)
//	GET  /apps/new                  create form
//	POST /apps                      create (re-renders with a flash)
//	GET  /apps/:id                  config form
//	PUT  /apps/:id                  update (re-renders with a flash)
//	GET  /apps/:id/listing          marketplace listing form
//	PUT  /apps/:id/listing          save listing (re-renders with a flash)
//	GET  /apps/:id/availability     availability tab
//	PUT  /apps/:id/availability     enable / disable (toggles the fixture)
//	PUT  /apps/:id/sudo             sudo visibility save
//	GET  /apps/:id/actions          actions list
//	GET  /apps/:id/actions/new      new action form
//	GET  /apps/:id/actions/:aid     edit action form
//	POST /apps/:id/actions          store action (re-renders with a flash)
//	POST /apps/:id/actions/:aid/disable | enable   sudo toggle (mutates fixture)
//	POST /apps/:id/actions/groups/:channel/:schema/reorder  order=<ids> from the deck
//	GET  /apps/:id/workflows        example workflows: a section per channel, a deck per schema
//	GET  /apps/:id/workflows/new    new workflow form (?channel=&schema= preselect)
//	POST /apps/:id/workflows        create / update a workflow
//	GET  /apps/:id/workflows/:wid   edit workflow form
//	POST /apps/:id/workflows/:wid/delete
//	POST /apps/:id/workflows/groups/:channel/:schema/reorder  order=<ids> from the deck
//	GET  /apps/:id/documents        example documents: a section per channel, a deck per schema
//	GET  /apps/:id/documents/new    new document form (?channel=&schema= preselect)
//	POST /apps/:id/documents        create / update a document
//	GET  /apps/:id/documents/:did   edit document form
//	POST /apps/:id/documents/:did/delete
//	POST /apps/:id/documents/groups/:channel/:schema/reorder  order=<ids> from the deck
//	GET  /apps/:id/integration      delivery, credentials, OAuth and enrollment
//	PUT  /apps/:id/integration      save delivery / OAuth / enrollment (re-renders with a flash)
//	POST /apps/:id/integration/rotate-secret        new HTTP signing secret, shown once
//	POST /apps/:id/integration/reset-client-secret  new OAuth client secret, shown once
//	POST /apps/:id/integration/credentials          issue NATS creds (one-time result modal)
//	POST /apps/:id/integration/credentials/:cid/revoke
func Register(e *echo.Echo) {
	g := e.Group(Prefix + "/apps")
	g.GET("", index)
	g.GET("/new", newApp)
	g.POST("", create)
	g.GET("/:id", edit)
	g.PUT("/:id", update)
	g.GET("/:id/listing", listing)
	g.PUT("/:id/listing", updateListing)
	g.GET("/:id/availability", availability)
	g.PUT("/:id/availability", updateAvailability)
	g.PUT("/:id/sudo", updateSudo)
	g.GET("/:id/actions", actionIndex)
	g.GET("/:id/actions/new", actionNew)
	g.GET("/:id/actions/:aid", actionEdit)
	g.POST("/:id/actions", actionStore)
	g.POST("/:id/actions/:aid/disable", actionSetDisabled(true))
	g.POST("/:id/actions/:aid/enable", actionSetDisabled(false))
	g.POST("/:id/actions/groups/:channel/:schema/reorder", actionReorder)
	g.GET("/:id/workflows", workflowIndex)
	g.GET("/:id/workflows/new", workflowNew)
	g.POST("/:id/workflows", workflowStore)
	g.GET("/:id/workflows/:wid", workflowEdit)
	g.POST("/:id/workflows/:wid/delete", workflowDelete)
	g.POST("/:id/workflows/groups/:channel/:schema/reorder", workflowReorder)
	g.GET("/:id/documents", documentIndex)
	g.GET("/:id/documents/new", documentNew)
	g.POST("/:id/documents", documentStore)
	g.GET("/:id/documents/:did", documentEdit)
	g.POST("/:id/documents/:did/delete", documentDelete)
	g.POST("/:id/documents/groups/:channel/:schema/reorder", documentReorder)
	g.GET("/:id/integration", integration)
	g.PUT("/:id/integration", updateIntegration)
	g.POST("/:id/integration/rotate-secret", rotateSecret)
	g.POST("/:id/integration/reset-client-secret", resetClientSecret)
	g.POST("/:id/integration/credentials", issueCredentials)
	g.POST("/:id/integration/credentials/:cid/revoke", revokeCredential)
}

// mainClass gives every page's content area breathing room at the bottom so
// the last section is not flush against the viewport edge when scrolled.
const mainClass = "pb-16"

// tooltipDelay is how long the pointer rests on a field's (i) before its
// tooltip opens, in milliseconds. Without it the cards pop up on every pass
// of the mouse.
const tooltipDelay = 300

// baseURL and orgURL mirror the access helpers: there, orgURL(org, "apps", id)
// yields "/admin/<org>/apps/<id>"; here the org segment is dropped.
func baseURL(_ *Org, paths ...string) string {
	return path.Join(append([]string{Prefix}, paths...)...)
}

func orgURL(org *Org, paths ...string) string {
	return baseURL(org, paths...)
}

// navAttrs are the HTMX attributes that make an element navigate to url by
// swapping the app container and pushing the URL, exactly as access does on
// every tab, row and breadcrumb.
func navAttrs(url string) templ.Attributes {
	return templ.Attributes{
		"hx-get":      url,
		"hx-target":   popui.AppTarget,
		"hx-push-url": "true",
	}
}

func render(c echo.Context, t templ.Component) error {
	ctx := htmx.WithContext(c)
	ctx = flash.WithContext(ctx, c)
	ctx = withSudo(ctx, c.QueryParam("sudo") != "0")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return t.Render(ctx, c.Response().Writer)
}

func success(c echo.Context, text string) {
	flash.SetMessage(c, &flash.Message{Type: "success", Text: text})
}

func findApp(c echo.Context) (*Application, error) {
	app := sampleApps.Get(c.Param("id"))
	if app == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "application not found")
	}
	return app, nil
}

func index(c echo.Context) error {
	showDeactivated := c.QueryParam("deactivated") != ""
	return render(c, AppsIndex(sampleOrg, sampleApps, showDeactivated))
}

func newApp(c echo.Context) error {
	return render(c, AppsEdit(sampleOrg, newApplicationForm(), nil))
}

func edit(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppsEdit(sampleOrg, ApplicationModelToForm(app), nil))
}

// create binds the posted form and renders it back as if it had been stored:
// a rev and timestamps appear, the OAuth secret disappears, and the flash
// confirms. Nothing is persisted.
func create(c echo.Context) error {
	form := new(ApplicationForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	if form.ID == "" || form.Name == "" {
		return render(c, AppsEdit(sampleOrg, form, errors.New("id and name are required")))
	}
	now := timeToString(time.Now())
	form.Rev = "1-" + randomHex(6)
	form.CreatedAt = now
	form.UpdatedAt = now
	// Land on the Integration tab, where the one-time OAuth and signing
	// secrets are shown once. The posted form stands in for the stored app.
	app := &Application{ID: form.ID, Name: form.Name, Transport: form.transportOrDefault(), APIURL: form.APIURL}
	success(c, "Application created. Copy the secrets below now; they will not be shown again.")
	return render(c, AppIntegration(sampleOrg, app, form, true, true, &CredentialCollection{}, nil, nil))
}

func update(c echo.Context) error {
	if _, err := findApp(c); err != nil {
		return err
	}
	form := new(ApplicationForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	form.ID = c.Param("id")
	form.UpdatedAt = timeToString(time.Now())
	success(c, "Application updated.")
	return render(c, AppsEdit(sampleOrg, form, nil))
}

func availability(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppsEditAvailability(sampleOrg, ApplicationModelToForm(app), nil))
}

// updateAvailability is the one handler that mutates the in-memory fixture,
// so the Enable/Disable button flips and the index tabs move the app between
// lists. Restarting the server resets it.
func updateAvailability(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	if c.FormValue("disabled") == "true" {
		now := time.Now()
		app.DisabledAt = &now
		success(c, "Application disabled.")
	} else {
		app.DisabledAt = nil
		success(c, "Application enabled.")
	}
	return render(c, AppsEditAvailability(sampleOrg, ApplicationModelToForm(app), nil))
}

func updateSudo(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	posted := new(ApplicationForm)
	if err := c.Bind(posted); err != nil {
		return err
	}
	form := ApplicationModelToForm(app)
	form.Visibility = posted.Visibility
	form.GateFeature = posted.GateFeature
	form.AllowOrgIDs = posted.AllowOrgIDs
	success(c, "Saved")
	return render(c, AppsEditAvailability(sampleOrg, form, nil))
}

func actionsFor(app *Application) *ActionCollection {
	if ac := sampleActions[app.ID]; ac != nil {
		return ac
	}
	return &ActionCollection{}
}

func actionIndex(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppActions(sampleOrg, app, actionsFor(app), nil))
}

func actionNew(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppActionEdit(sampleOrg, app, &Action{ApplicationId: app.ID}, nil))
}

func actionEdit(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	act := actionsFor(app).Get(c.Param("aid"))
	if act == nil {
		return echo.NewHTTPError(http.StatusNotFound, "action not found")
	}
	return render(c, AppActionEdit(sampleOrg, app, act, nil))
}

func actionStore(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	form := new(ActionForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	form.Results = bindResults(c.Request().PostForm)

	if !app.HasChannel(form.Channel) {
		err := errors.New("unknown channel \"" + form.Channel + "\" for application")
		return render(c, AppActionEdit(sampleOrg, app, form.ToAction(app.ID), err))
	}
	if form.isNew() && form.ActionVerb == "" {
		err := errors.New("action segment is required")
		return render(c, AppActionEdit(sampleOrg, app, form.ToAction(app.ID), err))
	}
	action := form.ToAction(app.ID)
	if action.Category == "" && len(app.Categories) > 0 {
		action.Category = app.Categories[0]
	}
	now := timeToString(time.Now())
	if action.Rev == "" {
		action.Rev = "1-" + randomHex(6)
		action.CreatedAt = now
	}
	action.UpdatedAt = now
	success(c, "Saved")
	return render(c, AppActionEdit(sampleOrg, app, action, nil))
}

// actionSetDisabled flips an action's disabled flag in the fixtures, the way
// the sudo checkbox on the action form does in access.
func actionSetDisabled(disabled bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		app, err := findApp(c)
		if err != nil {
			return err
		}
		ac := actionsFor(app)
		act := ac.Get(c.Param("aid"))
		if act == nil {
			return echo.NewHTTPError(http.StatusNotFound, "action not found")
		}
		if disabled {
			act.DisabledAt = timeToString(time.Now())
			success(c, "Action disabled.")
		} else {
			act.DisabledAt = ""
			success(c, "Action enabled.")
		}
		return render(c, AppActions(sampleOrg, app, ac, nil))
	}
}

// actionReorder receives a deck's hidden order input after a drag and answers
// 204 so the deck keeps its reorder mode (hx-swap=none).
func actionReorder(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	if order := c.FormValue("order"); order != "" {
		actionsFor(app).Reorder(strings.Split(order, ","))
	}
	return c.NoContent(http.StatusNoContent)
}

func credentialsFor(app *Application) *CredentialCollection {
	if cc := sampleCredentials[app.ID]; cc != nil {
		return cc
	}
	return &CredentialCollection{}
}

func integration(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return renderIntegration(c, app, ApplicationModelToForm(app), nil, nil)
}

// updateIntegration overlays the posted delivery, OAuth and enrollment fields
// on the stored app and re-renders, as the app form does.
func updateIntegration(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	posted := new(ApplicationForm)
	if err := c.Bind(posted); err != nil {
		return err
	}
	form := ApplicationModelToForm(app)
	form.Transport = posted.Transport
	form.APIURL = posted.APIURL
	form.ClientID = posted.ClientID
	form.RedirectURL = posted.RedirectURL
	form.Scopes = posted.Scopes
	form.SettingsSchema = posted.SettingsSchema
	form.AllowedOrigins = posted.AllowedOrigins
	if form.Transport == TransportHTTP && form.APIURL == "" {
		return renderIntegration(c, app, form, nil, errors.New("an API URL is required for HTTP delivery"))
	}
	success(c, "Integration settings saved.")
	return renderIntegration(c, app, form, nil, nil)
}

// rotateSecret mints a new HTTP signing secret and shows it once.
func rotateSecret(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	form := ApplicationModelToForm(app)
	form.SigningSecret = "whsec_" + randomHex(24)
	success(c, "Signing secret rotated. Copy it now; it will not be shown again.")
	return renderIntegration(c, app, form, nil, nil)
}

// resetClientSecret replaces a lost OAuth client secret and shows the new
// one once; the old one stops working as soon as the hash is overwritten.
func resetClientSecret(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	form := ApplicationModelToForm(app)
	form.ClientSecret = "cs_" + randomHex(24)
	form.ClientSecretHash = "sha256:" + randomHex(16)
	success(c, "OAuth client secret reset. Copy it now; it will not be shown again.")
	return renderIntegration(c, app, form, nil, nil)
}

func renderIntegration(c echo.Context, app *Application, form *ApplicationForm, issued *Credential, renderErr error) error {
	return render(c, AppIntegration(sampleOrg, app, form, true, true, credentialsFor(app), issued, renderErr))
}

func issueCredentials(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	label := c.FormValue("label")
	if label == "" {
		label = "Untitled"
	}
	issued := issuedCredential(app, label)
	return renderIntegration(c, app, ApplicationModelToForm(app), issued, nil)
}

func revokeCredential(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	success(c, "Credential revoked.")
	return renderIntegration(c, app, ApplicationModelToForm(app), nil, nil)
}
