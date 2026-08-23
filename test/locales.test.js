const fs = require("fs");
const path = require("path");

const root = path.join(__dirname, "..", "client");
const main = fs.readFileSync(path.join(root, "i18n.go"), "utf8");
const extra = fs.readFileSync(path.join(root, "i18n_extra.go"), "utf8");
const settings = fs.readFileSync(path.join(root, "settings.go"), "utf8");

// How many page keys each locale is still missing today. The numbers may only
// go down: a locale that falls further behind fails the check.
const SETTINGS_DEBT = { ru: 0, uk: 58, de: 69, fr: 69, es: 69, it: 69, pl: 69 };

function keysIn(text) {
  return new Set([...text.matchAll(/"([A-Za-z_][A-Za-z_0-9.]*)":/g)].map((m) => m[1]));
}

function blockOf(text, header, close) {
  const start = text.indexOf(header);
  if (start < 0) return null;
  const end = text.indexOf(close, start + header.length);
  if (end < 0) return null;
  return text.slice(start + header.length, end);
}

function sectionOf(text, header) {
  const start = text.indexOf(header);
  if (start < 0) return "";
  const next = text.indexOf("\nvar ", start + header.length);
  return next < 0 ? text.slice(start) : text.slice(start, next);
}

const settingsSection = sectionOf(main, "var settingsStrings = map[string]map[string]string{");
const msgsSection = sectionOf(main, "var msgs = map[string]map[string]string{");

function settingsBlock(locale) {
  const inline = blockOf(settingsSection, `\t"${locale}": {`, "\n\t},");
  if (inline) return inline;
  return blockOf(extra, `settingsStrings["${locale}"] = map[string]string{`, "\n\t}");
}

function msgsBlock(locale) {
  const inline = blockOf(msgsSection, `\t"${locale}": {`, "\n\t},");
  if (inline) return inline;
  return blockOf(extra, `msgs["${locale}"] = map[string]string{`, "\n\t}");
}

const failures = [];
function check(name, got, want) {
  const ok = JSON.stringify(got) === JSON.stringify(want);
  console.log(`${ok ? "PASS" : "FAIL"}  ${name}${ok ? "" : ` (got ${JSON.stringify(got)}, want ${JSON.stringify(want)})`}`);
  if (!ok) failures.push(name);
}

const en = keysIn(settingsBlock("en"));
const enMsgs = keysIn(msgsBlock("en"));
const page = new Set([...settings.matchAll(/\{\{(S_[A-Z_0-9]+)\}\}/g)].map((m) => m[1]));

check("every key the page asks for exists in English", [...page].filter((k) => !en.has(k)), []);
check("English carries the whole page", page.size > 100, true);

for (const locale of Object.keys(SETTINGS_DEBT)) {
  const keys = keysIn(settingsBlock(locale));
  const unknown = [...keys].filter((k) => !en.has(k));
  check(`${locale} invents no settings keys`, unknown, []);

  const missing = [...page].filter((k) => !keys.has(k));
  const ok = missing.length <= SETTINGS_DEBT[locale];
  console.log(
    `${ok ? "PASS" : "FAIL"}  ${locale} covers the page: ${page.size - missing.length}/${page.size}` +
      (ok ? "" : ` — worse than the recorded ${page.size - SETTINGS_DEBT[locale]}/${page.size}`)
  );
  if (!ok) failures.push(`${locale} settings coverage`);

  const msgs = keysIn(msgsBlock(locale));
  const unknownMsgs = [...msgs].filter((k) => !enMsgs.has(k));
  check(`${locale} invents no message keys`, unknownMsgs, []);
}

if (failures.length) {
  console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
  process.exit(1);
}
console.log("\nall locale checks passed");
