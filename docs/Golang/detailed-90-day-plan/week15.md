# Week 15: Production Hardening (Days 106–112)

> **Phase 4 continuation:** Software that works on your laptop is not production software. This week adds the operational characteristics that separate a prototype from a service that can run in production for years.

**Why this week matters:** SRE teams measure you by uptime, deployment safety, and operational controllability. Every exercise this week implements a capability your SRE team will demand before approving a production launch.

---

> **Background in this wiki:** Linux guide [Ch 74 (signals and PID 1)](../../os-linux/real-life-os-guide.md#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1), [Ch 78 (file descriptors)](../../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections), [Ch 81 (pprof leak hunting)](../../os-linux/real-life-os-guide.md#chapter-81-profiling-go-on-linux-pprof-the-execution-tracer-and-perf), [Ch 82 (image, systemd, Kubernetes)](../../os-linux/real-life-os-guide.md#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes); TCP/IP guide [Ch 57 (draining)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries); HTTPS guide [Ch 19 (graceful shutdown)](../../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown), [Ch 21 (circuit breakers)](../../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers), [Ch 23 (SLOs)](../../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services). Security in Depth: [Chapter 52 (tamper-evident audit logs, Go lab)](../../security/real-life-security-guide-v1.md#chapter-52-telemetry-architecture).

## Day 106 — Health Checks & Readiness Probes

| | |
|---|---|
| **Learning objective** | Build a composable health check framework with per-subsystem status |
| **Key packages** | `net/http`, `context`, `encoding/json`, `sync`, `time` |
| **Security context** | Kubernetes routes traffic based on readiness probes — misconfigured probes cause user-visible failures |

### Key Concepts

**Kubernetes probe types:**

| Probe | Question | Failure action |
|-------|----------|---------------|
| **Liveness** | Is the process alive and not deadlocked? | Kill + restart container |
| **Readiness** | Is the service ready for traffic? | Remove from load balancer pool |
| **Startup** | Has the service finished startup? | Don't kill during slow init |

**Rule:** Liveness checks only detect **catastrophic** failure (deadlock, panic). Don't check upstream dependencies in liveness — a slow database should not cause your pod to restart.

### Exercise — Health Check Framework

```go
package health

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// CheckFunc is a function that checks if a subsystem is healthy.
// It must respect context cancellation and return within the deadline.
type CheckFunc func(ctx context.Context) error

type checkEntry struct {
    name    string
    fn      CheckFunc
    timeout time.Duration
    // For readiness: is this required for traffic?
    required bool
}

// Registry holds all registered health checks
type Registry struct {
    mu     sync.RWMutex
    checks []checkEntry
}

func NewRegistry() *Registry {
    return &Registry{}
}

// Register adds a check. Use required=true for readiness-gating checks.
func (r *Registry) Register(name string, fn CheckFunc, timeout time.Duration, required bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.checks = append(r.checks, checkEntry{name: name, fn: fn, timeout: timeout, required: required})
}

type CheckResult struct {
    Status  string `json:"status"` // "ok" or "fail"
    Error   string `json:"error,omitempty"`
    Latency string `json:"latency_ms"`
}

type HealthResponse struct {
    Status string                 `json:"status"` // "ok" or "degraded" or "unhealthy"
    Checks map[string]CheckResult `json:"checks"`
    Took   string                 `json:"took_ms"`
}

func (r *Registry) run(ctx context.Context) HealthResponse {
    r.mu.RLock()
    checks := make([]checkEntry, len(r.checks))
    copy(checks, r.checks)
    r.mu.RUnlock()

    start := time.Now()
    results := make(map[string]CheckResult, len(checks))
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, c := range checks {
        wg.Add(1)
        go func(c checkEntry) {
            defer wg.Done()
            checkCtx, cancel := context.WithTimeout(ctx, c.timeout)
            defer cancel()

            checkStart := time.Now()
            err := c.fn(checkCtx)
            latency := time.Since(checkStart)

            result := CheckResult{
                Latency: fmt.Sprintf("%.1f", float64(latency.Microseconds())/1000.0),
            }
            if err != nil {
                result.Status = "fail"
                result.Error = err.Error()
            } else {
                result.Status = "ok"
            }

            mu.Lock()
            results[c.name] = result
            mu.Unlock()
        }(c)
    }
    wg.Wait()

    // Aggregate status
    status := "ok"
    for _, c := range checks {
        res := results[c.name]
        if res.Status == "fail" && c.required {
            status = "unhealthy"
            break
        } else if res.Status == "fail" {
            if status == "ok" {
                status = "degraded"
            }
        }
    }

    return HealthResponse{
        Status: status,
        Checks: results,
        Took:   fmt.Sprintf("%.1f", float64(time.Since(start).Microseconds())/1000.0),
    }
}

// LivezHandler handles GET /healthz (liveness — is process alive and not deadlocked?)
// Only runs lightweight checks — never checks upstream dependencies.
func (r *Registry) LivezHandler(w http.ResponseWriter, req *http.Request) {
    // Liveness is simple: if this handler runs, the process is alive.
    // Run only non-required checks with very short timeouts.
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ReadyzHandler handles GET /readyz (readiness — ready for traffic?)
// Runs all required checks. Fails if any required check fails.
func (r *Registry) ReadyzHandler(w http.ResponseWriter, req *http.Request) {
    ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
    defer cancel()

    resp := r.run(ctx)
    w.Header().Set("Content-Type", "application/json")
    if resp.Status == "unhealthy" {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    json.NewEncoder(w).Encode(resp)
}

// StartupHandler handles GET /startup — allows slow initialization.
// After startup completes, this always returns 200.
type StartupHandler struct {
    ready chan struct{}
    once  sync.Once
}

func NewStartupHandler() *StartupHandler {
    return &StartupHandler{ready: make(chan struct{})}
}

func (sh *StartupHandler) MarkReady() {
    sh.once.Do(func() { close(sh.ready) })
}

func (sh *StartupHandler) Handler(w http.ResponseWriter, r *http.Request) {
    select {
    case <-sh.ready:
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ready"}`))
    default:
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte(`{"status":"starting"}`))
    }
}

