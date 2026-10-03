# Week 12: Phase 3 Review + Integration (Days 82–90)

**Phase:** 3 — L7 Networking + WAF + API Security
**Goal:** Wire all Phase 2 and Phase 3 components together into a single production-grade security proxy. Add structured logging, hot-reload config, integration tests, fuzz tests, and load testing. Day 90 capstone: Unified Security Proxy.

---

> **Background in this wiki:** [OSI guide, Part 15 (architecture decisions by layer)](../../networking/real-life-example-osi.md#part-15-architecture-decisions-by-layer-the-senior-and-manager-view) and HTTPS guide [Ch 18 (end-to-end diagnosis)](../../v2-https/real-life-guide-v1.md#chapter-18-capstone-diagnose-a-slow-broken-or-unsafe-https-request-end-to-end).

## Day 82 — Middleware Ordering & Pipeline Design

### Learning Objectives
- Understand why middleware order is a security architecture decision
- Build a configurable middleware pipeline with load-from-JSON ordering
- Log which middleware rejected each request (for incident response)

### Correct Middleware Order (Security-First)

```
Request flow (each layer can short-circuit):
  1. PanicRecoverer          — outermost: catch panics from ALL other middleware
  2. RequestID               — inject before any logging
  3. CORS Preflight          — handle before auth (CORS is pre-auth by design)
  4. Global Connection Limit — protect server resources
  5. Per-IP Rate Limit       — DDoS mitigation before auth (before you know who it is)
  6. TLS/mTLS Check          — verify transport security
  7. Bot/Scanner Check       — reject known scanners early
  8. Authentication          — verify identity
  9. Per-User Rate Limit     — quota enforcement (after you know who)
 10. WAF                     — inspect payload (after auth so identity appears in WAF logs)
 11. JSON Schema Validation  — structural validation
 12. Input Validation        — semantic validation
 13. Security Headers        — always set on responses
 14. PII Scrubber            — mask before logging
 15. Request Logger          — log after all decisions made (includes final status)
 16. Application Handler
```

### Why Order Matters

```
Wrong order: WAF before Auth
  Problem: WAF blocks legitimate traffic before user can authenticate
  Worse:   WAF logs don't include user identity → harder to investigate false positives

Wrong order: Rate limit after Auth
  Problem: Unauthenticated requests can exhaust rate limit budget for real users
  Attack:  Attacker sends 10,000 unauthenticated requests → legitimate users get 429

Wrong order: Security headers inside handler
  Problem: Panicked handler never sets security headers → XSS unprotected

Correct:
  IP rate limit → Auth → User rate limit
  WAF after Auth → WAF logs include user identity
  PanicRecoverer outermost → Security headers always set
```

### Exercise: Configurable Pipeline

```go
package pipeline

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
)

type PipelineConfig struct {
    Middleware []string `json:"middleware"` // names in order
    Required   []string `json:"required"`   // must be present
}

type Registry struct {
    middleware map[string]func(http.Handler) http.Handler
}

func NewRegistry() *Registry { return &Registry{middleware: make(map[string]func(http.Handler) http.Handler)} }

func (r *Registry) Register(name string, mw func(http.Handler) http.Handler) {
    r.middleware[name] = mw
}

func (r *Registry) Build(configFile string, handler http.Handler) (http.Handler, error) {
    data, err := os.ReadFile(configFile)
    if err != nil { return nil, err }

    var cfg PipelineConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("invalid pipeline config: %w", err)
    }

    // Validate required middleware are present
    nameSet := make(map[string]bool)
    for _, name := range cfg.Middleware { nameSet[name] = true }
    for _, req := range cfg.Required {
        if !nameSet[req] {
            return nil, fmt.Errorf("required middleware %q not in pipeline", req)
        }
    }

    // Build chain (reverse order — first in config = outermost)
    chain := handler
    for i := len(cfg.Middleware) - 1; i >= 0; i-- {
        name := cfg.Middleware[i]
        mw, ok := r.middleware[name]
        if !ok {
            return nil, fmt.Errorf("unknown middleware: %q", name)
        }
        inner := chain
        mwName := name
        chain = mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Track which middleware handles each request
            ctx := context.WithValue(r.Context(), pipelineStageKey, mwName)
            inner.ServeHTTP(w, r.WithContext(ctx))
        }))
    }

    return chain, nil
}
```

**pipeline-config.json:**
```json
{
  "middleware": [
    "panic_recoverer",
    "request_id",
    "cors",
    "ip_rate_limit",
    "auth",
    "user_rate_limit",
    "waf",
    "schema_validation",
    "security_headers",
    "pii_scrubber",
    "request_logger"
  ],
  "required": [
    "panic_recoverer",
    "auth",
    "request_logger"
  ]
}
```

---

## Day 83 — Threat Intelligence Aggregation

### Learning Objectives
- Build a shared threat score per IP that aggregates signals from all security layers
- Implement score decay over time (old signals become less relevant)
- Connect WAF, rate limiter, and auth to the shared ThreatIntel store

### Exercise: ThreatIntel Store

```go
package threatintel

import (
    "fmt"
    "sync"
    "time"
)

type SignalType string

const (
    SignalWAFHit        SignalType = "waf_hit"
    SignalRateLimit     SignalType = "rate_limit"
    SignalAuthFailure   SignalType = "auth_failure"
    SignalScanner       SignalType = "scanner"
    SignalDDoS          SignalType = "ddos"
)

type Signal struct {
    Type   SignalType
    Weight float64
    At     time.Time
}

type IPRecord struct {
    mu      sync.Mutex
    signals []Signal
}

type ThreatIntel struct {
    mu      sync.RWMutex
    records map[string]*IPRecord
    halfLife time.Duration // score halves every halfLife
}

func New(halfLife time.Duration) *ThreatIntel {
    t := &ThreatIntel{
        records:  make(map[string]*IPRecord),
        halfLife: halfLife,
    }
    go t.cleanup()
    return t
}

// ReportSignal adds a threat signal for an IP
func (ti *ThreatIntel) ReportSignal(ip string, t SignalType, weight float64) {
    ti.mu.Lock()
    rec, ok := ti.records[ip]
    if !ok {
        rec = &IPRecord{}
        ti.records[ip] = rec
    }
    ti.mu.Unlock()

    rec.mu.Lock()
    rec.signals = append(rec.signals, Signal{Type: t, Weight: weight, At: time.Now()})
    rec.mu.Unlock()

    fmt.Printf("[THREAT] %s: signal=%s weight=%.1f\n", ip, t, weight)
}

// Score returns the current threat score for an IP (with time decay)
func (ti *ThreatIntel) Score(ip string) float64 {
    ti.mu.RLock()
    rec, ok := ti.records[ip]
    ti.mu.RUnlock()
    if !ok { return 0 }

    rec.mu.Lock()
    defer rec.mu.Unlock()

    now := time.Now()
    var total float64
    for _, s := range rec.signals {
        age := now.Sub(s.At)
        // Exponential decay: weight * 2^(-age/halfLife)
        decayFactor := math.Pow(2, -float64(age)/float64(ti.halfLife))
        total += s.Weight * decayFactor
    }
    return total
}

// Action returns the recommended action based on current score
func (ti *ThreatIntel) Action(ip string) string {
    score := ti.Score(ip)
    switch {
    case score >= 80: return "block"
    case score >= 50: return "challenge"
    case score >= 25: return "monitor"
    default:          return "allow"
    }
}

func (ti *ThreatIntel) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        cutoff := time.Now().Add(-24 * time.Hour) // clean signals older than 24h
        ti.mu.Lock()
        for ip, rec := range ti.records {
            rec.mu.Lock()
            var fresh []Signal
            for _, s := range rec.signals {
                if s.At.After(cutoff) { fresh = append(fresh, s) }
            }
            rec.signals = fresh
            if len(fresh) == 0 { delete(ti.records, ip) }
            rec.mu.Unlock()
        }
        ti.mu.Unlock()
    }
}

// Integration: ThreatIntel middleware
func (ti *ThreatIntel) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := realClientIP(r)
        action := ti.Action(ip)
        switch action {
        case "block":
            http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
            return
        case "challenge":
            // In production: return a PoW challenge (Day 97)
            // For now: add delay
            time.Sleep(200 * time.Millisecond)
        }
        next.ServeHTTP(w, r)
    })
}
```

### Real-World Context
A single WAF hit might be a false positive. But 5 WAF hits + 10 rate limit violations + 3 auth failures in 10 minutes = almost certainly an attacker. The aggregated score with time decay is how Cloudflare's threat scoring works. The 10-minute half-life means an IP that cleans up its act can recover — preventing permanent false positive lockouts.

---

## Day 84 — Structured Logging with slog

### Learning Objectives
- Migrate all `fmt.Printf` to `log/slog` with JSON output
- Propagate correlation IDs through context
- Configure log levels for production (suppress DEBUG)

### Exercise: Complete slog Integration

```go
package logging

import (
    "context"
    "log/slog"
    "os"
)

type ctxKey string
const loggerKey ctxKey = "logger"

// NewProductionLogger creates a JSON logger for production
func NewProductionLogger(level slog.Level) *slog.Logger {
    opts := &slog.HandlerOptions{
        Level: level,
        // Add source location to all log records
        AddSource: level <= slog.LevelDebug,
    }
    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// WithRequestContext injects request-specific fields into the logger
func WithRequestContext(ctx context.Context, logger *slog.Logger, requestID, clientIP, method, path string) (context.Context, *slog.Logger) {
    enriched := logger.With(
        "request_id", requestID,
        "client_ip",  clientIP,
        "method",     method,
        "path",       path,
    )
    ctx = context.WithValue(ctx, loggerKey, enriched)
    return ctx, enriched
}

// FromContext retrieves the request-scoped logger
func FromContext(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return l
    }
    return slog.Default()
}

// LoggingMiddleware injects a request-scoped logger
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := GetRequestID(r.Context())
            ctx, reqLogger := WithRequestContext(r.Context(), logger,
                requestID, realClientIP(r), r.Method, r.URL.Path)

            start := time.Now()
            capture := &statusCapture{ResponseWriter: w, statusCode: 200}
            next.ServeHTTP(capture, r.WithContext(ctx))

            reqLogger.Info("request_complete",
                "status",     capture.statusCode,
                "latency_ms", time.Since(start).Milliseconds(),
                "bytes",      capture.bytesWritten,
            )
        })
    }
}
```

**Usage in handlers:**
```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    logger := logging.FromContext(r.Context())
    // logger already has request_id, client_ip, method, path
    logger.Info("processing request", "user_id", getUserID(r))
    logger.Debug("cache lookup", "cache_key", "user:42", "hit", true)
    logger.Error("upstream error", "upstream", "db", "error", err)
}
```

**Sample JSON output:**
```json
{"time":"2025-01-10T10:00:00Z","level":"INFO","msg":"request_complete","request_id":"a3f8c2d1","client_ip":"10.0.0.1","method":"POST","path":"/api/auth/login","status":200,"latency_ms":47,"bytes":256}
```

### Real-World Context
Splunk, Datadog, and Elastic all parse structured JSON logs automatically. The correlation `request_id` field ties together WAF logs, auth logs, and application logs for a single request — critical for forensics after a security incident. `slog` is stdlib since Go 1.21, replacing the `zerolog`/`zap` ecosystem for new services.

---

## Day 85 — Configuration & Hot Reload

### Learning Objectives
- Build a unified config system covering WAF, rate limits, CORS, and auth
- Hot-reload on `SIGHUP` with validation before swap
- Reject invalid config without disrupting service

### Exercise: Hot-Reload Config System

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type GatewayConfig struct {
    WAF struct {
        Mode      string         `json:"mode"`
        Threshold int            `json:"threshold"`
    } `json:"waf"`
    RateLimit struct {
        PerIPBurst    float64 `json:"per_ip_burst"`
        PerIPRate     float64 `json:"per_ip_rate"`
        PerUserBurst  float64 `json:"per_user_burst"`
        PerUserRate   float64 `json:"per_user_rate"`
    } `json:"rate_limit"`
    CORS struct {
        AllowedOrigins []string `json:"allowed_origins"`
    } `json:"cors"`
    TrustedCIDRs []string `json:"trusted_cidrs"`
}

func (c *GatewayConfig) Validate() error {
    switch c.WAF.Mode {
    case "off", "detection", "blocking":
    default:
        return fmt.Errorf("waf.mode must be off|detection|blocking, got %q", c.WAF.Mode)
    }
    if c.RateLimit.PerIPRate <= 0 {
        return fmt.Errorf("rate_limit.per_ip_rate must be > 0")
    }
    if len(c.CORS.AllowedOrigins) == 0 {
        return fmt.Errorf("cors.allowed_origins must not be empty")
    }
    return nil
}

type ConfigStore struct {
    mu       sync.RWMutex
    current  *GatewayConfig
    path     string
    onChange []func(*GatewayConfig)
}

func Load(path string) (*ConfigStore, error) {
    s := &ConfigStore{path: path}
    if err := s.reload(); err != nil {
        return nil, err
    }
    return s, nil
}

func (s *ConfigStore) Get() *GatewayConfig {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.current
}

func (s *ConfigStore) OnChange(fn func(*GatewayConfig)) {
    s.onChange = append(s.onChange, fn)
}

func (s *ConfigStore) reload() error {
    data, err := os.ReadFile(s.path)
    if err != nil { return fmt.Errorf("read config: %w", err) }

    var cfg GatewayConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return fmt.Errorf("parse config: %w", err)
    }

    // Validate before swap — never apply invalid config
    if err := cfg.Validate(); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }

    s.mu.Lock()
    s.current = &cfg
    s.mu.Unlock()

    // Notify all listeners (WAF, rate limiter, CORS)
    for _, fn := range s.onChange {
        go fn(&cfg)
    }

    fmt.Printf("[CONFIG] reloaded from %s (waf.mode=%s)\n", s.path, cfg.WAF.Mode)
    return nil
}

