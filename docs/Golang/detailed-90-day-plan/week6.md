# Week 6: TLS & Secure Transport (Days 38–44)

**Phase:** 2 — Networking Foundations + L4
**Goal:** Build TLS from the certificate level up. Generate certs in Go, enforce mTLS, implement HMAC request signing, and build the crypto primitives that your API security gateway uses.

---

> **Background in this wiki:** HTTPS guide [Ch 4 (TLS 1.3)](../../v2-https/real-life-guide-v1.md#chapter-4-tls-https-the-encrypted-tunnel) and [Ch 16 (certificate operations)](../../v2-https/real-life-guide-v1.md#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation); TCP/IP guide [Ch 29 (`tlsinspect`)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted) and [Ch 56 (`certreload`: zero-downtime rotation)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation); Go guide [§35](../real-life-golang-guide.md#35-tls-in-go-crypto-tls-mutual-tls) and [§50 (mTLS)](../real-life-golang-guide.md#50-mtls-client-certificate-authentication-end-to-end).

## Day 38 — TLS Handshake

### Learning Objectives
- Understand what happens during a TLS handshake (why cert expiry causes incidents)
- Inspect live TLS connections: version, cipher, expiry
- Build a cert expiry monitor — the tool that prevents your next outage

### Key Concepts

```
TLS 1.3 Handshake (simplified):
  Client → ClientHello (supported ciphers, TLS version, random)
  Server ← ServerHello (chosen cipher, TLS version, random)
  Server ← Certificate (server's cert chain)
  Server ← CertificateVerify (signature proves server has the private key)
  Server ← Finished
  Client → Finished
  ─────── Handshake complete (~1 RTT ───────
  Client ⟺ Server: encrypted application data

Why cert expiry causes outages:
  - Browser/client rejects expired cert → connection fails
  - No graceful fallback — hard failure
  - "Certificate expired" is in the top 5 causes of security incidents

Certificate fields you need to know:
  Subject     → who the cert is for (CommonName, SANs)
  Issuer      → who signed it (CA)
  NotBefore   → valid from
  NotAfter    → valid until (expiry)
  SANs        → Subject Alternative Names (DNS names and IPs the cert covers)
  KeyUsage    → what the cert can be used for
```

### Exercise: TLS Inspector + Cert Monitor

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "net"
    "time"
)

type CertInfo struct {
    Host        string
    TLSVersion  string
    CipherSuite string
    Subject     string
    Issuer      string
    NotAfter    time.Time
    SANs        []string
    DaysLeft    int
    Warning     string
}

func tlsVersionName(v uint16) string {
    switch v {
    case tls.VersionTLS10: return "TLS 1.0 (INSECURE)"
    case tls.VersionTLS11: return "TLS 1.1 (DEPRECATED)"
    case tls.VersionTLS12: return "TLS 1.2"
    case tls.VersionTLS13: return "TLS 1.3"
    default: return fmt.Sprintf("unknown(0x%x)", v)
    }
}

func inspectTLS(host string) (*CertInfo, error) {
    conf := &tls.Config{
        ServerName: host,
        // InsecureSkipVerify: false — NEVER skip verification in production
    }

    conn, err := tls.DialWithDialer(
        &net.Dialer{Timeout: 10 * time.Second},
        "tcp", host+":443", conf,
    )
    if err != nil {
        return nil, fmt.Errorf("TLS connect: %w", err)
    }
    defer conn.Close()

    state := conn.ConnectionState()
    cert := state.PeerCertificates[0] // leaf cert (server's cert)
    daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

    info := &CertInfo{
        Host:        host,
        TLSVersion:  tlsVersionName(state.Version),
        CipherSuite: tls.CipherSuiteName(state.CipherSuite),
        Subject:     cert.Subject.CommonName,
        Issuer:      cert.Issuer.CommonName,
        NotAfter:    cert.NotAfter,
        SANs:        cert.DNSNames,
        DaysLeft:    daysLeft,
    }

    if daysLeft < 7 {
        info.Warning = fmt.Sprintf("CRITICAL: cert expires in %d days!", daysLeft)
    } else if daysLeft < 30 {
        info.Warning = fmt.Sprintf("WARNING: cert expires in %d days", daysLeft)
    }

    return info, nil
}

