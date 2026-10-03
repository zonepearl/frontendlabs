# Week 4: Advanced Go Patterns (Days 22–30)

**Phase:** 1 — Go Fundamentals
**Goal:** Learn the patterns that separate junior from senior Go code. By end of week you can profile CPU, write meaningful benchmarks, and build production-grade servers with graceful shutdown.

---

## Day 22 — Generics (Go 1.18+)

### Learning Objectives
- Write generic data structures used across your entire security proxy
- Build a `RingBuffer[T]` used for metrics, log entries, and connection histories
- Understand constraints: when `comparable` vs `any`

### Key Concepts

```go
// Generic function — works for any ordered numeric type
func Min[T int | float64 | int64](a, b T) T {
    if a < b { return a }
    return b
}

// Generic type — RingBuffer[T] works for any T
type RingBuffer[T any] struct {
    items []T
    head  int
    count int
    size  int
}

// Type constraints
type Number interface {
    int | int8 | int16 | int32 | int64 | float32 | float64
}

// comparable — for map keys and equality checks
func Contains[T comparable](slice []T, item T) bool {
    for _, s := range slice {
        if s == item { return true }
    }
    return false
}
```

### Exercise: Generic RingBuffer for the Security Stack

```go
package main

import "fmt"

// RingBuffer is a fixed-size circular buffer — O(1) add, O(n) scan
// Used for: last-N request log, latency histogram, connection history
type RingBuffer[T any] struct {
    items []T
    head  int
    count int
    size  int
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
    return &RingBuffer[T]{items: make([]T, size), size: size}
}

func (rb *RingBuffer[T]) Add(item T) {
    rb.items[rb.head] = item
    rb.head = (rb.head + 1) % rb.size
    if rb.count < rb.size {
        rb.count++
    }
}

func (rb *RingBuffer[T]) All() []T {
    result := make([]T, rb.count)
    for i := 0; i < rb.count; i++ {
        result[i] = rb.items[(rb.head-rb.count+i+rb.size)%rb.size]
    }
    return result
}

func (rb *RingBuffer[T]) Len() int { return rb.count }

// Result[T] — type-safe wrapper for (value, error)
type Result[T any] struct {
    Value T
    Err   error
}

func OK[T any](v T) Result[T]     { return Result[T]{Value: v} }
func Err[T any](e error) Result[T] { return Result[T]{Err: e} }

func (r Result[T]) Unwrap() (T, error) { return r.Value, r.Err }

func main() {
    // Latency ring buffer (last 1000 request latencies for p99 calc)
    latencies := NewRingBuffer[int64](1000)
    for _, ms := range []int64{5, 12, 8, 45, 3, 100, 7, 23} {
        latencies.Add(ms)
    }
    fmt.Printf("Last %d latencies: %v\n", latencies.Len(), latencies.All())

    // Log entry ring buffer
    type LogEntry struct{ IP, Path string; Status int }
    logs := NewRingBuffer[LogEntry](100)
    logs.Add(LogEntry{"10.0.0.1", "/api/users", 200})
    logs.Add(LogEntry{"10.0.0.2", "/api/admin", 403})
    for _, e := range logs.All() {
        fmt.Printf("%s %s %d\n", e.IP, e.Path, e.Status)
    }

    // Result type for auth operations
    authResult := OK("user:42")
    if user, err := authResult.Unwrap(); err == nil {
        fmt.Println("authenticated:", user)
    }
}
```

### Real-World Context
The `RingBuffer[int64]` stores the last 1000 request latencies. A simple sort of `.All()` gives you p50, p95, p99. This is your home-grown metrics system — no Prometheus library needed. You'll use it in the Week 4 mini project and connect it to the observability layer in Week 9.

---

## Day 23 — sync.Pool & Memory Efficiency

### Learning Objectives
- Reduce garbage collection pressure in hot paths using `sync.Pool`
- Benchmark allocations with `-benchmem`
- Understand when NOT to use `sync.Pool`

