# Week 5: Networking Concepts + TCP Basics (Days 31–37)

**Phase:** 2 — Networking Foundations + L4
**Goal:** Move from application code to the network. Build a TCP server, understand the OSI model from code, and implement the bidirectional proxy that is the foundation of every L4 load balancer and security gateway.

---

> **Background in this wiki:** the [OSI guide](../../networking/real-life-example-osi.md) (one request through all seven layers); TCP/IP guide [Ch 19 (ports and sockets)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-19-ports-and-sockets-which-program-gets-the-data), [Ch 20 (UDP)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-20-udp-fire-and-forget), [Ch 21–23 (handshake, reliability, closing)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake) with the `dialcheck`, `stream`, `udpping`, and `closewait` labs; Linux guide [Ch 73 (goroutines vs threads)](../../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime) and [Ch 78 (file descriptors and the netpoller)](../../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections).

## Day 31 — The OSI Model in Practice

### Learning Objectives
- Map Go's `net` package to OSI layers
- Trace a complete request from DNS to data delivery
- Understand where your code operates and where the kernel handles things

### OSI Model — Where Go Code Lives

```
Layer 7 — Application   → net/http, WebSocket, gRPC
Layer 6 — Presentation  → encoding/json, crypto/tls (TLS record layer)
Layer 5 — Session       → TLS handshake, HTTP/2 stream management
Layer 4 — Transport     → net.Conn (TCP), net.PacketConn (UDP)
Layer 3 — Network       → net.IP, net.IPNet (Go reads; kernel routes)
Layer 2 — Data Link     → Kernel handles; Go doesn't touch
Layer 1 — Physical      → Hardware; Go doesn't touch

Your proxy sits at L4-L7:
  Client → [L4 TCP] → Proxy → [L7 HTTP] → Backend
```

### Exercise: Full Connection Trace

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "time"
)

func traceConnection(host string, port string) {
    target := host + ":" + port

    // Step 1: DNS resolution
    t0 := time.Now()
    addrs, err := net.LookupHost(host)
    dnsLatency := time.Since(t0)
    if err != nil {
        fmt.Printf("DNS FAIL: %v\n", err)
        return
    }
    fmt.Printf("[DNS]  %s → %v (%v)\n", host, addrs, dnsLatency)

    // Step 2: TCP connection (3-way handshake)
    t1 := time.Now()
    conn, err := net.DialTimeout("tcp", target, 5*time.Second)
    tcpLatency := time.Since(t1)
    if err != nil {
        fmt.Printf("TCP FAIL: %v\n", err)
        return
    }
    defer conn.Close()
    fmt.Printf("[TCP]  connected to %s (%v) — local=%s\n", conn.RemoteAddr(), tcpLatency, conn.LocalAddr())

    // Step 3: Raw HTTP/1.0 request (L7 — application protocol over TCP)
    t2 := time.Now()
    fmt.Fprintf(conn, "GET / HTTP/1.0\r\nHost: %s\r\nConnection: close\r\n\r\n", host)

    // Step 4: Read response
    reader := bufio.NewReader(conn)
    statusLine, _ := reader.ReadString('\n')
    httpLatency := time.Since(t2)
    fmt.Printf("[HTTP] response: %s", statusLine)
    fmt.Printf("[HTTP] first-byte latency: %v\n", httpLatency)
    fmt.Printf("[TOTAL] %.1fms DNS + %.1fms TCP + %.1fms TTFB\n",
        float64(dnsLatency.Microseconds())/1000,
        float64(tcpLatency.Microseconds())/1000,
        float64(httpLatency.Microseconds())/1000,
    )
}

func main() {
    fmt.Println("=== Tracing connection to example.com ===")
    traceConnection("example.com", "80")
}
```

### Real-World Context
DNS latency is often the biggest surprise in production. A 50ms DNS lookup before every new connection = 50ms added to cold-start latency. This is why your load balancer (Week 7) caches DNS results with TTL. "Why is our p99 latency 300ms?" is often "why does DNS resolution take 300ms for 1% of requests?"

---

## Day 32 — TCP Fundamentals

### Learning Objectives
- Build a TCP echo server and observe connection states
- Understand the 3-way handshake and why it matters for SYN floods
- Use `netstat` / `ss` to observe socket states during development

### Key Concepts

```
TCP 3-Way Handshake:
  Client → SYN(seq=1000)        → Server  [SYN_SENT / SYN_RCVD]
  Client ← SYN-ACK(seq=2000)   ← Server
  Client → ACK(seq=2001)        → Server  [ESTABLISHED]

