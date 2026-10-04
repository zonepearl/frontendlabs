# Go — The Complete Field Guide (Beginner → Expert)

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/go/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> A practical, example-driven path through Go: what the language actually
> gives you, why it's shaped the way it is, and how to use it to build real
> terminal tools, GUI apps, network services, and security infrastructure.
>
> Every concept is paired with runnable code, a "Real-world example" or "War
> story," and — at the end of each Part — a full mini-project you can build
> and run today. Read it once top to bottom; keep it as a lookup guide.

---

> **The series:** 1 OS → 2 Networking → 3 Security → 4 HTTPS walkthrough, with Go alongside.
>
> **This guide runs alongside every step.** Start the series at [step 1: Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md).
>
> **Languages track:** Go → [C](../c-lang/real-life-c-guide.md) → [Rust](../rust-lang/real-life-rust-guide.md).
>
> [The full series map](#the-series-os-networking-security-https-in-go).

---

## How to use this guide

- **Beginner (Parts I):** the Go toolchain, types, control flow, functions,
  slices/maps, structs, interfaces, errors, packages. Ends with two terminal
  mini-projects.
- **Intermediate (Parts II):** goroutines, channels, `sync`, `context`,
  testing, generics, `net/http`, and the standard library you'll use every
  day. Ends with a worker-pool tool and an LRU cache service.
- **Advanced (Part III):** the scheduler, the memory model and GC,
  concurrency patterns, advanced generics, raw networking, TLS, profiling.
  Ends with a TCP chat server and a mini API gateway.
- **Real-world builds (Parts IV–VI):** production CLIs, a job queue, a log
  analyzer, GUI apps (Fyne, Wails), and security-focused services (JWT auth,
  a WAF middleware, mTLS) — the same kind of capstones as the
  networking/security/AI field guides in this Wiki.
- **Expert (Part VII):** what's actually in memory behind an interface, a
  slice, a map, and a goroutine; API design idioms; a capstone database
  engine.
- **Professional Go (Part VIII):** debugging, refactoring, serious testing,
  production databases and HTTP APIs, dependency/release management,
  observability, architecture, and container/Kubernetes operations.
- **Reference (Appendices):** idioms, pitfalls, glossary, cheat sheets.

Conventions:
- Code is Go, targeting Go 1.22+ (the guide calls out version-specific
  behavior, e.g. the pre-1.22 loop-variable capture bug).
- "**Real-world example**" boxes ground a language concept in something you
  will actually write. "**War story**" boxes describe an incident or design
  trade-off from real production systems.
- This guide is the *theory + architecture* companion to
  [`golang-90-day-plan.md`](./golang-90-day-plan.md) (day-by-day exercises)
  and [`projects/lru_cache.md`](./projects/lru_cache.md) (a worked project
  spec). Where useful, this guide cross-references them instead of repeating
  them.

---

## The series: OS → networking → security → HTTPS, in Go

This guide is one of seven, designed to be read as **one course** in this order:

```
   1  OS & Linux  ──▶  2  Networking  ──▶  3  Security  ──▶  4  HTTPS walkthrough
                         2a OSI map          3a From Zero        (one request through
                         2b TCP/IP           3b In Depth          every guide, in Go)
   ════════════════════  Go, alongside every step  ════════════════════
```

| Step | Guide | What it gives you | Hands-on |
|---|---|---|---|
| 1 | [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md) | the machine every request starts and ends on: processes, memory, files, sockets, containers, Kubernetes | 9 Go labs ([Part 20](../os-linux/real-life-os-guide.md#part-20-systems-programming-in-go-the-os-from-inside-a-program)) |
| 2a | [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md) | one click followed through all seven layers; the map for everything after it | concept map, 1 hour |
| 2b | [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md) | how machines talk: addressing, routing, TCP, TLS, packet capture, network operations | 23 Go labs ([§0.8](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself)) |
| 3a | [Security from Zero](../security/real-life-guide.md) | encoding vs hashing vs encryption, TLS, PKI, SSH, the OWASP Top 10, threat modelling | OpenSSL labs and a vulnerable app to break |
| 3b | [Security Engineering in Depth](../security/real-life-security-guide-v1.md) | cloud and Kubernetes security, distributed authorization, advanced web attacks, data protection, detection | 10 Go labs + a `govulncheck` exercise ([§0.8](../security/real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)) |
| 4 | [The HTTPS Request Lifecycle](../v2-https/real-life-guide-v1.md) | one request end to end, then the server side built in Go; Chapter 25 traces one request through every guide | 13 Go labs ([§0.6](../v2-https/real-life-guide-v1.md#0-6-the-go-labs-build-the-lifecycle-yourself), [Ch 25](../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)) |
| ∥ | **Go — The Complete Field Guide** ← you are here | the language behind every lab, plus the 120-day plan's multi-week projects | [120-day plan](golang-90-day-plan.md) |
| L2 | [C — The Complete Field Guide](../c-lang/real-life-c-guide.md) | manual memory, UB, sockets, and the kernel: what Go's runtime hides | mini-projects in every Part |
| L3 | [Rust — The Complete Field Guide](../rust-lang/real-life-rust-guide.md) | C's control with compiler-checked memory safety; rebuilds this table's labs in Rust | [Rust labs (§0.4)](../rust-lang/real-life-rust-guide.md#0-4-the-rust-labs-what-you-build-for-each-guide-in-the-series) |

**Why this order.** Every network connection is a file descriptor owned by
a process, so the OS comes first. Networking comes next, because every attack and
defence in the security guides assumes you can follow a packet. Security comes
third, because it needs both. The HTTPS walkthrough comes last because a
single HTTPS request uses all of it, and [its final chapter](../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)
runs one request through every guide in a Go program. Go runs alongside the
whole way: each guide's "Build it in Go" labs turn its chapters into code you
can run.

**Reading paths by role:**

- **Backend engineer:** 1 (Phases 1–2, Part 20) → 2a → 2b (Parts 1–7) → 3a (Parts 1–6) →
  4 (all, including Part 9) → 3b (Parts 5–6).
- **SRE / platform:** 1 (all) → 2a → 2b (all) → 3b (Parts 2–4, 8) → 4 (Parts 8–10).
- **Security engineer:** 1 (Phases 1, 3) → 2a → 2b (Parts 3–8) → 3a → 3b → 4 (Chapters 12, 19–22, 25).
- **Engineering manager:** 2a → 1 (Chapters 1–3, 44, 66, 72) → 3b (Part 9) →
  4 (Chapters 13, 18, 23, 25).

Chapters link to each other directly, with **"Across the series"** notes where
another guide covers a topic in more depth.

### What each guide builds in Go

| Guide | Go labs | Builds on these chapters here |
|---|---|---|
| [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md) ([index](../os-linux/real-life-os-guide.md#0-8-the-go-labs-the-os-from-inside-a-program)) | goroutines vs OS threads, PID 1 reaping, GOMAXPROCS and GOMEMLIMIT in containers, atomic writes, `ps` from `/proc`, a container runtime, pprof hunting | 27, 28, 33, 37, 57, 58, 67 |
| [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md) ([index](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself)) | DNS client from raw bytes, UDP loss meter, pcap decoder, TLS inspector, CLOSE_WAIT and Nagle reproductions, draining load balancer, PROXY protocol, blackbox prober | 15–18, 22, 24, 34–35 |
| [Security Engineering in Depth](../security/real-life-security-guide-v1.md) ([index](../security/real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)) | attack-path graphs, envelope encryption, SPIFFE mTLS, JWTs and token exchange, ReBAC, race conditions, `os.Root` uploads, detections as code, tamper-evident logs, `govulncheck` | 17, 22, 29, 35, 47–50, 60, 64 |
| [The HTTPS Request Lifecycle](../v2-https/real-life-guide-v1.md) ([index](../v2-https/real-life-guide-v1.md#0-6-the-go-labs-build-the-lifecycle-yourself)) | hardened server, Slowloris, reverse proxy and smuggling, CDN cache, idempotent retries and breakers, SSRF guard, SLO alerts, SSE, and the whole-series walkthrough | 18, 24, 29, 35, 36, 47, 62 |

In total: **55 programs**, each run and verified. The [120-day plan](golang-90-day-plan.md) grows
the same ideas into multi-week projects.

## Table of contents

**Part I — Foundations: the Go way**
1. What Go is, and why it exists
2. The toolchain and workspace
3. Variables, types, zero values, constants
4. Control flow: `if`, `switch`, `for`, `defer`
5. Functions: multiple returns, variadics, closures
6. Arrays and slices: the internals that explain the bugs
7. Maps: internals and idioms
8. Structs and methods: value vs. pointer receivers, embedding
9. Interfaces: implicit satisfaction, `any`, type switches
10. Error handling: the pattern that shapes all Go code
11. Packages and modules
12. Strings, runes, bytes, and Unicode
13. Terminal project: a system-info CLI and unit converter
14. Terminal project: a JSON config server

**Part II — Intermediate: concurrency and the standard library**
15. Goroutines: Go's answer to "what is a thread"
16. Channels: pipelines, `select`, and the closing convention
17. The `sync` package: mutexes, `WaitGroup`, `Once`, atomics
18. `context`: cancellation, deadlines, and request-scoped values
19. Testing in Go: table-driven tests and subtests
20. Benchmarks and a first look at `pprof`
21. Generics: type parameters and constraints
22. I/O: `io.Reader`/`Writer`, `bufio`, and `os`
23. `encoding/json` and reflection, briefly
24. `net/http` fundamentals: client, server, middleware
25. Terminal project: a concurrent URL health-checker (worker pool)
26. Terminal project: an LRU cache library and CLI cache server

**Part III — Advanced: systems Go**
27. The Go memory model and the garbage collector
28. The scheduler deep dive: GMP and work-stealing
29. Concurrency patterns: fan-in/fan-out, pipelines, rate limiting, `errgroup`
30. Advanced error handling: `errors.Join`, custom types, panic boundaries
31. Advanced generics: constraints, generic algorithms, inference limits
32. Reflection and `unsafe`: when, and why rarely
33. `cgo` and syscalls, briefly
34. Networking deep dive: `net.Conn`, TCP/UDP, framing your own protocol
35. TLS in Go: `crypto/tls`, mutual TLS
36. Structured logging and observability: `log/slog`, metrics, tracing
37. Performance tuning: `pprof`, `trace`, and benchmark methodology
38. Terminal project: a TCP chat server with rooms
39. Terminal project: a rate-limited reverse proxy (mini API gateway)

**Part IV — Real-world builds: terminal**
40. Building production CLIs: `flag`, subcommands, config precedence
41. Capstone: a distributed-style job queue with persistence
42. Capstone: a streaming log analyzer for multi-GB files

**Part V — Real-world builds: GUI**
43. GUI options in Go: Fyne, Gio, Wails, Walk, and "just serve HTML"
44. Building a Fyne desktop app: a live system-resource monitor
45. Building a Wails app: a Go backend behind a web-based to-do board
46. Packaging and distributing GUI apps

**Part VI — Security-focused Go**
47. Secure coding practices in Go
48. Auth service: password hashing and JWT issuing/verification
49. A WAF-style middleware: rate limiting and IP blocklists
50. mTLS: client-certificate authentication end to end

**Part VII — Expert: internals and design**
51. What's really behind an interface value
52. What's really behind a slice, a map, and a string
53. Goroutine stacks and the scheduler, revisited
54. Designing APIs: the Go proverbs, applied
55. Dependency injection without a framework: functional options
56. Capstone: a tiny embedded database engine (Bitcask-style)
57. Shipping it: cross-compilation, Docker, systemd, graceful shutdown

**Part VIII — Professional Go: debugging, refactoring, operations**
58. Debugging Go in production: Delve, pprof, traces, stack dumps, GODEBUG
59. Refactoring Go codebases: package seams, APIs, and compatibility
60. Testing beyond tables: fuzzing, httptest, golden files, integration tests
61. Production database access with `database/sql` and `pgx`
62. Production HTTP APIs: validation, timeouts, middleware, shutdown
63. Modules and releases: MVS, private modules, versioning, binaries
64. Professional workflow: linting, vuln checks, CI, and quality gates
65. Observability in Go services: logs, metrics, traces, health, profiling
66. Architecture and package boundaries: services, workers, and seams
67. Running Go in containers and Kubernetes

**Appendices**
- A. Go proverbs and idioms cheat sheet
- B. Common pitfalls and how to catch them
- C. Glossary
- D. Further reading
- E. Command cheat sheet

---

# Part I — Foundations: the Go way

The language itself, in order: what problem it was built to solve, the
toolchain, and then the pieces every Go program is made of — types, control
flow, the error-handling pattern that shapes all Go code, and just enough of
slices and maps' internals to explain the bugs you'll eventually hit. Two
terminal projects at the end put it all to work before Part II adds
concurrency on top.

## 1. What Go is, and why it exists

### 1.1 The problem being solved

Go was created at Google in 2007–2009 by Robert Griesemer, Rob Pike, and Ken
Thompson, in response to a very specific pain: Google's C++ services took
tens of minutes to build, dependency graphs were tangled, and the languages
available (C++, Java) made concurrent, networked systems hard to write
safely. They wanted a language with the safety and readability of a
garbage-collected language, the fast compilation of an interpreted one, and
concurrency as a first-class citizen instead of a library bolted on top.

### 1.2 The design goals, and how they show up in the syntax

| Goal | Consequence in the language |
|---|---|
| Fast compilation | No header files, no generics until 1.18 (kept simple even then), strict unused-import/variable errors so the compiler never guesses |
| Readability at scale | One way to format code (`gofmt`), no operator overloading, no implicit conversions |
| Safe concurrency | Goroutines + channels as language primitives, not a library; the race detector is built into the toolchain |
| Simple deployment | Single-file binaries are common; pure-Go builds are often static, while cgo builds may depend on system libraries |
| Explicitness | Errors are values, not exceptions; nothing happens "automatically" except `defer`, GC, and slice growth |

### 1.3 What Go deliberately does not have

No exceptions (only `panic`/`recover`, reserved for truly exceptional
conditions), no classes or inheritance (only structs, interfaces, and
composition), no operator overloading, no implicit numeric conversions, and
—until you explicitly opt in with generics—no generic containers. This is
not an oversight; it is the single biggest reason a large, multi-team Go
codebase stays readable after five years. A Go file you've never seen still
looks like every other Go file.

### 1.4 War story

A team migrating a Python API gateway to Go expected the rewrite to take
three months. It took six weeks, and the resulting binary went from a
150MB Docker image with a Python interpreter and dependencies to an 18MB
`FROM scratch` image with a single static binary. The surprise wasn't
performance — it was that code review got *faster*, because there was
no longer a debate about "should this be a decorator, a metaclass, or a
mixin." In Go there's usually one obvious way to write a given piece of
logic, and reviewers spend their time on the logic instead of the style.

```bash
# The whole "hello world" loop
mkdir hello && cd hello
go mod init example.com/hello
cat > main.go <<'EOF'
package main

import "fmt"

func main() {
	fmt.Println("hello, go")
}
EOF
go run main.go
```

---

## 2. The toolchain and workspace

### 2.1 The commands you'll type every day

| Command | What it does |
|---|---|
| `go run .` | Compile to a temporary binary and run it |
| `go build` | Produce a binary for the current OS/arch |
| `go test ./...` | Run every test in the module |
| `go vet ./...` | Static analysis for correctness bugs (e.g. `Printf` format mismatches) |
| `go fmt ./...` | Reformat every file to the one canonical style |
| `go mod init <module>` | Start a new module, write `go.mod` |
| `go mod tidy` | Add missing requirements, remove unused ones |
| `go get pkg@version` | Add/upgrade a dependency |
| `go install pkg@version` | Build and install a binary into `$GOPATH/bin` |
| `go doc pkg.Symbol` | Read documentation without leaving the terminal |

### 2.2 Modules, not GOPATH

Pre-2019 Go required your code to live inside `$GOPATH/src`. Modules
(`go.mod`) decoupled that: any directory can be a module root. `go.mod`
pins your module's own import path and its dependencies' versions;
`go.sum` pins their cryptographic hashes so a build is reproducible even if
a dependency's tag is force-pushed upstream.

```go
// go.mod
module github.com/kumaran/toolbelt

go 1.22

require (
	github.com/spf13/cobra v1.8.0
)
```

### 2.3 Cross-compilation, usually simple for pure-Go code

Go's compiler targets are selected with two environment variables, and
pure-Go programs usually do not need a separate compiler toolchain per
target. If your build uses cgo, native libraries, or OS-specific files, you
may need a matching C cross-compiler and target libraries.

```bash
GOOS=linux   GOARCH=amd64 go build -o toolbelt-linux-amd64 .
GOOS=darwin  GOARCH=arm64 go build -o toolbelt-mac-arm64  .
GOOS=windows GOARCH=amd64 go build -o toolbelt.exe        .
```

### 2.4 Real-world example

A pure-Go CLI you build on your Apple Silicon laptop, tested with
`go test ./...`, can often be cross-compiled into a Linux amd64 binary and
dropped into a small container image with no language runtime dependency —
no `python3`, no `node_modules`, no JVM. This deployment story is one
reason so much infrastructure tooling (Docker, Kubernetes, Terraform,
Prometheus) is written in Go.

---

## 3. Variables, types, zero values, constants

### 3.1 Declaring variables

```go
var count int              // zero value: 0
var name string            // zero value: ""
var ready bool             // zero value: false
var price float64          // zero value: 0.0
x := 42                    // short form; type inferred, only inside functions
var y int = 42             // equivalent, explicit
```

The **zero value** rule is load-bearing: every declared-but-unassigned
variable in Go gets a well-defined, usable value — `0`, `""`, `false`,
`nil` — never garbage memory. A `struct` zero value is every field
zeroed recursively. This eliminates an entire class of "uninitialized
variable" bugs common in C.

### 3.2 Basic types

`bool`, `string`, `int`/`int8`/`int16`/`int32`/`int64`,
`uint`/`uint8`(`byte`)/`uint16`/`uint32`/`uint64`, `float32`/`float64`,
`complex64`/`complex128`, `rune` (alias for `int32`, a Unicode code point).
`int` and `uint` are platform-width (64-bit on virtually everything you'll
target today) — use them for general-purpose counting, and use the sized
variants only when the width matters (wire formats, hashing, bit twiddling).

### 3.3 Constants and `iota`

Constants are untyped until used, which lets `const Pi = 3.14159` work as a
`float64` in one expression and (if it fit) an `int` in another. `iota`
generates enumerations without hand-numbering:

```go
type Unit int

const (
	Byte Unit = 1 << (10 * iota) // iota = 0 -> 1
	KB                            // iota = 1 -> 1024
	MB                            // iota = 2 -> 1024*1024
	GB                            // iota = 3 -> 1024*1024*1024
)

func (u Unit) String() string {
	switch u {
	case KB:
		return "KB"
	case MB:
		return "MB"
	case GB:
		return "GB"
	default:
		return "B"
	}
}
```

### 3.4 Type conversion is always explicit

```go
var i int = 42
var f float64 = float64(i) // no implicit int->float64
var u uint = uint(f)       // no implicit float64->uint either
```

There is no silent widening or narrowing anywhere in Go. If a bug turns an
`int32` into an `int64` unexpectedly, it's because you wrote the
conversion, not because the compiler inferred one.

### 3.5 Real-world example

The `Unit` + `iota` pattern above is the backbone of Day 2 in the 90-day
plan's "network unit converter" exercise — the same trick scales to HTTP
status classes, log levels (`DEBUG`/`INFO`/`WARN`/`ERROR`), and protocol
opcodes, anywhere you need a small closed set of named integers with a
`String()` method for readable logs.

---

## 4. Control flow: `if`, `switch`, `for`, `defer`

### 4.1 `if` with an initializer

```go
if v, err := parse(input); err != nil {
	return fmt.Errorf("parse: %w", err)
} else {
	use(v)
}
```

`v` and `err` are scoped to the `if`/`else` block only — a deliberate
narrowing that keeps error-carrying temporaries from leaking into the rest
of the function.

### 4.2 `switch` without fallthrough

Unlike C, Go's `switch` does **not** fall through by default — each case
implicitly breaks. Use the `fallthrough` keyword if you actually want the
old behavior. `switch` also works without an expression at all (a cleaner
`if/else if` chain) and with a type assertion (`switch v := x.(type)`,
covered in Chapter 9).

```go
switch status := resp.StatusCode; {
case status >= 500:
	log.Println("server error")
case status >= 400:
	log.Println("client error")
case status >= 300:
	log.Println("redirect")
default:
	log.Println("ok")
}
```

### 4.3 `for` is the only loop

There's no `while`, no `do-while` — `for` covers every case:

```go
for i := 0; i < 10; i++ { }      // classic
for cond { }                      // while
for { }                           // infinite (use break)
for i, v := range items { }       // range over slice/array/map/string/channel
for range ch { }                  // discard the value, just drain the channel
```

### 4.4 The loop-variable capture trap (fixed in Go 1.22)

Before Go 1.22, `for`-loop variables were shared across iterations, which
made this a classic bug:

```go
// Go < 1.22: prints "3 3 3" — all goroutines captured the same `i`
for i := 0; i < 3; i++ {
	go func() { fmt.Println(i) }()
}
```

Go 1.22+ gives each iteration its own `i`, so this now correctly prints
`0 1 2` (in some order). **Know both behaviors** — you will read code
written for both, and `go.mod`'s `go` directive controls which semantics
apply.

### 4.5 `defer`: cleanup that always runs, in LIFO order

```go
func processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close() // runs when processFile returns, however it returns

	// ... use f ...
	return nil
}
```

Three rules to internalize: deferred calls run in **LIFO** order; their
**arguments are evaluated immediately** when `defer` executes, not when the
deferred call finally runs; and `defer` runs even if the function panics —
which is exactly how `recover()` (Chapter 30) gets a chance to run.

```go
func demo() {
	for i := 0; i < 3; i++ {
		defer fmt.Println("deferred:", i) // captures i's value *now*
	}
	fmt.Println("done")
}
// prints: done / deferred: 2 / deferred: 1 / deferred: 0
```

### 4.6 Real-world example

Many HTTP handlers that open a database transaction follow this shape:
check `Begin`/`BeginTx`, `defer tx.Rollback()`, then `tx.Commit()` on the
success path. If the handler returns early on any error, the deferred
`Rollback()` cleans up the still-open transaction; after a successful
`Commit`, `Rollback` returns an error such as "transaction already
committed" and is typically ignored.

---

## 5. Functions: multiple returns, variadics, closures

### 5.1 Multiple return values and the `(value, error)` idiom

```go
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
```

This single pattern — return the real result *and* an explicit `error` —
is why Go code reads the same everywhere: every call site immediately
answers "did this work?" without a hidden control-flow jump. There's no
exception to catch three call frames up; the caller either handles the
error right there or explicitly propagates it.

### 5.2 Named returns

```go
func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // "naked" return sends back the current x, y
}
```

Named returns are most useful combined with `defer` to modify a return
value after the fact — e.g. wrapping an error with context right before
the function actually returns:

```go
func doWork() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("doWork: %w", err)
		}
	}()
	return someOperation()
}
```

### 5.3 Variadic functions

```go
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

sum(1, 2, 3)       // nums = []int{1,2,3}
sum(existingSlice...) // spread a slice into variadic args
```

### 5.4 Closures

Functions are values; they close over variables by reference, not by
value:

```go
func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

next := counter()
next() // 1
next() // 2
```

### 5.5 Real-world example

A rate limiter is naturally a closure: `NewLimiter(rps int) func() bool`
returns a function that closes over a token bucket and internal timers,
so callers never see (or can misuse) the bucket's internal state — only
the decision. This is the functional-options-adjacent pattern you'll use
constantly for middleware (Chapter 24) and for the WAF-style rate limiter
in Chapter 49.

---

## 6. Arrays and slices: the internals that explain the bugs

### 6.1 Arrays are fixed-size value types

```go
var a [3]int          // [0 0 0], part of the containing struct/stack frame
b := a                // b is a full COPY of a
b[0] = 99              // does not affect a
```

Arrays are rarely used directly in Go; they exist mostly as the backing
store beneath slices, or where a fixed size is a semantic guarantee (e.g.
`[32]byte` for a SHA-256 digest).

### 6.2 A slice is a three-word header

This is the single most important mental model in this chapter. A slice
value is *not* the data — it's a small struct:

```go
type sliceHeader struct {
	ptr *T   // pointer into a backing array
	len int  // number of elements currently visible
	cap int  // capacity of the backing array from ptr onward
}
```

That means **slicing a slice does not copy the underlying data** — it
creates a new header pointing into the same array:

```go
original := []int{1, 2, 3, 4, 5}
view := original[1:3]     // len=2, cap=4, shares backing array
view[0] = 99              // original[1] is now 99 too!
```

### 6.3 `append` and the growth trap

`append` only allocates a new backing array when the current one lacks
capacity. If it *doesn't* need to grow, it writes into the existing array
— which can silently corrupt a slice some other variable still references:

```go
a := make([]int, 3, 5) // len=3, cap=5 — room to grow without reallocating
b := append(a, 99)     // still within cap, so b shares a's backing array
a = append(a, 42)      // ALSO fits — overwrites the slot b just wrote to!
fmt.Println(b)          // [.. .. .. 42] — NOT 99. Classic aliasing bug.
```

The fix, whenever you hand a sub-slice to code that might append to it, is
to force a fresh allocation with a **full slice expression**
(`s[low:high:max]`, setting `cap == len`) or with `append(s[:0:0], s...)`.

### 6.4 Growth strategy and pre-allocation

Go's runtime roughly doubles capacity for small slices and grows by ~1.25x
for large ones (implementation detail, not a spec guarantee — don't rely
on exact numbers, but do rely on "occasional reallocation + copy"). If you
know the final size, `make([]T, 0, n)` avoids all of the intermediate
reallocations:

```go
// Without pre-allocation: several reallocations as the slice grows
results := []Result{}
for _, item := range items {
	results = append(results, process(item))
}

// With pre-allocation: one allocation, zero reallocations
results := make([]Result, 0, len(items))
for _, item := range items {
	results = append(results, process(item))
}
```

### 6.5 The `range` loop copies each element

```go
type Big struct{ data [1024]byte }
items := []Big{{}, {}, {}}
for _, item := range items {
	item.data[0] = 1 // mutates the COPY, not items[i]
}
// Fix: index into the slice directly, or range with pointers
for i := range items {
	items[i].data[0] = 1
}
```

### 6.6 War story

A caching layer took a `[]byte` from a `sync.Pool`, sliced off a header,
handed the rest to a downstream handler as the "body," and returned the
original buffer to the pool for reuse — while the downstream handler was
still reading from a slice that *aliased the same backing array* the pool
was about to overwrite for the next request. The bug appeared as
intermittent, unreproducible byte corruption in ~1 in 50,000 requests
under load, and took two engineers a full day with the race detector and
`-fsanitize`-style buffer poisoning before the aliasing was spotted. The
fix was one line: copy the body (`bytes.Clone` in modern Go) before
returning the buffer to the pool. **Lesson: slices are views, not
copies — treat any slice you didn't allocate yourself as borrowed, not
owned, unless the API contract says otherwise.**

---

## 7. Maps: internals and idioms

### 7.1 Declaration and the comma-ok idiom

```go
m := make(map[string]int)
m["a"] = 1

v, ok := m["b"] // ok=false, v=0 (the zero value) — key absent, not an error
if ok {
	use(v)
}

delete(m, "a") // no-op if key absent, never panics
```

Reading a missing key never panics — you get the zero value. This is a
frequent bug source for beginners who forget to check `ok` and silently
treat "absent" as "present with value 0."

### 7.2 Nil maps: readable, not writable

```go
var m map[string]int // nil, not an empty map
v := m["x"]           // fine — returns 0
m["x"] = 1            // PANIC: assignment to entry in nil map
```

Always initialize with `make` (or a composite literal) before writing.

### 7.3 Iteration order is intentionally randomized

Go's map iteration order is randomized *by design*, specifically to stop
anyone from depending on it. If you need sorted output, collect keys into
a slice and sort them:

```go
keys := make([]string, 0, len(m))
for k := range m {
	keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
	fmt.Println(k, m[k])
}
```

### 7.4 Maps are not safe for concurrent use

A concurrent map write (or a write racing a read) panics with `fatal
error: concurrent map read and map write` — Go detects this at runtime
rather than silently corrupting the map. Guard with a `sync.RWMutex`
(Chapter 17) or use `sync.Map` for specific high-read/rare-write patterns.

### 7.5 Structs, not strings, as keys

Any comparable type can be a map key — including structs — which gives
you composite keys without string concatenation:

```go
type CacheKey struct {
	UserID int
	Region string
}
cache := make(map[CacheKey]Result)
cache[CacheKey{UserID: 42, Region: "us-east"}] = result
```

### 7.6 Real-world example

Day 4 of the 90-day plan's ring-buffer request log uses exactly this
pattern: a slice for the fixed-size ring buffer plus a `map[string][]int`
indexing IPs to positions in the ring, giving O(1) "show me the last N
requests from this IP" without scanning the whole buffer.

---

## 8. Structs and methods: value vs. pointer receivers, embedding

### 8.1 Struct basics and tags

```go
type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	TLS  bool   `json:"tls,omitempty"`
}
```

Struct tags are metadata strings parsed by reflection at runtime
(`encoding/json`, `encoding/xml`, validators, ORMs all read them) — they
have zero effect on the compiler beyond being valid string literals.

### 8.2 Value vs. pointer receivers — the rule that matters

```go
type Counter struct{ n int }

func (c Counter) IncByValue()  { c.n++ }   // mutates a COPY — no visible effect
func (c *Counter) IncByPointer() { c.n++ } // mutates the real struct
```

**Rule of thumb:** if *any* method on a type needs a pointer receiver
(to mutate state, or to avoid copying something large), make *all*
methods on that type use pointer receivers, for consistency and because
Go's method set rules mean a value of type `T` only gets the pointer
receiver methods promoted automatically when it is *addressable* — a map
value, for instance, is not addressable and this trips people up:

```go
counters := map[string]Counter{"a": {}}
counters["a"].IncByPointer() // compile error — map value not addressable
```

### 8.3 Embedding: composition, not inheritance

```go
type Logger struct{ Prefix string }

func (l Logger) Log(msg string) { fmt.Println(l.Prefix, msg) }

type Server struct {
	Logger // embedded — promotes Log() onto Server
	Host   string
}

s := Server{Logger: Logger{Prefix: "[server]"}, Host: "localhost"}
s.Log("starting up") // calls the embedded Logger's method directly
```

Embedding promotes fields and methods to the outer type, but it is *not*
polymorphism — there's no dynamic dispatch, no `super`, and the embedded
type has no idea it's embedded. It's syntactic sugar for "has-a with
automatic delegation," which is exactly the composition-over-inheritance
philosophy Go pushes everywhere.

### 8.4 Real-world example

Day 5 of the 90-day plan asks you to embed a `Logger` struct into a
`Server` struct — this is the idiom behind virtually every Go SDK's
"client with options" pattern: an HTTP client embeds a base transport,
adds retry/auth logic, and callers use it exactly like the thing it
wraps, because embedding promoted all the underlying methods.

---

## 9. Interfaces: implicit satisfaction, `any`, type switches

### 9.1 Interfaces are satisfied implicitly

There is no `implements` keyword. A type satisfies an interface simply by
having the right methods — which means you can define an interface *after*
the type already exists, even in a different package than the type, and
even define it just for a test:

```go
type Storage interface {
	Set(key, value string) error
	Get(key string) (string, bool)
	Delete(key string) error
}

type MemStorage struct{ data map[string]string }

func (m *MemStorage) Set(k, v string) error { m.data[k] = v; return nil }
func (m *MemStorage) Get(k string) (string, bool) { v, ok := m.data[k]; return v, ok }
func (m *MemStorage) Delete(k string) error { delete(m.data, k); return nil }

var s Storage = &MemStorage{data: map[string]string{}} // compiles: shape matches
```

### 9.2 Small interfaces are idiomatic — "accept interfaces, return structs"

The standard library's most reused interfaces are one or two methods:
`io.Reader` (`Read([]byte) (int, error)`), `io.Writer`
(`Write([]byte) (int, error)`), `fmt.Stringer` (`String() string`),
`error` (`Error() string`), `sort.Interface`. Small interfaces mean any
type can satisfy many of them at once, and any function that *accepts* an
interface argument works with anything that shape — this is how
`io.Copy(dst io.Writer, src io.Reader)` works identically whether the
source is a file, a network socket, an in-memory buffer, or a gzip
decompressor.

### 9.3 `any` and type assertions/switches

`any` is an alias for `interface{}` — a value of any type at all, with no
methods guaranteed. Getting the concrete type back out requires an
assertion:

```go
func describe(v any) string {
	switch x := v.(type) {
	case int:
		return fmt.Sprintf("int: %d", x)
	case string:
		return fmt.Sprintf("string: %q", x)
	case error:
		return fmt.Sprintf("error: %v", x)
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("unhandled type %T", x)
	}
}
```

A single (non-switch) assertion has a two-value "comma-ok" form that
never panics:

```go
if s, ok := v.(string); ok {
	use(s)
}
s := v.(string) // panics if v is not a string — only use when you're certain
```

### 9.4 The nil-interface trap

An interface value is `nil` only if *both* its type and value are nil.
A `nil` pointer stored inside an interface is **not** a nil interface —
this is one of the most common "why doesn't `err != nil` work" bugs:

```go
type MyError struct{}
func (e *MyError) Error() string { return "boom" }

func doWork() error {
	var e *MyError = nil
	return e // returns a non-nil error interface wrapping a nil *MyError!
}

err := doWork()
fmt.Println(err == nil) // false! surprising to everyone the first time
```

The fix: functions that return `error` should `return nil` explicitly
when there's no error, never a typed nil pointer.

### 9.5 Real-world example

Day 6's `Storage` interface, wrapped by a `CachedStorage` that adds
metrics, is the canonical Go decorator: `CachedStorage` embeds a
`Storage` field, implements the same interface, and adds behavior around
each call — no framework, no annotations, just composition plus
interfaces. This exact shape reappears in Chapter 24 as HTTP middleware
and in Chapter 49 as the WAF middleware chain.

---

## 10. Error handling: the pattern that shapes all Go code

### 10.1 `error` is just an interface

```go
type error interface {
	Error() string
}
```

That's the entire contract. `errors.New("message")` and
`fmt.Errorf("...")` both return values satisfying it.

### 10.2 Wrapping with `%w`, and unwrapping with `errors.Is`/`errors.As`

```go
var ErrNotFound = errors.New("not found")

func lookup(key string) (string, error) {
	if _, ok := store[key]; !ok {
		return "", fmt.Errorf("lookup %q: %w", key, ErrNotFound)
	}
	return store[key], nil
}

_, err := lookup("missing")
if errors.Is(err, ErrNotFound) {
	// true even though the error message has extra context wrapped around it
}
```

`errors.Is` walks the `Unwrap()` chain checking for a target sentinel;
`errors.As` walks the chain looking for a target *type* so you can pull
out a custom error's fields:

```go
type ValidationError struct {
	Field string
	Msg   string
}
func (e *ValidationError) Error() string { return e.Field + ": " + e.Msg }

var verr *ValidationError
if errors.As(err, &verr) {
	fmt.Println("bad field:", verr.Field)
}
```

### 10.3 Sentinel errors vs. custom error types

Use a **sentinel** (`var ErrX = errors.New(...)`) when callers only need
to know *which* known condition occurred. Use a **custom type** when
callers need structured data about the failure (which field, what limit
was exceeded, etc). Both compose through wrapping.

### 10.4 What Go deliberately does *not* do: exceptions for control flow

There is no `try/catch`. `panic`/`recover` exist (Chapter 30) but are
reserved for programmer errors and unrecoverable states (index out of
range, nil dereference) — not for "user typed an invalid email." That
case is an `error` return, checked immediately at the call site. This is
why Go code has so many `if err != nil { return err }` lines: it's the
visible cost of making every failure path an explicit, reviewable branch
instead of an invisible jump.

### 10.5 Real-world example

Day 7's multi-layer error system (`ErrNotFound`, `ErrStorageFull`,
`ErrInvalidKey`, each wrapping context) is precisely the pattern a
production HTTP handler uses to map internal errors to status codes:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	err := doOperation(r)
	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidKey):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case err != nil:
		http.Error(w, "internal error", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
```

---

## 11. Packages and modules

### 11.1 Visibility is capitalization

There's no `public`/`private` keyword. An identifier starting with an
uppercase letter is exported (visible outside the package); lowercase is
package-private. This single rule replaces access modifiers entirely.

### 11.2 One package per directory, package name ≠ import path segment always

```
myapp/
  go.mod              module github.com/kumaran/myapp
  config/
    config.go          package config
  storage/
    storage.go          package storage
    memory.go           package storage  (same dir, same package)
  server/
    server.go            package server
  internal/
    auth/
      auth.go            package auth — importable ONLY from within myapp
```

`internal/` is a compiler-enforced convention: any package under a
directory named `internal` can only be imported by code rooted at the
parent of that `internal` directory. It's Go's answer to "private to
this module, but shared across several of its own packages."

### 11.3 `init()` and import cycles

Every package can define one or more `init()` functions, run automatically
before `main()`, in dependency order — used for registering drivers
(`database/sql` drivers register themselves this way) but easy to overuse
into "spooky action at a distance." Prefer explicit initialization
(a `New()` constructor) unless you're building a plugin-registration
pattern. Go also forbids import cycles at compile time — if package `a`
imports `b` and `b` imports `a`, the build fails immediately, forcing you
to extract the shared piece into a third package.

### 11.4 Real-world example

Day 9's refactor of the Week 1 project into `config/`, `storage/`,
`server/`, `health/` packages is the shape of nearly every non-trivial Go
service: each package owns one responsibility, `server` depends on
`config` and `storage` but not vice versa, and the dependency graph is a
DAG you can read straight off the import statements — no hidden coupling
through global state.

---

## 12. Strings, runes, bytes, and Unicode

### 12.1 A Go string is UTF-8 bytes, not characters

```go
s := "héllo"
fmt.Println(len(s))          // 6 — 'é' is 2 bytes in UTF-8, not 1
fmt.Println(len([]rune(s)))  // 5 — rune count matches human intuition
```

Indexing a string with `s[i]` gives you a **byte**, not a character.
`for i, r := range s` gives you rune-by-rune iteration with `i` as the
*byte* offset of each rune's start (which can skip values for multi-byte
runes) — this trips up almost everyone once.

### 12.2 Strings are immutable; conversions to `[]byte` copy

```go
b := []byte(s) // allocates and copies
s2 := string(b) // allocates and copies back
```

Because of that copy cost, hot paths that build strings incrementally
should use `strings.Builder` (amortized-growth, like `append` on a
slice) instead of repeated `+` concatenation:

```go
var sb strings.Builder
for _, part := range parts {
	sb.WriteString(part)
	sb.WriteByte('\n')
}
result := sb.String() // one final allocation, not N
```

### 12.3 The `strings`, `strconv`, and `unicode/utf8` packages

`strings.Split`, `strings.Join`, `strings.Contains`, `strings.TrimSpace`,
`strconv.Itoa`/`Atoi`, `strconv.ParseFloat`, and `utf8.RuneCountInString`
cover the overwhelming majority of text-munging you'll do. Reach for
regular expressions (`regexp`) only when a fixed-pattern parse genuinely
isn't enough — they're slower and harder to read than they look.

### 12.4 Real-world example

A CLI that accepts `--size=10MB` needs to split the numeric prefix from
the unit suffix, `strconv.ParseFloat` the number, and look up the unit in
the `iota`-based `Unit` type from Chapter 3 — three standard-library calls,
no regex needed, and it handles `"1.5GB"` just as easily as `"10MB"`.

---

## 13. Terminal project: a system-info CLI and unit converter

**Goal:** tie together Chapters 1–12 into one small, real binary you run
from your shell today.

```go
// cmd/sysinfo/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

type ByteUnit int64

const (
	Byte ByteUnit = 1
	KB            = Byte * 1024
	MB            = KB * 1024
	GB            = MB * 1024
)

func humanize(bytes uint64) string {
	switch {
	case bytes >= uint64(GB):
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= uint64(MB):
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= uint64(KB):
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func main() {
	jsonOut := flag.Bool("json", false, "print as JSON")
	flag.Parse()

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	report := struct {
		Host      string `json:"host"`
		GoVersion string `json:"go_version"`
		CPUs      int    `json:"cpus"`
		Goroutines int   `json:"goroutines"`
		HeapAlloc string `json:"heap_alloc"`
		Timestamp string `json:"timestamp"`
	}{
		Host:       host,
		GoVersion:  runtime.Version(),
		CPUs:       runtime.NumCPU(),
		Goroutines: runtime.NumGoroutine(),
		HeapAlloc:  humanize(mem.HeapAlloc),
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	if *jsonOut {
		printJSON(report)
		return
	}
	fmt.Printf("host:       %s\n", report.Host)
	fmt.Printf("go version: %s\n", report.GoVersion)
	fmt.Printf("cpus:       %d\n", report.CPUs)
	fmt.Printf("goroutines: %d\n", report.Goroutines)
	fmt.Printf("heap:       %s\n", report.HeapAlloc)
	fmt.Printf("time:       %s\n", report.Timestamp)
}
```

```bash
go run ./cmd/sysinfo
go run ./cmd/sysinfo --json | jq .
```

This is deliberately the same exercise as Day 1 of the 90-day plan, but
carried one step further with a `--json` flag — your first taste of the
`flag` package that Chapter 40 builds into full subcommand CLIs.

---

## 14. Terminal project: a JSON config server

**Goal:** a small `net/http` service — the Week-1 mini-project from the
90-day plan, built out fully here as a reference implementation you can
diff your own attempt against.

```go
// config/config.go
package config

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Store { return &Store{data: map[string]string{}} }

func (s *Store) All() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}
```

```go
// server/server.go
package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example.com/configserver/config"
)

type Server struct {
	store     *config.Store
	startedAt time.Time
}

func New(store *config.Store) *Server {
	return &Server{store: store, startedAt: time.Now()}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /config", s.handleGetAll)
	mux.HandleFunc("PUT /config/{key}", s.handleSet)
	mux.HandleFunc("GET /health", s.handleHealth)
	return mux
}

func (s *Server) handleGetAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.All())
}

func (s *Server) handleSet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if strings.TrimSpace(key) == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}
	var body struct{ Value string `json:"value"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	s.store.Set(key, body.Value)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"uptime": time.Since(s.startedAt).String(),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return
	}
}
```

```go
// main.go
package main

import (
	"log"
	"net/http"

	"example.com/configserver/config"
	"example.com/configserver/server"
)

func main() {
	store := config.New()
	srv := server.New(store)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv.Routes()))
}
```

```bash
curl -X PUT localhost:8080/config/feature_x -d '{"value":"on"}'
curl localhost:8080/config
curl localhost:8080/health
```

Notice the shape: `config` knows nothing about HTTP, `server` knows
nothing about how config is stored — swap `config.Store` for a
Redis-backed implementation later and `server` doesn't change. That
separation is worth more than any framework.

---

# Part II — Intermediate: concurrency and the standard library

Part I gave you the language running on a single thread. This part is where
Go actually distinguishes itself: goroutines, channels, and the `sync`
primitives that make "just spin up a thousand of these" a reasonable
sentence, plus the standard-library tools (testing, generics, `net/http`)
you'll reach for daily. It closes with two real concurrent tools, not toy
examples — the worker pool and LRU cache in Chapters 25–26 are worth
building yourself, not just reading.

## 15. Goroutines: Go's answer to "what is a thread"

### 15.1 Goroutines are not OS threads

A goroutine starts at ~2KB of stack (growable, shrinkable, managed by the
Go runtime) versus a typical OS thread's 1–8MB fixed stack. The Go
runtime multiplexes many goroutines onto a small number of OS threads
(the GMP model, Chapter 28) — which is why spawning 100,000 goroutines is
routine in Go, while 100,000 OS threads would exhaust most systems.

```go
go doSomething()          // fire-and-forget — but you own its lifecycle now
go func() {
	doSomethingElse()
}()
```

### 15.2 The goroutine leak

Starting a goroutine is one line. *Making sure it exits* is the actual
engineering problem. A goroutine blocked forever on a channel receive
that nothing will ever send to is a leak — it holds its stack and
whatever it closed over for the life of the process:

```go
func leaky() {
	ch := make(chan int) // unbuffered, no sender registered
	go func() {
		v := <-ch // blocks forever — ch is never written to or closed
		fmt.Println(v)
	}()
	// leaky() returns; the goroutine above never will
}
```

The fix is usually the same shape: give long-lived goroutines an explicit
way to be told to stop (a `context.Context`, a `done` channel) and know
what event ends each goroutine you start.

### 15.3 `runtime.NumGoroutine()` as a diagnostic

```go
fmt.Println("active goroutines:", runtime.NumGoroutine())
```

In production, a steadily climbing goroutine count exposed on a `/debug`
metric is usually the first visible symptom of a leak, well before memory
pressure becomes an incident.

### 15.4 War story

A payments service exposed `runtime.NumGoroutine()` as a Prometheus
gauge "just in case." Eighteen months later, the number had been slowly
climbing for weeks — invisible in latency or error-rate dashboards — and
the on-call engineer only caught it because the gauge crossed an
arbitrary alert threshold. The root cause: an HTTP client call without a
timeout, hitting a downstream service that had started hanging instead of
erroring, leaving one goroutine parked per stuck request, forever. The
fix was a `context.WithTimeout` on every outbound call (Chapter 18) — a
one-line change per call site, but only findable because the metric
existed at all. **Lesson: goroutine count is a cheap, high-signal metric
almost nobody adds until after the first leak incident.**

**Applied in the series:** [Linux guide, Chapter 73](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime) measures what goroutines cost in OS threads: 20,000 waiting goroutines used 16 threads, while 100 blocked in raw syscalls pushed it to 108.


---

## 16. Channels: pipelines, `select`, and the closing convention

### 16.1 Unbuffered vs. buffered

An unbuffered channel (`make(chan int)`) synchronizes sender and
receiver — a send blocks until a receive is ready, and vice versa. A
buffered channel (`make(chan int, 10)`) lets sends proceed without a
waiting receiver, up to the buffer size, decoupling producer and
consumer speed:

```go
unbuffered := make(chan int)   // hand-off, rendezvous
buffered := make(chan int, 10) // queue with capacity 10
```

### 16.2 `select`

`select` waits on multiple channel operations at once, proceeding with
whichever is ready first (randomly, if several are ready simultaneously):

```go
select {
case v := <-ch1:
	fmt.Println("from ch1:", v)
case ch2 <- 42:
	fmt.Println("sent to ch2")
case <-time.After(2 * time.Second):
	fmt.Println("timed out")
default:
	fmt.Println("nothing ready, non-blocking")
}
```

### 16.3 The closing convention

- Only the **sender** should close a channel — never the receiver.
- Closing signals "no more values are coming"; a receive on a closed
  channel returns immediately with the zero value and `ok == false`.
- Sending on a closed channel panics; closing an already-closed channel
  panics.

```go
func producer(out chan<- int) {
	defer close(out) // sender closes when done
	for i := 0; i < 5; i++ {
		out <- i
	}
}

func consumer(in <-chan int) {
	for v := range in { // range exits cleanly when in is closed
		fmt.Println(v)
	}
}
```

Directional channel types (`chan<- int` send-only, `<-chan int`
receive-only) in function signatures are cheap, compiler-enforced
documentation of intent — use them whenever a function only needs one
direction.

### 16.4 Pipelines

Channels compose into pipelines, each stage a goroutine reading from an
input channel and writing to an output channel:

```go
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

for result := range square(generate(1, 2, 3, 4)) {
	fmt.Println(result) // 1, 4, 9, 16
}
```

### 16.5 Real-world example

Day 13's job queue — one producer, five workers, one results channel —
is the buffered-channel worker-pool pattern used in many services, from
image resizing to batch email senders. Chapter 25 builds the fuller
version with cancellation and timeouts.

---

## 17. The `sync` package: mutexes, `WaitGroup`, `Once`, atomics

### 17.1 `sync.Mutex` and `sync.RWMutex`

```go
type SafeCounter struct {
	mu sync.Mutex
	n  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n[key]++
}
```

`sync.RWMutex` allows any number of concurrent readers OR one writer
(never both) — use it when reads vastly outnumber writes, exactly as Day
14's per-endpoint request counter does.

### 17.2 `sync.WaitGroup`

```go
var wg sync.WaitGroup
for _, url := range urls {
	wg.Add(1)
	go func(u string) {
		defer wg.Done()
		check(u)
	}(url)
}
wg.Wait() // blocks until all Done() calls have happened
```

**Common bug:** calling `wg.Add(1)` *inside* the goroutine instead of
before `go` starts it — this creates a race where `Wait()` might return
before all goroutines have even registered themselves.

### 17.3 `sync.Once`

```go
var once sync.Once
var instance *Config

func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig() // guaranteed to run exactly once, ever
	})
	return instance
}
```

### 17.4 The `sync/atomic` package

For simple counters, atomics avoid mutex overhead entirely:

```go
var requests atomic.Int64
requests.Add(1)
current := requests.Load()
```

### 17.5 Run with `-race` regularly

```bash
go test -race ./...
go run -race main.go
```

The race detector instruments memory accesses and catches races that
actually execute during that run. It cannot prove a program is race-free,
but running it in CI and on representative workloads finds many bugs that
otherwise show up only under load.

### 17.6 Real-world example

Day 14's concurrent request counter, hammered by 100 goroutines and run
under `-race`, is a tiny reproduction of the bug class behind "works on my
machine, corrupts data under load in prod" — and it's why many Go teams
run at least part of their test suite with `-race` before merging.

**Applied in the series:** the TCP/IP guide's [`lbdrain`](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries) and the HTTPS guide's [`resilience`](../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers) labs use `atomic.Bool`, `atomic.Uint64`, and mutex-guarded maps in real servers.


---

## 18. `context`: cancellation, deadlines, and request-scoped values

### 18.1 Why context exists

A single incoming HTTP request might fan out to a database query, two
downstream API calls, and a cache lookup. If the client disconnects, or a
deadline expires, every one of those in-flight operations needs to hear
about it *now* — not run to completion pointlessly. `context.Context` is
the value threaded through every layer of a call stack to carry exactly
that signal.

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel() // ALWAYS defer cancel, even if you expect the timeout to fire first

req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
if err != nil {
	return err
}
resp, err := http.DefaultClient.Do(req)
if err != nil {
	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("request timed out")
	}
	return err
}
defer resp.Body.Close()
return nil
```

### 18.2 The four context constructors

| Constructor | Cancels when |
|---|---|
| `context.Background()` | Never — the root of every context tree |
| `context.WithCancel(parent)` | You call the returned `cancel()` |
| `context.WithTimeout(parent, d)` | Duration `d` elapses, or manual cancel |
| `context.WithDeadline(parent, t)` | Wall-clock time `t` is reached |

Cancellation propagates **downward only**: cancelling a parent cancels
every context derived from it; a child can never cancel its parent.

### 18.3 Respecting cancellation in your own code

Context does nothing by itself — code has to actually check `ctx.Done()`:

```go
func worker(ctx context.Context, jobs <-chan Job) {
	for {
		select {
		case <-ctx.Done():
			return // stop working, ctx.Err() explains why
		case job, ok := <-jobs:
			if !ok {
				return
			}
			process(job)
		}
	}
}
```

### 18.4 Context values — use sparingly

`context.WithValue` carries request-scoped metadata (a request ID, an
authenticated user) across API boundaries that can't take an explicit
parameter. It is **not** a substitute for passing real parameters — if a
function needs a value to do its job, put it in the signature; reserve
context values for cross-cutting concerns like tracing IDs.

```go
type requestIDKey struct{}

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func requestIDFrom(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}
```

Using an unexported struct type as the key (instead of a `string`)
prevents collisions with context keys defined in other packages.

### 18.5 Real-world example

Week 2's concurrent URL checker cancels all remaining checks if the total
job takes longer than 30 seconds — that's `context.WithTimeout` wrapping
the whole batch, passed down into every worker's individual HTTP call, so
one shared deadline governs a fan-out of N concurrent requests without any
manual timer bookkeeping.

**Applied in the series:** `context` cancels dials and probes in the TCP/IP guide's [`dialcheck`](../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake) and stops streaming handlers on disconnect in the HTTPS guide's [SSE lab](../v2-https/real-life-guide-v1.md#chapter-24-streaming-over-http-server-sent-events-websockets-and-long-lived-connections).


---

## 19. Testing in Go: table-driven tests and subtests

### 19.1 The convention

Test files end in `_test.go`, sit beside the code they test, and test
functions take `*testing.T`:

```go
// storage_test.go
package storage

import "testing"

func TestSet(t *testing.T) {
	s := New()
	s.Set("key", "value")
	got, ok := s.Get("key")
	if !ok || got != "value" {
		t.Fatalf("got (%q, %v), want (%q, true)", got, ok, "value")
	}
}
```

### 19.2 Table-driven tests

The idiomatic way to cover many cases without duplicating test logic:

```go
func TestParseSize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"plain bytes", "100", 100, false},
		{"kilobytes", "10KB", 10 * 1024, false},
		{"empty input", "", 0, true},
		{"garbage unit", "10XB", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
```

`t.Run` gives each case its own named subtest — failures report exactly
which case broke (`TestParseSize/garbage_unit`), and `-run
TestParseSize/kilo` re-runs just that one.

### 19.3 Test helpers, `t.Helper()`, and fixtures

```go
func newTestStore(t *testing.T) *Store {
	t.Helper() // failures blame the caller's line, not this one
	s := New()
	t.Cleanup(func() { s.Close() }) // runs after the test, even on failure
	return s
}
```

### 19.4 Real-world example

Day 10's table-driven tests for the storage interface (empty key,
oversized value, concurrent reads, deleting a missing key) is exactly the
shape of a real PR's test suite — one table, one assertion block, four
edge cases documented as data instead of four near-duplicate test
functions.

---

## 20. Benchmarks and a first look at `pprof`

### 20.1 Writing a benchmark

```go
func BenchmarkMapLookup(b *testing.B) {
	m := make(map[int]bool, 10000)
	for i := 0; i < 10000; i++ {
		m[i] = true
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[i%10000]
	}
}
```

```bash
go test -bench=. -benchmem ./...
# BenchmarkMapLookup-8   50000000   23.1 ns/op   0 B/op   0 allocs/op
```

`b.N` is chosen automatically by the testing framework, running long
enough for a stable measurement. `-benchmem` adds allocation counts —
often the more actionable number, since allocations drive GC pressure
more than raw CPU time in typical services.

### 20.2 Comparing implementations

Day 11's exercise — linear scan vs. map lookup vs. sorted-slice binary
search — is the right way to answer "which is faster" with data instead
of intuition. The honest answer is usually: map wins past a few dozen
items; linear scan wins below that because of map's per-lookup hashing
overhead and worse cache locality.

### 20.3 `pprof`, the five-minute version

```go
import _ "net/http/pprof" // registers /debug/pprof/* handlers as a side effect

go func() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30   # CPU
go tool pprof http://localhost:6060/debug/pprof/heap                  # memory
(pprof) top10
(pprof) web   # opens an SVG call graph, requires graphviz
```

Chapter 37 goes deeper; for now, know that `pprof` is *always one import
away* in Go — no separate profiler installation, no sampling agent to
configure.

### 20.4 Real-world example

"It's slow" is not actionable. "CPU profile shows 40% of time in
`regexp.MustCompile` because it's being called per-request instead of
once at package init" is a one-line fix, and it's the kind of finding
`pprof top10` surfaces in under two minutes against a running service.

---

## 21. Generics: type parameters and constraints

### 21.1 The problem generics solve

Before Go 1.18, a generic-shaped function like "sum any slice of numbers"
required either code generation, `interface{}` plus reflection (slow,
unsafe at compile time), or a copy-pasted function per type. Generics add
compile-time-checked type parameters:

```go
func Sum[T int | int64 | float64](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

Sum([]int{1, 2, 3})          // T inferred as int
Sum([]float64{1.5, 2.5})     // T inferred as float64
```

### 21.2 Constraints

`T int | int64 | float64` is an inline **type set** constraint. The
standard library's `golang.org/x/exp/constraints` (and the built-in
`cmp.Ordered` in modern Go) name common ones:

```go
import "cmp"

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
```

`any` is the loosest constraint (no operations allowed beyond assignment
and comparison-to-nil for interface types); a constraint like
`comparable` permits `==`/`!=`, which is what makes a generic map key
type possible.

### 21.3 A generic data structure

```go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}

intStack := &Stack[int]{}
intStack.Push(1)
intStack.Push(2)
v, _ := intStack.Pop() // v == 2
```

### 21.4 When *not* to reach for generics

If an interface with one or two methods already expresses what you need
(Chapter 9), prefer it — generics shine for **data structures and
algorithms** parameterized purely by type (a generic `Stack[T]`, a
generic `Map`/`Filter`/`Reduce` over slices), not as a default replacement
for interfaces or as a way to avoid writing two similar 5-line functions.

### 21.5 Real-world example

A generic LRU cache — `Cache[K comparable, V any]` — is strictly better
than the `map[string]*Node` version from
[`projects/lru_cache.md`](./projects/lru_cache.md): the same code now
works for `Cache[int, User]` or `Cache[string, []byte]` with zero
duplication and full compile-time type safety. Chapter 26 builds this
version.

---

## 22. I/O: `io.Reader`/`Writer`, `bufio`, and `os`

### 22.1 The two interfaces underneath everything

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

Files, network connections, in-memory buffers, gzip streams, TLS
connections, and `os.Stdin`/`os.Stdout` all implement these two
interfaces — which is why `io.Copy(dst, src)` works identically no matter
what `dst` and `src` actually are.

### 22.2 `bufio` for line-oriented and buffered work

```go
f, err := os.Open("access.log")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

scanner := bufio.NewScanner(f)
for scanner.Scan() {
	line := scanner.Text()
	process(line)
}
if err := scanner.Err(); err != nil {
	log.Fatal(err)
}
```

`bufio.Scanner`'s default buffer caps at 64KB per line — for genuinely
huge lines, use `bufio.Reader.ReadString('\n')` or raise
`scanner.Buffer(...)` explicitly.

### 22.3 Streaming instead of loading whole files

```go
// BAD for a multi-GB file: os.ReadFile loads it all into memory
data, err := os.ReadFile("huge.log")
if err != nil {
	log.Fatal(err)
}
_ = data

// GOOD: stream line by line, constant memory regardless of file size
f, err := os.Open("huge.log")
if err != nil {
	log.Fatal(err)
}
defer f.Close()
scanner := bufio.NewScanner(f)
for scanner.Scan() {
	handleLine(scanner.Text())
}
if err := scanner.Err(); err != nil {
	log.Fatal(err)
}
```

### 22.4 Real-world example

Chapter 42's log analyzer processes multi-gigabyte files in constant
memory using exactly this streaming pattern — the difference between
`os.ReadFile` and `bufio.Scanner` is the difference between a tool that
OOMs on a large input and one that doesn't care how large the input is.

**Applied in the series:** the TCP/IP guide's [`pcapread`](../networking/tcp-ip/real-life-guide-v1.md#chapter-32-wireshark-reading-a-conversation) decodes binary capture files with `bufio`, `io.ReadFull`, and `encoding/binary`; the Linux guide's [`atomicwrite`](../os-linux/real-life-os-guide.md#chapter-77-files-that-survive-crashes-page-cache-fsync-and-atomic-replacement) makes file writes durable and atomic.


---

## 23. `encoding/json` and reflection, briefly

### 23.1 Marshal and Unmarshal

```go
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

data, err := json.Marshal(User{ID: 1, Name: "Alice"})
if err != nil {
	log.Fatal(err)
}
// {"id":1,"name":"Alice"}  -- Email omitted because omitempty + zero value

var u User
if err := json.Unmarshal(data, &u); err != nil { // pointer required, Unmarshal writes into it
	log.Fatal(err)
}
```

### 23.2 Streaming encode/decode for large payloads

```go
if err := json.NewEncoder(w).Encode(response); err != nil { // writes directly to the ResponseWriter
	return err
}
if err := json.NewDecoder(r.Body).Decode(&request); err != nil { // reads directly from the request body
	return err
}
```

Streaming avoids materializing the whole JSON document as a `[]byte`
first — meaningful for large request/response bodies.

### 23.3 Why `Marshal`/`Unmarshal` need reflection

`json.Marshal` doesn't know your struct's shape at compile time — it
walks the struct's fields at runtime via the `reflect` package, reading
tags and field types. This is also why JSON (un)marshaling is
comparatively slow versus a hand-written encoder, and why hot paths in
performance-sensitive services sometimes hand-roll marshaling or use
code-generation tools (`easyjson`, `ffjson`) instead.

### 23.4 A two-line taste of `reflect`

```go
v := reflect.ValueOf(someStruct)
t := reflect.TypeOf(someStruct)
for i := 0; i < t.NumField(); i++ {
	fmt.Println(t.Field(i).Name, "=", v.Field(i).Interface())
}
```

Reflection is powerful and dangerous: it defeats compile-time type
checking, and code that leans on it heavily is harder to read and often
slower. Chapter 32 covers when it's actually the right tool.

---

## 24. `net/http` fundamentals: client, server, middleware

### 24.1 The server side

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)

srv := &http.Server{
	Addr:         ":8080",
	Handler:      mux,
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout:  120 * time.Second,
}
log.Fatal(srv.ListenAndServe())
```

**Always** set explicit timeouts on `http.Server` — the zero-value
defaults are "no timeout," which turns a slow or malicious client into a
resource leak (this is exactly the kind of gap the WAF chapter, 49,
defends against from the outside).

### 24.2 The client side

```go
client := &http.Client{Timeout: 5 * time.Second}
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
if err != nil {
	return err
}
resp, err := client.Do(req)
if err != nil {
	return err
}
defer resp.Body.Close() // ALWAYS — otherwise connections leak from the pool
```

Reuse a single `http.Client` (it's safe for concurrent use and pools
connections internally) rather than constructing one per request.

### 24.3 Middleware: the `http.Handler` wrapping pattern

```go
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func withRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

handler := withLogging(withRecover(mux))
log.Fatal(http.ListenAndServe(":8080", handler))
```

Every middleware in Go is a function that takes a `Handler` and returns a
`Handler` — no framework needed, just closures composing in order. This
exact shape reappears as the WAF middleware chain in Chapter 49.

### 24.4 Real-world example

A production Go API almost always layers: recover → request-ID injection
→ logging → auth → rate limiting → the actual mux — each layer a
one-function `Handler`-wrapping `Handler`, composed once at startup. No
runtime dependency injection container required.

**Applied in the series:** [HTTPS guide, Part 9](../v2-https/real-life-guide-v1.md#part-9-build-the-lifecycle-in-go): a production server with every timeout set ([Ch 19](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown)), a reverse proxy ([Ch 20](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling)), and phase-by-phase request tracing with `net/http/httptrace` ([Ch 17](../v2-https/real-life-guide-v1.md#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs)).


---

## 25. Terminal project: a concurrent URL health-checker (worker pool)

**Goal:** Week 2's mini-project, built out in full — a worker pool with a
bounded concurrency limit and an overall deadline.

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	URL      string
	Status   int
	Latency  time.Duration
	Err      error
}

func checkOne(ctx context.Context, client *http.Client, url string) Result {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	resp, err := client.Do(req)
	if err != nil {
		return Result{URL: url, Err: err, Latency: time.Since(start)}
	}
	defer resp.Body.Close()
	return Result{URL: url, Status: resp.StatusCode, Latency: time.Since(start)}
}

func worker(ctx context.Context, client *http.Client, urls <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urls {
		select {
		case <-ctx.Done():
			return
		default:
			results <- checkOne(ctx, client, url)
		}
	}
}

func main() {
	concurrency := 10
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	urls := make(chan string)
	results := make(chan Result)
	client := &http.Client{Timeout: 5 * time.Second}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(ctx, client, urls, results, &wg)
	}

	go func() {
		defer close(urls)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case urls <- scanner.Text():
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.Err != nil {
			fmt.Printf("%-40s ERROR %v\n", r.URL, r.Err)
			continue
		}
		fmt.Printf("%-40s %d %v\n", r.URL, r.Status, r.Latency.Round(time.Millisecond))
	}
}
```

```bash
printf "https://go.dev\nhttps://golang.org\nhttp://localhost:9999\n" | go run .
```

The pipeline is: stdin → `urls` channel → N workers → `results` channel →
main prints as results arrive, unbuffered so results stream out
immediately rather than batching at the end. The whole batch is bounded
by one shared `context.WithTimeout`.

---

## 26. Terminal project: an LRU cache library and CLI cache server

**Goal:** implement the design from
[`projects/lru_cache.md`](./projects/lru_cache.md) as generic, reusable Go
— then wrap it in a tiny TCP server so it's a real, runnable service.

```go
// lru/lru.go
package lru

import "sync"

type node[K comparable, V any] struct {
	key        K
	value      V
	prev, next *node[K, V]
}

type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*node[K, V]
	head, tail *node[K, V] // head = most recently used
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	if capacity <= 0 {
		panic("lru: capacity must be positive")
	}
	return &Cache[K, V]{
		capacity: capacity,
		items:    make(map[K]*node[K, V], capacity),
	}
}

func (c *Cache[K, V]) unlink(n *node[K, V]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
}

func (c *Cache[K, V]) pushFront(n *node[K, V]) {
	n.prev, n.next = nil, c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.unlink(n)
	c.pushFront(n)
	return n.value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n, ok := c.items[key]; ok {
		n.value = value
		c.unlink(n)
		c.pushFront(n)
		return
	}
	n := &node[K, V]{key: key, value: value}
	c.items[key] = n
	c.pushFront(n)
	if len(c.items) > c.capacity {
		lru := c.tail
		c.unlink(lru)
		delete(c.items, lru.key)
	}
}

func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
```

```go
// lru/lru_test.go
package lru

import "testing"

func TestEviction(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3) // evicts "a", the least recently used

	if _, ok := c.Get("a"); ok {
		t.Fatal("expected 'a' to be evicted")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("got (%v, %v), want (2, true)", v, ok)
	}
}

func TestGetPromotesToMostRecentlyUsed(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")     // promotes "a"; "b" is now least recently used
	c.Put("c", 3) // evicts "b", not "a"

	if _, ok := c.Get("b"); ok {
		t.Fatal("expected 'b' to be evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("expected 'a' to survive")
	}
}
```

```bash
go test -race ./lru/...
```

Now wrap it in a minimal line-protocol TCP server (`GET key`, `PUT key
value`) — the same pattern that memcached and Redis expose, minus the
distributed part:

```go
// cmd/cacheserver/main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"example.com/cacheserver/lru"
)

func main() {
	cache := lru.New[string, string](1000)
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("cache server listening on :9090")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn, cache)
	}
}

func handleConn(conn net.Conn, cache *lru.Cache[string, string]) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		switch strings.ToUpper(fields[0]) {
		case "GET":
			if len(fields) != 2 {
				fmt.Fprintln(conn, "ERR usage: GET key")
				continue
			}
			if v, ok := cache.Get(fields[1]); ok {
				fmt.Fprintln(conn, "OK", v)
			} else {
				fmt.Fprintln(conn, "MISS")
			}
		case "PUT":
			if len(fields) < 3 {
				fmt.Fprintln(conn, "ERR usage: PUT key value")
				continue
			}
			cache.Put(fields[1], strings.Join(fields[2:], " "))
			fmt.Fprintln(conn, "OK")
		default:
			fmt.Fprintln(conn, "ERR unknown command")
		}
	}
}
```

```bash
go run ./cmd/cacheserver &
nc localhost 9090
PUT session:42 alice
GET session:42
GET missing
```

Each connection gets its own goroutine (`go handleConn(...)`); the shared
`*lru.Cache` is safe under concurrent access because every method locks
internally — this is the "one goroutine per connection, shared
synchronized state" pattern that describes most simple Go TCP servers,
extended fully in Chapter 38.

---

# Part III — Advanced: systems Go

You've used goroutines and channels since Part II; this part opens up what's
actually happening underneath them — the scheduler, the memory model,
garbage collection — plus the concurrency patterns, advanced generics, and
profiling tools you reach for once a program has to perform under load, not
just work correctly. It ends with a TCP chat server and a mini API gateway,
both built on the internals covered earlier in the part.

## 27. The Go memory model and the garbage collector

### 27.1 Stack vs. heap, and escape analysis

The compiler decides, per variable, whether it can live on the stack
(fast, freed automatically when the function returns, no GC involvement)
or must "escape" to the heap (because a pointer to it outlives the
function, e.g. it's returned, or captured by a goroutine/closure that
outlives the call).

```bash
go build -gcflags="-m" main.go
# main.go:12:6: moved to heap: x
```

```go
func onStack() int {
	x := 42
	return x // value copied out — x itself never escapes
}

