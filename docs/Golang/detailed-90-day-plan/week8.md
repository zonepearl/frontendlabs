# Week 8: Phase 2 Consolidation (Days 52–60)

**Phase:** 2 — Networking Foundations + L4
**Goal:** Slower pace, deeper understanding. Solidify TCP internals, add SNI routing, audit for races, benchmark the full L4 stack. End with the Phase 2 capstone: Secure TCP Gateway.

---

> **Background in this wiki:** TCP/IP guide [Ch 23 (CLOSE_WAIT, reproduced)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-23-tcp-part-3-closing-a-connection-and-the-states), [Ch 22 (framing a byte stream)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-22-tcp-part-2-how-it-never-loses-your-data), [Ch 56 (certificate rotation)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation), and [Ch 32 (reading a ClientHello: the basis of SNI routing)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-32-wireshark-reading-a-conversation); HTTPS guide [Ch 14 (HTTP/2 multiplexing and negotiation)](../../v2-https/real-life-guide-v1.md#chapter-14-http-version-negotiation-http-1-1-http-2-http-3-alt-svc-and-fallbacks).

## Day 52 — TCP Internals: What Goes Wrong

### Learning Objectives
- Diagnose `CLOSE_WAIT` and `TIME_WAIT` states in production
- Understand `FIN_WAIT_2` and when it indicates a bug
- Know the difference between normal and pathological socket states

### TCP State Machine Reference

```
Normal close sequence (server initiates):
  Server → FIN → Client   [server: FIN_WAIT_1 → FIN_WAIT_2]
  Client → ACK → Server   [client: CLOSE_WAIT]
  Client → FIN → Server   [server: TIME_WAIT → CLOSED]
  Server → ACK → Client   [client: LAST_ACK → CLOSED]

TIME_WAIT (normal):
  Lasts 2*MSL (60-120 seconds)
  Prevents duplicate packets from old connections confusing new ones
  SOLUTION: increase ephemeral port range (not your bug)

CLOSE_WAIT (application bug):
  Remote side sent FIN (closed their end)
  Your code received it but never called conn.Close()
  Connection stuck in CLOSE_WAIT forever
  SOLUTION: always defer conn.Close()

FIN_WAIT_2 (possible leak):
  You sent FIN, waiting for remote's FIN
  Remote is slow or never sends FIN
  Linux has FIN_WAIT_2 timeout (default 60s)
  SOLUTION: set ReadDeadline to force close
```

### Exercise: Reproduce and Fix CLOSE_WAIT

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "time"
)

// BUG: This server never closes the connection, causing CLOSE_WAIT
func buggyServer(ln net.Listener) {
    for {
        conn, err := ln.Accept()
        if err != nil { return }
        go func() {
            // MISSING: defer conn.Close()
            reader := bufio.NewReader(conn)
            line, _ := reader.ReadString('\n')
            conn.Write([]byte("ECHO: " + line))
            // Goroutine exits but conn is NOT closed
            // When client closes, conn enters CLOSE_WAIT
        }()
    }
}

// FIXED: Always defer conn.Close()
func fixedServer(ln net.Listener) {
    for {
        conn, err := ln.Accept()
        if err != nil { return }
        go func() {
            defer conn.Close() // FIXED: always close
            conn.SetDeadline(time.Now().Add(30 * time.Second))
            reader := bufio.NewReader(conn)
            line, _ := reader.ReadString('\n')
            conn.Write([]byte("ECHO: " + line))
        }()
    }
}

