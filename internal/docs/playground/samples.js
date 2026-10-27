// Sample JSON documents for the jsonui playground.
const SCHEMA = "https://raw.githubusercontent.com/invopop/popui.go/main/jsonui/schema/document.json";
window.SCHEMA_ID = SCHEMA;

const text = (t) => ({ type: "text", text: t });

window.SAMPLES = {
  // Data-driven: every value comes from "data" through "bind" (prop name →
  // data path), resolved by the renderer so the HTML is complete up front.
  details: {
    $schema: SCHEMA,
    data: {
      invoice: {
        code: "INV-2026-0142",
        issued: "Issued 27 Oct 2026",
        status: { label: "Sent", state: "success" },
        country: "Spain",
        regime: "VeriFactu",
        supplier: { name: "Provide One S.L." },
        customer: { name: "Sample Consumer S.A.", taxId: "B98602642" },
        uuid: "0190e6e5-6c5a-7a1e-9f1c-2d3c4b5a6f70",
        total: "€1,815.00",
      },
      quota: { used: 342, limit: 500 },
    },
    template: [
      {
        type: "card",
        children: [
          { type: "card_header", bind: { title: "invoice.code", subtitle: "invoice.issued" } },
          {
            type: "card_content",
            props: { class: "flex flex-col gap-4" },
            children: [
              {
                type: "button_group",
                children: [
                  { type: "status_badge", bind: { label: "invoice.status.label", status: "invoice.status.state" } },
                  { type: "tag_status", props: { status: "blue", dot: true }, bind: { label: "invoice.country" } },
                  { type: "tag_status", props: { status: "purple" }, bind: { label: "invoice.regime" } },
                ],
              },
              {
                type: "description_list",
                children: [
                  { type: "description_list_item", props: { label: "Supplier", inline: true }, bind: { value: "invoice.supplier.name" } },
                  { type: "description_list_item", props: { label: "Customer", inline: true }, bind: { value: "invoice.customer.name" } },
                  { type: "description_list_item", props: { label: "Tax ID", inline: true, mono: true }, bind: { value: "invoice.customer.taxId" } },
                  {
                    type: "description_list_item",
                    props: { label: "Document UUID", inline: true, mono: true, prefixLength: 8, suffixLength: 4 },
                    bind: { value: "invoice.uuid" },
                  },
                  { type: "description_list_item", props: { label: "Total", inline: true }, bind: { value: "invoice.total" } },
                ],
              },
            ],
          },
        ],
      },
      {
        type: "card",
        props: { class: "mt-4" },
        children: [
          { type: "card_header", props: { title: "Monthly quota" } },
          {
            type: "card_content",
            children: [
              {
                type: "card_progress_bar",
                props: { title: "Invoices", subtitle: "This month" },
                bind: { current: "quota.used", total: "quota.limit" },
              },
            ],
          },
        ],
      },
      {
        type: "button_group",
        props: { class: "mt-4", align: "right" },
        children: [
          { type: "button", props: { variant: "secondary" }, children: [text("Download PDF")] },
          { type: "button", props: { variant: "primary" }, children: [text("Send again")] },
        ],
      },
    ],
  },

  // Interactive: form controls are bound to "data" with "model", and events
  // run named actions with "on". The JSON holds no code: "alert" and "reset"
  // are built into jsonui.js, "save" is registered by this page's own JS.
  interactive: {
    $schema: SCHEMA,
    data: {
      customer: { name: "Acme S.L.", email: "billing@acme.com", country: "es", einvoice: true, notes: "" },
    },
    template: [
      {
        type: "form",
        props: { class: "flex flex-col gap-4" },
        on: {
          "submit.prevent": {
            action: "alert",
            message: "Saving {{customer.name}} <{{customer.email}}>\nCountry: {{customer.country}}\nE-invoices: {{customer.einvoice}}\nNotes: {{customer.notes}}",
          },
        },
        children: [
          { type: "title_group", props: { title: "Customer", description: "Edit the fields, then save." } },
          {
            type: "fieldset",
            props: { legend: "Identity", variant: "card" },
            children: [
              { type: "input", props: { name: "name", label: "Legal name", required: true }, model: "customer.name" },
              { type: "input", props: { name: "email", type: "email", label: "Email" }, model: "customer.email" },
              {
                type: "select",
                props: {
                  name: "country",
                  label: "Country",
                  options: [
                    { value: "es", label: "Spain" },
                    { value: "pt", label: "Portugal" },
                    { value: "fr", label: "France" },
                    { value: "de", label: "Germany" },
                  ],
                },
                model: "customer.country",
              },
            ],
          },
          {
            type: "fieldset",
            props: { legend: "Preferences", variant: "card" },
            children: [
              {
                type: "checkbox",
                props: { name: "einvoice", variant: "switch", label: "Send e-invoices", description: "Deliver through the national network." },
                model: "customer.einvoice",
              },
              { type: "textarea", props: { name: "notes", label: "Notes", rows: 3 }, model: "customer.notes" },
            ],
          },
          {
            type: "p",
            props: { class: "text-sm text-foreground-default-secondary" },
            children: [text("Preview: "), { type: "text", bind: "customer.name" }, text(" · "), { type: "text", bind: "customer.email" }],
          },
          {
            type: "button_group",
            props: { align: "right" },
            children: [
              { type: "button", props: { variant: "secondary", type: "button" }, on: { click: { action: "reset" } }, children: [text("Reset")] },
              { type: "button", props: { variant: "secondary", type: "button" }, on: { click: { action: "save" } }, children: [text("Send to host")] },
              { type: "button", props: { type: "submit" }, children: [text("Save")] },
            ],
          },
        ],
      },
    ],
  },

  form: {
    $schema: SCHEMA,
    template: [
      {
        type: "form",
        props: { id: "party-form", action: "#", class: "flex flex-col gap-4" },
        children: [
          {
            type: "title_group",
            props: { title: "Customer details", description: "Used on every invoice issued to this party." },
          },
          {
            type: "fieldset",
            props: { legend: "Identity", variant: "card" },
            children: [
              { type: "input", props: { name: "name", label: "Legal name", placeholder: "Acme S.L.", required: true } },
              { type: "input", props: { name: "tax_id", label: "Tax ID", value: "B98602642" } },
              {
                type: "select",
                props: {
                  name: "country",
                  label: "Country",
                  options: [
                    { value: "es", label: "Spain", selected: true },
                    { value: "pt", label: "Portugal" },
                    { value: "fr", label: "France" },
                    { value: "de", label: "Germany" },
                  ],
                },
              },
            ],
          },
          {
            type: "fieldset",
            props: { legend: "Contact", variant: "card" },
            children: [
              { type: "input", props: { name: "email", type: "email", label: "Email", placeholder: "billing@acme.com" } },
              { type: "textarea", props: { name: "notes", label: "Notes", rows: 3 } },
              {
                type: "checkbox",
                props: { name: "einvoice", variant: "switch", label: "Send e-invoices", description: "Deliver through the national network.", checked: true },
              },
            ],
          },
          {
            type: "button_group",
            props: { align: "right" },
            children: [
              { type: "button", props: { variant: "secondary", type: "reset" }, children: [text("Reset")] },
              { type: "button", props: { type: "submit" }, children: [text("Save")] },
            ],
          },
        ],
      },
    ],
  },
};