SYN Flood Attack:
  Attacker sends millions of SYNs without completing the handshake
  Server allocates state for each SYN → half-open connection queue fills
  Legitimate connections rejected → denial of service
  Defense: SYN cookies (kernel), connection rate limiting (your code, Day 94)

Connection States after Close:
  TIME_WAIT (2*MSL ≈ 60-120s) — normal, server-side after graceful close
  CLOSE_WAIT — application bug: peer closed but you haven't called conn.Close()
```

### Exercise: TCP Echo Server

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "sync/atomic"
    "time"
)

type TCPServer struct {
    listener net.Listener
    active   atomic.Int64
    total    atomic.Int64
}

func NewTCPServer(addr string) (*TCPServer, error) {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }
    return &TCPServer{listener: ln}, nil
}

func (s *TCPServer) Serve() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            if ne, ok := err.(net.Error); ok && !ne.Temporary() {
                return // listener closed
            }
            continue
        }
        s.total.Add(1)
        s.active.Add(1)
        go s.handleConn(conn)
    }
}

func (s *TCPServer) handleConn(conn net.Conn) {
    defer conn.Close()
    defer s.active.Add(-1)

    remote := conn.RemoteAddr()
    local := conn.LocalAddr()
    fmt.Printf("NEW  conn %s → %s\n", remote, local)

    // 30-second idle timeout — prevent zombie connections
    conn.SetDeadline(time.Now().Add(30 * time.Second))

    scanner := bufio.NewScanner(conn)
    var bytesExchanged int64

    for scanner.Scan() {
        line := scanner.Text()
        // Reset deadline on each message (sliding timeout)
        conn.SetDeadline(time.Now().Add(30 * time.Second))

        response := "ECHO: " + line + "\n"
        n, err := fmt.Fprint(conn, response)
        bytesExchanged += int64(len(line)) + int64(n)

        if err != nil {
            fmt.Printf("WRITE err %s: %v\n", remote, err)
            return
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("READ err %s: %v\n", remote, err)
    }
    fmt.Printf("CLOSE conn %s — %d bytes exchanged\n", remote, bytesExchanged)
}

func (s *TCPServer) Stats() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        fmt.Printf("STATS active=%d total=%d\n", s.active.Load(), s.total.Load())
    }
}

func main() {
    srv, err := NewTCPServer(":9090")
    if err != nil {
        fmt.Printf("listen error: %v\n", err)
        return
    }
    fmt.Println("TCP echo server on :9090")

    go srv.Stats()
    srv.Serve()
}
```

**Observe socket states:**
```bash
# In another terminal while server is running
netstat -an | grep 9090
# tcp4   0   0  *.9090   *.*     LISTEN
# tcp4   0   0  127.0.0.1.9090  127.0.0.1.50123  ESTABLISHED

# Connect multiple clients
for i in $(seq 1 10); do nc 127.0.0.1 9090 & done
# After killing nc clients, observe TIME_WAIT states
netstat -an | grep 9090 | grep TIME_WAIT
```

### Real-World Context
`TIME_WAIT` states are normal and harmless. If you see thousands of `CLOSE_WAIT` states, that's an application bug — your code accepted connections but didn't call `conn.Close()`. In a security proxy handling 10k connections/second, even a small leak is fatal.

---

## Day 33 — TCP Server Patterns

### Learning Objectives
- Implement connection limits to protect against resource exhaustion
- Set deadlines on connections to prevent slowloris-style attacks
- Track connection states for observability and incident response

### Exercise: Hardened TCP Server with Connection Limits

