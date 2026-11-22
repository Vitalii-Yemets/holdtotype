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

const TRM = [
  { id: "gemma-4-e4b", name: "Gemma 4 E4B", file: "gemma-4-E4B_q4_0-it.gguf", size: 4916, vram: 3015, card_need: 4096, topic: true, terms: true, rec: true, accuracy: 5, speed: 5, license: "Apache 2.0" },
  { id: "gemma-3-4b", name: "Gemma 3 4B", file: "gemma-3-4b-it-Q4_K_M.gguf", size: 2374, vram: 3022, card_need: 4096, topic: true, terms: true, rec: false, accuracy: 3, speed: 5, license: "Gemma" },
  { id: "translategemma-4b", name: "TranslateGemma 4B", file: "translategemma-4b-it-Q4_K_M.gguf", size: 2374, vram: 3022, card_need: 4096, topic: false, terms: false, rec: false, accuracy: 4, speed: 5, license: "Gemma" },
  { id: "eurollm-9b", name: "EuroLLM-9B", file: "EuroLLM-9B-Instruct-2512.Q4_K_M.gguf", size: 5324, vram: 5828, card_need: 8192, topic: true, terms: true, rec: false, accuracy: 4, speed: 4, license: "Apache 2.0" },
];

function trmView(S) {
  const st = (id) => (id === "gemma-3-4b" ? S.trm : (S.trms[id] = S.trms[id] || { state: "absent" }));
  const card = S.card === undefined ? 8192 : S.card;
  const rows = TRM.map((r) => Object.assign({}, r, { state: st(r.id).state, pct: st(r.id).pct || 0, err: st(r.id).err || "", fits: !S.cardWhy && card >= r.card_need - 256, reason: S.cardWhy || (card >= r.card_need - 256 ? "" : "small"), own: false, running: !!st(r.id).running }));
  S.llm.forEach((f) => rows.push({ id: "own:" + f.file, name: f.file.replace(/\.gguf$/, ""), file: f.file, size: f.size, vram: f.size + 1000, card_need: f.size + 1256, topic: false, terms: false, rec: false, accuracy: 0, speed: 0, own: true, state: "ready", pct: 0, err: "", fits: card >= f.size + 1000, running: false, license: "" }));
  const chosen = rows.find((r) => r.id === S.chosen && r.state === "ready") || rows.find((r) => r.rec && r.state === "ready") || rows.find((r) => !r.own && r.state === "ready");
  rows.forEach((r) => { r.chosen = !!chosen && r.id === chosen.id; });
  const shown = chosen || TRM.map((r) => rows.find((x) => x.id === r.id)).find((r) => r.id === S.chosen) || rows[0];
  return { id: shown.id, name: shown.name, state: shown.state, pct: shown.pct, err: shown.err, size: shown.size, fits: shown.fits, card: "", reason: shown.reason, card_mb: card, need_mb: shown.card_need, down: "", running: shown.running, models: rows };
}

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
    trms: {}, chosen: "",
    punct: { name: "Horizon Punctuation Small", state: "absent", pct: 0, err: "", size: 289, ram: 340, langs: 89, on: false, running: false },
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
        appTrModel: async (action) => {
          S.calls.push("trm:" + action);
          const at = action.indexOf(":");
          const act = at < 0 ? action : action.slice(0, at);
          const id = (at < 0 ? "" : action.slice(at + 1)) || "gemma-3-4b";
          const s = id === "gemma-3-4b" ? S.trm : (S.trms[id] = S.trms[id] || { state: "absent" });
          if (act === "delete") s.state = "absent";
          if (act === "download") s.state = "downloading";
          if (act === "cancel") s.state = "absent";
          return "";
        },
        appSave: async (f) => { const v = JSON.parse(f); S.saves.push(v); if (typeof v.subtitles_model === "string") S.chosen = v.subtitles_model; if (typeof v.subtitles_sentences === "boolean") S.punct.on = v.subtitles_sentences; return json({ ok: true, severity: "ok", message: "Saved" }); },
        appPunctModel: async (action) => {
          S.calls.push("punct:" + action);
          const p = S.punct;
          if (action.indexOf("download") === 0) { p.state = "downloading"; p.pct = 0; p.err = ""; p.turnOn = p.turnOn || action === "download:on"; }
          if (action === "cancel") { p.state = "absent"; p.turnOn = false; }
          if (action === "delete") { p.state = "absent"; p.running = false; return "Model deleted"; }
          return "";
        },
        appState: async () => json({ subtitles: { phase: "off" }, subs_tr: { choice: "quality", missing: "", pair: "", model: trmView(S), fast: { exe: true, pairs: [], offer: [], used: 0 } },
          subs_place: { position: "bottom", custom: {} }, hotkey: "ctrl+win", mic: "Mic", engine: "", llm: "a.gguf", ram: "", last: "", last_meta: "", ready: S.ready !== false, status: S.status || "Ready",
          status_line: S.statusLine || "Ready", active_model: "Nemotron 3.5", active_state: "ready", active_lang: "RU", assigned: [], installed_models: S.installedModels || [], loaded_now: Object.keys(S.loaded).length ? "Nemotron 3.5" : "",
          llm_ok: true, mic_ok: true, last_at: 0, last_app: "", remote: false, backend_err: "", post_err: "", stt_broken: false, whisper_now: false, enabled: true, autostart: false,
          disk_mb: 0, ram_mb: 16384, ram_free_mb: 9000, badges: { mic: "", models: "", system: "" },
          post_on: true, post_model: "a", post_size: 1066, post_installed: S.llm.length, post_remote: false, post_prompts: 1, punct: Object.assign({}, S.punct) }),
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
    ["asr:gigaam-v3", "asr:medium-q5_0", "asr:local:my-model", "asr:nemotron-3.5", "asr:small", "llm:a.gguf", "llm:b.gguf", "trm:gemma-3-4b"]);
  check("the total counts the models and their size", d.getElementById("lib_have_total").textContent, "8 models · 8.0 GB");
  check("a recognition model names the languages it dictates", use("asr:gigaam-v3"), "used for dictation: русский");
  check("auto-detect, the languages that follow it and the subtitles are named", use("asr:nemotron-3.5"), "used for dictation: auto-detect, en, es, fr, it, pl, uk, and for subtitles");
  check("a language of its own is named for its model", use("asr:small"), "used for dictation: deutsch");
  check("an unused model says so and where it is turned on", [row("asr:medium-q5_0").querySelector(".libidle").textContent, row("asr:medium-q5_0").querySelector(".libturn").textContent],
    ["not in use", "turn on: Languages →"]);
  check("the plate names the place in its tip", (row("asr:medium-q5_0").querySelector(".jump").dataset.tip || row("asr:medium-q5_0").querySelector(".jump").title), "The “Languages” tab, the “Which model serves which language” block");
  check("the editor in use is named for post-processing", use("llm:a.gguf"), "used for post-processing");
  check("an unused editor leads to post-processing", row("llm:b.gguf").querySelector(".libturn").textContent, "turn on: Post-processing →");
  check("the translation model in use is named for translation", use("trm:gemma-3-4b"), "used for translating subtitles");
  check("recognition models carry the info badge, accuracy and speed", [!!row("asr:gigaam-v3").querySelector(".pinfo"), row("asr:gigaam-v3").querySelectorAll(".mbar").length], [true, 2]);
  check("a hand-copied model has no bars", row("asr:local:my-model").querySelectorAll(".mbar").length, 0);
  check("editors and the translation model have no bars either", [row("llm:a.gguf").querySelectorAll(".mbar").length, row("trm:gemma-3-4b").querySelectorAll(".mbar").length], [0, 0]);
  check("the model in memory wears the chip and the unload icon", [row("asr:nemotron-3.5").querySelector(".pchip").textContent, !!row("asr:nemotron-3.5").querySelector('button[data-lib="unload"]')], ["in memory", true]);
  check("a model not in memory has neither", [!!row("asr:small").querySelector(".pchip"), !!row("asr:small").querySelector('button[data-lib="unload"]')], [false, false]);
  check("every installed model has its size and a cross", rows("lib_have_list").every((r) => r.querySelector(".psize").textContent && r.querySelector('button[data-lib="del"]')), true);
  check("the editor in memory says how much memory it takes", row("llm:a.gguf").querySelector(".pram").textContent, "≈2.7 GB RAM");
  check("the translation model has its plates and the terms plate", [...row("trm:gemma-3-4b").querySelectorAll(".libl2 > *")].map((e) => e.textContent), ["graphics card memory: 3 GB", "fits your graphics card", "Terms of use ↗"]);
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

  check("the models to download, by name", rows("lib_get_list").map((r) => r.dataset.key), ["asr:base", "asr:moonshine-uk", "asr:parakeet-v3", "trm:gemma-4-e4b", "trm:translategemma-4b", "trm:eurollm-9b", "punct"]);
  check("the hint says what they are for", d.getElementById("lib_get_hint").textContent, "For dictation and subtitles — models that turn speech into text. For translating subtitles — translation models, we recommend Gemma 4 E4B. A downloaded model moves to “Installed models”.");
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
  row("trm:gemma-3-4b").querySelector('button[data-lib="del"]').click(); await sleep(150);
  check("deleting the translation model in use names its use", modalText(), "The “Gemma 3 4B” model is used for translating subtitles. Delete it? Another model will have to be chosen there.");
  await answer(true); await refresh(); await sleep(200);
  check("the translation model goes to the ones to download", [!!d.querySelector('#lib_get_list .librow[data-key="trm:gemma-3-4b"]'), use("trm:gemma-3-4b")], [true, "for translating subtitles: the earlier model"]);
  check("and the hint names the recommended one", d.getElementById("lib_get_hint").textContent.includes("For translating subtitles — translation models, we recommend Gemma 4 E4B."), true);
  tab("subs"); await sleep(200);
  check("the quality card says the recommended one is not installed", d.getElementById("subs_q_size").textContent, "not installed · 4.8 GB");
  tab("library"); await refresh();
  row("trm:gemma-3-4b").querySelector('button[data-lib="dl"]').click(); await sleep(200);
  check("the arrow downloads that translation model", S.calls.includes("trm:download:gemma-3-4b"), true);

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
  check("Russian: the missing one is not installed in words", rrow("trm:gemma-3-4b").querySelector(".lmodel").textContent, "для перевода субтитров: прежняя модель");
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

  for (const keep of [true, false]) {
    const sdom = open("en", { subtitles_save: keep, subtitles_translate_on: true, subtitles_translator: "fast" });
    const sw = sdom.window, sd = sw.document, SS = sw.S;
    const sweepErrors = [];
    sw.addEventListener("error", (e) => sweepErrors.push(e.message));
    await sleep(700);
    const said = () => SS.saves.map((f) => f.subtitles_save);
    check("save switch " + keep + ": both switches show the setting when the window opens", [sd.getElementById("subtitles_save").checked, sd.getElementById("st_save_on").checked], [keep, keep]);
    const box = sd.getElementById("subtitles_save");
    box.checked = !keep;
    await sw.doSave();
    box.checked = keep;
    check("save switch " + keep + ": a save does not take the value from a box the page has not filled", SS.saves[SS.saves.length - 1].subtitles_save, keep);
    const own = new Set(["subtitles_save", "st_save_on"]);
    const pick = () => [...sd.querySelectorAll("#p-state input, #p-state select, #p-state textarea, #p-subs input, #p-subs select, #p-subs textarea")].filter((el) => !own.has(el.id));
    const before = SS.saves.length;
    let touched = 0;
    for (const el of pick()) {
      if (!el.isConnected) continue;
      if (el.type === "checkbox") el.checked = !el.checked;
      else if (el.type === "radio") el.checked = true;
      else if (el.tagName === "SELECT") { if (el.options.length > 1) el.selectedIndex = (el.selectedIndex + 1) % el.options.length; }
      else el.value = el.value ? el.value : "1";
      el.dispatchEvent(new sw.Event("input", { bubbles: true }));
      el.dispatchEvent(new sw.Event("change", { bubbles: true }));
      touched++;
      await sleep(140);
    }
    for (const el of [...sd.querySelectorAll("#p-state button, #p-subs button")]) {
      if (!el.isConnected) continue;
      el.click();
      touched++;
      await sleep(140);
    }
    await sleep(500);
    const sent = said().slice(before);
    check("save switch " + keep + ": every control of Status and Subtitles & translation was touched and saves were sent", [touched > 40, sent.length > 10], [true, true]);
    check("save switch " + keep + ": no save sent while touching them changes it", sent.filter((v) => v !== keep), []);
    check("save switch " + keep + ": both switches still show it", [sd.getElementById("subtitles_save").checked, sd.getElementById("st_save_on").checked], [keep, keep]);
    check("save switch " + keep + ": touching every control raised no JavaScript error", sweepErrors, []);
    sdom.window.close();
  }

  {
    const quality = { subtitles_translate_on: true, subtitles_translator: "quality", subtitles_language: "en", llm_model: "a.gguf" };
    const view = async (setup, cfg) => {
      const v = open("en", Object.assign({}, quality, cfg || {}), setup);
      const vw = v.window, vd = v.window.document;
      await sleep(700);
      vd.querySelector('.nav[data-p="subs"]').click();
      await vw.refreshModels(); await vw.refreshLLM(); await vw.refreshState(); await sleep(200);
      return { v, vw, vd, S: vw.S };
    };
    const cards = (vd) => [...vd.querySelectorAll("#subs_q_models .pcard")];
    const g6v = (vd, id) => vd.querySelector('#lib_get_list .librow[data-key="trm:' + id + '"]');
    const card = (vd, id) => vd.querySelector('#subs_q_models .pcard[data-id="' + id + '"]');
    const sw = (vd, id) => card(vd, id).querySelector("input.psw");
    const texts = (el, sel) => [...el.querySelectorAll(sel)].map((e) => e.textContent);

    let t = await view((S) => { S.trms["gemma-4-e4b"] = { state: "ready" }; S.trms["translategemma-4b"] = { state: "ready" }; });
    let vd = t.vd;
    const quals = [...vd.querySelectorAll("#subs_src_quality > .hint")].map((e) => e.textContent);
    check("models: the quality card says the topic and the terms are found by the translator itself", quals, ["S_SUBS_Q_HINT", "S_SUBS_Q_AUTO"]);
    check("models: the old plates of the card are gone", [!!vd.getElementById("subs_q_name"), !!vd.getElementById("subs_trm_terms"), !!vd.getElementById("subs_q_note"), vd.querySelectorAll("#subs_src_quality .sum .pmt").length], [false, false, false, 0]);
    check("models: the downloaded models and the own files are cards of the Languages scheme", cards(vd).map((c) => c.dataset.id), ["gemma-4-e4b", "gemma-3-4b", "translategemma-4b", "own:a.gguf", "own:b.gguf"]);
    const g4 = card(vd, "gemma-4-e4b");
    check("models: the recommended one is chosen: the card is current, its switch on, a tick and the chip", [g4.classList.contains("cur"), sw(vd, "gemma-4-e4b").checked, !!g4.querySelector(".pcur"), texts(g4, ".pchip")], [true, true, true, ["recommended"]]);
    check("models: a card has the parts of the Languages cards", ["ptop", "pname", "pnm", "pinfo", "mbars", "pmeta", "pact"].map((c) => !!g4.querySelector("." + c)), [true, true, true, true, true, true, true]);
    check("models: its tip, bars, plates and size", [g4.querySelector(".pinfo").dataset.tip, g4.querySelectorAll(".mbar").length, texts(g4, ".pmeta .pmt"), g4.querySelector(".pram").textContent], ["more precise with terms", 2, ["graphics card memory: 3 GB"], "4.8 GB on disk"]);
    const tg = card(vd, "translategemma-4b");
    check("models: TranslateGemma says it ignores the topic", [tg.querySelector(".pinfo").dataset.tip, texts(tg, ".pmeta .pmt")], ["the smoothest text; ignores the topic of the conversation", ["graphics card memory: 3 GB", "ignores the topic of the conversation"]]);
    const own = card(vd, "own:a.gguf");
    check("models: an own model has no bars and no memory plate, and says it is not checked", [own.querySelector(".pinfo").dataset.tip, own.querySelectorAll(".mbar").length, texts(own, ".pmeta .pmt"), own.querySelector(".pdesc").textContent], ["the file a.gguf from the models folder", 0, [], "not checked: the translation may not work"]);
    check("models: under the cards a line leads to the models to download", [vd.querySelector("#subs_q_models .jumpline span").textContent, vd.querySelector("#subs_q_models .jumpline .jump").textContent, vd.querySelector("#subs_q_models .jumpline .jump").dataset.jump], ["Download another model:", "Models: Models to download →", "library:lib_get"]);
    const saves = t.S.saves.length;
    sw(vd, "gemma-4-e4b").checked = false; sw(vd, "gemma-4-e4b").dispatchEvent(new t.vw.Event("change", { bubbles: true })); await sleep(300);
    check("models: the chosen one cannot be switched off", [sw(vd, "gemma-4-e4b").checked, t.S.saves.length], [true, saves]);
    sw(vd, "gemma-3-4b").checked = true; sw(vd, "gemma-3-4b").dispatchEvent(new t.vw.Event("change", { bubbles: true })); await sleep(400);
    check("models: switching on another one saves subtitles_model", t.S.saves[t.S.saves.length - 1].subtitles_model, "gemma-3-4b");
    await t.vw.refreshState(); await sleep(200);
    check("models: and it becomes the only one on", cards(vd).filter((c) => c.querySelector("input.psw").checked).map((c) => c.dataset.id), ["gemma-3-4b"]);
    sw(vd, "own:a.gguf").checked = true; sw(vd, "own:a.gguf").dispatchEvent(new t.vw.Event("change", { bubbles: true })); await sleep(400);
    await t.vw.refreshState(); await sleep(200);
    check("models: an own model can be chosen too", [t.S.saves[t.S.saves.length - 1].subtitles_model, sw(vd, "own:a.gguf").checked], ["own:a.gguf", true]);
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    const lrow = (k) => vd.querySelector('#p-library .librow[data-key="' + k + '"]');
    check("models: the own model chosen for translation says so on Models", lrow("llm:a.gguf").querySelector(".lmodel").textContent, "used for post-processing; used for translating subtitles");
    lrow("llm:a.gguf").querySelector('button[data-lib="del"]').click(); await sleep(150);
    check("models: and deleting it names the translation", vd.querySelector(".modal p").textContent, "The “a.gguf” model is used for post-processing; used for translating subtitles. Delete it? Another model will have to be chosen there.");
    vd.querySelector(".modal .btn.ghost").click(); await sleep(300);
    t.v.window.close();

    t = await view((S) => { S.trms["gemma-4-e4b"] = { state: "ready" }; S.trms["eurollm-9b"] = { state: "downloading", pct: 37 }; S.llm = []; });
    vd = t.vd;
    const eu = card(vd, "eurollm-9b");
    check("models: a model being downloaded is a card with its percent, a cross and a switch that cannot be used", [eu.querySelector(".mpct").textContent, !!eu.querySelector("button[data-a=cancel]"), sw(vd, "eurollm-9b").disabled], ["37%", true, true]);
    eu.querySelector("button[data-a=cancel]").click(); await sleep(200);
    check("models: the cross cancels that download", t.S.calls.includes("trm:cancel:eurollm-9b"), true);
    t.v.window.close();

    t = await view((S) => { S.card = 4096; S.trms["gemma-4-e4b"] = { state: "ready" }; S.trms["eurollm-9b"] = { state: "ready" }; S.trm.state = "absent"; S.llm = []; });
    vd = t.vd;
    const big = card(vd, "eurollm-9b");
    check("models: a model that does not fit has a switch that cannot be used and says why", [sw(vd, "eurollm-9b").disabled, big.querySelector(".pdesc").textContent], [true, "does not fit: it needs a graphics card with 8 GB or more, yours has 4 GB"]);
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    const lr = (k) => vd.querySelector('#p-library .librow[data-key="' + k + '"]');
    check("models: on Models the installed one that does not fit says so too", texts(lr("trm:eurollm-9b"), ".libl2 > *"), ["graphics card memory: 6 GB", "does not fit: it needs a graphics card with 8 GB or more, yours has 4 GB", "Terms of use ↗"]);
    check("models: and it does not lead to the subtitles card, where it cannot be switched on", [lr("trm:eurollm-9b").querySelector(".libidle").textContent, !!lr("trm:eurollm-9b").querySelector(".libturn")], ["not in use", false]);
    check("models: a model to download that fits has the arrow, one that does not has none", [!!lr("trm:gemma-3-4b").querySelector('button[data-lib="dl"]'), lr("trm:gemma-3-4b").parentNode.id], [true, "lib_get_list"]);
    t.v.window.close();

    t = await view((S) => { S.card = 6144; S.trm.state = "absent"; S.llm = []; });
    vd = t.vd;
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    const g6 = (k) => vd.querySelector('#lib_get_list .librow[data-key="' + k + '"]');
    check("models: EuroLLM on a 6 GB card names the card it needs, the others fit", [texts(g6("trm:eurollm-9b"), ".libl2 > *"), !!g6("trm:eurollm-9b").querySelector('button[data-lib="dl"]'), !!g6("trm:gemma-4-e4b").querySelector('button[data-lib="dl"]')], [["graphics card memory: 6 GB", "does not fit: it needs a graphics card with 8 GB or more, yours has 6 GB", "Terms of use ↗"], false, true]);
    vd.querySelector('.nav[data-p="subs"]').click(); await sleep(200);
    check("models: and the quality card still offers the recommended one", !!vd.getElementById("subs_q_get"), true);
    t.v.window.close();

    t = await view((S) => { S.card = 0; S.cardWhy = "integrated"; S.trm.state = "absent"; S.llm = []; });
    vd = t.vd;
    check("models: no model and integrated graphics: no download button, the reason, no plate to the models to download", [!!vd.getElementById("subs_q_get"), vd.querySelector("#subs_q_models .nowhy").textContent, !!vd.querySelector("#subs_q_models .jumpline")], [false, "This computer cannot run it: it has only integrated graphics, which share the computer's memory.", false]);
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    check("models: and on Models no translation model can be downloaded, each says why", ["gemma-4-e4b", "eurollm-9b"].map((id) => [!!g6v(vd, id).querySelector('button[data-lib="dl"]'), g6v(vd, id).querySelector(".libwarn").textContent]), [[false, "This computer cannot run it: it has only integrated graphics, which share the computer's memory."], [false, "This computer cannot run it: it has only integrated graphics, which share the computer's memory."]]);
    t.v.window.close();

    t = await view((S) => { S.card = 3072; S.trm.state = "absent"; S.llm = []; });
    vd = t.vd;
    check("models: no model and a 3 GB card: the numbers come from the model, no plate to the models to download", [!!vd.getElementById("subs_q_get"), !!vd.querySelector("#subs_q_models .jumpline"), vd.querySelector("#subs_q_models .nowhy").textContent], [false, false, "This computer cannot run it: it needs a graphics card with 4 GB of its own memory or more, yours has 3 GB."]);
    t.v.window.close();

    t = await view((S) => { S.trm.state = "absent"; S.llm = []; });
    vd = t.vd;
    check("models: with no model the card says so and offers the recommended one", [vd.querySelector("#subs_q_models .postnone .nonote").textContent, vd.getElementById("subs_q_get").textContent, vd.querySelector("#subs_q_models .jumpline .jump").textContent], ["There is no translation model yet.", "Download Gemma 4 E4B · 4.8 GB", "Models: Models to download →"]);
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    const gl = (k) => vd.querySelector('#lib_get_list .librow[data-key="' + k + '"]');
    check("models: every translation model is to download, each with its use, plates and size", ["gemma-4-e4b", "gemma-3-4b", "translategemma-4b", "eurollm-9b"].map((id) => [gl("trm:" + id).querySelector(".lmodel").textContent, texts(gl("trm:" + id), ".libl2 > *"), gl("trm:" + id).querySelector(".psize").textContent]), [
      ["for translating subtitles: more precise with terms", ["recommended", "graphics card memory: 3 GB", "fits your graphics card", "Terms of use ↗"], "4.8 GB"],
      ["for translating subtitles: the earlier model", ["graphics card memory: 3 GB", "fits your graphics card", "Terms of use ↗"], "2.3 GB"],
      ["for translating subtitles: the smoothest text; ignores the topic of the conversation", ["graphics card memory: 3 GB", "fits your graphics card", "Terms of use ↗"], "2.3 GB"],
      ["for translating subtitles: European languages; needs a graphics card with 8 GB or more", ["graphics card memory: 6 GB", "fits your graphics card", "Terms of use ↗"], "5.2 GB"]]);
    check("models: the tip of a row to download names the model", gl("trm:gemma-4-e4b").querySelector(".pinfo").dataset.tip, "a translation model on the graphics card: more precise with terms");
    const before = t.S.calls.length;
    for (const id of ["gemma-4-e4b", "gemma-3-4b", "translategemma-4b", "eurollm-9b"]) { gl("trm:" + id).querySelector("[data-terms]").click(); await sleep(50); }
    check("models: each terms plate opens the terms of its own model", t.S.calls.slice(before), ["trm:terms:gemma-4-e4b", "trm:terms:gemma-3-4b", "trm:terms:translategemma-4b", "trm:terms:eurollm-9b"]);
    vd.querySelector('.nav[data-p="subs"]').click(); await sleep(200);
    vd.getElementById("subs_q_get").click(); await sleep(200);
    check("models: the button downloads the recommended model", t.S.calls.includes("trm:download:gemma-4-e4b"), true);
    await t.vw.refreshState(); await sleep(200);
    check("models: and it becomes a card with its percent", [cards(vd).map((c) => c.dataset.id), sw(vd, "gemma-4-e4b").disabled], [["gemma-4-e4b"], true]);
    t.v.window.close();

    t = await view((S) => { S.trms["gemma-4-e4b"] = { state: "ready" }; S.llm = []; });
    vd = t.vd;
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    const hr = (k) => vd.querySelector('#lib_have_list .librow[data-key="' + k + '"]');
    check("models: the chosen installed one is used for translating, the other is not and leads to the subtitles card", [hr("trm:gemma-4-e4b").querySelector(".lmodel").textContent, hr("trm:gemma-3-4b").querySelector(".libidle").textContent, hr("trm:gemma-3-4b").querySelector(".libturn .jump").textContent], ["used for translating subtitles", "not in use", "Subtitles & translation →"]);
    check("models: the total counts every installed translation model", vd.getElementById("lib_have_total").textContent, "7 models · 9.8 GB");
    hr("trm:gemma-4-e4b").querySelector('button[data-lib="del"]').click(); await sleep(150);
    check("models: deleting the one in use names its use", vd.querySelector(".modal p").textContent, "The “Gemma 4 E4B” model is used for translating subtitles. Delete it? Another model will have to be chosen there.");
    vd.querySelector(".modal .btn.yes").click(); await sleep(300);
    check("models: and it is deleted by its id", t.S.calls.includes("trm:delete:gemma-4-e4b"), true);
    t.v.window.close();

    const sent = (vd) => ({
      row: vd.getElementById("subs_sent_row").className,
      sw: [vd.getElementById("subtitles_sentences").style.display !== "none", vd.getElementById("subtitles_sentences").disabled, vd.getElementById("subtitles_sentences").checked],
      dl: vd.getElementById("subs_sent_dl").style.display !== "none" ? [vd.getElementById("subs_sent_pct").textContent, vd.getElementById("subs_sent_x").title] : null,
      none: vd.getElementById("subs_sent_none") ? [vd.querySelector("#subs_sent_none .nonote").textContent, vd.querySelector("#subs_sent_none button").textContent] : null,
      last: vd.getElementById("subs_sent").lastElementChild.id,
    });
    t = await view((S) => { S.llm = []; });
    vd = t.vd;
    const label = vd.querySelector("#subs_sent_row > label");
    check("sentences: the row stands right under the translation language", [vd.getElementById("subs_tr_lang_row").nextElementSibling.firstElementChild.id, label.firstChild.textContent, label.querySelector(".sub").textContent], ["subs_sent_row", "S_SUBS_SENT", "S_SUBS_SENT_SUB"]);
    check("sentences: without the model the row says so and offers the download, the switch is hidden", sent(vd), { row: "row noline", sw: [false, true, false], dl: null, none: ["There is no model for this yet.", "Download Horizon Punctuation Small · 289 MB"], last: "subs_sent_none" });
    check("sentences: the empty state uses the general empty-state rule of the window", [vd.getElementById("subs_sent_none").className, vd.querySelector("#subs_sent_none .acts .mini.pri") !== null, vd.getElementById("subs_sent").parentNode.lastElementChild.id], ["postnone", true, "subs_sent"]);
    vd.querySelector("#subs_sent_none button").click(); await sleep(150);
    check("sentences: the button downloads the model and asks to turn the setting on after it", t.S.calls.filter((c) => c.indexOf("punct:") === 0), ["punct:download:on"]);
    await t.vw.refreshState(); await sleep(150);
    check("sentences: while it downloads the row shows the percent and a cross", sent(vd), { row: "row", sw: [false, true, false], dl: ["…", "Cancel the download"], none: null, last: "subs_sent_row" });
    t.S.punct.pct = 37; await t.vw.refreshState(); await sleep(150);
    check("sentences: the percent follows the download", sent(vd).dl, ["37%", "Cancel the download"]);
    vd.getElementById("subs_sent_x").click(); await sleep(150);
    check("sentences: the cross cancels the download and the row goes back to the offer", [t.S.calls.slice(-1)[0], sent(vd).none], ["punct:cancel", ["There is no model for this yet.", "Download Horizon Punctuation Small · 289 MB"]]);
    vd.querySelector("#subs_sent_none button").click(); await sleep(150);
    t.S.punct.state = "ready"; if (t.S.punct.turnOn) t.S.punct.on = true;
    await t.vw.refreshState(); await sleep(150);
    check("sentences: after the download the switch is there and on", sent(vd), { row: "row", sw: [true, false, true], dl: null, none: null, last: "subs_sent_row" });
    const savesBefore = t.S.saves.length;
    vd.getElementById("subtitles_lines").value = "4";
    vd.getElementById("subtitles_lines").dispatchEvent(new vd.defaultView.Event("change", { bubbles: true }));
    await sleep(400);
    check("sentences: a save of another setting leaves it alone", t.S.saves.slice(savesBefore).map((v) => v.subtitles_sentences), [null]);
    const swEl = vd.getElementById("subtitles_sentences");
    swEl.checked = false;
    swEl.dispatchEvent(new vd.defaultView.Event("change", { bubbles: true }));
    await sleep(400);
    check("sentences: turning the switch off saves subtitles_sentences right away", [t.S.saves.slice(-1)[0].subtitles_sentences, t.S.punct.on], [false, false]);
    await t.vw.refreshState(); await sleep(150);
    check("sentences: and the switch stays off", sent(vd).sw, [true, false, false]);
    swEl.checked = true;
    swEl.dispatchEvent(new vd.defaultView.Event("change", { bubbles: true }));
    await sleep(400);
    check("sentences: turning it on saves it too", t.S.saves.slice(-1)[0].subtitles_sentences, true);
    vd.getElementById("subtitles_translate_on").click(); await sleep(400);
    check("sentences: with translation off the row is dimmed and the switch cannot be used", [sent(vd).row, sent(vd).sw], ["row dimmed", [true, true, true]]);
    vd.getElementById("subtitles_translate_on").click(); await sleep(400);

    t.S.punct.running = true;
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    await t.vw.refreshState(); await sleep(150);
    const ph = vd.querySelector('#lib_have_list .librow[data-key="punct"]');
    check("sentences: on Models the installed model is among the installed ones, used, with its plates, size and tip", [ph.querySelector(".pnm").textContent, ph.querySelector(".lmodel").textContent, texts(ph, ".libl2 > *"), ph.querySelector(".psize").textContent, ph.querySelector(".pinfo").dataset.tip, !!ph.querySelector('button[data-lib="del"]'), !!vd.querySelector('#lib_get_list .librow[data-key="punct"]')],
      ["Horizon Punctuation Small", "used for subtitles: finds where a sentence ends", ["in memory", "languages: 89", "≈340 MB RAM", "Terms of use ↗"], "289 MB", "a punctuation model on the processor: finds by meaning where a sentence ends, in 89 languages", true, false]);
    check("sentences: the total counts it", vd.getElementById("lib_have_total").textContent, "7 models · 5.3 GB");
    ph.querySelector("[data-terms]").click(); await sleep(50);
    check("sentences: its terms plate opens its own terms", t.S.calls.slice(-1)[0], "punct:terms");
    ph.querySelector('button[data-lib="del"]').click(); await sleep(150);
    check("sentences: deleting it while the setting is on names its use", vd.querySelector(".modal p").textContent, "The “Horizon Punctuation Small” model is used for subtitles: finds where a sentence ends. Delete it? Subtitles will be translated by pauses in speech again.");
    vd.querySelector(".modal .btn.yes").click(); await sleep(400);
    const pg = vd.querySelector('#lib_get_list .librow[data-key="punct"]');
    check("sentences: after the delete it is a model to download with the arrow", [t.S.calls.slice(-1)[0], !!pg, pg && pg.querySelector(".lmodel").textContent, pg && !!pg.querySelector('button[data-lib="dl"]'), !!vd.querySelector('#lib_have_list .librow[data-key="punct"]')], ["punct:delete", true, "for subtitles: finds where a sentence ends", true, false]);
    vd.querySelector('.nav[data-p="subs"]').click(); await sleep(200);
    check("sentences: and the subtitles row is back to the offer", sent(vd), { row: "row noline", sw: [false, true, true], dl: null, none: ["There is no model for this yet.", "Download Horizon Punctuation Small · 289 MB"], last: "subs_sent_none" });
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    vd.querySelector('#lib_get_list .librow[data-key="punct"] button[data-lib="dl"]').click(); await sleep(150);
    check("sentences: the arrow on Models downloads it without turning the setting on", t.S.calls.slice(-1)[0], "punct:download");
    await t.vw.refreshState(); await sleep(150);
    check("sentences: on Models it shows its percent and a cross while it downloads", [vd.querySelector('#lib_get_list .librow[data-key="punct"] .mpct').textContent, !!vd.querySelector('#lib_get_list .librow[data-key="punct"] button[data-lib="cancel"]')], ["…", true]);
    vd.querySelector('#lib_get_list .librow[data-key="punct"] button[data-lib="cancel"]').click(); await sleep(150);
    check("sentences: and the cross there cancels it", t.S.calls.slice(-1)[0], "punct:cancel");
    t.S.punct.err = "no connection"; await t.vw.refreshState(); await sleep(150);
    check("sentences: a failed download is named on Models and in the row", [vd.querySelector('#lib_get_list .librow[data-key="punct"] .perr').textContent, vd.querySelector("#subs_sent_none .nonote").textContent], ["Not downloaded: no connection", "Not downloaded: no connection"]);
    t.v.window.close();

    t = await view((S) => { S.llm = []; S.punct.state = "ready"; S.punct.on = false; });
    vd = t.vd;
    vd.querySelector('.nav[data-p="library"]').click(); await sleep(200);
    const pi = vd.querySelector('#lib_have_list .librow[data-key="punct"]');
    check("sentences: installed but turned off it says what it is for", pi.querySelector(".lmodel").textContent, "for subtitles: finds where a sentence ends");
    pi.querySelector('button[data-lib="del"]').click(); await sleep(150);
    check("sentences: and deleting it asks the plain question", vd.querySelector(".modal p").textContent, "Delete the \"Horizon Punctuation Small\" model? It can be downloaded again.");
    t.v.window.close();
  }

  check("no JavaScript errors", errors, []);
  if (failures.length) {
    console.error(`\n${failures.length} check(s) failed: ${failures.join(", ")}`);
    process.exit(1);
  }
  console.log("\nall model checks passed");
  process.exit(0);
})().catch((e) => { console.error("harness crashed:", e.stack || e.message); process.exit(1); });