func main() {
    ln, _ := net.Listen("tcp", ":9090")
    fmt.Println("Buggy server on :9090")
    fmt.Println("After connecting a client and closing it:")
    fmt.Println("  netstat -an | grep 9090  # look for CLOSE_WAIT")

    // Run buggy server for 10s, then fixed
    go buggyServer(ln)
    time.Sleep(10 * time.Second)
    ln.Close()

    fmt.Println("\nNow running fixed server on :9090")
    ln2, _ := net.Listen("tcp", ":9090")
    go fixedServer(ln2)
    time.Sleep(10 * time.Second)
    ln2.Close()
}
```

**Observe socket states:**
```bash
# Run server, connect a client, then close client
go run . &
nc 127.0.0.1 9090
# type a line, press Enter, then Ctrl+C (close nc)
# Observe CLOSE_WAIT:
netstat -an | grep 9090 | grep CLOSE_WAIT
# After fix, this should not appear
```

### Real-World Context
In a load balancer or proxy, `CLOSE_WAIT` exhaustion is a common incident. You get reports: "the service stops accepting connections" — and `netstat` shows thousands of `CLOSE_WAIT` sockets. The fix is always `defer conn.Close()` plus a `SetDeadline`. This is an OWASP A05 issue (Security Misconfiguration — resource exhaustion).

---

## Day 53 — Binary Protocol Design

### Learning Objectives
- Design a length-prefixed binary protocol for internal services
- Use `encoding/binary` for byte-order-safe encoding
- Benchmark binary vs JSON for internal communication

### Key Concepts

```
Why binary protocols?
  JSON:   human-readable, verbose, allocates strings
  Binary: compact, fast, no text parsing overhead
  gRPC (protobuf) is 5-10x more efficient than JSON for same data

Length-prefix framing:
  [4-byte length][payload bytes]
  Receiver: read 4 bytes → parse length → read exactly that many bytes
  Alternative: delimiter-based (e.g., newlines) — fragile with binary data

Protocol versioning:
  Always include a version field — allows backward compatibility
  [1-byte version][1-byte opcode][4-byte length][payload]
```

### Exercise: Binary KV Protocol

```go
package main

import (
    "encoding/binary"
    "fmt"
    "io"
    "net"
)

// Protocol: [1B version][1B opcode][4B key-len][key][4B val-len][val]
const protocolVersion = 1

type Opcode byte
const (
    OpSet    Opcode = 0x01
    OpGet    Opcode = 0x02
    OpDelete Opcode = 0x03
    OpOK     Opcode = 0x10
    OpNotFound Opcode = 0x11
    OpError  Opcode = 0xFF
)

type Message struct {
    Version byte
    Opcode  Opcode
    Key     []byte
    Value   []byte
}

func WriteMessage(w io.Writer, msg Message) error {
    // Header: version + opcode
    if _, err := w.Write([]byte{msg.Version, byte(msg.Opcode)}); err != nil {
        return err
    }

    // Key: 4-byte length + data
    klen := make([]byte, 4)
    binary.BigEndian.PutUint32(klen, uint32(len(msg.Key)))
    if _, err := w.Write(klen); err != nil { return err }
    if _, err := w.Write(msg.Key); err != nil { return err }

    // Value: 4-byte length + data
    vlen := make([]byte, 4)
    binary.BigEndian.PutUint32(vlen, uint32(len(msg.Value)))
    if _, err := w.Write(vlen); err != nil { return err }
    if len(msg.Value) > 0 {
        if _, err := w.Write(msg.Value); err != nil { return err }
    }
    return nil
}

func ReadMessage(r io.Reader) (Message, error) {
    header := make([]byte, 2)
    if _, err := io.ReadFull(r, header); err != nil {
        return Message{}, fmt.Errorf("read header: %w", err)
    }

    msg := Message{Version: header[0], Opcode: Opcode(header[1])}

    if msg.Version != protocolVersion {
        return Message{}, fmt.Errorf("unsupported version %d", msg.Version)
    }

    readField := func() ([]byte, error) {
        lenBuf := make([]byte, 4)
        if _, err := io.ReadFull(r, lenBuf); err != nil {
            return nil, err
        }
        length := binary.BigEndian.Uint32(lenBuf)
        if length == 0 { return nil, nil }
        if length > 10*1024*1024 { // 10MB max field size
            return nil, fmt.Errorf("field too large: %d bytes", length)
        }
        data := make([]byte, length)
        _, err := io.ReadFull(r, data)
        return data, err
    }

    var err error
    msg.Key, err = readField()
    if err != nil { return Message{}, fmt.Errorf("read key: %w", err) }
    msg.Value, err = readField()
    if err != nil { return Message{}, fmt.Errorf("read value: %w", err) }

    return msg, nil
}

