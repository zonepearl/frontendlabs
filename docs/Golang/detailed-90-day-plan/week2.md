# Week 2: Control Flow, Packages & Testing (Days 8–14)

**Phase:** 1 — Go Fundamentals
**Goal:** Master Go's testing story and concurrency primitives. By Friday you'll write tests the way senior engineers do, and understand goroutines well enough to reason about race conditions.

---

## Day 8 — Control Flow & defer

### Learning Objectives
- Understand `defer` execution order and why it's essential for resource cleanup
- Use `switch` without fallthrough for clean dispatch
- Know when `for` replaces `while`, `do-while`, and `loop`

### Key Concepts

```go
// defer — runs when the surrounding function returns, LIFO order
// Arguments are evaluated when defer is called, not when it runs
func processRequest(conn net.Conn) {
    defer conn.Close()          // runs last: cleanup
    defer logRequest(time.Now()) // evaluated NOW, logged later

    // ... process ...
}

// Multiple defers — LIFO (last in, first out)
func example() {
    defer fmt.Println("3") // runs third
    defer fmt.Println("2") // runs second
    defer fmt.Println("1") // runs first
}
// Output: 1, 2, 3

// switch — no fallthrough by default (unlike C/Java)
switch r.Method {
case "GET", "HEAD":
    handleRead(w, r)
case "POST", "PUT", "PATCH":
    handleWrite(w, r)
case "DELETE":
    handleDelete(w, r)
default:
    http.Error(w, "method not allowed", 405)
}

// for — the only loop in Go
for i := 0; i < 10; i++ { }      // C-style
for condition { }                   // while-style
for { }                             // infinite loop
for i, v := range slice { }        // range
for k, v := range m { }            // map range
```

### Exercise: Request Lifecycle Simulator

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

var activeRequests atomic.Int64

func simulateRequest(id int) {
    activeRequests.Add(1)

    // Defers register in order, execute in reverse — critical for connection cleanup
    defer func() {
        activeRequests.Add(-1)
        fmt.Printf("[%d] counter decremented, active=%d\n", id, activeRequests.Load())
    }()
    defer func() { fmt.Printf("[%d] 3. log request complete\n", id) }()
    defer func() { fmt.Printf("[%d] 2. release connection to pool\n", id) }()
    defer func() { fmt.Printf("[%d] 1. close response body\n", id) }()

    fmt.Printf("[%d] processing...\n", id)
    time.Sleep(10 * time.Millisecond)
    fmt.Printf("[%d] done processing\n", id)
}

func main() {
    for i := 1; i <= 3; i++ {
        simulateRequest(i)
        fmt.Println("---")
    }
    fmt.Printf("final active requests: %d\n", activeRequests.Load())
}
```

**Expected output shows cleanup in reverse order:** close body → release conn → log → decrement

### Real-World Context
Leaked connections cause `CLOSE_WAIT` sockets (you'll debug these in Day 52). Every `net.Conn`, `http.Response.Body`, and `os.File` must be `defer`-closed. The order matters: release the connection *after* you close the body.

---

## Day 9 — Packages & Modules

### Learning Objectives
- Understand Go's flat package system and why circular imports are impossible
- Structure code for testability and clear ownership
- Know what `go.mod` and `go.sum` guarantee about dependency integrity

### Key Concepts

```
// Package naming rules
package main      // executable
package config    // matches directory name
package http      // stdlib exception: matches "net/http" last segment

// Exported vs unexported
type Handler struct { }    // exported — visible outside package
type handler struct { }    // unexported — package-private

func Process() { }         // exported
func validate() { }        // unexported

// internal/ — only parent packages can import
myproject/
├── internal/
│   └── crypto/     // only myproject/* can import this
└── api/
    └── handler.go  // can import internal/crypto
```

### Project Structure for Security Proxy

```
secproxy/
├── go.mod
├── main.go
├── config/
│   └── config.go       // Config struct, Load(), Validate()
├── storage/
│   └── storage.go      // Storage interface + MemoryStorage
├── server/
│   └── server.go       // HTTP server setup, middleware chain
├── health/
│   └── health.go       // /health endpoint handler
└── internal/
    └── ratelimit/
        └── ratelimit.go // not exported — implementation detail
```

### Exercise: Refactor Week 1 into Packages

```go
// config/config.go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "time"
)

type Config struct {
    Host        string        `json:"host"`
    Port        int           `json:"port"`
    MaxConns    int           `json:"max_connections"`
    WAFMode     string        `json:"waf_mode"`
    RateLimitRPS int          `json:"rate_limit_rps"`
}

