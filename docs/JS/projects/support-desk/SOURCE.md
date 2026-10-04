# Support Desk — full source

> The complete, tested source of Support Desk, generated from the files in
> `JS/projects/support-desk/` by `build-source-pages.mjs`. 15 files.

A small, realistic agentic web app: a signed-in customer chats with a support
agent that looks up orders and proposes refunds, which the customer approves
in the UI. It is the project for Parts XI, XII and XIV of the
[JavaScript field guide](../../real-life-JS-guide.md) (§63, §68, §74–§78).

- **Zero runtime dependencies**: `node:http` on the server, plain ES modules in the browser.
- Session-cookie auth (scrypt, HttpOnly, SameSite), CSRF tokens + Fetch Metadata.
- An agent endpoint that streams typed events (AG-UI style) over Server-Sent Events.
- An actionable-UI component registry with prop validation and human approval.
- Strict CSP with Trusted Types, rate limits, one fail-closed error handler.
- A 26-test Playwright suite: UI, API, security, clock, accessibility, visual.

## Run it

```bash
npm start            # http://127.0.0.1:3100
                     # alice@example.com / correct-horse-1   ·   bob@example.com / battery-staple-2
```

Try: "where is A-1002?", "list my orders", "refund A-1001", "where is B-2001?".

## Test it

```bash
npm install
npx playwright install chromium
npm test                                   # 26 tests, ~3s
npx playwright test --ui                   # watch mode with time-travel debugging
npx playwright test --repeat-each=10       # flake check (251/251 when written)
```

The tests start their own server on port 3101 with `NODE_ENV=test`, which
enables the `/__test__/orders` data factory. The "model" is a deterministic
stand-in so the app and the tests run offline; §63 shows how to swap in
WebLLM or an edge-hosted open-weight model behind the same event protocol.

## Files

- `package.json`
- `server.mjs`
- `public/index.html`
- `public/agent-ui.js`
- `public/api.js`
- `public/app.js`
- `public/styles.css`
- `playwright.config.js`
- `tests/a11y.spec.js`
- `tests/auth.setup.js`
- `tests/chat.spec.js`
- `tests/fixtures.js`
- `tests/login.spec.js`
- `tests/security.spec.js`
- `tests/session.spec.js`

## package.json

```json
{
  "name": "support-desk",
  "private": true,
  "type": "module",
  "description": "A conversational support agent app with auth, streaming and actionable UI — the E2E testing project of the JavaScript field guide (§74–§78).",
  "scripts": {
    "start": "node server.mjs",
    "test": "playwright test",
    "test:ui": "playwright test --ui",
    "report": "playwright show-report"
  },
  "devDependencies": {
    "@axe-core/playwright": "^4.13.0",
    "@playwright/test": "^1.63.0"
  }
}
```

## server.mjs

