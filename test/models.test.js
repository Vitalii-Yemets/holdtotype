const fs = require("fs");
const path = require("path");
const { JSDOM } = require("jsdom");

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const failures = [];
function check(name, actual, expected) {
  const ok = JSON.stringify(actual) === JSON.stringify(expected);
  console.log(`${ok ? "PASS" : "FAIL"}  ${name}${ok ? "" : ` (got ${JSON.stringify(actual)}, want ${JSON.stringify(expected)})`}`);
  if (!ok) failures.push(name);
}

const NEMO_LANGS = "ar,bg,cs,da,de,el,en,es,et,fi,fr,he,hi,hr,hu,it,ja,ko,lt,lv,mt,nb,nl,nn,pl,pt,ro,ru,sk,sl,sv,th,tr,uk,vi,zh";
const CATALOG = [
  { id: "base", name: "Base", desc: "fast", size: 142, ram: 273, engine: "whisper", langs: "*", auto: true, translate: true, speed: 5, accuracy: 2 },
  { id: "small", name: "Small", desc: "balanced", size: 466, ram: 921, engine: "whisper", langs: "*", auto: true, translate: true, speed: 3, accuracy: 3 },
  { id: "medium-q5_0", name: "Medium (q5)", desc: "recommended", size: 539, ram: 868, engine: "whisper", langs: "*", auto: true, translate: true, speed: 2, accuracy: 4 },
  { id: "gigaam-v3", name: "GigaAM v3", desc: "russian", size: 232, ram: 278, engine: "sherpa", langs: "ru", punct: true, speed: 5, accuracy: 5 },
  { id: "nemotron-3.5", name: "Nemotron 3.5", desc: "streaming", size: 651, ram: 781, engine: "stream", langs: NEMO_LANGS, auto: true, punct: true, speed: 3, accuracy: 5 },
  { id: "parakeet-v3", name: "Parakeet TDT 0.6B v3", desc: "european", size: 670, ram: 804, engine: "sherpa", langs: "en,es,fr,de,pl,ru,uk", auto: true, punct: true, speed: 5, accuracy: 4 },
  { id: "moonshine-uk", name: "Moonshine Base uk", desc: "ukrainian", size: 135, engine: "sherpa", langs: "uk", manual: true, link: "https://example.com/moonshine", speed: 5, accuracy: 3 },
  { id: "local:my-model", name: "my-model", desc: "found in the models folder", size: 900, engine: "whisper", langs: "*", custom: true },
];

function pageFor(lang, cfgPatch) {
  let html = fs.readFileSync(path.join(__dirname, lang === "en" ? "page.html" : "page-" + lang + ".html"), "utf8");
  const at = html.indexOf("const CFG = ") + "const CFG = ".length;
  const lineEnd = html.indexOf("\n", at);
  const raw = html.slice(at, lineEnd).replace(/;\s*$/, "");
  const cfg = Object.assign(JSON.parse(raw), cfgPatch);
  html = html.slice(0, at) + JSON.stringify(cfg) + ";" + html.slice(lineEnd);
  return html;
}

