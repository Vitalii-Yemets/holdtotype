package main

const setupPage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>{{TITLE}}</title><style>
:root{--bg:#0b0f0c;--panel:#0e1410;--line:#1d4a2b;--green:#3cff6e;--dim:#20a34a;--faint:#14803a;--amber:#ffb347;--glow:0 0 7px rgba(60,255,110,.55)}
*{margin:0;padding:0;box-sizing:border-box}
html,body{height:100%}
body{background:var(--bg);color:var(--green);font:14px Consolas,'Cascadia Mono',monospace;overflow:hidden;display:flex;flex-direction:column}
body::after{content:"";position:fixed;inset:0;pointer-events:none;background:repeating-linear-gradient(0deg,rgba(0,0,0,.10) 0 1px,transparent 1px 3px)}
.header{display:flex;align-items:center;gap:14px;padding:12px 12px 12px 20px;border-bottom:1px solid var(--line);box-shadow:0 1px 12px rgba(60,255,110,.12);cursor:default}
.header h1{font-size:16px;letter-spacing:3px;text-shadow:var(--glow);font-weight:700}
.logo svg{width:40px;height:40px;display:block;filter:drop-shadow(0 0 6px rgba(60,255,110,.6))}
.wave{animation:pulse 1.6s ease-in-out infinite;transform-box:fill-box;transform-origin:center}
@keyframes pulse{0%,100%{opacity:.35;transform:scale(.94)}50%{opacity:1;transform:scale(1)}}
@media (prefers-reduced-motion:reduce){.wave{animation:none;opacity:.75}}
.ver{color:var(--faint);font-size:12px}
button.cap{width:36px;height:30px;background:none;border:1px solid var(--line);color:var(--dim);font:14px Consolas,monospace;cursor:pointer;padding:0}
button.cap:hover{color:var(--green);border-color:var(--dim)}
button.cap.close:hover{color:#ff7b6b;border-color:#7a2e2e}
button:focus-visible,input:focus-visible,select:focus-visible,textarea:focus-visible{outline:1px solid var(--green);outline-offset:2px}
.body{flex:1;padding:18px 22px;display:flex;flex-direction:column;gap:12px}
.tagline{color:var(--dim);font-size:13px;line-height:1.5}
label.fld{color:var(--green);font-size:13px}
input[type=text]{width:100%;padding:8px 11px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;outline:none}
input[type=text]:focus{border-color:var(--dim);box-shadow:var(--glow)}
select{padding:8px 11px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;outline:none;cursor:pointer;max-width:100%;color-scheme:dark}
option{background:#0b100d;color:var(--green)}
option:checked{background:linear-gradient(#123f22,#123f22);color:var(--green)}
select,::picker(select){appearance:base-select}
::picker(select){background:#0b100d;border:1px solid var(--line);padding:2px;margin-top:2px;color:var(--green)}
::picker(select) option{padding:6px 10px;background:none;color:var(--dim);border:0;font:inherit;min-height:0}
::picker(select) option:hover,::picker(select) option:focus{background:#123f22;color:var(--green);outline:none}
::picker(select) option:checked{color:var(--green);text-shadow:var(--glow)}
option::checkmark{display:none}
select::picker-icon{color:var(--faint)}
select:open::picker-icon{transform:rotate(180deg)}
select:open{border-color:var(--dim)}
select:focus{border-color:var(--dim);box-shadow:var(--glow)}
button.ibtn{border:1px solid var(--line);background:none;color:var(--dim);cursor:pointer;padding:0 13px;line-height:0}
button.ibtn:hover{color:var(--green);border-color:var(--dim);box-shadow:var(--glow)}
.warn{color:var(--amber);font-size:12px}
.chk{display:flex;align-items:center;gap:9px;font-size:13px;cursor:pointer}
.chk input{appearance:none;-webkit-appearance:none;width:32px;height:17px;border:1px solid var(--line);background:none;position:relative;flex:none;margin:0;padding:0;cursor:pointer}
.chk input::after{content:"";position:absolute;top:2px;left:2px;width:11px;height:11px;background:var(--dim);transition:.15s}
.chk input:checked{border-color:var(--dim)}
.chk input:checked::after{left:17px;background:var(--green);box-shadow:0 0 7px rgba(60,255,110,.55)}
.chk input:focus-visible{outline:1px solid var(--green);outline-offset:2px}
button.btn{padding:11px 26px;border:1px solid var(--dim);background:#0d1a11;color:var(--green);font:inherit;cursor:pointer;letter-spacing:2px;text-transform:uppercase;font-size:13px}
button.btn.ghost{border-color:var(--line);background:none;color:var(--dim)}
button.btn.ghost:hover{color:var(--green);border-color:var(--dim)}
.foot{gap:8px}
button.btn:hover{background:#123f22;box-shadow:var(--glow)}
.bar{height:16px;border:1px solid var(--line);background:#08100b;position:relative;overflow:hidden}
.bar i{position:absolute;inset:0;width:0;background:linear-gradient(90deg,#123f22,var(--dim));box-shadow:var(--glow);transition:width .2s}
.plog{color:var(--dim);font-size:12px;min-height:16px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.err{color:#ff7b6b;font-size:12px;white-space:pre-wrap;user-select:text}
.done-ic{font-size:34px;text-shadow:var(--glow)}
.pathout{color:var(--dim);font-size:12px;user-select:text;word-break:break-all}
.step{display:none;flex-direction:column;gap:12px}
.step.on{display:flex}
.foot{margin-top:auto;display:flex;justify-content:flex-end}
</style></head><body>
<div class="header" onmousedown="if(event.button===0&&event.target.tagName!=='BUTTON')appDrag()">
 <div class="logo"><svg viewBox="0 0 64 64">
  <rect x="2" y="2" width="60" height="60" rx="12" fill="#0e1410" stroke="#1d4a2b" stroke-width="2"/>
  <g stroke="#3cff6e" stroke-width="4" fill="none" stroke-linecap="round">
   <rect x="26" y="12" width="12" height="20" rx="6" fill="#3cff6e"/>
   <path d="M19 27a13 13 0 0 0 26 0"/>
   <line x1="32" y1="40" x2="32" y2="46"/>
   <line x1="24" y1="49" x2="40" y2="49"/>
  </g>
  <g stroke="#3cff6e" stroke-width="2.5" fill="none" stroke-linecap="round">
   <path class="wave" d="M13 20a17 17 0 0 0 0 14" style="animation-delay:.2s"/>
   <path class="wave" d="M51 20a17 17 0 0 1 0 14" style="animation-delay:.6s"/>
  </g>
 </svg></div>
 <h1>{{APP}}</h1>
 <span class="ver" style="margin-left:auto">v{{VERSION}}</span>
 <button class="cap close" onclick="appClose()">&#10005;</button>
</div>

<div class="body">
 <div class="step on" id="st-opts">
  <div class="tagline">{{TAGLINE}}</div>
  {{UPDNOTE}}
  <label class="fld">{{PATH}}</label>
  <div style="display:flex;gap:8px">
   <input type="text" id="dir" style="flex:1;min-width:0">
   <button class="ibtn" onclick="pickDir()" title="{{PATH}}"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg></button>
  </div>
  <div id="freshopts" style="display:flex;flex-direction:column;gap:12px;{{FRESHSTYLE}}">
  <label class="fld">{{MODEL}}</label>
  <select id="model">{{MODELOPTS}}</select>
  <label class="chk"><input type="checkbox" id="shortcut" checked> {{SHORTCUT}}</label>
  <label class="chk"><input type="checkbox" id="autorun" checked> {{AUTORUN}}</label>
  <label class="chk"><input type="checkbox" id="updates" checked> {{UPDATES}}</label>
  </div>
  <label class="chk"><input type="checkbox" id="launch" checked> {{LAUNCH}}</label>
  <div class="err" id="operr"></div>
  <div class="foot"><button class="btn" onclick="startInstall()">{{INSTALL}}</button></div>
 </div>

 <div class="step" id="st-prog">
  <div class="tagline">{{PROG}}</div>
  <div class="bar"><i id="fill"></i></div>
  <div class="plog" id="plog"></div>
  <div class="foot" id="dlcancelrow" style="display:none"><button class="btn ghost" onclick="stopDownload()">{{DLCANCEL}}</button></div>
  <div class="err" id="perr"></div>
  <div class="foot" id="pfoot" style="display:none">
   <button class="btn" onclick="startInstall()">{{RETRY}}</button>
   <button class="btn ghost" onclick="backToOptions()">{{BACK}}</button>
  </div>
 </div>

 <div class="step" id="st-done">
  <div class="done-ic">&#10003;</div>
  <div><b>{{DONE}}</b></div>
  <div class="tagline">{{DONEAT}}</div>
  <div class="pathout" id="outdir"></div>
  <div class="warn" id="outwarn"></div>
  <div class="foot"><button class="btn" onclick="appFinish(document.getElementById('launch').checked)">{{FINISH}}</button></div>
 </div>
</div>

<script>
document.getElementById("dir").value = {{DEFDIR}};
const UPDATING = {{UPDATING}};
function show(id){
  document.querySelectorAll(".step").forEach(s=>s.classList.toggle("on", s.id===id));
}
async function pickDir(){
  const p = await appBrowse();
  if(p) document.getElementById("dir").value = p;
}
function backToOptions(){
  document.getElementById("pfoot").style.display = "none";
  document.getElementById("perr").textContent = "";
  show("st-opts");
}
function startInstall(){
  const dir = document.getElementById("dir").value.trim();
  if(!dir){
    show("st-opts");
    document.getElementById("operr").textContent = {{NODIR_JS}};
    document.getElementById("dir").focus();
    return;
  }
  document.getElementById("operr").textContent = "";
  document.getElementById("perr").textContent = "";
  document.getElementById("pfoot").style.display = "none";
  document.getElementById("fill").style.width = "0%";
  show("st-prog");
  const model = UPDATING ? "" : document.getElementById("model").value;
  document.getElementById("dlcancelrow").style.display = model ? "" : "none";
  appInstall(dir, document.getElementById("shortcut").checked, document.getElementById("autorun").checked,
    document.getElementById("updates").checked, model, UPDATING);
}
function stopDownload(){
  document.getElementById("dlcancelrow").style.display = "none";
  document.getElementById("plog").textContent = {{DLSTOPPING_JS}};
  appCancelDownload();
}
function setupProgress(pct, name){
  document.getElementById("fill").style.width = pct + "%";
  if(name) document.getElementById("plog").textContent = name;
}
function setupDone(err, warn, dir){
  document.getElementById("dlcancelrow").style.display = "none";
  if(err){
    document.getElementById("perr").textContent = err;
    document.getElementById("pfoot").style.display = "";
    return;
  }
  document.getElementById("outdir").textContent = dir;
  if(warn) document.getElementById("outwarn").textContent = warn === "cancelled" ? {{DLSTOPPED_JS}} : {{WARNMODEL_JS}};
  show("st-done");
}
setTimeout(()=>{ if(window.appReady) appReady(); }, 60);
</script>
</body></html>`