```javascript
// Support Desk — a small, realistic "agentic" web app with zero dependencies.
//
//   node server.mjs            -> http://127.0.0.1:3100
//
// What it demonstrates (each tagged with the guide chapter that explains it):
//   - session-cookie auth with scrypt password hashes, CSRF tokens   (§67)
//   - object-level access control: you can only touch your own orders (§65 A01)
//   - an agent endpoint that streams AG-UI-style events over SSE       (§59, §61)
//   - strict security headers: CSP + Trusted Types, COOP, nosniff     (§65)
//   - rate limits on login and on the (expensive) chat endpoint        (§66 LLM10)
//   - one central error handler: no stack traces leak to the client    (§65 A10)
//
// The "model" is a deterministic, rule-based stand-in so the app (and its
// Playwright suite, §74–§78) runs offline. §63 shows how to swap in WebLLM or an
// edge-hosted open-weight model behind the same event protocol.

import http from "node:http";
import { readFile } from "node:fs/promises";
import { randomBytes, scryptSync, timingSafeEqual } from "node:crypto";
import { extname, join, normalize, sep } from "node:path";
import { fileURLToPath } from "node:url";

const PORT = Number(process.env.PORT ?? 3100);
const HOST = process.env.HOST ?? "127.0.0.1";
const PROD = process.env.NODE_ENV === "production"; // adds `Secure` to cookies
const TEST = process.env.NODE_ENV === "test"; // enables the /__test__/ data factory
const PUBLIC_DIR = fileURLToPath(new URL("./public/", import.meta.url));
const SESSION_TTL_MS = 30 * 60 * 1000;

// ---------------------------------------------------------------- data ----

function hashPassword(password, salt = randomBytes(16)) {
  return { salt, hash: scryptSync(password, salt, 32) };
}

function verifyPassword(password, { salt, hash }) {
  const candidate = scryptSync(password, salt, 32);
  return timingSafeEqual(candidate, hash); // constant-time compare
}

const users = new Map([
  ["alice@example.com", { id: "u1", name: "Alice", email: "alice@example.com", pw: hashPassword("correct-horse-1") }],
  ["bob@example.com", { id: "u2", name: "Bob", email: "bob@example.com", pw: hashPassword("battery-staple-2") }],
]);
// Verifying against a dummy hash for unknown emails keeps the response time
// the same, so timing does not reveal which emails have accounts.
const DUMMY = hashPassword("dummy-password");

const orders = new Map(
  [
    { id: "A-1001", owner: "u1", item: "Noise-cancelling headphones", total: 199, status: "delivered" },
    { id: "A-1002", owner: "u1", item: "Mechanical keyboard", total: 129, status: "in transit" },
    { id: "B-2001", owner: "u2", item: "4K monitor", total: 449, status: "delivered" },
  ].map((o) => [o.id, o]),
);

const sessions = new Map(); // sid -> { userId, csrf, expires }

// ------------------------------------------------------------ helpers ----

const MIME = {
  ".html": "text/html; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".svg": "image/svg+xml",
  ".json": "application/json",
};

const SECURITY_HEADERS = {
  "Content-Security-Policy": [
    "default-src 'self'",
    "script-src 'self'",
    "style-src 'self'",
    "img-src 'self' data:",
    "connect-src 'self'",
    "base-uri 'none'",
    "form-action 'self'",
    "frame-ancestors 'none'",
    "object-src 'none'",
    "require-trusted-types-for 'script'",
    "trusted-types 'none'",
  ].join("; "),
  "X-Content-Type-Options": "nosniff",
  "Referrer-Policy": "strict-origin-when-cross-origin",
  "Cross-Origin-Opener-Policy": "same-origin",
  "Cross-Origin-Resource-Policy": "same-origin",
  "Permissions-Policy": "camera=(), microphone=(), geolocation=()",
};

class HttpError extends Error {
  constructor(status, message) {
    super(message);
    this.status = status;
  }
}

function log(event, fields = {}) {
  // One JSON line per security-relevant event: easy to ship to a SIEM (§65 A09).
  console.log(JSON.stringify({ t: new Date().toISOString(), event, ...fields }));
}

function send(res, status, body, headers = {}) {
  const isJSON = typeof body !== "string" && !Buffer.isBuffer(body);
  res.writeHead(status, {
    ...SECURITY_HEADERS,
    "Cache-Control": "no-store",
    ...(isJSON ? { "Content-Type": "application/json" } : {}),
    ...headers,
  });
  res.end(isJSON ? JSON.stringify(body) : body);
}

async function readJSON(req, limit = 4096) {
  if (!String(req.headers["content-type"] ?? "").startsWith("application/json")) {
    throw new HttpError(415, "Expected application/json");
  }
  let size = 0;
  const chunks = [];
  for await (const chunk of req) {
    size += chunk.length;
    if (size > limit) throw new HttpError(413, "Request body too large");
    chunks.push(chunk);
  }
  try {
    return JSON.parse(Buffer.concat(chunks).toString("utf8"));
  } catch {
    throw new HttpError(400, "Malformed JSON");
  }
}

function parseCookies(header = "") {
  return Object.fromEntries(
    header.split(";").map((p) => p.trim().split("=")).filter(([k, v]) => k && v).map(([k, v]) => [k, decodeURIComponent(v)]),
  );
}

function sessionCookie(sid, maxAgeSec) {
  return [`sid=${sid}`, "HttpOnly", "SameSite=Lax", "Path=/", `Max-Age=${maxAgeSec}`, PROD && "Secure"]
    .filter(Boolean)
    .join("; ");
}

function currentSession(req) {
  const sid = parseCookies(req.headers.cookie).sid;
  const s = sid && sessions.get(sid);
  if (!s) return null;
  if (s.expires < Date.now()) {
    sessions.delete(sid);
    return null;
  }
  s.expires = Date.now() + SESSION_TTL_MS; // sliding expiry
  return { sid, ...s, user: [...users.values()].find((u) => u.id === s.userId) };
}

function requireSession(req) {
  const s = currentSession(req);
  if (!s) throw new HttpError(401, "Not signed in");
  return s;
}

// State-changing requests must carry the CSRF token, and must not be
// cross-site according to Fetch Metadata (§65, §67).
function requireCsrf(req, session) {
  const site = req.headers["sec-fetch-site"];
  if (site && site !== "same-origin" && site !== "none") throw new HttpError(403, "Cross-site request blocked");
  const token = String(req.headers["x-csrf-token"] ?? "");
  const a = Buffer.from(token);
  const b = Buffer.from(session.csrf);
  if (a.length !== b.length || !timingSafeEqual(a, b)) throw new HttpError(403, "Missing or invalid CSRF token");
}

// Fixed-window rate limiter: good enough for one process (§62 shows the edge
// version). check() throws once a key is over its limit; hit() counts one.
function rateLimiter({ limit, windowMs }) {
  const hits = new Map();
  const live = (key) => {
    const h = hits.get(key);
    return h && h.reset > Date.now() ? h : null;
  };
  return {
    check(key) {
      const h = live(key);
      if (h && h.count >= limit) {
        const err = new HttpError(429, "Too many requests, slow down");
        err.headers = { "Retry-After": String(Math.ceil((h.reset - Date.now()) / 1000)) };
        throw err;
      }
    },
    hit(key) {
      const h = live(key);
      if (h) h.count++;
      else hits.set(key, { count: 1, reset: Date.now() + windowMs });
    },
  };
}
// Only *failed* logins count: a user who signs in correctly is never locked out.
const loginLimit = rateLimiter({ limit: Number(process.env.LOGIN_LIMIT ?? 5), windowMs: 60_000 });
const chatLimit = rateLimiter({ limit: Number(process.env.CHAT_LIMIT ?? 20), windowMs: 60_000 });

const publicOrder = ({ owner, ...o }) => o;

// --------------------------------------------------------------- agent ----
// A deterministic "model": it picks a tool from the user's words, runs it,
// and narrates. The event names follow the AG-UI style (§61): the UI never
// sees raw model output it has to parse; it sees typed events.

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

function* words(text) {
  for (const w of text.split(/(?<= )/)) yield w;
}

async function runAgent(user, message, emit, signal) {
  const id = message.match(/\b([A-Z]-\d{4})\b/i)?.[1]?.toUpperCase();
  const say = async (text) => {
    for (const w of words(text)) {
      if (signal.aborted) return;
      emit({ type: "text-delta", delta: w });
      await sleep(25);
    }
  };

  emit({ type: "run-started" });
  if (/refund/i.test(message) && id) {
    emit({ type: "tool-call", name: "lookupOrder", args: { orderId: id } });
    const order = orders.get(id);
    if (!order || order.owner !== user.id) {
      // Same answer for "missing" and "not yours": no enumeration (§65 A01).
      await say(`I couldn't find order ${id} on your account.`);
    } else if (order.status === "refund requested") {
      await say(`A refund for ${id} is already in progress.`);
    } else {
      await say(`I found ${id}. Refunds move money, so I need you to approve this one: `);
      // Excessive agency guard (§66 LLM06): the agent only *proposes*;
      // the human approves, and the server re-checks everything.
      emit({ type: "ui", component: "refund-approval", props: publicOrder(order) });
    }
  } else if (id) {
    emit({ type: "tool-call", name: "lookupOrder", args: { orderId: id } });
    const order = orders.get(id);
    if (!order || order.owner !== user.id) {
      await say(`I couldn't find order ${id} on your account.`);
    } else {
      await say(`Here is order ${id}: `);
      emit({ type: "ui", component: "order-card", props: publicOrder(order) });
    }
  } else if (/orders?|purchases?/i.test(message)) {
    emit({ type: "tool-call", name: "listOrders", args: {} });
    const mine = [...orders.values()].filter((o) => o.owner === user.id);
    await say(`You have ${mine.length} orders. `);
    for (const o of mine) emit({ type: "ui", component: "order-card", props: publicOrder(o) });
  } else {
    await say(
      `Hi ${user.name}! I can look up an order (try "where is A-1002?"), list your orders, or start a refund ("refund A-1001").`,
    );
  }
  emit({ type: "run-finished" });
}