type Store struct {
    mu        sync.RWMutex
    cfg       Config
    startedAt time.Time
}

func Load(path string) (*Store, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("config.Load: %w", err)
    }
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("config.Load: invalid JSON: %w", err)
    }
    if err := cfg.validate(); err != nil {
        return nil, fmt.Errorf("config.Load: %w", err)
    }
    return &Store{cfg: cfg, startedAt: time.Now()}, nil
}

func (c *Config) validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("port %d out of range", c.Port)
    }
    if c.MaxConns < 1 {
        return fmt.Errorf("max_connections must be >= 1")
    }
    switch c.WAFMode {
    case "off", "detection", "blocking":
    default:
        return fmt.Errorf("waf_mode must be off|detection|blocking, got %q", c.WAFMode)
    }
    return nil
}

func (s *Store) Get() Config {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.cfg
}

func (s *Store) Reload(path string) error {
    newStore, err := Load(path)
    if err != nil {
        return err
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    s.cfg = newStore.cfg
    return nil
}

func (s *Store) Uptime() time.Duration { return time.Since(s.startedAt) }
```

### Real-World Context
The `internal/` directory prevents external consumers of your library from depending on implementation details. `config.Store` with `Reload()` is exactly the pattern for hot-reloading WAF rules in Day 85 — `SIGHUP` calls `Reload()` and the new config is live within milliseconds.

---

## Day 10 — Testing in Go

### Learning Objectives
- Write table-driven tests that cover happy path, edge cases, and error cases
- Use `t.Run` for subtests — readable output, independent failures
- Run tests with `-race` to catch concurrency bugs early

### Key Concepts

```go
// Table-driven test pattern — the Go idiom
func TestStorage(t *testing.T) {
    tests := []struct {
        name    string
        key     string
        value   string
        wantErr bool
    }{
        {"normal set/get", "key1", "value1", false},
        {"empty key", "", "value", true},
        {"empty value", "key", "", false},
        {"long key", strings.Repeat("x", 1000), "v", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            store := NewMemoryStorage(100)
            err := store.Set(tt.key, tt.value)
            if (err != nil) != tt.wantErr {
                t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Exercise: Storage Tests

```go
// storage/storage_test.go
package storage_test

import (
    "sync"
    "testing"

    "github.com/yourname/secproxy/storage"
)

func TestMemoryStorage_SetGet(t *testing.T) {
    tests := []struct {
        name    string
        key     string
        value   string
        wantErr bool
    }{
        {"normal", "api_key_abc", "user:42", false},
        {"empty key", "", "value", true},
        {"overwrite", "key", "v1", false},
        {"unicode key", "用户/key", "value", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := storage.NewMemoryStorage(100)
            err := s.Set(tt.key, tt.value)
            if (err != nil) != tt.wantErr {
                t.Fatalf("Set() error=%v, wantErr=%v", err, tt.wantErr)
            }
            if err != nil {
                return
            }
            got, ok := s.Get(tt.key)
            if !ok {
                t.Fatalf("Get() key not found after Set()")
            }
            if got != tt.value {
                t.Errorf("Get() = %q, want %q", got, tt.value)
            }
        })
    }
}

func TestMemoryStorage_Concurrent(t *testing.T) {
    s := storage.NewMemoryStorage(10000)
    var wg sync.WaitGroup

    // 100 goroutines writing and reading concurrently
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            key := fmt.Sprintf("key-%d", id)
            if err := s.Set(key, "value"); err != nil {
                t.Errorf("Set: %v", err)
            }
            s.Get(key)
        }(i)
    }
    wg.Wait()
}

// Run with: go test -race -v ./storage/...
```

### Real-World Context
`go test -race` uses the Go race detector — it instruments memory accesses and catches concurrent map writes that would crash in production. Require `-race` in CI. The concurrent test above catches missing `sync.RWMutex` in `MemoryStorage`.

---

## Day 11 — Benchmarks & Examples

### Learning Objectives
- Write benchmarks to get real numbers, not intuitions
- Compare implementations: map vs linear scan vs binary search
- Read `ns/op` and `allocs/op` output to make data-driven decisions

### Key Concepts

```go
// Benchmark function signature — must start with Benchmark
func BenchmarkLookup(b *testing.B) {
    data := buildTestData(1000) // setup outside loop
    b.ResetTimer()              // don't count setup time

    for i := 0; i < b.N; i++ { // b.N is chosen by the framework
        _ = lookupLinear(data, "target")
    }
}

// Run: go test -bench=. -benchmem -benchtime=3s ./...
// Output:
// BenchmarkLookupLinear-8    5234891    229 ns/op    0 B/op    0 allocs/op
// BenchmarkLookupMap-8      78234512    15.3 ns/op   0 B/op    0 allocs/op
```

### Exercise: Compare Lookup Strategies

```go
// lookup_bench_test.go
package lookup_test

import (
    "sort"
    "testing"
)

// Linear scan: O(n) but simple
func lookupLinear(data []string, target string) bool {
    for _, s := range data {
        if s == target { return true }
    }
    return false
}

// Map: O(1) amortized, but higher memory
func lookupMap(data map[string]struct{}, target string) bool {
    _, ok := data[target]
    return ok
}

// Binary search: O(log n), sorted slice, good for static data
func lookupBinary(data []string, target string) bool {
    i := sort.SearchStrings(data, target)
    return i < len(data) && data[i] == target
}

func BenchmarkLookupLinear100(b *testing.B) {
    data := makeSlice(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        lookupLinear(data, "item-099")
    }
}

func BenchmarkLookupMap100(b *testing.B) {
    data := makeMap(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        lookupMap(data, "item-099")
    }
}

func BenchmarkLookupBinary100(b *testing.B) {
    data := makeSlice(100)
    sort.Strings(data)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        lookupBinary(data, "item-099")
    }
}
```

### Real-World Context
Your WAF's CIDR blocklist lookup runs on every request. At 10k rps, a 100ns lookup costs 1ms CPU per second. A 15ns map lookup costs 0.15ms. For static data (IP blocklist, API key allowlist), a pre-built map is the right choice. Binary search works for sorted CIDRs.

---

## Day 12 — Goroutines (Introduction)

### Learning Objectives
- Understand goroutines as lightweight user-space threads (not OS threads)
- Detect goroutine leaks — one of the most common production bugs
- Know when to use goroutines vs sequential code

### Key Concepts

```go
// Starting a goroutine — the "go" keyword
go func() {
    // runs concurrently
}()

// Goroutines are cheap: ~2KB stack, grows as needed
// OS threads: ~1MB stack, fixed
// Rule of thumb: 10,000 goroutines is normal; 1,000,000 is unusual but possible

// Goroutine leak — goroutine blocked forever with no way to exit
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // blocks forever — no one sends
        fmt.Println(val)
    }()
    // function returns, goroutine still running — LEAK
}

// Fix: use context for cancellation
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return // goroutine exits when context is cancelled
        }
    }()
}
```

### Exercise: Goroutine Lifecycle Observer

```go
package main

import (
    "fmt"
    "math/rand"
    "runtime"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    start := runtime.NumGoroutine()
    fmt.Printf("goroutines before: %d\n", start)

    // Start 1000 goroutines
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // Random sleep 1-100ms
            time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
        }(i)
    }

    // Monitor goroutine count every 50ms
    done := make(chan struct{})
    go func() {
        for {
            select {
            case <-done:
                return
            case <-time.After(50 * time.Millisecond):
                fmt.Printf("  active goroutines: %d\n", runtime.NumGoroutine())
            }
        }
    }()

    wg.Wait()
    close(done)

    // Give monitor goroutine time to exit
    time.Sleep(100 * time.Millisecond)
    end := runtime.NumGoroutine()
    fmt.Printf("goroutines after: %d (started with %d)\n", end, start)
    if end > start+5 {
        fmt.Printf("WARNING: possible goroutine leak! leaked %d goroutines\n", end-start)
    }
}
```

### Real-World Context
Every `net.Conn` handler in your TCP proxy spawns goroutines. Under a DDoS, you might get 100,000 simultaneous connections — that's 200,000 goroutines (read + write per conn). Goroutine leaks compound: a leak of 1 goroutine per request at 1000 rps = 1000 leaked goroutines/second = OOM in minutes.

---

## Day 13 — Channels

### Learning Objectives
- Use channels as typed message queues between goroutines
- Know when buffered vs unbuffered channels are appropriate
- Build a worker pool — the pattern used in concurrent HTTP checkers, DNS resolvers, and log processors

### Key Concepts

```go
// Unbuffered channel — send blocks until receiver is ready (synchronizes goroutines)
ch := make(chan int)

// Buffered channel — send only blocks when full (decouples producer/consumer)
ch := make(chan int, 100)

// Closing a channel signals no more values
close(ch)
for v := range ch { } // reads until closed and drained

// select — non-blocking multi-channel dispatch
select {
case msg := <-inbound:
    handle(msg)
case <-ctx.Done():
    return ctx.Err()
case <-time.After(5 * time.Second):
    return errors.New("timeout")
default:
    // non-blocking: runs if no channel is ready
}
```

### Exercise: Worker Pool for URL Checking

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"
)

type Job struct {
    URL string
}

type Result struct {
    URL     string
    Status  int
    Latency time.Duration
    Err     error
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    client := &http.Client{Timeout: 5 * time.Second}

    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                return // channel closed, worker exits
            }
            start := time.Now()
            resp, err := client.Get(job.URL)
            result := Result{URL: job.URL, Latency: time.Since(start), Err: err}
            if err == nil {
                result.Status = resp.StatusCode
                resp.Body.Close()
            }
            select {
            case results <- result:
            case <-ctx.Done():
                return
            }
        case <-ctx.Done():
            return
        }
    }
}

func main() {
    urls := []string{
        "https://example.com",
        "https://golang.org",
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/delay/6", // will timeout
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    jobs := make(chan Job, len(urls))
    results := make(chan Result, len(urls))

    // Start 3 workers
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go worker(ctx, i, jobs, results, &wg)
    }

    // Send all jobs, close channel when done
    for _, url := range urls {
        jobs <- Job{URL: url}
    }
    close(jobs)

    // Wait for workers, then close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for r := range results {
        if r.Err != nil {
            fmt.Printf("ERROR  %s: %v\n", r.URL, r.Err)
        } else {
            fmt.Printf("%d    %s (%v)\n", r.Status, r.URL, r.Latency.Round(time.Millisecond))
        }
    }
}
```

### Real-World Context
This worker pool pattern appears everywhere: health checkers (ping 50 backends every 5s), log processors (fan-out to multiple sinks), WAF rule testers (validate 1000 payloads against rules). The `ctx.Done()` case in `select` ensures workers exit cleanly when the parent context times out.

---

## Day 14 — sync Package

### Learning Objectives
- Use `sync.Mutex` vs `sync.RWMutex` correctly for shared state
- Coordinate goroutine completion with `sync.WaitGroup`
- Run the race detector: `go test -race`

### Key Concepts

```go
// sync.RWMutex — multiple concurrent readers OR one writer
// Use when reads are much more frequent than writes (e.g., config, IP blocklist)
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}
func (s *SafeMap) Get(k string) (int, bool) {
    s.mu.RLock()         // multiple goroutines can hold RLock simultaneously
    defer s.mu.RUnlock()
    v, ok := s.m[k]
    return v, ok
}
func (s *SafeMap) Set(k string, v int) {
    s.mu.Lock()          // exclusive: blocks all readers AND writers
    defer s.mu.Unlock()
    s.m[k] = v
}

