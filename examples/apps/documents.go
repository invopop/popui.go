package apps

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// DocumentForm is the document editor's fields.
type DocumentForm struct {
	ID          string `form:"id"`
	Title       string `form:"title"`
	Description string `form:"description"`
	Schema      string `form:"schema"`
	Channel     string `form:"channel"`
	Content     string `form:"content"`
}

func (f *DocumentForm) isNew() bool { return f.ID == "" }

func documentFormFrom(doc *Document) *DocumentForm {
	return &DocumentForm{
		ID:          doc.ID,
		Title:       doc.Title,
		Description: doc.Description,
		Schema:      doc.Schema,
		Channel:     doc.Channel,
		Content:     doc.Content,
	}
}

func (f *DocumentForm) toDocument() *Document {
	return &Document{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
		Schema:      f.Schema,
		Channel:     f.Channel,
		Content:     f.Content,
	}
}

func documentsFor(app *Application) *DocumentCollection {
	if dc := sampleDocuments[app.ID]; dc != nil {
		return dc
	}
	dc := &DocumentCollection{}
	sampleDocuments[app.ID] = dc
	return dc
}

// Handlers. Unlike the app forms these mutate the in-memory fixtures:
// reordering only makes sense if the new order survives the next request.

func documentIndex(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppDocuments(sampleOrg, app, documentsFor(app), nil))
}

func documentNew(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	form := &DocumentForm{Channel: c.QueryParam("channel"), Schema: c.QueryParam("schema")}
	if form.Schema == "" {
		form.Schema = "bill/invoice"
	}
	return render(c, AppDocumentEdit(sampleOrg, app, form, nil))
}

func documentEdit(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	_, doc := documentsFor(app).Find(c.Param("did"))
	if doc == nil {
		return echo.NewHTTPError(http.StatusNotFound, "document not found")
	}
	return render(c, AppDocumentEdit(sampleOrg, app, documentFormFrom(doc), nil))
}

func documentStore(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	dc := documentsFor(app)
	form := new(DocumentForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	switch {
	case strings.TrimSpace(form.Title) == "":
		return render(c, AppDocumentEdit(sampleOrg, app, form, errors.New("title is required")))
	case !app.HasChannel(form.Channel):
		return render(c, AppDocumentEdit(sampleOrg, app, form, errors.New("unknown channel \""+form.Channel+"\" for application")))
	}
	doc := form.toDocument()
	if form.isNew() {
		doc.ID = slugify(form.Title)
		if _, exists := dc.Find(doc.ID); exists != nil {
			doc.ID += "-" + randomHex(3)
		}
		dc.Add(doc)
		success(c, "Document created.")
	} else {
		// Replace in place; moving channel or schema re-homes it to that group.
		if _, prev := dc.Find(doc.ID); prev == nil {
			return echo.NewHTTPError(http.StatusNotFound, "document not found")
		}
		dc.Remove(doc.ID)
		dc.Add(doc)
		success(c, "Document saved.")
	}
	return render(c, AppDocumentEdit(sampleOrg, app, documentFormFrom(doc), nil))
}

func documentDelete(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	dc := documentsFor(app)
	dc.Remove(c.Param("did"))
	success(c, "Document removed.")
	return render(c, AppDocuments(sampleOrg, app, dc, nil))
}

// documentReorder receives a deck's hidden order input after a drag and
// answers 204 so the deck keeps its reorder mode (hx-swap=none).
func documentReorder(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	group := documentsFor(app).Group(c.Param("channel") + "/" + c.Param("schema"))
	if group == nil {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}
	if order := c.FormValue("order"); order != "" {
		group.Reorder(strings.Split(order, ","))
	}
	return c.NoContent(http.StatusNoContent)
}