// -------------------------------------------------------------- routes ----

const routes = {
  "GET /healthz": (req, res) => send(res, 200, { ok: true }),

  "POST /api/login": async (req, res) => {
    const { email = "", password = "" } = await readJSON(req);
    const ip = req.socket.remoteAddress;
    const limitKey = `${ip}|${String(email).toLowerCase()}`;
    loginLimit.check(limitKey);
    const user = users.get(String(email).toLowerCase());
    const ok = verifyPassword(String(password), user?.pw ?? DUMMY) && !!user;
    if (!ok) {
      loginLimit.hit(limitKey);
      log("login.failed", { email: String(email).slice(0, 100), ip });
      throw new HttpError(401, "Invalid email or password");
    }
    const sid = randomBytes(32).toString("base64url"); // new id on every login: no session fixation
    const csrf = randomBytes(32).toString("base64url");
    sessions.set(sid, { userId: user.id, csrf, expires: Date.now() + SESSION_TTL_MS });
    log("login.ok", { user: user.id, ip });
    send(res, 200, { user: { name: user.name, email: user.email }, csrf }, {
      "Set-Cookie": sessionCookie(sid, SESSION_TTL_MS / 1000),
    });
  },

  "GET /api/me": (req, res) => {
    const s = requireSession(req);
    send(res, 200, { user: { name: s.user.name, email: s.user.email }, csrf: s.csrf });
  },

  "POST /api/logout": (req, res) => {
    const s = requireSession(req);
    requireCsrf(req, s);
    sessions.delete(s.sid);
    send(res, 200, { ok: true }, { "Set-Cookie": sessionCookie("", 0) });
  },

  "GET /api/orders": (req, res) => {
    const s = requireSession(req);
    send(res, 200, [...orders.values()].filter((o) => o.owner === s.userId).map(publicOrder));
  },

  "POST /api/refunds": async (req, res) => {
    const s = requireSession(req);
    requireCsrf(req, s);
    const { orderId } = await readJSON(req);
    const order = orders.get(String(orderId));
    if (!order || order.owner !== s.userId) throw new HttpError(404, "Order not found");
    order.status = "refund requested";
    log("refund.requested", { user: s.userId, order: order.id });
    send(res, 200, publicOrder(order));
  },

  "POST /api/chat": async (req, res) => {
    const s = requireSession(req);
    requireCsrf(req, s);
    chatLimit.check(s.sid);
    chatLimit.hit(s.sid);
    const { message = "" } = await readJSON(req);
    if (typeof message !== "string" || !message.trim()) throw new HttpError(400, "Message is required");
    if (message.length > 2000) throw new HttpError(413, "Message too long");

    res.writeHead(200, {
      ...SECURITY_HEADERS,
      "Content-Type": "text/event-stream; charset=utf-8",
      "Cache-Control": "no-store",
      "X-Accel-Buffering": "no", // tell nginx-style proxies not to buffer the stream
    });
    const ctrl = new AbortController();
    res.on("close", () => ctrl.abort()); // client pressed Stop or went away
    const emit = (event) => !ctrl.signal.aborted && res.write(`data: ${JSON.stringify(event)}\n\n`);
    await runAgent(s.user, message, emit, ctrl.signal);
    res.end();
  },
};

// Test-data factory (§76): each test creates the orders it needs, so tests
// never depend on each other or on run order. Never mounted outside tests.
let nextTestOrder = 3000;
if (TEST) {
  routes["POST /__test__/orders"] = async (req, res) => {
    const s = requireSession(req);
    const { item = "Test item", total = 10, status = "delivered" } = await readJSON(req);
    const order = { id: `T-${nextTestOrder++}`, owner: s.userId, item, total, status };
    orders.set(order.id, order);
    send(res, 201, publicOrder(order));
  };
}

async function serveStatic(req, res) {
  const url = new URL(req.url, "http://x");
  const rel = url.pathname === "/" ? "index.html" : decodeURIComponent(url.pathname.slice(1));
  const file = normalize(join(PUBLIC_DIR, rel));
  if (!file.startsWith(PUBLIC_DIR) || file.includes(`${sep}.`)) throw new HttpError(404, "Not found"); // path traversal
  try {
    const body = await readFile(file);
    send(res, 200, body, { "Content-Type": MIME[extname(file)] ?? "application/octet-stream", "Cache-Control": "no-cache" });
  } catch {
    throw new HttpError(404, "Not found");
  }
}

export function createServer() {
  return http.createServer(async (req, res) => {
    const path = new URL(req.url, "http://x").pathname;
    const handler = routes[`${req.method} ${path}`];
    try {
      if (handler) await handler(req, res);
      else if (path.startsWith("/api/")) throw new HttpError(404, "Not found");
      else if (req.method === "GET" || req.method === "HEAD") await serveStatic(req, res);
      else throw new HttpError(405, "Method not allowed");
    } catch (err) {
      // Mishandled exceptions are an OWASP 2025 category (A10): fail closed,
      // log the details, return a generic message with a request id.
      const status = err instanceof HttpError ? err.status : 500;
      const requestId = randomBytes(6).toString("hex");
      if (status === 500) log("error", { requestId, path, message: err.message, stack: err.stack });
      if (res.headersSent) return res.end();
      send(res, status, { error: status === 500 ? "Something went wrong" : err.message, requestId }, err.headers);
    }
  });
}

