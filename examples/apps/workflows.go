package apps

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// WorkflowForm is the workflow editor's fields. Steps travel as the JSON
// readers see, edited in a monospace textarea.
type WorkflowForm struct {
	ID          string `form:"id"`
	Name        string `form:"name"`
	Description string `form:"description"`
	Channel     string `form:"channel"`
	Schema      string `form:"schema"`
	Content     string `form:"content"`
}

func (f *WorkflowForm) isNew() bool { return f.ID == "" }

func workflowFormFrom(w *Workflow) *WorkflowForm {
	return &WorkflowForm{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		Channel:     w.Channel,
		Schema:      w.Schema,
		Content:     w.Content(),
	}
}

// toWorkflow parses the steps and rescue arrays out of the content JSON; the
// name, description and schema fields on the form win over the JSON's.
func (f *WorkflowForm) toWorkflow() (*Workflow, error) {
	w := &Workflow{ID: f.ID, Name: f.Name, Description: f.Description, Channel: f.Channel, Schema: f.Schema}
	if strings.TrimSpace(f.Content) == "" {
		return w, nil
	}
	var body struct {
		Steps  []WorkflowStep `json:"steps"`
		Rescue []WorkflowStep `json:"rescue"`
	}
	if err := json.Unmarshal([]byte(f.Content), &body); err != nil {
		return nil, errors.New("content is not valid workflow JSON: " + err.Error())
	}
	w.Steps, w.Rescue = body.Steps, body.Rescue
	return w, nil
}

func workflowsFor(app *Application) *WorkflowCollection {
	if wc := sampleWorkflows[app.ID]; wc != nil {
		return wc
	}
	wc := &WorkflowCollection{}
	sampleWorkflows[app.ID] = wc
	return wc
}

// Handlers. As with documents, these mutate the in-memory fixtures so that
// reordering and adding survive the next request.

func workflowIndex(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppWorkflows(sampleOrg, app, workflowsFor(app), nil))
}

func workflowNew(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	form := &WorkflowForm{Channel: c.QueryParam("channel"), Schema: c.QueryParam("schema")}
	if form.Schema == "" {
		form.Schema = "bill/invoice"
	}
	return render(c, AppWorkflowEdit(sampleOrg, app, form, nil))
}

func workflowEdit(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	_, w := workflowsFor(app).Find(c.Param("wid"))
	if w == nil {
		return echo.NewHTTPError(http.StatusNotFound, "workflow not found")
	}
	return render(c, AppWorkflowEdit(sampleOrg, app, workflowFormFrom(w), nil))
}

func workflowStore(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	wc := workflowsFor(app)
	form := new(WorkflowForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	switch {
	case strings.TrimSpace(form.Name) == "":
		return render(c, AppWorkflowEdit(sampleOrg, app, form, errors.New("name is required")))
	case !app.HasChannel(form.Channel):
		return render(c, AppWorkflowEdit(sampleOrg, app, form, errors.New("unknown channel \""+form.Channel+"\" for application")))
	}
	w, err := form.toWorkflow()
	if err != nil {
		return render(c, AppWorkflowEdit(sampleOrg, app, form, err))
	}
	if form.isNew() {
		w.ID = slugify(form.Name)
		if _, exists := wc.Find(w.ID); exists != nil {
			w.ID += "-" + randomHex(3)
		}
		wc.Add(w)
		success(c, "Workflow created.")
	} else {
		// Replace in place; moving channel or schema re-homes it to that group.
		wc.Remove(w.ID)
		wc.Add(w)
		success(c, "Workflow saved.")
	}
	return render(c, AppWorkflowEdit(sampleOrg, app, workflowFormFrom(w), nil))
}

func workflowDelete(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	wc := workflowsFor(app)
	wc.Remove(c.Param("wid"))
	success(c, "Workflow removed.")
	return render(c, AppWorkflows(sampleOrg, app, wc, nil))
}

// workflowReorder receives a deck's hidden order input after a drag and
// answers 204 so the deck keeps its reorder mode (hx-swap=none).
func workflowReorder(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	group := workflowsFor(app).Group(c.Param("channel") + "/" + c.Param("schema"))
	if group == nil {
		return echo.NewHTTPError(http.StatusNotFound, "group not found")
	}
	if order := c.FormValue("order"); order != "" {
		group.Reorder(strings.Split(order, ","))
	}
	return c.NoContent(http.StatusNoContent)
}