func main() {
    server, client := net.Pipe()

    // Server: in-memory KV store
    go func() {
        defer server.Close()
        store := make(map[string][]byte)
        for {
            msg, err := ReadMessage(server)
            if err != nil { return }
            var resp Message
            resp.Version = protocolVersion
            switch msg.Opcode {
            case OpSet:
                store[string(msg.Key)] = msg.Value
                resp.Opcode = OpOK
            case OpGet:
                val, ok := store[string(msg.Key)]
                if ok {
                    resp.Opcode = OpOK
                    resp.Value = val
                } else {
                    resp.Opcode = OpNotFound
                }
            }
            WriteMessage(server, resp)
        }
    }()

    // Client: SET then GET
    WriteMessage(client, Message{Version: protocolVersion, Opcode: OpSet, Key: []byte("user:1"), Value: []byte(`{"name":"Alice","role":"admin"}`)})
    resp, _ := ReadMessage(client)
    fmt.Printf("SET response: opcode=0x%02X\n", resp.Opcode)

    WriteMessage(client, Message{Version: protocolVersion, Opcode: OpGet, Key: []byte("user:1")})
    resp2, _ := ReadMessage(client)
    fmt.Printf("GET response: value=%s\n", resp2.Value)

    client.Close()
}
```

### Real-World Context
Redis, Kafka, and MySQL all use binary protocols. Understanding framing prevents bugs like "reading a partial message" or "confusing the length with the data". The `io.ReadFull` call is critical — a plain `Read` might return less than requested on a busy network.

---

## Day 54 — TLS Certificate Rotation

### Learning Objectives
- Hot-reload TLS certificates without dropping connections
- Use `sync/atomic` for lock-free cert swapping
- Monitor cert expiry as part of server health

### Exercise: Zero-Downtime Cert Rotation

```go
package main

import (
    "crypto/tls"
    "fmt"
    "os"
    "sync/atomic"
    "time"
    "unsafe"
)

// CertReloader watches for cert file changes and atomically swaps certificates
type CertReloader struct {
    certPath    string
    keyPath     string
    certPointer unsafe.Pointer  // *tls.Certificate — atomic swap
    lastModTime time.Time
}

func NewCertReloader(certPath, keyPath string) (*CertReloader, error) {
    r := &CertReloader{certPath: certPath, keyPath: keyPath}
    if err := r.reload(); err != nil {
        return nil, err
    }
    return r, nil
}

func (r *CertReloader) reload() error {
    cert, err := tls.LoadX509KeyPair(r.certPath, r.keyPath)
    if err != nil {
        return fmt.Errorf("reload cert: %w", err)
    }

    // Parse to get expiry for logging
    x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
    if err == nil {
        daysLeft := int(time.Until(x509Cert.NotAfter).Hours() / 24)
        fmt.Printf("[CERT] loaded %s (expires in %d days)\n", x509Cert.Subject.CommonName, daysLeft)
    }

    // Atomic swap — zero downtime (no lock needed; pointer write is atomic on 64-bit)
    atomic.StorePointer(&r.certPointer, unsafe.Pointer(&cert))
    return nil
}

// GetCertificate is called by tls.Config for each TLS handshake
func (r *CertReloader) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
    ptr := atomic.LoadPointer(&r.certPointer)
    cert := (*tls.Certificate)(ptr)
    return cert, nil
}

// Watch polls for cert file changes and hot-reloads
func (r *CertReloader) Watch(stop <-chan struct{}) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-stop:
            return
        case <-ticker.C:
            info, err := os.Stat(r.certPath)
            if err != nil { continue }
            if info.ModTime().After(r.lastModTime) {
                r.lastModTime = info.ModTime()
                if err := r.reload(); err != nil {
                    fmt.Printf("[CERT] reload error: %v\n", err)
                } else {
                    fmt.Println("[CERT] reloaded successfully (no connections dropped)")
                }
            }
        }
    }
}

func main() {
    reloader, err := NewCertReloader("server.crt", "server.key")
    if err != nil {
        fmt.Println("ERROR:", err)
        return
    }

    stop := make(chan struct{})
    go reloader.Watch(stop)

    tlsConfig := &tls.Config{
        GetCertificate: reloader.GetCertificate, // called per-handshake
        MinVersion:     tls.VersionTLS12,
    }

    srv := &tls.Listener{} // attach tlsConfig
    _ = srv
    fmt.Println("TLS server with hot-reload on :8443")
    fmt.Println("Replace server.crt and server.key — reloads within 30s without restart")
    // In production: attach to a real server
    time.Sleep(24 * time.Hour)
    close(stop)
}
```

### Real-World Context
Let's Encrypt certificates expire every 90 days. Certbot and cert-manager automate renewal, but your server must reload without downtime. `GetCertificate` is the hook that Go's TLS implementation calls for each new handshake — new connections get the new cert while existing connections keep the old cert (they've already completed the handshake).

---

## Day 55 — SNI Routing

### Learning Objectives
- Route TLS connections by hostname before decrypting them (SNI routing)
- Serve multiple domains from one IP/port with different backends
- Understand how CDNs and multi-tenant proxies work

### Key Concepts

```
SNI (Server Name Indication):
  Extension to TLS ClientHello — client announces which hostname it's connecting to
  Server uses this to select the right certificate BEFORE the TLS handshake completes
  Without SNI: one IP can only serve one TLS certificate

