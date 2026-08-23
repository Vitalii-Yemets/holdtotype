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
    window.confirm = () => true;
    window.appLLM = async () => JSON.stringify(llmState);
    let micBadge = "Realtek";
    let restartFields = "";
    let micFails = false;
    let ruState = "ready";
    let otherState = "ready";
    let remote = false;
    window.micSelectCalls = 0;
    window.setMicFails = (v) => { micFails = v; };
    window.setModelStates = (ru, other) => { ruState = ru; otherState = other; };
    window.setRemote = (v) => { remote = v; };
    window.appState = async () =>
      JSON.stringify({ hotkey: "ctrl+win", mic: "Realtek", engine: "sherpa · gigaam-v3", llm: "model.gguf",
        ram: "8000 MB free", last: "hello", last_meta: "just now · 5 characters", ready: true, status: "Ready",
        status_line: "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free", ru_model: "gigaam-v3",
        other_model: "ggml-small.bin", llm_ok: true, mic_ok: true,
        restart_hint: restartFields, remote: remote,
        ru_state: ruState, other_state: otherState,
        badges: { mic: micBadge, models: "2", system: "" } });
    window.appRouting = async () =>
      JSON.stringify([
        { cond: "Speech in RU", engine: "gigaam-v3", why: "more accurate here" },
        { cond: "Other languages", engine: "ggml-small.bin", why: "99 languages" },
        { cond: "Translation", engine: "ggml-small.bin", why: "only Whisper translates" },
      ]);
    window.appModels = async () =>
      JSON.stringify([
        { id: "base", name: "Base", desc: "fast", size: 142, state: "absent", engine: "whisper", langs: "*" },
        { id: "small", name: "Small", desc: "balanced", size: 466, state: "active", engine: "whisper", langs: "*" },
        { id: "gigaam-v3", name: "GigaAM v3", desc: "russian", size: 232, state: "absent", engine: "sherpa", langs: "ru", punct: true },
      ]);
    window.appLLMSearch = async () =>
      JSON.stringify({ repos: [{ id: "org/Repo-GGUF", downloads: 1234, updated: "2026-01-01" }] });
    window.appLLMFiles = async () =>
      JSON.stringify({ files: [{ file: "q4.gguf", size: 4000, fit: "ok", need: 6166 }] });
    window.appLLMDel = async () => {
      llmState = { installed: [], downloads: [], ram: 16384, ram_free: 9000 };
      return "deleted";
    };
    window.appAdvise = async () => JSON.stringify({ primary: "gigaam-v3", companion: "small", text: "I recommend GigaAM v3.", ram: "8000 MB free" });
    for (const name of [
      "appLLMDlFile", "appLLMTest", "appHFPage", "appHFHome", "appRepoLink",
      "appAuthorLink", "appCapture", "appCaptureCombo", "appReload",
      "appPreviewSound", "appModelDl", "appMin", "appClose",
      "appDoUpdate", "appReady", "appJSError", "appCopyLast",
    ]) {
      window[name] = () => {};
    }
    window.dragCalls = 0;
    window.appDrag = () => { window.dragCalls++; };
    window.saveCalls = 0;
    window.appSave = async (json) => {
      window.saveCalls++;
      const f = JSON.parse(json);
      micBadge = f.mic_device_name ? f.mic_device_name.split(" ")[0] : "Realtek";
      if(Number(f.server_port) !== 8910) restartFields = "Restart HoldToType to apply: Port";
      return JSON.stringify({ ok: true, severity: "ok", message: "Saved" });
    };
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
  check("restart hint follows the app state", d.getElementById("st_pend").textContent, "");
  check("status bar names both models", d.getElementById("st_main").textContent, "Ready · gigaam-v3 + ggml-small.bin · 7.8 GB free");
  check("sidebar badges filled", [d.getElementById("badge_mic").textContent, d.getElementById("badge_models").textContent], ["Realtek", "2"]);
  check("status bar led lit", d.getElementById("st_led").classList.contains("on"), true);

  check("eight sections in the sidebar", d.querySelectorAll(".nav").length, 8);
  const everySetting = [
    "hotkey", "min_record_ms", "max_record_seconds", "auto_enter", "restore_clipboard",
    "type_mode", "overlay", "animation", "mic_device", "mic_bar", "beep", "sound_theme",
    "routing", "models", "language", "threads", "punctuation", "whisper_prompt", "profbody",
    "tr_default", "translate_target", "translate_ask", "translate_ask_seconds", "tr_hotkey",
    "tl_en", "ui_language", "upd_check", "check_updates", "server_autostart", "server_port",
    "server_exe", "server_url", "proc-models", "proc-search", "ver2",
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

  tab("system"); await sleep(30);
  check("service settings shown", [shown("system"), !!d.getElementById("server_url")], [true, true]);

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
  del.click(); await sleep(200);
  check("model list empty after delete", d.querySelectorAll('#proc-models input[name="llmmdl"]').length, 0);

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
  check("no unsaved-changes dialog left", !!d.getElementById("modalbg"), false);

  tab("system"); await sleep(80);
  const port = d.getElementById("server_port");
  port.value = "8999"; port.dispatchEvent(new w.Event("change", { bubbles: true }));
  await sleep(400);
  const pendText = d.getElementById("st_pend").textContent;
  check("the port says it needs a restart", pendText, "Restart HoldToType to apply: Port");
  check("the restart hint never mentions Save", /save/i.test(pendText), false);

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
