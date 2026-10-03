# Go SME Gap Analysis & Real-World Project Roadmap

Supplementary to `week1.md`–`week17.md`. Does not modify those files. Based on a
full read-through of all 17 weeks (120 days), reviewed against the bar of a
practicing Go SME across three domains: **(1) backend/distributed systems,
(2) infra/DevOps tooling, (3) high-performance/low-latency systems**.

## What this curriculum actually is

Despite the folder name, this is not a broad "learn Go in 90 days" plan — it's
a single continuous capstone: **build a from-scratch cloud-edge security proxy**
(WAF → DDoS mitigation → bot detection → production hardening → integration),
ending in a ~5,000-line single static binary. Weeks 1–6 lay Go foundations
*specifically in service of* that proxy (config/storage, concurrency, context/IO,
generics/profiling, TCP, TLS/crypto). Weeks 7–17 build the proxy itself, layer
by layer, with weeks 12–17 forming one unbroken integration arc.

**This is genuinely rare depth** in network security engineering — PKI/mTLS,
JA3 fingerprinting, 5 rate-limiting algorithms from scratch, SYN-flood/Slowloris
defense, WAF normalization bypass handling, FD-inheritance zero-downtime
restarts, W3C trace-context parsing. Most Go courses never go near any of this.

**The cost of that depth is breadth.** Because every day serves the one proxy,
entire categories of "everyday Go SME" work never appear at all, in 120 days:

| Never appears, anywhere in weeks 1–17 |
|---|
| `database/sql`, `pgx`, transactions, migrations — **zero persistence** |
| gRPC / protobuf |
| Message queues (Kafka/NATS/SQS), event-driven/async patterns |
| Kubernetes as *producer* — no `client-go`, `controller-runtime`, CRDs, operators (only as *consumer*: liveness/readiness semantics) |
| CLI tooling (`cobra`/`viper`) — despite this being a named target niche |
| Prometheus (`client_golang`) or OpenTelemetry SDKs — metrics/tracing are hand-rolled (`expvar`, custom ring-buffer spans) throughout |
| GC/memory tuning (`GOGC`, `GOMEMLIMIT`), escape analysis, `benchstat` |
| `go test`-based unit tests until **Week 12** — days 1–81 are almost entirely `func main()` demos, including a dedicated `-race`-auditing day with no `_test.go` file |
| Password hashing (bcrypt/argon2) — a striking gap given an entire "API Security" week |

## Topic map (quick reference)

| Week | Theme | Real-world shape |
|---|---|---|
| 1 | Types, structs, interfaces, errors | Toy-to-real transition; config/storage building blocks |
| 2 | Control flow, packages, testing, goroutines/channels | Strong (race detector, benchmarks, worker pools) |
| 3 | Context, I/O, JSON, regex, time | Strong (OWASP-aware, PCI/GDPR framing) |
| 4 | Generics, sync.Pool, pprof, graceful shutdown | Most senior-feeling early week; metrics self-rolled, not Prometheus |
| 5 | Raw TCP/UDP, DNS, connection pooling | Uncommonly deep L4 work |
| 6 | TLS, PKI, mTLS, HMAC, AES-GCM, ECDSA tokens | Strongest security content in weeks 1–6 |
| 7 | L4 load balancer (RR/weighted/least-conn), circuit breaker | Real, but zero `_test.go` |
| 8 | Binary protocols, cert rotation, race audit, load test | Best debugging-methodology week; docs promise SIGHUP that code never wires |
| 9 | HTTP/1.1 internals, timeouts, reverse proxy, WebSocket | Strong; no DB, no gRPC despite covering HTTP/2 |
| 10 | WAF (OWASP Top 10 detection) | Best security-engineering week; regex-chain has no perf story |
| 11 | API auth (API key/HMAC/ECDSA), validation, PII, CORS | Hand-rolled crypto/JWT with no "use a real library" caveat; no bcrypt |
| 12 | Integration, slog, SIGHUP reload, first real tests + fuzzing | Testing/observability finally treated as first-class |
| 13 | DDoS detection (5 rate-limit algos, SYN/Slowloris, PoW) | Production-level within its narrow domain |
| 14 | Bot detection (JA3, ASN classification, PoW) | Vendor-informed, realistic; JA3 is an acknowledged approximation |
| 15 | SRE hardening (health checks, FD-inherit restart, chaos middleware) | Most broadly transferable week of all 17 |
| 16 | Attack simulation, profiling, coverage, ADRs, fuzzing | Most senior-engineer-shaped week (ADR discipline, adversarial self-review) |
| 17 | Capstone assembly (single `main.go`, build-vs-buy analysis) | One large integration, not a portfolio |

## SME gap analysis by domain

### Domain 1 — Backend / distributed systems (largest gap)

- **No persistence anywhere.** Every store (API keys, abuse counters, rate-limit
  buckets, ThreatIntel scores) is an in-memory map. No `database/sql`, no
  connection pooling, no transactions, no migrations, no N+1/exhaustion failure
  modes — the incidents that actually dominate backend on-call.
