# Week 13: DDoS Detection & Mitigation (Days 91–98)

> **Phase 4 Goal:** Add the final security layers — DDoS mitigation and bot detection — to the production security proxy built in Phases 1–3.

**Why this week matters:** DDoS attacks are the most common attack against internet-facing services. Understanding the mechanics lets you evaluate DDoS vendors (Cloudflare, AWS Shield, Akamai), configure existing defenses correctly, and architect a first line of defense using nothing but the Go standard library.

---

> **Background in this wiki:** TCP/IP guide [Ch 21 (SYN floods and SYN cookies)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake), [Ch 52 (SYN and accept queues, overflowed on purpose)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack), [Ch 20 (UDP)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-20-udp-fire-and-forget); HTTPS guide [Ch 19 (Slowloris, measured)](../../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown), [Ch 21 (rate limiting)](../../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers), [Ch 23 (burn-rate alerting)](../../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services). Security in Depth: [Chapter 40B (rate limiting in depth)](../../security/real-life-security-guide-v1.md#chapter-40b-api-rate-limiting-in-depth) and [Chapter 28 (resilience as a security property)](../../security/real-life-security-guide-v1.md#chapter-28-resilience-as-a-security-property).

## Day 91 — DDoS Taxonomy & Attack Simulation

| | |
|---|---|
| **Learning objective** | Understand the four DDoS categories and what mitigation each requires |
| **Key packages** | `net`, `net/http`, `time`, `sync` |
| **Security context** | OWASP A05 — Security Misconfiguration; knowing attack types prevents deploying the wrong mitigation |

### Key Concepts

**DDoS attack categories:**

| Category | Mechanism | Go-side impact | Primary defense |
|----------|-----------|----------------|-----------------|
| **Volumetric** | Saturate bandwidth (Gbps) | Cannot defend at app layer | Upstream scrubbing (Cloudflare, AWS Shield) |
| **Protocol** | SYN flood, connection exhaustion | `net.Listen` backlog fills; `Accept` goroutines starve | Per-IP connection limits, TCP syncookies |
| **Application** | HTTP flood, slowloris, slow POST | Goroutine/memory exhaustion | Rate limiting, request timeouts |
| **Amplification** | UDP reflection (DNS, NTP, memcached) | Your service weaponized against others | Response size cap, per-IP rate limits |

### Exercise — Attack Pattern Simulators (localhost only)

```go
package main

import (
    "context"
    "fmt"
    "net"
    "net/http"
    "sync"
    "time"
)

// 1. TCP Connection Flood — open many connections, send no data
func simulateTCPFlood(target string, count int) {
    var wg sync.WaitGroup
    conns := make([]net.Conn, 0, count)
    var mu sync.Mutex

    for i := 0; i < count; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            conn, err := net.DialTimeout("tcp", target, 2*time.Second)
            if err != nil {
                return
            }
            mu.Lock()
            conns = append(conns, conn)
            mu.Unlock()
            // Hold the connection open — simulate half-open
            time.Sleep(30 * time.Second)
            conn.Close()
        }()
    }

    time.Sleep(2 * time.Second)
    fmt.Printf("TCP flood: opened %d connections to %s\n", len(conns), target)
    wg.Wait()
}

// 2. HTTP Flood — maximum request rate
func simulateHTTPFlood(target string, rps int, duration time.Duration) {
    client := &http.Client{Timeout: 5 * time.Second}
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    var sent, errs int64
    var mu sync.Mutex
    sem := make(chan struct{}, rps) // bound concurrency to target RPS

    ticker := time.NewTicker(time.Second / time.Duration(rps))
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            mu.Lock()
            fmt.Printf("HTTP flood: %d requests sent, %d errors\n", sent, errs)
            mu.Unlock()
            return
        case <-ticker.C:
            sem <- struct{}{}
            go func() {
                defer func() { <-sem }()
                resp, err := client.Get(target)
                mu.Lock()
                defer mu.Unlock()
                if err != nil {
                    errs++
                    return
                }
                resp.Body.Close()
                sent++
            }()
        }
    }
}

// 3. Slowloris — send headers one byte at a time
func simulateSlowloris(target string, connCount int) {
    // Parse host/port
    host, port := "127.0.0.1", "8080"
    _ = host
    _ = port

    var wg sync.WaitGroup
    for i := 0; i < connCount; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            conn, err := net.DialTimeout("tcp", target, 2*time.Second)
            if err != nil {
                return
            }
            defer conn.Close()

            // Send partial HTTP request — real headers drip in slowly
            conn.Write([]byte(fmt.Sprintf("GET /?id=%d HTTP/1.1\r\nHost: localhost\r\n", id)))

            // Drip one byte per second — Slowloris pattern
            headers := []byte("X-Slow: ")
            for j := 0; j < len(headers); j++ {
                time.Sleep(1 * time.Second)
                conn.Write(headers[j : j+1])
            }
            // Never send the final \r\n\r\n — connection stays open
            time.Sleep(30 * time.Second)
        }(i)
    }
    wg.Wait()
}

// DefendedServer demonstrates basic defenses against each attack type
func DefendedServer() *http.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    return &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,  // Slowloris defense
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      30 * time.Second,
        IdleTimeout:       60 * time.Second,
        MaxHeaderBytes:    1 << 16, // 64KB header limit
    }
}
```

### Real-World Context

- **Slowloris** killed Apache servers at PyCon 2009. nginx and Go servers are both resistant because they use goroutines per connection rather than one thread per connection — but goroutines still consume memory.
- **SYN flood** is handled by the OS with `tcp_syncookies` on Linux. Your Go code needs to limit connections *after* `Accept`, not prevent `Accept` from being called.
- **HTTP flood** is the most dangerous for Go services. Each request gets a goroutine. 100k req/s = 100k goroutines if your upstream is slow.

---

## Day 92 — Rate Limiting Algorithms Deep Dive

| | |
|---|---|
| **Learning objective** | Implement all 5 major rate limiting algorithms; understand their behavior under burst traffic |
| **Key packages** | `sync`, `sync/atomic`, `time`, `math` |
| **Security context** | Correct algorithm choice determines whether DDoS traffic is blocked cleanly or leaks through |

### Key Concepts

| Algorithm | Memory | Burst | Boundary issue | Best for |
|-----------|--------|-------|----------------|----------|
| Fixed window | O(1) | Yes | Double-rate at boundary | Simple API limits |
| Sliding window log | O(n) per key | No | None | Accurate audit |
| Sliding window counter | O(1) | Partial | Slight approximation | DDoS (scalable) |
| Token bucket | O(1) | Yes (up to B) | None | API platforms (Stripe, AWS) |
| Leaky bucket | O(1) | No | None | Traffic smoothing |

### Exercise — All 5 Algorithms

```go
package ratelimit

import (
    "sync"
    "sync/atomic"
    "time"
)

// 1. Fixed Window — fast but vulnerable to boundary bursts
type FixedWindow struct {
    mu        sync.Mutex
    limit     int
    count     int
    windowEnd time.Time
    period    time.Duration
}

func NewFixedWindow(limit int, period time.Duration) *FixedWindow {
    return &FixedWindow{limit: limit, period: period, windowEnd: time.Now().Add(period)}
}

func (fw *FixedWindow) Allow() bool {
    fw.mu.Lock()
    defer fw.mu.Unlock()
    now := time.Now()
    if now.After(fw.windowEnd) {
        fw.count = 0
        fw.windowEnd = now.Add(fw.period)
    }
    if fw.count >= fw.limit {
        return false
    }
    fw.count++
    return true
}

// 2. Sliding Window Log — accurate, memory-heavy (stores each request timestamp)
type SlidingWindowLog struct {
    mu      sync.Mutex
    limit   int
    period  time.Duration
    entries []time.Time
}

func NewSlidingWindowLog(limit int, period time.Duration) *SlidingWindowLog {
    return &SlidingWindowLog{limit: limit, period: period}
}

func (swl *SlidingWindowLog) Allow() bool {
    swl.mu.Lock()
    defer swl.mu.Unlock()
    now := time.Now()
    cutoff := now.Add(-swl.period)

    // Evict old entries
    valid := swl.entries[:0]
    for _, t := range swl.entries {
        if t.After(cutoff) {
            valid = append(valid, t)
        }
    }
    swl.entries = valid

    if len(swl.entries) >= swl.limit {
        return false
    }
    swl.entries = append(swl.entries, now)
    return true
}

// 3. Sliding Window Counter — approximation using two fixed windows
type SlidingWindowCounter struct {
    mu          sync.Mutex
    limit       int
    period      time.Duration
    curCount    int
    prevCount   int
    windowStart time.Time
}

func NewSlidingWindowCounter(limit int, period time.Duration) *SlidingWindowCounter {
    return &SlidingWindowCounter{limit: limit, period: period, windowStart: time.Now()}
}

func (swc *SlidingWindowCounter) Allow() bool {
    swc.mu.Lock()
    defer swc.mu.Unlock()
    now := time.Now()
    elapsed := now.Sub(swc.windowStart)

    if elapsed >= swc.period {
        // Shift windows
        swc.prevCount = swc.curCount
        swc.curCount = 0
        swc.windowStart = now
        elapsed = 0
    }

    // Weighted count: prev window contribution + current window
    ratio := 1.0 - elapsed.Seconds()/swc.period.Seconds()
    estimated := float64(swc.prevCount)*ratio + float64(swc.curCount)
    if int(estimated) >= swc.limit {
        return false
    }
    swc.curCount++
    return true
}

// 4. Token Bucket — allows burst up to capacity B, refills at rate R tokens/sec
type TokenBucket struct {
    mu       sync.Mutex
    tokens   float64
    capacity float64
    rate     float64 // tokens per second
    lastTime time.Time
}

func NewTokenBucket(capacity, rate float64) *TokenBucket {
    return &TokenBucket{
        tokens:   capacity,
        capacity: capacity,
        rate:     rate,
        lastTime: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    now := time.Now()
    elapsed := now.Sub(tb.lastTime).Seconds()
    tb.lastTime = now

    tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)
    if tb.tokens < 1.0 {
        return false
    }
    tb.tokens--
    return true
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// 5. Leaky Bucket — smooth constant outflow regardless of input burst
type LeakyBucket struct {
    capacity  int64 // queue depth
    current   atomic.Int64
    drainRate time.Duration // time between each drain
    done      chan struct{}
}

func NewLeakyBucket(capacity int64, drainRate time.Duration) *LeakyBucket {
    lb := &LeakyBucket{capacity: capacity, drainRate: drainRate, done: make(chan struct{})}
    go lb.drain()
    return lb
}

func (lb *LeakyBucket) drain() {
    ticker := time.NewTicker(lb.drainRate)
    defer ticker.Stop()
    for {
        select {
        case <-lb.done:
            return
        case <-ticker.C:
            for {
                cur := lb.current.Load()
                if cur <= 0 {
                    break
                }
                if lb.current.CompareAndSwap(cur, cur-1) {
                    break
                }
            }
        }
    }
}

func (lb *LeakyBucket) Allow() bool {
    for {
        cur := lb.current.Load()
        if cur >= lb.capacity {
            return false
        }
        if lb.current.CompareAndSwap(cur, cur+1) {
            return true
        }
    }
}

func (lb *LeakyBucket) Stop() { close(lb.done) }
```

### Real-World Context

- AWS API Gateway uses **token bucket** (burst capacity + sustained rate).
- Stripe uses **sliding window counter** per API key.
- DDoS mitigation layers prefer **token bucket** at L4 (per-IP) and **sliding window counter** at L7 (per path/endpoint) because token bucket allows legitimate burst while bounding steady-state rate.

---

## Day 93 — Per-IP Rate Limiting with CIDR Awareness

| | |
|---|---|
| **Learning objective** | Rate limit at both /32 (single IP) and /24 (subnet) granularity; handle IPv6 |
| **Key packages** | `net`, `sync`, `net/http`, `strings` |
| **Security context** | DDoS attackers control IP ranges (ASNs), not just single IPs |

### Key Concepts

- **X-Forwarded-For trust**: Only trust `XFF[len-1]` (last hop, set by your trusted L4 proxy). Never trust the full chain from untrusted clients.
- **IPv6 /48 grouping**: IPv6 bots spread across a /48. Group by first 48 bits.
- **Counter expiry**: Un-evicted stale counters grow memory unbounded. Background goroutine cleans up.

### Exercise — CIDR-Aware Rate Limiter

```go
package ddos

import (
    "fmt"
    "net"
    "net/http"
    "strings"
    "sync"
    "time"
)

type counterEntry struct {
    bucket  *TokenBucket
    lastSee time.Time
}

type CIDRRateLimiter struct {
    mu      sync.RWMutex
    ipLimit  int     // per /32 or /128
    subLimit int     // per /24 or /48
    rate     float64 // tokens/second
    entries  map[string]*counterEntry
    done     chan struct{}
}

func NewCIDRRateLimiter(ipLimit, subLimit int, rate float64) *CIDRRateLimiter {
    rl := &CIDRRateLimiter{
        ipLimit:  ipLimit,
        subLimit: subLimit,
        rate:     rate,
        entries:  make(map[string]*counterEntry),
        done:     make(chan struct{}),
    }
    go rl.cleanup()
    return rl
}

// realIP extracts the real client IP, trusting only one XFF hop
func realIP(r *http.Request) string {
    xff := r.Header.Get("X-Forwarded-For")
    if xff == "" {
        host, _, _ := net.SplitHostPort(r.RemoteAddr)
        return host
    }
    // Trust only the rightmost entry added by our L4 proxy
    parts := strings.Split(xff, ",")
    ip := strings.TrimSpace(parts[len(parts)-1])
    if net.ParseIP(ip) == nil {
        host, _, _ := net.SplitHostPort(r.RemoteAddr)
        return host
    }
    return ip
}

// subnetKey returns the /24 (IPv4) or /48 (IPv6) key
func subnetKey(ipStr string) string {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return ipStr
    }
    if ip4 := ip.To4(); ip4 != nil {
        // IPv4: mask to /24
        return fmt.Sprintf("%d.%d.%d.0/24", ip4[0], ip4[1], ip4[2])
    }
    // IPv6: mask to /48 (first 6 bytes)
    return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x::/48",
        ip[0], ip[1], ip[2], ip[3], ip[4], ip[5])
}

func (rl *CIDRRateLimiter) getBucket(key string, limit int) *TokenBucket {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    e, ok := rl.entries[key]
    if !ok {
        e = &counterEntry{bucket: NewTokenBucket(float64(limit), rl.rate)}
        rl.entries[key] = e
    }
    e.lastSee = time.Now()
    return e.bucket
}

func (rl *CIDRRateLimiter) Allow(r *http.Request) bool {
    ip := realIP(r)
    subnet := subnetKey(ip)

    // Both per-IP and per-subnet must allow
    if !rl.getBucket(ip, rl.ipLimit).Allow() {
        return false
    }
    if !rl.getBucket(subnet, rl.subLimit).Allow() {
        return false
    }
    return true
}

func (rl *CIDRRateLimiter) cleanup() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-rl.done:
            return
        case <-ticker.C:
            cutoff := time.Now().Add(-10 * time.Minute)
            rl.mu.Lock()
            for k, e := range rl.entries {
                if e.lastSee.Before(cutoff) {
                    delete(rl.entries, k)
                }
            }
            rl.mu.Unlock()
        }
    }
}

func (rl *CIDRRateLimiter) Stop() { close(rl.done) }

func (rl *CIDRRateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !rl.Allow(r) {
            w.Header().Set("Retry-After", "1")
            http.Error(w, `{"error":"rate_limit_exceeded"}`, http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

## Day 94 — SYN Flood Defense

| | |
|---|---|
| **Learning objective** | Detect and limit half-open TCP connections per source IP |
| **Key packages** | `net`, `sync`, `sync/atomic`, `time` |
| **Security context** | SYN floods exhaust server connection tables; defense requires limiting half-open state per IP |

### Key Concepts

**SYN flood mechanics:**
1. Attacker sends SYN packets with spoofed source IPs
2. Server allocates connection state, sends SYN-ACK
3. Attacker never sends ACK — server waits (default: 75s on Linux)
4. Server's half-open queue fills; legitimate connections rejected

**Go-side:** The OS handles SYN cookies. Your Go code limits connections *after* `Accept` returns — so you're limiting fully-established connections that connected but haven't sent data yet.

### Exercise — Half-Open Connection Limiter

```go
package ddos

import (
    "net"
    "sync"
    "sync/atomic"
    "time"
)

// HalfOpenTracker tracks connections that connected but haven't sent data
type HalfOpenTracker struct {
    mu          sync.Mutex
    perIP       map[string]int // count of half-open conns per IP
    maxPerIP    int
    drainAfter  time.Duration
    totalHalf   atomic.Int64
    totalBlocked atomic.Int64
}

func NewHalfOpenTracker(maxPerIP int, drainAfter time.Duration) *HalfOpenTracker {
    return &HalfOpenTracker{
        perIP:      make(map[string]int),
        maxPerIP:   maxPerIP,
        drainAfter: drainAfter,
    }
}

// TrackConn returns a monitored connection.
// Call this immediately after net.Listener.Accept().
// If the IP has too many half-open conns, returns false — caller should close conn.
func (t *HalfOpenTracker) TrackConn(conn net.Conn) (net.Conn, bool) {
    host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
    if err != nil {
        return conn, true // can't parse — let through
    }

    t.mu.Lock()
    count := t.perIP[host]
    if count >= t.maxPerIP {
        t.mu.Unlock()
        t.totalBlocked.Add(1)
        return conn, false
    }
    t.perIP[host] = count + 1
    t.mu.Unlock()
    t.totalHalf.Add(1)

    // Wrap the connection; when data arrives, remove from half-open
    return &trackedConn{
        Conn:    conn,
        tracker: t,
        ip:      host,
        timer:   time.AfterFunc(t.drainAfter, func() { t.release(host) }),
    }, true
}

func (t *HalfOpenTracker) release(ip string) {
    t.mu.Lock()
    defer t.mu.Unlock()
    if t.perIP[ip] > 0 {
        t.perIP[ip]--
        if t.perIP[ip] == 0 {
            delete(t.perIP, ip)
        }
    }
    t.totalHalf.Add(-1)
}

func (t *HalfOpenTracker) Stats() (half, blocked int64) {
    return t.totalHalf.Load(), t.totalBlocked.Load()
}

type trackedConn struct {
    net.Conn
    tracker  *HalfOpenTracker
    ip       string
    timer    *time.Timer
    released sync.Once
}

// Read — first data received, connection is no longer half-open
func (tc *trackedConn) Read(b []byte) (int, error) {
    tc.released.Do(func() {
        tc.timer.Stop()
        tc.tracker.release(tc.ip)
    })
    return tc.Conn.Read(b)
}

func (tc *trackedConn) Close() error {
    tc.released.Do(func() {
        tc.timer.Stop()
        tc.tracker.release(tc.ip)
    })
    return tc.Conn.Close()
}

// ProtectedListener wraps net.Listener with SYN flood defense
type ProtectedListener struct {
    net.Listener
    tracker *HalfOpenTracker
}

func NewProtectedListener(l net.Listener, maxHalfOpenPerIP int) *ProtectedListener {
    return &ProtectedListener{
        Listener: l,
        tracker:  NewHalfOpenTracker(maxHalfOpenPerIP, 3*time.Second),
    }
}

func (pl *ProtectedListener) Accept() (net.Conn, error) {
    for {
        conn, err := pl.Listener.Accept()
        if err != nil {
            return nil, err
        }
        tracked, ok := pl.tracker.TrackConn(conn)
        if !ok {
            conn.Close()
            continue // silently drop — don't return error
        }
        return tracked, nil
    }
}
```

---

## Day 95 — Slowloris Detection & Defense

| | |
|---|---|
| **Learning objective** | Detect connections dripping headers below a byte-rate threshold |
| **Key packages** | `net`, `net/http`, `sync`, `sync/atomic`, `time` |
| **Security context** | Slowloris exhausts goroutine pools by holding connections indefinitely |

### Exercise — Byte-Rate Connection Monitor

```go
package ddos

import (
    "net"
    "sync"
    "sync/atomic"
    "time"
)

// ByteRateConn wraps net.Conn and tracks bytes received per second
type ByteRateConn struct {
    net.Conn
    mu           sync.Mutex
    bytesInWindow int
    windowStart  time.Time
    windowSize   time.Duration

    // Slowloris detection thresholds
    minBytesPerSec int
    graceDeadline  time.Time // don't check until this time passes

    slowloris atomic.Bool
}

func NewByteRateConn(conn net.Conn, minBytesPerSec int, gracePeriod time.Duration) *ByteRateConn {
    return &ByteRateConn{
        Conn:           conn,
        windowStart:    time.Now(),
        windowSize:     time.Second,
        minBytesPerSec: minBytesPerSec,
        graceDeadline:  time.Now().Add(gracePeriod),
    }
}

func (c *ByteRateConn) Read(b []byte) (int, error) {
    n, err := c.Conn.Read(b)
    if n > 0 {
        c.mu.Lock()
        now := time.Now()
        if now.Sub(c.windowStart) >= c.windowSize {
            // Check rate for completed window
            if now.After(c.graceDeadline) {
                rate := c.bytesInWindow // bytes in past second
                if rate < c.minBytesPerSec {
                    c.slowloris.Store(true)
                }
            }
            c.bytesInWindow = n
            c.windowStart = now
        } else {
            c.bytesInWindow += n
        }
        c.mu.Unlock()
    }
    return n, err
}

func (c *ByteRateConn) IsSlowloris() bool { return c.slowloris.Load() }

// SlowlorisDetector wraps a listener, monitors byte rates, and closes suspicious conns
type SlowlorisDetector struct {
    net.Listener
    minBytesPerSec int
    gracePeriod    time.Duration
    checkInterval  time.Duration
    blockList      *ExpirableBlockList // reuse from Day 97

    mu    sync.Mutex
    conns map[*ByteRateConn]struct{}
    done  chan struct{}

    detected atomic.Int64
}

func NewSlowlorisDetector(l net.Listener, minBytesPerSec int) *SlowlorisDetector {
    sd := &SlowlorisDetector{
        Listener:       l,
        minBytesPerSec: minBytesPerSec,
        gracePeriod:    5 * time.Second,
        checkInterval:  2 * time.Second,
        conns:          make(map[*ByteRateConn]struct{}),
        done:           make(chan struct{}),
    }
    go sd.monitor()
    return sd
}

func (sd *SlowlorisDetector) Accept() (net.Conn, error) {
    conn, err := sd.Listener.Accept()
    if err != nil {
        return nil, err
    }
    monitored := NewByteRateConn(conn, sd.minBytesPerSec, sd.gracePeriod)
    sd.mu.Lock()
    sd.conns[monitored] = struct{}{}
    sd.mu.Unlock()
    return monitored, nil
}

func (sd *SlowlorisDetector) monitor() {
    ticker := time.NewTicker(sd.checkInterval)
    defer ticker.Stop()
    for {
        select {
        case <-sd.done:
            return
        case <-ticker.C:
            sd.mu.Lock()
            for c := range sd.conns {
                if c.IsSlowloris() {
                    sd.detected.Add(1)
                    ip, _, _ := net.SplitHostPort(c.RemoteAddr().String())
                    _ = ip // add to block list in production
                    c.Close()
                    delete(sd.conns, c)
                }
            }
            sd.mu.Unlock()
        }
    }
}

func (sd *SlowlorisDetector) removeConn(c *ByteRateConn) {
    sd.mu.Lock()
    delete(sd.conns, c)
    sd.mu.Unlock()
}
```

---

## Day 96 — Traffic Anomaly Detection (Z-Score + EMA)

| | |
|---|---|
| **Learning objective** | Flag IPs whose request rate deviates significantly from baseline using statistical methods |
| **Key packages** | `sync`, `math`, `time`, `net/http` |
| **Security context** | Static rate limits block burst-then-normal attacks poorly; statistical detection adapts to traffic patterns |

### Key Concepts

**Exponential Moving Average (EMA):**
- `EMA_new = alpha * value + (1 - alpha) * EMA_old`
- `alpha = 0.1` means recent values matter but history is preserved
- Tracks mean and variance without storing all historical data

**Z-score:** `(current_rate - mean) / stddev` — IPs > 3 sigma above mean are anomalous.

### Exercise — Statistical Traffic Analyzer

```go
package ddos

import (
    "math"
    "net/http"
    "sync"
    "time"
)

type ipStats struct {
    mu       sync.Mutex
    requests []time.Time // sliding window entries (last 60s)
    rate     float64     // current request rate (req/s)
}

func (s *ipStats) record(now time.Time, window time.Duration) float64 {
    s.mu.Lock()
    defer s.mu.Unlock()
    cutoff := now.Add(-window)
    valid := s.requests[:0]
    for _, t := range s.requests {
        if t.After(cutoff) {
            valid = append(valid, t)
        }
    }
    s.requests = append(valid, now)
    s.rate = float64(len(s.requests)) / window.Seconds()
    return s.rate
}

type TrafficAnalyzer struct {
    mu      sync.RWMutex
    ipStats map[string]*ipStats
    window  time.Duration

    // EMA-based global baseline
    emaMean float64
    emaVar  float64
    alpha   float64

    Threshold float64 // z-score threshold (default 3.0)
    Alerts    chan Alert
    done      chan struct{}
}

type Alert struct {
    IP        string
    Rate      float64
    ZScore    float64
    Timestamp time.Time
}

func NewTrafficAnalyzer() *TrafficAnalyzer {
    ta := &TrafficAnalyzer{
        ipStats:   make(map[string]*ipStats),
        window:    60 * time.Second,
        alpha:     0.1,
        Threshold: 3.0,
        Alerts:    make(chan Alert, 100),
        done:      make(chan struct{}),
    }
    go ta.cleanup()
    return ta
}

func (ta *TrafficAnalyzer) Record(ip string) {
    ta.mu.Lock()
    s, ok := ta.ipStats[ip]
    if !ok {
        s = &ipStats{}
        ta.ipStats[ip] = s
    }
    ta.mu.Unlock()

    now := time.Now()
    rate := s.record(now, ta.window)

    // Update global EMA
    ta.mu.Lock()
    diff := rate - ta.emaMean
    ta.emaMean += ta.alpha * diff
    ta.emaVar = (1-ta.alpha) * (ta.emaVar + ta.alpha*diff*diff)
    mean := ta.emaMean
    variance := ta.emaVar
    ta.mu.Unlock()

    // Compute z-score
    stddev := math.Sqrt(variance)
    if stddev < 0.01 {
        return // not enough data
    }
    z := (rate - mean) / stddev
    if z > ta.Threshold {
        select {
        case ta.Alerts <- Alert{IP: ip, Rate: rate, ZScore: z, Timestamp: now}:
        default:
        }
    }
}

func (ta *TrafficAnalyzer) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip, _, _ := net.SplitHostPort(r.RemoteAddr)
        ta.Record(ip)
        next.ServeHTTP(w, r)
    })
}

