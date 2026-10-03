# Week 1: Go Syntax, Types & Functions (Days 1–7)

**Phase:** 1 — Go Fundamentals
**Goal:** Go from zero to writing idiomatic Go confidently. Every networking and security system you build in Weeks 5–12 rests on this foundation.

---

## Day 1 — Go Toolchain & Your First Program

### Learning Objectives
- Understand the Go workspace model (`GOPATH` vs modules)
- Run the full build/test/format cycle
- Know where to look when engineers say "it works locally"

### Key Concepts

| Tool | What it does | When you use it |
|------|-------------|-----------------|
| `go build` | Compiles to binary | CI, deployment |
| `go run` | Compile + run | Development |
| `go test ./...` | Run all tests | Before every commit |
| `go mod tidy` | Sync `go.mod`/`go.sum` | After adding/removing imports |
| `go vet ./...` | Static analysis | CI required |
| `gofmt -w .` | Format code | Pre-commit hook |

### Setup

```bash
# Install Go (macOS)
brew install go

# Verify
go version  # go version go1.22.x darwin/arm64

# Create your workspace
mkdir -p ~/go/src/secproxy
cd ~/go/src/secproxy
go mod init github.com/yourname/secproxy
```

### Exercise: System Info Reporter

Build a program that prints a formatted system report. This is the kind of diagnostic tool your SREs run on a new server.

```go
// main.go
package main

import (
    "fmt"
    "os"
    "runtime"
    "time"
)

func main() {
    hostname, err := os.Hostname()
    if err != nil {
        hostname = "unknown"
    }

    fmt.Printf("=== System Report ===\n")
    fmt.Printf("Time:       %s\n", time.Now().Format(time.RFC3339))
    fmt.Printf("Hostname:   %s\n", hostname)
    fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
    fmt.Printf("CPUs:       %d\n", runtime.NumCPU())
    fmt.Printf("Go version: %s\n", runtime.Version())
    fmt.Printf("PID:        %d\n", os.Getpid())
}
```

**Run it:**
```bash
go run main.go
go build -o sysinfo .
./sysinfo
```

### Real-World Context
Every Go service at Cloudflare, Fastly, or your internal team prints something like this at startup. It's also in every incident runbook: "what version is running?" This program tells you.

---

## Day 2 — Variables, Types & Zero Values

### Learning Objectives
- Understand Go's static type system and why it prevents bugs
- Use `const` with `iota` for typed enumerations
- Know what zero values are and why they matter for network buffers

### Key Concepts

```go
// var declaration (explicit type)
var name string = "proxy"
var port int = 8080

// Short declaration (inferred type, only inside functions)
host := "localhost"
timeout := 30 * time.Second

// Constants — evaluated at compile time
const maxConnections = 10_000

// iota — auto-incrementing constant generator
type Protocol int

const (
    TCP  Protocol = iota // 0
    UDP                  // 1
    TLS                  // 2
    HTTP                 // 3
)
```

**Zero values** — every variable is initialized to its zero value. Critical for networking:

| Type | Zero value | Networking implication |
|------|-----------|----------------------|
| `int` | `0` | Safe default for counters |
| `bool` | `false` | Feature flags default off |
| `string` | `""` | Empty host means "not set" |
| `[]byte` | `nil` | Nil slice ≠ empty slice |
| `*T` | `nil` | Nil pointer — always check |

### Exercise: Network Unit Converter

```go
package main

import "fmt"

type ByteUnit int64

const (
    Byte     ByteUnit = 1
    Kilobyte ByteUnit = 1024
    Megabyte ByteUnit = 1024 * 1024
    Gigabyte ByteUnit = 1024 * 1024 * 1024
)

type DurationUnit int64

const (
    Millisecond DurationUnit = 1
    Second      DurationUnit = 1000
    Minute      DurationUnit = 60 * 1000
)

func convertBytes(value int64, from, to ByteUnit) float64 {
    return float64(value) * float64(from) / float64(to)
}

func main() {
    fmt.Printf("1 GB = %.0f MB\n", convertBytes(1, Gigabyte, Megabyte))
    fmt.Printf("512 MB = %.2f GB\n", convertBytes(512, Megabyte, Gigabyte))
    fmt.Printf("100 Mbps link: %.0f KB/s max throughput\n",
        convertBytes(100*1_000_000/8, Byte, Kilobyte))
}
```

### Real-World Context
When your team sets `MaxRequestBodySize = 10MB`, that's `ByteUnit` math. When you configure `IdleConnTimeout = 90 * time.Second`, that's duration constants. Typed constants prevent the `10 * 1024 * 1024` scattered everywhere.

