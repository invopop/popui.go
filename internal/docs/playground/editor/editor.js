// editor.js: the playground's CodeMirror editor, set up like gobl.dev's.
//
// The JSON Schema comes from the panel's own "$schema", resolved against
// the page URL and reloaded whenever it changes. When it's the "$id" of the
// schema published next to the playground (panel.json), that local copy is
// used instead, so panels can use the canonical address while the editor
// checks them against the schema that matches this build. Every step is
// logged to the console under [playground] to help debugging.
//
// Once mounted the editor is handed to the page through window.attachEditor.
//
// It's bundled with its dependencies into dist/editor.js by `npm run build`,
// which `popui serve` and `popui build` read; see package.json.

import { basicSetup, EditorView } from "codemirror";
import { EditorState, Compartment } from "@codemirror/state";
import { jsonSchema, updateSchema } from "@invopop/codemirror-json-schema";
import { vsCodeLight } from "@fsegurai/codemirror-theme-vscode-light";
import { vsCodeDark } from "@fsegurai/codemirror-theme-vscode-dark";

const log = (...args) => console.info("[playground]", ...args);

// SCHEMA_LIMIT is the panel size above which schema checks are switched
// off. Validation takes time in proportion to the panel, about a second
// at this size (the 200 card stress sample, ~380KB, takes ~1.4s), and runs
// again after each edit. Real panels are far below this.
const SCHEMA_LIMIT = 256 * 1024;

let schemaURL = null; // resolved from the panel's "$schema"
let schema; // fetched from schemaURL
let applied; // the schema the editor is using, if any

// localSchema is the schema published next to the playground, or null if it
// can't be loaded, in which case "$schema" is always fetched as given.
const localSchema = fetch("panel.json")
  .then((res) => {
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return res.json();
  })
  .then((s) => {
    log(`local panel.json loaded for ${s.$id}`);
    return s;
  })
  .catch((err) => {
    console.warn("[playground] local panel.json loading failed:", err);
    return null;
  });

// loadSchema reads "$schema" from the editor's text and fetches that schema,
// unless it is already loaded.
async function loadSchema(view) {
  const m = view.state.doc.toString().match(/"\$schema"\s*:\s*"([^"]+)"/);
  if (!m) return;
  const url = new URL(m[1], location.href).href;
  if (url === schemaURL) return;
  schemaURL = url;
  try {
    const local = await localSchema;
    if (local && local.$id === url) {
      schema = local;
      log(`schema ${url} replaced by the local panel.json`);
    } else {
      const res = await fetch(url);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      schema = await res.json();
      log(`schema loaded from ${url} (${Object.keys(schema.$defs || {}).length} definitions)`);
    }
    syncSchema(view);
  } catch (err) {
    schemaURL = null; // retry on the next change
    console.warn("[playground] schema loading failed:", url, err);
  }
}

// syncSchema applies the loaded schema, or removes it while the panel is
// larger than SCHEMA_LIMIT. It must not run during an editor update.
function syncSchema(view) {
  const size = view.state.doc.length;
  const large = size > SCHEMA_LIMIT;
  const note = document.getElementById("editor-note");
  if (note) {
    note.hidden = !large;
    note.textContent = `Schema checks are off for panels over ${SCHEMA_LIMIT / 1024}KB (this one is ${Math.round(size / 1024)}KB).`;
  }
  const want = large ? undefined : schema;
  if (want === applied) return;
  applied = want;
  updateSchema(view, want);
  if (large) log(`schema checks off: panel is ${Math.round(size / 1024)}KB`);
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
