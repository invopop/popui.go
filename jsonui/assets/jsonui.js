// jsonui.js: runs the actions declared in jsonui documents.
//
// A node's "on" handlers render as `x-on:<event>="$actions([...], $event)"`, so the
// JSON only names actions and their params; the behaviour lives here and in
// whatever the host page registers:
//
//   popui.actions.register("save", ({ params, data, event, el }) => { ... });
//
// String params may reference the document data with {{path}} placeholders,
// e.g. "Hello {{customer.name}}", resolved when the action runs. Two params
// work with every action:
//
//   into: data path to store the action's result in (skipped when the
//         action returns undefined, e.g. a cancelled file picker);
//   busy: data path set to true while the action runs.
//
// Form controls with a "model" show the error for their path, kept by
// document in $errors. Errors come from validating the data against the
// document's schemas with popuiValidate, from popui.wasm, as it changes, and
// from actions that throw an error with an "errors" list, such as a host
// whose server rejected the data:
//
//   throw Object.assign(new Error("invalid config"), {
//     errors: [{ path: "config.to", message: "Not a verified address" }],
//   });
//
// The action list stops there, and each error shows until its field is
// changed or the action runs again. Errors from schemas only show once their field has been changed
// or left, or after the "validate" action, which stops the list while there
// are any.
(function () {
  "use strict";

  const popui = (window.popui = window.popui || {});
  const registry = new Map();

  // Paths only follow the data's own properties: through inherited ones
  // such as "toString" or "__proto__" they would reach objects shared with
  // the host page, like Object.prototype.
  const own = (v, k) => v != null && typeof v === "object" && Object.hasOwn(v, k);

  const get = (obj, path) => path.split(".").reduce((v, k) => (own(v, k) ? v[k] : undefined), obj);

  const unsafeKeys = new Set(["__proto__", "prototype", "constructor"]);

  const set = (obj, path, value) => {
    const keys = String(path).split(".");
    if (keys.some((k) => unsafeKeys.has(k))) throw new Error(`[jsonui] unsafe data path "${path}"`);
    const last = keys.pop();
    const target = keys.reduce((v, k) => {
      if (!own(v, k)) throw new Error(`[jsonui] data has no "${path}"`);
      return v[k];
    }, obj);
    if (target == null || typeof target !== "object") throw new Error(`[jsonui] data has no "${path}"`);
    target[last] = value;
  };

  // busy flags stay true while any action using them runs, so overlapping
  // actions don't clear each other's flag. Counts are kept per data object
  // and path.
  const running = new WeakMap();
  const busy = (data, path, delta) => {
    if (!running.has(data)) running.set(data, new Map());
    const counts = running.get(data);
    const n = (counts.get(path) ?? 0) + delta;
    if (n > 0) counts.set(path, n);
    else counts.delete(path);
    set(data, path, n > 0);
  };

  // interpolate replaces {{path}} placeholders in strings, recursively.
  const interpolate = (value, data) => {
    if (typeof value === "string") {
      return value.replace(/\{\{\s*([\w$.]+)\s*\}\}/g, (_, p) => {
        const v = get(data, p);
        return v == null ? "" : typeof v === "object" ? JSON.stringify(v) : String(v);
      });
    }
    if (Array.isArray(value)) return value.map((v) => interpolate(v, data));
    if (value && typeof value === "object") {
      return Object.fromEntries(Object.entries(value).map(([k, v]) => [k, interpolate(v, data)]));
    }
    return value;
  };

  // snapshot returns a plain copy of the reactive data.
  const snapshot = (data) => JSON.parse(JSON.stringify(data));

  // STOP is returned by an action to end the list it runs in.
  const STOP = Symbol("stop");

  popui.actions = {
    STOP,
    register(name, fn) {
      registry.set(name, fn);
    },
    has(name) {
      return registry.has(name);
    },
    // run executes a list of actions in order, awaiting async ones.
    async run(list, ctx) {
      for (const { action, into, busy: busyPath, ...params } of list ?? []) {
        const fn = registry.get(action);
        if (!fn) {
          console.error(`[jsonui] unknown action "${action}"`);
          return;
        }
        if (busyPath) busy(ctx.data, busyPath, 1);
        retried(ctx.el, action);
        try {
          const result = await fn({ ...ctx, params: interpolate(params, ctx.data), snapshot: () => snapshot(ctx.data) });
          if (result === STOP) return;
          if (into && result !== undefined) set(ctx.data, into, result);
        } catch (err) {
          if (!Array.isArray(err?.errors)) throw err;
          reported(ctx.el, action, err.errors);
          return;
        } finally {
          if (busyPath) busy(ctx.data, busyPath, -1);
        }
      }
    },
  };

  // Built-in actions.
  popui.actions.register("alert", ({ params }) => window.alert(params.message ?? ""));
  popui.actions.register("set", ({ params, data }) => set(data, params.path, params.value));
  popui.actions.register("reset", ({ el, data }) => {
    // The document root's x-data is the original JSON object. Keys added
    // since are removed. Documents without data have no root and nothing
    // to reset.
    const root = el.closest("[data-popui-document]");
    if (!root) return;
    const initial = JSON.parse(root.getAttribute("x-data"));
    for (const key of Object.keys(data)) if (!(key in initial)) delete data[key];
    Object.assign(data, initial);
    Object.assign(errorState(root), { schema: {}, host: {}, reporter: undefined, touched: {}, revealed: false, checked: !root.getAttribute("x-jsonui-errors"), failure: "" });
  });
  popui.actions.register("dispatch", ({ params, el, snapshot }) => {
    el.dispatchEvent(new CustomEvent(params.event, { bubbles: true, detail: params.detail ?? snapshot() }));
  });

  popui.actions.register("validate", ({ params, el, data }) => {
    const root = el.closest("[data-popui-document]");
    if (!root) return;
    const s = errorState(root);
    validate(root, s, data);
    if (!params.silent) s.revealed = true;
    if (!s.checked || Object.keys(s.schema).length > 0) return STOP;
  });

  // Errors are kept for each document root: those found by validating the
  // data against its schemas, those reported by actions, and which fields
  // have been touched, by path.
  const errors = new WeakMap();
  const errorState = (root) => {
    let s = errors.get(root);
    if (!s) {
      s = window.Alpine.reactive({ schema: {}, host: {}, reporter: undefined, touched: {}, revealed: false, checked: !root.getAttribute("x-jsonui-errors"), failure: "" });
      errors.set(root, s);
    }
    return s;
  };

  // reported shows the errors an action reported, replacing any reported
  // before.
  const reported = (el, action, list) => {
    const root = el.closest("[data-popui-document]");
    if (!root) return console.error("[jsonui] errors reported outside a document", list);
    const host = {};
    for (const { path, message } of list) if (path && !(path in host)) host[path] = String(message ?? "");
    Object.assign(errorState(root), { host, reporter: action });
  };

  // retried clears the errors an action reported when it runs again.
  const retried = (el, action) => {
    const root = el.closest("[data-popui-document]");
    const s = root && errors.get(root);
    if (s?.reporter === action) Object.assign(s, { host: {}, reporter: undefined });
  };

  let unavailable = false;
  // validate updates the errors found in the data by the document's schemas.
  // It needs popuiValidate from popui.wasm.
  const validate = (root, s, data) => {
    const schemas = root.getAttribute("x-jsonui-errors");
    if (!schemas) {
      Object.assign(s, { schema: {}, checked: true, failure: "" });
      return;
    }
    s.checked = false;
    if (typeof window.popuiValidate !== "function") {
      if (!unavailable) console.warn("[jsonui] popui.wasm isn't loaded, so data isn't validated");
      unavailable = true;
      s.failure = "Unable to validate this form. Please try again.";
      return;
    }
    try {
      const res = window.popuiValidate(schemas, JSON.stringify(data));
      if (res.error) throw new Error(res.error);
      const found = {};
      for (const { path, message } of res.errors ?? []) found[path] = message;
      Object.assign(s, { schema: found, checked: true, failure: "" });
    } catch (err) {
      console.error("[jsonui] validation failed", err);
      s.failure = "Unable to validate this form. Please try again.";
    }
  };

  // The $errors API for a document. Errors reported by actions show as soon
  // as they're reported, those from schemas once their path is touched or
  // the document validated.
  const errorsAPI = (root, s) => ({
    shown(path) {
      if (path in s.host) return s.host[path];
      return (s.revealed || s.touched[path]) && path in s.schema ? s.schema[path] : "";
    },
    get checked() {
      return s.checked;
    },
    get valid() {
      return s.checked && Object.keys(s.schema).length === 0;
    },
    // unplaced lists the shown errors no field is showing.
    unplaced() {
      const fields = new Set([...root.querySelectorAll("[data-jsonui-field]")].map((f) => f.dataset.jsonuiField));
      const paths = new Set([...Object.keys(s.host), ...Object.keys(s.schema)]);
      const list = [...paths]
        .filter((p) => !fields.has(p))
        .map((path) => ({ path, message: this.shown(path) }))
        .filter((e) => e.message);
      if (s.revealed && s.failure) list.push({ path: "", message: s.failure });
      return list;
    },
  });
  const noErrors = { shown: () => "", checked: true, valid: true, unplaced: () => [] };

  // field finds the path of the field an event came from.
  const field = (e) => e.target?.closest?.("[data-jsonui-field]")?.dataset.jsonuiField;

  // Actions work on the document's own data, the x-data of its root, rather
  // than the merged scope at the element: components such as tabs add scopes
  // inside the document, and the host page may have scopes around it.
  // Documents without data have no root and use the element's scope.
  const install = (Alpine) => {
    Alpine.magic("actions", (el) => (list, event) => {
      const root = el.closest("[data-popui-document]");
      const data = root ? Alpine.closestDataStack(root)[0] : Alpine.$data(el);
      return popui.actions.run(list, { el, event, data });
    });

    Alpine.magic("errors", (el) => {
      const root = el.closest("[data-popui-document]");
      return root ? errorsAPI(root, errorState(root)) : noErrors;
    });

    // x-jsonui-errors is rendered on the root of documents with fields, with
    // the document's schemas, if any. The data is validated shortly after it
    // changes, a field is touched when it's changed or left, and an error an
    // action reported clears when its field is changed.
    Alpine.directive("jsonui-errors", (el, _, { effect, cleanup }) => {
      const s = errorState(el);
      const data = Alpine.closestDataStack(el)[0];
      let timer;
      effect(() => {
        JSON.stringify(data); // depend on all of it
        if (el.getAttribute("x-jsonui-errors")) s.checked = false;
        clearTimeout(timer);
        timer = setTimeout(() => validate(el, s, data), 100);
      });
      const touch = (e) => {
        const path = field(e);
        if (path) s.touched[path] = true;
      };
      const edit = (e) => {
        const path = field(e);
        if (path && path in s.host) delete s.host[path];
      };
      el.addEventListener("focusout", touch);
      el.addEventListener("change", touch);
      el.addEventListener("input", edit);
      el.addEventListener("change", edit);
      cleanup(() => {
        clearTimeout(timer);
        el.removeEventListener("focusout", touch);
        el.removeEventListener("change", touch);
        el.removeEventListener("input", edit);
        el.removeEventListener("change", edit);
      });
    });
  };

  if (window.Alpine) install(window.Alpine);
  else document.addEventListener("alpine:init", () => install(window.Alpine));
})();