// Example checks for common subsystems
func NewTCPCheck(addr string) CheckFunc {
    return func(ctx context.Context) error {
        d := net.Dialer{}
        conn, err := d.DialContext(ctx, "tcp", addr)
        if err != nil {
            return fmt.Errorf("tcp connect to %s: %w", addr, err)
        }
        conn.Close()
        return nil
    }
}

func NewHTTPCheck(url string) CheckFunc {
    client := &http.Client{Timeout: 2 * time.Second}
    return func(ctx context.Context) error {
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return err
        }
        resp, err := client.Do(req)
        if err != nil {
            return fmt.Errorf("http get %s: %w", url, err)
        }
        resp.Body.Close()
        if resp.StatusCode >= 500 {
            return fmt.Errorf("upstream returned %d", resp.StatusCode)
        }
        return nil
    }
}

func NewGoroutineLeakCheck(baseline int, maxMultiplier float64) CheckFunc {
    return func(ctx context.Context) error {
        current := runtime.NumGoroutine()
        limit := int(float64(baseline) * maxMultiplier)
        if current > limit {
            return fmt.Errorf("goroutine count %d exceeds %dx baseline (%d)", current, int(maxMultiplier), baseline)
        }
        return nil
    }
}
```

---

## Day 107 — Multi-Stage Graceful Shutdown

| | |
|---|---|
| **Learning objective** | Implement a 3-stage shutdown sequence with per-stage timeout budgets |
| **Key packages** | `net/http`, `context`, `os/signal`, `sync`, `sync/atomic`, `time` |
| **Security context** | `kill -9` drops in-flight requests; `SIGTERM` should complete in-flight work and close connections cleanly |

### Key Concepts

**3-stage shutdown:**
1. **Stage 1 (stop accepting new):** Close listener, stop accepting new connections. Budget: 0s (immediate).
2. **Stage 2 (drain in-flight):** Wait for active requests to complete. Budget: 30s.
3. **Stage 3 (force-close):** Force-close any remaining connections. Budget: 5s.

### Exercise — 3-Stage Shutdown

```go
package ops

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
)

// ShutdownManager orchestrates the 3-stage shutdown sequence
type ShutdownManager struct {
    server      *http.Server
    inFlight    atomic.Int64
    drainBudget time.Duration
    forceBudget time.Duration
    logger      *slog.Logger
    hooks       []func() // called in order during shutdown (e.g. close DB connections)
}

func NewShutdownManager(server *http.Server, logger *slog.Logger) *ShutdownManager {
    return &ShutdownManager{
        server:      server,
        drainBudget: 30 * time.Second,
        forceBudget: 5 * time.Second,
        logger:      logger,
    }
}

func (sm *ShutdownManager) AddHook(fn func()) {
    sm.hooks = append(sm.hooks, fn)
}

// TrackRequest wraps a handler, tracking in-flight request count
func (sm *ShutdownManager) TrackRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sm.inFlight.Add(1)
        defer sm.inFlight.Add(-1)
        next.ServeHTTP(w, r)
    })
}

// WaitForSignal blocks until SIGTERM or SIGINT, then executes shutdown sequence
func (sm *ShutdownManager) WaitForSignal() {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

    sig := <-sigCh
    sm.logger.Info("shutdown signal received", "signal", sig.String())

    sm.shutdown()
}