```go
package main

import (
    "fmt"
    "net"
    "sync/atomic"
    "time"
)

const (
    maxConnections = 1000
    idleTimeout    = 30 * time.Second
    readTimeout    = 60 * time.Second
)

type HardenedServer struct {
    listener   net.Listener
    active     atomic.Int64
    total      atomic.Int64
    rejected   atomic.Int64
}

func (s *HardenedServer) Serve() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return
        }

        // Connection limit check — reject immediately if at capacity
        current := s.active.Load()
        if current >= maxConnections {
            s.rejected.Add(1)
            fmt.Fprintf(conn, "ERR: server at capacity (%d/%d connections)\n", current, maxConnections)
            conn.Close()
            continue
        }

        s.active.Add(1)
        s.total.Add(1)
        go s.handleConn(conn)
    }
}

func (s *HardenedServer) handleConn(conn net.Conn) {
    defer conn.Close()
    defer s.active.Add(-1)

    remote := conn.RemoteAddr()

    // Set absolute deadline (protects against very slow clients)
    conn.SetDeadline(time.Now().Add(readTimeout))

    buf := make([]byte, 4096)
    for {
        // Reset idle timer on each read
        conn.SetReadDeadline(time.Now().Add(idleTimeout))

        n, err := conn.Read(buf)
        if err != nil {
            if ne, ok := err.(net.Error); ok && ne.Timeout() {
                fmt.Printf("IDLE_TIMEOUT %s\n", remote)
            }
            return
        }

        // Echo data back
        conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
        if _, err := conn.Write(buf[:n]); err != nil {
            return
        }
    }
}
```

### Real-World Context
Without `maxConnections`, a DDoS opening 1 million TCP connections would exhaust file descriptors (`ulimit -n`) and kill your server. Without `idleTimeout`, a Slowloris attack keeping 10,000 connections open but never sending data would do the same. Both are OWASP A05 (Security Misconfiguration).

---

## Day 34 — UDP

### Learning Objectives
- Build a statsd-compatible metrics receiver — understand UDP's fire-and-forget nature
- Understand why UDP is used for DDoS amplification attacks
- Handle packet loss gracefully

### Key Concepts

```
TCP vs UDP:
  TCP: connection-oriented, ordered, reliable, flow-controlled
  UDP: connectionless, unordered, unreliable, no flow control

Why UDP for metrics (statsd, syslog):
  - Sending metrics should never block the application
  - Lost metrics are acceptable; blocking is not
  - 10k packets/sec is normal; TCP overhead would be significant

Why UDP enables DDoS amplification:
  - Attacker spoofs source IP (UDP has no handshake to verify)
  - Sends small request to DNS/NTP server
  - Server sends large response to victim (the spoofed IP)
  - Amplification factor: DNS up to 70x, NTP up to 556x
```

### Exercise: statsd-Compatible Metrics Receiver

```go
package main

import (
    "fmt"
    "net"
    "strconv"
    "strings"
    "sync"
    "time"
)

type Metric struct {
    Name  string
    Value float64
    Type  string // c=counter, g=gauge, ms=timer
}

type StatsReceiver struct {
    conn     *net.UDPConn
    mu       sync.Mutex
    counters map[string]float64
    gauges   map[string]float64
    timers   map[string][]float64
    received int64
    dropped  int64
}

func NewStatsReceiver(addr string) (*StatsReceiver, error) {
    udpAddr, err := net.ResolveUDPAddr("udp", addr)
    if err != nil {
        return nil, err
    }
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return nil, err
    }
    // Set receive buffer — important for burst handling
    conn.SetReadBuffer(1024 * 1024) // 1MB
    return &StatsReceiver{
        conn:     conn,
        counters: make(map[string]float64),
        gauges:   make(map[string]float64),
        timers:   make(map[string][]float64),
    }, nil
}

// parseMetric parses "metric.name:value|type" (statsd format)
func parseMetric(packet string) (*Metric, error) {
    // Format: name:value|type[@sample_rate]
    colonIdx := strings.LastIndex(packet, ":")
    pipeIdx := strings.LastIndex(packet, "|")
    if colonIdx < 0 || pipeIdx < colonIdx {
        return nil, fmt.Errorf("invalid format: %q", packet)
    }

    name := packet[:colonIdx]
    valueStr := packet[colonIdx+1 : pipeIdx]
    metricType := strings.Split(packet[pipeIdx+1:], "@")[0]

    value, err := strconv.ParseFloat(valueStr, 64)
    if err != nil {
        return nil, fmt.Errorf("invalid value %q", valueStr)
    }

    return &Metric{Name: name, Value: value, Type: metricType}, nil
}

func (sr *StatsReceiver) Receive() {
    buf := make([]byte, 65507) // max UDP payload
    for {
        n, addr, err := sr.conn.ReadFromUDP(buf)
        if err != nil {
            return
        }
        _ = addr // could use for per-IP rate limiting

        // Parse potentially multiple metrics (newline-separated)
        for _, line := range strings.Split(string(buf[:n]), "\n") {
            line = strings.TrimSpace(line)
            if line == "" { continue }

            m, err := parseMetric(line)
            if err != nil {
                sr.dropped++
                continue
            }
            sr.record(m)
            sr.received++
        }
    }
}

func (sr *StatsReceiver) record(m *Metric) {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    switch m.Type {
    case "c":
        sr.counters[m.Name] += m.Value
    case "g":
        sr.gauges[m.Name] = m.Value
    case "ms":
        sr.timers[m.Name] = append(sr.timers[m.Name], m.Value)
    }
}

func (sr *StatsReceiver) Summary() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        sr.mu.Lock()
        fmt.Printf("=== Stats (received=%d dropped=%d) ===\n", sr.received, sr.dropped)
        for k, v := range sr.counters {
            fmt.Printf("  counter %s = %.0f\n", k, v)
        }
        for k, v := range sr.gauges {
            fmt.Printf("  gauge   %s = %.2f\n", k, v)
        }
        sr.mu.Unlock()
    }
}

func main() {
    recv, err := NewStatsReceiver(":8125")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    fmt.Println("statsd receiver on UDP :8125")
    go recv.Summary()
    recv.Receive()
}
```