How CDNs use SNI:
  1. All traffic hits the same IP (anycast)
  2. TLS ClientHello contains SNI (e.g., "acme.corp.com")
  3. CDN routes to the right customer's origin based on SNI
  4. CDN terminates TLS with the customer's cert
```

### Exercise: SNI Router

```go
package main

import (
    "crypto/tls"
    "fmt"
    "net"
    "strings"
)

type SNIRouter struct {
    // map from hostname → backend address
    routes   map[string]string
    reloader *CertReloader // from Day 54
}

func NewSNIRouter() *SNIRouter {
    return &SNIRouter{routes: make(map[string]string)}
}

func (r *SNIRouter) AddRoute(hostname, backendAddr string) {
    r.routes[hostname] = backendAddr
}

// GetCertificate picks the right cert based on SNI hostname
func (r *SNIRouter) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
    // In production: per-hostname cert store
    // Simplified: one cert that covers all hostnames via SANs
    return r.reloader.GetCertificate(hello)
}

// GetConfigForClient is called during TLS handshake — access to SNI before conn is established
func (r *SNIRouter) GetConfigForClient(hello *tls.ClientHelloInfo) (*tls.Config, error) {
    fmt.Printf("[SNI] hostname=%q from %s\n", hello.ServerName, hello.Conn.RemoteAddr())

    // Return per-hostname config (cert, min TLS version, etc.)
    return &tls.Config{
        GetCertificate: r.GetCertificate,
        MinVersion:     tls.VersionTLS12,
    }, nil
}

// Route returns the backend address for a given hostname
func (r *SNIRouter) Route(hostname string) (string, bool) {
    // Exact match
    if backend, ok := r.routes[hostname]; ok {
        return backend, true
    }
    // Wildcard match (*.example.com)
    parts := strings.SplitN(hostname, ".", 2)
    if len(parts) == 2 {
        wildcard := "*." + parts[1]
        if backend, ok := r.routes[wildcard]; ok {
            return backend, true
        }
    }
    return "", false
}

func handleSNIConn(conn net.Conn, router *SNIRouter) {
    defer conn.Close()

    // At this point, TLS handshake has completed
    // We know the SNI hostname from GetConfigForClient
    tlsConn, ok := conn.(*tls.Conn)
    if !ok { return }

    sni := tlsConn.ConnectionState().ServerName
    backend, found := router.Route(sni)
    if !found {
        fmt.Printf("[SNI] no route for %q\n", sni)
        return
    }

    // Connect to the appropriate backend and proxy
    upstream, err := net.Dial("tcp", backend)
    if err != nil {
        fmt.Printf("[SNI] backend %s unavailable: %v\n", backend, err)
        return
    }
    defer upstream.Close()

    fmt.Printf("[SNI] routing %q → %s\n", sni, backend)
    var wg sync.WaitGroup
    wg.Add(2)
    go func() { defer wg.Done(); io.Copy(upstream, conn) }()
    go func() { defer wg.Done(); io.Copy(conn, upstream) }()
    wg.Wait()
}

