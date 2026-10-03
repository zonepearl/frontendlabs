# Week 9: HTTP Internals (Days 61–67)

**Phase:** 3 — L7 Networking + WAF + API Security
**Goal:** Move from raw TCP to the application layer. Build an HTTP/1.1 parser from scratch, understand why HTTP request smuggling works, configure a hardened reverse proxy, and implement WebSocket from the wire up.

---

> **Background in this wiki:** TCP/IP guide [Ch 28 (`rawhttp`)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-28-http-how-the-web-actually-talks), [Ch 39 (HTTP framing)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-39-how-http-messages-are-framed-and-how-it-goes-wrong), [Ch 43 (WebSockets and SSE)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-43-websockets-and-server-sent-events); HTTPS guide [Ch 19 (production server and Slowloris)](../../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown), [Ch 20 (reverse proxy, X-Forwarded-For, smuggling)](../../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling), [Ch 24 (streaming)](../../v2-https/real-life-guide-v1.md#chapter-24-streaming-over-http-server-sent-events-websockets-and-long-lived-connections).

## Day 61 — HTTP/1.1 From Scratch

### Learning Objectives
- Parse HTTP manually using `bufio.Reader` — understand the wire format
- Know why HTTP request smuggling attacks are possible
- Implement keep-alive properly

### Key Concepts

```
HTTP/1.1 Wire Format:
  Request:
    GET /api/data?id=1 HTTP/1.1\r\n    ← request line
    Host: example.com\r\n               ← headers (key: value)
    Content-Type: application/json\r\n
    Content-Length: 42\r\n
    \r\n                                ← empty line = end of headers
    {"key":"value","other":"data"}      ← body (Content-Length bytes)

  Response:
    HTTP/1.1 200 OK\r\n
    Content-Type: application/json\r\n
    Content-Length: 19\r\n
    \r\n
    {"status":"success"}

HTTP Request Smuggling (why it works):
  Some proxies use Content-Length, some use Transfer-Encoding: chunked
  When they disagree on where the body ends, a second request can be injected
  Example: proxy reads 100 bytes as body; backend reads 20 bytes as body + 80 as new request
  OWASP A10 (Server-Side Request Forgery / HTTP Desync)
  Defense: normalize headers, reject ambiguous requests
```

### Exercise: Raw HTTP/1.1 Server

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
)

type HTTPRequest struct {
    Method  string
    Path    string
    Proto   string
    Headers map[string]string
    Body    []byte
}

// parseRequest parses a raw HTTP/1.1 request from bufio.Reader
func parseRequest(r *bufio.Reader) (*HTTPRequest, error) {
    // Request line: "GET /path HTTP/1.1"
    line, err := r.ReadString('\n')
    if err != nil { return nil, err }
    line = strings.TrimRight(line, "\r\n")

    parts := strings.SplitN(line, " ", 3)
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid request line: %q", line)
    }

    req := &HTTPRequest{
        Method:  parts[0],
        Path:    parts[1],
        Proto:   parts[2],
        Headers: make(map[string]string),
    }

    // Validate method (allowlist — reject unexpected methods early)
    validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true, "HEAD": true, "OPTIONS": true}
    if !validMethods[req.Method] {
        return nil, fmt.Errorf("invalid method: %q", req.Method)
    }

    // Headers: "Key: Value\r\n" until empty line
    for {
        line, err := r.ReadString('\n')
        if err != nil { return nil, err }
        line = strings.TrimRight(line, "\r\n")
        if line == "" { break } // empty line = end of headers

        colonIdx := strings.IndexByte(line, ':')
        if colonIdx < 0 {
            return nil, fmt.Errorf("invalid header: %q", line)
        }
        key := strings.TrimSpace(strings.ToLower(line[:colonIdx]))
        val := strings.TrimSpace(line[colonIdx+1:])

        // Security: limit header count and size
        if len(req.Headers) > 100 {
            return nil, fmt.Errorf("too many headers")
        }
        if len(key)+len(val) > 8192 {
            return nil, fmt.Errorf("header too large")
        }
        req.Headers[key] = val
    }

    // Body: read exactly Content-Length bytes
    // SECURITY: reject both Content-Length AND Transfer-Encoding (smuggling vector)
    _, hasTE := req.Headers["transfer-encoding"]
    clStr, hasCL := req.Headers["content-length"]
    if hasTE && hasCL {
        return nil, fmt.Errorf("ambiguous body framing: both Content-Length and Transfer-Encoding")
    }

    if hasCL {
        cl, err := strconv.ParseInt(clStr, 10, 64)
        if err != nil || cl < 0 {
            return nil, fmt.Errorf("invalid Content-Length: %q", clStr)
        }
        if cl > 10*1024*1024 { // 10MB limit
            return nil, fmt.Errorf("body too large: %d bytes", cl)
        }
        req.Body = make([]byte, cl)
        if _, err := r.Read(req.Body); err != nil {
            return nil, fmt.Errorf("read body: %w", err)
        }
    }

    return req, nil
}