// WatchSignal reloads config on SIGHUP
func (s *ConfigStore) WatchSignal(stop <-chan struct{}) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGHUP)
    defer signal.Stop(sigCh)

    for {
        select {
        case <-stop:
            return
        case sig := <-sigCh:
            fmt.Printf("[CONFIG] received %v, reloading...\n", sig)
            if err := s.reload(); err != nil {
                fmt.Printf("[CONFIG] reload error (keeping current): %v\n", err)
                // Critically: keep current config on error
            }
        }
    }
}

// WatchFile polls for file changes (alternative to SIGHUP)
func (s *ConfigStore) WatchFile(interval time.Duration, stop <-chan struct{}) {
    var lastMod time.Time
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-stop:
            return
        case <-ticker.C:
            info, err := os.Stat(s.path)
            if err != nil { continue }
            if info.ModTime().After(lastMod) {
                lastMod = info.ModTime()
                if err := s.reload(); err != nil {
                    fmt.Printf("[CONFIG] auto-reload error: %v\n", err)
                }
            }
        }
    }
}
```

**Test hot-reload:**
```bash
# Start server
go run . -config gateway.json &

# Change WAF mode to detection
sed -i 's/"blocking"/"detection"/' gateway.json

# Send SIGHUP to reload
kill -HUP $(pgrep -f 'go run')
# Server logs: [CONFIG] reloaded from gateway.json (waf.mode=detection)

