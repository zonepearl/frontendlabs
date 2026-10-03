# Week 11: API Security & Authentication (Days 75–81)

**Phase:** 3 — L7 Networking + WAF + API Security
**Goal:** Build the authentication and authorization layer for your API gateway. Cover API keys, HMAC request signing, JWT-like tokens, token bucket rate limiting, input validation, PII scrubbing, and CORS — the complete API security surface.

---

> **Background in this wiki:** HTTPS guide [Ch 12 (authentication, cookies, CORS)](../../v2-https/real-life-guide-v1.md#chapter-12-security-protecting-the-entire-lifecycle) and [Ch 21 (rate limiting by identity, idempotency keys)](../../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers); [OSI guide, Layer 3](../../networking/real-life-example-osi.md#part-5-layer-3-network) on why IP addresses aren't identities. Security in Depth: [Chapter 24 (token exchange, with a Go JWT lab)](../../security/real-life-security-guide-v1.md#chapter-24-token-exchange-and-delegation), [Chapter 25 (ReBAC, with a Go lab)](../../security/real-life-security-guide-v1.md#chapter-25-distributed-authorization-rbac-to-abac-to-rebac-zanzibar), and [Part 6A (API security)](../../security/real-life-security-guide-v1.md#part-6a-api-gateway-and-perimeter-security).

## Day 75 — API Authentication Patterns

### Learning Objectives
- Implement three auth patterns as middleware: API key, HMAC signed request, ECDSA token
- Return timing-consistent responses (prevent timing attacks on auth)
- Understand when to use each pattern

### Auth Pattern Decision Matrix

| Pattern | Use When | Security Model |
|---------|----------|----------------|
| API Key header | Simple service-to-service, low risk | Secret string in header |
| HMAC signed request | Webhook validation, AWS-style auth | Shared secret, covers body hash |
| ECDSA token (JWT-like) | User authentication, expiring tokens | Asymmetric, can verify without secret |
| mTLS certificate | Zero-trust service mesh | PKI, strongest |

### Exercise: Three Auth Middlewares

```go
package auth

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "crypto/subtle"
    "encoding/hex"
    "fmt"
    "net/http"
    "strings"
    "time"
)

type ctxKey string
const (
    subjectKey ctxKey = "auth_subject"
    methodKey  ctxKey = "auth_method"
)

func GetSubject(ctx context.Context) string {
    s, _ := ctx.Value(subjectKey).(string)
    return s
}

// ── Auth Pattern 1: API Key ──────────────────────────────────────────────────

type APIKeyStore interface {
    Lookup(key string) (subject string, ok bool)
}

type MemAPIKeyStore struct {
    keys map[string]string // key → subject
}

func (s *MemAPIKeyStore) Lookup(key string) (string, bool) {
    sub, ok := s.keys[key]
    return sub, ok
}

func APIKeyAuth(store APIKeyStore) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := r.Header.Get("X-API-Key")
            if key == "" {
                // Also check Authorization: Bearer <key>
                auth := r.Header.Get("Authorization")
                if strings.HasPrefix(auth, "Bearer ") {
                    key = auth[7:]
                }
            }

            if key == "" {
                denyAuth(w, r, "missing api key")
                return
            }

            subject, ok := store.Lookup(key)
            if !ok {
                // Add artificial delay to prevent timing enumeration
                // (all invalid keys take the same time)
                time.Sleep(5 * time.Millisecond)
                denyAuth(w, r, "invalid api key")
                return
            }

            ctx := context.WithValue(r.Context(), subjectKey, subject)
            ctx = context.WithValue(ctx, methodKey, "api_key")
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// ── Auth Pattern 2: HMAC Request Signing (AWS SigV4 style) ─────────────────

type HMACAuth struct {
    keys map[string][]byte // key-id → secret
}

func (h *HMACAuth) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Expected header: Authorization: HMAC-SHA256 key-id=xxx,timestamp=xxx,signature=xxx
        authHeader := r.Header.Get("Authorization")
        if !strings.HasPrefix(authHeader, "HMAC-SHA256 ") {
            denyAuth(w, r, "missing HMAC authorization")
            return
        }

        params := parseHMACParams(authHeader[12:])
        keyID := params["key-id"]
        tsStr := params["timestamp"]
        signature := params["signature"]

        if keyID == "" || tsStr == "" || signature == "" {
            denyAuth(w, r, "malformed HMAC authorization")
            return
        }

        // Check timestamp freshness (replay attack prevention)
        ts, err := strconv.ParseInt(tsStr, 10, 64)
        if err != nil || abs(time.Now().Unix()-ts) > 300 {
            denyAuth(w, r, "timestamp out of window")
            return
        }

        secret, ok := h.keys[keyID]
        if !ok {
            time.Sleep(5 * time.Millisecond) // constant time
            denyAuth(w, r, "unknown key id")
            return
        }

        // Compute expected signature
        // Sign: method + "\n" + path + "\n" + timestamp + "\n" + body-hash
        bodyHash := sha256BodyHash(r)
        signingString := fmt.Sprintf("%s\n%s\n%s\n%s",
            r.Method, r.URL.Path, tsStr, bodyHash)
        mac := hmac.New(sha256.New, secret)
        mac.Write([]byte(signingString))
        expected := hex.EncodeToString(mac.Sum(nil))

        if subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) != 1 {
            denyAuth(w, r, "signature mismatch")
            return
        }

        ctx := context.WithValue(r.Context(), subjectKey, "key:"+keyID)
        ctx = context.WithValue(ctx, methodKey, "hmac")
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func sha256BodyHash(r *http.Request) string {
    if r.Body == nil { return hex.EncodeToString(sha256.New().Sum(nil)) }
    body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
    restoreBody(r, body)
    h := sha256.Sum256(body)
    return hex.EncodeToString(h[:])
}

// ── Auth Pattern 3: ECDSA Token (from Day 44) ───────────────────────────────

func ECDSATokenAuth(svc *TokenService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            auth := r.Header.Get("Authorization")
            if !strings.HasPrefix(auth, "Bearer ") {
                denyAuth(w, r, "missing bearer token")
                return
            }
            token := auth[7:]
            claims, err := svc.Verify(token)
            if err != nil {
                denyAuth(w, r, "invalid token")
                return
            }
            ctx := context.WithValue(r.Context(), subjectKey, claims.Subject)
            ctx = context.WithValue(ctx, methodKey, "ecdsa_token")
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func denyAuth(w http.ResponseWriter, r *http.Request, reason string) {
    // Audit log the denial
    fmt.Printf("[AUTH] denied: %s from %s path=%s\n", reason, r.RemoteAddr, r.URL.Path)

    // Return 401 — NEVER leak the specific reason in the response body
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("WWW-Authenticate", "Bearer")
    w.WriteHeader(http.StatusUnauthorized)
    fmt.Fprintf(w, `{"error":"authentication required"}`)
}
```

### Real-World Context
The 5ms artificial delay in API key rejection prevents timing enumeration — without it, an attacker can tell if a key prefix is valid by measuring response time. The HMAC 5-minute timestamp window prevents replay attacks (reusing captured requests). These aren't theoretical — they're standard requirements in AWS, Stripe, and GitHub's API auth designs.

---

## Day 76 — Rate Limiting at L7

### Learning Objectives
- Implement per-API-key token bucket rate limiting
- Return standard rate limit headers (`X-RateLimit-*`)
- Build the per-entity sliding window from Day 21 into a complete middleware

### Key Concepts

```
Token Bucket Algorithm:
  - Bucket holds up to B tokens (burst capacity)
  - Tokens refill at rate R per second
  - Each request consumes 1 token
  - If empty: return 429

Rate limit headers (standard):
  X-RateLimit-Limit:     10        (requests per window)
  X-RateLimit-Remaining: 7         (tokens left)
  X-RateLimit-Reset:     1704067200 (Unix timestamp when bucket refills)
  Retry-After:           13        (seconds to wait before retry, on 429)

Per-entity vs global:
  Per-API-key: "Stripe API allows 100 req/s"
  Per-IP:      "DDoS mitigation, before auth"
  Per-user:    "Free plan 100/day, Pro 10000/day"
```

### Exercise: Token Bucket Rate Limiter Middleware

```go
package ratelimit

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

type TokenBucket struct {
    mu       sync.Mutex
    tokens   float64
    capacity float64
    rate     float64 // tokens per second
    lastFill time.Time
}

func NewTokenBucket(capacity, ratePerSec float64) *TokenBucket {
    return &TokenBucket{
        tokens:   capacity,
        capacity: capacity,
        rate:     ratePerSec,
        lastFill: time.Now(),
    }
}

// Take tries to consume one token. Returns (allowed, remaining, resetAfter).
func (tb *TokenBucket) Take() (bool, float64, time.Duration) {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastFill).Seconds()
    tb.lastFill = now

    // Refill tokens based on elapsed time
    tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)

    if tb.tokens < 1 {
        // Calculate when 1 token will be available
        waitSecs := (1 - tb.tokens) / tb.rate
        return false, 0, time.Duration(waitSecs * float64(time.Second))
    }

    tb.tokens--
    return true, tb.tokens, 0
}

type RateLimiter struct {
    mu       sync.Mutex
    buckets  map[string]*TokenBucket
    capacity float64
    rate     float64
    cleanup  *time.Ticker
}

func NewRateLimiter(capacity, ratePerSec float64) *RateLimiter {
    rl := &RateLimiter{
        buckets:  make(map[string]*TokenBucket),
        capacity: capacity,
        rate:     ratePerSec,
        cleanup:  time.NewTicker(5 * time.Minute),
    }
    go rl.cleanupLoop()
    return rl
}

func (rl *RateLimiter) getBucket(key string) *TokenBucket {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if _, ok := rl.buckets[key]; !ok {
        rl.buckets[key] = NewTokenBucket(rl.capacity, rl.rate)
    }
    return rl.buckets[key]
}

func (rl *RateLimiter) cleanupLoop() {
    for range rl.cleanup.C {
        rl.mu.Lock()
        // Remove buckets that are full (haven't been used recently)
        for key, bucket := range rl.buckets {
            bucket.mu.Lock()
            if bucket.tokens >= bucket.capacity*0.99 {
                delete(rl.buckets, key)
            }
            bucket.mu.Unlock()
        }
        rl.mu.Unlock()
    }
}

// Middleware returns a per-entity rate limiting handler
// keyFn extracts the rate limit key from the request (API key, user ID, IP)
func (rl *RateLimiter) Middleware(keyFn func(*http.Request) string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := keyFn(r)
            bucket := rl.getBucket(key)
            allowed, remaining, retryAfter := bucket.Take()

            // Always set rate limit headers
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", rl.capacity))
            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", remaining))

            if !allowed {
                w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusTooManyRequests)
                fmt.Fprintf(w, `{"error":"rate limit exceeded","retry_after":%d}`,
                    int(retryAfter.Seconds()))
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

func main() {
    // Rate limit by API key: 10 req/s burst 50
    limiter := NewRateLimiter(50, 10)

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, `{"data":"response"}`)
    })

    // Key extraction: use API key if authenticated, else IP
    keyFn := func(r *http.Request) string {
        if subject := auth.GetSubject(r.Context()); subject != "" {
            return "user:" + subject
        }
        return "ip:" + realClientIP(r)
    }

    http.ListenAndServe(":8080",
        Chain(handler,
            APIKeyAuth(store),
            limiter.Middleware(keyFn),
        ),
    )
}
```

### Real-World Context
Stripe's API allows 100 req/s. GitHub's API allows 5000 req/hour. These are token buckets. The `X-RateLimit-*` headers are de facto standard — clients use them to implement polite backoff. Without `Retry-After`, clients retry immediately (thundering herd). With it, they wait the exact right amount.

---

## Day 77 — Input Validation

### Learning Objectives
- Build a struct tag validation framework using `reflect`
- Return ALL validation errors (not just the first) — critical for UX
- Enforce field-level constraints as a security layer

### Key Concepts

```
Validation vs Sanitization:
  Validation: check if input is acceptable → reject if not (CORRECT)
  Sanitization: modify input to make it acceptable → risky (SQL injection sanitization often fails)

Allowlist vs Denylist:
  Allowlist: only accept known-good patterns → CORRECT
  Denylist: reject known-bad patterns → always incomplete
  Rule: always use allowlist validation

Tags: validate:"required,min=1,max=255,pattern=^[a-zA-Z0-9_]+$"
```

### Exercise: Struct Tag Validator

```go
package validation

import (
    "fmt"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
    msgs := make([]string, len(ve))
    for i, e := range ve {
        msgs[i] = e.Error()
    }
    return strings.Join(msgs, "; ")
}

// Validate validates a struct using `validate` field tags
// Tags: required, min=N, max=N, pattern=regex, email, oneof=a|b|c
func Validate(v interface{}) ValidationErrors {
    val := reflect.ValueOf(v)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    if val.Kind() != reflect.Struct {
        return nil
    }

    var errs ValidationErrors
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        fieldVal := val.Field(i)
        tag := field.Tag.Get("validate")
        if tag == "" { continue }

        jsonName := field.Tag.Get("json")
        if jsonName == "" { jsonName = field.Name }
        jsonName = strings.SplitN(jsonName, ",", 2)[0]

        errs = append(errs, validateField(jsonName, fieldVal, tag)...)
    }

    return errs
}

func validateField(name string, v reflect.Value, tag string) []ValidationError {
    var errs []ValidationError
    rules := strings.Split(tag, ",")

    strVal := ""
    if v.Kind() == reflect.String {
        strVal = v.String()
    }

    for _, rule := range rules {
        rule = strings.TrimSpace(rule)
        switch {
        case rule == "required":
            if strVal == "" {
                errs = append(errs, ValidationError{name, "required field is empty"})
            }

        case strings.HasPrefix(rule, "min="):
            n, _ := strconv.Atoi(rule[4:])
            if len(strVal) < n {
                errs = append(errs, ValidationError{name, fmt.Sprintf("minimum length %d, got %d", n, len(strVal))})
            }

        case strings.HasPrefix(rule, "max="):
            n, _ := strconv.Atoi(rule[4:])
            if len(strVal) > n {
                errs = append(errs, ValidationError{name, fmt.Sprintf("maximum length %d, got %d", n, len(strVal))})
            }

        case strings.HasPrefix(rule, "pattern="):
            pattern := rule[8:]
            re, err := regexp.Compile("^(?:" + pattern + ")$")
            if err != nil { continue }
            if !re.MatchString(strVal) {
                errs = append(errs, ValidationError{name, fmt.Sprintf("invalid format (expected pattern %s)", pattern)})
            }

        case rule == "email":
            emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
            if !emailRe.MatchString(strVal) {
                errs = append(errs, ValidationError{name, "invalid email address"})
            }

        case strings.HasPrefix(rule, "oneof="):
            options := strings.Split(rule[6:], "|")
            found := false
            for _, opt := range options {
                if strVal == opt { found = true; break }
            }
            if !found {
                errs = append(errs, ValidationError{name, fmt.Sprintf("must be one of: %s", strings.Join(options, ", "))})
            }
        }
    }
    return errs
}

// Example usage:
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=64,pattern=[a-zA-Z0-9_]+"`
    Email    string `json:"email"    validate:"required,email"`
    Role     string `json:"role"     validate:"required,oneof=viewer|editor|admin"`
    Password string `json:"password" validate:"required,min=8,max=128"`
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }

    if errs := Validate(req); len(errs) > 0 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnprocessableEntity)
        json.NewEncoder(w).Encode(map[string]any{
            "error":  "validation failed",
            "fields": errs,
        })
        return
    }
    // ... handle valid request
}
```

### Real-World Context
Input validation is OWASP A03 (Injection prevention). Every field that reaches a database, file system, or external API must be validated. The allowlist pattern (`[a-zA-Z0-9_]+`) is safer than a denylist (`[^<>'";&]`). Returning all errors at once (not just first) reduces round-trips for API consumers.

---

## Day 78 — JSON Schema Enforcement

### Learning Objectives
- Implement a JSON Schema validator for API request bodies
- Reject unknown fields (`additionalProperties: false`)
- Return structured validation errors

### Key Concepts

```
Why JSON Schema validation matters:
  - additionalProperties: false blocks parameter pollution
    {"username":"alice","password":"x","admin":true} → rejected
  - type enforcement catches coercion attacks
    {"amount":"100"} when amount should be number → rejected
  - required prevents partial submissions that bypass logic

JSON Schema subset to implement:
  type:              string|number|boolean|array|object|null
  required:          ["field1","field2"]
  minLength/maxLength for strings
  minimum/maximum for numbers
  pattern:           regex for strings
  additionalProperties: false
  items:             schema for array elements
  properties:        map of field name → schema
```

### Exercise: JSON Schema Validator

```go
package jsonschema

import (
    "encoding/json"
    "fmt"
    "math"
    "regexp"
)

type Schema struct {
    Type                 string             `json:"type"`
    Required             []string           `json:"required"`
    Properties           map[string]*Schema `json:"properties"`
    AdditionalProperties *bool              `json:"additionalProperties"`
    MinLength            *int               `json:"minLength"`
    MaxLength            *int               `json:"maxLength"`
    Minimum              *float64           `json:"minimum"`
    Maximum              *float64           `json:"maximum"`
    Pattern              string             `json:"pattern"`
    Enum                 []interface{}      `json:"enum"`
    Items                *Schema            `json:"items"`
}

type SchemaError struct {
    Path    string
    Message string
}

func (e SchemaError) Error() string {
    if e.Path != "" {
        return fmt.Sprintf("%s: %s", e.Path, e.Message)
    }
    return e.Message
}

// Validate validates data against schema, returning all errors
func (s *Schema) Validate(data interface{}, path string) []SchemaError {
    var errs []SchemaError

    // Type check
    if s.Type != "" {
        if !checkType(data, s.Type) {
            return []SchemaError{{path, fmt.Sprintf("expected %s, got %T", s.Type, data)}}
        }
    }

    switch val := data.(type) {
    case map[string]interface{}:
        // Check required fields
        for _, req := range s.Required {
            if _, ok := val[req]; !ok {
                errs = append(errs, SchemaError{path + "." + req, "required field missing"})
            }
        }

        // Validate properties
        for key, fieldSchema := range s.Properties {
            if v, ok := val[key]; ok {
                errs = append(errs, fieldSchema.Validate(v, path+"."+key)...)
            }
        }

        // Reject unknown fields if additionalProperties: false
        if s.AdditionalProperties != nil && !*s.AdditionalProperties {
            for key := range val {
                if _, ok := s.Properties[key]; !ok {
                    errs = append(errs, SchemaError{path + "." + key, "additional property not allowed"})
                }
            }
        }

    case string:
        if s.MinLength != nil && len(val) < *s.MinLength {
            errs = append(errs, SchemaError{path, fmt.Sprintf("min length %d, got %d", *s.MinLength, len(val))})
        }
        if s.MaxLength != nil && len(val) > *s.MaxLength {
            errs = append(errs, SchemaError{path, fmt.Sprintf("max length %d, got %d", *s.MaxLength, len(val))})
        }
        if s.Pattern != "" {
            re, err := regexp.Compile(s.Pattern)
            if err == nil && !re.MatchString(val) {
                errs = append(errs, SchemaError{path, fmt.Sprintf("pattern %q not matched", s.Pattern)})
            }
        }

    case float64:
        if s.Minimum != nil && val < *s.Minimum {
            errs = append(errs, SchemaError{path, fmt.Sprintf("minimum %g, got %g", *s.Minimum, val)})
        }
        if s.Maximum != nil && val > *s.Maximum {
            errs = append(errs, SchemaError{path, fmt.Sprintf("maximum %g, got %g", *s.Maximum, val)})
        }
    }

    return errs
}

func checkType(data interface{}, t string) bool {
    switch t {
    case "string":  _, ok := data.(string); return ok
    case "number":  _, ok := data.(float64); return ok
    case "boolean": _, ok := data.(bool); return ok
    case "array":   _, ok := data.([]interface{}); return ok
    case "object":  _, ok := data.(map[string]interface{}); return ok
    case "null":    return data == nil
    }
    return true
}
```

### Real-World Context
`additionalProperties: false` is the parameter pollution defense. In 2020, a password reset vulnerability was found where sending `{"email":"victim","admin":true}` was accepted because there was no schema validation. The extra field was silently processed. JSON Schema with `additionalProperties: false` blocks this class of attack entirely.

---

## Day 79 — Sensitive Data Protection

### Learning Objectives
- Mask sensitive fields in logs (PII protection, GDPR compliance)
- Scrub PII from response payloads before returning to clients
- Build log filtering for SIEM-safe audit trails

### Exercise: PII Scrubber Middleware

```go
package pii

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "regexp"
    "strings"
)

// sensitiveFields — mask these JSON field values in logs
var sensitiveFields = map[string]bool{
    "password": true, "passwd": true, "secret": true,
    "token":    true, "api_key": true, "apikey": true,
    "ssn":      true, "social_security_number": true,
    "card_number": true, "cvv": true, "pin": true,
    "credit_card": true, "bank_account": true,
}

// PII patterns — redact in response bodies
var (
    emailRe    = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
    ssrRe      = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`) // SSN: 123-45-6789
    cardRe     = regexp.MustCompile(`\b\d{4}[\s\-]?\d{4}[\s\-]?\d{4}[\s\-]?\d{4}\b`)
    phoneRe    = regexp.MustCompile(`\b\+?1?\s*\(?\d{3}\)?[\s.\-]?\d{3}[\s.\-]?\d{4}\b`)
)