func writeResponse(conn net.Conn, status int, body string) {
    statusText := map[int]string{200: "OK", 400: "Bad Request", 405: "Method Not Allowed", 413: "Payload Too Large"}
    text := statusText[status]
    if text == "" { text = "Unknown" }

    fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s",
        status, text, len(body), body)
}

func handleHTTPConn(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)

    req, err := parseRequest(reader)
    if err != nil {
        writeResponse(conn, 400, "bad request: "+err.Error())
        return
    }

    response := fmt.Sprintf("Method: %s\nPath: %s\nHost: %s\n",
        req.Method, req.Path, req.Headers["host"])
    writeResponse(conn, 200, response)
}

func main() {
    ln, _ := net.Listen("tcp", ":8080")
    fmt.Println("Raw HTTP/1.1 server on :8080")
    fmt.Println("Test: curl http://localhost:8080/api/test")
    for {
        conn, err := ln.Accept()
        if err != nil { return }
        go handleHTTPConn(conn)
    }
}
```

### Real-World Context
HTTP request smuggling (OWASP A10) exploits disagreement between proxy and backend on where one request ends and the next begins. Your proxy must: (1) reject requests with both `Content-Length` and `Transfer-Encoding`, (2) normalize `Transfer-Encoding` before forwarding. `net/http` handles this for you, but knowing why helps you evaluate WAF rules.

---

## Day 62 — net/http Server Configuration

### Learning Objectives
- Configure all four `http.Server` timeouts — each prevents a different attack
- Use `ConnState` callback to track connection lifecycle
- Understand why default `http.Server` is insecure for production

### Key Concepts

```
http.Server timeout fields and what they prevent:

  ReadHeaderTimeout: 5s  — Slowloris: attacker sends headers 1 byte/sec
  ReadTimeout:      30s  — Slow body: attacker sends body 1 byte/sec
  WriteTimeout:     30s  — Slow read: attacker reads response 1 byte/sec
  IdleTimeout:     120s  — Connection leeching: attacker holds keep-alive open
  MaxHeaderBytes:  64KB  — Header bomb: attacker sends huge headers

Without these:
  Slowloris + 10,000 connections = server out of file descriptors
  Default http.Server has NO timeouts — use ONLY with a reverse proxy in front
```

### Exercise: Hardened HTTP Server

```go
package main

import (
    "fmt"
    "net/http"
    "sync/atomic"
    "time"
)

type ConnStats struct {
    New     atomic.Int64
    Active  atomic.Int64
    Idle    atomic.Int64
    Closed  atomic.Int64
}

func (cs *ConnStats) Handler(conn net.Conn, state http.ConnState) {
    switch state {
    case http.StateNew:
        cs.New.Add(1)
        cs.Active.Add(1)
    case http.StateActive:
        cs.Idle.Add(-1)
    case http.StateIdle:
        cs.Active.Add(-1)
        cs.Idle.Add(1)
    case http.StateClosed, http.StateHijacked:
        cs.Active.Add(-1)
        cs.Idle.Add(-1)
        cs.Closed.Add(1)
    }
}

