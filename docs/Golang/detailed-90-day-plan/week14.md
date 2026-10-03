# Week 14: Bot Detection (Days 99–105)

> **Phase 4 continuation:** Add bot detection to the DDoS mitigation layer from Week 13. Bot traffic accounts for 40–50% of internet traffic. This week you build the signal framework used by DataDome, PerimeterX, and Cloudflare Bot Management.

**Why this week matters:** Bots scrape prices, credential-stuff accounts, abuse APIs, and skew analytics. Bot detection is a multi-signal problem — no single signal is reliable. This week you build a scoring system that aggregates HTTP fingerprinting, TLS fingerprinting, session behavior, and IP reputation.

---

> **Background in this wiki:** HTTPS guide [Ch 12 (security across the lifecycle)](../../v2-https/real-life-guide-v1.md#chapter-12-security-protecting-the-entire-lifecycle) and [Ch 21 (rate limiting and abuse controls)](../../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers). Security in Depth: [Chapter 40D (bots and edge enforcement)](../../security/real-life-security-guide-v1.md#chapter-40d-waap-bots-api-discovery-and-edge-enforcement) and [Chapter 74 (abuse and fraud)](../../security/real-life-security-guide-v1.md#chapter-74-abuse-fraud-bots-and-business-logic-defense).

## Day 99 — Bot Taxonomy & Signal Framework

| | |
|---|---|
| **Learning objective** | Design a pluggable signal-scoring system; classify bot types |
| **Key packages** | `sync`, `sync/atomic`, `time`, `net/http` |
| **Security context** | Different bot types require different mitigations — misclassification causes false positives |

### Key Concepts

**Bot taxonomy:**

| Type | Example | Detection signal | Mitigation |
|------|---------|-----------------|------------|
| Simple HTTP client | `curl`, `python-requests` | User-Agent, missing headers | PoW challenge |
| Headless browser | Puppeteer/Playwright | JA3 fingerprint, timing | JS challenge (browser required) |
| Residential proxy | Commercial bot farms | Session velocity, referrer | Rate limit + challenge |
| Human-like bot | Trained ML bots | Behavioral patterns over time | Honeypot, shadow-ban |
| Good bot | Googlebot, Bingbot | User-Agent + rDNS verify | Allowlist |

**Signal categories:**
- **Network signals:** IP ASN, datacenter vs residential, TLS fingerprint (JA3)
- **HTTP signals:** Header order/values, User-Agent consistency, Accept-Encoding
- **Behavioral signals:** Request velocity, path diversity, referrer chain, session depth

### Exercise — BotScorer Signal Framework

```go
package bot

import (
    "net/http"
    "sync"
    "time"
)

// SignalCategory groups related signals for reporting
type SignalCategory string

const (
    CategoryNetwork   SignalCategory = "network"
    CategoryHTTP      SignalCategory = "http"
    CategoryBehavior  SignalCategory = "behavioral"
    CategoryReputation SignalCategory = "reputation"
)

// Signal represents a single detection signal with its weight
type Signal struct {
    Name     string
    Weight   float64
    Category SignalCategory
    Reason   string // human-readable explanation
}

// Common signals
var (
    SigMissingAcceptLang   = Signal{"missing_accept_language", 10, CategoryHTTP, "Real browsers always send Accept-Language"}
    SigBadAcceptEncoding   = Signal{"bad_accept_encoding", 15, CategoryHTTP, "Accept-Encoding inconsistent with claimed browser"}
    SigDatacenterIP        = Signal{"datacenter_ip", 15, CategoryNetwork, "Source IP belongs to known hosting/datacenter ASN"}
    SigKnownBadJA3         = Signal{"known_bad_ja3", 30, CategoryNetwork, "TLS fingerprint matches curl/requests/scrapy"}
    SigHighVelocity        = Signal{"high_velocity", 20, CategoryBehavior, "Request rate > 10 req/s"}
    SigLowPathDiversity    = Signal{"low_path_diversity", 10, CategoryBehavior, "< 2 unique paths in session"}
    SigMissingReferrer     = Signal{"missing_referrer", 10, CategoryBehavior, "> 80% of requests have no Referer"}
    SigHoneypotHit         = Signal{"honeypot_hit", 100, CategoryBehavior, "Request to known honeypot endpoint"}
    SigUserAgentMismatch   = Signal{"user_agent_mismatch", 25, CategoryHTTP, "UA claims browser but headers don't match"}
)

// SessionScore tracks the aggregated score for a session/IP
type SessionScore struct {
    mu      sync.Mutex
    signals []Signal
    score   float64
    created time.Time
    lastSee time.Time
}

func (ss *SessionScore) AddSignal(s Signal) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    ss.signals = append(ss.signals, s)
    ss.score += s.Weight
    ss.lastSee = time.Now()
}

func (ss *SessionScore) Score() float64 {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    return ss.score
}

func (ss *SessionScore) Signals() []Signal {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    out := make([]Signal, len(ss.signals))
    copy(out, ss.signals)
    return out
}

// BotScorer maintains per-IP scores and dispatches detection signals
type BotScorer struct {
    mu       sync.RWMutex
    sessions map[string]*SessionScore
    detectors []Detector
    done     chan struct{}

    // Thresholds
    ChallengeScore float64 // score >= this → PoW challenge
    BlockScore     float64 // score >= this → block
}

// Detector is the interface for individual signal detectors
type Detector interface {
    Detect(r *http.Request, session *SessionScore)
}

func NewBotScorer(detectors ...Detector) *BotScorer {
    bs := &BotScorer{
        sessions:       make(map[string]*SessionScore),
        detectors:      detectors,
        ChallengeScore: 40,
        BlockScore:     70,
        done:           make(chan struct{}),
    }
    go bs.cleanup()
    return bs
}

func (bs *BotScorer) sessionFor(ip string) *SessionScore {
    bs.mu.RLock()
    s, ok := bs.sessions[ip]
    bs.mu.RUnlock()
    if ok {
        return s
    }
    bs.mu.Lock()
    defer bs.mu.Unlock()
    s, ok = bs.sessions[ip]
    if !ok {
        s = &SessionScore{created: time.Now(), lastSee: time.Now()}
        bs.sessions[ip] = s
    }
    return s
}

func (bs *BotScorer) Score(r *http.Request) (float64, []Signal) {
    ip := clientIP(r)
    session := bs.sessionFor(ip)

    for _, d := range bs.detectors {
        d.Detect(r, session)
    }

    return session.Score(), session.Signals()
}

// Action returns the mitigation action based on score
func (bs *BotScorer) Action(score float64) string {
    switch {
    case score >= bs.BlockScore:
        return "block"
    case score >= bs.ChallengeScore:
        return "challenge"
    default:
        return "allow"
    }
}

func (bs *BotScorer) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        score, _ := bs.Score(r)
        switch bs.Action(score) {
        case "block":
            http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
        case "challenge":
            // In production: redirect to PoW challenge page
            w.Header().Set("X-Bot-Score", fmt.Sprintf("%.0f", score))
            http.Error(w, `{"error":"challenge_required"}`, http.StatusTooManyRequests)
        default:
            next.ServeHTTP(w, r)
        }
    })
}

func (bs *BotScorer) cleanup() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-bs.done:
            return
        case <-ticker.C:
            cutoff := time.Now().Add(-30 * time.Minute)
            bs.mu.Lock()
            for ip, s := range bs.sessions {
                s.mu.Lock()
                stale := s.lastSee.Before(cutoff)
                s.mu.Unlock()
                if stale {
                    delete(bs.sessions, ip)
                }
            }
            bs.mu.Unlock()
        }
    }
}

func clientIP(r *http.Request) string {
    host, _, _ := net.SplitHostPort(r.RemoteAddr)
    return host
}
```

---

## Day 100 — HTTP Header Fingerprinting

| | |
|---|---|
| **Learning objective** | Score requests based on header consistency vs claimed User-Agent |
| **Key packages** | `net/http`, `strings` |
| **Security context** | Real browsers send specific headers in a consistent order; bots don't bother |

### Key Concepts

**Chrome 120 always sends these headers (in this order):**
```
:method, :authority, :scheme, :path  (HTTP/2 pseudo-headers)
sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform
Upgrade-Insecure-Requests
User-Agent
Accept
Sec-Fetch-Site, Sec-Fetch-Mode, Sec-Fetch-Dest
Referer (when navigating)
Accept-Encoding: gzip, deflate, br
Accept-Language
```

**curl sends:**
```
User-Agent, Accept, Host (only)
```

### Exercise — Header Consistency Detector

```go
package bot

import (
    "net/http"
    "strings"
)

// browserProfile defines expected headers for a claimed browser
type browserProfile struct {
    requiredHeaders  []string // must be present
    forbiddenValues  map[string]string // header → value that browser never sends
    minHeaderCount   int
}

var profiles = map[string]browserProfile{
    "Chrome": {
        requiredHeaders: []string{"Accept-Language", "Accept-Encoding"},
        forbiddenValues: map[string]string{
            "Accept-Encoding": "compress", // Chrome never sends compress
        },
        minHeaderCount: 6,
    },
    "Firefox": {
        requiredHeaders: []string{"Accept-Language", "Accept-Encoding"},
        forbiddenValues: map[string]string{},
        minHeaderCount: 5,
    },
    "Safari": {
        requiredHeaders: []string{"Accept-Language"},
        forbiddenValues: map[string]string{},
        minHeaderCount: 5,
    },
}

// claimedBrowser extracts the browser name from User-Agent
func claimedBrowser(ua string) string {
    ua = strings.ToLower(ua)
    switch {
    case strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg"):
        return "Chrome"
    case strings.Contains(ua, "firefox"):
        return "Firefox"
    case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
        return "Safari"
    case strings.Contains(ua, "curl"):
        return "curl"
    case strings.Contains(ua, "python"):
        return "python"
    case strings.Contains(ua, "go-http-client"):
        return "go-http"
    }
    return ""
}

// HeaderConsistencyDetector implements Detector
type HeaderConsistencyDetector struct{}

func (d *HeaderConsistencyDetector) Detect(r *http.Request, session *SessionScore) {
    ua := r.Header.Get("User-Agent")
    if ua == "" {
        session.AddSignal(Signal{"missing_ua", 20, CategoryHTTP, "No User-Agent header"})
        return
    }

    browser := claimedBrowser(ua)

    // Known non-browser clients — flag immediately
    switch browser {
    case "curl":
        session.AddSignal(Signal{"curl_ua", 25, CategoryHTTP, "curl User-Agent detected"})
        return
    case "python":
        session.AddSignal(Signal{"python_ua", 25, CategoryHTTP, "Python http client detected"})
        return
    case "go-http":
        session.AddSignal(Signal{"go_http_ua", 20, CategoryHTTP, "Go HTTP client detected"})
        return
    case "":
        session.AddSignal(Signal{"unknown_ua", 10, CategoryHTTP, "Unrecognized User-Agent"})
    }

    profile, ok := profiles[browser]
    if !ok {
        return
    }

    // Check required headers
    for _, h := range profile.requiredHeaders {
        if r.Header.Get(h) == "" {
            session.AddSignal(SigMissingAcceptLang)
        }
    }

    // Check forbidden values
    for header, forbidden := range profile.forbiddenValues {
        if strings.Contains(r.Header.Get(header), forbidden) {
            session.AddSignal(SigBadAcceptEncoding)
        }
    }

    // Check minimum header count
    if len(r.Header) < profile.minHeaderCount {
        session.AddSignal(Signal{
            "low_header_count", 15, CategoryHTTP,
            fmt.Sprintf("%s normally sends ≥%d headers, got %d", browser, profile.minHeaderCount, len(r.Header)),
        })
    }

    // Check Accept header consistency
    accept := r.Header.Get("Accept")
    if browser == "Chrome" && accept == "*/*" {
        // Chrome never sends Accept: */* for page navigation
        session.AddSignal(Signal{"bad_accept", 10, CategoryHTTP, "Chrome never sends Accept: */*"})
    }

    // User-Agent vs Sec-CH-UA consistency (HTTP/2 only)
    secCHUA := r.Header.Get("Sec-CH-UA")
    if browser == "Chrome" && r.ProtoMajor >= 2 && secCHUA == "" {
        session.AddSignal(SigUserAgentMismatch)
    }
}
```

---

## Day 101 — TLS Fingerprinting (JA3)

| | |
|---|---|
| **Learning objective** | Compute JA3 fingerprint from TLS ClientHello via `tls.ClientHelloInfo` |
| **Key packages** | `crypto/tls`, `crypto/md5`, `fmt`, `strings` |
| **Security context** | JA3 fingerprints the TLS stack, not the application — hard for bots to spoof without modifying TLS library |

### Key Concepts

**JA3 algorithm:**
1. Extract from ClientHello: `SSLVersion`, `Ciphers`, `Extensions`, `EllipticCurves`, `EllipticCurvePointFormats`
2. Join each list with `-` and sort (for extensions, preserve order)
3. Concatenate fields with `,`
4. MD5 hash the result

**Known JA3 hashes:**

| Client | JA3 Hash (representative) |
|--------|--------------------------|
| Chrome 120 | `b32309a26951912be7dba376398abc3b` |
| Firefox 120 | `a0e9f5d64349fb13191bc781f81f42e1` |
| curl 8.x | `7dc465e28f1d46d873bf9c2ae6d05e43` |
| Python requests | `6597c339e4df3929c3cff7cc61c2b1c9` |
| Scrapy | `0700c0e4ebf0d82de6a30e0c28d21c0e` |

### Exercise — JA3 Fingerprinter

```go
package bot

import (
    "crypto/md5"
    "crypto/tls"
    "fmt"
    "strings"
)

// JA3Fingerprint computes a JA3 hash from tls.ClientHelloInfo
func JA3Fingerprint(hello *tls.ClientHelloInfo) string {
    // SSLVersion (always TLS — represented as 0x0303 for TLS 1.2, etc.)
    // In Go's ClientHelloInfo, we use SupportedVersions
    version := uint16(0x0303) // TLS 1.2 as base
    if len(hello.SupportedVersions) > 0 {
        version = hello.SupportedVersions[0]
    }

    // Ciphers — exclude GREASE values (0xXaXa pattern)
    var ciphers []string
    for _, c := range hello.CipherSuites {
        if !isGREASE(c) {
            ciphers = append(ciphers, fmt.Sprintf("%d", c))
        }
    }

    // Extensions — from SupportedProtos, SignatureSchemes etc.
    // In practice, ClientHelloInfo doesn't expose raw extension IDs.
    // We use available fields as proxy.
    var extensions []string
    if len(hello.SupportedProtos) > 0 {
        extensions = append(extensions, "16") // ALPN extension ID
    }
    if len(hello.SupportedVersions) > 0 {
        extensions = append(extensions, "43") // supported_versions
    }
    if hello.ServerName != "" {
        extensions = append(extensions, "0") // SNI
    }
    // SignatureSchemes → signature_algorithms extension (13)
    if len(hello.SignatureSchemes) > 0 {
        extensions = append(extensions, "13")
    }
    // SupportedCurves → supported_groups (10)
    if len(hello.SupportedCurves) > 0 {
        extensions = append(extensions, "10")
    }

    // Elliptic curves
    var curves []string
    for _, c := range hello.SupportedCurves {
        if !isGREASECurve(uint16(c)) {
            curves = append(curves, fmt.Sprintf("%d", c))
        }
    }

    // Point formats (not directly exposed — use "0" as default)
    pointFormats := "0"

    // Construct JA3 string: SSLVersion,Ciphers,Extensions,EllipticCurves,EllipticCurvePointFormats
    ja3Str := fmt.Sprintf("%d,%s,%s,%s,%s",
        version,
        strings.Join(ciphers, "-"),
        strings.Join(extensions, "-"),
        strings.Join(curves, "-"),
        pointFormats,
    )

    hash := md5.Sum([]byte(ja3Str))
    return fmt.Sprintf("%x", hash)
}

func isGREASE(v uint16) bool {
    // GREASE values: 0x0a0a, 0x1a1a, 0x2a2a, ... 0xfafa
    return v&0x0f0f == 0x0a0a
}

func isGREASECurve(v uint16) bool {
    return v&0x0f0f == 0x0a0a
}

// Known JA3 hashes for common clients
var knownBotJA3 = map[string]Signal{
    "7dc465e28f1d46d873bf9c2ae6d05e43": {Name: "ja3_curl", Weight: 25, Category: CategoryNetwork, Reason: "curl TLS fingerprint"},
    "6597c339e4df3929c3cff7cc61c2b1c9": {Name: "ja3_python_requests", Weight: 25, Category: CategoryNetwork, Reason: "Python requests TLS fingerprint"},
    "0700c0e4ebf0d82de6a30e0c28d21c0e": {Name: "ja3_scrapy", Weight: 30, Category: CategoryNetwork, Reason: "Scrapy TLS fingerprint"},
}

var knownGoodJA3 = map[string]bool{
    "b32309a26951912be7dba376398abc3b": true, // Chrome 120
    "a0e9f5d64349fb13191bc781f81f42e1": true, // Firefox 120
}

// JA3Detector hooks into TLS handshake to compute fingerprint per connection
type JA3Detector struct {
    scorer *BotScorer
}

// GetConfigForClient returns a per-client TLS config that captures the ClientHello
func (d *JA3Detector) GetConfigForClient(hello *tls.ClientHelloInfo) (*tls.Config, error) {
    hash := JA3Fingerprint(hello)

    // Look up the session for this connection's remote address
    ip, _, _ := net.SplitHostPort(hello.Conn.RemoteAddr().String())
    session := d.scorer.sessionFor(ip)

    if sig, bad := knownBotJA3[hash]; bad {
        session.AddSignal(sig)
    } else if !knownGoodJA3[hash] {
        // Unknown fingerprint — add modest signal
        session.AddSignal(Signal{"ja3_unknown", 10, CategoryNetwork,
            fmt.Sprintf("Unknown JA3: %s", hash)})
    }

    // Store JA3 in session for debugging
    session.AddSignal(Signal{
        Name:     fmt.Sprintf("ja3:%s", hash),
        Weight:   0, // informational only
        Category: CategoryNetwork,
        Reason:   "JA3 fingerprint recorded",
    })

    return nil, nil // use default TLS config
}
```

---

## Day 102 — Session Behavioral Analysis

| | |
|---|---|
| **Learning objective** | Track session-level behavioral patterns that distinguish bots from humans |
| **Key packages** | `sync`, `time`, `net/http` |
| **Security context** | Bots interact with APIs programmatically — no navigation patterns, no referrer chain, too fast |

### Key Concepts

**Human vs bot behavior:**

| Signal | Human | Bot |
|--------|-------|-----|
| Requests/second | 0.1–2 | 10–1000 |
| Unique paths/session | 5–20 | 1–3 |
| Referer presence | 60–80% | < 10% |
| Session depth | Grows gradually | Constant single-endpoint |
| Inter-request timing | Variable (50ms–5s) | Fixed intervals |

### Exercise — Session Behavioral Tracker

```go
package bot

import (
    "net/http"
    "sync"
    "time"
)

type sessionData struct {
    mu              sync.Mutex
    requestTimes    []time.Time // sliding 60s window
    uniquePaths     map[string]struct{}
    refererCount    int
    requestCount    int
    created         time.Time
    lastSee         time.Time
}

func newSessionData() *sessionData {
    return &sessionData{
        uniquePaths: make(map[string]struct{}),
        created:     time.Now(),
        lastSee:     time.Now(),
    }
}

func (sd *sessionData) record(r *http.Request) {
    sd.mu.Lock()
    defer sd.mu.Unlock()

    now := time.Now()
    sd.requestCount++
    sd.lastSee = now
    sd.uniquePaths[r.URL.Path] = struct{}{}
    if r.Referer() != "" {
        sd.refererCount++
    }

    // Maintain sliding 60s window for velocity
    cutoff := now.Add(-60 * time.Second)
    valid := sd.requestTimes[:0]
    for _, t := range sd.requestTimes {
        if t.After(cutoff) {
            valid = append(valid, t)
        }
    }
    sd.requestTimes = append(valid, now)
}

// velocity returns requests per second over the last 60s
func (sd *sessionData) velocity() float64 {
    sd.mu.Lock()
    defer sd.mu.Unlock()
    if len(sd.requestTimes) == 0 {
        return 0
    }
    return float64(len(sd.requestTimes)) / 60.0
}

func (sd *sessionData) pathDiversity() int {
    sd.mu.Lock()
    defer sd.mu.Unlock()
    return len(sd.uniquePaths)
}

func (sd *sessionData) refererRate() float64 {
    sd.mu.Lock()
    defer sd.mu.Unlock()
    if sd.requestCount == 0 {
        return 0
    }
    return float64(sd.refererCount) / float64(sd.requestCount)
}

// BehavioralDetector implements Detector
type BehavioralDetector struct {
    mu       sync.RWMutex
    sessions map[string]*sessionData // keyed by IP (or session cookie)
}

func NewBehavioralDetector() *BehavioralDetector {
    return &BehavioralDetector{sessions: make(map[string]*sessionData)}
}

func (d *BehavioralDetector) getSession(ip string) *sessionData {
    d.mu.RLock()
    s, ok := d.sessions[ip]
    d.mu.RUnlock()
    if ok {
        return s
    }
    d.mu.Lock()
    defer d.mu.Unlock()
    s, ok = d.sessions[ip]
    if !ok {
        s = newSessionData()
        d.sessions[ip] = s
    }
    return s
}

func (d *BehavioralDetector) Detect(r *http.Request, score *SessionScore) {
    ip := clientIP(r)
    session := d.getSession(ip)
    session.record(r)

    // Only evaluate after enough data
    if session.requestCount < 5 {
        return
    }

    velocity := session.velocity()
    diversity := session.pathDiversity()
    refRate := session.refererRate()

    // High velocity signal
    if velocity > 10 {
        score.AddSignal(Signal{
            "high_velocity", 20, CategoryBehavior,
            fmt.Sprintf("%.1f req/s over last 60s", velocity),
        })
    }

    // Low path diversity
    if diversity < 2 && session.requestCount > 10 {
        score.AddSignal(Signal{
            "low_path_diversity", 10, CategoryBehavior,
            fmt.Sprintf("only %d unique paths in %d requests", diversity, session.requestCount),
        })
    }

    // Missing referer
    if refRate < 0.2 && session.requestCount > 10 {
        score.AddSignal(Signal{
            "low_referer_rate", 10, CategoryBehavior,
            fmt.Sprintf("referer present in only %.0f%% of requests", refRate*100),
        })
    }
}
```

---

## Day 103 — IP Reputation & ASN Classification

| | |
|---|---|
| **Learning objective** | Classify IP addresses as datacenter, residential, or known malicious |
| **Key packages** | `net`, `strings` |
| **Security context** | 60–70% of bot traffic originates from datacenter IPs; classification gives strong first-pass signal |

### Key Concepts

**ASN classification:** Each ASN (Autonomous System Number) maps to an organization. AWS = ASN 16509, GCP = ASN 15169, etc. IPs in datacenter ASNs are more likely bots.

**In production:** Use MaxMind GeoIP2 or ipinfo.io. For this exercise: hardcode representative CIDR ranges for major cloud providers.

### Exercise — IP Classifier

```go
package bot

import (
    "net"
)

type IPCategory string

const (
    IPCategoryDatacenter  IPCategory = "datacenter"
    IPCategoryHosting     IPCategory = "hosting"
    IPCategoryResidential IPCategory = "residential"
    IPCategoryUnknown     IPCategory = "unknown"
)

type cidrEntry struct {
    network  *net.IPNet
    provider string
    category IPCategory
}

// Representative CIDR ranges for major cloud providers
// In production: use a regularly updated database
var cloudCIDRs = []struct {
    cidr     string
    provider string
}{
    // AWS
    {"3.0.0.0/8", "AWS"},
    {"13.32.0.0/11", "AWS CloudFront"},
    {"52.0.0.0/8", "AWS"},
    {"54.0.0.0/8", "AWS"},
    // GCP
    {"8.34.208.0/20", "GCP"},
    {"34.0.0.0/9", "GCP"},
    {"35.184.0.0/13", "GCP"},
    // Azure
    {"13.64.0.0/11", "Azure"},
    {"20.0.0.0/8", "Azure"},
    {"40.64.0.0/10", "Azure"},
    // DigitalOcean
    {"104.16.0.0/12", "DigitalOcean"},
    {"159.89.0.0/16", "DigitalOcean"},
    // Linode/Akamai
    {"45.33.0.0/17", "Linode"},
    {"139.162.0.0/16", "Linode"},
    // Hetzner
    {"5.9.0.0/16", "Hetzner"},
    {"78.46.0.0/15", "Hetzner"},
}

var compiledCIDRs []cidrEntry

func init() {
    for _, c := range cloudCIDRs {
        _, network, err := net.ParseCIDR(c.cidr)
        if err != nil {
            continue
        }
        compiledCIDRs = append(compiledCIDRs, cidrEntry{
            network:  network,
            provider: c.provider,
            category: IPCategoryDatacenter,
        })
    }
}

func classifyIP(ipStr string) (IPCategory, string) {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return IPCategoryUnknown, ""
    }

    for _, entry := range compiledCIDRs {
        if entry.network.Contains(ip) {
            return entry.category, entry.provider
        }
    }

    // IPv6-only from unusual ASN — modest signal
    if ip.To4() == nil {
        return IPCategoryUnknown, "ipv6"
    }

    return IPCategoryResidential, ""
}

// IPReputationDetector implements Detector
type IPReputationDetector struct{}

func (d *IPReputationDetector) Detect(r *http.Request, session *SessionScore) {
    ip := clientIP(r)
    category, provider := classifyIP(ip)

    switch category {
    case IPCategoryDatacenter:
        session.AddSignal(Signal{
            "datacenter_ip", 15, CategoryNetwork,
            fmt.Sprintf("IP belongs to %s (datacenter)", provider),
        })
    case IPCategoryHosting:
        session.AddSignal(Signal{
            "hosting_ip", 20, CategoryNetwork,
            fmt.Sprintf("IP belongs to %s (hosting)", provider),
        })
    }
}
```

---

## Day 104 — Proof-of-Work & Invisible Challenges

| | |
|---|---|
| **Learning objective** | Issue and verify SHA256 PoW challenges; gate session cookies on completion |
| **Key packages** | `crypto/sha256`, `crypto/rand`, `crypto/hmac`, `net/http`, `encoding/hex` |
| **Security context** | PoW imposes computational cost on bots while allowing human users through (they solve it in JS) |

### Key Concepts

**PoW mechanics:**
- Server issues `{challenge: random_hex, difficulty: 16, expires: unix_timestamp}`
- Client finds `nonce` such that `SHA256(challenge + nonce)` has ≥ 16 leading zero bits
- Expected attempts: 2^16 = 65,536 (takes ~65ms on modern CPU)
- Difficulty=20 → ~1 million attempts → 1 second — too slow for human-facing challenges

**Session cookie issuance:** After PoW solved, issue signed cookie. Subsequent requests verify cookie rather than re-challenging.

### Exercise — PoW Challenge Endpoint + Session Cookie

```go
package bot

import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type PowSystem struct {
    secret     []byte
    difficulty int
    cookieName string
    expiry     time.Duration
}

func NewPowSystem(secret []byte, difficulty int) *PowSystem {
    return &PowSystem{
        secret:     secret,
        difficulty: difficulty,
        cookieName: "_bot_challenge",
        expiry:     30 * time.Minute,
    }
}

// ChallengeHandler issues a PoW challenge
func (ps *PowSystem) ChallengeHandler(w http.ResponseWriter, r *http.Request) {
    raw := make([]byte, 16)
    rand.Read(raw)
    expires := time.Now().Add(5 * time.Minute).Unix()

    // HMAC-sign the challenge + expiry
    mac := hmac.New(sha256.New, ps.secret)
    mac.Write(raw)
    mac.Write([]byte(strconv.FormatInt(expires, 10)))
    sig := hex.EncodeToString(mac.Sum(nil)[:8])

    challengeHex := hex.EncodeToString(raw)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "challenge":  challengeHex + "." + sig, // embed sig for verification
        "difficulty": ps.difficulty,
        "expires":    expires,
    })
}

// VerifyHandler validates a solved PoW and issues a session cookie
func (ps *PowSystem) VerifyHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Challenge string `json:"challenge"`
        Nonce     string `json:"nonce"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid_request"}`, http.StatusBadRequest)
        return
    }

    parts := strings.SplitN(req.Challenge, ".", 2)
    if len(parts) != 2 {
        http.Error(w, `{"error":"invalid_challenge"}`, http.StatusBadRequest)
        return
    }
    challengeHex, storedSig := parts[0], parts[1]

    challengeBytes, err := hex.DecodeString(challengeHex)
    if err != nil {
        http.Error(w, `{"error":"invalid_challenge"}`, http.StatusBadRequest)
        return
    }

    // Verify HMAC (simplified — in production include expiry in MAC)
    mac := hmac.New(sha256.New, ps.secret)
    mac.Write(challengeBytes)
    expectedSig := hex.EncodeToString(mac.Sum(nil)[:8])
    if !hmac.Equal([]byte(storedSig), []byte(expectedSig)) {
        http.Error(w, `{"error":"invalid_challenge"}`, http.StatusBadRequest)
        return
    }

    // Verify nonce
    nonce, err := hex.DecodeString(req.Nonce)
    if err != nil {
        http.Error(w, `{"error":"invalid_nonce"}`, http.StatusBadRequest)
        return
    }

    h := sha256.New()
    h.Write(challengeBytes)
    h.Write(nonce)
    hash := h.Sum(nil)
    if leadingZeroBits(hash) < ps.difficulty {
        http.Error(w, `{"error":"insufficient_work"}`, http.StatusBadRequest)
        return
    }

    // Issue session cookie
    sessionToken := ps.issueSessionToken()
    http.SetCookie(w, &http.Cookie{
        Name:     ps.cookieName,
        Value:    sessionToken,
        Path:     "/",
        Expires:  time.Now().Add(ps.expiry),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "verified"})
}