func onHeap() *int {
	x := 42
	return &x // address escapes — x must live on the heap
}
```

Fewer heap allocations means less GC work — this is *the* lever behind
most Go performance tuning, more so than algorithmic micro-optimization.

### 27.2 The GC: concurrent, tri-color mark-and-sweep

Go's GC runs concurrently with your program (short "stop-the-world"
pauses measured in microseconds, not the multi-millisecond pauses of
older collectors), using a tri-color mark-and-sweep algorithm with a
write barrier to track pointer mutations during marking. You almost
never need to know the algorithm's internals day to day — you need to
know its **knobs**:

```bash
GOGC=100   # default: GC triggers when heap doubles since the last collection
GOGC=off   # disable GC entirely (only for short-lived batch tools!)
GOMEMLIMIT=512MiB  # soft memory cap; GC works harder as usage approaches it
```

### 27.3 Reducing GC pressure in practice

- Pre-allocate slices with `make([]T, 0, n)` when the size is known
  (Chapter 6).
- Reuse buffers with `sync.Pool` for high-churn, short-lived allocations.
- Pass large structs by pointer instead of copying them repeatedly.
- Avoid unnecessary `interface{}`/`any` boxing of small values in hot
  loops — boxing an `int` into an interface allocates.

```go
var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func handle(w http.ResponseWriter, r *http.Request) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	// ... use buf ...
}
```

### 27.4 War story

A JSON API service's p99 latency had a periodic 40ms spike every few
seconds under load — invisible in average latency, brutal for a
low-latency SLA. A CPU profile during the spike showed `runtime.gcBgMarkWorker`
dominating. The root cause: a hot path allocated a fresh `[]byte` buffer
per request to build a response, and under load the heap grew fast
enough to trigger GC cycles every couple of seconds. Switching that one
buffer to a `sync.Pool` cut allocations in the hot path by ~90% and the
p99 spikes disappeared. **Lesson: GC pauses are a downstream symptom —
the actual bug is almost always "too many allocations," found by
profiling the allocator (`pprof`'s heap profile), not by tuning `GOGC`.**

**Applied in the series:** [Linux guide, Chapter 76](../os-linux/real-life-os-guide.md#chapter-76-memory-limits-the-go-heap-gomemlimit-and-the-oom-killer): the GC's pacing meets a container memory limit. OOM-killed at 300 MB without `GOMEMLIMIT`, survived with it, at the cost of about 1,500 GCs per second.


---

## 28. The scheduler deep dive: GMP and work-stealing

### 28.1 The three letters: G, M, P

- **G** (Goroutine): your lightweight function-in-flight, with its own
  small growable stack.
- **M** (Machine): an OS thread — the thing the operating system actually
  schedules.
- **P** (Processor): a scheduling context, one per `GOMAXPROCS`, each
  holding a local run queue of Gs ready to execute. An M must hold a P to
  run Go code.

### 28.2 Why this design

`GOMAXPROCS` (default: `runtime.NumCPU()`) caps how many Gs can execute Go
code *simultaneously*, while thousands of Gs can be *parked* waiting on
channels, I/O, or locks. When an M would block on a syscall, the runtime
detaches it from its P and hands that P to another M, so blocking one
goroutine on I/O never blocks the other goroutines scheduled on the same
logical processor.

### 28.3 Work stealing

Each P has a local run queue; when a P's local queue empties, it steals
half the Gs from another P's queue rather than going idle. This keeps
all `GOMAXPROCS` processors busy under uneven workloads without a global
lock on every scheduling decision.

### 28.4 Cooperative preemption and why it (mostly) doesn't matter to you

Since Go 1.14, the scheduler can preempt a long-running goroutine even
without a function call or channel operation providing a natural
checkpoint (async preemption via signals) — meaning a tight
`for {}` CPU-bound loop no longer starves other goroutines on the same P
the way it could in much older Go versions.

### 28.5 Real-world example

Setting `GOMAXPROCS` explicitly used to matter a lot in containers:
older Go releases defaulted mostly from the host CPU count / affinity, not
the container's CPU quota, so a container capped at 2 CPUs but running on a
64-core host could create too much runnable parallelism and hit CPU
throttling. Go 1.25 added cgroup-aware default `GOMAXPROCS` behavior on
Linux (`GODEBUG=containermaxprocs=...` controls it), and
`go.uber.org/automaxprocs` covers older deployments — still worth checking
explicitly any time you deploy Go into Kubernetes.

**Applied in the series:** [Linux guide, Chapter 75](../os-linux/real-life-os-guide.md#chapter-75-cpu-limits-gomaxprocs-cgroups-and-throttling): container-aware GOMAXPROCS (Go 1.25+) measured under cgroup CPU limits, including 13.7 s of throttling when GOMAXPROCS was forced too high.


---

## 29. Concurrency patterns: fan-in/fan-out, pipelines, rate limiting, `errgroup`

### 29.1 Fan-out / fan-in

Distribute work across multiple goroutines (fan-out), then merge their
results back into one stream (fan-in):

```go
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(channels))
	for _, c := range channels {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
```

### 29.2 Rate limiting with `time.Ticker` or `golang.org/x/time/rate`

```go
limiter := rate.NewLimiter(rate.Limit(10), 1) // 10 events/sec, burst 1
for _, req := range requests {
	if err := limiter.Wait(ctx); err != nil {
		break // ctx cancelled
	}
	send(req)
}
```

### 29.3 `errgroup`: `WaitGroup` that also propagates the first error

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(ctx)
for _, url := range urls {
	url := url
	g.Go(func() error {
		return fetch(ctx, url)
	})
}
if err := g.Wait(); err != nil {
	return err // the first error any goroutine returned
}
```

`errgroup.WithContext` also cancels the shared `ctx` the moment any
goroutine returns a non-nil error — so a failing fetch stops the others
from doing pointless work, without you wiring that cancellation by hand.

### 29.4 Real-world example

A service that fans out to five downstream APIs and needs "all succeed
or report which failed, and don't wait for slow ones once one has
failed" is `errgroup` in five lines — this exact pattern replaces pages
of manual `WaitGroup` + mutex + error-slice bookkeeping.

**Applied in the series:** the HTTPS guide's [`resilience`](../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers) lab implements rate limiting, retries with jitter and a budget, and a circuit breaker.


---

## 30. Advanced error handling: `errors.Join`, custom types, panic boundaries

### 30.1 `errors.Join`: combining independent failures

When several independent operations can each fail and you want to report
*all* of them (not just the first):

```go
var errs []error
for _, task := range tasks {
	if err := task.Run(); err != nil {
		errs = append(errs, err)
	}
}
if len(errs) > 0 {
	return errors.Join(errs...) // one error, Unwrap()-able into each part
}
```

### 30.2 `panic`/`recover`: the boundary, not the flow

Reserve `panic` for programmer errors (invariant violations you'd rather
crash loudly on than silently corrupt data over) and use `recover()` only
at a deliberate boundary — typically once, at the top of a goroutine or
an HTTP middleware — never as a substitute for returning an `error`:

```go
func safeRun(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()
	fn()
}
```

A `recover()` only works if it's in a `defer` in the **same goroutine**
that panicked — a panic in a goroutine you didn't wrap will crash the
whole process, by design. This is why Chapter 24's `withRecover`
middleware exists per-request: without it, one panicking handler takes
down the entire server.

### 30.3 Custom error types with structured fields

```go
type RateLimitError struct {
	Limit     int
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit of %d exceeded, retry after %s", e.Limit, e.RetryAfter)
}

var err error = &RateLimitError{Limit: 100, RetryAfter: 30 * time.Second}
var rle *RateLimitError
if errors.As(err, &rle) {
	w.Header().Set("Retry-After", fmt.Sprint(int(rle.RetryAfter.Seconds())))
}
```

---

## 31. Advanced generics: constraints, generic algorithms, inference limits

### 31.1 Custom constraint interfaces

```go
type Numeric interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Average[T Numeric](nums []T) float64 {
	var sum T
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums))
}
```

The `~` ("approximation element") means "this type, or any type whose
*underlying* type is this" — so a `type Celsius float64` still satisfies
`~float64`, letting your generic function work with named types, not
just the exact built-ins.

### 31.2 Generic functions over slices: `Map`, `Filter`, `Reduce`

```go
func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

func Filter[T any](s []T, pred func(T) bool) []T {
	var out []T
	for _, v := range s {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

names := Map(users, func(u User) string { return u.Name })
active := Filter(users, func(u User) bool { return u.Active })
```

The standard library's `slices` and `maps` packages (Go 1.21+) already
provide `slices.Sort`, `slices.Contains`, `maps.Keys`, and similar — reach
for those before writing your own.

### 31.3 Where type inference runs out

Generic **methods** (as opposed to generic functions/types) don't exist —
you cannot add a new type parameter to a method, only use the type
parameters already declared on the receiver type. And inference across
multiple unrelated type parameters sometimes needs an explicit hint:

```go
// Won't infer U from context alone in many expressions — spell it out:
result := Map[User, string](users, func(u User) string { return u.Name })
```

### 31.4 Real-world example

Rewriting the `Storage` interface from Chapter 9 as
`Storage[K comparable, V any]` gives you compile-time-checked,
zero-reflection generic storage backends — the exact upgrade path from
the 90-day plan's `map[string]string`-based cache to the
[`projects/lru_cache.md`](./projects/lru_cache.md) design in Chapter 26.

---

## 32. Reflection and `unsafe`: when, and why rarely

### 32.1 What reflection is actually for

Reflection (`reflect` package) is the right tool when you're writing a
*generic* library that must operate on arbitrary, caller-defined types
it cannot know at compile time: `encoding/json`, `encoding/xml`,
dependency-injection containers, ORMs, struct validators. It is the wrong
tool for application-level logic, where a type switch or generics almost
always express intent more clearly and run faster.

### 32.2 `unsafe`: escaping Go's type system entirely

`unsafe.Pointer` lets you reinterpret memory across types, bypassing the
type checker and the memory safety guarantees Go otherwise provides. It
exists for a few legitimate uses — implementing extremely low-level
data-structure tricks in the standard library itself, some cgo
interop, and zero-copy conversions like `unsafe.String`/`unsafe.Slice`
(Go 1.20+, replacing the older raw-pointer-arithmetic idiom) — and it is
explicitly exempted from Go's backward-compatibility promise. If you
find yourself reaching for `unsafe` in application code, that's a strong
signal to look for a safe alternative first.

```go
// A legitimate, narrow use: zero-copy []byte -> string when you can
// prove the []byte will never be mutated afterward.
func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
```

### 32.3 Real-world example

`protobuf`, `gRPC`, and high-throughput serialization libraries use
`unsafe` internally to avoid copying bytes between a wire buffer and a Go
string — but they hide it behind a safe public API. That's the correct
shape: `unsafe` stays an implementation detail of a well-tested library,
never a tool reached for casually in a service handler.

---

## 33. `cgo` and syscalls, briefly

### 33.1 `cgo`: calling C from Go

```go
/*
#include <stdlib.h>
*/
import "C"

func randomC() int {
	return int(C.rand())
}
```

`cgo` lets Go call into existing C libraries, at real costs: it disables
cross-compilation convenience (you now need a C toolchain for the
target), it's slower per-call than a pure-Go function call (a few dozen
nanoseconds of overhead per cgo call, from stack-switching), and it
complicates static linking. Use it only when a mature C library
genuinely has no Go equivalent worth the switch (common with certain
crypto, image codecs, or hardware SDKs).