func main() {
    stats := &ConnStats{}
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "OK\n")
    })
    mux.HandleFunc("/conn-stats", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "new=%d active=%d idle=%d closed=%d\n",
            stats.New.Load(), stats.Active.Load(), stats.Idle.Load(), stats.Closed.Load())
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,

        // SECURITY: all timeouts must be set in production
        ReadHeaderTimeout: 5 * time.Second,   // Slowloris defense
        ReadTimeout:       30 * time.Second,  // slow body defense
        WriteTimeout:      30 * time.Second,  // slow client defense
        IdleTimeout:       120 * time.Second, // keep-alive limit
        MaxHeaderBytes:    64 * 1024,         // 64KB header limit

        // Connection lifecycle monitoring
        ConnState: stats.Handler,
    }

    fmt.Println("Hardened HTTP server on :8080")
    fmt.Println("Test Slowloris: pip install slowloris && slowloris -p 8080 localhost")
    fmt.Println("  With ReadHeaderTimeout=5s: server stays alive")
    fmt.Println("  Without timeout: server dies in ~30s")

    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        fmt.Println("ERROR:", err)
    }
}
```

### Real-World Context
The default `http.ListenAndServe(":8080", handler)` creates a server with NO timeouts. Behind a properly configured nginx or ALB, this is fine. Directly on the internet: Slowloris will take it down in minutes. PCI DSS and SOC 2 scans specifically check for missing timeout configurations.

---

## Day 63 — HTTP Middleware

### Learning Objectives
- Build a middleware chain for the security proxy
- Capture status codes and response size without breaking the handler interface
- Implement panic recovery — prevents one request from crashing the whole server

### Exercise: Security Middleware Chain

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "runtime/debug"
    "time"
)

// statusCapture wraps ResponseWriter to capture status code written by handler
type statusCapture struct {
    http.ResponseWriter
    statusCode  int
    bytesWritten int64
    headersSent  bool
}

func (sc *statusCapture) WriteHeader(code int) {
    sc.statusCode = code
    sc.headersSent = true
    sc.ResponseWriter.WriteHeader(code)
}

func (sc *statusCapture) Write(b []byte) (int, error) {
    if !sc.headersSent {
        sc.WriteHeader(http.StatusOK)
    }
    n, err := sc.ResponseWriter.Write(b)
    sc.bytesWritten += int64(n)
    return n, err
}

// Chain executes middlewares in order (first middleware is outermost)
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- {
        h = mws[i](h)
    }
    return h
}

// PanicRecoverer catches panics in handlers and returns 500
func PanicRecoverer(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    logger.Error("handler panic",
                        "error", fmt.Sprintf("%v", rec),
                        "stack", string(debug.Stack()),
                        "path", r.URL.Path,
                    )
                    // Only write 500 if headers not yet sent
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusInternalServerError)
                    fmt.Fprintf(w, `{"error":"internal server error"}`)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}

// RequestLogger logs each request as structured JSON
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            capture := &statusCapture{ResponseWriter: w, statusCode: 200}
            next.ServeHTTP(capture, r)
            logger.Info("request",
                "method", r.Method,
                "path", r.URL.Path,
                "status", capture.statusCode,
                "bytes", capture.bytesWritten,
                "latency_ms", time.Since(start).Milliseconds(),
                "ip", r.RemoteAddr,
                "user_agent", r.UserAgent(),
            )
        })
    }
}

// SecurityHeaders sets security response headers on all responses
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        h := w.Header()
        h.Set("X-Content-Type-Options", "nosniff")
        h.Set("X-Frame-Options", "DENY")
        h.Set("X-XSS-Protection", "0")             // modern browsers: use CSP instead
        h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        h.Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
        // HSTS: only set on HTTPS responses
        if r.TLS != nil {
            h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
        }
        next.ServeHTTP(w, r)
    })
}

// RequestID injects a unique ID into the request context and response header
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := generateRequestID()
        w.Header().Set("X-Request-ID", id)
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func main() {
    logger := slog.Default()

    appHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/panic" {
            panic("deliberate panic for testing")
        }
        fmt.Fprintf(w, `{"status":"ok","path":%q}`, r.URL.Path)
    })

    // Middleware order (outermost first):
    // PanicRecoverer → RequestID → SecurityHeaders → Logger → handler
    handler := Chain(appHandler,
        PanicRecoverer(logger),
        RequestID,
        SecurityHeaders,
        RequestLogger(logger),
    )

    srv := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadHeaderTimeout: 5 * time.Second,
    }

    fmt.Println("Server with middleware chain on :8080")
    fmt.Println("Test panic recovery: curl http://localhost:8080/panic")
    srv.ListenAndServe()
}
```

