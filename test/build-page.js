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
  language: "ru",
  beep: true,
  sound_theme: "speech",
  auto_enter: false,
  restore_clipboard: true,
  overlay: true,
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
  active_profiles: ["clean"],
  profiles: [
    { id: "clean", name: "Cleanup", prompt: "p1", hotkey: "" },
    { id: "formal", name: "Business", prompt: "p2", hotkey: "" },
  ],
  _version: "0.0.0-test",
  _tab: "general",
  _cpus: 8,
};

const strings = {
  nohot: "-", dl: "Download", del: 'Delete "x"', hint: "Applied on save",
  add: "Add profile", pname: "Name", pprompt: "Prompt", phot: "Hotkey",
  pset: "Set", pclr: "Clear", ptest: "Test", fitok: "fits", fitwarn: "tight",
  fitbad: "no RAM", ram: "RAM:", hfph: "model name", nollm: "no models",
  nollmp: "no models for prompts", upd: "updated", pedit: "Edit", pclose: "Collapse",
  confirmdel: 'Delete the "%s" model?', free: "free", updnone: "latest",
  updavail: "Version %s available.", updgo: "Update", upderr: "Check failed",
  upddl: "Downloading",
};

let html = src.slice(start, end);
html = html.split("{{CFG}}").join(JSON.stringify(cfg));
html = html.split("{{L_JSON}}").join(JSON.stringify(strings));
html = html.replace(/\{\{[A-Z_0-9]+\}\}/g, "X");

const out = path.join(__dirname, "page.html");
fs.writeFileSync(out, html, "utf8");
console.log("page.html written:", html.length, "chars");
