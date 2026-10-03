# 120-Day Go Engineering Plan: Novice → Advanced
## Target Persona: Software Engineer
### Networking (L4–L7) + Security (WAF, DDoS, Bot, API Protection)
### ~1 hour/day | Standard library only | Weekly mini projects

---

## Who This Plan Is For

You are a **software engineer** who already ships code in at least one language. You've built services, worked with HTTP APIs, and been on call for a production incident or two. But you haven't written much production Go, and networking/security protocols live in the "I know what they do, not how they work" category.

This plan fixes that. After 120 days you will:
- Write idiomatic, production-quality Go from scratch
- Understand networking from L4 (TCP/UDP) to L7 (HTTP, WebSocket, TLS) — by building it
- Build WAF, DDoS mitigation, bot detection, and API protection systems
- Have the hands-on experience to build, review, debug, and design security infrastructure

**Ground rules:**
- 1 hour/day. No exceptions. Protect this time.
- Standard library only (`net`, `crypto/tls`, `net/http`, `sync`, `context`, etc.)
- Every week ends with a mini project. Build it. Don't skip it.
- Each day has a "**Why this matters**" — connect the code to real engineering decisions you'd make as an EM

---

## The series: OS → networking → security → HTTPS: where this plan fits

This plan builds the software; the guides explain the systems it runs on.
Every networking and security week (5–17) opens with a **"Background in this
wiki"** note pointing at the matching chapters.

The guides are, designed to be read as **one course** in this order:

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
| ∥ | [Go — The Complete Field Guide](real-life-golang-guide.md) | the language behind every lab, plus the 120-day plan's multi-week projects | [120-day plan](golang-90-day-plan.md) |

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

## 120-Day Overview

| Phase | Days | Theme | EM Context |
|-------|------|-------|------------|
| Phase 1 | 1–30 | Go Fundamentals | Become a confident Go reader and writer |
| Phase 2 | 31–60 | Networking L4 + Secure Transport | Understand what your infra team builds |
| Phase 3 | 61–90 | L7 + WAF + API Security | Know how to protect your APIs |
| Phase 4 | 91–120 | DDoS + Bot + Integration + Capstone | Own the full security stack |

---

## Phase 1: Go Fundamentals (Days 1–30)

> **Phase Goal:** Go from "I can read Go" to "I can write Go without looking things up." This is the foundation. Don't rush it.

### Week 1: Go Syntax, Types & Functions (Days 1–7)

**Why this week matters:** Every Go service you'll build — microservice, proxy, security middleware — starts here. Writing and reviewing Go confidently requires fluency in these basics.

---

**Day 1 — Go toolchain & your first program**
- **Why:** Understanding `go build`, `go run`, `go test`, `go mod` is what you'll reference when engineers say "it works locally"
- **Topic:** Install Go, workspace setup, `go mod init`, `package main`, `func main()`, `fmt.Println`
- **Exercise:** Write a program that prints a formatted system info report: current time, hostname (`os.Hostname`), number of CPUs (`runtime.NumCPU`), Go version (`runtime.Version`)
- **Key packages:** `fmt`, `os`, `runtime`, `time`

**Day 2 — Variables, types & zero values**
- **Why:** Go's type system prevents a whole class of bugs you'd otherwise debug at runtime in dynamic languages
- **Topic:** `var`, `:=`, basic types (`string`, `int`, `bool`, `float64`), zero values, type conversion, `const`, `iota`
- **Exercise:** Build a unit converter: convert between MB/GB/TB, ms/s/min. Use typed constants with `iota` for units. No hardcoded magic numbers.

**Day 3 — Functions & multiple return values**
- **Why:** Go's `(value, error)` return pattern is in every library call. Understanding it is essential for reading Go code
- **Topic:** Function syntax, multiple return values, named returns, variadic functions, first-class functions
- **Exercise:** Build a `parseConfig(input string) (Config, error)` function that parses `key=value` pairs. Return structured errors. Write 5 test cases with `if err != nil` handling.

**Day 4 — Slices & maps**
- **Why:** These are the workhorse data structures in 90% of Go code. Misuse causes production memory leaks
- **Topic:** Slice internals (len/cap), `append`, slice of slice, map operations, nil vs empty slice/map, iteration with `range`
- **Exercise:** Build an in-memory request log: store last 1000 requests (IP, path, status, latency) in a slice with a fixed-size ring buffer. Query by IP using a map index.

**Day 5 — Structs & methods**
- **Why:** Go uses structs + methods instead of classes. Every service, every config, every request object is a struct.
- **Topic:** Struct definition, field tags, methods with value vs pointer receivers, struct embedding, anonymous fields
- **Exercise:** Model a `Server` struct with config fields. Add methods: `Start()`, `Stop()`, `Status() string`. Use pointer receivers. Embed a `Logger` struct.

**Day 6 — Interfaces**
- **Why:** Go interfaces enable the pluggable, testable architectures your best engineers design. Understanding them helps you spot good vs bad design in PRs
- **Topic:** Interface definition, implicit satisfaction, `interface{}` (any), type assertions, type switches
- **Exercise:** Define a `Storage` interface with `Set(key, value string)`, `Get(key string) (string, bool)`, `Delete(key string)`. Implement it with an in-memory map. Write a `CachedStorage` that wraps any `Storage` and adds metrics.

**Day 7 — Error handling**
- **Why:** Go's explicit error handling is why Go services tend to be more reliable. Knowing the patterns helps you require them in code review
- **Topic:** `error` interface, `errors.New`, `fmt.Errorf`, `errors.Is`, `errors.As`, wrapping with `%w`, sentinel errors
- **Exercise:** Build a multi-layer error system for the storage from Day 6: `ErrNotFound`, `ErrStorageFull`, `ErrInvalidKey`. Each wraps context. Write a handler that distinguishes each error type.

**Week 1 Mini Project: Configuration Server**
- Build an HTTP server (`net/http`) that serves runtime config as JSON
- Endpoints: `GET /config` (return all config), `PUT /config/{key}` (update a value), `GET /health` (return status + uptime)
- Use a struct for config, interface for storage, proper error responses
- Key packages: `net/http`, `encoding/json`, `os`, `time`
- **Engineer Takeaway:** This is the skeleton of every internal CLI and admin tool. You can now choose deliberately between `var` and `:=`, and between pointer and value receivers.

---

### Week 2: Control Flow, Packages & Testing (Days 8–14)

**Why this week matters:** Go's testing story is one of its strongest features. By end of week you'll write tests the way your senior engineers do.

---

**Day 8 — Control flow**
- **Why:** Understanding `defer` explains a lot of "why does this cleanup always run?" patterns in Go servers
- **Topic:** `if/else`, `switch` (no fallthrough by default), `for` (only loop in Go), `break`/`continue`, `defer` (LIFO, arguments evaluated immediately), `goto` (avoid)
- **Exercise:** Write a request lifecycle simulator: `defer` cleanup functions in the right order (close body → release connection → log → decrement counter). Observe execution order.

**Day 9 — Packages & modules**
- **Why:** `go mod` and package structure drive how Go code is organized. Knowing this lets you lay out a repo that scales past one package
- **Topic:** Package naming, exported vs unexported, `go.mod`, `go.sum`, internal packages, `init()` function, circular imports (why Go forbids them)
- **Exercise:** Refactor the Week 1 mini project into proper packages: `config/`, `storage/`, `server/`, `health/`. Each package has one clear responsibility. No circular dependencies.

**Day 10 — Testing in Go**
- **Why:** Go has testing built in — no framework needed. Understanding `go test` helps you set coverage requirements
- **Topic:** `testing.T`, table-driven tests, `t.Run` subtests, `t.Helper`, test file naming, `-run` flag, `-v` flag
- **Exercise:** Write table-driven tests for the storage interface from Day 6. Cover: empty key, value too long, concurrent reads, delete non-existent key. Run with `-v` and `-count=3`.

**Day 11 — Benchmarks & examples**
- **Why:** Benchmark-driven performance discussions beat intuition. Knowing how to read benchmarks changes how you evaluate "it's slow" reports
- **Topic:** `testing.B`, `b.N`, `b.ResetTimer`, `b.ReportAllocs`, `go test -bench`, `go test -benchmem`, example functions
- **Exercise:** Benchmark three implementations of the same lookup: linear scan slice, map lookup, sorted slice + binary search. Report ns/op and allocs/op. Which wins at 100 items? At 10,000?

**Day 12 — Goroutines (introduction)**
- **Why:** Go's concurrency model is why Go handles high traffic efficiently. This is the single most important Go concept for a networking/security stack
- **Topic:** `go` keyword, goroutine lifecycle, the fact that goroutines are not OS threads, `runtime.NumGoroutine()`, goroutine leaks
- **Exercise:** Write a program that starts 1000 goroutines, each sleeping for a random duration. Count active goroutines with `runtime.NumGoroutine()` every 100ms. Observe the ramp-up and drain. Ensure all goroutines exit cleanly.

**Day 13 — Channels**
- **Why:** Channels are how Go goroutines communicate safely. This pattern appears in every concurrent server you'll write
- **Topic:** Make channel, send/receive, buffered vs unbuffered, `close`, `range` over channel, `select`, `default`
- **Exercise:** Build a job queue: producer goroutine sends jobs into a buffered channel; 5 worker goroutines consume and process jobs; a results channel collects outcomes. Print summary when all done.