// MaskJSONBody replaces sensitive field values with "***"
func MaskJSONBody(data []byte) []byte {
    var parsed interface{}
    if err := json.Unmarshal(data, &parsed); err != nil {
        return data // not JSON, return as-is
    }
    masked := maskValue(parsed)
    result, err := json.Marshal(masked)
    if err != nil { return data }
    return result
}

func maskValue(v interface{}) interface{} {
    switch val := v.(type) {
    case map[string]interface{}:
        result := make(map[string]interface{}, len(val))
        for k, v := range val {
            if sensitiveFields[strings.ToLower(k)] {
                result[k] = "***"
            } else {
                result[k] = maskValue(v)
            }
        }
        return result
    case []interface{}:
        result := make([]interface{}, len(val))
        for i, v := range val {
            result[i] = maskValue(v)
        }
        return result
    default:
        return v
    }
}

// ScrubPIIFromBody replaces PII patterns with redacted markers
func ScrubPIIFromBody(data []byte) []byte {
    result := data
    result = emailRe.ReplaceAll(result, []byte("[EMAIL REDACTED]"))
    result = ssrRe.ReplaceAll(result, []byte("[SSN REDACTED]"))
    result = cardRe.ReplaceAll(result, []byte("[CARD REDACTED]"))
    result = phoneRe.ReplaceAll(result, []byte("[PHONE REDACTED]"))
    return result
}