func main() {
    hosts := []string{"example.com", "golang.org", "cloudflare.com"}
    for _, host := range hosts {
        info, err := inspectTLS(host)
        if err != nil {
            fmt.Printf("[%s] ERROR: %v\n", host, err)
            continue
        }
        fmt.Printf("[%s]\n", info.Host)
        fmt.Printf("  TLS:     %s\n", info.TLSVersion)
        fmt.Printf("  Cipher:  %s\n", info.CipherSuite)
        fmt.Printf("  Subject: %s\n", info.Subject)
        fmt.Printf("  Issuer:  %s\n", info.Issuer)
        fmt.Printf("  Expiry:  %s (%d days)\n", info.NotAfter.Format("2006-01-02"), info.DaysLeft)
        fmt.Printf("  SANs:    %v\n", info.SANs)
        if info.Warning != "" {
            fmt.Printf("  *** %s ***\n", info.Warning)
        }
        fmt.Println()
    }
}
```

### Real-World Context
Cert expiry is one of the top causes of security incidents. Facebook had a major outage in 2019 partly due to a cert issue. LinkedIn, Spotify, and many others have had cert expiry outages. Build a monitoring tool that checks every cert in your fleet daily and alerts at 30 days. This code is that tool.

---

## Day 39 — TLS Server in Go

### Learning Objectives
- Configure `tls.Config` to enforce security policy
- Require TLS 1.2+ — reject connections from legacy clients
- Log TLS metadata per connection for security auditing

### Key Concepts

```go
// tls.Config security settings
cfg := &tls.Config{
    MinVersion: tls.VersionTLS12,           // reject TLS 1.0, 1.1
    MaxVersion: tls.VersionTLS13,           // allow up to TLS 1.3

    // Cipher suite selection (TLS 1.2 only; TLS 1.3 ciphers are fixed)
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
        // NEVER include: RC4, DES, 3DES, NULL, EXPORT ciphers
    },

    // ConnState callback — called for each connection state transition
    // Use to log TLS metadata for security auditing
}

// NEVER in production:
// InsecureSkipVerify: true  — disables cert verification (OWASP A02)
```

### Exercise: Hardened TLS Echo Server

```go
package main

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "net"
    "time"
)

func buildTLSConfig(certFile, keyFile string) (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, fmt.Errorf("load cert: %w", err)
    }

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
        // Prefer TLS 1.3 (set automatically if both sides support it)

        // Called after TLS handshake completes — log security metadata
        VerifyConnection: func(cs tls.ConnectionState) error {
            fmt.Printf("[TLS] version=%s cipher=%s server_name=%s\n",
                tlsVersionName(cs.Version),
                tls.CipherSuiteName(cs.CipherSuite),
                cs.ServerName,
            )
            return nil
        },
    }, nil
}

func handleTLSConn(conn net.Conn) {
    defer conn.Close()

    tlsConn, ok := conn.(*tls.Conn)
    if !ok { return }

    // Force TLS handshake to complete (it's lazy by default)
    conn.SetDeadline(time.Now().Add(10 * time.Second))
    if err := tlsConn.Handshake(); err != nil {
        fmt.Printf("[TLS] handshake failed: %v\n", err)
        return
    }
    conn.SetDeadline(time.Time{}) // clear deadline after handshake

    state := tlsConn.ConnectionState()
    fmt.Printf("[CONN] %s via %s\n", conn.RemoteAddr(), tlsVersionName(state.Version))

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        fmt.Fprintf(conn, "ECHO: %s\n", scanner.Text())
    }
}

func main() {
    // Generate a self-signed cert for testing (see Day 40 for how)
    // For now, use: openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
    cfg, err := buildTLSConfig("cert.pem", "key.pem")
    if err != nil {
        fmt.Println("TLS config error:", err)
        return
    }

    ln, err := tls.Listen("tcp", ":8443", cfg)
    if err != nil {
        fmt.Println("listen error:", err)
        return
    }
    defer ln.Close()
    fmt.Println("TLS echo server on :8443")
    fmt.Println("Test: openssl s_client -connect localhost:8443")

    for {
        conn, err := ln.Accept()
        if err != nil { return }
        go handleTLSConn(conn)
    }
}
```

**Test with openssl:**
```bash
# Connect and inspect TLS details
openssl s_client -connect localhost:8443 -tls1_2
openssl s_client -connect localhost:8443 -tls1_3