if (process.argv[1] === fileURLToPath(import.meta.url)) {
  createServer().listen(PORT, HOST, () => console.log(`Support Desk on http://${HOST}:${PORT}`));
}
```

## public/index.html

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Support Desk</title>
  <link rel="stylesheet" href="/styles.css">
  <script type="module" src="/app.js"></script>
</head>
<body>
  <header class="topbar">
    <h1>Support Desk</h1>
    <div id="whoami" hidden>
      <span id="user-name"></span>
      <button id="logout" type="button" class="link">Sign out</button>
    </div>
  </header>

  <main id="login-view" hidden>
    <form id="login-form" class="card login" novalidate>
      <h2>Sign in</h2>
      <label>Email <input name="email" type="email" autocomplete="username" required></label>
      <label>Password <input name="password" type="password" autocomplete="current-password" required></label>
      <p id="login-error" class="error" role="alert"></p>
      <button type="submit">Sign in</button>
      <p class="hint">Demo: alice@example.com / correct-horse-1</p>
    </form>
  </main>

  <main id="app-view" class="layout" hidden>
    <aside class="card">
      <h2>Your orders</h2>
      <ul id="orders" aria-label="Your orders"></ul>
    </aside>
    <section class="card chat" aria-labelledby="chat-title">
      <h2 id="chat-title">Assistant</h2>
      <div id="chat-log" role="log" aria-live="polite" aria-label="Conversation"></div>
      <p id="chat-error" class="error" role="alert" hidden>
        <span></span> <button id="retry" type="button" class="link">Retry</button>
      </p>
      <form id="chat-form">
        <label class="sr-only" for="message">Message</label>
        <textarea id="message" name="message" rows="2" maxlength="2000"
                  placeholder="Ask about an order, e.g. where is A-1002?" required></textarea>
        <button id="send" type="submit">Send</button>
        <button id="stop" type="button" hidden>Stop</button>
      </form>
    </section>
  </main>

  <dialog id="session-dialog" aria-labelledby="session-title">
    <h2 id="session-title">Your session is about to expire</h2>
    <p>You've been inactive for a while. Stay signed in?</p>
    <button id="stay" type="button">Stay signed in</button>
  </dialog>
</body>
</html>
```

## public/agent-ui.js

```javascript
// The "actionable UI" layer (§61): the agent never sends HTML. It sends
// { component, props } and this registry decides what may be rendered,
// validates every prop, and wires actions to app code it controls.

// h(): build DOM with createElement + text nodes only. Under the page's
// Trusted Types policy any innerHTML assignment would throw (§65).
export function h(tag, attrs = {}, ...children) {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs)) {
    if (k.startsWith("on") && typeof v === "function") el.addEventListener(k.slice(2), v);
    else if (v === true) el.setAttribute(k, "");
    else if (v !== false && v != null) el.setAttribute(k, String(v));
  }
  for (const c of children.flat()) {
    if (c != null && c !== false) el.append(c instanceof Node ? c : document.createTextNode(String(c)));
  }
  return el;
}

const money = (n) => new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(n);

const orderSchema = { id: "string", item: "string", total: "number", status: "string" };

const registry = {
  "order-card": {
    schema: orderSchema,
    render: (o) =>
      h("article", { class: "order-card", "aria-label": `Order ${o.id}` },
        h("strong", {}, o.id), " ", o.item,
        h("div", { class: "meta" }, money(o.total), " · ", h("span", { class: "status" }, o.status))),
  },
  "refund-approval": {
    schema: orderSchema,
    render: (o, { actions }) => {
      const status = h("p", { class: "status", role: "status" });
      const approve = h("button", {
        type: "button",
        onclick: async () => {
          approve.disabled = decline.disabled = true;
          try {
            await actions.requestRefund(o.id);
            status.textContent = `Refund requested for ${o.id}.`;
          } catch (err) {
            status.textContent = `Refund failed: ${err.message}`;
            approve.disabled = decline.disabled = false;
          }
        },
      }, `Approve refund of ${money(o.total)}`);
      const decline = h("button", {
        type: "button", class: "secondary",
        onclick: () => {
          approve.disabled = decline.disabled = true;
          status.textContent = "Refund declined. Nothing was changed.";
        },
      }, "Decline");
      return h("article", { class: "order-card approval", "aria-label": `Approve refund for ${o.id}` },
        h("strong", {}, o.id), " ", o.item, h("div", { class: "actions" }, approve, decline), status);
    },
  },
};

function validate(schema, props) {
  if (typeof props !== "object" || props === null) return null;
  const clean = {};
  for (const [key, type] of Object.entries(schema)) {
    if (typeof props[key] !== type) return null; // wrong or missing: refuse the whole component
    clean[key] = props[key];
  }
  return clean; // unknown props are dropped
}

export function renderComponent({ component, props }, ctx) {
  const def = Object.hasOwn(registry, component) ? registry[component] : null;
  const clean = def && validate(def.schema, props);
  if (!clean) {
    console.warn("agent-ui: refused component", component);
    return null;
  }
  return def.render(clean, ctx);
}
```

## public/api.js

```javascript
// A tiny API client: JSON in/out, the CSRF header on every state-changing
// request, and errors that carry the HTTP status (§67).

let csrfToken = "";
export const setCsrf = (t) => { csrfToken = t; };

export class ApiError extends Error {
  constructor(status, message) {
    super(message);
    this.status = status;
  }
}

export async function api(path, { method = "GET", body, signal } = {}) {
  const res = await fetch(path, {
    method,
    signal,
    credentials: "same-origin",
    headers: {
      ...(body !== undefined && { "Content-Type": "application/json" }),
      ...(method !== "GET" && { "X-CSRF-Token": csrfToken }),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  if (!res.ok) {
    const { error = res.statusText } = await res.json().catch(() => ({}));
    throw new ApiError(res.status, error);
  }
  return res;
}

// Server-Sent Events over fetch (EventSource can't POST or send headers).
// Yields one parsed event object per `data:` frame (§59).
export async function* readEvents(res) {
  const reader = res.body.pipeThrough(new TextDecoderStream()).getReader();
  let buf = "";
  for (;;) {
    const { value, done } = await reader.read();
    if (done) return;
    buf += value;
    let end;
    while ((end = buf.indexOf("\n\n")) !== -1) {
      const frame = buf.slice(0, end);
      buf = buf.slice(end + 2);
      const data = frame
        .split("\n")
        .filter((line) => line.startsWith("data:"))
        .map((line) => line.slice(5).trimStart())
        .join("\n");
      if (data) yield JSON.parse(data);
    }
  }
}
```

## public/app.js