// PIIProtectedLogger wraps a handler, masks request body in logs, scrubs response
type piiCapture struct {
    http.ResponseWriter
    buf bytes.Buffer
}

func (p *piiCapture) Write(b []byte) (int, error) {
    p.buf.Write(b)
    return p.ResponseWriter.Write(b)
}

func PIIProtectionMiddleware(scrubResponses bool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Read and mask request body for logging
            body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
            r.Body.Close()

            maskedBody := MaskJSONBody(body)
            // Log maskedBody (not original body) — password="***"
            fmt.Printf("[REQ] path=%s body=%s\n", r.URL.Path, maskedBody)

            // Restore body for handler
            r.Body = io.NopCloser(bytes.NewReader(body))

            if scrubResponses {
                // Capture and scrub response
                capture := &piiCapture{ResponseWriter: w}
                next.ServeHTTP(capture, r)
                // Note: response already written; scrubbing would require buffering first
                // In production: use a buffering ResponseWriter (more complex)
            } else {
                next.ServeHTTP(w, r)
            }
        })
    }
}
```

### Real-World Context
GDPR Article 32 requires "appropriate technical measures" to protect personal data. PCI DSS Requirement 3 says cardholder data must never be logged. A common compliance audit finding is "passwords appear in application logs" — this middleware prevents it. The `sensitiveFields` map should be reviewed quarterly and aligned with your data classification policy.

---

## Day 80 — API Abuse Prevention

### Learning Objectives
- Implement credential stuffing defense: independent per-username + per-IP rate limiting
- Add constant-time auth responses (timing attack prevention)
- Build honeypot detection for automated scanners

### Exercise: Abuse Prevention System

```go
package abuse

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

