# TCP/IP — The Complete Field Guide (Beginner → Expert)

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/tcp-ip-reference/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> A practical, example-driven reference for the TCP/IP protocol suite: how packets
> are built, addressed, routed, delivered reliably, and debugged in the real world.
>
> Every concept is paired with a concrete scenario, real command output, or a
> packet-level walkthrough. Read it top to bottom once; keep it as a lookup guide
> forever.

---

## How to use this guide

- **Beginner (Parts I–IV):** mental models, the 4-layer stack, addressing, and
  what actually happens when you load a web page.
- **Intermediate (Parts V–IX):** IP addressing math, routing, DNS, DHCP, NAT, and
  the full life of a TCP connection.
- **Advanced (Parts X–XIV):** congestion control algorithms, TCP tuning, TLS/QUIC,
  packet capture, and performance forensics.
- **Expert (Parts XV–XVIII):** kernel offloads, data-center and cloud networking,
  security attacks and defenses, and structured troubleshooting playbooks.

Each chapter ends with **"Try it yourself"** commands (Linux/macOS) and a
**"War story"** — a realistic failure and its root cause.

---

## Table of contents

**Part I — Foundations**
1. What a network is, and why TCP/IP won
2. The TCP/IP model (and how it maps to OSI)
3. Encapsulation: the nested-envelope model
4. End-to-end walkthrough: what happens when you open `https://example.com`

**Part II — The Link Layer**
5. Ethernet, MAC addresses, and frames
6. ARP: turning IP addresses into MAC addresses
7. Switches, broadcast domains, and VLANs
8. MTU, fragmentation triggers, and jumbo frames
9. Wi-Fi and other link layers

**Part III — The Internet Layer (IPv4)**
10. The IPv4 address and header, field by field
11. Subnetting, CIDR, and VLSM with worked examples
12. Special-use addresses, private space, and NAT/PAT
13. Routing: tables, longest-prefix match, and default routes
14. ICMP: ping, traceroute, and Path MTU Discovery
15. IP fragmentation and reassembly in detail

**Part IV — The Internet Layer (IPv6)**
16. IPv6 addressing and notation
17. Neighbor Discovery, SLAAC, and DHCPv6
18. Dual-stack, Happy Eyeballs, and transition mechanisms

**Part V — Host configuration**
19. DHCP, from DISCOVER to lease renewal
20. DNS: resolution, records, caching, and modern transports

**Part VI — The Transport Layer: UDP**
21. Ports, sockets, and the connection 4-tuple
22. UDP: the header, and when "unreliable" is the right choice

**Part VII — The Transport Layer: TCP mechanics**
23. The TCP header, field by field
24. The three-way handshake and connection setup
25. Connection teardown, half-close, and TIME_WAIT
26. The TCP finite state machine

**Part VIII — TCP reliability**
27. Sequence numbers, cumulative ACKs, and retransmission
28. Retransmission timeout (RTO), RTT estimation, and Karn's algorithm
29. Fast retransmit, fast recovery, and SACK

**Part IX — TCP flow and congestion control**
30. Flow control: the sliding window, window scaling, and SWS
31. Congestion control: slow start and congestion avoidance
32. Reno, NewReno, CUBIC, and BBR
33. Nagle, delayed ACK, and the interaction bug
34. Keepalive, and detecting dead peers

**Part X — TCP performance engineering**
35. Bandwidth-delay product and buffer sizing
36. Throughput math: loss, latency, and the Mathis equation
37. Tuning knobs on Linux (sysctl reference)

**Part XI — Application-layer transports**
38. TLS: handshakes (1.0–1.3), records, resumption, termination, and offload
39. HTTP/1.1 vs HTTP/2 vs HTTP/3 (QUIC)
40. Load balancing: L4 vs L7, DSR, and connection handling

**Part XII — Observability and capture**
41. `tcpdump` and BPF filters
42. Reading a Wireshark trace like a detective
43. Socket statistics: `ss`, `/proc`, and connection states

**Part XIII — Performance forensics**
44. Diagnosing latency, loss, throughput, and bufferbloat
45. Emulating bad networks with `tc`/`netem`

**Part XIV — NAT traversal and real-time media**
46. NAT behaviors, STUN, TURN, and ICE
47. Why VoIP and video use UDP, jitter buffers, and RTP

**Part XV — Kernel and NIC internals**
48. The packet's path through the Linux kernel
49. Offloads: checksum, TSO/GSO/GRO, LRO, RSS, and XDP

**Part XVI — Scaling the network**
50. Anycast, ECMP, and BGP in one page
51. QoS, DSCP, and traffic shaping

**Part XVII — Security**
52. Spoofing, SYN floods, RST injection, and off-path attacks
53. Firewalls, stateful inspection, and conntrack
54. DDoS shapes and mitigations

**Part XVIII — Cloud & container networking**
55. Network namespaces, veth pairs, and bridges
56. Overlays: VXLAN, and how Kubernetes pods talk
57. Cloud VPC constructs and their TCP/IP reality

**Part XIX — HTTP in depth**
58. Message framing: where does a request end? (and request smuggling)
59. Caching: `Cache-Control`, `ETag`, conditional requests, `Vary`
60. Cookies and state
61. CORS: the browser's cross-origin rules
62. WebSockets and Server-Sent Events
63. Range requests, streaming, and compression
64. HTTP authentication and authorization on the wire
65. Proxies: forward, reverse, CONNECT, and interception

**Part XX — Monitoring and telemetry**
66. What to measure: the four signal families
67. Device and interface metrics: SNMP and its successors
68. Flow telemetry: NetFlow, IPFIX, and sFlow
69. Active probing, synthetic monitoring, and alerting that works

**Part XXI — Routing protocols in depth**
70. How routing protocols work: the two families
71. OSPF in depth
72. BGP in depth: attributes and path selection
73. BGP operations: peering, scaling, and failures
74. First-hop redundancy, multicast, and MPLS/VRF in brief

**Appendices**
- A. Troubleshooting playbooks (symptom → cause → fix)
- B. Port and protocol quick reference
- C. Subnetting cheat sheet
- D. `tcpdump`/`ss`/`ip` command cookbook
- E. Glossary

---

# Part I — Foundations

## 1. What a network is, and why TCP/IP won

### 1.1 The problem being solved

A network exists to move a chunk of data from a process on one machine to a
process on another machine, across links that are unreliable, of varying speed,
owned by different organizations, and constantly changing shape.

Break that sentence into the hard sub-problems:

| Sub-problem | Question it answers | TCP/IP mechanism |
|---|---|---|
| **Framing** | Where does one message end and the next begin on the wire? | Ethernet frames, PPP, etc. |
| **Local addressing** | Which physical device on this link? | MAC addresses, ARP |
| **Global addressing** | Which host on the planet? | IP addresses |
| **Routing** | Which path do I send it down? | Routing tables, BGP/OSPF |
| **Multiplexing** | Which *application* on that host? | Port numbers |
| **Reliability** | Did it arrive, intact, in order? | TCP sequence/ACK/retransmit |
| **Fairness/speed** | How fast can I send without collapsing the network? | TCP congestion control |
| **Naming** | I only know a name, not an address | DNS |

TCP/IP is not one protocol. It is a *layered set* of protocols where each layer
solves one of those sub-problems and hands the rest to the layer below.

### 1.2 A one-paragraph history (why it matters)

The suite came out of the US DARPA-funded ARPANET research in the 1970s. Vint
Cerf and Bob Kahn's 1974 paper split the original monolithic "Transmission
Control Program" into **IP** (best-effort delivery of packets) and **TCP**
(reliability on top). The design goals — in priority order — were: survive loss
of network hardware, support multiple transport services, accommodate many
different networks, allow distributed management, be cost-effective, allow easy
host attachment, and account for resources. **Security and accurate billing were
explicitly low priority.** That ranking explains almost every rough edge you will
ever hit: NAT exists because host attachment and address conservation mattered;
spoofing is easy because security was #7; congestion control had to be bolted on
after the 1986 "congestion collapse" of the NSFNET.

The competing standard, the OSI protocol suite, was designed by committee, was
more complete on paper, and lost — because TCP/IP had running code, was free with
Berkeley Unix (the "sockets" API, 1983), and was good enough. **"Rough consensus
and running code"** (the IETF motto) beat formal completeness.

### 1.3 The three ideas that make it work

1. **Packet switching, not circuit switching.** The phone network of the era set
   up a dedicated end-to-end circuit per call. TCP/IP chops data into
   independently-addressed **packets** that share links statistically. A link
   with 1 Gbit/s capacity can serve thousands of mostly-idle connections because
   they rarely all burst at once.

2. **The end-to-end principle.** Keep the network core dumb and fast; push
   intelligence (reliability, encryption, ordering) to the endpoints. A router
   in the middle just forwards packets toward their destination. It does *not*
   track your TCP connection, retransmit for you, or care what's inside. This is
   why you can invent QUIC (a new reliable transport) without upgrading a single
   router.

3. **The narrow waist.** Everything runs *over* IP, and IP runs *over* everything.
   Draw the protocol suite as an hourglass: many link layers at the bottom (Wi-Fi,
   Ethernet, LTE, fiber), many applications at the top (HTTP, DNS, SMTP, SSH),
   and a single skinny neck in the middle — **IP**. Agreeing on exactly one
   internetworking protocol is what lets any app talk over any medium.

```
   HTTP  DNS  SMTP  SSH  RTP  QUIC  ...        <- many applications
      \    \   |    |   /    /
             TCP        UDP                     <- a few transports
                \      /
                  IP                            <- THE narrow waist
                /      \
      Ethernet  Wi-Fi  LTE  DOCSIS  MPLS ...    <- many link layers
```

### 1.4 War story

A startup ran its own "reliable UDP" protocol between microservices to "avoid TCP
overhead." It worked in the lab. In production, p99 latency was 40x worse than a
sister service on plain TCP. Root cause: they had reimplemented retransmission but
not congestion control, so under load their traffic caused queue buildup at a
shared switch, which caused loss, which caused more retransmits — a miniature
congestion collapse. **Lesson:** the hard part of TCP is not the handshake, it's
the 35 years of congestion-control tuning you get for free.

**Try it yourself**
```bash
# See the layers in one packet
sudo tcpdump -i any -c1 -vvv -e 'tcp port 443'
# -e shows the Ethernet (link) header; -vvv expands IP and TCP headers.
```

---

## 2. The TCP/IP model (and how it maps to OSI)

### 2.1 The four layers

TCP/IP is usually drawn with **4 layers** (the original RFC 1122 "host
requirements" model). The 7-layer OSI model is a teaching tool; the wire only
cares about the 4.

```
+----------------------+   Data unit        Example protocols
| 4. Application        |   "message"/data   HTTP, DNS, TLS, SSH, SMTP, gRPC
+----------------------+
| 3. Transport          |   segment (TCP)    TCP, UDP, QUIC(*), SCTP
|                       |   datagram (UDP)
+----------------------+
| 2. Internet           |   packet/datagram  IPv4, IPv6, ICMP, IGMP
+----------------------+
| 1. Link (Network      |   frame            Ethernet, Wi-Fi (802.11), PPP, ARP
|    Access)            |   bits             (plus the physical medium)
+----------------------+
```
(*) QUIC technically rides *inside* UDP but functionally replaces TCP; more in
Chapter 39.

### 2.2 Mapping to OSI

| OSI # | OSI name | TCP/IP layer | What lives here in practice |
|---|---|---|---|
| 7 | Application | Application | HTTP semantics, DNS queries |
| 6 | Presentation | Application | TLS encryption, character encoding, gzip |
| 5 | Session | Application | TLS session resumption, SIP dialogs, RPC |
| 4 | Transport | Transport | TCP, UDP — ports, reliability |
| 3 | Network | Internet | IP addressing and routing |
| 2 | Data Link | Link | Ethernet/Wi-Fi framing, MAC, switching |
| 1 | Physical | Link | Copper, fiber, radio, voltage/light levels |

People still say "Layer 3" for IP, "Layer 4" for TCP/UDP, and "Layer 7" for
application content because that vocabulary is baked into product names:

- **L2 switch** — forwards Ethernet frames by MAC address.
- **L3 switch / router** — forwards IP packets by destination address.
- **L4 load balancer** — picks a backend using IP + port, never looks at payload.
- **L7 load balancer / proxy** — parses HTTP, routes by URL path or `Host` header.
- **"Layer 8 problem"** — a joke: the user.

### 2.3 What each layer promises the one above it

- **Link → Internet:** "I will deliver this frame to another interface *on this
  same physical network*, or drop it. No promises about order or success."
- **Internet → Transport:** "I will make a best effort to deliver this packet to
  the destination IP anywhere on the internet. It may be lost, duplicated,
  reordered, delayed, or corrupted (header checksum only in v4)."
- **Transport (UDP) → App:** "Same as IP, plus I'll tell you which port it's for
  and checksum the payload."
- **Transport (TCP) → App:** "You get a bidirectional byte stream. Bytes arrive
  exactly once, in order, or the connection breaks and you find out."

Every layer is a *contract*. Debugging is usually the act of finding which layer
broke its contract.

### 2.4 War story

An engineer swore "the network is dropping my data" because a large HTTP response
was truncated. Packet capture showed every TCP segment ACKed and the connection
cleanly closed with FIN. TCP kept its contract perfectly. The real bug: the
application server had a 30-second write timeout and stopped sending mid-body
while returning HTTP 200. **Lesson:** "the network" is the last place to look,
not the first. Verify the transport-layer contract held before blaming lower
layers.

**Try it yourself**
```bash
# Watch which layer a tool operates at
ip -br addr          # Layer 3 config (addresses)
ip -br link          # Layer 2 config (MAC, state)
ss -tan              # Layer 4 (TCP sockets and states)
curl -v https://example.com   # Layer 7 (HTTP request/response)
```

---

## 3. Encapsulation: the nested-envelope model

### 3.1 The core mechanic

As your data goes *down* the stack on the sender, each layer wraps it in a header
(and sometimes a trailer). This is **encapsulation**. On the receiver, each layer
strips its header and hands the rest up — **decapsulation**. The payload of one
layer is the entire packet of the layer above.

```
Application data:                          [ HTTP request bytes ]
                                                   |
TCP adds header (ports, seq, flags):    [ TCP hdr | HTTP request bytes ]   = "segment"
                                                   |
IP adds header (src/dst IP, TTL):  [ IP hdr | TCP hdr | HTTP request ]     = "packet"
                                                   |
Ethernet adds header + trailer:
[ Eth hdr | IP hdr | TCP hdr | HTTP request | Eth FCS ]                    = "frame"
                                                   |
                                            serialized as bits on the wire
```

Physical analogy: you write a letter (data), seal it in an envelope addressed to
a person (TCP port = which person in the building), put that envelope inside a
courier pouch addressed to a building (IP address), and hand the pouch to a
driver who only knows the next depot (MAC address / next hop). Each handler reads
only their own label.

### 3.2 Header sizes to memorize

| Header | Typical size | Notes |
|---|---|---|
| Ethernet II | 14 bytes (+ 4 FCS, + 4 per VLAN tag) | 6 dst MAC, 6 src MAC, 2 ethertype |
| IPv4 | 20 bytes (no options) | up to 60 with options |
| IPv6 | 40 bytes fixed | extension headers are separate |
| TCP | 20 bytes (no options) | commonly 32 with timestamps + SACK |
| UDP | 8 bytes | src port, dst port, length, checksum |

**Overhead example.** On classic Ethernet (MTU 1500), a full TCP/IPv4 segment
carries `1500 − 20 (IP) − 20 (TCP) = 1460` bytes of application data — the
**MSS** (Maximum Segment Size). With TCP timestamps enabled (the common case),
options eat 12 more bytes, so MSS is often **1448**. Sending 1 GB of data
therefore costs roughly `1e9/1448 × (14+20+32+4) ≈ 48 MB` of pure header/framing
overhead (~4.8%).

### 3.3 The same bytes, four names

The industry uses different nouns for "the thing at this layer," and mixing them
up causes confused conversations:

- **Frame** — link layer (has MAC addresses).
- **Packet** — network layer (has IP addresses). Also the generic word.
- **Datagram** — a self-contained packet with no connection state. Used for both
  IP packets and UDP messages.
- **Segment** — a TCP unit (part of the byte stream).
- **Message** — application layer.

### 3.4 Demultiplexing keys

How does the receiver know which upper protocol to hand the payload to?

- Ethernet header's **EtherType** field: `0x0800` = IPv4, `0x86DD` = IPv6,
  `0x0806` = ARP.
- IP header's **Protocol** field (v4) / **Next Header** (v6): `6` = TCP,
  `17` = UDP, `1` = ICMP, `58` = ICMPv6.
- TCP/UDP **destination port**: `443` → your HTTPS server process, `53` → DNS.

So a browser packet is delivered by the chain: NIC → (EtherType 0x0800) → IP
stack → (Protocol 6) → TCP → (dst port matches a listening socket's 4-tuple) →
your process's file descriptor.

### 3.5 War story

A team saw mysterious 4-byte-larger frames and MTU errors after a data-center
migration. Cause: the new top-of-rack switches added an **802.1Q VLAN tag** (4
bytes) to every frame, pushing some full-size packets to 1518 bytes on links
where an intermediate device enforced a strict 1514 limit. Fix: lower the host
MTU to 1496 or enable jumbo/baby-jumbo frames end to end. **Lesson:**
encapsulation overhead is not fixed; every tag, tunnel, and VPN steals bytes from
your payload budget (see Chapter 8 and 56).

**Try it yourself**
```bash
# Show encapsulation layers explicitly
sudo tcpdump -i any -c1 -e -n 'tcp port 443' -vv
# Read from outside in: MAC.. > MAC.., ethertype IPv4, IP.. > IP.., Flags [S], ...
```

---

## 4. End-to-end walkthrough: opening `https://example.com`

This chapter ties Parts II–XI together. Follow one HTTPS page load from a laptop
on home Wi-Fi. Every step names the chapter that covers it in depth.

### 4.1 Setup assumed

- Laptop: `192.168.1.50/24`, gateway `192.168.1.1`, DNS `192.168.1.1`.
- Home router does NAT to public IP `203.0.113.7`.
- `example.com` resolves to `93.184.216.34` (IPv4) via an authoritative server.

### 4.2 The steps

**Step 0 — Is the network up? (Ch. 5, 19)**
The laptop already has an IP because at join time it did **DHCP**:
DISCOVER → OFFER → REQUEST → ACK, learning address, mask, gateway, DNS, and lease
time. The NIC link is up (Layer 1/2 negotiated speed and duplex).

**Step 1 — Name resolution (Ch. 20)**
The app calls `getaddrinfo("example.com", "443")`. The stub resolver checks the
local cache, then sends a **DNS** query (UDP/53) to `192.168.1.1`:
`example.com. IN A?`. The home router is a forwarding resolver; it recurses (or
asks upstream), gets `A 93.184.216.34` with a TTL, caches it, and answers the
laptop. Total added latency: one or two RTTs, or ~0 ms if cached.

**Step 2 — Route lookup (Ch. 13)**
The kernel needs to send to `93.184.216.34`. It consults the routing table.
`93.184.216.34` is not in `192.168.1.0/24`, so it matches the **default route**
`0.0.0.0/0 via 192.168.1.1 dev wlan0`. Next hop = `192.168.1.1`.

**Step 3 — Link-layer address of the next hop (Ch. 6)**
The laptop needs the **MAC address** of `192.168.1.1`. It checks the ARP cache.
On a miss it broadcasts `ARP: who-has 192.168.1.1 tell 192.168.1.50`. The router
replies with its MAC. Now the laptop can build frames.

**Step 4 — TCP three-way handshake (Ch. 24)**
The kernel picks an ephemeral source port (say `51000`) and sends:
- `SYN` seq=`X`, options: MSS=1460, SACK-permitted, window scale=7, timestamps.
- Server replies `SYN, ACK` seq=`Y`, ack=`X+1`, its own options.
- Laptop sends `ACK` ack=`Y+1`. Connection is `ESTABLISHED`.

Each of these is one IP packet in one Ethernet frame. The frame's dst MAC is the
router's; the packet's dst IP is the server's. **The MAC changes at every hop;
the IP addresses do not** (until NAT — see next step).

**Step 5 — NAT at the home router (Ch. 12)**
The router rewrites the outbound packet's source from `192.168.1.50:51000` to
`203.0.113.7:62000` (a port it allocates), records the mapping in its NAT table,
recomputes checksums, and forwards it toward its ISP via *its* default route.
Return packets to `203.0.113.7:62000` get rewritten back to `192.168.1.50:51000`.

**Step 6 — Routing across the internet (Ch. 13, 50)**
Each router along the path does the same thing: strip the incoming Ethernet
frame, look up the destination IP in its (BGP-populated) forwarding table, decide
the next hop, decrement **TTL**, rewrite src/dst MAC for the outgoing link,
recompute the IP header checksum, transmit. Typically 10–20 hops. No router keeps
per-connection state.

**Step 7 — TLS handshake (Ch. 38)**
Over the now-open TCP byte stream, the client sends `ClientHello` (TLS 1.3:
supported ciphers, key share, SNI = `example.com`). Server responds with
`ServerHello`, certificate, and a key share. Both derive session keys. **One
round trip** in TLS 1.3 (or zero with resumption). After this, all TCP payload
bytes are encrypted `application_data` records.

**Step 8 — The HTTP request (Ch. 39)**
The client writes:
```
GET / HTTP/1.1
Host: example.com
User-Agent: ...
Accept: */*
```
This is a few hundred bytes → one TCP segment → one IP packet → one frame.

**Step 9 — The response and TCP reliability (Ch. 27–29)**
The server streams the HTML. TCP breaks it into MSS-sized segments with
increasing sequence numbers. The client ACKs (often one ACK per two segments —
**delayed ACK**). If a segment is lost, the client sends duplicate ACKs; after
three, the server does a **fast retransmit** without waiting for a timeout. The
server ramps its sending rate via **slow start** then **congestion avoidance /
CUBIC**.

**Step 10 — Rendering, then teardown (Ch. 25)**
Browser parses HTML, opens more connections (or HTTP/2 streams) for CSS/JS/images
— each repeating steps 4–9, but DNS and often TCP/TLS are cached/reused. When
idle, connections close: `FIN`/`ACK` in each direction (4 packets, or 3 if
piggybacked). The initiator sits in **TIME_WAIT** for 2×MSL (~60 s on Linux) to
absorb stragglers.

### 4.3 The latency budget

For a cold load from Europe to a US server (~90 ms RTT), the *unavoidable*
serialized round trips are roughly:

| Phase | Round trips | ~Time |
|---|---|---|
| DNS (uncached, one level) | 1 | 90 ms |
| TCP handshake | 1 | 90 ms |
| TLS 1.3 handshake | 1 | 90 ms |
| HTTP request + first byte | 1 | 90 ms |
| **Total to first byte** | **~4** | **~360 ms** |

This is why CDNs (Ch. 50 anycast), TLS resumption, HTTP/3's 0-RTT, and connection
reuse matter so much: each eliminated round trip saves one whole RTT.

### 4.4 War story

"The site is slow only for our Sydney office." Traceroute showed traffic from
Sydney going to a data center in Virginia. DNS was returning the US IP because
the office's recursive resolver was in Virginia, and the CDN used **DNS-based
geo-routing** keyed on resolver location, not client location. Fix: move the
resolver, or adopt EDNS Client Subnet, or use an anycast CDN. **Lesson:** every
step in Chapter 4 is a place performance can go wrong; know the whole chain.

**Try it yourself**
```bash
curl -w '
dns:      %{time_namelookup}s
connect:  %{time_connect}s
tls:      %{time_appconnect}s
ttfb:     %{time_starttransfer}s
total:    %{time_total}s
' -o /dev/null -s https://example.com
```

---
# Part II — The Link Layer

The link layer moves a frame between two interfaces that share a physical
network (a "broadcast domain"). It has no idea the internet exists. Its job:
framing, local addressing, and error *detection* (not correction).

## 5. Ethernet, MAC addresses, and frames

### 5.1 The MAC address

A **MAC (Media Access Control) address** is a 48-bit identifier burned into (or
assigned to) a network interface. Written as six hex octets:
`a4:83:e7:2b:19:0c`.

Structure:
```
 a4:83:e7 : 2b:19:0c
 |______|   |______|
 OUI (vendor)  NIC-specific
```
- First 24 bits = **OUI** (Organizationally Unique Identifier), assigned by IEEE
  to the manufacturer. `a4:83:e7` = Apple, `00:1a:11` = Google, `52:54:00` =
  QEMU/KVM virtual NICs.
- Two special bits in the first octet:
  - **I/G bit** (least significant bit of first octet): 0 = unicast, 1 =
    multicast. `ff:ff:ff:ff:ff:ff` = broadcast.
  - **U/L bit** (second-least-significant): 0 = globally unique (OUI-based),
    1 = locally administered. Randomized Wi-Fi MACs and most VMs set this bit —
    that's why a randomized phone MAC looks like `9a:...` or `d6:...` (second hex
    digit is 2, 6, A, or E).

MAC addresses are **flat** (no hierarchy, no aggregation) and only meaningful on
the local link. A packet's MAC addresses are rewritten by every router it
crosses; its IP addresses are not (barring NAT).

### 5.2 The Ethernet II frame format

```
+-----------+-----------+-----------+-------------------+---------+
| Dst MAC   | Src MAC   | EtherType | Payload           | FCS     |
| 6 bytes   | 6 bytes   | 2 bytes   | 46-1500 bytes     | 4 bytes |
+-----------+-----------+-----------+-------------------+---------+
        (preamble + SFD, 8 bytes, added by hardware, not counted)
```

- **EtherType** ≥ `0x0600` identifies the payload protocol (`0x0800` IPv4,
  `0x86DD` IPv6, `0x0806` ARP, `0x8100` 802.1Q VLAN tag). Values ≤ `0x05DC` mean
  the field is instead a *length* (old 802.3 framing).
- **Payload minimum 46 bytes**: frames shorter than 64 bytes total collide-detect
  incorrectly on legacy half-duplex Ethernet, so short payloads are **padded**.
  A 1-byte TCP ACK's payload is padded to 46; Wireshark shows "trailer" bytes.
- **FCS (Frame Check Sequence)**: a CRC-32 over the frame. On mismatch the NIC
  silently **drops** the frame and increments an error counter. Ethernet does not
  retransmit — that's TCP's job.
- **MTU** (Maximum Transmission Unit) = the max payload = **1500 bytes** for
  standard Ethernet. This single number propagates constraints all the way up to
  TCP's MSS.

### 5.3 With a VLAN tag (802.1Q)

```
| Dst MAC | Src MAC | 0x8100 | PCP/DEI/VID (2B) | EtherType | Payload | FCS |
```
The 4-byte tag carries a 12-bit **VLAN ID** (1–4094) and a 3-bit **PCP**
(priority, used for QoS — Ch. 51). This is why VLAN-tagged frames can reach 1518
or 1522 bytes and trip strict MTU checks.

### 5.4 How a host decides: local or remote?

Before sending an IP packet, the stack computes
`(destination_IP AND local_netmask) == (local_IP AND local_netmask)`:
- **Equal** → destination is on my link. ARP for the destination's own MAC, send
  directly.
- **Not equal** → destination is remote. ARP for the **gateway's** MAC, send the
  frame there; the IP header still holds the final destination IP.

Example: laptop `192.168.1.50/24` sending to `192.168.1.99` → same `/24`, direct.
Sending to `8.8.8.8` → different, send to gateway MAC.

### 5.5 Duplex, autonegotiation, and the classic mismatch

Modern Ethernet is **full duplex** (send and receive simultaneously, no
collisions). Autonegotiation picks speed (10/100/1000/…) and duplex. If one side
is hard-set to `100/full` and the other autonegotiates, the autoneg side falls
back to `100/half` → **duplex mismatch**: works at low load, throws
`CRC errors` + `late collisions` and collapses throughput under load. Always set
both ends the same way (usually: autoneg on both).

### 5.6 War story

A monitoring host showed 0.5% packet loss to one server, steady for months.
`ethtool -S eth0` on the switch-facing NIC showed `rx_crc_errors` climbing. A
half-crimped RJ45 connector was corrupting ~1 frame in 200. TCP hid it (retransmits)
until a backup job needed line-rate and throughput cratered to 6 Mbit/s.
**Lesson:** Layer 1 errors masquerade as Layer 4 performance problems. Check NIC
error counters early.

**Try it yourself**
```bash
ip -br link                       # interface names, MACs, up/down
ethtool eth0                      # speed, duplex, autoneg
ethtool -S eth0 | grep -Ei 'err|drop|crc|collision'   # hardware error counters
cat /sys/class/net/eth0/mtu      # link MTU
ip neigh                         # the ARP/neighbor cache
```

---

## 6. ARP: turning IP addresses into MAC addresses

### 6.1 Why it exists

You have the next hop's **IP** (from routing). To build the Ethernet frame you
need its **MAC**. **ARP (Address Resolution Protocol, RFC 826)** bridges Layer 3
to Layer 2 on IPv4 networks. (IPv6 uses NDP instead — Ch. 17.)

### 6.2 The exchange

```
Host A (192.168.1.50, MAC aa:aa) wants to reach 192.168.1.1

A --> BROADCAST (dst ff:ff:ff:ff:ff:ff):
      "ARP Request: Who has 192.168.1.1? Tell 192.168.1.50 (aa:aa)"

Router (192.168.1.1, MAC bb:bb) --> UNICAST to aa:aa:
      "ARP Reply: 192.168.1.1 is at bb:bb"

A caches: 192.168.1.1 -> bb:bb   (state REACHABLE, ~30 s on Linux, then STALE)
```

ARP frames use EtherType `0x0806` and are **not** IP packets — no IP header, no
TTL, cannot cross a router. ARP scope = one broadcast domain.

### 6.3 The ARP cache and its states (Linux `ip neigh`)

| State | Meaning |
|---|---|
| `REACHABLE` | Confirmed recently; use freely. |
| `STALE` | Entry exists but unconfirmed; used, then revalidated on next send. |
| `DELAY` | Waiting a moment for upper-layer (TCP ACK) confirmation before probing. |
| `PROBE` | Actively sending unicast ARP requests to reconfirm. |
| `FAILED` | No answer; will retry on demand. |
| `PERMANENT` | Static entry, never expires. |

### 6.4 Gratuitous ARP

A host sends an **unsolicited** ARP reply/request for *its own* IP:
- On boot or IP assignment → **duplicate address detection** (if someone replies,
  there's a conflict).
- After a failover → tell every host "the IP `10.0.0.10` now lives at *my* MAC" so
  switches and hosts update their tables immediately. This is how VRRP/keepalived,
  cluster VIPs, and cloud "elastic IP" reattachment redirect traffic in <1 s.

### 6.5 Proxy ARP

A router answers ARP requests for IPs that are not its own, on behalf of hosts it
can reach, making two subnets *look* like one. Mostly legacy; still seen in some
VPN and hotspot setups.

### 6.6 Attacks: ARP spoofing / cache poisoning

ARP has **no authentication**. Any host can send "10.0.0.1 is at my MAC" and
receivers will believe it. Result: attacker becomes a **man-in-the-middle** for
the gateway. Tools: `arpspoof`, `ettercap`, `bettercap`. Defenses:
- **Dynamic ARP Inspection (DAI)** on switches (validates ARP against the
  DHCP-snooping table).
- Static ARP entries for critical gateways.
- Encrypt everything end-to-end (TLS) so MITM yields ciphertext.
- Port security / 802.1X.

### 6.7 War story

Two servers were configured with the same IP after a provisioning script bug.
Symptoms: intermittent connection resets, `arp` output flapping between two MACs
for that IP, kernel log `neighbour ... is not permanent`. Every gratuitous ARP
from one server stole traffic from the other. **Lesson:** "IP address already in
use" and flapping `ip neigh` output = duplicate IP or a failover event, not
"random network weirdness."

**Try it yourself**
```bash
ip neigh show                       # current table + states
sudo tcpdump -i eth0 -n arp         # watch requests/replies live
ip neigh flush all                  # clear cache (forces re-ARP)
arping -I eth0 192.168.1.1          # manually probe a neighbor
sudo ip neigh add 192.168.1.1 lladdr bb:bb:bb:bb:bb:bb dev eth0 nud permanent
```

---

## 7. Switches, broadcast domains, and VLANs

### 7.1 What a switch actually does

An Ethernet **switch** is a Layer 2 device that forwards frames based on a
**MAC address table** (a.k.a. CAM/FIB table) it builds by *learning*:

1. Frame arrives on port 3 with src MAC `aa:aa` → switch records
   `aa:aa → port 3` (with an aging timer, default 300 s).
2. Frame's dst MAC `bb:bb` is looked up:
   - **Hit** → forward out only that one port (**unicast forwarding**).
   - **Miss** (unknown) → **flood** out every port except the ingress
     (**unknown-unicast flooding**), hoping for a reply that teaches the table.
3. Broadcast (`ff:ff:...`) or multicast → flooded within the VLAN.

A switch has no IP involvement, no TTL, and (ideally) wire-speed forwarding in
hardware (ASIC).

### 7.2 Collision domain vs broadcast domain

- **Collision domain:** set of interfaces that can collide with each other. With
  full-duplex switched Ethernet, each port is its own collision domain — so
  collisions are effectively extinct. (Hubs, now obsolete, put everyone in one.)
- **Broadcast domain:** set of interfaces a broadcast frame reaches. **A switch
  does not break up a broadcast domain; a router (or a VLAN boundary) does.**
  Every host in one broadcast domain shares ARP traffic, DHCP broadcasts, etc.
  Too large → "broadcast storm" risk and wasted CPU. Rule of thumb: keep a
  broadcast domain under ~500 hosts, or one `/23`.

### 7.3 VLANs: many virtual switches in one

A **VLAN (Virtual LAN, 802.1Q)** partitions one physical switch into multiple
isolated broadcast domains. Ports are assigned VLAN IDs:

- **Access port:** belongs to exactly one VLAN, sends/receives *untagged* frames
  (the endpoint doesn't know VLANs exist). "This port is in VLAN 20."
- **Trunk port:** carries many VLANs between switches (or to a hypervisor/router),
  frames are **tagged** with their VLAN ID. One "native VLAN" may be untagged.

Traffic between VLANs *must* go through a router or L3 switch
("**router-on-a-stick**" or SVI). This is where inter-subnet policy/firewalling
naturally lives.

### 7.4 Spanning Tree Protocol (STP)

Redundant links between switches create **Layer 2 loops** — and a broadcast frame
in a loop circulates forever (no TTL at L2), melting the network in seconds.
**STP (802.1D)** and successors (RSTP 802.1w, MSTP) elect a root bridge and
*block* redundant ports, leaving a loop-free tree; a blocked port unblocks if the
active path fails. Modern data centers often avoid STP entirely with
routed/leaf-spine designs (Ch. 50) or link aggregation.

### 7.5 Link aggregation (LACP / bonding)

Bundle N physical links into one logical link for bandwidth and redundancy
(`bond0`, `port-channel`, EtherChannel). Frames are hashed onto member links by
`src/dst MAC`, `src/dst IP`, or `L4 ports` — so a *single* TCP flow still rides
*one* physical link and can't exceed one link's speed. Aggregation adds capacity
for *many* flows, not for one.

### 7.6 War story

After adding a second uplink "for redundancy," a whole floor went dark every few
minutes. Someone had cabled two access switches together *and* left both plugged
into the core — a loop — and STP was disabled on the cheap switches. Broadcast
storm. Re-enabling RSTP (or just pulling one cable) fixed it instantly.
**Lesson:** never create a second L2 path without STP or a LAG. At L2 there is no
TTL to save you.

**Try it yourself**
```bash
bridge fdb show                     # Linux bridge MAC table (a software switch)
ip -d link show type vlan           # VLAN sub-interfaces
cat /proc/net/bonding/bond0         # bond status and active members
sudo tcpdump -i eth0 -e vlan        # show 802.1Q tags
```

---

## 8. MTU, fragmentation triggers, and jumbo frames

### 8.1 The MTU chain

**MTU** is the largest Layer-3 payload a link will carry. The *effective*
end-to-end MTU is the **minimum** of every link on the path — the **Path MTU
(PMTU)**.

| Link | Typical MTU |
|---|---|
| Standard Ethernet | 1500 |
| Ethernet + 1 VLAN tag on a strict device | effectively 1496 |
| PPPoE (DSL) | 1492 |
| IPsec tunnel (varies by mode/cipher) | ~1400–1438 |
| WireGuard | 1420 |
| GRE | 1476 |
| VXLAN over 1500 underlay | 1450 |
| Jumbo frames (data-center Ethernet) | 9000 |
| IPv6 minimum any link must support | 1280 |
| IPv4 minimum every host must reassemble | 576 |

### 8.2 What happens when a packet is too big

- **IPv4, DF (Don't Fragment) bit = 0:** the router fragments the packet
  (Ch. 15). Slow, fragile, disables some hardware offloads.
- **IPv4, DF = 1:** router **drops** it and sends back
  `ICMP type 3 code 4 "Fragmentation needed"` including the next-hop MTU.
- **IPv6:** routers **never** fragment. The source must; an oversize packet gets
  `ICMPv6 type 2 "Packet Too Big"`.

TCP always sets DF and relies on **PMTUD** (Path MTU Discovery): start at the
local MTU, shrink when a "too big" ICMP arrives, cache per-destination.

### 8.3 The PMTUD black hole

If a firewall **blocks all ICMP** (a depressingly common "hardening" mistake),
the "fragmentation needed" messages never arrive. Symptom: **the TCP handshake
succeeds** (small packets), then the connection **hangs** the instant a full-size
data packet is sent — it's silently dropped forever. Classic signature:
`curl` connects, sends the request, and stalls; small pages work, big ones don't;
`ping` works but `ping -s 1472 -M do` fails.

Fixes:
- Allow `ICMP type 3` (and ICMPv6 type 2) through firewalls. Always.
- **MSS clamping** on the tunnel/router:
  `iptables -t mangle -A FORWARD -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu`.
  This rewrites the MSS option in SYN packets so endpoints never send oversize
  segments in the first place.
- Linux `ip route ... mtu lock 1400` for a known-bad path.
- **PLPMTUD** (RFC 4821 / RFC 8899) — probe for MTU using the transport itself,
  no ICMP needed. QUIC uses this.

### 8.4 Jumbo frames

Setting MTU 9000 on every device in a data-center segment means ~6× fewer packets
per MB → less CPU, higher throughput, especially for storage (iSCSI, NFS) and
backups. Rules:
- **Every** device on the L2 segment must agree, including the switches. One
  device at 1500 that receives a 9000 frame just drops it.
- Doesn't help across the internet (path will have a 1500 link).
- Test: `ping -M do -s 8972 <peer>` (8972 + 8 ICMP + 20 IP = 9000).

### 8.5 War story

A new site-to-site VPN worked for SSH and ping but every `git clone` and file
copy froze at a few KB. `tcpdump` showed a 1500-byte packet leaving repeatedly
with no ACK. The firewall vendor's "recommended" ruleset dropped all ICMP.
Adding an allow rule for ICMP type 3, plus MSS clamping to 1400 on the tunnel
interface, fixed every stuck transfer. **Lesson:** "connects then hangs on big
transfers" ≈ MTU/PMTUD black hole, 95% of the time.

**Try it yourself**
```bash
# Find the path MTU to a host (largest non-fragmented ICMP payload)
ping -M do -s 1472 example.com          # 1472+8+20 = 1500; shrink until it passes
tracepath example.com                    # reports PMTU per hop, no root needed
ip route get 8.8.8.8                      # shows the MTU that will be used
ip route add 10.0.0.0/8 via 10.1.1.1 mtu 1400   # pin an MTU for a route
```

---

## 9. Wi-Fi and other link layers

IP does not care what the link is, but the link's quirks leak upward as latency,
loss, and jitter that TCP must cope with.

### 9.1 Wi-Fi (802.11) vs Ethernet — the differences that matter

| Property | Switched Ethernet | Wi-Fi |
|---|---|---|
| Medium | Dedicated per port | **Shared** radio, half-duplex |
| Access method | None needed (full duplex) | **CSMA/CA** — listen, backoff, maybe RTS/CTS |
| Link-layer ACKs | No (FCS drop only) | **Yes** — every unicast frame is ACKed and retried at L2 |
| Rate | Fixed, negotiated once | **Adapts constantly** (MCS) with signal quality |
| Loss | ~0 on healthy cable | Bursty; interference, distance, hidden nodes |
| Latency | ~microseconds, stable | 1–50 ms, **highly variable** (airtime contention) |
| Addressing | 2 MACs per frame | Often **3 or 4** (adds the AP/BSSID) |

Key consequences for TCP:

- **Wi-Fi retries hide loss from TCP** — good (fewer spurious congestion events)
  but retries *add latency and jitter*, which inflates RTT estimates and can
  cause **bufferbloat** in the AP's queue (Ch. 44).
- **Airtime is the real currency.** One slow/distant client transmitting at
  6 Mbit/s hogs airtime and starves fast clients ("performance anomaly").
- **Half-duplex + contention** means measured "100 Mbit/s Wi-Fi" delivers ~50–60
  Mbit/s of TCP goodput at best, less with many clients.
- **Roaming** between APs drops frames briefly; long-lived TCP connections
  survive because of retransmission, but real-time media hiccups.

### 9.2 Other link layers you'll meet

- **PPP / PPPoE:** point-to-point (DSL, some fiber, cellular). MTU 1492 for
  PPPoE. Carries authentication (PAP/CHAP) and IP assignment.
- **DOCSIS (cable):** shared upstream, asymmetric, historically bufferbloat-prone
  until "PIE"/AQM mandates.
- **Cellular (LTE/5G):** deep per-device buffers, RTT 20–80 ms that spikes on
  radio state transitions (idle→connected can add 100+ ms), NAT everywhere,
  frequent IP changes on handover. TCP sees this as latency and occasional
  connection death.
- **MPLS:** a "Layer 2.5" label-switching transport carriers use internally;
  invisible to your IP packets except as one fast hop and possible ICMP oddities.
- **Loopback (`lo`):** MTU 65536, no real link; used for host-local IPC over
  `127.0.0.1`/`::1`.

### 9.3 War story

A video-call app was blamed for "freezing," but only in one conference room. The
room's AP shared a channel with a 2.4 GHz wireless projector and a microwave
kitchen next door. `iw dev wlan0 station dump` showed the client rate collapsing
to 1 Mbit/s and 40% frame retries during freezes. Moving the AP to a clean 5 GHz
channel fixed it. **Lesson:** on Wi-Fi, "the app froze" is often "the RF
environment degraded"; measure PHY rate and retry percentage, not just ping.

**Try it yourself**
```bash
# Linux Wi-Fi diagnostics
iw dev wlan0 link                   # current BSSID, signal, tx bitrate
iw dev wlan0 station dump           # rx/tx bitrate, retries, failed, signal
watch -n1 'iw dev wlan0 link | grep -E "signal|bitrate"'
# macOS
/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I
```

---
# Part III — The Internet Layer (IPv4)

IP delivers a packet from any host to any other host, best-effort. No guarantees
of delivery, order, timing, or (for the payload) integrity. Everything above IP
exists to paper over those gaps.

## 10. The IPv4 address and header, field by field

### 10.1 The address

32 bits, written as four **dotted-decimal** octets: `192.168.1.50` =
`11000000.10101000.00000001.00110010`. Range `0.0.0.0`–`255.255.255.255`, ~4.29
billion values — far too few, which is why NAT (Ch. 12) and IPv6 (Part IV) exist.

An address has two parts, split by the **subnet mask**:
```
192.168.1.50 / 24
|__________|   |__|
   address    prefix length (bits of network portion)

netmask /24 = 255.255.255.0 = 11111111.11111111.11111111.00000000
network part:  192.168.1   (first 24 bits - identifies the link)
host part:     .50         (last 8 bits - identifies the host on that link)
```

- **Network address** (all host bits 0): `192.168.1.0` — names the subnet, not
  assignable to a host.
- **Broadcast address** (all host bits 1): `192.168.1.255` — reaches every host
  on the subnet.
- **Usable host range:** `192.168.1.1` – `192.168.1.254` → `2^8 − 2 = 254` hosts.
- The "−2" does **not** apply to `/31` links (RFC 3021, point-to-point, 2 usable)
  or `/32` (single host).

**Classful addressing is dead.** You may still hear "Class A/B/C" (the pre-1993
scheme where the first bits fixed the mask). Since **CIDR** (Classless
Inter-Domain Routing, 1993), the mask is always explicit. Ignore classes except
as historical vocabulary; the class-based default masks in some old tools are a
trap.

### 10.2 The IPv4 header

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Version|  IHL  |    DSCP   |ECN|          Total Length          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         Identification         |Flags|      Fragment Offset    |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Time to Live |    Protocol   |         Header Checksum        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       Source IP Address                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Destination IP Address                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Options (0-40 bytes, rare)                  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Bits | Purpose / real-world relevance |
|---|---|---|
| **Version** | 4 | `4`. |
| **IHL** | 4 | Header length in 32-bit words. `5` = 20 bytes (no options). |
| **DSCP** | 6 | QoS class (Ch. 51). `EF` = expedited forwarding for VoIP; `CS0` = best effort. |
| **ECN** | 2 | Explicit Congestion Notification. `00` not-capable, `10`/`01` ECT, `11` = **CE** ("congestion experienced") — a router marks this instead of dropping. |
| **Total Length** | 16 | Whole packet incl. header. Max 65535; real max ≈ PMTU. |
| **Identification** | 16 | Same value on all fragments of one original packet (Ch. 15). Also abused for OS fingerprinting and idle scans. |
| **Flags** | 3 | bit0 reserved(0); **DF** Don't Fragment; **MF** More Fragments. |
| **Fragment Offset** | 13 | Position of this fragment in the original, in 8-byte units. |
| **TTL** | 8 | Decremented by 1 per hop. At 0 the packet is dropped and `ICMP time-exceeded` returned (this powers `traceroute`). Typical initial values: Linux/macOS 64, Windows 128, some routers 255. |
| **Protocol** | 8 | Payload type: `1` ICMP, `6` TCP, `17` UDP, `47` GRE, `50` ESP (IPsec), `89` OSPF. |
| **Header Checksum** | 16 | Covers the **header only**. Recomputed at every hop (because TTL changes). Does *not* protect the payload — that's TCP/UDP's checksum. |
| **Source / Destination** | 32 each | The endpoints. Rewritten only by NAT. |

### 10.3 What a router does with this header

1. Verify header checksum; drop if bad.
2. Decrement TTL; if 0, drop + ICMP time-exceeded.
3. Longest-prefix-match the destination in the forwarding table → next hop +
   egress interface (Ch. 13).
4. If egress MTU < packet size: fragment (DF=0) or drop + ICMP (DF=1).
5. Recompute header checksum (TTL changed).
6. Resolve next-hop MAC (ARP), build new Ethernet frame, transmit.

The payload is never touched (except by NAT/firew-fixups). No per-flow state.

### 10.4 War story

A CDN reported that a customer's origin was "unreachable from Asia" but fine from
Europe. `traceroute` from Asia died after 12 hops at a specific router. That
router had a bug decrementing TTL by 2; packets with initial TTL 64 that took 32+
hops arrived with TTL 0 and were dropped. Raising the origin's default TTL via
`sysctl net.ipv4.ip_default_ttl=128` was the pragmatic workaround while the
carrier fixed the router. **Lesson:** TTL is a hop budget; unusually long paths
plus buggy hops can exhaust it.

**Try it yourself**
```bash
ip -4 addr show                       # your addresses and prefix lengths
ipcalc 192.168.1.50/24                # network, broadcast, host range
sudo tcpdump -i any -n -v 'ip[8]<5'   # capture packets with TTL < 5 (near-expiry)
sysctl net.ipv4.ip_default_ttl
```

---

## 11. Subnetting, CIDR, and VLSM with worked examples

### 11.1 The mask, three ways

`/24` ⇔ `255.255.255.0` ⇔ `11111111 11111111 11111111 00000000`.
Memorize the last-octet values:

| Prefix | Mask (last non-255 octet) | Block size | Usable hosts |
|---|---|---|---|
| /24 | 0 | 256 | 254 |
| /25 | 128 | 128 | 126 |
| /26 | 192 | 64 | 62 |
| /27 | 224 | 32 | 30 |
| /28 | 240 | 16 | 14 |
| /29 | 248 | 8 | 6 |
| /30 | 252 | 4 | 2 |
| /31 | 254 | 2 | 2 (point-to-point, RFC 3021) |
| /32 | 255 | 1 | 1 (host route) |

Larger blocks:

| Prefix | Mask | # of /24s | Addresses |
|---|---|---|---|
| /23 | 255.255.254.0 | 2 | 512 |
| /22 | 255.255.252.0 | 4 | 1024 |
| /20 | 255.255.240.0 | 16 | 4096 |
| /16 | 255.255.0.0 | 256 | 65536 |
| /8 | 255.0.0.0 | 65536 | 16.7 M |

### 11.2 The four questions for any address+mask

Given `172.16.20.100/22`:

1. **Block size** = `256 − 252 = 4` in the **third** octet (because /22 borrows
   into octet 3). So subnets step by 4 in the third octet:
   `172.16.0.0`, `172.16.4.0`, `172.16.8.0`, `172.16.12.0`, `172.16.16.0`,
   `172.16.20.0`, `172.16.24.0`, …
2. **Which subnet is `.20.100` in?** `20` falls in the `20–23` block →
   **network `172.16.20.0/22`**.
3. **Broadcast** = last address of the block = `172.16.23.255`.
4. **Usable range** = `172.16.20.1` – `172.16.23.254` (1022 hosts).

### 11.3 VLSM: carving a block to fit real needs

**Scenario.** You're given `10.10.0.0/24` for a small office and must serve:

| Segment | Hosts needed | Prefix that fits | Size |
|---|---|---|---|
| Staff LAN | 100 | /25 (126) | 128 |
| Wi-Fi guests | 50 | /26 (62) | 64 |
| Servers | 10 | /28 (14) | 16 |
| Management | 10 | /28 (14) | 16 |
| Router-to-router link | 2 | /30 (2) | 4 |

**Allocate largest first** (to avoid fragmentation), from `10.10.0.0`:

```
10.10.0.0/25    Staff       hosts .1  - .126     bcast .127
10.10.0.128/26  Wi-Fi       hosts .129 - .190    bcast .191
10.10.0.192/28  Servers     hosts .193 - .206    bcast .207
10.10.0.208/28  Management   hosts .209 - .222   bcast .223
10.10.0.224/30  P2P link    hosts .225 - .226    bcast .227
10.10.0.228 .. 10.10.0.255   free for growth
```

If you had allocated the /30 first at `10.10.0.0/30`, the next /25 could not start
until `10.10.0.128` and you'd waste `.4`–`.127`. **Order matters.**

### 11.4 CIDR route aggregation (supernetting)

A router advertising four contiguous /24s —
`203.0.112.0/24, 203.0.113.0/24, 203.0.114.0/24, 203.0.115.0/24` — can announce
them as a single `203.0.112.0/22`. `112` = `01110000`, `115` = `01110011`; the
first 22 bits are common. Aggregation is why the global BGP table is ~1M routes
instead of 16M. Break aggregation (announce a /24 out of your /22 from a
different site) and you've done "**more-specific**" traffic engineering.

### 11.5 The `/31` trick

For a point-to-point link (router↔router), a `/30` wastes 2 of 4 addresses on
network+broadcast. **RFC 3021** `/31` gives you exactly the 2 host addresses you
need; there is no broadcast on a P2P link. Widely supported; halves address burn
on backbones.

### 11.6 War story

An ops team added a new rack on `10.20.5.0/24`, gateway `10.20.5.1`. Hosts could
reach the gateway but not each other's neighbors two racks over. Cause: the core
was configured with `10.20.4.0/23` for the *old* racks, and a stale summary route
`10.20.4.0/22` on a different router black-holed part of the new range. **Lesson:**
subnet/summary boundaries must line up on bit boundaries *and* be consistent
across every device that routes them. Draw the bitmap.

**Try it yourself**
```bash
ipcalc 10.10.0.0/24                       # overview
ipcalc -s 10.10.0.0/24 /25 /26 /28 /28 /30   # split into VLSM subnets (some versions)
sipcalc 172.16.20.100/22
python3 -c "import ipaddress as i; n=i.ip_network('172.16.20.100/22', strict=False); \
print(n, n.network_address, n.broadcast_address, n.num_addresses)"
```

---

## 12. Special-use addresses, private space, and NAT/PAT

### 12.1 Ranges you must recognize

| Range | Name | Meaning |
|---|---|---|
| `10.0.0.0/8` | Private (RFC 1918) | Not routed on the internet |
| `172.16.0.0/12` | Private (RFC 1918) | `172.16.0.0`–`172.31.255.255` |
| `192.168.0.0/16` | Private (RFC 1918) | Home/office LANs |
| `100.64.0.0/10` | CGNAT (RFC 6598) | Carrier-grade NAT shared space |
| `127.0.0.0/8` | Loopback | `127.0.0.1` = this host; whole /8 is local |
| `169.254.0.0/16` | Link-local (APIPA) | Self-assigned when DHCP fails; not routed |
| `0.0.0.0/8` | "This network" | `0.0.0.0` = "any"/unspecified/default route |
| `255.255.255.255` | Limited broadcast | This subnet only, never forwarded |
| `224.0.0.0/4` | Multicast | `224.0.0.1` all-hosts, `224.0.0.251` mDNS |
| `240.0.0.0/4` | Reserved (experimental) | Mostly unusable |
| `192.0.2.0/24`, `198.51.100.0/24`, `203.0.113.0/24` | Documentation (TEST-NET) | Safe to use in docs/examples |
| `198.18.0.0/15` | Benchmarking | Lab/testing only |

### 12.2 Why NAT exists

IPv4 has ~3.7 billion usable addresses; there are ~15 billion connected devices.
**NAT (Network Address Translation)** lets many private hosts share one public
address. It also became a de-facto firewall (unsolicited inbound has nowhere to
go) — which is both a feature and the reason peer-to-peer is hard (Ch. 46).

### 12.3 The kind you actually use: PAT / NAPT / "masquerade"

**Port Address Translation** (one public IP, many hosts) maps the full 4-tuple:

```
NAT table entry (conntrack):
  inside:  192.168.1.50 : 51000  <->  8.8.8.8 : 53
  outside: 203.0.113.7  : 62000  <->  8.8.8.8 : 53

Outbound: rewrite src 192.168.1.50:51000 -> 203.0.113.7:62000, fix checksums
Inbound:  match dst 203.0.113.7:62000    -> rewrite dst 192.168.1.50:51000
```

The router allocates a unique outside port per flow, so two inside hosts both
using source port 51000 to the same server get different outside ports. State
times out (TCP established ~5 days on Linux conntrack default, UDP ~30 s, TCP
after FIN ~1 min).

### 12.4 NAT variants (and why P2P cares — see Ch. 46)

- **Full-cone:** once `inside:port → outside:port` exists, *any* external host can
  send to `outside:port`. Easiest to traverse.
- **Restricted-cone:** external host may reply only if the inside host sent to
  that host's IP first.
- **Port-restricted-cone:** … that host's IP *and* port.
- **Symmetric:** a *different* outside port per destination. Two different servers
  see two different source ports for the same inside socket. Hardest to traverse;
  needs TURN relay.
- **Static NAT / port forwarding:** fixed `outside:port → inside:host:port` for
  hosting a server behind NAT.
- **1:1 NAT (NAT44) / hairpin NAT:** map a whole public IP to a private host;
  hairpin = inside host reaching the public IP of another inside host.
- **CGNAT (RFC 6598):** your ISP NATs *you* inside `100.64.0.0/10`, so you're
  behind two layers of NAT. Breaks port forwarding entirely; this is why IPv6 or
  services like Tailscale matter.

### 12.5 What NAT breaks

- **Embedded addresses:** protocols that put IPs *inside the payload* — FTP
  active mode, SIP, H.323, PPTP GRE — need an **ALG** (Application Layer Gateway)
  on the router to rewrite them, and ALGs are famously buggy.
- **End-to-end IPsec AH** (authenticates the IP header → NAT invalidates it);
  needs NAT-T (UDP/4500 encapsulation).
- **Inbound connections** without explicit port forwarding.
- **Address-based logging/geo/abuse** — thousands of users behind one IP.
- **Connection count limits** — a busy office can exhaust ~64k ports to a single
  popular destination IP.

### 12.6 War story

A SaaS provider IP-allowlisted a customer's office "public IP." It worked for a
week, then broke randomly. The ISP had moved the customer to **CGNAT**, so their
public IP now changed hourly and was shared with strangers. Allowlisting became
impossible; the fix was a site-to-site VPN with a stable tunnel IP. **Lesson:**
never assume a client's source IP is stable or unique — NAT, CGNAT, mobile
networks, and proxies all break that assumption.

**Try it yourself**
```bash
curl -s https://api.ipify.org ; echo          # your public (post-NAT) IP
ip route get 8.8.8.8                            # which src IP/interface is used
sudo conntrack -L | head                        # Linux NAT/state table (conntrack-tools)
sudo iptables -t nat -L -n -v                   # NAT rules
# On a Linux router, classic masquerade:
sudo iptables -t nat -A POSTROUTING -o eth0 -s 192.168.1.0/24 -j MASQUERADE
```

---

## 13. Routing: tables, longest-prefix match, and default routes

### 13.1 The routing table

Every IP host — not just routers — has a routing table. Sending a packet =
looking up its destination and getting `(next hop, egress interface, source IP)`.

```
$ ip route
default via 192.168.1.1 dev wlan0 proto dhcp metric 600
10.8.0.0/24 dev tun0 proto kernel scope link src 10.8.0.2
192.168.1.0/24 dev wlan0 proto kernel scope link src 192.168.1.50 metric 600
```

Read each line as: *"to reach `<prefix>`, send via `<next hop>` out `<dev>`."*
- `scope link` / no `via` → destination is **directly connected**; ARP for the
  destination itself.
- `via X` → send to gateway `X` (ARP for `X`).
- `default` = `0.0.0.0/0` → matches everything not matched more specifically.

### 13.2 Longest-prefix match (LPM)

When multiple routes match, **the one with the longest prefix (most specific)
wins**, regardless of order or metric.

```
Routes:            Packet to 10.1.2.3 matches:
  0.0.0.0/0          yes  (/0)
  10.0.0.0/8         yes  (/8)
  10.1.0.0/16        yes  (/16)
  10.1.2.0/24        yes  (/24)   <-- WINNER (longest)
```

Metric/administrative distance only breaks ties **among routes of equal prefix
length** (or chooses between routing protocols). This is why a `/32` "host route"
always beats a subnet route, and why a leaked more-specific BGP prefix hijacks
traffic (Ch. 50).

### 13.3 How routers learn routes

| Source | Examples | Scope |
|---|---|---|
| **Directly connected** | interface has an IP in that subnet | local |
| **Static** | admin types `ip route add …` | small networks, stub links |
| **Interior Gateway Protocols (IGP)** | **OSPF**, IS-IS (link-state), **EIGRP** (Cisco) | within one organization ("autonomous system") |
| **RIP** | distance-vector, hop-count ≤ 15 | legacy/tiny only |
| **Exterior Gateway Protocol** | **BGP** (path-vector) | *between* autonomous systems — the internet's routing (Ch. 50) |

Link-state (OSPF/IS-IS): every router floods its local link info; all build an
identical map; each runs Dijkstra for shortest paths. Fast reconvergence, scales
to thousands of routers per area.

### 13.4 Forwarding vs routing (control plane vs data plane)

- **Control plane:** the protocols (BGP/OSPF) that *compute* the best paths →
  build the **RIB** (Routing Information Base). Runs on the CPU. Slow, complex.
- **Data plane:** the **FIB** (Forwarding Information Base), a compiled copy of
  the best routes pushed into hardware (TCAM/ASIC) that forwards packets at line
  rate. When people say "the router does 10 Gbit/s," they mean the FIB.

### 13.5 Policy routing (Linux)

Beyond destination-based routing, Linux supports **multiple tables** chosen by
`ip rule` (source address, fwmark, interface, ToS). Used for multi-homing,
VPN split tunneling, and "traffic from container X goes out uplink 2."

```
ip rule add from 192.168.50.0/24 table 100
ip route add default via 10.0.0.1 dev eth1 table 100
```

### 13.6 Asymmetric routing and its traps

Packets can take path A→B outbound and B→A return on a different path. That's
normal on the internet. It breaks:
- **Stateful firewalls** that see only one direction → drop the "unsolicited"
  replies.
- **Linux reverse-path filtering** (`rp_filter=1`): drops a packet if the reply
  wouldn't go back out the interface it arrived on. Set `rp_filter=2` (loose) for
  multi-homed hosts.

### 13.7 War story

A dual-homed server (two ISPs) worked for outbound but external users could reach
it only via ISP-A. Replies to ISP-B-sourced requests left via the default route
(ISP-A); ISP-A's edge dropped packets with an ISP-B source IP (**BCP 38** anti-
spoofing). Fix: policy routing — mark connections by inbound interface and route
replies back the same way. **Lesson:** on multi-homed hosts, the *return* path
must match the request path, or stateful devices and anti-spoofing filters will
silently eat your traffic.

**Try it yourself**
```bash
ip route                              # main table
ip route get 1.1.1.1                  # exact decision for one destination
ip rule                               # policy rules
ip route show table all               # every table
mtr -rwzbc100 example.com             # combined traceroute + loss/latency per hop
traceroute -n example.com             # UDP probes; add -I for ICMP, -T for TCP
```

---

## 14. ICMP: ping, traceroute, and Path MTU Discovery

### 14.1 What ICMP is for

**ICMP (Internet Control Message Protocol)** is IP's error-reporting and
diagnostic sidekick (IP protocol number `1`; ICMPv6 is `58`). It is **not** a
transport — no ports. An ICMP error message quotes the IP header + first 8 bytes
of the offending packet so the sender can match it to a socket.

### 14.2 Message types you must know

| Type/Code | Name | Triggered by | Seen as |
|---|---|---|---|
| 8 / 0 | Echo Request / Reply | `ping` | reachability + RTT |
| 3 / 0 | Destination unreachable: **net** | no route | "Destination Net Unreachable" |
| 3 / 1 | … : **host** | ARP fails for final hop | "Destination Host Unreachable" |
| 3 / 3 | … : **port** | UDP packet to a closed port | how `traceroute` detects the end; "connection refused" for UDP |
| 3 / 4 | **Fragmentation needed & DF set** | oversize packet, DF=1 | PMTUD signal (Ch. 8) |
| 3 / 13 | Communication administratively prohibited | firewall `REJECT` | `!X` in traceroute |
| 5 | Redirect | "use a better first-hop router" | often filtered (spoofing risk) |
| 11 / 0 | **Time exceeded in transit** (TTL hit 0) | `traceroute` | each hop reply |
| 11 / 1 | Fragment reassembly time exceeded | lost fragment | reassembly failures |
| 12 | Parameter problem | malformed header | rare |

### 14.3 How `ping` works

Sends ICMP Echo Request; the target's kernel replies with Echo Reply (no app
involved). Measures RTT and loss.

Caveats:
- Many hosts/firewalls **rate-limit or drop** ICMP → "ping fails" ≠ "host down."
- ICMP may be **deprioritized** by routers (control-plane policing) → ping
  latency/loss to a *router's own address* can look terrible while transit
  traffic through it is perfect. Always ping the **endpoint**, not an
  intermediate hop, for real numbers.
- `ping -s <size> -M do` tests PMTU (don't-fragment, set payload size).
- `ping -f` (flood) and large counts need care on production links.

### 14.4 How `traceroute` works

Send packets with **increasing TTL**: TTL 1 → first router drops it and returns
`ICMP time-exceeded` (revealing hop 1's IP); TTL 2 → hop 2; and so on until a
packet reaches the destination.

- **Unix traceroute:** UDP to high ports by default; destination replies with
  `ICMP port unreachable` (type 3/3) → done.
- **Windows `tracert`:** ICMP Echo.
- **`traceroute -T` / `tcptraceroute`:** TCP SYN to a real port (e.g. 443) —
  best for paths where UDP/ICMP is filtered but TCP/443 is open.

Reading it:
- `* * *` at a hop = that router doesn't send time-exceeded (or rate-limits it).
  **Not necessarily a problem** if later hops respond.
- Rising latency that *stays* high from some hop onward = real; a single hop with
  high latency but lower latency *after* it = that router just deprioritizes ICMP
  to itself. Ignore it.
- Different IPs per probe at the same hop = **ECMP** load balancing (Ch. 50). Use
  `traceroute` Paris-mode to keep the flow tuple constant.
- Asymmetry: traceroute shows the **forward** path only; the return path can
  differ and cause latency you can't see here.

### 14.5 `mtr` — the tool you should actually use

`mtr` runs continuous traceroute + ping, showing per-hop **loss %**, mean, and
jitter over time. Loss that appears at hop N *and every hop after* = real loss at
N. Loss at hop N only, with 0% after = cosmetic (ICMP rate limiting at N).

### 14.6 Do NOT block all ICMP

Blocking ICMPv4 entirely breaks PMTUD (Ch. 8 black hole). Blocking ICMPv6 breaks
**Neighbor Discovery** (Ch. 17) — IPv6 stops working. Minimum allow-list:
- v4: echo request/reply (optional), **type 3 (all)**, type 11.
- v6: **types 1, 2, 3, 4**, 128/129 (echo), and NDP types 133–137 on-link.

### 14.7 War story

"Traceroute shows 30% packet loss at hop 6 — the ISP has a bad router!" `mtr`
showed 30% at hop 6 and **0%** at hops 7–15 and the destination. Hop 6 was a core
router policing ICMP-to-self at 100 pps; transit traffic was untouched. The real
user-visible issue was elsewhere (an overloaded app server). **Lesson:** only
trust hop loss if it persists through every subsequent hop.

**Try it yourself**
```bash
ping -c5 -W1 example.com
ping -M do -s 1472 -c3 example.com          # PMTU probe
traceroute -n example.com
traceroute -T -p 443 -n example.com          # TCP traceroute to a real port
mtr -rwzbc100 example.com
sudo tcpdump -i any -n icmp or icmp6         # watch ICMP live
```

---

## 15. IP fragmentation and reassembly in detail

### 15.1 The mechanism (IPv4)

If a packet larger than the egress MTU must be sent and **DF=0**, the router (or
the source) splits it:

```
Original: ID=4711, len=4000 (3980 payload), DF=0
MTU=1500 -> payload chunks must be multiples of 8, <= 1480

Frag 1: ID=4711, offset=0,    MF=1, payload bytes 0-1479     (len 1500)
Frag 2: ID=4711, offset=185,  MF=1, payload bytes 1480-2959  (185*8 = 1480)
Frag 3: ID=4711, offset=370,  MF=0, payload bytes 2960-3979  (len 1040)
```

- All fragments share the **Identification**, source, destination, protocol.
- **Offset** is in 8-byte units → every fragment except the last has a payload
  length divisible by 8.
- **MF=1** on all but the last.
- Only the **first** fragment has the TCP/UDP header → stateless firewalls and
  load balancers that key on L4 ports **can't classify fragments 2..N**.
- **Reassembly happens only at the final destination**, which holds fragments in
  a buffer with a timer (`net.ipv4.ipfrag_time`, default 30 s). Miss one →
  `ICMP 11/1`, whole packet discarded, upper layer must resend everything.

### 15.2 IPv6 differs

Routers **never** fragment IPv6. Only the **source** may, using a **Fragment
extension header**. If the source doesn't do PMTUD and sends too big → the packet
is dropped with `ICMPv6 Packet Too Big`. Minimum IPv6 MTU is **1280**; a source
can always fall back to that.

### 15.3 Why fragmentation is considered harmful

1. **Loss amplification:** losing 1 of 3 fragments wastes the other 2 and forces
   a full resend. Effective loss rate triples.
2. **Firewall/LB evasion & confusion:** overlapping/tiny fragments were classic
   IDS-evasion attacks (**teardrop**, **ping of death**). Many networks now
   **drop all fragments**.
3. **Stateful device load:** reassembly needs memory and CPU; a flood of
   first-fragment-missing packets is a DoS (**fragment flood**).
4. **Offload defeat:** NIC segmentation/checksum offloads assume unfragmented
   packets.
5. **PMTUD black holes** (Ch. 8) when the ICMP is filtered.

### 15.4 What to do instead

- Let **TCP** discover MSS/PMTU and never fragment (it always sets DF).
- For **UDP** apps that send large messages (DNS over UDP, QUIC, VPN, RTP):
  - Keep datagrams ≤ ~1200 bytes to survive almost any path (QUIC uses 1200,
    DNS recommends EDNS buffer ≤ 1232).
  - Or do application-level PLPMTUD.
- On tunnels: **MSS clamp** for TCP; set a conservative tunnel MTU for UDP.
- DNS specifically: large UDP responses (DNSSEC) fragment → get dropped → clients
  retry over **TCP/53**. Ensure TCP/53 is allowed (a very common misconfig).

### 15.5 War story

A resolver behind a firewall failed to validate DNSSEC for certain zones.
`dig +dnssec` over UDP returned truncated or timed out; over TCP it worked. The
firewall dropped IP fragments, and the 2500-byte DNSKEY response fragmented at
1500. Fixes applied: allow TCP/53, set the resolver's EDNS UDP buffer to 1232,
and stop dropping fragments to the resolver. **Lesson:** "DNSSEC is broken" is
often "IP fragments are being dropped" plus "TCP/53 is blocked."

**Try it yourself**
```bash
# Force fragmentation and watch it
ping -s 4000 -c1 <peer>                 # 4000-byte ICMP -> multiple fragments
sudo tcpdump -i any -n 'ip[6:2] & 0x3fff != 0'   # capture only IPv4 fragments
sysctl net.ipv4.ipfrag_time net.ipv4.ipfrag_high_thresh
dig +dnssec +tcp DNSKEY org             # force DNS over TCP
```

---
# Part IV — The Internet Layer (IPv6)

IPv6 is not "IPv4 with more bits." Auto-configuration, the death of NAT and
broadcast, mandatory ICMPv6, and a fixed header change day-to-day operations.

## 16. IPv6 addressing and notation

### 16.1 The address

128 bits, written as eight groups of four hex digits:
`2001:0db8:0000:0000:0000:ff00:0042:8329`.

Compression rules:
1. Drop leading zeros per group: `2001:db8:0:0:0:ff00:42:8329`.
2. Replace **one** run of all-zero groups with `::`:
   `2001:db8::ff00:42:8329`. (Only once — `::` twice would be ambiguous.)
3. Hex is case-insensitive; lowercase is the norm (RFC 5952).

Prefix length is always CIDR: `2001:db8:abcd:12::/64`.

### 16.2 Structure of a typical unicast address

```
| 48 bits global routing prefix | 16 bits subnet | 64 bits interface ID |
|   2001:0db8:abcd               :  0012          :  0000:0000:0000:0001 |
|<---------- given to you by your ISP/RIR ------->|<-- you pick per host -->|
```

- The **/64** is sacred: SLAAC, many implementations, and lots of tooling assume
  every LAN is exactly a `/64`. Don't subnet longer than /64 for normal LANs.
- A site typically gets **/48** (65 536 subnets) or **/56** (256 subnets) from
  its ISP. A single home often gets a /56 or /60 via DHCPv6-PD (prefix
  delegation).
- **Interface ID** (last 64 bits): historically **EUI-64** (derived from the MAC:
  split it, insert `ff:fe`, flip the U/L bit) — but that leaks your MAC and lets
  you be tracked across networks, so modern hosts use **RFC 8981 temporary
  privacy addresses** (random, rotated daily) for outbound and often a stable
  random address for inbound.

### 16.3 Address types (there is no broadcast)

| Prefix | Type | Notes |
|---|---|---|
| `::1/128` | Loopback | equals IPv4 `127.0.0.1` |
| `::/128` | Unspecified | "no address yet" (used as source during DAD) |
| `fe80::/10` | **Link-local** | Auto-configured on **every** IPv6 interface, always present, never routed. Used by NDP, OSPFv3, and as next-hop. You must specify the zone: `ping6 fe80::1%eth0`. |
| `fc00::/7` (`fd00::/8` in practice) | **Unique Local (ULA)** | The IPv6 "private" range; generate a random /48. Not globally routed. |
| `2000::/3` | **Global unicast (GUA)** | Internet-routable. |
| `ff00::/8` | **Multicast** | Replaces broadcast. Scoped: `ff02::1` all-nodes (link), `ff02::2` all-routers, `ff02::fb` mDNS, `ff02::1:ffXX:XXXX` **solicited-node** (used by NDP). |
| `2002::/16` | 6to4 | legacy transition, deprecated |
| `2001:db8::/32` | Documentation | safe for examples |
| `64:ff9b::/96` | NAT64 well-known prefix | maps IPv4 into IPv6 (Ch. 18) |
| `::ffff:0:0/96` | IPv4-mapped | `::ffff:192.0.2.1` — how a dual-stack socket represents an IPv4 peer |

An interface normally has **several** addresses at once: one `fe80::` link-local,
maybe a ULA, one or more GUAs (stable + temporary). This is normal, not a
misconfig.

### 16.4 The IPv6 header (fixed 40 bytes)

```
|Version| Traffic Class |           Flow Label                  |
|         Payload Length          | Next Header |   Hop Limit    |
|                     Source Address (128 bits)                  |
|                  Destination Address (128 bits)                |
```

Differences from IPv4 that matter:
- **No header checksum** (link layer + transport checksums suffice; and IPv6
  requires the UDP checksum, unlike IPv4).
- **No IHL / options in the base header** — extension headers chain via **Next
  Header** (Hop-by-Hop, Routing, Fragment, ESP, AH, Destination). Firewalls
  sometimes drop long extension-header chains.
- **Hop Limit** = IPv4's TTL, renamed.
- **Flow Label** (20 bits): lets routers ECMP-hash a flow consistently without
  parsing L4; also used by QUIC.
- **No fragmentation fields** — moved to an extension header, source-only.

### 16.5 War story

An app "worked on IPv4, hung on IPv6" for some users. `getaddrinfo` returned a
AAAA, the app connected to a GUA, but the user's ISP had broken IPv6 (RA present,
no working upstream). The app didn't implement **Happy Eyeballs** (Ch. 18) so it
waited a full 75 s connect timeout before trying IPv4. Fix: use a client library
with Happy Eyeballs v2, or lower the connect timeout and race both families.
**Lesson:** in a dual-stack world, always race A and AAAA; never let one broken
family block the other.

**Try it yourself**
```bash
ip -6 addr show                         # all v6 addresses (note multiple per iface)
ip -6 route
ping6 -c3 2606:4700:4700::1111           # Cloudflare v6
ping6 fe80::1%eth0                        # link-local needs a zone id
sysctl net.ipv6.conf.all.use_tempaddr    # 2 = prefer privacy addresses
```

---

## 17. Neighbor Discovery, SLAAC, and DHCPv6

### 17.1 ICMPv6 Neighbor Discovery (NDP) replaces ARP + ICMP redirects + more

NDP (RFC 4861) runs over ICMPv6 and uses multicast instead of broadcast:

| Message | Type | Role |
|---|---|---|
| **Router Solicitation (RS)** | 133 | Host boot: "any routers here?" (to `ff02::2`) |
| **Router Advertisement (RA)** | 134 | Router announces prefix(es), MTU, flags, its own presence; sent periodically to `ff02::1` and on demand |
| **Neighbor Solicitation (NS)** | 135 | "Who has `2001:db8::5`?" (to the solicited-node multicast `ff02::1:ff00:0005`) — the ARP-request equivalent, plus **DAD** |
| **Neighbor Advertisement (NA)** | 136 | "That's me, at MAC …" — the ARP-reply equivalent |
| **Redirect** | 137 | "Use a better first hop" |

Because NDP *is* ICMPv6, **blocking ICMPv6 breaks IPv6 entirely** — no address
resolution, no router discovery.

### 17.2 SLAAC — StateLess Address AutoConfiguration

The default way IPv6 hosts get addresses; no server, no state:

1. Host creates a tentative **link-local** `fe80::<interface-id>`.
2. **DAD:** sends an NS for its own tentative address; if an NA comes back,
   there's a conflict → give up that address.
3. Host sends **RS**; router replies **RA** containing a `/64` prefix + flags.
4. Host forms a **GUA** = prefix + interface ID (random privacy ID + a stable
   one), runs DAD on each.
5. Default route = the RA's source **link-local** address.

RA flags decide what else is needed:

| Flag | Meaning |
|---|---|
| **A** (autonomous) | Use this prefix for SLAAC |
| **M** (managed) | Get your address from **DHCPv6** (stateful) instead/also |
| **O** (other config) | Use DHCPv6 only for *other* info (DNS, NTP) — "stateless DHCPv6" |
| **L** (on-link) | Hosts on this prefix are directly reachable |

Also: RDNSS/DNSSL options in the RA can carry DNS servers directly (RFC 8106), so
pure SLAAC networks can work with no DHCPv6 at all — *if* the client supports
RDNSS (Android historically did not, forcing RA+RDNSS or... arguments).

### 17.3 DHCPv6

- **Stateful (M=1):** server hands out specific addresses + options over
  UDP `546`(client)/`547`(server), using multicast `ff02::1:2`. Clients are keyed
  by **DUID** (not MAC), which trips up MAC-based reservations.
- **Stateless (O=1):** address via SLAAC, but DNS/NTP/domain via DHCPv6.
- **Prefix Delegation (DHCPv6-PD):** ISP delegates e.g. a `/56` to your router,
  which then sub-allocates `/64`s to your LANs via RA. This is how a home router
  gets routable space for every device without NAT.

### 17.4 No NAT, but still a firewall

IPv6 hosts usually have globally routable addresses. Security comes from a
**stateful firewall** (default-deny inbound, allow established/related), *not*
from address hiding. Home routers do this by default. Privacy addresses limit
tracking.

### 17.5 Rogue RA — the IPv6 equivalent of a rogue DHCP server

Any device can send RAs. A misconfigured host (Windows Internet Connection
Sharing, a hypervisor, a phone) advertising itself as a router will black-hole a
whole LAN's IPv6. Defenses: **RA Guard** on switches, and
`net.ipv6.conf.*.accept_ra` = 0 on servers that shouldn't listen.

### 17.6 War story

After enabling IPv6 on a corporate LAN, random Windows laptops lost internet for
minutes at a time. A developer's VM was periodically emitting RAs with a bogus
prefix and a high priority. Hosts preferred it, installed a default route to
nowhere, and only recovered when that RA aged out. Enabling **RA Guard** on the
access switches ended it. **Lesson:** treat RAs like DHCP — exactly one
authorized source, enforced at the switch.

**Try it yourself**
```bash
sudo tcpdump -i eth0 -n 'icmp6 && ip6[40] >= 133 && ip6[40] <= 137'   # NDP/RA/RS
ip -6 neigh                              # the IPv6 neighbor cache (ARP equivalent)
rdisc6 eth0                              # actively solicit and print RAs
sysctl net.ipv6.conf.eth0.accept_ra
radvdump                                 # decode received RAs in detail
```

---

## 18. Dual-stack, Happy Eyeballs, and transition mechanisms

### 18.1 Dual-stack

A host runs IPv4 and IPv6 simultaneously; each connection uses one or the other.
`getaddrinfo` returns both A and AAAA; the resolver/OS orders them (RFC 6724 —
generally "IPv6 first, but prefer the family with a matching source address
scope").

### 18.2 Happy Eyeballs (RFC 8305)

Naively "try AAAA, on failure try A" means a broken-IPv6 user eats a full TCP
timeout (tens of seconds) on every connection. **Happy Eyeballs v2**:

1. Start the AAAA and A DNS lookups in parallel; don't wait for both.
2. Start a TCP connect to the first address (prefer IPv6) as soon as it's known.
3. If it hasn't connected within a short delay (~250 ms), **start a parallel
   connect to an address of the other family**.
4. First connection to complete wins; cancel the rest.
5. Cache the winner's family per destination for a while.

Result: broken IPv6 costs ~250 ms once, not 30 s every time. Browsers, `curl`,
Go's dialer, and modern libc all do this. If you write raw socket code, do it
too.

### 18.3 Transition mechanisms (know the names)

| Mechanism | What it does | Where you meet it |
|---|---|---|
| **Dual-stack** | run both | the "correct" answer, when feasible |
| **NAT64 + DNS64** | IPv6-only client → IPv4 server. DNS64 synthesizes a AAAA in `64:ff9b::/96` from the A record; NAT64 gateway translates `64:ff9b::x.y.z.w` ↔ IPv4 | mobile carriers (T-Mobile US), IPv6-only cloud subnets, `iOS` app-store requirement |
| **464XLAT** | adds a client-side CLAT so IPv4-only *apps* work over an IPv6-only + NAT64 network | Android on IPv6-only cellular |
| **DS-Lite** | IPv4 client → IPv4 internet over an IPv6-only access network, with CGNAT at the ISP (AFTR) | some ISPs |
| **MAP-E / MAP-T** | stateless IPv4-over-IPv6 with per-CPE port ranges | some ISPs |
| **6in4 / 6to4 / Teredo / ISATAP** | tunnel IPv6 over IPv4 | mostly legacy; 6in4 tunnels (HE.net) still used by hobbyists |
| **NPTv6** | 1:1 stateless prefix translation (not NAT/PAT) | multihoming without PI space |

### 18.4 Operational gotchas

- **NAT64/DNS64 breaks IPv4 literals and DNSSEC validation on the client** —
  apps that hardcode `1.2.3.4` won't work without CLAT.
- **`localhost` may resolve to `::1` first** — a service bound only to `127.0.0.1`
  will refuse IPv6 loopback connections. Bind to both or to `::` (with
  `bindv6only=0`).
- **Firewall rules must be duplicated** for `ip6tables`/`nftables inet`. A common
  breach: v4 locked down, v6 wide open.
- **Logging/allowlists** must handle 128-bit addresses and the fact that one host
  rotates through many privacy addresses.

### 18.5 War story

An API allowlist keyed on client IP started rejecting a partner intermittently.
The partner was on an IPv6-first network; some requests came from a stable GUA
(allowlisted) and some from **rotating temporary privacy addresses** (not).
Solutions considered: allowlist the partner's whole `/64`, have them use a
stable-address egress proxy, or switch to mutual TLS auth. **Lesson:** IPv6
source addresses are many-per-host and time-varying by design; don't build
security on "the client IP."

**Try it yourself**
```bash
getent ahosts example.com               # see A and AAAA and the resolver's ordering
curl -v --connect-timeout 5 https://example.com    # watch which family wins
curl -4 ...  /  curl -6 ...              # force a family
python3 -c "import socket; print(socket.getaddrinfo('example.com', 443))"
```

---
# Part V — Host configuration

## 19. DHCP, from DISCOVER to lease renewal

### 19.1 The DORA exchange (IPv4)

DHCP runs over UDP: client port **68**, server port **67**. Before it has an
address, the client uses source `0.0.0.0` and destination `255.255.255.255`
(limited broadcast).

```
Client                                Server(s)
  |  DHCPDISCOVER (bcast) ------------>|   "anyone? my MAC is aa:bb:.., xid=0x1f2e"
  |<------------ DHCPOFFER (bcast/uni) |   "you can have 192.168.1.50, mask /24,
  |                                    |    gw .1, dns .1, lease 86400s, from me"
  |  DHCPREQUEST (bcast) ------------->|   "I accept 192.168.1.50 from server X"
  |<------------- DHCPACK              |   "confirmed; here are all options"
  |                                    |
  | (client runs ARP/NS for its new    |
  |  address = duplicate detection)    |
```

- **DISCOVER** is broadcast because the client has no address and no idea who the
  server is.
- Multiple servers may **OFFER**; the client picks one (usually the first) and
  names it in the REQUEST via the *Server Identifier* option; the others withdraw.
- **DHCPNAK** = "that address isn't valid (any more) on this segment" → client
  restarts from DISCOVER. Common after moving between networks.
- **DHCPDECLINE** = client's post-assignment ARP/DAD found the address already in
  use → back to DISCOVER, server marks it bad.
- **DHCPRELEASE** = graceful give-back (many clients skip it).
- **DHCPINFORM** = "I have a static IP, just give me options (DNS, etc.)."

### 19.2 The lease timers: T1, T2, and expiry

```
0                         T1 (50%)        T2 (87.5%)     Lease end (100%)
|---------------------------|---------------|---------------|
       normal use           |               |               |
                            |  RENEW (unicast to same server)
                                            |  REBIND (broadcast to ANY server)
                                                            |  address released,
                                                               back to DISCOVER
```

- At **T1** (default 50% of lease) the client **renews** by unicasting a
  REQUEST to its server. Usually silent and instant.
- If that fails, at **T2** (87.5%) it **rebinds** — broadcasts a REQUEST to any
  server.
- If nothing answers by lease end, the client **must stop using the address**.
  A well-behaved client then falls to `169.254.x.x` link-local (APIPA) — the
  "I have limited connectivity / self-assigned IP" state.

### 19.3 Key options (the numbers show up in captures)

| Option | Name | Notes |
|---|---|---|
| 1 | Subnet mask | |
| 3 | Router (default gateway) | |
| 6 | DNS servers | |
| 12 / 15 | Hostname / Domain name | |
| 42 | NTP servers | |
| 51 | Lease time | |
| 53 | **DHCP message type** | 1=DISCOVER 2=OFFER 3=REQUEST 5=ACK 6=NAK … |
| 54 | Server identifier | which server this is from |
| 55 | **Parameter Request List** | the client's wishlist; also an OS fingerprint |
| 82 | Relay Agent Information | inserted by the relay; carries the ingress port/circuit — used for IP-per-port policy and DHCP snooping |
| 119 | Domain search list | |
| 121 / 249 | Classless static routes | push extra routes via DHCP (249 = MS variant) |

### 19.4 DHCP across subnets: the relay

The DHCP server is usually not on the client's broadcast domain. The router runs
a **DHCP relay** (`ip helper-address` on Cisco, `dhcp-relay` elsewhere): it
catches the broadcast DISCOVER, rewrites it as unicast to the real server, stamps
the **giaddr** (gateway IP) so the server knows which subnet/pool to use, and
relays the answer back. This is why one central DHCP server can serve 200 VLANs.

### 19.5 Failure modes

- **Rogue DHCP server** (someone plugs in a home router): hands out its own IP as
  gateway → MITM or black hole. Defense: **DHCP snooping** on switches (only
  "trusted" ports may send OFFER/ACK) + option 82.
- **Pool exhaustion:** short leases + churn (conference Wi-Fi) or a client
  requesting many addresses. Symptom: new clients get APIPA. Fix: bigger pool,
  shorter lease, or per-MAC limits.
- **Lease/identity mismatch after cloning a VM:** two VMs with the same MAC fight
  over one lease → flapping. Regenerate MACs.
- **DHCPv6 uses DUID, not MAC:** a re-imaged host presents a new DUID and gets a
  *new* address, orphaning your reservation.

### 19.6 War story

Every Monday, a branch office's laptops took 5 minutes to get online. The DHCP
lease was 8 days and the pool was sized for the ~40 daily devices; over a busy
week ~90 distinct devices (visitors, phones) consumed leases that didn't expire
until after the weekend. Monday's staff hit an exhausted pool and fell to APIPA
until leases aged out. Fix: 8-hour lease for the guest VLAN, separate pool.
**Lesson:** lease time × device churn must be ≤ pool size, with margin.

**Try it yourself**
```bash
sudo tcpdump -i eth0 -n -v 'udp port 67 or udp port 68'   # watch DORA
# Linux (systemd):
networkctl status eth0
# dhclient with verbose lease info:
sudo dhclient -v eth0
cat /var/lib/dhcp/dhclient.leases        # or /var/lib/NetworkManager/*.lease
nmcli -f DHCP4 con show "Wired connection 1"
```

---

## 20. DNS: resolution, records, caching, and modern transports

DNS turns names into addresses (and much more). It is a distributed, cached,
hierarchical database. Almost every "the internet is slow/broken" incident has a
DNS chapter.

### 20.1 The name hierarchy

```
.                                   <- root
 +-- com.                           <- TLD (top-level domain)
      +-- example.com.              <- second-level (a "zone" you control)
           +-- www.example.com.     <- a record within the zone
           +-- api.example.com.
```

The trailing dot = the root; `www.example.com` and `www.example.com.` are the
same. Each zone is authoritatively served by a set of **name servers** and can
**delegate** sub-zones (e.g. `example.com` delegates `eng.example.com` to another
team's servers via `NS` records + glue).

### 20.2 The players

| Role | What it does |
|---|---|
| **Stub resolver** | The tiny client in your OS/app (`getaddrinfo`). Asks *one* recursive resolver, caches little or nothing. |
| **Recursive resolver** (a.k.a. caching resolver) | Does the legwork: root → TLD → authoritative, caches every answer by TTL, returns the final answer. Examples: your ISP's, `1.1.1.1`, `8.8.8.8`, a corporate `unbound`/`bind`. |
| **Authoritative name server** | Holds the actual zone data for domains it's responsible for. Never recurses. Examples: Route 53, NS1, Cloudflare, your own `bind`. |
| **Forwarder** | A resolver that just relays to another resolver (home routers do this). |

### 20.3 A full recursive resolution (cold cache)

Client wants `A` for `shop.example.com`:

```
stub -> recursive:  "A shop.example.com?"  (recursion desired)

recursive -> a ROOT server:      "A shop.example.com?"
root      -> recursive:          "I don't know, but ask .com servers: NS a.gtld-servers.net (+ glue A/AAAA)"

recursive -> a .COM server:      "A shop.example.com?"
.com      -> recursive:          "ask example.com's servers: NS ns1.example.com (+ glue)"

recursive -> ns1.example.com:    "A shop.example.com?"
ns1       -> recursive:          "A shop.example.com = 93.184.216.34, TTL 300"  (authoritative, AA bit set)

recursive caches every referral + the final answer, then:
recursive -> stub:               "93.184.216.34, TTL 300"
```

Next lookup of anything in `.com` skips the root; next lookup in `example.com`
skips `.com`; next lookup of `shop.example.com` within 300 s is served straight
from cache. This is why the *first* request to a new domain is slow and the rest
are instant.

### 20.4 Record types you must know

| Type | Purpose | Example |
|---|---|---|
| **A** | name → IPv4 | `example.com. 300 IN A 93.184.216.34` |
| **AAAA** | name → IPv6 | `example.com. 300 IN AAAA 2606:2800:220:1:...` |
| **CNAME** | alias → canonical name | `www.example.com. IN CNAME example.com.` — **cannot** coexist with other records at the same name, so never at the zone apex |
| **MX** | mail servers + priority | `example.com. IN MX 10 mail.example.com.` |
| **NS** | delegation: who's authoritative | `example.com. IN NS ns1.example.com.` |
| **SOA** | zone metadata: primary NS, admin email, serial, refresh/retry/expire, **minimum (negative-cache) TTL** | one per zone |
| **TXT** | arbitrary text | SPF, DKIM, domain-verification tokens |
| **PTR** | IP → name (reverse DNS) | in `in-addr.arpa` / `ip6.arpa`; used by mail servers, logging |
| **SRV** | service location: host + port + priority + weight | `_sip._tcp.example.com. IN SRV 10 60 5060 sipserver...` |
| **CAA** | which CAs may issue certs for this domain | issuance policy |
| **DS / DNSKEY / RRSIG / NSEC(3)** | DNSSEC chain of trust | signatures + delegation-signer + authenticated denial |
| **SVCB / HTTPS** | service binding: ALPN, port, IP hints, ECH config — lets a client learn "use HTTP/3 on port 443 at these IPs" in one lookup | `example.com. IN HTTPS 1 . alpn=h3,h2 ipv4hint=...` |
| **CNAME at apex workarounds** | `ALIAS`/`ANAME`/CNAME-flattening | provider-specific |

### 20.5 TTL and caching — the operational core

- Every record has a **TTL**; resolvers must not serve it longer than that.
- **Lowering TTL takes effect only after the *old* TTL expires.** To migrate an
  IP with minimal downtime: drop the TTL to 60 s *at least the old-TTL duration
  in advance*, then change the record, then raise the TTL back afterward.
- **Negative caching:** "no such name" (NXDOMAIN) and "no such record" are cached
  too, for the SOA `minimum` field's duration. A typo'd record that returns
  NXDOMAIN can haunt you for hours.
- Resolvers and OSes and browsers **each** cache. `dig` bypasses the OS/browser
  cache; your app may not.
- Some resolvers **clamp** absurd TTLs (cap at a day, floor at a few seconds) and
  **prefetch** popular records before expiry.

### 20.6 Transport: UDP first, TCP fallback, then encrypted

- Classic DNS: **UDP/53**, one query per packet. If the response is too big for
  the client's advertised **EDNS0** buffer (or > 512 bytes without EDNS), the
  server sets the **TC (truncated)** bit and the client **retries over TCP/53**.
  **You must allow TCP/53** — DNSSEC, large TXT, and many CDNs need it.
- **EDNS0**: an OPT pseudo-record that advertises a larger UDP buffer (commonly
  1232 today, to dodge fragmentation — Ch. 15), plus flags and options
  (**ECS** = EDNS Client Subnet, which passes a truncated client prefix to
  authoritative servers for geo-accurate answers).
- **DoT (DNS over TLS)**: port **853**. DNS inside TLS. Easy to identify/block by
  port.
- **DoH (DNS over HTTPS)**: port **443**, looks like normal web traffic. Used by
  browsers (often bypassing your OS resolver and your network's DNS policy —
  a real operational headache: your split-horizon/internal names stop resolving).
- **DoQ (DNS over QUIC)**: port 853/UDP, no head-of-line blocking.
- **DNSSEC** ≠ encryption. It's *authentication/integrity*: authoritative servers
  sign records (`RRSIG`), a chain of `DS`→`DNSKEY` from the root lets a
  validating resolver detect tampering. It does not hide queries.

### 20.7 Split-horizon / split-brain DNS

The same name resolves differently depending on who asks: internal clients get
`10.x` (private), external clients get the public IP. Implemented with views
(BIND), separate resolvers, or "internal" zones. Breaks when: a laptop uses DoH
and bypasses the internal resolver; a VPN doesn't push the internal search
domain; caching carries an internal answer to the coffee shop.

### 20.8 Resolution details that bite

- **`search` domains + `ndots`:** `resolv.conf` `search corp.example.com` and
  `options ndots:5` mean `db` becomes `db.corp.example.com`, and even
  `foo.bar.baz.qux` (4 dots < 5) gets search suffixes tried **first**. In
  Kubernetes, default `ndots:5` makes every external lookup do 4–5 failed
  queries first — a well-known latency/QPS amplifier. Fix: `dnsConfig` with
  `ndots:1` or use FQDNs with a trailing dot.
- **`/etc/hosts` and NSS order** (`/etc/nsswitch.conf`: `hosts: files dns`) win
  before DNS.
- **`getaddrinfo` vs `gethostbyname`:** the former is dual-stack and applies
  RFC 6724 sorting; use it.
- **Happy Eyeballs** (Ch. 18) races the A/AAAA connects, not just the lookups.
- **One resolver dead in `resolv.conf`:** default behavior tries them in order
  with a ~5 s timeout each → every lookup stalls 5 s. `options rotate` /
  `timeout:1` / `attempts:2` help; better, run a local caching resolver.

### 20.9 War story

A deploy "broke DNS" for a service mesh: p50 request latency jumped 150 ms.
`tcpdump` on port 53 showed **five** queries per outbound call —
`svc`, `svc.ns.svc.cluster.local`, `svc.svc.cluster.local`,
`svc.cluster.local`, `svc.example.com` — the first four NXDOMAIN. Root cause:
code used a bare short name with the default Kubernetes `ndots:5`. Fix: append
the trailing dot to make it a FQDN (`svc.ns.svc.cluster.local.`), cutting it to
one query. **Lesson:** `ndots` + `search` can multiply your DNS load 5×
invisibly; always check with a capture.

**Try it yourself**
```bash
dig +trace example.com                       # watch the full root->TLD->auth walk
dig +norecurse @a.root-servers.net example.com   # ask a root server directly
dig example.com A +short
dig example.com HTTPS                          # SVCB/HTTPS record
dig -x 93.184.216.34                           # reverse (PTR)
dig +dnssec example.com                        # see RRSIG; 'ad' flag = validated
dig @1.1.1.1 example.com                        # bypass local resolver
resolvectl status                              # systemd-resolved: per-link servers, DNSSEC, DoT
grep -E 'search|nameserver|options' /etc/resolv.conf
# Measure resolver latency:
for i in $(seq 5); do dig +noall +stats example$i.test @1.1.1.1 | grep 'Query time'; done
```

---
# Part VI — The Transport Layer: UDP

## 21. Ports, sockets, and the connection 4-tuple

### 21.1 Ports

A **port** is a 16-bit number (0–65535) that identifies which application
endpoint on a host a segment belongs to. It's the demultiplexing key at Layer 4.

| Range | Name | Use |
|---|---|---|
| 0–1023 | **Well-known / system** | Binding usually needs root/`CAP_NET_BIND_SERVICE`. 22 SSH, 53 DNS, 80 HTTP, 443 HTTPS, 123 NTP, 25 SMTP. |
| 1024–49151 | **Registered** | App vendors register these (3306 MySQL, 6379 Redis, 5432 Postgres). Not enforced. |
| 49152–65535 | **Dynamic / ephemeral** (IANA) | Source ports for outbound connections. Linux actually uses `net.ipv4.ip_local_port_range`, default `32768 60999`. |

### 21.2 The socket and the 4-tuple

A TCP connection is uniquely identified by the **4-tuple**:
```
(source IP, source port, destination IP, destination port)
```
Plus the protocol, it's a **5-tuple**. The kernel uses it as the hash key to
find the right connection for every incoming segment.

Consequences:
- One server `:443` can hold **millions** of simultaneous connections — they
  differ in the client IP/port, not the server side.
- A single client can open at most ~28 000 connections **to the same
  (dst IP, dst port)** before exhausting its ephemeral range — then `connect()`
  returns `EADDRNOTAVAIL`. More client IPs, or more destination IPs/ports, or
  `SO_REUSEADDR`/`IP_BIND_ADDRESS_NO_PORT` extend this.
- **`SO_REUSEPORT`**: multiple listening sockets (one per worker thread/process)
  can bind the *same* port; the kernel load-balances incoming connections across
  them. This is how nginx/HAProxy/Envoy scale across cores without a single
  accept lock.

### 21.3 Listening vs connected sockets

```
$ ss -tanp
State    Recv-Q Send-Q  Local Address:Port   Peer Address:Port
LISTEN   0      511     0.0.0.0:443          0.0.0.0:*          # accepts new conns
ESTAB    0      0       10.0.0.5:443         203.0.113.9:54321  # one client
ESTAB    0      0       10.0.0.5:443         203.0.113.9:54322  # same client, diff port
```

- **`LISTEN`** socket: bound to `(local IP or wildcard, port)`. `Send-Q` here =
  the **accept queue** (`backlog`) capacity; `Recv-Q` = completed connections
  waiting for `accept()`.
- Each `accept()` returns a **new** socket bound to the full 4-tuple.
- `backlog` (from `listen(fd, backlog)`) is capped by `net.core.somaxconn`
  (default 4096 on modern Linux). If it overflows, new SYNs are dropped →
  clients retry after 1 s → visible "sometimes takes 1s to connect."

### 21.4 War story

A payments service capped out at ~28 000 req/s to its downstream fraud API and
threw `cannot assign requested address`. It opened a fresh connection per request
to a single `fraud-api:443` VIP, exhausting ephemeral ports faster than
`TIME_WAIT` (60 s) cleared them. Fix: an HTTP connection pool with keep-alive
(reuse a few hundred connections). Throughput went to 200 k req/s.
**Lesson:** "connection per request" doesn't scale; the 4-tuple space to a single
destination is only ~28k wide.

**Try it yourself**
```bash
sysctl net.ipv4.ip_local_port_range net.core.somaxconn
ss -s                                    # summary: sockets by state
ss -tan state time-wait | wc -l          # how many TIME_WAIT
ss -tlnp                                  # listening TCP sockets + owning process
cat /proc/sys/net/ipv4/tcp_max_syn_backlog
```

---

## 22. UDP: the header, and when "unreliable" is the right choice

### 22.1 The entire header (8 bytes)

```
 0                   1                   2                   3
+--------+--------+--------+--------+--------+--------+--------+--------+
|      Source Port (16)             |    Destination Port (16)         |
+--------+--------+--------+--------+--------+--------+--------+--------+
|      Length (16, header+data)     |       Checksum (16)             |
+--------+--------+--------+--------+--------+--------+--------+--------+
```

That's it. No sequence numbers, no ACKs, no connection, no flow/congestion
control, no ordering. UDP adds exactly three things to raw IP:
1. **Ports** (application demultiplexing).
2. A **checksum** over pseudo-header + UDP header + data (optional in IPv4 —
   `0` means "not computed"; **mandatory in IPv6**).
3. A **length** field.

`sendto()` of N bytes → exactly one UDP datagram → (usually) one IP packet. No
merging, no splitting (unless IP fragments it). If it's lost, it's gone — the app
finds out only if it built its own detection.

### 22.2 "Connected" UDP sockets

You *can* `connect()` a UDP socket. It doesn't send anything; it just fixes the
default peer so you can use `send()/recv()` and — importantly — the kernel will
deliver **ICMP port-unreachable** errors back to that socket (as `ECONNREFUSED`
on the next syscall). An unconnected UDP socket silently swallows those.

### 22.3 When UDP is the right call

| Use case | Why UDP wins |
|---|---|
| **DNS queries** | one tiny request/response; a lost query is just re-asked; connection setup would triple the cost |
| **Real-time audio/video (RTP, WebRTC)** | a 200 ms-late retransmitted voice packet is useless — better to skip it; app does its own loss concealment and jitter buffer (Ch. 47) |
| **Live game state** | only the *latest* position matters; TCP's in-order guarantee would stall on an old lost packet (**head-of-line blocking**) |
| **QUIC / HTTP/3** | wants reliability + encryption + streams but implemented **in userspace** so it can evolve without kernel/middlebox upgrades; UDP is just the substrate (Ch. 39) |
| **DHCP, TFTP, SNMP, NTP, syslog** | simple, often broadcast/one-shot, on constrained devices |
| **Multicast / broadcast** | TCP is strictly point-to-point; UDP can go one-to-many (IPTV, service discovery, mDNS) |
| **VPN transports (WireGuard, OpenVPN-UDP)** | avoid **TCP-over-TCP meltdown**: two nested reliable transports fight each other's retransmit timers and collapse under loss |

### 22.4 What your UDP app must handle itself (if it needs them)

- **Loss detection & retransmission** (sequence numbers + ACKs, or FEC).
- **Ordering / reassembly** of multi-datagram messages.
- **Congestion control** — ethically mandatory. Uncontrolled UDP floods are how
  you cause collateral damage and get null-routed. QUIC implements NewReno/CUBIC/
  BBR; WebRTC uses **GCC** (Google Congestion Control); use a library, don't wing
  it.
- **Path MTU** — keep datagrams small (≤1200) or do PLPMTUD; you won't get the
  friendly "connects then hangs" — you'll just get silent drops.
- **NAT keepalives** — UDP NAT mappings expire in as little as 30 s; send a
  keepalive every ~15–25 s or inbound packets stop (Ch. 46).
- **Amplification hygiene** — never answer a small spoofed request with a big
  response to an unverified source (Ch. 54). DNS, NTP `monlist`, memcached, and
  SSDP became infamous DDoS amplifiers this way.

### 22.5 UDP framing: one datagram = one message

Unlike TCP (a byte stream where you must delimit messages yourself), each UDP
`recvfrom()` returns exactly one datagram. If your buffer is smaller than the
datagram, the rest is **discarded** (`MSG_TRUNC`). A zero-length UDP datagram is
valid and *is* delivered (useful as a ping/keepalive).

### 22.6 War story

A metrics agent shipped counters over UDP to a collector "for low overhead."
During incidents — exactly when metrics matter most — the collector's socket
receive buffer (`net.core.rmem_max`) overflowed under the surge and the kernel
dropped ~40% of datagrams silently (`netstat -su` → "receive buffer errors").
Dashboards under-reported the outage. Fixes: raise `SO_RCVBUF`, batch with
`recvmmsg()`, sample, and add sequence numbers to *measure* loss.
**Lesson:** UDP loss is invisible unless you instrument it; check
`nstat`/`netstat -su` for `InErrors`/`RcvbufErrors`.

**Try it yourself**
```bash
netstat -su                              # UDP stats: InErrors, RcvbufErrors, NoPorts
nstat -az | grep -i udp
ss -uap                                  # UDP sockets + Recv-Q backlog
# Quick UDP echo test:
nc -u -l 9999            # terminal 1
nc -u 127.0.0.1 9999    # terminal 2, type lines
sysctl net.core.rmem_max net.core.rmem_default
```

---
# Part VII — The Transport Layer: TCP mechanics

TCP gives an application a **reliable, ordered, bidirectional byte stream** over
unreliable IP. "Byte stream" means TCP does **not** preserve your message
boundaries — 3 writes of 100 bytes may arrive as one `read()` of 300 or five
reads of 60. Delimiting messages is the application's job (length prefixes,
newlines, HTTP `Content-Length`, etc.).

## 23. The TCP header, field by field

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Source Port          |       Destination Port        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                        Sequence Number                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Acknowledgment Number                       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| Data  |Rsvd |C|E|U|A|P|R|S|F|                                  |
| Offset|     |W|C|R|C|S|S|Y|I|            Window Size           |
|       |     |R|E|G|K|H|T|N|N|                                  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           Checksum            |         Urgent Pointer        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Options (0-40 bytes)                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Size | Meaning / real-world relevance |
|---|---|---|
| **Source / Dest Port** | 16 b each | With the IPs, forms the 4-tuple. |
| **Sequence Number** | 32 b | Byte offset of the **first data byte** in this segment within the stream. On the SYN it carries the **ISN** (Initial Sequence Number, randomized — see 24.4). |
| **Acknowledgment Number** | 32 b | "I have received every byte up to (ack−1); send me byte `ack` next." **Cumulative.** Valid only when ACK flag set (which is always, after the SYN). |
| **Data Offset** | 4 b | Header length in 32-bit words; 5 = 20 bytes, 15 = 60 max. Tells where data starts. |
| **Reserved** | 3 b | Must be 0. |
| **CWR / ECE** | 1 b each | **ECN** signaling: ECE = "I saw a CE-marked packet, slow down"; CWR = "I reduced my window." Works with the IP-header ECN bits (Ch. 10, 31). |
| **URG + Urgent Pointer** | 1 b + 16 b | "Urgent data" up to the pointer. Effectively **dead** (Telnet `Ctrl-C` heritage); often stripped by middleboxes; a minor security footgun. |
| **ACK** | 1 b | Ack number is meaningful. Set on everything except the initial SYN. |
| **PSH** | 1 b | "Push buffered data to the app now." Mostly advisory today. |
| **RST** | 1 b | **Reset**: abort the connection immediately, no graceful close. Sent for "connection refused" (SYN to a closed port), for data on a nonexistent connection, or by an app calling `close()` with `SO_LINGER 0`, or by a firewall forging it. |
| **SYN** | 1 b | Synchronize sequence numbers — connection setup. Consumes one sequence number. |
| **FIN** | 1 b | "I'm done sending." Graceful half-close. Consumes one sequence number. |
| **Window Size** | 16 b | **Receive window (rwnd)**: how many more bytes the *sender of this segment* is willing to receive right now (flow control, Ch. 30). Max 65535 **unless window scaling** multiplies it. |
| **Checksum** | 16 b | Covers pseudo-header (src/dst IP, protocol, TCP length) + TCP header + data. Mandatory. NICs usually offload it (Ch. 49) — so a local capture may show "checksum incorrect (offloaded)". |
| **Urgent Pointer** | 16 b | see URG. |

### 23.1 The options that actually matter (negotiated on the SYN)

| Option | Kind | Purpose |
|---|---|---|
| **MSS** | 2 | Max Segment Size this host will accept (Ch. 8). Each side states its own; effective = min. Typically 1460 (Ethernet) or 1448 (with timestamps). |
| **Window Scale (WScale)** | 3 | Left-shift the 16-bit window by 0–14 bits → windows up to ~1 GB. **Must be on both SYNs** or scaling is off for the whole connection. Essential for high bandwidth-delay paths (Ch. 35). |
| **SACK-Permitted** | 4 | "I support Selective ACK." Enables SACK blocks in later segments (Ch. 29). |
| **SACK** | 5 | In data segments: "I have bytes A–B and C–D but I'm missing the gap" — lets the sender retransmit only the holes. |
| **Timestamps (TSopt)** | 8 | Two 32-bit values (TSval/TSecr). Used for accurate **RTT measurement** on every segment (RTTM) and **PAWS** (Protection Against Wrapped Sequences) on fast links. Costs 12 bytes/segment (10 + 2 padding). |
| **MP_CAPABLE / MP_JOIN** | 30 | **Multipath TCP** — one logical connection across multiple paths (Wi-Fi + cellular). |
| **TCP Fast Open (TFO)** | 34 | Send data *in* the SYN on repeat visits (a cookie authorizes it), saving one RTT. Patchy middlebox support. |

### 23.2 Reading flags in a capture

`tcpdump` shows flags as letters in `[ ]`:

| Symbol | Flag |
|---|---|
| `S` | SYN |
| `S.` | SYN+ACK (`.` = ACK) |
| `.` | ACK only |
| `P.` | PSH+ACK (data) |
| `F.` | FIN+ACK (graceful close) |
| `R` / `R.` | RST / RST+ACK (abort) |
| `E`,`W` | ECE, CWR |

Example line:
```
10:00:00.100 IP 10.0.0.5.51000 > 93.184.216.34.443: Flags [S], seq 1000,
  win 64240, options [mss 1460,sackOK,TS val 111 ecr 0,nop,wscale 7], length 0
```

### 23.3 Sequence and acknowledgment arithmetic

- SYN and FIN each **consume one sequence number** even though they carry no
  data. So after `SYN seq=1000`, the first data byte is `seq=1001`.
- A pure ACK carries **no** data and does **not** advance the sender's sequence
  number (which is why a lost pure ACK is simply superseded by the next one — ACKs
  are not retransmitted).
- `tcpdump -S` shows absolute sequence numbers; by default it shows **relative**
  ones (starting at 0/1) which is far easier to read.

### 23.4 War story

A load balancer health check "passed" but real traffic failed. Capture showed the
LB's SYN had `wscale 7` and `sackOK`, the backend's SYN-ACK had **no options at
all** — an old middlebox was stripping TCP options on that path. Result: window
capped at 64 KB, no SACK; on the ~30 ms link, throughput ceilinged at
`65535 / 0.03 ≈ 17 Mbit/s` and any loss caused a full timeout. Routing the flow
around the option-stripping box restored gigabit. **Lesson:** always compare the
options on *both* SYNs; a middlebox that mangles them silently kneecaps TCP.

**Try it yourself**
```bash
sudo tcpdump -i any -n -c10 -S 'tcp[tcpflags] & (tcp-syn) != 0'    # all SYNs, abs seq
sudo tcpdump -i any -n 'tcp port 443' -vv | grep -m4 options
# Verify negotiated features for a live socket:
ss -tie dst 93.184.216.34    # shows wscale, sack, ts, cubic, rtt, cwnd, mss...
```

---

## 24. The three-way handshake and connection setup

### 24.1 The exchange

```
Client (active open)                         Server (passive open, LISTENing)
  |                                            |
  |  1. SYN   seq=C_isn                         |
  |            win, MSS, wscale, SACK-OK, TS  ->|   creates a mini-state entry
  |                                            |   (SYN queue), replies:
  |                                            |
  |<- 2. SYN,ACK  seq=S_isn  ack=C_isn+1        |
  |               win, MSS, wscale, SACK-OK, TS |
  |                                            |
  |  3. ACK   seq=C_isn+1  ack=S_isn+1        ->|   moves conn: SYN queue ->
  |                                            |   accept queue; app accept()s
  |                                            |
  |========== ESTABLISHED both sides ==========|
```

Why three (not two)? Each side must (a) announce its ISN and (b) get it
acknowledged, so both directions of the byte stream are synchronized. Two packets
could sync only one direction; four would be redundant because the middle
SYN+ACK combines the server's SYN with its ACK of the client's SYN.

### 24.2 The two queues on the server

- **SYN queue** (a.k.a. half-open / `tcp_max_syn_backlog`): connections that got a
  SYN and sent a SYN-ACK, awaiting the final ACK.
- **Accept queue** (a.k.a. completed / `backlog`, capped by `somaxconn`):
  handshake done, waiting for the app to `accept()`.

Overflow behavior (Linux):
- SYN queue full → new SYNs **dropped** (unless SYN cookies kick in).
- Accept queue full → depending on `net.ipv4.tcp_abort_on_overflow`, the final
  ACK is **ignored** (client later retransmits it; connection completes if the
  app catches up) or the connection is **RST**. Silent 1-second connect stalls
  under load are almost always accept-queue overflow.

### 24.3 What each side allocates

A **TCB** (Transmission Control Block) per connection: the 4-tuple, both ISNs,
current send/receive sequence numbers, windows, RTT estimators, congestion state,
retransmit queue, timers, socket buffers. Tens of KB to hundreds of KB
(dominated by socket buffers). A million connections is real but you must size
`tcp_mem`, file descriptors, and RAM for it.

### 24.4 The Initial Sequence Number must be unpredictable

Original TCP used a simple clock-based ISN → attackers could **guess** it and
inject data or spoof a whole connection (Morris, Mitnick). Modern stacks
(RFC 6528) compute `ISN = clock + hash(src IP, src port, dst IP, dst port,
secret_key)` — monotonic enough to avoid confusing old duplicate segments, random
enough to resist off-path injection. (Off-path RST/data injection still matters —
Ch. 52.)

### 24.5 SYN cookies — surviving a SYN flood

Under a **SYN flood** (Ch. 52), the SYN queue fills with half-open connections
from spoofed sources and legit clients can't connect. **SYN cookies**
(`net.ipv4.tcp_syncookies=1`, default) let the server keep **no state** for a
half-open connection: it encodes the essential SYN info (MSS bucket, timestamp,
tuple hash) *into* the SYN-ACK's ISN. If a real client returns the final ACK, the
server reconstructs the connection from `ack-1`. Cost: a few TCP options can't be
fully preserved (mitigated with the timestamp option). It only engages when the
queue overflows.

### 24.6 TCP Fast Open (TFO)

Saves the handshake RTT for repeat connections: the first visit, the server sends
a **TFO cookie**; on later visits the client puts `SYN + cookie + request data`
in one packet and the server can process the request before the handshake
finishes. Blocked or broken by some middleboxes/CDNs; used opportunistically.

### 24.7 Connection refused vs timeout vs unreachable — decode it

| Client sees | Wire event | Usual cause |
|---|---|---|
| `Connection refused` (`ECONNREFUSED`) | **RST** in response to SYN | nothing listening on that port; host is up |
| `Connection timed out` (`ETIMEDOUT`) | SYNs sent, **nothing** back (SYN retries then give up, ~1–2 min) | firewall `DROP`, host down, black hole, wrong route, PMTU black hole on the SYN-ACK |
| `No route to host` (`EHOSTUNREACH`) | local routing failure or `ICMP host unreachable` | no route, ARP fail for gateway, link down |
| `Network is unreachable` (`ENETUNREACH`) | no matching route at all | missing default route |
| Connect **succeeds** then immediate RST | SYN-ACK then RST | app crashed on accept, backlog overflow with `abort_on_overflow`, or an IPS resetting |

### 24.8 War story

A service saw ~0.1% of client connections take **exactly 1.0 s or 3.0 s** to
establish, matching Linux's SYN/ACK retransmit backoff. `ss -s` showed
`ListenOverflows` and `ListenDrops` climbing (`nstat -az | grep -i listen`). The
app's `listen()` backlog was 128 while `somaxconn` was 4096 and traffic was
bursty. Raising the app's backlog to 4096 and adding `SO_REUSEPORT` workers
eliminated the stalls. **Lesson:** 1s/3s/7s connect latencies = SYN or final-ACK
loss = check `nstat` for `TcpExtListenOverflows`/`TcpExtListenDrops` and the
app's actual backlog argument.

**Try it yourself**
```bash
nstat -az | grep -Ei 'listen|syncookie|synretrans'
ss -ltn                                  # Recv-Q = current accept-queue depth, Send-Q = its max
sysctl net.ipv4.tcp_max_syn_backlog net.core.somaxconn net.ipv4.tcp_syncookies
sysctl net.ipv4.tcp_abort_on_overflow
# Watch a handshake:
sudo tcpdump -i any -n '(tcp[tcpflags] & (tcp-syn|tcp-rst)) != 0' -c20
```

---

## 25. Connection teardown, half-close, and TIME_WAIT

### 25.1 The graceful four-way close

TCP closes **each direction independently** (it's two simplex streams):

```
Client                                     Server
  |  FIN, ACK   seq=u                     ->|   "I'm done sending"
  |<- ACK       ack=u+1                     |   (server may still send data)
  |                    ... server sends remaining response ...             |
  |<- FIN, ACK  seq=v                       |   "now I'm done too"
  |  ACK        ack=v+1                    ->|
  |                                         |
  | TIME_WAIT (2*MSL)                        | CLOSED
  | then CLOSED                              |
```

Often it's only **three** packets because the server piggybacks its FIN on the
ACK of the client's FIN (`FIN, ACK` in one segment).

### 25.2 Half-close

`shutdown(fd, SHUT_WR)` sends a FIN but keeps *reading*. Used by e.g. old HTTP
clients: "I've finished sending my request (FIN), now stream me the response."
The peer sees EOF on read but can still write. `close()` shuts down both
directions.

### 25.3 TIME_WAIT — why it exists and why it's usually fine

The side that sends the **last ACK** (usually whoever initiated the close) enters
**TIME_WAIT** for **2×MSL** (Maximum Segment Lifetime; Linux hardcodes MSL=60 s,
so TIME_WAIT ≈ 60 s; the sysctl `net.ipv4.tcp_fin_timeout` controls `FIN_WAIT_2`,
not this). Purposes:
1. **Absorb stragglers:** a delayed retransmission of the peer's FIN can still
   arrive; TIME_WAIT is there to ACK it instead of replying RST.
2. **Prevent 4-tuple reuse confusion:** stops a *new* connection with the same
   4-tuple from receiving *old* delayed segments.

TIME_WAIT sockets hold almost no memory and **don't** hold a file descriptor. On
a **client** making many outbound connections they can exhaust the ephemeral port
range to one destination (~28k / 60s ≈ 460 conn/s ceiling). On a **server** they
are generally harmless (the server side rarely initiates close, and when it does,
its "port" is fixed at 443 while the tuple varies by client).

### 25.4 Managing TIME_WAIT correctly

- **Best fix:** don't churn connections — use keep-alive / connection pools so you
  open few and reuse them.
- **`net.ipv4.tcp_tw_reuse = 1`** (client-side, outbound): lets a new *outbound*
  connection safely reuse a TIME_WAIT tuple when the TCP timestamp proves the old
  segments are older. Safe and recommended for busy clients.
- **`tcp_tw_recycle`**: **removed** in Linux 4.12. It broke NAT clients horribly
  (rejected connections from hosts sharing a NAT whose timestamps went
  "backwards"). Never enable it on anything where it still exists.
- Widen `ip_local_port_range`, add source IPs, or spread across more destination
  IPs/ports.
- `SO_LINGER {1,0}` makes `close()` send an **RST** instead of FIN → no
  TIME_WAIT, but you **lose** any unacked data and the peer sees a hard reset.
  Only for controlled cases (e.g. deliberately dumping idle connections).

### 25.5 CLOSE_WAIT — the bug indicator

If you see many sockets stuck in **CLOSE_WAIT**, that's *your* application's bug:
the peer sent FIN, the kernel ACKed it, but your code never called `close()` on
that fd. They pile up, leak file descriptors, and eventually `accept()` fails
with `EMFILE`. Fix the code path that forgets to close (error branches,
exceptions, missing `defer conn.Close()`).

### 25.6 RST — the abortive close

A RST tears the connection down immediately, discards buffered data, and the peer
gets `ECONNRESET` ("Connection reset by peer") on its next read/write. Causes:
- App set `SO_LINGER 0` and closed.
- Data arrived for a socket the app already closed.
- Segment arrived that doesn't match any connection (e.g. after one side
  rebooted, or a NAT mapping expired mid-connection).
- A firewall/IPS **forged** a RST to block the flow (Great Firewall, some
  corporate DLP, old carrier "TCP optimizers").

### 25.7 War story

An API gateway leaked file descriptors and fell over every ~6 hours with
`too many open files`. `ss -tan state close-wait | wc -l` grew monotonically.
The upstream HTTP client returned a response object but, on a specific 204
no-content path, never closed the body reader, so the underlying socket stayed
in CLOSE_WAIT. One-line fix (close the body in a `finally`/`defer`).
**Lesson:** rising CLOSE_WAIT = local missing `close()`; rising TIME_WAIT = you're
opening too many connections.

**Try it yourself**
```bash
ss -tan state time-wait  | wc -l
ss -tan state close-wait | wc -l
ss -tan state fin-wait-1 state fin-wait-2
sysctl net.ipv4.tcp_tw_reuse net.ipv4.tcp_fin_timeout
sudo tcpdump -i any -n 'tcp[tcpflags] & (tcp-fin|tcp-rst) != 0' -c20
```

---

## 26. The TCP finite state machine

### 26.1 The states

| State | Meaning |
|---|---|
| `CLOSED` | No connection. |
| `LISTEN` | Server waiting for SYNs. |
| `SYN_SENT` | Client sent SYN, waiting for SYN-ACK. |
| `SYN_RECV` | Server got SYN, sent SYN-ACK, waiting for final ACK (this is the SYN-queue state). |
| `ESTABLISHED` | Data transfer. |
| `FIN_WAIT_1` | We sent FIN, waiting for its ACK (or peer's FIN). |
| `FIN_WAIT_2` | Our FIN ACKed; waiting for peer's FIN. (`tcp_fin_timeout` bounds this ~60 s to avoid leaks against a peer that never closes.) |
| `CLOSING` | Both sent FIN ~simultaneously; waiting for ACK of ours. |
| `TIME_WAIT` | We sent the last ACK; wait 2×MSL. |
| `CLOSE_WAIT` | Peer sent FIN, we ACKed it; **waiting for our app to call `close()`**. |
| `LAST_ACK` | We (after CLOSE_WAIT) sent our FIN; waiting for its ACK. |

### 26.2 The diagram (active vs passive paths)

```
                              CLOSED
                    active open |        \ passive open
                   send SYN     |         \  (listen())
                                v          v
        rcv SYN,ACK          SYN_SENT     LISTEN
        send ACK   +-----------|              | rcv SYN
                   |           | (simultaneous| send SYN,ACK
                   |            open path)    v
                   |                       SYN_RECV
                   |          rcv ACK  /       |
                   v                  /        |
             ESTABLISHED <-----------+         |
        ______/   |    \______                 
close()/send FIN  | rcv FIN / send ACK          
       v          |         v                   
   FIN_WAIT_1     |     CLOSE_WAIT               
   |    | rcv ACK |         | close()/send FIN   
   |    v         |         v                    
   | FIN_WAIT_2   |     LAST_ACK                 
   |    | rcv FIN |         | rcv ACK            
rcv FIN | send ACK|         v                    
send ACK|         |      CLOSED                  
   v    v                                        
  CLOSING -> (rcv ACK) -> TIME_WAIT -> (2MSL) -> CLOSED
```

### 26.3 How to read states operationally

- Many `SYN_RECV` on a server = SYN flood or a client population that sends SYN
  but not the final ACK (broken NAT, PMTU black hole on the SYN-ACK).
- Many `SYN_SENT` on a client = destination unreachable / firewalled / overloaded.
- `FIN_WAIT_2` piling up = the **peer** isn't closing its side (peer app bug or
  peer overload). `tcp_fin_timeout` will eventually reap them.
- `CLOSE_WAIT` piling up = **local** app bug (didn't `close()`), §25.5.
- `LAST_ACK` stuck = the final ACK from the peer isn't arriving (loss, peer gone).
- `TIME_WAIT` large but stable = normal for a busy connection initiator; only act
  if you're near port exhaustion.

### 26.4 The kernel timers attached to a connection

| Timer | Fires when | Effect |
|---|---|---|
| **Retransmission (RTO)** | an unACKed segment's timer expires | resend oldest unacked segment, back off RTO, cut congestion window (Ch. 28, 31) |
| **Persist** | receiver advertised a **zero window** | send tiny "window probes" so a lost window-update doesn't deadlock the connection (Ch. 30) |
| **Keepalive** | connection idle for `tcp_keepalive_time` (default 7200 s) | send a probe; after `tcp_keepalive_probes` failures, kill the connection (Ch. 34) |
| **TIME_WAIT (2MSL)** | last ACK sent | reap the TCB |
| **FIN_WAIT_2** | our FIN ACKed but no peer FIN | reap after `tcp_fin_timeout` |
| **Delayed ACK** | data received, nothing to piggyback on | send a bare ACK after up to ~40–200 ms (Ch. 33) |

### 26.5 War story

A batch job's worker connections sat in `ESTABLISHED` forever after the remote
database VM was hard-killed (no FIN, no RST). The app blocked on `read()`
indefinitely; the job hung for hours. Default keepalive is **2 hours** idle +
9 probes × 75 s ≈ **2h11m** before the kernel gives up — far too slow. Setting
`TCP_KEEPIDLE=30, TCP_KEEPINTVL=10, TCP_KEEPCNT=3` on the socket (or an
application-level read deadline) made it fail in ~1 minute and retry.
**Lesson:** the OS will not promptly tell you a peer vanished; set per-socket
keepalive or app-level timeouts for anything long-lived.

**Try it yourself**
```bash
ss -tan | awk '{print $1}' | sort | uniq -c | sort -rn   # histogram of TCP states
watch -n1 'ss -tan state syn-recv | wc -l'
sysctl net.ipv4.tcp_keepalive_time net.ipv4.tcp_keepalive_intvl net.ipv4.tcp_keepalive_probes
ss -tie                                   # per-socket: state, timers, retrans count
```

---
# Part VIII — TCP reliability

## 27. Sequence numbers, cumulative ACKs, and retransmission

### 27.1 The model

Every byte in each direction has a sequence number. The receiver returns an
**acknowledgment number** = "the next byte I expect" = one past the highest
**contiguous** byte received. This is a **cumulative** ACK: `ack=5001` means
"I have everything through byte 5000," regardless of what arrived in between.

```
Sender sends 4 segments, 1000 B each, starting seq=1:
  seg A: seq 1     .. 1000
  seg B: seq 1001  .. 2000
  seg C: seq 2001  .. 3000   <-- LOST
  seg D: seq 3001  .. 4000

Receiver:
  gets A -> ACK 1001
  gets B -> ACK 2001
  gets D -> still ACK 2001   (can't ack past the hole; this is a "duplicate ACK")
```

### 27.2 Two ways a loss is detected

1. **Retransmission timeout (RTO):** the sender's timer for the oldest unacked
   segment expires with no progress → resend it. Slow (RTO ≥ 200 ms, often
   1 s on the first go) and drastic (congestion window collapses to 1 — Ch. 31).
2. **Duplicate ACKs → fast retransmit:** three ACKs carrying the *same* ack
   number tell the sender "a later segment arrived but there's a gap at `ack`."
   The sender resends that one segment **immediately**, without waiting for the
   RTO (Ch. 29). Much faster.

### 27.3 What the receiver does with out-of-order data

It **buffers** it (up to the receive window) and keeps ACKing the last in-order
byte. When the missing segment finally arrives, the ACK **jumps forward** past
all the buffered data in one step. Without SACK, the sender can't tell *which*
later segments arrived, so a naive sender might resend B, C, and D
("go-back-N"–ish). SACK (Ch. 29) fixes exactly this.

### 27.4 Spurious retransmissions

If an ACK is merely **delayed** (not lost) and the RTO is too tight, the sender
retransmits data the receiver already has. The receiver drops the duplicate and
ACKs normally, but the sender needlessly cut its window. Mitigations:
- **Karn's algorithm** (Ch. 28): don't measure RTT from a retransmitted segment.
- **TCP timestamps**: let the receiver echo *which* send the ACK corresponds to →
  the sender can tell a retransmit was spurious.
- **F-RTO** (`net.ipv4.tcp_frto`, on by default): a heuristic that detects a
  spurious timeout from the pattern of the next ACKs and **un-does** the window
  reduction.
- **DSACK** (duplicate SACK): the receiver explicitly reports "you sent me this
  block twice," letting the sender learn its RTO/reordering estimate was wrong.

### 27.5 Reordering vs loss

The internet reorders packets (ECMP path changes, link bonding, parallel switch
fabrics). Three dup-ACKs from *reordering* (not loss) trigger an unnecessary fast
retransmit. Linux adapts: `tcp_reordering` starts at 3 and **RACK-TLP**
(`tcp_recovery`, default on modern kernels) replaces the dup-ACK counting with a
**time-based** notion — "a segment is lost if a later-sent segment was ACKed and
enough time has passed" — which is far more robust to reordering.

### 27.6 War story

A storage replication link across two data centers showed 6% "retransmissions" in
`nstat` but zero actual packet loss on the carrier's reports. `tcpdump` +
Wireshark's *Expert Info* showed **spurious retransmissions** and lots of
**out-of-order** segments: the pair of 10G links was hashing the flow's packets
across *both* members inconsistently, reordering them. Fixing the LAG hash to be
per-flow (5-tuple) dropped "retransmissions" to ~0.02% and doubled throughput.
**Lesson:** high "retransmit" counters with no carrier loss = reordering; look at
out-of-order and spurious-retransmit stats, and check LACP/ECMP hashing.

**Try it yourself**
```bash
nstat -az | grep -Ei 'retrans|reorder|dsack|spurious|lostretransmit'
# e.g. TcpRetransSegs, TcpExtTCPSpuriousRTOs, TcpExtTCPDSACKRecv, TcpExtTCPReordering
ss -tie dst <peer>          # per-socket: retrans:X/Y  (current/total), reordering N
sysctl net.ipv4.tcp_recovery net.ipv4.tcp_frto net.ipv4.tcp_reordering
```

---

## 28. Retransmission timeout (RTO), RTT estimation, and Karn's algorithm

### 28.1 Why RTT estimation is hard

The RTO must be **long enough** that a normal ACK isn't mistaken for loss, but
**short enough** to recover quickly. RTT varies constantly (queueing, route
changes, Wi-Fi). TCP measures it and adapts.

### 28.2 The Jacobson/Karels algorithm (RFC 6298)

Maintain a smoothed RTT (`SRTT`) and a smoothed mean deviation (`RTTVAR`):

```
On first RTT sample R:
    SRTT   = R
    RTTVAR = R / 2

On each subsequent sample R:
    RTTVAR = (1 - 1/4) * RTTVAR + 1/4 * |SRTT - R|
    SRTT   = (1 - 1/8) * SRTT   + 1/8 * R

    RTO    = SRTT + max(G, 4 * RTTVAR)      # G = clock granularity
    RTO    = clamp(RTO, 1s (RFC min) .. 60s)
```

- Weighting new samples at 1/8 makes SRTT a low-pass filter — stable but slow to
  react.
- Including `4 * RTTVAR` means a **jittery** path gets a generously padded RTO
  (fewer spurious retransmits) while a **stable** path gets a tight RTO (fast
  recovery). This variance term is the key insight that fixed 1980s TCP.
- Linux uses a practical **minimum RTO of 200 ms** (`TCP_RTO_MIN`), not the RFC's
  1 s, and can be tuned per-route (`ip route ... rto_min 20ms` inside a
  data center).

### 28.3 Karn's algorithm — the retransmission ambiguity

If you send segment X, it times out, you resend X, then an ACK arrives — **was
that ACK for the original or the retransmission?** You can't tell, so you can't
compute a valid RTT sample. **Karn's algorithm:**
1. **Do not** take an RTT sample from any segment that was retransmitted.
2. **Do** apply **exponential backoff** to the RTO on each retransmit
   (`RTO = RTO * 2`, capped), and keep the backed-off RTO until a segment is
   ACKed *without* retransmission (then resume normal calculation).

Backoff is what makes TCP retreat politely from a badly congested or broken path
instead of hammering it: retries at roughly 1s, 2s, 4s, 8s, 16s… (Linux caps
around `tcp_retries2` ≈ 15 attempts ≈ 13–30 minutes before giving up on an
established connection; `tcp_syn_retries` ≈ 6 for the initial SYN ≈ ~127 s).

### 28.4 Timestamps make RTT measurement continuous

With the **TCP Timestamps** option, every segment carries `TSval`; the ACK echoes
it in `TSecr`. The sender computes `RTT = now - TSecr` on **every** ACK, not just
once per window, and Karn's ambiguity disappears (the echoed timestamp
disambiguates original vs retransmit). This gives a much better `SRTT`, and also
powers **PAWS** (rejecting old segments whose sequence numbers wrapped on a fast
link).

### 28.5 Tail loss and the RTO's worst case

If the **last** segments of a transfer are lost, there are no later segments to
generate dup-ACKs, so fast retransmit can't trigger — you eat a full RTO (often
200 ms+). This is deadly for small request/response RPCs. Fixes:
- **TLP (Tail Loss Probe):** ~2×SRTT after the last segment, send one probe
  (retransmit the last segment or a new one) to elicit ACKs and trigger fast
  recovery instead of an RTO. Part of **RACK-TLP**, on by default in Linux.
- **Early Retransmit** and lowering the dup-ACK threshold when the window is
  small.

### 28.6 War story

An RPC service's p99.9 latency had a hard cliff at **+200 ms**. Traces showed the
final segment of small responses occasionally lost; with no trailing data, the
client waited a full `TCP_RTO_MIN` (200 ms) before the server's RTO fired.
Enabling `tcp_early_retrans`/RACK-TLP on the servers (newer kernel) and, on the
low-RTT data-center path, setting `ip route ... rto_min 20ms` collapsed the p99.9
spike. **Lesson:** a latency histogram with a spike exactly at 200 ms / 1 s = an
RTO firing; look at tail-loss recovery (TLP) and `rto_min`.

**Try it yourself**
```bash
ss -tie dst <peer>     # shows rtt:<srtt>/<rttvar>  rto:<ms>  ato:<delayed-ack timeout>
nstat -az | grep -Ei 'TCPLossProbe|TCPTimeouts|TCPSpuriousRTOs|TCPRenoRecovery'
sysctl net.ipv4.tcp_timestamps net.ipv4.tcp_early_retrans net.ipv4.tcp_retries2
ip route show cache
```

---

## 29. Fast retransmit, fast recovery, and SACK

### 29.1 Fast retransmit

Rule: **on the 3rd duplicate ACK** (same ack number, no new data, zero window
change), retransmit the segment starting at that ack number **immediately** —
don't wait for the RTO.

Why 3 and not 1? To tolerate mild reordering: 1–2 dup-ACKs commonly just mean
packets arrived slightly out of order. 3 is a heuristic threshold (RACK replaces
it with time, §27.5).

### 29.2 Fast recovery (the "half, not zero" idea)

A timeout means "the network might be dead" → drop congestion window to 1
(restart slow start). But three dup-ACKs mean "packets are *still flowing*" (the
receiver got later segments), so the reaction is gentler:

**TCP Reno fast recovery (RFC 5681):**
1. On 3rd dup-ACK: `ssthresh = max(FlightSize / 2, 2*MSS)`; retransmit the lost
   segment; `cwnd = ssthresh + 3*MSS` (the 3 accounts for the 3 segments that
   *left* the network — that's what generated the dup-ACKs).
2. For each **additional** dup-ACK: `cwnd += MSS` ("inflate" — more segments have
   left the network, so you may inject one).
3. When a **new** ACK finally arrives (the retransmit got through):
   `cwnd = ssthresh` ("deflate") and resume **congestion avoidance** (linear
   growth), *not* slow start.

Net effect: one loss ≈ halve the rate, then keep going. This is the "sawtooth"
of classic TCP throughput.

### 29.3 Why plain Reno struggles with multiple losses per window

Reno's dup-ACK counting only tells the sender about **one** hole at a time.
Lose two segments in the same window and Reno often needs a second recovery
round (or an RTO). **NewReno** (RFC 6582) patches this: it remembers the highest
sequence sent when recovery began ("recovery point") and stays in fast recovery,
retransmitting one segment per RTT, until *all* data up to that point is ACKed —
no RTO needed for multiple losses, but still only one repair per RTT.

### 29.4 SACK — Selective Acknowledgment (RFC 2018), the real fix

Negotiated by **SACK-permitted** on both SYNs. Then a receiver with holes sends
**SACK blocks** in the ACK's options: up to 3–4 pairs of
`(left edge, right edge)` describing contiguous ranges it *has* received above the
cumulative ack.

```
Cumulative ACK = 2001
SACK: 3001-4001, 5001-6001
Meaning: "I have through 2000, AND 3001-4000, AND 5001-6000.
          I'm missing 2001-3000 and 4001-5000."
```

The sender now retransmits **exactly** the two missing ranges, in one RTT,
without guessing. Combined with a **scoreboard** (per-segment SACKed/lost state)
and **RACK** for loss detection, this is how modern Linux TCP recovers from bursty
loss efficiently.

- **DSACK** (RFC 2883): the receiver also reports ranges it received
  **twice** ("you retransmitted unnecessarily"), so the sender can raise its
  reordering estimate / undo a spurious `cwnd` cut.
- **FACK / RACK**: heuristics/algorithms for deciding, from SACK info, how much
  data is actually "lost" vs merely "in flight," to set `cwnd` and what to resend.

### 29.5 PRR — Proportional Rate Reduction (RFC 6937)

Replaces Reno's abrupt "halve then inflate/deflate" during recovery with a
smooth, ACK-clocked reduction that converges `cwnd` to `ssthresh` by the end of
recovery without the bursty inflate/deflate. Default in Linux. It's why modern
TCP recovers with fewer secondary losses.

### 29.6 War story

A backup over a satellite link (RTT ~600 ms, occasional 1–2% bursty loss) crawled
at 3 Mbit/s on a 50 Mbit/s pipe. Capture: **SACK was disabled** because a
"security" middlebox stripped the SACK-permitted option from the SYN. Every loss
event forced NewReno to repair one segment per 600 ms RTT — catastrophic on a
long fat pipe. Whitelisting the option through the middlebox (and enabling
`tcp_sack`, `tcp_dsack`) took the backup to ~40 Mbit/s. **Lesson:** on
high-latency or lossy paths, **no SACK = no throughput**; verify both SYNs carry
`sackOK`.

**Try it yourself**
```bash
sysctl net.ipv4.tcp_sack net.ipv4.tcp_dsack net.ipv4.tcp_fack
nstat -az | grep -Ei 'sack|dsack|renorecovery|sackrecovery|prr|tcplostretransmit'
# In a capture, Wireshark: filter  tcp.analysis.flags  - flags duplicate ACKs,
# fast retransmits, out-of-order, spurious retransmits, zero windows.
sudo tcpdump -i any -n 'tcp[tcpflags] & tcp-ack != 0' -v | grep -m5 'sack'
```

---
# Part IX — TCP flow and congestion control

Two independent limits govern how much unacknowledged data a TCP sender may have
"in flight":

```
bytes_in_flight  <=  min( rwnd , cwnd )

rwnd  = receiver's advertised window  -> FLOW control   (don't overrun the peer's buffer)
cwnd  = sender's congestion window    -> CONGESTION ctrl (don't overrun the network)
```

Flow control is a hard contract from the receiver. Congestion control is the
sender's own estimate of what the path can take.

## 30. Flow control: the sliding window, window scaling, and SWS

### 30.1 The sliding window

The receiver advertises `rwnd` in every ACK = free space in its socket receive
buffer. The sender may send up to `snd_una + min(rwnd, cwnd)` — i.e. bytes from
"oldest unacked" forward. As ACKs come in, the window **slides** right.

```
   sent+acked | sent, not yet acked |  may send now  | can't send yet
 --------------+---------------------+----------------+---------------->
             snd_una              snd_nxt        snd_una + wnd
              <-------- window (min(rwnd,cwnd)) -------->
```

If the app on the receiver stops calling `read()`, its buffer fills, `rwnd`
shrinks toward 0, and the sender throttles — automatically. This is why a slow
consumer naturally slows a fast producer over TCP ("backpressure for free").

### 30.2 Window scaling — the 64 KB wall

The header's Window field is **16 bits → max 65 535 bytes**. Max throughput is
`window / RTT`:

| RTT | Without scaling (64 KB) | Needed window for 1 Gbit/s |
|---|---|---|
| 1 ms (same rack) | 512 Mbit/s | 125 KB |
| 30 ms (regional) | 17 Mbit/s | 3.75 MB |
| 100 ms (intercontinental) | 5 Mbit/s | 12.5 MB |
| 300 ms (satellite) | 1.7 Mbit/s | 37.5 MB |

**Window Scale option** (SYN only): a shift count `s` (0–14); the real window is
`advertised << s`. With `s=7`, a 65 535 advertisement means ~8 MB. **Both SYNs
must carry it** or scaling is off for the whole connection. Linux enables it by
default (`net.ipv4.tcp_window_scaling=1`); the actual max is bounded by
`net.ipv4.tcp_rmem` / `tcp_wmem` and autotuning (Ch. 35).

Middleboxes that strip or rewrite the WScale option cause a nasty asymmetry: one
side scales, the other doesn't, and windows are misinterpreted by orders of
magnitude → stalls or resets.

### 30.3 Zero window and the persist timer

If the receiver's buffer is completely full it advertises **window = 0**. The
sender stops. When the app drains the buffer, the receiver sends a **window
update** ACK. But **ACKs are not retransmitted** — if that update is lost, both
sides wait forever. The **persist timer** prevents deadlock: the sender
periodically sends a 1-byte **zero-window probe**; the receiver responds with the
current (possibly still 0) window. Eventually the real window is learned.

Seeing frequent `TCP ZeroWindow` / `TCP Window Full` in Wireshark = **the
receiving application is too slow** (not reading fast enough). The network is
fine; fix the consumer (more worker threads, faster processing, bigger buffer).

### 30.4 Silly Window Syndrome (SWS)

Pathology where the window is repeatedly advanced by tiny amounts, so the
connection exchanges many small segments (huge header overhead). Two guards:
- **Receiver side (David Clark's fix):** don't advertise a larger window until it
  can grow by at least one MSS (or half the buffer).
- **Sender side (Nagle's algorithm, Ch. 33):** don't send a small segment while
  a previous small segment is still unacked; coalesce.

### 30.5 Buffer autotuning

Modern Linux dynamically sizes each connection's receive buffer between
`tcp_rmem` min/default/max based on the measured bandwidth-delay product
(`tcp_moderate_rcvbuf=1`). You usually shouldn't set `SO_RCVBUF`/`SO_SNDBUF`
manually — doing so **disables autotuning** and often *hurts*. Only raise the
`tcp_rmem`/`tcp_wmem` **max** (Ch. 35, 37) for long-fat-network workloads.

### 30.6 War story

A log shipper's throughput to a central collector was stuck at ~90 Mbit/s despite
a 10G network, RTT ~25 ms. `ss -tim` showed `rcv_space`/`rwnd` pinned near
3 MB — the collector app had called `setsockopt(SO_RCVBUF, 2MB)` "to bound memory,"
which **capped and froze** the window and disabled autotuning. `3MB / 0.025s ≈
960 Mbit/s` ceiling, minus overhead ≈ observed. Removing the manual `SO_RCVBUF`
let autotuning grow it to ~25 MB and throughput hit multi-Gbit/s.
**Lesson:** manual socket buffer sizes usually do more harm than good; tune the
sysctl ceilings and let autotuning work.

**Try it yourself**
```bash
ss -tim dst <peer>       # rcv_space, snd_wnd, rcv_wnd, mss, cwnd, bytes_acked...
sysctl net.ipv4.tcp_window_scaling net.ipv4.tcp_moderate_rcvbuf
sysctl net.ipv4.tcp_rmem net.ipv4.tcp_wmem
nstat -az | grep -Ei 'TCPWantZeroWindow|TCPZeroWindow|TCPBacklogDrop|PruneCalled'
# Wireshark filter: tcp.analysis.zero_window || tcp.analysis.window_full
```

---

## 31. Congestion control: slow start and congestion avoidance

### 31.1 The problem: congestion collapse

October 1986: the NSFNET backbone throughput dropped from 32 kbit/s to
**40 bit/s** — a 1000× collapse. Cause: senders retransmitting into an already
congested network, adding load exactly when they should back off, in a vicious
cycle. Van Jacobson's 1988 congestion-control algorithms (slow start, congestion
avoidance, fast retransmit) saved the internet and are still the skeleton of
every TCP.

**Core principle — conservation of packets / ACK clocking:** in equilibrium, a
sender injects a new packet only when an ACK signals one has left the network.
The ACK stream *is* the clock.

### 31.2 The variables

- **`cwnd`** — congestion window (bytes the sender believes the path can hold).
- **`ssthresh`** — slow-start threshold: the `cwnd` boundary between exponential
  and linear growth.
- **IW** — initial window. Historically 1–4 MSS; **RFC 6928 raised it to 10 MSS**
  (~14.6 KB), so a small HTTP response often fits in the very first burst.

### 31.3 Slow start (exponential)

```
cwnd starts at IW (10 MSS on Linux)
for each ACK received:  cwnd += MSS
```
Per RTT, `cwnd` roughly **doubles** (every segment in the window gets ACKed, each
adding an MSS). Growth: 10 → 20 → 40 → 80 … MSS. "Slow" is historical — it starts
small, but ramps *fast*. It ends when:
- `cwnd >= ssthresh` → switch to congestion avoidance, or
- a loss occurs → set `ssthresh = cwnd/2`, react (Ch. 29 / below).

### 31.4 Congestion avoidance (linear / AIMD)

Once `cwnd >= ssthresh`, grow gently:
```
per RTT:  cwnd += 1 MSS         (implemented per-ACK as cwnd += MSS*MSS/cwnd)
```
This is **Additive Increase**. On a loss: **Multiplicative Decrease**
(`cwnd *= 0.5` for Reno). Together, **AIMD** — the property that makes many TCP
flows converge to a roughly fair share of a bottleneck: additive increase nudges
everyone up equally; multiplicative decrease cuts the biggest flows the most.

```
cwnd
 |        /|        /|        /|
 |       / |       / |       / |      <- congestion avoidance (linear rise)
 |      /  |      /  |      /  |
 |     /   |     /   |     /   |
 |    /    v    /    v    /    v      <- loss: cut cwnd in half
 |   /  (slow   
 |  /   start)  
 |_/____________________________________ time
      "TCP sawtooth"
```

### 31.5 Reaction to loss, summarized

| Detection | ssthresh | cwnd after | Then |
|---|---|---|---|
| **RTO (timeout)** | `max(FlightSize/2, 2·MSS)` | **1 MSS** | slow start again |
| **3 dup-ACKs (fast retransmit)** | `max(FlightSize/2, 2·MSS)` | `ssthresh` (+3 for Reno inflate) | fast recovery → congestion avoidance |

A timeout is the expensive event: full restart. Everything in modern TCP (SACK,
RACK, TLP, PRR) is about **avoiding RTOs** and recovering via the fast path.

### 31.6 ECN — congestion signal without loss

With **ECN** (Explicit Congestion Notification) negotiated, a router experiencing
congestion **marks** the IP header (CE) instead of **dropping** the packet. The
receiver echoes this via the **ECE** TCP flag; the sender reduces `cwnd` as if
there'd been a loss and sets **CWR**. Benefit: congestion response **without**
losing/retransmitting data — lower latency, no goodput loss. Needs endpoints +
routers to cooperate; historically under-deployed, but **L4S** (Low Latency, Low
Loss, Scalable throughput; RFC 9330) and **DCTCP** in data centers use fine-
grained ECN aggressively. Linux: `net.ipv4.tcp_ecn` (0 off / 1 request /
2 accept-only).

### 31.7 Fairness, RTT bias, and why it's imperfect

- **RTT unfairness:** classic AIMD grows `cwnd` per RTT, so a **short-RTT** flow
  ramps faster and steals bandwidth from a **long-RTT** flow sharing the same
  bottleneck. CUBIC reduces (not eliminates) this by making growth a function of
  *time since last loss*, not RTT.
- **Number-of-flows unfairness:** an app opening 8 parallel connections gets ~8×
  the share of one. (This is literally why browsers open 6 connections per host
  and why "download accelerators" work — they're not more efficient, just
  greedier.)
- **Loss-based vs delay-based** flows don't share fairly (BBR vs CUBIC — Ch. 32).

### 31.8 War story

A video CDN moved a POP closer to users (RTT 80 ms → 12 ms). Average bitrate
*rose* even though bandwidth was unchanged — because with a 6.7× shorter RTT,
each flow's `cwnd` ramped 6.7× faster through slow start and recovered from loss
6.7× faster, so more of the pipe was actually used. **Lesson:** TCP throughput is
governed by RTT as much as by bandwidth; `Rate ≈ cwnd/RTT`, and cwnd growth
itself is per-RTT. Latency is a throughput feature.

**Try it yourself**
```bash
ss -tie dst <peer>      # cwnd:<segs>  ssthresh:<segs>  bytes_sent  bytes_retrans  send <bps>
sysctl net.ipv4.tcp_congestion_control       # active default (usually cubic)
sysctl net.ipv4.tcp_ecn
nstat -az | grep -Ei 'TCPHystartTrainDetect|TCPHystartDelayDetect|TCPSlowStartRetrans'
# Plot cwnd over time for a transfer:
ss -tie -o state established '( dst <peer> )'   # sample in a loop; or use `tcp_probe`/bpftrace
```

---

## 32. Reno, NewReno, CUBIC, and BBR

### 32.1 The family tree

| Algorithm | Signal it reacts to | Increase rule | Decrease on loss | Where it shines / hurts |
|---|---|---|---|---|
| **Tahoe** (1988) | loss (dup-ACK or RTO) | slow start / linear | → 1 (always slow start) | historical |
| **Reno** (1990) | loss | AIMD | ×0.5, fast recovery | one loss/window only |
| **NewReno** (1999) | loss | AIMD | ×0.5, stays in recovery for multiple holes | still 1 repair/RTT without SACK |
| **CUBIC** (2006, **Linux default**) | loss | **cubic function of time since last loss** | ×0.7 (β=0.7) | high BDP, RTT-fairer than Reno; still fills buffers → bufferbloat |
| **BBR** / **BBRv2/v3** (2016–, Google) | **measured bottleneck bandwidth & RTT** (model-based, mostly ignores loss) | probe BW/RTT in cycles | doesn't halve on loss | high throughput on lossy/long paths, low queueing; v1 could be unfair to CUBIC and starve shallow buffers |
| **Vegas** (1994) | **increasing RTT** (delay) | proactive | back off before loss | pure delay-based; loses to loss-based flows |
| **DCTCP** (2010) | **fraction of ECN-marked ACKs** | proportional | fine-grained | data-center only (needs ECN + shallow-buffer switches); tiny queues |
| **BBR** vs **CUBIC** coexistence | — | — | — | on a shared bottleneck, BBRv1 often takes more than its share; BBRv2/v3 add loss + ECN response to fix this |
| **Illinois, YeAH, H-TCP, Scalable, HighSpeed** | loss (+ delay hints) | aggressive high-BDP variants | varies | niche / legacy long-fat-network tuning |

### 32.2 CUBIC in one paragraph

After a loss at window `W_max`, CUBIC sets a cubic curve
`cwnd(t) = C·(t − K)³ + W_max`, where `K` is the time to grow back to `W_max`.
Just after loss the curve is **concave** — grows fast toward `W_max` (recover the
lost ground quickly). Near `W_max` it **flattens** (probe cautiously around the
last known ceiling). Past `W_max` it turns **convex** — accelerates to look for
new capacity. Because `t` is wall-clock time, two flows with different RTTs on the
same bottleneck grow similarly → better RTT fairness than Reno. Downside: still
**loss-based**, so it keeps pushing until the bottleneck buffer overflows —
filling deep buffers and adding latency (bufferbloat, Ch. 44).

### 32.3 BBR in one paragraph

BBR builds a **model** of the path: `BtlBw` (max delivery rate seen) and `RTprop`
(min RTT seen). It paces sending at `BtlBw` and keeps ~1 BDP in flight, so it
achieves high throughput **without** filling the buffer — latency stays near
`RTprop`. It periodically **probes**: briefly send faster to test for more
bandwidth (`ProbeBW`), and every ~10 s send slower to re-measure bare RTT
(`ProbeRTT`). It largely **ignores packet loss** as a congestion signal, which is
why it screams on lossy Wi-Fi / trans-oceanic links where CUBIC collapses — but
also why BBRv1 can be aggressive/unfair toward loss-based flows and misbehave in
very shallow buffers. **BBRv2/v3** add explicit responses to loss and ECN and are
friendlier.

### 32.4 Choosing one (Linux)

```bash
sysctl net.ipv4.tcp_available_congestion_control     # what's compiled/loaded
modprobe tcp_bbr                                      # load BBR
sysctl -w net.ipv4.tcp_congestion_control=bbr         # set default
sysctl -w net.core.default_qdisc=fq                   # BBR wants fq (or fq_codel) for pacing
```
- **Default (CUBIC):** fine for most; safe, fair with the internet.
- **BBR:** great for a CDN/edge/proxy serving users over lossy or long paths;
  test fairness if you share bottlenecks with your own CUBIC traffic.
- **DCTCP:** only inside a data center you fully control (ECN end-to-end, AQM
  switches).
- Per-connection override: `setsockopt(TCP_CONGESTION, "bbr")`, or per-route
  `ip route ... congctl bbr`.

### 32.5 Pacing

Bursting a whole `cwnd` back-to-back causes queue spikes and loss. **Pacing**
spreads segments over the RTT at roughly `cwnd/RTT`. BBR requires it; CUBIC
benefits. Linux does it via the **`fq`** qdisc or internal TSO-aware pacing
(`net.ipv4.tcp_pacing_ss_ratio` / `tcp_pacing_ca_ratio`).

### 32.6 War story

A backup service in Singapore uploading to `us-east-1` (RTT ~230 ms, ~0.3%
loss) maxed at ~25 Mbit/s with CUBIC on a 1 Gbit/s link — the Mathis formula
(Ch. 36) predicts `MSS/(RTT·√p) ≈ 1448/(0.23·0.055) ≈ 114 KB/s ≈ ~0.9 Mbit/s`
per flow, so even 25 needed dozens of parallel streams. Switching the senders to
**BBR** (+`fq` qdisc) took a *single* stream to ~700 Mbit/s because BBR treats
that 0.3% loss as noise, not a signal to halve. **Lesson:** on long, slightly
lossy paths, loss-based CC is throughput-limited by `√loss`; BBR or many parallel
CUBIC streams are the escape hatches.

**Try it yourself**
```bash
ss -tie | grep -oE '(cubic|bbr|reno|dctcp)[^ ]*'      # which CC each socket uses
cat /proc/sys/net/ipv4/tcp_congestion_control
# A/B test on one host:
iperf3 -c <server> -t30 -C cubic
iperf3 -c <server> -t30 -C bbr
```

---

## 33. Nagle, delayed ACK, and the interaction bug

### 33.1 Nagle's algorithm (sender-side)

Goal: stop tiny "tinygram" segments (1-byte payload + 40-byte headers) from
swamping the network (the 1980s "telnet over a congested link" problem).

Rule: **if there is unacknowledged data outstanding, buffer new small writes
until either (a) a full MSS worth accumulates, or (b) an ACK arrives.** If nothing
is unacked, send immediately.

Effect: back-to-back small `write()`s get coalesced into fewer, bigger segments —
great for bulk/throughput, **bad for latency-sensitive request/response**.

Disable per-socket with **`TCP_NODELAY`**. Almost every RPC framework, database
driver, game server, and low-latency system sets `TCP_NODELAY`. Redis, gRPC,
nginx upstreams, PostgreSQL, SSH interactive — all disable Nagle.

### 33.2 Delayed ACK (receiver-side)

Goal: reduce pure-ACK traffic and enable piggybacking. Rule: **don't ACK
immediately; wait up to ~40 ms (Linux; 200 ms on some OSes) hoping to (a) have
return data to piggyback the ACK on, or (b) receive a second segment so one ACK
covers both.** Standard says ACK at least every 2nd full segment.

### 33.3 The classic 40 ms (or 200 ms) stall

Nagle + delayed ACK can **deadlock each other** for a request/response pattern
where the sender writes the request in **two** chunks (e.g. header then body):

```
Client (Nagle ON)                       Server (delayed ACK ON)
  write(header)  -> sent (nothing unacked yet)
  write(body)    -> BUFFERED by Nagle (header is unacked!)
                                         got header; no data to reply with yet;
                                         waits for a 2nd segment or 40ms before ACK
   ... both sides wait ...
                            (40 ms later) server delayed-ACK fires -> ACK header
  ACK arrives -> Nagle releases body -> sent
                                         server gets body, processes, replies
```

Every such round trip eats a fixed **~40 ms** (or 200 ms on older stacks). Under
load this shows up as throughput that's suspiciously close to `N / 0.04` requests
per second, or p50 latency pinned near 40/200 ms.

**Fixes (any one):**
- **`TCP_NODELAY` on the client** (most common; disables Nagle). Do this for
  virtually all interactive/RPC sockets.
- **Write the whole request in one `write()`** / use writev / buffer in userspace
  then flush once. No small-then-small pattern → Nagle never engages.
- **`TCP_QUICKACK`** on the server (Linux, must be re-set; it's not sticky) to
  suppress delayed ACK.
- `TCP_CORK` when you *want* coalescing deliberately (e.g. sendfile headers +
  file): cork, write parts, uncork to flush as full segments — controlled Nagle.

### 33.4 `TCP_CORK` vs `TCP_NODELAY`

- `TCP_NODELAY` = "never wait, send now" (latency).
- `TCP_CORK` = "hold everything until I uncork or 200 ms pass, then send full
  segments only" (efficiency, e.g. HTTP response headers + body in minimal
  packets). They're opposites; nginx toggles both around `sendfile()`.

### 33.5 War story

A trading gateway's median order-ack latency was **40.1 ms** in staging,
sub-millisecond in production. Staging used a client library build **without**
`TCP_NODELAY`; it wrote a 26-byte header then a variable body in a second
`send()`. Nagle held the body until the server's delayed ACK fired 40 ms later.
Setting `TCP_NODELAY` (one line) dropped median to 0.3 ms. **Lesson:** a latency
floor at exactly 40 ms or 200 ms is Nagle × delayed-ACK until proven otherwise.

**Try it yourself**
```bash
strace -f -e trace=setsockopt yourapp 2>&1 | grep -E 'TCP_NODELAY|TCP_CORK|TCP_QUICKACK'
ss -tie dst <peer>          # 'ato:' is the delayed-ACK timeout in use
# Reproduce: nc without NODELAY vs a tiny python client that sets it:
python3 - <<'EOF'
import socket
s=socket.create_connection(("example.com",80))
s.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
EOF
```

---

## 34. Keepalive, and detecting dead peers

### 34.1 The problem

An idle `ESTABLISHED` TCP connection sends **nothing**. If the peer host crashes,
loses power, or a NAT/firewall silently drops the mapping, **neither side knows**
until it next tries to send data (and then eats retransmits + RTO for minutes) or
forever (if it's only reading).

### 34.2 TCP keepalive

When enabled (`SO_KEEPALIVE` on the socket), after the connection is **idle** for
`tcp_keepalive_time`, the kernel sends a **keepalive probe** (a segment with
`seq = snd_nxt − 1` and no data — deliberately "old" so the peer must ACK it).

| sysctl | Default | Per-socket option | Meaning |
|---|---|---|---|
| `net.ipv4.tcp_keepalive_time` | **7200** s (2 h) | `TCP_KEEPIDLE` | idle time before first probe |
| `net.ipv4.tcp_keepalive_intvl` | **75** s | `TCP_KEEPINTVL` | gap between probes |
| `net.ipv4.tcp_keepalive_probes` | **9** | `TCP_KEEPCNT` | unanswered probes before declaring the connection dead (→ `ETIMEDOUT`) |

Defaults ⇒ **~2 h 11 min** to detect a dead peer. Almost always too slow. The
sysctls are system-wide; the `TCP_KEEP*` socket options let each app choose
(e.g. 60 s idle, 10 s interval, 3 probes ⇒ ~90 s detection).

### 34.3 `TCP_USER_TIMEOUT` — the better knob

`TCP_USER_TIMEOUT` (ms) bounds how long **transmitted but unacked** data may
remain before the connection is killed — this also caps the retransmission
sequence for *active* connections (not just idle ones), overriding
`tcp_retries2`. Combined with keepalive it gives predictable failure timing:
"this connection will error out within N seconds of the peer becoming
unreachable, whether idle or busy." Widely used by databases and gRPC.

### 34.4 Application-level keepalive / heartbeats

Often better than TCP keepalive because:
- It traverses proxies/load balancers that terminate TCP (each hop has its own
  connection; TCP keepalive on segment A tells you nothing about segment C).
- It can carry liveness *of the application*, not just the socket (an app can be
  deadlocked while its kernel still ACKs).
- Protocol examples: HTTP/2 `PING` frames, gRPC keepalive pings, WebSocket
  ping/pong, SSH `ServerAliveInterval`, MQTT `PINGREQ`, Redis `TCP keepalive` +
  `PING`.

### 34.5 Keepalive and NAT/firewall idle timeouts

Stateful middleboxes drop idle flows: common TCP idle timeouts are 5 min
(consumer NAT), 30–60 min (enterprise firewalls), and AWS **NLB is 350 s**, ELB
idle default **60 s**. If your keepalive/heartbeat interval is **longer** than
the shortest middlebox timeout on the path, the flow gets silently reaped and the
next send fails. Rule: heartbeat interval **< min middlebox idle timeout**
(e.g. ≤ 60 s to be safe behind unknown LBs; ≤ 30 s for UDP).

### 34.6 War story

A microservice held a pool of 50 long-lived gRPC connections to a database proxy
through an AWS NLB. Every night during a traffic lull, the first queries after
~6 minutes of idle failed with `connection reset`, then recovered. NLB's 350 s
idle timeout silently dropped the flows; the client only found out on the next
RPC. Fix: gRPC keepalive `time=60s, timeout=20s, permitWithoutStream=true` so the
flows never idle past 350 s. **Lesson:** long-lived connections need a heartbeat
shorter than every idle timeout on the path, or they die quietly during quiet
periods.

**Try it yourself**
```bash
sysctl net.ipv4.tcp_keepalive_time net.ipv4.tcp_keepalive_intvl net.ipv4.tcp_keepalive_probes
ss -tieo dst <peer>        # 'timer:(keepalive,...)' shows keepalive countdown per socket
# ssh keepalive:
grep -E 'ServerAlive' ~/.ssh/config     # ServerAliveInterval 30 / ServerAliveCountMax 3
# Set aggressive keepalive in Python:
python3 - <<'EOF'
import socket
s = socket.create_connection(("example.com", 80))
s.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
s.setsockopt(socket.IPPROTO_TCP, socket.TCP_KEEPIDLE, 30)
s.setsockopt(socket.IPPROTO_TCP, socket.TCP_KEEPINTVL, 10)
s.setsockopt(socket.IPPROTO_TCP, socket.TCP_KEEPCNT, 3)
EOF
```

---
# Part X — TCP performance engineering

## 35. Bandwidth-delay product and buffer sizing

### 35.1 The BDP

```
BDP (bytes) = bottleneck_bandwidth (bytes/s) x round_trip_time (s)
```

It's the amount of data "in the pipe" when the pipe is full. To keep a link
100% utilized, the sender must be allowed to have **at least one BDP**
outstanding (unacked) — so `min(cwnd, rwnd, sndbuf) ≥ BDP`.

| Path | Bandwidth | RTT | BDP |
|---|---|---|---|
| Same rack | 25 Gbit/s | 0.1 ms | 312 KB |
| Same region (AZ↔AZ) | 10 Gbit/s | 1 ms | 1.25 MB |
| Cross-country (US) | 1 Gbit/s | 70 ms | 8.75 MB |
| Transatlantic | 1 Gbit/s | 90 ms | 11.25 MB |
| Home ↔ CDN | 200 Mbit/s | 15 ms | 375 KB |
| Satellite (GEO) | 100 Mbit/s | 600 ms | 7.5 MB |
| LTE | 50 Mbit/s | 50 ms | 312 KB |

### 35.2 Turning BDP into settings

For a single flow to fill a **cross-country 1 Gbit/s** path you need ~9 MB of
window. That requires:
1. **Window scaling** on (Ch. 30) — 16-bit window alone maxes at 64 KB.
2. **Receive buffer max** ≥ BDP: `net.ipv4.tcp_rmem` third value ≥ ~16 MB
   (headroom over 9 MB; autotuning also reserves ~50% for metadata/overhead — the
   effective window is roughly half of `SO_RCVBUF`).
3. **Send buffer max** ≥ BDP: `net.ipv4.tcp_wmem` third value similarly, plus
   `net.core.wmem_max`.
4. Enough **`cwnd`** — congestion control must be *allowed* to grow there, i.e.
   low enough loss (`√p` term, Ch. 36) or a model-based CC (BBR).

Rule of thumb for a WAN transfer host: set `tcp_rmem`/`tcp_wmem` max to
**2 × BDP_max** for your worst path, and `net.core.rmem_max`/`wmem_max` to match.
Don't set `SO_RCVBUF`/`SO_SNDBUF` in the app — that pins the value and disables
autotuning.

### 35.3 Buffers that are too *big*: bufferbloat

The opposite failure. An over-large buffer at the **bottleneck** (home router,
cable modem, VM host, a qdisc with a 1000-packet FIFO) lets a loss-based CC
(CUBIC) fill it with a full BDP+ of *standing queue*. Throughput stays high but
**latency explodes** — a full 1 s of buffering on a slow uplink. Every other flow
sharing that link (your video call, DNS, ACKs) now waits behind the bulk
transfer. Signature: `ping` RTT jumps from 20 ms to 200–1000 ms **only while a
big upload/download runs**.

Fixes:
- **AQM (Active Queue Management)** at the bottleneck: **`fq_codel`** (default
  qdisc on modern Linux) or **CAKE** — they keep the queue short by dropping/ECN-
  marking early and isolate flows so a bulk transfer can't starve a sparse one.
- **BBR** on the senders (paces to BtlBw, keeps ~0 standing queue).
- On a Linux router/gateway: `tc qdisc replace dev <wan> root cake bandwidth
  <slightly-below-line-rate>` (shaping just under line rate moves the queue from
  the dumb modem into the smart Linux qdisc).

### 35.4 The "50%" gotchas

- Linux reports `SO_RCVBUF`/`SO_SNDBUF` as **double** what you set
  (`setsockopt(...,X)` → `getsockopt` returns `2X`); the extra is bookkeeping
  reserve. Effective window ≈ `SO_RCVBUF / 2`.
- `net.ipv4.tcp_adv_win_scale` (default 1 on older, `-2` on newer kernels)
  controls that application-vs-window split.

### 35.5 War story

A media company's transcode farm pulled 4 GB masters from an object store
cross-region and each pull took ~9 minutes (~60 Mbit/s) on a 10 Gbit/s network,
RTT 32 ms. BDP ≈ 40 MB. `tcp_rmem` max was the distro default **6 MB** →
effective window ~3 MB → `3MB / 0.032s ≈ 750 Mbit/s` theoretical, less in
practice. Raising `tcp_rmem`/`tcp_wmem` max to 128 MB and `core.rmem_max`/
`wmem_max` to 128 MB took a single stream to ~6 Gbit/s (pull time ~1 min).
**Lesson:** for WAN bulk transfer, if `buffer_max < 2×BDP` you are buffer-bound
no matter the link speed. Compute BDP first.

**Try it yourself**
```bash
# Measure BDP inputs:
ping -c20 <peer> | tail -1                     # avg RTT
iperf3 -c <peer> -t20                          # throughput (single stream)
iperf3 -c <peer> -t20 -P8                      # with 8 parallel streams (window/CC limited?)
# Current ceilings:
sysctl net.ipv4.tcp_rmem net.ipv4.tcp_wmem net.core.rmem_max net.core.wmem_max
ss -tim dst <peer>                             # watch rcv_space / snd_wnd grow
tc -s qdisc show dev eth0                      # is fq_codel/fq/cake active? drops?
```

---

## 36. Throughput math: loss, latency, and the Mathis equation

### 36.1 The equations every network engineer should know

**Window-limited throughput (no loss):**
```
Throughput ~ Window / RTT
```

**Loss-limited throughput of a single loss-based TCP flow (Mathis et al., 1997):**
```
Throughput ~ (MSS / RTT) x (C / sqrt(p))

  p = packet loss probability
  C ~ 1.22  (for the standard AIMD sawtooth with periodic loss)
```

A more complete form (Padhye/PFTK, 1998) adds the RTO term for heavier loss:
```
Throughput ~ MSS /
   ( RTT*sqrt(2p/3)  +  RTO*min(1, 3*sqrt(3p/8))*p*(1 + 32p^2) )
```

### 36.2 What Mathis tells you (the intuition)

- Throughput is **inversely proportional to RTT** — double the latency, halve the
  speed.
- Throughput is **inversely proportional to √loss** — going from 0.01% to 1% loss
  (100×) cuts throughput **10×**.
- A loss-based flow **cannot** sustain a large window on a path with even modest
  loss, regardless of link capacity.

Worked example — 1 Gbit/s link, RTT 80 ms, MSS 1448:

| Loss `p` | Mathis throughput (1 flow) | % of 1 Gbit/s |
|---|---|---|
| 0 (window-limited by 8 MB) | ~800 Mbit/s | 80% |
| 0.001% (1e-5) | `1448/0.08 × 1.22/√1e-5` ≈ **69.8 Mbit/s** | 7% |
| 0.01% (1e-4) | ≈ **22 Mbit/s** | 2.2% |
| 0.1% (1e-3) | ≈ **7 Mbit/s** | 0.7% |
| 1% (1e-2) | ≈ **2.2 Mbit/s** | 0.2% |

This is why a "1 Gbit/s" trans-continental link delivers single-digit Mbit/s per
CUBIC stream once there's any real loss — and why the workarounds are: **many
parallel streams** (N streams ≈ N× throughput until they collectively cause
loss), **BBR** (models bandwidth, ignores that loss), **FEC/erasure coding**, or
**fixing the loss**.

### 36.3 Parallel streams

`N` independent flows sharing a path get roughly `N ×` a single flow's Mathis
rate — until their combined load pushes `p` up. Tools: `iperf3 -P`, `aria2c -x`,
GridFTP, `rclone --multi-thread-streams`, Aspera/UDT (UDP-based, bypasses the
√p problem entirely). This is also why a browser's 6 connections and HTTP/2's
multiplexing-over-one-connection behave differently under loss (HTTP/2's single
TCP connection suffers head-of-line blocking; HTTP/3 fixes it — Ch. 39).

### 36.4 Goodput vs throughput

- **Throughput:** bytes on the wire per second (includes headers, retransmits).
- **Goodput:** *useful application* bytes per second delivered.
- Overhead sources: 40–52 B headers per ~1448 B (~3.5%), retransmissions (2×p
  roughly), TLS record overhead (~0.5–1%), slow-start ramp (matters for short
  transfers), and connection setup RTTs.
- For a **short** transfer (e.g. 50 KB HTTP response), you never leave slow start;
  throughput math is irrelevant and **round trips dominate**: `time ≈ (handshake
  RTTs) + ceil(log2(response/IW)) × RTT`. Optimize RTT count, not bandwidth.

### 36.5 War story

A genomics lab moved ~2 TB/day between universities over a shared 10 Gbit/s
research link, RTT 55 ms, measured loss ~0.02%. Single `scp` ≈ 15 Mbit/s (Mathis
predicts ~12–20). They blamed the link. Solution stack: (1) `tcp_rmem`/`wmem` max
→ 256 MB, (2) switch senders to **BBR + fq**, (3) use a parallel tool
(`bbcp`/GridFTP, 16 streams). Result: ~8 Gbit/s aggregate, transfer window shrank
from all day to ~45 min. **Lesson:** don't argue with `MSS/(RTT·√p)`; engineer
around it with buffers + BBR + parallelism.

**Try it yourself**
```bash
# Estimate p from counters over a known interval:
nstat -az | grep -E 'TcpRetransSegs|TcpOutSegs'   # p ~ RetransSegs / OutSegs (rough)
# Compare 1 vs N streams (if Nx faster, you're loss/CC-limited, not link-limited):
iperf3 -c <peer> -t30 ; iperf3 -c <peer> -t30 -P16
# Compare CUBIC vs BBR on the same lossy path:
iperf3 -c <peer> -t30 -C cubic ; iperf3 -c <peer> -t30 -C bbr
```

---

## 37. Tuning knobs on Linux (sysctl reference)

> Apply with `sysctl -w key=value` (runtime) and persist in
> `/etc/sysctl.d/99-net.conf`. **Measure before and after.** Defaults are good for
> most hosts; tune deliberately for a known workload (WAN bulk transfer, high-
> connection-count edge proxy, low-latency RPC).

### 37.1 Core buffers & backlog

| sysctl | Typical default | When to raise | Notes |
|---|---|---|---|
| `net.core.rmem_max` / `wmem_max` | 208 KB–4 MB | WAN bulk transfer | absolute cap for `SO_RCVBUF`/`SO_SNDBUF` and TCP autotune max |
| `net.core.netdev_max_backlog` | 1000 | 10G+ NICs, softirq drops | per-CPU queue between NIC and stack; raise to 5000–30000 if `nstat` shows `netdev_backlog` drops |
| `net.core.somaxconn` | 4096 (older: 128) | busy accept()ers | must also pass a large `backlog` to `listen()` |
| `net.core.default_qdisc` | `fq_codel` | `fq` for BBR / pacing | bufferbloat control |

### 37.2 TCP memory & windows

| sysctl | Meaning |
|---|---|
| `net.ipv4.tcp_rmem` = `min default max` | per-socket **receive** buffer autotuning bounds (bytes). Raise `max` to ~2×BDP for WAN. |
| `net.ipv4.tcp_wmem` = `min default max` | per-socket **send** buffer autotuning bounds. |
| `net.ipv4.tcp_mem` = `low pressure high` (pages!) | **system-wide** TCP memory limits in **4 KB pages**. Under memory pressure the stack shrinks windows and prunes. Raise on hosts with millions of sockets. |
| `net.ipv4.tcp_moderate_rcvbuf` = 1 | enable receive-buffer autotuning (keep on). |
| `net.ipv4.tcp_window_scaling` = 1 | keep on always. |
| `net.ipv4.tcp_adv_win_scale` | app-data vs window split of the rcvbuf. |

### 37.3 Connection setup / teardown

| sysctl | Meaning |
|---|---|
| `net.ipv4.tcp_max_syn_backlog` | SYN-queue size (half-open). Raise to 8192–65536 on public-facing servers. |
| `net.ipv4.tcp_syncookies` = 1 | SYN-flood protection; keep on. |
| `net.ipv4.tcp_syn_retries` / `tcp_synack_retries` | initial SYN / SYN-ACK retransmit count (~6/5 → ~127 s / ~31 s). Lower for faster failover. |
| `net.ipv4.tcp_abort_on_overflow` = 0 | on accept-queue overflow, drop the ACK (client retries) rather than RST. Keep 0 unless you prefer fast failure. |
| `net.ipv4.tcp_fin_timeout` = 60 | bounds `FIN_WAIT_2`. |
| `net.ipv4.tcp_tw_reuse` = 2 | reuse TIME_WAIT for **outbound** connects when safe (timestamps). Set 1 on busy clients. **Never** `tcp_tw_recycle` (removed/broken). |
| `net.ipv4.ip_local_port_range` = `32768 60999` | ephemeral source ports. Widen to `1024 65535` on heavy outbound-connection hosts. |
| `net.ipv4.tcp_max_tw_buckets` | cap on TIME_WAIT sockets (excess are killed early). |
| `net.ipv4.tcp_max_orphans` | cap on orphaned (no fd) sockets; DoS guard. |

### 37.4 Reliability / recovery

| sysctl | Meaning |
|---|---|
| `net.ipv4.tcp_sack` / `tcp_dsack` = 1 | keep on — essential on lossy/long paths. |
| `net.ipv4.tcp_fack` | deprecated (RACK supersedes). |
| `net.ipv4.tcp_recovery` = 1 | enable **RACK-TLP** loss detection (default on modern kernels). |
| `net.ipv4.tcp_early_retrans` | early-retransmit / TLP behavior. |
| `net.ipv4.tcp_frto` = 2 | spurious-RTO detection & undo. |
| `net.ipv4.tcp_retries2` = 15 | ~13–30 min before killing an established connection. Lower (e.g. 8 ≈ ~100 s) for faster failover clusters; or use `TCP_USER_TIMEOUT` per socket. |
| `net.ipv4.tcp_slow_start_after_idle` = 1 | reset `cwnd` after an idle period. **Set 0** for long-lived keep-alive connections / HTTP/2 origins to avoid re-ramping after every idle gap. |
| `net.ipv4.tcp_notsent_lowat` | limit unsent bytes in the send buffer → lower latency for apps that write large buffers then want to react (e.g. adaptive video). Try `131072`. |

### 37.5 Congestion control / pacing / ECN

| sysctl | Meaning |
|---|---|
| `net.ipv4.tcp_congestion_control` | `cubic` (default) / `bbr` / `dctcp` (DC only). |
| `net.ipv4.tcp_available_congestion_control` | loaded modules. |
| `net.core.default_qdisc` = `fq` | needed for BBR pacing (or `fq_codel`). |
| `net.ipv4.tcp_ecn` = 2 | accept ECN if peer requests; `1` also actively requests it. |
| `net.ipv4.tcp_pacing_ss_ratio` / `tcp_pacing_ca_ratio` | pacing aggressiveness in slow start / congestion avoidance. |

### 37.6 Offloads, routing, misc

| sysctl / tool | Meaning |
|---|---|
| `net.ipv4.tcp_mtu_probing` = 1 (or 2) | **PLPMTUD**: probe for MTU when a black hole is suspected — mitigates ICMP-filtered paths (Ch. 8). `2` with a set base MSS is aggressive. |
| `net.ipv4.tcp_timestamps` = 1 | keep on (RTT accuracy + PAWS); `0` only if a middlebox mangles them. |
| `net.ipv4.tcp_fastopen` = 3 | enable TFO client+server. |
| `net.ipv4.ip_forward` = 1 | make the host a router. |
| `net.ipv4.conf.all.rp_filter` = 1/2 | reverse-path filter: `1` strict (drops asymmetric), `2` loose (multi-homed hosts). |
| `net.netfilter.nf_conntrack_max` | firewall state table size; raise on NAT gateways/LBs or new flows get dropped (`nf_conntrack: table full`). |
| `net.ipv4.neigh.default.gc_thresh1/2/3` | ARP/ND cache GC thresholds; raise on hosts with thousands of neighbors or you get `neighbour table overflow`. |
| `fs.file-max`, `ulimit -n` | fd ceilings — a million connections needs a million fds. |

### 37.7 A sample profile: WAN bulk-transfer host

```conf
# /etc/sysctl.d/99-wan-transfer.conf
net.core.rmem_max            = 268435456
net.core.wmem_max            = 268435456
net.ipv4.tcp_rmem            = 4096 131072 268435456
net.ipv4.tcp_wmem            = 4096 131072 268435456
net.ipv4.tcp_mem            = 786432 1048576 26777216
net.core.default_qdisc       = fq
net.ipv4.tcp_congestion_control = bbr
net.ipv4.tcp_slow_start_after_idle = 0
net.ipv4.tcp_mtu_probing     = 1
net.core.netdev_max_backlog  = 30000
```

### 37.8 A sample profile: high-connection-count edge proxy

```conf
net.core.somaxconn            = 65535
net.ipv4.tcp_max_syn_backlog  = 65535
net.ipv4.ip_local_port_range  = 1024 65535
net.ipv4.tcp_tw_reuse         = 1
net.ipv4.tcp_fin_timeout      = 15
net.ipv4.tcp_slow_start_after_idle = 0
net.netfilter.nf_conntrack_max = 2000000
fs.file-max                   = 3000000
# plus: ulimit -n 1048576 for the proxy process
```

### 37.9 War story

An SRE copied a "10 Gbit tuning" blog's sysctls onto every host in the fleet,
including tiny 512 MB edge boxes. `net.ipv4.tcp_mem` was set to huge values and
`tcp_rmem` max to 64 MB. Under a connection surge those boxes OOM-killed the proxy
because a few thousand connections × tens of MB of window blew past RAM.
**Lesson:** buffer maxima are *per socket*; on memory-constrained or
high-fan-in hosts, keep `tcp_rmem`/`wmem` max modest and cap total with
`tcp_mem`. Tune per role, not per blog post.

**Try it yourself**
```bash
sysctl -a --pattern 'net.(ipv4.tcp|core)'      # dump all current values
nstat -az                                       # the counters to watch before/after
watch -n1 'nstat -az | grep -E "ListenDrops|ListenOverflows|TCPBacklogDrop|PruneCalled|RcvPruned|OfoPruned|TW"'
cat /proc/net/sockstat                           # TCP mem pages in use vs tcp_mem
```
---


# Part XI — Application-layer transports

## 38. TLS: handshakes, versions, resumption, termination, and offload

TLS (Transport Layer Security) gives an application **confidentiality**
(eavesdroppers see ciphertext), **integrity** (tampering is detected), and
**authentication** (you're really talking to `example.com`, and optionally it's
really talking to *you*). It does **not** change TCP — it is simply the first
bytes the application writes into the byte stream, and everything after the
handshake is encrypted records. QUIC (Ch. 39, and §38.13 here) reuses the TLS 1.3
*handshake* but replaces the *record layer* with its own packet protection.

This chapter, section by section:
- 38.1 Where TLS sits — the record layer
- 38.2 Version history: SSL 2/3 → TLS 1.0/1.1/1.2/1.3 (what changed, what's dead)
- 38.3 Cryptographic building blocks (so the handshakes make sense)
- 38.4 The TLS 1.2 handshake, flight by flight
- 38.5 The TLS 1.3 handshake, flight by flight (+ HelloRetryRequest, key schedule)
- 38.6 0-RTT / early data, and its replay danger
- 38.7 Session resumption: session IDs, session tickets, 1.3 PSK, fleet key management
- 38.8 Certificates, chains, SNI, ECH, OCSP stapling, CT, CAA
- 38.9 The record layer in detail: AEAD, sizing, limits, KeyUpdate, alerts
- 38.10 mTLS (mutual TLS) and client authentication
- 38.11 TLS termination architectures: edge, passthrough, re-encrypt, mesh
- 38.12 TLS offload: software, session cache, kTLS, NIC inline TLS, QAT/HSM
- 38.13 QUIC's use of TLS 1.3
- 38.14 Interactions with TCP that bite
- 38.15 Debugging TLS + war stories

### 38.1 Where TLS sits — the record layer

```
[ TCP payload bytes ] = [ TLS record ][ TLS record ][ TLS record ] ...

TLS 1.2 record:  ContentType(1) | LegacyVersion(2) | Length(2) | payload(+MAC or AEAD tag)
TLS 1.3 record:  opaque_type=23 | legacy_version=0x0303 | Length(2) | ciphertext
                 (real content type is the last byte *inside* the decrypted blob)

ContentType: 20 change_cipher_spec, 21 alert, 22 handshake, 23 application_data,
             24 heartbeat (RFC 6520, usually disabled - see Heartbleed)
```

- One record carries **≤ 2^14 = 16384 bytes** of plaintext (`max_fragment_length`
  / `record_size_limit` extensions can lower it).
- **Records do not align to TCP segments.** A record can span many segments; a
  segment can hold several records. The receiver must buffer until it has a whole
  record before it can decrypt/authenticate it — this is why a single lost TCP
  segment stalls *every* TLS record queued behind it (TCP head-of-line blocking,
  §38.14).
- After the handshake, everything (including HTTP/2 frames, close signals, even
  the real content-type in 1.3) is inside `application_data` records. On the wire
  you see only `type 23`, a length, and ciphertext.

### 38.2 Version history — what changed and what is dead

| Version | Year | Status | Key points |
|---|---|---|---|
| **SSL 2.0** | 1995 | **Dead / forbidden** (RFC 6176) | no handshake integrity, weak MACs, one key for everything. Never enable. |
| **SSL 3.0** | 1996 | **Dead** (RFC 7568) | new design, but CBC padding is unauthenticated → **POODLE** (2014). Disabled everywhere. |
| **TLS 1.0** | 1999 (RFC 2246) | **Deprecated** (RFC 8996, 2021) | SSL 3.1 essentially. Vulnerable to **BEAST** (CBC IV), needs `1/n-1` splitting workaround. PCI-DSS banned it. |
| **TLS 1.1** | 2006 (RFC 4346) | **Deprecated** (RFC 8996) | fixed the BEAST IV issue (explicit per-record IV). No AEAD, no SHA-256, no modern extensions. Browsers removed support in 2020. Almost nothing needs it today. |
| **TLS 1.2** | 2008 (RFC 5246) | **Widely used, still fine** | adds **AEAD** ciphers (AES-GCM, ChaCha20-Poly1305), SHA-256+ PRF, negotiable hash for signatures, extensions ecosystem (SNI, ALPN, OCSP stapling, EC curves). Also still *permits* RSA key transport, CBC-mode ciphers, compression, renegotiation, and weak curves — so real-world security depends on **configuration**. |
| **TLS 1.3** | 2018 (RFC 8446) | **Preferred** | ground-up redesign: **1-RTT** (or 0-RTT), **forward secrecy mandatory** (ephemeral (EC)DHE only — no RSA key transport), only 5 AEAD cipher suites, encrypted handshake from ServerHello onward, **removed**: renegotiation, compression, static RSA, CBC, custom DH groups, `ChangeCipherSpec` semantics (kept as a dummy for middlebox compatibility), MD5/SHA-1, RC4, DSA, export ciphers. HKDF-based key schedule. |
| **TLS 1.3 + PQC** | 2024– | rolling out | hybrid key exchange (**X25519MLKEM768**, formerly X25519Kyber768) to resist "harvest now, decrypt later." Bigger ClientHello/ServerHello → see §38.14 MTU. |

**What CBC-mode attacks taught the design (why 1.3 is minimal):** BEAST, Lucky13,
POODLE, and the padding-oracle family all exploited "MAC-then-encrypt" CBC in
TLS ≤1.2. **Compression** attacks (**CRIME** on TLS compression, **BREACH** on
HTTP compression) killed TLS-level compression. **RSA key transport** with PKCS#1
v1.5 gave the **Bleichenbacher / ROBOT** oracles and has no forward secrecy.
TLS 1.3 removed *every one* of these options rather than trusting operators to
disable them.

**Practical policy today:** enable **TLS 1.2 + 1.3 only**; on 1.2 allow only
ECDHE key exchange + AES-GCM/ChaCha20-Poly1305; disable renegotiation,
compression, and session-ID caching you don't manage. Test with `testssl.sh` or
SSL Labs.

### 38.3 Cryptographic building blocks (quick reference)

| Piece | Job | Examples |
|---|---|---|
| **Key exchange / KEM** | agree a shared secret over a public channel, ideally with **forward secrecy** | ECDHE (X25519, secp256r1/384r1), FFDHE, ML-KEM (post-quantum), hybrid X25519MLKEM768 |
| **Authentication (signature)** | prove identity via the certificate's private key | RSA-PSS, RSA-PKCS1 (1.2 only), ECDSA (P-256/384), Ed25519 |
| **AEAD cipher** | encrypt + authenticate record payloads with one primitive | AES-128-GCM, AES-256-GCM, ChaCha20-Poly1305, AES-CCM |
| **Hash / KDF** | key derivation, transcript hash, HMAC | SHA-256, SHA-384; HKDF (extract+expand) in 1.3 |
| **Transcript hash** | bind every handshake message so tampering breaks `Finished` | running hash of all handshake bytes |

A **cipher suite** names the combination. TLS 1.2:
`TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256` = ECDHE key exchange, RSA cert auth,
AES-128-GCM AEAD, SHA-256 KDF/PRF. TLS 1.3 shortened it to just the AEAD+hash:
`TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`, `TLS_CHACHA20_POLY1305_SHA256`
(key exchange and signature are negotiated in separate extensions).

**Forward secrecy (a.k.a. PFS):** because ECDHE uses a *fresh ephemeral* key pair
per handshake, an attacker who records ciphertext today and steals the server's
long-term private key next year still cannot decrypt it. Static-RSA key transport
had no PFS — one stolen key retroactively decrypts years of traffic. This is the
single biggest reason to be on 1.3 (or 1.2-ECDHE-only).

### 38.4 The TLS 1.2 handshake, flight by flight

```
   Client                                        Server
   ------                                        ------

   ClientHello                    -------->
     client_random, session_id, cipher_suites[], compression=null
     extensions: server_name(SNI), supported_groups, ec_point_formats,
                 signature_algorithms, ALPN, status_request(OCSP),
                 renegotiation_info, [session_ticket], extended_master_secret

                                                 ServerHello
                                                   server_random, chosen cipher, session_id
                                                 Certificate            (server chain)
                                                 [CertificateStatus]    (stapled OCSP)
                                                 ServerKeyExchange      (server ECDHE pubkey,
                                                                         signed with the cert key)
                                                 [CertificateRequest]   (mTLS only)
                                  <--------       ServerHelloDone

   [Certificate]                                  (client chain, mTLS only)
   ClientKeyExchange                              (client ECDHE pubkey)
   [CertificateVerify]                            (client signs transcript, mTLS)
   [ChangeCipherSpec]                             (everything after this is encrypted)
   Finished                       -------->       (PRF over transcript, under new keys)

                                                 [NewSessionTicket]
                                                 [ChangeCipherSpec]
                                  <--------       Finished

   Application Data (HTTP request) <------->      Application Data (HTTP response)
```

Key computation:
```
premaster_secret = ECDHE(client_ephemeral_priv, server_ephemeral_pub)   # both sides get the same value
master_secret    = PRF(premaster_secret, "master secret"
                          [ or "extended master secret", session_hash ], client_random + server_random)[0..48]
key_block        = PRF(master_secret, "key expansion", server_random + client_random)
                 -> client_write_MAC_key, server_write_MAC_key (CBC suites only),
                   client_write_key, server_write_key,
                   client_write_IV, server_write_IV
```

Notes that matter operationally:
- **2 round trips** before the client can send the request (ClientHello→
  ServerHelloDone is RTT 1; ClientKeyExchange/Finished→server Finished is RTT 2).
  Plus the TCP handshake = **3 RTT to first byte** on a cold connection.
- **`ServerKeyExchange` is signed** by the certificate's private key — that's the
  *authentication*. If the signature or curve is weak, the handshake is weak even
  if the AEAD is strong.
- **Extended Master Secret (EMS, RFC 7627)** binds the master secret to the full
  handshake transcript, closing the **triple-handshake** attack. Enable it; some
  old stacks and middleboxes break resumption without it.
- **Renegotiation**: either side can restart a handshake mid-connection (client
  cert step-up, rekey). The 2009 **plaintext injection** flaw (CVE-2009-3555) was
  patched with the `renegotiation_info` extension (RFC 5746); still, disable
  client-initiated renegotiation (it's a cheap DoS: forces server asymmetric
  crypto).
- CBC-mode 1.2 suites are still standardized; prefer GCM/ChaCha20 and disable CBC
  where you can (Lucky13 timing).

### 38.5 The TLS 1.3 handshake, flight by flight

```
   Client                                        Server
   ------                                        ------

   ClientHello                    -------->
     client_random, legacy_session_id, cipher_suites=[AES_128_GCM_SHA256, ...]
     extensions:
       supported_versions = [TLS 1.3]            <- the real version lives here
       supported_groups   = [x25519, secp256r1, ...]
       key_share          = x25519 client ephemeral pubkey
       signature_algorithms = [ecdsa_secp256r1_sha256, ...]
       server_name(SNI), ALPN, [pre_shared_key],
       [early_data], [psk_key_exchange_modes], [cookie]

                   (server derives the (EC)DHE shared secret from key_share)

                                                 ServerHello
                                                   server_random, chosen cipher,
                                                   key_share = server ephemeral pubkey
                   (both sides now hold handshake_traffic keys; rest is ENCRYPTED)
                                                 {EncryptedExtensions}
                                                 {CertificateRequest}   (mTLS only)
                                                 {Certificate}          (chain, +OCSP, +SCT)
                                                 {CertificateVerify}    (sig over transcript)
                                  <--------       {Finished}            (HMAC over transcript)
                                                 [ChangeCipherSpec]     (dummy, middlebox compat)

   {Certificate}{CertificateVerify}              (mTLS only)
   {Finished}                      -------->
   -- 1-RTT: client may send Application Data with/after its Finished -->
   <------- 0.5-RTT: server may send Application Data right after its Finished

                                  <--------       {NewSessionTicket} x N   (for resumption)
```

- **1 round trip** to first byte (client can send the HTTP request in the same
  flight as its `Finished`). With TCP that's **2 RTT cold**, **1 RTT** if TCP
  Fast Open is used, **0 RTT** with a PSK + early data (§38.6).
- **The certificate and everything after ServerHello is encrypted** under
  *handshake* keys — a passive observer sees SNI (unless ECH) and little else, not
  the server cert.
- **HelloRetryRequest (HRR):** if the client's `key_share` used a group the
  server won't use, the server replies HRR naming an acceptable group; the client
  resends ClientHello with a new `key_share`. Adds **1 RTT**. Tune the client's
  default group (X25519) to match servers and avoid HRR. Also used to force a
  `cookie` for DoS resistance (stateless retry).
- **Key schedule (HKDF):**
  ```
  Early Secret     = HKDF-Extract(0, PSK or 0)
  (binder/ext keys derived here for 0-RTT)
  Handshake Secret = HKDF-Extract(DeriveSecret(Early, "derived"), (EC)DHE)
    -> client_handshake_traffic_secret, server_handshake_traffic_secret
  Master Secret    = HKDF-Extract(DeriveSecret(Handshake, "derived"), 0)
    -> client_application_traffic_secret_0, server_application_traffic_secret_0
    -> exporter_master_secret, resumption_master_secret
  ```
  Traffic keys/IVs come from `HKDF-Expand-Label` of the relevant secret. After
  the handshake, `KeyUpdate` ratchets the application secret forward
  (`application_traffic_secret_{N+1} = HKDF-Expand-Label(secret_N, "traffic upd")`).
- **`CertificateVerify`** = the server signs the transcript hash with the cert's
  private key. That single signature is the entire proof of identity in 1.3
  (no `ServerKeyExchange`).
- **`Finished`** = HMAC over the transcript with a key derived from the traffic
  secret; a mismatch means the handshake was tampered with (downgrade,
  bit-flip) → `decrypt_error` alert, connection aborted.

### 38.6 0-RTT / early data

On a repeat visit where the client holds a PSK (from a prior `NewSessionTicket`),
the client can send **`early_data`** (the HTTP request itself) in the **first
flight**, encrypted under a key derived from the PSK — **zero round trips** to
first byte.

```
ClientHello (+ pre_shared_key, + early_data ext)
{early_data: "GET /home HTTP/1.1 ..."}          <- sent immediately, PSK-encrypted
--------------------------------------------->
                          ServerHello ...(accepts or rejects early_data)... Finished
                          {EncryptedExtensions: early_data}   <- "I accepted it"
<---------------------------------------------
{end_of_early_data}{Finished}  then 1-RTT app data
```

**The replay problem:** 0-RTT data has **no freshness guarantee**. A network
attacker can capture the early-data flight and **resend it** to the server (or to
another server in the same cluster sharing ticket keys), N times. TLS itself
cannot stop this. Therefore:
- Only send **idempotent, side-effect-free** requests as early data (GET without
  auth-state changes). Never `POST /transfer-money`.
- Servers **must** implement anti-replay: single-use ticket tracking, or a
  bounded replay window using the ticket age + a cluster-wide strike register, or
  restrict 0-RTT to safe methods, or disable it. Cloudflare/CDNs typically allow
  0-RTT only for GET/HEAD and strip it for anything with cookies.
- `early_data` also can't change the negotiated parameters — ALPN and cipher must
  match the original session.

Benefit is real for latency (a whole RTT on every repeat connection), so it's
common on CDNs; just gate it hard on the server side.

### 38.7 Session resumption — and managing it across a fleet

Full handshakes cost round trips **and** asymmetric crypto (an ECDSA sign or RSA
decrypt per connection). Resumption skips the certificate + key-exchange proof by
reusing a secret from a previous session.

**Mechanism 1 — Session IDs (TLS ≤1.2, stateful).** Server invents a
`session_id`, stores `{id → master_secret, cipher, ...}` in a **server-side
cache**. Client presents the id in a later ClientHello; server looks it up and
both skip to an abbreviated handshake (1 RTT). Problems: the server must hold
state (RAM), and in a load-balanced fleet the *same* server must get the resumed
connection (sticky routing) or share the cache (Memcached/Redis) — otherwise it
falls back to a full handshake.

**Mechanism 2 — Session Tickets (RFC 5077, TLS 1.2; stateless).** The server
encrypts the session state into an opaque **ticket** with a server-held
**Session Ticket Encryption Key (STEK / "ticket key")** and sends
`NewSessionTicket`. The client stores the ticket and presents it later; any
server that holds the same STEK can decrypt it and resume — **no shared cache,
no stickiness needed.**

**Mechanism 3 — TLS 1.3 PSK resumption.** Unifies the above: after a handshake
the server sends one or more `NewSessionTicket` messages, each establishing a
**PSK** (identified by the ticket) with an associated `ticket_age_add` and
lifetime. A later ClientHello carries `pre_shared_key` (+ **binder**: an HMAC
proving the client knows the PSK) and, ideally, a fresh `key_share` too so the
resumed session *still has forward secrecy* (`psk_dhe_ke` mode). This is also the
gateway to **0-RTT** (§38.6).

**Fleet key management (the operational core):**
- All terminators that might receive a resumed connection must share the current
  **STEK / ticket key** set. Distribute out-of-band (config management, a secrets
  store, `nginx ssl_session_ticket_key` files, HAProxy `tls-ticket-keys`, Envoy
  SDS).
- **Rotate frequently** — daily or faster. A STEK that never rotates is a
  **forward-secrecy hole**: stealing today's STEK lets an attacker decrypt every
  session that resumed under it, including recorded past traffic (this is exactly
  the mistake that made Twitter/others' PFS "not really PFS" in ~2013, and what
  the 2018 discourse around long-lived ticket keys was about).
- Keep a **small ring**: e.g. accept tickets under keys `{N, N-1, N-2}`, encrypt
  new tickets only under `N`. Roll `N` on a timer; drop the oldest.
- Cap **ticket lifetime** (`ssl_session_timeout`, 1.3 `NewSessionTicket`
  lifetime) to hours, not days.
- If you can't manage rotation safely, **disabling tickets** and taking the full-
  handshake cost can be the more secure choice.
- CDNs solve this by terminating in one security domain per POP with automated
  key rotation.

**Resumption vs security summary:**

| Method | State | LB requirement | Forward secrecy of the resumed session |
|---|---|---|---|
| Session ID cache | server RAM | sticky or shared cache | yes (reuses PFS master secret; cache theft exposes it) |
| Session ticket (1.2) | STEK only | shared STEK | **only as good as STEK rotation** |
| 1.3 PSK, `psk_dhe_ke` | STEK/PSK | shared STEK | **yes** (fresh ECDHE each resume) |
| 1.3 PSK, `psk_ke` (no DHE) | STEK/PSK | shared STEK | no |
| 0-RTT early data | STEK/PSK | shared STEK | no (and replayable) |

### 38.8 Certificates, chains, SNI, ECH, OCSP, CT, CAA

**The chain.** The server sends its **leaf** cert + **intermediate(s)**; the
client already trusts a **root** in its trust store. Validation:
1. Build a path leaf → intermediate(s) → trusted root (signatures verify at each
   link).
2. Each cert in-date (`notBefore`/`notAfter`); clock skew and expiry are the #1
   real-world TLS outage.
3. Leaf's **Subject Alternative Name (SAN)** matches the hostname the client
   wanted (CN is ignored by modern clients). Wildcards match one label.
4. Key usage / extended key usage allows `serverAuth`.
5. Not revoked (see OCSP/CRL below).
6. Meets policy: **CT** (below), pinning if any, min key size, allowed sig algs.

**Missing intermediate** is a classic: works in browsers (they cache or fetch via
AIA) but fails in `curl`, Java, Go, mobile. Always serve the full chain; test
from a clean client.

**SNI (Server Name Indication).** The client puts the target hostname in a
ClientHello extension so one IP + one terminator can serve **many certs / many
sites** (virtual hosting over TLS). In TLS 1.2 and 1.3 alike, SNI is **plaintext**
on the wire — it's the main remaining metadata leak, and what censors and
corporate DLP filter on.

**ECH (Encrypted Client Hello).** Encrypts the sensitive parts of ClientHello
(including SNI and ALPN) to a public key the client fetches from the DNS
**`HTTPS`/`SVCB`** record (`ech=` param). The outer ClientHello carries a
decoy/"public" name; a fronting server decrypts the inner one. Rolling out
(Cloudflare, Firefox); needs DNS + provider support; the whole POP must share the
ECH key.

**Revocation.**
- **CRL** — big signed list of revoked serials; largely impractical at web scale.
- **OCSP** — client asks the CA "is serial X still valid?" Adds latency, a
  privacy leak (CA sees your browsing), and soft-fails (clients proceed if the
  OCSP responder is down → weak).
- **OCSP stapling** (`status_request`) — the **server** periodically fetches a
  signed, time-stamped "good" OCSP response from the CA and **staples** it into
  the handshake (`CertificateStatus` in 1.2, an extension in 1.3). No client→CA
  round trip, no privacy leak. **Operational care:** if your stapling cache
  expires and the CA's OCSP is unreachable, and you set **`OCSP Must-Staple`** on
  the cert, every handshake **fails hard**. Monitor staple freshness; keep a
  fallback; consider short-lived certs instead of Must-Staple.
- Trend: **short-lived certificates** (90 days → 47 days coming → hours) make
  revocation almost moot; automate with ACME.

**Certificate Transparency (CT).** CAs must log every issued cert to public,
append-only logs; the leaf carries **SCTs** (signed timestamps) proving it was
logged, delivered via the cert, a TLS extension, or stapled OCSP. Chrome
*requires* CT. Benefit: you can **monitor CT logs** (crt.sh, certspotter) for
unexpected certs on your domains — early warning of mis-issuance or compromise.

**CAA DNS records.** `example.com. CAA 0 issue "letsencrypt.org"` tells CAs which
CAs may issue for your domain; a compliant CA refuses otherwise. Cheap defense
against mis-issuance.

**Pinning.** HPKP (HTTP header) is **dead** (foot-gun: bricked sites). App-level
pinning (mobile apps pinning a key or intermediate) is still used but risky —
pin the **intermediate or a backup key**, not just the leaf, and ship an
out-of-band update path.

### 38.9 The record layer in detail

**AEAD per record.** For each record: `nonce = static_IV XOR record_sequence_number`
(1.3) or explicit/implicit IV (1.2); `AEAD_encrypt(key, nonce, plaintext,
additional_data = record header)`. Output = ciphertext ‖ 16-byte tag. Any
tampering (or a nonce reuse bug) → tag check fails → fatal `bad_record_mac`
alert. **Sequence numbers are implicit** (never on the wire) and reset to 0 on
each key change; reusing an (key, nonce) pair with GCM is catastrophic
(**forbidden**; a few CVEs came from stacks that did it under concurrency).

**Record size vs latency vs CPU.**
- Big records (16 KB) → fewer AEAD ops, less framing overhead (~0.4%), best bulk
  throughput.
- But the receiver can't use *any* of a record until it's *all* arrived → under
  loss/low bandwidth, a 16 KB record adds delay before the app sees the first
  byte. TTFB-sensitive servers use **dynamic record sizing**: small records
  (~1–4 KB) at the start of a connection / after idle, ramping to 16 KB once
  cwnd is large (nginx `ssl_buffer_size`, Cloudflare's dynamic sizing, Go's
  `DynamicRecordSizingDisabled=false`).
- 1.3 adds optional **record padding** (hide plaintext lengths) — costs bandwidth
  and CPU.

**Limits and rekeying.** Each AEAD has a safe data volume per key (AES-GCM: well
below 2^34.5 records / ~24 GB with a 16 KB record before you should rekey;
ChaCha20-Poly1305 is roomier). TLS 1.3 `KeyUpdate` lets either side ratchet keys
mid-connection without a full renegotiation; long-lived bulk connections (backups,
replication, video) should rekey periodically. `key_update_request` asks the peer
to update too.

**Alerts.** One-byte level (warning/fatal) + description. Common ones:
`close_notify` (graceful shutdown — its absence is a **truncation attack**
signal; libraries differ on whether they treat a missing `close_notify` as an
error), `handshake_failure` (no common params), `protocol_version`,
`bad_certificate` / `certificate_expired` / `unknown_ca`,
`bad_record_mac` (tamper or key desync), `decrypt_error` (bad Finished/sig),
`unrecognized_name` (SNI not served), `no_application_protocol` (ALPN mismatch),
`internal_error`. In 1.3 alerts after ServerHello are encrypted.

**Heartbeat (RFC 6520).** A keepalive/PMTU extension. OpenSSL's unbounded memcpy
on it was **Heartbleed** (2014, CVE-2014-0160) — leaked server memory including
private keys. Almost everyone disables the heartbeat extension now.

### 38.10 Mutual TLS (mTLS) / client authentication

Normally only the server proves identity. **mTLS** makes the **client** present a
certificate too:
- 1.2: server sends `CertificateRequest` (acceptable CAs, sig algs); client
  replies with `Certificate` + `CertificateVerify` (client signs the handshake
  transcript with its private key).
- 1.3: same messages, encrypted; can also happen **post-handshake**
  (`post_handshake_auth`) — e.g. only when the client hits a protected route.

Uses: service-to-service auth in a mesh/zero-trust network, API clients, VPNs,
device identity, admin access. Operational realities:
- You run a **private CA** (or SPIFFE/SPIRE, Vault, cert-manager) and must handle
  issuance, **short lifetimes**, rotation, and revocation for *every* workload.
- The terminator validates the client chain against **your** CA bundle (not the
  public roots) and typically passes the verified identity upstream as a header
  (`X-Client-Cert`, `X-SPIFFE-ID`) — which the backend must trust *only* from the
  terminator (network policy + the terminator stripping any client-supplied copy).
- mTLS + TLS termination interact: if an L7 proxy terminates, the **backend**
  can't see the client cert unless the proxy forwards it, or you use
  **passthrough** / **TLS in TLS**, or the mesh sidecar does mTLS and the app
  speaks plaintext on `localhost` (Istio/Linkerd model).
- SNI-based routing still works; client-cert-based routing needs the proxy to
  parse the client cert.

### 38.11 TLS termination architectures

"Where does the encrypted tunnel end?" determines who can see plaintext, who
holds private keys, and what each hop can route on.

```
Legend:  ==TLS==>  encrypted     -->  plaintext     ==mTLS==>  mutually-authenticated TLS

(A) EDGE / OFFLOAD TERMINATION
  client ==TLS==> [ LB / reverse proxy / CDN ]  --HTTP plaintext-->  backend
     - LB holds the cert + private key, does all crypto
     - LB can route by Host/path/header, add X-Forwarded-For, compress, cache
     - backend is simpler/faster; INTERNAL network must be trusted (or segmented)

(B) TLS PASSTHROUGH  (a.k.a. SNI routing / L4)
  client ============ TLS ============> [ L4 LB ] ============ TLS ============> backend
     - LB never decrypts; peeks at plaintext SNI in ClientHello to pick a backend
     - end-to-end encryption; backend holds the cert; LB can't see HTTP, can't
       add headers, can't do L7 routing or WAF
     - cert/key management is distributed to every backend

(C) RE-ENCRYPT / "TLS bridging"
  client ==TLS(public cert)==> [ LB terminates ] ==TLS(internal cert)==> backend
     - LB decrypts (L7 routing, WAF, observability), then opens a NEW TLS
       connection to the backend using an internal CA
     - two handshakes, two CC loops; most common in regulated environments
     - backend sees the LB as the TLS peer; client identity passed via header or mTLS

(D) SERVICE MESH (sidecar mTLS)
  app --plaintext--> [ sidecar ] ==mTLS==> [ sidecar ] --plaintext--> app
     - app code speaks plaintext to 127.0.0.1; sidecars do mTLS with workload
       identities (SPIFFE). Termination is "everywhere, automatically."
```

Trade-offs:
| | Edge offload (A) | Passthrough (B) | Re-encrypt (C) | Mesh (D) |
|---|---|---|---|---|
| L7 routing / WAF / cache at LB | yes | no | yes | yes (per hop) |
| End-to-end ciphertext | no (plaintext internally) | yes | yes (two legs) | yes (between sidecars) |
| Key mgmt burden | 1 place | every backend | LB + internal CA | automated CA, every pod |
| Backend CPU cost | lowest | full TLS | full TLS (internal) | full mTLS (sidecar) |
| Client-IP visibility to backend | header (XFF / PROXY protocol) | native (real src IP) | header | header from sidecar |
| Sees client cert (mTLS) | at LB | at backend | at LB (forward as header) | at sidecar |

**PROXY protocol.** With L4 passthrough or TCP-mode LBs, the backend would
otherwise see the LB's IP as the client. The **PROXY protocol** (v1 text, v2
binary) prepends one line/header to the TCP stream carrying the real
`src_ip:port` / `dst_ip:port` before the TLS ClientHello. Backend must be
configured to expect it *only* from trusted LB IPs (otherwise it's a spoofing
vector). AWS NLB, HAProxy, Envoy, nginx (`proxy_protocol`) all support it.

### 38.12 TLS offload: software, session cache, kTLS, NIC, QAT/HSM

"Offload" = moving the cost of TLS off the application/backend. Layers, cheapest
change first:

**1. Terminate at a dedicated tier (architectural offload).** Architecture (A)/(C)
above: the reverse proxy / LB / CDN does handshakes and bulk crypto so app
servers don't. This is "TLS offload" in the load-balancer sense. Benefits: fewer
places hold private keys, app servers scale on business logic, central cipher
policy, central cert automation.

**2. Cut handshake cost.** The expensive part of a *full* handshake is the
asymmetric op (ECDSA P-256 sign ≈ tens of µs; RSA-2048 decrypt/sign ≈ ~10× that).
Reduce the number of full handshakes:
- **Session resumption** (§38.7) — a resumed handshake is symmetric-only. Aim for
  a high resume rate; monitor it.
- **Keep-alive / connection reuse / HTTP/2** — amortize one handshake over many
  requests. `tcp_slow_start_after_idle=0` and long idle timeouts help
  (Ch. 34, 37).
- **OCSP stapling cache** shared across the tier (don't let every node hammer the
  CA).
- Prefer **ECDSA certs** over RSA for the leaf (much cheaper signs); you can dual-
  stack ECDSA + RSA and let the client pick.
- **TLS 1.3** (fewer messages, 1-RTT, cheaper).

**3. Bulk-crypto offload — kTLS (kernel TLS).** After the handshake (done in
userspace by OpenSSL/BoringSSL), the app hands the negotiated keys to the kernel
via `setsockopt(SOL_TLS, TLS_TX/TLS_RX, crypto_info)`. The **kernel** then does
AEAD on the socket's data path. Wins:
- **Zero-copy `sendfile()`** of static files/video straight from page cache →
  NIC, encrypted inline, without copying through userspace. Big for CDNs (Netflix
  built kTLS largely for this).
- Fewer syscalls/copies for proxied bulk data.
- Pairs with **NIC inline TLS** (see next) so the CPU does *nothing* for bulk
  records.
Support: Linux `tls` module, nginx `ssl_conf_command Options KTLS` /
build-time, HAProxy, and kTLS-aware stacks. RX-side kTLS is newer/less universal
than TX.

**4. NIC inline TLS ("TLS offload to hardware").** SmartNICs (Mellanox/NVIDIA
ConnectX-6 Dx+, some Intel, Chelsio) implement the record-layer AEAD **on the
NIC**. With kTLS, the kernel programs per-socket keys into the NIC; the NIC
encrypts on TX / decrypts+authenticates on RX as packets flow, and handles
resync on retransmission/reorder. The host CPU is freed from bulk symmetric
crypto entirely — matters at 100 GbE+ video/CDN scale. Falls back to software for
the handshake and for edge cases.

**5. Asymmetric-crypto accelerators — QAT / HSM.**
- **Intel QuickAssist (QAT)** offloads RSA/ECDSA/DH (and bulk crypto, and
  compression) to a co-processor. OpenSSL via the `qatengine`/provider; nginx and
  HAProxy can use it. Turns handshake-bound TLS front-ends (lots of *new*
  connections, low resume rate) from CPU-bound into I/O-bound. Also
  **`AES-NI`/`VAES`/`AVX-512`** and ARM crypto extensions are "offload" in the
  sense of dedicated CPU instructions — always compile your TLS lib to use them
  (10–20× AES throughput vs pure software).
- **HSM (Hardware Security Module)** — PKCS#11 / KMS. The private key **never
  leaves** the HSM; the front-end sends the to-be-signed handshake hash and gets
  a signature back. Used for compliance (FIPS 140-2/3), CA keys, and
  "keyless SSL" (Cloudflare Keyless: the CDN terminates TLS but the customer
  keeps the private key on-prem in their HSM; the CDN calls back for each
  signing op). Adds latency per full handshake (network round trip to the HSM) →
  resumption matters even more.

**6. Session-cache / ticket-key distribution as offload infrastructure.** A
shared, fast STEK store + high resume rate is effectively an offload: it removes
asymmetric ops from most connections. Treat ticket-key rotation as a first-class
job (§38.7).

**Sizing rule of thumb (order-of-magnitude, one modern x86 core, ECDSA P-256,
TLS 1.3):** full handshakes ≈ 10k–40k/s/core software, ≈ 100k+/s/core with QAT;
bulk AES-GCM with AES-NI ≈ several GB/s/core, ≈ line-rate with NIC inline. So a
handshake-heavy edge (many short connections) is **asymmetric-bound** → resume +
QAT/ECDSA; a bulk-heavy edge (video, downloads) is **symmetric/CPU-copy-bound**
→ kTLS + NIC offload + `sendfile`.

**What can go wrong with offload:**
- kTLS + `sendfile` + a middlebox that mangles → harder to debug (packets
  encrypted before your `tcpdump` on TX; capture on RX or use keylog).
- NIC offload resync bugs under heavy loss/reorder (rare, firmware-dependent).
- QAT/HSM as a **single point of failure / latency spike**; capacity-plan and
  have software fallback enabled.
- FIPS builds disable ChaCha20/X25519/Ed25519 in some configs → unexpected
  cipher/curve negotiation, mobile clients falling back.
- Offloaded checksums already make local captures show "bad checksum"
  (Ch. 49); offloaded TLS makes local TX captures show ciphertext you can't
  match to app writes — plan capture points accordingly.

### 38.13 QUIC's use of TLS 1.3

QUIC (Ch. 39) does **not** run TLS-over-a-stream. It uses the **TLS 1.3 handshake
state machine** to authenticate and derive keys, but:
- Handshake messages travel inside QUIC **`CRYPTO` frames** in special
  Initial/Handshake packets — there is **no TLS record layer**.
- **`ClientHello` fits in the first flight**; QUIC's own transport parameters
  (flow control limits, max streams, idle timeout, `max_udp_payload_size`,
  original DCID, etc.) are carried in a **TLS extension**
  (`quic_transport_parameters`), so transport + crypto negotiation are one shot.
- Keys derived from the TLS key schedule protect QUIC **packets** (header
  protection + AEAD payload protection), not records. There are distinct key
  phases: **Initial** (keys derived from a well-known salt + the client's
  Destination Connection ID — so anyone can decrypt Initial packets; this exists
  only for version negotiation/PMTU/DoS handling), then **Handshake**, then
  **1-RTT (application)**. The **Key Phase** bit in the short header signals
  `KeyUpdate` ratchets.
- **1-RTT** to first byte (same as TLS 1.3 over TCP, but QUIC also carries the
  request, so no separate TCP handshake) and **0-RTT** works the same way as
  §38.6 (PSK from a prior `NewSessionTicket`, same replay caveats — QUIC servers
  additionally enforce the **3× amplification limit** until the client's address
  is validated).
- **ALPN is mandatory** in QUIC (`h3`, `doq`, etc.).
- **No `ChangeCipherSpec`**, no downgrade to a record layer, most of the
  transport header is encrypted → far less middlebox ossification than TLS/TCP.
- **Retry / stateless reset**: a QUIC server under load can send a `Retry` packet
  with a token (like TLS 1.3 HRR + cookie) to validate the client address before
  committing state.
- Encryption is **not optional** in QUIC — there is no "plaintext QUIC." The only
  cleartext is a handful of header bits and the Initial packets (which are
  "encrypted" with public keys purely for tamper-evidence/versioning).

Practical upshot: your TLS knowledge transfers directly (cipher suites,
certificates, SNI/ECH, resumption, 0-RTT risk, mTLS), but "TLS termination" for
HTTP/3 means a **QUIC-aware** terminator that routes on the **Connection ID**
(not the 5-tuple, which changes on migration) and speaks UDP. Load balancers need
explicit HTTP/3 support (Ch. 39, 40).

### 38.14 Interactions with TCP that bite

- **Round-trip stack-up.** Cold HTTPS = TCP handshake (1 RTT) + TLS handshake
  (1.3: 1 RTT; 1.2: 2 RTT) before the first request byte → **2–3 RTT to first
  byte**. Resumption (1.2 abbreviated / 1.3 PSK) → 1 RTT. 0-RTT → 0. TCP Fast
  Open can save the TCP RTT but has poor middlebox support. This is the whole
  reason CDNs/edge termination exist — make those RTTs short (Ch. 4, 50).
- **Head-of-line blocking is TCP's, and TLS inherits it.** TLS records must be
  processed in order; one lost TCP segment blocks decryption of every record
  behind it, which blocks every HTTP/2 stream on that connection. QUIC fixes this
  by moving crypto below the streams (Ch. 39).
- **Record size vs a lost segment.** A 16 KB record split across ~11 segments is
  undecryptable until the *last* byte arrives; lose one segment and the whole
  record waits for the retransmit. Dynamic record sizing (§38.9) mitigates.
- **MTU / large ClientHello.** Many extensions + **hybrid post-quantum key
  shares** (X25519MLKEM768 adds ~1.2 KB) push ClientHello past one segment
  (>1460 B). Middleboxes/LBs that assume "ClientHello fits in one packet" (for
  SNI inspection) drop or mishandle it → handshake hangs from *some* networks
  only. Ensure multi-segment ClientHello reassembly; watch this as PQC rolls out.
- **PMTU black hole** (Ch. 8) hits TLS exactly like any bulk TCP flow: handshake
  (small packets) succeeds, first full-size `application_data` record stalls
  forever. Signature is identical; fix is the same (allow ICMP type 3 / MSS
  clamp).
- **`close_notify` and truncation.** A clean TLS shutdown sends `close_notify`
  before the TCP FIN. If the peer just sends FIN/RST, a strict client treats the
  response as **truncated** (possible truncation attack); a lax one accepts it.
  Mismatched strictness between client libs and a load balancer that RSTs idle
  connections → spurious "SSL truncated" / `ECONNRESET` errors at low rates.
- **Renegotiation vs middleboxes** (1.2): some proxies choke on mid-stream
  handshakes; 1.3 removed renegotiation entirely (use `KeyUpdate` /
  post-handshake auth).
- **TLS version/cipher downgrade.** Active attackers (or dumb middleboxes) strip
  extensions or force old versions. 1.3 bakes a **downgrade sentinel** into
  `server_random` and covers everything with the transcript hash + `Finished`, so
  downgrades are detected. Still: disable ≤1.1 so there's nothing to downgrade
  *to*.

### 38.15 Debugging TLS + war stories

**Try it yourself**
```bash
# Full handshake detail, chain, negotiated version/cipher/curve, resumption:
openssl s_client -connect example.com:443 -servername example.com -showcerts </dev/null
openssl s_client -connect example.com:443 -servername example.com -tls1_3 -brief </dev/null
openssl s_client -connect example.com:443 -servername example.com < /dev/null 2>&1 \
  | grep -E 'Protocol|Cipher|Server Temp Key|Verify|Session-ID|TLS session ticket|Early data'

# Test a specific version / force downgrade behavior:
openssl s_client -connect host:443 -tls1_2 ; openssl s_client -connect host:443 -tls1_1

# ALPN result (h2 / http/1.1):
openssl s_client -alpn h2,http/1.1 -connect example.com:443 -servername example.com </dev/null 2>&1 | grep ALPN

# OCSP staple present & valid?
openssl s_client -connect example.com:443 -servername example.com -status </dev/null 2>&1 \
  | grep -A17 'OCSP response'

# Resumption actually working? (second call should say "Reused")
( echo; sleep 1 ) | openssl s_client -connect host:443 -servername host -sess_out /tmp/s.pem >/dev/null 2>&1
( echo; sleep 1 ) | openssl s_client -connect host:443 -servername host -sess_in  /tmp/s.pem 2>&1 | grep -E 'Reused|New'

# Timing breakdown (TCP vs TLS vs first byte):
curl -w 'tcp:%{time_connect}  tls:%{time_appconnect}  ttfb:%{time_starttransfer}  total:%{time_total}\n' \
  -o /dev/null -s https://example.com

# Decrypt your own traffic in Wireshark (no private key needed for ECDHE/1.3):
SSLKEYLOGFILE=/tmp/keys.log curl https://example.com     # or set it for Chrome/Firefox
#   Wireshark > Preferences > Protocols > TLS > (Pre)-Master-Secret log filename = /tmp/keys.log

# HTTP/3 / QUIC:
curl -sI --http3 https://cloudflare.com | head -1
dig +short HTTPS example.com            # ech= and alpn=h3 hints

# Full posture audit:
testssl.sh https://example.com          # or the SSL Labs / Mozilla Observatory web tools
nmap --script ssl-enum-ciphers -p 443 example.com

# kTLS in use? (Linux)
grep -r . /proc/net/tls* 2>/dev/null ; ss -ti | grep -i tls ; lsmod | grep '^tls'
# NIC TLS offload capability:
ethtool -k eth0 | grep -i tls           # tls-hw-tx-offload / tls-hw-rx-offload
```

**War story 1 — the un-rotated ticket key.** A security review found that a
company's nginx fleet had `ssl_session_ticket_key` pointing at a file baked into
the base image **18 months earlier** and never rotated. Anyone who had ever pulled
that image (contractors, a former employee, a leaked registry) held a key that
could decrypt **every resumed session** across the fleet, retroactively, from
packet captures — the "forward secrecy" on the marketing page was fiction for
~60% of connections. Fix: a `systemd` timer generating a fresh key every 12 h,
distributed via the secrets store, nginx reloading with a 3-key ring
(`keyN`, `keyN-1`, `keyN-2`), ticket lifetime cut to 8 h. **Lesson:** session
tickets are only as forward-secret as your STEK rotation; treat ticket keys like
signing keys.

**War story 2 — Must-Staple + a flaky CA.** A team enabled OCSP **Must-Staple**
on their cert "for security." Months later the CA's OCSP responder had a 40-minute
outage; the servers' staple cache had already expired, so **every new TLS
handshake failed hard** and the site went fully down — a self-inflicted outage
worse than the revocation risk it defended against. Fix: removed Must-Staple,
kept stapling with a long cache + proactive refresh at 50% of validity + alerting
on staple age, and moved to 90-day auto-renewed certs. **Lesson:** Must-Staple
converts "CA availability" into "your availability"; only use it if you can
guarantee staple freshness.

**War story 3 — PQC ClientHello over a legacy LB.** After a client-library bump
enabled `X25519MLKEM768`, handshakes started **timing out from exactly one
partner's network**. Their aging L4 load balancer parsed SNI only from the first
TCP segment; the ~1.7 KB ClientHello spanned two segments and the LB dropped
connections whose ClientHello wasn't complete in packet one. Short-term: pin that
route to `X25519` only. Long-term: LB firmware update with multi-segment
ClientHello reassembly. **Lesson:** TLS lives *in* the TCP byte stream — anything
that assumes "the ClientHello is one packet" breaks as handshakes grow (PQC makes
this imminent for everyone).

**War story 4 — offload made it *slower*.** A video edge added QAT cards expecting
a throughput jump; p99 handshake latency got **worse** under load. QAT was
configured for RSA-2048 offload, but the certs were **ECDSA P-256** — already
cheap in software — so requests queued for a co-processor they didn't need while
AES-NI bulk crypto (the actual bottleneck at 40 GbE of video) got no help. Fix:
turn off QAT for the asymmetric path, enable **kTLS + `sendfile` + NIC inline
TLS** for the bulk path; CPU dropped ~35% and p99 normalized. **Lesson:** profile
first — handshake-bound and bulk-bound TLS want *different* offloads (§38.12).

---

## 39. HTTP/1.1 vs HTTP/2 vs HTTP/3 (QUIC)

### 39.1 The problem each version solves

| | HTTP/1.1 | HTTP/2 | HTTP/3 |
|---|---|---|---|
| Year | 1997 | 2015 (RFC 7540 / 9113) | 2022 (RFC 9114) |
| Transport | TCP (usually + TLS) | **one** TCP + TLS connection | **QUIC** over **UDP** (+ TLS 1.3 built in) |
| Concurrency | 1 request per connection at a time; browsers open ~6 per origin; **pipelining** exists but is broken by intermediaries | **streams** multiplexed on one connection; interleaved frames | streams multiplexed **without TCP HOL blocking** |
| Head-of-line blocking | **at HTTP layer** (request N blocks N+1 on that connection) | **at TCP layer** (one lost segment stalls *all* streams) | **eliminated** — each QUIC stream is independently reliable; a loss only stalls its own stream |
| Header compression | none (verbose, repeated) | **HPACK** | **QPACK** |
| Server push | no | yes (largely deprecated in practice) | yes (rare) |
| Connection setup | TCP + TLS = 2–3 RTT | same | **1-RTT**, or **0-RTT** resumption; combines transport+crypto handshake |
| Connection migration | no (new IP = new connection) | no | **yes** — a Connection ID survives IP/port change (Wi-Fi↔cellular) |
| Prioritization | via connection count | stream priority tree (H2) / scheme (RFC 9218) | RFC 9218 |

### 39.2 HTTP/1.1 mechanics you still deal with

- **Persistent connections** (`Connection: keep-alive`, default in 1.1). Reuse
  avoids repeated TCP+TLS handshakes — huge.
- **Chunked transfer encoding** (`Transfer-Encoding: chunked`) streams a body of
  unknown length; each chunk is `<hex size>\r\n<data>\r\n`, terminated by a
  `0\r\n\r\n`.
- **`Content-Length` vs chunked** mismatch, or two `Content-Length` headers, or
  `CL` + `TE` together → **request smuggling** when a front proxy and back server
  disagree on where the request ends (Ch. 52 relevance). Normalize at the edge.
- **Pipelining** (send request 2 before response 1 returns) is allowed but
  disabled in every major browser because buggy proxies mis-associate responses.
- **Connection limits:** ~6 per origin per browser → domain sharding was a common
  (now counterproductive with H2) hack.

### 39.3 HTTP/2 mechanics

- Everything is a **frame** (`HEADERS`, `DATA`, `SETTINGS`, `WINDOW_UPDATE`,
  `RST_STREAM`, `PING`, `GOAWAY`) tagged with a **stream ID** (odd = client-
  initiated). Frames from many streams interleave on the one TCP connection.
- **Flow control** exists *again* at the HTTP/2 layer (per-stream and
  per-connection `WINDOW_UPDATE`) — independent of TCP's window. A stuck reader
  on one stream can stall via HTTP/2 flow control even if TCP is fine.
- **HPACK** keeps a synchronized header table both sides; repeated headers cost
  ~1 byte. Dynamic-table desync or a hostile peer → the **HPACK bomb** DoS.
- **The TCP HOL problem:** all streams ride one TCP connection, so a single lost
  TCP segment blocks delivery of *every* stream's bytes behind it. On a clean
  network H2 crushes H1; on a lossy network (mobile) H2 can be *worse* than H1's
  6 connections (which fail independently). This is the entire motivation for
  HTTP/3.
- **`GOAWAY`**: graceful connection shutdown telling the client the highest
  stream it will process — clients replay higher streams on a new connection.
  Load balancers use this for draining.

### 39.4 HTTP/3 / QUIC mechanics

QUIC (RFC 9000) is a full reliable, multiplexed, encrypted transport in
**userspace**, over **UDP**:

- **Streams** are independent: loss on stream 3 doesn't delay stream 7. Solves
  TCP HOL blocking.
- **Crypto is integral:** the transport and TLS 1.3 handshakes are one; 1-RTT
  setup, 0-RTT resumption.
- **Connection ID:** the connection is identified by a CID in the QUIC header,
  **not** the 4-tuple. Change networks (new IP/port) and the connection
  **migrates** seamlessly — great for phones. Load balancers must route by CID,
  not 5-tuple.
- **Always encrypted, including most of the transport header** — middleboxes
  can't inspect/mangle sequence numbers, so QUIC can *evolve*
  ("ossification" resistance). Only a few bits are visible.
- **Its own congestion control** (NewReno/CUBIC/BBR reimplemented in userspace),
  **ACK ranges** (like SACK but richer, no separate "cumulative" limitation),
  **PLPMTUD** (no ICMP dependency).
- **Costs:** UDP is more expensive per-packet in many kernels/NICs (less
  offload — improving with UDP GSO/GRO and hardware QUIC offload); some networks
  rate-limit or block UDP/443 → clients fall back to HTTP/2. QUIC amplification
  limits: a server may send at most ~3× the bytes it received from an unvalidated
  client address.

### 39.5 How a browser picks a version

1. Connects via TCP+TLS, ALPN offers `h2, http/1.1` → gets HTTP/2 or 1.1.
2. Server response includes `Alt-Svc: h3=":443"; ma=86400` **or** the DNS `HTTPS`
   record advertises `alpn=h3` → browser tries **QUIC/HTTP/3** on the next
   navigation and caches that.
3. If UDP/443 is blocked/broken → falls back to h2 (Happy-Eyeballs-style racing
   in modern browsers).

### 39.6 War story

A news site enabled HTTP/2 and saw p75 page-load *improve* but p95 (mobile users)
*regress* by 20%. On lossy cellular, all 30 sub-resources shared one TCP
connection, and each packet loss stalled the whole pipe (TCP HOL). Rolling out
**HTTP/3** (via `Alt-Svc`) recovered p95 because QUIC streams fail independently,
and connection migration cut reconnects during cell handovers. **Lesson:** HTTP/2
trades many-connection resilience for single-connection efficiency; on lossy
links that trade can lose. HTTP/3 is the fix.

**Try it yourself**
```bash
curl -sI --http2 https://example.com | grep -i '^HTTP\|alt-svc'
curl -sI --http3 https://cloudflare.com | head -1        # needs curl built with HTTP/3
dig example.com HTTPS +short                               # alpn=h3 advertised?
nghttp -nv https://example.com                             # HTTP/2 frame-level trace
# See which protocol Chrome used: DevTools > Network > Protocol column (h2 / h3 / http/1.1)
```

---

## 40. Load balancing: L4 vs L7, DSR, and connection handling

### 40.1 L4 (transport) load balancing

Operates on the **5-tuple**. Picks a backend per **connection** and forwards
packets; never parses the payload. Fast, protocol-agnostic (works for any TCP/UDP
service), preserves end-to-end TLS.

Forwarding modes:
| Mode | How | Return path | Notes |
|---|---|---|---|
| **NAT / proxy** | LB rewrites dst IP (+ maybe src) to the backend | back through the LB | LB sees both directions; backend sees LB's IP as client unless PROXY protocol / `X-Forwarded-For` |
| **DSR (Direct Server Return) / L3-DSR** | LB rewrites only dst MAC (same subnet) or tunnels (IPIP/GRE); backend has the VIP on loopback and replies **directly** to the client | bypasses the LB | LB only handles inbound (often 10% of bytes); backend must not ARP for the VIP; huge scale for video/download |
| **Maglev / ECMP + consistent hashing** | routers ECMP-spread flows to a fleet of stateless LB nodes; each hashes the 5-tuple consistently to a backend | via LB or DSR | Google Maglev, Cilium, Katran (XDP). Survives LB node churn without breaking flows. |

**Connection consistency** is the hard part: if the LB set changes (scale, deploy,
failure), a flow's packets must keep hitting the *same* backend. Solutions:
consistent hashing (minimal disruption), flow tables synced across LB nodes, or
"stickiness" via source-IP hashing.

### 40.2 L7 (application) load balancing / reverse proxy

**Terminates** the client TCP+TLS connection, parses HTTP, and opens its **own**
connection(s) to backends. Two independent TCP connections, two congestion-control
loops.

Capabilities L4 can't do:
- Route by **URL path**, `Host`, header, cookie, method.
- **Connection multiplexing / pooling:** thousands of client connections → a
  small pool of backend keep-alive connections (protects backends from the
  ephemeral-port / handshake storm — Ch. 21 war story).
- **Retries, timeouts, circuit breaking, outlier ejection** per request.
- **Protocol translation:** HTTP/3 client ↔ HTTP/1.1 backend; gRPC-web ↔ gRPC.
- **Buffering** slow clients (absorb a slow-loris client so the backend isn't
  tied up — Ch. 52).
- **Observability:** per-request logs, traces, status codes.
- Inject `X-Forwarded-For` / `Forwarded` / **PROXY protocol** so backends learn
  the real client IP.

Cost: it's in the data path for every byte, needs CPU for TLS, and adds a hop of
latency + its own buffering/HOL characteristics. Examples: nginx, HAProxy, Envoy,
Traefik, ALB, Cloudflare.

### 40.3 Health checks

- **L4 check:** can I `connect()` to `backend:port`? (TCP SYN/ACK). Cheap, but
  "port open" ≠ "app healthy."
- **L7 check:** `GET /healthz` → expect `200` and maybe a body. Reflects real
  readiness (DB reachable, caches warm).
- **Passive / outlier detection:** eject a backend that returns 5xx / connection
  errors for real traffic, probe it, re-admit.
- **Failure math:** check interval × unhealthy-threshold = detection time. Too
  aggressive → flapping under load; too slow → users hit dead backends. Pair with
  fast client retries to a *different* backend.

### 40.4 Session persistence ("sticky sessions")

Route a given user to the same backend (for in-memory session state):
- **Source-IP hash** (L4) — breaks with NAT/CGNAT (many users → one backend) and
  mobile IP changes.
- **Cookie-based** (L7) — LB sets/reads a cookie naming the backend. Robust.
- Better: **externalize session state** (Redis, signed JWT) so any backend can
  serve any request and you don't need stickiness at all.

### 40.5 Draining and deploys

To remove a backend without dropping requests: mark it **draining** → LB stops
sending *new* connections/requests, existing ones finish (bounded by a drain
timeout), backend sends HTTP/2 `GOAWAY` / closes keep-alives gracefully, then it's
pulled. Skipping drain = a burst of `502`/`connection reset` on every deploy.

### 40.6 War story

A service behind an L7 proxy saw periodic `502 Bad Gateway` spikes exactly at
deploy time. The proxy kept a pool of keep-alive connections to backend pods;
when Kubernetes killed a pod, the proxy sent a request onto an already-half-closed
connection and got a RST → `502`. Fixes: (1) backend sends `Connection: close` /
HTTP/2 `GOAWAY` and a `preStop` sleep so the proxy observes readiness=false and
drains first; (2) proxy configured to **retry** idempotent requests once on a
connection error to another backend. 502s went to zero. **Lesson:** L7 LBs pool
backend connections; backends must announce shutdown (GOAWAY/close) and LBs must
retry connection-level errors, or every deploy sheds requests.

**Try it yourself**
```bash
# Which backend am I hitting? (if the app echoes hostname)
for i in $(seq 10); do curl -s https://example.com/whoami; done
# See the LB's added headers:
curl -s https://example.com/headers | grep -iE 'x-forwarded|forwarded|via|x-real-ip'
# L4 vs L7 clue: does end-to-end TLS work (L4) or does the cert belong to the LB (L7)?
openssl s_client -connect example.com:443 -servername example.com </dev/null 2>&1 | grep -E 'issuer|subject'
# HAProxy/Envoy admin stats if exposed:
curl -s http://lb:8404/stats\;csv   # HAProxy
curl -s http://lb:9901/clusters      # Envoy
```

---
# Part XII — Observability and capture

## 41. `tcpdump` and BPF filters

### 41.1 Anatomy of a capture command

```
sudo tcpdump -i eth0 -nn -vv -e -s0 -c 200 -w /tmp/cap.pcap 'tcp port 443 and host 93.184.216.34'
            +iface+  |   |   |  |   |   |   +-- output file       +------ BPF filter ------+
                     |   |   |  |   |   +-- stop after N packets
                     |   |   |  |   +-- full packet (default is full on modern tcpdump)
                     |   |   |  +-- print link-layer (MAC) header
                     |   |   +-- extra verbose (TTL, IP id, options)
                     |   +-- don't resolve hostnames (-n) or ports (-nn); avoids DNS stalls
```

Useful flags: `-A` (payload as ASCII), `-X` (hex+ASCII), `-tttt` (readable
timestamps), `-Q in|out` (direction), `--time-stamp-precision=nano`, `-G 60 -w
cap-%H%M%S.pcap` (rotate hourly), `-C 100` (rotate at 100 MB), `-Z user` (drop
privileges).

**Golden rules**
- Always `-w` to a file for anything you'll analyze; read later with `-r` or in
  Wireshark. Printing to the terminal drops packets under load.
- `-s0` + high rate can overflow the kernel capture buffer → "N packets dropped
  by kernel". Raise `-B 4096` (KiB) or snap shorter (`-s 128` for headers only).
- Filter as tightly as possible at capture time; CPU and disk are finite.
- Capture on **both ends** for anything involving loss/latency — you can only
  prove "sent here, not received there" with two captures.

### 41.2 BPF filter cookbook

```
host 10.0.0.5                       # to or from
src host / dst host 10.0.0.5
net 10.0.0.0/24
port 443            portrange 8000-8100
tcp port 443 and not host 10.1.1.1

# TCP flags (tcp[tcpflags] is a bitmask):
'tcp[tcpflags] & (tcp-syn) != 0'                 # any SYN (incl SYN-ACK)
'tcp[tcpflags] & (tcp-syn|tcp-ack) == tcp-syn'   # SYN only (connection attempts)
'tcp[tcpflags] & (tcp-rst) != 0'                 # resets
'tcp[tcpflags] & (tcp-fin) != 0'                 # FINs
'(tcp[tcpflags] & (tcp-syn|tcp-fin|tcp-rst)) != 0'   # setup/teardown only

# Payload / size:
'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420'  # TCP payload starts with "GET "
'ip[2:2] > 1400'                                  # IP total length > 1400
'greater 1400'
'ip[6:2] & 0x3fff != 0'                           # IPv4 fragments (MF set or offset != 0)
'tcp[13] & 32 != 0'                               # URG flag set (weird; middlebox check)

# Other:
'icmp or icmp6'
'arp'
'vlan 100'
'udp port 53'
'ip proto 47'                                     # GRE
'ether host aa:bb:cc:dd:ee:ff'
```

`tcpdump` filters are compiled BPF and run **in the kernel** — very cheap. Display
filters (Wireshark's `tcp.analysis.retransmission` etc.) run in userspace on the
already-captured data and are far more expressive; use capture filters to bound
volume, display filters to investigate.

### 41.3 Reading `tcpdump` TCP output

```
12:00:00.100000 IP 10.0.0.5.51000 > 93.184.216.34.443: Flags [S], seq 1000,
                win 64240, options [mss 1460,sackOK,TS val 1 ecr 0,nop,wscale 7], length 0
12:00:00.190000 IP 93.184.216.34.443 > 10.0.0.5.51000: Flags [S.], seq 500, ack 1001,
                win 65535, options [mss 1460,sackOK,TS val 9 ecr 1,nop,wscale 9], length 0
12:00:00.190100 IP 10.0.0.5.51000 > 93.184.216.34.443: Flags [.], ack 501, win 502, length 0
12:00:00.190200 IP 10.0.0.5.51000 > 93.184.216.34.443: Flags [P.], seq 1001:1518, ack 501,
                win 502, length 517           # TLS ClientHello (517 bytes)
```

- `Flags [S]` SYN, `[S.]` SYN-ACK, `[.]` bare ACK, `[P.]` PSH+ACK (data),
  `[F.]` FIN, `[R]` RST.
- `seq 1001:1518` = bytes 1001..1517 (`length 517`). Relative numbers by default.
- `win 502` after the handshake = **scaled** (502 × 2^7 ≈ 64 KB).
- `ack 501` = "I have everything through 500."
- `options [nop,nop,sack 1 {3001:4001}]` on a later ACK = a SACK block.
- `ecr` echoes the peer's `val` — that pair is your RTT sample source.

### 41.4 War story

"Packet loss between app and DB" — a single capture on the app server showed
retransmissions. Capturing **simultaneously on the DB** showed the "lost"
segments *did* arrive, were ACKed, but the ACKs never reached the app. The loss
was **one-directional** on the return path through a flaky LACP member. A single
capture would have blamed the wrong direction. **Lesson:** for loss/latency,
capture both ends with synchronized clocks (`-tttt`, or NTP/PTP), then compare
by sequence number.

**Try it yourself**
```bash
sudo tcpdump -i any -nn -c20 'tcp[tcpflags] & (tcp-syn|tcp-rst) != 0'
sudo tcpdump -i eth0 -nn -s128 -w /tmp/db.pcap -B 8192 'host 10.2.0.9 and tcp port 5432'
sudo tcpdump -r /tmp/db.pcap -nn -ttt 'tcp[tcpflags] & tcp-rst != 0'
# ring buffer for an intermittent issue (24 files x 100MB):
sudo tcpdump -i eth0 -nn -s0 -C 100 -W 24 -w /var/tmp/ring.pcap 'port 443'
```

---

## 42. Reading a Wireshark trace like a detective

### 42.1 First moves on any capture

1. **Statistics → Capture File Properties**: duration, packet rate, drops.
2. **Statistics → Conversations** (TCP tab): sort by bytes / by "Bits/s". Find
   the flow that matters. "Follow → TCP Stream" to see the app dialogue.
3. **Statistics → Protocol Hierarchy**: what's actually on this wire.
4. **Set the time display** to "Seconds Since Previous Displayed Packet"
   (View → Time Display Format) — gaps jump out.
5. Enable **`tcp.analysis`** expert flags (on by default): look at
   **Analyze → Expert Information**.

### 42.2 The display filters that find problems

| Filter | Finds |
|---|---|
| `tcp.analysis.retransmission` | retransmitted segments |
| `tcp.analysis.fast_retransmission` | fast retransmits (3 dup-ACKs) |
| `tcp.analysis.spurious_retransmission` | retransmit of data already ACKed (RTO too tight / reordering) |
| `tcp.analysis.out_of_order` | reordering |
| `tcp.analysis.duplicate_ack` | dup-ACKs (`tcp.analysis.duplicate_ack_num >= 3` = triggering) |
| `tcp.analysis.lost_segment` | a gap Wireshark inferred ("previous segment not captured") |
| `tcp.analysis.zero_window` | receiver said "stop, buffer full" |
| `tcp.analysis.window_full` | sender has filled the offered window |
| `tcp.analysis.ack_rtt > 0.2` | ACKs that took >200 ms (delayed-ACK or path) |
| `tcp.flags.reset == 1` | who reset, and when relative to the last data |
| `tcp.time_delta > 0.2 && tcp.flags.push == 1` | app "think time" gaps |
| `tcp.analysis.bytes_in_flight` | plot it (I/O Graph) to see cwnd behavior |
| `tcp.stream == N` | isolate one connection |

### 42.3 Signatures → diagnosis

| What you see | Likely cause |
|---|---|
| Steady retransmissions + `duplicate_ack` bursts, throughput sawtooths | real packet loss on the path; check both-end captures & link counters |
| `spurious_retransmission` + `out_of_order`, ~no true loss | packet reordering (ECMP/LACP hashing), RTO too aggressive |
| `zero_window` from the receiver, then `window_update` seconds later | **receiver app too slow** — not a network problem |
| Long `tcp.time_delta` gap *before* a PSH from the server, no missing packets | **server-side application latency** (slow query, GC pause) |
| Client SYN, SYN, SYN at 1s/3s, no SYN-ACK | SYN or SYN-ACK lost; firewall DROP; accept-queue overflow |
| Handshake OK, first full-size (~1500) data packet retransmitted forever | **PMTU black hole** (ICMP filtered) — Ch. 8 |
| `RST` immediately after `SYN-ACK` | backlog overflow w/ abort, or app crashed on accept |
| `RST` right after a normal data exchange, from a *third* pattern/TTL | forged RST by a middlebox/IPS |
| ACKs every ~40 ms, tiny segments, request in 2 parts | **Nagle × delayed-ACK** — Ch. 33 |
| `win` never grows past ~64 KB on a high-RTT path | **window scaling** stripped/off — Ch. 30 |
| Big `bytes_in_flight`, RTT climbing steadily with throughput flat | **bufferbloat** at the bottleneck — Ch. 44 |
| Duplicate IP ID / duplicate SYN with different TTL | two hosts sharing an IP, or a retransmit vs a spoofed packet |

### 42.4 The I/O Graph (your most powerful view)

`Statistics → I/O Graph`. Plot on the same time axis:
- `tcp.analysis.bytes_in_flight` (AVG) — proxy for `cwnd`.
- `tcp.analysis.retransmission` (COUNT) — loss events.
- `tcp.analysis.ack_rtt` (AVG/MAX) — RTT over time.
- Throughput (bytes, "Bits/s").

Now you can *see* slow start ramp, the sawtooth of CUBIC on loss, an RTO
(everything stops for ~200 ms–1 s then restarts from near zero), or RTT rising
with in-flight bytes (bufferbloat).

### 42.5 tshark for scripted analysis

```bash
# Per-connection retransmit counts:
tshark -r cap.pcap -q -z conv,tcp
# List every retransmission with time and stream:
tshark -r cap.pcap -Y tcp.analysis.retransmission \
  -T fields -e frame.time_relative -e tcp.stream -e ip.src -e ip.dst -e tcp.seq
# RTT distribution:
tshark -r cap.pcap -Y tcp.analysis.ack_rtt -T fields -e tcp.analysis.ack_rtt | \
  sort -n | awk '{a[NR]=$1} END{print "p50",a[int(NR*.5)],"p95",a[int(NR*.95)],"max",a[NR]}'
# Time-to-first-byte per HTTP request:
tshark -r cap.pcap -Y http.response -T fields -e http.time -e http.request.uri
```

### 42.6 War story

A dashboard showed 2% "network errors" calling an internal API. The capture had
**zero** retransmissions or resets. Filtering `http.response.code >= 500` and
"Follow TCP Stream" showed the server returned `503` with a JSON body in ~4 ms —
the API's own rate limiter. The "network" error was an application status code
the client's metrics mislabeled. **Lesson:** confirm at the packet layer whether
a failure is TCP (RST/timeout/retransmit) or application (a clean HTTP 4xx/5xx)
before opening a network ticket.

**Try it yourself**
```bash
tshark -r cap.pcap -q -z io,stat,1,'COUNT(tcp.analysis.retransmission)tcp.analysis.retransmission','AVG(tcp.analysis.bytes_in_flight)tcp.analysis.bytes_in_flight'
tshark -r cap.pcap -q -z expert
wireshark cap.pcap   # then: Analyze > Expert Information; Statistics > I/O Graph
```

---

## 43. Socket statistics: `ss`, `/proc`, and connection states

### 43.1 `ss` — the modern `netstat`

```bash
ss -tan                      # all TCP, numeric
ss -tanp                     # + owning process (needs root for others' sockets)
ss -tl                       # listening TCP; Recv-Q = accept backlog used, Send-Q = its max
ss -tan state established
ss -tan state time-wait
ss -tan '( sport = :443 )'
ss -tan dst 10.0.0.0/8
ss -s                        # summary counts by state / mem
ss -tie                      # -e extended + -i internal TCP info  <-- the good one
ss -tim                      # + memory (rmem/wmem/... per socket)
ss -tio                      # + timer info (retransmit/keepalive countdown)
ss -natp | grep -c ESTAB
```

### 43.2 Decoding `ss -tie` (the per-socket TCP report)

```
ESTAB 0 0 10.0.0.5:443 203.0.113.9:54321
     ts sack cubic wscale:7,9 rto:212 rtt:11.5/5.2 ato:40 mss:1448 pmtu:1500
     rcvmss:536 advmss:1448 cwnd:42 ssthresh:35 bytes_sent:1048576 bytes_acked:1047000
     bytes_received:900 segs_out:730 segs_in:20 data_segs_out:729
     send 42.3Mbps lastsnd:4 lastrcv:1200 lastack:4 pacing_rate 84.6Mbps
     delivery_rate 40.1Mbps busy:980ms retrans:0/3 dsack_dups:1 reordering:5
     rcv_rtt:12 rcv_space:65535 minrtt:10.9
```

| Field | Meaning / use |
|---|---|
| `ts sack cubic` | options negotiated + congestion control in use |
| `wscale:7,9` | send/recv window scale shifts (both nonzero = scaling active) |
| `rto:212` | current retransmission timeout (ms). Near `TCP_RTO_MIN` (200) = healthy low-RTT path |
| `rtt:11.5/5.2` | **SRTT / RTTVAR** (ms). Rising RTTVAR = jittery path |
| `minrtt` | best RTT ever seen — the path's true propagation delay; `rtt` ≫ `minrtt` = queueing/bufferbloat |
| `ato:40` | delayed-ACK timeout in effect (ms) |
| `mss` / `pmtu` / `advmss` | segment sizing; `pmtu < 1500` = a smaller link on the path |
| `cwnd` / `ssthresh` | congestion window / threshold, in **segments**. `cwnd` small + `retrans` climbing = lossy |
| `retrans:0/3` | current unrecovered / **total** retransmits this connection |
| `send X`, `delivery_rate Y` | sender's rate estimate vs actually delivered rate |
| `pacing_rate` | rate the qdisc/stack is pacing at (BBR/fq) |
| `bytes_acked` vs `bytes_sent` | gap = data in flight or lost |
| `rcv_space` | receiver window autotuning target; grows on fast paths |
| `busy:` | time the socket had data to send (not app-limited) |
| `app_limited` (flag) | throughput is bounded by the app not writing, **not** the network |

### 43.3 `/proc` and `nstat`

```bash
nstat -az                       # every TCP/IP/UDP/ICMP kernel counter, since boot, non-zero-friendly
nstat                           # deltas since last run (great for before/after)
cat /proc/net/snmp              # classic SNMP-style counters (Tcp:, Udp:, Ip:, Icmp:)
cat /proc/net/netstat           # TcpExt: (the interesting ones)
cat /proc/net/sockstat          # sockets in use, TCP mem pages vs tcp_mem
cat /proc/net/tcp /proc/net/tcp6  # raw socket table (hex); ss parses this
```

Counters worth alerting on:
| Counter (`nstat`/`netstat -s`) | Meaning |
|---|---|
| `TcpExtListenOverflows` / `ListenDrops` | accept-queue full → dropped connections (Ch. 24) |
| `TcpExtSyncookiesSent` | under SYN-queue pressure / SYN flood |
| `TcpExtTCPBacklogDrop` | packets dropped because the socket backlog was full |
| `TcpExtPruneCalled` / `RcvPruned` / `OfoPruned` | out-of-order/receive memory pressure → data discarded |
| `TcpRetransSegs` / `TcpExtTCPLostRetransmit` | retransmissions / retransmits themselves lost |
| `TcpExtTCPTimeouts` / `TCPSpuriousRTOs` | RTO fired / RTO was a mistake |
| `TcpExtTCPFastRetrans` / `TCPSACKRecovery` | fast-path recovery events (good — better than timeouts) |
| `UdpInErrors` / `UdpRcvbufErrors` / `UdpInDatagrams` | UDP receive drops (Ch. 22) |
| `IpReasmFails` / `IpReasmReqds` | IP fragment reassembly failures (Ch. 15) |
| `TcpExtTW` / `TCPTimeWaitOverflow` | TIME_WAIT churn / hitting `tcp_max_tw_buckets` |
| `TcpExtTCPAbortOnMemory` / `TCPAbortOnTimeout` | connections killed by memory limits / retransmit exhaustion |
| `nf_conntrack: table full` (dmesg) | firewall state table exhausted → new flows dropped |

### 43.4 War story

A service's error rate rose with traffic but CPU, memory, and app latency looked
fine. `nstat` before/after a load test showed `TcpExtListenOverflows +14000` and
`TcpExtTCPReqQFullDrop`. The app's framework defaulted `listen()` backlog to 128;
bursts overflowed the accept queue and clients saw connect timeouts that the app
never logged (they never became connections). Raising backlog to `somaxconn`
(4096) + adding `SO_REUSEPORT` workers cleared it. **Lesson:** connection-level
failures are invisible to app logs; watch `nstat`'s `Listen*` / `*Drop` counters
as first-class SLIs.

**Try it yourself**
```bash
watch -n1 'ss -s'
ss -tie state established | grep -E 'minrtt|retrans|app_limited' | head
nstat -n && sleep 10 && nstat            # 10-second delta of every counter
nstat -az | grep -E 'ListenOverflows|ListenDrops|TCPBacklogDrop|SyncookiesSent|PruneCalled'
cat /proc/net/sockstat                    # 'mem' pages vs `sysctl net.ipv4.tcp_mem`
```

---


# Part XIII — Performance forensics

## 44. Diagnosing latency, loss, throughput, and bufferbloat

### 44.1 The four questions, in order

When "the network is slow," decide which of these it is — the fixes are totally
different:

1. **Latency** (RTT too high or too variable) — pages feel sluggish, RPC p99 bad,
   throughput of short transfers poor.
2. **Loss** (packets dropped) — throughput of *bulk* transfers poor, retransmit
   counters climbing, video pixelation.
3. **Throughput ceiling** with no loss — bulk transfer plateaus well below link
   rate; usually window/BDP/CC-limited (Ch. 35–36) or app-limited.
4. **Bufferbloat** — throughput is *fine* but latency balloons *while* a transfer
   runs, wrecking everything else on the link.

### 44.2 Step-by-step triage

```
A. Is it the network at all?
   - ss -tie: is the socket `app_limited`? busy: low? -> it's the app, not the net.
   - Capture: long gap before server PSH, no missing packets -> server think-time.
   - Compare curl time_starttransfer vs time_total -> TTFB (server) vs body (transfer).

B. Measure baseline latency (idle):
   - ping / mtr to the ENDPOINT (not a router). Note min / avg / mdev.
   - `ss -tie` minrtt vs rtt: rtt >> minrtt means standing queue somewhere.

C. Measure latency UNDER LOAD (the bufferbloat test):
   - Start `mtr <peer>` or `ping <peer>` in one window.
   - In another, run a big transfer both directions:
       iperf3 -c <peer> -t30            (upload)
       iperf3 -c <peer> -t30 -R         (download)
   - If ping RTT jumps 5-50x during the transfer -> BUFFERBLOAT at the bottleneck.
   - Which direction spikes tells you which buffer (up = your uplink/modem).

D. Quantify loss:
   - mtr -rwzbc300 : loss% that PERSISTS to the final hop is real.
   - nstat delta:  TcpRetransSegs / TcpOutSegs  ~= loss rate (sender side).
   - Two-ended capture: segments present at sender, absent at receiver.

E. Throughput ceiling with no loss:
   - 1 stream vs `iperf3 -P8`: if 8x faster -> window/CC/app-limited, not link.
   - ss -tie: cwnd small & flat, app_limited absent, retrans 0 -> rwnd/buffer cap.
   - Check window scaling on both SYNs; raise tcp_rmem/wmem max to 2*BDP.
```

### 44.3 Bufferbloat, in depth

**Cause:** a device at the bottleneck (home router, DSL/cable modem, LTE modem,
VM host vNIC, a `pfifo` qdisc with `txqueuelen 1000`) has a large **dumb FIFO**.
A loss-based sender (CUBIC) keeps increasing `cwnd` until that buffer overflows;
in steady state it stays nearly full. A full 1 MB buffer draining at 8 Mbit/s =
**1000 ms** of added one-way delay for *every* packet, including your DNS, your
SSH keystrokes, your video call, and the ACKs of the very transfer causing it.

**Diagnosis:** RTT idle 15 ms → RTT during upload 600 ms. `ss -tie` shows
`rtt:600 minrtt:15`. `tc -s qdisc` shows a `pfifo_fast`/`fifo` with a big
backlog and no drops until overflow.

**Fixes (in order of leverage):**
1. **AQM at the bottleneck.** On a Linux gateway:
   `tc qdisc replace dev <wan> root cake bandwidth 19Mbit` (set ~90–95% of true
   line rate so the queue forms in CAKE, which manages it, not in the modem).
   Or `fq_codel` if `cake` unavailable. Consumer routers: enable "Smart Queue
   Management / SQM" (OpenWrt `sqm-scripts`).
2. **`fq_codel` as the host default qdisc** (`net.core.default_qdisc=fq_codel`) —
   helps a host that is itself the bottleneck, and isolates flows (a bulk flow
   can't delay a sparse flow's packets).
3. **BBR** on your senders — paces to estimated bandwidth, keeps standing queue
   near zero.
4. **Reduce `txqueuelen`** on the egress interface (e.g. `ip link set eth0
   txqueuelen 100`) and NIC ring buffers if absurdly large.
5. **ECN** end-to-end so the AQM can *mark* instead of *drop*.

### 44.4 Latency decomposition

`Total one-way delay = propagation + serialization + queueing + processing`

| Component | Formula / cause | Fixable by |
|---|---|---|
| **Propagation** | distance ÷ (~200 000 km/s in fiber) — ~5 µs/km; NYC↔LON ≈ 28 ms one way, physics | move closer (CDN/edge), fewer RTTs (TLS resumption, 0-RTT, keep-alive) |
| **Serialization** | bytes ÷ link_rate — a 1500 B packet on 1 Mbit/s = 12 ms; on 1 Gbit/s = 12 µs | faster links; smaller packets on slow links |
| **Queueing** | time behind other packets in buffers — the variable, controllable part | AQM, less cross-traffic, QoS/priority |
| **Processing** | routing lookup, NAT, firewall, crypto, VM exit | hardware offload, fewer hops/middleboxes |

### 44.5 App-limited vs network-limited (don't fix the wrong thing)

`ss -tie` reports `app_limited` when the socket *could* have sent more but the
application didn't `write()` fast enough. If you see it, no amount of buffer/CC
tuning helps — the bottleneck is your code (single-threaded producer, small
writes, blocking disk read, GC pauses, `Content-Length` computed by buffering the
whole body first). `busy:` time near the total connection time = network-limited;
much lower = app-limited.

### 44.6 War story

A CI system's artifact downloads were "slow" (3 min for 500 MB, ~22 Mbit/s) on a
"100 Mbit/s" office link. Idle ping to the artifact host: 12 ms. Ping *during*
download: **340 ms**, and every other user's web browsing crawled. Classic
bufferbloat in the office's cable modem. The office gateway (a Linux box) got
`tc qdisc replace dev wan0 root cake bandwidth 95Mbit`; download rose to
~90 Mbit/s **and** in-download ping stayed at ~20 ms. **Lesson:** "slow
downloads + everyone else's internet dies during them" = bufferbloat; the cure is
AQM/shaping at the bottleneck, not a bigger pipe.

**Try it yourself**
```bash
# The two-terminal bufferbloat test:
#   term1:
mtr -o "L S  A M" <peer>
#   term2:
iperf3 -c <peer> -t30 ; iperf3 -c <peer> -t30 -R
# Or a one-shot score:
# (browser) https://www.waveform.com/tools/bufferbloat  or  fast.com (shows loaded latency)
ss -tie dst <peer> | grep -oE 'minrtt:[0-9.]+|rtt:[0-9.]+/[0-9.]+|app_limited|busy:[0-9]+ms'
tc -s qdisc show dev eth0
```

---

## 45. Emulating bad networks with `tc`/`netem`

`tc` (traffic control) with the **`netem`** qdisc lets you inject latency, jitter,
loss, reordering, duplication, and corruption on an interface — essential for
testing how your app behaves on real-world links **before** users find out.

### 45.1 Basic recipes

```bash
IFACE=eth0    # apply to the EGRESS of this interface (affects packets leaving)

# 100 ms fixed latency:
sudo tc qdisc add dev $IFACE root netem delay 100ms

# 100 ms +/- 20 ms jitter, correlation 25% (successive delays are related):
sudo tc qdisc change dev $IFACE root netem delay 100ms 20ms 25%

# 1% random packet loss:
sudo tc qdisc change dev $IFACE root netem loss 1%

# Correlated (bursty) loss - Gilbert-Elliott-ish:
sudo tc qdisc change dev $IFACE root netem loss 1% 25%

# Reordering: 25% of packets sent 10ms early (needs a delay baseline):
sudo tc qdisc change dev $IFACE root netem delay 10ms reorder 25% 50%

# Duplication / corruption:
sudo tc qdisc change dev $IFACE root netem duplicate 1%
sudo tc qdisc change dev $IFACE root netem corrupt 0.1%

# Rate limit to 1 Mbit with a small buffer (emulate a slow uplink):
sudo tc qdisc change dev $IFACE root netem rate 1mbit limit 20

# Combine: "bad hotel wifi" = 80ms+/-40ms jitter, 2% bursty loss, 5Mbit:
sudo tc qdisc replace dev $IFACE root netem delay 80ms 40ms rate 5mbit loss 2% 25%

# Remove everything:
sudo tc qdisc del dev $IFACE root
```

### 45.2 Applying it to only some traffic

`netem` on `root` hits everything (including your SSH session — do this on a test
box or over a console). To scope it, use a classful qdisc + filter:

```bash
sudo tc qdisc add dev $IFACE root handle 1: prio
sudo tc qdisc add dev $IFACE parent 1:3 handle 30: netem delay 200ms loss 3%
# send only traffic to 10.0.0.9:5432 through the netem band:
sudo tc filter add dev $IFACE parent 1:0 protocol ip u32 \
     match ip dst 10.0.0.9/32 match ip dport 5432 0xffff flowid 1:3
```

Or, cleaner, use a **network namespace** / **veth pair** and shape the veth so
only the app-under-test is affected (Ch. 55).

### 45.3 Emulating asymmetric and both-direction impairment

`tc` only shapes **egress**. To impair the *return* path too, either run `netem`
on the peer's egress, or use an **IFB** (Intermediate Functional Block) device to
redirect ingress through a qdisc:

```bash
sudo modprobe ifb numifbs=1
sudo ip link set ifb0 up
sudo tc qdisc add dev $IFACE handle ffff: ingress
sudo tc filter add dev $IFACE parent ffff: protocol ip u32 match u32 0 0 \
     action mirred egress redirect dev ifb0
sudo tc qdisc add dev ifb0 root netem delay 50ms loss 1%
```

### 45.4 What to test with it

- **Timeouts & retries:** does the client back off, or hammer? Does it retry
  non-idempotent requests (bad)?
- **Connection pools:** do they detect and evict dead connections under loss?
- **Keepalive/heartbeat** intervals vs a 30 s "NAT" timeout.
- **HTTP/2 vs HTTP/3** under 2% loss (HOL blocking — Ch. 39).
- **TLS handshake** under high latency + loss (does it complete in bounded time?).
- **Congestion control choice:** `iperf3 -C cubic` vs `-C bbr` under
  `netem loss 1% delay 100ms`.
- **User-facing SLOs:** run the real UI through a "3G" profile.

### 45.5 Reference impairment profiles

| Profile | delay | jitter | loss | rate |
|---|---|---|---|---|
| Same DC | 0.2ms | 0.05ms | 0 | 10gbit |
| Same region | 2ms | 0.5ms | 0.001% | 1gbit |
| Cross-country | 40ms | 3ms | 0.01% | 200mbit |
| Transatlantic | 90ms | 5ms | 0.05% | 100mbit |
| Good 4G | 50ms | 15ms | 0.1% | 30mbit |
| Congested 4G / train | 120ms | 60ms | 1.5% (bursty 25%) | 3mbit |
| Hotel/café Wi-Fi | 80ms | 40ms | 2% (bursty) | 5mbit |
| GEO satellite | 600ms | 20ms | 0.5% | 20mbit |

### 45.6 War story

A mobile app's sync feature passed every test in the office and generated a storm
of duplicate writes in the field. Reproduced instantly with
`tc qdisc add dev eth0 root netem delay 300ms 100ms loss 3%`: the client's 2 s
request timeout fired while the server was still processing, the client retried a
**non-idempotent** `POST /transactions`, and both eventually succeeded.
Fixes: idempotency keys, exponential backoff with jitter, and a timeout ≥ the
server's p99 processing time + 2×RTT. **Lesson:** most "works in the office"
distributed bugs reproduce in one command with `netem`; add it to CI.

**Try it yourself**
```bash
sudo tc qdisc replace dev eth0 root netem delay 100ms 20ms loss 1%
ping -c10 <peer>            # see the injected delay/jitter/loss
curl -w 'ttfb %{time_starttransfer}\n' -o /dev/null -s https://<peer>/
sudo tc -s qdisc show dev eth0     # sent/dropped/reordered counters
sudo tc qdisc del dev eth0 root    # clean up!
```

---


# Part XIV — NAT traversal and real-time media

## 46. NAT behaviors, STUN, TURN, and ICE

### 46.1 The peer-to-peer problem

Two clients (a laptop on home Wi-Fi, a phone on cellular) both behind NAT want a
direct connection (video call, game, file transfer). Neither has a reachable
public `IP:port`. Neither can initiate to the other. You need **NAT traversal**.

### 46.2 NAT mapping and filtering behavior (RFC 4787 terms)

Two independent properties of a NAT, per internal socket:

**Mapping behavior** — does the same internal `IP:port` get the same external
`IP:port` for different destinations?
- **Endpoint-Independent Mapping (EIM):** yes — one external port regardless of
  destination. Traversal-friendly. (Most home routers today.)
- **Address-Dependent / Address-and-Port-Dependent Mapping:** a new external port
  per destination IP (or per IP+port). This is **"symmetric NAT"**. Hostile to
  traversal — the port the peer must target isn't the port anyone else observed.

**Filtering behavior** — who is allowed to send *in* to an existing mapping?
- **Endpoint-Independent Filtering:** anyone (once the mapping exists) — "full
  cone".
- **Address-Dependent Filtering:** only hosts the internal endpoint has sent to
  (any port).
- **Address-and-Port-Dependent Filtering:** only the exact `IP:port` the internal
  endpoint sent to — "port-restricted cone".

Also matters: **hairpinning** (can two internal hosts reach each other via the
external IP?) and **port preservation** / **port parity**.

### 46.3 STUN — discover your public mapping

**STUN (Session Traversal Utilities for NAT, RFC 8489)**: the client sends a
`Binding Request` to a public STUN server (UDP/3478, or over TLS/TCP). The server
replies with the **source `IP:port` it saw** — i.e. your external mapping for
that socket. Now the client knows its "server-reflexive" candidate.

- If your NAT has **EIM**, that reflexive `IP:port` is the same one your peer can
  send to → direct connection works after a **hole punch**.
- If your NAT is **symmetric**, the port STUN saw is *only* valid toward the STUN
  server; your peer will hit a different port → STUN alone fails → you need TURN.

STUN also does keepalives (refresh the mapping before it times out) and, in ICE,
connectivity checks.

### 46.4 Hole punching

```
A (behind NAT-A)                  Signaling server                 B (behind NAT-B)
  |-- STUN -> learns A_pub:pa --------|                              |
  |                                   |------ STUN -> B_pub:pb ------|
  |------ offer (A_pub:pa) ---------> relayed <---- answer (B_pub:pb)-|
  |                                                                   |
  |== A sends UDP to B_pub:pb (opens NAT-A mapping; NAT-B drops it) ==>|
  |<= B sends UDP to A_pub:pa (opens NAT-B mapping; now NAT-A allows) =|
  |== subsequent packets both ways pass - direct P2P path established ==|
```

Both sides must fire near-simultaneously (coordinated by the signaling channel) so
each NAT sees an *outbound* packet first and then permits the *inbound* reply.
Works when at least the mapping is endpoint-independent on both sides; often works
with one symmetric NAT via port prediction; fails with symmetric-on-both.

### 46.5 TURN — relay when direct fails

**TURN (Traversal Using Relays around NAT, RFC 8656)**: a public relay server
allocates a **relayed transport address** on itself; both peers send to the
relay, which forwards. Always works (it's just a well-connected middleman) but
costs bandwidth/latency and server capacity. ~10–20% of real-world WebRTC calls
end up relayed. TURN runs over UDP/3478, TCP/3478, or **TLS/5443/443** (to punch
through restrictive firewalls that only allow "HTTPS").

### 46.6 ICE — try everything, pick the best

**ICE (Interactive Connectivity Establishment, RFC 8445)** is the algorithm that
ties it together. Each peer gathers **candidates**:
- **host** — local `IP:port` (LAN, works if peers are on the same network / have
  IPv6).
- **server-reflexive (srflx)** — from STUN (your public mapping).
- **relayed** — from TURN.
- (**peer-reflexive** — discovered during checks.)

Peers exchange candidate lists via signaling (SDP), form **candidate pairs**,
prioritize them (host > srflx > relay; IPv6 often preferred), and run
**connectivity checks** (STUN binding requests) on every pair simultaneously.
The highest-priority pair that succeeds becomes the path; ICE can switch if it
degrades (**ICE restart**). **Trickle ICE** sends candidates as they're found
instead of waiting, cutting call setup time.

### 46.7 Why IPv6 quietly fixes most of this

With end-to-end IPv6, both peers have globally routable addresses; ICE's **host**
candidates connect directly (subject to each side's stateful firewall allowing
the hole-punched flow). No STUN mapping games, no TURN relay. This is a real,
measurable reason to deploy IPv6 for real-time apps.

### 46.8 Keepalives and timeouts

UDP NAT mappings expire fast — RFC 4787 says **≥ 2 min**, but real routers use
**30–60 s**, some **as low as 20 s**. Any P2P/media session must send a
keepalive (STUN binding indication, or an RTP no-op, or a DTLS heartbeat) every
**15–25 s** or inbound media stops when the mapping is reaped. TCP mappings live
longer (minutes to hours) but still need keepalive for long idle gaps.

### 46.9 War story

A WebRTC app worked for 80% of calls; the rest had one-way or no audio. Logs
showed those clients were behind **symmetric NAT** (corporate firewalls) and the
deployment had **no TURN server** ("STUN is enough"). ICE found only srflx
candidates that didn't pair. Standing up a TURN server (with TLS/443 transport
for the strictest networks) fixed connectivity for the remaining 20%.
**Lesson:** STUN handles ~80–90%; you **must** run TURN for the rest, and expose
it on 443/TLS for locked-down networks.

**Try it yourself**
```bash
# Classify your NAT / see your reflexive address:
stunclient stun.l.google.com 19302            # from 'stuntman'
# or:
turnutils_stunclient stun.l.google.com        # from coturn
# Watch NAT mapping lifetime: send one UDP packet out, then see how long return works
sudo conntrack -L -p udp                        # on a Linux NAT box: see mappings + timeouts
# WebRTC internals in Chrome:  chrome://webrtc-internals  (ICE candidate pairs, selected pair)
```

---

## 47. Why VoIP and video use UDP, jitter buffers, and RTP

### 47.1 The timeliness-over-reliability trade

For a live voice/video stream, a packet that arrives **after its playout deadline
is useless** — worse than useless, because waiting for its retransmission stalls
everything after it (**head-of-line blocking**). So real-time media:
- rides **UDP** (no retransmission, no in-order guarantee, no connection stall),
- accepts loss and **conceals** it (interpolate a lost 20 ms audio frame; freeze
  or blur a lost video block until the next keyframe),
- does its **own** congestion control tuned for low latency.

### 47.2 RTP / RTCP

**RTP (Real-time Transport Protocol, RFC 3550)** over UDP carries the media:

```
RTP header (12 B): V|P|X|CC|M|PT | sequence number(16) | timestamp(32) | SSRC(32)
```
- **Sequence number:** detect **loss** and **reordering** (RTP itself doesn't
  fix them; the app does).
- **Timestamp:** sampling instant (e.g. +160 per 20 ms G.711 frame) → drives
  playout timing and lets the receiver compute **jitter**; independent of arrival
  time.
- **SSRC:** identifies the stream source (each participant/camera).
- **Payload Type (PT):** codec (Opus, VP8/VP9/AV1, H.264…).
- **Marker (M):** e.g. first packet after silence, or last packet of a video
  frame.

**RTCP** (the control sibling, usually next UDP port up, or multiplexed via
**RTCP-MUX** on the same port) carries **Sender/Receiver Reports**: packet
counts, **cumulative loss**, **interarrival jitter**, and **round-trip time**
(via LSR/DLSR timestamps). These reports are the feedback signal for congestion
control and for the sender to switch bitrate/resolution.

Modern stacks also use **RTP Extensions**: transport-wide congestion control
feedback (**transport-cc**), abs-send-time, NACK (selective retransmit for
video), **RED/FEC** (forward error correction), **PLI/FIR** (please send a
keyframe).

### 47.3 The jitter buffer

The network delivers packets with **variable** delay (jitter). If you played each
packet the instant it arrived, audio would stutter. The **jitter buffer** (a.k.a.
de-jitter buffer) holds arriving packets briefly and releases them on a **smooth
clock** derived from RTP timestamps.

- **Too small** → late packets miss playout → dropouts/glitches.
- **Too large** → added mouth-to-ear latency (conversation feels laggy; >150 ms
  one-way is noticeable, >400 ms is bad).
- **Adaptive** jitter buffers grow during turbulence and shrink during calm,
  trading latency for concealment dynamically. Typical target: cover ~p95–p99 of
  observed jitter, often 20–120 ms.

Mouth-to-ear budget ≈ `capture + encode + packetize + network + jitter buffer +
decode + playout`. The jitter buffer and network are where the fight is.

### 47.4 Loss concealment & recovery

| Technique | Cost | Use |
|---|---|---|
| **PLC (Packet Loss Concealment)** — synthesize the missing audio frame | free | every voice codec (Opus in-band) |
| **FEC (Forward Error Correction)** — send redundant parity so N losses are recoverable without retransmit | bandwidth (10–50%) | Opus in-band FEC for audio; flexible FEC / RED for video |
| **NACK + selective retransmit** — receiver asks for a specific lost RTP seq; only useful if it can arrive before deadline | one RTT | video (short-RTT), never for live audio on long paths |
| **Keyframe request (PLI/FIR)** — after unrecoverable video loss, ask sender for a fresh intra frame | big bitrate spike | video |
| **Reference picture selection / temporal layers** — decode continues from an older good frame | complexity | SVC video (VP9/AV1) |

### 47.5 Congestion control for media

Media CC must keep **latency low** (can't fill buffers like CUBIC) and adapt the
**encoder bitrate**, not just the send window:
- **GCC (Google Congestion Control)** — WebRTC's default: a **delay-based**
  controller (watches one-way delay gradient from transport-cc feedback / abs-
  send-time) combined with a **loss-based** controller; output is a target
  bitrate the encoder must hit. Reacts *before* queues build.
- **NADA**, **SCReAM** — IETF alternatives, similar goals.
- On congestion the app **drops resolution/framerate/bitrate** (and may pause
  video to protect audio) — graceful degradation instead of freeze.

### 47.6 Transport choices in the wild

| Protocol | Transport | Notes |
|---|---|---|
| **WebRTC** | UDP + **DTLS-SRTP** (media encrypted), ICE for pathing, SCTP-over-DTLS for data channels | falls back to TURN/TCP/443 when UDP blocked |
| **SIP** | UDP/TCP/TLS 5060/5061 signaling; RTP/UDP media | classic VoIP/PBX |
| **RTP/SRTP** | UDP | media plane for SIP/WebRTC |
| **QUIC / RTP-over-QUIC / MoQ** | UDP | emerging: per-stream reliability lets you mix reliable + unreliable media |
| **RTMP** | **TCP** 1935 | legacy ingest (streamer → server); TCP is OK because it's one-way with a multi-second buffer, not conversational |
| **HLS / MPEG-DASH / LL-HLS** | **TCP/HTTP(S)** | *streaming* (not real-time) — 2–30 s (or ~2–5 s low-latency) glass-to-glass; TCP + big buffer is fine; scales on CDNs |
| **SRT / RIST** | UDP + ARQ | contribution-grade "reliable UDP" for broadcast backhaul; tunable latency budget |

Rule: **conversational** (two-way, <400 ms) → UDP + RTP + jitter buffer +
concealment. **One-way broadcast/VOD** where seconds of latency are fine → TCP +
HTTP + CDN.

### 47.7 War story

A telehealth product had "robotic audio" for clinics on a particular ISP. RTCP
Receiver Reports showed **2–4% loss** and **interarrival jitter ~90 ms** on that
path; the app's jitter buffer was fixed at 60 ms, so ~p80 of packets were "late"
and discarded. Two changes: (1) enable **Opus in-band FEC** + raise packet loss
perceptual robustness, (2) switch to an **adaptive** jitter buffer (60–200 ms).
Audio MOS went from ~2.8 to ~4.1 with only ~25 ms added latency.
**Lesson:** for real-time media, watch RTCP loss & jitter, and match the jitter
buffer to the measured jitter distribution — plus FEC for the residual loss.

**Try it yourself**
```bash
# Inspect RTP/RTCP in a capture:
#   Wireshark: Telephony > RTP > RTP Streams  (shows lost %, max jitter, delta)
#              select a stream > Analyze  > "Play Streams" for audio
tshark -r call.pcap -q -z rtp,streams
tshark -r call.pcap -Y rtcp -T fields -e rtcp.ssrc.fraction -e rtcp.ssrc.jitter
# Chrome: chrome://webrtc-internals  -> per-stream packetsLost, jitter, jitterBufferDelay,
#         framesDropped, availableOutgoingBitrate, googCurrentDelayMs
```

---
# Part XV — Kernel and NIC internals

## 48. The packet's path through the Linux kernel

### 48.1 Receive path (RX), top to bottom

```
 wire
  |
  v
NIC hardware
  - matches dst MAC (or promisc), verifies Ethernet FCS
  - RSS: hashes the packet's 5-tuple -> picks one of N RX queues
  - DMA-writes the frame into a ring-buffer slot (rx ring) in RAM
  - raises an IRQ (or is polled)
  |
  v
Hardware IRQ handler (tiny)  ->  schedules NAPI softirq (NET_RX_SOFTIRQ) on that CPU
  |
  v
NAPI poll (softirq, per RX queue / per CPU)
  - pulls a budget of frames from the ring
  - GRO: merges consecutive same-flow segments into one big super-frame
  - builds sk_buff (skb) for each
  |
  v
netif_receive_skb  ->  tc ingress (qdisc/clsact, XDP already ran earlier at driver)
  |
  v
IP layer:  ip_rcv  -> netfilter PREROUTING (conntrack, DNAT)
  - routing decision: local? -> ip_local_deliver ;  forward? -> ip_forward
  - netfilter: FORWARD or INPUT chains
  |
  v
L4:  tcp_v4_rcv / udp_rcv
  - 5-tuple hash -> find the socket
  - TCP: sequence/ACK processing, into the socket receive queue (or out-of-order queue)
  |
  v
Socket receive buffer (sk_rcvbuf)  ->  wakes the process
  |
  v
recv()/read() syscall copies bytes into user space (or splice/mmap/io_uring)
```

### 48.2 Transmit path (TX), top to bottom

```
send()/write()  -> copies user bytes into sk_sndbuf
  |
  v
TCP: segmentation into MSS units (or defer via TSO), builds headers, seq numbers,
     congestion window check, pacing timestamp
  |
  v
IP: route lookup (dst cache), fills IP header, netfilter OUTPUT/POSTROUTING (SNAT)
  |
  v
Neighbor (ARP/ND): resolve next-hop MAC, build Ethernet header
  |
  v
tc egress qdisc  (fq_codel / fq / cake / htb ...)  <- where shaping & AQM live
  |
  v
Driver xmit -> TX ring slot -> NIC DMA-reads the frame
  - NIC offloads: TSO splits the big segment into MTU-sized frames on the wire,
    computes L3/L4 checksums, optionally inserts VLAN tag
  |
  v
wire
```

### 48.3 Where packets get dropped (and the counter that shows it)

| Stage | Drop reason | Where it shows |
|---|---|---|
| NIC RX ring full | softirq/CPU can't keep up; RX ring too small | `ethtool -S`: `rx_no_buffer`, `rx_missed_errors`, `rx_fifo_errors`; `ip -s link` `RX dropped` |
| `netdev_max_backlog` full | per-CPU backlog between NAPI and stack overflowed | `/proc/net/softnet_stat` col 2 (drops); raise `net.core.netdev_max_backlog` |
| softirq budget exhausted | `time_squeeze` — poll ran out of budget/time | `/proc/net/softnet_stat` col 3; tune `net.core.netdev_budget`/`netdev_budget_usecs` |
| conntrack table full | NAT/stateful firewall state exhausted | `dmesg`: `nf_conntrack: table full, dropping packet`; raise `nf_conntrack_max` |
| socket rcvbuf full (UDP) | app too slow / buffer too small | `nstat` `UdpRcvbufErrors`; `ss -u` Recv-Q |
| socket rcvbuf / ofo prune (TCP) | memory pressure, lots of out-of-order | `nstat` `TCPRcvQDrop`, `PruneCalled`, `OfoPruned` |
| accept queue full | app not `accept()`ing fast enough | `nstat` `ListenOverflows`/`ListenDrops` |
| qdisc drop | shaper/AQM dropped (this is often *intentional* — CoDel) | `tc -s qdisc show` `dropped N` |
| tc/nftables/iptables rule | explicit policy DROP | rule counters (`iptables -nvL`, `nft list ruleset`) |

### 48.4 IRQ affinity, RPS/RFS, and XPS

- **RSS** (hardware) spreads RX across queues by 5-tuple hash → each queue's IRQ
  is pinned to a CPU (`/proc/interrupts`, `irqbalance` or manual
  `smp_affinity`). Goal: spread load, keep a flow on one CPU (cache locality),
  don't land NIC IRQs on the CPUs your app threads run on.
- **RPS** (software RSS) — spread to more CPUs when the NIC has few queues.
- **RFS** — steer a flow's processing to the CPU where the *application* thread
  that owns the socket runs (best cache behavior). `rps_flow_cnt`,
  `rps_sock_flow_entries`.
- **XPS** — pick the TX queue based on the sending CPU.
- **aRFS** — hardware-accelerated RFS (NIC steers based on the kernel's hint).

Misconfigured affinity = one CPU at 100% `si` (softirq) while others idle, and
packet drops at moderate pps. Check `mpstat -P ALL 1` for a lopsided `%soft`.

### 48.5 War story

A 25 GbE cache server dropped packets and spiked p99 at only ~4 Mpps. `mpstat`
showed **CPU0 at 100% softirq**, others idle. `irqbalance` had parked all 16 NIC
queue IRQs on CPU0-3, and RPS was off. Pinning each queue's IRQ to a distinct
core, enabling **RFS**, and moving the app off the IRQ cores took it to ~14 Mpps
with p99 halved and zero drops. **Lesson:** at high pps the bottleneck is usually
**per-CPU softirq**, not bandwidth; fix IRQ affinity + RFS before buying NICs.

**Try it yourself**
```bash
ethtool -S eth0 | grep -Ei 'drop|miss|err|nobuf|fifo'
ethtool -g eth0                          # ring sizes; ethtool -G eth0 rx 4096
ethtool -l eth0                          # queue counts;  -L to change
cat /proc/net/softnet_stat               # per-CPU: processed, dropped, time_squeeze (hex)
grep eth0 /proc/interrupts               # which CPUs take NIC IRQs
mpstat -P ALL 1                          # look for one CPU pegged on %soft
cat /sys/class/net/eth0/queues/rx-0/rps_cpus
```

---

## 49. Offloads: checksum, TSO/GSO/GRO, LRO, RSS, and XDP

### 49.1 The offloads and what they do

| Offload | Direction | What it does | Effect |
|---|---|---|---|
| **Checksum offload (CSO)** | TX/RX | NIC computes/verifies IP & TCP/UDP checksums | frees CPU; **local captures show "checksum incorrect"** because tcpdump sees the packet *before* the NIC fills it in — normal, ignore |
| **TSO (TCP Segmentation Offload)** | TX | kernel hands the NIC one big (up to 64 KB) TCP "segment"; NIC chops it into MTU-sized frames + per-frame headers/checksums | massively fewer trips through the stack; ~2–4× throughput per core |
| **GSO (Generic Segmentation Offload)** | TX | same idea done in software just before the driver — for NICs/protocols without hardware TSO (also **USO** for UDP, **GSO for GRE/VXLAN**) | most of TSO's benefit, hardware-independent |
| **GRO (Generic Receive Offload)** | RX | merges consecutive same-flow segments into one large skb before the stack processes them (strict rules so it's reversible) | fewer stack traversals; **also inflates what a capture shows as one 40 KB "segment"** |
| **LRO (Large Receive Offload)** | RX | hardware version of GRO, but **lossy/aggressive** (can't always be un-merged) | **disable LRO on routers/bridges/hypervisors** — it corrupts forwarding; GRO is the safe choice |
| **RSS (Receive Side Scaling)** | RX | hardware hashes 5-tuple → spreads flows across RX queues/CPUs | multicore scaling |
| **RFS/aRFS** | RX | steer a flow to the CPU running its socket's app thread | cache locality |
| **VLAN offload** | TX/RX | NIC inserts/strips 802.1Q tag | |
| **UDP GSO/GRO** | TX/RX | segmentation/aggregation for UDP — critical for **QUIC** performance | closes much of the QUIC-vs-TCP CPU gap |
| **Tunnel offloads** | TX/RX | checksum/TSO awareness *inside* VXLAN/GENEVE/GRE | needed for overlay networks (Ch. 56) to not tank throughput |
| **TLS offload (kTLS + NIC)** | TX/RX | kernel (and optionally NIC) does record encryption | high-throughput HTTPS/CDN |

### 49.2 Managing them with `ethtool -k`

```bash
ethtool -k eth0                       # show all offload settings
ethtool -K eth0 tso off gso off gro off   # turn off (for debugging / capture accuracy)
ethtool -K eth0 lro off               # keep LRO OFF on anything that forwards
ethtool -K eth0 rx-udp-gro-forwarding on   # helps forwarded QUIC/UDP throughput
```

**When to disable offloads:**
- **Accurate packet capture / IDS / timing analysis** — GRO/LRO/TSO merge or
  split packets so a capture no longer reflects the wire. Disable GRO+LRO (and
  optionally TSO/GSO) on the box while capturing, or capture on a mirror/TAP port.
- **A router / bridge / firewall / hypervisor host** — LRO must be off; GRO is
  usually fine and re-segmented correctly on forward.
- **Debugging weird MTU/checksum bugs.**
- Otherwise: **leave offloads on** — they're a large multiple of throughput per
  CPU.

### 49.3 XDP and AF_XDP (kernel bypass, sort of)

**XDP (eXpress Data Path)**: run an eBPF program **in the driver**, on the raw
RX descriptor, *before* an `sk_buff` is even allocated. Verdicts: `XDP_DROP`
(discard — line-rate DDoS filtering, tens of Mpps per core), `XDP_TX` (bounce
back out — load balancer, e.g. **Katran**, **Cilium**), `XDP_PASS` (into the
normal stack), `XDP_REDIRECT` (to another NIC or into an **AF_XDP** userspace
socket for full DPDK-like bypass without leaving the kernel framework).

Uses: volumetric DDoS scrubbing, L4 load balancing at the host, packet-level
observability, custom forwarding — all without the cost of the full netstack.

### 49.4 Full kernel bypass: DPDK / netmap

For the extreme end (NFV, trading, 100G+ single-box routing), **DPDK** takes the
NIC away from the kernel entirely: userspace poll-mode drivers, hugepages, busy-
polling cores, zero interrupts. You get ~100+ Mpps but you also reimplement ARP,
routing, TCP (or use a userspace stack like **VPP**, **F-Stack**, **mTCP**), and
you burn whole cores at 100%. Only worth it when the kernel genuinely can't keep
up and XDP isn't enough.

### 49.5 War story

An IDS was "missing attacks" that a PCAP replay clearly contained. On the live
tap host, **GRO** was on, so the IDS saw coalesced 40 KB pseudo-segments and its
signature engine (expecting wire-sized packets, checking TCP flags per segment)
mis-parsed streams. `ethtool -K eth0 gro off lro off` on the capture interface
restored detection. Throughput of the *monitored* traffic was unaffected (it's a
passive tap). **Lesson:** any tool that reasons about individual packets —
IDS/IPS, latency measurement, forensic capture — needs GRO/LRO **off** on its
capture path, or a hardware TAP that sees the true wire.

**Try it yourself**
```bash
ethtool -k eth0 | grep -E 'tcp-segmentation|generic-.*-offload|large-receive|rx-checksumming|scatter'
# Before a capture on the host itself:
sudo ethtool -K eth0 gro off lro off ; sudo tcpdump ... ; sudo ethtool -K eth0 gro on
# See GRO in action: capture with offloads on vs off, compare 'length' values
ip -s link show eth0                    # RX/TX packets, dropped, errors
bpftool prog show                        # any XDP/eBPF programs attached?
ip link show eth0 | grep -o 'xdp.*'
```

---

# Part XVI — Scaling the network

## 50. Anycast, ECMP, and BGP in one page

### 50.1 BGP — how the internet's routes are chosen

**BGP (Border Gateway Protocol, RFC 4271)** is the path-vector protocol that
glues ~75 000 **Autonomous Systems** (ASes — networks under one administrative
routing policy, each with an **ASN**) into one internet.

- Neighboring routers form a **TCP session (port 179)** and exchange **route
  advertisements**: `prefix → AS_PATH` (the list of ASes to traverse) plus
  attributes.
- A router picks the **best path** per prefix using a decision process roughly:
  highest `LOCAL_PREF` → shortest `AS_PATH` → lowest `MED` → eBGP over iBGP →
  lowest IGP cost to next hop → tie-breakers. **Policy (LOCAL_PREF, filters) beats
  path length** — the internet's shape is commercial, not shortest-path.
- **eBGP** between ASes; **iBGP** to distribute those routes inside an AS (needs
  full mesh or route reflectors).
- Relationships: **customer** (pays you, you announce them everywhere),
  **peer** (settlement-free, exchange only each other's customers), **transit
  provider** (you pay them for "the rest of the internet"). "Valley-free"
  routing.
- **Convergence is slow** (seconds to minutes) and there's **no built-in
  authentication** historically → route leaks and hijacks.

### 50.2 BGP incidents you should recognize

- **Prefix hijack:** an AS announces a prefix it doesn't own. Because **longest-
  prefix match wins**, announcing a `/24` out of someone's `/16` **pulls that
  traffic** globally (Pakistan/YouTube 2008; numerous crypto-DNS thefts).
- **Route leak:** an AS re-announces routes it shouldn't (e.g. transit routes to
  a peer), becoming an accidental transit and black-holing/degrading traffic
  (Level3 2017, others).
- **Defenses:** **RPKI + ROV** (Route Origin Validation — cryptographically
  asserts "AS X may originate prefix P"; routers drop ROA-invalid routes),
  **IRR filtering**, **max-prefix limits**, **BGPsec** (path validation, barely
  deployed), and MANRS best practices.

### 50.3 ECMP — Equal-Cost Multi-Path

When several next hops have equal cost, a router **hashes each flow's 5-tuple** to
pick one — spreading load while keeping every packet of a flow on **one** path
(so TCP doesn't see reordering). Consequences:
- A single TCP flow **cannot exceed one path's capacity**, even if 4 equal paths
  exist. Big transfers need multiple flows (or MPTCP).
- **`traceroute` shows different IPs per probe** at ECMP hops because each probe
  is a different flow tuple. Use **Paris traceroute** (keeps the tuple constant)
  to get a coherent path.
- A **hash polarization** bug (every router using the same hash) can send all
  traffic down one member. Good implementations seed the hash per device.
- **Flowlet switching / adaptive load balancing** (data-center fabrics) rebalance
  on idle gaps to avoid elephant-flow collisions.

### 50.4 Anycast — one address, many locations

The **same IP prefix is announced via BGP from many sites**. Each client's
packets are routed by normal BGP to the **topologically nearest** announcement.
Used by:
- **Root/TLD DNS** and public resolvers (`1.1.1.1`, `8.8.8.8`).
- **CDNs** — an anycast VIP fronts hundreds of POPs; you hit the closest.
- **DDoS scrubbing** — spread a volumetric attack across global capacity.

Caveats:
- **Stateless-friendly, stateful-fragile:** UDP DNS is ideal. For **TCP/TLS**,
  a mid-connection BGP reconvergence can re-anchor the flow to a *different* site
  that has no state → RST. In practice reconvergence is rare enough, and each POP
  terminates TLS locally, so it works — but long-lived connections (websockets)
  are more exposed. Some designs use anycast only to *find* a POP, then hand off
  to a unicast address.
- **"Nearest" is BGP-nearest, not geographically nearest** — a peering quirk can
  send Lisbon to London instead of Madrid.

### 50.5 How a CDN request actually gets routed (tying it together)

1. **DNS** returns either an **anycast** VIP (Cloudflare, some Google) or a
   **unicast** IP chosen by the CDN's DNS based on the resolver's location / EDNS
   Client Subnet (Akamai-style).
2. Packets reach the nearest/edge **POP** via BGP (+ ECMP inside).
3. The POP's **L4 LB** (Maglev/Katran/XDP, consistent-hash) picks an edge proxy.
4. The edge proxy **terminates TCP+TLS**, serves from cache, or opens a
   connection over the CDN's **private backbone** (optimized routes, warm
   connections, BBR) to the origin.
5. Result: your 3-RTT TLS setup is against a server ~10 ms away, not the origin
   ~150 ms away; only cache misses pay the long haul, and even those ride a tuned
   backbone.

### 50.6 War story

A company's `/24` for its API suddenly saw traffic vanish in three regions. A
small ISP had fat-fingered a static route redistribution and **originated their
`/24`**; with no RPKI ROA published, many networks accepted the bogus origin
(equal prefix length, shorter AS_PATH from those networks' view). Mitigation:
publish an **RPKI ROA** for the prefix, announce more-specific `/25`s temporarily
to out-compete the hijack, and contact upstreams. Long-term: ROAs for every
prefix + monitoring (BGPStream/Alerts). **Lesson:** if traffic disappears from
*some* of the internet but not all, suspect a BGP hijack/leak; RPKI ROAs are the
seatbelt.

**Try it yourself**
```bash
# Public looking glasses / RIPEstat:
curl -s "https://stat.ripe.net/data/routing-status/data.json?resource=93.184.216.0/24" | jq .
whois -h whois.radb.net -- '-i origin AS15169' | grep route | head   # prefixes an AS originates
mtr -z example.com          # -z shows the AS number per hop
traceroute -A example.com   # -A appends [ASxxxx]
# Is a prefix RPKI-valid? (Cloudflare's checker)
#   https://rpki.cloudflare.com/   or   `routinator`/`rpki-client` locally
dig +short CHAOS TXT id.server @1.1.1.1     # which anycast POP answered you
```

---

## 51. QoS, DSCP, and traffic shaping

### 51.1 Why QoS exists

On a congested link, a FIFO queue treats a bulk backup and a voice call
identically → the call suffers. **QoS** = deliberately treating some packets
better (lower latency, lower loss, guaranteed rate) at the expense of others.
QoS **only matters at a congestion point** and **only helps if the bottleneck is
a link you control** — it does nothing across the public internet (which ignores
your DSCP or remarks it to 0).

### 51.2 Marking: DSCP and the IP header

The IPv4 **DSCP** field (6 bits, top of the old ToS byte; IPv6 Traffic Class) +
2 ECN bits. Standard code points (RFC 4594):

| DSCP name | Value (dec) | Intended for | Queue treatment |
|---|---|---|---|
| **EF** (Expedited Forwarding) | 46 | VoIP media, real-time | strict-priority / low-latency queue, policed |
| **CS6 / CS7** | 48 / 56 | network control (routing protocols) | protected |
| **AF41–AF43** | 34/36/38 | interactive video | assured bandwidth, 3 drop precedences |
| **AF31–AF33** | 26/28/30 | signaling, "mission critical" data | assured |
| **AF21–AF23** | 18/20/22 | "transactional" data | assured |
| **AF11–AF13** | 10/12/14 | bulk / "high throughput" | assured, low priority |
| **CS0 / Default (BE)** | 0 | everything else | best effort |
| **LE** (Lower Effort, RFC 3662/8622) | 1 | scavenger (backups, updates) | first to drop |

Marking happens at the **trust boundary** (the switch/router port facing an
endpoint): either honor what the endpoint set (trusted phones), or **classify and
re-mark** based on ACLs/DPI/app. Internal core routers then just act on the mark.

### 51.3 Queueing disciplines (what the mark triggers)

- **Priority Queuing (PQ) / Low-Latency Queuing (LLQ):** EF gets a strict-priority
  queue (always served first) but **policed** to a cap so it can't starve
  everything.
- **CBWFQ / DWRR / HTB classes:** each class gets a **weight** → a guaranteed
  share of the link under congestion, can borrow idle capacity.
- **WRED / CoDel / PIE per class:** drop/mark early within a class to keep its
  queue short (and give TCP a congestion signal before the buffer is full).
- **Shaping vs policing:**
  - **Policing** = drop (or re-mark) anything over rate *now*. Bursty, cheap, TCP-
    unfriendly (causes retransmit storms).
  - **Shaping** = **buffer** and release at the target rate (token bucket). Smooth,
    adds latency, TCP-friendly. Use shaping toward the customer, policing at
    ingress against abuse.

### 51.4 Linux `tc` shaping (the practical bit)

```bash
# HTB: cap the interface at 90 Mbit, carve classes:
tc qdisc add dev eth0 root handle 1: htb default 30
tc class add dev eth0 parent 1:  classid 1:1  htb rate 90mbit ceil 90mbit
tc class add dev eth0 parent 1:1 classid 1:10 htb rate 20mbit ceil 90mbit prio 0  # VoIP
tc class add dev eth0 parent 1:1 classid 1:20 htb rate 50mbit ceil 90mbit prio 1  # interactive
tc class add dev eth0 parent 1:1 classid 1:30 htb rate 20mbit ceil 90mbit prio 2  # bulk/default
# low-latency AQM inside each class:
tc qdisc add dev eth0 parent 1:10 fq_codel
tc qdisc add dev eth0 parent 1:20 fq_codel
tc qdisc add dev eth0 parent 1:30 fq_codel
# classify by DSCP EF -> 1:10
tc filter add dev eth0 parent 1:0 protocol ip u32 match ip dsfield 0xb8 0xfc flowid 1:10
```

Or just use **CAKE**, which does shaping + fairness + DSCP-aware tins + AQM in one
line: `tc qdisc replace dev eth0 root cake bandwidth 90Mbit diffserv4 besteffort`.

### 51.5 802.1p (Layer 2 QoS) and where it applies

The 3-bit **PCP** field in a VLAN tag (Ch. 5) carries a priority class **on that
L2 segment only** — used by switches for per-port egress queues, and by Wi-Fi as
**WMM** access categories (voice/video/best-effort/background) that change the
radio backoff timing. DSCP (L3) survives across routers; PCP (L2) is rewritten
every hop.

### 51.6 War story

A branch office's VoIP was fine until 9am, then choppy. The 20 Mbit/s WAN link
saturated with cloud-backup and Windows-update traffic each morning; all traffic
shared one FIFO on the router, so voice packets queued behind 1500-byte bulk
frames — ~50–150 ms of jitter. Fix: on the WAN egress, **shape to 19 Mbit/s**
(just under line rate, to own the queue), put **DSCP EF** voice in an LLQ capped
at 3 Mbit/s, mark updates/backups **LE**, `fq_codel` inside each class. Jitter
dropped to <10 ms with no bandwidth added. **Lesson:** QoS is queue management at
*your* bottleneck; shape just below line rate so the queue forms where you can
prioritize, and scavenger-mark the bulk traffic.

**Try it yourself**
```bash
tc -s qdisc show dev eth0            # per-class sent/dropped/backlog
tc -s class show dev eth0
# See DSCP on packets:
sudo tcpdump -i eth0 -v 'ip and (ip[1] & 0xfc) != 0' | grep -o 'tos 0x[0-9a-f]*'
# Set DSCP on your own test traffic:
ping -Q 0xb8 <peer>                  # ToS/DSCP EF
iperf3 -c <peer> --tos 184           # 184 = 0xb8 = DSCP 46 (EF)
# Verify it survives the path (often it won't across the internet):
sudo tcpdump -i eth0 -v 'host <peer>' | grep tos
```

---

# Part XVII — Security

> This part is defensive: recognize attacks in traffic and counter them. TCP/IP
> was designed with security as priority #7 (Ch. 1), so most of these are
> structural, not bugs.

## 52. Spoofing, SYN floods, RST injection, and off-path attacks

### 52.1 IP source spoofing — the root enabler

Nothing in IP proves the source address. A host can put **any** source IP in a
packet it sends. This works for the attacker whenever they don't need to see the
reply (**one-way / blind** attacks) or they're **on-path** (can see replies).

**Defense — BCP 38 / BCP 84 (ingress/egress filtering):** every network drops
outbound packets whose source isn't from its own prefixes, and drops inbound
packets claiming to be *from* its own prefixes. **uRPF** (unicast Reverse Path
Forwarding) automates it: drop a packet if the route back to its source doesn't
point out the interface it arrived on (strict) or isn't in the table at all
(loose). Widely deployed but not universal — which is why spoofed-source DDoS
still exists.

### 52.2 SYN flood

Send many `SYN`s (usually spoofed source) and never complete the handshake. Each
fills a slot in the server's **SYN queue** (half-open connections), and the
server wastes a SYN-ACK + retransmits. Legit clients can't get a slot → denial of
service, at very low attacker bandwidth.

**Defenses:**
- **SYN cookies** (`net.ipv4.tcp_syncookies=1`, Ch. 24.5) — server keeps **no
  state** for half-open connections; reconstructs from the final ACK. The primary
  defense; engages only under overflow.
- Larger `tcp_max_syn_backlog`, shorter `tcp_synack_retries`.
- **SYN proxy** at a firewall/LB/XDP (`synproxy` in netfilter, or hardware) —
  completes the handshake on the server's behalf and only hands over *real*
  connections.
- Upstream **scrubbing** for volumetric versions.
- Rate-limit SYNs per source (weak vs spoofing) / per subnet.

### 52.3 RST injection & connection reset attacks

If an attacker can **guess or observe** the 4-tuple and a valid sequence number
window, they can forge a `RST` (or data) segment and the endpoint will act on it:
- **On-path** (same LAN, compromised router, nation-state middlebox): trivial —
  they see everything. The **Great Firewall** forges RSTs to both ends to kill
  connections; some ISPs used RST to throttle P2P.
- **Off-path (blind):** must guess the tuple (server IP/port known; client IP
  maybe known; client **port** and **sequence** are the entropy). RFC 5961 makes
  this much harder: an in-window-but-not-exact RST triggers a **challenge ACK**
  instead of an immediate reset, so the attacker needs a near-exact sequence
  number. (Ironically, the shared *rate limit* on challenge ACKs became a
  side-channel — CVE-2016-5696 — now randomized.)

**Defenses:** randomized ISNs (RFC 6528), RFC 5961 checks, ephemeral **source
port randomization** (more tuple entropy), and — the real fix — **TLS/SSH**: a
forged RST can tear the TCP connection but can't inject or read application data,
and the app just reconnects. **TCP-AO** (RFC 5925, replaces TCP-MD5) authenticates
segments for infrastructure sessions like BGP.

### 52.4 Session hijacking / data injection

On-path attacker injects application data into an unauthenticated stream (classic
`telnet`, `http`, unencrypted SMTP). Or predicts sequence numbers to do it blind
(historic). **Defense:** encrypt + authenticate everything (TLS, SSH, IPsec,
WireGuard). Assume the path is hostile.

### 52.5 TCP-based reconnaissance

- **SYN scan** (`nmap -sS`): SYN → `SYN-ACK` (open) / `RST` (closed) / nothing
  (filtered). Doesn't complete the handshake ("half-open").
- **FIN/NULL/Xmas scans:** rely on RFC-793 behavior (closed port → RST, open →
  nothing) to slip past stateless filters. Modern stacks/firewalls detect these.
- **Idle/zombie scan:** uses a third host's predictable IP **ID** field to scan
  fully blind. Defense: randomized IP ID (Linux does per-destination).
- **OS fingerprinting** (`nmap -O`, `p0f`): initial TTL, window size, MSS, option
  order/presence, DF bit → identifies the OS. Defense: normalization at a
  firewall (`scrub` in pf), or accept that it's low-risk.
- **Banner grabbing**, TLS `JA3`/`JA4` client fingerprints.

**Defenses:** default-deny firewall, rate-limit + alert on scan patterns
(`nf_conntrack` + `hashlimit`, fail2ban, an IDS), don't expose management ports,
and reduce attack surface (fewer listening sockets — audit with `ss -tlnp`).

### 52.6 Amplification / reflection (covered more in Ch. 54)

Spoof the victim's IP as source of a **small** UDP request to a service that
returns a **large** response → the victim is flooded, attacker's bandwidth
multiplied. Amplification factors: DNS ANY ~50×, NTP `monlist` ~550×, memcached
~10 000–50 000×, SSDP ~30×, CLDAP ~55×. **Defense:** BCP 38 (kills spoofing at
source), disable/restrict the abusable services (no open resolvers, NTP
`noquery`, memcached not on UDP/not on the internet), rate-limit responses
(DNS **RRL**).

### 52.7 Other classics to recognize

| Attack | Mechanism | Defense |
|---|---|---|
| **ARP spoofing** | forged ARP → LAN MITM | Dynamic ARP Inspection + DHCP snooping, static ARP for gateways, 802.1X |
| **Rogue DHCP / rogue RA** | hand out attacker as gateway/DNS | DHCP snooping, RA Guard |
| **DNS cache poisoning** (Kaminsky) | race a forged answer into a resolver | source-port + query-ID randomization (32-bit entropy), **DNSSEC**, DoT/DoH to a trusted resolver, **0x20** encoding |
| **DNS rebinding** | public name flips to `127.0.0.1`/RFC1918 after TTL → browser attacks internal services | resolver/app **rebinding protection** (reject private answers for public names), `Host` allow-listing |
| **TCP "sequence prediction"** | guess ISN → blind spoof a whole connection | RFC 6528 randomized ISN |
| **Slowloris / slow POST / R-U-Dead-Yet** | open many connections, send headers/body 1 byte at a time → exhaust worker slots | request header/body **timeouts**, min data-rate, connection limits per IP, an async/event-driven front proxy (nginx/HAProxy) that buffers |
| **TCP request smuggling** | front proxy & backend disagree on request boundaries (`CL` vs `TE`) | normalize/reject ambiguous framing at the edge, HTTP/2 end-to-end |
| **LAND / teardrop / ping-of-death** | malformed/overlapping packets crash old stacks | patched kernels; drop obviously bogus packets |

### 52.8 War story

A service saw thousands of connections stuck in `SYN_RECV` and legit users
timing out. `nstat` showed `TcpExtTCPReqQFullDoCookies` climbing — a SYN flood
with spoofed sources. `tcp_syncookies` was already on and doing its job at the
kernel, but the upstream firewall's own connection table was also filling. Adding
an **XDP SYN filter** on the edge hosts (drop SYNs exceeding a per-source rate,
validate with cookies before passing) plus asking transit to enable uRPF on the
customer link absorbed it. **Lesson:** SYN floods are cheap and old; have SYN
cookies on, know your `nstat` cookie counters, and have an XDP/scrubbing layer
ready.

**Try it yourself (defensive)**
```bash
ss -tlnp                                  # audit every listening socket - attack surface
nstat -az | grep -Ei 'syncookie|challengeack|embryonic|ReqQFull|PAWS|ListenDrops'
sysctl net.ipv4.tcp_syncookies net.ipv4.conf.all.rp_filter net.ipv4.tcp_rfc1337
sysctl net.ipv4.conf.all.accept_redirects net.ipv4.conf.all.send_redirects   # set 0
sysctl net.ipv4.conf.all.accept_source_route   # set 0
sudo nft add rule inet filter input tcp flags syn limit rate 200/second accept  # SYN rate-limit
```

---

## 53. Firewalls, stateful inspection, and conntrack

### 53.1 Packet filter generations

1. **Stateless ACL** — match per-packet on 5-tuple + flags. Fast, but can't tell
   a reply from an unsolicited packet, so you must open both directions and
   high-port ranges (fragile, leaky).
2. **Stateful firewall** — tracks **connections**. Allow `NEW` outbound, then
   `ESTABLISHED,RELATED` return traffic is permitted automatically. This is the
   default model everywhere now (netfilter/`nftables`, pf, iptables, cloud
   security groups).
3. **NGFW / L7 firewall** — adds app identification (DPI), TLS inspection
   (MITM with a corporate CA), IDS/IPS signatures, user identity.
4. **WAF** — HTTP-semantic rules (SQLi/XSS/path traversal, request smuggling,
   rate limits).

### 53.2 Linux conntrack — the state table

netfilter's **conntrack** assigns every flow a state entry (the same table NAT
uses, Ch. 12). Per-entry it stores the tuple(s), a **state**, and a **timeout**.

| Conntrack state | Meaning |
|---|---|
| `NEW` | first packet of a flow (a SYN, or first UDP packet) |
| `ESTABLISHED` | traffic seen in **both** directions |
| `RELATED` | a new flow spawned by an existing one (FTP data channel, ICMP error referencing a tracked flow) |
| `INVALID` | doesn't match any flow and isn't a valid `NEW` (e.g. an ACK/RST for an unknown connection, out-of-window) — usually **drop** these |
| `UNTRACKED` | explicitly exempted via `NOTRACK` (perf) |

TCP timeouts (defaults, `nf_conntrack_tcp_timeout_*`): `syn_sent` 120 s,
`established` **432000 s (5 days)**, `time_wait` 120 s, `close` 10 s,
`fin_wait` 120 s. UDP: `nf_conntrack_udp_timeout` 30 s, `*_stream` 120 s.

### 53.3 The failure modes

- **`nf_conntrack: table full, dropping packet`** — table hit
  `nf_conntrack_max`. New flows are silently dropped; existing ones fine. Causes:
  a scan/flood, a connection-churning app, too-long timeouts holding dead
  entries, or just growth. Fixes: raise `nf_conntrack_max` (and the hash
  `hashsize` = max/8), shorten `established` timeout for your workload, drop
  `INVALID` early, or `NOTRACK` high-volume stateless traffic (e.g. a DNS server's
  UDP/53) in the `raw` table.
- **Idle-timeout drops** (Ch. 34) — the firewall reaps a legit idle connection;
  the next packet is `INVALID` → RST or drop. Fix with keepalives shorter than
  the timeout.
- **Asymmetric routing** (Ch. 13) — the firewall sees only one direction, so the
  flow never reaches `ESTABLISHED` and return packets look `INVALID`. Fix routing
  or the flow dies.
- **Conntrack + ECMP/HA** — two firewalls in a path must **sync conntrack state**
  (`conntrackd`, or cloud stateful-LB) or a failover/rehash drops every existing
  flow.
- **`INVALID` + `-j DROP` vs default ACCEPT** — many rulesets forget to drop
  `INVALID`, letting crafted out-of-state ACK/RST packets through filters that
  only checked "not NEW."

### 53.4 A sane host firewall shape (nftables)

```nft
table inet filter {
  chain input {
    type filter hook input priority 0; policy drop;
    ct state established,related accept
    ct state invalid drop
    iif "lo" accept
    ip protocol icmp accept
    ip6 nexthdr ipv6-icmp accept          # DO NOT blanket-drop ICMPv6 (Ch. 14/17)
    tcp dport 22 ct state new limit rate 15/minute accept
    tcp dport { 80, 443 } accept
    # everything else: dropped by policy
  }
  chain forward { type filter hook forward priority 0; policy drop; }
  chain output  { type filter hook output  priority 0; policy accept; }
}
```

### 53.5 Cloud security groups vs NACLs

- **Security Groups** (AWS/GCP/Azure NSG): **stateful**, attached to an
  instance/ENI, **allow-only** (implicit deny), evaluated as a set. Return
  traffic auto-allowed.
- **Network ACLs** (AWS subnet-level): **stateless**, ordered rules, explicit
  allow **and** deny, and **you must open ephemeral return ports** (e.g.
  `1024-65535` inbound) because it won't track state.
- Both are distributed enforcement (no single choke point), so "connection table
  full" isn't your problem — but **connection tracking on an NLB/NAT Gateway
  is**: AWS NAT Gateway ~55 000 simultaneous connections **per destination**,
  and per-flow idle timeout 350 s.

### 53.6 War story

A batch platform intermittently failed to reach an internal API; `curl` hung then
`Connection reset`. The firewall logged `INVALID` drops. Cause: the API's TCP
`established` conntrack entry on the firewall timed out after a long quiet phase
of a job (no keepalive), then the job's next request arrived on the dead entry →
`INVALID` → drop, and a later packet → RST. Fix: `TCP_KEEPIDLE=120` on the client
(under the firewall's timeout) **and** shortening the job's inter-request gaps.
Also added `ct state invalid drop` with logging so it was diagnosable next time.
**Lesson:** stateful firewalls silently expire idle flows; long-lived clients
need keepalives, and you should **log** `INVALID` drops.

**Try it yourself**
```bash
sudo sysctl net.netfilter.nf_conntrack_count net.netfilter.nf_conntrack_max
sudo conntrack -L | wc -l
sudo conntrack -L -p tcp --state ESTABLISHED | head
sudo conntrack -S                         # per-CPU: found, invalid, insert_failed, drop
dmesg | grep -i nf_conntrack
sudo nft list ruleset          # or: iptables -nvL --line-numbers  (watch the counters)
sysctl net.netfilter.nf_conntrack_tcp_timeout_established
```

---

## 54. DDoS shapes and mitigations

### 54.1 The three layers of DDoS

| Layer | Examples | Measured in | Goal |
|---|---|---|---|
| **Volumetric (L3/4)** | UDP/ICMP flood, **amplification/reflection** (DNS, NTP, memcached, CLDAP, SSDP), carpet bombing | **bits/s** (Gbps–Tbps) | saturate the pipe / transit |
| **Protocol / state exhaustion (L3/4)** | **SYN flood**, ACK flood, RST flood, fragmented-packet flood, `nf_conntrack` exhaustion, TLS renegotiation | **packets/s**, **connections/s**, table entries | exhaust CPU, NIC pps, firewall/LB state |
| **Application (L7)** | HTTP(S) request flood, "cache-busting" random query strings, expensive-endpoint targeting (search, login, GraphQL), **Slowloris**, HTTP/2 Rapid Reset (CVE-2023-44487) | **requests/s**, concurrency | exhaust app workers, DB, downstream APIs |

### 54.2 Volumetric mitigation

- **You cannot absorb a 1 Tbps flood on a 10 Gbps uplink.** Mitigation has to
  happen **upstream / in the cloud**: an **anycast scrubbing network** (Ch. 50)
  spreads the attack across hundreds of Tbps of global capacity, filters, and
  forwards clean traffic (via GRE/L2 tunnel or by being your CDN/proxy).
- **BGP** tools: **RTBH** (Remotely Triggered Black Hole — announce the victim
  `/32` with a blackhole community so transit drops **all** its traffic; the
  target goes fully dark but the rest of your network survives) and
  **Flowspec** (push fine-grained drop/rate-limit filters — "drop UDP/53
  responses > 1400 bytes to this /24" — into upstream routers).
- **Kill amplification at the source:** BCP 38 everywhere; no open resolvers;
  disable NTP `monlist`; keep memcached off UDP and off the internet.
- On your edge: **XDP/eBPF** drop of obvious junk (spoofed ranges, malformed,
  disallowed protocols) at line rate before the stack; drop IP fragments if your
  services don't need them.

### 54.3 Protocol/state mitigation

- SYN cookies + SYN proxy + large backlogs (Ch. 52).
- Drop `INVALID` (Ch. 53); rate-limit `NEW` per source; raise `nf_conntrack_max`
  or `NOTRACK` stateless services.
- NIC/pps: RSS across all queues, XDP filtering, hardware flow steering; ensure
  you're not single-CPU-softirq-bound (Ch. 48).
- Fragmentation floods: `net.ipv4.ipfrag_high_thresh`, or drop fragments at XDP.

### 54.4 Application-layer mitigation

- **Rate limiting**: per-IP, per-token/API-key, per-URL, sliding-window or token-
  bucket, at the edge proxy/WAF/CDN. Careful with shared NAT/CGNAT IPs — prefer
  keying on auth identity or a proof-of-work/JS challenge for anonymous traffic.
- **Caching**: serve as much as possible from cache so the origin never sees the
  flood; defeat cache-buster query strings by normalizing/ignoring unknown
  params.
- **Challenge**: JS challenge, CAPTCHA, **Proof-of-Work**, or cryptographic
  attestation (Privacy Pass) for suspicious clients; keep good clients on a fast
  path.
- **Concurrency & timeouts**: per-connection request timeouts (Slowloris),
  request/second and in-flight caps per connection, cap HTTP/2 concurrent streams
  and **RST_STREAM rate** (Rapid Reset), body-size and header-count limits.
- **Load shedding & prioritization**: when overloaded, shed anonymous/low-value
  traffic first (return cheap 429/503 with `Retry-After`), protect logged-in
  users and health checks; **circuit-break** to downstreams so the DB isn't
  dragged down.
- **Autoscaling** is a *cost* amplifier if it scales into the attack — cap it and
  combine with shedding.

### 54.5 Anti-amplification in protocols you operate

- **DNS**: Response Rate Limiting (RRL), `minimal-responses`, refuse `ANY` /
  return `HINFO` (RFC 8482), disable recursion for the internet, EDNS buffer
  1232.
- **QUIC**: the 3× amplification limit and address validation via Retry tokens
  (built into the protocol — make sure your library enforces it).
- **NTP/SNMP/SSDP/memcached/CLDAP**: don't expose to the internet; if you must,
  disable the amplifying query types.

### 54.6 Preparation checklist (the real mitigation)

1. Contract/enable an upstream scrubbing/CDN provider **before** an attack;
   know how to trigger it (BGP swing, DNS change, always-on).
2. Know your numbers: uplink Gbps, edge pps ceiling, LB conn-table size, app
   RPS capacity, DB QPS ceiling.
3. Have **RTBH** and **Flowspec** procedures with your transit providers,
   tested.
4. Dashboards + alerts on: inbound bps/pps by protocol, SYN rate, conntrack
   count, `ListenDrops`, 5xx/429 rate, upstream latency.
5. Runbooks: "volumetric UDP" / "SYN flood" / "HTTP flood on /search" each with
   concrete first moves.
6. Rate limits and caching **already deployed** (not "we'll turn them on during
   the attack").
7. Separate control-plane/admin access paths so you can still operate.

### 54.7 War story

An e-commerce site was hit with ~180 Gbps of DNS + CLDAP reflection during a
sale. The transit link (2×10 Gbps) was saturated instantly — nothing on the host
mattered. The on-call triggered the pre-arranged **BGP swing** to the scrubbing
provider (announce prefixes from the scrubber, GRE tunnel back for clean
traffic); within ~4 minutes clean traffic flowed and the site recovered.
Post-mortem actions: move to always-on anycast proxy (no swing delay), publish
**Flowspec** templates with transit, and add RRL to the authoritative DNS.
**Lesson:** volumetric defense is an *upstream* capability you arrange in advance;
host tuning is irrelevant once your pipe is full.

**Try it yourself (observability)**
```bash
# Traffic by protocol, live:
sudo iftop -i eth0                        # or:  nload, bmon
sudo tcpdump -i eth0 -nn -c100000 -w /tmp/s.pcap ; \
  tshark -r /tmp/s.pcap -q -z io,phs      # protocol hierarchy: what's flooding?
nstat -az | grep -Ei 'InReceives|InDiscards|SyncookiesSent|ListenDrops|TCPBacklogDrop'
watch -n1 'cat /proc/net/softnet_stat'    # col2 drops, col3 time_squeeze rising = pps overload
sudo conntrack -S                          # insert_failed / drop rising = state exhaustion
# app edge:
tail -f access.log | awk '{print $1}' | sort | uniq -c | sort -rn | head   # top talkers
```

---


# Part XVIII — Cloud & container networking

## 55. Network namespaces, veth pairs, and bridges

### 55.1 The building block: the network namespace

A **network namespace** (netns) is an isolated copy of the entire network stack:
its own interfaces, IP addresses, routing table, ARP/neighbor table, conntrack,
iptables/nftables rules, and port space. A container is (mostly) **a process in
its own netns**. `localhost` inside it is genuinely separate; two containers can
both bind `:8080`.

```bash
ip netns add red
ip netns exec red ip addr           # only 'lo' (down) - a blank stack
ip netns exec red ip link set lo up
```

### 55.2 veth pairs — a virtual patch cable

A **veth** is a pair of interfaces; whatever enters one end exits the other. Put
one end in the container's netns, leave the other in the host, and you have a
link between them.

```bash
ip link add veth-red type veth peer name veth-red-host
ip link set veth-red netns red
ip -n red addr add 10.10.0.2/24 dev veth-red
ip -n red link set veth-red up
ip addr add 10.10.0.1/24 dev veth-red-host        # host end = the container's gateway
ip link set veth-red-host up
ip -n red route add default via 10.10.0.1
# NAT the container out to the internet:
sysctl -w net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -s 10.10.0.0/24 -o eth0 -j MASQUERADE
```

That is essentially **Docker's default `bridge` network**, minus the bridge.

### 55.3 The Linux bridge — a software switch

To connect **many** containers on one host, put every container's host-side veth
into a **bridge** (`docker0`, `cni0`, `br0`). The bridge does MAC learning and
forwarding exactly like Ch. 7's switch, in software.

```
 container A netns        container B netns
   veth-a (10.10.0.2)       veth-b (10.10.0.3)
      |                        |
   veth-a-host             veth-b-host
      \_______  docker0  _______/     (bridge, 10.10.0.1, also the gateway)
                   |
                  eth0  --(MASQUERADE / SNAT)-->  physical network
```

- Container↔container on the same host/bridge: pure L2, fast, no NAT.
- Container→outside: default route to the bridge IP → host routes → **SNAT** to
  the host's IP.
- Outside→container: **DNAT** (`docker -p 8080:80` inserts an
  `iptables`/`nftables` DNAT rule `hostIP:8080 → 10.10.0.2:80`).
- **`docker-proxy`**: a userspace fallback that also `accept()`s on the published
  port (handles hairpin/localhost cases); a source of confusion in captures.

### 55.4 What this costs / gotchas

- **Extra hops & MTU:** each veth + bridge + NAT adds a little latency and
  (for overlays, Ch. 56) **eats MTU** — pod MTU is often 1450 or 1370. A pod
  MTU mismatch is the #1 cause of "TLS handshake hangs / big responses stall" in
  Kubernetes (same PMTU black-hole signature as Ch. 8).
- **conntrack pressure:** every pod flow is a host conntrack entry; busy nodes
  need `nf_conntrack_max` raised (Ch. 53). `kube-proxy` in iptables/IPVS mode
  adds lots of rules/state.
- **SNAT port exhaustion:** many pods sharing the node IP to one external service
  → ephemeral-port/tuple exhaustion (Ch. 21). Fixes: SNAT to a range of IPs,
  per-pod egress IP, or a dedicated egress gateway.
- **`hairpin` / NodePort quirks:** pod reaching its own Service VIP, or
  NodePort + externalTrafficPolicy (`Cluster` SNATs and hides client IP;
  `Local` preserves client IP but only routes to pods on that node).
- **Bridge + `iptables`:** `br_netfilter` / `net.bridge.bridge-nf-call-iptables`
  must be set for Kubernetes network policy to see bridged traffic.

### 55.5 War story

Pods on one Kubernetes cluster could `curl` small endpoints but every `helm`
pull, image layer, and large API response **hung**. Node `eth0` MTU was 9000
(jumbo), the CNI set pod veth MTU to **9000** too, but the cloud's VPC path
capped at **8900** for cross-subnet traffic and dropped the DF packets while ICMP
"too big" was filtered by a security group. Fixes: set pod/CNI MTU to 8900 (or
1450 to be safe), and allow ICMP type 3 / ICMPv6 type 2 in the security group.
**Lesson:** in container networking, always verify the **end-to-end pod MTU**
against the real underlay path; "small works, big hangs" = MTU.

**Try it yourself**
```bash
ip netns list
ip -n <ns> addr ; ip -n <ns> route ; ip -n <ns> neigh
bridge link ; bridge fdb show               # bridge ports + MAC table
nsenter -t <pid> -n ip addr                  # enter a container's netns by PID
ip link show type veth
# In a pod: check MTU and path
cat /sys/class/net/eth0/mtu
ping -M do -s 1400 <other-pod-ip>            # shrink until it passes = pod path MTU
```

---

## 56. Overlays: VXLAN, and how Kubernetes pods talk

### 56.1 The pod-network requirement

Kubernetes demands: **every pod gets a unique IP**, and **every pod can reach
every other pod without NAT** (across nodes). The **underlay** (the VPC/physical
network) usually doesn't know about pod IPs, so the CNI plugin must bridge that
gap. Two broad approaches:

### 56.2 Approach A — routed / native (no encapsulation)

The pod CIDR is made routable on the underlay:
- **Cloud-native CNIs** (AWS VPC CNI, Azure CNI, GKE): pods get **real VPC IPs**
  from ENIs/secondary ranges. No overlay, MTU stays 1500 (or jumbo), the VPC
  router forwards pod↔pod directly. Cost: consumes VPC address space, ENI/IP
  limits per node.
- **BGP** (Calico in BGP mode, Cilium BGP, kube-router): each node advertises its
  local pod `/26` via BGP to the ToR/underlay. Pure L3, no encap, easy to
  troubleshoot with normal tools. Needs BGP on the fabric.

### 56.3 Approach B — overlay (encapsulation)

Wrap each pod packet inside a **new IP/UDP packet addressed node→node**, so the
underlay only ever sees node IPs.

**VXLAN** (the common one; Flannel default, Calico VXLAN mode, Cilium VXLAN):

```
[ outer Eth ][ outer IP: nodeA -> nodeB ][ outer UDP dport 4789 ]
   [ VXLAN header: 8 bytes, 24-bit VNI ]
      [ inner Eth ][ inner IP: podA -> podB ][ inner payload ]
```

- **VNI** (VXLAN Network Identifier, 24 bits → 16M segments) separates virtual
  networks — the multi-tenant successor to the 12-bit VLAN ID.
- A **VTEP** (VXLAN Tunnel Endpoint) on each node encapsulates/decapsulates. It
  learns "pod MAC/IP → remote node IP" mappings via a control plane (the CNI
  watching the K8s API), or historically multicast/flood-and-learn.
- **Overhead = 50 bytes** (14 inner Eth + 8 VXLAN + 8 UDP + 20 outer IP) → pod
  MTU **1450** on a 1500 underlay. **GENEVE** (Cilium default, OVN) is similar
  with TLV options → slightly more overhead. **IP-in-IP** (Calico) = 20 bytes →
  MTU 1480. **WireGuard** encryption (Calico/Cilium encryption) = ~60–80 bytes.

### 56.4 Performance considerations

- **Offloads must be tunnel-aware.** Without VXLAN-aware checksum/TSO/GRO
  (`ethtool -k`: `tx-udp_tnl-segmentation`, `rx-udp_tnl-...`), every encapsulated
  packet is handled segment-by-segment on the CPU → throughput can drop 3–5×.
  Modern NICs + kernels do this; verify on older ones.
- **Double conntrack, double qdisc**, extra copies. Overlay ~5–15% throughput and
  a few µs latency vs routed mode, more if offloads are missing.
- **eBPF datapaths** (Cilium) can bypass much of `iptables`/`kube-proxy` and
  parts of the bridge/veth path for a measurable win; **`eBPF host-routing`** and
  **Direct Server Return** for Services reduce hops further.
- **`ipvs` vs `iptables` kube-proxy:** iptables mode is O(rules) per packet —
  thousands of Services → latency and reload spikes; IPVS or eBPF scale better.

### 56.5 Service VIPs, DNS, and NetworkPolicy — the TCP/IP view

- A **ClusterIP Service** is a virtual IP with **no host**; `kube-proxy`/eBPF
  installs DNAT rules that rewrite `VIP:port` → a random ready pod IP, tracked in
  conntrack so replies un-DNAT. Load balancing is per-connection.
- **CoreDNS** runs as pods behind a Service; every pod's `/etc/resolv.conf`
  points at the DNS ClusterIP with `search <ns>.svc.cluster.local svc.cluster.local
  cluster.local` and **`ndots:5`** — see Ch. 20.8 for the query-amplification
  trap. `NodeLocal DNSCache` puts a caching resolver on each node to cut latency
  and conntrack churn.
- **NetworkPolicy** is enforced by the CNI (Calico/Cilium) as stateful L3/L4 (and
  L7 with Cilium) filtering on the veth/eBPF path — default-allow until the first
  policy selects a pod, then default-deny for that pod.

### 56.6 War story

After a node-pool upgrade, cross-node pod traffic throughput fell from ~9 Gbit/s
to ~1.5 Gbit/s; same-node was fine. The new node image lacked the NIC driver
feature flag for **VXLAN TSO/checksum offload**, so the CPU segmented every
encapsulated flow. `ethtool -k eth0 | grep udp_tnl` showed the offloads `off
[fixed]`. Rolling to an image with the proper driver (or switching the CNI to
**routed/BGP mode**, eliminating encap) restored line rate. **Lesson:** overlay
throughput depends entirely on **tunnel-aware NIC offloads**; check
`tx-udp_tnl-segmentation` after any kernel/driver/AMI change, or avoid encap with
a routed CNI.

**Try it yourself**
```bash
ethtool -k eth0 | grep -E 'udp_tnl|gre|checksum|segmentation'   # tunnel offloads
ip -d link show flannel.1        # or vxlan.calico, cilium_vxlan - shows VNI, dstport, local
bridge fdb show dev flannel.1     # VTEP forwarding entries (pod MAC -> remote node IP)
sudo tcpdump -i eth0 -nn 'udp port 4789 or udp port 6081'    # VXLAN / GENEVE on the wire
sudo tcpdump -i eth0 -nn 'udp port 4789' -w /tmp/vx.pcap ; wireshark /tmp/vx.pcap  # decodes inner packet
cat /sys/class/net/$(ip route get 1.1.1.1 | awk '{print $5; exit}')/mtu
kubectl exec -it <pod> -- cat /etc/resolv.conf                 # ndots, search domains
```

---

## 57. Cloud VPC constructs and their TCP/IP reality

### 57.1 The mental model

A cloud **VPC** is a **software-defined L3 network**. There is no real broadcast
domain, no real switch you can see; the hypervisor's virtual switch + a
distributed control plane emulate one. Almost every "weird" cloud-network
behavior comes from that emulation.

### 57.2 Construct → what it actually is

| Cloud construct | TCP/IP reality |
|---|---|
| **VPC / VNet** | a routing domain + an IP prefix you chose (e.g. `10.0.0.0/16`). Isolated from other VPCs unless peered. |
| **Subnet** | a slice of the VPC prefix pinned to **one AZ**. **Not** a broadcast domain — there's no ARP flooding; the platform answers ARP for the gateway and knows every instance's MAC/IP. The **first ~3 IPs + broadcast are reserved** by the provider (gateway, DNS, future use). |
| **Route table** | per-subnet. `local` route for the VPC CIDR is implicit and highest priority; you add routes to IGW / NAT GW / peering / TGW / VPN. Longest-prefix match as usual. |
| **Internet Gateway (IGW)** | 1:1 NAT between an instance's **public IP** and its private IP (the instance only ever *sees* its private IP; the mapping is stateless and done in the fabric). No bandwidth limit of its own. |
| **NAT Gateway** | managed **PAT** for private subnets → internet. **Stateful**, ~55k simultaneous connections **per unique destination** `IP:port`, idle timeout 350 s, charged per GB. Port exhaustion and the 350 s timeout are real app-visible limits (Ch. 21, 34). |
| **Security Group** | **stateful** allow-list on the ENI (Ch. 53.5). |
| **Network ACL** | **stateless** subnet filter — remember ephemeral return ports. |
| **ENI (elastic network interface)** | the instance's virtual NIC; carries the MAC, primary+secondary private IPs, SG bindings. Instance type caps **# ENIs and IPs/ENI** (→ pod density limits for VPC-CNI). |
| **ELB/ALB/NLB** | ALB = L7 proxy (terminates, new backend conn, adds `X-Forwarded-For`, HTTP/2 front). NLB = L4 (preserves client IP optionally, ~350 s idle, handles millions of flows, static IP per AZ). GWLB = transparent L3 for appliance insertion (GENEVE-encapsulates to firewalls). |
| **VPC Peering** | routes + fabric permission; **non-transitive** (A–B and B–C ≠ A–C). No encapsulation, MTU stays. |
| **Transit Gateway / vWAN hub** | a managed router connecting many VPCs/VPNs/Direct Connect; **transitive**; often caps MTU (e.g. 8500) and per-flow bandwidth (~5 Gbit/s per flow on AWS TGW — a *single* TCP flow can't go faster; use multiple flows). |
| **PrivateLink / Private Service Connect** | exposes one service across VPCs via an ENI in the consumer VPC + an NLB in the producer VPC; DNAT under the hood; no CIDR overlap concerns. |
| **VPN (site-to-site)** | IPsec tunnels; MTU ~1400 (do MSS clamping — Ch. 8); often ~1.25 Gbit/s per tunnel → use ECMP over multiple tunnels. |
| **Direct Connect / ExpressRoute / Interconnect** | private L2/L3 circuit to the cloud; supports jumbo (up to ~9001/8500); BGP for route exchange. |

### 57.3 MTU inside the cloud

| Path | Typical MTU |
|---|---|
| Within a VPC (same or peered), instance↔instance | **9001** (jumbo) on AWS/GCP |
| Over an Internet Gateway (to the internet) | **1500** |
| Over most VPN / some TGW attachments | **1300–8500** (varies; clamp MSS) |
| Cross-region VPC peering | often **1500** |
| Overlay CNI on top (Ch. 56) | subtract 50 (VXLAN) / 60–80 (encrypted) |

Set instance MTU and pod MTU to match the **narrowest** segment the traffic will
cross, and allow ICMP "too big" through every SG/NACL.

### 57.4 Cloud DNS specifics

- Each VPC has a **resolver at `VPC_base + 2`** (`.2` address) — a hidden
  recursive resolver that also serves private hosted zones and instance names.
- **Per-ENI DNS query throttle** (~1024 packets/s to the `.2` resolver on AWS).
  Chatty apps + `ndots:5` (Ch. 20.8) hit it → intermittent `SERVFAIL`/timeouts.
  Fix: local caching resolver / NodeLocal DNSCache, FQDNs, fewer lookups.
- **Split-horizon** via private hosted zones; on-prem resolution via inbound/
  outbound resolver endpoints.

### 57.5 Bandwidth & pps caps that surprise people

- **Per-instance aggregate** caps (e.g. "up to 10 Gbit/s" = burst, baseline
  lower on small types) and **per-flow** caps (AWS: ~5 Gbit/s within a placement
  group / same-AZ for a single flow, ~10 with ENA Express; less across AZ/region).
  → A single `scp`/`iperf3` stream can't show you the instance's real capacity;
  use `-P` parallel streams.
- **NAT GW / TGW / VPN per-flow** limits as above.
- **"Network burst credits"** on small instances — great in tests, throttled
  under sustained load.
- **PPS limits** independent of bandwidth — small-packet workloads (DNS, RTP,
  trading) hit these first.

### 57.6 War story

A data pipeline moved TBs between two VPCs over a **Transit Gateway**. A single
transfer capped at ~5 Gbit/s despite 25 Gbit/s instances; the team suspected disk
or CPU. It was the **TGW per-flow limit (~5 Gbit/s)**. Splitting the transfer
into 8 parallel TCP streams (`s5cmd`/`rclone --transfers`) aggregated to
~25 Gbit/s across the same TGW. **Lesson:** cloud networks impose **per-flow**
ceilings well below link speed at NAT GW / TGW / VPN / single-flow paths; design
bulk movers to use many parallel connections (Ch. 36).

**Try it yourself**
```bash
# Which resolver / MTU / route is in play on an instance:
resolvectl status 2>/dev/null || cat /etc/resolv.conf
ip route ; ip route get 8.8.8.8 ; cat /sys/class/net/eth0/mtu
ping -M do -s 8972 <peer-in-vpc>        # jumbo within VPC?
ping -M do -s 1472 <internet-host>      # 1500 to the internet
# Single vs parallel flow (reveals per-flow caps):
iperf3 -c <peer> -t20 ; iperf3 -c <peer> -t20 -P8
# NAT GW connection pressure (from instances behind it): watch for EADDRNOTAVAIL / timeouts
ss -tan state established dst <external-service> | wc -l
# AWS: instance-level network allowance exhaustion
ethtool -S eth0 | grep -E 'bw_in_allowance_exceeded|bw_out_allowance_exceeded|pps_allowance_exceeded|conntrack_allowance_exceeded'
```

---


# Part XIX — HTTP in depth

Part XI covered HTTP's *shape* (versions, framing, multiplexing). This part covers
the semantics you actually operate: framing hazards, caching, cookies, CORS,
streaming, auth, and proxies.

## 58. Message framing: where does a request end?

### 58.1 The two ways to delimit a body

```
Content-Length: 1234              exact byte count; simple; needs the length up front
Transfer-Encoding: chunked        a series of  <hex-size>CRLF<data>CRLF  blocks,
                                  terminated by  0CRLF CRLF ; for streamed/unknown length
```

HTTP/1.1 rules (RFC 9112): if **both** are present, `Transfer-Encoding` wins and
`Content-Length` **must be ignored/rejected**. A message with neither, on a
response, is delimited by connection close (HTTP/1.0 style) — which is why
`Connection: close` responses can't be pipelined.

### 58.2 Request smuggling — the framing attack

If a front proxy and a back-end **disagree** about where a request ends, an
attacker can hide a second request inside the first.

```
Attacker sends ONE TCP stream containing:

  POST / HTTP/1.1
  Host: x
  Content-Length: 6
  Transfer-Encoding: chunked

  0

  GET /admin HTTP/1.1
  X: y

Front proxy honours Transfer-Encoding -> sees the request end at "0\r\n\r\n"
                                          and forwards the remainder as a NEW request
Back end honours Content-Length: 6      -> consumes 6 bytes, leaves "GET /admin"
                                          prefixed onto the NEXT victim's request
```

Variants are named by which side does what: **CL.TE**, **TE.CL**, **TE.TE**
(obfuscated `Transfer-Encoding` headers so one side ignores it).

**Defences**
- Normalise at the edge: reject any request containing **both** headers; reject
  unknown/obfuscated `Transfer-Encoding` values.
- Use **HTTP/2 end to end** (its framing is explicit and length-prefixed) — but
  beware **H2.CL / H2.TE downgrade smuggling** when the proxy translates H2 to
  HTTP/1.1 downstream; the proxy must regenerate framing, never trust the
  client-supplied length.
- Disable connection reuse to the back end if you can afford it (removes the
  victim's request from the picture).
- Reject requests with duplicate `Content-Length`, or with `Content-Length` on a
  `GET` where your app doesn't expect one.

### 58.3 Other framing edge cases

| Case | Behaviour |
|---|---|
| `HEAD` response | Has `Content-Length` but **no body**. Proxies must not wait for one. |
| `204 No Content`, `304 Not Modified` | **Never** have a body, regardless of headers. |
| `1xx` interim (`100 Continue`, `103 Early Hints`) | Followed by a *second* status line. Clients must handle multiple responses to one request. |
| `Expect: 100-continue` | Client asks permission before sending a large body; server replies `100` or an error. Saves uploading a body that will be rejected. |
| Trailers | Headers sent *after* a chunked body (`Trailer:` announces them). Used by gRPC for status. |

### 58.4 Real-world example

A CDN in front of a Node.js origin: the CDN parsed `Transfer-Encoding` with a
trailing space (`Transfer-Encoding : chunked`) as *invalid and therefore absent*,
while the origin's parser trimmed the space and honoured it. Classic TE.TE. One
crafted request poisoned the next user's response. The fix was a WAF rule
rejecting any header name with whitespace before the colon, plus upgrading both
parsers to RFC 9112-strict behaviour.

**Try it yourself**
```bash
# see chunked framing on the wire
curl -sv --raw http://httpbin.org/stream/3 2>&1 | head -30

# force HTTP/1.0 (no chunked; body ends at connection close)
curl -sv --http1.0 http://example.com 2>&1 | grep -i 'connection\|content-length'

# 100-continue in action
curl -sv -H 'Expect: 100-continue' -d @bigfile http://httpbin.org/post 2>&1 | grep '< HTTP'
```

---

## 59. Caching: the highest-leverage HTTP feature

### 59.1 The three caches in the path

```
browser cache  ->  (shared) CDN / proxy cache  ->  origin
   private            shared                       source of truth
```

**Private** caches hold one user's data; **shared** caches serve many users. The
distinction drives every directive below.

### 59.2 `Cache-Control` directives that matter

**Response directives:**

| Directive | Meaning |
|---|---|
| `max-age=N` | Fresh for N seconds (relative to `Date`) |
| `s-maxage=N` | Overrides `max-age` **for shared caches only** |
| `public` | May be stored by shared caches even if normally not (e.g. with `Authorization`) |
| `private` | Browser only — **CDNs must not store it** |
| `no-cache` | May store, but **must revalidate** before every reuse. (Not "don't cache".) |
| `no-store` | Never write to disk/memory. The real "don't cache". |
| `must-revalidate` | Once stale, do **not** serve stale — revalidate or fail |
| `immutable` | Content will never change; don't revalidate even on reload |
| `stale-while-revalidate=N` | Serve stale for up to N s while refreshing in the background |
| `stale-if-error=N` | Serve stale for up to N s if the origin errors |

**The single most common bug:** `no-cache` does *not* mean "don't cache". If you
mean that, say `no-store`.

### 59.3 Validators and conditional requests

```
Origin response:            Later request:                 Origin reply:
  ETag: "abc123"              If-None-Match: "abc123"   ->   304 Not Modified (no body)
  Last-Modified: <date>       If-Modified-Since: <date> ->   304 Not Modified
```

- **`ETag`** is an opaque version token. **Strong** (`"abc"`) means byte-identical;
  **weak** (`W/"abc"`) means semantically equivalent (fine for caching, not for
  range requests).
- A `304` carries headers but **no body** — that's the whole saving.
- `If-Match` / `If-Unmodified-Since` are the *write* direction: optimistic
  concurrency ("only update if it's still the version I read"), answering `412
  Precondition Failed` on conflict.

### 59.4 `Vary` — the correctness trap

`Vary` tells caches which **request** headers change the response.

```
Vary: Accept-Encoding        cache gzip and identity separately   (almost always needed)
Vary: Accept-Language        one entry per language
Vary: Cookie                 -> effectively uncacheable in a shared cache
Vary: User-Agent             -> cache explosion; avoid
```

Get it wrong in either direction and you either serve the wrong variant to
someone, or destroy your hit rate. `Vary: *` means "never reuse".

### 59.5 Cache keys and CDN behaviour

A CDN's cache key defaults to roughly `(method, scheme, host, path, query)` plus
whatever `Vary` demands. Operationally you often want to **normalise** it:

- Strip tracking params (`utm_*`, `fbclid`) from the key so they don't shard your
  cache.
- Ignore unknown query params for static assets.
- Add a header to the key deliberately (e.g. device class) rather than
  `Vary: User-Agent`.

**Cache-busting pattern:** serve immutable assets with a content hash in the
filename and `Cache-Control: public, max-age=31536000, immutable`; serve the HTML
that references them with `no-cache`. The HTML revalidates cheaply; assets never
do.

### 59.6 Invalidation

| Method | Notes |
|---|---|
| **TTL expiry** | Simplest; you wait |
| **Purge / invalidate API** | CDN-specific; may take seconds to propagate globally |
| **Surrogate keys / cache tags** | Tag responses (`Surrogate-Key: product-42`), purge by tag. The best pattern for dynamic sites. |
| **Versioned URLs** | No invalidation needed at all — the URL changes |

`Surrogate-Control` is a CDN-only sibling of `Cache-Control` (stripped before
reaching the browser), letting you cache aggressively at the edge and not at all
in the browser.

### 59.7 Real-world example

An e-commerce site had a 4% CDN hit rate. Cause: the origin sent
`Set-Cookie` on **every** response (a session cookie, even for anonymous static
pages), and the CDN — correctly — refuses to share a response carrying
`Set-Cookie`. Fix: only set the cookie on pages that need it, mark static assets
`public, max-age=31536000, immutable` with hashed filenames, and add
`Surrogate-Key` tags for product pages. Hit rate went to 94%; origin load fell
~20x.

**Try it yourself**
```bash
curl -sI https://example.com | grep -iE 'cache-control|etag|last-modified|age|vary'

# conditional request -- watch the 304
ET=$(curl -sI https://example.com | awk -F': ' '/[Ee]tag/{print $2}' | tr -d '\r')
curl -sI -H "If-None-Match: $ET" https://example.com | head -1

# see whether a CDN hit or missed (header name varies by vendor)
curl -sI https://www.cloudflare.com | grep -iE 'cf-cache-status|age|x-cache'
```

---

## 60. Cookies and state

### 60.1 The mechanics

```
Server -> client:   Set-Cookie: sid=abc; Path=/; Secure; HttpOnly; SameSite=Lax; Max-Age=3600
Client -> server:   Cookie: sid=abc            (on every matching request, automatically)
```

Cookies are attached by the **browser**, automatically, based on scope — which is
precisely why CSRF exists.

### 60.2 The attributes, and what they actually do

| Attribute | Effect |
|---|---|
| `Domain=example.com` | Sent to `example.com` **and all subdomains**. Omitting `Domain` is *narrower* (host-only) and usually safer. |
| `Path=/app` | Sent only for paths under `/app`. **Not a security boundary** (same-origin JS can read across paths). |
| `Secure` | HTTPS only. Always set it. |
| `HttpOnly` | Not readable by JavaScript — mitigates XSS token theft |
| `SameSite=Strict` | Never sent on cross-site requests, including top-level navigation. Breaks "click a link from email and be logged in". |
| `SameSite=Lax` | Sent on top-level **GET** navigation only. **The modern default.** |
| `SameSite=None` | Sent on all cross-site requests — **requires `Secure`**. Needed for third-party embeds/SSO. |
| `Max-Age` / `Expires` | Persistent; without either, it's a session cookie |
| `__Host-` prefix | Browser *enforces*: `Secure`, `Path=/`, and **no** `Domain`. The strongest binding available. |
| `__Secure-` prefix | Browser enforces `Secure`. |

### 60.3 Scoping hazards

- **Cookies ignore ports and (mostly) scheme.** A cookie set by
  `http://example.com:8080` is sent to `https://example.com:443`. Cookies are
  **not** origin-scoped like the rest of the web platform.
- **Subdomain leakage:** `Domain=example.com` sends your session cookie to
  `untrusted-user-content.example.com`. Use host-only cookies, and isolate
  user content on a **different registrable domain**, not a subdomain.
- **Cookie jar size limits** (~4 KB per cookie, ~50 per domain) — exceed them and
  cookies are silently dropped.
- **Every cookie is sent on every request** to that scope, including images and
  API calls. Large cookies are a real, measurable performance cost.

### 60.4 CSRF, and how `SameSite` changed it

Classic CSRF: an attacker's page issues a form POST to your bank; the browser
helpfully attaches the session cookie.

Defences, in order of strength:
1. **`SameSite=Lax` or `Strict`** — kills most classic CSRF. Now the browser
   default for cookies without an explicit attribute.
2. **Anti-CSRF token** — a random value in the form, checked server-side. Still
   required for cross-site flows and for `SameSite=None` cookies.
3. **Check `Origin` / `Sec-Fetch-Site` headers** server-side.
4. **Use `Authorization: Bearer`** instead of cookies — headers are not attached
   automatically, so CSRF doesn't apply (but you then own token storage, and
   `localStorage` is XSS-exposed).

### 60.5 Real-world example

A SaaS moved its API to `api.example.com` while the app stayed at
`app.example.com`, and set session cookies with `Domain=example.com` so both
would work. Six months later a marketing microsite at
`promo.example.com` — run by an agency, with a stored-XSS hole — leaked every
logged-in user's session cookie. The fix: host-only `__Host-` cookies plus a
token exchange for the API, and moving third-party content off the registrable
domain entirely.

**Try it yourself**
```bash
curl -sI https://example.com | grep -i set-cookie
# inspect flags in the browser: DevTools > Application > Cookies
# (columns for Secure, HttpOnly, SameSite are exactly the attributes above)
```

---

## 61. CORS: the browser's cross-origin rules

### 61.1 What problem it solves

The **same-origin policy** stops `evil.com`'s JavaScript from reading
`bank.com`'s responses. But legitimate cross-origin APIs exist. **CORS is the
server's way of granting exceptions** — enforced entirely by the browser.

> **CORS is not server security.** It restricts what *browser JavaScript* may
> read. `curl`, mobile apps, and servers ignore it completely. Never rely on CORS
> to protect data.

An **origin** is `scheme://host:port`. Different port or scheme = different
origin.

### 61.2 Simple vs preflighted requests

```
SIMPLE REQUEST (no preflight):
  method is GET, HEAD, or POST
  AND only "CORS-safelisted" headers
  AND Content-Type is one of:
      application/x-www-form-urlencoded, multipart/form-data, text/plain

  -> browser sends the request directly, with Origin: header
  -> then CHECKS the response headers before letting JS read it
     (the request WAS sent -- side effects already happened)

PREFLIGHTED REQUEST:
  anything else (PUT/DELETE/PATCH, Content-Type: application/json,
  custom headers like Authorization or X-Api-Key)

  -> browser first sends OPTIONS:
       OPTIONS /resource
       Origin: https://app.example.com
       Access-Control-Request-Method: PUT
       Access-Control-Request-Headers: content-type, authorization
  -> server must answer with matching Access-Control-Allow-* headers
  -> only then does the browser send the real request
```

**`Content-Type: application/json` is what triggers a preflight for most APIs.**
That's why "my GET works but my POST doesn't" is the archetypal CORS complaint.

### 61.3 The response headers

| Header | Purpose |
|---|---|
| `Access-Control-Allow-Origin` | The allowed origin, or `*`. **Cannot be `*` when credentials are used.** |
| `Access-Control-Allow-Methods` | Methods allowed (preflight response) |
| `Access-Control-Allow-Headers` | Request headers allowed (preflight response) |
| `Access-Control-Allow-Credentials: true` | Permit cookies / HTTP auth. Requires an explicit origin. |
| `Access-Control-Expose-Headers` | Which **response** headers JS may read (by default only a safelist) |
| `Access-Control-Max-Age` | How long the browser may cache the preflight (seconds) |

### 61.4 The four failures you will actually hit

1. **`Allow-Origin: *` with `credentials: 'include'`** → blocked. You must echo
   the specific origin *and* send `Allow-Credentials: true` *and* add
   `Vary: Origin` (or a shared cache will serve the wrong origin's header).
2. **Preflight not handled.** Your framework routes `PUT /x` but returns 404/405
   for `OPTIONS /x`. Handle `OPTIONS` explicitly.
3. **Custom header not listed.** You send `X-Request-Id`; the preflight response
   doesn't list it in `Allow-Headers`. Blocked.
4. **Reading a response header fails.** JS can only see a safelisted set unless
   you add `Access-Control-Expose-Headers`.

Also: **an error response still needs CORS headers.** A 500 without
`Access-Control-Allow-Origin` shows up in the browser as a CORS error, hiding the
real failure. Add the headers in your error handler too.

### 61.5 Real-world example

An API behind a CDN echoed the request `Origin` into
`Access-Control-Allow-Origin` — correct — but forgot `Vary: Origin`. The CDN
cached the response for `https://app-a.example` and served it, headers and all,
to `https://app-b.example`, which the browser then blocked. Intermittent, only
under cache hits, and invisible in origin logs. Fix: one header.

**Try it yourself**
```bash
# simulate a preflight
curl -sI -X OPTIONS https://api.example.com/v1/items \
  -H 'Origin: https://app.example.com' \
  -H 'Access-Control-Request-Method: PUT' \
  -H 'Access-Control-Request-Headers: content-type' | grep -i 'access-control\|vary'

# a simple request
curl -sI https://api.example.com/v1/items -H 'Origin: https://app.example.com' \
  | grep -i 'access-control'
```

---

## 62. WebSockets and Server-Sent Events

### 62.1 WebSocket: the upgrade

A WebSocket starts life as an HTTP request and then **stops being HTTP**.

```
CLIENT:
  GET /chat HTTP/1.1
  Host: example.com
  Upgrade: websocket
  Connection: Upgrade
  Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
  Sec-WebSocket-Version: 13
  Sec-WebSocket-Protocol: chat            (optional subprotocol negotiation)
  Origin: https://example.com

SERVER:
  HTTP/1.1 101 Switching Protocols
  Upgrade: websocket
  Connection: Upgrade
  Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=      (SHA-1 of key + magic GUID)

-> the same TCP connection now carries WebSocket FRAMES, bidirectionally,
   until either side sends a Close frame.
```

The `Sec-WebSocket-Accept` handshake exists to prove the server *understood* the
upgrade (defeating cache-poisoning attacks against naive proxies), not for
security.

### 62.2 Frames and operational details

| Aspect | Detail |
|---|---|
| **Framing** | Opcodes: text(1), binary(2), close(8), ping(9), pong(10). Continuation frames for fragmentation. |
| **Masking** | Client-to-server frames **must** be XOR-masked with a random key (an anti-cache-poisoning measure). Server-to-client must not be. |
| **Ping/pong** | The protocol's keepalive. **Use it** — see below. |
| **Close** | A close frame carries a status code (1000 normal, 1001 going away, 1006 = abnormal/no close frame seen) |
| **Subprotocols** | `Sec-WebSocket-Protocol` negotiates an application protocol (e.g. `graphql-ws`) |
| **Compression** | `permessage-deflate` extension; watch CPU and CRIME-style risks with mixed secret/attacker data |
| **Over HTTP/2** | RFC 8441 `:protocol` extended CONNECT; over HTTP/3 similarly. Many stacks still downgrade to HTTP/1.1 for WS. |

### 62.3 The operational hazards

- **Idle timeouts kill you.** Every proxy, load balancer, and NAT has one
  (commonly 60 s, AWS ALB default 60 s, NLB 350 s). A WebSocket with no traffic
  is silently reaped. **Send a ping every 20–30 s**, from the server, and expect
  pongs.
- **`1006` means "connection died without a close frame"** — almost always a
  middlebox timeout, a network drop, or a crashed process. It carries no reason,
  so log your own context.
- **Load balancers must be configured for it.** L7 LBs need WebSocket support
  enabled; sticky routing matters because the connection is long-lived and
  stateful.
- **Scaling is stateful.** N servers × M connections; broadcasting requires a
  pub/sub bus (Redis, NATS) because a message must reach whichever server holds
  each connection.
- **Reconnect with backoff and jitter**, and resume with a cursor — a
  mass-disconnect otherwise produces a thundering herd.

### 62.4 Server-Sent Events: the simpler option

SSE is **plain HTTP** with a never-ending response body.

```
Request:   GET /events    Accept: text/event-stream

Response:  HTTP/1.1 200 OK
           Content-Type: text/event-stream
           Cache-Control: no-cache
           Connection: keep-alive

           event: price
           id: 42
           data: {"symbol":"ACME","px":10.5}

           : this is a comment, used as a heartbeat

           data: another message
```

| | WebSocket | SSE |
|---|---|---|
| Direction | **Bidirectional** | Server → client only |
| Protocol | Custom framing after upgrade | Plain HTTP |
| Reconnect | You implement it | **Automatic**, built into `EventSource` |
| Event IDs / resume | You implement it | `id:` + `Last-Event-ID` header, built in |
| Proxy/CDN friendliness | Needs explicit support | Works through anything that streams |
| Binary | Yes | Text only (base64 if you must) |
| Browser connection limit | Fine | 6 per origin on HTTP/1.1 (**use HTTP/2**) |

**Choose SSE** for one-way feeds (notifications, dashboards, LLM token streaming,
progress). **Choose WebSocket** when the client must push frequently (chat,
collaborative editing, games).

Also consider **long polling** as a fallback, and note that plain HTTP **chunked
streaming** (what `fetch` + `ReadableStream` consumes) is often enough.

### 62.5 Real-world example

A trading dashboard used WebSockets behind an ALB with the default 60-second idle
timeout and no application ping. During quiet markets, connections dropped every
minute with code `1006`, the client reconnected, and the reconnect storm at the
open took down the service. Fix: server-side ping every 25 s, ALB idle timeout
raised to 300 s, exponential backoff with jitter on reconnect, and `Last-Event-ID`-
style resume. Disconnects fell to near zero.

**Try it yourself**
```bash
# watch the upgrade handshake
curl -sv -N -H 'Connection: Upgrade' -H 'Upgrade: websocket' \
     -H 'Sec-WebSocket-Version: 13' -H 'Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==' \
     https://echo.websocket.org/ 2>&1 | grep -iE '^< HTTP|upgrade|sec-websocket'

# consume an SSE stream (note it never ends)
curl -N -H 'Accept: text/event-stream' https://stream.wikimedia.org/v2/stream/recentchange | head -20

# websocat is the netcat of websockets
websocat wss://echo.websocket.org
```

---

## 63. Range requests, streaming, and compression

### 63.1 Range requests

```
Server advertises:   Accept-Ranges: bytes
Client asks:         Range: bytes=1000-1999
Server replies:      206 Partial Content
                     Content-Range: bytes 1000-1999/50000
                     Content-Length: 1000
```

Uses: resumable downloads, video seeking, parallel segment fetching, reading a
file header without the whole file.

| Detail | Behaviour |
|---|---|
| `Range: bytes=1000-` | From byte 1000 to the end |
| `Range: bytes=-500` | The **last** 500 bytes |
| Multiple ranges | `multipart/byteranges` response — rarely worth supporting |
| Unsatisfiable range | `416 Range Not Satisfiable` |
| **`If-Range: <etag or date>`** | "Send the range only if the resource is unchanged; otherwise send the whole thing (200)." **Essential for correct resumption** — without it a resumed download can splice two different versions together. |
| Weak ETags | Must **not** be used with `If-Range` (they don't guarantee byte-identity) |

### 63.2 Streaming responses

Three mechanisms, often confused:

```
Chunked transfer-encoding   HTTP/1.1 body of unknown length; each chunk usable on arrival
HTTP/2 + /3 DATA frames     same effect, native to the protocol
Server-Sent Events          a chunked stream with an event format on top (Ch 62)
```

Operational gotchas: proxies and CDNs may **buffer** your stream (nginx
`proxy_buffering off`, `X-Accel-Buffering: no`), compression can add latency by
buffering, and many frameworks buffer the whole response by default.

### 63.3 Compression

```
Client:  Accept-Encoding: br, gzip, deflate, zstd
Server:  Content-Encoding: br
         Vary: Accept-Encoding                <-- MANDATORY, or caches serve the wrong body
```

| Codec | Notes |
|---|---|
| **gzip** | Universal; good ratio; cheap |
| **brotli (`br`)** | ~15–20% better than gzip on text; higher CPU at high levels — **pre-compress static assets at level 11**, use level 4–5 for dynamic |
| **zstd** | Fast, increasingly supported |
| **deflate** | Historically ambiguous implementations; avoid |

Rules:
- **Never compress already-compressed formats** (JPEG, PNG, MP4, zip) — you burn
  CPU for nothing.
- **Set a minimum size** (~1 KB); compressing tiny responses can make them bigger.
- **`Content-Encoding` vs `Transfer-Encoding`**: the former is a property of the
  *resource* (and is what `ETag` covers); the latter is hop-by-hop.
- **Compression + secrets = BREACH.** If a response contains both attacker-
  controlled input and a secret (like a CSRF token), compression ratio leaks the
  secret. Mitigate by not reflecting user input alongside secrets, masking tokens
  per-response, or disabling compression on sensitive endpoints.

### 63.4 Real-world example

A media API supported `Range` but not `If-Range`. A CDN purge mid-download meant a
client's resumed request returned bytes from a *newer* file version, producing
corrupted downloads for ~0.3% of users — reported as "random file corruption" for
months. Adding `If-Range` support (falling back to a full 200 on mismatch) fixed
it outright.

**Try it yourself**
```bash
curl -sI https://example.com | grep -i accept-ranges
curl -s -r 0-99 https://example.com -o /dev/null -w '%{http_code} %{size_download}\n'
curl -sI -H 'Range: bytes=0-99' https://example.com | grep -i 'content-range\|^HTTP'

# see compression actually applied
curl -sI -H 'Accept-Encoding: gzip, br' https://example.com | grep -i 'content-encoding\|vary'
# compare sizes
curl -s https://example.com | wc -c
curl -s -H 'Accept-Encoding: gzip' --compressed-no-decompress https://example.com 2>/dev/null | wc -c
```

---

## 64. HTTP authentication and authorization on the wire

### 64.1 The mechanisms

| Scheme | Wire format | Notes |
|---|---|---|
| **Basic** | `Authorization: Basic base64(user:pass)` | **Not encryption** — base64 is reversible. HTTPS only. Fine for internal/service use. |
| **Digest** | Challenge-response with nonces | Avoids sending the password, but weak by modern standards and rarely used |
| **Bearer** | `Authorization: Bearer <token>` | The dominant API pattern (OAuth 2.0, RFC 6750) |
| **API key** | Usually a custom header (`X-Api-Key`) or query param | Query params leak into logs, `Referer`, and browser history — **use a header** |
| **mTLS** | Client certificate at the TLS layer (Ch 38.10) | Strongest; identity established before HTTP starts |
| **Cookie session** | `Cookie:` (Ch 60) | Browser-native; brings CSRF concerns |
| **HMAC request signing** | Signature over method + path + body + timestamp (AWS SigV4 style) | Resists replay and tampering; no bearer token to steal |

### 64.2 The 401 challenge flow

```
Client:  GET /private
Server:  401 Unauthorized
         WWW-Authenticate: Bearer realm="api", error="invalid_token",
                           error_description="expired"
Client:  GET /private
         Authorization: Bearer eyJ...
Server:  200 OK
```

**`401` vs `403`:** `401` = "I don't know who you are (or your credentials are
bad) — authenticate and retry." `403` = "I know who you are and you may not do
this — retrying won't help." Sending `401` for an authorization failure causes
clients to loop through token refresh pointlessly.

`407 Proxy Authentication Required` is the same flow for a proxy, using
`Proxy-Authenticate` / `Proxy-Authorization`.

### 64.3 OAuth 2.0 / OIDC at the wire level

You don't need the whole spec, but you should recognise the shapes:

```
Authorization Code + PKCE  (the correct flow for web and mobile apps)
  1. browser -> /authorize?response_type=code&client_id=..&redirect_uri=..
                &scope=..&state=..&code_challenge=..&code_challenge_method=S256
  2. user authenticates at the identity provider
  3. redirect back with ?code=..&state=..        <- state defends against CSRF
  4. server-side: POST /token with code + code_verifier  -> access_token,
                                                            refresh_token, id_token
  5. call the API with  Authorization: Bearer <access_token>

Client Credentials  (machine-to-machine)
  POST /token  grant_type=client_credentials  -> access_token
```

Operational notes:
- **Access tokens should be short-lived** (minutes); refresh tokens long-lived and
  rotated on use.
- **JWTs are not encrypted** — they're signed. Anyone can read the payload
  (`base64url` decode it). Never put secrets in one.
- **Validate**: signature, `iss`, `aud`, `exp`, `nbf`, and algorithm (reject
  `alg: none` and algorithm confusion).
- **Revocation is the JWT weakness** — a stateless token is valid until it
  expires. Keep lifetimes short, or maintain a denylist.

### 64.4 Real-world example

An internal API used API keys passed as `?api_key=...`. The keys turned up in: the
CDN's access logs, the load balancer's logs, the `Referer` header sent to a
third-party analytics script on an error page, and a screenshot in a support
ticket. Rotating them took a week. Moving to `Authorization: Bearer` plus
short-lived tokens removed four of those exposure paths immediately.

**Try it yourself**
```bash
curl -sI https://httpbin.org/basic-auth/user/pass | head -3      # 401 + WWW-Authenticate
curl -s  -u user:pass https://httpbin.org/basic-auth/user/pass   # 200

# decode a JWT payload (it is NOT encrypted)
echo '<jwt>' | cut -d. -f2 | tr '_-' '/+' | base64 -d 2>/dev/null | python3 -m json.tool
```

---

## 65. Proxies: forward, reverse, CONNECT, and interception

### 65.1 The two directions

```
FORWARD PROXY (on the client's side; the client knows about it)
   client -> [proxy] -> any server on the internet
   uses: corporate egress control, filtering, caching, anonymity

REVERSE PROXY (on the server's side; the client does not know)
   client -> [proxy] -> your back-end servers
   uses: load balancing, TLS termination, caching, WAF   (Ch 40)
```

Part XI covered reverse proxies. This section is about the forward direction,
which developers meet as "why does this only fail on the corporate network?".

### 65.2 The `CONNECT` method

A forward proxy can't read HTTPS, so for TLS it builds a **blind TCP tunnel**:

```
Client -> proxy:   CONNECT example.com:443 HTTP/1.1
                   Host: example.com:443
                   Proxy-Authorization: Basic ...        (if required)
Proxy  -> client:  HTTP/1.1 200 Connection Established

-> from here the proxy blindly relays bytes; the TLS handshake happens
   end-to-end between client and origin. The proxy sees only the hostname
   (from CONNECT) and the traffic volume.
```

Consequences: the proxy can **allow/deny by hostname** but cannot inspect
content — unless it performs TLS interception (below).

### 65.3 How clients find a proxy

| Mechanism | How |
|---|---|
| Environment variables | `http_proxy`, `https_proxy`, `no_proxy` (respected by curl, most CLIs, many SDKs) |
| **PAC file** | A JavaScript function `FindProxyForURL(url, host)` returning `"PROXY host:port"` / `"DIRECT"`. Allows per-destination rules. |
| **WPAD** | Auto-discovery of the PAC file via DHCP option 252 or DNS. Convenient; historically a nice attack vector. |
| OS/browser settings | System proxy configuration |

**`no_proxy` is a classic footgun:** its syntax is inconsistent across tools
(leading dots, CIDR support, wildcard support all vary), so "it works in curl but
not in my app" is common.

### 65.4 TLS interception (corporate MITM)

To inspect HTTPS, an enterprise proxy terminates TLS and re-encrypts:

```
client ==TLS(corp CA cert)==> [proxy decrypts, inspects] ==TLS(real cert)==> origin
```

This requires installing the **corporate root CA** on every device. Effects you
will hit:

- **Certificate pinning breaks** (apps that pin to a specific CA/key refuse to
  connect).
- **Non-browser clients fail** because they use their own trust store: Python
  (`certifi`), Node, Java (`cacerts`), Go, Docker, package managers. Each needs
  the corporate CA added separately — this is the single most common cause of
  "SSL certificate problem: unable to get local issuer certificate" in corporate
  environments.
- **mTLS to external services breaks** (the proxy can't present your client
  cert).
- **HTTP/3 usually isn't supported** by the proxy, so QUIC gets blocked and
  clients fall back to TCP.

### 65.5 The `Forwarded` / `X-Forwarded-*` headers

Once a proxy is in the path, the back end sees the *proxy's* address. Identity is
carried in headers:

```
X-Forwarded-For: 203.0.113.9, 198.51.100.7    (client, then each proxy, left to right)
X-Forwarded-Proto: https
X-Forwarded-Host: example.com
Forwarded: for=203.0.113.9;proto=https;host=example.com    (RFC 7239, the standard one)
```

**Security rule:** these headers are **client-supplied and trivially spoofed**.
Only trust them from known proxy IPs, and *strip or overwrite* them at your edge.
Count from the right, using the number of proxies you actually operate. Getting
this wrong turns your rate limiter and IP allowlist into decoration.

For TCP-mode (L4) proxies that can't add headers, the **PROXY protocol**
(Ch 38.11) carries the original addresses ahead of the byte stream.

### 65.6 Real-world example

A team's CI pipeline failed only inside the corporate network with
`x509: certificate signed by unknown authority`. `curl` worked (it had been
configured with `--cacert`), the browser worked (OS trust store had the corporate
CA), but Go and Python did not. Root cause: TLS interception plus three separate
trust stores. Fix: bake the corporate CA into the CI container image and set
`SSL_CERT_FILE`/`NODE_EXTRA_CA_CERTS`/`REQUESTS_CA_BUNDLE`, plus a `no_proxy`
entry for internal registries.

**Try it yourself**
```bash
# use a proxy explicitly
export https_proxy=http://proxy.corp:3128
curl -v https://example.com 2>&1 | grep -iE 'CONNECT|Proxy|HTTP/1.1 200'

# what does my environment think?
env | grep -i proxy

# which CA store is a tool using?
curl -v https://example.com 2>&1 | grep -i 'CAfile\|CApath'
python3 -c "import certifi; print(certifi.where())"

# detect TLS interception: is the issuer your corporate CA?
echo | openssl s_client -connect example.com:443 2>/dev/null \
  | openssl x509 -noout -issuer
```

---
# Part XX — Monitoring and telemetry

Parts XII–XIII covered *investigating* a problem you already know about. This part
is about knowing first — what to measure continuously, with what mechanism, and
how to turn it into alerts that aren't noise.

## 66. What to measure: the four signal families

### 66.1 The taxonomy

```
+-------------------+---------------------------+--------------------------+
| FAMILY            | ANSWERS                   | MECHANISM                |
+-------------------+---------------------------+--------------------------+
| 1. DEVICE/LINK    | "is the pipe healthy?"    | SNMP, kernel counters,   |
|    metrics        | utilisation, errors,      | ethtool, /proc, gNMI     |
|                   | drops, temperature        |                          |
+-------------------+---------------------------+--------------------------+
| 2. FLOW records   | "who talked to whom,      | NetFlow, IPFIX, sFlow,   |
|                   | how much, over what?"     | conntrack, eBPF          |
+-------------------+---------------------------+--------------------------+
| 3. ACTIVE probes  | "what does a user         | ping, TWAMP, synthetic   |
|    (synthetic)    | actually experience?"     | HTTP checks, traceroute  |
+-------------------+---------------------------+--------------------------+
| 4. PACKETS        | "what EXACTLY happened?"  | tcpdump, mirrors/TAPs,   |
|                   |                           | continuous PCAP          |
+-------------------+---------------------------+--------------------------+
```

Cost and volume rise sharply down the list; specificity rises with it. A mature
setup runs **1 and 3 always**, **2 sampled always**, and **4 on demand or in a
short ring buffer**.

### 66.2 The signals that actually catch incidents

**Per interface (from Chapter 48's drop table):**

| Metric | Why | Alert on |
|---|---|---|
| `ifInOctets` / `ifOutOctets` | Utilisation | > 70–80% sustained (queueing starts well before 100%) |
| `ifInErrors` / CRC errors | Physical layer failing | **any sustained increase** |
| `ifInDiscards` / `ifOutDiscards` | Buffer exhaustion / policing | rate of change |
| Interface flaps (`ifOperStatus`) | Bad optic/cable/power | > 1 flap in an hour |
| Optical Tx/Rx power (DOM/DDM) | Failing SFP or dirty fibre | outside vendor range, or trending down |
| Queue drops per class | QoS starvation (Ch. 51) | any drop in the priority queue |

**Per host (Linux, Ch. 43/48):**

```
nstat / /proc/net/snmp   TcpRetransSegs / TcpOutSegs      -> loss proxy
                         TcpExtListenOverflows/Drops       -> accept queue full  (Ch 24)
                         TcpExtTCPBacklogDrop              -> socket backlog full
                         UdpRcvbufErrors                   -> UDP receive drops   (Ch 22)
                         IpReasmFails                      -> fragment problems   (Ch 15)
/proc/net/softnet_stat   col2 drops, col3 time_squeeze     -> softirq overload    (Ch 48)
conntrack -S             insert_failed, drop               -> state table full    (Ch 53)
ss -ti                   rtt, retrans, cwnd per socket     -> per-connection health
```

**Per service (the ones users feel):**
`TTFT`/`TTFB`, request rate, error rate (4xx vs 5xx separately), p50/p95/p99
latency, connection setup time, TLS handshake time, DNS resolution time.

### 66.3 RED, USE, and the golden signals

Three complementary frames, all useful:

| Frame | Measures | Best for |
|---|---|---|
| **USE** (Brendan Gregg) | **U**tilisation, **S**aturation, **E**rrors — per *resource* | Links, NICs, queues, CPUs |
| **RED** | **R**ate, **E**rrors, **D**uration — per *service* | Request-driven services |
| **Four golden signals** (Google SRE) | Latency, Traffic, Errors, Saturation | Anything user-facing |

**Saturation is the one people forget**, and it's the one that predicts trouble.
A link at 60% utilisation with a full queue is already hurting; utilisation alone
hides it. Measure queue depth and drops, not just throughput.

### 66.4 Percentiles, not averages

```
Average latency 50 ms sounds fine.
   p50 = 20 ms, p99 = 4,000 ms  -> 1% of users are having a terrible time.

Rules:
  * always chart p50, p95, p99 (and p99.9 if you have the volume)
  * NEVER average percentiles across instances or time buckets -- it is
    mathematically meaningless. Aggregate from histograms instead.
  * plot MAX for saturation metrics; a 5-minute average hides a 10-second stall
```

Round-number latency clusters are timers, not physics — 40 ms (Nagle/delayed ACK,
Ch. 33), 200 ms (RTO, Ch. 28), 1/3/7 s (SYN retries, Ch. 24), 5 s (dead DNS
server, Ch. 20).

### 66.5 Real-world example

A payments service alerted only on CPU and HTTP 5xx rate. Every few days it saw a
30-second spike of connection failures that never triggered anything, because the
failures happened *before* the request reached the application — the accept queue
was overflowing (`TcpExtListenOverflows`) during traffic bursts. Nothing in the
app logs, nothing in CPU. Adding `nstat` counters to the metrics pipeline
surfaced it in a day; the fix was a bigger `listen()` backlog. **Connection-level
failures are invisible to application logs — you must scrape kernel counters.**

**Try it yourself**
```bash
nstat -az | grep -Ei 'ListenOverflows|ListenDrops|TCPBacklogDrop|RetransSegs|OutSegs'
cat /proc/net/softnet_stat | head -4     # col2=drops col3=time_squeeze (hex)
ip -s link show eth0                      # RX/TX errors, dropped
ethtool -S eth0 | grep -Ei 'err|drop|crc|miss'
ss -ti state established | grep -oE 'rtt:[0-9.]+|retrans:[0-9]+/[0-9]+' | head
```

---

## 67. Device and interface metrics: SNMP and its successors

### 67.1 SNMP in one page

**Simple Network Management Protocol** — the lingua franca of device monitoring
since 1988. UDP/161 (queries), UDP/162 (traps).

```
MANAGER (your poller)                       AGENT (switch/router/server)
     |                                              |
     |  GET / GETNEXT / GETBULK  (UDP 161)  ------->|
     |<--------------------------- response          |
     |                                              |
     |<-- TRAP / INFORM (UDP 162, unsolicited) ------|   "my interface just went down"
```

**OIDs and MIBs.** Every value has a numeric **OID** (`1.3.6.1.2.1.2.2.1.10.3`),
and a **MIB** is the dictionary mapping names to OIDs
(`IF-MIB::ifInOctets.3`). The ones you'll use constantly:

| Name | OID root | Value |
|---|---|---|
| `sysUpTime` | `1.3.6.1.2.1.1.3` | Time since boot — **a reset means the device rebooted** |
| `ifDescr` / `ifName` / `ifAlias` | `1.3.6.1.2.1.2.2.1.2` / `31.1.1.1.1` / `31.1.1.1.18` | Interface naming; `ifAlias` is the human description |
| `ifHCInOctets` / `ifHCOutOctets` | `1.3.6.1.2.1.31.1.1.1.6/.10` | **64-bit** byte counters — use these, not the 32-bit `ifInOctets` |
| `ifInErrors` / `ifOutErrors` | `...2.2.1.14/.20` | Errors |
| `ifInDiscards` / `ifOutDiscards` | `...2.2.1.13/.19` | Drops (buffer/policy) |
| `ifOperStatus` / `ifAdminStatus` | `...2.2.1.8/.7` | Up/down, and whether a human shut it |

**The 32-bit counter trap:** a 32-bit octet counter wraps at ~4.3 GB. On a
10 Gbit/s link that is **under 4 seconds**. Always poll the 64-bit `ifHC*`
counters (SNMPv2c or later), and always compute rates as deltas with wrap
handling.

**Versions:** v1 (obsolete), **v2c** (community string in cleartext — treat the
community as a password sent in the clear), **v3** (authentication + encryption;
use it). Restrict SNMP to a management network regardless.

### 67.2 Polling model and its limits

```
poll interval 60 s  ->  you cannot see a 10-second microburst
                        (it averages into invisibility)
poll 5,000 devices x 50 interfaces x 8 metrics every 60 s
                     ->  33,000 SNMP GETs/second -- and SNMP is a
                         request/response protocol over UDP
```

Consequences: use **GETBULK**, spread polls, and accept that SNMP is a *minutes*
resolution tool. For sub-second visibility you need streaming telemetry.

### 67.3 Streaming telemetry (gNMI / NETCONF / OpenConfig)

The modern replacement for polling: the device **pushes** updates on subscription.

```
SNMP (pull):    manager asks every 60 s          text/ASN.1, UDP, per-OID
gNMI (push):    device streams on change or at   protobuf over gRPC/HTTP2/TLS,
                a chosen cadence (e.g. 1 s)      structured YANG paths
```

- **YANG** is the data model; **OpenConfig** is a vendor-neutral set of YANG
  models. **gNMI** is the transport (`Subscribe` with `ON_CHANGE` or
  `SAMPLE` modes).
- Advantages: sub-second resolution, event-driven (no missed transitions),
  efficient encoding, structured data, secure by default.
- Reality: vendor support and model coverage vary; SNMP remains the universal
  fallback.

### 67.4 Host-level exporters

For servers, skip SNMP:

| Tool | Provides |
|---|---|
| **node_exporter** (Prometheus) | `/proc/net/*`, interface counters, conntrack, sockets |
| **cAdvisor / kube-state-metrics** | container and pod network stats |
| **eBPF exporters** (`ebpf_exporter`, Cilium/Hubble, Pixie) | per-socket RTT, retransmits, TCP state transitions, DNS latency — **without packet capture** |

eBPF is the significant recent shift: you can measure per-connection RTT and
retransmits for *every* flow at low cost, which used to require sampling packets.

### 67.5 Real-world example

A 40 Gbit/s uplink showed 45% average utilisation and constant complaints of
"slow at 09:00". Five-minute SNMP averages hid **10-second microbursts** hitting
line rate every morning as thousands of clients synced. Switching that interface
to 1-second streaming telemetry made the bursts obvious; the fix was shaping and
staggering the sync window. **Averaging is how you hide the incident from
yourself.**

**Try it yourself**
```bash
# poll a device you own (v2c shown; prefer v3 in production)
snmpwalk  -v2c -c public 192.168.1.1 1.3.6.1.2.1.1          # system info
snmpwalk  -v2c -c public 192.168.1.1 IF-MIB::ifDescr
snmpget   -v2c -c public 192.168.1.1 IF-MIB::ifHCInOctets.3
snmpbulkwalk -v2c -c public 192.168.1.1 IF-MIB::ifTable      # efficient bulk read

# rate = (counter2 - counter1) / (t2 - t1) * 8  -> bits per second
```

---

## 68. Flow telemetry: NetFlow, IPFIX, and sFlow

### 68.1 What a flow record is

A **flow** is a unidirectional sequence of packets sharing a key — classically the
5-tuple (src IP, dst IP, src port, dst port, protocol) plus input interface and
ToS.

The device summarises each flow and exports a record:

```
start=10:00:01.204  end=10:00:04.881  duration=3.677s
src=192.168.1.50:51000  ->  dst=93.184.216.34:443   proto=TCP
packets=48  bytes=52104  tcp_flags=0x1b  tos=0  in_if=3  out_if=1
src_as=64512  dst_as=15133  next_hop=10.0.0.1
```

**This is the answer to "who is using my bandwidth?"** — a question metrics
cannot answer and packet capture answers too expensively.

### 68.2 The three technologies

| | **NetFlow v5** | **NetFlow v9 / IPFIX** | **sFlow** |
|---|---|---|---|
| Origin | Cisco, 1996 | Cisco v9 → IETF **IPFIX** (RFC 7011) | InMon, standard |
| Format | Fixed fields | **Template-based** (extensible: IPv6, MPLS, MAC, VLAN, URLs) | Fixed structures |
| Method | Device tracks every flow in a cache, exports on expiry | Same | **Packet sampling** (1-in-N) + interface counters |
| Accuracy | Exact byte counts (if not sampled) | Exact, or sampled | Statistical — accurate in aggregate, not per-flow |
| Cost on device | High (flow cache state) | High | **Low** (no state) |
| Latency to see data | Seconds to minutes (cache timers) | Same | **Near real-time** |
| IPv6 / MPLS | v5 cannot | Yes | Yes |

**Rule of thumb:** IPFIX/NetFlow v9 when you need accurate accounting (billing,
security forensics); sFlow when you need cheap, real-time, high-speed visibility
across many devices. Many networks run both.

### 68.3 The timers that confuse everyone

Flow records don't appear when the flow starts:

```
active timeout   (default often 60 s or 1-5 min)
    a long-lived flow is exported PERIODICALLY, in pieces
inactive/idle timeout  (default ~15-30 s)
    a finished flow is exported this long after the last packet
```

So a 10-minute download appears as several records, and a short flow appears
~15 s late. **Never build a real-time alert on flow data expecting sub-minute
freshness**, and always sum the pieces when accounting.

### 68.4 Sampling maths

At 1-in-1000 sampling, a 5,000-packet flow yields ~5 samples. Scale back up by
the sampling rate to estimate bytes. Implications:

- Aggregate totals are accurate; **individual small flows are invisible**.
- Security use (finding one connection to a C2 server) needs unsampled flow or
  packet data.
- Always record the sampling rate with the data — mixing rates silently corrupts
  totals.

### 68.5 Where to collect flows when you don't own routers

| Environment | Source |
|---|---|
| Cloud | **VPC Flow Logs** (AWS/GCP/Azure) — same idea, delivered to object storage or a log pipeline |
| Linux host | `conntrack` events, eBPF (Cilium **Hubble**, Pixie), `nfcapd` with a software exporter |
| Kubernetes | Hubble, Calico flow logs, service-mesh telemetry |
| Anywhere | `nfdump`/`nfsen`, `pmacct`, `GoFlow2`, `vflow`, or ntopng |

**VPC flow logs are the cloud-native equivalent** and answer exactly the same
questions — including the crucial `ACCEPT`/`REJECT` field, which tells you whether
a security group dropped the traffic.

### 68.6 What flow data is actually good for

1. **Top talkers / bandwidth attribution** — which host, app, or customer.
2. **Capacity planning** — traffic matrices between sites over months.
3. **Peering and transit decisions** — traffic per destination AS, which is why
   `src_as`/`dst_as` fields exist.
4. **DDoS detection** — sudden explosion in flows-per-second with tiny packet
   counts is a scan or SYN flood signature.
5. **Security forensics** — "did anything talk to this IP in the last 90 days?"
   Cheap to retain for months; packet capture is not.
6. **Verifying firewall rules** — cloud flow logs' REJECT records show what's
   being blocked and whether that's intended.

### 68.7 Real-world example

A university's 10 Gbit/s uplink saturated every evening. Interface metrics showed
only "100% utilised". IPFIX records showed 62% of traffic was a handful of hosts
talking to a single destination AS on high ports — unauthorised media servers in a
residence hall. Without flow data the only options were guesswork or a full packet
capture at line rate. With it: fifteen minutes to the answer, and a QoS policy
(Ch. 51) rather than an uplink upgrade.

**Try it yourself**
```bash
# software flow export/collection on a Linux box
sudo apt install nfdump softflowd
sudo softflowd -i eth0 -n 127.0.0.1:9995 -v 9
nfcapd -w -D -l /tmp/flows -p 9995 &
sleep 120
nfdump -R /tmp/flows -s srcip/bytes -n 10      # top talkers
nfdump -R /tmp/flows -s dstport/bytes -n 10    # top destination ports

# cloud equivalent
aws ec2 describe-flow-logs
# then query the logs in Athena/CloudWatch Insights, grouping by srcaddr/dstaddr/bytes
```

---

## 69. Active probing, synthetic monitoring, and alerting that works

### 69.1 Passive vs active

```
PASSIVE  (Ch 66-68)  observes real traffic
   + reflects real users      - only sees paths users actually take
                              - silent when there is no traffic (3am outage
                                goes unnoticed until morning)

ACTIVE   (this chapter)  generates synthetic traffic
   + constant baseline        - synthetic; may not match real user experience
   + tests paths before       - adds (small) load
     users do
   + detects "up but broken"
```

**You need both.** Active monitoring is what tells you a service is down at 3am;
passive tells you what real users are experiencing.

### 69.2 Layered synthetic checks

Design probes to isolate the layer that failed — mirroring Ch. 34's method:

```
L3   ping / ICMP           reachability + RTT + loss
L3   traceroute/mtr        path and where loss begins
L4   TCP connect to :443   is the port open? handshake time?
L6   TLS handshake         certificate valid? expiry? handshake time?
L7   HTTP GET /healthz     status code, latency, body content check
L7   full transaction      log in, search, add to cart, check out
DNS  resolve the name      from multiple resolvers and regions
```

A check at each layer means the alert *names the layer*. A single "site is down"
check tells you nothing about where.

### 69.3 Standardised network probing

| Protocol | Purpose |
|---|---|
| **ICMP echo** | Universal, but deprioritised by routers (Ch. 14) — never measure a router's *own* response as if it were transit performance |
| **TWAMP / OWAMP** (RFC 5357 / 4656) | Purpose-built two-way / one-way active measurement with timestamps; measures **one-way** delay and jitter separately, which ICMP cannot |
| **IP SLA / RPM** (vendor) | Router-embedded probes (ICMP, UDP jitter, HTTP, DNS) reporting into SNMP |
| **iperf3 / netperf** | On-demand throughput testing (not for continuous monitoring — it saturates the link) |
| **RIPE Atlas / probes** | Global measurement from thousands of external vantage points; free to use |

**One-way delay matters** because paths are asymmetric (Ch. 13). A 90 ms round
trip could be 20 ms out and 70 ms back — TWAMP shows that; ping cannot.

### 69.4 Where to probe from

```
inside the datacentre     -> tests the service, not the network to users
from each POP/region      -> catches regional routing/peering problems
from real user networks   -> external synthetic providers, RIPE Atlas,
                             or RUM (Real User Monitoring) in the page
```

**RUM** (real-user timing beaconed from browsers, via the Navigation Timing and
Resource Timing APIs) is the ground truth for user experience; synthetics are the
controlled, comparable baseline. Use synthetics for alerting and RUM for
prioritisation.

### 69.5 Alerting that isn't noise

```
ALERT ON SYMPTOMS, NOT CAUSES
   good:  "p99 latency > 2 s for 5 minutes"   (users are suffering)
   good:  "error rate > 1% for 3 minutes"
   poor:  "CPU > 80%"                          (so what? is anything broken?)
   poor:  "interface at 75%"                   (page a human at 3am for this?)

MAKE EVERY PAGE ACTIONABLE
   if the responder cannot do something about it right now, it is not a page.
   Everything else is a ticket or a dashboard.

USE MULTI-WINDOW BURN RATES (SRE style) rather than raw thresholds:
   page   if  error budget burns at 14.4x over 1 h  AND  6x over 5 min
   ticket if  it burns at 1x over 6 h
   -> fast burns page immediately; slow burns don't wake anyone

SUPPRESS DEPENDENT ALERTS
   if the core switch is down, do not send 400 alerts for what is behind it.
   Model dependencies, or alert on the root and summarise the rest.

ALWAYS ALERT ON:
   * certificate expiry (30/14/7 days) -- the most preventable outage there is
   * domain expiry
   * DNS resolution failures from an external vantage point
   * interface errors trending up (before they become loss)
   * conntrack / accept queue / socket exhaustion (Ch 66)
```

### 69.6 A minimal but complete stack

```
metrics        Prometheus (+ node_exporter, snmp_exporter, blackbox_exporter)
               or InfluxDB/Telegraf; Grafana for dashboards
flows          nfdump / GoFlow2 / ntopng, or cloud VPC flow logs -> object storage
synthetics     blackbox_exporter (ICMP/TCP/HTTP/DNS/TLS probes) + external provider
packets        tcpdump ring buffer on key hosts; a TAP/mirror to a capture box
                for critical segments
logs           syslog/journald -> Loki/Elasticsearch, correlated by timestamp
alerting       Alertmanager with dependency-aware routing and burn-rate rules
```

`blackbox_exporter` deserves a specific mention: one binary gives you ICMP, TCP
connect, HTTP(S) with content and status assertions, **TLS certificate expiry**,
and DNS checks — covering §69.2 almost entirely.

### 69.7 Real-world example

A company monitored their site with a single HTTP check from one cloud region. A
BGP route leak (Ch. 50) made the site unreachable from three European ISPs for
90 minutes; the check — which ran from the same continent as the servers —
stayed green throughout. Customers reported it. Adding probes from multiple
external vantage points (and a RUM beacon) meant the next such event was detected
in under two minutes by *both* the geographic spread of synthetic failures and a
regional drop in real-user traffic.

**Try it yourself**
```bash
# a quick multi-layer probe of one target
H=example.com
ping -c 5 $H | tail -2
nc -zv -w3 $H 443
echo | openssl s_client -connect $H:443 -servername $H 2>/dev/null \
  | openssl x509 -noout -enddate                    # certificate expiry
curl -o /dev/null -s -w 'dns %{time_namelookup} tcp %{time_connect} tls %{time_appconnect} ttfb %{time_starttransfer} code %{http_code}\n' https://$H

# certificate expiry in days (put this in monitoring)
END=$(echo | openssl s_client -connect $H:443 -servername $H 2>/dev/null \
      | openssl x509 -noout -enddate | cut -d= -f2)
echo "$(( ( $(date -d "$END" +%s 2>/dev/null || date -jf '%b %d %T %Y %Z' "$END" +%s) - $(date +%s) ) / 86400 )) days left"
```

---
# Part XXI — Routing protocols in depth

Chapter 13 covered routing tables and longest-prefix match; Chapter 50 gave BGP
and anycast in one page. This part is the detail: how routers *build* those
tables, and what goes wrong.

## 70. How routing protocols work: the two families

### 70.1 The problem

A router knows its directly-connected networks. Everything else must be learned.
Static routes don't scale and don't react to failure. So routers talk to each
other — and there are exactly two ways to do it.

### 70.2 Distance-vector vs link-state

```
DISTANCE-VECTOR                          LINK-STATE
"routing by rumour"                      "everyone gets the same map"

Each router tells its NEIGHBOURS         Each router floods a description of
the DISTANCE to every destination        ITS OWN LINKS to EVERY router

  A: "I can reach X, cost 3"               A: "I connect to B (cost 1) and C (cost 5)"
  B believes A, adds its own cost,         every router assembles the same
  advertises "X, cost 4"                   topology database, then runs
                                            Dijkstra's shortest-path algorithm
                                            itself

+ simple, low memory                     + fast, loop-free convergence
+ small routers can run it               + full topology visibility
- slow convergence                       - more CPU and memory
- COUNT-TO-INFINITY loops                - needs hierarchy (areas) to scale
- limited visibility                     - more complex to design

RIP, EIGRP*, BGP*                        OSPF, IS-IS
```

`*` EIGRP is "advanced distance-vector" (uses DUAL for loop-freedom); BGP is
**path-vector** — a distance-vector variant that carries the *whole AS path*,
which is how it detects loops (§72.3).

### 70.3 Count-to-infinity, and the fixes

The classic distance-vector failure:

```
A --- B --- C --- X          C loses X.
                             B still advertises "X, cost 2" (learned from C).
                             C believes B, installs "X via B, cost 3".
                             B now hears cost 3, updates to 4. And so on,
                             counting slowly to infinity while packets loop.
```

Mitigations you'll see named in exams and configs:

| Technique | What it does |
|---|---|
| **Split horizon** | Never advertise a route back out the interface you learned it on |
| **Poison reverse** | Advertise it back with infinite metric (explicitly "don't use me") |
| **Hold-down timer** | Ignore updates about a route for N seconds after it fails |
| **Triggered updates** | Send changes immediately rather than waiting for the periodic timer |
| **Max metric = infinity** | RIP caps hop count at 15; 16 means unreachable — bounds the counting |

Link-state protocols avoid the problem entirely: nobody relies on a neighbour's
summary, because everyone computes from the same map.

### 70.4 IGP vs EGP, and administrative distance

```
IGP (Interior Gateway Protocol)   inside ONE administrative domain
                                   goal: find the FASTEST path
                                   OSPF, IS-IS, EIGRP, RIP

EGP (Exterior Gateway Protocol)   BETWEEN administrative domains
                                   goal: enforce POLICY and business relationships
                                   BGP (the only one in use)
```

When a router learns the same prefix from multiple protocols, **administrative
distance** (a trustworthiness ranking, Cisco convention) breaks the tie *before*
metrics are even compared:

```
Connected            0        <- always wins
Static               1
eBGP                20
EIGRP (internal)    90
OSPF               110
IS-IS              115
RIP                120
iBGP               200
Unreachable        255
```

Note **eBGP (20) beats OSPF (110)** but **iBGP (200) loses to it** — a deliberate
design that trips people up. Longest-prefix match still applies *first*
(Ch. 13.2); administrative distance only decides between equal-length prefixes.

### 70.5 Convergence: the number that matters

```
convergence time = detect + propagate + compute + install (FIB)

detect     link down (fast: ~ms with carrier loss)
           or hello timeout (slow: seconds)  -> use BFD (below)
propagate  flood the change
compute    run SPF / DUAL / best-path
install    push into hardware forwarding tables
```

**BFD (Bidirectional Forwarding Detection)** deserves special mention: a
lightweight hello protocol running at sub-second (often 50–300 ms) intervals that
tells the routing protocol "the neighbour is gone" far faster than protocol hello
timers can. It's the standard way to get sub-second failover, and it works with
OSPF, IS-IS, and BGP alike.

**Try it yourself**
```bash
# FRRouting gives you real OSPF/BGP on Linux, free
sudo apt install frr
sudo vtysh -c 'show ip route'          # the RIB, with protocol codes and AD
ip route show proto ospf 2>/dev/null   # kernel FIB entries installed by OSPF
ip route show proto bgp  2>/dev/null
```

---

## 71. OSPF in depth

### 71.1 The model

**Open Shortest Path First** (RFC 2328 for OSPFv2/IPv4; RFC 5340 for OSPFv3/IPv6)
is the dominant enterprise link-state IGP.

```
1. Discover neighbours    (Hello packets, multicast 224.0.0.5)
2. Form ADJACENCIES       (not with everyone -- see DR/BDR)
3. Exchange the LINK-STATE DATABASE (LSDB) -- every router ends up with an
   IDENTICAL copy for its area
4. Run DIJKSTRA (SPF) on that database, from itself, to get shortest paths
5. Install the results in the routing table
6. Re-flood on any change; re-run SPF
```

OSPF runs **directly on IP, protocol number 89** — not TCP or UDP. It provides
its own reliability (acknowledged LSAs).

### 71.2 Neighbour states and adjacency

```
Down -> Init -> 2-Way -> ExStart -> Exchange -> Loading -> FULL
                  ^                                          ^
          "I see you seeing me"              databases synchronised
```

**Stuck states are diagnostic:**

| Stuck at | Almost always means |
|---|---|
| `Init` | Hellos are one-way — the other side isn't hearing you (ACL, unidirectional link) |
| `2-Way` (on a broadcast segment) | **Normal** — non-DR routers stay 2-Way with each other by design |
| `ExStart` | **MTU mismatch** between neighbours (Ch. 15 again) |
| `Loading` | Packet loss or LSA problems |
| Flapping | Hello/dead timer mismatch, or a genuinely unstable link |

**Hello/dead intervals must match** (default 10 s / 40 s on broadcast links), as
must the area ID, authentication, and the stub-area flag, or the adjacency never
forms.

### 71.3 DR and BDR

On a broadcast segment with `n` routers, full mesh adjacency would need
`n(n-1)/2` relationships. OSPF elects a **Designated Router** and **Backup DR**;
everyone else forms full adjacency **only with the DR/BDR** and stays `2-Way`
with each other.

```
Election: highest OSPF priority (default 1; 0 = never DR), tie-break on
          highest Router ID.
          NOT pre-emptive -- a better router arriving later does NOT take over.
```

DR/BDR does not exist on point-to-point links. On point-to-multipoint or NBMA
topologies, misconfigured network types are a classic source of "adjacency won't
form".

### 71.4 Areas and LSA types

OSPF scales by **hierarchy**. Every area must connect to **area 0** (the
backbone); inter-area traffic transits area 0.

```
        +------------------ Area 0 (backbone) ------------------+
        |                                                       |
   [ABR]|                                                       |[ABR]
        |                                                       |
   Area 1                                                    Area 2
   (SPF runs only within an area; ABRs summarise between them)
```

| LSA type | Name | Carries | Scope |
|---|---|---|---|
| **1** | Router LSA | This router's own links and costs | within the area |
| **2** | Network LSA | The segment, generated by the DR | within the area |
| **3** | Summary LSA | Inter-area prefixes, generated by an **ABR** | into another area |
| **4** | ASBR Summary | How to reach an ASBR | between areas |
| **5** | External LSA | Routes redistributed from outside OSPF, by an **ASBR** | flooded domain-wide |
| **7** | NSSA External | Externals inside a not-so-stubby area (converted to type 5 at the ABR) | within an NSSA |

**Area types** control which LSAs are allowed in, trading visibility for size:

| Area type | Blocks | Result |
|---|---|---|
| Standard | nothing | Full LSDB |
| **Stub** | type 5 (externals) | ABR injects a default route instead |
| **Totally stubby** | types 3, 4, 5 | Only a default route — smallest tables |
| **NSSA** | type 5, but allows type 7 | A stub area that still has its own external connection |

**Roles:** ABR = Area Border Router (touches area 0 and another area);
ASBR = Autonomous System Boundary Router (redistributes external routes in).

### 71.5 Cost, and the reference-bandwidth trap

```
cost = reference_bandwidth / interface_bandwidth

Default reference bandwidth = 100 Mbit/s. Therefore:
   100 Mbit link  -> cost 1
   1 Gbit link    -> cost 1     <-- same!
   10 Gbit link   -> cost 1     <-- same!
   100 Gbit link  -> cost 1     <-- same!
```

**Every link at or above 100 Mbit/s has cost 1 by default**, so OSPF cannot tell a
gigabit link from a 100-gigabit one, and picks by hop count instead. Fix:
`auto-cost reference-bandwidth 100000` (100 Gbit/s) — **and set it identically on
every router**, or you get asymmetric routing.

Equal-cost paths are load-balanced (ECMP, Ch. 50.3).

### 71.6 OSPFv3 differences

- Runs over **IPv6** (protocol 89 still), addresses in LSAs are separated from
  topology, so it can carry multiple address families.
- Adjacencies form over **link-local** addresses (`fe80::`).
- Authentication is delegated to **IPsec** rather than being built in.
- Router IDs are still 32-bit dotted-quad values, even in IPv6.

### 71.7 Real-world example

A campus network had two data centres connected by both a 10 Gbit/s dark-fibre
link (2 hops) and a 1 Gbit/s backup (1 hop). Traffic consistently chose the 1
Gbit/s path. Cause: default reference bandwidth, so both links cost 1, and the
1-hop path had a lower total cost. Raising the reference bandwidth to 100 Gbit/s
on every router made the fibre cost 10 and the backup cost 100 — traffic moved
immediately. **Check reference bandwidth on every OSPF network you inherit.**

**Try it yourself**
```bash
# FRRouting on two Linux boxes/containers
sudo vtysh
  configure terminal
   router ospf
    ospf router-id 1.1.1.1
    auto-cost reference-bandwidth 100000
    network 10.0.0.0/24 area 0
   exit
  end
  show ip ospf neighbor          # look for state FULL
  show ip ospf database          # the LSDB -- LSA types visible
  show ip ospf interface eth0    # cost, hello/dead, DR/BDR, network type
  show ip route ospf
```

---

## 72. BGP in depth: attributes and path selection

### 72.1 The model

**BGP-4** (RFC 4271) connects ~75,000 **autonomous systems**. It is a
**path-vector** protocol: advertisements carry the full list of ASes traversed.

```
Runs over TCP port 179  (so it inherits TCP reliability, and needs
                          IP reachability BEFORE it can peer)

eBGP   between different ASes   -- usually directly connected, TTL 1 by default
iBGP   within one AS            -- carries external routes across your network;
                                   peers are often multiple hops apart
```

**The iBGP full-mesh rule:** a route learned from one iBGP peer is **not**
re-advertised to another iBGP peer (that's how BGP prevents internal loops).
Therefore every iBGP speaker must peer with every other — `n(n-1)/2` sessions —
unless you use route reflectors or confederations (§73.2).

### 72.2 The best-path selection algorithm

When BGP has several paths to a prefix, it walks this list and stops at the first
tiebreak. **Memorise the first four; the rest are rare.**

```
 0. Prefer the path with a reachable NEXT_HOP        (if not reachable, discard)
 1. Highest WEIGHT                    (Cisco-proprietary, LOCAL to one router)
 2. Highest LOCAL_PREF                (iBGP-wide; THE knob for choosing outbound)
 3. Locally originated                (network / redistribute / aggregate)
 4. Shortest AS_PATH                  (the famous one -- but it is only #4)
 5. Lowest ORIGIN                     (IGP < EGP < Incomplete)
 6. Lowest MED                        (a hint TO a neighbour about inbound; only
                                       compared between paths from the SAME AS
                                       by default)
 7. eBGP over iBGP
 8. Lowest IGP metric to the next hop
 9. Oldest path (stability)
10. Lowest neighbour Router ID
11. Lowest neighbour IP address
```

**The crucial insight:** `LOCAL_PREF` (step 2) beats `AS_PATH` (step 4). Policy
beats topology. The internet's shape is commercial, not geographic — which is why
your packets sometimes cross a continent to reach a server two miles away.

### 72.3 The attributes that matter

| Attribute | Type | Direction | Use |
|---|---|---|---|
| **NEXT_HOP** | well-known mandatory | — | Where to send the packet. On iBGP it is **not** rewritten by default → the classic "route is in the table but unusable" problem. Fix: `next-hop-self`. |
| **AS_PATH** | well-known mandatory | outbound signal | Loop detection (**reject any path containing my own ASN**) and path length. **AS-path prepending** artificially lengthens it to make a path less attractive — the standard way to influence *inbound* traffic. |
| **LOCAL_PREF** | well-known discretionary | **iBGP only, internal** | Choose which exit *your* traffic uses. Higher wins. Default 100. |
| **MED** (MULTI_EXIT_DISC) | optional non-transitive | to a neighbour AS | "If you have two links to me, prefer this one." Lower wins. **A hint your neighbour may ignore.** |
| **ORIGIN** | well-known mandatory | — | How the route entered BGP: IGP(i) < EGP(e) < Incomplete(?) |
| **COMMUNITIES** | optional transitive | tagging | 32-bit tags (`64512:100`) used to signal policy: "don't advertise to peers", "prepend twice in Europe", "blackhole this". Large communities (RFC 8092) extend this to 96 bits for 4-byte ASNs. |
| **ATOMIC_AGGREGATE / AGGREGATOR** | — | — | Marks that summarisation lost detail |

**Well-known communities** you should recognise:
`NO_EXPORT` (don't send outside this AS), `NO_ADVERTISE` (don't send to anyone),
`NO_EXPORT_SUBCONFED`, and the widely-implemented **blackhole** community
`65535:666` (RFC 7999) — "drop all traffic to this prefix", used for DDoS
mitigation (Ch. 54.2).

### 72.4 Influencing traffic: the asymmetry

```
OUTBOUND traffic (yours leaving)   -> YOU control it, easily
      set LOCAL_PREF on routes learned from each provider

INBOUND traffic (coming to you)    -> you can only HINT
      * AS-path prepending (make one path look longer)
      * MED (only works between two links to the SAME neighbour)
      * communities your provider publishes (best method -- ask them)
      * advertise MORE-SPECIFIC prefixes out one link (blunt but effective,
        because longest-prefix match beats every BGP attribute)
```

**Nothing you do guarantees inbound path selection** — the remote network's own
`LOCAL_PREF` beats your prepending every time. This asymmetry is the root of most
multihoming frustration.

### 72.5 Real-world example

A company multihomed to two ISPs and prepended its ASN three times toward the
more expensive one, expecting traffic to shift. Nothing changed: a major content
network had a `LOCAL_PREF` policy preferring that ISP as a settlement-free peer,
and `LOCAL_PREF` is evaluated at step 2 — long before `AS_PATH` at step 4. The fix
was to use the ISP's published **community** for "de-prefer this route toward your
peers", which acts on the *remote* side's `LOCAL_PREF`.

**Try it yourself**
```bash
# public looking glasses show real BGP tables
# https://lg.he.net  or  https://stat.ripe.net
curl -s "https://stat.ripe.net/data/routing-status/data.json?resource=8.8.8.0/24" | python3 -m json.tool | head -40

whois -h whois.radb.net -- '-i origin AS15169' | grep -c '^route:'   # prefixes an AS originates
traceroute -A 1.1.1.1        # -A annotates each hop with its ASN

# in FRR
sudo vtysh -c 'show bgp ipv4 unicast 8.8.8.0/24'   # all paths + why one was chosen
sudo vtysh -c 'show bgp summary'
```

---

## 73. BGP operations: peering, scaling, and failures

### 73.1 Business relationships shape routing

```
CUSTOMER   pays you      -> you advertise their prefixes to EVERYONE
PEER       settlement-free -> you exchange only YOUR and YOUR CUSTOMERS' prefixes
PROVIDER   you pay them  -> they give you a default/full table; you advertise
                             only yourself and your customers to them

The "valley-free" rule that falls out of this:
   traffic goes  customer -> provider  (up),  then across a peer link,
   then  provider -> customer  (down).  It never goes UP after coming DOWN.
```

Violating this by accident is a **route leak** (§73.4). Connections happen at
**IXPs** (internet exchange points) via a shared fabric and often a **route
server** that redistributes among many participants.

### 73.2 Scaling iBGP

The full-mesh requirement (§72.1) is quadratic. Two solutions:

| Solution | How | Trade-off |
|---|---|---|
| **Route reflectors (RR)** | Designated routers *may* re-advertise iBGP routes to clients. Clients peer only with the RRs. | Simple and dominant. RRs reflect only their *best* path, which can hide alternatives and cause suboptimal routing; use redundant RRs and `add-path` to mitigate. |
| **Confederations** | Split the AS into sub-ASes that eBGP-peer with each other internally | More complex; rare outside very large networks |

Loop prevention within reflection uses `ORIGINATOR_ID` and `CLUSTER_LIST`.

### 73.3 Session states and why sessions won't come up

```
Idle -> Connect -> Active -> OpenSent -> OpenConfirm -> ESTABLISHED
```

| Symptom | Cause |
|---|---|
| Stuck **Idle** | No route to the peer, or administratively shut, or damped |
| Flapping **Active** | TCP 179 blocked, wrong peer IP, ACL, or the peer isn't configured for you |
| Stuck **OpenSent/OpenConfirm** | **AS number mismatch**, router-ID collision, unsupported capability, or MD5/TCP-AO password mismatch |
| Established but **no prefixes** | Filters/route-maps rejecting everything, or no `network` statement; check `show bgp neighbor x advertised-routes` and `received-routes` |
| Session drops every ~90 s | Hold timer expiry — usually MTU (large updates dropped) or CPU |

**eBGP is TTL 1 by default** (peers must be directly connected). For multihop
peering use `ebgp-multihop`, and prefer **GTSM** (RFC 5082, `ttl-security`) which
requires the packet to arrive with TTL ≥ 254 — a cheap anti-spoofing control.

**Protect the session:** TCP-MD5 (RFC 2385) or the modern **TCP-AO** (RFC 5925),
plus infrastructure ACLs so only your peers can reach port 179.

### 73.4 The three ways BGP breaks the internet

**1. Prefix hijack.** An AS announces a prefix it doesn't own. Because
**longest-prefix match wins** (Ch. 13.2), announcing a `/24` out of someone's
`/16` steals the traffic globally. (Pakistan/YouTube 2008; repeated crypto-DNS
thefts.)

**2. Route leak.** An AS re-announces routes it shouldn't — e.g. re-advertising
one provider's routes to another provider, becoming an accidental transit for
traffic it cannot carry. Usually a filter mistake, and usually causes congestion
and blackholing rather than theft.

**3. Instability / flapping.** A flapping link generates constant updates that
propagate globally. **Route flap damping** suppresses repeatedly-flapping prefixes
(historically too aggressively; modern RIPE-recommended parameters are gentler).

### 73.5 Defences

| Control | What it does |
|---|---|
| **Prefix filters / IRR** | Only accept prefixes a customer is registered (in RADB/RIPE) to originate. The oldest and still-essential control. |
| **RPKI + ROV** | A **ROA** cryptographically states "AS X may originate prefix P with max length L". Routers validate and **drop Invalid** routes. **The single highest-value modern defence** — publish ROAs for every prefix you own. |
| **Max-prefix limits** | Tear down a session that suddenly sends 500,000 routes instead of 50. Catches leaks fast. |
| **AS-path filters** | Reject paths containing your own ASN, or unreasonable lengths |
| **Bogon filtering** | Drop private/reserved/unallocated space |
| **BGPsec** | Cryptographic path validation (not just origin). Barely deployed. |
| **MANRS** | An industry programme bundling filtering, anti-spoofing (BCP 38), coordination, and validation |

**Peer locking** and **RPKI-based origin validation** together stop most real
incidents. If you run BGP and do nothing else, **publish ROAs and reject
Invalids**.

### 73.6 Real-world example

A regional ISP's new engineer redistributed the full BGP table into OSPF and back
into BGP without filters. For eleven minutes the ISP advertised ~400,000 prefixes
to both its upstream providers, claiming to be the best path to much of the
internet. Traffic that its 10 Gbit/s links could never carry arrived and was
dropped. Two things limited the damage: one upstream had a **max-prefix limit of
1,000** and shut the session in seconds; the other did not, and carried the leak
for the full duration. **Max-prefix limits are the cheapest insurance in
networking.**

**Try it yourself**
```bash
# is a prefix RPKI-valid?   https://rpki-validator.ripe.net  or:
curl -s "https://stat.ripe.net/data/rpki-validation/data.json?resource=AS15169&prefix=8.8.8.0/24" \
  | python3 -m json.tool | grep -i status

# who originates this prefix, and has it changed?
curl -s "https://stat.ripe.net/data/prefix-overview/data.json?resource=1.1.1.0/24" | python3 -m json.tool | head -30

# monitor your own prefixes for hijacks: BGPalerter, RIPE Atlas, or bgp.tools
```

---

## 74. First-hop redundancy, multicast, and MPLS/VRF in brief

Three topics you will meet, compressed to what you need to recognise them.

### 74.1 First-hop redundancy (VRRP / HSRP / GLBP)

**The problem:** hosts have one default gateway (Ch. 13). If that router dies,
everything on the subnet is isolated — even if a second router is sitting right
there.

**The solution:** two or more routers share a **virtual IP and virtual MAC**.
Hosts point at the virtual IP and never know which physical router is answering.

```
   Router A (Master, priority 110)  \
                                      >--  VIP 192.168.1.1, VMAC 00:00:5e:00:01:XX
   Router B (Backup, priority 100)  /

   Hosts:  default gateway 192.168.1.1  (the VIP)

   A sends VRRP advertisements (multicast 224.0.0.18, IP protocol 112)
   every ~1 s. If B stops hearing them for ~3 s, B takes over the VIP and
   VMAC, and sends a gratuitous ARP so switches relearn the port.
```

| Protocol | Notes |
|---|---|
| **VRRP** (RFC 5798) | The **standard**; VMAC `00:00:5e:00:01:<VRID>`. Use this. |
| **HSRP** | Cisco-proprietary; same idea, different timers and terminology (Active/Standby) |
| **GLBP** | Cisco; also load-balances by handing different hosts different virtual MACs |
| **CARP** | The BSD equivalent |
| **keepalived** | The common Linux implementation of VRRP; also used for service VIPs |

Practicalities: **preemption** decides whether the higher-priority router takes
back over when it returns (usually yes, with a delay so it can rebuild routing
first). **Interface tracking** lowers priority if the *uplink* fails, so you don't
fail over to a router that can't reach anything. In modern data centres,
**anycast gateways** on every leaf switch (with EVPN) are replacing VRRP.

### 74.2 Multicast

**The problem:** send one stream to 500 receivers without sending 500 copies.

```
UNICAST     500 copies leave the server            (bandwidth = 500x)
BROADCAST   everyone gets it, wanted or not        (only works on one LAN)
MULTICAST   ONE copy per LINK; routers replicate only where receivers exist
```

Addresses: **224.0.0.0/4** (IPv4), `ff00::/8` (IPv6). Well-known ones you'll see:
`224.0.0.1` (all hosts), `224.0.0.2` (all routers), `224.0.0.5/6` (OSPF),
`224.0.0.18` (VRRP), `224.0.0.251` (mDNS).

Two protocols do the work:

| Protocol | Scope | Job |
|---|---|---|
| **IGMP** (v2/v3) / **MLD** for IPv6 | Host ↔ local router | "I want to receive group G" (Join) / "I'm done" (Leave). **IGMPv3** adds source filtering (SSM). |
| **PIM** | Router ↔ router | Builds the distribution tree across the network |

PIM modes: **PIM-SM** (sparse mode — receivers are few and scattered; uses a
**Rendezvous Point** to bootstrap, then switches to a shortest-path tree),
**PIM-SSM** (source-specific; no RP needed, simplest and preferred for one-to-many
like IPTV), **PIM-DM** (dense mode — flood and prune; obsolete).

On switches, **IGMP snooping** is essential: without it, a switch treats multicast
as broadcast and floods every port, defeating the point.

**Where you'll actually meet multicast:** IPTV and video distribution, financial
market data feeds, `mDNS`/Bonjour service discovery, routing protocols
themselves, and PTP time sync. It is largely absent from the public internet and
from most cloud providers (AWS supports it only within Transit Gateway multicast
domains) — it lives in enterprise and provider networks.

### 74.3 MPLS and VRFs

**MPLS** (Multiprotocol Label Switching) forwards on a short **label** instead of
doing a longest-prefix IP lookup.

```
   Ingress router (LER):  push a LABEL onto the packet
   Core routers (LSR):    swap label -> label, forward   (fast, no IP lookup)
   Egress router (LER):   pop the label, forward as normal IP

   Header sits between L2 and L3 -- hence "layer 2.5"
   [ Ethernet | MPLS label(s) | IP | TCP | data ]
```

Labels are distributed by **LDP** (simple) or **RSVP-TE / SR** (traffic
engineering). Modern networks increasingly use **Segment Routing (SR-MPLS or
SRv6)**, which encodes the path as a stack of segments in the packet and removes
the need for a separate label-distribution protocol and per-flow state.

**Why it exists** (the original "faster than IP lookup" reason is obsolete — modern
hardware does IP lookups at line rate):

1. **Traffic engineering** — steer traffic along a chosen path, not just the IGP
   shortest one.
2. **Fast reroute (FRR)** — pre-computed backup paths, sub-50 ms failover.
3. **VPNs** — the dominant use today.

**VRF** (Virtual Routing and Forwarding) is the related idea: **multiple
independent routing tables in one router**.

```
   VRF "customer-A"   10.0.0.0/8 -> ...
   VRF "customer-B"   10.0.0.0/8 -> ...     <-- same prefix, different table,
   VRF "management"   ...                        no conflict
```

**MPLS L3VPN** combines them: each customer gets a VRF, and **MP-BGP** carries
their routes across the provider core, tagged with a **route distinguisher** (to
keep overlapping prefixes unique) and **route targets** (to control which VRFs
import which routes). This is how a carrier sells "private WAN" to thousands of
customers over one shared network.

**EVPN** is the modern successor for data centres: BGP-based control plane over a
VXLAN data plane (Ch. 56), providing L2 and L3 multi-tenancy without the legacy
flood-and-learn behaviour.

**On Linux**, VRFs are real and usable:

```bash
sudo ip link add vrf-red type vrf table 100
sudo ip link set vrf-red up
sudo ip link set eth1 master vrf-red
sudo ip route add default via 10.1.1.1 table 100
ip route show vrf vrf-red
```

### 74.4 Real-world example

A hospital ran patient monitoring over multicast and lost feeds intermittently on
one floor. The floor's access switch had been replaced with a model where **IGMP
snooping was disabled by default**, so multicast flooded every port — saturating
the 1 Gbit/s uplink shared with everything else, and causing drops that looked
like random monitor failures. Enabling snooping (and configuring an IGMP querier,
since there was no multicast router on that VLAN) resolved it. **Multicast without
snooping is broadcast.**

**Try it yourself**
```bash
# see multicast groups your machine has joined
ip maddr show
netstat -gn 2>/dev/null || ip -4 maddr

# watch IGMP and multicast on the wire
sudo tcpdump -i any -n 'igmp or (ip multicast)'

# mDNS is multicast you already run
sudo tcpdump -i any -n 'udp port 5353' -c 10
avahi-browse -a 2>/dev/null | head        # Linux
dns-sd -B _services._dns-sd._udp          # macOS

# VRF experiment (Linux, needs root)
sudo ip link add vrf-test type vrf table 200 && sudo ip link set vrf-test up
ip -br link show type vrf
sudo ip link del vrf-test
```

---
# Appendix A — Troubleshooting playbooks (symptom → cause → fix)

> General method: **work up the stack**. Link (carrier/errors) → IP
> (address/route/ARP) → transport (handshake/RST/retransmit) → app (status codes/
> latency). Capture on **both ends** for anything about loss or latency. Change
> one thing at a time. Note the exact error string — it encodes the layer
> (§24.7).

### A1. "Connection timed out" (no response at all)

| Check | Command | If... |
|---|---|---|
| Route exists? | `ip route get <dst>` | `unreachable` → add route / fix default GW |
| Next-hop ARP? | `ip neigh get <gw> dev <if>` | `FAILED` → L2 issue, gateway down, wrong VLAN |
| SYNs leaving? | `sudo tcpdump -ni any 'tcp[tcpflags]&tcp-syn!=0 and host <dst>'` | no SYN out → local firewall (`nft list ruleset`), routing |
| SYN-ACK back? | same capture | SYN out, nothing back → remote firewall `DROP`, host down, wrong port, ACL, security group |
| Only large packets fail? | `ping -M do -s 1472 <dst>` fails but small ok, handshake ok then hang | **PMTU black hole** — allow ICMP type 3 / clamp MSS / `tcp_mtu_probing=1` (Ch. 8) |
| Intermittent, some regions | `mtr -z <dst>` from each region; check RPKI | **BGP hijack/leak** (Ch. 50) |

### A2. "Connection refused" (`ECONNRESET`/`ECONNREFUSED` immediately)

- RST to your SYN → **nothing listening** on that port. Verify: `ss -tlnp` on the
  server; is it bound to `127.0.0.1` only (not `0.0.0.0`/`::`)? Right port? Right
  container/netns? Service crashed?
- Connect succeeds then instant RST → **accept-queue overflow**
  (`nstat | grep ListenOverflows`, raise `listen()` backlog + `somaxconn`), or
  app crashes on `accept()`, or an IPS resetting.

### A3. "Works then hangs on big transfers / TLS stalls"

Signature of **MTU/PMTUD black hole** (Ch. 8). `tcpdump` shows a ~1500 B packet
retransmitted with no ACK.
Fixes: unblock `ICMP type 3` / `ICMPv6 type 2`; MSS clamp on the tunnel
(`iptables ... TCPMSS --clamp-mss-to-pmtu`); `sysctl net.ipv4.tcp_mtu_probing=1`;
lower interface/pod MTU. In K8s: check pod MTU vs underlay (Ch. 55).

### A4. Intermittent latency spikes at exactly 200 ms / 1 s / 3 s

- **200 ms / 1 s flat floor** on request/response → **Nagle × delayed-ACK**
  (Ch. 33). Fix: `TCP_NODELAY` on the client, or write the whole request in one
  `send()`.
- **200 ms tail on the last segment of small responses** → **RTO after tail
  loss** (Ch. 28). Fix: RACK-TLP (newer kernel), `ip route ... rto_min 20ms` in
  the DC.
- **1 s / 3 s connect latency** → SYN or final-ACK loss / accept-queue overflow
  (`nstat | grep -E 'ListenDrops|ListenOverflows|SynRetrans'`).
- **5 s DNS stalls** → a dead server in `resolv.conf`; add `options timeout:1
  attempts:2 rotate` or run a local caching resolver (Ch. 20.8).

### A5. Throughput far below link rate, no packet loss

1. `ss -tie dst <peer>` → is it `app_limited`? → fix the app (single-threaded
   producer, small writes, computing `Content-Length` by buffering).
2. Window scaling on **both** SYNs? (`tcpdump` the handshake). Middlebox
   stripping options? (Ch. 23 war story).
3. `buffer_max ≥ 2 × BDP`? Compute `BDP = bandwidth × RTT` (Ch. 35). Raise
   `net.ipv4.tcp_rmem`/`tcp_wmem` **max** and `net.core.rmem_max`/`wmem_max`.
4. `iperf3 -P8` ≫ single stream? → CC/window-limited; consider BBR
   (Ch. 32) or accept parallelism.
5. Per-flow cloud cap (TGW/NAT GW/VPN/single-flow ~5 Gbit/s)? → parallel streams
   (Ch. 57).
6. RX drops? `ethtool -S`, `/proc/net/softnet_stat`, IRQ affinity (Ch. 48).

### A6. Throughput low **with** packet loss

- Confirm loss is real & directional: **two-ended capture**, compare by seq;
  `mtr` loss must persist to the final hop.
- Link errors: `ethtool -S` `rx_crc_errors`/`rx_errors`; check cabling/SFP/optics,
  duplex mismatch (Ch. 5).
- Reordering masquerading as loss? `nstat | grep -Ei 'reorder|dsack|spurious'`;
  fix LACP/ECMP hash to per-flow (Ch. 27).
- If loss is inherent (Wi-Fi, satellite, transoceanic): **SACK on** (verify),
  switch senders to **BBR**, or use parallel streams / FEC. Mathis:
  `throughput ≈ MSS/(RTT·√p)` (Ch. 36).
- Bufferbloat check: does RTT balloon under load? → AQM (`fq_codel`/`cake`) at
  the bottleneck (Ch. 44).

### A7. Connections die after N seconds/minutes of idle

**Middlebox idle timeout** reaping the flow (Ch. 34, 53). Common values: 350 s
(AWS NLB/NAT GW), 60 s (some ELB/proxies), 300 s (consumer NAT), 5–15 min
(enterprise FW). Fix: application heartbeat / TCP keepalive with **idle interval
shorter than the smallest timeout on the path** (e.g. 30–60 s), and
`TCP_USER_TIMEOUT` for bounded failure.

### A8. File-descriptor leak / `too many open files`

- `ss -tan state close-wait | wc -l` climbing → **your app isn't calling
  `close()`** on a code path (Ch. 25.5). Find it (error/exception branches).
- `ss -tan state fin-wait-2` climbing → the **peer** isn't closing; bounded by
  `tcp_fin_timeout`.
- Legit high count → raise `ulimit -n` / `LimitNOFILE=` / `fs.file-max`.

### A9. Lots of `TIME_WAIT` on a client

Normal if you open many short connections. Only act if near port exhaustion
(`connect()` → `EADDRNOTAVAIL`): enable `net.ipv4.tcp_tw_reuse=1`, widen
`ip_local_port_range`, and — the real fix — **use connection pooling / keep-alive**
(Ch. 21 war story). Never `tcp_tw_recycle`.

### A10. UDP app "losing data"

`nstat | grep -Ei 'UdpInErrors|UdpRcvbufErrors|UdpNoPorts'`. `RcvbufErrors` →
receiver too slow / buffer too small: raise `SO_RCVBUF` + `net.core.rmem_max`,
use `recvmmsg()`, add a drain thread. `InErrors`/`NoPorts` → nothing listening /
checksum failures. Remember UDP loss is **silent** — you must instrument it
(Ch. 22).

### A11. DNS is "slow" or "flaky"

- `dig +trace <name>` — where does it break? NXDOMAIN cached? (`SOA` minimum TTL).
- Capture UDP/53: **how many queries per lookup?** `ndots:5` + `search` → 4–5
  queries; use FQDN with trailing dot (Ch. 20.8).
- Large response truncated (`TC` flag) → is **TCP/53** allowed? IP fragments
  dropped? EDNS buffer too big? Set 1232 (Ch. 15).
- Per-resolver throttle (cloud `.2` resolver ~1024 pps/ENI) → NodeLocal DNSCache
  / local `unbound`.
- DoH in the browser bypassing internal resolver → internal names fail (Ch. 20.7).

### A12. TLS handshake fails / times out for *some* clients

- `openssl s_client -connect host:443 -servername host` — which step fails?
- Cert chain incomplete (missing intermediate)? `openssl ... -showcerts` and test
  from a machine with a fresh trust store.
- `ClientHello` spans two segments (PQC key shares, many extensions) and a
  middlebox/LB only reads the first packet → handshake hangs (Ch. 38 war story).
- Protocol/cipher mismatch (old client, TLS 1.0 disabled), SNI required but not
  sent, clock skew (cert "not yet valid").

### A13. Kubernetes: pod can't reach Service / another pod

- Pod→pod same node works, cross-node fails → **CNI overlay / route** issue:
  `tcpdump` for VXLAN/GENEVE (udp 4789/6081) on the node `eth0`; check tunnel
  offloads (Ch. 56 war story); check node-to-node SG allows the overlay port.
- Pod→Service fails but pod→podIP works → `kube-proxy`/eBPF DNAT rules or
  `Endpoints` empty (readiness failing); `iptables-save | grep <svc-clusterip>`.
- Big responses hang → **pod MTU** vs underlay (Ch. 55 war story).
- Cross-node blocked → **NetworkPolicy** default-deny; `cilium monitor` /
  `calicoctl` to see drops.
- DNS slow → `ndots:5`, CoreDNS pods, conntrack pressure (Ch. 56).

### A14. "The network is dropping packets" (general)

Before opening a carrier ticket, rule out:
- NIC/driver: `ethtool -S`, `/proc/net/softnet_stat`, ring sizes, IRQ affinity.
- Socket buffers / accept queue / conntrack table: `nstat -az | grep -Ei
  'drop|overflow|prune|nf_conntrack'`, `conntrack -S`, `dmesg`.
- qdisc drops (may be intentional AQM): `tc -s qdisc`.
- Firewall `INVALID`/policy drops: rule counters.
- Application 4xx/5xx mislabeled as "network" (Ch. 42 war story) — check at the
  packet layer whether it's RST/timeout (transport) or a clean HTTP status (app).

---

# Appendix B — Port and protocol quick reference

### B1. IP protocol numbers (IPv4 `Protocol` / IPv6 `Next Header`)

| # | Protocol |
|---|---|
| 1 | ICMP |
| 2 | IGMP |
| 4 | IPv4 (IP-in-IP encapsulation) |
| 6 | **TCP** |
| 17 | **UDP** |
| 41 | IPv6 (6in4 encapsulation) |
| 47 | GRE |
| 50 | ESP (IPsec) |
| 51 | AH (IPsec) |
| 58 | ICMPv6 |
| 89 | OSPF |
| 112 | VRRP |
| 132 | SCTP |

### B2. Common TCP/UDP ports

| Port | Proto | Service |
|---|---|---|
| 20/21 | TCP | FTP data / control |
| 22 | TCP | SSH / SCP / SFTP |
| 23 | TCP | Telnet (avoid) |
| 25 | TCP | SMTP (MTA↔MTA) |
| 53 | UDP+**TCP** | DNS (TCP for large/AXFR — must be open) |
| 67/68 | UDP | DHCP server / client |
| 69 | UDP | TFTP |
| 80 | TCP | HTTP |
| 88 | TCP/UDP | Kerberos |
| 110 / 143 | TCP | POP3 / IMAP |
| 123 | UDP | NTP |
| 137–139 | UDP/TCP | NetBIOS |
| 161/162 | UDP | SNMP / SNMP trap |
| 179 | TCP | BGP |
| 389 / 636 | TCP | LDAP / LDAPS |
| 443 | TCP+**UDP** | HTTPS (UDP = HTTP/3 / QUIC) |
| 445 | TCP | SMB / CIFS |
| 465 / 587 | TCP | SMTPS / SMTP submission (STARTTLS) |
| 500 / 4500 | UDP | IKE / IPsec NAT-T |
| 514 | UDP | syslog |
| 636 | TCP | LDAPS |
| 853 | TCP/UDP | DNS over TLS / over QUIC |
| 993 / 995 | TCP | IMAPS / POP3S |
| 1194 | UDP | OpenVPN |
| 1433 / 1521 | TCP | MS SQL / Oracle |
| 1812/1813 | UDP | RADIUS auth / accounting |
| 1935 | TCP | RTMP |
| 2049 | TCP | NFS |
| 3306 | TCP | MySQL/MariaDB |
| 3389 | TCP | RDP |
| 3478 / 5349 | UDP/TCP | STUN/TURN / TURN-TLS |
| 4789 | UDP | VXLAN |
| 5060 / 5061 | UDP/TCP | SIP / SIP-TLS |
| 5432 | TCP | PostgreSQL |
| 5671 / 5672 | TCP | AMQP-TLS / AMQP |
| 6081 | UDP | GENEVE |
| 6379 | TCP | Redis |
| 6443 | TCP | Kubernetes API server |
| 8080 / 8443 | TCP | HTTP alt / HTTPS alt |
| 9092 | TCP | Kafka |
| 9100 | TCP | Prometheus node_exporter |
| 11211 | TCP/UDP | memcached (never expose UDP) |
| 27017 | TCP | MongoDB |
| 51820 | UDP | WireGuard |

### B3. ICMP / ICMPv6 types worth knowing

| ICMPv4 | ICMPv6 | Meaning |
|---|---|---|
| 0 / 8 | 129 / 128 | Echo reply / request (ping) |
| 3 (codes 0–15) | 1 | Destination unreachable |
| 3/4 | 2 | **Fragmentation needed / Packet Too Big** (PMTUD — keep it!) |
| 11 | 3 | Time exceeded (traceroute) |
| 12 | 4 | Parameter problem |
| 5 | 137 | Redirect |
| — | 133/134 | Router Solicitation / Advertisement (NDP) |
| — | 135/136 | Neighbor Solicitation / Advertisement (NDP) |

### B4. EtherTypes

| Value | Protocol |
|---|---|
| 0x0800 | IPv4 |
| 0x0806 | ARP |
| 0x86DD | IPv6 |
| 0x8100 | 802.1Q VLAN tag |
| 0x88A8 | 802.1ad QinQ |
| 0x8847/0x8848 | MPLS unicast / multicast |
| 0x8864 | PPPoE session |
| 0x88CC | LLDP |

---

# Appendix C — Subnetting cheat sheet

### C1. Prefix ↔ mask ↔ hosts

| CIDR | Netmask | Wildcard | Block | Usable hosts |
|---|---|---|---|---|
| /8 | 255.0.0.0 | 0.255.255.255 | 16 777 216 | 16 777 214 |
| /16 | 255.255.0.0 | 0.0.255.255 | 65 536 | 65 534 |
| /20 | 255.255.240.0 | 0.0.15.255 | 4 096 | 4 094 |
| /21 | 255.255.248.0 | 0.0.7.255 | 2 048 | 2 046 |
| /22 | 255.255.252.0 | 0.0.3.255 | 1 024 | 1 022 |
| /23 | 255.255.254.0 | 0.0.1.255 | 512 | 510 |
| /24 | 255.255.255.0 | 0.0.0.255 | 256 | 254 |
| /25 | 255.255.255.128 | 0.0.0.127 | 128 | 126 |
| /26 | 255.255.255.192 | 0.0.0.63 | 64 | 62 |
| /27 | 255.255.255.224 | 0.0.0.31 | 32 | 30 |
| /28 | 255.255.255.240 | 0.0.0.15 | 16 | 14 |
| /29 | 255.255.255.248 | 0.0.0.7 | 8 | 6 |
| /30 | 255.255.255.252 | 0.0.0.3 | 4 | 2 |
| /31 | 255.255.255.254 | 0.0.0.1 | 2 | 2 (RFC 3021 P2P) |
| /32 | 255.255.255.255 | 0.0.0.0 | 1 | 1 (host route) |

Mask octet values: `/25→.128  /26→.192  /27→.224  /28→.240  /29→.248  /30→.252`.

### C2. The 4-step method (recap from Ch. 11)

Given `A.B.C.D /p`:
1. **Interesting octet** = `ceil(p/8)`; **block size** = `256 − mask_octet` in
   that octet.
2. **Network** = round `D` (or the interesting octet) **down** to a multiple of
   the block size.
3. **Broadcast** = network + block size − 1 (in that octet), host bits all 1.
4. **Usable** = network+1 … broadcast−1 (except /31, /32).

### C3. Worked examples

| Address/CIDR | Network | Broadcast | First–Last usable | Hosts |
|---|---|---|---|---|
| 192.168.1.130/26 | 192.168.1.128 | 192.168.1.191 | .129 – .190 | 62 |
| 10.0.0.45/27 | 10.0.0.32 | 10.0.0.63 | .33 – .62 | 30 |
| 172.16.20.100/22 | 172.16.20.0 | 172.16.23.255 | .20.1 – .23.254 | 1022 |
| 100.64.5.7/30 | 100.64.5.4 | 100.64.5.7 | .5 – .6 | 2 |
| 203.0.113.17/28 | 203.0.113.16 | 203.0.113.31 | .17 – .30 | 14 |

### C4. VLSM allocation rule

Sort requirements **largest first**, allocate from the top of the block, next
subnet starts at previous network + its block size. (Allocating small subnets
first fragments the space — Ch. 11.3.)

### C5. IPv6 quick facts

- LAN = **/64** always (SLAAC assumes it).
- Site from ISP: **/48** (65 536 /64s) or **/56** (256 /64s).
- P2P link: still use a /64 (or /127 per RFC 6164 to dodge a ping-pong attack).
- `::1` loopback, `fe80::/10` link-local (always present), `fd00::/8` ULA
  (private), `2000::/3` global, `ff00::/8` multicast, `2001:db8::/32` docs.
- No broadcast; no ARP (NDP via ICMPv6); routers never fragment.

---

# Appendix D — `tcpdump` / `ss` / `ip` command cookbook

### D1. `ip` (iproute2)

```bash
ip -br -c addr                     # brief, colored: interfaces + IPs
ip -4 addr show dev eth0
ip -6 addr show
ip link set eth0 up / down
ip link set eth0 mtu 9000
ip link set eth0 txqueuelen 1000
ip route                            # main routing table
ip route get 8.8.8.8               # exact forwarding decision (src, dev, gw, mtu)
ip route add 10.0.0.0/8 via 10.1.1.1 dev eth0
ip route add default via 192.168.1.1
ip route add 10.9.0.0/16 via 10.1.1.1 mtu 1400        # pin MTU
ip route flush cache
ip rule                             # policy routing rules
ip rule add from 192.168.50.0/24 table 100
ip route add default via 10.0.0.1 table 100
ip neigh                            # ARP / NDP cache + states
ip neigh flush all
ip -s link show eth0               # RX/TX packets, errors, dropped
ip -s -s link show eth0           # + detailed error breakdown
ip netns list ; ip netns exec red <cmd>
ip -d link show vxlan0            # decode VXLAN/GRE/bond details
ip maddr                           # multicast group memberships
```

### D2. `ss`

```bash
ss -tan                            # all TCP, numeric
ss -tanp                           # + process (root for others)
ss -tlnp                           # listening TCP + owner
ss -uanp                           # UDP
ss -s                              # summary by state
ss -tan state established
ss -tan state time-wait | wc -l
ss -tan '( dport = :443 or sport = :443 )'
ss -tan dst 10.0.0.0/8
ss -tie                            # extended + internal TCP info (rtt, cwnd, retrans, CC)
ss -tim                            # + memory per socket
ss -tio                            # + timers (retransmit/keepalive countdown)
ss -tan state syn-recv | wc -l    # SYN flood indicator
ss -tanp | awk '{print $1}' | sort | uniq -c    # state histogram
ss -K dst 10.1.2.3                # kill matching sockets (root; needs CONFIG_INET_DIAG_DESTROY)
```

### D3. `tcpdump`

```bash
# Capture to file, headers only, kernel buffer bumped:
sudo tcpdump -ni eth0 -s128 -B 8192 -w /tmp/c.pcap 'tcp port 443 and host 1.2.3.4'
# Rotating ring (24 x 100MB) for intermittent bugs:
sudo tcpdump -ni eth0 -s0 -C 100 -W 24 -w /tmp/ring.pcap 'port 443'
# Read back with filters:
tcpdump -nr /tmp/c.pcap 'tcp[tcpflags] & tcp-rst != 0'
# Handshakes & resets only:
sudo tcpdump -ni any '(tcp[tcpflags] & (tcp-syn|tcp-rst)) != 0'
# Retransmit hunting is easier in tshark:
tshark -nr /tmp/c.pcap -Y tcp.analysis.retransmission
tshark -nr /tmp/c.pcap -q -z conv,tcp
tshark -nr /tmp/c.pcap -q -z io,phs             # protocol hierarchy
tshark -nr /tmp/c.pcap -q -z expert
# HTTP requests with timing:
tshark -nr /tmp/c.pcap -Y http.request -T fields -e ip.dst -e http.host -e http.request.uri
# DNS:
sudo tcpdump -ni any -s0 'udp port 53' -w /tmp/dns.pcap
# ARP / NDP:
sudo tcpdump -ni eth0 arp
sudo tcpdump -ni eth0 'icmp6 && ip6[40] >= 133 && ip6[40] <= 137'
# VXLAN / GENEVE (decodes inner packet in Wireshark):
sudo tcpdump -ni eth0 'udp port 4789 or udp port 6081'
# Fragments only:
sudo tcpdump -ni any 'ip[6:2] & 0x3fff != 0'
# By TCP payload prefix ("GET "):
sudo tcpdump -ni any 'tcp[((tcp[12:1] & 0xf0) >> 2):4] = 0x47455420'
```

### D4. Diagnostics

```bash
mtr -rwzbc100 <host>                   # loss/latency per hop + ASN, report mode
traceroute -n <host> ; traceroute -T -p 443 -n <host> ; tracepath <host>
ping -c5 <host> ; ping -M do -s 1472 <host>       # reachability ; PMTU probe
nstat -az                              # every kernel net counter
nstat ; sleep 10 ; nstat               # deltas over 10s
cat /proc/net/snmp /proc/net/netstat /proc/net/sockstat
ethtool eth0 ; ethtool -S eth0 ; ethtool -k eth0 ; ethtool -g eth0 ; ethtool -l eth0
tc -s qdisc show dev eth0 ; tc -s class show dev eth0
conntrack -L ; conntrack -S ; sysctl net.netfilter.nf_conntrack_count
iperf3 -s   /   iperf3 -c <host> -t30 -P8 -C bbr
curl -w 'dns %{time_namelookup} conn %{time_connect} tls %{time_appconnect} ttfb %{time_starttransfer} total %{time_total}\n' -o /dev/null -s https://<host>
dig +trace <name> ; dig @1.1.1.1 <name> ; dig -x <ip> ; dig <name> HTTPS
openssl s_client -connect <host>:443 -servername <host> -brief </dev/null
ss -tie dst <host>                      # the single most useful per-connection view
```

### D5. Sysctl inspection

```bash
sysctl -a --pattern 'net\.(ipv4|core|ipv6)\.' | less
sysctl net.ipv4.tcp_congestion_control net.core.default_qdisc
sysctl net.ipv4.tcp_rmem net.ipv4.tcp_wmem net.core.rmem_max net.core.wmem_max
sysctl net.ipv4.ip_local_port_range net.core.somaxconn net.ipv4.tcp_max_syn_backlog
sysctl net.ipv4.tcp_syncookies net.ipv4.tcp_tw_reuse net.ipv4.tcp_slow_start_after_idle
sysctl net.netfilter.nf_conntrack_max net.netfilter.nf_conntrack_count
```

---

# Appendix E — Glossary

**AIMD** — Additive Increase, Multiplicative Decrease. TCP's congestion-window
control law; source of fairness and the "sawtooth."

**ALG** — Application Layer Gateway. NAT helper that rewrites addresses embedded
in payloads (FTP, SIP). Often buggy.

**Anycast** — one IP prefix announced from many locations; BGP routes each client
to the nearest.

**ARP** — Address Resolution Protocol. Maps an IPv4 address to a MAC on the local
link. IPv6 uses NDP.

**AQM** — Active Queue Management (CoDel, PIE, fq_codel, CAKE). Keeps bottleneck
queues short to fight bufferbloat.

**AS / ASN** — Autonomous System / its number. A network with one routing policy;
the unit BGP connects.

**BDP** — Bandwidth-Delay Product = bandwidth × RTT. Bytes "in flight" needed to
fill a path; the target window size.

**BGP** — Border Gateway Protocol. Inter-AS path-vector routing; runs the
internet's route table over TCP/179.

**BBR** — Bottleneck Bandwidth and RTT. Model-based congestion control that paces
to measured bandwidth and largely ignores loss.

**Bufferbloat** — excessive latency from oversized, unmanaged buffers at a
bottleneck being kept full by loss-based TCP.

**CGNAT** — Carrier-Grade NAT (RFC 6598, `100.64.0.0/10`). A second NAT layer at
the ISP; breaks port forwarding and IP-based identity.

**CIDR** — Classless Inter-Domain Routing. Explicit prefix lengths
(`10.0.0.0/8`); replaced address classes in 1993.

**conntrack** — the Linux connection-tracking table used by stateful firewalling
and NAT.

**CUBIC** — Linux's default congestion control; window grows as a cubic function
of time since the last loss.

**cwnd / rwnd / ssthresh** — congestion window / receive (advertised) window /
slow-start threshold.

**DHCP** — Dynamic Host Configuration Protocol. DORA exchange assigns IP, mask,
gateway, DNS, lease.

**DNS** — Domain Name System. Hierarchical, cached name→record database
(A/AAAA/CNAME/MX/NS/SOA/TXT/SRV/HTTPS/…).

**DSCP** — Differentiated Services Code Point. 6-bit QoS class in the IP header
(EF, AFxx, CS0…).

**ECMP** — Equal-Cost Multi-Path. Per-flow hashing across equal-cost next hops;
caps a single flow to one path.

**ECN** — Explicit Congestion Notification. Routers mark (not drop) packets to
signal congestion; endpoints react without loss.

**EUI-64** — deriving a 64-bit IPv6 interface ID from a 48-bit MAC (now usually
replaced by random privacy addresses).

**FIN / RST / SYN / ACK** — TCP flags: finish (graceful), reset (abort),
synchronize (setup), acknowledge.

**GRO / GSO / TSO / LRO** — Generic Receive / Generic Segmentation / TCP
Segmentation / Large Receive Offload. Batch packets to cut per-packet CPU.

**Happy Eyeballs** — algorithm that races IPv6 and IPv4 connection attempts so a
broken family costs ~250 ms, not a timeout.

**HOL blocking** — Head-of-Line blocking. One stalled item delays everything
queued behind it (TCP under HTTP/2; fixed by QUIC streams).

**ICMP / ICMPv6** — control/error messaging for IP (echo, unreachable, time
exceeded, packet-too-big). ICMPv6 also carries NDP.

**ISN** — Initial Sequence Number. Randomized per connection to resist off-path
injection.

**MSS / MTU** — Maximum Segment Size (TCP payload per segment) / Maximum
Transmission Unit (L3 payload per frame). `MSS ≈ MTU − IP − TCP headers`.

**Mathis equation** — `throughput ≈ MSS / (RTT · √loss)` for a loss-based TCP
flow.

**NAT / PAT** — Network Address (and Port) Translation. Many private hosts share
one public IP by rewriting the 4-tuple.

**NDP** — Neighbor Discovery Protocol (ICMPv6). Replaces ARP + router discovery +
redirects for IPv6 (types 133–137).

**PMTUD** — Path MTU Discovery. Learn the smallest link on a path via "packet too
big" ICMP; broken when ICMP is filtered (black hole).

**QUIC** — a reliable, multiplexed, always-encrypted transport over UDP;
foundation of HTTP/3. Independent streams, connection migration, 0-RTT.

**RACK-TLP** — Recent ACKnowledgment loss detection + Tail Loss Probe.
Time-based, reordering-robust replacement for dup-ACK counting.

**RPKI / ROA / ROV** — Resource PKI / Route Origin Authorization / Route Origin
Validation. Cryptographic defense against BGP prefix hijacks.

**RTO / SRTT / RTTVAR** — Retransmission Timeout and the smoothed RTT / RTT
variance estimators that set it.

**SACK / DSACK** — Selective ACK (report received ranges above a gap so the
sender retransmits only holes) / Duplicate SACK (report data received twice).

**SLAAC** — StateLess Address AutoConfiguration. IPv6 hosts self-assign addresses
from a Router Advertisement prefix.

**SNI** — Server Name Indication. The hostname in the TLS ClientHello so one IP
serves many certificates (plaintext unless ECH).

**STUN / TURN / ICE** — discover your NAT mapping / relay when direct fails / the
algorithm that tries all candidate paths. NAT traversal for P2P and WebRTC.

**TCB** — Transmission Control Block. Per-connection kernel state (tuple, seq
numbers, windows, timers, buffers).

**TIME_WAIT** — post-close state (2×MSL) on the side that sent the last ACK;
absorbs stragglers and prevents tuple confusion.

**TTL / Hop Limit** — per-hop decremented counter; at 0 the packet is dropped and
ICMP "time exceeded" returned (powers traceroute).

**uRPF / BCP 38** — Unicast Reverse Path Forwarding / the practice of ingress
filtering spoofed source addresses.

**VLAN / VXLAN** — Virtual LAN (802.1Q, 12-bit ID, one switch) / Virtual
Extensible LAN (24-bit VNI, MAC-in-UDP overlay across an IP underlay).

**Window scaling** — TCP option multiplying the 16-bit window by up to 2^14, so
windows can exceed 64 KB (required for high-BDP paths).

---

## Further reading

- **RFC 9293** — TCP (the 2022 consolidated spec, replaces RFC 793).
- **RFC 1122 / 1123** — Host Requirements (the "4-layer model").
- **RFC 5681** — TCP Congestion Control; **RFC 6582** NewReno; **RFC 2018** SACK;
  **RFC 6298** RTO; **RFC 7323** window scaling & timestamps; **RFC 8985**
  RACK-TLP; **RFC 3168** ECN; **RFC 9438** CUBIC.
- **RFC 8200** — IPv6; **RFC 4861** NDP; **RFC 4862** SLAAC; **RFC 8981** privacy
  addresses.
- **RFC 826** ARP; **RFC 792 / 4443** ICMP / ICMPv6; **RFC 1191 / 8899** PMTUD /
  PLPMTUD.
- **RFC 2131** DHCP; **RFC 1034 / 1035** DNS; **RFC 6891** EDNS0; **RFC 7858**
  DoT; **RFC 8484** DoH.
- **RFC 9000–9002** QUIC; **RFC 9114** HTTP/3; **RFC 9113** HTTP/2.
- **RFC 4787 / 5245 / 8445** NAT behavior / ICE; **RFC 8489** STUN; **RFC 8656**
  TURN; **RFC 3550** RTP.
- Books: *TCP/IP Illustrated, Vol. 1* (Fall & Stevens); *Computer Networks: A
  Systems Approach* (Peterson & Davie); *High Performance Browser Networking*
  (Grigorik, free online); *BPF Performance Tools* (Gregg).
- Sites: the `iproute2` and `ss` man pages; Cloudflare, APNIC, and Julia Evans
  networking blog posts; the Bufferbloat project (`bufferbloat.net`).

---

*End of guide. Keep it close; update the war stories with your own.*