### Key Concepts

```go
// sync.Pool — recycles objects to reduce GC pressure
// Pool contents are garbage collected between GC cycles (not permanent!)
var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)  // called when pool is empty
    },
}

// Get a buffer, use it, return it
func processRequest(body io.Reader) string {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()             // CRITICAL: reset before use
    defer bufPool.Put(buf)  // return after use

    io.Copy(buf, body)
    return buf.String()
}

// When NOT to use Pool:
// - Objects with long lifetimes (cache entries, connections)
// - Objects that hold references to other objects (pool contents are dropped on GC)
// - Small objects (allocation cost < pool overhead)
```

### Exercise: Buffer Pool Benchmark

```go
// pool_bench_test.go
package pool_test

import (
    "bytes"
    "io"
    "strings"
    "sync"
    "testing"
)

var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

// Without pool: allocates new buffer every call
func processWithAlloc(body string) string {
    var buf bytes.Buffer
    io.Copy(&buf, strings.NewReader(body))
    return buf.String()
}

// With pool: reuses buffers across calls
func processWithPool(body string) string {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf)
    io.Copy(buf, strings.NewReader(body))
    return buf.String()
}

func BenchmarkWithAlloc(b *testing.B) {
    body := strings.Repeat("x", 4096)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        processWithAlloc(body)
    }
}

func BenchmarkWithPool(b *testing.B) {
    body := strings.Repeat("x", 4096)
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        processWithPool(body)
    }
}

// Run: go test -bench=. -benchmem -count=3 ./...
// Expected:
// BenchmarkWithAlloc   500000   2340 ns/op   4096 B/op   1 allocs/op
// BenchmarkWithPool    800000   1480 ns/op      0 B/op   0 allocs/op
```

### Real-World Context
Your WAF body scanner reads and buffers every request body. At 10k rps with 4KB average body: without pool = 40MB/s of allocations = constant GC. With pool = 0 allocations in steady state. This is why high-performance proxies (Envoy, nginx) use pre-allocated buffer rings.

---

## Day 24 — Functional Patterns & Middleware

### Learning Objectives
- Build the middleware chain pattern used throughout your security proxy
- Use the options pattern for clean, extensible constructors
- Understand how `net/http` middleware composes

### Key Concepts

```go
// Middleware type — wraps a handler with additional behavior
type Middleware func(http.Handler) http.Handler

// Chain applies middleware in order: first middleware is outermost
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}
// Chain(handler, A, B, C) → request flows: A → B → C → handler → C → B → A

// Options pattern — extensible constructors without config struct in signature
type Server struct { port int; timeout time.Duration; tlsCert string }
type Option func(*Server)

func WithPort(port int) Option        { return func(s *Server) { s.port = port } }
func WithTimeout(d time.Duration) Option { return func(s *Server) { s.timeout = d } }
func WithTLS(cert string) Option      { return func(s *Server) { s.tlsCert = cert } }

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080, timeout: 30 * time.Second} // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### Exercise: Security Middleware Chain

```go
package main

import (
    "fmt"
    "log/slog"
    "net/http"
    "time"
)

// Middleware chain helpers

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := generateRequestID()
        w.Header().Set("X-Request-ID", id)
        r = r.WithContext(context.WithValue(r.Context(), requestIDKey, id))
        next.ServeHTTP(w, r)
    })
}

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            wrap := &statusWriter{ResponseWriter: w, code: 200}
            next.ServeHTTP(wrap, r)
            logger.Info("request",
                "method", r.Method,
                "path", r.URL.Path,
                "status", wrap.code,
                "latency_ms", time.Since(start).Milliseconds(),
                "client_ip", r.RemoteAddr,
            )
        })
    }
}

func PanicRecoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                http.Error(w, "internal server error", http.StatusInternalServerError)
                fmt.Printf("PANIC: %v\n", rec)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

type statusWriter struct {
    http.ResponseWriter
    code int
}
func (sw *statusWriter) WriteHeader(code int) {
    sw.code = code
    sw.ResponseWriter.WriteHeader(code)
}