// sync.Once — exactly one initialization, thread-safe
var (
    instance *Config
    once     sync.Once
)
func GetConfig() *Config {
    once.Do(func() { instance = loadConfig() })
    return instance
}
```

### Exercise: Thread-Safe Request Counter

```go
package main

import (
    "fmt"
    "sync"
    "testing"
)

// RequestCounter tracks request counts per endpoint.
// Reads (metrics scraping) are frequent; resets are rare.
type RequestCounter struct {
    mu       sync.RWMutex
    counts   map[string]int64
    total    int64
}

func NewRequestCounter() *RequestCounter {
    return &RequestCounter{counts: make(map[string]int64)}
}

func (rc *RequestCounter) Inc(path string) {
    rc.mu.Lock()
    defer rc.mu.Unlock()
    rc.counts[path]++
    rc.total++
}

func (rc *RequestCounter) Get(path string) int64 {
    rc.mu.RLock()
    defer rc.mu.RUnlock()
    return rc.counts[path]
}

func (rc *RequestCounter) Snapshot() map[string]int64 {
    rc.mu.RLock()
    defer rc.mu.RUnlock()
    copy := make(map[string]int64, len(rc.counts))
    for k, v := range rc.counts {
        copy[k] = v
    }
    return copy
}

func (rc *RequestCounter) Reset() {
    rc.mu.Lock()
    defer rc.mu.Unlock()
    rc.counts = make(map[string]int64)
    rc.total = 0
}