**Day 14 — sync package**
- **Why:** `sync.Mutex` and `sync.WaitGroup` are in every Go service. Misuse causes deadlocks. Knowing them helps you spot issues in review
- **Topic:** `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, `sync.Once`, `sync.Map`
- **Exercise:** Implement a thread-safe request counter: tracks requests per endpoint, reads are frequent (use `RWMutex`), resets happen rarely. Write a concurrent test with 100 goroutines hammering it. Run with `-race`.

**Week 2 Mini Project: Concurrent URL Checker**
- Build a tool that takes a list of URLs (from a file) and checks each one concurrently (HTTP GET, 5s timeout)
- Worker pool: configurable concurrency (default: 10)
- Results: status code, latency, error (if any), written to output file
- Graceful: if any URL check takes > 30s total, cancel remaining
- Key packages: `net/http`, `sync`, `context`, `os`, `time`
- **Engineer Takeaway:** This pattern (worker pool + channel pipeline + context cancellation) is in every batch job, health checker, and data pipeline you'll write.

---

### Week 3: Context, I/O & JSON (Days 15–21)

**Why this week matters:** `context` and `encoding/json` are in literally every production Go service. These are not optional topics.

---

**Day 15 — context package**
- **Why:** Context is how timeouts and cancellations propagate through your entire stack — from load balancer to database query. It's why Go services can be reliably cancelled
- **Topic:** `context.Background()`, `context.TODO()`, `WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`, `ctx.Done()`, `ctx.Err()`
- **Exercise:** Build a request simulator with 3 stages (auth → fetch data → format response). Each stage has its own timeout. Cancel everything if any stage fails or the parent context is cancelled. Use `select` on `ctx.Done()`.

**Day 16 — io.Reader & io.Writer**
- **Why:** Every file read, HTTP body, network socket, and log writer in Go is an `io.Reader` or `io.Writer`. Composing them is fundamental
- **Topic:** `io.Reader`, `io.Writer`, `io.Copy`, `io.TeeReader`, `io.LimitReader`, `io.MultiWriter`, `io.Pipe`, `bufio.Reader/Writer`
- **Exercise:** Build a request logger that: reads an HTTP request body, logs it (via `io.TeeReader`), limits to 10KB (via `io.LimitReader`), and passes it downstream. Don't read the body twice.

**Day 17 — encoding/json**
- **Why:** JSON is the wire format for APIs. Understanding Go's JSON encoder/decoder explains why engineers write struct tags the way they do
- **Topic:** `json.Marshal`, `json.Unmarshal`, struct tags (`json:"name,omitempty"`), `json.Decoder` (streaming), `json.Encoder`, custom `MarshalJSON`/`UnmarshalJSON`, `json.RawMessage`
- **Exercise:** Build a JSON API request/response pipeline: decode incoming JSON request (streaming, don't buffer full body), validate required fields, transform, encode response. Test with malformed JSON.

**Day 18 — File I/O & os package**
- **Why:** Config files, log files, PID files — services interact with the filesystem constantly. Knowing the patterns prevents bugs
- **Topic:** `os.Open`, `os.Create`, `os.ReadFile`, `os.WriteFile`, `os.Stat`, `filepath`, `os.MkdirAll`, `os.TempDir`, file permissions
- **Exercise:** Build a rotating log writer: write to `app-YYYY-MM-DD.log`, create new file daily, keep last 7 files (delete older ones), concurrent-safe with `sync.Mutex`.

**Day 19 — Strings & bytes**
- **Why:** Security-critical code manipulates strings constantly — URLs, headers, payloads. Knowing strings vs bytes prevents subtle bugs
- **Topic:** `strings` package (Builder, Split, TrimSpace, Contains, HasPrefix), `bytes` package, string immutability, `[]byte` conversions, `strings.Builder` for efficient concatenation
- **Exercise:** Build a URL normalizer: lowercase scheme and host, remove default ports (`:80`, `:443`), normalize path (remove double slashes, resolve `./` and `../`), sort query params. No regex yet.

**Day 20 — regexp**
- **Why:** WAF rules and log parsers rely on regex. Understanding Go's `regexp` package performance (compiled vs interpreted) matters for high-throughput security middleware
- **Topic:** `regexp.MustCompile`, `regexp.Compile`, `FindString`, `FindAllString`, `ReplaceAll`, `SubexpNames`, pre-compiling patterns, `regexp.MatchString` pitfall (re-compiles every call)
- **Exercise:** Build a log line parser using named capture groups: parse nginx/Apache access log format (`$ip - - [$time] "$method $path $proto" $status $bytes`). Compile all patterns at startup.

**Day 21 — time package**
- **Why:** Rate limiting, token expiry, metrics windows, log rotation — all time-dependent. Go's `time` package is rich and often misused
- **Topic:** `time.Now()`, `time.Duration`, `time.Since`, `time.Until`, `time.After`, `time.Tick` (leak risk), `time.NewTicker`, `time.NewTimer`, time zones, `time.Parse/Format`
- **Exercise:** Build a sliding-window rate limiter (purely time-based, no external state): each call checks if the count in the last N seconds exceeds a limit. Use `time.NewTicker` for cleanup. Thread-safe.

**Week 3 Mini Project: Structured Request Logger**
- Build an HTTP middleware that logs every request as structured JSON to a rotating file
- Log fields: timestamp, request ID (UUID-like, from `crypto/rand`), method, path, status, latency, client IP, user-agent, request size, response size
- Context: inject request ID into context at entry, retrieve at response write
- Rotation: new file daily, keep 14 days
- Key packages: `net/http`, `encoding/json`, `context`, `sync`, `time`, `os`, `crypto/rand`
- **Engineer Takeaway:** Audit logging is a compliance requirement. You now know what "structured logging" means in code, how to choose a log format, and why log rotation matters.

---

### Week 4: Advanced Go Patterns (Days 22–30)

**Why this week matters:** These are the patterns that separate junior Go from senior Go. Knowing them helps you review architecture proposals.

---

**Day 22 — Generics (Go 1.18+)**
- **Why:** Generics reduce code duplication in shared libraries. You'll write a generic cache or pool yourself, and know when generics are worth it
- **Topic:** Type parameters, constraints, `comparable`, `any`, generic functions, generic types
- **Exercise:** Build a generic `RingBuffer[T]` (fixed-size circular buffer). Use it with `string` (log entries) and `int` (metrics). Then build a generic `Result[T]` type wrapping `(T, error)`.

**Day 23 — sync.Pool & memory efficiency**
- **Why:** In high-traffic services (security proxies, WAF), allocations cause GC pressure. `sync.Pool` is how teams reduce that
- **Topic:** `sync.Pool`, GC interaction with Pool, when to use (reusable buffers), when not to use (long-lived objects), `bytes.Buffer` pooling pattern
- **Exercise:** Build a byte buffer pool. Write a benchmark comparing: (1) allocate new buffer each request, (2) use `sync.Pool`. Measure allocs/op and ns/op at 10k/req.

**Day 24 — Functional patterns**
- **Why:** Go middleware chains, option functions, and builder patterns are functional patterns. They appear everywhere in idiomatic Go
- **Topic:** Function types, middleware pattern (`func(Handler) Handler`), options pattern (`func(*Config)`), method chaining builders
- **Exercise:** Build a configurable HTTP server using the options pattern: `NewServer(WithPort(8080), WithTimeout(30*time.Second), WithTLS(certFile, keyFile))`. No config struct in the constructor signature.

**Day 25 — Embedding & composition**
- **Why:** Go uses composition over inheritance. Understanding embedding explains how interfaces get "extended" in Go
- **Topic:** Struct embedding, interface embedding, method promotion, shadowing, embedding vs field
- **Exercise:** Build a `MeteredConn` that embeds `net.Conn` and counts bytes read/written. Build a `TimedHandler` that embeds `http.Handler` and records latency. Both are transparent wrappers.

**Day 26 — expvar & runtime metrics**
- **Why:** `expvar` is built-in observability. Every production Go binary should expose `/debug/vars`. Expose it in every service you ship
- **Topic:** `expvar.Int`, `expvar.Float`, `expvar.Map`, `expvar.String`, `/debug/vars` HTTP endpoint, `net/http/pprof`
- **Exercise:** Add `expvar` metrics to the Week 3 request logger: requests total, errors total, p50/p99 latency (ring buffer of last 1000 latencies), bytes served total. Expose at `/debug/vars`.

**Day 27 — Graceful shutdown**
- **Why:** Abrupt server shutdowns drop in-flight requests. Every production HTTP server needs graceful shutdown. This is a common incident cause
- **Topic:** `os/signal`, `signal.NotifyContext`, `http.Server.Shutdown`, `sync.WaitGroup` for in-flight tracking, shutdown timeout budget
- **Exercise:** Build an HTTP server with 3-stage graceful shutdown: (1) stop accepting (10s), (2) wait for in-flight (30s), (3) force close. Test by sending requests mid-shutdown.

**Day 28 — Profiling**
- **Why:** When engineers say "the service is slow," you need to know what a flame graph means and how to generate one
- **Topic:** `go tool pprof`, CPU profile (`-cpuprofile`), heap profile (`-memprofile`), goroutine profile, `net/http/pprof` endpoints, reading flame graphs
- **Exercise:** Add pprof endpoints to your server. Write a load generator (100 goroutines, 10s). Collect a CPU profile. Open in pprof. Find the hottest function. Explain what the flame graph shows.

**Day 29 — Testing: httptest & net.Pipe**
- **Why:** Testing network code without real network calls is how you get fast, reliable CI
- **Topic:** `net/http/httptest.NewRecorder`, `httptest.NewServer`, `net.Pipe()`, table-driven HTTP handler tests
- **Exercise:** Write tests for the Week 3 middleware using `httptest`. Test: request logging, rotation trigger, concurrent requests, context cancellation mid-request. No real HTTP server.

**Day 30 — Review & code audit**
- Review all Phase 1 code with fresh eyes
- Apply: consistent error handling, no goroutine leaks (`go test -race`), `defer` cleanup, proper context propagation
- **Engineer Takeaway:** You can now read production Go code comfortably, write tests, benchmark performance, and profile CPU/memory. You're ready for networking.

**Week 4 Mini Project: In-Memory Metrics Server**
- Build a self-contained metrics collection + query server
- Ingest: `POST /metrics` accepts `{name, value, tags}` JSON
- Query: `GET /metrics/{name}?window=60s` returns min/p50/p95/p99/max over the time window
- Storage: ring buffer per metric name, concurrent-safe
- Observability: itself exposes `/debug/vars` and `/debug/pprof`
- Graceful shutdown with drain
- Key packages: `net/http`, `encoding/json`, `sync`, `time`, `expvar`, `net/http/pprof`
- **Engineer Takeaway:** This is a miniature Prometheus. You now know how metrics systems work inside, and why they use ring buffers, not unbounded slices.

---

## Phase 2: Networking Foundations + L4 (Days 31–60)

> **Phase Goal:** Understand how network communication actually works — from socket to TLS — by building it in Go. This is the knowledge gap between "I know TCP/IP" and "I understand what my network team deals with."

### Week 5: Networking Concepts + TCP Basics (Days 31–37)

**Why this week matters:** You've approved budget for CDNs, WAFs, and load balancers. This week explains what they actually do at the wire level.

---

**Day 31 — The OSI model in practice**
- **Why:** When an incident happens at L3 vs L7, the debugging approach is completely different. Knowing where your tools operate is essential for incident command
- **Topic:** OSI layers 1–7 (physical → application), where Go code lives (L4–L7), what happens between `net.Dial` and data delivery, packet encapsulation
- **Exercise:** Write a Go program that traces what happens when you connect to `example.com:80`: DNS resolution (`net.LookupHost`), TCP connection (`net.Dial`), raw HTTP/1.0 request (`conn.Write`), read response. Log every step with timing.

**Day 32 — TCP fundamentals**
- **Why:** TCP is the transport for HTTP, TLS, gRPC, and almost everything your services use. SYN floods, slowloris, and connection exhaustion all target TCP
- **Topic:** TCP 3-way handshake (SYN→SYN-ACK→ACK), connection states, sequence numbers, acknowledgements, flow control, congestion control (conceptual), `TIME_WAIT`
- **Exercise:** Write a TCP echo server (`net.Listen` + `net.Accept`). Handle 100 concurrent clients. After each echo, log: local addr, remote addr, bytes exchanged. Observe `netstat -an` while running.

**Day 33 — TCP server patterns**
- **Why:** Every microservice you run is a TCP server underneath. Understanding accept loops, connection limits, and timeouts helps you set sensible `http.Server` configs
- **Topic:** Accept loop, per-connection goroutine, connection lifecycle, `SetDeadline`, `SetReadDeadline`, `SetWriteDeadline`, half-close (`CloseWrite`)
- **Exercise:** Build a hardened TCP server: max 1000 concurrent connections (use `sync/atomic` counter), 30s idle timeout per connection, log connections by state (new/active/closing). Test with 2000 simultaneous clients.

**Day 34 — UDP**
- **Why:** DNS, QUIC, DTLS, and monitoring protocols (StatsD, syslog UDP) use UDP. Understanding it helps you evaluate QUIC-based CDN products
- **Topic:** UDP vs TCP (no handshake, no ordering, no guaranteed delivery), `net.PacketConn`, `net.UDPConn`, MTU, why UDP is used for DDoS amplification
- **Exercise:** Build a UDP-based statsd-like metrics receiver: listen on UDP, parse `metric.name:value|type` packets, aggregate (count/gauge/timer), print summary every 5s. Handle 10k packets/sec.

**Day 35 — I/O patterns for networking**
- **Why:** `io.Copy`, `io.TeeReader`, `io.LimitReader` — these are how Go efficiently moves bytes between connections. Knowing them explains how proxies work
- **Topic:** `io.Copy`, `io.TeeReader`, `io.LimitReader`, `io.MultiWriter`, `bufio` over connections, zero-copy concepts
- **Exercise:** Build a TCP traffic mirror: for each client connection, forward to upstream AND copy all traffic to a "tap" (write to a file). Client sees normal responses. Use `io.TeeReader`. No buffering of full request.

**Day 36 — Connection pools**
- **Why:** Connection pooling is why your services don't exhaust upstream connection limits. `http.Transport` does this automatically — but security proxies need custom pools
- **Topic:** Pool design (Get/Put/Close), idle timeout, max pool size, health checking, `sync.Pool` vs explicit pool, when connections go stale
- **Exercise:** Build a generic TCP connection pool: `Get(addr) (net.Conn, error)`, `Put(conn)`, `Close()`. Idle connections expire after 30s. Max 10 connections per address. Test with 50 concurrent goroutines.

**Day 37 — DNS in Go**
- **Why:** DNS is the first hop in every request. DNS failures cause incidents. DNS is also an amplification attack vector
- **Topic:** `net.LookupHost`, `net.LookupMX`, `net.LookupTXT`, `net.Resolver` (custom), DNS round-robin, TTL, negative caching
- **Exercise:** Build a DNS-aware connection dialer: resolve hostname to IPs (potentially multiple), try each in order, cache results with TTL, fall back to next IP on connection failure. Measure DNS resolution latency separately.

**Week 5 Mini Project: TCP Proxy with Traffic Inspector**
- Build a transparent TCP proxy: `client → proxy → upstream`
- Features:
  - Bidirectional forwarding (using goroutines + `io.Copy`)
  - Traffic tap: log all bytes in hex to a file (via `io.TeeReader`)
  - Connection tracking: active connections, total bytes in/out (via `expvar`)
  - Per-connection timeout (30s idle)
  - Connection pool to upstream (max 20 connections)
  - Graceful shutdown: drain existing connections
- Key packages: `net`, `io`, `sync`, `sync/atomic`, `context`, `expvar`, `time`
- **Engineer Takeaway:** A TCP proxy is the foundation of load balancers, WAFs, and security gateways. Cloudflare, Fastly, and your internal proxies all start here.

---

### Week 6: TLS & Secure Transport (Days 38–44)

**Why this week matters:** TLS is invisible when it works and catastrophic when it breaks. Cert expiry causes outages. mTLS is how zero-trust networks authenticate services. You need to understand this deeply.

---

**Day 38 — TLS handshake**
- **Why:** Every HTTPS request, every gRPC call, every mTLS service mesh connection starts with a TLS handshake. Knowing what fails (cert mismatch, expired cert, cipher mismatch) helps you triage incidents
- **Topic:** TLS 1.2 vs 1.3 handshake steps, certificate purpose, cipher suites, SNI (Server Name Indication), ALPN, session resumption
- **Exercise:** Write a TLS client that connects to `example.com:443` using `crypto/tls`. Print: TLS version, cipher suite, server cert chain (subject, issuer, expiry, SANs). Detect if cert expires in < 30 days.

**Day 39 — TLS server in Go**
- **Why:** Every internal service should use TLS. Understanding `tls.Config` options is how you set security policy
- **Topic:** `tls.Listen`, `tls.Config`, `MinVersion` (require TLS 1.2+), cipher suite selection, `InsecureSkipVerify` (never in prod — understand why), certificate loading
- **Exercise:** Build a TLS echo server using a self-signed cert. Client verifies the cert. Force TLS 1.3 only. Log the cipher suite and TLS version for each connection. Test with `openssl s_client`.

**Day 40 — Certificate generation in Go**
- **Why:** Internal services and dev environments need certs. Knowing how to generate them in Go means you can automate cert provisioning without `openssl` scripts
- **Topic:** `crypto/x509`, `crypto/rsa`, `crypto/ecdsa`, `encoding/pem`, `x509.Certificate` fields, `x509.CreateCertificate`, self-signed vs CA-signed
- **Exercise:** Write a Go program that generates: (1) root CA cert + key, (2) server cert signed by CA, (3) client cert signed by CA. Save all as PEM files. Verify the chain.

**Day 41 — Mutual TLS (mTLS)**
- **Why:** mTLS is the foundation of zero-trust architectures. Service meshes (Istio, Linkerd) and internal APIs use it. You'll configure and debug mTLS in these products
- **Topic:** `tls.RequireAndVerifyClientCert`, `ClientCAs` pool, client cert validation, `VerifyPeerCertificate` hook, mTLS vs API keys
- **Exercise:** Build an mTLS server: require client cert, verify it's signed by your CA, extract the client's Common Name (use it as identity), reject unknown clients. Build a matching client.

**Day 42 — Hashing & HMAC**
- **Why:** Request signing (AWS SigV4, webhook validation), token integrity, and content fingerprinting all use HMAC. You'll implement and review these in real systems
- **Topic:** `crypto/sha256`, `crypto/sha512`, `crypto/hmac`, `crypto/subtle.ConstantTimeCompare` (timing attack prevention), when to use which hash
- **Exercise:** Build a webhook request signing system: server signs outgoing webhooks with HMAC-SHA256. Client verifies. Demonstrate why `bytes.Equal` is wrong (timing attack) and `subtle.ConstantTimeCompare` is right.

**Day 43 — Symmetric encryption**
- **Why:** Session tokens, at-rest encryption, and encrypted logs all use AES. Understanding AES-GCM vs AES-CBC helps you make (and review) crypto choices
- **Topic:** `crypto/aes`, `cipher.NewGCM`, nonce/IV requirements, `crypto/rand` for nonce generation, authenticated encryption (AEAD), why GCM > CBC for new designs
- **Exercise:** Build an encrypted audit log: each log entry is AES-256-GCM encrypted with a per-entry nonce. Include a timestamp in the plaintext. Decrypt and verify on read. The key comes from an env var.

**Day 44 — Asymmetric crypto & token signing**
- **Why:** JWT, OAuth tokens, and mTLS certs all use asymmetric crypto. When a design says "sign tokens with ECDSA," you'll know what that means and why ECDSA P-256 > RSA 2048 for this
- **Topic:** `crypto/ecdsa`, `crypto/elliptic` (P-256), digital signatures, sign/verify, `crypto/rand` for key generation, why ECDSA for tokens
- **Exercise:** Build a token system: `Sign(claims map[string]string) (token string, err error)` using ECDSA P-256. `Verify(token string) (claims map[string]string, err error)`. Tokens expire after 1 hour. No JWT library.

**Week 6 Mini Project: mTLS Service Gateway**
- Build a gateway that enforces mTLS between clients and a backend
- Features:
  - Generate CA, server cert, and client certs in Go (at startup if not present)
  - Require client cert; extract identity from Common Name
  - Forward authenticated requests to upstream HTTP server
  - Inject `X-Client-Identity` header with the client's CN
  - Reject unauthorized clients with `403`
  - Hot-reload certs on `SIGHUP` without dropping connections
  - Signed session tokens (ECDSA) issued after successful mTLS auth
- Key packages: `crypto/tls`, `crypto/x509`, `crypto/ecdsa`, `net/http`, `os/signal`
- **Engineer Takeaway:** This is the core of a zero-trust service-to-service auth system. Istio, SPIFFE, and mutual TLS service meshes all implement this pattern. You've now built it from scratch.

---

### Week 7: L4 Load Balancing & Resilience (Days 45–51)

**Why this week matters:** Load balancers sit in front of every production service. Understanding how they work helps you evaluate AWS NLB vs ALB vs HAProxy vs your own solution.

---

**Day 45 — Load balancing algorithms**
- **Why:** Round-robin, least-connections, and consistent hashing have very different behavior under load. Knowing the trade-offs helps you make the right choice for your traffic pattern
- **Topic:** Round-robin (simple, stateless), weighted round-robin, least connections (stateful, better for variable-cost requests), random, power of two choices
- **Exercise:** Implement all 4 algorithms. Simulate 10,000 requests with variable request durations. Measure load distribution across 5 backends. Which has best distribution? Which has lowest average latency?

**Day 46 — Backend health checking**
- **Why:** A load balancer that sends traffic to dead backends is worse than no load balancer. Health checking design is a real production decision
- **Topic:** Active health checks (periodic dial), passive health checks (track errors in-flight), health check interval, failure threshold, recovery threshold, circuit breaker concept
- **Exercise:** Build an active health checker: every 5s, TCP-connect to each backend. Track consecutive failures. Remove backend after 3 failures. Re-add after 2 successes. Emit events when state changes.

**Day 47 — Connection draining**
- **Why:** When you deploy a new version or scale down, connections to old backends must finish gracefully. This is what "connection draining" in your load balancer settings does
- **Topic:** Drain state machine (active → draining → closed), in-flight connection tracking, drain timeout, force-close after timeout
- **Exercise:** Add drain support to your load balancer: mark backend as draining (stop sending new connections), wait for existing connections to finish (30s timeout), then close. Test by draining while 50 connections are active.

**Day 48 — TCP load balancer (L4)**
- **Why:** AWS NLB, HAProxy in TCP mode, and Cloudflare's Spectrum are L4 load balancers. Understanding why you'd use L4 vs L7 helps you spec the right tool
- **Topic:** L4 vs L7 load balancing trade-offs, no HTTP parsing at L4, connection-level decisions only, performance advantage of L4
- **Exercise:** Assemble a complete L4 TCP load balancer: accept client connections, pick backend (least-connections), open a pooled connection to backend, bidirectional proxy. Handle backend failure mid-connection.

**Day 49 — Rate limiting at L4**
- **Why:** Connection-level rate limiting is the first line of DDoS defense. Your network team deploys this before traffic even reaches the application
- **Topic:** Per-IP connection limits, global connection limits, token bucket for connection rate, `sync/atomic` for counters
- **Exercise:** Add L4 rate limiting to the load balancer: max 100 concurrent connections per source IP, max 1000 global. New connections beyond limits are immediately rejected (RST). Track and expose metrics.

**Day 50 — Observability at L4**
- **Why:** Without metrics, you can't answer "how is the load balancer performing?" during an incident
- **Topic:** `expvar` counters for: connections accepted/rejected/active, bytes in/out per backend, health check results, connection duration histogram
- **Exercise:** Add full observability to the load balancer. Expose via `/debug/vars`. Build a simple `/stats` endpoint that returns a human-readable summary of backend health and connection counts.

**Day 51 — Resilience patterns**
- **Why:** Circuit breakers, retries, and timeouts are how services stay up when dependencies are flaky. You'll add them to every service that calls another service
- **Topic:** Circuit breaker (closed/open/half-open states), retry with exponential backoff, timeout propagation via context, idempotency requirements for retries
- **Exercise:** Implement a circuit breaker wrapping a backend connection. Trips after 5 consecutive failures. Half-opens after 30s. Build retry logic with jitter: 3 retries, exponential backoff (100ms, 200ms, 400ms), only retry on network errors.

**Week 7 Mini Project: Production L4 Load Balancer**
- Build a complete L4 TCP load balancer:
  - Multiple backends (round-robin + least-connections, switchable)
  - Active health checking (TCP probe every 5s)
  - Connection draining (graceful backend removal)
  - Per-IP and global connection rate limiting
  - Circuit breaker per backend
  - mTLS client connections (from Week 6)
  - Full observability: `/debug/vars`, `/stats`, `/health`
  - Graceful shutdown with drain
- Key packages: `net`, `crypto/tls`, `sync`, `sync/atomic`, `context`, `expvar`, `time`
- **Engineer Takeaway:** AWS NLB bills per LCU-hour. You now understand what an LCU is, what the LB is doing for that money, and what each of its settings changes.

---

### Week 8: Review + Phase 2 Consolidation (Days 52–60)

**Why this week matters:** Consolidation is where knowledge becomes instinct. Slower pace, deeper understanding.

---

**Day 52 — TCP internals: what goes wrong**
- **Why:** `TIME_WAIT`, `CLOSE_WAIT`, and `FIN_WAIT` are states you'll see during production incidents. Understanding them prevents misdiagnosis
- **Topic:** TCP state machine deep dive, `TIME_WAIT` (normal, 2*MSL), `CLOSE_WAIT` (application bug), `FIN_WAIT_2`, `SO_LINGER`, `SO_REUSEADDR`
- **Exercise:** Deliberately create `CLOSE_WAIT` sockets (server doesn't call `conn.Close()`). Observe with `netstat`. Fix the bug. Observe `TIME_WAIT` states after graceful close. Why don't they cause problems?

**Day 53 — Binary protocol design**
- **Why:** Internal services sometimes use binary protocols for performance (gRPC uses protobuf, Kafka uses its own wire format). Understanding framing helps you review protocol proposals
- **Topic:** `encoding/binary`, length-prefix framing, fixed-header + variable-body, endianness, versioning
- **Exercise:** Design a binary protocol for a key-value store: `[1B version][1B opcode][4B key-len][key][4B value-len][value]`. Build encoder/decoder. Build a test server and client. Benchmark vs JSON.

**Day 54 — TLS certificate rotation**
- **Why:** Cert expiry is a common incident cause. Zero-downtime cert rotation requires `GetCertificate` callback, not just loading cert at startup
- **Topic:** `tls.Config.GetCertificate`, atomic cert swap, cert caching, cert expiry monitoring, OCSP stapling concept
- **Exercise:** Build a TLS server that hot-reloads its certificate: watch for a new cert file (poll `os.Stat`), swap atomically using `sync/atomic` + `unsafe.Pointer` (or `sync.RWMutex`). No connections dropped on reload.

**Day 55 — SNI routing**
- **Why:** SNI lets one IP/port serve multiple domains with different certs and backends. This is how Cloudflare, Nginx, and most CDNs work
- **Topic:** SNI (Server Name Indication), `tls.Config.GetCertificate` receives `ClientHelloInfo.ServerName`, routing table
- **Exercise:** Build an SNI router: one listener, multiple backends. Route based on the TLS SNI hostname to different upstream servers. Each hostname gets its own cert. Test with 3 different hostnames.

**Day 56 — Connection multiplexing concepts**
- **Why:** HTTP/2 multiplexes many requests over one TCP connection. gRPC does the same. Understanding this helps you size connection pools correctly
- **Topic:** HTTP/2 streams, framing, stream ID, flow control, HEADERS + DATA frames (conceptual), why H2 needs fewer connections than H1
- **Exercise:** Build a simple request multiplexer over TCP: multiple logical "streams" over one connection, identified by a 4-byte stream ID in each frame. Demultiplex on the other end.

**Day 57 — Error propagation in network code**
- **Why:** Network errors are different from application errors. Knowing `net.Error`, `os.ErrDeadlineExceeded`, and temporary vs permanent errors helps you write reliable retry logic
- **Topic:** `net.Error` interface (`Timeout()`, `Temporary()`), `io.EOF` vs `io.ErrUnexpectedEOF`, `os.ErrDeadlineExceeded`, `syscall.ECONNRESET`, error wrapping
- **Exercise:** Build an error classifier for network code: categorize every error from `net.Conn` operations into: timeout, connection reset, EOF (clean close), unexpected EOF (peer crash), other. Log each with appropriate severity.

**Day 58 — Race conditions in network servers**
- **Why:** Race conditions in network code cause data corruption and security vulnerabilities. `-race` is mandatory. Knowing common patterns helps you catch them in review
- **Topic:** Common races: closing connection from multiple goroutines, concurrent map access in connection state, deadline race with data read
- **Exercise:** Take the Week 5 TCP proxy. Run `go test -race`. Find and fix all races. Document each race: what was the bug? What was the fix? What could have gone wrong in production?

**Day 59 — Benchmark the full L4 stack**
- **Why:** "It's fast enough" needs a number. Knowing how to measure throughput and latency for your proxy is how you set SLOs
- **Topic:** Load testing with goroutines (no external tools), measuring: connections/sec, throughput (MB/s), p50/p95/p99 latency, max concurrent connections
- **Exercise:** Write a load test for the Week 7 load balancer: 1000 concurrent clients, each making 100 TCP round-trips. Measure and print: total requests, RPS, p50/p95/p99 latency, error rate.

**Day 60 — Phase 2 Review**
- Audit all Phase 2 projects for: goroutine leaks, missing `defer conn.Close()`, improper error handling
- Run `go test -race ./...` on all projects
- Document: what would break first under 10x traffic? What would you change?
- **Engineer Takeaway:** You now understand L4 networking deeply enough to design a solution, debug it during an incident, and review other people's designs.

**Week 8 Mini Project: Secure TCP Gateway**
- Combine all Phase 2 work into a secure TCP gateway:
  - mTLS client termination (from Week 6)
  - L4 load balancing to multiple backends (from Week 7)
  - Per-IP rate limiting
  - Circuit breaker per backend
  - SNI routing (different backends by hostname)
  - Health checking with expiry
  - Full observability (`/debug/vars`, `/stats`, `/health`)
  - Graceful shutdown with cert hot-reload on SIGHUP
- Key packages: `net`, `crypto/tls`, `crypto/x509`, `sync`, `sync/atomic`, `expvar`
- **Engineer Takeaway:** This is roughly what Cloudflare Spectrum does for TCP proxying. You've built it.

---

## Phase 3: L7 Networking + WAF + API Security (Days 61–90)

> **Phase Goal:** Move up to the application layer. Build HTTP middleware, a Web Application Firewall, and API security controls. These are the systems your platform security team owns.

### Week 9: HTTP Internals (Days 61–67)

**Why this week matters:** HTTP is the protocol for every API your company exposes. Understanding how it works at the wire level explains every HTTP security vulnerability.

---

**Day 61 — HTTP/1.1 from scratch**
- **Why:** Parsing HTTP manually reveals why HTTP request smuggling, header injection, and request splitting attacks work. You can't evaluate these risks without knowing the protocol
- **Topic:** Request line, headers, body, `Content-Length` vs `Transfer-Encoding: chunked`, response format, persistent connections (`Connection: keep-alive`), pipelining
- **Exercise:** Build an HTTP/1.1 server using only `net.Listen` — no `net/http`. Parse request line and headers manually using `bufio.Reader`. Serve `GET /` → `200 OK`. Handle `Connection: close` vs keep-alive.

**Day 62 — net/http server**
- **Why:** `net/http` is what you'll use for every Go service. Understanding its config options (timeouts, limits) prevents slowloris and resource exhaustion
- **Topic:** `http.Server`, `ServeMux`, `http.Handler`, `http.HandlerFunc`, `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, `MaxHeaderBytes`
- **Exercise:** Build a hardened HTTP server: set all timeouts (read header: 5s, read body: 30s, write: 30s, idle: 120s), limit max header size, log connection state transitions using `ConnState` callback.