func (ta *TrafficAnalyzer) cleanup() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ta.done:
            return
        case <-ticker.C:
            ta.mu.Lock()
            // Remove IPs with no recent activity (empty window)
            for ip, s := range ta.ipStats {
                s.mu.Lock()
                cutoff := time.Now().Add(-ta.window)
                fresh := false
                for _, t := range s.requests {
                    if t.After(cutoff) {
                        fresh = true
                        break
                    }
                }
                s.mu.Unlock()
                if !fresh {
                    delete(ta.ipStats, ip)
                }
            }
            ta.mu.Unlock()
        }
    }
}
```

---

## Day 97 — IP Blocking & Proof-of-Work Challenges

| | |
|---|---|
| **Learning objective** | Implement tiered response: PoW challenge → block with auto-expiry |
| **Key packages** | `crypto/sha256`, `crypto/rand`, `crypto/hmac`, `sync`, `net/http`, `encoding/hex` |
| **Security context** | Hard blocks cause false positives; PoW imposes cost on bots while letting legitimate traffic through |

### Key Concepts

**Proof-of-Work:** Server issues `{challenge, difficulty}`. Client must find `nonce` such that `SHA256(challenge + nonce)` has `difficulty` leading zero bits. Cost scales exponentially with difficulty — difficulty=16 requires ~65536 attempts on average.

**Tiered response:** Suspected (score 50-80) → 429 + PoW; Confirmed (score >80) → 403 block.

### Exercise — PoW Challenge System + IP Blocklist

```go
package ddos

