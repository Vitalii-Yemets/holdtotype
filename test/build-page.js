const fs = require("fs");
const path = require("path");

const root = path.join(__dirname, "..");
const src = fs.readFileSync(path.join(root, "client", "settings.go"), "utf8");

const marker = "const settingsPage = `";
const start = src.indexOf(marker) + marker.length;
const end = src.indexOf("`", start);
if (start < marker.length || end < 0) {
  console.error("settingsPage template not found in client/settings.go");
  process.exit(1);
}

const cfg = {
  hotkey: "ctrl+win",
  ui_language: "en",
  ui_level: "all",
  skin: "terminal",
  theme: "green",
  hotkey_mode: "hold",
  language: "ru",
  model: "models/ggml-large-v3-turbo-q5_0.bin",
  lang_models: { auto: "large-v3-turbo-q5_0", ru: "gigaam-v3" },
  beep: true,
  sound_theme: "speech",
  auto_enter: false,
  restore_clipboard: true,
  overlay: true,
  overlay_position: "bottom",
  overlay_monitor: "",
  overlay_custom: {},
  _monitors: [{ index: 0, w: 1920, h: 1080, primary: true, name: "DELL U2720Q" }, { index: 1, w: 2560, h: 1440, primary: false, name: "" }],
  overlay_text: true,
  type_mode: false,
  threads: 4,
  min_record_ms: 150,
  max_record_seconds: 60,
  server_autostart: true,
  check_updates: false,
  server_port: 8910,
  server_exe: "whisper-server.exe",
  server_url: "",
  stt_source: "local",
  whisper_prompt: "GitHub, Docker",
  translate_default: false,
  translate_hotkey: "",
  translate_target: "en",
  translate_ask: "never",
  translate_ask_seconds: 3,
  translate_ask_langs: ["en", "ru"],
  llm_model: "model.gguf",
  punctuation: "model",
  active_profiles: ["clean"],
  replacements: [],
  commands: [],
  history: false,
  history_keep_min: 10080,
  history_max: 200,
  history_skip: "",
  post_enabled: true,
  post_source: "local",
  post_api_url: "",
  post_api_model: "",
  post_api_timeout_s: 30,
  paste_delay_ms: 0,
  profiles: [
    { id: "clean", name: "Cleanup", prompt: "p1", hotkey: "" },
    { id: "formal", name: "Business", prompt: "p2", hotkey: "" },
  ],
  _mail: "holdtotype@outlook.com",
  _site: "https://holdtotype.com",
  _repo: "https://github.com/Vitalii-Yemets/holdtotype",
  _version: "0.0.0-test",
  _tab: "general",
  _wizard: false,
  _cpus: 8,
};

const i18nSrc = fs.readFileSync(path.join(root, "client", "i18n.go"), "utf8");

const NAMES = {
  "{app}": "HoldToType",
  "{exe}": "holdtotype.exe",
  "{setup}": "holdtotype-setup.exe",
  "{log}": "holdtotype.log",
  "{zip}": "holdtotype-portable.zip",
};

function goText(raw) {
  const map = { '"': '"', "\\": "\\", n: "\n", r: "\r", t: "\t" };
  const text = raw.replace(/\\(.)/g, (m, c) => (c in map ? map[c] : m));
  return text.replace(/\{(app|exe|setup|log|zip)\}/g, (m) => NAMES[m]);
}

const tableStart = i18nSrc.indexOf("var settingsStrings");
const enStart = tableStart < 0 ? -1 : i18nSrc.indexOf('\n\t"en": {', tableStart);
const enEnd = enStart < 0 ? -1 : i18nSrc.indexOf("\n\t},", enStart);
if (enStart < 0 || enEnd < 0) {
  console.error("the English settings strings were not found in client/i18n.go");
  process.exit(1);
}
const EN = {};
for (const m of i18nSrc.slice(enStart, enEnd).matchAll(/"(S_[A-Z0-9_]+)":\s*"((?:[^"\\]|\\.)*)"/g)) {
  EN[m[1]] = goText(m[2]);
}