### 33.2 Direct syscalls without `cgo`

The `syscall` and `golang.org/x/sys/unix` packages let you call the OS
directly in pure Go, preserving cross-compilation and static linking —
this is how Go implements `os.Open`, `net.Listen`, etc. under the hood,
and it's the preferred route over `cgo` whenever the functionality is
just "talk to the kernel," not "use this specific C library's algorithm."

### 33.3 Real-world example

Docker and Kubernetes are written in pure Go specifically so they remain
statically-linked, cross-compilable, dependency-free binaries — both
avoid `cgo` in their core paths for exactly this reason, reaching for
`golang.org/x/sys/unix` when they need raw namespace/cgroup syscalls.

**Applied in the series:** the TCP/IP guide's [`tcpinfo`](../networking/tcp-ip/real-life-guide-v1.md#chapter-33-ss-what-your-machine-s-connections-are-doing) reads `TCP_INFO` with `SyscallConn().Control`; the Linux guide's [`minicontainer`](../os-linux/real-life-os-guide.md#chapter-80-a-container-runtime-in-150-lines-of-go) uses `SysProcAttr` clone flags and `CgroupFD`.


---

## 34. Networking deep dive: `net.Conn`, TCP/UDP, framing your own protocol

### 34.1 `net.Conn`: the interface underneath every connection

```go
type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	SetDeadline(t time.Time) error
	// ... plus local/remote address accessors
}
```

TCP connections, TLS connections (Chapter 35), and Unix domain sockets
all implement `net.Conn` — protocol-level code written against the
interface works over any of them unchanged.

### 34.2 TCP server skeleton

```go
ln, err := net.Listen("tcp", ":9000")
if err != nil {
	log.Fatal(err)
}
for {
	conn, err := ln.Accept()
	if err != nil {
		continue
	}
	go handle(conn) // one goroutine per connection — cheap, Chapter 15
}
```

### 34.3 UDP: connectionless, message-oriented

```go
addr, err := net.ResolveUDPAddr("udp", ":9001")
if err != nil {
	log.Fatal(err)
}
conn, err := net.ListenUDP("udp", addr)
if err != nil {
	log.Fatal(err)
}
buf := make([]byte, 1500) // stay under typical MTU to avoid IP fragmentation
for {
	n, remote, err := conn.ReadFromUDP(buf)
	if err != nil {
		continue
	}
	packet := append([]byte(nil), buf[:n]...) // copy before reusing buf in the next loop
	go handleDatagram(conn, remote, packet)
}
```

TCP guarantees ordered, reliable delivery over a byte **stream** with no
built-in message boundaries; UDP delivers discrete **datagrams** with no
ordering or delivery guarantee at all. Choosing between them is choosing
whether your application layer needs to reinvent ordering/reliability
(UDP-based protocols like QUIC do) or can lean on TCP's.

### 34.4 TCP has no message boundaries — you must frame your own

A `Read` on a TCP `net.Conn` can return a partial message, multiple
messages concatenated, or anything in between — TCP only guarantees byte
order, not that one `Write` maps to one `Read`. Real protocols solve this
with either delimiters (`bufio.Scanner` with `\n`, as Chapter 26's cache
server does) or **length-prefixing**:

```go
func writeFramed(w io.Writer, payload []byte) error {
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if _, err := w.Write(header[:]); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func readFramed(r io.Reader) ([]byte, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(header[:])
	payload := make([]byte, size)
	_, err := io.ReadFull(r, payload)
	return payload, err
}
```

`io.ReadFull` is the detail that matters here: a plain `Read` can return
fewer bytes than requested even when more are coming — `ReadFull` loops
until it either fills the buffer or hits an error.

### 34.5 Real-world example

Redis's RESP protocol and HTTP/1.1's `Content-Length` header are both
solving exactly this framing problem at the application layer, because
TCP itself gives you none of it. Chapter 38's chat server uses
newline-delimited framing (simple, human-debuggable with `nc`); a binary
protocol handling arbitrary payloads (including payloads containing raw
newline bytes) needs length-prefixing instead.

**Applied in the series:** the TCP/IP guide's [Go labs](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself) use `net.Conn` for real protocols: a [byte-stream framing demo](../networking/tcp-ip/real-life-guide-v1.md#chapter-22-tcp-part-2-how-it-never-loses-your-data), a [DNS client](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses), [raw HTTP/1.1](../networking/tcp-ip/real-life-guide-v1.md#chapter-28-http-how-the-web-actually-talks), and the [PROXY protocol](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries).


---

## 35. TLS in Go: `crypto/tls`, mutual TLS

### 35.1 A TLS server with `net/http`

```go
srv := &http.Server{Addr: ":8443", Handler: mux}
log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
```

### 35.2 A TLS client that pins expectations explicitly

```go
transport := &http.Transport{
	TLSClientConfig: &tls.Config{
		MinVersion: tls.VersionTLS12,
		// InsecureSkipVerify: true  <-- NEVER in production; disables all
		// hostname/chain verification and defeats the point of TLS.
	},
}
client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
```

### 35.3 Generating a local development certificate

```bash
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem \
  -days 365 -nodes -subj "/CN=localhost"
```

Chapter 50 builds full mutual TLS (client certificates, not just server
certificates) for a real auth use case.

**Applied in the series:** [`tlsinspect`](../networking/tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted) diagnoses certificate failures, [`certreload`](../networking/tcp-ip/real-life-guide-v1.md#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation) rotates certificates with `GetCertificate`, and [`certcheck`](../v2-https/real-life-guide-v1.md#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation) audits a certificate inventory.


---

## 36. Structured logging and observability: `log/slog`, metrics, tracing

### 36.1 `log/slog` (standard library since Go 1.21)

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("request handled",
	"method", r.Method,
	"path", r.URL.Path,
	"status", status,
	"latency_ms", latency.Milliseconds(),
)
// {"time":"...","level":"INFO","msg":"request handled","method":"GET",...}
```

Structured (key-value, machine-parseable) logs replace `fmt.Printf`-style
free text specifically so log aggregators (Loki, Elasticsearch,
CloudWatch Insights) can filter and query on fields instead of grepping
strings.

### 36.2 Metrics: counters, gauges, histograms

```go
import "github.com/prometheus/client_golang/prometheus"

var requestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{Name: "http_request_duration_seconds"},
	[]string{"method", "path", "status"},
)

requestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(latency.Seconds())
```

### 36.3 Tracing: propagating a request ID/span across services

Combine the `context.WithValue` request-ID pattern from Chapter 18 with
`log/slog`'s attribute passing, and every log line for a given request —
across every service it touches — carries the same ID, turning "find
every log line for this one failed request across five microservices"
from a manual grep exercise into a single query.

### 36.4 Real-world example

The single highest-leverage observability change most young Go services
are missing is: structured logs with a request ID, `runtime.NumGoroutine()`
as a gauge (Chapter 15's war story), and `net/http/pprof` mounted on an
internal-only port (Chapter 20) — three cheap additions that turn "why is
prod slow" from a guessing game into a five-minute diagnosis.

**Applied in the series:** the HTTPS guide's [production server](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) logs every request with `log/slog`, and its [SLO chapter](../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services) turns request outcomes into burn-rate alerts.


---

## 37. Performance tuning: `pprof`, `trace`, and benchmark methodology

### 37.1 The four profile types that matter most

| Profile | Endpoint / flag | Answers |
|---|---|---|
| CPU | `/debug/pprof/profile?seconds=30` | Where is CPU time going? |
| Heap | `/debug/pprof/heap` | What's allocating memory, and how much is live? |
| Goroutine | `/debug/pprof/goroutine` | How many goroutines, and where are they stuck? |
| Block/Mutex | `/debug/pprof/block`, `/mutex` | Where is contention (lock waits) happening? |

### 37.2 `go tool trace`: the scheduler's-eye view

```bash
go test -trace=trace.out ./...
go tool trace trace.out
```

Where `pprof` aggregates "how much time in this function," `trace`
shows a literal timeline of every goroutine's state transitions
(running, runnable, blocked, GC) — the right tool when the question is
"why did this specific request stall for 200ms" rather than "what's
hot."

### 37.3 Benchmark methodology that doesn't lie to you

- Run with `-benchtime=3x` or longer for stable numbers; a single quick
  run on a laptop with other processes competing for CPU is noise.
- Always include `-benchmem`; allocation count is often the more
  actionable number than raw nanoseconds.
- Use `benchstat` (`golang.org/x/perf/cmd/benchstat`) to compare two runs
  statistically instead of eyeballing a percentage difference that might
  be within noise.

```bash
go test -bench=. -benchmem -count=10 ./... > old.txt
# make your change
go test -bench=. -benchmem -count=10 ./... > new.txt
benchstat old.txt new.txt
```

### 37.4 Real-world example

"We think the new caching layer made things faster" is a hypothesis, not
a fact, until `benchstat` on ten runs each of old vs. new shows a
statistically significant delta rather than two numbers that happen to
differ once.

**Applied in the series:** [Linux guide, Chapter 81](../os-linux/real-life-os-guide.md#chapter-81-profiling-go-on-linux-pprof-the-execution-tracer-and-perf): two planted bugs found with pprof. A regexp compiled per request turned out to be 86% of the handler's CPU, and a goroutine leak showed up as 489 goroutines stuck on one line.


---

## 38. Terminal project: a TCP chat server with rooms

**Goal:** apply Chapters 15–17, 27–29, and 34 to a real multi-client,
concurrent, stateful server — the natural "advanced" successor to the
cache server in Chapter 26.

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Room struct {
	mu      sync.Mutex
	clients map[string]net.Conn
}

func NewRoom() *Room { return &Room{clients: map[string]net.Conn{}} }

func (r *Room) Join(name string, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[name] = conn
	r.broadcastLocked(name, "* "+name+" joined")
}

func (r *Room) Leave(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, name)
	r.broadcastLocked(name, "* "+name+" left")
}

func (r *Room) Broadcast(from, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.broadcastLocked(from, msg)
}

// caller must hold r.mu
func (r *Room) broadcastLocked(from, msg string) {
	line := fmt.Sprintf("[%s] %s\n", from, msg)
	for name, conn := range r.clients {
		if name == from {
			continue
		}
		fmt.Fprint(conn, line)
	}
}

func handleClient(conn net.Conn, room *Room) {
	defer conn.Close()
	fmt.Fprint(conn, "enter your name: ")
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		return
	}

	room.Join(name, conn)
	defer room.Leave(name)

	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "" {
			continue
		}
		room.Broadcast(name, msg)
	}
}

func main() {
	room := NewRoom()
	ln, err := net.Listen("tcp", ":9100")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("chat server listening on :9100")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleClient(conn, room)
	}
}
```

```bash
go run . &
nc localhost 9100   # open in two terminals, pick different names, chat
```

This uses newline-delimited framing (Chapter 34) rather than
length-prefixing because it's a human-typed, line-oriented protocol —
`nc` is a perfectly good client. Every shared piece of state (the
`clients` map) is behind one mutex; every connection's read loop runs in
its own goroutine.

---

## 39. Terminal project: a rate-limited reverse proxy (mini API gateway)

**Goal:** combine `net/http/httputil.ReverseProxy`, the middleware chain
from Chapter 24, and rate limiting from Chapter 29 into a small gateway —
directly in the spirit of the networking/security field guides' WAF and
API-protection material.

```go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"golang.org/x/time/rate"
)

type ipLimiters struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
}

func newIPLimiters() *ipLimiters {
	return &ipLimiters{limiters: map[string]*rate.Limiter{}}
}

func (l *ipLimiters) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(rate.Limit(5), 10) // 5 req/s, burst 10, per IP
		l.limiters[ip] = lim
	}
	return lim
}

func withRateLimit(limiters *ipLimiters, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !limiters.get(ip).Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd // production: parse + validate against a trusted proxy list
	}
	return r.RemoteAddr
}

func main() {
	target, _ := url.Parse("http://localhost:8080") // your upstream service
	proxy := httputil.NewSingleHostReverseProxy(target)

	limiters := newIPLimiters()
	handler := withRateLimit(limiters, proxy)

	log.Println("gateway listening on :9443, proxying to", target)
	log.Fatal(http.ListenAndServe(":9443", handler))
}
```

`httputil.ReverseProxy` rewrites the request's `Host`/URL, forwards it to
`target`, and streams the response back — the rate limiter middleware
runs first and short-circuits with `429 Too Many Requests` before the
upstream ever sees an over-quota request. Chapter 49 extends this exact
shape into a fuller WAF-style middleware with IP blocklists and request
inspection.

**Applied in the series:** the [HTTPS guide's reverse proxy chapter](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling) covers `Rewrite` vs `Director` and `X-Forwarded-For` trust, and reproduces request smuggling's root cause.


---

# Part IV — Real-world builds: terminal

Theory and internals are done — Parts I–III gave you everything the next
three parts use. This one and the two after it are capstones: full,
opinionated builds rather than isolated examples. This one stays on the
terminal: a production-grade CLI, then two larger capstones (a job queue, a
multi-GB log analyzer) that combine concurrency, I/O, and error handling
into something closer to what you'd actually ship.

## 40. Building production CLIs: `flag`, subcommands, config precedence

### 40.1 The standard library's `flag` package

```go
var (
	verbose = flag.Bool("verbose", false, "enable verbose logging")
	port    = flag.Int("port", 8080, "port to listen on")
	config  = flag.String("config", "", "path to config file")
)

func main() {
	flag.Parse()
	args := flag.Args() // positional args after all flags
}
```

### 40.2 Subcommands without a framework

`flag.NewFlagSet` gives each subcommand its own independent flag scope —
this is literally how `cobra`-style CLIs (`git commit -m`, `docker run
-it`) are built underneath, and it's often enough on its own for an
internal tool:

```go
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: toolbelt <serve|migrate|version>")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.Int("port", 8080, "port")
		serveCmd.Parse(os.Args[2:])
		runServe(*port)
	case "migrate":
		migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
		dryRun := migrateCmd.Bool("dry-run", false, "preview only")
		migrateCmd.Parse(os.Args[2:])
		runMigrate(*dryRun)
	case "version":
		fmt.Println("toolbelt v1.0.0")
	default:
		fmt.Printf("unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}
```

For richer needs — nested subcommands, auto-generated help, shell
completion — reach for `github.com/spf13/cobra`, which is `flag.FlagSet`
composed at scale, not a different paradigm.

### 40.3 Config precedence: flags > env > file > defaults

The convention nearly every well-built CLI follows, highest priority
first:

```go
func resolvePort(flagVal *int, flagSet bool) int {
	if flagSet {
		return *flagVal // 1. explicit flag wins
	}
	if v := os.Getenv("APP_PORT"); v != "" { // 2. environment variable
		if p, err := strconv.Atoi(v); err == nil {
			return p
		}
	}
	if cfg.Port != 0 { // 3. config file value
		return cfg.Port
	}
	return 8080 // 4. hardcoded default
}
```

### 40.4 Exit codes and stderr vs. stdout

Reserve stdout for a program's actual output (so it composes with `|` and
scripting); send diagnostics, warnings, and errors to stderr; and return
a non-zero exit code (`os.Exit(1)`) on failure so shell scripts and CI
can detect it:

```go
if err := run(); err != nil {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
```

---

## 41. Capstone: a distributed-style job queue with persistence

**Goal:** combine channels, workers, `context`, `encoding/json`, and file
persistence into a small but real job queue — the kind of tool that
sits behind "process this batch of uploaded files" or "send this batch
of emails" in a real backend.

```go
// queue/job.go
package queue

type Status string

const (
	Pending   Status = "pending"
	Running   Status = "running"
	Succeeded Status = "succeeded"
	Failed    Status = "failed"
)

type Job struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
	Status  Status `json:"status"`
	Error   string `json:"error,omitempty"`
}
```

```go
// queue/queue.go
package queue

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

type Queue struct {
	mu       sync.Mutex
	jobs     map[string]*Job
	pending  chan *Job
	statePath string
}

func New(statePath string, bufferSize int) *Queue {
	q := &Queue{
		jobs:      map[string]*Job{},
		pending:   make(chan *Job, bufferSize),
		statePath: statePath,
	}
	_ = q.load() // teaching example: production code should return/log this error
	return q
}

func (q *Queue) Enqueue(id, payload string) *Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	j := &Job{ID: id, Payload: payload, Status: Pending}
	q.jobs[id] = j
	q.pending <- j
	_ = q.persistLocked()
	return j
}

func (q *Queue) Get(id string) (*Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	j, ok := q.jobs[id]
	return j, ok
}

func (q *Queue) Run(ctx context.Context, workers int, handler func(string) error) {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-q.pending:
					if !ok {
						return
					}
					q.process(j, handler)
				}
			}
		}()
	}
	wg.Wait()
}

func (q *Queue) process(j *Job, handler func(string) error) {
	q.setStatus(j, Running, "")
	if err := handler(j.Payload); err != nil {
		q.setStatus(j, Failed, err.Error())
		return
	}
	q.setStatus(j, Succeeded, "")
}

func (q *Queue) setStatus(j *Job, status Status, errMsg string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	j.Status = status
	j.Error = errMsg
	_ = q.persistLocked()
}

// caller must hold q.mu
func (q *Queue) persistLocked() error {
	f, err := os.Create(q.statePath)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(q.jobs); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func (q *Queue) load() error {
	f, err := os.Open(q.statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no state file yet — fresh start
		}
		return err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&q.jobs); err != nil {
		return err
	}
	for _, j := range q.jobs {
		if j.Status == Pending || j.Status == Running {
			j.Status = Pending
			q.pending <- j // re-enqueue anything that didn't finish last run
		}
	}
	return nil
}
```

```go
// cmd/jobqueue/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"example.com/jobqueue/queue"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	q := queue.New("jobs.json", 100)

	for i := 0; i < 5; i++ {
		q.Enqueue(fmt.Sprintf("job-%d", i), fmt.Sprintf("payload-%d", i))
	}

	q.Run(ctx, 3, func(payload string) error {
		log.Println("processing", payload)
		time.Sleep(500 * time.Millisecond)
		return nil
	})
}
```

This gives you: a buffered channel as the in-memory work queue,
persistence to disk on every state change (so a crash doesn't silently
lose which jobs were pending), and `signal.NotifyContext` (Go 1.16+) to
turn Ctrl-C into a clean, cancellable shutdown instead of an abrupt kill.
Scaling this to *actually* distributed (multiple processes, not just
multiple goroutines) means swapping the in-memory channel for a shared
backend — Redis, PostgreSQL with `SELECT ... FOR UPDATE SKIP LOCKED`, or a
message broker — while keeping the exact same `Job`/`Status` model.

---

## 42. Capstone: a streaming log analyzer for multi-GB files

**Goal:** parse and aggregate a large access log in constant memory,
concurrently, using everything from Chapters 22, 27, and 29.

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"sync"
)

