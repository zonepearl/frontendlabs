# Week 10: Web Application Firewall (Days 68–74)

**Phase:** 3 — L7 Networking + WAF + API Security
**Goal:** Build a production-grade WAF from scratch: OWASP rule set, multi-layer normalization, body inspection, scoring model, and hot-reload config. By end of week you understand every feature AWS WAF charges $5/month/rule for.

---

> **Background in this wiki:** HTTPS guide [Ch 12 (web attacks and defences)](../../v2-https/real-life-guide-v1.md#chapter-12-security-protecting-the-entire-lifecycle), [Ch 20 (request smuggling, reproduced)](../../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling), [Ch 22 (SSRF)](../../v2-https/real-life-guide-v1.md#chapter-22-outbound-requests-ssrf-and-safe-http-clients); TCP/IP guide [Ch 53 (firewalls, and policy as tests)](../../networking/tcp-ip/real-life-guide-v1.md#chapter-53-firewalls-in-the-real-world-host-firewalls-cloud-security-groups-nacls). Security in Depth: [Part 6 (advanced web attacks)](../../security/real-life-security-guide-v1.md#part-6-advanced-web-and-api-exploitation) and [Chapter 40C (WAF architecture and bypasses)](../../security/real-life-security-guide-v1.md#chapter-40c-waf-architecture-rules-and-bypasses).

## Day 68 — WAF Fundamentals & OWASP Top 10

### Learning Objectives
- Map each OWASP Top 10 category to a WAF rule
- Build the WAF skeleton: `Rule` interface, scoring middleware
- Understand detection mode vs blocking mode

### OWASP Top 10 → WAF Rule Mapping

| OWASP Category | WAF Defense | Your Day |
|---------------|-------------|----------|
| A01 Broken Access Control | Auth middleware, CORS | Day 67, 75 |
| A02 Cryptographic Failures | TLS enforcement, HSTS | Day 39, 67 |
| A03 Injection (SQLi, XSS, Cmd) | WAF rules Day 69-71 | Days 69-71 |
| A04 Insecure Design | Schema validation | Day 78 |
| A05 Security Misconfiguration | Security headers | Day 67 |
| A06 Vulnerable Components | N/A (code audit) | — |
| A07 Auth Failures | Rate limiting, MFA | Day 76, 80 |
| A08 Data Integrity Failures | HMAC signing, ECDSA | Day 42, 44 |
| A09 Logging Failures | Structured audit log | Day 18, 84 |
| A10 SSRF | URL allowlisting | Day 71 |

### Exercise: WAF Skeleton

```go
package waf

import (
    "fmt"
    "net/http"
    "sync"
    "sync/atomic"
)

// Rule defines a WAF detection rule
type Rule interface {
    ID() string                           // unique identifier, e.g., "SQLI-001"
    Category() string                     // "SQLI", "XSS", "TRAVERSAL", "CMDINJ"
    Match(r *http.Request, body []byte) (matched bool, score int, detail string)
}

// MatchResult records one rule match for audit logging
type MatchResult struct {
    RuleID    string
    Category  string
    Score     int
    Detail    string
    Location  string // "query", "header", "body"
    Raw       string // matched value (truncated)
}

// WAFConfig is hot-reloadable at runtime
type WAFConfig struct {
    Mode           string         // "off", "detection", "blocking"
    BlockThreshold int            // block if total score >= this
    CategoryThresh map[string]int // per-category thresholds
    TrustedCIDRs   []string       // bypass WAF completely
}

// WAFMiddleware is the main WAF entry point
type WAFMiddleware struct {
    rules    []Rule
    config   atomic.Value // *WAFConfig
    auditor  AuditWriter

    // Metrics
    requestsChecked atomic.Int64
    requestsBlocked atomic.Int64
    ruleHits        sync.Map // ruleID → count
}

func NewWAFMiddleware(rules []Rule, cfg *WAFConfig, auditor AuditWriter) *WAFMiddleware {
    m := &WAFMiddleware{rules: rules, auditor: auditor}
    m.config.Store(cfg)
    return m
}

func (waf *WAFMiddleware) Reload(cfg *WAFConfig) {
    waf.config.Store(cfg)
    fmt.Println("[WAF] config reloaded")
}

func (waf *WAFMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg := waf.config.Load().(*WAFConfig)

        // Check if this IP is in the trusted CIDR bypass list
        if waf.isTrusted(r, cfg) {
            next.ServeHTTP(w, r)
            return
        }

        // Read and buffer body for inspection (with size limit)
        body, err := readBody(r)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusRequestEntityTooLarge)
            fmt.Fprintf(w, `{"error":"request body too large"}`)
            return
        }
        restoreBody(r, body) // restore for downstream handler

        // Run all rules, collect scores
        var results []MatchResult
        totalScore := 0
        categoryScores := make(map[string]int)

        for _, rule := range waf.rules {
            matched, score, detail := rule.Match(r, body)
            if matched {
                loc, raw := findMatchLocation(r, body, detail)
                result := MatchResult{
                    RuleID:   rule.ID(),
                    Category: rule.Category(),
                    Score:    score,
                    Detail:   detail,
                    Location: loc,
                    Raw:      truncate(raw, 100),
                }
                results = append(results, result)
                totalScore += score
                categoryScores[rule.Category()] += score

                // Track rule hit metrics
                waf.ruleHits.Store(rule.ID(), 0)
            }
        }

        waf.requestsChecked.Add(1)

        // Determine action
        action := "allow"
        if len(results) > 0 {
            if waf.shouldBlock(cfg, totalScore, categoryScores) {
                action = "block"
            } else {
                action = "detect"
            }
        }

        // Audit log every match (regardless of action)
        if len(results) > 0 {
            waf.auditor.Log(r, results, action, totalScore)
        }

        if action == "block" && cfg.Mode == "blocking" {
            waf.requestsBlocked.Add(1)
            requestID, _ := r.Context().Value(requestIDKey).(string)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusForbidden)
            // SECURITY: never reveal which rule triggered — information disclosure
            fmt.Fprintf(w, `{"error":"forbidden","request_id":%q}`, requestID)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (waf *WAFMiddleware) shouldBlock(cfg *WAFConfig, total int, catScores map[string]int) bool {
    if total >= cfg.BlockThreshold {
        return true
    }
    for cat, score := range catScores {
        if thresh, ok := cfg.CategoryThresh[cat]; ok && score >= thresh {
            return true
        }
    }
    return false
}
```

### Real-World Context
AWS WAF, Cloudflare WAF, and ModSecurity all use this scoring model. `detection` mode (log but don't block) is how you tune a new WAF deployment — run for 2 weeks in detection, review false positives, then enable blocking. Requiring WAF mode to be in config means a SIGHUP can enable/disable blocking without a deployment.

---

## Day 69 — SQL Injection Detection

### Learning Objectives
- Implement SQLi detection with multiple techniques: keywords, quote balance, comment patterns
- Tune for zero false positives on 15 legitimate inputs
- Understand why WAF-only SQLi defense is insufficient (defense-in-depth)

### Key Concepts

```
SQLi Attack Categories:
  Classic:    ' OR '1'='1   (breaks out of string context)
  Union:      ' UNION SELECT username,password FROM users--
  Blind:      ' AND 1=1--  / ' AND 1=2-- (boolean-based)
  Time-based: ' AND SLEEP(5)-- (detects via response time)
  Stacked:    '; DROP TABLE users--

Common Evasion Techniques (why normalization matters):
  Case:       SeLeCt, sElEcT
  Comments:   SE/*comment*/LECT, UN/**/ION
  URL encode: %53%45%4C%45%43%54
  Whitespace: SELECT\t\nFROM (tab/newline instead of space)
```

### Exercise: SQLi Detection Rule

```go
package waf

import (
    "net/http"
    "regexp"
    "strings"
    "unicode"
)

// Pre-compiled at package init — NEVER inside Match()
var (
    // Keywords that appear in SQLi but rarely in legitimate data
    sqliKeywordsRe = regexp.MustCompile(`(?i)\b(UNION\s+(?:ALL\s+)?SELECT|INSERT\s+INTO|UPDATE\s+\w+\s+SET|DELETE\s+FROM|DROP\s+(?:TABLE|DATABASE)|CREATE\s+(?:TABLE|DATABASE|USER)|EXEC(?:UTE)?|CAST\s*\(|CONVERT\s*\(|INFORMATION_SCHEMA|SLEEP\s*\(|BENCHMARK\s*\(|LOAD_FILE\s*\(|INTO\s+OUTFILE)\b`)

    // Comment patterns used in SQLi
    sqliCommentRe = regexp.MustCompile(`(--|#|\/\*.*?\*\/)`)

    // Dangerous tautologies
    sqliTautologyRe = regexp.MustCompile(`(?i)'\s*(?:OR|AND)\s+['"]?[0-9a-z]+['"]?\s*=\s*['"]?[0-9a-z]+['"]?`)
)

type SQLiRule struct{}

func (r *SQLiRule) ID() string       { return "SQLI-001" }
func (r *SQLiRule) Category() string { return "SQLI" }

func (r *SQLiRule) Match(req *http.Request, body []byte) (bool, int, string) {
    inputs := extractInputs(req, body)
    for _, input := range inputs {
        normalized := normalizeSQLi(input)

        // Check 1: keywords (score: 5)
        if m := sqliKeywordsRe.FindString(normalized); m != "" {
            return true, 5, fmt.Sprintf("SQLi keyword: %q", m)
        }

        // Check 2: tautology patterns (score: 8)
        if m := sqliTautologyRe.FindString(normalized); m != "" {
            return true, 8, fmt.Sprintf("SQLi tautology: %q", m)
        }

        // Check 3: unbalanced quotes + SQL comment (score: 7)
        stripped := sqliCommentRe.ReplaceAllString(normalized, "")
        if unbalancedQuotes(stripped) && sqliCommentRe.MatchString(normalized) {
            return true, 7, "SQLi: unbalanced quotes with comment"
        }
    }
    return false, 0, ""
}

// normalizeSQLi prepares input for rule matching
func normalizeSQLi(input string) string {
    // URL decode
    decoded, _ := url.QueryUnescape(input)
    // Collapse whitespace (including tabs, newlines)
    decoded = strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) { return ' ' }
        return r
    }, decoded)
    // Collapse multiple spaces
    for strings.Contains(decoded, "  ") {
        decoded = strings.ReplaceAll(decoded, "  ", " ")
    }
    return strings.TrimSpace(decoded)
}

// unbalancedQuotes returns true if single quotes are not balanced
func unbalancedQuotes(s string) bool {
    count := 0
    for _, c := range s {
        if c == '\'' { count++ }
    }
    return count%2 != 0
}

// Test with these payloads (all should trigger):
var sqliTestPayloads = []string{
    "' OR '1'='1",
    "1' UNION SELECT * FROM users--",
    "admin'--",
    "1; DROP TABLE users--",
    "' OR 1=1#",
    "' UNION ALL SELECT NULL,NULL,NULL--",
    "1 AND SLEEP(5)--",
    "1' AND EXTRACTVALUE(1,CONCAT(0x7e,(SELECT version())))--",
}

// Legitimate inputs that must NOT trigger (tune to 0 false positives):
var sqliLegitimateInputs = []string{
    "O'Brien",                   // Irish name with apostrophe
    "It's a great product",      // natural language apostrophe
    "SELECT from dropdown",      // SELECT in natural language
    "DROP shipping address",     // DROP in context
    "2024-01-01",               // date
    "user@example.com",         // email
    "Product #123",             // hash in context
    "price > 100 AND category = 'electronics'", // URL query params
}
```

### Real-World Context
SQLi is still #3 in OWASP Top 10 despite 20 years of awareness. The 2017 Equifax breach was SQL injection. WAF as the last line of defense is correct — parameterized queries are the real fix, but WAF catches application bugs. The `normalizeSQLi` function is why ModSecurity has hundreds of normalization rules.

---

## Day 70 — XSS Detection

### Learning Objectives
- Detect reflected XSS, script injection, and event handler injection
- Normalize input before matching: HTML-unescape, URL-decode, unicode
- Build the PII scrubber integration (Day 79)

### Key Concepts

```
XSS Attack Vectors:
  <script>alert(1)</script>          — direct script injection
  <img src=x onerror=alert(1)>      — event handler injection
  javascript:alert(1)               — JavaScript URI scheme
  <svg onload=alert(1)>             — SVG attack
  <a href="javascript:void(0)">     — inline JS in href
  "><script>alert(1)</script>       — breaking out of attribute

Common Bypasses:
  <ScRiPt>alert(1)</sCrIpT>         — mixed case → normalize to lowercase
  %3Cscript%3E                      — URL encoding → URL decode first
  &lt;script&gt;                    — HTML entity encoding → HTML unescape first
  <scr\x00ipt>                      — null byte injection → remove null bytes
```

### Exercise: XSS Detection Rule

```go
package waf

import (
    "html"
    "net/http"
    "net/url"
    "regexp"
    "strings"
)

var (
    // Script tags in all forms
    xssScriptRe = regexp.MustCompile(`(?i)<\s*script[\s>]`)

    // JavaScript URI scheme
    xssJSURIRe = regexp.MustCompile(`(?i)javascript\s*:`)

    // Inline event handlers: onerror=, onload=, onclick=, etc.
    xssEventHandlerRe = regexp.MustCompile(`(?i)\bon[a-z]{2,20}\s*=`)

    // Dangerous tags
    xssDangerousTagsRe = regexp.MustCompile(`(?i)<\s*(iframe|frame|object|embed|applet|meta|link|base|form|input|button|select|textarea)\b`)

    // Data URI with script
    xssDataURIRe = regexp.MustCompile(`(?i)data:\s*text/html`)
)

type XSSRule struct{}

func (r *XSSRule) ID() string       { return "XSS-001" }
func (r *XSSRule) Category() string { return "XSS" }

func (r *XSSRule) Match(req *http.Request, body []byte) (bool, int, string) {
    inputs := extractInputs(req, body)
    for _, input := range inputs {
        normalized := normalizeXSS(input)

        if xssScriptRe.MatchString(normalized) {
            return true, 8, "XSS: script tag"
        }
        if xssJSURIRe.MatchString(normalized) {
            return true, 7, "XSS: javascript: URI"
        }
        if xssEventHandlerRe.MatchString(normalized) {
            return true, 6, "XSS: event handler injection"
        }
        if xssDangerousTagsRe.MatchString(normalized) {
            return true, 4, "XSS: dangerous HTML tag"
        }
        if xssDataURIRe.MatchString(normalized) {
            return true, 7, "XSS: data: URI with HTML"
        }
    }
    return false, 0, ""
}

// normalizeXSS applies multi-layer decoding before matching
// Order matters: URL decode → HTML unescape → lowercase → strip null bytes
func normalizeXSS(input string) string {
    // URL decode (handles %3C → <)
    decoded, err := url.QueryUnescape(input)
    if err != nil {
        decoded = input
    }

    // HTML unescape (handles &lt; → <)
    decoded = html.UnescapeString(decoded)

    // Double URL decode (handles %253C → %3C → <)
    if again, err := url.QueryUnescape(decoded); err == nil {
        decoded = again
    }

    // Remove null bytes (used to break pattern matching)
    decoded = strings.ReplaceAll(decoded, "\x00", "")

    // Lowercase for case-insensitive matching
    return strings.ToLower(decoded)
}

// Test payloads — all should trigger XSS detection:
var xssTestPayloads = []string{
    "<script>alert(1)</script>",
    "<ScRiPt>alert(1)</ScRiPt>",
    "%3Cscript%3Ealert(1)%3C/script%3E",
    "&lt;script&gt;alert(1)&lt;/script&gt;",
    "<img src=x onerror=alert(1)>",
    "javascript:alert(1)",
    "<svg onload=alert(1)>",
    "<iframe src=javascript:alert(1)>",
}

// Legitimate inputs that must NOT trigger:
var xssLegitInputs = []string{
    "Hello <World>",                           // angle bracket in text (no tag)
    "2 < 3 is true",                           // math comparison
    "onclick button behavior documentation",   // 'onclick' as a word
    "<b>bold</b> <i>italic</i>",               // common HTML in CMS
    "script writing techniques",               // 'script' as a word
}
```

### Real-World Context
XSS (OWASP A03) is the most common web vulnerability. The multi-layer normalization is why ModSecurity has 8 separate normalization rules that run before each match. An attacker who knows you check for `<script>` will try `%3Cscript%3E`. An attacker who knows you URL-decode will try `%253Cscript%253E` (double-encoded).

---

## Day 71 — Path Traversal & Command Injection

### Learning Objectives
- Detect path traversal in all encoded forms: raw, URL-encoded, double-encoded, unicode
- Detect command injection operators
- Understand why SSRF is a path traversal variant at the URL level

### Exercise: Path Traversal + Command Injection Rules

```go
package waf

import (
    "net/http"
    "net/url"
    "regexp"
    "strings"
)

var (
    // Path traversal in various encodings
    pathTraversalPatterns = []struct{ pattern, name string }{
        {`\.\./`, "dot-dot-slash"},
        {`\.\.\\`, "dot-dot-backslash"},
        {`%2e%2e%2f`, "url-encoded ../"},
        {`%2e%2e/`, "partial url-encoded ../"},
        {`\.%2e/`, "partial url-encoded ../"},
        {`%252e%252e%252f`, "double url-encoded ../"},
        {`%c0%ae%c0%ae/`, "unicode dot ../"},
        {`\x00`, "null byte"},
    }

    // Command injection operators
    cmdInjectionRe = regexp.MustCompile("(?i)(;|\\||&&|\\$\\(|`|\\$\\{|\\bexec\\b|\\bsystem\\b|\\bpassthru\\b|\\bshell_exec\\b|\\bpopen\\b|\\bproc_open\\b)")

    // SSRF patterns — internal IP access via URL parameters
    ssrfInternalIPRe = regexp.MustCompile(`(?:^|[^0-9])(?:127\.0\.0\.1|0\.0\.0\.0|localhost|169\.254\.\d+\.\d+|10\.\d+\.\d+\.\d+|172\.(?:1[6-9]|2\d|3[01])\.\d+\.\d+|192\.168\.\d+\.\d+)(?:[:/]|$)`)
)

type PathTraversalRule struct{}

func (r *PathTraversalRule) ID() string       { return "TRAV-001" }
func (r *PathTraversalRule) Category() string { return "TRAVERSAL" }

func (r *PathTraversalRule) Match(req *http.Request, body []byte) (bool, int, string) {
    // Check URL path + all query parameters + body
    inputs := []string{req.URL.Path, req.URL.RawQuery}
    inputs = append(inputs, extractInputs(req, body)...)

    for _, input := range inputs {
        normalized := normalizeTraversal(input)
        for _, p := range pathTraversalPatterns {
            re := regexp.MustCompile("(?i)" + p.pattern)
            if re.MatchString(normalized) {
                return true, 10, fmt.Sprintf("path traversal: %s in %q", p.name, input[:min(50, len(input))])
            }
        }
    }
    return false, 0, ""
}

func normalizeTraversal(input string) string {
    // Multiple decode passes to catch double-encoding
    result := input
    for i := 0; i < 3; i++ {
        decoded, err := url.PathUnescape(result)
        if err != nil || decoded == result { break }
        result = decoded
    }
    return strings.ToLower(result)
}

type CmdInjectionRule struct{}

func (r *CmdInjectionRule) ID() string       { return "CMDINJ-001" }
func (r *CmdInjectionRule) Category() string { return "CMDINJ" }

func (r *CmdInjectionRule) Match(req *http.Request, body []byte) (bool, int, string) {
    for _, input := range extractInputs(req, body) {
        if m := cmdInjectionRe.FindString(input); m != "" {
            return true, 9, fmt.Sprintf("command injection: %q", m)
        }
    }
    return false, 0, ""
}

type SSRFRule struct{}

func (r *SSRFRule) ID() string       { return "SSRF-001" }
func (r *SSRFRule) Category() string { return "SSRF" }

func (r *SSRFRule) Match(req *http.Request, body []byte) (bool, int, string) {
    // Check URL parameters that might contain target URLs
    for _, param := range []string{"url", "target", "redirect", "src", "href", "link"} {
        val := req.URL.Query().Get(param)
        if val == "" { continue }
        if ssrfInternalIPRe.MatchString(val) {
            return true, 9, fmt.Sprintf("SSRF: internal IP in %q param", param)
        }
    }
    return false, 0, ""
}
```

### Real-World Context
Path traversal (OWASP A01) was used in the 2021 Accellion FTA breach (100+ organizations). Command injection was the mechanism of the 2021 Log4Shell vulnerability (JNDI lookup injection). SSRF (OWASP A10) was used in the Capital One breach — attacker used metadata service (`169.254.169.254`) to steal AWS credentials.

---

## Day 72 — Request Body Inspection

### Learning Objectives
- Buffer and scan the request body without breaking the downstream handler
- Dispatch to specialized scanners by Content-Type
- Enforce body size limits to prevent resource exhaustion

### Exercise: Body Inspector

```go
package waf

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime"
    "net/http"
    "strings"
)

const maxBodySize = 1 * 1024 * 1024 // 1MB — scan limit

// readBody reads and buffers the request body, enforcing size limit.
// Returns (body, error). On success, body is also restored to r.Body.
func readBody(r *http.Request) ([]byte, error) {
    if r.Body == nil || r.Body == http.NoBody {
        return nil, nil
    }

    limited := io.LimitReader(r.Body, int64(maxBodySize)+1)
    body, err := io.ReadAll(limited)
    r.Body.Close()

    if int64(len(body)) > maxBodySize {
        return nil, fmt.Errorf("body exceeds %d byte scan limit", maxBodySize)
    }

    restoreBody(r, body)
    return body, err
}

func restoreBody(r *http.Request, body []byte) {
    r.Body = io.NopCloser(bytes.NewReader(body))
    r.ContentLength = int64(len(body))
}

// extractInputs returns all user-controlled values from the request
func extractInputs(r *http.Request, body []byte) []string {
    var inputs []string

    // URL query parameters
    for _, vals := range r.URL.Query() {
        inputs = append(inputs, vals...)
    }

    // Headers (scan non-standard headers — standard ones like User-Agent less risky)
    for header, vals := range r.Header {
        if isUserControlledHeader(header) {
            inputs = append(inputs, vals...)
        }
    }

    // Body — dispatch by Content-Type
    if len(body) > 0 {
        ct, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
        inputs = append(inputs, extractBodyInputs(ct, body)...)
    }

    return inputs
}

func isUserControlledHeader(header string) bool {
    // Headers that might carry user input
    userHeaders := map[string]bool{
        "X-Forwarded-For": true, "Referer": true, "User-Agent": true,
        "X-Real-IP": true, "Content-Disposition": true, "Cookie": true,
    }
    return userHeaders[http.CanonicalHeaderKey(header)]
}

// extractBodyInputs parses body values based on content type
func extractBodyInputs(contentType string, body []byte) []string {
    switch {
    case strings.HasPrefix(contentType, "application/json"):
        return extractJSONValues(body)
    case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
        values, _ := url.ParseQuery(string(body))
        var result []string
        for _, v := range values {
            result = append(result, v...)
        }
        return result
    case strings.HasPrefix(contentType, "text/"):
        return []string{string(body)}
    default:
        // For unknown types: scan as raw string
        return []string{string(body)}
    }
}

// extractJSONValues recursively extracts all string values from JSON
func extractJSONValues(data []byte) []string {
    var result []string

    var recurse func(v interface{})
    recurse = func(v interface{}) {
        switch val := v.(type) {
        case string:
            result = append(result, val)
        case map[string]interface{}:
            for _, v := range val {
                recurse(v)
            }
        case []interface{}:
            for _, v := range val {
                recurse(v)
            }
        }
    }

    var parsed interface{}
    if err := json.Unmarshal(data, &parsed); err == nil {
        recurse(parsed)
    }
    return result
}
```

### Real-World Context
Most WAFs only scan URL and headers by default — attackers know this and put payloads in the JSON body. Body scanning adds latency (you must buffer before forwarding). The 1MB limit prevents a 100MB JSON payload from stalling your scanner. The `Content-Type` dispatch is why WAFs have separate rules for `multipart/form-data`, JSON, and XML.

---

## Day 73 — WAF Normalization Pipeline

### Learning Objectives
- Build the multi-layer normalization pipeline that separates basic WAFs from production ones
- Test 10+ encoding bypass attempts against your normalized matcher
- Understand why ModSecurity has 8 normalization passes

### Exercise: Normalization Pipeline

```go
package waf

import (
    "html"
    "net/url"
    "regexp"
    "strings"
    "unicode"
)

var (
    sqlCommentRe    = regexp.MustCompile(`\/\*.*?\*\/`)
    multiSpaceRe    = regexp.MustCompile(`\s+`)
    nullByteRe      = regexp.MustCompile(`\x00`)
    htmlEntityRe    = regexp.MustCompile(`&#x?[0-9a-fA-F]+;?|&\w+;`)
)

// NormalizationPipeline applies layers of decoding before rule matching.
// Each layer decodes a different encoding that attackers might use.
func NormalizationPipeline(input string) string {
    result := input

    // Pass 1: Remove null bytes (used to truncate strings in C-based systems)
    result = nullByteRe.ReplaceAllString(result, "")

    // Pass 2: URL decode (handles %XX encoding)
    if decoded, err := url.QueryUnescape(result); err == nil {
        result = decoded
    }

    // Pass 3: HTML entity decode (handles &lt; &gt; &#60; &#x3c;)
    result = html.UnescapeString(result)

    // Pass 4: Second URL decode (catches double-encoded: %2520 → %20 → space)
    if decoded, err := url.QueryUnescape(result); err == nil {
        result = decoded
    }

    // Pass 5: Strip SQL comments (/**/ used to split keywords: UN/**/ION)
    result = sqlCommentRe.ReplaceAllString(result, " ")

    // Pass 6: Normalize whitespace (tab, newline used to split keywords)
    result = strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) { return ' ' }
        return r
    }, result)

    // Pass 7: Collapse multiple spaces
    result = multiSpaceRe.ReplaceAllString(result, " ")

    // Pass 8: Lowercase (for case-insensitive matching)
    result = strings.ToLower(strings.TrimSpace(result))

    return result
}

