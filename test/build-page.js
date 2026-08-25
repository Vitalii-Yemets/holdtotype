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
  beep: true,
  sound_theme: "speech",
  auto_enter: false,
  restore_clipboard: true,
  overlay: true,
  overlay_position: "bottom",
  overlay_text: true,
  animation: true,
  type_mode: false,
  threads: 4,
  min_record_ms: 150,
  max_record_seconds: 60,
  server_autostart: true,
  check_updates: false,
  server_port: 8910,
  server_exe: "whisper-server.exe",
  server_url: "",
  whisper_prompt: "",
  translate_default: false,
  translate_hotkey: "",
  pause_hotkey: "ctrl+alt+p",
  translate_target: "en",
  translate_ask: "never",
  translate_ask_seconds: 3,
  translate_ask_langs: ["en", "ru"],
  llm_model: "model.gguf",
  punctuation: "model",
  active_profiles: ["clean"],
  app_rules: [],
  replacements: [],
  commands: [],
  history: false,
  history_days: 7,
  history_max: 200,
  history_skip: "",
  paste_delay_ms: 0,
  profiles: [
    { id: "clean", name: "Cleanup", prompt: "p1", hotkey: "" },
    { id: "formal", name: "Business", prompt: "p2", hotkey: "" },
  ],
  _version: "0.0.0-test",
  _tab: "general",
  _wizard: false,
  _cpus: 8,
};

