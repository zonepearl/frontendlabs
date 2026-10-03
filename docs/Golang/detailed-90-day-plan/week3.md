# Week 3: Context, I/O & JSON (Days 15–21)

**Phase:** 1 — Go Fundamentals
**Goal:** Master the three packages in literally every production Go service: `context`, `io`, and `encoding/json`. These are the building blocks of your HTTP middleware chain.

---

## Day 15 — context Package

### Learning Objectives
- Propagate cancellation and deadlines through a call stack
- Understand why context is Go's answer to request-scoped data
- Build timeout chains: parent cancels children automatically

### Key Concepts

```
Request arrives
    └── ctx, cancel = context.WithTimeout(bg, 30s)
        ├── auth check ctx, _ = context.WithTimeout(ctx, 2s)
        ├── DB query    ctx, _ = context.WithTimeout(ctx, 5s)
        └── response    context.WithValue(ctx, "request_id", "abc123")

If auth takes >2s: auth context times out → only auth cancels
If parent >30s:    ALL children cancel simultaneously
```

```go
// Creating contexts
ctx := context.Background()           // root, never cancels
ctx := context.TODO()                 // placeholder during refactoring

// Cancellation
ctx, cancel := context.WithCancel(parent)
defer cancel()  // ALWAYS defer cancel to avoid context leak

// Timeout
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// Deadline
ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
defer cancel()

// Request-scoped values (use typed keys to avoid collisions)
type ctxKey string
const requestIDKey ctxKey = "request_id"

ctx = context.WithValue(ctx, requestIDKey, "abc-123")
id := ctx.Value(requestIDKey).(string)

// Check cancellation
select {
case <-ctx.Done():
    return ctx.Err() // context.DeadlineExceeded or context.Canceled
default:
    // not cancelled, proceed
}
```

### Exercise: Multi-Stage Request Pipeline

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type ctxKey string

const (
    requestIDKey ctxKey = "request_id"
    userIDKey    ctxKey = "user_id"
)

// authenticate validates credentials — timeout: 2s
func authenticate(ctx context.Context, token string) (string, error) {
    authCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    // Simulate DB lookup
    done := make(chan string, 1)
    go func() {
        time.Sleep(50 * time.Millisecond) // fast auth
        done <- "user:42"
    }()

    select {
    case userID := <-done:
        return userID, nil
    case <-authCtx.Done():
        return "", fmt.Errorf("authenticate: %w", authCtx.Err())
    }
}

// fetchData retrieves data — timeout: 5s
func fetchData(ctx context.Context, userID string) (map[string]string, error) {
    dataCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    done := make(chan map[string]string, 1)
    go func() {
        time.Sleep(100 * time.Millisecond)
        done <- map[string]string{"name": "Alice", "role": "admin"}
    }()

    select {
    case data := <-done:
        return data, nil
    case <-dataCtx.Done():
        return nil, fmt.Errorf("fetchData: %w", dataCtx.Err())
    }
}

func handleRequest(ctx context.Context, requestID, token string) error {
    // Inject request ID into context — available to all downstream functions
    ctx = context.WithValue(ctx, requestIDKey, requestID)

    // Stage 1: auth
    userID, err := authenticate(ctx, token)
    if err != nil {
        return fmt.Errorf("[%s] auth failed: %w", requestID, err)
    }
    ctx = context.WithValue(ctx, userIDKey, userID)
    fmt.Printf("[%s] authenticated as %s\n", requestID, userID)

    // Stage 2: fetch data
    data, err := fetchData(ctx, userID)
    if err != nil {
        return fmt.Errorf("[%s] data fetch failed: %w", requestID, err)
    }
    fmt.Printf("[%s] data: %v\n", requestID, data)
    return nil
}

func main() {
    // Normal request — 30s budget, well within timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := handleRequest(ctx, "req-001", "valid-token"); err != nil {
        fmt.Println("ERROR:", err)
    }

    // Tight deadline — will fail
    tightCtx, tightCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer tightCancel()
    if err := handleRequest(tightCtx, "req-002", "valid-token"); err != nil {
        fmt.Println("TIMEOUT:", err)
    }
}
```

### Real-World Context
Every HTTP handler in your API gateway receives a `context` from `net/http`. The context is cancelled when the client disconnects. Without checking `ctx.Done()`, your handler keeps running, consuming DB connections and CPU, even though no one will receive the response. Always propagate context.

---

## Day 16 — io.Reader & io.Writer

### Learning Objectives
- Compose readers and writers to inspect traffic without buffering it all
- Use `io.TeeReader` to mirror traffic — the foundation of your TCP traffic inspector (Week 5)
- Use `io.LimitReader` to enforce body size limits — protection against request body attacks

### Key Concepts

```
io.Reader interface:
    Read(p []byte) (n int, err error)