var logLineRe = regexp.MustCompile(`^(\S+) \S+ \S+ \[.*?\] "(\S+) (\S+).*?" (\d{3}) (\d+)`)

type stats struct {
	mu         sync.Mutex
	byStatus   map[string]int
	byPath     map[string]int
	totalBytes int64
}

func newStats() *stats {
	return &stats{byStatus: map[string]int{}, byPath: map[string]int{}}
}

func (s *stats) record(status, path string, bytes int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byStatus[status]++
	s.byPath[path]++
	s.totalBytes += bytes
}

func processLine(line string, s *stats) {
	m := logLineRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	status := m[4]
	path := m[3]
	var bytes int64
	fmt.Sscanf(m[5], "%d", &bytes)
	s.record(status, path, bytes)
}

func main() {
	path := flag.String("file", "", "log file to analyze")
	workers := flag.Int("workers", 4, "concurrent parsing workers")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer f.Close()

	s := newStats()
	lines := make(chan string, 1000)

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				processLine(line, s)
			}
		}()
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // handle long lines
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)
	wg.Wait()

	printReport(s)
}

func printReport(s *stats) {
	fmt.Println("=== Status codes ===")
	for status, count := range sortedByCount(s.byStatus) {
		fmt.Printf("%-6s %d\n", status, count)
	}
	fmt.Println("\n=== Top paths ===")
	top := sortedByCount(s.byPath)
	i := 0
	for path, count := range top {
		if i >= 10 {
			break
		}
		fmt.Printf("%-40s %d\n", path, count)
		i++
	}
	fmt.Printf("\nTotal bytes served: %d\n", s.totalBytes)
}