// Test bypass attempts — all should be normalized to detectable form:
var bypassAttempts = []struct {
    input, expectedNorm, attack string
}{
    {
        input:       "UN/**/ION SE/**/LECT",
        expectedNorm: "union select",
        attack:      "SQL comment evasion",
    },
    {
        input:       "%55NION%20%53ELECT",
        expectedNorm: "union select",
        attack:      "URL encoding",
    },
    {
        input:       "%2555NION",
        expectedNorm: "union",
        attack:      "Double URL encoding (%25 → % → %55 → U)",
    },
    {
        input:       "&#85;NION",
        expectedNorm: "union",
        attack:      "HTML entity encoding (&#85; = U)",
    },
    {
        input:       "sElEcT\t\n*\tFrOm\tusers",
        expectedNorm: "select * from users",
        attack:      "Mixed case + whitespace",
    },
    {
        input:       "\x00UNION SELECT",
        expectedNorm: "union select",
        attack:      "Null byte prefix",
    },
}

func TestNormalizationBypass() {
    for _, tc := range bypassAttempts {
        norm := NormalizationPipeline(tc.input)
        detected := strings.Contains(norm, "union") || strings.Contains(norm, "select")
        fmt.Printf("[%s] input=%q norm=%q detected=%v\n", tc.attack, tc.input, norm, detected)
    }
}
```

### Real-World Context
OWASP ModSecurity Core Rule Set (CRS) — the reference implementation used by nginx, Apache, and cloud WAFs — has 8 specific normalization transformation steps. Without normalization, a WAF that detects `UNION SELECT` is trivially bypassed by `%55NION SELECT`. This is exactly why attackers test their payloads against WAFs first.

---

## Day 74 — WAF Rule Scoring & Exceptions

### Learning Objectives
- Build a scoring WAF: block only when multiple rules match (reduce false positives)
- Implement per-IP and per-path exception lists
- Build the config hot-reload used for emergency WAF disable

### Exercise: Complete Scoring WAF

```go
package waf

