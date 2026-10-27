// editor.js: the playground's CodeMirror editor, set up like gobl.dev's.
//
// The JSON Schema comes from the document's own "$schema", resolved against
// the page URL (so the samples' relative "schema.json" loads the schema
// published next to the playground) and reloaded whenever it changes. Every
// step is logged to the console under [playground] to help debugging.
//
// Once mounted the editor is handed to the page through window.attachEditor;
// until then, or if anything here fails, the page keeps using its textarea.

import { basicSetup, EditorView } from "codemirror";
import { EditorState, Compartment } from "@codemirror/state";
import { jsonSchema, updateSchema } from "codemirror-json-schema";
import { vsCodeLight } from "@fsegurai/codemirror-theme-vscode-light";
import { vsCodeDark } from "@fsegurai/codemirror-theme-vscode-dark";

const log = (...args) => console.info("[playground]", ...args);

// SCHEMA_LIMIT is the document size above which schema checks are switched
// off. Validation in json-schema-library (10.x and 11.x) copies the whole
// evaluation path for every "if", so it slows down with the square of the
// document size: the 200 card stress sample takes ~30s and 1000 cards
// crashes the tab. Real documents are far below this.
const SCHEMA_LIMIT = 64 * 1024;

let schemaURL = null; // resolved from the document's "$schema"
let schema; // fetched from schemaURL
let applied; // the schema the editor is using, if any

// loadSchema reads "$schema" from the editor's text and fetches that schema,
// unless it is already loaded.
async function loadSchema(view) {
  const m = view.state.doc.toString().match(/"\$schema"\s*:\s*"([^"]+)"/);
  if (!m) return;
  const url = new URL(m[1], location.href).href;
  if (url === schemaURL) return;
  schemaURL = url;
  try {
    const res = await fetch(url);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    schema = await res.json();
    log(`schema loaded from ${url} (${Object.keys(schema.$defs || {}).length} definitions)`);
    syncSchema(view);
  } catch (err) {
    schemaURL = null; // retry on the next change
    console.warn("[playground] schema loading failed:", url, err);
  }
}

// syncSchema applies the loaded schema, or removes it while the document is
// larger than SCHEMA_LIMIT. It must not run during an editor update.
function syncSchema(view) {
  const size = view.state.doc.length;
  const large = size > SCHEMA_LIMIT;
  const note = document.getElementById("editor-note");
  if (note) {
    note.hidden = !large;
    note.textContent = `Schema checks are off for documents over ${SCHEMA_LIMIT / 1024}KB (this one is ${Math.round(size / 1024)}KB).`;
  }
  const want = large ? undefined : schema;
  if (want === applied) return;
  applied = want;
  updateSchema(view, want);
  if (large) log(`schema checks off: document is ${Math.round(size / 1024)}KB`);
}

const dark = matchMedia("(prefers-color-scheme: dark)");
const theme = new Compartment();
const listeners = [];
let timer;

const view = new EditorView({
  state: EditorState.create({
    doc: "",
    extensions: [
      basicSetup,
      theme.of(dark.matches ? vsCodeDark : vsCodeLight),
      EditorView.lineWrapping,
      jsonSchema(),
      EditorView.updateListener.of((update) => {
        if (!update.docChanged) return;
        // Check the size straight away, before the linter runs.
        queueMicrotask(() => syncSchema(update.view));
        clearTimeout(timer);
        timer = setTimeout(() => {
          loadSchema(update.view);
          listeners.forEach((fn) => fn());
        }, 300);
      }),
    ],
  }),
  parent: document.getElementById("code"),
});

dark.addEventListener("change", (e) => {
  view.dispatch({ effects: theme.reconfigure(e.matches ? vsCodeDark : vsCodeLight) });
});

window.playgroundEditor = view;
window.attachEditor({
  getValue: () => view.state.doc.toString(),
  setValue(text) {
    view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } });
    syncSchema(view);
    loadSchema(view);
  },
  onChange: (fn) => listeners.push(fn),
});
log("editor ready");