type AbuseDetector struct {
    // Credential stuffing: limit by username AND by IP independently
    failsByUser map[string]*counter
    failsByIP   map[string]*counter
    mu          sync.Mutex

    // Honeypot: any IP that hits a hidden endpoint gets scored
    honeypotsHit sync.Map // IP → hit count
}

type counter struct {
    count    int
    window   time.Time
    lockout  time.Time
}

func (ad *AbuseDetector) RecordAuthFailure(username, ip string) (lockedOut bool) {
    ad.mu.Lock()
    defer ad.mu.Unlock()

    now := time.Now()
    resetWindow := 15 * time.Minute
    lockoutDuration := 30 * time.Minute
    maxFails := 5

    for key, m := range map[string]*map[string]*counter{
        username: &ad.failsByUser,
        ip:       &ad.failsByIP,
    } {
        mp := *m
        if mp == nil {
            if key == username { ad.failsByUser = make(map[string]*counter) }
            if key == ip { ad.failsByIP = make(map[string]*counter) }
            mp = *m
        }

        if _, ok := mp[key]; !ok {
            mp[key] = &counter{window: now.Add(resetWindow)}
        }
        c := mp[key]

        if now.Before(c.lockout) {
            return true // currently locked out
        }
        if now.After(c.window) {
            c.count = 0
            c.window = now.Add(resetWindow)
        }
        c.count++
        if c.count >= maxFails {
            c.lockout = now.Add(lockoutDuration)
            fmt.Printf("[ABUSE] lockout: %s (after %d failures)\n", key, c.count)
            return true
        }
    }
    return false
}