func main() {
    router := NewSNIRouter()
    router.AddRoute("api.example.com", "127.0.0.1:8081")
    router.AddRoute("admin.example.com", "127.0.0.1:8082")
    router.AddRoute("*.staging.example.com", "127.0.0.1:8083")

    // Reloader from Day 54
    reloader, _ := NewCertReloader("server.crt", "server.key")
    router.reloader = reloader

    tlsConfig := &tls.Config{
        GetConfigForClient: router.GetConfigForClient,
    }

    ln, err := tls.Listen("tcp", ":443", tlsConfig)
    if err != nil { fmt.Println(err); return }
    defer ln.Close()

    fmt.Println("SNI router on :443")
    for {
        conn, err := ln.Accept()
        if err != nil { return }
        go handleSNIConn(conn, router)
    }
}
```

**Test:**
```bash
# Connect with specific hostname (SNI)
curl --cacert ca.crt --resolve api.example.com:443:127.0.0.1 https://api.example.com/
curl --cacert ca.crt --resolve admin.example.com:443:127.0.0.1 https://admin.example.com/
```

### Real-World Context
Cloudflare, Fastly, and every multi-tenant CDN uses SNI routing. When your company hosts multiple products on the same IP, SNI routing lets each one have its own TLS cert and backend. This is also how `nginx ssl_server_name` and AWS ALB host-based routing work.

---

## Day 56 — Connection Multiplexing Concepts

### Learning Objectives
- Understand HTTP/2 multiplexing vs HTTP/1.1 connection-per-request
- Build a minimal stream multiplexer over a single TCP connection
- Know why H2 needs fewer backend connections than H1

### Exercise: Simple Stream Multiplexer

```go
package main

import (
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "sync"
)

// Simple multiplexer: multiple logical streams over one TCP connection
// Frame: [4B stream-id][2B frame-type][4B payload-len][payload]
const (
    FrameData  uint16 = 0x01
    FrameClose uint16 = 0x02
)

type Frame struct {
    StreamID uint32
    Type     uint16
    Payload  []byte
}

func WriteFrame(w io.Writer, f Frame) error {
    hdr := make([]byte, 10)
    binary.BigEndian.PutUint32(hdr[0:], f.StreamID)
    binary.BigEndian.PutUint16(hdr[4:], f.Type)
    binary.BigEndian.PutUint32(hdr[6:], uint32(len(f.Payload)))
    if _, err := w.Write(hdr); err != nil { return err }
    _, err := w.Write(f.Payload)
    return err
}

func ReadFrame(r io.Reader) (Frame, error) {
    hdr := make([]byte, 10)
    if _, err := io.ReadFull(r, hdr); err != nil {
        return Frame{}, err
    }
    f := Frame{
        StreamID: binary.BigEndian.Uint32(hdr[0:]),
        Type:     binary.BigEndian.Uint16(hdr[4:]),
    }
    length := binary.BigEndian.Uint32(hdr[6:])
    if length > 0 {
        f.Payload = make([]byte, length)
        if _, err := io.ReadFull(r, f.Payload); err != nil {
            return Frame{}, err
        }
    }
    return f, nil
}

type Mux struct {
    conn    net.Conn
    streams map[uint32]chan []byte
    mu      sync.Mutex
    wmu     sync.Mutex
}

func NewMux(conn net.Conn) *Mux {
    m := &Mux{conn: conn, streams: make(map[uint32]chan []byte)}
    go m.demux()
    return m
}

func (m *Mux) demux() {
    for {
        frame, err := ReadFrame(m.conn)
        if err != nil { return }
        m.mu.Lock()
        ch, ok := m.streams[frame.StreamID]
        m.mu.Unlock()
        if ok {
            ch <- frame.Payload
        }
    }
}

func (m *Mux) OpenStream(id uint32) chan []byte {
    ch := make(chan []byte, 10)
    m.mu.Lock()
    m.streams[id] = ch
    m.mu.Unlock()
    return ch
}

func (m *Mux) Send(streamID uint32, data []byte) error {
    m.wmu.Lock()
    defer m.wmu.Unlock()
    return WriteFrame(m.conn, Frame{StreamID: streamID, Type: FrameData, Payload: data})
}