// blocked shows what strict mode rejects: Alpine directives in attributes,
// javascript: URLs and actions the host hasn't allowed. Untick "strict" to
// render it as a trusted document.
window.SAMPLES.blocked = {
  $schema: SCHEMA,
  data: { name: "Acme" },
  template: [
    {
      type: "card",
      children: [
        { type: "card_header", props: { title: "Untrusted document" } },
        {
          type: "card_content",
          props: { class: "flex flex-col gap-2" },
          children: [
            {
              type: "button",
              props: { attributes: { "@click": "alert('Script ran from document attributes: ' + document.title)" } },
              children: [text("Run inline script")],
            },
          ],
        },
      ],
    },
  ],
};

// stress builds n invoice cards, each with a header, badge and five rows.
window.SAMPLES.stress = (n = 200) => ({
  $schema: SCHEMA,
  template: Array.from({ length: n }, (_, i) => ({
    type: "card",
    props: { class: "mb-2" },
    children: [
      { type: "card_header", props: { title: `Invoice INV-${String(i + 1).padStart(5, "0")}`, subtitle: "Issued 27 Oct 2026" } },
      {
        type: "card_content",
        children: [
          { type: "status_badge", props: { label: i % 3 ? "Sent" : "Failed", status: i % 3 ? "success" : "failed" } },
          {
            type: "description_list",
            children: [
              { type: "description_list_item", props: { label: "Customer", value: `Customer ${i}`, inline: true } },
              { type: "description_list_item", props: { label: "Tax ID", value: `B${10000000 + i}`, inline: true, mono: true } },
              { type: "description_list_item", props: { label: "Currency", value: "EUR", inline: true } },
              { type: "description_list_item", props: { label: "Lines", value: String((i % 9) + 1), inline: true } },
              { type: "description_list_item", props: { label: "Total", value: `€${(i * 13.37).toFixed(2)}`, inline: true } },
            ],
          },
        ],
      },
    ],
  })),
});
