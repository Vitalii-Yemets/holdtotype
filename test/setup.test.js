const fs = require("fs");
const path = require("path");
const { JSDOM } = require("jsdom");

const root = path.join(__dirname, "..");
const src = fs.readFileSync(path.join(root, "client", "setup", "page.go"), "utf8");

const marker = "const setupPage = `";
const start = src.indexOf(marker) + marker.length;
const end = src.indexOf("`", start);
if (start < marker.length || end < 0) {
  console.error("setupPage template not found in client/setup/page.go");
  process.exit(1);
}

const values = {
  TITLE: "HoldToType — Setup",
  APP: "HoldToType",
  TAGLINE: "Voice to text at the cursor.",
  PATH: "Install folder",
  MODEL: "Recognition model",
  MODELOPTS: '<option value="small">Small (466 MB)</option><option value="">Do not download</option>',
  SHORTCUT: "Start Menu shortcut",
  AUTORUN: "Start with Windows",
  LAUNCH: "Launch after install",
  INSTALL: "INSTALL",
  FRESHSTYLE: "",
  UPDNOTE: "",
  UPDATING: "false",
  PROG: "Installing…",
  DONE: "Installed",
  DONEAT: "The app is installed in:",
  WARNMODEL_JS: '"The model was not downloaded."',
  FINISH: "FINISH",
  RETRY: "RETRY",
  BACK: "BACK",
  NODIR_JS: '"Choose the installation folder"',
  VERSION: "0.0.0-test",
  DEFDIR: '"C:\\\\Programs\\\\HoldToType"',
};

let html = src.slice(start, end);
for (const [key, value] of Object.entries(values)) {
  html = html.split(`{{${key}}}`).join(value);
}
const leftovers = [...html.matchAll(/\{\{([A-Z_]+)\}\}/g)].map((m) => m[1]);

const failures = [];
function check(name, got, want) {
  const ok = JSON.stringify(got) === JSON.stringify(want);
  console.log(`${ok ? "PASS" : "FAIL"}  ${name}${ok ? "" : ` (got ${JSON.stringify(got)}, want ${JSON.stringify(want)})`}`);
  if (!ok) failures.push(name);
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

const dom = new JSDOM(html, {
  runScripts: "dangerously",
  beforeParse(w) {
    w.installCalls = [];
    w.appInstall = (...args) => w.installCalls.push(args);
    w.appBrowse = async () => "";
    w.appFinish = () => {};
    w.appDrag = () => {};
    w.appReady = () => {};
  },
});

const w = dom.window;
const d = w.document;
const shown = (id) => d.getElementById(id).classList.contains("on");

(async () => {
  const errors = [];
  w.addEventListener("error", (e) => errors.push(e.message));
  await sleep(200);

  check("every placeholder is filled", leftovers, []);
  check("opens on the options step", shown("st-opts"), true);
  check("the default folder is prefilled", d.getElementById("dir").value, "C:\\Programs\\HoldToType");

  d.getElementById("dir").value = "   ";
  d.querySelector("#st-opts .btn").click();
  await sleep(60);
  check("an empty folder does not start the install", w.installCalls.length, 0);
  check("an empty folder is explained", d.getElementById("operr").textContent, "Choose the installation folder");
  check("an empty folder keeps you on the options", shown("st-opts"), true);

  d.getElementById("dir").value = "C:\\Programs\\HoldToType";
  d.querySelector("#st-opts .btn").click();
  await sleep(60);
  check("a real folder starts the install", w.installCalls.length, 1);
  check("the install shows progress", shown("st-prog"), true);
  check("the folder error is cleared", d.getElementById("operr").textContent, "");

  w.setupProgress(40, "ggml-small.bin");
  await sleep(20);
  check("progress is reported", d.getElementById("plog").textContent, "ggml-small.bin");

  w.setupDone("mkdir Z:: the system cannot find the path specified.", false, "");
  await sleep(60);
  check("a failed install says why", d.getElementById("perr").textContent.length > 0, true);
  check("a failed install offers a way on", d.getElementById("pfoot").style.display, "");
  check("a failed install is not the end", shown("st-done"), false);

  const buttons = [...d.querySelectorAll("#pfoot .btn")].map((b) => b.textContent);
  check("retry and back are offered", buttons, ["RETRY", "BACK"]);

  d.querySelectorAll("#pfoot .btn")[1].click();
  await sleep(60);
  check("back returns to the options", shown("st-opts"), true);
  check("back keeps what was typed", d.getElementById("dir").value, "C:\\Programs\\HoldToType");
  check("back clears the failure", d.getElementById("perr").textContent, "");

  d.querySelector("#st-opts .btn").click();
  await sleep(60);
  check("retrying installs again", w.installCalls.length, 2);
  check("the failure buttons hide while it runs", d.getElementById("pfoot").style.display, "none");

  w.setupDone("", false, "C:\\Programs\\HoldToType");
  await sleep(60);
  check("a good install ends on the done step", shown("st-done"), true);
  check("the done step names the folder", d.getElementById("outdir").textContent, "C:\\Programs\\HoldToType");

  check("the switches are real checkboxes", d.querySelectorAll(".chk input[type=checkbox]").length, 3);
  check("every switch sits inside its label", [...d.querySelectorAll(".chk input")].every((i) => i.closest("label")), true);

  check("no page errors", errors, []);

  if (failures.length) {
    console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
    process.exit(1);
  }
  console.log("\nall installer checks passed");
  dom.window.close();
  process.exit(0);
})().catch((e) => {
  console.error("harness crashed:", e.message);
  process.exit(1);
});