- **No gRPC/protobuf**, despite HTTP/2 (its transport) being covered in depth in
  Week 9, and despite gRPC being explicitly named in target-domain scope.
- **No message queues / event-driven architecture** — no idempotency-key
  patterns, no outbox pattern, no at-least-once delivery reasoning.
- **No real observability stack.** `expvar` + hand-rolled ring-buffer tracing is
  taught as if it *replaces* Prometheus/OpenTelemetry rather than as a stepping
  stone toward them.
- **Testing culture arrives 81 days late.** Weeks 1–11 build production-shaped
  code with zero `_test.go` files; Week 8 dedicates a full day to `-race`
  auditing with nothing to run `-race` against except `func main()`.
- ADR-003 (Week 16) explicitly *chooses* in-process `sync.Map` over Redis
  "because signals are ephemeral" — a deliberate, acknowledged gap that means
  the learner never touches a real distributed cache/coordination client.

### Domain 2 — Infra / DevOps tooling

- **Kubernetes is consumer-only knowledge.** Liveness/readiness/startup probes,
  SIGTERM drain semantics, and rollout behavior are referenced constantly as
  analogies — but there is no `client-go`, no `controller-runtime`, no CRDs,
  no informers/workqueues, no operator pattern anywhere in 120 days, despite
  this being one of three explicitly named target domains.
- **No CLI tooling** (`cobra`/`viper`/`pflag`). Every tool is a bare `main.go`
  reading JSON files and OS signals, never a `kubectl`/`docker`-style CLI.
- **No containerization** — no Dockerfile, multi-stage builds, or distroless
  images, even though "ships as one static binary" is a running theme that
  begs for a container story.
- **No IaC/GitOps** — CI is described in prose ("require `-race` in CI") but
  never expressed as an actual pipeline file anywhere in the 17 weeks.

### Domain 3 — High-performance / low-latency systems

- Foundations exist (`sync.Pool` intro, `pprof` registration, atomics, ring
  buffers) but **profiling is named far more than it's exercised** — `pprof`
  is "served automatically" in checklists without a walked CPU/heap profile
  session in most weeks that mention it.
- **No GC/memory tuning** (`GOGC`, `GOMEMLIMIT`, `debug.SetGCPercent`), no
  escape-analysis workflow (`-gcflags=-m`), no `benchstat`-driven optimization
  loop — benchmarking stays at raw `ns/op`/`allocs/op` numbers.
- **No lock-free/cache-conscious data-structure design** beyond basic CAS loops
  in the token/leaky-bucket limiters — no false-sharing, cache-line padding,
  or `GOMAXPROCS`/scheduler-internals discussion.
- The WAF (Week 10) is a real anti-pattern for a low-latency system: 5+
  sequential regexes per request, no compiled multi-pattern matcher
  (Aho-Corasick), no before/after benchmark of WAF overhead despite a load-test
  tool being built specifically to measure it.
- Everything stays at the `net/http`/stdlib-socket abstraction level — no
  `SO_REUSEPORT`, `TCP_NODELAY`, kernel buffer tuning, `io_uring`, or discussion
  of the microsecond-to-nanosecond budgets that trading/gaming SMEs live in.
  Week 16's own perf work measures the proxy at 100–500µs — a full order of
  magnitude looser than that floor.

## Real-world bridge projects

Consolidated from all three per-range reviews, de-duplicated, and sequenced.
Each closes a specific gap above using code the curriculum already built —
none of these ask you to start from zero.

### Track A — Backend / distributed systems

