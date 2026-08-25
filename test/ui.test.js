const fs = require("fs");
const path = require("path");
const { JSDOM } = require("jsdom");

const html = fs.readFileSync(path.join(__dirname, "page.html"), "utf8");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
function meterMoves(d, id) {
  const bars = [...d.querySelectorAll(`#${id} i`)];
  if (!bars.length) return false;
  return bars.some((b) => b.style.height && b.style.height !== "4px");
}
function searchFinds(w, d, needle) {
  const hits = w.searchMatches(needle.toLowerCase());
  return hits.some((el) => el.closest("#p-about"));
}

let llmState = {
  installed: [{ file: "model.gguf", size: 4929, active: true }],
  downloads: [],
  ram: 16384,
  ram_free: 9000,
};

const dom = new JSDOM(html, {
  runScripts: "dangerously",
  beforeParse(window) {
    window.appLLM = async () => JSON.stringify(llmState);
    let micBadge = "Realtek";
    let micFails = false;
    let ruState = "ready";
    let otherState = "ready";
    let remote = false;
    window.micSelectCalls = 0;
    window.setMicFails = (v) => { micFails = v; };
    window.setModelStates = (ru, other) => { ruState = ru; otherState = other; };
    window.setRemote = (v) => { remote = v; };
    let lastAt = 0;
    window.setLastAt = (v) => { lastAt = v; };
    let backendErr = "";
    window.setBackendErr = (v) => { backendErr = v; };
    window.retryCalls = 0;
    window.appRetryBackend = async () => { window.retryCalls++; backendErr = ""; };
    window.autorunCalls = [];
    window.wizardDone = 0;
    window.appAutorun = async () => window.autorunCalls.length > 0 && window.autorunCalls[window.autorunCalls.length - 1];
    window.appSetAutorun = async (on) => { window.autorunCalls.push(on); return on; };
    window.appWizardDone = async () => { window.wizardDone++; };
    window.appState = async () =>
      JSON.stringify({ hotkey: "ctrl+win", mic: "Realtek", engine: "sherpa · gigaam-v3", llm: "model.gguf",
        ram: "8000 MB free", last: "hello", last_meta: "just now · 5 characters", ready: true, status: "Ready",
        status_line: "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free", ru_model: "gigaam-v3",
        other_model: "ggml-small.bin", llm_ok: true, mic_ok: true, last_at: lastAt, last_app: "chrome.exe",
        remote: remote, backend_err: backendErr,
        ru_state: ruState, other_state: otherState,
        badges: { mic: micBadge, models: "2", system: "" } });
    window.appRouting = async () =>
      JSON.stringify([
        { cond: "Speech in RU", engine: "gigaam-v3", why: "more accurate here" },
        { cond: "Other languages", engine: "ggml-small.bin", why: "99 languages" },
        { cond: "Translation", engine: "ggml-small.bin", why: "only Whisper translates" },
      ]);
    let modelStates = { base: "absent", small: "active", "gigaam-v3": "absent" };
    window.dlCalls = [];
    window.cancelCalls = [];
    window.appModelDl = async (id) => { window.dlCalls.push(id); modelStates[id] = "downloading"; };
    window.finishDl = (id) => { modelStates[id] = "installed"; };
    window.appModelCancel = async (id) => { window.cancelCalls.push(id); modelStates[id] = "absent"; return true; };
    window.appModels = async () =>
      JSON.stringify([
        { id: "base", name: "Base", desc: "fast", size: 142, state: modelStates.base, pct: 12, engine: "whisper", langs: "*" },
        { id: "small", name: "Small", desc: "balanced", size: 466, state: modelStates.small, engine: "whisper", langs: "*", slot: true },
        { id: "gigaam-v3", name: "GigaAM v3", desc: "russian", size: 232, state: modelStates["gigaam-v3"], pct: 5, engine: "sherpa", langs: "ru", punct: true },
      ]);
    window.appLLMSearch = async () =>
      JSON.stringify({ repos: [{ id: "org/Repo-GGUF", downloads: 1234, updated: "2026-01-01" }] });
    window.appLLMFiles = async () =>
      JSON.stringify({ files: [{ file: "q4.gguf", size: 4000, fit: "ok", need: 6166 }] });
    window.appLLMDel = async () => {
      llmState = { installed: [], downloads: [], ram: 16384, ram_free: 9000 };
      return "deleted";
    };
    window.appAdvise = async () => JSON.stringify({
      primary: "gigaam-v3", companion: "small", text: "I recommend GigaAM v3.", ram: "8000 MB free",
      need: 232,
      plan: [
        { id: "gigaam-v3", name: "GigaAM v3", size: 232, installed: false },
        { id: "small", name: "Small", size: 466, installed: true },
      ],
    });
    for (const name of [
      "appLLMDlFile", "appLLMTest", "appHFPage", "appHFHome", "appRepoLink",
      "appAuthorLink", "appCapture", "appCaptureCombo", "appReload",
      "appPreviewSound", "appMin", "appClose",
      "appDoUpdate", "appReady", "appJSError",
    ]) {
      window[name] = () => {};
    }
    window.dragCalls = 0;
    window.appDrag = () => { window.dragCalls++; };
    window.resizeCalls = 0;
    window.maximized = false;
    window.resizeEdges = [];
    window.appResizeEdge = (edge) => { window.resizeCalls++; window.resizeEdges.push(edge); };
    window.appMaxRestore = async () => { window.maximized = !window.maximized; return window.maximized; };
    window.appMaximized = async () => window.maximized;
    window.saveCalls = 0;
    window.saveForms = [];
    window.lastSave = {};
    window.appSave = async (json) => {
      window.saveCalls++;
      const f = JSON.parse(json);
      micBadge = f.mic_device_name ? f.mic_device_name.split(" ")[0] : "Realtek";
      const message = Number(f.server_port) === 8910 ? "Saved" : "Restarting the recognizer…";
      window.lastSaveForm = f;
      window.saveForms.push(f);
      window.lastSave = { ok: true, severity: "ok", message };
      return JSON.stringify(window.lastSave);
    };
    window.replaceCalls = [];
    window.appTestText = async (t) => { window.replaceCalls.push(t); const out = t.replace(/git hub/gi, "GitHub"); return JSON.stringify({ text: out, cancelled: false }); };
    let histItems = [{ at: 1700000000000, text: "выложи на GitHub", app: "chrome.exe" }, { at: 1699999000000, text: "привет команде", app: "Telegram.exe" }];
    window.histQueries = [];
    window.histCopied = [];
    window.appHistory = async (q) => { window.histQueries.push(q); return JSON.stringify(histItems.filter(i => !q || i.text.includes(q) || i.app.includes(q))); };
    window.appHistoryClear = async () => { histItems = []; };
    window.copyFails = false;
    window.setCopyFails = (v) => { window.copyFails = v; };
    window.appHistoryCopy = async (at) => { window.histCopied.push(at);
      return JSON.stringify(window.copyFails ? { ok: false, text: "Could not copy: clipboard busy" } : { ok: true, text: "Copied" }); };
    window.lastCopied = 0;
    window.appCopyLast = async () => { window.lastCopied++;
      return JSON.stringify(window.copyFails ? { ok: false, text: "Could not copy: clipboard busy" } : { ok: true, text: "Copied" }); };
    window.histInserted = [];
    window.appHistoryInsert = async (at) => { window.histInserted.push(at); return JSON.stringify({ ok: true, text: "pasted into “Editor”" }); };
    window.modelChecks = 0;
    window.appCheckModels = async () => { window.modelChecks++; return JSON.stringify({ ok: false, text: "Damaged files: Small (ggml-small.bin)", rows: [] }); };
    window.listsExported = [];
    window.appListsExport = async (payload) => { window.listsExported.push(JSON.parse(payload)); return JSON.stringify({ ok: true, text: "saved to lists.json" }); };
    window.listsImported = [];
    window.appListsImport = async (payload) => {
      window.listsImported.push(JSON.parse(payload));
      return JSON.stringify({ ok: true, text: "added: 2, skipped: 0",
        replacements: [{ id: "r1", from: "git hub", to: "GitHub", whole: true, match_case: false }],
        commands: [{ id: "c1", phrase: "новая строка", action: "newline", text: "" }] });
    };
    window.micChecks = 0;
    window.appMicCheck = async () => { window.micChecks++; return JSON.stringify({ verdict: "quiet", text: "Too quiet: peak -32 dB", peak_db: -32, voice: 0.8, clip: 0 }); };
    window.delCalls = [];
    window.appModelDel = async (id, force) => { window.delCalls.push([id, force]); if(force) modelStates[id] = "absent"; return "ok"; };
    window.appMics = async () =>
      JSON.stringify([
        { id: "dev1", name: "Headset (USB)", default: false },
        { id: "dev2", name: "Webcam microphone", default: false },
      ]);
    window.appMicLevel = async () => 0.42;
    window.appMicSelect = async () => {
      window.micSelectCalls++;
      if(micFails) return JSON.stringify({ ok: false, severity: "error", message: "Microphone busy" });
      return JSON.stringify({ ok: true, severity: "ok" });
    };
    window.appUpdateStatus = async () => JSON.stringify({ current: "0.0.0", latest: "", url: "" });
    window.appCheckUpdate = async () => JSON.stringify({ current: "0.0.0", latest: "v0.0.0", newer: false });
  },
});