# Test: malicious request now goes through (detection mode)
curl "http://localhost:8080/search?q=' UNION SELECT * FROM users--"
# Returns 200 (logged but not blocked)
```

### Real-World Context
Emergency WAF disable without deployment: your WAF is blocking legitimate traffic from a new client. Instead of deploying code, `kill -HUP $(pidof proxy)` switches to detection mode in seconds. The validation step ensures you can't accidentally deploy an invalid config (empty CORS origins, typo in WAF mode) that breaks all CORS.

---

## Day 86 — Integration Testing the Full Chain

### Learning Objectives
- Write integration tests for the complete middleware stack
- Test security scenarios: auth bypass, WAF bypass, rate limit
- Use `httptest.NewTLSServer` for full TLS integration tests

### Exercise: Integration Test Suite

```go
// integration_test.go
package integration_test

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"
)

func setupTestGateway(t *testing.T) (*httptest.Server, *TokenService) {
    t.Helper()
    tokenService, _ := NewTokenService()

    handler := buildTestGateway(tokenService)
    ts := httptest.NewTLSServer(handler)
    t.Cleanup(ts.Close)
    return ts, tokenService
}

func TestIntegration_ValidRequest_Passes(t *testing.T) {
    ts, svc := setupTestGateway(t)
    token, _ := svc.Issue("user:42", []string{"read"})

    req, _ := http.NewRequest("GET", ts.URL+"/api/data", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := ts.Client().Do(req)
    if err != nil { t.Fatal(err) }
    if resp.StatusCode != 200 {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
}

func TestIntegration_NoAuth_Returns401(t *testing.T) {
    ts, _ := setupTestGateway(t)

    resp, err := ts.Client().Get(ts.URL + "/api/data")
    if err != nil { t.Fatal(err) }
    if resp.StatusCode != 401 {
        t.Errorf("expected 401, got %d", resp.StatusCode)
    }
}

func TestIntegration_SQLInjection_Blocked(t *testing.T) {
    ts, svc := setupTestGateway(t)
    token, _ := svc.Issue("user:42", []string{"read"})

    req, _ := http.NewRequest("GET", ts.URL+"/api/search?q=' UNION SELECT * FROM users--", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := ts.Client().Do(req)
    if err != nil { t.Fatal(err) }
    if resp.StatusCode != 403 {
        t.Errorf("expected 403 for SQLi, got %d", resp.StatusCode)
    }

    // Verify request ID in response (not rule details)
    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    if _, ok := body["request_id"]; !ok {
        t.Error("response should include request_id")
    }
    if _, ok := body["rule"]; ok {
        t.Error("response should NOT include rule details (information disclosure)")
    }
}

func TestIntegration_RateLimitEnforced(t *testing.T) {
    ts, svc := setupTestGateway(t)
    token, _ := svc.Issue("user:42", []string{"read"})

    client := ts.Client()
    var tooManyRequests int

    // Send 100 requests rapidly (limit is 10/s burst 50)
    for i := 0; i < 100; i++ {
        req, _ := http.NewRequest("GET", ts.URL+"/api/data", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        resp, err := client.Do(req)
        if err != nil { continue }
        if resp.StatusCode == 429 {
            tooManyRequests++
        }
    }

    if tooManyRequests == 0 {
        t.Error("expected some 429 responses, got none")
    }
    t.Logf("Got %d 429 responses out of 100 requests", tooManyRequests)
}

func TestIntegration_CORS_AllowedOrigin(t *testing.T) {
    ts, _ := setupTestGateway(t)

    req, _ := http.NewRequest("GET", ts.URL+"/api/data", nil)
    req.Header.Set("Origin", "https://app.example.com")

    resp, err := ts.Client().Do(req)
    if err != nil { t.Fatal(err) }

    acao := resp.Header.Get("Access-Control-Allow-Origin")
    if acao != "https://app.example.com" {
        t.Errorf("expected ACAO=https://app.example.com, got %q", acao)
    }
    vary := resp.Header.Get("Vary")
    if !strings.Contains(vary, "Origin") {
        t.Errorf("Vary header should include Origin, got %q", vary)
    }
}

func TestIntegration_CORS_DisallowedOrigin_NoHeader(t *testing.T) {
    ts, _ := setupTestGateway(t)

    req, _ := http.NewRequest("GET", ts.URL+"/api/data", nil)
    req.Header.Set("Origin", "https://evil.com")

    resp, err := ts.Client().Do(req)
    if err != nil { t.Fatal(err) }

    if acao := resp.Header.Get("Access-Control-Allow-Origin"); acao != "" {
        t.Errorf("evil.com should not get ACAO header, got %q", acao)
    }
}

func TestIntegration_OversizedBody_Rejected(t *testing.T) {
    ts, svc := setupTestGateway(t)
    token, _ := svc.Issue("user:42", []string{"write"})

    // 2MB body (above 1MB WAF scan limit)
    body := strings.NewReader(strings.Repeat("x", 2*1024*1024))
    req, _ := http.NewRequest("POST", ts.URL+"/api/data", body)
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := ts.Client().Do(req)
    if err != nil { t.Fatal(err) }
    if resp.StatusCode != 413 {
        t.Errorf("expected 413 for oversized body, got %d", resp.StatusCode)
    }
}
```

### Real-World Context
Integration tests catch middleware ordering bugs that unit tests miss. "Auth runs before WAF" is only verifiable if you have a test that sends a malicious request WITHOUT auth — unit tests of individual middleware won't catch this. Run integration tests in CI with `-race` flag.

---

## Day 87 — Fuzz Testing

### Learning Objectives
- Write fuzz tests for WAF normalization and JSON schema validator
- Find crashes and panics with `go test -fuzz`
- Fix all panics before deploying security-critical code

### Exercise: Fuzz the WAF

```go
// fuzz_test.go
package waf_test

import (
    "testing"
    "unicode/utf8"
)

// FuzzNormalizationPipeline finds inputs that panic the normalizer
func FuzzNormalizationPipeline(f *testing.F) {
    // Seed corpus — known interesting inputs
    seeds := []string{
        "",
        "\x00",
        "%",
        "%%",
        "%2e%2e%2f",
        "%252e%252e%252f",
        "' OR '1'='1",
        "<script>alert(1)</script>",
        strings.Repeat("A", 10000),  // very long input
        "\xff\xfe",                   // invalid UTF-8
        "&#x0000;",                   // null HTML entity
        "javascript:%0a%0dalert(1)", // newline-encoded JS
    }
    for _, s := range seeds {
        f.Add(s)
    }

    f.Fuzz(func(t *testing.T, input string) {
        // Must not panic on any input
        result := NormalizationPipeline(input)

        // Must return valid UTF-8 (so downstream string operations are safe)
        if !utf8.ValidString(result) {
            t.Errorf("normalization returned invalid UTF-8 for input %q", input)
        }

        // Must not return longer than input (normalization should not expand)
        if len(result) > len(input)*10 {
            t.Errorf("normalization expanded input suspiciously: %d → %d", len(input), len(result))
        }
    })
}

// FuzzJSONSchemaValidator finds inputs that panic the schema validator
func FuzzJSONSchemaValidator(f *testing.F) {
    seeds := [][]byte{
        []byte(`{}`),
        []byte(`{"username":"alice","email":"a@b.com","role":"admin","password":"secret123"}`),
        []byte(`null`),
        []byte(`[]`),
        []byte(`{{{`),                   // malformed JSON
        []byte(`{"x":` + strings.Repeat("[", 1000)),  // deeply nested
    }
    for _, s := range seeds {
        f.Add(s)
    }

    schema := &Schema{
        Type:     "object",
        Required: []string{"username", "email"},
        Properties: map[string]*Schema{
            "username": {Type: "string", MinLength: ptr(3), MaxLength: ptr(64)},
            "email":    {Type: "string"},
        },
    }

    f.Fuzz(func(t *testing.T, body []byte) {
        // Must not panic
        var data interface{}
        json.Unmarshal(body, &data)
        _ = schema.Validate(data, "$")
    })
}
```

**Run the fuzzer:**
```bash
# Fuzz for 30 seconds
go test -fuzz=FuzzNormalizationPipeline -fuzztime=30s ./...
go test -fuzz=FuzzJSONSchemaValidator -fuzztime=30s ./...

# If a crash is found, it's saved to testdata/fuzz/
# Run to reproduce:
go test -run=FuzzNormalizationPipeline/testdata/fuzz/...
```

### Real-World Context
Go's built-in fuzzer finds bugs in security code that hand-written tests miss. In 2022, fuzzing found ~30 vulnerabilities in the Go standard library. For security-critical code (WAF, parsers, crypto), fuzzing for 10 minutes before each release is best practice. The `testdata/fuzz/` directory stores crash-inducing inputs as regression tests.

---

## Day 88 — Load Testing the Security Stack

### Learning Objectives
- Measure the latency overhead of each security middleware
- Identify the bottleneck (usually WAF regex or rate limiter mutex)
- Generate the numbers needed to set SLOs

### Exercise: Security Middleware Overhead Measurement

```go
package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "sort"
    "sync"
    "time"
)

type OverheadMeasurement struct {
    Name      string
    Samples   []time.Duration
}

func (m *OverheadMeasurement) P99() time.Duration {
    if len(m.Samples) == 0 { return 0 }
    sorted := make([]time.Duration, len(m.Samples))
    copy(sorted, m.Samples)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
    return sorted[int(float64(len(sorted))*0.99)]
}

func measureOverhead(name string, handler http.Handler, requests int) *OverheadMeasurement {
    m := &OverheadMeasurement{Name: name}
    var mu sync.Mutex

    var wg sync.WaitGroup
    concurrency := 50
    perWorker := requests / concurrency

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < perWorker; j++ {
                w := httptest.NewRecorder()
                r := httptest.NewRequest("GET", "/api/data?q=legitimate+search+term", nil)
                r.Header.Set("Authorization", "Bearer valid-token")

                start := time.Now()
                handler.ServeHTTP(w, r)
                elapsed := time.Since(start)

                mu.Lock()
                m.Samples = append(m.Samples, elapsed)
                mu.Unlock()
            }
        }()
    }
    wg.Wait()
    return m
}

func main() {
    baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
    })

    configs := []struct{ name string; handler http.Handler }{
        {"baseline (no middleware)", baseHandler},
        {"+ RequestID", Chain(baseHandler, RequestID)},
        {"+ SecurityHeaders", Chain(baseHandler, RequestID, SecurityHeaders)},
        {"+ Auth", Chain(baseHandler, RequestID, SecurityHeaders, ecdsaAuth)},
        {"+ Rate Limit", Chain(baseHandler, RequestID, SecurityHeaders, ecdsaAuth, rateLimiter)},
        {"+ WAF", Chain(baseHandler, RequestID, SecurityHeaders, ecdsaAuth, rateLimiter, waf)},
        {"full stack", fullStack},
    }

    fmt.Println("=== Security Middleware Overhead ===")
    fmt.Printf("%-35s %8s %8s\n", "Stack", "P50 (µs)", "P99 (µs)")
    fmt.Println(strings.Repeat("-", 55))

    for _, cfg := range configs {
        m := measureOverhead(cfg.name, cfg.handler, 10000)
        samples := m.Samples
        sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
        p50 := samples[len(samples)/2]
        p99 := m.P99()
        fmt.Printf("%-35s %8.0f %8.0f\n",
            cfg.name,
            float64(p50.Microseconds()),
            float64(p99.Microseconds()),
        )
    }
}
```

**Expected output on developer machine:**
```
Stack                               P50 (µs)  P99 (µs)
-------------------------------------------------------
baseline (no middleware)                   8        45
+ RequestID                               12        55
+ SecurityHeaders                         15        60
+ Auth                                    28        85
+ Rate Limit                              35       100
+ WAF                                     95       280
full stack                               110       320
```

**WAF adds ~85µs p50, ~220µs p99.** At 1000 rps: 85ms CPU/second for WAF scanning. At 10000 rps: 850ms CPU/second. This is the data you need to justify WAF hardware or horizontal scaling.

---

## Day 89 — Goroutine Leak Audit

### Learning Objectives
- Find goroutine leaks in the complete security proxy
- Verify goroutines return to baseline after load test ends
- Fix common leak patterns: context cancellation, channel drain

### Exercise: Goroutine Baseline Checker

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

// LeakChecker verifies goroutine count returns to baseline
type LeakChecker struct {
    baseline int
    label    string
}

func NewLeakChecker(label string) *LeakChecker {
    // Wait for GC and goroutine scheduler to settle
    runtime.GC()
    time.Sleep(100 * time.Millisecond)
    return &LeakChecker{
        baseline: runtime.NumGoroutine(),
        label:    label,
    }
}

func (lc *LeakChecker) Check(t interface{ Errorf(string, ...interface{}) }, timeout time.Duration) {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        current := runtime.NumGoroutine()
        if current <= lc.baseline+5 { // allow 5 goroutines overhead
            return // OK
        }
        time.Sleep(100 * time.Millisecond)
    }
    current := runtime.NumGoroutine()
    t.Errorf("[%s] goroutine leak: baseline=%d current=%d delta=%d",
        lc.label, lc.baseline, current, current-lc.baseline)

    // Dump goroutine stack for debugging
    buf := make([]byte, 64*1024)
    n := runtime.Stack(buf, true)
    fmt.Printf("Goroutine dump:\n%s\n", buf[:n])
}

// Test goroutine leak in the full gateway
func TestGateway_NoGoroutineLeak(t *testing.T) {
    ts := setupTestGateway(t)
    checker := NewLeakChecker("gateway load test")

    // Run 1000 concurrent requests
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 10; j++ {
                resp, _ := ts.Client().Get(ts.URL + "/api/data")
                if resp != nil { resp.Body.Close() }
            }
        }()
    }
    wg.Wait()

    // After load: goroutine count should return to baseline within 5s
    checker.Check(t, 5*time.Second)
}
```