// Test with race detector: go test -race -v
func TestRequestCounter_Concurrent(t *testing.T) {
    rc := NewRequestCounter()
    var wg sync.WaitGroup
    paths := []string{"/api/v1/users", "/api/v1/auth", "/health"}

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            path := paths[id%len(paths)]
            for j := 0; j < 100; j++ {
                rc.Inc(path)
            }
        }(i)
    }
    wg.Wait()

    snap := rc.Snapshot()
    var total int64
    for _, count := range snap {
        total += count
    }
    if total != 10000 {
        t.Errorf("expected 10000 total, got %d", total)
    }
}

func main() {
    rc := NewRequestCounter()

    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            rc.Inc("/api/resource")
            rc.Inc("/health")
        }(i)
    }
    wg.Wait()

    for path, count := range rc.Snapshot() {
        fmt.Printf("%s: %d requests\n", path, count)
    }
}
```

**Run the test with race detector:**
```bash
go test -race -v ./...
```

### Real-World Context
Your rate limiter (Day 76) uses this exact pattern — a `sync.RWMutex`-protected map from API key → token bucket. Without `RWMutex`, the write lock serializes all requests. With it, 1000 concurrent reads go through simultaneously, blocked only when a new bucket is created.

---

## Week 2 Mini Project: Concurrent URL Health Checker

Build a production-quality health checking tool used by your load balancer (Week 7) to detect dead backends.

### Full Implementation

```go
// healthcheck/main.go
package main

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "sync"
    "time"
)