// ConstantTimeAuthHandler wraps an auth handler to always respond in ~200ms
// Prevents timing attacks that reveal whether username exists
func ConstantTimeAuthHandler(minDuration time.Duration, h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        h.ServeHTTP(w, r)
        elapsed := time.Since(start)
        if elapsed < minDuration {
            time.Sleep(minDuration - elapsed)
        }
    })
}

// HoneypotHandler marks any client that calls this endpoint as a scanner
func HoneypotHandler(detector *AbuseDetector) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ip := realClientIP(r)
        detector.honeypotsHit.Store(ip, true)
        fmt.Printf("[HONEYPOT] scanner detected: %s hit %s\n", ip, r.URL.Path)
        // Return plausible 200 — don't alert the scanner
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        fmt.Fprint(w, `{"status":"ok","data":[]}`)
    }
}

func (ad *AbuseDetector) IsScanner(ip string) bool {
    _, ok := ad.honeypotsHit.Load(ip)
    return ok
}
```

### Real-World Context
Credential stuffing (automated login attempts with breached credentials) accounts for billions of dollars in account takeover fraud. The Cloudflare 2023 report found 34% of all login traffic was credential stuffing. Independent per-username + per-IP rate limiting is required — otherwise an attacker can test 1000 usernames at 1 req/s (below IP threshold) or test 1 username from 1000 IPs. The honeypot endpoint is how you identify automated scanners — only humans (or broken code) hit undocumented endpoints.

---

## Day 81 — CORS & Preflight

### Learning Objectives
- Implement strict origin allowlist CORS (never wildcard for authenticated APIs)
- Handle OPTIONS preflight correctly with `Vary: Origin`
- Test CORS configuration with multiple allowed and disallowed origins

### Key Concepts

```
Same-Origin Policy:
  Browser blocks JS from reading cross-origin responses by default
  Without CORS: attacker.com JS cannot read your API response even with credentials