func main() {
    logger := slog.Default()

    // The handler — application logic
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Context().Value(requestIDKey).(string)
        fmt.Fprintf(w, `{"request_id":%q,"message":"hello"}`, id)
    })

    // Build the middleware chain
    // Order: PanicRecoverer → RequestID → Logger → handler
    // (outermost first — panic recovery must be outermost)
    chain := Chain(handler,
        PanicRecoverer,
        RequestID,
        Logger(logger),
    )

    srv := NewServer(
        WithPort(8080),
        WithTimeout(30*time.Second),
    )
    fmt.Printf("Server starting on port %d\n", srv.port)

    http.ListenAndServe(fmt.Sprintf(":%d", srv.port), chain)
}
```

### Real-World Context
Your complete security proxy has this middleware chain (Day 82): `IPAllowlist → RateLimit → Auth → WAF → SchemaValidation → Handler`. The order is a security decision: rate limiting before auth prevents auth bypass attempts from exhausting rate limit budgets; WAF after auth means you know the user identity in WAF logs.

---

## Day 25 — Embedding & Composition

### Learning Objectives
- Build transparent wrappers using struct embedding
- Understand method promotion and shadowing
- Create `MeteredConn` for byte-counting — used in your TCP proxy

### Key Concepts

```go
// Embedding vs field
type Animal struct { Name string }
func (a Animal) Speak() string { return "..." }

type Dog struct {
    Animal       // embedding — Dog.Speak() works, Dog.Name works
    Breed string
}

type Cat struct {
    pet Animal   // field — must use cat.pet.Speak(), cat.pet.Name
}

// Shadowing — Dog's own method takes precedence
func (d Dog) Speak() string { return "Woof" }

// Interface embedding
type ReadWriter interface {
    io.Reader
    io.Writer
}
```

### Exercise: MeteredConn — Transparent Byte Counter

```go
package main

import (
    "fmt"
    "net"
    "sync/atomic"
)

// MeteredConn wraps net.Conn and counts bytes.
// Transparent: all net.Conn methods still work without modification.
type MeteredConn struct {
    net.Conn           // embedding: all net.Conn methods promoted
    bytesRead    atomic.Int64
    bytesWritten atomic.Int64
}

func NewMeteredConn(conn net.Conn) *MeteredConn {
    return &MeteredConn{Conn: conn}
}

// Override Read to count bytes
func (mc *MeteredConn) Read(b []byte) (int, error) {
    n, err := mc.Conn.Read(b)
    mc.bytesRead.Add(int64(n))
    return n, err
}

// Override Write to count bytes
func (mc *MeteredConn) Write(b []byte) (int, error) {
    n, err := mc.Conn.Write(b)
    mc.bytesWritten.Add(int64(n))
    return n, err
}

func (mc *MeteredConn) Stats() (read, written int64) {
    return mc.bytesRead.Load(), mc.bytesWritten.Load()
}

// TimedHTTPHandler wraps http.Handler and records latency
type TimedHTTPHandler struct {
    http.Handler
    histogram *RingBuffer[int64]
}

func NewTimedHTTPHandler(h http.Handler, histSize int) *TimedHTTPHandler {
    return &TimedHTTPHandler{Handler: h, histogram: NewRingBuffer[int64](histSize)}
}

func (th *TimedHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    th.Handler.ServeHTTP(w, r)
    th.histogram.Add(time.Since(start).Milliseconds())
}

func (th *TimedHTTPHandler) P99() int64 {
    samples := th.histogram.All()
    if len(samples) == 0 { return 0 }
    sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
    idx := int(float64(len(samples)) * 0.99)
    if idx >= len(samples) { idx = len(samples) - 1 }
    return samples[idx]
}