import (
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "os"
    "sync"
    "time"
)

type WAFRuleConfig struct {
    Mode           string            `json:"mode"`           // off|detection|blocking
    BlockThreshold int               `json:"block_threshold"` // block if total score >= this
    CategoryThresholds map[string]int `json:"category_thresholds"`
    TrustedCIDRs   []string          `json:"trusted_cidrs"`
    ExceptionPaths []string          `json:"exception_paths"` // paths that bypass WAF
}

type ScoringWAF struct {
    rules       []Rule
    configFile  string
    config      *WAFRuleConfig
    configMu    sync.RWMutex
    trustedNets []*net.IPNet
    lastModTime time.Time
    auditor     AuditWriter
}

func NewScoringWAF(rules []Rule, configFile string, auditor AuditWriter) (*ScoringWAF, error) {
    w := &ScoringWAF{
        rules:      rules,
        configFile: configFile,
        auditor:    auditor,
    }
    if err := w.loadConfig(); err != nil {
        return nil, err
    }
    go w.watchConfig()
    return w, nil
}

func (w *ScoringWAF) loadConfig() error {
    data, err := os.ReadFile(w.configFile)
    if err != nil { return err }

    var cfg WAFRuleConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return fmt.Errorf("invalid WAF config: %w", err)
    }

    // Parse trusted CIDRs
    var nets []*net.IPNet
    for _, cidr := range cfg.TrustedCIDRs {
        _, ipNet, err := net.ParseCIDR(cidr)
        if err != nil {
            return fmt.Errorf("invalid CIDR %q: %w", cidr, err)
        }
        nets = append(nets, ipNet)
    }

    w.configMu.Lock()
    defer w.configMu.Unlock()
    w.config = &cfg
    w.trustedNets = nets

    fmt.Printf("[WAF] config loaded: mode=%s threshold=%d\n", cfg.Mode, cfg.BlockThreshold)
    return nil
}