func (sm *ShutdownManager) shutdown() {
    start := time.Now()

    // Stage 1: Stop accepting new connections (immediate)
    sm.logger.Info("shutdown stage 1: stopping acceptance of new connections")
    drainCtx, drainCancel := context.WithTimeout(context.Background(), sm.drainBudget)
    defer drainCancel()

    // http.Server.Shutdown stops accepting and waits for in-flight to complete
    shutdownDone := make(chan error, 1)
    go func() {
        shutdownDone <- sm.server.Shutdown(drainCtx)
    }()

    // Stage 2: Wait for in-flight requests to drain
    sm.logger.Info("shutdown stage 2: draining in-flight requests",
        "in_flight", sm.inFlight.Load(),
        "budget_s", sm.drainBudget.Seconds(),
    )

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

drainLoop:
    for {
        select {
        case <-drainCtx.Done():
            sm.logger.Warn("shutdown stage 2: drain budget exhausted",
                "remaining_in_flight", sm.inFlight.Load(),
            )
            break drainLoop
        case err := <-shutdownDone:
            if err != nil {
                sm.logger.Error("shutdown stage 2: server shutdown error", "error", err)
            } else {
                sm.logger.Info("shutdown stage 2: all requests drained",
                    "elapsed_ms", time.Since(start).Milliseconds(),
                )
            }
            break drainLoop
        case <-ticker.C:
            remaining := sm.inFlight.Load()
            if remaining == 0 {
                break drainLoop
            }
            sm.logger.Info("shutdown stage 2: waiting", "in_flight", remaining)
        }
    }

    // Stage 3: Force-close remaining connections
    sm.logger.Info("shutdown stage 3: force closing remaining connections")
    forceCtx, forceCancel := context.WithTimeout(context.Background(), sm.forceBudget)
    defer forceCancel()
    _ = forceCtx // In production: use this to force-close lingering connections

    // Run shutdown hooks (close DB pools, flush metrics, etc.)
    var wg sync.WaitGroup
    for i, hook := range sm.hooks {
        wg.Add(1)
        i, hook := i, hook
        go func() {
            defer wg.Done()
            sm.logger.Info("shutdown hook running", "index", i)
            hook()
        }()
    }

    hookDone := make(chan struct{})
    go func() {
        wg.Wait()
        close(hookDone)
    }()

    select {
    case <-hookDone:
        sm.logger.Info("shutdown complete", "total_elapsed_ms", time.Since(start).Milliseconds())
    case <-time.After(sm.forceBudget):
        sm.logger.Warn("shutdown hooks timed out", "budget_s", sm.forceBudget.Seconds())
    }
}
```

---

## Day 108 — Signal Handling & Operational Controls

| | |
|---|---|
| **Learning objective** | Handle SIGHUP (config reload), SIGUSR1 (log rotate), SIGUSR2 (goroutine dump) |
| **Key packages** | `os/signal`, `syscall`, `runtime`, `log/slog` |
| **Security context** | Ops teams must be able to reload security rules and rotate logs without restarts |

### Key Concepts

| Signal | Convention | Action |
|--------|-----------|--------|
| `SIGTERM` | Kubernetes, systemd | Graceful shutdown |
| `SIGINT` | Ctrl-C | Graceful shutdown |
| `SIGHUP` | Log rotation (historically) | Reload config |
| `SIGUSR1` | User-defined 1 | Log file rotation |
| `SIGUSR2` | User-defined 2 | Goroutine stack dump |

### Exercise — Signal Handler Dispatcher

```go
package ops

import (
    "bytes"
    "fmt"
    "log/slog"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "time"
)

// SignalDispatcher routes OS signals to registered handlers
type SignalDispatcher struct {
    configReloader func() error
    logRotator     func() error
    logger         *slog.Logger
    done           chan struct{}
}

func NewSignalDispatcher(logger *slog.Logger) *SignalDispatcher {
    return &SignalDispatcher{logger: logger, done: make(chan struct{})}
}

func (sd *SignalDispatcher) OnConfigReload(fn func() error) { sd.configReloader = fn }
func (sd *SignalDispatcher) OnLogRotate(fn func() error)    { sd.logRotator = fn }

// Start begins listening for signals. Call in a goroutine.
func (sd *SignalDispatcher) Start() {
    sigCh := make(chan os.Signal, 8)
    signal.Notify(sigCh,
        syscall.SIGHUP,
        syscall.SIGUSR1,
        syscall.SIGUSR2,
    )
    defer signal.Stop(sigCh)

    for {
        select {
        case <-sd.done:
            return
        case sig := <-sigCh:
            sd.dispatch(sig)
        }
    }
}