**Day 63 — HTTP middleware**
- **Why:** Every cross-cutting concern (auth, logging, rate limiting, WAF) is implemented as middleware. Understanding the pattern helps you review middleware ordering decisions
- **Topic:** Middleware as `func(http.Handler) http.Handler`, handler chain, middleware ordering (matters for auth before logging), `http.ResponseWriter` wrapping to capture status code
- **Exercise:** Build a middleware chain builder: `Chain(h, mw1, mw2, mw3)`. Implement: (1) request ID injector, (2) latency logger (captures status code), (3) panic recoverer. Verify ordering.

**Day 64 — HTTP reverse proxy**
- **Why:** Nginx, Envoy, and your internal API gateway are all reverse proxies. Understanding how they work helps you configure them correctly and debug routing issues
- **Topic:** `httputil.ReverseProxy`, `Director` (modify request before forwarding), `ModifyResponse` (modify response before returning), hop-by-hop headers, `X-Forwarded-For`, `X-Real-IP`
- **Exercise:** Build a reverse proxy: forward to upstream, rewrite Host header, strip sensitive upstream headers (`X-Internal-*`), inject `X-Request-ID`, `X-Forwarded-For`. Log: method, path, upstream latency, status.

**Day 65 — HTTP/2**
- **Why:** HTTP/2 is the default for gRPC and many modern APIs. H2 multiplexing changes how you think about connection limits and load balancing
- **Topic:** H2 vs H1 differences (multiplexing, header compression, binary framing), `net/http` H2 support (automatic with TLS), why H2 needs fewer connections, H2 cleartext (`h2c`)
- **Exercise:** Build a server that supports both HTTP/1.1 and HTTP/2 (auto-negotiate via ALPN). Log which protocol each request uses. Send 100 requests via H1 and H2 — measure latency difference with `time`.

