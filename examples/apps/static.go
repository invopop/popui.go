package apps

import (
	"context"
	"path"
	"sort"

	"github.com/a-h/templ"
)

// StaticPage is one GET page of the prototype rendered to a file, so the
// docs site's static build (and its Netlify previews) can serve the
// prototype without the dev server. Form submissions and the reorder
// endpoints need the server; the static copy is for browsing and review.
type StaticPage struct {
	// Path is the URL path, e.g. "examples/apps/gov-es/actions". The build
	// writes it as <Path>/index.html so the same URLs resolve on a static host.
	Path      string
	Component templ.Component
}

// StaticPages lists every page reachable by navigation from the index, for
// every fixture app, in a stable order.
func StaticPages() []StaticPage {
	org := sampleOrg
	var pages []StaticPage
	add := func(c templ.Component, segs ...string) {
		pages = append(pages, StaticPage{Path: path.Join(append([]string{"examples", "apps"}, segs...)...), Component: c})
	}
	add(AppsIndex(org, sampleApps, false))
	add(AppsEdit(org, newApplicationForm(), nil), "new")
	for _, app := range sampleApps.List {
		form := ApplicationModelToForm(app)
		add(AppsEdit(org, form, nil), app.ID)
		add(AppListing(org, app, listingFormFrom(listingFor(app)), nil), app.ID, "listing")
		add(AppsEditAvailability(org, form, nil), app.ID, "availability")
		add(AppIntegration(org, app, form, true, true, credentialsFor(app), nil, nil), app.ID, "integration")

		ac := actionsFor(app)
		add(AppActions(org, app, ac, nil), app.ID, "actions")
		add(AppActionEdit(org, app, &Action{ApplicationId: app.ID}, nil), app.ID, "actions", "new")
		for _, act := range ac.List {
			add(AppActionEdit(org, app, act, nil), app.ID, "actions", act.Id)
		}

		wc := workflowsFor(app)
		add(AppWorkflows(org, app, wc, nil), app.ID, "workflows")
		add(AppWorkflowEdit(org, app, &WorkflowForm{Schema: "bill/invoice"}, nil), app.ID, "workflows", "new")
		for _, g := range wc.Groups {
			for _, w := range g.Workflows {
				add(AppWorkflowEdit(org, app, workflowFormFrom(w), nil), app.ID, "workflows", w.ID)
			}
		}

		dc := documentsFor(app)
		add(AppDocuments(org, app, dc, nil), app.ID, "documents")
		add(AppDocumentEdit(org, app, &DocumentForm{Schema: "bill/invoice"}, nil), app.ID, "documents", "new")
		for _, g := range dc.Groups {
			for _, d := range g.Documents {
				add(AppDocumentEdit(org, app, documentFormFrom(d), nil), app.ID, "documents", d.ID)
			}
		}
	}
	sort.SliceStable(pages, func(i, j int) bool { return pages[i].Path < pages[j].Path })
	return pages
}

// StaticContext is the context the static pages render with: a sudo viewer
// and no HTMX request, so each file is a complete document.
func StaticContext() context.Context {
	return withSudo(context.Background(), true)
}