function open(lang, cfgPatch, setup) {
  const html = pageFor(lang, cfgPatch);
  const names = [...new Set(html.match(/\bapp[A-Z][A-Za-z]*/g) || [])];
  const S = {
    states: { base: "absent", small: "installed", "medium-q5_0": "installed", "gigaam-v3": "installed", "nemotron-3.5": "active", "parakeet-v3": "absent", "moonshine-uk": "absent", "local:my-model": "installed" },
    loaded: { "nemotron-3.5": true }, pct: {}, hidden: {},
    llm: [{ file: "a.gguf", size: 1066, need: 2743, active: true, loaded: true }, { file: "b.gguf", size: 1950, need: 3775, active: false, loaded: false }],
    llmDl: [],
    trm: { name: "Gemma 3 4B", state: "ready", size: 2374, fits: true, running: false, pct: 0 },
    calls: [], saves: [], searches: [],
  };
  if (setup) setup(S);
  const dom = new JSDOM(html, {
    runScripts: "dangerously",
    beforeParse(window) {
      window.S = S;
      const json = (v) => JSON.stringify(v);
      const stubs = {
        appModels: async () => json(CATALOG.filter((m) => !S.hidden[m.id]).map((m) => Object.assign({}, m, { state: S.states[m.id], pct: S.pct[m.id] || 0, loaded: S.states[m.id] !== "absent" && !!S.loaded[m.id] }))),
        appModelDl: async (id) => { S.calls.push("dl:" + id); S.states[id] = "downloading"; },
        appModelCancel: async (id) => { S.calls.push("cancel:" + id); if (S.states[id] === "downloading") S.states[id] = "absent"; return true; },
        appModelDel: async (id, force) => { S.calls.push("del:" + id + ":" + force); S.states[id] = "absent"; delete S.loaded[id]; return "Deleted"; },
        appUnloadEngines: async () => { S.calls.push("unload"); S.loaded = {}; S.llm.forEach((m) => { m.loaded = false; }); },
        appUnloadRecognition: async () => { S.calls.push("unloadrec"); S.loaded = {}; },
        appModelLink: (id) => { S.calls.push("link:" + id); },
        appLLM: async () => json({ installed: S.llm, downloads: S.llmDl, ram: 16384, ram_free: 9000 }),
        appLLMUnload: async () => { S.calls.push("llmunload"); S.llm.forEach((m) => { m.loaded = false; }); },
        appLLMDel: async (f) => { S.calls.push("llmdel:" + f); S.llm = S.llm.filter((m) => m.file !== f); return "Deleted"; },
        appLLMSearch: async (q) => { S.searches.push(q); return json({ repos: [{ id: "org/Repo-GGUF", downloads: 1234, updated: "2026-01-01" }] }); },
        appLLMFiles: async () => json({ files: [{ file: "c.gguf", size: 900, fit: "ok", need: 2550 }, { file: "a.gguf", size: 1066, fit: "ok", need: 2743 }] }),
        appLLMDlFile: async (repo, file, size) => { S.calls.push("llmdl:" + file); S.llmDl.push({ file, pct: 10, err: "" }); },
        appTrModel: async (action) => { S.calls.push("trm:" + action); if (action === "delete") S.trm.state = "absent"; if (action === "download") S.trm.state = "downloading"; return ""; },
        appSave: async (f) => { S.saves.push(JSON.parse(f)); return json({ ok: true, severity: "ok", message: "Saved" }); },
        appState: async () => json({ subtitles: { phase: "off" }, subs_tr: { choice: "quality", missing: "", pair: "", model: S.trm, fast: { exe: true, pairs: [], offer: [], used: 0 } },
          subs_place: { position: "bottom", custom: {} }, hotkey: "ctrl+win", mic: "Mic", engine: "", llm: "a.gguf", ram: "", last: "", last_meta: "", ready: S.ready !== false, status: S.status || "Ready",
          status_line: S.statusLine || "Ready", active_model: "Nemotron 3.5", active_state: "ready", active_lang: "RU", assigned: [], installed_models: S.installedModels || [], loaded_now: Object.keys(S.loaded).length ? "Nemotron 3.5" : "", week_line: "",
          llm_ok: true, mic_ok: true, last_at: 0, last_app: "", remote: false, backend_err: "", post_err: "", stt_broken: false, whisper_now: false, enabled: true, autostart: false,
          disk_mb: 0, ram_mb: 16384, ram_free_mb: 9000, week_count: 0, week_chars: 0, today_count: 0, week_apps: [], badges: { mic: "", models: "", system: "" },
          post_on: true, post_model: "a", post_size: 1066, post_installed: S.llm.length, post_remote: false, post_prompts: 1 }),
        appSubtitlesModel: async (l) => json({ lang: l, id: "nemotron-3.5", name: "Nemotron 3.5", engine: "stream", mode: "live", installed: true, whisper: false, remote: false }),
        appMics: async () => "[]",
        appMicLevel: async () => 0,
        appHistory: async () => "[]",
        appUpdateStatus: async () => json({ current: "0.0.0", latest: "", url: "" }),
        appAutorun: async () => false,
        appPostKeySet: async () => false,
        appCheckModels: async () => json({ ok: true, text: "All files are intact", rows: [] }),
        appMaximized: async () => false,
        appWindowWidth: async (w) => json({ prev: 700, width: w }),
        appHFHome: () => { S.calls.push("hfhome"); },
      };
      for (const n of names) window[n] = stubs[n] || (async () => "{}");
      for (const n of Object.keys(stubs)) window[n] = stubs[n];
    },
  });
  return dom;
}