```javascript
import { api, ApiError, readEvents, setCsrf } from "./api.js";
import { h, renderComponent } from "./agent-ui.js";

const $ = (sel) => document.querySelector(sel);
const IDLE_WARNING_MS = 15 * 60 * 1000;

const state = { user: null, controller: null, lastMessage: "" };

// ------------------------------------------------------------- auth ----

function show(view) {
  $("#login-view").hidden = view !== "login";
  $("#app-view").hidden = view !== "app";
  $("#whoami").hidden = view !== "app";
}

async function boot() {
  try {
    const { user, csrf } = await (await api("/api/me")).json();
    signedIn(user, csrf);
  } catch {
    show("login");
  }
}

function signedIn(user, csrf) {
  state.user = user;
  setCsrf(csrf);
  $("#user-name").textContent = user.name;
  show("app");
  loadOrders();
  if (!$("#chat-log").children.length) {
    addBubble("assistant").append(`Hi ${user.name}! Ask me about an order, or say "refund A-1001".`);
  }
  armIdleTimer();
  $("#message").focus();
}

$("#login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const form = e.currentTarget; // currentTarget is null again after the first await
  $("#login-error").textContent = "";
  try {
    const res = await api("/api/login", { method: "POST", body: Object.fromEntries(new FormData(form)) });
    const { user, csrf } = await res.json();
    form.reset();
    signedIn(user, csrf);
  } catch (err) {
    $("#login-error").textContent = err.message;
  }
});

$("#logout").addEventListener("click", async () => {
  await api("/api/logout", { method: "POST" }).catch(() => {});
  state.user = null;
  $("#chat-log").replaceChildren();
  show("login");
});

// ----------------------------------------------------- idle warning ----

let idleTimer;
function armIdleTimer() {
  clearTimeout(idleTimer);
  idleTimer = setTimeout(() => state.user && $("#session-dialog").showModal(), IDLE_WARNING_MS);
}
for (const type of ["pointerdown", "keydown"]) addEventListener(type, () => state.user && armIdleTimer(), { passive: true });

$("#stay").addEventListener("click", async () => {
  $("#session-dialog").close();
  await api("/api/me").catch(() => show("login")); // refreshes the server's sliding expiry
  armIdleTimer();
});

// ----------------------------------------------------------- orders ----

async function loadOrders() {
  const list = await (await api("/api/orders")).json();
  $("#orders").replaceChildren(
    ...list.map((o) => h("li", { "data-order": o.id }, h("strong", {}, o.id), ` ${o.item} — `, h("span", { class: "status" }, o.status))),
  );
}

const actions = {
  async requestRefund(orderId) {
    await api("/api/refunds", { method: "POST", body: { orderId } });
    await loadOrders();
  },
};

// ------------------------------------------------------------- chat ----

function addBubble(role, text = "") {
  const body = h("div", { class: "body" }, text);
  $("#chat-log").append(h("div", { class: `msg ${role}`, "data-role": role }, body));
  return body;
}

function setBusy(busy) {
  $("#send").disabled = busy;
  $("#stop").hidden = !busy;
  $("#chat-log").setAttribute("aria-busy", String(busy));
}

async function send(message) {
  state.lastMessage = message;
  $("#chat-error").hidden = true;
  const controller = (state.controller = new AbortController());
  setBusy(true);
  let bubble = addBubble("assistant");
  try {
    const res = await api("/api/chat", { method: "POST", body: { message }, signal: controller.signal });
    // Batch text into one DOM write per frame: a fast stream can deliver
    // hundreds of deltas a second and each write would re-layout (§59).
    let pending = "";
    let scheduled = false;
    const flush = () => {
      bubble.append(pending);
      pending = "";
      scheduled = false;
    };
    for await (const event of readEvents(res)) {
      if (event.type === "text-delta") {
        pending += event.delta;
        if (!scheduled) {
          scheduled = true;
          requestAnimationFrame(flush);
        }
      } else if (event.type === "ui") {
        if (pending) flush();
        const node = renderComponent(event, { actions });
        if (node) bubble.append(node);
      }
    }
    if (pending) flush();
  } catch (err) {
    if (err.name === "AbortError") {
      bubble.append(h("em", { class: "stopped" }, " (stopped)"));
    } else {
      bubble.closest(".msg").remove();
      $("#chat-error span").textContent =
        err instanceof ApiError && err.status === 429 ? "You're sending messages too fast." : "The assistant is unavailable.";
      $("#chat-error").hidden = false;
      if (err instanceof ApiError && err.status === 401) show("login");
    }
  } finally {
    setBusy(false);
    state.controller = null;
  }
}

$("#chat-form").addEventListener("submit", (e) => {
  e.preventDefault();
  const message = $("#message").value.trim();
  if (!message) return;
  addBubble("user", message); // textContent, so "<img onerror=…>" is just text
  $("#message").value = "";
  send(message);
});

$("#message").addEventListener("keydown", (e) => {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault();
    $("#chat-form").requestSubmit();
  }
});

$("#stop").addEventListener("click", () => state.controller?.abort());
$("#retry").addEventListener("click", () => send(state.lastMessage));

boot();
```

## public/styles.css