CORS Mistakes (OWASP A05):
  Access-Control-Allow-Origin: *   with credentials → CRITICAL
    Any site can make authenticated API calls using victim's cookies
  Reflect Origin without allowlist validation:
    if origin != "" { w.Header().Set("ACAO", origin) }  → VULNERABLE
  Vary: Origin missing → CDN caches wrong origin, serves to all

Correct CORS for authenticated APIs:
  Explicit allowlist (no wildcards)
  Vary: Origin (required for CDN correctness)
  Access-Control-Allow-Credentials: true only for specific trusted origins
```

### Exercise: Strict CORS Middleware

```go
package cors

import (
    "fmt"
    "net/http"
    "strings"
)

type Config struct {
    AllowedOrigins   []string // explicit list, no wildcards
    AllowedMethods   []string
    AllowedHeaders   []string
    ExposedHeaders   []string
    AllowCredentials bool
    MaxAge           int
}

type CORS struct {
    cfg            Config
    allowedOrigins map[string]bool
}

func New(cfg Config) *CORS {
    c := &CORS{
        cfg:            cfg,
        allowedOrigins: make(map[string]bool),
    }
    for _, o := range cfg.AllowedOrigins {
        c.allowedOrigins[o] = true
    }
    return c
}

func (c *CORS) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")

        // Vary: Origin must ALWAYS be set when CORS headers might be present
        // Without it, CDN caches response for one origin and serves to all others
        w.Header().Add("Vary", "Origin")

        if origin == "" {
            // Not a CORS request — serve normally
            next.ServeHTTP(w, r)
            return
        }

        // Validate origin against allowlist — never reflect without validation
        if !c.allowedOrigins[origin] {
            // Don't set ACAO header — browser will block the response
            fmt.Printf("[CORS] blocked origin: %q\n", origin)
            next.ServeHTTP(w, r)
            return
        }

        // Set CORS headers for allowed origin
        w.Header().Set("Access-Control-Allow-Origin", origin) // never wildcard
        if c.cfg.AllowCredentials {
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }
        if len(c.cfg.ExposedHeaders) > 0 {
            w.Header().Set("Access-Control-Expose-Headers", strings.Join(c.cfg.ExposedHeaders, ", "))
        }

        // Handle preflight OPTIONS
        if r.Method == http.MethodOptions {
            w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.cfg.AllowedMethods, ", "))
            w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.cfg.AllowedHeaders, ", "))
            if c.cfg.MaxAge > 0 {
                w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", c.cfg.MaxAge))
            }
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    cors := New(Config{
        AllowedOrigins:   []string{"https://app.example.com", "https://admin.example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        ExposedHeaders:   []string{"X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           3600,
    })

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, `{"ok":true}`)
    })

    fmt.Println("Testing CORS:")
    fmt.Println("  Allowed: curl -H 'Origin: https://app.example.com' http://localhost:8080/")
    fmt.Println("  Blocked: curl -H 'Origin: https://evil.com' http://localhost:8080/")

    http.ListenAndServe(":8080", cors.Middleware(handler))
}
```

**Test:**
```bash
# Should include ACAO header (allowed):
curl -s -I -H "Origin: https://app.example.com" http://localhost:8080/
# Access-Control-Allow-Origin: https://app.example.com