type CheckResult struct {
    URL       string        `json:"url"`
    Status    int           `json:"status,omitempty"`
    Latency   time.Duration `json:"latency_ms"`
    Healthy   bool          `json:"healthy"`
    Error     string        `json:"error,omitempty"`
    CheckedAt time.Time     `json:"checked_at"`
}

type Config struct {
    Concurrency int
    Timeout     time.Duration
    TotalBudget time.Duration
    OutputFile  string
}

func checkURL(ctx context.Context, url string, timeout time.Duration) CheckResult {
    result := CheckResult{
        URL:       url,
        CheckedAt: time.Now(),
    }

    reqCtx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    start := time.Now()
    req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
    if err != nil {
        result.Error = err.Error()
        result.Latency = time.Since(start)
        return result
    }

    resp, err := http.DefaultClient.Do(req)
    result.Latency = time.Since(start)
    if err != nil {
        result.Error = err.Error()
        return result
    }
    defer resp.Body.Close()

    result.Status = resp.StatusCode
    result.Healthy = resp.StatusCode >= 200 && resp.StatusCode < 400
    return result
}

func run(cfg Config, urls []string) []CheckResult {
    ctx, cancel := context.WithTimeout(context.Background(), cfg.TotalBudget)
    defer cancel()

    jobs := make(chan string, len(urls))
    results := make(chan CheckResult, len(urls))

    var wg sync.WaitGroup
    for i := 0; i < cfg.Concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case url, ok := <-jobs:
                    if !ok { return }
                    results <- checkURL(ctx, url, cfg.Timeout)
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    for _, url := range urls {
        jobs <- url
    }
    close(jobs)

    go func() {
        wg.Wait()
        close(results)
    }()

    var out []CheckResult
    for r := range results {
        out = append(out, r)
    }
    return out
}

func main() {
    cfg := Config{
        Concurrency: 10,
        Timeout:     5 * time.Second,
        TotalBudget: 30 * time.Second,
        OutputFile:  "results.json",
    }

    // Read URLs from stdin or args
    var urls []string
    if len(os.Args) > 1 {
        urls = os.Args[1:]
    } else {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            if url := scanner.Text(); url != "" {
                urls = append(urls, url)
            }
        }
    }

    if len(urls) == 0 {
        // Demo mode
        urls = []string{
            "https://example.com",
            "https://golang.org",
            "https://httpbin.org/status/200",
            "https://httpbin.org/status/500",
        }
    }

    results := run(cfg, urls)

    // Write JSON output
    f, err := os.Create(cfg.OutputFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "output error: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()
    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    enc.Encode(results)

    // Print summary
    var healthy, unhealthy int
    for _, r := range results {
        if r.Healthy { healthy++ } else { unhealthy++ }
        status := "OK "
        if !r.Healthy { status = "ERR" }
        fmt.Printf("[%s] %s — %dms\n", status, r.URL, r.Latency.Milliseconds())
    }
    fmt.Printf("\n%d healthy, %d unhealthy (results in %s)\n", healthy, unhealthy, cfg.OutputFile)
}
```

### Run It
```bash
go run . https://example.com https://golang.org https://httpbin.org/status/500

# Or from a file
cat urls.txt | go run .
```

### Engineer Takeaway
This pattern (worker pool + channel pipeline + context cancellation) is in every batch job, health checker, and data pipeline. Your load balancer in Week 7 runs exactly this — 5 workers, each probing a backend every 5 seconds, with context timeout to prevent hanging checks from blocking the health state machine.