import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

// ExpirableBlockList is a thread-safe IP blocklist with auto-expiry
type ExpirableBlockList struct {
    mu      sync.RWMutex
    entries map[string]time.Time
    done    chan struct{}
}

func NewBlockList() *ExpirableBlockList {
    bl := &ExpirableBlockList{
        entries: make(map[string]time.Time),
        done:    make(chan struct{}),
    }
    go bl.cleanup()
    return bl
}

func (bl *ExpirableBlockList) Block(ip string, duration time.Duration) {
    bl.mu.Lock()
    bl.entries[ip] = time.Now().Add(duration)
    bl.mu.Unlock()
}

func (bl *ExpirableBlockList) IsBlocked(ip string) bool {
    bl.mu.RLock()
    exp, ok := bl.entries[ip]
    bl.mu.RUnlock()
    return ok && time.Now().Before(exp)
}

func (bl *ExpirableBlockList) cleanup() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-bl.done:
            return
        case <-ticker.C:
            now := time.Now()
            bl.mu.Lock()
            for ip, exp := range bl.entries {
                if now.After(exp) {
                    delete(bl.entries, ip)
                }
            }
            bl.mu.Unlock()
        }
    }
}

// PoW challenge system
type PowChallenger struct {
    secret     []byte
    difficulty int // leading zero bits required
    expiry     time.Duration
}