**Send test metrics:**
```bash
# Send a counter
echo -n "proxy.requests.total:1|c" | nc -u -w1 127.0.0.1 8125

# Send a gauge
echo -n "proxy.connections.active:42|g" | nc -u -w1 127.0.0.1 8125

# Send a timer
echo -n "proxy.request.latency:45.3|ms" | nc -u -w1 127.0.0.1 8125
```

### Real-World Context
DDoS amplification via UDP: attacker sends 100-byte DNS query with spoofed source IP → DNS server sends 7000-byte response to victim. Your security stack defends against this by rate limiting UDP responses per source IP (Day 98). Understanding the attack vector is why you're building this from scratch.

---

## Day 35 — I/O Patterns for Networking

### Learning Objectives
- Implement bidirectional TCP proxying — the core of every proxy and load balancer
- Use `io.TeeReader` to tap traffic without breaking the data stream
- Understand zero-copy concepts for high-throughput proxies

### Exercise: TCP Traffic Mirror

```go
package main

import (
    "fmt"
    "io"
    "net"
    "os"
    "sync"
    "time"
)

// mirrorConn proxies client↔upstream and tees all traffic to a tap writer
func mirrorConn(client net.Conn, upstreamAddr string, tap io.Writer) {
    defer client.Close()

    // Connect to upstream
    upstream, err := net.DialTimeout("tcp", upstreamAddr, 5*time.Second)
    if err != nil {
        fmt.Fprintf(os.Stderr, "upstream connect failed: %v\n", err)
        return
    }
    defer upstream.Close()

    remote := client.RemoteAddr().String()
    fmt.Printf("PROXY %s → %s\n", remote, upstreamAddr)

    var wg sync.WaitGroup
    wg.Add(2)

    // client → upstream (tee to tap)
    go func() {
        defer wg.Done()
        // TeeReader: every byte read from client is also written to tap
        tee := io.TeeReader(client, prefixWriter(tap, fmt.Sprintf("[%s→UP] ", remote)))
        n, err := io.Copy(upstream, tee)
        fmt.Printf("SENT %s: %d bytes (err: %v)\n", remote, n, err)
        upstream.(*net.TCPConn).CloseWrite() // signal EOF to upstream
    }()

    // upstream → client (tee to tap)
    go func() {
        defer wg.Done()
        tee := io.TeeReader(upstream, prefixWriter(tap, fmt.Sprintf("[UP→%s] ", remote)))
        n, err := io.Copy(client, tee)
        fmt.Printf("RECV %s: %d bytes (err: %v)\n", remote, n, err)
        client.(*net.TCPConn).CloseWrite()
    }()

    wg.Wait()
    fmt.Printf("CLOSE %s\n", remote)
}

// prefixWriter adds a prefix to each Write call
type prefixWriter struct {
    w      io.Writer
    prefix string
}

func (pw prefixWriter) Write(p []byte) (int, error) {
    fmt.Fprintf(pw.w, "%s%d bytes: %q\n", pw.prefix, len(p), p)
    return len(p), nil
}

func main() {
    // Tap: write to file
    tapFile, err := os.Create("traffic-tap.log")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    defer tapFile.Close()

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    fmt.Println("Mirror proxy on :8080 → upstream :9090")
    fmt.Println("Traffic logged to traffic-tap.log")

    for {
        conn, err := ln.Accept()
        if err != nil { return }
        go mirrorConn(conn, "127.0.0.1:9090", tapFile)
    }
}
```

