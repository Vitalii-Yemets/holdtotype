const fs = require("fs");
const path = require("path");

const root = path.join(__dirname, "..", "client");
const main = fs.readFileSync(path.join(root, "i18n.go"), "utf8");
const extra = fs.readFileSync(path.join(root, "i18n_extra.go"), "utf8");
const settings = fs.readFileSync(path.join(root, "settings.go"), "utf8");

// How many page keys each locale is still missing today. The numbers may only
// go down: a locale that falls further behind fails the check.
const SETTINGS_DEBT = { ru: 0, uk: 0, de: 0, fr: 0, es: 0, it: 0, pl: 0 };

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
const placeholders = [...settings.matchAll(/{{(S_[A-Z_0-9]+)}}/g)].map((m) => m[1]);
const jsStrings = [...settings.matchAll(/"[a-z]+": "(S_[A-Z_0-9]+)"/g)].map((m) => m[1]);
const page = new Set([...placeholders, ...jsStrings]);

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
  const missingMsgs = [...enMsgs].filter((k) => !msgs.has(k));
  check(`${locale} says everything the program can say`, missingMsgs, []);
}

function helpOf(locale) {
  const block = settingsBlock(locale) || "";
  const at = block.indexOf('"S_HELP_HTML":');
  if (at < 0) return "";
  const rest = block.slice(at);
  const end = rest.search(/",\r?\n/);
  return end < 0 ? rest : rest.slice(0, end);
}
function helpShape(locale) {
  const h = helpOf(locale);
  const count = (re) => (h.match(re) || []).length;
  return {
    sections: [...h.matchAll(/id=\\"(s-[a-z]+)\\"/g)].map((m) => m[1]),
    subheads: count(/class=\\"hsub\\"/g),
    rows: count(/<tr>/g),
    items: count(/<li>/g),
    warnings: count(/class=\\"hwarn\\"/g),
  };
}
const enHelp = helpShape("en");
check("the English help has the subtitles section after post-processing", enHelp.sections.indexOf("s-subs") === enHelp.sections.indexOf("s-post") + 1, true);
for (const locale of Object.keys(SETTINGS_DEBT)) {
  check(`${locale} help has the same shape as the English one`, helpShape(locale), enHelp);
}

if (failures.length) {
  console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
  process.exit(1);
}
console.log("\nall locale checks passed");