func NewPowChallenger(secret []byte, difficulty int) *PowChallenger {
    return &PowChallenger{secret: secret, difficulty: difficulty, expiry: 5 * time.Minute}
}

type Challenge struct {
    Challenge  string `json:"challenge"`
    Difficulty int    `json:"difficulty"`
    Expires    int64  `json:"expires"`
}

func (pc *PowChallenger) Issue() Challenge {
    raw := make([]byte, 16)
    rand.Read(raw)
    expires := time.Now().Add(pc.expiry).Unix()

    // Sign challenge + expiry with HMAC so server can verify it issued this challenge
    mac := hmac.New(sha256.New, pc.secret)
    mac.Write(raw)
    binary.Write(mac, binary.BigEndian, expires)
    sig := mac.Sum(nil)

    payload := append(raw, sig[:8]...) // compact: 16 + 8 bytes
    return Challenge{
        Challenge:  hex.EncodeToString(payload),
        Difficulty: pc.difficulty,
        Expires:    expires,
    }
}

type SolveRequest struct {
    Challenge string `json:"challenge"`
    Nonce     string `json:"nonce"`
}

func (pc *PowChallenger) Verify(sr SolveRequest) bool {
    payload, err := hex.DecodeString(sr.Challenge)
    if err != nil || len(payload) < 24 {
        return false
    }
    raw := payload[:16]
    storedSig := payload[16:24]

    // Reconstruct expires — in production, embed it in the challenge payload
    // For simplicity here we skip expiry re-extraction (demonstrate concept)
    mac := hmac.New(sha256.New, pc.secret)
    mac.Write(raw)
    expectedSig := mac.Sum(nil)[:8]
    if !hmac.Equal(storedSig, expectedSig) {
        return false
    }

    nonce, err := hex.DecodeString(sr.Nonce)
    if err != nil {
        return false
    }

    // Check SHA256(challenge_raw + nonce) has `difficulty` leading zero bits
    h := sha256.New()
    h.Write(raw)
    h.Write(nonce)
    hash := h.Sum(nil)
    return leadingZeroBits(hash) >= pc.difficulty
}