# Try to connect with TLS 1.1 — should fail
openssl s_client -connect localhost:8443 -tls1_1
```

### Real-World Context
PCI DSS 4.0 requires TLS 1.2+ for all payment data transmissions. NIST SP 800-52 recommends TLS 1.2+ for government systems. `MinVersion: tls.VersionTLS12` is a one-line compliance control. Many teams skip this and get flagged in audits.

---

## Day 40 — Certificate Generation in Go

### Learning Objectives
- Generate a CA, server cert, and client cert entirely in Go
- Understand certificate fields: SANs, KeyUsage, ExtKeyUsage
- Build automated cert provisioning for dev/test environments

### Key Concepts

```
Certificate Chain:
  Root CA (self-signed, trust anchor)
    └── Server Cert (signed by CA, presented to clients)
    └── Client Cert (signed by CA, presented to server in mTLS)

Key fields:
  IsCA: true                    — can sign other certs
  KeyUsage: KeyCertSign         — allowed to sign certificates
  ExtKeyUsage: ServerAuth       — valid for server authentication
  ExtKeyUsage: ClientAuth       — valid for client authentication
  DNSNames: ["example.com"]     — Subject Alternative Names (SANs)
  IPAddresses: [127.0.0.1]      — IP SANs
```

### Exercise: Full PKI Generator

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "math/big"
    "net"
    "os"
    "time"
)

func generateKey() (*ecdsa.PrivateKey, error) {
    return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func savePEM(path, pemType string, data []byte) error {
    f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
    if err != nil { return err }
    defer f.Close()
    return pem.Encode(f, &pem.Block{Type: pemType, Bytes: data})
}

// generateCA creates a self-signed root CA certificate
func generateCA(commonName string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
    key, err := generateKey()
    if err != nil { return nil, nil, err }

    serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    template := &x509.Certificate{
        SerialNumber:          serial,
        Subject:               pkix.Name{CommonName: commonName, Organization: []string{"SecProxy CA"}},
        NotBefore:             time.Now().Add(-time.Minute), // backdate 1 min for clock skew
        NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
        IsCA:                  true,
        KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
        BasicConstraintsValid: true,
    }

    certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
    if err != nil { return nil, nil, err }

    cert, _ := x509.ParseCertificate(certDER)

    // Save CA cert and key
    savePEM("ca.crt", "CERTIFICATE", certDER)
    keyDER, _ := x509.MarshalECPrivateKey(key)
    savePEM("ca.key", "EC PRIVATE KEY", keyDER)

    fmt.Printf("Generated CA: %s (valid until %s)\n", commonName, template.NotAfter.Format("2006-01-02"))
    return cert, key, nil
}

// generateCert creates a certificate signed by the given CA
func generateCert(commonName string, isServer bool, sans []string, ips []net.IP, caCert *x509.Certificate, caKey *ecdsa.PrivateKey) error {
    key, err := generateKey()
    if err != nil { return err }

    serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

    extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
    if isServer {
        extKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
    }

    template := &x509.Certificate{
        SerialNumber: serial,
        Subject:      pkix.Name{CommonName: commonName},
        NotBefore:    time.Now().Add(-time.Minute),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        ExtKeyUsage:  extKeyUsage,
        DNSNames:     sans,
        IPAddresses:  ips,
    }

    certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, &key.PublicKey, caKey)
    if err != nil { return err }

    // Save cert and key
    prefix := "client"
    if isServer { prefix = "server" }
    savePEM(prefix+".crt", "CERTIFICATE", certDER)
    keyDER, _ := x509.MarshalECPrivateKey(key)
    savePEM(prefix+".key", "EC PRIVATE KEY", keyDER)

    kind := "client"
    if isServer { kind = "server" }
    fmt.Printf("Generated %s cert: %s (signed by CA)\n", kind, commonName)
    return nil
}

func main() {
    // Step 1: Generate root CA
    caCert, caKey, err := generateCA("SecProxy Root CA")
    if err != nil { fmt.Println("CA error:", err); return }

    // Step 2: Generate server cert with SANs
    err = generateCert("proxy.internal", true,
        []string{"proxy.internal", "localhost"},
        []net.IP{net.ParseIP("127.0.0.1")},
        caCert, caKey,
    )
    if err != nil { fmt.Println("server cert error:", err); return }

    // Step 3: Generate client cert (for mTLS)
    err = generateCert("service-backend", false,
        nil, nil, caCert, caKey,
    )
    if err != nil { fmt.Println("client cert error:", err); return }

    fmt.Println("\nGenerated files: ca.crt, ca.key, server.crt, server.key, client.crt, client.key")
    fmt.Println("Verify: openssl verify -CAfile ca.crt server.crt")
}
```