io.Writer interface:
    Write(p []byte) (n int, err error)

Implementations: net.Conn, os.File, bytes.Buffer, strings.Reader,
                 http.Request.Body, gzip.Reader, tls.Conn

Composing:
    io.TeeReader(r, w)    — reads from r, copies every byte to w
    io.LimitReader(r, n)  — reads at most n bytes, then returns io.EOF
    io.MultiWriter(w1,w2) — writes to all writers simultaneously
    io.Copy(dst, src)     — copy all bytes; returns (n, err)
    bufio.NewReader(r)    — adds buffering; enables ReadLine, ReadString
```

### Exercise: Request Body Inspector (WAF Foundation)

This is the core of your WAF body scanner (Day 72):

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

const maxBodySize = 10 * 1024 // 10KB

// InspectBody reads a request body, enforces size limit, logs it,
// and returns a new reader that can be read again by the downstream handler.
// This solves the "you can only read a body once" problem.
func InspectBody(body io.ReadCloser, log io.Writer) (io.ReadCloser, int64, error) {
    if body == nil {
        return http.NoBody, 0, nil
    }

    // LimitReader: stop reading after maxBodySize (attacker can't force OOM)
    limited := io.LimitReader(body, maxBodySize+1)

    // TeeReader: every byte read from 'limited' is also written to 'log'
    var buf bytes.Buffer
    tee := io.TeeReader(limited, &buf)

    // Read all — this drains the original body
    n, err := io.Copy(io.Discard, tee)
    body.Close()

    if n > maxBodySize {
        return nil, n, fmt.Errorf("request body exceeds %d bytes", maxBodySize)
    }

    // Log the body content (buf was written to by TeeReader)
    fmt.Fprintf(log, "BODY[%d bytes]: %s\n", buf.Len(), buf.String())

    // Wrap the buffer as a new ReadCloser for downstream handler
    return io.NopCloser(bytes.NewReader(buf.Bytes())), n, err
}

func main() {
    // Simulate an HTTP request body
    originalBody := io.NopCloser(strings.NewReader(`{"username":"admin","password":"secret"}`))

    var auditLog strings.Builder
    newBody, size, err := InspectBody(originalBody, &auditLog)
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    fmt.Printf("Inspected %d bytes\n", size)
    fmt.Printf("Audit log: %s\n", auditLog.String())

    // Downstream handler reads the body again — works because we restored it
    data, _ := io.ReadAll(newBody)
    fmt.Printf("Handler got: %s\n", data)
}
```

### Real-World Context
`io.TeeReader` is used by AWS WAF, Cloudflare, and ModSecurity to inspect request bodies. The critical requirement: the body must be readable by the WAF AND by the application handler. Without `io.NopCloser(bytes.NewReader(...))`, the body would be empty when it reaches your route handler.

---

## Day 17 — encoding/json

### Learning Objectives
- Use streaming JSON decoding to avoid buffering large request bodies
- Write custom marshalers for sensitive fields (mask passwords in logs)
- Handle unknown fields for security — `json.Decoder.DisallowUnknownFields()`

### Key Concepts

```go
// Struct tags control JSON field names and behavior
type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password,omitempty"` // omit if empty
    MFA      string `json:"mfa_code,omitempty"`
}

// Decoding — use decoder, not Unmarshal, for streaming
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()  // reject {"username":"a", "admin":true}
if err := dec.Decode(&req); err != nil {
    http.Error(w, "invalid request", 400)
    return
}

// Custom marshaling — hide sensitive fields
type LogSafeAuthRequest struct {
    Username string `json:"username"`
    Password string `json:"-"` // never marshaled
}