# Should NOT include ACAO header (blocked):
curl -s -I -H "Origin: https://evil.com" http://localhost:8080/
# (no Access-Control-Allow-Origin header)
```

### Real-World Context
The CORS wildcard vulnerability (`Access-Control-Allow-Origin: *` with credentials) allows any website to make authenticated API calls using the victim's session. In 2018, this class of misconfiguration affected dozens of major platforms. The fix is 3 lines: use an allowlist, never reflect without validation, always set `Vary: Origin`.

---

## Week 11 Mini Project: API Security Gateway

### Complete Security Layer Stack

```
Incoming Request
    │
    ├── CORS Preflight? → Return 204 with CORS headers
    │
    ├── Rate Limit by IP (before auth — DDoS protection)
    │
    ├── Authentication: API Key OR HMAC OR ECDSA Token
    │   └── Failure → 401 (constant time)
    │
    ├── Rate Limit by authenticated user (post-auth quotas)
    │
    ├── Abuse Detection: credential stuffing, honeypot check
    │
    ├── WAF (from Week 10)
    │
    ├── JSON Schema Validation (per route)
    │
    ├── Input Validation (struct tags)
    │
    ├── PII Scrubber (mask sensitive fields in logs)
    │
    ├── Security Headers (HSTS, CSP, X-Frame-Options)
    │
    └── Application Handler
