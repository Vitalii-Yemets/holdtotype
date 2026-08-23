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
  hotkey_mode: "hold",
  language: "ru",
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
  translate_target: "en",
  translate_ask: "never",
  translate_ask_seconds: 3,
  translate_ask_langs: ["en", "ru"],
  llm_model: "model.gguf",
  punctuation: "model",
  active_profiles: ["clean"],
  app_rules: [],
  replacements: [],
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
  confirmdel: 'Delete the "%s" model?', free: "free", updnone: "latest",
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
  wiznext: "Next", wizfinish: "Finish", wizwait: "Waiting for the first phrase…",
  wizheard: "Heard:", wizhave: "Everything you need is already downloaded",
  wiztry: "Put the caret in the field below, hold %s, say a phrase and let go.",
};

let html = src.slice(start, end);
html = html.split("{{CFG}}").join(JSON.stringify(cfg));
html = html.split("{{L_JSON}}").join(JSON.stringify(strings));
html = html.replace(/{{([A-Z_0-9]+)}}/g, "$1");

const out = path.join(__dirname, "page.html");
fs.writeFileSync(out, html, "utf8");
console.log("page.html written:", html.length, "chars");