### Real-World Context
Internal service certificates are often managed manually with `openssl` scripts — fragile and audit-unfriendly. Go-based cert generation lets you automate PKI for dev/test, generate short-lived certs for ephemeral services, and build your own internal CA without depending on external tooling.

---

## Day 41 — Mutual TLS (mTLS)

### Learning Objectives
- Implement zero-trust service authentication: both sides verify each other
- Extract client identity from certificates — the foundation of your service mesh
- Understand why mTLS > API keys for service-to-service auth

### Key Concepts

```
Standard TLS (one-way):
  Server proves its identity to client via certificate
  Client is anonymous (authenticated by other means: API key, JWT, etc.)

mTLS (mutual):
  Server proves its identity to client via certificate
  Client ALSO proves its identity via certificate
  Both signed by a shared CA → only authorized clients can connect

Zero Trust with mTLS:
  No API keys to rotate, lose, or steal
  Certificate expiry enforces rotation automatically
  Identity is cryptographic, not a secret string
  Used by: Istio, Linkerd, Cloudflare Access, Google BeyondCorp
```

### Exercise: mTLS Server with Identity Extraction

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

func loadCACertPool(caFile string) (*x509.CertPool, error) {
    caPEM, err := os.ReadFile(caFile)
    if err != nil { return nil, err }
    pool := x509.NewCertPool()
    if !pool.AppendCertsFromPEM(caPEM) {
        return nil, fmt.Errorf("failed to parse CA cert from %s", caFile)
    }
    return pool, nil
}

func buildMTLSServerConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil { return nil, err }

    caPool, err := loadCACertPool(caFile)
    if err != nil { return nil, err }

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caPool,
        ClientAuth:   tls.RequireAndVerifyClientCert, // enforce mTLS

        MinVersion: tls.VersionTLS12,

        // Called after mTLS handshake — extract client identity
        VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
            if len(verifiedChains) == 0 {
                return fmt.Errorf("no verified certificate chain")
            }
            // Additional custom checks beyond standard verification
            clientCert := verifiedChains[0][0]
            fmt.Printf("[mTLS] client identity: %s\n", clientCert.Subject.CommonName)
            return nil
        },
    }, nil
}

func clientIdentityFromRequest(r *http.Request) string {
    if r.TLS == nil || len(r.TLS.VerifiedChains) == 0 {
        return ""
    }
    return r.TLS.VerifiedChains[0][0].Subject.CommonName
}

func main() {
    tlsConfig, err := buildMTLSServerConfig("server.crt", "server.key", "ca.crt")
    if err != nil {
        fmt.Println("TLS config:", err)
        return
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        identity := clientIdentityFromRequest(r)
        if identity == "" {
            http.Error(w, "missing client certificate", http.StatusForbidden)
            return
        }

        // Inject identity into response for downstream services
        w.Header().Set("X-Client-Identity", identity)
        fmt.Fprintf(w, `{"authenticated_as":%q,"path":%q}`, identity, r.URL.Path)
    })

    srv := &http.Server{
        Addr:      ":8443",
        Handler:   mux,
        TLSConfig: tlsConfig,
    }

    fmt.Println("mTLS server on :8443")
    fmt.Println("Test: curl --cert client.crt --key client.key --cacert ca.crt https://localhost:8443/")
    srv.ListenAndServeTLS("", "") // certs already in TLSConfig
}
```

**Test:**
```bash
# First generate certs (Day 40)
go run ../day40-cert-gen/main.go

# Test with valid client cert — should succeed
curl --cert client.crt --key client.key --cacert ca.crt https://localhost:8443/

# Test without client cert — should fail (403)
curl --cacert ca.crt https://localhost:8443/