### Real-World Context
`io.TeeReader` is used by:
- **Security tools**: Wireshark-equivalent for HTTP traffic
- **Debugging proxies**: mitmproxy tees traffic to show you the payload
- **WAFs**: inspect the request body without consuming it (the body must still reach the backend)
- **Compliance**: PCI DSS requires recording transaction data for certain merchant levels

---

## Day 36 — Connection Pools

### Learning Objectives
- Build a generic TCP connection pool — the foundation of your load balancer
- Implement idle timeout and max pool size
- Understand why connection pools dramatically reduce connection overhead

### Key Concepts

```
Without pool: each request = DNS + TCP handshake + (TLS handshake) + request
              latency: 50ms DNS + 1ms TCP + 5ms TLS = 56ms overhead per request

With pool:    first request connects, subsequent requests reuse the connection
              latency: ~0ms overhead for subsequent requests
              trade-off: idle connections consume server-side resources

Pool design:
  Get(addr):  return idle connection or dial new one
  Put(conn):  return connection to pool for reuse
  Close():    close all idle connections
  Idle timer: close connections idle > threshold
```

### Exercise: TCP Connection Pool

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

type Pool struct {
    mu          sync.Mutex
    idle        map[string][]idleConn
    maxPerAddr  int
    idleTimeout time.Duration
    dialer      net.Dialer
}

type idleConn struct {
    conn    net.Conn
    idleSince time.Time
}

func NewPool(maxPerAddr int, idleTimeout time.Duration) *Pool {
    p := &Pool{
        idle:        make(map[string][]idleConn),
        maxPerAddr:  maxPerAddr,
        idleTimeout: idleTimeout,
        dialer:      net.Dialer{Timeout: 5 * time.Second},
    }
    go p.reaper()
    return p
}

// Get returns an idle connection or dials a new one
func (p *Pool) Get(addr string) (net.Conn, error) {
    p.mu.Lock()
    idles := p.idle[addr]
    for len(idles) > 0 {
        ic := idles[len(idles)-1]
        idles = idles[:len(idles)-1]
        p.idle[addr] = idles
        p.mu.Unlock()

        // Test if connection is still alive (1-byte read with immediate timeout)
        ic.conn.SetReadDeadline(time.Now())
        buf := make([]byte, 1)
        _, err := ic.conn.Read(buf)
        ic.conn.SetReadDeadline(time.Time{}) // clear deadline

        if err != nil && !isTimeout(err) {
            ic.conn.Close() // stale connection
            p.mu.Lock()
            idles = p.idle[addr]
            continue
        }

        return ic.conn, nil // healthy idle connection
    }
    p.mu.Unlock()

    // No idle connection — dial new
    return p.dialer.Dial("tcp", addr)
}

// Put returns a connection to the pool for reuse
func (p *Pool) Put(addr string, conn net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()

    idles := p.idle[addr]
    if len(idles) >= p.maxPerAddr {
        conn.Close() // pool full — close excess
        return
    }
    p.idle[addr] = append(idles, idleConn{conn: conn, idleSince: time.Now()})
}

// reaper closes connections that have been idle too long
func (p *Pool) reaper() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        p.mu.Lock()
        cutoff := time.Now().Add(-p.idleTimeout)
        for addr, idles := range p.idle {
            var live []idleConn
            for _, ic := range idles {
                if ic.idleSince.After(cutoff) {
                    live = append(live, ic)
                } else {
                    ic.conn.Close()
                }
            }
            p.idle[addr] = live
        }
        p.mu.Unlock()
    }
}

func (p *Pool) Close() {
    p.mu.Lock()
    defer p.mu.Unlock()
    for _, idles := range p.idle {
        for _, ic := range idles {
            ic.conn.Close()
        }
    }
    p.idle = make(map[string][]idleConn)
}

func isTimeout(err error) bool {
    ne, ok := err.(net.Error)
    return ok && ne.Timeout()
}