```

### Full Middleware Chain

```go
// gateway/api_security.go
func BuildAPIGateway(cfg *GatewayConfig) http.Handler {
    // Auth backends
    apiKeyStore := &MemAPIKeyStore{keys: cfg.APIKeys}
    tokenService, _ := NewTokenService()

    // Rate limiters
    ipLimiter := NewRateLimiter(50, 10)   // 50 burst, 10/s per IP
    userLimiter := NewRateLimiter(200, 50) // 200 burst, 50/s per user

    // WAF
    waf := NewScoringWAF(DefaultRules(), "waf-config.json", NewAuditLogger())

    // CORS
    cors := New(cfg.CORSConfig)

    // Abuse detector
    abuse := &AbuseDetector{}

    appHandler := buildRoutes()

    return Chain(appHandler,
        PanicRecoverer(slog.Default()),
        RequestID,
        cors.Middleware,
        ipLimiter.Middleware(func(r *http.Request) string { return "ip:"+realClientIP(r) }),
        ECDSATokenAuth(tokenService), // or APIKeyAuth(apiKeyStore)
        userLimiter.Middleware(func(r *http.Request) string { return "user:"+GetSubject(r.Context()) }),
        waf.Handler,
        PIIProtectionMiddleware(false),
        SecurityHeadersMiddleware(cfg.SecurityHeaders),
        RequestLogger(slog.Default()),
    )
}
```

### Engineer Takeaway
This is the security policy of your API platform codified in software. Every middleware maps to a specific OWASP category or compliance control. You can now review a pull request or design doc for an API service and say whether each layer is implemented correctly, and which attack it prevents.