const lStart = src.indexOf("lMap := map[string]string{");
const lEnd = src.indexOf("lJSON, _ := json.Marshal(lMap)");
if (lStart < 0 || lEnd < 0) {
  console.error("the L table was not found in client/settings.go");
  process.exit(1);
}
const strings = { nohot: "—" };
const noText = [];
for (const m of src.slice(lStart, lEnd).matchAll(/"([a-z0-9_]+)":\s*"(S_[A-Z0-9_]+)"/g)) {
  if (EN[m[2]] === undefined) {
    noText.push(m[1] + " → " + m[2]);
    continue;
  }
  strings[m[1]] = EN[m[2]];
}
if (noText.length) {
  console.error("no English text for: " + noText.join(", "));
  process.exit(1);
}

const HELP_HTML = [
  '<p class="wh">How it works</p>',
  '<p>Hold the shortcut, say a phrase, let go — the text lands at the caret.</p>',
  '<p class="wh">Overlay</p>',
  '<ul><li>Every answer carries a number: 1…9 pick one, Enter takes the highlighted one.</li></ul>',
  '<p class="wh">Tray and files</p>',
  '<ul><li><b>config.json</b> — all settings; edits made by hand apply through Re-read in the System section.</li></ul>',
  '<p class="wh">Install and portability</p>',
  '<p>The installer downloads nothing by default.</p>',
].join("");

