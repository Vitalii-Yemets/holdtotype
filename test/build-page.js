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

const THEME_LIST = {
  green: { bg: "#0b0f0c", panel: "#0e1410", line: "#1d4a2b", accent: "#3cff6e", dim: "#20a34a", faint: "#14803a", warn: "#ffb347", bad: "#ff7b6b", rgb: "60,255,110", glow: "0 0 7px rgba(60,255,110,.55)" },
  amber: { bg: "#100c0a", panel: "#17110d", line: "#4a3018", accent: "#ff9e2c", dim: "#b56a12", faint: "#8a4f0d", warn: "#ffd24a", bad: "#ff6b5b", rgb: "255,158,44", glow: "0 0 7px rgba(255,158,44,.55)" },
  blue: { bg: "#0b0e10", panel: "#0e1317", line: "#1d3a4a", accent: "#4cc3ff", dim: "#1c7fb8", faint: "#14608f", warn: "#ffb347", bad: "#ff7b6b", rgb: "76,195,255", glow: "0 0 7px rgba(76,195,255,.55)" },
  pink: { bg: "#100b0e", panel: "#170e14", line: "#4a1d3a", accent: "#ff6ec7", dim: "#b82f86", faint: "#8f2467", warn: "#ffb347", bad: "#ff6b6b", rgb: "255,110,199", glow: "0 0 7px rgba(255,110,199,.55)" },
};

let html = src.slice(start, end);
html = html.split("{{S_HELP_HTML}}").join(HELP_HTML);
html = html.split("{{THEME_LIST}}").join(JSON.stringify(THEME_LIST));
html = html.split("{{THEME_VARS}}").join("--bg:#0b0f0c;--panel:#0e1410;--line:#1d4a2b;--green:#3cff6e;--dim:#20a34a;--faint:#14803a;--amber:#ffb347;--bad:#ff7b6b;--rgb:60,255,110;--glow:0 0 7px rgba(60,255,110,.55)");
html = html.split("{{CFG}}").join(JSON.stringify(cfg));
html = html.split("{{L_JSON}}").join(JSON.stringify(strings));
html = html.replace(/{{([A-Z_0-9]+)}}/g, "$1");

const out = path.join(__dirname, "page.html");
fs.writeFileSync(out, html, "utf8");
console.log("page.html written:", html.length, "chars");