### Real-World Context
Middleware ordering is a security decision: `PanicRecoverer` must be outermost (catches panics in all downstream middleware); `RequestID` before `Logger` (so logs include the request ID); `SecurityHeaders` unconditionally (always set, regardless of handler outcome). OWASP A05 — missing security headers is a common audit finding.

---

## Day 64 — HTTP Reverse Proxy

### Learning Objectives
- Configure `httputil.ReverseProxy` with all production options
- Handle hop-by-hop headers correctly
- Inject `X-Forwarded-For` and `X-Request-ID` securely

### Key Concepts

```
Hop-by-hop headers — must NOT be forwarded:
  Connection, Keep-Alive, Transfer-Encoding, TE, Upgrade,
  Proxy-Authorization, Proxy-Authenticate

X-Forwarded-For security:
  Proxy APPENDS the client IP (never overwrites)
  X-Forwarded-For: client-ip, proxy1-ip, proxy2-ip
  Backend should use the FIRST entry in the chain
  Never trust X-Forwarded-For from untrusted clients — it can be spoofed
  Only trust the XFF chain from IPs you control

X-Real-IP:
  Proxy sets to the direct client IP (before any other proxies)
  Never trust X-Real-IP if set by the client
```

### Exercise: Production Reverse Proxy

```go
package main

import (
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "time"
)

// hopByHop headers must not be forwarded upstream
var hopByHopHeaders = map[string]bool{
    "connection":          true,
    "keep-alive":         true,
    "transfer-encoding":  true,
    "te":                 true,
    "upgrade":            true,
    "proxy-authorization": true,
    "proxy-authenticate": true,
    "trailers":           true,
}

// internalHeaders must not be set by clients (privilege escalation)
var sensitiveInternalHeaders = []string{
    "X-Internal-Role",
    "X-Admin",
    "X-Authenticated-Subject",
    "X-User-ID",
}

func NewSecureReverseProxy(upstream *url.URL) *httputil.ReverseProxy {
    proxy := httputil.NewSingleHostReverseProxy(upstream)

    // Director: modify request before sending to upstream
    proxy.Director = func(r *http.Request) {
        // Standard reverse proxy setup
        r.URL.Scheme = upstream.Scheme
        r.URL.Host = upstream.Host
        r.Host = upstream.Host

        // Strip hop-by-hop headers
        for header := range hopByHopHeaders {
            r.Header.Del(header)
        }

        // Strip client-supplied internal headers (security: prevent privilege escalation)
        for _, h := range sensitiveInternalHeaders {
            r.Header.Del(h)
        }

        // X-Forwarded-For: append real client IP, don't trust existing
        clientIP := r.RemoteAddr
        if host, _, err := net.SplitHostPort(clientIP); err == nil {
            clientIP = host
        }
        if prior := r.Header.Get("X-Forwarded-For"); prior != "" {
            r.Header.Set("X-Forwarded-For", prior+", "+clientIP)
        } else {
            r.Header.Set("X-Forwarded-For", clientIP)
        }

        // X-Forwarded-Proto
        if r.TLS != nil {
            r.Header.Set("X-Forwarded-Proto", "https")
        } else {
            r.Header.Set("X-Forwarded-Proto", "http")
        }

        // Inject request ID from context
        if id := r.Context().Value(requestIDKey); id != nil {
            r.Header.Set("X-Request-ID", id.(string))
        }
    }

    // ModifyResponse: clean up upstream response before returning to client
    proxy.ModifyResponse = func(resp *http.Response) error {
        // Strip upstream's internal headers from response
        resp.Header.Del("X-Powered-By")
        resp.Header.Del("Server")
        resp.Header.Del("X-AspNet-Version")

        // Add correlation header
        if id := resp.Request.Header.Get("X-Request-ID"); id != "" {
            resp.Header.Set("X-Request-ID", id)
        }
        return nil
    }

    // ErrorHandler: what to return when upstream is unreachable
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        fmt.Printf("[PROXY] upstream error: %v\n", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadGateway)
        fmt.Fprintf(w, `{"error":"upstream unavailable","request_id":%q}`,
            r.Context().Value(requestIDKey))
    }

    // Custom transport with timeouts
    proxy.Transport = &http.Transport{
        DialContext:           (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
        TLSHandshakeTimeout:   5 * time.Second,
        ResponseHeaderTimeout: 30 * time.Second,
        MaxIdleConns:          100,
        MaxIdleConnsPerHost:   20,
        IdleConnTimeout:       90 * time.Second,
    }

    return proxy
}

func main() {
    upstream, _ := url.Parse("http://127.0.0.1:9090")
    proxy := NewSecureReverseProxy(upstream)

    handler := Chain(proxy,
        PanicRecoverer(slog.Default()),
        RequestID,
        SecurityHeaders,
        RequestLogger(slog.Default()),
    )

    srv := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      60 * time.Second, // higher: includes upstream response time
        IdleTimeout:       120 * time.Second,
    }

    fmt.Println("Reverse proxy :8080 → :9090")
    srv.ListenAndServe()
}
```