**Day 66 — WebSocket**
- **Why:** WebSocket powers real-time features (dashboards, alerts, live config updates). Understanding the upgrade handshake explains why some WAF rules break WebSocket connections
- **Topic:** WebSocket upgrade (HTTP 101), `Sec-WebSocket-Key`/`Accept` handshake (SHA-1), frame format (opcode, mask, payload), ping/pong, close handshake
- **Exercise:** Build a WebSocket server from scratch (no library): parse the Upgrade request, compute `Sec-WebSocket-Accept`, read/write frames manually. Build a real-time event broadcaster: clients subscribe, server fans out messages.

**Day 67 — HTTP security headers**
- **Why:** Security headers (HSTS, CSP, CORS) are the browser-side security controls every service should set. You'll catch PRs that miss them
- **Topic:** HSTS, CSP (Content-Security-Policy), CORS (preflight, `Access-Control-Allow-*`), `X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy`
- **Exercise:** Build a security headers middleware: configurable CORS (allowlist-based, not `*`), HSTS with preload, CSP with per-request nonce (via `crypto/rand`), automatic OPTIONS preflight handling.

**Week 9 Mini Project: HTTP Reverse Proxy with Full Observability**
- Build a production-grade reverse proxy:
  - HTTP/1.1 and HTTP/2 (TLS, auto-negotiate)
  - Middleware chain: request ID → security headers → logging → proxy
  - Backend health checks (HTTP GET `/health` every 5s)
  - Request timeout enforcement (per-request context with deadline)
  - Response modification: strip internal headers, inject tracing headers
  - Metrics: RPS, error rate, p50/p95/p99 latency (ring buffer of 10k samples)
  - `/metrics` (expvar), `/health`, `/debug/pprof`
  - Graceful shutdown with connection drain
- Key packages: `net/http`, `net/http/httputil`, `crypto/tls`, `sync`, `expvar`, `time`
- **Engineer Takeaway:** Nginx at its core does exactly this. Every request through your company's API gateway goes through a chain like this.

---

### Week 10: Web Application Firewall (Days 68–74)

**Why this week matters:** WAF is often the first security control evaluated for compliance. Understanding how rules work helps you configure a WAF, write precise exceptions, and reason about false positive tradeoffs.

---

**Day 68 — WAF fundamentals & OWASP Top 10**
- **Why:** OWASP Top 10 shows up in SOC2, PCI, and ISO 27001 audits. Knowing how WAF rules map to OWASP categories helps you have informed compliance conversations
- **Topic:** OWASP Top 10 categories, WAF modes (detection vs blocking), where WAF sits in the stack, false positive vs false negative tradeoffs, paranoia levels
- **Exercise:** Build a WAF skeleton: `type Rule interface { ID() string; Match(r *http.Request) (matched bool, score int, detail string) }`. Build `WAFMiddleware(rules []Rule, threshold int) http.Handler`. Log matches, block if score > threshold.

**Day 69 — SQL injection detection**
- **Why:** SQLi is the #1 data breach vector. WAF rules are your last line of defense when app code has a bug. Understanding detection helps you evaluate rule quality
- **Topic:** SQLi payloads (`' OR '1'='1`, `UNION SELECT`, `--` comments, stacked queries), tokenization vs regex detection, false positive patterns (legitimate SQL in app names)
- **Exercise:** Implement SQLi detection: (1) keyword list (`UNION`, `SELECT`, `DROP`, `INSERT`, `--`), (2) quote-balance detection, (3) comment pattern detection. Test with 15 SQLi payloads + 15 legitimate inputs. Tune to 0 false positives on legitimate inputs.

**Day 70 — XSS detection**
- **Why:** XSS leads to session hijacking and account takeover. Understanding detection patterns helps you tune WAF rules vs CSP headers
- **Topic:** Reflected XSS (in response), stored XSS (in DB), DOM XSS (in JS), `<script>` injection, event handlers (`onerror=`, `onload=`), JavaScript URIs (`javascript:`), SVG attacks
- **Exercise:** Build XSS detection: scan URL params, headers, and JSON body. Normalize input first (URL-decode, HTML-unescape). Match: `<script`, `javascript:`, event handler patterns (`on[a-z]+=`). Test with 10 payloads + 10 legitimate HTML fragments.