---

## Day 3 — Functions & Multiple Return Values

### Learning Objectives
- Master Go's `(value, error)` pattern — used in every stdlib call
- Write composable functions for a config parser
- Handle errors explicitly at the call site

### Key Concepts

```go
// Multiple return values — the core Go pattern
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Named return values (use sparingly — good for short functions)
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums[1:] {
        if n < min { min = n }
        if n > max { max = n }
    }
    return // naked return
}

// Variadic functions
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// First-class functions — used extensively in middleware
type Handler func(request string) string
type Middleware func(Handler) Handler
```

### Exercise: Config Parser

This is the foundation of the WAF config system you'll build in Week 10.

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
)

type Config struct {
    Host     string
    Port     int
    MaxConns int
    Debug    bool
}

// parseConfig parses "key=value" pairs, one per line
func parseConfig(input string) (Config, error) {
    cfg := Config{
        Host:     "0.0.0.0",
        Port:     8080,
        MaxConns: 1000,
    }

    lines := strings.Split(input, "\n")
    for i, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            return Config{}, fmt.Errorf("line %d: invalid format %q (expected key=value)", i+1, line)
        }

        key := strings.TrimSpace(parts[0])
        val := strings.TrimSpace(parts[1])

        switch key {
        case "host":
            cfg.Host = val
        case "port":
            p, err := strconv.Atoi(val)
            if err != nil {
                return Config{}, fmt.Errorf("line %d: port must be integer, got %q", i+1, val)
            }
            if p < 1 || p > 65535 {
                return Config{}, fmt.Errorf("line %d: port %d out of range [1-65535]", i+1, p)
            }
            cfg.Port = p
        case "max_conns":
            n, err := strconv.Atoi(val)
            if err != nil || n < 1 {
                return Config{}, fmt.Errorf("line %d: max_conns must be positive integer", i+1)
            }
            cfg.MaxConns = n
        case "debug":
            cfg.Debug = val == "true"
        default:
            return Config{}, fmt.Errorf("line %d: unknown key %q", i+1, key)
        }
    }
    return cfg, nil
}

func main() {
    input := `
host=127.0.0.1
port=9090
max_conns=5000
debug=true
`
    cfg, err := parseConfig(input)
    if err != nil {
        fmt.Printf("config error: %v\n", err)
        return
    }
    fmt.Printf("Config: %+v\n", cfg)
}
```

### Real-World Context
Every Go service reads config this way. The `(value, error)` pattern is how you get compile-time-enforced error handling — unlike exceptions, you can't ignore a returned error silently.

---

## Day 4 — Slices & Maps

### Learning Objectives
- Understand slice internals to avoid silent bugs and memory leaks
- Use maps correctly in concurrent code (maps are not thread-safe)
- Build an in-memory request log — the foundation of your WAF audit trail

### Key Concepts

```go
// Slice internals: header = (pointer, length, capacity)
s := make([]int, 0, 10) // len=0, cap=10 — no allocation on append until cap exceeded

// append may reallocate — never assume the underlying array is the same
a := []int{1, 2, 3}
b := append(a, 4)  // b may share memory with a if cap allows, or not

// Slice of slice — shares memory
original := []byte("hello world")
sub := original[6:11]  // "world" — same backing array, not a copy
sub[0] = 'W'           // modifies original too!

// Map operations
m := make(map[string]int)
m["requests"] = 42
count, ok := m["requests"]  // ok=true when key exists
if !ok { /* key not present */ }
delete(m, "requests")

// Iterating — order is random (by design, to prevent relying on it)
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}
```

### Exercise: Ring Buffer Request Log

```go
package main

import (
    "fmt"
    "net"
    "time"
)

type Request struct {
    IP      net.IP
    Path    string
    Status  int
    Latency time.Duration
    At      time.Time
}

// RingBuffer holds the last N requests without unbounded growth.
// Critical for WAF and rate-limiting — you don't want to buffer all requests forever.
type RingBuffer struct {
    entries []Request
    head    int
    count   int
    size    int
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{entries: make([]Request, size), size: size}
}

func (rb *RingBuffer) Add(r Request) {
    rb.entries[rb.head] = r
    rb.head = (rb.head + 1) % rb.size
    if rb.count < rb.size {
        rb.count++
    }
}

func (rb *RingBuffer) All() []Request {
    result := make([]Request, rb.count)
    for i := 0; i < rb.count; i++ {
        idx := (rb.head - rb.count + i + rb.size) % rb.size
        result[i] = rb.entries[idx]
    }
    return result
}

