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
    window.appState = async () =>
      JSON.stringify({ hotkey: "ctrl+win", mic: "Realtek", engine: "sherpa · gigaam-v3", llm: "model.gguf", ram: "8000 MB free", last: "hello", ready: true, status: "Ready" });
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
      "appPreviewSound", "appModelDl", "appMin", "appClose", "appDrag",
      "appDoUpdate", "appReady", "appJSError", "appCopyLast",
    ]) {
      window[name] = () => {};
    }
    window.appSave = async () => "";
    window.appModelDel = async () => "ok";
    window.appMics = async () =>
      JSON.stringify([
        { id: "dev1", name: "Headset (USB)", default: false },
        { id: "dev2", name: "Webcam microphone", default: false },
      ]);
    window.appMicLevel = async () => 0.42;
    window.appMicSelect = async () => "";
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
  check("status engine shown", d.getElementById("state_engine").textContent, "sherpa · gigaam-v3");
  check("status bar filled", d.getElementById("st_main").textContent, "Ready");
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
  mic.value = "dev1"; mic.dispatchEvent(new w.Event("change")); await sleep(30);
  check("microphone selection kept", mic.value, "dev1");
  await sleep(260);
  check("input level meter moves", d.getElementById("mic_bar").style.width !== "" && d.getElementById("mic_bar").style.width !== "0%", true);

  mic.value = ""; mic.dispatchEvent(new w.Event("change")); await sleep(30);
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