**Day 71 — Path traversal & command injection**
- **Why:** Path traversal lets attackers read `/etc/passwd`. Command injection lets them run arbitrary code. These are critical severity vulnerabilities
- **Topic:** Path traversal patterns (`../`, `%2e%2e%2f`, `%252e%252e%252f` double-encoding), null bytes (`%00`), command injection (`;`, `|`, `` ` ``, `$(`, `&&`)
- **Exercise:** Build detection for path traversal (handle URL-encoding, double-encoding, unicode variants) and command injection. Add a normalization step: URL-decode → percent-decode → normalize. Test bypass attempts with encoded variants.

**Day 72 — Request body inspection**
- **Why:** Most WAFs only scan URL and headers by default. Body scanning is where complex attacks hide. Understanding the tradeoffs (latency vs coverage) is a product decision
- **Topic:** Reading body without breaking the handler (`io.TeeReader` + buffer), `Content-Type` dispatch (form, JSON, multipart), body size limits, streaming vs buffered inspection
- **Exercise:** Build a body inspector: buffer up to 1MB of body (reject larger), dispatch by Content-Type for scanning, restore body for downstream handler using `io.NopCloser(bytes.NewReader(...))`. Scan JSON values recursively.

**Day 73 — WAF normalization pipeline**
- **Why:** Attackers bypass WAF rules using encoding. Multi-layer normalization is what separates basic WAFs from production-grade ones
- **Topic:** Multi-layer decoding (URL → HTML → Unicode → lowercase), SQL comment stripping (`/**/`), whitespace normalization, null byte removal, Unicode confusables
- **Exercise:** Build a normalization pipeline: `input → URL-decode → HTML-unescape → remove SQL comments → collapse whitespace → lowercase`. Apply before every rule match. Test 10 encoding bypass attempts.

**Day 74 — WAF rule scoring & exceptions**
- **Why:** Block-mode WAFs cause false positive incidents. Scoring-based WAFs (block only if multiple rules match) reduce false positives. Exception lists are how you handle legitimate traffic that looks like attacks
- **Topic:** Score accumulation (each rule adds weight), group scores (SQLi group, XSS group), block threshold per group, exception lists (by IP, path, header value)
- **Exercise:** Build a scoring WAF: SQLi group (block if SQLi score > 5), XSS group (block if XSS score > 4), path traversal (block on any match). Exception list: trusted IPs bypass all rules. Test with mixed payloads.

**Week 10 Mini Project: Production WAF Middleware**
- Build a full WAF middleware:
  - Normalization pipeline (URL-decode → HTML-unescape → lowercase → comment strip)
  - Rules: SQLi (scoring), XSS (scoring), path traversal (immediate block), command injection (immediate block)
  - Body inspection (JSON + form, up to 1MB)
  - Allow-list: trusted CIDRs bypass WAF
  - Block response: `403` with request ID (no internal details leaked to attacker)
  - Audit log: every match logged with rule ID, matched value, normalized value, score, action
  - Config hot-reload (rule thresholds from JSON file, poll every 30s)
  - Detection mode: log but don't block (for tuning)
- Key packages: `net/http`, `regexp`, `strings`, `html`, `net`, `sync`, `encoding/json`
- **Engineer Takeaway:** AWS WAF, Cloudflare WAF, and ModSecurity all implement this scoring model. When a WAF blocks a legitimate request, you'll know how to find the rule and write a narrow exception, and what risk it carries.

---

### Week 11: API Security & Authentication (Days 75–81)

**Why this week matters:** API security is where most breaches happen today. OAuth2, API keys, rate limiting, and input validation are everyday concerns for anyone building APIs.

---

**Day 75 — API authentication patterns**
- **Why:** Auth is the most critical security control. Understanding the tradeoffs between API keys, HMAC signing, and JWTs helps you make the right architecture decisions
- **Topic:** API key (simple, no expiry), HMAC-signed requests (AWS SigV4 pattern — key + timestamp + body hash), signed tokens (JWT-like, expiry built in), mTLS (certificate as identity)
- **Exercise:** Implement all 3 auth patterns as middleware: (1) API key header, (2) HMAC-SHA256 request signing (include timestamp, method, path, body hash in signature), (3) ECDSA-signed token (from Phase 2). Each rejects unauthorized requests with `401`.

**Day 76 — Rate limiting at L7**
- **Why:** API rate limiting protects backend systems from abuse. Token bucket is the algorithm used by AWS API Gateway, Stripe, and most API platforms
- **Topic:** Token bucket (tokens refill at rate R, burst up to capacity B), per-user vs per-IP vs per-API-key limiting, rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`)
- **Exercise:** Build rate limiting middleware: per-API-key token bucket (10 req/s, burst 50). Respond `429` with `Retry-After` header. Use `sync.Map` for per-key buckets + `time.NewTicker` for background refill. Test with 100 concurrent clients.

**Day 77 — Input validation**
- **Why:** Unvalidated input is the root cause of injection attacks. Go struct tags for validation are in every serious API codebase
- **Topic:** Allowlist vs denylist validation, length limits, type coercion, required fields, pattern matching, nested struct validation, validation error aggregation
- **Exercise:** Build a struct validation framework using `reflect` and struct tags: `validate:"required,max=255,pattern=^[a-zA-Z]+$"`. Return all validation errors (not just first). Test with 20 valid + 20 invalid inputs.

**Day 78 — JSON schema enforcement**
- **Why:** API schema validation prevents malformed requests from reaching application code. OpenAPI/JSON Schema is the industry standard
- **Topic:** JSON Schema spec subset (`type`, `required`, `minLength`, `maxLength`, `minimum`, `maximum`, `pattern`, `enum`, `additionalProperties: false`)
- **Exercise:** Implement a JSON Schema validator (no library): validate incoming JSON request bodies against a schema defined in Go. Reject unknown fields (important for security — prevents parameter pollution). Return structured errors.

**Day 79 — Sensitive data protection**
- **Why:** PII in logs is a compliance violation (GDPR, CCPA). Response scrubbing and log masking are often compliance requirements. You'll review data handling designs
- **Topic:** PII detection patterns (email, phone, SSN, credit card), log field masking (`"password":"***"`), response field redaction, `encoding/json` streaming for large payloads
- **Exercise:** Build a PII scrubber: (1) log middleware that masks known-sensitive field names (`password`, `token`, `ssn`, `card_number`) in JSON bodies before logging, (2) response scrubber that redacts patterns (credit card regex, email regex) from JSON responses.

**Day 80 — API abuse prevention**
- **Why:** Credential stuffing, account enumeration, and scraping are attacks against APIs, not just web UIs. Understanding the patterns helps you design effective controls
- **Topic:** Credential stuffing (try leaked credentials), account enumeration (error message differences), brute force (repeated auth attempts), scraping (extract all data via API)
- **Exercise:** Build defenses: (1) failed auth rate limiting per username + per IP independently, (2) constant-time auth response (always ~200ms via `time.Sleep` — prevent timing attacks), (3) honeypot API endpoint that flags any IP that calls it.

**Day 81 — CORS & preflight**
- **Why:** CORS misconfigurations are in the OWASP Top 10. `Access-Control-Allow-Origin: *` on an authenticated API is a critical vulnerability
- **Topic:** Same-origin policy (why CORS exists), preflight (OPTIONS), `Access-Control-Allow-Origin` (allowlist, never `*` for auth'd APIs), `Access-Control-Allow-Credentials`, `Vary: Origin`
- **Exercise:** Build a CORS middleware: strict origin allowlist (no wildcard for credentialed requests), preflight caching (`Access-Control-Max-Age`), `Vary: Origin` header, reject disallowed origins with `403`. Test with 5 allowed + 5 disallowed origins.

**Week 11 Mini Project: API Security Gateway**
- Build a complete API security layer:
  - Auth: ECDSA-signed token validation + HMAC request signing support
  - Rate limiting: per-API-key token bucket, per-IP sliding window
  - JSON schema validation per route (loaded from config file)
  - Input validation (struct tags) for known request types
  - Response PII scrubber (mask sensitive fields in logs and responses)
  - Abuse prevention: failed auth rate limiting, honeypot endpoint
  - CORS with origin allowlist
  - Security headers (HSTS, CSP, X-Frame-Options)
  - WAF (from Week 10) integrated in the chain
  - Full audit log (request + masked body + response status + latency)
- Key packages: `net/http`, `crypto/ecdsa`, `crypto/hmac`, `encoding/json`, `sync`, `regexp`
- **Engineer Takeaway:** This is the security policy of your API platform codified in software. Every feature maps to a compliance control or a real attack vector.

---

### Week 12: Phase 3 Review + Integration (Days 82–90)

**Why this week matters:** Integration week. Wire the L4 gateway (Phase 2) to the HTTP/API gateway (Phase 3) and validate the complete stack under load.

---

**Day 82 — Middleware ordering & pipeline design**
- **Why:** The order of security middleware matters. Auth before WAF? Rate limit before auth? Wrong order creates vulnerabilities or poor UX (user gets WAF-blocked before seeing auth error)
- **Topic:** Correct ordering: IP allowlist → connection limits → rate limit → TLS/bot check → auth → WAF → schema validation → handler. Short-circuit semantics. Context propagation.
- **Exercise:** Build a configurable middleware pipeline that loads order from a JSON config. Validate that required middleware are present. Log which middleware handled (rejected) each request.

**Day 83 — Threat intelligence aggregation**
- **Why:** IP reputation is a shared signal across WAF, DDoS, and bot detection. A unified threat score per IP improves accuracy and reduces false positives
- **Topic:** IP reputation store, score aggregation from multiple signals (WAF hits, rate limit violations, failed auth), score decay over time, block threshold vs challenge threshold
- **Exercise:** Build a `ThreatIntel` store: accepts signals from any module (`ReportSignal(ip, type, weight)`), aggregates scores with time decay (score halves every 10 minutes), exposes `Score(ip) float64` and `Action(ip) string` (allow/challenge/block).

**Day 84 — Structured logging with slog**
- **Why:** `log/slog` is stdlib since Go 1.21. Structured logs are machine-parseable. Your observability platform (Datadog, Splunk) expects JSON. Use slog in every service you ship
- **Topic:** `slog.Logger`, `slog.Handler`, JSON handler (`slog.NewJSONHandler`), log levels, `slog.With` (add context fields), correlation ID propagation via context
- **Exercise:** Replace all `log.Printf` in your projects with `slog`. Add: request ID (from context), client IP, user identity (from auth middleware), module name. Output JSON. Add log level filtering (debug logs off in prod).

**Day 85 — Configuration & hot reload**
- **Why:** Changing WAF rules or rate limits without restarting a service is a production requirement. Understanding the patterns helps you design config management
- **Topic:** JSON config loading, env var override, config validation on load, hot reload via file watching (`os.Stat` polling or `SIGHUP`), atomic config swap (`sync.RWMutex` or `atomic.Value`)
- **Exercise:** Build a config system for the API gateway: JSON file with WAF thresholds, rate limits, CORS origins, auth config. Hot-reload on `SIGHUP`. Validate new config before swapping (reject invalid config, keep old).

**Day 86 — Integration testing the full chain**
- **Why:** Unit tests don't catch middleware ordering bugs. Integration tests using `net/http/httptest` or `net.Pipe` catch real issues
- **Topic:** `httptest.NewServer`, `httptest.NewTLSServer`, `net.Pipe` for TCP tests, test scenarios: auth bypass attempts, WAF evasion, rate limit exhaustion
- **Exercise:** Write integration tests for the full middleware chain: (1) valid request passes all layers, (2) SQL injection blocked by WAF, (3) rate limit enforced, (4) invalid token rejected, (5) oversized body rejected, (6) CORS blocked for unknown origin.

**Day 87 — Fuzz testing**
- **Why:** Fuzz testing finds crashes and security bugs that hand-written tests miss. Go's built-in fuzzer is production-grade and requires no setup
- **Topic:** `testing/fuzz`, `f.Add` seed corpus, `f.Fuzz`, what to fuzz (parsers, decoders, validators), interpreting results, fixing panics
- **Exercise:** Write fuzz tests for: (1) WAF normalization pipeline (does it panic on arbitrary input?), (2) JSON schema validator (does malformed JSON panic?), (3) HTTP request parser (Day 61). Run with `go test -fuzz -fuzztime 30s`.

**Day 88 — Load testing the security stack**
- **Why:** Security middleware adds latency. Knowing the overhead helps you justify the cost and set SLOs. "Our WAF adds 2ms p99" is a concrete business number
- **Topic:** Load generation with goroutines, measuring latency percentiles (ring buffer approach), connection reuse vs new connection per request, what to measure
- **Exercise:** Load test the full API gateway: 500 concurrent clients, 10,000 total requests, mix of valid + invalid (10% WAF triggers, 5% rate limit hits, 5% auth failures). Report: RPS, p50/p95/p99 latency, per-middleware overhead, error rate breakdown.

**Day 89 — Goroutine leak audit**
- **Why:** Goroutine leaks cause memory growth and eventually OOM. They're common in middleware that starts goroutines for timeouts or background tasks
- **Topic:** `runtime.NumGoroutine()`, goroutine dumps (`runtime.Stack`), leak patterns (forgotten `done` channel, missing `context` cancellation, goroutine with no exit path)
- **Exercise:** Add goroutine count tracking to all middleware. Run the load test (Day 88). After load test ends, check goroutine count — it should return to baseline within 10s. Find and fix any leaks.

**Day 90 — Phase 3 Review & architecture diagram**
- Review the complete Phase 3 stack
- Draw an architecture diagram (ASCII) showing all components, data flow, and security controls
- Identify: what would you use a managed product for vs build? (WAF rules database, TLS certificate management, IP reputation feeds)
- **Engineer Takeaway:** You can now design, build, and operate a full API security gateway. You understand what a product sold as a "WAF" or "API gateway" does, line by line.

**Week 12 Mini Project: Unified Security Proxy (L4 + L7)**
- Wire the Phase 2 L4 gateway to the Phase 3 API security layer:
  - L4: TLS termination, mTLS, SNI routing, L4 rate limiting, circuit breaker
  - L7: HTTP/2, reverse proxy, WAF, API auth, rate limiting, schema validation
  - Unified threat intel (`ThreatIntel` store, shared across L4 and L7)
  - Structured logging (slog, JSON, correlation IDs)
  - Hot-reload config (SIGHUP)
  - Full observability: `/debug/vars`, `/debug/pprof`, `/health`, `/metrics`
  - Graceful shutdown (3-stage)
- Key packages: All stdlib
- **Engineer Takeaway:** This is the architecture of Cloudflare's edge, AWS WAF + ALB, and Fastly's compute platform — implemented in ~2000 lines of Go.

---

## Phase 4: DDoS + Bot Detection + Capstone (Days 91–120)

> **Phase Goal:** Add the final security layers: DDoS mitigation and bot detection. Then build the capstone — a production-ready, single-binary security proxy.

### Week 13: DDoS Detection & Mitigation (Days 91–98)

**Why this week matters:** DDoS attacks are the most common attack against internet-facing services. Understanding the mechanics helps you configure DDoS protection services (Cloudflare, AWS Shield, Akamai) well and build your own first line of defense.

---

**Day 91 — DDoS taxonomy**
- **Why:** "We're under DDoS" has a very different response depending on the attack type. L3/L4 volumetric vs L7 application attacks require different mitigations
- **Topic:** Volumetric (bandwidth exhaustion), Protocol (SYN flood, connection exhaustion), Application (HTTP flood, slowloris, slow POST), Amplification (UDP reflection — DNS, NTP, memcached)
- **Exercise:** For each DDoS type, write a Go program that simulates the attack pattern (at small scale, against localhost): (1) TCP connection flood (10k connections, no data), (2) HTTP flood (1000 req/s), (3) slowloris (100 connections, headers drip at 1 byte/sec).

**Day 92 — Rate limiting algorithms (deep dive)**
- **Why:** Token bucket vs sliding window have different behavior under bursty traffic. The right algorithm for DDoS is different from the right algorithm for API rate limiting
- **Topic:** Fixed window (fast, inaccurate at boundaries), sliding window log (accurate, memory heavy), sliding window counter (approximation, memory efficient), token bucket (allows burst), leaky bucket (smooths traffic)
- **Exercise:** Implement all 5. Test with a traffic pattern: 100 req/s normally, spike to 1000 req/s for 5 seconds, back to 100. Which algorithms block the spike cleanly? Which have boundary issues? Which allow the most burst?

**Day 93 — Per-IP rate limiting with CIDR awareness**
- **Why:** DDoS attackers use IP ranges (ASNs), not just single IPs. Rate limiting by /24 or /16 blocks more attack traffic
- **Topic:** Per-IP counters (`sync.Map`), CIDR matching (`net.IPNet`), X-Forwarded-For extraction (with trust chain), IPv6 considerations (/48 prefix grouping), counter expiry (memory cleanup)
- **Exercise:** Build per-IP and per-/24-subnet rate limiting. Handle both IPv4 and IPv6. Extract real IP from `X-Forwarded-For` (trust only 1 hop). Clean up expired counters with a background goroutine. Test with 1000 unique IPs.

**Day 94 — SYN flood defense**
- **Why:** SYN floods are still the most common L4 DDoS. Understanding them helps you configure `net.ipv4.tcp_syncookies` and size your TCP backlog
- **Topic:** SYN cookie concept, half-open connection limits, accept queue limits, Go-side defense (limit concurrent `Accept` goroutines, connection per IP limit), backlog tuning
- **Exercise:** Implement connection-level SYN flood defense: max 5 half-open connections per source IP (track IPs that connected but haven't sent data yet). Drain half-open connections after 3s. Track and expose metrics.

**Day 95 — Slowloris defense**
- **Why:** Slowloris can take down Apache/nginx with a single machine. Go's `ReadHeaderTimeout` partially defends against it, but you need to detect and block the source IP too
- **Topic:** Slowloris mechanics (send headers slowly, keep connection open), detection (header arrival rate < threshold), `ReadHeaderTimeout` is not enough (attackers rotate), connection rate per IP
- **Exercise:** Build slowloris detection: track bytes received per second per connection. If a connection sends < 50 bytes/second during header phase for > 5 seconds, classify as slowloris, close connection, add source IP to block list (5-minute expiry).

**Day 96 — Traffic analysis & anomaly detection**
- **Why:** DDoS traffic often looks legitimate at the request level but anomalous at the traffic level (too many requests from one IP, unusual user-agents, single endpoint targeted). Statistical detection is more accurate than static thresholds
- **Topic:** Sliding window request rates (per-IP, per-path, global), baseline estimation, z-score anomaly detection, exponential moving average
- **Exercise:** Build a traffic analyzer: maintain per-IP request rate with sliding window (60s). Compute global mean and stddev using exponential moving average. Flag IPs > 3 stddev above mean. Emit alert with IP, rate, and deviation score.

**Day 97 — IP blocking & challenges**
- **Why:** Outright blocking causes false positives. Challenges (CAPTCHA-like, PoW) let legitimate traffic through while slowing bots. Understanding the tradeoff is a product decision
- **Topic:** IP blocklist (`sync.Map` + expiry), `429 Too Many Requests` with `Retry-After`, proof-of-work challenge (SHA256 puzzle), challenge verification, whitelist (trusted CIDRs bypass)
- **Exercise:** Build a tiered response: (1) suspected IPs (score 50–80): 429 + SHA256 PoW challenge, (2) confirmed attackers (score > 80): 403 block for 5 minutes, (3) trusted IPs: bypass. Implement PoW: server issues random challenge, client must find nonce such that SHA256(challenge+nonce) has N leading zero bits.

**Day 98 — UDP amplification defense**
- **Why:** NTP, DNS, and memcached amplification attacks are still common. Understanding the mechanics helps you design UDP services that can't be weaponized
- **Topic:** Amplification factor (response >> request size), rate limiting UDP responses per source IP, response size caps, requiring handshake (challenge-response) before large responses
- **Exercise:** Build a UDP service with anti-amplification: (1) response size cap (response ≤ 1.5× request size), (2) rate limit responses per source IP (10/sec), (3) reject requests that trigger > 2 rate limit violations (add to block list). Measure: can your service be used for amplification?

**Week 13 Mini Project: DDoS Mitigation Layer**
- Build a complete DDoS mitigation proxy:
  - Per-IP and per-/24 rate limiting (token bucket, sliding window)
  - Slowloris detection (byte rate per connection)
  - SYN flood defense (per-IP half-open connection limit)
  - Traffic anomaly detection (z-score, EMA baseline)
  - Tiered response: challenge (PoW) → block → whitelist
  - IP blocklist with expiry (`sync.Map`, background cleanup)
  - Real-time stats: `/debug/ddos` (top IPs, current block list, attack status, request rates by source)
  - CIDR whitelist (bypass mitigation for trusted ranges)
  - Integrate with `ThreatIntel` store from Phase 3
- Key packages: `net`, `net/http`, `sync`, `sync/atomic`, `crypto/sha256`, `time`, `net` (IPNet)
- **Engineer Takeaway:** Cloudflare Magic Transit, AWS Shield Advanced, and Akamai Prolexic implement this logic at 100Tbps scale. You've now built the core algorithms from scratch.

---

### Week 14: Bot Detection (Days 99–105)

**Why this week matters:** Bot traffic accounts for 40–50% of internet traffic. Bots scrape prices, credential-stuff accounts, and abuse APIs. Bot detection is a multi-signal problem that requires understanding both network behavior and application behavior.

---

**Day 99 — Bot taxonomy & signal framework**
- **Why:** Not all bots are bad (Googlebot is welcome, scrapers are not). Signal-based detection with configurable thresholds lets you tune aggression. Understanding signals helps you integrate and tune bot detection products (Cloudflare Bot Management, DataDome, PerimeterX)
- **Topic:** Bot types: simple HTTP clients, headless browsers (Puppeteer/Playwright), residential proxy bots, human-like bots. Signal categories: network signals (IP, ASN, TLS), HTTP signals (headers, timing), behavioral signals (session patterns)
- **Exercise:** Design a `BotSignal` scoring system: `type Signal struct { Name string; Weight float64; Category string }`. Implement `BotScorer` that accumulates signals for each request/session. Score > threshold → challenge/block.

**Day 100 — HTTP header fingerprinting**
- **Why:** Real browsers send headers in a consistent order with specific values. Bots often send wrong Accept headers, missing Accept-Language, or User-Agents that don't match the Accept-Encoding
- **Topic:** Browser header consistency (Chrome always sends headers in the same order), User-Agent parsing, Accept-Encoding consistency (Chrome never sends `compress`), Accept-Language presence
- **Exercise:** Build a header consistency scorer: for each User-Agent claiming to be a specific browser, check: (1) Accept-Encoding matches expected values, (2) Accept-Language is present, (3) header count is within normal range. Score anomalies.

**Day 101 — TLS fingerprinting (JA3)**
- **Why:** JA3 fingerprints the TLS ClientHello — cipher suites, extensions, elliptic curves. Real Chrome has a specific JA3. `curl` and Python `requests` have different JA3s. This is hard to fake without modifying the TLS stack
- **Topic:** TLS ClientHello structure (cipher suites, extensions list, elliptic curves, EC point formats), JA3 algorithm (sort + join → MD5), known JA3 hashes for Chrome/Firefox/curl
- **Exercise:** Implement JA3 fingerprinting: intercept TLS ClientHello via `tls.Config.GetConfigForClient` (which receives `*tls.ClientHelloInfo`). Compute JA3 hash. Compare against known hashes for major browsers. Flag non-browser TLS stacks.

**Day 102 — Session behavioral analysis**
- **Why:** Bots interact with APIs in non-human patterns: too fast, no referrer chain, no session depth, single endpoint focus. Session tracking exposes these patterns
- **Topic:** Session tracking (session ID via cookie or token), request velocity (requests/second), path diversity (how many unique paths in session), referrer chain (direct API calls vs navigation)
- **Exercise:** Build a session tracker: per session, track: request count, unique paths, requests/second (sliding window), referrer presence rate. Score sessions with: velocity > 10 req/s (+20), referrer missing > 80% (+15), path diversity < 2 paths (+10).

**Day 103 — IP reputation & ASN classification**
- **Why:** Bots disproportionately originate from cloud provider IPs (AWS, GCP, Azure, DigitalOcean) and residential proxy networks. IP classification gives a strong first-pass signal
- **Topic:** ASN-to-IP CIDR mapping (embedded hardcoded ranges), datacenter IP classification, residential proxy patterns, VPN exit node indicators
- **Exercise:** Build an IP classifier: embed CIDR ranges for major cloud providers (AWS: `3.0.0.0/8`, etc. — hardcode a representative subset). Score requests: datacenter IP (+15), known hosting provider (+20). Also detect: IPv6-only requests from unusual ASNs.

**Day 104 — Proof-of-work & invisible challenges**
- **Why:** PoW is a cost-imposing challenge — it slows bots without blocking legitimate users. It works because bots can't amortize computational cost across many requests the way they amortize network cost
- **Topic:** PoW challenge mechanics (SHA256 hash puzzle), difficulty levels (leading zero bits), challenge issuance (signed by server, time-limited), validation, re-challenge on failure
- **Exercise:** Build a PoW challenge system: `GET /challenge` returns `{challenge, difficulty, expires}` (challenge is random bytes, signed with HMAC-SHA256). Client sends `POST /verify` with `{challenge, nonce}`. Server verifies SHA256(challenge+nonce) has `difficulty` leading zero bits. Issue session cookie on success.

**Day 105 — Bot mitigation strategies**
- **Why:** Hard blocking bots causes false positives. Tiered mitigation (slow → challenge → block) reduces false positives while still raising the cost for bots
- **Topic:** Slow response (artificial `time.Sleep` for suspected bots), shadow-ban (return fake `200 OK` with empty data), honeypot (hidden endpoint that only bots call), hard block
- **Exercise:** Implement tiered bot mitigation: (1) score 20–40: add 500ms artificial delay, (2) score 40–70: serve PoW challenge, (3) score > 70: block with `403`, (4) honeypot: any IP that hits `/api/hidden` gets score +100. Implement shadow-ban endpoint that returns `200 OK` with plausible but empty data.

**Week 14 Mini Project: Bot Detection Engine**
- Build a complete bot detection middleware:
  - Header consistency fingerprinting
  - JA3-like TLS fingerprinting (via `tls.ClientHelloInfo`)
  - Session behavioral scoring (velocity, path diversity, referrer rate)
  - IP ASN classification (datacenter CIDR list)
  - PoW challenge issuance and verification
  - Honeypot endpoint (flags IP permanently)
  - Tiered mitigation: delay → challenge → block
  - Dashboard: `/debug/bots` (scores, flagged sessions, challenge success/fail rates, top bot IPs)
  - Integrate with `ThreatIntel` store
- Key packages: `crypto/tls`, `net`, `net/http`, `crypto/md5`, `crypto/sha256`, `sync`, `time`
- **Engineer Takeaway:** DataDome and PerimeterX sell bot detection as a product. You now understand their core algorithms, so you can integrate them, tune them, and know when a simple in-house detector is enough.

---

### Week 15: Production Hardening (Days 106–112)

**Why this week matters:** Software that works on your laptop is not production software. This week is about the operational characteristics that matter at scale.

---

**Day 106 — Health checks & readiness**
- **Why:** Kubernetes uses liveness and readiness probes to route traffic and restart unhealthy pods. Getting these right prevents cascading failures
- **Topic:** Liveness (is the process alive and not deadlocked?), readiness (is the service ready to serve traffic — upstream connections established?), startup probe (allow slow startup), `/healthz`, `/readyz`
- **Exercise:** Build a health check framework: each subsystem registers a `Check func(ctx context.Context) error`. `/healthz` runs all checks with 2s timeout per check, returns JSON with per-check status. `/readyz` checks upstream connectivity.

**Day 107 — Multi-stage graceful shutdown**
- **Why:** `kill -9` drops connections. `SIGTERM` should drain connections gracefully. Getting shutdown right prevents data loss and user-visible errors
- **Topic:** 3-stage shutdown: (1) stop accepting new connections/requests, (2) wait for in-flight to complete (with timeout budget), (3) force-close remaining. `sync.WaitGroup` for in-flight tracking. Signal handling.
- **Exercise:** Build a complete 3-stage shutdown sequence: `SIGTERM/SIGINT` triggers stage 1 (10s budget), then stage 2 (30s budget), then force-close. Log each stage. Test: send requests, send SIGTERM mid-flight, observe graceful drain.

**Day 108 — Signal handling & operational controls**
- **Why:** Operations teams need runtime controls without restarts: reload config, rotate logs, dump goroutine stack for debugging. These are table-stakes for production services
- **Topic:** `SIGHUP` (reload config), `SIGUSR1` (rotate logs), `SIGUSR2` (dump goroutine stack to log), `SIGTERM` (graceful shutdown), `signal.Notify` + goroutine dispatch
- **Exercise:** Add signal handlers to the gateway: `SIGHUP` → reload config from file, `SIGUSR1` → rotate log file, `SIGUSR2` → `runtime.Stack` dump to log. Test each signal.

**Day 109 — Observability: traces & spans**
- **Why:** Distributed tracing (Jaeger, Zipkin, Datadog APM) is how you debug latency across microservices. Understanding the W3C Trace Context standard helps you integrate with your observability stack
- **Topic:** W3C Trace Context (`traceparent` header format: `version-traceid-parentid-flags`), span lifecycle (start, add attributes, end), context propagation, in-memory trace store
- **Exercise:** Build a manual tracer: inject `traceparent` on ingress (or propagate incoming), create child spans for each middleware, record spans in memory (ring buffer of last 1000 traces). Expose at `/debug/traces` as JSON.

**Day 110 — Chaos engineering**
- **Why:** Systems fail in unexpected ways. Chaos engineering (random failures in test) reveals brittleness before production does
- **Topic:** Fault injection (random errors), latency injection (random delays), circuit breaker pattern (open after N failures, half-open after timeout, close after M successes)
- **Exercise:** Build a chaos middleware: configurable fault injection (5% chance of `500`, 2% chance of random delay 100–2000ms, 1% chance of connection close). Build a circuit breaker: open after 5 failures in 10s, half-open after 30s (pass 1 request), close after 2 successes.

**Day 111 — Memory & goroutine leak hunting**
- **Why:** Memory leaks and goroutine leaks turn a 99.9% uptime service into a service that needs weekly restarts. Production services must not leak
- **Topic:** Goroutine leak patterns (goroutine waiting on unbuffered channel with no sender, goroutine whose context is never cancelled), heap leak patterns (growing `sync.Map`, unbounded slices), `runtime.Stack`, `pprof` heap diff
- **Exercise:** Audit all projects: (1) run load test, (2) after load, check `runtime.NumGoroutine()` — should return to baseline in 10s, (3) compare heap before and after load test using `/debug/pprof/heap`. Find and fix any leaks.

**Day 112 — Zero-downtime deployment**
- **Why:** Deployments that drop connections cause user-visible errors. Zero-downtime deployment is a production requirement for high-traffic services
- **Topic:** Listener FD inheritance (old process passes listener FD to new process via `os.StartProcess`), `SO_REUSEPORT` (multiple processes share port), graceful handoff, connection draining
- **Exercise:** Implement FD inheritance: old server, on `SIGUSR1`, starts new server process with listener FD in `ExtraFiles`. New server starts serving. Old server drains and exits. No connections dropped.

**Week 15 Mini Project: Production-Hardened Gateway**
- Take the Week 12 Unified Security Proxy and add:
  - Health check framework (`/healthz`, `/readyz`) with per-subsystem checks
  - 3-stage graceful shutdown
  - Signal handlers (SIGHUP, SIGUSR1, SIGUSR2, SIGTERM)
  - W3C trace context propagation with span recording
  - Chaos middleware (configurable, off by default)
  - Goroutine and memory leak validation (automated test)
  - FD inheritance for zero-downtime restart
- Key packages: All stdlib
- **Engineer Takeaway:** This is what separates "software that works" from "production software." SRE teams live and die by these operational characteristics.

---

### Week 16: Advanced Integration + Capstone Prep (Days 113–119)

---

**Day 113 — End-to-end attack simulation**
- **Why:** You've built each defense in isolation. Testing the integrated stack against real attack scenarios is how you know it actually works
- **Topic:** Attack playbook: SQLi, XSS, path traversal, credential stuffing, SYN flood, HTTP flood, slowloris, bot scraping
- **Exercise:** Write an attack simulator: 8 goroutines, each running a different attack type against the full integrated stack. Verify: each attack is detected, correct mitigation applied (block vs challenge vs rate limit), no false positives for legitimate traffic interleaved.

**Day 114 — Performance profiling the full stack**
- **Why:** Security middleware adds overhead. Every millisecond of added latency is a business cost. Knowing the actual overhead helps you tune and justify
- **Topic:** Per-middleware latency measurement, allocation profiling (per middleware), CPU flame graph interpretation, hot path optimization with `sync.Pool`
- **Exercise:** Profile the full middleware chain under load (500 concurrent clients). Measure p50/p95/p99 latency contribution per layer. Find the 2 highest-latency middleware. Optimize each (reduce allocations, pre-compile regex). Re-measure.

**Day 115 — Security audit of your own code**
- **Why:** Adversarial review of your own code is how you build the habit of writing secure code from the start
- **Topic:** Timing attacks (auth paths), integer overflow (packet parsers), goroutine leak under attack (deliberate goroutine exhaustion), information leakage in error messages
- **Exercise:** Audit: (1) all auth paths use `subtle.ConstantTimeCompare`, (2) all packet parsers check bounds before indexing, (3) error messages to clients contain no internal details (IP, stack trace, file paths), (4) no goroutine is started without a clear exit condition.

**Day 116 — Test coverage audit**
- **Why:** Coverage numbers are metrics. Below 80% on security-critical code is a risk. Knowing what's untested helps you prioritize test writing
- **Topic:** `go test -cover`, `go test -coverprofile`, `go tool cover -html`, identifying uncovered security paths (bypass conditions, error paths)
- **Exercise:** Run coverage on all projects. For any security-critical file (WAF rules, auth, rate limiting, bot detection) below 80% — write tests to reach 80%. Focus on negative paths (what blocks the attack).

**Day 117 — Fuzz testing critical parsers**
- **Why:** WAF normalization pipelines and binary protocol parsers are exactly the code fuzzing was designed to test. A panic in your WAF is a denial-of-service
- **Topic:** `testing/fuzz` for parsers, seed corpus design, interpreting fuzz failures, fixing panics safely (bounds check, recover)
- **Exercise:** Write fuzz tests for: (1) WAF normalization pipeline, (2) JA3 TLS fingerprint parser, (3) binary protocol parser from Week 8. Run each for 60 seconds with `go test -fuzz -fuzztime=60s`. Fix any panics.

**Day 118 — Documentation & architecture review**
- **Why:** Code you can't explain to your teammates is a liability. Architecture documentation is how you transfer knowledge
- **Topic:** ASCII architecture diagrams, package-level documentation (`doc.go`), decision records (why token bucket vs sliding window for rate limiting?), runbook for each operational control
- **Exercise:** Write an architecture document for the full security proxy: ASCII diagram of all layers, description of each module, key design decisions with rationale, runbook: how to add a new WAF rule, how to block an IP, how to reload config.

**Day 119 — Capstone preparation**
- Final integration: wire all components
- Run: `go build ./...`, `go vet ./...`, `go test -race ./...`
- Fix all warnings and race conditions
- Run the attack simulator (Day 113) against the fully integrated system
- Final performance baseline: RPS at p99 < 10ms overhead

---

**Day 120 — Capstone: Cloud-Edge Security Proxy**

Build a production-ready, single-binary security proxy that demonstrates every skill from 120 days:

```
Internet → [L4: mTLS + SNI + Rate Limit + Circuit Breaker]
         → [L7: HTTP/2 + Bot Detection + DDoS Mitigation]
         → [Security: WAF + API Auth + Schema Validation]
         → [Observability: slog + expvar + pprof + Traces]
         → [Upstream: Reverse Proxy]
```

**Requirements checklist:**

**Transport & L4:**
- [ ] TLS termination (TLS 1.2 minimum, TLS 1.3 preferred)
- [ ] mTLS support (client cert required for internal services)
- [ ] SNI routing (different backends by hostname)
- [ ] Cert hot-reload (SIGHUP, no dropped connections)
- [ ] L4 connection rate limiting (per-IP, global)
- [ ] Circuit breaker per backend

**Bot Detection:**
- [ ] HTTP header consistency fingerprinting
- [ ] JA3-like TLS fingerprint (via `tls.ClientHelloInfo`)
- [ ] Session behavioral scoring (velocity, path diversity)
- [ ] IP ASN classification (datacenter vs residential)
- [ ] PoW challenge for suspected bots
- [ ] Honeypot endpoint

**DDoS Mitigation:**
- [ ] Token bucket rate limiting (per-IP, per-/24 subnet)
- [ ] Slowloris detection (byte rate per connection)
- [ ] SYN flood defense (per-IP half-open limit)
- [ ] Traffic anomaly detection (z-score, EMA)
- [ ] IP blocklist with auto-expiry
- [ ] Tiered response: challenge → block

**WAF:**
- [ ] Multi-layer normalization (URL-decode → HTML-unescape → lowercase)
- [ ] SQLi detection (scoring)
- [ ] XSS detection (scoring)
- [ ] Path traversal detection (immediate block)
- [ ] Command injection detection (immediate block)
- [ ] Body inspection (JSON + form, up to 1MB)
- [ ] CIDR allow-list bypass

**API Security:**
- [ ] ECDSA-signed token auth
- [ ] HMAC request signing validation
- [ ] Per-API-key rate limiting (token bucket)
- [ ] JSON schema validation per route
- [ ] Input validation (struct tags)
- [ ] Response PII scrubber
- [ ] CORS (strict origin allowlist)
- [ ] Security headers (HSTS, CSP, X-Frame-Options)

**Observability & Operations:**
- [ ] Structured logging (`log/slog`, JSON, correlation IDs)
- [ ] `expvar` metrics (requests, errors, latency percentiles, blocked IPs)
- [ ] pprof endpoints (`/debug/pprof`)
- [ ] Manual trace context (W3C `traceparent`)
- [ ] Health checks (`/healthz`, `/readyz`, per-subsystem)
- [ ] 3-stage graceful shutdown
- [ ] Signal handlers (SIGHUP: config reload, SIGUSR1: log rotate, SIGUSR2: goroutine dump)
- [ ] Zero-downtime restart (FD inheritance)
- [ ] Chaos middleware (configurable fault injection)

**Quality:**
- [ ] `go test -race ./...` — zero races
- [ ] Security-critical packages: >80% test coverage
- [ ] Fuzz tests: WAF normalization, JA3 parser, schema validator
- [ ] Load test: 500 concurrent clients, p99 overhead < 10ms

**Deliverable:** Single binary, ~4000–6000 lines of Go, zero external dependencies.

---

## Skills Matrix

| Skill | Phase | Week | Project |
|-------|-------|------|---------|
| Go syntax, types, functions | 1 | 1 | Config Server |
| Error handling, interfaces | 1 | 1–2 | Config Server, URL Checker |
| Goroutines, channels, sync | 1 | 2 | URL Checker, Metrics Server |
| Context, I/O, JSON | 1 | 3 | Request Logger |
| Generics, profiling, testing | 1 | 4 | Metrics Server |
| TCP/UDP networking | 2 | 5 | TCP Proxy |
| TLS, mTLS, certificates | 2 | 6 | mTLS Gateway |
| HMAC, AES, ECDSA | 2 | 6 | mTLS Gateway |
| L4 load balancing | 2 | 7 | L4 Load Balancer |
| Circuit breaker, health checks | 2 | 7 | L4 Load Balancer |
| HTTP/1.1, HTTP/2, WebSocket | 3 | 9 | HTTP Reverse Proxy |
| HTTP middleware chains | 3 | 9 | HTTP Reverse Proxy |
| WAF rule engine | 3 | 10 | WAF Middleware |
| SQLi/XSS/traversal detection | 3 | 10 | WAF Middleware |
| API auth (ECDSA, HMAC) | 3 | 11 | API Security Gateway |
| Rate limiting (token bucket) | 3 | 11 | API Security Gateway |
| Input validation, schemas | 3 | 11 | API Security Gateway |
| DDoS rate limiting algorithms | 4 | 13 | DDoS Mitigation Layer |
| Slowloris, SYN flood defense | 4 | 13 | DDoS Mitigation Layer |
| PoW challenges | 4 | 13–14 | DDoS + Bot Engine |
| Bot header fingerprinting | 4 | 14 | Bot Detection Engine |
| JA3 TLS fingerprinting | 4 | 14 | Bot Detection Engine |
| Session behavioral scoring | 4 | 14 | Bot Detection Engine |
| Health checks, shutdown | 4 | 15 | Production Gateway |
| Signal handling, hot-reload | 4 | 15 | Production Gateway |
| W3C trace context | 4 | 15 | Production Gateway |
| Zero-downtime restart | 4 | 15 | Production Gateway |
| Fuzz testing | 4 | 16 | Capstone |
| End-to-end attack simulation | 4 | 16 | Capstone |

---

## Engineering Decision Framework

At the end of this plan, you'll be able to answer these questions in design reviews and incidents with code-level confidence:

| Question | You'll know |
|----------|-------------|
| "Should we build or buy a WAF?" | The algorithms, their limits, and the maintenance cost |
| "Which DDoS protection settings do we actually need?" | What each layer does, what it costs in latency, and what you'd need to replicate it |
| "Why is our bot rate 40%?" | How to measure it, what signals to add, what to tune |
| "Can we do mTLS between services without a service mesh?" | Yes — and you've built it |
| "Why does our WAF have false positives on this endpoint?" | Exactly which rule matched and how to write the exception |
| "How do we do zero-downtime deployment for the security proxy?" | FD inheritance and the 3-stage drain pattern |
| "Our rate limiter is causing latency spikes" | Token bucket vs sliding window tradeoffs and how to profile it |

---

## Recommended Reading by Phase

> Read alongside the daily exercises — not before. Open the book when the exercise raises a question you can't answer from godoc alone.

---

### Phase 1: Go Fundamentals (Days 1–30)

**Primary — read in this order:**

| # | Book | Author | Why |
|---|------|--------|-----|
| 1 | **The Go Programming Language** | Donovan & Kernighan (Addison-Wesley, 2015) | The definitive Go reference. Chapters 1–5 (syntax, types, functions, interfaces) map directly to Weeks 1–2. Chapter 8 (goroutines/channels) maps to Week 2. Chapter 11 (testing) maps to Week 2. Read a chapter per week alongside the exercises. |
| 2 | **Learning Go** (2nd ed.) | Jon Bodner (O'Reilly, 2024) | Modern companion to TGoPL. Covers generics (Week 4), idiomatic error handling (Day 7), and `context` (Day 15) with more contemporary examples. Read after TGoPL for a second perspective on the same topics. |
| 3 | **Concurrency in Go** | Katherine Cox-Buday (O'Reilly, 2017) | Read during Week 2 (Days 12–14) and Week 4. Deep treatment of goroutine patterns, channel pipelines, `sync` primitives, and the context package. Directly backs the exercises on worker pools, fan-in/out, and cancellation. |

**Supplementary:**

| Book | Author | When to Read | Why |
|------|--------|--------------|-----|
| **100 Go Mistakes and How to Avoid Them** | Teiva Harsanyi (Manning, 2022) | After Week 4 | Catalog of real production bugs (slices, goroutine leaks, error wrapping, benchmarks). Read as a code review checklist for your Phase 1 projects. |
| **Go in Practice** | Butcher & Farina (Manning, 2016) | During Week 3–4 | Pattern-focused. Covers I/O composition, JSON streaming, and middleware patterns with concrete examples. |

**Free references (keep open daily):**
- Go specification: `go.dev/ref/spec` — the authoritative source for language behavior
- Standard library docs: `pkg.go.dev` — read the package source, not just the API docs
- Go blog: `go.dev/blog` — official articles on generics, memory model, profiling

---

### Phase 2: L4 Networking + TLS (Days 31–60)

**Primary — read in this order:**

| # | Book | Author | Why |
|---|------|--------|-----|
| 1 | **Network Programming with Go** | Jan Newmarch (Apress, 2020) | The only Go-specific networking book. Covers `net.Dial`, `net.Listen`, TCP/UDP servers, TLS in Go, and HTTP from first principles. Read Chapter 1–4 during Week 5, Chapter 8–9 (TLS) during Week 6. |
| 2 | **TCP/IP Illustrated, Volume 1: The Protocols** (2nd ed.) | W. Richard Stevens & Gary Wright (Addison-Wesley, 2011) | The definitive TCP/IP reference. Language-agnostic. Read Chapter 12–18 (TCP) during Week 5. Chapter 2–4 (IP, ARP) for context. You won't read this cover-to-cover — use it as a reference when an exercise raises a "but why does TCP do this?" question. |
| 3 | **Bulletproof TLS and PKI** (2nd ed.) | Ivan Ristić (Feisty Duck, 2022) | The definitive TLS and certificate reference. Author runs SSL Labs (ssllabs.com). Read Part 1 (TLS protocol) during Week 6 (Days 38–41). Part 2 (PKI, certificates) during Days 40–41. Essential background for the mTLS gateway project. Updated to cover TLS 1.3. |

**Supplementary:**

| Book | Author | When to Read | Why |
|------|--------|--------------|-----|
| **Cryptography Engineering** | Ferguson, Schneier & Kohno (Wiley, 2010) | During Week 6 (Days 42–44) | Deep dive on AES, HMAC, ECDSA, and random number generation from first principles. Explains *why* specific constructions (GCM vs CBC, ECDSA vs RSA) are chosen. Read chapters on symmetric encryption, MACs, and public-key crypto. |
| **Computer Networks: A Top-Down Approach** (8th ed.) | Kurose & Ross (Pearson, 2021) | Before Week 5 if networking feels shaky | The standard university networking textbook. Chapters 3 (transport layer, TCP) and 4 (network layer) fill gaps if OSI/TCP feels abstract. More accessible than Stevens. |
| **Unix Network Programming, Vol. 1** (3rd ed.) | W. Richard Stevens (Prentice Hall, 2003) | During Week 7 (connection pooling, SO_REUSEPORT) | The reference for socket programming. Chapter 7 (socket options), Chapter 16 (non-blocking I/O), and Appendix B (TCP state machine) are most relevant. |

**Free references:**
- RFC 793 (TCP), RFC 9293 (updated TCP) — the actual spec for what TCP does
- RFC 8446 (TLS 1.3) — read Section 2 (overview) and Section 4 (handshake) alongside Day 38
- RFC 5280 (X.509 certificates) — read Section 4 alongside Day 40

---

### Phase 3: L7 Networking + WAF + API Security (Days 61–90)

**Primary — read in this order:**

| # | Book | Author | Why |
|---|------|--------|-----|
| 1 | **HTTP: The Definitive Guide** | Gourley, Totty, Sayer et al. (O'Reilly, 2002) | Still the most comprehensive HTTP reference despite its age (HTTP/1.1 is unchanged). Read Chapter 1–4 (HTTP basics, connections, messages) before Week 9. Chapter 11 (proxies) before Day 64. Chapter 14 (security) before Week 11. Chapter 7 (caching) and 8 (integration points) for context. |
| 2 | **The Web Application Hacker's Handbook** (2nd ed.) | Stuttard & Pinto (Wiley, 2011) | Attack techniques that map directly to WAF rules. Read Chapter 9 (SQLi) before Day 69. Chapter 12 (XSS) before Day 70. Chapter 10 (code injection) before Day 71. Chapter 4 (mapping the application) for the attacker's perspective on WAF bypass. This is the offensive side of what you're building defensively. |
| 3 | **API Security in Action** | Neil Madden (Manning, 2020) | Deep treatment of API authentication and authorization. Chapter 3 (token-based auth) maps to Day 75. Chapter 4 (OAuth2) maps to Phase 3 week 11. Chapter 8 (rate limiting) maps to Day 76. Chapter 11 (mTLS) ties back to Phase 2. The best single book covering Week 11's topics. |

**Supplementary:**

| Book | Author | When to Read | Why |
|------|--------|--------------|-----|
| **The Tangled Web** | Michal Zalewski (No Starch Press, 2011) | During Week 9 (Days 67, 81) | Browser security model from first principles. Same-origin policy, CORS, XSS, and clickjacking explained from the attacker's perspective. Read before Day 67 (security headers) and Day 81 (CORS). Explains *why* CORS misconfigs are critical, not just that they are. |
| **OWASP Testing Guide** (v4.2) | OWASP Foundation (free PDF) | During Week 10 (Days 68–74) | Free. The reference for web application vulnerability testing. Use as a checklist when writing WAF rules — for each category you implement, the OWASP guide shows the full attack taxonomy and bypass techniques. Download at owasp.org. |
| **Web Application Security** | Andrew Hoffman (O'Reilly, 2020) | During Week 10 | More modern than WAHH. Good coverage of encoding bypasses and modern XSS variants. Chapters on recon and offense inform how you write defensive rules. |

**Free references:**
- RFC 9110 (HTTP Semantics) — the 2022 consolidation of all HTTP/1.1 RFCs
- RFC 9113 (HTTP/2) — read Section 5 (streams) alongside Day 65
- RFC 6455 (WebSocket) — read alongside Day 66
- OWASP Cheat Sheet Series (cheatsheetseries.owasp.org) — one cheat sheet per attack type

---

### Phase 4: DDoS + Bot Detection + Production (Days 91–120)

**Primary — read in this order:**

| # | Book | Author | Why |
|---|------|--------|-----|
| 1 | **Security Engineering** (3rd ed.) | Ross Anderson (Wiley, 2020) | **Free online** at cl.cam.ac.uk/~rja14/book.html. Comprehensive security engineering textbook by a Cambridge professor. Chapter 21 (network attack and defense) maps to Week 13 (DDoS). Chapter 22 (telecom) covers amplification attacks. Chapter 3 (psychology of security failures) is essential EM reading. Dip in by chapter, don't read linearly. |
| 2 | **Site Reliability Engineering** | Beyer, Jones, Petoff & Murphy (Google/O'Reilly) | **Free online** at sre.google/books. The production hardening bible. Chapter 8 (release engineering / zero-downtime), Chapter 13 (emergency response), Chapter 21 (managing load), and Chapter 29 (dealing with interrupts) map directly to Week 15. Read during Days 106–112. |
| 3 | **The Practice of Network Security Monitoring** | Richard Bejtlich (No Starch Press, 2013) | Traffic analysis and anomaly detection from a practitioner. Chapter 1–3 (NSM concepts) and Chapter 8–9 (analysis) inform the DDoS anomaly detection work in Week 13 (Days 95–96). Bejtlich's mental model of "collect, detect, respond" maps well to the `ThreatIntel` store design. |

**Supplementary:**

| Book | Author | When to Read | Why |
|------|--------|--------------|-----|
| **Silence on the Wire** | Michal Zalewski (No Starch Press, 2005) | During Week 14 (Bot Detection) | Passive fingerprinting and traffic analysis. The chapter on TCP/IP stack fingerprinting is the conceptual foundation for JA3 TLS fingerprinting (Days 101–103). Zalewski's approach to behavioral inference informs the session scoring design. |
| **The Art of Intrusion** | Kevin Mitnick (Wiley, 2005) | During Week 14 | Case studies of real attacks. Read for attacker mindset — not for technical depth but to understand how determined attackers think around bot detection and rate limiting controls. Useful context before Day 113 (attack simulation). |
| **Accelerate: The Science of Lean Software and DevOps** | Forsgren, Humble & Kim (IT Revolution, 2018) | During Week 16 | The data behind deployment frequency, MTTR, and change failure rate. Connects Phase 4's production hardening work (graceful shutdown, zero-downtime deployment) to organizational outcomes you care about as an EM. |
| **The Phoenix Project** | Kim, Behr & Spafford (IT Revolution, 2013) | Anytime | A novel about DevOps transformation. Not technical — but it provides narrative context for why operational characteristics (observability, graceful shutdown, fast deployments) matter at the organizational level. |

**Free references:**
- Google SRE Workbook (sre.google/workbook) — companion to the SRE book, more prescriptive
- NIST SP 800-38D — AES-GCM specification (relevant to Day 43 and the audit log)
- Cloudflare blog (blog.cloudflare.com) — articles on DDoS mitigation, JA3, bot detection written by people who operate these systems at scale
- IETF RFC 8941 (Structured Field Values for HTTP) — context for rate limit headers

---

### Cross-Phase Reading (All 120 Days)

These span multiple phases and are worth reading alongside the entire plan:

| Book | Author | Read When | Why |
|------|--------|-----------|-----|
| **Designing Data-Intensive Applications** | Martin Kleppmann (O'Reilly, 2017) | Chapters as relevant | Chapter 1 (reliability, scalability) for all phases. Chapter 8 (distributed systems) for rate limiting and state management in Phase 3–4. The mental model of "what can go wrong" improves every design decision in this plan. |
| **The Pragmatic Programmer** (20th anniversary ed.) | Hunt & Thomas (Addison-Wesley, 2019) | One tip per week | Craft principles (DRY, orthogonality, design by contract) that apply throughout. Read one section per week as a metacognitive complement to the technical exercises. |
| **A Philosophy of Software Design** | John Ousterhout (Yaknyam Press, 2018) | During Phase 3–4 integration weeks | Deep modules, shallow interfaces, and designing systems that are easy to reason about. Directly applicable when integrating WAF + DDoS + Bot + API into a coherent middleware chain. |

---

### Reading Schedule (Integrated with Daily Workflow)

The daily 1-hour budget is for coding. Read books separately — 20–30 minutes before bed or during commute.

```
Phase 1 (Days 1–30)
  Week 1–2    The Go Programming Language, Ch 1–5, 8
  Week 3–4    Learning Go (modern patterns) + Concurrency in Go

Phase 2 (Days 31–60)
  Week 5      Network Programming with Go, Ch 1–4
              TCP/IP Illustrated Vol. 1 (TCP chapters, as reference)
  Week 6      Bulletproof TLS and PKI, Part 1–2
              Cryptography Engineering (AES, HMAC, ECDSA chapters)
  Week 7–8    Unix Network Programming (socket options, as reference)

Phase 3 (Days 61–90)
  Week 9      HTTP: The Definitive Guide, Ch 1–4, 11
  Week 10     The Web Application Hacker's Handbook (SQLi, XSS, injection chapters)
              OWASP Testing Guide (as checklist)
  Week 11–12  API Security in Action, Ch 3–4, 8, 11
              The Tangled Web (CORS, XSS model)

Phase 4 (Days 91–120)
  Week 13     Security Engineering, Ch 21–22
              Practice of Network Security Monitoring, Ch 1–3, 8–9
  Week 14     Silence on the Wire (fingerprinting chapters)
  Week 15     Site Reliability Engineering, Ch 8, 13, 21, 29
  Week 16     Accelerate (production metrics and org outcomes)
```

---

## Key Standard Library Packages

```
net                  - TCP/UDP/Unix sockets, Dialer, Listener, IPNet
net/http             - HTTP server/client, middleware, handlers, httptest
net/http/httputil    - ReverseProxy, DumpRequest
crypto/tls           - TLS 1.2/1.3, ClientHelloInfo, Config, mTLS
crypto/x509          - Certificate parsing, generation, CA operations
crypto/ecdsa         - ECDSA signing/verification (tokens, mTLS)
crypto/rsa           - RSA key generation
crypto/elliptic      - P-256, P-384 curves
crypto/aes           - AES-256 encryption
crypto/cipher        - GCM, CBC modes
crypto/hmac          - HMAC-SHA256 (request signing, PoW, webhook)
crypto/sha256        - Hashing, PoW challenges, JA3
crypto/md5           - JA3 fingerprint hash (not for security use)
crypto/rand          - Cryptographically secure random (nonces, tokens)
crypto/subtle        - Constant-time comparison (auth paths)
sync                 - Mutex, RWMutex, WaitGroup, Once, Pool, Map
sync/atomic          - Lock-free counters (connection counts, flags)
context              - Cancellation, timeouts, value propagation
encoding/binary      - Packet/protocol parsing (big/little endian)
encoding/json        - JSON encode/decode, streaming, schema validation
encoding/hex         - Hex encoding for fingerprints, tokens
regexp               - WAF pattern matching (pre-compiled)
log/slog             - Structured logging, JSON handler (Go 1.21+)
net/http/pprof       - CPU/heap/goroutine profiling endpoints
expvar               - Runtime metrics exposure (/debug/vars)
testing              - Unit tests, table-driven tests
testing/fuzz         - Fuzz testing for parsers and validators
os/signal            - Signal handling (SIGTERM, SIGHUP, SIGUSR*)
runtime              - GC tuning, goroutine count, stack dumps
html                 - HTML entity unescaping (WAF normalization)
bufio                - Buffered I/O for protocol parsing
io                   - Reader/Writer composition, TeeReader, Pipe
time                 - Timeouts, rate limiting windows, TTL, metrics
container/heap       - Priority queue (task scheduler)
```

---

## Daily Workflow (~1 hour)

```
:00 – :10  Read       Review the topic. Read the stdlib godoc page.
            Ask: "When would I need this in a real service?"
:10 – :50  Code       Write the exercise. Start from scratch.
            No copy-paste. Type it. Understand each line.
:50 – :55  Test       go test -race ./...
            Did it pass? Did -race find anything?
:55 – :60  Reflect    What would you change?
            What would break under 10x load?
            What would you add before shipping?
```

## Weekly Rhythm

```
Mon–Thu   Daily topics + exercises (build instinct through repetition)
Fri       Mini project (integrate the week's concepts, ~90 min ok)
Sat       Code review: re-read your own mini project code
          Ask: what would a senior engineer change?
Sun       Rest. Or read one Go proposal/blog post (optional).
```

## Mindset for a Software Engineer

You're not learning Go just to add a language to your CV. You're learning it so you can:
- **Build network and security services** that hold up in production, not just in demos
- **Review designs and PRs** and know whether a proposed WAF or bot detection approach is sound
- **Use managed products well** by knowing what their algorithms actually do
- **Debug the hard problems** (stuck connections, TLS failures, false positives) instead of escalating them
- **Hold yourself to production standards** (test coverage, graceful shutdown, observability) from experience, not checklists

The exercises are deliberately narrow (1 hour each) and deep. Breadth comes from the mini projects.