func main() {
    pool := NewPool(10, 30*time.Second)
    defer pool.Close()

    // Simulate 50 concurrent requests to the same backend
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // In a real proxy, replace with your actual backend addr
            conn, err := pool.Get("example.com:80")
            if err != nil {
                fmt.Printf("[%d] dial error: %v\n", id, err)
                return
            }
            // Use connection...
            pool.Put("example.com:80", conn) // return to pool
        }(i)
    }
    wg.Wait()
    fmt.Println("Done")
}
```

### Real-World Context
`http.Transport` in Go's standard library is a connection pool. When you configure `MaxIdleConnsPerHost: 100`, that's the pool size. Your reverse proxy (Week 9) uses `httputil.ReverseProxy` which uses `http.Transport` under the hood. Building a pool from scratch shows you why those config options exist.

---

## Day 37 — DNS in Go

### Learning Objectives
- Resolve hostnames with TTL-aware caching
- Handle multi-IP responses and failover
- Measure DNS latency separately — it's often the hidden bottleneck

### Exercise: DNS-Aware Connection Dialer

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

type DNSCache struct {
    mu      sync.RWMutex
    entries map[string]cacheEntry
}

type cacheEntry struct {
    addrs  []string
    expiry time.Time
}

var dnsCache = &DNSCache{entries: make(map[string]cacheEntry)}

func (c *DNSCache) Lookup(host string) ([]string, error) {
    c.mu.RLock()
    if entry, ok := c.entries[host]; ok && time.Now().Before(entry.expiry) {
        c.mu.RUnlock()
        return entry.addrs, nil
    }
    c.mu.RUnlock()

    // Cache miss — resolve
    start := time.Now()
    addrs, err := net.LookupHost(host)
    dnsLatency := time.Since(start)
    fmt.Printf("[DNS] resolved %s → %v in %v\n", host, addrs, dnsLatency)

    if err != nil {
        return nil, err
    }

    c.mu.Lock()
    c.entries[host] = cacheEntry{
        addrs:  addrs,
        expiry: time.Now().Add(30 * time.Second), // 30s TTL (simplified)
    }
    c.mu.Unlock()

    return addrs, nil
}

// DialWithFallback tries each resolved IP until one connects
func DialWithFallback(host, port string) (net.Conn, string, error) {
    addrs, err := dnsCache.Lookup(host)
    if err != nil {
        return nil, "", fmt.Errorf("DNS lookup %s: %w", host, err)
    }

    dialer := net.Dialer{Timeout: 2 * time.Second}
    var lastErr error
    for _, addr := range addrs {
        target := net.JoinHostPort(addr, port)
        conn, err := dialer.Dial("tcp", target)
        if err == nil {
            return conn, addr, nil
        }
        lastErr = err
        fmt.Printf("[DIAL] %s failed: %v, trying next\n", target, err)
    }
    return nil, "", fmt.Errorf("all addresses for %s failed: %w", host, lastErr)
}

func main() {
    // First call: DNS miss, resolves
    conn, ip, err := DialWithFallback("example.com", "80")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    fmt.Printf("Connected to %s (via %s)\n", "example.com", ip)
    conn.Close()

    // Second call: DNS cache hit
    conn2, ip2, err := DialWithFallback("example.com", "80")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }
    fmt.Printf("Reconnected to %s (via %s, from cache)\n", "example.com", ip2)
    conn2.Close()
}
```

### Real-World Context
Multi-IP DNS responses are how CDNs do geographic load balancing. `dig example.com` returns different IPs depending on your location. Your load balancer must handle DNS changes — caching with TTL means it picks up new IPs when backends are added, and falls back to other IPs when one fails.

---

## Week 5 Mini Project: TCP Proxy with Traffic Inspector

### Architecture

```
Client → [TCP Proxy :8080] → [Upstream :9090]
                 ↓
          [Traffic Tap: tap.log]
          [Metrics: /debug/vars]
```

### Full Implementation