func (ps *PowSystem) issueSessionToken() string {
    raw := make([]byte, 16)
    rand.Read(raw)
    mac := hmac.New(sha256.New, ps.secret)
    mac.Write(raw)
    mac.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
    return hex.EncodeToString(raw) + "." + hex.EncodeToString(mac.Sum(nil)[:8])
}

// HasValidSession checks if request has a valid PoW-issued session cookie
func (ps *PowSystem) HasValidSession(r *http.Request) bool {
    cookie, err := r.Cookie(ps.cookieName)
    if err != nil {
        return false
    }
    parts := strings.SplitN(cookie.Value, ".", 2)
    if len(parts) != 2 {
        return false
    }
    // In production: verify HMAC and check expiry
    // For brevity: just check format
    _, err1 := hex.DecodeString(parts[0])
    _, err2 := hex.DecodeString(parts[1])
    return err1 == nil && err2 == nil
}
```

---

## Day 105 — Bot Mitigation Strategies

| | |
|---|---|
| **Learning objective** | Implement tiered bot mitigation: slow → challenge → block → honeypot → shadow-ban |
| **Key packages** | `net/http`, `time`, `sync`, `encoding/json` |
| **Security context** | Hard blocking causes false positives; tiered mitigation reduces impact on legitimate users while imposing maximum cost on bots |

### Key Concepts

**Mitigation tiers:**

| Score | Response | Bot effect | Legitimate user effect |
|-------|----------|------------|----------------------|
| 20–40 | +500ms delay | Reduces throughput | Barely noticeable |
| 40–70 | PoW challenge | Requires 65k SHA256 ops | Solved by browser JS in ~100ms |
| > 70 | 403 block | Request rejected | False positive risk |
| Honeypot | Permanent flag | IP permanently flagged | Users never find hidden endpoint |
| Shadow-ban | Fake 200 with empty data | Bot thinks it's working | Never applied to real users |

### Exercise — Tiered Bot Mitigation Middleware

```go
package bot

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

