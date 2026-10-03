# Week 16: Advanced Integration + Capstone Preparation (Days 113–119)

> **Phase 4 final stretch:** Test the integrated stack against real attack scenarios, profile and optimize, security-audit your own code, build test coverage, and prepare the capstone. After this week, Day 120 is pure assembly.

**Why this week matters:** Each previous week built a component in isolation. Integration is where bugs hide — middleware ordering, race conditions under load, unexpected interactions between WAF and bot detection, config reload timing issues. This week validates the entire stack.

---

> **Background in this wiki:** Linux guide [Ch 81 (profiling)](../../os-linux/real-life-os-guide.md#chapter-81-profiling-go-on-linux-pprof-the-execution-tracer-and-perf); HTTPS guide [Ch 20 (parser disagreements: the reason to fuzz parsers)](../../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling); TCP/IP guide [Ch 32 (`pcapread`, a binary parser worth fuzzing)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-32-wireshark-reading-a-conversation). Security in Depth: [Chapter 59 (`govulncheck` reachability)](../../security/real-life-security-guide-v1.md#chapter-59-vulnerability-management-at-scale) and [Chapter 58 (SSDLC gates)](../../security/real-life-security-guide-v1.md#chapter-58-building-an-ssdlc-product-security-program).

## Day 113 — End-to-End Attack Simulation

| | |
|---|---|
| **Learning objective** | Run 8 simultaneous attack types against the fully integrated stack; verify each is detected and mitigated |
| **Key packages** | `net/http`, `sync`, `context`, `testing` |
| **Security context** | Unit tests verify individual components; integration attack tests verify the system actually defends |

### Key Concepts

**Attack playbook for integration testing:**

| Attack | Expected detection | Expected mitigation |
|--------|-------------------|---------------------|
| SQL injection | WAF SQLiRule score ≥ threshold | 403 block, WAF audit log |
| XSS | WAF XSSRule score ≥ threshold | 403 block |
| Path traversal | WAF PathTraversalRule | Immediate 403 |
| Credential stuffing | AbuseDetector per-username limit | 429 + delay |
| HTTP flood | CIDRRateLimiter token bucket | 429 Retry-After |
| Slowloris | SlowlorisDetector byte rate | Connection close |
| Bot scraping | BotScorer (velocity + datacenter IP) | PoW challenge |
| Honeypot probe | HoneypotTracker | Permanent shadow-ban |

### Exercise — Attack Simulator

```go
package integration

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

type AttackResult struct {
    Name      string
    Sent      int
    Blocked   int
    Allowed   int
    BlockRate float64
}

// AttackSimulator runs attack scenarios against a test server
type AttackSimulator struct {
    baseURL string
    client  *http.Client
}

func NewAttackSimulator(baseURL string) *AttackSimulator {
    return &AttackSimulator{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 5 * time.Second,
            CheckRedirect: func(*http.Request, []*http.Request) error {
                return http.ErrUseLastResponse // don't follow redirects
            },
        },
    }
}

// 1. SQL Injection
func (s *AttackSimulator) SQLInjection(t *testing.T, n int) AttackResult {
    payloads := []string{
        "' OR '1'='1",
        "'; DROP TABLE users; --",
        "1 UNION SELECT * FROM passwords",
        "admin'--",
        "' OR 1=1--",
    }
    return s.runHTTPAttack(t, "sqli", n, func(i int) *http.Request {
        payload := payloads[i%len(payloads)]
        req, _ := http.NewRequest("GET", s.baseURL+"/api/users?id="+payload, nil)
        return req
    }, http.StatusForbidden)
}

// 2. XSS
func (s *AttackSimulator) XSS(t *testing.T, n int) AttackResult {
    payloads := []string{
        "<script>alert(1)</script>",
        "<img src=x onerror=alert(1)>",
        "javascript:alert(1)",
        "<svg onload=alert(1)>",
    }
    return s.runHTTPAttack(t, "xss", n, func(i int) *http.Request {
        payload := payloads[i%len(payloads)]
        req, _ := http.NewRequest("GET", s.baseURL+"/search?q="+payload, nil)
        return req
    }, http.StatusForbidden)
}

// 3. Path Traversal
func (s *AttackSimulator) PathTraversal(t *testing.T, n int) AttackResult {
    payloads := []string{
        "../../../etc/passwd",
        "..%2F..%2F..%2Fetc%2Fpasswd",
        "%2e%2e/%2e%2e/etc/passwd",
        "....//....//etc/passwd",
    }
    return s.runHTTPAttack(t, "path_traversal", n, func(i int) *http.Request {
        req, _ := http.NewRequest("GET", s.baseURL+"/files/"+payloads[i%len(payloads)], nil)
        return req
    }, http.StatusForbidden)
}

// 4. Credential Stuffing
func (s *AttackSimulator) CredentialStuffing(t *testing.T, n int) AttackResult {
    return s.runHTTPAttack(t, "credential_stuffing", n, func(i int) *http.Request {
        body, _ := json.Marshal(map[string]string{
            "username": fmt.Sprintf("user%d@example.com", i%100), // cycle through 100 usernames
            "password": fmt.Sprintf("leaked-pass-%d", i),
        })
        req, _ := http.NewRequest("POST", s.baseURL+"/api/auth/login", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        return req
    }, http.StatusTooManyRequests)
}

// 5. HTTP Flood
func (s *AttackSimulator) HTTPFlood(t *testing.T, n int, concurrency int) AttackResult {
    var blocked, allowed atomic.Int64
    sem := make(chan struct{}, concurrency)
    var wg sync.WaitGroup

    for i := 0; i < n; i++ {
        wg.Add(1)
        sem <- struct{}{}
        go func(i int) {
            defer wg.Done()
            defer func() { <-sem }()
            req, _ := http.NewRequest("GET", s.baseURL+"/api/data", nil)
            resp, err := s.client.Do(req)
            if err != nil {
                blocked.Add(1)
                return
            }
            resp.Body.Close()
            if resp.StatusCode == http.StatusTooManyRequests {
                blocked.Add(1)
            } else {
                allowed.Add(1)
            }
        }(i)
    }
    wg.Wait()

    b, a := int(blocked.Load()), int(allowed.Load())
    return AttackResult{
        Name: "http_flood", Sent: n, Blocked: b, Allowed: a,
        BlockRate: float64(b) / float64(n),
    }
}

// 6. Bot Scraping (high velocity, datacenter UA)
func (s *AttackSimulator) BotScraping(t *testing.T, n int) AttackResult {
    return s.runHTTPAttack(t, "bot_scraping", n, func(i int) *http.Request {
        req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/products?page=%d", s.baseURL, i), nil)
        req.Header.Set("User-Agent", "python-requests/2.28.0") // known bot UA
        // Deliberately missing Accept-Language, Accept-Encoding
        return req
    }, http.StatusTooManyRequests) // expect PoW challenge or block
}

// 7. Honeypot Probe
func (s *AttackSimulator) HoneypotProbe(t *testing.T) AttackResult {
    req, _ := http.NewRequest("GET", s.baseURL+"/api/hidden", nil)
    resp, err := s.client.Do(req)
    result := AttackResult{Name: "honeypot_probe", Sent: 1}
    if err == nil {
        resp.Body.Close()
        // Honeypot returns 200 (shadow-ban) — subsequent requests should be shadow-banned
        result.Blocked = 1
        result.BlockRate = 1.0
    }

    // Verify subsequent request from same IP gets shadow-banned
    req2, _ := http.NewRequest("GET", s.baseURL+"/api/products", nil)
    resp2, err2 := s.client.Do(req2)
    if err2 == nil {
        defer resp2.Body.Close()
        var data map[string]interface{}
        json.NewDecoder(resp2.Body).Decode(&data)
        if results, ok := data["results"]; ok {
            if arr, ok := results.([]interface{}); ok && len(arr) == 0 {
                t.Log("honeypot: subsequent request shadow-banned (empty results)")
            }
        }
    }
    return result
}

func (s *AttackSimulator) runHTTPAttack(t *testing.T, name string, n int, makeReq func(int) *http.Request, expectedBlockStatus int) AttackResult {
    var blocked, allowed atomic.Int64
    var wg sync.WaitGroup

    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            req := makeReq(i)
            resp, err := s.client.Do(req)
            if err != nil {
                blocked.Add(1)
                return
            }
            resp.Body.Close()
            if resp.StatusCode == expectedBlockStatus || resp.StatusCode == http.StatusForbidden {
                blocked.Add(1)
            } else {
                allowed.Add(1)
            }
        }(i)
    }
    wg.Wait()

    b, a := int(blocked.Load()), int(allowed.Load())
    result := AttackResult{Name: name, Sent: n, Blocked: b, Allowed: a, BlockRate: float64(b) / float64(n)}
    t.Logf("attack=%s sent=%d blocked=%d (%.0f%%) allowed=%d", name, n, b, result.BlockRate*100, a)
    return result
}

// RunFullPlaybook runs all attacks and verifies block rates
func RunFullPlaybook(t *testing.T, serverURL string) {
    sim := NewAttackSimulator(serverURL)

    results := []AttackResult{
        sim.SQLInjection(t, 50),
        sim.XSS(t, 50),
        sim.PathTraversal(t, 50),
        sim.CredentialStuffing(t, 200),
        sim.HTTPFlood(t, 500, 50),
        sim.BotScraping(t, 100),
        sim.HoneypotProbe(t),
    }

    // Verify attack detection rates
    for _, r := range results {
        if r.BlockRate < 0.8 {
            t.Errorf("FAIL: %s block rate %.0f%% < 80%% (%d/%d blocked)",
                r.Name, r.BlockRate*100, r.Blocked, r.Sent)
        } else {
            t.Logf("OK: %s block rate %.0f%%", r.Name, r.BlockRate*100)
        }
    }
}
```

---

## Day 114 — Performance Profiling the Full Stack

| | |
|---|---|
| **Learning objective** | Measure per-middleware latency contribution; optimize the two hottest paths |
| **Key packages** | `runtime/pprof`, `net/http/pprof`, `testing`, `time` |
| **Security context** | "Our WAF adds 2ms p99" is a concrete SLA number — you can only defend it if you've measured it |

### Key Concepts

**Expected overhead for a well-implemented stack:**

| Middleware | p50 | p95 | p99 |
|-----------|-----|-----|-----|
| Rate limiter (token bucket) | 1µs | 5µs | 15µs |
| IP classifier | 2µs | 8µs | 20µs |
| Bot scorer (header check) | 5µs | 15µs | 40µs |
| WAF (full normalization + rules) | 80µs | 150µs | 300µs |
| Auth (ECDSA token verify) | 10µs | 25µs | 50µs |
| Schema validator | 15µs | 40µs | 100µs |
| Full chain | ~115µs | ~250µs | ~500µs |

**The WAF dominates** because regex matching + multi-layer normalization is CPU-intensive. Optimization targets: pre-compile all regex (done), use sync.Pool for normalization buffers, short-circuit on first definitive block.

### Exercise — Per-Middleware Latency Profiler

```go
package profiling

import (
    "fmt"
    "net/http"
    "runtime/pprof"
    "sort"
    "sync"
    "sync/atomic"
    "time"
)

// MiddlewareTimer wraps a handler and records precise latency per layer
type MiddlewareTimer struct {
    name    string
    mu      sync.Mutex
    samples []int64 // nanoseconds
}

func NewMiddlewareTimer(name string) *MiddlewareTimer {
    return &MiddlewareTimer{name: name, samples: make([]int64, 0, 10000)}
}

func (mt *MiddlewareTimer) Wrap(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        elapsed := time.Since(start).Nanoseconds()
        mt.mu.Lock()
        mt.samples = append(mt.samples, elapsed)
        mt.mu.Unlock()
    })
}

// Percentile computes the Nth percentile of recorded samples
func (mt *MiddlewareTimer) Percentile(n float64) time.Duration {
    mt.mu.Lock()
    s := make([]int64, len(mt.samples))
    copy(s, mt.samples)
    mt.mu.Unlock()

    if len(s) == 0 {
        return 0
    }
    sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
    idx := int(float64(len(s)-1) * n / 100)
    return time.Duration(s[idx])
}

func (mt *MiddlewareTimer) Report() string {
    return fmt.Sprintf("%-30s p50=%6s p95=%6s p99=%6s (n=%d)",
        mt.name,
        mt.Percentile(50),
        mt.Percentile(95),
        mt.Percentile(99),
        len(mt.samples),
    )
}

// LoadGenerator creates concurrency goroutines each sending requests
type LoadGenerator struct {
    baseURL     string
    concurrency int
    total       int
    client      *http.Client
    requests    atomic.Int64
    errors      atomic.Int64
}

func NewLoadGenerator(baseURL string, concurrency, total int) *LoadGenerator {
    return &LoadGenerator{
        baseURL:     baseURL,
        concurrency: concurrency,
        total:       total,
        client:      &http.Client{Timeout: 10 * time.Second},
    }
}

func (lg *LoadGenerator) Run() (rps float64, errRate float64) {
    sem := make(chan struct{}, lg.concurrency)
    var wg sync.WaitGroup
    start := time.Now()

    for i := 0; i < lg.total; i++ {
        wg.Add(1)
        sem <- struct{}{}
        go func(i int) {
            defer wg.Done()
            defer func() { <-sem }()

            // Mix of request types: 80% valid, 10% WAF triggers, 5% rate limit, 5% auth fail
            var url string
            switch i % 20 {
            case 0, 1: // WAF trigger (10%)
                url = lg.baseURL + "/api/data?q=%27+OR+1%3D1--"
            case 2: // Rate limit (5%)
                for j := 0; j < 20; j++ { // burst
                    lg.client.Get(lg.baseURL + "/api/data")
                }
                return
            default: // Valid (85%)
                url = lg.baseURL + "/api/data"
            }

            resp, err := lg.client.Get(url)
            if err != nil {
                lg.errors.Add(1)
                return
            }
            resp.Body.Close()
            lg.requests.Add(1)
        }(i)
    }
    wg.Wait()
    elapsed := time.Since(start)

    total := float64(lg.requests.Load())
    errs := float64(lg.errors.Load())
    return total / elapsed.Seconds(), errs / (total + errs)
}

// CaptureProfile captures a CPU profile during load test
func CaptureProfile(duration time.Duration, outputPath string, fn func()) error {
    f, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        return err
    }
    defer pprof.StopCPUProfile()

    done := make(chan struct{})
    go func() {
        fn()
        close(done)
    }()

    select {
    case <-done:
    case <-time.After(duration):
    }
    return nil
}

// Optimization targets based on profiling results:
//
// 1. WAF normalization — reduce allocations with sync.Pool
//    Before: 12 allocs/op, 3.2KB/op
//    After:  3 allocs/op, 512B/op  (reuse normalization buffer)
//
// 2. Regex matching — use sync.Once + package-level vars (already done in Week 10)
//    Before: regexp.MustCompile() on each request
//    After:  pre-compiled at init time
//
// 3. IP lookup — replace sync.Map with sharded map (reduce contention)
//    Before: sync.Map single lock
//    After:  16-shard map, shard by IP[last byte] & 0x0f
```

---

## Day 115 — Security Audit of Your Own Code

| | |
|---|---|
| **Learning objective** | Systematically audit all security-critical code paths for common vulnerabilities |
| **Key packages** | `crypto/subtle`, `net/http` |
| **Security context** | Adversarial review of your own code is the last line of defense before deployment |

### Audit Checklist

**1. Timing attacks — all secret comparisons must use `subtle.ConstantTimeCompare`**

```go
// BAD — timing oracle: loop exits early on first mismatch
if apiKey == expected {
    // ...
}

// BAD — string comparison is not constant-time in Go
if hmac.Equal([]byte(provided), []byte(expected)) {} // hmac.Equal IS constant-time — OK

// GOOD — always takes same time regardless of where mismatch occurs
import "crypto/subtle"
if subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
    // reject
}

// GOOD — for HMAC comparison
if !hmac.Equal(computedMAC, expectedMAC) {
    // reject
}
```

**2. Integer overflow in packet parsers**

```go
// BAD — attacker sends length=4294967295, allocation panics or wraps
func parsePacket(data []byte) []byte {
    length := binary.BigEndian.Uint32(data[:4])
    payload := make([]byte, length) // panic if length > maxInt
    copy(payload, data[4:4+length]) // out of bounds if length > len(data)-4
    return payload
}

// GOOD — bounds check before allocation
const maxPacketSize = 1 << 20 // 1MB

func parsePacket(data []byte) ([]byte, error) {
    if len(data) < 4 {
        return nil, errors.New("truncated header")
    }
    length := binary.BigEndian.Uint32(data[:4])
    if length > maxPacketSize {
        return nil, fmt.Errorf("packet length %d exceeds maximum %d", length, maxPacketSize)
    }
    if int(length) > len(data)-4 {
        return nil, fmt.Errorf("packet claims %d bytes but only %d available", length, len(data)-4)
    }
    payload := make([]byte, length)
    copy(payload, data[4:4+length])
    return payload, nil
}
```

**3. Error messages must not leak internals**

```go
// BAD — leaks internal path
http.Error(w, fmt.Sprintf("file not found: /var/app/data/%s", filename), 404)

// BAD — leaks stack trace
http.Error(w, err.Error(), 500) // err might be: "dial tcp 10.0.0.5:5432 connection refused"

// GOOD — generic message to client, detailed log internally
requestID := r.Header.Get("X-Request-ID")
logger.Error("upstream connection failed",
    "request_id", requestID,
    "upstream", upstreamAddr,
    "error", err,
)
http.Error(w, fmt.Sprintf(`{"error":"internal_error","request_id":%q}`, requestID), 500)
```

**4. Goroutine exit conditions**

```go
// AUDIT CHECKLIST for every goroutine in your codebase:
// [ ] Does it listen on a channel? → Is there a corresponding sender or close?
// [ ] Does it poll in a for loop? → Is there a done/ctx.Done case?
// [ ] Is it started per-connection? → Is it cleaned up when connection closes?
// [ ] Is it a background worker? → Does it respond to graceful shutdown?

// Pattern to audit:
grep -n "^.*go func" **/*.go  // find every goroutine start
// For each: trace to the goroutine body and verify the exit condition
```

### Exercise — Security Audit Report Template

```markdown
## Security Audit: [Component Name]

**Date:** YYYY-MM-DD
**Auditor:** [Your name]
**Files reviewed:** [list of files]

### Finding 1: [Title]
- **Severity:** Critical / High / Medium / Low
- **Location:** file.go:line
- **Finding:** Description of the issue
- **Evidence:** Code snippet showing the vulnerability
- **Fix:** Code snippet showing the fix
- **Status:** Fixed / Accepted Risk / Won't Fix

### Summary
| Severity | Found | Fixed | Accepted |
|----------|-------|-------|---------|
| Critical | 0 | 0 | 0 |
| High     | 1 | 1 | 0 |
| Medium   | 3 | 3 | 0 |
| Low      | 2 | 1 | 1 |
```

---

## Day 116 — Test Coverage Audit

| | |
|---|---|
| **Learning objective** | Achieve ≥80% coverage on all security-critical files; focus on negative paths |
| **Key packages** | `testing`, `go test -cover`, `go tool cover` |
| **Security context** | Uncovered code in WAF, auth, and rate limiting is unverified security behavior |

### Exercise — Coverage Audit + Gap-Filling

```bash
# Run coverage for all packages
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | sort -k3 -rn | head -20

# Generate HTML report to see exact uncovered lines
go tool cover -html=coverage.out -o coverage.html

# Focus on security-critical packages
go test -coverprofile=waf.out ./internal/waf/...
go tool cover -func=waf.out
```

**Coverage targets:**

| Package | Target | Why |
|---------|--------|-----|
| `internal/waf` | ≥85% | Every undetected attack vector |
| `internal/auth` | ≥90% | Every auth bypass path |
| `internal/ratelimit` | ≥85% | Every limit boundary |
| `internal/bot` | ≥80% | Bot evasion paths |
| `internal/ddos` | ≥80% | DDoS evasion paths |

**Negative path tests to add (examples):**

```go
// WAF tests — negative paths that must be covered
func TestWAF_SQLiBypassAttempts(t *testing.T) {
    cases := []struct {
        name    string
        input   string
        blocked bool
    }{
        {"url_encoded_sqli", "%27+OR+%271%27%3D%271", true},
        {"double_encoded_sqli", "%2527+OR+1%3D1", true},
        {"comment_obfuscation", "1/**/OR/**/1=1", true},
        {"case_variation", "SeLeCt * FrOm users", true},
        {"normal_search", "SELECT the best laptop", false}, // legitimate
        {"sql_in_product_name", "SQL Server 2019 License", false}, // legitimate
    }
    // ...
}

// Auth tests — negative paths
func TestAuth_RejectExpiredToken(t *testing.T) { /* ... */ }
func TestAuth_RejectTamperedSignature(t *testing.T) { /* ... */ }
func TestAuth_RejectReplayedToken(t *testing.T) { /* ... */ }
func TestAuth_TimingConsistency(t *testing.T) {
    // Auth rejection should take same time as acceptance (within 10ms)
    validStart := time.Now()
    authMiddleware.ServeHTTP(rec, validReq)
    validTime := time.Since(validStart)

    invalidStart := time.Now()
    authMiddleware.ServeHTTP(rec2, invalidReq)
    invalidTime := time.Since(invalidStart)

    delta := validTime - invalidTime
    if delta < 0 { delta = -delta }
    if delta > 10*time.Millisecond {
        t.Errorf("auth timing difference %v > 10ms — potential timing oracle", delta)
    }
}
```

---

## Day 117 — Fuzz Testing Critical Parsers

| | |
|---|---|
| **Learning objective** | Write production-grade fuzz tests for WAF normalization, JA3 parser, and binary protocol parser |
| **Key packages** | `testing/fuzz` |
| **Security context** | A panic in your WAF is a denial-of-service — fuzz testing is the most effective way to find panic paths |

### Key Concepts

**What to fuzz:**
- Any code that **parses input** of arbitrary length/content
- Any code with **complex branching** on input content (WAF rules, protocol parsers)
- Any code where **wrong behavior has security impact** (auth, schema validator)

### Exercise — Fuzz Tests

```go
package fuzz_test

import (
    "testing"
    // Import your packages
)

// FuzzNormalizationPipeline — does the WAF normalizer panic on arbitrary input?
func FuzzNormalizationPipeline(f *testing.F) {
    // Seed corpus — known interesting inputs
    seeds := []string{
        "",
        "normal input",
        "%27 OR 1=1--",
        "<script>alert(1)</script>",
        "../../../../etc/passwd",
        string([]byte{0x00, 0x01, 0x02}), // null bytes
        string([]byte{0xff, 0xfe}),         // invalid UTF-8
        "%gg%hh",                            // invalid percent encoding
        "%%2527",                            // double-double encoded
        "a" + string(make([]byte, 10000)),  // very long input
        "a\r\n\r\nGET / HTTP/1.1",          // CRLF injection
    }
    for _, s := range seeds {
        f.Add(s)
    }

    f.Fuzz(func(t *testing.T, input string) {
        // Must not panic on any input
        result := normalize(input)
        // Result must be a valid string (non-nil, valid UTF-8 is not required)
        _ = result
        // Normalization must be idempotent: normalize(normalize(x)) == normalize(x)
        if normalize(result) != result {
            t.Errorf("normalization not idempotent for input %q: first=%q second=%q",
                input, result, normalize(result))
        }
    })
}

// FuzzJSONSchemaValidator — does malformed JSON panic the schema validator?
func FuzzJSONSchemaValidator(f *testing.F) {
    seeds := []string{
        `{}`,
        `{"name":"test"}`,
        `{"name":null}`,
        `{"name":true}`,
        `[1,2,3]`,
        `{"a":{"b":{"c":{"d":{}}}}}`, // deep nesting
        `{"\u0000":"value"}`,           // null byte in key
        `{"key":"` + string(make([]byte, 100000)) + `"}`, // huge value
        `{`,          // truncated
        `{"k":`,      // truncated
        `{} extra`,   // trailing garbage
    }
    for _, s := range seeds {
        f.Add(s)
    }

    schema := &JSONSchema{
        Type:       "object",
        Required:   []string{"name"},
        Properties: map[string]*JSONSchema{"name": {Type: "string", MaxLength: 255}},
        AdditionalProperties: boolPtr(false),
    }

    f.Fuzz(func(t *testing.T, input string) {
        // Must not panic
        errs := schema.Validate([]byte(input))
        // Just check it returns without panicking
        _ = errs
    })
}

// FuzzBinaryProtocolParser — does the binary parser panic on arbitrary bytes?
func FuzzBinaryProtocolParser(f *testing.F) {
    seeds := [][]byte{
        {},
        {0x00, 0x00, 0x00, 0x00},           // zero-length packet
        {0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'}, // valid packet
        {0xff, 0xff, 0xff, 0xff},             // max length
        {0x00, 0x00, 0x00, 0x01},             // length 1, no data
        make([]byte, 1024),                    // all zeros
    }
    for _, s := range seeds {
        f.Add(s)
    }

    f.Fuzz(func(t *testing.T, data []byte) {
        // Must not panic — parsePacket must return an error for invalid input
        payload, err := parsePacket(data)
        if err != nil {
            return // expected for most random input
        }
        // If no error: payload must be a valid slice
        if payload == nil {
            t.Error("parsePacket returned nil payload with nil error")
        }
        // Round-trip test: re-parsing the encoded payload should give same result
        // (only if your protocol has a serialization function)
    })
}

// Run fuzz tests:
// go test -fuzz=FuzzNormalizationPipeline -fuzztime=60s ./internal/waf/...
// go test -fuzz=FuzzJSONSchemaValidator -fuzztime=60s ./internal/validation/...
// go test -fuzz=FuzzBinaryProtocolParser -fuzztime=60s ./internal/protocol/...
```

---

## Day 118 — Documentation & Architecture Review

| | |
|---|---|
| **Learning objective** | Write an architecture document with ASCII diagrams, design decisions, and operational runbook |
| **Security context** | Code you can't explain to your team is a liability; documentation is the knowledge transfer mechanism |

### Full Architecture Diagram

```
                    ┌─────────────────────────────────────────────────────────────┐
                    │             CLOUD-EDGE SECURITY PROXY (go-proxy)            │
                    │                    Single binary, ~5000 LOC                  │
                    └─────────────────────────────────────────────────────────────┘

Internet             Layer                      What it does
   │
   ▼
[TCP Listener]     L4 Transport        TLS termination, mTLS client cert, SNI routing
   │                                   Cert hot-reload (SIGHUP, no drop)
   ▼
[ProtectedListener] L4 Defense        SYN flood: max 5 half-open per IP, 3s drain
   │
   ▼
[SlowlorisDetector] L4 Defense        Close conns sending < 50B/s during header phase
   │
   ▼
[JA3 Fingerprinter] Bot/TLS          Compute JA3 in GetConfigForClient; flag bad stacks
   │
   ▼
[CIDRRateLimiter]  DDoS/L4           Token bucket per /32 + per /24 subnet
   │
   ▼
[TrafficAnalyzer]  DDoS              EMA z-score: flag IPs > 3σ above baseline
   │
   ▼
[TieredResponder]  DDoS              Score < 50: allow | 50-80: PoW | >80: block
   │
   ▼
[HTTP Router]      L7                 net/http, HTTP/1.1 + HTTP/2
   │
   ▼
[SecurityHeaders]  OWASP             HSTS, CSP, X-Frame-Options, X-Content-Type
   │
   ▼
[RequestLogger]    Observability     slog JSON, trace context, request ID
   │
   ▼
[HMACAuth/ECDSA]   API Auth         Middleware chain: try HMAC → ECDSA → mTLS identity
   │
   ▼
[L7RateLimiter]    API              Per-API-key token bucket (10 req/s, burst 50)
   │
   ▼
[BotScorer]        Bot              Header consistency + IP reputation + behavior
   │
   ▼
[ChaosMiddleware]  Ops              Configurable fault injection (disabled by default)
   │
   ▼
[WAFMiddleware]    OWASP            SQLi + XSS + path traversal + SSRF + body inspect
   │
   ▼
[SchemaValidator]  API              JSON schema per route, additionalProperties:false
   │
   ▼
[PIIScrubber]      Compliance       Mask sensitive fields in logs + responses
   │
   ▼
[ReverseProxy]     L7               httputil.ReverseProxy, HTTP/2, per-backend CB
   │
   ▼
[ThreatIntel]      Aggregation      Cross-layer score with exp decay (t½ = 10min)
   │                                Feeds back into: DDoS tier, bot score, WAF bypass
   ▼
[Upstreams]
```

### Design Decision Records

```markdown
## ADR-001: Token Bucket for L7 Rate Limiting

**Context:** Need per-API-key rate limiting at L7.

**Decision:** Token bucket (capacity=50, rate=10/s).

**Alternatives considered:**
- Fixed window: boundary burst issue (user gets 2× limit at window boundary)
- Sliding window log: O(N) memory per key, N = requests in window
- Sliding window counter: O(1), slight approximation

**Rationale:** Token bucket allows legitimate burst (e.g., user refreshes page quickly)
while bounding steady-state rate. Used by AWS API Gateway and Stripe.

**Consequences:** Memory: O(1) per key. Users can burst to 50 requests then are limited to 10/s.

---

## ADR-002: Scoring WAF vs. Rule-Based WAF

**Context:** Should WAF block on first rule match or aggregate scores?

**Decision:** Scoring WAF with per-category thresholds.

**Rationale:** Rule-based WAF has high false positive rate on complex applications
(legitimate search queries contain SQL keywords). Scoring allows tuning sensitivity
per deployment. ModSecurity paranoia levels implement the same concept.

**Consequences:** Some attacks require multiple signals to block. Single-signal attacks
(path traversal) use immediate block regardless of score.

---

## ADR-003: In-Process ThreatIntel vs. External Store

**Context:** Share threat signals between L4 and L7 layers.

**Decision:** In-process sync.Map with exponential decay.

**Alternatives:** Redis (external), SQLite (embedded persistence).

**Rationale:** Single binary, zero external dependencies. Signals are ephemeral
(attack traffic is transient). Memory: ~1KB per IP × 10k IPs = 10MB max.

**Consequences:** Signals lost on restart. Not shared across multiple proxy instances.
```

### Operational Runbook

```markdown
## Runbook: Blocking an IP

1. Send signal: `kill -HUP $(pgrep go-proxy)`
2. Add to blocklist in config: `echo '{"blocked_cidrs": ["1.2.3.4/32"]}' > /etc/proxy/blocked.json`
3. Reload: `kill -HUP $(pgrep go-proxy)`
4. Verify: `curl http://localhost:8080/debug/vars | jq '.blocked_cidrs'`

## Runbook: Adding a WAF Rule Exception

1. Edit `/etc/proxy/waf-exceptions.json`
2. Add: `{"path": "/api/search", "skip_rules": ["sqli"]}`
3. Reload: `kill -HUP $(pgrep go-proxy)`
4. Monitor: `tail -f /var/log/proxy/waf-audit.log | jq 'select(.action=="allow")'`

## Runbook: Emergency Config Rollback

1. Stop the rollout: `kill -TERM $(pgrep go-proxy)` — graceful drain, 30s
2. Start old version: `systemctl start go-proxy@previous`
3. Verify health: `curl http://localhost:8080/healthz`
```

---

## Day 119 — Capstone Preparation & Final Integration

| | |
|---|---|
| **Learning objective** | Wire all components, pass all quality gates, establish performance baseline |
| **Key packages** | All stdlib |
| **Security context** | This is the pre-launch checklist — every gate must pass before Day 120 |

### Pre-Capstone Quality Gates

```bash
# 1. Build — must compile without errors
go build ./...

# 2. Vet — must pass all static analysis checks
go vet ./...

# 3. Race detector — zero data races
go test -race -timeout=2m ./...

# 4. Coverage — security packages must meet targets
go test -coverprofile=cover.out ./...
go tool cover -func=cover.out | awk '
    /internal\/waf/ { if ($3+0 < 85) print "FAIL: waf coverage " $3 }
    /internal\/auth/ { if ($3+0 < 90) print "FAIL: auth coverage " $3 }
    /internal\/ratelimit/ { if ($3+0 < 85) print "FAIL: ratelimit coverage " $3 }
'

# 5. Fuzz (quick run) — no panics on seed corpus
go test -run=FuzzNormalizationPipeline/seed -fuzztime=10s ./internal/waf/...
go test -run=FuzzJSONSchemaValidator/seed -fuzztime=10s ./internal/validation/...

# 6. Attack simulation — all attacks blocked at ≥80% rate
go test -v -run=TestFullAttackPlaybook ./integration/...

# 7. Performance baseline
# 500 concurrent clients, 10k requests, p99 overhead < 10ms
go test -run=TestPerfBaseline -v ./integration/...

# 8. Goroutine leak check
# After load test: goroutine count must return to baseline ± 5 within 10s
go test -run=TestGoroutineLeak -v ./integration/...

# 9. Memory baseline (no unbounded growth)
# Heap after 10k requests must be < heap + (n * 10KB)
go test -run=TestMemoryBaseline -v ./integration/...
```

### Component Wiring Checklist

```go
// main.go structure for Day 120 — all components wired here
package main

import (
    // Standard library only
    "context"
    "crypto/tls"
    "log/slog"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    // Internal packages
    "github.com/you/proxy/internal/auth"
    "github.com/you/proxy/internal/bot"
    "github.com/you/proxy/internal/chaos"
    "github.com/you/proxy/internal/ddos"
    "github.com/you/proxy/internal/health"
    "github.com/you/proxy/internal/ops"
    "github.com/you/proxy/internal/ratelimit"
    "github.com/you/proxy/internal/threat"
    "github.com/you/proxy/internal/tracing"
    "github.com/you/proxy/internal/validation"
    "github.com/you/proxy/internal/waf"
)

// Component checklist:
// [x] TLS config with mTLS + SNI routing (Week 6)
// [x] CertReloader (Week 8)
// [x] ProtectedListener + SlowlorisDetector (Week 13/15)
// [x] CIDRRateLimiter (Week 13)
// [x] TrafficAnalyzer + TieredResponder (Week 13)
// [x] JA3Detector (Week 14)
// [x] BotScorer + TieredBotMiddleware (Week 14)
// [x] HoneypotTracker (Week 14)
// [x] AuthMiddleware (HMAC + ECDSA + mTLS) (Week 11)
// [x] L7RateLimiter (Week 11)
// [x] WAFMiddleware (Week 10)
// [x] SchemaValidator (Week 11)
// [x] PIIScrubber (Week 11)
// [x] ThreatIntel store (Week 12)
// [x] slog JSON logger (Week 12)
// [x] Tracer (W3C trace context) (Week 15)
// [x] HealthRegistry (Week 15)
// [x] ShutdownManager (Week 15)
// [x] SignalDispatcher (Week 15)
// [x] ChaosMiddleware (Week 15)
// [x] ReverseProxy with circuit breaker per backend (Week 9, 15)
// [x] expvar + pprof endpoints (Week 4, 8)
// [x] FD inheritance for zero-downtime restart (Week 15)
```

---

## Week 16 Review

| Day | Skill | Real-World Use |
|-----|-------|----------------|
| 113 | End-to-end attack simulation | Verify the integrated stack actually defends |
| 114 | Per-middleware latency profiling + optimization | Establish SLA: "WAF adds 300µs p99" |
| 115 | Security self-audit (timing, overflow, info leak, goroutine exit) | Pre-launch adversarial review |
| 116 | Test coverage audit (≥80% on security packages) | Compliance evidence and confidence |
| 117 | Fuzz testing (WAF normalization, schema validator, protocol parser) | Find panic paths before attackers do |
| 118 | Architecture documentation + ADRs + operational runbook | Knowledge transfer and team enablement |
| 119 | Pre-capstone integration + all quality gates | Day 120 is assembly, not discovery |

### Engineer Takeaway

This week is what "engineering excellence" looks like in practice: attack simulation, profiling, self-audit, coverage, fuzzing, and documentation — all before shipping. Teams that skip this week ship faster but pay the cost in incidents. Teams that do this week ship slightly slower but run reliably for years.