func sortedByCount(m map[string]int) map[string]int {
	type kv struct{ k string; v int }
	pairs := make([]kv, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].v > pairs[j].v })
	out := make(map[string]int, len(pairs))
	for _, p := range pairs {
		out[p.k] = p.v
	}
	return out
}
```

```bash
go run . -file access.log -workers 8
```

The file is read line-by-line (never fully loaded), lines fan out across
worker goroutines through a buffered channel, and a single mutex guards
the shared aggregation maps — the same architecture as Chapter 25's URL
checker, applied to a CPU-bound (regex parsing) workload instead of an
I/O-bound one. (`sortedByCount`'s Go-map-preserves-insertion-illusion is
sidestepped by returning a slice-derived ordering for display, not by
relying on map order — see Chapter 7.3.)

---

# Part V — Real-world builds: GUI

Same capstone spirit as Part IV, aimed at the desktop instead of the
terminal. Go has no built-in GUI toolkit — a deliberate omission — so this
part is as much about evaluating the ecosystem (Fyne, Gio, Wails, Walk) as it
is about building with one; the two projects that follow use different
approaches (native Fyne, web-based Wails) so you can feel the trade-off
directly instead of taking it on faith.

## 43. GUI options in Go: Fyne, Gio, Wails, Walk, and "just serve HTML"

Go has no built-in GUI toolkit — this is a deliberate omission consistent
with the language's "small core, do one thing well" philosophy — but a
mature ecosystem fills the gap, with genuinely different trade-offs:

| Toolkit | Model | Best for |
|---|---|---|
| **Fyne** | Pure Go, custom-drawn widgets, cross-platform (desktop + mobile) | Native-feeling cross-platform desktop apps, single binary |
| **Gio** | Pure Go, immediate-mode rendering, GPU-accelerated | Highly custom UI, animations, mobile + desktop from one codebase |
| **Wails** | Go backend + a real web frontend (HTML/CSS/JS/React/Vue) in a native window | Teams that already know web frontend and want Electron-like DX without Electron's footprint |
| **Walk** | Thin wrapper over native Win32 widgets | Windows-only apps that should feel native on Windows |
| "Serve HTML" | `net/http` + `html/template`, opened in the user's normal browser | Zero-install internal tools, dashboards, anything where "it's a webpage" is fine |

### 43.1 The trade-off that decides most of these choices

Fyne and Gio give you one Go codebase and no embedded browser — often
smaller distributables, though GUI builds may still depend on platform
libraries. You're drawing your own widgets (Fyne looks clean and "good
enough"; pixel-perfect custom design leans toward Gio). Wails gives you
the full power of modern CSS/JS and any frontend framework you already
know, at the cost of shipping a WebView-based window (lighter than
Electron, since it reuses the OS's built-in WebView instead of bundling
Chromium, but still a browser engine under the hood).

### 43.2 Real-world example

An internal engineering tool that just needs a form, a table, and a
button to trigger a Go backend job is almost always fastest and easiest
to maintain as "serve HTML on localhost, open it in the browser" — no
GUI toolkit dependency at all. Reach for Fyne/Wails specifically when you
need a real desktop app: a system tray icon, native file dialogs, an
app users double-click from their Applications folder, or offline use
with no browser required.

---

## 44. Building a Fyne desktop app: a live system-resource monitor

**Goal:** a small, real Fyne app — a window showing live CPU/goroutine
counts, updating on a ticker, exactly the kind of "wrap a Go backend in a
native window" task Fyne is best at.

```bash
go get fyne.io/fyne/v2@latest
```

For a long-lived project, replace `@latest` with the release you chose
after testing and commit the resulting `go.mod`/`go.sum`.

```go
// main.go
package main

