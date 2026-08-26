"use strict";
const $ = (id) => document.getElementById(id);
function spec() {
  return { n1: Number($("n1").value), n2: Number($("n2").value), core_diameter_um: Number($("core").value), wavelength_nm: Number($("wl").value) };
}
function showError(msg) { $("err-text").textContent = msg; $("err-panel").hidden = false; }
function hideError() { $("err-panel").hidden = true; }
function fill(rows) {
  const tb = $("out-table"); tb.innerHTML = "";
  for (const [k,v] of rows) { const tr = document.createElement("tr"); tr.innerHTML = "<td>"+k+"</td><td>"+v+"</td>"; tb.appendChild(tr); }
}
async function loadExample() {
  hideError();
  try {
    const resp = await fetch("/api/example"); const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "示例失败");
    $("n1").value = data.n1; $("n2").value = data.n2; $("core").value = data.core_diameter_um; $("wl").value = data.wavelength_nm;
    $("hint").textContent = "已加载 SMF-1310 算例。";
  } catch (e) { showError(String(e)); }
}
async function runMode() {
  hideError();
  try {
    const resp = await fetch("/api/mode", { method:"POST", headers:{"Content-Type":"application/json"}, body: JSON.stringify(spec()) });
    const data = await resp.json();
    if (!resp.ok) throw new Error(data.error || "HTTP "+resp.status);
    $("out-panel").hidden = false;
    fill([["NA", data.na.toPrecision(6)], ["V", data.v.toPrecision(6)], ["截止波长 nm", data.cutoff_wavelength_nm.toPrecision(6)], ["单模", String(data.single_mode)], ["D_tot ps/(nm·km)", data.d_total.toPrecision(5)]]);
  } catch (e) { showError(String(e)); }
}
$("btn-example").addEventListener("click", loadExample);
$("btn-mode").addEventListener("click", runMode);