// json.RawMessage — defer parsing of a field
type Event struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"` // parse later based on Type
}
```

### Exercise: Secure JSON API Pipeline

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strings"
)

// LoginRequest — incoming from client
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    TOTPCode string `json:"totp_code,omitempty"`
}

// Validate enforces field rules
func (r *LoginRequest) Validate() error {
    if r.Username == "" {
        return fmt.Errorf("username required")
    }
    if len(r.Username) > 64 {
        return fmt.Errorf("username too long (max 64)")
    }
    if r.Password == "" {
        return fmt.Errorf("password required")
    }
    if len(r.Password) < 8 {
        return fmt.Errorf("password too short (min 8)")
    }
    return nil
}

// LoginResponse — sent to client
type LoginResponse struct {
    Token     string `json:"token"`
    ExpiresIn int    `json:"expires_in"` // seconds
}

// AuditEntry — logged (password masked)
type AuditEntry struct {
    Username  string `json:"username"`
    Password  string `json:"-"` // NEVER log
    Action    string `json:"action"`
    Success   bool   `json:"success"`
    Reason    string `json:"reason,omitempty"`
}

func processLogin(body string) {
    dec := json.NewDecoder(strings.NewReader(body))
    dec.DisallowUnknownFields() // OWASP A03: reject parameter pollution

    var req LoginRequest
    if err := dec.Decode(&req); err != nil {
        fmt.Printf("decode error: %v\n", err)
        return
    }

    if err := req.Validate(); err != nil {
        fmt.Printf("validation error: %v\n", err)
        return
    }

    // Audit log — password field is omitted via json:"-"
    audit := AuditEntry{
        Username: req.Username,
        Action:   "login",
        Success:  true,
    }
    auditJSON, _ := json.Marshal(audit)
    fmt.Printf("AUDIT: %s\n", auditJSON)

    // Response
    resp := LoginResponse{Token: "eyJ...", ExpiresIn: 3600}
    var buf bytes.Buffer
    json.NewEncoder(&buf).Encode(resp)
    fmt.Printf("RESPONSE: %s", buf.String())
}

func main() {
    fmt.Println("=== Valid request ===")
    processLogin(`{"username":"alice","password":"s3cr3tP@ss"}`)

    fmt.Println("=== Parameter pollution ===")
    processLogin(`{"username":"alice","password":"s3cr3tP@ss","admin":true}`)

    fmt.Println("=== Missing password ===")
    processLogin(`{"username":"alice"}`)
}
```

### Real-World Context
`DisallowUnknownFields()` blocks parameter pollution — an attacker adding `"admin": true` to gain elevated access. The `json:"-"` tag prevents passwords from appearing in logs. This is OWASP A02 (Cryptographic Failures) — even a logging mistake can expose credentials.

---

## Day 18 — File I/O & os Package

### Learning Objectives
- Write a rotating log file — compliance requirement for audit trails
- Use `os.Stat` for file watching — the mechanism for hot-reloading WAF rules
- Handle file permissions correctly — prevent world-readable log files

### Key Concepts

```go
// Reading files
data, err := os.ReadFile("config.json")   // small files: read all at once
f, err := os.Open("large.log")            // large files: stream with bufio
defer f.Close()
scanner := bufio.NewScanner(f)
for scanner.Scan() {
    line := scanner.Text()
}

// Writing files
err := os.WriteFile("output.json", data, 0644)   // simple write
f, err := os.OpenFile("audit.log",                // append mode
    os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

// File watching via polling
for {
    info, err := os.Stat(configPath)
    if err == nil && info.ModTime().After(lastModTime) {
        reload(configPath)
        lastModTime = info.ModTime()
    }
    time.Sleep(5 * time.Second)
}
```

### Exercise: Rotating Audit Log Writer

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

// RotatingLog writes to date-stamped files, rotates at midnight, keeps N days
type RotatingLog struct {
    mu        sync.Mutex
    dir       string
    keepDays  int
    current   *os.File
    currentDay string
}

func NewRotatingLog(dir string, keepDays int) (*RotatingLog, error) {
    if err := os.MkdirAll(dir, 0750); err != nil {
        return nil, err
    }
    rl := &RotatingLog{dir: dir, keepDays: keepDays}
    if err := rl.rotate(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLog) rotate() error {
    today := time.Now().Format("2006-01-02")
    if rl.currentDay == today {
        return nil
    }
    if rl.current != nil {
        rl.current.Close()
    }
    path := filepath.Join(rl.dir, fmt.Sprintf("audit-%s.log", today))
    f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        return err
    }
    rl.current = f
    rl.currentDay = today
    go rl.cleanup()
    return nil
}