func main() {
    server, client := net.Pipe()

    // Server: echo all streams
    serverMux := NewMux(server)
    go func() {
        for id := uint32(1); id <= 3; id++ {
            ch := serverMux.OpenStream(id)
            go func(streamID uint32, recv chan []byte) {
                for data := range recv {
                    serverMux.Send(streamID, []byte("ECHO:"+string(data)))
                }
            }(id, ch)
        }
    }()

    // Client: 3 concurrent streams over 1 connection
    clientMux := NewMux(client)
    var wg sync.WaitGroup
    for id := uint32(1); id <= 3; id++ {
        wg.Add(1)
        ch := clientMux.OpenStream(id)
        go func(streamID uint32, recv chan []byte) {
            defer wg.Done()
            clientMux.Send(streamID, []byte(fmt.Sprintf("hello from stream %d", streamID)))
            response := <-recv
            fmt.Printf("[stream %d] %s\n", streamID, response)
        }(id, ch)
    }
    wg.Wait()
    server.Close()
    client.Close()
}
```

### Real-World Context
HTTP/2 multiplexes 100+ requests over 1 TCP connection. Your L4 load balancer must NOT use round-robin per-connection for H2 — all 100 requests would go to the same backend. Your L7 proxy (Week 9) uses `httputil.ReverseProxy` which correctly handles H2 multiplexing.

---

## Day 57 — Error Propagation in Network Code

### Learning Objectives
- Classify network errors correctly for retry logic
- Know the difference between `io.EOF` (clean close) and `io.ErrUnexpectedEOF` (crash)
- Build an error classifier used by your retry and circuit breaker logic

### Exercise: Network Error Classifier

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "net"
    "os"
    "syscall"
)

type NetworkError struct {
    Kind    string
    Retry   bool   // should the caller retry?
    Severity string // "info", "warn", "error"
    Err     error
}

func ClassifyNetworkError(err error) NetworkError {
    if err == nil {
        return NetworkError{}
    }

    // Clean close — peer closed connection gracefully
    if err == io.EOF {
        return NetworkError{Kind: "EOF", Retry: false, Severity: "info", Err: err}
    }

    // Partial read — peer crashed mid-message
    if err == io.ErrUnexpectedEOF {
        return NetworkError{Kind: "UNEXPECTED_EOF", Retry: true, Severity: "warn", Err: err}
    }

    // Deadline exceeded (our timeout)
    if os.IsTimeout(err) {
        return NetworkError{Kind: "TIMEOUT", Retry: true, Severity: "warn", Err: err}
    }

    // Check net.Error interface
    var netErr net.Error
    if errors.As(err, &netErr) {
        if netErr.Timeout() {
            return NetworkError{Kind: "NET_TIMEOUT", Retry: true, Severity: "warn", Err: err}
        }
    }

    // System call errors
    var sysErr syscall.Errno
    if errors.As(err, &sysErr) {
        switch sysErr {
        case syscall.ECONNREFUSED:
            return NetworkError{Kind: "CONN_REFUSED", Retry: true, Severity: "error", Err: err}
        case syscall.ECONNRESET:
            return NetworkError{Kind: "CONN_RESET", Retry: false, Severity: "warn", Err: err}
        case syscall.EPIPE:
            return NetworkError{Kind: "BROKEN_PIPE", Retry: false, Severity: "warn", Err: err}
        case syscall.ETIMEDOUT:
            return NetworkError{Kind: "SYS_TIMEOUT", Retry: true, Severity: "warn", Err: err}
        }
    }

    return NetworkError{Kind: "UNKNOWN", Retry: false, Severity: "error", Err: err}
}

func main() {
    // Simulate different error types
    testErrors := []error{
        io.EOF,
        io.ErrUnexpectedEOF,
        &net.OpError{Op: "read", Net: "tcp", Err: syscall.ECONNRESET},
        &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED},
    }

    for _, err := range testErrors {
        classified := ClassifyNetworkError(err)
        fmt.Printf("Error: %-20v Kind: %-20s Retry: %-5v Severity: %s\n",
            err, classified.Kind, classified.Retry, classified.Severity)
    }
}
```

### Real-World Context
In your circuit breaker (Day 51), `ECONNRESET` should count as a failure. `io.EOF` (client disconnected) should NOT count as a backend failure — it's the client's issue. Misclassifying errors causes false circuit breaker trips and unnecessary retries.

---

## Day 58 — Race Conditions in Network Servers

### Learning Objectives
- Find races in the Week 5 TCP proxy using `-race`
- Understand the three most common race patterns in network code
- Fix races without over-locking

### Common Race Patterns

```
Race 1: Shared connection state
  Two goroutines (read + write) both check conn.closed
  Fix: sync/atomic or sync.Once for close

Race 2: Concurrent map writes
  Handler goroutine reads map while health checker writes
  Fix: sync.RWMutex around the map

Race 3: Deadline race
  Read goroutine sets read deadline
  Write goroutine also sets deadline (via SetDeadline)
  They race to set the deadline
  Fix: separate SetReadDeadline / SetWriteDeadline

Race 4: Goroutine closure variable capture
  for i := range items { go func() { fmt.Println(i) }() }
  All goroutines see the last value of i
  Fix: pass i as argument: go func(id int) { fmt.Println(id) }(i)
```

