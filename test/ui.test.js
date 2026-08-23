const fs = require("fs");
const path = require("path");
const { JSDOM } = require("jsdom");

const html = fs.readFileSync(path.join(__dirname, "page.html"), "utf8");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

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
        remote: remote,
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
    window.appModelCancel = async (id) => { window.cancelCalls.push(id); modelStates[id] = "absent"; return true; };
    window.appModels = async () =>
      JSON.stringify([
        { id: "base", name: "Base", desc: "fast", size: 142, state: modelStates.base, pct: 12, engine: "whisper", langs: "*" },
        { id: "small", name: "Small", desc: "balanced", size: 466, state: modelStates.small, engine: "whisper", langs: "*" },
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
      "appDoUpdate", "appReady", "appJSError", "appCopyLast",
    ]) {
      window[name] = () => {};
    }
    window.dragCalls = 0;
    window.appDrag = () => { window.dragCalls++; };
    window.saveCalls = 0;
    window.lastSave = {};
    window.appSave = async (json) => {
      window.saveCalls++;
      const f = JSON.parse(json);
      micBadge = f.mic_device_name ? f.mic_device_name.split(" ")[0] : "Realtek";
      const message = Number(f.server_port) === 8910 ? "Saved" : "Restarting the recognizer…";
      window.lastSaveForm = f;
      window.lastSave = { ok: true, severity: "ok", message };
      return JSON.stringify(window.lastSave);
    };
    window.replaceCalls = [];
    window.appTestReplace = async (t) => { window.replaceCalls.push(t); return t.replace(/git hub/gi, "GitHub"); };
    window.appModelDel = async () => "ok";
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
  check("status screen meter follows the microphone", d.getElementById("state_mic_bar").style.width !== "" && d.getElementById("state_mic_bar").style.width !== "0%", true);
  check("status bar names both models", d.getElementById("st_main").textContent, "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free");
  check("sidebar badges filled", [d.getElementById("badge_mic").textContent, d.getElementById("badge_models").textContent], ["Realtek", "2"]);
  check("status bar led lit", d.getElementById("st_led").classList.contains("on"), true);

  check("eight sections in the sidebar", d.querySelectorAll(".nav").length, 8);
  const everySetting = [
    "hotkey", "min_record_ms", "max_record_seconds", "auto_enter", "restore_clipboard",
    "type_mode", "overlay", "overlay_position", "overlay_text", "animation", "mic_device", "mic_bar", "beep", "sound_theme",
    "routing", "models", "language", "threads", "punctuation", "whisper_prompt", "profbody",
    "tr_default", "translate_target", "translate_ask", "translate_ask_seconds", "tr_hotkey",
    "tl_en", "ui_language", "upd_check", "check_updates", "server_autostart", "server_port",
    "server_exe", "server_url", "proc-models", "proc-search", "ver2", "autorun",
  ];
  const missing = everySetting.filter((id) => !d.getElementById(id));
  check("every setting is present in the new window", missing, []);

  tab("models"); await sleep(80);
  check("models section shown", shown("models"), true);
  check("recognition models listed", d.querySelectorAll('#p-models input[name="mdl"]').length, 3);
  check("model filters rendered", d.querySelectorAll(".fchip").length, 5);
  check("ram estimate shown", d.querySelectorAll("#p-models .mram").length, 3);
  check("routing panel rows", d.querySelectorAll("#routing .rrow").length, 3);
  check("routing shows engine", d.querySelectorAll("#routing .reng")[0].textContent, "gigaam-v3");
  check("engine tags rendered", d.querySelectorAll("#p-models .mtag").length, 3);
  check("russian engine tagged RU", d.querySelectorAll("#p-models .mtag")[2].textContent, "RU");

  tab("mic"); await sleep(120);
  const mic = d.getElementById("mic_device");
  check("microphone list has default plus devices", mic.options.length, 3);
  check("default option is localized", mic.options[0].textContent, "System default");
  check("system default selected initially", mic.value, "");
  mic.value = "dev1"; mic.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(30);
  check("microphone selection kept", mic.value, "dev1");
  await sleep(260);
  check("input level meter moves", d.getElementById("mic_bar").style.width !== "" && d.getElementById("mic_bar").style.width !== "0%", true);
  check("sidebar badge follows the chosen microphone", d.getElementById("badge_mic").textContent, "Headset");

  mic.value = ""; mic.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(30);
  tab("text"); await sleep(30);
  check("dictionary textarea present", !!d.getElementById("whisper_prompt"), true);
  check("punctuation modes offered", d.getElementById("punctuation").options.length, 3);

  tab("dictation"); await sleep(60);
  const ovpos = d.getElementById("overlay_position");
  check("the plate can be put in three places", ovpos.options.length, 3);
  check("it starts at the bottom of the screen", ovpos.value, "bottom");
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
  const autorun = d.getElementById("autorun");
  const autorunBefore = w.autorunCalls.length;
  autorun.checked = true; autorun.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("starting with Windows is a switch of its own", w.autorunCalls.length, autorunBefore + 1);
  check("it is not written into the config", w.autorunCalls[w.autorunCalls.length - 1], true);
  autorun.checked = false; autorun.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(200);
  check("and it can be turned back off", w.autorunCalls[w.autorunCalls.length - 1], false);

  tab("translate"); await sleep(30);
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

  tab("models"); await sleep(30);
  const del = d.querySelector('#proc-models button[data-a="ldel"]');
  check("active LLM model can be deleted", !!del, true);
  del.click(); await sleep(60);
  check("deleting a model asks first", !!d.querySelector(".modal-bg"), true);
  check("the question is asked in the app style, not by the browser", d.querySelectorAll(".modal .btn").length, 2);
  check("the way out comes first, the action second", [...d.querySelectorAll(".modal .btn")].map(b=>b.className), ["btn ghost", "btn yes"]);
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("the question closes with the answer", !!d.querySelector(".modal-bg"), false);
  check("model list empty after delete", d.querySelectorAll('#proc-models input[name="llmmdl"]').length, 0);


  tab("models"); await sleep(120);
  const pickAbsent = () => d.querySelector('#models input[name="mdl"][value="base"]');
  pickAbsent().click(); await sleep(80);
  check("picking a model that is not here asks first", !!d.querySelector(".modal-bg"), true);
  check("the question names the model and its size", d.querySelector(".modal p").textContent.includes("Base") && d.querySelector(".modal p").textContent.includes("142 MB"), true);
  d.querySelector(".modal .btn.ghost").click(); await sleep(250);
  check("saying no downloads nothing", w.dlCalls.length, 0);
  check("saying no puts the choice back", d.querySelector('#models input[name="mdl"][value="small"]').checked, true);

  pickAbsent().click(); await sleep(80);
  d.querySelector(".modal .btn.yes").click(); await sleep(250);
  check("saying yes starts the download", w.dlCalls, ["base"]);
  check("a download can be stopped", !!d.querySelector('#models button[data-a="cancel"][data-id="base"]'), true);
  d.querySelector('#models button[data-a="cancel"][data-id="base"]').click(); await sleep(250);
  check("stopping asks the program to stop", w.cancelCalls, ["base"]);
  check("a stopped download offers to start again", !!d.querySelector('#models button[data-a="dl"][data-id="base"]'), true);

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

  d.getElementById("wiz_next").click(); await sleep(300);
  check("the third step names the shortcut", d.getElementById("wiz_hot").textContent, "ctrl+win");
  check("the third step offers the microphones", d.getElementById("wiz_mic").options.length, 3);
  await sleep(300);
  check("the wizard meter follows the microphone", d.getElementById("wiz_micbar").style.width !== "" && d.getElementById("wiz_micbar").style.width !== "0%", true);

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