**Common goroutine leaks to fix:**

```go
// LEAK: Ticker not stopped
ticker := time.NewTicker(time.Second)
go func() {
    for range ticker.C { doWork() }  // goroutine stuck if never stopped
}()

// FIX:
ticker := time.NewTicker(time.Second)
go func() {
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C: doWork()
        case <-ctx.Done(): return  // exits when context is cancelled
        }
    }
}()

// LEAK: Goroutine waiting on channel with no sender
results := make(chan string)
go func() {
    r := <-results  // blocks forever if no one sends
    fmt.Println(r)
}()
// Caller never sends to results — goroutine leaks

// FIX:
results := make(chan string, 1)  // buffered, or sender always sends before returning
```

### Real-World Context
Goroutine leaks cause memory growth without obvious cause — gradual OOM. At 1000 rps with a 1-goroutine-per-request leak: +1000 goroutines/second, ~2KB each = +2MB/second memory growth = OOM in minutes. The `runtime.Stack` dump shows exactly which goroutine is stuck and on what.

---

## Day 90 — Phase 3 Review & Architecture Diagram

### Review Checklist

```
Phase 3 Code Audit:
[ ] All HTTP handlers have timeout context
[ ] WAF normalization covers URL + HTML + double-encoding
[ ] Rate limiters clean up expired buckets (no memory leak)
[ ] Auth failures return constant-time responses
[ ] CORS never reflects origin without allowlist check
[ ] PII masking covers all sensitive field names
[ ] Security headers set on all responses (including 4xx/5xx)

Performance Numbers (from Day 88):
[ ] Baseline latency measured
[ ] Per-middleware overhead documented
[ ] WAF overhead acceptable for SLO
[ ] Rate limiter doesn't serialize all requests (use sync.Map not single Mutex)

Tests:
[ ] go test -race ./... passes
[ ] Integration tests cover: valid req, sqli, xss, no auth, rate limit, cors
[ ] Fuzz tests run for 30s minimum, no crashes

Architecture:
[ ] ASCII diagram drawn
[ ] Each component mapped to OWASP control
[ ] Vendor decisions documented: what you'd buy vs build
```