# Test with wrong CA client cert — should fail (TLS handshake error)
```

### Real-World Context
mTLS is how Google's zero-trust BeyondCorp works. Instead of VPN (trust everything inside the network), each service has a cert. Compromising one service doesn't compromise others — the attacker gets that service's cert, not network-level access. Istio automates mTLS between all pods in Kubernetes.

---

## Day 42 — Hashing & HMAC

### Learning Objectives
- Implement HMAC-SHA256 request signing — the pattern used by AWS SigV4 and webhook validation
- Understand timing attacks and why `subtle.ConstantTimeCompare` is required
- Build a webhook validation system (OWASP API Security — A02 Broken Authentication)

### Key Concepts

```go
// Hashing — one-way, no key
hash := sha256.Sum256([]byte("data"))

// HMAC — keyed hash; proves the sender knows the shared secret
mac := hmac.New(sha256.New, []byte("secret-key"))
mac.Write([]byte("message"))
signature := mac.Sum(nil)

// TIMING ATTACK: bytes.Equal reveals how many bytes match
// Attacker can guess HMAC byte-by-byte by measuring response time
if bytes.Equal(computed, received) { }  // WRONG

// SAFE: constant-time compare — always takes the same time
if subtle.ConstantTimeCompare(computed, received) == 1 { }  // RIGHT
```

### Exercise: Webhook Signing System

```go
package main

import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "crypto/subtle"
    "encoding/hex"
    "fmt"
    "net/http"
    "strconv"
    "time"
)

const webhookSecret = "your-webhook-secret-key-here"

// SignWebhook adds HMAC-SHA256 signature to outgoing webhook
// Signature = HMAC-SHA256(key, "timestamp.body")
func SignWebhook(body []byte, secret string) (signature, timestamp string) {
    ts := strconv.FormatInt(time.Now().Unix(), 10)

    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(ts + "."))
    mac.Write(body)

    return "sha256=" + hex.EncodeToString(mac.Sum(nil)), ts
}

// VerifyWebhook validates an incoming webhook signature
// Returns error if invalid or replayed (timestamp > 5min old)
func VerifyWebhook(body []byte, signature, timestamp, secret string) error {
    // Check timestamp freshness — prevent replay attacks
    ts, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return fmt.Errorf("invalid timestamp")
    }
    age := time.Since(time.Unix(ts, 0))
    if age > 5*time.Minute || age < -time.Minute {
        return fmt.Errorf("webhook timestamp too old or too far in future: age=%v", age)
    }

    // Compute expected signature
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(timestamp + "."))
    mac.Write(body)
    expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

    // Constant-time compare — prevents timing attacks
    // Without this, an attacker can measure response time to guess signature byte-by-byte
    if subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) != 1 {
        return fmt.Errorf("signature mismatch")
    }
    return nil
}

// WebhookSender sends a signed webhook
func WebhookSender(targetURL string, payload []byte) error {
    sig, ts := SignWebhook(payload, webhookSecret)

    req, _ := http.NewRequest("POST", targetURL, bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Signature", sig)
    req.Header.Set("X-Webhook-Timestamp", ts)

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("webhook rejected: %d", resp.StatusCode)
    }
    return nil
}