const strings = {
  nohot: "-", dl: "Download", del: 'Delete "x"', mdlready: "Model downloaded",
  get: "Download", change: "Change", remotewarn: "Audio will be sent to this server.",
  ok: "Yes", cancel: "Cancel", dlask: 'The "%s" model is not downloaded (%s). Start downloading?',
  dlstart: "Download", dlcancel: "Cancel the download", nofound: "none",
  advprimary: "main", advcompanion: "second", advhave: "already here", advapply: "Apply",
  advask: "These will be downloaded: %s — %s in total. Start?", nofound: "none",
  remoteask: "Send audio to %s?", remotebadge: "REMOTE",
  add: "Add profile", pname: "Name", pprompt: "Prompt", phot: "Hotkey",
  pset: "Set", pclr: "Clear", ptest: "Test", fitok: "fits", fitwarn: "tight",
  fitbad: "no RAM", ram: "RAM:", hfph: "model name", nollm: "no models",
  nollmp: "no models for prompts", upd: "updated", pedit: "Edit", pclose: "Collapse",
  confirmdel: 'Delete the "%s" model?',
  delactive: 'Delete the active model "%s"? Recognition stops until you pick another one.',
  wizneedmodel: "Download a model first",
  free: "free", updnone: "latest",
  badgemodels: "Installed models", badgemiss: "A model is not downloaded",
  badgesystem: "Warnings need attention", badgehist: "Entries in history",
  advrolemain: "for the language you picked", advrolesecond: "for other languages and translation",
  updavail: "Version %s available.", updgo: "Update", upderr: "Check failed",
  upddl: "Downloading", micdefault: "System default", micquiet: "quiet",
  more: "%d more settings", less: "Collapse %d settings",
  pasteinh: "insertion: as set", enterinh: "Enter: as set", delaynone: "no delay", promptinh: "prompts: as set",
  ruleclip: "clipboard", ruletype: "character by character",
  ruleenteron: "with Enter", ruleenteroff: "without Enter", rulenoprompt: "no prompts",
  rulelast: "last insertion: %s", ruleempty: "No rules yet", ruledel: "Delete the rule",
  ruleprompts: "Prompts", ruleph: "chrome.exe, msedge.exe",
  replempty: "No replacements yet", repldel: "Delete the replacement", replwhole: "whole words",
  replcase: "case", replfromph: "git hub", repltoph: "GitHub",
  histempty: "No history yet", histcopy: "Copy", histask: "Delete the whole dictation history?", histclear: "Clear",
  micchecking: "Checking…", mchecking: "Checking models…", histinsert: "Paste",
  retry: "Try again", slotru: "Russian speech", slotother: "Other languages",
  cmdempty: "No commands yet", cmddel: "Delete the command", cmdph: "new line",
  cmdnewline: "line break", cmdparagraph: "new paragraph", cmdtext: "insert text", cmdcancel: "cancel the dictation",
  cmdtextph: "what to insert",
  cmdpnewline: "new line", cmdpparagraph: "new paragraph", cmdpcancel: "cancel",
  wiznext: "Next", wizfinish: "Finish", wizwait: "Waiting for the first phrase…",
  wizheard: "Heard:", wizhave: "Everything you need is already downloaded",
  wiztry: "Put the caret in the field below, hold %s, say a phrase and let go.",
  wndmax: "S_WND_MAX", wndrestore: "S_WND_RESTORE", wndmin: "S_WND_MIN", wndclose: "S_WND_CLOSE",
};

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
  editor: { bg: "#1e1e1e", panel: "#252526", line: "#3c3c3c", text: "#d4d4d4", hi: "#4fc1ff", dim: "#9d9d9d", faint: "#6e6e6e", warn: "#cca700", bad: "#f14c4c", rec: "#f14c4c", field: "#3c3c3c", soft: "#2d2d2d", navon: "#37373d", on: "#094771" },
  neon: { bg: "#150a22", panel: "#1d0e30", line: "#4a2472", text: "#f3b6e4", hi: "#46e0ff", dim: "#b06ee0", faint: "#7d4fae", warn: "#ffd24a", bad: "#ff4d7d", rec: "#ff4d7d", field: "#1e0f33", soft: "#2a1442", navon: "#2b1240", on: "#4a2472" },
};
const SKINS = {
  terminal: { palette: "green", font: 'Consolas,"Cascadia Mono",monospace', fs: "14px", r: "0px", bw: "1px", scan: "1", shadow: "none", glow: true, round: false, caps: true },
  editor: { palette: "editor", font: '"Cascadia Mono",Consolas,"Segoe UI",sans-serif', fs: "13px", r: "3px", bw: "1px", scan: "0", shadow: "0 10px 30px rgba(0,0,0,.45)", glow: false, round: false, caps: false },
  neon: { palette: "neon", font: '"Segoe UI Variable Display","Segoe UI",system-ui,sans-serif', fs: "13px", r: "14px", bw: "1px", scan: ".18", shadow: "0 18px 46px rgba(150,40,220,.35)", glow: true, round: true, caps: false },
};
function rgbOf(hex) {
  return [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16)).join(",");
}
function varsFor(skinId, colourId) {
  const s = SKINS[skinId];
  const p = PALETTES[s.palette === "green" ? colourId : s.palette];
  const rgb = rgbOf(p.text);
  return [
    "--wborder:" + (s.round ? "none" : s.bw + " solid " + p.line),
    "--bg:" + p.bg, "--panel:" + p.panel, "--line:" + p.line,
    "--green:" + p.text, "--hi:" + p.hi, "--dim:" + p.dim, "--faint:" + p.faint,
    "--amber:" + p.warn, "--bad:" + p.bad, "--rec:" + p.rec,
    "--rgb:" + rgb,
    "--field:" + p.field, "--soft:" + p.soft, "--navon:" + p.navon, "--on:" + p.on,
    "--titlebg:transparent", "--sidebg:transparent", "--keybg:transparent",
    "--btnbg:" + p.navon, "--btnfg:" + p.text, "--btnline:" + p.dim,
    "--selbg:" + p.text, "--selfg:" + p.bg,
    "--brandbg:none", "--brandclip:border-box", "--brandfill:currentColor",
    "--scrim:rgba(3,7,4,.78)",
    "--glow:" + (s.glow ? "0 0 7px rgba(" + rgb + ",.55)" : "none"),
    "--higlow:" + (s.glow ? "0 0 8px rgba(" + rgbOf(p.hi) + ",.6)" : "none"),
    "--iconglow:" + (s.glow ? "drop-shadow(0 0 6px rgba(" + rgbOf(p.hi) + ",.7))" : "none"),
    "--font:" + s.font, "--fs:" + s.fs,
    "--caps:" + (s.caps ? "uppercase" : "none"), "--ls:" + (s.caps ? "1px" : "0"),
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

let html = src.slice(start, end);
html = html.split("{{S_HELP_HTML}}").join(HELP_HTML);
html = html.split("{{THEME_LIST}}").join(JSON.stringify(THEME_LIST));
html = html.split("{{THEME_VARS}}").join(varsFor("terminal", "green"));
html = html.split("{{CFG}}").join(JSON.stringify(cfg));
html = html.split("{{L_JSON}}").join(JSON.stringify(strings));
html = html.replace(/{{([A-Z_0-9]+)}}/g, "$1");

const out = path.join(__dirname, "page.html");
fs.writeFileSync(out, html, "utf8");
console.log("page.html written:", html.length, "chars");
