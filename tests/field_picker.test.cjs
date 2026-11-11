const assert = require("node:assert/strict");
const { readFileSync } = require("node:fs");
const { test } = require("node:test");
const { runInNewContext } = require("node:vm");

// Loads popui.js against a minimal browser and Alpine adapter, and returns
// a field picker over the given fields, initialised.
function fieldPicker(init) {
  const listeners = {};
  const data = {};
  const Alpine = { data(name, fn) { data[name] = fn; }, directive() {}, magic() {}, store() {}, bind() {} };
  runInNewContext(readFileSync(`${__dirname}/../assets/popui.js`, "utf8"), {
    window: { Alpine, addEventListener() {} },
    document: {
      addEventListener(type, fn) { (listeners[type] ??= []).push(fn); },
      querySelectorAll: () => [],
      getElementById: () => null,
    },
    Alpine,
    CSS: { supports: () => true },
    MutationObserver: class { observe() {} },
    navigator: {},
    console,
    setTimeout,
    clearTimeout,
  });
  listeners["alpine:init"].forEach((fn) => fn());
  const picker = data.fieldPicker(init);
  picker.$watch = () => {};
  picker.init();
  return picker;
}

// Invoice (group) > lines (array) > meta (group) > key
const fields = [
  { n: "Invoice", f: 8 },
  { n: "lines", p: 0, f: 1 | 4 },
  { n: "meta", p: 1, f: 8 },
  { n: "key", p: 2 },
];

test("scope names a group under an array as the panel shows it", () => {
  const picker = fieldPicker({ fields, scope: "Invoice.lines[].meta" });
  assert.equal(picker.scoped, 2);
  assert.equal(picker.rootLabel, "meta");
  assert.equal(picker.entries.map((e) => e.path).join(","), "key");
});

test("scope without the array marker names no group", () => {
  const picker = fieldPicker({ fields, scope: "Invoice.lines.meta" });
  assert.equal(picker.scoped, -1);
  assert.equal(picker.entries.length, 4);
});
