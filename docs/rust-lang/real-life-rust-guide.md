# Rust — The Complete Field Guide (Beginner → Expert)

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/rust/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> A practical, example-driven path through Rust: what the language actually
> gives you, why it is shaped the way it is, and how to use it to build real
> CLIs, network services, systems tools, and security-sensitive code.
>
> Rust is the third language in this Wiki's languages track. The
> [Go guide](../Golang/real-life-golang-guide.md) shows how a garbage collector
> and a runtime make systems work comfortable. The
> [C guide](../c-lang/real-life-c-guide.md) shows what the machine is really
> doing, and how much discipline it takes to keep memory correct by hand.
> Rust moves the discipline from C into the compiler, without adding Go's
> garbage collector. This guide explains how it does that and what you have
> to learn in return.
>
> Every concept comes with runnable code and a **Real-world example** or
> **War story**. Each Part ends with projects you can build and run today.
> Read it once from top to bottom, then keep it as a reference.

---

> **The series:** [1 OS](../os-linux/real-life-os-guide.md) → [2 Networking](../networking/real-life-example-osi.md) → [3 Security](../security/real-life-guide.md) → [4 HTTPS walkthrough](../v2-https/real-life-guide-v1.md), with [Go](../Golang/real-life-golang-guide.md) and **Rust** alongside. Languages track: Go → [C](../c-lang/real-life-c-guide.md) → **Rust**.
>
> **You are here: languages track, step 3: Rust.** ← Previous: [C — The Complete Field Guide](../c-lang/real-life-c-guide.md). Start the series at [step 1: Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md).
>
> [The full series map](#0-2-the-series-os-networking-security-https-go-c-rust).

---

## How to use this guide

- **Part 0 (Start here):** setup, how this guide fits the series, and a
  one-page Go ↔ C ↔ Rust phrasebook. Read the phrasebook even if you skip
  everything else in Part 0.
- **Beginner (Part I):** the toolchain, types, expressions, **ownership,
  borrowing**, slices and strings, enums and pattern matching, errors,
  collections, modules. Chapters 6–8 are the core of the language. Part I
  ends with a **🔎 Checkpoint** and two terminal projects.
- **Intermediate (Part II):** traits, generics, trait objects, the standard
  traits, **lifetimes**, closures, iterators, smart pointers, testing, I/O.
  Ends with a checkpoint, a `grep` clone, and an LRU cache. The LRU cache
  shows why a doubly linked list is hard in Rust and what to do instead.
- **Advanced (Part III):** threads, `Send`/`Sync`, locks, atomics, channels,
  **async Rust and Tokio**, networking, `unsafe`, FFI with C, macros,
  performance. Ends with a thread pool and a Tokio chat server.
- **Systems Rust (Part IV):** the [OS guide](../os-linux/real-life-os-guide.md)
  rebuilt in Rust: syscalls, processes and signals, crash-safe files, `/proc`,
  `mmap`, custom allocators, and Linux namespaces.
- **Network and HTTP services (Part V):** the [networking](../networking/tcp-ip/real-life-guide-v1.md)
  and [HTTPS lifecycle](../v2-https/real-life-guide-v1.md) guides in Rust:
  a DNS client from raw bytes, HTTP from a bare socket, a production Axum
  server, safe outbound clients, TLS and mTLS with `rustls`, and a reverse proxy.
- **Security-focused Rust (Part VI):** what Rust prevents and what it does
  not, cryptography done correctly, parsing untrusted input, fuzzing, supply
  chain, and a JWT auth service. Ties to both
  [security guides](../security/real-life-guide.md).
- **Expert (Part VII):** memory layout, vtables, how `async` compiles,
  `Pin`, variance and drop check, memory ordering, API design. Ends with two
  capstones: a Redis-like key-value server, and one HTTPS request traced
  through every guide, in Rust.
- **Professional Rust (Part VIII):** debugging, workspaces, CI, release
  builds, containers and Kubernetes, and an anti-pattern catalog.
- **Appendices:** a compiler-error decoder, a cheat sheet, a glossary, a
  90-day plan, and further reading.

Conventions:

- Code targets **Rust 1.92+ with edition 2024**. Every complete program in
  this guide was compiled with `rustc 1.92.0`. Anything newer works.
  Third-party crate versions are the ones current when the guide was written.
  They are pinned in each `Cargo.toml` snippet.
- Snippets that are **wrong on purpose**, such as borrow-checker errors,
  start with `// DOES NOT COMPILE`. Each one is followed by the compiler's
  error and the fix. They were checked to fail with the error shown.
- **Real-world example** boxes connect a language feature to code you will
  actually write. **War story** boxes describe a real class of incident or
  design trade-off.
- **Across the series** notes link to the chapter in another guide that
  covers the same topic from the OS, network, security, Go, or C side.

---

## Table of contents

**Part 0 — Start here**
- 0.1 Setup: rustup, cargo, rust-analyzer, clippy
- 0.2 The series: OS → networking → security → HTTPS → Go → C → Rust
- 0.3 The Go ↔ C ↔ Rust phrasebook
- 0.4 The Rust labs: what you build for each guide in the series

**Part I — Foundations: the Rust way**
1. What Rust is, and why it exists
2. The toolchain: cargo, crates, editions, and what `cargo build` really does
3. Variables, mutability, shadowing, and scalar types
4. Control flow: everything is an expression
5. Functions, statements, and the unit and never types
6. Ownership: moves, copies, and drop
7. Borrowing: shared XOR mutable
8. Slices and strings: fat pointers and UTF-8
9. Structs, enums, and pattern matching
10. Error handling: `Result`, `?`, and panics
11. Collections: `Vec`, `HashMap`, and friends
12. Modules, crates, and visibility
13. 🔎 Checkpoint: ownership and borrowing
14. Terminal project: a word-frequency counter
15. Terminal project: a `/etc/hosts`-style config parser with real errors

**Part II — Intermediate: traits, lifetimes, and the standard library**
16. Traits and generics: shared behavior without inheritance
17. Trait objects: `dyn Trait` and dynamic dispatch
18. The standard traits you will implement every week
19. Lifetimes: naming how long a borrow is valid
20. Closures: `Fn`, `FnMut`, `FnOnce`
21. Iterators: lazy, composable, zero-cost
22. Smart pointers: `Box`, `Rc`, `Arc`, `Cell`, `RefCell`, `Weak`
23. Testing: unit, integration, doc tests, and property tests
24. I/O, files, processes, and application errors
25. 🔎 Checkpoint: traits and lifetimes
26. Terminal project: `rgrep`, a tested `grep` clone
27. Terminal project: an LRU cache, and why linked lists are hard in Rust

**Part III — Advanced: concurrency, async, unsafe**
28. Threads, and why data races do not compile
29. Shared state: `Mutex`, `RwLock`, atomics, `Condvar`, channels
30. Async Rust: futures, `poll`, and why async exists
31. Tokio in practice: tasks, `select!`, timeouts, cancellation, shutdown
32. Networking with `std` and Tokio: TCP, UDP, and framing
33. `unsafe` Rust: the five superpowers and the contract they carry
34. FFI: calling C from Rust and Rust from C
35. Macros: `macro_rules!` and what derive macros do
36. Performance: release profiles, benchmarks, profiling, allocations
37. Project: a thread pool from scratch
38. Project: a multi-room chat server with Tokio

**Part IV — Systems Rust: the OS guide in Rust**
39. Syscalls from Rust: `std`, `libc`, and `nix`
40. Processes, exit codes, and signals
41. Files that survive crashes: `fsync` and atomic replacement
42. Reading `/proc`: build your own `ps`
43. Memory: `mmap` and a counting global allocator
44. A container in 100 lines: namespaces with `nix`

**Part V — Network and HTTP services: the networking and HTTPS guides in Rust**
45. A DNS client from raw bytes
46. HTTP/1.1 from a bare `TcpListener`
47. A production HTTP server with Axum: limits, timeouts, tracing, shutdown
48. Outbound requests: timeouts, retries, idempotency, and an SSRF guard
49. TLS and mutual TLS with `rustls`
50. A reverse proxy with Hyper, and request-smuggling defenses

**Part VI — Security-focused Rust**
51. What Rust prevents, and what it does not
52. Cryptography done right: AEAD, password hashing, HMAC, constant-time, zeroize
53. Parsing untrusted input: newtypes, "parse, don't validate", resource bounds
54. Fuzzing and property testing
55. Supply chain: `Cargo.lock`, `cargo audit`, `cargo deny`, `build.rs`
56. Project: an auth service with Argon2 and JWTs

**Part VII — Expert: internals and design**
57. Memory layout: size, alignment, niches, `repr(C)`, fat pointers
58. Inside `Vec`, `String`, `Box`, and trait objects
59. How `async` compiles: state machines, `Pin`, and a tiny executor
60. Variance, `PhantomData`, and drop check
61. Memory ordering: `Relaxed`, `Acquire`/`Release`, `SeqCst`, and a spinlock
62. API design: newtypes, typestate, builders, errors, and semver
63. Capstone: a Redis-compatible key-value server with persistence
64. Capstone: one HTTPS request, every layer, every guide, in Rust

**Part VIII — Professional Rust**
65. Debugging: backtraces, `lldb`/`gdb`, `tokio-console`, Miri
66. Workspaces, CI quality gates, release builds, and cross-compilation
67. Running Rust in containers and Kubernetes
68. How not to write Rust: an anti-pattern catalog

**Appendices**
- A. Compiler-error decoder: the 15 errors you will actually see
- B. Cargo and tooling cheat sheet
- C. Glossary
- D. A 90-day plan
- E. Further reading and codebases worth reading

---

# Part 0 — Start here

## 0.1 Setup: rustup, cargo, rust-analyzer, clippy

```bash
# Install the toolchain manager (macOS/Linux). It installs rustc, cargo,
# rustfmt, clippy and the standard library docs.
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

rustc --version          # the compiler
cargo --version          # build tool + package manager + test runner
rustup update            # Rust ships a new stable release every 6 weeks
rustup component add rust-src rust-analyzer clippy rustfmt
rustup doc --std         # the standard library docs, offline
```

Editor: install **rust-analyzer** (VS Code extension "rust-analyzer").
It shows inferred types inline. When you are learning Rust that is the most
useful feature an editor can have, because most confusion is about what type
a value has and who owns it.

Lab machine: Parts I–III and V–VII run on macOS or Linux. Part IV (Systems
Rust) needs Linux, because `/proc`, namespaces, and `cgroups` exist only
there. Use the same lab VM or container you built in the
[OS guide §0.3](../os-linux/real-life-os-guide.md#0-3-build-your-lab):

```bash
# A disposable Linux box with Rust, from macOS:
docker run --rm -it --privileged -v "$PWD":/work -w /work rust:1 bash
```

## 0.2 The series: OS → networking → security → HTTPS → Go → C → Rust

This guide belongs to one course, read in this order:

```
   1  OS & Linux  ──▶  2  Networking  ──▶  3  Security  ──▶  4  HTTPS walkthrough
                         2a OSI map          3a From Zero        (one request through
                         2b TCP/IP           3b In Depth          every guide)
   ════════════════  Go and Rust, alongside every step  ════════════════

   Languages track:   Go  ──▶  C  ──▶  Rust  (you are here)
                      GC +     the       C's control, with the
                      runtime  machine   discipline checked by the compiler
```

| Step | Guide | What it gives you | What you rebuild here in Rust |
|---|---|---|---|
| 1 | [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md) | processes, memory, files, sockets, signals, containers | Part IV: syscalls, signals, `fsync`, `ps` from `/proc`, `mmap`, a namespace container |
| 2a | [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md) | one click through all seven layers | Chapter 64 walks the same layers |
| 2b | [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md) | addressing, TCP, UDP, DNS, TLS, packet capture | Ch 32, 45, 46: TCP/UDP servers, framing, a raw-bytes DNS client, HTTP from a socket |
| 3a | [Security from Zero](../security/real-life-guide.md) | hashing, AEAD, passwords, TLS, PKI, OWASP Top 10 | Ch 49, 52, 56: `rustls`, AES-GCM, Argon2, HMAC, JWT auth |
| 3b | [Security Engineering in Depth](../security/real-life-security-guide-v1.md) | supply chain, smuggling, SSRF, races, parsers | Ch 48, 50, 53–55: SSRF guard, smuggling-safe proxy, fuzzing, `cargo deny` |
| 4 | [The HTTPS Request Lifecycle](../v2-https/real-life-guide-v1.md) | one request end to end, then the server side | Ch 47–50, 64: production server, retries, proxy, the whole-series walkthrough |
| L1 | [Go — The Complete Field Guide](../Golang/real-life-golang-guide.md) | a GC'd, runtime-scheduled systems language | the phrasebook in §0.3; Go comparisons in every chapter |
| L2 | [C — The Complete Field Guide](../c-lang/real-life-c-guide.md) | manual memory, UB, the ABI, sockets, the kernel | Ch 6–7 (ownership replaces C's discipline), 33–34 (unsafe, FFI), 57 (layout) |
| L3 | **Rust — The Complete Field Guide** ← you are here | memory safety without a GC, fearless concurrency | — |

**Why Rust comes last in the languages track.** Rust's rules make the most
sense once you have seen both alternatives fail. In C you have seen a
use-after-free, a double free, and a data race. In Go you have seen GC pauses,
nil-pointer panics, and the slice aliasing bug in
[Go §6.6](../Golang/real-life-golang-guide.md#6-arrays-and-slices-the-internals-that-explain-the-bugs).
After that, the borrow checker stops looking like an obstacle. It is the C
guide's [ownership discipline](../c-lang/real-life-c-guide.md#10-dynamic-memory-the-malloc-family-ownership-and-the-bugs-it-enables),
written down as rules and enforced at compile time.

**Reading paths:**

- **From Go (backend engineer):** Part 0 → Ch 1–12 (slowly, especially 6–8)
  → 16–22 → 28–31 → 47–48 → 62 → 68.
- **From C (systems / embedded):** Part 0 → Ch 1, 6–10 → 16–19 → 22 →
  28–29 → 33–34 → Part IV → 57–61.
- **Security engineer:** Ch 1, 6–10 → 33 → Part VI → 49–50 → 55.
- **SRE / platform:** Ch 1–12 → 24 → 30–31 → Part IV → 47 → 65–67.

### Across the series, at a glance

| Topic | OS / network / security side | Go side | C side | Rust chapter |
|---|---|---|---|---|
| System calls | [OS Ch 2](../os-linux/real-life-os-guide.md#chapter-2-kernel-space-vs-user-space-and-the-system-call) | [OS Ch 73](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime) | [C §21](../c-lang/real-life-c-guide.md#21-sockets-a-tcp-echo-server-then-a-tiny-http-server) | 39 |
| Virtual memory, stack vs heap | [OS Ch 10](../os-linux/real-life-os-guide.md#chapter-10-virtual-memory-every-process-s-private-lie) | [Go §27](../Golang/real-life-golang-guide.md#27-the-go-memory-model-and-the-garbage-collector) | [C §10](../c-lang/real-life-c-guide.md#10-dynamic-memory-the-malloc-family-ownership-and-the-bugs-it-enables) | 6, 57 |
| Threads and races | [OS Ch 17](../os-linux/real-life-os-guide.md#chapter-17-race-conditions-at-the-os-level) | [Go §17](../Golang/real-life-golang-guide.md#17-the-sync-package-mutexes-waitgroup-once-atomics) | [C §20](../c-lang/real-life-c-guide.md#20-concurrency-posix-threads-mutexes-condvars-c11-threads-h-and-atomics) | 28–29 |
| Signals | [OS Ch 20](../os-linux/real-life-os-guide.md#chapter-20-signals-the-os-s-tap-on-the-shoulder) | [OS Ch 74](../os-linux/real-life-os-guide.md#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1) | [C §55](../c-lang/real-life-c-guide.md#55-signals-handling-asynchronous-events-safely) | 40 |
| Crash-safe files | [OS Ch 14](../os-linux/real-life-os-guide.md#chapter-14-the-i-o-stack-from-read-to-the-disk-platter) | [OS Ch 77](../os-linux/real-life-os-guide.md#chapter-77-files-that-survive-crashes-page-cache-fsync-and-atomic-replacement) | [C §12](../c-lang/real-life-c-guide.md#12-file-i-o-buffering-binary-vs-text-mode-errno) | 41 |
| Many connections | [OS Ch 19](../os-linux/real-life-os-guide.md#chapter-19-pipes-sockets-shared-memory-and-message-queues) | [OS Ch 78](../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections) | [C §23.2](../c-lang/real-life-c-guide.md#23-mini-projects-a-thread-pool-and-a-select-based-chat-server) | 30–32 |
| Containers | [OS Ch 44](../os-linux/real-life-os-guide.md#chapter-44-what-a-container-actually-is-namespaces-cgroups-a-filesystem) | [OS Ch 80](../os-linux/real-life-os-guide.md#chapter-80-a-container-runtime-in-150-lines-of-go) | — | 44 |
| DNS | [Net Ch 18](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses) | [Net §0.8](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself) | — | 45 |
| TCP | [Net Ch 21](../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake) | [Go §34](../Golang/real-life-golang-guide.md#34-networking-deep-dive-net-conn-tcp-udp-framing-your-own-protocol) | [C §21](../c-lang/real-life-c-guide.md#21-sockets-a-tcp-echo-server-then-a-tiny-http-server) | 32 |
| HTTP | [Net Ch 28](../networking/tcp-ip/real-life-guide-v1.md#chapter-28-http-how-the-web-actually-talks), [HTTPS Ch 5](../v2-https/real-life-guide-v1.md#chapter-5-http-the-conversation) | [Go §24](../Golang/real-life-golang-guide.md#24-net-http-fundamentals-client-server-middleware) | [C §21.2](../c-lang/real-life-c-guide.md#21-sockets-a-tcp-echo-server-then-a-tiny-http-server) | 46–47 |
| TLS / mTLS | [Sec Ch 28](../security/real-life-guide.md#chapter-28-the-tls-1-3-handshake-step-by-step), [Sec Ch 53](../security/real-life-guide.md#chapter-53-project-1-your-own-ca-plus-mutual-tls) | [Go §50](../Golang/real-life-golang-guide.md#50-mtls-client-certificate-authentication-end-to-end) | — | 49 |
| Password storage | [Sec Ch 10](../security/real-life-guide.md#chapter-10-password-storage-is-a-completely-different-problem) | [Go §48](../Golang/real-life-golang-guide.md#48-auth-service-password-hashing-and-jwt-issuing-verification) | — | 52, 56 |
| Request smuggling | [Sec 3b Ch 30](../security/real-life-security-guide-v1.md#chapter-30-http-request-smuggling-and-desync) | [HTTPS Ch 20](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling) | — | 50 |
| SSRF | [Sec 3b Ch 33](../security/real-life-security-guide-v1.md#chapter-33-ssrf-mastery) | [HTTPS Ch 22](../v2-https/real-life-guide-v1.md#chapter-22-outbound-requests-ssrf-and-safe-http-clients) | — | 48 |
| Retries, idempotency | [HTTPS Ch 21](../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers) | same | — | 48 |
| Undefined behavior | — | — | [C §28](../c-lang/real-life-c-guide.md#28-undefined-behavior-the-list-every-c-programmer-must-memorize) | 33, 51 |
| Supply chain | [Sec Ch 47](../security/real-life-guide.md#chapter-47-supply-chain-and-the-code-you-didn-t-write), [3b Ch 13](../security/real-life-security-guide-v1.md#chapter-13-securing-the-software-factory-provenance-and-admission) | [Go §64](../Golang/real-life-golang-guide.md#64-professional-workflow-linting-vuln-checks-ci-and-quality-gates) | — | 55 |
| Memory ordering | — | [Go §27](../Golang/real-life-golang-guide.md#27-the-go-memory-model-and-the-garbage-collector) | [C §38](../c-lang/real-life-c-guide.md#38-the-c11-memory-model-atomics-ordering-lock-free-basics) | 61 |

## 0.3 The Go ↔ C ↔ Rust phrasebook

Use this table when you know how to say something in Go or C and want the
Rust equivalent. Every row is explained properly later in the guide.

| Idea | Go | C | Rust |
|---|---|---|---|
| Who frees memory | the GC, eventually | you, with `free`, exactly once | the **owner**, when it goes out of scope (`Drop`), decided at compile time |
| Heap allocation | `new(T)`, `&T{}`, escape analysis | `malloc(sizeof(T))` | `Box::new(t)`, `Vec`, `String` |
| Pointer that may be null | `*T` (nil) | `T *` (NULL) | `Option<&T>` / `Option<Box<T>>`. A plain `&T` is never null |
| Read-only view | by convention | `const T *` | `&T` (shared borrow) |
| Exclusive writable access | by convention | by convention | `&mut T` (unique borrow, checked) |
| Growable array | `[]T` + `append` | hand-rolled vector | `Vec<T>` + `push` |
| View into an array | slice `s[a:b]` (can alias and grow) | `T *` + length | `&v[a..b]`, a fat pointer that cannot outlive or alias-write `v` |
| String | `string` (immutable bytes, UTF-8 by convention) | `char *`, NUL-terminated | `String` (owned) / `&str` (borrowed), **always valid UTF-8** |
| Map | `map[K]V` | hand-rolled hash table | `HashMap<K, V>`, `BTreeMap<K, V>` |
| Error | `(v, err)`, `if err != nil` | return code + `errno` | `Result<T, E>` + `?` |
| Unrecoverable failure | `panic` | `abort()`, or UB | `panic!` (unwinds by default) |
| Sum type / tagged union | interface + type switch | `struct { enum tag; union {...} }` | `enum` with data + `match` (exhaustive) |
| Polymorphism | interfaces (implicit) | function pointers, vtables by hand | traits: generics (static) or `dyn Trait` (vtable) |
| Generics | type parameters (GC shape stenciling) | macros / `void *` | generics, **monomorphized** like C++ templates |
| Cleanup | `defer` | `goto cleanup` | `Drop`, runs automatically, including on early return |
| Threads | goroutines (M:N, runtime) | `pthread_create` | `std::thread` (1:1 OS threads) |
| Lightweight tasks | goroutines | — (event loop by hand, `epoll`) | `async` + an executor such as Tokio |
| Channel | `chan T` | pipe or hand-rolled queue | `std::sync::mpsc`, `tokio::sync::mpsc` |
| Mutex | `sync.Mutex`, guards nothing in particular | `pthread_mutex_t` | `Mutex<T>`: the lock **owns** the data |
| Cancellation | `context.Context` | flag + `pthread_cancel` (rarely) | drop the future; `CancellationToken` |
| Data race | runtime race detector (`-race`) | UB | **compile error** (`Send`/`Sync`) |
| Escape hatch | `unsafe`, `cgo` | (the whole language) | `unsafe { }` blocks, `extern "C"` |
| Package manager | modules, `go.mod` | none (system packages, vendoring) | Cargo, `Cargo.toml` + `Cargo.lock` |
| Formatter / linter | `gofmt`, `go vet` | `clang-format`, `clang-tidy` | `rustfmt`, `clippy` |
| Tests | `go test`, `_test.go` | hand-rolled harness | `cargo test`, `#[test]`, doc tests |
| Vulnerability scan | `govulncheck` | — | `cargo audit` / `cargo deny` |

**The one-sentence version of each language's memory model:**

- **C:** every pointer is a promise you made, and the compiler believes it.
- **Go:** every pointer keeps its target alive, and the GC sorts it out.
- **Rust:** every value has one owner, and borrows cannot outlive it or
  conflict with it. The compiler checks both rules.

## 0.4 The Rust labs: what you build for each guide in the series

All labs live in one Cargo workspace, one binary per lab:

```bash
cargo new --vcs none rust-labs && cd rust-labs
mkdir -p src/bin                       # each lab is src/bin/<name>.rs
cargo run --bin wordfreq -- README.md  # run one lab
```

| Series guide | Rust labs in this guide |
|---|---|
| [OS & Linux](../os-linux/real-life-os-guide.md) | `rawsys` (Ch 39), `supervisor` and graceful SIGTERM (40), `atomicwrite` (41), `rps` from `/proc` (42), `countalloc` (43), `minibox` namespace container (44) |
| [Networking (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md) | `echo` threaded and async (32), length-prefixed framing (32), `udpping` (32), `dnsq` raw-bytes DNS client (45), `tinyhttp` (46) |
| [Security from Zero](../security/real-life-guide.md) | `seal` AES-256-GCM (52), Argon2id password store (52), HMAC webhook verify (52), `authsvc` JWT service (56) |
| [Security in Depth](../security/real-life-security-guide-v1.md) | `safefetch` SSRF guard (48), smuggling-safe proxy (50), fuzz target (54), `cargo deny` policy (55) |
| [HTTPS Lifecycle](../v2-https/real-life-guide-v1.md) | `api` production Axum server (47), retries with idempotency keys (48), TLS/mTLS servers (49), `revproxy` (50), `lifecycle` whole-series trace (64) |
| [Go guide](../Golang/real-life-golang-guide.md) | `lru` (27, compare Go §26), `chat` (38, compare Go §38), `kv` (63, compare Go §56) |
| [C guide](../c-lang/real-life-c-guide.md) | `pool` thread pool (37, compare C §23.1), FFI to a C library (34), `kv` RESP server (63, compare C §43) |

The crates used across the guide (versions as written):

```toml
[package]
name = "rust-labs"
version = "0.1.0"
edition = "2024"

[dependencies]
anyhow = "1"
thiserror = "2"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
tokio = { version = "1", features = ["full"] }
bytes = "1"
axum = "0.8"
tower-http = { version = "0.7", features = ["timeout", "limit", "trace"] }
hyper = "1"
hyper-util = "0.1"
http-body-util = "0.1"
tracing = "0.1"
tracing-subscriber = "0.3"
rustls = "0.23"
tokio-rustls = "0.26"
rcgen = "0.14"
aes-gcm = "0.11"
argon2 = "0.6"
hmac = "0.13"
sha2 = "0.11"
subtle = "2"
zeroize = "1"
rand = "0.10"
jsonwebtoken = { version = "11", features = ["rust_crypto"] }
libc = "0.2"
nix = { version = "0.31", features = ["process", "signal", "sched", "mount", "fs", "hostname", "feature"] }
memmap2 = "0.9"
signal-hook = "0.4"

[dev-dependencies]
proptest = "1"
```

---

# Part I — Foundations: the Rust way

Part I covers the language you need to read any Rust program. Chapters 1–5
move quickly if you know Go or C. Chapters 6–8 (ownership, borrowing, slices)
are where Rust differs from every language you know. Read those slowly, run
every snippet, and break them on purpose. Part I ends with a checkpoint and
two projects.

---

## 1. What Rust is, and why it exists

### 1.1 The problem being solved

For about forty years, systems software had two options:

1. **Manual memory management (C, C++).** Fast and predictable, with control
   over layout and no runtime. Every pointer is a promise the programmer
   made, and the compiler believes it. The [C guide's Part VIII](../c-lang/real-life-c-guide.md#44-memory-leaks-in-depth-every-shape-they-actually-take-in-production)
   shows how often those promises are broken.
2. **Garbage collection (Java, Go, C#).** Memory safe, but it needs a runtime,
   adds pauses and memory overhead, and gives up some control over layout.
   [Go §27](../Golang/real-life-golang-guide.md#27-the-go-memory-model-and-the-garbage-collector)
   shows how much engineering Go puts into keeping that cost small.

The cost of the first option is measured. Microsoft (2019) and the Chromium
project (2020) each reported that about **70% of their serious security bugs
were memory-safety bugs**: use-after-free, buffer overflow, double free, and
uninitialized reads. Those are exactly the bug classes in the
[C guide's UB list](../c-lang/real-life-c-guide.md#28-undefined-behavior-the-list-every-c-programmer-must-memorize).

Rust's claim is a third option: **memory safety and data-race freedom,
checked at compile time, with no garbage collector and no mandatory
runtime.** The compiler tracks who owns each value and who is borrowing it.
It rejects programs where those rules could be broken.

### 1.2 The design goals, and how they show up in the syntax

| Goal | Where you see it |
|---|---|
| Memory safety without GC | ownership, borrowing, lifetimes (Ch 6–8, 19) |
| No data races | `Send` and `Sync` traits; `Mutex<T>` owns its data (Ch 28–29) |
| Zero-cost abstractions | generics are monomorphized; iterators compile to the same loop you would write by hand (Ch 16, 21) |
| Explicitness | no implicit numeric conversions, no null, no exceptions, `mut` is opt-in |
| Make illegal states unrepresentable | enums with data plus exhaustive `match` (Ch 9) |
| Escape hatch with a fence around it | `unsafe` blocks: you can still do anything C can, in a marked, auditable place (Ch 33) |

### 1.3 What Rust deliberately does not have

- **No null.** Use `Option<T>`. The compiler forces you to handle `None`.
- **No exceptions.** Recoverable errors are values (`Result<T, E>`). Panics
  exist for bugs, not for control flow.
- **No inheritance.** Use traits and composition, as in Go.
- **No garbage collector.** Memory is freed at a point the compiler decides
  statically: when the owner goes out of scope.
- **No implicit numeric conversion.** `u8 + u32` is a compile error. Write
  `x as u32` or `u32::from(x)`.
- **No uninitialized variables** in safe code. Reading a variable before it
  is assigned is a compile error, not
  [C's UB](../c-lang/real-life-c-guide.md#28-undefined-behavior-the-list-every-c-programmer-must-memorize).

### 1.4 War story

In 2019 a common class of bug in a large C++ media pipeline looked like this:
a decoder kept a pointer into a buffer and a background thread reallocated
the buffer. The pointer now pointed into freed memory. The crash appeared in
roughly one in a million frames, in production only, with a stack trace that
pointed at innocent code. In Rust, keeping a `&[u8]` into a `Vec<u8>` while
another part of the program calls `push` on it is a **compile error**
(E0502, Ch 7). Sharing the `Vec` with another thread at all needs an
`Arc<Mutex<Vec<u8>>>`, or the compiler refuses. The bug is impossible to
write in safe Rust, so nobody has to find it.

**Lesson:** Rust does not make you a better programmer at runtime. It refuses
to compile a large class of programs that a careful reviewer would also
reject, so the careful review happens on every build.

```bash
cargo new hello && cd hello && cargo run
#    Compiling hello v0.1.0
#     Finished `dev` profile [unoptimized + debuginfo] target(s)
#      Running `target/debug/hello`
# Hello, world!
```

---

## 2. The toolchain: cargo, crates, editions, and what `cargo build` really does

### 2.1 The commands you will type every day

| Command | What it does |
|---|---|
| `cargo new app` / `cargo new --lib mylib` | create a binary or library crate |
| `cargo build` / `cargo build --release` | compile to `target/debug/` or `target/release/` |
| `cargo run -- args` | build and run; arguments after `--` go to your program |
| `cargo check` | type-check and borrow-check without generating code; the fastest feedback loop |
| `cargo test` | build and run unit, integration, and doc tests |
| `cargo clippy` | lints, about 800 of them; treat its advice as code review |
| `cargo fmt` | format; like `gofmt`, nobody argues about style |
| `cargo doc --open` | build HTML docs for your crate and every dependency |
| `cargo add serde --features derive` | add a dependency to `Cargo.toml` |
| `cargo tree` | show the dependency graph |
| `cargo expand` (plugin) | show code after macro expansion |

### 2.2 What `cargo build` actually does

The [C guide §2.1](../c-lang/real-life-c-guide.md#2-the-toolchain-the-four-stage-pipeline-and-your-first-program)
splits compilation into preprocess → compile → assemble → link. Rust's
pipeline has different stages, and the borrow checker sits in the middle:

```
 source (.rs)
   │  parse, expand macros            (no textual preprocessor; macros work on tokens)
   ▼
 AST ──▶ HIR                          (name resolution, type inference, trait resolution)
   ▼
 MIR                                   (the control-flow graph that the BORROW CHECKER runs on)
   │  borrow check, const eval, MIR optimizations
   ▼
 LLVM IR  ──▶  LLVM optimizer  ──▶  object files (.o)   ← one per "codegen unit"
   ▼
 linker (cc / ld / lld)  ──▶  executable, statically linked with Rust deps, dynamically with libc
```

Consequences you will notice:

- **The unit of compilation is the crate**, not the file. A crate is a
  library or binary built from a tree of modules. That is why `cargo check`
  type-checks your whole crate at once.
- **Generics are compiled per concrete type** (monomorphization), so
  generic-heavy code is fast but compiles slowly and produces bigger binaries.
- **Rust dependencies are linked statically** into your binary. `ldd
  target/release/app` on Linux shows only libc and friends.

```bash
cargo rustc --release -- --emit=asm      # see the assembly (target/release/deps/*.s)
cargo rustc -- -Zunpretty=mir            # MIR needs nightly; use https://play.rust-lang.org's "MIR" button
```

### 2.3 `Cargo.toml`, `Cargo.lock`, and editions

```toml
[package]
name = "hello"
version = "0.1.0"
edition = "2024"     # language edition: 2015, 2018, 2021, 2024

[dependencies]
serde = { version = "1", features = ["derive"] }   # "1" means ^1.0.0: any 1.x >= 1.0.0

[profile.release]
lto = "thin"         # link-time optimization across crates
codegen-units = 1    # slower build, faster binary
panic = "abort"      # smaller binary; no unwinding (see Ch 10.5)
```

- **`Cargo.lock`** records the exact version of every transitive dependency.
  Commit it for applications and for libraries (current Cargo guidance).
  It does for Rust what `go.sum` does for Go, and it is the input to
  `cargo audit` (Ch 55).
- **Editions** are opt-in language versions. They change syntax details, not
  the ABI. Crates from different editions link together freely, so the
  ecosystem never splits the way Python 2/3 did. Edition 2024 shipped in Rust
  1.85 (February 2025).

### 2.4 Cross-compilation

```bash
rustup target add x86_64-unknown-linux-musl     # fully static Linux binaries
cargo build --release --target x86_64-unknown-linux-musl
rustup target add aarch64-unknown-linux-gnu     # needs a cross linker; see Ch 66 for `cross`
```

Compare [Go §2.3](../Golang/real-life-golang-guide.md#2-the-toolchain-and-workspace)
(`GOOS=linux GOARCH=arm64 go build` simply works) and the
[C guide §41.1](../c-lang/real-life-c-guide.md#41-cross-compilation-and-build-systems-at-scale-cmake-meson)
(you need a full cross toolchain). Rust sits between them: the compiler can
target anything LLVM supports, but you still need a linker and a libc for
the target. Ch 66 covers `cross` and `cargo-zigbuild`, which solve this.

---

## 3. Variables, mutability, shadowing, and scalar types

### 3.1 Immutable by default

```rust
fn main() {
    let x = 5;          // immutable binding
    // x = 6;           // error[E0384]: cannot assign twice to immutable variable
    let mut y = 5;      // `mut` is opt-in, and greppable
    y += 1;

    let x = x * 2;      // SHADOWING: a brand-new variable that hides the old `x`
    let spaces = "   ";
    let spaces = spaces.len(); // shadowing can even change the type: &str -> usize

    println!("{x} {y} {spaces}"); // 10 6 3
}
```

`mut` belongs to the **binding**, not to the type. `let mut v = Vec::new()`
means "this variable may be reassigned or mutated through". Later you will
see the same distinction between `&T` and `&mut T`.

### 3.2 Scalar types, and the sizes that are actually guaranteed

| Rust | Size | Go | C |
|---|---|---|---|
| `i8 i16 i32 i64 i128` | exact | `int8 … int64` | `int8_t …` (`<stdint.h>`) |
| `u8 u16 u32 u64 u128` | exact | `uint8 … uint64` | `uint8_t …` |
| `isize usize` | pointer width | `int uint` | `ptrdiff_t size_t` |
| `f32 f64` | IEEE 754 | `float32 float64` | `float double` |
| `bool` | 1 byte, only 0 or 1 | `bool` | `_Bool` |
| `char` | **4 bytes**, a Unicode scalar value | `rune` | — (`char` is a byte) |

Unlike [C §3.1](../c-lang/real-life-c-guide.md#3-types-and-the-sizes-the-standard-actually-guarantees),
every integer type has a fixed size. **Indexes and lengths are always
`usize`**, so you will write `as usize` often at first.

### 3.3 Integer overflow: defined, but different in debug and release

In C, signed overflow is UB. In Go, it silently wraps. Rust's behavior is
defined in every case, and you choose it:

```rust
fn main() {
    let a: u8 = 250;

    // In a DEBUG build, `a + 10` panics: "attempt to add with overflow".
    // In a RELEASE build it wraps to 4 (two's complement), unless you set
    // `overflow-checks = true` in [profile.release].
    // When overflow matters, say what you want explicitly:
    println!("{:?}", a.checked_add(10));     // None
    println!("{}", a.wrapping_add(10));      // 4
    println!("{}", a.saturating_add(10));    // 255
    println!("{:?}", a.overflowing_add(10)); // (4, true)

    // Conversions are explicit. `as` truncates silently; `try_from` checks.
    let big: i64 = 300;
    let truncated = big as u8;                    // 44, silently
    let checked = u8::try_from(big);              // Err(TryFromIntError(()))
    println!("{truncated} {checked:?}");
}
```

**Real-world example.** Parsing a length field from a network packet
(Ch 45): `let len = u16::from_be_bytes([b0, b1]) as usize;` is safe because
`u16 → usize` can never lose data. Going the other way (`usize → u16` when
writing a packet) must use `u16::try_from(len)?`. With `as`, a 70,000-byte
payload would quietly be encoded as a 4,464-byte length, which is the
classic framing bug behind many C parser CVEs.

**Production setting:** many security-sensitive projects set
`overflow-checks = true` in `[profile.release]`. The cost is usually under
a few percent. Profile your own workload before deciding.

### 3.4 Tuples, arrays, and type inference

```rust
fn main() {
    let pair: (i32, &str) = (200, "OK");
    let (code, reason) = pair;               // destructuring
    println!("{code} {reason} {}", pair.0);

    let a = [0u8; 4];                        // [u8; 4], four zero bytes, on the stack
    let b = [1, 2, 3];                       // inferred [i32; 3]
    println!("{} {:?} {}", a.len(), b, b[2]);
    // b[3] would not compile here: the index is a constant, and the compiler knows the length.
    // For a runtime index, an out-of-bounds access PANICS. It never reads past the end.
    let i = std::env::args().count() + 5;
    println!("{:?}", b.get(i));             // None: the non-panicking accessor
}
```

Bounds checks are the reason Rust programs cannot have the
[C guide's §29 stack overflow](../c-lang/real-life-c-guide.md#29-buffer-overflows-safe-string-handling-and-cert-c-rules).
They usually cost little, because LLVM removes them when it can prove the
index is in range. Iterators (Ch 21) avoid them entirely.

---

## 4. Control flow: everything is an expression

### 4.1 `if`, `match`, and blocks produce values

```rust
fn classify(status: u16) -> &'static str {
    // `match` must be EXHAUSTIVE: the `_` arm is required here, because u16 has 65,536 values.
    match status {
        200 | 204 => "success",
        301..=308 => "redirect",
        400..=499 => "client error",
        500..=599 => "server error",
        _ => "unknown",
    }
}

fn main() {
    let n = 7;
    let parity = if n % 2 == 0 { "even" } else { "odd" }; // no ternary; `if` is an expression
    let squared = {
        let t = n * n;
        t + 0 // the last expression, with NO semicolon, is the block's value
    };
    println!("{parity} {squared} {}", classify(404));
}
```

### 4.2 Loops: `loop`, `while`, `for`, labels, and `break` with a value

```rust
fn main() {
    // `loop` is an infinite loop that can return a value through `break`.
    let mut attempts = 0;
    let port = loop {
        attempts += 1;
        let candidate = 8000 + attempts;
        if candidate % 3 == 0 {
            break candidate; // the loop's value
        }
    };

    // `for` iterates over anything that implements IntoIterator.
    for i in 0..3 { print!("{i} "); }        // 0 1 2   (0..=3 includes 3)
    for (i, ch) in "héllo".char_indices() { print!("[{i}:{ch}]"); }
    println!();

    // Labels break out of nested loops; no goto needed (compare C §4.2).
    'outer: for row in 0..10 {
        for col in 0..10 {
            if row * col == 42 {
                println!("found at {row},{col}");
                break 'outer;
            }
        }
    }
    println!("port={port} after {attempts} attempts");
}
```

### 4.3 `if let`, `let else`, and `while let`

These are shortcuts for a `match` where you care about one pattern:

```rust
fn parse_port(s: &str) -> u16 {
    // let-else: bind on success, or the else block must leave the function (return/break/panic).
    let Ok(port) = s.parse::<u16>() else {
        return 8080;
    };
    port
}

fn main() {
    let maybe_user: Option<&str> = Some("ada");
    if let Some(name) = maybe_user {
        println!("hello {name}");
    }

    let mut stack = vec![1, 2, 3];
    while let Some(top) = stack.pop() { // runs until pop() returns None
        print!("{top} ");
    }
    println!("{} {}", parse_port("443"), parse_port("http"));
}
```

`let else` is the Rust version of Go's `if err != nil { return }` guard
clause. It keeps the main path unindented.

---

## 5. Functions, statements, and the unit and never types

```rust
// Parameter and return types are ALWAYS written; inference stays inside function bodies.
fn area(w: u32, h: u32) -> u32 {
    w * h // no semicolon: this is the return value
}

fn log(msg: &str) {          // no `-> T` means the return type is `()`, the unit type
    eprintln!("[log] {msg}");
}

fn fail(msg: &str) -> ! {    // `!` is the NEVER type: this function never returns
    panic!("fatal: {msg}");
}

fn main() {
    log("starting");
    let a = area(3, 4);
    // Because `!` converts to any type, diverging calls fit in any expression:
    let b: u32 = if a > 0 { a } else { fail("zero area") };
    println!("{b}");
}
```

Statements end in `;` and produce `()`. Expressions produce values. The most
common beginner error is a stray semicolon on the last line of a function:

```rust
// DOES NOT COMPILE
fn area(w: u32, h: u32) -> u32 {
    w * h;   // statement, so the body evaluates to ()
}
fn main() { println!("{}", area(1, 2)); }
```

```text
error[E0308]: mismatched types
  |     fn area(w: u32, h: u32) -> u32 {
  |        ----                    ^^^ expected `u32`, found `()`
  |         w * h;
  |              - help: remove this semicolon to return this value
```

Rust has **no function overloading and no default arguments**. Use
different names (`new`, `with_capacity`), `Option` parameters, or a builder
(Ch 62).

---

## 6. Ownership: moves, copies, and drop

This chapter is the core of the language. Read it slowly.

### 6.1 The three rules

1. Each value has exactly **one owner** (a variable, a struct field, a `Vec` slot...).
2. When the owner goes out of scope, the value is **dropped**: its
   destructor runs and its heap memory is freed.
3. Ownership can be **moved** to a new owner. After a move, the old owner
   cannot be used.

That is the complete memory-management model. There is no GC, and no `free`
for you to forget or call twice.

### 6.2 What is on the stack and what is on the heap

```rust
fn main() {
    let n: u64 = 42;                       // 8 bytes, on the stack
    let s = String::from("hello");         // 24 bytes on the stack: (ptr, len, cap)
                                           // + 5 bytes on the heap: h e l l o
    let v: Vec<u32> = vec![1, 2, 3];       // same shape: (ptr, len, cap) -> heap [1,2,3]
    println!("{n} {s} {v:?}");
    println!("{}", std::mem::size_of::<String>()); // 24 on a 64-bit target
}   // <- `v` is dropped, then `s` (reverse declaration order): both heap buffers are freed here.
```

```
 stack frame of main                 heap
 ┌─────────────────────┐
 │ n   = 42            │
 │ s   ptr ────────────┼──────────▶  [h][e][l][l][o]
 │     len = 5         │
 │     cap = 5         │
 │ v   ptr ────────────┼──────────▶  [1][2][3]
 │     len = 3, cap = 3│
 └─────────────────────┘
```

This is the same layout as the C guide's
[dynamic array](../c-lang/real-life-c-guide.md#16-mini-projects-a-dynamic-array-library-csv-parser-and-a-json-parser)
and Go's [slice header](../Golang/real-life-golang-guide.md#6-arrays-and-slices-the-internals-that-explain-the-bugs).
The difference is what happens when you copy the stack part.

### 6.3 Move: copying the header transfers ownership

```rust
// DOES NOT COMPILE
fn main() {
    let s1 = String::from("hello");
    let s2 = s1;              // the 24-byte header is copied; ownership MOVES to s2
    println!("{s1} {s2}");    // s1 is no longer valid
}
```

```text
error[E0382]: borrow of moved value: `s1`
  |     let s1 = String::from("hello");
  |         -- move occurs because `s1` has type `String`, which does not implement the `Copy` trait
  |     let s2 = s1;
  |              -- value moved here
  |     println!("{s1} {s2}");
  |               ^^^^ value borrowed here after move
  = help: consider cloning the value if the performance cost is acceptable: `.clone()`
```

Why the compiler rejects this: if both `s1` and `s2` were valid, both would
free the same heap buffer when they went out of scope. That is a
[double free, C guide §10.2](../c-lang/real-life-c-guide.md#10-dynamic-memory-the-malloc-family-ownership-and-the-bugs-it-enables).
Go avoids it with a GC that keeps the buffer alive for both. Rust avoids it
by allowing only one owner.

A move is only a `memcpy` of the stack part (24 bytes here). The heap data
does not move. In optimized code the copy usually disappears too.

Passing a value to a function moves it, and returning a value moves it out:

```rust
fn consume(s: String) -> usize { s.len() }    // `s` is dropped at the end of consume
fn make() -> String { String::from("fresh") } // ownership moves out to the caller

fn main() {
    let a = String::from("data");
    let n = consume(a);          // `a` moved into consume; cannot use `a` after this line
    let b = make();
    println!("{n} {b}");
}
```

### 6.4 `Copy` and `Clone`

- **`Copy`** types are duplicated bit for bit on assignment, and the original
  stays valid. These are the types with no heap ownership: integers, floats,
  `bool`, `char`, shared references `&T`, and tuples or arrays of `Copy`
  types. A type that implements `Drop` can never be `Copy`.
- **`Clone`** is an explicit, possibly expensive deep copy: `s.clone()`
  allocates a new heap buffer.

```rust
#[derive(Debug, Clone, Copy)]
struct Point { x: i32, y: i32 }   // all fields are Copy, so the struct can derive Copy

fn main() {
    let p1 = Point { x: 1, y: 2 };
    let p2 = p1;                   // a copy, not a move
    println!("{p1:?} {p2:?}");     // both still valid

    let s1 = String::from("hi");
    let s2 = s1.clone();           // explicit deep copy: a second heap allocation
    println!("{s1} {s2}");
}
```

Cloning is visible in the code, so it is easy to find in review. When you see
`.clone()` in a hot loop, it is often the first thing to optimize (Ch 68).

### 6.5 Drop: deterministic destruction (RAII)

```rust
struct TempFile { path: String }

impl Drop for TempFile {
    fn drop(&mut self) {
        // Runs exactly once, when the owner goes out of scope, including on early
        // return or panic unwind. Rust's answer to Go's `defer` and C's `goto cleanup`.
        println!("removing {}", self.path);
        let _ = std::fs::remove_file(&self.path);
    }
}

fn work(fail: bool) -> Result<(), String> {
    let _tmp = TempFile { path: "/tmp/rust-guide-demo".into() };
    std::fs::write("/tmp/rust-guide-demo", b"scratch").map_err(|e| e.to_string())?;
    if fail {
        return Err("bailed out early".into()); // _tmp is still dropped here
    }
    println!("finished work");
    Ok(())
} // ...and here on the normal path

fn main() {
    println!("{:?}", work(false));
    println!("{:?}", work(true));
}
```

Files, sockets, mutex guards, database transactions, and TLS sessions all
release their resources in `Drop`. The C guide's
[§44.5 "every acquire needs a matching release"](../c-lang/real-life-c-guide.md#44-memory-leaks-in-depth-every-shape-they-actually-take-in-production)
is automatic in Rust: you cannot forget the release, and you cannot run it
twice.

**Watch the binding name.** `let _guard = lock()` keeps the guard alive until
the end of the scope. `let _ = lock()` drops it **immediately**, because `_`
is not a binding. This is a real source of "my lock does nothing" bugs.

**Leaks are still possible, and they are safe.** `Rc` reference cycles
(Ch 22) and `std::mem::forget` both leak memory. Rust guarantees no
use-after-free, not no leaks. Ch 51 explains why.

### 6.6 War story

A team ported a C daemon to Rust. The C code had a
`conn_close(conn); log_conn(conn);` sequence that read a freed struct on one
error path. It had shipped for six years, because the freed memory usually
still held the old bytes. The Rust port could not compile the same sequence:
`conn.close()` took `self` by value, so `conn` was moved, and the logging
line failed with E0382. The bug was found when the code was translated, not
in production. **Lesson: when an API consumes a value (`fn close(self)`),
the compiler stops anyone from using that value afterwards. That is the
typestate idea from Ch 62, in its simplest form.**

---

## 7. Borrowing: shared XOR mutable

Moving values everywhere would be impractical. Usually you want to let code
**use** a value without taking ownership. That is a **borrow**, written as a
reference.

### 7.1 The rule

At any moment, a value can have **either**:

- any number of **shared references** `&T` (read-only), **or**
- exactly **one mutable reference** `&mut T` (read-write),

and **no reference may outlive the value it points to.**

```rust
fn total_len(words: &[String]) -> usize {      // borrows; the caller keeps ownership
    words.iter().map(|w| w.len()).sum()
}

fn shout(words: &mut Vec<String>) {           // exclusive borrow; may modify
    for w in words.iter_mut() {
        *w = w.to_uppercase();
    }
    words.push("!".to_string());
}

fn main() {
    let mut words = vec!["hello".to_string(), "world".to_string()];
    let n = total_len(&words);   // shared borrow, ends when the call returns
    shout(&mut words);           // mutable borrow, ends when the call returns
    println!("{n} {words:?}");   // owner still valid: ["HELLO", "WORLD", "!"]
}
```

### 7.2 Why the rule exists: iterator invalidation

```rust
// DOES NOT COMPILE
fn main() {
    let mut v = vec![1, 2, 3];
    let first = &v[0];       // shared borrow into v's heap buffer
    v.push(4);               // needs &mut v, and push may REALLOCATE the buffer
    println!("{first}");     // first would point into freed memory
}
```

```text
error[E0502]: cannot borrow `v` as mutable because it is also borrowed as immutable
  |     let first = &v[0];
  |                  - immutable borrow occurs here
  |     v.push(4);
  |     ^^^^^^^^^ mutable borrow occurs here
  |     println!("{first}");
  |               ------- immutable borrow later used here
```

This is the same bug as Go's
[append aliasing trap](../Golang/real-life-golang-guide.md#6-arrays-and-slices-the-internals-that-explain-the-bugs)
and C's realloc-then-dangling-pointer bug. In Go it silently corrupts data.
In C it is UB. In Rust it does not compile.

"Shared XOR mutable" rules out a whole category of bugs at once:

| Bug | Why it cannot happen |
|---|---|
| Iterator invalidation | you cannot mutate a collection while a borrow of it is alive |
| Use-after-free via a reference | a reference cannot outlive its owner |
| Data races | a data race needs two threads with access, at least one writing; that is shared AND mutable at once (Ch 28) |
| Aliasing surprises for the optimizer | `&mut T` is guaranteed unique, so LLVM gets `noalias` information C only gets from `restrict` |

### 7.3 Borrows last until their last use (non-lexical lifetimes)

```rust
fn main() {
    let mut scores = vec![10, 20, 30];
    let max = scores.iter().max().copied();   // `.copied()` turns Option<&i32> into Option<i32>
    // The shared borrow from iter() ended at the end of the line above,
    // because `max` holds a copy, not a reference.
    scores.push(40);
    println!("{max:?} {scores:?}");

    let first = &scores[0];
    println!("{first}");       // last use of `first`: the borrow ends here
    scores.push(50);           // so this is fine
}
```

Most of the friction with the borrow checker goes away once you remember
that a borrow ends at its last use, not at the end of the block.

### 7.4 Dereferencing and auto-ref

```rust
fn bump(n: &mut i32) {
    *n += 1;                       // `*` dereferences: write through the reference
}

fn main() {
    let mut count = 0;
    bump(&mut count);
    let s = String::from("abc");
    let r = &s;
    println!("{} {}", r.len(), (*r).len()); // method calls auto-dereference: r.len() just works
    println!("{count}");
}
```

The `.` operator adds `&`, `&mut`, or `*` as needed to make a method call
type-check. That is why you rarely write `*` outside assignments.

### 7.5 The four ways to fix a borrow-checker error

When you hit E0499, E0502, or E0505, one of these almost always works:

1. **Shorten the borrow.** Copy or clone the small value you need
   (`let id = user.id;`) before the mutation.
2. **Reorder.** Finish reading before you start writing.
3. **Split the borrow.** Borrow disjoint fields
   (`let (a, b) = (&mut s.left, &mut s.right);`) or slices
   (`v.split_at_mut(mid)`).
4. **Restructure ownership.** Use indices instead of references (Ch 27), an
   entry API (Ch 11), or move the data instead of borrowing it.

Wrapping everything in `Rc<RefCell<T>>` also works, but it moves the checks
to runtime. Treat it as a last resort (Ch 22, 68).

---

## 8. Slices and strings: fat pointers and UTF-8

### 8.1 Slices: a pointer and a length

A slice `&[T]` is a borrowed view into contiguous elements. It is a **fat
pointer**: two words, (pointer, length). Unlike a Go slice it has no
capacity, so it cannot grow and cannot write past its end.

```rust
fn sum(xs: &[i64]) -> i64 { xs.iter().sum() }   // accepts arrays, Vecs, and sub-slices

fn main() {
    let arr = [1i64, 2, 3, 4, 5];
    let v = vec![10i64, 20, 30];
    println!("{} {} {}", sum(&arr), sum(&v), sum(&arr[1..3])); // 15 60 5

    let mut buf = [0u8; 8];
    let (head, tail) = buf.split_at_mut(4);   // two non-overlapping &mut slices: allowed
    head[0] = 0xAA;
    tail[0] = 0xBB;
    println!("{buf:02x?}");
    println!("{}", std::mem::size_of::<&[u8]>()); // 16: pointer + length
}
```

Write functions that take `&[T]` instead of `&Vec<T>`, and `&str` instead
of `&String`. They accept more kinds of input and give up nothing. Clippy
warns about the other form (`ptr_arg`).

### 8.2 `String` vs `&str`

| | `String` | `&str` |
|---|---|---|
| Owns its bytes? | yes (heap) | no (borrowed view) |
| Layout | (ptr, len, cap) | (ptr, len) |
| Growable? | yes | no |
| Literal | `String::from("x")`, `"x".to_string()` | `"x"`, which is a `&'static str` stored in the binary |
| Guarantee | **always valid UTF-8** | **always valid UTF-8** |

```rust
fn greet(name: &str) -> String {
    format!("hello, {name}")       // format! allocates a new String
}

fn main() {
    let owned: String = String::from("Ada");
    let borrowed: &str = &owned;   // &String coerces to &str (Deref, Ch 18)
    let lit: &'static str = "Lovelace";
    let mut full = greet(borrowed);
    full.push(' ');
    full.push_str(lit);
    println!("{full}");
}
```

### 8.3 UTF-8: why `s[0]` does not compile

```rust
fn main() {
    let s = "héllo";                 // 'é' takes 2 bytes in UTF-8
    println!("{}", s.len());         // 6 BYTES, not 5 characters
    println!("{}", s.chars().count()); // 5 chars (Unicode scalar values)
    // let c = s[0];                 // error[E0277]: the type `str` cannot be indexed by `{integer}`
    println!("{}", &s[0..1]);        // "h": byte-range slicing is allowed...
    // println!("{}", &s[1..2]);     // ...but this PANICS: byte 2 is inside 'é'
    println!("{:?}", s.get(1..2));   // None: the non-panicking version
    for (i, c) in s.char_indices() { print!("{i}:{c} "); }
    println!();
    println!("{:?}", s.as_bytes());  // the raw bytes, when you need them
}
```

Compare [Go §12](../Golang/real-life-golang-guide.md#12-strings-runes-bytes-and-unicode):
Go strings are UTF-8 by convention and can hold invalid bytes. Rust strings
are UTF-8 by **invariant**. `String::from_utf8(bytes)` returns a `Result`,
so invalid input is rejected at the boundary. For bytes that are not
necessarily text, such as file contents, network payloads, or paths, use
`Vec<u8>`/`&[u8]`, `OsString`, or `PathBuf`.

**Across the series:** Unicode normalisation bugs, for example two different
byte sequences that render as the same username, are covered in
[Security from Zero Ch 7](../security/real-life-guide.md#chapter-7-unicode-normalisation-and-the-bugs-they-cause).
Valid UTF-8 does not make strings normalised. Handle that separately with
the `unicode-normalization` crate.

---

## 9. Structs, enums, and pattern matching

### 9.1 Structs and `impl` blocks

```rust
#[derive(Debug, Clone, PartialEq)]
struct Endpoint {
    host: String,
    port: u16,
    tls: bool,
}

impl Endpoint {
    // An associated function (no `self`): Rust's constructor convention is `new`.
    fn new(host: impl Into<String>, port: u16) -> Self {
        Self { host: host.into(), port, tls: port == 443 }
    }
    // &self: read-only method. &mut self: mutating. self: consuming.
    fn url(&self) -> String {
        let scheme = if self.tls { "https" } else { "http" };
        format!("{scheme}://{}:{}", self.host, self.port)
    }
    fn with_tls(mut self) -> Self { self.tls = true; self }
}

struct Meters(f64);        // tuple struct: the NEWTYPE pattern (Ch 53, 62)
struct Marker;             // unit struct: zero bytes

fn main() {
    let e = Endpoint::new("example.com", 443);
    let local = Endpoint { host: "localhost".into(), ..e.clone() }; // struct update syntax
    let dev = Endpoint::new("dev.local", 8443).with_tls();
    println!("{} {} {}", e.url(), local.url(), dev.url());
    let Meters(m) = Meters(3.5);
    println!("{m} {}", std::mem::size_of::<Marker>());
}
```

Receiver choice is the Rust version of
[Go's value vs. pointer receiver rule](../Golang/real-life-golang-guide.md#8-structs-and-methods-value-vs-pointer-receivers-embedding),
with a third option. `self` by value **consumes** the struct, which builders
and state machines use.

### 9.2 Enums are sum types

A Rust `enum` is a tagged union, the
[C guide's §8.3 pattern](../c-lang/real-life-c-guide.md#8-structs-unions-enums-typedef-and-the-preprocessor),
except the compiler checks the tag for you:

```rust
#[derive(Debug)]
enum Message {
    Ping,                                   // no data
    Data { stream: u32, payload: Vec<u8> }, // named fields
    Close(u16, String),                     // tuple-like
}

fn handle(m: &Message) -> String {
    match m {
        Message::Ping => "pong".into(),
        Message::Data { stream, payload } if payload.is_empty() => format!("empty frame on {stream}"),
        Message::Data { stream, payload } => format!("{} bytes on {stream}", payload.len()),
        Message::Close(code, reason) => format!("close {code}: {reason}"),
    }
    // Remove any arm and you get error[E0004]: non-exhaustive patterns.
}

fn main() {
    let msgs = [
        Message::Ping,
        Message::Data { stream: 3, payload: vec![1, 2, 3] },
        Message::Data { stream: 5, payload: vec![] },
        Message::Close(1000, "bye".into()),
    ];
    for m in &msgs { println!("{}", handle(m)); }
}
```

Exhaustiveness is the main benefit. When you add a `Message::Resume`
variant, every `match` that does not handle it stops compiling, so the
compiler gives you the list of places to update.

### 9.3 `Option<T>` instead of null

```rust
fn find_user(id: u32) -> Option<&'static str> {
    match id {
        1 => Some("ada"),
        2 => Some("grace"),
        _ => None,
    }
}

fn main() {
    // You cannot use an Option<&str> as a &str. Handle the None case first.
    let name = find_user(3).unwrap_or("anonymous");
    let upper = find_user(1).map(|n| n.to_uppercase());       // Some("ADA")
    let len = find_user(2).map_or(0, |n| n.len());            // 5
    let both = find_user(1).zip(find_user(2));                // Some(("ada","grace"))
    println!("{name} {upper:?} {len} {both:?}");

    // Option<&T> costs no extra space: the null pointer is used as the None "niche".
    println!("{}", std::mem::size_of::<Option<&u64>>()); // 8, same as &u64 (Ch 57)
}
```

Tony Hoare called null his "billion-dollar mistake". Go's nil-pointer panic
and C's NULL dereference are two versions of it. `Option` costs nothing at
runtime for references and boxes. The only cost is that you have to handle
`None` before you use the value.

### 9.4 Patterns everywhere

```rust
fn main() {
    let packet = (4u8, [192u8, 168, 1, 10], 443u16);
    match packet {
        (4, [10, ..], _) => println!("private 10/8"),
        (4, [192, 168, ..], port @ (80 | 443)) => println!("LAN web on {port}"),
        (4, [a, b, c, d], p) => println!("{a}.{b}.{c}.{d}:{p}"),
        (v, ..) => println!("IPv{v}"),
    }

    let point = (3, -3);
    if let (x, y) = point && x == -y {          // let-chains: edition 2024
        println!("on the anti-diagonal");
    }
}
```

`..` skips the remaining elements, `@` binds a value while testing it, `|`
matches alternatives, and `if` adds a guard. Patterns also work in `let`,
function parameters, `for`, and closures.

**Real-world example.** HTTP routers, protocol state machines, and
compilers' AST passes are where Rust's enums are most useful. Ch 63's RESP
parser is one `enum Frame` plus a `match`.

---

## 10. Error handling: `Result`, `?`, and panics

### 10.1 `Result<T, E>`

```rust
use std::num::ParseIntError;

fn parse_kv(line: &str) -> Result<(String, u32), String> {
    let (k, v) = line.split_once('=').ok_or(format!("missing '=' in {line:?}"))?;
    let v = v.trim();
    let v: u32 = v.parse().map_err(|e: ParseIntError| format!("bad value {v:?}: {e}"))?;
    Ok((k.trim().to_string(), v))
}

fn main() {
    for line in ["workers = 8", "timeout", "retries = many"] {
        match parse_kv(line) {
            Ok((k, v)) => println!("{k} -> {v}"),
            Err(e) => println!("error: {e}"),
        }
    }
}
```

The **`?` operator** means "if this is `Err(e)`, return `Err(e.into())` from
the current function, otherwise unwrap the `Ok` value". It does the same job
as Go's `if err != nil { return err }`, in one character, and it calls
`From::from` to convert the error type on the way out.

| | Go | C | Rust |
|---|---|---|---|
| Signal failure | `return nil, err` | `return -1; errno = EINVAL` | `return Err(e)` |
| Propagate | `if err != nil { return err }` | `if (rc < 0) goto fail;` | `?` |
| Can you ignore it? | yes, `v, _ := f()` | yes, very easily | `#[must_use]` warning; you must write `let _ =` explicitly |
| Wrap with context | `fmt.Errorf("...: %w", err)` | — | `.context("...")` (anyhow), or a `source` field (thiserror) |
| Inspect | `errors.Is` / `errors.As` | compare `errno` | `match` on the enum, `downcast_ref` |

### 10.2 Custom error types with `thiserror` (libraries)

```rust
use thiserror::Error;

#[derive(Debug, Error)]
pub enum ConfigError {
    #[error("cannot read config file {path}")]
    Io { path: String, #[source] source: std::io::Error },
    #[error("line {line}: {msg}")]
    Syntax { line: usize, msg: String },
    #[error("unknown key {0:?}")]
    UnknownKey(String),
}

fn load(path: &str) -> Result<String, ConfigError> {
    std::fs::read_to_string(path).map_err(|source| ConfigError::Io { path: path.into(), source })
}

fn main() {
    match load("/definitely/missing.toml") {
        Err(ConfigError::Io { path, source }) => println!("io error on {path}: {source}"),
        Err(e) => println!("{e}"),
        Ok(s) => println!("{} bytes", s.len()),
    }
    println!("{}", ConfigError::Syntax { line: 3, msg: "expected '='".into() });
    println!("{}", ConfigError::UnknownKey("colour".into()));
}
```

### 10.3 `anyhow` for applications

```rust
use anyhow::{bail, Context, Result};

fn read_port(path: &str) -> Result<u16> {
    let text = std::fs::read_to_string(path)
        .with_context(|| format!("reading port file {path}"))?;
    let port: u16 = text.trim().parse().context("port file is not a number")?;
    if port < 1024 { bail!("refusing privileged port {port}"); }
    Ok(port)
}

fn main() {
    if let Err(e) = read_port("/nonexistent/port") {
        // {:#} prints the whole chain: "reading port file ...: No such file or directory (os error 2)"
        println!("{e:#}");
    }
}
```

**The rule most Rust teams use:** libraries define precise error enums with
`thiserror`, so callers can `match` on them. Binaries use `anyhow::Result`
and add `.context(...)` at each layer, so the final message reads like a
stack of reasons.

### 10.4 Panics are for bugs

`panic!`, `unwrap()`, `expect()`, out-of-bounds indexing, and integer
overflow in debug builds all panic. By default a panic **unwinds** the
thread's stack and runs `Drop` for every live value. If it reaches the top
of `main`, the process exits with code 101.

```rust
fn main() {
    let config_from_build: Option<&str> = Some("prod");
    // expect() documents WHY you believe this cannot fail. Prefer it to unwrap().
    let env = config_from_build.expect("build script always sets the environment");

    let r = std::panic::catch_unwind(|| {
        let v: Vec<i32> = Vec::new();
        v[0] // index out of bounds: panic
    });
    println!("{env}, caught panic: {}", r.is_err());
    // The default panic hook still prints "index out of bounds: the len is 0 but
    // the index is 0" to stderr. catch_unwind stops the unwinding, not the message.
}
```

`catch_unwind` is for boundaries, such as a thread pool keeping a worker
alive or an FFI function that must not unwind into C. It is not a
`try`/`catch` for normal error handling.

When to panic: a broken invariant (a bug), tests, prototypes, and startup
code that cannot continue. When to return `Result`: anything caused by input,
the environment, or the network. A panic in a request handler is a
**denial-of-service bug** (Ch 51). Axum and Tokio isolate panics per task,
but you should not rely on that.

### 10.5 `panic = "abort"` and FFI

With `panic = "abort"` in the profile, a panic kills the process instead of
unwinding. You get smaller binaries and no `catch_unwind`. **A panic must
never unwind across an `extern "C"` boundary** into C code. Since Rust 1.81
that aborts the process (Ch 34).

---

## 11. Collections: `Vec`, `HashMap`, and friends

| Need | Type | Go equivalent |
|---|---|---|
| growable array | `Vec<T>` | `[]T` |
| double-ended queue, ring buffer | `VecDeque<T>` | — (container/list) |
| hash map / set | `HashMap<K,V>`, `HashSet<T>` | `map[K]V`, `map[T]struct{}` |
| sorted map / set | `BTreeMap<K,V>`, `BTreeSet<T>` | — (sort the keys) |
| priority queue | `BinaryHeap<T>` (max-heap) | `container/heap` |

### 11.1 `Vec` growth and capacity

```rust
fn main() {
    let mut v: Vec<u32> = Vec::with_capacity(4); // one allocation up front, like make([]T, 0, n)
    for i in 0..10 {
        v.push(i);
        print!("{}/{} ", v.len(), v.capacity()); // capacity at least doubles when it runs out
    }
    println!();
    v.retain(|&x| x % 2 == 0);    // filter in place
    v.dedup();                    // remove consecutive duplicates
    v.truncate(3);
    v.shrink_to_fit();            // give memory back
    println!("{v:?} cap={}", v.capacity());
    let drained: Vec<u32> = v.drain(..1).collect();
    println!("{drained:?} {v:?}");
}
```

### 11.2 `HashMap` and the entry API

```rust
use std::collections::HashMap;

fn main() {
    let log = "GET /a\nGET /b\nPOST /a\nGET /a";
    let mut hits: HashMap<&str, u32> = HashMap::new();
    for line in log.lines() {
        let path = line.split_whitespace().nth(1).unwrap_or("?");
        *hits.entry(path).or_insert(0) += 1; // ONE hash lookup, no double borrow
    }
    let mut sorted: Vec<_> = hits.iter().collect();
    sorted.sort_by(|a, b| b.1.cmp(a.1).then(a.0.cmp(b.0)));
    println!("{sorted:?}"); // [("/a", 3), ("/b", 1)]

    if let Some(n) = hits.get_mut("/b") { *n += 10; }
    println!("{:?} {:?}", hits.get("/b"), hits.get("/zzz"));
}
```

The entry API exists because of borrowing. The obvious version,
"`if !map.contains_key(k) { map.insert(k, 0) }` then `get_mut`", does two
lookups. Versions that hold a reference from the first lookup while
inserting do not compile. `entry()` does one lookup and returns a handle to
the slot.

**Security note:** `HashMap` uses SipHash-1-3 with a random key per process
by default. That makes it resistant to **hash-flooding DoS**, where an
attacker sends many keys that collide. Go's maps are also randomly seeded.
C hash tables like the
[C guide's djb2 table](../c-lang/real-life-c-guide.md#18-data-structures-from-scratch-linked-lists-stacks-queues-hash-tables-trees)
are not. For trusted keys in hot code, a faster hasher (`ahash`, `FxHash`)
can be several times quicker. Never use those for keys an attacker controls.

---

## 12. Modules, crates, and visibility

```
myapp/
├── Cargo.toml
└── src/
    ├── main.rs          // crate root of the binary
    ├── lib.rs           // crate root of the library (a package can have both)
    ├── config.rs        // mod config;
    └── net/
        ├── mod.rs       // or net.rs next to a net/ directory
        └── frame.rs     // mod frame; inside net
```

```rust
mod net {
    pub mod frame {
        pub struct Frame { pub len: u32, checksum: u32 } // `checksum` is private to this module

        impl Frame {
            pub fn new(len: u32) -> Self { Self { len, checksum: len ^ 0xFFFF } }
            pub fn valid(&self) -> bool { self.checksum == self.len ^ 0xFFFF }
        }

        pub(crate) fn internal_helper() -> u8 { 7 } // visible anywhere in THIS crate only
    }
    pub use frame::Frame; // re-export: callers write net::Frame
}

use net::Frame;

fn main() {
    let f = Frame::new(42);
    // let bad = net::frame::Frame { len: 1, checksum: 0 }; // error[E0451]: field `checksum` is private
    println!("{} {} {}", f.len, f.valid(), net::frame::internal_helper());
}
```

- Everything is **private by default**. Use `pub`, `pub(crate)`, or
  `pub(super)` to widen visibility. Go uses capitalization for the same
  purpose ([Go §11.1](../Golang/real-life-golang-guide.md#11-packages-and-modules));
  C uses `static` ([C §13.1](../c-lang/real-life-c-guide.md#13-multi-file-programs-headers-static-translation-units-and-make)).
- A private field means outside code cannot construct the struct with a
  literal. It has to go through your constructor, which can enforce
  invariants. This is how `String` guarantees valid UTF-8.
- **Features** (`[features]` in `Cargo.toml`) are compile-time flags,
  like C's `#ifdef`, applied per crate: `#[cfg(feature = "tls")]`.

---

## 13. 🔎 Checkpoint: ownership and borrowing

Try each one before you open the answer.

**1. Predict: does it compile? If not, which line fails?**

```rust
// DOES NOT COMPILE
fn main() {
    let names = vec![String::from("a"), String::from("b")];
    for n in names {
        println!("{n}");
    }
    println!("{}", names.len());
}
```

<details><summary>Answer</summary>

The last line fails with E0382 (borrow of moved value). `for n in names`
calls `names.into_iter()`, which **consumes** the vector, so each `n` is an
owned `String`. Write `for n in &names` to iterate over `&String` and keep
`names` usable.
</details>

**2. Spot the bug:**

```rust
// DOES NOT COMPILE
fn longest(words: &Vec<String>) -> &String {
    let mut best = &String::new();
    for w in words { if w.len() > best.len() { best = w; } }
    best
}
fn main() { println!("{}", longest(&vec!["hi".into()])); }
```

<details><summary>Answer</summary>

`&String::new()` creates a temporary owned by `longest` itself. Rust
extends its life to the end of the function body, but no further. If
`words` is empty, `best` still points at that temporary when the function
returns, so the compiler rejects the `best` return with E0515 ("cannot return value
referencing temporary value"). The signature promises a reference into
`words`, and this path breaks that promise. Fix: start from
`words.first()`, return `Option<&str>` so the empty case has an answer,
and take `&[String]`:

```rust
fn longest(words: &[String]) -> Option<&str> {
    words.iter().map(String::as_str).max_by_key(|w| w.len())
}
fn main() { println!("{:?}", longest(&["hi".into(), "hello".into()])); }
```

(When several words share the maximum length, `max_by_key` returns the
last one.)
</details>

**3. Write it yourself:** `fn dedupe_sorted(v: &mut Vec<i32>)` that removes
duplicates from an unsorted vector in place, without cloning the vector.

<details><summary>One correct answer</summary>

```rust
fn dedupe_sorted(v: &mut Vec<i32>) {
    v.sort_unstable();
    v.dedup();
}
fn main() {
    let mut v = vec![3, 1, 3, 2, 1];
    dedupe_sorted(&mut v);
    println!("{v:?}"); // [1, 2, 3]
}
```

To keep the original order instead, use `retain` with a `HashSet` of
values you have already seen: `let mut seen = HashSet::new(); v.retain(|x| seen.insert(*x));`.
</details>

**4. Conceptual:** why is `let s2 = s1;` a move for `String` but a copy for
`(i32, bool)`?

<details><summary>Answer</summary>

`String` owns a heap buffer and implements `Drop`. Two bitwise copies would
both free that buffer, so the language transfers ownership instead.
`(i32, bool)` owns nothing and implements `Copy`, so a bitwise copy is a
complete, independent value. The rule is: **a type is `Copy` only if
copying its bytes is enough to duplicate it.**
</details>

**5. Predict the output:**

```rust
fn main() {
    let mut s = String::from("ab");
    let r = &mut s;
    r.push('c');
    s.push('d');
    println!("{s}");
}
```

<details><summary>Answer</summary>

`abcd`. It compiles because the mutable borrow `r` is last used at
`r.push('c')`. After that line the borrow is over, so `s` can be used
directly again (non-lexical lifetimes, §7.3). If you added `println!("{r}")`
at the end, it would fail with E0499.
</details>

---

## 14. Terminal project: a word-frequency counter

The same project as [C §9.2](../c-lang/real-life-c-guide.md#9-mini-projects-a-unit-converter-and-a-word-frequency-counter).
Compare the two: the Rust version has no hand-written hash table, no
`strdup`, no `free`, and no fixed-size buffer.

```rust
// src/bin/wordfreq.rs
// usage: cargo run --bin wordfreq -- [FILE] [-n TOP]   (reads stdin with no FILE)
use std::collections::HashMap;
use std::io::{self, BufRead, BufReader, Read};

fn count_words(reader: impl BufRead) -> io::Result<HashMap<String, usize>> {
    let mut counts = HashMap::new();
    for line in reader.lines() {
        let line = line?;                       // propagate I/O and invalid-UTF-8 errors
        for word in line
            .split(|c: char| !c.is_alphanumeric() && c != '\'')
            .filter(|w| !w.is_empty())
        {
            *counts.entry(word.to_lowercase()).or_insert(0) += 1;
        }
    }
    Ok(counts)
}

fn main() -> io::Result<()> {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let top = args.iter().position(|a| a == "-n")
        .and_then(|i| args.get(i + 1))
        .and_then(|n| n.parse().ok())
        .unwrap_or(10);
    let path = args.iter().find(|a| !a.starts_with('-') && a.parse::<usize>().is_err());

    let input: Box<dyn Read> = match path {
        Some(p) => Box::new(std::fs::File::open(p)?),
        None => Box::new(io::stdin()),
    };
    let counts = count_words(BufReader::new(input))?;

    let mut ranked: Vec<(&String, &usize)> = counts.iter().collect();
    ranked.sort_unstable_by(|a, b| b.1.cmp(a.1).then_with(|| a.0.cmp(b.0)));
    for (word, n) in ranked.into_iter().take(top) {
        println!("{n:>7} {word}");
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn counts_case_insensitively() {
        let c = count_words("The cat. the CAT's hat!".as_bytes()).unwrap();
        assert_eq!(c["the"], 2);
        assert_eq!(c["cat's"], 1);
        assert_eq!(c.get("dog"), None);
    }
}
```

```bash
cargo run --release --bin wordfreq -- -n 5 /usr/share/dict/words
cargo test --bin wordfreq
```

**Stretch goals:**

1. Avoid allocating a `String` per word. Count `&str` slices of a buffer
   that holds the whole file. You will meet lifetimes (Ch 19), because the
   map now borrows from the buffer.
2. Process files in parallel with scoped threads (Ch 28) and merge the maps.
3. Compare `HashMap` with `BTreeMap` timings on a 100 MB file.

---

## 15. Terminal project: a `/etc/hosts`-style config parser with real errors

This project uses enums, `Result`, custom errors, `FromStr`, and tests, all
of which a real config loader needs.

```rust
// src/bin/hostsparse.rs
use std::collections::BTreeMap;
use std::fmt;
use std::net::IpAddr;
use std::str::FromStr;

#[derive(Debug, PartialEq)]
pub enum ParseError {
    BadAddress { line: usize, value: String },
    NoHostnames { line: usize },
    BadHostname { line: usize, value: String },
}

impl fmt::Display for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::BadAddress { line, value } => write!(f, "line {line}: {value:?} is not an IP address"),
            Self::NoHostnames { line } => write!(f, "line {line}: address with no hostnames"),
            Self::BadHostname { line, value } => write!(f, "line {line}: invalid hostname {value:?}"),
        }
    }
}
impl std::error::Error for ParseError {}

#[derive(Debug, Default)]
pub struct Hosts {
    by_name: BTreeMap<String, Vec<IpAddr>>,
}

fn valid_hostname(h: &str) -> bool {
    // RFC 1123-ish: labels of [a-z0-9-], not starting/ending with '-', total <= 253
    h.len() <= 253
        && h.split('.').all(|label| {
            !label.is_empty()
                && label.len() <= 63
                && !label.starts_with('-')
                && !label.ends_with('-')
                && label.bytes().all(|b| b.is_ascii_alphanumeric() || b == b'-')
        })
}

impl FromStr for Hosts {
    type Err = ParseError;

    fn from_str(text: &str) -> Result<Self, Self::Err> {
        let mut hosts = Hosts::default();
        for (idx, raw) in text.lines().enumerate() {
            let line = idx + 1;
            let content = raw.split('#').next().unwrap_or("").trim(); // strip comments
            if content.is_empty() { continue; }

            let mut fields = content.split_whitespace();
            let addr_text = fields.next().expect("non-empty line has a first field");
            let addr: IpAddr = addr_text
                .parse()
                .map_err(|_| ParseError::BadAddress { line, value: addr_text.into() })?;

            let names: Vec<&str> = fields.collect();
            if names.is_empty() { return Err(ParseError::NoHostnames { line }); }
            for name in names {
                if !valid_hostname(name) {
                    return Err(ParseError::BadHostname { line, value: name.into() });
                }
                hosts.by_name.entry(name.to_ascii_lowercase()).or_default().push(addr);
            }
        }
        Ok(hosts)
    }
}

impl Hosts {
    pub fn lookup(&self, name: &str) -> &[IpAddr] {
        self.by_name.get(&name.to_ascii_lowercase()).map(Vec::as_slice).unwrap_or(&[])
    }
}

fn main() {
    let text = std::fs::read_to_string("/etc/hosts").unwrap_or_else(|_| "127.0.0.1 localhost".into());
    match text.parse::<Hosts>() {
        Ok(h) => println!("localhost -> {:?}", h.lookup("LOCALHOST")),
        Err(e) => eprintln!("error: {e}"),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_and_looks_up() {
        let h: Hosts = "127.0.0.1 localhost  # loopback\n::1 localhost ip6-localhost\n".parse().unwrap();
        assert_eq!(h.lookup("localhost").len(), 2);
        assert!(h.lookup("nope").is_empty());
    }

    #[test]
    fn reports_line_numbers() {
        let err = "127.0.0.1 ok\n999.1.1.1 bad\n".parse::<Hosts>().unwrap_err();
        assert_eq!(err, ParseError::BadAddress { line: 2, value: "999.1.1.1".into() });
        assert!(matches!("10.0.0.1 -bad-".parse::<Hosts>(), Err(ParseError::BadHostname { .. })));
    }
}
```

**Across the series:** this is the file your resolver reads before it sends
any DNS query (see
[Net Ch 18](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses)).
Ch 45 builds the DNS query that runs when the name is not in this file.

---

# Part II — Intermediate: traits, lifetimes, and the standard library

Part I gave you values and ownership. Part II covers **abstraction**: traits
(Rust's interfaces), generics, closures, and iterators. It also covers the
part of the borrow checker that Part I avoided, **lifetimes**. By the end you
can read most library code on crates.io.

---

## 16. Traits and generics: shared behavior without inheritance

### 16.1 Defining and implementing a trait

```rust
trait Shape {
    fn area(&self) -> f64;                       // required method
    fn name(&self) -> String { "shape".into() }  // default method, can be overridden
    fn describe(&self) -> String {               // defaults can call required methods
        format!("{} with area {:.2}", self.name(), self.area())
    }
}

struct Circle { r: f64 }
struct Rect { w: f64, h: f64 }

impl Shape for Circle {
    fn area(&self) -> f64 { std::f64::consts::PI * self.r * self.r }
    fn name(&self) -> String { "circle".into() }
}
impl Shape for Rect {
    fn area(&self) -> f64 { self.w * self.h }
}

fn main() {
    println!("{}", Circle { r: 1.0 }.describe()); // circle with area 3.14
    println!("{}", Rect { w: 2.0, h: 3.0 }.describe()); // shape with area 6.00
}
```

Unlike [Go interfaces](../Golang/real-life-golang-guide.md#9-interfaces-implicit-satisfaction-any-type-switches),
Rust traits are implemented **explicitly** with `impl Trait for Type`. You
can implement your trait for types you did not write, including `i32` and
`Vec<T>`. The **orphan rule** restricts this: either the trait or the type
must be defined in your crate. That rule stops two crates from providing
conflicting implementations.

### 16.2 Generic functions and trait bounds

```rust
use std::fmt::Display;

// "T can be any type that implements PartialOrd and Copy"
fn largest<T: PartialOrd + Copy>(items: &[T]) -> Option<T> {
    let mut it = items.iter().copied();
    let first = it.next()?;
    Some(it.fold(first, |best, x| if x > best { x } else { best }))
}

// `where` clauses read better when bounds get long.
fn report<K, V>(pairs: &[(K, V)]) -> String
where
    K: Display,
    V: Display + PartialOrd<i32>,
{
    pairs
        .iter()
        .map(|(k, v)| format!("{k}={v}{}", if *v > 100 { " (!)" } else { "" }))
        .collect::<Vec<_>>()
        .join(", ")
}

// `impl Trait` in argument position is shorthand for a generic parameter.
fn shout(msg: impl AsRef<str>) -> String { msg.as_ref().to_uppercase() }

fn main() {
    println!("{:?} {:?}", largest(&[3, 9, 2]), largest::<f64>(&[]));
    println!("{}", report(&[("latency_ms", 120), ("errors", 3)]));
    println!("{} {}", shout("hi"), shout(String::from("there")));
}
```

### 16.3 Monomorphization: what generics compile to

The compiler generates **a separate copy of a generic function for every
concrete type it is used with**. `largest::<i32>` and `largest::<f64>`
become two functions in the binary, each optimized for its type, with no
indirection at runtime.

| | Rust generics | Go generics | C |
|---|---|---|---|
| Implementation | monomorphization (one copy per type) | GC-shape stenciling + dictionaries | macros or `void *` |
| Runtime cost | none; can be inlined | sometimes an indirect call | macros: none; `void *`: casts and no type safety |
| Compile time / binary size | larger | smaller | — |

This is what "zero-cost abstraction" means: the generic code runs as fast as
code written by hand for each type. The cost moves to compile time and binary
size. When that matters, you can switch to `dyn Trait` (Ch 17), which uses
one copy of the code and dispatches through a vtable.

### 16.4 Associated types and generic traits

```rust
// A trait with an ASSOCIATED TYPE: each implementor picks exactly one `Output`.
trait Parser {
    type Output;
    fn parse(&self, input: &str) -> Option<Self::Output>;
}

struct PortParser;
impl Parser for PortParser {
    type Output = u16;
    fn parse(&self, input: &str) -> Option<u16> { input.parse().ok().filter(|p| *p != 0) }
}

struct CsvParser;
impl Parser for CsvParser {
    type Output = Vec<String>;
    fn parse(&self, input: &str) -> Option<Vec<String>> {
        Some(input.split(',').map(|s| s.trim().to_string()).collect())
    }
}

fn parse_all<P: Parser>(p: &P, inputs: &[&str]) -> Vec<Option<P::Output>> {
    inputs.iter().map(|i| p.parse(i)).collect()
}

fn main() {
    println!("{:?}", parse_all(&PortParser, &["80", "0", "x"]));  // [Some(80), None, None]
    println!("{:?}", parse_all(&CsvParser, &["a, b,c"]));
}
```

Use an **associated type** when each implementor has exactly one natural
choice, as with `Iterator::Item`. Use a **generic parameter** (`trait From<T>`)
when a type can implement the trait many times for different `T`.

### 16.5 Real-world example

`std::io::Write` is a trait. Files, sockets, `Vec<u8>`, `Stdout`,
`BufWriter<W>`, TLS streams, and gzip encoders all implement it. Code
written as `fn write_report<W: Write>(out: &mut W)` works with every one of
them and gets monomorphized for each. This is Go's
[`io.Writer` pattern](../Golang/real-life-golang-guide.md#22-i-o-io-reader-writer-bufio-and-os),
with static dispatch instead of an interface value.

---

## 17. Trait objects: `dyn Trait` and dynamic dispatch

Generics need the concrete type at compile time. Sometimes you only know it
at runtime, for example with a list of plugins or middleware loaded from
config. For that, use a **trait object**.

```rust
trait Middleware {
    fn name(&self) -> &str;
    fn handle(&self, path: &str) -> Result<(), String>;
}

struct Auth { token: String }
struct RateLimit { max_path_len: usize }

impl Middleware for Auth {
    fn name(&self) -> &str { "auth" }
    fn handle(&self, path: &str) -> Result<(), String> {
        if path.starts_with("/admin") && self.token != "s3cret" { Err("forbidden".into()) } else { Ok(()) }
    }
}
impl Middleware for RateLimit {
    fn name(&self) -> &str { "limit" }
    fn handle(&self, path: &str) -> Result<(), String> {
        if path.len() > self.max_path_len { Err("uri too long".into()) } else { Ok(()) }
    }
}

fn main() {
    // A Vec of DIFFERENT concrete types, all behind the same interface.
    let chain: Vec<Box<dyn Middleware>> = vec![
        Box::new(RateLimit { max_path_len: 32 }),
        Box::new(Auth { token: "wrong".into() }),
    ];
    for path in ["/home", "/admin/users", "/a-very-long-path-that-should-be-rejected"] {
        let verdict = chain.iter().try_for_each(|m| m.handle(path).map_err(|e| format!("{}: {e}", m.name())));
        println!("{path:<45} {verdict:?}");
    }
    println!("{}", std::mem::size_of::<&dyn Middleware>()); // 16: data pointer + vtable pointer
}
```

### 17.1 What a trait object is in memory

`&dyn Middleware` and `Box<dyn Middleware>` are **fat pointers**: one
pointer to the data and one pointer to a **vtable**, a static table of
function pointers for that concrete type.

```
 Box<dyn Middleware>                         vtable for Auth (static, one per type)
 ┌──────────────┐                            ┌─────────────────────┐
 │ data ptr ────┼──▶ Auth { token }          │ drop_in_place<Auth> │
 │ vtable ptr ──┼──────────────────────────▶ │ size, align         │
 └──────────────┘                            │ Auth::name          │
                                             │ Auth::handle        │
                                             └─────────────────────┘
```

This is the same structure as a
[Go interface value](../Golang/real-life-golang-guide.md#51-what-s-really-behind-an-interface-value)
(type word + data word), and the same thing the
[C guide's function-pointer structs](../c-lang/real-life-c-guide.md#11-pointers-to-pointers-function-pointers-pointers-to-arrays)
build by hand. One difference from Go: the vtable pointer sits in the
**reference**, not in the object, so a plain `Auth` value carries no type
header at all.

### 17.2 Static vs dynamic dispatch: how to choose

| | Generics `T: Trait` | Trait objects `dyn Trait` |
|---|---|---|
| Dispatch | static, inlinable | indirect call through the vtable |
| Mixed types in one collection | no | yes |
| Binary size | one copy per type | one copy |
| Typical use | hot paths, libraries | plugins, heterogeneous lists, reducing compile time |

### 17.3 Dyn compatibility (object safety)

Not every trait can become a `dyn Trait`. Methods that return `Self`, or
that have their own generic type parameters, cannot be put in a vtable,
because the caller would need to know the concrete type. For example,
`Clone` (`fn clone(&self) -> Self`) is not dyn-compatible. If one method
breaks this, add `where Self: Sized` to it to exclude it from the vtable and
keep the rest of the trait usable as `dyn`.

### 17.4 Downcasting

`dyn Any` can be downcast back to the concrete type: `if let Some(a) = obj.downcast_ref::<Auth>()`.
This is Rust's version of Go's type switch, and you should need it about as
rarely. Error types are the main exception: `anyhow::Error::downcast_ref`
is how you check for one specific error inside an error chain.

---

## 18. The standard traits you will implement every week

| Trait | What it gives you | Usually |
|---|---|---|
| `Debug` | `{:?}` formatting | `#[derive(Debug)]` on almost everything |
| `Display` | `{}` formatting, `.to_string()` | implement by hand for user-facing types |
| `Clone`, `Copy` | `.clone()`, implicit copy | derive |
| `PartialEq`, `Eq` | `==` | derive; `Eq` only if equality is reflexive (not for floats) |
| `Hash` | usable as a `HashMap` key | derive together with `Eq` |
| `PartialOrd`, `Ord` | `<`, sorting, `BTreeMap` keys | derive (orders by fields, top to bottom) |
| `Default` | `T::default()`, `..Default::default()` | derive or implement |
| `From<T>` / `Into<T>` | infallible conversion | implement `From`; you get `Into` for free |
| `TryFrom<T>` | fallible conversion | validation at boundaries (Ch 53) |
| `FromStr` | `"..".parse::<T>()` | config and CLI parsing |
| `AsRef<T>` | cheap reference conversion | flexible function parameters |
| `Deref` | `*x`, method forwarding | smart pointers only; do not use it to fake inheritance |
| `Drop` | destructor | resource handles |
| `Iterator`, `IntoIterator` | `for` loops, adapters | Ch 21 |
| `Send`, `Sync` | thread safety | implemented automatically (Ch 28) |

```rust
use std::fmt;
use std::ops::Add;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, PartialOrd, Ord, Default)]
struct Bytes(u64);

impl fmt::Display for Bytes {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        const UNITS: [&str; 5] = ["B", "KiB", "MiB", "GiB", "TiB"];
        let mut v = self.0 as f64;
        let mut i = 0;
        while v >= 1024.0 && i < UNITS.len() - 1 { v /= 1024.0; i += 1; }
        if i == 0 { write!(f, "{} B", self.0) } else { write!(f, "{v:.1} {}", UNITS[i]) }
    }
}

impl Add for Bytes {                    // operator overloading is just a trait
    type Output = Bytes;
    fn add(self, rhs: Bytes) -> Bytes { Bytes(self.0 + rhs.0) }
}

impl From<u64> for Bytes { fn from(n: u64) -> Self { Bytes(n) } }

fn main() {
    let total: Bytes = Bytes(1536) + 1_048_576.into(); // `.into()` uses From<u64>
    let mut sizes = vec![Bytes(10), Bytes(3), Bytes::default()];
    sizes.sort();                                      // Ord derive
    println!("{total} {:?} {}", sizes, Bytes(512));    // 1.0 MiB [Bytes(0), Bytes(3), Bytes(10)] 512 B
}
```

**`Eq` and `Hash` must agree.** If `a == b`, then `hash(a) == hash(b)`. If
you implement `PartialEq` by hand (for example, case-insensitive), implement
`Hash` by hand in the same way. Otherwise `HashMap` gives wrong answers
without any error.

---

## 19. Lifetimes: naming how long a borrow is valid

### 19.1 What a lifetime is

A lifetime is **not** how long a value lives. It is a name for the region of
code in which a **borrow** must stay valid. Most of the time the compiler
infers them. You write them only when a function or struct relates the
lifetimes of several references and the compiler cannot work out which input
the output borrows from.

```rust
// DOES NOT COMPILE
fn longer(a: &str, b: &str) -> &str {
    if a.len() >= b.len() { a } else { b }
}
fn main() { println!("{}", longer("x", "yy")); }
```

```text
error[E0106]: missing lifetime specifier
  |     fn longer(a: &str, b: &str) -> &str {
  |                  ----     ----     ^ expected named lifetime parameter
  = help: this function's return type contains a borrowed value, but the signature does
          not say whether it is borrowed from `a` or `b`
```

The fix states the relationship: "the result borrows from both inputs, so it
is valid only while both are":

```rust
fn longer<'a>(a: &'a str, b: &'a str) -> &'a str {
    if a.len() >= b.len() { a } else { b }
}

fn main() {
    let outer = String::from("long string");
    let result;
    {
        let inner = String::from("xyz");
        result = longer(&outer, &inner);
        println!("{result}"); // fine: `inner` is still alive here
    }
    // println!("{result}"); // E0597: `inner` does not live long enough
}
```

Lifetimes are **checked, never chosen**. Annotations do not make anything
live longer. They describe constraints, and the compiler checks that every
call site satisfies them.

### 19.2 The three elision rules

You rarely write lifetimes because the compiler fills them in using three
rules:

1. Each reference parameter gets its own lifetime.
2. If there is exactly one input lifetime, the output gets that lifetime.
3. If one of the parameters is `&self` or `&mut self`, the output gets
   `self`'s lifetime.

So `fn first_word(s: &str) -> &str` and `fn name(&self) -> &str` need no
annotations. `fn longer(a: &str, b: &str) -> &str` has two inputs and no
`self`, so none of the rules decide the output and you must write it.

### 19.3 Structs that borrow: zero-copy parsing

```rust
#[derive(Debug)]
struct RequestLine<'a> {   // this struct cannot outlive the buffer it borrows from
    method: &'a str,
    path: &'a str,
    version: &'a str,
}

fn parse_request_line(line: &str) -> Option<RequestLine<'_>> { // '_ = "the elided lifetime"
    let mut parts = line.split(' ');
    let (method, path, version) = (parts.next()?, parts.next()?, parts.next()?);
    if parts.next().is_some() || !version.starts_with("HTTP/") { return None; }
    Some(RequestLine { method, path, version })
}

fn main() {
    let raw = String::from("GET /index.html HTTP/1.1");
    let req = parse_request_line(&raw).expect("valid request line");
    println!("{req:?}");   // no allocation: every field points into `raw`
    // drop(raw);          // E0505: cannot move out of `raw` because it is borrowed
    println!("{}", req.path);
}
```

This is how fast parsers such as `httparse` (used by Hyper, Ch 50), `serde`
borrowed deserialization, and `nom` work. They allocate nothing, and the
compiler makes sure the parsed view cannot outlive the buffer it points
into. In C, the same design is a common source of use-after-free bugs.

### 19.4 `'static`: two meanings you must keep apart

- **`&'static T`**: a reference valid for the whole program, such as string
  literals, `static` items, and leaked boxes.
- **`T: 'static`** (a bound): "T contains no non-static borrows". **Owned
  types such as `String` and `Vec<u8>` satisfy it.** It does *not* mean
  "lives forever". `thread::spawn` and `tokio::spawn` require `T: 'static`
  because the new thread or task may outlive the current stack frame, so it
  must own everything it uses (Ch 28).

```rust
fn needs_static<T: 'static + std::fmt::Debug>(t: T) { println!("{t:?}"); }

fn main() {
    needs_static(String::from("owned")); // OK: String owns its data
    needs_static(42);                    // OK
    let local = String::from("x");
    // needs_static(&local);             // E0597: `local` does not live long enough
    needs_static(local.clone());         // the usual fix: give it an owned value
}
```

### 19.5 When you are fighting lifetimes, own the data instead

If a struct with lifetime parameters spreads `'a` through half your
codebase, store owned data (`String`, `Vec<u8>`, `Arc<str>`) instead. Use
borrowing structs for short-lived views such as parsers and iterators. Use
owned structs for anything that gets stored, sent to another thread, or
returned from an API. The extra allocation is usually cheaper than the
design cost.

---

## 20. Closures: `Fn`, `FnMut`, `FnOnce`

A closure is an anonymous struct holding its captured variables, plus an
implementation of one of three traits, chosen by **how the body uses its
captures**:

| Trait | Body does | Callable | Example |
|---|---|---|---|
| `Fn` | reads captures | many times, also concurrently | `|x| x + offset` |
| `FnMut` | mutates captures | many times, one call at a time | `|| { count += 1 }` |
| `FnOnce` | moves a capture out | once | `|| drop(buffer)` |

```rust
fn apply_twice(f: impl Fn(i32) -> i32, x: i32) -> i32 { f(f(x)) }
fn run_n(mut f: impl FnMut(), n: usize) { for _ in 0..n { f(); } }
fn run_once(f: impl FnOnce() -> String) -> String { f() }

fn make_adder(n: i32) -> impl Fn(i32) -> i32 {
    move |x| x + n // `move`: the closure takes ownership of `n` (here a copy)
}

fn main() {
    let offset = 10;
    println!("{}", apply_twice(|x| x + offset, 1)); // 21

    let mut count = 0;
    run_n(|| count += 1, 3);
    println!("count={count}"); // 3

    let greeting = String::from("hi");
    println!("{}", run_once(move || greeting + " there")); // consumes `greeting`

    let add5 = make_adder(5);
    let handlers: Vec<Box<dyn Fn(i32) -> i32>> = vec![Box::new(add5), Box::new(|x| x * 2)];
    println!("{:?}", handlers.iter().map(|h| h(10)).collect::<Vec<_>>()); // [15, 20]
}
```

**`move` closures** take ownership of what they capture. You need them
whenever the closure outlives the current stack frame: in threads, tasks,
returned closures, and stored callbacks. Without `move`, a closure captures
by reference. That is usually what you want locally, and it is rejected for
spawned work.

Compare [Go §5.4](../Golang/real-life-golang-guide.md#5-functions-multiple-returns-variadics-closures):
Go closures always capture by reference and the GC keeps captured variables
alive, which is how Go got its loop-variable capture bug before 1.22. Rust
forces you to choose between borrowing and moving, and checks the choice.

---

## 21. Iterators: lazy, composable, zero-cost

### 21.1 The trait

```rust
// fragment: the real definition, minus ~75 provided methods
pub trait Iterator {
    type Item;
    fn next(&mut self) -> Option<Self::Item>;
}
```

You implement `next`. Every adapter (`map`, `filter`, `zip`, `take`,
`chain`, `flat_map`, `fold`, `sum`, ...) comes for free.

### 21.2 `iter`, `iter_mut`, `into_iter`

| Call | Yields | The collection afterwards |
|---|---|---|
| `v.iter()` / `for x in &v` | `&T` | unchanged |
| `v.iter_mut()` / `for x in &mut v` | `&mut T` | modified in place |
| `v.into_iter()` / `for x in v` | `T` (owned) | consumed |

### 21.3 Adapters are lazy

```rust
#[derive(Debug)]
struct LogLine { status: u16, ms: u32, path: &'static str }

fn main() {
    let logs = [
        LogLine { status: 200, ms: 12, path: "/" },
        LogLine { status: 500, ms: 340, path: "/api/pay" },
        LogLine { status: 200, ms: 95, path: "/api/cart" },
        LogLine { status: 503, ms: 1200, path: "/api/pay" },
    ];

    // Nothing runs until a CONSUMER (sum, collect, for, count...) pulls values through.
    let error_rate = logs.iter().filter(|l| l.status >= 500).count() as f64 / logs.len() as f64;
    let slow_paths: Vec<&str> = logs.iter().filter(|l| l.ms > 100).map(|l| l.path).collect();
    let p_max = logs.iter().map(|l| l.ms).max().unwrap_or(0);
    let total_ok_ms: u32 = logs.iter().filter(|l| l.status < 400).map(|l| l.ms).sum();

    println!("error rate {:.0}%, slow {slow_paths:?}, max {p_max}ms, ok total {total_ok_ms}ms", error_rate * 100.0);

    // collect() is driven by the TARGET type: Result<Vec<_>, _> stops at the first error.
    let parsed: Result<Vec<u16>, _> = ["80", "443", "x"].iter().map(|s| s.parse::<u16>()).collect();
    println!("{parsed:?}"); // Err(ParseIntError { kind: InvalidDigit })
}
```

**Across the series:** this is the error-rate calculation behind the SLIs
in [HTTPS Ch 23](../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services).

### 21.4 Writing your own iterator

```rust
/// Splits a byte buffer into length-prefixed frames: [len:u8][payload...]...
struct Frames<'a> { buf: &'a [u8] }

impl<'a> Iterator for Frames<'a> {
    type Item = Result<&'a [u8], &'static str>;
    fn next(&mut self) -> Option<Self::Item> {
        let (&len, rest) = self.buf.split_first()?;     // None at end of input
        let len = len as usize;
        if rest.len() < len {
            self.buf = &[];                              // stop after reporting the error
            return Some(Err("truncated frame"));
        }
        let (frame, rest) = rest.split_at(len);
        self.buf = rest;
        Some(Ok(frame))
    }
}

fn main() {
    let wire = [3, b'a', b'b', b'c', 0, 2, b'h', b'i', 9, b'x'];
    for f in (Frames { buf: &wire }) {
        println!("{:?}", f.map(|b| String::from_utf8_lossy(b).into_owned()));
    }
}
```

The same framing problem appears in [Go §34](../Golang/real-life-golang-guide.md#34-networking-deep-dive-net-conn-tcp-udp-framing-your-own-protocol)
and in Ch 32 over a real socket.

### 21.5 Zero cost

`logs.iter().filter(..).map(..).sum()` compiles to the same machine code as
a hand-written `for` loop with an `if`. Adapters are generic structs
(`Filter<Map<Iter<..>>>`), and after inlining LLVM sees only the loop.
Iterators also avoid bounds checks, because the iterator already knows where
the slice ends. When performance matters, prefer iterators to manual
indexing.

---

## 22. Smart pointers: `Box`, `Rc`, `Arc`, `Cell`, `RefCell`, `Weak`

| Type | Ownership | Mutation | Thread-safe | Use for |
|---|---|---|---|---|
| `Box<T>` | single owner, on the heap | via `&mut` | if `T` is | recursive types, trait objects, large values |
| `Rc<T>` | shared, reference-counted | no (immutable) | **no** | graphs and trees on one thread |
| `Arc<T>` | shared, atomically counted | no | yes | sharing across threads and tasks |
| `Cell<T>` | — | `get`/`set` for `Copy` types | no | counters and flags behind `&` |
| `RefCell<T>` | — | borrow checking **at runtime** | no | interior mutability on one thread |
| `Mutex<T>`, `RwLock<T>` | — | locking at runtime | yes | interior mutability across threads (Ch 29) |
| `Weak<T>` | non-owning | — | with `Arc` | parent pointers, caches, breaking cycles |

### 22.1 `Box` for recursive types

```rust
#[derive(Debug)]
enum Expr {
    Num(f64),
    Add(Box<Expr>, Box<Expr>),  // without Box, Expr would have infinite size
    Mul(Box<Expr>, Box<Expr>),
}

fn eval(e: &Expr) -> f64 {
    match e {
        Expr::Num(n) => *n,
        Expr::Add(a, b) => eval(a) + eval(b),
        Expr::Mul(a, b) => eval(a) * eval(b),
    }
}

fn main() {
    use Expr::*;
    let e = Add(Box::new(Num(2.0)), Box::new(Mul(Box::new(Num(3.0)), Box::new(Num(4.0)))));
    println!("2 + 3 * 4 = {}", eval(&e)); // 14
}
```

### 22.2 `Rc`, `RefCell`, and `Weak`: a tree with parent pointers

```rust
use std::cell::RefCell;
use std::rc::{Rc, Weak};

#[derive(Debug)]
struct Node {
    name: String,
    parent: RefCell<Weak<Node>>,        // Weak: a child must not keep its parent alive (no cycle)
    children: RefCell<Vec<Rc<Node>>>,   // Rc: the parent owns its children
}

fn node(name: &str) -> Rc<Node> {
    Rc::new(Node { name: name.into(), parent: RefCell::new(Weak::new()), children: RefCell::new(vec![]) })
}

fn main() {
    let root = node("/");
    let etc = node("etc");
    *etc.parent.borrow_mut() = Rc::downgrade(&root);
    root.children.borrow_mut().push(Rc::clone(&etc)); // Rc::clone only bumps a counter

    let parent_name = etc.parent.borrow().upgrade().map(|p| p.name.clone());
    println!("etc's parent = {parent_name:?}");
    println!("strong(root)={} weak(root)={} strong(etc)={}",
        Rc::strong_count(&root), Rc::weak_count(&root), Rc::strong_count(&etc)); // 1 1 2

    // RefCell enforces shared-XOR-mutable at RUNTIME:
    let _reading = root.children.borrow();
    let attempt = root.children.try_borrow_mut();       // borrow_mut() here would PANIC
    println!("second mutable borrow allowed? {}", attempt.is_ok()); // false
}
```

If both directions used `Rc`, the parent and child would keep each other's
count above zero, and neither would ever be freed. That is a reference
cycle, and it is a memory leak that safe Rust allows. Rust's rules prevent
dangling pointers, not leaks. Use `Weak` for back-pointers, or use an arena
with indices (Ch 27).

### 22.3 Interior mutability and why it exists

The borrow rules apply to `&T` versus `&mut T`. Some designs need mutation
through a shared reference: a cache inside an otherwise read-only object, a
reference count, or a mock that records calls. `Cell`, `RefCell`, `Mutex`,
and atomics allow it by **moving the shared-XOR-mutable check to runtime**
(or making it unnecessary, in `Cell`'s case). All of them are built on
`UnsafeCell`, the only legal way in Rust to mutate through a `&`.

`Rc<RefCell<T>>` everywhere is a sign that you are writing a garbage-collected
object graph in Rust. It works, but every `borrow_mut()` can panic at
runtime, and you lose the compile-time guarantees that made Rust worth
choosing. See Ch 68.

---

## 23. Testing: unit, integration, doc tests, and property tests

### 23.1 Unit tests live next to the code

```rust
pub fn parse_duration(s: &str) -> Result<u64, String> {
    let (num, unit) = s.split_at(s.find(|c: char| !c.is_ascii_digit()).unwrap_or(s.len()));
    let n: u64 = num.parse().map_err(|_| format!("no number in {s:?}"))?;
    match unit {
        "ms" => Ok(n),
        "s" | "" => Ok(n * 1000),
        "m" => Ok(n * 60_000),
        _ => Err(format!("unknown unit {unit:?}")),
    }
}

#[cfg(test)]                 // compiled only for `cargo test`
mod tests {
    use super::*;

    #[test]
    fn units() {
        for (input, want) in [("250ms", 250), ("3s", 3000), ("2m", 120_000), ("7", 7000)] {
            assert_eq!(parse_duration(input), Ok(want), "input {input:?}"); // table-driven, like Go
        }
    }

    #[test]
    fn rejects_garbage() -> Result<(), String> {   // tests may return Result and use `?`
        assert!(parse_duration("fast").is_err());
        assert!(parse_duration("5h").unwrap_err().contains("unknown unit"));
        Ok(())
    }

    #[test]
    #[should_panic(expected = "overflow")]
    fn overflow_panics_in_debug() {
        let big = parse_duration("18446744073709551615s").unwrap(); // u64::MAX * 1000
        let _ = big;
    }
}
```

The `should_panic` test passes under `cargo test`, which builds in debug mode
with overflow checks on. It would fail under `cargo test --release`, which is
a reminder of §3.3: use `checked_mul` in real code.

### 23.2 Integration tests and doc tests

```
mycrate/
├── src/lib.rs
└── tests/               # each file is a separate crate that can only use your PUBLIC API
    └── http_roundtrip.rs
```

````rust
/// Adds two numbers.
///
/// ```
/// assert_eq!(mycrate::add(2, 2), 4);   // this example is COMPILED AND RUN by `cargo test`
/// ```
pub fn add(a: i32, b: i32) -> i32 { a + b }
````

Doc tests keep the examples in your documentation from going stale.

| Command | Runs |
|---|---|
| `cargo test` | everything |
| `cargo test parse_` | tests whose name contains `parse_` |
| `cargo test -- --nocapture` | show `println!` output from passing tests |
| `cargo test --doc` | doc tests only |
| `cargo nextest run` | a faster runner, one process per test (`cargo install cargo-nextest`) |

### 23.3 Property-based testing

A table-driven test checks the examples you thought of. A property test
generates thousands of inputs and checks a rule that must always hold:

```rust
#[cfg(test)]
mod props {
    use proptest::prelude::*;

    fn encode(s: &str) -> String {
        s.bytes().map(|b| format!("{b:02x}")).collect()
    }
    fn decode(h: &str) -> Option<String> {
        let bytes: Option<Vec<u8>> = (0..h.len()).step_by(2)
            .map(|i| h.get(i..i + 2).and_then(|p| u8::from_str_radix(p, 16).ok()))
            .collect();
        String::from_utf8(bytes?).ok()
    }

    proptest! {
        #[test]
        fn roundtrip(s in ".*") {                    // any Unicode string
            prop_assert_eq!(decode(&encode(&s)), Some(s));
        }
        #[test]
        fn decode_never_panics(h in "\\PC*") {       // arbitrary garbage input
            let _ = decode(&h);
        }
    }
}
```

When a property fails, `proptest` **shrinks** the input to the smallest case
that still fails. Ch 54 adds fuzzing, which applies the same idea with
coverage guidance.

---

## 24. I/O, files, processes, and application errors

### 24.1 `Read`, `Write`, `BufRead`, and buffering

```rust
use std::fs::File;
use std::io::{self, BufRead, BufReader, BufWriter, Write};

fn main() -> io::Result<()> {
    let path = std::env::temp_dir().join("rust-guide-io.txt");

    // Writes: wrap the File in a BufWriter, or every write! is one write(2) syscall.
    {
        let mut out = BufWriter::new(File::create(&path)?);
        for i in 0..1000 {
            writeln!(out, "line {i}")?;
        }
        out.flush()?; // BufWriter flushes on drop too, but drop cannot report errors. Flush explicitly.
    }

    // Reads: BufReader + lines() streams the file; memory use stays flat for any file size.
    let reader = BufReader::new(File::open(&path)?);
    let n = reader.lines().filter(|l| l.as_ref().is_ok_and(|l| l.ends_with('7'))).count();
    println!("{n} lines end in 7");

    // Whole-file helpers for small files:
    let text = std::fs::read_to_string(&path)?;
    println!("{} bytes", text.len());

    // stdout is LINE-buffered and behind a lock; lock it once for bulk output.
    let stdout = io::stdout();
    let mut lock = stdout.lock();
    writeln!(lock, "done")?;
    std::fs::remove_file(&path)
}
```

**War story.** A log-processing tool written in Rust was slower than the
Python script it replaced. It used `println!` in a loop over 50 million
lines. Each `println!` takes the stdout lock, and when stdout is a pipe, it
issues a `write` syscall for every line. Wrapping `stdout().lock()` in a
`BufWriter` made it about 30 times faster. Rust does what you write, and
buffering is something you have to ask for. The
[C guide §12.1](../c-lang/real-life-c-guide.md#12-file-i-o-buffering-binary-vs-text-mode-errno)
describes the opposite surprise: C's stdio buffers by default, so output
sometimes appears late.

### 24.2 Paths and `OsStr`

File names on Unix are arbitrary bytes, and on Windows they are UTF-16.
Neither is guaranteed to be valid UTF-8, so `Path`/`PathBuf` and
`OsStr`/`OsString` exist. Use `path.display()` to print a path, and
`path.to_str()` (which returns `Option<&str>`) when you need it as UTF-8.

### 24.3 Running processes

```rust
use std::process::{Command, Stdio};

fn main() -> std::io::Result<()> {
    // Arguments are passed as an argv ARRAY, not through a shell: no shell injection.
    let out = Command::new("ls").arg("-1").arg("/").output()?;
    println!("exit={} first entries: {:?}",
        out.status, String::from_utf8_lossy(&out.stdout).lines().take(3).collect::<Vec<_>>());

    // Pipe data into a child process's stdin:
    let mut child = Command::new("wc").arg("-l").stdin(Stdio::piped()).stdout(Stdio::piped()).spawn()?;
    {
        use std::io::Write;
        let mut stdin = child.stdin.take().expect("stdin was piped");
        stdin.write_all(b"a\nb\nc\n")?;
    } // dropping stdin closes the pipe, so wc sees EOF
    let result = child.wait_with_output()?;
    println!("wc says {}", String::from_utf8_lossy(&result.stdout).trim());
    Ok(())
}
```

`Command` never runs a shell unless you ask for one with
`Command::new("sh").arg("-c")`. Do that only with strings you fully control.
This is the safe default against OS command injection (see
[Security from Zero Ch 42](../security/real-life-guide.md#chapter-42-injection-when-data-becomes-code)).
Ch 40 covers exit codes and signals in depth.

---

## 25. 🔎 Checkpoint: traits and lifetimes

**1. Why does this fail, and what are two fixes?**

```rust
// DOES NOT COMPILE
fn spawn_logger(prefix: &str) -> std::thread::JoinHandle<()> {
    std::thread::spawn(|| println!("{prefix}: started"))
}
fn main() { spawn_logger("app").join().unwrap(); }
```

<details><summary>Answer</summary>

`thread::spawn` requires its closure to be `'static` (§19.4), because the
thread may outlive the caller's stack frame, and this closure borrows
`prefix`. The compiler reports E0373 (closure may outlive the current
function) and E0521 (borrowed data escapes the function). Fix 1: own the data with
`let prefix = prefix.to_string();` and then `move || ...`. Fix 2: if the
caller waits for the thread before returning, use `std::thread::scope`
(Ch 28), which allows borrowing.
</details>

**2. Generic or `dyn`?** You are writing a function that compresses data
with one of three algorithms chosen by a config value at startup, and it
runs once per request. Which do you use?

<details><summary>Answer</summary>

`Box<dyn Compressor>` (or an `enum` of the three). The algorithm is chosen
at runtime, so you need a single type that can hold any of them. A virtual
call once per request costs nothing measurable. Use generics when the type
is known at compile time and the call is in a tight inner loop.
</details>

**3. Predict:** does `fn first<'a>(x: &'a str, y: &str) -> &'a str { x }` compile? What about returning `y`?

<details><summary>Answer</summary>

Returning `x` compiles: the signature says the output borrows from `x`.
Returning `y` fails, because `y` has a different, unrelated lifetime. The
compiler suggests adding `'a` to `y` (E0621, explicit lifetime required).
This is the point of annotations: they are a contract that the function body
is checked against.
</details>

**4. Write it:** implement `Iterator` for `struct Countdown(u32)` that
yields `3, 2, 1` for `Countdown(3)`.

<details><summary>One correct answer</summary>

```rust
struct Countdown(u32);
impl Iterator for Countdown {
    type Item = u32;
    fn next(&mut self) -> Option<u32> {
        if self.0 == 0 { return None; }
        self.0 -= 1;
        Some(self.0 + 1)
    }
}
fn main() { println!("{:?}", Countdown(3).collect::<Vec<_>>()); } // [3, 2, 1]
```
</details>

---

## 26. Terminal project: `rgrep`, a tested `grep` clone

A small but real CLI: library and binary split, typed config, errors with
context, streaming I/O, and tests.

```rust
// src/bin/rgrep.rs   usage: rgrep [-i] [-n] [-v] [-c] PATTERN [FILE...]
use anyhow::{bail, Context, Result};
use std::fs::File;
use std::io::{self, BufRead, BufReader, BufWriter, Write};

#[derive(Debug, Default, PartialEq)]
struct Config {
    pattern: String,
    files: Vec<String>,
    ignore_case: bool,
    line_numbers: bool,
    invert: bool,
    count_only: bool,
}

fn parse_args<I: IntoIterator<Item = String>>(args: I) -> Result<Config> {
    let mut cfg = Config::default();
    let mut positional = Vec::new();
    for arg in args {
        match arg.as_str() {
            "-i" => cfg.ignore_case = true,
            "-n" => cfg.line_numbers = true,
            "-v" => cfg.invert = true,
            "-c" => cfg.count_only = true,
            s if s.starts_with('-') && s.len() > 1 => bail!("unknown flag {s}"),
            _ => positional.push(arg),
        }
    }
    let mut it = positional.into_iter();
    cfg.pattern = it.next().context("missing PATTERN")?;
    cfg.files = it.collect();
    Ok(cfg)
}

/// Core logic, generic over any reader and writer, so tests need no files.
fn search<R: BufRead, W: Write>(cfg: &Config, label: Option<&str>, input: R, out: &mut W) -> Result<usize> {
    let needle = if cfg.ignore_case { cfg.pattern.to_lowercase() } else { cfg.pattern.clone() };
    let mut matches = 0;
    for (i, line) in input.lines().enumerate() {
        let line = line?;
        let hay = if cfg.ignore_case { line.to_lowercase() } else { line.clone() };
        if hay.contains(&needle) != cfg.invert {
            matches += 1;
            if !cfg.count_only {
                if let Some(l) = label { write!(out, "{l}:")?; }
                if cfg.line_numbers { write!(out, "{}:", i + 1)?; }
                writeln!(out, "{line}")?;
            }
        }
    }
    Ok(matches)
}

fn main() -> Result<()> {
    let cfg = parse_args(std::env::args().skip(1))?;
    let stdout = io::stdout();
    let mut out = BufWriter::new(stdout.lock());
    let mut total = 0;
    if cfg.files.is_empty() {
        total += search(&cfg, None, io::stdin().lock(), &mut out)?;
    } else {
        let label_files = cfg.files.len() > 1;
        for path in &cfg.files {
            let f = File::open(path).with_context(|| format!("opening {path}"))?;
            total += search(&cfg, label_files.then_some(path.as_str()), BufReader::new(f), &mut out)?;
        }
    }
    if cfg.count_only { writeln!(out, "{total}")?; }
    out.flush()?;
    std::process::exit(if total > 0 { 0 } else { 1 }); // grep's convention: 1 = no match
}

#[cfg(test)]
mod tests {
    use super::*;

    fn run(args: &[&str], input: &str) -> (String, usize) {
        let cfg = parse_args(args.iter().map(|s| s.to_string())).unwrap();
        let mut out = Vec::new();
        let n = search(&cfg, None, input.as_bytes(), &mut out).unwrap();
        (String::from_utf8(out).unwrap(), n)
    }

    #[test]
    fn basic_and_flags() {
        let text = "Error: disk\ninfo: ok\nERROR: net\n";
        assert_eq!(run(&["Error"], text), ("Error: disk\n".into(), 1));
        assert_eq!(run(&["-i", "error"], text).1, 2);
        assert_eq!(run(&["-n", "-v", "-i", "error"], text).0, "2:info: ok\n");
        assert_eq!(run(&["-c", "o"], text), (String::new(), 2)); // case-sensitive: "ERROR" has no 'o'
    }

    #[test]
    fn arg_errors() {
        assert!(parse_args(Vec::<String>::new()).is_err());
        assert!(parse_args(["-z".to_string(), "x".into()]).is_err());
    }
}
```

**Stretch goals:** add regex support with the `regex` crate. It guarantees
linear-time matching, so it is safe against ReDoS (see
[Security in Depth Ch 40A](../security/real-life-security-guide-v1.md#chapter-40a-api-payload-security-in-depth)).
Add recursive directory search with `walkdir`. Then compare your tool with
`ripgrep`, which is the same idea taken to production.

---

## 27. Terminal project: an LRU cache, and why linked lists are hard in Rust

[Go §26](../Golang/real-life-golang-guide.md#26-terminal-project-an-lru-cache-library-and-cli-cache-server)
builds an LRU cache from a map plus `container/list`, a doubly linked list.
In Rust, a doubly linked list with plain references does not compile:
every node would be pointed to by two others (`prev` and `next`), which
gives it two owners and mutable aliasing. The three real options are:

1. `Rc<RefCell<Node>>` plus `Weak` back-pointers. It works, but it is
   verbose, has runtime overhead, and can panic.
2. Raw pointers with `unsafe`, which is what `std::collections::LinkedList`
   does internally.
3. **An arena: store nodes in a `Vec` and link them by index.** This is
   safe, fast, and cache-friendly. The `lru` crate and many ECS game engines
   use this approach.

Option 3 is idiomatic. Indices are "pointers" the borrow checker does not
track. You give up some static checking (a stale index can point to the
wrong node) but never memory safety (a bad index panics, it never reads
freed memory).

```rust
// src/bin/lru.rs
use std::collections::HashMap;
use std::hash::Hash;

const NIL: usize = usize::MAX;

struct Node<K, V> { key: K, val: V, prev: usize, next: usize }

pub struct Lru<K, V> {
    map: HashMap<K, usize>,   // key -> index into `nodes`
    nodes: Vec<Node<K, V>>,   // the arena; slots are reused when entries are evicted
    head: usize,              // most recently used
    tail: usize,              // least recently used
    cap: usize,
}

impl<K: Hash + Eq + Clone, V> Lru<K, V> {
    pub fn new(cap: usize) -> Self {
        assert!(cap > 0, "capacity must be positive");
        Self { map: HashMap::with_capacity(cap), nodes: Vec::with_capacity(cap), head: NIL, tail: NIL, cap }
    }

    fn unlink(&mut self, i: usize) {
        let (p, n) = (self.nodes[i].prev, self.nodes[i].next);
        if p != NIL { self.nodes[p].next = n } else { self.head = n }
        if n != NIL { self.nodes[n].prev = p } else { self.tail = p }
    }

    fn push_front(&mut self, i: usize) {
        self.nodes[i].prev = NIL;
        self.nodes[i].next = self.head;
        if self.head != NIL { self.nodes[self.head].prev = i }
        self.head = i;
        if self.tail == NIL { self.tail = i }
    }

    pub fn get(&mut self, key: &K) -> Option<&V> {
        let i = *self.map.get(key)?;
        self.unlink(i);
        self.push_front(i);
        Some(&self.nodes[i].val)
    }

    /// Inserts or updates; returns the evicted (key, value), if any.
    pub fn put(&mut self, key: K, val: V) -> Option<(K, V)> {
        if let Some(&i) = self.map.get(&key) {
            self.nodes[i].val = val;
            self.unlink(i);
            self.push_front(i);
            return None;
        }
        if self.nodes.len() < self.cap {
            let i = self.nodes.len();
            self.nodes.push(Node { key: key.clone(), val, prev: NIL, next: NIL });
            self.map.insert(key, i);
            self.push_front(i);
            return None;
        }
        // Full: reuse the tail slot for the new entry.
        let i = self.tail;
        self.unlink(i);
        let old_key = std::mem::replace(&mut self.nodes[i].key, key.clone());
        let old_val = std::mem::replace(&mut self.nodes[i].val, val);
        self.map.remove(&old_key);
        self.map.insert(key, i);
        self.push_front(i);
        Some((old_key, old_val))
    }

    pub fn len(&self) -> usize { self.map.len() }

    /// Keys from most to least recently used.
    pub fn keys(&self) -> Vec<&K> {
        let mut out = Vec::with_capacity(self.len());
        let mut i = self.head;
        while i != NIL { out.push(&self.nodes[i].key); i = self.nodes[i].next; }
        out
    }
}

fn main() {
    let mut c = Lru::new(2);
    c.put("a", 1);
    c.put("b", 2);
    c.get(&"a");                               // a is now most recent
    let evicted = c.put("c", 3);               // evicts b
    println!("evicted={evicted:?} keys={:?}", c.keys()); // evicted=Some(("b", 2)) keys=["c", "a"]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn evicts_least_recent() {
        let mut c = Lru::new(3);
        for (k, v) in [(1, "one"), (2, "two"), (3, "three")] { c.put(k, v); }
        assert_eq!(c.get(&1), Some(&"one"));
        assert_eq!(c.put(4, "four"), Some((2, "two")));
        assert_eq!(c.keys(), vec![&4, &1, &3]);
        assert_eq!(c.get(&2), None);
        assert_eq!(c.put(1, "uno"), None);        // update does not evict
        assert_eq!(c.len(), 3);
    }

    #[test]
    fn capacity_one() {
        let mut c = Lru::new(1);
        c.put('x', 0);
        assert_eq!(c.put('y', 1), Some(('x', 0)));
        assert_eq!(c.keys(), vec![&'y']);
    }
}
```

`get` takes `&mut self`, because reading an LRU cache changes its recency
order. To share the cache between threads, wrap it in `Mutex<Lru<K, V>>`
(Ch 29). A read-write lock gives no benefit here, because every read is also
a write. A Go version has the same constraint, but nothing stops a Go caller
from forgetting the lock. In Rust, sharing an `Lru` across threads without
the `Mutex` does not compile.

**Stretch goals:** add a TTL per entry. Shard the cache into N
`Mutex<Lru>` by key hash to reduce lock contention. Benchmark it against the
`lru` crate with `criterion` (Ch 36).

---

# Part III — Advanced: concurrency, async, unsafe

The ownership rules from Part I matter most in this part. The same rule,
shared XOR mutable, that prevents iterator invalidation also makes **data
races a compile error**. This Part also covers async Rust, which is how Rust
handles tens of thousands of connections, and `unsafe`, which is how Rust
talks to C and the kernel.

---

## 28. Threads, and why data races do not compile

### 28.1 Spawning and joining

```rust
use std::thread;
use std::time::Duration;

fn main() {
    let handles: Vec<thread::JoinHandle<u64>> = (0..4)
        .map(|id| {
            thread::spawn(move || {           // `move`: the thread owns its copy of `id`
                thread::sleep(Duration::from_millis(10 * id));
                (1..=1_000_000u64).filter(|n| n % (id + 2) == 0).count() as u64
            })
        })
        .collect();

    // join() returns Result: Err if the thread panicked.
    let results: Vec<u64> = handles.into_iter().map(|h| h.join().expect("worker panicked")).collect();
    println!("{results:?}");
    println!("cores available: {:?}", thread::available_parallelism()); // respects cgroup CPU limits on Linux
}
```

`std::thread` creates **one OS thread per spawn** (1:1), the same as
`pthread_create` in the [C guide §20](../c-lang/real-life-c-guide.md#20-concurrency-posix-threads-mutexes-condvars-c11-threads-h-and-atomics).
Each thread gets about 2 MiB of stack by default (Rust's default, configurable
with `thread::Builder::stack_size`). Go's goroutines are M:N and start at a few KiB
([Go §15](../Golang/real-life-golang-guide.md#15-goroutines-go-s-answer-to-what-is-a-thread)).
Rust's version of cheap tasks is `async` (Ch 30).

**Across the series:** what a thread is to the kernel, and what a context
switch costs, is in [OS Ch 7](../os-linux/real-life-os-guide.md#chapter-7-threads-sharing-everything-except-the-stack)
and [OS Ch 9](../os-linux/real-life-os-guide.md#chapter-9-interrupts-context-switches-and-why-they-re-not-free).

### 28.2 The data race that will not compile

```rust
// DOES NOT COMPILE
use std::thread;

fn main() {
    let mut counter = 0;
    let mut handles = vec![];
    for _ in 0..4 {
        handles.push(thread::spawn(|| {
            for _ in 0..1000 { counter += 1; } // four threads writing one variable
        }));
    }
    for h in handles { h.join().unwrap(); }
    println!("{counter}");
}
```

```text
error[E0499]: cannot borrow `counter` as mutable more than once at a time
error[E0373]: closure may outlive the current function, but it borrows `counter`
```

The C version of this program compiles and prints a different wrong number
on each run ([OS Ch 17](../os-linux/real-life-os-guide.md#chapter-17-race-conditions-at-the-os-level)
shows the same race with shell processes). The Go version compiles, and
`go run -race` catches it **if** the race happens during that run. Rust
rejects it at compile time, because four `&mut counter` borrows would exist
at the same time.

### 28.3 `Send` and `Sync`

Two marker traits, implemented automatically by the compiler, encode thread
safety in the type system:

- **`T: Send`**: a `T` can be **moved** to another thread. Almost
  everything is `Send`. `Rc<T>` is not: its reference count is non-atomic,
  so two threads cloning it would race on the counter.
- **`T: Sync`**: a `&T` can be **shared** between threads, that is, `&T`
  is `Send`. `RefCell<T>` and `Cell<T>` are not `Sync`, because their
  runtime borrow tracking is not atomic. `Mutex<T>` is `Sync`.

`thread::spawn` requires `F: Send + 'static`. That one bound is what makes
the next program fail:

```rust
// DOES NOT COMPILE
use std::rc::Rc;

fn main() {
    let shared = Rc::new(vec![1, 2, 3]);
    let s2 = Rc::clone(&shared);
    std::thread::spawn(move || println!("{s2:?}")).join().unwrap();
}
```

```text
error[E0277]: `Rc<Vec<i32>>` cannot be sent between threads safely
  = help: within `{closure}`, the trait `Send` is not implemented for `Rc<Vec<i32>>`
```

Replace `Rc` with `Arc`, which uses atomic counters, and it compiles. The
compiler does not know anything special about threads. Thread safety
follows from two traits and a bound on `spawn`.

### 28.4 Scoped threads: borrowing across threads

```rust
fn main() {
    let data: Vec<u64> = (1..=10_000).collect();
    let mut partial = [0u64; 4];

    // Threads spawned in a scope are guaranteed to be joined before the scope ends,
    // so they may BORROW local data. No Arc and no 'static needed.
    std::thread::scope(|s| {
        for (chunk, slot) in data.chunks(data.len() / 4).zip(partial.iter_mut()) {
            s.spawn(move || *slot = chunk.iter().sum());
        }
    });
    println!("{partial:?} total={}", partial.iter().sum::<u64>()); // total=50005000
}
```

`chunks` gives each thread a disjoint `&[u64]`, and `iter_mut` gives each
thread its own `&mut u64`. Shared XOR mutable holds, so no locks are needed.
For data-parallel loops, the `rayon` crate turns `.iter()` into
`.par_iter()` and handles the work-splitting for you.

---

## 29. Shared state: `Mutex`, `RwLock`, atomics, `Condvar`, channels

### 29.1 `Mutex<T>` owns the data it protects

```rust
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let hits: Arc<Mutex<HashMap<String, u32>>> = Arc::new(Mutex::new(HashMap::new()));

    let workers: Vec<_> = ["/a", "/b", "/a", "/c", "/a"]
        .into_iter()
        .map(|path| {
            let hits = Arc::clone(&hits);            // one Arc handle per thread
            thread::spawn(move || {
                let mut map = hits.lock().unwrap();  // MutexGuard<HashMap<..>>
                *map.entry(path.to_string()).or_insert(0) += 1;
            })                                       // guard dropped here: unlocked
        })
        .collect();
    for w in workers { w.join().unwrap(); }

    let map = hits.lock().unwrap();
    println!("{:?}", map.get("/a")); // Some(3)
}
```

In C and Go, a mutex protects whatever the programmer *remembers* to access
only while holding it. In Rust, `Mutex<T>` **contains** the `T`. The only
way to reach the data is `lock()`, which returns a guard, and the guard
unlocks when it is dropped. You cannot forget to lock, and you cannot forget
to unlock.

**Poisoning:** if a thread panics while holding the lock, the mutex is
marked poisoned and `lock()` returns `Err`. The data might be half-updated.
`.unwrap()` turns that into a panic. Use `.unwrap_or_else(|e| e.into_inner())`
if you know the data is still consistent.

### 29.2 What Rust does **not** prevent: deadlocks and logical races

```rust
// fragment
// Thread 1                       // Thread 2
let a = accounts.lock();          let b = audit.lock();
let b = audit.lock();   // waits  let a = accounts.lock();   // waits -> DEADLOCK
```

Deadlocks are memory-safe, so the compiler allows them. Use the same
disciplines as the [OS guide Ch 18](../os-linux/real-life-os-guide.md#chapter-18-locks-semaphores-and-deadlock):
a global lock order, holding locks for short scopes, and never calling
unknown code (callbacks, `.await`) while holding a lock. **Check-then-act
races** are also possible: `if map.lock().contains(k) { map.lock().insert(..) }`
takes the lock twice, and another thread can insert in between. Do the whole
operation under one guard.

### 29.3 `RwLock` and atomics

```rust
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Arc, RwLock};
use std::thread;

struct Metrics {
    requests: AtomicU64,                 // lock-free counter
    config: RwLock<String>,              // many readers OR one writer
}

fn main() {
    let m = Arc::new(Metrics { requests: AtomicU64::new(0), config: RwLock::new("v1".into()) });
    thread::scope(|s| {
        for _ in 0..8 {
            s.spawn(|| {
                for _ in 0..10_000 {
                    m.requests.fetch_add(1, Ordering::Relaxed); // a counter needs no ordering (Ch 61)
                    let _cfg = m.config.read().unwrap();         // readers do not block each other
                }
            });
        }
        s.spawn(|| *m.config.write().unwrap() = "v2".into());
    });
    println!("{} requests, config {}", m.requests.load(Ordering::Relaxed), m.config.read().unwrap());
}
```

### 29.4 Channels: "share memory by communicating", Rust edition

```rust
use std::sync::mpsc;
use std::thread;
use std::time::Duration;

fn main() {
    // sync_channel(N) is BOUNDED: senders block when the buffer is full (backpressure).
    let (tx, rx) = mpsc::sync_channel::<(usize, String)>(16);

    for id in 0..3 {
        let tx = tx.clone();                       // multi-producer
        thread::spawn(move || {
            for j in 0..2 {
                tx.send((id, format!("job {j}"))).unwrap();
                thread::sleep(Duration::from_millis(5));
            }
        });                                        // this thread's tx is dropped here
    }
    drop(tx); // drop the original sender, or the receiver loop below never ends

    for (id, msg) in rx {                          // ends when ALL senders are dropped
        println!("worker {id}: {msg}");
    }
}
```

| Go | Rust (`std`) | Rust (Tokio, Ch 31) |
|---|---|---|
| `make(chan T)` (unbuffered) | `sync_channel(0)` | — |
| `make(chan T, n)` | `sync_channel(n)` | `tokio::sync::mpsc::channel(n)` |
| `close(ch)` | drop every `Sender` | drop every `Sender` |
| `select` | — (use `crossbeam-channel`) | `tokio::select!` |
| broadcast | — | `tokio::sync::broadcast` |
| "latest value" | — | `tokio::sync::watch` |

Sending a value **moves** it. After `tx.send(buf)`, the sending thread can
no longer touch `buf`. In Go, sending a pointer or slice over a channel and
then mutating it is a common race. In Rust that is a compile error.

### 29.5 `Condvar`: waiting for a condition

`Condvar` works like `pthread_cond_t` in the C guide's producer/consumer
example: wait while holding a `MutexGuard` and loop on the condition,
because wakeups can be spurious. Most code should use a channel instead.
Ch 37 uses one inside a thread pool.

---

## 30. Async Rust: futures, `poll`, and why async exists

### 30.1 Why: the cost of a thread per connection

A server that spawns a thread per connection pays a few KiB to MiB of stack
and a kernel context switch per connection. At 10,000 mostly idle
connections, the threads cost far more than the work. The kernel's answer
is readiness notification, `epoll` on Linux and `kqueue` on macOS: one
thread asks "which of these 10,000 sockets is ready?" and handles only those.
The [C guide's `select()` chat server](../c-lang/real-life-c-guide.md#23-mini-projects-a-thread-pool-and-a-select-based-chat-server)
does this by hand. Go hides it inside the runtime's netpoller
([OS Ch 78](../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections)).
Rust exposes it as `async`/`.await`, with no built-in runtime. You pick one,
usually **Tokio**.

### 30.2 A `Future` is a state machine you poll

```rust
// fragment: the real trait in std::future
pub trait Future {
    type Output;
    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output>;
}
pub enum Poll<T> { Ready(T), Pending }
```

- `async fn fetch() -> String` does **not** run when you call it. It returns
  a future, a value that will produce a `String` when it is polled to
  completion.
- `.await` polls the inner future. If it returns `Pending`, the current
  future also returns `Pending` and gives the thread back to the executor.
- When the socket becomes ready, the reactor (epoll/kqueue) calls the
  **`Waker`** from `cx`. The executor then polls that task again.
- The compiler turns each `async fn` into an enum, with one variant per
  `.await` point holding the variables that are live there. Ch 59 shows the
  generated shape and builds an executor.

```rust
use std::time::{Duration, Instant};

async fn step(name: &str, ms: u64) -> String {
    tokio::time::sleep(Duration::from_millis(ms)).await; // yields; does not block the thread
    format!("{name} done")
}

#[tokio::main]
async fn main() {
    let t = Instant::now();
    let f = step("lazy", 10);                    // nothing has happened yet
    println!("created future after {:?}", t.elapsed());
    println!("{}", f.await);                     // now it runs

    let t = Instant::now();
    let (a, b, c) = tokio::join!(step("a", 100), step("b", 100), step("c", 100));
    println!("{a}, {b}, {c} in {:?}", t.elapsed()); // ~100ms, not 300ms: concurrent on ONE task
}
```

### 30.3 Go goroutines vs Rust async

| | Go goroutine | Rust async task |
|---|---|---|
| Runtime | built in, always present | a library you choose (Tokio, smol, embassy for embedded) |
| Stack | growable stack per goroutine | **no stack**: the future stores exactly the state it needs across `.await` points |
| Preemption | yes (async preemption since Go 1.14) | **cooperative**: a task runs until it hits `.await` |
| Blocking call inside | fine, the runtime adds threads | **stalls a worker thread**: use `spawn_blocking` |
| Function color | none; any function can block | `async fn` can only be awaited from async code |
| Cancellation | `context.Context`, by convention | drop the future; it stops at its next `.await` |

The cost of async Rust is the "function color" split and the cooperative
model. The benefit is that a task is a small heap object of known size,
often a few hundred bytes, with no stack or GC. One process can hold a
million of them.

---

## 31. Tokio in practice: tasks, `select!`, timeouts, cancellation, shutdown

### 31.1 Tasks, `JoinSet`, and `spawn_blocking`

```rust
use std::time::Duration;
use tokio::task::JoinSet;

async fn check(host: &'static str) -> Result<u128, String> {
    let start = std::time::Instant::now();
    tokio::time::sleep(Duration::from_millis(host.len() as u64 * 10)).await; // stand-in for I/O
    if host.contains("bad") { return Err(format!("{host}: unreachable")); }
    Ok(start.elapsed().as_millis())
}

#[tokio::main]
async fn main() {
    let mut set = JoinSet::new();
    for host in ["api.local", "db.local", "bad.local", "cache.local"] {
        set.spawn(async move { (host, check(host).await) }); // 'static + Send, like thread::spawn
    }
    while let Some(res) = set.join_next().await {            // results in COMPLETION order
        let (host, outcome) = res.expect("task panicked");
        println!("{host:<12} {outcome:?}");
    }

    // CPU-heavy or blocking work goes to a separate thread pool:
    let digest = tokio::task::spawn_blocking(|| (0..5_000_000u64).fold(0u64, |a, x| a.wrapping_add(x * x)))
        .await
        .unwrap();
    println!("digest {digest}");
}
```

This is Go's [URL health-checker worker pool](../Golang/real-life-golang-guide.md#25-terminal-project-a-concurrent-url-health-checker-worker-pool),
without the WaitGroup.

### 31.2 Timeouts and `select!`

```rust
use std::time::Duration;
use tokio::sync::mpsc;
use tokio::time::{sleep, timeout};

#[tokio::main]
async fn main() {
    // timeout() wraps any future; when time runs out, the inner future is DROPPED (cancelled).
    match timeout(Duration::from_millis(50), sleep(Duration::from_secs(10))).await {
        Ok(()) => println!("finished"),
        Err(_) => println!("timed out after 50ms"),
    }

    let (tx, mut rx) = mpsc::channel::<String>(8);
    tokio::spawn(async move {
        for i in 0..3 {
            tx.send(format!("msg {i}")).await.unwrap();
            sleep(Duration::from_millis(20)).await;
        }
    });

    let mut ticker = tokio::time::interval(Duration::from_millis(25));
    loop {
        tokio::select! {                      // wait on several futures; the first one ready wins
            Some(m) = rx.recv() => println!("got {m}"),
            _ = ticker.tick() => println!("tick"),
            else => break,                    // all branches disabled: channel closed
        }
        if rx.is_closed() && rx.is_empty() { break; }
    }
}
```

### 31.3 Cancellation is "drop the future"

In Go, cancellation is cooperative through `ctx.Done()`
([Go §18](../Golang/real-life-golang-guide.md#18-context-cancellation-deadlines-and-request-scoped-values)),
and a function that ignores `ctx` keeps running. In Rust, a future that is
dropped **stops at its current `.await` point** and runs `Drop` for
everything it holds. When `timeout` or `select!` picks another branch, the
losing future is dropped. That is simpler than Go's model, but it has one
sharp edge, **cancellation safety**:

```rust
// fragment
// BUG: if the timeout fires between the two writes, the peer sees half a message.
tokio::select! {
    _ = async { sock.write_all(header).await?; sock.write_all(body).await } => {}
    _ = sleep(deadline) => {}  // the other branch is dropped mid-write
}
```

The rule: any future you put in `select!` or `timeout` may stop at any
`.await`. Keep partial state somewhere that survives, such as the
connection's buffer, or encode and write the whole frame in one call. Tokio's
documentation lists which methods are cancel-safe (for example,
`mpsc::Receiver::recv` is; `AsyncReadExt::read_exact` is not).

### 31.4 Graceful shutdown

```rust
use std::time::Duration;
use tokio::sync::watch;
use tokio::task::JoinSet;

async fn worker(id: u32, mut shutdown: watch::Receiver<bool>) {
    loop {
        tokio::select! {
            _ = shutdown.changed() => {
                println!("worker {id}: draining and exiting");
                return;
            }
            _ = tokio::time::sleep(Duration::from_millis(30)) => println!("worker {id}: tick"),
        }
    }
}

#[tokio::main]
async fn main() {
    let (tx, rx) = watch::channel(false);
    let mut workers = JoinSet::new();
    for id in 0..2 { workers.spawn(worker(id, rx.clone())); }

    // In a real server: tokio::signal::ctrl_c().await, or a SIGTERM stream (Ch 40, 67).
    tokio::time::sleep(Duration::from_millis(70)).await;
    tx.send(true).unwrap();                           // tell everyone

    // Give them a deadline to finish, like Kubernetes' terminationGracePeriodSeconds.
    let drained = tokio::time::timeout(Duration::from_secs(5), async {
        while workers.join_next().await.is_some() {}
    }).await;
    println!("clean shutdown: {}", drained.is_ok());
}
```

`tokio_util::sync::CancellationToken` packages this pattern as a
tree-shaped cancel signal, the closest equivalent to `context.WithCancel`.

### 31.5 The two async mistakes everyone makes

1. **Blocking the executor.** `std::thread::sleep`, `std::fs::read` of a
   large file, heavy CPU work, or a blocking database driver inside an
   `async fn` stalls every task on that worker thread. Use
   `tokio::time::sleep`, `tokio::fs`, or `spawn_blocking`.
2. **Holding a `std::sync::MutexGuard` across `.await`.** The guard is not
   `Send`, so `tokio::spawn` rejects the future with a long error that
   mentions "future cannot be sent between threads safely". Even when it
   compiles, holding a lock across a suspension point invites deadlocks.
   Copy what you need out of the guard and drop it before `.await`. Use
   `tokio::sync::Mutex` only when you really must hold a lock across `.await`.

---

## 32. Networking with `std` and Tokio: TCP, UDP, and framing

### 32.1 A threaded TCP echo server with `std`

```rust
use std::io::{BufRead, BufReader, Write};
use std::net::{TcpListener, TcpStream};
use std::time::Duration;

fn handle(stream: TcpStream) -> std::io::Result<()> {
    let peer = stream.peer_addr()?;
    stream.set_read_timeout(Some(Duration::from_secs(30)))?; // idle clients cannot hold threads forever
    stream.set_nodelay(true)?;                               // disable Nagle for small interactive writes
    let mut writer = stream.try_clone()?;
    for line in BufReader::new(stream).lines() {
        let line = line?;
        writeln!(writer, "echo: {line}")?;
    }
    println!("{peer} disconnected");
    Ok(())
}

fn main() -> std::io::Result<()> {
    let listener = TcpListener::bind("127.0.0.1:7878")?;  // socket() + bind() + listen()
    println!("listening on {}", listener.local_addr()?);
    for stream in listener.incoming().take(1) {           // take(1) so the demo exits; remove for a real server
        let stream = stream?;                              // accept()
        std::thread::spawn(move || {
            if let Err(e) = handle(stream) { eprintln!("connection error: {e}"); }
        });
    }
    Ok(())
}
```

```bash
cargo run --bin echo &       # then:
printf 'hello\nworld\n' | nc 127.0.0.1 7878
```

Every call above maps one-to-one onto the
[C guide's §21.1 echo server](../c-lang/real-life-c-guide.md#21-sockets-a-tcp-echo-server-then-a-tiny-http-server):
`bind` performs `socket`/`setsockopt(SO_REUSEADDR)`/`bind`/`listen`, and
`incoming` loops over `accept`. The differences: the socket is closed by
`Drop`, errors are `Result`s instead of `-1` and `errno`, and the buffer
cannot overflow.

**Across the series:** what `bind`, `listen`, and `accept` do on the wire is
in [Net Ch 19](../networking/tcp-ip/real-life-guide-v1.md#chapter-19-ports-and-sockets-which-program-gets-the-data)
and [Ch 21](../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake).
Why `set_nodelay` matters is [Ch 26, the 40-millisecond mystery](../networking/tcp-ip/real-life-guide-v1.md#chapter-26-the-40-millisecond-mystery).

### 32.2 The same server with Tokio, plus length-prefixed framing

TCP is a **byte stream**, not a message stream. One `write` of 100 bytes
can arrive as two `read`s of 60 and 40, and two writes can arrive in one
read. Every protocol therefore needs **framing**: delimiters (`\n`, as in
HTTP/1 headers and Redis) or a length prefix (as in HTTP/2, TLS records,
gRPC, and Kafka).

```rust
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};

const MAX_FRAME: u32 = 1 << 20; // 1 MiB: never trust a length field from the network

async fn read_frame(s: &mut TcpStream) -> std::io::Result<Option<Vec<u8>>> {
    let len = match s.read_u32().await {              // 4-byte big-endian length
        Ok(n) => n,
        Err(e) if e.kind() == std::io::ErrorKind::UnexpectedEof => return Ok(None), // clean close
        Err(e) => return Err(e),
    };
    if len > MAX_FRAME {
        return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, format!("frame of {len} bytes")));
    }
    let mut buf = vec![0u8; len as usize];
    s.read_exact(&mut buf).await?;                    // loops until all `len` bytes have arrived
    Ok(Some(buf))
}

async fn write_frame(s: &mut TcpStream, payload: &[u8]) -> std::io::Result<()> {
    let len = u32::try_from(payload.len()).map_err(|_| std::io::Error::other("payload too large"))?;
    let mut out = Vec::with_capacity(4 + payload.len());
    out.extend_from_slice(&len.to_be_bytes());
    out.extend_from_slice(payload);
    s.write_all(&out).await                           // one write: cancel-safe at the frame level
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let listener = TcpListener::bind("127.0.0.1:0").await?;
    let addr = listener.local_addr()?;

    tokio::spawn(async move {
        loop {
            let (mut sock, _) = listener.accept().await.unwrap();
            tokio::spawn(async move {                 // one TASK per connection, not one thread
                while let Ok(Some(frame)) = read_frame(&mut sock).await {
                    let reply = frame.to_ascii_uppercase();
                    if write_frame(&mut sock, &reply).await.is_err() { break; }
                }
            });
        }
    });

    let mut client = TcpStream::connect(addr).await?;
    for msg in ["hello", "framed", "world"] {
        write_frame(&mut client, msg.as_bytes()).await?;
        let reply = read_frame(&mut client).await?.expect("server replied");
        println!("{}", String::from_utf8_lossy(&reply));
    }
    Ok(())
}
```

`MAX_FRAME` is a security control. Without it, a client sending the four
bytes `FF FF FF FF` makes the server allocate 4 GiB. Ch 53 treats every
length field from the network this way.

### 32.3 UDP

```rust
use std::net::UdpSocket;
use std::time::{Duration, Instant};

fn main() -> std::io::Result<()> {
    let server = UdpSocket::bind("127.0.0.1:0")?;
    let addr = server.local_addr()?;
    std::thread::spawn(move || {
        let mut buf = [0u8; 1500];                         // one datagram, up to a typical MTU
        while let Ok((n, peer)) = server.recv_from(&mut buf) {
            let _ = server.send_to(&buf[..n], peer);       // echo
        }
    });

    let client = UdpSocket::bind("127.0.0.1:0")?;
    client.set_read_timeout(Some(Duration::from_millis(500)))?; // UDP never tells you about loss
    let mut lost = 0;
    for seq in 0u32..5 {
        let t = Instant::now();
        client.send_to(&seq.to_be_bytes(), addr)?;
        let mut buf = [0u8; 4];
        match client.recv_from(&mut buf) {
            Ok(_) if u32::from_be_bytes(buf) == seq => println!("seq={seq} rtt={:?}", t.elapsed()),
            _ => lost += 1,
        }
    }
    println!("lost {lost}/5");
    Ok(())
}
```

This is a small version of the networking guide's UDP loss meter lab
([Net Ch 20](../networking/tcp-ip/real-life-guide-v1.md#chapter-20-udp-fire-and-forget)).

---

## 33. `unsafe` Rust: the five superpowers and the contract they carry

### 33.1 What `unsafe` unlocks

An `unsafe { }` block allows exactly five extra operations:

1. Dereferencing a raw pointer (`*const T`, `*mut T`).
2. Calling an `unsafe fn`, including every FFI function.
3. Reading or writing a `static mut`.
4. Implementing an `unsafe trait` (such as `Send` or `Sync`).
5. Reading a field of a `union`.

It does **not** turn off the borrow checker or type checking for ordinary
references. It means: "I, the programmer, have checked the preconditions
the compiler cannot check here." If you are wrong, the result is
**undefined behavior**, the same category of bug as in the
[C guide §28](../c-lang/real-life-c-guide.md#28-undefined-behavior-the-list-every-c-programmer-must-memorize).

### 33.2 Safe abstractions over unsafe code

The standard library is full of `unsafe`, wrapped in APIs that cannot be
misused. Here is the real idea behind `split_at_mut`, which the borrow
checker cannot verify on its own because it cannot tell that the two halves
do not overlap:

```rust
fn split_at_mut<T>(s: &mut [T], mid: usize) -> (&mut [T], &mut [T]) {
    let len = s.len();
    assert!(mid <= len, "mid out of bounds");   // the check that makes the unsafe below sound
    let ptr = s.as_mut_ptr();
    // SAFETY: `mid <= len`, so both ranges are inside the original allocation,
    // and [0, mid) and [mid, len) do not overlap, so the two &mut never alias.
    unsafe {
        (
            std::slice::from_raw_parts_mut(ptr, mid),
            std::slice::from_raw_parts_mut(ptr.add(mid), len - mid),
        )
    }
}

fn main() {
    let mut v = [1, 2, 3, 4, 5];
    let (a, b) = split_at_mut(&mut v, 2);
    a[0] = 10;
    b[0] = 30;
    println!("{v:?}"); // [10, 2, 30, 4, 5]
}
```

The pattern is: **check the invariant in safe code, write a `// SAFETY:`
comment that justifies each `unsafe` operation, and expose a safe function.**
Callers cannot cause UB however they call it. Clippy's
`undocumented_unsafe_blocks` lint enforces the comment.

### 33.3 UB in Rust is a short list, but it is real

- Dereferencing a dangling, null, or misaligned pointer.
- **Breaking the aliasing rules:** creating a `&mut T` while another
  reference to the same data is in use, even through raw pointers.
- Producing an invalid value: a `bool` that is not 0 or 1, a `char` outside
  Unicode, a null `&T`, a `str` that is not UTF-8, an enum with an invalid
  tag.
- Data races, which are possible only through `unsafe` or a wrong
  `unsafe impl Send`.
- Calling a foreign function with the wrong signature.

### 33.4 Miri: an interpreter that detects UB

```bash
rustup +nightly component add miri
cargo +nightly miri test        # runs your tests in an interpreter that checks every memory access
```

Miri catches use-after-free, out-of-bounds access, aliasing violations, and
memory leaks in code with `unsafe`. Think of it as AddressSanitizer plus
UBSan ([C §14](../c-lang/real-life-c-guide.md#14-debugging-and-sanitizers-gdb-lldb-valgrind-asan-ubsan))
with knowledge of Rust's aliasing model. Run it in CI for every crate that
contains `unsafe`.

**War story.** A popular crate's `unsafe impl Send` on a type containing an
`Rc` passed review and tests for months. When users began moving it between
Tokio worker threads, the non-atomic reference count raced, and objects were
occasionally freed twice. The RustSec advisory database has many entries
like this, all in `unsafe` code. Safe Rust cannot express this bug, so
reviewing a crate's `unsafe` blocks covers the part of the code where it
can occur. `cargo geiger` counts the unsafe code in your dependency tree.

---

## 34. FFI: calling C from Rust and Rust from C

### 34.1 Calling C

```rust
use std::ffi::{c_char, c_int, CStr, CString};

unsafe extern "C" {                            // edition 2024: extern blocks are `unsafe extern`
    fn strlen(s: *const c_char) -> usize;      // from libc, linked by default
    fn getpid() -> c_int;
    safe fn abs(x: c_int) -> c_int;            // declared `safe`: callable without an unsafe block
}

fn main() {
    let owned = CString::new("hello, C").expect("no interior NUL bytes"); // adds the trailing \0
    // SAFETY: `owned` is a valid NUL-terminated string that lives across the call.
    let n = unsafe { strlen(owned.as_ptr()) };
    // SAFETY: getpid has no preconditions.
    let pid = unsafe { getpid() };
    println!("strlen={n} pid={pid} abs={}", abs(-7));

    // C string -> Rust: borrow with CStr, then validate as UTF-8.
    let back: &CStr = owned.as_c_str();
    println!("{:?}", back.to_str());
}
```

| Rust | C |
|---|---|
| `CString` / `&CStr` | owned / borrowed `char *` with a NUL terminator |
| `*const T` / `*mut T` | `const T *` / `T *` |
| `#[repr(C)] struct` | `struct` with C layout (Ch 57) |
| `Option<extern "C" fn(..)>` | nullable function pointer |
| `std::ffi::c_int`, `c_long`, ... | `int`, `long`, ... (platform sizes) |

The [`libc`](https://crates.io/crates/libc) crate declares the whole C
library for each platform, so you rarely write declarations like the ones
above yourself. `bindgen` generates them from a C header automatically.

### 34.2 Building and linking your own C code

To reuse code from the C guide, such as its
[JSON parser](../c-lang/real-life-c-guide.md#16-mini-projects-a-dynamic-array-library-csv-parser-and-a-json-parser),
compile it with the `cc` crate from a build script:

```rust
// fragment: build.rs  (Cargo.toml: [build-dependencies] cc = "1")
fn main() {
    cc::Build::new().file("c/json.c").flag_if_supported("-std=c17").compile("cjson"); // libcjson.a
    println!("cargo::rerun-if-changed=c/json.c");
}
```

### 34.3 Exposing Rust to C: the opaque-pointer pattern

[C §40.1](../c-lang/real-life-c-guide.md#40-abi-linking-and-designing-a-stable-c-api)
designs a stable C API around opaque handles (`CURL *`). In Rust you create
the handle with `Box::into_raw` and destroy it with `Box::from_raw`:

```rust
use std::ffi::{c_char, CStr};

pub struct Counter { hits: u64, last: String }

#[unsafe(no_mangle)]                                     // keep the symbol name `counter_new`
pub extern "C" fn counter_new() -> *mut Counter {
    Box::into_raw(Box::new(Counter { hits: 0, last: String::new() })) // ownership passes to C
}

/// # Safety
/// `c` must come from `counter_new`, not be freed, and not be used concurrently.
/// `key` must be a valid NUL-terminated string.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn counter_hit(c: *mut Counter, key: *const c_char) -> u64 {
    // Never let a panic unwind into C: catch it and return an error value.
    std::panic::catch_unwind(std::panic::AssertUnwindSafe(|| {
        // SAFETY: the caller upholds the contract documented above.
        let (c, key) = unsafe { (&mut *c, CStr::from_ptr(key)) };
        c.hits += 1;
        c.last = key.to_string_lossy().into_owned();
        c.hits
    }))
    .unwrap_or(u64::MAX)
}

/// # Safety
/// `c` must come from `counter_new` and must not be used after this call.
#[unsafe(no_mangle)]
pub unsafe extern "C" fn counter_free(c: *mut Counter) {
    if !c.is_null() {
        // SAFETY: per the contract, `c` is a live pointer from Box::into_raw.
        drop(unsafe { Box::from_raw(c) });               // ownership returns to Rust; Drop runs
    }
}

fn main() {
    let c = counter_new();
    let key = std::ffi::CString::new("GET /").unwrap();
    // SAFETY: `c` is live, `key` is NUL-terminated.
    let n = unsafe { counter_hit(c, key.as_ptr()) };
    // SAFETY: last use of `c`.
    unsafe { counter_free(c) };
    println!("hits={n}");
}
```

Build it as a C library with `crate-type = ["cdylib", "staticlib"]` under
`[lib]`, and generate the header with `cbindgen`. This is how
Rust components go into C and C++ codebases: rustls's C API, librsvg,
the Rust code in Firefox, and Rust in the Linux kernel.

---

## 35. Macros: `macro_rules!` and what derive macros do

Rust macros operate on **syntax trees**, not text, so the problems of the
[C preprocessor](../c-lang/real-life-c-guide.md#8-structs-unions-enums-typedef-and-the-preprocessor)
(double evaluation, missing parentheses, name capture) do not occur.

```rust
use std::collections::HashMap;

// A literal syntax for maps: hashmap!{"a" => 1, "b" => 2}
macro_rules! hashmap {
    ($($k:expr => $v:expr),* $(,)?) => {{          // $(...),* repeats; $(,)? allows a trailing comma
        let mut m = HashMap::new();
        $( m.insert($k, $v); )*
        m
    }};
}

// Generate repetitive code: one newtype per ID kind, so a UserId cannot be passed as an OrderId.
macro_rules! id_type {
    ($($name:ident),+) => {
        $(
            #[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
            pub struct $name(pub u64);
        )+
    };
}
id_type!(UserId, OrderId);

fn main() {
    let ports = hashmap! { "http" => 80, "https" => 443, };
    let u = UserId(7);
    let o = OrderId(7);
    // assert_eq!(u, o);  // compile error: different types
    println!("{:?} {u:?} {o:?}", ports.get("https"));
}
```

**Procedural macros** are Rust programs that run inside the compiler and
take tokens in and produce tokens out. You use them all the time:
`#[derive(Debug, Serialize)]`, `#[tokio::main]`, `#[test]`. `cargo expand`
shows what they generate. Write your own only for real boilerplate across a
codebase. They add compile time, and in a dependency they are code that
**runs on your machine at build time**, which matters for supply-chain
security (Ch 55).

---

## 36. Performance: release profiles, benchmarks, profiling, allocations

### 36.1 Always measure release builds

Debug builds are often 10 to 50 times slower than release builds for
iterator-heavy code. Never benchmark `cargo run` without `--release`.

```toml
[profile.release]
lto = "fat"            # whole-program optimization; slower builds
codegen-units = 1
debug = "line-tables-only"   # symbols for perf/flamegraph without much size cost

[profile.bench]
inherits = "release"
```

### 36.2 Micro-benchmarking without fooling yourself

```rust
use std::hint::black_box;
use std::time::Instant;

fn sum_index(v: &[u64]) -> u64 { let mut s = 0; for i in 0..v.len() { s += v[i]; } s }
fn sum_iter(v: &[u64]) -> u64 { v.iter().sum() }

fn bench(name: &str, f: impl Fn() -> u64) {
    for _ in 0..3 { black_box(f()); }                  // warm-up: caches, branch predictor, CPU frequency
    let runs = 50;
    let t = Instant::now();
    for _ in 0..runs { black_box(f()); }               // black_box stops the optimizer from deleting the work
    println!("{name:<10} {:?}/run", t.elapsed() / runs);
}

fn main() {
    let v: Vec<u64> = (0..2_000_000).collect();
    bench("index", || sum_index(black_box(&v)));
    bench("iterator", || sum_iter(black_box(&v)));
}
```

On an Apple Silicon laptop this printed about 11 ms (index) and 7 ms
(iterator) per run in a debug build, and about 0.2 ms for **both** with
`--release`. The debug build was roughly 40 times slower, and in release
LLVM proved `i < v.len()` and removed the bounds check, so the two versions
compiled to nearly the same loop. A single run like this is noisy. Run it
several times before you trust any difference.

For real benchmarks use **`criterion`** (statistics, outlier detection,
comparison against a saved baseline) or **`divan`**. The methodology is the
same as [Go §37](../Golang/real-life-golang-guide.md#37-performance-tuning-pprof-trace-and-benchmark-methodology):
warm up, repeat, compare distributions rather than single runs, and change
one thing at a time.

### 36.3 Profiling

```bash
cargo install flamegraph
cargo flamegraph --bin api          # perf on Linux, dtrace/xctrace on macOS -> flamegraph.svg
perf record -g --call-graph dwarf target/release/api && perf report   # Linux, by hand
```

These are the same tools as [OS Ch 81](../os-linux/real-life-os-guide.md#chapter-81-profiling-go-on-linux-pprof-the-execution-tracer-and-perf)
and [C §42.1](../c-lang/real-life-c-guide.md#42-profiling-perf-cache-aware-layout-and-false-sharing).
A Rust binary is a native binary, so everything that works for C works here.
Use `dhat` or `heaptrack` for allocation profiling.

### 36.4 The usual wins, in the order they usually pay off

1. **Algorithm and I/O.** Buffer your writes (Ch 24), batch syscalls, and
   stream instead of loading whole files.
2. **Allocations.** `Vec::with_capacity`, reuse buffers (`buf.clear()`
   keeps the capacity), use `&str` instead of `String` in hot paths, and
   `Cow<str>` when you only sometimes need to allocate.
3. **Clones.** Search the hot path for `.clone()` and `.to_string()`.
4. **Data layout.** Use contiguous `Vec`s instead of pointer-chasing
   (`Vec<Box<T>>`, linked lists), and struct-of-arrays for scans
   ([C §42.2](../c-lang/real-life-c-guide.md#42-profiling-perf-cache-aware-layout-and-false-sharing)).
   Pad hot atomics to separate cache lines to avoid false sharing
   (`crossbeam_utils::CachePadded`).
5. **Allocator.** Switching the global allocator to `mimalloc` or
   `jemalloc` is often a 5–20% gain for allocation-heavy servers (Ch 43
   shows how to install one).
6. **Target CPU.** `RUSTFLAGS="-C target-cpu=native"` enables AVX2 and
   similar instructions, but the binary then runs only on CPUs like the
   build machine.

---

## 37. Project: a thread pool from scratch

The same design as the [C guide's §23.1 thread pool](../c-lang/real-life-c-guide.md#23-mini-projects-a-thread-pool-and-a-select-based-chat-server):
N workers take jobs from a shared queue. Compare the shutdown logic: in Rust
it lives in `Drop`, so a pool cannot be leaked with its threads still
running.

```rust
// src/bin/pool.rs
use std::sync::mpsc::{self, Receiver, Sender};
use std::sync::{Arc, Mutex};
use std::thread::{self, JoinHandle};

type Job = Box<dyn FnOnce() + Send + 'static>;   // any closure that can move to a thread and run once

pub struct ThreadPool {
    sender: Option<Sender<Job>>,                 // Option so Drop can take it and close the channel
    workers: Vec<JoinHandle<()>>,
}

impl ThreadPool {
    pub fn new(size: usize) -> Self {
        assert!(size > 0);
        let (sender, receiver) = mpsc::channel::<Job>();
        let receiver: Arc<Mutex<Receiver<Job>>> = Arc::new(Mutex::new(receiver));
        let workers = (0..size)
            .map(|id| {
                let rx = Arc::clone(&receiver);
                thread::Builder::new()
                    .name(format!("pool-{id}"))
                    .spawn(move || loop {
                        // The guard is a temporary, dropped at the end of this statement,
                        // so the lock is NOT held while the job runs.
                        let msg = rx.lock().unwrap_or_else(|e| e.into_inner()).recv();
                        match msg {
                            Ok(job) => {
                                // Isolate panics: a bad job must not kill the worker.
                                if std::panic::catch_unwind(std::panic::AssertUnwindSafe(job)).is_err() {
                                    eprintln!("pool-{id}: job panicked");
                                }
                            }
                            Err(_) => break, // channel closed: shut down
                        }
                    })
                    .expect("spawn worker")
            })
            .collect();
        Self { sender: Some(sender), workers }
    }

    pub fn execute<F: FnOnce() + Send + 'static>(&self, f: F) {
        self.sender.as_ref().expect("pool is running").send(Box::new(f)).expect("workers alive");
    }
}

impl Drop for ThreadPool {
    fn drop(&mut self) {
        drop(self.sender.take());           // close the channel: every recv() now returns Err
        for w in self.workers.drain(..) {
            let _ = w.join();               // wait for queued jobs to finish
        }
    }
}

fn main() {
    let results = Arc::new(Mutex::new(Vec::new()));
    {
        let pool = ThreadPool::new(4);
        for i in 0..8u64 {
            let results = Arc::clone(&results);
            pool.execute(move || {
                if i == 5 { panic!("job {i} failed"); }
                let fib = (0..i).fold((0u64, 1u64), |(a, b), _| (b, a + b)).0;
                results.lock().unwrap().push((i, fib));
            });
        }
    } // pool dropped: waits for all jobs
    let mut r = results.lock().unwrap().clone();
    r.sort();
    println!("{r:?}"); // 7 results; job 5 panicked and was isolated
}
```

```text
thread 'pool-3' panicked at src/bin/pool.rs:65:29:
job 5 failed
pool-3: job panicked
[(0, 0), (1, 1), (2, 1), (3, 2), (4, 3), (6, 8), (7, 13)]
```

**Stretch goals:** use a bounded `sync_channel` so `execute` applies
backpressure. Return a handle from `execute` that gives back the job's
result (a one-shot channel). Then compare with `rayon::ThreadPool`.

---

## 38. Project: a multi-room chat server with Tokio

The [C guide](../c-lang/real-life-c-guide.md#23-mini-projects-a-thread-pool-and-a-select-based-chat-server)
builds this with `select()` and
[Go §38](../Golang/real-life-golang-guide.md#38-terminal-project-a-tcp-chat-server-with-rooms)
with goroutines and channels. In Tokio each connection is a task, and each
room is a `broadcast` channel.

```rust
// src/bin/chat.rs    connect with: nc 127.0.0.1 6000
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use tokio::io::{AsyncBufReadExt, AsyncWriteExt, BufReader};
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::broadcast;

type Rooms = Arc<Mutex<HashMap<String, broadcast::Sender<(u64, String)>>>>;

fn room(rooms: &Rooms, name: &str) -> broadcast::Sender<(u64, String)> {
    let mut map = rooms.lock().unwrap();      // std Mutex is fine: never held across .await
    map.entry(name.to_string()).or_insert_with(|| broadcast::channel(256).0).clone()
}

async fn client(stream: TcpStream, id: u64, rooms: Rooms) -> std::io::Result<()> {
    let (reader, mut writer) = stream.into_split();
    let mut lines = BufReader::new(reader).lines();
    let mut nick = format!("guest{id}");
    let mut tx = room(&rooms, "lobby");
    let mut rx = tx.subscribe();
    writer.write_all(b"welcome. commands: /nick NAME, /join ROOM, /quit\n").await?;

    loop {
        tokio::select! {
            line = lines.next_line() => {
                let Some(line) = line? else { break };            // EOF: client left
                let line = line.trim();
                if line.len() > 512 { writer.write_all(b"line too long\n").await?; continue; }
                match line.split_once(' ') {
                    Some(("/nick", n)) if !n.is_empty() => nick = n.chars().take(24).collect(),
                    Some(("/join", r)) if !r.is_empty() => {
                        tx = room(&rooms, r);
                        rx = tx.subscribe();                      // the old receiver is dropped
                        writer.write_all(format!("joined {r}\n").as_bytes()).await?;
                    }
                    _ if line == "/quit" => break,
                    _ if !line.is_empty() => { let _ = tx.send((id, format!("{nick}: {line}"))); }
                    _ => {}
                }
            }
            msg = rx.recv() => match msg {
                Ok((from, text)) if from != id => writer.write_all(format!("{text}\n").as_bytes()).await?,
                Ok(_) => {}                                       // our own message
                Err(broadcast::error::RecvError::Lagged(n)) => {  // slow client: skip, do not block others
                    writer.write_all(format!("[missed {n} messages]\n").as_bytes()).await?;
                }
                Err(broadcast::error::RecvError::Closed) => break,
            },
        }
    }
    Ok(())
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let listener = TcpListener::bind("127.0.0.1:6000").await?;
    let rooms: Rooms = Arc::default();
    println!("chat on {}", listener.local_addr()?);
    let mut next_id = 0u64;
    loop {
        let (stream, peer) = listener.accept().await?;
        next_id += 1;
        let (id, rooms) = (next_id, Arc::clone(&rooms));
        tokio::spawn(async move {
            if let Err(e) = client(stream, id, rooms).await { eprintln!("{peer}: {e}"); }
        });
    }
}
```

Design notes that carry over to every real server:

- **Slow consumers.** `broadcast` has a fixed buffer. A client that cannot
  keep up gets `Lagged` and skips messages, so one slow client never stalls
  the room. Go chat servers often get this wrong with unbuffered channels,
  and one stuck client then blocks everyone.
- **Input limits.** Line length and nickname length are capped. Note that
  `next_line` still reads a whole line into memory before checking it, so a
  client can send a 1 GB line. For a hostile network, use
  `tokio_util::codec::LinesCodec::new_with_max_length`.
- **Locks and `.await`.** The `std::sync::Mutex` is held only inside
  `room()`, which never awaits. That is why a std mutex is correct here
  (§31.5).

---

# Part IV — Systems Rust: the OS guide in Rust

The [OS guide's Part 20](../os-linux/real-life-os-guide.md#part-20-systems-programming-in-go-the-os-from-inside-a-program)
looks at the kernel from inside a Go program. This Part does the same in
Rust. Rust has no runtime between you and the kernel: no scheduler, no GC,
no netpoller. `strace` of a Rust binary shows almost exactly the syscalls
your code makes. Chapters 42 and 44 need Linux. Use the container from §0.1.

---

## 39. Syscalls from Rust: `std`, `libc`, and `nix`

There are three levels, from safest to rawest:

| Level | Example | Use when |
|---|---|---|
| `std` | `std::fs::File::open`, `TcpStream::connect` | portable code; this covers most programs |
| `nix` / `rustix` | `nix::unistd::fork`, `nix::sched::unshare` | Unix-specific calls with Rust types and `Result` errors |
| `libc` | `libc::prctl(...)` | anything else; the raw C declarations, all `unsafe` |

```rust
use std::io::Write;

fn main() -> std::io::Result<()> {
    // Level 1: std. One write(2) per call here (stdout is a pipe or tty).
    std::io::stdout().write_all(b"std: hello\n")?;

    // Level 2: nix. A safe wrapper returning Result<_, Errno>.
    let pid = nix::unistd::getpid();
    let uname = nix::sys::utsname::uname().map_err(std::io::Error::from)?;
    println!("nix: pid={pid} kernel={:?}", uname.release());

    // Level 3: libc. The exact C call, with C error conventions.
    let msg = b"libc: hello\n";
    // SAFETY: fd 1 is open for the life of the process; the buffer is valid for msg.len() bytes.
    let n = unsafe { libc::write(1, msg.as_ptr().cast(), msg.len()) };
    if n < 0 { return Err(std::io::Error::last_os_error()); } // reads errno, like perror()

    // Resource limits, the same numbers as `ulimit -n` (OS Ch 43):
    let mut lim = libc::rlimit { rlim_cur: 0, rlim_max: 0 };
    // SAFETY: `lim` is a valid, writable rlimit struct.
    if unsafe { libc::getrlimit(libc::RLIMIT_NOFILE, &mut lim) } == 0 {
        println!("open files: soft={} hard={}", lim.rlim_cur, lim.rlim_max);
    }
    Ok(())
}
```

```bash
cargo build --release --bin rawsys
strace -f ./target/release/rawsys 2>&1 | tail -20      # Linux
# Look for: write(1, "std: hello\n", 11), uname(...), write(1, "libc: hello\n", 12), prlimit64(...)
```

**What a Rust `main` does before your code.** On Unix, the standard
library's startup code checks that file descriptors 0, 1, and 2 are open (it
opens `/dev/null` for any that are closed) and sets **`SIGPIPE` to
`SIG_IGN`**. The second one surprises people. A C program writing to a
closed pipe is killed by `SIGPIPE`. A Rust program gets an `EPIPE` error
instead, and `println!` panics on that error. This is the source of the
familiar `cargo run | head` panic: "failed printing to stdout: Broken pipe".
CLI tools that are meant to be piped either write with `writeln!` and handle
`ErrorKind::BrokenPipe` by exiting quietly, or restore the default handler
with `libc::signal(libc::SIGPIPE, libc::SIG_DFL)` at startup.

**Across the series:** what crossing into the kernel costs is in
[OS Ch 2](../os-linux/real-life-os-guide.md#chapter-2-kernel-space-vs-user-space-and-the-system-call).
How the Go runtime adds its own threads and syscalls on top is in
[OS Ch 73](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime).
A Rust binary has none of that runtime layer, which is one reason its
`strace` output is so short.

---

## 40. Processes, exit codes, and signals

### 40.1 Exit codes

| How the process ends | Exit code | Destructors run? |
|---|---|---|
| `main` returns `()` / `Ok(())` | 0 | yes |
| `main` returns `Err(e)` | 1, and `Error: {e:?}` is printed to stderr | yes |
| `std::process::exit(n)` | `n` | **no**: buffered writers are not flushed |
| panic (unwinding) | 101 | yes, for the panicking thread |
| `panic = "abort"`, or `std::process::abort()` | killed by `SIGABRT` (shell shows 134) | no |
| killed by signal N | shell shows `128 + N` | no |

```rust
use std::process::ExitCode;

fn run() -> Result<(), String> {
    let arg = std::env::args().nth(1).ok_or("usage: tool FILE")?;
    std::fs::metadata(&arg).map_err(|e| format!("{arg}: {e}"))?;
    Ok(())
}

fn main() -> ExitCode {
    // Returning ExitCode lets destructors run, unlike process::exit().
    match run() {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            eprintln!("error: {e}");
            ExitCode::from(2) // match the convention of your tool family (grep: 1 = no match, 2 = error)
        }
    }
}
```

[C §24.3](../c-lang/real-life-c-guide.md#24-building-production-clis-argument-parsing-config-precedence-exit-codes)
makes the same point: exit codes are part of your program's API, because
shell scripts and Kubernetes probes act on them.

### 40.2 Graceful SIGTERM handling

When Kubernetes stops a pod, or systemd stops a service, the process gets
`SIGTERM`, then `SIGKILL` after a grace period. A signal handler may only
call async-signal-safe functions (see
[C §55](../c-lang/real-life-c-guide.md#55-signals-handling-asynchronous-events-safely)),
so the safe pattern is to **set a flag in the handler and do the work in
normal code**. `signal-hook` implements exactly that:

```rust
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::time::Duration;

fn main() -> std::io::Result<()> {
    let term = Arc::new(AtomicBool::new(false));
    // The handler only stores `true`, which is async-signal-safe.
    signal_hook::flag::register(signal_hook::consts::SIGTERM, Arc::clone(&term))?;
    signal_hook::flag::register(signal_hook::consts::SIGINT, Arc::clone(&term))?;

    println!("pid {} working; send SIGTERM to stop", std::process::id());
    let mut processed = 0u64;
    while !term.load(Ordering::Relaxed) {
        processed += 1;                               // one unit of work
        std::thread::sleep(Duration::from_millis(100));
        if processed == 20 { break; }                 // demo: stop by itself after 2s
    }
    println!("shutting down cleanly after {processed} units: flushing, closing connections");
    Ok(())
}
```

```bash
cargo run --bin graceful & sleep 0.5; kill -TERM $!; wait $!; echo "exit=$?"   # exit=0
```

In Tokio, use `tokio::signal::unix::signal(SignalKind::terminate())` and
`.recv().await` inside the `select!` from §31.4.

### 40.3 PID 1: reaping zombies

In a container, your program may run as PID 1. PID 1 has two special
duties ([OS Ch 74](../os-linux/real-life-os-guide.md#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1)):
the kernel does not apply default signal actions to it, so a `SIGTERM` with
no handler is **ignored**, and it inherits orphaned processes, which become
zombies unless it calls `wait`. A minimal supervisor:

```rust
// linux-only
use nix::sys::signal::{kill, Signal};
use nix::sys::wait::{waitpid, WaitPidFlag, WaitStatus};
use nix::unistd::Pid;
use std::process::Command;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::time::Duration;

fn main() -> std::io::Result<()> {
    // As PID 1 we inherit orphans automatically. When we are NOT PID 1 (in a test, or under
    // a shell), ask the kernel to re-parent our orphaned descendants to us anyway.
    // SAFETY: PR_SET_CHILD_SUBREAPER takes an integer flag and no pointers.
    unsafe { libc::prctl(libc::PR_SET_CHILD_SUBREAPER, 1) };

    let term = Arc::new(AtomicBool::new(false));
    signal_hook::flag::register(signal_hook::consts::SIGTERM, Arc::clone(&term))?;

    // The shell exits after 1s and leaves two `sleep 2` processes orphaned.
    let child = Command::new("sh").args(["-c", "sleep 2 & sleep 2 & sleep 1; echo child done"]).spawn()?;
    let main_pid = Pid::from_raw(child.id() as i32);
    let mut main_exit: Option<i32> = None;
    let mut forwarded = false;

    loop {
        if term.load(Ordering::Relaxed) && !forwarded {
            let _ = kill(main_pid, Signal::SIGTERM);  // pass the shutdown request on to the app
            forwarded = true;
        }
        // Reap ANY exited child, including orphans re-parented to us.
        match waitpid(Pid::from_raw(-1), Some(WaitPidFlag::WNOHANG)) {
            Ok(WaitStatus::Exited(pid, code)) => {
                println!("reaped {pid} exit={code}{}", if pid == main_pid { " (main)" } else { " (orphan)" });
                if pid == main_pid { main_exit = Some(code); }
            }
            Ok(WaitStatus::Signaled(pid, sig, _)) => {
                println!("reaped {pid} killed by {sig:?}");
                if pid == main_pid { main_exit = Some(128 + sig as i32); }
            }
            Ok(_) => std::thread::sleep(Duration::from_millis(50)), // StillAlive: nothing to reap yet
            Err(nix::errno::Errno::ECHILD) => std::process::exit(main_exit.unwrap_or(0)), // all reaped
            Err(e) => return Err(e.into()),
        }
    }
}
```

```text
child done
reaped 602 exit=0 (main)
reaped 604 exit=0 (orphan)
reaped 603 exit=0 (orphan)
```

Without the `waitpid(-1, ..)` loop, the two orphaned `sleep` processes
would stay in the process table as zombies (`Z` in `ps`) until the
supervisor exits. In a long-running container that runs shell scripts or
health checks, zombies accumulate until the PID limit is reached.

In production use `tini` or `dumb-init` as PID 1, or Docker's `--init`.
Writing this once shows why they exist.

---

## 41. Files that survive crashes: `fsync` and atomic replacement

`std::fs::write(path, data)` truncates the file and then writes it. If the
machine crashes between the two steps, the file is empty or partly written.
Data you have written sits in the page cache until the kernel flushes it
([OS Ch 14](../os-linux/real-life-os-guide.md#chapter-14-the-i-o-stack-from-read-to-the-disk-platter)).
The crash-safe sequence is the one from
[OS Ch 77](../os-linux/real-life-os-guide.md#chapter-77-files-that-survive-crashes-page-cache-fsync-and-atomic-replacement):
write a temporary file in the same directory, `fsync` it, `rename` it over
the target, then `fsync` the directory.

```rust
use std::fs::{self, File};
use std::io::{self, Write};
use std::path::Path;

/// Replace `path` with `data` atomically: readers see the old file or the new one, never a mix.
pub fn atomic_write(path: &Path, data: &[u8]) -> io::Result<()> {
    let dir = path.parent().filter(|p| !p.as_os_str().is_empty()).unwrap_or(Path::new("."));
    let file_name = path.file_name().ok_or_else(|| io::Error::other("path has no file name"))?;
    // Same directory = same filesystem, so rename(2) is atomic.
    let tmp = dir.join(format!(".{}.tmp.{}", file_name.to_string_lossy(), std::process::id()));

    let result = (|| {
        let mut f = File::create(&tmp)?;
        f.write_all(data)?;
        f.sync_all()?;                 // fsync: the DATA is on stable storage
        drop(f);
        fs::rename(&tmp, path)?;       // atomic swap of the directory entry
        File::open(dir)?.sync_all()    // fsync the DIRECTORY: the rename itself is durable
    })();

    if result.is_err() {
        let _ = fs::remove_file(&tmp); // do not leave temporary files behind on failure
    }
    result
}

fn main() -> io::Result<()> {
    let path = std::env::temp_dir().join("rust-guide-state.json");
    atomic_write(&path, br#"{"version": 1}"#)?;
    atomic_write(&path, br#"{"version": 2}"#)?;
    println!("{}", fs::read_to_string(&path)?);
    fs::remove_file(path)
}
```

On macOS, `fsync` does not force the drive's own write cache to flush. Full
durability there needs `fcntl(F_FULLFSYNC)`. Databases such as SQLite issue
it, and Rust's `File::sync_all` uses it on Apple platforms.

The Go version of this lab is in
[OS Ch 77](../os-linux/real-life-os-guide.md#chapter-77-files-that-survive-crashes-page-cache-fsync-and-atomic-replacement).
The Rust version adds one thing: the closure plus `is_err` cleanup makes
sure the temporary file is removed on every failure path. Ch 63's key-value
store uses this function for snapshots.

---

## 42. Reading `/proc`: build your own `ps`

`/proc` exposes kernel process state as text files
([OS Ch 23](../os-linux/real-life-os-guide.md#chapter-23-everything-is-a-file-almost)).
Parsing it correctly takes care: the command name in `/proc/PID/stat` is
wrapped in parentheses and **can itself contain spaces and `)`**. A process
named `a) b` breaks naive parsers, and attackers have used that to make
monitoring tools misread process data.

```rust
// linux-only
use std::fs;

#[derive(Debug)]
struct Proc { pid: u32, ppid: u32, state: char, comm: String, rss_kib: u64, threads: u32, cmdline: String }

fn read_proc(pid: u32) -> Option<Proc> {
    let stat = fs::read_to_string(format!("/proc/{pid}/stat")).ok()?;
    // Format: "PID (COMM) STATE PPID ...". Split at the LAST ')' because COMM may contain ')'.
    let open = stat.find('(')?;
    let close = stat.rfind(')')?;
    let comm = stat[open + 1..close].to_string();
    let fields: Vec<&str> = stat[close + 2..].split(' ').collect(); // fields[0] = state
    let state = fields.first()?.chars().next()?;
    let ppid = fields.get(1)?.parse().ok()?;
    let threads = fields.get(17)?.parse().ok()?;                   // field 20 in `man 5 proc`

    let status = fs::read_to_string(format!("/proc/{pid}/status")).ok()?;
    let rss_kib = status.lines()
        .find_map(|l| l.strip_prefix("VmRSS:"))
        .and_then(|v| v.split_whitespace().next()?.parse().ok())
        .unwrap_or(0);                                             // kernel threads have no VmRSS

    let raw = fs::read(format!("/proc/{pid}/cmdline")).unwrap_or_default();
    let cmdline = raw.split(|&b| b == 0).filter(|a| !a.is_empty())
        .map(|a| String::from_utf8_lossy(a).replace('\n', " ")).collect::<Vec<_>>().join(" ");

    Some(Proc { pid, ppid, state, comm, rss_kib, threads, cmdline })
}

fn main() {
    let mut procs: Vec<Proc> = fs::read_dir("/proc").expect("Linux only")
        .filter_map(|e| e.ok()?.file_name().to_str()?.parse().ok())
        .filter_map(read_proc)                 // processes may exit while we scan: skip them
        .collect();
    procs.sort_by(|a, b| b.rss_kib.cmp(&a.rss_kib));

    println!("{:>7} {:>7} {} {:>9} {:>4}  COMMAND", "PID", "PPID", "S", "RSS(KiB)", "THR");
    for p in procs.iter().take(15) {
        let cmd = if p.cmdline.is_empty() { format!("[{}]", p.comm) } else { p.cmdline.clone() };
        println!("{:>7} {:>7} {} {:>9} {:>4}  {}", p.pid, p.ppid, p.state, p.rss_kib, p.threads,
            cmd.chars().take(60).collect::<String>());
    }
}
```

Compare this with the Go lab in
[OS Ch 79](../os-linux/real-life-os-guide.md#chapter-79-reading-proc-from-go-build-your-own-ps).
The `filter_map` chain handles the race where a process exits between
`read_dir` and `read_to_string`. Any of the three reads can fail, and
skipping that process is the correct response.

---

## 43. Memory: `mmap` and a counting global allocator

### 43.1 Memory-mapped files

```rust
use memmap2::Mmap;
use std::fs::File;
use std::io::Write;

fn main() -> std::io::Result<()> {
    let path = std::env::temp_dir().join("rust-guide-mmap.log");
    let mut f = File::create(&path)?;
    for i in 0..100_000 { writeln!(f, "event {i} status={}", if i % 97 == 0 { 500 } else { 200 })?; }
    drop(f);

    let file = File::open(&path)?;
    // SAFETY: the mapping is read-only, and no other process truncates this file while it is mapped.
    // If one did, reading the vanished pages would raise SIGBUS. That is why Mmap::map is `unsafe`.
    let map = unsafe { Mmap::map(&file)? };
    let bytes: &[u8] = &map;                       // the file, as a byte slice, with no read() copies

    let lines = bytes.iter().filter(|&&b| b == b'\n').count();
    let errors = bytes.split(|&b| b == b'\n').filter(|l| l.ends_with(b"500")).count();
    println!("{lines} lines, {errors} errors, {} bytes mapped", map.len());
    std::fs::remove_file(path)
}
```

Pages are loaded lazily through page faults
([OS Ch 11](../os-linux/real-life-os-guide.md#chapter-11-paging-page-faults-and-the-tlb)),
so mapping a 10 GB log costs nothing until you touch it. Note the `SAFETY`
comment: Rust's `&[u8]` promises the bytes will not change, but the kernel
lets another process change or truncate the underlying file. That is the
reason `Mmap::map` is `unsafe`. [C §57](../c-lang/real-life-c-guide.md#57-memory-mapped-files-mmap-for-real-files-msync-and-shared-mappings)
covers `msync` and shared writable mappings.

### 43.2 A counting global allocator

Every `Box`, `Vec`, and `String` allocates through one global allocator,
which you can replace. Wrapping the system allocator to count allocations is
a cheap, always-on way to find allocation-heavy code paths:

```rust
use std::alloc::{GlobalAlloc, Layout, System};
use std::sync::atomic::{AtomicUsize, Ordering::Relaxed};

struct Counting;
static ALLOCS: AtomicUsize = AtomicUsize::new(0);
static BYTES: AtomicUsize = AtomicUsize::new(0);

// SAFETY: we forward to System, which upholds GlobalAlloc's contract; we only add counters.
unsafe impl GlobalAlloc for Counting {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        ALLOCS.fetch_add(1, Relaxed);
        BYTES.fetch_add(layout.size(), Relaxed);
        unsafe { System.alloc(layout) }
    }
    unsafe fn dealloc(&self, ptr: *mut u8, layout: Layout) {
        unsafe { System.dealloc(ptr, layout) }
    }
}

#[global_allocator]
static GLOBAL: Counting = Counting;

fn measure<T>(label: &str, f: impl FnOnce() -> T) -> T {
    let (a0, b0) = (ALLOCS.load(Relaxed), BYTES.load(Relaxed));
    let out = f();
    println!("{label:<28} {:>6} allocs {:>9} bytes", ALLOCS.load(Relaxed) - a0, BYTES.load(Relaxed) - b0);
    out
}

fn main() {
    let words: Vec<String> = (0..10_000).map(|i| format!("w{i}")).collect();
    measure("join via format! in a loop", || {
        let mut s = String::new();
        for w in &words { s = format!("{s}{w},"); } // allocates a new String every iteration
        s.len()
    });
    measure("push_str with capacity", || {
        let mut s = String::with_capacity(words.iter().map(|w| w.len() + 1).sum());
        for w in &words { s.push_str(w); s.push(','); }
        s.len()
    });
}
```

```text
join via format! in a loop    19998 allocs 868126828 bytes
push_str with capacity            1 allocs     58890 bytes
```

The first version copies the whole string on every iteration. It allocated
868 MB in total, with quadratic growth, to build a 58 KB result. This is
the same bug as `s += x` in a Python or Java loop, and the counting
allocator makes it obvious.

The same `#[global_allocator]` line is how you switch to `mimalloc` or
`jemalloc` in production (`static GLOBAL: mimalloc::MiMalloc = mimalloc::MiMalloc;`).
The [C guide's Ch 39](../c-lang/real-life-c-guide.md#39-writing-a-real-allocator-free-lists-coalescing-and-dlmalloc-s-ideas)
explains what such an allocator does internally.

---

## 44. A container in 100 lines: namespaces with `nix`

A container is a process with its own namespaces (what it can see), cgroups
(what it can use), and root filesystem
([OS Ch 44](../os-linux/real-life-os-guide.md#chapter-44-what-a-container-actually-is-namespaces-cgroups-a-filesystem)).
The Go version of this lab is
[OS Ch 80](../os-linux/real-life-os-guide.md#chapter-80-a-container-runtime-in-150-lines-of-go).
Run this one as root inside the `--privileged` lab container.

```rust
// linux-only
// usage: minibox HOSTNAME CMD [ARGS...]     e.g. minibox box1 sh
use nix::mount::{mount, MsFlags};
use nix::sched::{unshare, CloneFlags};
use nix::sys::wait::{waitpid, WaitStatus};
use nix::unistd::{execvp, fork, pipe, sethostname, ForkResult};
use std::ffi::CString;
use std::fs::File;
use std::io::{Read, Write};

fn setup_cgroup(pid: i32) -> std::io::Result<()> {
    // cgroup v2: a child group with a 64 MiB memory cap and a 20-process cap.
    let cg = "/sys/fs/cgroup/minibox";
    std::fs::create_dir_all(cg)?;
    for (file, value) in [("memory.max", "67108864".to_string()), ("pids.max", "20".into()), ("cgroup.procs", pid.to_string())] {
        // memory.max only exists if the parent delegates the memory controller (cgroup.subtree_control).
        std::fs::write(format!("{cg}/{file}"), value)
            .map_err(|e| std::io::Error::new(e.kind(), format!("{cg}/{file}: {e}")))?;
    }
    Ok(())
}

fn child(hostname: &str, argv: &[CString]) -> nix::Result<()> {
    sethostname(hostname)?;                                   // visible only inside the new UTS namespace
    // Make every mount private so our /proc mount does not propagate back to the host.
    mount(None::<&str>, "/", None::<&str>, MsFlags::MS_REC | MsFlags::MS_PRIVATE, None::<&str>)?;
    // A fresh /proc that shows only this PID namespace: `ps` inside sees PID 1 = our command.
    mount(Some("proc"), "/proc", Some("proc"), MsFlags::empty(), None::<&str>)?;
    execvp(&argv[0], argv)?;                                  // replaces this process; returns only on error
    unreachable!()
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    if args.len() < 2 { eprintln!("usage: minibox HOSTNAME CMD [ARGS...]"); std::process::exit(2); }
    let argv: Vec<CString> = args[1..].iter().map(|a| CString::new(a.as_str()).expect("no NUL")).collect();

    // New UTS (hostname), PID, and mount namespaces. CLONE_NEWPID affects our CHILDREN:
    // the next process we fork becomes PID 1 of the new namespace.
    unshare(CloneFlags::CLONE_NEWUTS | CloneFlags::CLONE_NEWPID | CloneFlags::CLONE_NEWNS)
        .expect("unshare (are you root, in a --privileged container?)");

    // The child must not run the command until the parent has put it in the cgroup.
    // Otherwise the command can start (and fork its own children) before the limits apply.
    let (ready_r, ready_w) = pipe().expect("pipe");

    // SAFETY: the program is single-threaded at this point, so fork() is sound.
    match unsafe { fork() }.expect("fork") {
        ForkResult::Child => {
            drop(ready_w);
            let mut go = [0u8; 1];
            let _ = File::from(ready_r).read(&mut go);   // blocks until the parent writes (or exits: EOF)
            if let Err(e) = child(&args[0], &argv) { eprintln!("child setup failed: {e}"); std::process::exit(127); }
        }
        ForkResult::Parent { child } => {
            drop(ready_r);
            if let Err(e) = setup_cgroup(child.as_raw()) { eprintln!("cgroup limits skipped: {e}"); }
            let _ = File::from(ready_w).write_all(b"1"); // release the child
            match waitpid(child, None) {
                Ok(WaitStatus::Exited(_, code)) => std::process::exit(code),
                Ok(WaitStatus::Signaled(_, sig, _)) => std::process::exit(128 + sig as i32),
                other => { eprintln!("unexpected wait result: {other:?}"); std::process::exit(1) }
            }
        }
    }
}
```

```bash
# --cgroupns=host so the memory and pids controllers are delegated to new groups
docker run --rm -it --privileged --cgroupns=host -v "$PWD":/work -w /work rust:1 bash
cargo build --release --bin minibox
./target/release/minibox box1 sh -c 'hostname; echo "pid=$$"; ls /proc | grep -E "^[0-9]+$" | tr "\n" " "; cat /proc/self/cgroup'
# box1
# pid=1                    <- PID 1 in its own namespace
# 1 3 4 5                  <- only this namespace's processes are visible in the new /proc
# 0::/minibox              <- running inside the limited cgroup (because of the sync pipe)
hostname                   # the host's hostname is unchanged
cat /sys/fs/cgroup/minibox/memory.max   # 67108864
```

Without `--cgroupns=host`, Docker gives the container its own cgroup
namespace in which the memory controller is not delegated to child groups.
`memory.max` then does not exist, and the program prints
`cgroup limits skipped: /sys/fs/cgroup/minibox/memory.max: Permission denied`.
Real runtimes handle this by moving processes into a leaf group and
enabling controllers in `cgroup.subtree_control`. The "no internal
processes" rule of cgroup v2 requires this.

**The pipe is not optional.** The first version of this lab, written
without it, printed the Docker cgroup instead of `/minibox`. The parent
wrote the child's PID into `cgroup.procs` *after* `fork`, and by then the
child had already started `sh`. A process that starts before its limits are
applied can fork children outside them. `runc` solves this the same way,
with a sync pipe between parent and child.

Why `fork` is `unsafe`: after `fork` in a multi-threaded program, the child
contains only the calling thread. Any lock another thread held at that
moment stays locked forever in the child, including the allocator's lock.
`Command` uses `posix_spawn` or a careful `fork`+`exec` for that reason.
Call `fork` yourself only in a single-threaded program, as here.

What this toy omits, and real runtimes (`runc`, `youki`, which is written
in Rust) add: `pivot_root` into an image root filesystem, network
namespaces with veth pairs
([OS Ch 69](../os-linux/real-life-os-guide.md#chapter-69-advanced-linux-networking-namespaces-routing-nftables-and-packet-paths)),
user namespaces, dropping capabilities, and seccomp filters
([Security in Depth Ch 15](../security/real-life-security-guide-v1.md#chapter-15-container-internals-and-isolation)).

---

# Part V — Network and HTTP services: the networking and HTTPS guides in Rust

This Part follows one request down the stack and back up, in the order of the
[HTTPS lifecycle guide](../v2-https/real-life-guide-v1.md): a DNS lookup,
HTTP over a raw socket, a production server, the outbound clients that
server calls, TLS, and finally a reverse proxy in front of all of it. Each
chapter names the guide chapter it implements. Read that chapter for the
protocol, and this one for the code.

---

## 45. A DNS client from raw bytes

**Implements:** [Net Ch 18](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses)
and [HTTPS Ch 2](../v2-https/real-life-guide-v1.md#chapter-2-dns-finding-the-address).
The networking guide's first Go lab builds the same tool.

A DNS message is a 12-byte header, then questions, then resource records.
Names are length-prefixed labels (`3www7example3com0`), and in answers they
are usually **compressed**: a byte with the top two bits set (`0xC0`) is a
pointer to an earlier offset in the message. That pointer is what makes
parsing DNS dangerous. A malicious response can point a name at itself, and
a naive parser loops forever. Several C resolvers have had CVEs from exactly
this.

```rust
// src/bin/dnsq.rs   usage: dnsq NAME [A|AAAA|CNAME] [SERVER]
use std::net::{Ipv4Addr, Ipv6Addr, UdpSocket};
use std::time::Duration;

const TYPE_A: u16 = 1;
const TYPE_CNAME: u16 = 5;
const TYPE_AAAA: u16 = 28;

fn build_query(id: u16, name: &str, qtype: u16) -> Result<Vec<u8>, String> {
    let mut q = Vec::with_capacity(512);
    q.extend_from_slice(&id.to_be_bytes());
    q.extend_from_slice(&0x0100u16.to_be_bytes());  // flags: RD (recursion desired)
    q.extend_from_slice(&1u16.to_be_bytes());       // QDCOUNT = 1
    q.extend_from_slice(&[0u8; 6]);                 // ANCOUNT, NSCOUNT, ARCOUNT = 0
    for label in name.trim_end_matches('.').split('.') {
        let len = u8::try_from(label.len()).ok().filter(|&l| (1..=63).contains(&l))
            .ok_or_else(|| format!("bad label {label:?}"))?;
        q.push(len);
        q.extend_from_slice(label.as_bytes());
    }
    q.push(0);                                      // root label ends the name
    q.extend_from_slice(&qtype.to_be_bytes());
    q.extend_from_slice(&1u16.to_be_bytes());       // class IN
    Ok(q)
}

/// A bounds-checked cursor. Every read returns Err instead of panicking on short input.
struct Reader<'a> { buf: &'a [u8], pos: usize }

impl<'a> Reader<'a> {
    fn take(&mut self, n: usize) -> Result<&'a [u8], String> {
        let end = self.pos.checked_add(n).filter(|&e| e <= self.buf.len()).ok_or("truncated message")?;
        let s = &self.buf[self.pos..end];
        self.pos = end;
        Ok(s)
    }
    fn u16(&mut self) -> Result<u16, String> { Ok(u16::from_be_bytes(self.take(2)?.try_into().unwrap())) }
    fn u32(&mut self) -> Result<u32, String> { Ok(u32::from_be_bytes(self.take(4)?.try_into().unwrap())) }

    /// Reads a possibly-compressed name. Defends against pointer loops and oversized names.
    fn name(&mut self) -> Result<String, String> {
        let mut labels = Vec::new();
        let mut pos = self.pos;
        let mut jumped = false;
        let mut jumps = 0;
        let mut total = 0;
        loop {
            let len = *self.buf.get(pos).ok_or("name runs past end")? as usize;
            if len & 0xC0 == 0xC0 {                                   // compression pointer
                let lo = *self.buf.get(pos + 1).ok_or("truncated pointer")? as usize;
                if !jumped { self.pos = pos + 2; }                      // the cursor resumes after the FIRST pointer
                jumped = true;
                jumps += 1;
                if jumps > 16 { return Err("compression pointer loop".into()); }
                pos = ((len & 0x3F) << 8) | lo;
                continue;
            }
            if len == 0 {
                if !jumped { self.pos = pos + 1; }
                break;
            }
            let label = self.buf.get(pos + 1..pos + 1 + len).ok_or("label runs past end")?;
            total += len + 1;
            if total > 255 { return Err("name longer than 255 bytes".into()); }
            labels.push(String::from_utf8_lossy(label).into_owned());
            pos += 1 + len;
        }
        Ok(labels.join("."))
    }
}

#[derive(Debug, PartialEq)]
enum Answer { A(Ipv4Addr), Aaaa(Ipv6Addr), Cname(String), Other(u16) }

fn parse_response(msg: &[u8], want_id: u16) -> Result<Vec<(String, u32, Answer)>, String> {
    let mut r = Reader { buf: msg, pos: 0 };
    let id = r.u16()?;
    if id != want_id { return Err(format!("ID mismatch: {id} != {want_id} (spoofed or stale reply)")); }
    let flags = r.u16()?;
    if flags & 0x8000 == 0 { return Err("not a response".into()); }
    if flags & 0x0200 != 0 { return Err("truncated (TC): retry over TCP".into()); }
    match flags & 0x000F {
        0 => {}
        3 => return Err("NXDOMAIN: name does not exist".into()),
        rcode => return Err(format!("server error, RCODE={rcode}")),
    }
    let (qd, an) = (r.u16()?, r.u16()?);
    r.take(4)?;                                       // skip NSCOUNT, ARCOUNT
    for _ in 0..qd { r.name()?; r.take(4)?; }          // skip the echoed questions

    let mut out = Vec::new();
    for _ in 0..an {
        let name = r.name()?;
        let (rtype, _class, ttl) = (r.u16()?, r.u16()?, r.u32()?);
        let rdlen = r.u16()? as usize;
        let rdata_start = r.pos;
        let rdata = r.take(rdlen)?;
        let ans = match (rtype, rdlen) {
            (TYPE_A, 4) => Answer::A(Ipv4Addr::from(<[u8; 4]>::try_from(rdata).unwrap())),
            (TYPE_AAAA, 16) => Answer::Aaaa(Ipv6Addr::from(<[u8; 16]>::try_from(rdata).unwrap())),
            (TYPE_CNAME, _) => Answer::Cname(Reader { buf: msg, pos: rdata_start }.name()?), // may point backwards
            (t, _) => Answer::Other(t),
        };
        out.push((name, ttl, ans));
    }
    Ok(out)
}

fn main() -> Result<(), String> {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let name = args.first().map(String::as_str).unwrap_or("example.com");
    let qtype = match args.get(1).map(String::as_str) { Some("AAAA") => TYPE_AAAA, Some("CNAME") => TYPE_CNAME, _ => TYPE_A };
    let server = args.get(2).map(String::as_str).unwrap_or("1.1.1.1:53");

    let id: u16 = rand::random();                     // random ID + random source port: anti-spoofing
    let query = build_query(id, name, qtype)?;
    let sock = UdpSocket::bind("0.0.0.0:0").map_err(|e| e.to_string())?;
    sock.set_read_timeout(Some(Duration::from_secs(3))).map_err(|e| e.to_string())?;
    sock.send_to(&query, server).map_err(|e| e.to_string())?;

    let mut buf = [0u8; 1232];                        // the EDNS-safe UDP size; classic DNS max is 512
    let (n, from) = sock.recv_from(&mut buf).map_err(|e| format!("no reply: {e}"))?;
    println!(";; {n} bytes from {from}, id {id}");
    for (owner, ttl, ans) in parse_response(&buf[..n], id)? {
        println!("{owner:<30} {ttl:>6}  {ans:?}");
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn query_layout() {
        let q = build_query(0xBEEF, "a.bc", TYPE_A).unwrap();
        assert_eq!(&q[..2], &[0xBE, 0xEF]);
        assert_eq!(&q[12..], &[1, b'a', 2, b'b', b'c', 0, 0, 1, 0, 1]);
        assert!(build_query(1, &"x".repeat(64), TYPE_A).is_err());
    }

    #[test]
    fn rejects_pointer_loop() {
        // header (id=7, response, 1 answer) + an answer whose name points at itself (offset 12)
        let mut m = vec![0, 7, 0x81, 0x80, 0, 0, 0, 1, 0, 0, 0, 0];
        m.extend_from_slice(&[0xC0, 12]);
        assert_eq!(parse_response(&m, 7).unwrap_err(), "compression pointer loop");
    }

    #[test]
    fn rejects_truncated_garbage() {
        for len in 0..12 { assert!(parse_response(&[0xFF; 12][..len], 0xFFFF).is_err()); }
    }
}
```

```bash
cargo run --bin dnsq -- example.com
cargo run --bin dnsq -- www.github.com A 8.8.8.8:53   # a CNAME chain
cargo test --bin dnsq
```

```text
;; 61 bytes from 1.1.1.1:53, id 61600
example.com                       112  A(104.20.23.154)
example.com                       112  A(172.66.147.243)
;; 62 bytes from 1.1.1.1:53, id 23339
www.github.com                   3581  Cname("github.com")
github.com                         41  A(20.207.73.82)
```

Your addresses, TTLs, and IDs will differ. Servers usually compress a
CNAME target such as `github.com` into a pointer to the `github.com` suffix
of the question. That is why the parser decodes CNAME data with a fresh
`Reader` over the whole message, not just the record's bytes.

Look at what the safe parser buys you. Every read goes through `take`,
which checks bounds with `checked_add`. Malformed input becomes an `Err`,
never a panic or an out-of-bounds read. The two security checks are
explicit and tested: the pointer-loop limit and the 255-byte name limit.
The same parser in C needs every one of these checks done by hand, and
DNS-parsing CVEs show what happens when one is missed.

---

## 46. HTTP/1.1 from a bare `TcpListener`

**Implements:** [Net Ch 28](../networking/tcp-ip/real-life-guide-v1.md#chapter-28-http-how-the-web-actually-talks),
[HTTPS Ch 5](../v2-https/real-life-guide-v1.md#chapter-5-http-the-conversation),
and the [C guide's §21.2](../c-lang/real-life-c-guide.md#21-sockets-a-tcp-echo-server-then-a-tiny-http-server),
with the limits a real server needs.

You will use a framework (next chapter). Writing HTTP/1.1 by hand once
shows you what every framework does about **message framing**. Framing
mistakes are behind request smuggling, which Ch 50 covers.

```rust
// src/bin/tinyhttp.rs
use std::io::{BufRead, BufReader, Read, Write};
use std::net::{TcpListener, TcpStream};
use std::time::Duration;

const MAX_HEAD: u64 = 8 * 1024;      // request line + headers
const MAX_BODY: usize = 1024 * 1024;

struct Request { method: String, path: String, headers: Vec<(String, String)>, body: Vec<u8>, keep_alive: bool }

enum HttpError { BadRequest(&'static str), HeadTooLarge, BodyTooLarge, NotImplemented(&'static str), Closed }

impl Request {
    fn header(&self, name: &str) -> Option<&str> {
        self.headers.iter().find(|(k, _)| k.eq_ignore_ascii_case(name)).map(|(_, v)| v.as_str())
    }
}

fn read_request(reader: &mut BufReader<TcpStream>) -> Result<Request, HttpError> {
    let mut budget = MAX_HEAD;
    let mut read_line = |reader: &mut BufReader<TcpStream>| -> Result<String, HttpError> {
        let mut line = String::new();
        let n = reader.by_ref().take(budget).read_line(&mut line).map_err(|_| HttpError::Closed)?;
        if n == 0 { return Err(HttpError::Closed); }
        if !line.ends_with('\n') { return Err(HttpError::HeadTooLarge); } // budget ran out mid-line
        budget -= n as u64;
        Ok(line.trim_end_matches(['\r', '\n']).to_string())
    };

    let request_line = read_line(reader)?;
    let mut parts = request_line.split(' ');
    let (method, path, version) = match (parts.next(), parts.next(), parts.next(), parts.next()) {
        (Some(m), Some(p), Some(v), None) if p.starts_with('/') => (m.to_string(), p.to_string(), v),
        _ => return Err(HttpError::BadRequest("malformed request line")),
    };
    if version != "HTTP/1.1" && version != "HTTP/1.0" { return Err(HttpError::BadRequest("bad version")); }

    let mut headers = Vec::new();
    loop {
        let line = read_line(reader)?;
        if line.is_empty() { break; }                                  // blank line ends the head
        let (k, v) = line.split_once(':').ok_or(HttpError::BadRequest("header without ':'"))?;
        if k.is_empty() || k.ends_with([' ', '\t']) {                  // "Content-Length : 5" is a smuggling trick
            return Err(HttpError::BadRequest("whitespace before ':'"));
        }
        headers.push((k.to_string(), v.trim().to_string()));
    }
    let mut req = Request { method, path, headers, body: Vec::new(), keep_alive: version == "HTTP/1.1" };
    if let Some(c) = req.header("connection") {
        req.keep_alive = !c.eq_ignore_ascii_case("close") && (version == "HTTP/1.1" || c.eq_ignore_ascii_case("keep-alive"));
    }

    // FRAMING: exactly one way to know where the body ends.
    let te = req.header("transfer-encoding").is_some();
    let cls: Vec<&str> = req.headers.iter().filter(|(k, _)| k.eq_ignore_ascii_case("content-length")).map(|(_, v)| v.as_str()).collect();
    if te && !cls.is_empty() { return Err(HttpError::BadRequest("both Transfer-Encoding and Content-Length")); }
    if te { return Err(HttpError::NotImplemented("chunked bodies")); }
    if cls.len() > 1 { return Err(HttpError::BadRequest("multiple Content-Length headers")); }
    if let Some(cl) = cls.first() {
        if cl.is_empty() || !cl.bytes().all(|b| b.is_ascii_digit()) { return Err(HttpError::BadRequest("bad Content-Length")); }
        let len: usize = cl.parse().map_err(|_| HttpError::BodyTooLarge)?;
        if len > MAX_BODY { return Err(HttpError::BodyTooLarge); }
        req.body = vec![0; len];
        reader.read_exact(&mut req.body).map_err(|_| HttpError::Closed)?;
    }
    Ok(req)
}

fn respond(w: &mut TcpStream, status: &str, body: &[u8], keep_alive: bool) -> std::io::Result<()> {
    let conn = if keep_alive { "keep-alive" } else { "close" };
    write!(w, "HTTP/1.1 {status}\r\nContent-Length: {}\r\nContent-Type: text/plain\r\nConnection: {conn}\r\n\r\n", body.len())?;
    w.write_all(body)
}

/// Closing a socket that still has unread input makes the kernel send RST, and the RST can
/// destroy the response before the client reads it. Half-close, then drain briefly (nginx's "lingering close").
fn lingering_close(writer: &mut TcpStream, mut reader: BufReader<TcpStream>) -> std::io::Result<()> {
    writer.shutdown(std::net::Shutdown::Write)?;                     // send FIN after the response
    reader.get_ref().set_read_timeout(Some(Duration::from_secs(2)))?;
    let _ = std::io::copy(&mut reader.by_ref().take(64 * 1024), &mut std::io::sink()); // discard what is left
    Ok(())
}

fn handle(stream: TcpStream) -> std::io::Result<()> {
    stream.set_read_timeout(Some(Duration::from_secs(10)))?;          // Slowloris defense (HTTPS Ch 19)
    let mut writer = stream.try_clone()?;
    let mut reader = BufReader::new(stream);
    loop {                                                             // keep-alive: many requests per connection
        let req = match read_request(&mut reader) {
            Ok(r) => r,
            Err(HttpError::Closed) => return Ok(()),
            Err(e) => {
                let (status, msg) = match e {
                    HttpError::BadRequest(m) => ("400 Bad Request", m),
                    HttpError::HeadTooLarge => ("431 Request Header Fields Too Large", "head too large"),
                    HttpError::BodyTooLarge => ("413 Content Too Large", "body too large"),
                    HttpError::NotImplemented(m) => ("501 Not Implemented", m),
                    HttpError::Closed => unreachable!(),
                };
                // After a framing error we cannot know where the next request starts: always close.
                respond(&mut writer, status, msg.as_bytes(), false)?;
                return lingering_close(&mut writer, reader);
            }
        };
        let (status, body) = match (req.method.as_str(), req.path.as_str()) {
            ("GET", "/") => ("200 OK", b"hello from tinyhttp\n".to_vec()),
            ("POST", "/echo") => ("200 OK", req.body.clone()),
            (_, "/" | "/echo") => ("405 Method Not Allowed", b"method not allowed\n".to_vec()),
            _ => ("404 Not Found", b"not found\n".to_vec()),
        };
        respond(&mut writer, status, &body, req.keep_alive)?;
        if !req.keep_alive { return Ok(()); }
    }
}

fn main() -> std::io::Result<()> {
    let listener = TcpListener::bind("127.0.0.1:0")?;
    let addr = listener.local_addr()?;
    std::thread::spawn(move || {
        for s in listener.incoming().flatten() {
            std::thread::spawn(move || { let _ = handle(s); });
        }
    });

    // A raw client: send bytes exactly as written, print the status line of each response.
    let send = |raw: &str| -> std::io::Result<Vec<String>> {
        let mut s = TcpStream::connect(addr)?;
        s.write_all(raw.as_bytes())?;
        s.shutdown(std::net::Shutdown::Write)?;
        let mut out = String::new();
        s.read_to_string(&mut out)?;
        Ok(out.lines().filter(|l| l.starts_with("HTTP/1.1")).map(String::from).collect())
    };
    println!("pipelined:  {:?}", send("GET / HTTP/1.1\r\nHost: x\r\n\r\nPOST /echo HTTP/1.1\r\nHost: x\r\nContent-Length: 5\r\nConnection: close\r\n\r\nhello")?);
    println!("smuggling:  {:?}", send("POST /echo HTTP/1.1\r\nHost: x\r\nContent-Length: 4\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\n")?);
    println!("dup CL:     {:?}", send("POST /echo HTTP/1.1\r\nHost: x\r\nContent-Length: 5\r\nContent-Length: 6\r\n\r\nhello!")?);
    println!("space-colon:{:?}", send("POST /echo HTTP/1.1\r\nHost: x\r\nContent-Length : 5\r\n\r\nhello")?);
    println!("huge head:  {:?}", send(&format!("GET / HTTP/1.1\r\nX-Big: {}\r\n\r\n", "a".repeat(10_000)))?);
    Ok(())
}
```

```text
pipelined:  ["HTTP/1.1 200 OK", "HTTP/1.1 200 OK"]
smuggling:  ["HTTP/1.1 400 Bad Request"]
dup CL:     ["HTTP/1.1 400 Bad Request"]
space-colon:["HTTP/1.1 400 Bad Request"]
huge head:  ["HTTP/1.1 431 Request Header Fields Too Large"]
```

**Why `lingering_close` exists.** The first version of this server simply
returned after sending the 431. The "huge head" client then failed with
`Connection reset by peer` and never saw the status line. The server had
closed its socket with about 2 KB of the client's request still unread in
the receive buffer. In that situation the kernel sends RST instead of FIN,
and a RST tells the peer to throw away any data it has not read yet,
including our response
([Net Ch 23](../networking/tcp-ip/real-life-guide-v1.md#chapter-23-tcp-part-3-closing-a-connection-and-the-states)).
Any server that rejects a request before reading all of it must half-close
and drain, as above.

Each rejection above is a real request-smuggling technique
([Security in Depth Ch 30](../security/real-life-security-guide-v1.md#chapter-30-http-request-smuggling-and-desync)).
Smuggling works when two HTTP parsers, usually a proxy and a backend,
disagree about where one request ends. The defense is to be strict: one
framing method, no duplicates, no whitespace before the colon, and close the
connection after any framing error. Rust's types do not give you this.
These are protocol rules, and you have to implement them. What Rust does
guarantee is that none of these hostile inputs can overflow a buffer.

---

## 47. A production HTTP server with Axum: limits, timeouts, tracing, shutdown

**Implements:** [HTTPS Ch 19](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown)
(the Go version) and [Go §62](../Golang/real-life-golang-guide.md#62-production-http-apis-validation-timeouts-middleware-shutdown).

Axum is a thin layer over Hyper (the HTTP implementation) and Tower (a
middleware abstraction: `Service<Request> -> Future<Response>`). Its main
idea is **extractors**: a handler's parameter types say what it needs from
the request, and Axum parses and validates them before the handler runs.

```rust
// src/bin/api.rs
use axum::extract::{Path, State};
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::routing::{get, post};
use axum::{Json, Router};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Arc, RwLock};
use std::time::Duration;
use tower_http::{limit::RequestBodyLimitLayer, timeout::TimeoutLayer, trace::TraceLayer};

#[derive(Clone, Default)]
struct AppState {                                   // cloned per request: keep it cheap (Arcs inside)
    users: Arc<RwLock<HashMap<u64, User>>>,
    next_id: Arc<AtomicU64>,
}

#[derive(Clone, Serialize)]
struct User { id: u64, name: String, email: String }

#[derive(Deserialize)]
#[serde(deny_unknown_fields)]                       // reject {"name":..,"is_admin":true} (mass assignment)
struct NewUser { name: String, email: String }

/// One error type for the whole API, mapped to HTTP in exactly one place.
enum ApiError { NotFound, Invalid(String) }

impl IntoResponse for ApiError {
    fn into_response(self) -> Response {
        let (status, msg) = match self {
            ApiError::NotFound => (StatusCode::NOT_FOUND, "not found".to_string()),
            ApiError::Invalid(m) => (StatusCode::UNPROCESSABLE_ENTITY, m),
        };
        (status, Json(serde_json::json!({ "error": msg }))).into_response()
    }
}

async fn healthz() -> &'static str { "ok" }

async fn create_user(State(st): State<AppState>, Json(body): Json<NewUser>) -> Result<(StatusCode, Json<User>), ApiError> {
    let name = body.name.trim();
    if name.is_empty() || name.len() > 100 { return Err(ApiError::Invalid("name must be 1-100 characters".into())); }
    if !body.email.contains('@') || body.email.len() > 254 { return Err(ApiError::Invalid("invalid email".into())); }
    let id = st.next_id.fetch_add(1, Ordering::Relaxed) + 1;
    let user = User { id, name: name.into(), email: body.email.to_lowercase() };
    st.users.write().unwrap().insert(id, user.clone());   // lock held briefly, never across .await
    Ok((StatusCode::CREATED, Json(user)))
}

async fn get_user(State(st): State<AppState>, Path(id): Path<u64>) -> Result<Json<User>, ApiError> {
    st.users.read().unwrap().get(&id).cloned().map(Json).ok_or(ApiError::NotFound)
}

fn app(state: AppState) -> Router {
    Router::new()
        .route("/healthz", get(healthz))
        .route("/users", post(create_user))
        .route("/users/{id}", get(get_user))               // axum 0.8 path syntax: {id}
        .layer(TraceLayer::new_for_http())                 // one span per request, with latency
        .layer(TimeoutLayer::with_status_code(StatusCode::SERVICE_UNAVAILABLE, Duration::from_secs(10)))
        .layer(RequestBodyLimitLayer::new(64 * 1024))      // 413 for bodies over 64 KiB
        .with_state(state)
}

async fn shutdown_signal() {
    let ctrl_c = async { tokio::signal::ctrl_c().await.expect("install Ctrl-C handler") };
    #[cfg(unix)]
    let term = async {
        tokio::signal::unix::signal(tokio::signal::unix::SignalKind::terminate())
            .expect("install SIGTERM handler")
            .recv()
            .await;
    };
    #[cfg(not(unix))]
    let term = std::future::pending::<()>();
    tokio::select! { _ = ctrl_c => {}, _ = term => {} }
    tracing::info!("shutdown signal received, draining connections");
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    tracing_subscriber::fmt().with_target(false).init();
    let listener = tokio::net::TcpListener::bind("127.0.0.1:3000").await?;
    tracing::info!("listening on {}", listener.local_addr()?);
    axum::serve(listener, app(AppState::default()))
        .with_graceful_shutdown(shutdown_signal())        // stop accepting, finish in-flight requests
        .await
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::body::{to_bytes, Body};
    use axum::http::Request;
    use tower::ServiceExt;                                // .oneshot(): call the Router without a socket

    async fn call(app: &Router, method: &str, uri: &str, body: &str) -> (StatusCode, String) {
        let req = Request::builder().method(method).uri(uri)
            .header("content-type", "application/json").body(Body::from(body.to_string())).unwrap();
        let resp = app.clone().oneshot(req).await.unwrap();
        let status = resp.status();
        let bytes = to_bytes(resp.into_body(), usize::MAX).await.unwrap();
        (status, String::from_utf8_lossy(&bytes).into_owned())
    }

    #[tokio::test]
    async fn user_lifecycle_and_limits() {
        let app = app(AppState::default());
        let (s, b) = call(&app, "POST", "/users", r#"{"name":"Ada","email":"ADA@x.io"}"#).await;
        assert_eq!(s, StatusCode::CREATED);
        assert!(b.contains(r#""email":"ada@x.io""#));
        assert_eq!(call(&app, "GET", "/users/1", "").await.0, StatusCode::OK);
        assert_eq!(call(&app, "GET", "/users/99", "").await.0, StatusCode::NOT_FOUND);
        assert_eq!(call(&app, "GET", "/users/abc", "").await.0, StatusCode::BAD_REQUEST);  // Path<u64> rejected
        assert_eq!(call(&app, "POST", "/users", r#"{"name":"x","email":"x@y","is_admin":true}"#).await.0,
                   StatusCode::UNPROCESSABLE_ENTITY);                                    // deny_unknown_fields
        assert_eq!(call(&app, "POST", "/users", r#"{"name":"","email":"a@b"}"#).await.0, StatusCode::UNPROCESSABLE_ENTITY);
        let huge = format!(r#"{{"name":"{}","email":"a@b"}}"#, "x".repeat(100_000));
        assert_eq!(call(&app, "POST", "/users", &huge).await.0, StatusCode::PAYLOAD_TOO_LARGE);
    }
}
```

```bash
cargo run --bin api &
curl -s localhost:3000/users -H 'content-type: application/json' -d '{"name":"Ada","email":"ada@example.com"}'
# {"id":1,"name":"Ada","email":"ada@example.com"}
curl -s -o /dev/null -w '%{http_code}\n' localhost:3000/users/42      # 404
kill -TERM %1                                                          # graceful drain, exit 0
cargo test --bin api
```

**What each layer defends against**, mapped to the HTTPS guide:

| Control | Threat | Guide |
|---|---|---|
| `RequestBodyLimitLayer` | memory exhaustion from huge bodies | [HTTPS Ch 12](../v2-https/real-life-guide-v1.md#chapter-12-security-protecting-the-entire-lifecycle) |
| `TimeoutLayer` | slow handlers holding resources | [HTTPS Ch 19](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) |
| `deny_unknown_fields` + validation | mass assignment, bad input | [Sec Ch 45](../security/real-life-guide.md#chapter-45-broken-access-control) |
| typed `Path<u64>` | injection through path parameters | [Sec Ch 42](../security/real-life-guide.md#chapter-42-injection-when-data-becomes-code) |
| graceful shutdown | dropped in-flight requests during deploys | [OS Ch 82](../os-linux/real-life-os-guide.md#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes) |
| `TraceLayer` | being unable to tell what happened | [HTTPS Ch 17](../v2-https/real-life-guide-v1.md#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs) |

**One gap to know about:** `axum::serve` does not expose a **header read
timeout**. A Slowloris client that sends one header byte every few seconds
is limited only by the OS and by `TimeoutLayer`, and `TimeoutLayer` starts
counting only after the headers have been parsed. Production deployments
either sit behind a proxy that enforces it (nginx, Envoy, or a cloud load
balancer), or serve the Router with `hyper-util`'s connection builder,
which has `header_read_timeout`. Ch 50 shows that builder.

---

## 48. Outbound requests: timeouts, retries, idempotency, and an SSRF guard

**Implements:** [HTTPS Ch 21](../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers)
and [Ch 22](../v2-https/real-life-guide-v1.md#chapter-22-outbound-requests-ssrf-and-safe-http-clients);
[Security in Depth Ch 33](../security/real-life-security-guide-v1.md#chapter-33-ssrf-mastery).

### 48.1 Retries with backoff, jitter, and an idempotency key

A retry is only safe if repeating the request cannot do the work twice.
`GET` is safe by definition. A `POST /payments` is safe only if the server
deduplicates by an **idempotency key** that the client generates **once**,
before the first attempt. The demo below simulates the worst case: the
server charges the card, then the response is lost (here, a 503 sent after
the charge).

```rust
// src/bin/retries.rs
use axum::extract::State;
use axum::http::{HeaderMap, StatusCode};
use axum::{routing::post, Router};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::Duration;

#[derive(Default)]
struct Ledger { charges: u32, attempts: u32, seen: HashMap<String, String> }

async fn pay(State(l): State<Arc<Mutex<Ledger>>>, headers: HeaderMap) -> (StatusCode, String) {
    let Some(key) = headers.get("idempotency-key").and_then(|v| v.to_str().ok()).map(String::from) else {
        return (StatusCode::BAD_REQUEST, "Idempotency-Key required".into());
    };
    let mut l = l.lock().unwrap();
    l.attempts += 1;
    if let Some(prev) = l.seen.get(&key) {
        return (StatusCode::CREATED, prev.clone());          // replay the stored result; no second charge
    }
    l.charges += 1;
    let receipt = format!("receipt-{}", l.charges);
    l.seen.insert(key, receipt.clone());
    if l.attempts == 1 {
        return (StatusCode::SERVICE_UNAVAILABLE, "lost response".into()); // charged, but the client sees a failure
    }
    (StatusCode::CREATED, receipt)
}

fn backoff(attempt: u32) -> Duration {
    let cap_ms = 2_000u64;
    let exp_ms = 100u64.saturating_mul(1 << attempt.min(10)).min(cap_ms);
    Duration::from_millis(rand::random_range(0..=exp_ms))     // "full jitter": avoid synchronized retry storms
}

fn retryable(status: StatusCode) -> bool {
    matches!(status, StatusCode::TOO_MANY_REQUESTS | StatusCode::BAD_GATEWAY | StatusCode::SERVICE_UNAVAILABLE | StatusCode::GATEWAY_TIMEOUT)
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let ledger = Arc::new(Mutex::new(Ledger::default()));
    let app = Router::new().route("/payments", post(pay)).with_state(Arc::clone(&ledger));
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await?;
    let url = format!("http://{}/payments", listener.local_addr()?);
    tokio::spawn(async move { axum::serve(listener, app).await });

    let client = reqwest::Client::builder()
        .connect_timeout(Duration::from_secs(2))
        .timeout(Duration::from_secs(5))                       // total per-attempt deadline
        .build()?;
    let key = format!("{:032x}", rand::random::<u128>());       // generated ONCE, reused by every attempt

    let mut attempt = 0;
    let receipt = loop {
        let result = client.post(&url).header("Idempotency-Key", &key).body("amount=500").send().await;
        match result {
            Ok(r) if r.status().is_success() => break r.text().await?,
            Ok(r) if retryable(r.status()) && attempt < 4 => println!("attempt {attempt}: {}, retrying", r.status()),
            Err(e) if (e.is_timeout() || e.is_connect()) && attempt < 4 => println!("attempt {attempt}: {e}, retrying"),
            Ok(r) => anyhow::bail!("giving up: {}", r.status()),
            Err(e) => return Err(e.into()),
        }
        tokio::time::sleep(backoff(attempt)).await;
        attempt += 1;
    };
    let l = ledger.lock().unwrap();
    println!("{receipt}; server saw {} attempts and charged {} time(s)", l.attempts, l.charges);
    Ok(())
}
```

```text
attempt 0: 503 Service Unavailable, retrying
receipt-1; server saw 2 attempts and charged 1 time(s)
```

Without the key, the retry would have charged the card twice. In a real
service, store the key with the result in the same database transaction as
the charge, and expire keys after a retention period. Add a **retry
budget** (for example, retries may be at most 10% of requests) so a
struggling dependency does not receive three times its normal load. That is
the retry-storm failure described in
[HTTPS Ch 21](../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers).

### 48.2 An SSRF guard that survives DNS rebinding

A server that fetches user-supplied URLs (webhooks, link previews, image
imports) can be pointed at `http://169.254.169.254/` (cloud credentials),
`http://localhost:6379/` (an internal Redis), or a private network. Checking
the hostname string is not enough. `evil.example` can resolve to
`127.0.0.1`, or resolve to a public IP during your check and a private one
when the HTTP client connects (DNS rebinding). The fix: **resolve once,
check every address, and connect only to the address you checked.**

```rust
// src/bin/safefetch.rs
use std::net::{IpAddr, SocketAddr};
use std::time::Duration;

fn is_forbidden(ip: IpAddr) -> bool {
    match ip {
        IpAddr::V4(v4) => {
            let o = v4.octets();
            v4.is_private() || v4.is_loopback() || v4.is_link_local()  // 10/8 172.16/12 192.168/16 127/8 169.254/16
                || v4.is_unspecified() || v4.is_broadcast() || v4.is_documentation() || v4.is_multicast()
                || o[0] == 0                                           // 0.0.0.0/8 ("this network")
                || (o[0] == 100 && (o[1] & 0xC0) == 64)                // 100.64/10 carrier-grade NAT
                || (o[0] == 198 && (o[1] & 0xFE) == 18)                // 198.18/15 benchmarking
        }
        IpAddr::V6(v6) => {
            let s = v6.segments();
            v6.is_loopback() || v6.is_unspecified() || v6.is_multicast()
                || (s[0] & 0xFE00) == 0xFC00                           // fc00::/7 unique local
                || (s[0] & 0xFFC0) == 0xFE80                           // fe80::/10 link local
                || v6.to_ipv4_mapped().is_some_and(|v4| is_forbidden(IpAddr::V4(v4))) // ::ffff:127.0.0.1
        }
    }
}

async fn safe_get(raw_url: &str, max_bytes: usize) -> anyhow::Result<String> {
    let url = reqwest::Url::parse(raw_url)?;
    anyhow::ensure!(matches!(url.scheme(), "http" | "https"), "scheme {} not allowed", url.scheme());
    anyhow::ensure!(url.username().is_empty() && url.password().is_none(), "credentials in URL not allowed");
    let host = url.host_str().ok_or_else(|| anyhow::anyhow!("no host"))?;
    let port = url.port_or_known_default().unwrap_or(80);
    let bare = host.trim_start_matches('[').trim_end_matches(']');   // IPv6 literals come bracketed

    // 1. Resolve ONCE, with a deadline.
    let addrs: Vec<SocketAddr> = tokio::time::timeout(Duration::from_secs(2), tokio::net::lookup_host((bare, port)))
        .await??.collect();
    // 2. EVERY address must be allowed: attackers return [public, private] and hope you check only the first.
    anyhow::ensure!(!addrs.is_empty(), "{host} did not resolve");
    if let Some(bad) = addrs.iter().find(|a| is_forbidden(a.ip())) {
        anyhow::bail!("{host} resolves to forbidden address {}", bad.ip());
    }
    // 3. PIN the connection to the vetted address, so a second DNS answer cannot redirect it.
    let client = reqwest::Client::builder()
        .resolve(bare, addrs[0])
        .redirect(reqwest::redirect::Policy::none())   // a redirect to http://127.0.0.1 would skip all checks
        .connect_timeout(Duration::from_secs(3))
        .timeout(Duration::from_secs(10))
        .build()?;
    let mut resp = client.get(url).send().await?;
    anyhow::ensure!(resp.status().is_success(), "status {}", resp.status());

    // 4. Cap the body while streaming. Content-Length can be missing or false.
    let mut body = Vec::new();
    while let Some(chunk) = resp.chunk().await? {
        anyhow::ensure!(body.len() + chunk.len() <= max_bytes, "response larger than {max_bytes} bytes");
        body.extend_from_slice(&chunk);
    }
    Ok(String::from_utf8_lossy(&body).into_owned())
}

#[tokio::main]
async fn main() {
    for url in [
        "http://169.254.169.254/latest/meta-data/",   // cloud metadata
        "http://localhost:6379/",                      // resolves to loopback
        "http://[::ffff:127.0.0.1]/",                  // IPv4-mapped IPv6 loopback
        "http://0x7f000001/",                          // 127.0.0.1 in hex: the URL parser normalizes it
        "http://10.1.2.3/admin",                       // private network
        "file:///etc/passwd",                          // wrong scheme
        "http://user:pw@example.com/",                 // credentials in the URL
    ] {
        match safe_get(url, 1 << 20).await {
            Ok(body) => println!("ALLOWED {url} ({} bytes)", body.len()),
            Err(e) => println!("blocked {url:<42} {e}"),
        }
    }
}
```

```text
blocked http://169.254.169.254/latest/meta-data/   169.254.169.254 resolves to forbidden address 169.254.169.254
blocked http://localhost:6379/                     localhost resolves to forbidden address ::1
blocked http://[::ffff:127.0.0.1]/                 [::ffff:7f00:1] resolves to forbidden address ::ffff:127.0.0.1
blocked http://0x7f000001/                         127.0.0.1 resolves to forbidden address 127.0.0.1
blocked http://10.1.2.3/admin                      10.1.2.3 resolves to forbidden address 10.1.2.3
blocked file:///etc/passwd                         scheme file not allowed
blocked http://user:pw@example.com/                credentials in URL not allowed
```

Note the `0x7f000001` line. The WHATWG URL parser that `reqwest::Url` uses
normalizes hex, octal, and short-form IPv4 hosts before you see them, so the
check runs on the real address. Hand-written string checks such as
`host.starts_with("127.")` miss all of these encodings. This is why every
SSRF guide says to parse first and then check addresses, never strings.

A production version should also apply the check to every redirect hop
(implement `redirect::Policy::custom`), and run the fetcher in a network
segment with egress firewall rules as the second layer
([Security in Depth Ch 8](../security/real-life-security-guide-v1.md#chapter-8-cloud-network-security)).

---

## 49. TLS and mutual TLS with `rustls`

**Implements:** [Security from Zero Ch 28–29](../security/real-life-guide.md#chapter-28-the-tls-1-3-handshake-step-by-step)
and [Ch 53](../security/real-life-guide.md#chapter-53-project-1-your-own-ca-plus-mutual-tls)
(your own CA plus mTLS), [Go §50](../Golang/real-life-golang-guide.md#50-mtls-client-certificate-authentication-end-to-end).

`rustls` is a TLS library written in Rust. It supports only TLS 1.2 and
1.3 with modern cipher suites, so there is no SSLv3, RC4, or renegotiation
to disable. Its cryptography comes from a provider: `aws-lc-rs` by default,
or `ring`. The program below builds a whole PKI in memory: a CA, a server
certificate, and a client certificate. It then runs an mTLS server and
connects to it twice, once with a client certificate and once without.

```rust
// src/bin/mtls.rs
use rcgen::{BasicConstraints, CertificateParams, DnType, IsCa, Issuer, KeyPair, KeyUsagePurpose};
use rustls::pki_types::{CertificateDer, PrivateKeyDer, PrivatePkcs8KeyDer, ServerName};
use rustls::server::WebPkiClientVerifier;
use rustls::{ClientConfig, RootCertStore, ServerConfig};
use std::sync::Arc;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};
use tokio_rustls::{TlsAcceptor, TlsConnector};

struct Pki { ca: CertificateDer<'static>, server: (CertificateDer<'static>, PrivateKeyDer<'static>), client: (CertificateDer<'static>, PrivateKeyDer<'static>) }

fn leaf(name: &str, san: Vec<String>, issuer: &Issuer<'_, KeyPair>) -> anyhow::Result<(CertificateDer<'static>, PrivateKeyDer<'static>)> {
    let key = KeyPair::generate()?;                                   // ECDSA P-256 by default
    let mut params = CertificateParams::new(san)?;
    params.distinguished_name.push(DnType::CommonName, name);
    let cert = params.signed_by(&key, issuer)?;
    Ok((cert.der().clone(), PrivatePkcs8KeyDer::from(key.serialize_der()).into()))
}

fn make_pki() -> anyhow::Result<Pki> {
    let ca_key = KeyPair::generate()?;
    let mut ca_params = CertificateParams::new(Vec::<String>::new())?;
    ca_params.distinguished_name.push(DnType::CommonName, "Rust Guide Lab CA");
    ca_params.is_ca = IsCa::Ca(BasicConstraints::Unconstrained);
    ca_params.key_usages = vec![KeyUsagePurpose::KeyCertSign, KeyUsagePurpose::CrlSign];
    let ca_cert = ca_params.self_signed(&ca_key)?;
    let issuer = Issuer::new(ca_params, ca_key);
    Ok(Pki {
        ca: ca_cert.der().clone(),
        server: leaf("localhost", vec!["localhost".into()], &issuer)?,
        client: leaf("billing-service", vec!["billing.internal".into()], &issuer)?,
    })
}

async fn connect(addr: std::net::SocketAddr, roots: Arc<RootCertStore>, identity: Option<(CertificateDer<'static>, PrivateKeyDer<'static>)>) -> anyhow::Result<String> {
    let builder = ClientConfig::builder().with_root_certificates(roots);
    let config = match identity {
        Some((cert, key)) => builder.with_client_auth_cert(vec![cert], key)?,
        None => builder.with_no_client_auth(),
    };
    let tcp = TcpStream::connect(addr).await?;
    let mut tls = TlsConnector::from(Arc::new(config)).connect(ServerName::try_from("localhost")?, tcp).await?;
    tls.write_all(b"ping").await?;
    let mut buf = Vec::new();
    tls.read_to_end(&mut buf).await?;     // in TLS 1.3, a rejected client cert surfaces here, after the handshake
    Ok(String::from_utf8(buf)?)
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let pki = make_pki()?;
    let mut roots = RootCertStore::empty();
    roots.add(pki.ca.clone())?;           // trust ONLY our lab CA, not the system store
    let roots = Arc::new(roots);

    let verifier = WebPkiClientVerifier::builder(Arc::clone(&roots)).build()?; // client certs REQUIRED
    let server_cfg = ServerConfig::builder()
        .with_client_cert_verifier(verifier)
        .with_single_cert(vec![pki.server.0.clone()], pki.server.1.clone_key())?;
    let acceptor = TlsAcceptor::from(Arc::new(server_cfg));

    let listener = TcpListener::bind("127.0.0.1:0").await?;
    let addr = listener.local_addr()?;
    tokio::spawn(async move {
        loop {
            let (tcp, _) = listener.accept().await.unwrap();
            let acceptor = acceptor.clone();
            tokio::spawn(async move {
                match acceptor.accept(tcp).await {
                    Ok(mut tls) => {
                        let (_, conn) = tls.get_ref();
                        let info = format!("pong over {:?} / {:?}, client chain of {} cert(s)",
                            conn.protocol_version().unwrap(),
                            conn.negotiated_cipher_suite().unwrap().suite(),
                            conn.peer_certificates().map_or(0, |c| c.len()));
                        let mut buf = [0u8; 4];
                        let _ = tls.read_exact(&mut buf).await;
                        let _ = tls.write_all(info.as_bytes()).await;
                        let _ = tls.shutdown().await;
                    }
                    Err(e) => eprintln!("server: handshake rejected: {e}"),
                }
            });
        }
    });

    println!("with client cert:    {:?}", connect(addr, Arc::clone(&roots), Some(pki.client)).await);
    println!("without client cert: {:?}", connect(addr, roots, None).await.map_err(|e| e.to_string()));
    Ok(())
}
```

```text
with client cert:    Ok("pong over TLSv1_3 / TLS13_AES_256_GCM_SHA384, client chain of 1 cert(s)")
server: handshake rejected: peer sent no certificates
without client cert: Err("received fatal alert: CertificateRequired")
```

What to take from this:

- **The client fails after the handshake, not during it.** In TLS 1.3 the
  client sends its (empty) certificate in its final flight and considers
  the handshake complete. The server's rejection arrives as an alert on
  the first read. This is the same TLS 1.3 behavior as in Go and OpenSSL,
  and it confuses everyone the first time they debug mTLS.
- **`with_root_certificates(roots)` trusts only your CA.** Service-to-service
  mTLS should never trust the public web PKI. This is the model behind
  SPIFFE/SPIRE ([Security in Depth Ch 23](../security/real-life-security-guide-v1.md#chapter-23-service-identity-spiffe-and-spire)).
- In production, load PEM files with `rustls::pki_types::pem::PemObject`
  (`CertificateDer::pem_file_iter`, `PrivateKeyDer::from_pem_file`), and
  reload certificates without a restart by implementing
  `ResolvesServerCert`. Short-lived certificates need reloading
  ([HTTPS Ch 16](../v2-https/real-life-guide-v1.md#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation)).

---

## 50. A reverse proxy with Hyper, and request-smuggling defenses

**Implements:** [HTTPS Ch 8](../v2-https/real-life-guide-v1.md#chapter-8-proxies-the-middlemen),
[Ch 20](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling),
and [Go §39](../Golang/real-life-golang-guide.md#39-terminal-project-a-rate-limited-reverse-proxy-mini-api-gateway).

A reverse proxy is where framing bugs and identity bugs meet. This one uses
Hyper directly, the library under Axum, reqwest, and Tonic, so you can see
the connection-level settings that frameworks hide.

```rust
// src/bin/revproxy.rs
use http_body_util::{combinators::BoxBody, BodyExt, Full};
use hyper::body::{Bytes, Incoming};
use hyper::header::{HeaderName, HeaderValue};
use hyper::{Request, Response, StatusCode, Uri};
use hyper_util::client::legacy::{connect::HttpConnector, Client};
use hyper_util::rt::{TokioExecutor, TokioIo, TokioTimer};
use std::net::SocketAddr;
use std::time::Duration;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};

type Body = BoxBody<Bytes, hyper::Error>;

// RFC 9110 §7.6.1: these describe ONE hop and must never be forwarded.
const HOP_BY_HOP: [&str; 8] = ["connection", "keep-alive", "proxy-connection", "te", "trailer",
                               "transfer-encoding", "upgrade", "proxy-authorization"];

fn text(status: StatusCode, msg: &'static str) -> Response<Body> {
    let mut r = Response::new(Full::new(Bytes::from(msg)).map_err(|never| match never {}).boxed());
    *r.status_mut() = status;
    r
}

async fn proxy(mut req: Request<Incoming>, peer: SocketAddr, upstream: &'static str,
               client: Client<HttpConnector, Incoming>) -> Result<Response<Body>, hyper::Error> {
    // 1. Remove hop-by-hop headers, including any the client NAMED in its Connection header.
    let named: Vec<HeaderName> = req.headers().get_all("connection").iter()
        .filter_map(|v| v.to_str().ok())
        .flat_map(|v| v.split(',').map(|s| s.trim().to_ascii_lowercase()))
        .filter_map(|s| HeaderName::from_bytes(s.as_bytes()).ok())
        .collect();
    for h in HOP_BY_HOP.iter().map(|h| HeaderName::from_static(h)).chain(named) {
        req.headers_mut().remove(h);
    }

    // 2. Identity: this proxy is the EDGE, so client-supplied forwarding headers are lies. Overwrite them.
    req.headers_mut().remove("forwarded");
    req.headers_mut().insert("x-forwarded-for", HeaderValue::from_str(&peer.ip().to_string()).unwrap());
    req.headers_mut().insert("x-forwarded-proto", HeaderValue::from_static("http"));

    // 3. Point the request at the upstream.
    let path = req.uri().path_and_query().map_or("/", |p| p.as_str());
    *req.uri_mut() = match format!("http://{upstream}{path}").parse::<Uri>() {
        Ok(u) => u,
        Err(_) => return Ok(text(StatusCode::BAD_REQUEST, "bad uri")),
    };

    // 4. Deadline for the upstream. The body streams through; nothing is buffered whole.
    match tokio::time::timeout(Duration::from_secs(10), client.request(req)).await {
        Ok(Ok(resp)) => {
            let (mut parts, body) = resp.into_parts();
            for h in HOP_BY_HOP { parts.headers.remove(h); }
            Ok(Response::from_parts(parts, body.boxed()))
        }
        Ok(Err(_)) => Ok(text(StatusCode::BAD_GATEWAY, "upstream error")),
        Err(_) => Ok(text(StatusCode::GATEWAY_TIMEOUT, "upstream timeout")),
    }
}

async fn backend(listener: TcpListener) {
    // A tiny upstream that reports what it received.
    loop {
        let (tcp, _) = listener.accept().await.unwrap();
        tokio::spawn(async move {
            let svc = hyper::service::service_fn(|req: Request<Incoming>| async move {
                let xff = req.headers().get("x-forwarded-for").cloned();
                let secret = req.headers().contains_key("x-internal-secret");
                let len = req.into_body().collect().await?.to_bytes().len();
                let msg = format!("backend: xff={xff:?} body={len}B internal-secret-leaked={secret}");
                Ok::<_, hyper::Error>(Response::new(Full::new(Bytes::from(msg))))
            });
            let _ = hyper::server::conn::http1::Builder::new().serve_connection(TokioIo::new(tcp), svc).await;
        });
    }
}

async fn raw(addr: SocketAddr, request: &str) -> String {
    let mut s = TcpStream::connect(addr).await.unwrap();
    s.write_all(request.as_bytes()).await.unwrap();
    let mut out = Vec::new();
    let _ = tokio::time::timeout(Duration::from_secs(2), s.read_to_end(&mut out)).await;
    let text = String::from_utf8_lossy(&out);
    let status = text.lines().next().unwrap_or("<no response>").to_string();
    let body = text.split("\r\n\r\n").nth(1).unwrap_or("").to_string();
    format!("{status} | {body}")
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let be = TcpListener::bind("127.0.0.1:0").await?;
    let upstream: &'static str = Box::leak(be.local_addr()?.to_string().into_boxed_str());
    tokio::spawn(backend(be));

    let client: Client<HttpConnector, Incoming> = Client::builder(TokioExecutor::new())
        .pool_idle_timeout(Duration::from_secs(30))
        .build_http();
    let front = TcpListener::bind("127.0.0.1:0").await?;
    let front_addr = front.local_addr()?;
    tokio::spawn(async move {
        loop {
            let (tcp, peer) = front.accept().await.unwrap();
            let client = client.clone();
            tokio::spawn(async move {
                let svc = hyper::service::service_fn(move |req| proxy(req, peer, upstream, client.clone()));
                let _ = hyper::server::conn::http1::Builder::new()
                    .timer(TokioTimer::new())
                    .header_read_timeout(Duration::from_secs(5))  // Slowloris defense, which axum::serve lacks
                    .max_buf_size(16 * 1024)                      // caps the request head
                    .serve_connection(TokioIo::new(tcp), svc)
                    .await;
            });
        }
    });

    println!("spoofed XFF:  {}", raw(front_addr, "GET / HTTP/1.1\r\nHost: x\r\nX-Forwarded-For: 10.0.0.1\r\nConnection: close\r\n\r\n").await);
    println!("Connection-named header: {}", raw(front_addr,
        "GET / HTTP/1.1\r\nHost: x\r\nX-Internal-Secret: 1\r\nConnection: close, X-Internal-Secret\r\n\r\n").await);
    println!("CL + TE:      {}", raw(front_addr,
        "POST / HTTP/1.1\r\nHost: x\r\nContent-Length: 6\r\nTransfer-Encoding: chunked\r\nConnection: close\r\n\r\n0\r\n\r\nGET /admin HTTP/1.1\r\n\r\n").await);
    println!("bad CL:       {}", raw(front_addr, "POST / HTTP/1.1\r\nHost: x\r\nContent-Length: 5, 6\r\n\r\nhello!").await);
    println!("space-colon:  {}", raw(front_addr, "POST / HTTP/1.1\r\nHost: x\r\nContent-Length : 5\r\n\r\nhello").await);
    Ok(())
}
```

The output of the demo, run as written:

```text
spoofed XFF:  HTTP/1.1 200 OK | backend: xff=Some("127.0.0.1") body=0B internal-secret-leaked=false
Connection-named header: HTTP/1.1 200 OK | backend: xff=Some("127.0.0.1") body=0B internal-secret-leaked=false
CL + TE:      HTTP/1.1 200 OK | backend: xff=Some("127.0.0.1") body=0B internal-secret-leaked=false
bad CL:       HTTP/1.1 400 Bad Request |
space-colon:  HTTP/1.1 400 Bad Request |
```

Read the `CL + TE` line carefully. Hyper did not reject the request.
RFC 9112 §6.3 says that when both headers are present, `Transfer-Encoding`
wins and `Content-Length` must be ignored. Hyper followed that rule: it read
a zero-length chunked body, and because the request said
`Connection: close`, it never read the smuggled `GET /admin` that followed.
The proxy then forwarded the request with **its own** framing. The
hop-by-hop `Transfer-Encoding` was stripped, and Hyper's client re-framed
the body. So the backend never sees the ambiguous message, and both
parsers agree. That is the real defense against desync: **terminate and
re-frame at every hop, with one strict parser.** Proxies that pass raw
bytes through, or that use two different parsers, are the ones that get
smuggled through. Some proxies go further and reject CL+TE outright with
400, as Ch 46 does. Both behaviors are allowed by the RFC. Forwarding the
ambiguous message unchanged is not.

The other rows show the remaining controls. The client's
`X-Forwarded-For: 10.0.0.1` was replaced, so the backend sees the real
peer, and an IP allow-list based on XFF cannot be bypassed through this
proxy
([HTTPS Ch 20](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling)).
The header the client listed in `Connection:` was stripped, which closes a
known trick for deleting a header that a proxy adds and the backend trusts.
Malformed lengths and whitespace before the colon get a 400 from Hyper's
parser before your code runs.

**Stretch goals:** add a per-client-IP token-bucket rate limiter
([Security in Depth Ch 40B](../security/real-life-security-guide-v1.md#chapter-40b-api-rate-limiting-in-depth)),
round-robin over several upstreams with passive health checks
([HTTPS Ch 9](../v2-https/real-life-guide-v1.md#chapter-9-load-balancers-rate-limiting)),
and terminate TLS using the acceptor from Ch 49.

---

# Part VI — Security-focused Rust

Rust removes the bug classes behind most of the
[C guide's Part V](../c-lang/real-life-c-guide.md#28-undefined-behavior-the-list-every-c-programmer-must-memorize).
It does not remove the rest of the
[OWASP Top 10](../security/real-life-guide.md#chapter-41-how-the-web-decides-what-to-trust).
This Part is precise about where the line falls, then covers the security
work that is still yours: cryptography, input handling, testing hostile
input, and dependencies.

---

## 51. What Rust prevents, and what it does not

| Bug class | In safe Rust | Notes |
|---|---|---|
| Buffer overflow / out-of-bounds read | **prevented** | bounds checks; a bad index panics |
| Use-after-free, double free, dangling pointers | **prevented** | ownership and borrowing |
| Data races | **prevented** | `Send`/`Sync` |
| Uninitialized memory reads | **prevented** | definite-assignment checking |
| Null-pointer dereference | **prevented** | `Option` |
| Type confusion | **prevented** | no unchecked casts outside `unsafe` |
| Integer overflow | **defined, not prevented** | panics in debug, wraps in release by default (§3.3) |
| Panics (`unwrap`, indexing, overflow) | **possible** | a remotely triggerable panic is a DoS bug |
| Deadlocks, logical race conditions | **possible** | memory-safe, therefore allowed (§29.2) |
| Memory leaks | **possible** | `Rc` cycles, `mem::forget`, unbounded caches |
| Injection (SQL, shell, template, log) | **possible** | use parameterized APIs; types can help (§53) |
| Broken access control, auth logic | **possible** | the most common serious web bug, in any language |
| Timing side channels | **possible** | `==` on secrets leaks timing (§52.3) |
| Secrets in logs | **possible** | `#[derive(Debug)]` prints every field (below) |
| Memory-safety bugs in `unsafe` or C dependencies | **possible** | audit `unsafe`, run Miri (§33.4) |
| Malicious or vulnerable dependencies | **possible** | Ch 55 |

**The real-world numbers.** Google reported that as Android shifted new
native code to Rust and other memory-safe languages, memory-safety bugs
fell from 76% of Android's vulnerabilities in 2019 to 24% in 2024. Google
had reported in 2022 that, by then, no memory-safety vulnerabilities had
been found in Android's Rust code.
The bugs that remain in Rust codebases are the right-hand column of this
table: logic, authorization, injection, and DoS through panics.

### 51.1 Secrets leaking through `Debug`

```rust
use std::fmt;

#[derive(Debug)]
struct LoginBad { user: String, password: String }

/// A wrapper whose Debug and Display never reveal the value.
#[derive(Clone)]
struct Secret<T>(T);
impl<T> Secret<T> { fn expose(&self) -> &T { &self.0 } }   // the only way in: greppable in review
impl<T> fmt::Debug for Secret<T> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result { f.write_str("Secret([REDACTED])") }
}

#[derive(Debug)]
struct LoginGood { user: String, password: Secret<String> }

fn main() {
    let bad = LoginBad { user: "ada".into(), password: "hunter2".into() };
    let good = LoginGood { user: "ada".into(), password: Secret("hunter2".into()) };
    println!("{bad:?}");   // LoginBad { user: "ada", password: "hunter2" }  <- ends up in your logs
    println!("{good:?}");  // LoginGood { user: "ada", password: Secret([REDACTED]) }
    println!("{}", good.password.expose().len());
}
```

`tracing::info!(?request, "login")` debug-prints the whole struct, so one
`#[derive(Debug)]` on a request type can send passwords to your log
pipeline. The `secrecy` crate provides a production version of `Secret<T>`
(`SecretString`) that also zeroizes the value on drop.
[Security from Zero Ch 49](../security/real-life-guide.md#chapter-49-secrets-management)
covers the operational side.

### 51.2 Panics are a denial-of-service surface

`let n: usize = header.parse().unwrap()`, `&body[..4]`, `a / b` with
`b == 0`, and `x + y` overflowing in debug builds are all panics an attacker
can trigger. Rules for code that handles untrusted input:

- Use `get`, `checked_*`, `?`, and `split_at_checked` instead of indexing
  and `unwrap`.
- Turn on `clippy::unwrap_used`, `clippy::expect_used`,
  `clippy::indexing_slicing`, and `clippy::arithmetic_side_effects` for
  parser modules (`#![deny(...)]` at the top of the module).
- Fuzz the parser (Ch 54). Fuzzers find panics within minutes.

---

## 52. Cryptography done right: AEAD, password hashing, HMAC, constant-time, zeroize

The rule from [Security from Zero](../security/real-life-guide.md#chapter-15-aead-encryption-and-authentication-together)
applies in every language: **do not design cryptography, and do not
assemble primitives yourself.** Use an AEAD for encryption, a password hash
for passwords, and HMAC for message authentication, from audited crates.

| Need | Crate | Notes |
|---|---|---|
| Symmetric encryption | `aes-gcm`, `chacha20poly1305` (RustCrypto); `aws-lc-rs`, `ring` | AEAD only. Nonces must never repeat under one key |
| Password hashing | `argon2` | Argon2id; `scrypt` and `bcrypt` are acceptable alternatives |
| MAC | `hmac` + `sha2` | verify in constant time |
| Signatures | `ed25519-dalek`, `p256`, `aws-lc-rs` | |
| TLS | `rustls` | Ch 49 |
| Random | `rand` (`OsRng` / `rand::random`), `getrandom` | never a seeded PRNG for secrets |
| Constant time, wiping memory | `subtle`, `zeroize` | |

### 52.1 AES-256-GCM with associated data

```rust
use aes_gcm::aead::{Aead, AeadCore, Generate, Key, KeyInit, Payload};
use aes_gcm::{Aes256Gcm, Nonce};

const NONCE_LEN: usize = 12;

/// Output format: nonce (12 bytes) || ciphertext || tag (16 bytes).
fn seal(key: &Key<Aes256Gcm>, plaintext: &[u8], aad: &[u8]) -> Vec<u8> {
    let cipher = Aes256Gcm::new(key);
    let nonce = Nonce::generate();                      // 96 random bits per message
    let ct = cipher.encrypt(&nonce, Payload { msg: plaintext, aad }).expect("encryption cannot fail for in-memory input");
    let mut out = Vec::with_capacity(NONCE_LEN + ct.len());
    out.extend_from_slice(&nonce);
    out.extend_from_slice(&ct);
    out
}

fn open(key: &Key<Aes256Gcm>, sealed: &[u8], aad: &[u8]) -> Result<Vec<u8>, &'static str> {
    let (nonce, ct) = sealed.split_at_checked(NONCE_LEN).ok_or("too short")?;
    let nonce = Nonce::try_from(nonce).map_err(|_| "bad nonce")?;
    Aes256Gcm::new(key)
        .decrypt(&nonce, Payload { msg: ct, aad })
        .map_err(|_| "authentication failed")           // one opaque error: no oracle for attackers
}

fn main() {
    let key = Key::<Aes256Gcm>::generate();             // in production: from a KMS, never hard-coded
    // AAD binds the ciphertext to its context: this blob is user 42's SSN, nothing else.
    let sealed = seal(&key, b"123-45-6789", b"user:42/field:ssn");
    println!("sealed {} bytes", sealed.len());          // 12 + 11 + 16 = 39
    println!("{:?}", open(&key, &sealed, b"user:42/field:ssn").map(String::from_utf8));
    println!("{:?}", open(&key, &sealed, b"user:43/field:ssn"));   // moved to another row: rejected

    let mut tampered = sealed.clone();
    tampered[15] ^= 1;                                   // flip one ciphertext bit
    println!("{:?}", open(&key, &tampered, b"user:42/field:ssn"));
}
```

```text
sealed 39 bytes
Ok(Ok("123-45-6789"))
Err("authentication failed")
Err("authentication failed")
```

Random 96-bit nonces are safe up to about 2³² messages per key. Beyond that,
the chance of a nonce repeating becomes significant, and a repeated nonce
breaks GCM completely. High-volume systems rotate keys or use
envelope encryption with a fresh data key per object
([Security in Depth Ch 9](../security/real-life-security-guide-v1.md#chapter-9-data-security-in-the-cloud-kms-and-envelope-encryption)).

### 52.2 Password hashing with Argon2id

```rust
use argon2::password_hash::{phc::PasswordHash, PasswordHasher, PasswordVerifier};
use argon2::{Algorithm, Argon2, Params, Version};

fn hasher() -> Argon2<'static> {
    // OWASP's minimum for Argon2id: 19 MiB memory, 2 iterations, 1 lane. Tune upward to about
    // 100–500 ms per hash on your production hardware.
    Argon2::new(Algorithm::Argon2id, Version::V0x13, Params::new(19 * 1024, 2, 1, None).unwrap())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let stored = hasher().hash_password(b"correct horse battery staple")?.to_string();
    println!("{stored}");      // $argon2id$v=19$m=19456,t=2,p=1$<salt>$<hash>: params and salt travel with the hash

    let parsed = PasswordHash::new(&stored)?;
    let t = std::time::Instant::now();
    println!("right password: {}", hasher().verify_password(b"correct horse battery staple", &parsed).is_ok());
    println!("wrong password: {}", hasher().verify_password(b"Tr0ub4dor&3", &parsed).is_ok());
    println!("two verifications took {:?}", t.elapsed()); // slow on purpose: that is what stops offline cracking
    Ok(())
}
```

In a debug build the two verifications took about 430 ms. Optimized builds
are several times faster, which is why you tune the parameters against a
`--release` binary on production hardware. Argon2 is deliberately slow and
memory-hard, which is the point. In an
async server, run `hash_password`/`verify_password` inside
`spawn_blocking` (§31.1). Otherwise every login stalls a Tokio worker for
the length of a hash. [Security from Zero Ch 10](../security/real-life-guide.md#chapter-10-password-storage-is-a-completely-different-problem)
explains why fast hashes such as SHA-256 are wrong for passwords.

### 52.3 HMAC webhook verification, in constant time

```rust
use hmac::{Hmac, KeyInit, Mac};
use sha2::Sha256;
use subtle::ConstantTimeEq;

type HmacSha256 = Hmac<Sha256>;

fn hex_decode(s: &str) -> Option<Vec<u8>> {
    if s.len() % 2 != 0 { return None; }
    (0..s.len()).step_by(2).map(|i| u8::from_str_radix(s.get(i..i + 2)?, 16).ok()).collect()
}

/// Verifies a GitHub-style header: "X-Hub-Signature-256: sha256=<hex>".
fn verify_webhook(secret: &[u8], body: &[u8], header: &str) -> bool {
    let Some(sig) = header.strip_prefix("sha256=").and_then(hex_decode) else { return false };
    let mut mac = HmacSha256::new_from_slice(secret).expect("HMAC accepts any key length");
    mac.update(body);
    mac.verify_slice(&sig).is_ok()      // constant-time comparison inside
}

fn sign(secret: &[u8], body: &[u8]) -> String {
    let mut mac = HmacSha256::new_from_slice(secret).unwrap();
    mac.update(body);
    let tag = mac.finalize().into_bytes();
    format!("sha256={}", tag.iter().map(|b| format!("{b:02x}")).collect::<String>())
}

fn main() {
    let secret = b"whsec_lab_only";
    let body = br#"{"action":"push","ref":"main"}"#;
    let header = sign(secret, body);
    println!("{header}");
    println!("valid:      {}", verify_webhook(secret, body, &header));
    println!("tampered:   {}", verify_webhook(secret, br#"{"action":"push","ref":"prod"}"#, &header));
    println!("garbage:    {}", verify_webhook(secret, body, "sha256=zz"));

    // Comparing secrets yourself? Use subtle, never `==` (which can stop at the first differing byte).
    let (token, presented) = (b"api-key-7f3a".as_slice(), b"api-key-7f3b".as_slice());
    println!("token match: {}", bool::from(token.ct_eq(presented)));
}
```

`==` on byte slices returns as soon as it finds a difference, so response
time reveals how many leading bytes were correct. Over many requests an
attacker can recover a token byte by byte. `verify_slice` and
`subtle::ConstantTimeEq` take the same time whatever the input.
[Security from Zero Ch 11](../security/real-life-guide.md#chapter-11-hmac-proving-a-message-wasn-t-tampered-with)
covers HMAC itself.

### 52.4 Wiping secrets from memory with `zeroize`

```rust
use zeroize::{Zeroize, ZeroizeOnDrop};

#[derive(Zeroize, ZeroizeOnDrop)]
struct SessionKeys { enc: [u8; 32], mac: [u8; 32] }

fn main() {
    let mut pin = String::from("4921");
    pin.zeroize();                                    // overwritten with zeros, with writes the optimizer may not remove
    println!("pin after zeroize: {pin:?}");           // ""
    let keys = SessionKeys { enc: [7; 32], mac: [9; 32] };
    println!("using keys {} {}", keys.enc[0], keys.mac[0]);
}   // `keys` is wiped here, automatically
```

Zeroizing reduces how long a secret sits in memory, where it could end up
in a core dump, swap, or a heartbleed-style over-read. It is not a
guarantee. Moves and `Vec` reallocation can leave old copies behind, so
allocate secret buffers once at their final size.

---

## 53. Parsing untrusted input: newtypes, "parse, don't validate", resource bounds

### 53.1 Parse, don't validate

A `validate(&str) -> bool` function leaves the value as a `String`, and
nothing stops the next developer from skipping the call. **Parsing** turns
untrusted input into a type that can only hold valid values. After the
boundary, the type is the proof that validation happened.

```rust
use serde::Deserialize;
use std::fmt;

#[derive(Debug, Clone, PartialEq, Eq, Hash, Deserialize)]
#[serde(try_from = "String")]                        // deserialize a String, then run TryFrom on it
pub struct Username(String);                        // private field: the only way to make one is TryFrom

impl TryFrom<String> for Username {
    type Error = String;
    fn try_from(s: String) -> Result<Self, String> {
        let ok_len = (3..=32).contains(&s.len());
        let ok_chars = s.bytes().all(|b| b.is_ascii_lowercase() || b.is_ascii_digit() || b == b'_');
        let ok_start = s.bytes().next().is_some_and(|b| b.is_ascii_lowercase());
        if ok_len && ok_chars && ok_start { Ok(Username(s)) } else { Err(format!("invalid username {s:?}")) }
    }
}
impl fmt::Display for Username {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result { f.write_str(&self.0) }
}

#[derive(Debug, Deserialize)]
#[serde(deny_unknown_fields)]
struct SignupRequest {
    username: Username,                               // validated during deserialization
    #[serde(default)]
    marketing_opt_in: bool,
}

fn create_home_dir(user: &Username) -> String {
    // Path traversal ("../../etc") is impossible here: the TYPE guarantees [a-z0-9_] only.
    format!("/home/{user}")
}

fn main() {
    for body in [
        r#"{"username":"ada_l"}"#,
        r#"{"username":"../../etc/passwd"}"#,
        r#"{"username":"Ada"}"#,
        r#"{"username":"ada_l","is_admin":true}"#,
    ] {
        match serde_json::from_str::<SignupRequest>(body) {
            Ok(req) => println!("ok   -> {} (opt-in {})", create_home_dir(&req.username), req.marketing_opt_in),
            Err(e) => println!("rejected: {e}"),
        }
    }
}
```

```text
ok   -> /home/ada_l (opt-in false)
rejected: invalid username "../../etc/passwd" at line 1 column 31
rejected: invalid username "Ada" at line 1 column 18
rejected: unknown field `is_admin`, expected `username` or `marketing_opt_in` at line 1 column 30
```

The same idea works for every input that matters: `Email`, `TenantId`,
`NonEmptyVec`, `Port`, `Cents(u64)` instead of `f64` for money, and an
`HtmlEscaped` type that only an escaping function can produce. When
`create_home_dir` takes `&Username` instead of `&str`, the compiler stops
any unvalidated string from reaching it.

### 53.2 Bounds on every resource

Untrusted input controls how much memory and CPU your parser uses unless
you bound it. Check each of these:

| Resource | Bound | Where you have seen it |
|---|---|---|
| Body size | `RequestBodyLimitLayer`, a `MAX_FRAME` check | §32.2, §47 |
| Header size | `max_buf_size`, `MAX_HEAD` | §46, §50 |
| Length fields | compare to a maximum **before** allocating | §32.2 |
| Nesting depth | `serde_json` stops at 128 levels by default; keep that limit | — |
| Collection counts | `Vec::with_capacity(n)` with `n` taken from input is an allocation bomb; cap `n` | §45 (DNS counts are not trusted) |
| Regex | the `regex` crate guarantees linear time (no ReDoS) | §26 |
| Decompression | cap output bytes, not input bytes (zip bombs) | — |
| Time | timeouts on reads, handlers, and outbound calls | §31, §47, §48 |

### 53.3 Injection: let the types build the query

```rust
// fragment: sqlx (async, compile-time checked SQL)
// BUG: SQL injection. Rust's safety does not help with string formatting.
let q = format!("SELECT * FROM users WHERE email = '{email}'");
sqlx::query(&q).fetch_one(&pool).await?;

// FIX: a bound parameter. The value never becomes SQL text.
sqlx::query_as::<_, User>("SELECT * FROM users WHERE email = $1")
    .bind(&email)
    .fetch_one(&pool)
    .await?;

// sqlx::query! checks the SQL against your real schema AT COMPILE TIME.
let user = sqlx::query!("SELECT id, email FROM users WHERE email = $1", email).fetch_one(&pool).await?;
```

[Security from Zero Ch 42](../security/real-life-guide.md#chapter-42-injection-when-data-becomes-code)
covers injection in general. The Rust-specific point is that `format!` is
string concatenation, so the language does not save you from SQL injection.
APIs that take parameters separately (`sqlx`, `Command::arg` in §24.3, and
auto-escaping templates such as `askama`) do.

---

## 54. Fuzzing and property testing

### 54.1 Coverage-guided fuzzing with `cargo-fuzz`

A fuzzer runs your parser millions of times with mutated inputs and keeps
the inputs that reach new code paths. In Rust it mostly finds **panics**
(slice indexing, `unwrap`, overflow) and infinite loops. In `unsafe` code it
also finds memory bugs, because the target is compiled with
AddressSanitizer. Same idea as
[C §31's libFuzzer chapter](../c-lang/real-life-c-guide.md#31-fuzzing-a-real-parsing-bug-with-libfuzzer).

```bash
cargo install cargo-fuzz          # needs a nightly toolchain: rustup toolchain install nightly
cargo fuzz init                   # creates fuzz/ with its own Cargo.toml
cargo fuzz add dns_parse
```

```rust
// fragment: fuzz/fuzz_targets/dns_parse.rs
#![no_main]
use libfuzzer_sys::fuzz_target;

fuzz_target!(|data: &[u8]| {
    // The only property: the parser must never panic or hang, whatever the bytes.
    let _ = rust_labs::dns::parse_response(data, u16::from_be_bytes([data.first().copied().unwrap_or(0), 0]));
});
```

```bash
cargo +nightly fuzz run dns_parse -- -max_total_time=300   # 5 minutes
# On a crash: fuzz/artifacts/dns_parse/crash-<hash> holds the input. Reproduce it with:
cargo +nightly fuzz run dns_parse fuzz/artifacts/dns_parse/crash-<hash>
```

Move the parser into `src/lib.rs` (as `pub mod dns`) so both the binary
and the fuzz target can use it. Try it on the Ch 45 parser with the
`jumps > 16` check removed. The fuzzer finds the pointer loop, reported as a
timeout, within seconds.

### 54.2 Property tests for invariants

Fuzzing asks "does it crash?". Property tests check that a rule holds, for
example that parsing a value you serialized gives back the same value:

```rust
#[derive(Debug, Clone, PartialEq)]
struct Frame { stream: u32, payload: Vec<u8> }

fn encode(f: &Frame) -> Vec<u8> {
    let mut out = f.stream.to_be_bytes().to_vec();
    out.extend_from_slice(&(f.payload.len() as u32).to_be_bytes());
    out.extend_from_slice(&f.payload);
    out
}

fn decode(b: &[u8]) -> Option<(Frame, &[u8])> {
    let (stream, rest) = b.split_at_checked(4)?;
    let (len, rest) = rest.split_at_checked(4)?;
    let len = u32::from_be_bytes(len.try_into().ok()?) as usize;
    let (payload, rest) = rest.split_at_checked(len)?;
    Some((Frame { stream: u32::from_be_bytes(stream.try_into().ok()?), payload: payload.to_vec() }, rest))
}

#[cfg(test)]
mod tests {
    use super::*;
    use proptest::prelude::*;

    proptest! {
        #[test]
        fn roundtrip(stream: u32, payload in proptest::collection::vec(any::<u8>(), 0..2048)) {
            let f = Frame { stream, payload };
            let wire = encode(&f);
            prop_assert_eq!(decode(&wire), Some((f, &[][..])));
        }

        #[test]
        fn never_panics_and_never_overreads(bytes in proptest::collection::vec(any::<u8>(), 0..64)) {
            if let Some((f, rest)) = decode(&bytes) {
                prop_assert_eq!(8 + f.payload.len() + rest.len(), bytes.len());
            }
        }
    }
}

fn main() {
    let f = Frame { stream: 1, payload: b"hi".to_vec() };
    println!("{:?}", decode(&encode(&f)));
}
```

Use both. Property tests run in milliseconds in `cargo test` on every
commit. Fuzzers run for hours in CI or on OSS-Fuzz, which accepts Rust
projects.

---

## 55. Supply chain: `Cargo.lock`, `cargo audit`, `cargo deny`, `build.rs`

A typical Axum service has 150 to 300 transitive dependencies. Each one is
code you run, and some run **at build time**: `build.rs` scripts and
procedural macros execute on the developer's laptop and in CI with full
access to the filesystem and network. That is the main supply-chain risk
specific to Rust
([Security from Zero Ch 47](../security/real-life-guide.md#chapter-47-supply-chain-and-the-code-you-didn-t-write),
[Security in Depth Ch 12–13](../security/real-life-security-guide-v1.md#chapter-12-ci-cd-as-an-attack-surface)).

### 55.1 The controls, in the order to adopt them

```bash
# 1. Reproducible resolution: commit Cargo.lock, and make CI fail if it would change.
cargo build --locked

# 2. Known vulnerabilities, from the RustSec advisory database.
cargo install cargo-audit && cargo audit

# 3. Policy: advisories + licenses + banned crates + duplicate versions + allowed sources.
cargo install cargo-deny && cargo deny init && cargo deny check

# 4. See what you actually pull in, and why.
cargo tree -e features -i openssl-sys   # who enables this crate?
cargo tree -d                          # crates present in more than one version

# 5. Human review of third-party code, recorded and shared: cargo-vet (Mozilla, Google).
cargo install cargo-vet && cargo vet init
```

A starting `deny.toml`:

```toml
[advisories]
version = 2
yanked = "deny"                      # yanked versions usually mean "do not use"

[licenses]
version = 2
allow = ["MIT", "Apache-2.0", "BSD-3-Clause", "ISC", "Unicode-3.0", "Zlib"]

[bans]
multiple-versions = "warn"
deny = [{ name = "openssl", reason = "use rustls" }]

[sources]
unknown-registry = "deny"
unknown-git = "deny"                 # no surprise git dependencies
```

### 55.2 Habits that reduce exposure

- **Turn off default features** you do not use:
  `reqwest = { version = "0.13", default-features = false, features = ["json", "rustls"] }`.
  Fewer features means fewer crates and less build-time code.
- **Watch for typosquatting** when you run `cargo add`. Check download
  counts, the repository link, and the owners on crates.io. Malicious
  crates imitating popular names have been published and removed several
  times.
- **Pin and verify CI tools** (`cargo install --locked --version x.y.z`),
  just like dependencies.
- **Sandbox builds** of untrusted code, in a container with no secrets and
  no network after `cargo fetch`.
- **Prefer pure-Rust crates** where they are mature (`rustls` over
  `openssl`), which removes a C toolchain and a C attack surface from the
  build.

The Go guide's equivalent is
[`govulncheck` in §64](../Golang/real-life-golang-guide.md#64-professional-workflow-linting-vuln-checks-ci-and-quality-gates).
One difference: `govulncheck` checks whether your code actually calls the
vulnerable function, while `cargo audit` reports any vulnerable crate
version in your lockfile. Expect more findings from `cargo audit`, and
triage them.

---

## 56. Project: an auth service with Argon2 and JWTs

**Implements:** [Go §48](../Golang/real-life-golang-guide.md#48-auth-service-password-hashing-and-jwt-issuing-verification)
in Rust, using
[Security from Zero Ch 46](../security/real-life-guide.md#chapter-46-authentication-sessions-and-federated-identity)
and [Security in Depth Ch 24](../security/real-life-security-guide-v1.md#chapter-24-token-exchange-and-delegation).

This project brings Part VI together: validated input types, Argon2id in
`spawn_blocking`, protection against user enumeration, strict JWT
validation, and an Axum extractor that makes authentication part of the
handler's signature.

```rust
// src/bin/authsvc.rs
use argon2::password_hash::{phc::PasswordHash, PasswordHasher, PasswordVerifier};
use argon2::Argon2;
use axum::extract::{FromRequestParts, State};
use axum::http::{request::Parts, StatusCode};
use axum::routing::{get, post};
use axum::{Json, Router};
use jsonwebtoken::{decode, encode, Algorithm, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::{Arc, LazyLock, RwLock};
use std::time::{SystemTime, UNIX_EPOCH};

const ISSUER: &str = "authsvc.rust-guide.local";
const AUDIENCE: &str = "api.rust-guide.local";
const TOKEN_TTL_SECS: u64 = 15 * 60;              // short-lived access tokens

#[derive(Clone)]
struct AppState { users: Arc<RwLock<HashMap<String, String>>>, jwt_secret: Arc<[u8]> }

#[derive(Deserialize)]
#[serde(deny_unknown_fields)]
struct Credentials { username: String, password: String }

#[derive(Serialize, Deserialize)]
struct Claims { sub: String, iss: String, aud: String, exp: u64, iat: u64 }

type ApiResult<T> = Result<T, (StatusCode, &'static str)>;

fn now() -> u64 { SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs() }

/// A real hash of a random password: verifying against it costs the same as a real user.
static DUMMY_HASH: LazyLock<String> = LazyLock::new(|| {
    Argon2::default().hash_password(&rand::random::<[u8; 16]>()).unwrap().to_string()
});

async fn register(State(st): State<AppState>, Json(c): Json<Credentials>) -> ApiResult<StatusCode> {
    let valid_name = (3..=32).contains(&c.username.len()) && c.username.bytes().all(|b| b.is_ascii_alphanumeric());
    if !valid_name { return Err((StatusCode::UNPROCESSABLE_ENTITY, "invalid username")); }
    if !(12..=128).contains(&c.password.len()) { return Err((StatusCode::UNPROCESSABLE_ENTITY, "password must be 12-128 chars")); }
    let hash = tokio::task::spawn_blocking(move || {        // ~tens of ms of CPU: keep it off the async workers
        Argon2::default().hash_password(c.password.as_bytes()).map(|h| h.to_string())
    }).await.unwrap().map_err(|_| (StatusCode::INTERNAL_SERVER_ERROR, "hash failed"))?;
    let mut users = st.users.write().unwrap();
    if users.contains_key(&c.username) { return Err((StatusCode::CONFLICT, "username taken")); }
    users.insert(c.username, hash);
    Ok(StatusCode::CREATED)
}

async fn login(State(st): State<AppState>, Json(c): Json<Credentials>) -> ApiResult<Json<serde_json::Value>> {
    let stored = st.users.read().unwrap().get(&c.username).cloned();
    let user_exists = stored.is_some();
    let hash = stored.unwrap_or_else(|| DUMMY_HASH.clone());
    // Always run a full Argon2 verification, so "no such user" and "wrong password" take the same time.
    let password_ok = tokio::task::spawn_blocking(move || {
        PasswordHash::new(&hash).is_ok_and(|h| Argon2::default().verify_password(c.password.as_bytes(), &h).is_ok())
    }).await.unwrap();
    if !(user_exists && password_ok) {
        return Err((StatusCode::UNAUTHORIZED, "invalid credentials"));  // one message for both cases
    }
    let iat = now();
    let claims = Claims { sub: c.username, iss: ISSUER.into(), aud: AUDIENCE.into(), iat, exp: iat + TOKEN_TTL_SECS };
    let token = encode(&Header::new(Algorithm::HS256), &claims, &EncodingKey::from_secret(&st.jwt_secret))
        .map_err(|_| (StatusCode::INTERNAL_SERVER_ERROR, "token error"))?;
    Ok(Json(serde_json::json!({ "access_token": token, "token_type": "Bearer", "expires_in": TOKEN_TTL_SECS })))
}

/// Extractor: a handler that takes `AuthUser` cannot run without a valid token.
struct AuthUser { name: String }

impl FromRequestParts<AppState> for AuthUser {
    type Rejection = (StatusCode, &'static str);
    async fn from_request_parts(parts: &mut Parts, st: &AppState) -> Result<Self, Self::Rejection> {
        let token = parts.headers.get("authorization").and_then(|v| v.to_str().ok())
            .and_then(|v| v.strip_prefix("Bearer "))
            .ok_or((StatusCode::UNAUTHORIZED, "missing bearer token"))?;
        let mut v = Validation::new(Algorithm::HS256);    // pin the algorithm: no "alg": "none", no RS/HS confusion
        v.set_issuer(&[ISSUER]);
        v.set_audience(&[AUDIENCE]);
        v.set_required_spec_claims(&["exp", "iss", "aud", "sub"]);
        v.leeway = 30;                                    // seconds of clock skew allowed
        let data = decode::<Claims>(token, &DecodingKey::from_secret(&st.jwt_secret), &v)
            .map_err(|_| (StatusCode::UNAUTHORIZED, "invalid token"))?;
        Ok(AuthUser { name: data.claims.sub })
    }
}

async fn me(user: AuthUser) -> String { format!("hello, {}", user.name) }

fn app() -> Router {
    let secret: [u8; 32] = rand::random();              // production: from a secret manager, rotated
    let state = AppState { users: Arc::default(), jwt_secret: Arc::from(secret.as_slice()) };
    Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .route("/me", get(me))
        .with_state(state)
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let listener = tokio::net::TcpListener::bind("127.0.0.1:3001").await?;
    println!("authsvc on {}", listener.local_addr()?);
    axum::serve(listener, app()).await
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::body::{to_bytes, Body};
    use axum::http::Request;
    use tower::ServiceExt;

    async fn call(app: &Router, method: &str, uri: &str, body: &str, bearer: Option<&str>) -> (StatusCode, String) {
        let mut b = Request::builder().method(method).uri(uri).header("content-type", "application/json");
        if let Some(t) = bearer { b = b.header("authorization", format!("Bearer {t}")); }
        let resp = app.clone().oneshot(b.body(Body::from(body.to_string())).unwrap()).await.unwrap();
        let status = resp.status();
        (status, String::from_utf8_lossy(&to_bytes(resp.into_body(), 1 << 20).await.unwrap()).into_owned())
    }

    #[tokio::test]
    async fn full_flow() {
        let app = app();
        let creds = r#"{"username":"ada","password":"correct horse battery"}"#;
        assert_eq!(call(&app, "POST", "/register", creds, None).await.0, StatusCode::CREATED);
        assert_eq!(call(&app, "POST", "/register", creds, None).await.0, StatusCode::CONFLICT);

        let (s, body) = call(&app, "POST", "/login", creds, None).await;
        assert_eq!(s, StatusCode::OK);
        let token = serde_json::from_str::<serde_json::Value>(&body).unwrap()["access_token"].as_str().unwrap().to_string();
        assert_eq!(call(&app, "GET", "/me", "", Some(&token)).await, (StatusCode::OK, "hello, ada".into()));

        // Same status and message for an unknown user and a wrong password: no user enumeration.
        let wrong_pw = call(&app, "POST", "/login", r#"{"username":"ada","password":"wrong password!!"}"#, None).await;
        let no_user = call(&app, "POST", "/login", r#"{"username":"bob","password":"whatever12345"}"#, None).await;
        assert_eq!(wrong_pw, no_user);

        // Tampered token: change one character of the signature.
        let mut bad = token.clone();
        let last = bad.pop().unwrap();
        bad.push(if last == 'A' { 'B' } else { 'A' });
        assert_eq!(call(&app, "GET", "/me", "", Some(&bad)).await.0, StatusCode::UNAUTHORIZED);
        assert_eq!(call(&app, "GET", "/me", "", None).await.0, StatusCode::UNAUTHORIZED);
    }
}
```

```bash
cargo test --bin authsvc
```

**What this does not do yet**, and where the guides cover it:

- **Refresh tokens and revocation.** A 15-minute access token cannot be
  revoked before it expires. Pair it with a server-side refresh token that
  can be.
- **Asymmetric signing.** With HS256, anyone who can verify a token can
  also mint one. If other services verify tokens, sign with Ed25519 or
  ES256 and publish a JWKS
  ([Security in Depth Ch 24](../security/real-life-security-guide-v1.md#chapter-24-token-exchange-and-delegation)).
- **Rate limiting** on `/login` per account and per IP, so that the
  expensive Argon2 hash is not itself a DoS lever
  ([Security in Depth Ch 40B](../security/real-life-security-guide-v1.md#chapter-40b-api-rate-limiting-in-depth)).
- **Persistent storage**, with a database and parameterized queries (§53.3).

---

# Part VII — Expert: internals and design

Parts I–VI covered how to use Rust. This Part covers how it works
underneath: what is in memory, what `async` compiles to, what the atomic
orderings promise, and how to design APIs that others can use safely. It
ends with two capstones that use the whole guide.

---

## 57. Memory layout: size, alignment, niches, `repr(C)`, fat pointers

```rust
use std::mem::{align_of, size_of};
use std::num::NonZeroU32;

#[allow(dead_code)]
struct Default_ { a: u8, b: u32, c: u16 }          // Rust may reorder fields

#[allow(dead_code)]
#[repr(C)]
struct CLayout { a: u8, b: u32, c: u16 }           // C order, C padding

#[allow(dead_code)]
enum Shape { Circle(f32), Rect(f32, f32) }

fn main() {
    println!("{:<28}{:>5}{:>7}", "type", "size", "align");
    macro_rules! row { ($t:ty) => { println!("{:<28}{:>5}{:>7}", stringify!($t), size_of::<$t>(), align_of::<$t>()) } }
    row!(Default_);            // 8: reordered to b, c, a + 1 byte of padding
    row!(CLayout);             // 12: a, 3 pad, b, c, 2 pad (same as the C guide §8.1)
    row!(&u8);                 // 8: thin pointer
    row!(&[u8]);               // 16: pointer + length
    row!(&str);                // 16
    row!(&dyn std::fmt::Debug);// 16: pointer + vtable
    row!(Box<u64>);            // 8
    row!(Option<Box<u64>>);    // 8: None uses the null niche
    row!(Option<u32>);         // 8: needs a separate tag
    row!(Option<NonZeroU32>);  // 4: 0 is the niche
    row!(Option<bool>);        // 1: values 2..=255 are niches
    row!(Result<u32, ()>);     // 8
    row!(Shape);               // 12: tag + the larger variant
    row!(String);              // 24
    row!(Vec<u8>);             // 24
    row!(());                  // 0: zero-sized type
}
```

```text
type                         size  align
Default_                        8      4
CLayout                        12      4
&u8                             8      8
&[u8]                          16      8
&str                           16      8
&dyn std::fmt::Debug           16      8
Box<u64>                        8      8
Option<Box<u64>>                8      8
Option<u32>                     8      4
Option<NonZeroU32>              4      4
Option<bool>                    1      1
Result<u32, ()>                 8      4
Shape                          12      4
String                         24      8
Vec<u8>                        24      8
()                              0      1
```

Three rules come out of this table:

1. **Rust's default layout is unspecified.** The compiler may reorder
   fields to minimize padding. `Default_` is 8 bytes, where the C layout of
   the same fields is 12. Use `#[repr(C)]` for anything that crosses an FFI
   boundary or is read from raw bytes, as in the
   [C guide §8.1](../c-lang/real-life-c-guide.md#8-structs-unions-enums-typedef-and-the-preprocessor).
   Never `transmute` a default-layout struct.
2. **Niches make `Option` free for pointer-like types.** References,
   `Box`, `NonNull`, and `NonZero*` have an invalid bit pattern (0), and
   `None` uses it. That is how `Option<&T>` has the same ABI as a nullable
   C pointer, so FFI code can use it.
3. **Zero-sized types cost nothing.** `HashSet<T>` is a
   `HashMap<T, ()>`, and the `()` values take no memory. Marker types in
   typestate APIs (Ch 62) take no space either.

---

## 58. Inside `Vec`, `String`, `Box`, and trait objects

### 58.1 Build a `Vec` to understand `Vec`

`Vec<T>` is a pointer, a capacity, and a length, plus a small amount of
`unsafe` code that keeps this invariant: **the first `len` slots are
initialized, and the slots from `len` to `cap` are allocated but
uninitialized.** Here is a working version:

```rust
use std::alloc::{self, Layout};
use std::ops::{Deref, DerefMut};
use std::ptr::{self, NonNull};

pub struct MiniVec<T> { ptr: NonNull<T>, cap: usize, len: usize }

// SAFETY: MiniVec owns its elements, exactly like Vec<T>, so it is Send/Sync when T is.
unsafe impl<T: Send> Send for MiniVec<T> {}
unsafe impl<T: Sync> Sync for MiniVec<T> {}

impl<T> MiniVec<T> {
    pub const fn new() -> Self {
        assert!(size_of::<T>() != 0, "zero-sized types are not supported in this sketch");
        Self { ptr: NonNull::dangling(), cap: 0, len: 0 }   // no allocation until the first push
    }

    fn grow(&mut self) {
        let new_cap = if self.cap == 0 { 4 } else { self.cap.checked_mul(2).expect("capacity overflow") };
        let new_layout = Layout::array::<T>(new_cap).expect("allocation too large");
        // SAFETY: new_layout has non-zero size (T is not a ZST and new_cap > 0). For realloc, the old
        // pointer was allocated by us with exactly Layout::array::<T>(self.cap).
        let raw = unsafe {
            if self.cap == 0 {
                alloc::alloc(new_layout)
            } else {
                alloc::realloc(self.ptr.as_ptr().cast(), Layout::array::<T>(self.cap).unwrap(), new_layout.size())
            }
        };
        self.ptr = NonNull::new(raw.cast()).unwrap_or_else(|| alloc::handle_alloc_error(new_layout));
        self.cap = new_cap;
    }

    pub fn push(&mut self, value: T) {
        if self.len == self.cap { self.grow(); }
        // SAFETY: len < cap, so the slot is allocated, and it is uninitialized, so write (not assign).
        unsafe { ptr::write(self.ptr.as_ptr().add(self.len), value) };
        self.len += 1;
    }

    pub fn pop(&mut self) -> Option<T> {
        if self.len == 0 { return None; }
        self.len -= 1;
        // SAFETY: slot `len` was initialized; after decrementing len we treat it as moved out.
        Some(unsafe { ptr::read(self.ptr.as_ptr().add(self.len)) })
    }

    pub fn capacity(&self) -> usize { self.cap }
}

impl<T> Deref for MiniVec<T> {
    type Target = [T];
    fn deref(&self) -> &[T] {
        // SAFETY: the first `len` elements are initialized; ptr is non-null and aligned (dangling is fine for len 0).
        unsafe { std::slice::from_raw_parts(self.ptr.as_ptr(), self.len) }
    }
}
impl<T> DerefMut for MiniVec<T> {
    fn deref_mut(&mut self) -> &mut [T] {
        // SAFETY: as above, and &mut self guarantees exclusive access.
        unsafe { std::slice::from_raw_parts_mut(self.ptr.as_ptr(), self.len) }
    }
}

impl<T> Drop for MiniVec<T> {
    fn drop(&mut self) {
        // SAFETY: drop exactly the initialized elements, then free the buffer with the layout we allocated.
        unsafe {
            ptr::drop_in_place(ptr::slice_from_raw_parts_mut(self.ptr.as_ptr(), self.len));
            if self.cap != 0 { alloc::dealloc(self.ptr.as_ptr().cast(), Layout::array::<T>(self.cap).unwrap()); }
        }
    }
}

fn main() {
    let mut v = MiniVec::new();
    for i in 0..10 { v.push(format!("item-{i}")); }   // heap-owning elements exercise Drop
    v[0].push_str("-edited");                          // DerefMut to [String]
    let last = v.pop();                                // pop first: printing v[0] and calling pop() in one
                                                       // println! would borrow v twice (E0502, Ch 7)
    println!("len={} cap={} first={} last={last:?}", v.len(), v.capacity(), v[0]);
    println!("{}", v.iter().filter(|s| s.ends_with('3')).count()); // slice methods come free via Deref
}
```

```bash
cargo +nightly miri run --bin minivec    # checks every unsafe access; it reports leaks, UB, and double drops
```

Each `unsafe` block has a `SAFETY` comment that refers to the invariant.
The public API (`push`, `pop`, indexing through `Deref`) cannot break that
invariant, so safe code using `MiniVec` cannot cause UB. That is the whole
technique behind the standard library's collections. The real `Vec` adds
zero-sized-type support, allocator parameters, and a lot of performance
work. The Rustonomicon builds it step by step.

### 58.2 `String` and `Box`

- `String` is a `Vec<u8>` plus the invariant "the bytes are valid UTF-8".
  Every method that could break that invariant either checks it or is
  `unsafe` (`from_utf8_unchecked`, `as_mut_vec`).
- `Box<T>` is a non-null pointer to a heap allocation that it owns. The
  compiler knows how to move a value out of it with `*boxed`, which no
  other type can do.

### 58.3 Trait objects vs Go interfaces

| | Rust `&dyn Trait` / `Box<dyn Trait>` | Go interface value |
|---|---|---|
| Size | 2 words: data pointer + vtable pointer | 2 words: type/itab pointer + data pointer |
| Where the vtable pointer lives | in the reference, not in the object | in the interface value |
| Small values | always behind a pointer | stored indirectly, usually on the heap |
| Nil trap | does not exist (a reference is never null) | an interface holding a nil pointer `!= nil` ([Go §9.4](../Golang/real-life-golang-guide.md#9-interfaces-implicit-satisfaction-any-type-switches)) |
| Multiple traits | `dyn A + B` only with auto traits; otherwise use a supertrait | any method set |

---

## 59. How `async` compiles: state machines, `Pin`, and a tiny executor

### 59.1 What the compiler generates

```rust
// fragment: what you write
async fn fetch_and_count(url: String) -> usize {
    let conn = connect(&url).await;
    let body = conn.read_all().await;
    body.len()
}

// fragment: roughly what the compiler generates (simplified)
enum FetchAndCount {
    Start { url: String },
    WaitingConnect { url: String, fut: ConnectFuture },
    WaitingRead { fut: ReadAllFuture },        // `url` is no longer live, so it is not stored
    Done,
}
impl Future for FetchAndCount {
    type Output = usize;
    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<usize> {
        // loop { match state { Start => create the connect future, move to WaitingConnect;
        //                      WaitingConnect => poll it; Pending => return Pending; Ready(conn) => ...
        //                      WaitingRead => poll it; Ready(body) => return Ready(body.len()) } }
    }
}
```

The size of the future is the size of its largest state. Variables that are
alive across an `.await` are stored in it. Everything else stays on the
stack during a `poll` call. This is why you can have a million tasks: each
one is exactly as large as the state it needs.

### 59.2 Why `Pin` exists

A state can hold a reference to another variable in the **same** state:

```rust
// fragment
async fn read_header(sock: &mut TcpStream) {
    let mut buf = [0u8; 1024];
    sock.read(&mut buf).await;   // the read future holds &mut buf, and buf lives in the same state machine
}
```

If the future were moved in memory after that borrow was created, the
reference would point at the old location. `Pin<&mut Self>` is a promise
that the value **will not move again** once it has been polled. That is
why `poll` takes `self: Pin<&mut Self>`, and why you see `Box::pin(fut)` or
`std::pin::pin!(fut)` when you store or poll futures by hand. Types that
are safe to move anyway, which is almost all of them, implement `Unpin`,
and for those `Pin` changes nothing.

### 59.3 A complete executor in 40 lines

```rust
use std::future::Future;
use std::pin::Pin;
use std::sync::atomic::{AtomicU32, Ordering};
use std::sync::{Arc, Mutex};
use std::task::{Context, Poll, Wake, Waker};
use std::thread;
use std::time::{Duration, Instant};

/// A waker that unparks the thread running block_on.
struct ThreadWaker(thread::Thread);
impl Wake for ThreadWaker {
    fn wake(self: Arc<Self>) { self.0.unpark(); }
}

fn block_on<F: Future>(fut: F) -> F::Output {
    let mut fut = std::pin::pin!(fut);                 // pin on the stack: it will never move again
    let waker = Waker::from(Arc::new(ThreadWaker(thread::current())));
    let mut cx = Context::from_waker(&waker);
    loop {
        match fut.as_mut().poll(&mut cx) {
            Poll::Ready(v) => return v,
            Poll::Pending => thread::park(),           // sleep until someone calls wake()
        }
    }
}

static POLLS: AtomicU32 = AtomicU32::new(0);

/// A timer future: the "reactor" is one helper thread per timer (Tokio uses one timer wheel instead).
struct Delay { deadline: Instant, waker: Option<Arc<Mutex<Waker>>> }

impl Future for Delay {
    type Output = &'static str;
    fn poll(mut self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<&'static str> {
        POLLS.fetch_add(1, Ordering::Relaxed);
        if Instant::now() >= self.deadline { return Poll::Ready("timer fired"); }
        match &self.waker {
            Some(w) => *w.lock().unwrap() = cx.waker().clone(), // always keep the LATEST waker
            None => {
                let w = Arc::new(Mutex::new(cx.waker().clone()));
                self.waker = Some(Arc::clone(&w));
                let deadline = self.deadline;
                thread::spawn(move || {
                    thread::sleep(deadline.saturating_duration_since(Instant::now()));
                    w.lock().unwrap().wake_by_ref();     // "the event happened: poll me again"
                });
            }
        }
        Poll::Pending
    }
}

fn main() {
    let t = Instant::now();
    let out = block_on(async {
        let a = Delay { deadline: Instant::now() + Duration::from_millis(50), waker: None }.await;
        let b = Delay { deadline: Instant::now() + Duration::from_millis(50), waker: None }.await;
        format!("{a}, {b}")
    });
    println!("{out} after {:?}, {} polls", t.elapsed(), POLLS.load(Ordering::Relaxed));
}
```

```text
timer fired, timer fired after 103.1ms, 4 polls
```

Four polls: each `Delay` is polled once (it returns `Pending` and arms the
timer) and once more after its wake-up. Nothing busy-waits. Tokio has the
same structure: a run queue of tasks, wakers that push a task back onto
the queue, and a reactor (epoll/kqueue plus a timer wheel) that calls the
wakers. It adds multiple worker threads with work stealing, which is the
same idea as Go's scheduler
([Go §28](../Golang/real-life-golang-guide.md#28-the-scheduler-deep-dive-gmp-and-work-stealing)).

---

## 60. Variance, `PhantomData`, and drop check

You need these when you write generic containers or `unsafe` code. You can
skip this chapter until you do.

**`PhantomData<T>`** tells the compiler that a type logically contains a
`T` when no field actually stores one. That affects three things: variance,
auto traits (`Send`/`Sync`), and drop checking.

```rust
use std::marker::PhantomData;

/// A typed ID: Id<User> and Id<Order> are different types with the same 8-byte representation.
struct Id<T> { raw: u64, _kind: PhantomData<fn() -> T> }
// `fn() -> T` instead of `T`: the Id does not OWN a T, so it is Send + Sync and Copy whatever T is.

// Derive would add `T: Clone` bounds we do not want, so implement by hand.
impl<T> Clone for Id<T> { fn clone(&self) -> Self { *self } }
impl<T> Copy for Id<T> {}
impl<T> PartialEq for Id<T> { fn eq(&self, o: &Self) -> bool { self.raw == o.raw } }
impl<T> Id<T> { fn new(raw: u64) -> Self { Id { raw, _kind: PhantomData } } }

struct User;
struct Order;

fn load_user(id: Id<User>) -> String { format!("user #{}", id.raw) }

fn main() {
    let u: Id<User> = Id::new(7);
    let o: Id<Order> = Id::new(7);
    println!("{}", load_user(u));
    // load_user(o);                         // error[E0308]: expected `Id<User>`, found `Id<Order>`
    println!("{} bytes, same raw value: {}", size_of::<Id<User>>(), u.raw == o.raw);
}
```

**Variance** decides whether `Foo<&'long T>` can be used where
`Foo<&'short T>` is expected. `&'a T` is covariant (a longer-lived
reference can stand in for a shorter one). `&'a mut T` is **invariant** in
`T`. If it were not, you could write a short-lived reference into a slot
that expects a long-lived one, and that is a use-after-free.
`PhantomData<T>` makes your type covariant in `T`, `PhantomData<fn(T)>`
makes it contravariant, and `PhantomData<*mut T>` or `PhantomData<Cell<T>>`
makes it invariant. The Rustonomicon's chapter on subtyping has the full
table.

**Drop check** stops a `Drop` implementation from reading data that has
already been dropped. If your `Drop` impl for `Container<T>` might use a
`T`, `PhantomData<T>` tells the compiler that dropping a `Container<T>`
drops `T` values, so borrows inside `T` must still be alive at that point.

---

## 61. Memory ordering: `Relaxed`, `Acquire`/`Release`, `SeqCst`, and a spinlock

CPUs and compilers reorder memory operations. An atomic ordering states
which reorderings are not allowed. Rust uses the C++20/C11 model, so the
[C guide §38](../c-lang/real-life-c-guide.md#38-the-c11-memory-model-atomics-ordering-lock-free-basics)
applies one-to-one.

| Ordering | Guarantees | Use for |
|---|---|---|
| `Relaxed` | atomicity only; no ordering with other memory | counters, statistics |
| `Release` (store) / `Acquire` (load) | everything written before the Release store is visible to the thread whose Acquire load reads that store | publishing data, locks |
| `AcqRel` | both, for read-modify-write operations | `compare_exchange` in lock-free structures |
| `SeqCst` | Acquire/Release plus one global order of all `SeqCst` operations | when you cannot prove a weaker ordering is enough |

### 61.1 Publishing data with Release/Acquire

```rust
use std::sync::atomic::{AtomicBool, AtomicU64, Ordering};
use std::thread;

static DATA: AtomicU64 = AtomicU64::new(0);
static READY: AtomicBool = AtomicBool::new(false);

fn main() {
    let producer = thread::spawn(|| {
        DATA.store(42, Ordering::Relaxed);          // (1) write the payload
        READY.store(true, Ordering::Release);       // (2) publish: (1) cannot move after this
    });
    let consumer = thread::spawn(|| {
        while !READY.load(Ordering::Acquire) {      // (3) once this sees true...
            std::hint::spin_loop();
        }
        DATA.load(Ordering::Relaxed)                // (4) ...this is guaranteed to see 42
    });
    producer.join().unwrap();
    println!("consumer saw {}", consumer.join().unwrap());
}
```

With `Relaxed` on both `READY` operations, a weakly ordered CPU such as an
ARM server or an Apple Silicon Mac may legally let the consumer see `READY
== true` and `DATA == 0`. On x86 the hardware hides most of these
reorderings, so bugs like this can pass every test on a developer laptop
and then fail on Graviton. Test lock-free code with
[`loom`](https://crates.io/crates/loom), which explores the possible
interleavings and orderings systematically.

### 61.2 A spinlock with the right orderings

```rust
use std::cell::UnsafeCell;
use std::ops::{Deref, DerefMut};
use std::sync::atomic::{AtomicBool, Ordering};

pub struct SpinLock<T> { locked: AtomicBool, value: UnsafeCell<T> }

// SAFETY: the lock gives one thread at a time access to `value`, so sharing &SpinLock across
// threads is sound whenever T itself may be sent between threads.
unsafe impl<T: Send> Sync for SpinLock<T> {}

pub struct Guard<'a, T> { lock: &'a SpinLock<T> }

impl<T> SpinLock<T> {
    pub const fn new(value: T) -> Self { Self { locked: AtomicBool::new(false), value: UnsafeCell::new(value) } }

    pub fn lock(&self) -> Guard<'_, T> {
        // Acquire on success: we must see every write made by the previous holder before its Release.
        while self.locked.compare_exchange_weak(false, true, Ordering::Acquire, Ordering::Relaxed).is_err() {
            while self.locked.load(Ordering::Relaxed) { std::hint::spin_loop(); } // spin on reads, not on CAS
        }
        Guard { lock: self }
    }
}

impl<T> Deref for Guard<'_, T> {
    type Target = T;
    // SAFETY: holding the Guard means we hold the lock: no other reference to `value` exists.
    fn deref(&self) -> &T { unsafe { &*self.lock.value.get() } }
}
impl<T> DerefMut for Guard<'_, T> {
    fn deref_mut(&mut self) -> &mut T { unsafe { &mut *self.lock.value.get() } }
}
impl<T> Drop for Guard<'_, T> {
    fn drop(&mut self) { self.lock.locked.store(false, Ordering::Release); } // publish our writes
}

fn main() {
    static COUNTER: SpinLock<u64> = SpinLock::new(0);
    std::thread::scope(|s| {
        for _ in 0..8 {
            s.spawn(|| for _ in 0..100_000 { *COUNTER.lock() += 1; });
        }
    });
    println!("{}", *COUNTER.lock()); // 800000, every time
}
```

The guard is what makes this a Rust API rather than a C API. You cannot
reach the data without locking, and you cannot forget to unlock. Do not use
spinlocks in user space in practice. A preempted lock holder makes every
waiter burn a full time slice. `std::sync::Mutex` spins briefly and then
sleeps on a futex, which is what you want
([OS Ch 18](../os-linux/real-life-os-guide.md#chapter-18-locks-semaphores-and-deadlock)).
Spinlocks belong in kernels and interrupt handlers
([C §52](../c-lang/real-life-c-guide.md#52-interrupts-and-kernel-concurrency-top-bottom-halves-spinlocks-vs-mutexes)).

---

## 62. API design: newtypes, typestate, builders, errors, and semver

The Go guide's [§54 "Go proverbs, applied"](../Golang/real-life-golang-guide.md#54-designing-apis-the-go-proverbs-applied)
has a Rust counterpart: the [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/).
The techniques below are the ones that most often separate a good crate
from an awkward one.

### 62.1 Typestate: invalid call orders do not compile

```rust
use std::marker::PhantomData;

struct Disconnected;
struct Connected;

struct Client<State> { addr: String, session: Option<u32>, _state: PhantomData<State> }

impl Client<Disconnected> {
    fn new(addr: &str) -> Self { Client { addr: addr.into(), session: None, _state: PhantomData } }
    fn connect(self) -> Result<Client<Connected>, String> {     // consumes the disconnected client
        if self.addr.is_empty() { return Err("no address".into()); }
        Ok(Client { addr: self.addr, session: Some(7), _state: PhantomData })
    }
}

impl Client<Connected> {
    fn send(&mut self, msg: &str) -> usize { println!("[{}#{:?}] {msg}", self.addr, self.session); msg.len() }
    fn close(self) -> Client<Disconnected> { Client { addr: self.addr, session: None, _state: PhantomData } }
}

fn main() -> Result<(), String> {
    let mut c = Client::new("db:5432").connect()?;
    c.send("SELECT 1");
    let c = c.close();
    // c.send("again");   // error[E0599]: no method named `send` found for `Client<Disconnected>`
    let _ = c;
    Ok(())
}
```

The states are zero-sized (Ch 57), so the type checking costs nothing at
runtime. `reqwest::RequestBuilder`, the `hyper` connection builder, and
embedded HAL crates (a GPIO pin configured as input has no `set_high`
method) all use this pattern.

### 62.2 Builders for many optional settings

```rust
use std::time::Duration;

#[derive(Debug)]
pub struct ServerConfig { addr: String, max_body: usize, timeout: Duration, tls: bool }

#[derive(Default)]
pub struct ServerConfigBuilder { addr: Option<String>, max_body: Option<usize>, timeout: Option<Duration>, tls: bool }

impl ServerConfigBuilder {
    pub fn addr(mut self, a: impl Into<String>) -> Self { self.addr = Some(a.into()); self }
    pub fn max_body(mut self, n: usize) -> Self { self.max_body = Some(n); self }
    pub fn timeout(mut self, t: Duration) -> Self { self.timeout = Some(t); self }
    pub fn tls(mut self, on: bool) -> Self { self.tls = on; self }
    pub fn build(self) -> Result<ServerConfig, String> {
        let max_body = self.max_body.unwrap_or(1 << 20);
        if max_body > 64 << 20 { return Err("max_body above 64 MiB".into()); } // validate once, here
        Ok(ServerConfig {
            addr: self.addr.ok_or("addr is required")?,
            max_body,
            timeout: self.timeout.unwrap_or(Duration::from_secs(30)),
            tls: self.tls,
        })
    }
}

fn main() {
    let cfg = ServerConfigBuilder::default().addr("0.0.0.0:443").tls(true).build();
    println!("{cfg:?}");
    println!("{:?}", ServerConfigBuilder::default().max_body(1 << 30).build());
}
```

This plays the same role as Go's functional options
([Go §55](../Golang/real-life-golang-guide.md#55-dependency-injection-without-a-framework-functional-options)).
If you combine it with typestate (a `Builder<NoAddr>` without a `build`
method), "addr is required" becomes a compile error instead of a runtime
one.

### 62.3 Design checklist for a public crate

- **Accept general types, return concrete ones:** `impl AsRef<Path>`,
  `impl Into<String>`, and `&[T]` as inputs, concrete types as outputs (the
  same advice as Go's "accept interfaces, return structs").
- **Errors:** one error enum per module or crate, implementing
  `std::error::Error` (`thiserror`), and `#[non_exhaustive]` so that adding
  a variant is not a breaking change.
- **`#[non_exhaustive]` on public structs** whose fields may grow, and
  private fields plus constructors wherever there is an invariant.
- **Implement the common traits** (`Debug`, `Clone`, `PartialEq`,
  `Default`, `Send`/`Sync` where possible). Users cannot add them to your
  types later, because of the orphan rule.
- **Sealed traits** (a public trait with a supertrait in a private module)
  when you want users to call a trait's methods but not implement it, so
  you can add methods later.
- **Semver:** `cargo semver-checks` detects breaking changes before you
  publish. Making a type no longer `Send` is a breaking change, and so is
  adding a variant to an exhaustive public enum.
- **Documentation:** a runnable example on every public item (doc tests,
  §23.2), and `#![warn(missing_docs)]`.

---

## 63. Capstone: a Redis-compatible key-value server with persistence

This is the capstone of the [C guide (§43)](../c-lang/real-life-c-guide.md#43-capstone-a-mini-redis-like-key-value-store-with-a-wire-protocol)
and close to the [Go guide's (§56)](../Golang/real-life-golang-guide.md#56-capstone-a-tiny-embedded-database-engine-bitcask-style),
written in Rust. It speaks enough of the RESP2 protocol for `redis-cli` to
work, and it combines most of the guide: Tokio tasks (31), framing with
limits (32, 53), enums and `match` (9), a std `Mutex` never held across
`.await` (31.5), expiry, atomic snapshots (41), and graceful shutdown (31.4).

```rust
// src/bin/kv.rs     run: cargo run --bin kv -- serve     then: redis-cli -p 6380 SET greeting hello
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::io;
use std::path::Path;
use std::sync::{Arc, Mutex};
use std::time::{Duration, SystemTime, UNIX_EPOCH};
use tokio::io::{AsyncBufReadExt, AsyncReadExt, AsyncWriteExt, BufReader};
use tokio::net::{TcpListener, TcpStream};

const MAX_ARGS: usize = 1024;
const MAX_BULK: usize = 8 * 1024 * 1024;             // 8 MiB per argument
const MAX_LINE: usize = 64;                          // "*123\r\n" / "$456\r\n" headers are short

#[derive(Clone, Serialize, Deserialize)]
struct Entry { value: Vec<u8>, expires_at_ms: Option<u64> }

type Db = Arc<Mutex<HashMap<Vec<u8>, Entry>>>;

fn now_ms() -> u64 { SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_millis() as u64 }

enum Reply { Ok, Simple(&'static str), Bulk(Option<Vec<u8>>), Int(i64), Err(String) }

impl Reply {
    fn encode(self) -> Vec<u8> {
        match self {
            Reply::Ok => b"+OK\r\n".to_vec(),
            Reply::Simple(s) => format!("+{s}\r\n").into_bytes(),
            Reply::Bulk(None) => b"$-1\r\n".to_vec(),
            Reply::Bulk(Some(v)) => [format!("${}\r\n", v.len()).as_bytes(), &v, b"\r\n"].concat(),
            Reply::Int(n) => format!(":{n}\r\n").into_bytes(),
            Reply::Err(e) => format!("-ERR {e}\r\n").into_bytes(),
        }
    }
}

/// Reads one header line ("*3", "$5"), bounded so a client cannot send an endless line.
async fn read_header(r: &mut BufReader<TcpStream>, prefix: u8) -> io::Result<Option<usize>> {
    let mut line = Vec::new();
    let n = (&mut *r).take(MAX_LINE as u64).read_until(b'\n', &mut line).await?;
    if n == 0 { return Ok(None); }                                    // clean EOF between commands
    let bad = || io::Error::new(io::ErrorKind::InvalidData, "protocol error");
    if !line.ends_with(b"\r\n") || line.first() != Some(&prefix) { return Err(bad()); }
    let num = std::str::from_utf8(&line[1..line.len() - 2]).map_err(|_| bad())?;
    num.parse::<usize>().map(Some).map_err(|_| bad())
}

async fn read_command(r: &mut BufReader<TcpStream>) -> io::Result<Option<Vec<Vec<u8>>>> {
    let Some(argc) = read_header(r, b'*').await? else { return Ok(None) };
    if argc == 0 || argc > MAX_ARGS { return Err(io::Error::new(io::ErrorKind::InvalidData, "bad arg count")); }
    let mut args = Vec::with_capacity(argc);                       // argc is bounded, so this is safe
    for _ in 0..argc {
        let len = read_header(r, b'$').await?.ok_or(io::ErrorKind::UnexpectedEof)?;
        if len > MAX_BULK { return Err(io::Error::new(io::ErrorKind::InvalidData, "argument too large")); }
        let mut buf = vec![0u8; len + 2];
        r.read_exact(&mut buf).await?;
        if !buf.ends_with(b"\r\n") { return Err(io::Error::new(io::ErrorKind::InvalidData, "missing CRLF")); }
        buf.truncate(len);
        args.push(buf);
    }
    Ok(Some(args))
}

fn execute(db: &Db, args: &[Vec<u8>]) -> Reply {
    let cmd = String::from_utf8_lossy(&args[0]).to_ascii_uppercase();
    let mut map = db.lock().unwrap();                                // synchronous critical section
    let live = |e: &Entry| e.expires_at_ms.is_none_or(|t| t > now_ms());
    match (cmd.as_str(), args.len()) {
        ("PING", 1) => Reply::Simple("PONG"),
        ("PING", 2) | ("ECHO", 2) => Reply::Bulk(Some(args[1].clone())),
        ("GET", 2) => Reply::Bulk(map.get(&args[1]).filter(|e| live(e)).map(|e| e.value.clone())),
        ("SET", 3) | ("SET", 5) => {
            let expires_at_ms = if args.len() == 5 {
                let unit = String::from_utf8_lossy(&args[3]).to_ascii_uppercase();
                let Some(n) = std::str::from_utf8(&args[4]).ok().and_then(|s| s.parse::<u64>().ok()).filter(|&n| n > 0) else {
                    return Reply::Err("invalid expire time".into());
                };
                match unit.as_str() {
                    "EX" => Some(now_ms().saturating_add(n.saturating_mul(1000))),
                    "PX" => Some(now_ms().saturating_add(n)),
                    _ => return Reply::Err("syntax error".into()),
                }
            } else { None };
            map.insert(args[1].clone(), Entry { value: args[2].clone(), expires_at_ms });
            Reply::Ok
        }
        ("DEL", n) if n >= 2 => Reply::Int(args[1..].iter().filter(|k| map.remove(*k).is_some()).count() as i64),
        ("EXISTS", n) if n >= 2 => Reply::Int(args[1..].iter().filter(|k| map.get(*k).is_some_and(|e| live(e))).count() as i64),
        ("INCR", 2) => {
            let current = map.get(&args[1]).filter(|e| live(e)).map(|e| e.value.clone());
            let n = match current.as_deref().map(std::str::from_utf8) {
                None => 0,
                Some(Ok(s)) => match s.parse::<i64>() { Ok(n) => n, Err(_) => return Reply::Err("value is not an integer".into()) },
                Some(Err(_)) => return Reply::Err("value is not an integer".into()),
            };
            let Some(next) = n.checked_add(1) else { return Reply::Err("increment would overflow".into()) };
            map.insert(args[1].clone(), Entry { value: next.to_string().into_bytes(), expires_at_ms: None });
            Reply::Int(next)
        }
        ("DBSIZE", 1) => Reply::Int(map.values().filter(|e| live(e)).count() as i64),
        ("SAVE", 1) => {
            let snapshot: Vec<(Vec<u8>, Entry)> = map.iter().filter(|(_, e)| live(e)).map(|(k, e)| (k.clone(), e.clone())).collect();
            drop(map);                                                    // do not hold the lock during disk I/O
            match serde_json::to_vec(&snapshot).map_err(io::Error::other).and_then(|b| atomic_write(Path::new(SNAPSHOT), &b)) {
                Ok(()) => Reply::Ok,
                Err(e) => Reply::Err(format!("save failed: {e}")),
            }
        }
        _ => Reply::Err(format!("unknown command or wrong number of arguments for '{cmd}'")),
    }
}

const SNAPSHOT: &str = "kv-snapshot.json";

fn atomic_write(path: &Path, data: &[u8]) -> io::Result<()> {        // Ch 41, condensed
    let tmp = path.with_extension("tmp");
    { let mut f = std::fs::File::create(&tmp)?; io::Write::write_all(&mut f, data)?; f.sync_all()?; }
    std::fs::rename(&tmp, path)?;
    std::fs::File::open(path.parent().filter(|p| !p.as_os_str().is_empty()).unwrap_or(Path::new(".")))?.sync_all()
}

fn load(path: &Path) -> HashMap<Vec<u8>, Entry> {
    let Ok(bytes) = std::fs::read(path) else { return HashMap::new() };
    serde_json::from_slice::<Vec<(Vec<u8>, Entry)>>(&bytes).map(|v| v.into_iter().collect()).unwrap_or_default()
}

async fn client(stream: TcpStream, db: Db) -> io::Result<()> {
    let mut r = BufReader::new(stream);
    while let Some(args) = read_command(&mut r).await? {
        let reply = execute(&db, &args).encode();                     // lock taken and released inside
        r.get_mut().write_all(&reply).await?;
    }
    Ok(())
}

async fn serve(listener: TcpListener, db: Db) {
    let sweeper_db = Arc::clone(&db);
    tokio::spawn(async move {                                         // active expiry, once per second
        let mut tick = tokio::time::interval(Duration::from_secs(1));
        loop { tick.tick().await; let t = now_ms(); sweeper_db.lock().unwrap().retain(|_, e| e.expires_at_ms.is_none_or(|x| x > t)); }
    });
    loop {
        let Ok((stream, peer)) = listener.accept().await else { continue };
        let db = Arc::clone(&db);
        tokio::spawn(async move {
            if let Err(e) = client(stream, db).await {
                if e.kind() != io::ErrorKind::UnexpectedEof { eprintln!("{peer}: {e}"); }
            }
        });
    }
}

async fn demo(addr: std::net::SocketAddr) -> io::Result<()> {
    let mut s = TcpStream::connect(addr).await?;
    let cmd = |parts: &[&str]| -> Vec<u8> {
        let mut out = format!("*{}\r\n", parts.len()).into_bytes();
        for p in parts { out.extend(format!("${}\r\n{p}\r\n", p.len()).into_bytes()); }
        out
    };
    for parts in [&["PING"][..], &["SET", "greeting", "hello"], &["GET", "greeting"], &["INCR", "hits"], &["INCR", "hits"],
                  &["SET", "session", "abc", "PX", "100"], &["EXISTS", "greeting", "session", "nope"], &["INCR", "greeting"], &["SAVE"]] {
        s.write_all(&cmd(parts)).await?;
        let mut buf = vec![0u8; 512];
        let n = s.read(&mut buf).await?;
        println!("{:<40} -> {:?}", parts.join(" "), String::from_utf8_lossy(&buf[..n]));
    }
    tokio::time::sleep(Duration::from_millis(150)).await;            // let `session` expire
    s.write_all(&cmd(&["GET", "session"])).await?;
    let mut buf = vec![0u8; 64];
    let n = s.read(&mut buf).await?;
    println!("{:<40} -> {:?}", "GET session (after 150ms)", String::from_utf8_lossy(&buf[..n]));
    s.write_all(b"*99999999\r\n").await?;                             // hostile header: server drops the connection
    println!("{:<40} -> closed: {}", "*99999999", s.read(&mut buf).await? == 0);
    Ok(())
}

#[tokio::main]
async fn main() -> io::Result<()> {
    let db: Db = Arc::new(Mutex::new(load(Path::new(SNAPSHOT))));
    if std::env::args().nth(1).as_deref() == Some("serve") {
        let listener = TcpListener::bind("127.0.0.1:6380").await?;
        println!("kv listening on {}, {} keys loaded", listener.local_addr()?, db.lock().unwrap().len());
        tokio::select! {
            _ = serve(listener, Arc::clone(&db)) => {}
            _ = tokio::signal::ctrl_c() => println!("shutting down"),
        }
        return Ok(());
    }
    let listener = TcpListener::bind("127.0.0.1:0").await?;
    let addr = listener.local_addr()?;
    tokio::spawn(serve(listener, db));
    demo(addr).await?;
    std::fs::remove_file(SNAPSHOT)
}
```

```text
PING                                     -> "+PONG\r\n"
SET greeting hello                       -> "+OK\r\n"
GET greeting                             -> "$5\r\nhello\r\n"
INCR hits                                -> ":1\r\n"
INCR hits                                -> ":2\r\n"
SET session abc PX 100                   -> "+OK\r\n"
EXISTS greeting session nope             -> ":2\r\n"
INCR greeting                            -> "-ERR value is not an integer\r\n"
SAVE                                     -> "+OK\r\n"
GET session (after 150ms)                -> "$-1\r\n"
127.0.0.1:55288: bad arg count            <- the server's stderr
*99999999                                -> closed: true
```

**Compare it with the C version.** The C capstone needs a hand-written hash
table, manual buffer management with a growing input buffer, careful
`realloc` error paths, and `select()` bookkeeping. Here the memory-safety
work is gone. The remaining code is protocol logic, limits, and
concurrency design, which is where the remaining bugs will be.

**Stretch goals:**

1. **Append-only log** (Redis AOF, Bitcask in the Go guide): write each
   mutating command to a file and replay it on startup. Compare `fsync`
   per command with once per second, the durability vs throughput
   trade-off from [OS Ch 77](../os-linux/real-life-os-guide.md#chapter-77-files-that-survive-crashes-page-cache-fsync-and-atomic-replacement).
2. **Sharding:** replace `Mutex<HashMap>` with 16 shards selected by key
   hash and measure the throughput with `redis-benchmark -p 6380 -t set,get -c 50`.
3. **Pipelining:** the server already supports it, because it reads
   commands in a loop. Batch the replies into one `write` per read buffer
   and measure again.
4. **TLS:** wrap the listener with the Ch 49 acceptor. `redis-cli --tls`
   will connect.

---

## 64. Capstone: one HTTPS request, every layer, every guide, in Rust

The [HTTPS guide's final chapter](../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)
traces one request through the whole series in Go. This is the Rust
version: a `curl -w`-style tool that performs each phase **by hand**, times
it, and prints which guide explains it.

```rust
// src/bin/lifecycle.rs     usage: lifecycle https://example.com/
use rustls::pki_types::ServerName;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::io::{AsyncBufReadExt, AsyncWriteExt, BufReader};
use tokio::net::TcpStream;
use tokio::time::timeout;

fn phase(name: &str, took: Duration, detail: &str, guide: &str) {
    println!("{name:<10} {:>9.2?}  {detail:<55} {guide}", took);
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let raw = std::env::args().nth(1).unwrap_or_else(|| "https://example.com/".into());
    let url = reqwest::Url::parse(&raw)?;
    anyhow::ensure!(url.scheme() == "https", "https URLs only");
    let host = url.host_str().ok_or_else(|| anyhow::anyhow!("no host"))?.to_string();
    let port = url.port_or_known_default().unwrap_or(443);
    let path = match url.query() { Some(q) => format!("{}?{q}", url.path()), None => url.path().to_string() };
    let total = Instant::now();

    // 1. DNS: name -> addresses (Net Ch 18, HTTPS Ch 2). The OS resolver reads /etc/hosts first (Ch 15).
    let t = Instant::now();
    let addrs: Vec<_> = timeout(Duration::from_secs(5), tokio::net::lookup_host((host.as_str(), port))).await??.collect();
    let addr = *addrs.iter().find(|a| a.is_ipv4()).or(addrs.first()).ok_or_else(|| anyhow::anyhow!("no addresses"))?;
    phase("dns", t.elapsed(), &format!("{host} -> {addr} ({} addresses)", addrs.len()), "Net Ch 18 · HTTPS Ch 2 · Rust Ch 45");

    // 2. TCP: three-way handshake = one round trip (Net Ch 21, HTTPS Ch 3). Measures RTT to the edge.
    let t = Instant::now();
    let tcp = timeout(Duration::from_secs(5), TcpStream::connect(addr)).await??;
    tcp.set_nodelay(true)?;
    phase("tcp", t.elapsed(), &format!("{} -> {}", tcp.local_addr()?, tcp.peer_addr()?), "Net Ch 21 · OS Ch 19 · Rust Ch 32");

    // 3. TLS 1.3: one more round trip; certificate chain validated against the Mozilla root set.
    let t = Instant::now();
    let mut roots = rustls::RootCertStore::empty();
    roots.extend(webpki_roots::TLS_SERVER_ROOTS.iter().cloned());
    let mut cfg = rustls::ClientConfig::builder().with_root_certificates(roots).with_no_client_auth();
    cfg.alpn_protocols = vec![b"http/1.1".to_vec()];             // we speak HTTP/1.1 by hand below
    let connector = tokio_rustls::TlsConnector::from(Arc::new(cfg));
    let mut tls = timeout(Duration::from_secs(5), connector.connect(ServerName::try_from(host.clone())?, tcp)).await??;
    let (_, conn) = tls.get_ref();
    let detail = format!("{:?}, {:?}, ALPN {:?}, chain of {}",
        conn.protocol_version().unwrap(), conn.negotiated_cipher_suite().unwrap().suite(),
        conn.alpn_protocol().map(String::from_utf8_lossy), conn.peer_certificates().map_or(0, |c| c.len()));
    phase("tls", t.elapsed(), &detail, "Sec Ch 28-29 · HTTPS Ch 4 · Rust Ch 49");

    // 4. HTTP/1.1 request; time to first byte = network RTT + server (and CDN/proxy) processing.
    let t = Instant::now();
    let req = format!("GET {path} HTTP/1.1\r\nHost: {host}\r\nUser-Agent: rust-guide-lifecycle/1.0\r\nAccept: */*\r\nConnection: close\r\n\r\n");
    tls.write_all(req.as_bytes()).await?;
    let mut reader = BufReader::new(tls);
    let mut status = String::new();
    timeout(Duration::from_secs(10), reader.read_line(&mut status)).await??;
    phase("ttfb", t.elapsed(), status.trim(), "Net Ch 28 · HTTPS Ch 5 · Rust Ch 46");

    // 5. Response headers: what the infrastructure in between left behind.
    let mut interesting = Vec::new();
    loop {
        let mut line = String::new();
        if reader.read_line(&mut line).await? == 0 || line == "\r\n" { break; }
        let lower = line.to_ascii_lowercase();
        for h in ["server:", "via:", "age:", "cache-control:", "cf-cache-status:", "x-cache:", "alt-svc:", "strict-transport-security:"] {
            if lower.starts_with(h) { interesting.push(line.trim().to_string()); }
        }
    }
    for h in &interesting { println!("{:<21}{h}", ""); }
    println!("{:<21}(proxies, CDN cache: HTTPS Ch 8, 10, 15 · HSTS: HTTPS Ch 12 · HTTP/3 via Alt-Svc: HTTPS Ch 14)", "");

    // 6. Body: read to EOF (Connection: close), bounded.
    let t = Instant::now();
    let mut body = Vec::new();
    let mut limited = tokio::io::AsyncReadExt::take(reader, 10 << 20);
    tokio::io::AsyncReadExt::read_to_end(&mut limited, &mut body).await?;
    phase("body", t.elapsed(), &format!("{} bytes", body.len()), "OS Ch 14 (bytes land in a socket buffer, then your Vec)");
    phase("total", total.elapsed(), "", "");
    Ok(())
}
```

```bash
cargo run --release --bin lifecycle -- https://example.com/
```

```text
dns          77.95ms  example.com -> 172.66.147.243:443 (2 addresses)         Net Ch 18 · HTTPS Ch 2 · Rust Ch 45
tcp          14.24ms  192.168.1.4:55293 -> 172.66.147.243:443                 Net Ch 21 · OS Ch 19 · Rust Ch 32
tls          62.97ms  TLSv1_3, TLS13_AES_256_GCM_SHA384, ALPN Some("http/1.1"), chain of 4 Sec Ch 28-29 · HTTPS Ch 4 · Rust Ch 49
ttfb         18.36ms  HTTP/1.1 200 OK                                         Net Ch 28 · HTTPS Ch 5 · Rust Ch 46
                     Server: cloudflare
                     Age: 12297
                     cf-cache-status: HIT
                     alt-svc: h3=":443"; ma=86400
                     (proxies, CDN cache: HTTPS Ch 8, 10, 15 · HSTS: HTTPS Ch 12 · HTTP/3 via Alt-Svc: HTTPS Ch 14)
body          9.04µs  589 bytes                                               OS Ch 14 (bytes land in a socket buffer, then your Vec)
total       173.64ms
```

**Reading your own numbers.** The run above was a debug build on a home
connection, and yours will differ. Here, `cf-cache-status: HIT` and
`Age: 12297` mean a Cloudflare edge served the page from cache, so `ttfb`
(18 ms) is about one round trip, the same as `tcp` (14 ms). The `alt-svc`
header advertises HTTP/3 for the next visit.

- **`tcp` ≈ one RTT, and `tls` ≈ one more RTT plus crypto.** Here `tls`
  is about four times `tcp`. Part of that is the unoptimized debug build
  doing the signature checks, so rerun with `--release` before blaming
  the network. If it stays high, look at certificate chain size or a TLS
  terminator far from the edge.
- **`ttfb` minus one RTT is server time.** On a CDN hit, as above, it is
  close to zero. On a miss it includes the trip to the origin
  ([HTTPS Ch 10](../v2-https/real-life-guide-v1.md#chapter-10-cdn-content-near-you)).
- **The `dns` phase** is often near zero on the second run because the OS
  or the resolver cached it ([Net Ch 18](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses)).
- **Debug a slow phase** with the tools from
  [HTTPS Ch 17](../v2-https/real-life-guide-v1.md#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs):
  `dig` for dns, `tcpdump`/`ss -ti` for tcp
  ([Net Ch 31, 33](../networking/tcp-ip/real-life-guide-v1.md#chapter-31-tcpdump-watching-packets-go-by)),
  and `openssl s_client` for tls.

Every phase in this program is one of the guides. DNS and TCP come from the
networking guide. The socket and the bytes in the kernel buffer come from
the OS guide. Certificate validation and AEAD record encryption come from
the security guides. Status codes, headers, and caches come from the HTTPS
guide. The program is about 100 lines because Rust's libraries handle each
layer correctly and the types keep the layers apart.

---

# Part VIII — Professional Rust

This Part covers what happens once the code exists: finding bugs in
running programs, structuring a growing codebase, putting quality gates in
CI, and running Rust services in containers. It ends with the anti-patterns
that code review catches most often.

---

## 65. Debugging: backtraces, `lldb`/`gdb`, `tokio-console`, Miri

### 65.1 Panics and backtraces

```bash
RUST_BACKTRACE=1 cargo run      # backtrace on panic
RUST_BACKTRACE=full cargo run   # with every frame, including std internals
RUST_LIB_BACKTRACE=1            # anyhow: capture a backtrace when an error is created
```

In production, install a panic hook that sends the panic to your logs, so
it is not lost on a stderr nobody reads:

```rust
fn main() {
    let default_hook = std::panic::take_hook();
    std::panic::set_hook(Box::new(move |info| {
        let location = info.location().map(|l| format!("{}:{}", l.file(), l.line())).unwrap_or_default();
        let msg = info.payload().downcast_ref::<&str>().copied()
            .or_else(|| info.payload().downcast_ref::<String>().map(String::as_str))
            .unwrap_or("<non-string panic>");
        // In a service: tracing::error!(%location, msg, "panic") so it reaches your log pipeline.
        eprintln!(r#"{{"level":"error","event":"panic","location":"{location}","message":"{msg}"}}"#);
        default_hook(info);
    }));
    let v: Vec<u8> = Vec::new();
    let _ = std::panic::catch_unwind(|| v[3]);
}
```

### 65.2 A debugger session

```bash
cargo build                         # debug info is on by default in the dev profile
rust-lldb target/debug/kv           # or rust-gdb on Linux; both add pretty-printers for Vec, String, Option...
(lldb) b kv.rs:120                  # breakpoint at a line
(lldb) r serve
(lldb) frame variable args          # prints Vec<Vec<u8>> readably thanks to the Rust formatters
(lldb) bt                           # backtrace
```

VS Code with rust-analyzer plus the CodeLLDB extension gives you the same
thing with a UI: click "Debug" above any `fn main` or `#[test]`. Delve in
the [Go guide §58](../Golang/real-life-golang-guide.md#58-debugging-go-in-production-delve-pprof-traces-stack-dumps-godebug)
and gdb in the [C guide §14.4](../c-lang/real-life-c-guide.md#14-debugging-and-sanitizers-gdb-lldb-valgrind-asan-ubsan)
are the equivalents.

### 65.3 Async debugging

A stuck async service usually has one of three causes: a task blocked in
synchronous code (§31.5), a lock held across `.await`, or a future that is
never woken. Two tools help:

- **`tracing` spans** around each request and each important `.await`, so
  the logs show where a request stopped making progress.
- **`tokio-console`**, which works like `top` for tasks. It shows every
  task, how long it has been idle or busy, and its wake-up count. Enable it
  with the `console-subscriber` crate and `RUSTFLAGS="--cfg tokio_unstable"`.
  A task with a large "busy" time and no `.await` is blocking the executor.

### 65.4 The rest of the toolbox

| Problem | Tool |
|---|---|
| UB in `unsafe` code, data races in `unsafe` | `cargo +nightly miri test` (§33.4) |
| Memory errors in C dependencies | `RUSTFLAGS=-Zsanitizer=address cargo +nightly test` |
| Concurrency orderings | `loom` (§61.1) |
| What a macro or `derive` generated | `cargo expand` |
| Why the binary is large | `cargo bloat --release` |
| Which line regressed performance | `cargo flamegraph` (§36.3), then `git bisect run cargo bench` |
| Compile time | `cargo build --timings` (an HTML report of what took how long) |

---

## 66. Workspaces, CI quality gates, release builds, and cross-compilation

### 66.1 A workspace for a real service

```
myservice/
├── Cargo.toml            # [workspace] only
├── Cargo.lock            # ONE lockfile for every crate
├── crates/
│   ├── domain/           # pure types and logic: no I/O, fast to test
│   ├── storage/          # database access; depends on domain
│   ├── api/              # axum handlers; depends on domain + storage
│   └── cli/              # admin tool; reuses domain + storage
└── deny.toml
```

```toml
# Cargo.toml (workspace root)
[workspace]
resolver = "3"
members = ["crates/*"]

[workspace.package]
edition = "2024"
rust-version = "1.85"          # MSRV: the oldest compiler you support; CI should test it

[workspace.dependencies]       # one version of each dependency for the whole repo
tokio = { version = "1", features = ["full"] }
serde = { version = "1", features = ["derive"] }

[workspace.lints.rust]
unsafe_code = "forbid"         # opt back in per crate, with review, where unsafe is needed

[workspace.lints.clippy]
unwrap_used = "warn"
dbg_macro = "deny"

[profile.release]
lto = "thin"
codegen-units = 1
debug = "line-tables-only"     # useful backtraces and profiles from release builds
```

Member crates then write `tokio.workspace = true` and `[lints] workspace = true`.
Splitting into crates is also the main way to cut compile time. Cargo
compiles crates in parallel and only recompiles the crates that changed.
This mirrors the package boundaries in
[Go §66](../Golang/real-life-golang-guide.md#66-architecture-and-package-boundaries-services-workers-and-seams).
The difference is that Rust enforces the dependency direction, because
crates cannot have dependency cycles.

### 66.2 CI pipeline

```yaml
# .github/workflows/ci.yml
name: ci
on: [push, pull_request]
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with: { components: "clippy, rustfmt" }
      - uses: Swatinem/rust-cache@v2
      - run: cargo fmt --all --check
      - run: cargo clippy --workspace --all-targets --locked -- -D warnings
      - run: cargo test --workspace --locked
      - run: cargo doc --workspace --no-deps --locked
        env: { RUSTDOCFLAGS: "-D warnings" }
  supply-chain:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: EmbarkStudios/cargo-deny-action@v2      # advisories, licenses, bans, sources (§55)
  msrv:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@1.85
      - run: cargo check --workspace --locked
  miri:                                              # only if the workspace contains unsafe code
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@nightly
        with: { components: miri }
      - run: cargo miri test -p storage
```

Pin third-party actions to a commit SHA rather than a tag in real
pipelines. Tags can be moved
([Security in Depth Ch 12](../security/real-life-security-guide-v1.md#chapter-12-ci-cd-as-an-attack-surface)).

### 66.3 Cross-compiling

```bash
# Fully static Linux binaries (no glibc version problems on old servers):
rustup target add x86_64-unknown-linux-musl
cargo build --release --target x86_64-unknown-linux-musl

# Any target without installing a cross C toolchain: zig acts as the linker.
cargo install cargo-zigbuild
cargo zigbuild --release --target aarch64-unknown-linux-gnu.2.17   # also pins the minimum glibc

# Or run the build in a prepared container per target:
cargo install cross && cross build --release --target armv7-unknown-linux-gnueabihf
```

Static musl binaries are convenient, but musl's default allocator is
noticeably slower than glibc's for multi-threaded, allocation-heavy servers.
Install `mimalloc` as the global allocator (§43.2) when you ship musl
builds.

---

## 67. Running Rust in containers and Kubernetes

### 67.1 A small, cache-friendly image

```dockerfile
# syntax=docker/dockerfile:1
FROM rust:1 AS build
WORKDIR /src
# Build dependencies in their own layer: rebuilt only when Cargo.toml/Cargo.lock change.
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo 'fn main() {}' > src/main.rs && cargo build --release --locked && rm -rf src
COPY src ./src
RUN touch src/main.rs && cargo build --release --locked

FROM gcr.io/distroless/cc-debian12:nonroot
COPY --from=build /src/target/release/api /api
USER nonroot:nonroot
EXPOSE 3000
ENTRYPOINT ["/api"]
```

The final image contains glibc, CA certificates, and your binary: about
30 MB, with no shell and no package manager to exploit. With a musl build,
`FROM scratch` also works (add CA certificates yourself if the program
makes outbound TLS calls through the system store. `webpki-roots` builds
them into the binary). `cargo-chef` is a more robust version of the
dependency-layer trick above.

### 67.2 Behaving well under orchestration

| Concern | What to do in Rust | Across the series |
|---|---|---|
| CPU limits | `std::thread::available_parallelism()` reads the cgroup CPU quota on Linux. Tokio sizes its worker pool from it. Go needed `GOMAXPROCS` fixes for years | [OS Ch 75](../os-linux/real-life-os-guide.md#chapter-75-cpu-limits-gomaxprocs-cgroups-and-throttling) |
| Memory limits | no GC heap to tune; memory is what you allocate. Bound caches and queues, and watch RSS against `memory.max`. The OOM killer sends SIGKILL with no warning | [OS Ch 76](../os-linux/real-life-os-guide.md#chapter-76-memory-limits-the-go-heap-gomemlimit-and-the-oom-killer) |
| SIGTERM | graceful shutdown (§31.4, §47); finish within `terminationGracePeriodSeconds` | [OS Ch 74](../os-linux/real-life-os-guide.md#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1) |
| PID 1 | handle signals explicitly or use `tini` (§40.3) | same |
| Probes | `/healthz` (liveness: the process works) separate from `/readyz` (readiness: dependencies are reachable; returns 503 while draining) | [Go §67](../Golang/real-life-golang-guide.md#67-running-go-in-containers-and-kubernetes) |
| Logs | JSON to stdout: `tracing_subscriber::fmt().json()` (feature `json`) | [Go §65](../Golang/real-life-golang-guide.md#65-observability-in-go-services-logs-metrics-traces-health-profiling) |
| Metrics and traces | `metrics` + a Prometheus exporter, or OpenTelemetry via `tracing-opentelemetry` | [HTTPS Ch 23](../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services) |
| Security context | `runAsNonRoot`, `readOnlyRootFilesystem`, drop all capabilities. A static Rust binary needs none of them | [Security in Depth Ch 15](../security/real-life-security-guide-v1.md#chapter-15-container-internals-and-isolation) |

**Real-world example.** Teams that rewrite a Go or JVM service in Rust
usually report two operational changes. Memory use becomes flat and
predictable, because there is no GC heap growing to its limit. Tail latency
(p99, p999) improves more than median latency, because there are no
collector pauses. Discord's 2020 post about rewriting its "Read States"
service from Go to Rust describes exactly this pattern of latency spikes
from GC every two minutes, which disappeared after the rewrite. The
counterpoint: Go's GC has improved since then, and a well-tuned Go service
is often good enough. Measure before you rewrite.

---

## 68. How not to write Rust: an anti-pattern catalog

Each entry shows a bad form, the problem, and a fix. These are the comments
code reviewers leave most often on Rust from people coming from Go, C, or
Java.

**1. `.clone()` to silence the borrow checker**

```rust
// fragment
// BAD: copies the whole Vec on every call, just to satisfy the compiler
fn total(items: Vec<Item>) -> u64 { items.iter().map(|i| i.price).sum() }
let t = total(cart.items.clone());
// GOOD: borrow
fn total(items: &[Item]) -> u64 { items.iter().map(|i| i.price).sum() }
let t = total(&cart.items);
```

Cloning is sometimes the right call (small values, or handing owned data to
a thread). Do it on purpose, never just to make an error go away.

**2. `unwrap()` on input you do not control**

```rust
// fragment
// BAD: a malformed header crashes the request task (DoS, §51.2)
let len: usize = req.headers()["content-length"].to_str().unwrap().parse().unwrap();
// GOOD: errors are values
let len: usize = req.headers().get("content-length")
    .and_then(|v| v.to_str().ok()).and_then(|s| s.parse().ok())
    .ok_or(ApiError::BadRequest("content-length"))?;
```

**3. `Rc<RefCell<T>>` (or `Arc<Mutex<T>>`) everywhere**

Wrapping every object in shared mutable pointers rebuilds a garbage-collected
object graph, with runtime borrow panics and lock contention. Usually the
data has one natural owner. Pass `&mut` down the call stack, use indices
into a `Vec` (Ch 27), or send messages to a task that owns the state (an
actor, Ch 38).

**4. Blocking inside `async`**

```rust
// fragment
// BAD: stalls a Tokio worker thread, and every task scheduled on it
async fn handler() -> String { std::fs::read_to_string("big.json").unwrap() }
// GOOD
async fn handler() -> std::io::Result<String> { tokio::fs::read_to_string("big.json").await }
// For CPU work (hashing, compression, image processing): tokio::task::spawn_blocking (§31.1)
```

**5. A `std::sync::MutexGuard` held across `.await`**

```rust
// fragment
// BAD: the guard lives across the await (not Send; deadlock-prone)
let mut cache = state.cache.lock().unwrap();
let fresh = fetch(&key).await;
cache.insert(key, fresh);
// GOOD: lock, copy out, unlock, await, lock again
let cached = state.cache.lock().unwrap().get(&key).cloned();
let value = match cached { Some(v) => v, None => {
    let fresh = fetch(&key).await;
    state.cache.lock().unwrap().insert(key, fresh.clone());
    fresh
}};
```

**6. `String` parameters everywhere**

`fn greet(name: String)` forces every caller to allocate. Take `&str`
(or `impl AsRef<str>`), and take `String` only when the function stores
the value. Then the caller decides whether to clone or move.

**7. Stringly-typed errors**

`Result<T, String>` cannot be matched on, has no source chain, and is
impossible to handle selectively. Use `thiserror` enums in libraries and
`anyhow` with `.context()` in binaries (§10.2–10.3).

**8. Indexing loops instead of iterators**

```rust
// fragment
// BAD: bounds check per access, and an off-by-one waiting to happen
for i in 0..v.len() { if v[i] > max { max = v[i]; } }
// GOOD
let max = v.iter().copied().max();
```

**9. `unsafe` for performance before measuring**

`get_unchecked` to avoid a bounds check that LLVM had already removed. A
hand-rolled `transmute` where `from_le_bytes` exists. Each `unsafe` block
needs a profile that justifies it and a `SAFETY` comment that proves it
(§33.2). Most of the time the safe version compiles to the same code.

**10. Giant generic signatures with no `dyn` escape**

Seven type parameters threaded through every function make compile errors
unreadable and compile times long. When the dispatch is not in a hot loop,
`Box<dyn Trait>` or an `enum` is simpler and just as fast in practice
(§17.2).

**11. Ignoring `#[must_use]` results**

`let _ = tx.send(msg);` and `let _ = file.flush();` discard errors on
purpose. Make sure that is your intent. Writing the reason in a comment
(`// receiver gone: client disconnected, nothing to do`) turns a silent
bug into a documented decision.

**12. Unbounded channels and collections fed by the network**

`mpsc::unbounded_channel()`, a `Vec` that grows per message, or a
`HashMap` cache without eviction can be filled by an attacker or a slow
consumer until the OOM killer stops the process. Use bounded channels for
backpressure, cap sizes, and evict (Ch 27).

---

# Appendices

## Appendix A — Compiler-error decoder: the 15 errors you will actually see

| Code | Message (short) | What it means | Usual fix | Chapter |
|---|---|---|---|---|
| E0382 | use of moved value | you used a value after giving it away | borrow (`&x`), clone, or restructure | 6 |
| E0499 | cannot borrow as mutable more than once | two `&mut` alive at once | shorten one borrow, split the borrow | 7 |
| E0502 | cannot borrow as mutable because also borrowed as immutable | reading and writing overlap | finish reading first; copy the value out | 7 |
| E0505 | cannot move out because it is borrowed | moved a value while a reference to it is alive | end the borrow first | 19 |
| E0506 | cannot assign to borrowed value | wrote to something being read | reorder | 7 |
| E0597 | does not live long enough | a reference outlives its owner | move the owner to an outer scope, or own the data | 19 |
| E0515 | cannot return reference to local | returned `&` to a value that dies on return | return an owned value | 13 |
| E0716 | temporary value dropped while borrowed | `&String::new()`-style temporaries | bind the temporary with `let` first | 13 |
| E0106 | missing lifetime specifier | the output borrows from an ambiguous input | add `<'a>`, or return an owned type | 19 |
| E0373 | closure may outlive the current function | a thread or task closure borrows a local | `move`, and own the data | 20, 28 |
| E0277 | trait bound not satisfied | wrong type, or missing `Send`/`Clone`/`Display`... | implement or derive the trait, change the type (`Rc` → `Arc`) | 16, 28 |
| E0308 | mismatched types | often a stray `;` or a missing `.into()`/`?` | read the expected vs found types | 5 |
| E0599 | no method found | missing trait import, or the wrong type/state | `use` the trait; check typestate | 16, 62 |
| E0004 | non-exhaustive patterns | a `match` misses a case | add the arm (do not reach for `_` by reflex) | 9 |
| E0384 | cannot assign twice to immutable variable | missing `mut` | `let mut` | 3 |

`rustc --explain E0502` prints a full explanation with examples for any
code.

## Appendix B — Cargo and tooling cheat sheet

```bash
# Project
cargo new app | cargo new --lib lib | cargo init
cargo add serde --features derive | cargo remove serde | cargo update -p tokio
cargo tree | cargo tree -d | cargo tree -i openssl-sys -e features

# Build & run
cargo check | cargo build --release | cargo run --bin NAME -- ARGS
cargo build --target x86_64-unknown-linux-musl | cargo build --timings

# Quality
cargo fmt | cargo clippy --all-targets -- -D warnings | cargo doc --open
cargo test | cargo test NAME -- --nocapture | cargo test --doc | cargo nextest run

# Safety and supply chain
cargo audit | cargo deny check | cargo vet | cargo geiger | cargo +nightly miri test
cargo semver-checks

# Performance
cargo bench | cargo flamegraph --bin NAME | cargo bloat --release

# Introspection
cargo expand | rustc --explain E0499 | rustup doc --std
```

## Appendix C — Glossary

- **Borrow:** a reference (`&T` or `&mut T`) that uses a value without owning it.
- **Borrow checker:** the compiler pass that enforces shared XOR mutable and
  "no reference outlives its owner", running on MIR.
- **Crate:** a compilation unit (library or binary). **Package:** a
  `Cargo.toml` with one or more crates.
- **Drop:** the destructor. It runs when the owner goes out of scope.
- **Fat pointer:** a two-word pointer: slice (ptr + len) or trait object
  (ptr + vtable).
- **Future:** a value that produces a result when polled to completion.
  Lazy until polled.
- **Interior mutability:** mutation through `&T`, made safe by runtime
  checks or atomics (`Cell`, `RefCell`, `Mutex`, atomics).
- **Lifetime:** a compile-time name for the region where a borrow must be
  valid.
- **Monomorphization:** compiling a generic function once per concrete type.
- **Move:** transferring ownership. A bitwise copy, after which the source
  can no longer be used.
- **Niche:** an invalid bit pattern of a type that an enum uses as its tag,
  for example `None` in `Option<&T>`.
- **Orphan rule:** you can implement a trait for a type only if your crate
  defines the trait or the type.
- **`Pin`:** a guarantee that a value will not move in memory again; needed
  for self-referential futures.
- **`Send` / `Sync`:** auto traits for "can be moved to another thread" and
  "can be shared between threads by reference".
- **Typestate:** encoding an object's state in its type, so invalid
  operations do not compile.
- **UB (undefined behavior):** what happens if `unsafe` code breaks the
  rules. The compiler may assume it never happens.
- **Zero-cost abstraction:** an abstraction that compiles to the same code
  you would write by hand.

## Appendix D — A 90-day plan

Each week: read the chapters, run every snippet, break each one on purpose,
and do the stretch goals. Put your work in the `rust-labs` workspace from
§0.4.

| Weeks | Chapters | Build | Done when |
|---|---|---|---|
| 1 | 0–5 | setup; port five small Go or C programs you have written before | `cargo clippy` is clean without `#[allow]` |
| 2–3 | 6–9, 13 | the checkpoint, plus a `Matrix` type with `Add`/`Mul` and a tokenizer for arithmetic | you can explain each E0382/E0502 you hit without looking it up |
| 4 | 10–12, 14–15 | `wordfreq`, `hostsparse`; a `thiserror` error enum for one of them | errors have context and nothing panics on bad input |
| 5–6 | 16–22, 25 | `rgrep`, `lru`; implement `Iterator` for one of your own types | the LRU passes property tests (§23.3) |
| 7 | 23–24 | add unit, integration, doc, and property tests to everything so far | `cargo test` covers each project's error paths |
| 8–9 | 28–31 | `pool`, then the health checker (§31.1) and graceful shutdown (§31.4) | `tokio-console` shows no tasks blocking the executor |
| 10 | 32, 38, 45–46 | `echo`, framing, `chat`, `dnsq`, `tinyhttp` | each one survives the hostile inputs from §46 and §53 |
| 11 | 39–44 | all Part IV labs on Linux, then rewrite one of the OS guide's Go labs | `minibox` runs a shell as PID 1 in its own cgroup |
| 12 | 47–50 | `api`, `retries`, `safefetch`, `mtls`, `revproxy` | the §50 demo output matches yours, and you can explain each line |
| 13 | 51–56 | `authsvc`; fuzz `dnsq` for an hour; add `cargo deny` | no `unwrap` in request paths; `cargo deny check` passes |
| 14–15 | 57–64 | `MiniVec` under Miri, `kv`, `lifecycle`; one `kv` stretch goal | `redis-benchmark` runs against `kv`; Miri is clean |
| Ongoing | 65–68 | a workspace with CI from §66.2, deployed as a container with §67.2 settings | the service drains cleanly on `kubectl rollout restart` |

## Appendix E — Further reading and codebases worth reading

**Books (free online unless noted):**

- *The Rust Programming Language* ("the book"): the official introduction.
- *Rust for Rustaceans*, Jon Gjengset (No Starch, paid): the best
  intermediate-to-expert book; it complements Parts II and VII.
- *Rust Atomics and Locks*, Mara Bos (O'Reilly, free online): Ch 61 in
  depth, by a member of the Rust library team.
- *The Rustonomicon*: unsafe Rust, layout, variance, and building `Vec`
  (Ch 33, 57–60).
- *Asynchronous Programming in Rust* and the Tokio tutorial: Ch 30–31.
- *Zero To Production In Rust*, Luca Palmieri (paid): a production web
  service from scratch, in the style of Ch 47–56.
- The Rust API Guidelines and the Rust Reference.

**Codebases to read, in rough order of difficulty:**

| Codebase | Read it for | Related chapters |
|---|---|---|
| `ripgrep` (BurntSushi) | CLI design, fast I/O, workspace structure | 24, 26, 36 |
| `mini-redis` (Tokio team) | a teaching-quality async server; compare with Ch 63 | 31, 38, 63 |
| `axum` examples directory | every web-service pattern, small and runnable | 47, 56 |
| `hyper` | strict HTTP/1 and HTTP/2 parsing; Ch 50's framing behavior | 46, 50 |
| `rustls` | protocol state machines in types; careful crypto API design | 49, 52, 62 |
| `std::collections` source (`Vec`, `HashMap`) | `unsafe` done well, with SAFETY comments | 33, 58 |
| `tokio` scheduler and `mio` | the reactor and work stealing from Ch 59 | 30, 59 |
| `youki` (OCI runtime) | Ch 44 with all the hard parts filled in | 39–44 |
| `sled` / `redb` | storage engines, crash safety, page caches | 41, 63 |
| Rust for Linux (`rust/kernel/` in the kernel tree) | Rust in kernel context, safe wrappers over C | 33, 34; [C Part IX](../c-lang/real-life-c-guide.md#47-two-worlds-user-space-hardware-i-o-vs-writing-a-kernel-driver) |

**Where to go next in the series:** if you have not done so, work through
the [OS guide](../os-linux/real-life-os-guide.md),
[networking](../networking/tcp-ip/real-life-guide-v1.md),
[security](../security/real-life-guide.md), and
[HTTPS lifecycle](../v2-https/real-life-guide-v1.md) guides with this one
open next to them. Every lab in those guides that is built in Go can also
be built in Rust, and §0.4 lists the ones this guide has already rewritten.
Doing the rest yourself is the best way to practise.