func (w *ScoringWAF) watchConfig() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        info, err := os.Stat(w.configFile)
        if err != nil { continue }
        if info.ModTime().After(w.lastModTime) {
            w.lastModTime = info.ModTime()
            if err := w.loadConfig(); err != nil {
                fmt.Printf("[WAF] config reload error: %v\n", err)
            } else {
                fmt.Println("[WAF] config hot-reloaded")
            }
        }
    }
}

func (w *ScoringWAF) isTrusted(r *http.Request) bool {
    w.configMu.RLock()
    defer w.configMu.RUnlock()

    ip := realClientIP(r)
    parsedIP := net.ParseIP(ip)
    if parsedIP == nil { return false }

    for _, network := range w.trustedNets {
        if network.Contains(parsedIP) {
            return true
        }
    }
    return false
}

func (w *ScoringWAF) isExemptPath(path string) bool {
    w.configMu.RLock()
    defer w.configMu.RUnlock()
    for _, exempt := range w.config.ExceptionPaths {
        if strings.HasPrefix(path, exempt) {
            return true
        }
    }
    return false
}
```

**Config file: waf-config.json**
```json
{
  "mode": "blocking",
  "block_threshold": 8,
  "category_thresholds": {
    "SQLI": 5,
    "XSS": 4,
    "TRAVERSAL": 1,
    "CMDINJ": 1
  },
  "trusted_cidrs": [
    "10.0.0.0/8",
    "172.16.0.0/12"
  ],
  "exception_paths": [
    "/api/content-editor",
    "/api/html-preview"
  ]
}
```

### Real-World Context
Block threshold tuning is what WAF vendors charge professional services for. A threshold too low → high false positive rate → security team gets alert fatigue. Too high → attackers score below threshold. The exception paths for `/api/content-editor` is real — rich text editors legitimately send HTML that WAF rules trigger on.

---

## Week 10 Mini Project: Production WAF Middleware

### Feature Checklist
- [x] Normalization pipeline (URL → HTML → SQL comments → whitespace → lowercase)
- [x] SQLi detection (keyword + tautology + comment + quote balance)
- [x] XSS detection (tags + events + JS URI)
- [x] Path traversal (all encoded forms)
- [x] Command injection
- [x] SSRF detection
- [x] Body inspection (JSON + form, up to 1MB)
- [x] Scoring model (per-rule + per-category thresholds)
- [x] Exception lists (trusted CIDRs, exempt paths)
- [x] Detection mode vs blocking mode
- [x] Audit log: rule ID, matched value (normalized), score, action
- [x] Hot-reload config every 30s (no restart needed)
- [x] Response: 403 with request ID only (no rule details leaked)

### Test the WAF

```bash
# Build and start
go run . &

# Should block (SQLi):
curl "http://localhost:8080/search?q=' UNION SELECT * FROM users--"
# Response: {"error":"forbidden","request_id":"abc123"}

# Should block (XSS):
curl "http://localhost:8080/comment?text=<script>alert(1)</script>"
# Response: {"error":"forbidden","request_id":"def456"}

# Should block (path traversal):
curl "http://localhost:8080/files?path=../../../../etc/passwd"

# Should allow (legitimate):
curl "http://localhost:8080/search?q=John+O'Brien+SELECT+all+products"
# Response: 200 OK (normalized input doesn't hit threshold)

# Switch to detection mode (hot-reload):
sed -i 's/"blocking"/"detection"/' waf-config.json
sleep 35  # wait for reload
# Now malicious requests are logged but not blocked
```

### Engineer Takeaway
AWS WAF costs $5/month per rule and $0.60/million requests. Cloudflare WAF starts at $20/month. You've now built the core algorithms that power both. You can now tune a WAF instead of fighting it: write precise rule exceptions ("SQLI-001 triggers on O'Brien in search because of the apostrophe"), and reason about the trade-off between false positives and false negatives.