```css
:root {
  --bg: #f6f7f9; --card: #fff; --ink: #1d2330; --muted: #5b6475; --line: #dde1e8;
  --accent: #2f5bea; --accent-ink: #fff; --danger: #b42318; --ok: #067647;
  color-scheme: light dark;
  font: 15px/1.5 system-ui, -apple-system, "Segoe UI", sans-serif;
}
@media (prefers-color-scheme: dark) {
  :root { --bg: #11141a; --card: #1a1f28; --ink: #e7eaf0; --muted: #9aa3b2; --line: #2c3340; --accent: #7c9bff; --accent-ink: #0b1020; --danger: #ff8a80; --ok: #6ee7a8; }
}
* { box-sizing: border-box; }
body { margin: 0; background: var(--bg); color: var(--ink); }
.topbar { display: flex; justify-content: space-between; align-items: center; padding: 12px 20px; border-bottom: 1px solid var(--line); background: var(--card); }
.topbar h1 { font-size: 1.1rem; margin: 0; }
.card { background: var(--card); border: 1px solid var(--line); border-radius: 12px; padding: 16px 20px; }
.card h2 { font-size: 1rem; margin: 0 0 12px; }
.login { max-width: 360px; margin: 64px auto; display: grid; gap: 12px; }
label { display: grid; gap: 4px; font-weight: 500; }
input, textarea { font: inherit; padding: 8px 10px; border: 1px solid var(--line); border-radius: 8px; background: var(--bg); color: var(--ink); }
button { font: inherit; padding: 8px 14px; border: 0; border-radius: 8px; background: var(--accent); color: var(--accent-ink); cursor: pointer; }
button:disabled { opacity: .55; cursor: not-allowed; }
button.secondary { background: transparent; color: var(--ink); border: 1px solid var(--line); }
button.link { background: none; color: var(--accent); padding: 0 4px; text-decoration: underline; }
.error { color: var(--danger); min-height: 1.5em; margin: 0; }
.hint { color: var(--muted); font-size: .85rem; margin: 0; }
.layout { display: grid; grid-template-columns: minmax(220px, 300px) 1fr; gap: 16px; padding: 16px; max-width: 1100px; margin: 0 auto; }
@media (max-width: 720px) { .layout { grid-template-columns: 1fr; } }
#orders { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
#orders li { padding: 8px; border: 1px solid var(--line); border-radius: 8px; }
.status { color: var(--muted); }
.chat { display: grid; grid-template-rows: auto 1fr auto auto; min-height: 70vh; }
#chat-log { display: flex; flex-direction: column; gap: 10px; overflow-y: auto; }
.msg { max-width: 85%; padding: 8px 12px; border-radius: 12px; white-space: pre-wrap; overflow-wrap: anywhere; }
.msg.user { align-self: flex-end; background: var(--accent); color: var(--accent-ink); }
.msg.assistant { align-self: flex-start; background: var(--bg); border: 1px solid var(--line); }
.order-card { margin-top: 8px; padding: 10px; border: 1px solid var(--line); border-radius: 10px; background: var(--card); white-space: normal; }
.order-card .meta { color: var(--muted); font-size: .9rem; }
.order-card .actions { display: flex; gap: 8px; margin-top: 8px; flex-wrap: wrap; }
.stopped { color: var(--muted); }
#chat-form { display: grid; grid-template-columns: 1fr auto auto; gap: 8px; margin-top: 12px; }
dialog { border: 1px solid var(--line); border-radius: 12px; background: var(--card); color: var(--ink); }
.sr-only { position: absolute; width: 1px; height: 1px; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
```

## playwright.config.js

```javascript
// Playwright config for Support Desk (§75).
//   npm test                         run everything
//   npx playwright test --ui         watch mode with time-travel debugging
//   npx playwright test --grep @a11y run one tag
import { defineConfig, devices } from "@playwright/test";

const PORT = 3101; // not the dev port: tests always get a fresh server in test mode

export default defineConfig({
  testDir: "tests",
  fullyParallel: true,
  forbidOnly: !!process.env.CI, // a stray test.only must fail CI, not silently skip the suite
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : undefined,
  reporter: process.env.CI ? [["github"], ["html", { open: "never" }]] : [["list"]],
  use: {
    baseURL: `http://127.0.0.1:${PORT}`,
    trace: "on-first-retry", // a full time-travel trace of any test that flaked
    screenshot: "only-on-failure",
  },
  expect: {
    toHaveScreenshot: { maxDiffPixelRatio: 0.01, animations: "disabled" },
  },
  webServer: {
    command: "node server.mjs",
    url: `http://127.0.0.1:${PORT}/healthz`,
    // Every test shares Alice's session, so the per-session chat limit would
    // throttle the suite itself. Tests of the limiter use their own keys.
    env: { PORT: String(PORT), NODE_ENV: "test", CHAT_LIMIT: "10000" },
    reuseExistingServer: false,
  },
  projects: [
    // 1. Sign in once per user and save the cookies (§75).
    { name: "setup", testMatch: /auth\.setup\.js/ },
    // 2. Everything else starts already signed in as Alice.
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"], storageState: "playwright/.auth/alice.json" },
      dependencies: ["setup"],
    },
    {
      name: "mobile",
      use: { ...devices["Pixel 7"], storageState: "playwright/.auth/alice.json" },
      dependencies: ["setup"],
      grep: /@mobile/,
    },
  ],
});
```

## tests/a11y.spec.js

```javascript
// Automated accessibility checks with axe-core (§77). They catch roughly a
// third of real issues — missing labels, contrast, ARIA misuse — for free.
import AxeBuilder from "@axe-core/playwright";
import { test, expect } from "./fixtures.js";

test("the chat screen has no detectable a11y violations @a11y", async ({ chat, page }) => {
  await chat.ask("where is A-1001?");
  await chat.waitForReply();
  const results = await new AxeBuilder({ page }).withTags(["wcag2a", "wcag2aa", "wcag22aa"]).analyze();
  expect(results.violations).toEqual([]);
});
```

## tests/auth.setup.js

```javascript
// Runs once before the other projects: sign in through the real UI, then
// save cookies + storage to a file every other test starts from (§75).
import { test as setup, expect } from "@playwright/test";
import { login, USERS } from "./fixtures.js";

setup("sign in as Alice", async ({ page }) => {
  await login(page, USERS.alice);
  await expect(page.getByRole("list", { name: "Your orders" })).toBeVisible();
  await page.context().storageState({ path: "playwright/.auth/alice.json" });
});
```

## tests/chat.spec.js

```javascript
import { test, expect, sse } from "./fixtures.js";

test("a reply streams in and the Stop button shows while it does", async ({ chat }) => {
  await chat.ask("hello");
  await expect(chat.stop).toBeVisible(); // still streaming
  await chat.waitForReply();
  await expect(chat.stop).toBeHidden();
  await expect(chat.lastReply()).toContainText("I can look up an order");
});

test("looking up an order renders an order card @mobile", async ({ chat }) => {
  await chat.ask("where is A-1002?");
  await chat.waitForReply();
  // An ARIA snapshot asserts structure and accessible names in one go (§77).
  await expect(chat.lastReply().getByRole("article")).toMatchAriaSnapshot(`
    - article "Order A-1002":
      - strong: A-1002
      - text: /Mechanical keyboard/
  `);
});

test("refund needs a human approval, then updates the order list", async ({ chat, seedOrder, page }) => {
  // The server is shared by every test (and by retries), so this test makes
  // its own order instead of mutating A-1001/A-1002 that others rely on.
  const order = await seedOrder({ item: "Desk lamp", total: 35 });
  await page.reload(); // pick up the new order in the sidebar
  await chat.ask(`please refund ${order.id}`);
  await chat.waitForReply();
  const card = chat.lastReply().getByRole("article", { name: `Approve refund for ${order.id}` });
  await expect(card).toBeVisible();
  await card.getByRole("button", { name: "Approve refund of $35.00" }).click();
  await expect(card.getByRole("status")).toHaveText(`Refund requested for ${order.id}.`);
  await expect(chat.orders.locator(`[data-order="${order.id}"]`)).toContainText("refund requested");
});