func (rl *RotatingLog) cleanup() {
    cutoff := time.Now().AddDate(0, 0, -rl.keepDays)
    entries, _ := os.ReadDir(rl.dir)
    for _, e := range entries {
        if e.IsDir() { continue }
        info, err := e.Info()
        if err != nil { continue }
        if info.ModTime().Before(cutoff) {
            os.Remove(filepath.Join(rl.dir, e.Name()))
        }
    }
}

func (rl *RotatingLog) Write(entry string) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if err := rl.rotate(); err != nil {
        return err
    }
    _, err := fmt.Fprintf(rl.current, "%s %s\n", time.Now().Format(time.RFC3339), entry)
    return err
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("./logs", 7)
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    defer log.Close()

    for i := 0; i < 5; i++ {
        log.Write(fmt.Sprintf(`{"event":"request","ip":"10.0.0.%d","path":"/api/data","status":200}`, i))
    }
    fmt.Println("Wrote 5 audit log entries")
}
```

### Real-World Context
GDPR Article 30 requires audit logs of data access. PCI DSS requires logs retained for 1 year, accessible for 3 months. The `0600` permission (owner read/write only) prevents other users from reading audit logs. File permissions are often an audit finding.

---

## Day 19 — Strings & Bytes

### Learning Objectives
- Normalize URLs and headers safely — prevent WAF bypass via encoding tricks
- Use `strings.Builder` for efficient string construction in hot paths
- Understand the difference between `string` and `[]byte` for security code

### Key Concepts

```go
// string is immutable — every concatenation allocates
// Use strings.Builder for building strings in loops
var sb strings.Builder
for _, part := range parts {
    sb.WriteString(part)
    sb.WriteByte('/')
}
result := sb.String() // one allocation at the end

// strings vs bytes: convert only when needed (costs allocation)
b := []byte(s)  // string → bytes copy
s := string(b)  // bytes → string copy

// Key functions for security code
strings.ToLower(s)                    // case normalization
strings.TrimSpace(s)                  // remove leading/trailing whitespace
strings.HasPrefix(s, "../")           // path traversal check
strings.Contains(s, "UNION SELECT")   // naive SQLi check
strings.ReplaceAll(s, "//", "/")      // normalize double slashes
```

### Exercise: URL Normalizer (WAF Pre-Processing)

```go
package main

import (
    "fmt"
    "net/url"
    "path"
    "sort"
    "strings"
)

// NormalizeURL prepares a URL for WAF rule matching.
// Attackers use encoding tricks to bypass naive string matching.
func NormalizeURL(rawURL string) (string, error) {
    // Step 1: Parse
    u, err := url.Parse(rawURL)
    if err != nil {
        return "", fmt.Errorf("invalid URL: %w", err)
    }

    // Step 2: Lowercase scheme and host
    u.Scheme = strings.ToLower(u.Scheme)
    u.Host = strings.ToLower(u.Host)

    // Step 3: Remove default ports
    host := u.Hostname()
    port := u.Port()
    switch {
    case u.Scheme == "http" && port == "80":
        u.Host = host
    case u.Scheme == "https" && port == "443":
        u.Host = host
    }

    // Step 4: Normalize path — resolve ./ and ../
    // url.Parse already handles %XX decoding, but path.Clean handles traversal
    cleaned := path.Clean(u.Path)
    if strings.HasSuffix(u.Path, "/") && cleaned != "/" {
        cleaned += "/"
    }
    u.Path = cleaned

    // Step 5: Sort query params (canonical order for matching)
    if u.RawQuery != "" {
        params := u.Query()
        var keys []string
        for k := range params {
            keys = append(keys, k)
        }
        sort.Strings(keys)
        var sb strings.Builder
        for i, k := range keys {
            if i > 0 { sb.WriteByte('&') }
            for j, v := range params[k] {
                if j > 0 { sb.WriteByte('&') }
                sb.WriteString(url.QueryEscape(k))
                sb.WriteByte('=')
                sb.WriteString(url.QueryEscape(v))
            }
        }
        u.RawQuery = sb.String()
    }

    return u.String(), nil
}