import (
	"fmt"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("System Monitor")
	w.Resize(fyne.NewSize(360, 200))

	goroutines := widget.NewLabel("")
	heapAlloc := widget.NewLabel("")
	cpus := widget.NewLabel(fmt.Sprintf("CPUs: %d", runtime.NumCPU()))

	content := container.NewVBox(
		widget.NewLabelWithStyle("Live System Stats", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		cpus,
		goroutines,
		heapAlloc,
	)
	w.SetContent(content)

	// Fyne widgets must only be updated from the main/UI goroutine's
	// event loop — this ticker goroutine only computes values; it hands
	// the actual widget update back via fyne.Do (Fyne 2.5+) / a channel.
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			g := runtime.NumGoroutine()
			heap := mem.HeapAlloc

			fyne.Do(func() {
				goroutines.SetText(fmt.Sprintf("Goroutines: %d", g))
				heapAlloc.SetText(fmt.Sprintf("Heap: %.2f MB", float64(heap)/1024/1024))
			})
		}
	}()

	w.ShowAndRun()
}
```

```bash
go run .
```

The one rule that matters for *any* GUI toolkit, not just Fyne: **never
mutate UI widgets from a background goroutine directly.** GUI toolkits
assume single-threaded access to their widget tree; Fyne's `fyne.Do`
(and, in older Fyne, a channel back to the main loop) is exactly the same
"hand data across a boundary, let the owning thread do the mutation"
discipline you already know from channels in Chapter 16 — the UI thread
*is* just another goroutine with an exclusivity rule.

---

## 45. Building a Wails app: a Go backend behind a web-based to-do board

**Goal:** show the "Go backend + web frontend in a native window" model,
where your Go code is called directly from JavaScript with no HTTP layer
in between — Wails auto-generates the JS bindings for exported Go
methods.

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails init -n todoboard -t vanilla
cd todoboard
```

For repeatable CI builds, pin the CLI to a specific release instead of
installing `@latest`.

```go
// app.go
package main

import (
	"context"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type App struct {
	ctx   context.Context
	mu    sync.Mutex
	todos []Todo
	nextID int
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// Exported methods become callable directly from JS as
// window.go.main.App.AddTodo(...) etc, via Wails' generated bindings.

func (a *App) AddTodo(text string) []Todo {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.nextID++
	a.todos = append(a.todos, Todo{ID: a.nextID, Text: text})
	return a.todos
}

func (a *App) ToggleTodo(id int) []Todo {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := range a.todos {
		if a.todos[i].ID == id {
			a.todos[i].Done = !a.todos[i].Done
		}
	}
	return a.todos
}

func (a *App) ListTodos() []Todo {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.todos
}
```

```javascript
// frontend/src/main.js
import { AddTodo, ToggleTodo, ListTodos } from '../wailsjs/go/main/App';

async function render() {
  const todos = await ListTodos();
  document.getElementById('list').innerHTML = todos
    .map(t => `<li class="${t.done ? 'done' : ''}" data-id="${t.id}">${t.text}</li>`)
    .join('');
}

document.getElementById('add').addEventListener('click', async () => {
  const input = document.getElementById('text');
  await AddTodo(input.value);
  input.value = '';
  render();
});

document.getElementById('list').addEventListener('click', async (e) => {
  const id = e.target.dataset.id;
  if (id) { await ToggleTodo(Number(id)); render(); }
});

render();
```

```bash
wails dev    # hot-reloading dev window
wails build  # produces a single native binary per OS with the frontend embedded
```

The mental model: your Go struct's exported methods *are* the API — no
REST layer, no JSON-over-HTTP boilerplate, no CORS. Wails serializes
arguments and return values across the JS↔Go boundary automatically. This
is the fastest path from "I know Go and basic HTML/JS" to a shippable
native desktop app.

---

## 46. Packaging and distributing GUI apps

### 46.1 Fyne: `fyne package`

```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin -icon icon.png   # produces a .app bundle
fyne package -os windows -icon icon.png  # produces a .exe with an embedded icon
fyne package -os linux -icon icon.png    # produces a portable Linux binary
```

As with Wails, pin the packaging tool version in CI once you've validated
the release.

### 46.2 Wails: `wails build`

```bash
wails build -platform darwin/universal,windows/amd64,linux/amd64
```

Wails cross-compiling to Windows/Linux from macOS (and vice versa)
requires the target platform's WebView runtime considerations — check
Wails' current cross-build docs before assuming a single command covers
every OS combination; native builds per-OS in CI remain the most reliable
path for a real release pipeline.

### 46.3 Code signing and notarization (macOS), and why it matters

An unsigned macOS `.app` triggers Gatekeeper's "unidentified developer"
warning; unsigned Windows executables trigger SmartScreen warnings. For
anything beyond a personal tool, budget time for an Apple Developer ID
(`codesign` + `notarytool`) and a Windows code-signing certificate — this
is packaging friction that has nothing to do with Go itself but blocks
real-world distribution just as much as any code bug would.

### 46.4 Real-world example

The reason Wails apps and Electron apps both bundle noticeably larger
binaries than a pure-Fyne app is the embedded web engine: Wails reuses
the OS's system WebView (WebView2 on Windows, WKWebView on macOS,
WebKitGTK on Linux) rather than bundling its own Chromium like Electron
— meaningfully smaller, but still a browser engine, versus Fyne/Gio's
"just draw the pixels ourselves in Go" approach with no embedded browser
at all.

---

# Part VI — Security-focused Go

A different angle on everything so far: the same language, aimed
specifically at not shipping vulnerabilities. Secure coding practices first,
then three services that put them into practice — password hashing and JWT
issuance, a WAF-style rate-limiting middleware, and a full mutual-TLS setup —
the security guide in this Wiki's companion volume, applied to Go's own
standard library (`crypto/tls`, `net/http` middleware) rather than covered
in the abstract.

## 47. Secure coding practices in Go

### 47.1 Input validation at the boundary

Validate and normalize every external input (HTTP body, query params,
CLI args, file uploads) at the point it enters your system — not deep
inside business logic where it's easy to forget a path. Go's explicit
error returns make "validate, then return early" natural:

```go
func validateEmail(s string) error {
	if len(s) == 0 || len(s) > 254 {
		return fmt.Errorf("email length invalid")
	}
	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}
	return nil
}
```

### 47.2 SQL injection: always use placeholders

```go
// NEVER — string-formats untrusted input directly into SQL
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", userInput)

// ALWAYS — the driver escapes/parameterizes the value safely
row := db.QueryRow("SELECT * FROM users WHERE email = ?", userInput)
```

`database/sql`'s placeholder syntax (`?` for MySQL/SQLite, `$1` for
Postgres) sends the value separately from the query text — the database
driver never treats it as SQL syntax, closing the injection vector
entirely. There is no legitimate reason to string-format untrusted input
into a query in Go.

### 47.3 Command injection: avoid the shell entirely

```go
// DANGEROUS if userInput is attacker-controlled — invokes a shell that
// interprets ;, |, `` etc.
exec.Command("sh", "-c", "ping "+userInput)

// SAFE — arguments are passed directly to execve, never interpreted by a shell
exec.Command("ping", "-c", "4", userInput)
```

`exec.Command`'s variadic arguments go straight to the OS's process-exec
call, bypassing shell parsing entirely, as long as you don't explicitly
route through `sh -c`.

### 47.4 Secrets: never in source, never logged

```go
// Read from environment or a secrets manager, never hardcode
apiKey := os.Getenv("API_KEY")

// Redact before logging anything that might contain a secret
logger.Info("api call", "key_prefix", apiKey[:4]+"...")
```

Go's `go vet` and third-party tools (`gitleaks`, `trufflehog`) catch
committed secrets; the cheaper fix is simply never letting them near
source control or structured logs in the first place.

### 47.5 Real-world example

Chapter 48–50 build three of the most common security-adjacent Go
services from scratch: an auth service issuing/verifying JWTs, a WAF-style
rate-limiting middleware, and a full mutual-TLS handshake — the same
territory covered conceptually in this Wiki's networking/security field
guides, here implemented directly in Go.

**Applied in the series:** the HTTPS guide's [SSRF guard](../v2-https/real-life-guide-v1.md#chapter-22-outbound-requests-ssrf-and-safe-http-clients) blocks internal addresses at connect time, and its [Slowloris lab](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) shows why server timeouts are a security control.


---

## 48. Auth service: password hashing and JWT issuing/verification

### 48.1 Never store plaintext or naively-hashed passwords

```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
```

`bcrypt` (or `argon2id` for stronger, memory-hard hashing) is
deliberately slow and salted per-call — a single unsalted SHA-256 hash of
a password is crackable via precomputed rainbow tables in seconds; bcrypt
at a reasonable cost factor takes meaningfully longer per guess, at scale.

### 48.2 Issuing a JWT

```go
import "github.com/golang-jwt/jwt/v5"

var signingKey = []byte(os.Getenv("JWT_SECRET")) // never hardcode this

func issueToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
```

### 48.3 Verifying a JWT — and the pitfalls that matter

```go
func verifyToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		// CRITICAL: verify the signing method yourself. Never trust the
		// token's own "alg" header blindly — that's the classic
		// "alg: none" / algorithm-confusion JWT attack.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	sub, _ := claims["sub"].(string)
	return sub, nil
}
```

### 48.4 Auth middleware

```go
func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := verifyToken(tokenStr)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

### 48.5 Real-world example

Short-lived access tokens (minutes, as above) paired with a separate,
longer-lived refresh token stored server-side (so it can be revoked) is
the standard production pattern — a stolen access token expires quickly
on its own, and revocation only has to be checked at the much
less-frequent refresh step.

---

## 49. A WAF-style middleware: rate limiting and IP blocklists

**Goal:** extend Chapter 39's gateway into something closer to a real
Web Application Firewall layer — inspecting and gating requests before
they reach the application, the same territory as this Wiki's WAF/DDoS
security material, implemented as ordinary Go middleware.

```go
package waf

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

type Firewall struct {
	mu        sync.RWMutex
	blocklist map[string]bool
	limiters  map[string]*rate.Limiter
	rps       rate.Limit
	burst     int
}

func New(rps float64, burst int) *Firewall {
	return &Firewall{
		blocklist: map[string]bool{},
		limiters:  map[string]*rate.Limiter{},
		rps:       rate.Limit(rps),
		burst:     burst,
	}
}

func (f *Firewall) Block(ip string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blocklist[ip] = true
}

func (f *Firewall) isBlocked(ip string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.blocklist[ip]
}

func (f *Firewall) limiterFor(ip string) *rate.Limiter {
	f.mu.Lock()
	defer f.mu.Unlock()
	lim, ok := f.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(f.rps, f.burst)
		f.limiters[ip] = lim
	}
	return lim
}

// suspiciousPatterns is a minimal illustration — a real WAF uses a
// maintained ruleset (e.g. OWASP CRS), not a hardcoded list.
var suspiciousPatterns = []string{"../", "<script", "UNION SELECT", "; DROP TABLE"}

