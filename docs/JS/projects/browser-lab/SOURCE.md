# Browser Lab — full source

> The complete, tested source of Browser Lab, generated from the files in
> `JS/projects/browser-lab/` by `build-source-pages.mjs`. 26 files.

A live, animated tour of how the browser and the JavaScript engine work,
built with plain JavaScript and CSS (no framework, no build step, no runtime
dependencies). It is Capstone II of the
[JavaScript field guide](../../real-life-JS-guide.md) (§79–§84).

Nine stations, each with missions that tick themselves off:

| | Station | What's live |
|---|---|---|
| 00 | The big picture | process/thread map, an animated "life of a click", your browser's capabilities |
| 01 | URL to bytes | Navigation Timing, an HTTP/1.1 vs HTTP/2 request burst, TLS 1.3 vs 1.2, ALPN/cipher, SSE |
| 02 | Rendering | pipeline stages per CSS property, layout thrashing measured, DOM vs render tree |
| 03 | Event loop | step-through simulator, prediction quiz, verified against the real engine |
| 04 | Inside V8 | tier-up/deopt, real bytecode, hidden classes, inline caches, megamorphic benchmark |
| 05 | Threads | visible jank, Web Worker fix, copy vs transfer, SharedArrayBuffer race vs Atomics |
| 06 | Storage | IndexedDB events as they fire, the auto-commit pitfall, sync vs async storage |
| 07 | Core Web Vitals | live LCP/INP/CLS gauges, break and fix CLS and INP, Long Animation Frames |
| 08 | Debugging | captured errors/rejections/CSP violations, a findable leak, User Timing, DevTools recipes |

## Run it

```bash
npm run cert    # once: self-signed localhost cert via openssl (or use mkcert)
npm start       # http://localhost:8088  (HTTP/1.1)   https://localhost:8443  (HTTP/2 over TLS 1.3)
```

Open both URLs and run the Network station's request burst in each.

## Test it

```bash
npm install
npx playwright install chromium
npm test        # 46 tests, over HTTP/1.1 and HTTP/2
```

## Files

- `package.json`
- `server.mjs`
- `scripts/cert.mjs`
- `public/index.html`
- `public/css/lab.css`
- `public/js/lib/dom.js`
- `public/js/lib/errors.js`
- `public/js/lib/loop-model.js`
- `public/js/lib/missions.js`
- `public/js/lib/primes.js`
- `public/js/lib/vitals.js`
- `public/js/main.js`
- `public/js/stations/debug.js`
- `public/js/stations/engine.js`
- `public/js/stations/eventloop.js`
- `public/js/stations/network.js`
- `public/js/stations/overview.js`
- `public/js/stations/render.js`
- `public/js/stations/storage.js`
- `public/js/stations/threads.js`
- `public/js/stations/vitals.js`
- `public/js/workers/counter.js`
- `public/js/workers/echo.js`
- `public/js/workers/primes.js`
- `playwright.config.js`
- `tests/lab.spec.js`

## package.json

```json
{
  "name": "browser-lab",
  "private": true,
  "type": "module",
  "description": "A live, animated tour of how the browser and the JavaScript engine work — the capstone of the JavaScript field guide (§79–§84). Plain JS and CSS, no build step.",
  "scripts": {
    "cert": "node scripts/cert.mjs",
    "start": "node server.mjs",
    "test": "playwright test"
  },
  "devDependencies": {
    "@playwright/test": "^1.63.0"
  }
}
```

## server.mjs

```javascript
// Browser Lab server — zero dependencies (§79).
//
//   npm run cert   one-time: a self-signed cert for localhost (or use mkcert)
//   npm start      http://localhost:8088   HTTP/1.1, plain text
//                  https://localhost:8443  HTTP/2 over TLS 1.3 (when certs/ exists)
//
// Open both and run the Network station's "fire 18 requests" experiment in
// each: HTTP/1.1 shows a 6-connection staircase, HTTP/2 one multiplexed burst.

import http from "node:http";
import http2 from "node:http2";
import { readFile } from "node:fs/promises";
import { existsSync, readFileSync } from "node:fs";
import { extname, join, normalize, sep } from "node:path";
import { fileURLToPath } from "node:url";

const HTTP_PORT = Number(process.env.HTTP_PORT ?? 8088);
const HTTPS_PORT = Number(process.env.HTTPS_PORT ?? 8443);
const HOST = process.env.HOST ?? "localhost";
const ROOT = fileURLToPath(new URL(".", import.meta.url));
const PUBLIC_DIR = join(ROOT, "public") + sep;
const CERT = join(ROOT, "certs", "cert.pem");
const KEY = join(ROOT, "certs", "key.pem");

const MIME = {
  ".html": "text/html; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".svg": "image/svg+xml",
  ".json": "application/json",
};

// Cross-origin isolation (COOP + COEP) unlocks SharedArrayBuffer and
// performance.measureUserAgentSpecificMemory() for the Threads and Debugging
// stations (§39, §83). The CSP is strict: no inline script, no eval — the
// Debugging station shows what a violation looks like.
const HEADERS = {
  "Content-Security-Policy":
    "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data: blob:; " +
    "worker-src 'self' blob:; connect-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'",
  "Cross-Origin-Opener-Policy": "same-origin",
  "Cross-Origin-Embedder-Policy": "require-corp",
  "Cross-Origin-Resource-Policy": "same-origin",
  "X-Content-Type-Options": "nosniff",
  "Referrer-Policy": "no-referrer",
};

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

function send(res, status, body, headers = {}) {
  res.writeHead(status, { ...HEADERS, ...headers });
  res.end(body);
}

const json = (res, status, data, headers = {}) =>
  send(res, status, JSON.stringify(data), { "Content-Type": "application/json", "Cache-Control": "no-store", ...headers });

// What the browser actually negotiated and sent, from the server's side.
function connectionInfo(req) {
  const socket = req.stream?.session?.socket ?? req.socket;
  return {
    httpVersion: req.httpVersion,
    alpn: socket.alpnProtocol || null,
    tls: socket.getProtocol?.() ?? null,
    cipher: socket.getCipher?.()?.name ?? null,
    remotePort: socket.remotePort,
    headers: req.headers,
  };
}

const routes = {
  // A deliberately slow endpoint for the waterfall experiment.
  "/api/delay": async (req, res, url) => {
    const ms = Math.min(Number(url.searchParams.get("ms")) || 300, 5000);
    const start = performance.now();
    await sleep(ms);
    // Server-Timing shows up in DevTools and in PerformanceResourceTiming.serverTiming.
    json(res, 200, { ms, id: url.searchParams.get("id") }, {
      "Server-Timing": `wait;dur=${(performance.now() - start).toFixed(1)};desc="simulated work"`,
    });
  },
  "/api/whoami": (req, res) => json(res, 200, connectionInfo(req)),
  // A server-sent event stream: one tick every 500ms, ten ticks.
  "/api/ticks": async (req, res) => {
    res.writeHead(200, { ...HEADERS, "Content-Type": "text/event-stream", "Cache-Control": "no-store" });
    for (let i = 1; i <= 10 && !res.destroyed; i++) {
      res.write(`id: ${i}\ndata: ${JSON.stringify({ tick: i, at: Date.now() })}\n\n`);
      await sleep(500);
    }
    res.end();
  },
};

async function handler(req, res) {
  const url = new URL(req.url, "http://x");
  try {
    if (routes[url.pathname]) return await routes[url.pathname](req, res, url);
    const rel = url.pathname === "/" ? "index.html" : decodeURIComponent(url.pathname.slice(1));
    const file = normalize(join(PUBLIC_DIR, rel));
    if (!file.startsWith(PUBLIC_DIR)) return json(res, 404, { error: "Not found" });
    const body = await readFile(file);
    send(res, 200, body, { "Content-Type": MIME[extname(file)] ?? "application/octet-stream", "Cache-Control": "no-cache" });
  } catch (err) {
    if (err.code === "ENOENT" || err.code === "EISDIR") return json(res, 404, { error: "Not found" });
    console.error(err);
    if (!res.headersSent) json(res, 500, { error: "Something went wrong" });
  }
}

http.createServer(handler).listen(HTTP_PORT, HOST, () =>
  console.log(`Browser Lab  HTTP/1.1  http://${HOST}:${HTTP_PORT}`),
);

if (existsSync(CERT) && existsSync(KEY)) {
  // allowHTTP1: browsers that don't offer "h2" in ALPN still get a page.
  http2
    .createSecureServer({ cert: readFileSync(CERT), key: readFileSync(KEY), allowHTTP1: true }, handler)
    .listen(HTTPS_PORT, HOST, () => console.log(`Browser Lab  HTTP/2    https://${HOST}:${HTTPS_PORT}`));
} else {
  console.log("No certs/ found: run `npm run cert` to also serve HTTP/2 over TLS on :" + HTTPS_PORT);
}
```

## scripts/cert.mjs

```javascript
// Creates certs/cert.pem + certs/key.pem: a self-signed RSA-2048 certificate
// for localhost. Browsers will warn once; mkcert (https://github.com/FiloSottile/mkcert)
// installs a local CA instead, so there is no warning at all.
import { execFileSync } from "node:child_process";
import { existsSync, mkdirSync } from "node:fs";

mkdirSync("certs", { recursive: true });
if (existsSync("certs/cert.pem") && existsSync("certs/key.pem")) {
  console.log("certs/ already exists");
} else {
  execFileSync("openssl", [
    // RSA, not EC: macOS's LibreSSL writes EC keys that some TLS stacks reject.
    "req", "-x509", "-newkey", "rsa:2048",
    "-nodes", "-days", "365", "-subj", "/CN=localhost",
    "-addext", "subjectAltName=DNS:localhost,IP:127.0.0.1",
    "-keyout", "certs/key.pem", "-out", "certs/cert.pem",
  ], { stdio: "inherit" });
  console.log("wrote certs/cert.pem and certs/key.pem");
}
```

## public/index.html

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Browser Lab</title>
  <meta name="description" content="A live, animated tour of how the browser and JavaScript engine work.">
  <link rel="stylesheet" href="/css/lab.css">
  <!-- modulepreload: fetch the router's imports in parallel with it (§26) -->
  <link rel="modulepreload" href="/js/lib/dom.js">
  <link rel="modulepreload" href="/js/lib/missions.js">
  <script type="module" src="/js/main.js"></script>
</head>
<body>
  <a class="skip" href="#stage">Skip to the station</a>
  <div class="shell">
    <nav class="sidebar" aria-label="Stations">
      <div class="brand">
        <span class="logo" aria-hidden="true">⚡</span>
        <div>
          <strong>Browser Lab</strong>
          <small id="proto">loading…</small>
        </div>
      </div>
      <ol id="station-nav"></ol>
      <div class="progress">
        <div class="bar"><span id="progress-fill"></span></div>
        <small id="progress-text"></small>
        <button id="reset-missions" type="button" class="link">reset missions</button>
      </div>
    </nav>
    <main id="stage" tabindex="-1">
      <p class="lede">Loading…</p>
    </main>
  </div>
</body>
</html>
```

## public/css/lab.css

