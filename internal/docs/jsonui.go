package docs

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/jsonui"
)

// jsonuiDetailsExample is a data-driven card: every value comes from "data"
// through "bind".
const jsonuiDetailsExample = `{
  "data": {
    "invoice": {
      "code": "INV-2026-0142",
      "issued": "Issued 27 Oct 2026",
      "status": "Sent",
      "customer": "Sample Consumer S.A.",
      "total": "€1,815.00"
    }
  },
  "components": [
    {
      "is": "card",
      "children": [
        { "is": "card_header", "bind": { "title": "invoice.code", "subtitle": "invoice.issued" } },
        {
          "is": "card_content",
          "children": [
            {
              "is": "description_list",
              "children": [
                { "is": "description_list_item", "props": { "label": "Status", "inline": true }, "bind": { "value": "invoice.status" } },
                { "is": "description_list_item", "props": { "label": "Customer", "inline": true }, "bind": { "value": "invoice.customer" } },
                { "is": "description_list_item", "props": { "label": "Total", "inline": true }, "bind": { "value": "invoice.total" } }
              ]
            }
          ]
        }
      ]
    }
  ]
}`

// jsonuiFormExample binds form controls to the data with "model" and runs
// built-in actions from "on", one of them through a named handler.
const jsonuiFormExample = `{
  "data": { "customer": { "name": "Acme S.L.", "email": "billing@acme.com", "einvoice": true } },
  "handlers": {
    "save": { "action": "alert", "message": "Saving {{customer.name}} <{{customer.email}}>" }
  },
  "components": [
    {
      "is": "form",
      "on": { "submit.prevent": "save" },
      "children": [
        { "is": "input", "props": { "name": "name", "label": "Legal name" }, "model": "customer.name" },
        { "is": "input", "props": { "name": "email", "type": "email", "label": "Email" }, "model": "customer.email" },
        { "is": "checkbox", "props": { "name": "einvoice", "variant": "switch", "label": "Send e-invoices" }, "model": "customer.einvoice" },
        { "is": "description", "children": [{ "is": "text", "text": "Hello, " }, { "is": "text", "bind": "customer.name" }] },
        {
          "is": "button_group",
          "props": { "align": "right" },
          "children": [
            { "is": "button", "props": { "variant": "secondary", "type": "button" }, "on": { "click": { "action": "reset" } }, "children": [{ "is": "text", "text": "Reset" }] },
            { "is": "button", "children": [{ "is": "text", "text": "Save" }] }
          ]
        }
      ]
    }
  ]
}`

const jsonuiGoExample = `import "github.com/invopop/popui.go/jsonui"

// Strict by default: only the built-in actions plus the ones listed here.
opts := jsonui.Options{
	Actions: append(jsonui.BuiltinActions, "ui.notify"),
}

// Render writes the HTML for a document straight to a response...
if err := opts.Render(ctx, body, w); err != nil {
	return err // invalid or unsafe documents are rejected with an error
}

// ...or Build returns a templ.Component to use inside other templates.
c, err := opts.Build(body)`

const jsonuiAssetsExample = `// jsonui.js runs the actions documents reference, and must be on the page
// alongside PopUI's own assets and Alpine.js.
e.StaticFS("/_jsonui", echo.MustSubFS(jsonui.Assets, "assets"))
// <script src="/_jsonui/jsonui.js"></script>`

const jsonuiListsExample = `{
  "data": {
    "config": { "connection_id": "" },
    "connections": [],
    "loading": { "connections": false }
  },
  "handlers": {
    "load": { "action": "app.request", "path": "/connections", "into": "connections", "busy": "loading.connections" }
  },
  "init": "load",
  "components": [
    {
      "is": "select",
      "props": { "label": "Connection" },
      "model": "config.connection_id",
      "bind": { "disabled": "loading.connections" },
      "children": [
        { "is": "option", "props": { "value": "", "label": "Default" } },
        { "is": "option", "each": "connections", "bind": { "value": "item.id", "label": "item.name" } }
      ]
    },
    { "is": "p", "show": "loading.connections", "children": [{ "is": "text", "text": "Loading…" }] }
  ]
}`

const jsonuiActionsExample = `<script src="/_jsonui/jsonui.js"></script>
<script>
  // Host actions receive the interpolated params, the document's reactive
  // data, a snapshot() of it, the DOM event and the element.
  popui.actions.register("ui.notify", ({ params }) => showToast(params.message));

  // An action's return value is stored at its "into" path, and its "busy"
  // path is true while a promise it returns is pending.
  popui.actions.register("asset.choose", ({ params }) => chooseFile(params.accept));
</script>`

const jsonuiBrowserExample = `<link rel="stylesheet" href="https://assets.invopop.com/popui/v0.75.0/popui.css">
<script src="https://assets.invopop.com/popui/v0.75.0/popui.js"></script>
<script src="https://assets.invopop.com/popui/v0.75.0/jsonui.js"></script>
<script src="https://cdn.jsdelivr.net/npm/alpinejs@3.15.0/dist/cdn.min.js" defer></script>
<script src="https://assets.invopop.com/popui/v0.75.0/wasm_exec.js"></script>
<script type="module">
  const go = new Go();
  const ready = new Promise((r) => addEventListener("popui:ready", r, { once: true }));
  const wasm = fetch("https://assets.invopop.com/popui/v0.75.0/popui.wasm");
  const { instance } = await WebAssembly.instantiateStreaming(wasm, go.importObject);
  go.run(instance);
  await ready;

  const { html, error } = popuiRender(JSON.stringify(doc), {
    actions: ["alert", "set", "reset", "dispatch", "ui.notify"],
  });
  if (error) throw new Error(error);
  document.getElementById("app").innerHTML = html; // Alpine initialises it
</script>`

// jsonuiRender renders a docs example document with the default, strict
// options, so the examples also show what strict mode accepts.
func jsonuiRender(doc string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return jsonui.Options{}.Render(ctx, []byte(doc), w)
	})
}
