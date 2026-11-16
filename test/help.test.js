const fs = require("fs");
const path = require("path");
const { JSDOM } = require("jsdom");

const root = path.join(__dirname, "..", "client");
const goSrc = fs.readFileSync(path.join(root, "i18n.go"), "utf8") + fs.readFileSync(path.join(root, "i18n_extra.go"), "utf8");
const settings = fs.readFileSync(path.join(root, "settings.go"), "utf8");
const LANGS = ["en", "ru", "uk", "de", "fr", "es", "it", "pl"];
const NAMES = { "{app}": "HoldToType", "{exe}": "holdtotype.exe", "{setup}": "holdtotype-setup.exe", "{log}": "holdtotype.log", "{zip}": "holdtotype-portable.zip" };

const TABS = { state: "s-state", system: "s-system", library: "s-library", mic: "s-mic", history: "s-hist", dictation: "s-dict", models: "s-models", text: "s-rules", post: "s-post", subs: "s-subs" };

const SKIP = {
  S_SUBS_OFF: "the heading of the subtitles block on its tab is their state, not a setting",
};

const EXTRA = {
  system: { keys: ["S_THREADS", "S_AUTOSTART", "S_PORT", "S_SERVEREXE"], why: "the rows of the recognition server dialog are built by the page script" },
  mic: { keys: ["S_MIC_REFRESH"], why: "the refresh button of the microphone list is built by the page script" },
  models: { keys: ["S_LANG_KEY"], why: "the shortcut row of an unfolded language is built by the page script" },
};

const failures = [];
let passed = 0;
function check(name, actual, expected) {
  const a = JSON.stringify(actual), e = JSON.stringify(expected);
  if (a === e) { passed++; console.log("PASS  " + name); }
  else { failures.push(name); console.log("FAIL  " + name + " (got " + a + ", want " + e + ")"); }
}

function goText(raw) {
  const map = { '"': '"', "\\": "\\", n: "\n", r: "\r", t: "\t" };
  return raw.replace(/\\(.)/g, (m, c) => (c in map ? map[c] : m)).replace(/\{(app|exe|setup|log|zip)\}/g, (m) => NAMES[m]);
}
function blockStart(lang) {
  let s = goSrc.indexOf('settingsStrings["' + lang + '"] = map[string]string{');
  if (s < 0) s = goSrc.indexOf('\n\t"' + lang + '": {', goSrc.indexOf("var settingsStrings"));
  return s;
}
function stringsOf(lang) {
  const s = blockStart(lang);
  const e = goSrc.indexOf("\n\t}", s + 10);
  const out = {};
  for (const m of goSrc.slice(s, e).matchAll(/"(S_[A-Z0-9_]+)":\s*"((?:[^"\\]|\\.)*)"/g)) out[m[1]] = goText(m[2]);
  return out;
}
function helpOf(lang) {
  const key = '"S_HELP_HTML":';
  const piece = /\s*(\+)?\s*"((?:[^"\\]|\\.)*)"/y;
  let at = goSrc.indexOf(key, blockStart(lang)) + key.length, out = "";
  for (;;) {
    piece.lastIndex = at;
    const m = piece.exec(goSrc);
    if (!m || (out && !m[1])) break;
    out += goText(m[2]);
    at = piece.lastIndex;
  }
  return out;
}

const norm = (s) => String(s).replace(/<[^>]+>/g, " ").replace(/[«»“”„"‚‘’'`]/g, "").replace(/[  \s]+/g, " ").trim().toLowerCase();

const marker = "const settingsPage = `";
const pageSrc = settings.slice(settings.indexOf(marker) + marker.length, settings.indexOf("`", settings.indexOf(marker) + marker.length));
const page = new JSDOM(pageSrc).window.document;
const keyIn = (text) => { const m = /\{\{(S_[A-Z0-9_]+)\}\}/.exec(text || ""); return m ? m[1] : null; };
const ownText = (el) => [...el.childNodes].filter((n) => n.nodeType === 3).map((n) => n.textContent).join(" ");

function namesOf(tab) {
  const p = page.getElementById("p-" + tab);
  const out = new Map();
  const add = (el, kind) => {
    const k = keyIn(ownText(el));
    if (k && !(k in SKIP)) out.set(k, kind);
  };
  const all = (sel, kind) => p.querySelectorAll(sel).forEach((el) => { if (!el.closest(".hshot")) add(el, kind); });
  all(".row > label, .row > div > label:not(.sizetile)", "row");
  all("label.blklbl", "block");
  all(".srchead", "card");
  all("button", "button");
  all("label.stsw > span", "switch");
  ((EXTRA[tab] || {}).keys || []).forEach((k) => out.set(k, "script"));
  return out;
}

function sectionOf(help, id) {
  const at = help.indexOf('id="' + id + '"');
  if (at < 0) return "";
  const start = help.lastIndexOf("<p", at);
  const next = help.indexOf('<p class="wh"', at);
  return help.slice(start, next < 0 ? help.length : next);
}

const report = {};
for (const lang of LANGS) {
  const S = stringsOf(lang);
  const help = helpOf(lang);
  for (const [tab, id] of Object.entries(TABS)) {
    const sec = sectionOf(help, id);
    const text = norm(sec);
    const names = namesOf(tab);
    const missing = [];
    const labels = [];
    for (const [k, kind] of names) {
      const v = S[k];
      if (typeof v !== "string" || !v.trim() || /%[sd]/.test(v)) continue;
      labels.push(norm(v));
      if (!text.includes(norm(v))) missing.push(kind + " " + k + " «" + v + "»");
    }
    const strange = [];
    for (const m of sec.matchAll(/<tr><td>(.*?)<\/td>/g)) {
      const first = norm(m[1]);
      const parts = first.split(/\s+[·/]\s+/);
      if (!parts.every((p) => labels.some((l) => l === p || l.startsWith(p + " ") || p.startsWith(l)))) strange.push("«" + m[1].replace(/<[^>]+>/g, "") + "»");
    }
    report[lang + " " + tab] = { missing, strange };
    check(lang + ": the help section of the " + tab + " tab names every row, block, card, switch and button of the tab", missing, []);
    check(lang + ": every row of the help tables of the " + tab + " tab names something on the tab", strange, []);
  }
}

if (process.env.HELP_REPORT) fs.writeFileSync(process.env.HELP_REPORT, JSON.stringify(report, null, 1));
if (failures.length) {
  console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
  process.exit(1);
}
console.log("\nall help checks passed");