```css
/* Browser Lab — plain CSS with custom properties. No framework, no build. */
:root {
  --bg: #f4f5f8; --surface: #ffffff; --surface-2: #eef0f5; --ink: #161a23; --muted: #5d6677; --line: #d9dde6;
  --accent: #3355ff; --accent-soft: #e3e8ff; --good: #0a8a4c; --warn: #b26b00; --bad: #c62828;
  --js: #c79500; --style: #8e44ad; --layout: #1f7a8c; --paint: #d35400; --composite: #2e86de;
  --task: #3355ff; --micro: #0a8a4c; --webapi: #b26b00;
  --radius: 12px; --mono: ui-monospace, "SF Mono", Menlo, Consolas, monospace;
  color-scheme: light dark;
  font: 15px/1.55 system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
}
@media (prefers-color-scheme: dark) {
  :root { --bg: #0e1117; --surface: #161b24; --surface-2: #1d2430; --ink: #e6e9ef; --muted: #98a2b3; --line: #2a3240;
    --accent: #7d95ff; --accent-soft: #232d52; --good: #4ade80; --warn: #fbbf24; --bad: #f87171;
    --task: #7d95ff; --micro: #4ade80; --webapi: #fbbf24; }
}
* { box-sizing: border-box; }
body { margin: 0; background: var(--bg); color: var(--ink); }
button { font: inherit; cursor: pointer; border: 1px solid var(--line); background: var(--surface); color: var(--ink);
  padding: 7px 13px; border-radius: 9px; transition: transform .1s, background .2s; }
button:hover { background: var(--surface-2); }
button:active { transform: scale(.97); }
button.primary { background: var(--accent); color: #fff; border-color: transparent; }
button.danger { color: var(--bad); }
button.link { border: 0; background: none; color: var(--accent); padding: 0; text-decoration: underline; }
button:disabled { opacity: .5; cursor: not-allowed; }
input, select { font: inherit; padding: 6px 9px; border: 1px solid var(--line); border-radius: 8px; background: var(--surface); color: var(--ink); }
code, pre, .mono { font-family: var(--mono); font-size: .86rem; }
.skip { position: absolute; left: -999px; }
.skip:focus { left: 12px; top: 12px; z-index: 10; background: var(--surface); padding: 8px; }

/* ---------------------------------------------------------- layout ---- */
.shell { display: grid; grid-template-columns: 290px 1fr; min-height: 100vh; }
.sidebar { position: sticky; top: 0; height: 100vh; overflow-y: auto; padding: 18px 14px; border-right: 1px solid var(--line);
  background: var(--surface); display: flex; flex-direction: column; gap: 14px; }
.brand { display: flex; gap: 10px; align-items: center; }
.brand small { display: block; color: var(--muted); font-family: var(--mono); font-size: .75rem; }
.logo { font-size: 1.6rem; animation: pulse 2.4s ease-in-out infinite; }
@keyframes pulse { 50% { transform: scale(1.15) rotate(-8deg); } }
#station-nav { list-style: none; margin: 0; padding: 0; display: grid; gap: 2px; }
#station-nav a { display: grid; grid-template-columns: 22px 22px 1fr auto; gap: 8px; align-items: center; padding: 7px 8px;
  border-radius: 8px; color: inherit; text-decoration: none; }
#station-nav a:hover { background: var(--surface-2); }
#station-nav a[aria-current="page"] { background: var(--accent-soft); font-weight: 600; }
#station-nav .num { color: var(--muted); font-family: var(--mono); font-size: .75rem; }
.badge { font-size: .72rem; font-family: var(--mono); color: var(--muted); background: var(--surface-2); border-radius: 99px; padding: 1px 7px; }
.badge.complete { background: var(--good); color: #fff; }
.progress { margin-top: auto; display: grid; gap: 4px; }
.progress .bar { height: 6px; background: var(--surface-2); border-radius: 99px; overflow: hidden; }
.progress .bar span { display: block; height: 100%; width: 0; background: linear-gradient(90deg, var(--accent), var(--good)); transition: width .6s; }
#stage { padding: 26px clamp(16px, 3vw, 40px) 60px; max-width: 1200px; outline: none; }
@media (max-width: 860px) {
  .shell { grid-template-columns: 1fr; }
  .sidebar { position: static; height: auto; }
}

/* --------------------------------------------------------- station ---- */
.station > header h2 { margin: 0 0 4px; font-size: 1.6rem; }
.lede { color: var(--muted); margin: 0 0 18px; max-width: 75ch; }
.panel { background: var(--surface); border: 1px solid var(--line); border-radius: var(--radius); padding: 16px 18px; margin: 14px 0; }
.panel h3 { margin: 0 0 10px; font-size: 1.02rem; }
.panel p { max-width: 80ch; }
.row { display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
.grid-2 { display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 14px; }
.note { color: var(--muted); font-size: .88rem; }
.missions { list-style: none; padding: 10px 14px; margin: 0 0 6px; display: flex; gap: 6px 18px; flex-wrap: wrap;
  background: var(--surface); border: 1px dashed var(--line); border-radius: var(--radius); font-size: .9rem; }
.missions li::before { content: "◻ "; color: var(--muted); }
.missions li.done { color: var(--good); }
.missions li.done::before { content: "✅ "; }
#toasts { position: fixed; right: 18px; bottom: 18px; display: flex; flex-direction: column; gap: 8px; z-index: 20; }
.toast { background: var(--ink); color: var(--bg); padding: 10px 16px;
  border-radius: 10px; box-shadow: 0 8px 30px #0004; animation: toast-in .35s ease-out; }
@keyframes toast-in { from { transform: translateY(20px); opacity: 0; } }
.readouts { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 10px; }
.readout { background: var(--surface-2); border-radius: 10px; padding: 10px 12px; }
.readout b { display: block; font-size: 1.25rem; font-family: var(--mono); }
.readout small { color: var(--muted); }
.kv { display: grid; grid-template-columns: max-content 1fr; gap: 4px 14px; font-size: .9rem; margin: 0; }
.kv dt { color: var(--muted); }
.kv dd { margin: 0; font-family: var(--mono); overflow-wrap: anywhere; }
.log { font-family: var(--mono); font-size: .82rem; background: var(--surface-2); border-radius: 8px; padding: 8px 10px;
  max-height: 260px; overflow: auto; margin: 8px 0 0; list-style: none; }
.log li { padding: 2px 0; border-bottom: 1px solid var(--line); animation: pop .3s ease-out; }
.log li:last-child { border-bottom: 0; }
@keyframes pop { from { opacity: 0; transform: translateX(-6px); } }

/* --------------------------------------------------------- overview --- */
.arch { display: grid; grid-template-columns: repeat(auto-fit, minmax(210px, 1fr)); gap: 12px; }
.proc { border: 2px solid var(--line); border-radius: var(--radius); padding: 10px 12px; background: var(--surface); position: relative; transition: border-color .3s, box-shadow .3s; }
.proc h4 { margin: 0 0 6px; font-size: .95rem; }
.proc ul { margin: 0; padding-left: 18px; font-size: .86rem; color: var(--muted); }
.proc.lit { border-color: var(--accent); box-shadow: 0 0 0 4px var(--accent-soft); }
.proc .thread { border: 1px solid var(--line); border-radius: 8px; padding: 4px 8px; margin: 6px 0 0; font-size: .85rem; transition: background .3s; }
.proc .thread.lit { background: var(--accent-soft); }
.tour-caption { min-height: 3em; font-size: 1.02rem; margin-top: 12px; }

/* --------------------------------------------------------- pipeline --- */
.pipeline { display: flex; gap: 6px; flex-wrap: wrap; margin: 10px 0; }
.stage-pill { padding: 6px 14px; border-radius: 99px; border: 2px solid var(--line); font-weight: 600; font-size: .88rem; opacity: .35; transition: opacity .25s, transform .25s; }
.stage-pill.on { opacity: 1; transform: translateY(-2px); }
.stage-pill.skipped { text-decoration: line-through; }
.stage-pill[data-stage="js"] { border-color: var(--js); }
.stage-pill[data-stage="style"] { border-color: var(--style); }
.stage-pill[data-stage="layout"] { border-color: var(--layout); }
.stage-pill[data-stage="paint"] { border-color: var(--paint); }
.stage-pill[data-stage="composite"] { border-color: var(--composite); }
.stage-pill.flash { animation: flash .5s; }
@keyframes flash { 30% { box-shadow: 0 0 0 6px var(--accent-soft); } }
.track { position: relative; height: 70px; background: var(--surface-2); border-radius: 10px; overflow: hidden; }
.mover { position: absolute; left: 8px; top: 11px; width: 48px; height: 48px; border-radius: 10px; background: var(--accent); }
.bars { display: grid; gap: 1px; max-height: 140px; overflow: hidden; }
.bars div { height: 2px; background: var(--accent); width: 100px; }
.tree { font-family: var(--mono); font-size: .84rem; margin: 0; padding-left: 16px; }
.tree .gone { color: var(--muted); text-decoration: line-through; }
.tree .pseudo { color: var(--paint); }

/* -------------------------------------------------------- waterfall --- */
.waterfall { display: grid; gap: 3px; font-family: var(--mono); font-size: .76rem; }
.wf-row { display: grid; grid-template-columns: 70px 1fr 60px; gap: 8px; align-items: center; }
.wf-track { position: relative; height: 14px; background: var(--surface-2); border-radius: 4px; }
.wf-bar { position: absolute; top: 2px; height: 10px; display: flex; border-radius: 3px; overflow: hidden; transform-origin: left; animation: grow .6s ease-out both; }
.wf-bar span { height: 100%; }
.seg-wait { background: #9aa3b2; } .seg-ttfb { background: var(--good); } .seg-dl { background: var(--accent); }
@keyframes grow { from { transform: scaleX(0); } }
.legend { display: flex; gap: 14px; font-size: .8rem; color: var(--muted); flex-wrap: wrap; margin: 6px 0; }
.legend i { display: inline-block; width: 12px; height: 10px; border-radius: 2px; margin-right: 4px; vertical-align: -1px; }
.phases { display: grid; gap: 6px; }
.phase { display: grid; grid-template-columns: 150px 1fr 70px; gap: 10px; align-items: center; font-size: .88rem; }
.phase .fill { height: 12px; border-radius: 6px; background: var(--accent); min-width: 3px; transform-origin: left; animation: grow .7s ease-out both; }

/* ----------------------------------------------------- tls sequence --- */
.seq { position: relative; display: grid; grid-template-columns: 1fr 1fr; gap: 0 20px; margin-top: 8px; }
.seq .actor { text-align: center; font-weight: 700; padding: 6px; border-bottom: 2px solid var(--line); }
.seq .msgs { grid-column: 1 / -1; display: grid; gap: 6px; padding-top: 8px; }
.msg-line { position: relative; padding: 6px 10px; border-radius: 8px; font-size: .86rem; width: 70%; animation: fly .5s ease-out both; }
.msg-line.c2s { justify-self: start; background: var(--accent-soft); border-left: 4px solid var(--accent); }
.msg-line.s2c { justify-self: end; background: var(--surface-2); border-right: 4px solid var(--good); text-align: right; }
.msg-line.c2s::after { content: "→"; position: absolute; right: -18px; top: 6px; color: var(--accent); }
.msg-line.s2c::before { content: "←"; position: absolute; left: -18px; top: 6px; color: var(--good); }
.msg-line small { display: block; color: var(--muted); }
.msg-line .lock { color: var(--good); }
@keyframes fly { from { opacity: 0; transform: translateX(var(--from, -30px)); } }
.msg-line.s2c { --from: 30px; }
.rtt { grid-column: 1 / -1; text-align: center; font-size: .78rem; color: var(--muted); border-top: 1px dashed var(--line); margin-top: 2px; }

/* ------------------------------------------------------- event loop --- */
.loop-grid { display: grid; grid-template-columns: minmax(280px, 1.1fr) 1fr; gap: 14px; }
@media (max-width: 1000px) { .loop-grid { grid-template-columns: 1fr; } }
.code { background: var(--surface-2); border-radius: 10px; padding: 10px 0; margin: 0; overflow-x: auto; counter-reset: ln; }
.code span { display: block; padding: 0 12px; white-space: pre; transition: background .2s; }
.code span::before { counter-increment: ln; content: counter(ln); display: inline-block; width: 2em; color: var(--muted); }
.code span.hl { background: var(--accent-soft); box-shadow: inset 3px 0 var(--accent); }
.boxes { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.box { border: 1px solid var(--line); border-radius: 10px; padding: 8px; min-height: 92px; background: var(--surface); }
.box h4 { margin: 0 0 6px; font-size: .78rem; text-transform: uppercase; letter-spacing: .05em; color: var(--muted); }
.box.stack .items { display: flex; flex-direction: column-reverse; } /* newest frame on top */
.token { font-family: var(--mono); font-size: .8rem; padding: 4px 8px; border-radius: 6px; margin: 3px 0; animation: drop .35s cubic-bezier(.2,1.4,.5,1) both; overflow-wrap: anywhere; }
.box.stack .token { background: var(--accent-soft); border: 1px solid var(--accent); }
.box.webapis .token { background: color-mix(in srgb, var(--webapi) 18%, transparent); border: 1px solid var(--webapi); }
.box.tasks .token { background: color-mix(in srgb, var(--task) 15%, transparent); border: 1px solid var(--task); }
.box.micro .token { background: color-mix(in srgb, var(--micro) 15%, transparent); border: 1px solid var(--micro); }
.box.console .token { background: var(--surface-2); }
@keyframes drop { from { opacity: 0; transform: translateY(-10px) scale(.9); } }
.loop-note { min-height: 3.2em; padding: 10px 12px; border-radius: 10px; background: var(--accent-soft); margin: 10px 0; }
.wheel { width: 54px; height: 54px; border-radius: 50%; border: 5px solid var(--line); border-top-color: var(--accent); flex: none; }
.wheel.spin { animation: spin .8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.quiz { display: flex; gap: 6px; flex-wrap: wrap; }
.quiz button.picked { opacity: .35; }
.answer { font-family: var(--mono); font-size: .85rem; }
.ok { color: var(--good); font-weight: 600; } .bad { color: var(--bad); font-weight: 600; }

/* ----------------------------------------------------------- engine --- */
.tiers { display: flex; gap: 6px; flex-wrap: wrap; align-items: center; }
.tier { padding: 8px 12px; border-radius: 10px; border: 2px solid var(--line); text-align: center; min-width: 110px; transition: all .3s; }
.tier small { display: block; color: var(--muted); font-size: .75rem; }
.tier.on { border-color: var(--accent); background: var(--accent-soft); transform: scale(1.06); }
.arrow { color: var(--muted); }
.heat { height: 10px; border-radius: 99px; background: var(--surface-2); overflow: hidden; }
.heat span { display: block; height: 100%; width: 0; background: linear-gradient(90deg, #3aa0ff, #ffb000, #ff3d00); transition: width .3s; }
.shapes ul { list-style: none; margin: 0; padding-left: 20px; border-left: 2px solid var(--line); }
.shapes > ul { border-left: 0; padding-left: 0; }
.shape { display: inline-flex; gap: 6px; align-items: center; font-family: var(--mono); font-size: .82rem; padding: 3px 8px; margin: 3px 0; border-radius: 6px; background: var(--surface-2); animation: drop .3s both; }
.shape .count { background: var(--accent); color: #fff; border-radius: 99px; padding: 0 6px; font-size: .72rem; }
.ic { font-weight: 700; }
.ic.monomorphic { color: var(--good); } .ic.polymorphic { color: var(--warn); } .ic.megamorphic { color: var(--bad); }

/* ---------------------------------------------------------- threads --- */
canvas { width: 100%; height: auto; display: block; background: var(--surface-2); border-radius: 10px; }
.lanes { display: grid; gap: 6px; }
.lane { display: grid; grid-template-columns: 110px 1fr; gap: 8px; align-items: center; font-size: .85rem; }
.lane .rail { height: 26px; border-radius: 6px; background: var(--surface-2); position: relative; overflow: hidden; }
.lane .work { position: absolute; top: 3px; bottom: 3px; border-radius: 4px; background: var(--paint); animation: grow .3s both; }
.lane .work.worker { background: var(--good); }

/* ------------------------------------------------------------ vitals --- */
.gauges { display: grid; grid-template-columns: repeat(auto-fit, minmax(170px, 1fr)); gap: 12px; }
.gauge { border-radius: var(--radius); padding: 12px 14px; background: var(--surface); border: 2px solid var(--line); transition: border-color .4s; }
.gauge b { font-size: 1.7rem; font-family: var(--mono); display: block; }
.gauge .meter { height: 8px; border-radius: 99px; background: linear-gradient(90deg, var(--good) 0 33%, var(--warn) 33% 66%, var(--bad) 66%); position: relative; margin-top: 8px; }
.gauge .needle { position: absolute; top: -4px; width: 4px; height: 16px; background: var(--ink); border-radius: 2px; transition: left .5s cubic-bezier(.3,1.5,.5,1); }
.gauge.good { border-color: var(--good); } .gauge.needs-improvement { border-color: var(--warn); } .gauge.poor { border-color: var(--bad); }
.gauge small { color: var(--muted); }
.shift-banner { background: var(--warn); color: #000; padding: 60px 22px; border-radius: 10px; margin-bottom: 10px; font-weight: 600; }
table { border-collapse: collapse; width: 100%; font-size: .86rem; }
th, td { text-align: left; padding: 5px 8px; border-bottom: 1px solid var(--line); vertical-align: top; }
th { color: var(--muted); font-weight: 600; }
.sample { border: 1px dashed var(--line); border-radius: 8px; padding: 8px 12px; margin: 8px 0; }
.sample h4 { margin: 0; }
.sample h4::before { content: "★ "; color: var(--paint); }
.sample span { margin-right: 10px; }
[role="tab"][aria-selected="true"] { background: var(--accent); color: #fff; border-color: transparent; }
.log li[data-kind="bad"], .log li[data-kind="error"], .log li[data-kind="unhandledrejection"] { color: var(--bad); }
.log li[data-kind="csp"] { color: var(--warn); }
.log li[data-kind="ok"] { color: var(--good); }
.shift-banner.expanded { padding-bottom: 240px; }
```

## public/js/lib/dom.js

```javascript
// Tiny DOM helpers shared by every station.

// h("div", { class: "x", onclick: fn }, "text", childNode, [more]) -> Element.
// Text always goes in as text nodes, never as HTML.
export function h(tag, attrs = {}, ...children) {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs ?? {})) {
    if (k.startsWith("on") && typeof v === "function") el.addEventListener(k.slice(2), v);
    else if (k === "style" && typeof v === "object") Object.assign(el.style, v);
    else if (v === true) el.setAttribute(k, "");
    else if (v !== false && v != null) el.setAttribute(k, String(v));
  }
  for (const c of children.flat(Infinity)) {
    if (c != null && c !== false) el.append(c instanceof Node ? c : document.createTextNode(String(c)));
  }
  return el;
}

export const $ = (sel, root = document) => root.querySelector(sel);
export const $$ = (sel, root = document) => [...root.querySelectorAll(sel)];
export const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
export const nextFrame = () => new Promise((r) => requestAnimationFrame(r));
export const ms = (n) => (n == null || Number.isNaN(n) ? "—" : n < 10 ? `${n.toFixed(1)} ms` : `${Math.round(n)} ms`);
export const bytes = (n) =>
  n == null ? "—" : n < 1024 ? `${n} B` : n < 1048576 ? `${(n / 1024).toFixed(1)} KB` : n < 1073741824 ? `${(n / 1048576).toFixed(1)} MB` : `${(n / 1073741824).toFixed(1)} GB`;

// A station section: title, a one-line lede, and a body.
export function station(title, lede, ...body) {
  return h("article", { class: "station" }, h("header", {}, h("h2", {}, title), h("p", { class: "lede" }, lede)), ...body);
}

export function panel(title, ...body) {
  return h("section", { class: "panel", "aria-label": title }, h("h3", {}, title), ...body);
}

export function button(label, onclick, attrs = {}) {
  return h("button", { type: "button", onclick, ...attrs }, label);
}

// Busy-wait on purpose: the Rendering, Vitals and Threads stations need to
// block the main thread to show what blocking does.
export function blockFor(duration) {
  const end = performance.now() + duration;
  while (performance.now() < end);
}

export function storage(key, fallback) {
  // localStorage can throw (private mode, blocked site data): never let that break the lab.
  return {
    get() {
      try {
        const v = localStorage.getItem(key);
        return v == null ? fallback : JSON.parse(v);
      } catch {
        return fallback;
      }
    },
    set(v) {
      try {
        localStorage.setItem(key, JSON.stringify(v));
      } catch {}
    },
  };
}
```

## public/js/lib/errors.js

```javascript
// Global error capture, installed at page load (§83): what an error-tracking
// SDK such as Sentry does first, minus the network upload.
export const captured = [];
const listeners = new Set();
export const onCaptured = (fn) => (listeners.add(fn), () => listeners.delete(fn));

function record(kind, detail) {
  captured.unshift({ kind, at: performance.now(), ...detail });
  captured.length = Math.min(captured.length, 30);
  listeners.forEach((fn) => fn(captured));
}

addEventListener("error", (e) =>
  record("error", { message: e.message, where: `${e.filename?.split("/").pop()}:${e.lineno}:${e.colno}`, stack: e.error?.stack }),
);
addEventListener("unhandledrejection", (e) =>
  record("unhandledrejection", { message: String(e.reason?.message ?? e.reason), stack: e.reason?.stack }),
);
// Fired by the browser whenever the page's CSP blocks something.
addEventListener("securitypolicyviolation", (e) =>
  record("csp", { message: `${e.effectiveDirective} blocked ${e.blockedURI || "inline"}`, where: `${e.sourceFile?.split("/").pop() ?? "?"}:${e.lineNumber}` }),
);
// Deprecations and browser interventions, reported by the browser itself.
if ("ReportingObserver" in window) {
  new ReportingObserver((reports) => {
    // csp-violation reports duplicate the event above; they show the other channel.
    for (const r of reports) {
      const b = r.body ?? {};
      record(`report:${r.type}`, { message: b.message ?? `${b.effectiveDirective ?? ""} ${b.blockedURL ?? ""}`.trim() });
    }
  }, { buffered: true }).observe();
}
```

## public/js/lib/loop-model.js