### Real-World Context
Stripping `X-Internal-Role` that clients might set prevents privilege escalation — OWASP A01 (Broken Access Control). Setting `X-Forwarded-For` correctly is required for IP-based rate limiting and WAF to see the real client IP. If your proxy overwrites (not appends) XFF, you lose the audit trail.

---

## Day 65 — HTTP/2

### Learning Objectives
- Enable HTTP/2 on your TLS server (it's automatic with TLS)
- Understand why H2 changes connection pool sizing assumptions
- Measure latency difference between H1 and H2

### Key Concepts

```
HTTP/1.1 vs HTTP/2:
  H1: 1 request per connection at a time (pipelining is broken in practice)
      browsers open 6 connections per origin to get parallelism
  H2: multiple requests multiplexed over 1 connection (streams)
      browser needs only 1 connection per origin
      header compression (HPACK): repeated headers cost ~10 bytes vs ~200 bytes

Why H2 changes backend pool sizing:
  H1: 100 concurrent requests → 100 backend connections
  H2: 100 concurrent requests → 1-10 backend connections (fewer, multiplexed)
  Your connection pool maxIdleConnsPerHost should be LOWER for H2 backends

H2 gotchas:
  - Requires TLS (h2c cleartext is non-standard)
  - Go's net/http enables H2 automatically over TLS
  - Check: log `r.Proto` to see "HTTP/2.0" vs "HTTP/1.1"
```

### Exercise: H2-Enabled Server

```go
package main

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "time"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        proto := r.Proto
        fmt.Fprintf(w, `{"protocol":%q,"path":%q}`, proto, r.URL.Path)
    })

    // TLS is required for H2; Go enables H2 automatically
    // To DISABLE H2: set TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler))
    srv := &http.Server{
        Addr:              ":8443",
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
        // H2 enabled automatically when serving TLS
    }

    fmt.Println("H1/H2 server on :8443")
    fmt.Println("H1 test: curl --http1.1 --cacert ca.crt https://localhost:8443/")
    fmt.Println("H2 test: curl --http2 --cacert ca.crt https://localhost:8443/")

    if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
        fmt.Println("ERROR:", err)
    }
}
```

```bash
# Compare H1 vs H2 latency (100 requests)
time curl --http1.1 --cacert ca.crt -s -o /dev/null https://localhost:8443/ &
time curl --http2 --cacert ca.crt -s -o /dev/null https://localhost:8443/
```

### Real-World Context
Go's `net/http` package enables HTTP/2 automatically when you call `ListenAndServeTLS`. Most teams don't notice — it just works. But your WAF middleware must work correctly with H2 (it does, because the WAF sees `http.Request` regardless of H1 or H2 underneath).

---

## Day 66 — WebSocket

### Learning Objectives
- Build a WebSocket server from the wire format (no library)
- Understand the upgrade handshake — why some WAF rules break WebSocket
- Implement real-time event broadcasting

### Key Concepts

```
WebSocket Upgrade:
  Client sends HTTP Upgrade request:
    GET /ws HTTP/1.1
    Host: example.com
    Upgrade: websocket
    Connection: Upgrade
    Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
    Sec-WebSocket-Version: 13

  Server responds with 101 Switching Protocols:
    HTTP/1.1 101 Switching Protocols
    Upgrade: websocket
    Connection: Upgrade
    Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
    (Accept = base64(SHA1(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11")))

  After this: raw WebSocket frames (not HTTP)

WebSocket Frame:
  [1B: FIN+RSV+Opcode][1B: MASK+Payload-len][0-8B: extended len][0-4B: mask][payload]
```

### Exercise: WebSocket Server (No Library)

```go
package main

import (
    "bufio"
    "crypto/sha1"
    "encoding/base64"
    "encoding/binary"
    "fmt"
    "io"
    "net"
    "net/http"
    "sync"
)

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func computeAcceptKey(key string) string {
    h := sha1.New()
    h.Write([]byte(key + wsGUID))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type WSConn struct {
    conn net.Conn
    mu   sync.Mutex
}

// ReadFrame reads one WebSocket frame
func (ws *WSConn) ReadFrame() (opcode byte, payload []byte, err error) {
    b := make([]byte, 2)
    if _, err := io.ReadFull(ws.conn, b); err != nil {
        return 0, nil, err
    }

    opcode = b[0] & 0x0F
    masked := (b[1] & 0x80) != 0
    length := int64(b[1] & 0x7F)

    switch length {
    case 126:
        ext := make([]byte, 2)
        io.ReadFull(ws.conn, ext)
        length = int64(binary.BigEndian.Uint16(ext))
    case 127:
        ext := make([]byte, 8)
        io.ReadFull(ws.conn, ext)
        length = int64(binary.BigEndian.Uint64(ext))
    }

    if length > 10*1024*1024 { // 10MB limit
        return 0, nil, fmt.Errorf("frame too large: %d bytes", length)
    }

    var mask [4]byte
    if masked {
        io.ReadFull(ws.conn, mask[:])
    }

    payload = make([]byte, length)
    io.ReadFull(ws.conn, payload)

    if masked {
        for i := range payload {
            payload[i] ^= mask[i%4]
        }
    }

    return opcode, payload, nil
}

// WriteFrame writes a text frame to the WebSocket connection
func (ws *WSConn) WriteFrame(opcode byte, payload []byte) error {
    ws.mu.Lock()
    defer ws.mu.Unlock()

    header := make([]byte, 2, 10)
    header[0] = 0x80 | opcode // FIN=1, opcode
    length := len(payload)
    switch {
    case length <= 125:
        header[1] = byte(length)
    case length <= 65535:
        header[1] = 126
        ext := make([]byte, 2)
        binary.BigEndian.PutUint16(ext, uint16(length))
        header = append(header, ext...)
    default:
        header[1] = 127
        ext := make([]byte, 8)
        binary.BigEndian.PutUint64(ext, uint64(length))
        header = append(header, ext...)
    }

    ws.conn.Write(header)
    _, err := ws.conn.Write(payload)
    return err
}

// Hub broadcasts messages to all connected WebSocket clients
type Hub struct {
    mu      sync.RWMutex
    clients map[*WSConn]bool
}

func (h *Hub) Add(ws *WSConn) {
    h.mu.Lock()
    h.clients[ws] = true
    h.mu.Unlock()
}

func (h *Hub) Remove(ws *WSConn) {
    h.mu.Lock()
    delete(h.clients, ws)
    h.mu.Unlock()
}

func (h *Hub) Broadcast(msg []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    for ws := range h.clients {
        go ws.WriteFrame(0x01, msg) // 0x01 = text frame
    }
}

var hub = &Hub{clients: make(map[*WSConn]bool)}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    // Validate WebSocket upgrade headers
    if r.Header.Get("Upgrade") != "websocket" {
        http.Error(w, "not a websocket upgrade", 400)
        return
    }

    key := r.Header.Get("Sec-WebSocket-Key")
    if key == "" {
        http.Error(w, "missing Sec-WebSocket-Key", 400)
        return
    }

    // Hijack the underlying TCP connection
    hijacker, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "hijack not supported", 500)
        return
    }
    conn, rw, err := hijacker.Hijack()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Send 101 Switching Protocols
    fmt.Fprintf(conn, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: %s\r\n\r\n",
        computeAcceptKey(key))

    ws := &WSConn{conn: conn}
    _ = rw
    hub.Add(ws)
    defer hub.Remove(ws)

    fmt.Printf("[WS] client connected: %s\n", conn.RemoteAddr())

    // Read loop
    for {
        opcode, payload, err := ws.ReadFrame()
        if err != nil {
            fmt.Printf("[WS] client disconnected: %v\n", err)
            return
        }
        switch opcode {
        case 0x01: // text
            fmt.Printf("[WS] message: %s\n", payload)
            hub.Broadcast(payload) // broadcast to all clients
        case 0x08: // close
            return
        case 0x09: // ping
            ws.WriteFrame(0x0A, payload) // pong
        }
    }
}

func main() {
    http.HandleFunc("/ws", wsHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, `<html><body>
<script>
var ws = new WebSocket("ws://localhost:8080/ws");
ws.onmessage = function(e) { document.body.innerHTML += "<p>" + e.data + "</p>"; };
ws.onopen = function() { ws.send("hello from browser"); };
</script></body></html>`)
    })

    fmt.Println("WebSocket server on :8080/ws")
    http.ListenAndServe(":8080", nil)
}
```

### Real-World Context
WAF rules that block `<script>` tags will break legitimate WebSocket applications that send HTML content. Understanding the WebSocket handshake (HTTP Upgrade → raw framing) explains why WAFs need separate WebSocket-aware rules. Your WAF (Week 10) must not inspect WebSocket frames as HTTP bodies.

---

## Day 67 — HTTP Security Headers

### Learning Objectives
- Implement a complete security headers middleware
- Build allowlist-based CORS (never wildcard for credentialed APIs)
- Generate per-request CSP nonces for inline script whitelisting

### Exercise: Complete Security Headers Middleware

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"
    "strings"
)

