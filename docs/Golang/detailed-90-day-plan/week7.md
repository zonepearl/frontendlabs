# Week 7: L4 Load Balancing & Resilience (Days 45–51)

**Phase:** 2 — Networking Foundations + L4
**Goal:** Build a production-grade L4 TCP load balancer with health checking, circuit breakers, and connection draining. This is what AWS NLB and HAProxy do — you'll build it from scratch.

---

> **Background in this wiki:** TCP/IP guide [Ch 57 (production load balancing)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries): the `lbdrain` lab measured 0 errors with draining vs 24 with `kill -9`, plus a PROXY protocol parser; [Ch 46 (`prober`: latency histograms)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-46-where-the-numbers-come-from-counters-flows-and-probes); HTTPS guide [Ch 9](../../v2-https/real-life-guide-v1.md#chapter-9-load-balancers-rate-limiting) and [Ch 21 (rate limits, retries, breakers)](../../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers).

## Day 45 — Load Balancing Algorithms

### Learning Objectives
- Implement round-robin, weighted round-robin, least-connections, and power-of-two
- Measure load distribution under unequal request costs
- Know which algorithm to recommend for different traffic patterns

### Key Concepts

```
Algorithm         | State     | Best for
─────────────────────────────────────────────────────
Round-robin       | none      | Homogeneous backends, equal request cost
Weighted RR       | weights   | Backends with different capacity
Least-connections | per-conn  | Variable request duration (API + slow DB queries)
Random            | none      | Simple, scales well, similar to RR
Power-of-two      | sampled   | Best for cache backends (reduces stale connections)

Real usage:
  AWS NLB:         Flow hash (5-tuple) — same client always hits same backend
  nginx upstream:  Round-robin by default, least_conn available
  Cloudflare:      Least-connections with health checking
  HAProxy:         Least-connections by default
```

### Exercise: Four Algorithms, One Interface

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
)

type Backend struct {
    Addr        string
    Weight      int
    ActiveConns atomic.Int64
}

type Balancer interface {
    Pick(backends []*Backend) *Backend
}

// RoundRobin — stateless modulo
type RoundRobin struct{ counter atomic.Uint64 }

func (rb *RoundRobin) Pick(backends []*Backend) *Backend {
    if len(backends) == 0 { return nil }
    idx := rb.counter.Add(1) - 1
    return backends[idx%uint64(len(backends))]
}

// WeightedRoundRobin — higher weight = more requests
type WeightedRoundRobin struct {
    mu      sync.Mutex
    current int
    gcd     int
    maxW    int
    cw      int // current weight
}

func gcd(a, b int) int {
    for b != 0 { a, b = b, a%b }
    return a
}

func (wrr *WeightedRoundRobin) Pick(backends []*Backend) *Backend {
    wrr.mu.Lock()
    defer wrr.mu.Unlock()
    n := len(backends)
    if n == 0 { return nil }

    // Compute GCD and max weight
    g := backends[0].Weight
    maxW := backends[0].Weight
    for _, b := range backends[1:] {
        g = gcd(g, b.Weight)
        if b.Weight > maxW { maxW = b.Weight }
    }

    for {
        wrr.current = (wrr.current + 1) % n
        if wrr.current == 0 {
            wrr.cw -= g
            if wrr.cw <= 0 {
                wrr.cw = maxW
            }
        }
        if backends[wrr.current].Weight >= wrr.cw {
            return backends[wrr.current]
        }
    }
}

// LeastConnections — pick the backend with fewest active connections
type LeastConnections struct{}

func (lc *LeastConnections) Pick(backends []*Backend) *Backend {
    if len(backends) == 0 { return nil }
    best := backends[0]
    for _, b := range backends[1:] {
        if b.ActiveConns.Load() < best.ActiveConns.Load() {
            best = b
        }
    }
    return best
}

// PowerOfTwo — sample 2 random backends, pick the less-loaded
type PowerOfTwo struct{}