const PALETTES = {
  green: { bg: "#0b0f0c", panel: "#0e1410", line: "#1d4a2b", text: "#3cff6e", hi: "#3cff6e", dim: "#20a34a", faint: "#14803a", warn: "#ffb347", bad: "#ff7b6b", rec: "#ff5b4d", field: "#08100b", soft: "#12241a", navon: "#101d14", on: "#123f22" },
  amber: { bg: "#100c0a", panel: "#17110d", line: "#4a3018", text: "#ff9e2c", hi: "#ff9e2c", dim: "#b56a12", faint: "#8a4f0d", warn: "#ffd24a", bad: "#ff6b5b", rec: "#ff5b4d", field: "#120c07", soft: "#2a1a0d", navon: "#22160c", on: "#402611" },
  blue: { bg: "#0b0e10", panel: "#0e1317", line: "#1d3a4a", text: "#4cc3ff", hi: "#4cc3ff", dim: "#1c7fb8", faint: "#14608f", warn: "#ffb347", bad: "#ff7b6b", rec: "#ff5b4d", field: "#070f14", soft: "#12222c", navon: "#101c24", on: "#123a52" },
  pink: { bg: "#100b0e", panel: "#170e14", line: "#4a1d3a", text: "#ff6ec7", hi: "#ff6ec7", dim: "#b82f86", faint: "#8f2467", warn: "#ffb347", bad: "#ff6b6b", rec: "#ff5b4d", field: "#120810", soft: "#2a1222", navon: "#22101c", on: "#40183a" },
  editor: { bg: "#1e1e1e", panel: "#252526", line: "#3c3c3c", text: "#d4d4d4", hi: "#4fc1ff", dim: "#9d9d9d", faint: "#6e6e6e", warn: "#cca700", bad: "#f14c4c", rec: "#f14c4c", field: "#3c3c3c", card: "#252526", soft: "#2d2d2d", navon: "#37373d", on: "#094771", btnbgh: "#1177bb", btn2bg: "#0e639c", btn2fg: "#ffffff", btn2line: "transparent", btn2bgh: "#1177bb", focus: "#007fd4", dangerbg: "#a1260d", dangerfg: "#ffffff", dangerbgh: "#c42b1c", dot: "#89d185" },
  neon: { bg: "#150a22", panel: "#1d0e30", line: "#4a2472", text: "#f3b6e4", hi: "#46e0ff", dim: "#b06ee0", faint: "#7d4fae", warn: "#ffd24a", bad: "#ff4d7d", rec: "#ff4d7d", field: "#1e0f33", card: "#231039", soft: "#2a1442", navon: "#2b1240", on: "#4a2472", btnbgh: "linear-gradient(90deg,rgba(255,95,200,.34),rgba(70,224,255,.28))", btn2bg: "linear-gradient(90deg,rgba(255,95,200,.20),rgba(70,224,255,.16))", btn2fg: "#ffffff", btn2line: "#4a2472", btn2bgh: "linear-gradient(90deg,rgba(255,95,200,.34),rgba(70,224,255,.28))", dangerbg: "linear-gradient(90deg,rgba(255,77,125,.30),rgba(255,77,125,.14))", dangerfg: "#ffd7e3", dangerbgh: "linear-gradient(90deg,rgba(255,77,125,.46),rgba(255,77,125,.24))", dot: "#46e0ff", focus: "#46e0ff" },
  soft: { bg: "#ddc9d5", panel: "#e7d5df", line: "#c9a5b9", text: "#8f2f60", hi: "#b8407c", dim: "#9d6885", faint: "#ae8b9e", warn: "#7d4a00", bad: "#ab2445", rec: "#c74a6c", field: "#f4e8ee", card: "#eadbe4", soft: "#d2b9c7", navon: "#d8c0ce", on: "#e6c6d8", btnbgh: "#ad3b71", btn2bg: "#d5b6c8", btn2fg: "#6d2349", btn2line: "#bf9bb0", btn2bgh: "#cba9be", dangerbg: "#ab2445", dangerfg: "#ffffff", dangerbgh: "#8f1c39", dot: "#2e7d4f", ok: "#2e7d4f", focus: "#b8407c" },
  paper: { bg: "#e7eaee", panel: "#f3f5f7", line: "#ced5dc", text: "#1f2328", hi: "#0969da", dim: "#59636e", faint: "#818b98", warn: "#9a6700", bad: "#cf222e", rec: "#cf222e", field: "#ffffff", card: "#f7f9fa", soft: "#dde2e8", navon: "#dfe5eb", on: "#ddf4ff", ok: "#1a7f37", btnbgh: "#0860ca", btn2bg: "#d8dee5", btn2fg: "#24292f", btn2line: "#bcc5ce", btn2bgh: "#c9d1da", dangerbg: "#cf222e", dangerfg: "#ffffff", dangerbgh: "#a40e26", dot: "#1a7f37", focus: "#0969da" },
};
const SKINS = {
  terminal: { palette: "green", font: '"IBM Plex Mono",Consolas,monospace', fs: "14px", r: "0px", bw: "1px", scan: "1", shadow: "none", glow: true, round: false, caps: true, fieldpad: "6px 10px", ctlfs: "12.5px", wb: "700", level: "bars", barr: "0", mark: "mic" },
  editor: { palette: "editor", font: '"Segoe UI Variable Text","Segoe UI",system-ui,sans-serif', fs: "13px", r: "8px", dotr: "50%", badger: "10px", panelr: "6px", switchr: "999px", radius: 4, bw: "1px", scan: "0", shadow: "0 10px 30px rgba(0,0,0,.45)", glow: false, round: false, caps: false, fieldpad: "6px 11px", ctlfs: "12.5px", wb: "600", level: "flat", barr: "0", mark: "mic" },
  neon: { palette: "neon", font: '"IBM Plex Sans","Segoe UI",system-ui,sans-serif', fs: "15px", r: "14px", dotr: "50%", badger: "999px", panelr: "20px", switchr: "999px", ctlh: "36px", bw: "1px", scan: ".18", shadow: "0 18px 46px rgba(150,40,220,.35)", glow: true, round: true, caps: false, fieldpad: "9px 13px", ctlfs: "13.5px", wb: "600", level: "bars", barr: "99px", mark: "mic" },
  soft: { palette: "soft", font: '"Comic Sans MS","Segoe UI Variable Display","Segoe UI",sans-serif', fs: "15px", r: "16px", dotr: "50%", badger: "999px", panelr: "22px", switchr: "999px", ctlh: "36px", barr: "99px", bw: "1px", scan: "0", shadow: "0 14px 34px rgba(255,140,190,.28)", glow: false, round: true, caps: false, fieldpad: "9px 13px", ctlfs: "13.5px", wb: "600", level: "dots", barr: "99px", mark: "face" },
  paper: { palette: "paper", font: '"Segoe UI",system-ui,sans-serif', fs: "14px", r: "10px", dotr: "50%", badger: "999px", panelr: "12px", switchr: "999px", bw: "1px", scan: "0", shadow: "0 8px 24px rgba(31,35,40,.12)", glow: false, round: true, caps: false, fieldpad: "7px 11px", ctlfs: "12.5px", wb: "600", level: "bars", barr: "2px", mark: "mic" },
};
function rgbOf(hex) {
  return [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16)).join(",");
}
function lumaOf(hex) {
  const c = rgbOf(hex).split(",").map(Number);
  return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2];
}
function blend(fg, bg, t) {
  const hex2 = (v) => v.toString(16).padStart(2, "0");
  const mix = (i) => hex2(Math.round(parseInt(fg.slice(i, i + 2), 16) * (1 - t) + parseInt(bg.slice(i, i + 2), 16) * t));
  return "#" + mix(1) + mix(3) + mix(5);
}
function varsFor(skinId, colourId) {
  const s = SKINS[skinId];
  const p = PALETTES[s.palette === "green" ? colourId : s.palette];
  const LBL = { green: "#f2fff5", amber: "#ffe9c9", blue: "#e4f6ff", pink: "#ffe6f4", editor: "#e0e0e0", neon: "#ffffff", soft: "#7d2b56", paper: "#111418" };
  const lbl = LBL[s.palette === "green" ? colourId : s.palette] || p.text;
  const rgb = rgbOf(p.text);
  const barr = s.barr || "0";
  const lvl = s.level === "dots" ? { w: "10px", r: "50%" } : s.level === "flat" ? { w: "3px", r: "0" } : { w: "4px", r: barr };
  return [
    "--wborder:" + (s.round ? "none" : s.bw + " solid " + p.line),
    "--bg:" + p.bg, "--panel:" + p.panel, "--line:" + p.line,
    "--green:" + p.text, "--hi:" + p.hi, "--dim:" + p.dim, "--faint:" + p.faint,
    "--amber:" + p.warn, "--bad:" + p.bad, "--rec:" + p.rec,
    "--rgb:" + rgb,
    "--field:" + p.field, "--card:" + (p.card || p.field), "--soft:" + p.soft, "--navon:" + p.navon, "--on:" + p.on,
    "--titlebg:transparent", "--sidebg:transparent", "--keybg:transparent",
    "--btnbg:" + p.navon, "--btnfg:" + p.text, "--btnline:" + p.dim,
    "--btnbgh:" + (p.btnbgh || p.navon), "--btn2bg:" + (p.btn2bg || "transparent"),
    "--btn2fg:" + (p.btn2fg || p.dim), "--btn2line:" + (p.btn2line || p.dim),
    "--btn2bgh:" + (p.btn2bgh || p.btn2bg || "transparent"), "--focus:" + (p.focus || p.dim),
    "--dotr:" + (s.dotr || "0"), "--badger:" + (s.badger || "calc(" + s.r + " * .4)"),
    "--panelr:" + (s.panelr || s.r), "--switchr:" + (s.switchr || "calc(" + s.r + " * .8)"),
    "--ctlh:" + (s.ctlh || "30px"),
    "--dangerbg:" + (p.dangerbg || "transparent"), "--dangerfg:" + (p.dangerfg || p.bad),
    "--dangerbgh:" + (p.dangerbgh || blend(p.bad, p.bg, 0.8)),
    "--btnbo:" + (s.caps ? '"[ "' : '""'), "--btnbc:" + (s.caps ? '" ]"' : '""'),
    "--selbg:" + p.text, "--selfg:" + p.bg,
    "--brandbg:none", "--brandclip:border-box", "--brandfill:currentColor",
    "--scrim:rgba(3,7,4,.78)",
    "--ok:" + (p.ok || p.hi), "--scheme:" + (lumaOf(p.bg) > 140 ? "light" : "dark"),
    "--badbg:" + blend(p.bad, p.bg, 0.8), "--badline:" + blend(p.bad, p.bg, 0.52),
    "--lvlw:" + lvl.w, "--lvlr:" + lvl.r,
    "--markmic:" + (s.mark === "face" ? "none" : "block"), "--markface:" + (s.mark === "face" ? "block" : "none"),
    "--lbl:" + lbl, "--lblglow:" + (s.glow ? "0 0 7px rgba(" + rgbOf(lbl) + ",.4)" : "none"),
    "--glow:" + (s.glow ? "0 0 7px rgba(" + rgb + ",.55)" : "none"),
    "--higlow:" + (s.glow ? "0 0 8px rgba(" + rgbOf(p.hi) + ",.6)" : "none"),
    "--iconglow:" + (s.glow ? "drop-shadow(0 0 6px rgba(" + rgbOf(p.hi) + ",.7))" : "none"),
    "--amberglow:" + (s.glow ? "0 0 6px rgba(" + rgbOf(p.warn) + ",.5)" : "none"),
    "--badglow:" + (s.glow ? "0 0 7px rgba(" + rgbOf(p.bad) + ",.5)" : "none"),
    "--badfilter:" + (s.glow ? "drop-shadow(0 0 4px rgba(" + rgbOf(p.bad) + ",.5))" : "none"),
    "--font:" + s.font, "--fs:" + s.fs,
    "--caps:" + (s.caps ? "uppercase" : "none"), "--ls:" + (s.caps ? "1px" : "0"),
    "--fieldpad:" + s.fieldpad, "--ctlfs:" + s.ctlfs, "--wb:" + s.wb,
    "--flicker:" + (s.caps ? "flicker 6s infinite" : "none"),
    "--r:" + s.r, "--barr:" + (parseInt(s.r, 10) >= 10 ? "99px" : "0"),
    "--bw:" + s.bw, "--scan:" + s.scan, "--shadow:" + s.shadow, "--brandls:.18em",
  ].join(";");
}
const THEME_LIST = {};
for (const id of ["green", "amber", "blue", "pink"]) {
  THEME_LIST["terminal:" + id] = { skin: "terminal", colour: id, accent: PALETTES[id].text, vars: varsFor("terminal", id) };
}
THEME_LIST.editor = { skin: "editor", colour: "editor", accent: PALETTES.editor.text, vars: varsFor("editor") };
THEME_LIST.neon = { skin: "neon", colour: "neon", accent: PALETTES.neon.text, vars: varsFor("neon") };
THEME_LIST.soft = { skin: "soft", colour: "soft", accent: PALETTES.soft.text, vars: varsFor("soft") };
THEME_LIST.paper = { skin: "paper", colour: "paper", accent: PALETTES.paper.text, vars: varsFor("paper") };

let html = src.slice(start, end);
html = html.split("{{S_HELP_HTML}}").join(HELP_HTML);
html = html.split("{{THEME_LIST}}").join(JSON.stringify(THEME_LIST));
html = html.split("{{THEME_VARS}}").join(varsFor("terminal", "green"));
html = html.split("{{SKIN}}").join("terminal");
html = html.split("{{CFG}}").join(JSON.stringify(cfg));
html = html.split("{{L_JSON}}").join(JSON.stringify(strings));
html = html.replace(/{{([A-Z_0-9]+)}}/g, "$1");

const out = path.join(__dirname, "page.html");
fs.writeFileSync(out, html, "utf8");
console.log("page.html written:", html.length, "chars");
