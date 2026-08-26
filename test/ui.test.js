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
        status_line: "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free",
        active_model: "GigaAM v3", active_state: ruState, active_lang: "RU",
        assigned: [
          { model: "GigaAM v3", langs: "RU", state: ruState, current: true },
          { model: "Small", langs: "Detect itself, EN", state: otherState, current: false },
        ],
        installed_models: ["Small", "my-model"],
        llm_ok: true, mic_ok: true, last_at: lastAt, last_app: "chrome.exe",
        remote: remote, backend_err: backendErr,
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
    window.appSetPostKey = async (k) => { window.postKeys.push(k); postKeySet = !!k; return JSON.stringify({ ok: true, severity: "ok", message: "Saved" }); };
    window.appPostKeySet = async () => postKeySet;
    window.folderOpens = 0;
    window.appOpenModelsFolder = () => { window.folderOpens++; };
    window.linkOpens = [];
    window.appModelLink = (id) => { window.linkOpens.push(id); };
    window.appModels = async () =>
      JSON.stringify([
        { id: "base", name: "Base", desc: "fast", size: 142, state: modelStates.base, pct: 12, engine: "whisper", langs: "*", auto: true, translate: true, speed: 5, accuracy: 2 },
        { id: "small", name: "Small", desc: "balanced", size: 466, state: modelStates.small, engine: "whisper", langs: "*", auto: true, translate: true, serves: ["auto"], loaded: true, speed: 3, accuracy: 3 },
        { id: "medium-q5_0", name: "Medium (q5)", desc: "recommended", size: 539, state: modelStates["medium-q5_0"], pct: 5, engine: "whisper", langs: "*", auto: true, translate: true, speed: 2, accuracy: 4 },
        { id: "gigaam-v3", name: "GigaAM v3", desc: "russian", size: 232, state: modelStates["gigaam-v3"], pct: 5, engine: "sherpa", langs: "ru", punct: true, serves: ["ru"], speed: 5, accuracy: 5 },
        { id: "moonshine-uk", name: "Moonshine Base uk", desc: "ukrainian", size: 135, state: modelStates["moonshine-uk"], engine: "sherpa", langs: "uk", manual: true, link: "https://example.com/moonshine", speed: 5, accuracy: 3 },
        { id: "local:my-model", name: "my-model", desc: "found in the models folder", size: 900, state: "installed", engine: "whisper", langs: "*", custom: true },
      ]);
    window.appLLMSearch = async () =>
      JSON.stringify({ repos: [{ id: "org/Repo-GGUF", downloads: 1234, updated: "2026-01-01" }] });
    window.appLLMFiles = async () =>
      JSON.stringify({ files: [{ file: "q4.gguf", size: 4000, fit: "ok", need: 6166 }, { file: "q8.gguf", size: 9000, fit: "bad", need: 13000 }] });
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
  for (const id of ["hotkey", "pause_hotkey", "tr_hotkey"]) {
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
  check("the models badge explains itself", d.getElementById("badge_models").title, "Installed models");
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
    ["state", "system", "mic", "history", "models", "text", "dictation", "post", "about", "help", "contacts"]);
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
    "type_mode", "overlay", "overlay_position", "overlay_text", "animation", "mic_device", "mic_bar", "beep", "sound_theme",
    "models", "presets", "munload", "language", "threads", "punctuation", "whisper_prompt", "profbody", "dict_model", "tr_engine", "mlang", "mfind",
    "tr_default", "translate_target", "translate_ask", "translate_ask_seconds", "tr_hotkey",
    "tl_en", "ui_language", "upd_check", "check_updates", "server_autostart", "server_port",
    "server_exe", "server_url", "proc-models", "proc-search", "ver2", "autorun", "pause_hotkey",
  ];
  const missing = everySetting.filter((id) => !d.getElementById(id));
  check("every setting is present in the new window", missing, []);

  tab("models"); await sleep(80);
  check("models section shown", shown("models"), true);
  check("recognition models listed", d.querySelectorAll("#models .mcard").length, 6);
  check("no radio buttons left in the library", d.querySelectorAll('#models input[type="radio"]').length, 0);
  check("the list is split by state, not by language", [...d.querySelectorAll("#models .mslot")].map(h=>h.dataset.slot), ["inst", "avail"]);
  check("the installed shelf holds what is on disk", [...d.querySelectorAll('#models .mcard[data-slot="inst"]')].map(r=>r.dataset.id), ["small", "local:my-model"]);
  check("the rest wait under available", [...d.querySelectorAll('#models .mcard[data-slot="avail"]')].map(r=>r.dataset.id), ["base", "medium-q5_0", "gigaam-v3", "moonshine-uk"]);
  check("the active model wears its pill", [...d.querySelectorAll("#models .mpill.on")].length, 1);
  check("and its whole card is lit", [...d.querySelectorAll("#models .mcard.on")].map(r=>r.dataset.id), ["small"]);
  check("a card assigned to a language says which", d.querySelector('#models .mcard[data-id="gigaam-v3"] .mpill').textContent, "RU");
  check("every card measures itself in two bars", d.querySelectorAll("#models .mcard .mbar").length, 10);
  check("the bars are filled to the model, not all alike", [...d.querySelectorAll(String.raw`#models .mcard[data-id="base"] .mtrack i`)].map(i=>i.style.width), ["40%", "100%"]);
  check("a hand-copied model shows no bars — its powers are unknown", d.querySelectorAll('#models .mcard[data-id="local:my-model"] .mbar').length, 0);
  check("and wears a question mark instead of languages", d.querySelector('#models .mcard[data-id="local:my-model"] .mtag').textContent, "?");
  check("the licensed model offers no download button", !!d.querySelector('#models .mcard[data-id="moonshine-uk"] button[data-a="dl"]'), false);
  check("it explains the licence instead", d.querySelector('#models .mcard[data-id="moonshine-uk"]').textContent.includes("licence forbids"), true);
  const linkBtn = d.querySelector('#models .mcard[data-id="moonshine-uk"] button[data-a="link"]');
  check("and offers the source link", !!linkBtn, true);
  linkBtn.click(); await sleep(60);
  check("the link opens the model page", w.linkOpens, ["moonshine-uk"]);
  check("the models folder can be opened for hand-copied files", !!d.querySelector('#p-models button[onclick="appOpenModelsFolder()"]'), true);
  const eject = d.querySelector('#models .mcard[data-id="small"] button[data-a="unload"]');
  check("the loaded model offers to leave the memory", !!eject, true);
  eject.click(); await sleep(150);
  check("and the program is asked to unload it", w.unloadCalls, 1);
  check("there is a plain unload row as well", !!d.getElementById("munload"), true);
  check("the language filter is one control", [...d.getElementById("mlang").options].map(o=>o.value), ["all", "multi", "punct", "fit", "ru"]);
  check("the library head stays put while the shelf scrolls", w.getComputedStyle(d.querySelector("#p-models .libhead")).position, "sticky");
  check("every language owns a preset row", d.querySelectorAll("#presets .prow").length, 9);
  const autoSel = d.querySelector("#presets .prow select");
  check("the auto row offers only models that detect the language", [...autoSel.options].map(o=>o.value), ["base", "small", "medium-q5_0"]);
  const ruRow = d.querySelectorAll("#presets .prow")[1];
  check("the russian row is marked as the current language", ruRow.className.includes("cur"), true);
  check("the russian row holds its assigned model", ruRow.querySelector("select").value, "gigaam-v3");
  check("a model that only recognizes says so in the list", [...ruRow.querySelectorAll("option")].find(o=>o.value==="gigaam-v3").textContent.includes("recognition only"), true);
  const mfind = d.getElementById("mfind");
  mfind.value = "giga"; mfind.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(150);
  check("the name search narrows the list", [...d.querySelectorAll("#models .mcard:not(.hidden)")].map(r=>r.dataset.id), ["gigaam-v3"]);
  mfind.value = ""; mfind.dispatchEvent(new w.Event("input", { bubbles: true })); await sleep(150);
  const saveBefore = w.saveCalls;
  d.querySelector('#models .mcard[data-id="small"]').click(); await sleep(200);
  check("a click on an installed row is the choice", w.saveCalls > saveBefore, true);
  check("and it pins the model to the current language", w.lastSaveForm.lang_models.ru, "small");
  check("the mark carries both shapes", [d.querySelectorAll(".mk.mic").length, d.querySelectorAll(".mk.face").length], [2, 2]);
  w.applyThemeVars("soft");
  check("the soft design shows the face", [d.documentElement.style.getPropertyValue("--markmic"), d.documentElement.style.getPropertyValue("--markface")], ["none", "block"]);
  w.applyThemeVars("terminal:green");
  check("and every other one keeps the microphone", [d.documentElement.style.getPropertyValue("--markmic"), d.documentElement.style.getPropertyValue("--markface")], ["block", "none"]);

  const recLangs = [...d.getElementById("language").options].map(o=>o.value);
  check("italian can be dictated too", recLangs.includes("it"), true);
  check("ram estimate shown", d.querySelectorAll("#p-models .mram").length, 6);
  check("the routing table is gone", !!d.getElementById("routing"), false);
  check("engine tags rendered", d.querySelectorAll("#p-models .mtag").length >= 3, true);
  check("russian engine tagged RU", d.querySelector('#models .mcard[data-id="gigaam-v3"] .mtag').textContent, "RU");
  check("the language moved in with the dictation", !!d.querySelector("#p-dictation #language"), true);
  check("the threads moved in with the server", !!d.querySelector("#p-system #threads"), true);
  check("dictation names the model it uses", d.getElementById("dict_model").textContent, "Small");
  check("translation names its engine too", d.getElementById("tr_engine").textContent, "Translation is done by Small");
  check("the honest CPU line is on the page", d.getElementById("p-models").textContent.includes("S_CPU_LINE"), true);
  w.setModelState("small", "installed"); w.setModelState("gigaam-v3", "active");
  await w.refreshModels(); await sleep(60);
  check("a model that cannot translate says Whisper will step in", d.getElementById("tr_engine").textContent.includes("does not translate"), true);
  w.setModelState("small", "active"); w.setModelState("gigaam-v3", "absent");
  await w.refreshModels(); await sleep(60);
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
  const dots = d.querySelectorAll("#ovscheme .ovdot");
  check("the plate offers eight scheme positions", dots.length, 8);
  check("it starts at the bottom of the screen", d.querySelector('#ovscheme .ovdot[data-pos="bottom"]').className.includes("on"), true);
  check("the pause hotkey is shown", d.getElementById("pause_hotkey").textContent, "ctrl+alt+p");
  d.getElementById("pause_clear").click(); await sleep(300);
  check("clearing the pause hotkey is saved", w.lastSaveForm.pause_hotkey, "");
  d.querySelector('#ovscheme .ovdot[data-pos="top-right"]').click(); await sleep(300);
  check("a dot click moves the plate and saves itself", w.lastSaveForm.overlay_position, "top-right");
  check("the miniature follows the chosen dot", d.getElementById("ovmini").style.left, "88%");
  const caretBox = d.getElementById("ov_caret");
  caretBox.checked = true; caretBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("following the cursor is the ninth choice", w.lastSaveForm.overlay_position, "caret");
  check("and the scheme steps aside while it is on", d.getElementById("ovscheme").className.includes("off"), true);
  caretBox.checked = false; caretBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  const ovbox = d.getElementById("ovscheme");
  const mini = d.getElementById("ovmini");
  ovbox.getBoundingClientRect = () => ({ left: 0, top: 0, width: 176, height: 99 });
  mini.dispatchEvent(new w.Event("pointerdown"));
  mini.dispatchEvent(new w.MouseEvent("pointermove", { clientX: 44, clientY: 25 }));
  mini.dispatchEvent(new w.Event("pointerup")); await sleep(300);
  check("dragging the miniature makes the spot its own", w.lastSaveForm.overlay_position, "custom");
  check("and the fraction is remembered for this screen", w.lastSaveForm.overlay_custom["1920x1080"], { x: 0.25, y: 0.253 });
  d.querySelector('#ovscheme .ovdot[data-pos="bottom"]').click(); await sleep(300);
  check("with two monitors the screen choice appears", d.getElementById("ovmon_row").style.display, "");
  check("the cursor screen is the first offer", d.getElementById("overlay_monitor").options[0].textContent, "The screen with the cursor");
  check("each monitor is named by its resolution", d.getElementById("overlay_monitor").options[1].textContent.includes("1920×1080"), true);
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

  tab("dictation"); await sleep(30);
  check("the old turbo warning is gone — honesty lives in the engine line", !!d.getElementById("tr_warn"), false);
  check("the target language carries the honest note", !!d.querySelector("#p-dictation .row label .sub.warn"), true);
  check("translation lives with the other controls now", !!d.querySelector("#p-dictation #translate_target"), true);
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
  const rlang = rep.querySelector(".rlang");
  check("a rule can be pinned to a language", !!rlang, true);
  check("but serves every language by default", rlang.value, "");
  rlang.value = "ru"; rlang.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(250);
  check("the pinned language is saved with the rule", w.lastSaveForm.replacements[0].lang, "ru");
  rlang.value = ""; rlang.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(250);
  check("the dictionary card says only Whisper reads it", d.getElementById("dict_whisper_note").textContent.includes("S_DICT_WHISPER_ONLY"), true);
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
  check("opening it lists the files that fit this computer", d.querySelectorAll('#hf_results button[data-repo]').length, 1);
  check("and the button now says it is open", d.querySelector(".hfrepo").getAttribute("aria-expanded"), "true");
  check("the too-big files are counted, not silently dropped", d.getElementById("hf_results").textContent.includes("hidden: 1"), true);
  const fitBox = d.getElementById("hf_fit");
  check("the fit filter is a visible choice", !!fitBox && fitBox.checked, true);
  fitBox.checked = false; fitBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(100);
  check("turning it off shows everything", d.querySelectorAll('#hf_results button[data-repo]').length, 2);
  fitBox.checked = true; fitBox.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(100);


  tab("models"); await sleep(120);
  const pickAbsent = () => d.querySelector('#models .mcard[data-id="base"]');
  pickAbsent().click(); await sleep(80);
  check("picking a model that is not here asks first", !!d.querySelector(".modal-bg"), true);
  check("the question names the model and its size", d.querySelector(".modal p").textContent.includes("Base") && d.querySelector(".modal p").textContent.includes("142 MB"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no downloads nothing", w.dlCalls.length, 0);
  check("saying no puts the choice back", !!d.querySelector('#models .mcard[data-id="small"] .mpill.on'), true);

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
  check("agreeing pins the language to the model at once", w.lastSaveForm.lang_models.ru, "base");
  w.finishDl("base"); await sleep(1400);
  check("and the program says the model is ready", d.getElementById("st_saved").textContent, "Model downloaded");

  const activeDel = () => d.querySelector('#models .mcard[data-id="small"] button[data-a="del"]');
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
  check("about keeps the version card to itself", d.querySelectorAll("#p-about .card").length, 1);
  tab("help"); await sleep(60);
  check("the guide lives on its own page", shown("help") && d.querySelectorAll("#p-help .card").length === 1, true);
  tab("contacts"); await sleep(60);
  check("and so does the author card", shown("contacts") && d.querySelectorAll("#p-contacts .card").length === 2, true);
  check("the prompts stand on their own page too", !!d.querySelector("#p-post #profbody"), true);
  check("out of the expert fold", !!d.querySelector("#p-post .card[data-adv]"), false);
  tab("post"); await sleep(60);
  const pu = d.getElementById("post_api_url");
  check("an external post-processing server can be set", !!pu, true);
  check("but nothing is filled in by default", pu.value, "");
  check("and no warning shows while it is empty", d.getElementById("postapi_warn").textContent, "");
  pu.value = "https://api.example.com/v1"; pu.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("pointing the prompts outward asks first", !!d.querySelector(".modal-bg"), true);
  check("and the question names the address", d.querySelector(".modal p").textContent.includes("api.example.com"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no puts the address back", pu.value, "");
  pu.value = "https://api.example.com/v1"; pu.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  d.querySelector(".modal .btn.yes").click(); await sleep(350);
  check("saying yes applies the address", w.lastSaveForm.post_api_url, "https://api.example.com/v1");
  check("and the honest warning appears", d.getElementById("postapi_warn").textContent.includes("Recognized text"), true);
  const pk = d.getElementById("post_api_key_new");
  pk.value = "sk-secret";
  d.getElementById("postapi_keysave").click(); await sleep(200);
  check("the key goes to the program, not into the config form", w.postKeys, ["sk-secret"]);
  check("the field forgets the key at once", pk.value, "");
  await sleep(100);
  check("and the row says a key is saved", d.getElementById("postapi_keystate").textContent, "key saved");
  check("no later save carries the key along", Object.keys(w.lastSaveForm).includes("post_api_key"), false);
  pu.value = ""; pu.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("clearing the address needs no question", [!!d.querySelector(".modal-bg"), w.lastSaveForm.post_api_url], [false, ""]);
  check("the app rules moved in with the other rules", !!d.querySelector("#p-text #rulesbody"), true);

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
  check("and the heading itself is what is highlighted", !!d.querySelector("#p-text .sect.hit"), true);
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
    ["input[type=text],input[type=number],select{", "border-radius:calc(var(--r) * .55)"],
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
  ];
  for (const [sel, want] of skinned) {
    const at = css.indexOf(sel);
    const rule = at < 0 ? "" : css.slice(at, css.indexOf("}", at));
    check(`${sel} follows the skin (${want})`, rule.includes(want), true);
  }
  const fieldRules = css.match(/(?:input[type=text]|.row select|.rulerow input|.replrow input|.replcheck input|.advq select|.wizrow select)[^}]*{[^}]*}/g) || [];
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