```javascript
// A model of the HTML event loop (§12, §81): it "executes" a snippet
// described as a list of operations and records a snapshot after every
// step: the call stack, timers held by the browser, the task queue, the
// microtask queue, and the console.
//
// Operations (`at` is a substring of the source line, for highlighting):
//   { at, log: "text" }
//   { at, timeout: { delay, label, body } }       setTimeout(fn, delay)
//   { at, then:    { label, body, next? } }       Promise.resolve().then(fn)[.then(next)]
//   { at, micro:   { label, body } }              queueMicrotask(fn)
//   { at, call:    { name, body } }               fn() — sync call (async fns too)
//   { at, await:   { label, rest } }              await <settled value>: the rest of the
//                                                 function becomes one microtask

export function simulate(ops, source = "") {
  const lines = source.split("\n");
  const lineOf = (at) => (at ? lines.findIndex((l) => l.includes(at)) : -1);
  const S = { stack: [], webapis: [], tasks: [], micro: [], console: [], now: 0 };
  const steps = [];
  let seq = 0;

  const snap = (kind, note, at) =>
    steps.push({
      kind,
      note,
      line: lineOf(at),
      stack: [...S.stack],
      webapis: S.webapis.map((t) => `${t.label} (${t.delay}ms)`),
      tasks: S.tasks.map((t) => t.label),
      micro: S.micro.map((m) => m.label),
      console: [...S.console],
    });

  function exec(body) {
    for (const op of body) {
      if ("log" in op) {
        S.console.push(op.log);
        snap("log", `console.log prints "${op.log}"`, op.at);
      } else if (op.timeout) {
        const t = op.timeout;
        S.webapis.push({ ...t, due: S.now + t.delay, seq: seq++ });
        snap("webapi", `setTimeout hands "${t.label}" to the browser's timer (${t.delay}ms). JS moves on immediately.`, op.at);
      } else if (op.then) {
        S.micro.push(op.then);
        snap("micro", `The promise is already resolved, so "${op.then.label}" goes straight onto the microtask queue.`, op.at);
      } else if (op.micro) {
        S.micro.push(op.micro);
        snap("micro", `queueMicrotask puts "${op.micro.label}" on the microtask queue.`, op.at);
      } else if (op.call) {
        S.stack.push(`${op.call.name}()`);
        snap("push", `Calling ${op.call.name}() pushes a new frame onto the call stack.`, op.at);
        exec(op.call.body);
        S.stack.pop();
        snap("pop", `${op.call.name}() returns: its frame is popped.`, op.at);
      } else if (op.await) {
        S.micro.push({ label: op.await.label, body: op.await.rest });
        snap("micro", `await suspends the function. Its continuation "${op.await.label}" is queued as a microtask.`, op.at);
        return; // the rest of this body runs later, as that microtask
      }
    }
  }

  function drainMicrotasks() {
    if (!S.micro.length) return;
    snap("checkpoint", "The call stack is empty: microtask checkpoint. Drain the WHOLE microtask queue, including microtasks queued along the way.");
    while (S.micro.length) {
      const m = S.micro.shift();
      S.stack.push(m.label);
      snap("run-micro", `Run microtask "${m.label}".`, m.at);
      exec(m.body);
      S.stack.pop();
      if (m.next) {
        S.micro.push(m.next);
        snap("micro", `"${m.label}" returned, which resolves the next promise in the chain: "${m.next.label}" is queued.`);
      } else {
        snap("pop", `Microtask "${m.label}" done.`);
      }
    }
  }

  function moveExpiredTimers() {
    const due = S.webapis.filter((t) => t.due <= S.now).sort((a, b) => a.due - b.due || a.seq - b.seq);
    for (const t of due) {
      S.webapis.splice(S.webapis.indexOf(t), 1);
      S.tasks.push(t);
      snap("task-queued", `Timer "${t.label}" expired: its callback joins the task (macrotask) queue.`);
    }
  }

  S.stack.push("<script>");
  snap("task", "The whole <script> is the first task. It runs to completion on the call stack.");
  exec(ops);
  S.stack.pop();
  snap("pop", "The script is done and the call stack is empty.");
  drainMicrotasks();

  for (let guard = 0; guard < 500; guard++) {
    moveExpiredTimers();
    if (!S.tasks.length) {
      if (!S.webapis.length) break;
      S.now = Math.min(...S.webapis.map((t) => t.due));
      continue;
    }
    const task = S.tasks.shift();
    S.stack.push(task.label);
    snap("task", `The event loop takes ONE task from the task queue: "${task.label}".`);
    exec(task.body);
    S.stack.pop();
    snap("pop", `Task "${task.label}" done.`);
    drainMicrotasks();
  }
  snap("idle", "Every queue is empty. The event loop sleeps until the next event, timer or network response.");
  return { steps, output: S.console };
}
```

## public/js/lib/missions.js

```javascript
// Missions: two or three small challenges per station, ticked off
// automatically when you do the thing. Progress persists per browser.
import { h, storage } from "./dom.js";

export const MISSIONS = {
  overview: [["tour", "Play the tour of a click, start to finish"]],
  network: [
    ["waterfall", "Fire the request burst and read the waterfall"],
    ["whoami", "See which HTTP version and TLS cipher you negotiated"],
    ["tls", "Step through the TLS 1.3 handshake"],
  ],
  render: [
    ["composite", "Animate with transform: skip Layout and Paint"],
    ["thrash", "Measure layout thrashing against batched reads"],
  ],
  eventloop: [
    ["predict", "Predict a snippet's output correctly"],
    ["verify", "Run a snippet for real and match the model"],
  ],
  engine: [
    ["shapes", "Make two objects with different hidden classes"],
    ["megamorphic", "Push a call site to megamorphic"],
  ],
  threads: [
    ["jank", "Freeze the ball with main-thread work, then fix it with a worker"],
    ["race", "Lose updates in a data race, then fix it with Atomics"],
  ],
  storage: [
    ["idb", "Write and read back a record in IndexedDB"],
    ["inactive", "Trigger a TransactionInactiveError on purpose"],
  ],
  vitals: [
    ["cls", "Push CLS past 0.1"],
    ["inp", "Make INP poor, then make it good with yielding"],
  ],
  debug: [
    ["errors", "Catch an uncaught error and an unhandled rejection"],
    ["csp", "Trigger a CSP violation report"],
    ["leak", "Create a memory leak and find it"],
  ],
};

const store = storage("browser-lab:missions", {});
const done = store.get();
const listeners = new Set();

export const missions = {
  isDone: (station, id) => !!done[`${station}:${id}`],
  complete(station, id) {
    const key = `${station}:${id}`;
    if (done[key]) return;
    done[key] = Date.now();
    store.set(done);
    const label = MISSIONS[station]?.find(([m]) => m === id)?.[1];
    if (label) toast(`Mission complete: ${label}`);
    listeners.forEach((fn) => fn());
  },
  progress(station) {
    const list = station ? MISSIONS[station] : Object.entries(MISSIONS).flatMap(([s, l]) => l.map(([id]) => [s, id]));
    const n = station ? list.filter(([id]) => done[`${station}:${id}`]).length : list.filter(([s, id]) => done[`${s}:${id}`]).length;
    return { done: n, total: list.length };
  },
  onChange: (fn) => (listeners.add(fn), () => listeners.delete(fn)),
  reset() {
    for (const k of Object.keys(done)) delete done[k];
    store.set(done);
    listeners.forEach((fn) => fn());
  },
};

// The mission checklist shown at the top of each station.
export function missionList(station) {
  const ul = h("ul", { class: "missions", "aria-label": "Missions" });
  const render = () =>
    ul.replaceChildren(
      ...MISSIONS[station].map(([id, label]) =>
        h("li", { class: missions.isDone(station, id) ? "done" : "", "data-mission": id }, label),
      ),
    );
  render();
  missions.onChange(render);
  return ul;
}

function toast(text) {
  let stack = document.getElementById("toasts");
  if (!stack) document.body.append((stack = h("div", { id: "toasts", role: "status" })));
  const el = h("div", { class: "toast" }, "🏅 ", text);
  stack.append(el);
  setTimeout(() => el.remove(), 3500);
}
```

## public/js/lib/primes.js

```javascript
// Deliberately slow trial division: a stand-in for any CPU-heavy work
// (parsing a big file, image filters, diffing, search indexing).
export function countPrimes(limit) {
  let count = 0;
  for (let n = 2; n < limit; n++) {
    let prime = true;
    for (let d = 2; d * d <= n; d++) {
      if (n % d === 0) {
        prime = false;
        break;
      }
    }
    if (prime) count++;
  }
  return count;
}
```

## public/js/lib/vitals.js

```javascript
// Core Web Vitals, collected from page load onward — the same algorithms the
// web-vitals library uses, in ~100 lines (§83). Started by main.js on load,
// so the numbers cover the whole visit, not just the Vitals station.

export const THRESHOLDS = {
  LCP: [2500, 4000], // ms
  INP: [200, 500], // ms
  CLS: [0.1, 0.25], // unitless
  FCP: [1800, 3000],
  TTFB: [800, 1800],
};

export const rating = (name, v) =>
  v == null ? "none" : v <= THRESHOLDS[name][0] ? "good" : v <= THRESHOLDS[name][1] ? "needs-improvement" : "poor";

export const vitals = { LCP: null, INP: null, CLS: 0, FCP: null, TTFB: null, lcpElement: null, interactions: 0, last: null };
export const loafs = []; // recent long animation frames
const listeners = new Set();
export const onVitals = (fn) => (listeners.add(fn), () => listeners.delete(fn));
const emit = () => listeners.forEach((fn) => fn(vitals));

const supported = PerformanceObserver.supportedEntryTypes ?? [];
const observe = (type, cb, opts = {}) => {
  if (!supported.includes(type)) return false;
  new PerformanceObserver((list) => cb(list.getEntries())).observe({ type, buffered: true, ...opts });
  return true;
};
export const support = {
  lcp: supported.includes("largest-contentful-paint"),
  cls: supported.includes("layout-shift"),
  inp: supported.includes("event"),
  loaf: supported.includes("long-animation-frame"),
};

// TTFB: from the navigation entry.
const nav = performance.getEntriesByType("navigation")[0];
if (nav) vitals.TTFB = nav.responseStart - nav.startTime;

// FCP
observe("paint", (entries) => {
  for (const e of entries) if (e.name === "first-contentful-paint") vitals.FCP = e.startTime;
  emit();
});

// LCP: the last candidate reported before the first input wins.
observe("largest-contentful-paint", (entries) => {
  const last = entries.at(-1);
  vitals.LCP = last.startTime;
  vitals.lcpElement = last.element;
  emit();
});

// CLS: shifts are grouped into "session windows" (gap < 1s, window < 5s);
// CLS is the largest window. Shifts right after input don't count.
let windowValue = 0;
let windowEntries = [];
observe("layout-shift", (entries) => {
  for (const e of entries) {
    if (e.hadRecentInput) continue;
    const first = windowEntries[0];
    const last = windowEntries.at(-1);
    if (last && e.startTime - last.startTime < 1000 && e.startTime - first.startTime < 5000) {
      windowValue += e.value;
      windowEntries.push(e);
    } else {
      windowValue = e.value;
      windowEntries = [e];
    }
    vitals.CLS = Math.max(vitals.CLS, windowValue);
  }
  emit();
});

// INP: the slowest interaction (or the 98th percentile once there are
// 50+ interactions). One interaction = all events sharing an interactionId.
const byInteraction = new Map();
const updateINP = (entries) => {
  for (const e of entries) {
    if (!e.interactionId) continue;
    const d = Math.max(byInteraction.get(e.interactionId) ?? 0, e.duration);
    byInteraction.set(e.interactionId, d);
    vitals.last = { name: e.name, duration: d, target: e.target?.textContent?.trim().slice(0, 40) };
  }
  const all = [...byInteraction.values()].sort((a, b) => b - a);
  vitals.interactions = all.length;
  vitals.INP = all.length ? all[Math.min(all.length - 1, Math.floor(all.length / 50))] : null;
  emit();
};
observe("event", updateINP, { durationThreshold: 16 });
observe("first-input", updateINP);

// Long Animation Frames: which frames were slow, and which scripts made them slow.
observe("long-animation-frame", (entries) => {
  for (const e of entries) {
    loafs.unshift({
      start: e.startTime,
      duration: e.duration,
      blocking: e.blockingDuration,
      scripts: e.scripts.map((s) => ({ invoker: s.invoker, source: s.sourceURL?.split("/").pop(), duration: s.duration })),
    });
  }
  loafs.length = Math.min(loafs.length, 15);
  emit();
});

export const resetINP = () => {
  byInteraction.clear();
  updateINP([]);
};
```

## public/js/main.js

```javascript
// The lab's router: one ES module per station, loaded with dynamic import()
// on first visit (§16). Vitals and error capture start here, at page load.
import { h, $, ms } from "./lib/dom.js";
import { missions, MISSIONS } from "./lib/missions.js";
import "./lib/vitals.js";
import "./lib/errors.js";

export const STATIONS = [
  { id: "overview", icon: "🗺️", title: "The big picture" },
  { id: "network", icon: "🌐", title: "URL to bytes: TLS & HTTP/2" },
  { id: "render", icon: "🎨", title: "The rendering pipeline" },
  { id: "eventloop", icon: "🔁", title: "The event loop" },
  { id: "engine", icon: "⚙️", title: "Inside V8" },
  { id: "threads", icon: "🧵", title: "Workers & threads" },
  { id: "storage", icon: "🗄️", title: "IndexedDB & storage" },
  { id: "vitals", icon: "📈", title: "Core Web Vitals" },
  { id: "debug", icon: "🐞", title: "Debugging lab" },
];

const stage = $("#stage");
let unmount = null;

function renderNav(active) {
  $("#station-nav").replaceChildren(
    ...STATIONS.map((s, i) => {
      const p = missions.progress(s.id);
      return h("li", {},
        h("a", { href: `#/${s.id}`, "aria-current": s.id === active ? "page" : false },
          h("span", { class: "num" }, String(i).padStart(2, "0")),
          h("span", { class: "icon", "aria-hidden": "true" }, s.icon),
          h("span", { class: "label" }, s.title),
          h("span", { class: `badge ${p.done === p.total ? "complete" : ""}`, title: "missions done" }, `${p.done}/${p.total}`)));
    }),
  );
  const all = missions.progress();
  $("#progress-fill").style.width = `${(all.done / all.total) * 100}%`;
  $("#progress-text").textContent = `${all.done} of ${all.total} missions`;
}