func main() {
    tests := []string{
        "HTTP://Example.COM:80/api/../admin/./users?b=2&a=1",
        "https://example.com:443/api//v1///resource",
        "https://example.com/path/%2e%2e/etc/passwd", // path traversal attempt
    }

    for _, raw := range tests {
        norm, err := NormalizeURL(raw)
        if err != nil {
            fmt.Printf("ERROR: %v\n", err)
            continue
        }
        fmt.Printf("IN:  %s\nOUT: %s\n\n", raw, norm)
    }
}
```

### Real-World Context
OWASP A03 attacks (Injection) rely on encoding to bypass WAF rules. `%2e%2e%2f` is `../` URL-encoded. `%252e%252e%252f` is double-encoded. Your WAF normalization pipeline (Day 73) applies multiple decoding passes before rule matching.

---

## Day 20 — regexp

### Learning Objectives
- Pre-compile regex patterns at startup — never inside a hot path
- Use named capture groups for structured log parsing
- Build the pattern-matching foundation for WAF rules (Day 69-71)

### Key Concepts

```go
// WRONG: re-compiles on every call — O(n) compilation per request
func isEmail(s string) bool {
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, s)
    return matched
}

// RIGHT: compile once at package init
var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isEmail(s string) bool {
    return emailRe.MatchString(s)
}

// Named capture groups — extract structured data
var logRe = regexp.MustCompile(
    `^(?P<ip>\S+) \S+ \S+ \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) \S+" (?P<status>\d+) (?P<bytes>\d+)`,
)
```

### Exercise: Nginx Log Parser + PII Detector

```go
package main

import (
    "bufio"
    "fmt"
    "regexp"
    "strings"
)

// Pre-compiled patterns — NEVER compile inside a loop or handler
var (
    // Nginx combined log format
    nginxRe = regexp.MustCompile(
        `^(?P<ip>\S+) \S+ \S+ \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<path>[^"]+) HTTP/[^"]*" (?P<status>\d+) (?P<bytes>\d+)`,
    )

    // WAF-relevant patterns (will expand in Week 10)
    sqlKeywordsRe  = regexp.MustCompile(`(?i)\b(UNION|SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|EXEC|EXECUTE)\b`)
    xssPatternRe   = regexp.MustCompile(`(?i)(<script|javascript:|on\w+=|<iframe)`)
    pathTraversalRe = regexp.MustCompile(`(\.\.[\\/]|%2e%2e)`)

    // PII patterns — for log scrubbing (Day 79)
    emailRe    = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
    creditCardRe = regexp.MustCompile(`\b\d{4}[\s\-]?\d{4}[\s\-]?\d{4}[\s\-]?\d{4}\b`)
)

type LogEntry struct {
    IP     string
    Time   string
    Method string
    Path   string
    Status string
    Bytes  string
    Flags  []string
}

func parseNginxLog(line string) (*LogEntry, bool) {
    m := nginxRe.FindStringSubmatch(line)
    if m == nil { return nil, false }

    names := nginxRe.SubexpNames()
    entry := &LogEntry{}
    for i, name := range names {
        switch name {
        case "ip":     entry.IP = m[i]
        case "time":   entry.Time = m[i]
        case "method": entry.Method = m[i]
        case "path":   entry.Path = m[i]
        case "status": entry.Status = m[i]
        case "bytes":  entry.Bytes = m[i]
        }
    }

    // Scan path for attack patterns
    if sqlKeywordsRe.MatchString(entry.Path) {
        entry.Flags = append(entry.Flags, "SQL_INJECTION")
    }
    if xssPatternRe.MatchString(entry.Path) {
        entry.Flags = append(entry.Flags, "XSS")
    }
    if pathTraversalRe.MatchString(entry.Path) {
        entry.Flags = append(entry.Flags, "PATH_TRAVERSAL")
    }

    return entry, true
}

// MaskPII replaces PII patterns with redacted markers
func MaskPII(s string) string {
    s = emailRe.ReplaceAllString(s, "[EMAIL]")
    s = creditCardRe.ReplaceAllString(s, "[CARD]")
    return s
}

