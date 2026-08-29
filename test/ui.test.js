const fs = require("fs");
const path = require("path");
const { JSDOM } = require("jsdom");

const html = fs.readFileSync(path.join(__dirname, "page.html"), "utf8");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
function meterMoves(d, id) {
  const bars = [...d.querySelectorAll(`#${id} i`)];
  if (!bars.length) return false;
  return bars.some((b) => b.classList.contains("on"));
}
function searchFinds(w, d, needle) {
  const hits = w.searchMatches(needle.toLowerCase());
  return hits.some((el) => el.closest("#p-help"));
}

let llmState = {
  installed: [{ file: "model.gguf", size: 4929, active: true, loaded: true }],
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
        status_line: "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free",
        active_model: "GigaAM v3", active_state: ruState, active_lang: "RU",
        assigned: [
          { model: "GigaAM v3", langs: "RU", state: ruState, current: true },
          { model: "Small", langs: "Detect itself, EN", state: otherState, current: false },
        ],
        installed_models: ["Small", "my-model"],
        loaded_now: "GigaAM v3", week_line: "12 dictations · 3400 characters",
        llm_ok: true, mic_ok: true, last_at: lastAt, last_app: "chrome.exe",
        remote: remote, backend_err: backendErr, post_err: postErr,
        badges: { mic: micBadge, models: "2", system: "" } });
    let modelStates = { base: "absent", small: "active", "medium-q5_0": "absent", "gigaam-v3": "absent", "moonshine-uk": "absent" };
    window.dlCalls = [];
    window.cancelCalls = [];
    window.appModelDl = async (id) => { window.dlCalls.push(id); modelStates[id] = "downloading"; };
    window.finishDl = (id) => { modelStates[id] = "installed"; };
    window.setModelState = (id, st) => { modelStates[id] = st; };
    window.appModelCancel = async (id) => { window.cancelCalls.push(id); modelStates[id] = "absent"; return true; };
    window.unloadCalls = 0;
    window.appUnloadEngines = async () => { window.unloadCalls++; };
    window.postKeys = [];
    let postKeySet = false;
    let postErr = "";
    window.setPostErr = (v) => { postErr = v; };
    window.postTests = [];
    let postTestOK = true;
    window.setPostTestOK = (v) => { postTestOK = v; };
    window.appPostTest = async (url, model, key, timeout) => {
      window.postTests.push([url, model, key, timeout]);
      if(!postTestOK){
        postErr = "no credits left";
        return JSON.stringify({ ok: false, severity: "error", message: "no credits left" });
      }
      postErr = "";
      return JSON.stringify({ ok: true, severity: "ok", message: "The server answered" });
    };
    window.appSetPostKey = async (k) => { window.postKeys.push(k); postKeySet = !!k; return JSON.stringify({ ok: true, severity: "ok", message: "Saved" }); };
    window.appPostKeySet = async () => postKeySet;
    window.folderOpens = 0;
    window.appOpenModelsFolder = () => { window.folderOpens++; };
    window.linkOpens = [];
    window.appModelLink = (id) => { window.linkOpens.push(id); };
    window.appModels = async () =>
      JSON.stringify([
        { id: "base", name: "Base", desc: "fast", size: 142, state: modelStates.base, pct: 12, engine: "whisper", langs: "*", auto: true, translate: true, speed: 5, accuracy: 2 },
        { id: "small", name: "Small", desc: "balanced", size: 466, ram: 921, state: modelStates.small, engine: "whisper", langs: "*", auto: true, translate: true, serves: ["auto"], loaded: true, speed: 3, accuracy: 3 },
        { id: "medium-q5_0", name: "Medium (q5)", desc: "recommended", size: 539, state: modelStates["medium-q5_0"], pct: 5, engine: "whisper", langs: "*", auto: true, translate: true, speed: 2, accuracy: 4 },
        { id: "gigaam-v3", name: "GigaAM v3", desc: "russian", size: 232, state: modelStates["gigaam-v3"], pct: 5, engine: "sherpa", langs: "ru", punct: true, serves: ["ru"], speed: 5, accuracy: 5 },
        { id: "moonshine-uk", name: "Moonshine Base uk", desc: "ukrainian", size: 135, state: modelStates["moonshine-uk"], engine: "sherpa", langs: "uk", manual: true, link: "https://example.com/moonshine", speed: 5, accuracy: 3 },
        { id: "local:my-model", name: "my-model", desc: "found in the models folder", size: 900, state: "installed", engine: "whisper", langs: "*", custom: true },
      ]);
    window.appLLMSearch = async () =>
      JSON.stringify({ repos: [{ id: "org/Repo-GGUF", downloads: 1234, updated: "2026-01-01" }] });
    window.appLLMFiles = async () =>
      JSON.stringify({ files: [{ file: "q4.gguf", size: 4000, fit: "ok", need: 6166 }, { file: "q8.gguf", size: 9000, fit: "bad", need: 13000 }] });
    window.llmUnloads = 0;
    window.appLLMUnload = async () => {
      window.llmUnloads++;
      llmState.installed.forEach(m => { m.loaded = false; });
    };
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
      "appAuthorLink", "appReload",
      "appPreviewSound", "appMin", "appClose",
      "appDoUpdate", "appReady", "appJSError",
    ]) {
      window[name] = () => {};
    }
    window.captureCalls = 0;
    window.appCapture = () => { window.captureCalls++; };
    window.appCaptureCombo = () => { window.captureCalls++; };
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
  for (const id of ["hotkey"]) {
    const el = d.getElementById(id);
    check("the " + id + " shortcut is one control, not a chip beside a button", el.tagName, "BUTTON");
    check("and " + id + " says what pressing it does", !!el.dataset.tip, true);
    const before = w.captureCalls;
    el.click();
    check("and pressing " + id + " asks for a new combination", w.captureCalls > before, true);
  }
  check("no separate set button is left beside them", [...d.querySelectorAll("[id60_set]")].map(e=>e.id), []);
  check("status names the model behind the current language", d.getElementById("state_active").textContent, "GigaAM v3");
  check("and the language stands beside it", d.getElementById("state_active_lang").textContent, "RU");
  check("the models in use are listed", d.querySelectorAll("#state_assigned .arow").length, 2);
  check("each with the languages it serves", d.querySelector("#state_assigned .arow .alangs").textContent, "RU");
  check("the current one is marked", d.querySelector("#state_assigned .arow").className.includes("on"), true);
  check("installed models are listed in one line", d.getElementById("state_installed").textContent, "Small, my-model");
  check("the memory row says what is loaded right now", d.getElementById("state_loaded").textContent, "GigaAM v3");
  check("a week of dictations is summed up", d.getElementById("state_week").textContent, "12 dictations · 3400 characters");
  check("and the row is shown when there is something to sum", d.getElementById("state_week_row").style.display, "");
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
  const micBar = d.getElementById("state_mic_bar");
  const lit = (el) => [...el.querySelectorAll("i")].filter(b=>b.classList.contains("on")).length;
  const settle = (el, lvl) => { for (let i = 0; i < 14; i++) w.paintMeter(el, lvl); };
  settle(micBar, 0.005);
  check("a quiet room leaves the meter dark", lit(micBar), 0);
  w.paintMeter(micBar, 0.6);
  check("a loud phrase lights it up", lit(micBar), micBar.querySelectorAll("i").length);
  check("silence never leaves the floor", w.micHeard(0.004), 0);
  check("room noise barely stirs it", w.micHeard(0.012) < 0.2, true);
  check("a quiet voice already shows", w.micHeard(0.05) > 0.35, true);
  check("an ordinary one reaches the middle", w.micHeard(0.15) > 0.7, true);
  check("a loud one nearly fills the meter", w.micHeard(0.5) > 0.85, true);
  check("the status meter is told to fill its card", micBar.classList.contains("grow"), true);
  Object.defineProperty(micBar, "clientWidth", { value: 200, configurable: true });
  settle(micBar, 0.05);
  check("a wider card gets more bars, not wider ones", micBar.querySelectorAll("i").length, 33);
  check("the meter fills to the loudness instead of drawing a wave", lit(micBar), Math.round(w.micHeard(0.05) * 33));
  check("and what is lit is one run from the left", [...micBar.querySelectorAll("i")].findIndex(b=>!b.classList.contains("on")), lit(micBar));
  const held = lit(micBar);
  w.paintMeter(micBar, 0);
  check("it lets go slowly rather than blinking out", lit(micBar) > 0 && lit(micBar) < held, true);
  settle(micBar, 0);
  check("and goes dark when the room falls quiet", lit(micBar), 0);
  Object.defineProperty(micBar, "clientWidth", { value: 60, configurable: true });
  w.paintMeter(micBar, 0.3);
  check("a narrower card gives the bars back", micBar.querySelectorAll("i").length, 10);
  const rowBar = d.getElementById("mic_bar");
  check("the meter on the microphone screen fills its row too", rowBar.classList.contains("grow"), true);
  check("and its row is the one marked for it", rowBar.closest(".row").classList.contains("lvlrow"), true);
  check("the word for a quiet room stands before the meter, not after it", [...rowBar.closest(".row").children].map(el=>el.id || el.tagName.toLowerCase()), ["label", "mic_hint", "mic_bar"]);
  Object.defineProperty(rowBar, "clientWidth", { value: 320, configurable: true });
  w.paintMeter(rowBar, 0.05);
  check("so it counts its own segments as well", rowBar.querySelectorAll("i").length, 53);
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
  check("the models badge explains itself", d.getElementById("badge_models").dataset.tip, "Installed models");
  check("icon buttons carry a name for screen readers", d.getElementById("mic_refresh").getAttribute("aria-label"), "S_MIC_REFRESH");
  check("and the same words are the hover hint", d.getElementById("mic_refresh").dataset.tip, "S_MIC_REFRESH");
  check("the browser's own tooltip is out of the way", d.getElementById("mic_refresh").hasAttribute("title"), false);
  check("the open section is announced", d.querySelector(".nav[data-p=state]").getAttribute("aria-selected"), "true");

  tab("system"); await sleep(200);
  const themeSel = d.getElementById("theme");
  const skinSel = d.getElementById("skin");
  check("design and colour are two separate choices", [!!skinSel, !!themeSel], [true, true]);
  check("five designs are offered", [...skinSel.options].map(o=>o.value), ["terminal", "editor", "neon", "soft", "paper"]);
  check("the two light ones come last", [...skinSel.options].slice(-2).map(o=>o.value), ["soft", "paper"]);
  for (const light of ["soft", "paper"]) {
    w.applyThemeVars(light);
    const root = d.documentElement.style;
    const lum = (v) => { const h = v.trim().replace("#",""); return 0.2126*parseInt(h.slice(0,2),16) + 0.7152*parseInt(h.slice(2,4),16) + 0.0722*parseInt(h.slice(4,6),16); };
    check(`the ${light} design paints a light ground`, lum(root.getPropertyValue("--bg")) > 140, true);
    check(`and writes on it in dark ink`, lum(root.getPropertyValue("--green")) < 110, true);
    check(`and asks for no halo`, [root.getPropertyValue("--glow"), root.getPropertyValue("--higlow"), root.getPropertyValue("--amberglow"), root.getPropertyValue("--badglow")], ["none", "none", "none", "none"]);
    check(`and tells the system controls it is light`, root.getPropertyValue("--scheme"), "light");
    check(`and no scanlines`, root.getPropertyValue("--scan"), "0");
    check(`and still fills in the warning surfaces`, /^#[0-9a-f]{6}$/.test(root.getPropertyValue("--badbg").trim()), true);
  }
  w.applyThemeVars("terminal:green");
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
  check("controls take their height from the design", d.documentElement.style.getPropertyValue("--fieldpad"), "6px 11px");
  check("and one type size for every field", d.documentElement.style.getPropertyValue("--ctlfs"), "12.5px");
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

  check("eleven sections in the sidebar", d.querySelectorAll(".nav").length, 11);
  check("in three groups", d.querySelectorAll(".ngrp").length, 3);
  check("in the agreed order", [...d.querySelectorAll(".nav")].map(b=>b.dataset.p),
    ["state", "system", "mic", "history", "dictation", "models", "text", "post", "help", "about", "contacts"]);
  check("no mode switch in the header", !!d.querySelector(".lvlsw"), false);
  check("no disclosure buttons anywhere", d.querySelectorAll(".moreb").length, 0);
  check("nothing is folded away", d.querySelectorAll("[data-adv].hidden").length, 0);

  tab("about"); await sleep(200);
  const toc = d.querySelector("#p-help .toc");
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
    "type_mode", "overlay", "overlay_position", "overlay_text", "mic_device", "mic_bar", "beep", "sound_theme",
    "langlist", "munload", "language", "threads", "punctuation", "dictbody", "profbody",
    "tr_default", "translate_target", "translate_ask", "translate_ask_seconds",
    "tl_en", "ui_language", "upd_check", "check_updates", "server_autostart", "server_port",
    "server_exe", "server_url", "llm_sum", "llm_catalog", "ver2", "autorun", "post_enabled",
  ];
  const missing = everySetting.filter((id) => !d.getElementById(id));
  check("every setting is present in the new window", missing, []);

  tab("models"); await sleep(80);
  check("models section shown", shown("models"), true);
  const rowFor = (lang) => d.querySelector('#langlist .lrow[data-lang="'+lang+'"]');
  check("every language owns a row", d.querySelectorAll("#langlist .lrow").length, 9);
  check("no old library, filters or advisor are left", [d.getElementById("models"), d.getElementById("presets"), d.getElementById("mlang"), d.getElementById("mfind"), d.getElementById("advisor")].filter(Boolean).length, 0);
  check("a language with its own model is bright and names it", !rowFor("ru").className.includes("dim") && rowFor("ru").querySelector(".lmodel").textContent === "GigaAM v3", true);
  check("a language without its own model is dim", rowFor("de").className.includes("dim"), true);
  check("and says honestly whom it follows", rowFor("de").querySelector(".lmodel").textContent, "as Auto-detect · Medium (q5)");
  check("no picker is open before a click", d.querySelectorAll("#langlist .lpick").length, 0);
  rowFor("ru").click(); await sleep(120);
  check("clicking a language unfolds the choice under it", d.querySelectorAll("#langlist .lpick").length, 1);
  const cards = () => [...d.querySelectorAll("#langlist .lpick .pcard")];
  check("only models that serve this language are offered", cards().map(c=>c.dataset.id).sort().join("+"), "base+gigaam-v3+local:my-model+medium-q5_0+small");
  check("no radio buttons anywhere in the choice", d.querySelectorAll('#langlist input[type="radio"]').length, 0);
  const ruOrder = cards().map(c=>c.dataset.id).join("+");
  check("the assigned model wears the mark", d.querySelector('#langlist .pcard[data-id="gigaam-v3"]').className.includes("cur"), true);
  check("no chip claims assignment — the switch says it", d.querySelector('#langlist .pcard[data-id="gigaam-v3"] .pchip').textContent, "recommended");
  check("the model story hides behind an info badge", d.querySelector('#langlist .pcard[data-id="gigaam-v3"] .pinfo').dataset.tip, "russian");
  check("and is no longer printed on the card", d.querySelector('#langlist .pcard[data-id="gigaam-v3"] .ptop').textContent.includes("russian"), false);
  check("every catalog model carries the badge", cards().filter(c => c.dataset.id !== "local:my-model").every(c => c.querySelector(".pinfo")), true);
  check("every known model measures itself in two bars", d.querySelectorAll("#langlist .mbar").length, 8);
  check("the bars are filled to the model, not all alike", [...d.querySelector('#langlist .pcard[data-id="base"]').querySelectorAll(".mtrack i")].map(i=>i.style.width), ["40%", "100%"]);
  check("a hand-copied model shows no bars — its powers are unknown", d.querySelector('#langlist .pcard[data-id="local:my-model"]').querySelectorAll(".mbar").length, 0);
  check("and admits its languages are unknown", d.querySelector('#langlist .pcard[data-id="local:my-model"] .pmeta').textContent.includes("languages: unknown"), true);
  check("a multilingual model counts its languages in words", d.querySelector('#langlist .pcard[data-id="small"] .pmeta').textContent.includes("languages: 99"), true);
  check("a model that translates says so in words", d.querySelector('#langlist .pcard[data-id="small"] .pmeta').textContent.includes("translates to English"), true);
  check("the words recognition-only are gone from the page", d.getElementById("p-models").textContent.includes("recognition only"), false);
  check("ram is estimated where it is known", d.querySelector('#langlist .pcard[data-id="small"] .pmeta').textContent.includes("≈921 MB RAM"), true);
  check("the models folder can be opened for hand-copied files", !!d.querySelector('#p-models button[onclick="appOpenModelsFolder()"]'), true);
  const eject = d.querySelector('#langlist .pcard[data-id="small"] button[data-a="unload"]');
  check("the loaded model offers to leave the memory", !!eject, true);
  eject.click(); await sleep(150);
  check("and the program is asked to unload it", w.unloadCalls, 1);
  check("there is a plain unload row as well", !!d.getElementById("munload"), true);
  const saveBefore = w.saveCalls;
  d.querySelector('#langlist .pcard[data-id="small"]').click(); await sleep(200);
  check("a click on the card itself chooses nothing", w.saveCalls, saveBefore);
  const flick = async (id) => { const s = d.querySelector('#langlist .pcard[data-id="'+id+'"] input.psw'); s.checked = true; s.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200); };
  await flick("small");
  check("the switch is the choice", w.saveCalls > saveBefore, true);
  check("and it pins the model to this language", w.lastSaveForm.lang_models.ru, "small");
  check("the picker stays open and the mark moves", d.querySelector('#langlist .pcard.cur').dataset.id, "small");
  check("but the list keeps its order", cards().map(c=>c.dataset.id).join("+"), ruOrder);
  const backBtn = d.querySelector("#langlist .lpick button.mini:last-child");
  check("the way back to Auto-detect is offered", backBtn.textContent, "Back to Auto-detect");
  backBtn.click(); await sleep(200);
  check("taking it clears the language's own model", "ru" in (w.lastSaveForm.lang_models || {}), false);
  check("and the row dims to the inherited truth", rowFor("ru").className.includes("dim"), true);
  await flick("small");
  check("choosing again is one flick of the switch", w.lastSaveForm.lang_models.ru, "small");
  const swFor = (id) => d.querySelector('#langlist .pcard[data-id="'+id+'"] input.psw');
  check("every model card carries a switch", d.querySelectorAll('#langlist .pcard input.psw').length, 5);
  check("the chosen card's switch is on and the others are off", [swFor("small").checked, swFor("base").checked, swFor("local:my-model").checked], [true, false, false]);
  swFor("local:my-model").checked = true; swFor("local:my-model").dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("a switch turned on is the choice too", w.lastSaveForm.lang_models.ru, "local:my-model");
  const swOff = swFor("local:my-model");
  swOff.checked = false; swOff.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("turning it off returns the language to Auto-detect", "ru" in (w.lastSaveForm.lang_models || {}), false);
  check("and the row dims to the inherited truth again", rowFor("ru").className.includes("dim"), true);
  await flick("small");
  check("the choice comes back with the switch", w.lastSaveForm.lang_models.ru, "small");
  rowFor("uk").click(); await sleep(120);
  check("the licensed model's switch is locked", d.querySelector('#langlist .pcard[data-id="moonshine-uk"] input.psw').disabled, true);
  check("switching languages swaps the picker", cards().map(c=>c.dataset.id).includes("moonshine-uk"), true);
  check("the licensed model offers no download button", !!d.querySelector('#langlist .pcard[data-id="moonshine-uk"] button[data-a="dl"]'), false);
  check("it explains the licence instead", d.querySelector('#langlist .pcard[data-id="moonshine-uk"]').textContent.includes("licence forbids"), true);
  d.querySelector('#langlist .pcard[data-id="moonshine-uk"]').click(); await sleep(150);
  check("and cannot be picked at all", "uk" in (w.lastSaveForm.lang_models || {}), false);
  const linkBtn = d.querySelector('#langlist .pcard[data-id="moonshine-uk"] button[data-a="link"]');
  check("and offers the source link", !!linkBtn, true);
  linkBtn.click(); await sleep(60);
  check("the link opens the model page", w.linkOpens, ["moonshine-uk"]);
  rowFor("ru").click(); await sleep(120);
  check("the mark carries both shapes", [d.querySelectorAll(".mk.mic").length, d.querySelectorAll(".mk.face").length], [2, 2]);
  w.applyThemeVars("soft");
  check("the soft design shows the face", [d.documentElement.style.getPropertyValue("--markmic"), d.documentElement.style.getPropertyValue("--markface")], ["none", "block"]);
  w.applyThemeVars("terminal:green");
  check("and every other one keeps the microphone", [d.documentElement.style.getPropertyValue("--markmic"), d.documentElement.style.getPropertyValue("--markface")], ["block", "none"]);

  const recLangs = [...d.getElementById("language").options].map(o=>o.value);
  check("italian can be dictated too", recLangs.includes("it"), true);
  check("the routing table is gone", !!d.getElementById("routing"), false);
  check("the language moved in with the models", !!d.querySelector("#p-models #language"), true);
  check("the threads moved in with the server", !!d.querySelector("#p-system #threads"), true);
  check("the engine hint line is gone with the manual", !!d.getElementById("tr_engine"), false);
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
  check("the dictionary shows its words as chips", d.querySelectorAll("#dictbody .cmdchip").length, 2);
  check("a word chip carries no pencil", d.querySelector("#dictbody .cmdchip .redit"), null);
  check("with a whisper model no warning shows", d.getElementById("dict_warn").style.display, "none");
  d.getElementById("dict_add").click(); await sleep(120);
  check("adding a word opens a titled dialog", d.querySelector(".modal .mtitle").textContent, "Adding a word");
  const dw = d.querySelector(".modal .dword");
  dw.value = "endpoint, webhook, Docker"; dw.dispatchEvent(new w.Event("input", { bubbles: true }));
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("commas bring several words at once, without duplicates", d.querySelectorAll("#dictbody .cmdchip").length, 4);
  check("and the whole set is saved as the prompt", w.lastSaveForm.whisper_prompt, "GitHub, Docker, endpoint, webhook");
  d.querySelectorAll("#dictbody .cmdchip .rdel")[3].click(); await sleep(300);
  check("the cross removes a word and saves", w.lastSaveForm.whisper_prompt, "GitHub, Docker, endpoint");
  w.setModelState("small", "installed"); w.setModelState("gigaam-v3", "active"); await w.refreshModels(); await sleep(80);
  check("a non-whisper model raises the yellow warning", d.getElementById("dict_warn").style.display, "");
  check("and it names the model", d.getElementById("dict_warn").textContent.includes("GigaAM v3"), true);
  w.setModelState("gigaam-v3", "absent"); w.setModelState("small", "active"); await w.refreshModels(); await sleep(80);
  check("bringing whisper back clears the warning", d.getElementById("dict_warn").style.display, "none");
  check("punctuation modes offered", d.getElementById("punctuation").options.length, 3);

  tab("dictation"); await sleep(60);
  const zones = d.querySelectorAll("#ovscheme .ovzone");
  check("the screen is divided into nine places", zones.length, 9);
  check("the plate settings live under one heading", [!!d.querySelector("#ovcard .blklbl"), !!d.querySelector("#ovcard .hint")], [true, true]);
  const ovMaster = d.getElementById("overlay");
  ovMaster.checked = false; ovMaster.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("turning the plate off disables the rest of the block", [d.getElementById("ov_free").disabled, d.getElementById("ov_caret").disabled, d.getElementById("overlay_text").disabled], [true, true, true]);
  check("and dims their rows", d.getElementById("ov_free").closest(".row").className.includes("dimmed"), true);
  check("the screen steps aside as well", d.getElementById("ovscheme").className.includes("off"), true);
  ovMaster.checked = true; ovMaster.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("turning it back on wakes them up", [d.getElementById("ov_free").disabled, d.getElementById("ov_free").closest(".row").className.includes("dimmed")], [false, false]);
  check("the plate starts at the bottom of the screen", d.getElementById("ovmini").style.top, "85%");
  check("the screen is drawn as a monitor with a stand", [!!d.querySelector("#ovscheme .ovcase"), !!d.querySelector("#ovscheme .ovneck"), !!d.querySelector("#ovscheme .ovbase")], [true, true, true]);
  check("the pause option is gone", [!!d.getElementById("pause_hotkey"), !!d.getElementById("pause_clear")], [false, false]);
  d.querySelector('#ovscheme .ovzone[data-pos="top-right"]').click(); await sleep(300);
  check("clicking a place moves the plate and saves itself", w.lastSaveForm.overlay_position, "top-right");
  check("the miniature follows the chosen place", d.getElementById("ovmini").style.left, "88%");
  d.querySelector('#ovscheme .ovzone[data-pos="center"]').click(); await sleep(300);
  check("the middle of the screen is offered too", [w.lastSaveForm.overlay_position, d.getElementById("ovmini").style.left, d.getElementById("ovmini").style.top], ["center", "50%", "50%"]);
  const caretBox = d.getElementById("ov_caret");
  caretBox.checked = true; caretBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("following the cursor is the ninth choice", w.lastSaveForm.overlay_position, "caret");
  check("and the scheme steps aside while it is on", d.getElementById("ovscheme").className.includes("off"), true);
  caretBox.checked = false; caretBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  const crt = d.getElementById("ovcrt");
  const mini = d.getElementById("ovmini");
  crt.getBoundingClientRect = () => ({ left: 0, top: 0, width: 156, height: 80 });
  const freeBox = d.getElementById("ov_free");
  check("a spot of your own is offered as a switch", freeBox.type, "checkbox");
  check("and it is off while a place is chosen", freeBox.checked, false);
  check("dragging does nothing while the switch is off", (() => { const before = w.saveCalls; crt.dispatchEvent(new w.MouseEvent("pointerdown", { clientX: 40, clientY: 20 })); return w.saveCalls === before; })(), true);
  freeBox.checked = true; freeBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("turning it on hands the screen over to you", w.lastSaveForm.overlay_position, "custom");
  check("the places step aside", d.getElementById("ovcrt").className.includes("free"), true);
  check("and the hint changes to say so", d.getElementById("ovpos_sub").textContent, "drag the plate with the mouse — it lands anywhere");
  crt.dispatchEvent(new w.MouseEvent("pointerdown", { clientX: 39, clientY: 20 }));
  crt.dispatchEvent(new w.MouseEvent("pointermove", { clientX: 39, clientY: 20 }));
  crt.dispatchEvent(new w.Event("pointerup")); await sleep(300);
  check("dragging remembers the fraction for this screen", w.lastSaveForm.overlay_custom["1920x1080"], { x: 0.25, y: 0.25 });
  freeBox.checked = false; freeBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("turning it off returns the last chosen place", w.lastSaveForm.overlay_position, "center");
  freeBox.checked = true; freeBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("and your spot is not lost while it was off", [w.lastSaveForm.overlay_position, d.getElementById("ovmini").style.left], ["custom", "25%"]);
  freeBox.checked = false; freeBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  d.querySelector('#ovscheme .ovzone[data-pos="bottom"]').click(); await sleep(300);
  check("with two monitors the screen choice appears", d.getElementById("ovmon_row").style.display, "");
  check("the cursor screen is the first offer", d.getElementById("overlay_monitor").options[0].textContent, "The screen with the cursor");
  check("a monitor goes by the name it reports", d.getElementById("overlay_monitor").options[1].textContent, "DELL U2720Q (1920×1080) ★");
  check("a nameless one keeps its number", d.getElementById("overlay_monitor").options[2].textContent, "Screen 2 (2560×1440)");
  const ovtext = d.getElementById("overlay_text");
  check("showing the text on the plate is a switch", ovtext.type, "checkbox");
  ovtext.checked = false; ovtext.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("turning it off is saved", w.lastSaveForm.overlay_text, false);
  ovtext.checked = true; ovtext.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);

  tab("dictation"); await sleep(60);
  check("the app-rules block is gone", [!!d.getElementById("rulesbody"), !!d.getElementById("rule_add"), !!d.getElementById("rule_last")], [false, false, false]);

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

  tab("dictation"); await sleep(30);
  check("the old turbo warning is gone — honesty lives in the engine line", !!d.getElementById("tr_warn"), false);
  check("the target language carries the honest note", !!d.querySelector("#p-models .row label .sub.warn"), true);
  check("translation moved in with the languages", !!d.querySelector("#p-models #translate_target"), true);
  const trd = d.getElementById("tr_default");
  const ask = d.getElementById("translate_ask");
  const state = () => [
    trd.disabled,
    d.getElementById("translate_target").disabled,
    ask.disabled,
    d.getElementById("translate_ask_seconds").disabled,
    d.getElementById("tl_en").disabled,
    d.getElementById("tl_de").disabled,
  ];
  check("the switch is off out of the box and everything below sleeps", [trd.checked, ...state()], [false, false, true, true, true, true, true]);
  trd.checked = true; trd.dispatchEvent(new w.Event("change")); await sleep(100);
  check("turning it on wakes the block, quiet mode locks other langs", state(), [false, false, false, true, false, true]);
  ask.value = "always"; ask.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(120);
  check("ask-every-time frees only the reachable dialog langs", state(), [false, false, false, true, false, true]);
  check("whisper reaches English alone — other targets are locked", [...d.getElementById("translate_target").options].filter(o=>o.disabled).map(o=>o.value).includes("de"), true);
  ask.value = "timeout"; ask.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(120);
  check("timeout adds its seconds", state(), [false, false, false, false, false, true]);
  ask.value = "never"; ask.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(150);
  check("going quiet with several languages asks first", !!d.querySelector(".modal-bg"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(200);
  check("saying no keeps the asking mode", ask.value, "timeout");
  ask.value = "never"; ask.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(150);
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("agreeing locks the list to the default language", [ask.value, d.getElementById("tl_ru").disabled, d.getElementById("tl_ru").checked, d.getElementById("tl_en").disabled], ["never", true, true, false]);
  check("and the quiet choice is saved as translate-right-away", [w.lastSaveForm.translate_default, w.lastSaveForm.translate_ask], [true, "never"]);
  ask.value = "always"; ask.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(120);
  const tlEn = d.getElementById("tl_en");
  tlEn.checked = false; tlEn.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(150);
  check("the default output language cannot be unchecked", !!d.querySelector(".modal-bg"), true);
  check("the guard speaks with a single button", d.querySelectorAll(".modal .btn").length, 1);
  d.querySelector(".modal .btn.yes").click(); await sleep(150);
  check("and the checkbox springs back", tlEn.checked, true);
  w.setModelState("gigaam-v3", "installed"); await w.refreshModels(); await sleep(80);
  const gflick = async () => { const s = d.querySelector('#langlist .pcard[data-id="gigaam-v3"] input.psw'); s.checked = true; s.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(150); };
  await gflick();
  check("a mute model warns it will silence translation", d.querySelector(".modal p").textContent.includes("cannot translate"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("declining keeps the model and the translation", [w.lastSaveForm.lang_models.ru, trd.checked], ["small", true]);
  check("and the switch snaps back", d.querySelector('#langlist .pcard[data-id="gigaam-v3"] input.psw').checked, false);
  await gflick();
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("agreeing assigns the model and turns translation off", [w.lastSaveForm.lang_models.ru, w.lastSaveForm.translate_default, w.lastSaveForm.translate_ask], ["gigaam-v3", false, "never"]);
  check("the switch is down and locked while the mute model works", [trd.checked, trd.disabled], [false, true]);
  check("and the reason is written beside the switch", d.getElementById("tr_unavail").textContent.includes("GigaAM v3"), true);
  const ssw = d.querySelector('#langlist .pcard[data-id="small"] input.psw');
  ssw.checked = true; ssw.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(250);
  check("bringing a translating model back unlocks the switch", trd.disabled, false);
  w.setModelState("gigaam-v3", "absent"); await w.refreshModels(); await sleep(80);

  tab("text"); await sleep(80);
  check("with no replacements the list says so", d.querySelector("#replbody .ruleempty").textContent, "No replacements yet");
  d.getElementById("repl_add").click(); await sleep(120);
  check("adding opens a dialog with both sides", [!!d.querySelector(".modal .rfrom"), !!d.querySelector(".modal .rto")], [true, true]);
  check("the dialog introduces itself", d.querySelector(".modal .mtitle").textContent, "Adding a replacement");
  check("the fields come first in the form", !!d.querySelector(".modal .fmrow .rfrom") && d.querySelector(".modal .fmrow") === d.querySelector(".modal .rfrom").closest(".fmrow"), true);
  check("a new pair is added, not saved", d.querySelector(".modal .btn.yes").textContent, "Add");
  check("every label explains itself in plain words", d.querySelectorAll(".modal .fmrow label .sub").length, 3);
  check("and the explanations are real text", d.querySelector(".modal .fmrow label .sub").textContent.length > 0, true);
  check("whole words is on by default", d.querySelector(".modal .rwhole").checked, true);
  check("case is off by default", d.querySelector(".modal .rcase").checked, false);
  const rfx = d.querySelector(".modal .rfrom");
  check("an empty field hides its cross", rfx.parentNode.querySelector(".clearx").style.display, "none");
  rfx.value = "abc"; rfx.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(60);
  check("a filled field grows a cross inside", rfx.parentNode.querySelector(".clearx").style.display, "");
  rfx.parentNode.querySelector(".clearx").click(); await sleep(60);
  check("the cross empties the field", rfx.value, "");
  const rlang = d.querySelector(".modal .rlang");
  check("a rule can be pinned to a language", !!rlang, true);
  check("but serves every language by default", rlang.value, "");
  d.querySelector(".modal .rfrom").value = "гит хаб";
  d.querySelector(".modal .rto").value = "GitHub";
  rlang.value = "ru";
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("saving the dialog adds the pair", [w.lastSaveForm.replacements[0].from, w.lastSaveForm.replacements[0].to, w.lastSaveForm.replacements[0].lang], ["гит хаб", "GitHub", "ru"]);
  check("its flags are saved too", [w.lastSaveForm.replacements[0].whole, w.lastSaveForm.replacements[0].match_case], [true, false]);
  check("the pair stands as a chip with a pencil and a cross", [!!d.querySelector("#replbody .cmdchip .redit"), !!d.querySelector("#replbody .cmdchip .rdel")], [true, true]);
  d.querySelector("#replbody .cmdchip .redit").click(); await sleep(120);
  check("the pencil reopens the dialog filled", d.querySelector(".modal .rfrom").value, "гит хаб");
  check("an existing pair is saved, not added", d.querySelector(".modal .btn.yes").textContent, "Save");
  check("and the title says editing", d.querySelector(".modal .mtitle").textContent, "Editing the replacement");
  d.querySelector(".modal .rlang").value = "";
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("editing updates the pair", w.lastSaveForm.replacements[0].lang, "");
  check("the try-a-phrase field is gone", !!d.getElementById("repl_test"), false);
  d.querySelector("#replbody .rdel").click(); await sleep(250);
  check("a replacement can be deleted", w.lastSaveForm.replacements.length, 0);

  check("with no commands the list says so", d.querySelector("#cmdbody .ruleempty").textContent, "No commands yet");
  check("the preset button left the toolbar", !!d.getElementById("cmd_preset"), false);
  d.getElementById("cmd_add").click(); await sleep(120);
  check("the add dialog offers the usual ones as a third button", d.getElementById("fm_extra").textContent, "Add the usual ones");
  check("and grows wider to fit three buttons", d.querySelector(".modal").className.includes("wide"), true);
  d.getElementById("fm_extra").click(); await sleep(300);
  check("the usual commands can be added in one click", d.querySelectorAll("#cmdbody .cmdchip").length, 3);
  check("and the dialog closes itself", !!d.querySelector(".modal-bg"), false);
  check("they carry the phrases of this language", w.lastSaveForm.commands.map((c) => c.phrase), ["new line", "new paragraph", "cancel"]);
  check("and the actions that go with them", w.lastSaveForm.commands.map((c) => c.action), ["newline", "paragraph", "cancel"]);
  d.getElementById("cmd_add").click(); await sleep(120);
  d.getElementById("fm_extra").click(); await sleep(300);
  check("adding them twice changes nothing", d.querySelectorAll("#cmdbody .cmdchip").length, 3);
  d.querySelector("#cmdbody .cmdchip .redit").click(); await sleep(120);
  check("a command with no text keeps its field hidden", d.querySelector(".modal .ctext").closest(".fmrow").style.display, "none");
  const cact = d.querySelector(".modal .caction");
  cact.value = "text"; cact.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(60);
  check("choosing insert-text reveals the field", d.querySelector(".modal .ctext").closest(".fmrow").style.display, "");
  d.querySelector(".modal .ctext").value = ":)";
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("the action is saved with its text", [w.lastSaveForm.commands[0].action, w.lastSaveForm.commands[0].text], ["text", ":)"]);
  d.querySelectorAll("#cmdbody .rdel")[0].click(); await sleep(300);
  check("a command can be deleted", w.lastSaveForm.commands.length, 2);

  check("the file buttons are gone", [!!d.getElementById("lists_export"), !!d.getElementById("lists_import")], [false, false]);
  const cfilter = d.getElementById("cmd_filter");
  cfilter.value = "paragraph"; cfilter.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(80);
  check("the finder narrows the chips down", d.querySelectorAll("#cmdbody .cmdchip").length, 1);
  check("and shows its cross", d.getElementById("cmd_filter_clear").style.display, "");
  cfilter.value = "zzz"; cfilter.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(80);
  check("a fruitless search says so with the words used", d.querySelector("#cmdbody .ruleempty").textContent, "Nothing found: “zzz”");
  d.getElementById("cmd_filter_clear").click(); await sleep(80);
  check("the cross brings every chip back", d.querySelectorAll("#cmdbody .cmdchip").length, 2);
  check("and the not-found note is gone", !!d.querySelector("#cmdbody .ruleempty"), false);

  const pb = d.getElementById("profbody");
  check("prompts listed", pb.querySelectorAll("input.profcb").length, 2);

  const pencil = () => pb.querySelector('button[data-a="edit"]');
  pencil().click(); await sleep(150);
  check("the prompt editor opens as a titled dialog", d.querySelector(".modal .mtitle").textContent, "Editing the prompt");
  check("it carries name, prompt and a try-it field", [!!d.querySelector(".modal .pfname"), !!d.querySelector(".modal .pfprompt"), !!d.querySelector(".modal .pfsample")], [true, true, true]);
  d.querySelector(".modal .pfname").value = "Cleanup+";
  d.querySelector(".modal .pfprompt").value = "tidy the text";
  d.querySelector(".modal .btn.yes").click(); await sleep(350);
  check("saving keeps the new name and prompt", [w.lastSaveForm.profiles[0].name, w.lastSaveForm.profiles[0].prompt], ["Cleanup+", "tidy the text"]);
  check("and the list shows the new name", pb.querySelector(".prow .pnm").textContent, "Cleanup+");
  pencil().click(); await sleep(150);
  d.querySelector(".modal .btn.ghost").click(); await sleep(200);
  check("cancelling changes nothing", w.lastSaveForm.profiles[0].name, "Cleanup+");
  check("every prompt has a handle to drag it by", pb.querySelectorAll(".prow .grip").length, 2);
  const beforeOrder = w.lastSaveForm.profiles.map(p => p.id).join("+");
  const plist = pb.querySelector(".plist");
  const prows = [...pb.querySelectorAll(".prow")];
  prows.forEach((r, i) => { r.getBoundingClientRect = () => ({ top: 100 + i * 40, bottom: 140 + i * 40, height: 40, left: 0, width: 600 }); });
  plist.getBoundingClientRect = () => ({ top: 100, bottom: 180, height: 80, left: 0, width: 600 });
  prows[0].querySelector(".grip").dispatchEvent(new w.MouseEvent("pointerdown", { bubbles: true, clientY: 110, button: 0 }));
  check("dragging lifts a copy and shows where it lands", [d.querySelectorAll(".prow.ghost").length, d.querySelectorAll(".dropline").length], [1, 1]);
  d.dispatchEvent(new w.MouseEvent("pointermove", { bubbles: true, clientY: 175 }));
  d.dispatchEvent(new w.MouseEvent("pointerup", { bubbles: true })); await sleep(350);
  check("dropping swaps the order and saves it", w.lastSaveForm.profiles.map(p => p.id).join("+") !== beforeOrder, true);
  check("and the helpers are cleaned up", [d.querySelectorAll(".prow.ghost").length, d.querySelectorAll(".dropline").length], [0, 0]);
  d.getElementById("profadd").click(); await sleep(150);
  check("adding opens the same dialog, titled anew", d.querySelector(".modal .mtitle").textContent, "New prompt");
  d.querySelector(".modal .btn.ghost").click(); await sleep(200);

  const pdel = () => pb.querySelector('button[data-a="pdel"]');
  const profsBefore = w.lastSaveForm.profiles ? w.lastSaveForm.profiles.length : 2;
  pdel().click(); await sleep(200);
  check("deleting a prompt asks first", !!d.querySelector(".modal-bg"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no keeps the prompt", pb.querySelectorAll("input.profcb").length, 2);
  pdel().click(); await sleep(200);
  d.querySelector(".modal .btn.yes").click(); await sleep(400);
  check("saying yes removes it and saves at once", w.lastSaveForm.profiles.length, profsBefore - 1);

  tab("post"); await sleep(30);
  check("the local model is summed up in the card", !!d.querySelector("#p-post #llm_sum"), true);
  check("and the catalog waits behind a button", !!d.querySelector("#p-post #llm_catalog"), true);
  check("the summary names the model in use", [...d.querySelectorAll("#llm_sum .sumv")][0].textContent, "model.gguf");
  check("nothing of it is left on the languages page", !!d.querySelector("#p-models #llm_catalog"), false);
  d.getElementById("llm_catalog").click(); await sleep(200);
  check("the catalog opens as a titled dialog", d.querySelector(".modal .mtitle").textContent, "Model catalog");
  check("installed models are listed with a switch each", d.querySelectorAll("#proc-models input.llmpick").length, 1);
  check("the model in use has its switch on", d.querySelector("#proc-models input.llmpick").checked, true);
  check("the installed models live under their own label", d.querySelector("#proc-models .blklbl").textContent, "Installed models");
  check("a warm model says it sits in memory", d.querySelector("#proc-models .mrow .mstate").textContent, "in memory");
  const ejectBtn = () => d.querySelector('#proc-models .mrow button[aria-label="Unload from memory"]');
  check("and offers to free it", !!ejectBtn() && !ejectBtn().disabled, true);
  ejectBtn().click(); await sleep(250);
  check("ejecting asks the program to drop the model", w.llmUnloads, 1);
  check("the row then says the model is only on disk", d.querySelector("#proc-models .mrow .mstate").textContent, "on disk");
  check("and the eject button goes quiet", ejectBtn().disabled, true);
  check("nothing is searched yet, and the count says so", d.getElementById("hf_count").textContent, "no search yet");
  check("the empty result area invites a search", !!d.querySelector("#hf_results .searchempty"), true);
  const del = d.querySelector("#proc-models button[data-a=\"ldel\"]");
  check("active LLM model can be deleted", !!del, true);
  del.click(); await sleep(60);
  check("deleting a model asks first", !!d.querySelector(".modal-bg"), true);
  const lastModal = () => [...d.querySelectorAll(".modal")].pop();
  check("the question is asked in the app style, not by the browser", lastModal().querySelectorAll(".btn").length, 2);
  check("the way out comes first, the action second", [...lastModal().querySelectorAll(".btn")].map(b=>b.className), ["btn ghost", "btn yes"]);
  check("the question keeps the focus on the way out", d.activeElement.className, "btn ghost");
  check("the question is a dialog for the reader", d.querySelector(".modal").getAttribute("aria-modal"), "true");
  check("the window behind is out of reach while it is open", !!d.querySelector(".content[inert]"), true);
  d.dispatchEvent(new w.KeyboardEvent("keydown", { key: "Tab", bubbles: true }));
  check("Tab moves between the two answers and no further", d.activeElement.className, "btn yes");
  d.dispatchEvent(new w.KeyboardEvent("keydown", { key: "Escape", bubbles: true }));
  await sleep(150);
  check("Escape closes it as a no", d.querySelectorAll(".modal").length, 1);
  check("and the catalog behind it stays open", !!d.querySelector(".modal.llmcat"), true);
  del.click(); await sleep(150);
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("the question closes with the answer", d.querySelectorAll(".modal").length, 1);
  check("model list empty after delete", d.querySelectorAll("#proc-models input.llmpick").length, 0);
  d.getElementById("hf_q").value = "qwen";
  d.getElementById("hf_go").click(); await sleep(300);
  check("searching Hugging Face is a button, not a decoration", d.getElementById("hf_go").tagName, "BUTTON");
  check("a finished search counts what it found", d.getElementById("hf_count").textContent, "found 1");
  check("the search button stands next to the field, not inside it", !d.querySelector(".srchbox").contains(d.getElementById("hf_go")), true);
  check("memory is spelled out in words", d.getElementById("hf_ramline").textContent, "Memory available: 8.8 GB of 16 GB");
  check("and the colour key sits on its own line", d.getElementById("hf_fitline").textContent, "●fits●tight●no RAM");
  const repo = d.querySelector(".hfrepo");
  check("a found repository is a button too", !!repo, true);
  check("and it says whether it is open", repo.getAttribute("aria-expanded"), "false");
  repo.click(); await sleep(300);
  check("opening it lists the files that fit this computer", d.querySelectorAll('#hf_results button[data-repo]').length, 1);
  check("and the button now says it is open", d.querySelector(".hfrepo").getAttribute("aria-expanded"), "true");
  check("the too-big files are counted, not silently dropped", d.getElementById("hf_results").textContent.includes("hidden: 1"), true);
  const fitBox = d.getElementById("hf_fit");
  check("the fit filter is a visible choice", !!fitBox && fitBox.checked, true);
  fitBox.checked = false; fitBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(100);
  check("turning it off shows everything", d.querySelectorAll('#hf_results button[data-repo]').length, 2);
  fitBox.checked = true; fitBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(100);
  d.querySelector(".modal.llmcat .btn.ghost").click(); await sleep(200);
  check("closing the catalog leaves the card behind", [!!d.querySelector(".modal.llmcat"), !!d.getElementById("llm_sum")], [false, true]);


  tab("models"); await sleep(120);
  const pickAbsent = async () => { const s = d.querySelector('#langlist .pcard[data-id="base"] input.psw'); s.checked = true; s.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(80); };
  check("the list remembers which language was open", !!d.querySelector('#langlist .pcard[data-id="base"]'), true);
  await pickAbsent();
  check("picking a model that is not here asks first", !!d.querySelector(".modal-bg"), true);
  check("the question wears a title", d.querySelector(".modal .mtitle").textContent, "Downloading a model");
  check("the question names the model and its size", d.querySelector(".modal p").textContent.includes("Base") && d.querySelector(".modal p").textContent.includes("142 MB"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no downloads nothing", w.dlCalls.length, 0);
  check("saying no keeps the old choice", w.lastSaveForm.lang_models.ru, "small");

  await pickAbsent();
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("saying yes starts the download", w.dlCalls, ["base"]);
  check("agreeing pins the language to the model at once", w.lastSaveForm.lang_models.ru, "base");
  check("the running download is announced above the languages", d.querySelector("#langlist .dlline").textContent.includes("Base"), true);
  check("a download can be stopped", !!d.querySelector('#langlist button[data-a="cancel"][data-id="base"]'), true);
  d.querySelector('#langlist button[data-a="cancel"][data-id="base"]').click(); await sleep(250);
  check("stopping asks the program to stop", w.cancelCalls, ["base"]);
  const dlIcon = () => d.querySelector('#langlist button[data-a="dl"][data-id="base"]');
  check("a stopped download offers the arrow again", !!dlIcon(), true);
  const dlBefore = w.dlCalls.length;
  const ruBefore = w.lastSaveForm.lang_models.ru;
  dlIcon().click(); await sleep(250);
  check("the arrow downloads without choosing", w.dlCalls.length, dlBefore + 1);
  check("and the language keeps its model", w.lastSaveForm.lang_models.ru, ruBefore);
  w.finishDl("base"); await sleep(1400);
  check("and the program says the model is ready", d.getElementById("st_saved").textContent, "Model downloaded");

  const activeDel = () => d.querySelector('#langlist .pcard[data-id="small"] button[data-a="del"]');
  check("the model in use can be removed too — that is the way out of a full disk", !!activeDel(), true);
  activeDel().click(); await sleep(150);
  check("removing the model in use warns what it costs", d.querySelector(".modal p").textContent.includes("Recognition stops"), true);
  d.querySelector(".modal .btn.yes").click(); await sleep(300);
  check("and the program is told this was meant", w.delCalls[w.delCalls.length - 1], ["small", true]);
  tab("about"); await sleep(60);
  check("about section shown", shown("about"), true);
  check("about keeps the version card to itself", d.querySelectorAll("#p-about .card").length, 1);
  tab("help"); await sleep(60);
  check("the guide lives on its own page", shown("help") && d.querySelectorAll("#p-help .card").length === 1, true);
  tab("contacts"); await sleep(60);
  check("and so does the author card", shown("contacts") && d.querySelectorAll("#p-contacts .card").length === 2, true);
  check("the prompts stand on their own page too", !!d.querySelector("#p-post #profbody"), true);
  check("out of the expert fold", !!d.querySelector("#p-post .card[data-adv]"), false);
  tab("post"); await sleep(60);
  const master = d.getElementById("post_enabled");
  check("post-processing has one master switch", master.type, "checkbox");
  check("with no model left the heading warns about it", d.getElementById("post_warn").textContent, "on, but no model is picked");
  master.checked = false; master.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("turning it off is saved", w.lastSaveForm.post_enabled, false);
  check("and the cards below dim out", d.getElementById("post_model_card").className.includes("offdim") && d.getElementById("post_prompts_card").className.includes("offdim"), true);
  master.checked = true; master.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("turning it back lights them up", d.getElementById("post_model_card").className.includes("offdim"), false);
  check("the local model is the source out of the box", d.getElementById("src_local").className.includes("on"), true);
  check("the external card is marked as the idle one", d.getElementById("src_api").className.includes("idle"), true);
  const flipSrc = async (id) => { const s = d.getElementById(id); s.checked = true; s.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200); };
  check("each source card carries a switch", d.querySelectorAll("#p-post input.srcpick").length, 2);
  check("the local one is on out of the box", d.getElementById("pick_local").checked, true);
  await flipSrc("pick_api");
  check("flipping the switch moves the work to the server", w.lastSaveForm.post_source, "api");
  check("an unset server is flagged beside the heading", d.getElementById("post_warn").textContent, "on, but the server is not set up");
  await flipSrc("pick_local");
  check("back on local the warning names the missing model", d.getElementById("post_warn").textContent, "on, but no model is picked");
  await flipSrc("pick_api");
  check("and the marks trade places", d.getElementById("src_api").className.includes("on") && d.getElementById("src_local").className.includes("idle"), true);
  check("the switches follow the choice", [d.getElementById("pick_api").checked, d.getElementById("pick_local").checked], [true, false]);
  await flipSrc("pick_local");
  check("switching back is one flick too", w.lastSaveForm.post_source, "local");
  check("the external server is summed up, not spelled out", !!d.getElementById("api_sum"), true);
  check("while nothing is set the card says so", d.querySelector("#api_sum .sumv").textContent, "not set up — post-processing runs locally");
  check("and the button offers to set it up", d.getElementById("api_edit").textContent, "Set up");
  check("no warning shows while it is empty", d.getElementById("postapi_warn").textContent, "");
  d.getElementById("api_edit").click(); await sleep(150);
  check("the settings live in a titled dialog", d.querySelector(".modal .mtitle").textContent, "External server");
  d.querySelector(".modal .apiurl").value = "https://api.example.com/v1";
  d.querySelector(".modal .apimodel").value = "gpt-4.1-mini";
  d.querySelector(".modal .apikey").value = "sk-secret";
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("pointing the prompts outward asks first", !!d.querySelector(".modal-bg"), true);
  check("and the question names the address", d.querySelector(".modal p").textContent.includes("api.example.com"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no leaves the card untouched", [w.lastSaveForm.post_api_url, w.postKeys.length], ["", 0]);
  d.getElementById("api_edit").click(); await sleep(150);
  d.querySelector(".modal .apiurl").value = "https://api.example.com/v1";
  d.querySelector(".modal .apimodel").value = "gpt-4.1-mini";
  d.querySelector(".modal .apikey").value = "sk-secret";
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  d.querySelector(".modal .btn.yes").click(); await sleep(400);
  check("saying yes applies the address", w.lastSaveForm.post_api_url, "https://api.example.com/v1");
  check("the key goes to the program, not into the config form", w.postKeys, ["sk-secret"]);
  check("no later save carries the key along", Object.keys(w.lastSaveForm).includes("post_api_key"), false);
  check("the honest warning appears", d.getElementById("postapi_warn").textContent.includes("Recognized text"), true);
  check("the summary now lists what is set", [...d.querySelectorAll("#api_sum .sumv")].map(e=>e.textContent), ["https://api.example.com/v1", "gpt-4.1-mini", "key saved", "30 s"]);
  check("and the warning by the heading clears itself", d.getElementById("post_warn").style.display, "none");
  const pcb = () => d.querySelector("#profbody input.profcb");
  pcb().checked = false; pcb().dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("with no prompt checked the heading says the work is idle", d.getElementById("post_warn").textContent, "on, but no prompt is checked");
  pcb().checked = true; pcb().dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("checking one back clears it", d.getElementById("post_warn").style.display, "none");
  check("and the button switches to changing it", d.getElementById("api_edit").textContent, "Change");
  check("setting an address puts the work on that server", w.lastSaveForm.post_source, "api");
  check("the card offers a test beside Change, to its left", [...d.querySelectorAll("#src_api .acts button")].map(b => b.id), ["api_test", "api_edit"]);
  w.setPostTestOK(false);
  d.getElementById("api_test").click(); await sleep(300);
  check("a failed test says why, under the honest warning", d.getElementById("postapi_err").textContent, "no credits left");
  check("and the test goes to the address that is set", w.postTests[w.postTests.length - 1][0], "https://api.example.com/v1");
  w.setPostTestOK(true);
  d.getElementById("api_test").click(); await sleep(300);
  check("an answering server clears the red line", d.getElementById("postapi_err").textContent, "");
  w.setPostErr("Cleanup: no credits left");
  await w.refreshState(); await sleep(60);
  check("a failure during dictation shows up there too", d.getElementById("postapi_err").textContent, "Cleanup: no credits left");
  w.setPostErr("");
  await w.refreshState(); await sleep(60);
  check("and a later good run wipes it", d.getElementById("postapi_err").textContent, "");
  d.getElementById("api_edit").click(); await sleep(150);
  check("the dialog keeps the test out of it", !!d.getElementById("fm_test"), false);
  check("a saved key shows as dots, so it is clear one is there", d.querySelector(".modal .apikey").value, "••••••••••••");
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  d.getElementById("api_edit").click(); await sleep(150);
  check("a saved key can be deleted from the dialog", d.getElementById("fm_extra").textContent, "Delete the key");
  d.querySelector(".modal .btn.ghost").click(); await sleep(200);
  d.getElementById("api_edit").click(); await sleep(150);
  d.querySelector(".modal .apiurl").value = "";
  d.querySelector(".modal .btn.yes").click(); await sleep(350);
  check("clearing the address needs no question", [!!d.querySelector(".modal-bg"), w.lastSaveForm.post_api_url], [false, ""]);

  const before = w.saveCalls;
  const sw = d.getElementById("auto_enter");
  sw.checked = !sw.checked; sw.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(220);
  check("a toggle applies itself, no Save needed", w.saveCalls > before, true);

  check("no permanent mode line in the status bar", !!d.getElementById("st_level"), false);
  check("no switching from the status bar", !!d.getElementById("st_levelbtn"), false);
  tab("contacts"); await sleep(30);
  check("the contacts page carries the mail address", d.querySelector("#p-contacts .val").textContent, "holdtotype@outlook.com");
  check("about carries the repository button", !!d.querySelector('#p-about button[onclick="appRepoLink()"]'), true);
  tab("history"); await sleep(30);
  const skipNew = d.getElementById("hist_skip_new");
  const skipStore = d.getElementById("history_skip");
  const savesBeforeSkip = w.saveCalls;
  skipNew.value = "game.exe";
  d.getElementById("hist_skip_add").click(); await sleep(220);
  check("a skipped program is added with a button, not a comma", skipStore.value, "game.exe");
  check("and it appears as a chip", [...d.querySelectorAll("#hist_skip_list .skipchip")].length, 1);
  check("and the change saves itself", w.saveCalls > savesBeforeSkip, true);
  d.querySelector("#hist_skip_list .chipx").click(); await sleep(220);
  check("the chip's cross takes it away", skipStore.value, "");

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
  check("and the heading itself is what is highlighted", !!d.querySelector("#p-text .blklbl.hit"), true);
  omni.value = ""; omni.dispatchEvent(new w.Event("input")); await sleep(60);

  check("no save button left", !!d.querySelector(".footer"), false);

  const drag0 = w.dragCalls;
  const down = (el) => el.dispatchEvent(new w.MouseEvent("mousedown", { bubbles: true, button: 0 }));
  down(d.getElementById("omni"));
  down(d.querySelector(".cap"));
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
  check("a missing model is not green", d.getElementById("state_active_led").className.includes("on"), false);
  check("a missing model warns", d.getElementById("state_active_led").className.includes("warn"), true);
  check("a missing model offers to download it", d.getElementById("state_active_btn").textContent, "Download");
  check("a missing model is named in the list of models in use", d.querySelector("#state_assigned .arow .amiss").textContent, "not installed");
  w.setModelStates("ready", "ready"); await sleep(1700);
  check("an installed model is green", d.getElementById("state_active_led").className.includes("on"), true);

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
  check("the wizard proposes the default model, Whisper Medium", d.querySelectorAll("#wiz_plan .advrow").length, 1);
  check("and names it", d.querySelector("#wiz_plan .advname").textContent.includes("Medium"), true);
  check("the download button carries the size", d.getElementById("wiz_dl").textContent.includes("539 MB"), true);
  check("the download can be put off", d.getElementById("wiz_dl_skip").style.display, "");
  const wizDlBefore = w.dlCalls.length;
  d.getElementById("wiz_dl").click(); await sleep(200);
  check("the wizard downloads what the plan is missing", w.dlCalls.length > wizDlBefore, true);
  check("the wizard shows the download running", d.getElementById("wiz_dlrow").style.display, "");
  check("the wizard does not let you walk past a download", d.getElementById("wiz_next").disabled, true);
  w.finishDl("medium-q5_0"); await sleep(1200);
  check("a finished download opens the way on", d.getElementById("wiz_next").disabled, false);

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

  w.setModelState("medium-q5_0", "absent");
  w.wizStart(); await sleep(150);
  d.getElementById("wiz_next").click(); await sleep(250);
  check("the model step blocks again while nothing is installed", d.getElementById("wiz_next").disabled, true);
  d.getElementById("wiz_dl_skip").click(); await sleep(100);
  check("skipping the download opens the way on", d.getElementById("wiz_next").disabled, false);
  check("and warns that dictation needs a model", d.getElementById("wiz_skipnote").style.display, "");
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
    ["input[type=text],input[type=number],input[type=password],select{", "border-radius:calc(var(--r) * .55)"],
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
    [".modal-bg{", "background:var(--scrim)"],
    [".content{", "scrollbar-gutter:stable both-edges"],
    [".hotkey-val{", "background:var(--keybg)"],
    [".scard .led.on{", "background:var(--ok)"],
    [".miclevel i{", "background:var(--soft)"],
    [".miclevel i{", "width:var(--lvlw,4px)"],
    [".miclevel i.on{", "background:var(--hi)"],
    [".wizlvl i.on{", "background:var(--hi)"],
    [".miclevel.grow{", "width:100%"],
    [".mock-dot{", "background:var(--rec)"],
    ["button.cap.close:hover{", "background:var(--badbg)"],
    ["button.mini.danger:hover{", "border-color:var(--badline)"],
    [".toast{", "text-shadow:var(--amberglow)"],
    ["button.iconbtn.danger:hover{", "filter:var(--badfilter)"],
    ["button.mini{", "border:1px solid var(--btnline)"],
    ["button.mini::before{", "content:var(--btnbo)"],
    ["button.mini::after{", "content:var(--btnbc)"],
    ["button.btn::before{", "content:var(--btnbo)"],
    ["button.btn::after{", "content:var(--btnbc)"],
  ];
  for (const [sel, want] of skinned) {
    const at = css.indexOf(sel);
    const rule = at < 0 ? "" : css.slice(at, css.indexOf("}", at));
    check(`${sel} follows the skin (${want})`, rule.includes(want), true);
  }
  const fieldRules = css.match(/(?:input[type=text]|.row select|.rulerow input|.replrow input|.advq select|.wizrow select)[^}]*{[^}]*}/g) || [];
  const strays = fieldRules.filter(r => /(?:^|;|{)(?:padding|font-size):s*(?!var(--fieldpad)|var(--ctlfs))/.test(r));
  check("no field sets its own height", strays, []);
  const literals = [...new Set((css.match(/#[0-9a-f]{6}/g) || []))].sort();
  check("no colour is written into the stylesheet by hand", literals, []);
  check("no face is nailed to Consolas outside the skin", /font:[^;]*Consolas/.test(css), false);
  check("native dropdowns follow the skin", css.includes("color-scheme:var(--scheme,dark)"), true);
  check("and no dark scheme is nailed down", /color-scheme:s*dark/.test(css), false);

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