// HoneypotTracker permanently flags any IP that hits a honeypot endpoint
type HoneypotTracker struct {
    mu      sync.RWMutex
    flagged map[string]time.Time // IP → when flagged
}

func NewHoneypotTracker() *HoneypotTracker {
    return &HoneypotTracker{flagged: make(map[string]time.Time)}
}

func (ht *HoneypotTracker) Flag(ip string) {
    ht.mu.Lock()
    ht.flagged[ip] = time.Now()
    ht.mu.Unlock()
}

func (ht *HoneypotTracker) IsFlagged(ip string) bool {
    ht.mu.RLock()
    _, ok := ht.flagged[ip]
    ht.mu.RUnlock()
    return ok
}

// HoneypotHandler — any IP that calls this gets permanently flagged
// The path should be hidden (never linked from real pages, never in sitemap)
func (ht *HoneypotTracker) HoneypotHandler(scorer *BotScorer) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := clientIP(r)
        ht.Flag(ip)

        // Add maximum score signal to BotScorer
        session := scorer.sessionFor(ip)
        session.AddSignal(SigHoneypotHit)

        // Return a plausible 200 — don't reveal that this is a honeypot
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "data":  []interface{}{},
            "total": 0,
            "page":  1,
        })
    })
}

// ShadowBanHandler returns plausible but empty data for confirmed bots
func ShadowBanHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Return 200 with empty/fake data
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "results": []interface{}{},
            "count":   0,
        })
        // Do NOT call next — backend never sees the request
    })
}