func main() {
    // Create a TCP connection pair for testing
    server, client := net.Pipe()
    defer server.Close()
    defer client.Close()

    // Wrap client connection with metering
    metered := NewMeteredConn(client)

    // Write some data
    go func() {
        metered.Write([]byte("GET /api/data HTTP/1.1\r\nHost: example.com\r\n\r\n"))
    }()

    // Read from server side
    buf := make([]byte, 1024)
    n, _ := server.Read(buf)
    server.Write(buf[:n]) // echo back

    // Read response
    metered.Read(buf)

    read, written := metered.Stats()
    fmt.Printf("Bytes read: %d, written: %d\n", read, written)
}
```

### Real-World Context
`MeteredConn` is used in your TCP proxy (Week 5) to count bytes per connection and expose them via `expvar`. This is how your load balancer tracks backend bandwidth usage — billing, capacity planning, and anomaly detection (a connection suddenly transferring 10x normal volume).

---

## Day 26 — expvar & Runtime Metrics

### Learning Objectives
- Expose runtime metrics via `/debug/vars` — zero-dependency observability
- Build the metrics foundation for your proxy's `/health` endpoint
- Understand why self-instrumentation matters before adding Prometheus

### Key Concepts

```go
// expvar — built-in, zero-dependency metrics
import "expvar"

var (
    reqTotal  = expvar.NewInt("requests_total")
    reqErrors = expvar.NewInt("requests_errors")
    reqLatency = expvar.NewFloat("requests_latency_p99_ms")
    connActive = expvar.NewInt("connections_active")
)

// Increment
reqTotal.Add(1)
connActive.Add(1)
// Decrement on close
defer connActive.Add(-1)

// Map for per-endpoint stats
var endpointStats = expvar.NewMap("endpoints")
func recordEndpoint(path string) {
    endpointStats.Add(path, 1)
}

// Accessed at: GET /debug/vars
// Returns JSON: {"requests_total": 1234, "connections_active": 42, ...}
```

### Exercise: Instrumented Proxy Metrics

```go
package main

import (
    "expvar"
    "net/http"
    _ "net/http/pprof"   // registers /debug/pprof automatically
    "time"
)

// ProxyMetrics collects metrics for the security proxy
type ProxyMetrics struct {
    ReqTotal    *expvar.Int
    ReqErrors   *expvar.Int
    ReqBlocked  *expvar.Int    // WAF blocks
    ConnActive  *expvar.Int
    ConnTotal   *expvar.Int
    BytesIn     *expvar.Int
    BytesOut    *expvar.Int
    ByEndpoint  *expvar.Map    // requests per path
    ByStatus    *expvar.Map    // requests per status code

    // latency ring buffer (not expvar — compute percentiles ourselves)
    latencies   *RingBuffer[int64]
    latMu       sync.Mutex
}

func NewProxyMetrics() *ProxyMetrics {
    m := &ProxyMetrics{
        ReqTotal:   expvar.NewInt("proxy_requests_total"),
        ReqErrors:  expvar.NewInt("proxy_requests_errors"),
        ReqBlocked: expvar.NewInt("proxy_requests_blocked"),
        ConnActive: expvar.NewInt("proxy_connections_active"),
        ConnTotal:  expvar.NewInt("proxy_connections_total"),
        BytesIn:    expvar.NewInt("proxy_bytes_in"),
        BytesOut:   expvar.NewInt("proxy_bytes_out"),
        ByEndpoint: expvar.NewMap("proxy_by_endpoint"),
        ByStatus:   expvar.NewMap("proxy_by_status"),
        latencies:  NewRingBuffer[int64](10000),
    }

    // Register computed P99 as expvar.Func (re-computed on each scrape)
    expvar.Publish("proxy_latency_p99_ms", expvar.Func(func() any {
        return m.percentile(99)
    }))
    expvar.Publish("proxy_latency_p95_ms", expvar.Func(func() any {
        return m.percentile(95)
    }))

    return m
}