const w = dom.window;
const d = w.document;
const failures = [];

function check(name, actual, expected) {
  const ok = JSON.stringify(actual) === JSON.stringify(expected);
  console.log(`${ok ? "PASS" : "FAIL"}  ${name}${ok ? "" : ` (got ${JSON.stringify(actual)}, want ${JSON.stringify(expected)})`}`);
  if (!ok) failures.push(name);
}

(async () => {
  const errors = [];
  w.addEventListener("error", (e) => errors.push(e.message));
  await sleep(500);

  check("page script executed", d.getElementById("ver").textContent, "0.0.0-test");

  const tab = (p) => d.querySelector(`.nav[data-p="${p}"]`).click();
  const shown = (p) => d.getElementById(`p-${p}`).classList.contains("active");

  check("opens on the status screen", shown("state"), true);
  check("status hotkey shown", d.getElementById("state_hotkey").textContent, "ctrl+win");
  check("status shows russian model", d.getElementById("state_ru").textContent, "gigaam-v3");
  check("status shows other-language model", d.getElementById("state_other").textContent, "ggml-small.bin");
  check("last dictation carries details", d.getElementById("state_last_meta").textContent, "just now · 5 characters");
  const lastCS = w.getComputedStyle(d.getElementById("state_last"));
  check("last dictation is clamped to its row", [lastCS.display, lastCS.overflow, lastCS.textOverflow, lastCS.whiteSpace], ["block", "hidden", "ellipsis", "nowrap"]);
  check("full dictation text kept on hover", d.getElementById("state_last").title, "hello");
  check("status screen meter follows the microphone", meterMoves(d, "state_mic_bar"), true);
  check("the status meter is drawn as bars", d.querySelectorAll("#state_mic_bar i").length, 7);
  const micCard = d.getElementById("state_mic").closest(".scard");
  check("a long device name is kept to one line", w.getComputedStyle(d.getElementById("state_mic").closest(".v")).textOverflow, "ellipsis");
  check("and the whole name waits under the pointer", micCard.dataset.tip, "Realtek");
  w.showTip(micCard);
  check("the hint shows what it was given", d.getElementById("tip").textContent, "Realtek");
  check("and it is dressed like the rest of the window", d.getElementById("tip").classList.contains("on"), true);
  w.hideTip();
  check("it goes away when the pointer leaves", d.getElementById("tip").classList.contains("on"), false);
  for (let i = 0; i < 8; i++) w.paintMeter(d.getElementById("state_mic_bar"), 0.005);
  check("a quiet room leaves the meter flat", [...d.querySelectorAll("#state_mic_bar i")].every(b=>b.style.height === "4px"), true);
  w.paintMeter(d.getElementById("state_mic_bar"), 0.6);
  check("a loud phrase raises it", [...d.querySelectorAll("#state_mic_bar i")].some(b=>parseInt(b.style.height) > 8), true);
  for (const edge of ["t", "b", "l", "r", "tl", "tr", "bl", "br"]) {
    d.querySelector(".rsz." + edge).dispatchEvent(new w.MouseEvent("mousedown", { button: 0, bubbles: true }));
  }
  check("every edge and corner asks the window to resize", w.resizeCalls, 8);
  check("and each asks for its own direction", w.resizeEdges, [12, 15, 10, 11, 13, 14, 16, 17]);
  d.getElementById("cap_max").click();
  await sleep(150);
  check("the button fills the screen", w.maximized, true);
  check("and turns into a restore button", d.getElementById("cap_max").title, "S_WND_RESTORE");
  d.getElementById("cap_max").click();
  await sleep(150);
  check("pressing it again gives the size back", w.maximized, false);
  check("status bar names both models", d.getElementById("st_main").textContent, "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free");
  check("sidebar badges filled", [d.getElementById("badge_mic").textContent, d.getElementById("badge_models").textContent], ["Realtek", "2"]);
  check("status bar led lit", d.getElementById("st_led").classList.contains("on"), true);
  check("the post-processing card offers what the state calls for", d.getElementById("state_llm_btn").textContent, "Change");
  check("copying the last dictation is offered", d.getElementById("state_copy").disabled, false);
  check("the models badge explains itself", d.getElementById("badge_models").title, "Installed models");
  check("icon buttons carry a name for screen readers", d.getElementById("mic_refresh").getAttribute("aria-label"), "S_MIC_REFRESH");
  check("and the same words are the hover hint", d.getElementById("mic_refresh").dataset.tip, "S_MIC_REFRESH");
  check("the browser's own tooltip is out of the way", d.getElementById("mic_refresh").hasAttribute("title"), false);
  check("the open section is announced", d.querySelector(".nav[data-p=state]").getAttribute("aria-selected"), "true");

  tab("system"); await sleep(200);
  const themeSel = d.getElementById("theme");
  const skinSel = d.getElementById("skin");
  check("design and colour are two separate choices", [!!skinSel, !!themeSel], [true, true]);
  check("three designs are offered", [...skinSel.options].map(o=>o.value), ["terminal", "editor", "neon"]);
  check("the colours are the terminal's four", [...themeSel.options].map(o=>o.value), ["green", "amber", "blue", "pink"]);
  check("the terminal design is the one in use", [skinSel.value, themeSel.value], ["terminal", "green"]);
  check("every colour has a swatch", d.querySelectorAll("#theme_swatches .swatch").length, 4);
  check("the swatch of the current colour is marked", d.querySelector("#theme_swatches .swatch.on").dataset.theme, "green");

  d.querySelector('#theme_swatches .swatch[data-theme="amber"]').click();
  await sleep(200);
  check("picking a swatch picks the colour", themeSel.value, "amber");
  check("the window repaints at once", d.documentElement.style.getPropertyValue("--green"), "#ff9e2c");
  check("and the warning colour comes with it", d.documentElement.style.getPropertyValue("--amber"), "#ffd24a");
  check("the design stays the terminal one", d.documentElement.style.getPropertyValue("--font").includes("Consolas"), true);

  skinSel.value = "editor";
  skinSel.dispatchEvent(new w.Event("change", { bubbles: true }));
  await sleep(200);
  check("a design brings its own font", d.documentElement.style.getPropertyValue("--font").includes("Cascadia"), true);
  check("and its own corners", d.documentElement.style.getPropertyValue("--r"), "3px");
  check("and turns the halo off", d.documentElement.style.getPropertyValue("--glow"), "none");
  check("and its own colours, not the picked one", d.documentElement.style.getPropertyValue("--green"), "#d4d4d4");
  check("with its accent kept apart from its text", d.documentElement.style.getPropertyValue("--hi"), "#4fc1ff");
  check("and no capitals shouted at the buttons", d.documentElement.style.getPropertyValue("--caps"), "none");
  check("controls take their height from the design", d.documentElement.style.getPropertyValue("--ctlpad"), "5px 9px");
  check("down to the surfaces you type into", d.documentElement.style.getPropertyValue("--field"), "#3c3c3c");
  check("so the colour row steps aside", d.getElementById("colour_row").style.display, "none");

  skinSel.value = "terminal";
  skinSel.dispatchEvent(new w.Event("change", { bubbles: true }));
  await sleep(200);
  check("back on the terminal the colour row returns", d.getElementById("colour_row").style.display, "");
  check("with the colour that was picked", d.documentElement.style.getPropertyValue("--green"), "#ff9e2c");
  check("square corners", d.documentElement.style.getPropertyValue("--r"), "0px");
  check("and its own window border", d.documentElement.style.getPropertyValue("--wborder"), "1px solid #4a3018");
  check("the surfaces follow the colour too", d.documentElement.style.getPropertyValue("--field"), "#120c07");
  await sleep(700);
  check("both choices are written into the settings", w.saveForms.some(f=>f.theme === "amber" && f.skin === "terminal"), true);
  d.querySelector('#theme_swatches .swatch[data-theme="green"]').click();
  await sleep(700);
  check("the closed sections are not", d.querySelector(".nav[data-p=models]").getAttribute("aria-selected"), "false");
  check("each screen is a panel with a name", [d.getElementById("p-state").getAttribute("role"), d.getElementById("p-state").getAttribute("aria-label")], ["tabpanel", "S_NAV_STATE"]);
  w.setHotkey("win+l", "win+l is taken by Windows: lock the computer. Dictation may never start");
  check("a shortcut Windows already uses is called out", d.getElementById("st_saved").textContent, "win+l is taken by Windows: lock the computer. Dictation may never start");
  check("and it is marked as a warning", d.getElementById("st_saved").className, "stsaved warn");
  w.setHotkey("ctrl+win", "");

  check("nine sections in the sidebar", d.querySelectorAll(".nav").length, 9);

  tab("about"); await sleep(200);
  const toc = d.querySelector("#p-about .toc");
  const helpCard = toc && toc.closest(".card");
  check("the help opens with a table of contents", !!toc, true);
  check("every help heading is listed", toc ? toc.querySelectorAll("a").length : 0, helpCard ? helpCard.querySelectorAll(".wh").length : -1);
  check("the contents link to the headings", toc ? toc.querySelector("a").getAttribute("href") : "", "#" + (helpCard ? helpCard.querySelector(".wh").id : ""));
  check("the help is reachable from search", searchFinds(w, d, "config.json"), true);

  tab("history"); await sleep(300);
  check("history is a screen of its own", shown("history"), true);
  check("history is off until asked for", d.getElementById("history").checked, false);
  check("the programs to skip are never folded away", !!d.querySelector("#hist_skip_row[data-adv]"), false);
  check("the dictations are listed", d.querySelectorAll("#histbody .histrow").length, 2);
  check("each entry names the program", d.querySelector("#histbody .histmeta").textContent.includes("chrome.exe"), true);
  check("each entry carries the text", d.querySelector("#histbody .histtext").textContent, "выложи на GitHub");
  d.querySelector("#histbody .mini").click(); await sleep(200);
  check("an entry can be copied", w.histCopied, [1700000000000]);
  check("copying says it copied", d.getElementById("st_saved").textContent, "Copied");
  w.setCopyFails(true);
  d.querySelectorAll("#histbody .histrow")[0].querySelectorAll("button")[0].click(); await sleep(250);
  check("a failed copy says so instead of success", d.getElementById("st_saved").textContent, "Could not copy: clipboard busy");
  check("and it is marked as an error", d.getElementById("st_saved").className.includes("bad"), true);
  w.setCopyFails(false);
  d.querySelectorAll("#histbody .histrow")[0].querySelectorAll("button")[1].click(); await sleep(250);
  check("an entry can be pasted back", w.histInserted, [1700000000000]);
  check("the program says where it went", d.getElementById("st_saved").textContent, "pasted into “Editor”");
  const hfind = d.getElementById("hist_find");
  hfind.value = "Telegram"; hfind.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(400);
  check("search reaches the program itself", w.histQueries[w.histQueries.length - 1], "Telegram");
  check("search narrows the list", d.querySelectorAll("#histbody .histrow").length, 1);
  hfind.value = ""; hfind.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(400);
  d.getElementById("hist_clear").click(); await sleep(150);
  check("clearing the history asks first", !!d.querySelector(".modal-bg"), true);
  d.querySelector(".modal .btn.yes").click(); await sleep(400);
  check("after clearing the list says it is empty", d.querySelector("#histbody .histempty").textContent, "No history yet");
  const everySetting = [
    "hotkey", "min_record_ms", "max_record_seconds", "auto_enter", "restore_clipboard",
    "type_mode", "overlay", "overlay_position", "overlay_text", "animation", "mic_device", "mic_bar", "beep", "sound_theme",
    "routing", "models", "language", "threads", "punctuation", "whisper_prompt", "profbody",
    "tr_default", "translate_target", "translate_ask", "translate_ask_seconds", "tr_hotkey",
    "tl_en", "ui_language", "upd_check", "check_updates", "server_autostart", "server_port",
    "server_exe", "server_url", "proc-models", "proc-search", "ver2", "autorun", "pause_hotkey",
  ];
  const missing = everySetting.filter((id) => !d.getElementById(id));
  check("every setting is present in the new window", missing, []);

  tab("models"); await sleep(80);
  check("models section shown", shown("models"), true);
  check("recognition models listed", d.querySelectorAll('#p-models input[name^="mdl-"]').length, 3);
  check("the list is split into two slots", [...d.querySelectorAll("#models .mslot")].map(h=>h.dataset.slot), ["ru", "other"]);
  check("the russian slot holds the russian engine", d.querySelectorAll('#models .mrow[data-slot="ru"]').length, 1);
  check("every other language shares the second slot", d.querySelectorAll('#models .mrow[data-slot="other"]').length, 2);
  check("model filters rendered", d.querySelectorAll(".fchip").length, 5);
  const recLangs = [...d.getElementById("language").options].map(o=>o.value);
  check("italian can be dictated too", recLangs.includes("it"), true);
  check("ram estimate shown", d.querySelectorAll("#p-models .mram").length, 3);
  check("routing panel rows", d.querySelectorAll("#routing .rrow").length, 3);
  check("routing shows engine", d.querySelectorAll("#routing .reng")[0].textContent, "gigaam-v3");
  check("engine tags rendered", d.querySelectorAll("#p-models .mtag").length, 3);
  check("russian engine tagged RU", d.querySelector('#models .mrow[data-slot="ru"] .mtag').textContent, "RU");
  d.getElementById("mcheck").click(); await sleep(250);
  check("installed models can be checked", w.modelChecks, 1);
  check("the check says what it found", d.getElementById("mcheck_out").textContent, "Damaged files: Small (ggml-small.bin)");
  check("a damaged model is marked as bad", d.getElementById("mcheck_out").className.includes("bad"), true);

  tab("mic"); await sleep(120);
  const mic = d.getElementById("mic_device");
  check("microphone list has default plus devices", mic.options.length, 3);
  check("default option is localized", mic.options[0].textContent, "System default");
  check("system default selected initially", mic.value, "");
  mic.value = "dev1"; mic.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(30);
  check("microphone selection kept", mic.value, "dev1");
  await sleep(260);
  check("input level meter moves", meterMoves(d, "mic_bar"), true);
  check("sidebar badge follows the chosen microphone", d.getElementById("badge_mic").textContent, "Headset");

  const micBtn = d.getElementById("mic_check");
  micBtn.click(); await sleep(300);
  check("the microphone can be checked on demand", w.micChecks, 1);
  check("the verdict is shown in words", d.getElementById("mic_verdict").textContent, "Too quiet: peak -32 dB");
  check("a bad verdict is marked as such", d.getElementById("mic_verdict").className, "micverdict bad");

  mic.value = ""; mic.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(30);
  tab("text"); await sleep(30);
  check("dictionary textarea present", !!d.getElementById("whisper_prompt"), true);
  check("punctuation modes offered", d.getElementById("punctuation").options.length, 3);

  tab("dictation"); await sleep(60);
  const ovpos = d.getElementById("overlay_position");
  check("the plate can be put in three places", ovpos.options.length, 3);
  check("it starts at the bottom of the screen", ovpos.value, "bottom");
  check("the pause hotkey is shown", d.getElementById("pause_hotkey").textContent, "ctrl+alt+p");
  d.getElementById("pause_clear").click(); await sleep(300);
  check("clearing the pause hotkey is saved", w.lastSaveForm.pause_hotkey, "");
  ovpos.value = "caret"; ovpos.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("moving the plate is saved", w.lastSaveForm.overlay_position, "caret");
  ovpos.value = "bottom"; ovpos.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  const ovtext = d.getElementById("overlay_text");
  check("showing the text on the plate is a switch", ovtext.type, "checkbox");
  ovtext.checked = false; ovtext.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("turning it off is saved", w.lastSaveForm.overlay_text, false);
  ovtext.checked = true; ovtext.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);

  tab("dictation"); await sleep(60);
  check("with no rules the list says so", d.querySelector("#rulesbody .ruleempty").textContent, "No rules yet");
  check("the last program is offered as a rule", d.getElementById("rule_last").textContent, "last insertion: chrome.exe");
  d.getElementById("rule_last").click(); await sleep(300);
  check("one click makes a rule for it", d.querySelectorAll("#rulesbody .rulerow").length, 1);
  check("the rule carries the program", w.lastSaveForm.app_rules[0].match, "chrome.exe");
  const rrow = d.querySelector("#rulesbody .rulerow");
  const rpaste = rrow.querySelector(".rpaste");
  check("a rule inherits until told otherwise", [rpaste.value, rpaste.options[0].textContent], ["", "insertion: as set"]);
  rpaste.value = "type"; rpaste.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("the insertion method is saved", w.lastSaveForm.app_rules[0].paste_mode, "type");
  const renter = rrow.querySelector(".renter");
  renter.value = "off"; renter.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("Enter can be turned off for one program", w.lastSaveForm.app_rules[0].auto_enter, "off");
  const rdelay = rrow.querySelector(".rdelay");
  rdelay.value = "250"; rdelay.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("the delay is saved as a number", w.lastSaveForm.app_rules[0].delay_ms, 250);
  const rprof = rrow.querySelector(".rprof");
  check("prompts of the rule offer every profile", rprof.options.length, 4);
  rprof.value = "-"; rprof.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("a program can be left without prompts", [w.lastSaveForm.app_rules[0].use_profiles, w.lastSaveForm.app_rules[0].profiles], [true, []]);
  rprof.value = "formal"; rprof.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("or given its own prompt", w.lastSaveForm.app_rules[0].profiles, ["formal"]);
  d.getElementById("rule_add").click(); await sleep(100);
  check("more rules can be added", d.querySelectorAll("#rulesbody .rulerow").length, 2);
  d.querySelectorAll("#rulesbody .rdel")[1].click(); await sleep(300);
  check("and deleted", d.querySelectorAll("#rulesbody .rulerow").length, 1);
  check("deleting leaves the other rule alone", w.lastSaveForm.app_rules.length, 1);

  tab("system"); await sleep(30);
  check("service settings shown", [shown("system"), !!d.getElementById("server_url")], [true, true]);
  const su = d.getElementById("server_url");
  const savesBeforeURL = w.saveCalls;
  su.value = "http://192.168.0.5:9000"; su.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(350);
  check("the remote address asks before it is applied", !!d.querySelector(".modal-bg"), true);
  check("and nothing is saved while the question is open", w.saveCalls, savesBeforeURL);
  d.querySelector(".modal .btn.ghost").click(); await sleep(350);
  check("saying no puts the field back", su.value, "");
  check("saying no saves nothing", w.saveCalls, savesBeforeURL);
  su.value = "http://192.168.0.5:9000"; su.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(350);
  d.querySelector(".modal .btn.yes").click(); await sleep(400);
  check("saying yes applies the address", w.lastSaveForm.server_url, "http://192.168.0.5:9000");
  const beepBox = d.getElementById("beep");
  su.value = "http://sneaky.example";
  beepBox.checked = !beepBox.checked; beepBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(350);
  check("an unconfirmed address never rides along with another setting", w.lastSaveForm.server_url, "http://192.168.0.5:9000");
  beepBox.checked = !beepBox.checked; beepBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(350);
  su.value = ""; su.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(350);
  check("clearing the address needs no question", [!!d.querySelector(".modal-bg"), w.lastSaveForm.server_url], [false, ""]);
  const autorun = d.getElementById("autorun");
  const autorunBefore = w.autorunCalls.length;
  autorun.checked = true; autorun.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("starting with Windows is a switch of its own", w.autorunCalls.length, autorunBefore + 1);
  check("it is not written into the config", w.autorunCalls[w.autorunCalls.length - 1], true);
  autorun.checked = false; autorun.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("and it can be turned back off", w.autorunCalls[w.autorunCalls.length - 1], false);

  tab("translate"); await sleep(30);
  check("the turbo warning is shown for a turbo model", d.getElementById("tr_warn").style.display, "block");
  check("the target language carries the honest note", !!d.querySelector("#p-translate .row label .sub.warn"), true);
  const trd = d.getElementById("tr_default");
  const ask = d.getElementById("translate_ask");
  const state = () => [
    d.getElementById("translate_target").disabled,
    ask.disabled,
    d.getElementById("translate_ask_seconds").disabled,
    d.getElementById("tl_en").disabled,
  ];
  trd.checked = true; trd.dispatchEvent(new w.Event("change"));
  check("always-translate disables ask controls", state(), [false, true, true, true]);
  trd.checked = false; trd.dispatchEvent(new w.Event("change"));
  ask.value = "never"; ask.dispatchEvent(new w.Event("change"));
  check("never mode disables target and dialog langs", state(), [true, false, true, true]);
  ask.value = "always"; ask.dispatchEvent(new w.Event("change"));
  check("always-ask enables dialog langs", state(), [true, false, true, false]);
  ask.value = "timeout"; ask.dispatchEvent(new w.Event("change"));
  check("timeout mode keeps target editable", state(), [false, false, false, false]);

  tab("text"); await sleep(80);
  check("with no replacements the list says so", d.querySelector("#replbody .ruleempty").textContent, "No replacements yet");
  d.getElementById("repl_add").click(); await sleep(100);
  const rep = d.querySelector("#replbody .replrow");
  check("a replacement row has both sides", [!!rep.querySelector(".rfrom"), !!rep.querySelector(".rto")], [true, true]);
  check("whole words is on by default", rep.querySelector(".rwhole").checked, true);
  check("case is off by default", rep.querySelector(".rcase").checked, false);
  const rfrom = rep.querySelector(".rfrom"), rto = rep.querySelector(".rto");
  rfrom.value = "гит хаб"; rfrom.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(250);
  rto.value = "GitHub"; rto.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(250);
  check("the replacement is saved", [w.lastSaveForm.replacements[0].from, w.lastSaveForm.replacements[0].to], ["гит хаб", "GitHub"]);
  check("its flags are saved too", [w.lastSaveForm.replacements[0].whole, w.lastSaveForm.replacements[0].match_case], [true, false]);
  const rtest = d.getElementById("repl_test");
  rtest.value = "push to git hub"; rtest.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(500);
  check("the test field asks the program itself", w.replaceCalls[w.replaceCalls.length - 1], "push to git hub");
  check("and shows what would come out", d.getElementById("repl_out").textContent, "push to GitHub");
  d.querySelector("#replbody .rdel").click(); await sleep(250);
  check("a replacement can be deleted", w.lastSaveForm.replacements.length, 0);

  check("with no commands the list says so", d.querySelector("#cmdbody .ruleempty").textContent, "No commands yet");
  d.getElementById("cmd_preset").click(); await sleep(300);
  check("the usual commands can be added in one click", d.querySelectorAll("#cmdbody .replrow").length, 3);
  check("they carry the phrases of this language", w.lastSaveForm.commands.map((c) => c.phrase), ["new line", "new paragraph", "cancel"]);
  check("and the actions that go with them", w.lastSaveForm.commands.map((c) => c.action), ["newline", "paragraph", "cancel"]);
  d.getElementById("cmd_preset").click(); await sleep(300);
  check("adding them twice changes nothing", d.querySelectorAll("#cmdbody .replrow").length, 3);
  const cmdRow = d.querySelector("#cmdbody .replrow");
  check("a command with no text keeps its field hidden", cmdRow.querySelector(".ctext").style.display, "none");
  const cact = cmdRow.querySelector(".caction");
  cact.value = "text"; cact.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("choosing insert-text reveals the field", cmdRow.querySelector(".ctext").style.display, "");
  check("the action is saved", w.lastSaveForm.commands[0].action, "text");
  d.querySelectorAll("#cmdbody .rdel")[0].click(); await sleep(300);
  check("a command can be deleted", w.lastSaveForm.commands.length, 2);

  d.getElementById("lists_export").click(); await sleep(250);
  check("the lists can be saved to a file", w.listsExported.length, 1);
  check("what is on the page goes into the file", Array.isArray(w.listsExported[0].commands), true);
  check("saving reports back", d.getElementById("st_saved").textContent, "saved to lists.json");
  d.getElementById("lists_import").click(); await sleep(300);
  check("the lists can be loaded from a file", w.listsImported.length, 1);
  check("loading brings the replacements in", w.lastSaveForm.replacements.map((r) => r.from), ["git hub"]);
  check("and the commands with them", w.lastSaveForm.commands.map((c) => c.phrase), ["новая строка"]);
  check("the rows appear on the page", d.querySelectorAll("#replbody .replrow").length, 1);

  const pb = d.getElementById("profbody");
  check("prompts listed", pb.querySelectorAll("input.profcb").length, 2);

  const pencil = () => pb.querySelector('button[data-a="edit"]');
  pencil().click(); await sleep(60);
  check("prompt editor opens", !!d.getElementById("pf_name"), true);
  pencil().click(); await sleep(60);
  check("prompt editor closes on second click", !!d.getElementById("pf_name"), false);
  pencil().click(); await sleep(60);
  d.getElementById("pf_close").click(); await sleep(60);
  check("prompt editor closes via collapse button", !!d.getElementById("pf_name"), false);

  const pdel = () => pb.querySelector('button[data-a="pdel"]');
  const profsBefore = w.lastSaveForm.profiles ? w.lastSaveForm.profiles.length : 2;
  pdel().click(); await sleep(200);
  check("deleting a prompt asks first", !!d.querySelector(".modal-bg"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no keeps the prompt", pb.querySelectorAll("input.profcb").length, 2);
  pdel().click(); await sleep(200);
  d.querySelector(".modal .btn.yes").click(); await sleep(400);
  check("saying yes removes it and saves at once", w.lastSaveForm.profiles.length, profsBefore - 1);

  tab("models"); await sleep(30);
  const del = d.querySelector('#proc-models button[data-a="ldel"]');
  check("active LLM model can be deleted", !!del, true);
  del.click(); await sleep(60);
  check("deleting a model asks first", !!d.querySelector(".modal-bg"), true);
  check("the question is asked in the app style, not by the browser", d.querySelectorAll(".modal .btn").length, 2);
  check("the way out comes first, the action second", [...d.querySelectorAll(".modal .btn")].map(b=>b.className), ["btn ghost", "btn yes"]);
  check("the question keeps the focus on the way out", d.activeElement.className, "btn ghost");
  check("the question is a dialog for the reader", d.querySelector(".modal").getAttribute("aria-modal"), "true");
  check("the window behind is out of reach while it is open", !!d.querySelector(".content[inert]"), true);
  d.dispatchEvent(new w.KeyboardEvent("keydown", { key: "Tab", bubbles: true }));
  check("Tab moves between the two answers and no further", d.activeElement.className, "btn yes");
  d.dispatchEvent(new w.KeyboardEvent("keydown", { key: "Escape", bubbles: true }));
  await sleep(150);
  check("Escape closes it as a no", !!d.querySelector(".modal-bg"), false);
  check("and the window behind comes back", !!d.querySelector(".content[inert]"), false);
  del.click(); await sleep(150);
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("the question closes with the answer", !!d.querySelector(".modal-bg"), false);
  check("model list empty after delete", d.querySelectorAll('#proc-models input[name="llmmdl"]').length, 0);

  d.getElementById("hf_q").value = "qwen";
  d.getElementById("hf_go").click(); await sleep(300);
  check("searching Hugging Face is a button, not a decoration", d.getElementById("hf_go").tagName, "BUTTON");
  const repo = d.querySelector(".hfrepo");
  check("a found repository is a button too", !!repo, true);
  check("and it says whether it is open", repo.getAttribute("aria-expanded"), "false");
  repo.click(); await sleep(300);
  check("opening it lists the files", d.querySelectorAll('#hf_results button[data-repo]').length, 1);
  check("and the button now says it is open", d.querySelector(".hfrepo").getAttribute("aria-expanded"), "true");


  tab("models"); await sleep(120);
  const pickAbsent = () => d.querySelector('#models input[value="base"]');
  pickAbsent().click(); await sleep(80);
  check("picking a model that is not here asks first", !!d.querySelector(".modal-bg"), true);
  check("the question names the model and its size", d.querySelector(".modal p").textContent.includes("Base") && d.querySelector(".modal p").textContent.includes("142 MB"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no downloads nothing", w.dlCalls.length, 0);
  check("saying no puts the choice back", d.querySelector('#models input[value="small"]').checked, true);

  pickAbsent().click(); await sleep(80);
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("saying yes starts the download", w.dlCalls, ["base"]);
  check("a download can be stopped", !!d.querySelector('#models button[data-a="cancel"][data-id="base"]'), true);
  d.querySelector('#models button[data-a="cancel"][data-id="base"]').click(); await sleep(250);
  check("stopping asks the program to stop", w.cancelCalls, ["base"]);
  check("a stopped download offers to start again", !!d.querySelector('#models button[data-a="dl"][data-id="base"]'), true);

  pickAbsent().click(); await sleep(80);
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("the row shows the download running", !!d.querySelector('#models button[data-a="cancel"][data-id="base"]'), true);
  const savesBeforeDl = w.saveForms.length;
  w.finishDl("base"); await sleep(1400);
  check("a finished download is applied by itself", w.saveForms.slice(savesBeforeDl).map((f) => f.model_id), ["base"]);
  check("and the program says the model is ready", d.getElementById("st_saved").textContent, "Model downloaded");

  const activeDel = () => d.querySelector('#models .mrow input[value="small"]').parentElement.querySelector('button[data-a="del"]');
  check("the model in use can be removed too — that is the way out of a full disk", !!activeDel(), true);
  activeDel().click(); await sleep(150);
  check("removing the model in use warns what it costs", d.querySelector(".modal p").textContent.includes("Recognition stops"), true);
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("and the program is told this was meant", w.delCalls[w.delCalls.length - 1], ["small", true]);

  tab("models"); await sleep(60);
  d.getElementById("adv_open").click(); await sleep(30);
  d.getElementById("adv_go").click(); await sleep(150);
  check("the recommendation is a plan, not a sentence", d.querySelectorAll("#adv_out .advrow").length, 2);
  check("the plan says which one is already here", [...d.querySelectorAll("#adv_out .advstate")].map(s=>s.className.includes("ok")), [false, true]);
  const dlBefore = w.dlCalls.length;
  d.querySelector("#adv_out button.mini").click(); await sleep(80);
  check("applying the plan asks about the download", !!d.querySelector(".modal-bg"), true);
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("applying the plan downloads what is missing", w.dlCalls.length > dlBefore, true);
  tab("about"); await sleep(60);
  check("about section shown", shown("about"), true);
  check("about carries version, help and author", d.querySelectorAll("#p-about .card").length, 3);

  const before = w.saveCalls;
  const sw = d.getElementById("auto_enter");
  sw.checked = !sw.checked; sw.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(220);
  check("a toggle applies itself, no Save needed", w.saveCalls > before, true);

  d.querySelector('.lvlb[data-l="simple"]').click(); await sleep(320);
  check("switching the mode reports a plain save", d.getElementById("st_saved").textContent, "Saved");
  check("simple mode hides advanced rows", d.querySelectorAll("#p-dictation .row[data-adv].hidden").length > 0, true);
  check("disclosure button offered", !!d.querySelector("#p-dictation .moreb"), true);
  d.querySelector("#p-dictation .moreb").click(); await sleep(60);
  check("disclosure reveals them", d.querySelectorAll("#p-dictation .row[data-adv].hidden").length, 0);
  check("no permanent mode line in the status bar", !!d.getElementById("st_level"), false);
  check("no switching from the status bar", !!d.getElementById("st_levelbtn"), false);
  check("simple mode folds the expert text blocks away", d.querySelectorAll("#p-text .card[data-adv].hidden").length, 3);
  check("but punctuation and the dictionary stay in sight", d.querySelectorAll("#p-text .card:not([data-adv])").length >= 2, true);

  const omni = d.getElementById("omni");
  omni.value = "S_PORT"; omni.dispatchEvent(new w.Event("input")); await sleep(120);
  check("search jumps to the section", shown("system"), true);
  check("search highlights the row", d.querySelectorAll(".row.hit").length, 1);
  check("search says how many it found", /^1\/\d+$/.test(d.getElementById("ocount").textContent), true);
  const firstHit = d.querySelector(".hit");
  const total = Number(d.getElementById("ocount").textContent.split("/")[1]);
  omni.dispatchEvent(new w.KeyboardEvent("keydown", { key: "Enter", bubbles: true })); await sleep(120);
  check("Enter walks to the next match", total > 1 ? d.querySelector(".hit") !== firstHit : d.querySelector(".hit") === firstHit, true);
  omni.value = "zzzznothing"; omni.dispatchEvent(new w.Event("input")); await sleep(120);
  check("a search with no matches says so", d.getElementById("ocount").textContent, "none");
  check("a search with no matches highlights nothing", d.querySelectorAll(".hit").length, 0);
  omni.value = "S_SEC_CMD"; omni.dispatchEvent(new w.Event("input")); await sleep(200);
  check("search finds a heading, not only a setting row", shown("text"), true);
  check("and the heading itself is what is highlighted", !!d.querySelector("#p-text .sect.hit"), true);
  omni.value = ""; omni.dispatchEvent(new w.Event("input")); await sleep(60);

  d.querySelector('.lvlb[data-l="all"]').click(); await sleep(80);
  check("no save button left", !!d.querySelector(".footer"), false);

  const drag0 = w.dragCalls;
  const down = (el) => el.dispatchEvent(new w.MouseEvent("mousedown", { bubbles: true, button: 0 }));
  down(d.getElementById("omni"));
  down(d.querySelector('.lvlb[data-l="simple"]'));
  check("title bar controls are not swallowed by the window drag", w.dragCalls, drag0);
  down(d.querySelector(".header h1"));
  check("empty title bar still drags the window", w.dragCalls, drag0 + 1);
  check("version sits in the status bar", !!d.querySelector("#statusbar #ver"), true);
  check("version left the title bar", !!d.querySelector(".header #ver"), false);


  tab("mic"); await sleep(120);
  const micSel = d.getElementById("mic_device");
  micSel.value = "dev1"; micSel.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  const savesBefore = w.saveCalls;
  w.setMicFails(true);
  micSel.value = "dev2"; micSel.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("a microphone that refuses is rolled back", micSel.value, "dev1");
  check("the refusal stays on screen", d.getElementById("mic_err").textContent, "Microphone busy");
  check("a refused microphone is not saved", w.saveCalls, savesBefore);
  w.setMicFails(false);

  w.setModelStates("missing", "downloading"); await sleep(1700);
  check("a missing model is not green", d.getElementById("state_ru_led").className.includes("on"), false);
  check("a missing model warns", d.getElementById("state_ru_led").className.includes("warn"), true);
  check("a missing model offers to download it", d.getElementById("state_ru_btn").textContent, "Download");
  check("a model being downloaded is not green", d.getElementById("state_other_led").className.includes("on"), false);
  w.setModelStates("ready", "ready"); await sleep(1700);
  check("an installed model is green", d.getElementById("state_ru_led").className.includes("on"), true);

  w.setRemote(true); await sleep(1700);
  check("remote recognition is announced", d.getElementById("st_remote").textContent, "REMOTE");
  w.setRemote(false); await sleep(1700);
  check("local recognition says nothing", d.getElementById("st_remote").textContent, "");

  check("with everything working the failure card stays away", d.getElementById("state_backend").style.display, "none");
  w.setBackendErr("port 8910 is busy"); await sleep(1700);
  check("a recognizer that will not start says why", d.getElementById("state_backend_text").textContent, "port 8910 is busy");
  check("and the card is on screen", d.getElementById("state_backend").style.display !== "none", true);
  d.getElementById("state_retry").click(); await sleep(1800);
  check("the card offers to try again, and it reaches the program", w.retryCalls, 1);
  check("a recognizer that came up takes the card away", d.getElementById("state_backend").style.display, "none");

  tab("dictation"); await sleep(80);
  const enterRow = d.getElementById("auto_enter").closest(".row");
  check("the label is tied to its switch", enterRow.querySelector("label").htmlFor, "auto_enter");
  const wasOn = d.getElementById("auto_enter").checked;
  enterRow.querySelector("label").click(); await sleep(200);
  check("clicking the text flips the switch", d.getElementById("auto_enter").checked, !wasOn);
  enterRow.querySelector("label").click(); await sleep(200);
  check("no unsaved-changes dialog left", !!d.getElementById("modalbg"), false);

  tab("system"); await sleep(80);
  const port = d.getElementById("server_port");
  port.value = "8999"; port.dispatchEvent(new w.Event("change", { bubbles: true }));
  await sleep(400);
  check("changing the port restarts the recognizer instead of asking", w.lastSave.message, "Restarting the recognizer…");
  check("nothing in the status bar asks for a restart", !!d.getElementById("st_pend"), false);

  check("the wizard stays out of the way until it is asked for", d.getElementById("wiz").classList.contains("on"), false);
  w.wizStart(); await sleep(150);
  check("the wizard covers the window", d.getElementById("wiz").classList.contains("on"), true);
  check("the wizard opens on the welcome step", d.getElementById("wz0").classList.contains("on"), true);
  check("the wizard walks five steps", d.querySelectorAll("#wizdots i").length, 5);
  check("there is no way back from the first step", d.getElementById("wiz_back").style.display, "none");
  check("the interface language is offered right away", d.getElementById("wiz_ui").value, "en");

  d.getElementById("wiz_next").click(); await sleep(250);
  check("the second step is the model", d.getElementById("wz1").classList.contains("on"), true);
  check("the wizard proposes a plan, not a list", d.querySelectorAll("#wiz_plan .advrow").length, 2);
  check("the plan says what is already here", [...d.querySelectorAll("#wiz_plan .advstate")].map((s) => s.className.includes("ok")), [false, true]);
  check("the download button carries the size", d.getElementById("wiz_dl").textContent.includes("232 MB"), true);
  const wizDlBefore = w.dlCalls.length;
  d.getElementById("wiz_dl").click(); await sleep(200);
  check("the wizard downloads what the plan is missing", w.dlCalls.length > wizDlBefore, true);
  check("the wizard shows the download running", d.getElementById("wiz_dlrow").style.display, "");
  check("the wizard does not let you walk past a download", d.getElementById("wiz_next").disabled, true);
  w.finishDl("gigaam-v3"); await sleep(1200);
  check("a finished download opens the way on", d.getElementById("wiz_next").disabled, false);
  const appliedIds = w.saveForms.slice(-6).map((f) => f.model_id);
  check("and both models of the plan are applied, not just the first",
    ["gigaam-v3", "small"].every((id) => appliedIds.includes(id)), true);

  d.getElementById("wiz_next").click(); await sleep(300);
  check("the third step names the shortcut", d.getElementById("wiz_hot").textContent, "ctrl+win");
  check("the third step offers the microphones", d.getElementById("wiz_mic").options.length, 3);
  await sleep(300);
  check("the wizard meter follows the microphone", meterMoves(d, "wiz_micbar"), true);

  d.getElementById("wiz_next").click(); await sleep(250);
  check("the fourth step waits for a phrase", d.getElementById("wiz_tryout").textContent, "Waiting for the first phrase…");
  check("the fourth step offers a place to dictate into", !!d.getElementById("wiz_try"), true);
  w.setLastAt(1700000000000); await sleep(1000);
  check("the wizard repeats what it heard", d.getElementById("wiz_tryout").textContent, "Heard: hello");

  d.getElementById("wiz_next").click(); await sleep(250);
  check("the last step offers to start with Windows", !!d.getElementById("wiz_auto"), true);
  check("the last button finishes instead of going on", d.getElementById("wiz_next").textContent, "Finish");
  check("there is nothing left to skip", d.getElementById("wiz_skip").style.display, "none");
  d.getElementById("wiz_auto").checked = true;
  d.getElementById("wiz_next").click(); await sleep(250);
  check("finishing hides the wizard", d.getElementById("wiz").classList.contains("on"), false);
  check("finishing is remembered", w.wizardDone, 1);
  check("the autostart answer is passed on", w.autorunCalls[w.autorunCalls.length - 1], true);
  check("finishing lands on the status screen", shown("state"), true);

  w.wizStart(); await sleep(150);
  d.getElementById("wiz_skip").click(); await sleep(250);
  check("skipping closes the wizard too", d.getElementById("wiz").classList.contains("on"), false);
  check("skipping is remembered, so it does not ask again", w.wizardDone, 2);

  const styleAt = html.indexOf("<style>");
  const css = html.slice(html.indexOf("}", html.indexOf(":root{", styleAt)), html.indexOf("</style>"));
  const shaped = [
    [".scard{", "border-radius:var(--r)"],
    [".hero{", "border-radius:var(--r)"],
    [".modal{", "border-radius:var(--r)"],
    ["button.btn{", "border-radius:calc(var(--r) * .5)"],
    ["button.mini{", "border-radius:calc(var(--r) * .5)"],
    ["input[type=text],input[type=number],select{", "border-radius:calc(var(--r) * .55)"],
    ["input[type=text],input[type=number],select{", "padding:var(--fieldpad)"],
    [".row select,.row input[type=text]{", "padding:var(--ctlpad)"],
    ["input[type=checkbox]{", "border-radius:calc(var(--r) * .8)"],
  ];
  for (const [sel, want] of shaped) {
    const at = css.indexOf(sel);
    const rule = at < 0 ? "" : css.slice(at, css.indexOf("}", at));
    check(`${sel} takes its corners from the skin`, rule.includes(want), true);
  }
  const skinned = [
    ["body::after{", "opacity:var(--scan)"],
    [".header h1{", "animation:var(--flicker)"],
    [".header{", "background:var(--titlebg)"],
    [".snav{", "background:var(--sidebg)"],
    ["button.btn{", "background:var(--btnbg)"],
    ["button.btn{", "text-transform:var(--caps)"],
    [".lvlb.on{", "background:var(--selbg)"],
    [".modal-bg{", "background:var(--scrim)"],
    [".hotkey-val{", "background:var(--keybg)"],
    [".scard .led.on{", "background:var(--hi)"],
    [".miclevel i{", "background:var(--hi)"],
    [".mock-dot{", "background:var(--rec)"],
  ];
  for (const [sel, want] of skinned) {
    const at = css.indexOf(sel);
    const rule = at < 0 ? "" : css.slice(at, css.indexOf("}", at));
    check(`${sel} follows the skin (${want})`, rule.includes(want), true);
  }
  const literals = [...new Set((css.match(/#[0-9a-f]{6}/g) || []))].sort();
  check("no colour is written into the stylesheet by hand", literals, ["#3c1212", "#7a2e2e"]);
  check("no face is nailed to Consolas outside the skin", /font:[^;]*Consolas/.test(css), false);

  check("no page errors", errors, []);

  if (failures.length) {
    console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
    process.exit(1);
  }
  console.log("\nall UI checks passed");
  dom.window.close();
  process.exit(0);
})().catch((e) => {
  console.error("harness crashed:", e.message);
  process.exit(1);
});