(async () => {
  const errors = [];
  const dom = open("en", {
    lang_models: { auto: "nemotron-3.5", ru: "gigaam-v3", de: "small" }, language: "ru", subtitles_language: "en",
    subtitles_translate_on: true, subtitles_translator: "quality", llm_model: "a.gguf", post_enabled: true, post_source: "local",
  });
  const w = dom.window, d = w.document, S = w.S;
  w.addEventListener("error", (e) => errors.push(e.message));
  await sleep(700);
  const tab = (p) => d.querySelector(`.nav[data-p="${p}"]`).click();
  const shown = (p) => d.getElementById(`p-${p}`).classList.contains("active");
  const rows = (id) => [...d.querySelectorAll("#" + id + " .librow")];
  const row = (key) => d.querySelector('#p-library .librow[data-key="' + key + '"]');
  const use = (key) => row(key).querySelector(".lmodel").textContent;
  const refresh = async () => { await w.refreshModels(); await w.refreshLLM(); await w.refreshState(); await sleep(120); };
  const lastSave = () => S.saves[S.saves.length - 1];
  const modalText = () => { const p = d.querySelector(".modal p"); return p ? p.textContent : null; };
  const answer = async (yes) => { d.querySelector(".modal .btn." + (yes ? "yes" : "ghost")).click(); await sleep(400); };

  tab("library"); await refresh();
  check("the models tab opens", shown("library"), true);
  check("four blocks in their order", [...d.querySelectorAll("#p-library > .card")].map((c) => c.id), ["lib_have", "lib_get", "lib_hf", "lib_own"]);
  check("no lead line above the blocks", d.querySelector("#p-library > .hint"), null);
  check("the installed models come in their order: recognition by name, editors, the translation model", rows("lib_have_list").map((r) => r.dataset.key),
    ["asr:gigaam-v3", "asr:medium-q5_0", "asr:local:my-model", "asr:nemotron-3.5", "asr:small", "llm:a.gguf", "llm:b.gguf", "trm"]);
  check("the total counts the models and their size", d.getElementById("lib_have_total").textContent, "8 models · 8.0 GB");
  check("a recognition model names the languages it dictates", use("asr:gigaam-v3"), "used for dictation: русский");
  check("auto-detect, the languages that follow it and the subtitles are named", use("asr:nemotron-3.5"), "used for dictation: auto-detect, en, es, fr, it, pl, uk, and for subtitles");
  check("a language of its own is named for its model", use("asr:small"), "used for dictation: deutsch");
  check("an unused model says so and where it is turned on", [row("asr:medium-q5_0").querySelector(".libidle").textContent, row("asr:medium-q5_0").querySelector(".libturn").textContent],
    ["not in use", "turn on: Languages →"]);
  check("the plate names the place in its tip", (row("asr:medium-q5_0").querySelector(".jump").dataset.tip || row("asr:medium-q5_0").querySelector(".jump").title), "The “Languages” tab, the “Which model serves which language” block");
  check("the editor in use is named for post-processing", use("llm:a.gguf"), "used for post-processing");
  check("an unused editor leads to post-processing", row("llm:b.gguf").querySelector(".libturn").textContent, "turn on: Post-processing →");
  check("the translation model in use is named for translation", use("trm"), "used for translating subtitles");
  check("recognition models carry the info badge, accuracy and speed", [!!row("asr:gigaam-v3").querySelector(".pinfo"), row("asr:gigaam-v3").querySelectorAll(".mbar").length], [true, 2]);
  check("a hand-copied model has no bars", row("asr:local:my-model").querySelectorAll(".mbar").length, 0);
  check("editors and the translation model have no bars either", [row("llm:a.gguf").querySelectorAll(".mbar").length, row("trm").querySelectorAll(".mbar").length], [0, 0]);
  check("the model in memory wears the chip and the unload icon", [row("asr:nemotron-3.5").querySelector(".pchip").textContent, !!row("asr:nemotron-3.5").querySelector('button[data-lib="unload"]')], ["in memory", true]);
  check("a model not in memory has neither", [!!row("asr:small").querySelector(".pchip"), !!row("asr:small").querySelector('button[data-lib="unload"]')], [false, false]);
  check("every installed model has its size and a cross", rows("lib_have_list").every((r) => r.querySelector(".psize").textContent && r.querySelector('button[data-lib="del"]')), true);
  check("the editor in memory says how much memory it takes", row("llm:a.gguf").querySelector(".pram").textContent, "≈2.7 GB RAM");
  check("the translation model has its plates and the terms plate", [...row("trm").querySelectorAll(".libl2 > *")].map((e) => e.textContent), ["graphics card with 4 GB of its own memory or more", "Terms of use ↗"]);
  const checkRow = d.getElementById("mcheck").closest(".row");
  check("the check row closes the block", [checkRow.querySelector("label").firstChild.textContent, [...d.querySelectorAll("#lib_have > .row")].pop() === checkRow], ["S_MCHECK", true]);
  const pageCss = [...d.querySelectorAll("style")].map((s) => s.textContent).join(" ");
  check("in a narrow window a model row goes in lines: name and size, then the text, then the plates and the bars", ["#lib_have_list,#lib_get_list{container:libs/inline-size}", "@container libs (max-width:819px){", ".librow .libmain,.librow .libl1{display:contents}", ".librow .pname{grid-column:1;grid-row:1;flex:none;width:auto}", ".librow .pact{grid-column:2;grid-row:1}", ".librow .lmodel{grid-column:1/-1;grid-row:2}", ".librow .libl2{grid-column:1;grid-row:3;min-width:0}", ".librow .libbars{grid-column:2;grid-row:3}", ".librow .libmain>:not(.libl1):not(.libl2){grid-column:1/-1}"].filter((r) => !pageCss.includes(r)), []);
  check("the wide row stays as it was", [".libl1{display:flex;align-items:baseline;gap:12px}", ".librow .pname{flex:0 0 238px;width:238px;", ".libbars{flex:none;width:160px}"].filter((r) => !pageCss.includes(r)), []);
  const idleSel = d.getElementById("models_idle_minutes"), idleRow = d.getElementById("idle_row");
  check("unloading when idle is a row with a list right before the check row", [idleRow.nextElementSibling === checkRow, idleRow.classList.contains("noline"), idleRow.querySelector("label").firstChild.textContent, idleRow.querySelector(".sub").textContent, idleSel.tagName],
    [true, false, "S_IDLE_UNLOAD", "S_IDLE_UNLOAD_SUB", "SELECT"]);
  check("the list offers never, 5, 10 and 30 minutes and an hour", [...idleSel.options].map((o) => [o.value, o.textContent]), [["0", "S_IDLE_NEVER"], ["5", "S_IDLE_5"], ["10", "S_IDLE_10"], ["30", "S_IDLE_30"], ["60", "S_IDLE_60"]]);
  check("without the key in the config nothing is unloaded", idleSel.value, "0");
  idleSel.value = "10"; idleSel.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(450);
  check("a choice is saved as minutes", lastSave().models_idle_minutes, 10);
  idleSel.value = "60"; idleSel.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(450);
  check("an hour is saved as 60", lastSave().models_idle_minutes, 60);
  idleSel.value = "0"; idleSel.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(450);
  check("and never as 0", lastSave().models_idle_minutes, 0);
  check("the search finds the row", w.searchMatches("s_idle_unload").some((el) => el === idleRow), true);

  check("the models to download, by name", rows("lib_get_list").map((r) => r.dataset.key), ["asr:base", "asr:moonshine-uk", "asr:parakeet-v3"]);
  check("the hint says what they are for", d.getElementById("lib_get_hint").textContent, "For dictation and subtitles — models that turn speech into text. A downloaded model moves to “Installed models”.");
  check("a model to download shows its description, size and the arrow", [use("asr:base"), row("asr:base").querySelector(".psize").textContent, !!row("asr:base").querySelector('button[data-lib="dl"]')], ["fast", "142 MB", true]);
  check("the licence-bound one offers its own link and explains why", [!!row("asr:moonshine-uk").querySelector('button[data-lib="dl"]'), row("asr:moonshine-uk").querySelector('button[data-lib="link"]').textContent.trim().startsWith("Download yourself"), !!row("asr:moonshine-uk").querySelector(".pdesc")], [false, true, true]);

  check("the licence-bound button stands with its note, and the right column keeps the size with empty icon places", [!!row("asr:moonshine-uk").querySelector('.libnote button[data-lib="link"]'), row("asr:moonshine-uk").querySelector(".pact .psize").textContent, row("asr:moonshine-uk").querySelectorAll(".pact button").length, row("asr:moonshine-uk").querySelectorAll(".pact .libslot").length], [true, "135 MB", 0, 2]);
  row("asr:base").querySelector('button[data-lib="dl"]').click(); await sleep(150); S.pct.base = 37; await refresh();
  check("downloading shows the percent and a cancel cross", [row("asr:base").querySelector(".mpct").textContent, !!row("asr:base").querySelector('button[data-lib="cancel"]')], ["37%", true]);
  S.states.base = "installed"; await refresh();
  check("a downloaded model moves to the installed ones", [!!d.querySelector('#lib_get_list .librow[data-key="asr:base"]'), !!d.querySelector('#lib_have_list .librow[data-key="asr:base"]')], [false, true]);

  S.calls.length = 0;
  row("asr:nemotron-3.5").querySelector('button[data-lib="unload"]').click(); await sleep(300); await refresh();
  check("the unload icon of a recognition model unloads recognition only", S.calls.slice(), ["unloadrec"]);
  check("and the chip and the icon go", [!!row("asr:nemotron-3.5").querySelector(".pchip"), !!row("asr:nemotron-3.5").querySelector('button[data-lib="unload"]')], [false, false]);
  check("while the editor stays in memory", [!!row("llm:a.gguf").querySelector(".pchip"), !!row("llm:a.gguf").querySelector('button[data-lib="unload"]')], [true, true]);
  S.calls.length = 0;
  row("llm:a.gguf").querySelector('button[data-lib="unload"]').click(); await sleep(300); await refresh();
  check("the unload icon of the editor unloads the editor only", S.calls.slice(), ["llmunload"]);
  S.loaded = { "nemotron-3.5": true }; S.llm[0].loaded = true; await refresh();
  S.calls.length = 0;
  tab("state"); await w.refreshState(); await sleep(150);
  d.getElementById("state_unload").click(); await sleep(300);
  check("the Unload button on Status unloads everything", S.calls.slice(), ["unload"]);
  tab("library"); await refresh();
  check("and every chip goes", [!!row("asr:nemotron-3.5").querySelector(".pchip"), !!row("llm:a.gguf").querySelector(".pchip")], [false, false]);

  tab("models"); await sleep(100);
  const lrow = (l) => d.querySelector('#langlist .lrow[data-lang="' + l + '"]');
  lrow("ru").click(); await sleep(150);
  const cards = () => [...d.querySelectorAll("#langlist .lpick .pcard")].map((c) => c.dataset.id);
  check("an unfolded language offers its installed models only", cards(), ["base", "gigaam-v3", "medium-q5_0", "local:my-model", "nemotron-3.5", "small"]);
  check("without download, unload or delete icons", d.querySelectorAll('#langlist .lpick button[data-a="dl"], #langlist .lpick button[data-a="unload"], #langlist .lpick button[data-a="del"]').length, 0);
  check("the chosen one wears the tick", d.querySelector('#langlist .pcard[data-id="gigaam-v3"] .pcur').textContent, "✓");
  check("and a plate leads to the models to download", [d.querySelector("#langlist .jumpline > span").textContent, d.querySelector("#langlist .jumpline .jump").textContent], ["Download another model:", "Models: Models to download →"]);
  lrow("pl").click(); await sleep(150);
  check("a missing recommended model stands last and says it is not installed", [cards().slice(-1)[0], d.querySelector('#langlist .pcard[data-id="parakeet-v3"] .psize').textContent], ["parakeet-v3", "not installed · 670 MB"]);
  const pk = d.querySelector('#langlist .pcard[data-id="parakeet-v3"] input.psw');
  pk.checked = true; pk.dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(150);
  check("turning it on asks the existing download question", d.querySelector(".modal .mtitle").textContent, "Downloading a model");
  check("the question names the model and its size", modalText().includes("Parakeet TDT 0.6B v3") && modalText().includes("670 MB"), true);
  await answer(true);
  check("and starts the download", S.calls.includes("dl:parakeet-v3"), true);
  S.pct["parakeet-v3"] = 12; await refresh();
  check("the card shows the percent with a cancel cross", [d.querySelector('#langlist .pcard[data-id="parakeet-v3"] .mpct').textContent, !!d.querySelector('#langlist .pcard[data-id="parakeet-v3"] button[data-a="cancel"]')], ["12%", true]);
  d.querySelector('#langlist .pcard[data-id="parakeet-v3"] button[data-a="cancel"]').click(); await sleep(300);
  check("cancelling takes the language back", "pl" in lastSave().lang_models, false);
  S.hidden["nemotron-3.5"] = true; S.hidden["parakeet-v3"] = true; S.hidden.base = true; S.hidden.small = true; S.hidden["medium-q5_0"] = true; S.hidden["local:my-model"] = true;
  await refresh(); lrow("uk").click(); await sleep(150);
  check("with nothing that fits, the language says so and still leads to the downloads", [d.querySelector("#langlist .lpick .nonote").textContent, !!d.querySelector("#langlist .lpick .jumpline .jump")], ["There are no suitable models on this computer.", true]);
  S.hidden = {}; await refresh();

  tab("library"); await refresh();
  const saves0 = S.saves.length;
  row("asr:local:my-model").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting a free model asks the existing question", modalText(), 'Delete the "my-model" model? It can be downloaded again.');
  await answer(true); await refresh();
  check("and deletes it without touching the languages", [S.calls.includes("del:local:my-model:true"), S.saves.length], [true, saves0]);

  row("asr:small").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting a model in use names its use", modalText(), "The “Small” model is used for dictation: deutsch. Delete it? The Auto-detect model will work in its place.");
  await answer(false);
  check("saying no keeps everything", [S.calls.includes("del:small:true"), !!row("asr:small")], [false, true]);
  row("asr:small").querySelector('button[data-lib="del"]').click(); await sleep(150);
  await answer(true); await refresh();
  check("its language goes back to auto-detect", ["de" in lastSave().lang_models, S.calls.includes("del:small:true")], [false, true]);
  tab("models"); await sleep(120);
  check("and the language row follows auto-detect", lrow("de").querySelector(".lmodel").textContent, "as Auto-detect · Nemotron 3.5");

  tab("library"); await refresh();
  row("asr:nemotron-3.5").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting the auto-detect model names the one auto-detect moves to", modalText(), "The “Nemotron 3.5” model is used for dictation: auto-detect, de, en, es, fr, it, pl, uk, and for subtitles. Delete it? Auto-detect will move to Medium (q5).");
  await answer(true); await refresh();
  check("auto-detect moves to that model", lastSave().lang_models.auto, "medium-q5_0");
  check("the model is deleted after the move", S.calls.indexOf("del:nemotron-3.5:true") > -1, true);
  check("the new auto-detect model now names its use", use("asr:medium-q5_0").startsWith("used for dictation: auto-detect"), true);

  row("asr:base").querySelector('button[data-lib="del"]').click(); await sleep(150);
  await answer(true); await refresh();
  row("asr:medium-q5_0").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting the last model auto-detect could use says dictation will stop", modalText(), "The “Medium (q5)” model is used for dictation: auto-detect, de, en, es, fr, it, pl, uk, and for subtitles. Delete it? There is no other recognition model on this computer — dictation will not work until you download one.");
  await answer(true); await refresh();
  check("auto-detect is left without a model of its own", "auto" in lastSave().lang_models, false);
  check("the model it had is among the ones to download again", !!d.querySelector('#lib_get_list .librow[data-key="asr:medium-q5_0"]'), true);

  row("llm:a.gguf").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting the editor in use names its use", modalText(), "The “a.gguf” model is used for post-processing. Delete it? Another model will have to be chosen there.");
  await answer(true); await refresh();
  tab("post"); await sleep(250);
  check("post-processing loses its choice and warns", [d.getElementById("llm_pick").value, d.getElementById("post_warn").style.display, d.getElementById("post_warn").textContent.length > 0], ["", "", true]);
  check("the list offers what is left", [...d.getElementById("llm_pick").options].map((o) => o.value), ["", "b.gguf"]);
  d.getElementById("llm_pick").value = "b.gguf"; d.getElementById("llm_pick").dispatchEvent(new w.Event("change", { bubbles: true })); await sleep(300);
  check("choosing from the list saves the editor", lastSave().llm_model_file, "b.gguf");
  check("the card names the size", [...d.querySelectorAll("#llm_sum .sumrow")].slice(0, 2).map((r) => r.querySelector(".sumk").textContent + "=" + (r.querySelector(".sumv") ? r.querySelector(".sumv").textContent : r.querySelector("select").value)), ["model=b.gguf", "size=1.9 GB"]);
  tab("library"); await refresh();
  row("llm:b.gguf").querySelector('button[data-lib="del"]').click(); await sleep(150);
  await answer(true); await refresh();
  tab("post"); await sleep(250);
  check("with no editor left the card says so and leads to the search", [d.querySelector("#llm_sum .postnone .nonote").textContent, d.querySelector("#llm_sum .postnone .jump").textContent], ["There is no editor model on this computer.", "Models: Search on Hugging Face →"]);

  tab("library"); await refresh();
  row("trm").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting the translation model in use names its use", modalText(), "The “Gemma 3 4B” model is used for translating subtitles. Delete it? Another model will have to be chosen there.");
  await answer(true); await refresh(); await sleep(200);
  check("the translation model goes to the ones to download", [!!d.querySelector('#lib_get_list .librow[data-key="trm"]'), use("trm")], [true, "for translating subtitles"]);
  check("and the hint names it", d.getElementById("lib_get_hint").textContent.includes("For translating subtitles — Gemma 3 4B."), true);
  tab("subs"); await sleep(200);
  check("the quality card says it is not installed", d.getElementById("subs_q_size").textContent, "not installed · 2.3 GB");
  tab("library"); await refresh();
  row("trm").querySelector('button[data-lib="dl"]').click(); await sleep(200);
  check("the arrow downloads the translation model", S.calls.includes("trm:download"), true);

  const q = d.getElementById("hf_q");
  check("before a search only the field is shown", ["hf_ramline", "hf_fitline", "hf_results"].map((id) => d.getElementById(id).style.display), ["none", "none", "none"]);
  check("the magnifier is a button inside the field, and there is no other search button", [d.getElementById("hf_find").closest(".srchbox") !== null, !!d.getElementById("hf_go")], [true, false]);
  q.value = "qwen"; q.dispatchEvent(new w.Event("input", { bubbles: true }));
  check("the cross appears with text", d.getElementById("hf_clr").style.display, "block");
  q.dispatchEvent(new w.KeyboardEvent("keydown", { key: "Enter", bubbles: true })); await sleep(250);
  check("Enter searches", [S.searches, d.getElementById("hf_results").style.display], [["qwen"], ""]);
  d.querySelector(".hfrepo").click(); await sleep(250);
  const acts = [...d.querySelectorAll("#hf_results .hfact")].map((e) => e.querySelector("button") ? e.querySelector("button").getAttribute("aria-label") : e.textContent);
  check("found files show the download arrow", acts, ["Download", "Download"]);
  d.querySelector('#hf_results button[data-file="c.gguf"]').click(); await sleep(300);
  check("downloading a file shows the percent and a cancel cross", [d.querySelector("#hf_results .hfact .mpct").textContent, !!d.querySelector("#hf_results button[data-hfcancel]")], ["10%", true]);
  d.getElementById("hf_clr").click(); await sleep(150);
  check("the cross clears and hides the results", [q.value, d.getElementById("hf_results").style.display], ["", "none"]);
  d.getElementById("hf_find").click(); await sleep(250);
  check("the magnifier searches too", S.searches.length, 2);

  tab("models"); await sleep(100);
  lrow("ru").click(); await sleep(150);
  d.querySelector("#langlist .jumpline .jump").click(); await sleep(150);
  check("a plate opens the models tab", shown("library"), true);
  tab("post"); await sleep(250);
  d.querySelector("#llm_sum .postnone .jump").dispatchEvent(new w.KeyboardEvent("keydown", { key: "Enter", bubbles: true })); await sleep(150);
  check("Enter on a plate works too", shown("library"), true);
  check("plates are focusable links", [d.querySelector("#langlist .jump").getAttribute("role"), d.querySelector("#langlist .jump").tabIndex], ["link", 0]);
  check("no underlined links in these places", d.querySelectorAll("#p-library .lnk, #src_local .lnk, #subs_src_quality .lnk, #langlist .lnk").length, 0);
  const hits = w.searchMatches("s_lib_hf");
  check("settings search finds the blocks of the models tab", hits.some((e) => e.closest("#p-library")), true);
  check("and the models by name", w.searchMatches("gigaam").some((e) => e.closest("#p-library")), true);
  w.searchSettings("gigaam"); w.showHit(0); await sleep(100);
  check("a found model row is marked like a found setting", [shown("library"), d.querySelector("#p-library .librow.hit") && d.querySelector("#p-library .librow.hit").dataset.key], [true, "asr:gigaam-v3"]);
  check("and the mark has a rule of its own", [...d.styleSheets].some((s) => [...s.cssRules].some((r) => r.selectorText === ".librow.hit" && /inset/.test(r.style.boxShadow))), true);
  w.searchSettings("");

  for (const [v, want] of [[30, "30"], [5, "5"], [7, "0"]]) {
    const one = open("en", { models_idle_minutes: v });
    await sleep(500);
    check("a saved " + v + " shows as " + want, one.window.document.getElementById("models_idle_minutes").value, want);
    one.window.close();
  }
  {
    const src = fs.readFileSync(path.join(__dirname, "..", "client", "i18n.go"), "utf8").replace(/\r\n/g, "\n");
    const block = (lang) => {
      const at = src.indexOf('\n\t"' + lang + '": {', src.indexOf("var settingsStrings"));
      const out = {};
      for (const m of src.slice(at, src.indexOf("\n\t},", at)).matchAll(/"(S_[A-Z0-9_]+)":\s*"((?:[^"\\]|\\.)*)"/g)) out[m[1]] = JSON.parse('"' + m[2] + '"');
      return out;
    };
    const tray = (lang) => {
      const m = src.match(new RegExp('"' + lang + '": \\{[^]*?"status\\.parked":\\s*"([^"]*)"'));
      return m && m[1];
    };
    const ruS = block("ru"), enS = block("en");
    check("Russian: the row texts", [ruS.S_IDLE_UNLOAD, ruS.S_IDLE_UNLOAD_SUB, ruS.S_IDLE_NEVER, ruS.S_IDLE_5, ruS.S_IDLE_10, ruS.S_IDLE_30, ruS.S_IDLE_60],
      ["Выгружать модели при простое", "если выбранное время не было ни диктовки, ни субтитров, модели освобождают память; следующая диктовка загрузит модель заново, и текст появится на несколько секунд позже", "не выгружать", "через 5 минут", "через 10 минут", "через 30 минут", "через час"]);
    check("English: the row texts", [enS.S_IDLE_UNLOAD, enS.S_IDLE_UNLOAD_SUB, enS.S_IDLE_NEVER, enS.S_IDLE_5, enS.S_IDLE_10, enS.S_IDLE_30, enS.S_IDLE_60],
      ["Unload models when idle", "if there was neither dictation nor subtitles for the chosen time, the models free the memory; the next dictation loads the model again, and the text appears a few seconds later", "never", "after 5 minutes", "after 10 minutes", "after 30 minutes", "after an hour"]);
    check("the status after an unload tells the truth", [tray("ru"), tray("en")], ["Модели выгружены — следующая диктовка загрузит их заново", "Models unloaded — the next dictation loads them again"]);
    check("no text promises a wake-up any more, and the line about unloading by itself is gone", [/разбуд|wake/i.test(src), "S_STATE_LOADED_SUB" in ruS || "S_STATE_LOADED_SUB" in enS], [false, false]);
  }

  {
    const na = open("en", { lang_models: { auto: "parakeet-v3", ru: "gigaam-v3" }, language: "ru", subtitles_language: "auto", llm_model: "a.gguf" });
    const nw = na.window, nd = nw.document;
    await sleep(700);
    nd.querySelector('.nav[data-p="library"]').click();
    await nw.refreshModels(); await nw.refreshLLM(); await nw.refreshState(); await sleep(150);
    nd.querySelector('#p-library .librow[data-key="asr:gigaam-v3"] button[data-lib="del"]').click(); await sleep(150);
    check("with the auto-detect model not installed the question says dictation in those languages stops", nd.querySelector(".modal p").textContent,
      "The “GigaAM v3” model is used for dictation: русский. Delete it? The Auto-detect model is not installed — dictation in these languages will not work until you download it.");
    nd.querySelector(".modal .btn.ghost").click(); await sleep(200);
    check("and saying no deletes nothing", nw.S.calls.some((c) => c.startsWith("del:")), false);
    na.window.close();
  }

  {
    const old = open("en", { subtitles_translate_on: true, subtitles_translation_only: true, subtitles_layout: "columns", llm_model: "a.gguf" });
    const ow = old.window, od = ow.document;
    await sleep(700);
    const on = (sel) => [...od.querySelectorAll(sel)].map((c) => c.classList.contains("on"));
    check("old settings with only the translation show the third card chosen, on the tab and on Status", [on("#subs_lays .lay"), on("#st_lays .stlay"), od.getElementById("subtitles_translation_only").checked], [[false, false, true], [false, false, true], true]);
    old.window.close();
  }

  const ru = open("ru", { lang_models: { auto: "nemotron-3.5", ru: "gigaam-v3" }, language: "ru", subtitles_language: "auto", llm_model: "a.gguf" });
  const rw = ru.window, rd = rw.document;
  await sleep(700);
  rd.querySelector('.nav[data-p="library"]').click();
  await rw.refreshModels(); await rw.refreshLLM(); await rw.refreshState(); await sleep(150);
  const rrow = (k) => rd.querySelector('#p-library .librow[data-key="' + k + '"]');
  check("Russian: sizes in МБ and ГБ", [rrow("asr:gigaam-v3").querySelector(".psize").textContent, rrow("llm:b.gguf").querySelector(".psize").textContent], ["232 МБ", "1.9 ГБ"]);
  check("Russian: memory in words", rrow("llm:a.gguf").querySelector(".pram").textContent, "≈2.7 ГБ памяти");
  check("Russian: the total counts in Russian", rd.getElementById("lib_have_total").textContent, "8 моделей · 8.0 ГБ");
  rw.S.llm = rw.S.llm.slice(0, 1); rw.S.trm.state = "absent"; await rw.refreshLLM(); await rw.refreshState(); await sleep(150);
  check("Russian: the total follows what is left", rd.getElementById("lib_have_total").textContent, "6 моделей · 3.8 ГБ");
  check("Russian: plural forms", [1, 2, 5, 21, 22].map((n) => rw.pluralCount(n)), ["1 модель", "2 модели", "5 моделей", "21 модель", "22 модели"]);
  check("Russian: the missing one is not installed in words", rrow("trm").querySelector(".lmodel").textContent, "для перевода субтитров");
  check("English: plural forms", [1, 2].map((n) => w.pluralCount(n)), ["1 model", "2 models"]);
  const memLine = async (win, n) => { win.S.installedModels = Array.from({ length: n }, (_, i) => "m" + i); await win.refreshState(); await sleep(120); return win.document.getElementById("state_mem_sub").textContent; };
  check("Russian: the memory line on Status counts the models on disk with the right form", [await memLine(rw, 1), await memLine(rw, 2), await memLine(rw, 5), await memLine(rw, 21)].map((t) => t.split("·")[1].trim().split(",")[0]),
    ["на диске 1 модель", "на диске 2 модели", "на диске 5 моделей", "на диске 21 модель"]);
  {
    const why = "No model for Russian: GigaAM v3 is not installed — download it on the Models tab or pick another model for this language on the Languages tab";
    w.S.ready = false; w.S.status = why; w.S.statusLine = why;
    await w.refreshState(); await sleep(150);
    const st = d.getElementById("st_main");
    check("the status bar gives the same reason as Status when the model is missing, cut with dots and whole in its tip", [st.textContent, st.title, d.getElementById("st_led").classList.contains("on")], [why, why, false]);
    check("the status bar text is cut with an ellipsis instead of pushing the version away", [...d.querySelectorAll("style")].some((s) => s.textContent.includes(".statusbar #st_main{min-width:0;overflow:hidden;text-overflow:ellipsis}")), true);
    w.S.ready = true; w.S.status = ""; w.S.statusLine = "";
    await w.refreshState(); await sleep(150);
    check("and the tip follows the text when all is well again", [st.textContent, st.title], ["Ready", "Ready"]);
  }
  check("English: and so does the English one", [await memLine(w, 1), await memLine(w, 2)].map((t) => t.split("·")[1].trim().split(",")[0]), ["on disk: 1 model", "on disk: 2 models"]);

  {
    const lit = () => [...d.querySelectorAll(".jumphit")];
    const active = () => d.querySelector(".page.active").id;
    const css = [...d.querySelectorAll("style")].map((s) => s.textContent).join(" ");
    check("a jump flash looks like a search hit and fades, still without motion when it is asked for", [css.includes(".jumphit{animation:jumpfade 2.2s ease-in forwards}"), css.includes("@keyframes jumpfade{0%,65%{box-shadow:inset 2px 0 0 var(--hi),inset 0 0 0 9999px var(--navon)}"), css.includes("@media (prefers-reduced-motion:reduce){.jumphit{animation:none;"), css.includes(".jumphit:focus{outline:none}")], [true, true, true, true]);
    w.jumpTo("models:langlist"); await sleep(50);
    const head = d.getElementById("langlist").closest(".card").querySelector(".blklbl");
    check("a jump to a block opens its tab and lights its heading, which takes the focus", [active(), lit().length === 1 && lit()[0] === head, d.activeElement === head], ["p-models", true, true]);
    await sleep(2600);
    check("the light goes out after about two seconds", lit().length, 0);
    w.jumpTo("models:langlist"); await sleep(50);
    check("a second jump to the same place lights it again", lit()[0] === head, true);
    w.jumpTo("dictation:hotkey"); await sleep(50);
    check("a jump to one setting lights its row, not the block, and focuses its control", [active(), lit().length === 1 && lit()[0] === d.getElementById("hotkey").closest(".row"), d.activeElement === d.getElementById("hotkey")], ["p-dictation", true, true]);
    w.jumpTo("post:src_local"); await sleep(50);
    check("a jump to a choice card lights the card", [active(), lit()[0] === d.getElementById("src_local")], ["p-post", true]);
    w.jumpTo("rec:lang=ru"); await sleep(50);
    const ruRow = d.querySelector('#langlist .lrow[data-lang="ru"]');
    check("the program's own way in, a language without its model, opens that language on Languages and lights its row", [active(), lit()[0] === ruRow, ruRow.classList.contains("open")], ["p-models", true, true]);
    w.jumpTo("state"); await sleep(50);
    check("a jump with no place finer than a tab lights nothing", [active(), lit().length], ["p-state", 0]);
    tab("post"); await refresh();
    const plate = d.querySelector("#p-post [data-jump]") || d.querySelector("[data-jump]");
    const plateTab = plate.dataset.jump.split(":")[0];
    plate.click(); await sleep(50);
    check("a jump plate goes through the same jump", [active(), lit().length], ["p-" + plateTab, 1]);
    tab("state"); await w.refreshState(); await sleep(150);
    const go = [...d.querySelectorAll("#state_quick .stcard .stgo")];
    go[0].click(); await sleep(50);
    check("the arrow of a quick card on Status lights the setting it stands for", [active(), lit()[0] === d.getElementById("mic_device").closest(".row")], ["p-mic", true]);
    tab("state"); await sleep(100);
    d.querySelector('[data-goto="subs:subs_main"]').click(); await sleep(50);
    check("the arrow of the subtitles card lights the subtitles block", [active(), lit()[0] === d.getElementById("subs_main").querySelector(".blklbl")], ["p-subs", true]);
    tab("state"); await sleep(100);
    d.getElementById("state_srv_go").click(); await sleep(50);
    check("the server button of a red card lights the recognition server block", [active(), lit()[0] === d.getElementById("stt_srv_card").querySelector(".blklbl")], ["p-system", true]);
    tab("subs"); await sleep(100);
    w.subsGotoModels(); await sleep(50);
    check("Open Languages on the subtitles tab lights the language it is about", [active(), lit()[0] && lit()[0].classList.contains("lrow")], ["p-models", true]);
    tab("help"); await sleep(150);
    const tocLink = d.querySelector("#helptoc a:nth-child(3)");
    if (tocLink) {
      tocLink.click(); await sleep(50);
      check("the help contents go through the same jump", lit()[0] && lit()[0].id === tocLink.getAttribute("href").slice(1), true);
    }
    w.searchSettings("s_subs_window"); await sleep(50);
    check("search lights a block heading it found, and a jump light gives way to it", [d.querySelector(".hit") && d.querySelector(".hit").classList.contains("blklbl"), lit().length], [true, 0]);
    w.searchSettings("s_subs_who_quality"); await sleep(50);
    check("search lights a choice card it found by its title", !!(d.querySelector(".hit") && d.querySelector(".hit").classList.contains("srccard")), true);
    check("found headings, hints and cards wear the same light as found rows", css.includes(".blklbl.hit,.card>.hint.hit,.srccard>.hint.hit,.srccard.hit,"), true);
    w.searchSettings("");
  }

  check("no JavaScript errors", errors, []);
  if (failures.length) {
    console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
    process.exit(1);
  }
  console.log("\nall model checks passed");
  process.exit(0);
})().catch((e) => { console.error("harness crashed:", e.stack || e.message); process.exit(1); });