test("declining a refund changes nothing", async ({ chat }) => {
  await chat.ask("refund A-1001");
  await chat.waitForReply();
  const card = chat.lastReply().getByRole("article", { name: "Approve refund for A-1001" });
  await card.getByRole("button", { name: "Decline" }).click();
  await expect(card.getByRole("status")).toHaveText("Refund declined. Nothing was changed.");
  await expect(chat.orders.locator('[data-order="A-1001"]')).toContainText("delivered");
});

test("you cannot look up someone else's order", async ({ chat }) => {
  await chat.ask("where is B-2001?");
  await chat.waitForReply();
  await expect(chat.lastReply()).toHaveText("I couldn't find order B-2001 on your account.");
});

test("Stop aborts the stream mid-reply", async ({ chat, page }) => {
  // Hold the response for a second, so there is time to press Stop mid-request.
  await page.route("**/api/chat", async (route) => {
    const body = sse(
      { type: "run-started" },
      ...Array.from({ length: 50 }, (_, i) => ({ type: "text-delta", delta: `word${i} ` })),
      { type: "run-finished" },
    );
    await new Promise((r) => setTimeout(r, 1000)); // hold the response
    await route.fulfill({ contentType: "text/event-stream", body }).catch(() => {}); // the page may have aborted
  });
  await chat.ask("tell me a long story");
  await chat.stop.click();
  await expect(chat.lastReply()).toContainText("(stopped)");
  await expect(chat.send).toBeEnabled();
});

test.describe("with a mocked agent", () => {
  test("a server error shows Retry, and Retry recovers", async ({ chat, page }) => {
    let calls = 0;
    await page.route("**/api/chat", (route) =>
      ++calls === 1
        ? route.fulfill({ status: 500, json: { error: "boom" } })
        : route.fulfill({ contentType: "text/event-stream", body: sse({ type: "text-delta", delta: "Back online." }) }),
    );
    await chat.ask("hello?");
    await expect(page.getByRole("alert")).toContainText("The assistant is unavailable.");
    await page.getByRole("button", { name: "Retry" }).click();
    await expect(chat.lastReply()).toHaveText("Back online.");
    await expect(page.getByRole("alert")).toBeHidden();
  });

  test("model output is rendered as text, never as HTML", async ({ chat, page }) => {
    const dialogs = [];
    page.on("dialog", (d) => (dialogs.push(d.message()), d.dismiss()));
    const payload = `<img src=x onerror="alert('xss')">`;
    await page.route("**/api/chat", (route) =>
      route.fulfill({ contentType: "text/event-stream", body: sse({ type: "text-delta", delta: payload }) }),
    );
    await chat.ask("hi");
    await expect(chat.lastReply()).toHaveText(payload);
    await expect(chat.lastReply().locator("img")).toHaveCount(0);
    expect(dialogs).toEqual([]);
  });

  test("the UI refuses components it does not know or props that fail validation", async ({ chat, page }) => {
    await page.route("**/api/chat", (route) =>
      route.fulfill({
        contentType: "text/event-stream",
        body: sse(
          { type: "ui", component: "iframe", props: { src: "https://evil.example" } },
          { type: "ui", component: "order-card", props: { id: "X-1", item: "Thing", total: "free", status: "?" } },
          { type: "ui", component: "order-card", props: { id: "X-2", item: "Real thing", total: 5, status: "ok" } },
        ),
      }),
    );
    await chat.ask("show me");
    await chat.waitForReply();
    await expect(chat.lastReply().getByRole("article")).toHaveCount(1);
    await expect(chat.lastReply().getByRole("article")).toHaveAccessibleName("Order X-2");
  });
});
```

## tests/fixtures.js

```javascript
// Shared fixtures (§76): page objects and helpers every spec can ask for by
// name. Playwright builds them lazily, only for the tests that use them.
import { test as base, expect } from "@playwright/test";

export const USERS = {
  alice: { email: "alice@example.com", password: "correct-horse-1" },
  bob: { email: "bob@example.com", password: "battery-staple-2" },
};

// A page object: the one place that knows how the chat UI is built.
// Tests talk in user intent ("ask", "lastReply"), not CSS selectors.
export class ChatPage {
  constructor(page) {
    this.page = page;
    this.input = page.getByRole("textbox", { name: "Message" });
    this.send = page.getByRole("button", { name: "Send" });
    this.stop = page.getByRole("button", { name: "Stop" });
    this.log = page.getByRole("log", { name: "Conversation" });
    this.orders = page.getByRole("list", { name: "Your orders" });
  }
  async goto() {
    await this.page.goto("/");
    await expect(this.input).toBeVisible();
  }
  async ask(text) {
    await this.input.fill(text);
    await this.send.click();
  }
  lastReply() {
    return this.log.locator('[data-role="assistant"]').last();
  }
  async waitForReply() {
    // aria-busy flips back to false when the stream ends.
    await expect(this.log).toHaveAttribute("aria-busy", "false");
  }
}

export async function login(page, { email, password }) {
  await page.goto("/");
  await page.getByLabel("Email").fill(email);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: "Sign in" }).click();
}

// Build an SSE body by hand, for tests that mock the agent (§77).
export const sse = (...events) => events.map((e) => `data: ${JSON.stringify(e)}\n\n`).join("");

export const test = base.extend({
  chat: async ({ page }, use) => {
    const chat = new ChatPage(page);
    await chat.goto();
    await use(chat);
  },
  // seedOrder({ item, total }) -> a fresh order owned by the signed-in user.
  seedOrder: async ({ page }, use) => {
    await use(async (order = {}) => {
      const res = await page.request.post("/__test__/orders", { data: order });
      expect(res.status()).toBe(201);
      return res.json();
    });
  },
});

export { expect };
```

## tests/login.spec.js

```javascript
import { test, expect } from "@playwright/test";
import { login, USERS } from "./fixtures.js";

// These tests are about signing in, so they start signed out.
test.use({ storageState: { cookies: [], origins: [] } });

