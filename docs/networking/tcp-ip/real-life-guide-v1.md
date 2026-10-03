# Networking from Zero — A Beginner's Guide to TCP/IP (v1)

> **You do not need to know anything about networking to read this.** You need to
> be able to open a terminal and type commands. Everything else is explained the
> first time it appears.
>
> There is a companion reference guide (`real-life-example.md`) that is denser and
> assumes more. Use that *after* this one, as a lookup manual.

---

> **The series:** 1 OS → 2 Networking → 3 Security → 4 HTTPS walkthrough, with Go alongside.
>
> **You are here: step 2b, Networking: in depth.** ← Previous: [The OSI Model, One Click at a Time](../real-life-example-osi.md). Next: [Security from Zero](../../security/real-life-guide.md) →
>
> [The full series map](#0-7-the-series-os-networking-security-https).

---

## What you will be able to do at the end

1. Explain what happens, step by step, between pressing Enter on a URL and seeing
   a page — with no hand-waving.
2. Read an IP address and subnet mask and say exactly which addresses are on that
   network. (Subnetting stops being scary.)
3. Use `ping`, `traceroute`, `dig`, `curl`, `ss`, and `tcpdump` deliberately
   rather than hopefully.
4. Look at a packet capture and describe what the two machines are saying to each
   other.
5. Diagnose the five most common network problems by their fingerprints, instead
   of guessing.
6. Explain TCP's handshake, reliability, flow control, and congestion control —
   and why a connection is slow.
7. Build a small production-like lab with DNS, TLS, load balancing, firewall
   rules, containers/Kubernetes-style networking, monitoring, and packet
   captures.
8. Write your own network tools in Go: a DNS client, a pcap decoder, a
   TLS inspector, a load balancer with draining, a blackbox prober, and more
   (§0.8 lists all 23 labs).
9. Know what to read/watch next.

---

## Contents

**Part 0 — Start here** *(read this first)*
- 0.1 Who this is for · 0.2 **The learning flow** · 0.3 The roadmap ·
  0.4 Set up your machine · 0.5 The tiny bit of background you need ·
  0.6 How to not get stuck · 0.7 **The series: OS → networking → security → HTTPS** ·
  0.8 **The Go labs: setup and index**

**Part 1 — What is a network, really?**
1. The problem networks solve
2. Layers: the postal-service idea
3. The four layers of TCP/IP
4. Our running example: what happens when you open a website

**Part 2 — Your local network (the link layer)**
5. MAC addresses and Ethernet frames
6. ARP: finding your neighbour
7. Switches, and what "broadcast domain" means
8. Wi-Fi: same idea, different physics

**Part 3 — IP addresses and getting around the world**
9. What an IP address actually is
10. **Subnet masks: the part everyone finds hard**
11. Subnetting practice (with answers)
12. Private addresses and NAT
13. Routing: how a packet finds its way
14. ICMP: what `ping` and `traceroute` really do
15. MTU: the maximum size of a packet
16. IPv6: why it exists and how it differs

**Part 4 — Getting connected automatically**
17. DHCP: how your device gets an address
18. DNS: turning names into addresses

**Part 5 — Ports and the transport layer**
19. Ports and sockets: which program gets the data?
20. UDP: fire and forget
21. **TCP part 1: the three-way handshake**
22. TCP part 2: how it never loses your data
23. TCP part 3: closing a connection, and the states

**Part 6 — Making TCP fast (and why yours isn't)**
24. Flow control: the sliding window
25. Congestion control: sharing the road
26. The 40-millisecond mystery (Nagle meets delayed ACK)
27. Latency, bandwidth, and why distance costs money

**Part 7 — The application layer**
28. HTTP: how the web actually talks
29. TLS: how it gets encrypted
30. HTTP/2 and HTTP/3

**Part 8 — Seeing and fixing the network**
31. `tcpdump`: watching packets go by
32. Wireshark: reading a conversation
33. `ss`: what your machine's connections are doing
34. **A troubleshooting method that works**
35. The five classic problems and their fingerprints

**Part 9 — Modern networking**
36. Load balancers, CDNs, and anycast
37. Cloud and container networking
38. Security basics

**Part 10 — HTTP in depth**
39. How HTTP messages are framed (and how it goes wrong)
40. **Caching: the biggest performance lever there is**
41. Cookies: how the web remembers you
42. CORS: why the browser blocked your request
43. WebSockets and Server-Sent Events
44. Compression, range requests, and proxies

**Part 11 — Watching your network (monitoring)**
45. What to measure, and why averages lie
46. Where the numbers come from: counters, flows, and probes
47. Alerts that wake you for the right reasons

**Part 12 — How routers learn routes**
48. Two ways to share directions
49. OSPF: finding your way inside one network
50. BGP: how the whole internet agrees on directions
51. When a router dies: VRRP and friends

**Part 13 — Professional networking: operations, debugging, and design**
52. Host networking internals: namespaces, veth, bridges, conntrack
53. Firewalls in the real world: host firewalls, cloud security groups, NACLs
54. VPNs and tunnels: WireGuard, IPsec, GRE, overlays, and MTU traps
55. DNS operations: authoritative DNS, delegation, split-horizon, outages
56. TLS and certificate operations: expiry, chains, SNI, rotation
57. Production load balancing: L4 vs L7, health checks, draining, retries
58. Kubernetes networking: Services, Ingress, NetworkPolicy, CNI
59. Network performance engineering: QoS, packet loss, shaping, saturation
60. Incident response: runbooks, packet evidence, and escalation
61. Network design/refactoring: migrations, diagrams, IPAM, change safety
62. Capstone lab: build, break, observe, and repair a production-like network

**Part 14 — Where to go next**
63. Honest gaps and a six-month plan
64. Project ideas by level
65. **The series map: where every topic lives**

**Appendices**
- A. Glossary (plain language)
- B. Command cheat sheet
- C. Subnetting cheat sheet + practice answers
- D. Ports and protocols reference
- E. Reading list (books, blogs, videos, courses) + how facts were checked
- F. Answers to "Check yourself"

---

# Part 0 — Start here: how to use this guide

Read this part. It's short and it will save you weeks.

## 0.1 Who this guide is for

| You are... | Works for you? |
|---|---|
| A developer who has never thought about networking | **Yes** — this is the target reader |
| A student starting a networking or sysadmin course | Yes |
| Someone studying for CompTIA Network+ / CCNA | Yes, as the "why" behind the syllabus |
| A support/SRE person who wants to stop guessing | Yes |
| A network engineer wanting reference material | Use the companion guide instead |

**Assumed knowledge:** you can open a terminal and run commands. You know that a
website has an address. That's it.

**Not assumed:** the OSI model, binary, hexadecimal, IP addresses, ports,
protocols, or any acronym at all. Every term is defined the first time it's used.

---

## 0.2 The learning flow (use this for every chapter)

Reading about networking does not teach you networking. **Running commands and
watching real packets does.** Every chapter is built to support this loop.

```
   +------------------------------------------------------------------+
   |                                                                  |
   |   1. READ        2. EXPLAIN      3. PRACTICE     4. CHECK        |
   |   (10-20 min)    (5 min)         (15-45 min)     (5 min)         |
   |                                                                  |
   |   Read the  -->  Close the  -->  RUN the     --> Answer the      |
   |   chapter        guide and       commands on     "Check          |
   |   once,          say it in       YOUR OWN        yourself"       |
   |   slowly.        your own        machine.        questions from  |
   |                  words.          Change things.  memory.         |
   |                       |                              |           |
   |                       | can't explain it?            | got one   |
   |                       +---> re-read that section      | wrong?    |
   |                                                       |           |
   |                       +-------------------------------+           |
   |                       v                                           |
   |   5. CONNECT: write ONE line in your notes:                       |
   |      "X exists because Y was a problem, and it works by Z."       |
   |                                                                  |
   +------------------------------------------------------------------+
                                   |
                                   v
                 every ~5 chapters: CAPTURE something real
                 (open tcpdump/Wireshark and watch it happen)
```

### Why each step matters

- **Read once, slowly.** Don't re-read a paragraph three times. Push on — the next
  paragraph usually answers your question.
- **Explain out loud.** The highest-value step. If you can't explain it simply,
  you don't understand it yet, and you find that out in 5 minutes rather than
  5 weeks.
- **Practice on your own machine.** Networking is the most *observable* subject in
  computing — you can literally watch every byte. A concept you have seen in a
  packet capture is a concept you own forever.
- **Check yourself from memory.** Don't look back first.
- **Connect.** One line per chapter in your own `notes.md`. After 40 chapters you
  have a map of networking that *you* wrote.

### The rules that keep you moving

1. **Time-box confusion to 20 minutes.** Write the exact question down, mark it
   `TODO`, move on. It will usually make sense two chapters later.
2. **Never skip the Practice sections.** Reading Part 5 without watching a real
   TCP handshake in `tcpdump` is how people "learn TCP" three times and still
   can't explain it.
3. **Don't chase every link.** Further Reading is for later or when stuck.
4. **One pass, then depth.** Get through the whole guide shallowly before going
   deep on any part.
5. **Keep a `notes.md`** with: your one-line summaries, your `TODO` questions, and
   every command that surprised you.

---

## 0.3 The roadmap

Three routes. All use the same chapters.

### Route A — "I want to understand how the internet works" (2 weeks, reading + light practice, ~1 h/day)

| Day | Read | You'll be able to |
|---|---|---|
| 1 | Part 0, Part 1 | Describe the layer model and the URL journey |
| 2 | Ch 5–8 | Explain MAC addresses, ARP, switches |
| 3–4 | Ch 9–11 | **Read any IP/mask and do subnetting** |
| 5 | Ch 12–14 | NAT, routing, ping/traceroute |
| 6 | Ch 15–16 | MTU, IPv6 |
| 7 | Ch 17–18 | DHCP and DNS |
| 8 | Ch 19–20 | Ports, sockets, UDP |
| 9–10 | Ch 21–23 | **TCP end to end** |
| 11 | Ch 24–27 | Why connections are slow |
| 12 | Ch 28–30 | HTTP, TLS, HTTP/3 |
| 13 | Ch 31–33 | tcpdump, Wireshark, ss |
| 14 | Ch 34–35 | **Troubleshoot like a professional** |

### Route B — "I need this for my job" (6 weeks, ~1–2 h/day)

Weeks 1–2: Route A, **doing every Practice section**.
Week 3: Part 8 (tools) + capture your own traffic daily.
Week 4: Part 3 again, deeply — subnet a real network on paper; Project 1.
Week 5: Part 5–6 again — capture and annotate a full TCP connection; Project 2.
Week 6: Part 9 + Project 3 (break something on purpose and diagnose it).

### Route C — "I want to go deep / change careers" (12 weeks, ~2 h/day)

| Week | Focus | Deliverable |
|---|---|---|
| 1 | Part 0–1 + tools installed | `notes.md` started, first capture taken |
| 2 | Part 2 (link layer) | Annotated capture of ARP + an Ethernet frame |
| 3 | Ch 9–11 | 30 subnetting problems solved correctly, on paper |
| 4 | Ch 12–16 | A network diagram of your own home/office network |
| 5 | Ch 17–18 | A full DNS resolution traced with `dig +trace` |
| 6 | Ch 19–23 | **A packet-by-packet writeup of one TCP connection** |
| 7 | Ch 24–27 | A measured throughput/latency experiment with `tc` |
| 8 | Ch 28–30 | A TLS handshake decoded in Wireshark |
| 9 | Ch 31–33 | Your own tcpdump/Wireshark cheat sheet |
| 10 | Ch 34–35 | Diagnose 5 deliberately-broken scenarios |
| 11 | Ch 36–38 | Build a small lab (containers/VMs) with routing between them |
| 12 | Consolidate | A blog post teaching one topic; read the companion guide |

Then: CompTIA Network+ or CCNA if you want the credential, *Computer Networking:
A Top-Down Approach* for theory, and the RFCs for the truth.

**Beyond all three routes:** each is time-boxed and stops around Part 8 or 9.
Parts 10–13 continue the same material at the same depth whenever you have
more time — HTTP internals (caching, cookies, CORS, streaming), monitoring at
scale, how routers actually learn routes (OSPF/BGP), and the guide's own
honest wrap-up. Nothing about them requires a fixed schedule; read them
chapter by chapter as the topic becomes relevant to what you're working on.

---

## 0.4 Set up your machine (20 minutes)

Networking tools are mostly already installed. Let's check and fill the gaps.

### Step 1 — Check what you have

Run each of these. You want a version number, not "command not found".

```bash
ping -V 2>/dev/null || ping -c1 127.0.0.1   # reachability test
traceroute --version 2>/dev/null || traceroute -V
dig -v                                       # DNS lookups
curl --version                               # HTTP client
ss -V 2>/dev/null || netstat --version       # socket list (Linux: ss)
tcpdump --version                            # packet capture
nc -h 2>&1 | head -2                         # netcat: manual connections
ip -V 2>/dev/null || ifconfig -h 2>&1|head -1
```

### Step 2 — Install what's missing

**macOS** (most tools ship built in):
```bash
# Homebrew: https://brew.sh  (install it first if you don't have it)
brew install iproute2mac mtr bind    # 'ip' command, mtr, and dig/host
brew install --cask wireshark        # GUI packet analyser
```
`ping`, `traceroute`, `netstat`, `tcpdump`, `nc`, `curl`, `dig` are already there.
macOS has `netstat` rather than `ss`.

**Linux (Debian/Ubuntu):**
```bash
sudo apt update
sudo apt install -y iproute2 iputils-ping traceroute mtr-tiny dnsutils \
                    curl netcat-openbsd tcpdump wireshark whois
```

**Windows:** the friendliest path is **WSL2** (Windows Subsystem for Linux) —
install Ubuntu from the Microsoft Store, then use the Linux commands above.
Native Windows equivalents also exist: `ping`, `tracert`, `nslookup`,
`netstat -ano`, `curl`, and Wireshark.

### Step 3 — Let yourself capture packets

Packet capture needs privileges.

**Linux** — either use `sudo tcpdump ...`, or grant your user permission once:
```bash
sudo usermod -aG wireshark $USER      # then LOG OUT and back in
```

**macOS** — `sudo tcpdump ...` works. Wireshark will prompt to install ChmodBPF.

**Test it:**
```bash
sudo tcpdump -i any -c 5 -n
```
You should see five lines of traffic scroll past. If you see
`tcpdump: <interface>: You don't have permission to capture`, fix privileges above.

### Step 4 — Know your own network

Run this and **save the output in `notes.md`** — you'll refer to it constantly:

```bash
# Linux
ip -br addr ; echo ---- ; ip route ; echo ---- ; cat /etc/resolv.conf

# macOS
ifconfig | grep -A3 '^en' ; echo ---- ; netstat -rn | head -20
echo ---- ; scutil --dns | grep nameserver | head
```

You are looking for four things:
1. **Your IP address** (something like `192.168.1.50`)
2. **Your subnet mask / prefix** (`/24` or `255.255.255.0`)
3. **Your default gateway** (usually `192.168.1.1`)
4. **Your DNS server(s)**

Don't worry about what they mean yet — Parts 3 and 4 explain each one. But having
your own numbers in front of you makes every example concrete.

### A safety note

Everything in this guide is done on **your own machine and your own network**.
Capturing traffic on networks you don't control, or scanning machines you don't
own, can be illegal and is definitely rude. Practise at home.

---

## 0.5 The tiny bit of background you need

Four ideas. Skim now; each is re-explained where it's used.

### Idea 1 — Bits, bytes, and why 255 keeps appearing

A **bit** is a 0 or a 1. A **byte** is 8 bits.

With 8 bits you can count from `00000000` to `11111111`:

```
  binary     decimal
  00000000  =   0
  00000001  =   1
  00000010  =   2
  ...
  11111111  = 255      <-- the biggest 8-bit number
```

That's why IP addresses are made of numbers **0 to 255**, and why you'll see
`255` constantly in subnet masks. There is nothing magic about it — it's just
"all 8 bits switched on".

The place values of the 8 bits (memorise these; subnetting becomes easy):

```
   128   64   32   16    8    4    2    1
    |    |    |    |     |    |    |    |
    1    1    0    0     0    0    0    0   = 128 + 64 = 192
    1    1    1    1     1    1    1    1   = 255
    1    0    0    0     0    0    0    0   = 128
```

### Idea 2 — Hexadecimal, in 30 seconds

**Hex** counts in 16s instead of 10s, using `0-9` then `a-f`:

```
  hex:      0 1 2 3 4 5 6 7 8 9  a  b  c  d  e  f
  decimal:  0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
```

Two hex digits = exactly one byte (`ff` = 255). That's why hardware addresses
look like `a4:83:e7:2b:19:0c` — six bytes, written as six pairs of hex digits.
You rarely need to convert; you just need to not be scared of it.

### Idea 3 — Client and server

- A **server** is a program that waits for connections and answers them.
- A **client** is a program that starts a connection and asks for something.

Your browser is a client. `example.com`'s web server is a server. The same
physical machine can run both. That's the entire distinction.

### Idea 4 — A protocol is just an agreed format

A **protocol** is a set of rules two programs agree on: what the message looks
like, what order things happen in, what each field means. Like agreeing that a
letter has an address on the outside and a greeting at the top.

"TCP/IP" is a *family* of such agreements, which we'll meet one at a time.

### That's genuinely all

No maths beyond counting. If you can read the place-value table in Idea 1, you can
do everything in this guide, including subnetting.

### Further reading (only if you want more background)

- **Video (~7 min):** "Binary Numbers" — search Computerphile or Khan Academy.
  Only if Idea 1 felt shaky.
- **Article:** "Hexadecimal" on Wikipedia — the first section is enough.
- **Book (delightful, non-technical):** *Code: The Hidden Language of Computer
  Hardware and Software* — Charles Petzold. Explains bits, bytes, and how
  computers work from first principles. Not required, but wonderful.

---

## 0.6 How to not get stuck

| Situation | Do this |
|---|---|
| "I don't understand this paragraph" | Read the next two. Usually resolved. |
| "Still lost after 20 minutes" | Write the question in `notes.md`, mark `TODO`, move on. Return in 2 chapters. |
| "The command doesn't work" | Read the **last line** of the error. Check if it needs `sudo`. Check your OS (macOS vs Linux differ). |
| "I understand but can't explain it" | You don't understand it yet. Re-read only the "How it works" section, then try again. |
| "There are too many acronyms" | Appendix A is a plain-language glossary. Bookmark it now. |
| "This feels too basic" | Skip to Practice. If you can do it, move on. |
| "Should I memorise port numbers?" | No. Appendix D exists. You'll absorb the common ones naturally. |
| "Do I need to learn the OSI 7-layer model?" | For exams, yes. For understanding, the 4-layer model in Chapter 3 is better. We cover both. |

### A note on honesty

This guide marks where it is simplifying:

> **Simplified:** the real story has more detail; here's the direction it goes.

> **Debated:** engineers genuinely disagree; here are the sides.

Facts, numbers, and RFC references were checked against primary sources —
see Appendix E.

---

## 0.7 The series: OS → networking → security → HTTPS

This guide is one of seven, designed to be read as **one course** in this order:

```
   1  OS & Linux  ──▶  2  Networking  ──▶  3  Security  ──▶  4  HTTPS walkthrough
                         2a OSI map          3a From Zero        (one request through
                         2b TCP/IP           3b In Depth          every guide, in Go)
   ════════════════════  Go, alongside every step  ════════════════════
```

| Step | Guide | What it gives you | Hands-on |
|---|---|---|---|
| 1 | [Operating Systems, Linux, and Containers](../../os-linux/real-life-os-guide.md) | the machine every request starts and ends on: processes, memory, files, sockets, containers, Kubernetes | 9 Go labs ([Part 20](../../os-linux/real-life-os-guide.md#part-20-systems-programming-in-go-the-os-from-inside-a-program)) |
| 2a | [The OSI Model, One Click at a Time](../real-life-example-osi.md) | one click followed through all seven layers; the map for everything after it | concept map, 1 hour |
| 2b | **Networking from Zero (TCP/IP)** ← you are here | how machines talk: addressing, routing, TCP, TLS, packet capture, network operations | 23 Go labs ([§0.8](#0-8-the-go-labs-build-the-network-tools-yourself)) |
| 3a | [Security from Zero](../../security/real-life-guide.md) | encoding vs hashing vs encryption, TLS, PKI, SSH, the OWASP Top 10, threat modelling | OpenSSL labs and a vulnerable app to break |
| 3b | [Security Engineering in Depth](../../security/real-life-security-guide-v1.md) | cloud and Kubernetes security, distributed authorization, advanced web attacks, data protection, detection | 10 Go labs + a `govulncheck` exercise ([§0.8](../../security/real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)) |
| 4 | [The HTTPS Request Lifecycle](../../v2-https/real-life-guide-v1.md) | one request end to end, then the server side built in Go; Chapter 25 traces one request through every guide | 13 Go labs ([§0.6](../../v2-https/real-life-guide-v1.md#0-6-the-go-labs-build-the-lifecycle-yourself), [Ch 25](../../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)) |
| ∥ | [Go — The Complete Field Guide](../../Golang/real-life-golang-guide.md) | the language behind every lab, plus the 120-day plan's multi-week projects | [120-day plan](../../Golang/golang-90-day-plan.md) |

**Why this order.** Every network connection is a file descriptor owned by
a process, so the OS comes first. Networking comes next, because every attack and
defence in the security guides assumes you can follow a packet. Security comes
third, because it needs both. The HTTPS walkthrough comes last because a
single HTTPS request uses all of it, and [its final chapter](../../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide)
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

## 0.8 The Go labs: build the network tools yourself

Reading about a protocol teaches you its shape. Writing a program that speaks
it shows you the details: the bytes, the timeouts, the error codes, the
states. From Part 3 onwards, many chapters end with a **"Build it in Go"**
section: a small, complete program (usually 50–200 lines) that turns the
chapter into something you can run, break, and extend.

### Why Go

Go's standard library contains a complete network stack API (`net`,
`net/netip`, `crypto/tls`, `net/http`) with no extra dependencies. Its
goroutines make "one connection per goroutine" natural. And it's the language
much of today's network infrastructure is written in: Kubernetes,
Docker, Caddy, Traefik, CoreDNS, Consul, Prometheus, and most Envoy control
planes. You don't need to know Go well to start. Each lab explains what's
specific to Go, and the [Go guide](../../Golang/real-life-golang-guide.md) covers the language itself.

### Setup (5 minutes)

```bash
# Install Go 1.22 or newer (the labs were tested with Go 1.26):
#   macOS:  brew install go        Linux: https://go.dev/doc/install
go version

# One module holds every lab; each lab is a folder with a main.go
mkdir -p ~/netlabs && cd ~/netlabs
go mod init netlabs
mkdir dnsquery            # then save the lab's code as dnsquery/main.go
go run ./dnsquery example.com
```

Labs marked **(Linux)** use Linux-only kernel features. On macOS, run them in
a container. This works even on Docker Desktop:

```bash
docker run --rm -it --cap-add NET_ADMIN -v "$PWD":/labs -w /labs golang:1.26 bash
apt-get update && apt-get install -y iproute2 tcpdump   # ss, tc, tcpdump inside the container
```

### What Go's `net` package does underneath

Every lab uses a few calls. Knowing what each one asks the kernel to do
connects the code to the chapters:

| Go call | System calls underneath (Linux) | What happens on the wire / in the kernel |
|---|---|---|
| `net.Dial("tcp", addr)` | `socket()`, non-blocking `connect()`, then wait in epoll | SYN sent; returns after SYN-ACK + ACK (Chapter 21) |
| `net.Listen("tcp", addr)` | `socket()`, `setsockopt(SO_REUSEADDR)`, `bind()`, `listen(somaxconn)` | Socket enters LISTEN; kernel now completes handshakes by itself |
| `ln.Accept()` | `accept4()` | Takes one finished connection off the accept queue (Chapter 52) |
| `conn.Read(buf)` | `read()`; if no data, the goroutine parks in the **netpoller** (epoll/kqueue) | Drains the socket's receive buffer, which opens the window (Chapter 24) |
| `conn.Write(buf)` | `write()`; parks if the send buffer is full | Copies into the send buffer; TCP sends when cwnd/rwnd allow |
| `conn.SetDeadline(t)` | none (a runtime timer) | Read/Write fail with a timeout error after `t` |
| `conn.Close()` | `close()` | FIN, or **RST** if unread data was still waiting |
| `net.Dial("udp", addr)` | `socket()`, `connect()` | **Nothing is sent.** Only a route and a source address are chosen |

Three Go defaults that matter later: TCP connections have **`TCP_NODELAY`
on** (Nagle off, Chapter 26), **TCP keepalives on** (every 15 s), and
listeners set **`SO_REUSEADDR`**, so a restarted server can rebind its port
even while old connections sit in TIME_WAIT (Chapter 23).

To see the system calls any lab makes, build it and trace it:

```bash
go build -o /tmp/lab ./dialcheck
strace -f -e trace=network /tmp/lab example.com:443     # Linux
```

### Index of the Go labs in this guide

| Lab | Chapter | What you build |
|---|---|---|
| `subnet` | 10, 11, 61 | Subnet calculator, subnet splitter, overlap checker, IPAM allocator |
| `whoami` | 12 | Your private vs public address: NAT made visible |
| `dnsquery` | 18 | A DNS client from raw bytes, with message compression |
| `bindlab` | 19 | Bind conflicts, loopback vs wildcard, `SO_REUSEPORT` |
| `udpping` | 20 | UDP echo + loss/reorder/jitter meter |
| `dialcheck` | 21 | Classifies connected / refused / timed out / DNS failure |
| `stream` | 22 | Proof that TCP has no message boundaries, and length-prefix framing |
| `closewait` | 23 | A CLOSE_WAIT leak, half-close with `CloseWrite` |
| `zerowindow` | 24 | A slow reader that drives the window to zero |
| `nagle` | 26 | The 40 ms Nagle + delayed-ACK stall, on demand |
| `rawhttp` | 28 | HTTP/1.1 by hand over TCP/TLS, decoding chunked framing |
| `tlsinspect` | 29, 56 | A scriptable `openssl s_client` with failure diagnosis |
| `pcapread` | 31–32 | A pcap decoder: Ethernet/SLL → IP → TCP/UDP |
| `tcpinfo` (Linux) | 33 | `ss -ti` from inside your program: rtt, cwnd, retransmits |
| `netcheck` | 34 | The troubleshooting ladder as one command |
| `prober` | 46–47 | A blackbox exporter with a Prometheus latency histogram |
| `acceptq` | 52 | Overflowing a listen queue on purpose |
| `policytest` | 53 | Firewall policy as executable tests |
| `dnsdiff` | 55 | Compare answers across resolvers |
| `certreload` | 56 | Zero-downtime certificate rotation |
| `lbdrain` | 57 | Health-checked load balancer + gracefully draining backends |
| `proxyproto` | 57 | PROXY protocol parser: the real client IP behind an L4 load balancer |
| `bdp` | 59 | The bandwidth-delay product, measured |

Larger projects that build on these, such as an L4 load balancer, a TLS
reverse proxy, a WAF, and a DDoS-resistant edge proxy, are in the
[Go 120-day plan](../../Golang/golang-90-day-plan.md), Weeks 5–17. The individual chapters link to them.

# Part 1 — What is a network, really?

Machine set up, route chosen — time to start. Before any protocol or packet,
there's one question worth answering slowly: what problem is a network
actually solving, and why did it end up built the way it is? Everything in
the following twelve parts is an answer to some version of that question.

## Chapter 1 — The problem networks solve

### In one sentence

A network's job is to get a chunk of data from a program on one machine to a
program on another machine, over links that are slow, unreliable, owned by
strangers, and constantly changing.

### The problem

Say you want to send a photo from your laptop to a friend's phone in another
country. Sounds simple. Now list what could go wrong:

- Your laptop and their phone have **no wire between them**. There are maybe
  fifteen machines in between, owned by six different companies.
- Those machines are **shared** with millions of other people doing the same
  thing right now.
- Any link along the way can **drop your data** silently — a cable is unplugged,
  a router is overloaded, a Wi-Fi signal is interfered with.
- Data can arrive **out of order**, or **twice**, or **corrupted**.
- The path can **change mid-transfer** because a router somewhere failed.
- The machines run **different operating systems** made by different companies
  decades apart.
- Nobody is in charge of the whole thing.

And yet it works, hundreds of billions of times a day. The rest of this guide is
how.

### The idea, in plain language

Break the giant problem into small problems and solve each one separately, with a
separate agreement (a **protocol**) for each.

| Small problem | The question it answers | Who solves it |
|---|---|---|
| **Framing** | Where does one message stop and the next begin on this wire? | Ethernet, Wi-Fi (Part 2) |
| **Local addressing** | Which physical device on *this* cable? | MAC addresses, ARP (Part 2) |
| **Global addressing** | Which machine on the *planet*? | IP addresses (Part 3) |
| **Routing** | Which way do I send it? | Routers, routing tables (Part 3) |
| **Which program?** | Machine has 50 programs — which one gets this? | Port numbers (Part 5) |
| **Reliability** | Did it arrive, complete and in order? | TCP (Part 5) |
| **Speed and fairness** | How fast can I send without breaking things? | TCP congestion control (Part 6) |
| **Names** | I only know `example.com`, not a number | DNS (Part 4) |

**"TCP/IP" is the name for this whole family of agreements.** It's not one
protocol; it's about a dozen working together, and this guide takes them one at a
time.

### Worked example: the same problem, in the physical world

Getting a parcel from your house to a friend abroad:

```
YOU                                                     YOUR FRIEND
 |                                                            ^
 | write a note                                               | reads note
 v                                                            |
 put it in an ENVELOPE with their NAME                  opens envelope
 |                                                            ^
 v                                                            |
 put envelope in a BOX with their FULL ADDRESS          opens box
 |                                                            ^
 v                                                            |
 hand box to the LOCAL COURIER                          local courier
 |                                                       delivers
 v                                                            ^
 local depot -> national hub -> plane -> foreign hub -> local depot
```

Notice three things that map exactly onto networking:

1. **Each layer only reads its own label.** The plane's loading crew doesn't read
   the note. The courier doesn't care what's in the envelope.
2. **The "next hop" changes constantly; the final address doesn't.** Every depot
   asks "who do I hand this to *next*?", but the address on the box never changes.
3. **You only had to know the local courier.** You didn't need to know the flight
   number or the foreign postal system.

That is the whole design of the internet. Chapter 2 makes it precise.

### Common confusions

- **"Is the internet the same as the web?"** No. The **internet** is the network
  (this guide). The **web** is one application that runs on it (HTTP, Part 7).
  Email, video calls, and games are other applications on the same internet.
- **"Who owns the internet?"** Nobody. It's tens of thousands of independent
  networks that agree to interconnect using these protocols. That's why the
  protocols matter so much — they're the only thing holding it together.

### Check yourself

1. Name three things that can go wrong when data crosses a network.
2. Why split the problem into separate protocols instead of one big one?
3. What's the difference between the internet and the web?

*(Answers: Appendix F.)*

### Further reading

- **Video (~9 min):** "How does the internet work?" — Code.org's series with
  Vint Cerf (who co-invented TCP/IP). Start here; it's genuinely excellent.
- **Article:** "How Does the Internet Work?" — Cloudflare Learning Center. Short,
  clear, no jargon.
- **Book (very readable):** *Tubes: A Journey to the Center of the Internet* —
  Andrew Blum. The *physical* internet: cables, buildings, oceans. Non-technical
  and fascinating.
- **Video (~13 min):** "The Internet: Packets, Routing & Reliability" — Code.org.

---

## Chapter 2 — Layers: the postal-service idea

### In one sentence

Networking is organised in **layers**, where each layer does one job and treats
everything above it as "cargo it doesn't need to understand".

### The problem

You could design a network as one giant program that handles everything: wires,
addressing, retransmission, encryption, web pages. People tried. It's
unmaintainable, and every new technology (Wi-Fi, fibre, 5G) would require
rewriting everything.

### The idea, in plain language

Split the work into stacked layers. Rules:

1. **Each layer only talks to the layer directly above and below it.**
2. **Each layer treats the layer above's data as opaque cargo.**
3. **Each layer can be swapped out** without touching the others.

That third rule is the payoff: Wi-Fi replaced Ethernet cables without changing a
single web browser. HTTP/3 was invented without upgrading a single router.

### How it actually works: wrapping

As your data goes **down** the layers, each layer wraps it in a **header** — a
small block of its own information, stuck on the front.

```
Your data:                                   [ "Hello" ]

Transport layer adds its header:      [ TCP hdr | "Hello" ]
   (says: which program on each end, and this is chunk #5)

Internet layer adds its header:  [ IP hdr | TCP hdr | "Hello" ]
   (says: from this machine, to that machine)

Link layer adds its header:  [ Eth hdr | IP hdr | TCP hdr | "Hello" ]
   (says: from this network card, to that network card, on this cable)

                             ---> onto the wire as electricity/light/radio
```

At the other end, the same thing happens in reverse — each layer peels off its
own header and hands the rest up. This is called **encapsulation** (going down)
and **decapsulation** (coming up).

The postal analogy again, precisely:

```
   NETWORKING                          POST
   ----------                          ----
   your data           <--->   the letter you wrote
   TCP header          <--->   "To: Sarah, Accounts Dept"  (which person)
   IP header           <--->   "To: 12 High St, Berlin"    (which building)
   Ethernet header     <--->   "Next stop: Frankfurt depot" (next hop only)
```

**The key insight:** the Ethernet header is *rewritten at every hop* (each depot
writes a new "next stop"), while the IP header stays the same the whole journey
(the building address never changes). You'll see this again in Chapter 13.

### Worked example: see a real packet's layers

Run this (it captures one packet and shows all its layers):

```bash
sudo tcpdump -i any -c 1 -e -n 'tcp port 443'
```

Then, in another terminal (or just browse the web), trigger some traffic:

```bash
curl -s https://example.com > /dev/null
```

You'll see something like:

```
14:22:01.123456 aa:bb:cc:dd:ee:ff > 11:22:33:44:55:66, ethertype IPv4 (0x0800),
    length 74: 192.168.1.50.51000 > 93.184.216.34.443: Flags [S], seq 1234567,
    win 64240, length 0
    \_______________________/  \_________________________________/
       LINK LAYER (Ethernet)      INTERNET LAYER (IP)  +  TRANSPORT (TCP)
       who's next door             who globally         which program
```

Read it left to right and you are reading the layers from outside in. You don't
need to understand every field yet — just see that they're nested.

### Common confusions

- **"Is this the OSI model?"** The OSI model has 7 layers and is taught in
  certifications. TCP/IP has 4 and is what actually runs. Chapter 3 shows both
  and how they line up.
- **"Do all 4 layers exist in every packet?"** Essentially yes for normal
  traffic. Some packets (like ARP, Chapter 6) skip the IP layer.
- **"Why is there overhead?"** Each header is real bytes on the wire. A typical
  packet spends ~54 bytes on headers to carry up to 1460 bytes of data — about
  4% overhead. That's the price of the design, and it's a bargain.

### Check yourself

1. What does "encapsulation" mean?
2. Which header gets rewritten at every hop, and which stays the same?
3. Why does layering mean HTTP/3 could be invented without upgrading routers?

### Further reading

- **Article:** "What is the OSI Model?" — Cloudflare Learning Center. Clear
  diagrams; read it after Chapter 3.
- **Video (~11 min):** "The OSI Model Explained" — Practical Networking
  (Ed Harmoush). Ed's whole "Networking Fundamentals" series is the best free
  video course on this material. Start it now and follow along with this guide.
- **Article:** "Networking: A Primer" — Julia Evans (jvns.ca). Her zines and
  blog posts on networking are unusually friendly and honest about what's
  confusing.

---

## Chapter 3 — The four layers of TCP/IP

### In one sentence

Link → Internet → Transport → Application: four layers, each with a job, each
with its own protocols and its own name for a chunk of data.

### The four layers

```
+----------------------------------------------------------------------+
| 4. APPLICATION    "What does this data MEAN?"                        |
|                                                                       |
|    HTTP (web), DNS (names), TLS (encryption), SMTP (email),          |
|    SSH (remote login), DHCP (auto-config)                            |
|    Unit of data: a MESSAGE                                            |
+----------------------------------------------------------------------+
| 3. TRANSPORT      "Which PROGRAM, and did it all arrive?"            |
|                                                                       |
|    TCP  - reliable, ordered, connection-based                        |
|    UDP  - fast, unreliable, no connection                            |
|    Unit of data: a SEGMENT (TCP) or DATAGRAM (UDP)                   |
+----------------------------------------------------------------------+
| 2. INTERNET       "Which MACHINE, anywhere in the world?"            |
|                                                                       |
|    IP (IPv4 and IPv6), ICMP (ping/errors)                            |
|    Unit of data: a PACKET                                             |
+----------------------------------------------------------------------+
| 1. LINK           "Which device on THIS cable/radio?"                |
|                                                                       |
|    Ethernet, Wi-Fi (802.11), ARP, and the physical wires/radio       |
|    Unit of data: a FRAME                                              |
+----------------------------------------------------------------------+
```

### What each layer promises the one above it

This framing is worth memorising — **troubleshooting is usually finding out which
layer broke its promise.**

| Layer | Its promise |
|---|---|
| **Link** | "I'll deliver this frame to another device *on this same physical network*, or drop it. No promises about order or success." |
| **Internet** | "I'll make a *best effort* to get this packet to that IP address anywhere on Earth. It might be lost, duplicated, reordered, or delayed. Good luck." |
| **Transport (UDP)** | "Same as IP, but I'll tell you which program it's for and check it wasn't corrupted." |
| **Transport (TCP)** | "You get a reliable, ordered stream of bytes. They arrive exactly once, in order — or the connection breaks and I tell you." |
| **Application** | "I know what these bytes mean: a web page, a DNS answer, an email." |

Notice that **IP promises almost nothing.** All the reliability you experience is
built on top by TCP. That's deliberate — it keeps the middle of the network
simple and fast (see "the end-to-end principle" below).

### The OSI model (for exams)

You'll meet the 7-layer OSI model in certifications. Here's how it maps:

| OSI # | OSI name | TCP/IP layer | What lives there really |
|---|---|---|---|
| 7 | Application | Application | HTTP requests, DNS queries |
| 6 | Presentation | Application | TLS encryption, compression, character encoding |
| 5 | Session | Application | TLS session reuse, RPC sessions |
| 4 | Transport | Transport | **TCP, UDP** |
| 3 | Network | Internet | **IP** |
| 2 | Data Link | Link | Ethernet/Wi-Fi framing, MAC addresses |
| 1 | Physical | Link | Copper, fibre, radio waves, voltages |

People say "Layer 3" for IP and "Layer 7" for HTTP because product names use it:

- **L2 switch** — forwards by MAC address (Chapter 7)
- **L3 switch / router** — forwards by IP address (Chapter 13)
- **L4 load balancer** — picks a server using IP + port, never reads the content
- **L7 load balancer / proxy** — reads the HTTP request and routes by URL
- **"Layer 8 problem"** — an old joke meaning "the user"

> **Simplified:** the mapping above is approximate. TLS doesn't fit neatly into
> OSI layer 6, and nobody in practice cares. Use the 4-layer model to think, and
> the 7-layer names to communicate.

### The two big design ideas

**1. Packet switching.** The phone network used to build a dedicated circuit for
each call, reserved end to end. TCP/IP instead chops data into independently
addressed **packets** that share links. A 1 Gbit/s link can serve thousands of
mostly-idle connections, because they rarely all burst at once. Much more
efficient — at the cost of variability.

**2. The end-to-end principle.** Keep the middle of the network **dumb and fast**;
put the intelligence at the endpoints. A router in the middle just forwards
packets toward their destination — it does not track your connection, retransmit
for you, or care what's inside. This is why you can invent a whole new transport
protocol (like QUIC, Chapter 30) without upgrading a single router on the
internet.

**The hourglass.** Draw the protocol family and it looks like this:

```
  HTTP   DNS   SMTP   SSH   video   games   ...     <- many applications
      \     \    |     |    /      /
                TCP        UDP                      <- a few transports
                  \        /
                     IP                             <- ONE narrow waist
                   /    \
   Ethernet   Wi-Fi   4G/5G   fibre   cable   ...   <- many physical links
```

**Everything runs over IP, and IP runs over everything.** Agreeing on exactly one
internetworking protocol in the middle is what lets any app work over any medium.

### Practice (10 min)

Identify the layer for each of these, then check against Appendix F:

1. Your laptop's Wi-Fi card
2. `192.168.1.50`
3. Your browser asking for `/index.html`
4. Port 443
5. `a4:83:e7:2b:19:0c`
6. Retransmitting a lost chunk of a file
7. `example.com`
8. A router deciding which cable to send a packet down

Then run this and label which layer each piece of output belongs to:

```bash
ip -br addr        # or on macOS: ifconfig
ip route           # or on macOS: netstat -rn
ss -tan | head     # or on macOS: netstat -an | head
curl -sI https://example.com | head -3
```

### Common confusions

- **"Frame vs packet vs segment — do they matter?"** They're the same bytes at
  different layers, and using the right word makes conversations much clearer.
  Frame = link, packet = internet, segment = TCP. "Packet" is also used loosely
  as a catch-all; that's fine in conversation.
- **"Which layer is a firewall?"** Depends on the firewall. Simple ones work at
  L3/L4 (IPs and ports). Modern ones also inspect L7 (HTTP content). Chapter 38.
- **"Why is IP so unreliable? Fix it!"** It's a feature. Making IP reliable would
  make routers stateful, slow, and expensive, and would penalise apps (like video
  calls) that would rather drop a packet than wait for it.

### Check yourself

1. Name the four layers, bottom to top, and one protocol at each.
2. What does IP *not* promise?
3. What is the end-to-end principle, and what does it buy us?
4. Why is the diagram called an hourglass?

### Further reading

- **Video series (free, ~4 h total, highly recommended):** "Networking
  Fundamentals" — Practical Networking on YouTube. Episodes 1–5 cover exactly
  this part. The single best free video resource for this guide.
- **Article:** "What is the OSI model?" and "What is the Internet Protocol?" —
  Cloudflare Learning Center.
- **Book (the standard textbook, very readable):** *Computer Networking: A
  Top-Down Approach* — Kurose & Ross. Chapter 1 covers this material properly. If
  you buy one networking textbook, buy this one.
- **Video (~10 min):** "TCP/IP Model Explained" — PowerCert Animated Videos.
  Simple animations, good for a first pass.

---

## Chapter 4 — Our running example: what happens when you open a website

This is the map for the whole guide. Skim it now; you won't understand every step
yet. Each step names the chapter that explains it. **Come back and re-read this
chapter after Part 7** — it'll feel completely different.

### The setup

- Your laptop on home Wi-Fi, address `192.168.1.50`
- Your router at `192.168.1.1`
- You type `https://example.com` and press Enter

### The journey

**Step 0 — You already have an address.** (Chapter 17)
When you joined the Wi-Fi, your laptop asked "does anyone here give out
addresses?" and your router replied with your IP address, the router's own
address (your **default gateway**), and which **DNS server** to use. That
conversation is called **DHCP**.

**Step 1 — Turn the name into a number.** (Chapter 18)
Your laptop has no idea what `example.com` means. It asks a **DNS** server:
"what's the IP address for example.com?" The answer comes back:
`93.184.216.34`. This might take one round trip, or zero if it's cached.

**Step 2 — Decide where to send it.** (Chapter 13)
`93.184.216.34` is not on your home network, so your laptop consults its
**routing table** and finds the rule "anything I don't recognise, send to
`192.168.1.1`" — the **default route**.

**Step 3 — Find the router's hardware address.** (Chapter 6)
To put a packet on the Wi-Fi, your laptop needs the router's **MAC address**, not
its IP. It shouts to everyone on the local network: "who has `192.168.1.1`?" The
router answers with its MAC address. This is **ARP**.

**Step 4 — Open a connection.** (Chapter 21)
Your laptop and the web server perform a **three-way handshake**: `SYN` →
`SYN-ACK` → `ACK`. Three small packets, one round trip, and now there's a TCP
connection. Your laptop picks a random local **port number** (say 51000) so it
can tell this connection apart from all the others.

**Step 5 — Your router rewrites the address.** (Chapter 12)
Your home router performs **NAT**: it swaps your private `192.168.1.50` for its
own public address, remembers the swap in a table, and forwards the packet on.
Everything coming back gets swapped in reverse.

**Step 6 — Cross the internet.** (Chapter 13)
Ten to twenty **routers** each do the same simple thing: look at the destination
IP, consult a table, decide the next hop, decrement a counter (**TTL**), rewrite
the link-layer header, forward. None of them knows about your connection.

**Step 7 — Set up encryption.** (Chapter 29)
Over the fresh TCP connection, your browser and the server perform a **TLS
handshake**: agree on ciphers, exchange keys, and the server proves its identity
with a **certificate**. That's the padlock in the address bar. Roughly one more
round trip.

**Step 8 — Finally, ask for the page.** (Chapter 28)
Your browser sends an **HTTP** request — a few hundred bytes of text:
`GET / HTTP/1.1`, `Host: example.com`. All encrypted now.

**Step 9 — The response streams back.** (Chapter 22, 24, 25)
The server sends the HTML. TCP chops it into segments, numbers them, and the two
sides use **acknowledgements** to make sure nothing is lost. TCP's **congestion
control** ramps the speed up carefully rather than blasting at full rate.

**Step 10 — Repeat, then hang up.** (Chapter 23)
The browser parses the HTML, finds CSS/JS/images, and fetches them — reusing the
same connection where possible. When idle, the connection closes with `FIN`
packets, and one side waits in a state called **TIME_WAIT** for about a minute.

### The latency budget (why distance costs you)

On a cold connection to a server ~90 ms away (round trip), the *unavoidable*
sequential round trips are:

| Phase | Round trips | Time |
|---|---|---|
| DNS lookup (not cached) | 1 | 90 ms |
| TCP handshake | 1 | 90 ms |
| TLS handshake (TLS 1.3) | 1 | 90 ms |
| HTTP request → first byte | 1 | 90 ms |
| **Total before you see anything** | **~4** | **~360 ms** |

That's before a single byte of the page arrives. **This table explains the entire
existence of CDNs, connection reuse, caching, and HTTP/3** (Parts 6, 7, 9). Your
bandwidth is irrelevant here — only the number of round trips and the distance
matter.

### Practice (10 min) — measure it on your own machine

```bash
curl -w '
  dns lookup:   %{time_namelookup}s
  tcp connect:  %{time_connect}s
  tls setup:    %{time_appconnect}s
  first byte:   %{time_starttransfer}s
  total:        %{time_total}s
' -o /dev/null -s https://example.com
```

Each number is *cumulative from the start*. So:
- DNS took `time_namelookup`
- TCP handshake took `time_connect − time_namelookup`
- TLS took `time_appconnect − time_connect`
- The server's thinking + first byte took `time_starttransfer − time_appconnect`

**Now run it again immediately.** DNS should drop to nearly zero (it's cached).
Run it against a server far away and
compare (try `https://www.gov.au` or `https://www.japan.go.jp`).
**You have just measured the physics of the internet.**

### Common confusions

- **"That's a lot of steps for one page."** It is — and modern browsers do many
  of them in parallel and cache aggressively. But every step is real, and when a
  page is slow, it's slow at one of these steps. Knowing the list is how you find
  out which.
- **"Why so many round trips just to say hello?"** Historical layering, plus
  security. HTTP/3 (Chapter 30) merges the TCP and TLS handshakes to cut this from
  ~3 round trips to 1.

### Check yourself

1. List the ten steps from memory, in order.
2. Which step turns a name into a number?
3. Which two steps involve a handshake?
4. If a page loads slowly, name three different steps that could be responsible.

### Further reading

- **The famous question:** "What happens when you type google.com into a browser
  and press enter?" — search for the GitHub repo `alex/what-happens-when`. It's a
  legendary, extremely detailed answer. Read the summary now, the full version
  after Part 7.
- **Video (~12 min):** "What happens when you type a URL" — several good versions;
  search that phrase. ByteByteGo's is a good visual one.
- **Tool:** browser DevTools → Network tab → click any request → **Timing**.
  It shows exactly these phases for real page loads. Open it now and look.

---

### End of Part 1 — Milestone check

- [ ] I can name the four layers and one protocol at each
- [ ] I can explain encapsulation using the postal analogy
- [ ] I know which header changes at every hop and which doesn't
- [ ] I can list the main steps of loading a web page
- [ ] **I have run the `curl -w` timing command and understood the output**

---

# Part 2 — Your local network (the link layer)

This is the bottom layer: getting a frame from one device to another **on the same
cable or the same Wi-Fi**. It has no idea the internet exists.

## Chapter 5 — MAC addresses and Ethernet frames

### In one sentence

Every network card has a permanent-ish hardware address called a **MAC address**,
and the link layer wraps your data in a **frame** addressed from one MAC to
another.

### The problem

Ten devices are plugged into the same switch, or associated with the same Wi-Fi
access point. When one sends data, how does the right device pick it up and the
others ignore it? They need names — and IP addresses aren't available yet at this
level (that's a layer up).

### The idea, in plain language

Give every network card a **factory-assigned serial number**. Put "from" and "to"
serial numbers on every frame. Each card listens, checks the "to" field, and only
processes frames addressed to it.

### How it actually works: the MAC address

A **MAC address** (Media Access Control) is **48 bits = 6 bytes**, written as six
pairs of hex digits:

```
   a4:83:e7:2b:19:0c
   \______/  \______/
      |          |
      |          +-- serial number, chosen by that manufacturer
      +------------- OUI: identifies the MANUFACTURER
                     (a4:83:e7 = Apple, 00:1a:11 = Google, 52:54:00 = a VM)
```

Three things to know:

1. **`ff:ff:ff:ff:ff:ff` is the broadcast address** — "everyone on this network,
   listen up". You'll see it in Chapter 6.
2. **MAC addresses are flat.** There's no structure like "all European devices
   start with X". You can't route by them across the internet — there'd be a
   billion-row table. This is exactly why IP addresses exist (Chapter 9).
3. **They can be changed.** Phones and laptops now use *randomised* MAC addresses
   on Wi-Fi to stop shops tracking you between visits. If your MAC's second hex
   digit is `2`, `6`, `a`, or `e`, it's a locally-assigned (often random) one.

### How it actually works: the Ethernet frame

```
+-----------+-----------+-----------+---------------------+---------+
| Dest MAC  | Src MAC   | EtherType | Payload             | FCS     |
| 6 bytes   | 6 bytes   | 2 bytes   | 46 - 1500 bytes     | 4 bytes |
+-----------+-----------+-----------+---------------------+---------+
   who it's    who sent   what's       your actual data     error
   for         it         inside       (the IP packet)      check
```

- **EtherType** tells the receiver what's inside so it can hand it to the right
  code: `0x0800` = IPv4, `0x86DD` = IPv6, `0x0806` = ARP.
- **Payload max 1500 bytes** — this number is called the **MTU** and it matters a
  lot (Chapter 15).
- **Payload min 46 bytes** — short frames get padded. A tiny 1-byte message still
  costs a 64-byte frame.
- **FCS (Frame Check Sequence)** is a checksum. If it doesn't match, the card
  **silently throws the frame away**. Ethernet does not retransmit — that's TCP's
  job, three layers up.

### Worked example: what your machine actually has

```bash
# Linux
ip -br link

# macOS
ifconfig | grep -E '^\w|ether'
```

You'll see something like:

```
lo               UNKNOWN        00:00:00:00:00:00      <- loopback (fake, internal)
eth0             UP             a4:83:e7:2b:19:0c      <- wired
wlan0            UP             9a:1f:04:c8:22:7e      <- Wi-Fi (note the 'a' -> randomised)
```

Look up the first three bytes of your MAC at `macvendors.com` — it should name
your laptop's manufacturer (unless it's randomised).

### Practice (20 min)

**A. See a real frame with its layers labelled:**

```bash
sudo tcpdump -i any -c 3 -e -n
```

The `-e` flag shows the Ethernet header. You'll see
`aa:bb:cc:dd:ee:ff > 11:22:33:44:55:66, ethertype IPv4 (0x0800)` at the front of
each line. That's a real frame header.

**B. Watch the MAC change but the IP stay the same.** This is the key insight of
Chapter 2, made visible:

```bash
# capture traffic to a website
sudo tcpdump -i any -c 5 -e -n 'host example.com'
```

Every packet you send to `example.com` has:
- **Destination MAC** = your *router's* MAC (the next hop), always
- **Destination IP** = `93.184.216.34` (the final destination), always

The MAC says "next door"; the IP says "eventually". Write that in `notes.md`.

**C. Check for hardware errors** (useful troubleshooting skill):

```bash
# Linux
ip -s link show eth0        # look for RX/TX 'errors' and 'dropped'
ethtool -S eth0 2>/dev/null | grep -Ei 'err|drop|crc'

# macOS
netstat -i                   # Ierrs / Oerrs columns
```

Non-zero and *climbing* CRC errors means a physical problem: a bad cable, a
loose connector, interference. This is often the real cause of "the network is
slow" and almost nobody checks it.

### Common confusions

- **"MAC address vs IP address?"** MAC = which device on *this* cable, assigned
  by the manufacturer, never routed. IP = which device *in the world*, assigned by
  the network you're on, routed globally. You need both, always.
- **"Is a MAC address permanent?"** Traditionally yes (burned into the card).
  In practice modern devices randomise them for privacy, and you can change them
  in software.
- **"What if two devices have the same MAC?"** On the same network segment,
  chaos — frames go to the wrong place intermittently. It happens with cloned VMs.

### Check yourself

1. How many bytes is a MAC address, and what do the first three bytes tell you?
2. What does the EtherType field do?
3. What happens to a frame whose FCS doesn't match?
4. When you send a packet to a website, whose MAC address is in the destination
   field?

### Further reading

- **Video (~10 min):** "Ethernet Frames Explained" — Practical Networking.
- **Article:** "What is a MAC address?" — Cloudflare Learning Center.
- **Book:** *Computer Networking: A Top-Down Approach* (Kurose & Ross), Chapter 6
  (the Link Layer).
- **Tool:** `macvendors.com` — paste any MAC address, see the manufacturer.

---

## Chapter 6 — ARP: finding your neighbour

### In one sentence

Before your machine can send a frame, it must translate the next hop's **IP
address** into a **MAC address** — and it does that by shouting a question to
everyone on the local network.

### The problem

You want to reach `192.168.1.1` (your router). You know its IP. But to build an
Ethernet frame you need its **MAC**. You've never spoken to it before. How do you
find out?

### The idea, in plain language

Shout.

```
YOU (to everyone on the network):
    "Who has 192.168.1.1? Tell 192.168.1.50!"

EVERYONE ELSE:  ...not me... not me... not me...

THE ROUTER (directly back to you):
    "192.168.1.1 is at bb:bb:bb:bb:bb:bb"

YOU: *writes it down in a notebook so you don't have to ask again*
```

That's **ARP** — Address Resolution Protocol. The notebook is the **ARP cache**.

### How it actually works

```
1. Your machine checks its ARP cache. Got 192.168.1.1? Use it, done.

2. Cache miss -> build an ARP REQUEST and send it to the BROADCAST MAC
   (ff:ff:ff:ff:ff:ff), so every device on the segment receives it.

      "Who has 192.168.1.1?  Tell 192.168.1.50 (aa:aa:aa:aa:aa:aa)"

3. Every device looks at it. Devices whose IP doesn't match ignore it.

4. The matching device sends an ARP REPLY -- UNICAST, straight back to you:

      "192.168.1.1 is at bb:bb:bb:bb:bb:bb"

5. You cache it (typically for tens of seconds to a few minutes) and now you can
   build frames.
```

Two important details:

- **ARP frames are not IP packets.** They have EtherType `0x0806`, no IP header,
  no TTL. **They cannot cross a router.** ARP only works within one local network.
- **ARP has no security whatsoever.** Any device can reply "that IP is at *my*
  MAC" and be believed. That's **ARP spoofing** — the classic local-network
  attack (Chapter 38).

### The "am I local?" decision (important!)

Before ARPing, your machine asks: *is this destination on my own network, or do I
need the router?*

```
Your address:      192.168.1.50 / 24     ("/24" explained in Chapter 10)
Your network is:   192.168.1.*

Sending to 192.168.1.99?
   -> same network -> ARP for 192.168.1.99 itself, send directly

Sending to 93.184.216.34?
   -> different network -> ARP for the GATEWAY (192.168.1.1), send there
      (the packet's destination IP is still 93.184.216.34 -- only the frame
       is addressed to the router)
```

**This decision is why subnet masks matter**, and it's the bridge into Chapter 10.

### Worked example: watch it happen live

```bash
# 1. Look at your current ARP cache
ip neigh                    # Linux
arp -a                      # macOS / Linux / Windows

# 2. Clear one entry so you can watch it be rebuilt (Linux)
sudo ip neigh flush all

# 3. In one terminal, start watching for ARP:
sudo tcpdump -i any -n arp

# 4. In another terminal, ping your router (use YOUR gateway address)
ping -c 2 192.168.1.1
```

In the tcpdump window you'll see exactly this:

```
ARP, Request who-has 192.168.1.1 tell 192.168.1.50, length 28
ARP, Reply 192.168.1.1 is-at bb:bb:bb:bb:bb:bb, length 46
```

**That is one of the most satisfying things in this guide.** You just watched a
protocol from 1982 (RFC 826) work in front of you.

### The ARP cache states (Linux)

```bash
ip neigh
192.168.1.1 dev wlan0 lladdr bb:bb:bb:bb:bb:bb REACHABLE
192.168.1.99 dev wlan0 lladdr cc:cc:cc:cc:cc:cc STALE
```

| State | Meaning |
|---|---|
| `REACHABLE` | Confirmed recently — use it |
| `STALE` | Probably still fine; will be re-verified on next use |
| `DELAY` / `PROBE` | Currently re-checking |
| `FAILED` | No answer — the device is gone or not responding |
| `PERMANENT` | Manually pinned, never expires |

### Practice (25 min)

1. Run the live capture above. Save the output in `notes.md`.
2. `ip neigh` (or `arp -a`) — how many neighbours does your machine know? Can you
   identify them (router, printer, phone, TV)?
3. Ping a device on your network that's turned **off**. Watch tcpdump: you'll see
   ARP requests repeated with no reply, and `ping` reports "Destination Host
   Unreachable". **That error means ARP failed** — the machine isn't there. This
   is a genuinely useful diagnostic to recognise.
4. Ping a device on the *internet* and watch ARP. You'll see... nothing new,
   because your router's MAC is already cached. Flush the cache and try again.

### Common confusions

- **"Why not just use MAC addresses for everything?"** Because they're flat.
  A router would need a table with every MAC address on Earth. IP addresses are
  hierarchical (Chapter 9), so a router can say "everything starting with 93.184
  goes that way" in one line.
- **"Why not just use IP for everything?"** Because the physical network only
  understands MACs. IP is a fiction layered on top; ARP is the glue.
- **"Does IPv6 use ARP?"** No — it uses **NDP** (Neighbor Discovery Protocol),
  which does the same job using ICMPv6 messages. Chapter 16.
- **"Destination Host Unreachable vs Request Timed Out?"** *Host Unreachable*
  usually means ARP failed on the last hop (device absent). *Timed out* means the
  packet went somewhere and nothing came back. Different problems.

### Check yourself

1. What does ARP translate, and in which direction?
2. Why is an ARP request sent to the broadcast address but the reply sent
   directly?
3. Can an ARP request cross a router? Why or why not?
4. When you browse the web, whose MAC does your laptop ARP for?

### Further reading

- **Video (~8 min):** "ARP Explained" — Practical Networking. Excellent
  animation of exactly this.
- **Article:** "How ARP Works" — search PracticalNetworking.net's ARP article.
- **RFC 826** — the original ARP spec from 1982. It is *three pages long* and
  surprisingly readable. A good first RFC to try.
- **Blog:** Julia Evans, "How does ARP work?" and her networking zines (jvns.ca).

---

## Chapter 7 — Switches, and what "broadcast domain" means

### In one sentence

A switch learns which device is on which port and sends each frame only where it
needs to go — but broadcasts still reach everyone, and that group is called a
**broadcast domain**.

### The problem

Old networks (with **hubs**) copied every frame to every port. Every device saw
everyone's traffic: wasteful, slow, and completely insecure. With 20 devices, 19
copies of every frame were pointless.

### The idea, in plain language

The switch **learns by listening**. Every time a frame arrives, it notes "device
with MAC `aa:aa` is on port 3". Next time a frame is addressed to `aa:aa`, it
sends it out port 3 only.

### How it actually works

```
   Switch's MAC address table (it builds this itself, by watching):

     MAC address          Port
     -----------          ----
     aa:aa:aa:aa:aa:aa      3       <- learned when device on port 3 sent a frame
     bb:bb:bb:bb:bb:bb      7
     cc:cc:cc:cc:cc:cc      1

   For each arriving frame:
     1. LEARN:  note the SOURCE MAC and the port it came in on
     2. LOOK UP the DESTINATION MAC:
          - found?      -> send out ONLY that port          (unicast forwarding)
          - not found?  -> send out ALL ports except the incoming one (flooding)
          - broadcast?  -> send out ALL ports except the incoming one
```

Entries expire (typically after 5 minutes) so devices can move.

### Collision domain vs broadcast domain

Two terms that sound similar and confuse everyone.

| Term | Meaning | Today |
|---|---|---|
| **Collision domain** | Set of devices that could interfere with each other if they transmit simultaneously | Each switch port is its own collision domain. **Collisions are effectively extinct** on modern wired networks. |
| **Broadcast domain** | Set of devices that receive each other's broadcasts (like ARP) | **This still matters a lot.** |

**The rule to remember:**

```
A SWITCH does NOT break up a broadcast domain.
A ROUTER (or a VLAN boundary) DOES.
```

So if you plug 200 devices into a chain of switches, all 200 are in one broadcast
domain — every ARP request from any of them reaches all of them. That's fine at
200; it's a problem at 2,000. Rule of thumb: keep a broadcast domain under a few
hundred devices.

### VLANs, briefly

A **VLAN** (Virtual LAN) lets one physical switch behave as several separate
switches. Ports are tagged with a VLAN number; devices in VLAN 10 cannot see
broadcasts from VLAN 20 — they're separate broadcast domains — and traffic
between them must pass through a router.

This is how one office switch safely carries guest Wi-Fi, staff machines, VoIP
phones, and CCTV without them seeing each other.

### The loop danger (why STP exists)

If you accidentally connect two switches with **two** cables, a broadcast frame
can circle forever — because, unlike IP packets, **Ethernet frames have no TTL
(hop counter)**. Within seconds the network saturates. This is a **broadcast
storm**, and it takes down whole offices.

**STP** (Spanning Tree Protocol) prevents it: switches talk to each other, detect
loops, and deliberately *block* redundant ports until they're needed. If you ever
add a redundant link between switches, make sure STP is on.

### Practice (20 min)

**A. See a software switch.** Docker creates one on your machine:

```bash
# Linux, if you have Docker
bridge link                 # show bridge ports
bridge fdb show | head -20  # the MAC address table -- exactly as described above
```

**B. Prove that a switch filters traffic.** Ping between two devices on your
network while capturing on a *third* device. With a switch, you should see
almost none of their traffic (only broadcasts). With an old hub, you'd see
everything.

**C. Watch broadcasts:**

```bash
sudo tcpdump -i any -n 'broadcast or multicast'
```

Leave it running for a minute. You'll see ARP requests, DHCP, mDNS
(`.local` name discovery), SSDP (device discovery), and more. **Every device on
your network is doing this constantly.** That's the broadcast domain at work.

### Common confusions

- **"Switch vs router?"** Switch = forwards by **MAC**, within one network,
  layer 2. Router = forwards by **IP**, between networks, layer 3. Your home
  "router" box is actually a router + switch + Wi-Fi access point + firewall +
  DHCP server in one case.
- **"Is my Wi-Fi access point a switch?"** Functionally similar — it bridges
  wireless devices into the wired network and forwards by MAC.
- **"Managed vs unmanaged switch?"** Unmanaged = plug in and it works.
  Managed = you can configure VLANs, STP, port security, monitoring.

### Check yourself

1. How does a switch learn which device is on which port?
2. What does a switch do with a frame whose destination MAC it doesn't know?
3. Does a switch break up a broadcast domain? What does?
4. Why is a Layer 2 loop more dangerous than a Layer 3 (IP) loop?

### Further reading

- **Video (~12 min):** "How Switches Work" — Practical Networking.
- **Video:** "Spanning Tree Protocol Explained" — Practical Networking or
  NetworkChuck.
- **Article:** "What is a network switch?" — Cloudflare Learning Center.
- **Book:** *Network Warrior* — Gary A. Donahue. Practical, real-world switching
  and routing from someone who has clearly suffered. Chapters on switches and
  VLANs are excellent.

---

## Chapter 8 — Wi-Fi: same idea, different physics

### In one sentence

Wi-Fi does the same job as Ethernet — deliver frames between devices on a local
network — but over a shared, noisy, half-duplex radio channel, which changes its
behaviour in ways you will feel.

### The problem

On a cable, you have a dedicated wire and can send and receive simultaneously. On
radio, **everyone shares the same air**, you can't listen while you're
transmitting, and the signal degrades with distance and interference.

### The key differences (and why you care)

| | Wired Ethernet | Wi-Fi |
|---|---|---|
| **Medium** | Your own dedicated cable | **Shared** air, half-duplex |
| **Taking turns** | Not needed | **CSMA/CA**: listen, wait a random time, then transmit |
| **Acknowledgements** | None at this layer | **Every frame is ACKed and retried** at the link layer |
| **Speed** | Fixed once negotiated | **Changes constantly** with signal quality |
| **Loss** | Nearly zero on a good cable | Bursty; interference, distance, walls, microwaves |
| **Latency** | Microseconds, very stable | 1–50 ms, **highly variable** |

Three consequences you will actually experience:

**1. Wi-Fi hides loss from TCP — mostly.** Because the link layer retries failed
frames itself, TCP often doesn't see a loss. Good (TCP doesn't panic), but those
retries **add latency and jitter**, which is why Wi-Fi feels less "crisp" for
video calls even at high speeds.

**2. Airtime is the real currency, not bandwidth.** One distant device
transmitting slowly hogs the air, slowing everyone. This is called the
"performance anomaly", and it's why one person in the far bedroom can make the
whole house's Wi-Fi feel bad.

**3. Advertised speeds are fantasy.** "Wi-Fi 6, 1200 Mbps" is a theoretical
combined figure. Because the channel is shared and half-duplex, real TCP
throughput is typically **40–60%** of the negotiated rate, less with multiple
devices.

### The other links you'll meet

| Link type | Notes |
|---|---|
| **Fibre / cable to the home** | High speed; cable is shared upstream in the neighbourhood |
| **Mobile (4G/5G)** | 20–80 ms latency that *spikes* when the radio wakes from idle; deep buffers; your IP changes as you move |
| **VPN tunnels** | Your data wrapped inside another packet — which **steals bytes from the MTU** (Chapter 15). A very common cause of weird bugs. |
| **Loopback (`lo`)** | A fake interface inside your own machine, address `127.0.0.1`. Traffic never touches a wire. |

### Practice (20 min)

**A. Look at your actual Wi-Fi quality:**

```bash
# Linux
iw dev wlan0 link                # signal strength and current bitrate
iw dev wlan0 station dump | grep -E 'signal|bitrate|retries|failed'

# macOS
system_profiler SPAirPortDataType | grep -A10 "Current Network"
# or: wdutil info    (may need sudo)
```

Two numbers matter most:
- **Signal**: better than −60 dBm is good; worse than −75 dBm is trouble.
  (It's negative — closer to zero is stronger.)
- **Tx bitrate**: if it's far below your Wi-Fi standard's maximum, you're far
  from the access point or there's interference.

**B. Watch it degrade.** Run this while walking away from your router:

```bash
# Linux
watch -n1 'iw dev wlan0 link | grep -E "signal|bitrate"'
```

Watch the bitrate drop in steps as you move. **That's the radio adapting**, and
it explains why your video call gets worse in the far room.

**C. Compare latency stability.** Ping your router over Wi-Fi and (if you can)
over a cable:

```bash
ping -c 50 192.168.1.1
```

Look at the summary line: `min/avg/max/mdev`. **`mdev` is the jitter.** On a cable
it'll be a fraction of a millisecond. On Wi-Fi it may be several milliseconds and
occasionally much worse. That variability is why Wi-Fi feels different even when
"fast".

### Common confusions

- **"More bars = faster?"** Bars show signal, not speed. A strong signal with
  heavy interference or 30 other devices can still be slow.
- **"5 GHz is always better than 2.4 GHz?"** 5 GHz is faster and less congested
  but has **shorter range** and is blocked more by walls. 2.4 GHz reaches further
  but is crowded (and shared with microwaves, Bluetooth, baby monitors).
- **"Wi-Fi is a network."** Wi-Fi is a *link layer* — a way to carry frames. IP,
  TCP, and everything above are identical over Wi-Fi and Ethernet.

### Check yourself

1. Why does Wi-Fi acknowledge every frame when Ethernet doesn't?
2. Why can one slow device degrade everyone's Wi-Fi?
3. What does `mdev` in ping output tell you, and why is it bigger on Wi-Fi?

### Further reading

- **Video:** "How Wi-Fi Works" — Practical Networking, or Branch Education's
  visually spectacular version.
- **Article:** "802.11 Wireless Fundamentals" — search Cisco or Cloudflare's
  explainer.
- **Book:** *802.11 Wireless Networks: The Definitive Guide* — Matthew Gast.
  The reference, if you go deep on wireless.
- **Tool:** WiFi Explorer (macOS) or WiFiAnalyzer (Android) — see all nearby
  networks and channel congestion. Genuinely eye-opening in an apartment block.

---

### End of Part 2 — Milestone check

- [ ] I can explain what a MAC address is and what its first 3 bytes mean
- [ ] **I have watched an ARP request and reply in tcpdump**
- [ ] I can explain how a switch learns, and what it does with an unknown MAC
- [ ] I know the difference between a collision domain and a broadcast domain
- [ ] I know why a Layer 2 loop is dangerous
- [ ] I have checked my own Wi-Fi signal and bitrate

---

# Part 3 — IP addresses and getting around the world

This part contains **Chapter 10**, the one topic most beginners bounce off.
It is not hard. It has just always been taught badly. Go slowly and do the
practice.

## Chapter 9 — What an IP address actually is

### In one sentence

An IP address is a number identifying a machine's *location on the network*, and
it's split into two parts: which **network** it's on, and which **host** it is on
that network.

### The problem

MAC addresses can't scale (Chapter 6): they're flat, so a global router would need
a table listing every device on Earth. We need addresses with **structure** — so a
router can say "everything starting with 93.184 goes left" in a single rule.

### The idea, in plain language

Use a **postcode**, not a name.

```
   MAC address  is like  "Sarah Chen"
        -> to deliver, you'd need a list of everyone on Earth

   IP address   is like  "SW1A 1AA" / "10001"
        -> the postcode itself tells you the routing
        -> the sorting office only needs to know "SW.. goes to London"
```

An IP address encodes **where you are**, not **who you are**. Move your laptop to
a café and its MAC stays the same while its IP changes — because its *location*
changed.

### How it actually works: the address

An **IPv4** address is **32 bits**, written as four numbers 0–255 separated by
dots:

```
   192  .  168  .   1   .  50
   \_/     \_/     \_/     \_/
    |       |       |       |
   8 bits  8 bits  8 bits  8 bits     =  32 bits total

   In binary:  11000000 . 10101000 . 00000001 . 00110010
```

Each of those four numbers is called an **octet** (8 bits). This is why they max
out at 255 (Chapter 0.5).

**32 bits gives about 4.3 billion addresses.** There are far more than 4.3 billion
connected devices, which is the whole reason NAT (Chapter 12) and IPv6
(Chapter 16) exist.

### The two halves

Here's the crucial idea, and everything in Chapter 10 follows from it:

```
   192.168.1.50
   \_________/ \/
        |       |
   NETWORK part  HOST part
   "which street"  "which house"
```

- **Every device on the same network shares the same network part.**
- **Each device has a different host part.**
- A router only needs to know about *networks*, not individual hosts. That's the
  scaling win.

But where exactly is the dividing line? **That's what the subnet mask tells you**
— and it's Chapter 10.

### Two addresses you never assign to a device

In every network, two addresses are reserved:

```
   Network 192.168.1.x  (with a /24 mask, explained next chapter)

   192.168.1.0     <- the NETWORK ADDRESS. Names the network itself.
                      All host bits = 0. Never assign it to a device.

   192.168.1.1
   192.168.1.2      <- usable addresses for actual devices
   ...
   192.168.1.254

   192.168.1.255   <- the BROADCAST ADDRESS. "everyone on this network".
                      All host bits = 1. Never assign it to a device.
```

So a "256-address" network gives you **254 usable** addresses. That "minus 2" is
the source of endless confusion; now you know exactly why.

### Special addresses to recognise

| Address / range | Meaning |
|---|---|
| `127.0.0.1` | **Loopback** — "this machine". Never leaves your computer. `localhost`. |
| `0.0.0.0` | "Any address" / "unspecified". A server binding to `0.0.0.0` means "listen on all my addresses". |
| `255.255.255.255` | Broadcast to this local network only |
| `169.254.x.x` | **Self-assigned** — "I asked for an address and nobody answered". A DHCP failure signal (Chapter 17). |
| `10.x.x.x`, `172.16-31.x.x`, `192.168.x.x` | **Private** — usable inside your network, never routed on the internet (Chapter 12) |
| `224.x.x.x` – `239.x.x.x` | **Multicast** — "everyone subscribed to this group" |

**Seeing `169.254.something` on a device is an instant diagnosis:** it never got a
DHCP address. That single fact will save you time for the rest of your career.

### Practice (15 min)

**A. Find and understand your own address:**

```bash
# Linux
ip -4 addr show
# macOS
ipconfig getifaddr en0 ; ifconfig en0 | grep 'inet '
```

Write down your address and note which of the private ranges it falls in.

**B. See the loopback do nothing:**

```bash
ping -c 2 127.0.0.1
sudo tcpdump -i any -c 5 -n 'host 127.0.0.1'   # in another terminal, then re-ping
```

Traffic to `127.0.0.1` appears on the `lo` interface and **never touches your
network card**.

**C. Convert to binary yourself** (do at least two by hand before checking):

```bash
python3 -c "print('.'.join(format(int(o),'08b') for o in '192.168.1.50'.split('.')))"
```

Use the place-value table from Chapter 0.5: `192 = 128 + 64` so `11000000`.

### Common confusions

- **"Is my IP address permanent?"** Usually not. Home devices get a new one
  periodically from DHCP (Chapter 17). Your *public* IP may also change.
- **"Why does my machine have several IP addresses?"** One per interface
  (Ethernet, Wi-Fi, VPN, Docker) plus loopback. Normal and expected.
- **"What's the difference between my IP and my public IP?"** Your device has a
  *private* IP (`192.168.x.x`); your router has a *public* one shared by
  everyone in your house. Chapter 12. Try: `curl ifconfig.me`.
- **"IP address vs domain name?"** The name (`example.com`) is a label for
  humans; DNS translates it to an address (Chapter 18).

### Check yourself

1. How many bits is an IPv4 address, and why do the octets stop at 255?
2. What are the two parts of an IP address?
3. Why can't you assign the network address or the broadcast address to a device?
4. What does it mean if a device has the address `169.254.10.7`?

### Further reading

- **Video (~10 min):** "IP Addresses Explained" — PowerCert Animated Videos, or
  Practical Networking's "IPv4 Addressing".
- **Article:** "What is an IP address?" — Cloudflare Learning Center.
- **Tool:** `curl ifconfig.me` or visit `whatismyipaddress.com` to see your
  public address.

---

## Chapter 10 — Subnet masks: the part everyone finds hard

> **This is the most important chapter in Part 3.** Read it twice if you need to.
> Once it clicks, it stays clicked forever.

### In one sentence

The subnet mask marks where the **network part** of an address ends and the
**host part** begins — and that single line determines who your machine can talk
to directly.

### Why it exists

From Chapter 9: an IP address has a network part and a host part. But
`192.168.1.50` doesn't say where the split is. Is the network `192`?
`192.168`? `192.168.1`?

The mask answers that. It has to be stated separately, which is why you always
see them together:

```
   192.168.1.50 / 24              or       192.168.1.50
                                            255.255.255.0
```

Those two lines mean **exactly the same thing**. Learn to read both.

### The idea, in plain language

Line up the mask under the address. **Where the mask is 1, that's network. Where
it's 0, that's host.**

```
   address:   192  .  168  .   1   .   50
   binary:  11000000.10101000.00000001.00110010

   mask /24: 11111111.11111111.11111111.00000000
             \_______________________/ \________/
                 NETWORK (24 ones)      HOST (8 zeros)

   So:  network part = 192.168.1
        host part    = .50
```

**`/24` literally means "the first 24 bits are the network".** That's all the
slash number is: a count of 1-bits in the mask. This is called **CIDR notation**
(Classless Inter-Domain Routing).

### The four numbers you can now work out

For any address + mask, you can determine:

```
   NETWORK ADDRESS   -> set every host bit to 0     (names the network)
   BROADCAST ADDRESS -> set every host bit to 1     ("everyone here")
   FIRST USABLE      -> network address + 1
   LAST USABLE       -> broadcast address - 1
   HOW MANY HOSTS    -> 2^(host bits) - 2
```

For `192.168.1.50/24`:

```
   host bits = 32 - 24 = 8

   network:    192.168.1.0        (host bits all 0)
   broadcast:  192.168.1.255      (host bits all 1)
   usable:     192.168.1.1  to  192.168.1.254
   count:      2^8 - 2 = 254 hosts
```

### The shortcut: block size (use this, not binary)

Converting to binary every time is slow. **Here is the method professionals
actually use.**

**Step 1 — Memorise this table.** It's the only thing you need to memorise in
this entire guide.

| CIDR | Mask (last non-255 octet) | **Block size** | Usable hosts |
|---|---|---|---|
| /24 | 255.255.255.**0** | **256** | 254 |
| /25 | 255.255.255.**128** | **128** | 126 |
| /26 | 255.255.255.**192** | **64** | 62 |
| /27 | 255.255.255.**224** | **32** | 30 |
| /28 | 255.255.255.**240** | **16** | 14 |
| /29 | 255.255.255.**248** | **8** | 6 |
| /30 | 255.255.255.**252** | **4** | 2 |

**Block size = 256 − (the mask number).** That's the whole trick.
`/26` → mask 192 → block size `256 − 192 = 64`.

**Step 2 — The four-step method.**

```
   1. Find the block size:        256 - mask value
   2. Count up in blocks from 0 until you pass your address;
      the network address is the block you landed IN.
   3. Broadcast = next network - 1
   4. Usable = everything between them
```

### Worked example 1: `192.168.1.130/26`

```
   /26 -> mask 255.255.255.192 -> block size = 256 - 192 = 64

   Count in 64s:  0,  64,  128,  192,  (256)
                            ^^^
   130 falls between 128 and 192, so it's in the 128 block.

   NETWORK:    192.168.1.128
   BROADCAST:  192.168.1.191      (next block 192, minus 1)
   USABLE:     192.168.1.129  to  192.168.1.190
   HOSTS:      62
```

Check: is `192.168.1.200` on the same network as `192.168.1.130`? Count: 200 is in
the 192 block, 130 is in the 128 block. **Different networks** — those two devices
cannot talk directly; they need a router.

### Worked example 2: `10.0.0.45/27`

```
   /27 -> mask 255.255.255.224 -> block size = 256 - 224 = 32

   Count in 32s:  0, 32, 64, 96, 128, ...
                     ^^
   45 falls between 32 and 64.

   NETWORK:    10.0.0.32
   BROADCAST:  10.0.0.63
   USABLE:     10.0.0.33  to  10.0.0.62
   HOSTS:      30
```

### Worked example 3: when the mask isn't in the last octet

`172.16.20.100/22`

```
   /22 = 255.255.252.0
   The "interesting" octet is the THIRD one (mask value 252).
   Block size = 256 - 252 = 4  -- and we count in the THIRD octet.

   Count in 4s in the third octet:  0, 4, 8, 12, 16, 20, 24, ...
                                                     ^^
   20 falls exactly on a boundary, so 20 IS the network.

   NETWORK:    172.16.20.0
   BROADCAST:  172.16.23.255     (next network is 172.16.24.0, minus 1)
   USABLE:     172.16.20.1  to  172.16.23.254
   HOSTS:      2^10 - 2 = 1022
```

**How to find the interesting octet:** it's the one where the mask is neither 255
nor 0. `/22` = `255.255.252.0` → third octet.

### Why this matters in practice: "can I reach it directly?"

Every time your machine sends a packet, it does this calculation (Chapter 6):

```
   Is the destination's network part the same as MY network part?

     YES -> it's a neighbour. ARP for it and send directly.
     NO  -> it's remote. Send to the default gateway (the router).
```

**Get the mask wrong and this breaks in confusing ways.** Classic symptom: you can
ping some machines on your network but not others, and the internet works fine.
Almost always a mask mismatch.

### Practice (40 min) — do these on paper first

**Set A — read the mask.** For each, give network, broadcast, first/last usable,
and host count. *Answers in Appendix C.*

```
   1.  192.168.10.75  /24
   2.  192.168.10.75  /26
   3.  10.5.3.200     /25
   4.  172.20.15.100  /28
   5.  10.0.0.6       /30
   6.  192.168.4.170  /27
   7.  172.16.35.42   /23
   8.  10.10.10.10    /22
```

**Set B — are these two on the same network?** *Answers in Appendix C.*

```
   9.  192.168.1.10/24   and  192.168.1.200/24
   10. 192.168.1.10/26   and  192.168.1.200/26
   11. 10.0.0.100/25     and  10.0.0.130/25
   12. 172.16.4.5/22     and  172.16.7.250/22
```

**Set C — check your answers with a tool** (only *after* doing them by hand):

```bash
# Install if needed: apt install ipcalc  /  brew install ipcalc
ipcalc 192.168.10.75/26

# Or use Python, which is everywhere:
python3 - <<'EOF'
import ipaddress
for cidr in ["192.168.10.75/24", "192.168.10.75/26", "10.5.3.200/25",
             "172.20.15.100/28", "10.0.0.6/30", "192.168.4.170/27",
             "172.16.35.42/23", "10.10.10.10/22"]:
    n = ipaddress.ip_network(cidr, strict=False)
    hosts = list(n.hosts())
    print(f"{cidr:20} network={n.network_address}  broadcast={n.broadcast_address}"
          f"  usable={hosts[0]}-{hosts[-1]}  count={len(hosts)}")
EOF
```

**Set D — check your own network:**

```bash
ip -4 addr show          # find your address and prefix, e.g. 192.168.1.50/24
ipcalc <your address>/<your prefix>
```

How many devices can your home network hold? Which addresses are usable?

> **Do not move on until you can do Set A without looking.** Everything in
> Chapters 12–13 assumes it, and every networking job interview asks it.

### Build it in Go (40 min) — a subnet calculator you can trust

Do the paper exercises above first. This program will check your answers,
and writing it teaches you what the "four numbers" really are: **bit
operations on a 32-bit (or 128-bit) number**.

Go's `net/netip` package represents addresses as small, comparable values
(`netip.Addr`, `netip.Prefix`). It handles both IPv4 and IPv6 and gives you
`Masked()`, `Contains()`, and `Overlaps()` for free. The program adds the parts
that `netip` leaves to you: the broadcast address, the usable host range, and
splitting a block into smaller ones.

Save as `subnet/main.go`:

```go
// subnet: a subnet calculator and tiny IPAM tool built on net/netip.
//
//	go run ./subnet info  192.168.10.77/26
//	go run ./subnet split 10.0.0.0/16 24          # carve /24s out of a /16
//	go run ./subnet overlap 10.0.0.0/16 10.0.128.0/20 172.16.0.0/12
//	go run ./subnet alloc 10.0.0.0/16 22 10.0.0.0/24 10.0.4.0/22
package main

import (
	"fmt"
	"math/big"
	"net/netip"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: subnet info|split|overlap|alloc <args>")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "info":
		err = info(os.Args[2])
	case "split":
		err = split(os.Args[2], os.Args[3])
	case "overlap":
		err = overlap(os.Args[2:])
	case "alloc":
		err = alloc(os.Args[2], os.Args[3], os.Args[4:])
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// info prints what Chapter 10 teaches you to work out by hand.
func info(s string) error {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		return err
	}
	host := p.Addr() // the address as typed, e.g. 192.168.10.77
	p = p.Masked()   // zero the host bits -> the network address
	first, last := p.Addr(), lastAddr(p)
	fmt.Printf("input        %s\n", host)
	fmt.Printf("network      %s\n", p)
	if host.Is4() {
		fmt.Printf("netmask      %s\n", mask4(p.Bits()))
		fmt.Printf("broadcast    %s\n", last)
	}
	total := new(big.Int).Lsh(big.NewInt(1), uint(host.BitLen()-p.Bits()))
	usable := new(big.Int).Set(total)
	if host.Is4() && p.Bits() <= 30 { // network + broadcast are not usable hosts
		usable.Sub(usable, big.NewInt(2))
		first, last = first.Next(), last.Prev()
	}
	fmt.Printf("host range   %s - %s\n", first, last)
	fmt.Printf("addresses    %s total, %s usable\n", total, usable)
	fmt.Printf("private?     %v\n", host.IsPrivate())
	return nil
}

// split carves a prefix into equal smaller prefixes (Chapter 11's design exercise).
func split(s, bitsStr string) error {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		return err
	}
	bits, err := strconv.Atoi(bitsStr)
	if err != nil || bits < p.Bits() || bits > p.Addr().BitLen() {
		return fmt.Errorf("new prefix length must be between /%d and /%d", p.Bits(), p.Addr().BitLen())
	}
	n := 0
	for sub := netip.PrefixFrom(p.Masked().Addr(), bits); p.Contains(sub.Addr()); {
		if n < 8 {
			fmt.Println(sub)
		}
		n++
		next := lastAddr(sub).Next()
		if !next.IsValid() { // ran off the end of the address space
			break
		}
		sub = netip.PrefixFrom(next, bits)
	}
	if n > 8 {
		fmt.Printf("... %d subnets in total\n", n)
	}
	return nil
}

// overlap reports every pair of prefixes that overlap -- the check that would
// have saved the "two companies both picked 10.0.0.0/16" migration.
func overlap(args []string) error {
	ps, err := parseAll(args)
	if err != nil {
		return err
	}
	clean := true
	for i := range ps {
		for j := i + 1; j < len(ps); j++ {
			if ps[i].Overlaps(ps[j]) {
				fmt.Printf("OVERLAP  %s  <->  %s\n", ps[i], ps[j])
				clean = false
			}
		}
	}
	if clean {
		fmt.Println("no overlaps")
	}
	return nil
}

// alloc finds the first free /bits block inside pool that overlaps none of
// the already-used prefixes. This is the core of every IPAM system.
func alloc(poolStr, bitsStr string, usedStr []string) error {
	pool, err := netip.ParsePrefix(poolStr)
	if err != nil {
		return err
	}
	bits, err := strconv.Atoi(bitsStr)
	if err != nil {
		return err
	}
	used, err := parseAll(usedStr)
	if err != nil {
		return err
	}
	for cand := netip.PrefixFrom(pool.Masked().Addr(), bits); pool.Contains(cand.Addr()); {
		free := true
		for _, u := range used {
			if cand.Overlaps(u) {
				free = false
				break
			}
		}
		if free {
			fmt.Println("allocated", cand)
			return nil
		}
		next := lastAddr(cand).Next()
		if !next.IsValid() {
			break
		}
		cand = netip.PrefixFrom(next, bits)
	}
	return fmt.Errorf("pool %s has no free /%d", pool, bits)
}

func parseAll(ss []string) ([]netip.Prefix, error) {
	var out []netip.Prefix
	for _, s := range ss {
		p, err := netip.ParsePrefix(s)
		if err != nil {
			return nil, err
		}
		out = append(out, p.Masked())
	}
	return out, nil
}

// lastAddr sets every host bit to 1: the broadcast address for IPv4.
func lastAddr(p netip.Prefix) netip.Addr {
	b := p.Masked().Addr().AsSlice()
	for i := p.Bits(); i < len(b)*8; i++ {
		b[i/8] |= 1 << (7 - uint(i%8))
	}
	a, _ := netip.AddrFromSlice(b)
	return a
}

func mask4(bits int) netip.Addr {
	var b [4]byte
	for i := 0; i < bits; i++ {
		b[i/8] |= 1 << (7 - uint(i%8))
	}
	return netip.AddrFrom4(b)
}
```

Run it against the worked examples:

```text
$ go run ./subnet info 192.168.10.77/26
input        192.168.10.77
network      192.168.10.64/26
netmask      255.255.255.192
broadcast    192.168.10.127
host range   192.168.10.65 - 192.168.10.126
addresses    64 total, 62 usable
private?     true

$ go run ./subnet split 10.0.0.0/16 24
10.0.0.0/24
10.0.1.0/24
...
... 256 subnets in total
```

**What to notice in the code:**

- `p.Masked()` *is* "AND the address with the mask". `lastAddr` *is* "set
  every host bit to 1". The block-size shortcut from this chapter is the same
  bit arithmetic done in your head.
- The usable count subtracts 2 only for IPv4 prefixes of `/30` or larger.
  A `/31` (point-to-point links, RFC 3021) and a `/32` (a single host) don't
  sacrifice a network/broadcast address, and IPv6 has no broadcast at all.
- `math/big` is needed because a single IPv6 `/64` holds 2⁶⁴ addresses,
  which doesn't fit in a signed 64-bit integer.

**Exercises:**

1. Run `info` on all of Set A and Set B from the practice above, and compare
   with your paper answers.
2. Add a `contains` command: `subnet contains 192.168.1.0/25 192.168.1.200`
   prints yes/no. This is the "can I reach it directly?" test from this
   chapter, the same decision your OS makes before choosing ARP or the
   gateway.
3. Add `-json` output. You'll reuse it in Chapter 61's IPAM exercise.

### Common confusions

- **"/24 vs 255.255.255.0 — which should I use?"** They're identical. `/24` is
  the modern shorthand and is what you'll mostly see. Be able to read both.
- **"Why 'minus 2'?"** Network address and broadcast address are reserved
  (Chapter 9). Exception: a `/31` is used for point-to-point links between two
  routers and gives you 2 usable addresses (RFC 3021), and `/32` means one single
  host.
- **"What about Class A/B/C?"** Historical. Before 1993, the mask was implied by
  the first octet. **CIDR replaced this** and masks are now always explicit.
  Certification exams still mention classes; real networks don't care. If someone
  says "a class C network", they usually just mean a `/24`.
- **"Smaller /number = bigger network?"** Yes, and it trips everyone up. `/24` is
  256 addresses; `/16` is 65,536. **Fewer network bits = more host bits = bigger
  network.**

### Check yourself

1. What does the `/24` in `192.168.1.50/24` actually count?
2. How do you calculate block size from a mask value?
3. Are `10.1.1.100/26` and `10.1.1.140/26` on the same network? Show your working.
4. Why is a `/25` network smaller than a `/24`?

### Further reading

- **Video (~20 min, the best free explanation):** "Subnetting Mastery" series —
  Practical Networking (Ed Harmoush) on YouTube. Watch the whole series; it uses
  exactly the block-size method above and drills it properly.
- **Video:** "Subnetting Made Easy" — NetworkChuck. More energetic, same content.
- **Practice site:** `subnettingpractice.com` or `subnetting.org` — endless
  generated problems with answers. **Do 20 of these.** It's the only way.
- **Book:** *CompTIA Network+ Study Guide* or Todd Lammle's *CCNA Study Guide* —
  both have excellent, drill-heavy subnetting chapters.
- **Tool:** `ipcalc` on the command line; `jodies.de/ipcalc` in a browser.

---

## Chapter 11 — Subnetting practice: designing a real network

### In one sentence

Given a block of addresses and a list of departments, split the block so each gets
what it needs without waste — this is called **VLSM**, and it's a standard
interview task.

### The problem

Your company has been given `192.168.10.0/24` (256 addresses). You must split it
between:

| Department | Devices needed |
|---|---|
| Staff laptops | 100 |
| Guest Wi-Fi | 50 |
| Servers | 12 |
| Printers/IoT | 10 |
| Link between two routers | 2 |

You cannot just give each a quarter — the sizes are wildly different. You need
**Variable Length Subnet Masking (VLSM)**: different mask sizes for different
subnets.

### The method

```
   STEP 1: For each requirement, find the SMALLEST subnet that fits.
           Remember: usable = 2^(host bits) - 2

   STEP 2: Sort them LARGEST FIRST.  (This is the crucial step -- see below.)

   STEP 3: Allocate from the start of your block, in order,
           stepping forward by each subnet's block size.
```

### Step 1 — size each subnet

| Need | Smallest fit | Why | Block size |
|---|---|---|---|
| 100 staff | **/25** | /26 gives 62 (too small), /25 gives 126 | 128 |
| 50 guests | **/26** | /27 gives 30 (too small), /26 gives 62 | 64 |
| 12 servers | **/28** | /29 gives 6 (too small), /28 gives 14 | 16 |
| 10 printers | **/28** | same reasoning | 16 |
| 2 for the link | **/30** | exactly 2 usable | 4 |

### Step 2 & 3 — allocate, largest first

```
   Start at 192.168.10.0

   1. STAFF      /25   block 128   ->  192.168.10.0   -  192.168.10.127
                                       network .0     broadcast .127
                                       usable  .1  to .126        (126 addresses)

   2. GUESTS     /26   block  64   ->  192.168.10.128 -  192.168.10.191
                                       usable  .129 to .190       (62)

   3. SERVERS    /28   block  16   ->  192.168.10.192 -  192.168.10.207
                                       usable  .193 to .206       (14)

   4. PRINTERS   /28   block  16   ->  192.168.10.208 -  192.168.10.223
                                       usable  .209 to .222       (14)

   5. ROUTER LINK /30  block   4   ->  192.168.10.224 -  192.168.10.227
                                       usable  .225 to .226       (2)

   LEFT OVER: 192.168.10.228 - 192.168.10.255  (28 addresses for growth)
```

### Why "largest first" matters

Watch what happens if you allocate the `/30` first:

```
   1. ROUTER LINK /30  ->  192.168.10.0  -  192.168.10.3

   2. STAFF /25 needs a block of 128, and a /25 must START on a multiple
      of 128.  The next multiple of 128 after .3 is .128.

      -> addresses .4 through .127 (124 addresses) are STRANDED and unusable.
```

**Subnets must start on a boundary that is a multiple of their own block size.**
That's why big ones go first. Remember this; it's the point of the exercise.

### Practice (45 min)

**Problem 1.** You have `10.20.0.0/24`. Allocate for:
- Sales: 60 devices
- Engineering: 25 devices
- Management: 10 devices
- Wi-Fi guests: 100 devices
- Two point-to-point router links: 2 devices each

Give the network address, mask, usable range, and broadcast for each.

**Problem 2.** You have `172.16.0.0/22` (1024 addresses). Allocate for:
- Floor 1: 200 devices
- Floor 2: 200 devices
- Floor 3: 100 devices
- Servers: 50 devices
- Management VLAN: 20 devices

**Problem 3 — reverse engineering.** Given these three subnets, what is the
smallest single block that contains all of them?

```
   192.168.4.0/24
   192.168.5.0/24
   192.168.6.0/24
   192.168.7.0/24
```

*(Answers to all three: Appendix C.)*

**Verify with a tool:**

```bash
python3 - <<'EOF'
import ipaddress
block = ipaddress.ip_network("192.168.10.0/24")
# subnet_of / supernet helpers are handy for checking your work
for pfx in [25, 26, 28, 28, 30]:
    sub = next(block.subnets(new_prefix=pfx))
    print(pfx, sub, list(sub.hosts())[0], list(sub.hosts())[-1], sub.broadcast_address)
    # then manually advance -- this is the exercise
EOF
```

### Build it in Go (30 min) — check your design with code

Use the `subnet` tool from Chapter 10 to check the design you just did by
hand:

```text
$ go run ./subnet overlap 10.0.0.0/25 10.0.0.128/26 10.0.0.192/27 10.0.0.224/30
no overlaps

$ go run ./subnet overlap 10.0.0.0/16 10.0.128.0/20 172.16.0.0/12
OVERLAP  10.0.0.0/16  <->  10.0.128.0/20

$ go run ./subnet alloc 10.0.0.0/16 22 10.0.0.0/24 10.0.4.0/22
allocated 10.0.8.0/22
```

`alloc` is "largest first" turned into code: walk the pool in steps of the
requested size and return the first block that overlaps nothing already
used. Because it only ever proposes **aligned** blocks (`/22`s start on
multiples of 4 in the third octet), it can't produce an invalid subnet. That
is the mistake people make most often by hand.

**Exercise:** write `plan` that takes a pool and a list of host counts
(`subnet plan 10.0.0.0/24 100 50 20 2`), sorts them largest first, works out
the prefix length each needs (hosts + 2, rounded up to a power of two), and
allocates them in order with `alloc`. You have just automated this chapter.
Chapter 61 grows it into an IPAM.

### Common confusions

- **"Can I start a /26 at .100?"** No. A `/26` (block 64) must start at 0, 64,
  128, or 192. Subnets align to their own size.
- **"I have addresses left over — did I do it wrong?"** Leftovers are fine and
  usually good (room to grow). *Stranded* addresses in the middle are the mistake.
- **"Why not give everyone a /24?"** In private networks you often can and do —
  address space is free. VLSM matters when space is limited (public addresses,
  cloud VPCs with fixed CIDRs, or exams).

### Check yourself

1. Why do you allocate the largest subnet first?
2. What's the smallest subnet that fits 30 devices? 60? 2?
3. Can a `/28` subnet start at `192.168.1.20`? Why not?

### Further reading

- **Practice:** `subnettingpractice.com` — generate VLSM problems until it's
  boring. That's the goal.
- **Video:** "VLSM Explained" — Practical Networking or Keith Barker.
- **Book:** Todd Lammle's *CCNA Study Guide* — the subnetting/VLSM chapters are
  drill-based and excellent, even if you never sit the exam.

---

## Chapter 12 — Private addresses and NAT

### In one sentence

There aren't enough IPv4 addresses in the world, so your home devices use
"private" addresses that mean nothing on the internet, and your router rewrites
them on the way out — that rewriting is called **NAT**.

### The problem

IPv4 has ~4.3 billion addresses (Chapter 9). There are far more devices than that.
Your house alone might have 20. Yet it works. How?

### The idea, in plain language

Like phone extensions in an office.

```
   The office has ONE public phone number:   +44 20 7946 0000
   Inside, everyone has an extension:        101, 102, 103...

   Outbound: you call out; the other party sees the OFFICE number.
   Inbound:  the receptionist remembers who called out and connects
             the return call to the right extension.
```

Your router is the receptionist. Your devices have extensions (private IPs).

### The private address ranges (memorise these)

These three ranges (defined in **RFC 1918**) are reserved for internal use and are
**never routed on the public internet**:

```
   10.0.0.0     to  10.255.255.255      (10.0.0.0/8)      ~16.7 million addresses
   172.16.0.0   to  172.31.255.255      (172.16.0.0/12)   ~1 million
   192.168.0.0  to  192.168.255.255     (192.168.0.0/16)  ~65 thousand
```

You'll also meet:
```
   100.64.0.0/10    Carrier-Grade NAT (CGNAT) -- your ISP NATs you too
```

**Everyone uses the same private addresses.** There are millions of homes with a
`192.168.1.1` router. That's fine, because those addresses never leave the house.

### How NAT actually works

Say your laptop (`192.168.1.50`) opens a connection to a web server:

```
   INSIDE YOUR HOUSE                  |          THE INTERNET
                                      |
   laptop 192.168.1.50                |
        |                             |
        |  from: 192.168.1.50 : 51000 |
        |  to:   93.184.216.34 : 443  |
        v                             |
   +---------+                        |
   | ROUTER  |  rewrites the source   |
   | (NAT)   |----------------------->|  from: 203.0.113.7 : 62000
   +---------+                        |  to:   93.184.216.34 : 443
        ^                             |
        |  and remembers:             |
        |  203.0.113.7:62000  <-->  192.168.1.50:51000
        |                             |
        |  reply arrives for          |
        |  203.0.113.7:62000          |
        |  -> looks it up             |
        |  -> rewrites back to        |
        |     192.168.1.50:51000      |
```

The router keeps a **NAT table** of these mappings. When several devices connect
to the same server, the router gives each a different outside **port number** so
it can tell the replies apart. (Ports are Chapter 19 — for now, just a number that
distinguishes connections.)

Because it translates ports as well as addresses, this is properly called
**PAT** (Port Address Translation) or **NAT overload**. Almost everyone just says
"NAT".

### What NAT gives you, and what it costs

**Good:**
- Many devices share one public address (the whole point).
- **An accidental firewall**: unsolicited traffic from the internet has no
  mapping to look up, so it's dropped. This is why home devices aren't
  immediately attacked.

**Bad:**
- **Incoming connections don't work** without explicit configuration. You can't
  run a server at home without **port forwarding**.
- **Peer-to-peer is hard.** Two devices both behind NAT can't easily connect —
  which is why video calls need extra machinery (STUN/TURN servers).
- **The internet stopped being end-to-end.** Any protocol that puts an IP address
  *inside* its data (old FTP, SIP) breaks unless the router specially understands
  it.
- **CGNAT makes it worse.** Many ISPs now NAT *you* as well, so you're behind two
  layers and your public IP is shared with strangers. Port forwarding becomes
  impossible.

### Practice (20 min)

**A. See both of your addresses:**

```bash
# your PRIVATE address (what your device thinks it is)
ip -4 addr show | grep inet        # Linux
ipconfig getifaddr en0             # macOS

# your PUBLIC address (what the internet sees)
curl -s ifconfig.me ; echo
```

They're different. **Everyone in your house shows the same public one.** Try it on
your phone (on Wi-Fi) to confirm.

**B. Find your router's NAT table** (if it's a Linux box or you can log into your
home router's web UI — look for "NAT", "connections", or "active sessions"):

```bash
# On a Linux machine acting as a router:
sudo conntrack -L | head
```

**C. Prove NAT blocks inbound.** Start a simple server on your laptop:

```bash
python3 -m http.server 8000
```

- From another device *in your house*: `http://<your-private-ip>:8000` → works.
- From outside (mobile data, or ask a friend): `http://<your-public-ip>:8000` →
  **fails**, because no NAT mapping exists. That's the accidental firewall.

**D. Check whether you're behind CGNAT:** compare `curl ifconfig.me` with the WAN
address shown in your router's admin page. If your router's WAN address is in
`100.64.x.x`–`100.127.x.x`, you are behind carrier-grade NAT.

### Build it in Go (15 min) — see your own NAT

Two questions decide whether you're behind NAT: *what address does my
machine think it has?* and *what address does the internet see?* This program
asks both.

```go
// whoami: show the NAT in front of you by comparing the address your machine
// THINKS it has with the address the internet SEES.
//
//	go run ./whoami
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func main() {
	// 1. Which local address would the kernel use to reach the internet?
	//    "Dialing" UDP sends no packets -- it only asks the kernel to pick a
	//    route and a source address, which we can then read back.
	c, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		fmt.Println("no route to the internet:", err)
		return
	}
	local := c.LocalAddr().(*net.UDPAddr).IP
	c.Close()
	fmt.Println("private (what you think you are): ", local)

	// 2. Ask a server on the internet which address our packets arrived from.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.ipify.org", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("could not reach api.ipify.org:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64))
	public := net.ParseIP(strings.TrimSpace(string(body)))
	fmt.Println("public  (what the internet sees): ", public)

	// 3. Interpret.
	switch {
	case public == nil:
		fmt.Println("unexpected answer:", string(body))
	case local.Equal(public):
		fmt.Println("=> no NAT: your machine holds a public address directly")
	case isCGNAT(local):
		fmt.Println("=> carrier-grade NAT: even your 'router' address is shared (100.64.0.0/10)")
	default:
		fmt.Println("=> NAT: a router rewrote your source address on the way out")
	}
}

func isCGNAT(ip net.IP) bool {
	_, cgnat, _ := net.ParseCIDR("100.64.0.0/10")
	return cgnat.Contains(ip)
}
```

```text
$ go run ./whoami
private (what you think you are):  192.168.1.2
public  (what the internet sees):  122.171.18.134
=> NAT: a router rewrote your source address on the way out
```

**What to notice:** `net.Dial("udp", ...)` **sends no packet**. UDP has no
handshake, so "connecting" a UDP socket only asks the kernel to pick a
route and a source address. Reading `LocalAddr()` afterwards is a standard
trick for finding "my outbound IP" without parsing routing tables.

**Exercises:**

1. Run it on home Wi-Fi, then on a phone hotspot. On mobile networks you'll
   often get the carrier-grade NAT result (`100.64.0.0/10`), or a public
   address that thousands of other customers share at the same moment.
2. Run it inside a Docker container (`docker run --rm -v "$PWD":/w -w /w golang:1.26 go run ./whoami`).
   You'll see an *extra* layer of NAT: container → host → router.
3. Why does this matter for software? Read the "rate limiting by IP" point in
   the [OSI guide's Layer 3 chapter](../real-life-example-osi.md#part-5-layer-3-network), then decide what *your*
   service should use as a client identity.

### Common confusions

- **"NAT is a firewall."** It behaves like one for inbound traffic, but it isn't
  designed as a security control and shouldn't be your only one. Use an actual
  firewall (Chapter 38).
- **"Why does my private IP change?"** DHCP leases (Chapter 17). Set a
  reservation in your router if you need it stable.
- **"Two networks both use 192.168.1.0/24 — can I VPN between them?"** Not
  easily; the addresses collide. This is a very common real-world headache, and
  the reason experienced people pick unusual ranges like `10.37.129.0/24` for
  home labs.
- **"Does IPv6 use NAT?"** Generally no — there are enough addresses that every
  device can have a public one. Security comes from a stateful firewall instead
  (Chapter 16).

### Check yourself

1. Name the three private IPv4 ranges.
2. What does the router store in its NAT table, and why does it need port numbers?
3. Why can't someone connect *to* your laptop from the internet by default?
4. What is CGNAT and why is it annoying?

### Further reading

- **Video (~12 min):** "NAT Explained" — Practical Networking. Also their
  "What is CGNAT?"
- **Article:** "What is NAT?" — Cloudflare Learning Center.
- **RFC 1918** — two pages, defines the private ranges. Worth skimming.
- **Article:** "The trouble with CGNAT" — search APNIC blog; good context on why
  IPv6 matters.

---

## Chapter 13 — Routing: how a packet finds its way

### In one sentence

Every machine has a **routing table** — a list of "to reach *this* range of
addresses, send to *that* next hop" — and a packet is passed from router to
router, each making one local decision, until it arrives.

### The problem

Your packet needs to reach a server 8,000 km away. No single machine knows the
whole path. How does it get there?

### The idea, in plain language

**Nobody knows the whole route. Everybody knows the next step.**

```
   You, driving, asking at each junction:
      "Which way to Rome?"  ->  "Take the second exit."
   You don't need the whole route. Just the next turn, repeatedly.
```

Each router asks one question: *given this destination address, which of my
interfaces should I send it out of?* Then it's someone else's problem.

### How it actually works: the routing table

Look at yours right now:

```bash
ip route              # Linux
netstat -rn           # macOS / BSD
route print           # Windows
```

You'll see something like:

```
default via 192.168.1.1 dev wlan0
192.168.1.0/24 dev wlan0 proto kernel scope link src 192.168.1.50
```

Read each line as: **"To reach *this destination*, send it *this way*."**

| Line | Means |
|---|---|
| `192.168.1.0/24 dev wlan0 ... src 192.168.1.50` | "Addresses in 192.168.1.x are on my own Wi-Fi network — send **directly** (ARP for them)." |
| `default via 192.168.1.1 dev wlan0` | "**Anything else**, hand to 192.168.1.1 and let it worry." |

`default` means `0.0.0.0/0` — "every possible address". It's the catch-all, and
it's why your home network needs only two rules to reach the entire internet.

### The one rule: longest prefix match

If several routes match, **the most specific one wins** — the one with the biggest
`/number`.

```
   A packet for 10.1.2.3 arrives. The table contains:

      0.0.0.0/0        <- matches (matches everything)
      10.0.0.0/8       <- matches
      10.1.0.0/16      <- matches
      10.1.2.0/24      <- matches   <=== WINNER (most specific)

   The /24 route is chosen, regardless of the order they're listed in.
```

Think of it as "the most specific instruction wins", like a company policy where
a department rule overrides the company-wide rule.

### What a router does to each packet

```
   1. Check the packet isn't corrupt (header checksum)
   2. Decrement the TTL ("hop counter"). If it hits 0: DROP it and
      send an error back  (this is what makes traceroute possible -- Ch 14)
   3. Look up the destination IP -> longest prefix match -> next hop + interface
   4. Rewrite the LINK-LAYER header (new source/destination MAC for the next hop)
   5. Send it
```

Note step 4: **the MAC addresses change at every hop; the IP addresses don't.**
That's the Chapter 2 insight, made concrete. And note what a router does *not* do:
it keeps no memory of your connection, doesn't retransmit, doesn't look inside.
That's the end-to-end principle (Chapter 3).

### How routers learn routes

| Source | How | Where |
|---|---|---|
| **Directly connected** | "I have an interface on this network" | automatic |
| **Static** | A human types `ip route add ...` | small networks, home |
| **OSPF / IS-IS** | Routers inside one organisation flood each other with link info and each computes shortest paths | company networks |
| **BGP** | Routers *between* organisations exchange "I can reach these address ranges via this path" | **this is what runs the internet** |

**BGP** deserves a note: there are ~75,000 independent networks (called
*autonomous systems*) on the internet, and BGP is how they tell each other which
addresses they can reach. It's held together by trust and business agreements more
than by technology, which is why occasional BGP mistakes take large parts of the
internet offline.

### Practice (30 min)

**A. Read your own routing table** and explain each line out loud.

```bash
ip route          # or: netstat -rn
```

**B. Ask your machine to decide** (Linux — this is a great tool):

```bash
ip route get 8.8.8.8
ip route get 192.168.1.99      # substitute an address on your own network
ip route get 127.0.0.1
```

It prints exactly which route it chose, which interface, and which source
address. Compare all three answers — one goes via the gateway, one goes direct,
one stays local.

**C. See the hops for real:**

```bash
traceroute -n example.com          # macOS/Linux
tracert example.com                # Windows
mtr -rwzbc 50 example.com          # better: continuous, shows loss per hop
```

Each line is one router along the path. Chapter 14 explains exactly how this
works and — importantly — how to read it without being misled.

**D. Watch TTL decrease.** Ping something far away and look at the TTL in the
reply:

```bash
ping -c 1 example.com
# "64 bytes from ...: icmp_seq=0 ttl=52 time=91 ms"
```

Most systems start TTL at 64. A reply arriving with `ttl=52` means it crossed
about `64 − 52 = 12` routers on the way back.

### Common confusions

- **"Default gateway vs router?"** The default gateway is *the router your machine
  sends everything unknown to*. Same box, different role name.
- **"Does the router know the whole path?"** No. It knows one step. This is the
  central idea of the chapter.
- **"Why does traceroute show different paths each run?"** Load balancing across
  multiple equal-cost links (ECMP). Normal.
- **"Can the path back be different from the path there?"** Yes, and it usually
  is. This is called **asymmetric routing** and it confuses both humans and
  firewalls.

### Check yourself

1. What does `default via 192.168.1.1` mean?
2. Given routes `10.0.0.0/8` and `10.1.2.0/24`, which is used for `10.1.2.50`?
3. What does a router change in the packet, and what does it leave alone?
4. What is BGP for?

### Further reading

- **Video (~15 min):** "How Routers Work" and "Routing Tables Explained" —
  Practical Networking.
- **Video (~10 min):** "BGP Explained" — search Cloudflare's or NetworkChuck's.
- **Article:** "What is BGP?" — Cloudflare Learning Center. Also read their
  post-mortems of real BGP incidents; they're gripping and educational.
- **Book:** *Computer Networking: A Top-Down Approach*, Chapter 5 (the Network
  Layer / routing).
- **Interactive:** `bgp.he.net` — look up any IP address or company and see how
  the internet routes to it.

---

## Chapter 14 — ICMP: what `ping` and `traceroute` really do

### In one sentence

**ICMP** is the internet's error-and-diagnostics messenger, and `ping` and
`traceroute` are two clever tricks built on it.

### What ICMP is

IP itself has no way to say "that didn't work". **ICMP** (Internet Control Message
Protocol) fills that gap. It is *not* a transport protocol — it has no ports and
carries no user data. It just reports conditions.

Messages you should recognise:

| Type | Name | When you see it |
|---|---|---|
| 8 / 0 | Echo Request / Reply | `ping` |
| 3, code 1 | Destination **host** unreachable | ARP failed on the last hop — device is off/absent |
| 3, code 3 | Destination **port** unreachable | Nothing is listening on that UDP port |
| 3, code 4 | **Fragmentation needed** | Packet too big — critical for Chapter 15 |
| 3, code 13 | Administratively prohibited | A firewall rejected it |
| 11 | **Time exceeded** | TTL hit zero — this is what powers traceroute |

### How `ping` works

Dead simple: send an **Echo Request**, the target's operating system replies with
an **Echo Reply**, measure the round trip.

```bash
ping -c 5 example.com
```
```
64 bytes from 93.184.216.34: icmp_seq=1 ttl=52 time=91.2 ms
64 bytes from 93.184.216.34: icmp_seq=2 ttl=52 time=90.8 ms
...
--- example.com ping statistics ---
5 packets transmitted, 5 received, 0% packet loss
round-trip min/avg/max/stddev = 90.8/91.4/92.1/0.5 ms
```

What each part tells you:
- **`time=91.2 ms`** — round trip. Mostly determined by *distance*.
- **`icmp_seq`** — sequence number. Gaps mean lost packets.
- **`ttl=52`** — how much hop budget was left (Chapter 13).
- **`stddev` / `mdev`** — **jitter**. Low = stable. High = congested or wireless.
- **`0% packet loss`** — the headline number.

**Important caveats** (these trip people up constantly):

1. **"Ping fails" does not mean "host is down".** Many hosts and firewalls block
   or rate-limit ICMP deliberately. Always confirm with something else.
2. **Never judge a router by pinging the router itself.** Routers deprioritise
   traffic addressed *to them* — a router can show 30% ping loss while forwarding
   your real traffic perfectly.
3. **Ping measures round trip, not one way.** They're often not symmetric.

### How `traceroute` works (the clever trick)

`traceroute` exploits the TTL counter from Chapter 13:

```
   Send a packet with TTL = 1
      -> the FIRST router decrements it to 0, drops it,
         and sends back "Time Exceeded"
      -> that error reveals router 1's address!

   Send a packet with TTL = 2
      -> router 1 passes it (TTL 2->1), router 2 drops it
      -> reveals router 2

   Send TTL = 3, 4, 5, ... until a packet finally reaches the destination.
```

Each row of traceroute output is one router that had to drop your packet.

```bash
traceroute -n example.com
```
```
 1  192.168.1.1        1.2 ms   1.1 ms   1.3 ms      <- your router
 2  10.20.30.1         8.4 ms   8.2 ms   8.9 ms      <- your ISP
 3  * * *                                            <- didn't reply (see below)
 4  62.115.11.22      22.1 ms  21.8 ms  22.4 ms
...
12  93.184.216.34     91.2 ms  90.9 ms  91.5 ms      <- destination
```

### How to read traceroute without being misled

This is one of the most misinterpreted tools in networking. Three rules:

**1. `* * *` usually means nothing is wrong.** It means that router chose not to
reply. If later hops respond fine, the path is fine.

**2. High latency at one hop that *doesn't* persist is fake.** Look at this:

```
 5   10.0.0.1     150 ms      <- looks terrible!
 6   10.0.0.2      22 ms      <- but the NEXT hop is fast
 7   destination   24 ms
```

Hop 5 isn't slow — it just deprioritises replies to itself. **Only believe a
latency or loss figure if it persists on every subsequent hop.**

**3. You only see the forward path.** The return path may be completely different,
and a problem there will show up in your numbers without being visible in the
list.

### `mtr` — use this instead

`mtr` combines ping and traceroute, running continuously:

```bash
mtr -rwzbc 100 example.com     # 100 packets, report mode
```

It shows loss% and latency **per hop over time**, which makes rule 2 above
obvious: loss that appears at hop 5 and vanishes at hop 6 is cosmetic; loss that
appears at hop 5 and continues to the destination is real.

### Do not block all ICMP

A depressingly common "hardening" mistake. If you block ICMP entirely:
- You break **Path MTU Discovery** (Chapter 15) — connections hang mysteriously.
- On IPv6 you break the network completely (Chapter 16) — ICMPv6 does address
  resolution.

Block what you must, but always allow **type 3** (destination unreachable) and
**type 11** (time exceeded).

### Practice (30 min)

**A. Ping a few destinations and compare:**

```bash
ping -c 10 192.168.1.1              # your router: should be <5 ms, near-zero jitter
ping -c 10 1.1.1.1                  # a big public DNS: typically 5-30 ms
ping -c 10 example.com              # a specific server, maybe far away
```

Record the min/avg/max/mdev for each in `notes.md`. **You are measuring
distance.** Roughly: light in fibre travels ~200 km per millisecond, and a round
trip doubles it — so 90 ms round trip ≈ 9,000 km of fibre, plus routing overhead.

**B. Traceroute somewhere far away** and count the hops:

```bash
traceroute -n www.gov.au        # or any distant site
```

Can you spot where it leaves your country? (Hostnames often hint at cities —
try without `-n`.)

**C. Run `mtr` for two minutes** to a site you use daily, and identify: which hop
is your router, which is your ISP, where loss (if any) begins.

**D. Watch ICMP in tcpdump:**

```bash
sudo tcpdump -i any -n icmp
# then, in another terminal:  ping -c 3 1.1.1.1
```

You'll see the request and reply pairs.

**E. Trigger an error deliberately:**

```bash
# ask for a UDP port nothing is listening on
sudo tcpdump -i any -n icmp &
nc -u -w1 127.0.0.1 9   # or any closed UDP port
```
You should see an ICMP **port unreachable**.

### Common confusions

- **"Ping is blocked so the server is down."** No — test the actual service:
  `curl`, `nc -zv host port`, or `traceroute -T -p 443 host` (TCP traceroute).
- **"traceroute shows 20 hops so it's slow."** Hop count barely matters; distance
  and congestion do.
- **"ICMP is a security risk, block it all."** See above. Block selectively.

### Check yourself

1. How does traceroute discover each router along the path?
2. You see `* * *` at hop 4 but hops 5–12 respond normally. Is there a problem?
3. Why shouldn't you trust high latency at a single middle hop?
4. Name two things that break if you block all ICMP.

### Further reading

- **Video:** "How Traceroute Works" — Practical Networking. Excellent animation of
  the TTL trick.
- **Article:** "Understanding MTR output" — search the Linode or DigitalOcean
  guide; both are clear.
- **Article:** "ICMP explained" — Cloudflare Learning Center.
- **Blog:** Julia Evans, "How does traceroute work?" (jvns.ca).

---

## Chapter 15 — MTU: the maximum size of a packet

### In one sentence

Every link has a maximum packet size (usually **1500 bytes**), and when a packet
is too big for the next link, something has to give — this causes one of the most
confusing bugs in networking.

### The problem

You send a 1 MB file. It can't go as one packet — Ethernet's payload maxes out at
1500 bytes (Chapter 5). So it's split into ~700 packets. Fine. But what if a link
somewhere along the path has a *smaller* limit than yours?

### The key numbers

```
   MTU  (Maximum Transmission Unit)  = the biggest IP packet a link will carry
   MSS  (Maximum Segment Size)       = the biggest chunk of DATA inside it

   Standard Ethernet:   MTU = 1500
                        - 20 bytes IP header
                        - 20 bytes TCP header
                        ---------------------
                        MSS = 1460 bytes of your actual data
```

Common MTUs you'll meet:

| Link | MTU |
|---|---|
| Standard Ethernet / Wi-Fi | **1500** |
| PPPoE (some DSL) | 1492 |
| WireGuard VPN | 1420 |
| IPsec VPN | ~1400 |
| Cloud overlay networks (VXLAN) | 1450 |
| Data-centre "jumbo frames" | 9000 |
| IPv6 minimum every link must support | 1280 |

Notice the VPN entries. **A VPN wraps your packet inside another packet**, and
that wrapper costs bytes. This is why VPNs cause MTU problems so often.

### What happens when a packet is too big

Two possibilities:

```
   CASE 1: the sender said "you may split this" (DF bit not set)
     -> the router chops it into FRAGMENTS
     -> slow, fragile, and often blocked by firewalls

   CASE 2: the sender said "DO NOT FRAGMENT" (DF bit set -- TCP always does this)
     -> the router DROPS the packet
     -> and sends back ICMP "Fragmentation Needed", saying the correct size
     -> the sender shrinks its packets and retries
```

Case 2 is called **Path MTU Discovery (PMTUD)**, and it works well — *as long as
that ICMP message gets through*.

### The classic bug: the PMTU black hole

Someone blocks all ICMP on a firewall "for security" (Chapter 14). Now:

```
   1. TCP handshake works fine        (those packets are tiny)
   2. Your request works fine          (also small)
   3. The server sends a full-size 1500-byte response packet
   4. A router on the path can't fit it, drops it, sends ICMP...
   5. ...which is BLOCKED. The sender never learns.
   6. The sender retransmits the same too-big packet. Forever.

   RESULT: the connection CONNECTS and then HANGS.
```

**Learn this fingerprint. It is one of the most common "impossible" bugs:**

> *"SSH connects then freezes."*
> *"Small web pages load, big ones hang."*
> *"`git clone` stalls at a few KB."*
> *"It works on the office Wi-Fi but not over the VPN."*

All four are usually MTU.

### How to diagnose it in 30 seconds

`ping` can send a specific size with "don't fragment" set:

```bash
# Linux
ping -M do -s 1472 -c 2 example.com     # 1472 + 8 (ICMP) + 20 (IP) = 1500

# macOS
ping -D -s 1472 -c 2 example.com
```

- **Works?** Path MTU is at least 1500. Good.
- **"Message too long" / "Frag needed"?** Something on the path is smaller.

Now **binary search** for the real limit:

```bash
for s in 1472 1452 1422 1372 1272; do
  echo -n "payload $s: "
  ping -M do -s $s -c 1 -W 2 example.com >/dev/null 2>&1 && echo OK || echo TOO BIG
done
```

The largest size that works, plus 28, is your path MTU.

Or just use:
```bash
tracepath example.com        # Linux -- reports the MTU per hop, no root needed
```

### How to fix it

| Fix | When |
|---|---|
| **Allow ICMP type 3** through your firewalls | Always. This is the correct fix. |
| **MSS clamping** on the router/VPN — rewrite the MSS in the handshake so both ends agree to smaller packets | Standard on VPN gateways |
| **Lower the interface MTU** on the affected machine (`ip link set eth0 mtu 1400`) | Quick local workaround |
| **Enable `tcp_mtu_probing`** (Linux) so TCP finds the limit itself without ICMP | Good defensive setting |

### Practice (25 min)

**A. Find your own path MTU** to three destinations: your router, a nearby site, a
distant site. Record them.

**B. Break it on purpose** (safe, and instructive):

```bash
# Linux -- lower your MTU, see what happens, then put it back
ip link show eth0 | grep mtu           # note the current value
sudo ip link set eth0 mtu 1200
curl -s -o /dev/null -w '%{http_code}\n' https://example.com
sudo ip link set eth0 mtu 1500         # PUT IT BACK
```

Small MTU still works (just less efficiently) because PMTUD is functioning. The
*breakage* comes from a mismatch plus blocked ICMP.

**C. See fragmentation:**

```bash
sudo tcpdump -i any -n 'ip[6:2] & 0x3fff != 0'   # captures only fragments
# in another terminal:
ping -s 4000 -c 1 192.168.1.1      # a 4000-byte ping must be fragmented
```

**D. If you use a VPN**, run the MTU test with the VPN on and off. Compare.

### Common confusions

- **"Bigger MTU is always better."** Only if *every* device on the path agrees.
  Jumbo frames (9000) are great inside a data centre and useless across the
  internet, where something will be 1500.
- **"MTU vs MSS?"** MTU is the whole packet including headers; MSS is just the
  data. `MSS = MTU − 40` for normal TCP over IPv4.
- **"Why doesn't this fix itself?"** It does — via PMTUD — unless someone blocks
  the ICMP that makes PMTUD work.

### Check yourself

1. What's the standard Ethernet MTU, and what's the resulting TCP MSS?
2. What does a router do with a too-big packet that has the DF bit set?
3. Describe the symptoms of a PMTU black hole.
4. Why do VPNs so often cause MTU problems?

### Further reading

- **Article (excellent):** "Path MTU Discovery" — Cloudflare blog, and their
  post on "Fixing an old hack" about MTU issues in the wild.
- **Video:** "MTU and MSS Explained" — Practical Networking.
- **Blog:** search "MTU black hole troubleshooting" — many good war stories;
  reading a few builds pattern recognition fast.

---

## Chapter 16 — IPv6: why it exists and how it differs

### In one sentence

IPv6 replaces the 32-bit address with a 128-bit one — enough addresses for
everything forever — and simplifies several things along the way, but it's not
just "IPv4 with more digits".

### Why it exists

IPv4 ran out of addresses. Officially. The global pool was exhausted between 2011
and 2019 depending on region. NAT (Chapter 12) bought time, at the cost of making
the internet more complicated and less peer-to-peer.

```
   IPv4:  32 bits  ->  4,300,000,000 addresses
   IPv6: 128 bits  ->  340,282,366,920,938,463,463,374,607,431,768,211,456
```

That's more addresses than there are grains of sand on Earth, by a very large
margin. The design goal was "never do this again".

### Reading an IPv6 address

They look intimidating and are simpler than they appear.

```
   Full:        2001:0db8:0000:0000:0000:ff00:0042:8329
                \___/ \___/ \___/ \___/ \___/ \___/ \___/ \___/
                 8 groups of 4 hex digits = 128 bits

   Rule 1: drop leading zeros in each group
                2001:db8:0:0:0:ff00:42:8329

   Rule 2: replace ONE run of all-zero groups with "::"
                2001:db8::ff00:42:8329

   (You may only use "::" once, or it would be ambiguous.)
```

The structure, for a typical address:

```
   2001:db8:abcd : 0012 : 0000:0000:0000:0001
   \____________/  \__/   \___________________/
     given to you  your     the device part
     by your ISP   subnet   (64 bits -- huge)
      (48 bits)   (16 bits)

   Almost every IPv6 LAN is a /64. That's the convention -- don't fight it.
```

### The types you'll see

| Starts with | Name | Meaning |
|---|---|---|
| `::1` | Loopback | Same as IPv4's `127.0.0.1` |
| `fe80::` | **Link-local** | Auto-generated on **every** IPv6 interface, always present, never routed off the local link. Used for the plumbing. |
| `fd00::` | Unique Local (ULA) | The IPv6 "private" range, like `192.168.x.x` |
| `2000::` – `3fff::` | Global unicast | Real, internet-routable addresses |
| `ff00::` | Multicast | IPv6 has **no broadcast at all** — multicast replaces it |

**Your device will have several IPv6 addresses at once** — a link-local, maybe a
ULA, and one or more global ones (often including temporary "privacy" addresses
that rotate daily). This is normal, not a misconfiguration.

### What's genuinely different

| | IPv4 | IPv6 |
|---|---|---|
| **Getting an address** | DHCP server hands them out | Usually **SLAAC**: the router advertises the network prefix, and the device makes up its own address. No server needed. |
| **Finding neighbours** | ARP | **NDP** (Neighbor Discovery), which runs over ICMPv6 |
| **Broadcast** | Yes | **None** — multicast only |
| **NAT** | Universal | Normally not used — every device gets a public address |
| **Fragmentation** | Routers can do it | **Only the sender** may fragment |
| **Header** | Variable length, has a checksum | Fixed 40 bytes, **no checksum** (lower layers handle it) |
| **Minimum MTU** | 576 | **1280** |

**The most important operational consequence:** because NDP runs on ICMPv6,
**blocking ICMPv6 breaks IPv6 entirely**. You cannot treat it like IPv4's ICMP.

### Running both at once

Almost everything today is **dual-stack** — it speaks both. When you connect to a
site, your machine gets both an IPv4 (`A` record) and IPv6 (`AAAA` record) answer
from DNS, and uses **Happy Eyeballs**: it starts connecting over both roughly
simultaneously and uses whichever answers first. That way, broken IPv6 costs you
~250 ms instead of a 30-second timeout.

### Practice (20 min)

**A. Do you have IPv6?**

```bash
ip -6 addr show                 # Linux
ifconfig | grep inet6           # macOS

curl -6 -s https://ifconfig.co ; echo    # your public IPv6 (fails if you have none)
curl -4 -s https://ifconfig.co ; echo    # your public IPv4
```

Also visit `test-ipv6.com` in a browser — it gives a clear score.

**B. Note that you always have a link-local address**, even with no IPv6 internet:

```bash
ip -6 addr show | grep fe80
```

Every IPv6-capable interface has one, always.

**C. Ping over IPv6:**

```bash
ping6 -c 3 2606:4700:4700::1111        # Cloudflare's IPv6 DNS
ping6 -c 3 ipv6.google.com

# link-local needs an interface specified, note the % syntax:
ping6 -c 2 fe80::1%eth0
```

**D. Watch Neighbor Discovery (IPv6's ARP):**

```bash
sudo tcpdump -i any -n 'icmp6 && ip6[40] >= 133 && ip6[40] <= 137'
```

Leave it running a minute. You'll see Router Advertisements and Neighbor
Solicitations — the IPv6 equivalents of Chapter 6's ARP.

**E. Compare DNS answers:**

```bash
dig example.com A +short
dig example.com AAAA +short
```

### Common confusions

- **"Do I need to learn IPv6?"** Yes. Mobile networks are largely IPv6-only
  internally, cloud providers use it, and it's in every certification. But you can
  learn it *after* you're solid on IPv4 — the concepts transfer.
- **"Why does my machine have five IPv6 addresses?"** By design (see above).
- **"Is IPv6 less secure without NAT?"** Different, not worse. You use a
  **stateful firewall** (default-deny inbound) instead of relying on NAT's
  accidental hiding. Your home router does this already.
- **"Can IPv4 talk to IPv6?"** Not directly. They're separate protocols. Devices
  run both (dual-stack), or use translation gateways (NAT64).

### Check yourself

1. Why does IPv6 exist?
2. What does `::` mean in an address, and why can you only use it once?
3. What is `fe80::1` and where would you find it?
4. Why is blocking ICMPv6 catastrophic when blocking ICMPv4 is merely bad?

### Further reading

- **Video series:** "IPv6 Fundamentals" — Practical Networking. Ed's IPv6 series
  is the gentlest good introduction available.
- **Site:** `test-ipv6.com` — test your own connectivity, with explanations.
- **Book:** *IPv6 Fundamentals* — Rick Graziani. Written for people who already
  know IPv4; clear and practical.
- **Article:** "Understanding IPv6" — APNIC blog has consistently good IPv6
  material.
- **Video:** "IPv6 - The Basics" — NetworkChuck (energetic, good for a first pass).

---

### End of Part 3 — Milestone check

- [ ] I can explain the two parts of an IP address
- [ ] **I can do subnetting (Set A) without looking anything up**
- [ ] I can design a VLSM allocation, largest first
- [ ] I can name the three private ranges and explain NAT
- [ ] I can read a routing table and explain longest prefix match
- [ ] **I have run traceroute/mtr and can read the output without being fooled**
- [ ] I can diagnose an MTU problem from its symptoms
- [ ] I can read an IPv6 address and know what `fe80::` means

**This is the hardest part of the guide.** If the subnetting box isn't ticked, go
back to Chapter 10 and do 20 practice problems. Everything else builds on it.

---

# Part 4 — Getting connected automatically

Two protocols do all the work you never think about: one gives your device an
address, the other turns names into addresses.

## Chapter 17 — DHCP: how your device gets an address

### In one sentence

When you join a network, your device shouts "does anyone here hand out
addresses?", and a **DHCP server** replies with an address, a gateway, and DNS
servers — a lease that expires and gets renewed.

### The problem

Someone has to assign your laptop an IP address, tell it the subnet mask, tell it
which router to use, and tell it where to look up names. Doing this by hand for
every device in an office (or a coffee shop) is impossible.

### The idea, in plain language

```
   YOU (new to the network, shouting to everyone):
       "Hello? Does anyone give out addresses?"

   DHCP SERVER (your router):
       "Yes! You can have 192.168.1.50."

   YOU (shouting again, so other servers know you've chosen):
       "Great, I'll take 192.168.1.50 from you."

   DHCP SERVER:
       "Confirmed. It's yours for 24 hours. Here's everything else you need:
        your mask is /24, your gateway is 192.168.1.1, your DNS is 192.168.1.1."
```

Four messages. This is called **DORA**, from the message names.

### How it actually works: DORA

```
   Client                                            Server
   (has NO address yet, so it uses 0.0.0.0
    and shouts to 255.255.255.255)

     |--- 1. DISCOVER  (broadcast) ------------------->|
     |    "Anyone out there? My MAC is aa:bb:cc:.."    |
     |                                                 |
     |<-- 2. OFFER  ---------------------------------- |
     |    "Have 192.168.1.50, mask /24, gw .1,         |
     |     DNS .1, lease 86400 seconds"                |
     |                                                 |
     |--- 3. REQUEST  (broadcast) -------------------->|
     |    "I accept 192.168.1.50, from YOU specifically"|
     |    (broadcast so any other DHCP servers know    |
     |     to withdraw their offers)                   |
     |                                                 |
     |<-- 4. ACK  ------------------------------------ |
     |    "Confirmed. It's yours."                     |
```

DHCP runs over **UDP**, ports **67** (server) and **68** (client). The first two
messages are broadcasts because the client doesn't have an address yet — the
chicken-and-egg problem solved by shouting.

### The lease: why your address changes

You don't own the address; you **lease** it.

```
   0%                    50% (T1)          87.5% (T2)         100%
   |----------------------|-----------------|-----------------|
   lease starts       RENEW: quietly    REBIND: ask ANY     lease expires:
   (e.g. 24 hours)    ask the same      server, by          you MUST stop
                      server for more   broadcast           using the address
                      time
```

Normally you renew at 50% and never notice. If your router reboots or you move
networks, the process starts over.

### The instant diagnosis: `169.254.x.x`

If nothing answers the DISCOVER, your device gives up and assigns itself an
address from `169.254.0.0/16` (called **APIPA** or link-local). Windows shows
"limited connectivity"; macOS shows "self-assigned IP".

> **Seeing `169.254.something` means: DHCP failed.** Either there's no DHCP
> server reachable, the cable/Wi-Fi isn't really connected, the DHCP pool is
> exhausted, or something is blocking broadcasts. This is one of the highest-value
> diagnostic facts in networking.

### What else DHCP gives you

Beyond the address, the server sends **options**:

| Option | What it is |
|---|---|
| 1 | Subnet mask |
| 3 | **Default gateway** (your router) |
| 6 | **DNS servers** |
| 51 | Lease time |
| 15 / 119 | Domain name / search list |
| 42 | NTP (time) servers |

So one exchange configures nearly everything.

### Things that go wrong

- **Rogue DHCP server.** Someone plugs a home router into the office network and
  it starts handing out addresses. Devices get the wrong gateway and lose
  connectivity, or worse, route through the attacker. Enterprise switches prevent
  this with **DHCP snooping** (only trusted ports may send OFFERs).
- **Pool exhaustion.** A conference with 500 devices and a 100-address pool and
  a 7-day lease: after the first day, nobody new can connect. Symptom: new devices
  get APIPA. Fix: bigger pool or shorter lease.
- **Duplicate address.** Two devices, same address (usually one statically
  configured inside the DHCP range). Intermittent, maddening connectivity.

### Practice (25 min)

**A. Watch DORA happen live.** This is very satisfying.

```bash
# Terminal 1 -- start capturing
sudo tcpdump -i any -n -v 'udp port 67 or udp port 68'

# Terminal 2 -- force a renewal
# Linux (NetworkManager):
sudo nmcli device reapply <your-interface>
#   or, more forcefully:
sudo dhclient -r && sudo dhclient
# macOS:
sudo ipconfig set en0 DHCP
```

You'll see the DISCOVER / OFFER / REQUEST / ACK sequence, with all the options
listed. Save a copy in `notes.md`.

**B. Look at your current lease:**

```bash
# Linux
nmcli -f DHCP4 connection show "<connection name>"
cat /var/lib/dhcp/dhclient.leases 2>/dev/null | tail -30

# macOS
ipconfig getpacket en0
```

Find: your address, lease time, gateway, DNS servers.

**C. Log into your router's admin page** (usually `http://192.168.1.1`) and find
the DHCP settings. Look at:
- The address pool range
- The lease time
- The list of currently connected devices

Try setting a **reservation** for one device (bind its MAC to a fixed address) —
this is how you give a printer or server a stable address without configuring it
manually.

### Common confusions

- **"Static IP vs DHCP reservation?"** A static IP is configured *on the device*
  (and the DHCP server doesn't know about it — risk of conflict). A reservation is
  configured *on the server* (device still uses DHCP, always gets the same
  address). **Reservations are almost always the better choice.**
- **"Why is DHCP UDP and not TCP?"** You can't do a TCP handshake without an
  address. UDP broadcast is the only option available.
- **"Does DHCP work across routers?"** Broadcasts don't cross routers, so a
  **DHCP relay** on the router forwards them to a central server. That's how one
  server can serve a whole company.

### Check yourself

1. What are the four DHCP messages, in order?
2. Why are the first messages broadcast?
3. What does a `169.254.x.x` address tell you?
4. Name three things besides an IP address that DHCP provides.

### Further reading

- **Video (~10 min):** "DHCP Explained" — Practical Networking or PowerCert.
- **Article:** "What is DHCP?" — Cloudflare Learning Center.
- **RFC 2131** — the DHCP spec. Long, but the state diagram is instructive.
- **Hands-on:** set up `dnsmasq` on a spare machine or Raspberry Pi and run your
  own DHCP + DNS server. Half a day, enormously clarifying.

---

## Chapter 18 — DNS: turning names into addresses

### In one sentence

DNS is a giant, distributed, cached phone book that turns `example.com` into
`93.184.216.34` — and it is involved in more outages than almost anything else.

### The problem

You remember `example.com`. The network needs `93.184.216.34`. Something must
translate — for hundreds of millions of names, changing constantly, for billions
of queries per second, with no single machine holding the list.

### The idea, in plain language

**Delegate.** Nobody knows everything; everybody knows who to ask next.

```
   You:  "Who is www.example.com?"

   ROOT SERVER:        "No idea. But I know who runs .com -- ask them."
   .COM SERVER:        "No idea. But I know who runs example.com -- ask them."
   EXAMPLE.COM SERVER: "That's me. www.example.com is 93.184.216.34."
```

Three questions, each to a more specific authority. Exactly like finding a person
by asking the country, then the city, then the street.

### The name hierarchy

Read a domain name **right to left**:

```
       www . example . com .
        |      |        |   |
        |      |        |   +-- the ROOT (usually invisible; the trailing dot)
        |      |        +------ TOP-LEVEL DOMAIN (TLD): com, org, uk, io
        |      +--------------- the DOMAIN somebody registered
        +---------------------- a HOST or subdomain within it
```

### The four players

| Role | What it does | Example |
|---|---|---|
| **Stub resolver** | The tiny client inside your OS. Asks one server, caches a little. | Your laptop |
| **Recursive resolver** | Does all the legwork (root → TLD → authoritative), caches everything | Your router, your ISP, `1.1.1.1`, `8.8.8.8` |
| **Root servers** | Know who runs each TLD. 13 logical addresses, hundreds of physical machines. | `a.root-servers.net` |
| **Authoritative server** | Holds the real answers for a specific domain | `ns1.example.com` |

**Your device almost never talks to root servers.** It asks one recursive
resolver, which does everything else.

### Caching is what makes it work

Every answer comes with a **TTL** (time to live) — how long you may remember it.

```
   First lookup of example.com:     3 or 4 round trips, maybe 100+ ms
   Next lookup within the TTL:      0 round trips -- served from cache
```

Caching happens at every level: your browser, your OS, your router, your ISP's
resolver. This is why the internet doesn't collapse.

**The operational consequence you must know:** if you change a DNS record, the old
answer stays cached everywhere until its TTL expires. So **before** a migration,
lower the TTL (to e.g. 60 seconds) and wait for the *old* TTL to pass. Then make
the change. Then raise it again.

### The record types you'll actually meet

| Type | Meaning | Example |
|---|---|---|
| **A** | name → IPv4 address | `example.com. A 93.184.216.34` |
| **AAAA** | name → IPv6 address | `example.com. AAAA 2606:2800:...` |
| **CNAME** | this name is an alias for that name | `www.example.com. CNAME example.com.` |
| **MX** | mail servers for this domain | `example.com. MX 10 mail.example.com.` |
| **NS** | which servers are authoritative | `example.com. NS ns1.example.com.` |
| **TXT** | arbitrary text — used for SPF, DKIM, domain verification | |
| **PTR** | address → name (reverse lookup) | |
| **SOA** | admin info for the zone, including the negative-cache TTL | |

### Practice (35 min) — this is the most useful tool chapter in Part 4

**A. Basic lookups:**

```bash
dig example.com                     # full output -- read the ANSWER section
dig example.com +short              # just the answer
dig example.com A +short
dig example.com AAAA +short
dig example.com MX +short
dig example.com NS +short
dig example.com TXT +short
```

**B. Watch the whole delegation chain.** This is the single best DNS learning
command:

```bash
dig +trace example.com
```

You will literally see it start at the root servers, get referred to `.com`, then
to `example.com`'s own servers, then get the answer. **Run this now.** It makes
the whole chapter concrete.

**C. See caching in action:**

```bash
dig example.com | grep -E "Query time|ANSWER SECTION" -A2
# run it again immediately
dig example.com | grep "Query time"
```

The second query time should be dramatically lower — and watch the TTL count
*down* between runs.

**D. Ask a specific server** (bypasses your local cache):

```bash
dig @1.1.1.1 example.com +short          # Cloudflare
dig @8.8.8.8 example.com +short          # Google
dig @ns1.example.com example.com +short  # straight to the source (authoritative)
```

If different servers give different answers, you've found a propagation or
split-horizon issue.

**E. Reverse lookup:**

```bash
dig -x 93.184.216.34 +short
```

**F. Watch DNS on the wire:**

```bash
sudo tcpdump -i any -n 'udp port 53'
# in another terminal:
dig example.com
```

You'll see the query go out and the response come back — usually one packet each.

**G. See your own resolver settings:**

```bash
cat /etc/resolv.conf              # Linux
scutil --dns | head -20           # macOS
resolvectl status                 # Linux with systemd-resolved
```

### DNS over UDP, TCP, and encrypted

- Normally **UDP port 53** — one small packet each way, fast.
- If the answer is too big, the server sets a "truncated" flag and the client
  **retries over TCP port 53**. **You must allow TCP/53** through firewalls; not
  doing so causes bizarre intermittent failures.
- **DoT** (DNS over TLS, port 853) and **DoH** (DNS over HTTPS, port 443) encrypt
  your queries. DoH is hard to distinguish from normal web traffic — good for
  privacy, awkward for network administrators who rely on internal DNS.

### The classic DNS problems

| Symptom | Likely cause |
|---|---|
| "It works on my machine but not theirs" | DNS cached differently; different resolver |
| Site works by IP but not by name | DNS failure — test with `dig` |
| Change made but not visible | TTL hasn't expired; check with `dig` against multiple resolvers |
| Intermittent slow page loads | One of the configured DNS servers is dead; each lookup waits ~5 s before trying the next |
| Works externally, fails internally (or vice versa) | Split-horizon DNS, or a VPN not pushing internal DNS |
| "DNS_PROBE_FINISHED_NXDOMAIN" | The name genuinely doesn't exist (or a typo) |

**A useful habit:** when anything network-related breaks, test DNS *first* with
`dig`, and test raw connectivity separately with an IP address. That immediately
splits the problem in half.

### Build it in Go (60 min) — a DNS client from raw bytes

`dig` is a wrapper around a surprisingly small protocol. This program builds a
DNS query **byte by byte**, sends it in one UDP datagram, and decodes the
reply, including the "message compression" pointers that confuse everyone the
first time. It uses no DNS library: only `encoding/binary` and a UDP socket.

The wire format you're implementing (RFC 1035):

```
 HEADER (12 bytes)   ID | flags (QR, RD, RA, TC, rcode) | QDCOUNT | ANCOUNT | NSCOUNT | ARCOUNT
 QUESTION            name as labels: 7 e x a m p l e 3 c o m 0 | QTYPE (A=1) | QCLASS (IN=1)
 ANSWER records      name (often a 2-byte POINTER: 0xC0 0x0C) | TYPE | CLASS | TTL | RDLENGTH | RDATA
```

Save as `dnsquery/main.go`:

```go
// dnsquery: build a DNS query by hand, send it over UDP, and decode the answer.
// No DNS library -- just the bytes RFC 1035 describes.
//
//	go run ./dnsquery example.com
//	go run ./dnsquery -type AAAA -server 8.8.8.8 example.com
//	go run ./dnsquery -type CNAME www.github.com
package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"strings"
	"time"
)

var qtypes = map[string]uint16{"A": 1, "NS": 2, "CNAME": 5, "MX": 15, "TXT": 16, "AAAA": 28}

func main() {
	server := flag.String("server", "1.1.1.1", "recursive resolver to ask")
	qtype := flag.String("type", "A", "record type: A, AAAA, CNAME, NS, MX, TXT")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("usage: dnsquery [-server IP] [-type T] name")
		os.Exit(2)
	}
	t, ok := qtypes[strings.ToUpper(*qtype)]
	if !ok {
		fmt.Println("unsupported type", *qtype)
		os.Exit(2)
	}

	id := uint16(rand.N(1 << 16))
	query := buildQuery(id, flag.Arg(0), t)
	fmt.Printf("query: %d bytes  % x ...\n", len(query), query[:12])

	conn, err := net.Dial("udp", net.JoinHostPort(*server, "53"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second)) // UDP has no retries: we own the timeout

	start := time.Now()
	if _, err := conn.Write(query); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := make([]byte, 1232) // the EDNS-era "safe" UDP DNS size
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("no answer:", err)
		os.Exit(1)
	}
	fmt.Printf("reply: %d bytes from %s in %v\n\n", n, *server, time.Since(start).Round(time.Microsecond))
	if err := printReply(buf[:n], id); err != nil {
		fmt.Println("bad reply:", err)
		os.Exit(1)
	}
}

// buildQuery writes the 12-byte header plus one question.
func buildQuery(id uint16, name string, qtype uint16) []byte {
	b := make([]byte, 12, 512)
	binary.BigEndian.PutUint16(b[0:], id)
	binary.BigEndian.PutUint16(b[2:], 0x0100) // flags: RD=1 ("recursion desired")
	binary.BigEndian.PutUint16(b[4:], 1)      // QDCOUNT = 1 question
	// a name is a sequence of length-prefixed labels: 7example3com0
	for _, label := range strings.Split(strings.TrimSuffix(name, "."), ".") {
		b = append(b, byte(len(label)))
		b = append(b, label...)
	}
	b = append(b, 0)
	b = binary.BigEndian.AppendUint16(b, qtype)
	b = binary.BigEndian.AppendUint16(b, 1) // class IN (internet)
	return b
}

func printReply(msg []byte, wantID uint16) error {
	if len(msg) < 12 {
		return errors.New("short header")
	}
	if id := binary.BigEndian.Uint16(msg[0:]); id != wantID {
		return fmt.Errorf("ID mismatch: sent %d got %d (spoofed or stale reply?)", wantID, id)
	}
	flags := binary.BigEndian.Uint16(msg[2:])
	rcode := flags & 0xF
	qd, an := binary.BigEndian.Uint16(msg[4:]), binary.BigEndian.Uint16(msg[6:])
	fmt.Printf("flags: QR=%d AA=%d TC=%d RD=%d RA=%d  rcode=%s  answers=%d\n",
		flags>>15&1, flags>>10&1, flags>>9&1, flags>>8&1, flags>>7&1, rcodeName(rcode), an)
	if flags>>9&1 == 1 {
		fmt.Println("TC=1: answer truncated -- a real resolver now retries over TCP")
	}

	off := 12
	for i := 0; i < int(qd); i++ { // skip the echoed question
		_, next, err := readName(msg, off)
		if err != nil {
			return err
		}
		off = next + 4
	}
	for i := 0; i < int(an); i++ {
		name, next, err := readName(msg, off)
		if err != nil {
			return err
		}
		if next+10 > len(msg) {
			return errors.New("truncated record header")
		}
		typ := binary.BigEndian.Uint16(msg[next:])
		ttl := binary.BigEndian.Uint32(msg[next+4:])
		rdlen := int(binary.BigEndian.Uint16(msg[next+8:]))
		rdata := next + 10
		if rdata+rdlen > len(msg) {
			return errors.New("truncated rdata")
		}
		fmt.Printf("%-30s TTL=%-6d %-5s %s\n", name, ttl, typeName(typ), rdataString(msg, typ, rdata, rdlen))
		off = rdata + rdlen
	}
	return nil
}

// readName decodes a possibly-compressed name. A byte >= 0xC0 starts a
// 2-byte pointer to an earlier offset: that is DNS "message compression".
func readName(msg []byte, off int) (string, int, error) {
	var labels []string
	next := -1 // where parsing continues after the name (set at first pointer)
	for hops := 0; hops < 20; hops++ {
		if off >= len(msg) {
			return "", 0, errors.New("name runs past end")
		}
		l := int(msg[off])
		switch {
		case l == 0:
			if next < 0 {
				next = off + 1
			}
			return strings.Join(labels, ".") + ".", next, nil
		case l&0xC0 == 0xC0:
			if off+1 >= len(msg) {
				return "", 0, errors.New("bad pointer")
			}
			if next < 0 {
				next = off + 2
			}
			off = int(binary.BigEndian.Uint16(msg[off:]) & 0x3FFF)
		default:
			if off+1+l > len(msg) {
				return "", 0, errors.New("label runs past end")
			}
			labels = append(labels, string(msg[off+1:off+1+l]))
			off += 1 + l
		}
	}
	return "", 0, errors.New("compression loop")
}

func rdataString(msg []byte, typ uint16, off, n int) string {
	switch typ {
	case 1, 28: // A, AAAA
		return net.IP(msg[off : off+n]).String()
	case 2, 5: // NS, CNAME
		name, _, err := readName(msg, off)
		if err != nil {
			return "?"
		}
		return name
	case 15: // MX: preference + name
		name, _, _ := readName(msg, off+2)
		return fmt.Sprintf("%d %s", binary.BigEndian.Uint16(msg[off:]), name)
	case 16: // TXT: one or more length-prefixed strings
		var parts []string
		for i := off; i < off+n; {
			l := int(msg[i])
			if i+1+l > off+n {
				break
			}
			parts = append(parts, fmt.Sprintf("%q", msg[i+1:i+1+l]))
			i += 1 + l
		}
		return strings.Join(parts, " ")
	}
	return fmt.Sprintf("% x", msg[off:off+n])
}

func typeName(t uint16) string {
	for k, v := range qtypes {
		if v == t {
			return k
		}
	}
	return fmt.Sprint(t)
}

func rcodeName(r uint16) string {
	switch r {
	case 0:
		return "NOERROR"
	case 2:
		return "SERVFAIL"
	case 3:
		return "NXDOMAIN"
	case 5:
		return "REFUSED"
	}
	return fmt.Sprint(r)
}
```

Real output:

```text
$ go run ./dnsquery example.com
query: 29 bytes  49 8d 01 00 00 01 00 00 00 00 00 00 ...
reply: 61 bytes from 1.1.1.1 in 17.618ms

flags: QR=1 AA=0 TC=0 RD=1 RA=1  rcode=NOERROR  answers=2
example.com.                   TTL=179    A     172.66.147.243
example.com.                   TTL=179    A     104.20.23.154

$ go run ./dnsquery www.github.com
www.github.com.                TTL=3040   CNAME github.com.
github.com.                    TTL=8      A     20.207.73.82

$ go run ./dnsquery nonexistent-zzzz.example
flags: QR=1 AA=0 TC=0 RD=1 RA=1  rcode=NXDOMAIN  answers=0
```

**What to notice:**

- **The whole query is 29 bytes.** DNS is fast because it's one small
  datagram each way, with no handshake.
- **`AA=0`**: the answer came from the resolver's *cache*, not from the
  authoritative server. Run it twice and watch the TTL count down. That's
  the cache ageing in real time.
- **The ID check matters for security.** A reply with the wrong ID is
  rejected. Off-path DNS spoofing attacks (Kaminsky, 2008) work by guessing
  this 16-bit ID, which is why resolvers also randomise the source port.
- **The CNAME chain is in one reply.** The resolver followed `www.github.com
  → github.com` for you and returned both records.
- **`TC=1`** means the answer didn't fit in the UDP packet; real clients
  retry over TCP port 53. Our 1232-byte buffer is the size the DNS community
  agreed in 2020 avoids IP fragmentation on almost every path.

**Exercises:**

1. Point it at a server that doesn't exist (`-server 192.0.2.1`). You get a
   timeout, not an error, because UDP can't tell "no server" from "slow
   server". This is why a dead resolver in `/etc/resolv.conf` adds a fixed
   5-second delay (the "suspiciously round number" from Chapter 26).
2. Query an **authoritative** server directly: find it with
   `go run ./dnsquery -type NS example.com`, resolve it, then ask it with
   `-server`. Check that `AA=1`.
3. Add **retries with exponential backoff** (1 s, 2 s, 4 s) and query
   two servers in parallel, using whichever answers first. That is what your
   OS stub resolver does.
4. Stretch: implement the **iterative** algorithm. Start at a root server
   (`198.41.0.4`), follow NS referrals in the *authority* section, and
   reproduce `dig +trace` yourself.

**Across the series:** the [HTTPS guide's DNS chapter](../../v2-https/real-life-guide-v1.md#chapter-2-dns-finding-the-address)
follows this lookup from the browser's point of view (browser cache, OS
cache, DoH, HTTPS records). The Go standard library's resolver is
`net.Resolver`; Chapter 55's `dnsdiff` lab uses it to query many servers at
once.

### Common confusions

- **"DNS is slow, that's why my site is slow."** Usually only the *first* lookup.
  Check with the `curl -w` timing from Chapter 4 — `time_namelookup` tells you
  exactly.
- **"`nslookup` vs `dig`?"** `dig` is better: more detail, more scriptable.
  `nslookup` is the Windows default and is fine for quick checks.
- **"Changing DNS servers makes the internet faster."** It can reduce lookup
  latency slightly, but it does not increase bandwidth. The main reasons to change
  are privacy, filtering, and reliability.
- **"A CNAME at the top of my domain doesn't work."** Correct — a CNAME can't
  coexist with other records, and the zone apex needs NS and SOA records.
  Providers offer `ALIAS`/`ANAME` workarounds.

### Check yourself

1. Explain the delegation chain from root to `www.example.com`.
2. What is a TTL and why does it matter before a migration?
3. What's the difference between a recursive resolver and an authoritative server?
4. When does DNS use TCP?
5. What does `dig +trace` show you?

### Further reading

- **Zine (best beginner resource, genuinely):** "How DNS Works" — Julia Evans
  (jvns.ca). She also has an excellent free comic/zine on DNS. Start here.
- **Interactive:** `messwithdns.net` — Julia Evans' free playground where you
  create real DNS records and see what happens. Spend 30 minutes here; it's the
  fastest way to make DNS click.
- **Video (~15 min):** "DNS Explained" — Practical Networking, or Computerphile's
  "How DNS Works".
- **Article:** "What is DNS?" — Cloudflare Learning Center (their whole DNS
  section is excellent).
- **Book:** *DNS and BIND* — Cricket Liu. The definitive reference if you go deep.

---

### End of Part 4 — Milestone check

- [ ] I can name the four DHCP messages and what each does
- [ ] I know what a `169.254.x.x` address means instantly
- [ ] **I have watched DORA in tcpdump**
- [ ] I can explain the DNS delegation chain
- [ ] **I have run `dig +trace` and understood the output**
- [ ] I know why you lower a TTL before a migration

---

# Part 5 — Ports and the transport layer

IP got the packet to the right *machine*. Now: which *program* on that machine,
and did everything arrive?

## Chapter 19 — Ports and sockets: which program gets the data?

### In one sentence

A **port** is a number that says which program on a machine a packet is for, and a
connection is identified by the combination of both machines' addresses and ports.

### The problem

Your laptop has one IP address and is running fifty programs — a browser with
twelve tabs, an email client, a chat app, a music streamer. A packet arrives.
Which program gets it?

### The idea, in plain language

Back to the office building analogy from Chapter 2:

```
   IP ADDRESS  = the building's street address     "12 High Street"
   PORT NUMBER = the room/department number         "Room 443"
```

The postman gets the letter to the building using the street address. The
receptionist gets it to the right room using the room number.

### The numbers

A port is a **16-bit number**: 0 to 65535.

| Range | Name | Used for |
|---|---|---|
| 0 – 1023 | **Well-known** | Standard services. Binding to these usually needs admin rights. |
| 1024 – 49151 | **Registered** | Applications that registered a number (MySQL 3306, Postgres 5432) |
| 49152 – 65535 | **Ephemeral** | Temporary, random ports your machine picks when it *starts* a connection |

The well-known ports worth knowing by heart (the rest are in Appendix D):

```
    22  SSH            53  DNS           80  HTTP        443  HTTPS
    25  SMTP (mail)   123  NTP (time)   3306 MySQL      5432  PostgreSQL
    67/68 DHCP        3389 RDP          6379 Redis      27017 MongoDB
```

### The four-part identity of a connection

This is the key concept. A connection is identified by **four things together**:

```
   (your IP, your port, their IP, their port)

   Example -- two browser tabs to the same website:

     tab 1:  192.168.1.50 : 51000  <-->  93.184.216.34 : 443
     tab 2:  192.168.1.50 : 51001  <-->  93.184.216.34 : 443
                            ^^^^^
                    only this differs -- and that's enough
```

This is called the **4-tuple** (or 5-tuple, including the protocol). It's why:

- **One server on port 443 can handle millions of simultaneous connections** —
  they all differ in the client's address and port.
- Your machine can have many connections to the same server at once.
- The receiving machine can always tell them apart.

The **server** listens on a fixed, well-known port (443). The **client** uses a
random ephemeral port so its replies can be sorted out.

### Listening vs connected

```bash
ss -tuln      # Linux: listening TCP+UDP sockets
netstat -an | grep LISTEN    # macOS
```

```
State    Local Address:Port     Peer Address:Port
LISTEN   0.0.0.0:22            0.0.0.0:*        <- SSH, waiting for anyone
LISTEN   127.0.0.1:5432        0.0.0.0:*        <- Postgres, LOCAL ONLY
ESTAB    192.168.1.50:51000    93.184.216.34:443  <- an actual connection
```

Note the difference between `0.0.0.0:22` and `127.0.0.1:5432`:

- **`0.0.0.0`** means "listen on **all** my addresses" — reachable from the
  network.
- **`127.0.0.1`** means "listen on loopback **only**" — reachable **only from
  this machine**.

> **This distinction causes an enormous number of "why can't I connect to my
> server?" problems.** Your database is running, the port is open, but it's bound
> to `127.0.0.1` so nothing outside the machine can reach it. Check this first,
> always.

### Practice (25 min)

**A. See what's listening on your machine:**

```bash
# Linux
ss -tulnp

# macOS
sudo lsof -iTCP -sTCP:LISTEN -n -P
netstat -an | grep LISTEN
```

For each one, ask: *do I know what that is, and should it be reachable from the
network?* This is a genuinely useful security habit.

**B. See your active connections:**

```bash
ss -tn state established        # Linux
netstat -an | grep ESTABLISHED  # macOS
```

Open a few browser tabs and run it again. Watch the count grow.

**C. Be a server and a client yourself** — this makes ports concrete:

```bash
# Terminal 1: listen on port 9999
nc -l 9999

# Terminal 2: connect to it and type something
nc 127.0.0.1 9999
# whatever you type appears in terminal 1. You just built a chat app.

# Terminal 3: watch the connection exist
ss -tn '( sport = :9999 or dport = :9999 )'
```

**D. Prove the binding-address rule:**

```bash
# Bind to loopback only
python3 -m http.server 8000 --bind 127.0.0.1
# From this machine:      curl http://127.0.0.1:8000   -> works
# From another device:    curl http://<your-ip>:8000    -> FAILS

# Now bind to everything
python3 -m http.server 8000 --bind 0.0.0.0
# From another device: works.
```

**E. Check whether a remote port is open:**

```bash
nc -zv example.com 443        # "succeeded" or "refused"/timeout
nc -zv example.com 444        # almost certainly nothing there
```

### Build it in Go (25 min) — the binding rules, tried for real

Three rules from this chapter cause real outages: **one listener per
address:port**, **a loopback bind only accepts loopback traffic**, and the
less well-known **`SO_REUSEPORT`**, which lets several sockets share one port on
purpose.

```go
//go:build linux || darwin

// bindlab: Chapter 19's binding rules, tried for real.
//
//	go run ./bindlab conflict     # two listeners on one port -> EADDRINUSE
//	go run ./bindlab reuseport    # SO_REUSEPORT: four listeners share one port
//	go run ./bindlab loopback     # 127.0.0.1 vs 0.0.0.0, and what each accepts
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: bindlab conflict|reuseport|loopback")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "conflict":
		conflict()
	case "reuseport":
		reusePort()
	case "loopback":
		loopback()
	}
}

func conflict() {
	a, err := net.Listen("tcp", "127.0.0.1:7300")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer a.Close()
	_, err = net.Listen("tcp", "127.0.0.1:7300")
	fmt.Println("second listen:", err)
	fmt.Println("EADDRINUSE?   ", errors.Is(err, syscall.EADDRINUSE))
}

// reusePort opens four sockets on the SAME port with SO_REUSEPORT set before
// bind(). The kernel then spreads incoming connections across them -- the
// trick nginx, Envoy, and HAProxy use to scale accept() across CPU cores.
func reusePort() {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			err := c.Control(func(fd uintptr) {
				opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
			})
			if err != nil {
				return err
			}
			return opErr
		},
	}
	counts := make([]int, 4)
	var mu sync.Mutex
	for i := range counts {
		ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:7301")
		if err != nil {
			fmt.Println("listener", i, err)
			return
		}
		go func() {
			for {
				c, err := ln.Accept()
				if err != nil {
					return
				}
				mu.Lock()
				counts[i]++
				mu.Unlock()
				c.Close()
			}
		}()
	}
	for i := 0; i < 400; i++ {
		if c, err := net.Dial("tcp", "127.0.0.1:7301"); err == nil {
			c.Close()
		}
	}
	time.Sleep(200 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("connections accepted per listener:", counts)
	fmt.Println("(Linux hashes each connection to one socket -> roughly even;")
	fmt.Println(" macOS/BSD semantics differ -> often all on one listener)")
}

func loopback() {
	lo, _ := net.Listen("tcp", "127.0.0.1:7302")
	all, _ := net.Listen("tcp", "0.0.0.0:7303")
	defer lo.Close()
	defer all.Close()

	// Find a non-loopback address of this machine to test "from the network".
	c, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		fmt.Println("no network address to test with:", err)
		return
	}
	lan := c.LocalAddr().(*net.UDPAddr).IP.String()
	c.Close()

	for _, t := range []string{"127.0.0.1:7302", lan + ":7302", "127.0.0.1:7303", lan + ":7303"} {
		conn, err := net.DialTimeout("tcp", t, time.Second)
		if err != nil {
			fmt.Printf("%-22s FAIL  %v\n", t, err)
			continue
		}
		conn.Close()
		fmt.Printf("%-22s ok\n", t)
	}
	fmt.Println("\nA socket bound to 127.0.0.1 only accepts traffic addressed to 127.0.0.1.")
	fmt.Println("No firewall rule can change that: fix the bind address, not the firewall.")
}
```

```text
$ go run ./bindlab conflict
second listen: listen tcp 127.0.0.1:7300: bind: address already in use
EADDRINUSE?    true

$ go run ./bindlab loopback
127.0.0.1:7302         ok
192.168.1.2:7302       FAIL  dial tcp 192.168.1.2:7302: connect: connection refused
127.0.0.1:7303         ok
192.168.1.2:7303       ok

$ go run ./bindlab reuseport          # on Linux
connections accepted per listener: [103 105 104 88]
$ go run ./bindlab reuseport          # on macOS
connections accepted per listener: [0 0 0 400]
```

**What to notice:**

- **"Refused" from your own LAN address.** The 127.0.0.1 listener exists, but
  the connection was addressed to 192.168.1.2:7302, where nothing listens.
  This is the most common "works on my machine" bug in container and VM
  setups: the app binds `127.0.0.1` and everything outside the host gets
  refused. Fix the bind, not the firewall.
- **`ListenConfig.Control`** runs your code on the raw file descriptor
  *after* `socket()` but *before* `bind()`. That's the only moment socket
  options like `SO_REUSEPORT` can be set. This hook is how Go programs
  reach any socket option the `net` package doesn't wrap.
- **The same code behaves differently on Linux and macOS.** Linux load-balances
  new connections across `SO_REUSEPORT` sockets by hashing the 4-tuple.
  BSD-derived systems historically deliver to one socket. Always test
  kernel behaviour on the kernel you deploy to.

**A real incident this explains:** during testing for this guide, a backend
bound to `127.0.0.1:8081` was killed. Instead of getting "connection
refused", the clients started receiving `404`s. Another program (Docker
Desktop) was listening on the *wildcard* `*:8081`. The kernel gives each
connection to the **most specific** matching listener. While the backend was
alive, it was more specific. Once it died, the wildcard listener got the
traffic. If you see "wrong service answered", check `ss -ltnp` /
`lsof -nP -iTCP -sTCP:LISTEN` for overlapping binds.

**Exercises:**

1. Start `python3 -m http.server 8000 --bind 0.0.0.0`, then try to
   `net.Listen("tcp", "127.0.0.1:8000")` in Go. Does it conflict? (On Linux it
   does: a wildcard bind covers every address.)
2. Run `reuseport` under `strace -f -e trace=setsockopt,bind,listen` and
   find the four `setsockopt(... SO_REUSEPORT ...)` calls.
3. Read how [Linux's Chapter 37 on network views](../../os-linux/real-life-os-guide.md#chapter-37-how-a-linux-box-sees-the-network) shows
   these listeners with `ss -ltnp`, then find which process owns each.

### Common confusions

- **"Is a port a physical thing?"** No. It's purely a number in a packet header.
  Nothing plugs into it.
- **"Port 80 vs port 8080?"** Both are just numbers. 80 is the *convention* for
  HTTP; 8080 is a common alternative that doesn't need admin rights.
- **"Open port vs listening port?"** "Listening" = a program is waiting on it.
  "Open" usually means a firewall permits traffic to it. You need both.
- **"Connection refused vs timeout?"** *Refused* = the machine is there and
  actively said no (nothing listening). *Timeout* = nothing came back at all
  (firewall dropping silently, or host unreachable). **Very different diagnoses.**

### Check yourself

1. Why does a connection need four numbers to identify it, not two?
2. What's the difference between binding to `0.0.0.0` and `127.0.0.1`?
3. If a server listens on 443 and a million clients connect, how does it tell
   them apart?
4. What's the difference between "connection refused" and "connection timed out"?

### Further reading

- **Video:** "Ports and Sockets Explained" — Practical Networking.
- **Article:** "What is a computer port?" — Cloudflare Learning Center.
- **Reference:** the IANA port registry (search "IANA service name port number
  registry") — the official list.
- **Blog:** Julia Evans' zine "How to be a Wizard Programmer" / her networking
  tools posts cover `ss`, `netstat`, and `lsof` well.

---

## Chapter 20 — UDP: fire and forget

### In one sentence

UDP adds ports and a checksum to IP, and **nothing else** — no connection, no
reliability, no ordering — which makes it fast and perfect for some jobs.

### The entire UDP header

That's not a joke. This is all of it — 8 bytes:

```
+------------------+------------------+
| Source Port (2B) | Dest Port (2B)   |
+------------------+------------------+
| Length (2B)      | Checksum (2B)    |
+------------------+------------------+
|          your data...               |
+-------------------------------------+
```

Compare with TCP's 20+ bytes and you see the philosophy: **UDP gets out of the
way.**

### What UDP does and doesn't do

| UDP does | UDP does NOT do |
|---|---|
| Say which program it's for (ports) | Establish a connection |
| Check the data isn't corrupted | Guarantee delivery |
| Preserve message boundaries | Guarantee order |
| Deliver *immediately* | Retransmit lost data |
| | Control its sending rate |

**Send one `sendto()` of 200 bytes → exactly one 200-byte datagram arrives (or
doesn't).** No merging, no splitting. That's different from TCP and sometimes
exactly what you want.

### When "unreliable" is the right answer

This surprises beginners: why would you ever choose a protocol that can lose your
data?

| Use case | Why UDP wins |
|---|---|
| **DNS lookups** (Ch 18) | One tiny question, one tiny answer. Setting up a connection would triple the cost. Lost? Just ask again. |
| **Voice and video calls** | A voice packet that arrives 200 ms late is **useless** — you'd rather skip it than pause the call waiting for it. Retransmission actively hurts. |
| **Online games** | Only the *latest* position matters. Re-sending where the player was 300 ms ago is pointless. |
| **DHCP** (Ch 17) | You don't have an address yet — you can't do a TCP handshake. |
| **NTP** (time sync) | Tiny, frequent, and lateness would defeat the purpose. |
| **Video streaming over QUIC / HTTP-3** | Wants reliability but implemented *its own way*, in the application (Ch 30). |
| **VPNs** (WireGuard, OpenVPN-UDP) | Avoids "TCP inside TCP", where two retransmission systems fight each other and collapse. |

**The pattern:** UDP is right when *lateness is worse than loss*, or when the
message is so small that a connection isn't worth it.

### What your app must do itself

If you use UDP and you *do* need reliability, you must build it:
sequence numbers, acknowledgements, retransmission, ordering, and — importantly —
**rate control**. An application that blasts UDP as fast as it can is a bad
network citizen and will cause congestion for everyone.

### Practice (20 min)

**A. Send UDP by hand:**

```bash
# Terminal 1: listen
nc -u -l 9999

# Terminal 2: send
nc -u 127.0.0.1 9999
# type lines; they appear in terminal 1
```

Notice: unlike the TCP version in Chapter 19, **there's no connection**. Kill
terminal 1 and keep typing in terminal 2 — it doesn't complain. The data just goes
nowhere. *That's UDP.*

**B. See UDP in the wild.** Your machine is doing this constantly:

```bash
sudo tcpdump -i any -n udp
```

Leave it a minute. You'll see DNS (port 53), mDNS (5353, `.local` discovery),
maybe NTP (123), DHCP (67/68), and QUIC (443/UDP — that's HTTP/3, Chapter 30).

**C. Compare TCP and UDP DNS:**

```bash
sudo tcpdump -i any -n 'port 53' &
dig example.com                # uses UDP -- 2 packets
dig +tcp example.com           # forces TCP -- handshake + query + teardown
```

Count the packets in each case. **That's the cost of a connection**, visible.

**D. Check UDP loss statistics on your machine:**

```bash
netstat -su | head -20      # look for "packet receive errors"
nstat -az | grep -i udp     # Linux
```

Because UDP loss is silent, these counters are the only way to see it.

### Build it in Go (40 min) — measure what UDP doesn't promise

The "What your app must do itself" list above is easier to remember once
you've had to implement it. This lab is an echo server plus a client that
stamps every datagram with a sequence number, then measures **loss,
reordering, duplication, RTT percentiles, and jitter**. These are the exact
metrics a video-call or game client tracks.

```go
// udpping: a UDP echo server and a client that measures what UDP does NOT
// promise you: loss, reordering, duplication, and jitter.
//
//	go run ./udpping -serve :9999
//	go run ./udpping -to 127.0.0.1:9999 -n 200 -every 10ms
package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"slices"
	"syscall"
	"time"
)

func main() {
	serve := flag.String("serve", "", "run the echo server on this address")
	to := flag.String("to", "", "server address to probe")
	n := flag.Int("n", 100, "number of probes")
	every := flag.Duration("every", 20*time.Millisecond, "interval between probes")
	size := flag.Int("size", 64, "datagram size in bytes (try 1472 and 1473 with a 1500 MTU)")
	flag.Parse()
	switch {
	case *serve != "":
		server(*serve)
	case *to != "":
		client(*to, *n, *every, *size)
	default:
		flag.Usage()
		os.Exit(2)
	}
}

func server(addr string) {
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("echoing UDP on", pc.LocalAddr())
	buf := make([]byte, 65535)
	for {
		n, from, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		pc.WriteTo(buf[:n], from) // no connection, no state: reply to whoever sent it
	}
}

func client(addr string, n int, every time.Duration, size int) {
	conn, err := net.Dial("udp", addr) // "connected" UDP: the kernel filters replies for us
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	if size < 16 {
		size = 16
	}

	sent := make(map[uint64]time.Time, n)
	var rtts []time.Duration
	var lastSeq uint64
	received, reordered, dups := 0, 0, 0
	seen := make(map[uint64]bool, n)

	// Receiver goroutine: every reply carries the sequence number we sent.
	type reply struct {
		seq uint64
		at  time.Time
	}
	replies := make(chan reply, n)
	go func() {
		buf := make([]byte, 65535)
		warned := false
		for {
			m, err := conn.Read(buf)
			if errors.Is(err, syscall.ECONNREFUSED) {
				// UDP has no handshake, so how can it be "refused"? The server's
				// kernel answered with ICMP "port unreachable" and ours relayed it.
				if !warned {
					fmt.Println("read: connection refused -> ICMP port unreachable: nothing listens there")
					warned = true
				}
				continue
			}
			if err != nil {
				return // socket closed
			}
			if m >= 8 {
				replies <- reply{binary.BigEndian.Uint64(buf), time.Now()}
			}
		}
	}()

	payload := make([]byte, size)
	tick := time.NewTicker(every)
	defer tick.Stop()
	for seq := uint64(1); seq <= uint64(n); seq++ {
		<-tick.C
		binary.BigEndian.PutUint64(payload, seq)
		sent[seq] = time.Now()
		if _, err := conn.Write(payload); err != nil {
			fmt.Printf("seq %d: write error: %v\n", seq, err) // e.g. "message too long" past the MTU
		}
	drain:
		for {
			select {
			case r := <-replies:
				switch {
				case seen[r.seq]:
					dups++
				case r.seq < lastSeq:
					reordered++
					fallthrough
				default:
					seen[r.seq] = true
					received++
					rtts = append(rtts, r.at.Sub(sent[r.seq]))
					lastSeq = max(lastSeq, r.seq)
				}
			default:
				break drain
			}
		}
	}
	time.Sleep(time.Second) // let stragglers arrive
	for len(replies) > 0 {
		r := <-replies
		if !seen[r.seq] {
			seen[r.seq] = true
			received++
			rtts = append(rtts, r.at.Sub(sent[r.seq]))
		}
	}

	loss := 100 * float64(n-received) / float64(n)
	fmt.Printf("%d sent, %d received, %.1f%% loss, %d reordered, %d duplicated\n",
		n, received, loss, reordered, dups)
	if len(rtts) == 0 {
		return
	}
	slices.Sort(rtts)
	var sum, jitter float64
	for i, r := range rtts {
		sum += float64(r)
		if i > 0 {
			jitter += math.Abs(float64(r - rtts[i-1]))
		}
	}
	fmt.Printf("rtt min/p50/p99/max = %v / %v / %v / %v   avg=%v\n",
		rtts[0], rtts[len(rtts)/2], rtts[len(rtts)*99/100], rtts[len(rtts)-1],
		time.Duration(sum/float64(len(rtts))))
}
```

```text
$ go run ./udpping -serve 127.0.0.1:9999 &
$ go run ./udpping -to 127.0.0.1:9999 -n 100 -every 2ms
100 sent, 100 received, 0.0% loss, 0 reordered, 0 duplicated
rtt min/p50/p99/max = 52.792µs / 132.083µs / 1.360292ms / 1.360292ms   avg=144.76µs

$ go run ./udpping -to 127.0.0.1:9998 -n 3      # nothing listening
read: connection refused -> ICMP port unreachable: nothing listens there
3 sent, 0 received, 100.0% loss, 0 reordered, 0 duplicated

$ go run ./udpping -to 127.0.0.1:9999 -n 3 -size 70000
seq 1: write error: ... write: message too long
```

**What to notice:**

- **"Connection refused" on a protocol without connections.** The server's
  kernel answered with an ICMP *port unreachable* (Chapter 14). Because the
  client socket is "connected", the kernel hands that error back to the
  next `Read`. An unconnected socket (`ListenPacket`) would never see it.
- **`message too long`**: a UDP datagram can't exceed 65,507 bytes of
  payload, and anything over the path MTU is fragmented at the IP layer
  (Chapter 15). Real UDP protocols (DNS, QUIC) keep datagrams around 1,200
  bytes to avoid fragmentation completely.
- **One goroutine sends, another receives.** UDP gives you no
  request/response pairing. The sequence number *is* your pairing, and the
  same is true of QUIC packet numbers and RTP in voice calls.

**Exercises:**

1. **(Linux)** Add impairment and measure it:
   `sudo tc qdisc add dev lo root netem delay 20ms 5ms loss 2% reorder 10%`,
   run the client against `127.0.0.1`, then remove it with
   `sudo tc qdisc del dev lo root`. Do your numbers match what you
   configured?
2. Find the largest `-size` that works between two real machines on your LAN.
   It should be 1472 bytes (1500 MTU − 20 IP − 8 UDP) unless you add a
   don't-fragment option.
3. Add a **retransmit**: if a reply doesn't arrive within 3× the median
   RTT, resend it. Congratulations, you've started reinventing TCP. Stop
   and read Chapter 22.

### Common confusions

- **"UDP is unreliable so it's worse."** It's *simpler*. For the right job it is
  strictly better. Choosing TCP for a voice call would make the call worse.
- **"UDP is faster."** Per-packet, yes — no handshake, no waiting for
  acknowledgements. But UDP has no congestion control, so on a congested network,
  a well-behaved TCP flow may actually deliver more usable data.
- **"Nobody uses UDP for important things."** HTTP/3 — a growing share of all web
  traffic — runs over UDP.

### Check yourself

1. What are the only two things UDP adds to IP?
2. Why is UDP the right choice for a voice call?
3. Why does DHCP have to use UDP?
4. What must your application implement if it uses UDP but needs reliability?

### Further reading

- **Video:** "TCP vs UDP" — Practical Networking or PowerCert. Almost every
  networking channel has one; they're all fine.
- **Article:** "What is UDP?" — Cloudflare Learning Center.
- **Book:** Kurose & Ross, Chapter 3 (Transport Layer) — the standard treatment.

---

## Chapter 21 — TCP part 1: the three-way handshake

### In one sentence

Before sending any data, TCP performs a three-packet greeting that establishes
both sides' starting positions and confirms both directions work.

### The problem

UDP just throws data. TCP promises reliable, ordered delivery. To make that
promise, both sides need to agree on where the conversation starts and confirm
they can actually hear each other.

### The idea, in plain language

A phone call:

```
   YOU:   "Hello?"                         <-- can you hear me?
   THEM:  "Hello, yes I can hear you."     <-- I hear you; can you hear me?
   YOU:   "Great, I hear you too."         <-- confirmed. Now let's talk.
```

Three messages. After that, both sides *know* the other can hear them. Two
wouldn't be enough (the first speaker would never know if their "hello" was
heard).

### How it actually works

```
   CLIENT                                              SERVER
   (your laptop)                                    (web server, LISTENing)
        |                                                  |
        |  1. SYN                                          |
        |     "Let's talk. My byte numbering starts at X."  |
        |------------------------------------------------->|
        |                                                  |
        |  2. SYN + ACK                                    |
        |     "OK. My numbering starts at Y.                |
        |      And I acknowledge your X."                   |
        |<-------------------------------------------------|
        |                                                  |
        |  3. ACK                                          |
        |     "I acknowledge your Y."                       |
        |------------------------------------------------->|
        |                                                  |
        |============ CONNECTION ESTABLISHED ==============|
        |         (data can now flow both ways)            |
```

- **SYN** = "synchronise" — here's my starting sequence number.
- **ACK** = "acknowledge" — I received up to here.
- The middle packet does both jobs at once, which is why it's three packets and
  not four.

**Cost: one round trip** before any data moves. That's the first RTT in your
Chapter 4 latency budget.

### What's in the handshake besides hello

The SYN packets also negotiate options both sides will use:

| Option | What it agrees |
|---|---|
| **MSS** | The biggest chunk of data each will accept (Chapter 15) |
| **Window scale** | Lets the window be bigger than 64 KB (Chapter 24) — essential for speed |
| **SACK permitted** | Enables smarter loss recovery (Chapter 22) |
| **Timestamps** | Better round-trip measurement |

**If a middlebox strips these options, the connection still works but is much
slower.** Worth remembering when a link is mysteriously slow.

### The three outcomes of trying to connect

This is diagnostic gold — learn these three:

```
   1. SYN  ->  SYN-ACK  ->  ACK           = CONNECTED
      Something is listening and let you in.

   2. SYN  ->  RST                         = CONNECTION REFUSED
      The machine is there and reachable, but NOTHING IS LISTENING
      on that port. (RST = "reset" = go away.)
      -> the service is down, or you have the wrong port

   3. SYN  ->  (silence)  -> SYN retry -> (silence) ... = TIMED OUT
      Nothing came back at all.
      -> a firewall is silently dropping, or the host is unreachable,
         or the route is broken
```

**Refused ≠ timeout.** Refused means you *reached* the machine. Timeout means you
possibly didn't. This distinction saves hours.

### Practice (30 min) — capture your first handshake

This is the most important practice exercise in Part 5.

```bash
# Terminal 1: capture ONLY handshake packets to a site
sudo tcpdump -i any -n 'tcp port 443 and host example.com'

# Terminal 2:
curl -s https://example.com > /dev/null
```

You will see, right at the top:

```
IP 192.168.1.50.51000 > 93.184.216.34.443: Flags [S],  seq 1234567, win 64240,
      options [mss 1460,sackOK,TS val ...,nop,wscale 7], length 0
IP 93.184.216.34.443 > 192.168.1.50.51000: Flags [S.], seq 7654321, ack 1234568,
      win 65535, options [mss 1460,sackOK,...], length 0
IP 192.168.1.50.51000 > 93.184.216.34.443: Flags [.],  ack 7654322, win 502, length 0
```

**Read it:**
- `[S]` = SYN. `[S.]` = SYN+ACK (the `.` means ACK). `[.]` = ACK only.
- `seq 1234567` = "my numbering starts here"
- `ack 1234568` = "I got your 1234567, send me 1234568 next" (note: +1)
- `options [mss 1460, sackOK, wscale 7]` = the negotiation described above

**Congratulations — you have now seen the foundation of the internet with your own
eyes.** Save this capture in `notes.md`.

**Now produce the other two outcomes:**

```bash
# CONNECTION REFUSED -- connect to a port with nothing on it, locally
sudo tcpdump -i any -n 'tcp port 9999' &
nc -v 127.0.0.1 9999
# you'll see: SYN then RST, and nc says "Connection refused"

# TIMEOUT -- connect to a filtered port on the internet
sudo tcpdump -i any -n 'host example.com and tcp port 81' &
nc -v -w 5 example.com 81
# you'll see: SYN ... SYN ... SYN (retries), no answer
```

**Watch the retry timing** on the timeout case: roughly 1 s, then 2 s, then 4 s.
That's **exponential backoff**, and it's why a firewalled port takes so long to
fail while a refused port fails instantly.

### Build it in Go (20 min) — a handshake outcome classifier

The three outcomes above (connected, refused, timed out) each surface as a
different error value in code. Production code that treats them all as
"error" throws away the most useful diagnostic information it has. This
program tells them apart:

```go
// dialcheck: try a TCP connection and say WHICH of the handshake outcomes
// from Chapter 21 happened -- connected, refused, timed out, or something else.
//
//	go run ./dialcheck example.com:443 127.0.0.1:9 example.com:81 nosuchhost.invalid:80
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {
	for _, target := range os.Args[1:] {
		fmt.Printf("%-28s %s\n", target, check(target, 3*time.Second))
	}
}

func check(target string, timeout time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", target)
	took := time.Since(start).Round(time.Millisecond)
	if err == nil {
		defer conn.Close()
		return fmt.Sprintf("CONNECTED in %v  (%s -> %s)", took, conn.LocalAddr(), conn.RemoteAddr())
	}

	var dnsErr *net.DNSError
	var netErr net.Error
	switch {
	case errors.As(err, &dnsErr):
		return fmt.Sprintf("DNS FAILURE (never sent a SYN): %v", dnsErr.Err)
	case errors.Is(err, syscall.ECONNREFUSED):
		return fmt.Sprintf("REFUSED in %v: host answered with RST -> nothing listening", took)
	case errors.Is(err, syscall.EHOSTUNREACH), errors.Is(err, syscall.ENETUNREACH):
		return fmt.Sprintf("UNREACHABLE in %v: no route, or an ICMP unreachable came back", took)
	case errors.As(err, &netErr) && netErr.Timeout():
		return fmt.Sprintf("TIMED OUT after %v: SYNs got no answer -> filtered/dropped", took)
	default:
		return fmt.Sprintf("OTHER after %v: %v", took, err)
	}
}
```

```text
$ go run ./dialcheck example.com:443 127.0.0.1:9 example.com:81 nosuchhost.invalid:80
example.com:443              CONNECTED in 93ms  (192.168.1.2:62150 -> 172.66.147.243:443)
127.0.0.1:9                  REFUSED in 0s: host answered with RST -> nothing listening
example.com:81               TIMED OUT after 3s: SYNs got no answer -> filtered/dropped
nosuchhost.invalid:80        DNS FAILURE (never sent a SYN): no such host
```

**What to notice:**

- **`errors.Is(err, syscall.ECONNREFUSED)`** works because Go wraps the kernel's
  error number inside `*net.OpError` → `*os.SyscallError`. `errors.Is` and
  `errors.As` unwrap the chain. See the
  [Go guide's error-handling chapter](../../Golang/real-life-golang-guide.md#10-error-handling-the-pattern-that-shapes-all-go-code) for the pattern.
- **Connect time ≈ 1 RTT.** The 93 ms "CONNECTED" figure is the handshake
  you captured above. Your app can measure RTT without ICMP at all, which
  is useful where `ping` is blocked.
- **The timeout is the caller's decision.** Without the `context.WithTimeout`,
  the kernel would keep retrying SYNs for about 2 minutes on Linux
  (`net.ipv4.tcp_syn_retries = 6`). Always put a deadline on dials.

**Exercises:**

1. Lower the timeout to 500 ms and dial `example.com:81` again. Then raise it
   to 10 s and count the SYN retransmissions with `tcpdump`. Do they follow
   1 s, 2 s, 4 s?
2. Make `check` run on a list of targets **concurrently** with a
   `sync.WaitGroup`. A filtered port shouldn't hold up the rest. Chapter 53's
   `policytest` lab is this idea grown up.
3. Add an IPv6-vs-IPv4 column: dial `tcp4` and `tcp6` separately. Hosts
   whose IPv6 path is broken are why browsers use *Happy Eyeballs* (racing
   both).

**Across the series:** the Go plan's
[Day 32 TCP echo server](../../Golang/detailed-90-day-plan/week5.md#day-32-tcp-fundamentals) builds the *server* side of this
handshake, and [Day 52](../../Golang/detailed-90-day-plan/week8.md#day-52-tcp-internals-what-goes-wrong) uses these states to debug
CLOSE_WAIT leaks (see also Chapter 23 here).

### Common confusions

- **"Why not two packets?"** After two, the *client* knows the server heard it,
  but the *server* doesn't know the client heard the reply. Three is the minimum
  for both to be sure.
- **"Is the handshake encrypted?"** No — TCP is plaintext. Encryption (TLS) is
  negotiated *after*, on top (Chapter 29).
- **"What's a SYN flood?"** An attack that sends many SYNs and never completes
  the handshake, filling the server's table of half-open connections. Modern
  defence: **SYN cookies** (Chapter 38).

### Check yourself

1. Why does the handshake need three packets rather than two?
2. What do SYN and ACK mean?
3. You send a SYN and get an RST. What does that tell you?
4. You send a SYN and get nothing, three times. What does that tell you?
5. Name two things negotiated during the handshake besides "hello".

### Further reading

- **Video (~12 min):** "TCP 3-Way Handshake" — Practical Networking. Watch after
  doing the capture; it'll confirm what you saw.
- **Article:** "What is TCP?" and "TCP handshake" — Cloudflare Learning Center.
- **RFC 9293** — the current TCP specification (August 2022; it replaced the
  original RFC 793 from 1981). Long, but section 3.5 on connection establishment
  is readable and authoritative.
- **Interactive:** Wireshark's "Follow TCP Stream" on your own capture — see
  Chapter 32.

---

## Chapter 22 — TCP part 2: how it never loses your data

### In one sentence

TCP numbers every byte, the receiver acknowledges what it has, and the sender
re-sends anything not acknowledged — that's the whole reliability mechanism.

### The problem

The network (IP) will lose, duplicate, reorder, and delay packets. TCP must
deliver a perfect, ordered byte stream on top of that.

### The idea, in plain language

Numbered pages in a manuscript sent by post:

```
   SENDER: numbers every page and posts them.
   RECEIVER: "I have everything up to page 500. Send me 501 next."
   SENDER: notices page 501 was never acknowledged -> posts it again.
```

Two mechanisms: **sequence numbers** (which byte is this?) and
**acknowledgements** (which bytes do I have?).

### How it actually works

**Sequence numbers count bytes, not packets.**

```
   Sender transmits 4 segments of 1000 bytes each, starting at byte 1:

     segment A:  bytes    1 - 1000
     segment B:  bytes 1001 - 2000
     segment C:  bytes 2001 - 3000     <-- gets LOST
     segment D:  bytes 3001 - 4000
```

**Acknowledgements are cumulative** — "I have everything up to here":

```
   Receiver gets A  ->  sends ACK 1001    ("I have through 1000, send me 1001")
   Receiver gets B  ->  sends ACK 2001
   Receiver gets D  ->  sends ACK 2001    ("still waiting for 2001!")   <- DUPLICATE
```

Notice: the receiver *has* segment D and holds onto it, but it can't acknowledge
past the gap. It keeps repeating the same ACK.

### Two ways the sender notices a loss

**1. Duplicate ACKs → fast retransmit (the quick way)**

```
   The sender sees ACK 2001 ... again ... and again.
   After the THIRD duplicate ACK, it concludes: "2001 must be lost."
   -> retransmits segment C immediately, without waiting.

   When C arrives, the receiver finally jumps:  ACK 4001
   (because it already had D buffered)
```

**2. Timeout → retransmission (the slow way)**

If nothing comes back at all (e.g. the *last* segment was lost, so there are no
later segments to trigger duplicate ACKs), the sender waits for a timer — the
**RTO** (retransmission timeout) — and then re-sends.

The RTO is calculated from measured round-trip times, with a margin for
variability. It's typically a minimum of ~200 ms, which is *slow* compared with
fast retransmit. **This is why a lost final packet can cost you 200 ms** while a
lost middle packet costs almost nothing.

### SACK: telling the sender exactly what's missing

Plain cumulative ACKs can only say "I'm stuck at 2001" — not "but I also have
3001–4000". **Selective Acknowledgement (SACK)**, negotiated in the handshake,
lets the receiver say precisely which ranges it has:

```
   ACK 2001, SACK [3001-4001]
   = "I need 2001-3000, but don't bother re-sending 3001-4000, I have it."
```

Without SACK, on a lossy link, the sender re-sends far more than necessary.
**On high-latency or lossy paths, SACK makes an enormous difference** — a
satellite or intercontinental link with SACK stripped can be 10x slower.

### Practice (30 min)

**A. See sequence numbers and ACKs in a real transfer:**

```bash
sudo tcpdump -i any -n -c 40 'tcp port 443 and host example.com'
# in another terminal:
curl -s https://example.com > /dev/null
```

Follow the `seq` and `ack` numbers down the capture. Watch `seq` climb by the
length of each data packet, and `ack` climb to match.

**B. Confirm SACK was negotiated** — look at the SYN packets from Chapter 21 for
`sackOK`. Then check your system's setting:

```bash
sysctl net.ipv4.tcp_sack        # Linux; should be 1
```

**C. Cause real retransmissions and watch them.** (Linux only; this deliberately
degrades your network for a moment — do it on a test machine or be ready to undo.)

```bash
# Add 10% packet loss to outgoing traffic
sudo tc qdisc add dev eth0 root netem loss 10%

# Now download something and capture
sudo tcpdump -i eth0 -n 'host example.com' -w /tmp/loss.pcap &
curl -s https://example.com > /dev/null
sudo pkill tcpdump

# IMPORTANT: remove the impairment
sudo tc qdisc del dev eth0 root
```

Open `/tmp/loss.pcap` in Wireshark. It will highlight retransmissions and
duplicate ACKs in black/red. **You have just watched TCP heal itself.**

**D. Check your machine's retransmission counters:**

```bash
# Linux
nstat -az | grep -Ei 'retrans|TCPLostRetransmit|TCPTimeouts'
# macOS
netstat -s -p tcp | grep -i retrans
```

A small number of retransmissions is normal. A *rising percentage* of segments
retransmitted (over ~1%) indicates real trouble.

### Build it in Go (25 min) — TCP is a byte stream, not a message stream

Everything in this chapter (sequence numbers, ACKs, retransmissions) works on
**bytes**, not on your messages. TCP doesn't know or care where one of your
`Write` calls ended and the next began. This surprises nearly every developer
the first time they write a protocol by hand:

```go
// stream: prove that TCP is a byte stream with no message boundaries, then fix
// it with length-prefix framing.
//
//	go run ./stream
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

var messages = []string{"hello", "how are you?", "bye"}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0") // port 0: let the kernel pick a free port
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("== naive: one Write per message, one Read per message?")
	go naiveSender(ln.Addr().String())
	naiveReceiver(ln)

	fmt.Println("\n== framed: 4-byte length prefix before every message")
	go framedSender(ln.Addr().String())
	framedReceiver(ln)
}

func naiveSender(addr string) {
	c, _ := net.Dial("tcp", addr)
	defer c.Close()
	for _, m := range messages {
		c.Write([]byte(m)) // three separate Write calls...
	}
}

func naiveReceiver(ln net.Listener) {
	c, _ := ln.Accept()
	defer c.Close()
	time.Sleep(50 * time.Millisecond) // let all three writes land in the receive buffer
	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			return
		}
		fmt.Printf("Read returned %2d bytes: %q\n", n, buf[:n]) // ...usually ONE Read
	}
}

func framedSender(addr string) {
	c, _ := net.Dial("tcp", addr)
	defer c.Close()
	w := bufio.NewWriter(c)
	for _, m := range messages {
		binary.Write(w, binary.BigEndian, uint32(len(m)))
		w.WriteString(m)
	}
	w.Flush() // one syscall for everything: framing makes batching safe
}

func framedReceiver(ln net.Listener) {
	c, _ := ln.Accept()
	defer c.Close()
	r := bufio.NewReader(c)
	for {
		var size uint32
		if err := binary.Read(r, binary.BigEndian, &size); err != nil {
			return // io.EOF: the sender closed cleanly
		}
		if size > 1<<20 {
			fmt.Println("refusing absurd frame size", size) // never trust a length from the wire
			return
		}
		msg := make([]byte, size)
		if _, err := io.ReadFull(r, msg); err != nil { // ReadFull loops over short reads
			fmt.Println("truncated frame:", err)
			return
		}
		fmt.Printf("message: %q\n", msg)
	}
}
```

```text
$ go run ./stream
== naive: one Write per message, one Read per message?
Read returned 20 bytes: "hellohow are you?bye"

== framed: 4-byte length prefix before every message
message: "hello"
message: "how are you?"
message: "bye"
```

**What to notice:**

- **Three writes arrived as one read.** The bytes were all in the receive
  buffer by the time the reader looked, so one `Read` returned all of them.
  On a slow or lossy network the opposite also happens: one write can arrive
  across several reads. Both are correct TCP behaviour.
- **`io.ReadFull`** keeps reading until the buffer is full. `Read` alone may
  legally return fewer bytes, and assuming it won't is a classic production
  bug that only shows up under load or on bad networks.
- **Validate lengths that come from the wire.** A peer that sends a length of
  4 GB would otherwise make you allocate 4 GB. The `1<<20` check is a
  one-line denial-of-service fix.
- Real protocols choose either **delimiters** (HTTP/1.1 headers end with
  `\r\n\r\n`; Redis uses `\r\n`) or **length prefixes** (HTTP/2 frames,
  gRPC messages, TLS records, Postgres, Kafka). Chapter 39 shows what happens
  when two parsers disagree about framing.

**Exercises:**

1. Remove the `time.Sleep` in `naiveReceiver` and run it 20 times. Does the
   split change from run to run? (It depends on timing, which is exactly the
   problem.)
2. Switch the framed version to newline delimiters with `bufio.Scanner`. Then
   send a message that *contains* a newline and watch it break. That's why
   binary protocols use length prefixes.
3. The [Go guide's framing section](../../Golang/real-life-golang-guide.md#34-networking-deep-dive-net-conn-tcp-udp-framing-your-own-protocol) and the Go plan's
   [binary KV protocol (Day 53)](../../Golang/detailed-90-day-plan/week8.md#day-53-binary-protocol-design) take this to a full
   protocol with versioning and message types. Build the KV server.

### Common confusions

- **"Does TCP retransmit every lost packet?"** It retransmits any *data* not
  acknowledged. Pure ACK packets are never retransmitted — a later ACK supersedes
  them.
- **"Do sequence numbers start at 1?"** No — they start at a **random** value
  chosen for each connection, for security (so an attacker can't guess them). The
  `-S` flag in tcpdump shows the real absolute numbers; by default it shows
  friendly relative ones.
- **"Retransmissions mean the network is broken."** A few are completely normal.
  It's the *rate* and the *trend* that matter.

### Check yourself

1. Do TCP sequence numbers count packets or bytes?
2. What is a duplicate ACK, and what does the sender do after three of them?
3. Why is a lost *last* packet more expensive than a lost middle packet?
4. What does SACK add, and when does it matter most?

### Further reading

- **Video:** "TCP Reliability" / "TCP Retransmissions Explained" — Practical
  Networking; also Wireshark's own YouTube channel has good retransmission
  analysis videos.
- **Article:** "TCP Analysis" chapter of the Wireshark User Guide — explains
  every "Expert Info" flag you'll see.
- **Book:** *TCP/IP Illustrated, Volume 1* — Fall & Stevens. **The** reference for
  this material. Dense but authoritative; dip into the TCP chapters.
- **Book (gentler):** Kurose & Ross Chapter 3, sections on reliable data transfer.

---

## Chapter 23 — TCP part 3: closing a connection, and the states

### In one sentence

Closing is a polite four-packet goodbye (each side closes its own direction), and
the connection then sits in **TIME_WAIT** for about a minute to catch stragglers.

### How closing works

TCP connections are really **two independent one-way streams**, so each direction
is closed separately.

```
   CLIENT                                          SERVER
        |                                              |
        |  1. FIN     "I'm done sending."              |
        |--------------------------------------------->|
        |                                              |
        |  2. ACK     "Got it."                        |
        |<---------------------------------------------|
        |         (server may still be sending!)       |
        |                                              |
        |  3. FIN     "Now I'm done too."              |
        |<---------------------------------------------|
        |                                              |
        |  4. ACK     "Got it."                        |
        |--------------------------------------------->|
        |                                              |
        |  [ TIME_WAIT for ~60 seconds ]           [ CLOSED ]
        |  then CLOSED                                 |
```

Often it's only **three** packets, because the server combines its ACK and FIN
into one.

There's also the impolite version: **RST** (reset). One packet, no negotiation,
connection dead immediately. Any buffered data is lost. You'll see "Connection
reset by peer" at the application.

### TIME_WAIT: why your machine holds on

The side that closes *first* waits about **60 seconds** in `TIME_WAIT` before
fully releasing the connection. Two reasons:

1. **Catch stragglers.** If the final ACK was lost, the other side will re-send
   its FIN. `TIME_WAIT` is there to answer it rather than replying with a
   confusing RST.
2. **Prevent confusion.** Stop a *brand-new* connection using the same four
   numbers from receiving delayed packets belonging to the *old* one.

> **On Linux this is hardcoded at 60 seconds** (`2 × MSL`) and is **not**
> controlled by `net.ipv4.tcp_fin_timeout` — that setting governs a different
> state (`FIN_WAIT_2`). This is a very common misconception; you'll see plenty of
> blog posts get it wrong.

**Is a lot of TIME_WAIT bad?** On a *server*, usually harmless. On a *client*
making thousands of short connections per second, it can exhaust the ephemeral
port range (~28,000 ports to a single destination). The real fix is **connection
reuse** (keep-alive), not tuning kernel timers.

### The states you'll actually see

```bash
ss -tan | awk '{print $1}' | sort | uniq -c | sort -rn
```

| State | Meaning | If you see lots of them |
|---|---|---|
| `LISTEN` | A server waiting for connections | Normal |
| `ESTAB` | Connection open, data can flow | Normal |
| `SYN-SENT` | Sent SYN, waiting for reply | Destination unreachable/firewalled |
| `SYN-RECV` | Got SYN, waiting for final ACK | Possible SYN flood, or clients disappearing |
| `TIME-WAIT` | Waiting out the 60 s after closing | Normal for busy clients; see above |
| `CLOSE-WAIT` | The peer closed; **waiting for the local app to close** | **A bug in your application** — it forgot to close the socket. Leaks file descriptors. |
| `FIN-WAIT-2` | We closed; waiting for the peer to close | The **peer** app isn't closing properly |

**The `CLOSE_WAIT` rule is worth memorising:** growing `CLOSE_WAIT` means *your*
code isn't closing sockets. Growing `FIN_WAIT_2` means *their* code isn't.

### Practice (25 min)

**A. Watch a full close:**

```bash
sudo tcpdump -i any -n 'tcp port 443 and host example.com'
curl -s https://example.com > /dev/null
```

At the end of the capture, find the `[F.]` (FIN+ACK) packets and the final `[.]`.

**B. Watch a reset:**

```bash
sudo tcpdump -i any -n 'tcp port 9999' &
nc -v 127.0.0.1 9999          # nothing is listening
```
You'll see `[R.]` — an RST.

**C. See the state histogram on your own machine:**

```bash
ss -tan | awk '{print $1}' | sort | uniq -c | sort -rn
```

Then open a lot of browser tabs and run it again.

**D. Create TIME_WAIT deliberately:**

```bash
for i in $(seq 1 20); do curl -s -o /dev/null http://example.com; done
ss -tan state time-wait | wc -l
# wait 60 seconds, run again -- they disappear
```

**E. Find any CLOSE_WAIT on your system** (and identify the guilty program):

```bash
ss -tanp state close-wait          # Linux
lsof -i -sTCP:CLOSE_WAIT           # macOS/Linux
```

If you find some that persist, you've found a real bug in some software.

### Build it in Go (30 min) — manufacture a CLOSE_WAIT leak, then fix it

TIME_WAIT is normal. **CLOSE_WAIT that keeps growing is always a bug in your
program**: the peer has closed, the kernel has acknowledged its FIN, and the
socket is waiting for *your code* to call `close()`. It will wait forever. Each
leaked socket holds a file descriptor, and eventually the process hits its
open-files limit and stops accepting anything (`too many open files`).

```go
// closewait: manufacture the CLOSE_WAIT leak from Chapter 23, watch it with
// ss/lsof, then fix it -- and see a clean half-close with CloseWrite.
//
//	go run ./closewait -leak     # server forgets to Close
//	go run ./closewait           # server closes properly
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	leak := flag.Bool("leak", false, "simulate the bug: never close accepted connections")
	flag.Parse()

	ln, err := net.Listen("tcp", "127.0.0.1:7070")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("server pid %d on %s\n", os.Getpid(), ln.Addr())
	fmt.Println("in another terminal run:")
	fmt.Println("  Linux:  watch -n1 \"ss -tan state close-wait '( sport = :7070 )'\"")
	fmt.Println("  macOS:  watch -n1 \"netstat -an -p tcp | grep 7070\"")

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go handle(c, *leak)
		}
	}()

	// Ten clients connect, send a request, half-close, read the reply, leave.
	for i := 0; i < 10; i++ {
		c, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(c, "request %d\n", i)
		c.(*net.TCPConn).CloseWrite()                   // send FIN: "I'm done talking, but still listening"
		reply, _ := bufio.NewReader(c).ReadString('\n') // the server can still answer after our FIN
		fmt.Printf("client %d got %q\n", i, reply)
		c.Close()
	}
	fmt.Println("clients finished. Server stays up 60s so you can look; Ctrl-C to quit.")
	time.Sleep(60 * time.Second)
}

func handle(c net.Conn, leak bool) {
	req, _ := io.ReadAll(c) // returns when the client's FIN arrives (EOF)
	fmt.Fprintf(c, "ok: %s", req)
	if leak {
		// BUG: no c.Close(). The client's FIN was acknowledged, so the kernel
		// moves our side to CLOSE_WAIT and waits for US to close. It will wait
		// forever, holding a file descriptor and kernel memory.
		return
	}
	c.Close()
}
```

Real output on Linux:

```text
$ go run ./closewait -leak &
$ ss -tan state close-wait '( sport = :7070 )'
Recv-Q Send-Q Local Address:Port Peer Address:Port
0      0          127.0.0.1:7070    127.0.0.1:55228
0      0          127.0.0.1:7070    127.0.0.1:55258
...                                                  <- 10 sockets, and they never go away

$ go run ./closewait            # the fixed version
$ ss -tan state close-wait '( sport = :7070 )' | tail -n +2 | wc -l
0
```

**What to notice:**

- **The client side is in `FIN_WAIT_2`** (run `netstat -an | grep 7070` and
  you'll see both halves). The client sent its FIN and got an ACK, and is
  now waiting for the server's FIN, which never comes.
- **`CloseWrite()` is a half-close**: it sends FIN ("I've finished sending") but
  keeps reading. The server's `io.ReadAll` returns when that FIN arrives,
  and the server can *still reply*. HTTP/1.0 used exactly this pattern, and
  so do `nc -N` and many proxies.
- **In Go, the fix is almost always a missing `defer c.Close()`**, often on
  an error path, or a `resp.Body` from `net/http` that was never closed. That
  last one leaks the same way on the *client* side.

**Exercises:**

1. While `-leak` runs, check the file descriptors: `ls /proc/$(pgrep closewait)/fd | wc -l`.
   Then lower the limit with `ulimit -n 64` and change the loop to 100
   clients. What error do you get?
   ([Linux guide, Chapter 43](../../os-linux/real-life-os-guide.md#chapter-43-resource-limits-ulimit-and-your-first-taste-of-cgroups) covers `ulimit`.)
2. Write the same bug with `net/http`: a client that does
   `resp, _ := http.Get(url)` and never closes `resp.Body`. Watch the
   connections pile up in `ESTABLISHED` instead. Which side is now
   leaking?
3. The Go plan's [Day 52 exercise](../../Golang/detailed-90-day-plan/week8.md#day-52-tcp-internals-what-goes-wrong) has a production-style
   version with a leak detector. Do it next.

### Common confusions

- **"I have 20,000 TIME_WAIT sockets, is that bad?"** Usually not on a server.
  Check whether you're near port exhaustion before "fixing" it. Never enable
  `tcp_tw_recycle` (it's removed from modern Linux and broke NAT users badly).
- **"Why does the connection stay around after I close the app?"** TIME_WAIT is
  kernel state, not application state; it outlives the process deliberately.
- **"FIN vs RST?"** FIN = polite "I've finished, here's everything". RST = "stop
  immediately", data may be lost.

### Check yourself

1. Why does closing take four packets when opening takes three?
2. Which side enters TIME_WAIT, and for roughly how long?
3. What does a growing number of `CLOSE_WAIT` sockets tell you?
4. Does `tcp_fin_timeout` control the TIME_WAIT duration?

### Further reading

- **Article (definitive):** "Coping with the TCP TIME-WAIT state on busy Linux
  servers" — Vincent Bernat. Corrects most of the internet's bad advice on this
  topic. Read it if you ever tune a busy server.
- **Video:** "TCP Connection Termination" — Practical Networking.
- **Diagram:** search for the "TCP state transition diagram" (the classic one from
  RFC 793 / 9293). Print it and put it on your wall.

---

### End of Part 5 — Milestone check

- [ ] I can explain the 4-tuple and why it identifies a connection
- [ ] I know the difference between binding to `0.0.0.0` and `127.0.0.1`
- [ ] I can explain when UDP is the *right* choice
- [ ] **I have captured a TCP three-way handshake and read it**
- [ ] I can tell "refused" from "timeout" and know what each means
- [ ] I can explain sequence numbers, ACKs, and fast retransmit
- [ ] I know what `CLOSE_WAIT` piling up means

---

# Part 6 — Making TCP fast (and why yours isn't)

TCP works. This part is about why it's sometimes *slow*, and it explains most
real-world performance complaints.

## Chapter 24 — Flow control: the sliding window

### In one sentence

The receiver continuously tells the sender "you may send this many more bytes
before waiting for me", which stops a fast sender from drowning a slow receiver.

### The problem

A powerful server sends to a small IoT device. If the server blasts at full speed,
the device's memory fills and everything is dropped — wasting the whole transfer.
The receiver needs a way to say "slow down".

### The idea, in plain language

A conveyor belt with a bucket at the end:

```
   SENDER ---> [ items on the belt ] ---> RECEIVER's bucket
                                              |
                          "my bucket has room for 8 more" <---+
```

The receiver announces how much room it has left. The sender never sends more than
that.

### How it actually works

Every TCP packet carries a **window** field: *"I can accept this many more bytes
right now."*

```
   Bytes:  ...already sent    |  sent, not yet  |   may send    | can't send
           and acknowledged   |   acknowledged  |    now        |    yet
        ---------------------+-----------------+---------------+------------->
                             ^                                 ^
                        oldest unacked                  end of window

        <---------------- WINDOW (what the receiver allows) ---->
```

As acknowledgements arrive, the whole window **slides** to the right — hence
"sliding window". If the receiving application stops reading, the window shrinks
toward zero and the sender pauses automatically. **You get backpressure for free.**

### The 64 KB wall (and why window scaling matters)

The window field in the TCP header is only **16 bits** — maximum 65,535 bytes.

Now here's the crucial formula:

```
   maximum throughput  =  window size  /  round-trip time
```

| Round-trip time | Max speed with a 64 KB window |
|---|---|
| 1 ms (same building) | ~500 Mbit/s |
| 30 ms (same country) | **~17 Mbit/s** |
| 100 ms (intercontinental) | **~5 Mbit/s** |
| 300 ms (satellite) | **~1.7 Mbit/s** |

**Read that table again.** On a 100 ms link, a 64 KB window caps you at 5 Mbit/s
*no matter how fast your connection is*. You could have a 10 Gbit/s fibre and get
5 Mbit/s.

The fix is the **window scale** option, negotiated in the handshake (Chapter 21).
It multiplies the window by a power of two, allowing windows up to about 1 GB.
It's on by default everywhere — **but if a middlebox strips it, you get the table
above**, and this is a real and maddening cause of "the link is fast but transfers
are slow".

### Bandwidth-delay product

The window you *need* is called the **bandwidth-delay product** (BDP):

```
   BDP = bandwidth x round-trip time     ("how much data fits in the pipe")

   Example: 1 Gbit/s link, 80 ms round trip
          = 1,000,000,000 bits/s x 0.08 s
          = 80,000,000 bits = 10 MB

   You need a ~10 MB window to keep that pipe full.
```

This is why long-distance high-speed transfers need big buffers, and why default
settings that are fine locally are terrible internationally.

### Zero window

If the receiving application stops reading entirely, the window hits zero and the
sender stops. The sender then sends periodic tiny "window probes" to ask "any room
yet?", so a lost window update can't deadlock the connection.

**Seeing "TCP ZeroWindow" in Wireshark means: the receiving application is too
slow.** It is *not* a network problem. Very useful to know — it immediately points
you at the app, not the network.

### Practice (25 min)

**A. See the window in your own traffic:**

```bash
sudo tcpdump -i any -n 'tcp port 443' -c 20
```

Look at the `win` value on each packet. After the handshake, the numbers are
**scaled** — multiply by 2^(scale factor from the handshake) for the real value.

**B. Look at live connection details** (Linux — this is a great command):

```bash
ss -tin
```

Look for `wscale:7,7` (window scaling active), `rtt:` (measured round-trip time),
`cwnd:` (Chapter 25), and `send X Mbps` (its own throughput estimate).

**C. Prove the window/RTT formula.** Compare a download from something nearby with
something far away:

```bash
curl -o /dev/null -s -w 'speed: %{speed_download} bytes/s   time: %{time_total}s\n' \
  https://speed.cloudflare.com/__down?bytes=10000000
```

Then compare `ping` times to servers at different distances. The relationship
between latency and achievable throughput will be visible.

**D. Check your system's buffer limits:**

```bash
sysctl net.ipv4.tcp_rmem net.ipv4.tcp_wmem net.core.rmem_max   # Linux
sysctl net.inet.tcp.recvspace net.inet.tcp.sendspace           # macOS
```

The third number in `tcp_rmem` is the maximum receive buffer. If it's smaller than
your BDP, that's your speed limit.

### Build it in Go (25 min) — drive the window to zero

Flow control is easiest to understand from the sender's side: **when the
receiver stops reading, the sender's `Write` stops returning.** This lab
connects a fast writer to a deliberately slow reader (4 KB every 500 ms) with
small buffers, so the window fills within a second.

```go
// zerowindow: a receiver that reads slowly, so you can watch flow control
// (Chapter 24) fill the receive buffer, shrink the window to zero, and make
// the sender's Write block.
//
//	go run ./zerowindow
//	# meanwhile:  ss -tnoi '( sport = :7071 or dport = :7071 )'   (Linux)
//	#             netstat -an -p tcp | grep 7071                   (macOS: watch Recv-Q / Send-Q)
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:7071")
	if err != nil {
		panic(err)
	}

	go func() { // the SLOW receiver: 4 KB every 500 ms
		c, _ := ln.Accept()
		c.(*net.TCPConn).SetReadBuffer(64 * 1024) // a small window, so it fills fast
		buf := make([]byte, 4096)
		for {
			time.Sleep(500 * time.Millisecond)
			if _, err := c.Read(buf); err != nil {
				return
			}
		}
	}()

	c, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		panic(err)
	}
	c.(*net.TCPConn).SetWriteBuffer(64 * 1024)

	chunk := make([]byte, 16*1024)
	total := 0
	for i := 0; i < 64; i++ {
		start := time.Now()
		n, err := c.Write(chunk) // returns once the bytes are in OUR send buffer
		if err != nil {
			fmt.Println(err)
			return
		}
		total += n
		if wait := time.Since(start); wait > 50*time.Millisecond {
			fmt.Printf("write #%-2d BLOCKED %-6v (total %4d KB) <- receiver's window is full\n",
				i, wait.Round(time.Millisecond), total/1024)
		} else {
			fmt.Printf("write #%-2d instant        (total %4d KB)\n", i, total/1024)
		}
	}
}
```

The program's own output shows writes going from instant to blocked. The
kernel's view is the interesting part (Linux, during the run):

```text
$ ss -tnoi '( sport = :7071 or dport = :7071 )'
ESTAB 79872  0        127.0.0.1:7071   127.0.0.1:51062           <- receiver: 78 KB unread (Recv-Q)
ESTAB 0      128000   127.0.0.1:51062  127.0.0.1:7071  timer:(persist,701ms,0)
     ... backoff:4 ... rwnd_limited:6008ms(100.0%) ... notsent:128000
```

**How to read it:**

- **Receiver `Recv-Q 79872`**: bytes the kernel holds that the application
  hasn't `Read` yet. The window it advertises shrinks by exactly this amount.
- **Sender `Send-Q 128000` / `notsent:128000`**: bytes the sender has accepted
  from `Write` but can't send, because the window is closed.
- **`timer:(persist,...)`**: the **zero-window probe** timer from this
  chapter. The sender periodically asks "is there room yet?", backing off
  each time (`backoff:4`).
- **`rwnd_limited:...(100.0%)`**: Linux is telling you directly that the whole
  connection was limited by the *receiver's window*, not by the network.
  When you see this in production, look at the receiving application,
  not the network.

**Exercises:**

1. Change the reader's sleep to 0 and watch every write become instant.
   The network didn't change; the application did.
2. Remove the `SetReadBuffer`/`SetWriteBuffer` calls. Linux autotunes the
   buffers upwards, so it takes longer to block. How many KB get through
   before the first blocked write?
3. In a real service this looks like a consumer that's too slow, for example
   a log shipper behind a slow disk. Find the equivalent in the
   [Linux guide's I/O stack chapter](../../os-linux/real-life-os-guide.md#chapter-14-the-i-o-stack-from-read-to-the-disk-platter): a slow disk
   produces exactly this kind of backpressure.

### Common confusions

- **"Flow control vs congestion control?"** Flow control protects the **receiver**
  (Chapter 24). Congestion control protects the **network** (Chapter 25). Both
  limit the sender; the sender obeys whichever is smaller.
- **"Bigger window is always better?"** For throughput on a long link, yes. But
  oversized buffers *in the network* cause bufferbloat (Chapter 27).

### Check yourself

1. What does the window field tell the sender?
2. Why does a 64 KB window limit you to ~5 Mbit/s on a 100 ms link?
3. What does "TCP ZeroWindow" tell you about where the problem is?
4. Calculate the BDP for a 100 Mbit/s link with 40 ms RTT.

### Further reading

- **Video:** "TCP Flow Control" — Practical Networking.
- **Article:** "TCP Window Scaling" — search for the Cloudflare or Packetlife
  explainer.
- **Tool:** `ss -tin` on Linux — the single most informative command for live TCP.

---

## Chapter 25 — Congestion control: sharing the road

### In one sentence

Nobody tells TCP how fast the network is, so it **probes**: speed up until packets
are lost, back off, and repeat — which is why your download speed graph looks like
a sawtooth.

### The problem

Flow control (Chapter 24) protects the *receiver*. But what protects the *network*?
If every sender blasts at full speed, routers' queues overflow, packets are lost,
everyone retransmits, and it gets worse — a feedback loop.

This actually happened. In **October 1986**, the NSFNET backbone's throughput
collapsed from 32 kbit/s to **40 bit/s** — a 1000x drop — because of exactly this.
The fix, invented by Van Jacobson in 1988, is what we still use today.

### The idea, in plain language

Driving in fog. You don't know the speed limit, so you accelerate gently until
something goes wrong, then slow down sharply, then start accelerating again.

**Packet loss is the signal.** TCP treats a lost packet as "the network is full".

### How it actually works

The sender keeps a second limit alongside the receiver's window, called the
**congestion window** (`cwnd`). It sends the **smaller** of the two.

```
   bytes in flight  <=  min( receiver's window , congestion window )
                            \______________/     \_______________/
                             protects the         protects the
                             RECEIVER              NETWORK
```

Two phases:

```
   1. SLOW START  (despite the name, this is the FAST ramp-up)
      Start small (about 10 packets).
      Every round trip, DOUBLE the window.
      10 -> 20 -> 40 -> 80 -> 160 ...
      Keep doubling until a packet is lost, or a threshold is reached.

   2. CONGESTION AVOIDANCE (the careful part)
      Now increase by only ONE packet per round trip.
      Creep upward, probing for more capacity.

   ON PACKET LOSS:
      Loss detected by duplicate ACKs -> cut the window roughly in half
                                          and continue (fast recovery)
      Loss detected by timeout        -> cut all the way back to the start
                                          (something is badly wrong)
```

The result is the famous **sawtooth**:

```
   speed
     ^
     |        /|      /|      /|
     |       / |     / |     / |      <- creep up (congestion avoidance)
     |      /  |    /  |    /  |
     |     /   |   /   |   /   |
     |    /    v  /    v  /    v      <- loss: cut in half
     |   /
     |  /  <- slow start (doubling)
     +-----------------------------------> time
```

### Why this matters to you

**1. Short connections never reach full speed.** A small web page finishes during
slow start. This is why *round trips* dominate web performance, not bandwidth —
and why connection reuse (Chapter 28) matters so much.

**2. Distance is brutal.** The window grows *per round trip*, so a connection with
100 ms RTT ramps up 10x slower than one with 10 ms RTT. Same bandwidth, 10x
slower to get going.

**3. A little loss is catastrophic over distance.** For classic
loss-based TCP, throughput falls roughly with the square root of the loss rate.
Going from 0.01% loss to 1% loss (100x worse) cuts your speed by about **10x**.
On a long path, even 0.1% loss can cap a single connection at a few Mbit/s.

**4. This is why downloads use multiple connections.** Each gets its own
congestion window, so N connections get roughly N times the share.

### The algorithms by name

You'll see these names; you don't need the maths.

| Name | Idea | Where |
|---|---|---|
| **Reno / NewReno** | The classic sawtooth described above | Historical, still the baseline |
| **CUBIC** | Grows as a cubic curve of *time* since the last loss — fairer to long-distance connections | **Default on Linux, Android, and most of the internet** |
| **BBR** | Ignores loss; instead *measures* the bottleneck bandwidth and round-trip time and paces to match | Google, YouTube, many CDNs. Much better on lossy/long paths. |

**BBR is worth knowing about** because it solves point 3 above: on a long,
slightly-lossy path where CUBIC gets a few Mbit/s, BBR can get hundreds.

### Practice (25 min)

**A. See your congestion window live:**

```bash
# Linux: start a big download, then watch in another terminal
curl -o /dev/null https://speed.cloudflare.com/__down?bytes=100000000 &
watch -n0.5 "ss -tin | grep -A1 cloudflare | head -4"
```

Watch `cwnd:` grow. That's slow start happening in front of you.

**B. Check which algorithm you're using:**

```bash
sysctl net.ipv4.tcp_congestion_control          # Linux
sysctl net.ipv4.tcp_available_congestion_control
```

**C. Watch the sawtooth.** Add a little loss and download something, then graph
throughput in Wireshark (`Statistics → I/O Graph`):

```bash
sudo tc qdisc add dev eth0 root netem loss 0.5%
curl -o /dev/null https://speed.cloudflare.com/__down?bytes=50000000
sudo tc qdisc del dev eth0 root          # UNDO
```

**D. Prove that parallel connections beat one.** This is a striking demo:

```bash
# one connection
time curl -s -o /dev/null https://speed.cloudflare.com/__down?bytes=50000000

# four connections in parallel
time ( for i in 1 2 3 4; do
    curl -s -o /dev/null https://speed.cloudflare.com/__down?bytes=12500000 &
  done; wait )
```

Same total bytes. On a long or slightly lossy path, the parallel version is often
much faster. **That is congestion control, measured.**

### Common confusions

- **"Why does TCP need to guess? Can't the network tell it?"** There's a mechanism
  for that (ECN — routers mark packets instead of dropping them), but it needs
  cooperation across the whole path, so deployment has been slow. Loss remains the
  universal signal.
- **"Is packet loss always bad?"** For loss-based congestion control, loss is the
  *only* signal available. Zero loss forever would mean TCP never learns the
  limit. A small amount is normal and healthy.
- **"My connection is 1 Gbit/s but the download is 20 Mbit/s."** Very likely
  congestion control + distance + a little loss, or a window that's too small
  (Chapter 24). Test with parallel connections to find out which.

### Check yourself

1. What signal does classic TCP use to detect congestion?
2. Why is it called "slow start" when it doubles every round trip?
3. Why does a short web page never reach full speed?
4. Why do download managers use multiple connections?

### Further reading

- **Video:** "TCP Congestion Control" — Practical Networking; also
  Computerphile's "Congestion Control".
- **Article:** "TCP BBR" — the Google blog post introducing it (search "BBR
  congestion-based congestion control Google blog"). Very readable.
- **Historical (fascinating):** Van Jacobson's 1988 paper "Congestion Avoidance
  and Control". The paper that saved the internet, and readable by a beginner.
- **Book:** Kurose & Ross Chapter 3, congestion control sections.

---

## Chapter 26 — The 40-millisecond mystery

### In one sentence

Two individually-sensible TCP optimisations — Nagle's algorithm and delayed
acknowledgements — can deadlock each other, adding a fixed 40 ms delay to every
request, and this bug is everywhere.

### The two optimisations

**Nagle's algorithm (sender side).** Stops tiny packets flooding the network:

> *"If I have data still unacknowledged, don't send another small packet — wait
> until I have a full-sized chunk, or until the outstanding data is
> acknowledged."*

Sensible: sending a 1-byte payload with 40 bytes of headers is 97% waste.

**Delayed ACK (receiver side).** Stops pure-acknowledgement packets flooding the
network:

> *"Don't acknowledge immediately. Wait up to ~40 ms in case I have data to send
> back — then I can piggyback the ACK on it for free."*

Also sensible.

### How they deadlock

Now put them together, with an application that sends a request in **two writes**
(very common — a header, then a body):

```
   CLIENT (Nagle on)                          SERVER (delayed ACK on)
   write(header)  -> sent immediately
                     (nothing unacknowledged yet)
                                       ------> got the header.
                                              Nothing to reply with yet.
                                              "I'll wait ~40 ms before ACKing,
                                               in case I get more."

   write(body)    -> Nagle: "the header is
                     unacknowledged and this
                     is small -- I'll WAIT."

        ... BOTH SIDES ARE NOW WAITING FOR EACH OTHER ...

   (40 ms later) the delayed-ACK timer fires  ------> ACK

   ACK arrives -> Nagle releases the body  ------>  now the server can work
```

**Every single request costs an extra 40 milliseconds.** The network is fine. The
CPUs are idle. Nothing is broken. It's just two politeness rules colliding.

### The fingerprint

> **If your latency graph has a hard floor at 40 ms (or 200 ms on some systems),
> and the network round-trip time is under 1 ms, suspect this immediately.**

Typical symptoms:
- Requests take suspiciously *exactly* 40 ms
- Throughput caps at a suspiciously round number like 25 requests/second
  (= 1 / 0.04)
- It only happens on small requests; big ones are fine
- It's worse locally than over the internet (because the real RTT is tiny in
  comparison)

### The fix

Any one of these:

1. **Set `TCP_NODELAY` on the socket** (disables Nagle). This is what virtually
   every RPC framework, database driver, and game server does. One line:

```python
import socket
s = socket.create_connection(("example.com", 80))
s.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
```

2. **Write the whole request in one call.** Build the message in memory, then one
   `send()`. Nagle never engages because there's no small-then-small pattern.
   This is the cleanest fix.

3. `TCP_QUICKACK` on the receiver (Linux) to suppress delayed ACKs.

### Practice (20 min)

**A. Check whether your tools set it.** Many do:

```bash
# Linux -- watch a program's socket options as it runs
strace -f -e trace=setsockopt curl -s https://example.com 2>&1 | grep -i nodelay
```

**B. Measure a suspiciously round latency.** If you have a local service, time
many small requests:

```bash
for i in $(seq 1 20); do
  curl -s -o /dev/null -w '%{time_total}\n' http://localhost:8000/
done | sort -n
```

Values clustered at 0.04 s with a local server = the fingerprint.

**C. Recognise the sibling problem.** Also watch for delays that are exact
multiples of common timers: 200 ms (retransmission timeout, Chapter 22),
1 s / 3 s / 7 s (SYN retries, Chapter 21), 5 s (a dead DNS server, Chapter 18).
**Suspiciously round numbers in latency are always a timer, never a coincidence.**

### Build it in Go (20 min) — the 40 ms mystery, on demand

In Go, the fix above is already applied. Every TCP connection Go creates has
**`TCP_NODELAY` set by default**. That's why Go services rarely hit this bug,
and why it's confusing when they talk to a peer written in something else. To
see the problem, this lab turns Nagle back *on* with `SetNoDelay(false)` and
uses the classic "write-write-read" pattern:

```go
// nagle: reproduce Chapter 26's "40 ms mystery" on purpose.
//
// Go sets TCP_NODELAY on every TCP connection by default (Nagle OFF), so
// Go programs rarely hit this -- unless someone turns Nagle back on, or the
// peer is a program in another language. Here we turn it on deliberately.
//
//	go run ./nagle
package main

import (
	"bufio"
	"fmt"
	"net"
	"slices"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go echoServer(ln)

	for _, noDelay := range []bool{true, false} {
		c, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			panic(err)
		}
		tc := c.(*net.TCPConn)
		tc.SetNoDelay(noDelay) // false = Nagle's algorithm ON
		r := bufio.NewReader(c)

		var lat []time.Duration
		for i := 0; i < 20; i++ {
			start := time.Now()
			// The classic bad pattern: header and body in TWO small writes,
			// then wait for a reply ("write-write-read").
			c.Write([]byte("HEADER:"))
			c.Write([]byte("body\n"))
			if _, err := r.ReadString('\n'); err != nil {
				panic(err)
			}
			lat = append(lat, time.Since(start))
		}
		c.Close()
		slices.Sort(lat)
		fmt.Printf("TCP_NODELAY=%-5v  p50=%-10v max=%v\n", noDelay, lat[len(lat)/2], lat[len(lat)-1])
	}
	fmt.Println("\nIf NODELAY=false is ~40ms (Linux) or higher, you just met Nagle + delayed ACK.")
	fmt.Println("Fix: one write per message (bufio.Writer + Flush), or keep NODELAY on.")
}

func echoServer(ln net.Listener) {
	for {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		go func() {
			defer c.Close()
			r := bufio.NewReader(c)
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					return
				}
				c.Write([]byte(line))
			}
		}()
	}
}
```

```text
$ go run ./nagle                  # Linux
TCP_NODELAY=true   p50=5.709µs    max=49.25µs
TCP_NODELAY=false  p50=43.952667ms max=46.074958ms
```

**Four orders of magnitude**, caused by one socket option and two small
writes. On macOS loopback the effect often doesn't appear, because delayed
ACK behaves differently there. Run it on Linux (or in the container from
Section 0.8) to see the real fingerprint.

**What to notice:**

- The fix in the output, **"one write per message"**, is the cleanest one.
  `bufio.Writer` collects the header and body and `Flush()` sends them in a
  single `write()`, so Nagle never has a small packet to hold back.
- `SetNoDelay` is a thin wrapper around
  `setsockopt(fd, IPPROTO_TCP, TCP_NODELAY, ...)`, the same call the Python
  snippet above makes.

**Exercises:**

1. Keep `SetNoDelay(false)` but wrap the connection in a `bufio.Writer` and
   `Flush()` after the body. Does the 40 ms disappear?
2. Run the server under `strace -e trace=setsockopt` and confirm Go set
   `TCP_NODELAY` on the *accepted* connection too.
3. Throughput is 1 / 44 ms ≈ 23 requests per second per connection. Find the
   "suspiciously round number" in that.

### Common confusions

- **"Should I always disable Nagle?"** For interactive/request-response traffic,
  yes. For bulk transfer, leaving it on is slightly more efficient. Most modern
  libraries set `TCP_NODELAY` by default now.
- **"Is this still a real problem?"** Less than it was, because frameworks
  default to `TCP_NODELAY` — but it still appears regularly in custom protocols,
  older software, and hand-rolled socket code.

### Check yourself

1. What does Nagle's algorithm do, and why?
2. What does delayed ACK do, and why?
3. Describe how they deadlock.
4. Name two ways to fix it.

### Further reading

- **Article (the classic):** "It's always TCP_NODELAY. Every damn time." —
  search that phrase; several good writeups exist.
- **Article:** John Nagle's own comments on Hacker News about the interaction —
  search "Nagle delayed ACK Hacker News"; the algorithm's inventor explains why
  the combination is the problem, not his algorithm.
- **Docs:** `man 7 tcp` on Linux — the `TCP_NODELAY`, `TCP_CORK`, and
  `TCP_QUICKACK` sections.

---

## Chapter 27 — Latency, bandwidth, and bufferbloat

### In one sentence

Bandwidth is how *wide* the pipe is; latency is how *long* it is — and for most
things you actually notice, latency matters far more.

### The two numbers, and the analogy

```
   BANDWIDTH  = how many lanes the motorway has
                (bits per second: 100 Mbit/s, 1 Gbit/s)

   LATENCY    = how long the motorway is
                (milliseconds: 5 ms, 90 ms, 600 ms)
```

A truck full of hard drives has enormous bandwidth and terrible latency. A phone
call has tiny bandwidth and needs low latency.

**Upgrading bandwidth does nothing for latency.** Going from 100 Mbit/s to
1 Gbit/s does not make a website load 10x faster, because most of a page load is
round trips (Chapter 4), not bytes.

### Where latency comes from

```
   total delay =  propagation  +  serialisation  +  queueing  +  processing
```

| Component | Cause | Can you fix it? |
|---|---|---|
| **Propagation** | Distance ÷ speed of light in fibre (~200,000 km/s ≈ 5 µs per km) | Only by moving closer — this is what CDNs do |
| **Serialisation** | Time to clock the bits onto the wire: a 1500-byte packet takes 12 ms at 1 Mbit/s but 12 µs at 1 Gbit/s | Faster links |
| **Queueing** | Waiting behind other packets in a router's buffer | **Yes — this is the controllable one** |
| **Processing** | Routing lookups, NAT, encryption, VM overhead | Better hardware, fewer hops |

**Physics floor:** London to New York is ~5,600 km. Light in fibre does ~200 km/ms,
so one way is ~28 ms, round trip ~56 ms minimum — and real paths aren't straight,
so ~70–90 ms is typical. **No amount of money makes that faster.** Only moving the
content closer does.

### Bufferbloat: the one you can actually fix

Here's a problem you have probably experienced without knowing its name.

**The symptom:** when someone starts a big upload or download, everything else on
the network becomes horrible. Video calls stutter, web pages crawl, games become
unplayable — even though the download itself is going fine.

**The cause:** somewhere at the bottleneck (usually your home router or modem)
there is a **large, dumb buffer**. TCP's congestion control (Chapter 25) increases
speed until it sees loss — but with a huge buffer, packets aren't dropped, they're
*queued*. So TCP keeps speeding up and fills the buffer, and now **every** packet —
including your video call's, and your DNS lookups — has to wait behind a megabyte
of bulk download.

```
   Idle:         ping = 20 ms
   During a big upload:  ping = 800 ms      <-- bufferbloat
```

**How to test it (do this now — it takes 2 minutes):**

```bash
# Terminal 1: watch latency continuously
ping 1.1.1.1

# Terminal 2: saturate your connection
curl -o /dev/null https://speed.cloudflare.com/__down?bytes=200000000
# and for the upload direction, use a speed test site
```

Watch the ping times in terminal 1. **If they jump from 20 ms to hundreds of
milliseconds, you have bufferbloat.** Most home connections do.

Or use a browser test: **waveform.com/tools/bufferbloat** — it measures loaded vs
unloaded latency and gives you a grade.

**How to fix it:** use **Smart Queue Management (SQM)** — a modern queueing
algorithm (`fq_codel` or `cake`) at the bottleneck, which keeps the queue short
and shares it fairly between flows.

- Many consumer routers have an SQM or "bufferbloat" setting — turn it on.
- OpenWrt has excellent SQM support.
- On a Linux router: `tc qdisc replace dev <wan> root cake bandwidth 90Mbit`
  (set slightly *below* your real line rate, so the queue forms in the smart
  device rather than the dumb modem).

You typically give up ~5–10% of peak throughput and get a **10x improvement in
latency under load**. It's the single best network upgrade most people can make,
and it's free.

### Practice (25 min)

**A. Measure your own physics:**

```bash
ping -c 10 <your router>          # ~1 ms
ping -c 10 1.1.1.1                 # your ISP + a bit
ping -c 10 <something far away>    # distance
```

Divide the far one by 2 and multiply by 200 km/ms to estimate the fibre distance.
Compare with a map. It's usually surprisingly close.

**B. Do the bufferbloat test above.** Record your idle and loaded latency in
`notes.md`.

**C. Prove bandwidth ≠ speed for web pages:**

```bash
curl -w 'dns %{time_namelookup}  tcp %{time_connect}  tls %{time_appconnect}  ttfb %{time_starttransfer}  total %{time_total}\n' \
     -o /dev/null -s https://example.com
```

Notice how much of `total` happens *before* any content arrives. That portion is
pure latency and is unaffected by your bandwidth.

**D. Check your current queue discipline (Linux):**

```bash
tc qdisc show
```
If it says `pfifo_fast`, you have a dumb FIFO queue. `fq_codel` or `cake` is what
you want.

### Common confusions

- **"I upgraded to gigabit but browsing feels the same."** Correct, and expected.
  Browsing is latency-bound. You'd notice on large downloads.
- **"Bufferbloat is my ISP's fault."** Usually the buffer is in *your* router or
  modem, which means you can fix it.
- **"More bandwidth fixes congestion."** Only until it's used. Bufferbloat happens
  at any speed; it's about queue management, not capacity.

### Check yourself

1. What's the difference between bandwidth and latency?
2. Why doesn't upgrading your connection speed make websites load much faster?
3. What is bufferbloat and how do you test for it?
4. Why does the fix involve deliberately shaping to *below* your line rate?

### Further reading

- **Test:** `waveform.com/tools/bufferbloat` — run it now.
- **Site:** `bufferbloat.net` — the project that identified and named the problem.
- **Video:** "Bufferbloat" — search for Jim Gettys' talks; he discovered it.
- **Article:** "Latency Is Everywhere and It Costs You Sales" — a classic post on
  why latency matters commercially.
- **Book:** *High Performance Browser Networking* — Ilya Grigorik. **Free online**
  at hpbn.co. Chapters 1–2 on latency and TCP are perfect follow-ups to Part 6.
  Highly recommended.

---

### End of Part 6 — Milestone check

- [ ] I can explain the sliding window and the `window / RTT` formula
- [ ] I know why window scaling matters on long links
- [ ] I can describe slow start and the congestion sawtooth
- [ ] I know why parallel connections are faster on a lossy long path
- [ ] I recognise the 40 ms fingerprint
- [ ] **I have tested my own connection for bufferbloat**

---

# Part 7 — The application layer

The top of the stack: what the bytes actually *mean*.

## Chapter 28 — HTTP: how the web actually talks

### In one sentence

HTTP is a simple text conversation — the client sends a request describing what it
wants, the server sends back a status code and the content — and it runs on top of
the TCP connection from Part 5.

### The problem

You have a reliable byte stream to a server (thanks, TCP). Now you need an agreed
way to say "give me the home page" and to reply "here it is" or "that doesn't
exist".

### The idea, in plain language

HTTP is **text you could type by hand**. That's not an exaggeration — you'll do
exactly that in the Practice section, and it's the fastest way to understand it.

### A complete request and response

```
   CLIENT SENDS:
   ------------------------------------------------
   GET /index.html HTTP/1.1
   Host: example.com
   User-Agent: curl/8.4.0
   Accept: */*
   <blank line>
   ------------------------------------------------
     ^      ^          ^
   method  path     version

   SERVER REPLIES:
   ------------------------------------------------
   HTTP/1.1 200 OK
   Content-Type: text/html; charset=UTF-8
   Content-Length: 1256
   Cache-Control: max-age=604800
   <blank line>
   <!doctype html>
   <html>...
   ------------------------------------------------
       ^     ^   ^
   version code reason
```

Three parts to each: a **first line**, some **headers** (`Name: value`), a
**blank line**, then optionally a **body**. That's the whole protocol.

### The methods

| Method | Means | Changes data? | Safe to retry? |
|---|---|---|---|
| **GET** | Give me this | No | Yes |
| **POST** | Here's some data, do something with it | Yes | **No** — may create duplicates |
| **PUT** | Store this at this exact location | Yes | Yes (same result if repeated) |
| **PATCH** | Partially update this | Yes | Not necessarily |
| **DELETE** | Remove this | Yes | Yes (already-gone is still gone) |
| **HEAD** | Like GET but headers only, no body | No | Yes |
| **OPTIONS** | What can I do here? | No | Yes |

### The status codes

You only need the shape, plus a handful of specific ones:

```
   1xx  Information       (rare)
   2xx  SUCCESS           200 OK, 201 Created, 204 No Content
   3xx  REDIRECT          301 Moved Permanently, 302 Found, 304 Not Modified
   4xx  YOU messed up     400 Bad Request, 401 Unauthorized, 403 Forbidden,
                          404 Not Found, 429 Too Many Requests
   5xx  THE SERVER messed up  500 Internal Error, 502 Bad Gateway,
                              503 Unavailable, 504 Gateway Timeout
```

**The 4xx/5xx split is the most useful thing here:** 4xx means *your request was
wrong*; 5xx means *the server broke*. That immediately tells you who has to fix
it.

Two worth extra attention:
- **502 / 504** almost always mean a **proxy or load balancer** couldn't reach or
  didn't get a timely answer from the *real* server behind it (Chapter 36).
- **304 Not Modified** is caching working correctly — the server says "you already
  have the current version".

### Headers that matter

| Header | Why it matters |
|---|---|
| `Host:` | **Required in HTTP/1.1.** Lets one IP address serve many websites — the server uses it to decide which site you want. |
| `Content-Type:` | What the body is (`text/html`, `application/json`, `image/png`) |
| `Content-Length:` | How many bytes the body is, so the receiver knows when it's done |
| `Cache-Control:` / `ETag:` | How long this may be cached; how to check if it changed |
| `Connection: keep-alive` | Reuse this TCP connection for the next request (**huge** — see below) |
| `Cookie:` / `Set-Cookie:` | How state is maintained across a stateless protocol |
| `Authorization:` | Credentials |
| `Accept-Encoding: gzip` | "I can handle compressed responses" |

### The single biggest performance idea: connection reuse

HTTP is *stateless* — each request stands alone. But the **TCP connection** under
it doesn't have to be thrown away.

```
   WITHOUT keep-alive (old HTTP/1.0):
     open TCP (1 RTT) -> TLS (1 RTT) -> request -> response -> CLOSE
     ... repeat all of that for every image, stylesheet, script.

   WITH keep-alive (HTTP/1.1 default):
     open TCP (1 RTT) -> TLS (1 RTT) -> request -> response
                                     -> request -> response
                                     -> request -> response ...
```

Given the latency budget from Chapter 4, this is the difference between a page
loading in 400 ms and 4 seconds. **Every modern client does this by default** —
and it's why an API client that opens a fresh connection per call is so slow.

### Practice (30 min) — speak HTTP by hand

**A. Type HTTP yourself.** This makes the whole chapter click:

```bash
# Connect to a web server manually (unencrypted port 80)
nc example.com 80
```

Then **type exactly this**, including the blank line at the end:

```
GET / HTTP/1.1
Host: example.com
Connection: close

```

(Press Enter twice.) The server replies with headers and HTML. **You just were a
web browser.**

**B. The encrypted version** (port 443 needs TLS, so use `openssl`):

```bash
openssl s_client -connect example.com:443 -servername example.com -quiet
# then type the same GET request
```

**C. Use curl properly** — `-v` shows you everything:

```bash
curl -v https://example.com 2>&1 | head -40
```

Lines starting with `>` are what you sent; `<` is what came back; `*` is curl's
own commentary (DNS, TCP, TLS).

**D. Explore status codes:**

```bash
curl -sI https://example.com                    # HEAD -- headers only
curl -sI https://example.com/nothing-here       # 404
curl -sI -L https://google.com                  # -L follows redirects; watch the 301
curl -s -o /dev/null -w '%{http_code}\n' https://example.com
```

**E. See keep-alive save round trips:**

```bash
# two requests on ONE connection
curl -s -o /dev/null -w 'req1 %{time_total}\n' \
     -o /dev/null -w 'req2 %{time_total}\n' \
     https://example.com https://example.com
```
The second is much faster — no TCP or TLS handshake needed.

**F. Watch HTTP on the wire** (use plain HTTP so it's readable):

```bash
sudo tcpdump -i any -n -A 'tcp port 80 and host example.com'
# in another terminal:
curl -s http://example.com > /dev/null
```
`-A` prints the payload as text. You'll see the request and response in plain
sight — a good reminder of why HTTPS exists.

### Build it in Go (35 min) — HTTP by hand, over TCP and TLS

`nc` let you type HTTP. This program does the same, then solves the problem
`nc` hides from you: **where does the response body end?** HTTP/1.1 has three
answers (`Content-Length`, `Transfer-Encoding: chunked`, or "when the server
closes"), and every HTTP client has to implement all three.

```go
// rawhttp: speak HTTP/1.1 by hand over a TCP (or TLS) connection and decode
// the response framing yourself -- Content-Length vs chunked.
//
//	go run ./rawhttp http://neverssl.com/
//	go run ./rawhttp https://example.com/
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: rawhttp URL")
		os.Exit(2)
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	port := u.Port()
	if port == "" {
		port = map[string]string{"http": "80", "https": "443"}[u.Scheme]
	}
	addr := net.JoinHostPort(u.Hostname(), port)

	// Layer 4: a plain TCP connection...
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	// ...wrapped in TLS for https. Everything after this line is identical.
	if u.Scheme == "https" {
		tc := tls.Client(conn, &tls.Config{ServerName: u.Hostname(), NextProtos: []string{"http/1.1"}})
		if err := tc.Handshake(); err != nil {
			fmt.Println("TLS:", err)
			os.Exit(1)
		}
		conn = tc
	}
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// Layer 7: an HTTP request is just text with CRLF line endings.
	path := u.RequestURI()
	req := "GET " + path + " HTTP/1.1\r\n" +
		"Host: " + u.Host + "\r\n" + // mandatory in HTTP/1.1: one IP, many sites
		"User-Agent: rawhttp-lab/1.0\r\n" +
		"Accept-Encoding: identity\r\n" + // ask for no compression so we can read it
		"Connection: close\r\n" +
		"\r\n" // the empty line ends the headers
	fmt.Printf(">>> request (%d bytes)\n%s", len(req), req)
	io.WriteString(conn, req)

	r := bufio.NewReader(conn)
	status, _ := r.ReadString('\n')
	fmt.Printf("<<< %s", status)
	headers := map[string]string{}
	for {
		line, err := r.ReadString('\n')
		if err != nil || line == "\r\n" {
			break
		}
		fmt.Print("    ", line)
		k, v, _ := strings.Cut(strings.TrimRight(line, "\r\n"), ":")
		headers[strings.ToLower(k)] = strings.TrimSpace(v)
	}

	// Framing: how do we know where the body ends? (Chapter 39)
	var body []byte
	switch {
	case strings.EqualFold(headers["transfer-encoding"], "chunked"):
		body, err = readChunked(r)
		fmt.Printf("\nbody framing: chunked -> %d bytes after de-chunking\n", len(body))
	case headers["content-length"] != "":
		n, _ := strconv.Atoi(headers["content-length"])
		body = make([]byte, n)
		_, err = io.ReadFull(r, body)
		fmt.Printf("\nbody framing: Content-Length: %d\n", n)
	default:
		body, err = io.ReadAll(r) // only valid because we sent Connection: close
		fmt.Printf("\nbody framing: read until the server closed (%d bytes)\n", len(body))
	}
	if err != nil {
		fmt.Println("body error:", err)
	}
	fmt.Printf("first 120 bytes: %q\n", body[:min(120, len(body))])
}

// readChunked decodes "<hex size>\r\n<data>\r\n ... 0\r\n\r\n".
func readChunked(r *bufio.Reader) ([]byte, error) {
	var out []byte
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return out, err
		}
		sizeStr, _, _ := strings.Cut(strings.TrimSpace(line), ";") // ignore chunk extensions
		size, err := strconv.ParseInt(sizeStr, 16, 64)
		if err != nil {
			return out, fmt.Errorf("bad chunk size %q", sizeStr)
		}
		if size == 0 {
			return out, nil // (trailers would follow; we ignore them)
		}
		chunk := make([]byte, size)
		if _, err := io.ReadFull(r, chunk); err != nil {
			return out, err
		}
		out = append(out, chunk...)
		r.Discard(2) // the CRLF after each chunk
	}
}
```

```text
$ go run ./rawhttp http://example.com/
>>> request (112 bytes)
GET / HTTP/1.1
Host: example.com
User-Agent: rawhttp-lab/1.0
Accept-Encoding: identity
Connection: close

<<< HTTP/1.1 200 OK
    Content-Type: text/html; charset=utf-8
    Transfer-Encoding: chunked
    Server: cloudflare
    Age: 13523
    cf-cache-status: HIT
    alt-svc: h3=":443"; ma=86400

body framing: chunked -> 577 bytes after de-chunking
```

**What to notice:**

- **Going from HTTP to HTTPS is a four-line change.** `tls.Client(conn, ...)`
  wraps the TCP connection, and every line after it is unchanged, because a
  `*tls.Conn` is also a `net.Conn`. That's the layering of Chapter 3,
  expressed as an interface.
- **The response headers tell a story.** `cf-cache-status: HIT` and `Age: 13523`
  mean a CDN answered from cache, a response about 3.75 hours old (Chapter 40).
  `alt-svc: h3=":443"` advertises HTTP/3 for next time (Chapter 30).
- **Chunked framing** exists because a server that streams a response doesn't
  know its length in advance. Each chunk is `<hex length>\r\n<data>\r\n`, and
  a zero-length chunk ends the body.

**Exercises:**

1. Remove `Connection: close` and send **two requests** on the same
   connection, reading the first response fully (using its framing) before
   writing the second. You've just implemented keep-alive, the
   "single biggest performance idea" above.
2. Send a `HEAD` request and keep the `Content-Length` branch. Why does it
   hang? (HEAD responses carry a `Content-Length` but no body. Clients must
   know the *method* to parse the response.)
3. Compare with the real thing: the Go plan's
   [raw HTTP/1.1 server (Day 61)](../../Golang/detailed-90-day-plan/week9.md#day-61-http-1-1-from-scratch) builds the server side,
   and the [HTTPS guide's HTTP chapter](../../v2-https/real-life-guide-v1.md#chapter-5-http-the-conversation) covers
   HTTP/2 and HTTP/3 framing.

### Common confusions

- **"HTTP vs HTTPS?"** Same protocol; HTTPS runs it inside TLS encryption
  (Chapter 29). Port 80 vs 443 by convention.
- **"Is HTTP stateless?"** The protocol is; each request is independent. Cookies
  and tokens are how applications *simulate* state on top.
- **"Why do I need the `Host:` header?"** Because one server IP hosts hundreds of
  sites. Without it, the server doesn't know which one you want. It's mandatory in
  HTTP/1.1.
- **"GET vs POST for an API?"** GET for reading (safe, cacheable, repeatable),
  POST for actions that change things. Putting a state-changing operation in a GET
  causes real bugs when a browser or crawler prefetches it.

### Check yourself

1. What are the three parts of an HTTP request?
2. What's the difference between a 4xx and a 5xx status code?
3. Why is the `Host:` header required?
4. What does keep-alive save, and why does it matter so much?

### Further reading

- **Book (free online, superb):** *High Performance Browser Networking* — Ilya
  Grigorik, hpbn.co. Chapters on HTTP/1.1, HTTP/2, and TLS. **The best free
  resource for Part 7.**
- **Docs:** MDN Web Docs — "HTTP" section. The reference for methods, headers, and
  status codes. Bookmark it.
- **Video:** "HTTP Explained" — search Practical Networking or Hussein Nasser
  (whose whole channel on network engineering for developers is excellent).
- **Tool:** your browser's DevTools → Network tab. You already have the best HTTP
  inspector ever built, installed.

---

## Chapter 29 — TLS: how it gets encrypted

### In one sentence

TLS runs *inside* the TCP connection: the two sides negotiate keys, the server
proves its identity with a certificate, and everything after that is encrypted —
which is what the padlock means.

### The problem

You proved it yourself in Chapter 28's practice: plain HTTP is readable by anyone
on the path. You need three things:

1. **Confidentiality** — nobody can read it
2. **Integrity** — nobody can change it undetected
3. **Authentication** — you're really talking to `example.com`, not an impostor

The third is the hardest and the most important. Encryption without authentication
is useless — you'd have a perfectly secure channel to an attacker.

### Where TLS sits

```
      HTTP  (your request)
        |
      TLS   (encrypts it)          <-- this chapter
        |
      TCP   (reliable delivery)    <-- Part 5
        |
      IP    (routing)              <-- Part 3
```

TLS doesn't change TCP at all. It's just the first bytes the application sends
after the TCP handshake completes.

### How the handshake works (TLS 1.3)

```
   CLIENT                                              SERVER
     |                                                    |
     |  1. ClientHello                                    |
     |     "I support these ciphers.                       |
     |      Here's my half of a key.                       |
     |      I want to talk to example.com" (SNI)           |
     |     "I can speak HTTP/2 or HTTP/1.1"  (ALPN)        |
     |--------------------------------------------------->|
     |                                                    |
     |  2. ServerHello                                    |
     |     "Let's use this cipher. Here's my half."        |
     |     {Certificate}  <- proves it IS example.com      |
     |     {CertificateVerify}  <- proves it holds the key |
     |     {Finished}                                      |
     |<---------------------------------------------------|
     |     (everything in {} is ALREADY ENCRYPTED)         |
     |                                                    |
     |  3. {Finished}                                      |
     |     + your first encrypted HTTP request             |
     |--------------------------------------------------->|
```

**One round trip** in TLS 1.3 (TLS 1.2 needed two). On a repeat visit it can be
**zero** using a resumption ticket.

Three things worth noting:

- **The two "halves of a key"** combine into a shared secret that neither could
  compute alone, and which an eavesdropper cannot derive even having seen both
  halves. That's Diffie-Hellman key exchange. It also gives **forward secrecy**:
  the keys are fresh each time, so stealing the server's long-term key later
  doesn't decrypt recorded past sessions.
- **SNI** (the hostname you want) is sent **in plaintext** — it's how one IP serves
  many sites, and it's the main thing still visible to observers. (Encrypted
  ClientHello is being deployed to fix this.)
- **ALPN** negotiates HTTP/2 vs HTTP/1.1 during the handshake, saving a round trip.

### Certificates: how trust works

The certificate is the interesting part.

```
   A certificate says:  "the holder of THIS public key is example.com"
   ...signed by a Certificate Authority (CA).

   Your device ships with a list of ~150 trusted CAs (the "trust store").

   Chain of trust:
       example.com's certificate
            signed by ->  an intermediate CA
                              signed by -> a ROOT CA
                                              already trusted by your device
```

Your browser checks:
1. Does the chain lead to a root I trust?
2. Does the name match the site I asked for?
3. Is it within its validity dates?
4. Has it been revoked?

**Failing any of these is what produces the scary browser warning** — and you
should take it seriously, because a valid-looking site with a bad certificate is
exactly what an attack looks like.

### The versions

| Version | Status |
|---|---|
| SSL 2.0 / 3.0 | **Dead.** Broken. Never enable. |
| TLS 1.0 / 1.1 | **Deprecated** (formally, in 2021). Disable. |
| **TLS 1.2** | Fine if configured well. Still widely used. |
| **TLS 1.3** | Preferred. Faster (1 round trip), simpler, removed all the broken options. |

(Note: people still say "SSL" out of habit — e.g. "SSL certificate" — but the
protocol has been TLS since 1999.)

### Practice (30 min)

**A. Inspect a real handshake:**

```bash
openssl s_client -connect example.com:443 -servername example.com </dev/null 2>&1 | head -50
```

Look for: `Protocol`, `Cipher`, the certificate chain (`0 s:` is the site,
`1 s:` is the intermediate), and `Verify return code: 0 (ok)`.

**B. Read the certificate itself:**

```bash
echo | openssl s_client -connect example.com:443 -servername example.com 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates -ext subjectAltName
```

You'll see who it's for, who signed it, when it expires, and which hostnames it
covers.

**C. Measure the handshake cost:**

```bash
curl -w 'tcp: %{time_connect}s   tls: %{time_appconnect}s   (tls cost: computed below)\n' \
     -o /dev/null -s https://example.com
```
`time_appconnect − time_connect` is what TLS cost you.

**D. Watch the handshake in tcpdump.** You'll see the ClientHello and ServerHello
in the clear, then everything becomes opaque:

```bash
sudo tcpdump -i any -n 'tcp port 443 and host example.com' -c 20
curl -s https://example.com > /dev/null
```

**E. See a certificate failure safely:**

```bash
curl -v https://expired.badssl.com 2>&1 | grep -i -A3 'certificate'
curl -v https://wrong.host.badssl.com 2>&1 | grep -i 'subject\|SSL'
```
`badssl.com` is a site specifically built to demonstrate every kind of TLS
failure. Explore it in a browser too — very instructive.

**F. Check your own site's configuration** (if you have one): use
`ssllabs.com/ssltest` — a free, thorough grading tool.

### Build it in Go (30 min) — a scriptable `openssl s_client`

This program completes a real TLS handshake and prints what was negotiated:
version, cipher, ALPN, and the certificate chain with days to expiry. More
usefully, it turns the four classic certificate failures into **specific
diagnoses** instead of one generic "x509 error".

```go
// tlsinspect: connect, complete a TLS handshake, and print what was negotiated
// and who the server claims to be -- a scriptable `openssl s_client`.
//
//	go run ./tlsinspect example.com
//	go run ./tlsinspect -sni github.com example.com:443   # SNI the edge has no cert for
//	go run ./tlsinspect expired.badssl.com self-signed.badssl.com wrong.host.badssl.com
package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	sni := flag.String("sni", "", "override the SNI/hostname to verify (default: the host)")
	flag.Parse()
	exit := 0
	for _, target := range flag.Args() {
		if !strings.Contains(target, ":") {
			target += ":443"
		}
		if err := inspect(target, *sni); err != nil {
			fmt.Printf("%s\n  FAILED: %s\n\n", target, explain(err))
			exit = 1
		}
	}
	os.Exit(exit)
}

func inspect(target, sni string) error {
	host, _, _ := net.SplitHostPort(target)
	if sni == "" {
		sni = host
	}
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	start := time.Now()
	conn, err := tls.DialWithDialer(dialer, "tcp", target, &tls.Config{
		ServerName: sni,                        // sent as SNI AND used to verify the cert
		NextProtos: []string{"h2", "http/1.1"}, // ALPN: what we'd like to speak next
	})
	if err != nil {
		return err
	}
	defer conn.Close()
	st := conn.ConnectionState()

	fmt.Printf("%s  (TCP+TLS in %v)\n", target, time.Since(start).Round(time.Millisecond))
	fmt.Printf("  version  %s\n", tls.VersionName(st.Version))
	fmt.Printf("  cipher   %s\n", tls.CipherSuiteName(st.CipherSuite))
	fmt.Printf("  ALPN     %q\n", st.NegotiatedProtocol)
	fmt.Printf("  resumed  %v\n", st.DidResume)
	for i, c := range st.PeerCertificates {
		role := "leaf"
		if i > 0 {
			role = "intermediate"
		}
		days := int(time.Until(c.NotAfter).Hours() / 24)
		fmt.Printf("  [%d] %-12s CN=%s\n", i, role, c.Subject.CommonName)
		fmt.Printf("      issuer        %s\n", c.Issuer.CommonName)
		fmt.Printf("      valid until   %s  (%d days left)\n", c.NotAfter.Format("2006-01-02"), days)
		if i == 0 {
			fmt.Printf("      SANs          %s\n", strings.Join(c.DNSNames, ", "))
			if days < 21 {
				fmt.Println("      WARNING: renew soon")
			}
		}
	}
	fmt.Printf("  verified chain length: %d (ends at a root in YOUR trust store)\n\n", len(st.VerifiedChains[0]))
	return nil
}

// explain turns crypto/x509 errors into the fingerprints from Chapter 56.
func explain(err error) string {
	var hostErr x509.HostnameError
	var unknown x509.UnknownAuthorityError
	var invalid x509.CertificateInvalidError
	switch {
	case errors.As(err, &hostErr):
		return "NAME MISMATCH: certificate is valid for " + strings.Join(hostErr.Certificate.DNSNames, ", ")
	case errors.As(err, &unknown):
		return "UNKNOWN AUTHORITY: self-signed, private CA, or missing intermediate"
	case errors.As(err, &invalid) && invalid.Reason == x509.Expired:
		return "EXPIRED (or your clock is wrong): " + invalid.Detail
	case strings.Contains(err.Error(), "handshake failure"):
		return "HANDSHAKE REFUSED by server: often an SNI it has no certificate for (" + err.Error() + ")"
	default:
		return err.Error()
	}
}
```

```text
$ go run ./tlsinspect example.com
example.com:443  (TCP+TLS in 82ms)
  version  TLS 1.3
  cipher   TLS_AES_128_GCM_SHA256
  ALPN     "h2"
  [0] leaf         CN=example.com
      issuer        Cloudflare TLS Issuing ECC CA 3
      valid until   2026-12-25  (83 days left)
      SANs          example.com, *.example.com
  [1] intermediate CN=Cloudflare TLS Issuing ECC CA 3
  ...
$ go run ./tlsinspect expired.badssl.com self-signed.badssl.com wrong.host.badssl.com
expired.badssl.com:443
  FAILED: EXPIRED (or your clock is wrong): "*.badssl.com" certificate is expired
self-signed.badssl.com:443
  FAILED: UNKNOWN AUTHORITY: self-signed, private CA, or missing intermediate
wrong.host.badssl.com:443
  FAILED: NAME MISMATCH: certificate is valid for *.badssl.com, badssl.com
$ go run ./tlsinspect -sni github.com example.com:443
  FAILED: HANDSHAKE REFUSED by server: often an SNI it has no certificate for
```

**What to notice:**

- **`ServerName` does two jobs**: it's sent as **SNI** (so a shared server
  knows which certificate to present) *and* it's the name the certificate is
  checked against. The last example asks a Cloudflare edge for a name it
  doesn't host, and the edge refuses the handshake before sending any
  certificate.
- **A wildcard covers exactly one label.** `*.example.com` matches
  `wrong.example.com`, so `-sni wrong.example.com example.com` *passes*. It
  does not match `a.b.example.com` or `example.com` itself, which is why the
  SAN list above includes both.
- **ALPN `"h2"`**: the server agreed during the handshake to speak HTTP/2
  next, with no extra round trip (Chapter 30).
- **The verified chain ends at a root in *your* trust store**, the
  operating system's on macOS and Linux. That's why "works on my laptop,
  fails in the container" happens: minimal container images often have no
  CA bundle at all.

**Exercises:**

1. Add an `-expiry-days 30` flag that exits non-zero if any leaf expires
   within 30 days, and run it from cron against all your endpoints. That
   is a certificate-expiry monitor. Chapter 56 turns it into an operational
   practice.
2. Connect twice in one process with a shared `tls.Config` that has a
   `ClientSessionCache` (`tls.NewLRUClientSessionCache(10)`). Does
   `resumed` become `true` the second time? Not if you close straight
   after the handshake. In TLS 1.3 the server sends its session ticket
   *after* the handshake, so the client only receives it if it reads from
   the connection. Send a `HEAD` request and read the reply on the first
   connection, and the second one resumes. (Tested: `resumed=false` without
   the read, `resumed=true` with it.)
3. Go deeper with the Go plan: [TLS inspector + cert monitor (Day 38)](../../Golang/detailed-90-day-plan/week6.md#day-38-tls-handshake),
   [a full PKI generator (Day 40)](../../Golang/detailed-90-day-plan/week6.md#day-40-certificate-generation-in-go), and
   [mTLS with identity extraction (Day 41)](../../Golang/detailed-90-day-plan/week6.md#day-41-mutual-tls-mtls).

**Across the series:** the [HTTPS guide's TLS chapter](../../v2-https/real-life-guide-v1.md#chapter-4-tls-https-the-encrypted-tunnel)
explains the handshake from the browser's side, and the
[Go guide's TLS chapter](../../Golang/real-life-golang-guide.md#35-tls-in-go-crypto-tls-mutual-tls) covers `crypto/tls` server
configuration.

### Common confusions

- **"The padlock means the site is safe."** It means the *connection* is
  encrypted and the *name* matches. A phishing site can have a perfectly valid
  certificate for its own (misleading) domain.
- **"HTTPS hides everything."** It hides the URL path, headers, and content. It
  does **not** hide which server you're talking to (IP), or the hostname (SNI, for
  now), or the size and timing of your traffic.
- **"Do I need to buy a certificate?"** No — **Let's Encrypt** issues them free
  and automatically. There's no security difference for standard certificates.
- **"Certificate expired" outages** are extremely common and entirely preventable.
  Automate renewal and monitor expiry dates.

### Check yourself

1. What three guarantees does TLS provide?
2. Why is authentication the most important one?
3. What is SNI, and what does it leak?
4. Walk through what your browser checks in a certificate.

### Further reading

- **Book (free):** *High Performance Browser Networking*, Chapter 4 (Transport
  Layer Security). Excellent and readable.
- **Site:** `badssl.com` — every kind of TLS failure, safely demonstrated.
- **Tool:** `ssllabs.com/ssltest` — grade any HTTPS site and see exactly why.
- **Video:** "TLS Handshake Explained" — Practical Networking has a superb series
  on TLS and cryptography basics; also Computerphile's "Diffie Hellman" video for
  the key-exchange intuition (uses paint mixing — genuinely clarifying).
- **Article:** "A Complete Overview of SSL/TLS and its Cryptographic System" or
  Cloudflare's "What happens in a TLS handshake?"

---

## Chapter 30 — HTTP/2 and HTTP/3

### In one sentence

HTTP/2 sends many requests at once over a single TCP connection, and HTTP/3 moves
to a new transport over UDP to fix the one problem HTTP/2 couldn't.

### The problem HTTP/2 solved

HTTP/1.1 can only handle **one request at a time per connection**. A modern web
page needs 50–100 resources. So browsers opened ~6 connections per site and queued
everything — each connection paying its own TCP + TLS handshake, and each request
waiting for the one before it ("head-of-line blocking" at the HTTP layer).

### What HTTP/2 does

```
   HTTP/1.1                            HTTP/2
   6 connections, one request           ONE connection, many
   at a time each                       requests interleaved

   conn1: [req A......][req D...]       conn: [A][B][C][A][B][D][C][A]...
   conn2: [req B........]                     \_______________________/
   conn3: [req C.....]                         all mixed together,
   ...                                         reassembled at the other end
```

Key features:
- **Multiplexing** — many requests in flight on one connection, interleaved.
- **Header compression (HPACK)** — repeated headers (cookies, user-agent) cost
  ~1 byte instead of hundreds.
- **Binary framing** instead of text — faster to parse, less ambiguous.
- **Stream priorities** — tell the server the CSS matters more than the footer
  image.

Result: usually much faster, and it made "domain sharding" and sprite sheets
obsolete.

### The problem HTTP/2 couldn't solve

Everything is on **one TCP connection**. And TCP guarantees *in-order* delivery of
the whole byte stream. So:

```
   If ONE packet is lost, TCP holds back EVERYTHING behind it --
   including data belonging to completely unrelated requests.

   [req A data][req B data][LOST][req C data][req D data]
                             ^
                    everything after this waits,
                    even though C and D arrived fine
```

This is **TCP head-of-line blocking**. On a clean network HTTP/2 is great. On a
lossy one (mobile, congested Wi-Fi), it can be *worse* than HTTP/1.1's six
independent connections, because with six connections a loss only stalls one of
them.

### What HTTP/3 does

HTTP/3 runs over **QUIC**, a new transport protocol built on **UDP**.

Wait — UDP? The unreliable one? Yes, but QUIC implements its own reliability *per
stream*:

```
   QUIC streams are INDEPENDENT.
   A lost packet for stream B stalls ONLY stream B.
   Streams A, C, and D carry on.
```

Plus:

| Feature | Benefit |
|---|---|
| **Encryption and transport handshake merged** | 1 round trip to set up (or 0 on a repeat visit) instead of TCP+TLS's 2–3 |
| **Connection ID instead of the 4-tuple** | Your connection **survives changing networks** — walk out of Wi-Fi onto mobile and the download continues |
| **Always encrypted** | No unencrypted QUIC exists |
| **Implemented in the application, not the kernel** | Can be improved without OS upgrades — which is *why* it's built on UDP |

**Why build it on UDP?** Because middleboxes on the internet only really
understand TCP and UDP. Inventing a genuinely new protocol number would have been
blocked everywhere. UDP is the escape hatch that let a new transport be deployed.

### Which version are you using?

Browsers negotiate automatically: TLS's ALPN picks HTTP/2, and an `Alt-Svc`
header or DNS `HTTPS` record advertises HTTP/3 for next time. If UDP is blocked,
it falls back.

### Practice (20 min)

**A. Check what a site supports:**

```bash
curl -sI --http2 https://example.com | head -1        # "HTTP/2 200"
curl -sI --http1.1 https://example.com | head -1      # "HTTP/1.1 200"
curl -sI --http3 https://cloudflare.com | head -1     # needs curl built with HTTP/3
```

**B. Look for the HTTP/3 advertisement:**

```bash
curl -sI https://cloudflare.com | grep -i alt-svc
dig cloudflare.com HTTPS +short          # look for alpn="h3"
```

**C. See QUIC traffic on your machine.** You're probably already using it:

```bash
sudo tcpdump -i any -n 'udp port 443'
```
Browse a Google or Cloudflare-fronted site. That UDP traffic on 443 is HTTP/3.

**D. Compare in your browser.** DevTools → Network → right-click the column
headers → enable **Protocol**. Reload a few sites and see `h2`, `h3`, and
`http/1.1` in the wild.

**E. See multiplexing:** in DevTools, look at the waterfall for an HTTP/2 site.
Notice many requests starting at nearly the same moment, rather than in batches of
six.

### Common confusions

- **"HTTP/3 uses UDP so it's unreliable."** QUIC adds reliability on top, per
  stream. It's as reliable as TCP, with better loss behaviour.
- **"Should I upgrade my site to HTTP/2?"** If you use a CDN or a modern web
  server, you probably already have it. It's essentially free.
- **"Is HTTP/1.1 dead?"** No — it's simple, universal, and fine for APIs and
  server-to-server calls. It's browsers loading complex pages that benefit most
  from H2/H3.

### Check yourself

1. What was the main limitation of HTTP/1.1 that HTTP/2 fixed?
2. What is TCP head-of-line blocking, and why does it hurt HTTP/2 specifically?
3. Why is QUIC built on UDP rather than being a brand-new protocol?
4. What does "connection migration" mean and why is it useful?

### Further reading

- **Book (free):** *High Performance Browser Networking*, Chapters on HTTP/2 —
  hpbn.co. Still the clearest explanation.
- **Article:** "HTTP/3: From root to tip" and "The Road to QUIC" — Cloudflare
  blog. Excellent and readable.
- **Video:** "HTTP/1 vs HTTP/2 vs HTTP/3" — Hussein Nasser (he has a whole
  playlist on this) or ByteByteGo's animated version.
- **Spec:** RFC 9114 (HTTP/3) and RFC 9000 (QUIC) — for later.

---

### End of Part 7 — Milestone check

- [ ] **I have typed an HTTP request by hand with `nc`**
- [ ] I know the meaning of the 2xx/3xx/4xx/5xx families
- [ ] I can explain why keep-alive matters
- [ ] I can explain the three guarantees of TLS and how certificates work
- [ ] I have inspected a real certificate with `openssl`
- [ ] I can explain head-of-line blocking and how HTTP/3 avoids it

---

# Part 8 — Seeing and fixing the network

Networking is the most observable subject in computing: you can watch every byte.
This part turns you from someone who guesses into someone who looks.

## Chapter 31 — `tcpdump`: watching packets go by

### In one sentence

`tcpdump` prints every packet matching a filter, and once you can read its output
you can answer questions no log file can.

### The anatomy of a command

```
   sudo tcpdump -i any -n -c 20 'tcp port 443 and host example.com'
        \______/ \____/ \/ \___/ \_______________________________/
           |      which  no  stop            the FILTER
        needs   interface DNS  after         (in quotes)
        root              lookups 20
```

**Flags you'll use constantly:**

| Flag | Does |
|---|---|
| `-i any` | capture on all interfaces (or `-i eth0` for one) |
| `-n` | **don't resolve hostnames** — always use this; DNS lookups slow you down and pollute the capture |
| `-nn` | also don't resolve port names (shows `443` not `https`) |
| `-c 20` | stop after 20 packets |
| `-e` | show the Ethernet (MAC) header |
| `-A` | show the payload as text (great for HTTP) |
| `-X` | show payload as hex + text |
| `-w file.pcap` | **write to a file** for Wireshark |
| `-r file.pcap` | read a saved file |
| `-t` / `-tttt` | timestamp formats (`-tttt` is human-readable) |
| `-s 0` | capture the whole packet (default on modern versions) |

### The filters you actually need

```bash
host 192.168.1.50              # to or from this address
src host 192.168.1.50          # only from
dst host 192.168.1.50          # only to
net 192.168.1.0/24             # a whole network
port 443                       # to or from this port
portrange 8000-8100
tcp / udp / icmp / arp         # by protocol
tcp port 443 and host example.com          # combine with and / or / not
not port 22                    # exclude your own SSH session (do this on servers!)

# TCP flags -- extremely useful:
'tcp[tcpflags] & (tcp-syn) != 0'                # any SYN (connection attempts)
'tcp[tcpflags] & (tcp-syn|tcp-ack) == tcp-syn'  # SYNs only, not SYN-ACKs
'tcp[tcpflags] & (tcp-rst) != 0'                # resets (refused/aborted)
'tcp[tcpflags] & (tcp-fin) != 0'                # connection closes
```

### Reading the output

```
14:32:07.891234 IP 192.168.1.50.51000 > 93.184.216.34.443: Flags [S],
      seq 1234567, win 64240, options [mss 1460,sackOK,wscale 7], length 0
\____________/    \_____________/   \______________/  \___/  \_______/
   timestamp        source            destination      flags   details
                  IP.port             IP.port
```

**The flags shorthand** (memorise this — it's most of reading tcpdump):

| Symbol | Means |
|---|---|
| `[S]` | SYN — opening a connection |
| `[S.]` | SYN+ACK — accepting (the `.` means ACK) |
| `[.]` | ACK only — no data, just acknowledging |
| `[P.]` | PUSH+ACK — **actual data** |
| `[F.]` | FIN+ACK — closing politely |
| `[R]` / `[R.]` | RST — aborted / refused |

So a healthy connection reads: `[S]` `[S.]` `[.]` `[P.]` `[P.]` ... `[F.]` `[.]`

### The golden rules

1. **Always use `-n`.** Otherwise tcpdump does DNS lookups which generate more
   traffic which appears in your capture. Confusing and slow.
2. **Filter tightly.** On a busy machine, an unfiltered capture is unreadable and
   can drop packets.
3. **On a remote server, exclude your own SSH session**: `and not port 22`.
   Otherwise you capture yourself capturing, forever.
4. **For anything you'll analyse, write to a file** (`-w`) and open it in
   Wireshark. Terminal output is for quick looks.
5. **To prove where packets are lost, capture on both ends** simultaneously. A
   single capture can only tell you what arrived *there*.

### Practice (35 min) — the essential recipes

**A. A complete connection, start to finish:**

```bash
sudo tcpdump -i any -n -tttt 'host example.com' -c 30
# in another terminal:
curl -s https://example.com > /dev/null
```
Identify: the DNS query, the handshake, the TLS negotiation, data, and the close.

**B. All connection attempts leaving your machine:**

```bash
sudo tcpdump -i any -n 'tcp[tcpflags] & (tcp-syn|tcp-ack) == tcp-syn'
```
Leave it running. **You will be surprised how much your machine talks to.**

**C. Find who's being refused:**

```bash
sudo tcpdump -i any -n 'tcp[tcpflags] & tcp-rst != 0'
```

**D. Read HTTP in plain text:**

```bash
sudo tcpdump -i any -n -A 'tcp port 80' -c 40
curl -s http://example.com > /dev/null
```

**E. Save a capture for Wireshark:**

```bash
sudo tcpdump -i any -n -s0 -w /tmp/capture.pcap 'host example.com'
# ... generate traffic ...
# Ctrl-C, then:
wireshark /tmp/capture.pcap        # or open it in the GUI
```

**F. A rotating capture for an intermittent problem** (leave it running for hours):

```bash
sudo tcpdump -i any -n -s0 -C 100 -W 10 -w /tmp/ring.pcap 'port 443'
# keeps 10 files of 100 MB each, overwriting the oldest
```

### Common confusions

- **"tcpdump shows nothing."** Check: right interface (`-i any` is safest), filter
  syntax (needs quotes), and that traffic is actually happening.
- **"I see the request but not the response."** Often the response goes via a
  different interface. Use `-i any`.
- **"Checksum incorrect" warnings.** Normal — your network card computes checksums
  *after* tcpdump sees the packet. Ignore them for outgoing traffic.
- **"Can I capture on Wi-Fi?"** Yes for your own traffic. Capturing *other
  people's* Wi-Fi traffic needs monitor mode and is usually illegal.

### Check yourself

1. What does `-n` do and why should you always use it?
2. What do the flags `[S]`, `[S.]`, `[P.]`, and `[R]` mean?
3. Why capture on both ends to diagnose packet loss?
4. What filter shows only new connection attempts?

### Further reading

- **Cheat sheet:** search "tcpdump cheat sheet" — Daniel Miessler's is the classic
  and very good. Print it.
- **Man page:** `man pcap-filter` — the full filter syntax, better documented than
  you'd expect.
- **Blog:** Julia Evans' "tcpdump" zine page and her networking posts (jvns.ca).
- **Video:** "tcpdump Tutorial" — search; several good 20-minute walkthroughs.

---

## Chapter 32 — Wireshark: reading a conversation

### In one sentence

Wireshark is tcpdump with a screen: it decodes every protocol, reassembles
conversations, flags problems automatically, and draws graphs.

### When to use which

| Use `tcpdump` when | Use Wireshark when |
|---|---|
| You're on a server via SSH | You're on a desktop |
| You want a quick look | You want to *analyse* |
| You're capturing to a file | You're reading a file |
| Scripting/automation | Following a conversation, seeing statistics |

**The standard workflow:** capture with `tcpdump -w` on the server, copy the file,
open it in Wireshark on your laptop.

### The five things to do first with any capture

**1. `Statistics → Conversations` → TCP tab.**
Sort by bytes. This immediately tells you *who is talking the most*. Half of all
network investigations end here.

**2. `Statistics → Protocol Hierarchy`.**
What's actually on this wire, by percentage. Instantly spots "why is 60% of this
DNS?"

**3. Right-click a packet → `Follow → TCP Stream`.**
Reassembles the entire conversation into readable text. For HTTP, you see the
whole request and response. **This is Wireshark's best feature.**

**4. `Analyze → Expert Information`.**
Wireshark's automatic problem detector. It lists retransmissions, duplicate ACKs,
zero windows, resets, and malformed packets — categorised by severity. **Start
here when investigating a problem.**

**5. `View → Time Display Format → Seconds Since Previous Displayed Packet`.**
Turns the time column into "how long was the gap?" — which makes delays leap out.

### The display filters that find problems

Type these in the filter bar (note: **display** filters use a different syntax from
tcpdump's **capture** filters — this trips everyone up):

```
tcp.analysis.retransmission        # packets sent again
tcp.analysis.duplicate_ack         # receiver saying "I'm still missing something"
tcp.analysis.zero_window           # receiver's buffer is FULL (app too slow)
tcp.analysis.lost_segment          # a gap Wireshark noticed
tcp.analysis.flags                 # ALL of the above at once -- start here
tcp.flags.reset == 1               # who reset the connection, and when
http.response.code >= 400          # HTTP errors
dns.flags.rcode != 0               # failed DNS lookups
tcp.stream eq 3                    # isolate one conversation
tcp.time_delta > 0.2               # gaps longer than 200 ms
ip.addr == 192.168.1.50            # one host
```

### Reading the fingerprints

| What you see | What it means |
|---|---|
| Many retransmissions + duplicate ACKs | **Real packet loss** on the path |
| Retransmissions but *no* loss reported by the network | Possible **reordering** (check "out-of-order" too) |
| `ZeroWindow` from one side | That side's **application** is too slow to read. Not a network problem. |
| A long gap, then a PSH from the server, nothing missing | **Server-side application latency** — it was thinking |
| SYN, SYN, SYN with no reply | Firewall dropping, or host down |
| Handshake fine, then one big packet retransmitted forever | **MTU black hole** (Chapter 15) |
| RST immediately after SYN-ACK | Server accepted then aborted — often an overloaded backlog |
| ACKs arriving ~40 ms after data, tiny packets | **Nagle × delayed ACK** (Chapter 26) |
| Window never grows past ~64 KB on a long link | **Window scaling stripped** (Chapter 24) |

### The I/O Graph (your most powerful view)

`Statistics → I/O Graph`. Add these as separate lines on the same time axis:

- Throughput (bits/s)
- `tcp.analysis.retransmission` (COUNT)
- `tcp.analysis.bytes_in_flight` (AVG) — a proxy for the congestion window

Now you can literally *see* slow start ramping, the sawtooth on loss, and a
timeout (everything stops for 200 ms, then restarts from near zero).

### Practice (40 min)

**A. Capture and analyse your own web browsing:**

1. Start a capture in Wireshark on your main interface.
2. Load a website.
3. Stop the capture.
4. Do the five steps above in order.
5. Find the DNS query, the TCP handshake, the TLS handshake, and the data.
6. `Follow → TCP Stream` on the DNS query and on an HTTP request.

**B. Deliberately create problems and see their fingerprints** (Linux):

```bash
# 1. Packet loss
sudo tc qdisc add dev eth0 root netem loss 5%
#    ... capture a download ... then look at Expert Information
sudo tc qdisc del dev eth0 root

# 2. High latency
sudo tc qdisc add dev eth0 root netem delay 200ms
#    ... capture ... notice how slowly the window grows
sudo tc qdisc del dev eth0 root

# 3. Both, like a bad mobile connection
sudo tc qdisc add dev eth0 root netem delay 150ms 50ms loss 2%
sudo tc qdisc del dev eth0 root      # ALWAYS undo
```

**This exercise is worth an entire textbook chapter.** You are creating the exact
conditions you'll later have to diagnose, and seeing what they look like.

**C. Practise on sample captures.** The Wireshark wiki has a large library of
example captures (search "Wireshark SampleCaptures"). Open a few and try to work
out what happened before reading the description.

### Build it in Go (60 min) — decode a capture file yourself

Wireshark's dissectors look like magic until you write one. A `.pcap` file is
a 24-byte header followed by records, and each record is a 16-byte header plus
the raw bytes of one frame. Decoding a frame means peeling headers off in the
order of Chapter 3: link layer → IP → TCP/UDP. This decoder handles Ethernet
(with VLAN tags), Linux "cooked" captures from `-i any` (v1 and v2), macOS
loopback, IPv4, IPv6, TCP, UDP, and ICMP, and it flags likely
retransmissions.

```go
// pcapread: decode a tcpdump capture file yourself -- pcap header, Ethernet,
// IPv4/IPv6, TCP, UDP -- and print a tcpdump-like line per packet, plus a
// summary that flags handshakes, resets, and likely retransmissions.
//
//	sudo tcpdump -i any -w web.pcap -c 200 'tcp port 443'   # Linux "any" = SLL headers
//	sudo tcpdump -i en0 -w web.pcap -c 200 'tcp port 443'   # Ethernet headers
//	go run ./pcapread web.pcap
package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strings"
	"time"
)

const (
	linkEthernet = 1   // DLT_EN10MB
	linkRaw      = 101 // DLT_RAW: starts at the IP header
	linkSLL      = 113 // DLT_LINUX_SLL: tcpdump -i any on older Linux
	linkSLL2     = 276 // DLT_LINUX_SLL2: tcpdump -i any on current Linux
	linkNull     = 0   // DLT_NULL: macOS loopback (lo0)
)

type flowKey struct {
	src, dst netip.AddrPort
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: pcapread file.pcap")
		os.Exit(2)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)

	// Global header: 24 bytes. The magic number tells us byte order and
	// whether timestamps are micro- or nanoseconds.
	var gh [24]byte
	if _, err := io.ReadFull(r, gh[:]); err != nil {
		fmt.Println("not a pcap file:", err)
		os.Exit(1)
	}
	var bo binary.ByteOrder = binary.LittleEndian
	nano := false
	switch binary.LittleEndian.Uint32(gh[:]) {
	case 0xa1b2c3d4:
	case 0xa1b23c4d:
		nano = true
	case 0xd4c3b2a1:
		bo = binary.BigEndian
	case 0x4d3cb2a1:
		bo, nano = binary.BigEndian, true
	default:
		fmt.Println("unknown magic (pcapng? convert with: editcap -F pcap in.pcapng out.pcap)")
		os.Exit(1)
	}
	link := bo.Uint32(gh[20:])
	fmt.Printf("pcap v%d.%d  snaplen=%d  linktype=%d\n\n", bo.Uint16(gh[4:]), bo.Uint16(gh[6:]), bo.Uint32(gh[16:]), link)

	seen := map[flowKey]map[uint32]bool{} // per-direction sequence numbers carrying data
	var count, syns, rsts, retrans int
	var first time.Time
	for {
		var rh [16]byte
		if _, err := io.ReadFull(r, rh[:]); err != nil {
			break // io.EOF: done
		}
		sec, frac := bo.Uint32(rh[0:]), bo.Uint32(rh[4:])
		capLen, origLen := bo.Uint32(rh[8:]), bo.Uint32(rh[12:])
		data := make([]byte, capLen)
		if _, err := io.ReadFull(r, data); err != nil {
			fmt.Println("truncated record")
			break
		}
		ns := int64(frac) * 1000
		if nano {
			ns = int64(frac)
		}
		ts := time.Unix(int64(sec), ns)
		if first.IsZero() {
			first = ts
		}
		count++

		line, info, err := decode(link, data)
		if err != nil {
			fmt.Printf("%9.6f  [%d bytes] %v\n", ts.Sub(first).Seconds(), origLen, err)
			continue
		}
		if info.tcp {
			if info.flags&0x02 != 0 {
				syns++
			}
			if info.flags&0x04 != 0 {
				rsts++
			}
			if info.payload > 0 { // same direction, same seq, carrying data again = retransmission
				k := flowKey{info.src, info.dst}
				if seen[k] == nil {
					seen[k] = map[uint32]bool{}
				}
				if seen[k][info.seq] {
					retrans++
					line += "   <-- RETRANSMISSION?"
				}
				seen[k][info.seq] = true
			}
		}
		fmt.Printf("%9.6f  %s\n", ts.Sub(first).Seconds(), line)
	}
	fmt.Printf("\n%d packets, %d SYN, %d RST, %d likely retransmissions\n", count, syns, rsts, retrans)
}

type pktInfo struct {
	src, dst netip.AddrPort
	tcp      bool
	flags    byte
	seq      uint32
	payload  int
}

// decode peels the layers off one packet, outermost first.
func decode(link uint32, b []byte) (string, pktInfo, error) {
	var etherType uint16
	switch link {
	case linkEthernet: // dst MAC(6) src MAC(6) EtherType(2)
		if len(b) < 14 {
			return "", pktInfo{}, errors.New("short ethernet")
		}
		etherType, b = binary.BigEndian.Uint16(b[12:]), b[14:]
		if etherType == 0x8100 && len(b) >= 4 { // 802.1Q VLAN tag
			etherType, b = binary.BigEndian.Uint16(b[2:]), b[4:]
		}
	case linkSLL: // 16-byte "cooked" header, protocol in the last 2 bytes
		if len(b) < 16 {
			return "", pktInfo{}, errors.New("short SLL")
		}
		etherType, b = binary.BigEndian.Uint16(b[14:]), b[16:]
	case linkSLL2: // 20-byte "cooked v2" header, protocol in the FIRST 2 bytes
		if len(b) < 20 {
			return "", pktInfo{}, errors.New("short SLL2")
		}
		etherType, b = binary.BigEndian.Uint16(b[0:]), b[20:]
	case linkNull: // 4-byte host-order address family
		if len(b) < 4 {
			return "", pktInfo{}, errors.New("short null header")
		}
		etherType, b = 0x0800, b[4:]
		if b[0]>>4 == 6 {
			etherType = 0x86DD
		}
	case linkRaw:
		etherType = 0x0800
		if len(b) > 0 && b[0]>>4 == 6 {
			etherType = 0x86DD
		}
	default:
		return "", pktInfo{}, fmt.Errorf("linktype %d not supported", link)
	}

	var src, dst netip.Addr
	var proto byte
	var ttl byte
	switch etherType {
	case 0x0800: // IPv4: header length is in 32-bit words
		if len(b) < 20 {
			return "", pktInfo{}, errors.New("short IPv4")
		}
		ihl := int(b[0]&0x0F) * 4
		total := int(binary.BigEndian.Uint16(b[2:]))
		ttl, proto = b[8], b[9]
		src, dst = netip.AddrFrom4([4]byte(b[12:16])), netip.AddrFrom4([4]byte(b[16:20]))
		if total < len(b) {
			b = b[:total] // drop Ethernet padding
		}
		b = b[ihl:]
	case 0x86DD: // IPv6: fixed 40-byte header
		if len(b) < 40 {
			return "", pktInfo{}, errors.New("short IPv6")
		}
		proto, ttl = b[6], b[7]
		src, dst = netip.AddrFrom16([16]byte(b[8:24])), netip.AddrFrom16([16]byte(b[24:40]))
		b = b[40:]
	case 0x0806:
		return "ARP", pktInfo{}, nil
	default:
		return fmt.Sprintf("ethertype 0x%04x", etherType), pktInfo{}, nil
	}

	switch proto {
	case 6: // TCP
		if len(b) < 20 {
			return "", pktInfo{}, errors.New("short TCP")
		}
		sp, dp := binary.BigEndian.Uint16(b[0:]), binary.BigEndian.Uint16(b[2:])
		seq, ack := binary.BigEndian.Uint32(b[4:]), binary.BigEndian.Uint32(b[8:])
		off := int(b[12]>>4) * 4
		flags := b[13]
		win := binary.BigEndian.Uint16(b[14:])
		payload := max(0, len(b)-off)
		info := pktInfo{netip.AddrPortFrom(src, sp), netip.AddrPortFrom(dst, dp), true, flags, seq, payload}
		return fmt.Sprintf("TCP %s > %s [%s] seq %d ack %d win %d len %d ttl %d",
			info.src, info.dst, tcpFlags(flags), seq, ack, win, payload, ttl), info, nil
	case 17: // UDP
		if len(b) < 8 {
			return "", pktInfo{}, errors.New("short UDP")
		}
		sp, dp := binary.BigEndian.Uint16(b[0:]), binary.BigEndian.Uint16(b[2:])
		return fmt.Sprintf("UDP %s > %s len %d", netip.AddrPortFrom(src, sp), netip.AddrPortFrom(dst, dp), len(b)-8), pktInfo{}, nil
	case 1, 58:
		return fmt.Sprintf("ICMP %s > %s type %d", src, dst, b[0]), pktInfo{}, nil
	}
	return fmt.Sprintf("IP proto %d %s > %s", proto, src, dst), pktInfo{}, nil
}

// tcpFlags renders flags the way tcpdump does: S=SYN F=FIN R=RST P=PSH .=ACK
func tcpFlags(f byte) string {
	var s strings.Builder
	for _, x := range []struct {
		bit byte
		c   string
	}{{0x02, "S"}, {0x01, "F"}, {0x04, "R"}, {0x08, "P"}, {0x10, "."}} {
		if f&x.bit != 0 {
			s.WriteString(x.c)
		}
	}
	return s.String()
}
```

Capture something and decode it (this is a real capture of `curl https://example.com`
on Linux):

```text
$ sudo tcpdump -U -i any -w web.pcap 'tcp port 443' & curl -s -o /dev/null https://example.com; kill %1
$ go run ./pcapread web.pcap
pcap v2.4  snaplen=262144  linktype=276

 0.000000  TCP 172.17.0.2:46276 > 104.20.23.154:443 [S] seq 299092458 ack 0 win 65495 len 0 ttl 64
 0.023721  TCP 104.20.23.154:443 > 172.17.0.2:46276 [S.] seq 2742249108 ack 299092459 win 65408 len 0 ttl 63
 0.023741  TCP 172.17.0.2:46276 > 104.20.23.154:443 [.] seq 299092459 ack 2742249109 win 64 len 0 ttl 64
 0.025817  TCP 172.17.0.2:46276 > 104.20.23.154:443 [P.] seq 299092459 ack 2742249109 win 64 len 1571 ttl 64
 0.026065  TCP 104.20.23.154:443 > 172.17.0.2:46276 [.] seq 2742249109 ack 299094030 win 4083 len 0 ttl 63
...
21 packets, 2 SYN, 0 RST, 0 likely retransmissions
```

**Read it the way this chapter taught you:**

- Lines 1–3 are the **three-way handshake** from Chapter 21. Check the `+1`:
  the SYN-ACK acknowledges `299092459` = the client's ISN + 1.
- **RTT = 23.7 ms**, the gap between the SYN and the SYN-ACK. You just measured
  latency from a file.
- Line 4 has `len 1571`: the TLS **ClientHello**, which is bigger than you
  might expect because it carries key shares and extensions. Line 5's ACK
  number (`299094030` = `299092459 + 1571`) confirms the server received
  exactly that many bytes.
- **`win 64` is not 64 bytes.** It's the raw header field. Both sides
  negotiated window scaling in the SYNs (Chapter 24), so the real window is
  `64 << wscale`. Decoding the scale needs the SYN's TCP options, which is
  exercise 2.
- **`ttl 63`** on packets from the server: it set TTL 64, and one router (here
  the Docker bridge's NAT) decremented it on the way.

**Exercises:**

1. Capture a download with loss injected
   (`sudo tc qdisc add dev eth0 root netem loss 15%`). In testing, the
   decoder flagged the retransmitted segment:
   `... [P.] seq 1861632025 ... len 9904 ttl 63   <-- RETRANSMISSION?`.
   Check it against Wireshark's `tcp.analysis.retransmission` filter.
2. Parse the **TCP options** in SYN packets (kind 2 = MSS, 3 = window scale,
   4 = SACK permitted, 8 = timestamps) and print the real window size.
3. Add a per-connection summary keyed by the 4-tuple: bytes in each
   direction, duration, handshake RTT. That's Wireshark's
   *Statistics → Conversations* view.
4. Pull the **SNI** out of the ClientHello (the first payload byte is `0x16`
   for a TLS handshake; the SNI is extension type 0). This is how firewalls
   and L4 proxies route TLS without decrypting it. The Go plan's
   [SNI router (Day 55)](../../Golang/detailed-90-day-plan/week8.md#day-55-sni-routing) builds on exactly this.

### Common confusions

- **"Capture filter vs display filter?"** Capture filters (`tcp port 443`) decide
  what gets *recorded* — they use tcpdump syntax. Display filters
  (`tcp.port == 443`) decide what gets *shown* — different syntax, much more
  powerful. Capture broadly, filter narrowly.
- **"I can't read HTTPS."** Correct — it's encrypted. You can see the handshake,
  addresses, sizes, and timing, but not content. (You *can* decrypt your own
  browser's traffic by setting `SSLKEYLOGFILE` — search for it; a great exercise.)
- **"Wireshark says checksum errors everywhere."** Checksum offload again. Disable
  the checksum validation in preferences, or ignore it.

### Check yourself

1. Name the five things to do first with a new capture.
2. What does `tcp.analysis.zero_window` tell you about where the problem is?
3. What's the difference between a capture filter and a display filter?
4. Which menu shows you Wireshark's automatic problem detection?

### Further reading

- **Book (the standard):** *Practical Packet Analysis* — Chris Sanders. Written
  for beginners, uses real captures, and is genuinely enjoyable. **If you buy one
  book for Part 8, buy this.**
- **Book (deeper):** *Wireshark Network Analysis* — Laura Chappell. The official
  certification study guide; encyclopaedic.
- **Free:** the Wireshark User Guide and the Wireshark Wiki's SampleCaptures page.
- **Video:** Chris Greer's YouTube channel — packet analysis walkthroughs, very
  practical. Also the official Wireshark channel (SharkFest talks).
- **Blog:** Julia Evans, "How I use Wireshark" (jvns.ca).

---

## Chapter 33 — `ss`: what your machine's connections are doing

### In one sentence

`ss` (or `netstat` on macOS) lists every socket on your machine — what's
listening, what's connected, and, with the right flags, detailed TCP health for
each connection.

### The commands you'll use

```bash
ss -tuln          # LISTENING TCP and UDP sockets (t=tcp, u=udp, l=listening, n=numeric)
ss -tunp          # add the owning PROCESS (needs sudo for other users' sockets)
ss -tan           # ALL TCP sockets with states
ss -tn state established        # just the live connections
ss -s             # a summary count by state
ss -tin           # ** detailed TCP internals -- the good one **
ss -tan '( dport = :443 )'      # filter
```

macOS equivalents:
```bash
netstat -an | grep LISTEN
netstat -an -p tcp
sudo lsof -iTCP -sTCP:LISTEN -n -P        # with process names
```

### The most useful command: `ss -tin`

```
ESTAB  0  0   192.168.1.50:51000   93.184.216.34:443
     cubic wscale:7,7 rto:204 rtt:3.5/1.2 mss:1448 cwnd:32
     bytes_sent:5120 bytes_acked:5120 bytes_received:104832
     send 105.9Mbps  retrans:0/2  rcv_rtt:4  rcv_space:65535
```

| Field | Meaning |
|---|---|
| `cubic` | which congestion control algorithm (Chapter 25) |
| `wscale:7,7` | window scaling active (Chapter 24) — good |
| `rtt:3.5/1.2` | **measured round-trip time / variation**, in ms |
| `rto:204` | retransmission timeout (Chapter 22) |
| `cwnd:32` | **congestion window**, in packets (Chapter 25) |
| `retrans:0/2` | currently unrecovered / **total retransmits this connection** |
| `send 105.9Mbps` | the kernel's own estimate of this connection's rate |
| `mss:1448` | max segment size in use (Chapter 15) |

**This one command answers "why is this connection slow?"** more often than
anything else: check `rtt` (distance), `retrans` (loss), `cwnd` (congestion), and
`wscale` (is scaling on?).

### The state histogram: a 5-second health check

```bash
ss -tan | awk '{print $1}' | sort | uniq -c | sort -rn
```

```
   1247 ESTAB          <- normal
    340 TIME-WAIT      <- normal for a busy client
     18 LISTEN         <- your services
      4 SYN-SENT       <- trying to connect somewhere (is it failing?)
     92 CLOSE-WAIT     <- ** SUSPICIOUS: an app isn't closing sockets **
```

Refer back to Chapter 23 for what each state means. This histogram is the fastest
way to spot a socket leak.

### Practice (25 min)

**A. Audit what's listening on your machine:**

```bash
ss -tulnp            # Linux
sudo lsof -iTCP -sTCP:LISTEN -n -P   # macOS
```

For each entry, ask: *what is this, and should it be reachable from the network?*
Anything on `0.0.0.0` is exposed to your local network. **Do this on any server you
run.**

**B. Watch connections appear and disappear:**

```bash
watch -n1 'ss -tn state established | head -20'
```
Then browse the web and watch them come and go.

**C. Investigate one connection in depth:**

```bash
# start a big download
curl -o /dev/null https://speed.cloudflare.com/__down?bytes=100000000 &
# watch its internals
watch -n0.5 "ss -tin | grep -A2 cloudflare"
```
Watch `cwnd` grow (slow start!), `rtt` settle, and `retrans` hopefully stay at 0.

**D. Find the process using a port:**

```bash
sudo ss -tulnp | grep :8000       # Linux
sudo lsof -i :8000                 # macOS/Linux
```
This is the answer to "address already in use".

**E. Check the kernel's TCP counters:**

```bash
nstat -az | grep -Ei 'retrans|timeout|listendrop|listenoverflow'  # Linux
netstat -s -p tcp | head -30                                       # macOS
```

`ListenOverflows` climbing means a server's accept queue is full — a real and
often-missed cause of intermittent connection failures.

### Build it in Go (30 min, Linux) — `ss -ti` from inside your program

`ss -tin` reads its numbers from the kernel with `getsockopt(TCP_INFO)`. Your
own program can make the same call on its own sockets and log the kernel's
view of each connection: RTT, congestion window, retransmissions. That turns
"the network is slow" into numbers. This lab uploads 40 MB through Go's
HTTP client and samples `TCP_INFO` every 250 ms.

```go
//go:build linux

// tcpinfo: read the kernel's per-connection TCP state (the same data `ss -ti`
// shows) from inside your own program while it UPLOADS -- the sender is the
// side whose congestion window (cwnd) grows and shrinks.
//
//	go run ./tcpinfo https://speed.cloudflare.com/__up
//	# then add loss and run again:  sudo tc qdisc add dev eth0 root netem loss 2%
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: tcpinfo URL")
		os.Exit(2)
	}
	conns := make(chan *net.TCPConn, 1)
	tr := &http.Transport{
		ForceAttemptHTTP2: false,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := (&net.Dialer{}).DialContext(ctx, network, addr)
			if tc, ok := c.(*net.TCPConn); ok {
				conns <- tc // keep a handle on the raw TCP socket under the HTTP client
			}
			return c, err
		},
	}
	const size = 40 << 20 // 40 MB of zeros
	done := make(chan error, 1)
	go func() {
		resp, err := (&http.Client{Transport: tr}).Post(os.Args[1], "application/octet-stream",
			io.LimitReader(zeros{}, size))
		if err == nil {
			resp.Body.Close()
		}
		done <- err
	}()
	tc := <-conns

	fmt.Println("  time   rtt(ms) rttvar  cwnd(segs)  ssthresh  unacked  lost  total_retrans")
	start := time.Now()
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case err := <-done:
			fmt.Printf("done in %v (err=%v)\n", time.Since(start).Round(time.Millisecond), err)
			return
		case <-tick.C:
			ti, err := getTCPInfo(tc)
			if err != nil {
				fmt.Println("getsockopt:", err)
				return
			}
			ssthresh := fmt.Sprint(ti.Snd_ssthresh)
			if ti.Snd_ssthresh >= 1<<30 {
				ssthresh = "inf" // no loss seen yet: still in slow start
			}
			fmt.Printf("%6.2fs  %7.2f %6.2f  %10d  %8s  %7d  %4d  %13d\n",
				time.Since(start).Seconds(), float64(ti.Rtt)/1000, float64(ti.Rttvar)/1000,
				ti.Snd_cwnd, ssthresh, ti.Unacked, ti.Lost, ti.Total_retrans)
		}
	}
}

// getTCPInfo calls getsockopt(fd, IPPROTO_TCP, TCP_INFO). SyscallConn gives
// us the file descriptor without taking it away from Go's network poller.
func getTCPInfo(c *net.TCPConn) (*syscall.TCPInfo, error) {
	raw, err := c.SyscallConn()
	if err != nil {
		return nil, err
	}
	var info syscall.TCPInfo
	var sysErr error
	err = raw.Control(func(fd uintptr) {
		size := uint32(syscall.SizeofTCPInfo)
		_, _, e := syscall.Syscall6(syscall.SYS_GETSOCKOPT, fd,
			syscall.IPPROTO_TCP, syscall.TCP_INFO,
			uintptr(unsafe.Pointer(&info)), uintptr(unsafe.Pointer(&size)), 0)
		if e != 0 {
			sysErr = e
		}
	})
	if err != nil {
		return nil, err
	}
	return &info, sysErr
}

type zeros struct{}

func (zeros) Read(p []byte) (int, error) { clear(p); return len(p), nil }
```

```text
$ go run ./tcpinfo https://speed.cloudflare.com/__up
  time   rtt(ms) rttvar  cwnd(segs)  ssthresh  unacked  lost  total_retrans
  ...one row every 250 ms...

$ sudo tc qdisc add dev eth0 root netem loss 2%     # then run it again
```

Run it on a real Linux host or VM. Inside Docker Desktop on a Mac, the
connection is proxied by the VM's network layer, so RTT and cwnd describe a
hop to the proxy, not the path to Cloudflare. What to look for:

- **Without loss:** `ssthresh` shows `inf` at first and `cwnd` climbs quickly.
  That's slow start. `total_retrans` stays at 0.
- **With 2% loss:** `total_retrans` starts counting, `ssthresh` becomes a
  real number the first time loss is detected, and `cwnd` drops to around
  it and then grows more slowly. You're watching congestion control react.
- **`rtt` vs `rttvar`:** a large variance relative to RTT means jitter. In
  testing, a lossy run started with `rtt` 814 ms and `rttvar` 480 ms: the
  first samples included retransmission delays.

**What to notice:**

- **`SyscallConn().Control()`** hands you the raw file descriptor *without*
  taking it away from Go's runtime. It's the safe escape hatch for socket
  options and `getsockopt` calls that the `net` package doesn't wrap. (Calling
  `File()` instead duplicates the descriptor and puts it into blocking mode.)
- **The custom `DialContext`** is how you get hold of the TCP socket under a
  high-level HTTP client. The same hook is used for custom DNS, SOCKS proxies,
  and the SSRF protection in the HTTPS guide.
- **`ssthresh = inf`** means no loss has happened yet, so the connection is still in
  slow start (Chapter 25). After a loss, `ssthresh` is set to roughly half
  of `cwnd`, and the congestion window grows linearly from there.
- **Which side to measure:** `cwnd` describes the *sending* side. On a
  download your client's `cwnd` stays at 10 because it hardly sends
  anything. To study a slow download, measure on the server.

**Exercise:** make it a library. Write `func LogTCPInfo(c net.Conn, every time.Duration)`
and call it from a real service on its slowest outbound connection, emitting
`rtt`, `cwnd`, and `total_retrans` as structured log fields
([Go guide: `log/slog`](../../Golang/real-life-golang-guide.md#36-structured-logging-and-observability-log-slog-metrics-tracing)). When someone next says "the database
is slow", you'll know whether the network was involved.

### Common confusions

- **"`ss` vs `netstat`?"** `ss` is the modern Linux replacement — faster and shows
  more. `netstat` is deprecated on Linux but still standard on macOS/BSD.
- **"Why can't I see other users' processes?"** Use `sudo`.
- **"`Recv-Q` and `Send-Q` on a LISTEN socket?"** On a listening socket these mean
  something different: `Recv-Q` is how many completed connections are waiting to be
  accepted, and `Send-Q` is the maximum (the backlog). If `Recv-Q` is at the max,
  your application isn't accepting fast enough.

### Check yourself

1. What does `ss -tuln` show?
2. In `ss -tin`, what do `rtt`, `cwnd`, and `retrans` tell you?
3. What does a large number of `CLOSE-WAIT` sockets indicate?
4. How do you find which process is using port 8080?

### Further reading

- **Man pages:** `man ss` — genuinely good, with examples.
- **Article:** search "ss command examples" — many good cheat sheets.
- **Blog:** Brendan Gregg's Linux performance tools posts (brendangregg.com) —
  where `ss` fits in the wider observability toolkit.

---

## Chapter 34 — A troubleshooting method that works

### In one sentence

Work **up the layers** from the bottom, prove each one before moving on, and
change one thing at a time — this beats guessing every single time.

### The method

```
   START: "the network is broken"

   0. DEFINE IT.  What exactly fails? From where? To what? Since when?
      Does it fail EVERY time or sometimes? What changed recently?
      -> "sometimes" and "what changed" solve half of all problems

   1. LAYER 1-2: is the link up?
        ip -br link           (is it UP? right speed/duplex?)
        ip -s link            (errors/drops climbing?)
      -> if down: cable, port, Wi-Fi association, driver

   2. LAYER 3a: do I have a sane address?
        ip -br addr
      -> 169.254.x.x means DHCP FAILED (Chapter 17)
      -> no address at all: DHCP or config

   3. LAYER 3b: can I reach my own gateway?
        ip route ; ping <gateway>
      -> fails: local network problem, wrong mask, ARP failure

   4. LAYER 3c: can I reach the internet BY IP?
        ping 1.1.1.1
      -> fails but gateway works: routing, ISP, or firewall

   5. LAYER 4-7a: does DNS work?
        dig example.com ; dig @1.1.1.1 example.com
      -> ping by IP works but by name doesn't = DNS. Always test both.

   6. LAYER 4b: can I reach the SERVICE?
        nc -zv host port          (open? refused? timeout?)
      -> refused = nothing listening
      -> timeout = firewall or unreachable

   7. LAYER 7: does the APPLICATION work?
        curl -v https://host/path
      -> now you're looking at HTTP status codes, TLS, and app logs

   8. If still stuck: CAPTURE PACKETS (Chapters 31-32) on both ends.
```

### The single most useful diagnostic split

```
   ping <ip address>      works?
   ping <hostname>        fails?
         -> IT IS DNS.  (Chapter 18)
```

Run this **first**, always. It splits the problem space in half in five seconds.
There's a reason "it's always DNS" is the oldest joke in the industry.

### The principles

1. **Bottom-up, not guess-first.** Prove the link before blaming the app.
2. **Change one thing at a time**, and write down what you changed. Networking
   problems are notorious for two people "fixing" things simultaneously and
   confusing each other.
3. **Compare against something that works.** Another machine, another network,
   another user. The difference *is* the answer.
4. **"It works from here"** — establish *where* it works and where it doesn't. The
   boundary is where the problem lives.
5. **Timing is a clue.** Failures at exact intervals (30 s, 60 s, 5 min) are
   timeouts or leases, not randomness. Round numbers in latency (40 ms, 200 ms,
   1 s, 5 s) are timers (Chapter 26).
6. **Both ends.** For loss or latency, capture on both sides. One capture can only
   prove what arrived, not what was sent.
7. **Don't skip "what changed?"** Most breakage follows a change: a deploy, an
   update, a config edit, a certificate expiry, a new firewall rule.

### A one-command first look

Keep this handy — it answers most of steps 1–5 at once:

```bash
echo "=== LINK ===";    ip -br link 2>/dev/null || ifconfig | grep -E '^\w|status'
echo "=== ADDR ===";    ip -br addr 2>/dev/null || ifconfig | grep 'inet '
echo "=== ROUTE ===";   ip route 2>/dev/null || netstat -rn | head -5
echo "=== DNS CFG ==="; cat /etc/resolv.conf 2>/dev/null | grep -v '^#'
echo "=== GATEWAY ==="; ping -c2 -W2 $(ip route | awk '/default/{print $3; exit}') 2>/dev/null
echo "=== INTERNET ==="; ping -c2 -W2 1.1.1.1
echo "=== DNS ===";     dig +short +time=2 example.com
echo "=== HTTP ===";    curl -s -o /dev/null -w '%{http_code} in %{time_total}s\n' https://example.com
```

### Practice (30 min) — break things on purpose

**The single best way to learn troubleshooting is to cause the failures yourself,
so you recognise them later.** Do these on a spare machine or VM, and undo each
one.

```bash
# 1. Break DNS
sudo cp /etc/resolv.conf /tmp/resolv.backup
echo "nameserver 192.0.2.1" | sudo tee /etc/resolv.conf   # a black hole
# symptom: ping by IP works, by name doesn't
sudo cp /tmp/resolv.backup /etc/resolv.conf                # UNDO

# 2. Break the default route
ip route | grep default                                    # note it first!
sudo ip route del default
# symptom: local works, internet doesn't, "Network is unreachable"
sudo ip route add default via <your gateway>               # UNDO

# 3. Break the MTU
sudo ip link set eth0 mtu 900
# symptom: small things work, large downloads stall
sudo ip link set eth0 mtu 1500                             # UNDO

# 4. Simulate loss and latency
sudo tc qdisc add dev eth0 root netem loss 10% delay 200ms
sudo tc qdisc del dev eth0 root                            # UNDO

# 5. Block a port
sudo iptables -A OUTPUT -p tcp --dport 443 -j DROP
# symptom: TIMEOUT (not refused) -- feel the difference
sudo iptables -D OUTPUT -p tcp --dport 443 -j DROP         # UNDO
```

**For each one:** before undoing it, run the diagnostic script above and note
which step first shows the failure. Write it in `notes.md`. That table becomes
your personal troubleshooting reference.

### Build it in Go (35 min) — the troubleshooting ladder as one command

The method in this chapter is "climb the layers and stop at the first one
that fails". That's an algorithm, so here it is as a program. Each rung only
runs if the one below it passed, and the program tells you which layer to
investigate.

```go
// netcheck: Chapter 34's troubleshooting ladder as one program. It climbs the
// stack -- route/interface, DNS, TCP, TLS, HTTP -- stops at the first rung
// that fails, and tells you which layer to look at.
//
//	go run ./netcheck https://example.com/
//	go run ./netcheck https://example.com:81/      # watch it stop at TCP
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

type step struct {
	layer string
	run   func() (string, error)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: netcheck URL")
		os.Exit(2)
	}
	u, err := url.Parse(os.Args[1])
	if err != nil || u.Hostname() == "" {
		fmt.Println("bad URL:", os.Args[1])
		os.Exit(2)
	}
	host, port := u.Hostname(), u.Port()
	if port == "" {
		port = map[string]string{"http": "80", "https": "443"}[u.Scheme]
	}
	var ips []net.IP
	var addr string

	steps := []step{
		{"L3 route", func() (string, error) {
			c, err := net.Dial("udp", "1.1.1.1:53") // sends nothing; just asks the routing table
			if err != nil {
				return "", fmt.Errorf("no route to the internet: %w", err)
			}
			defer c.Close()
			src := c.LocalAddr().(*net.UDPAddr).IP
			// Which interface owns that source address? That's the one that matters.
			ifs, _ := net.Interfaces()
			for _, i := range ifs {
				as, _ := i.Addrs()
				for _, a := range as {
					if n, ok := a.(*net.IPNet); ok && n.IP.Equal(src) {
						return fmt.Sprintf("default route via %s (up=%v), source %s", i.Name, i.Flags&net.FlagUp != 0, src), nil
					}
				}
			}
			return "default route exists, source " + src.String(), nil
		}},
		{"L7 DNS", func() (string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			start := time.Now()
			addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return "", err
			}
			for _, a := range addrs {
				ips = append(ips, a.IP)
			}
			return fmt.Sprintf("%s -> %v in %v", host, ips, time.Since(start).Round(time.Millisecond)), nil
		}},
		{"L4 TCP", func() (string, error) {
			addr = net.JoinHostPort(ips[0].String(), port)
			start := time.Now()
			c, err := net.DialTimeout("tcp", addr, 3*time.Second)
			if err != nil {
				return "", err
			}
			c.Close()
			return fmt.Sprintf("connected to %s in %v (~1 RTT)", addr, time.Since(start).Round(time.Millisecond)), nil
		}},
		{"L5/6 TLS", func() (string, error) {
			if u.Scheme != "https" {
				return "skipped (plain http)", nil
			}
			start := time.Now()
			c, err := tls.DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "tcp", addr,
				&tls.Config{ServerName: host})
			if err != nil {
				return "", err
			}
			defer c.Close()
			st := c.ConnectionState()
			days := int(time.Until(st.PeerCertificates[0].NotAfter).Hours() / 24)
			return fmt.Sprintf("%s, cert valid %d more days, TCP+TLS %v",
				tls.VersionName(st.Version), days, time.Since(start).Round(time.Millisecond)), nil
		}},
		{"L7 HTTP", func() (string, error) {
			client := &http.Client{Timeout: 5 * time.Second,
				CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
			start := time.Now()
			resp, err := client.Get(u.String())
			if err != nil {
				return "", err
			}
			resp.Body.Close()
			if resp.StatusCode >= 500 {
				return "", fmt.Errorf("server error %s", resp.Status)
			}
			return fmt.Sprintf("%s in %v (server: %q)", resp.Status, time.Since(start).Round(time.Millisecond), resp.Header.Get("Server")), nil
		}},
	}

	for _, s := range steps {
		detail, err := s.run()
		if err != nil {
			fmt.Printf("FAIL  %-11s %v\n\n=> Everything below %s works. Investigate %s first.\n", s.layer, err, s.layer, s.layer)
			os.Exit(1)
		}
		fmt.Printf("ok    %-11s %s\n", s.layer, detail)
	}
	fmt.Println("\n=> every layer answered. If users still complain, look at latency, not reachability.")
}
```

```text
$ go run ./netcheck https://example.com/
ok    L3 route    default route via en0 (up=true), source 192.168.1.2
ok    L7 DNS      example.com -> [104.20.23.154 172.66.147.243 2606:4700:...] in 2ms
ok    L4 TCP      connected to 104.20.23.154:443 in 38ms (~1 RTT)
ok    L5/6 TLS    TLS 1.3, cert valid 83 more days, TCP+TLS 96ms
ok    L7 HTTP     200 OK in 514ms (server: "cloudflare")

$ go run ./netcheck https://example.com:81/
ok    L3 route    default route via en0 (up=true), source 192.168.1.2
ok    L7 DNS      example.com -> [...] in 3ms
FAIL  L4 TCP      dial tcp 172.66.147.243:81: i/o timeout

=> Everything below L4 TCP works. Investigate L4 TCP first.
```

**What to notice:**

- **The order isn't the OSI order.** DNS (an L7 protocol) runs before TCP,
  because you need the address before you can connect. The
  [OSI guide's "layers are not steps in time"](../real-life-example-osi.md#2-5-layers-are-not-steps-in-time) section explains
  why this is correct.
- **Each rung is the smallest test that can prove its layer.** TCP is checked
  with a bare connect, not an HTTP request. A failing HTTP request alone can't
  tell you whether DNS, TCP, TLS, or the app is at fault.
- **Timing on every rung.** "Everything works, but TCP took 900 ms" is also a
  finding: SYN loss and a retry (Chapter 21).

**Exercises:**

1. Break each layer in turn, using the "break things on purpose" practice
   above (a bad resolver, no default route, a blocked port), and check that
   the tool names the right layer each time.
2. Add an **MTU rung**: download something larger than 1,500 bytes and fail
   if the small request works but the large one stalls (Chapter 15's
   fingerprint).
3. Add `-json` output and run it from cron on three networks (home, office,
   cloud VM). You now have a minimal synthetic monitor. Chapter 46's
   `prober` lab turns that into proper metrics.

**Across the series:** the same ladder appears as `curl -w` timings in the
[OSI guide's debugging playbook](../real-life-example-osi.md#part-14-debugging-by-layer-the-incident-playbook) and in the
[HTTPS guide's debugging chapter](../../v2-https/real-life-guide-v1.md#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs). The
[Linux guide's Appendix C](../../os-linux/real-life-os-guide.md#appendix-c-troubleshooting-playbook-the-server-is-slow) is the host-side
equivalent for "the server is slow".

### Check yourself

1. What's the first thing you test, and why?
2. What does "ping by IP works, by name doesn't" tell you?
3. Why is "refused" a more informative failure than "timeout"?
4. Why capture on both ends?

### Further reading

- **Book:** *Network Troubleshooting Tools* — Joseph Sloan (O'Reilly), or the
  troubleshooting chapters of *Network Warrior*.
- **Article:** search "systematic network troubleshooting OSI bottom-up" — the
  approach is standard and well documented.
- **Practice:** set up two VMs or containers and deliberately misconfigure the
  network between them. Nothing teaches faster.

---

## Chapter 35 — The five classic problems and their fingerprints

A reference table. Come back to this whenever something breaks.

### 1. DNS failure

```
   SYMPTOMS:  ping <ip> works;  ping <name> fails
              "could not resolve host" / NXDOMAIN
              Works for some people, not others
              Worked yesterday, broke after a change
   CHECK:     dig <name>
              dig @1.1.1.1 <name>       (compare -- different answers = caching/propagation)
              cat /etc/resolv.conf
   CAUSES:    dead DNS server, wrong resolver, TTL not expired after a change,
              split-horizon/VPN, typo in a record
```

### 2. MTU / fragmentation black hole

```
   SYMPTOMS:  Connects, then HANGS
              SSH logs in then freezes
              Small pages load, big ones don't
              git clone stalls at a few KB
              Happens over VPN but not without
   CHECK:     ping -M do -s 1472 <host>       (fails? MTU is under 1500)
              tracepath <host>
   CAUSES:    VPN/tunnel overhead + a firewall blocking ICMP type 3
   FIX:       allow ICMP type 3; MSS clamping; lower the MTU
```

### 3. Firewall dropping (vs rejecting)

```
   SYMPTOMS:  Connection TIMES OUT (long wait, no error)
              vs "connection refused" (instant)
   CHECK:     nc -zv host port
              sudo tcpdump -n 'host X and port Y'   -- do you see SYNs going out
                                                       and nothing coming back?
   CAUSES:    firewall rule, security group, missing port forward,
              service bound to 127.0.0.1 only
   KEY:       TIMEOUT = silently dropped.  REFUSED = reached the host,
              nothing listening.
```

### 4. Packet loss

```
   SYMPTOMS:  Slow transfers despite a fast link
              Video calls stutter; voice breaks up
              Retransmissions in Wireshark
   CHECK:     mtr -rwzbc 100 <host>
              -> loss at hop N that CONTINUES to the destination = real
              -> loss at hop N that disappears after = cosmetic (ICMP rate limit)
              ip -s link          (interface errors?)
              ethtool -S eth0 | grep -i err
   CAUSES:    bad cable/connector, Wi-Fi interference, congested link,
              duplex mismatch, failing hardware
```

### 5. Latency / bufferbloat

```
   SYMPTOMS:  Everything is fine until someone downloads something
              Video calls break during backups
              Ping jumps from 20 ms to 500+ ms under load
   CHECK:     ping <somewhere> while running a big transfer
              waveform.com/tools/bufferbloat
   CAUSES:    a big dumb buffer at the bottleneck (usually your router/modem)
   FIX:       enable SQM / fq_codel / cake on the router
```

### The bonus one: it's the application

```
   SYMPTOMS:  Network metrics are all fine.  No loss, low latency,
              handshakes complete, TLS fine.  But it's still slow.
   CHECK:     curl -w '... ttfb %{time_starttransfer} total %{time_total}'
              -> a big gap between connect and first byte = the SERVER is thinking
              Wireshark: a long gap with no missing packets = application delay
              ss -tin: ZeroWindow = the RECEIVER's app is too slow
   REMEMBER:  "The network is slow" is usually a database query.
```

### The quick reference table

| Symptom | First suspect | First command |
|---|---|---|
| Name doesn't resolve | DNS | `dig <name>` |
| Connects then hangs on big transfers | MTU | `ping -M do -s 1472 <host>` |
| Times out | Firewall dropping | `nc -zv host port` |
| Refused instantly | Nothing listening | `ss -tulnp` on the server |
| Slow bulk transfer | Loss or window | `mtr`, `ss -tin` |
| Everything slow *during* downloads | Bufferbloat | ping while transferring |
| Latency floor at exactly 40/200 ms | A timer (Nagle, RTO) | Chapter 26 |
| 502 / 504 errors | Backend unreachable from the proxy | check the proxy's upstream |
| Works for some users only | DNS, routing, or MTU | compare a working and failing client |
| Broke after a change | The change | `what changed?` |

### Practice

Take the five deliberate breakages from Chapter 34, and for each one write down,
in `notes.md`:
1. The exact symptom a user would report
2. The first command you'd run
3. What that command showed

**That document is your troubleshooting playbook**, and it's worth more than any
chapter of this guide because you wrote it.

### Further reading

- **Article:** search "network troubleshooting flowchart" for visual versions of
  Chapter 34's method.
- **Community:** r/networking's wiki, and the Server Fault questions tagged
  `networking` — reading real diagnoses is excellent practice.

---

### End of Part 8 — Milestone check

- [ ] I can write a tcpdump command with a filter from memory
- [ ] I can read tcpdump flags (`[S]`, `[S.]`, `[P.]`, `[R]`)
- [ ] **I have captured traffic and analysed it in Wireshark**
- [ ] I know the five first steps for any capture
- [ ] I can read `ss -tin` and say why a connection is slow
- [ ] **I have deliberately broken my network 5 ways and diagnosed each**
- [ ] I have written my own troubleshooting playbook

---

# Part 9 — Modern networking

A tour of what sits on top of everything you've learned. These are overviews, not
deep dives — each links to where to go deeper.

## Chapter 36 — Load balancers, CDNs, and anycast

### The problem

One server can't handle a popular website, and even if it could, it would be
slow for people on the other side of the world (Chapter 27: physics).

### Load balancers

A **load balancer** sits in front of several servers and spreads requests between
them.

```
                    +--> server 1
   users --> [ LB ] +--> server 2
                    +--> server 3
```

Two kinds, and the distinction matters:

| | **L4 load balancer** | **L7 load balancer** |
|---|---|---|
| Decides using | IP address + port only | The actual HTTP request |
| Can route by | Nothing about the content | URL path, `Host` header, cookies |
| Sees your data? | No — traffic passes through encrypted | Yes — it terminates TLS and re-encrypts |
| Speed | Very fast | Slower, but far more capable |
| Examples | AWS NLB, IPVS, HAProxy in TCP mode | nginx, HAProxy in HTTP mode, Envoy, AWS ALB |

**Health checks** are the important operational piece: the load balancer
continuously tests each server (`GET /healthz`) and stops sending traffic to ones
that fail. Getting this right (interval, threshold, what "healthy" means) is most
of running a reliable service.

**Why you get 502/504 errors:** the load balancer couldn't reach the server
behind it (502) or waited too long for a reply (504). **The error comes from the
load balancer, not the server** — which is why your application logs may show
nothing.

### CDNs

A **Content Delivery Network** puts copies of your content in hundreds of
locations worldwide, so users connect to something nearby.

Recall Chapter 4's latency budget: 4 round trips before the first byte. If the
server is 90 ms away that's 360 ms; if a CDN edge is 10 ms away it's 40 ms.
**Same content, 9x faster, purely from distance.**

CDNs also: cache static files, absorb traffic spikes, terminate TLS close to the
user, and provide DDoS protection.

### Anycast

Here's the clever bit that makes CDNs work: **the same IP address is announced
from many locations at once**, and internet routing (BGP, Chapter 13) naturally
delivers your packets to the nearest one.

```
   1.1.1.1 exists in London, Frankfurt, Singapore, Sao Paulo...
   You send a packet to 1.1.1.1 and it goes to whichever is closest to YOU.
```

This is how `1.1.1.1` and `8.8.8.8` answer in a few milliseconds from anywhere on
Earth. It also spreads a DDoS attack across the whole global capacity instead of
one data centre.

### Practice (15 min)

```bash
# See a CDN in action -- ask different resolvers, get different answers
dig @1.1.1.1 www.cloudflare.com +short
dig @8.8.8.8 www.cloudflare.com +short

# See which anycast location answers YOU
dig CHAOS TXT id.server @1.1.1.1 +short        # Cloudflare tells you the city!

# Compare latency to a CDN vs a single-location server
ping -c 5 1.1.1.1
ping -c 5 <a server you know is far away>

# Look at the headers a CDN adds
curl -sI https://www.cloudflare.com | grep -iE 'server|cf-ray|age|cache'
```

### Further reading

- **Article:** "What is a CDN?" and "What is anycast?" — Cloudflare Learning
  Center (they are, unsurprisingly, very good at explaining this).
- **Video:** "Load Balancing Explained" — ByteByteGo or Hussein Nasser.
- **Book:** *Site Reliability Engineering* (Google, free at sre.google) —
  the load balancing chapters.

---

## Chapter 37 — Cloud and container networking

### The idea

Everything you've learned still applies — cloud and container networks are built
from the same IP, TCP, routing, and NAT. They just add a layer of software
abstraction with new names.

### The translation table

| Cloud/container term | What it really is |
|---|---|
| **VPC / VNet** | A private IP address range you chose (e.g. `10.0.0.0/16`), plus routing |
| **Subnet** | A subnet (Chapter 10), usually pinned to one data centre |
| **Route table** | A routing table (Chapter 13) |
| **Internet Gateway** | NAT/routing to the public internet (Chapter 12) |
| **NAT Gateway** | NAT for private machines to reach out (Chapter 12) |
| **Security Group** | A **stateful** firewall attached to a machine — allow rules only |
| **Network ACL** | A **stateless** firewall on a subnet — so you must open return ports too |
| **Elastic/Public IP** | A public address 1:1 mapped to your private one |
| **Container network** | A software switch (Chapter 7) + veth "virtual cables" + NAT |
| **Kubernetes Service** | A virtual IP that gets rewritten (DNAT) to a real pod |
| **Overlay network (VXLAN)** | Your packet wrapped inside another packet to cross the real network |

**The one thing that bites everyone:** overlay networks (VXLAN, VPN, most
Kubernetes CNIs) wrap your packet in another packet, which **eats into the MTU**
(Chapter 15). Pod MTU is often 1450 or 1400, not 1500. If large requests hang in
Kubernetes but small ones work, check the MTU. This is the single most common
"impossible" container networking bug.

### Stateful vs stateless firewalls — a crucial cloud distinction

```
   SECURITY GROUP (stateful):
      You allow inbound port 443.  Return traffic is AUTOMATICALLY allowed.
      Simple.

   NETWORK ACL (stateless):
      You allow inbound port 443.  You must ALSO allow OUTBOUND on the
      ephemeral port range (1024-65535), or replies are blocked.
      This catches people out constantly.
```

### Practice (20 min)

```bash
# If you have Docker: look at the software network it built
ip -br addr | grep docker              # the bridge (a virtual switch)
ip route | grep docker                 # the route to container-land
sudo iptables -t nat -L -n | head -30  # the NAT rules Docker added

# Run two containers and watch them talk
docker run -d --name a alpine sleep 600
docker run -d --name b alpine sleep 600
docker exec a ip -br addr              # container a's private IP
docker exec a ping -c2 <b's ip>

# Check a container's MTU -- often NOT 1500
docker exec a cat /sys/class/net/eth0/mtu
```

### Further reading

- **Docs:** your cloud provider's VPC documentation — AWS's "VPC User Guide" is
  the most thorough and the concepts transfer.
- **Article:** "A Guide to the Kubernetes Networking Model" — search; several
  good ones. Also "Container Networking From Scratch" (a talk by Kristen Jacobs)
  which builds it up from `ip netns` commands.
- **Hands-on:** build a container network by hand with `ip netns`, `veth`, and
  `iptables`. Half a day, and it demystifies Docker completely.

---

## Chapter 38 — Security basics

### The uncomfortable truth

TCP/IP was designed in the 1970s for a small network of trusted research
institutions. **Security was explicitly a low priority.** Almost every security
problem you'll meet comes from that decision, and every fix is bolted on.

### The attacks worth understanding

| Attack | How it works | Defence |
|---|---|---|
| **Eavesdropping** | Anyone on the path reads your traffic (you saw this with `tcpdump -A` in Chapter 28) | **Encrypt everything** — TLS, SSH, VPN |
| **ARP spoofing** | Attacker on your LAN claims to be the router (Chapter 6 — ARP has no authentication) | Dynamic ARP Inspection on switches; encryption makes it useless |
| **DNS spoofing/poisoning** | Attacker gets you a fake answer, sending you to their server | DNSSEC, DoT/DoH, TLS certificate checks catch it |
| **SYN flood** | Thousands of SYNs, never completed, filling the server's half-open table (Chapter 21) | **SYN cookies** (on by default in Linux) |
| **DDoS** | Overwhelm with traffic from many sources | Upstream scrubbing/CDN — you cannot absorb it yourself |
| **Man-in-the-middle** | Attacker sits in the path and relays/modifies | TLS with **certificate validation** — this is why the browser warning matters |
| **Port scanning** | Probing to find open services | Firewall, close unused ports, monitor |
| **Amplification** | Small spoofed request to a service that sends a huge reply to the victim | Don't run open DNS/NTP resolvers; source-address filtering (BCP 38) |

### Firewalls

A firewall decides which packets to allow. Three generations:

1. **Stateless (packet filter)** — checks each packet's addresses/ports
   independently. Simple, but you have to open return traffic manually.
2. **Stateful** — tracks connections. "Allow new outbound, allow established
   return traffic." **This is what everything uses now.**
3. **Application-aware (NGFW/WAF)** — also inspects HTTP content, blocks SQL
   injection, etc.

**The one rule:** `DROP` (silently discard) vs `REJECT` (send back a refusal).
DROP makes attackers wait and reveals less; REJECT is friendlier for debugging
internal services. This is exactly the timeout-vs-refused distinction from
Chapter 21 — now you know why you see one or the other.

A sane default firewall shape:

```
   INBOUND:   default DENY
              allow established/related  (stateful -- return traffic)
              allow loopback
              allow ICMP types 3 and 11   <-- do NOT block all ICMP (Ch 14, 15)
              allow the specific services you actually run (e.g. 22, 443)
   OUTBOUND:  usually allow all (tighten in high-security environments)
```

### Practical hygiene

```bash
# 1. What is exposed on this machine? (do this on every server you run)
ss -tulnp

# 2. Is the firewall on?
sudo ufw status              # Ubuntu
sudo firewall-cmd --state    # RHEL/Fedora
sudo iptables -L -n          # raw
# macOS: System Settings > Network > Firewall

# 3. Are services bound to 0.0.0.0 that shouldn't be?
ss -tuln | grep '0.0.0.0'
```

**The single highest-value security habit in this guide:** run `ss -tulnp` on any
machine you're responsible for, and for each line ask *"should the world be able
to reach that?"* Databases bound to `0.0.0.0` with default passwords are how a
huge fraction of breaches begin.

### Practice (20 min)

```bash
# See plaintext traffic (why TLS exists)
sudo tcpdump -i any -A 'tcp port 80' -c 30 &
curl -s http://example.com > /dev/null
# your request and the response are fully readable

# Compare with HTTPS
sudo tcpdump -i any -A 'tcp port 443 and host example.com' -c 30 &
curl -s https://example.com > /dev/null
# handshake visible; content is gibberish

# Audit your own machine
ss -tulnp
```

### Further reading

- **Course (free):** Professor Messer's CompTIA Security+ videos — solid,
  vendor-neutral network security fundamentals.
- **Book:** *Practical Packet Analysis* has a good security-analysis chapter.
- **Site:** OWASP (owasp.org) for application-layer security; the OWASP Top 10 is
  essential reading for anyone building web services.
- **Practice (legal!):** TryHackMe or HackTheBox — hands-on labs you're allowed to
  attack. Never practise on systems you don't own.

---

# Part 10 — HTTP in depth

Part 7 taught you what an HTTP request looks like. This part is what you actually
need to *operate* the web: caching, cookies, CORS, and the streaming protocols.
**If you build web applications, this is the most immediately useful part of the
guide.**

## Chapter 39 — How HTTP messages are framed (and how it goes wrong)

### In one sentence

The receiver needs to know where one request ends and the next begins, and HTTP
gives two ways to say it — which is exactly why they can be made to disagree.

### The problem

You're reading bytes off a TCP connection (Part 5). TCP gives you a stream with no
message boundaries. So HTTP must mark them itself.

### The two ways

```
   METHOD 1: COUNT THE BYTES
   ------------------------------------
   POST /upload HTTP/1.1
   Host: example.com
   Content-Length: 27
                                <-- blank line
   {"name":"Ada","age":36}
   \_____________________/
     exactly 27 bytes; then the message is over


   METHOD 2: CHUNKS (when you don't know the length in advance)
   ------------------------------------
   HTTP/1.1 200 OK
   Transfer-Encoding: chunked

   1a                           <-- 0x1a = 26 bytes of data next
   {"status":"still working"}
   9                            <-- 9 bytes next
   more data
   0                            <-- a zero-size chunk means THE END
                                <-- blank line
```

**Chunked** is what lets a server start sending before it knows the total size —
streaming a log, generating a report, or an LLM producing tokens one at a time.

### The rule when both appear

If a message has **both** `Content-Length` and `Transfer-Encoding: chunked`, the
standard says: **`Transfer-Encoding` wins, and `Content-Length` must be
ignored** — or better, the message rejected.

Why does the standard have to say that? Because if two servers in a chain
disagree, you get a genuine security hole.

### Request smuggling, in plain language

Imagine two people reading the same letter through a slot, and disagreeing about
where the letter ends.

```
   A single stream of bytes arrives containing what LOOKS like one request.

   The front proxy reads it using Transfer-Encoding
        -> it thinks the request ends HERE  ------+
                                                  |
   The back-end server reads it using Content-Length
        -> it thinks the request ends HERE  --+   |
                                              |   |
   The bytes in between belong to nobody -- so the back end treats them as
   the BEGINNING OF THE NEXT PERSON'S REQUEST.
```

The attacker's leftover bytes get glued onto the front of the next customer's
request. That can steal their session, poison a cache, or bypass access controls.

**You don't need to memorise the variants** (CL.TE, TE.CL, TE.TE). You need to
know the defence:

- **Reject any request containing both headers**, at the edge.
- **Reject weird `Transfer-Encoding` values** — including ones with odd spacing
  like `Transfer-Encoding : chunked`, which one parser may accept and another
  ignore.
- Prefer **HTTP/2 end-to-end**, where framing is explicit and length-prefixed.
- Keep your proxy and your server on current versions; this class of bug is
  found and fixed regularly.

### Other framing rules worth knowing

| Case | Rule |
|---|---|
| `HEAD` response | Has a `Content-Length` but **no body at all**. Don't wait for one. |
| `204 No Content` and `304 Not Modified` | **Never** have a body, whatever the headers say |
| `100 Continue` | An interim reply meaning "go ahead, send the body" — so a client can ask permission before uploading 2 GB. You'll see *two* status lines for one request. |
| No `Content-Length`, no chunked (HTTP/1.0) | The body ends when the connection closes |

### Practice (20 min)

**A. Watch chunked encoding happen:**

```bash
curl -sv --raw http://httpbin.org/stream/3 2>&1 | tail -30
```

The `--raw` flag stops curl from decoding chunks, so you see the hex sizes and the
terminating `0`.

**B. Compare with a fixed-length response:**

```bash
curl -sI https://example.com | grep -iE 'content-length|transfer-encoding'
```

**C. Watch a streamed response arrive over time** (chunked in action):

```bash
curl -sN https://httpbin.org/drip?duration=5\&numbytes=10
```
Bytes trickle in over 5 seconds rather than arriving at once.

**D. See `100 Continue`:**

```bash
head -c 1000000 /dev/urandom > /tmp/big
curl -sv -H 'Expect: 100-continue' --data-binary @/tmp/big \
     https://httpbin.org/post 2>&1 | grep '< HTTP'
```
You should see `HTTP/1.1 100 Continue` and then `HTTP/1.1 200 OK`.

### Common confusions

- **"Is chunked slower?"** No — it lets the server start sending immediately,
  which usually makes the user experience *faster*.
- **"Why would a proxy and server disagree?"** Different codebases, different
  authors, different decades, different tolerance for malformed input. That's why
  strictness at the edge matters.
- **"Does this affect HTTP/2 and HTTP/3?"** Their native framing is unambiguous.
  But a proxy that *translates* HTTP/2 to HTTP/1.1 downstream can reintroduce the
  problem if it trusts a client-supplied length.

### Check yourself

1. What are the two ways HTTP marks the end of a message body?
2. When would you need chunked encoding?
3. In plain words, how does request smuggling work?
4. What should an edge proxy do with a request containing both headers?

### Further reading

- **Article (the definitive one):** "HTTP request smuggling" — PortSwigger Web
  Security Academy. Free, interactive labs included. Genuinely excellent.
- **Docs:** MDN, "Transfer-Encoding" and "Content-Length".
- **Spec:** RFC 9112 §6 — message body length rules. Short and precise.

---

## Chapter 40 — Caching: the biggest performance lever there is

### In one sentence

The fastest request is the one you never make, and HTTP has a rich, precise
system for saying how long a response may be reused.

### The problem

Every request costs round trips (Chapter 4's latency budget), bandwidth, and
server work. Most content doesn't change between requests. Caching removes the
request entirely.

### The three caches

```
   BROWSER CACHE          CDN / PROXY CACHE           ORIGIN SERVER
   (private -- one user)  (shared -- many users)      (the source of truth)
        |                        |                          |
   fastest, free            fast, near the user         slow, expensive
```

The **private vs shared** distinction runs through everything below: a response
containing your bank balance may be cached by *your browser* but must never be
cached by a *shared* CDN.

### `Cache-Control`: the directives that matter

```
   Cache-Control: public, max-age=31536000, immutable
                  \____/  \______________/  \_______/
                  who may   how long fresh   never even
                  store it  (seconds)        revalidate
```

| Directive | Means |
|---|---|
| `max-age=N` | Fresh for N seconds |
| `s-maxage=N` | Same, but **for shared caches only** (overrides `max-age` there) |
| `public` | Shared caches may store it |
| `private` | **Browser only** — CDNs must not store it |
| `no-cache` | You may store it, but **check with me before every reuse** |
| `no-store` | **Never store it anywhere.** This is the real "don't cache". |
| `must-revalidate` | Once stale, don't serve it — check or fail |
| `immutable` | This will never change; don't revalidate even on a reload |
| `stale-while-revalidate=N` | Serve the stale copy for up to N s while fetching a fresh one in the background |

> **The single most common mistake in web development:** using `no-cache` when
> you mean `no-store`. `no-cache` **permits** storage — it just requires
> revalidation. If the response contains something secret, you want `no-store`.

### Revalidation: the 304

When a cached copy goes stale, the client doesn't re-download blindly — it asks
"has this changed?"

```
   Original response:                Later request:
     ETag: "v3-abc123"                 If-None-Match: "v3-abc123"
     Last-Modified: Tue, 01 ...        If-Modified-Since: Tue, 01 ...

   Server, if unchanged:
     HTTP/1.1 304 Not Modified
     (headers only -- NO BODY)

   -> the client reuses what it already has.
      A 2 MB image cost ~200 bytes to confirm.
```

**`ETag`** is a version token — usually a hash of the content. **`Last-Modified`**
is a timestamp (1-second resolution, so `ETag` is more precise).

`304` is not an error. Seeing lots of them in your logs means caching is *working*.

### `Vary`: the header that keeps caches honest

A cache stores one response per URL. But what if the response depends on a request
*header*?

```
   Same URL, different responses:
     Accept-Encoding: gzip     ->  compressed bytes
     Accept-Encoding: (none)   ->  plain bytes

   Without "Vary: Accept-Encoding", a cache might hand compressed
   bytes to a client that can't decompress them.
```

```
   Vary: Accept-Encoding     almost always needed if you compress
   Vary: Accept-Language     one cached copy per language
   Vary: Cookie              -> effectively uncacheable in a shared cache
   Vary: User-Agent          -> cache explosion; avoid
```

### The pattern that actually works

This is the standard modern setup, and it's worth copying exactly:

```
   ASSETS (JS, CSS, images) -- give them a content hash in the filename:
       /static/app.9f3c2a1.js
       Cache-Control: public, max-age=31536000, immutable
       -> cached for a year, never revalidated.
          Content changed? The FILENAME changes, so it's a different URL.

   HTML -- the file that references those assets:
       Cache-Control: no-cache
       -> always revalidated (cheap: a 304 if unchanged),
          so a deploy is visible immediately.

   API responses (personal data):
       Cache-Control: private, no-store
```

This gives you near-perfect caching *and* instant deploys, with no purging.

### Practice (30 min)

**A. Look at real caching headers:**

```bash
curl -sI https://example.com | grep -iE 'cache-control|etag|last-modified|age|vary|expires'
```

Try a few sites and compare. Look at a big site's static assets versus its HTML.

**B. Do a conditional request yourself and get a 304:**

```bash
URL=https://example.com
ET=$(curl -sI $URL | awk -F': ' 'tolower($1)=="etag"{print $2}' | tr -d '\r')
echo "ETag is: $ET"
curl -sI -H "If-None-Match: $ET" $URL | head -1      # expect: 304 Not Modified
```

**C. See a CDN hit vs miss:**

```bash
curl -sI https://www.cloudflare.com | grep -iE 'cf-cache-status|age'
curl -sI https://www.cloudflare.com | grep -iE 'cf-cache-status|age'   # run twice
```
`Age:` tells you how many seconds the cache has held it.

**D. In your browser:** DevTools → Network. Reload a page and look at the **Size**
column — entries saying "(disk cache)" or "(memory cache)" cost zero network. Then
tick **Disable cache** and reload; watch the load time and transfer size change.
**That difference is what caching is worth.**

### Common confusions

- **"`no-cache` means don't cache."** It doesn't. `no-store` does.
- **"Why is my update not showing?"** Something along the path cached it with a
  long `max-age` and no revalidation. This is why hashed filenames exist.
- **"Should I use `Expires` or `Cache-Control`?"** `Cache-Control` — it's newer and
  wins when both are present. `Expires` is an absolute date and suffers from clock
  skew.
- **"Does HTTPS prevent caching?"** No. Browsers and your own CDN cache HTTPS
  responses normally. Only *unrelated* intermediaries can't.

### Check yourself

1. What's the difference between `no-cache` and `no-store`?
2. What is an `ETag` and what does a `304` save you?
3. Why does `Vary: Accept-Encoding` matter?
4. Why do modern sites put a hash in asset filenames?

### Further reading

- **Article (best on the topic):** "Caching best practices & max-age gotchas" —
  Jake Archibald. Explains the hashed-filename pattern definitively.
- **Docs:** MDN, "HTTP caching" — thorough and accurate.
- **Book (free):** *High Performance Browser Networking*, the HTTP chapters.
- **Tool:** `redbot.org` — paste a URL, get an expert analysis of its caching
  headers and what's wrong with them. Genuinely useful.

---

## Chapter 41 — Cookies: how the web remembers you

### In one sentence

HTTP is stateless, so the server hands the browser a small piece of text and the
browser sends it back on every matching request — automatically, which is both the
feature and the danger.

### How they work

```
   Server -> browser:
     Set-Cookie: session=abc123; Path=/; Secure; HttpOnly; SameSite=Lax; Max-Age=3600

   Browser -> server, on EVERY matching request from then on:
     Cookie: session=abc123
```

The word **automatically** is the important one. The browser attaches cookies
based on scope, without the page asking. That's what makes sessions work — and
what makes CSRF possible.

### The attributes, and what each actually protects

| Attribute | Effect | Why you want it |
|---|---|---|
| `Secure` | HTTPS only | Stops the cookie being sent in the clear. **Always set it.** |
| `HttpOnly` | JavaScript can't read it | If your site has an XSS bug, the attacker still can't steal the session token |
| `SameSite=Lax` | Not sent on cross-site requests, except top-level navigation | **Blocks most CSRF.** This is the modern browser default. |
| `SameSite=Strict` | Never sent cross-site at all | Stronger, but breaks "click the link in the email and be logged in" |
| `SameSite=None` | Sent everywhere — **requires `Secure`** | Only for genuine third-party use (embeds, SSO) |
| `Domain=example.com` | Sent to that domain **and all its subdomains** | **Omitting `Domain` is narrower and safer** |
| `Path=/app` | Only for URLs under that path | Convenience, **not security** |
| `Max-Age` / `Expires` | Survives browser restart | Without either, it's gone when the browser closes |
| `__Host-` name prefix | Browser *enforces* `Secure`, `Path=/`, and no `Domain` | The strongest available binding |

### The trap: cookies are not origin-scoped

Everything else on the web is scoped to an **origin** (`scheme://host:port`).
Cookies are not:

```
   Set by:  http://example.com:8080
   Sent to: https://example.com:443     <-- yes, really.

   Cookies ignore the PORT, and (mostly) the SCHEME.
```

And the `Domain` attribute is worse:

```
   Set-Cookie: session=abc; Domain=example.com

   -> sent to  app.example.com        (intended)
   -> sent to  blog.example.com       (probably fine)
   -> sent to  promo.example.com      (run by an agency... with an XSS bug)
```

**Rule:** don't set `Domain` unless you truly need subdomain sharing, and never
host untrusted content on a subdomain of your main site.

### CSRF, and why `SameSite` changed everything

```
   You are logged into your bank (cookie stored).
   You visit evil.com, which contains:

       <form action="https://bank.com/transfer" method="POST">
         <input name="to" value="attacker"><input name="amount" value="1000">
       </form>
       <script>document.forms[0].submit()</script>

   The browser HELPFULLY attaches your bank cookie. The bank sees a
   perfectly authenticated request.
```

Defences, strongest first:
1. **`SameSite=Lax` or `Strict`** — the browser simply doesn't send the cookie.
   Now the default, which killed most classic CSRF overnight.
2. **Anti-CSRF token** — a random value in the form, verified server-side. Still
   needed for `SameSite=None` flows.
3. **Check the `Origin` header** server-side on state-changing requests.
4. **Use `Authorization: Bearer` instead of cookies** — headers aren't attached
   automatically, so CSRF doesn't apply. (But you then have to store the token
   somewhere, and `localStorage` is readable by XSS. Pick your trade-off.)

### Practice (20 min)

**A. See cookies being set:**

```bash
curl -sI https://github.com | grep -i set-cookie
curl -sI https://example.com | grep -i set-cookie
```

**B. Inspect them properly in the browser:** DevTools → **Application** →
Cookies. The columns are literally the attributes above. Look at a site you're
logged into: is the session cookie `Secure`? `HttpOnly`? What's its `SameSite`?

**C. Watch a cookie round-trip:**

```bash
curl -s -c /tmp/jar.txt https://httpbin.org/cookies/set?flavour=chocolate > /dev/null
cat /tmp/jar.txt                                        # the stored cookie
curl -s -b /tmp/jar.txt https://httpbin.org/cookies     # sent back automatically
```

**D. Measure the cost.** Cookies are sent on **every** request to that scope —
including images and API calls. Check the size:

```bash
curl -sI https://<a site you use> | grep -i set-cookie | wc -c
```
Multiply by every request on a page load. Large cookies are a real performance
problem, which is why static assets are often served from a cookie-less domain.

### Common confusions

- **"Cookies are insecure."** Cookies with `Secure`, `HttpOnly`, and `SameSite`
  are a solid session mechanism — often safer than tokens in `localStorage`,
  which JavaScript (and therefore XSS) can read.
- **"`Path` protects my cookie."** It doesn't. JavaScript on the same origin can
  read cookies across paths.
- **"Third-party cookies are being removed — does that break sessions?"** No. Your
  own site's cookies (first-party) are unaffected. It's cross-site tracking
  cookies that are going away.

### Check yourself

1. Why does the browser send cookies automatically, and what problem does that
   create?
2. What does `HttpOnly` protect against? What about `SameSite`?
3. Why is omitting the `Domain` attribute safer than setting it?
4. Name two defences against CSRF.

### Further reading

- **Docs:** MDN, "Using HTTP cookies" and "Set-Cookie" — accurate and complete.
- **Article:** "SameSite cookies explained" — web.dev. Clear, with diagrams.
- **Article:** OWASP "Cross-Site Request Forgery Prevention Cheat Sheet".
- **Labs:** PortSwigger Web Security Academy — free CSRF labs you can actually
  exploit.

---

## Chapter 42 — CORS: why the browser blocked your request

### In one sentence

Browsers stop JavaScript on one site from reading responses from another site, and
CORS is how a server grants permission — enforced **only** in the browser.

### The thing to understand first

> **CORS is not server security.** It restricts what **browser JavaScript** may
> *read*. `curl`, mobile apps, servers, and scripts ignore it entirely. Never use
> CORS to protect data — use authentication.

CORS exists to protect *users* from *malicious websites*, not to protect your API
from attackers.

### What an "origin" is

```
   https://app.example.com:443
   \___/   \_____________/ \_/
   scheme      host       port      -- ALL THREE must match

   https://app.example.com    vs  http://app.example.com     DIFFERENT (scheme)
   https://app.example.com    vs  https://api.example.com    DIFFERENT (host)
   https://app.example.com    vs  https://app.example.com:8443  DIFFERENT (port)
```

### The two kinds of request

```
   SIMPLE REQUEST -- sent straight away
     * method is GET, HEAD, or POST
     * AND Content-Type is one of:
         text/plain, application/x-www-form-urlencoded, multipart/form-data
     * AND no custom headers

     -> browser SENDS it, then checks the response headers before letting
        your JavaScript READ the result.
        (Note: the request already happened. Side effects already occurred.)


   PREFLIGHTED REQUEST -- an extra round trip first
     * anything else: PUT, DELETE, PATCH
     * OR Content-Type: application/json          <-- this is the usual trigger
     * OR custom headers (Authorization, X-Api-Key, ...)

     -> browser first sends:
          OPTIONS /resource
          Origin: https://app.example.com
          Access-Control-Request-Method: PUT
          Access-Control-Request-Headers: content-type, authorization

     -> server must reply with matching permissions
     -> ONLY THEN does the real request go
```

**`Content-Type: application/json` is what makes nearly every modern API call
preflighted.** That's why "my GET works but my POST doesn't" is the classic CORS
complaint.

### The headers the server must send

```
   Access-Control-Allow-Origin: https://app.example.com
   Access-Control-Allow-Methods: GET, POST, PUT, DELETE
   Access-Control-Allow-Headers: Content-Type, Authorization
   Access-Control-Allow-Credentials: true        (only if using cookies)
   Access-Control-Expose-Headers: X-Request-Id   (headers JS may READ)
   Access-Control-Max-Age: 86400                 (cache the preflight)
   Vary: Origin                                  (if you echo the origin back)
```

### The four errors you will actually hit

**1. `Allow-Origin: *` together with credentials.**
If the request sends cookies (`credentials: 'include'`), the wildcard is
forbidden. You must echo the **specific** origin, add
`Access-Control-Allow-Credentials: true`, **and** add `Vary: Origin` (or a cache
will serve one origin's headers to another).

**2. The preflight isn't handled.**
Your framework routes `PUT /items` but returns 404 or 405 for `OPTIONS /items`.
The real request never happens. Handle `OPTIONS` explicitly.

**3. A custom header isn't listed.**
You send `X-Request-Id`; the preflight response doesn't include it in
`Access-Control-Allow-Headers`. Blocked.

**4. You can't read a response header.**
By default JavaScript can only see a small safelist. Anything else needs
`Access-Control-Expose-Headers`.

**And one that wastes hours:** *an error response still needs CORS headers.* If
your server throws a 500 without them, the browser reports a **CORS error**, hiding
the real failure. Add the headers in your error handler too.

### Practice (25 min)

**A. Simulate a preflight from the command line:**

```bash
curl -sI -X OPTIONS https://httpbin.org/anything \
  -H 'Origin: https://app.example.com' \
  -H 'Access-Control-Request-Method: PUT' \
  -H 'Access-Control-Request-Headers: content-type' \
  | grep -iE 'access-control|^HTTP'
```

**B. A simple request:**

```bash
curl -sI https://httpbin.org/get -H 'Origin: https://app.example.com' \
  | grep -i access-control
```

**C. See it fail and succeed in a browser.** Open DevTools console on any site
and run:

```javascript
// this will be BLOCKED by CORS (no permission headers from example.com)
fetch('https://example.com').then(r => r.text()).then(console.log).catch(console.error)

// this will WORK (httpbin sends Access-Control-Allow-Origin: *)
fetch('https://httpbin.org/get').then(r => r.json()).then(console.log)
```

Read the error message in the console carefully — it tells you exactly which
header was missing.

**D. Prove that CORS is browser-only:** the same `https://example.com` request
that the browser blocked works perfectly from `curl`. **The data was never
protected; only the browser's JavaScript was restricted.**

### Common confusions

- **"CORS is blocking my server-to-server call."** It isn't. Only browsers enforce
  CORS. If your server call fails, it's a firewall, DNS, or auth problem.
- **"I'll just disable CORS."** You can't — it's the browser's rule. You can only
  make the *server* grant permission. (Browser flags/extensions that disable it
  are for local debugging only.)
- **"A proxy fixes CORS."** Yes, and that's a legitimate pattern: your front end
  calls your own origin, which forwards to the third-party API server-side. Same
  origin, no CORS involved.
- **"Why did it work in Postman?"** Postman isn't a browser and doesn't enforce
  CORS.

### Check yourself

1. What is an "origin", exactly?
2. Why does sending JSON usually trigger a preflight?
3. Why can't you use `Access-Control-Allow-Origin: *` with cookies?
4. Does CORS protect your API from an attacker with `curl`?

### Further reading

- **Docs (start here):** MDN, "Cross-Origin Resource Sharing (CORS)" — includes a
  flowchart and a list of every error message with its cause.
- **Article:** "CORS in 100 Seconds" (Fireship) for a quick video, then MDN for
  the detail.
- **Tool:** `test-cors.org` — interactively fire cross-origin requests and see
  what happens.

---

## Chapter 43 — WebSockets and Server-Sent Events

### In one sentence

Normal HTTP is the client asking and the server answering; these two let the
**server push** data to the client — one bidirectionally, one simply.

### The problem

You want live updates: a chat message, a stock price, a build status, tokens
streaming from an AI model. Plain HTTP means the client has to keep asking
("polling"), which is wasteful and laggy.

### WebSocket: a two-way pipe

A WebSocket **starts as HTTP and then stops being HTTP**.

```
   CLIENT:                                    SERVER:
   GET /chat HTTP/1.1
   Host: example.com
   Upgrade: websocket
   Connection: Upgrade
   Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
   Sec-WebSocket-Version: 13
                        ------------------->
                                              HTTP/1.1 101 Switching Protocols
                                              Upgrade: websocket
                                              Connection: Upgrade
                        <-------------------  Sec-WebSocket-Accept: s3pPLM...

   === the SAME TCP connection now carries WebSocket messages, both ways,
       until someone closes it ===
```

Status `101` is the handshake succeeding. After that there are no more HTTP
requests — just messages.

### Server-Sent Events: a one-way stream

SSE is far simpler: it's **an HTTP response that never ends**.

```
   GET /events           Accept: text/event-stream

   HTTP/1.1 200 OK
   Content-Type: text/event-stream
   Cache-Control: no-cache

   data: {"price": 10.50}

   data: {"price": 10.52}

   id: 42
   event: alert
   data: {"msg": "threshold crossed"}

   : this is a comment -- often used as a heartbeat

   ... connection stays open, more events arrive ...
```

In the browser it's three lines:

```javascript
const es = new EventSource('/events');
es.onmessage = e => console.log(JSON.parse(e.data));
```

### Which one to use

| | **WebSocket** | **Server-Sent Events** |
|---|---|---|
| Direction | **Both ways** | Server → client only |
| Protocol | Custom, after upgrade | **Plain HTTP** |
| Reconnects automatically | No — you write it | **Yes, built in** |
| Resume after a drop | You build it | **Built in** (`id:` + `Last-Event-ID`) |
| Works through proxies/CDNs | Needs explicit support | Works through anything that streams |
| Binary data | Yes | Text only |
| Browser connection limit | Fine | 6 per origin on HTTP/1.1 — **use HTTP/2** |
| Complexity | Higher | **Very low** |

**Rule of thumb:** if the client only *receives*, use SSE. If the client sends
frequently too (chat, collaborative editing, games), use WebSocket.

Most "we need WebSockets" requirements are actually satisfied by SSE with far
less code. LLM token streaming, notifications, progress bars, and live dashboards
are all one-directional.

### The hazard that bites everyone: idle timeouts

A long-lived connection with no traffic gets **silently killed** by something in
the middle — a load balancer, a proxy, a NAT (Chapters 12 and 23). Common
timeouts: 60 seconds on many load balancers, 350 seconds on AWS NLB, a few minutes
on home NAT.

```
   SYMPTOM: the connection dies roughly every N seconds when idle,
            with WebSocket close code 1006 ("closed abnormally, no
            close frame") and no useful error message.

   FIX:     send a heartbeat more often than the shortest timeout on the path.
              * WebSocket: use protocol ping/pong frames, every 20-30 s
              * SSE: send a comment line (":\n\n") every 15-30 s
            AND raise the load balancer's idle timeout.
```

**Close code `1006` almost always means a middlebox timeout or a network drop**,
not an application error.

Also plan for: **reconnect with exponential backoff *and jitter*** (otherwise
every client reconnects simultaneously and takes you down), and **resume from a
cursor** so a reconnect doesn't lose or duplicate messages.

### Practice (25 min)

**A. Watch a WebSocket handshake:**

```bash
curl -sv -N \
  -H 'Connection: Upgrade' -H 'Upgrade: websocket' \
  -H 'Sec-WebSocket-Version: 13' \
  -H 'Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==' \
  https://echo.websocket.org/ 2>&1 | grep -iE '^< HTTP|upgrade|sec-websocket'
```
Look for `101 Switching Protocols`.

**B. Actually use a WebSocket:**

```bash
# install: apt install websocat  /  brew install websocat
websocat wss://echo.websocket.org
# type anything; it echoes back. Ctrl-C to quit.
```

**C. Consume a real SSE stream** (Wikipedia publishes every edit, live):

```bash
curl -N -H 'Accept: text/event-stream' \
  https://stream.wikimedia.org/v2/stream/recentchange | head -20
```
Note that it never finishes — that's the point. `Ctrl-C` to stop.

**D. In the browser** (DevTools console on any page):

```javascript
const es = new EventSource('https://stream.wikimedia.org/v2/stream/recentchange');
es.onmessage = e => console.log(JSON.parse(e.data).title);
setTimeout(() => es.close(), 10000);   // stop after 10 seconds
```

**E. See it in DevTools:** Network tab → filter by **WS** to inspect WebSocket
frames, or click an `eventsource` request to see the message stream.

### Common confusions

- **"WebSockets are faster than HTTP."** For *frequent bidirectional* messages,
  yes — no per-message headers or round trips. For request/response, HTTP/2 is
  comparable and far simpler.
- **"SSE is old / deprecated."** No — it's a stable standard, and it's what most
  AI chat interfaces use for token streaming.
- **"Why does my WebSocket keep dying at exactly 60 seconds?"** Idle timeout.
  Add a heartbeat. (Round numbers are always a timer — Chapter 26's lesson again.)
- **"Can I use WebSockets through a corporate proxy?"** Sometimes. Many block the
  upgrade. SSE, being ordinary HTTP, usually passes.

### Check yourself

1. What HTTP status code indicates a successful WebSocket upgrade?
2. Name three things SSE gives you for free that you'd have to build with
   WebSockets.
3. What does close code 1006 usually mean?
4. Why does a heartbeat fix it?

### Further reading

- **Docs:** MDN, "The WebSocket API" and "Using server-sent events". Both have
  runnable examples.
- **Article:** "WebSockets vs Server-Sent Events" — several good comparisons;
  search that phrase and read one that includes the proxy considerations.
- **Video:** Hussein Nasser's WebSocket and SSE videos — practical and honest
  about the operational problems.
- **Tool:** `websocat` — the netcat of WebSockets, invaluable for debugging.

---

## Chapter 44 — Compression, range requests, and proxies

Three smaller topics that round out practical HTTP.

### Compression

```
   Client says what it can handle:   Accept-Encoding: br, gzip
   Server says what it used:         Content-Encoding: br
   Server MUST also send:            Vary: Accept-Encoding
```

| Codec | Notes |
|---|---|
| **gzip** | Universal, cheap, good ratio. The safe default. |
| **brotli (`br`)** | ~15–20% smaller than gzip on text. Expensive at high levels — **pre-compress static files at maximum level**, use a low level for dynamic responses. |
| **zstd** | Fast, growing support |

Rules that matter:
- **Don't compress already-compressed files** (JPEG, PNG, MP4, ZIP). You burn CPU
  and may make them bigger.
- **Set a minimum size** (~1 KB). Compressing a 50-byte response makes it larger.
- **The `Vary` header is mandatory**, or a cache will serve compressed bytes to a
  client that asked for plain ones.

### Range requests

Lets a client ask for **part** of a resource.

```
   Server advertises:  Accept-Ranges: bytes
   Client asks:        Range: bytes=1000-1999
   Server replies:     206 Partial Content
                       Content-Range: bytes 1000-1999/50000
```

This is how: video seeking works, downloads resume after a drop, and download
managers fetch several pieces in parallel.

**The one detail that matters for correctness:** use `If-Range`.

```
   Range: bytes=5000-
   If-Range: "v3-abc123"        <-- "only if it's still the version I started"

   -> unchanged?  206 with the requested range
   -> changed?    200 with the WHOLE new file (safe restart)
```

Without `If-Range`, a resumed download can splice together bytes from two
different versions of a file and produce silent corruption.

### Proxies — the forward kind

Chapter 36 covered **reverse** proxies (in front of servers). A **forward** proxy
sits in front of *clients* — typically a corporate egress gateway.

```
   FORWARD:  your machine -> [corporate proxy] -> the internet
   REVERSE:  the internet -> [load balancer]   -> the company's servers
```

For HTTPS, a forward proxy can't read the traffic, so it builds a blind tunnel:

```
   CLIENT -> PROXY:   CONNECT example.com:443 HTTP/1.1
   PROXY  -> CLIENT:  HTTP/1.1 200 Connection Established
   ... then it just relays bytes; TLS is end-to-end ...
```

The proxy sees the hostname and the volume, but not the content.

**TLS interception** is when a company *does* want to see the content: the proxy
terminates TLS with its own certificate and re-encrypts to the real server. This
requires installing a **corporate root CA** on every device, and it causes a very
specific, very common problem:

> **"SSL certificate problem: unable to get local issuer certificate"** — but only
> at work, and only in some tools.

The reason: the corporate CA got installed in the OS trust store (so browsers and
`curl` work), but **Python, Node, Java, Go, and Docker each use their own trust
store**. Each needs the CA added separately.

Also learn the environment variables, because everything respects them:

```bash
export http_proxy=http://proxy.corp:3128
export https_proxy=http://proxy.corp:3128
export no_proxy=localhost,127.0.0.1,.internal.example.com
```

### Practice (20 min)

```bash
# compression -- see it applied, and check Vary
curl -sI -H 'Accept-Encoding: gzip, br' https://example.com \
  | grep -iE 'content-encoding|vary|content-length'

# compare sizes
curl -s https://example.com | wc -c
curl -s -H 'Accept-Encoding: gzip' https://example.com --output - | wc -c

# range requests
curl -sI https://example.com | grep -i accept-ranges
curl -s -r 0-99 https://example.com | wc -c            # expect ~100 bytes
curl -sI -H 'Range: bytes=0-99' https://example.com | grep -iE '^HTTP|content-range'

# what does your environment think about proxies?
env | grep -i proxy

# which CA bundle is each tool using?
curl -v https://example.com 2>&1 | grep -i 'CAfile'
python3 -c "import certifi; print(certifi.where())"

# are you behind TLS interception? check who signed the certificate
echo | openssl s_client -connect example.com:443 2>/dev/null | openssl x509 -noout -issuer
```

If that last command names your employer instead of a public CA, your traffic is
being intercepted.

### Check yourself

1. Why is `Vary: Accept-Encoding` required when you compress?
2. Which files should you *not* compress, and why?
3. What does `If-Range` prevent?
4. Why does `curl` work but Python fail on a corporate network?

### Further reading

- **Docs:** MDN — "Content-Encoding", "Range requests", "Proxy servers and
  tunneling".
- **Article:** "Brotli vs gzip" benchmarks — search for a recent one; the
  trade-off between compression level and CPU is the interesting part.
- **Article:** search "corporate TLS interception certificate trust store" for
  practical fixes per language.

---

### End of Part 10 — Milestone check

- [ ] I can explain the two ways HTTP frames a message body
- [ ] I know the difference between `no-cache` and `no-store`
- [ ] **I have triggered a 304 with a conditional request**
- [ ] I can name the four cookie attributes that matter and what each protects
- [ ] I can explain why a JSON POST triggers a CORS preflight
- [ ] I know when to choose SSE over WebSockets
- [ ] I can explain why a WebSocket dies every 60 seconds

---

# Part 11 — Watching your network (monitoring)

Part 8 was about investigating a problem you already knew about. This part is
about **finding out first**.

## Chapter 45 — What to measure, and why averages lie

### In one sentence

Measure what users feel (latency, errors, throughput) and what predicts failure
(errors, drops, saturation) — and always look at percentiles, never averages.

### The four kinds of signal

```
   1. DEVICE / LINK METRICS    "is the pipe healthy?"
      utilisation, errors, drops, interface up/down
      -> cheap, constant, from SNMP or kernel counters

   2. FLOW RECORDS             "who is talking to whom, and how much?"
      summaries of each conversation
      -> answers "what is using my bandwidth?"

   3. ACTIVE PROBES            "what would a user experience right now?"
      ping, TCP connect, HTTP checks, DNS lookups
      -> works even when there is no real traffic (3am outages)

   4. PACKET CAPTURE           "what EXACTLY happened?"
      -> expensive; on demand or in a small rolling buffer
```

Cost rises as you go down; precision rises with it. Run **1 and 3 always**,
**2 sampled**, **4 when you need it**.

### The metrics that catch real incidents

**On a link or interface:**

| Metric | Alert when |
|---|---|
| Utilisation (in/out) | Above ~70–80% sustained — queueing starts long before 100% |
| **Errors / CRC** | **Any sustained increase.** This is a failing cable or optic. |
| Discards / drops | Rising — buffers are overflowing |
| Interface flaps | More than once an hour |

**On a Linux host** (these are the ones people forget):

```bash
nstat -az | grep -Ei 'ListenOverflows|ListenDrops|TCPBacklogDrop|RetransSegs'
```

| Counter | Means |
|---|---|
| `TcpExtListenOverflows` | **The accept queue is full** — connections are being dropped *before* your app sees them |
| `TcpRetransSegs` / `TcpOutSegs` | Ratio ≈ your packet loss rate |
| `UdpRcvbufErrors` | UDP packets dropped because the app read too slowly |

> **This is the most valuable point in the chapter:** connection-level failures
> are **invisible in application logs**, because the connection never became a
> request. If you don't scrape kernel counters, you will never see them.

**On a service:** request rate, error rate (**4xx and 5xx counted separately** —
Chapter 28), latency percentiles, and the timing breakdown from Chapter 4 (DNS /
TCP / TLS / first byte).

### Why averages lie

```
   "Average latency: 50 ms"        sounds completely fine.

   Reality:
      p50 (half of users)   =    20 ms
      p95                    =   180 ms
      p99 (1 in 100)         = 4,000 ms   <-- these people are furious

   The average hid the entire problem.
```

**Rules:**
- Always chart **p50, p95, p99**.
- **Never average a percentile** across servers or time buckets — it's
  mathematically meaningless. Aggregate from histograms.
- For saturation, plot the **maximum**. A five-minute average hides a
  ten-second stall completely.

### Watch for round numbers

Chapter 26 taught you this and it applies to monitoring too. If a latency chart
has a hard floor or a cluster at a suspiciously round value, it's a **timer**, not
physics:

```
   40 ms    Nagle + delayed ACK      (Chapter 26)
   200 ms   TCP retransmission timer (Chapter 22)
   1/3/7 s  SYN retries              (Chapter 21)
   5 s      a dead DNS server        (Chapter 18)
   60 s     an idle timeout          (Chapters 23, 43)
```

### Practice (20 min)

```bash
# the kernel counters nobody watches
nstat -az | grep -Ei 'ListenOverflows|ListenDrops|RetransSegs|OutSegs|UdpRcvbufErrors'

# interface errors (should be zero and STAY zero)
ip -s link show eth0
ethtool -S eth0 2>/dev/null | grep -Ei 'err|drop|crc'

# take two samples 10 seconds apart -- rates matter, not totals
nstat -n; sleep 10; nstat

# per-connection health (Chapter 33)
ss -ti state established | grep -oE 'rtt:[0-9.]+|retrans:[0-9]+/[0-9]+' | head
```

Write your baseline numbers into `notes.md`. **You cannot spot an anomaly without
knowing what normal looks like.**

### Check yourself

1. Name the four kinds of monitoring signal.
2. Why can't application logs show you accept-queue overflows?
3. Why is an average latency figure misleading?
4. You see latency clustered at exactly 200 ms. What does that suggest?

### Further reading

- **Book (free):** *Site Reliability Engineering* (Google) — sre.google. The
  "Monitoring Distributed Systems" chapter defines the four golden signals.
- **Article:** Brendan Gregg's "USE Method" (brendangregg.com) — Utilisation,
  Saturation, Errors, applied per resource.
- **Article:** "How NOT to measure latency" — Gil Tene. On why averages and
  naive percentiles mislead. Watch the talk version if you prefer.

---

## Chapter 46 — Where the numbers come from: counters, flows, and probes

### In one sentence

Devices expose counters (SNMP), summarise conversations (flow records), and you
generate your own test traffic (probes) — three different tools for three
different questions.

### SNMP: reading a device's counters

**Simple Network Management Protocol** is how you ask a switch, router, printer,
or UPS for its numbers. It's from 1988 and it's everywhere.

```
   YOUR MONITORING SERVER              THE DEVICE
        |   "what is ifInOctets on interface 3?"   |
        |------------------- UDP 161 ------------->|
        |<------------------- 4823910284 ----------|
        |                                          |
        |<---- TRAP (UDP 162): "interface 3 down!" |   (unsolicited alert)
```

Every value has a numeric address called an **OID**, and a **MIB** is the
dictionary that gives them names. The handful you'll actually use:

| Name | Tells you |
|---|---|
| `sysUpTime` | Time since boot — **if it resets, the device rebooted** |
| `ifDescr` / `ifAlias` | Interface name / the human description someone typed |
| `ifHCInOctets` / `ifHCOutOctets` | Bytes in/out — **the 64-bit versions** |
| `ifInErrors` / `ifOutErrors` | Errors |
| `ifOperStatus` | Up or down |

> **The classic trap:** the original 32-bit byte counters wrap around at ~4.3 GB.
> On a 10 Gbit/s link that's **under 4 seconds**. Always poll the 64-bit `ifHC*`
> counters, or your graphs will be nonsense.

**Security note:** SNMP v1 and v2c send the "community string" (effectively a
password) in **plain text**. Use v3 (which has real authentication and
encryption), and restrict SNMP to a management network.

**A limitation to understand:** if you poll every 60 seconds, you're seeing
60-second averages. A 10-second burst that saturates the link **averages into
invisibility**. That's why "the graph says 45% but users say it's slow" happens.

### Flow records: who is using the bandwidth?

Counters tell you *how much*. Flow records tell you **who and what**.

```
   A "flow" is one conversation. The router summarises it:

   10:00:01 -> 10:00:04   192.168.1.50:51000 -> 93.184.216.34:443  TCP
                          48 packets, 52,104 bytes
```

Collect millions of those and you can answer: *who are my top talkers? what
protocol is eating the uplink? has anything ever connected to this address?*

| Technology | How it works | Trade-off |
|---|---|---|
| **NetFlow / IPFIX** | Router tracks every flow and exports a record when it ends | Accurate byte counts; costs router memory |
| **sFlow** | Samples 1 packet in N, plus interface counters | Cheap and real-time; statistical, so small flows are invisible |
| **VPC Flow Logs** | The cloud equivalent — same idea, written to storage | Also shows **ACCEPT/REJECT**, so you can see what a security group blocked |

**A timing detail that confuses people:** flow records appear *after* the flow
ends (or after a periodic timer, often 60 seconds). So flow data is minutes
behind — never build a real-time alert on it.

### Active probes: testing before users do

Passive monitoring only sees traffic that exists. If your service is down at 3am
with no users, **nothing is broken from the monitoring's point of view**. Active
probes fix that by generating their own traffic.

Design them in **layers**, so the alert tells you where the failure is:

```
   ping the host             -> is it reachable at all?          (Layer 3)
   TCP connect to :443       -> is the port open?                (Layer 4)
   TLS handshake             -> is the certificate valid?        (Layer 6)
   HTTP GET /healthz         -> does the app respond correctly?  (Layer 7)
   full user journey         -> does log in -> search -> buy work?
   DNS lookup                -> does the name still resolve?
```

**Probe from multiple places.** A check running in the same data centre as your
servers will not notice that a routing problem has made you unreachable from
Europe.

### Practice (25 min)

**A. SNMP against your own router** (if it supports it — most home routers do,
sometimes needing enabling in the admin page):

```bash
sudo apt install snmp   # or: brew install net-snmp
snmpwalk -v2c -c public 192.168.1.1 1.3.6.1.2.1.1        # system info
snmpwalk -v2c -c public 192.168.1.1 IF-MIB::ifDescr      # interface names
snmpget  -v2c -c public 192.168.1.1 IF-MIB::ifHCInOctets.2
```

To turn a counter into a rate: take two readings, subtract, divide by the seconds
between them, multiply by 8 for bits/second.

**B. Generate and read your own flow records:**

```bash
sudo apt install softflowd nfdump
sudo softflowd -i eth0 -n 127.0.0.1:9995 -v 9
nfcapd -w -D -l /tmp/flows -p 9995 &
# ... browse the web for two minutes ...
nfdump -R /tmp/flows -s srcip/bytes -n 10       # your top talkers
nfdump -R /tmp/flows -s dstport/bytes -n 10     # top ports
```

**You just built the thing that answers "what's using my bandwidth?"**

**C. Write a layered probe script:**

```bash
#!/bin/bash
H=${1:-example.com}
echo -n "DNS:  "; dig +short +time=2 $H | head -1 || echo FAIL
echo -n "PING: "; ping -c2 -W2 $H >/dev/null 2>&1 && echo ok || echo FAIL
echo -n "TCP:  "; nc -zv -w3 $H 443 >/dev/null 2>&1 && echo ok || echo FAIL
echo -n "TLS:  "; echo | openssl s_client -connect $H:443 -servername $H 2>/dev/null \
                   | openssl x509 -noout -enddate || echo FAIL
echo -n "HTTP: "; curl -s -o /dev/null -w '%{http_code} in %{time_total}s\n' https://$H
```

Run it against a few sites. **This is a monitoring system in twelve lines** — the
real ones just do it continuously and store the results.

### Build it in Go (45 min) — a blackbox prober with a latency histogram

"Active probes" from this chapter are what Prometheus's `blackbox_exporter`
does. Writing a small one shows you why monitoring systems store
**histograms**, not averages (Chapter 45): a histogram is a set of cumulative
counters, cheap to store and aggregate, from which percentiles can be
estimated later.

```go
// prober: a tiny blackbox exporter (Chapters 45-47). It probes TCP targets on
// a schedule and exposes Prometheus-format metrics -- a latency histogram,
// success/failure counters -- on /metrics.
//
//	go run ./prober -listen :9115 example.com:443 1.1.1.1:53 10.255.255.1:80
//	curl -s localhost:9115/metrics
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Bucket upper bounds in seconds. Choose them around your SLO, not at random.
var buckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5}

type series struct {
	counts   []uint64 // cumulative-later; raw per-bucket here
	sum      float64
	total    uint64
	failures uint64
}

type registry struct {
	mu sync.Mutex
	m  map[string]*series
}

func (r *registry) observe(target string, d time.Duration, ok bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.m[target]
	if s == nil {
		s = &series{counts: make([]uint64, len(buckets)+1)}
		r.m[target] = s
	}
	if !ok {
		s.failures++
		return
	}
	sec := d.Seconds()
	i := sort.SearchFloat64s(buckets, sec) // first bucket >= sec; len(buckets) = +Inf
	s.counts[i]++
	s.sum += sec
	s.total++
}

// ServeHTTP writes the Prometheus text exposition format by hand -- it is
// just lines of "name{labels} value".
func (r *registry) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# HELP probe_tcp_connect_seconds Time to complete the TCP handshake.")
	fmt.Fprintln(w, "# TYPE probe_tcp_connect_seconds histogram")
	targets := make([]string, 0, len(r.m))
	for t := range r.m {
		targets = append(targets, t)
	}
	sort.Strings(targets)
	for _, t := range targets {
		s := r.m[t]
		var cum uint64
		for i, le := range buckets {
			cum += s.counts[i]
			fmt.Fprintf(w, "probe_tcp_connect_seconds_bucket{target=%q,le=\"%g\"} %d\n", t, le, cum)
		}
		cum += s.counts[len(buckets)]
		fmt.Fprintf(w, "probe_tcp_connect_seconds_bucket{target=%q,le=\"+Inf\"} %d\n", t, cum)
		fmt.Fprintf(w, "probe_tcp_connect_seconds_sum{target=%q} %g\n", t, s.sum)
		fmt.Fprintf(w, "probe_tcp_connect_seconds_count{target=%q} %d\n", t, s.total)
	}
	fmt.Fprintln(w, "# TYPE probe_failures_total counter")
	for _, t := range targets {
		fmt.Fprintf(w, "probe_failures_total{target=%q} %d\n", t, r.m[t].failures)
	}
}

func main() {
	listen := flag.String("listen", ":9115", "address for /metrics")
	every := flag.Duration("every", 5*time.Second, "probe interval")
	timeout := flag.Duration("timeout", 2*time.Second, "per-probe timeout")
	flag.Parse()

	reg := &registry{m: map[string]*series{}}
	for _, target := range flag.Args() {
		go func() {
			for {
				start := time.Now()
				c, err := net.DialTimeout("tcp", target, *timeout)
				if err == nil {
					c.Close()
				}
				reg.observe(target, time.Since(start), err == nil)
				time.Sleep(*every)
			}
		}()
	}
	http.Handle("/metrics", reg)
	log.Printf("probing %v every %v; metrics on http://localhost%s/metrics", flag.Args(), *every, *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
```

```text
$ go run ./prober -every 1s example.com:443 1.1.1.1:53 10.255.255.1:80 &
$ curl -s localhost:9115/metrics | grep -E 'example.com|failures'
probe_tcp_connect_seconds_bucket{target="example.com:443",le="0.05"} 5
probe_tcp_connect_seconds_bucket{target="example.com:443",le="+Inf"} 6
probe_tcp_connect_seconds_count{target="example.com:443"} 6
probe_failures_total{target="1.1.1.1:53"} 0
probe_failures_total{target="10.255.255.1:80"} 2
probe_failures_total{target="example.com:443"} 0
```

**What to notice:**

- **Buckets are cumulative.** `le="0.05"} 5` means 5 of the 6 probes finished
  in 50 ms or less. To find the p99 you look for the first bucket that holds
  99% of `_count`, which is exactly what PromQL's `histogram_quantile()` does.
- **Choose bucket edges around your objective.** If the SLO is "connect under
  100 ms", you need edges just below and above 0.1. Otherwise the histogram
  can't tell you whether you're meeting it.
- **Failures are a separate counter.** Mixing timeouts into the latency
  histogram would hide them in the `+Inf` bucket. Alert on the failure rate
  (Chapter 47) and track latency separately.
- **The exposition format is plain text.** No client library is needed to
  *expose* metrics, though in production you'd use
  `github.com/prometheus/client_golang` (see the
  [Go guide's observability chapter](../../Golang/real-life-golang-guide.md#65-observability-in-go-services-logs-metrics-traces-health-profiling)).

**Exercises:**

1. Point a local Prometheus at it (`scrape_configs: [{job_name: probe, static_configs: [{targets: ['localhost:9115']}]}]`)
   and graph `histogram_quantile(0.99, rate(probe_tcp_connect_seconds_bucket[5m]))`.
2. Add a **TLS** probe type that records handshake time separately from TCP
   connect time, and an **HTTP** probe that checks the status code. Now you
   can tell which layer got slow.
3. Write the Chapter 47 alert for it: "more than 10% of probes to a target
   failed for 5 minutes", plus a runbook link.

### Check yourself

1. What question do counters answer that flow records can't, and vice versa?
2. Why must you use the 64-bit interface counters?
3. Why does a 60-second polling interval hide microbursts?
4. Why do you need active probes as well as passive monitoring?

### Further reading

- **Docs:** Prometheus `blackbox_exporter` — does ICMP, TCP, HTTP, DNS, and TLS
  certificate-expiry probing in one binary. Read its config examples.
- **Docs:** Prometheus `node_exporter` and `snmp_exporter`.
- **Tool:** `ntopng` — a free web UI that turns flow data into readable graphs.
  Excellent for learning what flow records can tell you.
- **Article:** "SNMP explained" — search for a recent practical guide; the
  concepts (OID, MIB, community, trap) are stable.

---

## Chapter 47 — Alerts that wake you for the right reasons

### In one sentence

Alert on **symptoms users feel**, make every page actionable, and suppress the
alerts that are just consequences of another alert.

### The problem

The natural instinct is to alert on everything you can measure. Do that and you
get **alert fatigue**: hundreds of notifications, all ignored, including the one
that mattered.

### The rules

**1. Alert on symptoms, not causes.**

```
   GOOD:  "p99 latency above 2 s for 5 minutes"       users are suffering
   GOOD:  "error rate above 1% for 3 minutes"
   GOOD:  "the site is unreachable from 2 of 5 regions"

   POOR:  "CPU above 80%"           So what? Is anything actually broken?
   POOR:  "interface at 75%"        Is this worth waking someone at 3am?
   POOR:  "disk at 70%"             That's a ticket, not a page.
```

High CPU with happy users is not an incident. It might be worth a dashboard or a
ticket — not a page.

**2. Every page must be actionable.**
If the person woken up can't do anything about it right now, it should not be a
page. Make it a ticket or leave it on a dashboard.

**3. Use duration, not instants.**
"Above threshold **for 5 minutes**" eliminates the vast majority of false alarms
from momentary spikes.

**4. Suppress dependent alerts.**
If the core switch is down, you do not want 400 separate alerts for everything
behind it. Alert on the root cause and summarise the rest.

**5. Distinguish "page" from "ticket".**

```
   PAGE (wake a human now):   users are affected and it won't self-heal
   TICKET (look tomorrow):    a disk will fill in three weeks;
                              a redundant link failed but the backup is carrying it
```

### The alerts everybody should have

Regardless of what you run:

```
   [ ] TLS CERTIFICATE EXPIRY -- warn at 30, 14, and 7 days.
       The most preventable outage in existence, and it still happens weekly
       to large companies.
   [ ] Domain name expiry
   [ ] The service is unreachable FROM OUTSIDE (probe from elsewhere)
   [ ] Error rate above your normal baseline
   [ ] p99 latency above your normal baseline
   [ ] DNS resolution failing
   [ ] Interface errors increasing (this precedes packet loss -- Chapter 35)
   [ ] Disk full on anything that writes logs
   [ ] Backups did not complete
```

### A monitoring stack you could actually build

```
   METRICS      Prometheus + node_exporter (hosts) + snmp_exporter (devices)
   PROBES       blackbox_exporter (ICMP / TCP / HTTP / DNS / TLS expiry)
   DASHBOARDS   Grafana
   ALERTS       Alertmanager -> email / Slack / PagerDuty
   FLOWS        nfdump or ntopng, or your cloud's VPC flow logs
   LOGS         journald/syslog -> Loki, correlated by timestamp
   PACKETS      a tcpdump ring buffer on important hosts (Chapter 31)
```

All of the above is free and open source, and it runs comfortably on one small
server for a modest environment.

### Practice (25 min)

**A. Write the certificate-expiry check** — the single highest-value alert:

```bash
#!/bin/bash
for H in example.com yoursite.com; do
  END=$(echo | openssl s_client -connect $H:443 -servername $H 2>/dev/null \
        | openssl x509 -noout -enddate | cut -d= -f2)
  # macOS uses a different date syntax; try both
  EXP=$(date -d "$END" +%s 2>/dev/null || date -jf '%b %d %T %Y %Z' "$END" +%s)
  DAYS=$(( (EXP - $(date +%s)) / 86400 ))
  printf '%-25s %4d days\n' "$H" "$DAYS"
  [ $DAYS -lt 14 ] && echo "  *** WARNING: renew soon ***"
done
```

Put it in `cron`. You have now prevented a category of outage.

**B. Establish your baselines.** For a service you care about, record over a
normal day: typical latency (p50/p95/p99), typical error rate, typical
throughput. **Without a baseline, no threshold you choose is meaningful.**

**C. Design alerts for the five problems in Chapter 35.** For each — DNS failure,
MTU black hole, firewall drop, packet loss, bufferbloat — write down: what
symptom would you alert on, and what would the alert message say to help the
responder?

### Common confusions

- **"More monitoring is better."** More *signal* is better. More *alerts* is
  usually worse.
- **"We should alert on every error."** Alert on error **rate** relative to
  normal. Some errors are normal.
- **"Monitoring is for big companies."** A cron job that checks your site from
  outside and emails you is monitoring, and it's better than nothing by a huge
  margin.

### Check yourself

1. Why alert on symptoms rather than causes?
2. What's the difference between a page and a ticket?
3. Why require a threshold to be exceeded *for a duration*?
4. Name three alerts every service should have.

### Further reading

- **Book (free):** *Site Reliability Engineering* (Google), chapters on
  "Monitoring Distributed Systems" and "Being On-Call". Also *The Site Reliability
  Workbook* for the practical alerting recipes (multi-window burn rates).
- **Article:** "My Philosophy on Alerting" — Rob Ewaschuk. A short, widely-cited
  memo that shaped modern alerting practice.
- **Docs:** Prometheus "Alerting rules" and Alertmanager documentation.

---

### End of Part 11 — Milestone check

- [ ] I can name the four kinds of monitoring signal
- [ ] I know why application logs miss connection-level failures
- [ ] I can explain why averages hide problems
- [ ] **I have read SNMP counters or kernel counters from a real device**
- [ ] **I have written a layered probe script**
- [ ] I have a certificate-expiry check running

---

# Part 12 — How routers learn routes

Chapter 13 showed you the routing table. This part is about how it gets filled in
without a human typing every line.

## Chapter 48 — Two ways to share directions

### In one sentence

Routers either tell each other **how far away** things are (distance-vector) or
**what they're connected to** (link-state), and everything else follows from that
choice.

### The problem

A router knows the networks it's directly plugged into. Everything else has to be
learned. Typing static routes works for three routers and is impossible for three
hundred — and static routes don't react when a link fails.

### The two approaches

```
   DISTANCE-VECTOR                    LINK-STATE
   "routing by rumour"                "everyone gets the same map"

   Each router tells its NEIGHBOURS   Each router describes ITS OWN LINKS
   how far away everything is.        to EVERY other router.

   A: "I can reach X, it costs 3"     A: "I'm connected to B (cost 1)
   B: "OK, so for me it costs 4"          and C (cost 5)"
   B tells C: "X costs 4"             Everyone builds an IDENTICAL map,
                                      then each calculates its own shortest
                                      paths from that map.

   + simple, low memory               + fast, reliable convergence
   + fine for small networks          + you can see the whole topology
   - slow to react to failures        - more CPU and memory
   - can form temporary LOOPS         - needs structure to scale

   RIP, EIGRP                         OSPF, IS-IS
```

### Why "routing by rumour" is a real problem

```
   A --- B --- C --- X

   C loses its link to X.
   But B is still telling everyone "I can reach X, cost 2"
       (B learned that FROM C, but doesn't remember that).
   C hears B and thinks: "great, I'll go via B, cost 3."
   B hears C's new cost and updates to 4. Then 5. Then 6.

   Meanwhile packets for X bounce between B and C until the
   count reaches "infinity" (16 hops in RIP) and both give up.
```

This is **count-to-infinity**. Distance-vector protocols patch it with rules like
*split horizon* ("never tell a neighbour about a route you learned from them")
and hold-down timers. Link-state protocols simply don't have the problem, because
nobody relies on a neighbour's summary.

### Which protocol goes where

```
   INSIDE one organisation  (an "IGP")     BETWEEN organisations  (an "EGP")
   goal: find the FASTEST path             goal: enforce BUSINESS POLICY

   OSPF   -- most common in enterprises    BGP -- the only one
   IS-IS  -- common in large ISPs
   EIGRP  -- Cisco networks
   RIP    -- teaching and tiny networks
```

That distinction matters more than it sounds: **an IGP optimises for speed; BGP
optimises for money.** That's Chapter 50.

### When routes conflict

If a router learns the same destination from two protocols, it needs a tiebreak.
Most vendors use **administrative distance** — a trust ranking, lower = more
trusted:

```
   Connected      0     <- always wins; it's directly attached
   Static         1     <- a human said so
   eBGP          20
   OSPF         110
   RIP          120
   iBGP         200
```

**Remember from Chapter 13:** longest-prefix match is checked *first*. A more
specific route always wins, regardless of which protocol found it. Administrative
distance only decides between routes of the *same* prefix length.

### Practice (15 min)

```bash
# what protocol installed each of your routes?
ip route show                 # Linux -- look for "proto kernel/static/dhcp/ospf/bgp"
ip route show proto static
netstat -rn                    # macOS

# install a static route and watch it appear (then remove it)
sudo ip route add 192.0.2.0/24 via 192.168.1.1
ip route show 192.0.2.0/24
sudo ip route del 192.0.2.0/24
```

**Optional, and very worthwhile:** install **FRRouting** (`sudo apt install frr`)
on two VMs or containers. It gives you real OSPF and BGP for free, and you can run
every example in the next three chapters.

### Check yourself

1. What's the core difference between distance-vector and link-state?
2. Explain count-to-infinity in your own words.
3. What's the difference between an IGP and an EGP, in terms of goals?
4. Which is checked first: longest-prefix match, or administrative distance?

### Further reading

- **Video:** "Routing Protocols Explained" and the OSPF/BGP series — Practical
  Networking. Ed's routing playlist is the best free introduction.
- **Book:** *Computer Networking: A Top-Down Approach*, Chapter 5.
- **Book:** *Network Warrior* (Donahue) — practical routing chapters written by
  someone who has debugged this at 3am.

---

## Chapter 49 — OSPF: finding your way inside one network

### In one sentence

Every OSPF router floods a description of its own links to every other router,
they all build the identical map, and each calculates its own shortest paths from
it.

### How it works, step by step

```
   1. HELLO      Routers multicast Hello packets to find neighbours.
                 Both sides must agree on: hello timer, dead timer, area ID,
                 authentication, and subnet mask -- or no relationship forms.

   2. ADJACENCY  Neighbours synchronise their databases and become "FULL".

   3. FLOODING   Each router sends an LSA ("Link State Advertisement")
                 describing ITS OWN links and their costs.
                 These get flooded to everyone.

   4. DATABASE   Every router now holds an IDENTICAL copy of the map
                 (the "link-state database").

   5. SPF        Each router runs Dijkstra's shortest-path algorithm
                 on that map, from itself.

   6. INSTALL    The results go into the routing table.

   Any change -> flood the update -> everyone recalculates.
```

OSPF runs directly on IP as protocol 89 — not TCP, not UDP.

### Cost: how OSPF picks a path

```
   cost = reference bandwidth / link bandwidth
   (lower cost = preferred)
```

And here is the single most important practical fact about OSPF:

> **The default reference bandwidth is 100 Mbit/s. So every link at or above
> 100 Mbit/s gets cost 1.**

```
   100 Mbit link   -> cost 1
   1 Gbit link     -> cost 1     <-- identical!
   10 Gbit link    -> cost 1     <-- identical!
   100 Gbit link   -> cost 1     <-- identical!
```

Out of the box, OSPF **cannot tell a gigabit link from a 100-gigabit link** and
will choose by hop count instead. The fix is one command
(`auto-cost reference-bandwidth 100000`) — but it must be set **identically on
every router**, or you get inconsistent routing.

**Check this on every OSPF network you inherit.** It is a very common real-world
misconfiguration.

### Areas: how OSPF scales

Flooding every link change to every router doesn't scale past a few hundred
routers. OSPF divides the network into **areas**:

```
              +-------------- AREA 0 (the backbone) --------------+
              |                                                   |
        [ABR] |                                                   | [ABR]
              |                                                   |
          AREA 1                                              AREA 2

   * SPF runs only WITHIN an area -- a change in Area 1 doesn't make
     Area 2 recalculate.
   * ALL areas must connect to Area 0. Traffic between areas goes through it.
   * An ABR (Area Border Router) sits in two areas and summarises between them.
```

You can also make an area **stub** (block external routes, inject a default route
instead) to keep its routing tables tiny.

### Neighbour states, and what "stuck" means

You'll see these in `show ip ospf neighbor`. The healthy end state is **FULL**.

| Stuck at | Almost always means |
|---|---|
| `Init` | Hellos are only going one way — an ACL or a one-way link |
| `2-Way` | **Normal** on a shared LAN between two non-designated routers |
| `ExStart` | **MTU mismatch** between the two routers (Chapter 15 again!) |
| Flapping | Timer mismatch, or a genuinely unstable link |

**`ExStart` = check your MTUs** is a genuinely useful piece of knowledge.

### Practice (30 min)

If you installed FRRouting on two machines:

```bash
sudo vtysh
  configure terminal
   router ospf
    ospf router-id 1.1.1.1
    auto-cost reference-bandwidth 100000
    network 10.0.0.0/24 area 0
   exit
  end

  show ip ospf neighbor        # want to see state FULL
  show ip ospf interface eth0  # cost, timers, area, network type
  show ip ospf database        # the map every router shares
  show ip route ospf           # what OSPF installed
```

Then **break it deliberately** (the best way to learn):
1. Change the hello timer on one side only → the adjacency drops. Watch
   `show ip ospf neighbor`.
2. Set different MTUs on each side → watch it stick at `ExStart`.
3. Unplug one of two paths → watch OSPF reconverge and the route change.

### Check yourself

1. What does each OSPF router flood, and what does every router end up with?
2. Why does the default reference bandwidth cause problems?
3. What must all areas connect to?
4. An adjacency is stuck at `ExStart`. What do you check?

### Further reading

- **Video:** "OSPF Explained" series — Practical Networking. Thorough and clear.
- **Book:** *Routing TCP/IP, Volume 1* — Jeff Doyle. The classic reference.
- **Book:** any CCNA study guide — OSPF is a large part of the syllabus and the
  material is drill-oriented.
- **Lab:** FRRouting, or Cisco Packet Tracer / GNS3 for a full virtual network.

---

## Chapter 50 — BGP: how the whole internet agrees on directions

### In one sentence

Around 75,000 independent networks tell each other which addresses they can reach
and through which path — and they choose between paths based on **business
relationships**, not distance.

### The problem BGP solves

OSPF works inside one organisation where everyone cooperates and shares a goal.
Between organisations, that's false: your ISP doesn't want to carry your
competitor's traffic for free, and nobody wants to publish their internal
topology.

### Autonomous systems

```
   An AUTONOMOUS SYSTEM (AS) is one network under one administrative control,
   identified by a number (an ASN).

   AS15169 = Google        AS13335 = Cloudflare        AS3356 = Lumen

   BGP's job: let ASes tell each other which IP ranges they can reach.
```

BGP runs over **TCP port 179** — so, unlike OSPF, the two routers must already be
able to reach each other before BGP can start.

### How it advertises

```
   "I can reach 93.184.216.0/24, and the path is: AS64500 AS3356 AS15133"
                \______________/                 \___________________/
                  the prefix                       the AS_PATH -- the list
                                                   of networks to traverse
```

The **AS_PATH** does two jobs:
1. **Loop prevention** — if a router sees its own ASN in the path, it rejects the
   route. Simple and effective.
2. **A tiebreak** — a shorter path is *usually* preferred.

### The crucial insight: policy beats distance

When BGP has several paths to the same prefix, it walks a list of tiebreaks. The
first four are what matter:

```
   1. Highest WEIGHT          (local to one router, vendor-specific)
   2. Highest LOCAL_PREF      <-- THE knob. Set by YOUR network's policy.
   3. Locally originated
   4. Shortest AS_PATH        <-- the famous one... but it's only FOURTH
   ... (several more, rarely reached)
```

**`LOCAL_PREF` beats `AS_PATH`.** In other words, a network will happily send your
traffic along a longer path because that path is cheaper or contractually
preferred.

> **This is why your packets sometimes cross a continent to reach a server two
> miles away.** The internet's shape is commercial, not geographic.

### Business relationships

```
   CUSTOMER   pays you        -> you advertise them to EVERYONE
   PEER       free exchange   -> you exchange only your own and your
                                 customers' routes
   PROVIDER   you pay them    -> they give you everything;
                                 you advertise only yourself and your customers
```

These relationships, encoded as routing policy, determine the actual path your
packets take.

### Influencing traffic (and the asymmetry)

```
   OUTBOUND (traffic you send)      -> YOU control it easily.
        Set LOCAL_PREF higher on routes from your preferred provider.

   INBOUND (traffic sent to you)    -> you can only HINT.
        * AS-path prepending: advertise with your own ASN repeated,
          making the path look longer and less attractive
        * ask your provider for a "community" tag that adjusts their policy
        * advertise a MORE SPECIFIC prefix out one link
          (longest-prefix match beats every BGP attribute)
```

**You cannot force inbound path selection**, because the remote network's own
`LOCAL_PREF` is evaluated before your prepending. This asymmetry frustrates
everyone who multihomes for the first time.

### The three ways BGP breaks the internet

**1. Prefix hijack.** An AS announces address space it doesn't own. Because
**longest-prefix match wins** (Chapter 13), announcing a `/24` inside someone's
`/16` pulls their traffic globally. This has taken down YouTube (2008) and been
used to steal cryptocurrency by hijacking DNS providers.

**2. Route leak.** An AS accidentally re-announces routes it shouldn't — becoming
an unintended transit for traffic it can't carry. Usually a filter mistake;
usually causes congestion and blackholes rather than theft.

**3. Instability.** A flapping link generates updates that ripple globally.

### The defences

| Control | What it does |
|---|---|
| **RPKI + route validation** | A cryptographic record ("ROA") states which AS may announce which prefix. Routers reject invalid announcements. **The most important modern defence — publish ROAs for anything you own.** |
| **Prefix filters** | Only accept prefixes a customer is registered to originate |
| **Max-prefix limits** | Automatically shut a session that suddenly sends 500,000 routes instead of 500. **The cheapest insurance in networking.** |
| **MANRS** | An industry programme bundling these practices |

### Practice (20 min)

You don't need to run BGP to explore it — the global table is public.

```bash
# who announces this prefix, and through which path?
curl -s "https://stat.ripe.net/data/prefix-overview/data.json?resource=1.1.1.0/24" \
  | python3 -m json.tool | head -30

# what prefixes does an AS announce?
whois -h whois.radb.net -- '-i origin AS13335' | grep -c '^route:'

# see the AS path your own traffic takes
traceroute -A 1.1.1.1          # -A shows the ASN at each hop
mtr -z 8.8.8.8                  # -z also shows ASNs

# is a prefix RPKI-valid?
curl -s "https://stat.ripe.net/data/rpki-validation/data.json?resource=AS13335&prefix=1.1.1.0/24" \
  | python3 -m json.tool | grep -i status
```

Also explore **bgp.he.net** in a browser: look up any company and see its ASN,
its prefixes, and who it peers with. It's genuinely interesting.

### Check yourself

1. What is an autonomous system?
2. What two jobs does the AS_PATH do?
3. Why does `LOCAL_PREF` beating `AS_PATH` matter?
4. Why is a prefix hijack effective, in terms of Chapter 13's rules?
5. What is RPKI for?

### Further reading

- **Article:** "What is BGP?" and "BGP hijacking" — Cloudflare Learning Center.
  Their incident post-mortems are also excellent reading.
- **Video:** "BGP Explained" — Practical Networking or NetworkChuck.
- **Site:** `bgp.he.net` — explore the real global routing table.
- **Book:** *BGP* (Iljitsch van Beijnum) or *Internet Routing Architectures* —
  when you're ready for depth.
- **News:** follow BGP incident write-ups; they're the best teaching material
  available and they happen regularly.

---

## Chapter 51 — When a router dies: VRRP and friends

### In one sentence

Your devices have exactly one default gateway, so two routers share a fake IP
address between them and one takes over instantly if the other fails.

### The problem

Every device on your network has one default gateway (Chapter 13). If that router
dies, **everything is cut off** — even if there's a perfectly good second router
sitting next to it, because nobody is configured to use it.

You could reconfigure every device. You could rely on DHCP to hand out a new
gateway (slow, and only on lease renewal). Neither is acceptable.

### The idea

Two routers pretend to be **one** router.

```
   Router A  (priority 110)  \
                               >---  VIRTUAL IP:  192.168.1.1
   Router B  (priority 100)  /       VIRTUAL MAC: 00:00:5e:00:01:01

   Every device on the LAN uses 192.168.1.1 as its gateway.
   They have no idea two routers exist.

   Router A is MASTER: it answers ARP for the virtual IP and forwards traffic.
   Router B is BACKUP: it sits silently, listening.

   A sends a "still alive" message every second.
   If B misses three of them (~3 seconds), B takes over:
       - claims the virtual IP and virtual MAC
       - sends a gratuitous ARP so switches update their tables (Chapter 6!)
       - starts forwarding

   Devices notice nothing except a ~3 second pause.
```

Because the **MAC address moves too**, devices don't even need to re-ARP — their
existing ARP cache entry stays valid.

### The protocols

| Name | Notes |
|---|---|
| **VRRP** | The open standard (RFC 5798). Use this. |
| **HSRP** | Cisco's proprietary equivalent. Same idea, different names (Active/Standby). |
| **GLBP** | Cisco; also shares the *load* between routers, not just failover |
| **keepalived** | The common Linux implementation of VRRP — also widely used for service VIPs, not just routers |

### Two details that matter in practice

**1. Track your uplink.** A router that's alive but has lost its own internet
connection is worse than useless — it'll happily be Master and blackhole
everything. **Interface tracking** lowers its priority when the uplink fails, so
the other router takes over.

**2. Preemption.** When the failed router comes back, should it take over again?
Usually yes — but with a **delay**, so it has time to rebuild its routing tables
first. Otherwise it becomes Master while it still has no routes.

### Where you'll meet this

- Any office or data centre with redundant routers or firewalls
- Load balancers in a highly-available pair (`keepalived` + HAProxy is a classic)
- Your cloud provider does the equivalent invisibly — the "virtual router" in a
  VPC is already redundant
- Modern data centres increasingly use **anycast gateways** (every switch answers
  for the same gateway address) instead, but VRRP is still everywhere

### Practice (15 min)

```bash
# see VRRP on the wire if your network uses it
sudo tcpdump -i any -n 'vrrp or proto 112'

# a virtual MAC always starts 00:00:5e:00:01:xx -- check your ARP cache
ip neigh | grep -i '00:00:5e:00:01'
arp -a | grep -i '00:00:5e'

# what IS your gateway's MAC? Is it a virtual one?
ip neigh show $(ip route | awk '/default/{print $3; exit}')
```

If your gateway's MAC starts `00:00:5e:00:01`, you're behind a VRRP pair right
now.

**Thought exercise for `notes.md`:** your gateway fails over. Which of these
survive, and which break?
1. An in-progress file download (TCP)
2. A video call (UDP)
3. A ping
4. An SSH session

*(Answer: all of them survive if failover is fast enough. TCP retransmits the
missing packets — Chapter 22 — and the SSH session continues. The video call
drops a few hundred milliseconds of audio. This is precisely why TCP's
reliability exists.)*

### Check yourself

1. What problem does VRRP solve?
2. What two things do the routers share?
3. Why does the virtual MAC matter, not just the virtual IP?
4. Why should a router track its uplink interface?

### Further reading

- **Video:** "VRRP / HSRP Explained" — Practical Networking.
- **Docs:** `keepalived` documentation — the practical Linux implementation, with
  configuration examples you can run.
- **RFC 5798** — the VRRP spec; the state machine section is readable.

---

### End of Part 12 — Milestone check

- [ ] I can explain distance-vector vs link-state
- [ ] I know why the OSPF default reference bandwidth is a problem
- [ ] I know that `ExStart` means "check the MTU"
- [ ] I can explain why `LOCAL_PREF` beating `AS_PATH` shapes the internet
- [ ] I can explain a BGP prefix hijack using longest-prefix match
- [ ] I understand how VRRP makes a gateway failure invisible

---

# Part 13 — Professional networking: operations, debugging, and design

> Parts 1–12 teach the protocols. This part teaches the work: operating
> networks safely, collecting evidence during incidents, changing designs
> without breaking users, and debugging failures that cross hosts, firewalls,
> cloud platforms, containers, DNS, TLS, and load balancers.

## Chapter 52 — Host networking internals: namespaces, veth, bridges, conntrack

### In one sentence

Most "cloud/container networking" is Linux building blocks composed together:
interfaces, routes, namespaces, virtual Ethernet pairs, bridges, NAT,
connection tracking, and the queues between the NIC and your program.

### The problem

A packet reaches a server and the application never sees it. Or the
application is "up" but clients time out. Or a node starts dropping
connections at exactly the same load every day. These failures happen
*inside* the host, after the network has delivered the packet. To debug them
you need a map of the host's own packet path.

### The mental model

On one Linux host:

- a **network namespace** is a separate copy of interfaces, routes, ARP/ND
  tables, firewall rules, and sockets. Every container has one;
- a **veth pair** is a virtual cable with two ends. Put one end in a
  namespace and the other on the host;
- a **bridge** is a software switch (Chapter 7) that connects many veth ends;
- **netfilter** (driven by `nftables` or `iptables`) is a set of hooks where
  rules can accept, drop, reject, or rewrite packets;
- **conntrack** remembers every flow so NAT and stateful firewalling work;
- **queues** sit between the hardware and your code, and every one of them
  can overflow.

### How it actually works: the receive path

```
 NIC ── DMA into RX ring buffer ── IRQ → NAPI poll (softirq)
                                          │  GRO merges segments; RSS/RPS spread flows across CPUs
                                          ▼
                            ┌──── netfilter PREROUTING ────┐   conntrack lookup, DNAT (port-forwards,
                            │  (raw → conntrack → mangle → nat) │   kube-proxy Service IPs)
                            └──────────────┬───────────────┘
                              routing decision: for me, or forward?
                     ┌─────────────────────┴──────────────────────┐
                     ▼ local                                      ▼ forward (router, container host)
               INPUT hook (filter)                          FORWARD hook (filter)
                     │                                            │
           TCP: SYN queue → handshake → ACCEPT QUEUE              POSTROUTING (SNAT/MASQUERADE)
                     │                                            │
           accept() → socket RECEIVE BUFFER → read()              out the other interface
```

Every box is somewhere a packet can disappear, and each has a counter:

| Where | Symptom when it overflows | Where to look |
|---|---|---|
| NIC RX ring | `rx_missed_errors`, `rx_no_buffer` | `ethtool -S eth0` |
| softirq backlog | drops under very high packet rates | `/proc/net/softnet_stat` (2nd column) |
| conntrack table | `nf_conntrack: table full, dropping packet` | `dmesg`, `conntrack -S` |
| SYN queue | SYN floods; SYN cookies kick in | `nstat TcpExtSyncookiesSent` |
| **accept queue** | slow connects (1 s, 3 s retries) while the app looks idle | `ss -ltn` (Recv-Q vs Send-Q), `nstat TcpExtListenOverflows` |
| socket receive buffer | zero window (Chapter 24) | `ss -tm`, `nstat TcpExtTCPRcvQDrop` |

### Practice (30 min): build a container network by hand

Run this only on a Linux lab machine or VM. You're building what Docker's
default bridge network builds:

```bash
# 1. A namespace ("the container") and a virtual cable
sudo ip netns add app
sudo ip link add veth-host type veth peer name veth-app
sudo ip link set veth-app netns app

# 2. A bridge ("docker0") on the host, with the host end plugged in
sudo ip link add br0 type bridge
sudo ip addr add 10.200.0.1/24 dev br0
sudo ip link set br0 up
sudo ip link set veth-host master br0
sudo ip link set veth-host up

# 3. Inside the namespace: an address, loopback, and a default route
sudo ip netns exec app ip addr add 10.200.0.2/24 dev veth-app
sudo ip netns exec app ip link set lo up
sudo ip netns exec app ip link set veth-app up
sudo ip netns exec app ip route add default via 10.200.0.1

ping -c 2 10.200.0.2                          # host -> "container" works

# 4. Let the namespace reach the internet: forwarding + NAT (MASQUERADE)
sudo sysctl -w net.ipv4.ip_forward=1
sudo iptables -t nat -A POSTROUTING -s 10.200.0.0/24 ! -o br0 -j MASQUERADE
sudo ip netns exec app ping -c 2 1.1.1.1      # works only after step 4

# 5. Watch conntrack record the NAT translation
sudo conntrack -L -p icmp 2>/dev/null | head
```

Clean up:

```bash
sudo ip netns del app
sudo ip link del br0
sudo iptables -t nat -D POSTROUTING -s 10.200.0.0/24 ! -o br0 -j MASQUERADE
```

Step 5 prints a line like this (from a real run):

```text
icmp 1 29 src=10.200.0.2 dst=1.1.1.1 type=8 ... src=1.1.1.1 dst=172.17.0.2 type=0 ...
```

Read it as two halves. The **original** direction (`src=10.200.0.2
dst=1.1.1.1`) is what the namespace sent. The **reply** direction
(`src=1.1.1.1 dst=172.17.0.2`) is what conntrack expects back: replies come
to the *host's* address, and conntrack rewrites them back to `10.200.0.2`.
That one table entry is the whole of NAT.

Step 4 is the one people forget. Without IP forwarding the host won't route
the namespace's packets, and without MASQUERADE the replies have nowhere to
go: `10.200.0.2` is a private address the internet can't route back to
(Chapter 12).

### What to inspect during real incidents

```bash
ip -br addr; ip route; ip neigh          # addresses, routes, ARP
ss -s                                    # socket totals by state
ss -ltn                                  # listeners: Recv-Q vs Send-Q = accept queue
sudo nft list ruleset                    # firewall + NAT rules (or: iptables-save)
sudo conntrack -S                        # per-CPU conntrack stats: look at "drop", "insert_failed"
sysctl net.netfilter.nf_conntrack_count net.netfilter.nf_conntrack_max
nstat -az | grep -Ei 'overflow|drop|prune|collapse'
```

**Fingerprint:** if `ss` shows the service listening on `127.0.0.1` only,
no firewall or route change will make it reachable from another host. The
program is bound to loopback (Chapter 19's `bindlab`).

### Build it in Go (20 min) — overflow an accept queue on purpose

The accept queue is the least understood queue in this list. The kernel
completes handshakes **without your program's involvement** and parks
finished connections until the program calls `Accept()`. A program that stops
accepting (deadlocked, all workers busy, a long GC pause in older runtimes)
keeps "accepting" connections until the queue is full. After that, new SYNs
are silently dropped and clients see 1-second, 3-second connect delays while
your dashboards show an idle process.

```go
// acceptq: fill a listening socket's accept queue on purpose.
//
// The kernel completes TCP handshakes on its own and parks finished
// connections in the accept queue until the program calls Accept(). A program
// that stops calling Accept() (deadlocked, GC-paused, all workers busy) still
// "accepts" connections -- until the queue is full. Then new SYNs are dropped
// and clients see slow connects (SYN retries at 1 s, 3 s, ...).
//
//	go run ./acceptq
//	# Linux, in another terminal:  ss -ltn 'sport = :7400'    (Recv-Q = queued, Send-Q = queue size)
//	#                               nstat -az TcpExtListenOverflows TcpExtListenDrops
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:7400") // Go asks for backlog = net.core.somaxconn
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("listening, but NEVER calling Accept()")

	var conns []net.Conn
	for i := 1; i <= 5000; i++ {
		start := time.Now()
		c, err := net.DialTimeout("tcp", ln.Addr().String(), 2*time.Second)
		took := time.Since(start)
		if err != nil {
			fmt.Printf("connection %d: %v after %v  <- queue full, SYN dropped\n", i, err, took.Round(time.Millisecond))
			break
		}
		conns = append(conns, c)
		if took > 500*time.Millisecond {
			fmt.Printf("connection %d took %v  <- SYN was dropped and retried\n", i, took.Round(time.Millisecond))
		}
		if i%1000 == 0 {
			fmt.Printf("%d connections completed the handshake without Accept()\n", i)
		}
	}
	fmt.Printf("\n%d connections are waiting in the accept queue.\n", len(conns))
	fmt.Println("Look now with ss/netstat; exiting in 30s.")
	time.Sleep(30 * time.Second)
}
```

Real output on Linux:

```text
$ sysctl net.core.somaxconn
net.core.somaxconn = 4096
$ go run ./acceptq
listening, but NEVER calling Accept()
1000 connections completed the handshake without Accept()
...
4000 connections completed the handshake without Accept()
connection 4098: dial tcp 127.0.0.1:7400: i/o timeout after 2.001s  <- queue full, SYN dropped

$ ss -ltn 'sport = :7400'
State  Recv-Q Send-Q Local Address:Port
LISTEN 4097   4096       127.0.0.1:7400       <- 4097 waiting, queue size 4096
$ nstat -az TcpExtListenOverflows TcpExtListenDrops
TcpExtListenOverflows           2
TcpExtListenDrops               2
```

**What to notice:** for a **LISTEN** socket, `ss` reuses its columns:
`Recv-Q` is the current accept-queue length and `Send-Q` is its maximum. Go
asks for a backlog of `net.core.somaxconn`, so the queue held 4,096 (plus
one, a long-standing Linux off-by-one) before dropping. `ListenOverflows`
is the counter to alert on. It means "my process wasn't calling accept() fast
enough".

**Exercises:**

1. Add a goroutine that calls `Accept()` once every 10 ms. What's the
   steady-state queue length now? (Little's law: arrival rate × time in
   queue.)
2. Lower `somaxconn` to 128 (`sudo sysctl -w net.core.somaxconn=128`), rerun,
   then restore it. Older kernels defaulted to 128, which is why "raise
   somaxconn" appears in every old tuning guide.
3. Read how containers get their own copies of these sysctls in the
   [Linux guide's container chapter](../../os-linux/real-life-os-guide.md#chapter-44-what-a-container-actually-is-namespaces-cgroups-a-filesystem). `somaxconn`
   is per network namespace.

### Common confusions

- **"Conntrack is only for firewalls."** Conntrack is also what makes NAT,
  Docker port publishing, and Kubernetes Services work. A full conntrack table
  breaks all of them at once, even with no firewall rules at all.
- **"iptables and nftables are different firewalls."** Both are front ends to
  the same netfilter hooks in the kernel. Modern distributions translate
  `iptables` commands into nftables rules (`iptables -V` shows `nf_tables`).
- **"A namespace is a container."** A namespace is one isolation feature. A
  container is several namespaces plus cgroups plus a filesystem
  ([Linux guide, Chapter 44](../../os-linux/real-life-os-guide.md#chapter-44-what-a-container-actually-is-namespaces-cgroups-a-filesystem)).

### Check yourself

1. Draw the path of a packet from the NIC to `read()` and name two places it
   can be dropped.
2. Why does a namespace need both IP forwarding and MASQUERADE on the host to
   reach the internet?
3. `ss -ltn` shows `Recv-Q 4097 Send-Q 4096` for a listener. What's wrong,
   and is it the network's fault?
4. What breaks when the conntrack table is full?

### Further reading

- **Book:** *Linux Kernel Networking* (Rami Rosen): the receive path,
  netfilter, and conntrack, chapter by chapter.
- **Article:** "Monitoring and Tuning the Linux Networking Stack: Receiving
  Data" (Packagecloud blog): the definitive walk through the NIC → socket
  path with every counter.
- **Across the series:** [Linux guide, Chapter 69](../../os-linux/real-life-os-guide.md#chapter-69-advanced-linux-networking-namespaces-routing-nftables-and-packet-paths)
  (namespaces, nftables, packet paths) and [Chapter 50](../../os-linux/real-life-os-guide.md#chapter-50-container-networking-how-containers-talk-to-the-world)
  (Docker networking) show these pieces from the operator's side.

---

## Chapter 53 — Firewalls in the real world: host firewalls, cloud security groups, NACLs

### In one sentence

Firewalls are policy checkpoints. Debugging them means finding which
checkpoint made which decision, and operating them means expressing policy so
that it can be reviewed, tested, and rolled back.

### The problem

A packet from a client to a service in the cloud can pass through five or more
independent policy points: the client's corporate firewall, the cloud
provider's network ACL, the instance's security group, the host's `nftables`,
and a Kubernetes NetworkPolicy. Each has different rules, an owner, and logs
(or none). "The firewall is blocking it" isn't a diagnosis until you can say
*which one*.

### The three layers you will meet

| Layer | Examples | Typical behavior |
|---|---|---|
| Host firewall | `nftables`, `iptables`, Windows Defender Firewall, macOS PF | Runs on the machine itself; stateful via conntrack |
| Cloud security group | AWS SG, GCP firewall rules, Azure NSG | Stateful, attached to a VM/NIC/workload; allow-rules only |
| Network ACL / appliance | AWS NACL, corporate firewall, router ACL | Subnet or perimeter level; often **stateless**; ordered allow *and* deny |

Do not say "the firewall" until you know which one.

### How it actually works: stateful vs stateless

A **stateful** firewall uses connection tracking (Chapter 52). Once it has
allowed the first packet of a flow, it recognises every later packet in both
directions as part of an `ESTABLISHED` connection. You write one rule ("allow
TCP 443 inbound"), and replies are allowed automatically.

A **stateless** filter judges every packet on its own. To allow inbound HTTPS
you need *two* rules: inbound to port 443, **and outbound from 443 to the
client's ephemeral port range** (Chapter 19). Forgetting the return rule is the
classic NACL outage: the SYN gets in, the SYN-ACK is dropped on the way out,
and the client sees a timeout.

```
 STATEFUL (security group)            STATELESS (network ACL)
 inbound:  tcp 443 from 0.0.0.0/0     inbound:  100 allow tcp 443        from 0.0.0.0/0
 outbound: (replies allowed           outbound: 100 allow tcp 1024-65535 to   0.0.0.0/0   <- the return path!
            automatically)                       * deny all
```

### A real host firewall, line by line

This is a complete default-deny nftables policy for a web server. It was
syntax-checked and tested in a lab:

```nft
# /etc/nftables.d/server.nft -- a default-deny host firewall for a web server
table inet filter {
    chain input {
        type filter hook input priority 0; policy drop;

        ct state established,related accept   # replies to traffic we allowed
        ct state invalid drop                 # packets that fit no known flow
        iif "lo" accept                       # local processes talking to each other

        meta l4proto { icmp, ipv6-icmp } accept   # ping, PMTU discovery (Chapter 15!)
        tcp dport 22 ip saddr 10.0.0.0/8 accept   # SSH only from the internal network
        tcp dport { 80, 443 } accept              # the actual service

        tcp dport 8080 reject with tcp reset      # fail FAST for a retired port
        counter comment "default drop"            # count what the policy drops
    }
}
```

```bash
sudo nft -c -f server.nft     # check syntax WITHOUT applying: always do this first
sudo nft -f server.nft        # apply atomically: the whole file or nothing
sudo nft list ruleset         # what's actually loaded
```

Things that make this policy production-grade:

- **`policy drop` plus `ct state established,related accept` first.** Default
  deny, with replies allowed cheaply before any other rule runs.
- **ICMP is allowed.** Blocking it breaks Path MTU Discovery (Chapter 15) and
  makes failures much harder to debug.
- **Source restrictions on admin ports.** SSH is open only to the internal
  range, not the internet.
- **A deliberate `reject` for a retired port**, so old clients fail in 1 ms
  instead of hanging.
- **A counter on the default drop**, so you can see whether the policy is
  dropping anything at all.

### Drop vs reject: measured

The same lab, connecting to the server with `nc -z -w 3` from `10.9.0.1`:

```text
port 443       2 ms  succeeded                            <- allowed, listening
port 8080      1 ms  failed: Connection refused           <- "reject with tcp reset"
port 9000   3012 ms  timed out                            <- listening, but policy DROP
port 22        2 ms  failed: Connection refused           <- ALLOWED (10.0.0.0/8), nothing listening
```

- **Reject** sends an error (TCP RST or ICMP unreachable), so clients fail
  fast.
- **Drop** says nothing, so clients retry SYNs until their own timeout.
- **Port 22** is the trap: the firewall *allowed* it, and the refusal came
  from the kernel because no SSH daemon was running. "Refused" proves you
  got through every firewall on the path.

This is why "connection refused" and "connection timed out" are different
fingerprints, not two wordings for the same thing.

### Real troubleshooting workflow

```bash
# 1. Is anything listening, and on which address?
ss -tulpn | grep ':443'

# 2. Can localhost reach it? (bypasses every network firewall)
curl -vk https://127.0.0.1:443/

# 3. Can another host reach the port?
nc -vz server.example.com 443

# 4. Does the packet arrive at the server?
sudo tcpdump -ni any 'tcp port 443'

# 5. Which rule matched? Watch counters change while you retry
sudo nft list ruleset | grep -n counter
sudo nft monitor trace     # after adding: meta nftrace set 1 to a rule in a prerouting chain
```

Interpretation:

| Observation | Meaning |
|---|---|
| Packet never arrives (step 4 silent) | route, upstream firewall, cloud SG/NACL, wrong IP |
| SYN arrives, no SYN-ACK leaves | host firewall, service not listening on that IP, local policy |
| SYN-ACK leaves, client never gets it | stateless ACL missing the return rule, asymmetric routing |
| Handshake completes, app fails later | not a firewall problem; move up the stack |

### Build it in Go (30 min) — firewall policy as executable tests

Firewall rules drift: someone adds a "temporary" rule and forgets it, or a
migration loses a path. The fix is to write the *intent* down as tests and run
them after every change, from every network segment that matters. This
program reads `allow`/`deny` expectations and checks them in parallel:

```go
// policytest: firewall policy as executable tests (Chapter 53).
//
// Write down what SHOULD be reachable and what should NOT, then let the
// program prove it -- from every network segment that matters, before and
// after every firewall change. Exit code 1 on any mismatch, so it fits in CI.
//
//	cat > policy.txt <<'EOF'
//	# expect   target               why
//	allow      example.com:443      public website
//	deny       example.com:22       SSH must not be exposed
//	allow      1.1.1.1:53           DNS resolver
//	EOF
//	go run ./policytest policy.txt
package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

type rule struct {
	line   int
	expect string // "allow" or "deny"
	target string
	why    string
}

type result struct {
	rule
	outcome string // connected | refused | timeout | error
	ok      bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: policytest policy.txt")
		os.Exit(2)
	}
	rules, err := load(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	results := make([]result, len(rules))
	var wg sync.WaitGroup
	for i, r := range rules {
		wg.Add(1)
		go func() { // probes run in parallel: a "deny" waits for a timeout
			defer wg.Done()
			out := probe(r.target, 3*time.Second)
			// "deny" is satisfied by refused OR timeout: either way, no connection.
			ok := (r.expect == "allow") == (out == "connected")
			results[i] = result{r, out, ok}
		}()
	}
	wg.Wait()

	failed := 0
	for _, r := range results {
		mark := "PASS"
		if !r.ok {
			mark, failed = "FAIL", failed+1
		}
		fmt.Printf("%s  line %-3d expect %-5s %-24s got %-9s  %s\n", mark, r.line, r.expect, r.target, r.outcome, r.why)
	}
	fmt.Printf("\n%d rules, %d failed\n", len(results), failed)
	if failed > 0 {
		os.Exit(1)
	}
}

func probe(target string, timeout time.Duration) string {
	c, err := net.DialTimeout("tcp", target, timeout)
	var ne net.Error
	switch {
	case err == nil:
		c.Close()
		return "connected"
	case errors.Is(err, syscall.ECONNREFUSED):
		return "refused" // a REJECT rule, or nothing listening
	case errors.As(err, &ne) && ne.Timeout():
		return "timeout" // a DROP rule, or no route
	default:
		return "error"
	}
}

func load(path string) ([]rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var rules []rule
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 || (fields[0] != "allow" && fields[0] != "deny") {
			return nil, fmt.Errorf("line %d: want 'allow|deny host:port [reason]'", n)
		}
		rules = append(rules, rule{n, fields[0], fields[1], strings.Join(fields[2:], " ")})
	}
	return rules, sc.Err()
}
```

```text
$ cat policy.txt
# expect target why
allow example.com:443 public website
deny example.com:22 SSH must not be exposed
allow 1.1.1.1:53 DNS resolver
allow 127.0.0.1:1 this one should fail

$ go run ./policytest policy.txt
PASS  line 2   expect allow example.com:443          got connected  public website
PASS  line 3   expect deny  example.com:22           got timeout    SSH must not be exposed
PASS  line 4   expect allow 1.1.1.1:53               got connected  DNS resolver
FAIL  line 5   expect allow 127.0.0.1:1              got refused    this one should fail

4 rules, 1 failed
exit status 1
```

**What to notice:** "deny" passes on either *refused* or *timeout*, because
both mean no connection was made. The output still records which one
happened, and that tells you whether something listens behind the rule. The
probes run concurrently because each "deny" waits for its full timeout; done
one after another, 50 rules would take minutes.

**Exercises:**

1. Run the same policy file from three places: your laptop, a VM in the same
   VPC, and a VM in a different subnet. The results *should* differ. That
   difference is your segmentation, verified.
2. Add a `-junit` flag that writes JUnit XML, and run it in CI after every
   Terraform change to security groups.
3. Extend the format to `allow-from <cidr>` assertions by running the tool
   inside network namespaces (Chapter 52) with different source addresses.

### Safety rule

Before changing firewall policy, write down:

- source IP/prefix;
- destination IP;
- protocol/port;
- direction;
- expected rollback command;
- **how you'll keep your own access** (on a remote host, schedule an
  automatic rollback first: `sudo nft -f server.nft; sleep 60 && sudo nft -f previous.nft &`,
  then cancel it once you've confirmed you can still connect).

Most firewall incidents are caused by "temporary" broad rules that nobody
removes, or by narrow rules that forgot one path.

### Common confusions

- **"I allowed it in the security group, so it's allowed."** Every other layer
  (NACL, host firewall, NetworkPolicy) must allow it too. Policy is the
  intersection of all of them.
- **"The firewall must be dropping replies."** With stateful firewalls,
  replies are allowed automatically. If replies vanish, look for a
  *stateless* ACL or asymmetric routing (Chapter 13).
- **"REJECT is less secure than DROP."** Hiding a closed port gains you very
  little: scanners notice drops too. Use reject inside your own networks,
  where fast failure helps your engineers. Drop at the internet edge if
  you want to slow scanners down.

### Check yourself

1. Why does a stateless NACL need an outbound rule for inbound HTTPS to work?
2. A client gets "connection refused" from your server. Which firewalls can
   you rule out?
3. What do `policy drop` and `ct state established,related accept` each do,
   and why does their order matter?
4. Why should firewall intent be written as tests rather than only as rules?

### Further reading

- **Docs:** the nftables wiki (wiki.nftables.org): "Quick reference" and
  "Ruleset debug/tracing".
- **Docs:** AWS VPC documentation, "Compare security groups and network ACLs".
- **Across the series:** the [Linux guide's hardening chapter](../../os-linux/real-life-os-guide.md#chapter-71-linux-security-hardening-selinux-apparmor-auditd-capabilities-and-least-privilege)
  covers host-level least privilege, and the
  [HTTPS guide's security chapter](../../v2-https/real-life-guide-v1.md#chapter-12-security-protecting-the-entire-lifecycle) covers
  application-layer controls (WAF, rate limiting) that sit above these
  firewalls. The Go plan builds a WAF in [Week 10](../../Golang/detailed-90-day-plan/week10.md).

---

## Chapter 54 — VPNs and tunnels: WireGuard, IPsec, GRE, overlays, and MTU traps

### In one sentence

A tunnel puts one packet inside another packet. The price is extra headers,
new routing rules, a smaller effective MTU, and one more place where things
can break silently.

### The problem

You need two private networks to talk across the internet, or containers on
different hosts to share one flat network, or remote laptops to reach internal
services. You can't change the internet's routing, so you **encapsulate**:
wrap your packet in another packet whose addresses the internet *can* route,
and unwrap it at the far end.

### The tunnel family

| Technology | Encapsulation | Encrypted? | What it is good for |
|---|---|---|---|
| WireGuard | IP inside **UDP** (port 51820 by default) | Yes (Noise protocol, ChaCha20-Poly1305) | Simple site-to-site or host VPN, modern crypto, small config |
| IPsec (ESP) | IP inside ESP (IP protocol 50), or UDP 4500 behind NAT | Yes | Enterprise/firewall VPNs, cloud site-to-site VPNs; standard but complex |
| GRE / IP-in-IP | IP inside IP (protocol 47 / 4) | No | Simple encapsulation, often combined with IPsec |
| VXLAN / Geneve | **Ethernet frame** inside UDP (4789 / 6081) | No | Data-centre and container overlays; many virtual L2 networks over one L3 underlay |
| SSH tunnel | TCP inside TCP | Yes | Quick ad-hoc port forwarding, not a network design |

### How it actually works

**WireGuard** in four ideas:

1. Each peer has a **key pair**. Peers are identified only by public key;
   there are no usernames or certificates.
2. **`AllowedIPs` is both routing table and access list.** Outbound: "send
   packets for 10.20.0.0/24 to this peer". Inbound: "accept packets from
   this peer only if their source is in 10.20.0.0/24". The WireGuard paper
   calls this *cryptokey routing*.
3. It's **UDP and connectionless**. A handshake happens on demand and is
   refreshed every two minutes. When a peer's IP address changes (a laptop
   moving from Wi-Fi to 4G), the tunnel follows it automatically.
4. Behind NAT, set **`PersistentKeepalive = 25`** so the NAT mapping
   (Chapter 12) doesn't expire while the tunnel is idle.

**IPsec** has two phases: **IKE** (UDP 500, or 4500 behind NAT) authenticates
the peers and negotiates keys, creating **Security Associations**. **ESP** then
encrypts each packet. *Tunnel mode* wraps the whole original packet (VPNs).
*Transport mode* protects only the payload (host to host). Most "IPsec is
down" incidents are mismatched IKE proposals, expired pre-shared keys, or
UDP 500/4500 blocked somewhere.

**VXLAN** carries a whole Ethernet frame, so the two ends look like they're
on the same switch even if they're in different racks or regions. A 24-bit
**VNI** separates up to 16 million virtual networks. This is how many
Kubernetes CNIs and cloud VPCs build "one flat network" on top of routed
infrastructure ([OSI guide, Layer 2](../real-life-example-osi.md#part-4-layer-2-data-link)).

### The MTU tax

Every tunnel adds headers, and those bytes come out of the inner packet's
space:

| Tunnel (over IPv4, 1500-byte path) | Overhead | Inner MTU |
|---|---|---|
| IP-in-IP | 20 | 1480 |
| GRE | 24 | 1476 |
| VXLAN | 50 | 1450 |
| WireGuard | 60 (IPv4) / 80 (IPv6) | 1420 (default, safe for both) |
| IPsec ESP tunnel mode | ~50–73 (varies with cipher, padding) | ~1400–1438 |
| PPPoE (not a tunnel you chose, but the same effect) | 8 | 1492 |

If the inner packet is still 1500 bytes, something must shrink, fragment, or
fail. Common fingerprints:

- SSH works, file transfer stalls.
- Small HTTP responses work, large ones hang.
- VPN connects, but some sites or APIs fail.
- `ping` works until you add "don't fragment" and a larger size.

**The two fixes:**

1. **Set the inner MTU correctly** on the tunnel interface (WireGuard
   defaults to 1420 for this reason).
2. **Clamp TCP MSS** on the router so TCP never *tries* to send oversized
   segments. This is one rule, and it's the standard fix for PPPoE and
   site-to-site VPNs:

```bash
sudo iptables -t mangle -A FORWARD -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu
```

MSS clamping only helps TCP. UDP protocols (QUIC, DNS, games) still depend on
a correct MTU or working PMTU discovery.

### Practice (30 min): two tunnels between two namespaces

Run on a Linux lab machine or a privileged container
(`docker run --rm -it --privileged debian:bookworm`, then
`apt-get install -y iproute2 iputils-ping wireguard-tools`). Two namespaces,
`a` and `b`, connected by a veth "internet":

```bash
ip netns add a; ip netns add b
ip link add va type veth peer name vb; ip link set va netns a; ip link set vb netns b
ip -n a addr add 192.0.2.1/24 dev va; ip -n a link set va up; ip -n a link set lo up
ip -n b addr add 192.0.2.2/24 dev vb; ip -n b link set vb up; ip -n b link set lo up

# Tunnel 1: IP-in-IP (no encryption)
ip -n a link add tun0 type ipip local 192.0.2.1 remote 192.0.2.2
ip -n b link add tun0 type ipip local 192.0.2.2 remote 192.0.2.1
ip -n a addr add 10.10.0.1/30 dev tun0; ip -n a link set tun0 up
ip -n b addr add 10.10.0.2/30 dev tun0; ip -n b link set tun0 up
ip -n a link show tun0                                # mtu 1480
ip netns exec a ping -c1 -M do -s 1452 10.10.0.2      # 1452 + 8 ICMP + 20 IP = 1480: fits
ip netns exec a ping -c1 -M do -s 1453 10.10.0.2      # one byte more: "message too long"

# Tunnel 2: WireGuard (encrypted)
cd /tmp; wg genkey | tee a.key | wg pubkey > a.pub; wg genkey | tee b.key | wg pubkey > b.pub
ip -n a link add wg0 type wireguard
ip netns exec a wg set wg0 private-key a.key listen-port 51820 \
  peer $(cat b.pub) allowed-ips 10.20.0.2/32 endpoint 192.0.2.2:51820
ip -n b link add wg0 type wireguard
ip netns exec b wg set wg0 private-key b.key listen-port 51820 \
  peer $(cat a.pub) allowed-ips 10.20.0.1/32 endpoint 192.0.2.1:51820
ip -n a addr add 10.20.0.1/24 dev wg0; ip -n a link set wg0 up
ip -n b addr add 10.20.0.2/24 dev wg0; ip -n b link set wg0 up

ip netns exec a ping -c2 10.20.0.2
ip netns exec a wg show wg0                           # latest handshake, transfer counters
ip netns exec a ping -c1 -M do -s 1392 10.20.0.2      # 1392 + 28 = 1420: the WireGuard MTU
```

All of this was run while writing this guide. The MTUs came out at exactly 1480
and 1420, and `wg show` reported `latest handshake: 1 second ago`.

Now capture on the "internet" link while you ping through each tunnel:
`ip netns exec a tcpdump -ni va`. The IP-in-IP packets show the inner ICMP
in plain text. The WireGuard packets are opaque UDP datagrams to port 51820.
That's the difference encryption makes.

### Design advice

- Document the underlay path, overlay prefixes, and MTU for every tunnel.
- Prefer **route-based** VPNs (a tunnel interface plus routes) over
  policy-based ones when you need several networks and failover.
- Monitor tunnel state *and* loss and latency through it, not just "is the
  process running?". For WireGuard, alert when the latest handshake is older
  than 3 minutes on a tunnel that should be busy.
- Keep tunnel routes specific. Default-route VPNs are easy to misuse and
  make every outage look like a VPN outage.
- Test with **real payload sizes** after every change: a large HTTP download
  through the tunnel, not just ping.

### Common confusions

- **"VXLAN is a VPN."** It's an overlay with no encryption. Traffic between
  hosts is readable unless you add IPsec/WireGuard underneath, or mTLS
  above it.
- **"WireGuard is connection-oriented like TCP."** It's UDP with no
  connection state you can see. "Connected" just means a recent handshake.
- **"Fragmentation will sort it out."** Most tunnels set DF on the outer
  packet or can't fragment the inner one, and ICMP "too big" messages are often
  filtered. Assume fragmentation will *not* save you.

### Check yourself

1. What does `AllowedIPs` do in each direction?
2. Why does WireGuard default to an MTU of 1420 instead of 1500?
3. What does MSS clamping fix, and what can't it fix?
4. You can SSH through a new VPN but `git clone` hangs. What's your first
   hypothesis and your first test?

### Further reading

- **Paper:** "WireGuard: Next Generation Kernel Network Tunnel" (Jason
  Donenfeld), short and readable, including cryptokey routing.
- **RFC 7348** (VXLAN), **RFC 4301** (IPsec architecture).
- **Across the series:** the [OSI guide's byte-budget section](../real-life-example-osi.md#part-11-x-ray-of-one-packet)
  totals the overhead layer by layer, and the
  [Linux guide's advanced networking chapter](../../os-linux/real-life-os-guide.md#chapter-69-advanced-linux-networking-namespaces-routing-nftables-and-packet-paths) shows
  where tunnel interfaces fit into the host's routing.

---

## Chapter 55 — DNS operations: authoritative DNS, delegation, split-horizon, outages

### In one sentence

Using DNS and operating DNS are different jobs. Operating it means delegation,
TTL planning, authoritative servers, split views, DNSSEC, and changes made so
that they can be rolled back.

### The problem

DNS changes look trivial (edit one record) and fail in non-obvious ways. The
old answer stays in caches for hours, one resolver disagrees with another,
a typo in a delegation takes a whole domain offline, or a DNSSEC key rollover
turns every lookup into `SERVFAIL`. Several of the internet's biggest outages
were DNS outages.

### Recursive vs authoritative

- **Recursive resolver**: answers clients by walking and caching the DNS tree.
  Examples: your router, a corporate resolver, `1.1.1.1`, `8.8.8.8`.
- **Authoritative server**: owns a zone and publishes its records. It never
  asks anyone else.

Debug both sides separately:

```bash
dig www.example.com                      # via your configured resolver (cached)
dig +trace www.example.com               # walk the tree yourself: root -> .com -> example.com
dig @8.8.8.8 www.example.com             # a specific public resolver
dig @ns1.example.com www.example.com +norecurse   # straight from the source, no cache
```

The `flags:` line tells you who answered: `aa` (authoritative answer) appears
only on answers that come straight from the authoritative server.

### How it actually works: delegation, glue, and the SOA

A domain works because its **parent** points at its nameservers:

```
 .com zone (run by the registry)          example.com zone (run by you or your DNS provider)
 ─────────────────────────────           ─────────────────────────────────────────────────
 example.com.  NS  ns1.example.com.       example.com.  SOA ns1.example.com. hostmaster.example.com.
 example.com.  NS  ns2.example.com.                         2026100301 ; serial
 ns1.example.com. A 192.0.2.53  <- GLUE                    7200       ; refresh (secondaries poll)
 ns2.example.com. A 198.51.100.53 <- GLUE                  900        ; retry
                                                            1209600    ; expire
                                                            300        ; minimum = NEGATIVE-cache TTL
                                          example.com.  NS  ns1.example.com.
                                          www           A   203.0.113.10   ; TTL 300
```

- **Delegation** = NS records in the *parent* zone. You change them at your
  **registrar**, not in your own zone file. If the parent's NS and your
  zone's NS disagree, resolvers behave unpredictably ("lame delegation").
- **Glue** = the A/AAAA records of nameservers that live *inside* the zone
  they serve. Without glue, finding `ns1.example.com` would require asking
  `ns1.example.com`.
- **SOA serial**: secondaries copy the zone only when the serial increases.
  Forgetting to bump it is a classic "my change isn't visible on half the
  servers" bug. The `YYYYMMDDnn` convention makes it hard to forget.
- **Negative caching**: an `NXDOMAIN` is cached too, for the SOA's minimum
  field or the SOA record's TTL, whichever is lower (RFC 2308). If you query
  a name *before* creating it, resolvers remember that it doesn't exist.

### TTL is your rollback timer

If a record has a 1-hour TTL, clients and resolvers can keep the old answer
for up to about an hour. A safe migration timeline:

```
 T-48h  lower TTL 3600 -> 60          (wait at least the OLD TTL for caches to pick it up)
 T-0    change the record              (most clients follow within ~60 s)
 T+1h   verify traffic moved; keep the old target running
 T+24h  raise TTL back to 3600         (cheaper, more resilient to DNS outages)
 T+7d   decommission the old target    (some clients ignore TTLs: JVMs, embedded devices)
```

Long TTLs aren't bad: they make you **resilient** when your DNS provider is
down, because cached answers keep working. Short TTLs make you **agile**.
Choose per record.

### Split-horizon DNS

The same name may intentionally resolve differently depending on where you
ask from:

```text
api.example.com from office VPN  -> 10.10.5.20
api.example.com from internet    -> 203.0.113.20
```

This is useful, but it creates "works for me" incidents. Always record
**which resolver answered**. The same effect happens unintentionally with
GeoDNS and CDNs, which return different addresses by location on purpose.

### DNSSEC in one picture

DNSSEC adds signatures so a validating resolver can prove an answer wasn't
forged:

```
 root (trust anchor) --DS--> .com DNSKEY --DS--> example.com DNSKEY --RRSIG--> www A 203.0.113.10
```

Each parent publishes a **DS** record (a hash of the child's key). When the chain
breaks (an expired signature, or a key rollover where the DS at the
registrar wasn't updated), validating resolvers return **`SERVFAIL`** and the
domain disappears for every user of those resolvers. `dig +dnssec` shows the
`ad` (authenticated data) flag, and `dig +cd` (checking disabled) tells you
whether DNSSEC is the cause: if `+cd` works and a normal query fails, the
problem is the signature chain.

### Build it in Go (25 min) — compare resolvers in one command

Most DNS incidents come down to "different resolvers give different
answers". This program asks several resolvers at once, groups identical
answers, and flags disagreement:

```go
// dnsdiff: ask several resolvers the same question and show where the answers
// disagree -- split-horizon, stale caches, half-finished migrations, and
// "works for me" incidents (Chapter 55).
//
//	go run ./dnsdiff www.github.com
//	go run ./dnsdiff -servers system,1.1.1.1,8.8.8.8,9.9.9.9 api.example.com
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

func main() {
	servers := flag.String("servers", "system,1.1.1.1,8.8.8.8,9.9.9.9", "comma-separated resolvers; 'system' = your OS config")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("usage: dnsdiff [-servers a,b,c] name")
		os.Exit(2)
	}
	name := flag.Arg(0)
	list := strings.Split(*servers, ",")

	type answer struct {
		addrs []string
		took  time.Duration
		err   error
	}
	answers := make([]answer, len(list))
	var wg sync.WaitGroup
	for i, s := range list {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			start := time.Now()
			addrs, err := resolverFor(s).LookupHost(ctx, name)
			slices.Sort(addrs)
			answers[i] = answer{addrs, time.Since(start), err}
		}()
	}
	wg.Wait()

	groups := map[string][]string{} // identical answer sets -> which resolvers gave them
	for i, a := range answers {
		key := strings.Join(a.addrs, " ")
		if a.err != nil {
			key = "ERROR: " + a.err.Error()
		}
		groups[key] = append(groups[key], list[i])
		fmt.Printf("%-10s %-8v %s\n", list[i], a.took.Round(time.Millisecond), key)
	}
	if len(groups) == 1 {
		fmt.Println("\nall resolvers agree")
		return
	}
	fmt.Printf("\n%d different answers. Possible reasons: split-horizon DNS, a record changed\n", len(groups))
	fmt.Println("within the TTL (caches disagree), GeoDNS/CDN steering, or a broken resolver.")
}

// resolverFor returns a Go resolver that sends every query to one server.
// PreferGo uses Go's own DNS client, so the Dial hook decides the server.
func resolverFor(server string) *net.Resolver {
	if server == "system" {
		return net.DefaultResolver
	}
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, net.JoinHostPort(server, "53"))
		},
	}
}
```

```text
$ go run ./dnsdiff www.github.com
system     94ms     20.207.73.82
1.1.1.1    81ms     20.207.73.82
8.8.8.8    80ms     20.207.73.82
9.9.9.9    81ms     20.205.243.166

2 different answers. Possible reasons: split-horizon DNS, a record changed
within the TTL (caches disagree), GeoDNS/CDN steering, or a broken resolver.
```

That disagreement is real and harmless: GitHub steers clients by location, and
Quad9's resolver egresses from a different region. During a migration, the
same output shows you exactly which caches still hold the old address.

**What to notice:** `net.Resolver{PreferGo: true, Dial: ...}` makes Go use its
own DNS client and send every query to the server *you* choose. The same
`Dial` hook is how you'd run DNS over TCP or through a SOCKS proxy.

**Exercises:**

1. Add your company's internal resolver to `-servers` and query an internal
   name from inside and outside the VPN. You've just mapped your split horizon.
2. Add an authoritative mode: look up the zone's NS records, then query each
   authoritative server directly. All of them should return the same SOA
   serial. If they don't, a secondary is behind.
3. During your next DNS change, run it in a loop every 30 s and log when
   each resolver switches. Compare that with the TTL you set.

### DNS incident checklist

- Does the domain delegate to the nameservers you think? (`dig NS example.com @a.gtld-servers.net`)
- Did you query the authoritative servers directly?
- Do all authoritative servers return the same SOA serial?
- Is DNSSEC involved? Does `+cd` change the result?
- Did a CNAME point outside your control (a SaaS provider, a CDN)?
- Did an internal resolver cache a stale answer or a negative answer?
- Are clients using DoH/DoT and bypassing your intended resolver?
- Has the **domain registration** expired? It happens to large companies too.

### Real-world outages

- **October 2016, Dyn:** a Mirai-botnet DDoS on a major DNS provider made
  dozens of large sites unreachable for hours. The lesson: use **two
  independent DNS providers** for critical domains.
- **October 2021, Facebook:** a routing change withdrew the BGP routes to
  Facebook's own authoritative DNS servers. With no route to the nameservers,
  every resolver in the world got `SERVFAIL` for facebook.com, and internal
  tools that depended on the same DNS went down too. The lesson: don't make
  your recovery tools depend on the system that's down.

### Common confusions

- **"DNS propagation takes 48 hours."** There's no propagation. There are
  caches, and they expire on the TTL you set. Plan the TTL and the change is
  predictable.
- **"I can put a CNAME at the zone apex."** The standard forbids a CNAME at
  `example.com` because it would conflict with the SOA and NS records there.
  Providers offer ALIAS/ANAME/"CNAME flattening" to work around it.
- **"SERVFAIL means the server is down."** It means the resolver couldn't
  get a valid answer: upstream unreachable, a lame delegation, or a broken
  DNSSEC chain.

### Check yourself

1. Where do you change a domain's delegation, and why there?
2. What are glue records, and when are they required?
3. You created a record five minutes ago but a resolver still returns
   NXDOMAIN. Why, and how long until it doesn't?
4. Write the TTL timeline for moving `api.example.com` to a new load balancer.
5. How do you tell a DNSSEC failure from a nameserver outage?

### Further reading

- **Book:** *DNS and BIND* (Liu & Albitz), still the best explanation of
  zones, delegation, and operations.
- **Tool:** DNSViz (dnsviz.net) visualises the delegation and DNSSEC chain
  for any domain and highlights what's broken.
- **Across the series:** the [HTTPS guide's DNS chapter](../../v2-https/real-life-guide-v1.md#chapter-2-dns-finding-the-address)
  covers DNS from the browser's side, including the HTTPS/SVCB records that
  advertise HTTP/3, and Chapter 18 here builds a DNS client from raw bytes.

---

## Chapter 56 — TLS and certificate operations: expiry, chains, SNI, rotation

### In one sentence

Most TLS outages aren't broken encryption. They're expired certificates,
wrong names, missing intermediates, SNI mismatches, or clients trusting a
different set of CAs than the one you tested with, and every one of them is
an operational failure.

### The problem

Chapter 29 explained how TLS works. In production, the hard part is the
**inventory**: dozens or hundreds of certificates, on load balancers, CDNs,
ingress controllers, internal services, and mobile apps, each with its own
expiry date, renewal method, and owner. One forgotten certificate takes down
the service that uses it, without warning, on the day it expires.

### How it actually works: the chain

A server sends its **leaf** certificate plus the **intermediates**. The client
builds a path from the leaf to a **root** it already trusts:

```
 leaf: CN=example.com          signed by -> intermediate: "Cloudflare TLS Issuing ECC CA 3"
 intermediate                  signed by -> ... -> root in the CLIENT's trust store
```

Three facts explain most chain failures:

1. **The server must send the intermediates.** Browsers sometimes fetch a
   missing one themselves (AIA fetching) or have it cached, so "works in
   Chrome" proves nothing. `curl`, Java, Go, Python, and Android often don't
   fetch it, and they fail.
2. **Trust stores differ.** The OS, the browser, the JVM (`cacerts`),
   Python's `certifi`, and minimal containers (often *no* CA bundle) each have
   their own list of roots.
3. **Roots expire too.** On 30 September 2021 an old root that Let's Encrypt
   certificates chained to (DST Root CA X3) expired. Up-to-date clients
   switched to the newer root. Old Android versions, old OpenSSL, and
   embedded devices started failing on certificates that were themselves
   perfectly valid.

### Certificate lifetimes are shrinking

In 2025 the CA/Browser Forum voted to cut the maximum lifetime of public TLS
certificates in stages: **200 days** from March 2026, **100 days** from March
2027, and **47 days** from March 2029. At 47 days, manual renewal stops being
possible in practice. Every certificate needs **automated issuance** (ACME:
Let's Encrypt, ZeroSSL, cloud certificate managers) and **automated
deployment** to every place it's used.

### Inspect a real endpoint

```bash
openssl s_client -connect example.com:443 -servername example.com -showcerts </dev/null
echo | openssl s_client -connect example.com:443 -servername example.com 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates -ext subjectAltName
curl -vI https://example.com/ 2>&1 | grep -E 'SSL|subject|expire|issuer'
```

Look for:

- certificate `Not After` date;
- the Subject Alternative Name list containing the hostname (the CN is
  ignored by modern clients);
- issuer and chain: does `-showcerts` print the intermediates?
- TLS version and cipher;
- whether SNI (`-servername`) changes the certificate. Try with and without it.

### Common fingerprints

| Symptom | Likely cause |
|---|---|
| Browser says expired | leaf cert expired, or the client's clock is wrong |
| Works in browser, fails in CLI/app | missing intermediate, or a different CA store |
| Works for some old devices, not others | root/cross-sign expiry; device trust store too old |
| Fails by IP, works by name | certificate names are DNS names; an IP isn't in the SAN |
| Wrong certificate behind load balancer | missing/wrong SNI routing, default certificate served |
| Handshake refused without a certificate | the edge has no certificate for that SNI at all |
| Mobile app fails after rotation | pinned cert/public key not updated |
| Internal service fails after "routine" CA rotation | clients still trust only the old private CA |

### Build it in Go (30 min) — rotate a certificate with zero downtime

Restarting a server to load a new certificate drops connections. Go's
`tls.Config.GetCertificate` hook avoids that: it's called **during every
handshake**, so the server can switch certificates by swapping a pointer.
New connections get the new certificate, and existing connections finish on
the old one.

```go
// certreload: rotate a TLS certificate with zero downtime (Chapter 56).
//
// The server never holds "the" certificate. tls.Config.GetCertificate is
// called during every handshake, so swapping an atomic pointer changes the
// certificate for all NEW connections while existing ones carry on untouched.
//
//	go run ./certreload            # generates two self-signed certs and rotates between them
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type certStore struct {
	cur atomic.Pointer[tls.Certificate]
}

// GetCertificate runs inside every handshake. It can also choose by SNI
// (hello.ServerName) -- that is how one IP serves many certificates.
func (s *certStore) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return s.cur.Load(), nil
}

func main() {
	store := &certStore{}
	store.cur.Store(selfSigned("generation-1"))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Handler:   http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok\n")) }),
		TLSConfig: &tls.Config{GetCertificate: store.GetCertificate, MinVersion: tls.VersionTLS12},
	}
	go srv.ServeTLS(ln, "", "") // empty paths: certificates come from GetCertificate
	addr := ln.Addr().String()

	fmt.Println("serving:", whoServes(addr))
	// In production this would be triggered by a file watcher, SIGHUP, or a
	// secret-manager callback after renewal. Load and VALIDATE first, then swap.
	next := selfSigned("generation-2")
	store.cur.Store(next)
	fmt.Println("rotated. new connections now see:", whoServes(addr))
}

func whoServes(addr string) string {
	// InsecureSkipVerify only because our demo certs are self-signed;
	// never do this in real clients.
	c, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return err.Error()
	}
	defer c.Close()
	leaf := c.ConnectionState().PeerCertificates[0]
	return fmt.Sprintf("CN=%s serial=%s expires=%s", leaf.Subject.CommonName, leaf.SerialNumber, leaf.NotAfter.Format(time.DateOnly))
}

func selfSigned(cn string) *tls.Certificate {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	serial, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: cn},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(90 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	return &tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}
```

```text
$ go run ./certreload
serving: CN=generation-1 serial=2779963337554967105 expires=2027-01-01
rotated. new connections now see: CN=generation-2 serial=2047191123025242675 expires=2027-01-01
```

**What to notice:**

- **`atomic.Pointer`** makes the swap safe while thousands of handshakes run
  concurrently, with no lock on the hot path.
- **Validate before you swap.** In production, load the new pair with
  `tls.LoadX509KeyPair`, parse the leaf, and check that the names, expiry,
  and key match. Only then store it. A bad certificate loaded hot is an
  instant outage.
- **`hello.ServerName`** is the SNI. Serving many domains from one IP means a
  map lookup here, which is what Caddy, Traefik, and Kubernetes ingress
  controllers do.
- Go's `golang.org/x/crypto/acme/autocert` package plugs ACME into this same
  hook: it obtains and renews Let's Encrypt certificates automatically.

**Exercises:**

1. Add a file watcher: reload when `cert.pem`/`key.pem` change on disk (or on
   `SIGHUP`), with validation and a log line showing old and new serial
   numbers.
2. Combine it with Chapter 29's `tlsinspect`: run `tlsinspect -expiry-days 30`
   against every endpoint in a list. That list is your **certificate
   inventory**.
3. Do the Go plan's [zero-downtime rotation (Day 54)](../../Golang/detailed-90-day-plan/week8.md#day-54-tls-certificate-rotation) and
   [mTLS with identity extraction (Day 41)](../../Golang/detailed-90-day-plan/week6.md#day-41-mutual-tls-mtls) next.

### Rotation rule

Certificate rotation is a deployment, not a file copy. Have:

- an **inventory** of every termination point (CDN, LB, ingress, app, mobile
  pins), the names on each certificate, and its owner;
- **expiry monitoring** at least 30 days out (14 when lifetimes reach 47 days),
  alerting from *outside* your network, on the certificate actually served;
- staging validation with old *and* new clients;
- a rollback path (keep the previous certificate and key);
- tests from outside and inside the network.

### Common confusions

- **"We use Let's Encrypt, so renewal is automatic."** Issuance is automatic.
  *Deployment* to the load balancer, the CDN, and the Java keystore may not be.
  Monitor the certificate that's actually served, not the one in the
  directory.
- **"Revoking a certificate stops it from working."** Revocation checking is
  inconsistent. Browsers rely mostly on pushed revocation lists, and in 2025
  Let's Encrypt stopped running OCSP altogether in favour of CRLs. Short
  lifetimes are the real revocation mechanism.
- **"Certificate pinning makes us safer."** It also gives you an outage
  whenever the pinned key changes. If you must pin, pin a CA's or
  intermediate's key, keep a backup pin, and have a remote kill switch.

### Check yourself

1. Why can a site work in a browser but fail from `curl` or a Java service?
2. What does `GetCertificate` let a Go server do that loading a certificate at
   startup can't?
3. What will 47-day certificates force your team to change?
4. List four places a certificate for `api.example.com` might be deployed.

### Further reading

- **Tool:** SSL Labs Server Test (ssllabs.com/ssltest) for chain, protocol,
  and client-compatibility analysis.
- **Tool:** crt.sh, a Certificate Transparency search. Find every certificate
  ever issued for your domains, including ones you didn't know about.
- **Across the series:** the [HTTPS guide's certificate operations chapter](../../v2-https/real-life-guide-v1.md#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation)
  covers ACME, CAA, and CT monitoring, and the
  [Go guide's mTLS chapter](../../Golang/real-life-golang-guide.md#50-mtls-client-certificate-authentication-end-to-end) builds client-certificate
  authentication end to end.

---

## Chapter 57 — Production load balancing: L4 vs L7, health checks, draining, retries

### In one sentence

A load balancer is both a traffic director and a failure amplifier. Good
health checks, draining, and retry policy make a backend failure invisible.
Bad ones turn one slow backend into a global outage.

### The problem

Chapter 36 explained what a load balancer is. Running one raises harder
questions. How does it decide a backend is dead, and how fast? What happens to
in-flight requests during a deploy? Who retries a failed request, and how many
times? And how does the backend find out who the client really was?

### L4 vs L7

| Type | Sees | Connection handling | Typical use |
|---|---|---|---|
| L4 | IPs, ports, TCP/UDP | Forwards or splices the client's connection; balances **per connection** | raw TCP, TLS passthrough, very high throughput, non-HTTP protocols |
| L7 | HTTP method/path/headers/status | Terminates the client connection, opens its own to backends; balances **per request** | routing by hostname/path, auth, compression, HTTP retries, canaries |

The per-connection vs per-request difference is critical for **long-lived
connections**. HTTP/2 and gRPC clients send *all* their requests over one
connection. Behind an L4 balancer, that client is pinned to one backend
forever: new pods get no traffic and hot pods melt. gRPC needs an L7 balancer,
a service mesh, or client-side balancing
([OSI guide, Part 15](../real-life-example-osi.md#part-15-architecture-decisions-by-layer-the-senior-and-manager-view)).

### How it actually works: choosing a backend

| Algorithm | How it picks | Good for | Weakness |
|---|---|---|---|
| Round robin | Next in the list | Identical backends, uniform requests | Ignores load; a slow backend gets the same share |
| Least connections / least outstanding requests | Fewest in-flight | Variable request cost | Needs shared state across LB instances |
| **Power of two choices** | Pick 2 at random, take the less loaded | Large fleets, many LB instances | Slightly less even than a global view |
| Consistent hashing (ring hash, Maglev) | Hash of a key (client IP, user ID) → backend | Cache affinity, sticky sessions | Hot keys overload one backend |

"Power of two choices" is the quiet favourite of modern proxies (Envoy,
NGINX's `random two`). Picking two at random and taking the better one gives
almost the balance of a global "least loaded" decision without any
coordination between load balancer instances.

### Health checks must test what matters

Bad:

```text
GET / -> 200
```

Better:

```text
GET /readyz -> 200 only if this instance can serve real traffic right now
```

- **Readiness** should fail when the instance should be *removed from
  rotation*: starting up, draining, or unable to reach a dependency it can't
  work without.
- **Liveness** should fail only when the process should be *restarted*:
  deadlocked or corrupted. Never make liveness depend on a database.
  Otherwise a database outage restarts every pod at once, and they all
  stampede the database when it comes back.
- **Hysteresis:** eject after N consecutive failures and restore after M
  successes, so one slow probe doesn't flap a backend in and out.
- **Passive health checks (outlier detection):** also eject backends based on
  *real* traffic, such as 5 consecutive 5xx errors. Active probes run every few
  seconds, while real traffic notices failures in milliseconds.

### Draining

Before shutting down a backend:

1. mark it not-ready (fail `/readyz`);
2. wait for the load balancer to stop sending new traffic (at least
   *failure threshold × check interval*);
3. let in-flight requests finish, up to a deadline;
4. exit.

Skipping draining creates random user-visible errors on every deploy.

### Build it in Go (45 min) — a draining backend behind a health-checking LB

This lab has two roles in one binary: a **backend** that implements the four
draining steps, and a **round-robin L7 load balancer** that health-checks
`/readyz` with hysteresis. You can kill a backend mid-traffic and count the
errors.

```go
// lbdrain: Chapter 57 in one program -- backends with /readyz and graceful
// draining, and a round-robin load balancer with active health checks.
// Kill a backend gracefully mid-traffic and count the errors (there should be 0).
//
//	go build -o lbdrain ./lbdrain
//	./lbdrain backend 127.0.0.1:9081 &
//	./lbdrain backend 127.0.0.1:9082 &
//	./lbdrain lb 127.0.0.1:9080 http://127.0.0.1:9081 http://127.0.0.1:9082 &
//	for i in $(seq 1 2000); do curl -s -o /dev/null -w '%{http_code}\n' localhost:9080/; done | sort | uniq -c &
//	kill -TERM %1      # graceful: readiness fails first, then in-flight requests finish
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: lbdrain backend ADDR | lbdrain lb ADDR BACKEND_URL...")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "backend":
		backend(os.Args[2])
	case "lb":
		loadBalancer(os.Args[2], os.Args[3:])
	}
}

// ---------------------------------------------------------------- backend --

func backend(addr string) {
	var draining atomic.Bool
	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if draining.Load() {
			http.Error(w, "draining", http.StatusServiceUnavailable) // "stop sending me new work"
			return
		}
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond) // pretend to work
		fmt.Fprintf(w, "hello from %s\n", addr)
	})
	srv := &http.Server{Addr: addr, Handler: mux,
		ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Printf("backend %s up (pid %d)", addr, os.Getpid())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	<-ctx.Done()

	// The drain sequence from Chapter 57:
	log.Printf("backend %s: SIGTERM -> failing readiness", addr)
	draining.Store(true)        // 1. mark not-ready
	time.Sleep(3 * time.Second) // 2. wait for the LB's health checks to notice
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil { // 3. stop accepting, finish in-flight
		log.Printf("forced close: %v", err)
	}
	log.Printf("backend %s: drained, exiting", addr) // 4. exit
}

// ----------------------------------------------------------- load balancer --

type upstream struct {
	url     *url.URL
	proxy   *httputil.ReverseProxy
	healthy atomic.Bool
}

func loadBalancer(addr string, targets []string) {
	var ups []*upstream
	for _, t := range targets {
		u, err := url.Parse(t)
		if err != nil {
			log.Fatal(err)
		}
		up := &upstream{url: u, proxy: httputil.NewSingleHostReverseProxy(u)}
		up.healthy.Store(true)
		ups = append(ups, up)
		go healthCheck(up)
	}

	var next atomic.Uint64
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Round-robin over HEALTHY upstreams only.
		for range ups {
			up := ups[next.Add(1)%uint64(len(ups))]
			if up.healthy.Load() {
				up.proxy.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "no healthy upstream", http.StatusServiceUnavailable)
	})
	log.Printf("lb on %s -> %v", addr, targets)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// healthCheck polls /readyz. Two failures to eject, two successes to restore:
// hysteresis stops one slow probe from flapping a backend in and out.
func healthCheck(up *upstream) {
	client := &http.Client{Timeout: time.Second}
	fails, oks := 0, 0
	for range time.Tick(time.Second) {
		resp, err := client.Get(up.url.String() + "/readyz")
		ok := err == nil && resp.StatusCode == http.StatusOK
		if resp != nil {
			resp.Body.Close()
		}
		if ok {
			oks, fails = oks+1, 0
		} else {
			fails, oks = fails+1, 0
		}
		switch {
		case fails == 2 && up.healthy.Load():
			up.healthy.Store(false)
			log.Printf("lb: EJECT %s", up.url)
		case oks == 2 && !up.healthy.Load():
			up.healthy.Store(true)
			log.Printf("lb: RESTORE %s", up.url)
		}
	}
}
```

Build it once, start two backends and the load balancer, send traffic, and
kill one backend **gracefully** while the traffic is flowing. (Build a binary
rather than using `go run`, so the signal goes to the program itself and not
to the `go` command.)

```text
$ go build -o lbdrain ./lbdrain
$ ./lbdrain backend 127.0.0.1:9081 &        # job %1
$ ./lbdrain backend 127.0.0.1:9082 &
$ ./lbdrain lb 127.0.0.1:9080 http://127.0.0.1:9081 http://127.0.0.1:9082 &
$ (for i in $(seq 1 300); do curl -s -o /dev/null -w '%{http_code}\n' localhost:9080/; done | sort | uniq -c) &
$ kill -TERM %1
backend 127.0.0.1:9081: SIGTERM -> failing readiness
lb: EJECT http://127.0.0.1:9081
backend 127.0.0.1:9081: drained, exiting
 300 200                                        <- zero errors
```

Now do the same with `kill -9` (no draining):

```text
 276 200
  24 502                                        <- every in-flight or not-yet-ejected request fails
```

**What to notice:**

- **The 3-second sleep between "not ready" and `Shutdown`** is the most
  important line. Without it, the backend stops accepting connections while
  the LB still thinks it's healthy. Kubernetes has exactly this race: pod
  removal from Service endpoints is asynchronous with SIGTERM, which is why
  production pods use a `preStop` sleep
  ([Linux guide, Chapter 60](../../os-linux/real-life-os-guide.md#chapter-60-deployment-patterns-rolling-updates-health-checks-rollbacks)).
- **`srv.Shutdown(ctx)`** closes the listener, closes *idle* keep-alive
  connections, and waits for active requests to finish, up to the context
  deadline.
- **The 24 errors under `kill -9` were 502s** from `httputil.ReverseProxy`,
  which couldn't reach the dead backend until the health check ejected it
  about 2 seconds later. Passive outlier detection would have cut that
  window to one failed request.

### Who is the client? X-Forwarded-For and the PROXY protocol

Behind any proxy, the backend's TCP peer is the proxy. An **L7** proxy can add
a header (`X-Forwarded-For: 203.0.113.7`). An **L4** proxy can't change HTTP,
because it may not even be HTTP or may be encrypted. Instead it uses the
**PROXY protocol**: one line of text, sent before the client's bytes.

```go
// proxyproto: recover the real client address behind an L4 load balancer.
//
// An L4 proxy (HAProxy, AWS NLB with proxy protocol enabled, Envoy) opens its
// OWN connection to your backend, so RemoteAddr() is the proxy, not the
// client. With the PROXY protocol, the proxy writes one header line first:
//
//	PROXY TCP4 203.0.113.7 10.0.1.15 61001 443\r\n
//
//	go run ./proxyproto            # runs a server and a fake proxy against it
package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

// readProxyV1 parses a PROXY protocol v1 header. Only call it on listeners
// that ONLY your load balancer can reach: otherwise any client can claim to
// be any IP address.
func readProxyV1(r *bufio.Reader) (src, dst netip.AddrPort, err error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return src, dst, err
	}
	if len(line) > 107 || !strings.HasSuffix(line, "\r\n") { // 107 = spec maximum
		return src, dst, errors.New("malformed PROXY header")
	}
	f := strings.Fields(strings.TrimSuffix(line, "\r\n"))
	if len(f) < 2 || f[0] != "PROXY" {
		return src, dst, errors.New("missing PROXY header")
	}
	if f[1] == "UNKNOWN" { // e.g. the LB's own health checks
		return src, dst, nil
	}
	if len(f) != 6 || (f[1] != "TCP4" && f[1] != "TCP6") {
		return src, dst, fmt.Errorf("unsupported PROXY header %q", line)
	}
	sip, err1 := netip.ParseAddr(f[2])
	dip, err2 := netip.ParseAddr(f[3])
	sp, err3 := strconv.ParseUint(f[4], 10, 16)
	dp, err4 := strconv.ParseUint(f[5], 10, 16)
	if err := errors.Join(err1, err2, err3, err4); err != nil {
		return src, dst, fmt.Errorf("bad PROXY header: %w", err)
	}
	return netip.AddrPortFrom(sip, uint16(sp)), netip.AddrPortFrom(dip, uint16(dp)), nil
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 2; i++ {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			r := bufio.NewReader(c) // keep using r: it may have buffered request bytes
			src, _, err := readProxyV1(r)
			if err != nil {
				fmt.Println("rejected:", err)
				c.Close()
				continue
			}
			req, _ := r.ReadString('\n')
			fmt.Printf("TCP peer (the proxy): %s\nreal client:          %s\nrequest:              %q\n",
				c.RemoteAddr(), src, req)
			c.Close()
		}
	}()

	// Pretend to be the load balancer: header first, then the client's bytes.
	c, _ := net.Dial("tcp", ln.Addr().String())
	fmt.Fprint(c, "PROXY TCP4 203.0.113.7 10.0.1.15 61001 443\r\nGET / HTTP/1.1\r\n")
	c.Close()
	// And a client that skipped the proxy and talks to the backend directly.
	c, _ = net.Dial("tcp", ln.Addr().String())
	fmt.Fprint(c, "GET / HTTP/1.1\r\n")
	c.Close()
	<-done
}
```

```text
$ go run ./proxyproto
TCP peer (the proxy): 127.0.0.1:52937
real client:          203.0.113.7:61001
request:              "GET / HTTP/1.1\r\n"
rejected: missing PROXY header
```

**Both mechanisms are trivially forgeable** by anyone who can reach the
backend directly. Accept PROXY headers or `X-Forwarded-For` only from your
load balancer's addresses, and for `X-Forwarded-For` take the address added
by *your* proxy (the right-most trusted entry), never the left-most, which
the client wrote.

### Retry trap

Retries are useful only when:

- the operation is idempotent or has an idempotency key;
- the retry count is bounded, and **only one layer retries** (client, LB,
  and service mesh each retrying 3 times means 27 attempts);
- backoff with jitter exists;
- a **retry budget** caps retries at a fraction of normal traffic (say 10%),
  so retries can't multiply load during an outage;
- the total deadline still matters: a retry that can't finish in time is
  pure waste.

Blind retries turn "one slow backend" into "three times the traffic during
an outage."

### Timeouts must nest

Every hop needs a timeout, and they must fit inside each other:

```
 client 10s  >  CDN 9s  >  LB idle 60s / request 8s  >  app handler 7s  >  DB query 5s
```

And for idle **keep-alive** connections, the rule is reversed: the *backend*
must keep idle connections open **longer** than the load balancer
(e.g. app 65 s > LB 60 s), so the LB always closes first. Otherwise the LB
occasionally sends a request on a connection the backend is closing, and the
user gets a 502 that appears in no application log. The
[OSI guide's Layer 5 war story](../real-life-example-osi.md#part-7-layer-5-session) walks through this exact
incident.

### Common confusions

- **"The load balancer is highly available, so we're fine."** The LB is often
  the most reliable part. Bad health checks and missing draining cause
  most LB-related incidents.
- **"Sticky sessions are a feature."** They're a workaround for stateful
  servers, and they break even distribution and deploys. Move the state out
  ([OSI guide, Part 7](../real-life-example-osi.md#part-7-layer-5-session)).
- **"502 and 504 mean the app crashed."** 502: the proxy got a bad or no
  response (connection refused or reset, crash, keep-alive race). 504: the
  backend didn't answer in time. Check the LB log against the app log: if the
  app never logged the request, the problem is between them.

### Check yourself

1. Why does gRPC behind an L4 load balancer overload some backends?
2. What should a readiness check test, and what must a liveness check *not*
   test?
3. Write the four draining steps, and say what goes wrong if you skip step 2.
4. Why must the backend's keep-alive timeout be longer than the LB's?
5. How should a backend decide whether to trust `X-Forwarded-For`?

### Further reading

- **Article:** "Load Balancing" chapter of Google's *Site Reliability
  Engineering* book (free online), covering subsetting, weighted round robin, and
  the "lame duck" state (draining).
- **Paper:** "Maglev: A Fast and Reliable Software Network Load Balancer"
  (Google, NSDI 2016).
- **Go plan:** build the full versions: [four LB algorithms behind one
  interface (Day 45)](../../Golang/detailed-90-day-plan/week7.md#day-45-load-balancing-algorithms), [active health checking (Day 46)](../../Golang/detailed-90-day-plan/week7.md#day-46-backend-health-checking),
  [connection draining (Day 47)](../../Golang/detailed-90-day-plan/week7.md#day-47-connection-draining), and a
  [complete L4 load balancer (Day 48)](../../Golang/detailed-90-day-plan/week7.md#day-48-tcp-load-balancer-l4).
- **Across the series:** the [HTTPS guide's load balancer chapter](../../v2-https/real-life-guide-v1.md#chapter-9-load-balancers-rate-limiting)
  covers rate limiting and global load balancing.

---

## Chapter 58 — Kubernetes networking: Services, Ingress, NetworkPolicy, CNI

### In one sentence

Kubernetes networking is normal networking with more layers: pod IPs, Service
virtual IPs, kube-proxy or eBPF rules, ingress controllers, cluster DNS, and a
CNI plugin, each with its own failure modes.

### The problem

A request that fails inside a cluster could have died at any of six layers,
most of them invisible to `kubectl get pods`. Teams new to Kubernetes tend to
restart pods until the problem goes away. Teams that understand the packet
path find the broken layer in minutes.

### The model: four rules

Kubernetes requires every network plugin (CNI) to provide:

1. Every **pod gets its own IP address**.
2. Pods can reach all other pods **without NAT**.
3. Nodes can reach all pods (and vice versa) without NAT.
4. A pod sees itself with the same IP that others use to reach it.

How the CNI meets these rules varies: routing pod CIDRs directly (Calico in
BGP mode, cloud VPC CNIs that give pods real VPC addresses) or an overlay
(VXLAN/Geneve, Chapter 54). Overlays explain many MTU problems in clusters.

### The packet path to remember

```text
client
  -> cloud LoadBalancer / Ingress controller / Gateway   (L7: host + path routing)
  -> Service (a virtual ClusterIP, e.g. 10.96.12.34:80)   (exists only as rules!)
  -> kube-proxy rules (iptables / IPVS / nftables) or eBPF (Cilium): DNAT to ...
  -> EndpointSlice: the list of READY pod IPs
  -> Pod IP:containerPort (through the node's veth into the pod namespace)
```

The key insight: **a ClusterIP isn't bound to any interface.** No process
listens on it, and you can't `ping` it in most clusters. It exists only as
NAT rules on every node that rewrite the destination to a pod IP chosen from the
EndpointSlice. That's why a Service with zero ready endpoints answers with
"connection refused" or a timeout, depending on the implementation.

Debug in this order, or reverse it from inside the pod.

### How it actually works: the details that cause incidents

**1. Cluster DNS and `ndots:5`.** A pod's `/etc/resolv.conf` usually looks
like:

```text
nameserver 10.96.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

`ndots:5` means "a name with fewer than 5 dots is tried with every search
suffix *first*". So `api.github.com` (2 dots) becomes, in order:

```text
api.github.com.default.svc.cluster.local   -> NXDOMAIN
api.github.com.svc.cluster.local           -> NXDOMAIN
api.github.com.cluster.local               -> NXDOMAIN
api.github.com                             -> answer
```

That's 4 lookups (×2 for A and AAAA) for every external name. On a busy
service this is a large share of CoreDNS load and adds latency. Fixes: use a
trailing dot (`api.github.com.`), lower `ndots` in the pod's `dnsConfig`, or
run NodeLocal DNSCache.

**2. The conntrack DNS race.** glibc sends the A and AAAA queries in parallel
from the same UDP socket. Both packets race through conntrack's NAT on the
node, one can lose the insert, and it is silently dropped. The client then
waits its full **5-second** DNS timeout. If a cluster shows requests taking
"exactly 5 s more" now and then, suspect this. Mitigations: NodeLocal
DNSCache, `options single-request-reopen`, or a newer kernel/CNI that fixes
the race.

**3. `externalTrafficPolicy`.** For `LoadBalancer`/`NodePort` Services:

| | `Cluster` (default) | `Local` |
|---|---|---|
| Which nodes accept traffic | all | only nodes running a ready pod |
| Extra hop | possibly (node → another node's pod) | no |
| Client source IP | **lost** (SNAT to the node's IP) | **preserved** |
| Balance | even across pods | uneven if pods are unevenly spread |

If your app logs show node IPs instead of client IPs, this is why.

**4. Headless Services** (`clusterIP: None`) skip the virtual IP. DNS returns
the pod IPs directly. StatefulSets use them for stable per-pod names, and gRPC
clients use them for client-side load balancing (Chapter 57's long-lived
connection problem).

**5. NetworkPolicy.** If no policy selects a pod, all traffic is allowed. Once
any policy selects a pod, only traffic that some policy allows gets through.
And the CNI must *enforce* NetworkPolicy, which not every CNI does. A
default-deny starting point for a namespace:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: shop
spec:
  podSelector: {}                 # every pod in the namespace
  policyTypes: [Ingress, Egress]
  egress:                         # ...but still allow DNS, or nothing will resolve
  - to:
    - namespaceSelector: {}
      podSelector:
        matchLabels: {k8s-app: kube-dns}
    ports:
    - {protocol: UDP, port: 53}
    - {protocol: TCP, port: 53}
```

Forgetting the DNS exception is the most common NetworkPolicy outage:
everything "times out", because nothing can resolve a name.

### Commands that answer real questions

```bash
kubectl get pod -o wide                          # pod IPs and nodes
kubectl get svc,endpointslices -l app=web        # does the Service have READY endpoints?
kubectl describe ingress web                     # rules, backends, events
kubectl get networkpolicy -A                     # anything selecting these pods?
kubectl logs -n kube-system -l k8s-app=kube-dns  # CoreDNS errors
```

A throwaway debugging pod with every network tool in it:

```bash
kubectl run -it --rm netshoot --image=nicolaka/netshoot -- bash
# inside:
cat /etc/resolv.conf                              # search domains and ndots
dig web.default.svc.cluster.local +short          # Service ClusterIP
dig +search web                                   # same, via search path
curl -sv http://web:80/healthz                    # through the Service
curl -sv http://<pod-ip>:8080/healthz             # bypass the Service: is it the pod or the Service?
```

On a node (or a `kind` node container): `iptables-save | grep KUBE-SVC | head`
shows the NAT rules that *are* the Services.

### Practice (45 min): see a Service is just NAT rules

With `kind` ([Linux guide, Chapter 53](../../os-linux/real-life-os-guide.md#chapter-53-why-kubernetes-exists-the-problem-docker-alone-doesn-t-solve) has the
setup):

```bash
kind create cluster --name net
kubectl create deployment web --image=nginx --replicas=3
kubectl expose deployment web --port=80
kubectl get svc web                                      # note the CLUSTER-IP
docker exec net-control-plane iptables-save | grep -E 'KUBE-SVC|KUBE-SEP' | grep -m6 web
kubectl scale deployment web --replicas=1
docker exec net-control-plane iptables-save | grep -c KUBE-SEP   # fewer endpoint rules
kind delete cluster --name net
```

You'll see one `KUBE-SVC-…` chain for the Service and one `KUBE-SEP-…`
("service endpoint") chain per pod, chosen with `--probability` rules: random
load balancing implemented in the firewall.

### Common fingerprints

| Symptom | Likely cause |
|---|---|
| Service has no endpoints | label selector doesn't match pods, or pods not Ready |
| Pod can call by IP but not by name | CoreDNS down, wrong namespace in the name, NetworkPolicy blocking port 53 |
| Random +5 s latency on outbound calls | conntrack DNS race; `ndots:5` amplification |
| Ingress returns 404 | host/path rule mismatch, wrong ingress class |
| Connection times out only across namespaces | NetworkPolicy |
| Client IPs in logs are all node IPs | `externalTrafficPolicy: Cluster` SNAT |
| Large responses fail in cluster | overlay MTU / PMTU issue (Chapter 54) |
| New pods get no gRPC traffic | long-lived HTTP/2 connections through a ClusterIP (L4) Service |

### Build it in Go (20 min) — client-side balancing over a headless Service

Use what you built in Chapters 18 and 57: resolve a **headless** Service with
`net.DefaultResolver.LookupHost(ctx, "web-headless.default.svc.cluster.local")`
to get every ready pod IP, then round-robin requests across them yourself and
re-resolve every 30 s. This is what gRPC's `dns:///` resolver with the
`round_robin` policy does, and it fixes the "new pods get no traffic" row
above. Run it in-cluster from a `golang` pod, scale the Deployment, and watch
the IP list change.

### Common confusions

- **"I can ping the Service IP."** Usually you can't. There's no interface
  behind it to answer ICMP. Test with the real port (`curl`, `nc -vz`).
- **"Kubernetes load-balances my gRPC calls."** A ClusterIP Service balances
  *connections*. One long-lived HTTP/2 connection carries all calls to one pod.
- **"NetworkPolicy is enforced everywhere."** Only if the CNI supports it.
  Test a deny rule before relying on it.

### Check yourself

1. Why can't a ClusterIP be pinged, and what actually makes it work?
2. Explain `ndots:5` and the cost of resolving `api.github.com` from a pod.
3. When would you choose `externalTrafficPolicy: Local`, and what does it cost?
4. Write a NetworkPolicy that denies all egress except DNS.

### Further reading

- **Docs:** kubernetes.io, "Services, Load Balancing, and Networking" and
  "Debug Services".
- **Article:** "Racy conntrack and DNS lookup timeouts" (Weave Works blog), the
  original write-up of the 5-second DNS timeout.
- **Across the series:** the [Linux guide's Kubernetes chapters](../../os-linux/real-life-os-guide.md#chapter-56-services-and-how-traffic-finds-a-pod)
  cover Services and Pods from the workload side, and the
  [Go guide's Kubernetes chapter](../../Golang/real-life-golang-guide.md#67-running-go-in-containers-and-kubernetes) covers what a Go service
  must do to behave well there (probes, graceful shutdown, resource limits).

---

## Chapter 59 — Network performance engineering: QoS, packet loss, shaping, saturation

### In one sentence

Performance work means separating latency, loss, bandwidth, queueing, and
application delay, then fixing whichever one is actually the limit, instead
of guessing.

### The problem

"The network is slow" can mean a dozen different things. A 1 Gb/s link that
gives one transfer 40 Mb/s might be lossy, high-latency, window-limited,
shaped, congested, or perfectly healthy with a slow application at one end.
Each needs a different fix, and buying more bandwidth fixes only one of them.

### The four questions

1. Is there packet loss?
2. Is latency high even without load?
3. Does latency increase only under load? (Queueing / bufferbloat.)
4. Is throughput limited by bandwidth, window size, CPU, disk, or the app?

### How it actually works: three laws to carry around

**1. Bandwidth-delay product (Chapter 24).** One TCP connection can have at
most one window of unacknowledged data in flight, so:

```
 max throughput of one connection  ≈  window / RTT
 window needed to fill a link      =  bandwidth × RTT
 1 Gb/s × 100 ms = 12.5 MB in flight        (a 64 KB window gives ~5 Mb/s)
```

**2. Little's law.** For any stable queue: *items in the system = arrival
rate × time each spends in it*. A server handling 2,000 requests/s at 50 ms
each has ~100 requests in flight, so it needs ~100 concurrent workers or
connections. If latency doubles, so does concurrency, which is why a slow
dependency exhausts connection pools.

**3. Mathis' formula for loss-limited TCP (Reno/CUBIC-style).**
`throughput ≲ (MSS / RTT) × (1.22 / √loss)`. With MSS 1460, RTT 50 ms, and
1% loss, that's about 2.9 Mb/s, however fast the link. Loss hurts most
on long paths. BBR (Chapter 25) is much less sensitive to random loss, which
is why it helps on lossy long-haul and mobile links.

### Useful tests

```bash
mtr -ezbw -c 100 example.com                 # per-hop loss and latency (look at the LAST hop)
iperf3 -c server.example.com -t 20           # single-stream throughput
iperf3 -c server.example.com -P 8            # 8 streams: if much faster, you're window-limited
iperf3 -c server.example.com -R              # reverse direction (server -> client)
ss -tin dst server.example.com               # rtt, cwnd, retrans, and "rwnd_limited"/"sndbuf_limited"
curl -w '@curl-format.txt' -o /dev/null -s https://example.com/
```

`curl-format.txt`:

```text
dns=%{time_namelookup}
connect=%{time_connect}
tls=%{time_appconnect}
ttfb=%{time_starttransfer}
total=%{time_total}
```

Interpreting `iperf3`: if **one stream is slow but eight together fill the
link**, the limit is per-connection: a window or buffer that's too small, or
loss. If eight streams are also slow, the limit is the path or an endpoint's
CPU or NIC.

### Build it in Go (30 min) — measure the bandwidth-delay product

This lab sends data as fast as possible to a sink and reports throughput. It
can pin the socket buffer (`SetWriteBuffer`, which disables Linux's
autotuning) and open parallel streams, which lets you reproduce the first law
on demand:

```go
// bdp: feel the bandwidth-delay product (Chapters 24 and 59).
//
// One TCP connection can have at most one window of data "in flight". If the
// window (bounded by socket buffers) is smaller than bandwidth x RTT, the pipe
// is never full and throughput = window / RTT, no matter how fast the link is.
//
//	go run ./bdp -serve :7500
//	go run ./bdp -to HOST:7500 -buf 65536            # small buffer
//	go run ./bdp -to HOST:7500                       # kernel autotuning
//	go run ./bdp -to HOST:7500 -streams 4 -buf 65536 # parallel streams
//
// Add delay first to make the effect obvious (Linux, on the client):
//
//	sudo tc qdisc add dev eth0 root netem delay 50ms
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	serve := flag.String("serve", "", "run the sink server on this address")
	to := flag.String("to", "", "server to send to")
	buf := flag.Int("buf", 0, "socket send buffer in bytes (0 = kernel autotuning)")
	streams := flag.Int("streams", 1, "parallel connections")
	secs := flag.Int("secs", 5, "test duration in seconds")
	flag.Parse()

	if *serve != "" {
		ln, err := net.Listen("tcp", *serve)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("sink listening on", ln.Addr())
		for {
			c, err := ln.Accept()
			if err != nil {
				continue
			}
			go func() { io.Copy(io.Discard, c); c.Close() }()
		}
	}
	if *to == "" {
		flag.Usage()
		os.Exit(2)
	}

	rtt := measureRTT(*to)
	var total atomic.Int64
	deadline := time.Now().Add(time.Duration(*secs) * time.Second)
	var wg sync.WaitGroup
	for i := 0; i < *streams; i++ {
		c, err := net.Dial("tcp", *to)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if *buf > 0 {
			c.(*net.TCPConn).SetWriteBuffer(*buf) // pins the size: disables autotuning
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer c.Close()
			chunk := make([]byte, 128*1024)
			for time.Now().Before(deadline) {
				n, err := c.Write(chunk)
				total.Add(int64(n))
				if err != nil {
					return
				}
			}
		}()
	}
	wg.Wait()

	mbps := float64(total.Load()) * 8 / float64(*secs) / 1e6
	fmt.Printf("streams=%d buf=%d rtt~%v  ->  %.1f Mbit/s\n", *streams, *buf, rtt.Round(time.Millisecond), mbps)
	if *buf > 0 && rtt > 0 {
		// Linux doubles the SO_SNDBUF you ask for (half is bookkeeping overhead),
		// so roughly 2*buf bytes can be in flight per connection.
		ceiling := float64(2**buf) * 8 / rtt.Seconds() / 1e6 * float64(*streams)
		fmt.Printf("theory: streams x 2*buf / RTT ~ %.1f Mbit/s\n", ceiling)
	}
}

// measureRTT takes the fastest of three TCP handshakes: one handshake ~ 1 RTT.
func measureRTT(addr string) time.Duration {
	best := time.Hour
	for i := 0; i < 3; i++ {
		start := time.Now()
		c, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		best = min(best, time.Since(start))
		c.Close()
	}
	return best
}
```

Real results between two containers, with 50 ms of delay added by `netem`:

```text
$ go run ./bdp -to sink:7500 -buf 65536
streams=1 buf=65536 rtt~55ms  ->  18.4 Mbit/s
theory: streams x 2*buf / RTT ~ 19.0 Mbit/s

$ go run ./bdp -to sink:7500                    # kernel autotuning
streams=1 buf=0 rtt~52ms  ->  230.7 Mbit/s

$ go run ./bdp -to sink:7500 -buf 65536 -streams 4
streams=4 buf=65536 rtt~51ms  ->  72.9 Mbit/s
theory: streams x 2*buf / RTT ~ 82.0 Mbit/s
```

**What to notice:**

- **The theory line predicts the measurement** to within a few percent. The
  network was the same in all three runs. Only the window changed.
- **Linux doubles the `SO_SNDBUF` you ask for** (half is bookkeeping), so
  `2 × buf` is in flight. This surprises everyone who reads `man 7 socket`
  for the first time.
- **Setting a buffer size *disables* autotuning.** Autotuning reached 230 Mbit/s
  without help. Hard-coding "big" buffers in an application often makes
  things *worse* than leaving the kernel alone. Tune the kernel's ceilings
  (`net.ipv4.tcp_wmem`/`tcp_rmem` max) for long fat networks, not each socket.
- **Parallel streams multiply the window.** That's why download managers, `aws
  s3 cp`, and `rclone` open many connections on high-RTT paths.

**Exercises:**

1. Add `-delay` values of 10, 50, and 150 ms (with `tc netem`) and plot
   throughput for the pinned 64 KB buffer. It should fall as 1/RTT.
2. Add 1% loss and compare CUBIC with BBR by switching
   `sudo sysctl -w net.ipv4.tcp_congestion_control=bbr` on the sender
   between runs (load it first with `sudo modprobe tcp_bbr` if needed). Which
   degrades less? Compare both with Mathis' prediction above.
3. Run it across the internet between two cloud regions. How close does a
   single stream get to the instance's advertised bandwidth?

### Queues, bufferbloat, and shaping

Latency that rises *only under load* is a queue filling somewhere: a home
router, a cable modem, a VPN gateway, a cloud NAT. Large, dumb FIFO buffers
cause **bufferbloat**: hundreds of milliseconds of queueing delay that
make video calls stutter while a download runs.

- **Active queue management** (`fq_codel`, `cake`) keeps queues short by
  dropping or marking early and isolating flows from each other. Most Linux
  distributions now default to `fq_codel`; check with `tc qdisc show`.
- **Shaping** means sending slightly *below* the bottleneck rate, so the queue
  forms on a device you control (with good AQM) instead of in the ISP's
  modem. `tc qdisc add dev eth0 root cake bandwidth 90mbit`.

**QoS does not create bandwidth.** It decides which packets suffer first when
a link is full. It's useful for voice, video, and control traffic, but it
can't fix an undersized link or an application that sends too much.

**Bufferbloat rule:** if ping latency explodes during an upload or download,
you need queue management or shaping near the bottleneck, not a faster DNS
server.

### Host tuning that actually matters

```bash
sysctl net.ipv4.tcp_congestion_control net.core.default_qdisc   # e.g. cubic, fq_codel
sysctl net.ipv4.tcp_rmem net.ipv4.tcp_wmem                       # min default max (bytes)
sysctl net.core.somaxconn net.ipv4.tcp_max_syn_backlog           # accept/SYN queues (Chapter 52)
sysctl net.ipv4.ip_local_port_range                              # ephemeral ports (Chapter 19)
```

Change these only with a measurement showing they're the limit: before,
after, and in a canary.

### Common confusions

- **"More bandwidth will fix it."** Not if the limit is RTT × window, loss,
  queueing, or the application. Measure first.
- **"Average latency is fine."** Bufferbloat and loss show up in p99, not the
  mean (Chapter 45).
- **"iperf3 shows 900 Mb/s, so the app should get that."** iperf3 measures the
  network with an ideal application. Your app adds serialisation, TLS,
  syscalls per small write, disk, and locks.

### Check yourself

1. Compute the window needed for 500 Mb/s at 80 ms RTT.
2. Using Little's law, how many DB connections does a service doing 300
   queries/s at 40 ms each need at minimum?
3. What do you conclude if 1 iperf3 stream gives 50 Mb/s and 8 streams give 400 Mb/s?
4. Why can setting `SO_SNDBUF` in an app reduce throughput?

### Further reading

- **Book:** *High Performance Browser Networking* (Ilya Grigorik, free online),
  chapters 1–2 on latency, bandwidth, and TCP.
- **Article:** "Bufferbloat: Dark Buffers in the Internet" (Gettys & Nichols,
  ACM Queue) and bufferbloat.net.
- **Across the series:** the [Linux guide's observability chapter](../../os-linux/real-life-os-guide.md#chapter-70-observability-and-performance-perf-flame-graphs-ebpf-and-opentelemetry)
  covers host-side profiling for when the bottleneck is the application, not the
  network, and the [Go guide's performance chapter](../../Golang/real-life-golang-guide.md#37-performance-tuning-pprof-trace-and-benchmark-methodology) covers
  profiling Go programs.

---

## Chapter 60 — Incident response: runbooks, packet evidence, and escalation

### In one sentence

During an outage, your job is to narrow the blast radius, preserve evidence,
communicate clearly, and change one thing at a time. After it, your job is to
make the same failure easier to detect and harder to repeat.

### The problem

Network incidents are unusually prone to chaos. Symptoms appear far from the
cause ("the app is down" when a route changed), several teams own pieces of
the path, evidence disappears when someone restarts something, and everyone
has a theory. A simple, practised structure beats individual brilliance.

### Roles, even for small incidents

| Role | Job | Not their job |
|---|---|---|
| **Incident commander (IC)** | Owns the incident: decides priorities, assigns work, decides when to mitigate | Debugging. The IC coordinates and doesn't type commands |
| **Operations lead(s)** | Investigate and change systems, one change at a time, announced first | Talking to customers |
| **Communications lead** | Status updates to stakeholders and customers on a fixed cadence (e.g. every 30 min) | Speculating on root cause |
| **Scribe** (often the IC early on) | A timestamped log of observations, hypotheses, and actions | Filtering what seems unimportant |

On a two-person team, one person is IC and comms and the other operates. The
point is that *someone* is deciding and *someone* is writing it down.

### First 10 minutes

```text
1. What changed?          deploys, config, firewall, DNS, certificates, provider status pages
2. Who is affected?       one user, one AZ, one region, one ISP, everyone? since when?
3. Which layer?           DNS, routing, firewall, TLS, load balancer, app, dependency?
4. Is there a safe rollback?   rolling back a recent change beats debugging it live
5. What evidence do we have before changing state?
```

Scope is the fastest discriminator. "Only users on one ISP" points at routing
or peering. "Only one AZ" points at infrastructure. "Everyone, since 09:05" points
at whatever changed at 09:05.

### Evidence bundle

Collect before and after any change, with UTC timestamps:

```bash
date -u
ip addr; ip route; ip rule
ss -s; ss -tin | head -50
nstat -az | grep -Ei 'retrans|overflow|drop|reset'
dig +short service.example.com; dig +short service.example.com @1.1.1.1
curl -v --max-time 10 -w '\ndns=%{time_namelookup} tcp=%{time_connect} tls=%{time_appconnect} ttfb=%{time_starttransfer}\n' \
     -o /dev/null https://service.example.com/healthz
sudo tcpdump -ni any -s 0 -w incident-$(date -u +%H%M).pcap 'host 203.0.113.10'
```

For intermittent problems, leave a **ring-buffer capture** running so the
evidence exists when the problem next happens:

```bash
sudo tcpdump -ni eth0 -s 128 -C 100 -W 20 -w /var/tmp/ring.pcap 'tcp port 443'
# 20 files × 100 MB, oldest overwritten; -s 128 keeps headers, not payloads
```

Captures can contain credentials and personal data, so treat them as
sensitive and delete them when the investigation ends.

### Hypothesis discipline

Write every hypothesis in the incident log in this form:

```text
09:24  H3: backend subnet NACL drops return traffic since rule change CR-1234
       test: tcpdump on LB shows SYN-ACK from backend? on backend shows SYN arriving?
       result 09:29: SYN arrives at backend, SYN-ACK leaves backend, never reaches LB -> SUPPORTED
```

One hypothesis, one test, one result. It stops the team from changing three
things at once and losing track of what actually fixed it.

### Escalation template

```text
Impact: 30% of users in us-east cannot reach api.example.com
Started: 2026-09-19 09:12 UTC
Evidence: SYN reaches load balancer; no SYN-ACK from backend subnet
Not DNS: authoritative and public resolvers return expected IP
Not client-side: reproduced from two networks
Recent changes: firewall rule update at 09:05 UTC
Rollback: revert rule change CR-1234
Need: network/firewall owner to confirm drop counters
```

Good escalation is specific. "Network is broken" is not actionable. "SYN
arrives, SYN-ACK leaves, never reaches the LB, since 09:05" is.

### After the incident: the blameless postmortem

```text
Summary          one paragraph, written for someone outside the team
Impact           who, how many, how long, which SLOs were burned
Timeline (UTC)   detection, escalation, mitigation, resolution, with evidence links
Root cause       the mechanism, not a person ("a NACL rule lacked the ephemeral return range")
Contributing     why it wasn't caught: no policy tests, no alert on SYN-ACK drops, ...
What went well   keep doing these
Action items     each with an owner and a date: prevent, detect sooner, mitigate faster
```

"Blameless" doesn't mean nobody is accountable. It means the system that let
a reasonable person make the mistake is what gets fixed. For network
incidents, the best action items are usually **tests** (Chapter 53's
`policytest`), **alerts** (Chapter 47), and **runbooks** that would have cut the
time to diagnosis.

### Common confusions

- **"Root-cause first, then fix."** Mitigate first (roll back, fail over,
  shed load), then find the root cause, using evidence you captured before
  mitigating.
- **"The status page can wait."** Customers notice silence. Regular updates,
  even "still investigating", reduce support load and build trust.
- **"Restarting fixed it, so we're done."** A restart is a mitigation. Without
  evidence, the incident will happen again
  ([Linux guide, Chapter 72](../../os-linux/real-life-os-guide.md#chapter-72-incident-response-and-refactoring-on-real-linux-systems)).

### Check yourself

1. Name the four incident roles and the one thing each should *not* do.
2. Why capture evidence before rolling back?
3. Rewrite "the database is unreachable" as a specific escalation.
4. What makes an action item good?

### Further reading

- **Book:** Google's *Site Reliability Engineering*, chapters "Managing
  Incidents" and "Postmortem Culture" (free online).
- **Guide:** PagerDuty Incident Response documentation (response.pagerduty.com),
  an openly published, practical incident process.
- **Across the series:** the [OSI guide's debugging playbook](../real-life-example-osi.md#part-14-debugging-by-layer-the-incident-playbook),
  the [HTTPS guide's capstone diagnosis](../../v2-https/real-life-guide-v1.md#chapter-18-capstone-diagnose-a-slow-broken-or-unsafe-https-request-end-to-end), and the
  [Linux guide's incident chapter](../../os-linux/real-life-os-guide.md#chapter-72-incident-response-and-refactoring-on-real-linux-systems) each give the
  layer-by-layer investigation steps that this process coordinates.

---

## Chapter 61 — Network design/refactoring: migrations, diagrams, IPAM, change safety

### In one sentence

Network refactoring means changing paths, prefixes, policies, or dependencies
without losing the ability to explain the system or roll it back.

### The problem

Networks outlive the decisions that shaped them. Two companies merge and both
used `10.0.0.0/16`. A "temporary" VPN becomes the production path. A
firewall rule nobody understands blocks a migration. Unlike application code,
network state is spread across devices, cloud consoles, and DNS, and there's
often no test suite. Refactoring it safely needs written artefacts and a
staged process.

### The artefacts professionals keep current

- **Diagram:** L3 boundaries, routing domains, firewalls, load balancers,
  VPNs, NAT, internet edges. Ideally kept as code (Mermaid, D2, diagrams.net
  files in git) so it's reviewed alongside changes.
- **IPAM:** who owns each prefix, gateway, VLAN/VPC/subnet, DHCP pool, and
  reserved range. A spreadsheet is a start; NetBox or a cloud IPAM is better.
- **Policy matrix:** source → destination → protocol/port → reason → owner.
  Chapter 53's `policytest` turns it into tests.
- **Runbooks:** how to test, fail over, and roll back.
- **Change log:** when and why routes, firewall rules, DNS, and TLS changed.

### IP address planning: decisions that last a decade

- **Allocate from one plan.** Give each environment and region a large,
  non-overlapping block (for example `10.0.0.0/12` for prod, `10.16.0.0/12` for
  staging), and carve VPCs/VNets from it. Every future peering, VPN,
  and acquisition depends on this.
- **Leave room.** Subnets are hard to resize. Kubernetes in particular
  consumes addresses fast (pod CIDRs, service CIDRs).
- **Avoid the popular defaults** (`10.0.0.0/16`, `192.168.0.0/24`,
  `172.17.0.0/16` for Docker) for anything that will ever be connected to
  something else. Someone else will have used them.
- **Plan IPv6 too.** A `/56` per site or VPC and `/64` per subnet is
  conventional. IPv6 makes address *scarcity* go away but not address *planning*.

### Build it in Go (30 min) — grow the subnet tool into an IPAM check

The `subnet` program from Chapter 10 already has the two operations every IPAM
is built on: **overlap detection** and **first-fit allocation**:

```text
$ go run ./subnet overlap 10.0.0.0/16 10.0.128.0/20 172.16.0.0/12
OVERLAP  10.0.0.0/16  <->  10.0.128.0/20

$ go run ./subnet alloc 10.0.0.0/16 22 10.0.0.0/24 10.0.4.0/22
allocated 10.0.8.0/22
```

**Exercises (a small, real IPAM):**

1. Read allocations from a CSV kept in git
   (`prefix,owner,environment,purpose`) and make `overlap` fail CI when a pull
   request adds a conflicting prefix. That single check prevents the most
   expensive network mistake there is.
2. Add `alloc -owner team-payments -purpose "eks pods"`, which appends the new
   row to the CSV. Allocation now leaves a record of who asked and why.
3. Add a `check-cloud` command that compares the CSV with reality (for
   example, the output of `aws ec2 describe-subnets` saved as JSON) and reports
   prefixes in the cloud but not in IPAM, and the reverse. That's **drift
   detection**.

### Migration pattern

```text
1. Observe current traffic and dependencies (flow logs, DNS query logs, captures).
2. Add the new path in parallel.
3. Test with one client or canary prefix.
4. Move a small percentage (weighted DNS, LB weights, BGP prepending).
5. Monitor (error rate, latency, and the OLD path's traffic, which should fall).
6. Move the rest.
7. Keep rollback until the old path is demonstrably unused.
8. Remove the old path later, deliberately, as its own change.
```

Step 7 needs evidence: "nothing has used the old path for 7 days" from flow
logs or counters, not "we think everything moved".

### Change safety

- **Peer review** every network change, as you would code, including firewall
  and DNS changes in the console. Infrastructure-as-code (Terraform,
  Pulumi) makes review possible.
- **Pre-flight and post-flight tests.** Run the policy tests, `netcheck`, and
  a real transaction before and after. If the "before" run fails, stop: you
  can't verify a change against a baseline that's already broken.
- **Smallest blast radius first:** one AZ, one region, one canary prefix.
- **Automatic rollback** for anything that could cut off your own access
  (Chapter 53's safety rule).
- **Change freezes** around critical business events, with an exception
  process for real emergencies.

### Common migration failures

| Failure | Prevention |
|---|---|
| Hidden dependency on old IP | flow logs, DNS logs, packet captures; grep configs for literal IPs |
| Firewall allows old subnet only | policy matrix and pre-flight tests from the new subnet |
| TTL too high for fast rollback | lower TTL ahead of migration (Chapter 55) |
| Asymmetric routing | trace both directions, check stateful firewalls |
| MTU changes in new path | test large packets and real application payloads (Chapter 54) |
| Overlapping CIDRs discovered mid-project | central IPAM with overlap checks in CI |

### Common confusions

- **"It's just a subnet change."** Every prefix appears in firewall rules,
  route tables, security groups, allow-lists at partners, monitoring, and
  sometimes hard-coded in applications. Search for it before changing it.
- **"We'll document it after the migration."** The diagram and IPAM *are* the
  plan. Without them you can't review the change.

### Check yourself

1. Why is IP planning a leadership decision and not only a network one?
2. What evidence lets you remove an old path safely?
3. Design a CI check that would have prevented two VPCs from overlapping.
4. List the pre-flight tests you'd run before moving a subnet behind a new firewall.

### Further reading

- **Book:** *Network Programmability and Automation* (O'Reilly), on IPAM,
  source of truth, and automation.
- **Tool:** NetBox (netbox.dev), the most widely used open-source IPAM/DCIM
  source of truth.
- **Across the series:** the [OSI guide's IP planning war story](../real-life-example-osi.md#part-5-layer-3-network)
  shows the cost of getting this wrong, and the
  [Linux guide's fleet and drift chapter](../../os-linux/real-life-os-guide.md#chapter-66-production-linux-mental-models-fleet-drift-and-failure-domains) applies the same
  source-of-truth idea to servers.

---

## Chapter 62 — Capstone lab: build, break, observe, and repair a production-like network

### Goal

Build a small lab that contains the pieces you will meet in real work:
clients, DNS, a load balancer, two app servers, a database subnet, firewall
rules, TLS, monitoring, and packet capture.

### Minimum topology

```text
client subnet
   |
firewall/router
   |
load balancer
  / \
app1 app2
   |
database subnet
```

You can implement this with VMs, Linux network namespaces, Docker/Podman,
or cloud free-tier resources. The value is not the tool; the value is
watching traffic cross each boundary.

### Required exercises

1. Draw the diagram before building.
2. Assign subnets and document them.
3. Make DNS names for the load balancer and database.
4. Terminate TLS at the load balancer.
5. Add health checks and take one app down.
6. Capture a successful request at client, load balancer, and app.
7. Break DNS and diagnose it.
8. Break the firewall and distinguish drop vs reject.
9. Break MTU with a tunnel or lowered interface MTU.
10. Fill the uplink and observe latency under load.
11. Rotate a certificate.
12. Write a one-page incident report for one failure.

### Passing standard

You are done when you can answer, from evidence:

- Where did the packet stop?
- Which layer failed?
- Which command proved it?
- What fixed it?
- What would have alerted earlier?
- What rollback would have restored service fastest?

This capstone turns the guide from reading into operational muscle memory.

### The Go variant: build the lab from your own tools

If you've done the "Build it in Go" sections, you can build most of this
capstone from programs you wrote, which means you understand every moving
part:

| Lab component | Your Go program | Chapter |
|---|---|---|
| Load balancer with health checks and draining | `lbdrain lb` + `lbdrain backend` ×2 | 57 |
| TLS termination with hot certificate rotation | wrap `lbdrain lb` in `certreload`'s `GetCertificate` | 56 |
| Client IP preservation through an L4 hop | `proxyproto` in front of the backends | 57 |
| Firewall verification from each subnet | `policytest` with one policy file per subnet | 53 |
| Synthetic monitoring and latency histograms | `prober` scraped by Prometheus | 46–47 |
| "Which layer is broken?" in one command | `netcheck` | 34 |
| DNS consistency across resolvers | `dnsdiff` | 55 |
| Packet evidence | `tcpdump -w`, then `pcapread` | 31–32 |

Put each component in its own network namespace (Chapter 52) or container,
connect them with a bridge, apply the nftables policy from Chapter 53 to the
"firewall" namespace, and run the twelve required exercises. The Go plan's
capstone, a [cloud-edge security proxy (Week 17)](../../Golang/detailed-90-day-plan/week17.md), is the
production-grade version of the same idea, with a WAF, rate limiting, and
DDoS defences in front.


---

# Part 14 — Where to go next

**You've reached the end of the core material.** Parts 1–12 took you from
"what problem does a network solve" to reading packet captures, diagnosing
real failures, and understanding BGP and OSPF — the full arc this guide set
out to teach. This last part is a checkpoint: what you can genuinely do now,
what's honestly still missing, and where to point that six-month plan you've
been meaning to make.

## Chapter 63 — Honest gaps and a six-month plan

### What you can now do

- Explain the whole path from URL to pixels
- Subnet confidently
- Read packet captures
- Diagnose the five classic failure modes systematically
- Understand why things are slow, not just that they are

That is more than many working developers and support engineers know.

### Honest gaps this guide left

| Gap | Where to fill it |
|---|---|
| Routing protocol **configuration** (Part 12 covers the concepts, not the CLI) | CCNA material; *Routing TCP/IP* (Doyle); FRRouting docs |
| Switching in depth (STP, VLANs, trunking, port security) | CCNA; *Network Warrior* |
| Email protocols (SMTP flow, STARTTLS, SPF/DKIM/DMARC) | *Postfix: The Definitive Guide*; the DMARC.org guides |
| Physical layer (cabling, fibre types, SFPs, PoE) | CompTIA Network+; vendor cabling guides |
| Multicast (IGMP, PIM) and MPLS/VRF | CCNP material; *Routing TCP/IP, Vol 2* |
| IPsec internals (IKE, SAs, tunnel vs transport) | *IPsec VPN Design*; your firewall vendor's docs |
| Wireless in depth (RF, channels, roaming, 802.1X) | *802.11 Wireless Networks* (Gast); CWNA |
| Network security in depth | Security+; *Practical Packet Analysis*; OWASP |
| Cloud networking in depth | Your provider's docs and certification |
| Automation (Ansible, Python for networks, NetBox) | *Network Programmability and Automation* (O'Reilly) |
| Theory and maths (queueing, algorithms) | Kurose & Ross; Peterson & Davie |
| Voice/video, QoS, SIP | Specialist material once you need it |

### A six-month plan

```
MONTH 1-2 : CEMENT THE FUNDAMENTALS
  * Re-read Parts 3, 5, and 6 of this guide -- they'll read differently now
  * Work through Kurose & Ross chapters 1-5 (or the free lectures)
  * Do 50 subnetting problems until it's boring
  Deliverable: a network diagram of your home/office with every address explained

MONTH 3 : PACKETS
  * Read "Practical Packet Analysis" (Chris Sanders) cover to cover
  * Analyse one capture per day: your own traffic, then the Wireshark samples
  Deliverable: a personal Wireshark cheat sheet and 5 annotated captures

MONTH 4 : BUILD A LAB
  * Two or three VMs (or containers, or a Raspberry Pi)
  * Configure: static addresses, routing between subnets, a firewall, DHCP, DNS
  * Break each piece deliberately and fix it
  Deliverable: a working multi-subnet lab you built from nothing

MONTH 5 : GO WIDE OR DEEP -- choose ONE
    (a) Certification -- CompTIA Network+ then CCNA (structured, employable)
    (b) Systems/SRE -- Linux networking, containers, Kubernetes networking
    (c) Security -- Security+, packet forensics, TryHackMe
    (d) Performance -- HPBN, tuning, BBR, load testing
  Deliverable: a certification, or a written project

MONTH 6 : TEACH IT
  * Write up something you found hard (subnetting? MTU? TIME_WAIT?)
  * Answer questions on r/networking or Server Fault
  Deliverable: a published article
```

**The best habit:** whenever something on a network confuses you, *capture it*.
Twenty minutes with `tcpdump` teaches more than an hour of searching.

---

## Chapter 64 — Project ideas by level

Each project lists what it builds on. The Go projects reuse the labs from
this guide. The Go plan references point to step-by-step versions.

**Beginner**
1. Map your home network: every device, its IP, MAC, and what it does
   (Chapters 5, 6, 9).
2. Extend `netcheck` (Chapter 34) to report link, address, gateway, DNS, and
   internet reachability, with JSON output, and run it from cron.
3. Capture and fully annotate one web page load, packet by packet
   (Chapters 4, 31–32). Then decode the same capture with `pcapread`.
4. Set up a Raspberry Pi as a DNS server (Pi-hole) and watch your own queries
   (Chapter 18). Use `dnsdiff` to compare it with public resolvers.
5. Measure and fix bufferbloat on your home connection (Chapters 27, 59).

**Intermediate**
6. Build a two-subnet lab with a Linux router between them, and route between
   them by hand (Chapters 13, 52).
7. Set up a WireGuard VPN and investigate its MTU behaviour (Chapter 54).
8. Write a TCP client and server **in Go** with a length-prefixed binary
   protocol, and watch it in Wireshark (Chapter 22's `stream` lab; Go plan
   [Day 53: binary KV protocol](../../Golang/detailed-90-day-plan/week8.md#day-53-binary-protocol-design)).
9. Put a load balancer (nginx, HAProxy, or your own `lbdrain`) in front of two
   web servers with health checks. Kill one and watch the failover in a
   capture (Chapter 57).
10. Compare HTTP/1.1, HTTP/2, and HTTP/3 performance on the same content under
    simulated loss (`tc netem`) (Chapters 30, 59).
11. Write a **UDP-based reliable transfer** in Go: sequence numbers, ACKs,
    retransmission timer, sliding window. Then compare its throughput with TCP's
    under 2% loss (Chapters 20, 22, 24).

**Advanced**
12. Build a container network from scratch with `ip netns` and `veth`, with no
    Docker (Chapter 52).
13. Set up BGP between two routers (FRRouting or BIRD) in a lab (Chapter 50).
14. Extend `pcapread` into a mini-Wireshark: TCP option decoding, conversations,
    retransmission and RTT analysis, TLS SNI extraction (Chapter 32).
15. Instrument a Go service with per-connection `TCP_INFO` metrics (rtt,
    retransmits, cwnd) exported to Prometheus, and build a dashboard
    (Chapters 33, 46).
16. Run a full IPv6-only lab with NAT64/DNS64 and see what breaks (Chapter 16).
17. Build an **L4 load balancer in Go** with PROXY protocol, health checks, and
    connection draining (Chapter 57; Go plan [Day 48](../../Golang/detailed-90-day-plan/week7.md#day-48-tcp-load-balancer-l4)).
18. Build an **SNI router**: accept TLS connections, read the ClientHello
    without decrypting, and forward to a backend chosen by hostname
    (Chapters 29, 32; Go plan [Day 55](../../Golang/detailed-90-day-plan/week8.md#day-55-sni-routing)).
19. Build a **DDoS-aware edge**: SYN-flood and Slowloris detection, per-CIDR
    rate limiting (Go plan [Week 13](../../Golang/detailed-90-day-plan/week13.md)), tested with traffic you
    generate in your own lab.

**Capstone (multi-week):** the Go plan's [cloud-edge security proxy](../../Golang/detailed-90-day-plan/week17.md)
combines TLS termination, L7 routing, a WAF, rate limiting, observability,
and graceful operations. It covers almost everything in Parts 5–13 of this
guide in one program.

---

## Chapter 65 — The series map: where every topic lives

This guide is step 2b of the series (OS → networking → security → HTTPS). It's
strongest on L2–L4 and on seeing protocols on the wire. Use this table to jump to
the guide that covers each topic from a different angle.

| Topic | This guide | OSI guide | Linux guide | Security guides | HTTPS guide | Go |
|---|---|---|---|---|---|---|
| The whole path of one request | [Ch 4](#chapter-4-our-running-example-what-happens-when-you-open-a-website) | [Part 10](../real-life-example-osi.md#part-10-the-full-replay-a-millisecond-timeline) | — | — | [Ch 0.3](../../v2-https/real-life-guide-v1.md#0-3-the-complete-flow-at-a-glance), [Ch 25](../../v2-https/real-life-guide-v1.md#chapter-25-one-https-request-every-layer-every-guide) | `walkthrough` lab |
| Layers and encapsulation | [Ch 2–3](#chapter-3-the-four-layers-of-tcp-ip) | [Part 2](../real-life-example-osi.md#part-2-the-mental-model), [Part 11](../real-life-example-osi.md#part-11-x-ray-of-one-packet) | — | — | [Ch 6](../../v2-https/real-life-guide-v1.md#chapter-6-osi-layers-the-network-stack) | — |
| Ethernet, ARP, switching | [Ch 5–7](#chapter-5-mac-addresses-and-ethernet-frames) | [Part 4](../real-life-example-osi.md#part-4-layer-2-data-link) | [Ch 37](../../os-linux/real-life-os-guide.md#chapter-37-how-a-linux-box-sees-the-network) | — | — | — |
| IP addressing, subnets, NAT | [Ch 9–12](#chapter-9-what-an-ip-address-actually-is) | [Part 5](../real-life-example-osi.md#part-5-layer-3-network) | — | [In Depth Ch 8](../../security/real-life-security-guide-v1.md#chapter-8-cloud-network-security) | [Ch 1](../../v2-https/real-life-guide-v1.md#chapter-1-how-the-internet-works), [Ch 7](../../v2-https/real-life-guide-v1.md#chapter-7-nat-vip-bgp-anycast-vrrp-global-traffic-routing) | `subnet`, `whoami` |
| Routing, BGP, anycast, VRRP | [Ch 13](#chapter-13-routing-how-a-packet-finds-its-way), [Ch 48–51](#chapter-48-two-ways-to-share-directions) | [Part 5](../real-life-example-osi.md#part-5-layer-3-network) | [Ch 69](../../os-linux/real-life-os-guide.md#chapter-69-advanced-linux-networking-namespaces-routing-nftables-and-packet-paths) | — | [Ch 7](../../v2-https/real-life-guide-v1.md#chapter-7-nat-vip-bgp-anycast-vrrp-global-traffic-routing) | — |
| DNS | [Ch 18](#chapter-18-dns-turning-names-into-addresses), [Ch 55](#chapter-55-dns-operations-authoritative-dns-delegation-split-horizon-outages) | [Part 9](../real-life-example-osi.md#part-9-layer-7-application) | [Ch 37](../../os-linux/real-life-os-guide.md#chapter-37-how-a-linux-box-sees-the-network) | — | [Ch 2](../../v2-https/real-life-guide-v1.md#chapter-2-dns-finding-the-address) | `dnsquery`, `dnsdiff` |
| Sockets, threads, descriptors | [Ch 19](#chapter-19-ports-and-sockets-which-program-gets-the-data) | [Part 6](../real-life-example-osi.md#part-6-layer-4-transport) | [Ch 73](../../os-linux/real-life-os-guide.md#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime), [Ch 78](../../os-linux/real-life-os-guide.md#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections) | — | — | [Go §34](../../Golang/real-life-golang-guide.md#34-networking-deep-dive-net-conn-tcp-udp-framing-your-own-protocol) |
| TCP in depth | [Ch 21–27](#chapter-21-tcp-part-1-the-three-way-handshake) | [Part 6](../real-life-example-osi.md#part-6-layer-4-transport) | — | — | [Ch 3](../../v2-https/real-life-guide-v1.md#chapter-3-tcp-reliable-delivery) | [Go plan Week 5](../../Golang/detailed-90-day-plan/week5.md) |
| TLS and certificates | [Ch 29](#chapter-29-tls-how-it-gets-encrypted), [Ch 56](#chapter-56-tls-and-certificate-operations-expiry-chains-sni-rotation) | [Part 8](../real-life-example-osi.md#part-8-layer-6-presentation) | — | [From Zero](../../security/real-life-guide.md), [In Depth Ch 23](../../security/real-life-security-guide-v1.md#chapter-23-service-identity-spiffe-and-spire) | [Ch 4](../../v2-https/real-life-guide-v1.md#chapter-4-tls-https-the-encrypted-tunnel), [Ch 16](../../v2-https/real-life-guide-v1.md#chapter-16-certificate-and-domain-operations-acme-rotation-ct-monitoring-and-emergency-revocation) | [Go §35](../../Golang/real-life-golang-guide.md#35-tls-in-go-crypto-tls-mutual-tls), [Week 6](../../Golang/detailed-90-day-plan/week6.md) |
| HTTP/1.1, /2, /3 | [Ch 28](#chapter-28-http-how-the-web-actually-talks), [Ch 30](#chapter-30-http-2-and-http-3), [Ch 39–44](#chapter-39-how-http-messages-are-framed-and-how-it-goes-wrong) | [Part 9](../real-life-example-osi.md#part-9-layer-7-application) | — | [In Depth Ch 30](../../security/real-life-security-guide-v1.md#chapter-30-http-request-smuggling-and-desync) | [Ch 5](../../v2-https/real-life-guide-v1.md#chapter-5-http-the-conversation), [Ch 14](../../v2-https/real-life-guide-v1.md#chapter-14-http-version-negotiation-http-1-1-http-2-http-3-alt-svc-and-fallbacks) | [Go §24](../../Golang/real-life-golang-guide.md#24-net-http-fundamentals-client-server-middleware), [Week 9](../../Golang/detailed-90-day-plan/week9.md) |
| Proxies, load balancers, CDNs | [Ch 36](#chapter-36-load-balancers-cdns-and-anycast), [Ch 57](#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries) | [Part 15](../real-life-example-osi.md#part-15-architecture-decisions-by-layer-the-senior-and-manager-view) | — | [In Depth Ch 31](../../security/real-life-security-guide-v1.md#chapter-31-web-cache-poisoning-and-cache-deception) | [Ch 8–10](../../v2-https/real-life-guide-v1.md#chapter-8-proxies-the-middlemen), [Ch 20](../../v2-https/real-life-guide-v1.md#chapter-20-reverse-proxies-in-go-the-client-s-identity-framing-and-request-smuggling) | [Week 7](../../Golang/detailed-90-day-plan/week7.md) |
| Host networking, namespaces, conntrack | [Ch 52](#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack) | [Part 12](../real-life-example-osi.md#part-12-inside-the-server-from-the-nic-to-your-handler) | [Ch 44](../../os-linux/real-life-os-guide.md#chapter-44-what-a-container-actually-is-namespaces-cgroups-a-filesystem), [Ch 69](../../os-linux/real-life-os-guide.md#chapter-69-advanced-linux-networking-namespaces-routing-nftables-and-packet-paths), [Ch 80](../../os-linux/real-life-os-guide.md#chapter-80-a-container-runtime-in-150-lines-of-go) | [In Depth Ch 15](../../security/real-life-security-guide-v1.md#chapter-15-container-internals-and-isolation) | — | — |
| Containers and Kubernetes networking | [Ch 37](#chapter-37-cloud-and-container-networking), [Ch 58](#chapter-58-kubernetes-networking-services-ingress-networkpolicy-cni) | — | [Ch 50](../../os-linux/real-life-os-guide.md#chapter-50-container-networking-how-containers-talk-to-the-world), [Ch 56](../../os-linux/real-life-os-guide.md#chapter-56-services-and-how-traffic-finds-a-pod) | [In Depth Ch 19](../../security/real-life-security-guide-v1.md#chapter-19-kubernetes-network-security) | — | [Go §67](../../Golang/real-life-golang-guide.md#67-running-go-in-containers-and-kubernetes) |
| Firewalls and network security | [Ch 38](#chapter-38-security-basics), [Ch 53](#chapter-53-firewalls-in-the-real-world-host-firewalls-cloud-security-groups-nacls) | — | [Ch 71](../../os-linux/real-life-os-guide.md#chapter-71-linux-security-hardening-selinux-apparmor-auditd-capabilities-and-least-privilege) | [In Depth Ch 8](../../security/real-life-security-guide-v1.md#chapter-8-cloud-network-security), [Ch 33](../../security/real-life-security-guide-v1.md#chapter-33-ssrf-mastery) | [Ch 12](../../v2-https/real-life-guide-v1.md#chapter-12-security-protecting-the-entire-lifecycle), [Ch 22](../../v2-https/real-life-guide-v1.md#chapter-22-outbound-requests-ssrf-and-safe-http-clients) | [Weeks 10–14](../../Golang/detailed-90-day-plan/week10.md) |
| Performance and capacity | [Ch 24–27](#chapter-24-flow-control-the-sliding-window), [Ch 59](#chapter-59-network-performance-engineering-qos-packet-loss-shaping-saturation) | [Part 10](../real-life-example-osi.md#part-10-the-full-replay-a-millisecond-timeline) | [Ch 41](../../os-linux/real-life-os-guide.md#chapter-41-reading-load-average-vmstat-and-iostat-correctly), [Ch 75–76](../../os-linux/real-life-os-guide.md#chapter-75-cpu-limits-gomaxprocs-cgroups-and-throttling) | [In Depth Ch 28](../../security/real-life-security-guide-v1.md#chapter-28-resilience-as-a-security-property) | [Ch 13](../../v2-https/real-life-guide-v1.md#chapter-13-traffic-direction-sli-slo), [Ch 21](../../v2-https/real-life-guide-v1.md#chapter-21-resilience-between-services-rate-limits-retries-idempotency-and-circuit-breakers) | [Go §37](../../Golang/real-life-golang-guide.md#37-performance-tuning-pprof-trace-and-benchmark-methodology) |
| Monitoring, detection, alerting | [Ch 45–47](#chapter-45-what-to-measure-and-why-averages-lie) | — | [Ch 34](../../os-linux/real-life-os-guide.md#chapter-34-logs-journalctl-and-where-output-really-goes), [Ch 70](../../os-linux/real-life-os-guide.md#chapter-70-observability-and-performance-perf-flame-graphs-ebpf-and-opentelemetry) | [In Depth Ch 51–52](../../security/real-life-security-guide-v1.md#chapter-51-the-detection-engineering-lifecycle) | [Ch 23](../../v2-https/real-life-guide-v1.md#chapter-23-slis-slos-and-error-budgets-for-https-services) | [Go §65](../../Golang/real-life-golang-guide.md#65-observability-in-go-services-logs-metrics-traces-health-profiling) |
| Debugging and incidents | [Ch 34–35](#chapter-34-a-troubleshooting-method-that-works), [Ch 60](#chapter-60-incident-response-runbooks-packet-evidence-and-escalation) | [Part 14](../real-life-example-osi.md#part-14-debugging-by-layer-the-incident-playbook) | [Ch 42](../../os-linux/real-life-os-guide.md#chapter-42-strace-lsof-and-answering-why-is-this-stuck), [Ch 72](../../os-linux/real-life-os-guide.md#chapter-72-incident-response-and-refactoring-on-real-linux-systems) | [In Depth Ch 55](../../security/real-life-security-guide-v1.md#chapter-55-incident-response-for-cloud-and-containers) | [Ch 17–18](../../v2-https/real-life-guide-v1.md#chapter-17-debugging-the-lifecycle-browser-devtools-curl-openssl-packet-capture-and-logs) | [Go §58](../../Golang/real-life-golang-guide.md#58-debugging-go-in-production-delve-pprof-traces-stack-dumps-godebug) |

## Continue the series

**Next, step 3a (Security: foundations): [Security from Zero](../../security/real-life-guide.md).** It covers encoding vs hashing vs encryption, TLS, PKI, SSH, the OWASP Top 10, threat modelling. Hands-on: OpenSSL labs and a vulnerable app to break.

---

# Appendix A — Glossary (plain language)

**ACK** — Acknowledgement. "I received everything up to byte N."

**Anycast** — The same IP address announced from many locations; routing sends you
to the nearest.

**ARP** — Address Resolution Protocol. Finds the MAC address for an IP address on
your local network, by broadcasting a question.

**Bandwidth** — How much data per second a link can carry. *Not* the same as
speed (see latency).

**BGP** — Border Gateway Protocol. How independent networks tell each other which
addresses they can reach. Runs the internet's routing.

**Broadcast** — A message to everyone on the local network. `255.255.255.255` or
the all-ones host address.

**Broadcast domain** — The set of devices that receive each other's broadcasts. A
router (or VLAN) breaks it up; a switch does not.

**Bufferbloat** — Excessive latency caused by oversized buffers filling up at the
bottleneck. Test it: ping while downloading.

**CDN** — Content Delivery Network. Copies of content in many locations so users
connect to something nearby.

**CIDR** — The `/24` notation. Counts how many bits of the address are the network
part.

**Congestion control** — TCP's mechanism for finding a safe sending rate by
speeding up until loss occurs. Protects the *network*.

**Default gateway** — The router your machine sends everything it doesn't
recognise to.

**DHCP** — The protocol that automatically gives your device an address, gateway,
and DNS servers. Four messages: DISCOVER, OFFER, REQUEST, ACK.

**DNS** — Turns names (`example.com`) into addresses. Hierarchical and heavily
cached.

**Encapsulation** — Each layer wrapping the layer above's data in its own header.

**Ephemeral port** — A temporary, high-numbered port your machine picks when it
starts an outbound connection.

**Ethernet** — The dominant wired link-layer technology. Frames, MAC addresses,
1500-byte MTU.

**FIN** — The TCP flag meaning "I've finished sending."

**Firewall** — Decides which packets to allow. *Stateful* ones track connections
and auto-allow return traffic.

**Flow control** — TCP's sliding window; stops a fast sender overwhelming a slow
*receiver*.

**Frame** — A chunk of data at the link layer (has MAC addresses).

**Handshake** — The three-packet exchange (SYN, SYN-ACK, ACK) that opens a TCP
connection.

**HTTP** — The protocol of the web. Text-based request/response.

**HTTPS** — HTTP inside TLS encryption.

**ICMP** — The internet's error and diagnostics messenger. `ping` and `traceroute`
are built on it.

**IP address** — A number identifying a machine's location on the network. IPv4 is
32 bits (`192.168.1.50`), IPv6 is 128 bits.

**Jitter** — Variation in latency. Shown as `mdev`/`stddev` in ping output.

**Latency** — Time for data to travel there and back. Dominated by distance.

**Load balancer** — Spreads requests across multiple servers. L4 uses IP/port; L7
reads the HTTP request.

**Loopback** — `127.0.0.1` / `::1`. Traffic to yourself that never leaves the
machine.

**MAC address** — A 48-bit hardware address on a network card. Only meaningful on
the local network.

**MSS** — Maximum Segment Size. The biggest chunk of *data* in a TCP packet.
Usually 1460.

**MTU** — Maximum Transmission Unit. The biggest packet a link will carry.
Usually 1500.

**Multicast** — A message to a group of interested devices (not everyone).

**NAT** — Network Address Translation. Rewrites private addresses to a shared
public one. How your whole house shares one internet address.

**Packet** — A chunk of data at the internet layer (has IP addresses).

**Ping** — Sends an ICMP echo request and measures the round trip.

**Port** — A 16-bit number identifying which program on a machine a packet is
for.

**Private address** — `10.x`, `172.16-31.x`, `192.168.x`. Usable internally, never
routed on the internet.

**Protocol** — An agreed format and set of rules for a conversation.

**QUIC** — A modern transport protocol built on UDP; the foundation of HTTP/3.

**RST** — TCP's abrupt "reset" — connection aborted immediately.

**Round-trip time (RTT)** — How long a packet takes to get there and back.

**Router** — Forwards packets between different networks, using IP addresses.

**Routing table** — The list of "to reach this range, send it that way" rules.

**Segment** — A chunk of data at the TCP layer.

**Socket** — One end of a connection: an IP address plus a port.

**Subnet mask** — Marks where the network part of an address ends. `/24` or
`255.255.255.0`.

**Switch** — Forwards frames between devices on the same network, using MAC
addresses.

**SYN** — The TCP flag that starts a connection.

**TCP** — Reliable, ordered, connection-based transport. Retransmits lost data.

**TIME_WAIT** — A state a closed connection sits in for ~60 s to catch late
packets.

**TLS** — Encrypts and authenticates a connection. The "S" in HTTPS.

**Traceroute** — Discovers the routers along a path by sending packets with
increasing TTL.

**TTL** — Time To Live. A hop counter; decremented by each router. At zero the
packet is dropped.

**UDP** — Fast, connectionless, unreliable transport. Ports and a checksum, nothing
more.

**VLAN** — Virtual LAN. Splits one physical switch into separate broadcast
domains.

**VPN** — An encrypted tunnel; your packets wrapped inside other packets.

**Window** — How many bytes the receiver will accept before waiting for
acknowledgement.

---

# Appendix B — Command cheat sheet

```
### WHAT AM I? ####################################################
ip -br addr / ip -br link          your addresses and interfaces (Linux)
ifconfig                            same (macOS/BSD)
ip route / netstat -rn              your routing table
cat /etc/resolv.conf                your DNS servers
curl ifconfig.me                    your PUBLIC address
ip neigh / arp -a                   your ARP cache (local neighbours)

### CAN I REACH IT? ###############################################
ping -c 5 <host>                    basic reachability + latency + loss
ping -M do -s 1472 <host>           MTU test (Linux; -D on macOS)
traceroute -n <host>                the path, hop by hop
traceroute -T -p 443 <host>         TCP traceroute (works when ICMP is filtered)
mtr -rwzbc 100 <host>               ** best tool: loss + latency per hop **
nc -zv <host> <port>                is this port open? refused? timeout?

### DNS ###########################################################
dig <name>                          full answer
dig <name> +short                   just the address
dig +trace <name>                   ** watch the whole delegation chain **
dig @1.1.1.1 <name>                 ask a specific server (bypass cache)
dig -x <ip>                         reverse lookup
dig <name> A / AAAA / MX / NS / TXT  specific record types

### WHAT'S ON MY MACHINE? #########################################
ss -tuln                            what's LISTENING
ss -tulnp                           ... with the process (sudo for others')
ss -tan                             all TCP sockets and their states
ss -tin                             ** detailed TCP health per connection **
ss -tan | awk '{print $1}' | sort | uniq -c | sort -rn    state histogram
lsof -i :8080                       what's using this port (macOS/Linux)
nstat -az | grep -i retrans         kernel TCP counters (Linux)
netstat -s -p tcp                   TCP stats (macOS)

### HTTP ##########################################################
curl -v https://host                verbose: DNS, TCP, TLS, headers
curl -sI https://host               headers only
curl -s -o /dev/null -w '%{http_code}\n' https://host
curl -w 'dns %{time_namelookup} tcp %{time_connect} tls %{time_appconnect} ttfb %{time_starttransfer} total %{time_total}\n' -o /dev/null -s https://host
nc host 80                          type HTTP by hand

### TLS ###########################################################
openssl s_client -connect host:443 -servername host </dev/null
echo | openssl s_client -connect host:443 -servername host 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates

### PACKETS #######################################################
sudo tcpdump -i any -n -c 20                          quick look
sudo tcpdump -i any -n 'host X and port Y'            filtered
sudo tcpdump -i any -n -e arp                         watch ARP
sudo tcpdump -i any -n -A 'tcp port 80'               readable HTTP
sudo tcpdump -i any -n 'tcp[tcpflags] & tcp-syn != 0'  connection attempts
sudo tcpdump -i any -n 'tcp[tcpflags] & tcp-rst != 0'  resets
sudo tcpdump -i any -n -s0 -w file.pcap 'host X'      save for Wireshark
sudo tcpdump -r file.pcap -n                          read it back

### SIMULATE PROBLEMS (Linux; ALWAYS UNDO) ########################
sudo tc qdisc add dev eth0 root netem delay 200ms
sudo tc qdisc add dev eth0 root netem loss 5%
sudo tc qdisc add dev eth0 root netem delay 150ms 50ms loss 2%
sudo tc qdisc del dev eth0 root                       <-- UNDO
```

---

# Appendix C — Subnetting cheat sheet and answers

## The table to memorise

| CIDR | Mask | Block size | Usable hosts |
|---|---|---|---|
| /24 | 255.255.255.0 | 256 | 254 |
| /25 | 255.255.255.128 | 128 | 126 |
| /26 | 255.255.255.192 | 64 | 62 |
| /27 | 255.255.255.224 | 32 | 30 |
| /28 | 255.255.255.240 | 16 | 14 |
| /29 | 255.255.255.248 | 8 | 6 |
| /30 | 255.255.255.252 | 4 | 2 |
| /31 | 255.255.255.254 | 2 | 2 (point-to-point) |
| /32 | 255.255.255.255 | 1 | 1 (single host) |

Larger networks:

| CIDR | Mask | Contains | Usable hosts |
|---|---|---|---|
| /23 | 255.255.254.0 | 2 × /24 | 510 |
| /22 | 255.255.252.0 | 4 × /24 | 1,022 |
| /21 | 255.255.248.0 | 8 × /24 | 2,046 |
| /20 | 255.255.240.0 | 16 × /24 | 4,094 |
| /16 | 255.255.0.0 | 256 × /24 | 65,534 |
| /8 | 255.0.0.0 | 65,536 × /24 | 16,777,214 |

**Method:** block size = 256 − mask value. Count in blocks from 0 until you pass
your address. The block you're in is the network; the next block minus 1 is the
broadcast.

## Chapter 10, Set A — answers

| # | Given | Network | Broadcast | Usable range | Hosts |
|---|---|---|---|---|---|
| 1 | 192.168.10.75 /24 | 192.168.10.0 | 192.168.10.255 | .1 – .254 | 254 |
| 2 | 192.168.10.75 /26 | 192.168.10.64 | 192.168.10.127 | .65 – .126 | 62 |
| 3 | 10.5.3.200 /25 | 10.5.3.128 | 10.5.3.255 | .129 – .254 | 126 |
| 4 | 172.20.15.100 /28 | 172.20.15.96 | 172.20.15.111 | .97 – .110 | 14 |
| 5 | 10.0.0.6 /30 | 10.0.0.4 | 10.0.0.7 | .5 – .6 | 2 |
| 6 | 192.168.4.170 /27 | 192.168.4.160 | 192.168.4.191 | .161 – .190 | 30 |
| 7 | 172.16.35.42 /23 | 172.16.34.0 | 172.16.35.255 | 172.16.34.1 – 172.16.35.254 | 510 |
| 8 | 10.10.10.10 /22 | 10.10.8.0 | 10.10.11.255 | 10.10.8.1 – 10.10.11.254 | 1022 |

*Working for #7:* /23 → mask 255.255.**254**.0 → block size 256−254 = **2**, in the
third octet. Count in 2s: 0, 2, 4 … 34, 36. `35` falls in the **34** block.

*Working for #8:* /22 → mask 255.255.**252**.0 → block size **4** in the third
octet. Count: 0, 4, **8**, 12. `10` falls in the **8** block.

## Chapter 10, Set B — answers

| # | Question | Answer | Why |
|---|---|---|---|
| 9 | 192.168.1.10/24 & 192.168.1.200/24 | **Same** | Both in 192.168.1.0–255 |
| 10 | 192.168.1.10/26 & 192.168.1.200/26 | **Different** | Block 64: `.10` is in the 0-block, `.200` is in the 192-block |
| 11 | 10.0.0.100/25 & 10.0.0.130/25 | **Different** | Block 128: `.100` in 0–127, `.130` in 128–255 |
| 12 | 172.16.4.5/22 & 172.16.7.250/22 | **Same** | Block 4 in octet 3: both in the 4-block (4–7) |

## Chapter 11 — answers

**Problem 1** — `10.20.0.0/24`, largest first:

| Order | Department | Need | Prefix | Network | Usable | Broadcast |
|---|---|---|---|---|---|---|
| 1 | Wi-Fi guests | 100 | /25 | 10.20.0.0 | .1 – .126 | 10.20.0.127 |
| 2 | Sales | 60 | /26 | 10.20.0.128 | .129 – .190 | 10.20.0.191 |
| 3 | Engineering | 25 | /27 | 10.20.0.192 | .193 – .222 | 10.20.0.223 |
| 4 | Management | 10 | /28 | 10.20.0.224 | .225 – .238 | 10.20.0.239 |
| 5 | Link 1 | 2 | /30 | 10.20.0.240 | .241 – .242 | 10.20.0.243 |
| 6 | Link 2 | 2 | /30 | 10.20.0.244 | .245 – .246 | 10.20.0.247 |

Used: 248 of 256. Left over: `10.20.0.248` – `.255`.

**Problem 2** — `172.16.0.0/22` (that's 172.16.0.0 – 172.16.3.255):

| Order | Area | Need | Prefix | Network | Usable |
|---|---|---|---|---|---|
| 1 | Floor 1 | 200 | /24 | 172.16.0.0 | 172.16.0.1 – 172.16.0.254 |
| 2 | Floor 2 | 200 | /24 | 172.16.1.0 | 172.16.1.1 – 172.16.1.254 |
| 3 | Floor 3 | 100 | /25 | 172.16.2.0 | 172.16.2.1 – 172.16.2.126 |
| 4 | Servers | 50 | /26 | 172.16.2.128 | 172.16.2.129 – 172.16.2.190 |
| 5 | Management | 20 | /27 | 172.16.2.192 | 172.16.2.193 – 172.16.2.222 |

Used: 736 of 1024. Left over: `172.16.2.224` – `172.16.3.255`.

**Problem 3** — the four /24s `192.168.4.0` through `192.168.7.0` combine into
**`192.168.4.0/22`**.

*Why:* in binary, the third octet values 4, 5, 6, 7 are `00000100`, `00000101`,
`00000110`, `00000111`. The first **six** bits are identical (`000001`), so
16 + 6 = **22** bits are common → `/22`.

---

# Appendix D — Ports and protocols reference

## Ports worth recognising

| Port | Protocol | Service |
|---|---|---|
| 20 / 21 | TCP | FTP (data / control) |
| **22** | TCP | **SSH / SFTP** |
| 23 | TCP | Telnet (unencrypted — avoid) |
| **25** | TCP | SMTP (mail between servers) |
| **53** | **UDP + TCP** | **DNS** (TCP for large answers — must be open!) |
| 67 / 68 | UDP | **DHCP** (server / client) |
| 69 | UDP | TFTP |
| **80** | TCP | **HTTP** |
| 110 / 143 | TCP | POP3 / IMAP (mail retrieval) |
| **123** | UDP | NTP (time sync) |
| 161 / 162 | UDP | SNMP (monitoring) |
| 179 | TCP | BGP |
| 389 / 636 | TCP | LDAP / LDAPS |
| **443** | **TCP + UDP** | **HTTPS** (UDP = HTTP/3 / QUIC) |
| 445 | TCP | SMB (Windows file sharing) |
| 465 / 587 | TCP | SMTP submission (sending mail) |
| 500 / 4500 | UDP | IPsec VPN |
| 853 | TCP/UDP | DNS over TLS / QUIC |
| 993 / 995 | TCP | IMAPS / POP3S |
| 1194 | UDP | OpenVPN |
| 1433 / 1521 | TCP | MS SQL Server / Oracle |
| 3306 | TCP | MySQL / MariaDB |
| 3389 | TCP | RDP (Windows remote desktop) |
| 5432 | TCP | PostgreSQL |
| 5353 | UDP | mDNS (`.local` discovery) |
| 6379 | TCP | Redis |
| 8080 / 8443 | TCP | HTTP / HTTPS alternates |
| 27017 | TCP | MongoDB |
| 51820 | UDP | WireGuard VPN |

## Protocol numbers (in the IP header)

| Number | Protocol |
|---|---|
| 1 | ICMP |
| 6 | **TCP** |
| 17 | **UDP** |
| 47 | GRE (tunnelling) |
| 50 | ESP (IPsec) |
| 58 | ICMPv6 |
| 89 | OSPF |

## ICMP types worth knowing

| Type | Meaning |
|---|---|
| 0 / 8 | Echo Reply / Request (**ping**) |
| 3 | Destination unreachable (code 1 = host, 3 = port, **4 = fragmentation needed**) |
| 11 | **Time exceeded** (powers **traceroute**) |

## Address ranges

| Range | Meaning |
|---|---|
| `10.0.0.0/8` | Private (RFC 1918) |
| `172.16.0.0/12` | Private (`172.16` – `172.31`) |
| `192.168.0.0/16` | Private |
| `100.64.0.0/10` | Carrier-grade NAT |
| `127.0.0.0/8` | Loopback |
| `169.254.0.0/16` | **Self-assigned — DHCP failed** |
| `224.0.0.0/4` | Multicast |
| `192.0.2.0/24`, `198.51.100.0/24`, `203.0.113.0/24` | Reserved for documentation |

---

# Appendix E — Reading list and fact-checking notes

## If you only do five things

1. **Watch:** Practical Networking's "Networking Fundamentals" series (free,
   YouTube) — alongside Parts 1–5 of this guide.
2. **Watch + drill:** Practical Networking's "Subnetting Mastery" series, then do
   30 problems at `subnettingpractice.com`.
3. **Read:** *High Performance Browser Networking* (Ilya Grigorik) — **free at
   hpbn.co**. Parts 6 and 7 of this guide, done properly.
4. **Read:** *Practical Packet Analysis* (Chris Sanders) — for Part 8.
5. **Do:** the projects in Chapter 64, starting with mapping your own network.

## Books

| Book | For |
|---|---|
| *Computer Networking: A Top-Down Approach* — Kurose & Ross | **The** standard university textbook. Readable, comprehensive. If you buy one, buy this. |
| *High Performance Browser Networking* — Ilya Grigorik (**free**: hpbn.co) | TCP performance, TLS, HTTP/2. Excellent and modern. |
| *Practical Packet Analysis* — Chris Sanders | Wireshark, for beginners, with real captures |
| *TCP/IP Illustrated, Vol. 1* — Fall & Stevens | The deep reference. Dense; dip into it. |
| *Network Warrior* — Gary Donahue | Practical, real-world, opinionated. Great for switching/routing. |
| *CompTIA Network+ Study Guide* | Structured coverage of everything, exam-focused |
| *Wireshark Network Analysis* — Laura Chappell | The packet-analysis reference |
| *802.11 Wireless Networks* — Matthew Gast | If you go deep on wireless |

## Free video courses

- **Practical Networking (Ed Harmoush)** — youtube.com/practicalnetworking.
  "Networking Fundamentals" and "Subnetting Mastery". **The best free video
  material for this guide**; the pacing and diagrams are excellent.
- **Professor Messer** — free full CompTIA Network+ course. Structured, exam-aligned,
  completely free.
- **NetworkChuck** — energetic, beginner-friendly, good for motivation and a first
  pass on subnetting/DNS/OSI.
- **PowerCert Animated Videos** — short animated explainers, good for a first look
  at any single topic.
- **Hussein Nasser** — networking *for developers*: HTTP, TCP, TLS, databases.
  Excellent if you're a programmer.
- **Chris Greer** — packet analysis walkthroughs; pairs with Part 8.

## Blogs and sites

- **jvns.ca (Julia Evans)** — networking explained honestly and kindly. Her free
  zines and `messwithdns.net` are outstanding for beginners.
- **Cloudflare Learning Center** (`cloudflare.com/learning`) — short, accurate,
  well-illustrated articles on almost every topic in this guide. Excellent first
  stop.
- **Cloudflare Blog** — deep technical posts and incident post-mortems.
- **APNIC Blog** — routing, IPv6, DNS, measurement.
- **Wireshark Wiki** — sample captures and protocol references.
- **bufferbloat.net** — for Chapter 27.

## Hands-on practice

- `subnettingpractice.com` — endless generated subnetting problems
- `messwithdns.net` — a real DNS playground (Julia Evans)
- `badssl.com` — every kind of TLS failure, safely
- `test-ipv6.com` — check your IPv6
- `waveform.com/tools/bufferbloat` — bufferbloat grading
- **Cisco Packet Tracer** (free with registration) or **GNS3** — build virtual
  networks
- **TryHackMe / HackTheBox** — legal hands-on security labs

## Specifications (for when you want the truth)

- **RFC 9293** — TCP. Published August 2022; replaced RFC 793 (1981).
- **RFC 8200** — IPv6.
- **RFC 1918** — private address ranges. Two pages.
- **RFC 826** — ARP. Three pages, and a good first RFC.
- **RFC 2131** — DHCP.
- **RFC 1034 / 1035** — DNS.
- **RFC 9110–9114** — HTTP semantics, HTTP/1.1, HTTP/2, HTTP/3.
- **RFC 9000** — QUIC.

Read RFCs at `rfc-editor.org` or `datatracker.ietf.org`.

## How the facts in this guide were checked

Historical claims, protocol details, and specific numbers were verified against
primary sources. Specifically confirmed:

- **RFC 9293** (Transmission Control Protocol) was published **August 2022** and
  obsoletes **RFC 793** (1981) — a 41-year gap.
- **Linux TIME_WAIT is hardcoded at 60 seconds** (`TCP_TIMEWAIT_LEN`, defined as
  2 × MSL) and is **not** controlled by `net.ipv4.tcp_fin_timeout`, which governs
  the `FIN_WAIT_2` state instead. (This is widely misreported online; Chapter 23
  states it correctly.)
- **RFC 1918 private ranges**: `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`.
- **Standard Ethernet MTU 1500**, giving **MSS 1460** after 20-byte IP and
  20-byte TCP headers.
- All **subnetting answers in Appendix C** were computed and independently
  verified.
- The **1986 NSFNET congestion collapse** (throughput falling from 32 kbit/s to
  40 bit/s) and Van Jacobson's 1988 congestion-control response.

**What changes and should be re-checked before you rely on it:** command syntax
and flags across OS versions, specific tool availability, cloud provider
behaviour, and anything about "current" adoption levels of IPv6 or HTTP/3. The
protocol fundamentals in Parts 1–7 are stable and have been for decades.

---

# Appendix F — Answers to "Check yourself"

**Ch 1** (1) Loss, reordering, duplication, corruption, delay, path changes.
(2) Each problem can be solved once, independently, and each piece can be replaced
without rewriting everything else. (3) The internet is the network; the web is one
application (HTTP) running on it.

**Ch 2** (1) Each layer wrapping the layer above's data in its own header.
(2) The link-layer (Ethernet) header is rewritten at every hop; the IP header
stays the same end to end. (3) Because routers only read the IP layer — a new
transport protocol above it needs no changes below it.

**Ch 3** (1) Link (Ethernet/Wi-Fi), Internet (IP), Transport (TCP/UDP),
Application (HTTP/DNS). (2) It promises nothing about delivery, order, timing, or
duplication — best effort only. (3) Keep the network core simple and put
intelligence at the endpoints; it lets you innovate at the edges without changing
the middle. (4) Many link layers at the bottom, many applications at the top, and
exactly one protocol (IP) in the narrow middle.

**Ch 3 Practice** 1 = Link · 2 = Internet · 3 = Application · 4 = Transport ·
5 = Link · 6 = Transport (TCP) · 7 = Application (DNS name) · 8 = Internet.

**Ch 4** (1) DHCP, DNS, route lookup, ARP, TCP handshake, NAT, routing across the
internet, TLS handshake, HTTP request, response and close. (2) Step 1, DNS.
(3) TCP (step 4) and TLS (step 7). (4) DNS resolution, TCP/TLS setup (distance),
server processing time, or the transfer itself.

**Ch 5** (1) 6 bytes; the first three identify the manufacturer (the OUI).
(2) Tells the receiver which protocol is inside so it can hand the payload to the
right code. (3) It's silently discarded — the link layer does not retransmit.
(4) Your router's (the next hop), not the website's.

**Ch 6** (1) An IP address into a MAC address, on the local network only.
(2) You don't know who has it, so you must ask everyone; but the answerer knows
exactly who asked. (3) No — ARP frames have no IP header and no TTL, so routers
don't forward them. (4) Your default gateway's.

**Ch 7** (1) By noting the source MAC of every frame and the port it arrived on.
(2) Floods it out every port except the one it came in on. (3) No; a router (or a
VLAN boundary) does. (4) Ethernet frames have no TTL, so a looping broadcast
circulates forever and saturates the network.

**Ch 8** (1) Radio is unreliable and half-duplex, so the link layer retries
locally rather than leaving it to TCP. (2) Airtime is shared; a device
transmitting slowly occupies the channel for longer, starving everyone else.
(3) Jitter — the variability of latency. It's larger on Wi-Fi because of
contention and retries.

**Ch 9** (1) 32 bits; each octet is 8 bits, and the largest 8-bit number is 255.
(2) A network part and a host part. (3) They're reserved: all-zero host bits names
the network itself; all-one host bits is the broadcast address. (4) DHCP failed —
it never got an address.

**Ch 10** (1) How many leading bits of the address are the network part.
(2) 256 minus the mask value. (3) Block size 64; `.100` is in the 64-block
(64–127), `.140` is in the 128-block (128–191) — **different networks**.
(4) `/25` has more network bits, so fewer host bits, so fewer addresses.

**Ch 11** (1) Because a subnet must start on a multiple of its own block size;
allocating small ones first strands unusable gaps. (2) 30 → /27; 60 → /26;
2 → /30. (3) No — a `/28` has block size 16, so it must start at 0, 16, 32, 48…
`.20` isn't a multiple of 16.

**Ch 12** (1) `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`. (2) The mapping
between an inside address+port and an outside address+port; ports let it
distinguish many devices sharing one public address. (3) There's no existing NAT
mapping for an unsolicited inbound packet, so the router has nowhere to send it.
(4) Your ISP NATs you as well, so your "public" address is shared and changeable —
port forwarding becomes impossible.

**Ch 13** (1) "Send anything I don't otherwise recognise to 192.168.1.1."
(2) The `/24` — longest prefix match wins. (3) Changes the link-layer (MAC)
header and decrements the TTL; leaves the IP addresses and payload alone.
(4) Exchanging reachability information *between* independent networks — it's how
the global routing table is built.

**Ch 14** (1) By sending packets with TTL 1, 2, 3… and collecting the "time
exceeded" errors each router sends back. (2) No — that router simply chose not to
reply. (3) Routers deprioritise traffic addressed to themselves; only believe it
if the latency persists on all later hops. (4) Path MTU Discovery (causing hangs)
and, on IPv6, address resolution (breaking the network entirely).

**Ch 15** (1) MTU 1500; MSS 1460 after 20 bytes of IP and 20 of TCP header.
(2) Drops it and sends back an ICMP "fragmentation needed" message with the
correct size. (3) The connection establishes, small requests work, then it hangs
when a full-size packet is sent. (4) They add an outer header, reducing the space
available for your data below the standard 1500.

**Ch 16** (1) IPv4's 4.3 billion addresses ran out. (2) One run of all-zero
groups, compressed; only once, because otherwise you couldn't tell how many groups
each `::` replaced. (3) A link-local address — every IPv6 interface has one
automatically; it's used for local plumbing and is never routed. (4) Neighbor
Discovery (IPv6's ARP) runs over ICMPv6, so blocking it stops address resolution
entirely.

**Ch 17** (1) DISCOVER, OFFER, REQUEST, ACK. (2) The client has no address yet, so
it can't send a unicast packet; and it doesn't know where the server is.
(3) DHCP failed. (4) Subnet mask, default gateway, DNS servers (also lease time,
domain, NTP).

**Ch 18** (1) Ask a root server → referred to the `.com` servers → referred to
`example.com`'s servers → get the answer. (2) How long an answer may be cached;
you must lower it *before* a change and wait out the old TTL, or people keep
getting the old answer. (3) The recursive resolver does the legwork and caches;
the authoritative server holds the real data for one domain. (4) When the answer
is too large for a UDP packet (the truncated flag is set). (5) The full delegation
chain from the root down to the authoritative answer.

**Ch 19** (1) Because many connections share the same server IP and port; only the
client's address and port distinguish them. (2) `0.0.0.0` listens on all
interfaces (reachable from the network); `127.0.0.1` listens only on loopback
(reachable only from that machine). (3) By the 4-tuple — each client has a
different source IP and/or port. (4) Refused = the host replied with RST (nothing
listening). Timeout = nothing came back at all (dropped, or host unreachable).

**Ch 20** (1) Port numbers and a checksum. (2) A late voice packet is useless;
waiting for a retransmission would make the call worse, not better. (3) The client
has no IP address yet, so it can't complete a TCP handshake. (4) Sequence
numbers, acknowledgements, retransmission, ordering, and rate control.

**Ch 21** (1) After two, the client knows the server heard it, but the server
doesn't yet know the client heard the reply. (2) SYN = "synchronise, here's my
starting sequence number"; ACK = "I received up to here". (3) The host is
reachable but nothing is listening on that port. (4) Something is silently
dropping your packets (a firewall), or the host is unreachable. (5) MSS, window
scaling, SACK support, timestamps.

**Ch 22** (1) Bytes. (2) An ACK repeating the same number, meaning "I'm still
missing that byte"; after three, the sender fast-retransmits immediately without
waiting for a timeout. (3) There are no later packets to generate duplicate ACKs,
so recovery waits for the retransmission timer (~200 ms+). (4) It tells the sender
exactly which ranges arrived, so only the true gaps are resent — critical on lossy
or high-latency paths.

**Ch 23** (1) Each direction is closed independently, so each side sends its own
FIN and gets its own ACK. (2) The side that closes first; about 60 seconds on
Linux. (3) Your application received the peer's FIN but never called `close()` —
a socket leak in your code. (4) No — it controls `FIN_WAIT_2`. TIME_WAIT is
hardcoded.

**Ch 24** (1) How many more bytes the receiver is willing to accept right now.
(2) Because throughput ≈ window ÷ RTT, and 65,535 bytes ÷ 0.1 s ≈ 5 Mbit/s.
(3) The receiving *application* isn't reading fast enough — it's not a network
problem. (4) 100 Mbit/s × 0.04 s = 4 Mbit = **500 KB**.

**Ch 25** (1) Packet loss. (2) It *starts* small; it then doubles every round
trip, which is exponential, not slow. (3) It finishes during the slow-start ramp
before reaching full speed. (4) Each connection gets its own congestion window, so
N connections claim roughly N times the bandwidth share.

**Ch 26** (1) Buffers small writes until a full segment accumulates or outstanding
data is acknowledged, to avoid tiny wasteful packets. (2) Delays the
acknowledgement briefly, hoping to piggyback it on returning data. (3) The sender
withholds the second small write waiting for an ACK; the receiver withholds the
ACK waiting for data to piggyback on — both wait until the timer fires.
(4) Set `TCP_NODELAY`, or write the whole message in a single call.

**Ch 27** (1) Bandwidth is capacity per second; latency is delay. (2) Page loads
are dominated by round trips, which depend on distance, not capacity.
(3) Excessive latency from oversized buffers filling at the bottleneck; test by
pinging while running a large transfer. (4) So the queue forms in your smart
queue-management device rather than in the dumb modem buffer you can't control.

**Ch 28** (1) A request line, headers, a blank line, and an optional body.
(2) 4xx = the client's request was wrong; 5xx = the server failed. (3) One IP
address hosts many websites; the server needs to know which one you want.
(4) It reuses the TCP (and TLS) connection, saving one or two round trips per
request.

**Ch 29** (1) Confidentiality, integrity, authentication. (2) Without it you could
have a perfectly encrypted connection to an attacker. (3) Server Name Indication —
the hostname you're requesting, sent in plaintext so one IP can serve many sites.
(4) Chain leads to a trusted root; the name matches; it's within its validity
dates; it isn't revoked.

**Ch 30** (1) Only one request at a time per connection, forcing browsers to open
many connections. (2) TCP delivers bytes strictly in order, so one lost packet
stalls *all* multiplexed streams on that connection. (3) A brand-new protocol
number would be blocked by middleboxes; UDP passes everywhere. (4) The connection
is identified by an ID rather than the IP/port 4-tuple, so it survives moving
between networks.

**Ch 31** (1) Disables name resolution — otherwise tcpdump generates DNS traffic
that appears in its own capture and slows it down. (2) SYN; SYN+ACK; PSH+ACK
(data); RST (reset/refused). (3) A single capture only shows what arrived there;
comparing both ends proves whether a packet was sent and lost.
(4) `'tcp[tcpflags] & (tcp-syn|tcp-ack) == tcp-syn'`.

**Ch 32** (1) Conversations; Protocol Hierarchy; Follow TCP Stream; Expert
Information; set the time column to "since previous displayed packet".
(2) The *receiving application* is too slow to read its buffer — not a network
problem. (3) Capture filters (tcpdump syntax) decide what's recorded; display
filters (Wireshark syntax) decide what's shown. (4) Analyze → Expert Information.

**Ch 33** (1) Listening TCP and UDP sockets, numerically. (2) `rtt` = measured
round-trip time (distance); `cwnd` = congestion window (how much it's allowed to
send); `retrans` = retransmissions (loss). (3) The application isn't closing
sockets — a leak. (4) `sudo ss -tulnp | grep :8080` or `sudo lsof -i :8080`.

**Ch 34** (1) Ping by IP versus by name — it splits the problem into "DNS" or
"not DNS" instantly. (2) It's DNS. (3) Refused proves you reached the host, so the
network path works and only the service is missing. (4) To distinguish "never
sent" from "sent but lost in transit".

**Ch 52** (1) NIC → RX ring → softirq/GRO → netfilter PREROUTING (conntrack,
DNAT) → routing → INPUT → TCP (SYN queue, accept queue) → socket receive buffer
→ `read()`. Drops can happen at the ring (`rx_missed`), conntrack (table full),
the accept queue (`ListenOverflows`), or the receive buffer. (2) Forwarding lets
the host route the namespace's packets at all. MASQUERADE rewrites the private
source to the host's address so replies can find their way back. (3) The accept
queue is full: the application isn't calling `accept()` fast enough. It's a
process problem, not a network one. (4) Every new flow that needs tracking: NAT,
Docker port publishing, Kubernetes Services, and stateful firewall rules all
start dropping new connections.

**Ch 53** (1) A stateless filter doesn't remember the inbound SYN, so the
outbound SYN-ACK to the client's ephemeral port needs its own allow rule.
(2) All of them on the path: "refused" means the SYN reached the host and the
kernel answered. (3) `policy drop` makes anything not explicitly allowed be dropped;
the `established,related` rule lets replies through. It goes first so that the
vast majority of packets (replies) are accepted after one cheap lookup.
(4) Tests can be run after every change and from every segment, so drift and
forgotten paths are caught automatically.

**Ch 54** (1) Outbound it's a routing table (which peer to send a destination
to); inbound it's an access list (which source addresses a peer may use).
(2) WireGuard over IPv6 adds 80 bytes (40 IP + 8 UDP + 32 WireGuard); 1500 − 80
= 1420 works over either IP version. (3) It rewrites the MSS in TCP SYNs so TCP
never sends segments too big for the tunnel; it can't help UDP or protocols that
don't negotiate MSS. (4) MTU: small packets (SSH keystrokes) fit, large ones
(a clone's data) don't. Test with `ping -M do -s` and a large download through
the tunnel.

**Ch 55** (1) At the registrar, because delegation is NS records in the
*parent* zone, which the registry publishes. (2) A/AAAA records for nameservers
named inside the zone they serve (`ns1.example.com` for `example.com`); without
them the lookup would be circular. (3) Negative caching: the earlier NXDOMAIN is
cached for the lower of the SOA minimum and the SOA record's TTL. (4) Lower the
TTL to ~60 s at least one old-TTL period before, switch, verify, keep the old LB
running, raise the TTL after a day, and remove the old LB after a week. (5) A
DNSSEC failure resolves with `dig +cd` but SERVFAILs without it; a nameserver
outage fails both ways and `dig @ns1... +norecurse` times out.

**Ch 56** (1) The server isn't sending its intermediate certificate (browsers
may fetch or cache it; other clients don't), or the client uses a different
trust store. (2) Change the certificate for new handshakes without restarting
or dropping existing connections, and choose a certificate per SNI name.
(3) Fully automated issuance, deployment, and monitoring; no manual renewals
anywhere, including load balancers, CDNs, keystores, and appliances. (4) For
example: CDN edge, cloud load balancer, Kubernetes ingress, the service's own
TLS listener, a Java keystore, a mobile app pin.

**Ch 57** (1) gRPC multiplexes all calls over a long-lived HTTP/2 connection,
and an L4 balancer balances connections, so each client sticks to one backend.
(2) Readiness: "can I serve real traffic now?", including critical
dependencies and draining state. Liveness must not depend on external
dependencies, or a dependency outage restarts every instance. (3) Fail
readiness; wait for the LB to notice; finish in-flight requests within a
deadline; exit. Skipping step 2 means the LB keeps sending new requests to a
server that has stopped accepting them, so they fail. (4) So the LB always
closes idle connections first and never sends a request on a connection the
backend is closing (the 502 race). (5) Trust it only from known proxy addresses,
and use the right-most entry added by your own trusted proxy.

**Ch 58** (1) It isn't assigned to any interface. It exists only as DNAT rules
(iptables/IPVS/nftables) or eBPF programs on each node that rewrite it to a
ready pod IP. (2) Names with fewer than 5 dots try every search domain first:
`api.github.com` costs three NXDOMAIN lookups (each for A and AAAA) before the
real one. (3) When you need the client's source IP or want to avoid the extra
hop; the cost is uneven balance and only nodes running pods receiving traffic.
(4) A policy with `podSelector: {}`, `policyTypes: [Egress]`, and a single
egress rule allowing UDP/TCP 53 to the kube-dns pods.

**Ch 59** (1) 500 Mb/s × 0.08 s = 40 Mb = 5 MB. (2) 300 × 0.04 = 12 concurrent
connections (plus headroom). (3) Each connection is limited per-connection
(window, buffers, or loss), not by the path's capacity. (4) Setting it pins
the buffer and disables Linux's autotuning, which would usually grow it much
larger.

**Ch 60** (1) IC (doesn't debug), operations (doesn't do external comms),
communications (doesn't speculate on cause), scribe (doesn't filter). (2) A
rollback changes state and often destroys the evidence (connections, counters,
logs in memory) needed to find the root cause. (3) For example: "Since 14:02 UTC,
orders-service in eu-west cannot open TCP connections to db-1:5432 (SYNs time
out); other services reach db-1; a security-group change was applied at 14:00;
need the DB team to check SG and host firewall counters." (4) It has an owner,
a date, and addresses prevention, detection, or mitigation of the *class* of
failure.

**Ch 61** (1) Address plans constrain every future connection, acquisition,
and cloud expansion for a decade or more, and fixing overlaps costs quarters of
engineering time. (2) Flow logs or counters showing zero traffic on the old
path for an agreed period. (3) Keep all allocations in a file in git, and have
CI run an overlap check (`subnet overlap`) on every change, failing the build
on any overlap. (4) Policy tests from the subnet's new location, `netcheck` to
key dependencies, a large-payload transfer (MTU), and a real business
transaction.

---

*End of guide.*

**What to do right now:** pick a route from §0.3, open `notes.md`, and start
Chapter 1. Read one chapter, explain it out loud, run the commands, and answer the
check-yourself questions from memory.

The tools and the fashions change. The four layers do not.