**A1. Persist the capstone's auth store to Postgres.**
Replace `MemAPIKeyStore` (Week 11) and the abuse-lockout counters with
`pgx` + connection pooling + `golang-migrate` migrations; add bcrypt password
hashing (closes the glaring gap in Week 11's `CreateUserRequest`); wrap
lockout-counter writes in real transactions. *Mirrors:* Auth0/Okta-style
credential store. *Forces:* the DB gap, transaction semantics, durability
across restarts.

**A2. gRPC-ify the stack with real tracing.**
Build 3 gRPC services (orders/inventory/payments or auth/waf-admin/threat-intel)
with protobuf codegen, unary + streaming RPCs, and gRPC interceptors (the RPC
analog of the HTTP middleware chain already mastered in Week 9). Export spans
via the real OpenTelemetry SDK to Jaeger/Tempo — not the hand-rolled W3C
trace-context ring buffer from Week 15. *Forces:* the total gRPC gap and real
distributed tracing.

**A3. Distributed rate limiter — close ADR-003 for real.**
Rebuild the Week 13 token bucket as a Redis-backed limiter using atomic Lua
scripts (`EVAL`) and client-side consistent hashing across a small Redis
cluster; back it with real `client_golang` Prometheus metrics and a Grafana
dashboard instead of `expvar`. *Mirrors:* the exact gap the curriculum's own
ADR-003 explicitly deferred. *Forces:* first real external-datastore
integration plus the metrics stack a real SRE org expects.

**A4. Test-and-fuzz retrofit.**
Convert the two best fuzz targets already built — the Week 9 hand-rolled
HTTP/1.1 parser and the Week 8 binary KV protocol parser — into real
`_test.go` table-driven tests plus `go test -fuzz` corpora, and wire a GitHub
Actions pipeline (`vet` + `-race` + fuzz + bench) against the Week 8/12
gateways. *Forces:* closing the 81-day testing-culture gap.

### Track B — Infra / DevOps tooling

**B1. Kubernetes operator for the proxy's live config.**
Build a `RateLimitPolicy` or `BackendPool` CRD with `client-go`/
`controller-runtime`: a reconcile loop that replaces the SIGHUP/polling-file
reload (Weeks 8, 12, 15) with real informers + workqueues, feeding live Pod
endpoints into the Week 7 load balancer's backend list. Add leader election
for HA. *Mirrors:* cert-manager, ingress-nginx, AWS Load Balancer Controller.
*Forces:* the single largest gap in the entire curriculum — Kubernetes as
producer, not just a source of analogies.

**B2. `edgectl` — a Cobra/Viper control-plane CLI.**
Build a `kubectl`-style CLI (`edgectl block <ip>`, `edgectl waf tail`,
`edgectl threat top --format=json`) talking to a small gRPC/REST admin API
exposed by the capstone proxy, with shell completion and a `goreleaser`
release pipeline. *Forces:* the CLI-tooling gap despite it being a named
target niche.

**B3. Containerize and ship.**
Multi-stage Dockerfile producing a distroless/scratch image of the Week 17
capstone binary, plus a GitHub Actions pipeline that builds, scans (`trivy`),
and pushes it, and a minimal Helm chart wiring the liveness/readiness probes
Week 15 already built into real Kubernetes manifests. *Forces:* the
container/IaC gap; cheap to do, high payoff for "shippable" credibility.

### Track C — High-performance / low-latency systems

**C1. Make the WAF fast, and prove it.**
Add `testing.B` benchmarks for the Week 10 normalization pipeline and rule
matcher; introduce `sync.Pool` for the per-request buffers it currently
allocates; replace the sequential regex chain with a compiled Aho-Corasick
multi-pattern scanner for single-pass detection; drive the before/after
comparison with `pprof` + `benchstat`. *Forces:* the GC/allocation/profiling
literacy the curriculum only gestures at, using code you already own.

**C2. Nanosecond-budget matching engine.**
A separate, deliberately harder project: an in-memory limit-order-book
matching engine (price-time priority) fed by a lock-free SPSC ring buffer,
with `GOMAXPROCS` pinning, `GOGC`/`GOMEMLIMIT` tuning, zero-allocation hot
path validated by `testing.B` + `benchstat` + escape analysis
(`-gcflags=-m`), and p99.9 tail-latency tracking under load. *Mirrors:* a
trading-system tick handler. *Forces:* the real nanosecond-to-low-microsecond
discipline that Week 16's 100–500µs HTTP profiling doesn't reach.

**C3. Authoritative UDP game server.**
Extend the Week 13 anti-amplification UDP work into a real bidirectional,
tick-based (30–60Hz) protocol with client-side prediction, server
reconciliation, and snapshot interpolation, enforcing a strict per-tick frame
budget measured via `pprof`. Optionally orchestrate a fleet with Agones on
Kubernetes, tying Track B and Track C together. *Forces:* real low-latency
UDP protocol design under hard per-tick budgets.

## Suggested sequencing (a "Week 18+" continuation)

1. **B3** (containerize) — cheapest, immediately makes the existing capstone
   demoable and deployable.
2. **A4** (test/fuzz retrofit) — pays down risk before adding more surface area.
3. **A1** (Postgres persistence + bcrypt) — closes the most embarrassing gap
   (no DB, no password hashing) with the smallest new-concept load.
4. **B1** (Kubernetes operator) — the single highest-leverage project for the
   "infra tooling" domain; reuses the config-reload logic already built.
5. **A3** (distributed rate limiter) — directly resolves the curriculum's own
   ADR-003, and introduces Redis + real Prometheus in one project.
6. **A2** (gRPC + OpenTelemetry) — closes the RPC and tracing gaps together.
7. **B2** (`edgectl` CLI) — lower effort, high resume/portfolio value once
   there's a real admin API (from A2/B1) to point it at.
8. **C1** (WAF performance pass) — cheap, uses existing code, teaches the
   profiling discipline needed before attempting C2.
9. **C2 or C3** (low-latency capstone) — pick based on target sub-domain
   (trading vs. gaming); this is the deepest, most specialized addition and
   should come last.

## Bottom line

Weeks 1–17 make you genuinely strong in **network/security-layer Go** — TLS,
raw TCP/UDP, WAF/DDoS/bot-detection engineering, and SRE-grade hardening —
areas most Go curricula never touch. To be a *complete* SME across backend,
infra tooling, and high-performance systems (not just edge security), the
nine projects above are the shortest path: they reuse the capstone's own code
and directly target the exact things it never had a reason to build.
