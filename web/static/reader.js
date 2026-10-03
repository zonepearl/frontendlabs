// reader.js: the book reader. No dependencies; every feature degrades
// gracefully if localStorage is unavailable (private windows, blocked storage).
(function () {
  "use strict";
  var doc = document.documentElement;
  var body = document.body;
  var page = body.dataset.page;

  // ------------------------------------------------------------ storage --
  function load(key, fallback) {
    try { var v = localStorage.getItem(key); return v ? JSON.parse(v) : fallback; } catch (e) { return fallback; }
  }
  function save(key, value) {
    try { localStorage.setItem(key, JSON.stringify(value)); } catch (e) {}
  }

  // -------------------------------------------------------- preferences --
  var prefs = load("fl:prefs", {});
  var THEMES = ["paper", "light", "night"];
  var WIDTHS = ["narrow", "normal", "wide"];
  var themeBtn = document.getElementById("themeBtn");
  function applyPrefs() {
    doc.dataset.theme = new URLSearchParams(location.search).get("theme") || prefs.theme || doc.dataset.theme || "paper";
    doc.style.setProperty("--font-scale", prefs.font || 1);
    doc.dataset.width = prefs.width || "normal";
    if (themeBtn) themeBtn.textContent = doc.dataset.theme[0].toUpperCase() + doc.dataset.theme.slice(1);
  }
  function setPref(k, v) { prefs[k] = v; save("fl:prefs", prefs); applyPrefs(); }
  function cycleTheme() { setPref("theme", THEMES[(THEMES.indexOf(doc.dataset.theme) + 1) % THEMES.length]); }
  function font(delta) { setPref("font", Math.min(1.4, Math.max(0.8, Math.round(((prefs.font || 1) + delta) * 100) / 100))); }
  function cycleWidth() { setPref("width", WIDTHS[(WIDTHS.indexOf(doc.dataset.width || "normal") + 1) % WIDTHS.length]); }
  applyPrefs();
  themeBtn && themeBtn.addEventListener("click", cycleTheme);
  document.getElementById("fsUp").addEventListener("click", function () { font(0.05); });
  document.getElementById("fsDown").addEventListener("click", function () { font(-0.05); });
  document.getElementById("widthBtn").addEventListener("click", cycleWidth);

  // ------------------------------------------------------- contents (TOC) --
  var toc = document.getElementById("toc");
  var scrim = document.getElementById("scrim");
  var desktop = matchMedia("(min-width: 1024px)");
  function tocOpen() { return desktop.matches ? !body.classList.contains("toc-hidden") : body.classList.contains("toc-open"); }
  function toggleToc(force) {
    var open = force === undefined ? !tocOpen() : force;
    if (desktop.matches) { body.classList.toggle("toc-hidden", !open); save("fl:tocHidden", !open); }
    else { body.classList.toggle("toc-open", open); scrim.classList.toggle("hidden", !open); }
  }
  if (load("fl:tocHidden", false)) body.classList.add("toc-hidden");
  document.getElementById("tocBtn").addEventListener("click", function () { toggleToc(); });
  scrim.addEventListener("click", function () { toggleToc(false); });
  toc.addEventListener("click", function (e) { if (e.target.closest("a") && !desktop.matches) toggleToc(false); });

  var filter = document.getElementById("tocFilter");
  filter.addEventListener("input", function () {
    var q = filter.value.trim().toLowerCase();
    toc.querySelectorAll("#tocList li").forEach(function (li) {
      li.style.display = !q || li.textContent.toLowerCase().indexOf(q) >= 0 ? "" : "none";
    });
  });

  // ---------------------------------------- chapters: spy, progress, resume --
  var chapters = Array.prototype.slice.call(document.querySelectorAll("section.chapter"));
  var links = {};
  toc.querySelectorAll("a[data-id]").forEach(function (a) { links[a.dataset.id] = a; });
  var where = document.getElementById("where");
  var progress = document.getElementById("progress");
  var current = null;

  function setCurrent(ch) {
    if (ch === current) return;
    if (current && links[current.dataset.ch]) links[current.dataset.ch].classList.remove("is-current");
    current = ch;
    if (!ch) return;
    var a = links[ch.dataset.ch];
    if (a) {
      a.classList.add("is-current");
      var r = a.getBoundingClientRect(), tr = toc.getBoundingClientRect();
      if (r.top < tr.top + 60 || r.bottom > tr.bottom - 20) a.scrollIntoView({ block: "center" });
    }
    var h = ch.querySelector("h2");
    if (where && h) where.textContent = h.textContent;
  }

  var ticking = false;
  function onScroll() {
    if (ticking) return;
    ticking = true;
    requestAnimationFrame(function () {
      ticking = false;
      var max = doc.scrollHeight - innerHeight;
      progress.style.width = (max > 0 ? (scrollY / max) * 100 : 0) + "%";
      var ch = null;
      for (var i = 0; i < chapters.length; i++) {
        if (chapters[i].getBoundingClientRect().top <= 120) ch = chapters[i]; else break;
      }
      setCurrent(ch);
      remember();
    });
  }
  addEventListener("scroll", onScroll, { passive: true });

  // resume: the chapter and position within it, per page; plus "last read" for the home page
  var resumeKey = "fl:pos:" + page;
  var saveTimer = null;
  function remember() {
    clearTimeout(saveTimer);
    saveTimer = setTimeout(function () {
      if (!current || scrollY < 400) return;
      var h = current.querySelector("h2");
      var r = current.getBoundingClientRect();
      var pos = { ch: current.dataset.ch, title: h ? h.textContent : "", frac: Math.max(0, -r.top / Math.max(1, r.height)), t: Date.now() };
      save(resumeKey, pos);
      var recent = load("fl:recent", []).filter(function (x) { return x.page !== page; });
      recent.unshift({ page: page, guide: body.dataset.guide, title: body.dataset.guideTitle, icon: body.dataset.emoji, chapter: pos.title, ch: pos.ch, t: pos.t });
      save("fl:recent", recent.slice(0, 6));
    }, 400);
  }

  var resume = load(resumeKey, null);
  var resumeBox = document.getElementById("resume");
  if (resume && !location.hash && document.getElementById("ch-" + resume.ch)) {
    var btn = document.createElement("button");
    btn.className = "resume-btn";
    btn.innerHTML = "Continue where you left off <span></span>";
    btn.querySelector("span").textContent = "· " + resume.title;
    btn.addEventListener("click", function () {
      var ch = document.getElementById("ch-" + resume.ch);
      scrollTo(0, ch.offsetTop + ch.offsetHeight * resume.frac - 70);
    });
    resumeBox.appendChild(btn);
  }

  // previous / next chapter at the end of each chapter
  chapters.forEach(function (ch, i) {
    var nav = document.createElement("nav");
    nav.className = "chapnav";
    function link(target, label, cls) {
      var a = document.createElement("a");
      a.href = "#" + target.dataset.ch;
      a.className = cls;
      a.innerHTML = "<small></small><span></span>";
      a.querySelector("small").textContent = label;
      a.querySelector("span").textContent = target.querySelector("h2").textContent;
      return a;
    }
    nav.appendChild(i > 0 ? link(chapters[i - 1], "← Previous", "prev") : document.createElement("span"));
    if (i + 1 < chapters.length) nav.appendChild(link(chapters[i + 1], "Next →", "next"));
    ch.appendChild(nav);
  });

  // ------------------------------------------------ code blocks: copy --
  document.querySelectorAll("#content pre").forEach(function (pre) {
    var wrap = document.createElement("div");
    wrap.className = "codewrap";
    pre.parentNode.insertBefore(wrap, pre);
    wrap.appendChild(pre);
    var lang = (pre.querySelector("code") || {}).className || "";
    var m = /language-([\w-]+)/.exec(lang);
    if (m) wrap.dataset.lang = m[1];
    var b = document.createElement("button");
    b.className = "copy";
    b.textContent = "Copy";
    b.addEventListener("click", function () {
      navigator.clipboard.writeText(pre.innerText).then(function () {
        b.textContent = "Copied"; b.classList.add("done");
        setTimeout(function () { b.textContent = "Copy"; b.classList.remove("done"); }, 1200);
      });
    });
    wrap.appendChild(b);
  });

  // ------------------------------------------ checklists that remember --
  var checks = load("fl:checks:" + page, {});
  document.querySelectorAll('#content input[type="checkbox"]').forEach(function (cb, i) {
    cb.disabled = false;
    if (checks[i] !== undefined) cb.checked = checks[i];
    cb.addEventListener("change", function () { checks[i] = cb.checked; save("fl:checks:" + page, checks); });
  });

  // ---------------------------------------------------------- keyboard --
  addEventListener("keydown", function (e) {
    if (e.metaKey || e.ctrlKey || e.altKey || /input|textarea|select/i.test(e.target.tagName)) return;
    var i = chapters.indexOf(current);
    switch (e.key) {
      case "t": toggleToc(); break;
      case "d": cycleTheme(); break;
      case "w": cycleWidth(); break;
      case "+": case "=": font(0.05); break;
      case "-": font(-0.05); break;
      case "n": if (chapters[i + 1]) chapters[i + 1].scrollIntoView(); break;
      case "p": if (i > 0) chapters[i - 1].scrollIntoView(); else if (current) current.scrollIntoView(); break;
      case "/": e.preventDefault(); toggleToc(true); filter.focus(); break;
      case "Escape": if (!desktop.matches) toggleToc(false); break;
      default: return;
    }
  });

  // print: open collapsed answers so they appear in the PDF
  var closed = [];
  addEventListener("beforeprint", function () {
    closed = Array.prototype.slice.call(document.querySelectorAll("#content details:not([open])"));
    closed.forEach(function (d) { d.open = true; });
  });
  addEventListener("afterprint", function () { closed.forEach(function (d) { d.open = false; }); });

  onScroll();
})();
