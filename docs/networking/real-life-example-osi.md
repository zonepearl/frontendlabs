# The OSI Model, One Click at a Time

> **One real request, followed through all seven layers.**
> You click **"Place order"** on a shopping site from your laptop on home Wi-Fi.
> About 250 ms later you see "Order confirmed". This guide follows that click
> from radio waves to JSON and back, one layer at a time.
>
> It's written for senior engineers and engineering managers. You should come
> away able to explain the model, read a packet capture, find which layer an
> incident lives in, and make architecture decisions with the network in mind.

---

> **The series:** 1 OS → 2 Networking → 3 Security → 4 HTTPS walkthrough, with Go alongside.
>
> **You are here: step 2a, Networking: the map.** ← Previous: [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md). Next: [Networking from Zero (TCP/IP)](tcp-ip/real-life-guide-v1.md) →
>
> [The full series map](#the-series-os-networking-security-https).

---

## How to read this guide

- **First pass (about 45 minutes):** read Parts 1–2, then the *"In our story"* and
  *"Why it matters to you"* boxes in each layer chapter, then Part 10 (the timeline).
- **Second pass (deep):** read every layer chapter in full, then Parts 11–15.
- **Day to day:** Part 14 (debugging playbook) and Part 16 (cheat sheet).

Each layer chapter uses the same structure:

| Section | What you get |
|---|---|
| **Job** | The layer's job, in one sentence |
| **In our story** | What happens to *your click* at this layer |
| **How it works** | The theory, at the depth you'd need in a design review |
| **Try it yourself** | Commands for macOS and Linux |
| **What breaks** | Symptom → cause table |
| **War story** | A realistic production failure |
| **Why it matters to you** | Takeaways for decisions, reviews, and incidents |

---

## The series: OS → networking → security → HTTPS

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
| 2a | **The OSI Model, One Click at a Time** ← you are here | one click followed through all seven layers; the map for everything after it | concept map, 1 hour |
| 2b | [Networking from Zero (TCP/IP)](tcp-ip/real-life-guide-v1.md) | how machines talk: addressing, routing, TCP, TLS, packet capture, network operations | 23 Go labs ([§0.8](tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself)) |
| 3a | [Security from Zero](../security/real-life-guide.md) | encoding vs hashing vs encryption, TLS, PKI, SSH, the OWASP Top 10, threat modelling | OpenSSL labs and a vulnerable app to break |
| 3b | [Security Engineering in Depth](../security/real-life-security-guide-v1.md) | cloud and Kubernetes security, distributed authorization, advanced web attacks, data protection, detection | 10 Go labs + a `govulncheck` exercise ([§0.8](../security/real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)) |
| 4 | [The HTTPS Request Lifecycle](../v2-https/real-life-guide-v1.md) | one request end to end, then the server side built in Go; Chapter 25 traces one request through every guide | 13 Go labs ([§0.6](../v2-https/real-life-guide-v1.md#0-6-the-go-labs-build-the-lifecycle-yourself), [Ch 25](../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)) |
| ∥ | [Go — The Complete Field Guide](../Golang/real-life-golang-guide.md) | the language behind every lab, plus the 120-day plan's multi-week projects | [120-day plan](../Golang/golang-90-day-plan.md) |

**Why this order.** Every network connection is a file descriptor owned by
a process, so the OS comes first. Networking comes next, because every attack and
defence in the security guides assumes you can follow a packet. Security comes
third, because it needs both. The HTTPS walkthrough comes last because a
single HTTPS request uses all of it, and [its final chapter](../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)
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

## Table of contents

1. The scenario: one click and the people and devices involved
2. The mental model: seven layers, encapsulation, and the honest truth
3. Layer 1: Physical. Turning bits into radio, light, and voltage
4. Layer 2: Data Link. Getting a frame to the next box
5. Layer 3: Network. Getting a packet across the planet
6. Layer 4: Transport. Turning packets into reliable conversations
7. Layer 5: Session. Keeping a conversation going
8. Layer 6: Presentation. Agreeing on what the bytes mean
9. Layer 7: Application. What you actually asked for
10. The full replay: a millisecond timeline of the click
11. X-ray of one packet: byte budget and hop-by-hop rewrites
12. Inside the server: from the NIC to your handler
13. The return trip and the goodbye
14. Debugging by layer: the incident playbook
15. Architecture decisions by layer: the senior and manager view
16. Cheat sheet
17. Self-check questions (with answers)

---

## Part 1 — The scenario

### The story

It's 9:14 pm. You're on the sofa with your laptop, connected to home Wi-Fi.
Your cart on `shop.example.com` holds a pair of running shoes. You click
**Place order**.

The browser's JavaScript runs:

```js
await fetch("https://shop.example.com/api/v1/orders", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "Idempotency-Key": "7f3c9a2e-1b4d-4c8e-9f00-2a6b8d1e5c43",
  },
  body: JSON.stringify({ cartId: "c_8812", paymentMethodId: "pm_77" }),
});
```

From your point of view that's one line of code. Underneath, it involves about
a dozen devices, three or four organisations, two continents' worth of routing
tables, and every layer of the OSI model, *twice*: once on the way out and once
on the way back.

### The cast of characters

All addresses use reserved documentation ranges, so they're safe to copy.

| # | Who | What it is | Key addresses |
|---|---|---|---|
| 1 | **Your laptop** | The client: browser, OS kernel, Wi-Fi chip | MAC `3c:22:fb:aa:bb:01`, IP `192.168.1.23` |
| 2 | **Home router** | Wi-Fi access point + switch + router + NAT + DHCP server + DNS forwarder in one box | LAN MAC `9c:53:22:00:00:01`, LAN IP `192.168.1.1`, WAN IP `203.0.113.7` |
| 3 | **ONT** (Optical Network Terminal) | The small box that turns the fibre's light into Ethernet | (no IP you care about) |
| 4 | **ISP access network** | OLT, aggregation switches, BNG (Broadband Network Gateway) | Your gateway as seen from the ISP |
| 5 | **ISP core and the Internet** | Routers that speak BGP to other networks | Many hops |
| 6 | **Cloud edge** | The cloud provider's border routers | |
| 7 | **Load balancer (ALB)** | An L7 proxy that terminates TLS and routes by path | Public `198.51.100.20`, private `10.0.1.15` |
| 8 | **Kubernetes node** | A VM running your pods | `10.0.2.40` |
| 9 | **`orders-service` pod** | Your code | `10.0.2.77:8080` |
| 10 | **Postgres** | The database | `10.0.3.10:5432` |
| 11 | **DNS resolvers** | Recursive resolver, root, `.com` TLD, `example.com` authoritative servers | Various |

### The map

```
 YOUR HOME                         ISP                          INTERNET              CLOUD REGION (VPC 10.0.0.0/16)
┌─────────────────────────┐   ┌────────────────────┐   ┌────────────────────┐   ┌──────────────────────────────────────────┐
│ Laptop                  │   │                    │   │                    │   │  ┌─────────┐   ┌──────────┐   ┌────────┐ │
│ 192.168.1.23            │   │ OLT ─ Aggregation  │   │  Transit / peering │   │  │   ALB   │──▶│ K8s node │──▶│  Pod   │ │
│     │  Wi-Fi (radio)    │   │        │           │   │   routers (BGP)    │   │  │ L7 proxy│   │10.0.2.40 │   │10.0.2.77│ │
│     ▼                   │   │       BNG ─ Core ──┼──▶│                    │──▶│  │10.0.1.15│   └──────────┘   └───┬────┘ │
│ Home router 192.168.1.1 │   │                    │   │                    │   │  └─────────┘                      │      │
│ NAT → 203.0.113.7       │   └────────────────────┘   └────────────────────┘   │   ▲ 198.51.100.20         ┌───────▼────┐ │
│     │  Ethernet (copper)│          ▲ fibre (light)                            │   │                       │ Postgres   │ │
│     ▼                   │          │                                          │   │                       │ 10.0.3.10  │ │
│ ONT ════════════════════┼──────────┘                                          │                           └────────────┘ │
└─────────────────────────┘                                                     └──────────────────────────────────────────┘
```

Keep this map in mind. Every layer chapter comes back to it.

---

## Part 2 — The mental model

### 2.1 The seven layers in one table

| # | Layer | Its one job | Unit of data (PDU) | "Address" it uses | Lives in |
|---|---|---|---|---|---|
| 7 | **Application** | Do the thing the user wanted (fetch a page, place an order, resolve a name) | Message | URL, hostname, API path | Your code, browser, nginx, ALB |
| 6 | **Presentation** | Agree on what the bytes *mean*: encoding, serialisation, compression, encryption | Message | — | JSON/protobuf libraries, TLS, gzip |
| 5 | **Session** | Set up, maintain, resume, and end a dialogue | Message | Session ID, cookie, TLS ticket | TLS, HTTP/2, app session logic |
| 4 | **Transport** | Process-to-process delivery; reliability, ordering, flow and congestion control | **Segment** (TCP) / **Datagram** (UDP) | **Port** | OS kernel |
| 3 | **Network** | Get a packet from any host to any host, across networks | **Packet** | **IP address** | OS kernel, routers |
| 2 | **Data Link** | Get a frame to the *next device* on the same link | **Frame** | **MAC address** | NIC driver/firmware, switches, Wi-Fi AP |
| 1 | **Physical** | Turn bits into signals on a medium and back | **Bit / symbol** | — | NIC hardware, cables, radios, optics |

Mnemonic (top-down): **A**ll **P**eople **S**eem **T**o **N**eed **D**ata **P**rocessing.
Bottom-up: **P**lease **D**o **N**ot **T**hrow **S**ausage **P**izza **A**way.

### 2.2 The rule that makes it work: each layer talks only to its neighbours and its peer

```
   Laptop                                                   Server
┌─────────────┐        "peer" conversation (logical)     ┌─────────────┐
│ L7  HTTP    │ ◀─────────────────────────────────────▶  │ L7  HTTP    │
│ L6  TLS/JSON│ ◀─────────────────────────────────────▶  │ L6  TLS/JSON│
│ L5  session │ ◀─────────────────────────────────────▶  │ L5  session │
│ L4  TCP     │ ◀─────────────────────────────────────▶  │ L4  TCP     │
│ L3  IP      │ ◀──▶ router ◀──▶ router ◀──▶ router ◀──▶ │ L3  IP      │
│ L2  Wi-Fi   │ ◀──▶ AP │ Eth ◀──▶ │ ... each link ... ◀▶│ L2  Ethernet│
│ L1  radio   │ ~~~~ ~~ │ ─── ──── │ ═══ fibre ═══ ───── │ L1  copper  │
└─────────────┘                                          └─────────────┘
```

Two ideas carry most of the weight:

1. **Vertical (service):** each layer uses only the layer below it and serves
   only the layer above it. TCP doesn't know whether it's running over Wi-Fi or
   fibre. HTTP doesn't know about MAC addresses.
2. **Horizontal (protocol):** each layer talks *logically* to the same layer on
   the other side. L4 to L7 are **end-to-end**: only the two endpoints read them.
   L1 to L3 are **hop-by-hop**: every router and switch along the way reads and
   rewrites them.

> **The single most useful sentence in networking:**
> *MAC addresses change at every hop; IP addresses stay the same end to end
> (unless NAT rewrites them); ports and everything above belong to the endpoints
> (unless a proxy terminates the connection).*

### 2.3 Encapsulation: envelopes inside envelopes

On the way down, each layer wraps the data from above in its own header (and
sometimes a trailer). On the way up, each layer removes its own wrapper and
passes the rest upward.

```
L7  HTTP request          [ POST /api/v1/orders  headers  {"cartId":"c_8812",...} ]
L6  TLS record        [TLS hdr][ ..........encrypted HTTP bytes.......... ][auth tag]
L4  TCP segment   [TCP hdr][ ..................TLS record.................... ]
L3  IP packet  [IP hdr][ ......................TCP segment.......................]
L2  Frame [MAC hdr][ ............................IP packet.......................... ][FCS]
L1  Bits  ▁▃▇▅▂▇▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃▅▇▂▁▃
```

Compare it to posting a letter. You write it (L7), translate it into a language
both of you read and seal it (L6), number the pages so they can be put back in
order (L4), put a city-and-street address on the envelope (L3), and then the
local van takes it to the next sorting office (L2) by road (L1). Each sorting
office puts it on a *different* van (new L2), but the street address (L3)
doesn't change.

### 2.4 The honest truth about OSI vs real life

The Internet doesn't run on OSI protocols. It runs on the **TCP/IP suite**,
which has four or five layers. OSI is the **vocabulary** engineers use to talk
about it. When someone says "that's an L7 load balancer" or "L2 adjacency",
they're using OSI terms for TCP/IP things.

| OSI | TCP/IP model | Real protocols in our story |
|---|---|---|
| 7 Application | Application | HTTP/1.1, HTTP/2, DNS, DHCP¹, the Postgres wire protocol |
| 6 Presentation | Application | TLS (encryption), JSON, UTF-8, gzip/brotli, HPACK |
| 5 Session | Application | TLS handshake and resumption, HTTP/2 streams, cookies, connection pools |
| 4 Transport | Transport | TCP, UDP, (QUIC sits on top of UDP) |
| 3 Network | Internet | IPv4, IPv6, ICMP, (BGP runs over TCP but *controls* L3) |
| 2 Data Link | Link | Ethernet (802.3), Wi-Fi (802.11), ARP², VLAN (802.1Q) |
| 1 Physical | Link | 1000BASE-T copper, 802.11ax radio, GPON fibre |

¹ DHCP is an L7 protocol (over UDP) whose *purpose* is to configure L3.

² ARP sits between L2 and L3. It's carried directly in Ethernet frames (L2), and its job is to map L3 addresses to L2 addresses.

**Where the model is fuzzy, and why that's fine:**
- **TLS** is variously called L4.5, L5, L6, or L7. Its handshake behaves like a
  session (L5), and its record encryption behaves like presentation (L6).
- **QUIC / HTTP/3** combines L4, L5, and L6 into one protocol that runs in user
  space over UDP.
- **Routing protocols** (BGP, OSPF) carry L3 *control* information, but BGP is
  itself carried over TCP.

Use the model for what it's good at: **locating a problem, splitting
responsibility, and naming things precisely.** Don't argue about which layer TLS
belongs to in a design review. Ask *"which component terminates it, and what
can that component see?"*

### 2.5 Layers are not steps in time

People often assume L1 happens first and L7 happens last. Here's the actual
order of events for our click:

```
Long before the click (seconds to hours earlier):
  L1   Wi-Fi radio links up with the router
  L2   Wi-Fi association + WPA3 authentication
  L3   DHCP gives the laptop 192.168.1.23 (DHCP itself is L7-over-UDP!)
  L2/3 ARP learns the router's MAC

At the click (t = 0):
  L7   DNS lookup for shop.example.com      ← an L7 protocol, needed BEFORE...
  L4   TCP handshake                        ← ...the L4 connection can start
  L5/6 TLS handshake
  L7   HTTP POST
       ... and every single one of those uses L1–L3 for every packet.
```

The layers describe **nesting**, not chronology. The chapters below go
bottom-up, as you asked, but they keep pointing out *when* each thing happens.

---

## Part 3 — Layer 1: Physical

### Job

Turn a stream of 1s and 0s into something that can travel through a medium
(radio waves, voltage on copper, pulses of light), and turn it back into bits
at the other end with as few errors as possible.

### In our story

Your click crosses **three completely different physical media** before it
even leaves the neighbourhood, and several more after that:

| Hop | Medium | Standard | What's physically happening |
|---|---|---|---|
| Laptop → router | Radio, 5 GHz band | Wi-Fi 6 (802.11ax) | Bits are encoded as changes in the amplitude and phase of a radio wave (1024-QAM: each symbol carries 10 bits), spread over many sub-carriers (OFDM/OFDMA) |
| Router → ONT | Copper, Cat6 cable | 1000BASE-T | Four twisted pairs at once, each carrying 5-level voltage symbols (PAM-5) at 125 Mbaud, in both directions simultaneously |
| ONT → ISP | Single-mode fibre | GPON | Pulses of laser light: 1490 nm downstream (2.488 Gb/s shared), 1310 nm upstream (1.244 Gb/s shared) |
| ISP core / Internet | Fibre | 100G/400G Ethernet over DWDM | Many wavelengths ("colours") on one fibre, each carrying 100–400 Gb/s using PAM4 or coherent optics |
| Inside the cloud data centre | Fibre + copper DACs | 25G/100G Ethernet | Short-reach optics or direct-attach copper |

### How it works

**1. Bits become symbols.** Physical media can't carry "a 1". They carry a
*symbol*: a voltage level, a light intensity, or a radio wave with a particular
amplitude and phase. The more distinct symbols a link can reliably tell apart,
the more bits each symbol carries:

```
  2 levels  (NRZ)    → 1 bit / symbol
  4 levels  (PAM4)   → 2 bits / symbol
  1024-QAM (Wi-Fi 6) → 10 bits / symbol   ← needs a very clean signal (high SNR)
```

That's why your Wi-Fi speed drops as you walk away from the router. The
**signal-to-noise ratio (SNR)** falls, so the radios automatically fall back to
simpler modulations (a lower **MCS** index) that carry fewer bits per symbol but
survive noise better.

**2. Clocking and line coding.** The receiver has to know *when* to sample. Line
codes (8b/10b, 64b/66b) and scrambling make sure the signal changes often enough
for the receiver to recover the sender's clock from the data itself.

**3. Error correction.** Modern links add **Forward Error Correction (FEC)**:
redundant bits that let the receiver *fix* small errors without retransmitting.
This adds a little latency and costs nothing at higher layers. Anything FEC
can't fix becomes a bad frame at L2.

**4. The three delays L1 controls.** Most latency discussions mix these up:

| Delay | Formula | Our numbers |
|---|---|---|
| **Propagation** (speed of light in the medium) | distance ÷ (~200,000 km/s in fibre) | ≈ **5 µs per km**. A region 1,500 km away costs ≥ 7.5 ms one way, ≥ 15 ms round trip, *before any equipment* |
| **Serialisation** (time to put the bits on the wire) | frame size ÷ link rate | 1,500 B at 1 Gb/s = **12 µs**. At 10 Mb/s upstream = **1.2 ms** |
| **Duplex / medium access** | Wi-Fi is half-duplex and shared: devices take turns | Under contention it adds anywhere from 1 ms to 50 ms or more |

**Propagation delay can't be engineered away.** You can buy more bandwidth.
You can't buy a faster speed of light. That's the physical reason CDNs, edge
regions, and "deploy close to users" exist.

**5. Wi-Fi is the weakest link in most user journeys.** It's shared, half-duplex,
subject to interference (microwaves, neighbours' networks on the same channel),
and it retransmits at L2 when frames are lost. That turns packet loss into
**latency jitter** that higher layers see as random slowness.

### Try it yourself

```bash
# macOS: Wi-Fi signal (RSSI), noise, channel, PHY mode, transmit rate
sudo wdutil info | sed -n '/WIFI/,/BLUETOOTH/p'
system_profiler SPAirPortDataType | head -40

# macOS: end-to-end capacity + responsiveness under load ("bufferbloat")
networkQuality

# Linux: wired link speed, duplex, auto-negotiation
ethtool eth0

# Linux: L1/L2 error counters. Rising CRC/rx_errors = physical problem
ip -s link show eth0
ethtool -S eth0 | grep -Ei 'err|crc|drop|fcs'

# Linux: Wi-Fi signal and bitrate
iw dev wlan0 link
```

What healthy looks like: Wi-Fi RSSI better than −65 dBm, noise around −90 dBm
(an SNR of 25 dB or more), wired links at their full speed in full duplex, and
**zero** CRC errors that keep increasing.

### What breaks

| Symptom you see higher up | Physical cause |
|---|---|
| "Wi-Fi is slow in the bedroom", video calls freeze | Low SNR → low MCS rate + L2 retransmissions |
| TCP throughput is fine at night and terrible at 8 pm | Channel congestion (shared medium), or the shared GPON/cable segment is saturated |
| Steady trickle of CRC errors, slow transfers | Damaged cable, dirty fibre connector, failing optic (check DOM/DDM light levels) |
| A link that keeps going up and down | Bad cable, auto-negotiation fight, overheated transceiver |
| A server NIC stuck at 100 Mb/s | Cable with a broken pair. 1000BASE-T needs all 4 pairs and silently falls back to 100BASE-TX on 2 |

### War story: "The database is slow on Tuesdays"

A team spent two weeks tuning queries on a database whose p99 latency spiked at
random. In the end someone ran `ethtool -S` on the database host and found
`rx_crc_errors` rising steadily. A fibre patch lead had been crushed under a
rack door. About 0.1% of frames were corrupted, which meant 0.1% packet loss,
which made TCP retransmit after a 200 ms or longer timeout. The "Tuesday"
pattern was the weekly batch job: more traffic meant more corrupted frames.
**Lesson:** when the latency *tail* is bad and the *median* is fine, check for
packet loss first, and check the error counters at the bottom of the stack
before profiling the top.

### Why it matters to you

- **Latency has a physical floor.** If your SLO needs a p50 under 50 ms
  worldwide, no amount of code optimisation will get you there from one region.
  That's an architecture decision (CDN, multi-region, edge compute), not a
  performance ticket.
- **Packet loss usually starts physical, but shows up as an application tail-latency
  problem.** Make "check the interface error counters" step 1 of your runbooks.
- **"Bandwidth" and "latency" are different budgets.** More bandwidth doesn't make
  a chatty API with 30 sequential calls any faster.
- In the cloud you never touch L1. You just need to know it exists. Choosing
  an instance type with "up to 10 Gbps" networking *is* an L1/L2 capacity
  decision, and "up to" often means burst credits.

### Go deeper in the series

- **Physical and Wi-Fi in more detail:** [TCP/IP guide, Chapter 8 (Wi-Fi)](tcp-ip/real-life-guide-v1.md#chapter-8-wi-fi-same-idea-different-physics)
  and [Chapter 27 (latency, bandwidth, bufferbloat)](tcp-ip/real-life-guide-v1.md#chapter-27-latency-bandwidth-and-bufferbloat).
- **Interface counters on a server:** [Linux guide, Chapter 37](../os-linux/real-life-os-guide.md#chapter-37-how-a-linux-box-sees-the-network).


---

## Part 4 — Layer 2: Data Link

### Job

Get a **frame** from one device to the **next device on the same link**,
identify both by **MAC address**, and detect corrupted frames. L2 never sees
past the next hop.

### In our story

Before you even clicked, L2 had already done a lot of work, and during the
click it does a little more at **every hop**.

**Earlier that evening (bootstrapping the link):**

1. **Association.** The laptop's Wi-Fi chip scans, finds the SSID
   `HomeNet-5G` (the router's radio has a MAC/BSSID), and *associates*.
2. **Authentication (WPA3-SAE + 4-way handshake).** The laptop and router
   prove they both know the passphrase *without sending it*, and derive
   per-session encryption keys. **Wi-Fi frames are encrypted at L2**, so a
   neighbour sniffing the air sees only ciphertext.
3. **ARP for the gateway.** Once DHCP (Part 5) has given the laptop an IP
   address and told it "your gateway is `192.168.1.1`", the laptop needs the
   gateway's *MAC* to build frames:

```
Laptop  → broadcast ff:ff:ff:ff:ff:ff  ARP: "Who has 192.168.1.1? Tell 192.168.1.23"
Router  → 3c:22:fb:aa:bb:01            ARP: "192.168.1.1 is at 9c:53:22:00:00:01"
```

The answer is cached (the **ARP table**) for minutes, so the click itself
doesn't trigger a new ARP request.

**At the click: one frame, rebuilt at every hop.** The laptop's first outgoing
packet to the cloud is wrapped like this:

```
Wi-Fi hop (802.11 frame, as transmitted over the air):
  Addr1 (receiver / BSSID)  = 9c:53:22:00:00:01   router radio
  Addr2 (transmitter)       = 3c:22:fb:aa:bb:01   laptop
  Addr3 (final L2 dest)     = 9c:53:22:00:00:01   router (gateway)
  [ encrypted payload: IP packet src 192.168.1.23 → dst 198.51.100.20 ]

Router WAN → ONT → ISP (Ethernet II frame):
  Dst MAC = <ISP BNG's MAC>
  Src MAC = <router WAN MAC>
  EtherType = 0x0800 (IPv4)
  [ IP packet, now src 203.0.113.7 (after NAT) → dst 198.51.100.20 ]
  FCS (CRC-32)

... a new frame on every subsequent link ...

Cloud: final hop to ALB
  Dst MAC = <ALB ENI MAC>   Src MAC = <VPC virtual router MAC>
```

**The IP packet inside stays essentially the same the whole way (only NAT and
TTL change it). The L2 envelope is thrown away and rebuilt at every router.**

### How it works

**1. The Ethernet II frame.** This is the most common structure in networking:

```
┌──────────┬─────┬─────────┬─────────┬───────────┬───────────┬──────────────┬─────┬─────┐
│ Preamble │ SFD │ Dst MAC │ Src MAC │ 802.1Q tag│ EtherType │   Payload    │ FCS │ IFG │
│  7 B     │ 1 B │  6 B    │  6 B    │ 4 B (opt) │   2 B     │ 46–1500 B    │ 4 B │12 B │
└──────────┴─────┴─────────┴─────────┴───────────┴───────────┴──────────────┴─────┴─────┘
   L1 sync                                        0x0800 IPv4                  CRC   gap
                                                  0x86DD IPv6
                                                  0x0806 ARP
```

- **MAC address:** 48 bits. The first 24 bits are the vendor's **OUI**. Modern
  phones and laptops *randomise* their MAC per network for privacy, which
  breaks "identify devices by MAC" schemes.
- **EtherType** tells the receiver which L3 protocol is inside. This is the
  hand-off point between layers.
- **FCS** is a CRC-32. A frame that fails the check is **silently dropped**. L2
  doesn't retransmit on Ethernet; TCP (L4) has to notice the gap. *Wi-Fi does*
  retransmit at L2, because the air loses so many frames.
- **MTU (1500 bytes)** is the maximum *payload*. It's the most important number
  that crosses layers: IP and TCP size themselves to fit it.

**2. Switches learn, flood, and forward.**

```
Frame arrives on port 3 from MAC A, destined for MAC B
 1. LEARN:   "MAC A lives on port 3"     (MAC/CAM table)
 2. LOOKUP:  do I know where MAC B is?
       yes → FORWARD out that one port
       no  → FLOOD out every port except 3  (B will answer, and I'll learn)
 3. Broadcast (ff:ff:ff:ff:ff:ff) → always flood
```

Everything that receives a given broadcast is in one **broadcast domain**. ARP,
DHCP discovery, and some service discovery depend on broadcast, which is why
it doesn't scale: a broadcast domain with 5,000 hosts is constantly busy with
ARP traffic.

**3. VLANs (802.1Q)** add a 4-byte tag with a 12-bit VLAN ID (1–4094), so one
physical switch can act as many *separate* broadcast domains (for example
"guest Wi-Fi" and "home"). VLANs are the L2 version of network segmentation.

**4. Loops are fatal at L2.** Ethernet frames have **no TTL**. If you create a
loop between switches, broadcast frames circulate forever and multiply into a
**broadcast storm** that takes down the whole segment within seconds. **Spanning
Tree Protocol (STP)** blocks redundant links to prevent this. Modern data centres
avoid big L2 domains altogether (leaf-spine with L3 routing down to the rack).

**5. L2 in the cloud is mostly an illusion.** Your VPC *looks* like Ethernet,
but:
- There's no real broadcast. The hypervisor or smart NIC (for example AWS
  Nitro) answers ARP itself and delivers frames directly.
- The "virtual router" at `.1` of every subnet is software in the host.
- Kubernetes overlays (Flannel, Calico in VXLAN mode) put a **whole L2 frame
  inside a UDP packet** (VXLAN, port 4789) to build a virtual L2 network on top
  of L3. That costs **50 bytes** per packet, which is why pod MTUs are often
  1450 instead of 1500.

### Try it yourself

```bash
# ARP/neighbour table: who's on my LAN and what's their MAC?
arp -a                       # macOS / Linux
ip neigh show                # Linux

# Watch ARP in real time (then run: ping 192.168.1.1 in another terminal)
sudo tcpdump -i en0 -n -e arp            # -e prints the MAC (L2) headers

# See the L2 header of every packet to the shop
sudo tcpdump -i en0 -n -e host shop.example.com

# Interface MTU
ifconfig en0 | grep mtu      # macOS
ip link show eth0            # Linux

# Linux bridge / switch MAC table (on a Linux router or K8s node)
bridge fdb show | head
```

### What breaks

| Symptom | L2 cause |
|---|---|
| Traffic to an IP goes to the *old* server for minutes after a failover | Stale ARP caches. The new owner must send a **gratuitous ARP** (keepalived/VRRP does this) |
| Whole office network dies within seconds of someone plugging in a cable | L2 loop → broadcast storm (STP disabled or misconfigured) |
| Two hosts on "the same subnet" can't reach each other | Different VLANs; trunk port not carrying the VLAN |
| Small requests work, large responses hang | **MTU mismatch** (for example a VXLAN overlay with MTU 1500 inside). Large frames get dropped |
| Intermittent connectivity, log shows "MAC flapping" | Duplicate MAC (cloned VM) or a loop |
| Device keeps getting a new IP / gets blocked by MAC filtering | MAC randomisation on the client |

### War story: "Big JSON responses time out, small ones are fine"

After moving to a new Kubernetes CNI, `GET /api/orders?limit=10` worked but
`?limit=500` hung until timeout. The overlay used VXLAN (50 bytes of overhead)
but the pod interfaces were still set to MTU 1500. Small responses fit in one
frame. Large responses produced full 1500-byte packets, which became 1550
bytes after encapsulation, exceeded the underlying 1500 MTU, and were dropped,
because the ICMP "too big" messages that should have triggered Path MTU
Discovery were filtered. **Fix:** set the pod MTU to 1450 (or enable jumbo
frames, 9001 in AWS, on the underlay). **Lesson:** *"works for small payloads,
hangs for large ones" is almost always MTU.*

### Why it matters to you

- **"Same subnet" means "same L2 domain."** Hosts can talk directly without a
  router. Anything that relies on that (multicast discovery, VRRP failover,
  some legacy clustering) **won't work the same way in the cloud**, so check
  for it during migrations.
- **Failover designs that move an IP need L2 cooperation** (gratuitous ARP). In
  the cloud you use the provider's API instead (move an Elastic IP or ENI),
  which takes seconds, not milliseconds. Put that in your RTO.
- **MTU is a cross-team contract.** Overlays, VPNs, tunnels, and PPPoE all eat
  into it. Write it down for every network path.
- **L2 encryption (WPA3, MACsec) protects only one hop.** It doesn't replace
  TLS.

### Go deeper in the series

- **Ethernet, ARP, switches, VLANs:** [TCP/IP guide, Chapters 5–7](tcp-ip/real-life-guide-v1.md#chapter-5-mac-addresses-and-ethernet-frames),
  with ARP captured live.
- **MTU and overlays:** [TCP/IP Chapter 15](tcp-ip/real-life-guide-v1.md#chapter-15-mtu-the-maximum-size-of-a-packet) and the
  tunnels lab in [Chapter 54](tcp-ip/real-life-guide-v1.md#chapter-54-vpns-and-tunnels-wireguard-ipsec-gre-overlays-and-mtu-traps), where a WireGuard tunnel measured
  exactly 1,420 bytes.
- **Bridges and veth pairs, built by hand:** [TCP/IP Chapter 52](tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack).


---

## Part 5 — Layer 3: Network

### Job

Deliver a **packet** from any host to any other host *across many networks*,
using hierarchical **IP addresses** and **routing**. L3 is "best effort": it can
drop, duplicate, or reorder packets and doesn't tell anyone. Reliability is the
job of a higher layer.

### In our story

**Earlier that evening, DHCP gave the laptop its L3 identity** (the "DORA"
exchange, over UDP 68 → 67):

```
Laptop → broadcast   DHCP DISCOVER  "I'm 3c:22:fb:aa:bb:01, I need an address"
Router → laptop      DHCP OFFER     "Try 192.168.1.23"
Laptop → broadcast   DHCP REQUEST   "I'll take 192.168.1.23"
Router → laptop      DHCP ACK       IP 192.168.1.23/24, gateway 192.168.1.1,
                                    DNS 192.168.1.1, lease 24h
```

**At the click, the laptop's kernel makes a routing decision** for
`198.51.100.20` (the ALB's IP, which DNS supplied; see Part 9):

```
Is 198.51.100.20 inside my subnet 192.168.1.0/24?   No.
→ Use the default route: send it to the gateway 192.168.1.1
→ (L2) wrap it in a frame addressed to the gateway's MAC
```

**The packet then travels about 12–18 router hops:**

```
 hop  device                          what it does with our packet
 ───  ──────────────────────────────  ───────────────────────────────────────────────
  1   Home router 192.168.1.1         TTL 64→63, NAT: src 192.168.1.23:54012
                                         → 203.0.113.7:61001, forward to ISP
  2   ISP BNG                         TTL 63→62, route lookup, forward
 3–6  ISP aggregation/core            longest-prefix match on 198.51.100.0/24
                                         → "learned via BGP from the cloud's network"
 7–12 Transit / peering exchange      BGP-chosen path between autonomous systems
13–15 Cloud border + backbone         into the region, towards the ALB's subnet
 16   ALB                             destination reached; hand to L4
```

### How it works

**1. The IPv4 header (20 bytes, field by field):**

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
┌───────┬───────┬───────────┬───┬───────────────────────────────┐
│Version│  IHL  │   DSCP    │ECN│         Total Length          │  4 | 5 | 0 | 0 | 1500
├───────┴───────┴───────────┴───┼─────┬─────────────────────────┤
│        Identification         │Flags│    Fragment Offset      │  0x1c46 | DF | 0
├───────────────┬───────────────┼─────┴─────────────────────────┤
│      TTL      │   Protocol    │        Header Checksum        │  64 | 6 (TCP) | ...
├───────────────┴───────────────┴───────────────────────────────┤
│                    Source IP: 192.168.1.23                    │
├───────────────────────────────────────────────────────────────┤
│                 Destination IP: 198.51.100.20                 │
└───────────────────────────────────────────────────────────────┘
```

The fields that matter in real life:

| Field | Why you care |
|---|---|
| **TTL** | Decremented at every router. At 0 the packet is dropped and an ICMP "Time Exceeded" is sent back. That's what makes routing loops survivable, and it's how `traceroute` works. Default 64 (Linux/macOS), 128 (Windows) |
| **Protocol** | 6 = TCP, 17 = UDP, 1 = ICMP. The hand-off to L4 |
| **DF (Don't Fragment)** | Set by modern stacks. A router that can't forward a too-big DF packet drops it and sends "ICMP Fragmentation Needed". **Path MTU Discovery** depends on that ICMP getting back to the sender |
| **DSCP** | Quality-of-service marking. Respected inside enterprise/carrier networks, mostly ignored across the public Internet |
| **Src/Dst** | The end-to-end identity, *except where NAT rewrites it* |

**2. Routing = longest-prefix match.** Every router has a table of prefixes and
picks the **most specific** one that matches:

```
Destination 198.51.100.20
  0.0.0.0/0         → upstream            (matches, /0)
  198.51.100.0/22   → peer A              (matches, /22)
  198.51.100.0/24   → peer B              (matches, /24)  ← WINNER: most specific
```

**3. BGP: how the Internet agrees on routes.** The Internet is about 75,000
**Autonomous Systems** (ASes: ISPs, clouds, large companies). Each one
*announces* the prefixes it owns to its neighbours via **BGP**. Your ISP knows
how to reach `198.51.100.0/24` because the cloud provider's AS announced it and
that announcement spread. BGP chooses paths by *policy and business
relationships*, not by shortest latency. That's why traffic sometimes takes
odd detours.

**4. NAT: the home router's trick.** IPv4 has only about 4.3 billion addresses,
so your whole home shares **one** public IP. The router keeps a translation
table:

```
Inside (private)            Outside (public)          Remote
192.168.1.23:54012   ⇄   203.0.113.7:61001   ⇄   198.51.100.20:443
192.168.1.40:51877   ⇄   203.0.113.7:61002   ⇄   142.250.x.x:443   (your phone)
```

Consequences for software engineers:
- NAT is **stateful**. Idle mappings expire (often after 30 s to 5 min for TCP),
  which silently breaks long-idle connections. That's why WebSockets and gRPC
  streams need **keepalives**.
- Many mobile and some broadband users sit behind **Carrier-Grade NAT** (the
  `100.64.0.0/10` range): thousands of customers share one IP. **Rate-limiting
  or blocking by IP punishes innocent users.**
- Servers can't initiate connections to clients behind NAT. That's why push
  notifications, WebRTC (STUN/TURN), and "phone home" agents exist.

**5. ICMP: L3's messaging channel.** `ping` (Echo), `traceroute` (Time Exceeded),
and Path MTU Discovery (Fragmentation Needed) all depend on it. **Blocking all
ICMP is a classic misconfiguration** that causes the MTU black-hole from Part 4.

**6. IPv6 in one paragraph.** 128-bit addresses, so every device can have a
global address and **no NAT is needed**. A fixed 40-byte header, no
fragmentation by routers, and ARP is replaced by Neighbour Discovery (ICMPv6).
Browsers try IPv6 and IPv4 in parallel (**Happy Eyeballs**) and use whichever
connects first. If your AAAA record points to a broken IPv6 path, some users
get a 250 ms penalty on every new connection.

**7. L3 in the cloud.**

| Cloud concept | OSI reality |
|---|---|
| VPC `10.0.0.0/16` | Your private L3 address space |
| Subnet `10.0.2.0/24` | An L3 prefix, pinned to one availability zone |
| Route table | Static routes: `0.0.0.0/0 → internet gateway` or `→ NAT gateway` |
| Internet gateway | 1:1 NAT between a public IP and an instance's private IP |
| NAT gateway | Many-to-one NAT (like your home router) for outbound traffic from private subnets |
| Security group | **Stateful** L3/L4 firewall: allow `tcp/443 from 0.0.0.0/0`; return traffic is automatic |
| Network ACL | **Stateless** L3/L4 filter per subnet: you must allow the return (ephemeral) ports too |
| VPC peering / Transit Gateway | Routing between L3 address spaces. **Requires non-overlapping CIDRs** |

### Try it yourself

```bash
# My IP, subnet, gateway
ipconfig getifaddr en0; netstat -rn | grep default      # macOS
ip addr show; ip route show                              # Linux
ipconfig getpacket en0                                   # macOS: full DHCP lease

# Which route would the kernel use for this destination?
route -n get 198.51.100.20                               # macOS
ip route get 198.51.100.20                               # Linux

# My public IP (what the server sees after NAT)
curl -s https://ifconfig.me; echo

# The path, hop by hop (TTL trick). mtr = traceroute + live loss/latency stats
traceroute -n shop.example.com
mtr -rwc 50 shop.example.com

# Find the path MTU: largest packet that passes with DF set
ping -D -s 1472 shop.example.com      # macOS (-D sets DF)
ping -M do -s 1472 shop.example.com   # Linux. 1472 + 8 ICMP + 20 IP = 1500
```

### What breaks

| Symptom | L3 cause |
|---|---|
| `ping 192.168.1.1` works, `ping 1.1.1.1` doesn't | Default route missing, or upstream/ISP down |
| VPC peering request rejected / some services unreachable | **Overlapping CIDRs** (both teams picked `10.0.0.0/16`) |
| Connections hang only for some users / large payloads | PMTU black-hole (ICMP blocked) |
| Traffic works one way, replies vanish | **Asymmetric routing** through a stateful firewall that only sees half the flow |
| "Packets to region X take a crazy path" | BGP policy; route leak; ISP peering dispute |
| Outbound calls from private subnet fail | No route to the NAT gateway; NAT gateway in the wrong AZ/subnet |
| Rate limiter blocks a whole office/mobile carrier | Many users behind one NAT/CGNAT IP |

### War story: "We can't connect the acquisition's network"

Company A bought Company B. Both had built their entire cloud estate in
`10.0.0.0/16`. VPC peering, Transit Gateway, and VPN all require non-overlapping
address space. The fix took **two quarters**: re-IP'ing hundreds of services,
redoing firewall rules and DNS, and running a NAT translation layer in the
meantime. **Lesson for managers:** IP address planning is a company-wide,
decade-long decision. Allocate CIDRs centrally (IPAM) *from day one*, even if
you're "just a startup".

### Why it matters to you

- **An IP address is not a user identity.** NAT, CGNAT, proxies, and VPNs
  make sure of that. Design rate limits and fraud rules accordingly.
- **IP plans are architecture.** Overlapping CIDRs are one of the most expensive
  network mistakes to fix.
- **Don't block ICMP wholesale.** Allow at least "Fragmentation Needed" / "Packet
  Too Big".
- **Know which of your firewalls are stateful and which are stateless.**
  Security groups and NACLs behave very differently, and teams lose days to
  that difference.
- **Long-lived idle connections die in NAT tables you don't control.** Add
  application-level keepalives shorter than the smallest NAT timeout on the path.

### Go deeper in the series

- **Subnetting with exercises**, plus a Go subnet calculator and IPAM allocator:
  [TCP/IP guide, Chapters 10–11](tcp-ip/real-life-guide-v1.md#chapter-10-subnet-masks-the-part-everyone-finds-hard).
- **See your own NAT** with the `whoami` lab: [TCP/IP Chapter 12](tcp-ip/real-life-guide-v1.md#chapter-12-private-addresses-and-nat).
- **Routing protocols** (OSPF, BGP, VRRP): [TCP/IP Part 12](tcp-ip/real-life-guide-v1.md#part-12-how-routers-learn-routes);
  anycast and global routing: [HTTPS guide, Chapter 7](../v2-https/real-life-guide-v1.md#chapter-7-nat-vip-bgp-anycast-vrrp-global-traffic-routing).
- **IP planning and overlapping CIDRs** (the war story above, with tooling):
  [TCP/IP Chapter 61](tcp-ip/real-life-guide-v1.md#chapter-61-network-design-refactoring-migrations-diagrams-ipam-change-safety).


---

## Part 6 — Layer 4: Transport

### Job

Turn the unreliable, host-to-host packets of L3 into **process-to-process**
communication, identified by **ports**. TCP adds reliability, ordering, flow
control, and congestion control. UDP adds only ports and a checksum.

### In our story

The browser asks the kernel to `connect()` to `198.51.100.20:443`. The kernel
picks an **ephemeral source port** (54012) and performs the **three-way
handshake**:

```
  Laptop 192.168.1.23:54012                                ALB 198.51.100.20:443
        │                                                          │
   t=0  │──── SYN  seq=3001829117  win=65535                    ──▶│
        │          options: MSS=1460, SACK-permitted,              │
        │                   WScale=6, Timestamps                   │
        │                                                          │
  t=40ms│◀─── SYN-ACK seq=901120044 ack=3001829118  MSS=1460 ... ──│
        │                                                          │
  t=40ms│──── ACK  ack=901120045   ─────────────────────────────▶ │  (connection ESTABLISHED)
        │──── [immediately followed by TLS ClientHello] ─────────▶ │
```

One full **round trip (RTT ≈ 40 ms)** has passed and *no application data has
moved yet*. That's the cost of reliability, and the reason connection reuse
matters so much.

The connection is now identified everywhere by its **5-tuple**:

```
(protocol=TCP, src=192.168.1.23, sport=54012, dst=198.51.100.20, dport=443)
```

…which the home router's NAT rewrites to `(TCP, 203.0.113.7, 61001,
198.51.100.20, 443)` for everyone on the Internet side.

**Important:** the ALB is an L7 proxy, so **this TCP connection ends at the
ALB.** The ALB opens a *separate* TCP connection from `10.0.1.15:<ephemeral>`
to the pod `10.0.2.77:8080`, and it usually reuses one it already has open.
The pod never sees your laptop's TCP connection (or your IP, except in the
`X-Forwarded-For` header).

### How it works

**1. TCP vs UDP.**

| | TCP | UDP |
|---|---|---|
| Connection | Yes (handshake, state) | No |
| Reliability and ordering | Yes: sequence numbers, ACKs, retransmissions | No: the app handles loss |
| Flow/congestion control | Yes | No (the app must be polite) |
| Header | 20–60 B | 8 B |
| Message boundaries | **No**: it's a byte stream | Yes: one send = one datagram |
| Used by | HTTP/1.1, HTTP/2, TLS, Postgres, SSH | DNS, DHCP, QUIC/HTTP/3, video calls, games, VXLAN |

> "TCP is a byte stream" means: if you `send()` 100 bytes twice, the other side
> may `recv()` 200 bytes once, or 37 then 163. **Every TCP protocol must frame its
> own messages** (HTTP uses `Content-Length` / chunked encoding; gRPC uses a
> length prefix). Forgetting this is a classic bug in hand-written protocols.

**2. The TCP header and the fields that matter:**

```
Src port 54012 | Dst port 443
Sequence number   — byte offset of this segment's first byte in the stream
Ack number        — "I've received everything up to byte N-1, send N next"
Flags             — SYN, ACK, FIN (graceful close), RST (abort), PSH, ECE/CWR (ECN)
Window            — receive window: how many more bytes I can buffer (flow control)
Options           — MSS, Window Scale, SACK, Timestamps
```

**3. Reliability.** Every byte has a sequence number. The receiver ACKs what it
has received. If an ACK doesn't arrive within the **retransmission timeout
(RTO)**, or the sender sees **three duplicate ACKs** (fast retransmit), it sends
the data again. Linux's minimum RTO is **200 ms**, and the initial RTO before
any RTT measurement (the SYN) is **1 second**. So:

> **One lost packet in the wrong place costs 200 ms to 1 s or more.** A lost SYN
> adds a full second to connection setup. This is why 0.1% packet loss
> damages p99 latency far more than you'd expect.

**4. Flow control vs congestion control.** People mix these up:

| | Flow control | Congestion control |
|---|---|---|
| Protects | The **receiver** (don't overflow its buffer) | The **network** (don't overflow routers in between) |
| Mechanism | Receiver advertises a window (`rwnd`) | Sender estimates a congestion window (`cwnd`) |
| Effective send limit | `min(rwnd, cwnd)` | |

**Slow start:** a new connection starts with `cwnd = 10 segments` (~14 KB) and
roughly doubles every RTT until it sees loss (CUBIC, the Linux default) or until
it models the path's bandwidth (BBR). Consequence:

```
Transferring 200 KB on a fresh connection, RTT 40 ms, ignoring loss:
  RTT 1: 14 KB    RTT 2: 28 KB    RTT 3: 56 KB    RTT 4: 112 KB   → ~4 RTTs ≈ 160 ms
On a warmed-up, reused connection: ~1 RTT ≈ 40 ms
```

**5. Bandwidth-delay product (BDP).** To fill a pipe you need `bandwidth × RTT`
bytes in flight. 1 Gb/s × 100 ms = **12.5 MB**. If socket buffers or windows are
smaller than that, a single connection *can't* use the bandwidth you're paying
for. That's why cross-region bulk copies use parallel streams.

**6. Closing.** A graceful close is FIN → ACK → FIN → ACK. The side that closes
**first** enters **TIME_WAIT** for 2×MSL (60 s on Linux) to absorb stray
packets. A busy client that opens and closes thousands of short connections
can run out of **ephemeral ports** (Linux default range 32768–60999, about 28k)
because of TIME_WAIT. **RST** is an abort: "this connection doesn't exist" or
"I'm giving up". It's what you see as `ECONNRESET` / "connection reset by peer".

**7. Ports and sockets in the kernel.** A server `listen()`s on `:8080`. The
kernel completes the handshake *without involving your app* and queues the
connection in the **accept backlog**. Your app calls `accept()` to get it. If
your app is too slow to accept, the backlog fills and new SYNs get dropped.
To clients that looks like a 1–3 s connect delay.

**8. L4 load balancers vs L7 proxies (a preview of Part 15).** An **L4 load
balancer** (AWS NLB, IPVS, Maglev) routes on the 5-tuple and passes the TCP
stream through without reading it. An **L7 proxy** (ALB, nginx, Envoy) *ends*
your TCP connection, reads the HTTP, and starts a new connection to the backend.

### Try it yourself

```bash
# All TCP connections with state, process, and per-connection internals (Linux)
ss -tanp
ss -tni dst 198.51.100.20         # shows rtt, cwnd, retrans, mss per connection

# macOS equivalent (less detail)
netstat -anv -p tcp | head -30
lsof -nP -iTCP -sTCP:ESTABLISHED

# Capture just the handshake and teardown packets
sudo tcpdump -i en0 -n 'host shop.example.com and (tcp[tcpflags] & (tcp-syn|tcp-fin|tcp-rst) != 0)'

# Is a port reachable? (L4 check, no HTTP involved)
nc -vz shop.example.com 443

# System-wide retransmission and drop counters
netstat -s | grep -iE 'retrans|listen|overflow|reset'

# Ephemeral port range and TIME_WAIT count (Linux)
sysctl net.ipv4.ip_local_port_range
ss -tan state time-wait | wc -l
```

### What breaks

| Symptom | L4 cause |
|---|---|
| `Connection refused` (immediately) | Host reachable, **nothing listening** on that port → kernel sends RST |
| Connect **hangs then times out** | Firewall/security group silently dropping SYNs, wrong route, or host down |
| Random 1 s or 3 s connect delays | SYN dropped (backlog overflow, conntrack full, packet loss) → SYN retransmit after 1 s |
| `connection reset by peer` mid-request | Peer process crashed, idle timeout reaped the connection, or a middlebox injected an RST |
| `EADDRNOTAVAIL` / "cannot assign requested address" | Ephemeral port exhaustion (no connection pooling, TIME_WAIT pile-up, NAT gateway port limits) |
| Single stream won't exceed ~50 Mb/s cross-region | BDP > window/buffer size |
| Requests take exactly +40 ms (or +200 ms) | **Nagle's algorithm + delayed ACK** interacting. Set `TCP_NODELAY` for request/response protocols |
| `nf_conntrack: table full, dropping packet` in kernel log | Linux connection-tracking table exhausted on a busy node/NAT box |

### War story: "Black Friday, outbound payments start failing"

The checkout service called a payment provider for every order. It opened a
**new HTTPS connection per request** and closed it immediately. All outbound
traffic went through one cloud NAT gateway, which supports about 55,000
simultaneous connections *per unique destination*. At peak traffic, TIME_WAIT
plus in-flight connections exhausted that limit. New connections failed with
timeouts and `EADDRNOTAVAIL`, and payments failed. **Fix:** an HTTP client with
a connection pool and keep-alive (100 long-lived connections instead of tens of
thousands of short ones), plus more NAT IPs. **Lesson:** connection reuse isn't
a micro-optimisation. It's a *capacity* property of your system. Review every
outbound HTTP client for pooling.

### Why it matters to you

- **Every new TCP connection costs 1 RTT, plus slow start.** Every new TLS
  connection costs another RTT. Connection pooling and keep-alive are the
  cheapest latency wins available.
- **Packet loss → retransmission timeouts → tail latency.** Your p99 is often a
  network story.
- **"Refused" and "timed out" point to different teams.** Refused means the host
  answered, so the service is down. Timeout means something dropped the packet,
  so look at firewalls, routing, and security groups.
- **Timeouts must be layered deliberately** (connect timeout < request timeout <
  upstream LB idle timeout). See the war story in Part 7.
- **Ports are a finite resource** on clients, NAT gateways, and conntrack tables.
  High-fan-out services need pooling *by design*.

### Go deeper in the series

Each idea in this part has a Go lab in the [TCP/IP guide](tcp-ip/real-life-guide-v1.md):

| Idea here | Lab there |
|---|---|
| Refused vs timed out | `dialcheck`, [Ch 21](tcp-ip/real-life-guide-v1.md#chapter-21-tcp-part-1-the-three-way-handshake) |
| TCP is a byte stream | `stream`, [Ch 22](tcp-ip/real-life-guide-v1.md#chapter-22-tcp-part-2-how-it-never-loses-your-data) |
| TIME_WAIT / CLOSE_WAIT | `closewait`, [Ch 23](tcp-ip/real-life-guide-v1.md#chapter-23-tcp-part-3-closing-a-connection-and-the-states) |
| Flow control | `zerowindow`, [Ch 24](tcp-ip/real-life-guide-v1.md#chapter-24-flow-control-the-sliding-window) |
| Nagle + delayed ACK (the +40 ms) | `nagle`, [Ch 26](tcp-ip/real-life-guide-v1.md#chapter-26-the-40-millisecond-mystery) (measured: 44 ms) |
| Accept queue full ("slow connects, idle app") | `acceptq`, [Ch 52](tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack) |
| Bandwidth-delay product | `bdp`, [Ch 59](tcp-ip/real-life-guide-v1.md#chapter-59-network-performance-engineering-qos-packet-loss-shaping-saturation) |

And the operating-system side, where 10,000 waiting connections need about
16 threads in Go: [Linux guide, Chapters 73 and 78](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime).


---

## Part 7 — Layer 5: Session

### Job

Establish, maintain, synchronise, resume, and end a **dialogue** between two
applications: something that lasts longer than a single packet or single
request.

### The honest truth

TCP/IP has **no dedicated session layer protocol**. Session responsibilities
are spread across TLS, HTTP, and your application. That's exactly why session
problems are so often *nobody's* problem until they become an incident.

### In our story

Five different "sessions" are active at once when you click Place order:

| "Session" | Lifetime | What establishes it | What it buys you |
|---|---|---|---|
| **TCP connection** | Seconds to minutes | 3-way handshake | Reliable byte stream (strictly L4, but the browser treats it as a reusable session) |
| **TLS session** | Hours (via tickets) | TLS handshake; resumption with **session tickets / PSK** | Skip the expensive handshake work next time |
| **HTTP/2 connection + streams** | Minutes | ALPN negotiation (`h2`), then SETTINGS frames | Many concurrent requests (*streams*) multiplexed over one TCP+TLS connection |
| **Browser connection pool** | Until idle timeout | Browser policy | Your click reuses the connection opened when the cart page loaded, so *no* new TCP or TLS handshake is needed |
| **Application login session** | Days | Login → `Set-Cookie: session=…` or a JWT | The server knows *who you are* on every stateless HTTP request |

So in the *common* case, your click skips the TCP and TLS handshakes entirely:
the browser already has a warm HTTP/2 connection to `shop.example.com` from
loading the checkout page. The POST becomes a new **stream** (stream ID 15, say)
on that existing connection.

Two more session-level behaviours happen on the server side:
- The ALB keeps its own **pool of keep-alive connections** to the pods, which
  is a separate session layer behind the proxy.
- If the ALB has **sticky sessions** enabled, a cookie (`AWSALB=…`) pins you to
  one pod.

### How it works

**1. Session resumption (TLS).** After a full handshake the server gives the
client an encrypted **session ticket**. On reconnecting, the client presents it
and both sides derive keys without repeating the certificate exchange. TLS 1.3
even allows **0-RTT** data, where the request goes in the very first flight,
but 0-RTT data **can be replayed by an attacker**, so it's only safe for
idempotent requests (GET, never POST /orders).

**2. Multiplexing (HTTP/2).** One connection carries many independent
**streams**, each with its own ID and flow control. That removes HTTP/1.1's
"one request at a time per connection" limit. The catch is **TCP-level
head-of-line blocking**: one lost TCP packet stalls *all* streams until it's
retransmitted. HTTP/3 (QUIC over UDP) fixes this by doing loss recovery per
stream.

**3. Dialogue control and checkpointing.** These are classic OSI session
functions. Modern equivalents:
- **Resumable uploads** (tus protocol, S3 multipart): restart from the last
  checkpoint after a disconnect.
- **WebSocket / gRPC streams**: long-lived, full-duplex dialogues with their own
  keepalive pings.
- **Database sessions**: Postgres connection state (transactions, `SET` variables,
  prepared statements) is session state. That's why *transaction-mode* poolers
  like PgBouncer break `SET` and session-level prepared statements.

**4. Application sessions: stateful vs stateless.**

| Approach | Where state lives | Trade-off |
|---|---|---|
| Server-side session (cookie holds an ID → Redis) | Shared store | Easy revocation; extra lookup per request; the store is a dependency |
| Sticky sessions (cookie → specific instance) | Instance memory | Simple; **breaks on deploys/scale-in, uneven load** |
| Token (JWT holds the claims) | The client | No lookup; **hard to revoke**; token size added to every request |

### Try it yourself

```bash
# See ALPN (h2 vs http/1.1), TLS version, and whether the TLS session was reused
curl -sv -o /dev/null https://shop.example.com/ 2>&1 | grep -E 'ALPN|SSL connection|HTTP/'

# Two requests in one curl invocation: watch the connection get reused
curl -sv -o /dev/null -o /dev/null https://shop.example.com/ https://shop.example.com/health 2>&1 \
  | grep -E 'Connected to|Re-using|left intact'

# TLS session resumption test (look for "Reused" on later connections)
openssl s_client -connect shop.example.com:443 -servername shop.example.com -reconnect </dev/null 2>/dev/null \
  | grep -E '^(New|Reused)'

# Browser: chrome://net-internals/#sockets  and  DevTools → Network → "Connection ID" column
```

### What breaks

| Symptom | Session-level cause |
|---|---|
| Sporadic **502s** from the load balancer, especially after idle periods | **Keep-alive timeout mismatch**: the backend closes an idle connection the LB still thinks is open (see war story) |
| Users logged out during a deploy | Sessions stored in instance memory; sticky sessions pointed at terminated instances |
| One pod at 90% CPU, others idle | Sticky sessions or long-lived HTTP/2 / gRPC connections pinning traffic |
| WebSocket drops every ~60 s | An idle timeout somewhere on the path (LB, NAT, corporate proxy) is shorter than your ping interval |
| "Prepared statement does not exist" errors | Transaction-mode connection pooler + session-level DB features |

### War story: "The 502s nobody could reproduce"

A Node.js service behind an ALB returned about 0.05% 502s with no errors in its
own logs. Cause: the ALB's idle timeout was **60 s**, but Node's
`server.keepAliveTimeout` defaulted to **5 s**. After 5 s of idleness, Node
closed the backend connection. Sometimes the ALB had *just* sent a new request
on that same connection at that exact moment. The request hit a closing socket,
and the ALB returned 502. **Fix:** make the *server's* keep-alive timeout
**longer** than the load balancer's idle timeout (for example 65 s vs 60 s), so
the LB always closes first. **Lesson (and a good interview question):** in any
proxy chain, *the side that sends requests on a reused connection must be the
one that closes idle connections first*, so every hop's idle timeout must be
shorter than the timeout of the hop behind it.

### Why it matters to you

- **Connection reuse is the biggest session-layer performance lever.** A warm
  connection avoids 2 RTTs (TCP + TLS) and TCP slow start.
- **Timeouts along a proxy chain are a design artefact.** Document them per hop
  (client → CDN → LB → sidecar → app → DB).
- **Prefer stateless app servers with an external session store or tokens.**
  Sticky sessions are technical debt that shows up at deploy time.
- **Long-lived connections defeat load balancing.** gRPC/HTTP/2 clients on an
  L4 balancer send *everything* to one backend. See Part 15.

### Go deeper in the series

- **The keep-alive 502 race and timeout nesting:** [TCP/IP guide, Chapter 57](tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries).
- **TLS session resumption, tested from Go:** [TCP/IP Chapter 29](tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted)
  (it resumes only after the client has *read* the session ticket).
- **Long-lived sessions (SSE, WebSockets)** and the timeouts that cut them:
  [HTTPS guide, Chapter 24](../v2-https/real-life-guide-v1.md#chapter-24-streaming-over-http-server-sent-events-websockets-and-long-lived-connections).


---

## Part 8 — Layer 6: Presentation

### Job

Make sure both sides agree on **what the bytes mean**: character encoding, data
serialisation, compression, and **encryption**.

### In our story

Four presentation transformations happen to your order before it hits the wire:

```
JS object { cartId: "c_8812", paymentMethodId: "pm_77" }
   │  1. SERIALISE   JSON.stringify → text
   ▼
'{"cartId":"c_8812","paymentMethodId":"pm_77"}'
   │  2. ENCODE      text → UTF-8 bytes (43 bytes)
   ▼
7b 22 63 61 72 74 49 64 22 3a 22 63 5f 38 38 31 32 22 ...
   │  3. COMPRESS    headers compressed with HPACK (HTTP/2); body too small to gzip
   ▼
   │  4. ENCRYPT     TLS 1.3 record: AES-128-GCM with this session's keys
   ▼
17 03 03 00 5a  <90 bytes of ciphertext + 16-byte authentication tag>
```

On the server, the ALB **decrypts** (TLS terminates there). It then either
re-encrypts towards the pod (end-to-end TLS) or sends plaintext inside the VPC.
Finally the pod's framework **decompresses, decodes, and deserialises** the
bytes back into an object.

### How it works

**1. The TLS 1.3 handshake** (when there's no warm connection to reuse):

```
Client                                                       Server (ALB)
  │── ClientHello ───────────────────────────────────────────────▶│
  │     supported versions: TLS 1.3                               │
  │     cipher suites: TLS_AES_128_GCM_SHA256, CHACHA20_POLY1305  │
  │     key_share: X25519 public key        ← key exchange starts NOW
  │     SNI: shop.example.com               ← which certificate do I want (plaintext!)
  │     ALPN: h2, http/1.1                  ← which L7 protocol next
  │                                                               │
  │◀── ServerHello (key_share) ───────────────────────────────────│
  │◀── {EncryptedExtensions, Certificate, CertificateVerify, ──── │   ← already encrypted
  │     Finished}                                                 │
  │                                                               │
  │── {Finished} ── {HTTP request} ──────────────────────────────▶│   1 RTT total
```

The four things TLS gives you:

| Property | Mechanism |
|---|---|
| **Confidentiality** | Symmetric encryption (AES-GCM / ChaCha20) with keys from an ephemeral **ECDHE** exchange (X25519) |
| **Integrity** | AEAD authentication tag on every record. Any tampering → connection fails |
| **Authentication** | Certificate chain: `shop.example.com` → intermediate CA → a root CA in your OS/browser trust store; server proves it owns the private key (CertificateVerify) |
| **Forward secrecy** | Ephemeral keys: stealing the server's private key later doesn't decrypt recorded past traffic |

TLS 1.2 needed **2 RTTs**. TLS 1.3 needs **1**. QUIC combines the transport and
TLS handshakes into **1** as well.

**What's still visible on the wire:** the destination IP, the port, **SNI**
(the hostname, unless Encrypted Client Hello is used), packet sizes, and
timing. Encryption hides *content*, not *metadata*.

**2. Serialisation formats.**

| Format | Strengths | Gotchas |
|---|---|---|
| JSON | Human-readable, universal | No integer type: **JS loses precision above 2⁵³−1**, so send 64-bit IDs as strings. No binary, no dates (agree on ISO-8601), verbose |
| Protobuf / gRPC | Compact, typed, schema evolution | Needs schema discipline: never reuse field numbers; not human-readable |
| Avro / Parquet | Schema registry, analytics | Writer and reader schema compatibility rules |
| XML / SOAP | Mature validation | Heavy; XXE attacks if parsers are misconfigured |

**3. Character encoding.** Always UTF-8 on the wire. Declare it
(`Content-Type: application/json; charset=utf-8`). Mojibake (`CafÃ©` instead of
`Café`) means one side decoded UTF-8 bytes as Latin-1. Database column
encodings, CSV exports, and email are where this still happens.

**4. Compression.** `Accept-Encoding: gzip, br` → `Content-Encoding: br`.
Brotli/gzip typically shrink JSON 70–90%. HTTP/2 compresses *headers* with
**HPACK** (HTTP/3 uses QPACK). Security note: compressing secrets together with
attacker-controlled input under encryption enables CRIME/BREACH-style attacks.

### Try it yourself

```bash
# Full TLS handshake details: protocol, cipher, certificate chain, expiry
openssl s_client -connect shop.example.com:443 -servername shop.example.com </dev/null 2>/dev/null \
  | grep -E 'Protocol|Cipher|subject=|issuer=|Verify return'

# Certificate dates and names (expiry monitoring in one line)
echo | openssl s_client -connect shop.example.com:443 -servername shop.example.com 2>/dev/null \
  | openssl x509 -noout -dates -subject -ext subjectAltName

# Did the server compress the response?
curl -s -o /dev/null -H 'Accept-Encoding: br, gzip' -w '%{size_download} bytes\n' https://shop.example.com/
curl -sI -H 'Accept-Encoding: br, gzip' https://shop.example.com/ | grep -i content-encoding

# Inspect bytes / encoding of a payload
printf 'Café' | xxd        # c3 a9 = é in UTF-8
```

### What breaks

| Symptom | Presentation cause |
|---|---|
| Sudden total outage, browsers show `NET::ERR_CERT_DATE_INVALID` | **Expired certificate** |
| Works in Chrome, fails in `curl`/Java/Android | **Missing intermediate certificate** (browsers sometimes fetch it themselves; other clients don't) |
| Some clients fail TLS, others fine | Clock skew on the client; old clients without TLS 1.2+/SNI; cipher mismatch |
| Wrong certificate served | SNI not sent (old client, raw IP), or the wrong default cert on the LB |
| Order IDs off by one in the UI | 64-bit integer in JSON parsed into a JS `Number` |
| `Café` shows as `CafÃ©` | UTF-8 decoded as Latin-1 somewhere in the pipeline |
| gRPC consumer crashes after a deploy | Protobuf field number reused or type changed |

### War story: "The certificate that expired on a Saturday"

A company's API certificate was renewed manually every year by one engineer,
who left. On a Saturday the certificate expired. Every mobile app version
pinned to that certificate's public key also broke, and fixing that needed an
app-store release that took four days. **Fixes:** automated issuance and
renewal (ACME / cloud certificate manager), expiry monitoring alerting 30 days
ahead, and pinning the *CA/intermediate* key rather than the leaf, with a backup
pin. **Lesson for managers:** certificate expiry is one of the most common
*preventable* outages. Track it as an owned, monitored asset, not tribal
knowledge.

### Why it matters to you

- **Decide where TLS terminates, and say so.** At the CDN? The LB? The pod
  (end-to-end / mTLS via a service mesh)? Every terminating hop can read
  plaintext. That's a compliance (PCI/HIPAA) and threat-model decision.
- **Encoding and serialisation are API contracts.** Put them in your API
  guidelines: UTF-8, ISO-8601 timestamps in UTC, string IDs, explicit nulls vs
  missing fields, protobuf evolution rules.
- **Automate certificates end to end.** A manual certificate is an incident on
  a timer.

### Go deeper in the series

- **TLS 1.3 in depth:** [HTTPS guide, Chapter 4](../v2-https/real-life-guide-v1.md#chapter-4-tls-https-the-encrypted-tunnel).
- **Inspect any server's TLS from Go**, with named failure modes: `tlsinspect`,
  [TCP/IP Chapter 29](tcp-ip/real-life-guide-v1.md#chapter-29-tls-how-it-gets-encrypted).
- **Certificate operations** (inventory, chains, ACME, 47-day lifetimes):
  [TCP/IP Chapter 56](tcp-ip/real-life-guide-v1.md#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation) and [HTTPS Chapter 16](../v2-https/real-life-guide-v1.md#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation).
- **Serialisation and framing in Go:** [Go guide §23 (JSON)](../Golang/real-life-golang-guide.md#23-encoding-json-and-reflection-briefly) and
  [§34 (framing your own protocol)](../Golang/real-life-golang-guide.md#34-networking-deep-dive-net-conn-tcp-udp-framing-your-own-protocol).
- **The cryptography underneath** (hashing vs encryption, AEAD, certificates, running a CA): [Security from Zero](../security/real-life-guide.md); envelope encryption with per-tenant keys, as a Go lab: [Security Engineering in Depth, Chapter 9](../security/real-life-security-guide-v1.md#chapter-9-data-security-in-the-cloud-kms-and-envelope-encryption).


---

## Part 9 — Layer 7: Application

### Job

Provide the service the user actually wanted, using a protocol that both
applications understand: HTTP to place an order, DNS to find a server, the
Postgres wire protocol to save it.

### In our story

Three L7 protocols take part in your one click: **DNS**, **HTTP**, and the
**Postgres protocol**.

**Step A: DNS (assuming nothing is cached).** The browser needs an IP for
`shop.example.com`:

```
Browser cache?        miss
OS resolver cache?    miss
 │
 ▼ UDP 53 → home router 192.168.1.1 (DNS forwarder) → ISP recursive resolver
                                                         │
     ┌───────────────────────────────────────────────────┘
     │ 1. Root server       "who handles .com?"           → a.gtld-servers.net ...
     │ 2. .com TLD server   "who handles example.com?"    → ns1.example-dns.net ...
     │ 3. Authoritative     "shop.example.com?"
     │       → CNAME shop-alb-123.region.elb.example.net
     │       → A     198.51.100.20   TTL 60
     ▼
Browser now has 198.51.100.20 (cached for 60 s)
```

DNS is an L7 protocol carried over **UDP** (L4), falling back to TCP for large
responses. Increasingly it's encrypted (DNS-over-HTTPS/TLS). A cold lookup can
take 20–100 ms or more. A cached one takes 0 ms.

**Step B: the HTTP request.** Inside the TLS tunnel, on HTTP/2 stream 15:

```http
:method: POST
:scheme: https
:authority: shop.example.com
:path: /api/v1/orders
content-type: application/json
content-length: 43
cookie: session=s%3A9d8f...; AWSALB=...
idempotency-key: 7f3c9a2e-1b4d-4c8e-9f00-2a6b8d1e5c43
user-agent: Mozilla/5.0 ...
accept-encoding: gzip, br

{"cartId":"c_8812","paymentMethodId":"pm_77"}
```

**Step C: L7 routing at the load balancer.** Because the ALB decrypts the
traffic and parses HTTP, it can make **content-aware** decisions that an L4
device can't:

```
Rule: Host = shop.example.com AND Path = /api/v1/orders*  → target group "orders-service"
Adds: X-Forwarded-For: 203.0.113.7      (your public IP, which is otherwise lost)
      X-Forwarded-Proto: https
      X-Amzn-Trace-Id: Root=1-6703...   (distributed-tracing correlation ID)
Picks a healthy pod: 10.0.2.77:8080 (least outstanding requests)
Sends over an existing keep-alive connection, HTTP/1.1 (or HTTP/2, gRPC)
```

**Step D: your code.** The `orders-service` handler:

```
1. Authenticate: look up the session cookie (Redis) → user u_4411
2. Idempotency: have we seen key 7f3c9a2e...? No → record it
3. Business logic: price the cart, call the payment provider (another full
   outbound L1–L7 journey through the NAT gateway!)
4. Persist: INSERT INTO orders ... over a pooled Postgres connection
   (Postgres wire protocol, L7, over TCP 5432, often TLS)
5. Respond: 201 Created  {"orderId":"o_99120371","status":"confirmed"}
```

Every outbound call in step D repeats the *entire* OSI journey on a different
path. A single user click usually fans out into 5–50 internal L7 conversations.

### How it works

**1. HTTP's model:** a stateless request/response with method, target, headers,
and body. Stateless means every request carries everything the server needs
(cookies, tokens), which is what lets any instance handle any request and makes
horizontal scaling possible.

**2. Methods and idempotency matter to the network layers below.**

| Method | Safe? | Idempotent? | Can a proxy/client retry automatically? |
|---|---|---|---|
| GET, HEAD | Yes | Yes | Yes |
| PUT, DELETE | No | Yes | Yes (with care) |
| **POST** | No | **No** | **No**, unless you add an **Idempotency-Key** |

When a TCP connection drops *after* the request was sent but *before* the
response arrived, the client can't tell whether the order was created. L4 can't
answer that question. **Only L7 idempotency can.** That's why our request carries
an `Idempotency-Key`.

**3. Status codes: who generated them?** In a proxy chain, the same code means
different things depending on *which hop* produced it:

| Code | From your app | From the load balancer / proxy |
|---|---|---|
| 400/401/403/404 | App logic | Rare (WAF block → 403) |
| 500 | Your code threw an exception | — |
| **502 Bad Gateway** | — | The upstream sent an invalid response or **closed the connection** (crash, keep-alive race, protocol mismatch) |
| **503 Service Unavailable** | Deliberate load shedding | **No healthy targets** (all failing health checks) |
| **504 Gateway Timeout** | — | Upstream **didn't respond** in time (slow DB, deadlock, network drop) |

> **Incident tip:** always check whether a 5xx came from the *app's* access log
> or only the *LB's* log. A 502/504 that appears only in LB logs means the
> request never got a response from your code. It's a network, timeout, or crash
> problem, not a logic bug.

**4. The L7 infrastructure you already run:**

| Component | L7 superpower |
|---|---|
| CDN | Caching by URL/headers, edge TLS termination, close to users |
| WAF | Blocks malicious payloads (SQLi, XSS) by reading the HTTP content |
| API gateway | Auth, rate limiting per API key, request transformation |
| L7 load balancer / ingress | Routing by host/path/header, retries, health checks |
| Service mesh sidecar (Envoy) | mTLS, per-request load balancing, retries, circuit breaking, telemetry |

**5. HTTP versions.**

| | HTTP/1.1 | HTTP/2 | HTTP/3 |
|---|---|---|---|
| Transport | TCP | TCP | **QUIC over UDP** |
| Concurrency | One request per connection at a time (browsers open about 6 connections) | Many streams on one connection | Many streams, **no TCP head-of-line blocking** |
| Headers | Plain text | HPACK compressed | QPACK compressed |
| Handshake (new connection) | TCP + TLS = 2 RTT | TCP + TLS = 2 RTT | **1 RTT** (0-RTT on resumption) |
| Connection migration (Wi-Fi → 4G) | No (new 4-tuple = new connection) | No | **Yes** (connection IDs) |

### Try it yourself

```bash
# DNS: the answer, the full delegation chain, and the TTL
dig shop.example.com +noall +answer
dig +trace shop.example.com
dig @1.1.1.1 shop.example.com A +short    # ask a specific resolver
scutil --dns | head -20                   # macOS: which resolver am I using?

# HTTP: full request/response with headers
curl -v https://shop.example.com/api/v1/health

# Force a specific IP (bypass DNS) while keeping SNI/Host correct: tests one backend/LB
curl -v --resolve shop.example.com:443:198.51.100.20 https://shop.example.com/api/v1/health

# HTTP/3, if your curl supports it
curl --http3 -sI https://shop.example.com/

# Browser: DevTools → Network → click a request → Timing tab
```

### What breaks

| Symptom | L7 cause |
|---|---|
| "Works for me" but not for users after a DNS change | DNS TTL/caching. Old records live until they expire (and some resolvers ignore TTLs) |
| `NXDOMAIN` / `SERVFAIL` | Missing record, broken delegation, DNSSEC validation failure |
| Duplicate orders | Client or proxy retried a non-idempotent POST after a timeout |
| 503 from LB, app looks healthy | Health-check path misconfigured / returning non-200 / wrong port |
| 413 Payload Too Large / 431 headers too large | Proxy limits (body size, header size: watch for huge cookies/JWTs) |
| Real client IP is the LB's IP in logs | Not reading `X-Forwarded-For` (or reading it insecurely, trusting spoofed values) |
| CORS errors in the browser | Missing/incorrect `Access-Control-Allow-*` headers. A *browser* policy at L7 |

### War story: "The retry storm"

A slow database caused `orders-service` latency to rise from 100 ms to 3 s. The
ALB timed out at 2 s and returned 504. The mobile app retried 3 times. The API
gateway *also* retried 2 times. Each hop's retries multiplied: one user tap
became up to 3 × 3 = 9 requests to an already overloaded database, which then
fell over completely. **Fixes:** retry at only one layer, use exponential
backoff with jitter, set retry budgets, require idempotency keys for POSTs, and
add a circuit breaker. **Lesson:** retries are an L7 *policy* that multiplies
load across the whole stack. Somebody has to own the end-to-end retry design.

### Why it matters to you

- **L7 is where business semantics live**: idempotency, retries, auth, rate
  limits. Lower layers can't make those decisions for you.
- **DNS is part of your deployment and failover system.** TTLs decide how fast
  you can move traffic. Lower them *before* a migration.
- **Know which hop generates each 5xx.** It cuts incident time in half.
- **Every L7 proxy terminates the connection.** Client IP, TLS, and timeouts all
  restart at each hop.

### Go deeper in the series

- **DNS from raw bytes** (`dnsquery`) and **DNS operations** (`dnsdiff`):
  [TCP/IP Chapters 18 and 55](tcp-ip/real-life-guide-v1.md#chapter-18-dns-turning-names-into-addresses).
- **HTTP versions and how HTTP/3 is discovered** (`protoprobe`): [HTTPS Chapter 14](../v2-https/real-life-guide-v1.md#chapter-14-http-version-negotiation-http-1-1-http-2-http-3-alt-svc-and-fallbacks).
- **Caching done right** (`edgecache`, including a cache leaking one user's
  data to another): [HTTPS Chapter 15](../v2-https/real-life-guide-v1.md#chapter-15-cache-correctness-browser-cdn-proxy-and-application-caches).
- **Idempotency keys and retries, measured** (`resilience`): [HTTPS Chapter 21](../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers).
- **502/504 from a real reverse proxy** (`revproxy`): [HTTPS Chapter 20](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling).
- **Attacking L7:** request smuggling, cache poisoning, SSRF, XXE, races, and
  uploads, each with a lab: [Security Engineering in Depth, Part 6](../security/real-life-security-guide-v1.md#part-6-advanced-web-and-api-exploitation).


---

## Part 10 — The full replay: a millisecond timeline

Two versions of the same click: **cold** (first visit, nothing cached) and
**warm** (what usually happens: the page you're on already has a connection).

Assumptions: Wi-Fi + fibre, RTT to the region = **40 ms**, server work = **90 ms**.

### Cold path

```
 t (ms)  Layer   Event
 ──────  ─────   ─────────────────────────────────────────────────────────────────
      0  L7      fetch() called. Browser checks DNS cache: miss
      0  L4/L3   DNS query over UDP → router → ISP resolver
                   (each packet: L3 IP → L2 Wi-Fi frame → L1 radio, then Ethernet, then light)
     25  L7      DNS answer: 198.51.100.20 (resolver had .com cached; went to authoritative)
     25  L4      TCP SYN →
     65  L4      ← SYN-ACK. ACK →. Connection ESTABLISHED          (1 RTT)
     65  L5/L6   TLS ClientHello (SNI, ALPN h2, key share) →
    105  L5/L6   ← ServerHello + cert + Finished. Keys ready       (1 RTT)
    105  L5      HTTP/2 connection preface + SETTINGS
    105  L7      POST /api/v1/orders (HEADERS + DATA frames, encrypted) →
    125          request arrives at ALB (½ RTT)
    126  L7      ALB: decrypt, parse, route, add X-Forwarded-For, pick pod
    127  L4      ALB → pod over an existing keep-alive connection (~0.5 ms in-VPC)
    128  L7      orders-service: auth (Redis 1 ms), idempotency check, payment call (60 ms),
                 INSERT into Postgres (5 ms), build response
    215  L7      201 Created → ALB → encrypt → TCP → IP → ... →
    235  L7      Browser receives response (½ RTT), decrypts, parses JSON
    240  L7      UI shows "Order confirmed"
```

### Warm path (connection reused)

```
      0  L7      fetch() → reuse the existing HTTP/2 connection, new stream ID
      0  L7      POST sent immediately
     20          arrives at ALB (½ RTT)
    110          server done (90 ms)
    130  L7      response arrives. UI updates.
```

### Where the time goes

```
COLD (240 ms)                                   WARM (130 ms)
DNS        ████                     25 ms
TCP        ███████                  40 ms
TLS        ███████                  40 ms
Network    ███████                  40 ms        ███████        40 ms
Server     ███████████████          90 ms        ███████████████ 90 ms
Client       █                       5 ms
```

**What a senior engineer reads from this:**

- **About 45% of the cold request is setup** (DNS + TCP + TLS) before your code
  even runs. Connection reuse, DNS prefetching (`<link rel="preconnect">`), and
  HTTP/3 attack exactly this.
- **Network RTT appears at least once in every request**, and once per
  *sequential* round trip. A page that makes 5 dependent API calls pays 5 RTTs.
  Batch calls, or move the work closer to users.
- **Server time is the only part you fully control in code**, and inside it the
  payment call (another network journey) dominates.

---

## Part 11 — X-ray of one packet

### 11.1 The byte budget of one full-size packet

A full-size segment of your response leaving the ALB on a standard
1500-byte-MTU Ethernet link:

```
 Layer   Header                              Bytes   Running total on the wire
 ─────   ──────────────────────────────────  ─────   ─────────────────────────
  L1     Preamble + SFD                          8
  L2     Ethernet header (dst, src, type)       14
  L3     IPv4 header                            20   ┐
  L4     TCP header + timestamps option          32   │  MTU = 1500
  L6     TLS record header + tag + type (≈22)   (in the payload, amortised over a 16 KB record)
  L7     HTTP/2 frame header (9)                 (in the payload)
         Application payload                  1448   ┘
  L2     FCS                                     4
  L1     Inter-frame gap                        12
                                              ─────
                                              1538 bytes on the wire for 1448 bytes of TCP payload
```

**Efficiency ≈ 94%.** With tiny packets (a 40-byte ACK) the ratio flips: a
64-byte minimum frame plus 20 bytes of L1 overhead carries almost nothing.
That's why **packets per second (pps)**, not just bits per second, is a real
limit on load balancers, firewalls, and NICs.

**Overheads that quietly reduce your MTU:**

| Technology | Overhead | Effective inner MTU |
|---|---|---|
| PPPoE (many DSL/fibre ISPs) | 8 B | 1492 |
| VXLAN overlay (K8s, VMware) | 50 B | 1450 |
| IPsec VPN (tunnel mode) | ~50–73 B | ~1400–1430 |
| WireGuard | 60 B (IPv4) / 80 B (IPv6) | 1420 typical |
| GRE | 24 B | 1476 |

### 11.2 Hop-by-hop: what each device reads and rewrites

Follow the first request packet. ✎ = rewritten; 👁 = read but unchanged; — = not looked at.

| Hop | L2 MACs | L3 Src IP | L3 Dst IP | TTL | L4 ports | L5–L7 (TLS/HTTP) |
|---|---|---|---|---|---|---|
| Laptop sends | laptop → router | 192.168.1.23 | 198.51.100.20 | 64 | 54012 → 443 | encrypted |
| Wi-Fi AP (part of router) | ✎ 802.11 → 802.3 | — | — | — | — | — |
| Home router (route + NAT) | ✎ new frame | ✎ **203.0.113.7** | 👁 | ✎ 63 | ✎ **61001** → 443 | — |
| ISP / Internet routers (×12) | ✎ new frame each hop | 👁 | 👁 | ✎ −1 each | — | — |
| Stateful firewalls along the way | ✎ | 👁 | 👁 | ✎ | 👁 (track state) | — |
| **ALB (L7 proxy): connection ENDS** | ✎ | — | — | — | ends at 443 | **decrypts, reads, routes** |
| **ALB opens a new connection** | new | **10.0.1.15** | **10.0.2.77** | 64 (fresh) | **new** ephemeral → 8080 | new HTTP request (+ X-Forwarded-For) |
| K8s node / CNI | ✎ (veth) | 👁 | 👁 (or DNAT if via a Service IP) | ✎ | 👁 | — |
| Pod | — | — | — | — | — | your handler reads it |

Three insights in that table that come up in architecture reviews:

1. **Routers don't care about anything above L3.** That's what lets the
   Internet carry any protocol.
2. **NAT breaks end-to-end addressing.** The server sees `203.0.113.7`, shared by
   your whole household.
3. **An L7 proxy is a full endpoint.** Two separate TCP connections, two separate
   TLS sessions, and the client's identity survives only in headers you
   deliberately add.

---

## Part 12 — Inside the server: from the NIC to your handler

"Moving up the stack" on the receiving side means a real sequence of hardware
and kernel events. Here's what happens on the Kubernetes node hosting
`orders-service`:

```
 L1  Signal arrives at the (virtual) NIC; PHY decodes symbols → bits
 L2  NIC checks FCS, matches dst MAC, DMA-copies the frame into a RAM ring buffer
     ── NIC raises an interrupt; kernel switches to polling (NAPI) to batch packets
     ── (offloads: the NIC may verify checksums, combine segments (GRO/LRO),
         and spread flows across CPU cores by hashing the 5-tuple (RSS))
 L2→3 Driver builds an sk_buff; EtherType 0x0800 → hand to IPv4
 L3  Validate header; routing decision: "this is for a local pod address"
     ── netfilter/iptables hooks: conntrack lookup, kube-proxy DNAT
         (if addressed to a Service ClusterIP → rewritten to a pod IP)
     ── forwarded across the node's bridge/veth pair into the pod's network namespace
 L4  TCP: find the socket by 5-tuple; check seq numbers; queue in the socket's
     receive buffer; send ACK (possibly delayed)
     ── wakes the process blocked in epoll/kqueue: "socket readable"
 L5-7 User space: the runtime's event loop reads bytes → TLS decrypt (if mTLS) →
     HTTP parser → router → your handler(req)
```

Each step has a queue, and **every queue can overflow and drop**:

| Queue | Overflow symptom | Where to look |
|---|---|---|
| NIC ring buffer | `rx_missed` / `rx_no_buffer` counters | `ethtool -S` |
| Kernel backlog (`netdev_max_backlog`) | drops in `/proc/net/softnet_stat` | softnet stats |
| Conntrack table | "table full, dropping packet" | `dmesg`, `conntrack -S` |
| SYN / accept backlog | `ListenOverflows`, `ListenDrops` | `netstat -s`, `nstat` |
| Socket receive buffer | `RcvbufErrors`, zero-window advertisements | `ss -tm`, `nstat` |
| Application worker pool / event loop | latency rises, then timeouts | APM, event-loop lag metrics |

> **Rule of thumb:** when a service is "slow but CPU is low", look for one of these
> queues being full, or for the app waiting on a downstream call.

### Go deeper in the series

- **The kernel receive path, netfilter hooks, and queue overflows**, with a
  lab that fills an accept queue to 4,097: [TCP/IP guide, Chapter 52](tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack).
- **From the program's side:** threads vs the netpoller
  ([Linux Chapter 73](../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime)), file-descriptor limits
  ([Chapter 78](../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections)), and profiling a live Go service
  ([Chapter 81](../os-linux/real-life-os-guide.md#chapter-81-profiling-go-on-linux-pprof-the-execution-tracer-and-perf)).


---

## Part 13 — The return trip and the goodbye

### The response goes back down and up again

The pod writes `201 Created` + JSON. The data goes **down** the pod's stack (L7 →
L1), crosses the VPC, goes **up** the ALB's stack to L7 (the ALB reads the
response, may compress it, logs it), then **down** again into your TLS
connection. It crosses the Internet (routers handle only L1–L3), and the home
router reverses the NAT using its table (`203.0.113.7:61001` →
`192.168.1.23:54012`). The laptop's stack takes it **up** from L1 to L7, and
JavaScript's `await` resolves.

Two details that surprise people:

- **The return path may differ from the forward path.** BGP routing is
  per-direction. Your traceroute shows only *your* direction.
- **The NAT table entry is the only reason the reply finds you.** If it had
  expired (an idle connection), the reply would be dropped and you'd see a
  timeout.

### Ending the conversation

Nothing closes after your click. **That's on purpose.** The browser keeps the
HTTP/2 connection open for the next click. Eventually:

```
After ~minutes idle:  browser or ALB sends HTTP/2 GOAWAY        (L7/L5 says "no more streams")
                      TLS close_notify alert                    (L6 says "end of encrypted data")
                      TCP FIN → ACK, FIN → ACK                  (L4 graceful close)
                      closer holds TIME_WAIT for 60 s           (L4 hygiene)
                      NAT mapping expires in the home router    (L3 state cleaned up)
                      ARP entry for gateway stays (refreshed)   (L2)
                      Wi-Fi stays associated                    (L1/L2)
```

**Graceful shutdown on deploys** depends on getting this sequence right: stop
accepting new connections, send GOAWAY / `Connection: close`, finish in-flight
requests, *then* exit. Kubernetes adds a twist: removing a pod from load
balancer targets is **asynchronous** with sending it SIGTERM, so a pod must keep
serving for a few seconds after SIGTERM (a `preStop` sleep) or users get 502s
on every deploy.

### Go deeper in the series

- **Graceful shutdown, measured:** a draining backend lost 0 of 300 requests
  and `kill -9` lost 24 ([TCP/IP Chapter 57](tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries)).
- **Shutdown in Kubernetes**, with the `preStop` sleep this section describes:
  [Linux guide, Chapter 82](../os-linux/real-life-os-guide.md#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes).


---

## Part 14 — Debugging by layer: the incident playbook

### 14.1 Three strategies

| Strategy | Start at | Use when |
|---|---|---|
| **Bottom-up** | L1 (link, cables, errors) and work upward | "Nothing works"; new network/hardware; one host is affected |
| **Top-down** | L7 (the error message, logs) and work downward | One app/endpoint is broken while others work |
| **Divide and conquer** | L3/L4 (`ping`, `nc`) | Most of the time: one test splits the stack in half |

### 14.2 The ladder: one command per layer

Run these from the failing client, in order. **The first one that fails tells
you the layer.**

```bash
# L1/L2: do I have a link and an address?
ifconfig en0 | grep -E 'status|inet '        # Linux: ip -br addr
# L2/L3: can I reach my gateway?
ping -c 3 $(route -n get default | awk '/gateway/{print $2}')   # Linux: ip route | awk '/default/{print $3}'
# L3: can I reach the Internet by IP (no DNS)?
ping -c 3 1.1.1.1
# L7 (DNS): can I resolve the name?
dig +short shop.example.com
# L3: can I reach the target IP? (ICMP may be blocked; don't trust "fail" alone)
ping -c 3 198.51.100.20 ; mtr -rwc 20 198.51.100.20
# L4: is the port open?
nc -vz shop.example.com 443
# L5/L6: does TLS negotiate, with a valid certificate?
openssl s_client -connect shop.example.com:443 -servername shop.example.com </dev/null
# L7: does the app answer correctly?
curl -sv https://shop.example.com/api/v1/health
```

### 14.3 One `curl` command that splits latency by layer

This is the most useful single command in this guide:

```bash
curl -o /dev/null -s -w '
  DNS lookup  (L7 DNS)      : %{time_namelookup}s
  TCP connect (L4)          : %{time_connect}s
  TLS done    (L5/L6)       : %{time_appconnect}s
  Request sent              : %{time_pretransfer}s
  First byte  (server work) : %{time_starttransfer}s
  Total                     : %{time_total}s
  HTTP %{http_code}, %{size_download} bytes, remote %{remote_ip}, HTTP/%{http_version}
' https://shop.example.com/api/v1/health
```

The values are **cumulative**. Subtract consecutive values to get each phase:

| Large gap between… | Means | Look at |
|---|---|---|
| 0 → `namelookup` | Slow DNS | Resolver, TTLs, DNS provider |
| `namelookup` → `connect` | High RTT or SYN loss/retries (a gap ≈1 s or ≈3 s = retransmitted SYN) | Network path, firewall, backlog |
| `connect` → `appconnect` | Slow TLS (OCSP, large cert chain, CPU-starved TLS terminator, extra RTT) | LB/TLS config |
| `pretransfer` → `starttransfer` | **Server think time** (the app or its dependencies) | APM traces, DB |
| `starttransfer` → `total` | Large body / low throughput / loss | Payload size, compression, BDP |

### 14.4 Symptom → layer quick map

| What the user/engineer says | Most likely layer(s) | First check |
|---|---|---|
| "No internet at all" | L1–L3 | Link, IP address, gateway ping |
| "Some sites work, others don't" | L3 DNS/L7, MTU | `dig`, PMTU ping test |
| "It's slow, but only sometimes / only p99" | L1/L2 loss → L4 retransmits | `mtr` loss, `ss -ti` retrans, interface errors |
| "Connection refused" | L4 (nothing listening) | Process up? Right port? |
| "Connection timed out" | L3/L4 filtering | Security groups, NACLs, routes |
| "Connection reset" | L4/L5 | Idle timeouts, crashes, keep-alive mismatch |
| "SSL/certificate error" | L6 | Expiry, chain, SNI, client clock |
| "502 / 504" | L4–L7 between proxy and app | LB logs vs app logs, timeouts chain |
| "500 / wrong data" | L7 (or L6 serialisation) | App logs, traces |
| "Works locally, fails in Kubernetes" | L3 (network policy), L7 (DNS/service names), L2 (MTU) | `kubectl exec` + the ladder from inside the pod |

### 14.5 Packet capture: ground truth

When logs disagree, the packets decide:

```bash
# Capture on a server, filtered to one client, full packets, to a file
sudo tcpdump -i any -nn -s0 -w /tmp/case123.pcap host 203.0.113.7 and port 443

# On Kubernetes: capture inside a pod's network namespace
kubectl debug -it pod/orders-7d9f --image=nicolaka/netshoot --target=orders -- tcpdump -nn -i any port 8080
```

Open it in **Wireshark** and use *Statistics → Conversations*, *Analyze → Expert
Information* (retransmissions, zero windows, resets), and *Statistics → TCP
Stream Graphs*. For TLS, export browser keys with `SSLKEYLOGFILE=/tmp/keys.log`
and point Wireshark at that file to see decrypted HTTP.

### Go deeper in the series

- **The ladder as a program:** `netcheck` stops at the first failing layer
  ([TCP/IP Chapter 34](tcp-ip/real-life-guide-v1.md#chapter-34-a-troubleshooting-method-that-works)).
- **`curl -w` in Go:** `httptrace` times each phase and shows connection reuse
  ([HTTPS Chapter 17](../v2-https/real-life-guide-v1.md#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs)).
- **Decode a capture yourself:** `pcapread` ([TCP/IP Chapter 32](tcp-ip/real-life-guide-v1.md#chapter-32-wireshark-reading-a-conversation)).
- **Running the incident:** [TCP/IP Chapter 60](tcp-ip/real-life-guide-v1.md#chapter-60-incident-response-runbooks-packet-evidence-and-escalation) and
  [Linux Chapter 72](../os-linux/real-life-os-guide.md#chapter-72-incident-response-and-refactoring-on-real-linux-systems).
- **Security incidents in cloud and containers**, and detections as tested Go code: [Security Engineering in Depth, Chapters 51 and 55](../security/real-life-security-guide-v1.md#chapter-55-incident-response-for-cloud-and-containers).


---

## Part 15 — Architecture decisions by layer: the senior and manager view

### 15.1 L4 vs L7 load balancing

| | L4 (NLB, IPVS, Maglev) | L7 (ALB, nginx, Envoy, ingress) |
|---|---|---|
| Sees | IPs, ports, TCP/UDP | Full HTTP: host, path, headers, cookies, body |
| Terminates TCP/TLS | No (pass-through; optional TLS termination) | Yes |
| Client IP preserved | Yes (or via PROXY protocol) | Only via `X-Forwarded-For` |
| Latency / throughput | Lowest overhead, millions of pps | More CPU per request |
| Routing | Per **connection** | Per **request** |
| Features | Static IPs, any TCP/UDP protocol | Path routing, retries, auth, WAF, canaries, header rewriting |
| Pick when | Non-HTTP protocols, extreme throughput, static IPs, TLS pass-through for compliance | HTTP APIs and web apps (most of the time) |

**The gRPC trap:** gRPC uses long-lived HTTP/2 connections. Behind an **L4**
balancer, each client opens one connection to one backend and sends *every*
request there, so new pods get no traffic and hot pods melt. You need
**per-request L7 balancing** (Envoy/service mesh, an L7 LB with HTTP/2
support, or client-side load balancing with a lookaside resolver).

### 15.2 Where to terminate TLS

```
Option A: edge only        Client ══TLS══ CDN/LB ──plain── app        simplest; plaintext inside VPC
Option B: re-encrypt       Client ══TLS══ LB ══TLS══ app              compliance-friendly; certs on apps
Option C: pass-through     Client ════════TLS════════ app (L4 LB)     LB can't do L7 routing
Option D: mesh mTLS        Client ══TLS══ LB ══mTLS══ sidecar ─ app   zero-trust; identity per service
```

Ask: who can see plaintext? Where do certificates live, and who rotates them?
What does the compliance scope (PCI, HIPAA) require?

### 15.3 Ownership map: who owns which layer

| Layer | Typically owned by | Typical vendor / product touchpoints |
|---|---|---|
| L1–L2 | Cloud provider / data-centre / IT networking | Instance network performance, Direct Connect / ExpressRoute links |
| L3 | Platform/network team | VPC design, IPAM, route tables, Transit Gateway, VPN, DDoS protection (L3/L4) |
| L4 | Platform team | NLB, security groups, NAT gateways, conntrack/port limits |
| L5–L6 | Platform + security | Certificates, TLS policy, mTLS/mesh, KMS |
| L7 | Product engineering teams (+ platform for shared infra) | API gateway, ALB/ingress rules, WAF rules, CDN config, DNS records |

Gaps in this table cause incidents. **The most common orphaned items:**
certificate renewal, DNS records, timeout chains across proxies, NAT gateway
capacity, and MTU on overlays and VPNs. Assign each one an owner explicitly.

### 15.4 Design checklist by layer

Use this in design reviews for any new service:

- **L1/L2:** Which regions/AZs? What RTT do users get? Any L2 assumptions
  (multicast, VRRP) that won't hold in the cloud? MTU across overlays and VPNs?
- **L3:** Non-overlapping CIDR from central IPAM? Public vs private subnets? Egress
  through NAT, with capacity planned? IPv6 plan? Which ICMP is allowed?
- **L4:** TCP or UDP/QUIC? Connection pooling on every client? Connect/read
  timeouts set explicitly? Expected connections per second and concurrent
  connections vs NAT and conntrack limits?
- **L5:** Stateless servers? Where is session state stored? Idle timeouts
  documented per hop, each hop shorter than the next? Graceful shutdown and
  connection draining?
- **L6:** TLS terminated where, and who can see plaintext? Certificates
  automated and monitored? Serialisation format and schema-evolution rules?
  UTF-8, timestamp, and ID conventions?
- **L7:** Idempotency for writes? Retry policy owned at exactly one layer, with
  backoff, jitter, and a budget? Rate limiting by what key (not just IP)?
  Health checks that test real dependencies, but not too aggressively? DNS
  TTLs for failover?

### 15.5 Cost lives at layers too

- **Data transfer** (L3): cross-AZ, cross-region, and Internet egress are billed
  per GB. A chatty service talking across AZs can cost more than its compute.
- **NAT gateway processing** (L3/L4): billed per GB. Route traffic to cloud
  services through VPC endpoints instead.
- **L7 features** (ALB LCUs, WAF rules, API gateway per-request pricing): grow
  with requests, not just bytes.

### Go deeper in the series

- **L4 vs L7, health checks, draining, the PROXY protocol:** [TCP/IP Chapter 57](tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries);
  the Go plan builds an [L4 load balancer (Day 48)](../Golang/detailed-90-day-plan/week7.md#day-48-tcp-load-balancer-l4).
- **TLS termination, X-Forwarded-For trust, request smuggling:** [HTTPS Chapter 20](../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling).
- **SLOs and error budgets** for the decisions above: [HTTPS Chapter 23](../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services).
- **Kubernetes networking pitfalls** (`ndots:5`, conntrack DNS race, source-IP
  loss): [TCP/IP Chapter 58](tcp-ip/real-life-guide-v1.md#chapter-58-kubernetes-networking-services-ingress-networkpolicy-cni).
- **Packaging the service** (15 MB image, systemd, Kubernetes manifest):
  [Linux Chapter 82](../os-linux/real-life-os-guide.md#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes).
- **The security review of the same architecture:** attack paths and chokepoints ([Security in Depth, Chapter 1](../security/real-life-security-guide-v1.md#chapter-1-from-vulnerabilities-to-attack-paths)) and the senior-engineer review playbook ([Chapter 77](../security/real-life-security-guide-v1.md#chapter-77-security-architecture-review-the-senior-engineer-playbook)).


---

## Part 16 — Cheat sheet

| Layer | PDU | Address | Devices / software | Key protocols | Cloud equivalent | Go-to tools | Classic failures |
|---|---|---|---|---|---|---|---|
| **7 Application** | Message | URL, hostname | Browser, app, ALB, API gateway, CDN, WAF | HTTP, DNS, gRPC, SMTP, Postgres | ALB, API Gateway, CloudFront, Route 53 | `curl -v`, `dig`, DevTools | 5xx, retry storms, DNS caching, duplicate POSTs |
| **6 Presentation** | Message | — | TLS libraries, serialisers | TLS, JSON, protobuf, UTF-8, gzip/brotli | ACM/certificate manager, KMS | `openssl s_client`, `xxd` | Expired/missing-chain certificates, mojibake, precision loss |
| **5 Session** | Message | Session ID, ticket, cookie | TLS, HTTP/2, connection pools | TLS resumption, HTTP/2 streams, WebSocket | Sticky sessions, ElastiCache sessions | `curl -v` (Re-using), `openssl -reconnect` | Keep-alive 502s, sticky imbalance, idle disconnects |
| **4 Transport** | Segment / datagram | Port | Kernel, L4 LB, firewall | TCP, UDP, QUIC | NLB, security groups | `ss -ti`, `nc -vz`, `netstat -s` | Refused vs timeout, port exhaustion, retransmits, backlog drops |
| **3 Network** | Packet | IP address | Router, L3 switch, NAT, firewall | IPv4/6, ICMP, BGP, OSPF | VPC, subnets, route tables, NAT/IGW, TGW, NACL | `ip route`, `traceroute`, `mtr`, `ping` | Overlapping CIDRs, no route, PMTU black-hole, asymmetric routing |
| **2 Data Link** | Frame | MAC address | Switch, Wi-Fi AP, bridge, NIC | Ethernet, 802.11, ARP, VLAN, STP, VXLAN | ENI, (hidden) virtual switch | `arp -a`, `ip neigh`, `tcpdump -e` | Stale ARP, loops/broadcast storms, VLAN mismatch, MTU |
| **1 Physical** | Bit / symbol | — | Cables, optics, radios, NIC PHY, ONT | 1000BASE-T, 802.11ax PHY, GPON, DWDM | Instance network bandwidth | `ethtool`, `wdutil`, `iw`, `networkQuality` | CRC errors, low SNR, duplex mismatch, distance (latency floor) |

**Numbers to remember**

| Number | Meaning |
|---|---|
| ~5 µs/km | Light in fibre. 1,000 km ≈ 10 ms RTT minimum |
| 1500 | Standard Ethernet MTU (bytes). 1460 TCP MSS (1448 with timestamps) |
| 9001 | AWS in-VPC jumbo MTU |
| 64 / 128 | Default TTL (Linux/macOS / Windows) |
| 1 s, 200 ms | Initial TCP RTO (SYN retry) / Linux minimum RTO |
| 10 segments | TCP initial congestion window (~14 KB) |
| 60 s | Linux TIME_WAIT |
| 32768–60999 | Linux ephemeral ports (~28k) |
| 1 RTT | TLS 1.3 handshake (TLS 1.2: 2 RTT; QUIC: 1 RTT including transport) |
| 2⁵³−1 | Largest safe integer in a JavaScript Number |

---

## Part 17 — Self-check questions

<details>
<summary><b>1. Your request crosses 15 routers. Which addresses change, and which stay the same?</b></summary>

MAC addresses change at **every** hop (each link has a new L2 frame). The IP
addresses stay the same end to end, **except** where NAT rewrites them (your
home router rewrites the source IP and port). TTL decrements at every router.
Everything from L4 up is untouched by routers, but it's *terminated* by any L7
proxy (the ALB), which starts a brand-new connection with new IPs and ports.
</details>

<details>
<summary><b>2. A user reports "connection refused" and another reports "connection timed out" for the same service. What does each tell you?</b></summary>

**Refused:** the packet reached the host and the kernel replied with RST
because nothing is listening on that port (process down, wrong port, bound to
127.0.0.1 only). **Timed out:** the SYN got no answer at all. Something dropped
it silently (security group, NACL, firewall, missing route, host down). Different
causes, often different owners.
</details>

<details>
<summary><b>3. Why do small API responses work but large ones hang after enabling a VPN/overlay?</b></summary>

The MTU was reduced by the encapsulation overhead, but endpoints still send
1500-byte packets with DF set. The oversized packets are dropped, and the ICMP
"fragmentation needed" message that would trigger Path MTU Discovery is blocked
or never generated. Fix the inner MTU, clamp TCP MSS, or allow the ICMP.
</details>

<details>
<summary><b>4. Why does 0.5% packet loss hurt p99 latency so much more than p50?</b></summary>

Most requests see no loss (p50 unaffected), but a lost packet costs a
retransmission timeout of at least 200 ms (Linux minimum RTO), and a lost SYN
costs 1 s. Those rare, huge penalties concentrate in the tail. HTTP/2 makes it
worse: one lost TCP packet stalls every multiplexed stream on that connection.
</details>

<details>
<summary><b>5. You're behind an ALB and see intermittent 502s with nothing in your app's logs. What's your first hypothesis?</b></summary>

A keep-alive / idle-timeout race: your server closes idle connections sooner
than the load balancer does, so the LB occasionally sends a request on a socket
the server is closing. Make the server's keep-alive timeout longer than the LB's
idle timeout. Other candidates: pod crashes/OOM kills, or pods terminating
during deploys without connection draining.
</details>

<details>
<summary><b>6. gRPC traffic all goes to one pod even though you have 10. Why, and what's the fix?</b></summary>

gRPC multiplexes all requests over a long-lived HTTP/2 connection. An L4
balancer (or a Kubernetes ClusterIP Service through kube-proxy) balances
*connections*, not *requests*, so each client sticks to one pod. Use L7
per-request balancing (Envoy/service mesh, an HTTP/2-aware L7 LB) or
client-side load balancing.
</details>

<details>
<summary><b>7. Why must a POST /orders carry an idempotency key when TCP already guarantees reliable delivery?</b></summary>

TCP guarantees delivery of *bytes on a live connection*. It can't tell the
client whether the server *processed* a request if the connection died before
the response arrived. Only the application knows whether the order exists. An
idempotency key lets the client safely retry, and the server deduplicates.
That's end-to-end reliability, which by its nature belongs at L7.
</details>

<details>
<summary><b>8. Where does TLS sit in the OSI model, and why is the better question "who terminates it"?</b></summary>

Its handshake acts like L5 (session setup and resumption), and its record
protection acts like L6 (encryption). It runs on top of L4. In practice, the
important facts are *which component decrypts traffic* (and can therefore read,
route, and log it), where the certificates live, and who rotates them. Those
answers drive security, compliance, and routing capabilities.
</details>

<details>
<summary><b>9. Your service rate-limits per client IP. A big customer complains that their whole office is blocked. What happened?</b></summary>

The whole office reaches you through one NAT IP (or a corporate proxy, or
CGNAT), so all their users share one rate-limit bucket. Rate-limit on
authenticated identity (API key, user ID, tenant) and use IP only as a coarse
abuse signal.
</details>

<details>
<summary><b>10. A cold request takes 240 ms, a warm one 130 ms. Name three changes that reduce the cold cost without touching server code.</b></summary>

(1) `preconnect` / DNS prefetch so TCP and TLS happen before the click; (2) HTTP/3
(QUIC) to merge the transport and TLS handshakes, and TLS session resumption /
0-RTT for safe requests; (3) terminate TCP/TLS at a CDN edge close to the user,
so the handshakes cost a short RTT and the edge reuses warm connections to the
origin. Also: longer DNS TTLs where safe, and keeping connections alive longer.
</details>

---

### Where to go next

- **The rest of the series**, in order: this guide is step 2a. Next is
  [Networking from Zero (TCP/IP)](tcp-ip/real-life-guide-v1.md) (L1–L4 in depth), then the security
  pair, [Security from Zero](../security/real-life-guide.md) and [Security Engineering in Depth](../security/real-life-security-guide-v1.md),
  and finally [The HTTPS Request Lifecycle](../v2-https/real-life-guide-v1.md), whose
  [Chapter 25](../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide) traces one request through every guide in Go.
  The machines at both ends are covered in step 1,
  [Operating Systems, Linux, and Containers](../os-linux/real-life-os-guide.md).
- **Hands-on:** capture your own "Place order" click with Wireshark
  (`SSLKEYLOGFILE` enabled), and label every packet with its layer and the
  matching step in Part 10. Doing this once is worth more than reading any
  guide twice.
- **References:** RFC 9293 (TCP), RFC 8446 (TLS 1.3), RFC 9000 (QUIC), RFC 9113
  (HTTP/2), RFC 9114 (HTTP/3), RFC 1918 (private IPv4), RFC 6598 (CGNAT space).

---

## Continue the series

**Next, step 2b (Networking: in depth): [Networking from Zero (TCP/IP)](tcp-ip/real-life-guide-v1.md).** It covers how machines talk: addressing, routing, TCP, TLS, packet capture, network operations. Hands-on: 23 Go labs ([§0.8](tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself)).