func (p *PowerOfTwo) Pick(backends []*Backend) *Backend {
    if len(backends) == 0 { return nil }
    if len(backends) == 1 { return backends[0] }
    i := rand.Intn(len(backends))
    j := rand.Intn(len(backends) - 1)
    if j >= i { j++ }
    a, b := backends[i], backends[j]
    if a.ActiveConns.Load() <= b.ActiveConns.Load() {
        return a
    }
    return b
}

func main() {
    backends := []*Backend{
        {Addr: "10.0.0.1:8080", Weight: 1},
        {Addr: "10.0.0.2:8080", Weight: 2},
        {Addr: "10.0.0.3:8080", Weight: 1},
    }

    // Simulate load with variable connection counts
    backends[1].ActiveConns.Store(50)
    backends[2].ActiveConns.Store(10)

    algorithms := map[string]Balancer{
        "round-robin":        &RoundRobin{},
        "weighted-rr":        &WeightedRoundRobin{},
        "least-connections":  &LeastConnections{},
        "power-of-two":       &PowerOfTwo{},
    }

    for name, algo := range algorithms {
        counts := make(map[string]int)
        for i := 0; i < 1000; i++ {
            b := algo.Pick(backends)
            counts[b.Addr]++
        }
        fmt.Printf("\n%s (1000 reqs):\n", name)
        for addr, count := range counts {
            fmt.Printf("  %s: %d (%.1f%%)\n", addr, count, float64(count)/10)
        }
    }
}
```

### Real-World Context
AWS NLB costs $0.006/LCU-hour. Understanding these algorithms helps you answer: "Why does the NLB keep sending traffic to the slow backend?" (Round-robin ignores backend latency — use Least-connections for APIs with variable response time.)

---

## Day 46 — Backend Health Checking

### Learning Objectives
- Build an active health checker with configurable failure thresholds
- Emit state change events — hooks into your alert system
- Understand the difference between active and passive health checking

### Key Concepts

```
Active health check:
  Every N seconds: TCP connect or HTTP GET /health
  Track consecutive successes/failures
  Transition: healthy → unhealthy after K failures
  Transition: unhealthy → healthy after M successes

Passive health check:
  Observe in-flight traffic
  Track error rate over a time window
  Transition based on observed error rate

When to use which:
  Active: more reliable, detects failures before real traffic hits backend
  Passive: lower overhead, but detects failures only during real traffic
  Production: use both
```

### Exercise: Active Health Checker

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

type HealthState int

const (
    Healthy   HealthState = iota
    Unhealthy
    Unknown
)

func (s HealthState) String() string {
    switch s {
    case Healthy:   return "HEALTHY"
    case Unhealthy: return "UNHEALTHY"
    default:        return "UNKNOWN"
    }
}

type BackendHealth struct {
    addr             string
    state            HealthState
    consecutiveFails int
    consecutiveOK    int
    lastCheck        time.Time
    mu               sync.RWMutex
}

type HealthChecker struct {
    backends    []*BackendHealth
    interval    time.Duration
    failThresh  int // consecutive failures to mark unhealthy
    okThresh    int // consecutive successes to mark healthy
    onStateChange func(addr string, old, new HealthState)
}

func NewHealthChecker(addrs []string, interval time.Duration) *HealthChecker {
    hc := &HealthChecker{
        interval:   interval,
        failThresh: 3,
        okThresh:   2,
        onStateChange: func(addr string, old, new HealthState) {
            fmt.Printf("[HEALTH] %s: %s → %s\n", addr, old, new)
        },
    }
    for _, addr := range addrs {
        hc.backends = append(hc.backends, &BackendHealth{addr: addr, state: Unknown})
    }
    return hc
}

func (hc *HealthChecker) checkOne(b *BackendHealth) {
    // TCP health check: can we connect?
    conn, err := net.DialTimeout("tcp", b.addr, 2*time.Second)
    if err == nil {
        conn.Close()
    }

    b.mu.Lock()
    defer b.mu.Unlock()
    b.lastCheck = time.Now()

    oldState := b.state
    if err != nil {
        b.consecutiveFails++
        b.consecutiveOK = 0
        if b.consecutiveFails >= hc.failThresh {
            b.state = Unhealthy
        }
    } else {
        b.consecutiveOK++
        b.consecutiveFails = 0
        if b.consecutiveOK >= hc.okThresh {
            b.state = Healthy
        }
    }

    if b.state != oldState && hc.onStateChange != nil {
        go hc.onStateChange(b.addr, oldState, b.state)
    }
}

func (hc *HealthChecker) Run(stop <-chan struct{}) {
    ticker := time.NewTicker(hc.interval)
    defer ticker.Stop()
    for {
        select {
        case <-stop:
            return
        case <-ticker.C:
            for _, b := range hc.backends {
                go hc.checkOne(b)
            }
        }
    }
}

func (hc *HealthChecker) HealthyBackends() []string {
    var healthy []string
    for _, b := range hc.backends {
        b.mu.RLock()
        if b.state == Healthy {
            healthy = append(healthy, b.addr)
        }
        b.mu.RUnlock()
    }
    return healthy
}

func main() {
    checker := NewHealthChecker(
        []string{"127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083"},
        5*time.Second,
    )

    stop := make(chan struct{})
    go checker.Run(stop)

    // Run for 30s, check healthy backends every 10s
    for i := 0; i < 3; i++ {
        time.Sleep(10 * time.Second)
        fmt.Printf("Healthy backends: %v\n", checker.HealthyBackends())
    }
    close(stop)
}
```