async function route() {
  const id = location.hash.replace(/^#\/?/, "") || "overview";
  const meta = STATIONS.find((s) => s.id === id) ?? STATIONS[0];
  renderNav(meta.id);
  unmount?.();
  unmount = null;
  const t0 = performance.now();
  // Each station is its own module: watch them arrive in the Network panel.
  const mod = await import(`./stations/${meta.id}.js`);
  const loadMs = performance.now() - t0;
  const view = await mod.mount({ meta, loadMs });
  stage.replaceChildren(view.el);
  unmount = view.unmount ?? null;
  document.title = `${meta.title} · Browser Lab`;
  stage.focus({ preventScroll: true });
  stage.scrollTo?.(0, 0);
  // A User Timing mark: it shows up in DevTools' Performance panel (§83).
  performance.mark("station-shown", { detail: { station: meta.id, loadMs } });
}

addEventListener("hashchange", route);
missions.onChange(() => renderNav(location.hash.replace(/^#\/?/, "") || "overview"));
$("#reset-missions").addEventListener("click", () => missions.reset());

const nav = performance.getEntriesByType("navigation")[0];
$("#proto").textContent = `${nav?.nextHopProtocol || "?"} · ${location.protocol === "https:" ? "TLS" : "plain-text"} · ${ms(nav?.responseStart)} TTFB`;

route();

// Expose a tiny API for the Playwright suite and for poking around in the console.
window.lab = { missions, MISSIONS, STATIONS };
```

## public/js/stations/debug.js

```javascript
// Station 8: debugging (§50, §83). Errors are captured by lib/errors.js from
// page load; this station triggers them on purpose and shows what you get.
import { h, station, panel, button, ms, bytes } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";
import { captured, onCaptured } from "../lib/errors.js";

const leaked = []; // module scope: survives leaving the station, like a real leak

function measureWork() {
  performance.mark("parse:start");
  const data = JSON.parse(JSON.stringify(Array.from({ length: 50000 }, (_, i) => ({ i, s: `item ${i}` }))));
  performance.mark("parse:end");
  performance.mark("sort:start");
  data.sort((a, b) => b.s.localeCompare(a.s));
  performance.mark("sort:end");
  const m1 = performance.measure("parse 50k items", { start: "parse:start", end: "parse:end", detail: { devtools: { track: "Browser Lab", color: "primary" } } });
  const m2 = performance.measure("sort 50k items", { start: "sort:start", end: "sort:end", detail: { devtools: { track: "Browser Lab", color: "secondary" } } });
  return [m1, m2];
}

export function mount() {
  const list = h("ol", { class: "log", "aria-label": "Captured errors" });
  const render = () => {
    list.replaceChildren(...captured.map((c) => h("li", { "data-kind": c.kind }, `[${c.kind}] ${c.message}${c.where ? ` @ ${c.where}` : ""}`)));
    const kinds = new Set(captured.map((c) => c.kind));
    if (kinds.has("error") && kinds.has("unhandledrejection")) missions.complete("debug", "errors");
    if (kinds.has("csp")) missions.complete("debug", "csp");
  };
  const off = onCaptured(render);
  render();

  const memOut = h("dl", { class: "kv" });
  const leakInfo = h("p", { "aria-live": "polite" }, `Leaked nodes held: ${leaked.length.toLocaleString()}`);
  const measureMemory = async () => {
    memOut.replaceChildren(h("dt", {}, "Measuring…"), h("dd", {}, "(can take a few seconds: it waits for a garbage collection)"));
    if (crossOriginIsolated && performance.measureUserAgentSpecificMemory) {
      const r = await performance.measureUserAgentSpecificMemory();
      memOut.replaceChildren(h("dt", {}, "Total (all realms)"), h("dd", {}, bytes(r.bytes)),
        ...r.breakdown.filter((b) => b.bytes).slice(0, 5).flatMap((b) => [h("dt", {}, b.types.join(", ") || "other"), h("dd", {}, bytes(b.bytes))]));
    } else if (performance.memory) {
      memOut.replaceChildren(h("dt", {}, "JS heap used (Chrome-only, coarse)"), h("dd", {}, bytes(performance.memory.usedJSHeapSize)));
    } else {
      memOut.replaceChildren(h("dt", {}, "Not available"), h("dd", {}, "this browser exposes no memory API; use DevTools → Memory"));
    }
    if (leaked.length) missions.complete("debug", "leak");
  };

  const marksOut = h("ol", { class: "log", "aria-label": "Measures" });

  const el = station(
    "Debugging lab",
    "Production bugs rarely announce themselves. Here you trigger the common ones on purpose (uncaught errors, unhandled rejections, CSP violations, leaks, slow code) and see what the browser gives you to catch them.",
    missionList("debug"),
    h("div", { class: "grid-2" },
      panel("Trigger something",
        h("div", { class: "row" },
          button("Throw an uncaught error", () => setTimeout(() => { throw new Error("Boom from a timer callback"); })),
          button("Reject a promise, never catch it", () => { Promise.reject(new Error("Nobody handled me")); }),
          button("Violate the CSP (eval)", () => {
            try {
              new Function("return 1")(); // script-src has no 'unsafe-eval'
            } catch (err) {
              console.info("blocked as expected:", err.name);
            }
          }),
          button("Pause in the debugger", () => {
            debugger; // only stops when DevTools is open
          })),
        h("p", { class: "note" }, "These go to window 'error', 'unhandledrejection' and 'securitypolicyviolation' listeners installed at page load: exactly what an error tracker hooks into.")),
      panel("Captured (newest first)", list)),
    h("div", { class: "grid-2" },
      panel("A memory leak you can find",
        h("p", {}, "Each click creates 10,000 DOM nodes, never attaches them, and keeps them in a module-level array. They can never be collected: detached DOM nodes."),
        h("div", { class: "row" },
          button("Leak 10k nodes", () => {
            for (let i = 0; i < 10000; i++) leaked.push(h("div", { class: "leak" }, `leaked ${i}`));
            leakInfo.textContent = `Leaked nodes held: ${leaked.length.toLocaleString()}`;
          }, { class: "danger" }),
          button("Free them", () => {
            leaked.length = 0;
            leakInfo.textContent = "Leaked nodes held: 0 (collectable at the next GC)";
          }),
          button("Measure memory", measureMemory)),
        leakInfo, memOut,
        h("p", { class: "note" }, "Find it: DevTools → Memory → Heap snapshot → filter \"Detached\". Expand a HTMLDivElement → Retainers shows the array called leaked in debug.js.")),
      panel("User Timing: your own marks in the Performance panel",
        button("Run measured work", () => {
          for (const m of measureWork()) marksOut.prepend(h("li", {}, `${m.name}: ${ms(m.duration)}`));
        }, { class: "primary" }),
        marksOut,
        h("p", { class: "note" }, "Record in DevTools → Performance, click the button, stop. Your measures appear in the Timings track, and in a custom \"Browser Lab\" track in Chrome (the detail.devtools extension)."))),
    panel("DevTools recipes worth knowing",
      h("table", {},
        h("thead", {}, h("tr", {}, h("th", {}, "Problem"), h("th", {}, "Tool"), h("th", {}, "How"))),
        h("tbody", {}, [
          ["Which code changed this element?", "DOM breakpoints", "Elements → right-click node → Break on → subtree/attribute modifications"],
          ["Log without editing code", "Logpoints", "Sources → right-click line number → Add logpoint"],
          ["Stop only for one bad value", "Conditional breakpoint", "right-click line number → condition like id === 'A-1002'"],
          ["Who sent this request?", "Initiator + fetch breakpoints", "Network → Initiator column; Sources → XHR/fetch Breakpoints"],
          ["Where did this async call come from?", "Async stack traces", "Call Stack panel shows the await/then chain, on by default"],
          ["Why is this slow?", "Performance panel", "Record, then read the flame chart bottom-up; long tasks have red corners"],
          ["What repaints? What shifts?", "Rendering tab", "Paint flashing, Layout Shift Regions, Core Web Vitals overlay"],
          ["What would break offline / on 3G?", "Network conditions", "Throttling presets, request blocking, offline mode"],
          ["Fix prod CSS/JS without a deploy", "Local overrides", "Sources → Overrides → select a folder, then edit and save"],
          ["Unused code", "Coverage", "More tools → Coverage → reload: red = never executed"],
        ].map((r) => h("tr", {}, r.map((c) => h("td", {}, c))))))),
  );
  return { el, unmount: off };
}
```

## public/js/stations/engine.js

```javascript
// Station 4: inside V8 (§49, §81). The tier and shape diagrams are models
// (V8's real thresholds depend on feedback and budgets and change between
// versions); the benchmark at the bottom is real.
import { h, station, panel, button, ms } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";

const TIERS = [
  ["parse", "Parser", "source → AST", "lazy: inner functions are only pre-parsed until first called"],
  ["ignition", "Ignition", "interpreter", "runs compact bytecode and collects type feedback"],
  ["sparkplug", "Sparkplug", "baseline JIT", "bytecode → machine code with no optimization; very fast to compile"],
  ["maglev", "Maglev", "mid-tier JIT", "quick optimizing compiler (Chrome 117+) that uses the feedback"],
  ["turbofan", "TurboFan", "top-tier JIT", "speculative, heavily optimized code for the hottest functions"],
];
// Illustrative call counts for the animation only.
const STEPS = [[0, "ignition"], [500, "sparkplug"], [5000, "maglev"], [50000, "turbofan"]];

const BYTECODE = `// node --print-bytecode --print-bytecode-filter=add  (V8 12.4)
function add(a, b) { return a + b; }

[generated bytecode for function: add]
Parameter count 3        // a, b, and the receiver (this)
Register count 0
   21 S> @ 0 : 0b 04       Ldar a1          // accumulator = b
   30 E> @ 2 : 3b 03 00    Add a0, [0]      // accumulator = a + accumulator; [0] = feedback slot
   34 S> @ 5 : ae          Return           // return the accumulator`;

function tiersPanel() {
  let calls = 0;
  let tier = "ignition";
  let deopts = 0;
  const add = (a, b) => a + b; // a real function, really called
  const tierEls = TIERS.map(([id, name, short, desc]) => h("div", { class: "tier", "data-tier": id, title: desc }, name, h("small", {}, short)));
  const heat = h("span");
  const status = h("p", { "aria-live": "polite" });
  const render = (msg) => {
    tierEls.forEach((t) => t.classList.toggle("on", t.dataset.tier === tier));
    heat.style.width = `${Math.min(100, (Math.log10(calls + 1) / Math.log10(60000)) * 100)}%`;
    status.textContent = msg ?? `add() called ${calls.toLocaleString()} times · running in ${TIERS.find((t) => t[0] === tier)[1]} · deopts: ${deopts}`;
  };
  const call = (n) => {
    let s = 0;
    for (let i = 0; i < n; i++) s = add(s, 1);
    calls += n;
    tier = STEPS.filter(([c]) => calls >= c).at(-1)[1];
    render();
  };
  const deopt = () => {
    add("oops", 1); // a string where feedback said "always small integers"
    const was = tier;
    if (was === "maglev" || was === "turbofan") {
      deopts++;
      tier = "ignition";
      calls = 0;
      render(`💥 Deoptimized! ${TIERS.find((t) => t[0] === was)[1]}'s code assumed numbers; a string broke the assumption, so V8 threw that code away and went back to Ignition. It can re-optimize later with wider feedback.`);
    } else render("Called with a string. Still in an unoptimized tier, so nothing to throw away: the feedback just gets more general.");
  };
  render();
  return panel("From source to machine code",
    h("div", { class: "tiers" }, tierEls.flatMap((t, i) => (i ? [h("span", { class: "arrow" }, "→"), t] : [t]))),
    h("div", { class: "heat", "aria-hidden": "true" }, heat),
    status,
    h("div", { class: "row" }, button("Call add() ×1,000", () => call(1000)), button("×10,000", () => call(10000)), button("×100,000", () => call(100000), { class: "primary" }), button("Call add('oops', 1)", deopt, { class: "danger" })),
    h("p", { class: "note" }, "Tier-up thresholds here are illustrative. Try it for real: node --trace-opt --trace-deopt file.js prints every optimization and bailout."),
    h("pre", { class: "code" }, BYTECODE.split("\n").map((l) => h("span", {}, l))));
}

function shapesPanel() {
  // A transition tree: each node is a hidden class (V8 calls them Maps).
  const root = { props: [], children: new Map(), count: 0, id: 0 };
  let nextId = 1;
  const dictionary = { count: 0 };
  const seenByIC = new Set();
  const tree = h("div", { class: "shapes" });
  const ic = h("p", { "aria-live": "polite" });

  const shapeFor = (props) => {
    let node = root;
    for (const p of props) {
      if (!node.children.has(p)) node.children.set(p, { props: [...node.props, p], children: new Map(), count: 0, id: nextId++ });
      node = node.children.get(p);
    }
    return node;
  };
  const icState = (n) => (n === 0 ? "uninitialized" : n === 1 ? "monomorphic" : n <= 4 ? "polymorphic" : "megamorphic");
  const draw = () => {
    const walk = (node) =>
      h("li", {},
        h("span", { class: "shape" }, `Map${node.id} {${node.props.join(", ")}}`, node.count ? h("span", { class: "count" }, `×${node.count}`) : null),
        node.children.size ? h("ul", {}, [...node.children.values()].map(walk)) : null);
    tree.replaceChildren(h("ul", {}, walk(root)),
      dictionary.count ? h("p", {}, h("span", { class: "shape" }, "dictionary mode (hash table) ", h("span", { class: "count" }, `×${dictionary.count}`))) : "");
    const state = icState(seenByIC.size);
    ic.replaceChildren("Inline cache at ", h("code", {}, "getX(obj)"), ` has seen ${seenByIC.size} shape(s): `, h("span", { class: `ic ${state}` }, state));
    const shapes = [];
    const collect = (n) => (n.count && shapes.push(n), n.children.forEach(collect));
    collect(root);
    if (shapes.length >= 2) missions.complete("engine", "shapes");
    if (state === "megamorphic") missions.complete("engine", "megamorphic");
  };
  const make = (props, label) => () => {
    const node = shapeFor(props);
    node.count++;
    if (props.includes("x")) seenByIC.add(node.id);
    draw();
    log.prepend(h("li", {}, label));
  };
  const log = h("ol", { class: "log", "aria-label": "Objects created" });
  let extra = 0;
  draw();
  return panel("Hidden classes and inline caches",
    h("p", {}, "V8 gives every object a hidden class describing its layout. Objects built the same way share one, so property access compiles to a fixed-offset load. The order you add properties matters."),
    h("div", { class: "row" },
      button("{ x, y }", make(["x", "y"], "const o = { x: 1, y: 2 }")),
      button("{ y, x }", make(["y", "x"], "const o = { y: 2, x: 1 }   // different order → different Map")),
      button("{ x } then o.y = …", make(["x", "y"], "const o = { x: 1 }; o.y = 2   // same transitions as { x, y }")),
      button("{ x, y, z }", make(["x", "y", "z"], "const o = { x, y, z }")),
      button("random extra prop", () => make([`p${++extra}`, "x"], `const o = { p${extra}: 0, x: 1 }   // yet another shape`)()),
      button("delete o.x", () => {
        dictionary.count++;
        draw();
        log.prepend(h("li", {}, "delete o.x   // the object drops to slow dictionary mode"));
      }, { class: "danger" })),
    h("div", { class: "grid-2" }, tree, h("div", {}, ic, log)),
    h("p", { class: "note" }, "Verify in Node: node --allow-natives-syntax -e \"const a={x:1};a.y=2;const b={y:2};b.x=1;console.log(%HaveSameMap(a,b))\" prints false."));
}

// Two identical function bodies, so each keeps its own type feedback.
function sumMono(objs, n) {
  let s = 0;
  for (let i = 0; i < n; i++) s += objs[i & (objs.length - 1)].x;
  return s;
}
function sumMega(objs, n) {
  let s = 0;
  for (let i = 0; i < n; i++) s += objs[i & (objs.length - 1)].x;
  return s;
}

function benchPanel() {
  const out = h("div", { class: "readouts" });
  const run = () => {
    const N = 20_000_000;
    const same = Array.from({ length: 8 }, () => ({ x: 1, y: 2 }));
    const mixed = Array.from({ length: 8 }, (_, i) => {
      const o = {};
      o[`p${i}`] = 0; // a different first property → 8 different hidden classes
      o.x = 1;
      return o;
    });
    sumMono(same, 10_000); // warm up both
    sumMega(mixed, 10_000);
    let t = performance.now();
    sumMono(same, N);
    const mono = performance.now() - t;
    t = performance.now();
    sumMega(mixed, N);
    const mega = performance.now() - t;
    out.replaceChildren(
      h("div", { class: "readout" }, h("small", {}, "1 shape (monomorphic)"), h("b", {}, ms(mono))),
      h("div", { class: "readout" }, h("small", {}, "8 shapes (megamorphic)"), h("b", {}, ms(mega))),
      h("div", { class: "readout" }, h("small", {}, "slowdown"), h("b", {}, `${(mega / mono).toFixed(1)}×`)));
  };
  return panel("Benchmark: the cost of megamorphic access (real)",
    h("p", {}, "Same loop, 20 million reads of .x. Left: 8 objects with one hidden class. Right: 8 objects with 8 different hidden classes."),
    button("Run benchmark (blocks ~1s)", run, { class: "primary" }), out,
    h("p", { class: "note" }, "Microbenchmarks lie a little: run it a few times. The lesson holds: hot code that sees many object shapes can't use a fast inline cache."));
}

function timeline() {
  const items = [
    ["2017", "TurboFan + Ignition replace Crankshaft/Full-codegen in V8 5.9"],
    ["2018", "Site Isolation ships in Chrome after Spectre; SharedArrayBuffer gated behind cross-origin isolation (§39, §51)"],
    ["2019", "await gets faster: awaiting a native promise takes 1 microtask tick instead of 3 (spec change + V8 7.2)"],
    ["2021", "Sparkplug baseline compiler (Chrome 91)"],
    ["2023", "Maglev mid-tier compiler (Chrome 117); Turboshaft begins replacing TurboFan's backend"],
    ["2024", "INP replaces FID as a Core Web Vital; Long Animation Frames API (Chrome 123); scheduler.yield() (Chrome 129)"],
    ["2025", "Iterator helpers, Set methods, Promise.try, RegExp.escape, Float16Array in ES2025"],
    ["2025–26", "Temporal ships (Firefox 139, Chrome 144) and using/await using ships (Chrome 134, Firefox 141); both finished, slated for ES2027"],
    ["2026", "ES2026: Array.fromAsync, Error.isError, Math.sumPrecise, Uint8Array base64/hex, Iterator.concat, Map getOrInsert"],
  ];
  return panel("Recent engine and event-loop milestones",
    h("table", {}, h("tbody", {}, items.map(([y, t]) => h("tr", {}, h("th", {}, y), h("td", {}, t))))));
}

export function mount() {
  return {
    el: station(
      "Inside V8",
      "JavaScript is compiled at runtime, several times. V8 starts fast with an interpreter, watches what types actually flow through your code, and recompiles the hot parts with assumptions built in. When an assumption breaks, it bails out.",
      missionList("engine"),
      tiersPanel(),
      shapesPanel(),
      benchPanel(),
      timeline(),
    ),
  };
}
```

## public/js/stations/eventloop.js

```javascript
// Station 3: the event loop (§12–§15, §81). Each snippet exists twice:
// as real code (`run`) and as a model (`ops`). The model drives the
// animation; "Run it for real" executes `run` and checks they agree.
import { h, station, panel, button, sleep, storage } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";
import { simulate } from "../lib/loop-model.js";

export const SNIPPETS = [
  {
    id: "basics",
    title: "Sync, microtasks, tasks",
    run: (console) => {
      console.log("script start");
      setTimeout(() => console.log("timeout"), 0);
      Promise.resolve()
        .then(() => console.log("promise 1"))
        .then(() => console.log("promise 2"));
      console.log("script end");
    },
    ops: [
      { at: '"script start"', log: "script start" },
      { at: "setTimeout", timeout: { delay: 0, label: "timeout cb", body: [{ at: '"timeout"', log: "timeout" }] } },
      { at: ".then(() => console.log(\"promise 1\"))", then: { label: "then: promise 1", body: [{ at: '"promise 1"', log: "promise 1" }],
        next: { label: "then: promise 2", body: [{ at: '"promise 2"', log: "promise 2" }] } } },
      { at: '"script end"', log: "script end" },
    ],
  },
  {
    id: "await",
    title: "async / await",
    run: (console) => {
      async function load() {
        console.log("load() starts");
        await null;
        console.log("load() resumes");
      }
      console.log("start");
      load();
      queueMicrotask(() => console.log("microtask"));
      console.log("end");
    },
    ops: [
      { at: 'console.log("start")', log: "start" },
      { at: "load();", call: { name: "load", body: [
        { at: '"load() starts"', log: "load() starts" },
        { at: "await null", await: { label: "load() continuation", rest: [{ at: '"load() resumes"', log: "load() resumes" }] } },
      ] } },
      { at: "queueMicrotask", micro: { label: "microtask cb", body: [{ at: 'log("microtask")', log: "microtask" }] } },
      { at: 'console.log("end")', log: "end" },
    ],
  },
  {
    id: "nested",
    title: "Microtasks inside tasks",
    run: (console) => {
      setTimeout(() => {
        console.log("timeout 1");
        Promise.resolve().then(() => console.log("promise in timeout 1"));
      }, 0);
      setTimeout(() => console.log("timeout 2"), 0);
      queueMicrotask(() => {
        console.log("microtask 1");
        queueMicrotask(() => console.log("microtask 2"));
      });
      console.log("sync");
    },
    ops: [
      { at: "setTimeout(() => {", timeout: { delay: 0, label: "timeout 1 cb", body: [
        { at: '"timeout 1"', log: "timeout 1" },
        { at: "Promise.resolve().then", then: { label: "then in timeout 1", body: [{ at: '"promise in timeout 1"', log: "promise in timeout 1" }] } },
      ] } },
      { at: '"timeout 2"', timeout: { delay: 0, label: "timeout 2 cb", body: [{ at: '"timeout 2"', log: "timeout 2" }] } },
      { at: "queueMicrotask(() => {", micro: { label: "microtask 1 cb", body: [
        { at: '"microtask 1"', log: "microtask 1" },
        { at: 'queueMicrotask(() => console.log("microtask 2"))', micro: { label: "microtask 2 cb", body: [{ at: '"microtask 2"', log: "microtask 2" }] } },
      ] } },
      { at: '"sync"', log: "sync" },
    ],
  },
  {
    id: "chain",
    title: "await on an async function, timers with delays",
    run: (console) => {
      async function a() {
        console.log("a: before await");
        await b();
        console.log("a: after await");
      }
      async function b() {
        console.log("b runs");
      }
      setTimeout(() => console.log("timeout 10ms"), 10);
      setTimeout(() => console.log("timeout 0ms"), 0);
      a();
      Promise.resolve().then(() => console.log("then"));
      console.log("sync");
    },
    ops: [
      { at: '"timeout 10ms"', timeout: { delay: 10, label: "10ms cb", body: [{ at: '"timeout 10ms"', log: "timeout 10ms" }] } },
      { at: '"timeout 0ms"', timeout: { delay: 0, label: "0ms cb", body: [{ at: '"timeout 0ms"', log: "timeout 0ms" }] } },
      { at: "  a();", call: { name: "a", body: [
        { at: '"a: before await"', log: "a: before await" },
        { at: "await b()", call: { name: "b", body: [{ at: '"b runs"', log: "b runs" }] } },
        { at: "await b()", await: { label: "a() continuation", rest: [{ at: '"a: after await"', log: "a: after await" }] } },
      ] } },
      { at: "Promise.resolve().then", then: { label: "then cb", body: [{ at: 'log("then")', log: "then" }] } },
      { at: 'console.log("sync")', log: "sync" },
    ],
  },
];

// The real snippet's source, from Function.prototype.toString (which returns
// the exact source text), minus the wrapper and the common indentation.
export function sourceOf(fn) {
  const lines = fn.toString().split("\n").slice(1, -1);
  const indent = Math.min(...lines.filter((l) => l.trim()).map((l) => l.match(/^ */)[0].length));
  return lines.map((l) => l.slice(indent)).join("\n");
}

// Run the real code with a console that records instead of printing.
export async function runForReal(snippet) {
  const out = [];
  snippet.run({ log: (s) => out.push(s) });
  await sleep(60); // long enough for every timer in the snippets to fire
  return out;
}

const score = storage("browser-lab:quiz", {});

function snippetView(snippet) {
  const source = sourceOf(snippet.run);
  const { steps, output } = simulate(snippet.ops, source);
  let i = 0;
  let timer = null;

  const code = h("pre", { class: "code", "aria-label": "Snippet source" }, source.split("\n").map((l) => h("span", {}, l || " ")));
  const box = (cls, title) => h("div", { class: `box ${cls}` }, h("h4", {}, title), h("div", { class: "items" }));
  const boxes = {
    stack: box("stack", "Call stack"),
    webapis: box("webapis", "Web APIs (timers)"),
    micro: box("micro", "Microtask queue"),
    tasks: box("tasks", "Task queue"),
    console: box("console", "Console"),
  };
  const note = h("div", { class: "loop-note", "aria-live": "polite" });
  const wheel = h("div", { class: "wheel", title: "the event loop" });
  const counter = h("span", { class: "note mono" });

  const render = () => {
    const s = steps[i];
    // Patch the tokens in place: re-inserting a node restarts its CSS
    // animation, so kept tokens stay put and only new arrivals animate.
    const fill = (key, items) => {
      const holder = boxes[key].querySelector(".items");
      const pool = [...holder.children];
      const wanted = items.map((t) => {
        const i = pool.findIndex((n) => n.textContent === t);
        return i === -1 ? h("div", { class: "token" }, t) : pool.splice(i, 1)[0];
      });
      pool.forEach((n) => n.remove()); // tokens that left this box
      wanted.forEach((n, i) => {
        if (holder.children[i] !== n) holder.insertBefore(n, holder.children[i] ?? null);
      });
    };
    fill("stack", s.stack);
    fill("webapis", s.webapis);
    fill("micro", s.micro);
    fill("tasks", s.tasks);
    fill("console", s.console);
    [...code.children].forEach((ln, n) => ln.classList.toggle("hl", n === s.line));
    note.textContent = s.note;
    wheel.classList.toggle("spin", s.kind === "task" || s.kind === "checkpoint");
    counter.textContent = `step ${i + 1} / ${steps.length}`;
    prev.disabled = i === 0;
    next.disabled = i === steps.length - 1;
  };
  const go = (n) => {
    i = Math.max(0, Math.min(steps.length - 1, n));
    render();
  };
  const stop = () => {
    clearInterval(timer);
    timer = null;
    playBtn.textContent = "▶ Play";
  };
  const prev = button("◀ Back", () => (stop(), go(i - 1)));
  const next = button("Step ▶", () => (stop(), go(i + 1)));
  const playBtn = button("▶ Play", () => {
    if (timer) return stop();
    if (i === steps.length - 1) go(0);
    playBtn.textContent = "⏸ Pause";
    timer = setInterval(() => (i === steps.length - 1 ? stop() : go(i + 1)), 1100);
  }, { class: "primary" });

  // --- prediction quiz
  const picks = [];
  const answer = h("div", { class: "answer", "aria-live": "polite" });
  const shuffled = [...output].sort(() => Math.random() - 0.5);
  const quiz = h("div", { class: "quiz" },
    shuffled.map((line) => {
      const b = button(line, () => {
        if (b.classList.contains("picked")) return;
        b.classList.add("picked");
        picks.push(line);
        answer.textContent = picks.map((p, n) => `${n + 1}. ${p}`).join("   ");
        if (picks.length === output.length) {
          const right = picks.every((p, n) => p === output[n]);
          answer.append(h("div", { class: right ? "ok" : "bad" }, right ? "✔ Correct! Now step through to see why." : `✘ Not quite. Actual: ${output.join(" → ")}`));
          const s = score.get();
          s[snippet.id] = right;
          score.set(s);
          if (right) missions.complete("eventloop", "predict");
        }
      });
      return b;
    }));
  const resetQuiz = button("reset", () => {
    picks.length = 0;
    answer.textContent = "";
    quiz.querySelectorAll("button").forEach((b) => b.classList.remove("picked"));
  }, { class: "link" });

  // --- run for real
  const real = h("div", { class: "answer" });
  const realBtn = button("🧪 Run it for real", async () => {
    const out = await runForReal(snippet);
    const match = out.join("|") === output.join("|");
    real.replaceChildren(h("div", {}, `Real:  ${out.join(" → ")}`), h("div", {}, `Model: ${output.join(" → ")}`),
      h("div", { class: match ? "ok" : "bad", "data-match": String(match) }, match ? "✔ The model matches your browser." : "✘ Mismatch: your engine disagrees with the model!"));
    if (match) missions.complete("eventloop", "verify");
  });

  render();
  const el = h("div", { class: "loop-grid" },
    h("div", {},
      code,
      panel("1. Predict the console output", h("p", { class: "note" }, "Click the lines in the order you think they print."), quiz, resetQuiz, answer),
      panel("3. Check against your browser", realBtn, real)),
    h("div", {},
      h("div", { class: "row" }, wheel, prev, playBtn, next, counter),
      note,
      h("div", { class: "boxes" }, boxes.stack, boxes.webapis, boxes.micro, boxes.tasks),
      boxes.console));
  return { el, stop };
}

export function mount() {
  const holder = h("div", {});
  let current = null;
  const tabs = h("div", { class: "row", role: "tablist" });
  const select = (s) => {
    current?.stop();
    current = snippetView(s);
    holder.replaceChildren(panel(`2. Step through: ${s.title}`, current.el));
    tabs.querySelectorAll("button").forEach((b) => b.setAttribute("aria-selected", String(b.dataset.id === s.id)));
  };
  tabs.append(...SNIPPETS.map((s) => {
    const b = button(s.title, () => select(s), { role: "tab", "data-id": s.id });
    return b;
  }));
  select(SNIPPETS[0]);
  return {
    el: station(
      "The event loop",
      "One thread, one call stack, two queues that matter. After every task the engine drains the entire microtask queue, and only then takes the next task. Predict, step through, then run the real code to check.",
      missionList("eventloop"),
      tabs,
      holder,
      panel("The rules, in one place",
        h("ol", {},
          h("li", {}, "Run one task (a script, a timer callback, an event handler) to completion. Nothing interrupts it."),
          h("li", {}, "When the call stack empties, drain ALL microtasks: promise reactions, await continuations, queueMicrotask. Microtasks queued meanwhile run too."),
          h("li", {}, "If a frame is due: run requestAnimationFrame callbacks, then style → layout → paint."),
          h("li", {}, "Take the next task. Repeat forever.")),
        h("p", { class: "note" }, "Why setTimeout(fn, 0) waits: it is a task, and tasks only run after the current task and all microtasks. Nested timers are clamped to ≥4ms after 5 levels, and background tabs throttle timers to ≥1s."))),
    unmount: () => current?.stop(),
  };
}
```

## public/js/stations/network.js

```javascript
// Station 1: from URL to bytes (§25–§27, §80). Everything on this page is
// measured from the real requests your browser just made.
import { h, station, panel, button, sleep, ms } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";

// ------------------------------------------------ navigation phases ----
function navigationPhases() {
  const n = performance.getEntriesByType("navigation")[0];
  if (!n) return h("p", {}, "No navigation entry available.");
  const phases = [
    ["Redirects", n.redirectEnd - n.redirectStart],
    ["Service worker / cache", n.domainLookupStart - n.fetchStart],
    ["DNS lookup", n.domainLookupEnd - n.domainLookupStart],
    ["TCP connect", (n.secureConnectionStart || n.connectEnd) - n.connectStart],
    ["TLS handshake", n.secureConnectionStart ? n.connectEnd - n.secureConnectionStart : 0],
    ["Request → first byte", n.responseStart - n.requestStart],
    ["Download HTML", n.responseEnd - n.responseStart],
    ["Parse → DOMContentLoaded", n.domContentLoadedEventEnd - n.responseEnd],
  ];
  const max = Math.max(...phases.map(([, v]) => v), 1);
  return h("div", { class: "phases" },
    phases.map(([label, v], i) =>
      h("div", { class: "phase" },
        h("span", {}, label),
        h("div", {}, h("div", { class: "fill", style: { width: `${Math.max(0, (v / max) * 100)}%`, animationDelay: `${i * 90}ms` } })),
        h("span", { class: "mono" }, ms(Math.max(0, v))))),
    h("p", { class: "note" },
      `Protocol: ${n.nextHopProtocol || "?"} · transfer ${n.transferSize} B (0 means it came from cache) · `,
      n.secureConnectionStart ? "a fresh TLS connection" : location.protocol === "https:" ? "TLS connection reused" : "no TLS (plain HTTP)",
      ". On localhost DNS and TCP are near zero; on a real site they are often 50–300ms each, which is why connection reuse and preconnect matter."));
}

// ----------------------------------------------------- TLS 1.3 ----
const TLS13 = [
  ["c2s", "ClientHello", "versions, cipher suites, + key_share (an ECDHE public key, guessed up front)"],
  ["s2c", "ServerHello", "chosen suite + server key_share → both sides now derive the handshake keys"],
  ["s2c", "🔒 EncryptedExtensions, Certificate, CertificateVerify, Finished", "already encrypted: the cert chain is hidden from on-path observers"],
  ["c2s", "🔒 Finished + first HTTP request", "the client validated the cert chain and sends data right away: 1 RTT total"],
  ["s2c", "🔒 HTTP response", "application data, encrypted with traffic keys"],
];
const TLS12 = [
  ["c2s", "ClientHello", "versions, cipher suites, no key share yet"],
  ["s2c", "ServerHello, Certificate, ServerKeyExchange, ServerHelloDone", "certificate travels in the clear"],
  ["c2s", "ClientKeyExchange, ChangeCipherSpec, Finished", "second round trip"],
  ["s2c", "ChangeCipherSpec, Finished", "now both sides have keys"],
  ["c2s", "🔒 HTTP request", "data only after 2 RTTs"],
  ["s2c", "🔒 HTTP response", ""],
];

function tlsPanel() {
  const msgs = h("div", { class: "msgs", "aria-live": "polite" });
  let running = false;
  const play = async (script, rtts) => {
    if (running) return;
    running = true;
    msgs.replaceChildren();
    for (const [dir, title, detail] of script) {
      msgs.append(h("div", { class: `msg-line ${dir}` }, title, h("small", {}, detail)));
      await sleep(900);
    }
    msgs.append(h("div", { class: "rtt" }, `${rtts} before the first byte of HTTP data can be sent`));
    running = false;
    if (script === TLS13) missions.complete("network", "tls");
  };
  return panel("The TLS handshake, step by step",
    h("div", { class: "row" },
      button("▶ TLS 1.3 (1-RTT)", () => play(TLS13, "1 round trip"), { class: "primary" }),
      button("▶ TLS 1.2 (2-RTT)", () => play(TLS12, "2 round trips"))),
    h("div", { class: "seq" }, h("div", { class: "actor" }, "💻 Browser"), h("div", { class: "actor" }, "🖥️ Server"), msgs),
    h("p", { class: "note" }, "With QUIC (HTTP/3) the transport and TLS handshakes are merged, so a new connection costs 1 RTT in total, and resumed connections can send 0-RTT data (replayable, so only for idempotent requests)."));
}

// ------------------------------------------------- whoami ----
function whoamiPanel() {
  const out = h("div", {}, h("p", { class: "note" }, "Ask the server what it saw: the negotiated protocol, cipher, and every header your browser sent."));
  const go = async () => {
    const info = await (await fetch("/api/whoami")).json();
    const interesting = Object.entries(info.headers).filter(([k]) => /^(:|sec-|priority|accept-encoding|user-agent|cookie|referer|origin)/.test(k));
    out.replaceChildren(
      h("dl", { class: "kv" },
        h("dt", {}, "HTTP version"), h("dd", {}, info.httpVersion),
        h("dt", {}, "ALPN"), h("dd", {}, info.alpn ?? "— (no TLS)"),
        h("dt", {}, "TLS version"), h("dd", {}, info.tls ?? "—"),
        h("dt", {}, "Cipher"), h("dd", {}, info.cipher ?? "—"),
        h("dt", {}, "Client port"), h("dd", {}, String(info.remotePort))),
      h("table", {}, h("thead", {}, h("tr", {}, h("th", {}, "Request header"), h("th", {}, "Value"))),
        h("tbody", {}, interesting.map(([k, v]) => h("tr", {}, h("td", { class: "mono" }, k), h("td", { class: "mono" }, String(v)))))),
      h("p", { class: "note" }, "Headers starting with \":\" are HTTP/2 pseudo-headers. sec-fetch-* are Fetch Metadata, which servers use to reject cross-site requests (§65). \"priority\" is the HTTP extensible-priorities header (RFC 9218)."));
    missions.complete("network", "whoami");
  };
  return panel("What did we negotiate?", button("Ask the server", go, { class: "primary" }), out);
}

// ------------------------------------------------- waterfall ----
function waterfallPanel() {
  const out = h("div", {});
  const count = h("input", { type: "number", value: "18", min: "1", max: "40", "aria-label": "Number of requests", style: { width: "70px" } });
  const delay = h("input", { type: "number", value: "600", min: "50", max: "3000", step: "50", "aria-label": "Server delay (ms)", style: { width: "80px" } });

  const fire = async () => {
    fireBtn.disabled = true;
    const run = Math.random().toString(36).slice(2, 8);
    const n = Number(count.value);
    out.replaceChildren(h("p", {}, `Firing ${n} requests…`));
    const t0 = performance.now();
    await Promise.all(Array.from({ length: n }, (_, i) => fetch(`/api/delay?ms=${delay.value}&id=${i}&run=${run}`, { cache: "no-store" }).then((r) => r.json())));
    const wall = performance.now() - t0;
    await sleep(50); // let the resource timing entries land
    const entries = performance.getEntriesByType("resource").filter((e) => e.name.includes(`run=${run}`))
      .sort((a, b) => a.startTime - b.startTime);
    const start = Math.min(...entries.map((e) => e.startTime));
    const end = Math.max(...entries.map((e) => e.responseEnd));
    const pct = (t) => ((t - start) / (end - start)) * 100;
    const proto = entries[0]?.nextHopProtocol || "?";
    const conns = assignConnections(entries, proto);
    out.replaceChildren(
      h("div", { class: "readouts" },
        h("div", { class: "readout" }, h("small", {}, "Protocol"), h("b", {}, proto)),
        h("div", { class: "readout" }, h("small", {}, "Wall time"), h("b", {}, ms(wall))),
        h("div", { class: "readout" }, h("small", {}, "Connections used (approx)"), h("b", {}, String(conns.lanes))),
        h("div", { class: "readout" }, h("small", {}, "Server time each"), h("b", {}, ms(entries[0]?.serverTiming?.[0]?.duration)))),
      h("div", { class: "legend" },
        h("span", {}, h("i", { class: "seg-wait" }), "queued / stalled"),
        h("span", {}, h("i", { class: "seg-ttfb" }), "waiting for server (TTFB)"),
        h("span", {}, h("i", { class: "seg-dl" }), "download")),
      h("div", { class: "waterfall", role: "img", "aria-label": `Waterfall of ${entries.length} requests over ${proto}` },
        entries.map((e, i) => {
          const segs = [
            ["seg-wait", e.requestStart - e.startTime],
            ["seg-ttfb", e.responseStart - e.requestStart],
            ["seg-dl", e.responseEnd - e.responseStart],
          ];
          const total = e.responseEnd - e.startTime;
          return h("div", { class: "wf-row" },
            h("span", {}, `#${i} c${conns.of.get(e)}`),
            h("div", { class: "wf-track" },
              h("div", { class: "wf-bar", style: { left: `${pct(e.startTime)}%`, width: `${Math.max(0.5, pct(e.responseEnd) - pct(e.startTime))}%`, animationDelay: `${i * 40}ms` } },
                segs.map(([cls, d]) => h("span", { class: cls, style: { width: `${(Math.max(0, d) / total) * 100}%` } })))),
            h("span", {}, ms(total)));
        })),
      h("p", { class: "note" }, proto === "h2" || proto === "h3"
        ? "HTTP/2: every request shares one connection as interleaved streams, so they all start at once. The grey \"queued\" segments are gone."
        : "HTTP/1.1: a connection carries one request at a time and browsers open at most 6 per origin, so the requests run in waves. Open the https:// URL (npm run cert first) to compare with HTTP/2."));
    fireBtn.disabled = false;
    missions.complete("network", "waterfall");
  };
  const fireBtn = button("🚀 Fire the burst", fire, { class: "primary" });
  return panel("Experiment: many requests at once",
    h("p", {}, "Each request hits /api/delay, which waits on the server and then answers. Watch how the protocol schedules them."),
    h("div", { class: "row" }, h("label", {}, "requests ", count), h("label", {}, "server delay ms ", delay), fireBtn),
    out);
}

// Group requests into connections the simple way: a request reuses the
// first lane that was free when it started sending. HTTP/2 is always one lane.
function assignConnections(entries, proto) {
  const of = new Map();
  if (proto === "h2" || proto === "h3") {
    entries.forEach((e) => of.set(e, 1));
    return { lanes: 1, of };
  }
  const lanes = [];
  for (const e of entries) {
    let i = lanes.findIndex((free) => free <= e.requestStart + 1);
    if (i === -1) i = lanes.push(0) - 1;
    lanes[i] = e.responseEnd;
    of.set(e, i + 1);
  }
  return { lanes: lanes.length, of };
}

function ssePanel() {
  const log = h("ol", { class: "log", "aria-label": "Server-sent events" });
  let es;
  const start = () => {
    es?.close();
    log.replaceChildren();
    es = new EventSource("/api/ticks");
    es.onmessage = (e) => {
      const { tick, at } = JSON.parse(e.data);
      log.prepend(h("li", {}, `id=${e.lastEventId} tick ${tick} · server→browser ${Date.now() - at}ms`));
      if (tick === 10) es.close();
    };
    es.onerror = () => es.close();
  };
  return {
    el: panel("Streaming: Server-Sent Events", h("p", { class: "note" }, "One long-lived response; the server writes an event every 500ms. Over HTTP/2 it is just another stream on the shared connection, so it doesn't use up one of the 6 HTTP/1.1 connection slots."), button("Open the stream", start), log),
    close: () => es?.close(),
  };
}

export function mount() {
  const sse = ssePanel();
  const el = station(
    "URL to bytes: TLS & HTTP/2",
    "Before any JavaScript runs, the browser has to find the server, open a connection, secure it and ask for the page. These numbers come from your browser's own Resource Timing and Navigation Timing entries.",
    missionList("network"),
    panel("How this page arrived", navigationPhases()),
    waterfallPanel(),
    h("div", { class: "grid-2" }, tlsPanel(), whoamiPanel()),
    sse.el,
  );
  return { el, unmount: sse.close };
}
```

## public/js/stations/overview.js

```javascript
// Station 0: the browser's process architecture, and a guided tour of what
// happens between a click and the pixels changing (§18, §22, §79).
import { h, station, panel, button, sleep, ms, bytes } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";

const PROCS = [
  { id: "browser", title: "Browser process", threads: ["UI thread", "IO thread"], notes: ["Tabs, address bar, permissions", "Routes input to the right renderer"] },
  { id: "network", title: "Network service", threads: ["IO thread"], notes: ["DNS, TCP/QUIC, TLS", "HTTP cache, cookies"] },
  { id: "renderer", title: "Renderer process (one per site)", threads: ["Main thread", "Compositor thread", "Raster workers"], notes: ["V8 + Blink: JS, DOM, style, layout, paint", "Sandboxed: no direct disk or network"] },
  { id: "gpu", title: "GPU process", threads: ["Viz / display compositor"], notes: ["Draws composited layers", "Talks to the OS graphics stack"] },
];

const TOUR = [
  ["browser", "UI thread", "1. The OS delivers your click to the browser process. It hit-tests the tab and forwards the event over IPC to that tab's renderer."],
  ["renderer", "Compositor thread", "2. The renderer's compositor thread gets it first. If no listener is on that spot (or listeners are passive), it could scroll on its own without the main thread."],
  ["renderer", "Main thread", "3. There is a click listener, so the event is queued as a task for the main thread. It runs when the current task finishes. A long task here is what INP measures."],
  ["renderer", "Main thread", "4. Your handler runs in V8, then the microtask checkpoint, then (at the next frame) style → layout → paint record. Same thread, one after another."],
  ["renderer", "Raster workers", "5. Paint produced a list of draw commands. Raster threads turn tiles of it into bitmaps (or GPU textures)."],
  ["renderer", "Compositor thread", "6. The compositor thread assembles layers into a compositor frame. Transform/opacity animations live here, off the main thread."],
  ["gpu", "Viz / display compositor", "7. The GPU process draws the frame and hands it to the OS at the next vsync, usually within 16.7ms at 60Hz."],
];

function facts() {
  const nav = performance.getEntriesByType("navigation")[0];
  const conn = navigator.connection;
  const rows = [
    ["Engine (brands)", navigator.userAgentData?.brands?.map((b) => `${b.brand} ${b.version}`).join(", ") ?? navigator.userAgent.slice(0, 80)],
    ["Logical CPU cores", navigator.hardwareConcurrency],
    ["Device memory (approx)", navigator.deviceMemory ? `${navigator.deviceMemory} GB` : "not exposed"],
    ["Page protocol", nav?.nextHopProtocol || "unknown"],
    ["Secure context", String(isSecureContext)],
    ["Cross-origin isolated", String(crossOriginIsolated)],
    ["WebGPU", "gpu" in navigator ? "available" : "not available"],
    ["Network (estimate)", conn ? `${conn.effectiveType}, ~${conn.downlink} Mbps, rtt ~${conn.rtt}ms` : "not exposed"],
    ["Device pixel ratio", devicePixelRatio],
    ["Perf entry types", (PerformanceObserver.supportedEntryTypes ?? []).length],
  ];
  const dl = h("dl", { class: "kv" }, rows.map(([k, v]) => [h("dt", {}, k), h("dd", {}, String(v))]));
  navigator.storage?.estimate?.().then(({ usage, quota }) => dl.append(h("dt", {}, "Storage used / quota"), h("dd", {}, `${bytes(usage)} / ${bytes(quota)}`)));
  if (nav) dl.append(h("dt", {}, "This page: TTFB / DOMContentLoaded"), h("dd", {}, `${ms(nav.responseStart)} / ${ms(nav.domContentLoadedEventEnd)}`));
  return dl;
}

export function mount() {
  const boxes = PROCS.map((p) =>
    h("div", { class: "proc", "data-proc": p.id },
      h("h4", {}, p.title),
      h("ul", {}, p.notes.map((n) => h("li", {}, n))),
      p.threads.map((t) => h("div", { class: "thread", "data-thread": t }, "🧵 ", t))),
  );
  const caption = h("p", { class: "tour-caption", "aria-live": "polite" }, "Press play to follow one click through the browser.");
  let cancelled = false;

  const play = async () => {
    playBtn.disabled = true;
    for (const [proc, thread, text] of TOUR) {
      if (cancelled) return;
      boxes.forEach((b) => b.classList.toggle("lit", b.dataset.proc === proc));
      boxes.forEach((b) => b.querySelectorAll(".thread").forEach((t) => t.classList.toggle("lit", b.dataset.proc === proc && t.dataset.thread === thread)));
      caption.textContent = text;
      await sleep(2300);
    }
    boxes.forEach((b) => b.classList.remove("lit"));
    caption.textContent = "That round trip is usually a few milliseconds. When it isn't, the culprit is almost always step 3 or 4: the main thread was busy.";
    playBtn.disabled = false;
    missions.complete("overview", "tour");
  };
  const playBtn = button("▶ Play the tour of a click", play, { class: "primary" });

  const el = station(
    "The big picture",
    "A browser is a small operating system: several sandboxed processes, each with its own threads, talking over IPC. Every other station zooms into one box on this map.",
    missionList("overview"),
    panel("Processes and threads", h("div", { class: "arch" }, boxes), h("div", { class: "row" }, playBtn), caption),
    h("div", { class: "grid-2" },
      panel("Your browser, right now", facts(), h("p", { class: "note" }, "All read live from Web APIs. Cross-origin isolation is on because the server sends COOP + COEP headers.")),
      panel("Where to go next",
        h("ol", {},
          h("li", {}, h("a", { href: "#/network" }, "URL to bytes"), ": DNS, TLS 1.3, HTTP/1.1 vs HTTP/2, measured."),
          h("li", {}, h("a", { href: "#/render" }, "Rendering"), ": style → layout → paint → composite, and how to skip stages."),
          h("li", {}, h("a", { href: "#/eventloop" }, "Event loop"), ": step through tasks and microtasks, then check against the real thing."),
          h("li", {}, h("a", { href: "#/engine" }, "Inside V8"), ": bytecode, tiers, hidden classes, inline caches."),
          h("li", {}, h("a", { href: "#/threads" }, "Threads"), ": workers, transfer, SharedArrayBuffer races."),
          h("li", {}, h("a", { href: "#/storage" }, "Storage"), ": IndexedDB transactions you can watch."),
          h("li", {}, h("a", { href: "#/vitals" }, "Core Web Vitals"), ": break LCP/INP/CLS and fix them."),
          h("li", {}, h("a", { href: "#/debug" }, "Debugging"), ": errors, CSP reports, leaks, DevTools recipes.")))),
  );
  return { el, unmount: () => (cancelled = true) };
}
```

## public/js/stations/render.js

```javascript
// Station 2: the rendering pipeline (§21, §22, §81).
import { h, station, panel, button, nextFrame, ms } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";

const STAGES = [
  ["js", "JavaScript"],
  ["style", "Style"],
  ["layout", "Layout"],
  ["paint", "Paint"],
  ["composite", "Composite"],
];

// Which stages a change to each property triggers in Blink/WebKit/Gecko
// (with the element on its own compositor layer for transform/opacity).
const EXPERIMENTS = [
  { prop: "left", label: "left (geometry)", stages: ["js", "style", "layout", "paint", "composite"],
    apply: (el, t) => (el.style.left = `${8 + t * 380}px`),
    why: "Changing geometry means positions of this box (and possibly everything around it) must be recomputed, then repainted." },
  { prop: "background-color", label: "background-color (paint only)", stages: ["js", "style", "paint", "composite"],
    apply: (el, t) => (el.style.backgroundColor = `hsl(${Math.round(t * 300)} 80% 55%)`),
    why: "No geometry changed, so layout is skipped, but the pixels inside the box must be repainted every frame." },
  { prop: "transform", label: "transform (composite only)", stages: ["js", "style", "composite"],
    apply: (el, t) => (el.style.transform = `translateX(${t * 380}px) rotate(${t * 360}deg)`),
    why: "The box's layer is already painted. The compositor just moves the bitmap: no layout, no paint, and it keeps running even if the main thread is busy (for CSS animations)." },
];

function pipeline() {
  const pills = STAGES.map(([id, name]) => h("span", { class: "stage-pill", "data-stage": id }, name));
  return {
    el: h("div", { class: "pipeline", "aria-label": "Rendering stages" }, pills),
    show(active) {
      pills.forEach((p) => {
        const on = active.includes(p.dataset.stage);
        p.classList.toggle("on", on);
        p.classList.toggle("skipped", !on);
        if (on) {
          p.classList.remove("flash");
          void p.offsetWidth; // restart the CSS animation (a deliberate forced layout!)
          p.classList.add("flash");
        }
      });
    },
  };
}

function animationPanel() {
  const pipe = pipeline();
  const mover = h("div", { class: "mover", "aria-hidden": "true" });
  const why = h("p", { "aria-live": "polite" }, "Pick a property to animate.");
  const stats = h("p", { class: "note mono" });
  let running = false;

  const run = async (exp) => {
    if (running) return;
    running = true;
    mover.style.cssText = ""; // reset all three properties
    mover.style.willChange = exp.prop === "transform" ? "transform" : "auto"; // promote to its own layer
    pipe.show(exp.stages);
    why.textContent = exp.why;
    const frames = [];
    const start = performance.now();
    let last = start;
    for (;;) {
      const now = await nextFrame();
      frames.push(now - last);
      last = now;
      const t = Math.min(1, (now - start) / 1500);
      exp.apply(mover, t);
      if (t === 1) break;
    }
    const avg = frames.reduce((a, b) => a + b, 0) / frames.length;
    stats.textContent = `${frames.length} frames · avg ${avg.toFixed(1)}ms per frame · worst ${Math.max(...frames).toFixed(1)}ms`;
    running = false;
    if (exp.prop === "transform") missions.complete("render", "composite");
  };

  return panel("Animate a property, watch which stages run",
    pipe.el,
    h("div", { class: "row" }, EXPERIMENTS.map((e) => button(e.label, () => run(e)))),
    h("div", { class: "track" }, mover),
    why, stats,
    h("p", { class: "note" }, "See it for real: DevTools → More tools → Rendering → tick \"Paint flashing\" and \"Layer borders\", then run each animation. Green flashes are repaints."));
}

function thrashPanel() {
  const N = 400;
  const bars = h("div", { class: "bars", "aria-hidden": "true" }, Array.from({ length: N }, () => h("div")));
  const result = h("div", { class: "readouts" });
  const times = {};

  const reset = () => [...bars.children].forEach((b) => (b.style.width = "100px"));
  const show = () => {
    result.replaceChildren(
      ...Object.entries(times).map(([k, v]) => h("div", { class: "readout" }, h("small", {}, k), h("b", {}, ms(v)))),
    );
    if (times["interleaved (thrash)"] != null && times["batched reads, then writes"] != null) missions.complete("render", "thrash");
  };

  const thrash = () => {
    reset();
    const t0 = performance.now();
    for (const b of bars.children) {
      const w = b.offsetWidth; // READ: forces layout, because the previous write dirtied it
      b.style.width = `${w + 1}px`; // WRITE: dirties layout again
    }
    void bars.offsetWidth;
    times["interleaved (thrash)"] = performance.now() - t0;
    show();
  };
  const batched = () => {
    reset();
    const t0 = performance.now();
    const widths = [...bars.children].map((b) => b.offsetWidth); // all reads: ONE layout
    [...bars.children].forEach((b, i) => (b.style.width = `${widths[i] + 1}px`)); // all writes
    void bars.offsetWidth; // the one layout the browser would do anyway
    times["batched reads, then writes"] = performance.now() - t0;
    show();
  };

  return panel("Experiment: forced synchronous layout (layout thrashing)",
    h("p", {}, `${N} bars. Each run reads every bar's width and makes it 1px wider. The only difference is the order of reads and writes.`),
    h("div", { class: "row" }, button("Interleave read/write", thrash, { class: "danger" }), button("Batch reads, then writes", batched, { class: "primary" })),
    result, bars,
    h("pre", { class: "code" }, h("span", {}, "// thrash: N layouts"), h("span", {}, "for (const b of bars) { b.style.width = b.offsetWidth + 1 + 'px'; }"),
      h("span", {}, "// batched: 1 layout"), h("span", {}, "const ws = bars.map(b => b.offsetWidth);"), h("span", {}, "bars.forEach((b, i) => b.style.width = ws[i] + 1 + 'px');")));
}

function renderTreePanel() {
  const sample = h("div", { class: "sample", id: "rt-sample" },
    h("h4", {}, "Title"),
    h("p", { hidden: true }, "display:none paragraph"),
    h("span", { style: { visibility: "hidden" } }, "visibility:hidden span"),
    h("span", {}, "visible span"));
  const out = h("div", { class: "grid-2" });
  const describe = () => {
    const domItems = [];
    const renderItems = [];
    const walk = (el, depth) => {
      const cs = getComputedStyle(el);
      const name = `${el.tagName.toLowerCase()}${el.hidden ? " [hidden]" : ""}`;
      domItems.push(h("li", { style: { marginLeft: `${depth * 14}px` } }, name));
      if (cs.display === "none") {
        renderItems.push(h("li", { class: "gone", style: { marginLeft: `${depth * 14}px` } }, `${name} (display:none → no box)`));
        return; // and none of its children get boxes either
      }
      const before = getComputedStyle(el, "::before").content;
      renderItems.push(h("li", { style: { marginLeft: `${depth * 14}px` } }, `${name}${cs.visibility === "hidden" ? " (box, but invisible)" : ""}`));
      if (before && before !== "none" && before !== "normal") renderItems.push(h("li", { class: "pseudo", style: { marginLeft: `${(depth + 1) * 14}px` } }, `::before ${before}`));
      for (const c of el.children) walk(c, depth + 1);
    };
    walk(sample, 0);
    out.replaceChildren(
      h("div", {}, h("h4", {}, "DOM tree"), h("ul", { class: "tree" }, domItems)),
      h("div", {}, h("h4", {}, "Render tree (layout boxes)"), h("ul", { class: "tree" }, renderItems)));
  };
  // getComputedStyle() on a node that isn't in the document returns empty
  // values, so wait until the station has actually been attached.
  const whenConnected = () => (sample.isConnected ? describe() : requestAnimationFrame(whenConnected));
  requestAnimationFrame(whenConnected);
  return panel("DOM tree vs render tree", h("p", {}, "The render tree only has nodes that produce boxes: display:none disappears (with all its children), visibility:hidden keeps its box, and ::before/::after appear even though they are not in the DOM. Computed live with getComputedStyle()."), sample, out);
}

export function mount() {
  return {
    el: station(
      "The rendering pipeline",
      "Every frame, the main thread may run JavaScript, recompute styles, lay out boxes and record paint commands. Then the compositor thread assembles layers on the GPU. Fast UIs skip as many stages as they can.",
      missionList("render"),
      animationPanel(),
      thrashPanel(),
      renderTreePanel(),
    ),
  };
}
```

## public/js/stations/storage.js

```javascript
// Station 6: IndexedDB and friends (§43–§45, §82). Every request and
// transaction event is logged as it fires, so you can see the lifecycle.
import { h, station, panel, button, sleep, ms, bytes } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";

const DB_NAME = "browser-lab";
const t0 = performance.now();

function openDB(log) {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, 1);
    log("indexedDB.open('browser-lab', 1) → request created (nothing has happened yet)", "req");
    req.onupgradeneeded = () => {
      log("upgradeneeded: version 0 → 1, creating object store 'notes' + index 'byCreated'", "evt");
      const store = req.result.createObjectStore("notes", { keyPath: "id", autoIncrement: true });
      store.createIndex("byCreated", "created");
    };
    req.onsuccess = () => {
      log("open success: we have a connection", "ok");
      resolve(req.result);
    };
    req.onerror = () => reject(req.error);
    req.onblocked = () => log("blocked: another tab holds an older version open", "bad");
  });
}

// Wrap a request in a promise, logging its success.
const done = (req, log, label) =>
  new Promise((resolve, reject) => {
    req.onsuccess = () => {
      log(`${label} → success`, "ok");
      resolve(req.result);
    };
    req.onerror = () => reject(req.error);
  });

export function mount() {
  const logEl = h("ol", { class: "log", "aria-label": "IndexedDB event log" });
  const log = (text, kind = "") => logEl.prepend(h("li", { "data-kind": kind }, `+${ms(performance.now() - t0)}  ${text}`));
  const list = h("ul", { class: "tree", "aria-label": "Stored notes" });
  const title = h("input", { "aria-label": "Note text", value: "Hello from IndexedDB" });
  let dbp = openDB(log);

  const refresh = async () => {
    const db = await dbp;
    const tx = db.transaction("notes", "readonly");
    const items = [];
    await new Promise((resolve) => {
      tx.objectStore("notes").index("byCreated").openCursor(null, "prev").onsuccess = (e) => {
        const cursor = e.target.result;
        if (cursor) {
          items.push(cursor.value);
          cursor.continue(); // each step is another async success event
        } else resolve();
      };
    });
    list.replaceChildren(...items.map((n) => h("li", {}, `#${n.id} · ${n.text} · ${new Date(n.created).toLocaleTimeString()}`)));
    return items;
  };

  const add = async () => {
    const db = await dbp;
    const tx = db.transaction("notes", "readwrite");
    log("transaction('notes', 'readwrite') → active", "evt");
    tx.oncomplete = () => log("transaction complete: the write is durable", "ok");
    const id = await done(tx.objectStore("notes").add({ text: title.value, created: Date.now() }), log, "store.add()");
    await new Promise((r) => tx.addEventListener("complete", r));
    const items = await refresh();
    if (items.some((n) => n.id === id)) missions.complete("storage", "idb");
  };

  const clear = async () => {
    const db = await dbp;
    const tx = db.transaction("notes", "readwrite");
    await done(tx.objectStore("notes").clear(), log, "store.clear()");
    refresh();
  };

  // The classic bug: awaiting something that is not an IDB request inside a
  // transaction. The transaction auto-commits when it has no pending
  // requests at the end of a task, so the next request throws.
  const inactive = async () => {
    const db = await dbp;
    const tx = db.transaction("notes", "readwrite");
    const store = tx.objectStore("notes");
    log("transaction started; now awaiting a 10ms timer inside it…", "evt");
    tx.oncomplete = () => log("…the transaction already auto-committed (no pending requests)", "evt");
    await sleep(10);
    try {
      store.add({ text: "never written", created: Date.now() });
    } catch (err) {
      log(`${err.name}: ${err.message}`, "bad");
      missions.complete("storage", "inactive");
    }
  };

  const blocking = async () => {
    const big = "x".repeat(2 * 1024 * 1024);
    let t = performance.now();
    try {
      localStorage.setItem("browser-lab:big", big);
    } catch (err) {
      log(`localStorage: ${err.name}`, "bad");
    }
    const lsMs = performance.now() - t;
    localStorage.removeItem("browser-lab:big");
    const db = await dbp;
    t = performance.now();
    const tx = db.transaction("notes", "readwrite");
    tx.objectStore("notes").put({ id: -1, text: big, created: 0 });
    const idbMainMs = performance.now() - t; // time the main thread was busy issuing the write
    await new Promise((r) => (tx.oncomplete = r));
    const idbTotal = performance.now() - t;
    const tx2 = db.transaction("notes", "readwrite");
    tx2.objectStore("notes").delete(-1);
    log(`2 MB write: localStorage blocked the main thread ${ms(lsMs)}; IndexedDB blocked it ${ms(idbMainMs)} (finished ${ms(idbTotal)} later, off-thread)`, "ok");
  };

  const quota = h("dl", { class: "kv" });
  (async () => {
    const est = await navigator.storage?.estimate?.();
    const persisted = await navigator.storage?.persisted?.();
    quota.append(
      h("dt", {}, "Origin usage"), h("dd", {}, bytes(est?.usage)),
      h("dt", {}, "Origin quota"), h("dd", {}, bytes(est?.quota)),
      h("dt", {}, "Persistent?"), h("dd", {}, String(persisted)),
    );
  })();

  refresh();
  return {
    el: station(
      "IndexedDB & storage",
      "IndexedDB is a transactional, asynchronous object database in every browser. Requests fire events; transactions commit on their own when nothing is pending. Watch it happen.",
      missionList("storage"),
      h("div", { class: "grid-2" },
        panel("Notes in IndexedDB",
          h("div", { class: "row" }, title, button("Add note", add, { class: "primary" }), button("Clear", clear)),
          list),
        panel("Event log (newest first)", logEl)),
      h("div", { class: "grid-2" },
        panel("Pitfall: the transaction that commits behind your back",
          h("p", {}, "await a timer or fetch() inside a transaction and it auto-commits. Your next request throws TransactionInactiveError."),
          button("Trigger it", inactive, { class: "danger" })),
        panel("Sync vs async storage",
          h("p", {}, "localStorage is synchronous: a big write blocks the main thread. IndexedDB does its I/O off the main thread."),
          button("Write 2 MB to each", blocking), quota)),
      panel("Inspect it", h("p", { class: "note" }, "DevTools → Application → IndexedDB → browser-lab → notes. Right-click → Refresh after adding notes. Storage → \"Clear site data\" wipes everything for this origin.")),
    ),
    unmount: async () => (await dbp).close(),
  };
}
```

## public/js/stations/threads.js

```javascript
// Station 5: workers and threads (§37–§42, §82).
import { h, station, panel, button, ms, bytes } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";
import { countPrimes } from "../lib/primes.js";

const worker = (name) => new Worker(new URL(`../workers/${name}.js`, import.meta.url), { type: "module" });

// A bouncing ball drawn every frame, plus a frame-time graph. If the main
// thread is blocked, both freeze: that is jank, made visible.
function ballPanel() {
  const canvas = h("canvas", { width: 720, height: 200, "aria-label": "Animated ball and frame-time graph", role: "img" });
  const ctx = canvas.getContext("2d");
  const fps = h("b", {}, "—");
  const worst = h("b", {}, "—");
  const frames = [];
  let x = 40, y = 60, vx = 4, vy = 3, last = performance.now(), raf = 0, worstGap = 0;
  const color = getComputedStyle(document.documentElement).getPropertyValue("--accent").trim() || "#3355ff";
  const tick = (now) => {
    const dt = now - last;
    last = now;
    frames.push(dt);
    if (frames.length > 180) frames.shift();
    worstGap = Math.max(worstGap, dt);
    x += vx * (dt / 16.7);
    y += vy * (dt / 16.7);
    if (x < 14 || x > 706) vx = -vx, (x = Math.max(14, Math.min(706, x)));
    if (y < 14 || y > 110) vy = -vy, (y = Math.max(14, Math.min(110, y)));
    ctx.clearRect(0, 0, 720, 200);
    ctx.fillStyle = color;
    ctx.beginPath();
    ctx.arc(x, y, 12, 0, Math.PI * 2);
    ctx.fill();
    // frame-time graph: one bar per frame, red when over the 16.7ms budget
    frames.forEach((f, i) => {
      ctx.fillStyle = f > 50 ? "#e53935" : f > 17.5 ? "#f5a623" : "#2eb872";
      const hgt = Math.min(70, f);
      ctx.fillRect(i * 4, 200 - hgt, 3, hgt);
    });
    ctx.fillStyle = "#888";
    ctx.fillRect(0, 200 - 16.7, 720, 1);
    fps.textContent = `${Math.round(1000 / (frames.slice(-30).reduce((a, b) => a + b, 0) / Math.min(30, frames.length)))}`;
    worst.textContent = ms(worstGap);
    raf = requestAnimationFrame(tick);
  };
  raf = requestAnimationFrame(tick);
  return {
    el: h("div", {}, canvas, h("div", { class: "readouts" },
      h("div", { class: "readout" }, h("small", {}, "frames / second"), fps),
      h("div", { class: "readout" }, h("small", {}, "worst frame gap"), worst))),
    resetWorst: () => (worstGap = 0),
    worst: () => worstGap,
    stop: () => cancelAnimationFrame(raf),
  };
}

function jankPanel(ball) {
  const limit = h("select", { "aria-label": "How much work" },
    h("option", { value: "3000000" }, "primes below 3M (~0.3s)"),
    h("option", { value: "6000000", selected: true }, "primes below 6M (~1s)"),
    h("option", { value: "12000000" }, "primes below 12M (~3s)"));
  const out = h("ol", { class: "log", "aria-label": "Results" });
  const lanes = h("div", { class: "lanes" },
    h("div", { class: "lane" }, h("span", {}, "Main thread"), h("div", { class: "rail", "data-lane": "main" })),
    h("div", { class: "lane" }, h("span", {}, "Worker thread"), h("div", { class: "rail", "data-lane": "worker" })));
  const mark = (lane, startPct, widthPct) => {
    const rail = lanes.querySelector(`[data-lane="${lane}"]`);
    rail.replaceChildren(h("div", { class: `work ${lane === "worker" ? "worker" : ""}`, style: { left: `${startPct}%`, width: `${widthPct}%` } }));
  };
  let blockedOnMain = false;
  const onMain = () => {
    ball.resetWorst();
    out.prepend(h("li", {}, "main thread: computing… (the ball will freeze)"));
    mark("main", 5, 90);
    lanes.querySelector('[data-lane="worker"]').replaceChildren();
    // Let the message paint first, then block.
    setTimeout(() => {
      const t0 = performance.now();
      const count = countPrimes(Number(limit.value));
      const took = performance.now() - t0;
      out.prepend(h("li", {}, `main thread: ${count.toLocaleString()} primes in ${ms(took)}. The page could not paint or respond for that long.`));
      requestAnimationFrame(() => requestAnimationFrame(() => {
        blockedOnMain = ball.worst() > 100;
        out.prepend(h("li", {}, `worst frame gap: ${ms(ball.worst())}`));
      }));
    }, 50);
  };
  const inWorker = () => {
    ball.resetWorst();
    const w = worker("primes");
    mark("worker", 5, 90);
    mark("main", 5, 1);
    const t0 = performance.now();
    w.postMessage({ limit: Number(limit.value) });
    w.onmessage = ({ data }) => {
      out.prepend(h("li", {}, `worker: ${data.count.toLocaleString()} primes in ${ms(data.ms)} (round trip ${ms(performance.now() - t0)}), worst frame gap ${ms(ball.worst())}. The ball never stopped.`));
      w.terminate();
      if (blockedOnMain && ball.worst() < 100) missions.complete("threads", "jank");
    };
  };
  return panel("Experiment: the same work on two threads",
    h("div", { class: "row" }, limit, button("Run on the main thread", onMain, { class: "danger" }), button("Run in a Web Worker", inWorker, { class: "primary" })),
    lanes, out);
}

function transferPanel() {
  const out = h("ol", { class: "log", "aria-label": "Transfer results" });
  const run = async (transfer) => {
    const w = worker("echo");
    const buffer = new ArrayBuffer(64 * 1024 * 1024);
    const t0 = performance.now();
    w.postMessage({ buffer }, transfer ? [buffer] : []);
    const posted = performance.now() - t0;
    const size = await new Promise((r) => (w.onmessage = (e) => r(e.data)));
    out.prepend(h("li", {},
      `${transfer ? "transfer" : "copy (structured clone)"}: postMessage took ${ms(posted)} on the main thread; worker got ${bytes(size)}; sender's buffer is now ${bytes(buffer.byteLength)}`));
    w.terminate();
  };
  return panel("Copy vs transfer a 64 MB buffer",
    h("p", {}, "postMessage copies by default. Listing the buffer as transferable moves ownership instead: no copy, and the sender's buffer is detached (length 0)."),
    h("div", { class: "row" }, button("Copy it", () => run(false)), button("Transfer it", () => run(true), { class: "primary" })), out);
}

function racePanel() {
  const out = h("ol", { class: "log", "aria-label": "Race results" });
  const run = async (atomic) => {
    if (!crossOriginIsolated) {
      out.prepend(h("li", {}, "SharedArrayBuffer needs cross-origin isolation (COOP + COEP headers). Serve the lab with npm start."));
      return;
    }
    const sab = new SharedArrayBuffer(4);
    const N = 1_000_000;
    const workers = Array.from({ length: 4 }, () => worker("counter"));
    const t0 = performance.now();
    await Promise.all(workers.map((w) => new Promise((r) => ((w.onmessage = r), w.postMessage({ sab, n: N, atomic })))));
    const value = new Int32Array(sab)[0];
    workers.forEach((w) => w.terminate());
    const lost = 4 * N - value;
    out.prepend(h("li", {}, `${atomic ? "Atomics.add" : "counter[0]++"}: expected ${(4 * N).toLocaleString()}, got ${value.toLocaleString()} → ${lost.toLocaleString()} updates lost (${ms(performance.now() - t0)})`));
    raceState[atomic ? "atomicOk" : "lost"] ||= atomic ? lost === 0 : lost > 0;
    if (raceState.lost && raceState.atomicOk) missions.complete("threads", "race");
  };
  const raceState = {};
  return panel("Experiment: a data race on shared memory",
    h("p", {}, "Four workers each add 1 to the same counter a million times, through one SharedArrayBuffer. ", h("code", {}, "counter[0]++"), " is load, add, store: two threads can load the same value and one increment vanishes."),
    h("p", { class: "note" }, `crossOriginIsolated = ${crossOriginIsolated}`),
    h("div", { class: "row" }, button("Race with counter[0]++", () => run(false), { class: "danger" }), button("Fix with Atomics.add", () => run(true), { class: "primary" })), out);
}

export function mount() {
  const ball = ballPanel();
  return {
    el: station(
      "Workers & threads",
      "Your JavaScript, style, layout and paint share one main thread. Any long computation there freezes everything. Workers are real OS threads with their own event loop; they talk to the page by message passing or, with care, shared memory.",
      missionList("threads"),
      panel("The main thread's heartbeat", ball.el, h("p", { class: "note" }, "Green bars fit the 16.7ms frame budget, orange ones miss it, red ones are long frames.")),
      jankPanel(ball),
      h("div", { class: "grid-2" }, transferPanel(), racePanel()),
    ),
    unmount: ball.stop,
  };
}
```

## public/js/stations/vitals.js

```javascript
// Station 7: Core Web Vitals (§35, §83). The numbers are collected from page
// load by lib/vitals.js; this station shows them and lets you break them.
import { h, station, panel, button, blockFor, ms } from "../lib/dom.js";
import { missions, missionList } from "../lib/missions.js";
import { vitals, onVitals, rating, THRESHOLDS, loafs, support, resetINP } from "../lib/vitals.js";

const fmt = (name, v) => (v == null ? "—" : name === "CLS" ? v.toFixed(3) : ms(v));

function gauge(name, desc) {
  const value = h("b", {}, "—");
  const needle = h("span", { class: "needle" });
  const el = h("div", { class: "gauge", "data-metric": name }, h("small", {}, name), value, h("div", { class: "meter" }, needle), h("small", {}, desc));
  return {
    el,
    update(v) {
      value.textContent = fmt(name, v);
      const [good, poor] = THRESHOLDS[name];
      // place the needle: 0–33% good, 33–66% needs improvement, 66–100% poor
      const pos = v == null ? 0 : v <= good ? (v / good) * 33 : v <= poor ? 33 + ((v - good) / (poor - good)) * 33 : Math.min(98, 66 + ((v - poor) / poor) * 33);
      needle.style.left = `${pos}%`;
      el.className = `gauge ${rating(name, v)}`;
    },
  };
}

// Give the browser a chance to paint between chunks of work.
const yieldToMain = () =>
  globalThis.scheduler?.yield ? scheduler.yield() : new Promise((r) => setTimeout(r, 0));

export function mount() {
  const gauges = {
    LCP: gauge("LCP", "Largest Contentful Paint: when the main content showed up"),
    INP: gauge("INP", "Interaction to Next Paint: worst input → paint delay"),
    CLS: gauge("CLS", "Cumulative Layout Shift: how much things jumped"),
    FCP: gauge("FCP", "First Contentful Paint (diagnostic)"),
    TTFB: gauge("TTFB", "Time to First Byte (diagnostic)"),
  };
  const lastEl = h("p", { class: "note", "aria-live": "polite" });
  const loafBody = h("tbody");
  const slot = h("div", {}); // layout shifts get injected here
  let sawPoorINP = false;

  const render = () => {
    for (const [k, g] of Object.entries(gauges)) g.update(vitals[k]);
    lastEl.textContent = vitals.last ? `Last interaction: ${vitals.last.name} on "${vitals.last.target}" took ${ms(vitals.last.duration)} to next paint · ${vitals.interactions} interactions so far` : "No interactions measured yet.";
    loafBody.replaceChildren(...loafs.slice(0, 8).map((l) =>
      h("tr", {}, h("td", {}, ms(l.start)), h("td", {}, ms(l.duration)), h("td", {}, ms(l.blocking)),
        h("td", { class: "mono" }, l.scripts.map((s) => `${s.invoker} (${s.source || "?"}) ${ms(s.duration)}`).join("; ") || "—"))));
    if (vitals.CLS > 0.1) missions.complete("vitals", "cls");
    if (vitals.last && vitals.last.duration > 500) sawPoorINP = true;
  };
  const off = onVitals(render);

  const status = h("span", { class: "note", "aria-live": "polite" });
  const shift = () => {
    status.textContent = "An \"ad\" will load in 800ms…";
    // Shifts within 500ms of an input don't count toward CLS, so this
    // simulates what really hurts: content that arrives late, on its own,
    // and then resizes itself. Both shifts land in one session window.
    setTimeout(() => {
      const ad = h("div", { class: "shift-banner" }, "📢 A late-loading ad with no space reserved for it");
      slot.prepend(ad);
      status.textContent = "Shifted! Everything below the ad jumped down…";
      setTimeout(() => {
        ad.classList.add("expanded");
        ad.append(h("div", {}, "…and then the ad resized itself. Shifted again."));
        status.textContent = "Shifted twice, 400ms apart: one session window, so the two scores add up.";
      }, 400);
    }, 800);
  };
  const fixShift = () => {
    slot.replaceChildren();
    status.textContent = "Banners removed. The fix in real life: reserve space (min-height or aspect-ratio) before content arrives.";
  };

  const slowClick = (e) => {
    e.currentTarget.textContent = "Working… (350ms blocking)";
    blockFor(350); // the handler hogs the main thread: no paint until it ends
    e.currentTarget.textContent = "Slow handler (350ms)";
  };
  const yieldingClick = async (e) => {
    const btn = e.currentTarget;
    btn.textContent = "Working in chunks…";
    for (let i = 0; i < 35; i++) {
      blockFor(10); // the same 350ms of work…
      await yieldToMain(); // …but the browser can paint after the first chunk
    }
    btn.textContent = "Same work, yielding";
    // The next paint came right after the first 10ms chunk, so this
    // interaction's INP contribution is small even though the work is the same.
    setTimeout(() => {
      if (sawPoorINP && vitals.last && vitals.last.duration < 200) missions.complete("vitals", "inp");
    }, 300);
  };

  render();
  const el = station(
    "Core Web Vitals",
    "Google's three user-centric metrics: loading (LCP), responsiveness (INP) and visual stability (CLS). Good means p75 of real users under 2.5s, 200ms and 0.1. These gauges measure this visit, with the same algorithms as the web-vitals library.",
    missionList("vitals"),
    slot, // late banners land here, above everything: shifts only count when visible
    h("div", { class: "gauges" }, Object.values(gauges).map((g) => g.el)),
    panel("Break CLS",
      h("div", { class: "row" }, button("Load a late banner", shift, { class: "danger" }), button("Remove banners", fixShift), status),
      h("p", {}, "Each banner lands at the top of the page, pushing everything below it down. CLS scores how much of the viewport moved and how far. Content that moves off-screen doesn't count, which is why real ads at the top hurt most."),
      !support.cls ? h("p", { class: "bad" }, "This browser doesn't report layout shifts.") : ""),
    panel("Break INP, then fix it",
      h("p", {}, "Both buttons do 350ms of work. The first does it in one block, so the click can't paint for 350ms. The second splits it into 10ms chunks and yields between them with ", h("code", {}, "scheduler.yield()"), " (or setTimeout as a fallback)."),
      h("div", { class: "row" },
        button("Slow handler (350ms)", slowClick, { class: "danger" }),
        button("Same work, yielding", yieldingClick, { class: "primary" }),
        button("reset INP", () => (resetINP(), render()), { class: "link" })),
      lastEl,
      !support.inp ? h("p", { class: "bad" }, "This browser doesn't expose Event Timing, so INP can't be measured here.") : ""),
    panel("Long animation frames: which script made the frame slow",
      support.loaf
        ? h("table", {}, h("thead", {}, h("tr", {}, h("th", {}, "start"), h("th", {}, "duration"), h("th", {}, "blocking"), h("th", {}, "scripts"))), loafBody)
        : h("p", {}, "This browser doesn't support the Long Animation Frames API (Chromium 123+)."),
      h("p", { class: "note" }, "In production, send these with the vitals so you know which handler to fix, not just that INP is bad.")),
    panel("Report vitals from real users",
      h("pre", { class: "code" }, `import { onLCP, onINP, onCLS } from "web-vitals/attribution";

function send(metric) {
  const body = JSON.stringify({ name: metric.name, value: metric.value,
    rating: metric.rating, id: metric.id, page: location.pathname,
    debug: metric.attribution });          // which element / script / phase
  // sendBeacon survives the page being closed; fetch keepalive is the fallback
  navigator.sendBeacon?.("/rum", body) ||
    fetch("/rum", { body, method: "POST", keepalive: true });
}
onLCP(send); onINP(send); onCLS(send);`.split("\n").map((l) => h("span", {}, l)))),
  );
  return { el, unmount: off };
}
```

## public/js/workers/counter.js

```javascript
// Increments a counter in shared memory n times (§39).
// atomic=false: load, add, store — three steps another thread can interleave with.
// atomic=true:  Atomics.add — one indivisible read-modify-write.
self.onmessage = ({ data: { sab, n, atomic } }) => {
  const counter = new Int32Array(sab);
  if (atomic) for (let i = 0; i < n; i++) Atomics.add(counter, 0, 1);
  else for (let i = 0; i < n; i++) counter[0]++;
  self.postMessage("done");
};
```

## public/js/workers/echo.js

```javascript
// Receives a buffer and reports its size: used to compare copy vs transfer.
self.onmessage = ({ data }) => self.postMessage(data.buffer.byteLength);
```

## public/js/workers/primes.js

```javascript
// A module worker (§38): same function, different thread.
import { countPrimes } from "../lib/primes.js";

self.onmessage = ({ data: { limit } }) => {
  const t0 = performance.now();
  const count = countPrimes(limit);
  self.postMessage({ count, ms: performance.now() - t0 });
};
```

## playwright.config.js

```javascript
// The lab's own E2E suite (§84): every station, over HTTP/1.1 and HTTP/2.
import { defineConfig, devices } from "@playwright/test";

const HTTP = 18080;
const HTTPS = 18443;

export default defineConfig({
  testDir: "tests",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? [["github"], ["html", { open: "never" }]] : [["list"]],
  use: { trace: "on-first-retry", ...devices["Desktop Chrome"], viewport: { width: 1360, height: 900 } },
  webServer: {
    command: "node scripts/cert.mjs && node server.mjs",
    url: `http://localhost:${HTTP}/`,
    env: { HTTP_PORT: String(HTTP), HTTPS_PORT: String(HTTPS) },
    reuseExistingServer: false,
  },
  projects: [
    { name: "http1", use: { baseURL: `http://localhost:${HTTP}` } },
    // Self-signed cert: tell the browser to accept it (tests only!).
    { name: "h2", use: { baseURL: `https://localhost:${HTTPS}`, ignoreHTTPSErrors: true } },
  ],
});
```

## tests/lab.spec.js

```javascript
import { test, expect } from "@playwright/test";

const STATIONS = ["overview", "network", "render", "eventloop", "engine", "threads", "storage", "vitals", "debug"];

// Fail any test on an error the page did not throw on purpose.
test.beforeEach(async ({ page }) => {
  page.on("pageerror", (err) => {
    if (!/Boom from a timer|Nobody handled me/.test(err.message)) throw err;
  });
});

const open = async (page, id) => {
  await page.goto(`/#/${id}`);
  await expect(page.locator(".station h2")).toBeVisible();
};

for (const id of STATIONS) {
  test(`station ${id} renders`, async ({ page }) => {
    await open(page, id);
    await expect(page.getByRole("link", { name: new RegExp(id === "eventloop" ? "event loop" : "", "i") }).first()).toBeVisible();
    await expect(page.getByRole("list", { name: "Missions" })).toBeVisible();
  });
}

test("the server is cross-origin isolated", async ({ page }) => {
  await open(page, "overview");
  expect(await page.evaluate(() => crossOriginIsolated)).toBe(true);
});

test("the event-loop model matches the real engine for every snippet", async ({ page }) => {
  await open(page, "eventloop");
  for (const tab of await page.getByRole("tab").all()) {
    await tab.click();
    await page.getByRole("button", { name: "Run it for real" }).click();
    await expect(page.locator("[data-match]")).toHaveAttribute("data-match", "true");
  }
  await expect(page.locator('[data-mission="verify"]')).toHaveClass(/done/);
});

test("stepping through a snippet fills the queues", async ({ page }) => {
  await open(page, "eventloop");
  const next = page.getByRole("button", { name: "Step ▶" });
  for (let i = 0; i < 4; i++) await next.click();
  await expect(page.locator(".box.webapis .token")).toHaveCount(1);
  await expect(page.locator(".box.micro .token")).toHaveCount(1);
});

test("request burst: protocol decides how many connections", async ({ page }, testInfo) => {
  await open(page, "network");
  await page.getByLabel("Number of requests").fill("12");
  await page.getByLabel("Server delay (ms)").fill("300");
  await page.getByRole("button", { name: "Fire the burst" }).click();
  const readout = (label) => page.locator(".readout", { hasText: label }).locator("b");
  if (testInfo.project.name === "h2") {
    await expect(readout("Protocol")).toHaveText("h2");
    await expect(readout("Connections used")).toHaveText("1");
  } else {
    await expect(readout("Protocol")).toHaveText("http/1.1");
    // 6 connections per origin, give or take one the page already had open
    await expect(readout("Connections used")).toHaveText(/^[2-7]$/);
  }
  await expect(page.locator(".wf-row")).toHaveCount(12);
});

test("whoami reports TLS 1.3 over HTTP/2", async ({ page }, testInfo) => {
  test.skip(testInfo.project.name !== "h2", "TLS only on the https project");
  await open(page, "network");
  await page.getByRole("button", { name: "Ask the server" }).click();
  await expect(page.locator(".kv dd").first()).toHaveText("2.0");
  await expect(page.getByText("TLSv1.3")).toBeVisible();
});

test("layout thrashing is slower than batching", async ({ page }) => {
  await open(page, "render");
  await page.getByRole("button", { name: "Interleave read/write" }).click();
  await page.getByRole("button", { name: "Batch reads, then writes" }).click();
  const value = async (label) => parseFloat(await page.locator(".readout", { hasText: label }).locator("b").innerText());
  expect(await value("interleaved")).toBeGreaterThan(await value("batched"));
});

test("the render tree drops display:none but keeps ::before", async ({ page }) => {
  await open(page, "render");
  await expect(page.locator(".tree .gone")).toContainText("p [hidden] (display:none → no box)");
  await expect(page.locator(".tree .pseudo")).toContainText("::before");
});

test("Atomics fixes the data race", async ({ page }) => {
  await open(page, "threads");
  await page.getByRole("button", { name: "Fix with Atomics.add" }).click();
  await expect(page.getByText(/Atomics\.add: expected 4,000,000, got 4,000,000 → 0 updates lost/)).toBeVisible({ timeout: 15000 });
});

test("a worker keeps the page responsive", async ({ page }) => {
  await open(page, "threads");
  await page.getByLabel("How much work").selectOption("3000000");
  await page.getByRole("button", { name: "Run in a Web Worker" }).click();
  await expect(page.getByText(/The ball never stopped/)).toBeVisible({ timeout: 15000 });
});

test("IndexedDB: write, read back, and the inactive-transaction pitfall", async ({ page }) => {
  await open(page, "storage");
  await page.getByLabel("Note text").fill("written by Playwright");
  await page.getByRole("button", { name: "Add note" }).click();
  await expect(page.getByRole("list", { name: "Stored notes" })).toContainText("written by Playwright");
  await page.getByRole("button", { name: "Trigger it" }).click();
  await expect(page.getByRole("list", { name: "IndexedDB event log" })).toContainText("TransactionInactiveError");
});

test("a late banner pushes CLS past 0.1", async ({ page }) => {
  await open(page, "vitals");
  await page.getByRole("button", { name: "Load a late banner" }).click();
  await expect(page.locator(".shift-banner")).toHaveCount(1);
  await expect(page.locator('[data-metric="CLS"]')).toHaveClass(/needs-improvement|poor/);
  await expect(page.locator('[data-mission="cls"]')).toHaveClass(/done/);
});

test("a 350ms click handler shows up as a slow interaction", async ({ page }) => {
  await open(page, "vitals");
  await page.getByRole("button", { name: "Slow handler (350ms)" }).click();
  await expect(page.getByText(/Last interaction: .* took (3[5-9]\d|[4-9]\d\d) ms/)).toBeVisible();
});

test("errors, rejections and CSP violations are captured", async ({ page }) => {
  await open(page, "debug");
  await page.getByRole("button", { name: "Throw an uncaught error" }).click();
  await page.getByRole("button", { name: "Reject a promise, never catch it" }).click();
  await page.getByRole("button", { name: "Violate the CSP (eval)" }).click();
  const list = page.getByRole("list", { name: "Captured errors" });
  await expect(list).toContainText("[error] Uncaught Error: Boom from a timer callback");
  await expect(list).toContainText("[unhandledrejection] Nobody handled me");
  await expect(list).toContainText("[csp] script-src blocked eval");
  await expect(page.locator('[data-mission="csp"]')).toHaveClass(/done/);
});

test("missions persist across reloads", async ({ page }) => {
  await open(page, "engine");
  await page.getByRole("button", { name: "{ x, y }", exact: true }).click();
  await page.getByRole("button", { name: "{ y, x }", exact: true }).click();
  await expect(page.locator('[data-mission="shapes"]')).toHaveClass(/done/);
  await page.reload();
  await expect(page.locator('[data-mission="shapes"]')).toHaveClass(/done/);
  await expect(page.locator("#progress-text")).toContainText("1 of");
});
```