// TieredBotMiddleware applies mitigation based on bot score
type TieredBotMiddleware struct {
    scorer   *BotScorer
    pow      *PowSystem
    honeypot *HoneypotTracker

    // Thresholds
    delayScore    float64 // score >= this → add artificial delay
    challengeScore float64 // score >= this → PoW challenge
    blockScore    float64 // score >= this → 403
    shadowScore   float64 // score >= this (confirmed) → shadow ban

    delayAmount time.Duration
}

func NewTieredBotMiddleware(scorer *BotScorer, pow *PowSystem, honeypot *HoneypotTracker) *TieredBotMiddleware {
    return &TieredBotMiddleware{
        scorer:         scorer,
        pow:            pow,
        honeypot:       honeypot,
        delayScore:     20,
        challengeScore: 40,
        blockScore:     70,
        shadowScore:    85,
        delayAmount:    500 * time.Millisecond,
    }
}

func (m *TieredBotMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := clientIP(r)

        // Honeypot flagged → shadow ban (permanent)
        if m.honeypot.IsFlagged(ip) {
            ShadowBanHandler(next).ServeHTTP(w, r)
            return
        }

        // Skip challenge if they have a valid PoW session
        if m.pow.HasValidSession(r) {
            next.ServeHTTP(w, r)
            return
        }

        score, _ := m.scorer.Score(r)

        switch {
        case score >= m.shadowScore:
            // Confirmed bot — shadow ban
            ShadowBanHandler(next).ServeHTTP(w, r)

        case score >= m.blockScore:
            // High confidence bot — hard block
            http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)

        case score >= m.challengeScore:
            // Medium confidence — serve PoW challenge
            w.Header().Set("X-Bot-Challenge", "required")
            http.Redirect(w, r, "/challenge", http.StatusFound)

        case score >= m.delayScore:
            // Low confidence — slow down
            time.Sleep(m.delayAmount)
            next.ServeHTTP(w, r)

        default:
            next.ServeHTTP(w, r)
        }
    })
}
```

---

## Week 14 Mini Project: Bot Detection Engine

**Build a complete bot detection middleware** that integrates all week's signals.

### Architecture

```
Request
    │
    ▼
