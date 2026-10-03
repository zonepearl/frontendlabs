// home.js: search across every guide and chapter, theme toggle, and
// "continue reading" from the reader's saved positions.
(function () {
  "use strict";
  function load(key, fallback) {
    try { var v = localStorage.getItem(key); return v ? JSON.parse(v) : fallback; } catch (e) { return fallback; }
  }
  var doc = document.documentElement;
  var prefs = load("fl:prefs", {});
  var THEMES = ["paper", "light", "night"];
  var themeBtn = document.getElementById("theme");
  function label() { themeBtn.textContent = doc.dataset.theme[0].toUpperCase() + doc.dataset.theme.slice(1); }
  function cycle() {
    prefs.theme = THEMES[(THEMES.indexOf(doc.dataset.theme) + 1) % THEMES.length];
    doc.dataset.theme = prefs.theme;
    try { localStorage.setItem("fl:prefs", JSON.stringify(prefs)); } catch (e) {}
    label();
  }
  label();
  themeBtn.addEventListener("click", cycle);

  // ---------------------------------------------------- continue reading --
  var recent = load("fl:recent", []);
  if (recent.length) {
    var list = document.getElementById("continue-list");
    recent.slice(0, 4).forEach(function (r) {
      var a = document.createElement("a");
      a.href = r.page + "#" + r.ch;
      a.className = "flex items-center gap-3 rounded-xl border border-rule bg-card px-4 py-3 no-underline shadow-sm hover:border-accent";
      a.innerHTML = '<span class="text-2xl"></span><span class="min-w-0"><span class="block text-sm font-semibold text-ink"></span><span class="block max-w-xs truncate text-xs text-muted"></span></span>';
      var icon = r.icon || "";
      if (icon.charAt(0) === "%") { try { icon = decodeURIComponent(icon); } catch (e) {} } // saved by an older version
      a.children[0].textContent = icon;
      a.children[1].children[0].textContent = r.title;
      a.children[1].children[1].textContent = r.chapter;
      list.appendChild(a);
      var card = document.querySelector('[data-guide="' + r.guide + '"] .resume-badge');
      if (card) card.classList.remove("hidden");
    });
    document.getElementById("continue").classList.remove("hidden");
  }

  // -------------------------------------------------------------- search --
  var input = document.getElementById("search");
  var results = document.getElementById("results");
  var index = null;
  function ensureIndex() {
    if (index) return Promise.resolve(index);
    return fetch("search.json").then(function (r) { return r.json(); }).then(function (j) { index = j; return j; });
  }
  function render(q) {
    var terms = q.toLowerCase().split(/\s+/).filter(Boolean);
    if (!terms.length) { results.classList.add("hidden"); return; }
    var hits = index.filter(function (e) {
      var hay = (e.t + " " + e.g).toLowerCase();
      return terms.every(function (t) { return hay.indexOf(t) >= 0; });
    }).slice(0, 40);
    results.innerHTML = "";
    if (!hits.length) {
      results.innerHTML = '<p class="px-3 py-2 text-sm text-muted">No matches.</p>';
    }
    hits.forEach(function (h) {
      var a = document.createElement("a");
      a.href = h.u;
      a.className = "flex items-start gap-3 rounded-lg px-3 py-2 no-underline hover:bg-accent-soft";
      a.innerHTML = '<span class="text-lg"></span><span class="min-w-0"><span class="block text-sm text-ink"></span><span class="block truncate text-xs text-faint"></span></span>';
      a.children[0].textContent = h.i;
      a.children[1].children[0].textContent = h.t;
      a.children[1].children[1].textContent = h.g;
      results.appendChild(a);
    });
    results.classList.remove("hidden");
  }
  input.addEventListener("input", function () { ensureIndex().then(function () { render(input.value); }); });
  input.addEventListener("focus", function () { ensureIndex(); if (input.value) render(input.value); });
  document.addEventListener("click", function (e) { if (!e.target.closest("#search, #results")) results.classList.add("hidden"); });
  addEventListener("keydown", function (e) {
    if (e.target === input) { if (e.key === "Escape") { input.value = ""; results.classList.add("hidden"); input.blur(); } return; }
    if (e.key === "/") { e.preventDefault(); input.focus(); }
    if (e.key === "d") cycle();
  });
})();