func leadingZeroBits(b []byte) int {
    count := 0
    for _, by := range b {
        if by == 0 {
            count += 8
            continue
        }
        for mask := byte(0x80); mask > 0; mask >>= 1 {
            if by&mask != 0 {
                return count
            }
            count++
        }
        break
    }
    return count
}

// TieredResponder applies challenge or block based on score
type TieredResponder struct {
    blockList  *ExpirableBlockList
    challenger *PowChallenger
    whitelist  []string // trusted CIDR prefixes
}

func (tr *TieredResponder) Handle(score float64, ip string, w http.ResponseWriter, r *http.Request, next http.Handler) {
    if tr.isTrusted(ip) {
        next.ServeHTTP(w, r)
        return
    }
    if tr.blockList.IsBlocked(ip) {
        http.Error(w, `{"error":"blocked"}`, http.StatusForbidden)
        return
    }
    switch {
    case score > 80:
        tr.blockList.Block(ip, 5*time.Minute)
        http.Error(w, `{"error":"blocked"}`, http.StatusForbidden)
    case score > 50:
        // Issue PoW challenge
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusTooManyRequests)
        json.NewEncoder(w).Encode(tr.challenger.Issue())
    default:
        next.ServeHTTP(w, r)
    }
}

func (tr *TieredResponder) isTrusted(ip string) bool {
    for _, cidr := range tr.whitelist {
        _, network, err := net.ParseCIDR(cidr)
        if err != nil {
            continue
        }
        if network.Contains(net.ParseIP(ip)) {
            return true
        }
    }
    return false
}
```

---

## Day 98 — UDP Amplification Defense

| | |
|---|---|
| **Learning objective** | Build a UDP service that cannot be used as an amplification vector |
| **Key packages** | `net`, `sync`, `time`, `sync/atomic` |
| **Security context** | Amplification factor = response_bytes / request_bytes; services with factor >1 can be weaponized |

### Key Concepts

**Amplification attack:** Attacker sends small UDP packet with spoofed source IP (victim). Service replies with large response to victim. Factor 10x = 10Gbps attack from 1Gbps attacker.

**Defenses:**
1. **Response size cap**: Never send more bytes than received × 1.5
2. **Per-source rate limiting**: Max 10 UDP responses per second per source IP
3. **Challenge-response**: Require a handshake before sending large responses

### Exercise — Anti-Amplification UDP Server

```go
package ddos