// WebhookReceiver validates incoming webhooks
func WebhookReceiver(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB limit
    if err != nil {
        http.Error(w, "read error", 500)
        return
    }

    sig := r.Header.Get("X-Webhook-Signature")
    ts := r.Header.Get("X-Webhook-Timestamp")

    if err := VerifyWebhook(body, sig, ts, webhookSecret); err != nil {
        // Log the attempt — could be an attack
        fmt.Printf("WEBHOOK INVALID: %v from %s\n", err, r.RemoteAddr)
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    fmt.Printf("WEBHOOK VALID: %s\n", body)
    w.WriteHeader(http.StatusNoContent)
}

func main() {
    http.HandleFunc("/webhook", WebhookReceiver)

    fmt.Println("Webhook receiver on :8080/webhook")
    // Test: go to Day 42 exercise to send test webhooks
    http.ListenAndServe(":8080", nil)
}
```

### Real-World Context
GitHub, Stripe, Shopify, and Twilio all use HMAC-SHA256 webhook signing with this exact pattern. When your team integrates with any webhook-based service, verify the signature. Without it, an attacker can forge events: "order shipped", "payment received", "user admin escalated".

---

## Day 43 — Symmetric Encryption (AES-GCM)

### Learning Objectives
- Encrypt sensitive data at rest with AES-256-GCM
- Understand why GCM > CBC (authenticated encryption)
- Build an encrypted audit log — PII encrypted in transit and at rest

### Key Concepts

```
AES-GCM (Galois/Counter Mode):
  - Authenticated encryption: provides both confidentiality AND integrity
  - If ciphertext is tampered with, Decrypt() returns an error
  - Requires a unique nonce for every encryption (never reuse!)
  - 96-bit (12-byte) nonce: standard for GCM
  - 256-bit (32-byte) key: AES-256

AES-CBC (Cipher Block Chaining):
  - Older mode — no authentication (attacker can flip bits without detection)
  - Requires padding (PKCS#7) — vulnerable to padding oracle attacks
  - Avoid for new designs

Key derivation:
  - Never use passwords directly as AES keys
  - Use PBKDF2, scrypt, or Argon2 to derive a key from a password
  - For service-to-service: use a random 32-byte key from crypto/rand
```

### Exercise: Encrypted Audit Log

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "time"
)

// AuditRecord represents a security event
type AuditRecord struct {
    Time      time.Time `json:"time"`
    RequestID string    `json:"request_id"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Resource  string    `json:"resource"`
    IP        string    `json:"client_ip"`
    Result    string    `json:"result"` // allow, deny, block
}

// EncryptedAuditLog writes AES-256-GCM encrypted audit records
type EncryptedAuditLog struct {
    key  []byte   // 32 bytes (AES-256)
    file *os.File
}

func NewEncryptedAuditLog(path string, keyHex string) (*EncryptedAuditLog, error) {
    // Derive key from hex string (in production: load from secrets manager)
    key := make([]byte, 32)
    if n, err := fmt.Sscanf(keyHex, "%x", &key); n != 1 || err != nil {
        // Generate a random key if not provided
        if _, err := io.ReadFull(rand.Reader, key); err != nil {
            return nil, err
        }
        fmt.Printf("Generated key: %x\n", key)
    }

    f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil { return nil, err }

    return &EncryptedAuditLog{key: key, file: f}, nil
}

// Write encrypts and appends an audit record
func (l *EncryptedAuditLog) Write(record AuditRecord) error {
    // Serialize record
    plaintext, err := json.Marshal(record)
    if err != nil { return err }

    // Create AES-256-GCM cipher
    block, err := aes.NewCipher(l.key)
    if err != nil { return err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return err }

    // Generate random nonce (must be unique per encryption)
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return err }

    // Encrypt: nonce + ciphertext + GCM auth tag
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

    // Write length-prefixed encrypted record
    // Format: [4 bytes length][nonce+ciphertext+tag]
    length := make([]byte, 4)
    binary.BigEndian.PutUint32(length, uint32(len(ciphertext)))
    l.file.Write(length)
    l.file.Write(ciphertext)
    return nil
}

// ReadAll decrypts all audit records from the log file
func ReadAll(path string, key []byte) ([]AuditRecord, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close()

    block, err := aes.NewCipher(key)
    if err != nil { return nil, err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return nil, err }

    var records []AuditRecord
    lengthBuf := make([]byte, 4)
    for {
        if _, err := io.ReadFull(f, lengthBuf); err != nil {
            if err == io.EOF { break }
            return nil, err
        }
        length := binary.BigEndian.Uint32(lengthBuf)
        ciphertext := make([]byte, length)
        if _, err := io.ReadFull(f, ciphertext); err != nil { return nil, err }

        nonceSize := gcm.NonceSize()
        if len(ciphertext) < nonceSize { return nil, fmt.Errorf("ciphertext too short") }
        nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

        plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
        if err != nil { return nil, fmt.Errorf("decrypt failed (tampered?): %w", err) }

        var record AuditRecord
        if err := json.Unmarshal(plaintext, &record); err != nil { return nil, err }
        records = append(records, record)
    }
    return records, nil
}