// IPIndex provides O(1) lookup of requests by source IP
type RequestLog struct {
    ring    *RingBuffer
    byIP    map[string][]int // IP -> slice of ring indices (approximate)
}

func NewRequestLog(size int) *RequestLog {
    return &RequestLog{
        ring: NewRingBuffer(size),
        byIP: make(map[string][]int),
    }
}

func (rl *RequestLog) Add(r Request) {
    rl.ring.Add(r)
}

func main() {
    log := NewRequestLog(1000)

    // Simulate requests
    for i := 0; i < 5; i++ {
        log.Add(Request{
            IP:      net.ParseIP("192.168.1.1"),
            Path:    fmt.Sprintf("/api/resource/%d", i),
            Status:  200,
            Latency: time.Duration(i+1) * time.Millisecond,
            At:      time.Now(),
        })
    }

    for _, r := range log.ring.All() {
        fmt.Printf("[%s] %s %s %d %v\n", r.At.Format(time.RFC3339), r.IP, r.Path, r.Status, r.Latency)
    }
}
```

### Real-World Context
WAFs, IDS systems, and rate limiters all use ring buffers. An unbounded slice would OOM a production server. The map index enables O(1) "show all requests from IP X" — critical for incident response.

---

## Day 5 — Structs & Methods

### Learning Objectives
- Model real-world entities with structs
- Choose between value and pointer receivers correctly
- Use struct embedding for transparent wrappers (the pattern used in `MeteredConn` in Week 7)

### Key Concepts

```go
// Struct with field tags (used by encoding/json, database drivers)
type ServerConfig struct {
    Host     string        `json:"host"`
    Port     int           `json:"port"`
    Timeout  time.Duration `json:"timeout_ms"` // custom unmarshal needed
    TLSCert  string        `json:"tls_cert,omitempty"` // omit if empty
}

// Value receiver — method gets a copy, cannot modify
func (c ServerConfig) Address() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Pointer receiver — method can modify, avoids copy for large structs
func (c *ServerConfig) SetDefaults() {
    if c.Host == "" { c.Host = "0.0.0.0" }
    if c.Port == 0  { c.Port = 8080 }
    if c.Timeout == 0 { c.Timeout = 30 * time.Second }
}

// Rule: if ANY method uses a pointer receiver, use pointer receivers for ALL methods on that type
```

### Exercise: Server Struct

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

type Logger struct {
    prefix string
}

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), l.prefix, msg)
}

type Server struct {
    Logger               // embedded — Server.Log() works directly
    host      string
    port      int
    startedAt time.Time
    running   atomic.Bool
}

func NewServer(host string, port int) *Server {
    return &Server{
        Logger: Logger{prefix: "SERVER"},
        host:  host,
        port:  port,
    }
}

func (s *Server) Start() error {
    if s.running.Swap(true) {
        return fmt.Errorf("server already running")
    }
    s.startedAt = time.Now()
    s.Log(fmt.Sprintf("starting on %s:%d", s.host, s.port))
    return nil
}

func (s *Server) Stop() {
    if !s.running.Swap(false) {
        return
    }
    s.Log(fmt.Sprintf("stopped after %v", time.Since(s.startedAt)))
}

func (s *Server) Status() string {
    if s.running.Load() {
        return fmt.Sprintf("running for %v on %s:%d", time.Since(s.startedAt), s.host, s.port)
    }
    return "stopped"
}

func main() {
    srv := NewServer("0.0.0.0", 8080)
    srv.Start()
    fmt.Println(srv.Status())
    time.Sleep(100 * time.Millisecond)
    srv.Stop()
    fmt.Println(srv.Status())
}
```

### Real-World Context
The embedded `Logger` pattern is how Go services compose behavior without inheritance. Your WAF middleware, rate limiter, and proxy will all use this pattern. `atomic.Bool` for `running` prevents races — you'll need this when tests run `Start()` and `Stop()` concurrently.

---

## Day 6 — Interfaces

### Learning Objectives
- Define behavior with interfaces, not data
- Write code against interfaces to enable testing (mock the `Storage` interface in tests)
- Understand why Go interfaces are implicit (no `implements` keyword)

### Key Concepts

```go
// Interface = set of method signatures
// Any type that has these methods satisfies the interface — automatically
type Storage interface {
    Set(key, value string) error
    Get(key string) (string, bool)
    Delete(key string)
}

// io.Writer is an interface — net.Conn, os.File, bytes.Buffer all satisfy it
// This is why io.Copy works with any two types

// Type assertion — check the concrete type at runtime
var s Storage = NewMemoryStorage()
if ms, ok := s.(*MemoryStorage); ok {
    // access MemoryStorage-specific methods
    ms.Len()
}

// Type switch — dispatch on concrete type
switch v := s.(type) {
case *MemoryStorage:
    fmt.Println("memory, size:", v.Len())
case *RedisStorage:
    fmt.Println("redis, addr:", v.Addr())
}
```

