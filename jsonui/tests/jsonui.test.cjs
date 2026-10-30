const assert = require("node:assert/strict");
const { readFileSync } = require("node:fs");
const { test } = require("node:test");
const { runInNewContext } = require("node:vm");

// Exercise the public action and magic APIs with a minimal Alpine adapter.
// No browser or third-party dependencies are needed.
function setup(schemas = '{"config":{"type":"string"}}') {
  const magic = {};
  const directives = {};
  const timers = new Map();
  let nextTimer = 0;
  const data = { config: "valid" };
  const root = {
    getAttribute: (name) => name === "x-jsonui-errors" ? schemas : JSON.stringify({ config: "valid" }),
    closest() { return this; },
    querySelectorAll: () => [],
    addEventListener() {},
    removeEventListener() {},
  };
  const window = {
    Alpine: {
      reactive: (s) => s,
      magic: (name, fn) => { magic[name] = fn; },
      directive: (name, fn) => { directives[name] = fn; },
      closestDataStack: () => [data],
    },
  };
  runInNewContext(readFileSync(`${__dirname}/../assets/jsonui.js`, "utf8"), {
    window,
    console: { warn() {}, error() {} },
    setTimeout(fn) { const id = ++nextTimer; timers.set(id, fn); return id; },
    clearTimeout(id) { timers.delete(id); },
  });
  let saves = 0;
  window.popui.actions.register("save", () => { saves++; });
  return {
    window, data, root,
    errors: magic.errors(root),
    get saves() { return saves; },
    run: (list = [{ action: "validate" }, { action: "save" }]) => window.popui.actions.run(list, { el: root, data }),
    track() {
      let effect;
      directives["jsonui-errors"](root, {}, { effect(fn) { effect = fn; fn(); }, cleanup() {} });
      return () => effect();
    },
    flush() { for (const [id, fn] of timers) { timers.delete(id); fn(); } },
  };
}

test("missing validator stops Save and shows a panel error; retry recovers", async () => {
  const form = setup();
  assert.equal(form.errors.checked, false);
  assert.equal(form.errors.valid, false);
  await form.run();
  assert.equal(form.saves, 0);
  assert.match(form.errors.unplaced()[0].message, /Unable to validate/);
  form.window.popuiValidate = () => ({ errors: [] });
  await form.run();
  assert.equal(form.saves, 1);
  assert.equal(form.errors.checked, true);
  assert.equal(form.errors.valid, true);
  assert.equal(form.errors.unplaced().length, 0);
});

for (const [name, failure] of [["error", () => ({ error: "compile failure" })], ["exception", () => { throw new Error("runtime failure"); }]]) {
  test(`validator ${name} stops Save after a previously valid check`, async () => {
    const form = setup();
    form.window.popuiValidate = () => ({ errors: [] });
    await form.run();
    form.window.popuiValidate = failure;
    await form.run();
    assert.equal(form.saves, 1);
    assert.equal(form.errors.checked, false);
    assert.equal(form.errors.valid, false);
    assert.match(form.errors.unplaced()[0].message, /Unable to validate/);
  });
}

test("silent validation blocks without revealing panel errors", async () => {
  const form = setup();
  await form.run([{ action: "validate", silent: true }, { action: "save" }]);
  assert.equal(form.saves, 0);
  assert.equal(form.errors.unplaced().length, 0);
});

test("schema errors block Save and keep friendly/custom messages", async () => {
  const form = setup();
  form.window.popuiValidate = () => ({ errors: [{ path: "config", message: "Enter a recipient" }] });
  await form.run();
  assert.equal(form.saves, 0);
  assert.equal(form.errors.checked, true);
  assert.equal(form.errors.valid, false);
  assert.equal(form.errors.shown("config"), "Enter a recipient");
});

test("data changes invalidate the previous result during debounce", async () => {
  const form = setup();
  form.window.popuiValidate = () => ({ errors: [] });
  const changed = form.track();
  form.flush();
  assert.equal(form.errors.valid, true);
  form.data.config = 42;
  changed();
  assert.equal(form.errors.checked, false);
  assert.equal(form.errors.valid, false);
  form.window.popuiValidate = () => ({ errors: [{ path: "config", message: "Must be a string" }] });
  form.flush();
  assert.equal(form.errors.checked, true);
  assert.equal(form.errors.valid, false);
});

test("reset clears validation failure and requires a fresh check", async () => {
  const form = setup();
  await form.run();
  await form.run([{ action: "reset" }]);
  assert.equal(form.errors.unplaced().length, 0);
  assert.equal(form.errors.checked, false);
  assert.equal(form.errors.valid, false);
});

test("background validation failure remains hidden until explicit validation", async () => {
  const form = setup();
  form.track();
  form.flush();
  assert.equal(form.errors.unplaced().length, 0);
  assert.equal(form.errors.valid, false);
  await form.run();
  assert.equal(form.errors.unplaced().length, 1);
});

test("panels without schemas can save without WASM", async () => {
  const form = setup("");
  assert.equal(form.errors.checked, true);
  assert.equal(form.errors.valid, true);
  await form.run();
  assert.equal(form.saves, 1);
});