### Exercise: Race Audit Checklist

```bash
# Run with race detector
go test -race -v ./...
go build -race . && ./proxy  # run the binary with race detection

# Common findings and fixes:
# DATA RACE: concurrent map read and write
# Fix: add sync.RWMutex

# DATA RACE: conn.closed written from two goroutines
# Fix: use atomic.Bool

# DATA RACE: per-connection stats map
# Fix: use sync.Map or per-connection struct with single owner goroutine
```

**Annotated fixes for the Week 5 proxy:**

```go
// Before (race): both goroutines write to same conn field
type connection struct {
    conn     net.Conn
    closed   bool        // RACE: read and written from two goroutines
    bytesIn  int64       // RACE: written from read goroutine, read from stats goroutine
    bytesOut int64
}

// After (fixed):
type connection struct {
    conn     net.Conn
    closed   atomic.Bool // atomic — safe from multiple goroutines
    bytesIn  atomic.Int64  // atomic counter — no lock needed
    bytesOut atomic.Int64
}
```

### Real-World Context
The race detector finds bugs that appear 1 in 10,000 times in production — intermittent corruptions and occasional panics. Running `-race` in CI is mandatory for any networked service. A data race in a proxy can silently corrupt responses — security-relevant because it could mix up authenticated and unauthenticated responses.

---

## Day 59 — Benchmark the Full L4 Stack

### Learning Objectives
- Write a load test in pure Go — measure RPS, latency percentiles
- Know the numbers: what throughput can your proxy handle?
- Use the results to set SLOs and identify bottlenecks

### Exercise: Load Test the L4 Load Balancer

```go
package main

import (
    "fmt"
    "net"
    "sort"
    "sync"
    "sync/atomic"
    "time"
)

type LoadTestResult struct {
    TotalRequests int64
    Errors        int64
    Duration      time.Duration
    Latencies     []time.Duration
}

func (r *LoadTestResult) Print() {
    rps := float64(r.TotalRequests) / r.Duration.Seconds()
    errorRate := float64(r.Errors) / float64(r.TotalRequests) * 100

    sort.Slice(r.Latencies, func(i, j int) bool { return r.Latencies[i] < r.Latencies[j] })
    pct := func(p float64) time.Duration {
        if len(r.Latencies) == 0 { return 0 }
        idx := int(p / 100 * float64(len(r.Latencies)))
        if idx >= len(r.Latencies) { idx = len(r.Latencies) - 1 }
        return r.Latencies[idx]
    }

    fmt.Printf("\n=== Load Test Results ===\n")
    fmt.Printf("Duration:   %v\n", r.Duration.Round(time.Millisecond))
    fmt.Printf("Total:      %d requests\n", r.TotalRequests)
    fmt.Printf("RPS:        %.0f req/s\n", rps)
    fmt.Printf("Error rate: %.2f%% (%d errors)\n", errorRate, r.Errors)
    fmt.Printf("Latency p50: %v\n", pct(50))
    fmt.Printf("Latency p95: %v\n", pct(95))
    fmt.Printf("Latency p99: %v\n", pct(99))
    fmt.Printf("Latency max: %v\n", pct(100))
}

func runLoadTest(proxyAddr string, concurrency, requestsPerClient int) LoadTestResult {
    var (
        totalRequests atomic.Int64
        errors        atomic.Int64
        mu            sync.Mutex
        allLatencies  []time.Duration
    )

    start := time.Now()
    var wg sync.WaitGroup
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            var localLatencies []time.Duration
            for j := 0; j < requestsPerClient; j++ {
                t := time.Now()
                conn, err := net.DialTimeout("tcp", proxyAddr, 2*time.Second)
                if err != nil {
                    errors.Add(1)
                    continue
                }
                msg := []byte("ping\n")
                conn.Write(msg)
                buf := make([]byte, 16)
                conn.SetReadDeadline(time.Now().Add(2 * time.Second))
                conn.Read(buf)
                conn.Close()
                localLatencies = append(localLatencies, time.Since(t))
                totalRequests.Add(1)
            }
            mu.Lock()
            allLatencies = append(allLatencies, localLatencies...)
            mu.Unlock()
        }()
    }
    wg.Wait()

    return LoadTestResult{
        TotalRequests: totalRequests.Load(),
        Errors:        errors.Load(),
        Duration:      time.Since(start),
        Latencies:     allLatencies,
    }
}

func main() {
    fmt.Println("Load testing L4 proxy at :8080")
    fmt.Println("(Start the proxy from Week 5/7 first)")
    fmt.Println("Concurrency: 100, 100 requests each = 10,000 total")

    result := runLoadTest(":8080", 100, 100)
    result.Print()
}
```

