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

const strings = {
  nohot: "-", dl: "Download", del: 'Delete "x"', mdlready: "Model downloaded",
  get: "Download", change: "Change", remotewarn: "Audio will be sent to this server.",
  ok: "Yes", cancel: "Cancel", dlask: 'The "%s" model is not downloaded (%s). Start downloading?',
  dlstart: "Download", dlcancel: "Cancel the download", nofound: "none",
  asauto: "as Auto-detect",
  recchip: "recommended", backauto: "Back to Auto-detect",
  langsc: "languages: %d", langsq: "languages: unknown",
  tren: "translates to English", translist: "translates: %s", dlgoing: "downloading:",
  trunavail: "unavailable — %s cannot translate", trlockok: "Got it", trconfirm: "Confirm",
  trlockmsg: "%s cannot be removed from the list — it is the default text output language. Pick another default language, and then %s can be excluded.",
  tronedlg: "Several languages are checked, but without asking the translation will always go into one — %s. The others stay checked but get disabled.",
  trmodeldlg: "%s cannot translate. If you continue, translation will be turned off and unavailable while this model works.",
  srcused: "in use", hfgo: "Search",
  remoteask: "Send audio to %s?", remotebadge: "REMOTE",
  helptocshow: "Show the contents — the window gets wider", helptochide: "Hide the contents and restore the window width",
  unitmin: "minutes", unithour: "hours", unitday: "days",
  depwhisper: "Runs Whisper models: a separate program next to the app that turns audio into text.", depllama: "post-processing, GGUF models", depsherpa: "engine for GigaAM and friends",
  depggml: "tensor library", deponnx: "runs the models inside sherpa-onnx", depwebview: "the settings window on WebView2",
  depwv2rt: "the Windows component that draws that window", depmalgo: "microphone capture", depminiaudio: "the audio layer inside malgo",
  depws: "talks to sherpa-server", depxsys: "WinAPI calls from Go", depwinloader: "DLL loading inside go-webview2",
  depplex: "the interface font", dephf: "where the models are downloaded from",
  srvdown: "Recognition is unavailable", srvdownwhy: "the remote recognition server is not set up — set its address in the settings", srvdowngo: "Open the server settings", retry: "Try again", srvwarnnow: "Dictation does not work right now: no address for the remote server.", srvwarnlater: "As soon as a Whisper model is picked, recognition will not work: no address for the remote server.",
  srvkthreads: "CPU threads", srvkport: "port", srvkauto: "autostart", srvkfile: "file", srvkaddr: "address", srvkcheck: "last check",
  srvnear: "whisper-server.exe next to the app", srvnoaddr: "not set", srvnocheck: "never checked",
  srvon: "on", srvoff: "off", srvlocaldlg: "Local recognition server", srvaddr: "Server address",
  srvaddrsub: "the address of whisper-server on the other machine, port included",
  threads: "CPU threads", threadssub: "more threads is not always faster",
  srvport: "Port", srvportsub: "the recognizer restarts itself",
  srvautostart: "Start whisper-server automatically", srvautostartsub: "turn off if you start the server yourself",
  serverexe: "whisper-server path", serverexesub: "filled in for you", profedit: "Edit",
  add: "Add profile", pname: "Name", pprompt: "Prompt", profadd: "Add",
  noprompts: "No prompts yet", pdrag: "drag to reorder", pnameph: "name it",
  ptestph: "type a phrase", pfnew: "New prompt", pfedit: "Editing the prompt",
  ptest: "Test", fitok: "fits", fitwarn: "tight",
  fitbad: "no RAM", ram: "RAM:", hfph: "model name", nollm: "no models",
  upd: "updated", pedit: "Edit", pclose: "Collapse",
  llmcatalog: "Model catalog",
  llmblock: "Installed models", llmnonehint: "No model is installed yet — download a found one with the arrow and it will appear here",
  llminmem: "in memory", llmondisk: "on disk", llmeject: "Unload from memory",
  llmfound: "found %d", llmnosearch: "no search yet", llmsearchhint: "Type a model name and press it", llmpickwait: "Available once the model is downloaded", llminstalled: "installed", llmsummodel: "model",
  llmsumsize: "size", llmsumcount: "installed", llmsumram: "memory",
  dlgclose: "Close", llmnopick: "not picked",
  confirmdel: 'Delete the "%s" model?',
  delactive: 'Delete the active model "%s"? Recognition stops until you pick another one.',
  wizneedmodel: "Download a model first",
  free: "free", updnone: "latest",
  ramavail: "Memory available: %s GB of %s GB", ramof: "%s GB of %s GB",
  badgemodels: "Installed models", badgemiss: "A model is not downloaded",
  badgesystem: "Warnings need attention", badgehist: "Entries in history",
  updavail: "Version %s available.", updgo: "Update", upderr: "Check failed",
  upddl: "Downloading", micdefault: "System default", micquiet: "quiet",
  more: "%d more settings", less: "Collapse %d settings",
  replempty: "No replacements yet", repldel: "Delete the replacement", replwhole: "whole words",
  replcase: "case", replfromph: "git hub", repltoph: "GitHub",
  repllang: "Rule language", repllangall: "all languages", listnothing: "Nothing found: “%s”", replwholefull: "Whole words only",
  replcasefull: "Match letter case", cmdaction: "Action",
  fmtrepladd: "Adding a replacement", fmtrepledit: "Editing the replacement",
  fmtcmdadd: "Adding a command", fmtcmdedit: "Editing the command",
  mtdel: "Deleting a model", mtdelprompt: "Deleting a prompt", mtdl: "Downloading a model",
  mttroff: "Turning translation off", mttrone: "Translating without asking", mttrlock: "Default output language",
  mtremote: "Remote server", mtpost: "External server", mthist: "Clearing the history",
  mtreset: "Resetting the settings", mtexe: "Server path",
  fmtdictadd: "Adding a word", dictempty: "No words yet",
  dictnomodel: "The current model %s does not support the dictionary — only Whisper models read it.",
  dictaddph: "a word, or several separated by commas",
  ovposscheme: "click the screen — the plate lands there",
  ovposdrag: "drag the plate with the mouse — it lands anywhere",
  fmadd: "Add", fmsave: "Save",
  tiprepllang: "fires only in the chosen language", tipreplcase: "case matters when on",
  tipreplwhole: "matches whole words only", tipcmdaction: "what saying the phrase does",
    stsummary: "Summary", stquick: "Quick settings", stmodels: "Models", stusage: "This week",
  stready: "Ready to dictate", stoff: "Switched off in the tray", stoffsub: "the shortcut does nothing until you switch it back on",
  stenable: "Switch on", stgoto: "Open this setting on its own tab", sthotkeygo: "Change the shortcut",
  stupdlast: "Version %s — the latest", stupdhave: "Version %s is available",
  stmem: "%s GB free of %s", stmemsub: "in memory: %s · on disk: %d models, %s GB", stmemnone: "nothing",
  stlang: "Speech language", stasr: "Recognition", ston: "on", stoffw: "off", stonf: "on", stofff: "off",
  stactive: "active", stidle: "not started", stdisk: "on disk, %s",
  stusagesub: "%d characters · %d today · %d characters on average", stnoweek: "no dictations this week",
  stautorunsub: "the app will not start by itself", stoverlaysub: "shown while recording",
  sthint: "hold the shortcut and speak", stcheck: "Check", strecog: "recognized by %s", stver: "Version %s", stlatest: "latest", navmic: "Microphone", overlay: "Plate on screen", beep: "Sound cue",
  autorun: "Start with Windows", postenable: "Enable post-processing", berropen: "Open", updcheck: "Check",
  histtill: "till %s", histtill1: "till tomorrow", histtillfull: "Will be deleted on %s — kept for %d days",
  histempty: "No history yet", histcopy: "Copy", histask: "Delete the whole dictation history?", histclear: "Clear",
  micchecking: "Checking…", mchecking: "Checking models…", histinsert: "Paste",
  retry: "Try again",
  acc: "accuracy", spd: "speed",
  recauto: "Detect itself",
  notforlang: "%s does not recognize this language", notinstalled: "not installed",
  manualnote: "Cannot be downloaded from the app — the licence forbids redistribution.",
  manuallink: "Download yourself", unload: "Unload", unloaded: "Unloaded",
  hffit: "only those that fit this computer", hfhidden: "hidden: %s",
  ovmoncursor: "The screen with the cursor", ovmonn: "Screen %d",
  postwarn: "Recognized text will go to this address.",
  postask: "Send recognized text to %s?",
  postkeyset: "key saved", postkeynone: "no key",
  apisumurl: "address", apisummodel: "model", apisumkey: "key", apisumtimeout: "timeout",
  apisumstate: "state", apinomodel: "not set", apinone: "not set up — post-processing runs locally",
  apisetup: "Set up", apiedit: "Change", apikeydel: "Delete the key", apidlg: "External server",
  apitest: "Test connection", apitestrun: "Checking…", apitestok: "The server answered",
  apiclear: "Clear", apiclearask: "Delete the address, the model and the key of the external server?",
  postnomodel: "on, but no model is picked", postnoprompt: "on, but no prompt is checked", postnoapi: "on, but the server is not set up",
  postbad: "the server did not answer: %s",
  postapiurl: "Address", postapimodel: "Model",
  postapikey: "API key", postapitimeout: "Response timeout", secshort: "s",
  skipadddlg: "Adding a program", skipeditdlg: "Editing the program", skipname: "Program name", skipnamesub: "The file name without a path",
  skipopen: "Programs open right now", skiprefresh: "Refresh the list", skippicked: "%d of %d picked", skipnone: "Nothing picked", skipempty: "The list is empty",
  cmdempty: "No commands yet", cmddel: "Delete the command", cmdph: "new line",
  cmdnewline: "line break", cmdparagraph: "new paragraph", cmdtext: "insert text", cmdcancel: "cancel the dictation",
  cmdtextph: "what to insert",
  cmdpnewline: "new line", cmdpparagraph: "new paragraph", cmdpcancel: "cancel", cmdpreset: "Add the usual ones",
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
  editor: { bg: "#1e1e1e", panel: "#252526", line: "#3c3c3c", text: "#d4d4d4", hi: "#4fc1ff", dim: "#9d9d9d", faint: "#6e6e6e", warn: "#cca700", bad: "#f14c4c", rec: "#f14c4c", field: "#3c3c3c", card: "#252526", soft: "#2d2d2d", navon: "#37373d", on: "#094771", btnbgh: "#1177bb", btn2bg: "#0e639c", btn2fg: "#ffffff", btn2line: "transparent", btn2bgh: "#1177bb", focus: "#007fd4", dangerbg: "#a1260d", dangerfg: "#ffffff", dangerbgh: "#c42b1c", dot: "#89d185" },
  neon: { bg: "#150a22", panel: "#1d0e30", line: "#4a2472", text: "#f3b6e4", hi: "#46e0ff", dim: "#b06ee0", faint: "#7d4fae", warn: "#ffd24a", bad: "#ff4d7d", rec: "#ff4d7d", field: "#1e0f33", card: "#231039", soft: "#2a1442", navon: "#2b1240", on: "#4a2472", btnbgh: "linear-gradient(90deg,rgba(255,95,200,.34),rgba(70,224,255,.28))", btn2bg: "linear-gradient(90deg,rgba(255,95,200,.20),rgba(70,224,255,.16))", btn2fg: "#ffffff", btn2line: "#4a2472", btn2bgh: "linear-gradient(90deg,rgba(255,95,200,.34),rgba(70,224,255,.28))", dangerbg: "linear-gradient(90deg,rgba(255,77,125,.30),rgba(255,77,125,.14))", dangerfg: "#ffd7e3", dangerbgh: "linear-gradient(90deg,rgba(255,77,125,.46),rgba(255,77,125,.24))", dot: "#46e0ff", focus: "#46e0ff" },
  soft: { bg: "#f3e2ea", panel: "#fbeef4", line: "#e9c3d8", text: "#a83a72", hi: "#c04a86", dim: "#c47ba3", faint: "#c9a2b8", warn: "#c07a1f", bad: "#c02a50", rec: "#e2597b", field: "#fffafc", card: "#fdf4f8", soft: "#ecd5e0", navon: "#f0d9e6", on: "#f6cfe2", btnbgh: "#c8467f", btn2bg: "#eccbdc", btn2fg: "#7d2b56", btn2line: "#dfb0c9", btn2bgh: "#e3bcd1", dangerbg: "#c02a50", dangerfg: "#ffffff", dangerbgh: "#a92345", dot: "#4f9d69", focus: "#c04a86" },
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