func (m *ProxyMetrics) RecordRequest(path, status string, latencyMS int64, blocked bool) {
    m.ReqTotal.Add(1)
    m.ByEndpoint.Add(path, 1)
    m.ByStatus.Add(status, 1)
    if blocked { m.ReqBlocked.Add(1) }

    m.latMu.Lock()
    m.latencies.Add(latencyMS)
    m.latMu.Unlock()
}

func (m *ProxyMetrics) percentile(p float64) float64 {
    m.latMu.Lock()
    samples := m.latencies.All()
    m.latMu.Unlock()

    if len(samples) == 0 { return 0 }
    sorted := make([]int64, len(samples))
    copy(sorted, samples)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
    idx := int(float64(len(sorted)) * p / 100)
    if idx >= len(sorted) { idx = len(sorted) - 1 }
    return float64(sorted[idx])
}

func main() {
    metrics := NewProxyMetrics()

    // Simulate requests
    for i := 0; i < 100; i++ {
        latency := int64(rand.Intn(100) + 1)
        metrics.RecordRequest("/api/data", "200", latency, false)
    }
    metrics.RecordRequest("/api/admin", "403", 5, true)

    http.ListenAndServe(":8080", nil) // /debug/vars and /debug/pprof served automatically
}
```

### Real-World Context
`/debug/vars` is the first observability tool you add to any Go service. When an incident happens and your Datadog is down, `curl http://service:8080/debug/vars` still gives you request rates, error counts, and latency percentiles. Every internal service should expose this endpoint — not to the internet, but behind your VPN or internal load balancer.

---

## Day 27 — Graceful Shutdown

### Learning Objectives
- Implement 3-stage graceful shutdown: stop accepting → drain in-flight → force close
- Handle `SIGTERM` and `SIGINT` correctly
- Prevent connection drops during deployment

### Key Concepts

```go
// Signal handling
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
defer stop()

// http.Server graceful shutdown
srv := &http.Server{Addr: ":8080", Handler: handler}
go srv.ListenAndServe()

<-ctx.Done()  // wait for signal
shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(shutdownCtx)  // stops accepting, waits for in-flight to complete
```

### Exercise: 3-Stage Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os/signal"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
)

type GracefulServer struct {
    srv        *http.Server
    inFlight   sync.WaitGroup
    reqCount   atomic.Int64
    shutting   atomic.Bool
}

func (g *GracefulServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if g.shutting.Load() {
        // Stage 1: stop accepting new requests (return 503 with Retry-After)
        w.Header().Set("Retry-After", "10")
        http.Error(w, "service shutting down", http.StatusServiceUnavailable)
        return
    }

    g.inFlight.Add(1)
    defer g.inFlight.Done()
    g.reqCount.Add(1)

    // Simulate work — in a real proxy, this forwards the request
    time.Sleep(100 * time.Millisecond)
    fmt.Fprintf(w, `{"request":%d}`, g.reqCount.Load())
}

func (g *GracefulServer) Shutdown() {
    fmt.Println("shutdown: stage 1 — stopping new requests")
    g.shutting.Store(true)

    // Stage 2: wait for in-flight with 30s budget
    done := make(chan struct{})
    go func() {
        g.inFlight.Wait()
        close(done)
    }()

    fmt.Println("shutdown: stage 2 — draining in-flight requests")
    select {
    case <-done:
        fmt.Println("shutdown: all requests drained cleanly")
    case <-time.After(30 * time.Second):
        fmt.Println("shutdown: drain timeout, force-closing remaining")
    }

    // Stage 3: shut down the HTTP server
    fmt.Println("shutdown: stage 3 — closing server")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    g.srv.Shutdown(ctx)
}

func main() {
    gs := &GracefulServer{}
    gs.srv = &http.Server{Addr: ":8080", Handler: gs}

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()

    go func() {
        fmt.Println("server listening on :8080")
        if err := gs.srv.ListenAndServe(); err != http.ErrServerClosed {
            fmt.Printf("server error: %v\n", err)
        }
    }()

    <-ctx.Done()
    fmt.Println("signal received")
    gs.Shutdown()
    fmt.Println("shutdown complete")
}
```

**Test it:**
```bash
# Start server
go run .