func looksLikeAttack(r *http.Request) bool {
	target := r.URL.RawQuery + " " + r.URL.Path
	for _, p := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(target), strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func (f *Firewall) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if f.isBlocked(ip) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if looksLikeAttack(r) {
			f.Block(ip) // escalate: one bad request blocks this IP going forward
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if !f.limiterFor(ip).Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

```go
func main() {
	fw := waf.New(10, 20) // 10 req/s, burst 20, per source IP
	mux := http.NewServeMux()
	mux.HandleFunc("/", appHandler)

	log.Fatal(http.ListenAndServe(":8080", fw.Middleware(mux)))
}
```

This is illustrative, not production-grade — a real deployment puts this
kind of logic behind (or alongside) a dedicated reverse proxy/CDN layer,
uses a maintained signature set instead of a hardcoded pattern list, and
persists the blocklist outside process memory so a restart doesn't wipe
it. The point worth internalizing is architectural: **the firewall is
just another `http.Handler`-wrapping-`http.Handler` middleware**, composed
in front of your application exactly like logging or auth.

---

## 50. mTLS: client-certificate authentication end to end

**Goal:** go one level past Chapter 35's server-only TLS — require the
*client* to present a trusted certificate too, the mechanism behind
service-to-service auth in zero-trust architectures.

```bash
# 1. Create a private CA
openssl req -x509 -newkey rsa:4096 -days 3650 -nodes \
  -keyout ca-key.pem -out ca-cert.pem -subj "/CN=My Test CA"

# 2. Server certificate, signed by the CA
openssl req -newkey rsa:2048 -nodes -keyout server-key.pem -out server-req.pem \
  -subj "/CN=localhost"
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem \
  -CAcreateserial -out server-cert.pem -days 365

# 3. Client certificate, also signed by the same CA
openssl req -newkey rsa:2048 -nodes -keyout client-key.pem -out client-req.pem \
  -subj "/CN=trusted-client"
openssl x509 -req -in client-req.pem -CA ca-cert.pem -CAkey ca-key.pem \
  -CAcreateserial -out client-cert.pem -days 365
```

```go
// server.go
func main() {
	caCert, err := os.ReadFile("ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert, // reject any client with no valid cert
		ClientCAs:  caPool,
	}

	srv := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cn := r.TLS.PeerCertificates[0].Subject.CommonName
			fmt.Fprintf(w, "hello, %s — you presented a valid client cert\n", cn)
		}),
	}
	log.Fatal(srv.ListenAndServeTLS("server-cert.pem", "server-key.pem"))
}
```

```go
// client.go
func main() {
	cert, err := tls.LoadX509KeyPair("client-cert.pem", "client-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := os.ReadFile("ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		log.Fatal("failed to parse CA certificate")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caPool,
			},
		},
	}
	resp, err := client.Get("https://localhost:8443/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
```

`tls.RequireAndVerifyClientCert` is the setting that turns this from
"encrypted" into "mutually authenticated" — without it, TLS still
encrypts the connection, but the server accepts *any* client, which is
plain TLS, not mTLS. `r.TLS.PeerCertificates[0].Subject.CommonName` is
how the server learns *who* connected — the certificate itself is the
credential, no password or token exchanged at all.

---

# Part VII — Expert: internals and design

The last part goes underneath abstractions you've been using without
question since Part I: what an interface, a slice, a map, and a goroutine
stack actually are in memory. Then it turns from internals to judgment — API
design idioms and dependency injection without a framework — before closing
with a capstone embedded database engine that draws on nearly everything
earlier in the guide, and a final chapter on actually shipping a Go binary.

## 51. What's really behind an interface value

### 51.1 Two words, not one

An interface value is a two-word pair: a pointer to type information (the
concrete type's method table, called an **itab** when the interface has
methods, or nothing extra for the empty interface's **eface**), and a
pointer/value holding the underlying data.

```go
var w io.Writer = os.Stdout
// conceptually: { itab: *itab-for-(*os.File,io.Writer), data: pointer-to-os.Stdout }
```

This is exactly why the nil-interface trap from Chapter 9.4 exists: a
`nil` `*MyError` stored in an `error` interface produces a two-word value
`{itab: *itab-for-(*MyError,error), data: nil}` — the *data* pointer is
nil, but the interface value as a whole is not, because its type pointer
is set.

### 51.2 Dynamic dispatch, roughly

Calling a method on an interface value looks up the method in the itab
(a small, per-concrete-type-and-interface-pair table built once, cached),
then calls through that function pointer with the data word as the
receiver — conceptually similar to a C++ vtable call, but resolved
per-(type, interface)-pair rather than per-type alone.

### 51.3 Real-world example

This is *why* interface method calls have a small but real cost versus a
direct concrete-type method call (an extra pointer indirection through
the itab) — usually irrelevant, but the reason performance-critical inner
loops sometimes avoid an interface-typed hot path in favor of a concrete
generic type parameter (Chapter 21), which the compiler can inline
directly with no itab lookup at all.

---

## 52. What's really behind a slice, a map, and a string

### 52.1 Slice — covered fully in Chapter 6

A three-word header: pointer, length, capacity. Worth repeating here
because it's the single most consequential piece of "what's really
there" knowledge in the whole language.

### 52.2 Map — a hash table of buckets, with a twist

Go's `map` is implemented as an array of buckets, each holding a small
fixed number of key/value pairs plus overflow bucket pointers for
collisions, with the array growing (and old buckets migrating
incrementally, not all at once) as the map fills. The "twist": Go
deliberately randomizes the starting bucket for iteration precisely to
prevent programs from ever depending on a stable order (Chapter 7.3) —
this is enforced by the implementation, not just documented as
unspecified.

### 52.3 String — a two-word, read-only view

```go
type stringHeader struct {
	ptr *byte // pointer to an immutable byte array
	len int
}
```

No capacity field, because strings never grow in place — every operation
that "modifies" a string (concatenation, slicing rules aside) allocates a
new one. Slicing a string (`s[2:5]`) *does* share the underlying byte
array without copying, exactly like slicing a slice (Chapter 6.2) — one
more reason large substrings of a huge string you no longer need the rest
of can pin that whole backing array in memory unless you explicitly copy
the substring out.

### 52.4 Real-world example

A service parsing huge log lines and storing only a short substring
(say, a request ID) per entry can accidentally retain the *entire original
line's byte array* in memory for every stored substring, because slicing
a string doesn't copy. The fix mirrors Chapter 6's slice-aliasing fix:
`strings.Clone(sub)` (Go 1.18+) forces an actual copy, letting the
original large string be garbage collected.

---

## 53. Goroutine stacks and the scheduler, revisited

### 53.1 Growable stacks

A goroutine starts with a small stack (historically 2KB) and grows by
copying to a larger allocation when it would overflow — checked cheaply
at function-call prologues. This is why launching a goroutine per
connection (Chapters 26, 38) scales to tens of thousands of concurrent
connections without the multi-megabyte-per-thread cost a traditional
one-thread-per-connection server pays.

### 53.2 What actually blocks a goroutine (and what doesn't)

Channel operations, mutex waits, and blocking syscalls park a goroutine
without consuming an OS thread while parked (Chapter 28.2's syscall
detachment). CPU-bound work with no goroutine-yielding operation used to
be a scheduling risk before cooperative preemption existed everywhere
(Chapter 28.4); today it mostly is not, though a truly pathological
tight loop with function inlining that removes preemption checkpoints
remains a rare edge case worth knowing about, not a daily concern.

### 53.3 Real-world example

The reason "spawn 50,000 goroutines to handle 50,000 concurrent
WebSocket connections" is a completely normal, load-tested Go
architecture — and the equivalent in a thread-per-connection language
would need a very large machine or an event-loop rewrite — comes directly
from this stack-growth-plus-scheduler design, not from any particular
library.

---

## 54. Designing APIs: the Go proverbs, applied

### 54.1 "Accept interfaces, return structs"

A function that *accepts* a narrow interface (`io.Reader`, not
`*os.File`) works with anything satisfying that shape; a function that
*returns* a concrete struct gives callers full access to every field and
method, instead of forcing them through whatever subset an interface
happened to expose. Returning an interface is occasionally right (when
you deliberately want to hide the concrete type, e.g. a constructor for
an internal implementation detail) but should be the exception, not the
default.

### 54.2 "The bigger the interface, the weaker the abstraction"

A one-or-two-method interface (`io.Reader`, `Storage.Get`/`Set`) can be
satisfied by almost anything, which is exactly what makes it reusable. A
ten-method interface can usually only be satisfied by one real
implementation and a hand-written mock — at which point it isn't really
abstracting anything, just adding indirection.

### 54.3 "Don't just check errors, handle them gracefully"

```go
// Weak: propagates without adding anything a caller couldn't already see
if err != nil {
	return err
}

// Better: adds context that helps whoever reads this error three layers up
if err != nil {
	return fmt.Errorf("loading config from %s: %w", path, err)
}
```

### 54.4 Avoid premature abstraction

A single interface with one implementation and no second implementation
in sight is usually just ceremony — Go's implicit interface satisfaction
(Chapter 9.1) means you can introduce the interface *later*, exactly when
a second implementation (or a test double) actually shows up, with zero
cost paid up front. This mirrors the general engineering principle: three
similar concrete structs are better than one premature interface built to
anticipate a variation that may never arrive.

### 54.5 Real-world example

The `Storage` interface from Chapter 9 earns its keep specifically
because it *does* get a second implementation (`CachedStorage` wrapping
it) — that's the signal an interface was the right call, retroactively.
If it had stayed at one implementation forever, a plain
`*MemStorage` struct with no interface would have been equally correct
and slightly simpler.

---

## 55. Dependency injection without a framework: functional options

### 55.1 The problem: configurable constructors without exploding parameter lists

```go
// Unworkable past 3-4 options: which bool is which? what's the order?
func NewServer(host string, port int, timeout time.Duration, tls bool, maxConns int) *Server
```

### 55.2 The functional options pattern

```go
type Server struct {
	host     string
	port     int
	timeout  time.Duration
	tls      bool
	maxConns int
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}
func WithTLS() Option {
	return func(s *Server) { s.tls = true }
}
func WithMaxConns(n int) Option {
	return func(s *Server) { s.maxConns = n }
}

func NewServer(host string, port int, opts ...Option) *Server {
	s := &Server{host: host, port: port, timeout: 30 * time.Second, maxConns: 100}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

srv := NewServer("localhost", 8080,
	WithTLS(),
	WithTimeout(5*time.Second),
)
```

Every option is self-documenting at the call site, new options are
purely additive (existing call sites never need to change), and sensible
defaults live in one place. This is how most well-designed Go SDKs
(cloud provider clients, gRPC's `DialOption`, `http.Client`-adjacent
libraries) expose configuration.

### 55.3 Dependency injection as plain constructor parameters

Go doesn't need (and mostly avoids) DI *frameworks* — passing
dependencies explicitly as constructor parameters, typed as interfaces
where you want substitutability (Chapter 54.1), *is* dependency
injection:

```go
func NewOrderService(db Database, cache Cache, logger *slog.Logger) *OrderService {
	return &OrderService{db: db, cache: cache, logger: logger}
}
```

Tests substitute a fake `Database`/`Cache` implementation; production
wiring happens once, explicitly, in `main()`. No reflection-based
container, no annotations, no magic — the entire dependency graph is
readable by following function calls.

---

## 56. Capstone: a tiny embedded database engine (Bitcask-style)

**Goal:** the expert-tier capstone — an append-only log-structured
key-value store with an in-memory index, the same core design as
Bitcask/early LevelDB, built with nothing but the standard library. This
pulls together file I/O (Ch. 22), binary encoding (Ch. 34.4), generics
(Ch. 21/26), and `sync` (Ch. 17).

### 56.1 The design

Writes always **append** to a single log file (`data.log`) — never
seek-and-overwrite — which makes writes fast (sequential disk I/O) and
crash-safe (a half-written record at the tail is simply the last thing to
recover, never corrupting earlier data). An in-memory map indexes each
key to the **byte offset** of its most recent record in the log, so reads
are one seek plus one read, not a scan.

```go
package bitcask

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"
)

var ErrNotFound = errors.New("key not found")

type entryHeader struct {
	KeyLen   uint32
	ValueLen uint32
	Tombstone uint8 // 1 = this record represents a delete
}

const headerSize = 4 + 4 + 1

type indexEntry struct {
	offset int64
	size   int64
}

type DB struct {
	mu    sync.RWMutex
	file  *os.File
	index map[string]indexEntry
}

func Open(path string) (*DB, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	db := &DB{file: f, index: map[string]indexEntry{}}
	if err := db.rebuildIndex(); err != nil {
		_ = f.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) rebuildIndex() error {
	if _, err := db.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	var offset int64
	for {
		header := make([]byte, headerSize)
		if _, err := io.ReadFull(db.file, header); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return err
		}
		keyLen := binary.BigEndian.Uint32(header[0:4])
		valLen := binary.BigEndian.Uint32(header[4:8])
		tombstone := header[8]

		key := make([]byte, keyLen)
		if _, err := io.ReadFull(db.file, key); err != nil {
			return err
		}
		recordSize := int64(headerSize) + int64(keyLen) + int64(valLen)

		if tombstone == 1 {
			delete(db.index, string(key))
			if _, err := db.file.Seek(int64(valLen), io.SeekCurrent); err != nil {
				return err
			}
		} else {
			db.index[string(key)] = indexEntry{
				offset: offset + int64(headerSize) + int64(keyLen),
				size:   int64(valLen),
			}
			if _, err := db.file.Seek(int64(valLen), io.SeekCurrent); err != nil {
				return err
			}
		}
		offset += recordSize
	}
	return nil
}

func (db *DB) Put(key, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	offset, err := db.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	if err := writeRecord(db.file, key, value, 0); err != nil {
		return err
	}
	db.index[string(key)] = indexEntry{
		offset: offset + int64(headerSize) + int64(len(key)),
		size:   int64(len(value)),
	}
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.RLock()
	entry, ok := db.index[string(key)]
	db.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}

	db.mu.RLock()
	defer db.mu.RUnlock()
	value := make([]byte, entry.size)
	if _, err := db.file.ReadAt(value, entry.offset); err != nil {
		return nil, err
	}
	return value, nil
}

func (db *DB) Delete(key []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, err := db.file.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if err := writeRecord(db.file, key, nil, 1); err != nil {
		return err
	}
	delete(db.index, string(key))
	return nil
}

func writeRecord(w io.Writer, key, value []byte, tombstone uint8) error {
	header := make([]byte, headerSize)
	binary.BigEndian.PutUint32(header[0:4], uint32(len(key)))
	binary.BigEndian.PutUint32(header[4:8], uint32(len(value)))
	header[8] = tombstone
	if _, err := w.Write(header); err != nil {
		return err
	}
	if _, err := w.Write(key); err != nil {
		return err
	}
	_, err := w.Write(value)
	return err
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.file.Close()
}
```

### 56.2 What this teaches, and where it stops being a toy

This gets you: `O(1)` writes (always append), `O(1)` reads (one map
lookup + one `ReadAt`, no seek-and-scan), and crash recovery by replaying
the log on startup. What it deliberately leaves out, and what real
engines (LevelDB, RocksDB, Bitcask itself) add on top: **compaction** (the
log grows forever otherwise — old, overwritten/deleted records are never
reclaimed here), **concurrent readers during compaction**, **checksums per
record** (to detect a torn write after an unclean shutdown, not just
truncate at EOF), and a **write-ahead-log + separate index snapshot** to
avoid replaying the entire log on every restart. Each of those is a
natural next exercise once this base version passes its own tests.

```go
func TestPutGetDelete(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(dir + "/data.log")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := db.Put([]byte("k1"), []byte("v1")); err != nil {
		t.Fatal(err)
	}
	v, err := db.Get([]byte("k1"))
	if err != nil || string(v) != "v1" {
		t.Fatalf("got (%s, %v)", v, err)
	}

	if err := db.Delete([]byte("k1")); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get([]byte("k1")); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRecoversIndexOnReopen(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/data.log"
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Put([]byte("k1"), []byte("v1")); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := reopened.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	v, err := reopened.Get([]byte("k1"))
	if err != nil || string(v) != "v1" {
		t.Fatalf("index not rebuilt correctly: got (%s, %v)", v, err)
	}
}
```

---

## 57. Shipping it: cross-compilation, Docker, systemd, graceful shutdown

### 57.1 A minimal, secure Docker image

```dockerfile
# Dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/server

FROM scratch
COPY --from=build /app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/app"]
```

`CGO_ENABLED=0` makes pure-Go builds avoid libc/cgo dependencies, which is
what makes `FROM scratch` — a base image with *nothing* in it, not even a
shell — practical. The result is often small and has far less OS/userland
attack surface: no shell to pop, no package manager, and fewer unrelated
libraries with their own CVEs.

### 57.2 Graceful shutdown

```go
func main() {
	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("forced shutdown:", err)
	}
}
```

`srv.Shutdown` stops accepting new connections immediately, but lets
in-flight requests finish (up to the passed context's deadline) before
returning — the difference between a rolling deploy that drops active
requests and one that doesn't.

### 57.3 systemd unit for a bare-metal/VM deployment

```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=My Go Service
After=network.target

[Service]
ExecStart=/opt/myapp/myapp
Restart=on-failure
User=myapp
AmbientCapabilities=CAP_NET_BIND_SERVICE
Environment=APP_PORT=8080

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now myapp
journalctl -u myapp -f
```

### 57.4 Real-world example

The complete "boring but correct" production checklist for a typical Go
service built across this guide: small self-contained binary/container (57.1), explicit
`http.Server` timeouts (24.1), graceful shutdown on SIGTERM (57.2),
structured logs with a request ID (36.1), `/debug/pprof` mounted
internally (20.3), and `runtime.NumGoroutine()` exposed as a metric
(15.3). None of these are exotic — together, they're the difference
between a service that degrades observably and recovers cleanly, and one
that's a mystery every time it misbehaves.

**Applied in the series:** [Linux guide, Chapter 82](../os-linux/real-life-os-guide.md#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes): a tested 14.8 MB distroless image, a hardened systemd unit (exposure score 3.1 vs 9.6 unhardened), and a Kubernetes manifest, each setting traced to the chapter that explains it.


---

# Part VIII — Professional Go: debugging, refactoring, operations

> This Part is what turns "I know Go" into "I can own Go systems." The
> language is only half the job; the professional skill is diagnosing
> production behavior, changing code safely, proving behavior with tests,
> shipping repeatable releases, and operating services with enough telemetry
> that incidents become debuggable instead of mysterious.

## 58. Debugging Go in production: Delve, pprof, traces, stack dumps, GODEBUG

### 58.1 Read the panic stack trace first

A Go panic already tells you the failing goroutine, the call stack, and
usually the exact nil pointer / bounds / type assertion problem. Start at
the first frame in your code, not at `runtime/panic.go`.

```text
panic: runtime error: invalid memory address or nil pointer dereference

goroutine 42 [running]:
example.com/app/user.(*Service).Get(...)
	/app/user/service.go:87
example.com/app/httpapi.handleGetUser(...)
	/app/httpapi/users.go:33
```

**Debugging rule:** fix the earliest invalid assumption, not the last
function on the stack. If `handleGetUser` passed a nil service, the bug is
probably construction/wiring, not the line that finally dereferenced it.

### 58.2 Delve for local debugging

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/server -- --config ./dev.yaml
```

Useful commands inside Delve:

```text
break user/service.go:87
continue
goroutines
goroutine 42
stack
locals
print req.UserID
next
step
```

Use Delve when a failure is easy to reproduce locally and the question is
"which branch/value got me here?" Use logs, traces, and pprof when the
failure only appears under concurrency, load, or production traffic.

### 58.3 Goroutine dumps for hangs and leaks

When a Go service is "stuck," you want to know where goroutines are parked.
Expose `net/http/pprof` on an internal-only port:

```go
import _ "net/http/pprof"

go func() {
	log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
}()
```

Then capture:

```bash
curl -s http://127.0.0.1:6060/debug/pprof/goroutine?debug=2 > goroutines.txt
go tool pprof http://127.0.0.1:6060/debug/pprof/profile?seconds=30
go tool trace trace.out
```

Patterns to recognize:

| Symptom | Likely cause |
|---|---|
| Many goroutines blocked on `<-ch` | producer exited, channel never closed, or missing cancellation |
| Many goroutines in `net/http.(*persistConn).roundTrip` | outbound HTTP calls lack timeout or downstream is slow |
| Many goroutines waiting on `sync.Mutex.Lock` | lock contention or code doing slow work while holding a lock |
| Heap grows while goroutine count grows | goroutine leak retaining request data, buffers, or closures |

### 58.4 `GODEBUG`, `GOTRACEBACK`, and runtime knobs

```bash
GOTRACEBACK=all ./server         # include all goroutines in crash output
GODEBUG=gctrace=1 ./server       # print GC cycle summaries
GOMEMLIMIT=512MiB ./server       # soft memory limit for the Go runtime
GOGC=50 ./server                 # collect more aggressively than default
```

These are diagnostic and operational tools, not magic fixes. If `gctrace`
shows constant GC work, the long-term fix is usually to reduce allocations
or retained heap, then set `GOMEMLIMIT`/`GOGC` intentionally for the
deployment environment.

**War story.** A service "randomly froze" every few hours. CPU was low,
memory was fine, and logs stopped. A goroutine dump showed hundreds of
goroutines blocked on `jobs <- item`: the worker pool had returned after
one worker hit an error, but the producer kept sending forever. The fix was
not "increase channel buffer"; it was `errgroup.WithContext`, so failure
canceled the whole pipeline and all senders selected on `ctx.Done()`.

---

## 59. Refactoring Go codebases: package seams, APIs, and compatibility

### 59.1 Refactor around package boundaries, not files

In Go, the package is the real unit of design. Files are just organization.
Before splitting a package, ask:

- Does the new package have a small, coherent purpose?
- Can it avoid importing the package that imports it? If not, you are
  creating a cycle.
- Is the API smaller than the implementation it hides?
- Would a caller outside this directory understand the exported names?

The common refactor is extracting a lower-level package:

```text
cmd/server      -> wires dependencies, parses config, starts HTTP server
internal/httpapi -> handlers, middleware, request/response types
internal/users   -> user service and business rules
internal/store   -> database implementation
```

Avoid extracting `internal/common` or `internal/utils` early. Those names
usually mean the package has no clear owner and will accumulate unrelated
helpers.

### 59.2 Introduce interfaces at the consumer side

Prefer this:

```go
// internal/users/service.go
type Store interface {
	User(ctx context.Context, id string) (User, error)
	Save(ctx context.Context, u User) error
}

type Service struct {
	store Store
}
```

Not this:

```go
// internal/postgres/store.go
type StoreInterface interface {
	User(ctx context.Context, id string) (User, error)
	Save(ctx context.Context, u User) error
}
```

The consumer knows the narrow behavior it needs. The implementation should
not define a broad interface just because it exists.

### 59.3 Preserve API compatibility deliberately

For exported packages, changing any of these is a breaking change:

- removing an exported name;
- changing a function signature;
- changing interface methods;
- changing struct fields callers set directly;
- changing documented error behavior.

Safer migration pattern:

```go
// Old API kept temporarily.
func New(addr string) *Client {
	return NewWithOptions(Options{Addr: addr})
}

type Options struct {
	Addr    string
	Timeout time.Duration
	Logger  *slog.Logger
}

func NewWithOptions(opts Options) *Client {
	// new implementation
}
```

Deprecate, migrate call sites, then remove only at the next major version.
For public modules, that often means semantic import versioning (`/v2`).

### 59.4 Mechanical refactoring checklist

```bash
go test ./...
go test -race ./...
go vet ./...
staticcheck ./...
gofmt -w .
go mod tidy
```

Run these before and after large moves. If the change is truly mechanical,
commit it separately from behavior changes; future reviewers and debuggers
will thank you.

---

## 60. Testing beyond tables: fuzzing, httptest, golden files, integration tests

### 60.1 `httptest` for real handler tests

```go
func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(handleHealth)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("Content-Type = %q", ct)
	}
}
```

Use `httptest.NewServer` when you need a real URL and network client:

```go
srv := httptest.NewServer(http.HandlerFunc(handle))
defer srv.Close()

resp, err := http.Get(srv.URL)
if err != nil {
	t.Fatal(err)
}
defer resp.Body.Close()
```

### 60.2 Golden files for large outputs

Golden files are useful when the output is too large to read comfortably in
the test body: rendered HTML, generated config, CLI output, SQL, JSON.

```go
func TestRender(t *testing.T) {
	got := RenderReport(sampleReport())
	want, err := os.ReadFile("testdata/report.golden")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), got); diff != "" {
		t.Fatal(diff)
	}
}
```

Keep golden files under `testdata/`; the Go tool ignores that directory for
normal package builds.

### 60.3 Fuzz tests for parsers and boundary-heavy code

```go
func FuzzParseSize(f *testing.F) {
	f.Add("10KB")
	f.Add("")
	f.Add("999999999999999999999TB")

	f.Fuzz(func(t *testing.T, input string) {
		_, _ = ParseSize(input) // must not panic
	})
}
```

Run:

```bash
go test -fuzz=FuzzParseSize ./...
```

Fuzzing shines for decoders, parsers, validators, path normalization, and
anything that handles untrusted input.

### 60.4 Integration tests without slowing every unit test

Use build tags or environment variables for tests that need Docker,
Postgres, Redis, cloud credentials, or network access:

```go
//go:build integration

func TestPostgresStore(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}
	// start/use real database
}
```

```bash
go test ./...
go test -tags=integration ./...
```

The professional habit is a pyramid: many fast unit tests, fewer
integration tests, and a small number of end-to-end smoke tests that prove
the deployed binary actually works.

### 60.5 Testing time and concurrency

Avoid tests that depend on real sleeps:

```go
time.Sleep(500 * time.Millisecond) // flaky: machine speed and scheduler dependent
```

Prefer injected clocks/timers or synchronization:

```go
done := make(chan struct{})
go func() {
	defer close(done)
	doWork()
}()

