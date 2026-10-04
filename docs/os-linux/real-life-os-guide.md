# Operating Systems, Linux, and Containers — A Practical Guide

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/linux/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

A go-to reference for understanding what's actually happening under your
application: what the operating system is doing for you, how to work Linux
confidently from the command line, and how Docker and Kubernetes turn "a
process on a machine" into "a workload on a fleet." Written so a software
engineer, an engineering manager, or someone starting out in SRE can all read
it and come away able to reason about a real system, not just recite
definitions.

Three phases, each beginner → intermediate → expert:

```
   PHASE 1   OS CONCEPTS         what the kernel is actually doing for
                                 every program you've ever run
   PHASE 2   LINUX FUNDAMENTALS  commands, filesystem, and processes,
                                 hands-on, until they're second nature
   PHASE 3   CONTAINERS          Docker and Kubernetes, built from the
                                 OS primitives in Phase 1 -- not magic
   PART 20   GO ON LINUX         the same kernel behaviour seen from inside
                                 a program: threads, limits, PID 1, fsync,
                                 /proc, and a container runtime you build
```

---

> **The series:** 1 OS → 2 Networking → 3 Security → 4 HTTPS walkthrough, with Go alongside.
>
> **You are here: step 1, OS & Linux.** Next: [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md) →
>
> [The full series map](#0-7-the-series-os-networking-security-https).

---

## What you will be able to do at the end

1. Explain what happens, at the OS level, between typing a command and seeing
   output — processes, memory, scheduling, I/O, all of it.
2. Read a "server is slow" incident and know which Linux tool answers which
   question (CPU? memory? disk? network?) instead of guessing.
3. Work the Linux command line fluently: navigate, search, process text,
   manage files and permissions, and write a shell script that won't
   embarrass you.
4. Explain what a container actually is (not "a lightweight VM" — the real
   mechanism) and build, run, and debug one with Docker.
5. Read and reason about a Kubernetes manifest, understand what each
   component of the cluster is doing, and diagnose why a pod won't start.
6. Connect an incident, a slow deploy, or a cost spike back to the specific
   OS/Linux/container concept causing it — the skill that actually
   separates "knows the commands" from "understands the system."
7. Know, for every topic, one real, well-documented incident or system where
   it mattered — because concepts stick when they're attached to a story.

---

## How this guide is organised

Every chapter follows the same shape:

```
   In one sentence      the whole idea, compressed
   Why this matters      the cost/reliability/incident reason to care
                         (useful even if you'll never type the command)
   How it actually works  the mechanism, explained from first principles
   Real-world example     a genuine, named incident or system -- not a
                         toy story
   Try it yourself        hands-on, on your own machine or a free VM
   Common mistakes         where people misunderstand this, and why
   Check yourself          questions (answers in Appendix E)
   Further reading         a book, a doc, a talk -- picked for clarity
```

Honesty markers, used throughout:

```
   Simplified:  a true statement that leaves out detail you don't need yet
   Debated:     practitioners genuinely disagree here
   Dangerous:   don't run this outside a throwaway VM/container
```

---

## Contents

**Part 0 — Start here**
- 0.1 Who this is for
- 0.2 The learning flow
- 0.3 Build your lab (a free, disposable Linux box)
- 0.4 The little computer-science you need
- 0.5 How to read a manifest/command you've never seen
- 0.6 How not to get stuck
- 0.7 **The series: OS → networking → security → HTTPS**
- 0.8 **The Go labs: the OS from inside a program**

**PHASE 1 — Operating system concepts, in depth**

**Part 1 — What an operating system actually does**
1. The OS as a resource manager and a liar (abstraction)
2. Kernel space vs. user space, and the system call
3. A program's whole life, from `./app` to exit code

**Part 2 — Processes, threads, and the CPU**
4. What a process actually is
5. Process states and the process lifecycle
6. `fork()` and `exec()`: how new processes are born
7. Threads: sharing everything except the stack
8. Scheduling: who gets the CPU, and for how long
9. Interrupts, context switches, and why they're not free

**Part 3 — Memory management**
10. Virtual memory: every process's private lie
11. Paging, page faults, and the TLB
12. The memory hierarchy: registers to disk, and why it exists
13. Memory allocation and the OOM killer

**Part 4 — I/O and storage**
14. The I/O stack: from `read()` to the disk platter
15. Filesystems: inodes, directories, and journaling
16. Caching, buffering, and disk scheduling

**Part 5 — Concurrency and inter-process communication**
17. Race conditions at the OS level
18. Locks, semaphores, and deadlock
19. Pipes, sockets, shared memory, and message queues
20. Signals: the OS's tap on the shoulder

**Part 6 — Boot and the kernel**
21. From power button to login prompt
22. Init systems, kernel modules, and device drivers

**PHASE 2 — Linux fundamentals**

**Part 7 — The filesystem hierarchy and navigation**
23. Everything is a file (almost)
24. The Filesystem Hierarchy Standard, and why `/etc` isn't random
25. Permissions, ownership, and links

**Part 8 — Essential commands, in depth**
26. File and directory operations you'll use daily
27. Text processing: `grep`, `sed`, `awk`, and pipes
28. Redirection and the pipeline as a design pattern
29. Finding things: `find`, `locate`, `which`, `grep -r`
30. Package management: how software gets onto a Linux box

**Part 9 — Process management in practice**
31. Watching and controlling processes: `ps`, `top`, `htop`, `kill`
32. Job control: foreground, background, and `nohup`
33. systemd and services: how "always running" actually works
34. Logs: `journalctl` and where output really goes

**Part 10 — Users, permissions, and basic hardening**
35. Users, groups, and `sudo`
36. The permission model, beyond `chmod 777`

**Part 11 — Networking on Linux**
37. How a Linux box sees the network
38. The tools: `ss`, `curl`, `dig`, `tcpdump`, and a firewall's basics

**Part 12 — Shell scripting and automation**
39. Bash scripting fundamentals
40. Cron, environment variables, and glueing things together

**Part 13 — Performance and troubleshooting**
41. Reading load average, `vmstat`, and `iostat` correctly
42. `strace`, `lsof`, and answering "why is this stuck"
43. Resource limits: `ulimit` and your first taste of cgroups

**PHASE 3 — Containerization: Docker and Kubernetes**

**Part 14 — Container fundamentals**
44. What a container actually is (namespaces + cgroups + a filesystem)
45. Containers vs. virtual machines, honestly compared
46. Images and layers: why containers start in milliseconds

**Part 15 — Docker in depth**
47. Running your first containers, and understanding what happened
48. Writing a good Dockerfile
49. Volumes and data: what survives a container's death
50. Container networking: how containers talk to the world
51. Docker Compose: running more than one container together
52. Registries and image security basics

**Part 16 — Kubernetes fundamentals**
53. Why Kubernetes exists (the problem Docker alone doesn't solve)
54. The control plane and the node: Kubernetes' architecture
55. Pods, Deployments, and ReplicaSets
56. Services and how traffic finds a pod
57. Configuration: ConfigMaps, Secrets, and environment
58. Scheduling and scaling

**Part 17 — Kubernetes in practice**
59. Reading cluster state and debugging a pod that won't start
60. Deployment patterns: rolling updates, health checks, rollbacks
61. Helm and packaging Kubernetes applications

**Part 18 — Capstone projects**
62. Project 1: build and profile a small program, OS-level
63. Project 2: containerize a real application properly
64. Project 3: deploy that application to Kubernetes with health checks
65. Project 4: diagnose and fix a planted incident, end to end

**Part 19 — Advanced Linux operations and expert practice**
66. Production Linux mental models: fleet, drift, and failure domains
67. systemd beyond start/stop: dependencies, timers, sandboxing, and debugging
68. Storage operations: LVM, RAID, filesystem growth, backups, and recovery
69. Advanced Linux networking: namespaces, routing, nftables, and packet paths
70. Observability and performance: perf, flame graphs, eBPF, and OpenTelemetry
71. Linux security hardening: SELinux/AppArmor, auditd, capabilities, and least privilege
72. Incident response and refactoring on real Linux systems

**Part 20 — Systems programming in Go: the OS from inside a program**
73. Go meets the kernel: system calls, threads, and the runtime
74. Processes from Go: exec, exit codes, signals, and PID 1
75. CPU limits: GOMAXPROCS, cgroups, and throttling
76. Memory limits: the Go heap, GOMEMLIMIT, and the OOM killer
77. Files that survive crashes: page cache, fsync, and atomic replacement
78. File descriptors and the netpoller: holding thousands of connections
79. Reading /proc from Go: build your own ps
80. A container runtime in 150 lines of Go
81. Profiling Go on Linux: pprof, the execution tracer, and perf
82. Capstone: shipping a well-behaved Go service (image, systemd, Kubernetes)

**Part 21 — Where to go next**
83. Honest gaps and a study plan
84. Staying current

**Appendices**
- A. Glossary (plain language)
- B. Command cheat sheet
- C. Troubleshooting playbook ("the server is slow" flowchart)
- D. Further reading, and how this guide's facts were verified
- E. Answers to "Check yourself"

---

# Part 0 — Start here

## 0.1 Who this is for

| You are… | This guide fits? |
|---|---|
| A software engineer who's never had to think below the application layer | **Yes — start at Chapter 1** |
| An engineering manager who wants to understand what your team means by "OOM killed the pod" | **Yes** — read the "Why this matters" boxes even if you skip the hands-on parts |
| Someone moving into an SRE/platform role | **Yes** — this is close to the actual on-call curriculum |
| Already comfortable with `ps`, `grep`, and Docker | Skim Phase 1-2, start deep at Part 14 |
| Looking for kernel-internals/driver-development depth | This guide explains mechanisms clearly but stops short of writing kernel code — see Part 19 for where that lives |

## 0.2 The learning flow

Three phases, and each one **builds the vocabulary the next one assumes**:

```
   PHASE 1 (OS concepts)     teaches you what a "process," "page fault," or
        |                    "file descriptor" actually IS
        v
   PHASE 2 (Linux)           teaches you the COMMANDS that ask the kernel
        |                    about those exact things (`ps`, `free`, `lsof`)
        v
   PHASE 3 (Containers)      shows you that a "container" is just Linux
        |                    primitives from Phase 1, used in a clever way
        |                    -- and Docker/Kubernetes commands are, again,
        |                    just asking the kernel (via a daemon) to do
        v                    things you already understand.
   PART 20 (Go on Linux)     has you write the programs that make those
                             primitives visible -- and build a container
                             runtime yourself (optional, needs basic Go).
```

Skipping Phase 1 and going straight to `docker run` is possible — plenty of
engineers do exactly that — but "why did my container get OOM-killed" or "why
is this pod stuck in `CrashLoopBackOff`" stops being mysterious only once you
know what's underneath. That's the whole bet this guide makes: **go down one
level, and everything above it gets easier to debug, not harder to learn.**

The loop for every chapter:

```
   READ the concept  ->  see the REAL-WORLD example  ->  TRY IT on your own
   box  ->  answer the CHECK YOURSELF questions  ->  move on
```

Don't skip "try it yourself." Reading about `fork()` is not the same as
watching a process count double.

## 0.3 Build your lab

You need a disposable Linux machine you're not afraid to break. Any of these
work; pick the cheapest one to start:

```
   [ ] A cloud free-tier VM (AWS/GCP/Azure/Oracle all have a permanently-free
       small instance) running Ubuntu 22.04/24.04 or Debian.
   [ ] A local VM: UTM (Mac, free), VirtualBox, or Multipass
       (`multipass launch` gives you a throwaway Ubuntu VM in under a
       minute).
   [ ] WSL2 on Windows -- a real Linux kernel, good enough for almost
       everything in Phases 1-2.
   [ ] Docker Desktop / Docker Engine + Kubernetes -- for Phase 3. `kind`
       (Kubernetes in Docker) or `minikube` gives you a free, disposable
       cluster on your laptop.

   Minimum install for the whole guide:
     build-essential (or gcc), python3, htop, strace, lsof, net-tools,
     iproute2, dnsutils, curl, git, docker, kind or minikube, kubectl
```

```bash
# Ubuntu/Debian, one shot:
sudo apt update && sudo apt install -y build-essential python3 htop strace \
  lsof net-tools iproute2 dnsutils curl git
```

**Never practise destructive commands (`rm -rf`, `kill -9 1`, disk
formatting, `fork()` bombs) anywhere except this throwaway box.**

## 0.4 The little computer science you need

Four ideas, and you're equipped for everything that follows.

```
   1. A COMPUTER ONLY DOES ONE THING AT A TIME, PER CORE.
      Everything that looks "simultaneous" -- 200 browser tabs, a hundred
      server processes -- is actually the CPU switching between them very
      fast (Chapter 8). "Concurrency" (juggling many things) is not the
      same as "parallelism" (doing many things at the exact same instant,
      which needs multiple cores).

   2. MEMORY IS A LADDER, AND EACH RUNG IS ORDERS OF MAGNITUDE SLOWER.
      CPU register (~1 cycle) -> CPU cache (~few cycles) -> RAM (~100s of
      cycles) -> SSD (~10,000s of cycles) -> spinning disk / network
      (~millions of cycles). Most of performance engineering is "keep the
      hot data as high on this ladder as possible" (Chapter 12).

   3. EVERYTHING IS A NUMBER, AND NUMBERS ARE CHEAP TO PASS AROUND.
      A process is identified by a number (a PID). A file, once opened, is
      identified by a number (a file descriptor). This is why Unix tools
      compose so well -- they're all just passing small numbers and
      streams of bytes around.

   4. THE KERNEL IS THE ONLY THING TOUCHING THE HARDWARE.
      Your program never talks to the disk, the network card, or memory
      management hardware directly. It ASKS the kernel to do it (a system
      call, Chapter 2) and waits. Every "how does X work" question in this
      guide eventually resolves to "the kernel did it, here's how."
```

## 0.5 How to read a manifest/command you've never seen

You will constantly meet commands and config you haven't memorised. The
practised move, every time:

```
   1. Read the THING being acted on (a file? a process? a container? a
      Kubernetes object?).
   2. Read the VERB (create, list, delete, describe, exec).
   3. Check `--help` or `man <command>` for the flags you don't recognise
      -- don't guess. `man ls`, `kubectl explain pod.spec`.
   4. Run it in your lab, read the ACTUAL output, and connect it back to
      the mental model from Phase 1.
```

## 0.6 How not to get stuck

```
   [ ] Time-box confusion to 20 minutes, then move on -- some things click
       two chapters later, once you've seen the NEXT layer.
   [ ] `man <command>` and `<command> --help` before a search engine --
       building that reflex now pays off for the rest of your career.
   [ ] When a command fails, READ THE ERROR MESSAGE FULLY before acting.
       Linux/kernel error messages are usually precise.
   [ ] Break things on purpose in your lab. You cannot hurt anything by
       experimenting on a disposable VM. Use that freedom.
```

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
| 1 | **Operating Systems, Linux, and Containers** ← you are here | the machine every request starts and ends on: processes, memory, files, sockets, containers, Kubernetes | 9 Go labs ([Part 20](#part-20-systems-programming-in-go-the-os-from-inside-a-program)) |
| 2a | [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md) | one click followed through all seven layers; the map for everything after it | concept map, 1 hour |
| 2b | [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md) | how machines talk: addressing, routing, TCP, TLS, packet capture, network operations | 23 Go labs ([§0.8](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself)) |
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

## 0.8 The Go labs: the OS from inside a program

The commands in this guide show the operating system from the *outside*:
`ps`, `strace`, `free`, `docker stats`. **Part 20** shows it from the
*inside*, with Go programs that make the kernel's behaviour visible: how many
OS threads a goroutine really costs, what happens when a container's memory
limit meets Go's garbage collector, why PID 1 must reap zombies, how to write
a file that survives a power cut, and a working container runtime in about
150 lines.

Every lab was run on Linux (in Docker on an Apple-silicon Mac) while writing
this guide, and the outputs shown are real. Most labs are **Linux-only**
(`//go:build linux`). On a Mac, cross-compile and run them in a container:

```bash
mkdir -p ~/oslabs && cd ~/oslabs && go mod init oslabs
mkdir threads    # save the lab as threads/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/threads ./threads   # GOARCH=amd64 on Intel
docker run --rm -v "$PWD/bin":/b debian:bookworm /b/threads
```

| Lab | Chapter | What it shows |
|---|---|---|
| `threads` | 73 | Goroutines vs OS threads: sleeping, network-blocked, syscall-blocked |
| `pid1` | 74 | Zombies piling up under a PID 1 that doesn't reap, and three fixes |
| `gomaxprocs` | 75 | Go's container-aware GOMAXPROCS, and CPU throttling when it's wrong |
| `memlimit` | 76 | An OOM kill at 300 MB, prevented by GOMEMLIMIT, and what that costs |
| `atomicwrite` | 77 | Durable, atomic file replacement, and what `fsync` costs |
| `fdlimit` | 78 | RLIMIT_NOFILE, EMFILE, and Go raising its own limit |
| `gops` | 79 | `ps`/`top` rebuilt from `/proc` |
| `minicontainer` | 80 | Namespaces + chroot + cgroup v2 = a container |
| `profiled` | 81 | A CPU hot spot and a goroutine leak, found with pprof |
| Dockerfile, unit, manifest | 82 | Shipping a Go service: 15 MB distroless image, systemd, Kubernetes |

# Part 1 — What an operating system actually does

Lab set up, learning flow read — time to start. Before processes, memory, or
any specific subsystem: what is an operating system actually *for*, and why
does it need to lie to every program running on it? Everything in the next
fifteen parts is a specific instance of that lie and the machinery that
maintains it.

## Chapter 1 — The OS as a resource manager and a liar

### In one sentence

An operating system's whole job is to let many programs share one computer's
hardware safely, by *pretending* to each program that it has the machine to
itself.

### Why this matters

Every "why is my server slow" investigation, every container resource limit,
every "noisy neighbour" incident traces back to this one idea: programs don't
actually get dedicated hardware — they get a convincing illusion of it, arbitrated
by the OS, and that arbitration can go wrong or run out.

### How it actually works

```
   WITHOUT AN OS: your program would need to know EXACTLY which physical
   RAM address is free, EXACTLY how to talk to this specific disk
   controller, and EXACTLY how to avoid another program overwriting its
   memory at the same instant. Every program would need to be its own
   tiny operating system. This was, roughly, how computing worked before
   the 1960s-70s: one program, one computer, no sharing.

   WITH AN OS, three lies are told to every program, consistently:
     1. "You have the CPU to yourself."
        (Reality: the scheduler, Chapter 8, is switching between you and
         every other runnable program many times per second.)
     2. "You have a huge, private block of memory, starting at address 0."
        (Reality: virtual memory, Chapter 10, maps your program's private
         address space onto scattered, shared physical RAM -- and part of
         it might not even be in RAM right now.)
     3. "Files and devices are just streams of bytes you can read/write."
        (Reality: a "file" might be on a spinning disk, an SSD, a network
         share, or not even really exist, like /dev/random.)

   The OS maintains these lies by being the ONLY thing with direct
   hardware access (Chapter 2) and mediating every request. This is what
   "operating system" means: it OPERATES the hardware ON BEHALF of every
   program that wants to use it.
```

### Real-world example

```
   Open a modern browser and look at your process list (Chapter 31) --
   Chrome alone might show 50-150 separate OS processes for a browser
   window with a handful of tabs. Google's own engineering documentation
   explains this was a DELIBERATE architecture decision: each tab (and
   each extension, and the browser UI itself) runs in its OWN OS process,
   specifically so the OS's process-isolation guarantee (Chapter 4) means
   one crashed or malicious tab can't read another tab's memory or bring
   down the whole browser.
   This is "the OS as resource manager" used as a SECURITY boundary, not
   just a performance one: Chrome is deliberately paying the OS's process-
   creation and context-switching cost (Chapter 9) in exchange for the
   isolation guarantee only the OS can provide. You'll see this exact
   trade-off again with containers in Phase 3 -- isolation has a cost, and
   engineers pay it on purpose, constantly.
```

### Try it yourself

```bash
# See how many processes are "really" running on your machine right now
ps aux | wc -l

# Watch the illusion of "you have the CPU to yourself" break down --
# start a CPU-heavy loop, then watch top show it competing for time
yes > /dev/null &
top    # note the %CPU for the `yes` process; press 'q' to quit
kill %1   # stop the background job
```

### Common mistakes

- **"More cores means my single-threaded program runs faster."** A program
  using one thread can only ever use one core at a time, no matter how many
  cores the machine has (Chapter 7 explains why, and when more cores DO help).
- **Thinking "virtual memory" means "not real."** It's real memory your
  program can actually use — "virtual" describes the ADDRESSING scheme, not
  whether the memory works (Chapter 10).
- **Assuming isolation is free.** Every layer of isolation (processes,
  containers, VMs) costs CPU and memory to maintain. More isolation is a
  trade-off, not a free win.

### Check yourself

1. Name the three "lies" the OS tells every program, and what's really
   happening behind each one.
2. Why couldn't programs safely share a computer before operating systems
   existed?
3. Why did Chrome's engineers choose a multi-process architecture, and what
   did they pay for that choice?

*(Answers: Appendix E.)*

### Further reading

- **Book (free):** *Operating Systems: Three Easy Pieces* (Remzi and Andrea
  Arpaci-Dusseau, `ostep.org`) — the best modern, free OS textbook; this
  guide's Phase 1 follows its spirit closely.
- **Article:** "Multi-process Architecture" — the Chromium project's own
  design documentation on why Chrome uses one process per tab.
- **Video:** Computerphile, "What is an Operating System?" — a clear, short
  visual introduction.

---

## Chapter 2 — Kernel space vs. user space, and the system call

### In one sentence

The kernel is the one piece of software allowed to touch hardware directly,
everything else (your shell, your browser, your database) runs in "user
space" and must formally *ask* the kernel to do anything that touches the
outside world.

### Why this matters

When you read that a vulnerability lets an attacker "escape to the kernel" or
"gain root," this is the boundary being crossed. It's also why a crashing
application usually doesn't take down your whole machine, but a crashing
kernel (a "kernel panic") does.

### How it actually works

```
   THE CPU ITSELF enforces this boundary, via "privilege rings" (or
   "modes") built into the hardware -- it's not just a software policy
   the kernel chooses to follow, it's physically enforced by the
   processor:
     RING 0 / KERNEL MODE   full access: can execute any instruction,
                            touch any memory, talk to any device.
     RING 3 / USER MODE     restricted: cannot execute privileged
                            instructions or touch hardware directly.
   Your shell, your text editor, your web server -- all of it runs in
   user mode. Only the kernel runs in kernel mode.

   SO HOW DOES A PROGRAM READ A FILE, IF IT CAN'T TOUCH THE DISK?
   It makes a SYSTEM CALL (a "syscall"): a formal request to the kernel,
   using a special CPU instruction that deliberately, safely switches the
   processor from user mode to kernel mode, runs a specific bit of kernel
   code to service the request, and switches back.

     your program:  open("data.txt", O_RDONLY)
                         |  (this C library function wraps a syscall)
                         v
     CPU switches to kernel mode (a "trap")
                         |
     kernel: checks permissions, finds the file, sets up a file
             descriptor (Chapter 23), switches back to user mode
                         |
     your program: receives a file descriptor number, continues

   THIS IS SLOWER than a plain function call -- crossing the boundary has
   real, measurable overhead (a "context switch," Chapter 9) -- which is
   exactly why well-written programs try to minimise the NUMBER of
   syscalls they make (buffering writes instead of writing one byte at a
   time, for example).
```

**The syscalls you'll meet constantly**, so the names stop being scary:

```
   read() / write()      move bytes between a program and a file/socket
   open() / close()      get/release a file descriptor
   fork() / exec()       create a new process / replace a process's code
                        (Chapter 6)
   mmap()                 map a file or memory region into your address
                        space directly
   socket() / connect()   networking (a socket is just another kind of
                        file descriptor)
   wait()                 a parent process waiting for a child to finish
```

### Real-world example

```
   CVE-2016-5195, nicknamed "Dirty COW," is a textbook illustration of
   why this boundary has to be enforced PERFECTLY. It was a RACE
   CONDITION (Chapter 17) in how the Linux kernel handled a specific
   memory-management optimisation called copy-on-write, reachable
   through an ordinary, unprivileged sequence of syscalls
   (`open()`/`mmap()`/`write()`, timed carefully against each other).
   An attacker with only a normal, unprivileged user-space login could
   exploit the race to WRITE TO MEMORY THAT SHOULD HAVE BEEN READ-ONLY --
   including files they had no permission to modify, like the system's
   own password database -- effectively promoting themselves to root.
   The bug had existed, unnoticed, in the kernel for close to a DECADE
   before it was found and fixed in 2016, precisely because the
   kernel/user-space boundary is exactly the kind of code where a subtle
   timing bug is both extremely hard to find and extremely damaging once
   found -- it's the one piece of software every single program on the
   machine implicitly trusts completely.
```

### Try it yourself

```bash
# See the actual syscalls a simple command makes -- this demystifies
# "system call" faster than any explanation
strace -c ls > /dev/null      # -c summarises counts per syscall
strace ls 2>&1 | head -20     # or see them one by one, in order

# Watch the boundary in real time: this counts context switches
# (roughly, crossings between your program and the kernel)
/usr/bin/time -v ls > /dev/null 2>&1 || time -v ls  # (flags vary by system)
```

### Build it in Go
Go programs make system calls like any other program, but the Go *runtime*
makes many on your behalf (`futex`, `epoll_pwait`, `clone`, `mmap`). Run any Go
binary under `strace -f -c` and compare the summary with a C program's.
[Chapter 73](#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime) explains what you're seeing and measures
when a goroutine costs an OS thread.

### Common mistakes

- **"User mode is for regular users, kernel mode is for root."** No — this is
  about the CPU's privilege level for a given piece of CODE, not about which
  human is logged in. Even `root`'s shell runs in user mode; it just has
  permission to ask the kernel to do more things.
- **Thinking syscalls are instant.** They're fast, but far slower than a
  normal function call, because of the mode switch. This is why high-
  performance code (and Chapter 16's buffering) works hard to reduce syscall
  count.
- **Confusing "kernel" with "operating system."** The kernel is the core;
  "the OS" usually also includes user-space tools (a shell, system daemons,
  package managers) built on top of it.

### Check yourself

1. What physically enforces the separation between kernel mode and user
   mode?
2. Name four syscalls and what each one does.
3. Why are syscalls slower than ordinary function calls?
4. In the Dirty COW example, what made this specific class of bug so
   dangerous?

### Further reading

- **Book (free):** *OSTEP*, the "Direct Execution" and "Limited Direct
  Execution" chapters — the clearest explanation of the privilege-ring
  mechanism available.
- **Reference:** `man 2 syscalls` on any Linux box lists every syscall the
  kernel supports, in one page.
- **Advisory:** the original Dirty COW disclosure (`dirtycow.ninja`) — a
  clear, non-hyped technical writeup of the bug.

---

## Chapter 3 — A program's whole life, from `./app` to exit code

### In one sentence

Running a program is a specific, well-defined sequence — the shell asks the
kernel to create a new process, load your code into it, and start it running;
the kernel tracks it until it finishes; and everyone who cares can find out
how it ended.

### Why this matters

"Why didn't my script's error get caught," "why is this process a zombie,"
and "how does Docker know my container exited" are all answered by
understanding this one lifecycle.

### How it actually works

```
   1. YOU TYPE: ./myapp arg1 arg2

   2. YOUR SHELL calls fork() (Chapter 6) -- the kernel clones the
      shell's own process into a near-identical CHILD process.

   3. THE CHILD calls exec() -- this REPLACES the child's code and
      memory with `myapp`'s code, but keeps the same process ID. The
      kernel loads the executable from disk, sets up its memory layout
      (Chapter 10), and starts running its `main()` function.

   4. THE PROCESS RUNS, doing whatever it does, making syscalls
      (Chapter 2) whenever it needs the kernel's help, getting scheduled
      on and off the CPU (Chapter 8) the whole time.

   5. THE PROCESS ENDS, one of three ways:
        - it returns from main() / calls exit(N) deliberately
        - it's killed by a SIGNAL (Chapter 20) -- e.g. you pressed
          Ctrl-C, or something sent SIGKILL
        - it crashes (an illegal memory access, division by zero, etc.
          -- the kernel delivers this as a signal too)
      Either way, the kernel records an EXIT STATUS: a small number, by
      convention 0 for "success" and non-zero for "something went wrong"
      (the SPECIFIC non-zero value is up to the program).

   6. THE PARENT (your shell) can call wait() to retrieve that exit
      status. Until it does, the finished process sits in the process
      table as a ZOMBIE -- it's done running, but the kernel keeps its
      exit status around until the parent asks for it. (A parent that
      never calls wait() leaks zombies -- Chapter 5 covers this properly.)
```

```bash
   $ ./myapp arg1 arg2
   $ echo $?              # the shell's own wait() already happened;
                          # $? holds the exit status it collected
   0
```

### Real-world example

```
   Every CI/CD pipeline in existence (GitHub Actions, Jenkins, GitLab CI)
   is built on EXACTLY this mechanism: a build step is "successful" if
   and only if the process it ran exited with status 0. This is why a
   test script that finds a failure but forgets to `exit 1` will show as
   a PASSING build even though tests actually failed -- a extremely
   common, real, and costly class of CI bug. Conversely, a script that
   `exit`s non-zero for a condition that ISN'T actually a failure (a
   common one: `grep` returns exit status 1 when it finds NO matches,
   which is often not an error at all) can fail a build for no real
   reason. Understanding "the exit code IS the entire signal" explains
   both failure modes instantly, instead of them feeling like CI being
   flaky.
```

### Try it yourself

```bash
# See exit codes in action
true; echo $?     # 0 -- the `true` command always succeeds
false; echo $?    # 1 -- the `false` command always "fails"
ls /nope; echo $?  # some non-zero number -- ls failed, file doesn't exist

# Watch a zombie process appear (harmless, in a controlled demo)
( sleep 2 & )      # subshell forks a child then exits without waiting
ps aux | grep defunct   # briefly, you may catch the zombie before
                        # your shell reaps it

# See fork+exec happen for real
strace -f -e trace=fork,execve,exit_group bash -c 'echo hi' 2>&1 | head -20
```

### Common mistakes

- **Assuming exit code 0 always means "worked correctly."** It only means
  "the program itself believes it finished without error" — that's only as
  reliable as the program's own logic.
- **Ignoring `$?` in scripts.** If you don't check it, a failed step silently
  continues as if nothing went wrong (Chapter 39 covers `set -e` and proper
  error handling).
- **Confusing a zombie process with a hung process.** A zombie has ALREADY
  finished; it's just waiting to be acknowledged. A hung process is still
  running (or stuck). `top`/`ps` show zombies with state `Z` — very different
  from a runaway process eating CPU.

### Check yourself

1. Walk through the steps from typing a command to seeing its output.
2. What's the difference between `fork()` and `exec()`, at a high level?
3. What is a zombie process, and why does it exist at all instead of just
   disappearing?
4. Why can a CI pipeline report "success" even when your tests actually
   failed?

### Further reading

- **Book (free):** *OSTEP*, "Process API" chapter — walks through
  `fork`/`exec`/`wait` with real code.
- **Reference:** `man 2 exit`, `man 2 wait`, `man 1 bash` (the "EXIT STATUS"
  section).
- **Article:** "Unix process states explained" — search for a clear diagram
  of running/sleeping/zombie/stopped states (Chapter 5 covers this in full).

---

# Part 2 — Processes, threads, and the CPU

Part 1 established that the OS multiplexes hardware among programs that each
believe they own it. This part is the first and biggest instance of that lie
made concrete: what a process actually is, how one CPU is shared between
many of them, and where a thread fits into that picture.

## Chapter 4 — What a process actually is

### In one sentence

A process is a running program plus everything the OS tracks to keep it
isolated and resumable: its own memory space, its own file descriptors, and a
record of exactly where it is in its own execution.

### Why this matters

"How many processes can this server handle" is one of the oldest, most
practical capacity-planning questions in computing, and answering it well
requires knowing what a process actually costs.

### How it actually works

```
   A process is NOT just "the code running." The kernel maintains a
   whole record for it (a "process control block") including:
     PID           a unique process ID number
     MEMORY MAP     its private virtual address space (Chapter 10)
     FILE
     DESCRIPTORS    a table of everything it has open (files, sockets,
                    pipes -- Chapter 23)
     REGISTERS      the CPU's exact state when this process was last
                    paused, so it can resume perfectly
     PARENT PID     who created it (for wait()/zombies, Chapter 3)
     PERMISSIONS     which user it runs as, what it's allowed to do
     SIGNAL HANDLERS custom behaviour for signals it's registered for
                    (Chapter 20)

   Creating a process is NOT CHEAP: the kernel has to allocate and set up
   all of the above. This cost is exactly why the industry has spent
   decades inventing lighter-weight alternatives -- threads (Chapter 7),
   thread pools, event loops, and eventually containers (Phase 3), all of
   which exist partly to avoid paying the full cost of a brand-new
   process for every unit of work.
```

### Real-world example

```
   The "C10K problem" (a name coined by Dan Kegel in a widely-read 1999
   essay) asked: can a single server handle 10,000 CONCURRENT
   connections? The traditional answer at the time -- Apache's "prefork"
   model, spawning one whole OS PROCESS per incoming connection -- hit a
   hard wall well before 10,000, because each process's memory
   footprint and scheduling overhead (Chapter 8) added up brutally at
   that scale.
   This is precisely why Igor Sysoev built NGINX (first released 2004),
   using an EVENT-DRIVEN model that handles thousands of connections
   with a small, FIXED number of worker processes instead of one per
   connection -- the same instinct as this chapter's "processes aren't
   free" lesson, engineered into one of the world's most-used web
   servers. Today's async/event-loop frameworks (Node.js, async Python,
   Go's goroutines) all descend from the same realisation: for
   high-concurrency I/O-bound work, a process (or even a thread) per
   task is too expensive.
```

### Try it yourself

```bash
# See what the kernel actually tracks for one process
echo $$                          # your shell's own PID
cat /proc/$$/status | head -20   # the kernel's own record of it
ls /proc/$$/fd/                  # its open file descriptors right now

# Estimate the real cost of a process on your machine
ps aux | awk '{sum+=$6} END {print sum/1024 " MB total RSS across all processes"}'
```

### Common mistakes

- **Thinking a process is "just code."** The kernel-tracked state (memory map,
  file descriptors, permissions) is just as much "the process" as the code
  itself.
- **Assuming more processes is always fine.** Each one has real memory and
  scheduling overhead; "spawn a process per request" doesn't scale the way
  "spawn a process per user session, ever" does.
- **Confusing PID reuse with process identity.** PIDs get RECYCLED after a
  process exits — a PID you saw an hour ago may belong to a totally different
  process now.

### Check yourself

1. Name four things the kernel tracks for a process, beyond "the code."
2. Why is creating a process more expensive than, say, calling a function?
3. What problem was NGINX's event-driven architecture specifically designed
   to solve?

*(Answers: Appendix E.)*

### Further reading

- **Essay:** "The C10K problem," Dan Kegel (`kegel.com/c10k.html`) — the
  original, still-relevant essay that shaped a generation of server design.
- **Book (free):** *OSTEP*, "Processes" chapter.
- **Reference:** `man 5 proc` — the full documentation of everything under
  `/proc/<pid>/`.

---

## Chapter 5 — Process states and the process lifecycle

### In one sentence

At any moment, every process is in exactly one of a handful of states —
running, ready to run, waiting on something, stopped, or dead-but-unclaimed —
and reading those states off `ps`/`top` is the single fastest way to diagnose
"what is this process actually doing right now."

### Why this matters

"The process is using 0% CPU but the request is hanging" and "the process
shows 100% CPU and won't respond" are DIFFERENT problems with different
causes, and the process state tells you which one you're looking at before
you dig any deeper.

### How it actually works

```
   THE CORE STATES (Linux ps/top codes in parentheses):
     RUNNING (R)          actually executing on a CPU right now, OR
                          ready and waiting for the scheduler to give it
                          a turn (both are shown as "R" -- "runnable").
     SLEEPING /
     INTERRUPTIBLE (S)     waiting for something (input, a timer, a lock)
                          and CAN be woken early by a signal. This is
                          the state of almost every idle, well-behaved
                          process on your machine right now.
     UNINTERRUPTIBLE
     SLEEP (D)             waiting for something -- usually disk or
                          network I/O -- and CANNOT be interrupted, not
                          even by SIGKILL, until the I/O completes or
                          times out. A process stuck here for a long
                          time is one of the most frustrating real
                          production problems (see below).
     STOPPED (T)           paused deliberately (Ctrl-Z, or SIGSTOP) --
                          not running, not scheduled, just frozen until
                          resumed (SIGCONT).
     ZOMBIE (Z)            finished running; just waiting for its parent
                          to collect its exit status (Chapter 3).

   THE TRANSITIONS, roughly:
     new process -> RUNNING (well, "runnable" -- waiting its turn)
     RUNNING -> SLEEPING     (it asked for I/O, or called sleep())
     SLEEPING -> RUNNING     (the thing it was waiting for happened)
     RUNNING -> STOPPED      (a stop signal arrived)
     STOPPED -> RUNNING      (a continue signal arrived)
     RUNNING -> ZOMBIE       (it exited)
     ZOMBIE -> gone          (the parent called wait())
```

### Real-world example

```
   Uninterruptible sleep (D state) is a well-known, recurring real
   operations headache, especially with NETWORK-ATTACHED STORAGE (NFS
   being the classic case, but any slow/unresponsive block storage can
   trigger it). If an NFS server becomes unreachable while a process is
   mid-read, that process can enter D state waiting for the I/O to
   complete -- and by design, the kernel won't interrupt an in-flight I/O
   operation for signal delivery, INCLUDING SIGKILL. The process becomes
   completely unkillable until the I/O either completes or the mount
   itself times out. This is documented extensively in sysadmin/SRE
   folklore and vendor knowledge bases (NetApp, Red Hat, and others all
   have support articles specifically about "D state process won't die")
   precisely because it's such a common, confusing real incident: `kill
   -9` visibly does nothing, and the only real fixes are waiting out the
   storage timeout or, in the worst case, rebooting the box.
```

### Try it yourself

```bash
# Watch a process sleep (S) vs run (R)
sleep 30 &                 # background it
ps -o pid,stat,cmd -p $!   # STAT column shows S

yes > /dev/null &
ps -o pid,stat,cmd -p $!   # STAT column shows R
kill %1 %2                 # clean up both background jobs

# Watch STOPPED and CONT in real time
sleep 100 &
kill -STOP %1; ps -o pid,stat,cmd -p $!   # T
kill -CONT %1; ps -o pid,stat,cmd -p $!   # back to S
kill %1
```

### Common mistakes

- **"0% CPU means it's not doing anything."** It might be doing plenty of
  work — just work that involves WAITING (on disk, network, a lock), not
  crunching numbers. Check the state, not just the CPU percentage.
- **Panicking at a `D`-state process immediately.** Brief D states are
  completely normal (any disk read passes through it briefly). It's only a
  problem when a process is STUCK there for an extended period.
- **Trying `kill -9` repeatedly on a D-state process.** It won't work, by
  design — SIGKILL cannot interrupt uninterruptible sleep. Find and fix the
  underlying I/O problem instead.

### Check yourself

1. List the core process states and what each one means.
2. Why can't `kill -9` stop a process stuck in D state?
3. Given "the process shows 0% CPU and the request is hanging," what's your
   first diagnostic step?

### Further reading

- **Reference:** `man 1 ps` (the "PROCESS STATE CODES" section) — the
  authoritative list for your system.
- **Book (free):** *OSTEP*, "Process API" and "Scheduling" chapters.
- **Article:** search "linux D state process troubleshooting" for current,
  practical writeups — this remains one of the most commonly searched real
  ops problems.

---

## Chapter 6 — `fork()` and `exec()`: how new processes are born

### In one sentence

Linux creates new processes in two separate steps — `fork()` clones the
current process, `exec()` replaces a process's code with something else
entirely — and understanding why they're SEPARATE explains a surprising
amount of how Unix actually works.

### Why this matters

Shells, servers, and container runtimes all use this exact pattern. Knowing
it means "how does bash run a program" and "how does a container start" stop
being separate mysteries — they're the same mechanism.

### How it actually works

```
   fork()   makes an almost-exact COPY of the calling process: same code,
            same memory contents (initially), same open file descriptors
            -- but a NEW PID. Both the original (parent) and the copy
            (child) continue running from the EXACT SAME LINE of code,
            right after the fork() call, and can tell themselves apart
            by fork()'s RETURN VALUE:
              in the parent:  fork() returns the child's PID
              in the child:   fork() returns 0
              (if it fails:   fork() returns -1)

   exec()   REPLACES the calling process's code and memory with a
            DIFFERENT program entirely, keeping the same PID and open
            file descriptors. It doesn't return on success -- the calling
            program is simply gone, replaced.

   WHY TWO SEPARATE STEPS, rather than one "run this other program" call?
   Because fork() ALONE is useful too -- and because the gap BETWEEN
   fork() and exec() is exactly where a shell does useful setup work
   (redirecting a file descriptor for `>`, changing directory, setting
   environment variables) on the CHILD ONLY, before replacing it, without
   ever affecting the parent shell itself:

     pid_t pid = fork();
     if (pid == 0) {
         // child: I am a copy of the shell right now
         // ... redirect stdout to a file, if the user wrote `> out.txt`
         // ... then:
         execve("/bin/ls", args, environ);   // now I become `ls`
     } else {
         // parent (the shell): keep running, maybe wait() for the child
     }
```

### Real-world example

```
   PostgreSQL's architecture is a clean, real-world illustration of
   fork() used WITHOUT a following exec(): PostgreSQL's documentation
   describes its "process per connection" model, where the main
   `postmaster` process FORKS a new backend process for each client
   connection -- and that child does NOT exec() into a different
   program. It simply continues running as a full copy of the postmaster,
   now dedicated to servicing one connection, benefiting from
   copy-on-write (Chapter 10) so the fork itself is cheap even though
   postgres is a large program.
   This is a DELIBERATE architectural trade-off, openly discussed in
   PostgreSQL's own documentation and by its developers: process-per-
   connection gives strong ISOLATION (one connection's crash can't
   corrupt another's memory, unlike a thread-per-connection model) at
   the cost of the per-process overhead this Part keeps returning to --
   which is exactly why connection poolers like PgBouncer exist, to
   avoid forking a whole new backend for every short-lived client
   connection in high-throughput applications.
```

### Try it yourself

```bash
# Watch fork+exec happen live, for the simplest possible command
strace -f -e trace=clone,fork,vfork,execve bash -c 'ls' 2>&1

# Write and run a tiny C program that forks (skip if you don't have gcc)
cat > forkdemo.c << 'EOF'
#include <stdio.h>
#include <unistd.h>
int main() {
    pid_t pid = fork();
    if (pid == 0) {
        printf("child: my PID is %d\n", getpid());
    } else {
        printf("parent: I created child PID %d\n", pid);
    }
    return 0;
}
EOF
gcc -o forkdemo forkdemo.c && ./forkdemo
```

### Build it in Go
Go has no `fork()` in its public API (forking a multi-threaded runtime is
unsafe); `os/exec` uses a combined fork+exec (`clone` + `execve`) under the hood.
[Chapter 74](#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1) covers exit codes, process groups,
timeouts, and the PID 1 duty of reaping orphans, with a lab that produces
zombies on purpose.

### Common mistakes

- **Thinking `fork()` runs a different program.** It doesn't — it duplicates
  the SAME program. `exec()` is the step that changes what code is running.
- **Assuming the parent and child run in a guaranteed order.** After
  `fork()`, the OS scheduler (Chapter 8) decides which runs first — your code
  must not assume one happens before the other unless you explicitly
  synchronise them.
- **Forgetting that file descriptors are shared (initially) after fork.** Both
  parent and child have the SAME open files/sockets right after forking —
  a common source of subtle bugs if both try to write to the same one.

### Check yourself

1. What does `fork()` do, and what does `exec()` do — and why are they
   separate calls?
2. In the parent process, what does `fork()` return? In the child?
3. Why does PostgreSQL fork a new process per connection instead of using
   threads, and what does that trade off?

### Further reading

- **Book (free):** *OSTEP*, "Process API" chapter — walks through real
  `fork()`/`exec()` code examples in detail.
- **Reference:** `man 2 fork`, `man 2 execve`.
- **Docs:** PostgreSQL's own documentation, "Overview of PostgreSQL Internals"
  — the process-per-connection model, explained by its own maintainers.

---

## Chapter 7 — Threads: sharing everything except the stack

### In one sentence

A thread is a separate, independently-scheduled sequence of execution that
lives INSIDE a process and shares that process's memory with every other
thread in it — which makes threads cheaper than processes and dramatically
easier to get wrong.

### Why this matters

"Thread-safe" and "race condition" (Chapter 17) are two of the most common
words in any serious engineering discussion, and they only make sense once
you know exactly what threads do and don't share.

### How it actually works

```
   A PROCESS can contain ONE OR MORE threads. All threads in the same
   process share:
     - the same memory (global variables, heap allocations)
     - the same open file descriptors
     - the same process ID (from the OS's perspective, mostly)

   Each thread gets its OWN:
     - stack (local variables, function call history)
     - CPU register state / program counter (where THIS thread currently
       is in the code)
     - thread ID

   WHY THIS MATTERS SO MUCH: because threads share memory, TWO THREADS
   CAN TOUCH THE SAME VARIABLE AT THE SAME TIME -- something two
   separate PROCESSES simply cannot do by accident, because their
   memory is isolated (Chapter 10). This is BOTH the main reason threads
   are useful (fast communication, no copying needed, no fork() cost)
   AND the main reason threaded code is notoriously hard to get right
   (Chapter 17's race conditions, Chapter 18's locks).

   THREADS VS. PROCESSES, the trade-off in one table:
     PROCESSES   isolated (safer), heavier to create, communicate via
                slower IPC mechanisms (Chapter 19)
     THREADS     share memory (faster, riskier), cheap to create, need
                explicit synchronisation (Chapter 18) to be safe
```

### Real-world example

```
   Node.js was DELIBERATELY built single-threaded for its JavaScript
   execution (I/O is handled separately, under the hood, via a thread
   pool and an event loop) -- a choice its creator, Ryan Dahl, explained
   publicly in his original 2009 talks introducing the project. The
   explicit reasoning: multi-threaded programming with shared mutable
   state is HARD TO GET RIGHT (exactly this chapter's warning), and for
   I/O-heavy server workloads, an event-driven single-threaded model
   sidesteps an entire category of race-condition and locking bugs
   almost entirely, at the cost of needing a different mental model
   (callbacks/promises/async-await) for writing non-blocking code.
   This is a genuinely influential, still-debated design choice
   (Debated: many engineers argue Node's model trades one kind of
   complexity for another, "callback hell" or now `async`/`await`
   chains) -- but it's a real, deliberate illustration of "threads are
   powerful but dangerous enough that entire platforms are designed
   specifically to minimise their use."
```

### Try it yourself

```bash
# See threads vs processes on your own machine right now
ps -eLf | wc -l     # count of THREADS across the whole system
ps aux | wc -l       # count of PROCESSES

# Pick one multi-threaded process and see its thread count
pgrep -f firefox 2>/dev/null | head -1 | xargs -I{} cat /proc/{}/status 2>/dev/null | grep Threads
# or, if you don't have firefox running, try:
cat /proc/self/status | grep Threads   # your shell has how many?
```

### Build it in Go
Goroutines are Go's answer to "threads are expensive". In a measured run,
20,000 goroutines (sleeping or blocked on sockets) used **16 OS threads**,
while 100 goroutines blocked in raw `read(2)` calls pushed it to **108**. See
[Chapter 73](#chapter-73-go-meets-the-kernel-system-calls-threads-and-the-runtime) for the lab, and the
[Go guide's scheduler chapter](../Golang/real-life-golang-guide.md#28-the-scheduler-deep-dive-gmp-and-work-stealing) for how the runtime maps
goroutines onto threads.

### Common mistakes

- **"Threads are always faster than processes."** For CPU-bound, isolated
  work, processes (which can run on separate cores, same as threads can) are
  often simpler and just as fast — threads win specifically when you need
  fast, frequent shared-memory communication.
- **Assuming shared memory means synchronised access.** The OS guarantees
  threads CAN see the same memory; it guarantees NOTHING about the order or
  safety of concurrent access — that's entirely on you and your locks
  (Chapter 18).
- **Thinking "single-threaded" means "can't handle concurrency."** Node.js
  handles huge numbers of CONCURRENT connections single-threadedly, via the
  event loop — concurrency and threading are related but distinct ideas
  (§0.4).

### Check yourself

1. What do threads in the same process share, and what does each thread keep
   private?
2. Why are threads cheaper to create than processes?
3. Why did Node.js's creator choose a single-threaded model for JavaScript
   execution?

### Further reading

- **Book (free):** *OSTEP*, "Concurrency" chapters (the introduction through
  locks) — the best treatment of why threading is hard, and how to do it
  safely.
- **Talk:** Ryan Dahl's original Node.js introduction talk (2009, JSConf) —
  the actual reasoning, from the source.
- **Reference:** `man 7 pthreads` — the POSIX threads API reference.

---

## Chapter 8 — Scheduling: who gets the CPU, and for how long

### In one sentence

With more runnable processes/threads than CPU cores (almost always true), the
kernel's **scheduler** decides, many times per second, which one runs next
and for how long — and the fairness and responsiveness of your entire machine
depends on how good that decision is.

### Why this matters

"Why does my high-priority job feel slow" and "why did this container get
throttled" both trace back to scheduling decisions — and `nice`/cgroup CPU
limits (Chapter 43) are literally you giving the scheduler instructions.

### How it actually works

```
   THE BASIC PROBLEM: you might have 4 CPU cores and 200 runnable
   processes. The scheduler has to pick which handful run RIGHT NOW, and
   for how long, before switching to the next batch -- fairly enough that
   nothing starves, responsively enough that interactive programs (your
   mouse cursor, your SSH session) don't feel laggy, and efficiently
   enough that the switching itself (Chapter 9) doesn't eat all your
   CPU.

   LINUX'S SCHEDULER (CFS -- the Completely Fair Scheduler, the default
   since 2007, itself succeeded by EEVDF in newer kernels) works, in
   spirit, like this:
     - every runnable task accumulates "virtual runtime" while it runs
     - the scheduler always picks the task with the LEAST virtual
       runtime to run next -- i.e., whoever has had the LEAST CPU time
       so far, relatively speaking, goes next.
     - this naturally converges toward fairness: no task can hog the CPU
       indefinitely, because the more it runs, the further behind
       (in scheduling priority) it falls relative to everyone else.

   NICE VALUES let you bias this, without breaking fairness entirely:
     `nice -n 10 mycommand`    lower priority (be nice to others,
                              range -20 to +19, higher number = lower
                              priority)
     `nice -n -10 mycommand`   higher priority (needs privilege)
   This doesn't GUARANTEE a task runs sooner -- it changes the WEIGHT the
   scheduler uses when deciding, so a "niced-up" process accumulates
   virtual runtime slower, and gets picked more often, relatively.
```

### Real-world example

```
   The move FROM Linux's older "O(1) scheduler" TO the Completely Fair
   Scheduler (CFS) in 2007 (kernel 2.6.23, engineered by Ingo Molnar) is
   a well-documented real case study in why scheduling fairness matters
   at scale. The O(1) scheduler's heuristics for distinguishing
   "interactive" from "CPU-hungry" tasks were complex, hard to tune, and
   produced noticeably unfair or laggy behaviour under certain real
   workloads -- users and kernel developers reported cases where
   interactive tasks (things you'd expect to feel snappy, like typing in
   a terminal) got starved by batch workloads. CFS replaced the
   heuristic-heavy approach with the much simpler "always run whoever
   has had the least CPU time so far" fairness principle described
   above, and became the long-standing Linux default specifically
   because it behaved more predictably and fairly across a huge range of
   real workloads without needing per-workload tuning.
```

### Try it yourself

```bash
# Watch nice values change real scheduling behaviour
nice -n 19 yes > /dev/null &     # low priority CPU hog
yes > /dev/null &                # normal priority CPU hog
top    # compare %CPU for the two `yes` processes; the niced one should
       # get noticeably less CPU time on a busy system. Ctrl-C both after.
kill %1 %2

# See the current nice value of every process
ps -eo pid,ni,cmd | head -20
```

### Build it in Go
A container's CPU limit is a CFS quota, not fewer CPUs, so programs that size
thread pools by CPU count get throttled. Go 1.25+ reads the cgroup limit and
sets GOMAXPROCS to match. [Chapter 75](#chapter-75-cpu-limits-gomaxprocs-cgroups-and-throttling) measures the
difference: forcing GOMAXPROCS=8 under a 1.5-CPU limit cost **13.7 s of
throttling in 2 s**.

### Common mistakes

- **Thinking `nice` guarantees priority.** It's a WEIGHT/bias for the
  scheduler's fairness algorithm, not a hard priority override — a heavily
  niced-down process still gets SOME CPU time, just less.
- **Assuming more CPU cores means the scheduler matters less.** With enough
  runnable tasks (very common on a busy server or a laptop with many
  background apps), scheduling decisions matter at any core count.
- **Confusing `nice` (CPU scheduling priority) with `ionice`(I/O scheduling
  priority) or cgroup CPU limits (Chapter 43, a hard cap, not just a bias).**

### Check yourself

1. In one sentence, what problem is the scheduler solving?
2. Roughly how does Linux's Completely Fair Scheduler decide what runs next?
3. Does a negative `nice` value guarantee a process runs immediately? Why or
   why not?

### Further reading

- **Book (free):** *OSTEP*, the scheduling chapters (MLFQ and CFS) — the
  clearest conceptual walkthrough available.
- **Documentation:** the Linux kernel's own `Documentation/scheduler/`
  directory (`sched-design-CFS.rst`) — written by CFS's own author.
- **Reference:** `man 1 nice`, `man 2 sched_setscheduler`.

---

## Chapter 9 — Interrupts, context switches, and why they're not free

### In one sentence

Every time the CPU stops running one thing to run another — because a
hardware device needs attention, or the scheduler decided it's someone else's
turn — it has to save and restore a full set of state, and that switch costs
real, measurable time that adds up at scale.

### Why this matters

This is the concrete reason "fewer, bigger operations beat many, small ones"
shows up everywhere in performance engineering: from batching database writes
to why NGINX's event loop beats thread-per-connection at high concurrency.

### How it actually works

```
   AN INTERRUPT is a signal from HARDWARE (a network card just received
   a packet, a disk finished a read, a timer fired) telling the CPU
   "stop what you're doing and handle this NOW." The CPU pauses whatever
   it was running, jumps to a small piece of kernel code (an interrupt
   handler) to deal with it, then resumes.

   A CONTEXT SWITCH is the broader mechanism: the kernel saves the FULL
   state of whatever was running (every CPU register, the program
   counter, memory-mapping information) into that process/thread's
   saved state, and loads a DIFFERENT one's saved state instead so IT
   can run. This happens:
     - when the scheduler (Chapter 8) decides it's someone else's turn
     - when a process makes a blocking syscall (Chapter 2) and has to
       wait
     - when a hardware interrupt needs kernel attention

   WHY IT'S NOT FREE:
     - saving/restoring register state takes real CPU cycles
     - it very likely INVALIDATES THE CPU CACHE (0.4's memory ladder) --
       the new process's data isn't in the fast, nearby cache, so its
       first few memory accesses are slow, cache-miss reads from RAM
     - on modern CPUs, it can also flush speculative-execution and
       branch-prediction state, adding further real overhead

   A rough, often-cited order of magnitude: a context switch costs low
   MICROSECONDS -- a few thousand CPU cycles' worth of pure overhead,
   before the newly-scheduled task does a single second of useful work.
   That sounds tiny, until you're doing it hundreds of thousands of
   times per second under heavy load.
```

### Real-world example

```
   This is precisely the problem NGINX was engineered to sidestep
   (Chapter 4's C10K story, revisited from this chapter's specific
   angle): a thread- or process-per-connection server handling 10,000
   simultaneous slow/idle connections pays a CONTEXT-SWITCH COST every
   time the scheduler cycles between them, even though almost all of
   them are just sitting there waiting for I/O with nothing to do.
   NGINX's event-driven design uses a small, fixed pool of worker
   processes, each handling MANY connections via non-blocking I/O and a
   single event loop -- collapsing what would be thousands of
   context switches per second down to a small, predictable number,
   which is a major, well-documented reason event-driven servers
   dramatically outperform thread-per-connection designs specifically
   at HIGH CONCURRENCY (not necessarily at low concurrency, where the
   simpler thread-per-connection model can be perfectly fine and easier
   to reason about).
```

### Try it yourself

```bash
# See the total number of context switches your system has done, and
# the rate right now
vmstat 1 5     # the "cs" column = context switches per second

# Compare context-switch rate: idle system vs many competing processes
vmstat 1 3                          # baseline
for i in $(seq 1 20); do yes > /dev/null & done   # spawn 20 CPU hogs
vmstat 1 3                          # watch "cs" jump
kill $(jobs -p)                     # clean up
```

### Common mistakes

- **Thinking interrupts and context switches are the same thing.** An
  interrupt is one specific TRIGGER for kernel attention; a context switch is
  the broader mechanism of swapping what's running, which can happen for
  several reasons.
- **Assuming "more threads/processes" scales linearly.** Beyond some point,
  the context-switch overhead of managing MORE runnable things starts eating
  into the actual work getting done — this is a real, measurable ceiling, not
  a theoretical one.
- **Ignoring cache effects.** The "hidden" cost of a context switch — the
  cold CPU cache afterward — is often bigger than the raw save/restore cost
  itself, and doesn't show up as directly in simple benchmarks.

### Check yourself

1. What is an interrupt, and what is a context switch — how do they relate?
2. Name two reasons a context switch is more expensive than its raw
   save/restore time suggests.
3. Why does an event-driven server (like NGINX) generate far fewer context
   switches than a thread-per-connection server at high concurrency?

### Further reading

- **Book (free):** *OSTEP*, "Limited Direct Execution" chapter — explains
  the trap/context-switch mechanism precisely.
- **Article:** "Context Switching, or Why Servers Aren't Cars" and similar
  engineering-blog explanations of context-switch overhead in production
  systems — search for recent, concrete-numbers writeups.
- **Reference:** `man 8 vmstat` — the tool used above, and what each column
  means.

---

# Part 3 — Memory management

Part 2 gave every process the illusion of its own CPU; this part gives it
the illusion of its own memory. Virtual memory and paging are how that
illusion is built, and this part also covers what happens once the illusion
runs out — swapping, and the OOM killer deciding who dies.

## Chapter 10 — Virtual memory: every process's private lie

### In one sentence

Every process gets its own private, seemingly-huge range of memory addresses
that has almost nothing to do with the actual physical RAM chips — the kernel
and the CPU together translate one into the other, invisibly, on every single
memory access.

### Why this matters

This is why one process can't (accidentally or maliciously) read another
process's memory, why your program can use memory addresses that don't
"physically" exist yet, and why "how much RAM does this process really use"
turns out to be a surprisingly subtle question (Chapter 31's `RSS` vs `VSZ`).

### How it actually works

```
   WITHOUT VIRTUAL MEMORY, every program would need to know the actual,
   physical RAM address of everything it touches -- and if two programs
   both tried to use the same physical address, disaster. This was a
   real problem in early computing.

   WITH VIRTUAL MEMORY, every process gets its OWN, PRIVATE address
   space -- typically the full range a 64-bit pointer can express. When
   your program reads or writes address `0x7fff1234`, the CPU's MEMORY
   MANAGEMENT UNIT (MMU) -- dedicated hardware, working with page tables
   the kernel maintains -- translates that VIRTUAL address into whatever
   PHYSICAL RAM address it actually corresponds to, if any, completely
   transparently to your program:

     Process A's virtual 0x1000  ->  physical RAM address 0x7A2000
     Process B's virtual 0x1000  ->  physical RAM address 0x3F1000
                                     (a DIFFERENT physical location --
                                      this is the whole isolation
                                      guarantee)

   THIS ALSO ENABLES SOME GENUINELY USEFUL TRICKS:
     - two processes can SHARE the same physical memory (a shared
       library like libc) while each believes it has its own private
       copy at its own virtual address.
     - a virtual address doesn't need to be backed by RAM AT ALL, right
       now -- it can be backed by disk (Chapter 11's paging), or not
       backed by anything until it's actually touched (this is how
       `malloc()`-ing a huge amount of memory can succeed instantly:
       nothing physical is allocated until you actually write to it).
```

### Real-world example

```
   Address Space Layout Randomization (ASLR) is a direct SECURITY
   application of virtual memory, standard in every modern OS since the
   mid-2000s (mainline Linux since kernel 2.6.12, 2005). Because a
   process's virtual addresses are just a mapping the kernel controls,
   the kernel can RANDOMISE where key regions (the stack, the heap,
   shared libraries) land in a process's virtual address space EACH TIME
   it runs. This specifically defeats a huge class of MEMORY-CORRUPTION
   EXPLOITS that rely on knowing exactly where, in memory, a piece of
   attacker-useful code or data will be -- without virtual memory's
   indirection, there would be nothing TO randomise; physical addresses
   would need to change, which is far harder to do safely. This is a
   genuinely important, still-active area of security engineering (ASLR
   bypasses and defeats are an ongoing arms race, referenced throughout
   the companion security guides in this wiki) built entirely on top of
   the OS concept this chapter teaches.
```

### Try it yourself

```bash
# See your own shell's virtual memory map -- notice how large and sparse
# the address ranges are compared to your machine's actual physical RAM
cat /proc/self/maps | head -20

# Compare virtual size (VSZ) to actual resident memory (RSS) for a
# running process -- VSZ is often dramatically larger
ps -eo pid,vsz,rss,cmd | sort -k2 -n -r | head -5

# See ASLR in action: run the same program twice, watch a library's
# load address change between runs
for i in 1 2; do
  grep libc /proc/self/maps 2>/dev/null | head -1
done
# (run in two separate shell invocations for a clean comparison, e.g.:
#  bash -c 'grep libc /proc/self/maps | head -1'  -- run this twice)
```

### Common mistakes

- **"Virtual memory means slower memory."** The translation happens in
  dedicated hardware (the MMU) and is extremely fast — the SLOW part only
  happens when a virtual address isn't currently backed by RAM at all
  (Chapter 11's page fault).
- **Assuming a process's memory usage equals its virtual address space
  size.** `VSZ` (virtual size) can be enormous and mostly meaningless; `RSS`
  (resident set size — memory ACTUALLY in RAM right now) is usually the
  number you care about.
- **Thinking two processes literally cannot share memory.** They can, on
  purpose (shared libraries, `mmap`-based shared memory, Chapter 19) — what
  virtual memory prevents is ACCIDENTAL, unintended sharing.

### Check yourself

1. What problem does virtual memory solve, in one sentence?
2. What translates a virtual address into a physical one, and where does
   that translation information live?
3. Why does ASLR depend on virtual memory existing at all?

### Further reading

- **Book (free):** *OSTEP*, the entire "Virtualization: Memory" section —
  the clearest available walkthrough of address translation, from first
  principles.
- **Reference:** `man 5 proc` (the `/proc/[pid]/maps` section).
- **Article:** "A Brief History of ASLR" or similar retrospectives on why
  and when major OSes adopted it.

---

## Chapter 11 — Paging, page faults, and the TLB

### In one sentence

Physical memory is managed in fixed-size chunks called **pages**, and when a
process touches a virtual address that isn't currently mapped to a physical
page, the CPU raises a **page fault** — which the kernel may resolve
instantly, slowly (by reading from disk), or fatally (by killing the process).

### Why this matters

"Why did my program's first request feel slow, but the second one was fast"
and "why does swapping make everything grind to a halt" are both, directly,
page-fault stories.

### How it actually works

```
   PAGES: physical RAM and virtual address space are both divided into
   fixed-size chunks -- PAGES -- commonly 4KB on Linux/x86 (larger "huge
   pages," 2MB or 1GB, exist for specific performance-sensitive
   workloads). The kernel's PAGE TABLE records, for every virtual page a
   process uses, WHERE (or whether) it's currently backed by a physical
   page.

   A PAGE FAULT happens when a process accesses a virtual address whose
   page ISN'T currently mapped to physical RAM. There are several very
   different flavours:

     MINOR FAULT (cheap, common, totally normal)
       the page exists somewhere accessible but isn't yet mapped into
       THIS process (e.g. a shared library page another process already
       loaded) -- the kernel just updates the mapping. Fast, invisible.

     MAJOR FAULT (expensive)
       the data has to be read FROM DISK -- either it's part of a
       memory-mapped file not yet loaded, or (the painful case) it was
       previously SWAPPED OUT to disk to free up RAM for something else,
       and now has to be swapped back IN. This is orders of magnitude
       slower than RAM access (Chapter 12's hierarchy).

     SEGMENTATION FAULT (fatal)
       the process touched a virtual address that ISN'T VALID AT ALL for
       it -- not "not yet loaded," but genuinely not part of its address
       space (a wild pointer, a buffer overrun). The kernel delivers
       SIGSEGV (Chapter 20) and the process usually dies. This is the
       infamous "segfault."

   THE TLB (Translation Lookaside Buffer) is a small, very fast CACHE,
   right on the CPU, of recent virtual-to-physical translations -- so the
   full page-table lookup (itself potentially several memory accesses)
   doesn't have to happen on EVERY single memory access. A "TLB miss"
   means falling back to the slower, full page-table walk.
```

### Real-world example

```
   THRASHING is the well-documented, classic failure mode this chapter's
   mechanism produces under memory pressure: when a system's actively-
   used memory (its "working set") exceeds available RAM, the kernel
   starts SWAPPING pages out to disk to make room, and swapping others
   BACK IN as they're needed -- generating a storm of expensive MAJOR
   PAGE FAULTS. If the working set is large enough relative to RAM, the
   system can spend almost ALL of its time swapping pages in and out
   and almost NONE of its time doing actual work -- CPU usage looks
   deceptively LOW (everything is stuck waiting on disk I/O, D-state,
   Chapter 5) while the system feels completely unresponsive. This
   remains one of the most common real production incidents reported in
   SRE/ops communities: "the server isn't doing much CPU-wise but is
   totally unresponsive" is, very often, a thrashing system -- and it's
   exactly why cloud providers and Kubernetes (Chapter 58) generally
   prefer to KILL a process/pod that's exceeding its memory budget
   (Chapter 13's OOM killer) rather than let it swap into oblivion.
```

### Try it yourself

```bash
# See page fault counts for a command -- minor vs major
/usr/bin/time -v ls / 2>&1 | grep -i "page faults"
# (on some systems: /usr/bin/time -v cat /proc/self/maps 2>&1 | grep fault)

# Watch major faults spike under real memory pressure (careful, only on
# your disposable VM -- this deliberately stresses memory)
free -h                          # baseline
vmstat 1 5                       # watch the "si"/"so" (swap in/out) columns

# See your system's page size
getconf PAGE_SIZE                # almost always 4096 on x86_64 Linux
```

### Common mistakes

- **Treating all page faults as bad.** Minor faults are constant, cheap, and
  completely normal — only SUSTAINED major faults (especially swap
  activity) are a red flag.
- **Confusing a segmentation fault with "the program is out of memory."**
  A segfault means an INVALID memory access (a bug), not memory exhaustion —
  those are different failures with different causes (Chapter 13 covers
  actual memory exhaustion).
- **Assuming more swap space fixes a thrashing system.** It usually just
  makes the thrashing last longer before the system becomes truly unusable —
  the real fix is more RAM or a smaller working set.

### Check yourself

1. What's the difference between a minor page fault and a major page fault?
2. What does the TLB cache, and why does it exist?
3. What is thrashing, and why does a thrashing system often show LOW CPU
   usage despite being unresponsive?

### Further reading

- **Book (free):** *OSTEP*, "Paging" and "Swapping" chapters.
- **Reference:** `man 2 mmap`, `man 8 vmstat`.
- **Article:** search "Linux memory thrashing troubleshooting" for current,
  concrete production writeups.

---

## Chapter 12 — The memory hierarchy: registers to disk, and why it exists

### In one sentence

Every layer of storage in a computer — CPU registers, cache, RAM, SSD, spinning
disk, the network — trades capacity for speed in the opposite direction, and
almost all of performance engineering is about keeping the data you need
right now as high up that ladder as possible.

### Why this matters

This single idea explains why caching exists at every layer of every system
you'll ever build, why "just add more RAM" fixes so many problems, and why
network calls are the slowest thing your code routinely does.

### How it actually works

```
   THE LADDER, fastest/smallest at the top, slowest/largest at the
   bottom (this exact table, in this exact spirit, is widely known in
   the industry as "latency numbers every programmer should know,"
   popularised by Google engineer Jeff Dean's internal talks and now
   cited constantly across performance-engineering material):

     CPU REGISTER          ~ sub-nanosecond          (a handful of bytes)
     L1 CPU CACHE          ~ 1 nanosecond             (tens of KB)
     L2 CPU CACHE          ~ few nanoseconds          (hundreds of KB)
     L3 CPU CACHE          ~ tens of nanoseconds      (tens of MB)
     MAIN MEMORY (RAM)     ~ 100 nanoseconds          (GBs)
     SSD (random read)     ~ 100 MICROSECONDS         (hundreds of GB-TB)
     SPINNING DISK (seek)  ~ 1-10 MILLISECONDS        (many TB)
     NETWORK (same DC)     ~ 0.5 MILLISECONDS round trip
     NETWORK (cross-region/
       continent)          ~ 50-150+ MILLISECONDS round trip

   NOTICE THE GAPS: RAM is roughly 100x slower than L1 cache. An SSD is
   roughly 1,000x slower than RAM. A cross-region network call can be
   MILLIONS of times slower than a register access. These aren't small
   differences -- they're differences of many ORDERS OF MAGNITUDE, and
   they compound: a program that "just" does one unnecessary disk read
   per request, at scale, can dominate your entire system's latency
   budget.

   THIS IS WHY CACHING EXISTS EVERYWHERE: the CPU caches recently-used
   RAM contents (hardware-managed, invisible to you). The OS caches
   recently-read disk blocks in spare RAM (Chapter 16's page cache).
   Your application caches expensive database queries in Redis. Your
   CDN caches HTTP responses close to users. It's the SAME PATTERN,
   repeated at every layer of the stack, because the SAME physics
   (this ladder) applies at every layer.
```

### Real-world example

```
   This exact table -- Jeff Dean's "Numbers Every Programmer Should
   Know" -- has circulated across the software industry for over a
   decade (originally from internal Google engineering talks, widely
   re-published and interactively visualised, e.g. at
   `colin-scott.github.io/personal_website/research/interactive_latency.html`)
   specifically because engineers KEPT making the same class of mistake:
   architecting a system as if a network call, a disk read, and a memory
   access all cost roughly the same. A single unnecessary synchronous
   network call inside a hot loop -- something that looks completely
   harmless in code review -- can be the difference between a service
   responding in single-digit milliseconds and one that takes whole
   seconds, purely because of where "network round trip" sits on this
   ladder relative to everything above it. Real, well-known performance
   incidents at scale are very often traced back to exactly this: doing
   an operation from a slow rung of the ladder when a faster one would
   have worked, repeated enough times (in a loop, per-request, per-item)
   for the multiplier to dominate everything else.
```

### Try it yourself

```bash
# Feel the RAM-vs-disk gap directly: write and read the same amount of
# data via RAM (tmpfs, usually mounted at /dev/shm) vs actual disk
dd if=/dev/zero of=/dev/shm/testfile bs=1M count=512 2>&1 | tail -1
dd if=/dev/zero of=./testfile-ondisk bs=1M count=512 2>&1 | tail -1
rm /dev/shm/testfile ./testfile-ondisk

# See your CPU's actual cache sizes
lscpu | grep -i cache
```

### Common mistakes

- **"RAM is basically infinite/free compared to disk, so it doesn't
  matter."** The GAPS at every level of this ladder matter, not just the
  RAM/disk one — cache-unfriendly code (poor memory access patterns) can be
  dramatically slower even when everything technically "fits in RAM."
- **Treating all "fast" operations as equally fast.** A local network call
  and a cross-region network call can differ by 100x or more — "it's just a
  network call" hides an enormous range.
- **Over-caching without eviction strategy.** A cache with no size limit or
  eviction policy just becomes a slow memory leak wearing a fast hat.

### Check yourself

1. List the memory hierarchy from fastest/smallest to slowest/largest.
2. Roughly how many orders of magnitude slower is a disk read than a RAM
   access?
3. Why does the SAME caching pattern appear at the CPU, OS, application, and
   CDN layers?

### Further reading

- **Reference:** the interactive "Latency Numbers Every Programmer Should
  Know" visualisation (`colin-scott.github.io/personal_website/research/
  interactive_latency.html`) — genuinely worth spending five minutes with.
- **Book (free):** *OSTEP*, the caching-related sections of the memory
  chapters.
- **Talk:** search for Jeff Dean's "Designs, Lessons and Advice from Building
  Large Distributed Systems" — the original context for these numbers.

---

## Chapter 13 — Memory allocation and the OOM killer

### In one sentence

When Linux runs critically low on memory and can't free any more, it doesn't
just let the system grind to a halt — it picks a process to kill, using a
scoring heuristic, to save the rest of the machine.

### Why this matters

"Why did my process just disappear with no error message" and "why did
Kubernetes restart my pod with status `OOMKilled`" are the two most common
ways engineers meet this mechanism for the first time, usually in production,
usually stressfully.

### How it actually works

```
   MEMORY ALLOCATION, briefly: when your program calls `malloc()` (or
   its language's equivalent -- `new`, a Python object, a Go allocation),
   it's asking a USER-SPACE allocator library for memory, which in turn
   asks the KERNEL for more address space via syscalls like `brk()` or
   `mmap()` (Chapter 2) when it needs more room. Critically: on Linux,
   by default, this request can SUCCEED even if the system doesn't
   actually have that much free memory right now -- this is called
   OVERCOMMIT. The kernel is betting that not every process will
   actually USE all the memory it asked for, all at once.

   THAT BET SOMETIMES FAILS. If enough processes actually DO try to use
   the memory they were promised, and the system genuinely runs out of
   both RAM and swap, the kernel has no way to honour the requests
   anymore. Rather than let the whole system deadlock or crash, the
   OOM KILLER (Out-Of-Memory killer) steps in and FORCIBLY KILLS A
   PROCESS to free memory.

   WHICH PROCESS GETS KILLED is decided by a "badness" score
   (`oom_score`, visible per-process in `/proc/<pid>/oom_score`), roughly
   weighted toward: how much memory a process is using (bigger =
   more likely to be picked), how long it's been running, and whether
   it's been given special treatment via `oom_score_adj` (a value from
   -1000, "never kill this," to +1000, "kill this first"). The kernel's
   goal is "free the most memory for the least collateral damage" -- but
   its notion of "damage" has no idea your database is more important
   than a stray debugging script, unless you've told it so explicitly.
```

### Real-world example

```
   Kubernetes' `OOMKilled` pod status (a status any engineer running
   containers will eventually see, documented extensively in Kubernetes'
   own docs and in countless real-world troubleshooting writeups) is
   THIS EXACT MECHANISM, one layer up: a container's memory usage is
   enforced via a cgroup memory limit (Chapter 43, Phase 3's Chapter 44),
   and when a container tries to exceed it, the KERNEL'S OOM KILLER --
   scoped to that cgroup -- kills the offending process inside the
   container, and Kubernetes reports it as `OOMKilled`.
   A well-known, recurring real production pain point: the OOM killer,
   left to its DEFAULT scoring, can kill the WRONG process on a
   multi-process host or inside a multi-process container -- for
   example, killing a critical, well-behaved database process because
   it happens to have a large (but legitimate, intentional) memory
   footprint, while leaving a genuinely leaking, smaller process alive.
   This is precisely why production runbooks and platform teams commonly
   set `oom_score_adj` explicitly for critical processes (or, in
   Kubernetes, set requests/limits and Quality-of-Service classes
   deliberately, Phase 3), rather than trusting the kernel's default
   heuristic to guess correctly under pressure.
```

### Try it yourself

```bash
# See the current OOM score for your own shell (lower = less likely to
# be killed first)
cat /proc/self/oom_score
cat /proc/self/oom_score_adj

# Watch memory overcommit in action (safe, small-scale demo)
cat /proc/sys/vm/overcommit_memory   # 0 = heuristic (default),
                                     # 1 = always allow, 2 = strict

# See what the OOM killer has done recently on your system, if anything
dmesg | grep -i "out of memory" || echo "no OOM events found (good!)"
journalctl -k | grep -i "killed process" 2>/dev/null | tail -5
```

```
   Dangerous: don't deliberately trigger the OOM killer on a shared or
   important machine. If you want to see it fire, do it on a throwaway
   VM with a small, fixed memory size, and expect to need to reboot it.
```

### Build it in Go
Garbage-collected runtimes and memory limits interact badly by default: Go
lets the heap grow to twice the live data before collecting. In
[Chapter 76](#chapter-76-memory-limits-the-go-heap-gomemlimit-and-the-oom-killer), a Go program with 160 MB of live data is
OOM-killed in a 300 MB container (exit 137) and survives with
`GOMEMLIMIT=250MiB`.

### Common mistakes

- **"My process crashed, it must be a bug in my code."** Check `dmesg`/
  `journalctl` for OOM kills FIRST — a silently-vanished process with no
  application-level error message is a classic OOM-kill signature.
- **Assuming `malloc()` succeeding means the memory is really available.**
  Overcommit means a successful allocation is a PROMISE, not a guarantee —
  the bill can come due later, unpredictably, possibly for a DIFFERENT
  process entirely.
- **Never setting `oom_score_adj` or container memory limits for critical
  processes.** Leaves you at the mercy of the kernel's generic heuristic
  during exactly the moment you can least afford a wrong guess.

### Check yourself

1. What is memory overcommit, and why does the kernel do it?
2. What does the OOM killer do, and roughly how does it choose a victim?
3. Why might a critical database process get killed by the OOM killer
   instead of a smaller, actually-leaking process?
4. How does this mechanism show up in Kubernetes?

### Further reading

- **Documentation:** the Linux kernel's own `Documentation/admin-guide/
  mm/overcommit-accounting.rst` and OOM-killer documentation.
- **Reference:** `man 5 proc` (the `oom_score`/`oom_score_adj` sections).
- **Docs:** Kubernetes' own documentation on "Resource Management for Pods
  and Containers" — the direct, real-world continuation of this chapter.

---

# Part 4 — I/O and storage

CPU and memory are the first two resources the OS manages on a process's
behalf; I/O is the third, and the slowest by orders of magnitude. This part
follows a single `read()` call all the way down through the kernel's I/O
stack to the physical disk, and back.

## Chapter 14 — The I/O stack: from `read()` to the disk platter

### In one sentence

A single `read()` or `write()` call passes through many distinct layers — the
C library, the kernel's virtual filesystem layer, a specific filesystem
driver, the block layer, and finally a device driver — and each layer adds
both useful abstraction and a small amount of real latency.

### Why this matters

"I wrote to the file, why is the data still not really saved" is one of the
most consequential misunderstandings in software engineering — it has caused
real, well-documented data loss.

### How it actually works

```
   THE STACK, roughly, for a file write:

     your program:      write(fd, buf, len)
                              |  (a syscall, Chapter 2)
     kernel VFS layer:  the Virtual File System -- a generic interface
                        so "write to a file" works the same whether it's
                        ext4, XFS, NFS, or a network filesystem
                              |
     filesystem driver: ext4 (or whichever) decides WHERE this data goes
                        on the underlying block device, updates its own
                        metadata (Chapter 15)
                              |
     page cache:         (Chapter 16) the write usually lands HERE first
                        -- in RAM -- and is marked "dirty," not yet
                        flushed to the physical device
                              |
     block layer:        eventually, the kernel decides to actually flush
                        dirty pages to disk; the block layer batches and
                        schedules these I/O requests (Chapter 16)
                              |
     device driver:      talks the actual protocol (SATA, NVMe, etc.) to
                        the physical device
                              |
     THE PHYSICAL DEVICE: an SSD's flash cells, or a spinning disk's
                        magnetic platter, finally, actually, stores the
                        bytes.

   THE CRITICAL, OFTEN-MISUNDERSTOOD POINT: your `write()` call RETURNING
   SUCCESSFULLY usually only means the data reached the PAGE CACHE (step
   4) -- it does NOT mean the data is safely on the physical device yet.
   If the machine loses power before the kernel flushes that dirty page,
   the write can be LOST, even though your program was told it
   succeeded. `fsync()` is the syscall that explicitly asks the kernel
   to flush a specific file's data all the way down to the device and
   wait for confirmation -- and even THEN, some storage hardware/
   virtualisation layers have historically lied about actually
   completing that flush, for performance reasons.
```

### Real-world example

```
   In 2009, a real, widely-discussed controversy broke out on the Linux
   Kernel Mailing List after users of the (then-new) ext4 filesystem
   started reporting that files they had JUST WRITTEN would sometimes be
   ZERO BYTES or contain GARBAGE after an unexpected power loss or
   crash -- even though their application code had called the standard
   `write()`/`close()` sequence that had "always worked" on the older
   ext3 filesystem. The underlying cause was exactly this chapter's
   layering: ext4's newer "delayed allocation" optimisation could delay
   actually writing data to disk for LONGER than ext3 typically had,
   combined with applications that were, technically, always relying on
   an UNDOCUMENTED, filesystem-specific timing behaviour rather than
   explicitly calling `fsync()` when they actually needed a durability
   guarantee. Ext4's own lead maintainer, Ted Ts'o, and Linus Torvalds
   both weighed in publicly (their mailing-list posts from the time are
   still findable and are a genuinely excellent, blunt explanation of
   this exact layering problem) -- the eventual resolution involved BOTH
   kernel-side changes (shortening the delay) AND, more fundamentally, a
   renewed, widely-repeated industry reminder: if your application needs
   a guarantee that data is durably on disk, you must call `fsync()`
   (or use a library/database that does it correctly for you) --
   "the write() call returned" was never actually that guarantee.
```

### Try it yourself

```bash
# See the layers in action: strace a simple file write and watch the
# actual syscalls (write, then, if the program calls it, fsync)
strace -e trace=write,fsync,close -o /tmp/trace.log \
  bash -c 'echo hello > /tmp/testfile.txt'
cat /tmp/trace.log

# Watch dirty (not-yet-flushed) pages accumulate and get flushed
cat /proc/meminfo | grep -i dirty
dd if=/dev/zero of=/tmp/bigfile bs=1M count=200 2>&1 | tail -1
cat /proc/meminfo | grep -i dirty    # likely non-zero right after
sleep 3
cat /proc/meminfo | grep -i dirty    # should shrink as the kernel flushes
rm /tmp/bigfile /tmp/testfile.txt
```

### Common mistakes

- **Assuming a successful `write()` means the data is safe on disk.** It
  usually only means the data reached the page cache — durability requires
  `fsync()` (or the equivalent your database/library already does for you).
- **Calling `fsync()` after every tiny write "to be safe."** This defeats the
  performance benefit of the page cache almost entirely — the real skill is
  knowing WHEN durability actually matters (a committed transaction) vs. when
  it doesn't (an intermediate scratch file).
- **Blaming the filesystem for data loss that's actually an application
  bug.** As the ext4 story shows, the real fix was usually in APPLICATIONS
  correctly requesting durability, not in the filesystem "misbehaving."

### Check yourself

1. List the layers a `write()` call passes through, from your program to the
   physical device.
2. Does a successful `write()` guarantee the data is safely on disk? What
   does guarantee that?
3. What was the real, underlying lesson from the 2009 ext4 delayed-allocation
   controversy?

### Further reading

- **Book (free):** *OSTEP*, the "File Systems" introduction and "Crash
  Consistency" chapters.
- **Reference:** `man 2 fsync`, `man 2 write`.
- **History:** search "ext4 delayed allocation data loss 2009 LKML" — the
  original mailing-list threads are genuinely worth reading firsthand.

---

## Chapter 15 — Filesystems: inodes, directories, and journaling

### In one sentence

A filesystem separates a file's actual **data** from its **metadata** (an
"inode" holding size, permissions, timestamps, and pointers to the data) from
its **name** (which just lives in a directory, pointing at an inode) — and
understanding this split explains several classic "wait, that shouldn't be
possible" filesystem behaviours.

### Why this matters

"`df` says I have space but I get 'no space left on device'," "why does
deleting a huge open log file not free disk space until the process
restarts," and "why can I rename a file across directories instantly" are all
answered by this one model.

### How it actually works

```
   THREE SEPARATE THINGS, commonly conflated as "a file":

     THE DATA           the actual bytes, stored in disk blocks
                        somewhere.
     THE INODE          a fixed-size metadata record: size, permissions,
                        owner, timestamps, and (crucially) POINTERS to
                        which disk blocks hold the data. Identified by an
                        INODE NUMBER (see it with `ls -i`).
     THE DIRECTORY
     ENTRY               just a NAME mapped to an inode number, living
                        inside a directory (which is itself just a
                        special kind of file). Multiple directory entries
                        (even in different directories) CAN point to the
                        SAME inode -- this is exactly what a HARD LINK
                        is (Chapter 25).

   THIS EXPLAINS SEVERAL THINGS THAT OTHERWISE SEEM WEIRD:
     - deleting a file (`rm`) doesn't necessarily delete the DATA -- it
       removes a directory entry and DECREMENTS the inode's "link
       count." The data is only actually freed when the link count AND
       the number of processes with the file open both reach zero. This
       is why a huge log file, deleted while a process still has it
       open, doesn't free disk space until that process closes it or
       restarts.
     - renaming/moving a file WITHIN the same filesystem is nearly
       instant, regardless of file size -- it's just updating a
       directory entry to point at the same inode, no data is actually
       moved. (Moving ACROSS filesystems is different -- it has to
       actually copy the data, then delete the original.)
     - a filesystem can run out of INODES before it runs out of DATA
       SPACE -- each filesystem has a FIXED number of inodes, set when
       it was created, independent of how much raw storage capacity
       exists.

   JOURNALING (ext4, XFS, and most modern filesystems use it): before
   making a metadata change, the filesystem first writes a small record
   of its INTENT to a journal (a dedicated log area). If the system
   crashes mid-operation, the filesystem can REPLAY the journal on
   next boot to finish or cleanly roll back the interrupted operation,
   instead of leaving the filesystem's own internal structures in an
   inconsistent, possibly-corrupted state. This is why modern Linux
   boxes almost never need the old, slow, full filesystem-consistency
   check after an unclean shutdown that was routine decades ago.
```

### Real-world example

```
   "No space left on device" while `df -h` clearly shows free space is
   a genuinely common, well-documented real production incident, and
   it's a DIRECT consequence of this chapter's inode/data split: every
   filesystem is created with a FIXED NUMBER of inodes (`df -i` shows
   this, separately from `df -h`'s DATA space), and a workload that
   creates a very large number of SMALL files -- a classic real
   example being a mail server storing each email as a separate tiny
   file (the Maildir format), or a web application writing huge numbers
   of small session/cache files -- can exhaust the inode count while
   still having plenty of raw disk space free. This specific failure
   mode is documented extensively across sysadmin/SRE knowledge bases
   (Red Hat, DigitalOcean, and countless "why can't I write this file
   when df says I have space" Stack Overflow threads) precisely because
   it's a genuinely counter-intuitive symptom until you understand that
   "disk space" and "inode count" are two SEPARATE, independently
   exhaustible resources.
```

### Try it yourself

```bash
# See inode usage vs. data usage -- two separate numbers
df -h .      # data space
df -i .      # inode count

# See a file's inode number, and watch a hard link share it
echo "hello" > original.txt
ls -i original.txt
ln original.txt hardlink.txt
ls -i hardlink.txt          # SAME inode number as original.txt
rm original.txt
cat hardlink.txt            # still works -- the data wasn't deleted,
                            # only one of its two directory entries was
rm hardlink.txt

# Watch "delete a file, space isn't freed until the process closes it"
python3 -c "
import time
f = open('/tmp/bigfile.tmp', 'w')
f.write('x' * 50_000_000)
f.flush()
import subprocess
subprocess.run(['rm', '/tmp/bigfile.tmp'])   # 'deleted' but still open
print('file removed from directory, but check: lsof | grep deleted')
time.sleep(5)
f.close()   # NOW the space is actually freed
"
```

### Build it in Go
Journaling protects the filesystem's own structures, not your file's contents.
[Chapter 77](#chapter-77-files-that-survive-crashes-page-cache-fsync-and-atomic-replacement) implements the write → fsync → rename →
fsync-directory recipe that databases and package managers use to make file
replacement atomic and durable, and measures its cost (about 1 ms per write on
Linux, versus 0.07 ms for a plain write that isn't safe).

### Common mistakes

- **"`df -h` shows free space, so I can't be out of disk space."** Check
  `df -i` too — inode exhaustion is a separate, real failure mode.
- **Assuming `rm`-ing a file always immediately frees disk space.** Not if
  another process still has it open — the space is freed only when the last
  reference (directory entry AND open file handle) goes away.
- **Confusing a hard link with a symbolic link (symlink).** A hard link
  shares the SAME inode (can't cross filesystems, invisible as a "link" at
  all); a symlink is a separate, small file that just contains a PATH
  string, and breaks if the target moves (Chapter 25 covers both properly).

### Check yourself

1. What are the three separate things commonly lumped together as "a file"?
2. Why doesn't deleting an open file immediately free its disk space?
3. Why can a filesystem run out of space for NEW FILES while `df -h` shows
   plenty of free bytes?
4. What does journaling protect against, and why did it reduce the need for
   long filesystem-check-on-boot delays?

### Further reading

- **Book (free):** *OSTEP*, "Files and Directories" and "Crash Consistency:
  FSCK and Journaling" chapters.
- **Reference:** `man 2 stat`, `man 1 df` (the `-i` flag), `man 2 link`.
- **Article:** search "linux inode exhaustion no space left on device" for
  current, practical troubleshooting writeups.

---

## Chapter 16 — Caching, buffering, and disk scheduling

### In one sentence

Linux aggressively uses "spare" RAM to cache recently-read and recently-
written disk data (the **page cache**), which makes your system faster and
your `free -h` output more confusing, in roughly equal measure.

### Why this matters

More production monitoring false alarms are caused by misreading `free -h`
than almost any other single Linux quirk — and understanding the page cache
fixes that instantly and permanently.

### How it actually works

```
   THE PAGE CACHE: whenever the kernel reads data from disk, it keeps a
   copy in RAM (in pages it isn't otherwise using) in case that same
   data is needed again soon -- which, for most real workloads, it often
   is (re-reading the same config file, the same popular database rows,
   the same commonly-loaded library). This is transparent to your
   program: you didn't ask for it, you don't manage it, but it makes
   repeated reads dramatically faster (Chapter 12's hierarchy: RAM vs.
   disk).

   BUFFERING WRITES works similarly, in the other direction (Chapter
   14): writes typically land in the page cache first (as "dirty"
   pages) and get flushed to the actual device by the kernel on its own
   schedule (or when explicitly `fsync()`-ed), batching and reordering
   them for efficiency rather than hitting the disk on every single
   write.

   THE KEY INSIGHT FOR READING `free -h` CORRECTLY: this cached data is
   NOT "wasted" or "unavailable" memory. The kernel will INSTANTLY
   reclaim page-cache pages the moment any process actually needs that
   RAM for something else -- there's no cost or delay to giving it back
   (unless it's DIRTY and needs flushing first). This is why modern
   `free -h` output distinguishes:
     "free"       truly, completely unused memory (often deliberately
                 near zero on a healthy, long-running Linux box -- the
                 kernel would rather use spare RAM for cache than leave
                 it idle)
     "available"  the number that actually matters: free memory PLUS
                 memory that could be reclaimed painlessly from cache if
                 needed -- this is "how much RAM can a new process
                 actually use right now."

   DISK SCHEDULING: the kernel's block layer (Chapter 14) doesn't just
   send I/O requests to the device in the order they arrive -- it can
   reorder and batch them (an I/O SCHEDULER) to reduce costly disk-head
   movement (on spinning disks) or maximise the parallel request queue
   depth (on SSDs/NVMe, which have very different performance
   characteristics and mostly use simpler, low-overhead schedulers like
   `mq-deadline` or `none` rather than the seek-minimising algorithms
   designed for spinning disks).
```

### Real-world example

```
   "The server shows almost no FREE memory, is it about to run out?"
   is one of the single most common, well-documented FALSE ALARMS in
   Linux operations -- widely discussed across sysadmin communities,
   monitoring-tool documentation (many observability platforms
   explicitly warn users about this in their own docs), and countless
   "why is my Linux server always at 95% memory usage" support threads.
   The near-universal answer: it's the page cache doing exactly what
   it's designed to do -- using otherwise-idle RAM to speed up disk
   access -- and it is NOT memory under pressure. A NAIVE monitoring
   alert configured on raw "free memory" percentage will fire constantly
   on a perfectly healthy, well-utilised Linux server, training whoever
   is on-call to ignore the alert entirely (a real, dangerous "alert
   fatigue" pattern) -- which is exactly why correctly-configured
   monitoring alerts on "AVAILABLE" memory (or, better, on actual signs
   of memory PRESSURE like swap activity or OOM kills, Chapter 13)
   instead of raw "free."
```

### Try it yourself

```bash
# See the page cache in action, and the free-vs-available distinction
free -h                       # note "available" is much bigger than "free"

# Watch the page cache grow as you read a large file
sync && echo 3 | sudo tee /proc/sys/vm/drop_caches > /dev/null 2>&1 || true
free -h
dd if=/dev/zero of=/tmp/cachetest bs=1M count=300 2>&1 | tail -1
dd if=/tmp/cachetest of=/dev/null bs=1M 2>&1 | tail -1   # first read
free -h                       # cache usage grew
dd if=/tmp/cachetest of=/dev/null bs=1M 2>&1 | tail -1   # second read --
                                                         # should be MUCH
                                                         # faster (served
                                                         # from cache)
rm /tmp/cachetest

# See your current I/O scheduler
cat /sys/block/*/queue/scheduler 2>/dev/null
```

### Common mistakes

- **Alerting on raw "free" memory percentage.** Almost guaranteed to produce
  false alarms on a healthy system — use "available" memory, or better,
  actual pressure signals (swap, OOM events).
- **Manually dropping caches "to free memory" as a routine practice.** It's
  a legitimate diagnostic tool occasionally, but doing it habitually just
  makes your NEXT disk read slower, for no real benefit — the cache was
  helping you.
- **Assuming a busy-looking `iostat` means your disk is the bottleneck.**
  Cross-check with CPU and memory numbers too — high disk utilisation with
  everything else idle is a very different story than high disk utilisation
  alongside CPU saturation.

### Check yourself

1. What is the page cache, and why does the kernel use "spare" RAM for it
   rather than leaving it empty?
2. What's the difference between "free" and "available" memory in `free -h`
   output, and which one should you actually monitor?
3. Why is spinning-disk I/O scheduling different from SSD/NVMe I/O
   scheduling?

### Further reading

- **Book (free):** *OSTEP*, I/O-related sections of the file-systems
  chapters.
- **Reference:** `man 1 free`, `man 8 iostat`, the kernel's
  `Documentation/admin-guide/sysctl/vm.rst`.
- **Article:** search "linux free available memory explained" — this exact
  confusion has been explained many times; find one with clear diagrams.

---

# Part 5 — Concurrency and inter-process communication

Parts 2–4 covered a single process's view of the machine. This part is what
happens once multiple processes — or multiple threads inside one — touch the
same resource at the same time: race conditions first, then the IPC
mechanisms (pipes, signals, shared memory) that let them cooperate on
purpose instead of corrupting each other by accident.

## Chapter 17 — Race conditions at the OS level

### In one sentence

A race condition happens when the correctness of a program depends on the
precise TIMING or ORDER of events that the OS does not guarantee — and
because that timing can vary run to run, the bug can be invisible in testing
and catastrophic in production.

### Why this matters

Race conditions are one of the hardest classes of bug to find, because they
often don't reproduce reliably — "it worked when I tested it" is exactly what
you'd expect from a race condition, not evidence it's fixed.

### How it actually works

```
   THE OS MAKES FEW TIMING GUARANTEES. When two threads (Chapter 7) or
   processes run "at the same time," the scheduler (Chapter 8) decides
   the EXACT interleaving of their instructions, and that interleaving
   can change from run to run -- based on system load, CPU count,
   scheduling decisions that have nothing to do with your program's
   logic.

   THE CLASSIC SHAPE OF A RACE CONDITION -- a "check then act" sequence
   that ISN'T ATOMIC (isn't guaranteed to happen as one uninterruptible
   step):

     Thread A: reads shared_counter (sees 5)
     Thread B: reads shared_counter (ALSO sees 5, before A writes back)
     Thread A: writes shared_counter = 5 + 1 = 6
     Thread B: writes shared_counter = 5 + 1 = 6      <- should be 7!

   Both threads read the SAME starting value because the read-modify-
   write wasn't protected -- one of the two increments is SILENTLY LOST.
   This can happen with shared memory (Chapter 7's threads), with
   files (two processes both checking "does this file exist" then
   both creating it), or with any shared, mutable resource touched from
   more than one thread of control without coordination.

   THIS IS WHY THIS GUIDE KEEPS RETURNING TO "ATOMIC" OPERATIONS
   (Chapter 18's locks, and this same theme appears throughout the
   companion security guide's race-condition chapter, in a web-request
   context) -- the fix is always some mechanism that makes the
   check-and-act sequence happen as ONE step nothing else can interleave
   with.
```

### Real-world example

```
   The Therac-25, a radiation therapy machine built in the 1980s, is
   the canonical, deeply-studied real-world case for why race conditions
   are taken so seriously in safety-critical engineering. Investigated
   in exhaustive detail by Nancy Leveson and Clark Turner (their 1993
   report remains the standard reference, taught in most software-
   engineering curricula), the machine's control software had a RACE
   CONDITION: if an operator entered treatment parameters QUICKLY enough
   -- specifically, by using keyboard edit commands to correct a mistake
   within about 8 seconds -- the software's checks for a dangerous
   configuration could be BYPASSED, because a shared variable
   (essentially a flag indicating whether safety checks had completed)
   could be read by one part of the program before another part had
   finished updating it. Between 1985 and 1987, this bug (and related
   issues in the same system) contributed to several patients receiving
   massive radiation overdoses, with fatal and severely injurious
   consequences.
   THE LESSON FOR THIS CHAPTER, stated as starkly as the real history
   warrants: the bug was HARD TO REPRODUCE (it depended on precise,
   fast operator timing, which most testers and most routine use never
   triggered) and existed in FIELDED, IN-USE equipment for YEARS before
   it was understood. "It passed testing" is not evidence a race
   condition doesn't exist -- it's evidence the specific interleaving
   that triggers it didn't happen to occur during testing.
```

### Try it yourself

```bash
# See a race condition happen live, harmlessly, in a tiny script
cat > race.py << 'EOF'
import threading

counter = 0
def increment():
    global counter
    for _ in range(100000):
        counter += 1   # NOT atomic: read, add 1, write back -- 3 steps

threads = [threading.Thread(target=increment) for _ in range(4)]
for t in threads: t.start()
for t in threads: t.join()
print(f"expected 400000, got {counter}")
EOF
python3 race.py   # run it a few times -- the result varies, and is
                  # usually LESS than 400000, because increments were lost
```

```
   Note: Python's Global Interpreter Lock (GIL) means this specific
   demo's "lost update" rate may be lower/less dramatic than the same
   code in a language without a GIL (like Go, Java, or C) -- but the
   race condition is real in ALL of them; only its VISIBILITY differs.
```

### Common mistakes

- **"It worked every time I tested it, so it's fine."** This is exactly what
  a race condition looks like from the outside — absence of failure in
  testing is weak evidence, not proof.
- **Assuming a single line of code is "one step" for the CPU/OS.** Even
  `counter += 1` is actually several separate machine instructions (read,
  add, write) — any of them can be interrupted between.
- **Only worrying about race conditions in explicitly multi-threaded code.**
  Two separate PROCESSES racing on a shared file, a shared database row, or a
  shared external resource is exactly the same class of bug.

### Check yourself

1. What makes a sequence of operations "not atomic," and why does that
   matter?
2. Why can a race condition pass testing repeatedly and still be a real bug?
3. In the Therac-25 case, what made the bug so hard to catch before it
   caused real harm?

*(Answers: Appendix E.)*

### Further reading

- **Report:** Leveson & Turner, "An Investigation of the Therac-25
  Accidents" (IEEE Computer, 1993) — the definitive, sobering account;
  freely findable online.
- **Book (free):** *OSTEP*, "Concurrency: An Introduction" chapter.
- **Guide:** this wiki's security guides cover the same underlying idea in a
  web-application context (check-then-act request races) if you want the
  application-layer version of this exact lesson.

---

## Chapter 18 — Locks, semaphores, and deadlock

### In one sentence

A lock lets you turn a multi-step operation into something that behaves as if
it were ONE uninterruptible step, at the cost of real risk: two locks taken in
inconsistent order can leave everyone waiting on each other, forever.

### Why this matters

Locks are the standard fix for Chapter 17's race conditions — and getting them
wrong (holding them too long, acquiring them in inconsistent order) creates a
different, equally real class of production incident: the deadlock.

### How it actually works

```
   A MUTEX (mutual exclusion lock) is the simplest form: only ONE
   thread/process can HOLD it at a time. Anyone else trying to acquire
   it BLOCKS (waits) until the holder releases it.

     lock(mutex)
       counter += 1     // now genuinely safe: nobody else can be in
                        // here at the same time
     unlock(mutex)

   A SEMAPHORE generalises this: it allows up to N holders at once
   (a mutex is just a semaphore with N=1) -- useful for things like
   "allow at most 5 concurrent connections to this resource."

   DEADLOCK is the classic failure mode of locking: two (or more) threads
   each hold a lock the OTHER one needs, and neither will ever release
   what it has:

     Thread A:  locks Resource1, then tries to lock Resource2 (blocked --
                B has it)
     Thread B:  locks Resource2, then tries to lock Resource1 (blocked --
                A has it)
     -> both wait FOREVER. Nothing will ever unblock either of them.

   FOUR CONDITIONS must ALL be true for deadlock to occur (the classic
   formulation, from Coffman et al., 1971) -- and breaking ANY ONE of
   them prevents it:
     1. MUTUAL EXCLUSION       resources can't be shared (inherent to
                              locks)
     2. HOLD AND WAIT          a thread holds one resource while waiting
                              for another
     3. NO PREEMPTION          a resource can't be forcibly taken away
     4. CIRCULAR WAIT          a cycle of threads each waiting on the
                              next

   THE MOST COMMON REAL-WORLD FIX: always acquire multiple locks in the
   SAME, CONSISTENT ORDER, everywhere in your codebase (e.g. always lock
   "lower resource ID before higher resource ID") -- this breaks
   condition 4, circular wait, and is simple enough to actually follow
   as a team-wide convention.
```

### Real-world example

```
   NASA's Mars Pathfinder rover, after landing successfully in July
   1997, began experiencing unexplained total system resets days into
   its mission -- a serious, publicly documented problem for an
   irreplaceable spacecraft millions of miles from Earth. The cause,
   diagnosed and famously explained afterward by the mission's lead
   developer Glenn Reeves, was a PRIORITY INVERSION bug, a close
   relative of this chapter's deadlock family: a LOW-PRIORITY task
   would occasionally acquire a mutex protecting a shared information
   bus, then get PREEMPTED by MEDIUM-priority tasks (which had nothing
   to do with that mutex at all) before it could finish and release the
   lock. Meanwhile, a HIGH-priority task that actually NEEDED that same
   mutex sat blocked, waiting -- and because the medium-priority tasks
   kept the LOW-priority lock-holder from ever running long enough to
   finish, the high-priority task could be starved long enough to trip
   a watchdog timer, causing the whole system to reset.
   Engineers on Earth DIAGNOSED and FIXED this remotely by enabling
   "priority inheritance" (a real-time-scheduling technique where a
   low-priority task temporarily BORROWS the priority of whoever is
   waiting on its lock, so it can finish and release quickly) --
   uploaded to the rover as a patch, resolving the resets. It remains
   one of the most famous, well-documented real demonstrations that
   locking bugs are not merely academic: they took down a Mars rover,
   and the fix is a named, standard technique specifically because this
   class of problem recurs.
```

### Try it yourself

```bash
# Fix the race condition from Chapter 17 with a proper lock, and see
# the difference
cat > lockfix.py << 'EOF'
import threading

counter = 0
lock = threading.Lock()

def increment():
    global counter
    for _ in range(100000):
        with lock:            # now the read-modify-write IS atomic
            counter += 1

threads = [threading.Thread(target=increment) for _ in range(4)]
for t in threads: t.start()
for t in threads: t.join()
print(f"expected 400000, got {counter}")   # now always correct
EOF
python3 lockfix.py   # run it several times -- always 400000 now
```

### Common mistakes

- **Holding a lock longer than necessary.** Do the minimum work inside the
  locked section — slow I/O or expensive computation inside a lock turns a
  small critical section into a bottleneck for every other thread waiting on
  it.
- **Acquiring locks in inconsistent order across different code paths.**
  This is THE most common real-world cause of production deadlocks —
  establish and enforce a consistent lock-ordering convention.
- **"My code doesn't have explicit locks, so it can't deadlock."** Database
  row locks, distributed locks (Redis, ZooKeeper), and other implicit
  locking mechanisms can deadlock exactly the same way — this isn't unique to
  hand-written mutexes.

### Check yourself

1. What does a mutex guarantee, and what's the difference between a mutex
   and a semaphore?
2. List the four conditions required for deadlock, and name one way to break
   the cycle.
3. In the Mars Pathfinder case, what specifically caused the high-priority
   task to be starved, and how was it fixed?

### Further reading

- **Book (free):** *OSTEP*, "Locks" and "Condition Variables" chapters.
- **Account:** Glenn Reeves' own account/interviews about the Mars
  Pathfinder priority-inversion bug — search "Mars Pathfinder priority
  inversion" for the original engineering writeups.
- **Reference:** `man 7 pthreads`, `man 3 pthread_mutex_lock`.

---

## Chapter 19 — Pipes, sockets, shared memory, and message queues

### In one sentence

When two separate processes need to exchange data — and they can't just share
memory the way threads do — the OS provides several distinct mechanisms
(pipes, sockets, shared memory segments, message queues), each with a
different trade-off between speed, simplicity, and how far apart the two ends
can be.

### Why this matters

`|` (the pipe operator) is the single most-used piece of Unix philosophy in
existence, and understanding what it actually IS (not just how to use it)
explains why Unix command-line tools compose so beautifully.

### How it actually works

```
   PIPES: a one-directional, in-kernel byte stream connecting two
   processes' file descriptors (Chapter 23) -- one writes, the other
   reads, and the kernel buffers the data between them. This is EXACTLY
   what `|` does in your shell:

     ls -la | grep ".txt"
     the shell: creates a pipe, forks two children (Chapter 6), connects
     `ls`'s stdout to the pipe's write end and `grep`'s stdin to the
     pipe's read end. `ls` never needs to know `grep` exists, or vice
     versa -- they just read/write file descriptors, and the KERNEL
     moves the bytes between them.

   Named pipes (FIFOs) work the same way but have a filesystem path,
   so unrelated processes (not just parent/child) can connect to them.

   SOCKETS: like pipes, but bidirectional and, crucially, can connect
   processes on DIFFERENT MACHINES (network sockets, TCP/UDP) or on
   the SAME machine (Unix domain sockets -- faster than network sockets
   for local-only communication, since they skip the whole network
   stack).

   SHARED MEMORY: the fastest IPC mechanism, because it skips COPYING
   data through the kernel entirely -- two processes map the SAME
   physical memory pages (Chapter 10) into their own virtual address
   spaces, and can read/write it directly. Fastest, but the RISKIEST:
   now you need explicit synchronisation (Chapter 18's locks) between
   processes, not just threads, because the OS gives you NO protection
   against both writing at once.

   MESSAGE QUEUES: the kernel maintains a queue of discrete MESSAGES
   (not a raw byte stream like a pipe) that processes can send to and
   receive from, with some ordering/priority semantics -- less commonly
   used directly today than the other three, but still real and
   available (POSIX message queues, System V message queues).
```

### Real-world example

```
   Docker's own architecture makes constant, direct use of a Unix
   domain socket: the Docker CLI you run (`docker ps`, `docker build`)
   talks to the Docker DAEMON over a Unix socket at
   `/var/run/docker.sock` by default (documented in Docker's own
   architecture docs) -- not a network socket, even though the
   communication LOOKS like a network API call (it's HTTP underneath).
   This is a deliberate choice, for two of this chapter's exact reasons:
   it's FASTER than a loopback TCP connection for same-machine
   communication, and it can use FILESYSTEM PERMISSIONS (Chapter 24) to
   control who's allowed to talk to the daemon at all -- which is
   exactly why "the docker.sock file's permissions" is a well-known,
   frequently-discussed real security consideration (anyone who can
   write to that socket can effectively control the whole Docker daemon,
   which usually means root-equivalent access to the host -- a fact
   documented extensively in container security guidance). You'll meet
   this exact socket again in Phase 3.
```

### Try it yourself

```bash
# See pipes in action, explicitly
mkfifo /tmp/mypipe
(echo "hello through a named pipe" > /tmp/mypipe &)
cat /tmp/mypipe
rm /tmp/mypipe

# See the anonymous pipe your shell creates for `|`
strace -f -e trace=pipe,pipe2 bash -c 'echo hi | cat' 2>&1 | grep pipe

# See the Docker socket, if you have Docker installed
ls -la /var/run/docker.sock 2>/dev/null || echo "Docker not installed yet -- see Phase 3"
```

### Build it in Go
Sockets are the IPC mechanism you'll use most. The
[TCP/IP guide's Go labs](../networking/tcp-ip/real-life-guide-v1.md#0-8-the-go-labs-build-the-network-tools-yourself) build TCP and UDP clients and
servers from scratch. [Chapter 22 there](../networking/tcp-ip/real-life-guide-v1.md#chapter-22-tcp-part-2-how-it-never-loses-your-data) shows why a
socket is a byte stream that needs message framing, the same lesson as this
chapter's pipes.

### Common mistakes

- **Assuming a pipe has unlimited buffer space.** Pipes have a FIXED kernel
  buffer size (commonly 64KB on Linux); a writer can BLOCK if the reader
  isn't keeping up and the buffer fills — this is exactly how `|` creates
  natural backpressure between commands.
- **Using shared memory "because it's fastest" without proper
  synchronisation.** The speed comes with zero built-in safety — you've
  opted into Chapter 17's and 18's exact problems, across process boundaries.
- **Forgetting Unix domain sockets are LOCAL-ONLY.** They cannot connect
  across machines — reach for a network socket the moment you need that.

### Check yourself

1. Name four IPC mechanisms and one key trade-off for each.
2. Why is a Unix domain socket typically faster than a TCP loopback
   connection for local-only communication?
3. Why does controlling who can write to `/var/run/docker.sock` matter so
   much for security?

### Further reading

- **Book (free):** *OSTEP*, IPC is lightly covered; supplement with Kerrisk's
  *The Linux Programming Interface* (the definitive, comprehensive reference)
  for pipes/sockets/shared memory/message queues in full depth.
- **Reference:** `man 7 pipe`, `man 7 unix`, `man 7 shm_overview`.
- **Docs:** Docker's own "Docker daemon socket option" documentation.

---

## Chapter 20 — Signals: the OS's tap on the shoulder

### In one sentence

A signal is a small, asynchronous notification the kernel (or another
process) can send to a process — "please stop," "please reload your
config," "you just tried to divide by zero" — and how a process chooses to
handle (or ignore) each one determines whether it shuts down gracefully or
just vanishes.

### Why this matters

Every graceful shutdown, every "why didn't my cleanup code run," and every
container orchestrator's stop sequence is built entirely on signals.

### How it actually works

```
   A SIGNAL is a number, delivered to a process, that INTERRUPTS its
   normal execution to run a handler (either a DEFAULT action the
   kernel provides, or a CUSTOM one your program registered). Common
   ones:

     SIGTERM (15)   "please terminate" -- the POLITE request. A
                    well-behaved program catches this, cleans up
                    (closes files, finishes in-flight work, flushes
                    buffers), and exits on its own terms. This is what
                    `kill <pid>` sends by DEFAULT (no flag needed).
     SIGKILL (9)    "terminate NOW, no negotiation." The kernel
                    forcibly ends the process -- it CANNOT be caught,
                    blocked, or ignored by the program (this is by
                    design: it's the guaranteed last resort). This is
                    what `kill -9 <pid>` sends. No cleanup code runs.
     SIGINT (2)     "interrupt" -- sent when you press Ctrl-C in a
                    terminal. Catchable, like SIGTERM.
     SIGHUP (1)     historically "the terminal hung up"; commonly
                    repurposed by daemons to mean "reload your
                    configuration without restarting."
     SIGSEGV (11)   "you accessed invalid memory" (Chapter 11's
                    segfault) -- delivered BY THE KERNEL when it
                    detects this, not sent deliberately by anyone.
     SIGCHLD (17)   "one of your child processes changed state" (e.g.
                    exited) -- how a parent process is notified without
                    having to constantly poll (Chapter 3's wait()).

   A PROCESS CAN, for most signals: use the DEFAULT action (often
   "terminate," sometimes "ignore," depending on the signal), install a
   CUSTOM HANDLER (run specific code when it arrives), or explicitly
   IGNORE it. SIGKILL and SIGSTOP are the two exceptions -- deliberately
   un-catchable and un-ignorable, so there's always a way to force a
   process to stop, no matter how badly it's misbehaving.
```

### Real-world example

```
   Kubernetes' pod termination sequence (documented in detail in
   Kubernetes' own official documentation) is this EXACT mechanism,
   used deliberately and by design: when a pod needs to stop (a
   deployment update, a scale-down, a node drain), Kubernetes sends
   SIGTERM to the container's main process FIRST, then waits up to a
   configurable "grace period" (30 seconds by default) for the process
   to exit on its own, and ONLY IF IT HASN'T exited by then, sends
   SIGKILL as a forced last resort.
   This is precisely why a well-known, extremely common real production
   bug is an application that does NOT properly handle SIGTERM
   (ignoring it, or handling it too slowly): during a rolling deploy or
   autoscaling event, in-flight requests get abruptly cut off when the
   grace period expires and SIGKILL arrives, because the application
   never used the polite warning it was given. This is documented
   extensively across cloud-native engineering blogs and is one of the
   most common causes of "why do we see a spike of errors during every
   deploy" investigations -- and the fix is almost always: catch
   SIGTERM in your application, stop accepting NEW work, finish
   IN-FLIGHT work, then exit cleanly, well within the grace period.
```

### Try it yourself

```bash
# Send signals and watch the difference
sleep 100 &
PID=$!
kill $PID           # SIGTERM (default) -- sleep has no custom handler,
                    # so it just uses the default action: terminate
wait $PID 2>/dev/null; echo "exit status: $?"

# Write a tiny program that catches SIGTERM and cleans up gracefully
cat > sigdemo.py << 'EOF'
import signal, time, sys

def handle_sigterm(signum, frame):
    print("caught SIGTERM -- cleaning up gracefully...")
    time.sleep(1)   # pretend to finish in-flight work
    print("done, exiting cleanly")
    sys.exit(0)

signal.signal(signal.SIGTERM, handle_sigterm)
print(f"PID {__import__('os').getpid()} running, send it SIGTERM")
while True:
    time.sleep(1)
EOF
python3 sigdemo.py &
PID=$!
sleep 1
kill $PID    # SIGTERM -- watch the graceful handler run
wait $PID
```

### Build it in Go
In Go, `signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt)` turns a
signal into a cancelled context, the cleanest way to start a graceful
shutdown. Complete examples: the
[HTTPS guide's production server](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown) (fail readiness,
drain, exit) and the [TCP/IP guide's draining backend](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries),
where a SIGTERM-aware backend lost 0 of 300 requests and `kill -9` lost 24.

### Common mistakes

- **"`kill` means force-kill."** Plain `kill <pid>` sends SIGTERM (a
  request), not SIGKILL (a command) — `kill -9` is the forceful one.
- **Never handling SIGTERM in a long-running service.** The process will be
  forcibly SIGKILLed after the grace period expires, with no chance to clean
  up — a direct, common cause of dropped requests during deploys.
- **Doing SLOW cleanup work in a SIGTERM handler.** If it takes longer than
  the grace period (Kubernetes' default: 30 seconds), you'll be SIGKILLed
  mid-cleanup anyway — keep the handler fast, or extend the grace period
  deliberately if the work genuinely needs longer.

### Check yourself

1. What's the difference between SIGTERM and SIGKILL, and why can't SIGKILL
   be caught or ignored?
2. What does `kill <pid>` (no flags) actually send?
3. Why does an application that ignores SIGTERM cause errors during a
   Kubernetes rolling deploy?

### Further reading

- **Reference:** `man 7 signal` — the authoritative list of every signal and
  its default action.
- **Book (free):** *OSTEP*'s process chapters mention signals; supplement
  with Kerrisk's *The Linux Programming Interface* for full depth.
- **Docs:** Kubernetes' own "Pod Lifecycle" documentation — the direct,
  real-world continuation of this chapter, covered again in Phase 3.

---

# Part 6 — Boot and the kernel

Everything so far assumed a system that's already running. This part goes
backward, to before any of it exists: what actually happens between pressing
the power button and seeing a login prompt, and where the kernel — the thing
implementing every mechanism from Parts 1–5 — fits into that sequence.

## Chapter 21 — From power button to login prompt

### In one sentence

Booting a computer is a chain of progressively-more-capable programs each
loading and handing off to the next — firmware, then a bootloader, then the
kernel, then user space — and every link in that chain is a place things can
go wrong, or (just as importantly) be attacked.

### Why this matters

"Why won't this server boot" (a genuinely common real incident) is a
diagnostic process that only makes sense once you know WHICH stage failed —
each stage fails differently and needs a completely different fix.

### How it actually works

```
   THE CHAIN, roughly, on a modern PC/server:

     1. FIRMWARE (UEFI, the modern successor to the old BIOS) runs
        FIRST, the instant power is applied -- built into the
        motherboard, independent of any disk. It initialises basic
        hardware, runs a Power-On Self-Test (POST), and looks for a
        BOOTLOADER to hand off to.

     2. THE BOOTLOADER (commonly GRUB on Linux) is a small program,
        itself loaded from disk by the firmware, whose entire job is
        to find and load the actual KERNEL into memory and start it.
        This is why you see a boot MENU (choose a kernel version, boot
        into recovery mode) at this stage -- the bootloader is where
        that choice happens, before the kernel is even running.

     3. THE KERNEL takes over: initialises itself, sets up memory
        management (Chapter 10), detects and initialises hardware
        (loading DEVICE DRIVERS as it goes, Chapter 22), mounts an
        initial, minimal root filesystem, and finally starts the very
        FIRST USER-SPACE PROCESS -- historically always PID 1.

     4. PID 1 (the INIT SYSTEM -- systemd on most modern Linux
        distributions, Chapter 22) takes over from here: it starts every
        other system service (networking, logging, your SSH daemon, your
        application) in the right order, and becomes the ultimate PARENT
        of every process on the machine (Chapter 3) -- if a process's
        original parent dies, PID 1 adopts it.

     5. Eventually, a LOGIN PROMPT (or, for a server, the point where
        SSH becomes reachable) appears -- the machine is now "up."

   EVERY STAGE HANDS OFF TRUST to the next: the firmware trusts the
   bootloader it loads, the bootloader trusts the kernel, the kernel
   trusts PID 1. This CHAIN OF TRUST is exactly what security features
   like UEFI SECURE BOOT are designed to protect (below).
```

### Real-world example

```
   UEFI SECURE BOOT exists specifically because of a well-documented,
   real attack category: BOOTKITS -- malware that infects the
   BOOTLOADER or early boot stages, running BEFORE the operating system
   (and therefore before any antivirus/endpoint-security software
   running IN the OS) has even started. A bootkit that compromises stage
   2 of this chain can control everything that loads afterward,
   including subverting the kernel itself, while remaining invisible to
   security tools that only start running once the OS is already up --
   a genuinely serious, well-documented real threat category (multiple
   real bootkit families have been publicly documented by security
   researchers over the years, targeting exactly this "trust chain,
   compromised at the earliest possible link" attack surface).
   Secure Boot's mechanism is a direct application of this chapter's
   "each stage hands off trust" model: the firmware CRYPTOGRAPHICALLY
   VERIFIES the bootloader's signature before running it, the bootloader
   verifies the kernel's signature before running IT, and so on --
   refusing to hand off execution to anything unsigned or tampered
   with, specifically closing the "compromise an early stage, own
   everything after it" attack this whole boot chain would otherwise be
   vulnerable to.
```

### Try it yourself

```bash
# See how your current machine actually booted
dmesg | head -30              # kernel's own boot-time log (may need sudo,
                              # and may only show recent entries on
                              # systems with dmesg buffer limits)
journalctl -b | head -30      # systemd's fuller boot log, if available

# See what your bootloader currently offers
cat /boot/grub/grub.cfg 2>/dev/null | grep -i "menuentry" | head -5 || \
  echo "GRUB config not found at the default path on this system"

# Time your own boot (useful, real diagnostic technique)
systemd-analyze 2>/dev/null || echo "systemd-analyze not available"
systemd-analyze blame 2>/dev/null | head -10   # which services were
                                               # slowest to start
```

### Common mistakes

- **Treating "won't boot" as one problem.** A firmware failure, a
  bootloader failure, a kernel panic, and an init-system failure all look
  different and need completely different fixes — figure out WHICH stage
  failed before troubleshooting.
- **Assuming Secure Boot is "just an annoying setting."** It's a genuine,
  meaningful defence against a real, documented attack class — disabling it
  casually (common when installing certain drivers or dual-booting) trades
  away real protection, which may or may not matter for your threat model.
- **Confusing the bootloader menu with the OS itself.** Seeing the GRUB menu
  means the FIRMWARE and BOOTLOADER stages succeeded — a subsequent failure
  is happening in the kernel or init stage, a different problem entirely.

### Check yourself

1. List the boot chain from power-on to login prompt, in order.
2. What is PID 1, and why does it matter that it becomes the "ultimate
   parent" of every process?
3. What specific attack does UEFI Secure Boot defend against, and how does
   it use this chapter's "chain of trust" idea to do it?

*(Answers: Appendix E.)*

### Further reading

- **Book (free):** *OSTEP* touches boot briefly; supplement with your
  distribution's own boot-process documentation (Arch Linux's wiki has an
  unusually clear, distro-agnostic explanation).
- **Reference:** `man 8 systemd-analyze`, `man 1 dmesg`, `man 5 journald.conf`.
- **Article:** search "UEFI Secure Boot explained bootkit" for current,
  concrete threat-model writeups.

---

## Chapter 22 — Init systems, kernel modules, and device drivers

### In one sentence

The init system (PID 1) starts and supervises every other process on the
machine in the right order, while kernel modules and device drivers are how
the kernel supports new hardware and features without needing to be entirely
recompiled for each one.

### Why this matters

`systemctl status`, `systemctl restart`, and `journalctl -u <service>` are
commands you will run constantly in any real Linux role — knowing what's
actually happening underneath makes them far more useful than memorised
incantations.

### How it actually works

```
   THE INIT SYSTEM'S JOB: start services (networking, logging, your
   database, your web server) in the correct ORDER (some depend on
   others being up first), keep track of what's supposed to be running,
   RESTART things that crash unexpectedly (if configured to), and clean
   up properly on shutdown.

   SYSTEMD (the default init system on most major Linux distributions
   today -- Debian, Ubuntu, RHEL/Fedora, and others) organises this
   around UNIT FILES: small config files describing a service, what it
   depends on, how to start/stop/restart it, and when it should run.

     $ systemctl status nginx        # is it running? recent log lines?
     $ systemctl start|stop|restart nginx
     $ systemctl enable nginx        # start automatically on boot
     $ journalctl -u nginx           # this service's own logs
                                     # (Chapter 34 covers this properly)

   KERNEL MODULES let the kernel load ADDITIONAL CODE -- typically
   device drivers, but also filesystem implementations and other
   optional functionality -- WITHOUT needing to be entirely recompiled.
   A module is loaded on demand (when the hardware/feature is actually
   needed) or explicitly:

     $ lsmod                    # what's currently loaded
     $ modprobe some_module     # load a module (and its dependencies)
     $ rmmod some_module        # unload one (if nothing's using it)

   A DEVICE DRIVER is the specific kernel code that knows how to talk to
   a particular piece of hardware (a network card, a disk controller, a
   GPU) -- translating the KERNEL's generic notion of "read some bytes
   from this block device" into the exact commands THAT SPECIFIC
   hardware chip understands. This is the layer that makes the rest of
   the OS's abstractions (Chapter 1's "lies") possible at all: the
   kernel doesn't need to know the details of every disk controller ever
   made, because each one's driver handles that translation.
```

### Real-world example

```
   The move from SysVinit (the older, decades-standard Unix init
   system, based on sequential shell scripts) TO systemd is one of the
   most publicly, extensively DEBATED changes in Linux's recent history
   -- Debian's own Technical Committee held a formal, documented,
   contentious vote on the matter in 2014 (its minutes and the public
   mailing-list debate are still findable and make for a genuinely
   interesting real look at how a major open-source project makes a
   consequential technical decision). systemd's designer, Lennart
   Poettering, argued publicly for its core advantages: PARALLELISED
   service startup (SysVinit's sequential scripts started services one
   at a time, in a fixed order, which was slow; systemd could start
   INDEPENDENT services simultaneously, dramatically speeding up boot
   times on modern multi-core hardware) and built-in SERVICE
   SUPERVISION (automatically restarting a crashed service, something
   SysVinit didn't do natively). Critics raised real, substantive
   concerns about scope creep (systemd grew to manage far more than
   just process startup -- logging, networking, and more) and
   architectural philosophy (Debated: whether bundling so much
   functionality into one project violates the traditional Unix
   philosophy of small, focused tools). Today, systemd is the default
   on most major distributions, but the debate itself remains a genuinely
   useful real case study in trade-offs between "simple and well
   understood" and "more capable but more complex."
```

### Try it yourself

```bash
# Explore your own init system
systemctl list-units --type=service --state=running | head -15
systemctl status ssh 2>/dev/null || systemctl status sshd 2>/dev/null

# See which services are set to auto-start on boot
systemctl list-unit-files --type=service --state=enabled | head -10

# See kernel modules currently loaded, and pick one to inspect
lsmod | head -10
modinfo $(lsmod | awk 'NR==2{print $1}') 2>/dev/null | head -10

# See how long your boot took, and what was slowest
systemd-analyze 2>/dev/null
systemd-analyze blame 2>/dev/null | head -5
```

### Common mistakes

- **Using `service <name> restart` and `systemctl restart <name>`
  interchangeably without knowing which init system you actually have.**
  Most modern distros alias the old commands to systemd equivalents, but
  don't assume this on every system you touch.
- **Forgetting `systemctl enable` vs `systemctl start` are different.**
  `start` runs it now; `enable` makes it start automatically on the NEXT
  boot. You often want both, and forgetting `enable` is a common real cause
  of "the service was fine, then the server rebooted and it wasn't running."
- **Assuming a loaded kernel module means the hardware is working
  correctly.** The module loading successfully means the DRIVER is present;
  it doesn't guarantee the underlying hardware itself is functioning
  correctly.

### Check yourself

1. What is a unit file, and what does `systemctl enable` do that `systemctl
   start` doesn't?
2. What is a kernel module, and why does the kernel support loading them
   dynamically instead of compiling everything in?
3. Name one concrete advantage systemd's designers argued for over
   SysVinit, and one substantive criticism raised against it.

### Further reading

- **Reference:** `man 1 systemctl`, `man 5 systemd.unit`, `man 8 modprobe`.
- **Docs:** `freedesktop.org`'s systemd documentation — the authoritative
  source, maintained by the project itself.
- **History:** search "Debian technical committee init system 2014" for the
  public debate and its resolution — a genuinely worthwhile real case study.

---

### End of Phase 1 — Milestone check

- [ ] I can explain what "the OS" is actually doing for every program I run,
      in plain language
- [ ] I understand the kernel/user-space boundary and what a syscall is
- [ ] I can trace a process's full life: fork, exec, running, exiting, and
      being reaped
- [ ] I know the difference between a thread and a process, and why that
      distinction matters for safety and performance
- [ ] I can read process states (`R`/`S`/`D`/`Z`/`T`) and use them to
      diagnose "what is this process actually doing"
- [ ] I understand virtual memory, page faults, and why `free -h`'s numbers
      often confuse people
- [ ] **I have watched a race condition produce a wrong answer, then fixed
      it with a lock**
- [ ] I know what a signal is, and specifically why SIGTERM matters for
      graceful shutdown
- [ ] I can describe the boot chain from power-on to login prompt

---

# PHASE 2 — Linux fundamentals

You now know what the OS is doing underneath. This phase is fluency: the
commands and mental models you'll use every single day, hands-on, until they
stop feeling like memorisation.

# Part 7 — The filesystem hierarchy and navigation

Parts 1–6 covered the OS's internals; this part and the next few turn
practical. It starts with Linux's own organizing idea — that (almost)
everything is a file — and the directory hierarchy that idea produces: the
map you need before the command-heavy parts that follow make any sense.

## Chapter 23 — Everything is a file (almost)

### In one sentence

Unix's defining design choice is that files, directories, devices, pipes, and
even kernel information are all accessed through the SAME simple interface —
open, read, write, close — which is exactly why Linux's command-line tools
compose so well together.

### Why this matters

Once you internalise this, half of Linux stops looking like a pile of
unrelated commands and starts looking like one consistent idea, applied
everywhere.

### How it actually works

```
   A FILE DESCRIPTOR (Chapter 2 introduced this) is just a small integer
   a process uses to refer to something it has open. Every process gets
   three, by convention, at startup:
     0   stdin    (standard input)
     1   stdout   (standard output)
     2   stderr   (standard error)
   But a file descriptor can point at FAR more than a regular file:
   a directory, a device (`/dev/sda`, your disk), a pipe (Chapter 19), a
   network socket, or a purely VIRTUAL file the kernel generates on the
   fly.

   THE VIRTUAL FILESYSTEMS ARE THE CLEAREST DEMONSTRATION of this
   philosophy taken all the way:
     /proc    a filesystem that doesn't correspond to anything on disk
             AT ALL -- it's the kernel exposing its OWN internal state
             (running processes, memory stats, CPU info) as if it were a
             directory of files, generated live, on read.
     /sys     similar idea, focused on hardware/device and kernel
             subsystem configuration.
     /dev     special files representing DEVICES -- `/dev/null` (discards
             anything written to it, always returns end-of-file when
             read), `/dev/random` (returns random bytes), `/dev/sda`
             (your actual disk, as a stream of bytes).

   THIS IS NOT JUST A CUTE ABSTRACTION -- it's WHY `ps`, `top`, and
   `free` (Chapter 31) work at all: they don't use some special, secret
   kernel API. They just READ FILES under /proc, the exact same `open`/
   `read`/`close` syscalls (Chapter 2) your own program would use to
   read a text file. You could write your own (much worse) version of
   `ps` in a few lines of shell script, entirely from `/proc`.
```

### Real-world example

```
   The "everything is a file" philosophy didn't originate with Linux --
   it traces back to Unix's earliest design (Ken Thompson and Dennis
   Ritchie at Bell Labs, early 1970s) and was pushed even FURTHER by a
   later Bell Labs research OS called Plan 9, whose designers extended
   the idea to networking, graphics, and even PROCESS CONTROL, all
   exposed as file operations. Linux's `/proc` filesystem (introduced in
   the early 1990s) is a direct descendant of this same philosophy,
   deliberately choosing to expose kernel internals as a browsable,
   readable filesystem rather than inventing a separate, bespoke query
   API. The practical payoff, today, is enormous and very real: any
   tool that can read a file -- `cat`, `grep`, a shell script, a
   monitoring agent, a config-management tool -- can ALSO inspect
   process info, memory stats, or CPU details, with zero special
   libraries, because it's just text files. Countless real monitoring
   and diagnostic tools (from decades-old shell scripts to modern
   observability agents) are built on nothing more exotic than reading
   `/proc` on a schedule.
```

### Try it yourself

```bash
# See your own process's file descriptors, live
ls -la /proc/self/fd/

# Watch `/dev/null` do exactly what it claims
echo "this vanishes" > /dev/null
cat /dev/null; echo "(nothing printed above -- always empty)"

# Read raw kernel information as plain text -- no special tool needed
cat /proc/cpuinfo | head -10
cat /proc/meminfo | head -5
cat /proc/version

# Write your OWN tiny "ps" using nothing but /proc and shell tools
for pid in /proc/[0-9]*; do
  [ -r "$pid/comm" ] && echo "$(basename $pid): $(cat $pid/comm 2>/dev/null)"
done | head -10
```

### Build it in Go
Everything `ps` and `top` display comes from files under `/proc`.
[Chapter 79](#chapter-79-reading-proc-from-go-build-your-own-ps) builds a `ps`/`top` in Go from
`/proc/<pid>/stat`, `status`, and `cmdline`, including the parsing trap where a
process name contains spaces.

### Common mistakes

- **Thinking `/proc` files are "real" files on disk.** They're generated
  live, on read, by the kernel — there's nothing to find if you look for them
  on the physical disk.
- **Assuming every kernel object behaves identically to a regular file.**
  Devices like `/dev/sda` can be read/written like a file, but writing the
  wrong bytes to the wrong offset can genuinely destroy data — this power
  comes with real risk.
- **Forgetting stdin/stdout/stderr ARE file descriptors 0/1/2.** Redirection
  syntax (`2>&1`, `>`, `<`) is just telling the shell which file descriptor
  points where — not a separate, unrelated concept from everything else in
  this chapter.

### Check yourself

1. What does "everything is a file" actually mean, concretely?
2. Name the three standard file descriptors every process starts with.
3. Why can a tool like `ps` be implemented entirely by reading files, with no
   special kernel API?

*(Answers: Appendix E.)*

### Further reading

- **Reference:** `man 5 proc` — the exhaustive list of what's exposed under
  `/proc`.
- **Book (free):** *OSTEP*, "Files and Directories" chapter.
- **History:** search "Plan 9 everything is a file Bell Labs" for the
  research-OS lineage of this idea.

---

## Chapter 24 — The Filesystem Hierarchy Standard, and why `/etc` isn't random

### In one sentence

Linux's top-level directory layout (`/etc`, `/var`, `/usr`, `/home`, and the
rest) follows a documented, decades-old convention — the Filesystem Hierarchy
Standard — specifically so that any tool, script, or engineer can predict
where something lives without having to ask.

### Why this matters

Knowing this layout by MEANING, not by memorised list, means a Linux box
you've never touched before is still navigable in minutes.

### How it actually works

```
   THE CORE DIRECTORIES, and WHY each one exists (this is the part worth
   actually remembering -- the meaning, not just the name):

     /etc      "editable text configuration" -- system-wide config files.
              Nothing here should be a compiled binary or a program's
              own data; it's SETTINGS.
     /var      "variable data" -- files that CHANGE while the system
              runs: logs (`/var/log`), databases, mail queues, caches.
              The opposite of "static."
     /usr      historically "Unix System Resources" (not "user"!) --
              the bulk of installed software: binaries (`/usr/bin`),
              libraries (`/usr/lib`), shared data. On modern systems,
              `/bin` and `/sbin` are often just SYMLINKS into
              `/usr/bin` (see the real-world example).
     /home     personal directories for each human user.
     /root     the root user's OWN home directory -- deliberately
              SEPARATE from `/home`, partly so root's home still exists
              and is usable even if `/home` is a separate, unmounted
              filesystem.
     /tmp      temporary files -- may be CLEARED on every reboot (often
              literally a RAM-backed filesystem, tmpfs) -- never store
              anything here you can't afford to lose.
     /opt      "optional" -- self-contained third-party software that
              doesn't want to scatter itself across `/usr`.
     /dev, /proc, /sys   the virtual filesystems from Chapter 23.
     /boot     the kernel, bootloader config, and other files needed
              extremely early in Chapter 21's boot process.

   THE OVERALL PRINCIPLE: separate STATIC from VARIABLE, separate
   SYSTEM-PROVIDED from USER-INSTALLED from USER-PERSONAL, and put
   things a script can DEPEND ON in a small number of PREDICTABLE
   places. This predictability is the entire point -- it's a decades-old
   social/technical CONVENTION (documented as the "Filesystem Hierarchy
   Standard," maintained by the Linux Foundation), not a technical
   requirement the kernel enforces.
```

### Real-world example

```
   Around 2012-2020 (adoption spread gradually across distributions),
   many major Linux distributions -- including Fedora (an early mover,
   around 2011-2012) and later Debian and Ubuntu -- carried out what's
   widely known as the "UsrMerge": consolidating `/bin`, `/sbin`, and
   `/lib` (which had historically been SEPARATE top-level directories,
   for a specific historical reason -- ensuring critical early-boot
   binaries were available even if `/usr` was on a separate filesystem
   that hadn't mounted yet) into being SYMLINKS pointing into
   `/usr/bin`, `/usr/sbin`, and `/usr/lib` respectively.
   This is a genuinely useful, real example of TWO things this chapter
   teaches: first, that the FHS is a living, evolving convention (not
   an immutable law) that distributions can and do change when the
   original justification (limited early-boot filesystems) stops being
   relevant to how modern systems actually boot (initramfs-based
   early boot made the original separation largely unnecessary); and
   second, that the CONVENTION'S PREDICTABILITY survived the change --
   `/bin/ls` still works exactly as every script assumed it would, even
   though it's now "really" `/usr/bin/ls` underneath, specifically
   because the migration preserved backward compatibility via symlinks
   rather than breaking decades of scripts that hard-coded the old
   paths.
```

### Try it yourself

```bash
# See the top-level layout, and notice which of these are symlinks
ls -la /

# Confirm whether YOUR system has done the UsrMerge
ls -la /bin    # if this shows "-> usr/bin", you're on a merged system

# Explore what actually lives in each key directory
ls /etc | head -10        # configuration
ls /var/log | head -10    # logs (Chapter 34)
du -sh /usr /var /etc /tmp 2>/dev/null   # relative sizes -- /usr is
                                        # usually by far the biggest
```

### Common mistakes

- **Storing important data in `/tmp`.** It can be wiped on reboot (and
  often is, by design) — use `/var` or a dedicated data directory instead.
- **Guessing paths instead of checking.** `which <command>`, `man
  hier` (the FHS's own man page), and `ls` beat memorised assumptions,
  especially across different distributions that may vary slightly.
- **Assuming the FHS is identical on every Unix-like system.** macOS and
  BSD systems follow RELATED but not identical conventions — don't assume
  Linux paths transfer directly.

### Check yourself

1. What's the meaningful difference between `/etc` and `/var`?
2. Why does `/root` exist separately from `/home`?
3. What was the "UsrMerge," and why did it happen?

### Further reading

- **Reference:** `man 7 hier` — the Filesystem Hierarchy Standard's own man
  page, on almost every Linux system already.
- **Standard:** the Linux Foundation's official FHS specification document.
- **History:** search "UsrMerge Fedora Debian" for the real technical
  rationale and migration writeups.

---

## Chapter 25 — Permissions, ownership, and links

### In one sentence

Every file has an owner, a group, and three sets of read/write/execute
permissions (for the owner, the group, and everyone else) — and two special
kinds of files, hard links and symbolic links, let a single piece of data be
reachable by more than one name.

### Why this matters

Permission mistakes are one of the most common real causes of both "why can't
I access this" AND serious security incidents — getting this model genuinely
right (not just cargo-culting `chmod 777`) matters.

### How it actually works

```
   THE PERMISSION MODEL: `ls -l` shows something like:
     -rwxr-xr--  1  alice  developers  1234  Jan 1 12:00  script.sh
      |||||||||
      |||||||`- OTHERS: read only (r--)
      ||||`----- GROUP: read + execute (r-x)
      `--------- OWNER (alice): read + write + execute (rwx)

   Three PERMISSIONS (read, write, execute), applied to three CATEGORIES
   of user (owner, group, everyone else). For a FILE: read = can view
   contents, write = can modify, execute = can RUN it as a program. For
   a DIRECTORY, the meanings shift slightly: read = can LIST its
   contents, execute = can ENTER it / access files inside by name (even
   without read), write = can CREATE/DELETE entries inside it.

     chmod u+x script.sh      # add execute for the owner (user)
     chmod g-w file.txt       # remove write for the group
     chmod 755 script.sh      # numeric form: owner=7(rwx), group=5(r-x),
                              # other=5(r-x) -- read each digit as a
                              # 3-bit binary sum: r=4, w=2, x=1

   OWNERSHIP: every file has an OWNING USER and an OWNING GROUP
   (`chown`, `chgrp`). This, combined with the permission bits above,
   is the FULL access-control model for a standard Unix file (Chapter
   36 covers more advanced access control beyond this).

   HARD LINKS vs. SYMBOLIC LINKS (both introduced briefly in Chapter
   15, formalised here):
     HARD LINK (`ln target linkname`)    another directory entry
       pointing at the SAME INODE. Indistinguishable from "the
       original" -- there IS no original, just multiple names for one
       piece of data. Can't cross filesystems, can't link a directory.
     SYMBOLIC LINK (`ln -s target linkname`)    a SEPARATE, tiny file
       that just CONTAINS A PATH STRING pointing at the target. Can
       cross filesystems, can link directories, but BREAKS ("dangling
       link") if the target is moved or deleted -- because it's just
       a path, not a direct pointer to the data.
```

### Real-world example

```
   The SUID (Set User ID) permission bit is a real, historically
   important, and still occasionally security-relevant extension of
   this model: a program with the SUID bit set runs with the
   PERMISSIONS OF ITS OWNER, not the permissions of whoever launched it.
   The textbook LEGITIMATE use, present on essentially every Unix-like
   system: `/usr/bin/passwd` (the program you run to change your own
   password) is owned by root and has the SUID bit set, specifically so
   an ORDINARY user can run it and have it write to `/etc/shadow` (the
   system password file, which ordinary users otherwise cannot write to
   at all) -- SUID is precisely the narrow, deliberate escape hatch that
   makes this possible without giving every user general write access
   to a security-critical file.
   The FLIP SIDE is real and well documented: a MISCONFIGURED SUID
   binary -- one that's SUID-root but has a bug, or is writable/
   replaceable by a non-root user, or simply shouldn't have SUID set at
   all -- is a classic, extremely common PRIVILEGE ESCALATION vector,
   the subject of a large fraction of "local privilege escalation" CVEs
   and CTF/security-training exercises (finding stray SUID binaries via
   `find / -perm -4000` is a standard first step in real penetration
   testing methodology, covered from the offensive side in this wiki's
   security guides). This is a genuinely good, concrete illustration of
   why understanding the permission model precisely -- not just "chmod
   777 and move on" -- is a real, practical security skill, not academic
   trivia.
```

### Try it yourself

```bash
# See permissions and ownership, and change them
touch demo.txt
ls -l demo.txt
chmod u+x demo.txt; ls -l demo.txt
chmod 644 demo.txt; ls -l demo.txt   # back to a normal, sane default

# See hard links vs symlinks in action
echo "original data" > original.txt
ln original.txt hardlink.txt
ln -s original.txt symlink.txt
ls -li original.txt hardlink.txt symlink.txt   # -i shows inode numbers:
                                               # hardlink SHARES one,
                                               # symlink has its OWN
rm original.txt
cat hardlink.txt      # still works
cat symlink.txt        # BROKEN -- "No such file or directory"
                       # (dangling symlink, exactly as this chapter warns)
rm -f hardlink.txt symlink.txt

# Find SUID binaries on your own lab machine (informational, safe)
find /usr/bin /usr/sbin -perm -4000 2>/dev/null
```

### Common mistakes

- **`chmod 777` as a reflex "fix" for permission errors.** It removes ALL
  access control on that file/directory — almost always the wrong fix; find
  the actual owner/group mismatch instead.
- **Confusing "no read permission" with "file doesn't exist."** A
  permission-denied error and a not-found error look similar in casual
  reading but mean very different things — read the actual error message.
- **Assuming a symlink and its target are "the same file" for permission
  purposes.** Symlink permissions are largely cosmetic on Linux (the
  TARGET's permissions are what actually apply when you follow it) — this
  trips people up constantly.

### Check yourself

1. Read `-rwxr-xr--` and explain exactly what each group of three characters
   means.
2. What's the difference between a hard link and a symbolic link, concretely?
3. What does the SUID bit do, and give both a legitimate and a risky real
   use of it.

### Further reading

- **Reference:** `man 1 chmod`, `man 1 chown`, `man 2 stat`.
- **Cheat sheet:** any "chmod calculator" or permission-bits reference table
  — useful until the octal math becomes second nature.
- **Guide:** this wiki's security guides cover privilege escalation via
  SUID/capabilities in much greater depth, from the offensive side.

---

# Part 8 — Essential commands, in depth

Part 7 gave you the map; this part is the toolbox. These are the commands
you'll type dozens of times a day, covered in enough depth that you
understand what each one actually does to the filesystem underneath —
not just which flags produce which output.

## Chapter 26 — File and directory operations you'll use daily

### In one sentence

`cp`, `mv`, `rm`, `mkdir`, and their close relatives look trivial and are
genuinely, permanently dangerous when used carelessly — because Linux, by
design, does not ask "are you sure?" before deleting your data.

### Why this matters

More real data loss traces back to a mistyped `rm` than almost any other
single command — knowing exactly what these commands do, and building the
right habits, is part of using Linux professionally.

### How it actually works

```bash
   cp source dest          copy a file
   cp -r sourcedir destdir copy a directory, RECURSIVELY (required for
                          directories -- plain `cp` refuses)
   mv source dest           move OR rename (same command; a rename is
                          just a move within the same directory --
                          Chapter 15's directory-entry model, again)
   rm file                  delete a file. NO confirmation, NO recycle
                          bin, NO undo, by default.
   rm -r directory           delete a directory and everything in it,
                          RECURSIVELY. Combined with -f (force, skip
                          confirmations) this is the single most
                          dangerous common command in daily use.
   mkdir -p a/b/c            create a/b/c, creating any missing parent
                          directories along the way (without -p, this
                          fails if `a` or `a/b` don't already exist)
   touch file                create an empty file, or update an
                          existing file's modification timestamp
```

**Habits that prevent real disasters:**

```bash
   # Always double-check a wildcard/variable BEFORE running rm with it
   echo rm -rf "$TARGET_DIR"/*      # print first, remove the `echo`
                                   # once you've confirmed it's correct

   # Prefer `rm -i` (interactive, asks per file) while learning
   alias rm='rm -i'                 # add to your shell config

   # For anything precious, copy before you touch it
   cp -a important_dir important_dir.bak
```

### Real-world example

```
   In 2015, Valve's Steam client for Linux shipped an uninstall script
   containing (roughly) the command `rm -rf "$STEAMROOT/"*`. The bug,
   widely reported and acknowledged by Valve at the time: under certain
   conditions, the `$STEAMROOT` variable could end up EMPTY (unset or
   blank) before this line ran -- and an empty variable in that position
   silently turns the command into `rm -rf /*`, recursively deleting
   EVERYTHING the user running it had permission to delete, starting
   from the root of the filesystem. Multiple users reported losing their
   entire home directories as a result.
   THE LESSON, and it is a genuinely common real bug PATTERN, not just a
   one-off: a destructive command that depends on a VARIABLE being
   correctly set, with NO validation that the variable is actually
   non-empty before the destructive part runs, is a loaded gun. This is
   exactly why production-grade scripts (Chapter 39 covers this
   properly) check that critical variables are set before using them in
   anything destructive, and why the "echo it first" habit above is
   worth building now, permanently, before you're the one writing the
   script that runs as root against a production filesystem.
```

### Try it yourself

```bash
# Safe practice, in your lab only
mkdir -p ~/lab/playground && cd ~/lab/playground
mkdir -p a/b/c
touch a/b/c/file1.txt a/b/file2.txt
cp -r a a_backup
mv a_backup a_renamed
find . -type f              # see the whole tree

# Deliberately reproduce (harmlessly) the class of bug from the
# real-world example -- see WHY checking your variable matters
TARGET=""                    # simulate the bug: empty variable
echo rm -rf "$TARGET"/*      # print it -- notice this would expand to
                             # `rm -rf /*` if TARGET is empty. NEVER
                             # remove the echo on a command like this
                             # without checking TARGET first.

rm -r a a_renamed            # clean up, safely, because we know exactly
                             # what's in this directory
```

### Common mistakes

- **Running a destructive command with a variable, without checking it's
  set first.** The Steam incident, exactly. Always validate, or `echo` first.
- **Assuming `rm` has an undo.** On most Linux filesystems, by default, it
  doesn't — the data is gone (Chapter 15: the inode's link count hits zero,
  the space is freed). Backups are the only real undo.
- **Using `mv` across filesystems and assuming it's instant.** Chapter 15
  explained why a same-filesystem move is nearly free — a cross-filesystem
  move has to actually copy the data first, then delete the original, and
  can be interrupted partway, unlike a same-filesystem rename.

### Check yourself

1. Why does `cp` require `-r` for directories but `rm` doesn't (by default)?
2. What specifically went wrong in the 2015 Steam Linux `rm -rf` incident?
3. What habit would have prevented it, and why is "echo it first" a good
   general practice for destructive commands?

*(Answers: Appendix E.)*

### Further reading

- **Reference:** `man 1 rm`, `man 1 cp`, `man 1 mv`, `man 1 mkdir`.
- **History:** search "Steam Linux rm -rf bug 2015" for Valve's own
  postmortem and community reports.
- **Practice:** `tldr rm` / `tldr cp` (the `tldr-pages` project) for
  quick, example-driven command references.

---

## Chapter 27 — Text processing: `grep`, `sed`, `awk`, and pipes

### In one sentence

`grep` finds lines, `sed` transforms them, and `awk` treats each line as
structured columns you can compute over — and chained together with pipes,
these three tools can replace an enormous amount of custom code.

### Why this matters

Logs are text. Config files are text. Command output is text. Fluency with
these three tools turns "I'd need to write a script for that" into a
one-liner you type in ten seconds, constantly, for the rest of your career.

### How it actually works

```bash
   grep "pattern" file        # print lines matching a pattern
   grep -i "pattern" file      # case-insensitive
   grep -v "pattern" file      # INVERT -- print lines that DON'T match
   grep -r "pattern" dir/      # search recursively through a directory
   grep -c "pattern" file      # just COUNT matching lines
   grep -E "regex|pattern" file # extended regex (alternation, etc.)

   sed 's/old/new/' file       # substitute the FIRST match per line
   sed 's/old/new/g' file      # substitute ALL matches per line (global)
   sed -n '5,10p' file         # print only lines 5 through 10

   awk '{print $1}' file       # print the FIRST whitespace-separated
                               # column of every line
   awk -F, '{print $2}' file    # same, but split on commas instead
   awk '$3 > 100 {print $1}' file  # print column 1, but only for lines
                               # where column 3 is greater than 100
```

**Chained together, this is where the real power shows up:**

```bash
   # Top 5 most common words in a text file
   cat file.txt | tr ' ' '\n' | sort | uniq -c | sort -rn | head -5

   # Every unique IP address in a web server log, with request counts
   awk '{print $1}' access.log | sort | uniq -c | sort -rn | head -10

   # Every ERROR line from today's log, with the message only
   grep "ERROR" app.log | awk -F'ERROR: ' '{print $2}'
```

### Real-world example

```
   A now-famous 1986 exchange, documented in the essay collection "A
   Quarter Century of Unix" and widely retold since, involved Donald
   Knuth (asked to demonstrate his "literate programming" style)
   writing a multi-PAGE, carefully structured program to solve a simple
   task: read a text file, find the N most frequently occurring words,
   print them sorted by frequency. Doug McIlroy -- one of Unix's original
   designers, and the person credited with inventing the PIPE itself --
   was asked to review it, and responded with a SIX-LINE shell pipeline
   that did the exact same job:

     tr -cs A-Za-z '\n' | tr A-Z a-z | sort | uniq -c | sort -rn | sed ${1}q

   (translate non-letters to newlines, one word per line -> lowercase
   everything -> sort -> count unique occurrences -> sort by count,
   descending -> take the top N)

   This became one of the most-cited illustrations of the Unix
   philosophy IN EXISTENCE precisely because it's real, it's short, and
   it demonstrates EXACTLY this chapter's lesson: small, focused tools
   (`tr`, `sort`, `uniq`), each doing ONE thing well, chained together
   with pipes (Chapter 19), can replace substantial custom programs --
   and the resulting pipeline is often easier to read, debug, and modify
   piece by piece than the equivalent bespoke code.
```

### Try it yourself

```bash
cd ~/lab/playground
cat > sample.log << 'EOF'
2024-01-01 10:00:01 INFO  request from 10.0.0.1
2024-01-01 10:00:02 ERROR request from 10.0.0.2 failed: timeout
2024-01-01 10:00:03 INFO  request from 10.0.0.1
2024-01-01 10:00:04 ERROR request from 10.0.0.3 failed: refused
2024-01-01 10:00:05 INFO  request from 10.0.0.1
EOF

grep ERROR sample.log
grep -c ERROR sample.log
sed 's/ERROR/PROBLEM/' sample.log
awk '{print $4}' sample.log             # just the "from"/IP-ish column
awk '/ERROR/ {print $4, $NF}' sample.log  # IP + last field, ERROR lines only

# Reproduce a version of the McIlroy pipeline yourself
echo "the quick brown fox the lazy dog the fox" | \
  tr ' ' '\n' | sort | uniq -c | sort -rn
```

### Common mistakes

- **Reaching for a full script before trying a one-line pipeline.** For
  text-shaped problems, try `grep`/`sed`/`awk`/`sort`/`uniq` first — it's
  often genuinely faster to write AND to read later.
- **Forgetting `sed`'s substitution is per-line by default.** Without the
  `g` flag, only the FIRST match on each line is replaced — a very common
  source of "why didn't this replace everything" confusion.
- **Using `awk` for arbitrarily complex logic.** It's brilliant for
  column-oriented, line-by-line processing — past a certain complexity, a
  proper script (Chapter 39) is more maintainable.

### Check yourself

1. What does `grep -v` do, and when is it useful?
2. Why does `sed 's/a/b/'` sometimes seem to "not work" on a line with
   multiple matches, and what fixes it?
3. What made the McIlroy/Knuth story such an enduring illustration of the
   Unix philosophy?

### Further reading

- **Book (free):** *The Unix Programming Environment* (Kernighan & Pike) —
  the classic, still-relevant text on exactly this philosophy.
- **Reference:** `man 1 grep`, `man 1 sed`, `man 1 gawk` (or `man 1 awk`).
- **Essay:** search "Doug McIlroy Knuth word count pipeline" for the full,
  original story.

---

## Chapter 28 — Redirection and the pipeline as a design pattern

### In one sentence

`>`, `>>`, `<`, and `|` are all the same underlying idea — rewiring which
file descriptor points where (Chapter 23) — and understanding that makes
every redirection trick predictable instead of memorised.

### Why this matters

Getting `2>&1` in the wrong place, or piping something that should have been
redirected, is one of the most common real shell-scripting bugs — and it's
entirely avoidable once the model clicks.

### How it actually works

```bash
   command > file        # redirect STDOUT (fd 1) to a file, OVERWRITING it
   command >> file        # same, but APPEND instead of overwrite
   command < file          # redirect STDIN (fd 0) to read FROM a file
   command 2> file         # redirect STDERR (fd 2) to a file
   command > file 2>&1     # redirect stdout to file, THEN point stderr
                          # at "wherever fd 1 currently points" (the
                          # file) -- ORDER MATTERS, this is the classic
                          # gotcha (see below)
   command 2>&1 > file     # THIS DOES SOMETHING DIFFERENT: stderr is
                          # pointed at wherever stdout CURRENTLY goes
                          # (the terminal), THEN stdout is redirected to
                          # the file -- stderr still goes to the
                          # terminal, not the file!
   command1 | command2      # connect command1's STDOUT directly to
                          # command2's STDIN, via a pipe (Chapter 19) --
                          # no file involved at all
```

**Why the order matters, explained properly:**

```
   `2>&1` means "make file descriptor 2 point at wherever file
   descriptor 1 CURRENTLY points, right now, at this moment in the
   command line" -- it's not a persistent link, it's a one-time copy of
   a destination, evaluated left to right.

     > file 2>&1     step 1: fd 1 -> file.  step 2: fd 2 -> (wherever
                     fd 1 now points) -> file. BOTH end up in the file.
     2>&1 > file      step 1: fd 2 -> (wherever fd 1 currently points)
                     -> the terminal (unchanged so far). step 2: fd 1 ->
                     file. fd 2 is STILL pointed at the terminal from
                     step 1 -- it never got updated again.
```

### Real-world example

```
   The widely-discussed real security debate around "curl | bash"
   (or `wget ... | sh`) style installation instructions -- a genuinely
   common real pattern, historically used by the official install
   instructions of several well-known tools -- is directly about THIS
   chapter's pipe mechanism, examined critically: piping a script
   DIRECTLY from a network download into a shell interpreter, with NO
   step where a human reviews what will actually run. Security
   researchers and engineering blog posts have repeatedly, publicly
   demonstrated real, concrete issues with this pattern: a server can
   detect it's being accessed by `curl` specifically (vs. a browser) and
   serve DIFFERENT content, a connection interrupted partway through can
   leave a shell executing a PARTIAL, syntactically-different script than
   intended, and a compromised download server or a machine-in-the-
   middle (this wiki's security guides cover this exact threat model)
   can inject arbitrary commands with no visible warning to the user at
   all.
   The safer, real-world-recommended alternative pattern -- also just
   this chapter's redirection, applied deliberately -- is: `curl -o
   install.sh <url>`, THEN actually read the file, THEN run it
   (`bash install.sh`) as a separate, reviewable step. Same underlying
   mechanism (Chapter 23's file descriptors and this chapter's
   redirection); very different risk profile, purely because of WHERE
   the human review step sits in the pipeline.
```

### Try it yourself

```bash
cd ~/lab/playground

# See the 2>&1 order gotcha for yourself
echo "stdout line" > /dev/null   # (setup, ignore)
{ echo "to stdout"; echo "to stderr" >&2; } > both.log 2>&1
cat both.log                     # BOTH lines are here

{ echo "to stdout"; echo "to stderr" >&2; } 2>&1 > stdout_only.log
cat stdout_only.log              # only "to stdout" -- stderr went to
                                 # your terminal instead, watch for it
                                 # printing above this line when you ran
                                 # the command

# See a pipeline moving data with no file at all
echo "hello world" | tr 'a-z' 'A-Z' | rev

rm -f both.log stdout_only.log
```

### Common mistakes

- **Writing `2>&1 > file` when you meant `> file 2>&1`.** The single most
  common redirection bug — order genuinely changes the behaviour, as shown
  above.
- **Piping when you meant to redirect, or vice versa.** `|` connects two
  COMMANDS; `>` connects a command to a FILE. Using `>` when you have two
  commands just overwrites a file named after your second "command."
- **Blindly running `curl <url> | bash`.** Download first, read it, then
  run it — the extra step costs seconds and closes a real, documented risk.

### Check yourself

1. What does `2>&1` actually do, precisely?
2. Why does `command > file 2>&1` behave differently from `command 2>&1 >
   file`?
3. What's the specific, documented risk with piping a downloaded script
   directly into a shell, and what's the safer alternative?

### Further reading

- **Reference:** `man 1 bash` (the "REDIRECTION" section) — the authoritative
  and surprisingly readable source.
- **Article:** search "curl pipe bash security considerations" for current,
  detailed writeups of the real risk and mitigations.
- **Book (free):** *The Linux Command Line* (William Shotts, free PDF) —
  Part 3 covers I/O redirection thoroughly.

---

## Chapter 29 — Finding things: `find`, `locate`, `which`, `grep -r`

### In one sentence

`find` searches the filesystem live (slower, always accurate, endlessly
flexible), `locate` searches a pre-built index (instant, can be stale), and
`which`/`type` tell you exactly which program a command name actually
resolves to — different tools for genuinely different questions.

### Why this matters

"Where is this config file," "which binary actually runs when I type this
command," and "what changed on this box in the last hour" are all real,
constant, and each has a specific right tool.

### How it actually works

```bash
   find /path -name "*.log"          # search BY NAME, live, recursive
   find /path -type f                # only regular files (-type d = dirs)
   find /path -mtime -1               # MODIFIED in the last 1 day
   find /path -size +100M             # LARGER than 100 megabytes
   find /path -perm -4000             # has the SUID bit set (Chapter 25)
   find /path -name "*.tmp" -delete    # find AND act -- delete matches
                                      # (careful -- test without -delete
                                      # first!)

   locate filename                    # instant, but reads a PRE-BUILT
                                     # INDEX (usually updated once daily
                                     # via a cron job, `updatedb`) --
                                     # can miss recently-created files

   which command                      # shows the FULL PATH of the
                                     # executable that would run, based
                                     # on your $PATH
   type command                        # similar, but also reveals if
                                     # it's a SHELL BUILTIN, alias, or
                                     # function, not just a file on disk

   grep -r "TODO" /path               # search file CONTENTS recursively
                                     # (different from `find`, which
                                     # searches by NAME/metadata, not
                                     # content)
```

### Real-world example

```
   `find` with time-based filters is a genuine, standard, real
   FIRST STEP in incident-response and intrusion-investigation
   workflows -- documented extensively across security and forensics
   guidance (including this wiki's own security guides): when
   investigating "was this box compromised, and when," a very common
   early command is something in the spirit of:
     find / -mtime -1 -type f 2>/dev/null
   ("show me every regular file modified in the last day") -- because
   an attacker who has planted a webshell, modified a config, or added
   a persistence mechanism (a cron job, a new SSH key, a systemd unit)
   has, almost by definition, MODIFIED FILES, and time-based `find`
   filters are one of the fastest ways to narrow "the entire filesystem"
   down to "the small number of things that actually changed recently"
   -- a genuinely practical, real application of a command that looks,
   on the surface, like a simple file-search utility.
```

### Try it yourself

```bash
cd ~/lab/playground

# Create some files with different ages (simulated) and sizes
touch old_file.txt
dd if=/dev/zero of=big_file.bin bs=1M count=5 2>&1 | tail -1
mkdir subdir && touch subdir/nested.txt

find . -type f                       # everything
find . -name "*.txt"                 # by name pattern
find . -size +1M                     # by size
find . -mtime -1                     # modified in the last day (should
                                     # be everything you just created)

which ls; type ls                    # compare the two
which cd; type cd                    # `which` may find nothing or be
                                     # confusing -- cd is a shell BUILTIN,
                                     # not a file on disk; `type` shows
                                     # this clearly

grep -r "TODO" . 2>/dev/null || echo "(no TODOs found, as expected)"

rm -rf big_file.bin old_file.txt subdir
```

### Common mistakes

- **Using `locate` and expecting it to find a file created five minutes
  ago.** Its index is only as fresh as the last `updatedb` run — use `find`
  when you need certainty about the CURRENT filesystem state.
- **Confusing `find`'s name search with content search.** `find -name
  "*.log"` finds files whose NAME matches; `grep -r "pattern"` finds files
  whose CONTENT matches — genuinely different questions.
- **Running `find ... -delete` without testing the search first.** Always
  run the `find` WITHOUT `-delete` first, read the output carefully, THEN
  add `-delete` once you're certain.

### Check yourself

1. What's the key trade-off between `find` and `locate`?
2. What's the difference between what `which` and `type` each tell you?
3. Why is a time-based `find` search a standard early step in incident
   investigation?

### Further reading

- **Reference:** `man 1 find`, `man 1 locate`, `man 1 type`.
- **Cheat sheet:** `tldr find` for quick, practical example patterns.
- **Guide:** this wiki's security guides cover forensic file-timeline
  analysis (Ch 55 in the advanced guide) in much greater depth.

---

## Chapter 30 — Package management: how software gets onto a Linux box

### In one sentence

A package manager installs software along with a formal record of exactly
what it depends on, so installing, upgrading, and removing software doesn't
silently break something else that quietly relied on the same files.

### Why this matters

"Dependency hell" was a real, painful, industry-wide problem before package
managers solved it properly — and knowing how your distro's package manager
actually works (not just the install command) matters the first time an
upgrade goes wrong.

### How it actually works

```bash
   # Debian/Ubuntu family (dpkg underneath, apt as the friendly frontend)
   sudo apt update                 # refresh the list of AVAILABLE packages
   sudo apt install nginx          # install a package (and its dependencies)
   sudo apt remove nginx           # remove it (config files may remain)
   sudo apt purge nginx            # remove it AND its config
   apt list --installed | grep nginx   # is it installed? which version?

   # RHEL/Fedora family (rpm underneath, dnf -- formerly yum -- as the
   # friendly frontend)
   sudo dnf install nginx
   sudo dnf remove nginx
   dnf list installed | grep nginx

   THE CORE IDEA, either family: a PACKAGE isn't just files -- it's files
   PLUS METADATA, including a list of what it DEPENDS ON (other
   packages, at specific minimum versions). The package manager reads
   this metadata and:
     - refuses to install a package if a dependency can't be satisfied
       (or, more usefully, AUTOMATICALLY installs the needed
       dependencies too)
     - refuses to REMOVE a package if something else still depends on it
       (without an explicit override)
     - tracks exactly what's installed, so it can be queried, upgraded,
       or cleanly removed as a unit, rather than "some files that are
       probably related to nginx, somewhere."
```

### Real-world example

```
   "Dependency hell" is a real, well-documented, widely-remembered
   industry term for the pain that existed BEFORE automatic dependency
   resolution was standard: installing one piece of software could
   require MANUALLY finding, downloading, and installing several OTHER
   pieces of software it depended on, each of which might themselves
   need something else, and upgrading ANY of them risked silently
   breaking something unrelated that happened to depend on the specific
   version you just replaced.
   Debian's `dpkg` (and later `apt`, layered on top specifically to
   solve automatic dependency resolution) and, separately, the
   introduction of `yum` on the RPM side (RPM itself, used by Red Hat
   from the mid-1990s, initially had NO automatic dependency
   resolution -- installing a package with unmet dependencies simply
   FAILED with an error, leaving the user to sort it out manually; `yum`
   was created specifically to add automatic resolution on top of RPM,
   analogous to what `apt` did for `dpkg`) are both well-documented,
   real engineering responses to exactly this problem. This history is
   worth knowing because it explains WHY package managers are built the
   way they are: the dependency-metadata-and-resolution machinery isn't
   bureaucratic overhead, it's the specific, hard-won fix for a genuinely
   painful problem the entire industry lived through.
```

### Try it yourself

```bash
# Explore what's already installed and why (Debian/Ubuntu example --
# adapt the commands if you're on an RPM-based distro)
apt list --installed 2>/dev/null | wc -l    # how much is on this box?

apt-cache depends curl 2>/dev/null | head -10   # what does `curl`
                                                # depend on?
apt-cache rdepends curl 2>/dev/null | head -10   # what DEPENDS ON curl?
                                                 # (reverse dependencies --
                                                 # this is WHY you can't
                                                 # just remove a package
                                                 # blindly)

# See a package's actual installed files
dpkg -L curl 2>/dev/null | head -10
```

### Common mistakes

- **Force-removing a package with `--force` flags to "fix" a dependency
  error.** This is how you actually END UP in a broken, inconsistent state —
  understand WHY the dependency conflict exists before overriding it.
- **Installing software by manually copying files instead of via the
  package manager.** Now the package manager has NO RECORD of it — future
  upgrades, dependency checks, and clean removal all become manual, error-
  prone work.
- **Not running `apt update` (refreshing the package LIST) before
  installing.** You may be attempting to install a version that's no longer
  available, or missing a newer one, from a stale local index.

### Check yourself

1. What extra information does a package carry, beyond the actual files?
2. Why did "dependency hell" happen before tools like `apt`/`yum` existed?
3. What's the difference between `apt remove` and `apt purge`?

### Further reading

- **Reference:** `man 8 apt`, `man 8 dnf`, `man 1 dpkg`, `man 8 rpm`.
- **History:** search "dependency hell history package management" for
  retrospectives on the problem these tools were built to solve.
- **Docs:** your specific distribution's own package-management
  documentation (Debian's, Ubuntu's, and Fedora's wikis are all thorough
  and well maintained).

---

# Part 9 — Process management in practice

Part 2 explained what a process is in theory. This part is the hands-on
counterpart — watching, controlling, and killing real processes on a real
system, using the tools you'll actually reach for during an incident, not
just to confirm you understand the concept.

## Chapter 31 — Watching and controlling processes: `ps`, `top`, `htop`, `kill`

### In one sentence

`ps` gives you a snapshot, `top`/`htop` give you a live, continuously-updating
view, and `kill` sends signals (Chapter 20) — together they're how you
actually SEE and ACT on Chapter 4's abstract "process" concept, in real time.

### Why this matters

"The server is slow, what's going on" almost always starts with one of these
three commands — they're the fastest way to turn a vague symptom into a
specific, actionable diagnosis.

### How it actually works

```bash
   ps aux                    # every process, classic BSD-style output
   ps -eo pid,ppid,cmd,%cpu,%mem --sort=-%cpu | head   # customised,
                             # sorted by CPU usage -- often more useful
                             # than the default columns
   ps -ef                    # every process, System V-style output
                             # (slightly different columns/conventions)

   top                        # live, auto-refreshing view. Press:
                              #   'q' to quit
                              #   'M' to sort by memory
                              #   'P' to sort by CPU (often the default)
                              #   'k' to kill a process, interactively
   htop                        # a much friendlier, colour, scrollable,
                              # mouse-clickable version of top (not
                              # always pre-installed -- worth adding)

   kill PID                    # send SIGTERM (Chapter 20) -- the polite
                              # request
   kill -9 PID                 # send SIGKILL -- the forceful, un-
                              # ignorable last resort
   kill -l                     # list every signal name/number available
   pkill -f "pattern"           # kill by matching the COMMAND LINE
                              # against a pattern, instead of a PID
   killall processname          # kill every process with this EXACT
                              # name
```

**Reading `ps`/`top` columns that confuse people at first:**

```
   %CPU    can exceed 100% on a multi-core machine -- it's the
           percentage of ONE CORE's capacity, so a process using two
           full cores shows ~200%.
   VSZ      virtual memory size (Chapter 10) -- often huge, often
           mostly meaningless on its own.
   RSS      resident set size -- memory ACTUALLY in physical RAM right
           now. Usually the number you actually care about.
   STAT     the process state (Chapter 5: R/S/D/Z/T, plus modifiers
           like `+` for foreground process group, `l` for multi-threaded)
```

### Real-world example

```
   A FORK BOMB is a genuine, well-documented, real demonstration of
   Chapter 4's "processes aren't free" lesson taken to its logical,
   destructive extreme -- a tiny piece of code (famously, in bash:
   `:(){ :|:& };:`) that defines a function which calls ITSELF TWICE,
   in the background, recursively, with no base case -- exponentially
   spawning new processes until the system's process table, memory, or
   CPU scheduling capacity is completely exhausted, typically making the
   machine unresponsive within seconds. This is widely documented across
   systems-security and Unix-history material as a classic, real
   denial-of-service technique (and a classic systems-administration
   cautionary tale) -- it's PRECISELY why real production systems set
   hard process-count limits per user (`ulimit -u`, Chapter 43) and why
   understanding tools like `ps`/`top`/`kill` -PRECISELY-, under
   pressure, matters: on a system actively being consumed by runaway
   process creation, you need to be able to IDENTIFY and KILL the
   right thing FAST, before the tools themselves become too starved of
   resources to even run.
```

```
   Dangerous: do NOT run a fork bomb outside a fully disposable,
   isolated VM you're prepared to hard-reboot. It can render a machine
   completely unresponsive, including to attempts to kill it.
```

### Try it yourself

```bash
# Explore process info safely
ps aux | head -10
ps -eo pid,ppid,cmd,%cpu,%mem --sort=-%mem | head -10

# Start something to practice killing
sleep 300 &
PID=$!
echo "started PID $PID"
ps -p $PID
kill $PID                       # SIGTERM
ps -p $PID 2>&1 || echo "gone"

# Practice pkill/killall (safe -- only affects OUR test processes)
sleep 300 & sleep 300 &
pkill -f "sleep 300"
jobs                             # should show them as Terminated

# Watch live process activity (press 'q' to exit)
top
```

### Common mistakes

- **Reading `%CPU` as capped at 100%.** On a multi-core machine, it's
  per-core — 300% means roughly three full cores' worth of work.
- **Reaching for `kill -9` as the default.** Try plain `kill` (SIGTERM)
  first, and give the process a moment to shut down cleanly (Chapter 20) —
  reserve `-9` for processes that genuinely won't respond.
- **Using `killall`/`pkill` with an overly broad pattern.** `pkill -f sleep`
  would match EVERY process with "sleep" anywhere in its command line,
  possibly including things you didn't intend — be as specific as possible.

### Check yourself

1. What's the difference between `ps` and `top`, fundamentally?
2. Why can `%CPU` show a number greater than 100?
3. What is a fork bomb, and why does understanding process tools matter
   specifically WHILE one is running?

### Further reading

- **Reference:** `man 1 ps`, `man 1 top`, `man 1 kill`, `man 1 pkill`.
- **Tool:** `htop`'s own documentation — worth installing and using daily
  once you've learned plain `top`.
- **History:** search "fork bomb history unix" for its documented origins
  and the standard defences against it.

---

## Chapter 32 — Job control: foreground, background, and `nohup`

### In one sentence

A shell can run commands in the foreground (you wait for it) or the
background (it runs while you keep typing), and understanding exactly what
happens to a background job when your terminal session ends is the
difference between a smooth long-running task and a lost one.

### Why this matters

Starting a multi-hour data migration over SSH, then having your connection
drop and the job silently die with it, is one of the most common, most
avoidable real production frustrations — and it's entirely explained by
Chapter 20's signals.

### How it actually works

```bash
   command &            # run in the BACKGROUND -- your shell keeps
                        # control immediately, doesn't wait
   jobs                  # list your shell's current background/stopped
                        # jobs
   fg %1                  # bring job 1 to the FOREGROUND (wait for it)
   bg %1                   # resume a STOPPED job in the background
   Ctrl-Z                  # SUSPEND (STOP, Chapter 5's `T` state) the
                        # current foreground job
   Ctrl-C                  # send SIGINT to the current foreground job
                        # (Chapter 20)

   nohup command &         # run in the background, AND ignore SIGHUP --
                        # so it survives your terminal/SSH session
                        # ending
   disown %1                # remove a job from your shell's job table
                        # WITHOUT killing it -- it keeps running, but
                        # your shell stops tracking it
```

**Why a background job can still die when you log out:** when your SSH
session (or terminal) closes, the shell process itself is typically sent
SIGHUP ("hang up" — a name dating back to actual telephone modems). By
default, many programs' handling of SIGHUP (or simply inheriting the shell's
own death) means background jobs started with a plain `&` can be killed too,
unless something protects them — which is exactly `nohup`'s job.

### Real-world example

```
   "I started a long-running migration/backup/build over SSH, my
   connection dropped, and the job died with it" is one of the single
   most common, most universally-experienced real frustrations in
   day-to-day Linux/ops work -- documented in essentially every "Linux
   tips" resource and endlessly rediscovered by engineers the hard way,
   usually mid-incident, at the worst possible time. The mechanism is
   exactly this chapter's: the SSH session closing sends SIGHUP toward
   the shell and, often, its children; a plain background job (`&` with
   no protection) can die right along with it.
   THE REAL, PRACTICAL FIXES, all addressing the same root cause:
     `nohup long_command &`        explicitly ignore SIGHUP
     `disown`                       detach an ALREADY-running job from
                                   the shell's job table
     `screen` / `tmux`              (properly) run the command inside a
                                   TERMINAL MULTIPLEXER session that
                                   keeps running independently of your
                                   SSH connection entirely -- you can
                                   disconnect and RE-ATTACH later, seeing
                                   the still-running session exactly as
                                   you left it. This is the tool most
                                   experienced engineers actually reach
                                   for by default for anything long-
                                   running, precisely because it survives
                                   network drops, not just deliberate
                                   logouts.
```

### Try it yourself

```bash
# See job control directly
sleep 100 &
jobs                    # shows [1]+ Running   sleep 100 &
fg                       # brings it to foreground -- now press Ctrl-C
                         # to actually stop it (or wait, or Ctrl-Z)

sleep 100 &
kill -STOP %1            # simulate Ctrl-Z from a script
jobs                      # shows Stopped
bg %1                      # resume it in the background
jobs                        # now Running again
kill %1                      # clean up

# See nohup protect a job from SIGHUP (simulated, safely)
nohup sleep 60 > /tmp/nohup_test.log 2>&1 &
echo "PID $! is protected from SIGHUP"
kill $!                      # clean up

# If tmux is installed, try it
tmux new -s demo 'echo "inside tmux"; sleep 5; echo done' 2>/dev/null || \
  echo "(tmux not installed -- worth adding for real long-running work)"
```

### Common mistakes

- **Assuming `&` alone is enough for a job that must survive logout.** It
  isn't, by default — pair it with `nohup`, `disown`, or better, run it
  inside `screen`/`tmux` from the start.
- **Forgetting `fg`/`bg` need a job NUMBER (`%1`) if you have more than
  one.** With only one background job, plain `fg`/`bg` works; with several,
  you need to specify which.
- **Starting a critical long-running task directly over a plain SSH session
  "because it'll probably be fine."** For anything that takes more than a
  few minutes, use `tmux`/`screen` as a default habit, not an afterthought.

### Check yourself

1. What signal typically threatens a background job when your SSH session
   closes, and why?
2. What does `nohup` actually do?
3. Why do experienced engineers often prefer `tmux`/`screen` over `nohup`
   for long-running work?

### Further reading

- **Reference:** `man 1 nohup`, `man 1 jobs`, `man 1 disown`, `man 1 tmux`.
- **Book (free):** *The Linux Command Line* — the job-control chapter covers
  this cleanly with worked examples.
- **Guide:** any good `tmux` cheat sheet — worth genuinely learning, not just
  skimming, given how often it pays off.

---

## Chapter 33 — systemd and services: how "always running" actually works

### In one sentence

A systemd service is a small, declarative file describing how to start a
program, when to start it, and (crucially) whether to restart it if it
crashes — turning "keep this running forever" from a manual chore into a
guarantee the init system enforces.

### Why this matters

The difference between a service that silently stays down after a 3am crash
and one that restarts itself automatically is, very often, a single line in
a unit file that nobody set.

### How it actually works

```ini
   # /etc/systemd/system/myapp.service -- a minimal, real unit file
   [Unit]
   Description=My Application
   After=network.target              # start AFTER networking is up

   [Service]
   ExecStart=/usr/bin/myapp --config /etc/myapp/config.yml
   Restart=always                    # THE critical line -- restart on
                                     # ANY exit, crash or otherwise
   RestartSec=5                       # wait 5 seconds before restarting
   User=myappuser                     # don't run as root unless you
                                     # genuinely need to

   [Install]
   WantedBy=multi-user.target         # start automatically at normal
                                     # boot
```

```bash
   sudo systemctl daemon-reload        # tell systemd to re-read unit
                                      # files after you edit one
   sudo systemctl enable myapp         # start automatically on boot
   sudo systemctl start myapp           # start it right now
   sudo systemctl status myapp          # is it running? recent log
                                      # lines (Chapter 34)?
   sudo systemctl restart myapp          # stop then start
```

**The `Restart=` options that matter most:**

```
   no          (the default if unset) never restart automatically
   always       restart no matter HOW it exited -- crash, clean exit,
               anything
   on-failure   restart only if it exited with a non-zero status/signal
               (Chapter 3) -- NOT if it exited cleanly (status 0)
```

### Real-world example

```
   A well-documented, extremely common real production failure pattern,
   discussed across countless postmortems and SRE writeups: a service
   is deployed and manually started (`systemctl start myapp`) but never
   actually `enable`d, and its unit file never sets `Restart=`. It runs
   fine for weeks or months. Then, one night, the process crashes --
   an unhandled exception, an out-of-memory kill (Chapter 13), a
   transient bug -- and because there's no `Restart=` directive, systemd
   makes NO attempt to bring it back. The service simply stays down,
   silently, until a human notices (often hours later, via a customer
   complaint or a downstream alert, rather than the actual crash itself)
   and manually restarts it.
   THIS EXACT GAP -- a real, common finding in real production
   incident postmortems across the industry -- is why `Restart=always`
   (or, more carefully considered, `Restart=on-failure` if a clean exit
   genuinely should stay stopped) is treated as close to a DEFAULT
   BEST PRACTICE for any service expected to run continuously, and why
   `systemctl enable` (survive a REBOOT) and setting a restart policy
   (survive a CRASH) are both treated as basic, non-optional hygiene for
   anything running in production -- two separate, easy-to-forget
   settings, each addressing a different real failure mode.
```

### Try it yourself

```bash
# Write and install a tiny real service (safe, in your lab)
cat > /tmp/heartbeat.sh << 'EOF'
#!/bin/bash
while true; do
  echo "heartbeat: $(date)"
  sleep 5
done
EOF
chmod +x /tmp/heartbeat.sh

sudo tee /etc/systemd/system/heartbeat.service > /dev/null << 'EOF'
[Unit]
Description=Heartbeat demo service
After=network.target

[Service]
ExecStart=/tmp/heartbeat.sh
Restart=always
RestartSec=2

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start heartbeat
sleep 6
sudo systemctl status heartbeat --no-pager   # see it running

# Prove Restart=always works: kill it, watch it come back
sudo kill $(pgrep -f heartbeat.sh)
sleep 4
sudo systemctl status heartbeat --no-pager   # should be running again,
                                             # with a new PID

# Clean up
sudo systemctl stop heartbeat
sudo systemctl disable heartbeat 2>/dev/null
sudo rm /etc/systemd/system/heartbeat.service /tmp/heartbeat.sh
sudo systemctl daemon-reload
```

### Common mistakes

- **`systemctl start` without `systemctl enable`.** Works fine until the
  next reboot, then the service simply doesn't come back — a very common
  real cause of "it worked yesterday, why is it down after maintenance."
- **No `Restart=` policy on a production service.** Silent, extended
  downtime after any crash, with no automatic recovery.
- **Forgetting `daemon-reload` after editing a unit file.** systemd caches
  unit file contents — your edits won't take effect until you tell it to
  re-read them.

### Check yourself

1. What does `Restart=always` do, and what's the risk of leaving it unset?
2. What's the difference between `systemctl start` and `systemctl enable`?
3. Why did the service in the real-world example stay down for hours after
   crashing?

### Further reading

- **Reference:** `man 5 systemd.service`, `man 5 systemd.unit`.
- **Docs:** `freedesktop.org`'s systemd documentation, "Converting SysVinit
  scripts" guide — useful even if you're writing units from scratch, for
  the concepts it explains along the way.
- **Practice:** write a unit file for a real small script you use — the best
  way to make this concrete.

---

## Chapter 34 — Logs: `journalctl` and where output really goes

### In one sentence

On a systemd-based system, most service output is captured centrally by the
**journal**, queryable with `journalctl` — but plenty of software still
writes plain text log files under `/var/log`, and knowing which is which (and
keeping either from filling your disk) is a real, constant operational skill.

### Why this matters

A disk that fills up with unrotated logs is one of the single most common,
most preventable real production outages in the industry's entire history —
and it's caused by a genuinely simple oversight.

### How it actually works

```bash
   journalctl                    # ALL journal entries, oldest first
   journalctl -u myapp            # just THIS service's log entries
   journalctl -f                   # FOLLOW -- live tail, like `tail -f`
   journalctl -b                    # just THIS boot session
   journalctl --since "1 hour ago"   # time-filtered
   journalctl -p err                # only ERROR-level and above

   # Traditional plain-text logs still live under /var/log -- not
   # everything goes through the journal (Chapter 24's FHS: /var is
   # "variable data," logs are the classic example)
   tail -f /var/log/syslog          # (Debian/Ubuntu) or
   tail -f /var/log/messages         # (RHEL/Fedora family) -- the
                                    # general system log
   ls /var/log/                      # see what else is being logged --
                                    # often includes per-application
                                    # logs, auth logs, and more
```

**Log rotation** (`logrotate`, the standard tool) exists specifically to
prevent logs from growing forever:

```
   /etc/logrotate.conf and /etc/logrotate.d/*  define, per log file/
   category: how often to rotate (daily/weekly), how many old copies to
   KEEP, whether to compress old copies, and when to DELETE the oldest
   ones entirely. Without this running correctly, a busy service's logs
   genuinely can, and regularly do, grow without bound.
```

### Real-world example

```
   "The disk filled up with log files and took down production" is one
   of the most common, most WIDELY DOCUMENTED real operational
   incidents across the entire software industry -- a search of almost
   any large company's public postmortem archive turns up at least one
   example, and it's a recurring theme in SRE training material
   precisely BECAUSE it keeps happening, in new forms, across different
   organisations, for decades. The shape is almost always the same: a
   service (often a NEW one, or one recently changed to log more
   verbosely for debugging) writes logs that were never added to
   `logrotate`'s configuration, or a misconfigured log rotation rule
   that silently stops working, or a runaway ERROR condition that
   causes MASSIVE log volume in a short period (an application stuck in
   a tight retry-and-log loop is a classic culprit) -- and because
   Linux behaves VERY badly when its root filesystem runs completely
   out of space (new file writes fail, and depending on what else lives
   on that filesystem, entirely unrelated services can start failing
   too), what began as "just a logging problem" becomes a full outage.
   THE PREVENTABLE FIX, universally recommended across incident
   postmortems on this exact topic: EVERY log-producing service gets an
   explicit logrotate configuration from day one (not added later, after
   the first incident), PLUS disk-usage monitoring with alerting well
   BEFORE the disk is actually full (Chapter 41 covers monitoring disk
   usage properly) -- treating log rotation as a deployment checklist
   item, not an afterthought.
```

### Try it yourself

```bash
# Explore your own system's logs
journalctl -u ssh --no-pager 2>/dev/null | tail -10 || \
  journalctl -u sshd --no-pager 2>/dev/null | tail -10

journalctl --since "10 minutes ago" --no-pager | tail -20

ls -la /var/log/ | head -15
du -sh /var/log/* 2>/dev/null | sort -rh | head -10   # biggest log
                                                       # consumers right
                                                       # now

# See your log rotation configuration
cat /etc/logrotate.conf 2>/dev/null | head -20
ls /etc/logrotate.d/ 2>/dev/null

# See how much disk space is currently used, and by what
df -h /
du -sh /var/log 2>/dev/null
```

### Common mistakes

- **Deploying a new, verbose service without checking it's covered by log
  rotation.** The single most common preventable cause of "disk filled with
  logs" incidents.
- **Only monitoring "is the service up," not "is the disk filling up."** A
  disk-full incident is often entirely predictable hours in advance if
  disk-usage trends are actually being watched.
- **Assuming `journalctl` captures EVERYTHING.** Plenty of software still
  writes its own plain-text log files directly to `/var/log`, entirely
  outside the journal — check both when troubleshooting.

### Check yourself

1. What's the difference between the systemd journal and traditional
   `/var/log` files?
2. What does logrotate actually do, and why does it exist?
3. Why does "the disk filled up with logs" so often cascade into a much
   bigger outage than just "logging stopped working"?

### Further reading

- **Reference:** `man 1 journalctl`, `man 8 logrotate`, `man 5 logrotate.conf`.
- **Docs:** `freedesktop.org`'s systemd journal documentation.
- **Article:** search your own company's (or any major tech company's)
  public postmortem archive for "disk full" incidents — genuinely
  instructive, real examples abound.

---

# Part 10 — Users, permissions, and basic hardening

Every command in Parts 7–9 ran as some user with some set of permissions,
mostly unexamined. This part makes that explicit: who's allowed to do what,
how `sudo` actually works underneath the prompt, and the baseline hardening
steps that follow directly once permissions are set up correctly.

## Chapter 35 — Users, groups, and `sudo`

### In one sentence

Every process runs AS a specific user (Chapter 4's permission field), Linux
groups let you grant the same access to many users at once, and `sudo` is the
standard, AUDITABLE way to run something as a different user (almost always
root) without actually logging in as them.

### Why this matters

`sudo` is one of the most-run commands on any real Linux system, and it's
also — because it's a gateway to root — a genuinely high-stakes piece of
software, as a real, severe vulnerability made clear.

### How it actually works

```bash
   whoami                    # who am I, right now
   id                         # my user ID, group ID, and all groups
                             # I belong to
   sudo command                # run ONE command as root (by default),
                             # after authenticating (usually your OWN
                             # password, not root's)
   sudo -u alice command        # run as a SPECIFIC other user, not
                             # necessarily root
   su - alice                   # actually SWITCH USER, starting a new
                             # login shell as alice (needs alice's
                             # password, unless you're already root)

   sudo -l                       # list what YOU are allowed to run with
                             # sudo -- often surprisingly restrictive,
                             # by design, on a well-configured system

   /etc/passwd                   # every user account, with basic info
   /etc/group                     # every group, and its members
   /etc/sudoers                    # who can sudo, and to run WHAT
                             # (edit with `visudo`, never directly --
                             # it validates syntax before saving,
                             # preventing you from locking yourself out
                             # with a typo)
```

**Why `sudo` is preferred over logging in as root directly**, as a real,
deliberate security design:

```
   - EVERY sudo invocation is LOGGED (who, what command, when) --
     accountability that a shared root login simply doesn't give you.
   - Access can be SCOPED: a user can be granted sudo rights to run only
     SPECIFIC commands, not unrestricted root access.
   - It requires the USER'S OWN credentials (by default), not a shared
     root password everyone has to know and rotate together.
```

### Real-world example

```
   In January 2021, the security firm Qualys publicly disclosed
   CVE-2021-3156, nicknamed "Baron Samedit" -- a HEAP-BASED BUFFER
   OVERFLOW vulnerability in `sudo` itself, present in the codebase for
   roughly a DECADE before it was found, affecting nearly every major
   Linux distribution's default `sudo` installation. The bug allowed
   ANY local, unprivileged user (no sudo permissions needed at all) to
   exploit a flaw in how `sudo` parsed command-line arguments (related
   to escaping backslash characters when invoked in a specific mode)
   to gain FULL ROOT ACCESS, with no password required.
   Because `sudo` is installed by default on the vast majority of
   real-world Linux systems -- servers, desktops, containers built on
   common base images -- this was an exceptionally widespread, severe
   vulnerability, and it triggered urgent, coordinated patching across
   the entire industry within days of disclosure. It's a genuinely
   sobering, well-documented real illustration of two things this
   Part keeps returning to: first, that the tools we trust MOST
   completely (a program whose entire purpose is granting root access)
   are exactly the ones where a bug is most consequential; and second,
   that "this code has existed and been widely used for ten years" is
   NOT the same as "this code has been proven safe" -- it can simply
   mean the specific bug hadn't been found yet.
```

### Try it yourself

```bash
# Explore your own user/group setup
id
groups
cat /etc/passwd | head -5
cat /etc/group | grep $(whoami)

# See what YOU'RE allowed to sudo (varies by system configuration)
sudo -l 2>/dev/null | head -20

# Check your sudo version against the CVE-2021-3156 fix (informational)
sudo --version | head -1
# (any reasonably current, patched Linux distro will already be safe --
#  this is just to see the version-checking habit, not to alarm you)
```

### Common mistakes

- **Sharing a root login/password among a team instead of using `sudo`.**
  Loses individual accountability entirely — nobody can tell WHO did what.
- **Editing `/etc/sudoers` directly with a text editor.** A syntax error can
  lock EVERYONE, including you, out of `sudo` — always use `visudo`, which
  validates before saving.
- **Assuming widely-used, long-standing software is inherently safe from
  serious bugs.** The Baron Samedit vulnerability sat undiscovered in `sudo`
  for roughly ten years — keep systems patched, always, regardless of how
  "trusted" a tool is.

### Check yourself

1. Why is `sudo` generally preferred over sharing a root login?
2. What was CVE-2021-3156, and why was it so significant?
3. Why should you always use `visudo` instead of editing `/etc/sudoers`
   directly?

### Further reading

- **Reference:** `man 8 sudo`, `man 8 visudo`, `man 5 sudoers`.
- **Advisory:** Qualys's original "Baron Samedit" technical writeup
  (search "CVE-2021-3156 Qualys") — detailed, well-explained, worth reading
  in full.
- **Practice:** run `sudo -l` on a real system you administer and confirm
  the scope actually matches what you intended to grant.

---

## Chapter 36 — The permission model, beyond `chmod 777`

### In one sentence

Standard owner/group/other permissions (Chapter 25) are a blunt instrument —
Linux CAPABILITIES let you grant a program a specific, narrow slice of root's
power instead of all of it, and ACLs let you grant fine-grained permissions to
specific individual users beyond just "owner" and "group."

### Why this matters

"This program needs a LITTLE bit of extra privilege, so I made it full SUID
root" is a real, common, and often unnecessary security trade-off — modern
Linux gives you much narrower tools for exactly this problem.

### How it actually works

```
   LINUX CAPABILITIES split "what root can do" into roughly 40
   independent, individually-grantable PRIVILEGES, instead of the
   old all-or-nothing model (either you're root, with EVERY privilege,
   or you're not, with NONE of them). A few examples:

     CAP_NET_RAW        create raw network sockets (needed by `ping`,
                        packet capture tools)
     CAP_NET_BIND_SERVICE  bind to a port below 1024 (historically
                        root-only)
     CAP_SYS_TIME         change the system clock
     CAP_CHOWN             change file ownership
     CAP_SYS_ADMIN          a large, famously broad grab-bag of
                        administrative operations (genuinely close to
                        "basically root" on its own -- treated with
                        real caution)

   You can grant a SPECIFIC capability to a specific binary, instead of
   making it SUID-root (Chapter 25):

     sudo setcap cap_net_raw+ep /path/to/mytool
     getcap /path/to/mytool          # confirm what was granted

   This means a program with a BUG can, at worst, misuse the ONE
   narrow capability it was granted -- not read every file on the
   system, not create arbitrary users, not do anything else root could
   do. This same mechanism (Linux capabilities) is exactly what
   container runtimes use to run containers with a much smaller,
   scoped set of privileges than genuine root, rather than either "full
   root" or "no privileges at all" (Phase 3 covers this properly).

   ACLs (Access Control Lists) extend the basic owner/group/other model
   to allow permissions for SPECIFIC, ADDITIONAL individual users or
   groups on a single file -- something the basic three-category model
   simply cannot express (e.g. "the owner and the finance group can
   both write to this file, but ALSO grant read-only to this one
   specific contractor's account, without adding them to the finance
   group at all"):

     setfacl -m u:bob:r-- shared_report.csv   # grant bob read-only,
                                              # specifically
     getfacl shared_report.csv                 # see all ACL entries
```

### Real-world example

```
   The `ping` command is a genuinely useful, concrete real illustration
   of this chapter's improvement over the blunt SUID model (Chapter
   25): `ping` NEEDS to create a raw network socket to send ICMP
   packets, an operation that historically required root privileges
   entirely -- which is exactly why, for decades, `/bin/ping` on most
   Unix-like systems was SUID-root: an ordinary user could run it, and
   it would run WITH root's full power, purely to get the one specific
   permission (creating a raw socket) it actually needed.
   On modern Linux systems, this is commonly implemented instead using
   the CAP_NET_RAW capability specifically (`getcap /bin/ping` on many
   current distributions shows exactly this, rather than the SUID bit)
   -- `ping` gets precisely the one privilege it actually requires, and
   NOTHING else. If a bug were ever found in `ping`'s code (and
   real vulnerabilities HAVE been found in ping implementations over
   the years, across various Unix systems, historically), the
   worst-case blast radius with capabilities is "misuse of raw socket
   creation" -- not "arbitrary code execution as full root," which is
   what a SUID-root `ping` bug would have granted. This is a small,
   concrete, verifiable-on-your-own-machine example of the exact
   principle -- least privilege, applied precisely instead of broadly --
   that runs through this entire guide and Phase 3's container security
   material.
```

### Try it yourself

```bash
# See ping's actual privilege model on your system
which ping
ls -l $(which ping)              # look for the SUID bit (an 's' in the
                                 # owner-execute position) -- may or may
                                 # not be present, depending on your distro
getcap $(which ping) 2>/dev/null  # or see the capability-based approach
                                  # -- many modern distros use this instead

# Try capabilities yourself, safely, on a harmless test binary
cat > /tmp/captest.c << 'EOF'
#include <stdio.h>
int main() { printf("hello from a capability-scoped binary\n"); return 0; }
EOF
gcc -o /tmp/captest /tmp/captest.c
getcap /tmp/captest                 # nothing granted yet
sudo setcap cap_net_raw+ep /tmp/captest
getcap /tmp/captest                 # now shows the granted capability
rm /tmp/captest /tmp/captest.c

# See ACLs in action
touch shared_report.csv
setfacl -m u:$(whoami):r-- shared_report.csv 2>/dev/null || \
  echo "(setfacl not installed -- try: sudo apt install acl)"
getfacl shared_report.csv 2>/dev/null
rm -f shared_report.csv
```

### Common mistakes

- **Reaching for SUID-root as the default fix for "this needs extra
  privilege."** Check whether a specific capability covers exactly what's
  needed first — it's a strictly smaller, safer grant.
- **Granting `CAP_SYS_ADMIN` casually.** It's broad enough to be close to
  "full root" in practice — treat it with the same caution as SUID-root
  itself, not as "just another capability."
- **Forgetting ACLs exist and fighting the owner/group/other model to
  express something it can't.** If you find yourself creating a new GROUP
  just to grant one person access to one file, an ACL is usually the more
  precise tool.

### Check yourself

1. What problem do Linux capabilities solve that SUID doesn't?
2. What does `CAP_NET_RAW` grant, and why does `ping` need it (or its
   SUID-root equivalent)?
3. What can an ACL express that basic owner/group/other permissions cannot?

### Further reading

- **Reference:** `man 7 capabilities`, `man 1 setcap`, `man 1 getcap`,
  `man 1 setfacl`, `man 5 acl`.
- **Guide:** this wiki's security guides cover container capability-dropping
  (`--cap-drop ALL`, adding back only what's needed) in operational depth —
  the direct Phase 3 continuation of this chapter.
- **Article:** search "linux capabilities vs suid" for clear comparative
  writeups.

---

# Part 11 — Networking on Linux

Everything until now happened on one machine. This part is how that machine
talks to others — the networking stack from a single Linux box's point of
view. It complements rather than repeats this Wiki's dedicated TCP/IP
guide; the focus here is what you configure and inspect on the host itself.

## Chapter 37 — How a Linux box sees the network

### In one sentence

A Linux machine's view of "the network" is built from a handful of concrete
kernel concepts — network interfaces, IP addresses, routes, and DNS
resolution — and every one of them can be inspected directly, which is
exactly how you'll debug connectivity problems for the rest of your career.

### Why this matters

"I can't reach the server" has a specific, findable cause every time — DNS,
routing, a closed port, a firewall rule — and knowing this layer means you
check the right thing first instead of guessing.

### How it actually works

```bash
   ip addr show                # every network INTERFACE and its
                               # assigned IP address(es)
   ip route show                 # the ROUTING TABLE -- where does
                               # traffic to a given destination
                               # actually go?
   cat /etc/resolv.conf           # which DNS server(s) this box asks

   THE CONCEPTS, briefly:
     INTERFACE     a network endpoint the kernel manages -- a physical
                  NIC (`eth0`, `enp0s3`), or a VIRTUAL one (the loopback
                  `lo`, a VPN tunnel, a virtual bridge -- Phase 3's
                  container networking is built entirely on virtual
                  interfaces).
     IP ADDRESS     a numeric address assigned to an interface -- how
                  packets know where to go.
     ROUTE           a rule: "to reach THIS range of addresses, send
                  traffic out THROUGH this interface, via this
                  gateway." Every machine has a DEFAULT route -- where
                  everything NOT otherwise covered goes (usually your
                  router).
     DNS RESOLUTION  translating a hostname ("example.com") into an IP
                  address, by asking a configured DNS server (or
                  checking `/etc/hosts` first, for manual overrides).

   THE LOOPBACK INTERFACE (`lo`, almost always IP address 127.0.0.1,
   "localhost") is worth understanding precisely: it's a VIRTUAL
   interface that loops traffic straight back to the SAME machine,
   entirely within the kernel, never touching real network hardware at
   all. This matters a great deal once you reach containers (Phase 3):
   EACH container typically gets its OWN, separate loopback interface
   and its own "localhost," entirely distinct from the host's -- which
   is exactly why a service listening on `localhost` INSIDE a container
   is invisible from OUTSIDE it, a genuinely common, real point of
   confusion for engineers new to Docker.
```

### Real-world example

```
   NETWORK NAMESPACES (a Linux kernel feature covered properly in Phase
   3, but worth previewing here) give EACH container its own, separate
   copy of everything this chapter just described: its own loopback
   interface, its own IP address(es), its own routing table -- fully
   isolated from the host's, by default. This is the DIRECT, concrete
   mechanism behind an extremely common real question every engineer
   asks the first time they use Docker: "I started a web server on
   `localhost:8080` inside my container, why can't I reach it from my
   browser on the host?"
   The answer is exactly this chapter's model, applied literally: the
   container's `localhost` is its OWN loopback interface, entirely
   separate from the host machine's `localhost` -- a service bound to
   it is reachable ONLY from processes inside that SAME network
   namespace (i.e., other processes in the same container), not from
   the host or any other container, until you EXPLICITLY publish/map
   a port (`docker run -p 8080:8080 ...`, covered in Chapter 50) to
   bridge the two separate network views together. Understanding
   Chapter 37's plain concepts -- interfaces, addresses, routes -- as
   REAL, INSPECTABLE kernel objects (not abstract networking theory) is
   precisely what makes this "obvious" instead of mysterious once you
   reach Phase 3.
```

### Try it yourself

```bash
# Inspect your own machine's network view
ip addr show
ip route show
cat /etc/resolv.conf

# See the loopback interface specifically
ip addr show lo
ping -c 3 127.0.0.1     # loops entirely within the kernel, no real
                        # network hardware involved at all

# Resolve a hostname manually, and see /etc/hosts take priority
getent hosts example.com    # uses the SAME resolution order your
                            # system actually uses (checks /etc/hosts
                            # first, then DNS)
cat /etc/hosts               # your manual overrides, if any
```

### Common mistakes

- **Assuming `127.0.0.1` is the SAME thing on every machine/container.**
  It's a loopback address, local to WHATEVER network namespace you're
  currently in — a container's `127.0.0.1` and the host's are different,
  isolated interfaces (Phase 3).
- **Troubleshooting "can't connect" by guessing instead of checking the
  layers in order.** DNS resolving correctly? Route existing to the
  destination? Port actually open there (Chapter 38)? Firewall blocking it?
  Check in that order, not randomly.
- **Editing `/etc/resolv.conf` by hand on a modern system and expecting it
  to stick.** Many modern distros manage this file automatically (via
  NetworkManager or systemd-resolved) and will overwrite manual edits —
  check your distro's actual DNS-management approach first.

### Check yourself

1. What does a routing table entry actually describe?
2. Why is the loopback interface's behaviour different inside a container
   than on the host?
3. What's the standard order of things to check when troubleshooting "I
   can't reach this service"?

### Further reading

- **Reference:** `man 8 ip`, `man 5 resolv.conf`, `man 5 hosts`.
- **Book (free):** *TCP/IP Illustrated* concepts are covered at a deeper,
  protocol level in this wiki's own networking guide — a natural next step.
- **Guide:** this wiki's `networking/tcp-ip` guide covers everything from DNS
  to TLS in far greater depth than this chapter attempts.

---

## Chapter 38 — The tools: `ss`, `curl`, `dig`, `tcpdump`, and a firewall's basics

### In one sentence

A small set of command-line tools — `ss` for local connections, `curl` for
making requests, `dig` for DNS, `tcpdump` for watching raw traffic, and
basic firewall commands — cover the overwhelming majority of real network
troubleshooting you'll ever do.

### Why this matters

These are the exact tools referenced constantly across incident postmortems
and on-call runbooks industry-wide — fluency here is a direct, practical
career skill, not academic knowledge.

### How it actually works

```bash
   ss -tulpn              # every LISTENING socket on this machine:
                         # protocol (t=tcp, u=udp), listening state,
                         # PID/process name. THE modern replacement for
                         # the older `netstat` (see below).
   ss -tn state established    # currently ESTABLISHED TCP connections

   curl -v https://example.com    # make an HTTP(S) request, -v shows
                                 # the full request/response, headers
                                 # and all
   curl -o file.txt <url>          # save output to a file instead of
                                 # printing it
   curl -I <url>                    # HEAD request only -- just headers,
                                 # fast way to check if something's
                                 # reachable and what it returns

   dig example.com                   # full DNS lookup, with all the
                                 # detail
   dig +short example.com             # just the answer -- the IP(s)
   dig example.com MX                  # a SPECIFIC record type (mail
                                 # servers, in this case)

   tcpdump -i any port 443              # watch raw packets live,
                                 # filtered to a specific port
                                 # (needs root/sudo)

   # basic firewall inspection (varies by distro/tooling -- ufw is
   # common on Ubuntu, firewalld on RHEL/Fedora, both are friendly
   # frontends on top of the kernel's netfilter/iptables/nftables)
   sudo ufw status               # Ubuntu-style
   sudo firewall-cmd --list-all   # RHEL/Fedora-style
```

### Real-world example

```
   `ss` replacing `netstat` is a genuinely well-documented, real tooling
   transition worth knowing the story of, because it explains why some
   tutorials you'll find online teach a command that's no longer the
   recommended default: `netstat` (along with `ifconfig`, `route`, and
   others) came from the older `net-tools` package, which -- as
   documented in its own project history and repeatedly discussed across
   Linux distribution mailing lists and release notes -- went through a
   LONG period with little to no active maintenance. `iproute2` (which
   provides `ss`, `ip`, and related modern commands) became the actively
   maintained, kernel-feature-current replacement, and most major
   distributions eventually stopped installing `net-tools` BY DEFAULT,
   years apart from each other, causing exactly the confusion you'd
   expect: countless still-circulating tutorials and Stack Overflow
   answers reference `netstat`/`ifconfig` commands that may not even be
   INSTALLED on a fresh, modern system.
   THE PRACTICAL LESSON: `ss` and `ip` are the current, actively-
   maintained standard (`ss` for socket/connection info, `ip` for
   interfaces/routes/addresses -- Chapter 37's commands) -- worth learning
   as your defaults, while still recognising `netstat`/`ifconfig` syntax
   when you encounter them in older material, exactly the way you'd
   recognise an older dialect of a language you're fluent in.
```

### Try it yourself

```bash
# See what's actually listening on your machine, and by what
sudo ss -tulpn | head -15

# Make some real requests and inspect them
curl -I https://example.com
curl -v https://example.com 2>&1 | head -20

# Real DNS lookups
dig +short example.com
dig example.com MX +short

# Watch real traffic (needs sudo; Ctrl-C to stop; generates traffic
# with curl in a second terminal, or just watch DNS lookups happen)
sudo timeout 5 tcpdump -i any port 53 -n 2>/dev/null &
dig example.com > /dev/null
wait

# Check your firewall's current state (read-only, safe)
sudo ufw status 2>/dev/null || sudo firewall-cmd --list-all 2>/dev/null || \
  sudo iptables -L -n 2>/dev/null | head -15
```

### Common mistakes

- **Following a tutorial that uses `netstat`/`ifconfig` and being confused
  they're missing.** Recognise the older syntax, but reach for `ss`/`ip` as
  your actual default on modern systems.
- **Using `curl` without `-v`/`-I` when troubleshooting, and only seeing the
  body.** The headers and status code are very often where the actual
  problem/answer is.
- **Assuming a firewall's default-deny stance is "broken" rather than
  intentional.** A properly hardened box denies everything not explicitly
  allowed — a connection failing might be the firewall doing exactly its job,
  not a bug.

### Check yourself

1. What replaced `netstat`, and why?
2. What does `dig +short` give you that plain `dig` doesn't, and when would
   you want the fuller output instead?
3. Given "I can't reach this web service," name the check-in-order sequence
   from Chapter 37 combined with this chapter's tools.

### Further reading

- **Reference:** `man 8 ss`, `man 1 curl`, `man 1 dig`, `man 8 tcpdump`.
- **Guide:** this wiki's `networking/tcp-ip` guide is the deep, dedicated
  reference for everything beyond this chapter's practical toolkit.
- **History:** search "net-tools deprecated iproute2" for the real story
  behind the tooling transition.

---

# Part 12 — Shell scripting and automation

You've now run dozens of commands by hand across Parts 7–11. This part is
about not doing that by hand anymore — enough bash to automate the
repetitive parts of everything covered so far, from file operations to
process management to basic networking checks.

## Chapter 39 — Bash scripting fundamentals

### In one sentence

A bash script is just a sequence of the same commands you'd type
interactively, saved to a file and given a few extra safety rails — and those
safety rails (`set -e`, quoting, checking exit status) are what separates a
script you can trust in production from one that fails silently.

### Why this matters

An unattended script that keeps running after something has already gone
wrong is one of the most dangerous, most common real sources of production
damage — automation multiplies whatever mistake you make, at machine speed.

### How it actually works

```bash
   #!/bin/bash              # the SHEBANG -- tells the OS which
                            # interpreter to run this file with
                            # (Chapter 6's exec(), essentially)
   set -e                    # EXIT IMMEDIATELY if any command fails
                            # (non-zero exit status, Chapter 3) --
                            # without this, a script just KEEPS GOING
                            # after a failure, often compounding it
   set -u                    # error out on any UNSET variable, instead
                            # of silently treating it as empty (this is
                            # EXACTLY the class of bug behind Chapter
                            # 26's Steam `rm -rf` incident)
   set -o pipefail            # a pipeline's exit status reflects the
                            # FIRST failing command, not just the last
                            # one -- without this, `false | true` reports
                            # SUCCESS, hiding a real failure upstream

   # a common, recommended combination:
   set -euo pipefail

   VARIABLE="value"             # no spaces around =, by convention
   echo "$VARIABLE"              # ALWAYS quote variable expansions --
                                 # unquoted expansion is a classic
                                 # source of bugs with spaces/globs

   if [ "$1" = "start" ]; then    # test/condition syntax
     echo "starting"
   elif [ "$1" = "stop" ]; then
     echo "stopping"
   else
     echo "usage: $0 start|stop" >&2
     exit 1
   fi

   for file in *.txt; do          # loop over files
     echo "processing $file"
   done

   result=$(some_command)          # capture command output into a
                                  # variable
   if [ $? -ne 0 ]; then           # $? = the PREVIOUS command's exit
                                  # status (Chapter 3) -- but prefer
                                  # `set -e` over manually checking
                                  # this everywhere
     echo "failed" >&2
     exit 1
   fi
```

### Real-world example

```
   Knight Capital Group's 2012 trading incident remains one of the
   most widely studied, most consequential real illustrations of why
   AUTOMATION AND DEPLOYMENT DISCIPLINE matters -- publicly documented
   via the SEC's own investigative order (a real, official regulatory
   document, not just press reporting). During a software deployment,
   NEW trading code was rolled out to SEVEN production servers, but --
   due to a deployment process failure -- ONE of the eight servers
   retained OLD, previously-repurposed code that reused a flag value
   with a different, dangerous meaning under the new system. When live
   trading began, that ONE server's stale code triggered a RUNAWAY
   sequence of unintended trades. Automated systems, running exactly as
   designed, executed the resulting bad instructions at high speed with
   NO human in the loop fast enough to intervene -- Knight Capital lost
   roughly $440 MILLION in about 45 minutes, and the firm itself was
   sold off shortly afterward as a direct result.
   THE LESSON THIS CHAPTER TAKES FROM IT, stated precisely: the failure
   wasn't one bad line of code in isolation -- it was a DEPLOYMENT
   PROCESS with no verification step confirming all servers actually
   ran the SAME, INTENDED code before going live, combined with
   automation that had no circuit breaker to catch and halt an
   obviously-anomalous pattern of behaviour once it started. `set -e`
   catching a script failure immediately, verifying deployment state
   before proceeding, and building in sanity checks BEFORE automation
   acts irreversibly are all, in spirit, the exact same discipline this
   incident demonstrates the cost of skipping -- at a scale ordinary
   engineers will rarely if ever match, but the underlying principle
   scales down to every cron job and deploy script you'll ever write.
```

### Try it yourself

```bash
cd ~/lab/playground

# See set -e catch a failure that would otherwise continue silently
cat > unsafe.sh << 'EOF'
#!/bin/bash
echo "step 1"
false                      # this "fails" but the script keeps going
echo "step 2 (ran anyway, even though step 1 failed!)"
EOF
chmod +x unsafe.sh
./unsafe.sh

cat > safe.sh << 'EOF'
#!/bin/bash
set -euo pipefail
echo "step 1"
false                      # this time, the script STOPS here
echo "step 2 (this line should never print)"
EOF
chmod +x safe.sh
./safe.sh; echo "exit status: $?"

# See set -u catch an unset variable (the Steam-bug class of mistake)
cat > unset_demo.sh << 'EOF'
#!/bin/bash
set -u
echo "TARGET is: ${TARGET}"
EOF
chmod +x unset_demo.sh
./unset_demo.sh    # errors out immediately, loudly -- exactly what you
                   # want, instead of silently treating TARGET as empty

rm unsafe.sh safe.sh unset_demo.sh
```

### Common mistakes

- **Writing a script without `set -euo pipefail`.** The single highest-
  leverage habit in this whole chapter — add it to every script, as a
  default, before you write anything else.
- **Leaving variable expansions unquoted (`$VAR` instead of `"$VAR"`).** A
  value containing spaces or shell-special characters can silently break the
  script's logic in ways that are hard to spot by eye.
- **Trusting a script "worked" because it didn't print an error.** Without
  `set -e`, a failed command's error can scroll past unnoticed while later
  steps run anyway, on bad assumptions.

### Check yourself

1. What does `set -e` do, and what's the risk of a script WITHOUT it?
2. Why does `set -u` specifically address the class of bug behind Chapter
   26's Steam `rm -rf` incident?
3. What was the actual root cause of the Knight Capital incident — not just
   "a bug," but specifically?

### Further reading

- **Book (free):** *The Linux Command Line* (Shotts) — the shell scripting
  chapters, thorough and example-driven.
- **Reference:** `man 1 bash` (the "SHELL BUILTIN COMMANDS," `set` section).
- **Tool:** `shellcheck` (`shellcheck.net`, or install locally) — a static
  analysis tool for shell scripts; run it on everything you write.
- **Filing:** the SEC's own administrative proceeding release on Knight
  Capital (search "SEC Knight Capital release") — a genuinely detailed,
  real regulatory account of exactly what went wrong.

---

## Chapter 40 — Cron, environment variables, and gluing things together

### In one sentence

`cron` runs a command on a schedule, completely unattended — which means it
also runs with a much SMALLER environment than your interactive shell,
silently, and that gap is the single most common reason "it works when I run
it myself" and "it fails under cron" are both true at once.

### Why this matters

A silently-failing scheduled job (a backup, a cleanup task, a report) can go
unnoticed for weeks — genuinely common, and genuinely costly when the job
turns out to have mattered.

### How it actually works

```bash
   crontab -e                   # edit YOUR OWN scheduled jobs
   crontab -l                    # list them

   # format: minute hour day-of-month month day-of-week   command
   0 2 * * *     /path/to/backup.sh      # every day at 2:00 AM
   */15 * * * *  /path/to/healthcheck.sh  # every 15 minutes
   0 0 1 * *     /path/to/monthly.sh       # midnight, first of every
                                          # month

   ENVIRONMENT VARIABLES are named values available to a process and
   everything it launches -- `$PATH` (where to look for commands),
   `$HOME`, `$USER`, and anything else you or a program has set:

     export MY_VAR="value"        # set it, and make it available to
                                 # any CHILD process this shell launches
     echo $MY_VAR                  # read it
     env                            # list every variable in the CURRENT
                                 # environment
     printenv PATH                  # just one specific variable

   WHY CRON JOBS BEHAVE DIFFERENTLY: cron runs your job with a MINIMAL,
   often DIFFERENT environment than your interactive login shell --
   frequently a much SHORTER `$PATH`, no shell customisations from your
   `.bashrc`, and none of the environment variables your normal
   terminal session has accumulated. A script that calls a command by
   NAME alone (relying on `$PATH` to find it) can work perfectly when
   YOU run it, and fail silently (or with a confusing "command not
   found") under cron, purely because cron's `$PATH` doesn't include
   wherever that command actually lives.
```

### Real-world example

```
   "The cron job has been silently failing for weeks and nobody
   noticed" is one of the most universally-experienced real operational
   patterns in the industry -- documented across countless sysadmin
   forums, blog postmortems, and "lessons learned" writeups, precisely
   because the FAILURE MODE ITSELF is silent by design: cron, by
   default, only EMAILS output if the system's mail delivery is
   actually configured (which, on a huge number of real servers, it
   simply isn't) -- so a script that starts failing, whether from an
   environment mismatch (this chapter's exact `$PATH` gotcha), a
   changed file location, an expired credential, or any other cause,
   can fail EVERY SINGLE RUN, indefinitely, with the job's log showing
   nothing but a string of quiet non-zero exit codes nobody is watching.
   THE STANDARD, WIDELY-RECOMMENDED REAL FIXES, all addressing the same
   underlying "silent by default" problem:
     - set `PATH` EXPLICITLY at the top of the crontab, or use full,
       absolute paths for every command inside the script itself,
       rather than relying on cron's minimal default `$PATH`.
     - redirect the job's output EXPLICITLY to a log file
       (`>> /var/log/myjob.log 2>&1`) so there's always somewhere to
       actually look.
     - add MONITORING that alerts if the job HASN'T reported success
       recently (a "dead man's switch" pattern -- if a heartbeat/success
       signal stops arriving, THAT absence is the alert), rather than
       relying on the job itself to notice and report its own failure,
       which is exactly the assumption that fails when the job is
       broken badly enough to not run at all.
```

### Try it yourself

```bash
# See the environment difference directly -- this is the whole lesson
echo "interactive PATH: $PATH"

# Simulate cron's minimal environment and see a script behave
# differently
env -i /bin/bash -c 'echo "minimal PATH: $PATH"; which curl || echo "curl not found!"'

# Write a small, properly defensive cron-style script
cat > ~/lab/playground/cronjob.sh << 'EOF'
#!/bin/bash
set -euo pipefail
LOG="/tmp/cronjob.log"
echo "$(date): job started" >> "$LOG"
/usr/bin/echo "doing the actual work" >> "$LOG"   # absolute path --
                                                   # doesn't rely on
                                                   # cron's $PATH at all
echo "$(date): job finished successfully" >> "$LOG"
EOF
chmod +x ~/lab/playground/cronjob.sh

# Install it to run every minute, briefly, to see it work (remove
# afterward)
(crontab -l 2>/dev/null; echo "* * * * * ~/lab/playground/cronjob.sh") | crontab -
sleep 65
cat /tmp/cronjob.log
crontab -l | grep -v cronjob.sh | crontab -    # remove the test entry
rm /tmp/cronjob.log ~/lab/playground/cronjob.sh
```

### Common mistakes

- **Calling commands by name in a cron script, relying on `$PATH`.** Use
  absolute paths, or explicitly set `PATH` at the top of the crontab.
- **Assuming "no error email" means "it's working."** Mail delivery often
  isn't configured at all — redirect output to a log file explicitly, and
  monitor for the job's ABSENCE of success, not just its presence of error.
- **Never testing a cron job's ACTUAL cron execution, only running the
  script manually.** The environment difference means "works when I run it"
  is genuinely weak evidence it'll work under cron — test it for real.

### Check yourself

1. Why can a script work perfectly when run manually but fail under cron?
2. What are two concrete, standard fixes for the "cron job silently fails"
   problem?
3. What is a "dead man's switch" monitoring pattern, and why does it fit
   this specific failure mode better than "alert on error output"?

### Further reading

- **Reference:** `man 5 crontab`, `man 1 crontab`, `man 1 env`.
- **Tool:** `crontab.guru` — an interactive tool for reading/writing cron
  schedule expressions, genuinely useful even for experienced engineers.
- **Article:** search "cron job silent failure best practices" for current,
  practical monitoring-pattern writeups.

---

# Part 13 — Performance and troubleshooting

This part pulls together CPU (Part 2), memory (Part 3), and I/O (Part 4)
into the actual skill you've been building toward: reading a struggling
system's vital signs correctly and knowing which of the three is really the
bottleneck, instead of guessing.

## Chapter 41 — Reading load average, `vmstat`, and `iostat` correctly

### In one sentence

"Load average" is one of the most widely MISREAD numbers in all of Linux
operations — on Linux specifically (unlike some other Unix systems), it
counts processes waiting on disk I/O, not just processes waiting for CPU, and
knowing that single fact prevents a huge fraction of false diagnoses.

### Why this matters

Misreading load average as "pure CPU demand" sends engineers chasing CPU
optimisation for problems that are actually disk or network I/O bottlenecks —
a real, common, and avoidable waste of incident-response time.

### How it actually works

```bash
   uptime                    # shows load average as three numbers:
                             # 1-minute, 5-minute, 15-minute averages
   # load average: 2.15, 1.98, 1.75

   WHAT THOSE NUMBERS MEAN: roughly, the average number of processes
   that were EITHER actively running on a CPU OR WAITING to (Chapter 5's
   "R" state) OR -- and this is the widely-misunderstood part, ON LINUX
   SPECIFICALLY -- in UNINTERRUPTIBLE SLEEP (Chapter 5's "D" state,
   almost always waiting on disk/network I/O), averaged over that time
   window. A load average of "4" on a 4-core machine, if it's ALL
   CPU-bound work, means the machine is fully, appropriately busy. The
   SAME number of "4," if it's mostly D-state processes stuck waiting on
   a slow, struggling disk, means something completely different --
   your CPUs might be almost entirely IDLE while load average looks
   identically "high."

   vmstat 1                    # live, per-second system stats. Key
                               # columns:
                               #   r  = processes waiting to RUN (CPU
                               #        pressure)
                               #   b  = processes BLOCKED, in
                               #        uninterruptible sleep (I/O
                               #        pressure -- Chapter 5's D state)
                               #   si/so = swap in/out (Chapter 11's
                               #        memory pressure)
                               #   us/sy/id/wa = %CPU in user time /
                               #        kernel time / idle / waiting on
                               #        I/O

   iostat -x 1                  # per-DEVICE disk I/O detail -- %util
                               # (how busy is this specific disk),
                               # await (average time per I/O request,
                               # in milliseconds -- Chapter 12's
                               # hierarchy, made concrete and
                               # measurable)
```

### Real-world example

```
   Brendan Gregg -- a widely-published, well-known performance engineer
   (formerly at Netflix and Sun/Oracle, among others, and author of
   several standard reference books on this exact topic) has written
   extensively, in detailed, technical, publicly-available material, on
   exactly this chapter's core confusion: many engineers coming from
   other Unix backgrounds (or just intuition) ASSUME "load average"
   means "CPU demand," full stop -- and Linux's kernel-level definition
   has, since a very early kernel version (a specific, documented
   change from the early 1990s), DELIBERATELY included processes
   waiting on uninterruptible I/O in the count, not just CPU-runnable
   ones. Gregg's writing on the topic (searchable, still current,
   genuinely worth reading in full) has become one of the standard
   references specifically BECAUSE this misunderstanding causes REAL,
   REPEATED, DOCUMENTED misdiagnosis in production incidents: an
   engineer sees "load average: 40" on an otherwise-idle-looking
   8-core box, assumes a CPU emergency, and starts investigating the
   wrong layer entirely -- when `vmstat`'s `b` column (or a look at
   `iostat`) would have immediately shown the REAL story: a struggling
   disk or an overwhelmed network filesystem, with the CPUs sitting
   nearly idle the entire time, waiting.
```

### Try it yourself

```bash
# Read your own machine's numbers, correctly, right now
uptime
nproc                          # how many CPU cores you actually have
                               # -- context for interpreting load average

vmstat 1 5                      # watch 'r' and 'b' columns specifically
                               # -- CPU-waiting vs I/O-waiting processes

iostat -x 1 3 2>/dev/null || sudo apt install -y sysstat  # (iostat is
                               # part of the `sysstat` package -- may
                               # need installing)

# See a CPU-bound load vs. an I/O-bound "load" side by side (safe demo)
for i in $(seq 1 4); do yes > /dev/null & done   # CPU-bound load
sleep 3; uptime; vmstat 1 2
kill $(jobs -p)                                   # clean up
```

### Common mistakes

- **Treating load average as a pure CPU metric on Linux.** It isn't — check
  `vmstat`'s `b` column or `iostat` before assuming a "high load" number
  means CPU pressure specifically.
- **Ignoring the number of CORES when judging "is this load high."** A load
  average of 8 is very different on a 2-core box (severely overloaded) than
  a 64-core one (barely used) — always read it relative to `nproc`.
- **Looking only at 1-minute load average.** The three numbers together
  (1/5/15-minute) tell you the TREND — rising, falling, or steady — which is
  often more diagnostically useful than any single snapshot.

### Check yourself

1. What does Linux's load average actually count, precisely — and how does
   that differ from "pure CPU demand"?
2. If `vmstat`'s `b` column is consistently non-zero, what kind of pressure
   does that indicate?
3. Why is "load average: 8" ambiguous without also knowing the core count?

*(Answers: Appendix E.)*

### Further reading

- **Article:** Brendan Gregg, "Linux Load Averages: Solving the Mystery" —
  search for it directly; the definitive, widely-cited explanation.
- **Book:** *Systems Performance* (2nd ed.), Brendan Gregg — the standard
  reference for this entire chapter's territory, in much greater depth.
- **Reference:** `man 8 vmstat`, `man 1 iostat`, `man 1 uptime`.

---

## Chapter 42 — `strace`, `lsof`, and answering "why is this stuck"

### In one sentence

`strace` shows you EVERY syscall a process makes, in real time, and `lsof`
shows you every file (in the broad Chapter 23 sense — files, sockets,
devices) a process has open — together, they turn "this process is doing
something mysterious" into "here's exactly what it's doing, right now."

### Why this matters

When a process is hung and you have NO idea why, these two tools are almost
always the fastest path from "no idea" to "found it" — genuinely one of the
highest-leverage skill pairs in all of Linux troubleshooting.

### How it actually works

```bash
   strace -p PID                 # attach to an ALREADY-RUNNING process
                                # and watch its syscalls live (Chapter
                                # 2) -- needs appropriate permissions
   strace command                 # run a NEW command under strace from
                                # the start
   strace -c command                # just a SUMMARY: which syscalls,
                                # how many times, how much time spent
   strace -f command                 # follow child processes too
                                # (Chapter 6's fork/exec)
   strace -e trace=open,read,write command  # filter to SPECIFIC
                                # syscalls, cutting the noise

   lsof -p PID                   # every file/socket/device THIS
                                # process has open
   lsof /path/to/file             # every process with THIS SPECIFIC
                                # file open -- "who's using this?"
   lsof -i :8080                    # what's listening on/using port
                                # 8080
   lsof +L1                         # DELETED files that are still
                                # open (Chapter 15's exact "deleted but
                                # not actually freed" scenario) --
                                # genuinely useful for finding disk
                                # space that `du` can't explain
```

**A practical diagnostic sequence, when a process seems "stuck":**

```
   1. ps -o pid,stat,cmd -p PID    check its STATE (Chapter 5) -- is
                                  it actually running (R), or stuck in
                                  D (uninterruptible I/O wait)?
   2. strace -p PID                 see the LAST syscall it made, and
                                  whether it's making progress at all,
                                  or genuinely blocked on one call
   3. lsof -p PID                    see WHAT it has open -- a file on
                                  a hung network mount? A socket to an
                                  unresponsive service?
```

### Real-world example

```
   `lsof +L1` (or the closely related `lsof | grep deleted`) is a
   genuinely standard, widely-documented, real fix for a specific,
   common production mystery: "`df` shows the disk is almost full, but
   `du` can't find anything that adds up to that much space" -- exactly
   the scenario Chapter 15 described conceptually (a deleted file whose
   space isn't freed because a process still has it open), made
   PRACTICALLY SOLVABLE with this one command. This exact combination
   -- `df` disagreeing with `du`, then `lsof` revealing a large,
   DELETED-but-still-open log file held by a long-running process
   (often one that was never restarted after a log-rotation tool
   deleted/rotated its log file out from under it) -- is documented
   across countless sysadmin troubleshooting guides and forum threads
   as one of the standard "disk space mystery" playbooks, precisely
   because it recurs constantly across real production systems whenever
   log rotation (Chapter 34) and long-running processes interact badly.
```

### Try it yourself

```bash
# See strace reveal exactly what a command does, syscall by syscall
strace -c ls / > /dev/null    # summary: which syscalls, how many
strace ls / 2>&1 | head -15    # the actual sequence, in order

# See lsof answer "who's using this port"
python3 -m http.server 8765 &
sleep 1
lsof -i :8765
kill %1

# Reproduce and then SOLVE Chapter 15's "deleted but not freed" mystery
python3 -c "
import time, os
f = open('/tmp/bigfile.tmp', 'w')
f.write('x' * 20_000_000)
f.flush()
os.remove('/tmp/bigfile.tmp')   # 'deleted' -- but still open
time.sleep(8)
" &
sleep 2
lsof +L1 2>/dev/null | grep bigfile   # there it is -- a deleted file
                                       # still consuming real disk space
wait
```

### Common mistakes

- **Reaching for `strace` on a process handling large volumes of traffic
  without filtering.** It can slow the process down noticeably and produce
  an overwhelming amount of output — use `-e trace=` to filter to what you
  actually care about.
- **Forgetting `lsof` needs elevated permissions to see OTHER users'
  processes.** `sudo lsof ...` is often necessary outside your own
  processes.
- **Only checking `du` when disk space doesn't add up.** `du` reports space
  used by files STILL VISIBLE in the directory tree — it can't see deleted-
  but-open files at all; that's exactly where `lsof +L1` comes in.

### Check yourself

1. What does `strace` show you, and what does `lsof` show you — how are
   they complementary?
2. What's the standard three-step sequence for diagnosing a "stuck" process?
3. Why can `df` and `du` disagree about how much disk space is actually
   used, and what tool resolves the mystery?

### Further reading

- **Reference:** `man 1 strace`, `man 8 lsof`.
- **Book:** *Systems Performance* (2nd ed.), Gregg — covers both tools (and
  many more) in the context of a full diagnostic methodology.
- **Article:** search "linux disk space deleted file still open lsof" for
  concrete, worked troubleshooting examples.

---

## Chapter 43 — Resource limits: `ulimit` and your first taste of cgroups

### In one sentence

`ulimit` caps what a SINGLE process (or shell session) can consume — open
files, memory, processes — while **cgroups** (control groups) let the kernel
enforce resource limits on a whole GROUP of processes together, and cgroups
are the exact mechanism that makes container resource limits (Phase 3)
possible at all.

### Why this matters

This chapter is the direct bridge into Phase 3: everything Docker and
Kubernetes do to limit a container's CPU/memory is built on cgroups, which
are built on ideas you've now already learned (Chapters 8, 10, 13).

### How it actually works

```bash
   ulimit -a                     # show ALL current limits for this
                                # shell (and anything it launches)
   ulimit -n                      # just open file descriptors
                                # (Chapter 23) -- a VERY commonly hit
                                # limit in real production systems
   ulimit -n 4096                  # raise it (for THIS shell session;
                                # permanent changes need config files --
                                # `/etc/security/limits.conf` on most
                                # systems)
   ulimit -u                        # max number of processes this user
                                # can create -- Chapter 31's fork-bomb
                                # defence, in practice

   CGROUPS (control groups) let the kernel group a set of PROCESSES
   together and apply RESOURCE LIMITS to the GROUP AS A WHOLE -- CPU
   time, memory, I/O bandwidth, process count -- enforced directly by
   the kernel, not just requested politely:

     /sys/fs/cgroup/...    cgroups are exposed as a VIRTUAL FILESYSTEM
                          (Chapter 23's exact philosophy again) -- you
                          create a cgroup by making a directory, and
                          configure it by writing to files inside it.

   THIS IS PRECISELY THE MECHANISM CONTAINERS USE (Phase 3, in full):
   when you run `docker run --memory=512m ...`, Docker is, under the
   hood, creating a cgroup for that container's processes and writing
   `512m` into that cgroup's memory-limit file. If the container's
   processes collectively try to exceed it, the KERNEL'S OOM killer
   (Chapter 13) -- scoped specifically to that cgroup -- kills a process
   inside it, exactly as Chapter 13 described, just now applied to a
   GROUP of processes instead of the whole machine.
```

### Real-world example

```
   "Too many open files" (often surfaced as the error `EMFILE` at the
   syscall level, or a plain "too many open files" message in
   application logs) is one of the most common, most WELL-DOCUMENTED
   real production incidents caused directly by an unraised `ulimit`:
   the default per-process open-file-descriptor limit on many Linux
   distributions has historically been a relatively low number (1024 is
   a very common historical default) -- entirely reasonable for a
   typical desktop process, but FAR too low for a high-connection-count
   server process, where EVERY open network connection ALSO consumes a
   file descriptor (Chapter 19/23).
   This is so well-documented and so common that major server software
   -- nginx, Redis, PostgreSQL, and Elasticsearch all explicitly, in
   their OWN official production-deployment documentation, instruct
   operators to RAISE the open-file-descriptor ulimit before running
   them under real load, specifically because the DEFAULT will cause
   real, production connection failures once traffic exceeds roughly a
   thousand concurrent connections -- a limit that's trivially reachable
   for any moderately successful real-world service. This is a genuine,
   concrete, and extremely common real-world case of Chapter 4's
   abstract "processes aren't free, and neither are the things they
   hold open" lesson becoming an actual, documented, easily Googled
   production outage.
```

### Try it yourself

```bash
# See your current limits
ulimit -a

# See the specific, commonly-hit one
ulimit -n

# Watch a process hit the open-file limit (safe, contained demo)
(
  ulimit -n 10          # deliberately, TEMPORARILY set it very low,
                        # in a SUBSHELL so it doesn't affect your real
                        # session
  python3 -c "
files = []
try:
    for i in range(20):
        files.append(open('/tmp/testfile_%d' % i, 'w'))
        print(f'opened file {i}')
except OSError as e:
    print(f'FAILED at file {len(files)}: {e}')
import os
for i in range(len(files)):
    os.remove('/tmp/testfile_%d' % i)
"
)

# Explore cgroups on your own system, read-only
ls /sys/fs/cgroup/ 2>/dev/null | head -15
cat /sys/fs/cgroup/memory.max 2>/dev/null || \
  echo "(cgroup v2 path may differ, or you may need sudo -- this is
  informational only)"
```

### Build it in Go
[Chapter 78](#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections) hits `EMFILE` ("too many open files") on
purpose and shows that Go quietly raises its own soft `RLIMIT_NOFILE` to the
hard limit (minus one) at startup. [Chapters 75–76](#chapter-75-cpu-limits-gomaxprocs-cgroups-and-throttling)
read cgroup v2 files (`cpu.max`, `cpu.stat`, `memory.max`) from inside a
container.

### Common mistakes

- **Running a high-connection-count service on default ulimits without
  checking.** A very common, avoidable real production outage — check your
  server software's OWN documentation for its recommended limits.
- **Confusing `ulimit` (per-process/per-user) with cgroups (per-group,
  kernel-enforced, what containers actually use).** They're related ideas,
  solving overlapping but distinct problems, at different granularities.
- **Setting a `ulimit` change in your current shell and expecting it to
  persist.** It only applies to that shell (and its children) — permanent
  changes need `/etc/security/limits.conf` or a systemd unit's `LimitNOFILE=`
  directive (Chapter 33).

### Check yourself

1. What's the difference between `ulimit` and a cgroup, in terms of scope?
2. Why does a high-connection-count server commonly need its open-file
   ulimit raised?
3. How does this chapter's cgroup mechanism directly become Phase 3's
   container memory limits?

### Further reading

- **Reference:** `man 1 ulimit` (as part of `man 1 bash`), `man 5
  limits.conf`, `man 7 cgroups`.
- **Docs:** nginx's, Redis's, or PostgreSQL's own official production
  deployment/tuning documentation — search any one of them for "ulimit" or
  "open files" to see the real, documented recommendation.
- **Book:** *Systems Performance* (2nd ed.), Gregg — the cgroups material
  connects directly into container performance analysis.

---

### End of Phase 2 — Milestone check

- [ ] I can navigate the Linux filesystem confidently and explain WHY
      `/etc`, `/var`, and `/usr` hold what they hold
- [ ] I'm fluent with `grep`/`sed`/`awk`/pipes for real text-processing tasks
- [ ] I understand file permissions, ownership, hard vs. symbolic links, and
      what SUID actually does
- [ ] I can watch and control processes with `ps`/`top`/`kill`, and manage
      background jobs that survive a disconnect
- [ ] I can write and manage a real systemd service, with a sensible
      restart policy
- [ ] I know where logs actually go, and why log rotation is a day-one
      deployment concern, not an afterthought
- [ ] **I can write a defensive bash script with `set -euo pipefail`, and
      explain exactly what each flag protects against**
- [ ] I understand the cron environment gap and how to guard against it
- [ ] **I can correctly read load average, `vmstat`, and `iostat` — and
      explain why "high load" doesn't always mean "CPU pressure" on Linux**
- [ ] I can use `strace` and `lsof` together to diagnose a stuck process or
      a disk-space mystery
- [ ] I understand `ulimit` and have a first working mental model of
      cgroups, ready for Phase 3

---

# PHASE 3 — Containerization: Docker and Kubernetes

Everything from here builds directly on Phases 1 and 2. A container is not a
new kind of virtual computer — it's an ordinary Linux process, given its own
view of the filesystem, its own process tree, its own network stack, and a
capped slice of resources, using EXACTLY the OS mechanisms you already know:
namespaces, cgroups, and the filesystem model. If a container ever does
something confusing, the answer is almost always "go back to Phase 1's
mental model and apply it here."

> **Going further:** this phase gets you to "I understand containers and
> Kubernetes". [Docker & Kubernetes — The Real-Life Field Guide](../K8s-docker/real-life-k8s-guide.md)
> continues to production: hardened images, probes and graceful shutdown,
> autoscaling, OpenTelemetry observability, GitOps, canary releases and
> disaster recovery, built as one system in a Linux lab VM.

# Part 14 — Container fundamentals

Every part so far was one OS running many processes with the illusion of
isolation (Part 1's lie, again). A container is that same illusion, turned
into a deliberate, hardened feature. This part shows that namespaces and
cgroups are the real mechanism behind it — not magic, and nothing you
haven't already learned the pieces of.

## Chapter 44 — What a container actually is (namespaces + cgroups + a filesystem)

### In one sentence

A container is a normal Linux process running with its own PID namespace
(so it sees only its own processes), its own network namespace (Chapter 37's
exact model, isolated), a cgroup capping its resource usage (Chapter 43), and
its own root filesystem — no new kernel feature invented specifically for
"containers" at all.

### Why this matters

"How is a container different from a process" and "why did my container get
OOM-killed" both stop being mysterious the moment you see containers as a
COMPOSITION of things you already learned in Phase 1, not a separate
technology.

### How it actually works

```
   FOUR INGREDIENTS, each one YOU ALREADY KNOW:

     1. PID NAMESPACE      the container's main process becomes PID 1
                          INSIDE its own namespace (Chapter 6's process
                          tree, but the container gets its OWN, private
                          root of that tree) -- it cannot see, and
                          cannot signal (Chapter 20), processes outside
                          its namespace.

     2. NETWORK NAMESPACE   its own network interfaces, IP address, and
                          routing table (Chapter 37) -- including its
                          own, SEPARATE `localhost`, exactly as Chapter
                          37's real-world example previewed.

     3. MOUNT NAMESPACE / ROOT FILESYSTEM   its own view of the
                          filesystem (Chapter 24's hierarchy), typically
                          built from a container IMAGE (Chapter 46) --
                          so it can have its OWN `/etc`, its own
                          installed packages, entirely separate from the
                          host's, even though it's ultimately running on
                          the SAME kernel.

     4. CGROUP             a resource limit (CPU, memory -- Chapters 8,
                          10, 13, 43) applied to the container's
                          process(es) as a group, enforced by the
                          kernel, not politely requested.

   PUT TOGETHER: a container process, from the KERNEL'S perspective, is
   still just an ordinary process -- `ps aux` on the HOST machine shows
   it, right alongside every other process, with its own real PID in the
   host's OWN, top-level PID namespace. It's just been given a
   deliberately RESTRICTED VIEW of the system (can't see other
   namespaces' processes/network) and a CAPPED SLICE of resources
   (cgroups) -- fundamentally the same isolation idea as Chapter 1's
   "the OS as a resource manager" and Chapter 10's virtual memory, just
   composed at a coarser, whole-process-environment level.
```

### Real-world example

```
   Docker (first publicly introduced by Solomon Hykes at a PyCon
   lightning talk in 2013) is often, inaccurately, credited with
   "inventing containers" -- the real, well-documented history is more
   interesting and more directly illustrates this chapter's point: the
   underlying Linux kernel mechanisms Docker relies on -- namespaces
   (introduced incrementally into the kernel starting in the early-to-
   mid 2000s) and cgroups (merged into the mainline kernel in 2007,
   developed originally at Google for exactly this resource-isolation
   purpose) -- PREDATED Docker by years. LXC (Linux Containers), an
   earlier project, already used these SAME primitives directly to
   provide container-like isolation before Docker existed.
   What Docker actually, genuinely contributed -- and this is well
   documented across its own project history and countless retrospective
   engineering writeups -- was a dramatically better USER EXPERIENCE: a
   simple, portable IMAGE FORMAT (Chapter 46) that could be built once
   and run identically anywhere, a straightforward CLI, and a
   registry/distribution model (Chapter 52) for sharing images. Docker's
   genuine innovation was PACKAGING AND ERGONOMICS around existing
   kernel primitives, not the isolation mechanism itself -- which is
   exactly why this chapter can teach you "what a container is" using
   ONLY concepts from Phase 1: the primitives really are that
   fundamental, and really did exist before the word "container" became
   an industry-wide household term.
```

### Try it yourself

```bash
# See a container process from the HOST's perspective -- it's just a
# normal process (requires Docker -- if not installed yet, install it
# now: get.docker.com has the official install script, or use your
# distro's package manager)
docker run -d --name pidtest nginx
docker inspect --format '{{.State.Pid}}' pidtest   # the REAL, host-
                                                    # level PID
ps -p $(docker inspect --format '{{.State.Pid}}' pidtest)  # there it
                                                            # is, in
                                                            # your
                                                            # NORMAL
                                                            # host `ps`
                                                            # output

# See the PID namespace difference: PID 1 inside vs. outside
docker exec pidtest ps aux    # inside the container, nginx is PID 1
ps -p $(docker inspect --format '{{.State.Pid}}' pidtest)  # outside,
                                                            # it's some
                                                            # other,
                                                            # larger PID

# See the cgroup this container was actually assigned
docker inspect --format '{{.HostConfig.Memory}}' pidtest
cat /sys/fs/cgroup/system.slice/docker-*.scope/memory.max 2>/dev/null | head -1 || \
  echo "(cgroup path varies by cgroup version/driver -- informational)"

docker rm -f pidtest
```

### Build it in Go
The best way to stop thinking of containers as magic is to build one.
[Chapter 80](#chapter-80-a-container-runtime-in-150-lines-of-go) does it in about 150 lines of Go: new UTS,
PID, mount, IPC, and network namespaces; `chroot` into a BusyBox filesystem;
a private `/proc`; and a cgroup v2 limit that stops a fork bomb at 20
processes. For how attackers break out of containers, and how to stop them,
see [Security Engineering in Depth, Chapter 15](../security/real-life-security-guide-v1.md#chapter-15-container-internals-and-isolation).

### Common mistakes

- **Thinking a container is "a lightweight VM."** Chapter 45 makes this
  precise, but the short version: there's no virtualised hardware, no
  separate kernel — it's a specially-isolated ordinary process, sharing the
  HOST's kernel directly.
- **Assuming "container" is a Linux kernel feature you can point to.**
  There's no `container` syscall or kernel object — it's a COMPOSITION of
  namespaces, cgroups, and a filesystem, orchestrated by tooling (Docker,
  containerd, etc.), exactly as this chapter describes.
- **Forgetting the container shares the HOST's kernel.** A kernel-level
  vulnerability or a container escape (a real, documented risk category,
  covered in this wiki's security guides) can cross this isolation boundary
  — it's strong, but it is not the same guarantee as separate hardware
  virtualization.

### Check yourself

1. Name the four Linux mechanisms that, composed together, make a container.
2. Why does `ps aux` on the HOST show a container's processes too?
3. What did Docker actually invent, versus what already existed in the
   Linux kernel before it?

*(Answers: Appendix E.)*

### Further reading

- **Talk:** Solomon Hykes' original 2013 Docker introduction (search "Docker
  PyCon 2013 lightning talk") — genuinely worth watching for the real
  history.
- **Reference:** `man 7 namespaces`, `man 7 cgroups`, `man 7 pid_namespaces`.
- **Article:** search "containers are just Linux processes" — this exact
  framing has been written up excellently, multiple times, by various
  container-runtime engineers.

---

## Chapter 45 — Containers vs. virtual machines, honestly compared

### In one sentence

A virtual machine virtualises the HARDWARE (each VM runs its own full
kernel, on emulated or hypervisor-managed hardware), while a container
virtualises the OPERATING SYSTEM (many containers share one real kernel) —
a genuinely different point in the isolation-vs-overhead trade-off, not
simply "containers are the new, better VMs."

### Why this matters

"Should this run in a container or a VM" is a real, recurring architecture
decision, and answering it well requires understanding the actual trade-off,
not a marketing simplification.

### How it actually works

```
   A VIRTUAL MACHINE, via a HYPERVISOR (software like KVM, Xen, or
   VMware that can run directly on hardware, or on top of a host OS):
     - virtualises actual HARDWARE -- CPU, memory, disk, network devices
       are all presented to the guest as if they were real.
     - each VM runs its OWN, COMPLETE, INDEPENDENT KERNEL -- genuinely
       a separate operating system instance, which can even be a
       DIFFERENT OS entirely from the host (a Linux host running a
       Windows VM, for example).
     - ISOLATION is very strong: a kernel bug or crash inside one VM
       cannot directly affect another VM or the host's own kernel,
       because they're not sharing one.
     - COST: a full kernel boot, a larger memory footprint per instance,
       and generally SLOWER startup (often tens of seconds or more) and
       lower density (fewer instances per physical machine, for the same
       hardware).

   A CONTAINER, as Chapter 44 established:
     - shares the HOST's SINGLE kernel -- every container on a machine
       is running on the SAME kernel, at the SAME time.
     - isolation comes from namespaces + cgroups, at the OS level, not
       hardware virtualization.
     - STARTUP is typically near-instant (milliseconds, Chapter 46
       explains why) since there's no kernel to boot -- you're just
       starting a process.
     - DENSITY is far higher for the same hardware -- no per-instance
       kernel memory overhead, much smaller per-instance footprint.
     - the TRADE-OFF: isolation is WEAKER than a VM's -- a sufficiently
       severe kernel vulnerability, or a container-escape technique
       (a real, actively-studied security category, covered in this
       wiki's security guides), can potentially cross between
       containers, or from a container to the host, in a way a
       properly-configured VM's hardware-level boundary generally
       resists more strongly.

   THIS IS WHY "SANDBOXED" CONTAINER RUNTIMES EXIST (gVisor, Kata
   Containers, covered in this wiki's advanced security guide) -- they
   sit deliberately IN BETWEEN: Kata runs each container in its own
   lightweight VM (closer to VM-level isolation, container-like
   ergonomics); gVisor intercepts and re-implements syscalls in
   userspace rather than passing them straight to the host kernel. Real
   engineering teams choose points along this spectrum based on how much
   they trust the WORKLOAD running inside.
```

### Real-world example

```
   The industry-wide shift toward running MORE workloads in containers
   rather than one-VM-per-application, over roughly the 2014-2020
   period, is a well-documented, real architectural trend, driven
   directly by this chapter's density and startup-time trade-offs:
   companies operating at real scale (running thousands of application
   instances) found that packing many containers onto shared hosts --
   instead of many separate, per-application VMs, each paying its own
   full-kernel memory and boot-time overhead -- meaningfully reduced
   infrastructure cost and dramatically improved how quickly new
   instances could be started (critical for AUTOSCALING responsively to
   real traffic spikes, since a VM that takes 30-60 seconds to boot
   responds to a traffic spike far more sluggishly than a container that
   starts in under a second). This is precisely why Kubernetes (Part 16)
   -- designed around exactly this rapid-start, high-density model --
   became the dominant orchestration layer for this style of workload,
   while VMs REMAIN the right, still-widely-used choice for workloads
   needing stronger isolation guarantees (genuinely multi-tenant,
   security-sensitive workloads; running a different OS entirely; very
   large, long-lived, stateful single instances) -- this is a real,
   ongoing trade-off decision, not a settled "containers won" story.
```

### Try it yourself

```bash
# Feel the startup-time difference directly (VM comparison is
# illustrative -- most readers won't have a VM handy to time, but the
# container side is immediate and real)
time docker run --rm alpine echo "started"    # note how fast this is
                                              # -- typically well under
                                              # a second once the image
                                              # is already pulled

# See how many containers you could plausibly run vs. how many full VMs
free -h                          # your available memory
docker run --rm alpine cat /proc/meminfo | head -1   # a minimal
                                                      # container's
                                                      # footprint is
                                                      # tiny -- no
                                                      # separate kernel
                                                      # to account for
```

### Common mistakes

- **"Containers are just lightweight VMs."** They're a fundamentally
  different isolation model (shared kernel vs. separate kernel) — the
  practical CONSEQUENCES (startup time, density) are similar in spirit, but
  the mechanism and the security trade-off are genuinely different.
- **Assuming containers are always the "better" choice.** For genuinely
  untrusted, multi-tenant, or security-critical isolation needs, a VM's
  stronger boundary is often still the right, deliberate choice — this
  wiki's advanced security guide covers exactly this trade-off in production
  contexts.
- **Not knowing sandboxed runtimes (gVisor, Kata) exist.** They're the real,
  practical middle ground when you want container-like ergonomics with
  closer-to-VM isolation strength.

### Check yourself

1. What does a VM virtualise, and what does a container virtualise — state
   the difference precisely.
2. Why do containers typically start dramatically faster than VMs?
3. What trade-off do sandboxed runtimes like gVisor/Kata occupy, and why
   would a team choose one?

### Further reading

- **Book (free):** *OSTEP*'s virtualization chapters give you the VM-side
  concepts (hypervisors); this guide's Phase 1 gives you the container-side
  ones.
- **Guide:** this wiki's advanced security guide covers gVisor, Kata, and
  container-isolation trade-offs in real operational depth.
- **Article:** search "containers vs virtual machines" for current,
  well-illustrated comparisons — many container-runtime vendors publish
  genuinely fair, technically accurate ones.

---

## Chapter 46 — Images and layers: why containers start in milliseconds

### In one sentence

A container image is built from a stack of read-only **layers**, each one a
set of filesystem changes, combined at runtime by a **union filesystem** — and
because unchanged layers can be cached and reused across many containers,
starting a new container almost never means copying a whole filesystem from
scratch.

### Why this matters

Understanding layers explains why `docker build` caching works the way it
does, why images can be surprisingly large or surprisingly small, and why
`docker pull` only downloads what's actually new.

### How it actually works

```
   AN IMAGE is a STACK of LAYERS, each one representing a set of
   filesystem CHANGES (files added, modified, or removed) relative to
   the layer below it. Every instruction in a Dockerfile that changes
   the filesystem (`RUN`, `COPY`, `ADD`) typically creates a NEW layer:

     FROM ubuntu:22.04          <- base layers (the OS itself)
     RUN apt-get install -y curl   <- a new layer: curl's files added
     COPY app.py /app/               <- a new layer: your one file added
     CMD ["python3", "/app/app.py"]  <- metadata only, no filesystem
                                     change, no new layer

   THESE LAYERS ARE READ-ONLY and, crucially, CACHEABLE and SHAREABLE:
   if two different images both start `FROM ubuntu:22.04`, they can
   literally SHARE that base layer on disk -- it's stored ONCE, not
   duplicated. When you `docker pull` an image, Docker only downloads
   layers it doesn't ALREADY have locally.

   A UNION FILESYSTEM (commonly OverlayFS on modern Linux) is what
   makes running a container from these layers actually WORK: it
   presents the STACK of read-only layers as if they were a SINGLE,
   unified filesystem, and adds ONE thin, WRITABLE layer on top,
   specific to THIS running container -- any file the container creates
   or modifies while running goes into THAT layer, leaving the
   underlying image layers completely untouched.

   THIS IS EXACTLY WHY CONTAINER STARTUP IS SO FAST: starting a new
   container from an image ALREADY ON DISK doesn't require copying
   gigabytes of filesystem -- it just needs to set up this thin
   writable layer on top of ALREADY-PRESENT, READ-ONLY layers, then
   start the process (Chapter 44). It's a direct, practical consequence
   of Chapter 15's "renaming/linking is cheap, copying data is
   expensive" lesson, applied at the whole-filesystem scale.
```

### Real-world example

```
   "Why did changing ONE LINE in my Dockerfile trigger a rebuild of
   EVERYTHING" is a genuinely common, well-documented real question new
   Docker users ask, and it's a direct, practical consequence of THIS
   chapter's layer-caching mechanism: Docker's build cache works, by
   design, LAYER BY LAYER, IN ORDER -- if you change (or invalidate) a
   layer, EVERY LAYER AFTER IT in the Dockerfile must be rebuilt too,
   even if those later instructions themselves didn't change at all,
   because each layer is built ON TOP of the previous one's exact
   output.
   This is EXACTLY why Docker's own official best-practices
   documentation, and essentially every serious Dockerfile style guide,
   recommends ORDERING instructions from LEAST-likely-to-change to
   MOST-likely-to-change: install your OS packages and dependencies
   FIRST (these change rarely), and COPY your actual, frequently-
   changing application code LAST. A Dockerfile that does the reverse --
   copying application code early, then installing dependencies -- means
   EVERY single code change (which might happen dozens of times a day
   during active development) invalidates the dependency-installation
   layer too, needlessly re-running a slow `apt-get`/`pip install`/`npm
   install` step on every single build, turning what should be a
   near-instant rebuild into a genuinely painful, minutes-long one. This
   real, common performance mistake -- and its well-documented fix -- is
   entirely explained by understanding layers as this chapter describes
   them, rather than treating Dockerfile ordering as arbitrary.
```

### Try it yourself

```bash
# See an image's actual layers
docker pull alpine
docker history alpine

# Build something yourself and watch layer caching in action
mkdir -p ~/lab/dockerdemo && cd ~/lab/dockerdemo
cat > Dockerfile << 'EOF'
FROM alpine:latest
RUN echo "installing dependencies (pretend)" && sleep 2
COPY app.txt /app.txt
CMD ["cat", "/app.txt"]
EOF
echo "version 1" > app.txt

docker build -t layerdemo .              # full build, note the timing
echo "version 2" > app.txt
docker build -t layerdemo .              # rebuild -- notice the FIRST
                                         # layer (the slow RUN/sleep)
                                         # is CACHED, only the COPY
                                         # layer (and after) rebuilds

# See the difference an image makes on disk
docker images | grep -E "alpine|layerdemo"
```

### Common mistakes

- **Putting frequently-changing files (application code) early in a
  Dockerfile.** Invalidates every subsequent layer's cache on every change —
  order from least-changing to most-changing instead.
- **Assuming a container's writable layer changes are permanent.** By
  default, they're LOST when the container is removed — Chapter 49's volumes
  are the mechanism for data that must actually persist.
- **Not realising shared base layers save real disk space and pull time.**
  Standardising on a small number of base images across your organisation is
  a genuine, practical optimisation, not just tidiness.

### Check yourself

1. What is a layer, and why are layers cacheable/shareable across images?
2. What does the writable layer contain, and what happens to it when a
   container is removed?
3. Why should a Dockerfile install dependencies BEFORE copying application
   code, in terms of build-cache behaviour?

### Further reading

- **Docs:** Docker's own official "Best practices for writing Dockerfiles" —
  the layer-ordering guidance directly, from the source.
- **Reference:** `man 8 mount` (the overlay filesystem type), or search
  "OverlayFS explained" for a clear technical breakdown.
- **Practice:** `docker history <image>` on a few real, popular public
  images — a genuinely instructive way to see real-world layer structure.

---

# Part 15 — Docker in depth

Part 14 showed you what a container actually is from the kernel's point of
view — namespaces, cgroups, a filesystem. This part is the practical layer
built on top of that mechanism: Docker specifically, running and
understanding real containers rather than the primitives underneath them.

## Chapter 47 — Running your first containers, and understanding what happened

### In one sentence

`docker run` does several distinct things in sequence — pull an image if
needed, create a container from it, and start the process inside — and
naming each step precisely turns "Docker magic" into a small, predictable
sequence you fully understand.

### Why this matters

Every Docker command you'll ever run is a variation on this same handful of
verbs — knowing them precisely, not just by muscle memory, is what lets you
debug when something doesn't behave as expected.

### How it actually works

```bash
   docker pull ubuntu:22.04       # download an IMAGE (Chapter 46) from
                                 # a registry (Chapter 52), if not
                                 # already present locally
   docker run ubuntu:22.04 echo "hi"    # CREATE a container from that
                                 # image, and START the process inside
                                 # it (this implicitly pulls first, if
                                 # needed)
   docker ps                      # list RUNNING containers
   docker ps -a                     # list ALL containers, including
                                 # stopped ones
   docker stop <container>           # send SIGTERM (Chapter 20!),
                                 # wait, then SIGKILL if it hasn't
                                 # exited -- EXACTLY Kubernetes' pod
                                 # termination sequence from Chapter
                                 # 20, one layer down
   docker rm <container>              # remove a STOPPED container
                                 # (and its writable layer, Chapter 46
                                 # -- any unsaved changes are gone)
   docker exec -it <container> bash    # run an ADDITIONAL process
                                 # inside an ALREADY-RUNNING container
                                 # -- useful for poking around, NOT
                                 # how the container's main process
                                 # started
   docker logs <container>              # see the STDOUT/STDERR
                                 # (Chapter 28's file descriptors 1/2)
                                 # of the container's main process

   FLAGS WORTH KNOWING IMMEDIATELY:
     -d          detached -- run in the background, return your
                terminal immediately (Chapter 32's background jobs,
                one layer up)
     -it          interactive + a pseudo-TTY -- lets you actually type
                into the container's process (needed for `docker exec
                ... bash` to feel like a real shell)
     --rm          automatically remove the container once it exits --
                good hygiene for throwaway/test runs
     --name X       give it a memorable name instead of a random one
     -p HOST:CONTAINER   publish a port (Chapter 50 covers this properly)
```

### Real-world example

```
   `docker run hello-world` is, quite literally, Docker's OWN official
   "does this actually work" onboarding test -- documented in Docker's
   getting-started material since the project's early public releases,
   and still the first command most engineers ever run with Docker
   installed. It's worth understanding EXACTLY what it does, because
   it's a genuinely complete, minimal demonstration of this chapter's
   whole sequence: Docker checks for the `hello-world` image locally
   (usually not found, the first time), PULLS it from Docker Hub
   (Chapter 52's registry), CREATES a new container from it, STARTS the
   tiny program inside (which just prints an explanatory message and
   exits immediately), and the container then sits, STOPPED (Chapter 3's
   exit -- the process finished, exit status 0), visible in `docker ps
   -a` but not `docker ps`.
   Running it once, and then deliberately walking through `docker ps
   -a`, `docker logs`, and `docker rm` on the resulting stopped
   container -- rather than just trusting the friendly printed message
   -- is a genuinely good way to make this chapter's whole vocabulary
   concrete on your very first real Docker command.
```

### Try it yourself

```bash
# The official one, examined properly
docker run hello-world
docker ps -a                     # see it, STOPPED, exit code 0
docker logs $(docker ps -a -q --filter ancestor=hello-world | head -1)

# Run something longer-lived, and practice the full lifecycle
docker run -d --name webtest nginx
docker ps                          # RUNNING
curl -s localhost:80 > /dev/null 2>&1 || echo "(not published yet --
  Chapter 50 covers -p properly; this container's port 80 isn't
  reachable from the host yet)"
docker exec webtest ls /usr/share/nginx/html    # poke around inside
docker logs webtest
docker stop webtest                 # SIGTERM, then SIGKILL if needed
docker ps -a                          # STOPPED, not gone
docker rm webtest                      # actually removed now
docker ps -a | grep webtest || echo "confirmed: gone"

# Clean up hello-world containers too
docker rm $(docker ps -a -q --filter ancestor=hello-world) 2>/dev/null
```

### Common mistakes

- **Confusing `docker stop` (graceful, Chapter 20's SIGTERM) with `docker
  kill` (immediate SIGKILL).** Use `stop` by default; reach for `kill` only
  when a container genuinely won't respond to the polite request.
- **Forgetting `docker rm` is needed after `docker stop`.** A stopped
  container still exists (and its writable layer still takes disk space)
  until explicitly removed.
- **Using `docker exec` and assuming it's "how the container started."**
  `exec` runs an ADDITIONAL process alongside the main one — killing your
  `exec`'d shell does NOT stop the container itself.

### Check yourself

1. List, in order, what `docker run` actually does the first time you run
   a given image.
2. What's the difference between `docker stop` and `docker rm`?
3. Why does killing a `docker exec`'d shell not stop the container?

*(Answers: Appendix E.)*

### Further reading

- **Docs:** Docker's own "Getting Started" guide — still one of the
  clearest, most concrete introductions available, directly from the source.
- **Reference:** `docker --help`, `docker run --help` (genuinely worth
  reading in full once).
- **Practice:** run `docker events` in one terminal while you run commands
  in another — watch the exact sequence of internal events fire in real time.

---

## Chapter 48 — Writing a good Dockerfile

### In one sentence

A Dockerfile is a recipe for building an image, layer by layer (Chapter 46),
and a genuinely GOOD one is small, cacheable, and — critically — doesn't run
its main process as root by default.

### Why this matters

A careless Dockerfile produces a bloated, slow-to-build, slow-to-deploy image
that also runs with more privilege than it needs — real, common, and
avoidable technical debt from day one.

### How it actually works

```dockerfile
   # A genuinely good, real-world-shaped example
   FROM python:3.12-slim                 # a SMALL base image, not the
                                        # full, much larger default

   WORKDIR /app

   COPY requirements.txt .                # dependencies FIRST (Chapter
                                        # 46's caching lesson) --
                                        # changes rarely
   RUN pip install --no-cache-dir -r requirements.txt

   COPY . .                                # application code LAST --
                                        # changes often

   RUN useradd --create-home appuser        # create a NON-ROOT user
   USER appuser                              # switch to it -- the
                                        # container's main process now
                                        # runs as appuser, NOT root

   EXPOSE 8000                                # documentation -- doesn't
                                        # actually publish the port
                                        # (Chapter 50 does that)
   CMD ["python3", "app.py"]
```

**The habits that separate a good Dockerfile from a careless one:**

```
   [ ] a SMALL base image (a `-slim`/`-alpine` variant, or a distroless
       image) -- smaller attack surface, faster pulls, less disk.
   [ ] dependencies installed BEFORE application code (Chapter 46).
   [ ] a `.dockerignore` file (works like `.gitignore`) -- excludes
       things like `.git/`, local virtual environments, and secrets
       from ever being COPY-ed into the image at all.
   [ ] a NON-ROOT user for the actual running process.
   [ ] PINNED versions (a specific base image tag, pinned dependency
       versions) rather than `latest`/unpinned -- reproducible builds,
       exactly this wiki's advanced security guide's supply-chain
       chapter's lesson, applied here.
   [ ] MULTI-STAGE builds for compiled languages -- build in one stage
       (with all the compilers/build tools), copy ONLY the final
       artifact into a clean, minimal final stage, so build tooling
       never ships in your production image at all.
```

### Real-world example

```
   "Don't run your container's process as root" is one of the most
   consistently, widely repeated pieces of real container-security
   guidance in the industry -- present in Docker's own official
   documentation, in every major container-security scanning tool's
   default findings, and covered in operational depth in this wiki's
   security guides' container-hardening material. The concrete, real
   reasoning: a container running as root, if COMPROMISED (through an
   application vulnerability, a malicious dependency, or any other real
   attack path), gives an attacker ROOT WITHIN the container's namespace
   -- and root, combined with ANY of the real, documented container-
   escape techniques (a misconfigured volume mount reaching sensitive
   host paths, a kernel vulnerability, an overly-permissive capability
   set), is a materially WORSE starting position for the attacker than
   a non-root user would have been. This is precisely why "USER
   appuser" (or the equivalent) is treated as close to a non-negotiable
   default in serious Dockerfile guidance, not an optional hardening
   step reserved for especially sensitive workloads -- the cost of
   adding it is a few lines; the benefit, in a real compromise, is
   substantial and well documented across container-security research.
```

### Try it yourself

```bash
cd ~/lab/dockerdemo

# See the root-vs-non-root difference directly
cat > Dockerfile.root << 'EOF'
FROM alpine:latest
CMD ["whoami"]
EOF
docker build -t rootdemo -f Dockerfile.root .
docker run --rm rootdemo             # prints "root"

cat > Dockerfile.nonroot << 'EOF'
FROM alpine:latest
RUN adduser -D appuser
USER appuser
CMD ["whoami"]
EOF
docker build -t nonrootdemo -f Dockerfile.nonroot .
docker run --rm nonrootdemo          # prints "appuser" -- NOT root

# Compare image sizes: slim base vs full
docker pull python:3.12-slim
docker pull python:3.12
docker images | grep python

rm -f Dockerfile.root Dockerfile.nonroot
```

### Build it in Go
Go is unusually well suited to small, safe images: a static binary
(`CGO_ENABLED=0`) needs no libc, so the final stage can be `distroless/static`
or even `scratch`. [Chapter 82](#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes) has a tested multi-stage
Dockerfile that produced a **14.8 MB** image that runs as non-root, read-only,
with all capabilities dropped, and with no shell to exploit.

### Common mistakes

- **Using `FROM ubuntu:latest` (or any `:latest` tag) in production.** Not
  reproducible — what "latest" points to changes over time, silently, out
  from under you. Pin a specific version.
- **`COPY . .` before installing dependencies.** Defeats Chapter 46's layer
  caching on every single code change — order matters.
- **Never creating a `.dockerignore`.** Without one, your build context
  (and potentially your final image) can include `.git` history, local
  secrets, or virtual environment junk you never intended to ship.

### Check yourself

1. Why does dependency installation belong BEFORE application code in a
   Dockerfile?
2. What real security benefit does running as a non-root user provide,
   specifically?
3. What does `.dockerignore` do, and why does it matter?

### Further reading

- **Docs:** Docker's own official "Best practices for writing Dockerfiles."
- **Guide:** this wiki's security guides' container-hardening chapters cover
  the full "harden a container" checklist (capabilities, read-only
  filesystems, and more) in operational depth.
- **Tool:** `hadolint` — a Dockerfile linter that catches many of these
  mistakes automatically; worth running in CI.

---

## Chapter 49 — Volumes and data: what survives a container's death

### In one sentence

A container's writable layer (Chapter 46) disappears the moment the
container is removed — a **volume** is Docker's mechanism for data that
needs to outlive any single container, stored OUTSIDE that disposable layer
entirely.

### Why this matters

"I ran `docker rm` on my database container and lost all my data" is one of
the single most common, most painful real mistakes new Docker users make —
and it's entirely, easily avoidable once this model is clear.

### How it actually works

```bash
   docker volume create mydata          # create a NAMED volume --
                                        # managed by Docker, stored
                                        # OUTSIDE any container's
                                        # writable layer
   docker run -v mydata:/var/lib/postgresql/data postgres  # MOUNT it
                                        # into the container at a
                                        # specific path -- anything the
                                        # container writes there goes
                                        # into the VOLUME, not its
                                        # disposable writable layer

   docker run -v /host/path:/container/path myimage    # a BIND MOUNT --
                                        # maps a SPECIFIC HOST DIRECTORY
                                        # directly into the container
                                        # (useful for local development
                                        # -- edit code on the host, see
                                        # changes immediately inside the
                                        # running container)

   docker volume ls                       # list all named volumes
   docker volume inspect mydata            # where does it actually
                                        # live on the host's disk?
   docker volume rm mydata                  # DELETE it -- and the
                                        # data with it. Deliberately
                                        # separate from removing any
                                        # container.
```

**The mental model that prevents the classic mistake:**

```
   `docker rm <container>`    removes the CONTAINER and its writable
                              layer (Chapter 46). Anything stored ONLY
                              there is GONE.
   `docker volume rm <vol>`    removes the VOLUME, separately, and
                              explicitly. Data in a properly-mounted
                              volume SURVIVES any number of `docker rm`
                              operations on containers that used it --
                              it's a genuinely separate lifecycle,
                              precisely so this can't happen by
                              accident.
```

### Real-world example

```
   "I ran a database in a container without a volume, then removed the
   container, and lost everything" is a genuinely, extremely common real
   mistake -- so common that Docker's own official documentation
   explicitly, prominently warns about exactly this scenario, and it's
   one of the most frequently asked/answered questions across every
   Docker-related community and support forum. The typical shape: an
   engineer runs `docker run postgres` (or any stateful service) WITHOUT
   a `-v` volume flag during initial exploration or local development,
   accumulates real, valuable data inside the container over hours or
   days of work, then runs a routine `docker rm` (perhaps as part of an
   unrelated cleanup, or to apply an image update) -- and the data,
   having lived ONLY in that container's disposable writable layer, is
   gone, with no warning and no recovery path.
   THE LESSON, and it generalises past Docker specifically: any
   STATEFUL workload running in a container MUST have its actual data
   directory mounted to a volume (or bind mount) from the very first
   time you run it, not added "later, once it matters" -- because by the
   time it obviously matters, the container holding the only copy of the
   data may already be gone.
```

### Try it yourself

```bash
# Reproduce the classic mistake, safely, with throwaway data
docker run -d --name nodata alpine sh -c "echo important data > /data.txt; sleep 300"
docker exec nodata cat /data.txt      # it's there
docker rm -f nodata                    # remove the container
docker run --rm alpine cat /data.txt 2>&1 || echo "gone forever -- exactly the mistake"

# Now do it correctly, with a volume
docker volume create demodata
docker run -d --name withdata -v demodata:/storage alpine sh -c "echo important data > /storage/data.txt; sleep 300"
docker exec withdata cat /storage/data.txt
docker rm -f withdata                   # remove the CONTAINER...
docker run --rm -v demodata:/storage alpine cat /storage/data.txt  # ...but
                                        # the data survived, in the volume

docker volume rm demodata
```

### Common mistakes

- **Running any stateful service without a volume "just for now."** By the
  time it holds real, valuable data, removing the container without one is
  an easy, catastrophic mistake to make.
- **Confusing a bind mount with a named volume.** A bind mount ties you to a
  SPECIFIC host path (great for local dev, awkward for portability); a named
  volume is managed by Docker and more portable across environments.
- **Forgetting `docker volume rm` is a SEPARATE, deliberate action.** This
  is a feature, not friction — it's specifically what prevents accidental
  data loss when cleaning up containers.

### Check yourself

1. What happens to a container's writable-layer data when the container is
   removed?
2. What's the difference between a named volume and a bind mount?
3. Why does Docker require a separate, explicit command to delete a volume,
   rather than deleting it automatically with its container?

### Further reading

- **Docs:** Docker's own official "Manage data in Docker" documentation —
  volumes, bind mounts, and tmpfs mounts covered thoroughly.
- **Reference:** `docker volume --help`.
- **Practice:** deliberately reproduce the "lost data" mistake once, in your
  lab, exactly as shown above — a lesson learned safely sticks far better
  than one read about.

---

## Chapter 50 — Container networking: how containers talk to the world

### In one sentence

By default, Docker creates a private, virtual network for your containers to
talk to EACH OTHER by name, and you must explicitly **publish** a port to
make a container reachable from OUTSIDE that network — both mechanisms are
Chapter 37's networking concepts, applied inside a container's own network
namespace.

### Why this matters

"Why can't I reach my container from my browser" and "why can my two
containers talk to each other by name" are both answered by understanding
exactly what Docker's networking layer is doing.

### How it actually works

```bash
   docker network ls                      # Docker creates a default
                                          # "bridge" network; you can
                                          # also make your own

   docker run -p 8080:80 nginx              # PUBLISH: map the HOST's
                                          # port 8080 to the
                                          # CONTAINER's port 80 --
                                          # WITHOUT this flag, the
                                          # container's port is only
                                          # reachable from OTHER
                                          # containers on the same
                                          # Docker network, never from
                                          # the host or outside world

   docker network create mynet               # a custom, USER-DEFINED
                                          # network
   docker run --network mynet --name web nginx
   docker run --network mynet --name db postgres
   # `web` can now reach `db` simply as "db" -- Docker runs a built-in
   # DNS resolver on user-defined networks, resolving CONTAINER NAMES
   # to their internal IP addresses automatically (Chapter 37's DNS
   # concept, provided for you)

   THE NETWORK DRIVERS, briefly:
     bridge (the default)    an isolated virtual network on this host;
                            containers on it can reach each other; the
                            outside world can reach IN only via
                            published ports.
     host                     the container uses the HOST's network
                            namespace DIRECTLY, no isolation at all --
                            faster, but loses the isolation Chapter 44
                            described.
     none                      no networking at all.
     overlay                    spans MULTIPLE HOSTS -- the foundation
                            for Docker Swarm, and conceptually a
                            preview of what Kubernetes' networking
                            model (Part 16) does at a larger scale.
```

### Real-world example

```
   The reason Docker Compose (Chapter 51) creates a DEDICATED,
   PROJECT-SPECIFIC network automatically, rather than putting
   everything on the shared default "bridge" network, is directly this
   chapter's DNS-by-container-name mechanism, made practical: Compose's
   own official documentation describes exactly this behaviour -- every
   service defined in a `docker-compose.yml` file can reach every OTHER
   service simply by its SERVICE NAME (e.g. an application container
   connecting to `postgresql://db:5432/...`, where `db` is just the
   name given to the database service in the same compose file), with
   NO manual IP address configuration, NO manual network creation, and
   NO hardcoded addresses anywhere -- because Compose automatically
   creates a user-defined network (this chapter's exact mechanism) for
   the whole project and gives every service that network's built-in
   name resolution for free.
   This is a genuinely practical, extremely common real pattern: almost
   every real multi-container local-development setup you'll encounter
   relies on this exact "reach the other service by its plain name"
   mechanism, and it stops feeling like magic the moment you know it's
   just Chapter 37's DNS resolution, running on a Docker-managed virtual
   network, automatically configured for you.
```

### Try it yourself

```bash
# See the difference publishing makes
docker run -d --name unpub nginx
curl -s localhost:80 > /dev/null 2>&1 && echo "reachable" || echo "NOT reachable (expected -- no -p flag)"
docker rm -f unpub

docker run -d --name pub -p 8080:80 nginx
sleep 1
curl -s localhost:8080 > /dev/null 2>&1 && echo "reachable now!" || echo "still not reachable -- check Docker is running"
docker rm -f pub

# See container-name DNS resolution on a user-defined network
docker network create demonet
docker run -d --name pingtarget --network demonet alpine sleep 300
docker run --rm --network demonet alpine ping -c 2 pingtarget    # reaches
                                                                 # it BY
                                                                 # NAME
docker rm -f pingtarget
docker network rm demonet
```

### Common mistakes

- **Forgetting `-p` and assuming a container should just be reachable.**
  Without publishing, it's isolated on Docker's internal network — invisible
  from the host or outside world by design.
- **Confusing `-p 8080:80` order.** It's always `HOST:CONTAINER` — a common,
  easy-to-transpose mistake.
- **Putting everything on the default "bridge" network and wondering why
  name-based resolution doesn't work.** Container-name DNS resolution only
  works on USER-DEFINED networks — create one (or let Compose do it for you).

### Check yourself

1. What does `-p 8080:80` actually do, and what's reachable without it?
2. Why can containers on a user-defined network reach each other by name,
   with no manual configuration?
3. Why does the default "bridge" network NOT support this name-based
   resolution?

### Further reading

- **Docs:** Docker's own official "Networking overview" and "Networking with
  standalone containers" documentation.
- **Reference:** `docker network --help`.
- **Guide:** this wiki's `networking/tcp-ip` guide covers DNS resolution
  itself in far greater depth than this chapter's Docker-specific summary.

---

## Chapter 51 — Docker Compose: running more than one container together

### In one sentence

Docker Compose lets you describe an entire multi-container application — every
service, its image, its ports, its volumes, its dependencies on other
services — in ONE declarative YAML file, and bring the whole thing up or down
with a single command.

### Why this matters

Almost no real application is a single container — Compose is how you
develop and test a realistic, multi-service setup locally, and it's the
direct conceptual predecessor to Kubernetes' declarative model (Part 16).

### How it actually works

```yaml
   # docker-compose.yml
   services:
     web:
       build: .                       # build FROM a local Dockerfile
       ports:
         - "8080:80"                   # Chapter 50's -p, declaratively
       depends_on:
         - db                           # start `db` first (doesn't wait
                                       # for it to be actually READY --
                                       # a common, real gotcha)
       environment:
         DATABASE_URL: "postgresql://db:5432/myapp"  # Chapter 50's
                                       # container-name DNS resolution,
                                       # used directly

     db:
       image: postgres:16               # or pull an existing image
       volumes:
         - dbdata:/var/lib/postgresql/data   # Chapter 49's volumes,
                                       # declaratively
       environment:
         POSTGRES_PASSWORD: "devpassword"

   volumes:
     dbdata:                            # declare the named volume
```

```bash
   docker compose up                     # build/pull images, create
                                        # the network, start everything,
                                        # in the RIGHT ORDER
   docker compose up -d                    # detached (background)
   docker compose down                      # stop and remove
                                        # everything (containers,
                                        # network -- volumes SURVIVE by
                                        # default, exactly Chapter 49's
                                        # deliberate separation)
   docker compose logs -f web                # follow ONE service's
                                        # logs
   docker compose ps                          # what's running, in
                                        # this project
```

### Real-world example

```
   `depends_on`'s actual behaviour is one of the most common, well-
   documented real points of confusion for engineers new to Compose:
   Docker Compose's own official documentation is explicit that, by
   default, `depends_on` only controls STARTUP ORDER -- it starts `db`
   BEFORE `web`, but it does NOT wait for the database to actually be
   READY to accept connections (a database process can take several
   real seconds to initialise after its CONTAINER has technically
   started, Chapter 3's process-vs-"actually ready" distinction). A
   very common real symptom: `web`'s application code tries to connect
   to `db` immediately on startup, the database isn't accepting
   connections yet, and the application crashes or errors out --
   despite `depends_on` being correctly configured, exactly as
   documented.
   THE REAL, DOCUMENTED FIX Compose's own docs recommend: either build
   RETRY LOGIC into your application's own startup/connection code
   (the generally-preferred, more robust real-world approach, since it
   also protects you in production, not just local dev), or use
   `depends_on`'s `condition: service_healthy` option paired with an
   explicit HEALTHCHECK definition on the dependency service, so
   Compose actually waits for a real readiness signal, not just "the
   container process started." This exact "started" vs. "actually
   ready" distinction reappears, in a more powerful form, as
   Kubernetes' readiness probes (Chapter 60) -- worth remembering this
   Compose gotcha as the direct conceptual ancestor of that later,
   more robust mechanism.
```

### Try it yourself

```bash
mkdir -p ~/lab/composedemo && cd ~/lab/composedemo
cat > docker-compose.yml << 'EOF'
services:
  web:
    image: alpine
    command: sh -c "echo waiting; sleep 3; ping -c 2 db"
    depends_on:
      - db
  db:
    image: alpine
    command: sh -c "sleep 300"
EOF

docker compose up
docker compose ps
docker compose logs web
docker compose down
```

### Common mistakes

- **Trusting `depends_on` alone for "the service I depend on is actually
  ready."** It only guarantees START ORDER — build real readiness checks
  (retries, healthchecks) into your application.
- **Forgetting `docker compose down` preserves volumes by default.** Good,
  deliberate safety (Chapter 49) — but surprising the first time if you
  expected a truly clean slate; add `-v` explicitly if you genuinely want
  volumes removed too.
- **Hardcoding secrets directly in `docker-compose.yml`.** Use an `.env`
  file (excluded from version control) or a proper secrets-management
  approach instead — a committed compose file with real credentials in it
  is a common, real leak.

### Check yourself

1. What does `depends_on` actually guarantee, and what does it NOT
   guarantee?
2. What happens to volumes when you run `docker compose down` (without
   `-v`)?
3. What's the more robust, production-relevant fix for the "started vs.
   ready" problem, and what Kubernetes mechanism does it foreshadow?

### Further reading

- **Docs:** Docker's own official Compose specification and "Control startup
  and shutdown order" documentation.
- **Reference:** `docker compose --help`.
- **Practice:** convert a small, real two-service idea (a web app + a
  database) into your own `docker-compose.yml` from scratch.

---

## Chapter 52 — Registries and image security basics

### In one sentence

A registry (Docker Hub being the most common default) stores and serves
container images, and because ANYONE can publish to most public registries,
treating every pulled image as untrusted-until-verified is a real, necessary
security practice, not paranoia.

### Why this matters

Pulling and running a container image is, functionally, downloading and
executing someone else's code with real privileges on your machine — the
same supply-chain scrutiny you'd apply to any dependency applies here.

### How it actually works

```bash
   docker pull nginx                  # pulls from Docker Hub by
                                     # default (implicitly:
                                     # docker.io/library/nginx:latest)
   docker pull ghcr.io/someorg/someimage:v1.2.3   # a DIFFERENT
                                     # registry, explicitly -- GitHub
                                     # Container Registry, in this
                                     # example
   docker tag myimage myregistry.example.com/myimage:v1   # tag an
                                     # image for a specific registry
   docker push myregistry.example.com/myimage:v1     # publish it
                                     # (needs `docker login` first)

   PRACTICES WORTH ADOPTING AS DEFAULT:
     [ ] PIN images by a specific TAG (or, more strongly, by DIGEST --
         `nginx@sha256:abc123...`) rather than `latest` -- reproducible,
         and immune to the tag being silently repointed later (this
         wiki's advanced security guide covers exactly this "mutable
         tag" risk in its supply-chain chapters).
     [ ] SCAN images for known vulnerabilities before running them in
         anything beyond a quick local test -- tools like `docker
         scout`, Trivy, or Grype check an image's installed packages
         against known-CVE databases.
     [ ] PREFER official/verified publishers on public registries
         (Docker Hub marks "Official Images" and "Verified Publisher"
         images distinctly) over an unknown, unverified account's image
         claiming to be the same thing.
     [ ] for anything running in production, prefer a PRIVATE registry
         (or a pull-through proxy/mirror) you control, rather than
         pulling directly from the public internet on every deploy.
```

### Real-world example

```
   Security researchers (multiple firms, including well-documented
   public research from Palo Alto Networks' Unit 42 and from Sysdig,
   among others, published across several separate studies over the
   years) have REPEATEDLY found MALICIOUS IMAGES living on Docker Hub's
   public registry -- images designed to look like legitimate,
   popular software (via names deliberately similar to well-known,
   real projects -- a "typosquatting" pattern directly analogous to the
   malicious-package problem covered in this wiki's security guides'
   supply-chain material) but actually containing CRYPTOMINING malware,
   or other unwanted payloads, that activate once the image is pulled
   and run.
   Because ANYONE can create a Docker Hub account and publish public
   images with essentially any name they choose, this is a persistent,
   REAL, ongoing risk category -- not a one-time, patched incident. The
   practical, industry-standard defence, reflected directly in this
   chapter's practice list: don't blindly `docker pull` and run
   whatever image a quick search or an unofficial-looking tutorial
   suggests; check for "Official Image"/"Verified Publisher" status,
   check the publisher's actual identity and reputation, and -- for
   anything beyond throwaway local experimentation -- scan the image
   before running it, exactly the way you'd scrutinise an unfamiliar
   third-party code dependency before adding it to a real project.
```

### Try it yourself

```bash
# See the difference between an official image and pinning by digest
docker pull nginx:latest
docker inspect nginx:latest --format '{{.RepoDigests}}'   # the actual,
                                                          # immutable
                                                          # digest
                                                          # underneath
                                                          # the mutable
                                                          # "latest" tag

# Pull by that exact digest instead of a mutable tag (copy the digest
# from the output above)
# docker pull nginx@sha256:<the digest you saw above>

# Scan an image, if you have a scanner installed (Docker Desktop
# includes `docker scout` by default on many installs)
docker scout quickview nginx:latest 2>/dev/null || \
  echo "(docker scout not available -- try installing Trivy:
  https://aquasecurity.github.io/trivy/ for an open-source alternative)"
```

### Common mistakes

- **Pulling images by `latest` (or any mutable tag) and trusting it stays
  the same.** It can be silently repointed to different content later — pin
  by digest for anything that matters.
- **Assuming a popular-sounding image name is the real, official one.**
  Typosquatting on registries is a real, documented, ongoing risk — verify
  the publisher.
- **Never scanning images before running them in anything beyond a quick,
  fully throwaway local test.** A five-second scan is cheap insurance
  against a genuinely real risk category.

### Check yourself

1. Why is pinning an image by digest stronger than pinning by tag?
2. What is registry typosquatting, and why does it work?
3. What does "Official Image"/"Verified Publisher" status on Docker Hub
   actually tell you?

### Further reading

- **Docs:** Docker Hub's own documentation on "Official Images" and
  "Verified Publishers" — what the badges actually mean and verify.
- **Tool:** Trivy (`aquasecurity.github.io/trivy`) — a free, widely-used,
  open-source image vulnerability scanner.
- **Research:** search "Docker Hub malicious images cryptomining research"
  for the real, published security-research findings on this topic.
- **Guide:** this wiki's advanced security guide covers image signing,
  provenance, and admission-time verification (going well beyond scanning
  alone) in full operational depth.

---

# Part 16 — Kubernetes fundamentals

Docker (Part 15) runs one container well, on one machine. This final part is
the problem that remains once you have hundreds of containers spread across
many machines — and why Kubernetes, rather than a bigger shell script built
on Part 12's material, is what actually fills that gap.

## Chapter 53 — Why Kubernetes exists (the problem Docker alone doesn't solve)

### In one sentence

Docker runs ONE container on ONE machine well; Kubernetes exists to answer
the much harder question — across a FLEET of machines, which containers run
where, what happens when a machine dies, how do they find each other, and
how does the whole thing scale — automatically, continuously, without a
human manually running `docker run` on the right server every time.

### Why this matters

Understanding WHY Kubernetes exists, not just its commands, is what lets you
reason about whether you actually need it, and what problem each of its
pieces (Parts 55-58) is specifically solving.

### How it actually works

```
   RUNNING ONE CONTAINER, on ONE machine, with `docker run`, you already
   fully understand (Part 15). Real production systems need FAR more:

     - MANY containers, across MANY machines (a fleet), because one
       machine's capacity isn't enough, and because a single machine is
       a single point of failure.
     - if a MACHINE DIES, the containers that were running on it need to
       be RESTARTED somewhere else, AUTOMATICALLY -- nobody wants a
       human paged at 3am to manually run `docker run` on a replacement
       server.
     - containers need to FIND EACH OTHER reliably, even as individual
       instances come and go, are rescheduled, or scale up and down
       (Chapter 56).
     - the whole fleet needs to SCALE -- more instances of a busy
       service during high traffic, fewer during quiet periods
       (Chapter 58) -- again, ideally without a human manually
       intervening every time.
     - deploying a NEW VERSION of an application across many running
       instances needs to happen WITHOUT DOWNTIME, and be safely
       REVERSIBLE if something's wrong (Chapter 60).

   KUBERNETES is, fundamentally, a system that takes a DECLARATIVE
   DESCRIPTION of what you want ("I want 5 copies of this container
   running, with these resource limits, exposed on this port") and
   CONTINUOUSLY, AUTOMATICALLY works to make reality match that
   description -- restarting failed instances, rescheduling them onto
   healthy machines, routing traffic to whichever instances are
   currently healthy -- all without a human re-issuing the request
   every time something changes. This "reconcile actual state toward
   desired state, continuously, forever" idea is the single most
   important concept to carry into every chapter that follows.
```

### Real-world example

```
   Kubernetes is not a from-scratch invention -- its own creators have
   been explicit, in a well-documented, published 2016 paper ("Borg,
   Omega, and Kubernetes: Lessons Learned from Three Container-
   Management Systems over a Decade," written by Google engineers and
   published through the ACM) about its direct lineage: Kubernetes was
   built by engineers who had spent YEARS operating Google's INTERNAL
   cluster-management systems -- Borg (which has run essentially all of
   Google's own production workloads, including Search and Gmail, since
   the mid-2000s) and its successor Omega -- and Kubernetes was
   deliberately designed to bring the SAME hard-won lessons from running
   containers at that massive, real, internal scale to the wider
   public, as an open-source project (Kubernetes was announced by
   Google in 2014).
   The paper itself is genuinely worth reading: it candidly describes
   both what worked and what Google's own engineers considered MISTAKES
   in Borg's original design, and how Kubernetes deliberately corrected
   for them (the shift toward Kubernetes' now-familiar declarative, API-
   driven model was, per the paper, a direct response to specific,
   real, operational pain points Google had experienced running Borg at
   scale for over a decade). This lineage is precisely why Kubernetes'
   design feels the way it does throughout Parts 16-17: it's the
   distilled, public version of one of the largest, longest-running
   real-world experiments in "how do you actually run containers
   reliably at scale" in the industry's history.
```

### Try it yourself

```bash
# You'll need a local cluster for the rest of this Part -- set one up
# now (kind = "Kubernetes IN Docker," runs a real cluster using
# containers as the "nodes," genuinely excellent for learning)
# Install kind: https://kind.sigs.k8s.io/docs/user/quick-start/
# Install kubectl: https://kubernetes.io/docs/tasks/tools/

kind create cluster --name learning
kubectl cluster-info
kubectl get nodes                  # your single-node "cluster" -- but
                                   # a REAL Kubernetes API server,
                                   # scheduler, and kubelet are all
                                   # genuinely running
kubectl get namespaces
```

### Common mistakes

- **Reaching for Kubernetes for a single-server hobby project.** It solves
  fleet-scale problems — genuinely excellent tooling, genuinely unnecessary
  complexity for something `docker compose` (Chapter 51) already handles
  well.
- **Thinking Kubernetes replaces Docker.** Kubernetes ORCHESTRATES
  containers across many machines; something still has to actually RUN
  them on each machine (a container runtime — historically often Docker's
  own runtime, now commonly `containerd` directly).
- **Assuming "declarative" means "instant."** Reconciliation toward desired
  state happens continuously, but it's not instantaneous — there's real,
  observable time between "I declared 5 replicas" and "5 replicas are
  actually running and healthy" (Chapter 59 covers watching this happen).

### Check yourself

1. What specific problems does Kubernetes solve that plain Docker, running
   on one machine, does not?
2. What is the "reconcile actual state toward desired state" idea, and why
   is it central to everything Kubernetes does?
3. What real, internal Google system directly inspired Kubernetes' design?

*(Answers: Appendix E.)*

### Further reading

- **Paper:** "Borg, Omega, and Kubernetes" (Google, 2016, published via
  ACM Queue) — genuinely readable, genuinely worth reading in full.
- **Docs:** Kubernetes' own official documentation, "Kubernetes Components"
  overview page — the authoritative starting point.
- **Tool:** `kind` (`kind.sigs.k8s.io`) — the local-cluster tool used
  throughout the rest of this Part.

---

## Chapter 54 — The control plane and the node: Kubernetes' architecture

### In one sentence

Every Kubernetes cluster splits into the **control plane** (the brain —
decides what SHOULD be running, where) and one or more **nodes** (the
muscle — actually RUN the containers), and this split is exactly why a
cluster can survive real, partial failures gracefully.

### Why this matters

"Why did my app keep running even though the control plane had an issue"
and "why doesn't the API server itself run my containers" are both answered
by this one architectural split.

### How it actually works

```
   THE CONTROL PLANE (the "brain," typically running on dedicated
   control-plane nodes):
     API SERVER          the front door -- EVERYTHING (kubectl, every
                        other control-plane component, every kubelet)
                        talks to Kubernetes ONLY through this REST API.
                        Nothing bypasses it.
     etcd                  a distributed, consistent key-value store --
                        the SOURCE OF TRUTH for the ENTIRE cluster's
                        state (every object's desired configuration,
                        current status).
     SCHEDULER              watches for newly-created pods with no node
                        assigned yet, and DECIDES which node they should
                        run on (based on resource availability,
                        constraints -- Chapter 58 covers this properly).
     CONTROLLER MANAGER      runs the actual "reconcile actual state
                        toward desired state" LOOPS (Chapter 53's core
                        idea) -- a Deployment controller noticing "I
                        should have 5 replicas but only 3 are running"
                        and creating 2 more, for example.

   THE NODE (the "muscle," where your actual containers run -- a
   cluster typically has many of these):
     KUBELET                 the agent running on EVERY node -- talks
                        to the API server, and makes sure the
                        containers THIS node has been assigned are
                        actually running, restarting them locally if
                        they crash (Chapter 33's systemd `Restart=`
                        idea, one layer up).
     CONTAINER RUNTIME         the actual thing running containers on
                        this node -- `containerd` (very commonly, these
                        days) or another compatible runtime -- doing
                        the literal Chapter 44/47 work.
     KUBE-PROXY                handles the networking rules that make
                        Services (Chapter 56) actually route traffic to
                        the right pods on this node.

   THE ARCHITECTURAL PAYOFF: because the CONTROL PLANE decides "what
   should run" and NODES independently, locally, keep "what's actually
   running" matching their OWN last-known assignment (via the kubelet),
   a TEMPORARY control-plane disruption does NOT immediately stop
   already-scheduled, already-running pods from continuing to run and
   serve traffic on healthy nodes -- the kubelet keeps enforcing the
   last instruction it received, even if it can't currently reach the
   API server for a NEW one. This graceful-degradation property is a
   deliberate, well-documented design characteristic, not an accident.
```

### Real-world example

```
   This exact resilience property -- already-running workloads
   surviving a control-plane disruption -- is explicitly, deliberately
   documented in Kubernetes' own official architecture documentation as
   intentional design, not a lucky side effect: the kubelet is built to
   operate with a degree of AUTONOMY specifically so that a control-
   plane issue (a temporary API-server outage, a networking partition
   between a node and the control plane, a slow etcd) doesn't
   IMMEDIATELY translate into an application-facing outage for workloads
   that are ALREADY scheduled and running. Real operators managing
   Kubernetes at scale rely on this property directly: it changes how
   URGENTLY a control-plane issue needs to be treated (still urgent --
   NEW deployments, scaling, and self-healing from node failures all
   stop working until it's fixed -- but not an instant, cluster-wide
   application outage the moment it happens), and it's a genuinely
   important distinction for anyone doing on-call for a Kubernetes-based
   system to understand precisely, rather than treating "control plane
   is unhealthy" and "my application is down" as automatically the same
   incident.
```

### Try it yourself

```bash
# See the control plane and node(s) in your own cluster
kubectl get nodes -o wide
kubectl get pods -n kube-system     # the control-plane components
                                    # themselves ARE pods (on a `kind`
                                    # cluster, running as containers) --
                                    # a genuinely nice, concrete
                                    # illustration that Kubernetes
                                    # itself is largely "Kubernetes,
                                    # bootstrapped, running Kubernetes"

kubectl get pods -n kube-system -o wide | grep -E "etcd|apiserver|scheduler|controller"

# See the API server directly (everything really does go through it)
kubectl get --raw /healthz
kubectl proxy &                      # expose the API server locally
sleep 1
curl -s localhost:8001/version       # a raw HTTP call to the API
                                     # server -- exactly what `kubectl`
                                     # itself does, under the hood
kill %1
```

### Common mistakes

- **Assuming the control plane runs your application containers.** It
  DECIDES what should run — the actual container processes run on nodes,
  managed by each node's kubelet.
- **Treating "control plane down" and "application down" as automatically
  the same incident.** Already-running workloads often keep serving
  traffic — verify what's ACTUALLY affected before escalating as a full
  outage.
- **Forgetting etcd is the actual source of truth.** Every other component
  is, in a sense, working FROM what's recorded in etcd — its health is
  foundational to the whole cluster's.

### Check yourself

1. What's the difference between the control plane and a node, in one
   sentence each?
2. What does the kubelet actually do, and why does it matter for
   resilience during a control-plane issue?
3. Why does "the control plane had an issue" not automatically mean "the
   application was down"?

### Further reading

- **Docs:** Kubernetes' own official "Kubernetes Components" documentation
  page — the authoritative architecture reference.
- **Reference:** `kubectl get pods -n kube-system` on any real cluster —
  literally see the architecture running.
- **Paper:** the Borg/Omega/Kubernetes paper (Chapter 53) covers the
  control-plane/node split's design rationale directly.

---

## Chapter 55 — Pods, Deployments, and ReplicaSets

### In one sentence

A **Pod** is the smallest deployable unit (one or more tightly-coupled
containers sharing a network namespace), a **ReplicaSet** keeps a specified
NUMBER of identical pods running, and a **Deployment** manages ReplicaSets
over time — enabling safe, versioned rollouts and rollbacks.

### Why this matters

You will write, read, and debug Deployment manifests constantly in any real
Kubernetes role — understanding the three-layer relationship (Deployment
manages ReplicaSets manages Pods) explains a huge amount of `kubectl`
output that otherwise looks redundant or confusing.

### How it actually works

```yaml
   # a minimal, real Deployment manifest
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: web-app
   spec:
     replicas: 3                     # I want THREE identical pods
     selector:
       matchLabels:
         app: web-app
     template:                        # the POD TEMPLATE -- what each
                                     # replica actually looks like
       metadata:
         labels:
           app: web-app
       spec:
         containers:
         - name: web
           image: nginx:1.25
           ports:
           - containerPort: 80
```

```
   THE THREE LAYERS, and why each exists:
     POD          one or more containers that share a network namespace
                (Chapter 44) and are always scheduled TOGETHER, on the
                SAME node. Pods are treated as EPHEMERAL and DISPOSABLE
                -- you almost never create one directly in real use;
                you describe what you want at the Deployment level and
                let Kubernetes create/destroy Pods as needed.
     REPLICASET    ensures a SPECIFIED NUMBER of pods matching a
                selector are running, AT ALL TIMES -- if one crashes or
                is deleted, the ReplicaSet's controller (Chapter 54)
                notices and creates a replacement.
     DEPLOYMENT     manages ReplicaSets OVER TIME, enabling ROLLING
                UPDATES (Chapter 60): when you change the pod template
                (a new image version, say), the Deployment creates a
                NEW ReplicaSet with the new template, gradually scales
                it UP while scaling the OLD ReplicaSet DOWN, and keeps
                the old ReplicaSet around (scaled to zero) so a ROLLBACK
                is fast if needed.
```

```bash
   kubectl apply -f deployment.yaml     # create or UPDATE, declaratively
   kubectl get deployments
   kubectl get replicasets               # you'll see one per version
                                        # you've ever deployed (kept
                                        # around, scaled to 0, for
                                        # rollback -- up to a
                                        # configurable history limit)
   kubectl get pods                       # the actual, currently
                                        # running containers
   kubectl scale deployment web-app --replicas=5   # imperatively change
                                        # the desired count
```

### Real-world example

```
   Niantic's Pokemon GO launch in July 2016 is a well-documented, real,
   dramatic case study in why AUTOMATIC, RELIABLE SCALING (exactly what
   Deployments and ReplicaSets provide the foundation for) matters:
   Google Cloud published its own technical account of the event,
   describing that actual player demand reached roughly FIFTY TIMES
   Niantic's original capacity planning estimates within the first
   weeks -- a genuinely enormous, largely unpredictable real spike in
   traffic, hitting infrastructure running on Google's cloud platform.
   The specific, published account focuses on Google's broader cloud
   infrastructure and its own internal orchestration systems (the
   direct Borg/Kubernetes lineage from Chapter 53) scaling the backend
   services to absorb genuinely unprecedented, unplanned real-world
   demand, with Niantic's own engineering team working closely with
   Google's infrastructure teams during the event.
   WHY THIS BELONGS HERE, regardless of the exact underlying
   orchestration technology at the time: it's a vivid, real, public
   illustration of EXACTLY the problem this chapter's three-layer model
   is built to solve -- automatically maintaining (and rapidly
   increasing) a correct number of healthy service instances under
   real, unpredictable, extreme demand, without every single scaling
   decision requiring a human to manually intervene. This is precisely
   the scenario a properly-configured Deployment (Chapter 58's
   autoscaling, layered on top) is designed to handle gracefully, and
   it's a genuinely useful real story to have in mind for WHY this
   chapter's abstractions exist, not just what their YAML syntax looks
   like.
```

### Try it yourself

```bash
# Deploy something real and watch the three layers
cat > deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
      - name: web
        image: nginx:1.25
        ports:
        - containerPort: 80
EOF

kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get replicasets
kubectl get pods -l app=web-app

# Kill a pod and watch the ReplicaSet replace it automatically
POD_TO_DELETE=$(kubectl get pods -l app=web-app -o jsonpath='{.items[0].metadata.name}')
test -n "$POD_TO_DELETE" && kubectl delete pod "$POD_TO_DELETE"
sleep 3
kubectl get pods -l app=web-app       # still 3 -- a new one appeared

# Update the image and watch a rolling update happen
kubectl set image deployment/web-app web=nginx:1.26
kubectl rollout status deployment/web-app
kubectl get replicasets               # TWO ReplicaSets now -- old
                                      # (scaled to 0) and new

kubectl delete -f deployment.yaml
```

### Common mistakes

- **Creating bare Pods directly for anything meant to be reliable.** No
  ReplicaSet/Deployment means nothing recreates it if it crashes or its node
  fails — use a Deployment (or another appropriate controller) for anything
  beyond a one-off debugging session.
- **Being confused by multiple ReplicaSets for one Deployment.** This is
  normal and intentional — old ones (scaled to 0) are kept for fast
  rollback, up to a configurable history limit.
- **Editing a Pod directly to "fix" something.** Pods created by a
  Deployment are managed BY it — a direct edit will likely be overwritten
  the next time the controller reconciles. Edit the Deployment instead.

### Check yourself

1. What does a Pod represent, and why is it usually NOT created directly?
2. What's the relationship between a Deployment, a ReplicaSet, and a Pod?
3. Why does `kubectl get replicasets` often show more than one entry for a
   single Deployment?

### Further reading

- **Docs:** Kubernetes' own official "Pods," "ReplicaSet," and "Deployments"
  documentation pages — the authoritative reference for every field.
- **Case study:** search "Google Cloud Pokemon Go case study" for the
  original, published technical account.
- **Reference:** `kubectl explain deployment.spec` — live, built-in
  documentation for any field, straight from your cluster.

---

## Chapter 56 — Services and how traffic finds a pod

### In one sentence

Pods are ephemeral — they're created, destroyed, and rescheduled with new IP
addresses constantly — so a **Service** provides a STABLE address and name
that automatically, continuously routes traffic to whichever healthy pods
currently match a label selector.

### Why this matters

"How does anything reliably talk to my constantly-changing pods" is one of
the first genuinely confusing questions for anyone new to Kubernetes, and
Services are the entire, elegant answer.

### How it actually works

```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: web-app-svc
   spec:
     selector:
       app: web-app              # matches the SAME label from Chapter
                                # 55's Deployment's pod template
     ports:
     - port: 80                   # the Service's OWN port
       targetPort: 80              # the port on the PODS to send
                                # traffic to
     type: ClusterIP               # the default -- internal-only
                                # (Chapter 37's networking, one layer
                                # up)
```

```
   THE CORE MECHANISM: a Service doesn't point at specific PODS by
   name or IP -- it uses a LABEL SELECTOR, continuously, live. Whatever
   pods CURRENTLY have matching labels are automatically included as
   valid destinations -- as pods are created, destroyed, or rescheduled
   (each getting a NEW IP, Chapter 55), the Service's routing updates
   automatically, with no manual reconfiguration ever needed. This is
   EXACTLY Docker Compose's container-name DNS resolution (Chapter 50),
   but far more dynamic -- resolving to a live, continuously-updated SET
   of healthy pods, not one fixed container.

   SERVICE TYPES, briefly:
     ClusterIP (default)   reachable only FROM WITHIN the cluster --
                          the normal way services talk to EACH OTHER.
     NodePort                 additionally exposes the service on a
                          specific port on EVERY node's own IP --
                          simple, but rarely the final answer for real
                          external traffic.
     LoadBalancer              (on a cloud provider) provisions an
                          actual, external cloud load balancer pointing
                          at the service -- the standard way to expose
                          something to the real internet.

   EVERY SERVICE also gets a DNS NAME, automatically, resolvable from
   anywhere in the cluster: `web-app-svc.default.svc.cluster.local`
   (or just `web-app-svc` from within the same namespace) -- Chapter
   37's DNS resolution, provided automatically by Kubernetes' own
   internal DNS (CoreDNS, typically).
```

### Real-world example

```
   The PROBLEM Services solve is explicit, foundational, and explained
   directly in Kubernetes' own official documentation: pods are, by
   design, EPHEMERAL. A pod can be deleted and recreated (a
   Deployment's rolling update, Chapter 60; a node failure and
   reschedule; a manual scale-down-then-up) and the REPLACEMENT pod
   gets a COMPLETELY NEW IP ADDRESS every time. Any system that tried
   to talk to "the web server" by remembering a specific pod's IP
   address would break constantly, unpredictably, every single time a
   pod was replaced for ANY reason -- a genuinely unworkable model at
   any real scale, and precisely the problem Kubernetes' designers
   (with the Borg lineage's years of real operational experience behind
   them, Chapter 53) built the Service abstraction specifically to
   solve from day one, not as a later addition.
   This is worth sitting with as a GENERAL PRINCIPLE that recurs
   throughout distributed systems, not just Kubernetes specifically:
   whenever INDIVIDUAL INSTANCES of something are expected to come and
   go, you need a STABLE LAYER OF INDIRECTION (a name, a label selector,
   a load balancer) in front of them -- pointing directly at an
   individual, ephemeral instance is a reliability bug waiting to
   happen, in Kubernetes or in any other real distributed system.
```

### Try it yourself

```bash
# Continue from Chapter 55's deployment (recreate if needed)
kubectl apply -f deployment.yaml

cat > service.yaml << 'EOF'
apiVersion: v1
kind: Service
metadata:
  name: web-app-svc
spec:
  selector:
    app: web-app
  ports:
  - port: 80
    targetPort: 80
EOF
kubectl apply -f service.yaml
kubectl get svc web-app-svc

# See the Service resolve to MULTIPLE pod IPs, live
kubectl get endpoints web-app-svc     # every pod IP currently backing
                                      # this service

# Prove it survives a pod being replaced
OLDPOD=$(kubectl get pods -l app=web-app -o jsonpath='{.items[0].metadata.name}')
test -n "$OLDPOD" && kubectl delete pod "$OLDPOD"
sleep 3
kubectl get endpoints web-app-svc      # a DIFFERENT set of IPs now --
                                       # the OLD pod's IP is gone, a
                                       # NEW pod's IP is present -- but
                                       # the SERVICE's own address/name
                                       # never changed at all

# Reach it from inside the cluster (DNS resolution in action)
kubectl run testcurl --image=curlimages/curl --rm -it --restart=Never -- \
  curl -s http://web-app-svc

kubectl delete -f service.yaml -f deployment.yaml
```

### Common mistakes

- **Hardcoding a specific pod's IP anywhere.** It WILL change — always go
  through a Service (or its DNS name) instead.
- **Forgetting the Service's `selector` must actually match the pod
  template's labels.** A typo here silently results in a Service with ZERO
  healthy endpoints — always check `kubectl get endpoints` when traffic
  isn't reaching pods as expected.
- **Using `NodePort` as the default choice for external traffic.** It works,
  but `LoadBalancer` (on a cloud provider) or an Ingress controller (Chapter
  59-adjacent territory) is almost always the better real-world answer for
  actual internet-facing traffic.

### Check yourself

1. What problem does a Service solve, specifically?
2. How does a Service know which pods to route traffic to, and does that
   update automatically?
3. What's the difference between `ClusterIP`, `NodePort`, and
   `LoadBalancer`?

### Further reading

- **Docs:** Kubernetes' own official "Service" documentation — thorough,
  and the authoritative reference for every service type and field.
- **Reference:** `kubectl explain service.spec`, `kubectl get endpoints`.
- **Concept:** search "stable indirection ephemeral instances distributed
  systems" for the broader pattern this chapter's lesson generalises into.

---

## Chapter 57 — Configuration: ConfigMaps, Secrets, and environment

### In one sentence

**ConfigMaps** hold non-sensitive configuration and **Secrets** hold
sensitive values, both stored separately from your container image so the
SAME image can run with different configuration in different environments —
but Secrets being base64-ENCODED, not encrypted, by default is a real,
commonly misunderstood security detail.

### Why this matters

Baking configuration (especially credentials) directly into a container
image means rebuilding the image for every environment, and — worse — that
image now contains a secret anyone who can pull it can read (Chapter 52).

### How it actually works

```yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: app-config
   data:
     LOG_LEVEL: "info"
     API_TIMEOUT: "30s"

   ---
   apiVersion: v1
   kind: Secret
   metadata:
     name: app-secret
   type: Opaque
   stringData:                      # plaintext HERE, in the manifest --
     DB_PASSWORD: "supersecret123"    # Kubernetes base64-ENCODES it on
                                    # storage (NOT encryption, Chapter
                                    # 6 of the advanced security guide
                                    # in this wiki covers this
                                    # precisely)
```

```yaml
   # using both in a pod spec
   spec:
     containers:
     - name: app
       image: myapp:1.0
       envFrom:
       - configMapRef:
           name: app-config           # every key becomes an env var
       env:
       - name: DB_PASSWORD
         valueFrom:
           secretKeyRef:
             name: app-secret
             key: DB_PASSWORD
```

```bash
   kubectl create configmap app-config --from-literal=LOG_LEVEL=info
   kubectl get configmap app-config -o yaml
   kubectl create secret generic app-secret --from-literal=DB_PASSWORD=supersecret123
   kubectl get secret app-secret -o jsonpath='{.data.DB_PASSWORD}' | base64 -d
      # ^ trivially decoded -- this is exactly the point: it's NOT
      #   encryption, it's encoding, exactly G1's Chapter 6 lesson
      #   ("Base64 is an envelope, not a safe"), applied directly to
      #   Kubernetes Secrets
```

### Real-world example

```
   This wiki's own advanced security guide (`security/real-life-
   security-guide-v1.md`, Part 4 Chapter 16) documents this EXACT,
   genuinely common, real misunderstanding at length: Kubernetes
   Secrets are, BY DEFAULT, stored in etcd (Chapter 54's source of
   truth) as PLAIN BASE64 -- not encrypted -- unless the cluster
   operator EXPLICITLY configures encryption-at-rest with a KMS/AESGCM
   provider. Anyone who can read etcd directly (or an etcd BACKUP, or a
   snapshot in an over-permissioned storage bucket) gets every "Secret"
   in the cluster in effectively plaintext, with one trivial `base64
   -d` step.
   This is a genuinely, repeatedly real source of security incidents
   and misconfigurations: engineers reasonably ASSUME "Secret," as a
   Kubernetes object TYPE, implies real encryption, precisely because
   the name strongly suggests it -- and the base64 encoding step, which
   IS applied, LOOKS superficially like protection to someone who
   doesn't already know Chapter 6's "encoding is not encryption"
   distinction cold. The practical, real fix -- covered in full
   operational depth in this wiki's advanced security guide -- is
   configuring genuine etcd encryption-at-rest, AND treating Kubernetes
   Secrets as one layer of a real secrets-management strategy (ideally
   backed by a proper external secrets manager, G1 Chapter 49) rather
   than the whole story on their own.
```

### Try it yourself

```bash
kubectl create configmap demo-config --from-literal=GREETING="hello from a ConfigMap"
kubectl create secret generic demo-secret --from-literal=PASSWORD="not-really-secret-until-encrypted"

# See the "encoding, not encryption" lesson directly
kubectl get secret demo-secret -o jsonpath='{.data.PASSWORD}'
echo    # (that's the base64-encoded value)
kubectl get secret demo-secret -o jsonpath='{.data.PASSWORD}' | base64 -d
echo    # trivially recovered -- exactly this chapter's point

cat > configpod.yaml << 'EOF'
apiVersion: v1
kind: Pod
metadata:
  name: config-demo
spec:
  containers:
  - name: demo
    image: alpine
    command: ["sh", "-c", "echo GREETING=$GREETING; echo PASSWORD=$PASSWORD; sleep 300"]
    env:
    - name: GREETING
      valueFrom:
        configMapKeyRef:
          name: demo-config
          key: GREETING
    - name: PASSWORD
      valueFrom:
        secretKeyRef:
          name: demo-secret
          key: PASSWORD
EOF
kubectl apply -f configpod.yaml
sleep 2
kubectl logs config-demo

kubectl delete -f configpod.yaml
kubectl delete configmap demo-config
kubectl delete secret demo-secret
```

### Common mistakes

- **Treating a Kubernetes Secret as inherently encrypted.** It's base64
  ENCODED by default — genuinely, commonly misunderstood, and a real
  security gap unless you explicitly configure etcd encryption-at-rest.
- **Baking configuration or credentials directly into a container image.**
  Means rebuilding the image per environment, and a credential leak the
  moment anyone pulls the image (Chapter 52).
- **Putting genuinely sensitive values in a ConfigMap instead of a
  Secret.** ConfigMaps get NO special treatment at all — always use a
  Secret (and real encryption-at-rest) for anything sensitive.

### Check yourself

1. What's the practical difference between a ConfigMap and a Secret, given
   that both are, by default, unencrypted at rest?
2. Why is "Kubernetes Secrets are base64, not encrypted" such a commonly
   real, commonly misunderstood point?
3. What's the actual, complete fix for storing genuinely sensitive values
   safely in Kubernetes?

### Further reading

- **Docs:** Kubernetes' own official "Secrets" documentation, specifically
  the "Risks" section — refreshingly honest about this exact limitation.
- **Guide:** this wiki's advanced security guide, Part 4 Chapter 16
  ("Kubernetes architecture and its trust boundaries") — the full, real
  operational treatment of this exact topic.
- **Reference:** `kubectl explain secret`, `kubectl explain configmap`.

---

## Chapter 58 — Scheduling and scaling

### In one sentence

The **scheduler** decides which node a new pod runs on (based on available
resources and constraints you can specify), and **autoscaling** — both of
pods (more replicas) and of nodes (more machines) — is how a cluster
responds to changing real demand without a human manually intervening every
time.

### Why this matters

Understanding scheduling explains why a pod sometimes stays stuck
`Pending`, and autoscaling is the concrete, practical payoff of everything
this Part has built toward: a system that genuinely handles real,
unpredictable load on its own.

### How it actually works

```yaml
   # resource REQUESTS and LIMITS -- the scheduler's primary inputs
   resources:
     requests:                # "I need AT LEAST this much" -- the
       cpu: "250m"              # SCHEDULER uses this to decide which
       memory: "256Mi"           # node has room
     limits:                   # "never let me use MORE than this" --
       cpu: "500m"                # enforced by the KERNEL via cgroups
       memory: "512Mi"             # (Chapter 43!) -- exceed the memory
                                 # limit and the OOM killer (Chapter 13)
                                 # acts, scoped to this pod's cgroup,
                                 # exactly as Chapter 44 described
```

```
   THE SCHEDULER'S JOB, each time a new pod needs a node: look at every
   node's AVAILABLE capacity (total minus what's already requested by
   other pods there), filter out nodes that don't meet any
   CONSTRAINTS (taints/tolerations, node selectors, affinity rules --
   ways to say "only run on nodes with this label" or "never run two of
   these on the same node"), and pick a node from what's left (using a
   scoring algorithm to balance the cluster reasonably).

   IF NO NODE HAS ENOUGH ROOM: the pod stays in `Pending` state
   (Chapter 59 covers diagnosing this) -- it's not an error exactly,
   it's the scheduler correctly reporting "I have nowhere to safely put
   this yet."

   AUTOSCALING, two distinct, complementary layers:
     HORIZONTAL POD AUTOSCALER (HPA)   watches a metric (commonly CPU
       or memory utilisation, or a custom metric) and automatically
       adjusts a Deployment's REPLICA COUNT within a range you define --
       more pods when busy, fewer when quiet. This is Chapter 55's
       ReplicaSet mechanism, driven automatically instead of manually.
     CLUSTER AUTOSCALER (on a cloud provider)   watches for pods stuck
       `Pending` because NO EXISTING NODE has room, and automatically
       provisions ADDITIONAL NODES (real cloud instances) to make room --
       then, just as importantly, removes nodes again once they're no
       longer needed, to control cost.
   TOGETHER: HPA scales PODS to match demand; the cluster autoscaler
   scales the underlying MACHINES to make sure there's actually room for
   however many pods HPA has decided are needed.
```

### Real-world example

```
   Handling large, predictable seasonal traffic spikes -- the Black
   Friday/Cyber Monday retail shopping period being the most widely-
   discussed, most publicly-documented recurring real example across
   the industry -- is one of the most common, real-world-validated use
   cases explicitly cited for exactly this chapter's two-layer
   autoscaling model. Many retailers and e-commerce platforms have
   publicly published engineering blog posts and conference talks
   describing their own Kubernetes-based autoscaling architecture
   specifically built to handle this kind of extreme, well-understood-
   in-advance-but-still-massive traffic multiplier -- scaling from
   normal baseline capacity up to many times that, automatically, for a
   short, intense period, then scaling back down once demand subsides,
   without needing to PERMANENTLY provision (and pay for) peak-level
   infrastructure year-round.
   This is precisely the real, practical payoff this entire Part has
   been building toward: a Deployment declares the DESIRED shape of an
   application (Chapter 55), a Service gives it a stable address
   (Chapter 56), configuration is externalised cleanly (Chapter 57),
   and autoscaling makes the whole thing RESPOND to real demand
   automatically -- the full, connected picture of why an organisation
   takes on Kubernetes' real complexity in the first place.
```

### Try it yourself

```bash
# See scheduling decisions and resource requests/limits directly
cat > scheduled.yaml << 'EOF'
apiVersion: v1
kind: Pod
metadata:
  name: resource-demo
spec:
  containers:
  - name: demo
    image: nginx
    resources:
      requests:
        cpu: "100m"
        memory: "64Mi"
      limits:
        cpu: "200m"
        memory: "128Mi"
EOF
kubectl apply -f scheduled.yaml
kubectl get pod resource-demo -o wide     # which node did it land on?
kubectl describe pod resource-demo | grep -A5 "Requests\|Limits"

# Deliberately request more than your (single-node kind) cluster has,
# and watch it stay Pending
cat > toobig.yaml << 'EOF'
apiVersion: v1
kind: Pod
metadata:
  name: too-big
spec:
  containers:
  - name: demo
    image: nginx
    resources:
      requests:
        cpu: "999"          # deliberately absurd -- no real node has
        memory: "999Gi"     # this much capacity
EOF
kubectl apply -f toobig.yaml
sleep 2
kubectl get pod too-big              # Pending
kubectl describe pod too-big | grep -A3 Events   # the scheduler
                                                  # explains WHY

kubectl delete -f scheduled.yaml -f toobig.yaml
```

### Common mistakes

- **Setting no resource requests/limits at all.** The scheduler has nothing
  to reason about, and a runaway pod (Chapter 13's OOM-killer story, one
  layer up) can starve its neighbours on the same node.
- **Setting `requests` and `limits` to the SAME value for everything,
  reflexively.** Sometimes appropriate (predictable workloads); often overly
  rigid — understand the actual trade-off (`limits` above `requests` allows
  "bursting" into spare capacity) rather than copy-pasting a pattern.
- **Assuming autoscaling reacts instantly.** Both HPA and cluster
  autoscaling take real, observable TIME (metric collection intervals,
  actual cloud-instance provisioning time) — plan for that lag, don't
  assume it's immediate.

### Check yourself

1. What's the difference between a resource `request` and a `limit`?
2. Why does a pod stay `Pending` when no node has enough capacity, and is
   that an error?
3. What's the difference between the Horizontal Pod Autoscaler and a
   Cluster Autoscaler, and why do you typically need both together?

### Further reading

- **Docs:** Kubernetes' own official "Assigning Pods to Nodes,"
  "Horizontal Pod Autoscaling," and "Cluster Autoscaler" documentation.
- **Reference:** `kubectl describe node <name>` — see a real node's actual
  allocated vs. available capacity.
- **Guide:** this wiki's advanced security guide covers multi-tenancy and
  node isolation (Part 4 Chapter 21) as a related, deeper topic once
  scheduling basics are solid.

---

# Part 17 -- Kubernetes in practice

You know the pieces now: Pods, Deployments, Services, ConfigMaps, scheduling.
This part is about the muscle memory of actually operating a cluster --
reading its state when something is wrong, rolling out changes safely, and
packaging applications so other people can install them without reading
your mind.

## Chapter 59 -- Reading cluster state and debugging a pod that won't start

### In one sentence
`kubectl describe`, `kubectl logs`, and `kubectl get events` are the
stethoscope, thermometer, and chart review of Kubernetes debugging, and
almost every "pod won't start" problem falls into one of about six buckets.

### Why this matters
A Deployment says "3 desired, 1 available" and nothing works, and the
instinct of someone new to Kubernetes is to restart everything or delete
the cluster and start over. That instinct is expensive and usually wrong.
Kubernetes tells you exactly what is wrong, in plain English, in the
`Events` section of `kubectl describe pod` -- you just have to know to look
there first, before logs, before dashboards, before Slack.

For a manager: the cost of a production incident is often dominated not by
the failure itself but by the time-to-diagnosis. An engineer who knows the
six buckets below finds the cause in two minutes. An engineer who does not
know them tries five unrelated things over 40 minutes. That gap, multiplied
across every incident a team has in a year, is a real and measurable amount
of engineering time -- which is why "can you debug a broken pod calmly" is
a legitimate, common interview question for SRE and platform roles.

### How it actually works
The debugging order that works, almost every time:

```
   1. kubectl get pods                 -- what STATUS is it in?
   2. kubectl describe pod <name>      -- read the Events at the bottom
   3. kubectl logs <name>              -- what did the app itself say?
   4. kubectl logs <name> --previous   -- what did it say before it
                                          last crashed? (for CrashLoopBackOff)
   5. kubectl get events --sort-by=.metadata.creationTimestamp
                                       -- cluster-wide event timeline
```

The six buckets, by `STATUS` column in `kubectl get pods`:

```
   Pending
     The Pod has not been scheduled to a node yet. Cause is almost
     always in `describe`'s Events: "Insufficient cpu", "Insufficient
     memory", "0/5 nodes are available: 5 node(s) had taint
     {node-role: master} that the pod didn't tolerate", or a
     PersistentVolumeClaim that has no matching PersistentVolume.
     Fix: right-size requests, add tolerations, or fix storage.

   ImagePullBackOff / ErrImagePull
     Kubernetes cannot pull the container image. Almost always a typo
     in the image name/tag, a private registry with no imagePullSecret
     configured, or (surprisingly common) the image tag was deleted or
     overwritten by a later build.
     Fix: `kubectl describe pod` shows the exact registry error.

   CrashLoopBackOff
     The container starts, then exits, repeatedly, with increasing
     backoff delay between restarts. The container image and pull are
     fine -- the *application* is crashing. This is the one where
     `kubectl logs --previous` is essential, because by the time you
     look, the current attempt's logs may be empty (crashed before
     printing anything) while the previous attempt's logs show the
     real stack trace.
     Common causes: missing environment variable / Secret, cannot
     reach a database on startup, misconfigured config file, the
     entrypoint script has a typo.

   Pending forever with "0/N nodes are available: N Insufficient memory"
     This is a scheduling problem, not a crash. Either lower the
     Pod's memory *request*, add nodes, or find and evict a
     resource-hogging Pod that's wasting reserved-but-unused capacity.

   OOMKilled (seen in `kubectl describe pod`, "Last State: Terminated,
   Reason: OOMKilled")
     The container exceeded its memory *limit* and the kernel's OOM
     killer (Chapter 11) killed it. Fix: raise the limit if the
     workload legitimately needs more memory, or fix a memory leak if
     it doesn't.

   Running, but Service traffic still fails (0/1 Ready)
     The container process is alive but the readiness probe is
     failing -- e.g. the app takes 45 seconds to warm up but the
     probe's initialDelaySeconds is 5. Kubernetes is correctly
     refusing to send traffic to a Pod that says it isn't ready;
     the fix is almost always in the probe configuration.
```

### Real-world example
A well-known and often-retold Kubernetes debugging story from the SRE
community involves a service that worked perfectly in staging and then
went into `CrashLoopBackOff` immediately after every production deploy.
`kubectl logs` on the current, crashed container showed nothing -- an
empty log, because the process died within milliseconds of starting,
before its logging library had even finished initializing. Only
`kubectl logs --previous`, combined with `describe pod`'s "Last State"
block, revealed the actual error: the production ConfigMap was missing a
key that staging's ConfigMap happened to have, so the app's config parser
threw an unhandled exception at import time, before any log line could be
written. The lesson generalized into a common piece of team tooling:
runbooks for "Pod won't start" now routinely open with the exact five-step
sequence above, in order, precisely because skipping straight to "check
the logs" misses the crash that happens before there are any logs to
check. Datadog, Google Cloud's SRE documentation, and the official
Kubernetes troubleshooting docs all converge on the same ordered checklist
independently, which is itself a signal that this is the sequence that
actually works in practice rather than a stylistic preference.

### Try it yourself
```bash
# Deliberately break a pod, then debug it using only the sequence above.
kubectl create deployment broken --image=nginx:this-tag-does-not-exist
kubectl get pods
# STATUS should show ImagePullBackOff or ErrImagePull within ~30s

kubectl describe pod -l app=broken | tail -20
# read the Events block -- it names the exact problem

kubectl delete deployment broken

# Now break it a different way: a container that crashes on startup.
kubectl run crasher --image=busybox --restart=Never -- sh -c "exit 1"
kubectl get pod crasher
kubectl logs crasher            # probably empty or minimal
kubectl describe pod crasher | grep -A5 "Last State"
kubectl delete pod crasher
```

### Common mistakes
- Jumping to `kubectl logs` before `kubectl describe`, and missing a
  scheduling or image-pull problem that never produced any application
  logs at all because the container never ran.
- Reading only the *current* logs of a crash-looping container and
  concluding "no errors," when `--previous` would show the actual crash.
- Treating `Pending` as "the cluster is broken" when it usually means
  "this specific Pod's resource requests do not fit anywhere right now" --
  a request/capacity problem, not an outage.
- Restarting the Deployment (`kubectl rollout restart`) as a first
  response before reading a single event -- this sometimes "fixes" a
  transient problem by luck, teaches the team nothing, and can mask a
  real bug that will recur.

### Check yourself
1. A Pod shows `STATUS: Pending` and `describe` says "Insufficient cpu."
   What are two different valid fixes, and what's the trade-off between
   them?
2. Why can `kubectl logs` show nothing at all for a container that is
   definitely crashing, and what command reveals the real error instead?
3. A Pod is `Running` and `1/1 Ready`, but a Service pointing at it times
   out. Name two places, other than the Pod's own health, that could be
   the actual cause.

### Further reading
- Kubernetes documentation, "Debug Running Pods" and "Debug Pods and
  ReplicationControllers" (kubernetes.io)
- Google Cloud, "Troubleshooting GKE clusters and workloads" documentation
- Sematext / Datadog, "Kubernetes troubleshooting: common issues and fixes"

## Chapter 60 -- Deployment patterns: rolling updates, health checks, rollbacks

### In one sentence
A Deployment's default rolling update replaces old Pods with new ones a
few at a time, gated entirely by readiness probes -- which means a
deployment is only as safe as the health check backing it, and a bad
health check turns "safe by design" into "outage by design."

### Why this matters
This is the single most consequential piece of Kubernetes configuration
most teams under-invest in. The rollout mechanics (`maxSurge`,
`maxUnavailable`, readiness gating) are genuinely well engineered and, when
paired with a correct readiness probe, make bad deploys self-limiting:
Kubernetes stops rolling forward the moment new Pods fail to become ready,
and traffic never reaches them. But if the readiness probe just checks
"is the process running" instead of "can this process actually serve a
request," the safety net has a hole in it exactly where it's needed most --
and the org discovers this precisely during the deploy that matters.

### How it actually works
```
   ROLLING UPDATE, Deployment with replicas: 4,
   maxSurge: 1, maxUnavailable: 1 (the defaults, roughly):

   Start:   [v1] [v1] [v1] [v1]                4 old, 0 new
   Step 1:  [v1] [v1] [v1] [v1] [v2]            surge: create 1 new
                                    ^ wait for readiness probe to pass
   Step 2:  [v1] [v1] [v1]      [v2] [v2]      1 old terminated, 1 more new
   Step 3:  [v1] [v1]          [v2] [v2] [v2]
   Step 4:  [v1]               [v2] [v2] [v2] [v2]
   Step 5:                     [v2] [v2] [v2] [v2]   done

   At every step, the Service only ever sends traffic to Pods that
   are Ready. A new Pod that fails its readiness probe simply never
   receives traffic and the rollout STALLS there -- it does not
   proceed to kill more old Pods. This is the core safety property.
```

Two different probes, often confused:

```
   livenessProbe:   "is this process healthy, or should it be
                      killed and restarted?" A failing liveness
                      probe causes a container restart. Wrong
                      liveness config -> death spiral: probe fails
                      under normal load -> container restarted ->
                      cold-starts and is slow -> fails probe again.

   readinessProbe:  "should this Pod currently receive traffic?"
                      A failing readiness probe removes the Pod from
                      the Service's endpoint list WITHOUT killing it.
                      This is the one that gates rollouts and is safe
                      to make strict, because failing it doesn't
                      destroy anything -- it just stops sending new
                      requests until the Pod says it's ready again.
```

Rollback is a first-class operation, not a manual redeploy:
```bash
kubectl rollout history deployment/api
kubectl rollout undo deployment/api              # back to previous revision
kubectl rollout undo deployment/api --to-revision=3
kubectl rollout status deployment/api            # watch a rollout live
```

### Real-world example
A widely cited pattern in postmortems (Monzo, Cloudflare, and many
smaller companies' public engineering blogs describe variants of this) is
the "healthy-looking crash": a service's `/healthz` endpoint returns 200
as long as the HTTP server thread is alive, but the service depends on a
database connection pool that can silently exhaust. The liveness probe
keeps passing because the HTTP server itself never dies -- it just returns
500s to every real request. Kubernetes has no idea anything is wrong: the
container isn't restarted (liveness passes) and the rollout isn't blocked
(readiness, if it exists at all, also just checks the HTTP server). The
fix pattern that shows up repeatedly across these postmortems is the same:
separate liveness (cheap, "is the process alive") from readiness
(meaningful, "can I actually reach my database and process a real
request"), and make the readiness probe do at least one lightweight
dependency check. Kubernetes's own documentation explicitly warns against
conflating the two probes for exactly this reason -- it is one of the most
common production misconfigurations in the ecosystem, common enough that
it appears in nearly every "Kubernetes gotchas" talk given at KubeCon.

### Try it yourself
```bash
kubectl create deployment demo --image=nginx --replicas=4
kubectl set image deployment/demo nginx=nginx:1.25
kubectl rollout status deployment/demo
kubectl rollout history deployment/demo

# simulate a bad deploy and watch Kubernetes refuse to finish rolling out
kubectl set image deployment/demo nginx=nginx:this-tag-does-not-exist
kubectl rollout status deployment/demo --timeout=20s
# it will report the rollout is stuck -- old Pods are still serving

kubectl rollout undo deployment/demo
kubectl rollout status deployment/demo
kubectl delete deployment demo
```

### Build it in Go
A Go service needs three things to make rolling updates seamless: a
`/readyz` that fails as soon as SIGTERM arrives, a delay before
`srv.Shutdown`, and `Shutdown` itself to drain in-flight requests. The
[TCP/IP guide's `lbdrain` lab](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries) proves the difference
under load, and [Chapter 82](#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes) here has the matching
Kubernetes manifest (probes, `preStop`, `terminationGracePeriodSeconds`).

### Common mistakes
- Using the same probe definition (or no readiness probe at all) for
  both liveness and readiness, so a dependency outage looks "healthy" to
  Kubernetes.
- Setting `initialDelaySeconds` too low for a slow-starting app (JVM
  services are the classic offender), causing the liveness probe to kill
  the container mid-startup, forever, before it ever gets a chance to
  finish booting.
- Setting `maxUnavailable: 0` and `maxSurge: 0` (accidentally, via a
  misunderstood template) which makes rollouts impossible to progress at
  all, and the Deployment just hangs.
- Assuming `kubectl rollout undo` also rolls back a database migration
  that shipped in the same release -- it only reverts the Deployment's Pod
  template, not any side effect the new code already caused.

### Check yourself
1. What is the actual, mechanical difference between a livenessProbe
   failure and a readinessProbe failure -- what does Kubernetes do in
   each case?
2. A service's `/healthz` always returns 200 even when its database is
   unreachable. What kind of incident does this specifically enable, and
   why doesn't Kubernetes catch it?
3. Why does a rolling update with `maxUnavailable: 1` on a Deployment with
   4 replicas stall, rather than fail outright, when new Pods can't
   become ready?

### Further reading
- Kubernetes documentation, "Configure Liveness, Readiness and Startup
  Probes"
- Kubernetes documentation, "Performing a Rolling Update"
- Honeycomb / Monzo engineering blogs, various postmortems on health
  check misconfiguration (search "readiness probe incident")

## Chapter 61 -- Helm and packaging Kubernetes applications

### In one sentence
Helm is a package manager for Kubernetes -- it turns a directory of
templated YAML into a versioned, install/upgrade/rollback-able "chart," for
the same reason `apt` exists instead of everyone compiling and copying
binaries by hand.

### Why this matters
A real application is rarely one Deployment. It's a Deployment, a
Service, a ConfigMap, maybe a Secret, an Ingress, a HorizontalPodAutoscaler
-- 6 to 15 YAML files that all need to agree with each other (same labels,
same names, same namespace) and that change slightly between environments
(dev uses 1 replica and no autoscaling, prod uses 5 and does). Copy-pasting
these YAML files between environments and hand-editing the differences is
exactly the class of manual, error-prone, "someone forgot to update the
image tag in staging" work that Helm exists to eliminate.

For a manager: Helm charts are also how most third-party software gets
installed onto Kubernetes today (Prometheus, cert-manager, ingress
controllers, databases) -- "just `helm install` it" is the standard onboarding
instruction across the ecosystem, so basic Helm literacy is close to a
baseline expectation for anyone operating a cluster, not a specialist skill.

### How it actually works
```
   A CHART is a directory:

     mychart/
       Chart.yaml          -- name, version, description
       values.yaml         -- default configuration values
       templates/
         deployment.yaml   -- Go-template YAML, referencing .Values.*
         service.yaml
         configmap.yaml
       charts/              -- bundled dependencies (sub-charts)

   templates/deployment.yaml might contain:

     replicas: {{ .Values.replicaCount }}
     image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"

   values.yaml (the defaults):
     replicaCount: 1
     image:
       repository: myapp
       tag: latest

   A production override, values-prod.yaml:
     replicaCount: 5
     image:
       tag: v2.4.1
```

The commands that matter day to day:
```bash
helm install myapp ./mychart -f values-prod.yaml   # first install
helm upgrade myapp ./mychart -f values-prod.yaml   # apply changes
helm rollback myapp 3                              # back to revision 3
helm history myapp                                 # see revisions
helm template ./mychart -f values-prod.yaml         # render YAML,
                                                     # don't apply --
                                                     # use this to review
                                                     # exactly what
                                                     # would be sent to
                                                     # the cluster
helm uninstall myapp
```

Every `helm upgrade` is a new numbered revision, tracked by Helm as a
Secret in the cluster -- which is what makes `helm rollback` possible: it
isn't magic, it's "re-apply the exact rendered YAML from revision N."

### Real-world example
The Kubernetes package ecosystem's default answer to "how do I run
Prometheus" or "how do I run an nginx ingress controller" is, almost
universally, a Helm chart from the project's own repository or from
Artifact Hub (the successor to the old central "stable" Helm chart
repository). This is not incidental: Helm charts became the de facto
distribution format precisely because hand-writing 10+ YAML files for
Prometheus's Deployment, Services, RBAC rules, ConfigMaps, and
PersistentVolumeClaims -- correctly, for your specific cluster's storage
class and resource limits -- is enough work that essentially nobody wants
to do it themselves. A commonly repeated cautionary lesson in platform
engineering, though, is `helm install`-ing a chart without reading
`values.yaml` first: charts often ship reasonable defaults for a generic
cluster, but a default like `persistence.enabled: false` on a database
chart or a `resources: {}` with no limits can silently mean "this loses
its data on every restart" or "this can consume unbounded memory on your
node" -- and because Helm makes installation so effortless, teams
routinely encounter these gaps only once the workload is already running
in production. The practical rule that has emerged across SRE teams is:
`helm template` and read the rendered output before ever running `helm
install` on anything you didn't author yourself.

### Try it yourself
```bash
# Requires helm installed and a cluster (minikube/kind is fine).
helm create demochart
ls demochart/templates/

# Render without installing -- read this output before ever installing
# a chart you didn't write
helm template demochart

helm install demo ./demochart
helm list
helm upgrade demo ./demochart --set replicaCount=3
helm history demo
helm rollback demo 1
helm uninstall demo
```

### Common mistakes
- Running `helm install` on an unfamiliar third-party chart without first
  running `helm template` to see exactly what it creates -- including
  RBAC permissions and resource limits (or their absence).
- Editing a running application's Kubernetes objects directly with
  `kubectl edit` after it was installed via Helm -- the next `helm
  upgrade` silently overwrites the manual change, because Helm has no
  idea it happened.
- Treating `values.yaml` defaults as production-appropriate without
  reading them, especially `replicaCount`, `resources`, and
  `persistence` settings.
- Not pinning chart versions (`helm install myapp repo/chart` without
  `--version`), so the same install command produces a different chart
  version, with different defaults, weeks later.

### Check yourself
1. What does `helm rollback myapp 3` actually do, mechanically -- what
   is Helm re-applying?
2. Why does hand-editing a Helm-managed Kubernetes object with `kubectl
   edit` cause problems later, even if the edit itself is correct?
3. Why is `helm template` recommended before `helm install` for a chart
   you didn't write yourself?

### Further reading
- Helm official documentation, "Helm Charts" and "Using Helm"
  (helm.sh/docs)
- Artifact Hub (artifacthub.io) -- the community index of published
  Helm charts
- CNCF, "Helm" project page and graduation announcement (Helm graduated
  as a CNCF project in 2020)

---

# Part 18 -- Capstone projects

Reading about OS internals, Linux commands, and Kubernetes gets you to
"I understand this." Building something end to end, and then deliberately
breaking it and fixing it, gets you to "I can do this under pressure,"
which is the actual bar for a beginner SRE or backend engineer role. These
four projects are designed to be done in order -- each one assumes the
skills from the previous one and the previous parts of this guide.

There is no single "correct" solution to any of these. What matters is
that you can explain every decision you made, in the vocabulary this guide
has given you: syscalls, page cache, cgroups limits, readiness probes,
image layers.

## Chapter 62 -- Project 1: build and profile a small program, OS-level

### In one sentence
Write a small program that does real work, then use the OS-level tools
from Parts 1-6 (`strace`, `time`, `/proc`, `ulimit`, `top`) to explain, in
kernel terms, exactly where its time and memory go.

### Why this matters
This project exists to make the abstract concrete: "page fault," "context
switch," and "syscall overhead" stop being vocabulary words and become
numbers you personally measured on your own machine. This is also close
to the actual daily work of performance engineering -- "why is this slow"
is answered by measurement, not by guessing, and the tools here are the
same ones used to debug production latency issues at any company running
Linux.

### How it actually works
The project, concretely:

```
   STEP 1 -- Write a program (any language) that:
     a) reads a large text file (a few hundred MB -- e.g. concatenate
        some log files, or download a public dataset)
     b) counts word frequency (a classic, CPU + I/O mixed workload)
     c) writes the top 20 words to stdout

   STEP 2 -- Measure it plainly first:
     /usr/bin/time -v ./wordcount bigfile.txt
     Read every line of the output: elapsed time, user time, system
     time, maximum resident set size, major/minor page faults,
     voluntary/involuntary context switches. Write down what each
     number means in your own words, using Chapters 4-13's vocabulary.

   STEP 3 -- Trace its syscalls:
     strace -c ./wordcount bigfile.txt
     Which syscall dominates the count? Which dominates total time?
     (Usually `read` calls dominate count; for a naive line-by-line
     reader, this can be enormous -- tens of thousands of small reads
     instead of a handful of large ones.)

   STEP 4 -- Change the I/O strategy and re-measure:
     Rewrite the file-reading part to read in large chunks (e.g. 1MB
     buffers) instead of line-by-line with an unbuffered reader, or
     use memory-mapping (mmap) instead of read(). Re-run STEP 2 and
     STEP 3. Compare: syscall count, elapsed time, major page faults.

   STEP 5 -- Constrain it with cgroups and watch it suffer honestly:
     systemd-run --scope -p MemoryMax=50M ./wordcount hugefile.txt
     If the file is bigger than 50MB and your program loads it all
     into memory at once, watch it get OOM-killed. Check `dmesg` or
     `journalctl -k` for the OOM killer's log line (Chapter 11).
     Then fix the program to stream the file instead of loading it
     all at once, and show it now completes under the same limit.
```

### Real-world example
This exact exercise -- profile, find the surprising bottleneck, fix it,
re-measure -- mirrors one of the most famous real performance
investigations in systems literature: the discovery that naive
line-by-line I/O (calling `read()` or `getline()` once per line) can be
10-100x slower than buffered, chunked reads, purely because of syscall
overhead, not because of disk speed. This is precisely the kind of
finding Brendan Gregg documents extensively in his USE Method and
"Systems Performance" methodology (used at Netflix, and taught widely in
the SRE community): most performance problems are not "the CPU is too
slow," they are "the program is making the kernel do far more work than
the task actually requires," and the only way to see that is to measure
syscalls and page faults directly instead of guessing from wall-clock time
alone.

### Try it yourself
This chapter's project *is* the "try it yourself" -- work through Steps 1
through 5 above on your own machine (a Linux VM or container is fine;
`systemd-run` cgroup limiting won't work in plain macOS/Windows, but
Docker's `--memory` flag gives you the equivalent, see Chapter 47).

### Common mistakes
- Measuring only wall-clock time (`time ./program`) and stopping there --
  wall-clock time alone doesn't tell you *why* it's slow, only *that* it
  is.
- Optimizing before measuring -- rewriting code based on a guess about
  what's slow, then finding the real bottleneck was somewhere else
  entirely.
- Running the "before" and "after" measurements under different
  conditions (different file, different machine load, cache warm vs
  cold) and drawing conclusions from a comparison that isn't
  apples-to-apples. Run each version multiple times and note whether the
  file is already in the page cache (Chapter 15) from a previous run.

### Check yourself
1. If `/usr/bin/time -v` shows a large number of "minor" page faults but
   zero "major" page faults, what does that tell you, and should you be
   worried?
2. Why would switching from line-by-line reads to 1MB chunked reads
   reduce the syscall count by orders of magnitude, and why does that
   matter even though both approaches read the same total bytes?
3. Your program gets OOM-killed under a 50MB memory limit even though the
   input file is only 30MB. What are two plausible reasons memory usage
   could exceed the file size?

### Further reading
- Brendan Gregg, "Systems Performance: Enterprise and the Cloud," 2nd
  edition (2020) -- especially the chapters on methodology and the USE
  Method
- `man 1 time`, `man 1 strace`, `man 1 systemd-run`

## Chapter 63 -- Project 2: containerize a real application properly

### In one sentence
Take an application with a real dependency (a database, or at least a
config file and a port) and containerize it following every practice from
Part 15 -- non-root user, multi-stage build, pinned base image, health
check, `.dockerignore`, and a docker-compose setup with a second service.

### Why this matters
Nearly every backend and SRE job posting today lists "Docker" as a
requirement, and nearly every take-home interview exercise or first-90-days
onboarding task is some version of "containerize this." The gap between
"I can write a Dockerfile that works" and "I can write a Dockerfile a
security-conscious team would approve in code review" is exactly the
content of Chapters 47-49 -- this project is where that content becomes a
reflex instead of a checklist you look up each time.

### How it actually works
```
   STEP 1 -- Pick or write a small app with a real dependency:
     A web API (any language) that reads/writes to a Postgres or Redis
     instance is ideal -- it forces you to deal with networking between
     containers, not just a single standalone process.

   STEP 2 -- Write a Dockerfile that gets every practice right:
     - Multi-stage build: a "builder" stage with full build tools,
       a final stage with only the runtime and the built artifact
     - Pinned, minimal base image (e.g. python:3.12-slim, not
       python:latest)
     - A dedicated, non-root user (USER appuser)
     - A .dockerignore excluding .git, node_modules/venv, secrets
     - A HEALTHCHECK instruction that hits a real readiness endpoint
     - Layers ordered so dependency installation is cached separately
       from application code copying (Chapter 48)

   STEP 3 -- Write docker-compose.yml with two services:
     the app, and its dependency (postgres or redis), using
     depends_on, a named volume for the database's data directory,
     and environment variables (not hardcoded secrets) for connection
     details.

   STEP 4 -- Verify the security properties, don't just assume them:
     docker inspect <container> --format '{{.Config.User}}'
       -- confirm it is NOT empty/root
     docker exec <container> whoami
       -- confirm the actual running user
     docker history <your-image>
       -- confirm no secrets baked into a layer
     docker images <your-image> --format '{{.Size}}'
       -- compare against a naive single-stage build's size

   STEP 5 -- Break the dependency on purpose and observe:
     docker compose stop postgres
     -- does your app crash hard, retry gracefully, or hang forever?
     Fix it to fail predictably (clear error, healthcheck reflects it)
     rather than hanging or crash-looping without useful logs.
```

### Real-world example
The size and security difference between a naive and a properly built
image is not a minor stylistic detail -- it is routinely a 5-10x size
difference and the difference between an image with zero known
non-root-exploitable paths and one where a container escape immediately
grants root on the host. A commonly cited real comparison: a naive
`FROM node` + `npm install` + copy-everything Dockerfile for a typical
Node.js service often produces an 900MB-1.2GB image running as root; the
equivalent multi-stage build on `node:20-alpine`, with a non-root user and
a `.dockerignore` excluding `node_modules` and `.git`, commonly lands
under 150MB. Snyk's annual container security reports have repeatedly
found that a large fraction of the most-pulled public images on Docker
Hub run as root by default and contain known-vulnerable packages baked
into unnecessary build tooling that a multi-stage build would have
discarded -- which is precisely the gap this project is designed to close
for you personally, on a project small enough to fully understand.

### Try it yourself
This chapter's project *is* the "try it yourself" -- work through Steps 1
through 5 above. If you don't have an existing app handy, a small Flask
or Express API with a single `/health` and a single `/items` endpoint
backed by Postgres is enough; the goal is the containerization discipline,
not the application's business logic.

### Common mistakes
- Declaring victory once `docker build` succeeds, without checking image
  size, running user, or what `docker history` reveals about layer
  contents.
- Using `latest` tags for the base image or the dependency image in
  compose, making the "same" setup produce different results days apart.
- Forgetting `depends_on` does not wait for Postgres to be *ready* --
  only for its container to have *started* (Chapter 49) -- and being
  surprised by intermittent connection failures on `docker compose up`.
- Putting real credentials in `docker-compose.yml` committed to git,
  instead of a `.env` file that is itself gitignored.

### Check yourself
1. Why does a multi-stage build typically produce a dramatically smaller
   final image than a single-stage build, even though both install the
   same build tools at some point?
2. `docker compose up` starts your app before Postgres is actually ready
   to accept connections, even with `depends_on` configured. Why, and
   what are two ways to fix it?
3. You confirm your Dockerfile has a `USER appuser` line, but `docker exec
   <container> whoami` still prints `root`. What's a plausible
   explanation?

### Further reading
- Snyk, "State of Open Source Security" / container image reports
  (annual, search current year's edition)
- Docker official documentation, "Best practices for writing Dockerfiles"
- [[Chapter 48]] and [[Chapter 49]] of this guide (multi-stage builds,
  volumes and networking)

## Chapter 64 -- Project 3: deploy that application to Kubernetes with health checks

### In one sentence
Take the containerized application from Project 2 and deploy it to a real
(even if local, via minikube or kind) Kubernetes cluster with correct
liveness/readiness probes, resource requests/limits, a Service, and a
ConfigMap/Secret for configuration -- then perform a rolling update and a
rollback on it.

### Why this matters
This is the project that proves you can operate, not just read about,
everything in Parts 14-17. It is also close to a literal recreation of the
most common first real task given to a new hire on a platform or SRE team:
"here's a containerized app someone else wrote, put it on the cluster
properly." Doing this once, deliberately and carefully, with your own
hands, converts a dozen chapters of reading into something you can discuss
concretely in an interview or a design review.

### How it actually works
```
   STEP 1 -- Set up a local cluster:
     minikube start        (or: kind create cluster)
     kubectl get nodes     -- confirm it's up

   STEP 2 -- Write the Kubernetes manifests (or a Helm chart, per
   Chapter 61, if you want the extra practice):
     - a Deployment for your app image, with:
         resources.requests and resources.limits set deliberately,
         not left blank
         a livenessProbe hitting a cheap endpoint
         a readinessProbe hitting an endpoint that actually checks
         the database connection
     - a ConfigMap for non-secret configuration
     - a Secret for the database password / connection string
     - a Service (ClusterIP is fine for this exercise) in front of
       the Deployment
     - if using a database, a separate Deployment (or StatefulSet) +
       Service + PersistentVolumeClaim for it

   STEP 3 -- Deploy and verify with the Chapter 59 debugging sequence,
   even if nothing is actually broken -- practice the habit:
     kubectl apply -f .
     kubectl get pods -w
     kubectl describe pod <name>
     kubectl logs <name>

   STEP 4 -- Perform a real rolling update:
     Change something visible (a response header, a log line) in your
     app, rebuild the image with a new tag, kubectl set image ..., and
     watch `kubectl rollout status` complete cleanly.

   STEP 5 -- Perform a rollback, on purpose, under a manufactured
   failure:
     Deploy a broken version (bad image tag, or a readiness probe
     pointed at a path that 404s) and confirm the rollout correctly
     STALLS rather than replacing every old Pod. Then:
     kubectl rollout undo deployment/<name>
     and confirm service was never actually interrupted for end users
     during either the bad deploy attempt or the rollback.
```

### Real-world example
This project sequence -- deploy, then deliberately break a rollout, then
roll back -- is essentially a compressed version of what companies call
"chaos engineering" or "game days" (a practice popularized publicly by
Netflix's Chaos Monkey and now standard at many companies operating
Kubernetes at scale). The entire value of Kubernetes's rolling-update
safety mechanism (Chapter 60) is unproven to you personally until you have
watched it actually refuse to finish rolling out a broken version, with
your own eyes, on your own cluster. Teams that skip this kind of
deliberate practice often discover the *first* time a rollout stalls in
production is also the first time anyone on the team has seen what a
stalled rollout looks like -- at 2 a.m., under incident pressure, which is
a strictly worse time to learn it.

### Try it yourself
This chapter's project *is* the "try it yourself" -- work through Steps 1
through 5 above, on minikube or kind. Budget a real afternoon for this one;
rushing it defeats the purpose.

### Common mistakes
- Leaving `resources.requests` unset "to get it working first," then
  never coming back to set it -- which is how clusters end up with
  unpredictable scheduling and noisy-neighbor problems in real
  production environments (Chapter 58).
- Writing a readiness probe that checks the same thing as the liveness
  probe (Chapter 60's central lesson), defeating the purpose of having
  both.
- Not actually watching the rollout fail before rolling back -- clicking
  straight to `rollout undo` skips the entire learning point of the
  exercise, which is seeing Kubernetes's safety mechanism work.
- Storing the database password in the ConfigMap instead of a Secret --
  functionally similar in this toy exercise, but the habit is worth
  building correctly from the start (Chapter 57).

### Check yourself
1. Why is it important to deliberately deploy a *broken* version in this
   exercise rather than only ever deploying working versions?
2. What's the difference in outcome between a readiness probe failure and
   a liveness probe failure during a rolling update, and why does that
   difference matter for user-facing availability?
3. If `resources.requests` is left unset entirely, what does the
   Kubernetes scheduler assume, and what problem can that cause on a
   busy cluster?

### Further reading
- Netflix Technology Blog, "The Netflix Simian Army" (introducing Chaos
  Monkey, 2011)
- Kubernetes documentation, "Kubernetes Basics" interactive tutorial
- [[Chapter 58]], [[Chapter 59]], [[Chapter 60]] of this guide

## Chapter 65 -- Project 4: diagnose and fix a planted incident, end to end

### In one sentence
Have a friend (or your future self, a week later, having forgotten the
details) deliberately break your Project 3 deployment in one specific way,
without telling you what they changed, and practice the full incident
process: detect, diagnose, mitigate, fix, and write a short postmortem.

### Why this matters
Every skill in this guide -- OS internals, Linux commands, container and
Kubernetes operation -- exists to be used under exactly this condition:
something is broken, you don't know what, and people are waiting on you.
This project is the only one in this guide that simulates that condition
honestly, because in Projects 1-3 you always already knew what you had
just changed. Practicing diagnosis on a problem *someone else* planted,
with no hints, is the closest a self-study guide can get to a real
on-call page -- and it is exactly what SRE and platform interviews
increasingly test for directly (major companies run live "debug this
broken cluster" exercises as part of their SRE interview loops).

### How it actually works
```
   SETUP (ask a friend, a study partner, or do this to your own
   deployment a week later after you've forgotten the details):
   Introduce exactly ONE of these planted failures, chosen at random,
   into the Project 3 Kubernetes deployment:

     a) Change the readiness probe's path to one that 404s
     b) Set a memory limit lower than the app actually needs, causing
        OOMKilled restarts under load
     c) Change a ConfigMap value the app needs (e.g. a database host
        name) to something wrong
     d) Change the Service's label selector so it no longer matches
        the Deployment's Pods
     e) Set an overly aggressive livenessProbe (very short
        initialDelaySeconds) causing restart-loop on a normally-slow-
        starting app
     f) Fill the node's disk (or a PVC) to force a storage-related
        failure

   YOUR PROCESS, once told only "the app is down, users are getting
   errors":

     1. DETECT   -- confirm the actual user-visible symptom first
                    (curl the Service / Ingress yourself)
     2. DIAGNOSE -- run the Chapter 59 sequence: get pods, describe,
                    logs, events. Form a hypothesis before changing
                    anything.
     3. MITIGATE -- stop the bleeding first. Often this means
                    `kubectl rollout undo` rather than root-causing
                    live, if a previous known-good revision exists.
     4. FIX      -- once mitigated, identify and fix the actual root
                    cause (not just the symptom).
     5. WRITE    -- a half-page postmortem: what broke, how you
                    detected it, what you tried that didn't help,
                    what actually fixed it, and one concrete change
                    (a monitor, an alert, a stricter code review
                    check) that would have caught this earlier.
```

### Real-world example
The five-step process above -- detect, diagnose, mitigate, fix, write --
is a direct, simplified mirror of the incident response process used at
Google, Amazon, and most companies with a mature SRE practice, as
documented in Google's *Site Reliability Engineering* book (the
"mitigate first, root-cause second" ordering in particular is one of its
most emphasized and most frequently learned-the-hard-way lessons: engineers
under pressure to find the "real" cause often delay an available, working
mitigation like a rollback, extending an outage's user impact for the
sake of intellectual completeness that could have waited). Blameless
postmortems -- the practice of writing up what happened without assigning
personal fault, focusing instead on what the system and process allowed to
happen -- is likewise a direct practice from that same book and from
Etsy's widely cited "Blameless PostMortems and a Just Culture" (2012)
engineering blog post, one of the most influential single blog posts in
the SRE field's history. Practicing the format here, on a low-stakes toy
incident, is what makes writing one under real pressure feel familiar
instead of foreign.

### Try it yourself
This chapter's project *is* the "try it yourself." If no study partner is
available, write all six planted failures (a-f above) on slips of paper,
draw one at random a week from now without looking, apply it to your
Project 3 deployment, and force yourself to diagnose it cold.

### Common mistakes
- Skipping straight to "fix the root cause" under time pressure instead
  of mitigating first -- in a real incident, this is the single most
  common mistake that extends outage duration.
- Guessing and changing multiple things simultaneously ("let me just
  restart everything and bump the memory limit and check the
  ConfigMap") -- this can accidentally fix the problem while leaving you
  unable to explain why, which defeats the exercise and, in a real
  incident, often introduces a second unrelated problem.
- Writing a postmortem that focuses on blame ("I should have checked
  this sooner") rather than systemic prevention ("we should have an
  alert that fires when readiness-probe failures exceed N for M
  minutes").
- Not timing yourself. Real incidents have a clock running; practicing
  without one removes the exact pressure this exercise exists to
  simulate.

### Check yourself
1. Why does "mitigate first, root-cause second" often mean choosing
   `kubectl rollout undo` over continuing to investigate, even if you
   haven't yet found the actual root cause?
2. What makes a postmortem "blameless" in practice -- what does it focus
   on instead of individual fault?
3. Of the six planted-failure categories listed, which ones would the
   Chapter 59 `describe pod` Events block surface almost immediately, and
   which ones would require checking application logs or external
   symptoms first?

### Further reading
- Google, *Site Reliability Engineering* (2016), Chapter 12,
  "Effective Troubleshooting," and Chapter 15, "Postmortem Culture:
  Learning from Failure" -- free online at sre.google/sre-book
- Etsy Engineering, "Blameless PostMortems and a Just Culture" (2012)
- PagerDuty, "Incident Response Documentation" (public incident response
  guide)

---

# Part 19 -- Advanced Linux operations and expert practice

> Expertise is not knowing every flag. It is being able to take a vague
> production symptom, form a hypothesis, choose the least-invasive tool, prove
> what is happening, make a reversible change, and leave the system easier to
> operate than you found it.

## Chapter 66 -- Production Linux mental models: fleet, drift, and failure domains

### In one sentence
Production Linux is not one machine you lovingly tune; it is a fleet of
mostly-similar machines where consistency, safe rollout, observability, and
rollback matter as much as any individual command.

### Why this matters
Most beginner Linux learning happens on one laptop or VM. Real incidents happen
across dozens, hundreds, or thousands of nodes: one kernel version, one sysctl,
one package mirror, one container runtime setting, or one clock-sync problem can
split the fleet into "works here, fails there." Senior operators therefore ask
two questions constantly:

```
   Is this host representative?
   If I change this, how do I know what else I changed?
```

### How it actually works
Think in layers:

```
   IMAGE / BASELINE     distro version, kernel, packages, hardening profile
   BOOT CONFIG          kernel cmdline, initramfs, systemd units, cloud-init
   RUNTIME STATE        running processes, open files, mounts, sockets, cgroups
   DECLARED CONFIG      Ansible/Chef/Puppet/Nix/Terraform/Kubernetes manifests
   OBSERVED BEHAVIOUR   metrics, logs, traces, alerts, user symptoms
```

When the layers disagree, you have drift. Drift is not automatically bad — an
emergency one-line fix is sometimes the right move — but untracked drift is how
the same incident returns weeks later.

### Real-world scenario
A payment API is slow only on half the nodes. Application logs are identical.
The expert path is not "restart everything." It is:

```bash
# Compare the boring facts first.
uname -a
cat /etc/os-release
systemctl --version
mount | sort
sysctl net.ipv4.tcp_congestion_control vm.swappiness
```

Then compare package versions, kernel cmdline, cgroup mode, container runtime
version, and node labels. In many real incidents, the cause is not in the
application at all: a partial AMI rollout, one pool using cgroups v2, one pool
with a different filesystem mount option, or one pool missing a kernel module.

### Try it yourself
On two Linux machines or two disposable VMs, intentionally make one small
configuration difference:

```bash
sysctl net.ipv4.tcp_congestion_control
systemctl list-units --failed
findmnt -no SOURCE,FSTYPE,OPTIONS /
```

Write a short "fleet diff" note: what is the same, what differs, and which
differences could affect application behaviour?

### Common mistakes
- Debugging only the broken host and never checking a known-good host.
- Fixing runtime state manually but never committing the change to the source
  of truth.
- Treating "works after restart" as a root cause. Restarting is a mitigation;
  it is rarely the explanation.

### Check yourself
1. What is configuration drift?
2. Why is comparing a bad host to a good host often faster than reading every
   log on the bad host?
3. What is the difference between a mitigation and a root cause?

## Chapter 67 -- systemd beyond start/stop: dependencies, timers, sandboxing, and debugging

### In one sentence
systemd is not only a service launcher; it is a dependency engine, logger,
supervisor, timer system, resource-control interface, and practical hardening
tool.

### Why this matters
Many production failures are "the process is fine, but the unit is wrong":
wrong working directory, missing environment, bad restart policy, dependency
cycle, timeout too short, service starts before the network/storage is truly
ready, or a hardening option blocks a path the app needs.

### How it actually works
Useful commands for real diagnosis:

```bash
systemctl status myapp.service
journalctl -u myapp.service -b --no-pager
systemctl cat myapp.service
systemd-analyze verify ./myapp.service
systemd-analyze critical-chain myapp.service
systemctl show myapp.service -p Restart -p User -p Group -p ExecStart
```

Useful unit-file ideas:

```ini
[Unit]
Description=Example API
After=network-online.target
Wants=network-online.target

[Service]
User=myapp
Group=myapp
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/bin/server
Restart=on-failure
RestartSec=5s
TimeoutStartSec=30s
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/myapp /var/log/myapp

[Install]
WantedBy=multi-user.target
```

Do not cargo-cult hardening directives. `ProtectSystem=strict` is excellent
when the write paths are explicit; it is an outage when the service legitimately
needs to write somewhere you forgot to list.

### Real-world scenario
A service starts manually but fails under systemd. The usual cause is not
"systemd is broken"; it is that the interactive shell had environment variables,
current directory, PATH, open files, or permissions the unit did not have.

### Try it yourself
Create a disposable service that prints its environment:

```bash
systemd-run --user --unit env-demo /usr/bin/env
journalctl --user -u env-demo --no-pager
systemctl --user status env-demo
```

Compare that environment with your interactive shell:

```bash
env | sort | head
```

### Common mistakes
- Using `Restart=always` for a program that exits successfully by design.
- Editing `/lib/systemd/system/...` instead of using an override under
  `/etc/systemd/system/...`.
- Forgetting `systemctl daemon-reload` after changing unit files.
- Adding hardening flags without testing realistic application behaviour.

### Check yourself
1. Why can a program work in your shell but fail as a systemd service?
2. What does `systemctl cat` show that `systemctl status` does not?
3. Why should hardening settings be tested like code?

## Chapter 68 -- Storage operations: LVM, RAID, filesystem growth, backups, and recovery

### In one sentence
Storage expertise is knowing the stack from block device to filesystem to
application durability, and testing restore paths before an outage proves they
do not work.

### Why this matters
Disk incidents are unforgiving. A CPU spike usually gives you time to think; a
full filesystem, failed disk, corrupt journal, or accidental deletion can turn
into data loss quickly. Experts know which operations are online, which need a
maintenance window, and which require a verified backup first.

### How it actually works
Common stack:

```
   physical disk / cloud volume
      -> partition table
      -> LUKS encryption, maybe
      -> RAID or cloud-level replication, maybe
      -> LVM physical volume / volume group / logical volume, maybe
      -> filesystem: ext4, XFS, Btrfs, ZFS, etc.
      -> mount options
      -> application write pattern and fsync behaviour
```

Safe read-only inventory commands:

```bash
lsblk -f
findmnt
df -hT
du -xh /var 2>/dev/null | sort -h | tail
sudo journalctl -p warning -b --no-pager | grep -Ei 'disk|ext4|xfs|nvme|i/o'
```

Growth workflow, conceptually:

```
   enlarge underlying disk/volume
   grow partition or PV
   grow LVM logical volume, if used
   grow filesystem
   verify with df/findmnt/application metrics
```

The exact commands differ by filesystem and cloud provider. Do not paste a
resize recipe into production unless you have confirmed the device names, stack
layout, filesystem type, and backup state.

### Real-world scenario
`/var` is full because container logs grew without rotation. The wrong fix is
blind deletion. The safer flow is:

```bash
df -hT
sudo du -xh /var 2>/dev/null | sort -h | tail -20
journalctl --disk-usage
docker system df 2>/dev/null || true
```

Then decide: rotate/truncate logs, fix log retention, move data, expand the
volume, or clean unused container artifacts. The expert result is not "space
free right now"; it is "space free, alert threshold corrected, retention fixed,
and restore still tested."

### Try it yourself
In a disposable VM, attach or create a small extra disk/loop device and practise:
identify it with `lsblk`, create a filesystem, mount it under `/mnt/demo`,
write data, unmount it, and remount it by UUID. Do this only in a lab.

### Common mistakes
- Confusing RAID, backup, snapshot, and replication. They solve different
  problems.
- Assuming a snapshot is restorable because it exists.
- Running cleanup commands before identifying what is consuming space.
- Growing storage without fixing the growth rate that filled it.

### Check yourself
1. Why is a tested restore more important than a configured backup job?
2. What information do you need before resizing a filesystem?
3. Why can deleting an open log file fail to free space immediately?

## Chapter 69 -- Advanced Linux networking: namespaces, routing, nftables, and packet paths

### In one sentence
Linux networking becomes understandable when you trace the packet path through
interfaces, routing tables, namespaces, conntrack, firewall rules, and the
application socket.

### Why this matters
"Network issue" is often a label for very different failures: DNS, routing,
firewall, MTU, NAT, conntrack exhaustion, TLS, application backlog, or a service
listening on the wrong address. Experts separate those layers quickly.

### How it actually works
Core inspection sequence:

```bash
ip addr
ip route
ip rule
ss -ltnp
ss -s
resolvectl status 2>/dev/null || cat /etc/resolv.conf
sudo nft list ruleset 2>/dev/null || sudo iptables -S
```

Packet-path mental model:

```
   process socket
      -> network namespace
      -> local routing decision
      -> firewall/conntrack/NAT hooks
      -> interface queue
      -> driver/NIC or virtual bridge/veth
      -> remote network
```

Container networking is the same model with extra namespaces and virtual
interfaces. A Kubernetes Pod has its own network namespace; Services add load
balancing and NAT/eBPF rules depending on the cluster implementation.

### Real-world scenario
An app listens on `127.0.0.1:8080` inside a VM. Health checks from the same host
work, but remote traffic fails. `ss -ltnp` reveals the binding is loopback-only.
No firewall change will make another machine reach a socket that is not bound to
a routable address.

### Try it yourself
Create a network namespace and prove loopback isolation:

```bash
sudo ip netns add demo
sudo ip netns exec demo ip link set lo up
sudo ip netns exec demo ip addr
sudo ip netns del demo
```

Only do namespace and firewall experiments in a disposable lab; a bad firewall
rule can lock you out of a remote host.

### Common mistakes
- Starting with `tcpdump` before checking `ss`, `ip route`, and DNS.
- Forgetting that `localhost` inside a container is the container, not the host.
- Treating nftables/iptables rules as the only firewall; cloud security groups,
  Kubernetes NetworkPolicy, and host firewalls can all participate.
- Ignoring MTU when VPNs, overlays, or Kubernetes CNIs are involved.

### Check yourself
1. Why can a service be reachable locally but unreachable remotely?
2. What does a network namespace isolate?
3. Why is packet capture powerful but not the first step for every issue?

## Chapter 70 -- Observability and performance: perf, flame graphs, eBPF, and OpenTelemetry

### In one sentence
Expert performance work starts with a symptom, measures the bottleneck, and
only then reaches for deeper tools like `perf`, flame graphs, or eBPF.

### Why this matters
Most bad performance fixes are plausible guesses. Linux gives you enough
visibility to avoid guessing, but the tools have different costs and levels of
intrusion. `top` and `vmstat` are cheap; `strace` can slow a busy process;
profilers and eBPF tools require care and usually elevated privileges.

### How it actually works
Use a ladder:

```
   1. User symptom: slow endpoint, timeout, error rate, cost spike
   2. Golden signals: latency, traffic, errors, saturation
   3. Host view: CPU, memory, disk, network, run queue
   4. Process view: threads, file descriptors, syscalls, allocations
   5. Code/path view: profiles, flame graphs, traces
   6. Kernel view: scheduler, block I/O, TCP, locks, eBPF probes
```

Examples:

```bash
uptime
vmstat 1 5
pidstat -durh 1 5 2>/dev/null || true
perf top 2>/dev/null
```

`perf` answers "where is CPU time going?" eBPF-based tools can answer questions
like "which process is opening this file?", "which TCP connects are slow?", or
"which kernel function is hot?" without permanently instrumenting the
application. OpenTelemetry answers a different question: how requests move
through distributed services.

### Real-world scenario
A service is slow and CPU is high. Logs say nothing. A CPU profile shows most
time in JSON serialization, not the database. The fix is a code/data-shape
change, not more replicas. Another day, the same symptom is caused by disk I/O
wait from synchronous logging. Same symptom, different measurement, different
fix.

### Try it yourself
Profile a harmless CPU-heavy command:

```bash
openssl speed sha256 2>/dev/null &
DEMO_PID=$!
ps -p "$DEMO_PID" -o pid,comm,%cpu
wait "$DEMO_PID"
```

If your lab has `perf` installed and permissions allow it, run:

```bash
perf stat openssl speed sha256
```

### Build it in Go
Go has profiling built in: import `net/http/pprof` and every running service can
produce CPU, heap, goroutine, and block profiles on demand.
[Chapter 81](#chapter-81-profiling-go-on-linux-pprof-the-execution-tracer-and-perf) plants two bugs (a regexp compiled on every
request, and a goroutine leak) and finds both: the regexp turned out to be
86% of the handler's CPU, and the leak showed up as 489 goroutines stuck on one
line.

### Common mistakes
- Optimizing the function you dislike instead of the function the profile
  identifies.
- Looking only at average latency and missing tail latency.
- Treating observability as "logs only." Metrics, traces, profiles, and events
  answer different questions.
- Running heavy tracing on production without understanding overhead and data
  sensitivity.

### Check yourself
1. Why should you define the symptom before picking a tool?
2. What does a flame graph help you see?
3. How are distributed traces different from host metrics?

## Chapter 71 -- Linux security hardening: SELinux/AppArmor, auditd, capabilities, and least privilege

### In one sentence
Linux hardening is layered least privilege: reduce what a process can read,
write, execute, call, bind, load, and persist, while keeping enough observability
to detect when those boundaries are tested.

### Why this matters
Unix permissions are necessary but not sufficient. Real services also need
service accounts, limited sudo, capabilities instead of root, syscall filters,
MAC policy, secret handling, patching, logging, and a way to prove the controls
are working.

### How it actually works
Hardening layers:

```
   identity        dedicated user/group, no shared admin account
   filesystem      ownership, mode bits, mount options, read-only paths
   privilege       sudoers scope, Linux capabilities, no_new_privs
   isolation       systemd sandboxing, namespaces, containers, chroot where useful
   MAC             SELinux or AppArmor policy
   kernel          seccomp, module loading rules, exploit mitigations
   audit           auditd/journald/SIEM events tied to real questions
```

Useful inspection:

```bash
id myapp 2>/dev/null || true
getcap -r /usr/bin /usr/sbin 2>/dev/null
sudo -l
aa-status 2>/dev/null || true
getenforce 2>/dev/null || true
systemctl show myapp.service -p NoNewPrivileges -p CapabilityBoundingSet 2>/dev/null
```

### Real-world scenario
A web service runs as root only because it binds port 80. Better options:

```
   - put a reverse proxy/load balancer in front;
   - bind a high port and let the proxy listen on 80/443;
   - grant only CAP_NET_BIND_SERVICE to the binary, if appropriate;
   - use systemd socket activation.
```

The expert question is not "can I make it work?" It is "what is the smallest
privilege that makes it work, and how will I notice if it tries to exceed that?"

### Try it yourself
In a lab, compare a service run as your user, as root, and as a dedicated
service account. List the files it can write and the ports it can bind. Then
remove one privilege and observe the failure mode in logs.

### Common mistakes
- Disabling SELinux/AppArmor because of one denial instead of understanding the
  denial and adjusting policy or application behaviour.
- Giving `NOPASSWD: ALL` in sudoers for convenience.
- Running containers as root and assuming "container" means "safe."
- Collecting audit logs without a detection question.

### Check yourself
1. Why are Linux capabilities safer than full root for narrow privileges?
2. What problem do SELinux/AppArmor solve that mode bits do not?
3. Why is logging a hardening control only when someone reviews or alerts on it?

### Across the series

- **The attacker's view of these controls:** [Security Engineering in Depth](../security/real-life-security-guide-v1.md)
  covers container internals and escapes ([Chapter 15](../security/real-life-security-guide-v1.md#chapter-15-container-internals-and-isolation)),
  Kubernetes RBAC, admission, network, and runtime security
  ([Chapters 17–20](../security/real-life-security-guide-v1.md#chapter-17-kubernetes-rbac-and-identity)), with a Go lab for workload identity
  over mTLS ([Chapter 23](../security/real-life-security-guide-v1.md#chapter-23-service-identity-spiffe-and-spire)).
- **Hardening applied to one service:** the systemd unit and Kubernetes manifest in
  [Chapter 82](#chapter-82-capstone-shipping-a-well-behaved-go-service-image-systemd-kubernetes) (exposure score 3.1 vs 9.6 unhardened).


## Chapter 72 -- Incident response and refactoring on real Linux systems

### In one sentence
The senior Linux workflow is: stabilize, preserve evidence, understand the
failure, fix the smallest safe thing, verify, then refactor the system so the
same class of failure is less likely or easier to diagnose next time.

### Why this matters
Real incidents are not exam questions. You may need to restore service before
you fully understand the root cause, but every emergency action should be
visible, reversible when possible, and followed by cleanup.

### How it actually works
Production-safe incident rhythm:

```
   1. State the user-visible impact.
   2. Freeze context: what changed, when, on which hosts?
   3. Capture volatile evidence before restarting: logs, process state,
      sockets, open files, resource usage.
   4. Mitigate using the smallest reversible action.
   5. Verify recovery from the user's perspective.
   6. Root-cause with evidence.
   7. Refactor: tests, alerts, runbooks, config, capacity, rollback path.
```

Evidence commands:

```bash
date -Is
uptime
free -h
df -hT
ss -s
systemctl --failed
journalctl -p warning -b --no-pager | tail -100
```

Refactoring examples:

```
   repeated disk-full incident
      -> add log rotation, quota/retention, filesystem alert, runbook

   repeated OOM kill
      -> measure working set, fix leak or limit, add memory dashboard,
         tune request/limit, add load test

   repeated manual service restart
      -> fix readiness check, dependency ordering, backoff, crash alert,
         and deploy rollback
```

### Real-world scenario
An API recovers every time someone restarts it, so the team adds a cron restart.
That hides the leak and destroys evidence. A better response is to capture RSS,
heap/profile data if available, request rate, open file count, and OOM/restart
events; then choose a temporary mitigation while preserving enough information
to fix the leak.

### Try it yourself
Write a one-page runbook for "service is slow" using Appendix C as a starting
point. For each step, write:

```
   command to run
   what healthy looks like
   what bad looks like
   what action is safe
   what action needs approval
```

### Common mistakes
- Restarting before capturing the evidence that would explain the incident.
- Treating a one-off manual fix as done without adding a test, alert, or
  configuration change.
- Making several changes at once and losing the ability to know what helped.
- Writing a postmortem that blames a person instead of improving the system.

### Check yourself
1. Why can a restart make root-cause analysis harder?
2. What is the difference between mitigation, remediation, and prevention?
3. What should be added to a runbook after an incident?

---

# Part 20 -- Systems programming in Go: the OS from inside a program

> Every chapter so far looked at the kernel from the outside, through `ps`,
> `strace`, `free`, and `docker stats`. This part looks from the inside. You
> write small Go programs that make the kernel's behaviour visible and
> measurable, and you see where a modern language runtime and the operating
> system meet: threads, signals, cgroup limits, durability, file
> descriptors, `/proc`, namespaces, and profiling. Basic Go is enough. The
> [Go guide](../Golang/real-life-golang-guide.md) covers the language. Every program was run on Linux
> while writing this guide, and the outputs shown are real.

## Chapter 73 -- Go meets the kernel: system calls, threads, and the runtime

### In one sentence
A Go program is a normal Linux process whose runtime multiplexes many
goroutines onto a few OS threads, and how a goroutine *waits* decides whether
it costs an OS thread.

### Why this matters
"Goroutines are cheap" is true, but not unconditionally. A goroutine that
sleeps or waits on a socket costs a few kilobytes. A goroutine stuck in a
blocking system call (a slow disk read, a cgo call, an NFS stall) pins a whole
OS thread, and the runtime creates more threads to keep running other
goroutines. Services have hit Go's 10,000-thread limit and crashed this way.
Knowing which waits are cheap explains thread counts in `top`, memory
overhead, and some very strange incidents.

### How it actually works
The scheduler has three kinds of object (G, M, P, covered in the
[Go guide's scheduler chapter](../Golang/real-life-golang-guide.md#28-the-scheduler-deep-dive-gmp-and-work-stealing)):

```
   G  goroutine     your code; thousands to millions
   M  machine       an OS thread (what `ps -L` and /proc/<pid>/status "Threads:" count)
   P  processor     a slot allowed to run Go code; there are GOMAXPROCS of them
```

How a goroutine waits determines what happens to its thread:

| The goroutine is... | What the runtime does | Thread cost |
|---|---|---|
| sleeping, waiting on a channel or mutex | parks the G; the M runs other Gs | none |
| reading/writing a **socket** or pipe | parks the G and registers the fd with **epoll** (the netpoller); one thread waits for all of them | none |
| in a **blocking system call** (file I/O, `syscall.Read`) | the M blocks inside the kernel; the P is handed to another M, which may be new | one thread per blocked call |
| in **cgo** | same as a blocking system call | one thread per call |

`strace -f -c ./yourprogram` on any Go binary shows the runtime's own system
calls: `futex` (threads sleeping and waking), `epoll_pwait` (the netpoller),
`clone` (new threads), `mmap` (heap growth), `rt_sigprocmask` (signal handling).

### Build it in Go (20 min) -- count the threads

```go
//go:build linux

// threads: how many OS threads does a Go program really use? It depends on
// HOW goroutines wait. Run it and compare the three cases.
//
//	go run ./threads
package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func threads() string {
	b, _ := os.ReadFile("/proc/self/status")
	for _, l := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(l, "Threads:") {
			return strings.TrimSpace(strings.TrimPrefix(l, "Threads:"))
		}
	}
	return "?"
}

func report(label string) {
	time.Sleep(500 * time.Millisecond) // let the runtime settle
	fmt.Printf("%-52s goroutines=%-6d OS threads=%s\n", label, runtime.NumGoroutine(), threads())
}

func main() {
	report("start")

	// 1. 10,000 goroutines in time.Sleep: parked by the Go scheduler, no thread each.
	for i := 0; i < 10000; i++ {
		go time.Sleep(time.Hour)
	}
	report("+10,000 goroutines sleeping")

	// 2. 10,000 goroutines blocked reading network sockets: parked in the
	//    netpoller (one epoll instance), still no thread each.
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func() { buf := make([]byte, 1); c.Read(buf) }() // waits forever
		}
	}()
	var conns []net.Conn
	for i := 0; i < 5000; i++ { // 5,000 client + 5,000 server goroutines
		c, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			fmt.Println("dial:", err)
			break
		}
		conns = append(conns, c)
		go func() { buf := make([]byte, 1); c.Read(buf) }()
	}
	report("+10,000 goroutines blocked on socket reads")

	// 3. 100 goroutines blocked in a RAW blocking system call: each one pins
	//    an OS thread, and the runtime starts new threads to keep running Go code.
	for i := 0; i < 100; i++ {
		var p [2]int
		syscall.Pipe(p[:])
		go func() { buf := make([]byte, 1); syscall.Read(p[0], buf) }() // blocks the thread itself
	}
	report("+100 goroutines blocked in raw read(2) syscalls")
	runtime.KeepAlive(conns)
}
```

Real output (Linux, 14 CPUs):

```text
start                                                goroutines=1      OS threads=5
+10,000 goroutines sleeping                          goroutines=10001  OS threads=16
+10,000 goroutines blocked on socket reads           goroutines=20002  OS threads=16
+100 goroutines blocked in raw read(2) syscalls      goroutines=20102  OS threads=108
```

**What to notice:**

- **20,000 goroutines, 16 threads.** Sleeping and network-blocked goroutines
  cost no threads. Ten thousand idle TCP connections are cheap in Go, which is
  why it's popular for proxies and gateways.
- **100 blocking system calls, +92 threads.** Each raw `read(2)` on a pipe
  that has no data pins its thread inside the kernel. (Go's `os.File` puts
  pipes in the netpoller too. We used `syscall.Read` to bypass it.)
- **Regular files can't use epoll on Linux**, so `os.ReadFile` on a slow disk or
  network filesystem blocks a thread for real. A thousand goroutines reading
  from a stalled NFS mount means about a thousand threads. The runtime crashes
  the program at 10,000 (`runtime/debug.SetMaxThreads`).

**Exercises:**

1. Run any lab under `strace -f -c -o trace.txt` and read the summary. Which
   system call is called most?
2. Run `GODEBUG=schedtrace=1000 ./threads` and read the scheduler's
   once-a-second report: `gomaxprocs`, `threads`, `idleprocs`, `runqueue`.
3. Lower the limit with `debug.SetMaxThreads(50)` and rerun. What happens at
   the syscall step?

### Real-world example
A log-shipping service written in Go read thousands of files from an NFS
share, one goroutine per file. When the NFS server stalled, every read blocked
its thread, the thread count climbed past 10,000, and the process died with
`runtime: program exceeds 10000-thread limit`. The fix was a semaphore (a
buffered channel of capacity 64) around file reads: bounded concurrency for
blocking work.

### Common mistakes
- Assuming every goroutine is cheap, including those blocked in file I/O or cgo.
- Unbounded fan-out over blocking operations. Use a worker pool or semaphore.
- Reading a high thread count in `top` as a leak without checking what the
  threads are blocked in (`cat /proc/<pid>/task/*/stack` as root, or a
  goroutine profile).

### Check yourself
1. What are G, M, and P?
2. Why don't 10,000 goroutines waiting on sockets need 10,000 threads?
3. Why can a slow disk create thousands of threads in a Go program?

### Further reading
- [Go guide, Chapter 28: the scheduler deep dive](../Golang/real-life-golang-guide.md#28-the-scheduler-deep-dive-gmp-and-work-stealing) and
  [Chapter 53: goroutine stacks, revisited](../Golang/real-life-golang-guide.md#53-goroutine-stacks-and-the-scheduler-revisited).
- "Scheduling In Go" (Ardan Labs, three-part series): the GMP model with diagrams.

## Chapter 74 -- Processes from Go: exec, exit codes, signals, and PID 1

### In one sentence
Go starts processes with a combined fork+exec, reports their exit status, and
can catch signals, but when your program is PID 1 in a container it also
inherits the kernel's least-known duty: reaping orphaned processes.

### Why this matters
Go programs routinely run other programs (`git`, `ffmpeg`, health-check
scripts) and are themselves run as a container's first process. Both roles
have traps: zombies that pile up, children that outlive a timeout, signals
that never arrive, and exit codes that hide what happened.

### How it actually works

**Running a program well** with `os/exec`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
cmd := exec.CommandContext(ctx, "ffmpeg", "-i", in, out)  // killed if ctx expires
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}     // own process group...
cmd.Cancel = func() error {                               // ...so we can kill ALL of it,
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)  // including grandchildren
}
out, err := cmd.CombinedOutput()
var ee *exec.ExitError
if errors.As(err, &ee) {
	fmt.Println("exit code", ee.ExitCode())               // -1 means "killed by a signal"
}
```

The `Setpgid` + negative-PID kill matters because killing only the direct
child leaves its children (a shell's pipeline, say) running as orphans.

**PID 1 is special.** When a process exits, its parent must `wait()` for it;
until then it's a **zombie** (Chapter 5): no memory, but a PID-table entry.
When a parent dies first, the kernel re-parents its children to **PID 1**. A
normal init (systemd, tini) loops on `wait()` to reap them. A Go server running
as PID 1 in a container doesn't, unless you write that loop.

Also: the kernel doesn't apply default signal actions to PID 1, so a PID 1
without a SIGTERM handler simply ignores `docker stop` and is SIGKILLed after
the grace period. Go installs handlers when you call `signal.Notify`.

### Build it in Go (20 min) -- make zombies, then stop making them

```go
//go:build linux

// pid1: the job nobody tells you your program has when it is PID 1 in a
// container -- reaping orphaned processes. Run it as a container's entrypoint.
//
//	docker run --rm -v "$PWD/bin":/b debian:bookworm /b/pid1           # zombies pile up
//	docker run --rm -v "$PWD/bin":/b debian:bookworm /b/pid1 -reap     # we reap them
//	docker run --rm --init -v "$PWD/bin":/b debian:bookworm /b/pid1    # tini reaps them
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	reap := flag.Bool("reap", false, "reap orphaned children (what an init process must do)")
	flag.Parse()
	fmt.Println("my PID is", os.Getpid())

	if *reap {
		go reaper()
	}
	for i := 0; i < 5; i++ {
		// sh starts a background sleep and exits at once. The sleep is now an
		// ORPHAN, and the kernel re-parents orphans to PID 1: us.
		cmd := exec.Command("sh", "-c", "sleep 0.2 & exit 0")
		cmd.Run() // we reap sh itself; nobody has waited for the sleep yet
	}
	time.Sleep(time.Second) // the orphaned sleeps exit in the meantime
	fmt.Println("zombie processes now:", zombies())
}

// reaper waits for ANY child whenever SIGCHLD arrives. In a real init this
// loop must not race with exec.Cmd.Wait for children you started yourself --
// one reason to use a tiny dedicated init (tini, dumb-init, docker --init).
func reaper() {
	ch := make(chan os.Signal, 16)
	signal.Notify(ch, syscall.SIGCHLD)
	for range ch {
		for {
			var ws syscall.WaitStatus
			pid, err := syscall.Wait4(-1, &ws, syscall.WNOHANG, nil)
			if pid <= 0 || err != nil {
				break
			}
		}
	}
}

// zombies scans /proc for processes in state Z (exited, never waited for).
func zombies() int {
	n := 0
	stats, _ := filepath.Glob("/proc/[0-9]*/stat")
	for _, f := range stats {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		// format: pid (comm) STATE ...  -- comm may contain spaces, so split after ')'
		if i := strings.LastIndex(string(b), ")"); i > 0 && strings.HasPrefix(string(b[i+2:]), "Z") {
			n++
		}
	}
	return n
}
```

Run it three ways (real output):

```text
$ docker run --rm -v "$PWD/bin":/b debian:bookworm /b/pid1
my PID is 1
zombie processes now: 5

$ docker run --rm -v "$PWD/bin":/b debian:bookworm /b/pid1 -reap
my PID is 1
zombie processes now: 0

$ docker run --rm --init -v "$PWD/bin":/b debian:bookworm /b/pid1
my PID is 7
zombie processes now: 0
```

**What to notice:**

- **Five zombies from five innocent shell commands.** Each `sh -c "sleep 0.2 &"`
  left an orphaned `sleep` that became ours. A service that runs scripts
  thousands of times a day in a container eventually runs out of PIDs.
- **`docker --init`** puts `tini` in as PID 1. It reaps zombies and forwards
  signals, and our program became PID 7. In Kubernetes, use an init binary in
  the image (`tini`, `dumb-init`), or share the pod's process namespace with
  `shareProcessNamespace: true` so the pause container reaps.
- **The `-reap` loop races with `exec.Cmd.Wait`**: it can collect a child that
  `cmd.Run` was waiting for, so `Run` returns "no child processes". That race
  is the main reason to use a dedicated init rather than reaping inside your
  application.

**Exercises:**

1. Run `/b/pid1` with `docker run` and, from another terminal,
   `docker exec <id> ps -eo pid,ppid,stat,comm`. Find the `Z` entries and their
   parent.
2. Remove the `signal.Notify` from Chapter 19 of the HTTPS guide's server and
   run it as PID 1. How long does `docker stop` take now? (10 s, Docker's
   default grace period, then SIGKILL.)
3. Write `runWithTimeout(cmd string, d time.Duration)` that kills the whole
   process group on timeout, and test it with
   `sh -c "sleep 100 & sleep 100"`.

### Common mistakes
- Running a Go binary as PID 1 that starts subprocesses, with no init.
- Killing only the child PID on timeout, leaving grandchildren running.
- Treating `ExitCode() == -1` as a normal failure (it means a signal killed it;
  check `ee.Sys().(syscall.WaitStatus).Signal()`).
- No SIGTERM handling in PID 1, so every shutdown is a SIGKILL.

### Check yourself
1. Why do orphans become PID 1's responsibility?
2. What does `docker run --init` change?
3. Why use `Setpgid` before killing a command on timeout?

### Further reading
- `man 2 wait`, `man 7 signal`; the tini README ("What is advantage of
  Tini?") explains the PID 1 problem in a page.
- Chapter 5 (process states) and Chapter 20 (signals) of this guide.

## Chapter 75 -- CPU limits: GOMAXPROCS, cgroups, and throttling

### In one sentence
A container's CPU limit is a time quota per scheduling period, not a smaller
number of CPUs, and a runtime that starts more busy threads than the quota
covers gets **throttled**, which shows up as latency spikes, not high CPU.

### Why this matters
`--cpus=1.5` or a Kubernetes `limits.cpu: 1500m` becomes cgroup v2's
`cpu.max = 150000 100000`: 150 ms of CPU time per 100 ms period, *spread over
any number of CPUs*. The process still sees every CPU on the host
(`nproc` = 14 below). Before Go 1.25, Go set GOMAXPROCS to that host CPU count,
burned through the quota in the first part of each period, and then stalled for
the rest. Tail latency rose while average CPU looked fine. Java, Node.js
worker pools, and anything else that sizes itself from the CPU count can make
the same mistake.

### How it actually works

```
cpu.max = 150000 100000          quota 150 ms per 100 ms period (= 1.5 CPUs)

GOMAXPROCS=8, all busy:
  |########--------------------------------| 8 threads x ~19 ms = 150 ms used by t=19ms
  |        ^ throttled for the remaining ~81 ms of every period
GOMAXPROCS=2, all busy:
  |####################################----| 2 threads share the quota evenly; little or no stall
```

Go 1.25 made GOMAXPROCS **container-aware**: on Linux it reads the cgroup CPU
limit, rounds it **up**, applies a **minimum of 2**, and re-checks
periodically in case the limit changes. Measured with Go 1.26: `--cpus=0.5`,
`1`, and `1.5` all gave GOMAXPROCS = 2, and `--cpus=2.2` gave 3. The
`GOMAXPROCS` environment variable still overrides it.

### Build it in Go (15 min) -- measure throttling

```go
// gomaxprocs: what does a Go program think it can use, versus what the
// container's cgroup actually allows?
//
//	go run ./gomaxprocs
//	docker run --rm --cpus=1.5 -v "$PWD":/w -w /w golang:1.26 go run ./gomaxprocs
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("runtime.NumCPU()   =", runtime.NumCPU(), " (CPUs the OS lets this process run on)")
	fmt.Println("GOMAXPROCS         =", runtime.GOMAXPROCS(0), " (OS threads running Go code at once)")
	if b, err := os.ReadFile("/sys/fs/cgroup/cpu.max"); err == nil {
		f := strings.Fields(string(b)) // "<quota> <period>" or "max <period>"
		fmt.Printf("cgroup cpu.max     = %s", b)
		if f[0] != "max" {
			fmt.Printf("                   = %s / %s µs per period -> the container's CPU limit\n", f[0], f[1])
		}
	} else {
		fmt.Println("cgroup cpu.max     = (not in a cgroup v2 container)")
	}

	// Burn CPU on every P for 2 seconds and see how much we got.
	before := throttled()
	var wg sync.WaitGroup
	stop := time.Now().Add(2 * time.Second)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(stop) {
			}
		}()
	}
	wg.Wait()
	after := throttled()
	if before != nil && after != nil {
		fmt.Printf("during 2s of busy work: throttled %d times, for %.0f ms in total\n",
			after[0]-before[0], float64(after[1]-before[1])/1000)
	}
}

// throttled reads nr_throttled and throttled_usec from cgroup v2's cpu.stat.
func throttled() []int64 {
	b, err := os.ReadFile("/sys/fs/cgroup/cpu.stat")
	if err != nil {
		return nil
	}
	var n, us int64
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) == 2 && f[0] == "nr_throttled" {
			fmt.Sscan(f[1], &n)
		}
		if len(f) == 2 && f[0] == "throttled_usec" {
			fmt.Sscan(f[1], &us)
		}
	}
	return []int64{n, us}
}
```

Real output (Go 1.26, host with 14 CPUs):

```text
$ docker run --rm -v "$PWD/bin":/b debian:bookworm /b/gomaxprocs                   # no limit
runtime.NumCPU()   = 14
GOMAXPROCS         = 14
cgroup cpu.max     = max 100000
during 2s of busy work: throttled 0 times, for 0 ms in total

$ docker run --rm --cpus=1.5 -v "$PWD/bin":/b debian:bookworm /b/gomaxprocs
runtime.NumCPU()   = 14  (CPUs the OS lets this process run on)
GOMAXPROCS         = 2   (OS threads running Go code at once)
cgroup cpu.max     = 150000 100000
during 2s of busy work: throttled 20 times, for 1011 ms in total

$ docker run --rm --cpus=1.5 -e GOMAXPROCS=8 -v "$PWD/bin":/b debian:bookworm /b/gomaxprocs
GOMAXPROCS         = 8
during 2s of busy work: throttled 21 times, for 13658 ms in total
```

**Read it:**

- **Go picked GOMAXPROCS = 2** for a 1.5-CPU limit (rounded up), while
  `NumCPU()` still said 14.
- **`throttled_usec` is summed across threads.** With 8 busy threads, 13.7 s of
  thread-time was spent frozen during a 2 s run. Each thread was stopped for
  most of every period. Requests arriving then wait up to 80 ms for no reason
  visible in the application.
- **Some throttling remains at GOMAXPROCS = 2**, because 2 threads on a 1.5-CPU
  quota still overrun it slightly. For latency-critical services, some teams
  set CPU *requests* and no CPU *limit*, letting the pod use idle CPU without
  quota stalls. Measure before choosing.

**Exercises:**

1. Watch throttling on a real Kubernetes pod:
   `kubectl exec <pod> -- cat /sys/fs/cgroup/cpu.stat` and compare
   `nr_throttled` before and after a load test.
2. Run with `--cpus=0.5`. What does Go choose for GOMAXPROCS, and how much
   throttling do you see? (Go never chooses fewer than 2, so a half-CPU
   container is throttled heavily by design. Weigh that against setting
   `GOMAXPROCS=1` explicitly.)
3. Build the same program with Go 1.24 (`docker run golang:1.24 ...`) and
   compare GOMAXPROCS under `--cpus=1.5`.

### Common mistakes
- Reading "CPU usage 60% of limit" as "plenty of headroom" while
  `nr_throttled` climbs. Throttling happens within each 100 ms period, not on
  average.
- Running Go older than 1.25 in containers without `automaxprocs` or an
  explicit `GOMAXPROCS`.
- Sizing worker pools from `runtime.NumCPU()` instead of `runtime.GOMAXPROCS(0)`.

### Check yourself
1. What do the two numbers in `cpu.max` mean?
2. Why can a container be throttled while its average CPU is below its limit?
3. What did Go 1.25 change, and how do you override it?

### Further reading
- Kernel documentation: "Control Group v2", the `cpu.max` and `cpu.stat`
  sections.
- Go 1.25 release notes: "Container-aware GOMAXPROCS".
- Chapter 58 (scheduling and resources) of this guide.

## Chapter 76 -- Memory limits: the Go heap, GOMEMLIMIT, and the OOM killer

### In one sentence
By default Go lets its heap grow to about twice the live data before
collecting, so in a container a program whose working set fits comfortably
can still be OOM-killed. `GOMEMLIMIT` tells the garbage collector about the
limit.

### Why this matters
"OOMKilled, exit code 137" (Chapter 13) is one of the most common container
failures, and for garbage-collected languages it often isn't a leak. The GC
was pacing itself as if memory were unlimited. Knowing the two knobs (`GOGC`
and `GOMEMLIMIT`) turns a mysterious crash into a configuration decision.

### How it actually works
- **`GOGC=100`** (the default): run the next GC when the heap has grown 100%
  beyond the live heap after the last GC. Live 160 MB means the next GC at
  about 320 MB.
- **`GOMEMLIMIT`** (Go 1.19+): a *soft* limit on the runtime's total memory.
  As memory approaches it, the GC runs more often, whatever `GOGC` says.
- **The cgroup limit** (`memory.max`): a *hard* limit enforced by the kernel. On
  reaching it, the kernel reclaims what it can and then OOM-kills.

Set `GOMEMLIMIT` to roughly 90% of the container's memory limit, leaving
room for non-heap memory: goroutine stacks, runtime metadata, and C
libraries in cgo programs.

### Build it in Go (20 min) -- an OOM kill and its cure

The program keeps 160 MB live (a "cache") and allocates 4 MB of garbage in a
tight loop, like a busy service. It's built on the host and run in a 300 MB
container. (Building *inside* that container fails: the Go compiler itself was
OOM-killed when we tried, which is another lesson.)

```go
// memlimit: why Go services get OOM-killed in containers, and how GOMEMLIMIT
// prevents it. The program keeps `-live` MB permanently reachable and churns
// through short-lived garbage, like a busy service with a cache.
//
//	docker run --rm --memory=300m -v "$PWD":/w -w /w golang:1.26 go run ./memlimit
//	docker run --rm --memory=300m -e GOMEMLIMIT=250MiB -v "$PWD":/w -w /w golang:1.26 go run ./memlimit
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var sink []byte // defeats escape analysis: garbage must really be allocated

func main() {
	live := flag.Int("live", 160, "MB kept alive for the whole run (the 'cache')")
	secs := flag.Int("secs", 8, "how long to run")
	flag.Parse()

	fmt.Printf("GOMEMLIMIT=%q  GOGC=%q  container limit=%s\n",
		os.Getenv("GOMEMLIMIT"), os.Getenv("GOGC"), cgroup("memory.max"))

	cache := make([][]byte, 0, *live)
	for i := 0; i < *live; i++ {
		b := make([]byte, 1<<20)
		for j := range b {
			b[j] = byte(j) // touch every page so it is really resident
		}
		cache = append(cache, b)
	}

	end := time.Now().Add(time.Duration(*secs) * time.Second)
	next := time.Now()
	for time.Now().Before(end) {
		sink = make([]byte, 4<<20) // 4 MB of garbage per iteration
		for j := 0; j < len(sink); j += 4096 {
			sink[j] = 1
		}
		if time.Now().After(next) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("heap=%4d MB  next GC at %4d MB  GCs=%4d  cgroup memory.current=%4d MB\n",
				m.HeapAlloc>>20, m.NextGC>>20, m.NumGC, atoiMB(cgroup("memory.current")))
			next = time.Now().Add(time.Second)
		}
	}
	runtime.KeepAlive(cache)
	fmt.Println("survived. limit from debug.SetMemoryLimit:", debug.SetMemoryLimit(-1)>>20, "MB")
}

func cgroup(file string) string {
	b, err := os.ReadFile("/sys/fs/cgroup/" + file)
	if err != nil {
		return "n/a"
	}
	return strings.TrimSpace(string(b))
}

func atoiMB(s string) int64 {
	var n int64
	fmt.Sscan(s, &n)
	return n >> 20
}
```

```text
$ docker run --rm --memory=300m --memory-swap=300m -v "$PWD/bin":/b debian:bookworm /b/memlimit
GOMEMLIMIT=""  GOGC=""  container limit=314572800
heap= 164 MB  next GC at  280 MB  GCs=   6  cgroup memory.current= 173 MB
Killed
exit code: 137

$ docker run --rm --memory=300m --memory-swap=300m -e GOMEMLIMIT=250MiB -v "$PWD/bin":/b debian:bookworm /b/memlimit
heap= 164 MB  next GC at  236 MB  GCs=   6  cgroup memory.current= 173 MB
heap= 228 MB  next GC at  236 MB  GCs=1508  cgroup memory.current= 256 MB
heap= 228 MB  next GC at  236 MB  GCs=3006  cgroup memory.current= 254 MB
...
heap= 180 MB  next GC at  236 MB  GCs=10666  cgroup memory.current= 248 MB
survived. limit from debug.SetMemoryLimit: 250 MB
exit code: 0
```

**What to notice:**

- **Exit code 137 = 128 + 9 (SIGKILL).** The kernel killed the process the
  moment the cgroup hit 300 MB. Go printed no error, because it never got the
  chance. "Next GC at 280 MB" plus runtime overhead was already too close to
  the limit.
- **With `GOMEMLIMIT` the target dropped to 236 MB**, and memory stayed under
  260 MB for the whole run.
- **It ran about 1,500 GCs per second.** The limit traded memory for CPU. If
  the live heap really needs most of the limit, `GOMEMLIMIT` turns an OOM kill
  into a "death spiral" of constant GC. Go caps GC CPU at about 50% to
  prevent the worst of it, but the service still slows down badly. The real fix
  is then a bigger limit or a smaller cache.

**Exercises:**

1. Try `-live 100` without `GOMEMLIMIT`. Does it survive? Where does
   `memory.current` peak?
2. Run with `GODEBUG=gctrace=1` and read one line: heap before → after → live,
   and the GC CPU percentage.
3. In Kubernetes, set `GOMEMLIMIT` from the manifest (Chapter 82's manifest
   does) and alert when `go_gc_duration_seconds` rises sharply, which is the
   early warning of a death spiral.

### Common mistakes
- Treating every OOM kill of a Go service as a memory leak.
- Setting `GOMEMLIMIT` equal to the container limit, which leaves no room for
  non-heap memory.
- Setting `GOGC=off` with `GOMEMLIMIT` without understanding that the program
  then collects *only* near the limit.
- Building inside the same memory-limited container that runs the service.

### Check yourself
1. With `GOGC=100` and 200 MB of live data, roughly when does the next GC run?
2. Why is `GOMEMLIMIT` called a soft limit?
3. What does exit code 137 tell you?

### Further reading
- "A Guide to the Go Garbage Collector" (go.dev/doc/gc-guide), with
  interactive graphs of `GOGC` and the memory limit.
- [Go guide, Chapter 27: the memory model and the garbage collector](../Golang/real-life-golang-guide.md#27-the-go-memory-model-and-the-garbage-collector).
- Chapter 13 (allocation and the OOM killer) of this guide.

## Chapter 77 -- Files that survive crashes: page cache, fsync, and atomic replacement

### In one sentence
A successful `write()` only means the data reached the page cache. Making a
change survive a crash takes `fsync`, and making a *replacement* all-or-nothing
takes write-to-temp, fsync, rename, then fsync of the directory.

### Why this matters
Config files, state files, caches, and embedded databases get written by
ordinary application code. After a power cut or kernel panic, the classic
results are a zero-length config file, a half-written JSON file the service
can't parse at startup, or "the save succeeded" but the data is gone.
Chapters 14–16 explained the page cache and journaling. This chapter is the
application side of that contract.

### How it actually works

```
os.WriteFile(path, data)       -> bytes in the PAGE CACHE; on disk "eventually" (seconds later)
f.Sync()   (fsync)             -> this file's data AND metadata on stable storage
os.Rename(tmp, path)           -> atomic: any reader sees the old file or the new one, never a mix
dir.Sync() (fsync on the dir)  -> the rename itself (a directory entry) is durable
```

- **Why not just `O_TRUNC` and write?** Between the truncate and the last write
  there's a window in which the file is empty or partial, and a crash inside
  that window leaves it that way.
- **Why the same directory?** `rename(2)` is only atomic within one filesystem.
  Across filesystems it's a copy.
- **Why fsync the directory?** The new name lives in the directory's data. Until
  that's flushed, a crash can bring back the old name, or no name at all.

### Build it in Go (25 min) -- `WriteFileAtomic`, and what safety costs

```go
// atomicwrite: replace a file so that, after a crash or power cut, readers
// see EITHER the complete old version OR the complete new version -- never a
// half-written mix. Also measures what durability costs.
//
//	go run ./atomicwrite
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WriteFileAtomic is the standard recipe used by databases, package managers,
// and config tools:
//  1. write to a temporary file IN THE SAME DIRECTORY (same filesystem)
//  2. fsync it                  -> the new bytes are on stable storage
//  3. rename over the target    -> atomic: the name points at old OR new
//  4. fsync the directory       -> the rename itself is on stable storage
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tmp-"+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // no-op after a successful rename
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil { // fsync(2)
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil { // rename(2): atomic within a filesystem
		return err
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync() // fsync the directory entry
}

func main() {
	dir, _ := os.MkdirTemp("", "atomic")
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "config.json")
	data := []byte(`{"version": 2, "replicas": 3}` + "\n")

	// What does each level of safety cost? 200 small writes each.
	measure := func(label string, write func(i int) error) {
		start := time.Now()
		for i := 0; i < 200; i++ {
			if err := write(i); err != nil {
				fmt.Println(label, err)
				return
			}
		}
		fmt.Printf("%-44s %8.3f ms per write\n", label, float64(time.Since(start).Microseconds())/200/1000)
	}
	measure("os.WriteFile (page cache only, NOT durable)", func(int) error {
		return os.WriteFile(path, data, 0o644)
	})
	measure("write + fsync (durable, NOT atomic)", func(int) error {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.Write(data); err != nil {
			return err
		}
		return f.Sync()
	})
	measure("WriteFileAtomic (durable AND atomic)", func(int) error {
		return WriteFileAtomic(path, data, 0o644)
	})
	b, _ := os.ReadFile(path)
	fmt.Printf("final content: %s", b)
}
```

Real output:

```text
== Linux (container, overlay filesystem)
os.WriteFile (page cache only, NOT durable)      0.070 ms per write
write + fsync (durable, NOT atomic)              0.699 ms per write
WriteFileAtomic (durable AND atomic)             1.047 ms per write

== macOS (APFS, Apple SSD)
os.WriteFile (page cache only, NOT durable)      0.040 ms per write
write + fsync (durable, NOT atomic)              4.624 ms per write
WriteFileAtomic (durable AND atomic)             9.540 ms per write
```

**What to notice:**

- **Durability costs 10× to 100×.** That's why databases batch many commits
  into one fsync ("group commit"), and why fsync-per-request designs don't
  scale.
- **macOS is slower for a good reason:** Go's `File.Sync` on macOS uses
  `F_FULLFSYNC`, which forces the drive's own cache to flush. Plain `fsync`
  there doesn't guarantee that.
- **Cloud and container numbers vary widely.** Network block storage often
  takes milliseconds per fsync, so measure on the volume you actually deploy to.
- **The `defer os.Remove(tmp)`** cleans up a failed attempt. After a successful
  rename it's a harmless no-op.

**Exercises:**

1. Use `WriteFileAtomic` in a program that saves its state on SIGTERM
   (Chapter 74), and kill it repeatedly in a loop. The state file must always
   parse.
2. Write 1,000 records with an fsync after each, then with one fsync per 100.
   That's the group-commit trade-off between throughput and how much can be
   lost.
3. Recreate Chapter 15's "deleted but still open" mystery from Go: open a file,
   `os.Remove` it, keep writing, and watch `df` and
   `ls -l /proc/<pid>/fd`.

### Common mistakes
- Assuming `Close()` implies durability (it doesn't).
- Writing the temporary file in `/tmp`, which is often a different filesystem
  (tmpfs), so the rename becomes a non-atomic copy.
- Ignoring the error from `Sync()` or `Close()`. On some filesystems that's
  where write errors are reported.
- Calling fsync for data that's fine to lose (caches), and paying for nothing.

### Check yourself
1. What does a successful `write()` guarantee?
2. Why must the temporary file be in the same directory as the target?
3. What does fsyncing the directory protect against?

### Further reading
- "Files are hard" (Dan Luu): a survey of what really happens to file data
  across crashes.
- `man 2 fsync`, `man 2 rename`; Chapters 14–16 of this guide.

## Chapter 78 -- File descriptors and the netpoller: holding thousands of connections

### In one sentence
Every socket, file, and pipe is a file descriptor, the per-process limit
(`RLIMIT_NOFILE`) caps how many a server can hold, and Go's netpoller lets one
OS thread wait on all of them with epoll.

### Why this matters
"Too many open files" (`EMFILE`) is the classic failure of a busy server or a
leaky one: every new connection fails while the existing ones keep working.
Chapter 43 introduced `ulimit -n`. Go adds a twist that confuses people
comparing a Go service with the shell that started it.

### How it actually works
- **Two limits:** the *soft* limit is what's enforced, and the *hard* limit is the
  ceiling an unprivileged process can raise its soft limit to.
- **Go raises its own soft limit at startup** to the hard limit **minus one**
  (since Go 1.19; the Go source explains that the −1 lets Go notice if another
  process later changes its limit with `prlimit`). Go then restores the
  *original* soft limit for any program it starts with `os/exec`, so old
  software that breaks above 1,024 descriptors isn't affected.
- **epoll:** the runtime registers every network descriptor with one epoll
  instance. A dedicated wait (`epoll_pwait`) wakes exactly the goroutines whose
  sockets are ready. That's why Chapter 73 measured 10,000 socket-blocked
  goroutines on 16 threads.

### Build it in Go (15 min) -- hit EMFILE on purpose

```go
//go:build linux

// fdlimit: every socket and file is a file descriptor, and the per-process
// limit (RLIMIT_NOFILE, `ulimit -n`) caps how many connections a server can
// hold. Go quietly raises the soft limit to the hard limit at startup.
//
//	docker run --rm --ulimit nofile=256:2048 -v "$PWD/bin":/b debian:bookworm /b/fdlimit
package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func main() {
	var rl syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rl)
	fmt.Printf("RLIMIT_NOFILE as Go sees it: soft=%d hard=%d\n", rl.Cur, rl.Max)
	fmt.Println("(Go raised the soft limit to hard-1 at startup; programs it execs get the original back)")

	var files []*os.File
	for {
		f, err := os.Open("/dev/null")
		if err != nil {
			fmt.Printf("open #%d failed: %v\n", len(files)+1, err)
			fmt.Println("EMFILE?", errors.Is(err, syscall.EMFILE), "-- 'too many open files'")
			break
		}
		files = append(files, f)
	}
	files[0].Close() // reading /proc/self/fd needs one free descriptor itself
	fds, _ := os.ReadDir("/proc/self/fd")
	fmt.Printf("descriptors open: %d (%d files + stdin/out/err + the runtime's epoll and wakeup fds)\n", len(fds), len(files)-1)
}
```

```text
$ docker run --rm --ulimit nofile=256:2048 -v "$PWD/bin":/b debian:bookworm sh -c 'ulimit -Sn; /b/fdlimit'
256                                                         <- the shell's soft limit
RLIMIT_NOFILE as Go sees it: soft=2047 hard=2048            <- Go raised it: hard - 1
(Go raised the soft limit to hard-1 at startup; programs it execs get the original back)
open #2042 failed: open /dev/null: too many open files
EMFILE? true -- 'too many open files'
descriptors open: 2047 (2040 files + stdin/out/err + the runtime's epoll and wakeup fds)
```

**What to notice:**

- **The shell said 256, and Go got 2,047.** If you size a Go service by
  `ulimit -n` in the shell, you'll be wrong. Read
  `/proc/<pid>/limits` for the real value.
- **The `ReadDir` of `/proc/self/fd` needs a descriptor too.** The first draft of
  this lab printed "0 open" because every descriptor was in use, so even the
  diagnostic failed. In production, logging libraries that open files hit the
  same wall during an `EMFILE` incident. Keep a little headroom.
- **Four descriptors belong to the runtime** (epoll and its wakeup pipe/eventfd).

**Exercises:**

1. Replace `os.Open("/dev/null")` with `net.Dial` to a local listener and count
   how many connections fit. Then set `LimitNOFILE=` in a systemd unit
   (Chapter 82) and check `/proc/<pid>/limits`.
2. Add a leak: open a file per HTTP request and "forget" to close it. Watch
   `ls /proc/<pid>/fd | wc -l` climb under load, then find the leak with
   `lsof -p <pid>`.
3. Use `ss -s` while running the TCP/IP guide's
   [`closewait` lab](../networking/tcp-ip/real-life-guide-v1.md#chapter-23-tcp-part-3-closing-a-connection-and-the-states): each CLOSE_WAIT socket is a leaked
   descriptor.

### Common mistakes
- Checking the shell's `ulimit -n` instead of the running process's
  `/proc/<pid>/limits`.
- Raising the limit to hide a descriptor leak.
- Forgetting that every client connection *and* every upstream connection a
  proxy holds counts.

### Check yourself
1. What's the difference between the soft and hard limits?
2. What does Go do to `RLIMIT_NOFILE` at startup, and for its children?
3. How does the netpoller avoid one thread per connection?

### Further reading
- `man 2 getrlimit`, `man 7 epoll`.
- The comment at the top of Go's `src/syscall/rlimit.go` (go.dev/issue/46279)
  explains the startup adjustment.

## Chapter 79 -- Reading /proc from Go: build your own ps

### In one sentence
`ps`, `top`, `free`, and most monitoring agents are parsers for files under
`/proc`. Reading them yourself shows exactly what each number means and where
the parsing traps are.

### Why this matters
When a container has no `ps`, when you need one number in a health check, or
when you're building an exporter or debugging an agent, you read `/proc`
directly. Chapter 23 showed the files from the shell. This chapter shows the
format details that break naive parsers.

### How it actually works

| File | What it holds | Trap |
|---|---|---|
| `/proc/<pid>/stat` | one line: state, CPU ticks, threads, start time... | field 2 is `(comm)` and **may contain spaces and parentheses**: split after the *last* `)` |
| `/proc/<pid>/status` | human-readable: `VmRSS`, `Threads`, `voluntary_ctxt_switches` | values carry units (`kB`) |
| `/proc/<pid>/cmdline` | argv, **NUL-separated** | empty for kernel threads |
| `/proc/<pid>/fd/` | one symlink per open descriptor | needs permission for other users' processes |
| `/proc/<pid>/limits` | effective rlimits | the truth about Chapter 78's limits |

CPU% requires **two samples**: `(ticks_after − ticks_before) / USER_HZ / interval`.
`USER_HZ` is 100 on practically every Linux system.

### Build it in Go (30 min) -- `gops`

```go
//go:build linux

// gops: `ps` and a little of `top`, written from /proc alone -- everything
// those tools show comes from these files.
//
//	go run ./gops            # top 10 processes by CPU over 1 second
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type proc struct {
	pid      int
	comm     string
	state    string
	threads  int
	rssKB    int64
	cpuTicks int64 // utime + stime, in clock ticks
	cmdline  string
}

// read parses /proc/<pid>/stat (one line, space-separated, but the command
// name is in parentheses and may itself contain spaces) and /proc/<pid>/status.
func read(pid int) (*proc, error) {
	b, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return nil, err
	}
	s := string(b)
	l, r := strings.IndexByte(s, '('), strings.LastIndexByte(s, ')')
	f := strings.Fields(s[r+2:]) // fields after the comm: f[0] is field 3 ("state")
	p := &proc{pid: pid, comm: s[l+1 : r], state: f[0]}
	utime, _ := strconv.ParseInt(f[11], 10, 64) // field 14
	stime, _ := strconv.ParseInt(f[12], 10, 64) // field 15
	p.threads, _ = strconv.Atoi(f[17])          // field 20
	p.cpuTicks = utime + stime

	st, _ := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	for _, line := range strings.Split(string(st), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			p.rssKB, _ = strconv.ParseInt(strings.Fields(line)[1], 10, 64)
		}
	}
	cl, _ := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid)) // NUL-separated argv
	p.cmdline = strings.TrimSpace(strings.ReplaceAll(string(cl), "\x00", " "))
	if p.cmdline == "" {
		p.cmdline = "[" + p.comm + "]" // kernel threads have no cmdline
	}
	return p, nil
}

func snapshot() map[int]*proc {
	out := map[int]*proc{}
	dirs, _ := filepath.Glob("/proc/[0-9]*")
	for _, d := range dirs {
		pid, _ := strconv.Atoi(filepath.Base(d))
		if p, err := read(pid); err == nil { // processes can vanish mid-scan: ignore
			out[pid] = p
		}
	}
	return out
}

func main() {
	const hz = 100 // USER_HZ: clock ticks per second on virtually all Linux systems
	interval := time.Second
	before := snapshot()
	time.Sleep(interval)
	after := snapshot()

	type row struct {
		p   *proc
		cpu float64
	}
	var rows []row
	for pid, p := range after {
		var delta int64
		if b, ok := before[pid]; ok {
			delta = p.cpuTicks - b.cpuTicks
		}
		rows = append(rows, row{p, float64(delta) / hz / interval.Seconds() * 100})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].cpu > rows[j].cpu })

	states := map[string]int{}
	for _, r := range rows {
		states[r.p.state]++
	}
	fmt.Printf("%d processes; states: %v  (R running, S sleeping, D disk wait, Z zombie)\n\n", len(rows), states)
	fmt.Printf("%7s %-1s %6s %5s %9s  %s\n", "PID", "S", "%CPU", "THR", "RSS(MB)", "COMMAND")
	for _, r := range rows[:min(10, len(rows))] {
		cmd := r.p.cmdline
		if len(cmd) > 50 {
			cmd = cmd[:50] + "..."
		}
		fmt.Printf("%7d %-1s %6.1f %5d %9.1f  %s\n", r.p.pid, r.p.state, r.cpu, r.p.threads, float64(r.p.rssKB)/1024, cmd)
	}
}
```

Real output (in a container with `--pid=host`, a busy loop running):

```text
311 processes; states: map[I:126 R:2 S:183]  (R running, S sleeping, D disk wait, Z zombie)

    PID S   %CPU   THR   RSS(MB)  COMMAND
  82877 R  101.0     1       0.5  sh -c (while :; do :; done) & sleep 0.5; /b/gops
  34206 S    0.0     1       3.5  nginx: worker process
    359 I    0.0     1       0.0  [kworker/11:1H-kblockd]
   1114 S    0.0     1      14.6  postgres: checkpointer
    295 S    0.0    20      88.9  /usr/bin/containerd --config /etc/containerd/conta...
```

**What to notice:**

- **`I` (idle)** is a kernel-thread state for idle workers: 126 of them here, all
  in square brackets because kernel threads have no command line.
- **101% CPU** is one fully busy core, plus sampling error. `top` would show the
  same thing.
- **Processes vanish mid-scan.** The code ignores read errors for that reason:
  any `/proc` walker must tolerate processes exiting between listing and reading.

**Exercises:**

1. Add a `-p PID` mode that prints the process's open descriptors (resolving
   `/proc/PID/fd/*` symlinks) and its `/proc/PID/limits`.
2. Add a zombie report: every `Z` process with its parent (field 4, `ppid`).
   Run it next to Chapter 74's `pid1` lab.
3. Export the numbers in Prometheus format (the TCP/IP guide's
   [`prober` lab](../networking/tcp-ip/real-life-guide-v1.md#chapter-46-where-the-numbers-come-from-counters-flows-and-probes) shows how).

### Common mistakes
- Splitting `/proc/<pid>/stat` on spaces from the start of the line.
- Computing CPU% from a single sample.
- Treating a failed read as fatal during a `/proc` scan.

### Check yourself
1. Why must the `stat` parser split after the last `)`?
2. How do you compute a process's CPU% from `/proc`?
3. Why do some processes have an empty `cmdline`?

### Further reading
- `man 5 proc`, the authoritative field-by-field reference.
- Chapter 23 of this guide.

## Chapter 80 -- A container runtime in 150 lines of Go

### In one sentence
A container is a process started with new namespaces (its own hostname, PID
numbering, mounts, IPC, and network), a different root filesystem, and a
cgroup that limits its resources. A short Go program can do all of it.

### Why this matters
Chapter 44 said "a container is namespaces + cgroups + a filesystem". Building
one makes that concrete, and it explains real behaviour: why PID 1 inside a
container is special (Chapter 74), why a container's `localhost` isn't the
host's, why `ps` inside shows so little, and why fork bombs are contained. This
lab follows the well-known "containers from scratch" exercise popularised by
Liz Rice.

### How it actually works

```
 run (parent, on the host)
  ├─ mkdir /sys/fs/cgroup/minicontainer-<pid>; write pids.max, memory.max
  └─ re-exec /proc/self/exe "child" with:
       Cloneflags: NEWUTS | NEWPID | NEWNS | NEWIPC | NEWNET
       CgroupFD:   the new cgroup              (Go 1.20+: born inside it, no race)
 child (PID 1 in its own PID namespace)
  ├─ sethostname("minicontainer")
  ├─ chroot(rootfs); chdir("/")
  ├─ mount proc on /proc                      (so ps sees only this namespace)
  └─ run the user's command
```

Real runtimes (runc, crun) add `pivot_root` instead of `chroot` (which can be
escaped by a root process), user namespaces, seccomp filters, capability
dropping, a network setup (veth pair + bridge, as in the
[TCP/IP guide's Chapter 52](../networking/tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack)), and OCI image unpacking.

### Build it in Go (45 min) -- `minicontainer`

```go
//go:build linux

// minicontainer: a container runtime in ~150 lines. Namespaces give the
// process its own view of the system; chroot gives it its own filesystem;
// a cgroup limits what it can consume. That's all a container is.
//
// Run as root on a Linux machine/VM (cgroup v2), with a root filesystem:
//
//	mkdir -p /tmp/rootfs && docker export $(docker create busybox) | tar -x -C /tmp/rootfs
//	go build -o minicontainer ./minicontainer
//	sudo ./minicontainer run /tmp/rootfs /bin/sh
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: minicontainer run <rootfs> <command> [args...]")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "run":
		run(os.Args[2], os.Args[3:])
	case "child":
		child(os.Args[2], os.Args[3:])
	}
}

// run is the "runtime": it re-executes this same binary as `child` inside new
// namespaces and a new cgroup, and waits for it.
func run(rootfs string, argv []string) {
	cg, err := makeCgroup()
	must(err)
	defer os.Remove(cg) // rmdir: only works once every process in it has exited
	cgfd, err := syscall.Open(cg, syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
	must(err)
	defer syscall.Close(cgfd)

	cmd := exec.Command("/proc/self/exe", append([]string{"child", rootfs}, argv...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // own hostname
			syscall.CLONE_NEWPID | // own PID numbering: the child is PID 1
			syscall.CLONE_NEWNS | // own mount table
			syscall.CLONE_NEWIPC | // own System V IPC / message queues
			syscall.CLONE_NEWNET, // own network stack: only a loopback, and it's down
		Unshareflags: syscall.CLONE_NEWNS, // don't propagate our mounts back to the host
		UseCgroupFD:  true,                // Go 1.20+: start the child directly inside the cgroup
		CgroupFD:     cgfd,
	}
	if err := cmd.Run(); err != nil {
		fmt.Println("container exited:", err)
	}
	if b, err := os.ReadFile(filepath.Join(cg, "pids.events")); err == nil {
		fmt.Printf("cgroup pids.events: %s", b) // "max N" = N forks refused by the limit
	}
}

// child runs INSIDE the new namespaces, as PID 1 of its own PID namespace.
func child(rootfs string, argv []string) {
	must(syscall.Sethostname([]byte("minicontainer")))
	must(syscall.Chroot(rootfs)) // production runtimes use pivot_root, which can't be escaped like chroot
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "/proc", "proc", 0, "")) // a /proc that only shows OUR PID namespace
	defer syscall.Unmount("/proc", 0)

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin", "PS1=container# "}
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}

// makeCgroup creates a cgroup v2 group with a process limit (a fork-bomb
// guard) and a memory limit. Requires root and the pids + memory controllers
// enabled in the parent's cgroup.subtree_control.
func makeCgroup() (string, error) {
	cg := filepath.Join("/sys/fs/cgroup", "minicontainer-"+strconv.Itoa(os.Getpid()))
	if err := os.Mkdir(cg, 0o755); err != nil {
		return "", fmt.Errorf("create cgroup (are you root, on cgroup v2?): %w", err)
	}
	for file, value := range map[string]string{
		"pids.max":   "20",       // at most 20 processes
		"memory.max": "64000000", // 64 MB
	} {
		if err := os.WriteFile(filepath.Join(cg, file), []byte(value), 0o644); err != nil {
			os.Remove(cg)
			return "", fmt.Errorf("%s: %w (enable the controller in the parent's cgroup.subtree_control)", file, err)
		}
	}
	return cg, nil
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
```

Run it as root on a Linux VM, or inside a privileged container as below:

```bash
mkdir -p rootfs && docker export $(docker create busybox) | tar -x -C rootfs
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/minicontainer ./minicontainer
docker run --rm -it --privileged --cgroupns=private -v "$PWD":/w debian:bookworm sh -c '
  mkdir /sys/fs/cgroup/init && echo $$ > /sys/fs/cgroup/init/cgroup.procs   # cgroup v2: move ourselves
  echo "+pids +memory" > /sys/fs/cgroup/cgroup.subtree_control                # into a leaf, enable controllers
  /w/bin/minicontainer run /w/rootfs /bin/sh'
```

(The two `cgroup` lines are needed only inside a container. cgroup v2 doesn't
let a group both contain processes and delegate controllers, so the shell moves
itself into a leaf group first. On a systemd host, run as root and skip them.)

Real output from a scripted run:

```text
== host view: hostname=3b127053433e pid=1
inside: hostname=minicontainer my-pid=7
PID   USER     TIME  COMMAND
    1 root      0:00 /proc/self/exe child /b/rootfs /bin/sh -c ...
    7 root      0:00 /bin/sh -c ...
    9 root      0:00 ps

1: lo: <LOOPBACK> mtu 65536 qdisc noop qlen 1000
2: tunl0@NONE: <NOARP> mtu 1480 qdisc noop qlen 1000
...                                                    (every interface DOWN: "noop")

fork test:
/bin/sh: line 0: can't fork: Resource temporarily unavailable
cgroup pids.events: max 1
```

**What to notice:**

- **Our `child` is PID 1**, and `ps` sees only three processes, because the
  freshly mounted `/proc` belongs to the new PID namespace.
- **The new network namespace has only a loopback, and it's down.** (The
  `tunl0`, `gre0`... devices appear automatically in every network namespace
  when those tunnel modules are loaded on the host, also down.) The container
  has no network until something creates a veth pair for it: Docker's job,
  [built by hand in the TCP/IP guide](../networking/tcp-ip/real-life-guide-v1.md#chapter-52-host-networking-internals-namespaces-veth-bridges-conntrack).
- **The fork bomb stopped at 20 processes**, and `pids.events` recorded it. That
  one line is why a container can't take down its host by forking.
- **`UseCgroupFD`/`CgroupFD`** start the child directly inside the cgroup,
  with no window in which it runs unlimited.

**Exercises:**

1. Inside the container run `hostname foo`, then check the host's hostname. The
   UTS namespace kept them apart.
2. Add `CLONE_NEWUSER` with a UID mapping so it runs without root: a
   *rootless* container.
3. Add networking: create a veth pair, move one end into the child's network
   namespace (`ip link set veth1 netns <pid>`), and give both ends addresses.
4. Copy Chapter 76's `memlimit` binary into the rootfs and run it with
   `-live 100`. The 64 MB `memory.max` kills it. Then read
   `memory.events` in the cgroup directory to see the `oom_kill` count.

### Common mistakes
- Thinking `chroot` alone is isolation. Without namespaces and cgroups it's a
  different view of the same system, and a root process can escape it.
- Forgetting to mount a fresh `/proc`, so `ps` shows the host's processes.
- Running untrusted code in a home-made runtime. This is for learning. Use
  runc, crun, or gVisor for real isolation.

### Check yourself
1. Which namespace makes the child PID 1, and which makes its hostname private?
2. Why is a new `/proc` mount needed inside the container?
3. What stops the fork bomb, and where would you see that it happened?

### Further reading
- Liz Rice, "Containers From Scratch" (talk and code), the classic version of this lab.
- *Container Security* (Liz Rice, O'Reilly), chapters on namespaces, cgroups,
  and what real runtimes add.
- Chapters 44–46 of this guide.

## Chapter 81 -- Profiling Go on Linux: pprof, the execution tracer, and perf

### In one sentence
Go has profiling built in: one import exposes CPU, heap, goroutine, block, and
mutex profiles from a running service, so you can find a hot spot or a leak in
minutes instead of guessing.

### Why this matters
Chapter 70 described the ladder from symptoms to profiles. For Go services,
the top of that ladder is unusually cheap: a CPU profile costs a few percent
while running and nothing when it isn't, and a goroutine dump shows every
goroutine's stack. Most Go performance incidents (an expensive call in a hot
path, a goroutine leak, lock contention) are found this way.

### How it actually works

| Profile | Answers | Endpoint |
|---|---|---|
| CPU | where is CPU time going? | `/debug/pprof/profile?seconds=30` |
| heap | what is holding memory (and what allocates most)? | `/debug/pprof/heap` |
| goroutine | how many goroutines, and where are they stuck? | `/debug/pprof/goroutine?debug=1` |
| block / mutex | where do goroutines wait on channels or locks? | `/debug/pprof/block`, `/mutex` (enable with `runtime.SetBlockProfileRate` / `SetMutexProfileFraction`) |
| trace | a timeline of scheduling, GC, syscalls | `/debug/pprof/trace?seconds=5` → `go tool trace` |

`perf` (Chapter 70) works on Go binaries too. Go keeps frame pointers on amd64
and arm64, so `perf record -g` gets proper call stacks, and it can show time
spent in the kernel, which pprof can't.

### Build it in Go (30 min) -- find two planted bugs

```go
// profiled: a service with two planted bugs -- a CPU hot spot and a goroutine
// leak -- and the profiling endpoints that find them.
//
//	go run ./profiled -load 8 &          # serve, and hammer ourselves with 8 clients
//	go tool pprof -top -seconds 5 http://localhost:6060/debug/pprof/profile
//	curl -s 'localhost:6060/debug/pprof/goroutine?debug=1' | head -20
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof/* on the default mux
	"regexp"
	"runtime"
	"time"
)

func main() {
	load := flag.Int("load", 0, "also run N client goroutines against ourselves")
	flag.Parse()

	// BUG 1: the regexp is compiled on EVERY request. Correct code compiles it
	// once, at package level: var emailRE = regexp.MustCompile(...)
	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		fmt.Fprintln(w, re.MatchString(r.URL.Query().Get("email")))
	})

	// BUG 2: a goroutine waits for a reply that never comes. Each request leaks
	// one goroutine (and everything it references), forever.
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		reply := make(chan string) // unbuffered, and nobody ever sends
		go func() {
			select {
			case msg := <-reply:
				log.Println(msg)
				// FIX: add `case <-time.After(5 * time.Second): return` or a context
			}
		}()
		fmt.Fprintln(w, "queued")
	})

	go func() { // a cheap "leak alarm": watch the goroutine count
		for range time.Tick(5 * time.Second) {
			log.Printf("goroutines: %d", runtime.NumGoroutine())
		}
	}()
	// Keep-alive pool big enough for every client goroutine. With Go's default
	// of 2 idle connections per host, the load generator itself would spend
	// its time opening and closing connections -- and dominate the profile.
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 64}}
	for i := 0; i < *load; i++ {
		go func() {
			time.Sleep(200 * time.Millisecond) // let the server start
			for n := 0; ; n++ {
				path := "/validate?email=user@example.com"
				if n%10 == 0 {
					path = "/notify"
				}
				if resp, err := client.Get("http://localhost:6060" + path); err == nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			}
		}()
	}
	log.Println("listening on :6060 -- profiles at /debug/pprof/")
	log.Fatal(http.ListenAndServe("localhost:6060", nil)) // NEVER expose pprof publicly
}
```

```text
$ go run ./profiled -load 8 &
$ go tool pprof -top -cum -seconds 5 http://localhost:6060/debug/pprof/profile
Duration: 5s, Total samples = 39.26s (784.98%)
     0.01s 0.025% 44.75%      4.47s 11.39%  main.main.func4          <- the load generator
         0     0% 53.36%      3.05s  7.77%  main.main.func1          <- the /validate handler
         0     0% 53.46%      2.62s  6.67%  regexp.MustCompile       <- 86% of the handler!
     0.01s 0.025% 53.49%      2.62s  6.67%  regexp.compile

$ curl -s 'localhost:6060/debug/pprof/goroutine?debug=1' | head -3
goroutine profile: total 493
489 @ 0x100ce1a7c 0x100c76a10 0x100c76584 0x100e94b88 0x100ce9094
#	0x100e94b87	main.main.func2.1+0x27	.../profiled/main.go:34      <- the leak, with its line number
```

**What to notice:**

- **Read the cumulative (`-cum`) column for your own functions.** The
  handler (`func1`) spent 3.05 s, and 2.62 s of that, **86%**, was compiling a
  regexp that should have been compiled once at startup. The rest of the
  profile is loopback networking, because the load generator runs in the same
  process.
- **The goroutine profile groups identical stacks**: "489 @ ..." means 489
  goroutines are stuck at the same place, `main.go:34`, the `select` waiting on
  a channel nobody sends to. The log line "goroutines: 243 → 491" in five seconds
  is the alarm. The profile is the diagnosis.
- **An earlier version of this lab generated load with `curl` in a shell
  loop**, and the CPU profile showed almost nothing: the bottleneck was starting
  `curl` processes, not the service. Then the in-process load generator used
  Go's default `MaxIdleConnsPerHost` of 2, and connection churn dominated the
  profile. Profiles describe the system *as it's loaded*. Make sure the load is
  realistic.
- **Never expose `/debug/pprof` publicly.** Profiles leak internals and are
  expensive. Bind to localhost or an admin port, as here.

**Exercises:**

1. Fix both bugs (compile the regexp once; add a timeout case to the
   `select`), rerun, and compare profiles: `go tool pprof -diff_base old.pb.gz new.pb.gz`.
2. Take a heap profile under load and run `top -inuse_space`. Where does the
   leaked goroutines' memory show up?
3. On Linux, run `perf record -g -p $(pgrep profiled) -- sleep 5` and
   `perf report`. How much time is in the kernel's TCP stack?

### Common mistakes
- Optimising from intuition instead of a profile.
- Profiling with unrealistic load (or none at all).
- Reading `flat` time only. Your function's total cost is in `cum`.
- Leaving pprof reachable from the internet.

### Check yourself
1. What's the difference between `flat` and `cum`?
2. Which profile finds a goroutine leak, and what does "489 @" mean?
3. What can `perf` show that pprof can't?

### Further reading
- [Go guide, Chapter 37 (pprof, trace, benchmarks)](../Golang/real-life-golang-guide.md#37-performance-tuning-pprof-trace-and-benchmark-methodology) and
  [Chapter 58 (debugging Go in production)](../Golang/real-life-golang-guide.md#58-debugging-go-in-production-delve-pprof-traces-stack-dumps-godebug).
- The Go blog: "Profiling Go Programs"; go.dev/doc/diagnostics.
- Go plan [Day 111: memory and goroutine leak hunting](../Golang/detailed-90-day-plan/week15.md#day-111-memory-goroutine-leak-hunting).

## Chapter 82 -- Capstone: shipping a well-behaved Go service (image, systemd, Kubernetes)

### In one sentence
A well-behaved Linux service is small, unprivileged, resource-aware,
observable, and stops gracefully. This chapter packages the HTTPS guide's Go
server three ways, with every setting traced back to a chapter that explains
why.

### Why this matters
Most production incidents with containerised Go services come down to
packaging and runtime settings, not the code: running as root, no memory limit
hint for the GC, the wrong GOMAXPROCS, PID 1 with no signal handling, a
`preStop` hook that calls a binary the image doesn't have, or a 1 GB image full
of tools an attacker can use. Every one of those is a line in the files below.

### The container image

```dockerfile
# syntax=docker/dockerfile:1
# ---- build stage: full Go toolchain, never shipped ----
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
# CGO_ENABLED=0: a static binary that needs no libc in the final image.
# -trimpath and -ldflags="-s -w": reproducible paths, smaller binary.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./https/httpsserver

# ---- final stage: no shell, no package manager, non-root ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
USER nonroot:nonroot
EXPOSE 8443
ENTRYPOINT ["/server", "-addr", ":8443", "-redirect", ""]
```

Built and run while writing this guide:

```text
$ docker build -f Dockerfile -t gosvc:lab .
$ docker images gosvc:lab --format '{{.Size}}'
14.8MB
$ docker run -d --rm -p 18443:8443 --read-only --cap-drop=ALL gosvc:lab
$ curl -sk https://127.0.0.1:18443/readyz
ready
$ docker exec <id> sh
exec: "sh": executable file not found in $PATH         <- nothing to exploit
```

| Line | Why | Chapter |
|---|---|---|
| multi-stage build | the toolchain never ships | [48](#chapter-48-writing-a-good-dockerfile) |
| `CGO_ENABLED=0` | static binary, no libc, runs on `distroless/static` | 48 |
| `:nonroot` + `USER` | no root inside the container | [36](#chapter-36-the-permission-model-beyond-chmod-777), [71](#chapter-71-linux-security-hardening-selinux-apparmor-auditd-capabilities-and-least-privilege) |
| `--read-only`, `--cap-drop=ALL` | the process can't modify its image or gain privileges | 71 |
| no shell in the image | `docker exec ... sh` fails, and so does an attacker's | 52 |

### The systemd unit (for VMs and bare metal)

```ini
# /etc/systemd/system/gosvc.service
[Unit]
Description=Go HTTPS service (gosvc)
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
# %d = the credentials directory: LoadCredential copies root-only key files there,
# readable by this service alone, even though it runs as an unprivileged dynamic user.
ExecStart=/opt/gosvc/server -addr :443 -redirect :80 -cert %d/tls.crt -key %d/tls.key
Restart=on-failure
RestartSec=2s
# Graceful shutdown: SIGTERM, then up to 30s to drain before SIGKILL.
KillSignal=SIGTERM
TimeoutStopSec=30s

# Identity: a throwaway user created at start, no shell, no home.
DynamicUser=yes
# Bind 80/443 without root: grant exactly one capability.
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
NoNewPrivileges=yes

# Resource limits (cgroup v2) -- and tell the Go runtime about them.
MemoryMax=512M
CPUQuota=200%
Environment=GOMEMLIMIT=460MiB
LimitNOFILE=65536

# Sandboxing: read-only OS, private /tmp, no access to /home, few syscalls.
ProtectSystem=strict
ProtectHome=yes
PrivateTmp=yes
PrivateDevices=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
SystemCallFilter=@system-service
SystemCallArchitectures=native
LoadCredential=tls.crt:/etc/gosvc/tls.crt
LoadCredential=tls.key:/etc/gosvc/tls.key

[Install]
WantedBy=multi-user.target
```

`systemd-analyze verify` passes this unit, and `systemd-analyze security`
rates its exposure **3.1 ("OK")**. A default unit with no sandboxing scores
about 9.6 ("UNSAFE").

| Setting | Why | Chapter |
|---|---|---|
| `DynamicUser=yes` + `AmbientCapabilities=CAP_NET_BIND_SERVICE` | ports 80/443 without root | [36](#chapter-36-the-permission-model-beyond-chmod-777), [71](#chapter-71-linux-security-hardening-selinux-apparmor-auditd-capabilities-and-least-privilege) |
| `LoadCredential` + `%d` | root-only key files, readable by this service alone | 67 |
| `MemoryMax` + `GOMEMLIMIT` (90%) | kernel limit, and the GC told about it | [76](#chapter-76-memory-limits-the-go-heap-gomemlimit-and-the-oom-killer) |
| `CPUQuota=200%` | Go 1.25+ sets GOMAXPROCS=2 from it | [75](#chapter-75-cpu-limits-gomaxprocs-cgroups-and-throttling) |
| `LimitNOFILE=65536` | room for connections | [78](#chapter-78-file-descriptors-and-the-netpoller-holding-thousands-of-connections) |
| `KillSignal=SIGTERM`, `TimeoutStopSec=30s` | drain before SIGKILL | [74](#chapter-74-processes-from-go-exec-exit-codes-signals-and-pid-1) |
| `ProtectSystem=strict`, `SystemCallFilter=@system-service` | read-only OS, restricted syscalls | [67](#chapter-67-systemd-beyond-start-stop-dependencies-timers-sandboxing-and-debugging) |

### The Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gosvc
spec:
  replicas: 3
  selector:
    matchLabels: {app: gosvc}
  template:
    metadata:
      labels: {app: gosvc}
    spec:
      terminationGracePeriodSeconds: 30     # > preStop sleep + drain time
      securityContext:
        runAsNonRoot: true
        seccompProfile: {type: RuntimeDefault}
      containers:
      - name: server
        image: registry.example.com/gosvc:1.4.2   # pin a version (or digest), never :latest
        ports:
        - containerPort: 8443
        env:
        - name: GOMEMLIMIT                   # ~90% of the 256Mi limit below (Chapter 76)
          value: "230MiB"
        resources:
          requests: {cpu: "500m", memory: "256Mi"}
          limits:   {cpu: "2",    memory: "256Mi"}   # GOMAXPROCS follows the CPU limit (Go 1.25+)
        readinessProbe:                       # "send me traffic?"  fails while draining
          httpGet: {path: /readyz, port: 8443, scheme: HTTPS}
          periodSeconds: 2
          failureThreshold: 2
        livenessProbe:                        # "restart me?"  never depends on other services
          httpGet: {path: /healthz, port: 8443, scheme: HTTPS}
          periodSeconds: 10
          failureThreshold: 3
        lifecycle:
          preStop:                            # give endpoint removal time to propagate
            sleep: {seconds: 5}               # K8s 1.30+; distroless has no /bin/sleep to exec
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities: {drop: ["ALL"]}
```

This manifest was checked against the real Kubernetes API types
(`k8s.io/api/apps/v1`, strict decoding). Writing it turned up two bugs
worth knowing about:

- The first draft's `preStop` ran `/server -sleep 5s`, but distroless images
  have **no `/bin/sleep`** and the server has no such flag, so the hook would
  fail on every pod termination. Kubernetes 1.30+ has a built-in
  `sleep` action for exactly this case.
- The first draft took `GOMEMLIMIT` from `limits.memory` via `resourceFieldRef`,
  which gives **100%** of the limit, not the 90% Chapter 76 recommends.

| Setting | Why | Chapter |
|---|---|---|
| `readinessProbe` on `/readyz` | fails while draining, so no new traffic | [60](#chapter-60-deployment-patterns-rolling-updates-health-checks-rollbacks) |
| `livenessProbe` on `/healthz` | never depends on other services | 60 |
| `preStop: sleep: 5` | lets endpoint removal propagate before SIGTERM | 60, [TCP/IP 57](../networking/tcp-ip/real-life-guide-v1.md#chapter-57-production-load-balancing-l4-vs-l7-health-checks-draining-retries) |
| `terminationGracePeriodSeconds: 30` | longer than preStop + drain | 74 |
| `limits.cpu: "2"` | GOMAXPROCS follows automatically (Go 1.25+) | 75 |
| `GOMEMLIMIT: 230MiB` with `limits.memory: 256Mi` | no surprise OOM kills | 76 |
| `runAsNonRoot`, `readOnlyRootFilesystem`, `drop: [ALL]`, `seccompProfile` | least privilege | 71 |
| pinned image tag | reproducible rollbacks | [52](#chapter-52-registries-and-image-security-basics) |

### Exercises
1. Deploy the manifest to `kind` (Chapter 53), run a load generator against
   the Service, and do a rolling update. Count errors (target: zero).
2. Remove `preStop` and repeat. How many requests fail during the rollout?
3. Lower `limits.memory` to 128Mi without changing `GOMEMLIMIT`. What happens,
   and what does `kubectl describe pod` say?
4. Add `/debug/pprof` on a separate, cluster-internal port and take a CPU
   profile with `kubectl port-forward`.

### Check yourself
1. Name three settings in these files that exist only because of the Go
   runtime, and the chapter that explains each.
2. Why does the `preStop` hook use `sleep:` rather than `exec:`?
3. What does `LoadCredential` solve that file permissions alone can't, for a
   `DynamicUser` service?

### Further reading
- [Go guide, Chapter 57 (shipping: Docker, systemd, graceful shutdown)](../Golang/real-life-golang-guide.md#57-shipping-it-cross-compilation-docker-systemd-graceful-shutdown)
  and [Chapter 67 (Go in containers and Kubernetes)](../Golang/real-life-golang-guide.md#67-running-go-in-containers-and-kubernetes).
- Go plan [Week 15: production hardening](../Golang/detailed-90-day-plan/week15.md).
- [HTTPS guide, Chapter 19](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown), the server being shipped.

# Part 21 -- Where to go next

You have covered, in depth, what an operating system does, how Linux
exposes that to you at the command line, and how containers and
Kubernetes package and run software at scale. You also covered the advanced
operational layer: fleet drift, systemd, storage, networking, observability,
hardening, and incident response. This last part is deliberately short: an
honest accounting of what still sits beyond this guide, and how to keep
learning once you close it.

## Chapter 83 -- Honest gaps and a study plan

### In one sentence
This guide gave you a strong practical foundation across OS concepts, Linux,
containers, Kubernetes, and production operations. That is still not "all of
Linux." Naming the remaining specialist depths precisely is more useful than
pretending any single guide can be complete.

### Why this matters
Knowing what you don't know is a professional skill in itself. A beginner
SRE who says "I understand the fundamentals of process scheduling and
containers, and I know that I have not yet gone deep on distributed
consensus or kernel networking internals" is more hireable and more
trustworthy than one who either overclaims mastery or has no sense of the
shape of the field at all. This chapter is a map of the edges.

### How it actually works
Honest gaps, organized by phase, with a concrete next resource for each:

```
   PHASE 1 (OS concepts) -- what this guide under-covered:

     - Kernel internals beyond the syscall boundary (how the scheduler,
       memory manager, and VFS are actually implemented in C).
       Next: "Linux Kernel Development" by Robert Love; read the
       actual kernel source for a subsystem you care about.

     - Filesystem internals in depth (how ext4's journal actually
       works, B-trees in Btrfs/XFS, copy-on-write).
       Next: "Operating Systems: Three Easy Pieces" (free online,
       ostep.org) -- covers this more rigorously than this guide did.

     - Real-time operating systems and embedded constraints.
       Not covered here at all -- a genuinely different field with
       different tradeoffs (predictability over throughput).

   PHASE 2 (Linux fundamentals) -- what this guide under-covered:

     - Advanced networking beyond host-level diagnosis: BGP, EVPN/VXLAN,
       traffic engineering, load-balancer internals, and routing at scale.
       Ch 69 gave you the Linux packet-path model; running a network team is
       a deeper specialization.
       Next: "TCP/IP Illustrated" by W. Richard Stevens; the
       Linux Advanced Routing & Traffic Control HOWTO (lartc.org).

     - Security hardening at dedicated security-engineering depth:
       SELinux/AppArmor policy authoring, kernel exploit mitigation research,
       audit frameworks, incident forensics, and compliance operations. Ch 71
       gave you the Linux operator's model; security is its own field.
       Next: the [[real-life-security-guide-v1]] companion guide in
       this wiki's security/ folder covers this ground in depth.

     - Distributed systems fundamentals (consensus, CAP theorem,
       replication) -- necessary for understanding *why* etcd and
       Kubernetes's control plane are built the way they are, only
       gestured at in Chapter 54.
       Next: "Designing Data-Intensive Applications" by Martin
       Kleppmann.

   PHASE 3 (Containers/Kubernetes) -- what this guide under-covered:

     - Kubernetes networking internals (CNI plugins, how Services
       actually route packets via iptables/IPVS, service meshes like
       Istio/Linkerd).
       Next: "Kubernetes Networking" chapters in "Kubernetes in
       Action" by Marko Luksa; the CNCF's Istio documentation.

     - Multi-cluster and multi-region operation, GitOps (ArgoCD/Flux),
       and progressive delivery (canary releases, feature flags at
       the infrastructure layer).
       Next: CNCF's own "Cloud Native Landscape" as a map of what
       exists; ArgoCD's official documentation as a hands-on start.

     - Observability platform engineering: Prometheus architecture and PromQL
       at scale, distributed tracing storage, high-cardinality control, log
       pipeline economics, and profiling as a service. Ch 70 gave you the
       troubleshooting ladder; operating the observability platform is a
       specialty.
       Next: "Distributed Systems Observability" by Cindy Sridharan
       (free online).

     - Cost and capacity planning at scale -- right-sizing a fleet of
       nodes, spot/preemptible instance strategy, multi-tenant
       cluster isolation.
       Next: your cloud provider's own well-architected / cost
       optimization documentation (AWS, GCP, and Azure each publish
       one).
```

A study plan, if you want one rather than just a map: pick ONE gap above
that's closest to your actual job's needs, go deep on it for a month using
the named resource, and come back to this guide's capstones (Part 18)
using what you learned as an added constraint (e.g. "redo Project 3, but
also add a NetworkPolicy" once you've studied Kubernetes networking).
Breadth-first learning (a little of everything) is how you got here;
depth-first learning (one thing at a time, thoroughly) is how you become
genuinely senior at any one of them.

### Real-world example
This gaps-and-resources structure mirrors how most senior engineers
actually describe their own learning path when asked in interviews or
mentoring conversations: nobody claims to have learned "all of Linux" or
"all of Kubernetes" from one source. The common pattern, visible across
countless engineering blogs and conference "how I learned X" talks, is a
broad foundational resource (often a book or course much like this guide)
followed by picking up genuine depth in two or three specific areas that
turned out to matter for their actual job -- one engineer goes deep on
container networking because their company runs a service mesh, another
goes deep on filesystem internals because they work on a database engine.
The map above is deliberately structured the same way: broad coverage
here, named depth resources for whichever direction your actual work
pulls you.

### Try it yourself
Pick one gap from the list above that's closest to something your current
or target job actually touches. Spend 30 minutes right now finding and
bookmarking the named resource (or an equivalent you prefer), and write
one sentence about *why* that gap matters for your specific situation.
That sentence is worth more than the bookmark -- it's what will make you
actually go back and use it.

### Common mistakes
- Trying to close every gap listed here before considering yourself
  "ready" -- this guide's whole premise is that breadth-first, practical
  competence is valuable on its own; depth can be added incrementally,
  driven by actual need.
- Choosing a depth topic based on what sounds impressive rather than what
  your actual job or target job needs.
- Reading about a gap area without doing the equivalent of this guide's
  "Try it yourself" sections for it -- reading "TCP/IP Illustrated"
  without ever running `tcpdump` on your own traffic teaches you far less
  than reading half of it and testing every claim.

### Across the series

Several of the gaps above are covered by the companion guides in this wiki:

- **Networking beyond Part 11:** [Networking from Zero (TCP/IP)](../networking/tcp-ip/real-life-guide-v1.md)
  goes from Ethernet frames to BGP, with packet-capture practice, and its
  [Part 13](../networking/tcp-ip/real-life-guide-v1.md#part-13-professional-networking-operations-debugging-and-design) covers host networking internals, firewalls,
  tunnels, and Kubernetes networking at professional depth.
- **What your servers serve:** [The HTTPS Request Lifecycle](../v2-https/real-life-guide-v1.md)
  covers TLS, HTTP, proxies, CDNs, and SLOs, with a Go
  [production HTTPS server](../v2-https/real-life-guide-v1.md#chapter-19-a-production-https-server-in-go-tls-timeouts-headers-and-shutdown).
- **One request, layer by layer:** [The OSI guide](../networking/real-life-example-osi.md), a short read that
  ties the stack together and works well for teaching others.
- **Writing the software:** the [Go guide](../Golang/real-life-golang-guide.md) and its
  [120-day plan](../Golang/golang-90-day-plan.md), whose later weeks build production network
  services that need exactly the Linux knowledge in this guide.
- **Attacking and defending these systems:** [Security from Zero](../security/real-life-guide.md) (SSH hardening,
  OpenSSL, a CA of your own) and [Security Engineering in Depth](../security/real-life-security-guide-v1.md): container
  isolation and escapes ([Chapter 15](../security/real-life-security-guide-v1.md#chapter-15-container-internals-and-isolation)), Kubernetes RBAC, admission,
  and runtime security ([Chapters 17–20](../security/real-life-security-guide-v1.md#chapter-17-kubernetes-rbac-and-identity)), and secrets
  ([Chapter 73](../security/real-life-security-guide-v1.md#chapter-73-secrets-management-and-credential-lifecycle)).

### Check yourself
1. Name one gap from the list above that is directly relevant to your own
   current job or the job you want next, and explain why in one sentence.
2. Why might "go deep on one gap for a month" be more valuable at this
   stage than "read a little about every gap"?
3. This guide mentioned but did not teach service meshes (Istio/Linkerd).
   Based on what you learned about Kubernetes Services in Chapter 56,
   what problem do you think a service mesh exists to solve that a plain
   Service does not?

### Further reading
- "Operating Systems: Three Easy Pieces," Remzi and Andrea
  Arpaci-Dusseau -- free at ostep.org
- Martin Kleppmann, "Designing Data-Intensive Applications" (2017)
- Cindy Sridharan, "Distributed Systems Observability" -- free online
- CNCF Cloud Native Landscape, landscape.cncf.io

## Chapter 84 -- Staying current

### In one sentence
Operating systems, Linux, and Kubernetes are all actively developed,
moving targets -- the specific commands and defaults in this guide are
accurate as of when it was written, and staying useful over the following
years means knowing which sources to check periodically, not re-reading
this guide.

### Why this matters
A guide like this one is a snapshot. Kubernetes and the Linux kernel both ship
regularly, and they periodically deprecate APIs, graduate features, or change
defaults (the cgroups v1-to-v2 transition is a good example of a change that
affected how memory and CPU limits were enforced in practice).
An engineer who treats what they learned once as permanently true will,
within a couple of years, be confidently wrong about specific details even
while being right about the underlying concepts -- and the underlying
concepts (processes, memory, namespaces, the control plane/node split)
are exactly what stays stable across those version changes, which is why
this guide focused on them.

### How it actually works
A short, durable list of places to check periodically, organized by how
often they change:

```
   RARELY CHANGES (concepts -- revisit if you feel rusty, not on a
   schedule):
     - The core OS concepts in Parts 1-6 of this guide. Processes,
       memory, syscalls, and the boot sequence have been stable in
       their fundamentals for decades and will likely remain so.

   CHANGES OVER YEARS (check every year or two):
     - Linux kernel release notes (kernelnewbies.org/LinuxChanges) --
       skimming major version summaries tells you what's newly
       possible or newly deprecated.
     - Your distribution's release notes (Ubuntu LTS, Debian stable,
       RHEL) if you operate a specific distribution professionally.

   CHANGES OVER MONTHS (check a few times a year, or when planning
   an upgrade):
     - Kubernetes release notes and the deprecated API migration
       guide (kubernetes.io/docs/reference/using-api/deprecation-guide)
       -- checking this before cluster upgrades is a professional habit;
       skipping it is a common way to break production.
     - Docker Engine release notes, if you build images as part of
       your job.
     - The CNCF's annual "Cloud Native Landscape" survey, for a sense
       of which tools in the ecosystem are gaining or losing adoption.

   CONTINUOUS (follow loosely, don't obsess):
     - Your own company's or a major cloud provider's incident
       postmortems and engineering blog -- reading how other
       real teams break things and fix them is, chapter for chapter,
       some of the most useful continuing education available, and
       it's exactly the kind of material this guide's "Real-world
       example" sections were built from.
     - CVE feeds relevant to your stack (e.g. a mailing list or RSS
       feed for the specific Linux distribution, container runtime,
       and Kubernetes distribution you actually run).
```

### Real-world example
The cgroups v1-to-v2 migration, referenced throughout Part 6 and Part 14
of this guide, is a concrete illustration of exactly the kind of change
that rewards staying current: for years, `docker stats` and Kubernetes
resource enforcement relied on cgroups v1's separate hierarchies per
resource type; cgroups v2's unified hierarchy (merged into the mainline
kernel starting around 4.5, but not the default in major distributions
and container runtimes until years later, roughly 2021-2022 for most
mainstream Linux distributions and container tooling) changed some
memory-accounting edge cases enough that teams upgrading their base OS
image sometimes saw containers get OOM-killed at different thresholds
than before, purely because of the underlying cgroup version switching
under them. Nothing about the *concept* of cgroups (Chapter 44) changed;
the concrete mechanics teams needed to re-verify after an OS upgrade did.
This is the recurring pattern behind essentially every "why did upgrading
the base image change our container's behavior" incident: the concepts
this guide taught you remain the right mental model; the specific
mechanics are worth re-checking against current documentation whenever
you cross a major version boundary.

### Try it yourself
Subscribe to, or bookmark for periodic checking, one source from each of
the "changes over months" and "continuous" categories above that's
relevant to your actual stack. If you don't currently run Kubernetes
professionally, substitute the release notes of whatever you do run
(a specific Linux distribution, a specific database, a specific cloud
provider's compute service) -- the habit matters more than the specific
source.

### Common mistakes
- Treating a guide, book, or course as a one-time investment rather than
  a foundation to be periodically checked against current sources --
  especially dangerous for fast-moving layers like Kubernetes.
- The opposite mistake: chasing every new tool and blog post without a
  stable foundation underneath, which produces someone who knows many
  buzzwords and few fundamentals. The concepts in Parts 1-19 of this
  guide are the antidote to this; keep them solid and let the
  fast-moving specifics rotate on top of them.
- Skipping a Kubernetes deprecated-API check before a version upgrade,
  then discovering in production that a manifest using a removed
  `apiVersion` no longer applies.
- Ignoring your own organization's postmortems as a learning source
  because they feel like "just work," when they are, chapter for
  chapter, as valuable as any external resource on this list.

### Check yourself
1. Why does this guide claim the core OS concepts (Parts 1-6) will stay
   useful longer than the specific Kubernetes commands in Part 17?
2. Name one category of source from the list above you would check
   before performing a Kubernetes version upgrade, and explain what bad
   outcome checking it helps you avoid.
3. What changed, mechanically, in the cgroups v1-to-v2 transition, and
   why did teams who upgraded their base OS image sometimes see
   different OOM-killing behavior as a result?

### Further reading
- kernelnewbies.org, "Linux Changes" -- kernel version summaries in
  plain language
- Kubernetes documentation, "Deprecated API Migration Guide"
  (kubernetes.io)
- CNCF, annual "Cloud Native Landscape" and "State of Cloud Native
  Development" reports

---

**That's the full arc** -- from "what is an OS actually for" in Part 1 to
debugging and improving production Linux systems in Part 19, by way of
processes, memory, I/O, concurrency, boot, the filesystem, the command line,
users and permissions, networking, scripting, performance, containers, and
Kubernetes. What follows is reference material, not more reading: a glossary,
a command cheat sheet, a troubleshooting playbook, and further reading -- check
back against these while you're actually working on a system, rather than
reading them front to back.

---
## Continue the series

**Next, step 2a (Networking: the map): [The OSI Model, One Click at a Time](../networking/real-life-example-osi.md).** It covers one click followed through all seven layers; the map for everything after it. Hands-on: concept map, 1 hour.

---

# Appendix A -- Glossary (plain language)

Terms are defined the way they're actually used in conversation, not the
way a spec defines them. Where a chapter covers a term in depth, it's
named so you can go back.

- **ACL (Access Control List)** -- an extra, finer-grained permission list
  on a file, beyond owner/group/other. Ch36.
- **Bind mount** -- a Docker mount that points at a specific path on the
  host filesystem, as opposed to a Docker-managed named volume. Ch49.
- **Boot loader** -- the small program (e.g. GRUB) that loads the kernel
  into memory and hands it control. Ch21.
- **Capability (Linux)** -- a narrow privilege carved out of root, such
  as binding low ports, changing ownership, or loading kernel modules.
  Ch71.
- **cgroup (control group)** -- a kernel mechanism that limits and
  accounts for a group of processes' use of CPU, memory, and I/O.
  The basis of container resource limits. Ch43, Ch44.
- **ConfigMap** -- a Kubernetes object holding non-secret configuration
  data, injected into Pods as environment variables or files. Ch57.
- **Container** -- a process (or group of processes) running with its
  own PID/network/mount namespaces and a cgroup limit, giving the
  illusion of an isolated machine while sharing the host kernel. Ch44.
- **Context switch** -- the CPU saving one process/thread's state and
  loading another's, so it can run something else. Ch9.
- **Control plane** -- the set of Kubernetes components (API server,
  etcd, scheduler, controller manager) that decide what *should* be
  running, as opposed to the nodes that actually run it. Ch54.
- **CVE** -- Common Vulnerabilities and Exposures, the standard public
  identifier for a specific named security vulnerability.
- **Daemon** -- a background process with no controlling terminal,
  typically started at boot and running for the system's lifetime
  (e.g. `sshd`, `dockerd`). Ch22.
- **Deployment** -- a Kubernetes object describing a desired Pod template
  and replica count; manages ReplicaSets, which manage Pods. Ch55.
- **Drift** -- the difference between intended configuration and actual
  runtime state across one host or a fleet. Ch66.
- **eBPF** -- a Linux kernel technology for safely running small programs
  attached to kernel/application events, widely used for observability,
  networking, and security tooling. Ch70.
- **etcd** -- the distributed key-value store that is Kubernetes'
  single source of truth for cluster state. Ch54.
- **exec()** -- the syscall family that replaces a process's running
  program with a different one, keeping the same PID. Ch6.
- **fork()** -- the syscall that creates a new process by duplicating
  the calling process. Ch6.
- **FHS (Filesystem Hierarchy Standard)** -- the convention that defines
  what `/etc`, `/var`, `/usr`, `/tmp`, etc. are each for. Ch24.
- **HPA (Horizontal Pod Autoscaler)** -- a Kubernetes controller that
  changes a Deployment's replica count based on observed metrics like
  CPU usage. Ch58.
- **Image (container)** -- a read-only, layered filesystem snapshot plus
  metadata that a container is started from. Ch46.
- **Inode** -- the on-disk data structure holding a file's metadata
  (permissions, size, pointers to data blocks) and identity, separate
  from its name. Ch15.
- **IPC (Inter-Process Communication)** -- any mechanism (pipes, sockets,
  shared memory, message queues) letting separate processes exchange
  data. Ch19.
- **Journaling (filesystem)** -- writing a log of intended changes before
  applying them, so a crash mid-write can be recovered from without a
  full filesystem scan. Ch15.
- **kubelet** -- the agent running on every Kubernetes node that talks to
  the container runtime and reports Pod status back to the control
  plane. Ch54.
- **Layer (image)** -- one immutable filesystem diff stacked inside a
  container image; layers are cached and shared across images. Ch46.
- **Liveness probe** -- a Kubernetes health check whose failure causes
  the container to be restarted. Ch60.
- **Load average** -- a rolling measure of the number of processes
  running or waiting for a resource (CPU or, on Linux, also I/O),
  averaged over 1/5/15 minutes. Ch41.
- **LVM (Logical Volume Manager)** -- a Linux storage layer that lets you
  group physical storage and carve it into logical volumes that can often
  be resized more flexibly than raw partitions. Ch68.
- **MAC (Mandatory Access Control)** -- policy-based access control, such
  as SELinux or AppArmor, enforced beyond ordinary Unix mode bits. Ch71.
- **Namespace (kernel)** -- a kernel mechanism that gives a process its
  own view of a global resource (PIDs, network interfaces, mounts,
  hostname) separate from the rest of the system. Ch44. (Not to be
  confused with a *Kubernetes* namespace, which is a logical grouping of
  objects inside one cluster.)
- **OOM killer** -- the kernel's Out-Of-Memory killer, which forcibly
  kills a process to free memory when the system is critically low.
  Ch13.
- **Page cache** -- RAM the kernel uses to cache recently read/written
  disk data, transparently, so re-reads don't touch the disk. Ch16.
- **Page fault** -- a CPU trap raised when a virtual memory access can't
  be immediately satisfied; "minor" is resolved in RAM, "major" requires
  disk I/O. Ch11.
- **PID 1** -- the first process started by the kernel after boot
  (traditionally `init`, now usually `systemd`); becomes the ultimate
  parent/reaper of every other process. Ch21.
- **Pod** -- the smallest deployable unit in Kubernetes: one or more
  containers that share a network namespace and are always scheduled
  together. Ch55.
- **Readiness probe** -- a Kubernetes health check whose failure removes
  a Pod from a Service's traffic without restarting it. Ch60.
- **ReplicaSet** -- a Kubernetes controller ensuring a specified number
  of identical Pod replicas are running; usually managed for you by a
  Deployment. Ch55.
- **Rolling update** -- replacing old Pods with new ones a few at a time,
  gated by readiness, so the workload stays available throughout. Ch60.
- **Secret (Kubernetes)** -- like a ConfigMap, but intended for sensitive
  values; base64-encoded, not encrypted, by default. Ch57.
- **Service (Kubernetes)** -- a stable virtual IP and DNS name that load
  balances traffic across a dynamic set of Pods, chosen by label
  selector. Ch56.
- **Signal** -- an asynchronous notification sent to a process by the
  kernel or another process (e.g. SIGTERM, SIGKILL, SIGHUP). Ch20.
- **SUID (Set User ID) bit** -- a permission bit that makes a program run
  with its *owner's* privileges rather than the invoking user's. Ch25.
- **Syscall (system call)** -- the only sanctioned way a user-space
  program asks the kernel to do something privileged. Ch2.
- **systemd unit** -- a declarative object managed by systemd: service,
  timer, socket, mount, target, and more. Ch33, Ch67.
- **TLB (Translation Lookaside Buffer)** -- a small, fast, on-CPU cache
  of recent virtual-to-physical address translations. Ch11.
- **Virtual memory** -- the abstraction giving every process its own
  private, contiguous-looking address space, actually backed by
  scattered physical RAM and disk via the page table. Ch10.
- **Volume (Docker/Kubernetes)** -- storage that persists independently
  of a container's lifecycle. Ch49.
- **Zombie process** -- a process that has exited but whose exit status
  hasn't yet been collected (`wait()`ed on) by its parent. Ch5.

# Appendix B -- Command cheat sheet

Grouped by task, not alphabetically -- so you can find "the command for
what I'm trying to do" rather than needing to already know its name.

```
PROCESSES
  ps aux                          list all processes
  ps -ef --forest                 list with a parent/child tree
  top / htop                      live process viewer
  kill -TERM <pid>                ask a process to exit gracefully
  kill -KILL <pid>  (or kill -9)  force-kill, unconditionally
  pkill -f <pattern>              kill by matching command line
  nohup ./long_task.sh &          survive terminal disconnect
  jobs / fg / bg                  manage foreground/background jobs
  strace -p <pid>                 trace a running process's syscalls
  lsof -p <pid>                   list a process's open files/sockets

FILES AND TEXT
  ls -la                          list files, including permissions
  cp -r src/ dst/                 copy a directory recursively
  find . -name "*.log" -mtime -1  files named *.log, modified <1 day ago
  grep -rn "TODO" .               recursive, line-numbered search
  sed 's/foo/bar/g' file          replace foo with bar, all occurrences
  awk '{print $1}' file           print the first whitespace-delimited field
  tail -f /var/log/app.log        follow a growing log file live
  chmod 640 file                  owner rw, group r, others none
  chown user:group file           change ownership

PACKAGES AND SERVICES
  sudo apt install <pkg> / sudo dnf install <pkg>    install a package
  sudo apt remove <pkg> / sudo apt purge <pkg>       remove (purge also drops config)
  systemctl status <svc>          is it running, and why/why not
  systemctl restart <svc>         restart a service now
  systemctl enable <svc>          make it start automatically at boot
  journalctl -u <svc> -f          follow a service's live logs

NETWORKING
  ip addr                         list network interfaces and IPs
  ss -tlnp                        list listening TCP sockets and owners
  curl -v https://host/path       make an HTTP request, verbosely
  dig +short example.com          resolve a hostname, short output
  tcpdump -i eth0 port 443        capture packets on a port

PERFORMANCE
  uptime                          load average (1/5/15 min)
  vmstat 1                        CPU/memory/swap snapshot every second
  iostat -x 1                     per-disk I/O stats every second
  free -h                         memory and swap usage
  df -h / du -sh *                disk space by filesystem / by directory
  ulimit -n                       current open-file limit

DOCKER
  docker build -t name:tag .      build an image from a Dockerfile
  docker run -d --name x image    run a container, detached
  docker ps -a                    list containers, including stopped
  docker logs -f <name>           follow a container's logs
  docker exec -it <name> sh       open a shell inside a running container
  docker compose up -d            start every service in compose.yaml
  docker system df                show Docker disk usage before cleanup
  docker system prune             remove unused Docker data after review

KUBERNETES
  kubectl get pods -o wide        list pods with node/IP info
  kubectl describe pod <name>     full detail, including Events
  kubectl logs <name> --previous  logs from the last crashed instance
  kubectl apply -f file.yaml      create/update from a manifest
  kubectl rollout status deploy/x watch a rollout to completion
  kubectl rollout undo deploy/x   roll back to the previous revision
  kubectl exec -it <pod> -- sh    shell inside a running pod
  helm install/upgrade/rollback   install, update, or revert a chart

ADVANCED OPERATIONS
  uname -a                        kernel/version fingerprint
  cat /etc/os-release             distro fingerprint
  systemctl cat <svc>             effective unit and overrides
  systemd-analyze critical-chain  boot/service dependency timing
  lsblk -f                        block devices and filesystems
  findmnt                         mounted filesystems and options
  journalctl --disk-usage         journald storage use
  ip route / ip rule              routing tables and policy routing
  nft list ruleset                nftables firewall rules
  ip netns list                   network namespaces
  perf stat <cmd>                 hardware/software performance counters
  getcap -r /usr/bin /usr/sbin    files with Linux capabilities
  aa-status / getenforce          AppArmor / SELinux mode, if installed
```

# Appendix C -- Troubleshooting playbook ("the server is slow")

A flowchart-style sequence for the single most common vague complaint in
this field. Work top to bottom; each step either finds the cause or
rules out a whole category.

```
"The server / service is slow." Start here.

1. IS IT ACTUALLY THE SERVER, OR THE NETWORK / CLIENT?
   curl -w "@curl-format.txt" the endpoint yourself, from somewhere
   close to the server. If IT is fast but users report slow, suspect
   DNS, a CDN, or a client-side/network path issue instead. (Ch37, Ch38)

2. CHECK THE FOUR BASELINE RESOURCES, IN THIS ORDER:

   a) CPU        uptime / vmstat 1
                  Load average far above core count + high 'us'/'sy'
                  time -> CPU-bound. (Ch41)

   b) MEMORY      free -h ; check "available", not "free"
                  Available near zero, swap climbing, or dmesg showing
                  "Out of memory: Killed process" -> memory pressure.
                  (Ch13, Ch16)

   c) DISK I/O     iostat -x 1 ; check %util and await
                   %util near 100 with high await -> disk-bound.
                   df -h AND df -i -- full data blocks or full inodes
                   both cause writes to fail. (Ch14, Ch15, Ch16)

   d) PROCESSES WAITING ON I/O   vmstat 1, the 'b' column
                   Consistently non-zero -> processes blocked on I/O,
                   not the CPU scheduler. (Ch41)

3. IF ALL FOUR LOOK FINE, SUSPECT THE APPLICATION ITSELF:
   strace -c -p <pid>            which syscalls dominate?
   lsof -p <pid>                  how many files/sockets open --
                                    approaching the ulimit? (Ch43)
   Check application-level logs and metrics (slow queries, GC pauses,
   thread pool exhaustion) -- this is now outside OS-level diagnosis.

4. IF IT'S A CONTAINER / POD SPECIFICALLY:
   docker stats <name>  or  kubectl top pod <name>
     hitting its memory/CPU LIMIT? -> OOMKilled or CPU-throttled.
     (Ch43, Ch58)
   kubectl describe pod <name> -- any restarts, and what reason?
   kubectl logs <name> --previous -- if it has restarted at all.
   (Ch59)

5. IF IT JUST STARTED BEING SLOW AFTER A CHANGE:
   What deployed, and when? (kubectl rollout history / your CI log)
   Mitigate first: roll back if a known-good previous version exists,
   BEFORE continuing to root-cause live. (Ch65)

6. ONCE MITIGATED: write down what you checked, in what order, and
   what you found at each step -- whether or not it turned out to be
   the cause. This is the raw material for the postmortem (Ch65) and
   for improving this exact playbook next time.
```

# Appendix D -- Further reading, and how this guide's facts were verified

Every "Real-world example" in this guide draws on one of three kinds of
source, chosen deliberately over inventing a plausible-sounding but
made-up scenario:

```
  1. NAMED, PUBLIC INCIDENTS with a primary source -- a postmortem,
     a CVE record, a court filing, or reporting from a named
     publication at the time (e.g. Knight Capital's SEC filing, the
     Mars Pathfinder JPL technical report, CVE-2016-5195 "Dirty COW,"
     CVE-2021-3156 "Baron Samedit," the Therac-25 case studies by
     Leveson and Turner).

  2. OFFICIAL PROJECT DOCUMENTATION AND SOURCE HISTORY -- the
     Kubernetes docs, Docker's own engineering blog, the Linux kernel's
     own changelogs and LKML archives, systemd's design documents.

  3. WIDELY CORROBORATED INDUSTRY PATTERNS -- practices and failure
     modes (like the "healthy-looking crash" readiness-probe problem,
     or the naive-line-by-line-I/O performance trap) that are described
     independently, in similar terms, by multiple unrelated
     engineering organizations' public writing (Google's SRE book,
     Netflix's tech blog, Cloudflare's, Monzo's, and others) -- treated
     as reliable precisely because they were arrived at independently.

  Where a fact's precision matters (a version number, a date, a specific
  percentage), this guide favors the vaguer-but-correct phrasing
  ("around 2021-2022," "roughly," "a widely cited figure") over a
  precise-sounding number it could not verify. If you plan to cite any
  of these examples yourself -- in a blog post, a talk, or an interview
  answer -- look up the named primary source above and read it directly;
  secondhand retellings, including this guide's, drift from the original
  facts over time.
```

Broad, durable resources worth owning, beyond the per-chapter "Further
reading" lists:

- "Operating Systems: Three Easy Pieces," Remzi H. Arpaci-Dusseau and
  Andrea C. Arpaci-Dusseau -- free at ostep.org
- W. Richard Stevens and Stephen A. Rago, "Advanced Programming in the
  UNIX Environment"
- Brendan Gregg, "Systems Performance: Enterprise and the Cloud," 2nd
  edition
- Google, "Site Reliability Engineering" and "The Site Reliability
  Workbook" -- both free at sre.google/books
- The official Kubernetes documentation at kubernetes.io/docs -- kept
  current, unlike any book, by design
- Julia Evans' zines and blog (jvns.ca) -- exceptionally clear,
  practically-minded writing on Linux internals and debugging

# Appendix E -- Answers to "Check yourself"

Answers are deliberately brief -- a correct direction and the key fact,
not a full re-explanation. If an answer doesn't make sense, that's a
signal to reread the chapter's "How it actually works" section, not a
sign you should have memorized it.

```
CHAPTER 1
1. (a) each program gets its own private memory -- really shared RAM
   partitioned by virtual memory; (b) each program appears to have the
   CPU to itself -- really time-sliced by the scheduler; (c) files/
   devices look uniform -- really very different hardware behind a
   common syscall interface.
2. Without an OS, any program could read/overwrite any other program's
   memory or monopolize the CPU/devices -- there was no enforced
   isolation.
3. Isolation and crash containment (one tab crashing doesn't take down
   the browser or other tabs) -- paid for with higher total memory
   usage, since less is shared between processes.

CHAPTER 2
1. The CPU itself, via a hardware-enforced privilege level (ring 0 vs.
   ring 3 on x86); user-mode code physically cannot execute privileged
   instructions.
2. E.g. read() (read bytes from a fd), write() (write bytes), fork()
   (create a process), open()/close() (manage file descriptors),
   mmap() (map memory).
3. Because they require a trap into kernel mode -- a context switch in
   privilege level, with real CPU overhead -- unlike an ordinary
   function call that stays in user space.
4. It was a race condition in the memory subsystem letting an
   unprivileged user write to files they should only be able to read,
   and it had existed in the kernel for roughly nine years before
   discovery -- dangerous because of both severity (arbitrary write)
   and how long it went unnoticed.

CHAPTER 3
1. Shell forks a child, the child execs the target program, the kernel
   loads and runs it, the program's stdout is inherited from the shell,
   the shell wait()s for the exit code, and the shell prints its next
   prompt.
2. fork() creates a near-identical copy of the calling process; exec()
   replaces the current process's program with a different one,
   keeping the same PID.
3. A zombie is a process that exited but whose exit status hasn't been
   collected yet by wait() -- it exists so the parent has a chance to
   read the exit code before the kernel discards it.
4. If the pipeline's last command succeeds even though an earlier
   command failed (without `pipefail`), the shell's exit code reflects
   only the last command -- masking the real failure.

CHAPTER 4
1. E.g. PID, memory mappings, open file descriptors, signal handlers/
   pending signals, scheduling priority, current working directory.
2. Creating a process means the kernel allocating and initializing a
   full set of these structures (and, without copy-on-write
   optimizations, copying memory) -- far more work than a function
   call, which just pushes a stack frame.
3. Handling many concurrent connections without the memory and
   context-switch overhead of one thread (or process) per connection.

CHAPTER 5
1. Running (on CPU), Ready/Runnable (waiting for CPU), Blocked/
   Uninterruptible (D state, waiting on I/O), Stopped (suspended),
   Zombie (exited, awaiting reap).
2. D state (uninterruptible sleep) means the process is waiting on a
   kernel-level I/O operation that ignores signals, including SIGKILL,
   until that I/O completes or times out.
3. Check its process state with `ps` or `top` (is it D-state, blocked
   on I/O?) before assuming it's a CPU or application logic problem.

CHAPTER 6
1. fork() duplicates the calling process; exec() replaces the current
   program image. They're separate because you often want to do
   setup (redirect fds, change environment) in the child AFTER
   forking but BEFORE the new program starts.
2. In the parent, fork() returns the child's PID; in the child, it
   returns 0.
3. Forking a process per connection gives strong isolation (a crashed
   connection handler can't corrupt another's memory) at the cost of
   higher per-connection overhead than threads, which share memory
   but risk one bug corrupting all connections.

CHAPTER 7
1. Threads share the same address space, open file descriptors, and
   signal handlers; each thread keeps its own stack, registers, and
   program counter.
2. No full memory-space duplication is needed -- just a new stack and
   kernel bookkeeping, versus copying/mapping an entire address space
   for a new process.
3. To avoid the complexity and bugs of shared-memory concurrency
   (locks, race conditions) for I/O-bound work, using a single-threaded
   event loop plus non-blocking I/O instead.

CHAPTER 8
1. Deciding which of many runnable processes/threads gets the CPU next,
   and for how long, so the system feels responsive and fair.
2. CFS models an idealized "perfectly fair" CPU that gives every
   runnable task an equal share, and picks the task with the least
   accumulated virtual runtime so far to run next.
3. No -- it makes the process MORE LIKELY to be scheduled sooner and
   get a larger share of CPU time, but the kernel can still delay it
   if higher-priority (lower nice) work is also runnable.

CHAPTER 9
1. An interrupt is a hardware or software signal that stops the CPU's
   current work to run a handler immediately; a context switch is the
   act of saving one task's state and loading another's -- interrupts
   are one common trigger for a context switch.
2. Cache and TLB contents built up for the old task are cold for the
   new one, causing slower memory access after the switch; and the
   scheduler itself consumes CPU cycles deciding what to run next.
3. Because it's driven by I/O readiness events rather than one
   OS thread per connection, so far fewer threads exist to be switched
   between in the first place.

CHAPTER 10
1. Giving each process its own private, contiguous-looking address
   space, so processes can't see or corrupt each other's memory and
   don't need to know about physical RAM layout.
2. The Memory Management Unit (MMU), using the page table (with the
   TLB as its cache) that the kernel maintains per process.
3. ASLR randomizes where a process's memory regions are placed in ITS
   virtual address space on each run -- only possible because that
   address space is a kernel-managed abstraction, not physical
   addresses.

CHAPTER 11
1. A minor fault is resolved by mapping already-resident RAM (fast); a
   major fault requires reading the page in from disk/swap (orders of
   magnitude slower).
2. The TLB caches recent virtual-to-physical address translations, so
   the CPU doesn't have to walk the page table on every memory access.
3. Thrashing is when the system spends most of its time swapping pages
   in and out rather than doing real work; CPU usage can look low
   because processes are mostly blocked waiting on disk I/O, not
   actually computing.

CHAPTER 12
1. CPU registers, CPU cache (L1/L2/L3), RAM, SSD/NVMe, spinning disk /
   network storage -- fastest and smallest to slowest and largest.
2. Roughly five to six orders of magnitude (RAM access ~100ns, SSD
   ~10-100 microseconds, spinning disk several milliseconds).
3. Because the same underlying trade-off -- a small fast layer in front
   of a large slow one, exploiting locality of access -- recurs at
   every scale of computing, from hardware to distributed systems.

CHAPTER 13
1. The kernel allows more memory to be "promised" via allocation calls
   than physically exists, betting that not all processes will use
   their full allocation simultaneously -- it does this to allow
   normal patterns like large sparse allocations and fork()'s
   copy-on-write.
2. It kills a process to free memory when the system is critically low
   and can't satisfy an allocation; it picks a victim using a
   badness/oom_score heuristic weighted heavily by memory usage.
3. The OOM killer's heuristic favors killing whichever process is using
   the MOST memory at that moment, which can be the database rather
   than a smaller process that is actually leaking, unless scores are
   deliberately adjusted.
4. Kubernetes sets a Pod's memory `limit` as its cgroup's hard cap --
   exceeding it triggers the same kernel OOM mechanism, reported back
   as the Pod status "OOMKilled."

CHAPTER 14
1. Program's write() -> VFS -> filesystem driver -> page cache -> block
   layer/I/O scheduler -> device driver -> physical disk.
2. No -- write() returning success usually only means the data reached
   the page cache. fsync() (or O_DIRECT/O_SYNC) is what guarantees it's
   physically durable on disk.
3. That "successfully written" and "durably persisted" are different
   guarantees, and applications that assumed the former implied the
   latter lost data when the assumption broke.

CHAPTER 15
1. The inode (metadata/identity), the data blocks (actual content), and
   the directory entry / filename (a pointer to the inode).
2. The inode and its data blocks aren't freed until the last reference
   to them -- including any process still holding the file open -- is
   gone, so an open-but-unlinked file keeps consuming space.
3. Because inodes are a separate, finite resource from data blocks --
   if all inodes are used (e.g. by millions of tiny files), no new file
   can be created even with free space, since there's nowhere to record
   its metadata.
4. Journaling protects against a crash mid-write leaving the filesystem
   structure itself inconsistent, by logging intended metadata changes
   before applying them; recovery then just replays the journal instead
   of scanning the entire disk for inconsistencies.

CHAPTER 16
1. RAM the kernel uses to cache recently accessed disk data; it uses
   otherwise-idle RAM because unused RAM provides zero benefit, while
   cached data can turn a disk read into a much faster RAM read.
2. "Free" is genuinely unused RAM; "available" also includes page-cache
   memory that can be reclaimed instantly if a process needs it --
   available is the number that reflects real memory pressure.
3. Spinning disks have real mechanical seek latency, so scheduling
   (e.g. elevator algorithms) reorders requests to minimize head
   movement; SSDs have no seek penalty, so schedulers instead focus on
   maximizing parallel queue depth.

CHAPTER 17
1. An operation is not atomic if it can be interrupted partway through
   by another thread/process seeing or changing shared state -- this
   matters because the interruption can leave data in a state neither
   party intended.
2. Because races depend on precise timing that may only manifest under
   specific load, thread interleavings, or hardware -- tests that don't
   reproduce that exact timing pass every time until production
   conditions trigger it.
3. The bug depended on a specific, rare sequence of rapid keystrokes
   within a race window that testers were unlikely to reproduce, and
   the software had no hardware interlock to catch the resulting
   over-radiation dose.

CHAPTER 18
1. A mutex guarantees mutual exclusion -- only one thread holds it at a
   time, and typically only the holder can release it; a semaphore
   maintains a count and allows up to N holders, and can be signaled by
   a different thread than the one that waited on it.
2. Mutual exclusion, hold-and-wait, no preemption, circular wait; break
   the cycle by enforcing a consistent lock-acquisition order across
   all threads.
3. A low-priority task held a mutex needed by a high-priority task,
   while medium-priority tasks kept preempting the low-priority one --
   priority inversion; fixed with priority inheritance, temporarily
   boosting the lock-holder's priority.

CHAPTER 19
1. E.g. pipes (simple, one-directional, related processes), sockets
   (bidirectional, works across machines, more overhead), shared memory
   (fastest, but no built-in synchronization), message queues
   (structured messages, kernel-managed delivery).
2. A Unix domain socket stays entirely in the kernel without going
   through the full TCP/IP network stack (no packetization, checksums,
   or loopback network device), so it has lower overhead.
3. Anyone who can write to the Docker socket can effectively run
   arbitrary commands as root on the host via the Docker daemon --
   it's equivalent to root access.

CHAPTER 20
1. SIGTERM asks a process to terminate gracefully and can be caught/
   handled; SIGKILL forces immediate termination and is delivered
   directly by the kernel, never reaching the process's own code, so it
   cannot be caught, blocked, or ignored.
2. SIGTERM.
3. If it ignores SIGTERM, Kubernetes will wait out the grace period and
   then send SIGKILL, which doesn't give the app a chance to finish
   in-flight requests or close connections cleanly -- causing dropped
   requests during deploys.

CHAPTER 21
1. Power on -> firmware (UEFI/BIOS) -> boot loader (e.g. GRUB) ->
   kernel loads and initializes -> kernel starts PID 1 (init/systemd)
   -> PID 1 brings up services -> login prompt.
2. PID 1 is the first process the kernel starts; it becomes the
   ultimate parent because any process whose parent dies gets
   "re-parented" to PID 1, which is responsible for reaping them so
   they don't become permanent zombies.
3. It defends against a compromised or malicious boot loader/kernel
   being loaded -- each stage cryptographically verifies the next
   before running it, forming an unbroken chain of trust from firmware
   to OS.

CHAPTER 22
1. A unit file describes how systemd should manage something (a
   service, mount, timer, etc.); `systemctl start` runs it now,
   `systemctl enable` additionally creates the symlink that makes it
   start automatically at future boots.
2. A kernel module is code that can be loaded into the running kernel
   without a reboot; dynamic loading avoids needing every possible
   driver compiled into every kernel, and lets hardware support be
   added or removed as needed.
3. Advantage: parallelized, dependency-aware startup made boot much
   faster than SysVinit's sequential scripts. Criticism: it's a large,
   complex system taking on much more than "just init," violating the
   Unix "do one thing" philosophy in some critics' view.

CHAPTER 23
1. Most kernel-exposed resources -- regular files, directories, devices,
   pipes, sockets, and process information (/proc) -- are all accessed
   through the same read/write/open file-descriptor interface.
2. stdin (0), stdout (1), stderr (2).
3. Because /proc exposes process information as regular files/
   directories (e.g. /proc/<pid>/status) -- ps just reads and parses
   those files, with no special kernel API required beyond the normal
   filesystem calls.

CHAPTER 24
1. /etc holds static configuration files; /var holds data that changes
   at runtime (logs, databases, spool queues, caches).
2. /root is the superuser's home directory, kept separate from /home so
   it remains accessible even if /home is a separate, unmounted, or
   network filesystem during early boot/recovery.
3. The UsrMerge moved /bin, /sbin, /lib (and variants) to be symlinks
   into their /usr equivalents, unifying what had become a historically
   arbitrary split between "essential early-boot" and "everything else"
   binaries, simplifying packaging and read-only-root setups.

CHAPTER 25
1. `-rwx r-x r--`: file type (-, regular file), owner has read/write/
   execute, group has read/execute, others have read-only.
2. A hard link is a second directory entry pointing at the SAME inode
   (indistinguishable from the "original," can't cross filesystems); a
   symbolic link is a separate small file containing a path string,
   which can cross filesystems and can point at nonexistent targets.
3. SUID makes a program run with its owner's (often root's) privileges
   regardless of who invokes it. Legitimate use: `passwd`, so a normal
   user can update the root-owned password file in a controlled way.
   Risky use: a SUID-root script or binary with a shell-injection or
   path bug becomes an instant privilege escalation.

CHAPTER 26
1. `cp` needs `-r` because copying a directory means recursively
   copying its contents; `rm` deletes a bare file with no recursion
   needed, but by convention still requires `-r` for directories --
   the asymmetry the question points at is that `rm -rf` easily deletes
   a whole tree with one flag combination, unlike cp's need for content
   duplication logic.
2. A script computed a path using an unset/emptied environment
   variable, which collapsed the intended path down to `rm -rf /`
   equivalent, deleting content from the root of the filesystem.
3. Using `set -u` (fail on unset variables) and, for any risky rm/mv/
   dd command, first replacing the destructive command with `echo` to
   print exactly what WOULD run -- catching a wrong or empty variable
   before it does damage.

CHAPTER 27
1. It inverts the match, printing lines that do NOT match the pattern
   -- useful for filtering out noise (e.g. `grep -v DEBUG` on a log).
2. `sed 's/a/b/'` without the trailing `/g` flag only replaces the
   FIRST match per line; adding `g` (global) replaces every match on
   each line.
3. It's a real, documented exchange between two giants of computer
   science (Doug McIlroy and Donald Knuth) where a short pipeline of
   small Unix tools solved a text-processing task as effectively as a
   much longer bespoke program -- a concrete demonstration of "small
   tools, composed" over "one big program."

CHAPTER 28
1. It redirects file descriptor 2 (stderr) to wherever file descriptor
   1 (stdout) currently points.
2. Order matters because redirections are applied left to right: `>
   file 2>&1` first points stdout at the file, then points stderr at
   wherever stdout NOW points (the file) -- both go to the file.
   `2>&1 > file` first points stderr at wherever stdout CURRENTLY
   points (the terminal), then redirects stdout to the file -- stderr
   stays on the terminal.
3. A script fetched via `curl | sh` can change between the moment you
   inspect it and the moment it executes (or serve different content
   based on request headers/IP), so what you reviewed isn't guaranteed
   to be what runs; the safer alternative is downloading it to a file,
   reading it, then executing that file.

CHAPTER 29
1. `find` searches the live filesystem in real time (always accurate,
   slower on large trees); `locate` searches a prebuilt index (very
   fast, but can be stale if the index hasn't been updated since a
   recent change).
2. `which` shows the path of the executable that would run based on
   PATH lookup; `type` also reports if the command is actually a shell
   builtin, alias, or function rather than an external binary.
3. Because filtering by modification/access time (`find / -mtime -1`)
   quickly surfaces exactly which files changed around the time of a
   suspected compromise or incident, narrowing an otherwise huge
   filesystem down to a short, relevant list.

CHAPTER 30
1. Dependency metadata (what else it needs), version info, install
   scripts (pre/post-install hooks), and a manifest of exactly which
   files it installs (enabling clean removal).
2. Without automated dependency resolution, installing one piece of
   software could require manually finding, building, and installing
   an unpredictable chain of other libraries at compatible versions --
   often conflicting with what other installed software needed.
3. `apt remove` uninstalls the package's files but leaves its
   configuration files in place; `apt purge` also removes those
   configuration files.

CHAPTER 31
1. `ps` gives a single point-in-time snapshot; `top` continuously
   refreshes, showing live, changing resource usage.
2. Because %CPU is measured per core available -- a process using two
   full cores on a multi-core machine shows 200%, not capped at 100%.
3. A fork bomb is a process that recursively spawns copies of itself
   until it exhausts the system's process table/resources; because it
   happens WHILE it's running, you need process tools (and often `kill`
   with the right scope, or a ulimit on max processes) that still work
   even as the system is being overwhelmed.

CHAPTER 32
1. SIGHUP (hang-up) -- historically sent when the controlling terminal
   closes, and by default it terminates the process unless handled or
   disowned.
2. It makes a process ignore SIGHUP and redirects its output to a file
   (nohup.out by default), so it survives the terminal closing.
3. tmux/screen keep a persistent session you can detach from and
   reattach to later (including to see live output and interact again),
   whereas nohup only lets the process survive -- you can't get an
   interactive session back.

CHAPTER 33
1. `Restart=always` tells systemd to automatically restart the service
   if it exits for any reason; leaving it unset means a crash leaves
   the service down until someone notices and restarts it manually.
2. `systemctl start` runs it now, one time; `systemctl enable` makes it
   start automatically on future boots (they're independent -- you can
   do either without the other).
3. Because `Restart=` wasn't configured, so after the process crashed,
   nothing brought it back up, and no alert fired until a person or
   monitor eventually noticed it was down.

CHAPTER 34
1. The journal is a structured, binary, indexed, queryable log store
   managed by systemd (`journalctl`); traditional `/var/log` files are
   plain text files each application/daemon writes to independently.
2. logrotate periodically compresses, renames, and eventually deletes
   old log files based on size/age rules, so logs don't grow
   unbounded and fill the disk.
3. Because a full disk isn't just "no more logs" -- once free space
   hits zero, writes fail system-wide, which can crash databases,
   block new connections, and break unrelated services that also need
   to write to disk, turning a logging issue into a much bigger outage.

CHAPTER 35
1. sudo logs who ran what and when, and can be scoped to specific
   commands per user, whereas a shared root login gives everyone
   unrestricted, unattributable access.
2. A critical sudo heap-overflow vulnerability (found in 2021,
   present for roughly a decade) that let a local unprivileged user
   gain root; significant because sudo is nearly universally installed
   and the flaw required no special configuration to exploit.
3. Because `visudo` validates the syntax before saving and locks the
   file during editing -- a syntax error written directly to
   /etc/sudoers can break sudo entirely, potentially locking everyone,
   including yourself, out of privileged access.

CHAPTER 36
1. Capabilities let a process be granted ONE specific privileged
   operation (e.g. binding to a low port, raw sockets) without needing
   full root -- SUID-root grants ALL of root's power just to get one
   privileged operation.
2. CAP_NET_RAW grants the ability to create raw sockets; ping needs it
   (via SUID or the capability directly) to construct raw ICMP packets,
   which is otherwise a privileged operation.
3. An ACL can grant specific permissions to specific additional users
   or groups beyond the single owner and single group that basic Unix
   permissions allow.

CHAPTER 37
1. Which network interface/gateway to send packets to, based on
   destination IP/subnet.
2. Inside a container, loopback (127.0.0.1) is private to that
   container's own network namespace -- it does NOT reach services in
   sibling containers or the host, unlike on a normal host where
   loopback reaches everything bound to it.
3. Check locally first (is the process even listening, `ss -tlnp`),
   then local network path (routing, firewall), then DNS resolution,
   then the remote host/service itself.

CHAPTER 38
1. `ss` replaced `netstat`, because it reads directly from kernel data
   structures and is significantly faster on systems with many open
   sockets.
2. `dig +short` gives just the resolved IP(s)/record value with no
   extra formatting; the fuller output additionally shows query
   timing, the authoritative server used, TTLs, and flags -- useful
   when diagnosing DNS-specific issues rather than just "what does this
   resolve to."
3. Check locally (is anything listening, `ss`), then routing/firewall
   locally, then DNS resolution (`dig`), then reachability to the
   remote host/port directly (`curl -v`, or `tcpdump` to see if packets
   even arrive).

CHAPTER 39
1. `set -e` makes the script exit immediately if any command fails;
   without it, a script continues past a failed command as if nothing
   went wrong, potentially compounding the error.
2. `set -u` treats referencing an unset variable as an error and exits
   immediately, which would have stopped a script from silently
   collapsing a path down to `/` the way Chapter 26's incident did.
3. A deployment script accidentally activated old, dead feature-flagged
   code on production servers because of a manual, error-prone,
   partially-completed deployment process across many servers, not
   simply "a bug in the trading logic" -- the failure was fundamentally
   a deployment-process failure.

CHAPTER 40
1. Cron jobs run with a minimal environment (no login shell, a bare-
   bones PATH, no interactive profile variables) -- a script relying on
   PATH entries or environment variables set in your interactive shell
   silently fails or behaves differently under cron.
2. Use absolute paths for every command and file, and explicitly set
   any required environment variables (or source a known environment
   file) at the top of the script rather than relying on inherited
   shell state.
3. A dead man's switch expects a periodic "I'm alive and I ran
   successfully" signal (e.g. a ping to a monitoring endpoint) and
   alerts when that signal STOPS arriving -- it catches a cron job that
   silently stops running entirely, which "alert on error output"
   can't, since a job that never runs produces no error output either.

CHAPTER 41
1. It counts processes that are running OR waiting for CPU, and on
   Linux also processes in uninterruptible sleep (waiting on I/O) --
   so a high load average can reflect I/O pressure, not just CPU
   demand.
2. I/O pressure -- processes are blocked waiting on disk/network I/O
   rather than waiting for CPU time.
3. Because "load average: 8" means something very different on a
   4-core machine (heavily oversubscribed) versus a 32-core machine
   (mostly idle) -- it must be compared against core count to be
   meaningful.

CHAPTER 42
1. `strace` shows the syscalls a process makes (and their arguments/
   return values/timing); `lsof` shows what files, sockets, and other
   descriptors a process currently has open -- one shows activity over
   time, the other shows current state.
2. Check what it's doing right now (`strace -p`), what it has open
   (`lsof -p`), and its process state (`ps`, is it D-state?) -- in
   that rough order.
3. `df` reports free space based on the filesystem's own accounting;
   `du` reports space actually used by files it can see and walk. They
   disagree when a file is deleted but still held open by a running
   process -- `lsof | grep deleted` resolves the mystery by finding
   that file.

CHAPTER 43
1. `ulimit` sets per-process resource limits, configured per shell/
   process and inherited by children; a cgroup limits and accounts
   for an entire GROUP of processes together, enforced by the kernel
   independent of any one process's own settings.
2. Because each incoming connection typically consumes at least one
   file descriptor (a socket), and the default per-process open-file
   ulimit (often 1024) is far lower than the number of concurrent
   connections a busy server needs to hold open.
3. A container's memory/CPU limit (`docker run --memory`, or a
   Kubernetes Pod's `resources.limits`) is implemented using exactly
   this chapter's cgroup mechanism -- the container runtime creates a
   cgroup and applies the limit to it.

CHAPTER 44
1. PID namespace, network namespace, mount namespace (with a root
   filesystem), and a cgroup.
2. Because a container's processes are still ordinary Linux processes
   running under the SAME host kernel, just with a namespaced view --
   the host's process table sees all of them, just as it would see any
   other process.
3. Docker did not invent namespaces or cgroups (both existed in the
   Linux kernel beforehand, used by tools like LXC); Docker's real
   contribution was packaging these primitives with an easy build/
   ship/run workflow (Dockerfile, image format, registry) that made
   containers accessible to ordinary developers.

CHAPTER 45
1. A VM virtualizes hardware -- each VM runs its own full kernel on
   emulated/virtualized hardware; a container virtualizes the OS's
   process view -- containers share one host kernel, isolated via
   namespaces and cgroups.
2. Because starting a container just means starting a new process
   (with namespaces applied) under an already-running kernel, whereas
   a VM must boot an entire separate kernel and OS from scratch.
3. They trade some of a container's speed/density advantage for
   stronger isolation (often via a VM-like boundary or intercepted
   syscalls) -- a team chooses this when running less-trusted code
   (e.g. multi-tenant workloads) where plain container isolation feels
   insufficient.

CHAPTER 46
1. A layer is one immutable filesystem diff in an image; layers are
   cacheable/shareable because multiple images built from a common base
   or common early steps can reuse the exact same layer instead of
   storing/downloading it again.
2. The writable layer holds all changes made while the container ran;
   it is deleted along with the container when the container is
   removed, unless that data lives in a volume.
3. Because Docker caches layers and invalidates a layer (and every
   layer after it) as soon as its inputs change -- putting rarely-
   changing dependency installation first means routine application
   code changes don't force a slow dependency reinstall on every build.

CHAPTER 47
1. Checks locally for the image; if missing, pulls it (all its layers)
   from the registry; creates a new container (namespaces, cgroup,
   writable layer) from the image; starts the container's main
   process.
2. `docker stop` sends SIGTERM (then SIGKILL after a timeout) to halt
   a running container, but the container and its filesystem still
   exist; `docker rm` deletes a (stopped) container and its writable
   layer entirely.
3. Because `docker exec` starts an ADDITIONAL process inside the
   container's namespaces -- the container's original main process
   (PID 1 inside the container) is unaffected by that shell exiting.

CHAPTER 48
1. Because Docker caches each instruction's layer and invalidates
   everything after the first changed layer -- installing dependencies
   before copying frequently-changing application code means routine
   code changes don't force a slow dependency reinstall.
2. It limits the blast radius of a container escape or a vulnerability
   in the containerized application -- an attacker who breaks out as a
   non-root user lands with far fewer host privileges than one who
   breaks out as root.
3. It tells Docker which files/directories to exclude from the build
   context sent to the daemon -- preventing large or sensitive files
   (like `.git`, `node_modules`, or credentials) from being
   unnecessarily included or accidentally baked into an image layer.

CHAPTER 49
1. It's deleted along with the container's writable layer -- gone
   permanently unless it was stored in a volume or bind mount instead.
2. A named volume is managed by Docker itself (Docker decides where it
   lives on the host, referenced by name); a bind mount points directly
   at a path you specify on the host filesystem.
3. Because volumes are meant to potentially outlive any one container
   (e.g. a database's data surviving a container recreate/upgrade) --
   automatically deleting them with their last container risks silent,
   irreversible data loss.

CHAPTER 50
1. It publishes the container's port 80 to the host's port 8080,
   making it reachable from outside the host; without it, the port is
   only reachable from other containers on the same Docker network (or
   not at all, depending on network mode).
2. Docker runs an embedded DNS server for user-defined networks that
   automatically resolves container names to their current IP,
   updating as containers are created/removed.
3. The default bridge network is a legacy setup that predates this
   embedded DNS feature and doesn't provide it -- user-defined networks
   were introduced specifically to add automatic name resolution and
   better isolation.

CHAPTER 51
1. It guarantees the dependency's CONTAINER has started, not that the
   service inside it is actually ready to accept connections/requests.
2. Named volumes are preserved after `docker compose down`; only
   containers and the default network are removed. `docker compose
   down -v` additionally removes volumes.
3. Building a real readiness/health check into the dependent service's
   startup logic (retry with backoff until the dependency responds) or
   using Compose's `healthcheck` + `condition: service_healthy` --
   this directly foreshadows Kubernetes readiness probes (Ch60).

CHAPTER 52
1. A tag (like `:latest` or `:v2`) can be reassigned to point at a
   different image later; a digest is a cryptographic hash of the
   image's exact content, so pinning by digest guarantees you always
   get the exact same bytes.
2. Registering a package/image name deliberately similar to a popular,
   trusted one (a typo or lookalike), hoping developers will
   accidentally pull the malicious one instead -- it works because
   people often type or copy image names without double-checking them
   carefully.
3. It indicates Docker (or a verified organization) has reviewed and
   vouches for that image's publisher/maintenance -- it is NOT a
   guarantee the image is free of vulnerabilities, which still requires
   active scanning.

CHAPTER 53
1. Scheduling containers across many machines, restarting failed
   containers automatically, load-balancing traffic across replicas,
   rolling out updates without downtime, and scaling based on demand --
   none of which plain Docker on one machine does for you.
2. Kubernetes controllers continuously compare the cluster's actual
   state to the desired state you declared, and take action to close
   any gap -- this loop is how self-healing, scaling, and rollouts all
   work, driven by the same underlying mechanism.
3. Google's internal cluster management system, Borg, directly inspired
   Kubernetes' design (several of Kubernetes' original creators had
   worked on Borg at Google).

CHAPTER 54
1. The control plane decides what SHOULD be running cluster-wide (API
   server, etcd, scheduler, controller manager); a node is a machine
   that actually RUNS the workloads (kubelet, container runtime,
   kube-proxy).
2. The kubelet runs on each node and keeps existing Pods running
   locally even if it temporarily can't reach the control plane --
   already-scheduled workloads keep running through a control-plane
   outage, just without new scheduling/scaling decisions.
3. Because already-running Pods, managed by the kubelet on each node,
   keep serving traffic independently of the control plane -- a
   control-plane issue mainly affects the ability to make NEW changes
   (deploys, scaling, rescheduling), not already-running workloads.

CHAPTER 55
1. A Pod represents one or more containers that are always scheduled
   and run together, sharing a network namespace; it's usually not
   created directly because Pods aren't self-healing on their own --
   if one dies, nothing recreates it unless a higher-level controller
   manages it.
2. A Deployment manages ReplicaSets (creating a new one on each update);
   a ReplicaSet ensures a specified number of identical Pod replicas
   exist at all times.
3. Because a rolling update creates a NEW ReplicaSet for the new
   version while scaling the OLD ReplicaSet down to zero rather than
   deleting it immediately -- the old ReplicaSet is kept around
   (scaled to 0) to support fast rollbacks.

CHAPTER 56
1. It gives a stable, unchanging network identity (an IP and DNS name)
   in front of a set of Pods whose individual IPs change constantly as
   they're created and destroyed.
2. It uses a label selector to continuously match Pods with matching
   labels, and yes, this updates automatically as matching Pods come
   and go -- no manual re-registration is needed.
3. ClusterIP is reachable only inside the cluster; NodePort additionally
   exposes a port on every node's own IP; LoadBalancer additionally
   provisions an external cloud load balancer routing to the service.

CHAPTER 57
1. Practically, they're used for different intent (config vs.
   sensitive data), Secrets get some extra protections (like not being
   shown in plain `kubectl get` output by default, and support for
   encryption-at-rest when the cluster is configured for it), but
   neither is truly encrypted by default -- both are equally readable
   to anyone with API access.
2. Because base64 encoding LOOKS like security (the values aren't
   plainly readable in raw manifest text) but is trivially reversible
   by anyone with basic tooling -- leading many people to falsely
   assume Secrets are automatically safe from anyone with cluster
   access.
3. Enabling encryption-at-rest for etcd, restricting RBAC access to
   Secrets, and/or using a dedicated external secrets manager (like
   Vault or a cloud KMS-backed secrets integration) rather than relying
   on Kubernetes Secrets alone.

CHAPTER 58
1. A request is what the scheduler guarantees is reserved for a Pod
   when placing it on a node; a limit is the hard ceiling the Pod is
   not allowed to exceed (enforced by the kernel's cgroup mechanism).
2. It's not an error -- it means no node currently has enough
   unreserved capacity to satisfy the Pod's requests; it's the
   scheduler correctly refusing to overcommit a node beyond what it
   promised other Pods.
3. HPA changes how many Pod REPLICAS exist, based on metrics; Cluster
   Autoscaler changes how many NODES exist, based on whether Pods are
   unschedulable due to lack of capacity. You typically need both
   together because HPA can create more Pods than existing nodes have
   room for, and only Cluster Autoscaler can add the capacity to
   actually run them.

CHAPTER 59
1. Lower the Pod's CPU request (may mean it gets throttled sooner) or
   add more/bigger nodes (costs more) -- the trade-off is resource
   sufficiency for THIS Pod versus cost/capacity for the cluster.
2. A container can crash before its logging library finishes
   initializing, leaving current logs empty; `kubectl logs --previous`
   shows the logs from the last (crashed) instance instead of the
   current, freshly-restarted one.
3. The Service's label selector doesn't actually match the Pod's
   labels, or the Pod's readiness probe is failing intermittently so
   it's cycling in and out of the Service's endpoint list.

CHAPTER 60
1. A failed liveness probe causes Kubernetes to RESTART the container;
   a failed readiness probe removes the Pod from the Service's traffic
   endpoints WITHOUT restarting it.
2. It enables an outage where the app is actually broken (can't serve
   real requests) but Kubernetes has no way to know, because both
   liveness and (if it exists) readiness are checking something too
   shallow (just "is the HTTP server alive") to reveal the real
   problem.
3. Because the rollout only proceeds to remove more old Pods once new
   Pods pass their readiness probe -- if new Pods never become ready,
   Kubernetes correctly stops there rather than continuing to replace
   working old Pods with broken new ones.

CHAPTER 61
1. It re-applies the exact rendered Kubernetes YAML that was generated
   and applied during revision 3's `helm install`/`upgrade` -- Helm
   stores each revision's rendered manifest.
2. Because the next `helm upgrade` re-applies Helm's own tracked state,
   silently overwriting the manual `kubectl edit` change with no
   warning that it's discarding something.
3. Because it shows exactly what Kubernetes objects, RBAC permissions,
   and resource settings the chart would create, letting you review
   for problems (missing limits, excessive permissions, disabled
   persistence) before anything is actually applied to your cluster.

CHAPTER 62
1. It tells you the program touched many pages that were already in
   RAM but not yet mapped into its own address space (common and
   cheap, e.g. shared libraries, copy-on-write pages) -- generally not
   a concern on its own, unlike major faults which indicate disk I/O.
2. Because each syscall has fixed kernel-transition overhead regardless
   of how much data it moves -- doing the same total I/O in fewer,
   larger reads amortizes that fixed overhead across far more bytes
   per call.
3. The program may load the file plus significant additional
   structures (e.g. a hash map of every word) into memory
   simultaneously, or it may buffer/duplicate data during processing --
   either can push total memory usage well past the raw file size.

CHAPTER 63
1. Because the final stage only copies the already-built artifact from
   the builder stage, discarding the builder stage's compilers, build
   tools, and intermediate files entirely -- those never become part
   of the final image's layers.
2. Because `depends_on` only waits for the dependency's CONTAINER to
   start, not for the database process inside it to finish its own
   startup and begin accepting connections; fix with an application-
   level retry/backoff on connect, or a Compose healthcheck combined
   with `condition: service_healthy`.
3. The base image's default user might already be root and the
   Dockerfile's `USER` instruction may be placed before another
   instruction (or a later stage in a multi-stage build) that resets
   it, or `docker exec` itself may have been run with an explicit
   `--user root` override.

CHAPTER 64
1. Because it's the only way to genuinely verify the rollout safety
   mechanism (readiness gating, stalled rollout, rollback) actually
   works, rather than assuming it does -- deploying only working
   versions never exercises the failure path at all.
2. A readiness probe failure keeps old, working Pods serving traffic
   while new ones are blocked -- users see no impact; a liveness probe
   failure restarts containers, which, if misconfigured, can cause
   repeated restarts and real user-visible disruption.
3. It assumes the Pod needs a very small, effectively negligible amount
   of each resource, which can cause the scheduler to over-pack a node,
   leading to real resource contention once the Pod's actual usage
   turns out to be much higher.

CHAPTER 65
1. Because restoring service quickly limits user impact, and a known-
   good previous version is a fast, low-risk action -- continuing to
   investigate live while users are affected extends the outage for
   the sake of understanding that can just as well happen after
   service is restored.
2. It focuses on what the SYSTEM and PROCESS allowed to happen (missing
   alerts, unclear runbooks, a gap in testing) rather than which
   individual made a mistake -- assuming good faith and treating the
   incident as a systems-design learning opportunity.
3. Config/selector/probe-path/scheduling-related failures (a, c, d, e)
   would generally show up quickly in `describe pod`'s Events; a slow
   memory leak leading to eventual OOMKilled (b) or a disk-fill (f)
   might require watching trends over time or checking application
   logs/metrics before the root cause is obvious.

CHAPTER 66
1. Configuration drift is the gap between the intended baseline and the
   actual state of one or more machines: package versions, sysctls,
   mounts, unit files, kernel versions, runtime state, or manual fixes.
2. A known-good host narrows the search space. Differences between good
   and bad hosts often reveal partial rollouts, missing packages,
   different kernels, changed mount options, or runtime drift faster
   than reading every log line on the bad host.
3. A mitigation reduces immediate impact; a root cause explains why the
   failure happened. Restarting a service may mitigate an outage, but it
   does not explain the leak, dependency failure, or bad rollout that
   made the restart necessary.

CHAPTER 67
1. systemd starts a service with a different environment, working
   directory, user, PATH, limits, dependencies, and sandboxing than your
   interactive shell.
2. `systemctl cat` shows the effective unit file plus drop-in overrides;
   `status` shows current state, recent logs, and process information.
3. Hardening settings change what the process can read, write, execute,
   and access. They can prevent real attacks, but they can also break
   legitimate behaviour if untested.

CHAPTER 68
1. A backup job only proves bytes were written somewhere. A tested
   restore proves you can recover the system or data when it matters.
2. You need the device, partition/LVM/RAID/encryption stack, filesystem
   type, mount point, current usage, backup state, and whether the
   filesystem supports online growth.
3. If a process still has the deleted file open, the directory entry is
   gone but the disk blocks remain allocated until that file descriptor
   is closed.

CHAPTER 69
1. The service may be bound only to loopback (`127.0.0.1`), blocked by a
   firewall/security group, routed incorrectly, or listening in a
   different network namespace.
2. A network namespace isolates interfaces, addresses, routes, firewall
   state, and sockets, giving a process its own view of the network.
3. Packet capture is powerful but can be noisy and intrusive. Often
   `ss`, routing, DNS, and firewall inspection identify the layer before
   you need packets.

CHAPTER 70
1. The symptom tells you which layer to measure. Without it, you can
   collect impressive data that does not answer the user-impacting
   question.
2. A flame graph shows where time is spent across call stacks, making
   hot paths visible without reading the whole codebase first.
3. Host metrics describe machine/resource behaviour; distributed traces
   show how one request moves through services and where time/errors
   occur along that path.

CHAPTER 71
1. Capabilities allow a narrow privilege, such as binding a low port,
   without granting every privilege that comes with full root.
2. Mode bits mostly answer owner/group/other access. SELinux/AppArmor
   can enforce policy based on process domain/profile, file labels,
   paths, capabilities, and allowed behaviours even when Unix mode bits
   would otherwise permit access.
3. Logs are useful only if they answer a detection or investigation
   question and someone or something reviews them. Unread logs are
   storage, not control.

CHAPTER 72
1. Restarting can erase volatile evidence: process state, memory growth,
   open files, socket state, stack traces, and the exact failure mode.
2. Mitigation reduces immediate impact; remediation fixes the specific
   defect; prevention changes tests, alerts, architecture, capacity, or
   process so the same class of failure is less likely or easier to
   catch.
3. Add the commands that worked, healthy/bad examples, known safe
   actions, escalation criteria, rollback steps, and any alert or
   dashboard links that would have shortened the incident.

CHAPTER 73
1. G = goroutine (your code), M = an OS thread, P = a processor slot
   allowed to run Go code; there are GOMAXPROCS Ps, any number of Gs,
   and as many Ms as currently needed.
2. A goroutine waiting on a socket is parked and its descriptor is
   registered with epoll (the netpoller); one thread waits for all of
   them and wakes only the goroutines whose sockets become ready.
3. Regular-file I/O can't use epoll on Linux, so a read on a stalled
   disk or NFS mount blocks its OS thread inside the kernel; the
   runtime starts more threads to keep running other goroutines.

CHAPTER 74
1. The kernel re-parents orphaned processes to PID 1 (of their PID
   namespace), so PID 1 is the only process that can wait() for them.
2. It runs a tiny init (tini) as PID 1, which reaps zombies and forwards
   signals to your program, which is no longer PID 1.
3. So the child and everything it started share a process group, and a
   kill to the negative PGID reaches grandchildren too.

CHAPTER 75
1. Quota and period: microseconds of CPU time allowed per period, e.g.
   150000 100000 = 150 ms per 100 ms = 1.5 CPUs.
2. Throttling is decided per 100 ms period: many threads can use the
   whole quota early in a period and then stall for the rest, even
   though the average over a second is below the limit.
3. Go 1.25+ reads the cgroup CPU limit and sets GOMAXPROCS to it, rounded
   up, minimum 2, re-checked periodically; the GOMAXPROCS environment
   variable overrides it.

CHAPTER 76
1. At roughly twice the live heap: about 400 MB.
2. It is a target the GC works towards (by collecting more often), not
   a cap the runtime enforces; live data larger than the limit still
   grows past it.
3. 137 = 128 + 9: the process was killed by SIGKILL, typically by the
   kernel's OOM killer when the cgroup reached memory.max.

CHAPTER 77
1. Only that the bytes were copied into the kernel's page cache; they
   may not be on stable storage yet.
2. rename(2) is atomic only within one filesystem; across filesystems
   it becomes a copy that can be interrupted.
3. Losing the rename itself in a crash -- the directory entry change
   is metadata that is only durable once the directory is fsynced.

CHAPTER 78
1. The soft limit is enforced; the hard limit is the ceiling an
   unprivileged process may raise its own soft limit to.
2. It raises its soft RLIMIT_NOFILE to the hard limit minus one, and
   restores the original soft limit for processes it starts.
3. All network descriptors are registered with one epoll instance;
   waiting goroutines are parked, and the runtime is woken only for
   descriptors that are ready.

CHAPTER 79
1. The command name (field 2) is in parentheses and may contain spaces
   and parentheses, so only the LAST ')' reliably marks its end.
2. Read utime+stime (fields 14 and 15) twice, subtract, divide by
   USER_HZ (100) and by the interval in seconds, times 100.
3. Kernel threads have no user-space command line; tools show their
   name from the comm field in square brackets instead.

CHAPTER 80
1. CLONE_NEWPID (the child is PID 1 of a new PID namespace); CLONE_NEWUTS
   (its own hostname).
2. /proc reflects the PID namespace it was mounted in; without a fresh
   mount, ps inside would show the host's processes.
3. The cgroup's pids.max = 20; the refusals are counted in pids.events
   ("max N") in the cgroup directory.

CHAPTER 81
1. flat = time spent in the function itself; cum = time in the function
   plus everything it called.
2. The goroutine profile; "489 @ ..." means 489 goroutines share exactly
   that stack -- here, all blocked at the same line.
3. Time inside the kernel (TCP stack, page faults, scheduler) and in
   non-Go code, using hardware performance counters.

CHAPTER 82
1. GOMEMLIMIT (Chapter 76), relying on the CPU limit for GOMAXPROCS
   (Chapter 75), a SIGTERM-aware drain with preStop and a grace period
   (Chapter 74), CGO_ENABLED=0 for a static distroless binary (Chapter 48).
2. Distroless images contain no /bin/sleep and the server has no sleep
   flag, so an exec hook would fail; the built-in sleep action needs
   nothing in the image.
3. With DynamicUser the UID changes on every start, so no file owner can
   be set in advance; LoadCredential copies the root-only files into a
   per-service directory only this service can read.

CHAPTER 83
(Personal/reflective -- no single correct answer; the value is in
articulating your own reasoning using this guide's vocabulary.)
3. A service mesh typically adds mutual TLS between services,
   fine-grained traffic routing/splitting (e.g. canary percentages),
   and detailed per-request observability -- none of which a plain
   Service's simple label-based load balancing provides on its own.

CHAPTER 84
1. Because process/memory/syscall fundamentals have been stable for
   decades and change slowly, while Kubernetes, container runtimes, and
   distributions regularly change APIs, defaults, and operational
   mechanics.
2. E.g. the Kubernetes deprecated API migration guide, checked before
   an upgrade -- avoids a manifest using a removed `apiVersion` silently
   failing to apply after the upgrade.
3. The underlying cgroup mechanism switched from cgroups v1's separate
   per-resource hierarchies to v2's single unified hierarchy; some
   memory-accounting edge cases changed as a result, so containers
   could hit OOM-killed at effectively different real memory thresholds
   than before, purely from the OS-level cgroup version underneath
   them changing, without any change to the container's own configured
   limit.
```