# Send some requests (in another terminal)
for i in $(seq 1 20); do curl -s http://localhost:8080/ & done

# Send SIGINT while requests are in-flight
kill -SIGINT $(pgrep -f 'go run')
```

### Real-World Context
Kubernetes sends `SIGTERM` when scaling down or deploying. Your pod has a `terminationGracePeriodSeconds` budget (default 30s). If your server doesn't drain in-flight requests, users see connection resets. Every production Go service must implement graceful shutdown.

---

## Day 28 — Profiling

### Learning Objectives
- Generate and read CPU and heap profiles
- Identify hot functions in a flame graph
- Know what to look for when engineers say "the WAF is slow"

### Key Concepts

```bash
# Add pprof endpoints (one import)
import _ "net/http/pprof"

# Collect CPU profile (30 seconds)
curl -s http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# Analyze
go tool pprof cpu.prof
(pprof) top10        # top 10 functions by CPU
(pprof) web          # open flame graph in browser (requires graphviz)
(pprof) list regexp  # source code for regexp package

# Heap profile
curl -s http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
(pprof) top          # top allocators
(pprof) allocs       # allocation profile
```

### Exercise: Profile the WAF Pattern Matcher

```go
package main

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"
    "regexp"
    "time"
)

// Simulate WAF pattern matching under load
var wafPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)\b(UNION|SELECT|INSERT|DELETE|DROP)\b`),
    regexp.MustCompile(`(?i)<script[\s>]`),
    regexp.MustCompile(`(?i)javascript:`),
    regexp.MustCompile(`\.\./`),
    regexp.MustCompile(`(?i)(cmd|exec|system|eval)\s*\(`),
}

func checkPayload(payload string) bool {
    for _, re := range wafPatterns {
        if re.MatchString(payload) {
            return true // blocked
        }
    }
    return false
}

func handler(w http.ResponseWriter, r *http.Request) {
    payload := r.URL.Query().Get("q")
    blocked := checkPayload(payload)
    if blocked {
        http.Error(w, "forbidden", 403)
        return
    }
    fmt.Fprint(w, "ok")
}

func main() {
    http.HandleFunc("/", handler)

    // Load generator: 500 goroutines hammering for 30s
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("starting load test...")
        payloads := []string{"normal search", "product name", "UNION SELECT", "<script>"}
        var wg sync.WaitGroup
        for i := 0; i < 500; i++ {
            wg.Add(1)
            go func(id int) {
                defer wg.Done()
                for j := 0; j < 100; j++ {
                    payload := payloads[j%len(payloads)]
                    http.Get(fmt.Sprintf("http://localhost:8080/?q=%s", url.QueryEscape(payload)))
                }
            }(i)
        }
        wg.Wait()
        fmt.Println("load test complete. check: go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30")
    }()

    fmt.Println("profiling server on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Profile it:**
```bash
go run . &
sleep 2
go tool pprof -http=:9090 http://localhost:8080/debug/pprof/profile?seconds=10
# Open http://localhost:9090 for flame graph
```

### Real-World Context
When your WAF adds more latency than expected, profiling will show you whether it's regex, JSON parsing, or mutex contention. A flame graph makes the answer obvious in 30 seconds — no guessing. This is the skill that separates engineers who fix performance issues from ones who guess.

---

## Day 29 — Testing: httptest & net.Pipe

### Learning Objectives
- Test HTTP handlers without starting a real server
- Test TCP protocol code using `net.Pipe()` — an in-memory bidirectional connection
- Write integration tests that cover your full middleware chain

### Key Concepts

```go
// httptest.NewRecorder — captures response without network
w := httptest.NewRecorder()
r := httptest.NewRequest("GET", "/api/data?id=1", nil)
handler.ServeHTTP(w, r)
// Assert
assert(w.Code == 200)
assert(w.Header().Get("X-Request-ID") != "")

// httptest.NewServer — real listening server on random port
ts := httptest.NewServer(handler)
defer ts.Close()
resp, err := http.Get(ts.URL + "/api/data")

// httptest.NewTLSServer — TLS server with self-signed cert
ts := httptest.NewTLSServer(handler)
defer ts.Close()
client := ts.Client()  // pre-configured to trust the test cert
resp, err := client.Get(ts.URL + "/api/data")

// net.Pipe — in-memory bidirectional TCP connection (no network)
server, client := net.Pipe()
```

### Exercise: Test the Logger Middleware

```go
// middleware/logger_test.go
package middleware_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestLoggerMiddleware(t *testing.T) {
    var logBuffer strings.Builder

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(201)
        w.Write([]byte(`{"created":true}`))
    })

    logged := LoggerMiddleware(&logBuffer)(handler)

    tests := []struct {
        name           string
        method, path   string
        wantStatus     int
        wantLogContains []string
    }{
        {
            name: "GET request logged",
            method: "GET", path: "/api/users",
            wantStatus: 201,
            wantLogContains: []string{"GET", "/api/users", "201", "request_id"},
        },
        {
            name: "request ID injected in response",
            method: "POST", path: "/api/items",
            wantStatus: 201,
            wantLogContains: []string{"POST", "/api/items"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            logBuffer.Reset()

            w := httptest.NewRecorder()
            r := httptest.NewRequest(tt.method, tt.path, nil)
            logged.ServeHTTP(w, r)

            if w.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
            }

            if w.Header().Get("X-Request-ID") == "" {
                t.Error("X-Request-ID header not set")
            }

            logEntry := logBuffer.String()
            for _, want := range tt.wantLogContains {
                if !strings.Contains(logEntry, want) {
                    t.Errorf("log missing %q in: %s", want, logEntry)
                }
            }

            // Validate log is valid JSON
            var entry map[string]any
            if err := json.Unmarshal([]byte(strings.TrimSpace(logEntry)), &entry); err != nil {
                t.Errorf("log is not valid JSON: %v", err)
            }
        })
    }
}

