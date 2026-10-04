# Real-Life Scale, Load & Performance Testing — In-Depth Guide

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/scale-perf/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

A practitioner's reference for reasoning about network behavior, latency, and load —
and for building simulations that prove a system will survive real traffic.
Covers open-source tools, hands-on simulations, and report generation, with code
in **C**, **Python**, **Go**, and **Rust** wherever the language fits the job best
(C/Rust for raw socket-level work, Go for concurrent HTTP load generation, Python
for orchestration, scripting, and analysis).

---

*This guide moves through three loose stages, in order: **§1–6** build the
mental model and vocabulary (where time goes, the metrics that matter, how
to model realistic traffic) with no tools yet; **§7–10** get hands-on
(specific tools, network-condition simulation, writing your own load
generator, turning results into reports); **§11–13** zoom back out to how
you'd operationalize this ongoing (SLOs, checklists, further reading). You
can stop after §6 if you just need the concepts, or skip to §7 if you
already know what you're measuring and just need a tool.*

## Table of Contents

1. [Mental Model: Where Time Goes](#1-mental-model-where-time-goes)
2. [Network Fundamentals That Determine Performance](#2-network-fundamentals-that-determine-performance)
3. [Latency & Load Testing Metrics — Deep-Dive Reference](#3-latency--load-testing-metrics--deep-dive-reference)
4. [Load & Traffic Modeling](#4-load--traffic-modeling)
5. [API & Endpoint Performance](#5-api--endpoint-performance)
6. [Real-World Scenarios & Case Studies](#6-real-world-scenarios--case-studies)
7. [Open-Source Load Testing Tools — Walkthroughs](#7-open-source-load-testing-tools--walkthroughs)
8. [Simulating Network Conditions](#8-simulating-network-conditions)
9. [Building Your Own Load Generators](#9-building-your-own-load-generators)
10. [Report Generation & Visualization](#10-report-generation--visualization)
11. [Metrics, SLIs/SLOs, and the RED/USE Methods](#11-metrics-slis-slos-and-the-red-use-methods)
12. [Checklists & Common Pitfalls](#12-checklists--common-pitfalls)
13. [Further Reading](#13-further-reading)

---

## 1. Mental Model: Where Time Goes

Every request a client makes crosses several layers, and load testing is really about
finding **which layer breaks first** as concurrency rises. Draw this stack mentally
before touching a tool:

```
Client                                                        Server
  |--DNS lookup-------------------------------------------------|
  |--TCP handshake (SYN, SYN-ACK, ACK)---------------------------|
  |--TLS handshake (ClientHello...Finished)----------------------|
  |--Request sent (serialize + write to socket)-------------------|
  |          [network transit: propagation + queueing delay]     |
  |                                              [server accept()]|
  |                                       [thread/goroutine pickup]|
  |                                         [business logic + I/O]|
  |                                          [DB / cache / RPC]   |
  |                                       [serialize response]    |
  |          [network transit back]                               |
  |--Response read + deserialize----------------------------------|
```

A load test result of "p99 = 800ms" is meaningless until you know which segment
dominates. The rest of this guide gives you the vocabulary and tools to decompose
that number.

---

## 2. Network Fundamentals That Determine Performance

### 2.1 TCP handshake and connection reuse cost

A cold TCP connection costs **1 RTT** before any bytes of your protocol move. Add TLS
1.2 (2 RTT) or TLS 1.3 (1 RTT, 0-RTT for resumption) and a client hitting a fresh
connection per request pays 2-3 RTTs of pure overhead. This is why **connection
pooling / keep-alive** is the single highest-leverage performance change most APIs
can make.

```
Fresh HTTPS/TLS1.2 connection over a 50ms RTT link:
  TCP handshake:   50ms
  TLS handshake:  100ms
  HTTP request:    50ms
  ---------------------
  Total before first byte: 200ms

Reused (keep-alive) connection, same link:
  HTTP request: 50ms   (4x faster, purely from reuse)
```

### 2.2 DNS

DNS resolution is often the invisible tail-latency source. A misconfigured resolver,
or a client re-resolving on every connection because it ignores TTLs (or a service
mesh sidecar doing DNS on every call), silently adds 20-200ms. Under load tests,
always resolve DNS once and reuse the IP, unless you're specifically testing DNS
behavior (e.g., DNS-based failover, GSLB).

### 2.3 TCP congestion control & the bandwidth-delay product

Throughput is capped by:

```
Throughput <= TCP_window_size / RTT
```

This is why a single TCP connection between two data centers 150ms apart, even on a
10 Gbps link, may only push a few Mbps until the congestion window grows (slow
start) or unless you use multiple parallel streams (HTTP/2 multiplexing, multiple
connections). This matters directly for load tests: **if your load generator opens
too few connections over a high-latency link, you will bottleneck the generator, not
the server under test.**

### 2.4 Nagle's algorithm & delayed ACK

Nagle's algorithm batches small writes; delayed ACK delays acknowledgments up to
~40ms waiting to piggyback data. Combined, they can produce a notorious ~40ms stall
on small, latency-sensitive request/response protocols. Fix: disable Nagle
(`TCP_NODELAY`) on latency-sensitive services — virtually every HTTP server and load
tool does this by default, but hand-rolled TCP clients (see the C example in
§9.1) must set it explicitly.

### 2.5 HTTP/1.1 vs HTTP/2 vs HTTP/3

| | Multiplexing | Head-of-line blocking | Handshake | Best for |
|---|---|---|---|---|
| HTTP/1.1 | No (needs 6 conns/host) | Yes, per-connection | TCP+TLS | legacy, simple |
| HTTP/2 | Yes, single TCP conn | Yes, at TCP layer (packet loss stalls all streams) | TCP+TLS | most APIs today |
| HTTP/3 (QUIC) | Yes, per-stream, over UDP | No (independent streams) | 1-RTT / 0-RTT | mobile/lossy networks |

Load testing implication: benchmarking an HTTP/2 API from a client that only opens
HTTP/1.1 connections will radically understate real client concurrency behavior (a
browser reuses one HTTP/2 connection for 100 requests; your tool must be configured
to do the same, or configured to open many connections if you're modeling many
distinct clients).

### 2.6 MTU, fragmentation, and packet loss

Standard Ethernet MTU is 1500 bytes. Large JSON/gRPC payloads that exceed the
path MTU get fragmented at the IP layer (or segmented at TCP layer, which is
normal and fine) — but on paths with a broken PMTU (common with some VPNs/tunnels),
fragmented packets get silently dropped, producing intermittent hangs that look like
"random" latency spikes. This is a classic hard-to-diagnose production incident —
worth explicitly testing with `ping -s 1472 -M do <host>` to check for PMTU
blackholes before blaming the application.

---

## 3. Latency & Load Testing Metrics — Deep-Dive Reference

### 3.1 Never trust the average

Latency distributions are heavy-tailed. Averaging hides the exact thing you care
about: how bad is the worst experience a real user hits?

```
Example: 1000 requests
  990 requests: 20ms
  10 requests:  2000ms   (GC pause, lock contention, slow DB query)

  Average  = (990*20 + 10*2000) / 1000 = 39.8ms   <- looks fine!
  p50      = 20ms
  p99      = 2000ms                                 <- the real story
```

Always report **p50, p90, p95, p99, p99.9**, plus max. For systems with SLAs, p99.9
or p99.99 is often the contractually meaningful number (1 in 1,000 or 1 in 10,000
users).

### 3.2 Tail latency amplification (fan-out)

If a single user request fans out to 20 backend calls, and each backend call has a
p99 of 100ms (i.e., 1% chance of being slow), the probability that *at least one* of
the 20 calls is slow is:

```
P(at least one slow) = 1 - (1 - 0.01)^20 ≈ 18%
```

So a "p99 of 100ms" backend turns into an effective p50-ish experience of slowness
for a fanned-out caller. This is why high-fan-out architectures (search, feed
aggregation, microservices with deep call graphs) must budget latency far more
conservatively per hop than the end-to-end SLA suggests. Mitigations: hedged
requests, timeouts + retries with jitter, request-scoped deadlines.

### 3.3 Little's Law — the load-testing Rosetta Stone

```
L = λ * W

L = average number of requests in the system (concurrency)
λ = average arrival rate (throughput, req/s)
W = average time a request spends in the system (latency)
```

This single equation tells you, given any two of {concurrency, throughput,
latency}, the third. It's how you sanity-check load test results:

```
Example: your tool reports throughput = 500 req/s, p50 latency = 200ms.
Expected concurrency in flight: L = 500 * 0.2 = 100 concurrent requests.

If your load generator was configured for 100 concurrent workers, and it produced
exactly ~500 req/s at ~200ms latency, the numbers are self-consistent — good.
If concurrency was configured at 1000 but throughput stayed at 500 req/s, your
server has hit a **concurrency ceiling** (thread pool, connection pool, DB pool
exhaustion) — requests are queueing somewhere, and the 200ms you're measuring is
now including queue wait time, not just service time.
```

### 3.4 Open-loop vs. closed-loop load generators

This is the most-misunderstood concept in load testing.

- **Closed-loop**: a fixed number of virtual users, each waits for a response before
  sending the next request (think of N people repeatedly hitting refresh). As the
  server slows down, the *arrival rate drops* to compensate — the tool self-throttles
  and hides how bad things really are. Most naive load-test scripts are
  accidentally closed-loop.
- **Open-loop**: requests arrive at a fixed rate (e.g., Poisson-distributed at 1000
  req/s) regardless of how fast the server responds. This matches real-world traffic
  — real users don't stop arriving because your server is slow. Open-loop testing
  reveals queueing collapse that closed-loop testing conceals.

**Rule of thumb**: for capacity planning and finding breaking points, use open-loop
(rate-based) load generation — e.g., `vegeta`, `k6` with `arrival-rate` executors,
`wrk2`. For simulating a fixed pool of real users browsing a site, closed-loop
(`ab`, plain `wrk`, `Locust` in its default mode) is fine.

### 3.5 Latency vs. response time vs. RTT — precise definitions

These three terms get used interchangeably in casual conversation but mean
different things, and mixing them up produces reports that don't compare
apples-to-apples:

| Term | Definition | Measured where |
|---|---|---|
| **Latency** (network sense) | One-way time for a packet to travel from sender to receiver | Network layer, hard to measure directly without synced clocks |
| **RTT (round-trip time)** | Time for a packet to go out and its acknowledgment/reply to come back | `ping`, TCP handshake timing |
| **Response time / request duration** | Wall-clock time from "client sends request" to "client finishes receiving response" — what load-testing tools actually report | Client-side, includes network + server processing |
| **TTFB (time to first byte)** | Time from request sent to the first byte of the response arriving | Reveals server-side processing delay before it starts streaming |
| **TTLB (time to last byte)** | Time until the *entire* response body is received | Matters for large payloads/downloads; TTLB − TTFB ≈ transfer time |
| **Service time** | Time the server itself spent actually processing (excludes network and any queueing before the server picked up the request) | Server-side instrumentation (APM/tracing), not visible to the client |
| **Queue/wait time** | Time a request spent waiting for a thread/connection/worker before service began | `response_time − service_time`, invisible unless the server reports it explicitly |

When a load report says "latency: 240ms," always ask: *is this response time
including network, or server-reported service time?* A load generator sitting in
the same AZ as the server and one running from a user's mobile device on the
other side of the planet will report wildly different "latency" for an identical
backend — because most of the number is §2's network fundamentals, not the
application.

### 3.6 Average (mean) — what it's actually good for

```
mean = (sum of all latencies) / (count of requests)
```

The mean is **not useless** — it's the right metric when you care about
*aggregate cost*, not the worst individual experience: total CPU-seconds
consumed, total time budget across a batch job, or cost-per-request billing
models. It is the *wrong* metric for user-experience SLAs precisely because it's
sensitive to outliers in one direction only (a few very slow requests drag it up)
while being insensitive to how many users actually experienced that slowness.
Never report a mean alone — always pair it with at least p50 and p99 so a reader
can tell whether the distribution is tight or heavy-tailed (§3.1's worked
example: mean 39.8ms with p99 of 2000ms is a very different system than a mean of
39.8ms with p99 of 45ms, even though the mean is identical).

### 3.7 Median (p50) vs. percentiles — computation methods matter

The median is "the middle value" — but *how* you find it, and how you find
p90/p95/p99/p99.9 above it, is not one universally agreed algorithm. Two common
methods, applied to 10 sorted latency samples (ms): `12, 15, 18, 20, 22, 25, 30,
45, 90, 250`:

**Nearest-rank** (used by many load-testing tools' simple implementations):
```
rank = ceil(p/100 * N)
p50 -> rank = ceil(0.50*10) = 5  -> 5th value = 22ms
p90 -> rank = ceil(0.90*10) = 9  -> 9th value = 90ms
```

**Linear interpolation** (numpy's default `np.percentile`, many stats libraries):
```
index = (N-1) * p/100
p50 -> index = 9*0.50 = 4.5 -> interpolate between 5th(22) and 6th(25) = 23.5ms
p90 -> index = 9*0.90 = 8.1 -> interpolate between 9th(90) and 10th(250) = 106ms
```

Notice p90 differs by 90ms vs. 106ms — a ~18% difference — purely from
methodology, with identical raw data. **Never compare a p99 from one tool against
a p99 from another tool without confirming they use the same computation
method**, and never compare against a hand-rolled percentile function without
checking which convention it follows.

At scale, sorting every raw sample is expensive (O(n log n) time, O(n) memory
for a multi-hour test generating millions of samples). Production-grade tools use
streaming approximate-percentile structures instead:

- **HdrHistogram** (used internally by `wrk2`, `gatling`, and available as a
  library for custom tools — see §9.4's use of it in Rust): buckets values into
  logarithmically-sized ranges, giving a fixed *relative* error (e.g., 3
  significant figures ⇒ ~0.1% error) across a huge dynamic range (microseconds to
  minutes) in O(1) memory and O(1) recording time per sample, regardless of how
  many samples you record.
- **t-digest**: a different approximate-quantile sketch, popular in
  monitoring/APM backends (Datadog, Prometheus's `summary` type historically),
  optimized for accuracy at the extreme tails (p99.9+) at the cost of coarser
  accuracy near the median.

Demonstrating the two methods directly in Python:

```python
import numpy as np

samples = [12, 15, 18, 20, 22, 25, 30, 45, 90, 250]

def nearest_rank_percentile(data, p):
    data = sorted(data)
    n = len(data)
    rank = max(1, int(np.ceil(p / 100 * n)))
    return data[rank - 1]

print("nearest-rank p90:", nearest_rank_percentile(samples, 90))       # 90
print("linear-interp p90:", np.percentile(samples, 90))                 # 106.0
```

### 3.8 Standard deviation, variance, and coefficient of variation

```
variance   = mean((x_i - mean)^2)
stddev     = sqrt(variance)
coeff. of variation (CV) = stddev / mean
```

Standard deviation quantifies spread around the mean, but on the heavy-tailed
distributions typical of latency data it's dominated by the same outliers that
make the mean misleading (§3.6) — it is rarely reported on its own for latency,
but the **coefficient of variation** (stddev normalized by the mean) is a useful
single number for comparing the *relative* stability of two systems with
different absolute latencies: a CV near 0 means very consistent latency, a CV
above ~1 signals a heavy tail or bimodal distribution worth investigating with a
full histogram rather than summary statistics alone.

### 3.9 Jitter — variation between consecutive requests

Percentiles describe the overall *shape* of a distribution but say nothing about
**order** — two systems can have identical p50/p95/p99 while one delivers steady,
predictable latency and the other swings wildly request-to-request. That
swinginess is **jitter**, and it matters independently of percentiles for
anything sensitive to *consistency* rather than just worst-case bound: VoIP,
video calls, multiplayer gaming, financial market-data feeds, and any UI that
does incremental/streaming rendering.

```
Same p50 (50ms), wildly different jitter:

System A (low jitter):  48, 51, 49, 52, 50, 49, 51, 50ms   — smooth, predictable
System B (high jitter): 5, 95, 10, 88, 3, 99, 15, 85ms     — same mean, unusable for real-time media
```

Two common ways to quantify it:

**1. Simple successive-difference jitter** (easy to compute from a load test's
raw latency series):
```
jitter = mean(|latency[i] - latency[i-1]|) for consecutive requests
```

**2. RFC 3550 (RTP) smoothed jitter estimate** — the formula VoIP/video systems
actually use internally, an exponentially-weighted moving average of the
inter-packet spacing difference between sender and receiver:

```
J(i) = J(i-1) + (|D(i-1,i)| - J(i-1)) / 16
```

where `D(i-1,i)` is the difference in relative transit time between two
consecutive packets. The `/16` smoothing constant (a 1/16 gain, per the RFC)
keeps the estimate stable against single outliers while still tracking real
trend changes.

```c
// rtp_jitter.c — RFC 3550 style jitter estimator over a stream of timestamps.
// send_ts[]/recv_ts[] are in the same clock units (e.g., RTP timestamp units).
#include <stdio.h>
#include <stdlib.h>
#include <math.h>

double compute_rfc3550_jitter(const double *send_ts, const double *recv_ts, int n) {
    double jitter = 0.0;
    double prev_transit = recv_ts[0] - send_ts[0];

    for (int i = 1; i < n; i++) {
        double transit = recv_ts[i] - send_ts[i];
        double d = fabs(transit - prev_transit);
        jitter += (d - jitter) / 16.0;   // RFC 3550 sec 6.4.1
        prev_transit = transit;
    }
    return jitter;
}

int main(void) {
    // simulated send times (steady 20ms cadence) vs. arrival times (network jitter)
    double send[] = {0, 20, 40, 60, 80, 100, 120, 140};
    double recv[] = {5, 30, 42, 90, 83, 140, 121, 200};
    int n = sizeof(send) / sizeof(send[0]);

    double j = compute_rfc3550_jitter(send, recv, n);
    printf("estimated jitter: %.3f (same units as timestamps)\n", j);
    return 0;
}
```

```bash
gcc -O2 -o rtp_jitter rtp_jitter.c -lm && ./rtp_jitter
```

**Simulating jitter for a load test**: `tc netem`'s second `delay` argument is
exactly this — `delay 100ms 20ms distribution normal` (§8.1) adds 100ms latency
with 20ms of jitter drawn from a normal distribution, letting you reproduce
System A/B above on a real network path and verify your client's jitter
buffering/retry logic under each condition.

### 3.10 Fan-out — the full math

§3.2 introduced the headline formula; here's the complete picture for both
fan-out shapes.

**Parallel fan-out (breadth)** — N independent calls issued concurrently, caller
waits for all N (a typical "gather" / `Promise.all` pattern):

```
P(at least one of N calls exceeds threshold T) = 1 - (1 - p)^N
```

where `p` is the single-call probability of exceeding T (e.g., p = 0.01 for a
p99 threshold).

| N (fan-out breadth) | P(at least one slow), p=1% | Effective experience |
|---|---|---|
| 1 | 1.0% | matches the backend's own p99 |
| 5 | 4.9% | |
| 10 | 9.6% | |
| 20 | 18.2% | roughly the backend's own p85 |
| 50 | 39.5% | roughly the backend's own p60 |
| 100 | 63.4% | worse than the backend's own median |

Beyond the failure-threshold framing, fan-out breadth also drags up the
*expected* latency even when nothing "fails" — the caller's total time is the
**maximum** of N samples (an order statistic), not the mean of one. For latency
distributions with an exponential-like tail, the expected maximum of N iid
samples grows roughly with `scale * ln(N)` — i.e., **fan-out breadth costs you
tail latency logarithmically, for free, even under perfectly nominal
conditions.** This is why search/aggregation systems that fan out to dozens or
hundreds of shards invest heavily in **hedged requests** (fire a duplicate
request to a second replica if the first hasn't returned within, say, the
p50 latency, and take whichever answers first) and in **tail-tolerant
timeouts + cancellation** rather than trying to make every single shard
uniformly fast.

**Sequential fan-out (depth)** — a chain of M dependent hops, each waiting for
the previous:

```
total_latency ≈ sum(latency_i) for i in 1..M          (means simply add)
total_variance ≈ sum(variance_i) for i in 1..M         (variances add for independent hops)
total_stddev = sqrt(total_variance)
```

Because variances (not standard deviations) add, the **relative** jitter of a
long chain shrinks somewhat as depth grows (central limit theorem smoothing),
but the **absolute** latency grows linearly with every hop — which is why deep
microservice call chains (A → B → C → D → E) are budgeted hop-by-hop in strict
latency budgets (each service typically gets no more than its allotted
milliseconds of the end-to-end SLA), rather than hoping it works out.

Most real systems are a mix of both shapes (a request fans out breadth-wise to
several services, each of which internally makes a sequential chain of calls to
its own dependencies) — model your load test's synthetic backend delays (e.g.,
via Toxiproxy, §8.2) to match the *actual* shape of your call graph, not a flat
single-hop assumption.

### 3.11 Throughput, RPS, QPS, TPS, and goodput — same-sounding, different things

These terms are used inconsistently across tools and teams; always state units
explicitly in a report rather than assuming "throughput" means the same thing to
every reader.

| Term | What it counts | Typical source |
|---|---|---|
| **RPS** (requests/sec) | Completed HTTP/RPC requests per second, any size | `wrk`, `hey`, `vegeta`, k6 `http_reqs` |
| **QPS** (queries/sec) | Historically DB/search query rate; often used interchangeably with RPS at the HTTP layer, but can mean something narrower (DB queries per second, which may be a multiple of the HTTP request rate if one request triggers several queries) | DB monitoring, search engines |
| **TPS** (transactions/sec) | A "transaction" may bundle multiple requests/steps into one logical unit (e.g., a JMeter/Locust "transaction" = login + browse + checkout counted as 1) | JMeter, Locust, payment systems |
| **Throughput** (network sense) | Bytes/sec or bits/sec moved — a bandwidth measurement, independent of request count | `iperf3` (§7.9), CDN/network dashboards |
| **Throughput** (application sense) | Loosely used to mean RPS — always confirm which sense a dashboard means | varies |
| **Goodput** | Throughput of only *successfully useful* data — excludes retransmissions, failed requests, and protocol overhead | the number that should actually drive capacity planning |

A concrete trap: a service under retry storms (§6.1, §6.3) can show *rising* raw
RPS/throughput at the load balancer while **goodput is flat or falling** — the
extra "throughput" is wasted work from clients retrying failed or timed-out
calls. Always track success-only throughput (goodput) alongside raw request
rate; a load test dashboard that only shows total RPS can make a degrading
system look like it's handling more traffic than it actually is.

Comparing numbers across tools: "500 TPS" from a JMeter test plan where each
transaction contains 5 HTTP samplers is **2500 RPS** at the HTTP layer — not
directly comparable to "500 RPS" reported by `wrk` against a single endpoint.
Always normalize to requests-per-second at a specific, named endpoint before
comparing results across tools or teams.

### 3.12 Concurrency, saturation, and queue depth

- **Concurrency / VUs (virtual users)**: the *closed-loop* knob most tools expose
  (`-c` in `wrk`/`hey`, `--users` in Locust/Goose). Not the same as
  requests-in-flight if virtual users include think-time/sleep between actions —
  a VU that spends 90% of its time sleeping contributes far less to
  requests-in-flight than a VU hammering back-to-back. Little's Law (§3.3) gives
  you the real requests-in-flight number: `L = λ * W`, independent of how many
  VUs are configured.
- **Saturation**: the fraction of a resource's capacity currently in use, and —
  more importantly — whether work is *queueing* for it (USE method, §11.2).
  Saturation is a **leading indicator**: queue depth on a thread pool or DB
  connection pool rises *before* p99 latency visibly degrades, because the first
  queued requests still complete reasonably fast — watch queue depth /
  pool-utilization metrics directly rather than waiting for the latency graph to
  confirm what's already happening.
- **Queue / wait time**: `response_time − service_time` (§3.5) — the portion of
  latency spent waiting rather than being processed. A load test that shows
  rising p99 with *flat* service-time-per-request (visible via APM/tracing) but
  *rising* queue time is diagnosing a concurrency-limit problem (thread pool,
  connection pool, semaphore) rather than a per-request performance problem —
  the fix is capacity/pooling, not code optimization.

### 3.13 Error rate and Apdex

```
error_rate = failed_requests / total_requests
```

Always segment error rate by status-code class (4xx client errors mixed into the
same bucket as 5xx server errors hides real degradation) and watch it
*alongside* latency — a service that starts returning fast 503s under overload
(load shedding working as designed, §5.3) will show a latency *improvement* and
an error-rate *spike* simultaneously; reading either metric alone gives the
wrong story.

**Apdex** (Application Performance Index) condenses latency into a single
satisfaction score using a target threshold `T`:

```
satisfied   = requests with response_time <= T
tolerating  = requests with T < response_time <= 4T
frustrated  = requests with response_time > 4T

Apdex = (satisfied + tolerating/2) / total_requests     // range: 0 (worst) to 1 (best)
```

Apdex trades away percentile granularity for a single trend-able number
executives/dashboards can watch over time — useful as a high-level health
indicator, but always keep the underlying percentiles available for actual
debugging; Apdex alone can't tell you *whether* the slow tail is 2xT or 50xT.

### 3.14 Don't average percentiles — merging results across nodes/shards

A subtle but common mistake: if you run a distributed load test (Locust workers,
k6 with multiple instances, or a fleet of servers each reporting their own p99),
you **cannot** average the per-node p99 values and call the result the global
p99 — percentiles are not linear, so `mean(p99_node1, p99_node2, ..., p99_nodeN)`
is mathematically meaningless and typically *underestimates* the true global
tail.

The correct approach is to merge the underlying **histograms** (not the derived
percentiles) and compute the percentile once from the combined data —
exactly what HdrHistogram is designed to support:

```go
// merge_histograms.go — combine per-node HdrHistogram data into one global view
package main

import (
	"fmt"
	"github.com/HdrHistogram/hdrhistogram-go"
)

func main() {
	// each node records its own latencies independently during the test
	node1 := hdrhistogram.New(1, 60_000, 3) // 1us..60s range, 3 significant figures
	node2 := hdrhistogram.New(1, 60_000, 3)

	for _, v := range []int64{20, 22, 25, 30, 45} {
		node1.RecordValue(v)
	}
	for _, v := range []int64{18, 21, 90, 250, 26} {
		node2.RecordValue(v)
	}

	// WRONG: (node1.ValueAtQuantile(99) + node2.ValueAtQuantile(99)) / 2
	// RIGHT: merge raw histograms, then compute the percentile once
	global := hdrhistogram.New(1, 60_000, 3)
	global.Merge(node1)
	global.Merge(node2)

	fmt.Printf("node1 p99=%d node2 p99=%d (do NOT average these)\n",
		node1.ValueAtQuantile(99), node2.ValueAtQuantile(99))
	fmt.Printf("global p99 (correct, from merged histogram)=%d\n",
		global.ValueAtQuantile(99))
}
```

Locust's distributed mode and k6's cloud/operator output already do this
merge internally before reporting aggregate percentiles — but any custom
multi-node harness (§9) must merge raw histograms itself, or report per-node
percentiles clearly labeled as such rather than a misleading blended average.

### 3.15 Worked example — reading a real report end-to-end

A typical `vegeta report` (§7.3, §10.1) output looks like this:

```
Requests      [total, rate, throughput]  30000, 500.02, 499.87
Duration      [total, attack, wait]      60.016s, 60s, 16.3ms
Latencies     [min, mean, 50, 90, 95, 99, max]  8.1ms, 24.3ms, 19.8ms, 38.1ms, 52.4ms, 187.6ms, 1.204s
Success       [ratio]                    99.62%
Status Codes  [code:count]               200:29886  503:114
```

Reading it with everything covered in this section:

- **Requests rate (500.02) vs. throughput (499.87)**: rate is what you *asked
  for* (open-loop target); throughput is what was actually *achieved*
  (successful completions/sec, §3.11's "goodput" sense) — the tiny gap here
  (0.15 req/s) is healthy; a large gap means the server couldn't keep up with
  the requested rate.
- **mean (24.3ms) vs. p50 (19.8ms)**: mean is pulled above the median (§3.6) —
  expected, and the gap is small, so the tail isn't yet extreme.
- **p99 (187.6ms) vs. max (1.204s)**: a big jump from p99 to max tells you a
  *very* small number of requests (≤1%, potentially just 1-2 out of 30,000) hit
  something pathological (GC pause, lock, cold cache) — worth cross-referencing
  against server-side APM traces at that exact timestamp, not something to
  over-index on if it's truly one outlier.
- **Success ratio (99.62%) and the 503 count (114)**: 114/30000 ≈ 0.38% error
  rate, consistent with the success ratio — if this were a flash-sale test
  (§6.1), 503s here might be *correct* load-shedding behavior rather than a bug;
  check what the 503s correlate with (a specific hot endpoint? a specific time
  window during the ramp?) before treating the error rate alone as pass/fail.
- **What's missing from this report**: jitter, saturation/queue-depth, and
  service-vs-wait-time breakdown — none of which `vegeta` captures on its own.
  Get those from server-side metrics (APM, Prometheus histograms, §11) captured
  during the same time window, and correlate by timestamp — a load tool's
  client-side report is only half the picture (§3.5).

### Quick-reference: metrics glossary

| Metric | One-line definition | Watch out for |
|---|---|---|
| Latency / response time | Time from request sent to response fully received | Client-side vs. server-side measurement differ (§3.5) |
| TTFB / TTLB | Time to first byte / time to last byte | Large gap ⇒ slow streaming or large payload, not slow processing |
| Mean / average | Sum ÷ count | Skewed by outliers; always pair with percentiles (§3.6) |
| Median / p50 | Middle value | Computation method (nearest-rank vs. interpolation) affects the exact number (§3.7) |
| p90 / p95 / p99 / p99.9 | Value below which X% of requests fall | Cannot be averaged across nodes — must merge histograms (§3.14) |
| Standard deviation / CV | Spread around the mean | Dominated by the same tail that skews the mean; CV better for comparing systems (§3.8) |
| Jitter | Variation between consecutive requests | Distinct from percentile spread — two systems can share p50/p99 with very different jitter (§3.9) |
| Fan-out amplification | P(≥1 of N slow) = 1-(1-p)^N | Grows fast with breadth even at low per-call failure rates (§3.10) |
| RPS / QPS / TPS | Requests / queries / transactions per second | Not interchangeable — a "transaction" may bundle multiple requests (§3.11) |
| Throughput / goodput | Bytes or requests per second; goodput = successful-only | Raw throughput can rise during retry storms while goodput falls (§3.11) |
| Concurrency / VUs | Virtual users or in-flight requests | VUs ≠ in-flight if there's think-time; use Little's Law to get the real number (§3.3, §3.12) |
| Saturation / queue depth | How full a resource is / how much work is queued | Leading indicator — rises before latency visibly degrades (§3.12) |
| Error rate | Failed ÷ total | Segment by status-code class; can rise *with* falling latency during load shedding (§3.13) |
| Apdex | Weighted satisfaction score from a threshold T | Single number hides *how* bad the frustrated tail is (§3.13) |

---

## 4. Load & Traffic Modeling

### 4.1 Traffic shapes seen in production

| Pattern | Description | Real example |
|---|---|---|
| Diurnal | Smooth daily sine wave | B2B SaaS, business-hours traffic |
| Flash spike | 10-100x traffic in seconds | Flash sale, celebrity tweet, news event |
| Thundering herd | Many clients retry simultaneously after an outage or cache expiry | Cache stampede, mass token refresh at midnight |
| Sawtooth | Steady growth then hard reset | Batch jobs, cron-triggered fan-out |
| Poisson / random arrivals | Independent, memoryless arrivals | General public web traffic baseline |
| Bursty / self-similar | Traffic with long-range dependence, not smooth Poisson | CDN request logs, video streaming (real internet traffic is provably NOT Poisson — see Leland et al.'s "self-similar nature of Ethernet traffic") |

### 4.2 Modeling arrival rates with Poisson / exponential inter-arrival

For a target rate λ req/s, inter-arrival times between requests should be drawn from
an exponential distribution (`Exp(λ)`) if you want to simulate independent random
users, not a fixed period (fixed period = "closed-loop" bursting in lockstep, which
under-stresses connection-acceptance queues).

```python
import random
def next_interarrival_seconds(rate_per_sec: float) -> float:
    return random.expovariate(rate_per_sec)
```

### 4.3 Ramp patterns for finding breaking points

- **Step load**: hold at 100 req/s for 1 min, then 200, then 400... until errors/
  latency degrade. Cheap and fast; identifies the rough ceiling.
- **Ramp (soak-in)**: linear increase over 10-30 min. Reveals gradual resource
  leaks (connection leaks, memory growth, GC pressure) that step load's short
  plateaus might miss.
- **Soak test**: fixed moderate load for hours. Catches memory leaks, log-disk
  fill-up, certificate/token expiry bugs, slow degradation from fragmentation.
- **Spike test**: baseline load, then a sudden 10x spike for 30-60s, then back to
  baseline. Tests autoscaling reaction time and overload protections (rate
  limiting, circuit breakers, load shedding).
- **Breakpoint / stress test**: increase until the system actually fails, to learn
  *how* it fails (graceful 503s + shed load? Or does it fall over and take
  dependencies with it — e.g., DB connection storm during retries?).

---

## 5. API & Endpoint Performance

### 5.1 Where API latency actually comes from

For a typical CRUD/REST endpoint, budget roughly:

```
Total p50 = 40ms example breakdown
  Load balancer / ingress:        1-2ms
  App framework routing:          <1ms
  Auth/JWT validation:            1-3ms   (RSA verify is much slower than HMAC)
  Business logic:                 2-5ms
  DB query (indexed):             3-10ms
  DB query (unindexed/full scan): 50-500ms+   <- usual culprit
  Serialization (JSON):           1-5ms   (grows with payload size + nesting)
  Network back to client:         RTT/2
```

### 5.2 Payload & serialization costs

- JSON is human readable but CPU-costly to (de)serialize at scale; Protocol
  Buffers / FlatBuffers / MessagePack cut CPU and payload size significantly for
  high-throughput internal APIs.
- Pagination and field selection (GraphQL, sparse fieldsets, `?fields=`) prevent
  "N+1 over-fetch" from blowing up payload size and serialization time under load.
- Compression (gzip/br) trades CPU for bandwidth — beneficial on WAN links, harmful
  on already-CPU-bound backends serving tiny payloads.

### 5.3 Connection-level tuning that dominates results

- **Keep-alive**: reuse TCP+TLS handshakes (§2.1). Verify your load tool is actually
  reusing connections (`wrk` and `k6` do by default over HTTP/1.1 keep-alive/HTTP/2;
  naive scripts using a new `requests.get()` session per call in Python do not).
- **Connection pool sizing** on both client and server: too small → queueing before
  a request even starts; too large → downstream resource exhaustion (DB max
  connections, file descriptor limits).
- **Timeouts**: every hop needs an explicit timeout shorter than its caller's, or a
  single slow dependency stalls the entire call chain (see §3.2). No timeout =
  eventual thread-pool/connection-pool exhaustion under any transient slowdown.
- **Backpressure & load shedding**: rate limiting (token bucket/leaky bucket),
  bulkheads (isolate resource pools per dependency), and circuit breakers prevent
  a single slow endpoint from cascading. Load tests should explicitly verify these
  kick in (you should see clean `429`/`503` responses under overload, not hangs or
  crashes).

### 5.4 gRPC / RPC specific concerns

gRPC multiplexes many calls over one HTTP/2 connection — great for efficiency, but
means a single misbehaving stream (huge message, slow client reader) can create
head-of-line blocking for other RPCs sharing that connection. Load-test gRPC
services with multiple connections per client (`ghz --connections`) to model
realistic multi-client concurrency, not one giant shared connection.

---

## 6. Real-World Scenarios & Case Studies

### 6.1 Flash sale / e-commerce (thundering herd on inventory)

**Scenario**: 200k users refresh a product page in the 10 seconds around a sale's
start time; all attempt to buy 5,000 units of a limited item.

**What breaks in the wild**: optimistic-locking retries storm the DB; cache
stampede on the product-detail cache key when it expires exactly at sale time;
payment-gateway calls queue up behind a DB row lock on the inventory counter.

**How to simulate**:
- Open-loop ramp from 100 → 50,000 virtual users over 15s (spike test, §4.3).
- Model realistic behavior: 95% of users only view the page (read-heavy GET), 5%
  attempt checkout (write, contended resource) — Locust's weighted `@task` or k6
  scenarios express this mix well (see §7).
- Verify: does the inventory count ever go negative (correctness under
  concurrency)? Does p99 for *browsing* degrade because checkout writes starve DB
  connections shared with reads? (Answer: use separate connection pools/read
  replicas.)

### 6.2 Ticket booking (seat/slot contention, hold-and-confirm)

**Scenario**: Concert ticket sale — thousands compete for the same 500 seats,
with a "hold for 5 minutes while paying" step.

**What breaks**: held-but-unpaid seats leak if the release-on-timeout job isn't
tested under load; DB deadlocks from concurrent updates on the same seat-map rows;
WebSocket/polling connections for live seat-map updates fan out to O(users) —
tests must include the read-side broadcast load, not just the write path.

**How to simulate**: multi-stage scenario — `select seat -> hold -> pay (with
injected 2-10s delay to mimic a payment gateway) -> confirm/expire`. This is a
stateful, multi-step scenario best expressed in Locust (Python, stateful `TaskSet`)
or k6 (JS, sequential steps in one VU iteration) rather than a stateless
single-endpoint hammer tool like `wrk`.

### 6.3 Payment gateway (tail latency + idempotency under retry storms)

**Scenario**: a client times out at 3s and retries; if the original request
actually succeeded server-side just after the client gave up, a naive retry
double-charges.

**What to load-test**: correctness under *retried* concurrent load — send the same
idempotency key twice concurrently and assert exactly one charge occurs. This is a
correctness-under-load test, not just a throughput test — write a small
Go/Python harness that fires N duplicate requests concurrently (see §9) and
asserts on the resulting ledger state.

### 6.4 Social feed fan-out (read-heavy, cache-dependent)

**Scenario**: celebrity account posts; millions of followers' feeds must reflect
it, either via fan-out-on-write (push to followers' feed caches) or
fan-out-on-read (compute at read time).

**Load test focus**: fan-out-on-write load-tests the write path's amplification
(1 post → millions of cache writes — test with a "hot user" whose follower count
skews the distribution, not just uniform random users); fan-out-on-read
load-tests read-path aggregation latency as follower/following graphs grow —
percentile latency should be measured *per user tier* (a load test that only
reports a global p99 will hide that celebrity-adjacent reads are 50x slower).

### 6.5 Video streaming (bandwidth + jittery mobile networks)

**Scenario**: adaptive bitrate streaming over mobile networks with variable
bandwidth, latency, and packet loss.

**Load test focus**: this is as much a *network condition* test as a throughput
test. Use `tc netem` (§8) to simulate 3G/4G profiles (high latency, jitter,
1-2% packet loss) while running normal HTTP range-request load against the CDN/
origin, and verify ABR (adaptive bitrate) logic downgrades quality instead of
stalling, and that origin shielding correctly absorbs a thundering herd of
cache-misses when a new episode drops.

### 6.6 IoT telemetry ingestion (connection churn, small payloads, at massive scale)

**Scenario**: millions of devices each send a small payload every 30s, with
frequent connect/disconnect (unlike a browser, IoT devices often don't keep
long-lived connections due to intermittent power/network).

**Load test focus**: here connection *setup* cost (§2.1) dominates over
in-connection throughput. Benchmark connections/sec, not just requests/sec on
warm connections — tools like `wrk` with `--latency` and low keep-alive, or a
custom Go/Rust harness that deliberately opens a fresh connection per request (see
§9), model this far better than default keep-alive-friendly tools.

---

## 7. Open-Source Load Testing Tools — Walkthroughs

Sections 1–6 gave you the vocabulary and the scenario patterns; from here on
it's hands-on. This section is a tool-by-tool walkthrough — pick the one
that matches your target (raw HTTP throughput, distributed scripted
scenarios, or raw TCP/UDP) rather than reading all of it linearly.

### 7.1 `wrk` (C) — high-performance HTTP benchmarking

Written in C on top of an epoll/kqueue event loop with a small Lua scripting
layer. Best for raw HTTP/1.1 throughput/latency measurement with minimal
generator overhead.

```bash
# install (macOS)
brew install wrk

# 30 seconds, 12 threads, 400 concurrent connections, print latency percentiles
wrk -t12 -c400 -d30s --latency http://localhost:8080/api/products
```

Custom Lua script to POST JSON with per-request randomized bodies (models varied
request content, not identical repeated requests):

```lua
-- post.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

request = function()
  local id = math.random(1, 100000)
  local body = string.format('{"user_id": %d, "item": "sku-%d"}', id, id % 500)
  return wrk.format(nil, nil, nil, body)
end
```

```bash
wrk -t8 -c200 -d60s -s post.lua --latency http://localhost:8080/api/orders
```

Note: plain `wrk` is a **closed-loop** tool (§3.4) — for open-loop, constant-rate
testing, use its fork **`wrk2`**, which adds `-R <rate>` and corrects for
"coordinated omission" (a bias where slow responses cause the tool to
under-sample exactly the slow period, making reported percentiles look
artificially good).

### 7.2 `hey` (Go) — simple, single-binary HTTP load generator

```bash
go install github.com/rakyll/hey@latest

hey -z 30s -c 200 -m POST -T "application/json" \
    -d '{"user_id":1,"item":"sku-1"}' \
    http://localhost:8080/api/orders
```

Output includes a latency histogram and status code distribution — good for quick
checks without writing a script.

### 7.3 `vegeta` (Go) — attack/report pipeline, built for open-loop rate testing

`vegeta` is designed around Little's Law-aware, constant-arrival-rate testing and
a separate `report` step, which maps directly onto §10 (report generation).

```bash
go install github.com/tsenart/vegeta@latest

# targets file: METHOD URL pairs
cat <<EOF > targets.txt
GET http://localhost:8080/api/products/1
GET http://localhost:8080/api/products/2
POST http://localhost:8080/api/orders
Content-Type: application/json
@body.json
EOF

# open-loop attack at fixed 500 req/s for 60s, save results
vegeta attack -targets=targets.txt -rate=500 -duration=60s > results.bin

# text report: min/mean/p50/p95/p99/max, success ratio
vegeta report results.bin

# histogram of latency buckets
vegeta report -type='hist[0,10ms,50ms,100ms,200ms,500ms,1s]' results.bin

# plot: interactive HTML latency-over-time chart
vegeta plot results.bin > plot.html
```

### 7.4 `k6` (Go core, JavaScript scripting) — modern scriptable load testing

The most widely adopted modern tool; scenarios, thresholds (pass/fail SLAs baked
into the test), and multiple executor types (constant VUs, ramping VUs,
constant-arrival-rate, ramping-arrival-rate) make it well-suited to every pattern
in §4.3.

```javascript
// flash_sale.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const checkoutLatency = new Trend('checkout_latency');

export const options = {
  scenarios: {
    browse: {
      executor: 'ramping-arrival-rate',
      startRate: 50,
      timeUnit: '1s',
      preAllocatedVUs: 2000,
      maxVUs: 5000,
      stages: [
        { target: 50, duration: '10s' },   // baseline
        { target: 5000, duration: '15s' }, // flash-sale spike
        { target: 5000, duration: '30s' }, // sustained peak
        { target: 0, duration: '10s' },    // drain
      ],
      exec: 'browse',
    },
    checkout: {
      executor: 'ramping-arrival-rate',
      startRate: 2,
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 500,
      stages: [
        { target: 2, duration: '10s' },
        { target: 300, duration: '15s' },
        { target: 300, duration: '30s' },
        { target: 0, duration: '10s' },
      ],
      exec: 'checkout',
    },
  },
  thresholds: {
    'http_req_duration{scenario:browse}': ['p(95)<300'],
    'http_req_duration{scenario:checkout}': ['p(99)<1500'],
    errors: ['rate<0.01'],
  },
};

export function browse() {
  const res = http.get('http://localhost:8080/api/products/42');
  check(res, { 'status is 200': (r) => r.status === 200 });
  errorRate.add(res.status !== 200);
}

export function checkout() {
  const payload = JSON.stringify({ user_id: __VU, item: 'sku-limited-42' });
  const res = http.post('http://localhost:8080/api/orders', payload, {
    headers: { 'Content-Type': 'application/json' },
  });
  checkoutLatency.add(res.timings.duration);
  check(res, { 'checkout ok or sold-out': (r) => [200, 409].includes(r.status) });
  errorRate.add(![200, 409].includes(res.status));
}
```

```bash
k6 run --out json=results.json flash_sale.js
```

k6's `thresholds` block is the key feature for CI gating: the run **exits non-zero**
if p95/p99 SLAs are violated, so this drops directly into a CI/CD pipeline as a
performance regression gate.

### 7.5 Locust (Python) — code-first, stateful user-behavior simulation

Best when scenarios are inherently stateful/sequential (login → browse → hold seat
→ pay, §6.2) and you want to express that logic in a real programming language
rather than a DSL.

```python
# locustfile.py
from locust import HttpUser, task, between, events
import random

class ShopperUser(HttpUser):
    wait_time = between(1, 3)   # closed-loop "think time" between actions

    def on_start(self):
        resp = self.client.post("/api/login", json={"user": f"user{random.randint(1,100000)}"})
        self.token = resp.json().get("token")
        self.client.headers.update({"Authorization": f"Bearer {self.token}"})

    @task(9)
    def browse_product(self):
        product_id = random.randint(1, 5000)
        self.client.get(f"/api/products/{product_id}", name="/api/products/[id]")

    @task(1)
    def checkout(self):
        payload = {"item": "sku-limited-42", "qty": 1}
        with self.client.post("/api/orders", json=payload, catch_response=True) as resp:
            if resp.status_code == 409:
                resp.success()  # sold out is an expected outcome, not a failure
            elif resp.status_code != 200:
                resp.failure(f"unexpected status {resp.status_code}")

@events.quitting.add_listener
def _(environment, **kwargs):
    if environment.stats.total.fail_ratio > 0.01:
        environment.process_exit_code = 1  # CI gate
```

```bash
pip install locust

# headless, distributed-ready, generates HTML + CSV reports (see §10.3)
locust -f locustfile.py --headless \
  --users 5000 --spawn-rate 200 --run-time 5m \
  --host http://localhost:8080 \
  --html report.html --csv results
```

For scale beyond one machine's CPU (Python's GIL limits single-process
concurrency), Locust supports a distributed mode:

```bash
# on the master
locust -f locustfile.py --master

# on N worker machines/containers
locust -f locustfile.py --worker --master-host=<master-ip>
```

### 7.6 Goose (Rust) — Locust-inspired, for CPU-efficient high-concurrency tests

When Python's per-process concurrency ceiling becomes the bottleneck (i.e., your
*load generator* becomes the thing under test instead of the target service),
Goose gives Locust-style ergonomics with Rust/Tokio throughput.

```rust
// Cargo.toml: goose = "0.17"
use goose::prelude::*;

async fn browse_product(user: &mut GooseUser) -> TransactionResult {
    let product_id = rand::random::<u32>() % 5000 + 1;
    let _goose = user.get(&format!("/api/products/{product_id}")).await?;
    Ok(())
}

async fn checkout(user: &mut GooseUser) -> TransactionResult {
    let body = serde_json::json!({ "item": "sku-limited-42", "qty": 1 });
    let _goose = user.post_json("/api/orders", &body).await?;
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), GooseError> {
    GooseAttack::initialize()?
        .register_scenario(
            scenario!("Shopper")
                .register_transaction(transaction!(browse_product).set_weight(9)?)
                .register_transaction(transaction!(checkout).set_weight(1)?),
        )
        .execute()
        .await?;
    Ok(())
}
```

```bash
cargo run --release -- --host http://localhost:8080 \
  --users 20000 --run-time 5m --report-file report.html
```

### 7.7 `ghz` (Go) — gRPC load testing

```bash
go install github.com/bojand/ghz/cmd/ghz@latest

ghz --insecure \
  --proto ./order.proto --call order.OrderService.CreateOrder \
  -d '{"user_id":1,"item":"sku-1"}' \
  -c 200 -n 100000 --connections 20 \
  --format html -o report.html \
  localhost:50051
```

`--connections 20` with `-c 200` models 200 concurrent callers spread over 20
HTTP/2 connections (10 streams/connection) — mirroring §5.4's guidance to avoid
funneling all concurrency through one connection.

### 7.8 Apache JMeter (Java) — GUI + enterprise protocol coverage

Still widely used where teams need a GUI test-plan builder, non-HTTP protocol
support (JDBC, JMS, LDAP, FTP), or existing organizational JMeter test-plan
investment. CLI mode for CI:

```bash
jmeter -n -t flash_sale.jmx -l results.jtl -e -o report_dir/
```

`-e -o` auto-generates an HTML dashboard report from the results file (see §10).

### 7.9 `iperf3` (C) — pure network layer, no application involved

Before blaming your application, prove what the network itself can do. `iperf3`
measures raw TCP/UDP throughput and jitter between two hosts with no HTTP/app
overhead — the baseline every other number in this guide should be compared
against.

```bash
# on the server
iperf3 -s

# on the client: 30s TCP throughput test
iperf3 -c server_ip -t 30

# UDP test with specified bandwidth, reports jitter + packet loss - useful for
# validating the "network" side of a video-streaming or VoIP load test (§6.5)
iperf3 -c server_ip -u -b 50M -t 30
```

If `iperf3` shows you can only push 20 Mbps between two regions, no amount of
API optimization will make a 100 Mbps application-level load test pass — fix
network path/MTU/routing first.

### 7.10 Tool selection cheat-sheet

| Need | Tool |
|---|---|
| Quick single-endpoint smoke test | `hey`, `ab` |
| Raw max-throughput HTTP/1.1 benchmark | `wrk` |
| Open-loop constant-rate testing, avoid coordinated omission | `wrk2`, `vegeta`, k6 arrival-rate executor |
| Complex, stateful, multi-step user journeys | Locust (Python), Goose (Rust), k6 (JS) |
| CI-gated performance regression testing | k6 (`thresholds`), Locust (`process_exit_code`) |
| gRPC services | `ghz` |
| Distributed, massive-scale generation | Locust distributed mode, k6 Cloud/Operator, Goose |
| Raw network baseline (no app layer) | `iperf3` |
| GUI test-plan / enterprise protocols | JMeter |
| Network condition simulation (latency/loss/jitter) | `tc netem`, Toxiproxy, Comcast |

---

## 8. Simulating Network Conditions

Load (concurrency/throughput) is only half the story — real users are on
imperfect networks. Simulating latency, jitter, and packet loss surfaces bugs
that a pristine localhost/datacenter test never will (timeout tuning, retry
storms, TCP behavior under loss).

### 8.1 `tc` + `netem` (Linux) — kernel-level network emulation

```bash
# add 100ms latency +/- 20ms jitter (normal distribution) on eth0 egress
sudo tc qdisc add dev eth0 root netem delay 100ms 20ms distribution normal

# add 2% packet loss on top
sudo tc qdisc change dev eth0 root netem delay 100ms 20ms loss 2%

# simulate a constrained mobile link: 100ms delay, 1% loss, 1mbit cap
sudo tc qdisc add dev eth0 root handle 1: netem delay 100ms loss 1%
sudo tc qdisc add dev eth0 parent 1: handle 10: tbf rate 1mbit burst 32kbit latency 400ms

# remove all shaping when done
sudo tc qdisc del dev eth0 root
```

For containerized targets, apply netem inside the container's network namespace
(or on a dedicated bridge) so you don't degrade the host's other traffic:

```bash
# apply to a specific docker container's interface
CID=$(docker inspect -f '{{.State.Pid}}' my_service_container)
sudo nsenter -t $CID -n tc qdisc add dev eth0 root netem delay 150ms 30ms loss 1%
```

Common real-world profiles to script as reusable presets:

| Profile | Command suffix |
|---|---|
| Good WiFi | `delay 5ms 2ms` |
| 4G | `delay 50ms 20ms loss 0.5%` |
| 3G | `delay 150ms 50ms loss 2% rate 750kbit` |
| Satellite | `delay 600ms 50ms loss 1%` |
| Congested intercontinental link | `delay 200ms 40ms loss 3% rate 5mbit` |

### 8.2 Toxiproxy — application-level fault injection proxy

Unlike `tc` (kernel-wide, needs root, per-interface), Toxiproxy sits as a TCP
proxy in front of a specific dependency (DB, downstream API) and lets you inject/
remove faults at runtime via HTTP API — ideal for automated test suites.

```bash
docker run -d --name toxiproxy -p 8474:8474 -p 5433:5433 ghcr.io/shopify/toxiproxy

# create a proxy: app connects to localhost:5433 instead of the real DB directly
curl -X POST http://localhost:8474/proxies \
  -d '{"name":"postgres","listen":"0.0.0.0:5433","upstream":"postgres_real:5432"}'

# inject 200ms latency with 50ms jitter on all upstream traffic
curl -X POST http://localhost:8474/proxies/postgres/toxics \
  -d '{"type":"latency","attributes":{"latency":200,"jitter":50}}'

# simulate the DB link going completely dark (timeout testing)
curl -X POST http://localhost:8474/proxies/postgres/toxics \
  -d '{"type":"timeout","attributes":{"timeout":0}}'
```

Run your standard load test (k6/Locust/wrk) against the app *while* Toxiproxy
injects DB latency — this directly validates the timeout/circuit-breaker/bulkhead
behavior described in §5.3, under real concurrent load rather than a single
manual request.

### 8.3 Comcast (Go) — simpler netem wrapper, cross-platform-ish CLI

```bash
go install github.com/tylertreat/comcast@latest
comcast --latency=150 --packet-loss=5% --target-bw=1000 --device=eth0
comcast --stop   # remove all rules
```

---

## 9. Building Your Own Load Generators

Off-the-shelf tools cover most needs, but building a minimal generator yourself
teaches you exactly what "requests/sec" costs at the socket level, and is
sometimes necessary for protocol-specific or highly custom scenarios (§6.3's
idempotency race test, for instance, needs custom assertion logic no generic
tool provides).

### 9.1 C — raw TCP client with precise per-request latency (epoll)

Demonstrates connection-per-request cost (§6.6, IoT scenario) using raw sockets,
`TCP_NODELAY` (§2.4), and `clock_gettime(CLOCK_MONOTONIC)` for latency
measurement unaffected by wall-clock adjustments.

```c
// tcp_bench.c — minimal HTTP/1.0 GET load generator, connection-per-request
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/tcp.h>

static double now_ms(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return ts.tv_sec * 1000.0 + ts.tv_nsec / 1e6;
}

static double do_request(const char *ip, int port, const char *path) {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    int one = 1;
    setsockopt(fd, IPPROTO_TCP, TCP_NODELAY, &one, sizeof(one)); // disable Nagle, see 2.4

    struct sockaddr_in addr = {0};
    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    inet_pton(AF_INET, ip, &addr.sin_addr);

    double t0 = now_ms();
    if (connect(fd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        perror("connect");
        close(fd);
        return -1;
    }

    char req[512];
    int n = snprintf(req, sizeof(req),
        "GET %s HTTP/1.0\r\nHost: %s\r\nConnection: close\r\n\r\n", path, ip);
    write(fd, req, n);

    char buf[4096];
    ssize_t r;
    while ((r = read(fd, buf, sizeof(buf))) > 0) { /* drain response */ }
    double elapsed = now_ms() - t0;

    close(fd);
    return elapsed;
}

int main(int argc, char **argv) {
    if (argc < 5) {
        fprintf(stderr, "usage: %s <ip> <port> <path> <num_requests>\n", argv[0]);
        return 1;
    }
    const char *ip = argv[1];
    int port = atoi(argv[2]);
    const char *path = argv[3];
    int total = atoi(argv[4]);

    double *latencies = malloc(sizeof(double) * total);
    int ok = 0;
    double sum = 0, max_lat = 0;

    for (int i = 0; i < total; i++) {
        double lat = do_request(ip, port, path);
        if (lat >= 0) {
            latencies[ok++] = lat;
            sum += lat;
            if (lat > max_lat) max_lat = lat;
        }
    }

    // sort for percentiles
    for (int i = 0; i < ok - 1; i++)
        for (int j = i + 1; j < ok; j++)
            if (latencies[j] < latencies[i]) {
                double tmp = latencies[i]; latencies[i] = latencies[j]; latencies[j] = tmp;
            }

    printf("requests=%d ok=%d\n", total, ok);
    printf("avg=%.2fms p50=%.2fms p95=%.2fms p99=%.2fms max=%.2fms\n",
        sum / ok,
        latencies[(int)(ok * 0.50)],
        latencies[(int)(ok * 0.95)],
        latencies[(int)(ok * 0.99)],
        max_lat);

    free(latencies);
    return 0;
}
```

```bash
gcc -O2 -o tcp_bench tcp_bench.c
./tcp_bench 127.0.0.1 8080 /api/health 1000
```

This is intentionally sequential (one connection at a time) to isolate
**per-connection setup cost**. For concurrent connection-per-request load, fork
worker processes or add an epoll event loop managing many sockets in parallel —
which is exactly what `wrk` does internally.

### 9.2 Python — asyncio open-loop load generator with latency histogram

Models a proper **open-loop** generator (§3.4): requests are scheduled at a fixed
Poisson rate regardless of response time, avoiding coordinated omission.

```python
# async_load.py
import asyncio
import time
import random
import statistics
import aiohttp

TARGET_URL = "http://localhost:8080/api/products/1"
TARGET_RATE_PER_SEC = 200      # open-loop arrival rate
DURATION_SEC = 30

latencies = []
errors = 0

async def fire_request(session: aiohttp.ClientSession):
    global errors
    start = time.perf_counter()
    try:
        async with session.get(TARGET_URL, timeout=aiohttp.ClientTimeout(total=5)) as resp:
            await resp.read()
            elapsed = (time.perf_counter() - start) * 1000
            if resp.status == 200:
                latencies.append(elapsed)
            else:
                errors += 1
    except Exception:
        errors += 1

async def generator(session: aiohttp.ClientSession):
    end_time = time.perf_counter() + DURATION_SEC
    tasks = []
    while time.perf_counter() < end_time:
        # schedule immediately; don't await -> true open-loop, overlapping in-flight requests
        tasks.append(asyncio.create_task(fire_request(session)))
        # Poisson inter-arrival, see 4.2
        await asyncio.sleep(random.expovariate(TARGET_RATE_PER_SEC))
    await asyncio.gather(*tasks)

def percentile(data, p):
    data = sorted(data)
    k = int(len(data) * p)
    return data[min(k, len(data) - 1)]

async def main():
    connector = aiohttp.TCPConnector(limit=0)  # no artificial cap; let the OS/server bound it
    async with aiohttp.ClientSession(connector=connector) as session:
        await generator(session)

    total = len(latencies) + errors
    print(f"requests={total} success={len(latencies)} errors={errors}")
    if latencies:
        print(f"avg={statistics.mean(latencies):.2f}ms "
              f"p50={percentile(latencies,0.50):.2f}ms "
              f"p95={percentile(latencies,0.95):.2f}ms "
              f"p99={percentile(latencies,0.99):.2f}ms "
              f"max={max(latencies):.2f}ms")

if __name__ == "__main__":
    asyncio.run(main())
```

```bash
pip install aiohttp
python async_load.py
```

### 9.3 Go — worker-pool HTTP load generator (closed-loop, concurrency-controlled)

Idiomatic Go concurrency pattern: N goroutines pulling from a shared job channel,
results aggregated on a separate channel — the pattern behind most Go-based load
tools (`hey`, `vegeta`).

```go
// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type result struct {
	latency time.Duration
	err     error
	status  int
}

func worker(client *http.Client, url string, jobs <-chan int, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for range jobs {
		start := time.Now()
		resp, err := client.Get(url)
		lat := time.Since(start)
		if err != nil {
			results <- result{latency: lat, err: err}
			continue
		}
		resp.Body.Close()
		results <- result{latency: lat, status: resp.StatusCode}
	}
}

func main() {
	url := os.Args[1]
	concurrency, _ := strconv.Atoi(os.Args[2])
	totalRequests, _ := strconv.Atoi(os.Args[3])

	// reuse connections: keep-alive transport tuned for high concurrency (see 5.3)
	transport := &http.Transport{
		MaxIdleConns:        concurrency * 2,
		MaxIdleConnsPerHost: concurrency * 2,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	jobs := make(chan int, totalRequests)
	results := make(chan result, totalRequests)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(client, url, jobs, results, &wg)
	}
	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	go func() { wg.Wait(); close(results) }()

	var latencies []float64
	errs, ok := 0, 0
	for r := range results {
		if r.err != nil || r.status >= 500 {
			errs++
		} else {
			ok++
			latencies = append(latencies, r.latency.Seconds()*1000)
		}
	}
	elapsed := time.Since(start)

	sort.Float64s(latencies)
	pct := func(p float64) float64 {
		if len(latencies) == 0 {
			return 0
		}
		idx := int(float64(len(latencies)) * p)
		if idx >= len(latencies) {
			idx = len(latencies) - 1
		}
		return latencies[idx]
	}

	fmt.Printf("total=%d ok=%d errors=%d duration=%s throughput=%.1f req/s\n",
		totalRequests, ok, errs, elapsed, float64(ok)/elapsed.Seconds())
	fmt.Printf("p50=%.2fms p95=%.2fms p99=%.2fms\n", pct(0.50), pct(0.95), pct(0.99))
}
```

```bash
go run main.go http://localhost:8080/api/products/1 200 50000
# args: url, concurrency, total requests
```

### 9.4 Rust — tokio async generator with `hdrhistogram` for accurate percentiles

Rust's async ecosystem (`tokio` + `reqwest`) plus `hdrhistogram` (the same
technique HdrHistogram/wrk2 use internally) gives accurate, low-overhead
percentile tracking without keeping every raw sample in memory.

```rust
// Cargo.toml
// [dependencies]
// tokio = { version = "1", features = ["full"] }
// reqwest = { version = "0.12", features = ["rustls-tls"] }
// hdrhistogram = "7"
// futures = "0.3"

use hdrhistogram::Histogram;
use std::sync::{Arc, Mutex};
use std::time::Instant;
use tokio::sync::Semaphore;

#[tokio::main]
async fn main() {
    let url = std::env::args().nth(1).unwrap();
    let concurrency: usize = std::env::args().nth(2).unwrap().parse().unwrap();
    let total: usize = std::env::args().nth(3).unwrap().parse().unwrap();

    let client = reqwest::Client::builder()
        .pool_max_idle_per_host(concurrency)
        .build()
        .unwrap();

    let sem = Arc::new(Semaphore::new(concurrency));   // bound in-flight requests
    let hist = Arc::new(Mutex::new(Histogram::<u64>::new(3).unwrap())); // 3 sig figs
    let errors = Arc::new(std::sync::atomic::AtomicUsize::new(0));

    let start = Instant::now();
    let mut handles = Vec::with_capacity(total);

    for _ in 0..total {
        let permit = sem.clone().acquire_owned().await.unwrap();
        let client = client.clone();
        let url = url.clone();
        let hist = hist.clone();
        let errors = errors.clone();

        handles.push(tokio::spawn(async move {
            let t0 = Instant::now();
            match client.get(&url).send().await {
                Ok(resp) if resp.status().is_success() => {
                    let _ = resp.bytes().await;
                    let micros = t0.elapsed().as_micros() as u64;
                    hist.lock().unwrap().record(micros).unwrap();
                }
                _ => {
                    errors.fetch_add(1, std::sync::atomic::Ordering::Relaxed);
                }
            }
            drop(permit);
        }));
    }

    for h in handles {
        let _ = h.await;
    }

    let elapsed = start.elapsed();
    let h = hist.lock().unwrap();
    let to_ms = |v: u64| v as f64 / 1000.0;

    println!(
        "total={} errors={} duration={:.2}s throughput={:.1} req/s",
        total,
        errors.load(std::sync::atomic::Ordering::Relaxed),
        elapsed.as_secs_f64(),
        h.len() as f64 / elapsed.as_secs_f64()
    );
    println!(
        "p50={:.2}ms p95={:.2}ms p99={:.2}ms p99.9={:.2}ms max={:.2}ms",
        to_ms(h.value_at_quantile(0.50)),
        to_ms(h.value_at_quantile(0.95)),
        to_ms(h.value_at_quantile(0.99)),
        to_ms(h.value_at_quantile(0.999)),
        to_ms(h.max())
    );
}
```

```bash
cargo run --release -- http://localhost:8080/api/products/1 500 100000
```

The `Semaphore`-bounded concurrency model here is deliberately **closed-loop**
(§3.4) — it caps in-flight requests at `concurrency`. To convert this into an
open-loop generator, replace the semaphore-gated spawn loop with a
`tokio::time::interval` ticking at the target rate, spawning a task per tick
without waiting on a permit (mirroring the Python asyncio example in §9.2).

---

## 10. Report Generation & Visualization

### 10.1 `vegeta` — text, histogram, and interactive HTML plots (built-in)

```bash
vegeta attack -targets=targets.txt -rate=500 -duration=60s | tee results.bin | vegeta report

vegeta report -type=json results.bin > report.json
vegeta plot results.bin > plot.html    # self-contained interactive latency-over-time HTML
```

### 10.2 k6 — JSON/InfluxDB output + Grafana dashboards, or a static HTML summary

For a single static report:

```bash
k6 run --out json=results.json flash_sale.js
```

```javascript
// add to flash_sale.js for an automatic HTML summary at the end of the run
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export function handleSummary(data) {
  return { "summary.html": htmlReport(data) };
}
```

For live dashboards during long-running tests, stream metrics to Prometheus/
InfluxDB and visualize in Grafana:

```bash
k6 run --out experimental-prometheus-rw flash_sale.js
# then point a Grafana dashboard (k6's official "k6 Prometheus" dashboard ID: 19665)
# at the same Prometheus instance for real-time percentile graphs during the run
```

### 10.3 Locust — built-in HTML + CSV, or programmatic access

```bash
locust -f locustfile.py --headless --users 5000 --spawn-rate 200 --run-time 5m \
  --host http://localhost:8080 --html report.html \
  --csv results   # produces results_stats.csv, results_stats_history.csv, results_failures.csv
```

`results_stats_history.csv` has a time-series row per reporting interval — feed
this into the Python report generator below for custom charts Locust's built-in
HTML doesn't provide (e.g., overlaying multiple test runs).

### 10.4 JMeter — HTML Dashboard Report

```bash
jmeter -n -t flash_sale.jmx -l results.jtl -e -o report_dir/
# report_dir/index.html: response-time percentiles over time, throughput,
# error rate, and APDEX score out of the box
```

### 10.5 Custom Python report generator — merge multiple tools' raw output into one chart

Useful when comparing runs across tools (e.g., a baseline `wrk` run vs. a k6 run
with `tc netem` network shaping applied), or when you need a report format your
tooling doesn't natively produce (e.g., attaching a chart to a PR comment).

```python
# generate_report.py
import json
import csv
import matplotlib.pyplot as plt

def load_vegeta_json(path):
    with open(path) as f:
        data = json.load(f)
    return {
        "p50": data["latencies"]["50th"] / 1e6,   # ns -> ms
        "p95": data["latencies"]["95th"] / 1e6,
        "p99": data["latencies"]["99th"] / 1e6,
        "throughput": data["throughput"],
        "success_ratio": data["success"],
    }

def load_locust_history_csv(path):
    timestamps, p95s, rps = [], [], []
    with open(path) as f:
        for row in csv.DictReader(f):
            timestamps.append(float(row["Timestamp"]))
            p95s.append(float(row["95%"]))
            rps.append(float(row["Requests/s"]))
    return timestamps, p95s, rps

def render_comparison_chart(runs: dict, out_path="report.png"):
    """runs: {"label": {"p50":.., "p95":.., "p99":.., "throughput":..}, ...}"""
    labels = list(runs.keys())
    p50s = [runs[l]["p50"] for l in labels]
    p95s = [runs[l]["p95"] for l in labels]
    p99s = [runs[l]["p99"] for l in labels]

    x = range(len(labels))
    width = 0.25
    fig, ax1 = plt.subplots(figsize=(10, 5))

    ax1.bar([i - width for i in x], p50s, width, label="p50")
    ax1.bar(x, p95s, width, label="p95")
    ax1.bar([i + width for i in x], p99s, width, label="p99")
    ax1.set_xticks(list(x))
    ax1.set_xticklabels(labels, rotation=20, ha="right")
    ax1.set_ylabel("Latency (ms)")
    ax1.set_title("Latency percentiles by test run")
    ax1.legend()

    plt.tight_layout()
    plt.savefig(out_path, dpi=150)
    print(f"wrote {out_path}")

if __name__ == "__main__":
    runs = {
        "baseline (no shaping)": load_vegeta_json("baseline.json"),
        "with tc netem 4G profile": load_vegeta_json("throttled.json"),
    }
    render_comparison_chart(runs)
```

```bash
pip install matplotlib
python generate_report.py
```

### 10.6 What a good report includes, regardless of tool

1. **Test parameters**: tool, version, target rate/concurrency, duration, ramp
   shape, environment (staging vs prod-like), network conditions applied.
2. **Headline SLI numbers**: p50/p95/p99/p99.9, error rate, throughput achieved
   vs. target — and whether thresholds/SLAs passed or failed (pass/fail, not just
   numbers).
3. **Time-series graphs**, not just aggregate numbers — a flat p99 line across
   the whole run reads very differently from one that spikes only during the
   ramp phase.
4. **Resource utilization overlay**: CPU/memory/DB connections/GC pauses on the
   server during the same window (Grafana dashboard alongside the load tool's
   output) — without this, you can see *that* latency degraded but not *why*.
5. **Comparison to the previous baseline run** — a single run's numbers are much
   less useful than a trend; store historical JSON/CSV output in version control
   or a metrics store to diff against.

---

## 11. Metrics, SLIs/SLOs, and the RED/USE Methods

The tools in §7–10 produce numbers; this section is how those numbers turn
into an ongoing practice instead of a one-off report. RED and USE are the
two standard lenses for deciding *what to measure continuously*, not just
what to measure during a single load test run.

### 11.1 RED method (for request-driven services)

- **R**ate — requests per second
- **E**rrors — failed requests per second (and error rate %)
- **D**uration — latency distribution (percentiles, §3.1)

### 11.2 USE method (for resources: CPU, DB connections, queues, disks)

- **U**tilization — % time the resource was busy
- **S**aturation — amount of queued work the resource couldn't service immediately
- **E**rrors — count of error events for that resource

Cross-reference RED (what users feel) against USE (why) during every load test:
a RED-method latency spike should always correlate with a USE-method saturation
signal somewhere (DB connection pool at 100%, CPU pegged, run queue growing) — if
it doesn't, you likely have a lock contention or external-dependency issue that
plain resource metrics won't show directly.

### 11.3 Defining SLOs from load test data

```
SLI: 99th percentile latency of GET /api/products/{id}
SLO: p99 < 300ms, measured over rolling 5 minutes, 99.9% of the time (error budget: 0.1%)
```

Bake the SLO directly into the load test as a pass/fail gate (k6 `thresholds`,
§7.4) so performance regressions fail CI the same way a broken unit test does,
rather than being discovered after a production incident.

---

## 12. Checklists & Common Pitfalls

**Before running any load test**
- [ ] Confirm you're not accidentally load-testing a shared staging environment
  other teams depend on, or (never) production without explicit sign-off.
- [ ] Warm up caches/JIT/connection pools deliberately, or explicitly include
  cold-start behavior as part of what you're measuring — decide which, don't
  let it happen by accident.
- [ ] Verify the load generator's own host isn't the bottleneck (check its CPU,
  socket/file-descriptor limits — `ulimit -n`, and outbound ephemeral port
  exhaustion under very high connection-per-request rates).
- [ ] Decide open-loop vs. closed-loop deliberately based on what you're
  modeling (§3.4) — don't default to whatever the tool ships with.
- [ ] Synchronize clocks (NTP) between load generator and server if correlating
  timestamps across systems.

**Common pitfalls**
- **Coordinated omission**: closed-loop tools under-sample exactly the slow
  periods, making reported percentiles look better than reality. Use `wrk2`,
  `vegeta`, or k6's arrival-rate executors for anything you'll make a capacity
  decision from.
- **DNS caching skew**: if the load generator resolves DNS once and the real
  clients don't (or vice versa), the two won't see comparable connection-setup
  costs.
- **Same-machine test skew**: running the load generator on the same host as
  the server steals CPU from the very thing you're measuring, and localhost
  network stack behaves nothing like a real network path.
- **Ignoring warm-up in results**: including the first 10-30s of a ramp (JIT
  warm-up, connection-pool fill, autoscaler lag) in your percentile calculation
  skews everything pessimistic in a way that isn't steady-state reality.
- **Testing only the happy path**: also fire malformed/oversized payloads,
  expired auth tokens, and duplicate idempotency keys under load — correctness
  under concurrency (§6.3) is as important as raw speed.
- **One global percentile hides per-tier behavior**: segment latency by
  endpoint, user tier, and payload size (§6.4) — a global p99 can look fine while
  a specific critical path is badly broken.
- **Not testing failure/recovery, only steady load**: spike tests and dependency
  fault-injection (§8.2) reveal how the system behaves *while breaking* and
  *while recovering* — both matter more operationally than steady-state numbers.

---

## 13. Further Reading

- Little's Law — Little, J.D.C. (1961), foundational queueing-theory result.
- "Self-Similarity in Ethernet Traffic" — Leland, Taqqu, Willinger, Wilson (1993)
  — why real traffic isn't Poisson.
- Gil Tene, "How NOT to Measure Latency" (talk) — the definitive explanation of
  coordinated omission, and the motivation behind HdrHistogram/wrk2.
- Google SRE Book, chapters on SLIs/SLOs and the four golden signals (latency,
  traffic, errors, saturation — the request-driven analogue of RED/USE here).
- `wrk2` — https://github.com/giltene/wrk2
- `k6` docs — https://k6.io/docs/
- `Locust` docs — https://docs.locust.io/
- `Goose` docs — https://book.goose.rs/
- `vegeta` — https://github.com/tsenart/vegeta
- `Toxiproxy` — https://github.com/Shopify/toxiproxy
- Linux `tc-netem(8)` man page.