**Expected results on a developer machine:**
```
Duration:     ~2-5 seconds
RPS:          2,000-5,000 req/s
Latency p50:  ~1ms
Latency p95:  ~5ms
Latency p99:  ~20ms
Error rate:   0%
```

### Real-World Context
"Our proxy adds 2ms p99 latency" is a concrete, defensible number. Without benchmarks, every performance discussion is guessing. Run this load test before and after adding WAF middleware in Week 10 — that gives you the exact overhead of security inspection.

---

## Day 60 — Phase 2 Review

### Review Checklist

```
Architecture Review:
[ ] Can you explain the mTLS flow from cert generation to identity extraction?
[ ] Can you explain why TIME_WAIT is normal and CLOSE_WAIT is a bug?
[ ] Can you describe the circuit breaker state machine?
[ ] Can you explain the difference between L4 and L7 load balancing?

Code Audit:
[ ] go test -race ./... passes on all Phase 2 projects
[ ] All connections have defer conn.Close()
[ ] All tickers have defer ticker.Stop()
[ ] Context propagated through all goroutine-spawning code
[ ] No unbounded goroutine spawning (max connections enforced)

Metrics:
[ ] Every component emits expvar metrics
[ ] /stats endpoint shows backend health
[ ] Latency percentiles computed from ring buffer

What Would Break at 10x Traffic?
[ ] Connection pool exhaustion (increase maxPerAddr)
[ ] DNS cache thrash (increase TTL)
[ ] Circuit breaker sensitivity (tune thresholds)
[ ] TLS handshake CPU (consider session tickets)
```

---

## Week 8 Mini Project: Secure TCP Gateway

Combine all Phase 2 work into a single, production-grade secure TCP gateway.

### Architecture

```
Internet Traffic
    │
    ▼
┌──────────────────────────────────────────────────┐
│              Secure TCP Gateway :443              │
│                                                  │
│  1. mTLS Termination (client cert required)      │
│  2. SNI Routing → select backend cluster         │
│  3. Per-IP Rate Limiting (100 conns/IP)          │
│  4. Circuit Breaker per backend                  │
│  5. Least-Connections → pool connection          │
│  6. Bidirectional proxy with metrics tap         │
│                                                  │
│  Admin API :9090                                 │
│    /stats   /debug/vars   /debug/pprof           │
│    SIGHUP → hot-reload cert                      │
│    SIGTERM → graceful drain                      │
└──────────────────────────────────────────────────┘
    │              │              │
    ▼              ▼              ▼
Backend A      Backend B      Backend C
```

### Feature Integration Map

| Feature | Implemented in | Integrated via |
|---------|---------------|----------------|
| mTLS termination | Day 41 | `tls.RequireAndVerifyClientCert` |
| SNI routing | Day 55 | `tls.Config.GetConfigForClient` |
| Cert hot-reload | Day 54 | `SIGHUP` → `CertReloader.reload()` |
| Connection pool | Day 36 | `Pool.Get()` / `Pool.Put()` |
| Health checker | Day 46 | `HealthChecker.Run()` in goroutine |
| Circuit breaker | Day 51 | `CircuitBreaker.Allow()` / `RecordFailure()` |
| Rate limiter | Day 49 | `IPLimiter.Allow()` / `Release()` |
| Observability | Day 50 | `expvar`, `/stats` handler |
| Graceful shutdown | Day 27 | `signal.NotifyContext` + drain |

### Engineer Takeaway
This is roughly what Cloudflare Spectrum does for TCP proxying. The entire implementation is ~500 lines of Go stdlib code. You've moved from "I read the architecture doc" to "I built it." L4 networking is now something you can debug and extend, not just configure.