### Real-World Context
Health checking is where most load balancer misconfigurations live. Typical mistakes: health check interval too long (60s → 60s of traffic to dead backend), failure threshold too high (3 failures × 5s interval = 15s of bad traffic), health endpoint too expensive (full DB check on each probe overwhelms the DB during recovery).

---

## Day 47 — Connection Draining

### Learning Objectives
- Implement graceful backend removal without dropping active connections
- Build the drain state machine: active → draining → removed
- Test with 50 concurrent active connections

### Exercise: Backend with Drain Support

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

type DrainState int32

const (
    StateActive   DrainState = iota
    StateDraining
    StateRemoved
)

type DrainableBackend struct {
    Addr     string
    state    atomic.Int32
    inFlight sync.WaitGroup
}

func (b *DrainableBackend) State() DrainState {
    return DrainState(b.state.Load())
}

// AcquireSlot returns true if a new connection can be sent to this backend.
// Returns false if draining or removed.
func (b *DrainableBackend) AcquireSlot() bool {
    if b.state.Load() != int32(StateActive) {
        return false
    }
    b.inFlight.Add(1)

    // Re-check state after Add (race: state may have changed between check and Add)
    if b.state.Load() != int32(StateActive) {
        b.inFlight.Done()
        return false
    }
    return true
}

func (b *DrainableBackend) ReleaseSlot() {
    b.inFlight.Done()
}

// Drain marks the backend as draining and waits for in-flight to complete
func (b *DrainableBackend) Drain(timeout time.Duration) bool {
    b.state.Store(int32(StateDraining))
    fmt.Printf("[DRAIN] %s: draining (waiting for in-flight connections)\n", b.Addr)

    done := make(chan struct{})
    go func() {
        b.inFlight.Wait()
        close(done)
    }()

    select {
    case <-done:
        b.state.Store(int32(StateRemoved))
        fmt.Printf("[DRAIN] %s: drained cleanly\n", b.Addr)
        return true
    case <-time.After(timeout):
        b.state.Store(int32(StateRemoved))
        fmt.Printf("[DRAIN] %s: drain timeout, force-removed\n", b.Addr)
        return false
    }
}