func main() {
    logs := `
192.168.1.1 - - [10/Jan/2025:10:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234
10.0.0.1 - - [10/Jan/2025:10:00:01 +0000] "GET /api/data?id=1 UNION SELECT * FROM users-- HTTP/1.1" 200 0
203.0.113.5 - - [10/Jan/2025:10:00:02 +0000] "POST /login?email=alice@example.com HTTP/1.1" 200 256
172.16.0.1 - - [10/Jan/2025:10:00:03 +0000] "GET /../../etc/passwd HTTP/1.1" 404 0
`
    scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(logs)))
    for scanner.Scan() {
        line := MaskPII(scanner.Text())
        entry, ok := parseNginxLog(line)
        if !ok { continue }

        flagStr := ""
        if len(entry.Flags) > 0 {
            flagStr = " [" + strings.Join(entry.Flags, ",") + "]"
        }
        fmt.Printf("%s %s %s %s%s\n", entry.IP, entry.Method, entry.Path, entry.Status, flagStr)
    }
}
```

### Real-World Context
Pre-compiled regex is non-negotiable in a WAF. At 10k rps, compiling a regex per request adds ~50µs per request = 500ms wasted CPU per second. Your WAF will have 20+ patterns; compile them once at startup in `var` declarations.

---

## Day 21 — time Package

### Learning Objectives
- Build a sliding-window rate limiter — the algorithm used by AWS API Gateway and Stripe
- Use `time.NewTicker` vs `time.Tick` (the latter leaks)
- Understand monotonic clock vs wall clock in Go

### Key Concepts

```go
// time.After — creates a timer channel, fires once
// WARNING: leaks if you don't receive from it
select {
case <-time.After(5 * time.Second):
    return errors.New("timeout")
case result := <-ch:
    return result
}

// time.NewTicker — reusable ticker, MUST Stop() to prevent leak
ticker := time.NewTicker(time.Second)
defer ticker.Stop()
for range ticker.C {
    // runs every second
}

// time.Tick — convenience wrapper, LEAKS (can't stop)
// Only use for package-level tickers that run forever
for range time.Tick(time.Second) { }  // OK at top-level, bad in functions

// Monotonic clock: time.Since, time.Until use it automatically
// Wall clock: time.Now().Format(...), time.Parse use it
start := time.Now()
// ... work ...
elapsed := time.Since(start) // uses monotonic — correct even across DST changes
```

### Exercise: Sliding-Window Rate Limiter

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// SlidingWindowLimiter implements the sliding window counter algorithm.
// More accurate than fixed-window at boundary conditions.
type SlidingWindowLimiter struct {
    mu         sync.Mutex
    limit      int           // max requests per window
    window     time.Duration // window size
    buckets    map[string]*windowState
    cleanupTick *time.Ticker
}

type windowState struct {
    count     int
    windowEnd time.Time
    prevCount int       // count from previous window (for sliding approximation)
    prevEnd   time.Time
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
    l := &SlidingWindowLimiter{
        limit:       limit,
        window:      window,
        buckets:     make(map[string]*windowState),
        cleanupTick: time.NewTicker(window * 2),
    }
    go l.cleanup()
    return l
}

func (l *SlidingWindowLimiter) Allow(key string) bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    now := time.Now()
    state, ok := l.buckets[key]
    if !ok {
        state = &windowState{windowEnd: now.Add(l.window)}
        l.buckets[key] = state
    }

    // Rotate window if expired
    if now.After(state.windowEnd) {
        state.prevCount = state.count
        state.prevEnd = state.windowEnd
        state.count = 0
        state.windowEnd = now.Add(l.window)
    }

    // Sliding window: weighted sum of current + fraction of previous window
    elapsed := now.Sub(state.prevEnd)
    prevWeight := 1.0 - float64(elapsed)/float64(l.window)
    if prevWeight < 0 { prevWeight = 0 }

    estimated := float64(state.prevCount)*prevWeight + float64(state.count)
    if estimated >= float64(l.limit) {
        return false
    }

    state.count++
    return true
}

func (l *SlidingWindowLimiter) cleanup() {
    for range l.cleanupTick.C {
        l.mu.Lock()
        cutoff := time.Now().Add(-l.window * 2)
        for key, state := range l.buckets {
            if state.windowEnd.Before(cutoff) {
                delete(l.buckets, key)
            }
        }
        l.mu.Unlock()
    }
}

func (l *SlidingWindowLimiter) Stop() { l.cleanupTick.Stop() }

func main() {
    limiter := NewSlidingWindowLimiter(10, time.Second)
    defer limiter.Stop()

    clientIP := "192.168.1.1"
    allowed, blocked := 0, 0

    for i := 0; i < 25; i++ {
        if limiter.Allow(clientIP) {
            allowed++
        } else {
            blocked++
        }
    }
    fmt.Printf("Burst test: %d allowed, %d blocked (limit: 10)\n", allowed, blocked)

    // Wait for window to reset
    time.Sleep(1100 * time.Millisecond)
    allowed, blocked = 0, 0
    for i := 0; i < 5; i++ {
        if limiter.Allow(clientIP) {
            allowed++
        } else {
            blocked++
        }
    }
    fmt.Printf("After reset: %d allowed, %d blocked\n", allowed, blocked)
}
```