func (sd *SignalDispatcher) dispatch(sig os.Signal) {
    sd.logger.Info("signal received", "signal", sig.String())
    switch sig {
    case syscall.SIGHUP:
        sd.handleConfigReload()
    case syscall.SIGUSR1:
        sd.handleLogRotate()
    case syscall.SIGUSR2:
        sd.handleGoroutineDump()
    }
}

func (sd *SignalDispatcher) handleConfigReload() {
    if sd.configReloader == nil {
        sd.logger.Warn("SIGHUP received but no config reloader registered")
        return
    }
    sd.logger.Info("reloading config")
    if err := sd.configReloader(); err != nil {
        // CRITICAL: keep serving with old config — never apply invalid config
        sd.logger.Error("config reload failed — keeping old config", "error", err)
        return
    }
    sd.logger.Info("config reloaded successfully")
}

func (sd *SignalDispatcher) handleLogRotate() {
    if sd.logRotator == nil {
        sd.logger.Warn("SIGUSR1 received but no log rotator registered")
        return
    }
    sd.logger.Info("rotating logs")
    if err := sd.logRotator(); err != nil {
        sd.logger.Error("log rotation failed", "error", err)
        return
    }
    sd.logger.Info("log rotation complete")
}

func (sd *SignalDispatcher) handleGoroutineDump() {
    buf := make([]byte, 1<<20) // 1MB buffer
    n := runtime.Stack(buf, true)
    dump := string(buf[:n])

    // Count goroutines
    lines := bytes.Split(buf[:n], []byte("\n"))
    goroutineCount := 0
    for _, line := range lines {
        if bytes.HasPrefix(line, []byte("goroutine ")) {
            goroutineCount++
        }
    }

    sd.logger.Info("goroutine dump",
        "goroutine_count", goroutineCount,
        "dump_bytes", n,
    )

    // Write full dump to a timestamped file
    filename := fmt.Sprintf("/tmp/goroutine-dump-%s.txt", time.Now().Format("20060102-150405"))
    if err := os.WriteFile(filename, []byte(dump), 0600); err != nil {
        sd.logger.Error("failed to write goroutine dump", "error", err)
        return
    }
    sd.logger.Info("goroutine dump written", "file", filename)
}

func (sd *SignalDispatcher) Stop() { close(sd.done) }
```

---

## Day 109 — Observability: Traces & Spans

| | |
|---|---|
| **Learning objective** | Implement W3C Trace Context propagation with in-memory span recording |
| **Key packages** | `net/http`, `context`, `encoding/hex`, `crypto/rand`, `encoding/json`, `sync` |
| **Security context** | Distributed tracing exposes security middleware latency per layer — required for SLA debugging |

### Key Concepts

**W3C Trace Context `traceparent` header format:**
```
traceparent: 00-{traceId 32 hex}-{parentId 16 hex}-{flags 2 hex}
             00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
```

- `traceId`: 128-bit random ID (same across entire request chain)
- `parentId`: 64-bit ID for current span
- `flags`: `01` = sampled

### Exercise — Manual Tracer with Ring Buffer

```go
package tracing

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

type contextKey int

const traceKey contextKey = 0

type SpanContext struct {
    TraceID  string
    SpanID   string
    ParentID string
    Sampled  bool
}

type Span struct {
    SpanContext
    Name       string
    StartTime  time.Time
    Duration   time.Duration
    Attributes map[string]string
    Error      string
}