import (
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "net"
    "sync"
    "sync/atomic"
    "time"
)

type udpRateEntry struct {
    mu      sync.Mutex
    tokens  float64
    lastSee time.Time
    blocked atomic.Int32
}

// AntiAmpUDPServer is a UDP server that enforces anti-amplification rules
type AntiAmpUDPServer struct {
    conn        *net.UDPConn
    maxFactor   float64 // max response_size / request_size ratio (e.g. 1.5)
    ratePerIP   float64 // max responses per second per source IP
    mu          sync.Mutex
    ipBuckets   map[string]*udpRateEntry
    blocked     map[string]time.Time
    ampAttempts atomic.Int64
}

func NewAntiAmpUDPServer(addr string, maxFactor, ratePerIP float64) (*AntiAmpUDPServer, error) {
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return nil, err
    }
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return nil, err
    }
    s := &AntiAmpUDPServer{
        conn:      conn,
        maxFactor: maxFactor,
        ratePerIP: ratePerIP,
        ipBuckets: make(map[string]*udpRateEntry),
        blocked:   make(map[string]time.Time),
    }
    return s, nil
}

func (s *AntiAmpUDPServer) Serve() {
    buf := make([]byte, 1500) // MTU-sized buffer
    for {
        n, src, err := s.conn.ReadFromUDP(buf)
        if err != nil {
            return
        }
        reqBytes := n
        ip := src.IP.String()

        if s.isBlocked(ip) {
            continue // silently drop — no amplification
        }
        if !s.allowRate(ip) {
            s.checkBlock(ip)
            continue
        }

        // Process request and generate response
        response := s.process(buf[:n])

        // Anti-amplification: cap response to maxFactor × request size
        maxRespSize := int(float64(reqBytes) * s.maxFactor)
        if len(response) > maxRespSize {
            s.ampAttempts.Add(1)
            response = response[:maxRespSize] // truncate
        }

        s.conn.WriteToUDP(response, src)
    }
}