select {
case <-done:
case <-time.After(2 * time.Second):
	t.Fatal("timed out")
}
```

Timeouts in tests are safety nets, not the synchronization mechanism.

---

## 61. Production database access with `database/sql` and `pgx`

### 61.1 `database/sql` is a pool, not a connection

```go
db, err := sql.Open("postgres", dsn)
if err != nil {
	return err
}
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(30 * time.Minute)

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
if err := db.PingContext(ctx); err != nil {
	return err
}
```

`sql.Open` validates arguments; it does not prove the database is reachable.
Use `PingContext` during startup.

### 61.2 Always pass context into database calls

```go
func (s *Store) User(ctx context.Context, id string) (User, error) {
	const q = `select id, email, created_at from users where id = $1`
	var u User
	if err := s.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("query user %s: %w", id, err)
	}
	return u, nil
}
```

Context ties DB work to request cancellation and deadlines. Without it, a
client can disconnect while the database query continues pointlessly.

### 61.3 Transactions: defer rollback, check commit

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	return err
}
defer tx.Rollback()

if _, err := tx.ExecContext(ctx, `update accounts set balance = balance - $1 where id = $2`, amount, from); err != nil {
	return err
}
if _, err := tx.ExecContext(ctx, `update accounts set balance = balance + $1 where id = $2`, amount, to); err != nil {
	return err
}
if err := tx.Commit(); err != nil {
	return err
}
```

The deferred rollback is cleanup for every early return. After `Commit`, it
does nothing useful and its error is normally ignored.

### 61.4 `pgx` for PostgreSQL-heavy services

`database/sql` is portable. `pgx` is often better when PostgreSQL is the
system of record and you want Postgres-specific features, better type
support, `COPY`, notifications, or direct pool control.

```go
pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
if err != nil {
	return err
}
defer pool.Close()

var email string
if err := pool.QueryRow(ctx, `select email from users where id=$1`, id).Scan(&email); err != nil {
	return err
}
```

### 61.5 Migrations are part of the app lifecycle

Use a migration tool (`golang-migrate`, `goose`, Atlas, Flyway, Liquibase)
and keep migrations in source control:

```text
migrations/
  0001_create_users.up.sql
  0001_create_users.down.sql
  0002_add_last_login.up.sql
  0002_add_last_login.down.sql
```

Never rely on "the app creates whatever tables it needs" for production
schema evolution. Schema changes need review, ordering, rollback planning,
and observability just like code changes.

---

## 62. Production HTTP APIs: validation, timeouts, middleware, shutdown

### 62.1 Limit and validate request bodies

```go
func decodeJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return false
	}
	return true
}
```

Production handlers should treat every boundary value as hostile: path
params, query params, headers, cookies, JSON, multipart files, and remote
addresses.

### 62.2 Map internal errors to stable HTTP responses

```go
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, ErrInvalid):
		http.Error(w, "bad request", http.StatusBadRequest)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
```

Do not leak database errors, stack traces, SQL strings, token parsing
details, or filesystem paths to clients. Log the internal detail; return a
stable external contract.

### 62.3 Client and transport timeouts

```go
transport := &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 20,
	IdleConnTimeout:     90 * time.Second,
	TLSHandshakeTimeout: 5 * time.Second,
}
client := &http.Client{
	Transport: transport,
	Timeout:   10 * time.Second,
}
```

Use one long-lived client per destination/configuration. Creating a new
client per request defeats connection pooling.

### 62.4 Graceful shutdown in real servers

```go
srv := &http.Server{
	Addr:         ":8080",
	Handler:      handler,
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout:  2 * time.Minute,
}

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

go func() {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}()

<-ctx.Done()
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(shutdownCtx); err != nil {
	log.Printf("graceful shutdown failed: %v", err)
}
```

Shutdown only works if handlers respect request contexts. If a handler
starts background work, it must either finish quickly, detach intentionally,
or be owned by a separate worker lifecycle.

**Applied in the series:** [HTTPS guide, Chapter 19](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) measures the four server timeouts against Slowloris: 200 held connections were dropped within 10 s.


---

## 63. Modules and releases: MVS, private modules, versioning, binaries

### 63.1 Minimal version selection, practically

Go modules use **minimal version selection**: builds use the minimum module
versions required by your dependency graph, not necessarily the newest
available versions.

```bash
go list -m all
go mod graph
go mod why -m github.com/some/dependency
go get github.com/acme/lib@v1.4.2
go mod tidy
```

This makes builds reproducible, but it also means dependency bumps should
be deliberate and reviewed like code.

### 63.2 Private modules

```bash
go env -w GOPRIVATE=github.com/mycompany/*
go env -w GONOSUMDB=github.com/mycompany/*
```

`GOPRIVATE` tells the Go command not to use the public checksum database or
module proxy for matching paths. Without it, private import paths may leak
to public infrastructure during module resolution.

### 63.3 Semantic import versioning

Public module v2+ paths include the major version:

```text
module github.com/acme/client/v2
```

Callers import:

```go
import "github.com/acme/client/v2"
```

If you do not want a `/v2`, avoid breaking exported APIs. Add new functions
and options instead of changing old signatures.

### 63.4 Build metadata and repeatable releases

```bash
go build -trimpath \
  -ldflags "-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" \
  -o dist/myapp ./cmd/myapp
```

```go
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)
```

Expose this in `myapp version` and logs. During an incident, "which commit
is running?" should be answerable in seconds.

### 63.5 Release checklist

- Tag source (`v1.2.3`).
- Build from a clean tree in CI.
- Run tests, race tests where feasible, linters, and `govulncheck`.
- Produce checksums for artifacts.
- Sign artifacts if your users install them directly.
- Publish release notes with breaking changes, migrations, and known issues.

Tools like GoReleaser are useful because they encode this checklist.

---

## 64. Professional workflow: linting, vuln checks, CI, and quality gates

### 64.1 Baseline commands

```bash
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
go test -cover ./...
govulncheck ./...
staticcheck ./...
```

`go vet` and `staticcheck` are not style police; they catch real bugs:
wrong `Printf` formats, impossible type assertions, copied locks,
unreachable code, ineffective assignments, lost cancel functions.

### 64.2 A minimal GitHub Actions workflow

```yaml
name: ci
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go mod download
      - run: gofmt -w . && git diff --exit-code
      - run: go vet ./...
      - run: go test ./...
      - run: go test -race ./...
      - run: govulncheck ./...
```

For large repos, run `-race` on the most important packages or nightly if
runtime is too high for every PR.

### 64.3 Keep generated and vendored code separate

```text
internal/api/openapi.gen.go  # generated; don't hand-edit
internal/api/handlers.go     # hand-written
```

Generated code should have a reproducible command:

```go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o openapi.gen.go openapi.yaml
```

If nobody can regenerate it, the generated file is technical debt.

---

## 65. Observability in Go services: logs, metrics, traces, health, profiling

### 65.1 Logs: structured and boring

```go
logger.InfoContext(ctx, "request completed",
	"method", r.Method,
	"path", r.URL.Path,
	"status", status,
	"duration_ms", time.Since(start).Milliseconds(),
	"request_id", requestIDFrom(ctx),
)
```

Good logs answer: what happened, to whom, where, when, how long, and with
which correlation ID. They should not leak secrets.

### 65.2 Metrics: RED and USE

For request/response services, start with RED:

- **Rate**: requests per second.
- **Errors**: error count/rate by route/status class.
- **Duration**: latency histogram, not just average.

For resources, use USE:

- **Utilization**: CPU, memory, disk, connection pools.
- **Saturation**: queue length, goroutines, blocked workers.
- **Errors**: failed DB calls, rejected jobs, dropped messages.

Avoid high-cardinality labels like user ID, raw URL, email, token, request
ID, or arbitrary error string.

### 65.3 Health, readiness, and liveness are different

```go
mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // process is alive
})

mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
})
```

Liveness answers "should the orchestrator restart me?" Readiness answers
"should I receive traffic?" Do not make liveness depend on every downstream
service, or a database outage can cause restart storms.

### 65.4 Tracing and pprof

Trace IDs connect logs across services; pprof explains where CPU, memory,
goroutines, blocking, and mutex contention are happening inside one
process. Keep pprof internal-only and access-controlled.

**Real-world example.** A p99 latency regression with normal CPU often
shows up in a block or mutex profile, not a CPU profile: goroutines are
waiting, not burning cycles. That distinction changes the fix from
"optimize code" to "reduce lock scope / increase pool / remove shared
state."

---

## 66. Architecture and package boundaries: services, workers, and seams

### 66.1 A boring service layout

```text
cmd/server/main.go
internal/config
internal/httpapi
internal/users
internal/postgres
internal/worker
internal/observability
migrations
```

`cmd/server` wires dependencies. `internal/httpapi` owns HTTP transport.
`internal/users` owns business rules. `internal/postgres` owns database
details. `internal/worker` owns background jobs. This is not a rule; it is
a starting shape that keeps transport, business logic, and persistence from
collapsing into one package.

### 66.2 Keep business logic away from HTTP

Prefer:

```go
func (s *Service) Register(ctx context.Context, email string) (User, error)
```

over:

```go
func (s *Service) Register(w http.ResponseWriter, r *http.Request)
```

The first can be called from HTTP, CLI, tests, workers, gRPC, or a
migration. The second hard-couples business logic to one transport.

### 66.3 Background workers are part of architecture

```go
type Worker struct {
	log   *slog.Logger
	store Store
	jobs  <-chan Job
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job, ok := <-w.jobs:
			if !ok {
				return nil
			}
			if err := w.handle(ctx, job); err != nil {
				w.log.ErrorContext(ctx, "job failed", "id", job.ID, "err", err)
			}
		}
	}
}
```

Workers need the same discipline as HTTP handlers: context, logging,
metrics, retries, idempotency, and graceful shutdown.

### 66.4 Idempotency and retries

If an operation may be retried, give it an idempotency key and store the
result:

```go
type PaymentRequest struct {
	ID             string // idempotency key from caller
	AccountID      string
	AmountCents    int64
}
```

Retries without idempotency are how "charge the card once" becomes "charge
the card three times."

---

## 67. Running Go in containers and Kubernetes

### 67.1 Container image basics

```dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /app ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
```

`scratch` is smallest; distroless is often more convenient because it
includes useful basics such as CA certificates and a non-root user variant.
For cgo-heavy apps, use a runtime image that contains the required shared
libraries.

### 67.2 Kubernetes probes and graceful shutdown

```yaml
readinessProbe:
  httpGet:
    path: /readyz
    port: 8080
  periodSeconds: 5
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  periodSeconds: 10
terminationGracePeriodSeconds: 30
```

Your Go process should handle SIGTERM, stop accepting new work, let
in-flight requests finish up to a deadline, and then exit. Kubernetes will
send SIGKILL after `terminationGracePeriodSeconds`.

### 67.3 Resource limits and Go runtime settings

```yaml
resources:
  requests:
    cpu: "250m"
    memory: "256Mi"
  limits:
    cpu: "1"
    memory: "512Mi"
env:
  - name: GOMEMLIMIT
    value: "450MiB"
```

Set `GOMEMLIMIT` below the container memory limit so the runtime starts
working harder before the kernel OOM-kills the process. Leave room for
non-Go memory: stacks, mmap, cgo, kernel buffers, and the process itself.

### 67.4 Configuration and secrets

Use environment variables, mounted files, or a secret manager. Avoid baking
environment-specific config into the binary or image.

```go
type Config struct {
	Addr        string
	DatabaseURL string
	LogLevel    string
}

func LoadConfig() Config {
	return Config{
		Addr:        getenv("ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    getenv("LOG_LEVEL", "info"),
	}
}
```

Validate config at startup. A service should fail fast if required config is
missing; partial startup with broken dependencies creates harder incidents.

### 67.5 Deployment checklist

- Runs as non-root.
- Handles SIGTERM gracefully.
- Has readiness and liveness endpoints.
- Uses explicit HTTP server timeouts.
- Sets memory limits and, where useful, `GOMEMLIMIT`.
- Exposes structured logs, metrics, and traces.
- Does not log secrets.
- Has pprof available only on an internal/admin path.
- Includes build version/commit in logs and `/version`.
- Has rollback instructions.

---

**That's the full arc** — from `go build`'s hello-world loop in Part I to a
capstone database engine, production-shippable service, and professional
debugging/refactoring/operations workflow in Part VIII. What follows is
reference material, not more reading: idioms, pitfalls, a glossary, and
cheat sheets to check back against while you're actually writing Go. For
structured day-by-day practice on top of this guide's theory, see the
companion [`golang-90-day-plan.md`](./golang-90-day-plan.md).

**Applied in the series:** [Linux guide, Chapters 74–76](../os-linux/real-life-os-guide.md#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1) (PID 1 and zombies, GOMAXPROCS, GOMEMLIMIT) and the [TCP/IP guide's Kubernetes networking chapter](../networking/tcp-ip/real-life-guide-v1.md#chapter-58-kubernetes-networking-services-ingress-networkpolicy-cni) (`ndots:5`, the conntrack DNS race).


---

## Continue the series

With the language in hand, follow the course from [step 1: Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md). Each guide's "Build it in Go" labs build on the chapters of this guide.

---

# Appendices

## Appendix A — Go proverbs and idioms cheat sheet

- **Don't communicate by sharing memory; share memory by communicating.**
  Prefer a channel handoff over a mutex-guarded shared variable when the
  design allows it — it makes ownership explicit instead of implicit.
- **The bigger the interface, the weaker the abstraction.** (Ch. 54.2)
- **Accept interfaces, return structs.** (Ch. 54.1)
- **A little copying is better than a little dependency.** Duplicating a
  small helper function across two packages beats importing a whole
  library for it.
- **Clear is better than clever.** Optimize for the next reader, not for
  showing off language features.
- **Errors are values.** Handle them like any other value — wrap, inspect,
  compare — not like a control-flow escape hatch.
- **Don't just check errors, handle them gracefully.** (Ch. 54.3)
- **gofmt's style is no one's favorite, yet gofmt is everyone's
  favorite.** Never bikeshed formatting; run `gofmt`/`goimports` and move
  on.

## Appendix B — Common pitfalls and how to catch them

| Pitfall | Symptom | Fix / detection |
|---|---|---|
| Loop variable capture (pre-1.22) | Goroutines all print the same value | Upgrade `go` directive, or copy the var explicitly each iteration |
| Slice aliasing via `append` | Silent data corruption under specific cap conditions | Full slice expressions (`s[:n:n]`) or explicit copy when handing off sub-slices |
| Nil interface holding a nil pointer | `err != nil` is true for a "nil" error | Return `nil` explicitly, never a typed nil pointer, from `error`-returning functions |
| Writing to a nil map | Panic: assignment to entry in nil map | Always `make()` a map before writing |
| Goroutine leak | Steadily climbing goroutine count | Every goroutine takes a `context.Context` or `done` channel; expose `NumGoroutine()` as a metric |
| Forgetting `defer resp.Body.Close()` | Leaked connections, eventual pool exhaustion | Always defer `Close()` immediately after checking the error |
| Data race on a shared map/counter | Rare panics or corrupted state under load | Run tests and staging traffic with `-race`; guard with `sync.RWMutex` or `sync/atomic` |
| No `http.Server` timeouts | Slow/malicious clients hold connections open indefinitely | Set `ReadTimeout`/`WriteTimeout`/`IdleTimeout` explicitly |
| `range`-copying large struct elements | Mutations inside a `range` loop silently don't stick | Index into the slice (`items[i].Field = ...`) instead of mutating the loop variable |

## Appendix C — Glossary

- **Escape analysis** — the compiler's decision about whether a value can
  stay on the stack or must move to the heap. (Ch. 27.1)
- **GMP** — Goroutine / Machine (OS thread) / Processor (scheduling
  context): the three abstractions behind Go's scheduler. (Ch. 28)
- **itab** — the interface table pairing a concrete type with an
  interface's method set, used for dynamic dispatch. (Ch. 51.1)
- **Zero value** — the well-defined default value every type gets when
  declared without explicit initialization. (Ch. 3.1)
- **Sentinel error** — a package-level `error` value (`io.EOF`,
  `sql.ErrNoRows`) compared with `errors.Is` to detect a specific known
  condition. (Ch. 10.3)
- **Functional options** — a constructor pattern using variadic
  `func(*T)` arguments for optional configuration. (Ch. 55.2)
- **Framing** — an application-layer scheme (delimiters or
  length-prefixing) for recovering message boundaries from a byte stream
  protocol like TCP. (Ch. 34.4)

## Appendix D — Further reading

- The Go specification — https://go.dev/ref/spec (short enough to read
  end to end once)
- *Effective Go* — https://go.dev/doc/effective_go
- *The Go Memory Model* — https://go.dev/ref/mem
- Standard library docs — https://pkg.go.dev/std (searchable, always the
  primary source over any blog post)
- `golang.org/x/time/rate`, `golang.org/x/sync/errgroup`,
  `golang.org/x/sys/unix` — the "extended standard library" used
  throughout Parts III–VI
- This Wiki: [`golang-90-day-plan.md`](./golang-90-day-plan.md) (day-by-day
  exercises), [`detailed-90-day-plan/`](./detailed-90-day-plan/) (expanded
  daily lessons), [`projects/lru_cache.md`](./projects/lru_cache.md) (the
  design spec Chapter 26 implements)

## Appendix E — Command cheat sheet

```bash
go run .                          # compile + run, no binary kept
go build -o bin/app .             # produce a binary
go test ./...                     # run all tests
go test -race ./...               # run all tests with the race detector
go test -bench=. -benchmem ./...  # run benchmarks with allocation stats
go test -cover ./...              # coverage percentage
go vet ./...                      # static correctness checks
gofmt -l .                        # list files that need formatting
go mod tidy                       # sync go.mod/go.sum with actual imports
GOOS=linux GOARCH=amd64 go build  # cross-compile
go tool pprof <profile>           # analyze a CPU/heap/goroutine profile
go tool trace <trace.out>         # analyze a scheduler trace
```