### ASCII Architecture Diagram

```
Internet
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Unified Security Proxy                       │
│                                                                 │
│  L4 Layer (TCP/TLS):                                           │
│    mTLS Termination ──── SNI Routing ──── L4 Rate Limit        │
│    Circuit Breaker ───── Connection Pool ── Health Check        │
│                                                                 │
│  L7 Layer (HTTP):                                               │
│    ┌─────────────────────────────────────────────────────────┐ │
│    │ Middleware Chain (in order)                              │ │
│    │                                                          │ │
│    │  PanicRecoverer → RequestID → CORS                      │ │
│    │     → IP Rate Limit → Bot Check                         │ │
│    │     → Auth (ECDSA Token / API Key / HMAC)               │ │
│    │     → User Rate Limit → ThreatIntel                     │ │
│    │     → WAF (Normalize → SQLi/XSS/Traversal/CmdInj)      │ │
│    │     → Schema Validation → Input Validation              │ │
│    │     → PII Scrubber → Security Headers                   │ │
│    │     → Request Logger → Application Handler              │ │
│    └─────────────────────────────────────────────────────────┘ │
│                                                                 │
│  Observability:                                                 │
│    slog JSON → SIEM │ expvar → Prometheus │ pprof → Profiler   │
│    /health  /stats  /debug/vars  /debug/pprof                  │
│                                                                 │
│  Operations:                                                    │
│    SIGHUP → hot-reload config (WAF mode, rate limits, CORS)    │
│    SIGTERM → 3-stage graceful drain                            │
│    SIGUSR2 → goroutine stack dump                              │
└─────────────────────────────────────────────────────────────────┘
    │              │              │
    ▼              ▼              ▼
Backend A      Backend B      Backend C
```