func main() {
    backend := &DrainableBackend{Addr: "10.0.0.1:8080"}

    // Simulate 50 concurrent connections
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            if !backend.AcquireSlot() {
                fmt.Printf("[conn %d] rejected (draining)\n", id)
                return
            }
            defer backend.ReleaseSlot()
            // Simulate request processing (variable duration)
            time.Sleep(time.Duration(id%500) * time.Millisecond)
        }(i)
    }

    // Drain after 100ms (while connections are still active)
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Active connections at drain time: state=ACTIVE\n")

    cleaned := backend.Drain(30 * time.Second)
    fmt.Printf("Drain result: clean=%v\n", cleaned)

    wg.Wait()
    fmt.Printf("Final state: %v\n", backend.State())
}
```

### Real-World Context
Kubernetes deployment rollovers use drain: when a new version deploys, old pods receive `SIGTERM`, stop accepting new requests, and wait for in-flight to complete. Without drain: users get connection reset errors during deploys. With drain: deploys are invisible to users. `terminationGracePeriodSeconds` is your drain timeout.

---

## Day 48 — TCP Load Balancer (L4)

### Learning Objectives
- Assemble a complete L4 TCP load balancer from the parts built this week
- Handle backend failure mid-connection gracefully
- Understand the difference between L4 and L7 load balancing

### Key Concepts

```
L4 vs L7 Load Balancing:
  L4 (TCP/UDP):
    - Decision at connection time (by IP:port or IP hash)
    - No HTTP parsing — passes bytes blindly
    - Faster (no decryption, no HTTP parsing)
    - Cannot route by URL path, host, or cookie
    - Used for: raw TCP, non-HTTP protocols, TLS passthrough

  L7 (HTTP):
    - Decision per-request (by URL, header, cookie)
    - Can terminate TLS, modify headers, compress, cache
    - Slightly slower (must parse HTTP)
    - Used for: APIs, web applications, content routing

  When to use L4:
    - Non-HTTP protocols (MySQL, Redis, gRPC with TLS passthrough)
    - Very high throughput (>1Gbps)
    - You need TLS passthrough (don't want the LB to see plaintext)
```

### Exercise: Complete L4 TCP Load Balancer

```go
package main

import (
    "fmt"
    "io"
    "net"
    "sync"
    "sync/atomic"
    "time"
)

type LB struct {
    listener  net.Listener
    backends  []*LBBackend
    balancer  Balancer
    active    atomic.Int64
    total     atomic.Int64
    rejected  atomic.Int64
    wg        sync.WaitGroup
}

type LBBackend struct {
    *BackendHealth
    *DrainableBackend
    pool *Pool
}

func NewLB(addr string, backends []string) (*LB, error) {
    ln, err := net.Listen("tcp", addr)
    if err != nil { return nil, err }

    lb := &LB{
        listener: ln,
        balancer: &LeastConnections{},
    }

    for _, bAddr := range backends {
        lb.backends = append(lb.backends, &LBBackend{
            BackendHealth:    &BackendHealth{addr: bAddr, state: Unknown},
            DrainableBackend: &DrainableBackend{Addr: bAddr},
            pool:             NewPool(20, 30*time.Second),
        })
    }

    return lb, nil
}

func (lb *LB) healthyBackends() []*Backend {
    var result []*Backend
    for _, b := range lb.backends {
        b.BackendHealth.mu.RLock()
        state := b.BackendHealth.state
        drainState := b.DrainableBackend.State()
        b.BackendHealth.mu.RUnlock()
        if state == Healthy && drainState == StateActive {
            result = append(result, &Backend{
                Addr:        b.Addr,
                ActiveConns: b.DrainableBackend.inFlight, // reuse counter
            })
        }
    }
    return result
}

func (lb *LB) Serve() {
    for {
        conn, err := lb.listener.Accept()
        if err != nil { return }

        // Pick a backend
        available := lb.healthyBackends()
        picked := lb.balancer.Pick(available)
        if picked == nil {
            conn.Close()
            lb.rejected.Add(1)
            continue
        }

        lb.active.Add(1)
        lb.total.Add(1)
        lb.wg.Add(1)
        go lb.proxy(conn, picked.Addr)
    }
}

func (lb *LB) proxy(client net.Conn, backendAddr string) {
    defer lb.wg.Done()
    defer lb.active.Add(-1)
    defer client.Close()

    upstream, err := net.DialTimeout("tcp", backendAddr, 5*time.Second)
    if err != nil {
        fmt.Printf("[LB] backend %s unavailable: %v\n", backendAddr, err)
        return
    }
    defer upstream.Close()

    // Set idle timeouts
    client.SetDeadline(time.Now().Add(5 * time.Minute))
    upstream.SetDeadline(time.Now().Add(5 * time.Minute))

    var once sync.Once
    closeAll := func() {
        once.Do(func() {
            client.Close()
            upstream.Close()
        })
    }

    var wg sync.WaitGroup
    wg.Add(2)

    // Bidirectional copy
    go func() {
        defer wg.Done()
        defer closeAll()
        io.Copy(upstream, client)
    }()
    go func() {
        defer wg.Done()
        defer closeAll()
        io.Copy(client, upstream)
    }()

    wg.Wait()
}

func main() {
    lb, err := NewLB(":8080", []string{
        "127.0.0.1:8081",
        "127.0.0.1:8082",
        "127.0.0.1:8083",
    })
    if err != nil { fmt.Println(err); return }

    fmt.Println("L4 Load Balancer on :8080")
    lb.Serve()
}
```

### Real-World Context
This is the core of HAProxy in TCP mode, AWS NLB, and Cloudflare Spectrum. The bidirectional `io.Copy` with `sync.Once` close is a subtle pattern: when one side closes, you want to immediately close the other side (not leave it hanging). The `sync.Once` prevents double-close panics.

---

## Day 49 — Rate Limiting at L4

### Learning Objectives
- Limit connections per source IP — first line of DDoS defense
- Use `sync/atomic` for lock-free counters in the accept loop
- Expose rejection metrics for anomaly detection

### Exercise: Per-IP Connection Rate Limiter

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "sync/atomic"
    "time"
)

type IPLimiter struct {
    mu           sync.Mutex
    perIP        map[string]*ipState
    globalActive atomic.Int64
    globalMax    int64
    perIPMax     int64
    cleanupTick  *time.Ticker
}

type ipState struct {
    active    atomic.Int64
    total     atomic.Int64
    lastSeen  time.Time
}

func NewIPLimiter(globalMax, perIPMax int64) *IPLimiter {
    l := &IPLimiter{
        perIP:       make(map[string]*ipState),
        globalMax:   globalMax,
        perIPMax:    perIPMax,
        cleanupTick: time.NewTicker(5 * time.Minute),
    }
    go l.cleanup()
    return l
}

func (l *IPLimiter) Allow(conn net.Conn) bool {
    // Global limit check (lock-free)
    if l.globalActive.Load() >= l.globalMax {
        return false
    }

    ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

    l.mu.Lock()
    state, ok := l.perIP[ip]
    if !ok {
        state = &ipState{}
        l.perIP[ip] = state
    }
    state.lastSeen = time.Now()
    l.mu.Unlock()

    // Per-IP limit check
    if state.active.Load() >= l.perIPMax {
        return false
    }

    l.globalActive.Add(1)
    state.active.Add(1)
    state.total.Add(1)
    return true
}

func (l *IPLimiter) Release(conn net.Conn) {
    ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
    l.globalActive.Add(-1)

    l.mu.Lock()
    if state, ok := l.perIP[ip]; ok {
        state.active.Add(-1)
    }
    l.mu.Unlock()
}

func (l *IPLimiter) cleanup() {
    for range l.cleanupTick.C {
        cutoff := time.Now().Add(-10 * time.Minute)
        l.mu.Lock()
        for ip, state := range l.perIP {
            if state.lastSeen.Before(cutoff) && state.active.Load() == 0 {
                delete(l.perIP, ip)
            }
        }
        l.mu.Unlock()
    }
}

func (l *IPLimiter) Stats() map[string]int64 {
    return map[string]int64{
        "global_active": l.globalActive.Load(),
        "global_max":    l.globalMax,
    }
}

// Integration with the TCP server
type RateLimitedServer struct {
    listener net.Listener
    limiter  *IPLimiter
}

func (s *RateLimitedServer) Serve() {
    for {
        conn, err := s.listener.Accept()
        if err != nil { return }

        if !s.limiter.Allow(conn) {
            // Immediately close — don't waste server resources
            conn.Close()
            fmt.Printf("[RATELIMIT] rejected %s\n", conn.RemoteAddr())
            continue
        }

        go func() {
            defer s.limiter.Release(conn)
            defer conn.Close()
            // handle connection...
            time.Sleep(100 * time.Millisecond)
        }()
    }
}

func main() {
    limiter := NewIPLimiter(1000, 100) // 1000 global, 100 per IP
    ln, _ := net.Listen("tcp", ":8080")
    srv := &RateLimitedServer{listener: ln, limiter: limiter}
    fmt.Println("Rate-limited server on :8080 (1000 global, 100/IP)")
    srv.Serve()
}
```

### Real-World Context
This is the first mitigation layer against DDoS. When Cloudflare detects a SYN flood, they drop packets at L3/L4 before they even reach the application. Your proxy can't do L3 packet dropping, but it can reject L4 connections immediately — "fail fast" to protect backend resources.

---

## Day 50 — Observability at L4

### Learning Objectives
- Add metrics to every component of your load balancer
- Build a `/stats` endpoint with human-readable backend status
- Know what to look at during an L4 incident

### Exercise: Load Balancer Observability

```go
package main

import (
    "encoding/json"
    "expvar"
    "fmt"
    "net/http"
    "time"
)

// L4Stats exposes all LB metrics
type L4Stats struct {
    ConnAccepted  *expvar.Int
    ConnActive    *expvar.Int
    ConnRejected  *expvar.Int
    BytesIn       *expvar.Int
    BytesOut      *expvar.Int
    ByBackend     *expvar.Map  // per-backend: connections, bytes
    HealthStates  *expvar.Map  // per-backend: health state
}

func NewL4Stats() *L4Stats {
    return &L4Stats{
        ConnAccepted: expvar.NewInt("lb_connections_accepted"),
        ConnActive:   expvar.NewInt("lb_connections_active"),
        ConnRejected: expvar.NewInt("lb_connections_rejected"),
        BytesIn:      expvar.NewInt("lb_bytes_in"),
        BytesOut:     expvar.NewInt("lb_bytes_out"),
        ByBackend:    expvar.NewMap("lb_by_backend"),
        HealthStates: expvar.NewMap("lb_backend_health"),
    }
}

func (s *L4Stats) RecordBackendHealth(addr string, state HealthState) {
    s.HealthStates.Set(addr, expvar.Func(func() any { return state.String() }))
}

// /stats — human-readable status page
func statsHandler(stats *L4Stats, checker *HealthChecker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        type backendStatus struct {
            Addr         string `json:"addr"`
            Health       string `json:"health"`
            ActiveConns  int64  `json:"active_connections"`
        }

        type response struct {
            Timestamp    string          `json:"timestamp"`
            Backends     []backendStatus `json:"backends"`
            GlobalActive int64           `json:"global_active_connections"`
            TotalAccepted int64          `json:"total_accepted"`
            TotalRejected int64          `json:"total_rejected"`
        }

        var backends []backendStatus
        for _, b := range checker.backends {
            b.mu.RLock()
            state := b.state.String()
            b.mu.RUnlock()
            backends = append(backends, backendStatus{
                Addr:   b.addr,
                Health: state,
            })
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response{
            Timestamp:    time.Now().UTC().Format(time.RFC3339),
            Backends:     backends,
            GlobalActive: stats.ConnActive.Value(),
            TotalAccepted: stats.ConnAccepted.Value(),
            TotalRejected: stats.ConnRejected.Value(),
        })
    }
}

func main() {
    stats := NewL4Stats()
    checker := NewHealthChecker(
        []string{"127.0.0.1:8081", "127.0.0.1:8082"},
        5*time.Second,
    )

    // Serve metrics and debug endpoints
    mux := http.NewServeMux()
    mux.HandleFunc("/stats", statsHandler(stats, checker))
    // /debug/vars served automatically by expvar
    // /debug/pprof served automatically by net/http/pprof
    go http.ListenAndServe(":9090", mux)

    fmt.Println("LB observability on :9090/stats and :9090/debug/vars")
    select {} // block
}
```

### Real-World Context
During a DDoS incident, the first thing you check is connection rates per backend. Is one backend getting all the traffic? Are connections being accepted and immediately closed (rejected)? Is `lb_connections_active` growing without bound? These metrics answer those questions in real-time.

---

## Day 51 — Resilience Patterns

### Learning Objectives
- Implement a circuit breaker that protects backends from cascading failures
- Build exponential backoff with jitter for retry logic
- Understand when to retry vs when to fail fast

### Key Concepts

```
Circuit Breaker States:
  Closed (normal):  requests pass through, track failures
  Open (tripped):   reject all requests immediately (fast fail)
                    protect backend from being overwhelmed during recovery
  Half-Open (probe): after timeout, let ONE request through
                     if success → close; if fail → open again

Retry with exponential backoff + jitter:
  Attempt 1: wait 100ms ± 10ms
  Attempt 2: wait 200ms ± 20ms
  Attempt 3: wait 400ms ± 40ms
  Max 3 retries

Jitter prevents "thundering herd": without jitter, all retrying clients
retry at the same time (100ms, 200ms...) and overload the recovering backend.
With jitter: retries are spread out over the window.
```

### Exercise: Circuit Breaker + Retry

```go
package main

import (
    "errors"
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

var ErrCircuitOpen = errors.New("circuit breaker open")

type CBState int32

const (
    CBClosed   CBState = iota
    CBOpen
    CBHalfOpen
)

type CircuitBreaker struct {
    state          atomic.Int32
    failures       atomic.Int64
    successes      atomic.Int64
    failThresh     int64
    successThresh  int64
    openDuration   time.Duration
    openedAt       time.Time
    mu             sync.Mutex
}

func NewCircuitBreaker(failThresh, successThresh int64, openDuration time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failThresh:    failThresh,
        successThresh: successThresh,
        openDuration:  openDuration,
    }
}

func (cb *CircuitBreaker) Allow() error {
    switch CBState(cb.state.Load()) {
    case CBClosed:
        return nil
    case CBOpen:
        cb.mu.Lock()
        defer cb.mu.Unlock()
        if time.Since(cb.openedAt) > cb.openDuration {
            cb.state.Store(int32(CBHalfOpen))
            cb.successes.Store(0)
            fmt.Println("[CB] OPEN → HALF-OPEN (probing)")
            return nil
        }
        return ErrCircuitOpen
    case CBHalfOpen:
        return nil // let one request through to probe
    }
    return nil
}

func (cb *CircuitBreaker) RecordSuccess() {
    cb.failures.Store(0)
    if CBState(cb.state.Load()) == CBHalfOpen {
        if cb.successes.Add(1) >= cb.successThresh {
            cb.state.Store(int32(CBClosed))
            fmt.Println("[CB] HALF-OPEN → CLOSED (recovered)")
        }
    }
}

func (cb *CircuitBreaker) RecordFailure() {
    if CBState(cb.state.Load()) == CBHalfOpen {
        cb.mu.Lock()
        cb.state.Store(int32(CBOpen))
        cb.openedAt = time.Now()
        cb.mu.Unlock()
        fmt.Println("[CB] HALF-OPEN → OPEN (probe failed)")
        return
    }

    failures := cb.failures.Add(1)
    if failures >= cb.failThresh {
        cb.mu.Lock()
        if CBState(cb.state.Load()) == CBClosed {
            cb.state.Store(int32(CBOpen))
            cb.openedAt = time.Now()
            fmt.Printf("[CB] CLOSED → OPEN (failed %d times)\n", failures)
        }
        cb.mu.Unlock()
    }
}

// RetryWithBackoff executes fn with exponential backoff + jitter
func RetryWithBackoff(maxAttempts int, base time.Duration, fn func() error) error {
    var lastErr error
    for attempt := 0; attempt < maxAttempts; attempt++ {
        if attempt > 0 {
            // Exponential backoff with ±20% jitter
            delay := base * (1 << uint(attempt-1))
            jitter := time.Duration(float64(delay) * 0.2 * (rand.Float64()*2 - 1))
            time.Sleep(delay + jitter)
        }

        lastErr = fn()
        if lastErr == nil {
            return nil
        }

        // Don't retry on circuit breaker open — it's already protecting the backend
        if errors.Is(lastErr, ErrCircuitOpen) {
            return lastErr
        }
        fmt.Printf("[RETRY] attempt %d failed: %v\n", attempt+1, lastErr)
    }
    return fmt.Errorf("all %d attempts failed: %w", maxAttempts, lastErr)
}

func main() {
    cb := NewCircuitBreaker(5, 2, 10*time.Second)
    failures := 0

    makeRequest := func() error {
        if err := cb.Allow(); err != nil {
            return err
        }
        failures++
        if failures <= 7 { // first 7 fail, then recover
            cb.RecordFailure()
            return fmt.Errorf("backend error")
        }
        cb.RecordSuccess()
        return nil
    }

    // Send 20 requests
    for i := 1; i <= 20; i++ {
        err := RetryWithBackoff(3, 100*time.Millisecond, makeRequest)
        state := CBState(cb.state.Load())
        stateStr := []string{"CLOSED", "OPEN", "HALF-OPEN"}[state]
        fmt.Printf("Request %d: err=%v state=%s\n", i, err, stateStr)
        time.Sleep(500 * time.Millisecond)
    }
}
```

### Real-World Context
Netflix's Hystrix (now Resilience4j) popularized circuit breakers. Without them, a slow database causes HTTP handlers to block, threads fill up, and your service crashes — this is cascading failure. With a circuit breaker, slow database → circuit opens → requests fail fast → service stays responsive. This is OWASP A05 mitigation.

---

## Week 7 Mini Project: Production L4 Load Balancer

### Full Architecture

```
                        ┌─────────────────────────┐
Clients ──TCP──▶       │   L4 Load Balancer :8080  │
                        │                           │
                        │  ┌──────────────────────┐ │
                        │  │ mTLS Termination      │ │
                        │  └──────────────────────┘ │
                        │  ┌──────────────────────┐ │
                        │  │ Per-IP Rate Limiter   │ │
                        │  └──────────────────────┘ │
                        │  ┌──────────────────────┐ │
                        │  │ Least-Conn Balancer   │ │
                        │  └──────────────────────┘ │
                        │  ┌──────────────────────┐ │
                        │  │ Circuit Breaker/Pool  │ │
                        │  └──────────────────────┘ │
                        │                           │
                        │  Observability: :9090     │
                        │  /stats /debug/vars       │
                        └─────────────────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
              Backend :8081   Backend :8082   Backend :8083
```

### Features Checklist
- [x] Least-connections balancing (switchable to round-robin)
- [x] Active health checking (TCP probe every 5s)
- [x] Connection draining (graceful backend removal)
- [x] Per-IP connection rate limiting (100/IP, 1000 global)
- [x] Circuit breaker per backend (trips at 5 failures, recovers after 30s)
- [x] Connection pool to backends (max 20 per backend)
- [x] Full metrics at `/debug/vars` and `/stats`
- [x] Graceful shutdown with drain

**Test the circuit breaker:**
```bash
# Kill backend :8082 while serving traffic
# Watch the circuit open, then test recovery when you restart it
curl http://localhost:9090/stats
```

### Engineer Takeaway
AWS NLB handles 100 million connections/day. You've now built the core algorithms that power it. The pricing model (per LCU = per connection + per bandwidth) makes sense now because you've implemented the primitives. When you configure a managed load balancer, you'll know what each setting does underneath. And when someone proposes building one in-house, you can explain why the answer is almost always "use the managed one", and when it isn't.