```go
// proxy/main.go
package main

import (
    "context"
    "expvar"
    "fmt"
    "io"
    "net"
    "net/http"
    _ "net/http/pprof"
    "os"
    "os/signal"
    "sync"
    "sync/atomic"
    "syscall"
    "time"
)

var (
    connAccepted  = expvar.NewInt("connections_accepted")
    connActive    = expvar.NewInt("connections_active")
    connRejected  = expvar.NewInt("connections_rejected")
    bytesInTotal  = expvar.NewInt("bytes_in_total")
    bytesOutTotal = expvar.NewInt("bytes_out_total")
)

const (
    maxConnections = 500
    idleTimeout    = 30 * time.Second
    upstreamAddr   = "127.0.0.1:9090"
    listenAddr     = ":8080"
)

type Proxy struct {
    listener net.Listener
    pool     *Pool
    tap      io.Writer
    active   atomic.Int64
    wg       sync.WaitGroup
    shutdown atomic.Bool
}

func NewProxy(listener net.Listener, upstreamAddr string, tap io.Writer) *Proxy {
    return &Proxy{
        listener: listener,
        pool:     NewPool(20, 30*time.Second),
        tap:      tap,
    }
}

func (p *Proxy) Serve(ctx context.Context) {
    for {
        conn, err := p.listener.Accept()
        if err != nil {
            select {
            case <-ctx.Done():
                return
            default:
                continue
            }
        }

        if p.active.Load() >= maxConnections {
            connRejected.Add(1)
            conn.Close()
            continue
        }

        connAccepted.Add(1)
        connActive.Add(1)
        p.active.Add(1)
        p.wg.Add(1)
        go p.handleConn(conn)
    }
}

func (p *Proxy) handleConn(client net.Conn) {
    defer client.Close()
    defer connActive.Add(-1)
    defer p.active.Add(-1)
    defer p.wg.Done()

    // Get upstream connection from pool
    upstream, err := p.pool.Get(upstreamAddr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "upstream connect: %v\n", err)
        return
    }

    remote := client.RemoteAddr().String()
    var bytesIn, bytesOut int64
    var innerWg sync.WaitGroup
    innerWg.Add(2)

    // client → upstream
    go func() {
        defer innerWg.Done()
        tapped := io.TeeReader(client,
            &hexLogger{w: p.tap, dir: fmt.Sprintf("IN  [%s]", remote)})
        n, _ := io.Copy(upstream, tapped)
        atomic.AddInt64(&bytesIn, n)
        bytesInTotal.Add(n)
        if tc, ok := upstream.(*net.TCPConn); ok {
            tc.CloseWrite()
        }
    }()

    // upstream → client
    go func() {
        defer innerWg.Done()
        tapped := io.TeeReader(upstream,
            &hexLogger{w: p.tap, dir: fmt.Sprintf("OUT [%s]", remote)})
        n, _ := io.Copy(client, tapped)
        atomic.AddInt64(&bytesOut, n)
        bytesOutTotal.Add(n)
        if tc, ok := client.(*net.TCPConn); ok {
            tc.CloseWrite()
        }
    }()

    innerWg.Wait()
    p.pool.Put(upstreamAddr, upstream)
}

func (p *Proxy) Drain() {
    p.wg.Wait()
}

type hexLogger struct {
    w   io.Writer
    dir string
}

func (h *hexLogger) Write(p []byte) (int, error) {
    preview := p
    if len(preview) > 64 {
        preview = preview[:64]
    }
    fmt.Fprintf(h.w, "%s %d bytes: %q\n", h.dir, len(p), preview)
    return len(p), nil
}

func main() {
    // Traffic tap
    tap, _ := os.Create("tap.log")
    defer tap.Close()

    ln, err := net.Listen("tcp", listenAddr)
    if err != nil {
        fmt.Println("listen:", err)
        return
    }

    proxy := NewProxy(ln, upstreamAddr, tap)

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()

    // Observability server
    go http.ListenAndServe(":9091", nil) // /debug/vars, /debug/pprof

    fmt.Printf("TCP proxy %s → %s\n", listenAddr, upstreamAddr)
    fmt.Println("Metrics: http://localhost:9091/debug/vars")

    go proxy.Serve(ctx)
    <-ctx.Done()

    fmt.Println("shutting down...")
    ln.Close()
    proxy.Drain()
    fmt.Println("shutdown complete")
}
```

### Engineer Takeaway
A TCP proxy is the foundation of load balancers (Week 7), WAFs (Week 10), and security gateways (Week 12). Cloudflare Spectrum, AWS NLB, and HAProxy in TCP mode all start with this: accept → bidirectional io.Copy → return to pool. The traffic tap is how compliance teams capture data for forensics without changing the application.