### Build vs Buy Decision Framework

| Component | Build | Buy/Use | Rationale |
|-----------|-------|---------|-----------|
| TLS termination | Done | Cloud LB | Cloud LB cheaper at scale |
| mTLS service mesh | Done | Istio/Linkerd | Complex at K8s scale |
| WAF rules database | Core done | CloudFlare/AWS WAF rules | Rules maintenance is expensive |
| DDoS absorption | Build mitigation | Cloudflare/AWS Shield | Physical network needed for Tbps |
| Cert management | Build rotation | Let's Encrypt + cert-manager | Automation at scale |
| Rate limiting | Done | Redis-backed Envoy | Multi-instance needs shared state |

---

## Week 12 Mini Project: Unified Security Proxy (L4 + L7)

### Complete Integration

```go
// main.go — all components wired together
func main() {
    // Config
    cfg, _ := config.Load("gateway.json")

    // PKI
    reloader, _ := cert.NewReloader("server.crt", "server.key")

    // Security components
    tokenSvc, _ := auth.NewTokenService()
    threatIntel := threatintel.New(10 * time.Minute)
    waf, _ := waf.NewScoringWAF(waf.DefaultRules(), "waf-config.json", audit.NewLogger())
    ipLimiter := ratelimit.NewRateLimiter(50, 10)
    cors := cors.New(buildCORSConfig(cfg.Get()))
    abuse := abuse.New()

    // Update components when config changes
    cfg.OnChange(func(c *config.GatewayConfig) {
        waf.Reload(&waf.Config{Mode: c.WAF.Mode, Threshold: c.WAF.Threshold})
    })

    // Build middleware chain
    appHandler := buildRoutes(tokenSvc)
    gatewayHandler := pipeline.Chain(appHandler,
        middleware.PanicRecoverer(logger),
        middleware.RequestID,
        cors.Middleware,
        ipLimiter.Middleware(ipKeyFn),
        abuse.CheckHoneypot,
        auth.ECDSATokenAuth(tokenSvc),
        threatIntel.Middleware,
        waf.Handler,
        middleware.SecurityHeaders,
        pii.Middleware,
        middleware.RequestLogger(logger),
    )

    // TLS server with mTLS + SNI + cert hot-reload
    tlsConfig := &tls.Config{
        GetCertificate: reloader.GetCertificate,
        MinVersion:     tls.VersionTLS12,
    }
    srv := &http.Server{
        Addr:              ":443",
        Handler:           gatewayHandler,
        TLSConfig:         tlsConfig,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      60 * time.Second,
        IdleTimeout:       120 * time.Second,
    }

    // Observability
    go http.ListenAndServe(":9090", observabilityMux())

    // Signal handling
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()
    go cfg.WatchSignal(ctx.Done())

    logger.Info("unified security proxy starting", "addr", ":443", "metrics", ":9090")
    go srv.ListenAndServeTLS("", "")

    <-ctx.Done()
    logger.Info("shutting down...")
    shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    srv.Shutdown(shutCtx)
    logger.Info("shutdown complete")
}
```

### Engineer Takeaway

You've built the architecture of Cloudflare's edge — L4 TCP termination, mTLS, SNI routing, circuit breakers, HTTP/2, WAF (OWASP Top 10), API authentication, rate limiting, and full observability — in ~3000 lines of Go standard library code.

More importantly, you now have the technical depth to:
- **Design and review** network and security code, knowing what each layer actually does
- **Choose between building and configuring** managed products (AWS WAF, Cloudflare, Akamai), because you know what they do inside
- **Debug incidents across the whole stack**: you've written every component
- **Fix false positives precisely** ("WAF rule SQLI-001 triggers on Irish names; here's the exact normalization path")
- **Set realistic SLOs**: you have benchmark numbers for each middleware layer

90 days of one hour each turned "black box" infrastructure into code you understand line by line.