[JA3Detector — hooks into TLS handshake via GetConfigForClient]
    │
    ▼
[HeaderConsistencyDetector — runs on each HTTP request]
    │
    ▼
[IPReputationDetector — CIDR classification]
    │
    ▼
[BehavioralDetector — session velocity, path diversity, referrer]
    │
    ▼
[BotScorer — aggregates signals → score]
    │
    ▼
[TieredBotMiddleware]
    ├── score 0–20  → allow
    ├── score 20–40 → 500ms delay + allow
    ├── score 40–70 → PoW challenge redirect
    ├── score 70–85 → 403 block
    └── score 85+   → shadow ban (200, empty data)
    │
    ▼
[HoneypotHandler → /api/hidden]
    │ any hit → +100 score, permanent flag
    ▼
[/debug/bots — stats dashboard]
```

### `/debug/bots` Response Format

```json
{
  "top_scored_ips": [
    {
      "ip": "1.2.3.4",
      "score": 85.0,
      "action": "shadow_ban",
      "signals": [
        {"name": "datacenter_ip", "weight": 15, "reason": "AWS"},
        {"name": "high_velocity", "weight": 20, "reason": "15.2 req/s"},
        {"name": "low_path_diversity", "weight": 10, "reason": "1 unique path in 50 requests"},
        {"name": "ja3_curl", "weight": 25, "reason": "curl TLS fingerprint"},
        {"name": "missing_accept_language", "weight": 10, "reason": "..."}
      ]
    }
  ],
  "honeypot_hits": 3,
  "flagged_ips": 3,
  "challenges_issued": 12,
  "sessions_tracked": 891,
  "action_breakdown": {
    "allow": 8234,
    "delay": 156,
    "challenge": 12,
    "block": 8,
    "shadow_ban": 3
  }
}
```

### Feature Checklist

- [ ] Header consistency fingerprinting (UA vs header set)
- [ ] JA3-like TLS fingerprint via `tls.ClientHelloInfo`
- [ ] Session behavioral scoring (velocity, path diversity, referrer rate)
- [ ] IP ASN classification (datacenter CIDR list for major cloud providers)
- [ ] PoW challenge issuance (`GET /challenge`) and verification (`POST /verify`)
- [ ] Session cookie issued after PoW solved (bypasses future challenges)
- [ ] Honeypot endpoint (`/api/hidden`) — permanent IP flag
- [ ] Shadow-ban handler (200 with empty data)
- [ ] Tiered mitigation: delay → challenge → block → shadow-ban
- [ ] `/debug/bots` dashboard
- [ ] Integration with `ThreatIntel` store from Phase 3

### Engineer Takeaway

DataDome and PerimeterX sell bot detection for $50k–$500k/year. You now understand their core algorithms: header fingerprinting, JA3, behavioral scoring, and tiered mitigation. You can integrate and tune these products with confidence, adjust false-positive thresholds, and judge when in-house detection is enough (internal APIs) and when you need a vendor's breadth (global IP reputation, browser JS challenges).

---

## Week 14 Review

| Day | Skill | Real-World Use |
|-----|-------|----------------|
| 99 | Bot taxonomy + signal scoring framework | Architect pluggable bot detection |
| 100 | HTTP header consistency fingerprinting | Detect bots claiming to be Chrome |
| 101 | JA3 TLS fingerprinting | Identify clients by their TLS stack |
| 102 | Session behavioral analysis | Catch bots that fake good headers |
| 103 | IP ASN/datacenter classification | First-pass bot signal from IP origin |
| 104 | PoW challenge system + session cookies | Impose CPU cost on bots |
| 105 | Tiered mitigation: delay/challenge/block/shadow-ban | Minimize false positives while blocking bots |
