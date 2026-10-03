# Week 17: Capstone — Cloud-Edge Security Proxy (Day 120)

> **Phase 4 Final:** Day 120 is assembly day. Every component was built and tested in Weeks 1–16. Today you wire them together into a single production-ready binary: the Cloud-Edge Security Proxy.

**Why this day matters:** This is the deliverable that demonstrates 120 days of work. ~4000–6000 lines of Go, zero external dependencies, every component built from scratch using the standard library. This is the architecture of Cloudflare's edge, AWS WAF + ALB, and Fastly's compute platform — in one binary.

---

> **Background in this wiki:** HTTPS guide [Part 9 (the server side, built in Go)](../../v2-https/real-life-guide-v1.md#part-9-build-the-lifecycle-in-go); TCP/IP guide [Part 13 (professional networking)](../../networking/tcp-ip/real-life-guide-v1.md#part-13-professional-networking-operations-debugging-and-design); Linux guide [Ch 82 (shipping the service)](../../os-linux/real-life-os-guide.md#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes); [OSI guide, Part 15](../../networking/real-life-example-osi.md#part-15-architecture-decisions-by-layer-the-senior-and-manager-view). Security in Depth: [Chapter 1 (attack paths, Go lab)](../../security/real-life-security-guide-v1.md#chapter-1-from-vulnerabilities-to-attack-paths) and [Chapter 77 (architecture review)](../../security/real-life-security-guide-v1.md#chapter-77-security-architecture-review-the-senior-engineer-playbook); and the whole-series request walkthrough, [HTTPS Chapter 25](../../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide).

## Day 120 — Cloud-Edge Security Proxy (Capstone)

### Architecture

```
                         CLOUD-EDGE SECURITY PROXY
                         ─────────────────────────
                         Single binary | ~5000 LOC
                         Zero external dependencies
                         Standard library only

Internet
   │
   │  TCP connection (port 443)
   ▼
┌─────────────────────────────────────────────────────┐
│ L4 TRANSPORT LAYER                                  │
│                                                     │
│  ProtectedListener                                  │
│    └── max 5 half-open connections per source IP    │
│    └── 3s half-open drain timeout                   │
│                                                     │
│  SlowlorisDetector                                  │
│    └── close conns < 50B/s during header phase      │
│                                                     │
│  TLS Termination (crypto/tls)                       │
│    └── TLS 1.2 min, TLS 1.3 preferred               │
│    └── mTLS: RequireAndVerifyClientCert option       │
│    └── SNI routing: hostname → backend pool         │
│    └── CertReloader: SIGHUP hot-swap, no drop       │
│    └── JA3 fingerprint via GetConfigForClient       │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│ DDOS MITIGATION LAYER                               │
│                                                     │
│  CIDRRateLimiter                                    │
│    └── token bucket per /32 (100 req/s, burst 200)  │
│    └── token bucket per /24 (500 req/s, burst 1000) │
│    └── IPv6 /48 grouping                            │
│                                                     │
│  TrafficAnalyzer                                    │
│    └── EMA baseline (α=0.1)                         │
│    └── z-score anomaly (flag > 3σ)                  │
│                                                     │
│  TieredResponder                                    │
│    └── trusted CIDR bypass                          │
│    └── score 50–80: PoW challenge (SHA256, d=16)    │
│    └── score > 80: 403 block (5min expiry)          │
│    └── score feeds into ThreatIntel                 │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│ BOT DETECTION LAYER                                 │
│                                                     │
│  BotScorer (pluggable signal system)                │
│                                                     │
│  Signal detectors:                                  │
│    ├── JA3 TLS fingerprint (curl/python/scrapy → +25│
│    ├── Header consistency (UA vs Accept/Sec-CH-UA)  │
│    ├── IP ASN/datacenter classification (+15)       │
│    └── Session behavior (velocity, path diversity,  │
│        referrer rate)                               │
│                                                     │
│  TieredBotMiddleware                                │
│    ├── score 0–20:  allow                           │
│    ├── score 20–40: +500ms artificial delay         │
│    ├── score 40–70: PoW challenge redirect          │
│    ├── score 70–85: 403 block                       │
│    └── score 85+:   shadow-ban (200, empty data)    │
│                                                     │
│  HoneypotHandler (/api/v0/internal)                 │
│    └── any hit → permanent shadow-ban               │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│ HTTP/L7 LAYER                                       │
│                                                     │
│  http.Server                                        │
│    └── HTTP/1.1 + HTTP/2 (via TLS)                  │
│    └── ReadHeaderTimeout: 5s (slowloris defense)    │
│    └── ReadTimeout: 30s                             │
│    └── WriteTimeout: 30s                            │
│    └── MaxHeaderBytes: 64KB                         │
│                                                     │
│  Middleware Chain (in order):                       │
│    1.  PanicRecoverer       (never crash)           │
│    2.  RequestID            (inject X-Request-ID)   │
│    3.  TracingMiddleware    (W3C traceparent)        │
│    4.  RequestLogger        (slog JSON)             │
│    5.  SecurityHeaders      (HSTS, CSP, etc.)       │
│    6.  CORSMiddleware       (strict origin allowlist│
│    7.  BodyLimiter          (reject > 1MB bodies)   │
│    8.  AuthMiddleware       (HMAC → ECDSA → mTLS)   │
│    9.  L7RateLimiter        (per-API-key, 10/s)     │
│    10. WAFMiddleware        (SQLi+XSS+traversal)    │
│    11. SchemaValidator      (per-route JSON schema) │
│    12. PIIScrubber          (log + response masking)│
│    13. ChaosMiddleware      (configurable, off/def) │
│    14. ReverseProxy         (to upstream backends)  │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│ CROSS-CUTTING CONCERNS                              │
│                                                     │
│  ThreatIntel store                                  │
│    └── aggregates signals from all layers           │
│    └── score decay: t½ = 10 minutes                 │
│    └── thresholds: challenge=50, block=80           │
│                                                     │
│  Observability                                      │
│    ├── /debug/vars (expvar: reqs, errors, latency)  │
│    ├── /debug/pprof (CPU, heap, goroutine profiles) │
│    ├── /debug/traces (last 1000 W3C trace spans)    │
│    ├── /debug/bots (bot scores, honeypot hits)      │
│    └── /debug/ddos (rate stats, block list)         │
│                                                     │
│  Health                                             │
│    ├── /healthz (liveness — always fast)            │
│    ├── /readyz  (readiness — checks upstreams)      │
│    └── /startup (startup probe)                     │
│                                                     │
│  Operations                                         │
│    ├── SIGTERM/SIGINT: 3-stage graceful shutdown    │
│    ├── SIGHUP: config hot-reload (validate first)   │
│    ├── SIGUSR1: log rotation + FD inheritance restart│
│    └── SIGUSR2: goroutine stack dump to /tmp        │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
                   Upstream backends
                   (any HTTP/HTTPS service)
```

---

### Complete `main.go`

```go
// cmd/proxy/main.go
// Cloud-Edge Security Proxy — capstone for 120-day Go security plan
// Standard library only. No external dependencies.

package main

import (
    "context"
    "crypto/tls"
    "expvar"
    "fmt"
    "log/slog"
    "net"
    "net/http"
    _ "net/http/pprof" // registers /debug/pprof handlers
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "time"

    // Replace these with your actual internal package paths
    "github.com/you/proxy/internal/auth"
    "github.com/you/proxy/internal/bot"
    "github.com/you/proxy/internal/chaos"
    "github.com/you/proxy/internal/ddos"
    "github.com/you/proxy/internal/health"
    "github.com/you/proxy/internal/ops"
    "github.com/you/proxy/internal/proxy"
    "github.com/you/proxy/internal/ratelimit"
    "github.com/you/proxy/internal/threat"
    "github.com/you/proxy/internal/tracing"
    "github.com/you/proxy/internal/validation"
    "github.com/you/proxy/internal/waf"
)

func main() {
    // ── 1. Logger ──────────────────────────────────────────────────────────
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    logger.Info("starting cloud-edge security proxy",
        "go_version", runtime.Version(),
        "pid", os.Getpid(),
    )

    // ── 2. Config ──────────────────────────────────────────────────────────
    cfg, err := ops.LoadConfig("config.json")
    if err != nil {
        logger.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    // ── 3. TLS ─────────────────────────────────────────────────────────────
    certReloader, err := newCertReloader(cfg.TLS.CertFile, cfg.TLS.KeyFile)
    if err != nil {
        logger.Error("failed to load TLS certificates", "error", err)
        os.Exit(1)
    }

    tlsConfig := &tls.Config{
        MinVersion:     tls.VersionTLS12,
        GetCertificate: certReloader.GetCertificate,
        NextProtos:     []string{"h2", "http/1.1"},
    }

    // mTLS support (optional, controlled by config)
    if cfg.TLS.RequireClientCert {
        tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
        tlsConfig.ClientCAs = loadCAPool(cfg.TLS.ClientCAFile)
    }

    // ── 4. Bot Detection ───────────────────────────────────────────────────
    botScorer := bot.NewBotScorer(
        &bot.HeaderConsistencyDetector{},
        &bot.IPReputationDetector{},
        bot.NewBehavioralDetector(),
    )

    ja3Detector := &bot.JA3Detector{Scorer: botScorer}
    tlsConfig.GetConfigForClient = ja3Detector.GetConfigForClient

    honeypot := bot.NewHoneypotTracker()
    powSystem := bot.NewPowSystem(cfg.PowSecret, 16)
    tieredBot := bot.NewTieredBotMiddleware(botScorer, powSystem, honeypot)

    // ── 5. DDoS ────────────────────────────────────────────────────────────
    threatIntel := threat.NewThreatIntel(10 * time.Minute)
    cidrRL := ddos.NewCIDRRateLimiter(100, 500, 10.0)
    trafficAnalyzer := ddos.NewTrafficAnalyzer()
    blockList := ddos.NewBlockList()
    tieredDDoS := &ddos.TieredResponder{
        BlockList:  blockList,
        Challenger: ddos.NewPowChallenger(cfg.PowSecret, 16),
        Whitelist:  cfg.TrustedCIDRs,
    }

    // ── 6. Tracing ─────────────────────────────────────────────────────────
    tracer := tracing.NewTracer(1000, 1.0)

    // ── 7. WAF ─────────────────────────────────────────────────────────────
    wafMiddleware, err := waf.NewMiddleware(cfg.WAFConfig, logger)
    if err != nil {
        logger.Error("failed to initialize WAF", "error", err)
        os.Exit(1)
    }

    // ── 8. Auth ────────────────────────────────────────────────────────────
    authMiddleware := auth.Chain(
        auth.HMACMiddleware(cfg.HMACSecret),
        auth.ECDSATokenMiddleware(cfg.ECDSAPublicKey),
        auth.MTLSIdentityMiddleware(), // extracts CN from client cert
    )

    // ── 9. Rate Limiting (L7) ──────────────────────────────────────────────
    l7RL := ratelimit.NewAPIKeyRateLimiter(10, 50) // 10/s, burst 50

    // ── 10. Validation ─────────────────────────────────────────────────────
    schemaValidator := validation.NewSchemaValidator(cfg.RouteSchemas)

    // ── 11. Chaos (disabled by default) ───────────────────────────────────
    chaosCfg := chaos.NewChaosConfig() // Enabled=false by default

    // ── 12. Reverse Proxy ──────────────────────────────────────────────────
    reverseProxy := proxy.NewSecureReverseProxy(cfg.Upstreams, logger)

    // ── 13. Health Checks ──────────────────────────────────────────────────
    healthRegistry := health.NewRegistry()
    startup := health.NewStartupHandler()

    for _, upstream := range cfg.Upstreams {
        u := upstream
        healthRegistry.Register(
            "upstream:"+u.Name,
            health.NewTCPCheck(u.Addr),
            2*time.Second,
            true, // required for readiness
        )
    }
    healthRegistry.Register(
        "goroutine_baseline",
        health.NewGoroutineLeakCheck(runtime.NumGoroutine(), 3.0),
        1*time.Second,
        false, // informational only
    )

    // ── 14. Middleware Chain ───────────────────────────────────────────────
    chain := buildChain(
        panicRecoverer(logger),
        requestIDMiddleware,
        tracer.Middleware,
        requestLogger(logger),
        securityHeadersMiddleware,
        corsMiddleware(cfg.AllowedOrigins),
        bodyLimiter(1<<20), // 1MB max body
        authMiddleware,
        l7RL.Middleware,
        tieredBot.Middleware,
        wafMiddleware.ServeHTTP,
        schemaValidator.Middleware,
        piiScrubberMiddleware,
        chaos.ChaosMiddleware(chaosCfg),
        reverseProxy.ServeHTTP,
    )

    // ── 15. Debug & Admin Mux ──────────────────────────────────────────────
    mux := http.NewServeMux()

    // DDoS + Bot path
    mux.Handle("/", ddosGate(cidrRL, trafficAnalyzer, tieredDDoS, threatIntel, chain))

    // Honeypot
    mux.Handle("/api/v0/internal", honeypot.HoneypotHandler(botScorer))

    // PoW challenge endpoints
    mux.HandleFunc("/challenge", powSystem.ChallengeHandler)
    mux.HandleFunc("/verify", powSystem.VerifyHandler)

    // Health
    mux.HandleFunc("/healthz", healthRegistry.LivezHandler)
    mux.HandleFunc("/readyz", healthRegistry.ReadyzHandler)
    mux.HandleFunc("/startup", startup.Handler)

    // Debug (restrict to internal networks in production)
    mux.HandleFunc("/debug/traces", tracer.TracesHandler)
    mux.HandleFunc("/debug/bots", botScorer.DebugHandler)
    mux.HandleFunc("/debug/ddos", func(w http.ResponseWriter, r *http.Request) {
        writeDDoSStats(w, cidrRL, blockList, trafficAnalyzer)
    })
    mux.HandleFunc("/debug/threat", threatIntel.DebugHandler)
    // /debug/vars and /debug/pprof/* registered by expvar + net/http/pprof imports

    // ── 16. Servers ────────────────────────────────────────────────────────

    // HTTPS server (main)
    httpServer := &http.Server{
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      30 * time.Second,
        IdleTimeout:       90 * time.Second,
        MaxHeaderBytes:    1 << 16,
    }

    // HTTP→HTTPS redirect server
    redirectServer := &http.Server{
        Addr:    cfg.HTTPAddr,
        Handler: http.HandlerFunc(httpsRedirect(cfg.HTTPSAddr)),
    }

    // ── 17. Listeners ──────────────────────────────────────────────────────
    var baseListener net.Listener
    if inherited, ok := ops.ListenerFromEnv(); ok {
        logger.Info("using inherited listener FD (zero-downtime restart)")
        baseListener = inherited
    } else {
        baseListener, err = net.Listen("tcp", cfg.HTTPSAddr)
        if err != nil {
            logger.Error("failed to listen", "addr", cfg.HTTPSAddr, "error", err)
            os.Exit(1)
        }
    }

    // Layer defenses onto the TCP listener
    protected := ddos.NewProtectedListener(baseListener, 5)
    slowlorisDefended := ddos.NewSlowlorisDetector(protected, 50)
    tlsListener := tls.NewListener(slowlorisDefended, tlsConfig)

    // ── 18. Shutdown Manager ───────────────────────────────────────────────
    shutdownMgr := ops.NewShutdownManager(httpServer, logger)
    shutdownMgr.AddHook(func() {
        cidrRL.Stop()
        trafficAnalyzer.Stop()
        blockList.Stop()
        botScorer.Stop()
        threatIntel.Stop()
        logger.Info("all background workers stopped")
    })
    httpServer.Handler = shutdownMgr.TrackRequest(mux)

    // ── 19. Signal Dispatcher ──────────────────────────────────────────────
    sigDispatcher := ops.NewSignalDispatcher(logger)

    sigDispatcher.OnConfigReload(func() error {
        newCfg, err := ops.LoadConfig("config.json")
        if err != nil {
            return err
        }
        if err := ops.ValidateConfig(newCfg); err != nil {
            return err
        }
        // Atomically swap hot-reloadable components
        wafMiddleware.Reload(newCfg.WAFConfig)
        l7RL.Reload(newCfg.RateLimitConfig)
        chaosCfg.Reload(newCfg.ChaosConfig)
        logger.Info("config reloaded")
        return nil
    })

    sigDispatcher.OnLogRotate(func() error {
        logger.Info("log rotation triggered")
        // In production: reopen log file FD
        return nil
    })

    // Handle SIGUSR1 for zero-downtime restart
    go func() {
        upgradeCh := make(chan os.Signal, 1)
        signal.Notify(upgradeCh, syscall.SIGUSR1)
        <-upgradeCh
        logger.Info("SIGUSR1: initiating zero-downtime restart")
        if err := ops.StartNewProcess(baseListener); err != nil {
            logger.Error("failed to start new process", "error", err)
            return
        }
        logger.Info("new process started, draining old process")
        shutdownMgr.GracefulShutdown()
    }()

    go sigDispatcher.Start()

    // ── 20. Expvar Metrics ─────────────────────────────────────────────────
    expvar.Publish("goroutine_count", expvar.Func(func() interface{} {
        return runtime.NumGoroutine()
    }))
    expvar.Publish("uptime_seconds", expvar.Func(func() interface{} {
        return time.Since(startTime).Seconds()
    }))

    // ── 21. Start ──────────────────────────────────────────────────────────
    errCh := make(chan error, 2)

    go func() {
        logger.Info("HTTPS server starting", "addr", cfg.HTTPSAddr)
        errCh <- httpServer.Serve(tlsListener)
    }()

    go func() {
        logger.Info("HTTP redirect server starting", "addr", cfg.HTTPAddr)
        errCh <- redirectServer.ListenAndServe()
    }()

    startup.MarkReady()
    logger.Info("proxy ready")

    // Wait for shutdown signal or server error
    go shutdownMgr.WaitForSignal()

    err = <-errCh
    if err != nil && err != http.ErrServerClosed {
        logger.Error("server error", "error", err)
        os.Exit(1)
    }

    logger.Info("shutdown complete")
}

var startTime = time.Now()

// ── Helpers ──────────────────────────────────────────────────────────────────

func ddosGate(
    rl *ddos.CIDRRateLimiter,
    ta *ddos.TrafficAnalyzer,
    tr *ddos.TieredResponder,
    ti *threat.ThreatIntel,
    next http.Handler,
) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        ta.Record(ip)
        if !rl.Allow(r) {
            ti.ReportSignal(ip, "rate_limit", 10)
            w.Header().Set("Retry-After", "1")
            http.Error(w, `{"error":"rate_limit_exceeded"}`, http.StatusTooManyRequests)
            return
        }
        score := ti.Score(ip)
        tr.Handle(score, ip, w, r, next)
    })
}

func buildChain(handlers ...func(http.Handler) http.Handler) http.Handler {
    // Final handler is always a 404
    var h http.Handler = http.NotFoundHandler()
    // Wrap in reverse order so first handler in list is outermost
    for i := len(handlers) - 1; i >= 0; i-- {
        h = handlers[i](h)
    }
    return h
}

func httpsRedirect(httpsAddr string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        host := r.Host
        if h, _, err := net.SplitHostPort(httpsAddr); err == nil {
            host = h
        }
        target := "https://" + host + r.URL.RequestURI()
        http.Redirect(w, r, target, http.StatusMovedPermanently)
    }
}
```

---

### Final Requirements Checklist

**Transport & L4**
- [ ] TLS termination (TLS 1.2 minimum, TLS 1.3 preferred)
- [ ] mTLS support (`RequireAndVerifyClientCert` via config)
- [ ] SNI routing (hostname → backend pool)
- [ ] Cert hot-reload (SIGHUP, no dropped connections)
- [ ] L4 half-open connection limiting (SYN flood defense)
- [ ] Slowloris detection (byte rate monitor, close slow conns)
- [ ] Circuit breaker per upstream backend

**DDoS Mitigation**
- [ ] Token bucket per /32 IP and per /24 subnet
- [ ] IPv6 /48 grouping
- [ ] EMA z-score traffic anomaly detection
- [ ] PoW challenge (score 50–80, SHA256, difficulty=16)
- [ ] IP block (score > 80, 5-minute expiry)
- [ ] CIDR whitelist (bypass all mitigations)
- [ ] `/debug/ddos` stats endpoint

**Bot Detection**
- [ ] JA3-like TLS fingerprint via `GetConfigForClient`
- [ ] HTTP header consistency fingerprinting
- [ ] IP ASN/datacenter classification
- [ ] Session behavioral scoring (velocity, path diversity, referrer)
- [ ] Tiered mitigation: delay → PoW → block → shadow-ban
- [ ] Honeypot endpoint (any hit → permanent shadow-ban)
- [ ] PoW challenge endpoint (`/challenge`) with HMAC-signed challenges
- [ ] Session cookie issued after PoW verification
- [ ] `/debug/bots` stats endpoint

**WAF (OWASP)**
- [ ] Multi-layer normalization (URL → HTML → null bytes → lowercase → whitespace)
- [ ] SQLi detection (scoring: keywords + tautology + comment/quote balance)
- [ ] XSS detection (scoring: script tags + event handlers + JS URIs)
- [ ] Path traversal detection (immediate block, all encoding variants)
- [ ] Command injection detection (immediate block)
- [ ] SSRF detection (internal IP/metadata access)
- [ ] Body inspection (JSON + form, up to 1MB via `io.LimitReader`)
- [ ] CIDR allow-list bypass (trusted internal callers)
- [ ] Audit log (every match: rule, value, normalized, score, action)
- [ ] Config hot-reload (thresholds from JSON, SIGHUP)

**API Security**
- [ ] HMAC-SHA256 request signing validation (AWS SigV4-style)
- [ ] ECDSA P-256 token auth (without JWT library)
- [ ] mTLS client cert identity extraction
- [ ] Per-API-key token bucket rate limiting (10/s, burst 50)
- [ ] Per-IP sliding window rate limiting
- [ ] JSON schema validation per route (`additionalProperties: false`)
- [ ] Struct tag input validation (required/min/max/pattern)
- [ ] Response PII scrubber (mask sensitive fields in logs + responses)
- [ ] Credential stuffing defense (per-username + per-IP fail counters)
- [ ] CORS strict origin allowlist (`Vary: Origin` always set)
- [ ] Security headers (HSTS, CSP, X-Frame-Options, X-Content-Type-Options)

**Observability & Operations**
- [ ] Structured logging (`log/slog`, JSON, request ID, client IP, user identity)
- [ ] `expvar` metrics (requests, errors, latency ring buffer, blocked IPs, goroutines)
- [ ] pprof endpoints (`/debug/pprof`) for CPU/heap/goroutine profiling
- [ ] W3C Trace Context propagation (`traceparent` header)
- [ ] In-memory span ring buffer (last 1000), `/debug/traces`
- [ ] Health checks (`/healthz`, `/readyz`, `/startup`, per-subsystem)
- [ ] 3-stage graceful shutdown (accept→drain→force)
- [ ] Signal handlers: SIGHUP config reload, SIGUSR1 log rotate + FD restart, SIGUSR2 goroutine dump
- [ ] Zero-downtime restart via FD inheritance
- [ ] Chaos middleware (configurable fault injection, off by default)
- [ ] ThreatIntel aggregated score with exponential decay

**Quality Gates**
- [ ] `go build ./...` — clean build
- [ ] `go vet ./...` — zero warnings
- [ ] `go test -race ./...` — zero data races
- [ ] Security packages ≥ 80% test coverage
- [ ] Fuzz tests: WAF normalization, schema validator, binary parser
- [ ] Attack simulation: ≥ 80% block rate for all 8 attack types
- [ ] Load test: 500 concurrent clients, p99 latency overhead < 10ms

---

### What You've Built

This capstone integrates every component from all 17 weeks:

| Week | Component | In Capstone |
|------|-----------|-------------|
| 1–4 | Go fundamentals: generics, sync, context, interfaces | Everywhere |
| 5 | TCP networking, connection pools | L4 listener |
| 6 | TLS, mTLS, PKI, ECDSA tokens | TLS layer, auth |
| 7 | L4 load balancing, circuit breaker, health checks | Upstream proxy |
| 8 | Binary protocol, SNI, cert hot-reload | TLS + routing |
| 9 | HTTP internals, reverse proxy, WebSocket | Reverse proxy |
| 10 | WAF: SQLi, XSS, traversal, normalization | WAF middleware |
| 11 | API auth, rate limiting, schema validation, PII | API security layer |
| 12 | ThreatIntel, slog, config hot-reload, integration tests | Cross-cutting |
| 13 | DDoS: rate limiting algorithms, SYN flood, slowloris, PoW | DDoS layer |
| 14 | Bot detection: JA3, headers, session behavior, shadow-ban | Bot layer |
| 15 | Health checks, graceful shutdown, signals, traces, FD inherit | Operations |
| 16 | Attack simulation, profiling, security audit, fuzz tests | Quality gates |

---

### Build vs. Buy Decision Table

| Component | Build cost | Vendor | Buy cost | Verdict |
|-----------|-----------|--------|----------|---------|
| WAF rules database | High (maintain rules) | AWS WAF, Cloudflare | $0.60/1M req | Buy for rule DB; build engine |
| DDoS volumetric | Can't build (bandwidth) | Cloudflare, AWS Shield | $3k–$100k/mo | Buy |
| DDoS L7 (algorithmic) | 2 weeks (this plan) | Included in WAF | N/A | Build |
| Bot detection (basic) | 1 week (this plan) | DataDome, PerimeterX | $50k+/yr | Build for internal APIs |
| Bot detection (browser JS) | Can't build (frontend required) | DataDome, PerimeterX | $50k+/yr | Buy for user-facing |
| API gateway | 1 week (this plan) | Kong, AWS API GW | $200/mo+ | Build for full control |
| TLS termination | 3 days (this plan) | Any reverse proxy | Free (nginx) | Build (you understand it now) |
| IP reputation feeds | Can't maintain | MaxMind, ipinfo.io | $24–$400/mo | Buy data, build integration |

---

### Engineer Takeaway

You have now built the core of what Cloudflare's edge, AWS WAF + ALB, and Fastly's compute platform do — in ~5000 lines of Go, zero external dependencies, using nothing but the standard library.

This doesn't mean you should replace Cloudflare with your own proxy in production. It means:

1. **You understand the products you use at code level.** You know exactly what a "WAF" is, what "bot management" means, and what "DDoS scrubbing" involves, because you built each component.

2. **You can back technical decisions with measurements.** "Our WAF adds 300µs p99" is a number you measured, not a vendor claim you're accepting on faith.

3. **You can architect security infrastructure correctly.** You know why middleware ordering matters, why HMAC is not the same as ECDSA, why token bucket and sliding window have different burst behavior, and why PoW challenges are better than hard blocks for DDoS.

4. **You can own security infrastructure code.** You've written every layer yourself. You can review PRs, spot subtle bugs (timing attacks, goroutine leaks, integer overflows), and make design decisions from experience.

5. **You can debug production incidents.** When users report "the WAF is blocking legitimate traffic," you know which normalization pass to examine. When someone says "we're getting slowloris'd," you know what to check in the connection metrics.

---

### 120-Day Skills Recap

```
Phase 1 (Weeks 1-4):    Go fundamentals, concurrency, testing, profiling
Phase 2 (Weeks 5-8):    TCP/UDP, TLS/mTLS, PKI, L4 load balancing, protocols
Phase 3 (Weeks 9-12):   HTTP internals, WAF, API security, observability
Phase 4 (Weeks 13-17):  DDoS mitigation, bot detection, production hardening, capstone

Total: 120 days × 1 hour = 120 hours of deliberate practice
Output: 1 production-ready security proxy, ~5000 lines of Go, 0 external dependencies
Outcome: a software engineer who can build, review, and operate security infrastructure
```