func newID(n int) string {
    b := make([]byte, n)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// Tracer records spans in a ring buffer
type Tracer struct {
    mu      sync.Mutex
    spans   []*Span
    head    int
    size    int
    sampleRate float64 // 0.0–1.0
}

func NewTracer(bufferSize int, sampleRate float64) *Tracer {
    return &Tracer{
        spans:      make([]*Span, bufferSize),
        size:       bufferSize,
        sampleRate: sampleRate,
    }
}

func (t *Tracer) record(s *Span) {
    t.mu.Lock()
    t.spans[t.head%t.size] = s
    t.head++
    t.mu.Unlock()
}

// StartSpan creates a new span, returning a finish function
func (t *Tracer) StartSpan(ctx context.Context, name string) (context.Context, func(attrs map[string]string, errMsg string)) {
    parent, _ := ctx.Value(traceKey).(*SpanContext)

    span := &Span{
        Name:       name,
        StartTime:  time.Now(),
        Attributes: make(map[string]string),
    }

    if parent != nil {
        span.TraceID = parent.TraceID
        span.ParentID = parent.SpanID
        span.Sampled = parent.Sampled
    } else {
        span.TraceID = newID(16)
        span.Sampled = true // simplification: sample all root spans
    }
    span.SpanID = newID(8)

    childCtx := context.WithValue(ctx, traceKey, &span.SpanContext)

    finish := func(attrs map[string]string, errMsg string) {
        span.Duration = time.Since(span.StartTime)
        for k, v := range attrs {
            span.Attributes[k] = v
        }
        if errMsg != "" {
            span.Error = errMsg
        }
        if span.Sampled {
            t.record(span)
        }
    }

    return childCtx, finish
}

// Middleware injects/propagates trace context
func (t *Tracer) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var sc *SpanContext

        // Try to propagate incoming traceparent
        tp := r.Header.Get("Traceparent")
        if tp != "" {
            sc = parseTraceparent(tp)
        }

        if sc == nil {
            sc = &SpanContext{
                TraceID: newID(16),
                SpanID:  newID(8),
                Sampled: true,
            }
        }

        ctx := context.WithValue(r.Context(), traceKey, sc)
        r = r.WithContext(ctx)

        // Set traceparent on response for downstream propagation
        w.Header().Set("Traceparent", formatTraceparent(sc))

        ctx, finish := t.StartSpan(ctx, "http.request")
        defer finish(map[string]string{
            "http.method": r.Method,
            "http.path":   r.URL.Path,
        }, "")
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}

func parseTraceparent(tp string) *SpanContext {
    // 00-{32hex}-{16hex}-{2hex}
    if len(tp) != 55 || tp[2] != '-' || tp[35] != '-' || tp[52] != '-' {
        return nil
    }
    return &SpanContext{
        TraceID:  tp[3:35],
        ParentID: tp[36:52],
        Sampled:  tp[53:55] == "01",
    }
}

func formatTraceparent(sc *SpanContext) string {
    flags := "00"
    if sc.Sampled {
        flags = "01"
    }
    spanID := sc.SpanID
    if spanID == "" {
        spanID = newID(8)
    }
    return fmt.Sprintf("00-%s-%s-%s", sc.TraceID, spanID, flags)
}