### Exercise: Storage Interface with Metrics Wrapper

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type Storage interface {
    Set(key, value string) error
    Get(key string) (string, bool)
    Delete(key string)
}

// MemoryStorage is the real implementation
type MemoryStorage struct {
    mu   sync.RWMutex
    data map[string]string
    maxSize int
}

func NewMemoryStorage(maxSize int) *MemoryStorage {
    return &MemoryStorage{data: make(map[string]string), maxSize: maxSize}
}

func (m *MemoryStorage) Set(key, value string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if len(m.data) >= m.maxSize && m.data[key] == "" {
        return fmt.Errorf("storage full: max %d entries", m.maxSize)
    }
    m.data[key] = value
    return nil
}

func (m *MemoryStorage) Get(key string) (string, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    v, ok := m.data[key]
    return v, ok
}

func (m *MemoryStorage) Delete(key string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.data, key)
}

// MetricsStorage wraps any Storage and counts operations.
// This is the Decorator pattern — same interface, added behavior.
type MetricsStorage struct {
    inner  Storage
    gets   atomic.Int64
    sets   atomic.Int64
    hits   atomic.Int64
    misses atomic.Int64
}

func NewMetricsStorage(inner Storage) *MetricsStorage {
    return &MetricsStorage{inner: inner}
}

func (m *MetricsStorage) Set(key, value string) error {
    m.sets.Add(1)
    return m.inner.Set(key, value)
}

func (m *MetricsStorage) Get(key string) (string, bool) {
    m.gets.Add(1)
    v, ok := m.inner.Get(key)
    if ok {
        m.hits.Add(1)
    } else {
        m.misses.Add(1)
    }
    return v, ok
}

func (m *MetricsStorage) Delete(key string) { m.inner.Delete(key) }

func (m *MetricsStorage) Stats() string {
    return fmt.Sprintf("gets=%d hits=%d misses=%d sets=%d hit_rate=%.1f%%",
        m.gets.Load(), m.hits.Load(), m.misses.Load(), m.sets.Load(),
        float64(m.hits.Load())/float64(max(m.gets.Load(), 1))*100)
}

func max(a, b int64) int64 {
    if a > b { return a }
    return b
}

func main() {
    base := NewMemoryStorage(100)
    store := NewMetricsStorage(base)

    store.Set("api_key_abc123", "user:42")
    store.Set("api_key_def456", "user:99")

    if v, ok := store.Get("api_key_abc123"); ok {
        fmt.Println("auth:", v)
    }
    store.Get("api_key_unknown") // miss

    fmt.Println(store.Stats())
}
```

### Real-World Context
Your API gateway (Week 11) will have an `AuthStore` interface. In production: Redis. In tests: `MemoryStorage`. The metrics wrapper is how you add observability without modifying the core implementation — the same pattern your team uses for database clients.

---

## Day 7 — Error Handling

### Learning Objectives
- Build a multi-level error hierarchy for network/security code
- Use `errors.Is` and `errors.As` for programmatic error handling
- Understand timing-safe comparisons (relevant to auth in Week 11)

### Key Concepts

```go
// Sentinel errors — compare with errors.Is
var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")

// Wrapped errors — preserve context while keeping the chain
return fmt.Errorf("storage.Get %q: %w", key, ErrNotFound)

// Check the chain
if errors.Is(err, ErrNotFound) { /* handle */ }

// Structured error types — use when you need data in the error
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
}