type CORSConfig struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
    AllowCredentials bool
    MaxAge          int
}

func generateNonce() string {
    b := make([]byte, 16)
    rand.Read(b)
    return base64.StdEncoding.EncodeToString(b)
}

// SecurityHeadersMiddleware sets all required security headers
func SecurityHeadersMiddleware(cors CORSConfig) func(http.Handler) http.Handler {
    allowedOriginSet := make(map[string]bool)
    for _, o := range cors.AllowedOrigins {
        allowedOriginSet[o] = true
    }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            h := w.Header()

            // Generate CSP nonce for this request (for inline scripts)
            nonce := generateNonce()

            // Content Security Policy
            h.Set("Content-Security-Policy",
                fmt.Sprintf("default-src 'self'; script-src 'self' 'nonce-%s'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; frame-ancestors 'none'", nonce))

            // Inject nonce into context for templates to use
            ctx := context.WithValue(r.Context(), cspNonceKey, nonce)

            // HSTS (only over HTTPS)
            if r.TLS != nil {
                h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
            }

            // Anti-clickjacking
            h.Set("X-Frame-Options", "DENY")

            // MIME sniffing prevention
            h.Set("X-Content-Type-Options", "nosniff")

            // Referrer policy
            h.Set("Referrer-Policy", "strict-origin-when-cross-origin")

            // Permissions policy (disable unused browser APIs)
            h.Set("Permissions-Policy", "accelerometer=(), camera=(), geolocation=(), microphone=(), payment=()")

            // CORS handling
            origin := r.Header.Get("Origin")
            if origin != "" {
                if allowedOriginSet[origin] {
                    h.Set("Access-Control-Allow-Origin", origin)
                    h.Set("Vary", "Origin") // critical: cache must vary by Origin
                    if cors.AllowCredentials {
                        h.Set("Access-Control-Allow-Credentials", "true")
                    }
                } else {
                    // Don't set CORS headers for disallowed origins
                    // Browser will block the request
                }

                // Preflight
                if r.Method == http.MethodOptions {
                    if allowedOriginSet[origin] {
                        h.Set("Access-Control-Allow-Methods", strings.Join(cors.AllowedMethods, ", "))
                        h.Set("Access-Control-Allow-Headers", strings.Join(cors.AllowedHeaders, ", "))
                        h.Set("Access-Control-Max-Age", fmt.Sprintf("%d", cors.MaxAge))
                    }
                    w.WriteHeader(http.StatusNoContent)
                    return
                }
            }

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func main() {
    corsConfig := CORSConfig{
        AllowedOrigins:   []string{"https://app.example.com", "https://admin.example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           3600,
    }

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        nonce := r.Context().Value(cspNonceKey).(string)
        fmt.Fprintf(w, `<html><body><script nonce="%s">console.log("CSP nonce works")</script></body></html>`, nonce)
    })

    srv := &http.Server{
        Addr:              ":8080",
        Handler:           SecurityHeadersMiddleware(corsConfig)(handler),
        ReadHeaderTimeout: 5 * time.Second,
    }

    fmt.Println("Security headers server on :8080")
    fmt.Println("Check headers: curl -I http://localhost:8080/")
    fmt.Println("Test CORS: curl -H 'Origin: https://evil.com' http://localhost:8080/")
    srv.ListenAndServe()
}
```

### Real-World Context
`Access-Control-Allow-Origin: *` on an authenticated API endpoint is OWASP A05 (critical). It lets any website make authenticated API requests using the victim's cookies/tokens. The `Vary: Origin` header is required to prevent CDN caching a response with the wrong origin header. CSP nonces prevent XSS in inline scripts.

---

## Week 9 Mini Project: HTTP Reverse Proxy with Full Observability

### Architecture Summary

```
Client Request
    ↓
┌──────────────────────────────────────────────────┐
│   HTTP Reverse Proxy (H1 + H2, TLS)              │
│                                                  │
│  Middleware Chain:                               │
│  1. PanicRecoverer                               │
│  2. RequestID                                    │
│  3. SecurityHeaders (CSP nonce, CORS, HSTS)      │
│  4. RequestLogger (structured JSON, slog)        │
│  5. UpstreamProxy → Backend health check         │
│                                                  │
│  Metrics Ring Buffer: p50/p95/p99 latency        │
│  Endpoints: /health /debug/vars /debug/pprof     │
└──────────────────────────────────────────────────┘
    ↓
Backend (HTTP :9090)
```

### Engineer Takeaway
Every Nginx, Envoy, and AWS ALB request goes through a chain like this. The middleware ordering is a security architecture decision: PanicRecoverer must be outermost to catch everything; RequestID before Logger so all logs are correlated. CORS misconfigurations are in the OWASP Top 10 every year — your allowlist CORS middleware eliminates the entire class of CORS wildcard vulnerabilities.