func main() {
    log, err := NewEncryptedAuditLog("audit.enc", "")
    if err != nil { fmt.Println("ERROR:", err); return }

    records := []AuditRecord{
        {Time: time.Now(), RequestID: "req-001", UserID: "user:42", Action: "login", Resource: "/auth/login", IP: "10.0.0.1", Result: "allow"},
        {Time: time.Now(), RequestID: "req-002", UserID: "", Action: "access", Resource: "/api/admin", IP: "203.0.113.1", Result: "deny"},
    }

    for _, r := range records {
        if err := log.Write(r); err != nil {
            fmt.Println("write error:", err)
        }
    }
    fmt.Println("Wrote encrypted audit log: audit.enc")
    fmt.Println("Content is encrypted — cat audit.enc shows only ciphertext")
}
```

### Real-World Context
HIPAA requires encryption of PHI at rest. PCI DSS requires encryption of cardholder data. GDPR requires appropriate technical measures for PII. AES-256-GCM with a random nonce per record is the correct pattern. The GCM auth tag detects tampering — if someone modifies the log file, decryption fails.

---

## Day 44 — Asymmetric Crypto & Token Signing

### Learning Objectives
- Build a JWT-like token system using ECDSA P-256 — no JWT library needed
- Understand why ECDSA P-256 > RSA-2048 for tokens (smaller, faster, same security)
- Implement the foundation of your API gateway's authentication

### Key Concepts

```
Symmetric (HMAC): same key signs and verifies → key must be shared with verifier
  Use for: webhook signing (you control both sides), session cookies

Asymmetric (ECDSA): private key signs, public key verifies → public key shareable
  Use for: tokens (anyone can verify, only you can issue)
  ECDSA P-256: 256-bit key = ~128-bit security = equivalent to RSA-3072
  RSA-2048: use for TLS certs (protocol requirement); prefer ECDSA for tokens

Token structure (custom, no JWT library):
  Header:  {"alg":"ES256","typ":"TOKEN"}
  Claims:  {"sub":"user:42","exp":1234567890,"iat":1234567890,"roles":["admin"]}
  Signature: ECDSA.Sign(SHA256(base64(header).base64(claims)))
  Token: base64url(header).base64url(claims).base64url(signature)
```

### Exercise: ECDSA Token System

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "math/big"
    "strings"
    "time"
)

// TokenClaims are the payload of a token
type TokenClaims struct {
    Subject   string   `json:"sub"`
    IssuedAt  int64    `json:"iat"`
    ExpiresAt int64    `json:"exp"`
    Roles     []string `json:"roles"`
}

type TokenService struct {
    privateKey *ecdsa.PrivateKey
    publicKey  *ecdsa.PublicKey
}

func NewTokenService() (*TokenService, error) {
    key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil { return nil, err }
    return &TokenService{privateKey: key, publicKey: &key.PublicKey}, nil
}

func b64(data []byte) string {
    return base64.RawURLEncoding.EncodeToString(data)
}

func unb64(s string) ([]byte, error) {
    return base64.RawURLEncoding.DecodeString(s)
}

// Issue signs and returns a token valid for 1 hour
func (ts *TokenService) Issue(subject string, roles []string) (string, error) {
    now := time.Now().Unix()
    claims := TokenClaims{
        Subject:   subject,
        IssuedAt:  now,
        ExpiresAt: now + 3600,
        Roles:     roles,
    }

    header := b64([]byte(`{"alg":"ES256","typ":"TOKEN"}`))
    claimsJSON, err := json.Marshal(claims)
    if err != nil { return "", err }
    payload := header + "." + b64(claimsJSON)

    // Sign SHA256 of the payload
    hash := sha256.Sum256([]byte(payload))
    r, s, err := ecdsa.Sign(rand.Reader, ts.privateKey, hash[:])
    if err != nil { return "", err }

    // Encode signature as r||s (each 32 bytes for P-256)
    sig := make([]byte, 64)
    rBytes := r.Bytes()
    sBytes := s.Bytes()
    copy(sig[32-len(rBytes):32], rBytes)
    copy(sig[64-len(sBytes):64], sBytes)

    return payload + "." + b64(sig), nil
}

// Verify validates a token and returns claims
func (ts *TokenService) Verify(token string) (*TokenClaims, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid token format")
    }

    // Verify signature
    payload := parts[0] + "." + parts[1]
    hash := sha256.Sum256([]byte(payload))

    sigBytes, err := unb64(parts[2])
    if err != nil || len(sigBytes) != 64 {
        return nil, fmt.Errorf("invalid signature encoding")
    }
    r := new(big.Int).SetBytes(sigBytes[:32])
    s := new(big.Int).SetBytes(sigBytes[32:])

    if !ecdsa.Verify(ts.publicKey, hash[:], r, s) {
        return nil, fmt.Errorf("invalid signature")
    }

    // Decode claims
    claimsJSON, err := unb64(parts[1])
    if err != nil { return nil, fmt.Errorf("invalid claims encoding") }

    var claims TokenClaims
    if err := json.Unmarshal(claimsJSON, &claims); err != nil {
        return nil, fmt.Errorf("invalid claims JSON")
    }

    // Check expiry
    if time.Now().Unix() > claims.ExpiresAt {
        return nil, fmt.Errorf("token expired")
    }

    return &claims, nil
}

func main() {
    svc, err := NewTokenService()
    if err != nil { fmt.Println("ERROR:", err); return }

    // Issue a token
    token, err := svc.Issue("user:42", []string{"read", "write"})
    if err != nil { fmt.Println("issue error:", err); return }
    fmt.Printf("Token: %s\n\n", token[:50]+"...")

    // Verify the token
    claims, err := svc.Verify(token)
    if err != nil { fmt.Println("verify error:", err); return }
    fmt.Printf("Verified: subject=%s roles=%v expires=%s\n",
        claims.Subject, claims.Roles,
        time.Unix(claims.ExpiresAt, 0).Format(time.RFC3339),
    )

    // Tamper with the token
    tampered := token[:len(token)-5] + "XXXXX"
    _, err = svc.Verify(tampered)
    fmt.Printf("Tampered token: %v\n", err)
}
```

