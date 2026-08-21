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
    window.appModels = async () =>
      JSON.stringify([
        { id: "base", name: "Base", desc: "fast", size: 142, state: "absent" },
        { id: "small", name: "Small", desc: "balanced", size: 466, state: "active" },
      ]);
    window.appLLMSearch = async () =>
      JSON.stringify({ repos: [{ id: "org/Repo-GGUF", downloads: 1234, updated: "2026-01-01" }] });
    window.appLLMFiles = async () =>
      JSON.stringify({ files: [{ file: "q4.gguf", size: 4000, fit: "ok", need: 6166 }] });
    window.appLLMDel = async () => {
      llmState = { installed: [], downloads: [], ram: 16384, ram_free: 9000 };
      return "deleted";
    };
    for (const name of [
      "appLLMDlFile", "appLLMTest", "appHFPage", "appHFHome", "appRepoLink",
      "appAuthorLink", "appCapture", "appCaptureCombo", "appReload",
      "appPreviewSound", "appModelDl", "appMin", "appClose", "appDrag",
      "appDoUpdate", "appReady", "appJSError",
    ]) {
      window[name] = () => {};
    }
    window.appSave = async () => "";
    window.appModelDel = async () => "ok";
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

  const tab = (p) => d.querySelector(`.tab[data-p="${p}"]`).click();
  const sub = (page, s) => d.querySelector(`#p-${page} .stab[data-s="${s}"]`).click();
  const visible = (id) => d.getElementById(id).classList.contains("on");

  tab("rec"); await sleep(80);
  check("recognition opens on Models", visible("rec-models"), true);
  check("whisper models listed", d.querySelectorAll('#rec-models input[name="mdl"]').length, 2);
  sub("rec", "dict"); await sleep(30);
  check("dictionary textarea is large", d.getElementById("whisper_prompt").getAttribute("rows"), "14");
  sub("rec", "server"); await sleep(30);
  check("server subtab shown", [visible("rec-server"), visible("rec-params")], [true, false]);

  sub("rec", "translate"); await sleep(30);
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

  tab("proc"); await sleep(80);
  check("post-processing opens on Models", visible("proc-models"), true);
  sub("proc", "prompts"); await sleep(30);
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

  sub("proc", "models"); await sleep(30);
  const del = d.querySelector('#proc-models button[data-a="ldel"]');
  check("active LLM model can be deleted", !!del, true);
  del.click(); await sleep(200);
  check("model list empty after delete", d.querySelectorAll('#proc-models input[name="llmmdl"]').length, 0);

  tab("about"); await sleep(60);
  check("about opens on Info", visible("about-info"), true);
  sub("about", "help"); await sleep(30);
  check("guide subtab shown", visible("about-help"), true);
  sub("about", "author"); await sleep(30);
  check("author subtab shown", visible("about-author"), true);

  check("no page errors", errors, []);

  if (failures.length) {
    console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
    process.exit(1);
  }
  console.log("\nall UI checks passed");
})().catch((e) => {
  console.error("harness crashed:", e.message);
  process.exit(1);
});