func TestNetPipe_ProtocolTest(t *testing.T) {
    // Test TCP protocol code without a real network
    server, client := net.Pipe()
    defer server.Close()
    defer client.Close()

    // Server goroutine: echo server
    go func() {
        buf := make([]byte, 1024)
        n, _ := server.Read(buf)
        server.Write(buf[:n])
    }()

    // Client sends and reads echo
    msg := []byte("hello proxy")
    client.Write(msg)
    response := make([]byte, len(msg))
    n, err := client.Read(response)
    if err != nil { t.Fatal(err) }
    if string(response[:n]) != string(msg) {
        t.Errorf("echo = %q, want %q", response[:n], msg)
    }
}
```

### Real-World Context
`httptest` tests are 100x faster than spinning up a real server. Your CI runs 500 middleware tests per second vs 5 per second with real HTTP. `net.Pipe()` lets you test your TCP proxy, protocol parser, and TLS handshake without touching the network stack.

---

## Day 30 — Phase 1 Review & Code Audit

### Learning Objectives
- Apply consistent error handling across all Phase 1 code
- Find and fix goroutine leaks with `-race`
- Document technical debt before moving to networking

### Review Checklist

```
For each file:
[ ] All errors checked (no _ for error returns)
[ ] All goroutines have an exit path (channel close, ctx.Done(), or timeout)
[ ] All net.Conn / os.File / http.Response.Body have defer Close()
[ ] No global mutable state without sync.Mutex or sync/atomic
[ ] No panic() in library code (only main and tests)
[ ] Contexts propagated through all function calls
[ ] Tests cover error paths, not just happy path
[ ] go test -race passes
[ ] go vet passes
```

**Run the full check:**
```bash
go vet ./...
go test -race ./...
go test -bench=. -benchmem ./...
```

---

## Week 4 Mini Project: In-Memory Metrics Server

Build a self-contained metrics collection and query server. This is the observability layer you'll wire into your security proxy.

```go
// metrics/server.go
package main