func (s *AntiAmpUDPServer) process(req []byte) []byte {
    // Example: echo service with SHA256 response header (simulated work)
    h := sha256.Sum256(req)
    resp := make([]byte, 4+len(req))
    binary.BigEndian.PutUint32(resp[:4], uint32(len(req)))
    copy(resp[4:], h[:])
    return resp
}

func (s *AntiAmpUDPServer) allowRate(ip string) bool {
    s.mu.Lock()
    entry, ok := s.ipBuckets[ip]
    if !ok {
        entry = &udpRateEntry{tokens: s.ratePerIP, lastSee: time.Now()}
        s.ipBuckets[ip] = entry
    }
    s.mu.Unlock()

    entry.mu.Lock()
    defer entry.mu.Unlock()
    now := time.Now()
    elapsed := now.Sub(entry.lastSee).Seconds()
    entry.lastSee = now
    entry.tokens = minF(s.ratePerIP, entry.tokens+elapsed*s.ratePerIP)
    if entry.tokens < 1 {
        return false
    }
    entry.tokens--
    return true
}

func (s *AntiAmpUDPServer) checkBlock(ip string) {
    entry, ok := s.ipBuckets[ip]
    if !ok {
        return
    }
    entry.blocked.Add(1)
    if entry.blocked.Load() >= 2 {
        s.mu.Lock()
        s.blocked[ip] = time.Now().Add(5 * time.Minute)
        s.mu.Unlock()
    }
}

