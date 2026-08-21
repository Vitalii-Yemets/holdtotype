package main

const setupPage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>{{TITLE}}</title><style>
:root{--bg:#0b0f0c;--panel:#0e1410;--line:#1d4a2b;--green:#3cff6e;--dim:#20a34a;--faint:#14803a;--amber:#ffb347;--glow:0 0 7px rgba(60,255,110,.55)}
*{margin:0;padding:0;box-sizing:border-box}
html,body{height:100%}
body{background:var(--bg);color:var(--green);font:14px Consolas,'Cascadia Mono',monospace;overflow:hidden;display:flex;flex-direction:column}
body::after{content:"";position:fixed;inset:0;pointer-events:none;background:repeating-linear-gradient(0deg,rgba(0,0,0,.16) 0 1px,transparent 1px 3px)}
.header{display:flex;align-items:center;gap:14px;padding:12px 12px 12px 20px;border-bottom:1px solid var(--line);box-shadow:0 1px 12px rgba(60,255,110,.12);cursor:default}
.header h1{font-size:16px;letter-spacing:3px;text-shadow:var(--glow);font-weight:700}
.logo svg{width:40px;height:40px;display:block;filter:drop-shadow(0 0 6px rgba(60,255,110,.6))}
.ver{color:var(--faint);font-size:12px}
button.cap{width:36px;height:30px;background:none;border:1px solid var(--line);color:var(--dim);font:14px Consolas,monospace;cursor:pointer;padding:0}
button.cap:hover{color:var(--green);border-color:var(--dim)}
button.cap.close:hover{color:#ff7b6b;border-color:#7a2e2e}
.body{flex:1;padding:18px 22px;display:flex;flex-direction:column;gap:12px}
.tagline{color:var(--faint);font-size:13px;line-height:1.5}
label.fld{color:var(--green);font-size:13px}
input[type=text]{width:100%;padding:8px 11px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;outline:none}
input[type=text]:focus{border-color:var(--dim);box-shadow:var(--glow)}
select{padding:8px 11px;border:1px solid var(--line);background:#08100b;color:var(--green);font:inherit;outline:none;cursor:pointer;max-width:100%}
select:focus{border-color:var(--dim);box-shadow:var(--glow)}
button.ibtn{border:1px solid var(--line);background:none;color:var(--dim);cursor:pointer;padding:0 13px;line-height:0}
button.ibtn:hover{color:var(--green);border-color:var(--dim);box-shadow:var(--glow)}
.warn{color:var(--amber);font-size:12px}
.chk{display:flex;align-items:center;gap:9px;font-size:13px;cursor:pointer}
.chk input{width:16px;height:16px;accent-color:var(--dim);cursor:pointer}
button.btn{padding:11px 26px;border:1px solid var(--dim);background:#0d1a11;color:var(--green);font:inherit;cursor:pointer;letter-spacing:2px;text-transform:uppercase;font-size:13px}
button.btn:hover{background:#123f22;box-shadow:var(--glow)}
.bar{height:16px;border:1px solid var(--line);background:#08100b;position:relative;overflow:hidden}
.bar i{position:absolute;inset:0;width:0;background:linear-gradient(90deg,#123f22,var(--dim));box-shadow:var(--glow);transition:width .2s}
.plog{color:var(--faint);font-size:12px;min-height:16px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.err{color:#ff7b6b;font-size:12px;white-space:pre-wrap;user-select:text}
.done-ic{font-size:34px;text-shadow:var(--glow)}
.pathout{color:var(--dim);font-size:12px;user-select:text;word-break:break-all}
.step{display:none;flex-direction:column;gap:12px}
.step.on{display:flex}
.foot{margin-top:auto;display:flex;justify-content:flex-end}
</style></head><body>
<div class="header" onmousedown="if(event.button===0&&event.target.tagName!=='BUTTON')appDrag()">
 <div class="logo"><svg viewBox="0 0 64 64" fill="none" stroke="#3cff6e" stroke-width="3" stroke-linecap="round">
  <rect x="24" y="8" width="16" height="26" rx="8"/><path d="M16 26a16 16 0 0 0 32 0"/>
  <line x1="32" y1="42" x2="32" y2="50"/><line x1="23" y1="53" x2="41" y2="53"/></svg></div>
 <h1>VOX&nbsp;TERMINAL</h1>
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
  <label class="chk"><input type="checkbox" id="autorun"> {{AUTORUN}}</label>
  </div>
  <label class="chk"><input type="checkbox" id="launch" checked> {{LAUNCH}}</label>
  <div class="foot"><button class="btn" onclick="startInstall()">{{INSTALL}}</button></div>
 </div>

 <div class="step" id="st-prog">
  <div class="tagline">{{PROG}}</div>
  <div class="bar"><i id="fill"></i></div>
  <div class="plog" id="plog"></div>
  <div class="err" id="perr"></div>
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
function startInstall(){
  const dir = document.getElementById("dir").value.trim();
  if(!dir) return;
  show("st-prog");
  appInstall(dir, document.getElementById("shortcut").checked, document.getElementById("autorun").checked,
    UPDATING ? "" : document.getElementById("model").value, UPDATING);
}
function setupProgress(pct, name){
  document.getElementById("fill").style.width = pct + "%";
  if(name) document.getElementById("plog").textContent = name;
}
function setupDone(err, warn, dir){
  if(err){
    document.getElementById("perr").textContent = err;
    return;
  }
  document.getElementById("outdir").textContent = dir;
  if(warn) document.getElementById("outwarn").textContent = {{WARNMODEL_JS}};
  show("st-done");
}
setTimeout(()=>{ if(window.appReady) appReady(); }, 60);
</script>
</body></html>`