// TracesHandler serves recent spans at /debug/traces
func (t *Tracer) TracesHandler(w http.ResponseWriter, r *http.Request) {
    t.mu.Lock()
    result := make([]*Span, 0, t.size)
    for i := 0; i < t.size; i++ {
        idx := (t.head - 1 - i + t.size*2) % t.size
        if t.spans[idx] != nil {
            result = append(result, t.spans[idx])
        }
    }
    t.mu.Unlock()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

---

## Day 110 — Chaos Engineering

| | |
|---|---|
| **Learning objective** | Build configurable fault injection and a circuit breaker |
| **Key packages** | `net/http`, `math/rand`, `sync/atomic`, `time` |
| **Security context** | Chaos testing exposes brittleness in security middleware before production does |

### Exercise — Chaos Middleware + Circuit Breaker

```go
package chaos

import (
    "math/rand"
    "net/http"
    "sync/atomic"
    "time"
)

// ChaosConfig controls fault injection rates (all zero by default — safe to deploy)
type ChaosConfig struct {
    ErrorRate   float64       // 0.0–1.0: probability of returning 500
    DelayRate   float64       // 0.0–1.0: probability of injecting delay
    MinDelay    time.Duration
    MaxDelay    time.Duration
    CloseRate   float64       // 0.0–1.0: probability of hijacking and closing connection
    Enabled     atomic.Bool   // master kill switch
}

func NewChaosConfig() *ChaosConfig {
    cfg := &ChaosConfig{
        ErrorRate: 0.05, MaxDelay: 2 * time.Second, MinDelay: 100 * time.Millisecond,
        DelayRate: 0.02, CloseRate: 0.01,
    }
    // Disabled by default — enable via CHAOS_ENABLED env var or config
    return cfg
}

// ChaosMiddleware injects random faults per ChaosConfig
func ChaosMiddleware(cfg *ChaosConfig, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !cfg.Enabled.Load() {
            next.ServeHTTP(w, r)
            return
        }

        roll := rand.Float64()

        // Error injection
        if roll < cfg.ErrorRate {
            http.Error(w, `{"error":"chaos_injected"}`, http.StatusInternalServerError)
            return
        }

        // Latency injection
        if roll < cfg.ErrorRate+cfg.DelayRate {
            delay := cfg.MinDelay + time.Duration(rand.Int63n(int64(cfg.MaxDelay-cfg.MinDelay)))
            time.Sleep(delay)
        }

        // Connection close injection — requires Hijacker
        if roll < cfg.ErrorRate+cfg.DelayRate+cfg.CloseRate {
            if h, ok := w.(http.Hijacker); ok {
                conn, _, err := h.Hijack()
                if err == nil {
                    conn.Close()
                    return
                }
            }
        }

        next.ServeHTTP(w, r)
    })
}

// CircuitBreaker — 3 states: closed (normal), open (blocking), half-open (testing)
type CircuitBreaker struct {
    failures    atomic.Int64
    successes   atomic.Int64
    state       atomic.Int32 // 0=closed, 1=open, 2=half-open
    lastFailure atomic.Int64 // unix nano

    failureThreshold int64
    successThreshold int64
    openDuration     time.Duration
    halfOpenAllowed  atomic.Int32 // permit count in half-open state
}

const (
    stateClosed   int32 = 0
    stateOpen     int32 = 1
    stateHalfOpen int32 = 2
)

func NewCircuitBreaker(failureThreshold, successThreshold int, openDuration time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: int64(failureThreshold),
        successThreshold: int64(successThreshold),
        openDuration:     openDuration,
    }
}

// Allow returns true if the request should proceed
func (cb *CircuitBreaker) Allow() bool {
    switch cb.state.Load() {
    case stateClosed:
        return true

    case stateOpen:
        // Check if enough time has passed to try half-open
        lastFail := time.Unix(0, cb.lastFailure.Load())
        if time.Since(lastFail) > cb.openDuration {
            if cb.state.CompareAndSwap(stateOpen, stateHalfOpen) {
                cb.halfOpenAllowed.Store(1) // allow one request to test
            }
            return cb.state.Load() == stateHalfOpen && cb.halfOpenAllowed.Add(-1) >= 0
        }
        return false

    case stateHalfOpen:
        return cb.halfOpenAllowed.Add(-1) >= 0
    }
    return true
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
    if cb.state.Load() == stateHalfOpen {
        if cb.successes.Add(1) >= cb.successThreshold {
            cb.successes.Store(0)
            cb.failures.Store(0)
            cb.state.Store(stateClosed)
        }
    }
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
    cb.lastFailure.Store(time.Now().UnixNano())
    if cb.failures.Add(1) >= cb.failureThreshold {
        cb.state.Store(stateOpen)
        cb.failures.Store(0)
        cb.successes.Store(0)
    }
}

func (cb *CircuitBreaker) State() string {
    switch cb.state.Load() {
    case stateClosed:
        return "closed"
    case stateOpen:
        return "open"
    case stateHalfOpen:
        return "half_open"
    }
    return "unknown"
}

// Middleware wraps an HTTP handler with circuit breaker protection
func (cb *CircuitBreaker) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !cb.Allow() {
            http.Error(w, `{"error":"circuit_open"}`, http.StatusServiceUnavailable)
            return
        }

        // Use statusCapture to detect failures
        sc := &statusCapture{ResponseWriter: w}
        next.ServeHTTP(sc, r)

        if sc.status >= 500 {
            cb.RecordFailure()
        } else {
            cb.RecordSuccess()
        }
    })
}

type statusCapture struct {
    http.ResponseWriter
    status int
}

func (sc *statusCapture) WriteHeader(code int) {
    sc.status = code
    sc.ResponseWriter.WriteHeader(code)
}
```

---

## Day 111 — Memory & Goroutine Leak Hunting

| | |
|---|---|
| **Learning objective** | Audit all middleware for goroutine leaks and heap growth under load |
| **Key packages** | `runtime`, `runtime/pprof`, `net/http/pprof`, `testing`, `sync` |
| **Security context** | Goroutine leaks are a denial-of-service vector — an attacker can exhaust memory by triggering leak paths |

### Key Concepts

**Common goroutine leak patterns:**

```go
// LEAK: goroutine blocks forever on unbuffered channel with no sender
go func() {
    result := <-ch  // if nobody sends, this goroutine lives forever
    process(result)
}()

// FIX: use context cancellation
go func() {
    select {
    case result := <-ch:
        process(result)
    case <-ctx.Done():
        return
    }
}()

// LEAK: goroutine started without exit condition
go func() {
    for {
        doWork() // if done channel is never closed, goroutine lives forever
    }
}()

// FIX: always provide exit
go func() {
    for {
        select {
        case <-done:
            return
        default:
            doWork()
        }
    }
}()
```

### Exercise — LeakChecker for Load Tests

```go
package ops

import (
    "fmt"
    "runtime"
    "testing"
    "time"
)

// LeakChecker verifies goroutine count returns to baseline after a test
type LeakChecker struct {
    baseline int
    t        *testing.T
}

// NewLeakChecker captures the current goroutine baseline.
// Call at the START of a test, before starting any goroutines.
func NewLeakChecker(t *testing.T) *LeakChecker {
    // Allow a brief settle period for any framework goroutines
    time.Sleep(10 * time.Millisecond)
    return &LeakChecker{baseline: runtime.NumGoroutine(), t: t}
}

// Check verifies goroutine count is within tolerance of baseline.
// Call at the END of a test, after stopping all goroutines.
func (lc *LeakChecker) Check() {
    lc.CheckWithTolerance(5, 5*time.Second)
}

func (lc *LeakChecker) CheckWithTolerance(tolerance int, waitFor time.Duration) {
    deadline := time.Now().Add(waitFor)
    for {
        current := runtime.NumGoroutine()
        if current <= lc.baseline+tolerance {
            lc.t.Logf("goroutine leak check: baseline=%d current=%d (OK)", lc.baseline, current)
            return
        }
        if time.Now().After(deadline) {
            // Capture stack for diagnosis
            buf := make([]byte, 1<<20)
            n := runtime.Stack(buf, true)
            lc.t.Errorf("goroutine leak: baseline=%d current=%d (leaked %d)\n%s",
                lc.baseline, current, current-lc.baseline, buf[:n])
            return
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// GoroutineDelta measures goroutine growth during a function
func GoroutineDelta(fn func()) (delta int, stacks string) {
    before := runtime.NumGoroutine()
    fn()
    time.Sleep(50 * time.Millisecond) // let goroutines start

    after := runtime.NumGoroutine()
    delta = after - before

    buf := make([]byte, 1<<20)
    n := runtime.Stack(buf, true)
    return delta, string(buf[:n])
}

// HeapDiff compares heap allocation before and after fn.
// Returns growth in bytes. Use this to detect unbounded memory growth.
func HeapDiff(fn func()) (heapGrowthBytes uint64) {
    runtime.GC()
    var before runtime.MemStats
    runtime.ReadMemStats(&before)

    fn()

    runtime.GC()
    var after runtime.MemStats
    runtime.ReadMemStats(&after)

    if after.HeapAlloc > before.HeapAlloc {
        return after.HeapAlloc - before.HeapAlloc
    }
    return 0
}

// Example usage in test:
// func TestMiddlewareNoLeak(t *testing.T) {
//     leak := NewLeakChecker(t)
//     // ... run 1000 requests through middleware ...
//     leak.Check() // fails if goroutines leak
// }
```

---

## Day 112 — Zero-Downtime Deployment via FD Inheritance

| | |
|---|---|
| **Learning objective** | Pass a listening socket FD from the old process to a new process without dropping connections |
| **Key packages** | `net`, `os`, `os/exec`, `syscall`, `os/signal` |
| **Security context** | Zero-downtime restarts are required for security patch deployments — you can't afford a restart gap |

### Key Concepts

**FD inheritance flow:**
1. Running process receives `SIGUSR1`
2. Running process starts new process, passing its listener FD via `cmd.ExtraFiles`
3. New process detects FD 3 (first ExtraFile), wraps it as `net.Listener`
4. New process starts serving; old process drains remaining connections and exits

### Exercise — FD Inheritance Server

```go
package ops

import (
    "fmt"
    "net"
    "net/http"
    "os"
    "os/exec"
    "os/signal"
    "strconv"
    "syscall"
    "time"
)

const fdEnvVar = "LISTENER_FD"

// ListenerFromEnv retrieves an inherited listener FD, if present.
// Returns (listener, true) if FD was inherited; (nil, false) if starting fresh.
func ListenerFromEnv() (net.Listener, bool) {
    fdStr := os.Getenv(fdEnvVar)
    if fdStr == "" {
        return nil, false
    }
    fd, err := strconv.Atoi(fdStr)
    if err != nil {
        return nil, false
    }

    // Wrap the FD as a net.Listener
    file := os.NewFile(uintptr(fd), "listener")
    if file == nil {
        return nil, false
    }
    defer file.Close()

    l, err := net.FileListener(file)
    if err != nil {
        return nil, false
    }
    return l, true
}

// StartServer starts the server, using an inherited listener if available
func StartServer(addr string, handler http.Handler) error {
    var l net.Listener
    var err error

    if inherited, ok := ListenerFromEnv(); ok {
        l = inherited
        fmt.Println("server: using inherited listener FD")
    } else {
        l, err = net.Listen("tcp", addr)
        if err != nil {
            return fmt.Errorf("listen %s: %w", addr, err)
        }
        fmt.Printf("server: listening on %s\n", addr)
    }

    server := &http.Server{
        Handler:      handler,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    // Watch for SIGUSR1 (upgrade signal)
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGUSR1)

    go func() {
        <-sigCh
        fmt.Println("server: received upgrade signal, starting new process")

        if err := startNewProcess(l); err != nil {
            fmt.Printf("server: failed to start new process: %v\n", err)
            return
        }

        // Gracefully drain: stop accepting, wait for in-flight
        fmt.Println("server: draining connections (old process)")
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        server.Shutdown(ctx)
    }()

    return server.Serve(l)
}

// startNewProcess starts a new instance of the current binary
// and passes the listener FD via ExtraFiles
func startNewProcess(l net.Listener) error {
    // Get the underlying *os.File from the listener
    tcpListener, ok := l.(*net.TCPListener)
    if !ok {
        return fmt.Errorf("listener is not a *net.TCPListener")
    }
    f, err := tcpListener.File()
    if err != nil {
        return fmt.Errorf("get listener file: %w", err)
    }
    defer f.Close()

    // os.Args[0] is the current binary
    cmd := exec.Command(os.Args[0], os.Args[1:]...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.ExtraFiles = []*os.File{f} // FD 3 (stdin=0, stdout=1, stderr=2, extraFiles[0]=3)
    cmd.Env = append(os.Environ(), fmt.Sprintf("%s=3", fdEnvVar))

    return cmd.Start()
}
```

---

## Week 15 Mini Project: Production-Hardened Gateway

**Take the Week 12 Unified Security Proxy and add all production-hardening features.**

### Feature Addition Checklist

- [ ] Health check framework
  - [ ] `/healthz` — liveness (always fast, no upstream checks)
  - [ ] `/readyz` — readiness (checks all required upstream connections)
  - [ ] `/startup` — startup probe (wait for initialization to complete)
  - [ ] Per-subsystem checks: WAF rule load, TLS cert validity, upstream TCP reachability
- [ ] 3-stage graceful shutdown
  - [ ] Stage 1: Stop accepting (immediate)
  - [ ] Stage 2: Drain in-flight (30s budget)
  - [ ] Stage 3: Force-close + run hooks (5s budget)
- [ ] Signal handlers
  - [ ] `SIGHUP` → reload config from file (validate first, keep old if invalid)
  - [ ] `SIGUSR1` → rotate log file
  - [ ] `SIGUSR2` → dump goroutine stacks to `/tmp/goroutine-dump-{timestamp}.txt`
  - [ ] `SIGTERM/SIGINT` → graceful shutdown
- [ ] W3C trace context propagation
  - [ ] Inject `traceparent` on ingress if missing
  - [ ] Propagate to upstream via `Director` in reverse proxy
  - [ ] Record spans for each middleware layer
  - [ ] `/debug/traces` — last 1000 traces as JSON
- [ ] Chaos middleware (disabled by default)
  - [ ] Configurable error rate, delay rate, connection close rate
  - [ ] Enable via config file (not env var — must be hot-reloadable)
- [ ] Leak validation (automated)
  - [ ] Post-startup goroutine baseline
  - [ ] Expose `goroutine_count` in `/debug/vars`
  - [ ] Alert (log) if count grows > 2× baseline
- [ ] Zero-downtime restart
  - [ ] `SIGUSR1` triggers FD inheritance restart
  - [ ] Old process drains; new process serves

### Operational Test Commands

```bash
# Send SIGHUP to reload config
kill -HUP $(pgrep gateway)

# Send SIGUSR1 to rotate logs
kill -USR1 $(pgrep gateway)

# Send SIGUSR2 to dump goroutines
kill -USR2 $(pgrep gateway)
cat /tmp/goroutine-dump-*.txt | head -50

# Trigger zero-downtime restart
kill -USR1 $(pgrep gateway)  # starts new process
# Old process drains; new process inherits listener

# Health checks
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

### Engineer Takeaway

This week implements what SRE teams call "Day 2 operations" — everything needed to run a service safely in production for months and years. SRE teams at Google, Stripe, and Cloudflare require all of these before a service enters production. A service without graceful shutdown, health checks, and signal handling is a liability, not an asset.

---

## Week 15 Review

| Day | Skill | Real-World Use |
|-----|-------|----------------|
| 106 | Health check framework (liveness/readiness/startup) | Kubernetes traffic routing and restart policies |
| 107 | 3-stage graceful shutdown | Zero dropped connections during deploys |
| 108 | Signal handlers (SIGHUP/SIGUSR1/SIGUSR2) | Runtime operational controls without restarts |
| 109 | W3C Trace Context propagation + span recording | Distributed tracing integration (Jaeger, Datadog) |
| 110 | Chaos middleware + circuit breaker | Production resilience testing |
| 111 | Goroutine and heap leak hunting | Prevent OOM-induced restarts |
| 112 | Zero-downtime deployment via FD inheritance | Deploy security patches without user impact |