func (s *AntiAmpUDPServer) isBlocked(ip string) bool {
    s.mu.Lock()
    exp, ok := s.blocked[ip]
    s.mu.Unlock()
    return ok && time.Now().Before(exp)
}

func minF(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

func (s *AntiAmpUDPServer) Stats() string {
    return fmt.Sprintf("amp_attempts=%d", s.ampAttempts.Load())
}
```

---

## Week 13 Mini Project: DDoS Mitigation Layer

**Build a complete DDoS mitigation proxy** that sits in front of any HTTP backend.

### Architecture

```
Internet
    │
    ▼
[TCP Listener — ProtectedListener (SYN flood defense)]
    │
    ▼
[SlowlorisDetector — byte rate monitor per connection]
    │
    ▼
[CIDRRateLimiter — per-IP + per-/24 token bucket]
    │
    ▼
[TrafficAnalyzer — EMA z-score anomaly detection]
    │
    ▼
[TieredResponder — score → allow / PoW challenge / block]
    │
    ▼
[ThreatIntel — aggregate scores from all modules]
    │
    ▼
[HTTP Reverse Proxy → Backend]
```

### Feature Checklist

- [ ] Per-IP token bucket (100 req/s, burst 200)
- [ ] Per-/24 subnet token bucket (500 req/s, burst 1000)
- [ ] IPv6 /48 grouping
- [ ] Slowloris detection (< 50 bytes/sec during header phase)
- [ ] SYN flood defense (max 5 half-open per IP, 3s drain)
- [ ] Traffic anomaly (z-score > 3.0 → flag)
- [ ] PoW challenge at score 50–80 (SHA256, difficulty=16)
- [ ] Block at score > 80 (5-minute expiry)
- [ ] CIDR whitelist (bypass all mitigations for trusted ranges)
- [ ] `/debug/ddos` endpoint: JSON stats (top IPs, block list, attack status, request rates)
- [ ] Integration with `ThreatIntel` store (signals feed back into score)
- [ ] Anti-amplification UDP service (response ≤ 1.5× request)

### `/debug/ddos` Response Format

```json
{
  "blocked_ips": ["1.2.3.4", "5.6.7.0/24"],
  "top_request_rates": [
    {"ip": "1.2.3.4", "rate_per_sec": 847.2, "z_score": 4.1}
  ],
  "slowloris_detected": 3,
  "syn_flood_blocked": 12,
  "pow_challenges_issued": 47,
  "pow_challenges_solved": 2,
  "amp_attempts_blocked": 0,
  "current_block_list_size": 8
}
```

### Key Packages

`net`, `net/http`, `sync`, `sync/atomic`, `crypto/sha256`, `crypto/hmac`, `crypto/rand`, `math`, `time`

### Engineer Takeaway

Cloudflare Magic Transit, AWS Shield Advanced, and Akamai Prolexic implement this logic at 100 Tbps scale. Their edge nodes run rate limiting at line rate using eBPF/XDP. The core algorithms — token bucket, z-score anomaly, PoW — are identical to what you've built. The difference is hardware, not ideas.

---

## Week 13 Review

| Day | Skill | Real-World Use |
|-----|-------|----------------|
| 91 | DDoS taxonomy + attack simulators | Know what "we're under DDoS" actually means |
| 92 | All 5 rate limiting algorithms | Choose the right algorithm for the right context |
| 93 | CIDR-aware per-IP/subnet rate limiting | Block entire ASNs with one rule |
| 94 | SYN flood defense (half-open tracking) | Understand why `tcp_syncookies` isn't enough |
| 95 | Slowloris detection (byte rate monitor) | Defend against single-machine takedown attacks |
| 96 | Z-score + EMA traffic anomaly detection | Adapt thresholds to traffic patterns automatically |
| 97 | PoW challenges + IP blocklist with expiry | Impose cost on bots without blocking legit users |
| 98 | UDP anti-amplification | Prevent your service from being used as a weapon |
