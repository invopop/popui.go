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
(function () {
  "use strict";

  const popui = (window.popui = window.popui || {});
  const registry = new Map();

  const get = (obj, path) => path.split(".").reduce((v, k) => (v == null ? undefined : v[k]), obj);

  // Paths through these keys would reach Object.prototype and pollute the
  // host page, so set refuses them.
  const unsafeKeys = new Set(["__proto__", "prototype", "constructor"]);

  const set = (obj, path, value) => {
    const keys = String(path).split(".");
    if (keys.some((k) => unsafeKeys.has(k))) throw new Error(`[jsonui] unsafe data path "${path}"`);
    const last = keys.pop();
    const target = keys.reduce((v, k) => v[k], obj);
    target[last] = value;
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

  popui.actions = {
    register(name, fn) {
      registry.set(name, fn);
    },
    has(name) {
      return registry.has(name);
    },
    // run executes a list of actions in order, awaiting async ones.
    async run(list, ctx) {
      for (const { action, into, busy, ...params } of list) {
        const fn = registry.get(action);
        if (!fn) {
          console.error(`[jsonui] unknown action "${action}"`);
          return;
        }
        if (busy) set(ctx.data, busy, true);
        try {
          const result = await fn({ ...ctx, params: interpolate(params, ctx.data), snapshot: () => snapshot(ctx.data) });
          if (into && result !== undefined) set(ctx.data, into, result);
        } finally {
          if (busy) set(ctx.data, busy, false);
        }
      }
    },
  };

  // Built-in actions.
  popui.actions.register("alert", ({ params }) => window.alert(params.message ?? ""));
  popui.actions.register("set", ({ params, data }) => set(data, params.path, params.value));
  popui.actions.register("reset", ({ el, data }) => {
    // The document root's x-data is the original JSON object. Keys added
    // since are removed.
    const initial = JSON.parse(el.closest("[data-popui-document]").getAttribute("x-data"));
    for (const key of Object.keys(data)) if (!(key in initial)) delete data[key];
    Object.assign(data, initial);
  });
  popui.actions.register("dispatch", ({ params, el, snapshot }) => {
    el.dispatchEvent(new CustomEvent(params.event, { bubbles: true, detail: params.detail ?? snapshot() }));
  });

  // Actions work on the document's own data, the x-data of its root, rather
  // than the merged scope at the element: components such as tabs add scopes
  // inside the document, and the host page may have scopes around it.
  // Documents without data have no root and use the element's scope.
  const install = (Alpine) =>
    Alpine.magic("actions", (el) => (list, event) => {
      const root = el.closest("[data-popui-document]");
      const data = root ? Alpine.closestDataStack(root)[0] : Alpine.$data(el);
      return popui.actions.run(list, { el, event, data });
    });

  if (window.Alpine) install(window.Alpine);
  else document.addEventListener("alpine:init", () => install(window.Alpine));
})();