// Extract the structured type
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println("invalid field:", ve.Field)
}
```

### Exercise: Security Error System

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors for the security stack
var (
    ErrNotFound      = errors.New("not found")
    ErrStorageFull   = errors.New("storage full")
    ErrInvalidKey    = errors.New("invalid key")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrRateLimited   = errors.New("rate limited")
    ErrBlocked       = errors.New("blocked")
)

// RequestError carries request context for audit logging
type RequestError struct {
    Code    int    // HTTP status to return
    Public  string // safe to show to client
    Internal string // log this, never expose
    Cause   error
}

func (e *RequestError) Error() string {
    return fmt.Sprintf("[%d] %s: %v", e.Code, e.Internal, e.Cause)
}

func (e *RequestError) Unwrap() error { return e.Cause }

func newAuthError(internal string, cause error) *RequestError {
    return &RequestError{
        Code:     401,
        Public:   "authentication required",  // never leak the real reason
        Internal: internal,
        Cause:    cause,
    }
}

func newBlockedError(reason string) *RequestError {
    return &RequestError{
        Code:     403,
        Public:   "forbidden",
        Internal: reason,
        Cause:    ErrBlocked,
    }
}

func authenticate(apiKey string) error {
    if apiKey == "" {
        return newAuthError("missing api key header", ErrUnauthorized)
    }
    if apiKey == "blocked-key" {
        return newBlockedError("api key on blocklist")
    }
    if apiKey != "valid-key" {
        return newAuthError(fmt.Sprintf("api key not found: %q", apiKey), ErrNotFound)
    }
    return nil
}

func handleRequest(apiKey string) {
    err := authenticate(apiKey)
    if err == nil {
        fmt.Println("request allowed")
        return
    }

    // Log the internal detail
    fmt.Printf("AUDIT: %v\n", err)

    // Return the public message to client
    var reqErr *RequestError
    if errors.As(err, &reqErr) {
        fmt.Printf("CLIENT RESPONSE: %d %s\n", reqErr.Code, reqErr.Public)
    }

    // Programmatic checks
    if errors.Is(err, ErrBlocked) {
        fmt.Println("ACTION: add to block list")
    }
}

func main() {
    handleRequest("valid-key")
    fmt.Println("---")
    handleRequest("")
    fmt.Println("---")
    handleRequest("blocked-key")
    fmt.Println("---")
    handleRequest("unknown-key")
}
```

### Real-World Context
Never return internal error details to API callers — that's information disclosure (OWASP A01). The `RequestError` pattern separates what you log from what you expose. `errors.Is` lets middleware handle auth errors differently from storage errors without string matching.

---

## Week 1 Mini Project: Configuration HTTP Server

Build an HTTP server that serves and updates runtime config as JSON. This is the skeleton of every internal admin tool — health endpoints, feature flags, WAF rule management.

### File Structure
```
configserver/
├── main.go
├── config/
│   └── config.go
├── storage/
│   └── storage.go
└── server/
    └── server.go
```

### Full Implementation

```go
// config/config.go
package config

import (
    "sync"
    "time"
)

type Config struct {
    mu       sync.RWMutex
    data     map[string]string
    startedAt time.Time
}

func New() *Config {
    return &Config{
        data:      make(map[string]string),
        startedAt: time.Now(),
    }
}

func (c *Config) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func (c *Config) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

func (c *Config) All() map[string]string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    copy := make(map[string]string, len(c.data))
    for k, v := range c.data {
        copy[k] = v
    }
    return copy
}

func (c *Config) Uptime() time.Duration {
    return time.Since(c.startedAt)
}
```

```go
// main.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/yourname/configserver/config"
)

func main() {
    cfg := config.New()
    cfg.Set("max_connections", "1000")
    cfg.Set("rate_limit_rps", "100")
    cfg.Set("waf_mode", "detection")

    mux := http.NewServeMux()

    // GET /config — return all config as JSON
    mux.HandleFunc("GET /config", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(cfg.All())
    })

    // PUT /config/{key} — update a value
    mux.HandleFunc("PUT /config/", func(w http.ResponseWriter, r *http.Request) {
        key := strings.TrimPrefix(r.URL.Path, "/config/")
        if key == "" {
            http.Error(w, "key required", http.StatusBadRequest)
            return
        }
        var body struct {
            Value string `json:"value"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }
        cfg.Set(key, body.Value)
        w.WriteHeader(http.StatusNoContent)
    })

    // GET /health — liveness + uptime
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]any{
            "status": "ok",
            "uptime": cfg.Uptime().String(),
            "time":   time.Now().UTC().Format(time.RFC3339),
        })
    })

    srv := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       10 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       120 * time.Second,
    }

    fmt.Println("config server listening on :8080")
    if err := srv.ListenAndServe(); err != nil {
        fmt.Printf("server error: %v\n", err)
    }
}
```

### Test It
```bash
# Start the server
go run ./...

# Get all config
curl http://localhost:8080/config

# Update a value
curl -X PUT http://localhost:8080/config/waf_mode \
  -H "Content-Type: application/json" \
  -d '{"value":"blocking"}'

# Health check
curl http://localhost:8080/health
```

### Engineer Takeaway
You can now make the everyday Go choices deliberately: `var` vs `:=`, pointer vs value receivers, and explicit error returns instead of exceptions. This Config Server is also the admin interface you'll add to your WAF in Week 10 — hot-reload rules without restarting the server.
