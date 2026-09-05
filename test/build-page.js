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
  green: { bg: "#0b0f0c", panel: "#0e1410", line: "#1d4a2b", text: "#3cff6e", hi: "#3cff6e", dim: "#20a34a", faint: "#14803a", warn: "#ffb347", bad: "#ff7b6b", rec: "#ff5b4d", field: "#08100b", soft: "#12241a", navon: "#101d14", on: "#123f22", keyfg: "#f2fff5" },
  amber: { bg: "#100c0a", panel: "#17110d", line: "#4a3018", text: "#ff9e2c", hi: "#ff9e2c", dim: "#b56a12", faint: "#8a4f0d", warn: "#ffd24a", bad: "#ff6b5b", rec: "#ff5b4d", field: "#120c07", soft: "#2a1a0d", navon: "#22160c", on: "#402611", keyfg: "#ffe9c9" },
  blue: { bg: "#0b0e10", panel: "#0e1317", line: "#1d3a4a", text: "#4cc3ff", hi: "#4cc3ff", dim: "#1c7fb8", faint: "#14608f", warn: "#ffb347", bad: "#ff7b6b", rec: "#ff5b4d", field: "#070f14", soft: "#12222c", navon: "#101c24", on: "#123a52", keyfg: "#e4f6ff" },
  pink: { bg: "#100b0e", panel: "#170e14", line: "#4a1d3a", text: "#ff6ec7", hi: "#ff6ec7", dim: "#b82f86", faint: "#8f2467", warn: "#ffb347", bad: "#ff6b6b", rec: "#ff5b4d", field: "#120810", soft: "#2a1222", navon: "#22101c", on: "#40183a", keyfg: "#ffe6f4" },
  editor: { bg: "#1b1f24", panel: "#20252b", line: "#2c333c", text: "#d7dce2", hi: "#4fc1ff", dim: "#9aa4b1", faint: "#6b7480", warn: "#e3b341", bad: "#e5534b", rec: "#e5534b", field: "#161a1f", card: "#20252b", soft: "#262c34", navon: "#2a3139", on: "#1e3a55", btnbgh: "#3a8fd8", btn2bg: "#262c34", btn2fg: "#d7dce2", btn2line: "#39414b", btn2bgh: "#2f3640", focus: "#4fc1ff", swbg: "#161a1f", swonbg: "#2f81c8", swknob: "#ffffff", swonline: "#4fc1ff", keyfg: "#4fc1ff", keyline: "#39414b", dangerbg: "#b23a3a", dangerfg: "#ffffff", dangerbgh: "#c44545", dot: "#7ed49b", ok: "#7ed49b" },
  neon: { bg: "#120a1e", panel: "#190f2b", line: "#3a2360", text: "#f1e4ff", hi: "#46e0ff", dim: "#b79bd6", faint: "#7e63a3", warn: "#ffd24a", bad: "#ff4d7d", rec: "#ff4d7d", field: "#150b25", card: "#190f2b", soft: "#24163d", navon: "#221540", on: "#2b1a4a", btnbgh: "linear-gradient(90deg,#ff7ad2,#6be7ff)", btn2bg: "#1f1236", btn2fg: "#f1e4ff", btn2line: "#4a2f78", btn2bgh: "#291a45", dangerbg: "#ff4d7d", dangerfg: "#1b0710", dangerbgh: "#ff6690", dot: "#5cf2c4", ok: "#5cf2c4", focus: "#46e0ff", swbg: "#150b25", swonbg: "#2b1a4a", swknob: "#46e0ff", swonline: "#46e0ff", keyfg: "#46e0ff", keyline: "#46e0ff", tagfg: "#ff8fd9", tagline: "#ff5fc8" },
  soft: { bg: "#f2eef8", panel: "#fbf9fe", line: "#e0d8ee", text: "#2b2438", hi: "#7c5cff", dim: "#6f6684", faint: "#9a92ab", warn: "#b26a00", bad: "#c93d64", rec: "#e2557a", field: "#ffffff", card: "#fbf9fe", soft: "#e6dff3", navon: "#e8e0f6", on: "#e8e0f6", btnbgh: "#6c4cf0", btn2bg: "#ebe5f7", btn2fg: "#4a3f66", btn2line: "transparent", btn2bgh: "#e1d9f2", dangerbg: "#e2557a", dangerfg: "#ffffff", dangerbgh: "#d1436a", dot: "#2e9e6a", ok: "#2e9e6a", focus: "#7c5cff", swbg: "#e6dff3", swonbg: "#7c5cff", swknob: "#ffffff", swonline: "#7c5cff", keyfg: "#7c5cff", keyline: "#d8cdf2", tagfg: "#6f6684", tagline: "transparent", tagbg: "#efe9fb" },
  paper: { bg: "#f4f1ea", panel: "#fbfaf6", line: "#d9d4c8", text: "#1d1d1b", hi: "#1f4fbf", dim: "#5d6067", faint: "#8b8e94", warn: "#8a5a00", bad: "#b3261e", rec: "#b3261e", field: "#ffffff", card: "#fbfaf6", soft: "#e9e5db", navon: "#e6e1d5", on: "#dfe6f7", ok: "#1f7a3f", btnbgh: "#1a45a8", btn2bg: "#e9e5db", btn2fg: "#1d1d1b", btn2line: "#cfc9bb", btn2bgh: "#dfdacd", dangerbg: "#b3261e", dangerfg: "#ffffff", dangerbgh: "#9a1f18", dot: "#1f7a3f", focus: "#1f4fbf", swbg: "#e9e5db", swonbg: "#1f4fbf", swknob: "#ffffff", swonline: "#1f4fbf", keyfg: "#1f4fbf", keyline: "#c7c1b3", tagfg: "#5d6067", tagline: "#cfc9bb", tagbg: "#ffffff" },
  fluent: { bg: "#202020", panel: "#2b2b2b", line: "#3a3a3a", text: "#ffffff", hi: "#60cdff", dim: "#c5c5c5", faint: "#8a8a8a", warn: "#fce100", bad: "#ff99a4", rec: "#ff6b6b", field: "#323232", card: "#2b2b2b", soft: "#323232", navon: "#2d2d2d", on: "#24404f", btnbgh: "#56b8e6", btn2bg: "#323232", btn2fg: "#ffffff", btn2line: "#3f3f3f", btn2bgh: "#3a3a3a", dangerbg: "#c42b1c", dangerfg: "#ffffff", dangerbgh: "#d33a2b", dot: "#6ccb5f", ok: "#6ccb5f", focus: "#60cdff", swbg: "transparent", swonbg: "#60cdff", swknob: "#000000", swonline: "#60cdff", keyfg: "#60cdff", keyline: "#454545", tagfg: "#c5c5c5", tagline: "transparent", tagbg: "#323232" },
  studio: { bg: "#1c1d1f", panel: "#242527", line: "#343538", text: "#e6e3dc", hi: "#ff9f43", dim: "#a19d94", faint: "#6f6c66", warn: "#ffd166", bad: "#e06060", rec: "#d64545", field: "#161718", card: "#242527", soft: "#2c2d30", navon: "#2c2d30", on: "#3a2a16", btnbgh: "#ffb066", btn2bg: "#2c2d30", btn2fg: "#e6e3dc", btn2line: "#45464a", btn2bgh: "#35363a", dangerbg: "#d64545", dangerfg: "#ffffff", dangerbgh: "#e05555", dot: "#7bd88f", ok: "#7bd88f", focus: "#ff9f43", swbg: "#161718", swonbg: "#3a2a16", swknob: "#ff9f43", swonline: "#ff9f43", keyfg: "#ff9f43", keyline: "#45464a", tagfg: "#ff9f43", tagline: "#5a4a33", tagbg: "#161718" },
};
const SKINS = {
  terminal: { palette: "green", font: '"IBM Plex Mono",Consolas,monospace', fs: "14px", r: "0px", bw: "1px", scan: "1", shadow: "none", glow: true, round: false, caps: true, fieldpad: "6px 10px", ctlfs: "12.5px", wb: "700", level: "bars", barr: "0", mark: "mic" },
  editor: { palette: "editor", font: '"IBM Plex Mono",Consolas,monospace', fs: "13px", r: "8px", dotr: "50%", badger: "6px", panelr: "8px", switchr: "999px", btnr: "6px", fieldr: "6px", cardr: "8px", keyr: "6px", bw: "1px", scan: "0", shadow: "0 10px 30px rgba(0,0,0,.45)", glow: false, round: true, caps: false, fieldpad: "6px 11px", ctlfs: "12.5px", wb: "600", level: "flat", barr: "1px", mark: "mic" },
  neon: { palette: "neon", font: '"IBM Plex Mono",Consolas,monospace', fs: "15px", r: "10px", dotr: "50%", badger: "999px", panelr: "14px", switchr: "999px", ctlh: "36px", btnr: "10px", fieldr: "10px", cardr: "14px", keyr: "10px", bw: "1px", scan: ".18", shadow: "0 18px 46px rgba(120,40,220,.35)", glow: true, round: true, caps: false, fieldpad: "8px 14px", ctlfs: "13.5px", wb: "600", level: "bars", barr: "3px", mark: "mic" },
  fluent: { palette: "fluent", font: '"IBM Plex Mono",Consolas,monospace', fs: "13px", r: "8px", dotr: "50%", badger: "4px", panelr: "8px", switchr: "999px", btnr: "4px", fieldr: "4px", cardr: "8px", keyr: "4px", bw: "1px", scan: "0", shadow: "0 32px 64px rgba(0,0,0,.5)", glow: false, round: true, caps: false, fieldpad: "6px 11px", ctlfs: "12.5px", wb: "600", level: "bars", barr: "2px", mark: "mic" },
  studio: { palette: "studio", font: '"IBM Plex Mono",Consolas,monospace', fs: "13px", r: "4px", dotr: "50%", badger: "3px", panelr: "4px", switchr: "999px", btnr: "3px", fieldr: "3px", cardr: "4px", keyr: "3px", bw: "1px", scan: "0", shadow: "0 24px 60px rgba(0,0,0,.6)", glow: false, round: true, caps: false, fieldpad: "6px 11px", ctlfs: "12.5px", wb: "600", level: "bars", barr: "1px", mark: "mic" },
  soft: { palette: "soft", font: '"IBM Plex Mono",Consolas,monospace', fs: "15px", r: "20px", dotr: "50%", badger: "999px", panelr: "20px", switchr: "999px", ctlh: "36px", btnr: "14px", fieldr: "14px", cardr: "20px", keyr: "14px", bw: "1px", scan: "0", shadow: "0 14px 34px rgba(124,92,255,.18)", glow: false, round: true, caps: false, fieldpad: "8px 14px", ctlfs: "13.5px", wb: "600", level: "dots", barr: "99px", mark: "mic" },
  paper: { palette: "paper", font: '"IBM Plex Mono",Consolas,monospace', fs: "14px", r: "8px", dotr: "50%", badger: "4px", panelr: "8px", switchr: "999px", btnr: "6px", fieldr: "6px", cardr: "8px", keyr: "6px", bw: "1px", scan: "0", shadow: "0 12px 30px rgba(40,35,20,.14)", glow: false, round: true, caps: false, fieldpad: "7px 11px", ctlfs: "12.5px", wb: "600", level: "bars", barr: "1px", mark: "mic" },
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
  const LBL = { green: "#f2fff5", amber: "#ffe9c9", blue: "#e4f6ff", pink: "#ffe6f4", editor: "#e6eaef", neon: "#ffffff", soft: "#241d33", paper: "#111110", fluent: "#ffffff", studio: "#f3f1ea" };
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
    "--btnr:" + (s.btnr || "calc(" + s.r + " * .5)"), "--fieldr:" + (s.fieldr || "calc(" + s.r + " * .55)"),
    "--cardr:" + (s.cardr || "calc(" + s.r + " * .6)"), "--keyr:" + (s.keyr || "calc(" + s.r + " * .6)"), "--keyfg:" + (p.keyfg || p.text), "--keyline:" + (p.keyline || p.line),
    "--tagfg:" + (p.tagfg || lbl), "--tagline:" + (p.tagline || p.line), "--tagbg:" + (p.tagbg || "transparent"),
    "--swbg:" + (p.swbg || "transparent"), "--swonbg:" + (p.swonbg || "transparent"), "--swknob:" + (p.swknob || p.hi), "--swonline:" + (p.swonline || p.dim),
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
THEME_LIST.fluent = { skin: "fluent", colour: "fluent", accent: PALETTES.fluent.text, vars: varsFor("fluent") };
THEME_LIST.studio = { skin: "studio", colour: "studio", accent: PALETTES.studio.text, vars: varsFor("studio") };

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