test("a failed sign-in shows a generic error", async ({ page }, testInfo) => {
  // Not Alice: failed logins count toward a per-email lockout, and repeated
  // runs would lock out the user every other test signs in as.
  const email = `nobody-${testInfo.repeatEachIndex}-${testInfo.retry}@example.com`;
  await login(page, { email, password: "nope" });
  // The server says the same thing for an unknown email and a wrong
  // password, so the form can't be used to discover who has an account.
  await expect(page.getByRole("alert")).toHaveText("Invalid email or password");
  await expect(page.getByRole("list", { name: "Your orders" })).toBeHidden();
});

test("signing in shows the user's own orders", async ({ page }) => {
  await login(page, USERS.bob);
  await expect(page.getByText("Bob", { exact: true })).toBeVisible();
  await expect(page.getByRole("list", { name: "Your orders" }).getByRole("listitem")).toHaveCount(1);
});

test("the session cookie is HttpOnly and SameSite", async ({ page, context }) => {
  await login(page, USERS.bob);
  await expect(page.getByRole("list", { name: "Your orders" })).toBeVisible();
  const sid = (await context.cookies()).find((c) => c.name === "sid");
  expect(sid).toMatchObject({ httpOnly: true, sameSite: "Lax" });
  // ...and so page JavaScript (and any XSS) cannot read it.
  expect(await page.evaluate(() => document.cookie)).not.toContain("sid=");
});

test("sign out ends the session on the server too", async ({ page }) => {
  await login(page, USERS.bob);
  await page.getByRole("button", { name: "Sign out" }).click();
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  expect((await page.request.get("/api/me")).status()).toBe(401);
});

test("login page looks right @visual", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
  await expect(page).toHaveScreenshot("login.png");
});
```

## tests/security.spec.js

```javascript
// API-level security tests: no browser page needed, just Playwright's
// request fixture, which shares cookies with the signed-in context (§77).
import { test, expect } from "@playwright/test";
import { USERS } from "./fixtures.js";

const csrfOf = async (request) => (await (await request.get("/api/me")).json()).csrf;

test("every response carries the security headers", async ({ request }) => {
  const res = await request.get("/");
  const h = res.headers();
  expect(h["content-security-policy"]).toContain("require-trusted-types-for 'script'");
  expect(h["content-security-policy"]).toContain("frame-ancestors 'none'");
  expect(h["x-content-type-options"]).toBe("nosniff");
  expect(h["cross-origin-opener-policy"]).toBe("same-origin");
});

test("a state-changing request without the CSRF token is refused", async ({ request }) => {
  const res = await request.post("/api/refunds", { data: { orderId: "A-1001" } });
  expect(res.status()).toBe(403);
});

test("a cross-site request is refused even with the token", async ({ request }) => {
  const res = await request.post("/api/refunds", {
    data: { orderId: "A-1001" },
    headers: { "x-csrf-token": await csrfOf(request), "sec-fetch-site": "cross-site" },
  });
  expect(res.status()).toBe(403);
});

test("Alice cannot refund Bob's order (no IDOR)", async ({ request }) => {
  const res = await request.post("/api/refunds", {
    data: { orderId: "B-2001" },
    headers: { "x-csrf-token": await csrfOf(request) },
  });
  expect(res.status()).toBe(404); // not 403: don't confirm the order exists
});

test("errors never leak stack traces", async ({ request }) => {
  const res = await request.post("/api/chat", {
    headers: { "x-csrf-token": await csrfOf(request), "content-type": "application/json" },
    data: Buffer.from("{not json"), // a string would be JSON-encoded into valid JSON
  });
  expect(res.status()).toBe(400);
  const body = await res.json();
  expect(body).toEqual({ error: "Malformed JSON", requestId: expect.stringMatching(/^[0-9a-f]{12}$/) });
});

test.describe("signed out", () => {
  test.use({ storageState: { cookies: [], origins: [] } });

  test("the API requires a session", async ({ request }) => {
    expect((await request.get("/api/orders")).status()).toBe(401);
  });

  test("login is rate-limited after 5 failed attempts", async ({ request }, testInfo) => {
    // A unique key per run: the limiter's memory outlives this test, and
    // retries or --repeat-each would otherwise start already locked out.
    const email = `ratelimit-${testInfo.workerIndex}-${testInfo.repeatEachIndex}-${testInfo.retry}@example.com`;
    const attempt = () => request.post("/api/login", { data: { email, password: "x" } });
    for (let i = 0; i < 5; i++) expect((await attempt()).status()).toBe(401);
    const blocked = await attempt();
    expect(blocked.status()).toBe(429);
    expect(Number(blocked.headers()["retry-after"])).toBeGreaterThan(0);
  });

  test("two users in two contexts see only their own data", async ({ browser }) => {
    // Two isolated browser contexts = two users at once, in one test.
    const [alice, bob] = await Promise.all([browser.newContext(), browser.newContext()]);
    for (const [ctx, user] of [[alice, USERS.alice], [bob, USERS.bob]]) {
      expect((await ctx.request.post("/api/login", { data: user })).ok()).toBe(true);
    }
    const ids = async (ctx) => (await (await ctx.request.get("/api/orders")).json()).map((o) => o.id);
    // Other tests seed extra orders for Alice, so assert ownership, not exact lists.
    const [aliceIds, bobIds] = await Promise.all([ids(alice), ids(bob)]);
    expect(aliceIds).toEqual(expect.arrayContaining(["A-1001", "A-1002"]));
    expect(aliceIds).not.toContain("B-2001");
    expect(bobIds).toEqual(["B-2001"]);
    await Promise.all([alice.close(), bob.close()]);
  });
});
```

## tests/session.spec.js

```javascript
// Time-dependent UI without waiting 15 real minutes: Playwright's clock (§77).
import { test, expect } from "./fixtures.js";

test("an idle session warns before it expires", async ({ page }) => {
  await page.clock.install(); // fake Date, setTimeout, setInterval... before the app loads
  await page.goto("/");
  await expect(page.getByRole("textbox", { name: "Message" })).toBeVisible();

  await page.clock.fastForward("14:59");
  await expect(page.getByRole("dialog")).toBeHidden();

  await page.clock.fastForward("00:02");
  const dialog = page.getByRole("dialog", { name: "Your session is about to expire" });
  await expect(dialog).toBeVisible();
  await dialog.getByRole("button", { name: "Stay signed in" }).click();
  await expect(dialog).toBeHidden();
});
```