### Real-World Context
This is the token system your API gateway uses in Day 75 and Week 11. ECDSA P-256 tokens are ~60% smaller than RSA-2048 tokens (important for mobile clients and high-frequency API calls). The private key stays on your auth server; the public key can be distributed to all services that need to verify tokens.

---

## Week 6 Mini Project: mTLS Service Gateway

Build a complete service gateway: mTLS client authentication + ECDSA token issuance + upstream forwarding.

```go
// gateway/main.go
package main

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "time"
)

func main() {
    // Load PKI (generated in Day 40)
    tlsConfig, err := buildMTLSServerConfig("server.crt", "server.key", "ca.crt")
    if err != nil { fmt.Println(err); return }

    tokenService, err := NewTokenService() // from Day 44
    if err != nil { fmt.Println(err); return }

    upstream, _ := url.Parse("http://127.0.0.1:9090")
    proxy := httputil.NewSingleHostReverseProxy(upstream)

    mux := http.NewServeMux()

    // Auth endpoint: validates mTLS, issues ECDSA token
    mux.HandleFunc("POST /auth/token", func(w http.ResponseWriter, r *http.Request) {
        identity := clientIdentityFromRequest(r)
        if identity == "" {
            http.Error(w, "client certificate required", http.StatusUnauthorized)
            return
        }

        token, err := tokenService.Issue(identity, []string{"read", "write"})
        if err != nil {
            http.Error(w, "token issue error", 500)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"token":%q,"expires_in":3600}`, token)
    })

    // Protected API: validate ECDSA token, inject identity header, proxy
    mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
            http.Error(w, "authorization required", http.StatusUnauthorized)
            return
        }

        claims, err := tokenService.Verify(authHeader[7:])
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        // Inject verified identity for upstream service
        r.Header.Set("X-Authenticated-Subject", claims.Subject)
        r.Header.Del("Authorization") // don't forward token to upstream

        proxy.ServeHTTP(w, r)
    })

    srv := &http.Server{
        Addr:              ":8443",
        Handler:           mux,
        TLSConfig:         tlsConfig,
        ReadHeaderTimeout: 5 * time.Second,
    }

    fmt.Println("mTLS gateway on :8443")
    fmt.Println("1. Get token: curl --cert client.crt --key client.key --cacert ca.crt https://localhost:8443/auth/token -X POST")
    fmt.Println("2. Use token: curl --cacert ca.crt -H 'Authorization: Bearer <token>' https://localhost:8443/api/data")

    srv.ListenAndServeTLS("", "")
}
```

### Engineer Takeaway
This is the core of a zero-trust service-to-service auth system. Istio, SPIFFE, and AWS Private CA all implement this pattern. You've built it from scratch in ~200 lines. The `X-Authenticated-Subject` header is how your backend services know *who* is calling them — without managing API keys or shared secrets.