import (
    "encoding/json"
    "expvar"
    "fmt"
    "math"
    "net/http"
    _ "net/http/pprof"
    "sort"
    "sync"
    "time"
)

type Sample struct {
    Value float64
    At    time.Time
}

type MetricStore struct {
    mu      sync.RWMutex
    metrics map[string]*RingBuffer[Sample]
    maxAge  time.Duration
}

func NewMetricStore(maxAge time.Duration) *MetricStore {
    return &MetricStore{
        metrics: make(map[string]*RingBuffer[Sample]),
        maxAge:  maxAge,
    }
}

func (ms *MetricStore) Record(name string, value float64) {
    ms.mu.Lock()
    defer ms.mu.Unlock()
    if _, ok := ms.metrics[name]; !ok {
        ms.metrics[name] = NewRingBuffer[Sample](10000)
    }
    ms.metrics[name].Add(Sample{Value: value, At: time.Now()})
}

type Stats struct {
    Count int     `json:"count"`
    Min   float64 `json:"min"`
    P50   float64 `json:"p50"`
    P95   float64 `json:"p95"`
    P99   float64 `json:"p99"`
    Max   float64 `json:"max"`
}

func (ms *MetricStore) Query(name string, window time.Duration) Stats {
    ms.mu.RLock()
    defer ms.mu.RUnlock()

    buf, ok := ms.metrics[name]
    if !ok { return Stats{} }

    cutoff := time.Now().Add(-window)
    var values []float64
    for _, s := range buf.All() {
        if s.At.After(cutoff) {
            values = append(values, s.Value)
        }
    }
    if len(values) == 0 { return Stats{} }

    sort.Float64s(values)
    pct := func(p float64) float64 {
        idx := int(math.Ceil(p/100*float64(len(values)))) - 1
        if idx < 0 { idx = 0 }
        return values[idx]
    }

    return Stats{
        Count: len(values),
        Min:   values[0],
        P50:   pct(50),
        P95:   pct(95),
        P99:   pct(99),
        Max:   values[len(values)-1],
    }
}

func main() {
    store := NewMetricStore(24 * time.Hour)

    // expvar for live metrics
    expvar.Publish("metric_names", expvar.Func(func() any {
        return len(store.metrics)
    }))

    mux := http.NewServeMux()

    // POST /metrics — ingest
    mux.HandleFunc("POST /metrics", func(w http.ResponseWriter, r *http.Request) {
        var body struct {
            Name  string  `json:"name"`
            Value float64 `json:"value"`
        }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
            http.Error(w, "invalid JSON", 400)
            return
        }
        store.Record(body.Name, body.Value)
        w.WriteHeader(http.StatusNoContent)
    })

    // GET /metrics/{name}?window=60s — query
    mux.HandleFunc("GET /metrics/", func(w http.ResponseWriter, r *http.Request) {
        name := strings.TrimPrefix(r.URL.Path, "/metrics/")
        windowStr := r.URL.Query().Get("window")
        window, err := time.ParseDuration(windowStr)
        if err != nil { window = time.Minute }

        stats := store.Query(name, window)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
    })

    fmt.Println("metrics server on :8080")
    srv := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
    }
    srv.ListenAndServe()
}
```

**Test it:**
```bash
# Ingest metrics
curl -X POST http://localhost:8080/metrics \
  -d '{"name":"request_latency_ms","value":45.2}'

# Query with 60s window
curl "http://localhost:8080/metrics/request_latency_ms?window=60s"
# {"count":1,"min":45.2,"p50":45.2,"p95":45.2,"p99":45.2,"max":45.2}
```

### Engineer Takeaway
This is a miniature Prometheus in ~200 lines. You now know why observability systems need ring buffers (bounded memory), percentile queries (p99 is what users experience, not average), and structured ingestion endpoints. This metrics server wires into every component you build in Phases 2 and 3.