### Real-World Context
This rate limiter is the foundation of Day 76 (per-API-key token bucket) and Day 93 (per-IP DDoS mitigation). The sliding window algorithm avoids the "double burst" problem of fixed windows — an attacker can't send 100% of limit at the end of window 1 + 100% at the start of window 2. AWS API Gateway uses exactly this approach.

---

## Week 3 Mini Project: Structured Request Logger Middleware

Build an HTTP middleware that logs every request as structured JSON to a rotating file. This is the audit trail for your WAF and API gateway.

### Full Implementation

```go
// middleware/logger.go
package middleware

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "net"
    "net/http"
    "time"
)

type ctxKey string
const requestIDKey ctxKey = "request_id"

// requestCapture wraps ResponseWriter to capture status code and response size
type requestCapture struct {
    http.ResponseWriter
    statusCode int
    size       int64
}

func (rc *requestCapture) WriteHeader(code int) {
    rc.statusCode = code
    rc.ResponseWriter.WriteHeader(code)
}

func (rc *requestCapture) Write(b []byte) (int, error) {
    n, err := rc.ResponseWriter.Write(b)
    rc.size += int64(n)
    return n, err
}

type LogEntry struct {
    Time       string  `json:"time"`
    RequestID  string  `json:"request_id"`
    Method     string  `json:"method"`
    Path       string  `json:"path"`
    Query      string  `json:"query,omitempty"`
    Status     int     `json:"status"`
    LatencyMS  float64 `json:"latency_ms"`
    ClientIP   string  `json:"client_ip"`
    UserAgent  string  `json:"user_agent,omitempty"`
    ReqSize    int64   `json:"req_size_bytes"`
    RespSize   int64   `json:"resp_size_bytes"`
}

func generateRequestID() string {
    b := make([]byte, 8)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func realIP(r *http.Request) string {
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
        // Take only the first IP (rightmost is added by our proxy)
        parts := strings.Split(forwarded, ",")
        return strings.TrimSpace(parts[0])
    }
    host, _, _ := net.SplitHostPort(r.RemoteAddr)
    return host
}

// Logger returns a middleware that logs each request to the provided writer
func Logger(log *RotatingLog) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            requestID := generateRequestID()

            // Inject request ID into context
            ctx := context.WithValue(r.Context(), requestIDKey, requestID)
            r = r.WithContext(ctx)

            // Inject request ID into response header
            w.Header().Set("X-Request-ID", requestID)

            // Wrap writer to capture status + size
            capture := &requestCapture{ResponseWriter: w, statusCode: 200}

            next.ServeHTTP(capture, r)

            entry := LogEntry{
                Time:      time.Now().UTC().Format(time.RFC3339Nano),
                RequestID: requestID,
                Method:    r.Method,
                Path:      r.URL.Path,
                Query:     r.URL.RawQuery,
                Status:    capture.statusCode,
                LatencyMS: float64(time.Since(start).Microseconds()) / 1000.0,
                ClientIP:  realIP(r),
                UserAgent: r.UserAgent(),
                ReqSize:   r.ContentLength,
                RespSize:  capture.size,
            }

            data, _ := json.Marshal(entry)
            log.Write(string(data))
        })
    }
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
```

### Engineer Takeaway
Structured JSON logging is what lets your SIEM (Splunk, Datadog, Elastic) correlate events across services. The `request_id` field ties together the WAF log, auth log, and application log for a single request — critical for forensics after a security incident. The `0600` file permission on logs is an audit requirement.
