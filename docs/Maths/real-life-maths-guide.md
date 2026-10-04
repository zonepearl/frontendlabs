# Real-Life Mathematics — A Practical Guide for Software / AI / Network / Security Engineers

> 📖 **Read this on [frontendlabs.xyz](https://frontendlabs.xyz/maths/)**: the official edition, with one page per chapter, search, and dark mode. <!-- frontendlabs-notice -->

> From "why does `0.1 + 0.2 != 0.3`" to deriving backpropagation, breaking down
> why RSA-2048 is safe, and computing how many /24s fit in a /16.
>
> **You do not need a math degree to start.** You need arithmetic, curiosity,
> and a terminal with Python. Every concept is tied to a real engineering
> decision you will actually make — not abstract theory for its own sake.

---

## What you will be able to do at the end

1. Read and reason about **Big-O**, entropy, and probability without hand-waving.
2. Explain **why floating point breaks money math**, and how to avoid it.
3. Compute **subnet ranges, bandwidth-delay products, and queue depths** by hand.
4. Derive **gradient descent and backpropagation** from the chain rule, not from memory,
   and explain why attention divides by `√d_k` and why `η < 2/λ_max`.
5. Read a formula like an API signature: name every symbol, check the units, and
   test it on simple numbers and extreme values.
6. Explain **why RSA and Diffie-Hellman are secure**, using the actual number theory.
7. Quantify **password strength, hash collisions, and detection system false-positive rates**.
8. Know which branch of math to reach for when a new problem shows up, and where
   to go deeper.

---

## How to use this guide

Each of the 29 chapters follows the same shape:

```
## N. Topic Name
**Level:** Beginner / Intermediate / Expert
**You'll use this for:** <the 1-line reason a senior engineer needs it>

- The idea (plain language)
- The formula / mechanism, with a table naming every symbol
- Worked example (by hand, with real numbers), directly under the formula
- The theorems behind it, and how AI and networking systems use it
- Code (Python, runnable)
- Real-life engineering ties (SWE / AI / Network / Security)
- Common mistake
```

Read top to bottom the first time. After that, use it as a reference — jump
straight to the chapter you need. Chapters build on each other roughly in
order, but each is self-contained enough to read alone.

**Tools you need:** Python 3 with the small scientific Python stack below. A
terminal. Every worked example below can be pasted into a Python REPL.

```bash
python3 -m venv venv && source venv/bin/activate
pip install numpy matplotlib scipy sympy scikit-learn pandas
python3
```

---

## Contents

**Part 0 — Start here**
- How to use this guide · Role-based learning paths · Notation cheat sheet · How to read a formula

**Part 1 — Beginner: the foundations everyone needs**
1. Number systems & bit-level math
2. Modular arithmetic
3. Boolean algebra & logic
4. Sets, relations, and functions
5. Exponents & logarithms
6. Sequences, series & growth rates
7. Geometry & trigonometry essentials
8. Basic probability (axioms, conditional probability, expectation)
9. Basic statistics (correlation, EWMA and TCP's retransmission timer)

**Part 2 — Intermediate: core engineering mathematics**
10. Linear algebra I — vectors & matrices (norms, neural layers, routing matrices, least squares)
11. Linear algebra II — eigenvalues, eigenvectors, determinants (PageRank, Markov chains, graph Laplacians)
12. Calculus I — limits, derivatives, integrals (Taylor, Newton's method, the Fundamental Theorem)
13. Calculus II — gradients, Jacobians, Hessians, the chain rule (backprop by hand, softmax gradient)
14. Probability distributions (CLT, LLN, concentration bounds, heavy tails, LLM sampling)
15. Combinatorics & counting
16. Bayes' theorem
17. Graph theory
18. Complexity & asymptotic analysis
19. Number theory basics

**Part 3 — Expert: domain-specific deep dives (senior level)**
20. Linear algebra for AI/ML (SVD, PCA, tensors, attention)
21. Optimization & backpropagation (learning-rate limits, Adam, Lagrange and TCP fairness)
22. Statistics for ML & experimentation (MLE, z-tests, power analysis)
23. Information theory (entropy, KL divergence, channel capacity)
24. Cryptographic mathematics (RSA, Diffie-Hellman, elliptic curves)
25. Network & systems mathematics (queueing, subnetting, congestion control)
26. Quantitative security (entropy, detection theory, risk)
27. Signal processing basics (Fourier, sampling)
28. Game theory & adversarial thinking
29. Control theory basics (feedback loops, PID)

**Part 4 — Putting it together**
- Role-based priority map · Formula cheat sheet · Practice labs

---

## Role-based learning paths

You don't need all 29 chapters at equal depth. Use this as a priority map —
**P1 = must know cold, P2 = should know, P3 = good to know**.

| Chapter | Software Eng. | AI Engineer | Network Eng. | Security Eng. |
|---|---|---|---|---|
| 1. Number systems | P1 | P2 | P1 | P1 |
| 2. Modular arithmetic | P2 | P2 | P2 | P1 |
| 3. Boolean algebra | P1 | P2 | P1 | P1 |
| 4. Sets & functions | P1 | P2 | P2 | P2 |
| 5. Exponents & logs | P1 | P1 | P2 | P1 |
| 6. Sequences & growth | P2 | P2 | P2 | P2 |
| 7. Geometry & trig | P3 | P2 | P2 | P3 |
| 8. Basic probability | P2 | P1 | P2 | P1 |
| 9. Basic statistics | P1 | P1 | P2 | P1 |
| 10–11. Linear algebra | P2 | P1 | P3 | P2 |
| 12–13. Calculus | P2 | P1 | P3 | P2 |
| 14. Distributions | P2 | P1 | P2 | P1 |
| 15. Combinatorics | P2 | P2 | P2 | P1 |
| 16. Bayes' theorem | P2 | P1 | P2 | P1 |
| 17. Graph theory | P1 | P2 | P1 | P2 |
| 18. Complexity | P1 | P1 | P2 | P2 |
| 19. Number theory | P2 | P2 | P2 | P1 |
| 20–21. ML math | P3 | P1 | P3 | P3 |
| 22. Stats for ML/experiments | P2 | P1 | P2 | P2 |
| 23. Information theory | P2 | P1 | P2 | P1 |
| 24. Crypto math | P2 | P3 | P2 | P1 |
| 25. Network/systems math | P2 | P2 | P1 | P2 |
| 26. Quantitative security | P2 | P2 | P2 | P1 |
| 27. Signal processing | P3 | P2 | P1 | P3 |
| 28. Game theory | P3 | P2 | P3 | P2 |
| 29. Control theory | P2 | P2 | P1 | P3 |

## Notation cheat sheet

| Symbol | Meaning | Symbol | Meaning |
|---|---|---|---|
| `∑` | sum | `∏` | product |
| `∈` | "is an element of" | `⊆` | "is a subset of" |
| `∀` | "for all" | `∃` | "there exists" |
| `log₂ x` | log base 2 (bits) | `ln x` | natural log (base *e*) |
| `x mod n` | remainder of x ÷ n | `⌊x⌋` | floor (round down) |
| `O(f(n))` | upper-bound growth rate | `≈` | approximately equal |
| `P(A\|B)` | probability of A given B | `E[X]` | expected value of X |
| `∇f` | gradient of f | `∂f/∂x` | partial derivative |
| `⊕` | XOR | `≡` | congruent (modular equality) |
| `ℝⁿ` | the space of lists of `n` real numbers | `x ∈ ℝⁿ` | "x is a vector with n real entries" |
| `‖x‖` | length (norm) of vector `x` | `Aᵀ` | transpose of `A` (rows ↔ columns) |
| `A⁻¹` | inverse of matrix `A` | `det(A)` | determinant of `A` |
| `I` | identity matrix (1s on the diagonal) | `⊙` | element-wise product |
| `∫ₐᵇ f(x) dx` | integral: area under `f` from `a` to `b` | `lim(h→0)` | the value approached as `h` shrinks to 0 |
| `X ~ N(μ, σ²)` | "X is drawn from a normal distribution with mean μ, variance σ²" | `∝` | "proportional to" |
| `argmax_θ f(θ)` | the `θ` that makes `f` largest (not the largest value itself) | `x̄` | sample mean of `x` |
| `Var(X)` | variance of X | `Cov(X,Y)` | covariance of X and Y |

The Greek letters you will meet over and over, and what they usually stand for:

| Letter | Name | Usually means | Letter | Name | Usually means |
|---|---|---|---|---|---|
| `α` | alpha | significance level, smoothing factor | `β` | beta | momentum/decay factor, Type II error |
| `δ` `Δ` | delta | a small change; `Δx` = "change in x" | `ε` | epsilon | a tiny tolerance or error |
| `η` | eta | learning rate | `θ` | theta | model parameters, or an angle |
| `λ` | lambda | eigenvalue, arrival rate, Lagrange multiplier | `μ` | mu | mean, or service rate in queueing |
| `σ` | sigma | standard deviation, or the sigmoid function | `Σ` | capital sigma | sum, or a covariance matrix |
| `ρ` | rho | utilization (queueing), correlation | `π` | pi | 3.14159…, or a stationary distribution |
| `φ` | phi | Euler's totient, or the normal CDF | `ω` | omega | angular frequency |

Context decides which meaning applies. `λ` is an arrival rate in Ch. 25 and an
eigenvalue in Ch. 11, and each chapter says which one it means.

## How to read a formula (the skill nobody teaches)

Most people bounce off mathematics because they try to read a formula the way
they read prose, left to right, all at once. Read it the way you read an
unfamiliar function signature instead:

1. **Find the output.** What is on the left of the `=`? That is the return value.
2. **Name every symbol.** Each letter is a parameter. Which are inputs you supply,
   which are constants, and which are dummy loop variables (like the `i` in `Σ`)?
3. **Check the units.** If the left side is "seconds", the right side must
   come out in seconds too. Unit-checking catches most misread formulas.
4. **Plug in the simplest numbers.** Set things to 0, 1 or 2 and see if the
   output makes sense.
5. **Push it to the extremes.** What happens as an input goes to 0 or to infinity?
   Extremes show what a formula is really about.

Here is the method on a formula from Ch. 25, the M/M/1 queue's average time in system:

```
W = 1 / (μ - λ)
```

| Symbol | Say it as | What it means | Units |
|---|---|---|---|
| `W` | "W" | average time a request spends waiting plus being served | seconds |
| `μ` | "mu" | how many requests the server *can* finish per second | requests/s |
| `λ` | "lambda" | how many requests *arrive* per second | requests/s |

- **Units:** `1 / (requests/s)` = seconds. ✓
- **Simple numbers:** `μ = 100`, `λ = 50` → `W = 1/50 s = 20 ms`.
- **Extremes:** as `λ → 0` (idle server), `W → 1/μ = 10 ms`, just the service
  time, which makes sense. As `λ → μ` (server at 100% utilization), the denominator
  goes to 0 and `W → ∞`. The formula is telling you that **queues explode
  near full utilization**, which is the whole lesson of Ch. 25 in one line.

Every formula in this guide is laid out the same way: the formula, a table
that names every symbol, then a worked example with real numbers right below it.
When a formula looks scary, go back to these five steps.

---

# Part 1 — Beginner: The Foundations Everyone Needs

Nine chapters, each following the same shape described above, needing
nothing beyond arithmetic to start. This is the foundation everything
later — linear algebra, calculus, cryptography, ML math — quietly assumes
you already have, whether or not anyone ever said so out loud.

## 1. Number Systems & Bit-Level Math

**Level:** Beginner
**You'll use this for:** many bugs involving memory, money, overflow, or "why is this number wrong."

### The idea

Computers store everything as binary (base 2). Humans think in decimal (base
10). Engineers read/write hex (base 16) because it maps cleanly to binary (1
hex digit = 4 bits) and is far more compact than binary for humans.

### Conversions

| Decimal | Binary | Hex |
|---|---|---|
| 10 | 1010 | 0xA |
| 255 | 11111111 | 0xFF |
| 256 | 100000000 | 0x100 |
| 4096 | 1000000000000 | 0x1000 |

```python
>>> bin(255)          # '0b11111111'
>>> hex(255)           # '0xff'
>>> int('ff', 16)      # 255
>>> int('11111111', 2) # 255
```

The table above shows the *answers*. Here's the actual mechanism behind
each arrow — the by-hand process, so a conversion is never a black box.

#### Positional notation: the one idea every conversion rests on

Every number system here is **positional** — each digit's value depends on
*where* it sits, not just what digit it is. In decimal, `243` means
`2×10² + 4×10¹ + 3×10⁰`. The base (10) tells you what each position is a
power of. Binary and hex work identically, just with base 2 and base 16
instead of base 10:

```
binary  1010  =  1×2³ + 0×2² + 1×2¹ + 0×2⁰  =  8 + 0 + 2 + 0  =  10
hex       FF  =  15×16¹ + 15×16⁰            =  240 + 15       =  255
```

Every conversion method below is just a different way of applying (decimal
→ other base) or undoing (other base → decimal) this same place-value idea.

#### Decimal → Binary: repeated division by 2

Divide by 2 repeatedly, writing down the **remainder** (always 0 or 1) each
time, until the quotient hits 0. Then read the remainders **bottom to
top** — that's your binary number.

```
Converting 10 to binary:
  10 ÷ 2 = 5   remainder 0     <- read bottom-up: 1 0 1 0
   5 ÷ 2 = 2   remainder 1
   2 ÷ 2 = 1   remainder 0
   1 ÷ 2 = 0   remainder 1     <- stop, quotient is 0

Reading remainders bottom-to-top: 1010  ✓ matches the table
```

**Why bottom-to-top?** Each division peels off the *least significant* bit
first (the 2⁰ place), the same way repeatedly dividing by 10 peels off a
decimal number's last digit first. The last remainder you compute is the
most significant bit, so it goes on the left.

```python
def decimal_to_binary(n):
    if n == 0:
        return "0"
    bits = []
    while n > 0:
        bits.append(str(n % 2))   # the remainder
        n //= 2                    # the quotient, fed into the next round
    return "".join(reversed(bits))  # remainders came out least-significant-first

print(decimal_to_binary(10))    # '1010'
print(decimal_to_binary(255))   # '11111111'
```

#### Binary → Decimal: multiply each digit by its place value, then add

Reverse of the above: number the bit positions from the right starting at
0, raise 2 to each position, multiply by the bit (0 or 1), and sum.

```
Converting 1010 to decimal:
  position:   3  2  1  0
  bit:        1  0  1  0
  value:   1×2³ + 0×2² + 1×2¹ + 0×2⁰
         =    8 +    0 +    2 +    0
         =   10  ✓
```

```python
def binary_to_decimal(bits):
    total = 0
    for position, bit in enumerate(reversed(bits)):
        total += int(bit) * (2 ** position)
    return total

print(binary_to_decimal("1010"))      # 10
print(binary_to_decimal("11111111"))  # 255
```

#### Decimal → Hex: repeated division by 16

Identical process to decimal → binary, just dividing by 16 instead of 2.
The catch: a remainder can now be 0–15, and hex only has single-character
digits `0-9`, so remainders 10–15 are written as letters `A-F`.

```
Converting 255 to hex:
  255 ÷ 16 = 15   remainder 15 -> 'F'    <- read bottom-up: F F
   15 ÷ 16 =  0   remainder 15 -> 'F'    <- stop, quotient is 0

Reading remainders bottom-to-top: FF  ✓ matches the table
```

```python
HEX_DIGITS = "0123456789ABCDEF"

def decimal_to_hex(n):
    if n == 0:
        return "0"
    digits = []
    while n > 0:
        digits.append(HEX_DIGITS[n % 16])   # remainder, mapped to a hex character
        n //= 16
    return "".join(reversed(digits))

print(decimal_to_hex(255))    # 'FF'
print(decimal_to_hex(4096))   # '1000'
```

#### Hex → Decimal: multiply each digit by its place value, then add

Same idea as binary → decimal, but each position is a power of 16 instead
of 2, and letter digits convert back to their 10–15 value first.

```
Converting FF to decimal:
  position:      1    0
  digit:         F    F   ->  15   15
  value:   15×16¹ + 15×16⁰
         =    240 +    15
         =    255  ✓
```

```python
def hex_to_decimal(hex_str):
    total = 0
    for position, char in enumerate(reversed(hex_str.upper())):
        total += HEX_DIGITS.index(char) * (16 ** position)
    return total

print(hex_to_decimal("FF"))     # 255
print(hex_to_decimal("1000"))   # 4096
```

#### Binary ↔ Hex: group into nibbles (no arithmetic needed at all)

**Why hex maps to binary so cleanly:** 16 = 2⁴, so every hex digit is
*exactly* 4 bits (a "nibble"), and every possible 4-bit pattern (`0000`
through `1111`) has exactly one hex digit (`0` through `F`). That means
converting between binary and hex is pure **regrouping** — no division, no
multiplication, just splitting into chunks of 4 and looking each chunk up:

```
Binary -> Hex: split into groups of 4 bits, starting from the right,
                then convert each group independently.

  11111111  ->  1111 1111  ->  F    F   ->  FF

Hex -> Binary: expand each hex digit into its 4-bit pattern, and
                concatenate.

  FF        ->  F     F    ->  1111 1111  ->  11111111
```

```
nibble lookup table (memorize this once, use it forever):
0000=0   0100=4   1000=8   1100=C
0001=1   0101=5   1001=9   1101=D
0010=2   0110=6   1010=A   1110=E
0011=3   0111=7   1011=B   1111=F
```

```python
def binary_to_hex_by_grouping(bits):
    bits = bits.zfill((len(bits) + 3) // 4 * 4)   # pad to a multiple of 4 on the left
    nibbles = [bits[i:i+4] for i in range(0, len(bits), 4)]
    return "".join(HEX_DIGITS[int(nibble, 2)] for nibble in nibbles)

print(binary_to_hex_by_grouping("11111111"))   # 'FF'
```

This is why MAC addresses (`00:1A:2B:3C:4D:5E`), IPv6 addresses, and memory
dumps are always shown in hex — it's binary in disguise, without 32
characters of 1s and 0s to read, and you can convert either direction in
your head by regrouping 4 bits at a time instead of doing any real math.

### Two's complement (how negative numbers actually work)

An 8-bit signed integer doesn't store a "minus sign" — it flips all bits of
the positive value and adds 1.

```
 5  = 00000101
-5  = 11111011   (flip bits: 11111010, add 1: 11111011)
```

**Why it matters:** signed integer overflow. An 8-bit signed int ranges from
-128 to 127. `127 + 1` wraps to `-128`. This is a *real, recurring bug class*:

- The 1996 Ariane 5 rocket exploded 37 seconds after launch because a 64-bit
  float was cast into a 16-bit signed integer and overflowed.
- Classic game bugs ("civilization gandhi nuke" isn't real, but integer
  wraparound score/health bugs are extremely common).
- **Security angle:** integer overflow is a documented vulnerability class
  (CWE-190). An attacker who can push a size/length calculation past
  `INT_MAX` can cause a buffer to be allocated smaller than expected,
  leading to a buffer overflow.

```python
# Simulate 8-bit signed overflow
def to_int8(n):
    n = n & 0xFF                 # keep lowest 8 bits
    return n - 256 if n >= 128 else n

print(to_int8(127 + 1))   # -128  <- wraps around
```

### IEEE-754 floating point (why `0.1 + 0.2 != 0.3`)

Floats store numbers as `sign × mantissa × 2^exponent`. Just like `1/3`
cannot be written exactly in decimal, `0.1` cannot be written exactly in
binary — it's a repeating fraction. The stored value is the *closest
representable approximation*.

```python
>>> 0.1 + 0.2
0.30000000000000004
>>> 0.1 + 0.2 == 0.3
False
```

**Real-life engineering ties:**
- **SWE:** never store money as `float`/`double`. Use integer cents, or a
  `Decimal` type (`from decimal import Decimal`). Comparing floats for
  equality is a classic bug — use `abs(a - b) < epsilon` instead.
- **AI:** float16/bfloat16 training (mixed precision) trades mantissa
  precision for speed/memory — understanding *why* explains why some
  models see NaN losses (exponent overflow) at low precision.
- **Network:** checksums and packet length fields are fixed-width integers;
  overflow/wraparound is a real parsing bug source (e.g., TCP sequence
  number wraparound at 2³², handled by "serial number arithmetic," RFC 1982).
- **Security:** integer overflow (CWE-190) and floating-point rounding
  errors in financial systems are exploitable classes — "salami slicing"
  attacks historically exploited rounding remainders.

**More real-life examples:**
- The Boeing 787 Dreamliner's generator control units had a documented
  32-bit counter overflow (FAA Airworthiness Directive 2015-19-07) that
  forced a shutdown after 248 days of continuous power — the same
  overflow class as the toy 8-bit example above, at safety-critical scale.
- The "Year 2038 problem": Unix time stored as a signed 32-bit integer
  overflows on January 19, 2038 — every system still storing timestamps
  this way inherits this chapter's wraparound bug on a fixed calendar date.
- In 2010, a Bitcoin integer-overflow bug let someone mint 184.4 billion
  BTC in a single transaction before the network was patched — a direct
  real-money consequence of fixed-width arithmetic.
- Payment APIs (e.g., Stripe) represent amounts as integer minor units
  (cents), not decimals, specifically to avoid the floating-point rounding
  drift shown in the `0.1 + 0.2` example above.

### Common mistake

Assuming `int` in your language is unbounded. Python's `int` is arbitrary
precision (auto-grows), but C, Java, Go's fixed-width ints are not — this
difference alone explains many "works in Python, crashes in C" ports.

---

## 2. Modular Arithmetic

**Level:** Beginner
**You'll use this for:** hashing, hash tables, load balancing, checksums, and the entire foundation of cryptography (Ch. 24).

### The idea

`a mod n` is the remainder when `a` is divided by `n`. Think of a clock:
`14:00 mod 12 = 2:00`. Numbers "wrap around" after reaching `n`.

```python
>>> 17 % 5     # 2
>>> -1 % 5     # 4 in Python (always non-negative); -1 in C/Java (sign of dividend)
```

**Careful:** the sign convention for negative numbers differs across
languages. Python and math convention give a non-negative result; C, Java,
and JavaScript give a result with the sign of the dividend. This is a real
source of off-by-one/wraparound bugs when porting code.

### How to compute `a mod n` by hand

`a mod n = a - n · ⌊a/n⌋` — divide, round the quotient *down* to the
nearest whole number (floor), multiply that back by `n`, and subtract from
`a`. Whatever's left is the remainder.

```
17 mod 5:
  17 ÷ 5 = 3.4  ->  floor to 3     (5 fits into 17 three whole times)
  3 × 5 = 15
  17 - 15 = 2                       -> 17 mod 5 = 2   ✓ matches the code above

-1 mod 5 (Python/math convention):
  -1 ÷ 5 = -0.2  ->  floor to -1    (floor rounds toward negative infinity,
                                      not toward zero — this is the crux of
                                      the sign difference vs. C/Java)
  -1 × 5 = -5
  -1 - (-5) = 4                     -> -1 mod 5 = 4    ✓ matches the code above
```

C-family languages instead **truncate** toward zero (`-1 ÷ 5 = -0.2 -> 0`,
giving `-1 - 0 = -1`), which is exactly why the same expression produces a
different sign in Python vs. C — the division step rounds differently, not
the subtraction.

### Congruence

`a ≡ b (mod n)` means `a` and `b` have the same remainder mod `n`. This is
the backbone of clock arithmetic, cyclic buffers, and (later) RSA.

### Worked example: round-robin load balancing

You have 5 backend servers (indices 0–4) and need to route request #`i` to
server `i mod 5`.

```python
servers = ["srv0", "srv1", "srv2", "srv3", "srv4"]
for request_id in [0, 1, 5, 6, 12]:
    print(request_id, "->", servers[request_id % 5])
# 0 -> srv0, 1 -> srv1, 5 -> srv0, 6 -> srv1, 12 -> srv2
```

### Worked example: hash table bucket index

A hash table with `n` buckets computes `bucket = hash(key) % n`. This is
*why* hash table sizes are often chosen as powers of 2 — `x % 2^k` is just
`x & (2^k - 1)`, a single fast bitwise AND instead of a division.

```python
n = 16               # 2^4 buckets
key_hash = 0b101101010111  # some hash value
bucket_via_mod = key_hash % n
bucket_via_mask = key_hash & (n - 1)
assert bucket_via_mod == bucket_via_mask
```

### Real-life engineering ties

- **SWE:** circular buffers, ring indices (`idx = (idx + 1) % capacity`),
  date/time math (day-of-week = `days_since_epoch % 7`).
- **Network:** TCP sequence numbers wrap at 2³² and are compared using
  *modular* ("serial number") arithmetic, not plain `<`/`>` (RFC 1982) —
  otherwise wraparound would make a new packet look "older" than an old one.
- **Security:** every symmetric/asymmetric cipher, hash function, and CRC
  checksum operates in modular arithmetic. Consistent hashing (used to
  shard data across nodes with minimal reshuffling when nodes are added) is
  modular arithmetic on a "hash ring."

**More real-life examples:**
- Distributed databases and caches (DynamoDB, Cassandra, memcached) use
  consistent hashing to shard keys across nodes with minimal data movement
  when nodes are added or removed.
- Credit card numbers, ISBNs, and IMEI numbers all carry a check digit
  computed with modular arithmetic (the Luhn algorithm) so a single
  mistyped digit is caught before the number is even looked up.
- Ethernet frames, ZIP files, and PNG images all carry a CRC32 checksum —
  a value computed over the data modulo a fixed generator polynomial — to
  catch accidental corruption in transit.
- Round-robin DNS and round-robin CPU scheduling both assign the next
  resource using `index mod n`, the same operation as the load-balancer
  example above.

### Common mistake

Forgetting that `%` in most C-family languages can return **negative**
results, breaking array-index code that assumes `0 <= result < n`. Fix:
`((a % n) + n) % n`.

---

## 3. Boolean Algebra & Logic

**Level:** Beginner
**You'll use this for:** every `if` statement, firewall rule, SQL `WHERE` clause, and digital circuit ever built.

### The idea

Boolean algebra has two values (`true`/`false`, `1`/`0`) and three core
operators: AND (`∧`), OR (`∨`), NOT (`¬`). Everything in digital logic —
CPUs, firewalls, query planners — reduces to combinations of these.

### Truth tables

| A | B | A AND B | A OR B | A XOR B |
|---|---|---|---|---|
| 0 | 0 | 0 | 0 | 0 |
| 0 | 1 | 0 | 1 | 1 |
| 1 | 0 | 0 | 1 | 1 |
| 1 | 1 | 1 | 1 | 0 |

XOR ("exactly one is true") is the odd one out and deserves special
attention: it's the building block of parity checks, checksums, and the
one-time pad cipher (Ch. 24) — because `X ⊕ K ⊕ K = X` (XOR-ing twice with
the same key returns the original).

### How to evaluate a compound expression by hand

For a multi-variable expression like `NOT(A AND B) OR C`, build the truth
table one operator at a time, innermost first — exactly how a computer
evaluates it, following operator precedence (NOT binds tightest, then AND,
then OR):

```
A B C | A AND B | NOT(A AND B) | NOT(A AND B) OR C
0 0 0 |    0     |      1       |         1
0 0 1 |    0     |      1       |         1
0 1 0 |    0     |      1       |         1
0 1 1 |    0     |      1       |         1
1 0 0 |    0     |      1       |         1
1 0 1 |    0     |      1       |         1
1 1 0 |    1     |      0       |         0
1 1 1 |    1     |      0       |         1
```

Each row plugs A, B, C into the *inner* expression first (column 4:
`A AND B`), then applies the next operator to that result (column 5), then
the last one (column 6) — a compound boolean expression is just several
truth-table lookups chained together, evaluated left-to-right by
precedence.

### De Morgan's laws

```
¬(A ∧ B) = ¬A ∨ ¬B
¬(A ∨ B) = ¬A ∧ ¬B
```

**Why it matters:** these let you simplify/rewrite conditions correctly.
`!(is_admin && is_active)` is the same as `!is_admin || !is_active` — a
transformation you'll do constantly when negating guard clauses, and get
*wrong* constantly if you don't know the law (a classic source of
authorization bugs — accidentally granting access because a negation was
distributed incorrectly).

### Worked example: firewall rule as a boolean expression

```
ALLOW if:  (src_ip IN trusted_range) AND (dest_port == 443) AND NOT (src_ip IN blocklist)
```
This is a literal boolean circuit. Firewall rule *ordering* matters because
most engines short-circuit (stop evaluating once the outcome is
determined) — same as `&&`/`||` in code.

### Bitwise flags

```python
READ, WRITE, EXECUTE = 0b100, 0b010, 0b001   # 4, 2, 1
perms = READ | WRITE          # 0b110 = 6  (like chmod 6xx)
can_write = bool(perms & WRITE)   # True
perms_no_write = perms & ~WRITE   # clear the WRITE bit
```

**Working the bit arithmetic by hand**, one column at a time — OR sets a
bit if *either* input has it, AND keeps a bit only where *both* inputs have
it:

```
  READ    = 100
  WRITE   = 010
  OR:       ---
  perms   = 110   (6)   <- each column: 1 if either input bit is 1

  perms   = 110
  WRITE   = 010
  AND:      ---
  result  = 010   (2, nonzero -> truthy -> can_write = True)
             <- each column: 1 only if BOTH input bits are 1

  ~WRITE (flip every bit of 010, in an 8-bit view) = 11111101
  perms      = 00000110
  AND:         --------
  result     = 00000100   (4 = READ only -> WRITE bit cleared)
```

This is literally how Unix file permissions, HTTP method flags, and
capability systems (Linux `CAP_*`) are implemented — a permission set is a
bitmask, and checking/granting/revoking permissions is AND/OR/NOT/XOR.

### Real-life engineering ties

- **SWE:** SQL `WHERE` clause optimization, guard-clause refactors, feature
  flag combination logic.
- **Network:** subnet mask arithmetic *is* bitwise AND (Ch. 25); ACL/firewall
  rule evaluation.
- **Security:** access-control logic (RBAC/ABAC rule engines), WAF rule
  chains, and — critically — **auditing boolean logic for that one `OR`
  that should have been `AND`** is one of the most common real-world
  authorization vulnerabilities (broken access control, OWASP #1).

**More real-life examples:**
- Unix file permissions (`chmod 754`) are three 3-bit boolean masks
  (`rwx`) combined with OR when granting and AND/NOT when checking — the
  exact bitwise-flag mechanism from this chapter, running on every file
  access on every Linux/macOS machine.
- AWS IAM and most cloud policy engines evaluate access as boolean
  expressions with an explicit precedence rule ("an explicit Deny always
  overrides an Allow") — misunderstanding this precedence is a common
  source of real cloud misconfiguration incidents.
- A CPU's Arithmetic Logic Unit (ALU) is physically built from AND/OR/
  XOR/NOT logic gates — every instruction a processor executes ultimately
  reduces to the truth tables in this chapter.
- CORS preflight logic in browsers, and WAF rule engines (e.g.,
  ModSecurity, Cloudflare rules), evaluate requests through chained
  boolean conditions — a misplaced `OR` in one of these rule sets is a
  real, recurring class of security misconfiguration.

### Common mistake

Operator precedence: `A OR B AND C` parses as `A OR (B AND C)`, not
`(A OR B) AND C`. Always parenthesize explicitly in security-critical logic
— an unparenthesized boolean expression is a real, recurring CVE root
cause.

---

## 4. Sets, Relations, and Functions

**Level:** Beginner
**You'll use this for:** databases, type systems, deduplication, and reasoning about API contracts.

### Sets

A set is an unordered collection of distinct elements. Core operations:

| Operation | Symbol | Meaning | Python |
|---|---|---|---|
| Union | `A ∪ B` | everything in A or B | `a \| b` |
| Intersection | `A ∩ B` | in both A and B | `a & b` |
| Difference | `A − B` | in A, not B | `a - b` |
| Subset | `A ⊆ B` | every element of A is in B | `a <= b` |

```python
admins = {"alice", "bob"}
active = {"alice", "carol"}
print(admins & active)   # {'alice'}          -> active admins
print(admins - active)   # {'bob'}            -> inactive admins
print(admins | active)   # {'alice','bob','carol'}
```

**How to compute these by hand:** list both sets, then scan for membership.

```
admins = {alice, bob}      active = {alice, carol}

Intersection (A ∩ B): keep an element only if it's in BOTH lists.
  alice -> in admins? yes. in active? yes.  -> keep
  bob   -> in admins? yes. in active? no.   -> drop
  carol -> in admins? no.                    -> drop
  result: {alice}

Difference (A - B): keep an element from A only if it's NOT in B.
  alice -> in active? yes -> drop
  bob   -> in active? no  -> keep
  result: {bob}

Union (A ∪ B): every element that appears in either list, no duplicates.
  result: {alice, bob, carol}
```

A hash-set implementation does exactly this scan, just in `O(1)` per
membership check instead of scanning a list — which is *why* the
data-structure choice mentioned in the common mistake below matters.

### Relations and functions

A **relation** is any set of pairs (a foreign key relationship, a graph
edge). A **function** is a relation where every input maps to *exactly one*
output — this constraint is why a database primary key must be unique, and
why a hash function must be deterministic (same input → same output, always).

- **Injective (one-to-one):** different inputs never collide on the same
  output. A perfect hash function would be injective; real hash functions
  aren't (hence collisions, Ch. 26).
- **Surjective (onto):** every possible output is reachable from some input.
- **Bijective:** both — this is what makes a function *invertible*.
  Encryption must be bijective (you must be able to decrypt); hashing must
  **not** be invertible by design (one-way).

### Real-life engineering ties

- **SWE:** SQL `JOIN`/`UNION`/`INTERSECT` are literally set operations;
  `DISTINCT` is "convert to a set."
- **AI:** train/validation/test splits must be **disjoint sets** — a
  leaking overlap between train and test sets is one of the most common,
  most silent causes of an ML model that looks great in evaluation and
  fails in production.
- **Network:** IP address *ranges* are sets; CIDR blocks are set membership
  tests (Ch. 25).
- **Security:** allowlists/denylists are sets; "is this IP in the
  blocklist" is a set-membership query, and the *data structure* backing
  that query (hash set vs. sorted array vs. bloom filter) is a real
  performance/security tradeoff at scale.

**More real-life examples:**
- Redis's native `SET` type and operations like `SINTERSTORE`/
  `SUNIONSTORE` are used in production to answer questions like "which
  users are both online and in this user's friend list" — a live
  intersection query at scale.
- Google Chrome's Safe Browsing and many CDNs use **Bloom filters** — a
  probabilistic, space-efficient set-membership structure — to check "have
  we seen this URL/hash before" without storing every value explicitly.
- Database foreign-key constraints formally encode a relation between two
  tables; a "dangling foreign key" bug is a relation whose function
  property (every child row maps to exactly one valid parent) has been
  violated.
- API contract/type systems (OpenAPI schemas, TypeScript types) define the
  *set* of valid inputs and outputs a function may accept and return — a
  type error is really a set-membership violation caught at compile time.

### Common mistake

Treating a `list` where you need a `set` — leads to accidental O(n)
membership checks (and, in security-critical allowlist code, accidental
duplicate/near-duplicate entries that create bypass gaps).

---

## 5. Exponents & Logarithms

**Level:** Beginner
**You'll use this for:** understanding Big-O, decibels/signal strength, information content, and every "how many bits" question.

### Exponent rules

```
x^a · x^b = x^(a+b)
x^a / x^b = x^(a-b)
(x^a)^b   = x^(ab)
x^0 = 1        x^(-a) = 1/x^a
```

### Logarithms — the inverse of exponentiation

`log_b(x) = y` means `b^y = x`. "How many times do I multiply `b` by itself
to reach `x`?"

```
log2(8)   = 3      (2^3 = 8)
log2(1024)= 10     (2^10 = 1024)
log10(1000)=3      (10^3 = 1000)
```

**Why `log₂` specifically matters to engineers:** it answers "how many bits
do I need?" `log2(n)` rounded up is the number of bits needed to represent
`n` distinct values. This is *the same number* that shows up as the depth
of a balanced binary search tree, the number of comparisons in binary
search, and the number of "bits of entropy" in a random value.

```python
import math
math.log2(1_000_000)   # ≈ 19.93 -> need 20 bits to enumerate a million things
math.log2(256)         # 8.0    -> a byte has 8 bits, 256 possible values
```

### How to estimate `log2(x)` by hand — bracket between powers of 2

You don't need a calculator to get a useful estimate: find the two powers
of 2 that `x` sits between, and interpolate.

```
Estimating log2(1,000,000):
  2^19 = 524,288        <- too small
  2^20 = 1,048,576      <- just above 1,000,000

  So log2(1,000,000) is between 19 and 20 — closer to 20, since
  1,000,000 is much nearer 1,048,576 than 524,288.
  (The exact value, ≈19.93, confirms the estimate.)
```

This "which two powers of 2 does it sit between" check is exactly how you
mentally answer "roughly how many bits do I need for a million things" —
you don't need the precise decimal, just the bracket, which is why
engineers eyeball `log2` constantly without reaching for a calculator.

### Worked example: why binary search is O(log n)

Each comparison halves the search space. Starting with `n` items, after `k`
halvings you have `n / 2^k` items left. Search ends when that's 1:
`n / 2^k = 1` → `2^k = n` → `k = log2(n)`. For a billion items, that's
`log2(10^9) ≈ 30` comparisons — not a billion.

### Log rules that matter in practice

```
log(a·b) = log(a) + log(b)      <- multiplication becomes addition
log(a/b) = log(a) - log(b)
log(a^k) = k·log(a)
```
This is *why* decibels (`dB = 10·log10(P/P_ref)`), pH, the Richter scale,
and `perplexity`/log-loss in ML all use logs: they compress a value that
spans many orders of magnitude into a manageable additive scale, and they
turn products of small probabilities (which underflow to 0 in floating
point) into sums.

```python
# Why ML loss functions use log-probabilities, not raw probabilities:
probs = [0.001] * 50          # 50 independent low-probability events
raw_product = 1.0
for p in probs: raw_product *= p
print(raw_product)            # 0.0  <- underflowed to zero!

log_sum = sum(math.log(p) for p in probs)
print(log_sum)                 # -345.4  <- still numerically meaningful
```

### Real-life engineering ties

- **SWE:** Big-O reasoning (`O(log n)` algorithms), pagination depth
  estimates.
- **AI:** cross-entropy loss, log-likelihood, perplexity, log-sum-exp trick
  for numerical stability in softmax (Ch. 22–23).
- **Network:** signal-to-noise ratio and channel capacity are log-based
  (Shannon-Hartley, Ch. 25/27); dB loss budgets in fiber links.
- **Security:** password/key entropy is measured in **bits**, i.e., `log2`
  of the keyspace size (Ch. 26).

**More real-life examples:**
- `git bisect` finds a bad commit among thousands in `O(log n)` steps by
  binary search — the same halving logic that makes binary search fast,
  applied to version history instead of a sorted array.
- Earthquake magnitude (Richter/moment scale) and sound intensity
  (decibels) are both logarithmic: a magnitude-7 earthquake releases
  roughly 32x more energy than a magnitude-6, not "1 more unit" of
  damage — exactly the compression effect `log` provides for values
  spanning huge ranges.
- Large language model evaluation reports **perplexity**, which is
  literally `2^(cross-entropy loss)` — a log-scale metric chosen for the
  same numerical-underflow reasons shown in the code example above.
- Database index design assumes `O(log n)` lookups (B-trees) specifically
  so lookups stay fast as a table grows from thousands to billions of
  rows — the log term is why adding 10x more rows barely changes query time.

### Common mistake

Confusing `log` base — most math textbooks default to `ln` (base *e*), most
CS contexts default to `log2`, most engineering/audio contexts default to
`log10`. Always check which base a formula assumes.

---

## 6. Sequences, Series & Growth Rates

**Level:** Beginner
**You'll use this for:** capacity planning, understanding "doubling," and recognizing exponential trouble early.

### Arithmetic vs. geometric

- **Arithmetic sequence:** constant difference. `2, 5, 8, 11, ...` (+3 each
  step). Sum of first `n` terms: `n/2 · (first + last)`.
- **Geometric sequence:** constant ratio. `2, 4, 8, 16, ...` (×2 each step).
  `n`th term: `a·r^(n-1)`. Sum of first `n` terms (r≠1): `a·(r^n - 1)/(r-1)`.

### Where these sum formulas actually come from

**Arithmetic sum — Gauss's pairing trick.** To sum `1, 2, 3, ..., 10`,
write the sum forwards and backwards underneath itself, and add
column-by-column:

```
   1  +  2  +  3  + ... +  9  + 10
 +10  +  9  +  8  + ... +  2  +  1
 ------------------------------------
  11  + 11  + 11  + ... + 11  + 11     <- every column sums to 11 (first+last)

10 columns × 11 = 110, but that's DOUBLE the sum (we added it to itself),
so the actual sum = 110 / 2 = 55.
```

Generalizing: `n` columns each summing to `(first + last)`, then halved —
exactly the formula `n/2 · (first + last)` given above.

**Geometric sum — the telescoping trick.** To sum `S = a + ar + ar² + ... +
ar^(n-1)`, multiply the whole sum by `r` and subtract:

```
        S =   a + ar + ar² + ... + ar^(n-1)
       rS =       ar + ar² + ... + ar^(n-1) + ar^n
  --------------------------------------------------
   S - rS =   a                              - ar^n     <- everything in the
                                                            middle cancels

  S(1 - r) = a(1 - r^n)
  S = a(1 - r^n) / (1 - r)     <- same as a(r^n - 1)/(r - 1), signs flipped
```

Both derivations use the same move: rewrite the sum a second way so that
almost everything cancels, leaving a formula with no loop required.

### The Rule of 72 (mental-math doubling time)

For any exponential growth at rate `r`% per period, time to double ≈
`72 / r`. Growing disk usage at 8%/month? Doubles in ~9 months. This is the
single fastest sanity check for "is this growth curve about to become a
problem."

```python
def doubling_time(rate_percent):
    return 72 / rate_percent

print(doubling_time(8))    # 9.0 months
```

### Worked example: viral growth / incident blast radius

A worm that infects 1 host, and each infected host infects 2 more per hour,
follows `hosts(t) = 3^t` (1 becomes 3 total: itself + 2 new, compounding).
After 10 hours: `3^10 = 59,049` hosts. This is why "time to detect and
contain" matters exponentially more than "time to detect" alone — an extra
hour of dwell time in an exponential spread isn't linear damage, it's
multiplicative.

```python
for t in range(0, 11):
    print(t, 3**t)
# by hour 10: 59,049 — by hour 15: 14,348,907
```

### Real-life engineering ties

- **SWE:** Fibonacci-style retry backoff, amortized array-doubling
  (`ArrayList`/`Vec` capacity growth), pagination cost.
- **AI:** dataset/model size scaling laws are often near-geometric;
  training cost vs. model size follows power-law curves.
- **Network:** exponential backoff for TCP retransmission and Wi-Fi CSMA/CA
  collision avoidance (Ch. 25) is literally a geometric sequence
  (`wait = base · 2^attempt`).
- **Security:** malware/worm propagation modeling, credential-stuffing
  attack scaling, and understanding why "small" unpatched vulnerabilities
  compound (an attacker chaining 3 low-severity bugs can have
  multiplicatively higher impact than any one alone).

**More real-life examples:**
- The 2016 Mirai botnet and the 2017 WannaCry ransomware both spread by
  having each infected device scan for and infect several more —
  textbook exponential growth, which is why "time to detect" mattered
  exponentially more than any single infection.
- Compound interest (`A = P(1+r)^t`) is the same geometric-sequence
  formula as this chapter's worm-spread example, just with a friendlier
  sign — "financial literacy" and "understanding exponential attack
  growth" are the same math with different subject matter.
- Cloud cost forecasting (e.g., "our S3 storage bill is growing 6%/month")
  uses the Rule of 72 to sanity-check budget requests: 6%/month doubles
  storage spend in about 12 months, a number FinOps teams use directly.
- AWS SDKs and gRPC clients implement retry backoff as a geometric
  sequence (`delay = base * 2^attempt`, often with jitter) specifically so
  retry storms don't recreate the exponential-growth problem they're
  trying to avoid.

### Common mistake

Assuming a curve that "looks linear" for the first few data points will
stay linear. Exponential curves look almost flat right up until they don't
— always check the *ratio* between consecutive periods, not just the
absolute delta.

---

## 7. Geometry & Trigonometry Essentials

**Level:** Beginner
**You'll use this for:** GPS/location math, graphics, antenna/RF angles, and vector similarity (used constantly in ML/search).

### Pythagorean theorem

`a² + b² = c²` for a right triangle. This generalizes directly to
**Euclidean distance** in n dimensions:

```
distance(P, Q) = sqrt( (p1-q1)² + (p2-q2)² + ... + (pn-qn)² )
```

```python
import math
def euclidean(p, q):
    return math.sqrt(sum((a-b)**2 for a, b in zip(p, q)))

euclidean([0,0], [3,4])         # 5.0  (classic 3-4-5 triangle)
euclidean([1,2,3], [4,6,3])     # distance in 3D
```

This is *the* similarity metric behind k-nearest-neighbors, clustering
(k-means), and "find similar embeddings" in vector databases.

### Sine, cosine, and the unit circle

For an angle `θ` on a unit circle: `x = cos(θ)`, `y = sin(θ)`. Radians, not
degrees, are the native unit in almost all math libraries (`2π` radians = 360°).

```python
math.radians(180)   # π
math.sin(math.pi/2) # 1.0
```

### Cosine similarity — the workhorse of ML/search

```
cos_sim(A, B) = (A · B) / (|A| |B|)
```
Measures the angle between two vectors, ignoring their magnitude — used
constantly to compare text/image embeddings ("how similar is this document
to that one") because it cares about *direction* (meaning) not length
(magnitude, which can vary with document length or brightness).

```python
import numpy as np
def cosine_similarity(a, b):
    a, b = np.array(a), np.array(b)
    return np.dot(a, b) / (np.linalg.norm(a) * np.linalg.norm(b))

cosine_similarity([1,1,0], [1,0,0])   # ≈ 0.707 -> 45° apart
```

**Computing that 0.707 by hand**, three small steps:

```
A = [1, 1, 0]     B = [1, 0, 0]

1. Dot product (A · B): multiply matching positions, sum them.
     (1×1) + (1×0) + (0×0) = 1 + 0 + 0 = 1

2. Magnitudes (|A|, |B|): Pythagoras (Ch. 7 above), sqrt of sum of squares.
     |A| = sqrt(1² + 1² + 0²) = sqrt(2) ≈ 1.414
     |B| = sqrt(1² + 0² + 0²) = sqrt(1) = 1.0

3. Divide: cos_sim = 1 / (1.414 × 1.0) = 1 / 1.414 ≈ 0.707
```

`cos(45°) ≈ 0.707` too — confirming the two vectors are 45° apart, which
checks out visually: `[1,1,0]` points diagonally, `[1,0,0]` points straight
along one axis.

### The Haversine formula (great-circle distance on a sphere)

Straight-line (Euclidean) distance is wrong for lat/long coordinates
because Earth is a sphere. The Haversine formula accounts for curvature:

```python
def haversine(lat1, lon1, lat2, lon2):
    R = 6371  # Earth radius in km
    p1, p2 = math.radians(lat1), math.radians(lat2)
    dp = math.radians(lat2 - lat1)
    dl = math.radians(lon2 - lon1)
    a = math.sin(dp/2)**2 + math.cos(p1)*math.cos(p2)*math.sin(dl/2)**2
    return 2 * R * math.asin(math.sqrt(a))

haversine(40.7128, -74.0060, 51.5074, -0.1278)  # NYC to London ≈ 5570 km
```

### Real-life engineering ties

- **SWE:** map/location features, "nearby" search, geofencing.
- **AI:** cosine similarity for embedding search (RAG/vector DBs), image
  transforms (rotation matrices, Ch. 10).
- **Network:** RF/antenna beam angles, satellite link elevation angles,
  signal path geometry for line-of-sight microwave links.
- **Security:** anomaly detection on "impossible travel" (login from NYC,
  then Tokyo, 10 minutes later) uses Haversine distance ÷ time to flag
  physically implausible logins.

**More real-life examples:**
- GPS receivers compute your position via **trilateration** — solving a
  system of sphere-distance equations (a direct extension of the
  Pythagorean/Euclidean distance idea) using signals from 4+ satellites.
- Ride-share and mapping apps (Uber, Google Maps) combine Haversine
  "as the crow flies" distance with graph-based road-network distance
  (Ch. 17) to estimate ETAs and match nearby drivers.
- Music and video recommendation systems (Spotify, YouTube) compare
  embedding vectors with cosine similarity at massive scale — billions of
  comparisons a day resting on the same formula as the 3-line example above.
- Responsive UI/graphics layout (aspect ratios, scaling, rotation of
  elements on screen) is applied trigonometry and matrix transforms
  (Ch. 10) running in every browser's rendering engine.

### Common mistake

Using Euclidean distance on latitude/longitude directly ("Pythagorean
distance on lat/long") — it's wrong by an amount that grows with distance
and gets *worse* near the poles, since a degree of longitude shrinks as you
move away from the equator.

---

## 8. Basic Probability

**Level:** Beginner
**You'll use this for:** reasoning about A/B tests, retry logic, and every "how likely is this" question, including security alert quality.

### Core definitions

- **Sample space (`Ω`):** every possible outcome.
- **Event (`A`):** a subset of outcomes you care about.
- **P(A):** a number from 0 to 1. `P(A) = 0` never happens, `P(A) = 1`
  always happens.

### Independence vs. dependence

Two events are **independent** if one doesn't affect the other's
probability: `P(A ∩ B) = P(A) · P(B)`. Coin flips are independent (the coin
has no memory) — this is why "I've flipped 5 heads in a row, tails is due"
(the gambler's fallacy) is simply wrong. Each flip is still 50/50.

### Worked example: retry probability

A flaky network call succeeds with probability `p = 0.7` per attempt,
independently. What's the probability all 3 retries fail?

```
P(fail once) = 1 - p = 0.3
P(fail 3 times in a row) = 0.3³ = 0.027   (2.7%)
P(at least one success in 3 tries) = 1 - 0.027 = 0.973
```

```python
p_fail = 1 - 0.7
p_all_fail = p_fail ** 3
print(1 - p_all_fail)   # 0.973
```

This is the actual math behind "3 retries with independent failure modes
gets you from 70% to 97.3% reliability" — and *why* retries must hit
independent failure domains (different replicas/AZs) to actually
compound; retrying the same overloaded server 3 times doesn't multiply
independent probabilities, because the failures are correlated.

### Union of events (why "or" isn't just addition)

```
P(A ∪ B) = P(A) + P(B) - P(A ∩ B)
```
You subtract the overlap so you don't double count it. Forgetting the
`-P(A ∩ B)` term is the single most common basic probability bug.

### The three axioms everything else is built on

In 1933 Andrey Kolmogorov showed that all of probability theory follows from
three rules. Every formula in this guide that involves `P(...)`, from Bayes to
the Poisson distribution, can be derived from these three.

```
1.  P(A) ≥ 0                                  for every event A
2.  P(Ω) = 1                                  something in the sample space happens
3.  P(A ∪ B) = P(A) + P(B)                    when A and B can't both happen (A ∩ B = ∅)
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `Ω` | "omega" | the sample space: every possible outcome |
| `A ∪ B` | "A union B" | A happens, or B happens, or both |
| `A ∩ B` | "A intersect B" | A and B both happen |
| `∅` | "empty set" | an impossible event |

**Example:** the complement rule, `P(not A) = 1 - P(A)`, isn't a separate
fact you memorize. `A` and `not A` can't both happen, and together they cover
`Ω`. So by axiom 3, `P(A) + P(not A) = P(Ω)`, and by axiom 2 that equals 1.
The retry example above (`1 - 0.027`) used exactly this.

### Conditional probability: updating on new information

```
P(A | B) = P(A ∩ B) / P(B)          (defined when P(B) > 0)
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `P(A \| B)` | "probability of A given B" | the probability of A once you know B happened |
| `P(A ∩ B)` | "probability of A and B" | both happen together |
| `P(B)` | "probability of B" | the condition, which becomes your new "whole world" |

Read it as **zooming in**. Once you know `B` happened, every outcome outside
`B` is gone. `B` becomes the new universe, and you ask what fraction of it
also lies inside `A`.

**Example:** from logs, 2% of requests are to `/checkout`, and 0.5% of *all*
requests are `/checkout` requests that failed. What is the failure rate
*of checkout requests*?

```
P(fail | checkout) = P(fail ∩ checkout) / P(checkout) = 0.005 / 0.02 = 0.25  -> 25%
```

A global failure rate that looks fine can hide one endpoint that fails a
quarter of the time. This is why SLO dashboards slice by route.

### The law of total probability: divide and conquer

If `B₁, B₂, …, Bₙ` split the world into non-overlapping cases that cover
everything:

```
P(A) = Σᵢ P(A | Bᵢ) · P(Bᵢ)
```

| Symbol | Meaning |
|---|---|
| `Bᵢ` | the `i`-th case (e.g. "request routed to zone i") |
| `P(Bᵢ)` | how often case `i` happens (the weights must sum to 1) |
| `P(A \| Bᵢ)` | how likely `A` is *inside* case `i` |
| `Σᵢ` | add up over every case `i = 1 … n` |

**Example (networking):** a global load balancer sends 50% of traffic to
zone 1, 30% to zone 2 and 20% to zone 3. The per-zone packet-loss rates are
0.1%, 0.2% and 1%. What is the overall loss rate?

```
P(loss) = 0.001·0.5 + 0.002·0.3 + 0.010·0.2
        = 0.0005 + 0.0006 + 0.0020
        = 0.0031   -> 0.31%
```

Now run it backwards with Bayes (Ch. 16): given that a packet was lost, how
likely is it that it went through zone 3? `0.0020 / 0.0031 ≈ 65%`. Zone 3
carries only 20% of traffic but causes 65% of the loss, so that's where to
look first. Total probability plus Bayes is the maths behind "which
component is responsible for most of the errors?"

### Random variables and expected value

A **random variable** `X` is a number whose value is decided by chance, such
as the number of retries, a request's latency, or the bytes in a flow. Its
**expected value** is the long-run average you'd see over many repetitions:

```
E[X] = Σₓ x · P(X = x)                 (discrete: weighted sum of outcomes)
E[X] = ∫ x · f(x) dx                    (continuous: weighted integral, Ch. 12)
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `E[X]` | "the expectation of X" | the probability-weighted average value |
| `x` | "little x" | one particular value `X` can take |
| `P(X = x)` | "probability X equals x" | how often that value occurs |
| `f(x)` | "the density of X at x" | the continuous version of `P(X = x)` (Ch. 14) |

**Example: how many attempts does a flaky call take, on average?** Each
attempt succeeds independently with probability `p`. The number of attempts
until the first success follows a *geometric* distribution, and its
expectation works out to:

```
E[attempts] = 1·p + 2·(1-p)·p + 3·(1-p)²·p + ...  =  1/p
```

With `p = 0.7`, that's `1/0.7 ≈ 1.43` attempts per call on average, so a
retrying client puts about **43% more load** on the backend than a
non-retrying one. A backend that fails *more* (smaller `p`) gets *more*
traffic. That feedback loop is how retry storms happen.

**Networking in the wild: the ETX routing metric.** In wireless mesh networks
(MIT's Roofnet, then the Linux OLSR and B.A.T.M.A.N. routing daemons), a
link's cost is its *expected transmission count*. If a data frame gets
through with probability `d_f` and its ACK comes back with probability
`d_r`, an attempt succeeds with probability `d_f · d_r`, so:

```
ETX = 1 / (d_f · d_r)
```

| Path | Calculation | Expected transmissions |
|---|---|---|
| 1 hop, a bad link (`d_f = d_r = 0.5`) | `1 / (0.5·0.5)` | **4.0** |
| 2 hops, two good links (`0.9` each way) | `2 × 1 / (0.9·0.9)` | **2.47** |

Hop-count routing picks the 1-hop path because it has fewer hops. ETX
routing picks the 2-hop path because, on average, it uses the air about 40%
less. The difference is just expected value.

### Linearity of expectation: the most useful "free lunch" in probability

```
E[X + Y] = E[X] + E[Y]         ALWAYS, even if X and Y are dependent
E[c · X] = c · E[X]            for any constant c
```

The surprising part is that this works **without independence**. That lets
you break a hard count into easy indicator pieces and add up their
expectations.

**Example: expected hash collisions.** Insert `n = 1000` keys into a hash
table with `m = 10,000` buckets. How many *pairs* of keys share a bucket, on
average? For each of the `C(n,2)` pairs (Ch. 15), let an indicator be 1 if
that pair collides. Each pair collides with probability `1/m`. So:

```
E[colliding pairs] = C(n,2) · (1/m) = (1000·999/2) / 10000 = 499,500 / 10,000 ≈ 50
```

That's about 50 collisions with the table only 10% full. Linearity of
expectation answered it in one line, with no need to work out the messy
joint distribution. The same argument gives the birthday bound in Ch. 24.

**AI in the wild: why dropout rescales.** During training, dropout zeroes
each neuron's output `h` with probability `1 - q` (so it keeps it with
probability `q`). The expected output becomes `E[kept h] = q·h + (1-q)·0 = q·h`.
At inference time nothing is dropped, so the next layer would suddenly see
inputs `1/q` times larger than it was trained on. "Inverted dropout", the
default in PyTorch, divides by `q` during training. By linearity,
`E[h/q · keep] = (1/q)·q·h = h`, so the expected value is the same in
training and inference.

### Variance of a random variable

```
Var(X) = E[(X - E[X])²] = E[X²] - (E[X])²
```

| Symbol | Meaning |
|---|---|
| `X - E[X]` | how far one outcome lands from the average |
| `(…)²` | squared, so misses in either direction count and big misses count more |
| `E[X²] - (E[X])²` | the shortcut form, easier to compute by hand |

**Example: a single request that fails with probability `p` (a Bernoulli
variable, 1 = fail, 0 = success).** `E[X] = p`, and since `X² = X` for 0/1
values, `E[X²] = p`. So `Var(X) = p - p² = p(1-p)`. This is largest at
`p = 0.5` and tiny near 0 or 1. That's why a rare-error metric (`p = 0.001`)
needs **huge** sample sizes before you can trust a change in it, which
Ch. 22 turns into a sample-size formula.

### Real-life engineering ties

- **SWE:** retry/circuit-breaker design, error budget math (SRE), load
  test result interpretation.
- **AI:** every classifier outputs a probability; understanding calibration
  (does "90% confident" actually mean right 90% of the time?) starts here.
- **Network:** packet loss modeling, redundant path availability
  (`1 - (1-p)^n` for n independent paths — Ch. 25 combines this with
  queueing).
- **Security:** **false positive rate** of a detection rule is a
  probability; understanding independence vs. correlation of signals is
  the difference between a SIEM rule that actually reduces risk and one
  that just generates alert fatigue (see Bayes, Ch. 16).

**More real-life examples:**
- Cloud architects use independent-failure probability math directly when
  deciding how many Availability Zones to span: `P(all zones down) = p1 ·
  p2 · p3` only holds if the zones truly fail independently — the same
  assumption checked in the common mistake below.
- Circuit-breaker libraries (Netflix Hystrix, resilience4j) are built on
  the same "probability of success across independent retries" reasoning
  as the retry example above, and track a rolling failure probability to
  decide when to "open" the breaker.
- CAPTCHA and bot-detection systems reason about the probability that a
  *sequence* of behaviors (mouse movement, timing, request pattern) is
  exhibited by a human vs. a script — a compound probability judgment, not
  a single yes/no rule.
- Credential-stuffing attacks exploit the (very non-independent) fact that
  many users reuse the same password across services — attackers rely on
  this correlation, not on independent per-site guessing, which is exactly
  why password reuse defeats per-site rate-limiting defenses.

### Common mistake

Assuming events are independent when they're not (e.g., "each server has a
1% chance of failing, so with 100 servers only 1 fails" ignores correlated
failures — a single power/network outage can take down all 100 at once).
This is the mathematical root of "the cloud availability zone problem."

---

## 9. Basic Statistics

**Level:** Beginner
**You'll use this for:** dashboards, SLOs, experiments, and "is this normal" judgment calls.

### Mean, median, mode

- **Mean:** sum ÷ count. Sensitive to outliers.
- **Median:** the middle value when sorted. Robust to outliers.
- **Mode:** the most frequent value.

**Why this matters for latency dashboards:** if your API has 99 requests at
50ms and 1 request at 50,000ms (a stuck connection), the **mean** jumps to
~549ms — misleadingly bad. The **median** stays at 50ms — correctly
reflecting "almost everyone had a fine experience." This is exactly why SRE
dashboards report **percentiles**, not averages.

### Percentiles (p50, p95, p99)

The value below which `X%` of observations fall. `p99 = 500ms` means 99% of
requests were faster than 500ms (and 1% were slower — that 1% is often
where your angriest users live).

```python
import numpy as np
latencies = [45, 48, 50, 52, 51, 49, 5000, 47, 53, 46]  # ms, one outlier
print("mean:", np.mean(latencies))     # 544.1  <- misleading
print("median:", np.median(latencies)) # 49.0   <- representative
print("p95:", np.percentile(latencies, 95))
print("p99:", np.percentile(latencies, 99))
```

### Variance and standard deviation

Variance measures spread: average of squared distances from the mean.
Standard deviation (`σ`, the square root of variance) puts that spread back
into the original units, which is why alerting thresholds are usually
phrased as "mean ± *k*·σ" ("3-sigma alert").

```
variance = (1/n) · Σ(xᵢ - mean)²
std_dev  = sqrt(variance)
```

```python
data = [10, 12, 9, 11, 60]   # one clear outlier
print(np.mean(data), np.std(data))
```

**Computing this by hand**, step by step:

```
data = [10, 12, 9, 11, 60]

1. Mean:  (10+12+9+11+60) / 5  =  102 / 5  =  20.4

2. Deviation from mean, each point (xi - mean):
     10 - 20.4 = -10.4
     12 - 20.4 =  -8.4
      9 - 20.4 = -11.4
     11 - 20.4 =  -9.4
     60 - 20.4 =  39.6

3. Square each deviation (removes negative signs, penalizes big gaps more):
     (-10.4)² = 108.16
     (-8.4)²  =  70.56
     (-11.4)² = 129.96
     (-9.4)²  =  88.36
     (39.6)²  = 1568.16

4. Average the squared deviations (this average IS the variance):
     (108.16+70.56+129.96+88.36+1568.16) / 5 = 1965.2 / 5 = 393.04

5. Standard deviation = sqrt(variance):
     sqrt(393.04) ≈ 19.82
```

Notice how step 3 makes the outlier (60) dominate the whole result — its
squared deviation (1568.16) is more than 10× any other point's, which is
exactly why one wild data point drags both variance and standard deviation
up so much.

### The normal (Gaussian) distribution and the 68-95-99.7 rule

For normally distributed data: ~68% of values fall within 1σ of the mean,
~95% within 2σ, ~99.7% within 3σ. This is *why* "3-sigma" is a common
default anomaly threshold — a point beyond 3σ has < 0.3% chance of
occurring "naturally," making it a reasonable (if crude) alert trigger.

### Sample vs. population: why some formulas divide by `n - 1`

The variance above divides by `n`. That is correct when your data **is** the
entire population. Usually it's only a *sample* (the last 1,000 requests,
not every request ever), and then you should divide by `n - 1`:

```
population variance:  σ² = (1/N)     · Σᵢ (xᵢ - μ)²
sample variance:      s² = (1/(n-1)) · Σᵢ (xᵢ - x̄)²
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `σ²`, `μ`, `N` | "sigma squared", "mu", "big N" | the true variance, true mean and size of the *whole* population |
| `s²`, `x̄`, `n` | "s squared", "x bar", "n" | the same three quantities estimated from a *sample* |
| `xᵢ` | "x sub i" | the `i`-th observation |
| `n - 1` | "degrees of freedom" | the number of deviations that are free to vary |

**Why `n - 1` (Bessel's correction)?** The deviations are measured from `x̄`,
which was computed from the same data. That makes the data look closer to its
own centre than it is to the true `μ`, so dividing by `n` *systematically
underestimates* the spread. Another way to see it: once you know `x̄`, the
deviations must sum to zero, so only `n - 1` of them are free.

**Example:** the same data `[10, 12, 9, 11, 60]` gives
`1965.2 / 5 = 393.04` with `n`, but `1965.2 / 4 = 491.3` (`s ≈ 22.17`) with
`n - 1`. NumPy's `np.std` divides by `n` by default (`ddof=0`) and pandas'
`.std()` divides by `n - 1` (`ddof=1`), so the "same" calculation can disagree
between two notebooks. Pass `ddof` explicitly.

### z-scores: putting anything on a common ruler

```
z = (x - μ) / σ
```

| Symbol | Meaning |
|---|---|
| `x` | the observation you want to judge |
| `μ`, `σ` | the mean and standard deviation of "normal" behaviour |
| `z` | how many standard deviations `x` sits from the mean (no units) |

**Example:** a host normally sends `μ = 200 MB/hour` of egress traffic with
`σ = 40 MB`. This hour it sent 380 MB:

```
z = (380 - 200) / 40 = 4.5
```

Under a normal model, a 4.5σ event happens about 3 times in a million hours.
That's worth paging for, since it might be data exfiltration. Dividing by `σ` is
also why ML pipelines **standardize** features before training. Packet sizes
(hundreds) and port numbers (thousands) end up on the same scale, so gradient
descent (Ch. 21) doesn't spend all its effort on whichever feature has the
biggest raw numbers.

### Covariance and correlation: do two things move together?

```
Cov(X, Y) = (1/(n-1)) · Σᵢ (xᵢ - x̄)(yᵢ - ȳ)

r = Cov(X, Y) / (s_X · s_Y)                 (Pearson correlation, always in [-1, 1])
```

| Symbol | Meaning |
|---|---|
| `(xᵢ - x̄)(yᵢ - ȳ)` | positive when both are above (or both below) their means at the same time |
| `Cov(X, Y)` | the average of those products: positive = move together, negative = move opposite |
| `s_X`, `s_Y` | the sample standard deviations, dividing them out removes units |
| `r` | `+1` = perfect positive line, `0` = no *linear* relation, `-1` = perfect negative line |

**Example (worked by hand):** five load-test runs, with `x` = payload size
(KB) and `y` = p50 latency (ms):

```
x = [1, 2, 3, 4, 5]     x̄ = 3
y = [2, 4, 5, 4, 5]     ȳ = 4

xᵢ - x̄:            -2  -1   0   1   2
yᵢ - ȳ:            -2   0   1   0   1
product:             4   0   0   0   2     sum = 6

Cov = 6 / 4 = 1.5
s_X² = (4+1+0+1+4)/4 = 2.5      s_Y² = (4+0+1+0+1)/4 = 1.5
r = 1.5 / sqrt(2.5 · 1.5) = 1.5 / 1.936 ≈ 0.77
```

That's a strong positive relationship. The best-fit line's slope is
`Cov/s_X² = 1.5/2.5 = 0.6 ms per KB`, and the intercept is
`ȳ - 0.6·x̄ = 2.2 ms`. That's linear regression in two lines. Ch. 10 derives
the same answer from matrices, which scales to thousands of features.

```python
import numpy as np
x = np.array([1, 2, 3, 4, 5]); y = np.array([2, 4, 5, 4, 5])
print(np.cov(x, y)[0, 1], np.corrcoef(x, y)[0, 1])   # 1.5  0.7746
print(np.polyfit(x, y, 1))                           # [0.6 2.2]  slope, intercept
```

**Remember:** `r = 0` does *not* mean "unrelated". `y = x²` on `x ∈ [-1, 1]`
has `r = 0` but is perfectly determined by `x`. Correlation only measures
*straight-line* relationships. Mutual information (Ch. 23) catches the rest.

### Exponentially weighted moving average (EWMA): statistics with a memory budget

A plain average needs every past sample. An EWMA keeps a **single number**
and nudges it toward each new observation:

```
S_new = (1 - α) · S_old + α · x_new
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `S` | "the smoothed value" | the running estimate |
| `x_new` | | the newest measurement |
| `α` | "alpha", the smoothing factor (0 < α ≤ 1) | how much to trust the new sample. A large `α` reacts fast, a small `α` is smooth but slow |

Unroll it and you'll see sample `k` steps old gets weight `α(1-α)ᵏ`, which is
geometric decay (Ch. 6). Its "memory" is roughly `1/α` samples.

**Networking in the wild: TCP's retransmission timer (RFC 6298).** Every TCP
connection in your kernel runs *two* EWMAs, one on round-trip time and one
on its variability, to decide how long to wait before retransmitting:

```
RTTVAR ← (1 - β) · RTTVAR + β · |SRTT - R|       β = 1/4
SRTT   ← (1 - α) · SRTT   + α · R                α = 1/8
RTO    =  SRTT + 4 · RTTVAR
```

| Symbol | Meaning |
|---|---|
| `R` | the newest RTT sample, measured from a data segment and its ACK |
| `SRTT` | smoothed RTT, the EWMA of `R` |
| `RTTVAR` | the EWMA of the absolute deviation, a cheap stand-in for standard deviation |
| `RTO` | retransmission timeout = mean + 4 × deviation, a "4-sigma" rule (Ch. 9 above) |

**Worked by hand.** The first sample `R = 100 ms` initializes `SRTT = 100`,
`RTTVAR = R/2 = 50`, so `RTO = 100 + 4·50 = 300 ms`. Then:

```
R = 120 ms:  RTTVAR = 0.75·50   + 0.25·|100   - 120| = 37.5   + 5      = 42.5
             SRTT   = 0.875·100 + 0.125·120          = 102.5
             RTO    = 102.5 + 4·42.5                 = 272.5 ms

R = 300 ms:  RTTVAR = 0.75·42.5  + 0.25·|102.5 - 300| = 31.875 + 49.375 = 81.25
             SRTT   = 0.875·102.5 + 0.125·300         = 127.19
             RTO    = 127.19 + 4·81.25                = 452.19 ms
```

One slow sample moves the *mean* by only 25 ms, but the *timeout* by 180 ms,
because the variance term reacts to the surprise. (Real stacks also clamp
`RTO` to a minimum: 1 s in the RFC, 200 ms on Linux.)

**AI in the wild:** the Adam optimizer (Ch. 21) is two EWMAs, one over
gradients (`β₁ = 0.9`, memory ≈ 10 steps) and one over squared gradients
(`β₂ = 0.999`, memory ≈ 1,000 steps). BatchNorm's "running mean" and
"running var" are EWMAs with momentum 0.1. Model-weight EMA, used to stabilize
diffusion models and LLM checkpoints, is the same formula applied to every
parameter.

### Real-life engineering ties

- **SWE:** SLO definitions (`p99 latency < 300ms`), capacity planning off
  percentile traffic, not average traffic.
- **AI:** every model evaluation metric is a statistic over a test set; a
  single accuracy number without variance/confidence interval is
  incomplete (Ch. 22).
- **Network:** jitter is literally the standard deviation of packet
  latency; SLA definitions for ISPs use percentile packet loss/latency.
- **Security:** baselining "normal" user/system behavior statistically is
  the foundation of anomaly-based intrusion detection — and its Achilles'
  heel (see "base rate fallacy," Ch. 16 and 26).

**More real-life examples:**
- Google's Site Reliability Engineering practice formalizes "error
  budgets" directly from percentile/SLO statistics — a team's freedom to
  ship risky changes is a function of how much of their p99-based budget
  remains this quarter.
- Monitoring systems like Prometheus and Datadog store latency as
  histograms specifically so they can compute p50/p95/p99 after the fact,
  rather than only an average — a direct product decision driven by the
  mean-vs-median problem shown above.
- Credit card fraud detection scores transactions against a customer's
  historical spending distribution and flags statistical outliers (large
  deviations from the customer's typical mean/variance) for review.
- Manufacturing and DevOps alike use "three-sigma" control charts
  (originally from Six Sigma manufacturing QA) to decide when a metric —
  defect rate or deploy failure rate — has drifted outside normal
  statistical variation and warrants investigation.

### Common mistake

Reporting only the mean/average of a skewed distribution (like latency,
which is almost always right-skewed with a long tail) — always pair it
with a median or percentile, and look at the shape, not just one number.

---

# Part 2 — Intermediate: Core Engineering Mathematics

Part 1 covered the math every engineer touches daily, almost without
noticing. This Part is where it stops being invisible: linear algebra,
calculus, and discrete math you'll reach for *by name* — in an ML paper, a
systems-design doc, a graph problem — rather than use unconsciously the way
you used arithmetic in Part 1.

## 10. Linear Algebra I — Vectors & Matrices

**Level:** Intermediate
**You'll use this for:** ML embeddings, graphics transforms, recommendation systems, and any "list of numbers as one object" problem.

### Vectors

A vector is an ordered list of numbers — a point in space, a feature list,
or a text embedding. Vectors have **magnitude** (length) and **direction**.

```
|v| = sqrt(v1² + v2² + ... + vn²)     <- this is just Pythagoras (Ch. 7) generalized
```

### Dot product

```
A · B = a1*b1 + a2*b2 + ... + an*bn
```
The dot product measures how much two vectors point in the same direction
— it's the numerator of cosine similarity (Ch. 7), and it's the single
operation a neuron performs on its inputs before an activation function
(Ch. 21).

```python
import numpy as np
a = np.array([1, 2, 3])
b = np.array([4, 5, 6])
np.dot(a, b)   # 1*4 + 2*5 + 3*6 = 32
```

### Matrices and matrix multiplication

A matrix is a grid of numbers — a batch of vectors, an image, or a linear
transformation. Multiplying a matrix by a vector *applies* that
transformation.

```
[a b] [x]   [ax + by]
[c d] [y] = [cx + dy]
```

**Rule that trips everyone up:** matrix multiplication is **not
commutative** (`AB ≠ BA` in general), and the inner dimensions must match:
`(m×n) · (n×p) = (m×p)`.

```python
A = np.array([[1, 2], [3, 4]])
B = np.array([[5, 6], [7, 8]])
A @ B          # matrix multiply (NOT A * B, which is element-wise!)
```

**Computing `A @ B` by hand, entry by entry:** each output entry is the dot
product of a *row* from `A` and a *column* from `B`.

```
A = [1 2]     B = [5 6]
    [3 4]         [7 8]

result[0][0] = (row 0 of A) · (col 0 of B) = 1×5 + 2×7 = 5+14 = 19
result[0][1] = (row 0 of A) · (col 1 of B) = 1×6 + 2×8 = 6+16 = 22
result[1][0] = (row 1 of A) · (col 0 of B) = 3×5 + 4×7 = 15+28 = 43
result[1][1] = (row 1 of A) · (col 1 of B) = 3×6 + 4×8 = 18+32 = 50

A @ B = [19 22]
        [43 50]
```

Swap the order (`B @ A`) and you'd pair *B*'s rows with *A*'s columns
instead — different numbers entirely, which is the mechanical reason
matrix multiplication isn't commutative.

**This exact bug (`*` vs. `@`/`.dot()`) is one of the most common silent
correctness bugs in numpy/PyTorch code** — `*` does element-wise
multiplication, `@` does true matrix multiplication, and both run without
error on compatible shapes, silently producing wrong numbers.

### Worked example: image as a matrix

A grayscale image is literally a 2D matrix of pixel intensities (0–255). A
color image is 3 stacked matrices (a tensor, Ch. 20) — one per RGB channel.
Rotating, blurring, and "convolutional" operations in computer vision are
all matrix/tensor operations applied to this grid.

### Worked example: a rotation matrix

Rotating a 2D point by angle `θ`:
```
[x']   [cos θ  -sin θ] [x]
[y'] = [sin θ   cos θ] [y]
```
```python
theta = math.radians(90)
R = np.array([[math.cos(theta), -math.sin(theta)],
              [math.sin(theta),  math.cos(theta)]])
R @ np.array([1, 0])   # -> [0, 1], a 90° rotation
```

### Norms: three ways to measure "how big"

```
‖x‖₁ = |x₁| + |x₂| + ... + |xₙ|               L1 / Manhattan norm
‖x‖₂ = sqrt(x₁² + x₂² + ... + xₙ²)            L2 / Euclidean norm (the |v| above)
‖x‖∞ = max(|x₁|, |x₂|, ..., |xₙ|)             L∞ / max norm
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `x ∈ ℝⁿ` | "x in R n" | a vector of `n` real numbers |
| `xᵢ` | "x sub i" | its `i`-th entry |
| `‖x‖ₚ` | "the p-norm of x" | a length. The subscript says which way of measuring |
| `\|xᵢ\|` | "absolute value of x i" | drop the sign |

**Example:** `x = (3, 4)` has `‖x‖₁ = 7`, `‖x‖₂ = 5`, `‖x‖∞ = 4`.
In a city grid you'd walk 7 blocks (L1), a crow would fly 5 (L2), and the
longest single leg of the walk is 4 (L∞).

- **AI:** L2 regularization ("weight decay") adds `λ‖w‖₂²` to the loss and
  shrinks all weights a little. L1 regularization adds `λ‖w‖₁` and drives
  many weights to exactly zero, so it selects features. Adversarial-robustness
  papers measure attack budgets in L∞: "no pixel changed by more than 8/255."
- **Network:** gradient clipping in distributed training, and "max link
  utilization" in traffic engineering, are both L∞ thinking: you care about
  the *worst* coordinate, not the total.

### The dot product, geometrically

The algebraic formula (multiply pairwise, then add) hides a geometric fact:

```
A · B = ‖A‖ · ‖B‖ · cos θ
```

| Symbol | Meaning |
|---|---|
| `‖A‖`, `‖B‖` | the L2 lengths of the two vectors |
| `θ` | the angle between them |
| `cos θ` | `1` if they point the same way, `0` if perpendicular, `-1` if opposite |

**Example:** `A = (3, 4)`, `B = (4, 3)`.

```
A · B = 3·4 + 4·3 = 24
‖A‖ = ‖B‖ = 5
cos θ = 24 / (5·5) = 0.96   ->   θ ≈ 16.3°
```

Rearranged, this *is* cosine similarity (Ch. 7). Why does the algebraic sum
equal a cosine? Expand `‖A - B‖²` both as a sum of squares and by the law of
cosines from trigonometry, and the cross term `-2·A·B` has to equal
`-2‖A‖‖B‖cos θ`.

**Projection, "how much of A lies along B":**

```
proj_B(A) = (A · B / B · B) · B  =  (24/25) · (4, 3)  =  (3.84, 2.88)
```

Projection is how PCA (Ch. 20) squashes data onto its principal axes, and
how the Gram-Schmidt process builds orthogonal bases.

### Matrix × vector = a weighted mix of columns

There are two equally correct ways to read `W·x`, and the second one is the
more useful:

```
W·x = x₁·(column 1 of W) + x₂·(column 2 of W) + ... + xₙ·(column n of W)
```

The vector `x` tells you *how much of each column to mix in*. The set of every
possible mix is called the **column space** (or **span**). It's every output
the matrix can ever produce.

**AI in the wild: one neural network layer, by hand.** A dense layer computes
`h = ReLU(W·x + b)`.

| Symbol | Shape | Meaning |
|---|---|---|
| `x` | `n × 1` | input features |
| `W` | `m × n` | weights. Row `i` is neuron `i`'s "template" |
| `b` | `m × 1` | biases, one per neuron |
| `ReLU(z)` | element-wise | `max(0, z)`: keep positives, zero out negatives |
| `h` | `m × 1` | the layer's output (activations) |

```
x = [1]     W = [0.5  -1  ]     b = [ 0.1]
    [2]         [2     0.25]        [-0.5]

W·x = [0.5·1 + (-1)·2 ] = [-1.5]       (each row dotted with x)
      [2·1   + 0.25·2 ]   [ 2.5]

W·x + b = [-1.4]     ReLU  ->   h = [0.0]
          [ 2.0]                    [2.0]
```

Neuron 1's template points away from this input, so it stays silent.
Neuron 2 matches it and fires. GPT-3 does the same thing with feed-forward
matrices of 49,152 × 12,288 and repeats it across 96 layers for every token
it generates. Inference is mostly `W·x`, which is why GPUs (matrix-multiply
machines) took over AI.

### Networking in the wild: the routing matrix

Put a network's routing into a matrix and many traffic-engineering questions
become `y = R·d`:

| Symbol | Meaning |
|---|---|
| `d` | demand vector: traffic of each flow (Mbps) |
| `R` | routing matrix: `R[link][flow] = 1` if the flow crosses that link |
| `y` | link-load vector: total traffic on each link |

Three flows over three links. Flow A uses links 1 and 2, flow B uses links 2
and 3, and flow C uses only link 3:

```
        A  B  C
R = [   1  0  0 ]  link 1          d = [10]  Mbps (A)
    [   1  1  0 ]  link 2              [20]       (B)
    [   0  1  1 ]  link 3              [ 5]       (C)

y = R·d = [10, 10+20, 20+5] = [10, 30, 25] Mbps
```

The **reverse** question is harder and more useful. SNMP gives you
link loads `y` for free, but you want the per-flow demand `d`. This is **network
tomography**, solving `R·d = y` for `d`. Here `R` is square and invertible, so
`d = R⁻¹y` recovers `[10, 20, 5]` exactly. Real networks have far more flows
than links, so `R` is wide and has no inverse. Operators then use
least squares (below) plus priors, which is what ISP traffic-matrix
estimation tools do.

### Counting paths with matrix powers

For a network's adjacency matrix `A` (Ch. 17), there's a remarkable theorem:

```
(Aᵏ)[i][j] = number of walks of exactly k hops from node i to node j
```

Four routers: `R1–R2`, `R1–R3`, `R2–R3`, `R3–R4`.

```
     R1 R2 R3 R4                  A² = [2 1 1 1]
A = [0  1  1  0]                       [1 2 1 1]
    [1  0  1  0]                       [1 1 3 0]
    [1  1  0  1]                       [1 1 0 1]
    [0  0  1  0]
```

`A²[R1][R4] = 1`: there's exactly one 2-hop walk (`R1→R3→R4`). The diagonal
`A²[i][i]` is node `i`'s degree, since every neighbour gives a walk out and
straight back. `A³[R1][R1] = 2` counts the triangle `R1→R2→R3→R1` in both
directions. Fraud-ring detection on payment graphs counts triangles with
exactly this trace-of-`A³` trick.

### Solving `A·x = b`, and the inverse

A system of linear equations is one matrix equation. For a `2 × 2` matrix the
inverse has a closed form worth knowing by heart:

```
A = [a b]       A⁻¹ = 1/(ad - bc) · [ d  -b]        x = A⁻¹·b
    [c d]                           [-c   a]
```

| Symbol | Meaning |
|---|---|
| `ad - bc` | the determinant (Ch. 11). If it's 0, divide-by-zero, no inverse |
| `A⁻¹` | the matrix that undoes `A`: `A⁻¹·A = I` |
| `I` | identity, `[[1,0],[0,1]]`, the "do nothing" matrix |

In code you almost never form `A⁻¹` explicitly. `np.linalg.solve(A, b)` is
faster and more accurate numerically (Ch. 1's floating-point issues compound
inside an explicit inverse).

### Least squares: fitting a line through noisy data

When there are more equations than unknowns, as with noisy measurements, no `x`
satisfies them all. **Least squares** picks the `w` that minimizes total
squared error, and calculus (Ch. 13: set the gradient to zero) gives a closed
form called the **normal equations**:

```
minimize  ‖X·w - y‖₂²     ==>     w = (Xᵀ X)⁻¹ Xᵀ y
```

| Symbol | Shape | Meaning |
|---|---|---|
| `X` | `n × p` | design matrix: one row per observation, a column of 1s for the intercept |
| `y` | `n × 1` | observed outputs |
| `w` | `p × 1` | the coefficients we want |
| `Xᵀ` | `p × n` | transpose of `X` |

**Worked example: latency vs. load.** Four measurements of load (thousand
req/s) and p50 latency (ms). We fit `latency = w₀ + w₁·load`:

```
load: 1, 2, 3, 4        latency: 12, 15, 19, 22

X = [1 1]     y = [12]       XᵀX = [ 4  10]      Xᵀy = [ 68]
    [1 2]         [15]             [10  30]            [187]
    [1 3]         [19]
    [1 4]         [22]

(XᵀX)⁻¹ = 1/(4·30 - 10·10) · [ 30 -10]  = 1/20 · [ 30 -10]
                             [-10   4]           [-10   4]

w₀ = (30·68 - 10·187) / 20 = 170 / 20 = 8.5 ms     (latency at zero load)
w₁ = (-10·68 + 4·187) / 20 =  68 / 20 = 3.4 ms     (per extra 1k req/s)
```

Forecast: at 6k req/s, latency ≈ `8.5 + 3.4·6 = 28.9 ms`. Be careful with that
number, because queueing (Ch. 25) bends this line upward near saturation.
A linear fit is only valid in the range you measured.

```python
X = np.array([[1, 1], [1, 2], [1, 3], [1, 4]]); y = np.array([12, 15, 19, 22])
w, *_ = np.linalg.lstsq(X, y, rcond=None)
print(w)   # [8.5 3.4]
```

Linear regression, the first model in every ML course, *is* this equation.
A neural network's last layer, trained with squared error, is solving it too.

### Real-life engineering ties

- **SWE:** any spreadsheet-like data (rows × columns) is a matrix; game/UI
  transforms (rotate, scale, translate) are matrix multiplications.
- **AI:** embeddings are vectors; a transformer's attention mechanism is,
  at its core, a sequence of matrix multiplications (Ch. 20–21); a neural
  network layer *is* `output = activation(W @ input + b)`.
- **Network:** representing a network topology as an **adjacency matrix**
  (Ch. 17) lets you use matrix operations to answer reachability questions.
- **Security:** feature vectors for ML-based malware/intrusion detection;
  representing permission matrices (subjects × objects) in access-control
  models.

**More real-life examples:**
- Every 3D game engine (Unity, Unreal) represents object position,
  rotation, and scale as 4x4 transform matrices, and renders a frame by
  multiplying a vertex through a chain of these matrices (model → view →
  projection) — a direct, large-scale use of Ch. 10's rotation-matrix idea.
- Search engines historically represented documents as TF-IDF vectors
  (one dimension per vocabulary word) and ranked relevance using vector
  operations like the dot product and cosine similarity from Ch. 7.
- Recommendation engines (Netflix, Amazon) represent the entire
  user-item interaction history as one giant sparse matrix, and factorize
  it (Ch. 20) to predict missing entries — "would this user like this item."
- Spreadsheet formulas that reference ranges (`SUMPRODUCT`, array formulas
  in Excel/Sheets) are literally performing vector/matrix operations under
  a friendlier UI.

### Common mistake

Confusing `*` (element-wise) with `@`/`np.dot`/`np.matmul` (true matrix
multiply) in numpy/PyTorch — both are legal, only one is usually correct,
and there's no error to warn you.

---

## 11. Linear Algebra II — Eigenvalues, Eigenvectors, Determinants

**Level:** Intermediate
**You'll use this for:** PageRank, PCA, graph/network centrality, and stability analysis.

### The idea

For a matrix `A`, an **eigenvector** `v` is a special vector that, when
transformed by `A`, only gets *scaled* (not rotated): `A·v = λ·v`. The
scalar `λ` is the **eigenvalue**. Intuitively: eigenvectors are the "natural
axes" of a transformation — the directions that don't get twisted.

```python
A = np.array([[2, 0], [0, 3]])
eigenvalues, eigenvectors = np.linalg.eig(A)
print(eigenvalues)   # [2. 3.]
```

That example is diagonal, which makes the eigenvalues visible by
inspection (2 and 3, right on the diagonal). Here's the general by-hand
method — the **characteristic equation** — for a matrix where they aren't
obvious, `A = [[2,1],[1,2]]`:

```
Start from the definition: A·v = λ·v  =>  (A - λI)·v = 0

For a non-zero v to solve this, (A - λI) must be non-invertible, i.e.
det(A - λI) = 0.

A - λI = [2-λ   1  ]
         [1     2-λ]

det(A - λI) = (2-λ)(2-λ) - (1)(1)
            = (2-λ)² - 1
            = 4 - 4λ + λ² - 1
            = λ² - 4λ + 3

Set it to 0 and factor:  λ² - 4λ + 3 = (λ-1)(λ-3) = 0
  ->  λ = 1  or  λ = 3
```

Every eigenvalue computation is this same move: build `A - λI`, take its
determinant (a formula in `λ`), and solve where that formula equals zero.
For a `2×2` matrix that's a quadratic, solvable by hand; for larger
matrices it becomes a higher-degree polynomial, which is why in practice
you use `np.linalg.eig` beyond toy sizes — but the mechanism is identical.

### Determinant

A single number summarizing a matrix's transformation: how much it scales
area/volume, and whether it flips orientation. `det(A) = 0` means the
matrix "collapses" space into a lower dimension — it's not invertible
(no unique solution exists for `Ax = b`).

```python
np.linalg.det(np.array([[1, 2], [3, 4]]))   # -2.0
np.linalg.det(np.array([[1, 2], [2, 4]]))   # 0.0 -> singular, not invertible
```

### Worked example: PageRank (and, more generally, graph centrality)

PageRank models a "random surfer" clicking links forever. The long-run
probability of being on each page is the **dominant eigenvector** of the
web's link matrix (normalized so columns sum to 1). This is the actual
mathematical mechanism behind Google's original ranking algorithm, and the
same idea (eigenvector centrality) is used today to find the most
"influential" node in *any* graph — social networks, citation networks, or
an internal service dependency graph.

```python
# Toy 3-page web: A links to B and C, B links to C, C links to A and B
M = np.array([[0,   0,   0.5],
              [0.5, 0,   0.5],
              [0.5, 1,   0  ]])
eigvals, eigvecs = np.linalg.eig(M)
dominant = eigvecs[:, np.argmax(eigvals)]
rank = np.abs(dominant) / np.sum(np.abs(dominant))
print(rank)   # relative importance of pages A, B, C
```

### The eigen-equation, symbol by symbol

```
A · v = λ · v          (v ≠ 0)
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `A` | "A" | a square `n × n` matrix, a transformation |
| `v` | "v", the eigenvector | a direction that `A` does not turn |
| `λ` | "lambda", the eigenvalue | how much `A` stretches that direction. `λ > 1` grows, `0 < λ < 1` shrinks, `λ < 0` flips |
| `v ≠ 0` | | the zero vector trivially satisfies the equation, so it doesn't count |

### Finding the eigenvectors, by hand

We found `λ = 1` and `λ = 3` for `A = [[2,1],[1,2]]` above. To get each
eigenvector, plug `λ` back into `(A - λI)·v = 0` and solve:

```
λ = 3:  A - 3I = [-1  1]   ->  -v₁ + v₂ = 0  ->  v₁ = v₂   ->  v = (1, 1)
                 [ 1 -1]

λ = 1:  A - 1I = [ 1  1]   ->   v₁ + v₂ = 0  ->  v₁ = -v₂  ->  v = (1, -1)
                 [ 1  1]
```

**Check:** `A·(1,1) = (2+1, 1+2) = (3,3) = 3·(1,1)`. ✓
Geometrically, `A` stretches the diagonal `(1,1)` by 3× and leaves the
anti-diagonal `(1,-1)` alone. Every other vector is a mix of those two, so
it gets partly stretched and partly left alone, which turns it.

Two theorems give you free sanity checks:

```
trace(A) = sum of eigenvalues        2 + 2 = 4 = 1 + 3   ✓
det(A)   = product of eigenvalues    2·2 - 1·1 = 3 = 1·3 ✓
```

### Diagonalization: why eigenvalues predict the future

If `A` has `n` independent eigenvectors, stack them as the columns of `P` and
the eigenvalues on the diagonal of `D`:

```
A = P · D · P⁻¹          and therefore          Aᵏ = P · Dᵏ · P⁻¹
```

| Symbol | Meaning |
|---|---|
| `P` | columns are the eigenvectors, a change into the "natural axes" |
| `D` | diagonal matrix of eigenvalues, so `A` is pure stretching in those axes |
| `Dᵏ` | trivial to compute: raise each diagonal entry to the `k`-th power |

Applying `A` a thousand times is hard. Raising each eigenvalue to the
thousandth power is easy. For our `A`:

```
A¹⁰ = P · [1¹⁰   0  ] · P⁻¹ = [ (3¹⁰+1)/2   (3¹⁰-1)/2 ] = [29525  29524]
          [0    3¹⁰ ]         [ (3¹⁰-1)/2   (3¹⁰+1)/2 ]   [29524  29525]
```

**The long-run behaviour of *any* repeated linear process is governed by its
largest eigenvalue.** If `|λ_max| > 1` it blows up, if `< 1` it dies out, and
if `= 1` it settles into a steady state. That one fact explains three very
different things:

- **AI, exploding/vanishing gradients in RNNs:** backprop through `T` time steps
  multiplies by the recurrent weight matrix `T` times, roughly `Wᵀ`. If
  `|λ_max(W)| = 1.1` and `T = 100`, gradients scale by `1.1¹⁰⁰ ≈ 13,780`. If
  it's `0.9`, they scale by `0.9¹⁰⁰ ≈ 0.00003`. LSTMs, gradient clipping and
  orthogonal initialization (all eigenvalues exactly magnitude 1) are each a fix
  for this.
- **Network, PageRank and Markov chains:** steady state = the eigenvector for
  `λ = 1` (below).
- **Control, feedback-loop stability:** an autoscaler or congestion
  controller is stable only if its update matrix has every eigenvalue inside
  the unit circle (Ch. 29).

### The spectral theorem: symmetric matrices are the friendly ones

> **Theorem.** If `A` is real and symmetric (`A = Aᵀ`), then all its
> eigenvalues are real, and it has a full set of **orthogonal** eigenvectors,
> so `A = Q·D·Qᵀ` with `Qᵀ = Q⁻¹`.

This matters because the most important matrices in practice are symmetric:
covariance matrices (PCA, Ch. 20), Hessians (Ch. 13, 21), graph Laplacians
(below), and kernel matrices. For all of them the eigenvalues are real numbers
you can sort and interpret. That's why the PCA code in Ch. 20 uses
`np.linalg.eigh`, the solver specialized for symmetric matrices.

### Power iteration, and PageRank properly

You rarely need *all* eigenvalues, just the dominant one. **Power iteration**
finds it by repeatedly multiplying and re-normalizing, because the
largest-eigenvalue component outgrows the rest (that's the diagonalization
argument above). Google's real PageRank adds a **damping factor**:

```
PR = (1 - d)/N · 1  +  d · M · PR
```

| Symbol | Meaning |
|---|---|
| `PR` | vector of page scores (sums to 1) |
| `M` | column-stochastic link matrix: `M[i][j] = 1/outlinks(j)` if `j` links to `i` |
| `d` | damping, ≈ 0.85: the probability the surfer follows a link instead of jumping to a random page |
| `N` | number of pages |
| `1` | a vector of all ones, the "teleport anywhere" term |

Toy web from above, `d = 0.85`, start uniform at `[1/3, 1/3, 1/3]`:

```
iter 1:  [0.192, 0.333, 0.475]
iter 2:  [0.252, 0.333, 0.415]
iter 3:  [0.226, 0.333, 0.440]
...
steady:  [0.234, 0.333, 0.433]     <- C ranks highest: both A and B link to it
```

The teleport term does two jobs. It guarantees a unique answer even when the
web has dead ends and closed loops (that's the Perron–Frobenius theorem: a
positive stochastic matrix has a unique dominant eigenvector with `λ = 1`).
It also makes convergence fast, because the error shrinks by a factor of
about `d = 0.85` per iteration, so 50 iterations reach about `0.85⁵⁰ ≈ 0.0003`.

### Markov chains: bursty packet loss, modelled

A **Markov chain** is a system that hops between states with fixed
probabilities. The transition matrix `P[i][j]` is the probability of moving
from state `i` to `j`, and each row sums to 1. In the long run the fraction
of time spent in each state is the **stationary distribution** `π`:

```
π · P = π        (π is a left eigenvector of P with eigenvalue 1)
Σᵢ πᵢ = 1
```

**Networking in the wild: the Gilbert–Elliott channel.** Real links don't
drop packets independently (Ch. 8's coin-flip model). They drop them in
*bursts*, during Wi-Fi interference or a congested queue. The classic model
has two states:

```
                      to Good   to Bad
P =  from Good   [    0.99      0.01  ]        Good: 0% loss
     from Bad    [    0.30      0.70  ]        Bad:  50% loss
```

Solve `π·P = π` for two states. Flow into Bad has to equal flow out of Bad,
so `π_G · 0.01 = π_B · 0.30`. Combined with `π_G + π_B = 1`:

```
π_B = 0.01 / (0.01 + 0.30) = 0.032     (the link is "Bad" 3.2% of the time)
average loss = 0.032 · 50% ≈ 1.6%
mean burst length = 1 / 0.30 ≈ 3.3 packets in a row in the Bad state
```

The long-run loss rate is the same 1.6% as an independent-loss model would
give, but the losses arrive in clumps. That's why forward-error-correction
schemes (used in video calls and QUIC experiments) **interleave** packets:
it spreads one burst across many FEC blocks so each block loses only one packet.
Netem, Linux's network emulator, ships this exact model
(`tc qdisc ... loss gemodel`).

The same `π·P = π` equation is the language-model view of text too. An n-gram
model is a Markov chain over words, and the "random surfer" of PageRank is a
Markov chain over web pages.

### The graph Laplacian: measuring how hard a network is to cut

For a graph with adjacency matrix `A` and degree matrix `D` (node degrees on
the diagonal), the **Laplacian** is:

```
L = D - A
```

Its eigenvalues are all `≥ 0`, and the smallest is always 0. The
**second-smallest eigenvalue `λ₂`** (the *algebraic connectivity* or *Fiedler
value*) measures how well-connected the graph is. `λ₂ = 0` exactly when the
network is already split in two, and a bigger `λ₂` means it takes more link
cuts to partition it.

Four routers, three topologies:

| Topology | Links | `λ₂` | Read it as |
|---|---|---|---|
| Line `R1–R2–R3–R4` | 3 | **0.59** | one cut in the middle splits it |
| Ring (add `R4–R1`) | 4 | **2.00** | needs two cuts to split |
| Full mesh | 6 | **4.00** | needs three cuts to isolate any node |

```python
def fiedler(edges, n=4):
    A = np.zeros((n, n))
    for i, j in edges:
        A[i, j] = A[j, i] = 1
    L = np.diag(A.sum(axis=1)) - A
    return np.sort(np.linalg.eigvalsh(L))[1]

print(fiedler([(0,1), (1,2), (2,3)]))          # 0.586  line
print(fiedler([(0,1), (1,2), (2,3), (3,0)]))   # 2.0    ring
```

Adding **one** link (line → ring) more than triples `λ₂`. That's a
quantitative answer to "which single link should we add for the most
resilience?": try each candidate and pick the largest jump in `λ₂`. The
eigenvector for `λ₂` (the *Fiedler vector*) has positive entries on one side of
the network's weakest cut and negative on the other. **Spectral clustering**
in ML and community detection in social graphs both partition data with that
sign pattern.

### Real-life engineering ties

- **SWE:** recommendation systems, search ranking.
- **AI:** Principal Component Analysis (PCA, Ch. 20) is literally "find the
  eigenvectors of the covariance matrix"; eigenvalues of a Hessian
  determine whether a training loss landscape point is a minimum, maximum,
  or saddle point.
- **Network:** identifying the most critical router/node in a topology
  (eigenvector/betweenness centrality) for redundancy planning; spectral
  graph analysis for detecting network partitions.
- **Security:** anomaly detection via PCA (projecting traffic onto the top
  eigenvectors and flagging what doesn't fit); social-graph analysis for
  identifying influential/central accounts in fraud rings.

**More real-life examples:**
- Academic citation-ranking metrics (Eigenfactor, and the ideas behind
  Google Scholar) rank journals/papers using the same eigenvector-
  centrality math as PageRank, applied to a citation graph.
- Social network platforms use graph centrality measures (eigenvector,
  PageRank-style, or betweenness centrality) to identify influential
  accounts for recommendation, ranking, and bot-network detection.
- Network engineers use centrality analysis on a topology graph to find
  single points of failure — the router/link whose removal would fragment
  the network the most — directly informing redundancy investment.
- Structural engineers compute the eigenvalues of a system's stiffness
  matrix to find its resonant frequencies — the 1940 Tacoma Narrows Bridge
  collapse is a famous real-world illustration of why eigenvalue analysis
  of physical systems matters.

### Common mistake

Assuming every matrix has "nice" real eigenvalues — many don't (some are
complex), and non-square matrices don't have eigenvalues at all (you need
Singular Value Decomposition instead, Ch. 20).

---

## 12. Calculus I — Limits, Derivatives, Rates of Change

**Level:** Intermediate
**You'll use this for:** understanding "rate of change" metrics, and as the prerequisite for all of ML optimization.

### The idea

A **derivative** measures instantaneous rate of change — how fast `f(x)` is
changing *right now*, not on average. If `f(x)` is position, `f'(x)` is
velocity.

```
f'(x) = lim(h→0) [f(x+h) - f(x)] / h
```

### Watching that limit converge, numerically

For `f(x) = x²` at `x = 3`, plug in shrinking values of `h` and watch
`[f(x+h) - f(x)] / h` settle down:

```
h = 1:      [f(3+1) - f(3)] / 1     = [16 - 9] / 1     = 7.0
h = 0.1:    [f(3.1) - f(3)] / 0.1   = [9.61 - 9] / 0.1  = 6.1
h = 0.01:   [f(3.01) - f(3)] / 0.01 = [9.0601-9]/0.01  = 6.01
h = 0.001:  ...                                          = 6.001

As h -> 0, the ratio converges to 6.
```

The power rule says `f'(x) = 2x`, so `f'(3) = 2×3 = 6` — matching exactly
what the shrinking-`h` calculation converges to. That's what the "limit"
in the definition means in practice: it's not a mysterious symbol, it's
"keep shrinking `h` and see what number the ratio approaches."

### Common derivative rules

```
d/dx (x^n) = n·x^(n-1)          power rule
d/dx (e^x) = e^x                 e^x is its own derivative
d/dx (ln x) = 1/x
d/dx (f(x)+g(x)) = f'(x) + g'(x) sum rule
d/dx (f(x)·g(x)) = f'(x)g(x) + f(x)g'(x)   product rule
```

### Worked example: marginal cost / rate-of-change intuition

If `cost(n) = 0.01·n² + 5n` is the cost of serving `n` requests/sec of
extra compute capacity, `cost'(n) = 0.02n + 5` tells you the *marginal*
cost of the *next* unit of capacity — not the average cost so far. This
distinction (marginal vs. average) is exactly the reasoning behind capacity
planning decisions like "is it worth adding one more node."

```python
import sympy as sp
n = sp.symbols('n')
cost = 0.01*n**2 + 5*n
marginal = sp.diff(cost, n)
print(marginal)             # 0.02*n + 5
print(marginal.subs(n, 100))  # marginal cost at n=100 requests/sec
```

### The derivative, symbol by symbol

```
f'(x) = lim(h→0) [f(x+h) - f(x)] / h        also written  df/dx,  d/dx f(x)
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `f'(x)` | "f prime of x" | the slope of `f` at the point `x` |
| `h` | "h" | a small step to the right of `x` |
| `f(x+h) - f(x)` | "rise" | how much the output changed over that step |
| `… / h` | "rise over run" | the average slope over the step, a *secant* line |
| `lim(h→0)` | "the limit as h goes to zero" | the value the average slope settles on as the step shrinks, the *tangent* line |
| `df/dx` | "d f d x" (Leibniz notation) | the same thing, written as "a tiny change in f per tiny change in x" |

**Units check:** if `f` is bytes and `x` is seconds, `f'` is bytes/second.
A derivative always has units of *output per input*.

### Deriving a rule instead of memorizing it

Here's the power rule for `f(x) = x²`, straight from the definition:

```
[f(x+h) - f(x)] / h = [(x+h)² - x²] / h
                    = [x² + 2xh + h² - x²] / h
                    = [2xh + h²] / h
                    = 2x + h          ->  as h → 0:   f'(x) = 2x
```

The `h²` term was "too small to matter", and that's the entire spirit of
calculus: zoom in far enough and every smooth curve looks like a straight
line.

Two more rules you'll need, plus the chain rule (Ch. 13 develops it fully):

```
d/dx [f(x)/g(x)] = [f'(x)·g(x) - f(x)·g'(x)] / g(x)²     quotient rule
d/dx f(g(x))     = f'(g(x)) · g'(x)                        chain rule
```

### AI in the wild: the sigmoid's derivative and vanishing gradients

The sigmoid squashes any number into a probability `(0, 1)`. It's the output
of every logistic regression and every binary classifier:

```
σ(x) = 1 / (1 + e^(-x))
```

Differentiate with the chain rule. The outer function is `u⁻¹`, the inner is
`1 + e^(-x)`:

```
σ'(x) = -1·(1 + e^(-x))⁻² · (-e^(-x))
      = e^(-x) / (1 + e^(-x))²
      = [1/(1 + e^(-x))] · [e^(-x)/(1 + e^(-x))]
      = σ(x) · (1 - σ(x))
```

The derivative is built from the function's own output, so backprop gets it
almost for free. It also has a dark side:

| `x` | `σ(x)` | `σ'(x)` |
|---|---|---|
| 0 | 0.5 | **0.25** (the maximum possible) |
| 4 | 0.982 | 0.018 |
| -4 | 0.018 | 0.018 |

Each sigmoid layer multiplies the gradient by **at most 0.25**. Through 10
layers that's `0.25¹⁰ ≈ 0.000001`, so the early layers barely learn. This
one-line derivative is *the* reason deep networks switched to ReLU, whose
derivative is exactly 1 for positive inputs (Ch. 21).

### Taylor's theorem: every smooth function is secretly a polynomial

> **Theorem (Taylor).** Near a point `a`, a smooth function equals its
> derivatives arranged as a polynomial:
> ```
> f(x) = f(a) + f'(a)(x-a) + f''(a)/2! · (x-a)² + f'''(a)/3! · (x-a)³ + ...
> ```

| Symbol | Meaning |
|---|---|
| `a` | the point you expand around, where you know everything |
| `x - a` | how far you've moved from it |
| `f''(a)` | second derivative: the curvature (how fast the slope changes) |
| `n!` | factorial, `n·(n-1)·…·1`, which shrinks the higher terms fast |

Keep only the first two terms and you get the **linear approximation**
`f(x) ≈ f(a) + f'(a)(x-a)`. That's the "zoom in until it's a line" idea as a
formula.

**Examples you've already relied on without knowing it:**

```
e^x       ≈ 1 + x + x²/2 + x³/6      at x = 0.1:  1.1051667  (true: 1.1051709)
ln(1 + x) ≈ x                         at x = 0.01: 0.01       (true: 0.00995)
```

The second one *derives* Ch. 6's Rule of 72. Money growing at rate `r`
doubles when `(1+r)ᵗ = 2`, so `t = ln 2 / ln(1+r) ≈ 0.693 / r`. That gives
"69.3 divided by the percentage rate", and 72 is used instead because it has
more divisors.

**Where Taylor shows up in AI:** gradient descent (Ch. 21) *is* the
first-order Taylor approximation. It assumes the loss is locally a plane and
steps downhill on it. Newton's method and second-order optimizers keep the
`f''` term and model the loss as a bowl. The GELU activation used in GPT and
BERT ships with a `tanh`-based approximation that's tuned to match the true
curve closely and runs faster on GPUs.

### Newton's method: square roots in four steps

To solve `f(x) = 0`, start with a guess, take the tangent line (the
linear approximation above), and jump to where *it* hits zero:

```
x_{n+1} = x_n - f(x_n) / f'(x_n)
```

**Example: compute `√2`**, which is the root of `f(x) = x² - 2` with `f'(x) = 2x`.
The update simplifies to `x_{n+1} = (x_n + 2/x_n) / 2`:

```
x₀ = 1
x₁ = (1 + 2/1) / 2           = 1.5
x₂ = (1.5 + 2/1.5) / 2       = 1.41666...
x₃ = (1.41667 + 2/1.41667)/2 = 1.4142157
x₄                           = 1.41421356237469    (true: 1.41421356237310)
```

The number of correct digits roughly **doubles** every step (1, 3, 6, 12).
That's called *quadratic convergence*. Your CPU's square-root and division
units, the famous Quake III "fast inverse square root" hack, and the
implied-volatility solvers in finance all finish with Newton steps.

### Optimization: where the derivative is zero

At a smooth minimum or maximum the tangent is flat, so `f'(x) = 0`. The
second derivative tells you which one you found:

```
f'(x*) = 0  and  f''(x*) > 0   ->  local minimum  (curves up, like a bowl)
f'(x*) = 0  and  f''(x*) < 0   ->  local maximum  (curves down, like a hill)
```

**Example: right-sizing a fleet.** Each server costs $2/hour, and the
latency penalty to the business falls as you add servers, about `800/n`
dollars/hour:

```
C(n)   = 2n + 800/n
C'(n)  = 2 - 800/n²   = 0   ->   n² = 400   ->   n* = 20 servers
C''(n) = 1600/n³ > 0                          ->   minimum ✓

C(10) = $100/h     C(20) = $80/h     C(40) = $100/h
```

At the optimum, the **marginal cost of one more server equals the marginal
saving from it** (`2 = 800/n²`). Every "find the sweet spot" engineering
trade-off has this shape: cache size vs. hit rate, batch size vs. latency,
replication factor vs. durability.

### Networking in the wild: how sensitive is TCP to packet loss?

The Mathis formula approximates steady-state throughput for a loss-based
TCP (Reno-style) flow:

```
throughput ≈ (MSS / RTT) · (C / √p)
```

| Symbol | Meaning | Example |
|---|---|---|
| `MSS` | maximum segment size (bytes per packet) | 1460 bytes |
| `RTT` | round-trip time | 100 ms |
| `p` | packet-loss probability | 0.0001 (0.01%) |
| `C` | a constant from the sawtooth geometry (Ch. 25), `√(3/2) ≈ 1.22` | |

```
throughput = (1460·8 bits / 0.1 s) · (1.22 / √0.0001)
           = 116,800 · 122.5  ≈ 14.3 Mbit/s
```

Now ask the calculus question: how sensitive is throughput to loss?

```
d(throughput)/dp = -½ · throughput / p
```

so a **1% relative increase in loss costs about 0.5% of throughput**.
Doubling the loss rate divides throughput by `√2` (−29%, down to
10.1 Mbit/s). Going to 1% loss cuts it 10× to 1.4 Mbit/s. And because RTT
sits in the denominator, the same loss rate hurts a transatlantic flow 10×
more than a metro one. These two derivatives are why loss-based congestion
control struggles on long, slightly lossy paths, and why Google built BBR,
which models bandwidth and RTT directly instead of reacting to loss.

### Integrals: the accumulated total

If the derivative answers "how fast right now?", the **integral** answers
"how much in total?". It adds up a rate over time, or a density over a
range, into an amount.

```
∫ₐᵇ f(x) dx  =  lim(n→∞) Σᵢ₌₁ⁿ f(xᵢ) · Δx        where Δx = (b - a)/n
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `∫` | "integral" | a stretched "S" for **S**um |
| `a`, `b` | "from a to b" | the start and end of the range |
| `f(x)` | "the integrand" | the rate or height being added up |
| `dx` | "d x" | an infinitely thin slice of the range. Its units multiply in |
| `Δx` | "delta x" | the width of one finite slice |
| `Σ f(xᵢ)·Δx` | Riemann sum | chop the area into `n` rectangles, height × width, and add them up |

**Units check:** if `f` is bytes/second and `x` is seconds, the integral
is bytes. Integration multiplies units by the units of `dx`, and
differentiation divides them.

**Watching the rectangles converge**, for `∫₀³ x² dx` using midpoint rectangles:

```
n = 3 rectangles:    8.75
n = 6:               8.9375
n = 30:              8.9975
n = 300:             8.999975       ->  converges to exactly 9
```

Same idea as the derivative's shrinking `h`: the "limit" just means "keep
slicing thinner and see what the sum approaches."

### The Fundamental Theorem of Calculus

The theorem that welded derivatives and integrals into one subject (Newton
and Leibniz, 1660s–1680s):

> **Theorem (FTC).** If `F' = f` (that is, `F` is an *antiderivative* of `f`), then
> ```
> ∫ₐᵇ f(x) dx = F(b) - F(a)
> ```
> And in the other direction: `d/dx ∫ₐˣ f(t) dt = f(x)`. Accumulating a rate,
> then asking how fast the total is growing, gives you back the rate.

Integration and differentiation are inverses, so every derivative rule
read backwards is an integral rule:

```
∫ xⁿ dx   = xⁿ⁺¹/(n+1) + C      (n ≠ -1)
∫ eˣ dx   = eˣ + C
∫ 1/x dx  = ln|x| + C
∫ e^(-λt) dt = -e^(-λt)/λ + C
```

(`C` is the "constant of integration". Constants have zero derivative, so an
antiderivative is only fixed up to a constant, and it cancels in `F(b) - F(a)`.)

**Example:** `∫₀³ x² dx = [x³/3]₀³ = 27/3 - 0 = 9`. That's the 9 the rectangles
were crawling toward, found exactly in one line.

### Networking in the wild: rates, totals, and Prometheus

Your monitoring stack is a calculus engine:

| Prometheus function | What it computes | Calculus |
|---|---|---|
| `rate(bytes_total[5m])` | bytes/second from an ever-growing counter | derivative of the counter |
| `increase(bytes_total[1h])` | bytes transferred in the last hour | integral of the rate |
| `deriv(queue_depth[10m])` | how fast a gauge is changing | derivative, by least-squares slope (Ch. 10) |

**Example 1, exact:** a transfer ramps up linearly, `r(t) = 2t` Gbit/s, for 10
seconds. Total data:

```
∫₀¹⁰ 2t dt = [t²]₀¹⁰ = 100 Gbit = 12.5 GB
```

Its *average* rate is 10 Gbit/s, half the 20 Gbit/s peak. If you capacity-plan a
link off the average, you'll saturate it at the end of every ramp.

**Example 2, from samples (the trapezoid rule):** you only have
scraped rates, one every 15 s: `100, 140, 180, 160, 120` Mbit/s. Approximate
each 15 s slice as a trapezoid (average of its two ends × width):

```
total ≈ 15 · [ (100+140)/2 + (140+180)/2 + (180+160)/2 + (160+120)/2 ]
      = 15 · [ 120 + 160 + 170 + 140 ]
      = 8,850 Mbit  ≈  1,106 MB  in one minute
```

That's how volume-based billing computes "total GB transferred" from sampled
rates. (Transit billed on the 95th percentile uses Ch. 9's percentiles
instead.) A **token-bucket rate limiter** is the same
integral run live: the bucket holds `∫ refill_rate dt - consumed`, capped at the
bucket size.

### AI in the wild: probability is integration

For a continuous random variable with density `f(x)` (Ch. 14), probabilities
*are* areas:

```
P(a ≤ X ≤ b) = ∫ₐᵇ f(x) dx            ∫₋∞^∞ f(x) dx = 1            E[X] = ∫ x·f(x) dx
```

**Example: the exponential distribution** (time to next failure, rate `λ`):

```
f(t) = λ·e^(-λt)        for t ≥ 0

Total probability:  ∫₀^∞ λe^(-λt) dt = [-e^(-λt)]₀^∞ = 0 - (-1) = 1   ✓

P(failure within T) = ∫₀ᵀ λe^(-λt) dt = 1 - e^(-λT)
```

With MTBF 720 hours (`λ = 1/720`) and `T = 100` h:
`1 - e^(-100/720) = 1 - e^(-0.139) ≈ 13%`. That's exactly what `expon.cdf` in
Ch. 14 returns, because the CDF *is* this integral.

**The Gaussian integral, a theorem that explains a constant:**

```
∫₋∞^∞ e^(-x²) dx = √π
```

(The proof squares the integral, switches to polar coordinates, and the
`π` falls out of the circle. It's one of the most elegant tricks in
mathematics.) This is *why* the normal distribution's formula carries that odd
`1/√(2π)`: it's the constant needed to make the total area exactly 1.

**ROC AUC is literally an integral.** The "area under the ROC curve" used to
score every classifier and intrusion detector (Ch. 26) is
`∫₀¹ TPR d(FPR)`, computed by the trapezoid rule over the curve's points.
For points `(0,0), (0.1,0.6), (0.3,0.85), (1,1)`:

```
AUC = 0.1·(0+0.6)/2 + 0.2·(0.6+0.85)/2 + 0.7·(0.85+1)/2
    = 0.030 + 0.145 + 0.6475 = 0.8225
```

```python
from scipy.integrate import quad
import numpy as np
print(quad(lambda x: x**2, 0, 3))                       # (9.0, ...)  area under x²
print(quad(lambda x: np.exp(-x**2), -np.inf, np.inf))   # (1.7724..., ...) = √π
print(np.trapezoid([0, 0.6, 0.85, 1], [0, 0.1, 0.3, 1]))     # 0.8225  ROC AUC
```

### Real-life engineering ties

- **SWE:** understanding "rate of change" dashboards (deploy velocity,
  error rate acceleration).
- **AI:** the derivative is literally what gradient descent follows
  downhill (Ch. 13, 21) — this chapter is the prerequisite for all of
  neural network training.
- **Network:** rate of change of queue depth or throughput signals
  congestion before it becomes an outage (used in AQM/CoDel algorithms,
  Ch. 25).
- **Security:** rate-of-change of failed login attempts (acceleration, not
  just count) is a stronger brute-force signal than a raw threshold.

**More real-life examples:**
- Cloud cost/capacity dashboards that show "spend rate" and "rate of
  change of spend" are literally reporting a metric and its derivative —
  a sudden jump in the *derivative* is a much earlier warning sign than
  watching total spend alone.
- Physics engines in games and robotics simulations integrate
  acceleration to get velocity and velocity to get position every frame —
  the same derivative chain (position → velocity → acceleration) run
  continuously at 60+ times a second.
- Newton's method, a numerical technique that uses a function's
  derivative to iteratively find where it crosses zero, is used inside
  many optimizers, root-finders, and implied-volatility calculations in
  finance.
- SRE teams track error-rate *velocity* (how fast the error rate itself
  is climbing) as an earlier, more sensitive incident-detection signal
  than a static error-rate threshold — a derivative-based alert, not a
  level-based one.

### Common mistake

Confusing a **rate** with a **total**. "Requests per second is climbing" is
a statement about the derivative; "we've served 1M requests today" is a
statement about the integral (the accumulated total, Ch. 13). Dashboards
that don't clearly label which one you're looking at cause real
misdiagnosis during incidents.

---

## 13. Calculus II — Gradients, Partial Derivatives, the Chain Rule

**Level:** Intermediate
**You'll use this for:** the entire mathematical basis of how neural networks learn (backpropagation, Ch. 21).

### Partial derivatives

When a function has multiple inputs, `∂f/∂x` measures the rate of change
with respect to `x`, holding every other variable fixed. For
`f(x, y) = x²y + y³`:

```
∂f/∂x = 2xy          (treat y as a constant)
∂f/∂y = x² + 3y²      (treat x as a constant)
```

**Getting there term by term.** `f(x,y) = x²y + y³` has two additive
terms, so differentiate each separately (sum rule, Ch. 12) and add:

```
∂f/∂x  (treat y as a fixed number, like "5"):
  d/dx (x²y) = y · d/dx(x²) = y · 2x = 2xy    <- power rule on x, y is just a constant multiplier
  d/dx (y³)  = 0                               <- no x in this term at all, so its rate of change w.r.t. x is 0
  sum: ∂f/∂x = 2xy + 0 = 2xy

∂f/∂y  (treat x as a fixed number):
  d/dy (x²y) = x² · d/dy(y) = x² · 1 = x²      <- x² is just a constant multiplier
  d/dy (y³)  = 3y²                              <- power rule on y
  sum: ∂f/∂y = x² + 3y²
```

Every partial derivative is an ordinary single-variable derivative in
disguise — you just freeze every other variable as if it were a plain
number first.

### The gradient

The **gradient** `∇f` is the vector of all partial derivatives — it points
in the direction of **steepest ascent** of `f`. Negate it, and you get the
direction of steepest *descent* — the entire idea behind gradient descent
(Ch. 21).

```python
import sympy as sp
x, y = sp.symbols('x y')
f = x**2 * y + y**3
gradient = [sp.diff(f, x), sp.diff(f, y)]
print(gradient)   # [2*x*y, x**2 + 3*y**2]
```

### The chain rule (the single most important rule in ML math)

If `y = f(g(x))`, then `dy/dx = f'(g(x)) · g'(x)`. In words: "the rate of
change of the outer function times the rate of change of the inner
function." This is the *entire* mathematical mechanism behind
backpropagation — a neural network is nothing but nested functions
(layer₁(layer₂(layer₃(...)))), and computing how the final loss changes
with respect to an early-layer weight requires applying the chain rule
across every layer in between.

```python
# y = (3x + 1)^2  =>  outer: u^2, inner: u = 3x+1
x = sp.symbols('x')
y = (3*x + 1)**2
print(sp.diff(y, x))          # 6*(3*x + 1)  -- sympy applies chain rule automatically
```

### Worked example: chain rule through two tiny "layers"

```
layer1: h = w1 * x
layer2: y = w2 * h
loss:   L = (y - target)^2
```
To update `w1`, you need `∂L/∂w1`, computed by chaining backwards:
```
∂L/∂w1 = ∂L/∂y · ∂y/∂h · ∂h/∂w1
```
This is precisely backpropagation — computing derivatives backward through
a computation graph, one chain-rule link at a time (fully worked out with
numbers in Ch. 21).

### The gradient, symbol by symbol

```
∇f(x) = [ ∂f/∂x₁,  ∂f/∂x₂,  ...,  ∂f/∂xₙ ]
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `∇` | "nabla" or "del" | the gradient operator: "take every partial derivative" |
| `∂` | "partial" | a curly `d` that warns you other variables exist and are being held fixed |
| `∂f/∂xᵢ` | "partial f partial x i" | the slope of `f` if you nudge only `xᵢ` |
| `∇f(x)` | "grad f at x" | a vector with one entry per input, the same shape as `x` |

**Example:** for `f(x,y) = x²y + y³` at the point `(1, 2)`:

```
∇f = (2xy, x² + 3y²) = (2·1·2, 1 + 3·4) = (4, 13)
```

The output is about 3× more sensitive to `y` than to `x` here. In ML terms,
`y`'s weight would get a 3× bigger update.

### Why the gradient points uphill (a two-line proof)

The slope in an arbitrary unit direction `u` is the **directional
derivative**:

```
D_u f = ∇f · u = ‖∇f‖ · ‖u‖ · cos θ = ‖∇f‖ · cos θ
```

(This uses the geometric dot product from Ch. 10, with `‖u‖ = 1`.) `cos θ` is
largest (= 1) when `θ = 0`, that is, when `u` points exactly along `∇f`. So
**the gradient is the direction of steepest ascent, and its length is
that steepest slope**. `-∇f` is steepest descent. Perpendicular to `∇f`
(`cos θ = 0`) the slope is zero, which is the direction of a contour line on
a map.

**Example:** at `(1, 2)` with `∇f = (4, 13)`, moving in direction
`u = (0.6, 0.8)` gives slope `4·0.6 + 13·0.8 = 12.8`. Moving along the
gradient itself gives `‖∇f‖ = √(16 + 169) ≈ 13.6`, which is steeper, as the
proof promised.

### The Jacobian: derivatives of vector-valued functions

When a function takes a vector in and gives a vector out, as every neural
network layer does, its derivative is a **matrix**:

```
       [ ∂y₁/∂x₁   ∂y₁/∂x₂   ...  ∂y₁/∂xₙ ]
J  =   [ ∂y₂/∂x₁   ∂y₂/∂x₂   ...  ∂y₂/∂xₙ ]        (m outputs × n inputs)
       [   ...                             ]
       [ ∂yₘ/∂x₁   ...            ∂yₘ/∂xₙ ]
```

Row `i` is the gradient of output `i`. The multivariable chain rule becomes
**matrix multiplication of Jacobians**:

```
if  y = f(x)  and  L = g(y),  then   ∂L/∂x = Jᵀ · ∂L/∂y
```

So backpropagation is a sequence of matrix-vector products, running the
forward pass's matrices in reverse, transposed. That's why a backward pass
costs about 2× a forward pass and runs on the same GPU matrix units.

### Backprop through a linear layer, by hand

For `y = W·x` (Ch. 10's layer, bias left out for clarity) with loss `L`, and
writing `δ = ∂L/∂y` (the "error signal" arriving from above), the two
gradients you need are:

```
∂L/∂W = δ · xᵀ          (an outer product: same shape as W)
∂L/∂x = Wᵀ · δ          (sent further back to the previous layer)
```

| Symbol | Meaning |
|---|---|
| `δ` | "delta", how much the loss would change per unit change of each output |
| `δ · xᵀ` | column × row = a matrix. Entry `[i][j] = δᵢ · xⱼ`: "error at output i × activity of input j" |
| `Wᵀ · δ` | the same weights, transposed, carry the blame backwards |

Same `W` and `x` as Ch. 10, target `t = (0, 2)`, loss `L = ½‖y - t‖²`:

```
x = (1, 2)       W = [0.5  -1  ]       y = W·x = (-1.5, 2.5)
                     [2     0.25]

L = ½·((-1.5-0)² + (2.5-2)²) = ½·(2.25 + 0.25) = 1.25
δ = y - t = (-1.5, 0.5)                                    (derivative of ½(y-t)²)

∂L/∂W = δ·xᵀ = [-1.5] · [1  2] = [-1.5  -3.0]
               [ 0.5]            [ 0.5   1.0]

∂L/∂x = Wᵀ·δ = [0.5   2  ] · [-1.5] = [ 0.25 ]
               [-1    0.25]   [ 0.5]   [ 1.625]
```

The weight connecting input 2 (the larger input, 2) to output 1 (the larger
error, -1.5) gets the biggest gradient, -3.0. Learning is assigning blame in
proportion to *activity × error*. That's the mathematical form of Hebb's
"neurons that fire together wire together". `loss.backward()` computes
exactly these two products for every layer.

### The most elegant gradient in deep learning: softmax + cross-entropy

Classifiers (and every LLM's next-token prediction) turn raw scores `z`
("logits") into probabilities with **softmax**, then score them with
**cross-entropy** against the true class `c`:

```
pᵢ = e^(zᵢ) / Σⱼ e^(zⱼ)                L = -ln(p_c)
```

| Symbol | Meaning |
|---|---|
| `zᵢ` | logit for class `i`, any real number |
| `e^(zᵢ)` | exponentiate to make it positive |
| `Σⱼ e^(zⱼ)` | normalizer, so the `pᵢ` sum to 1 |
| `p_c` | the probability the model gave the *correct* class |
| `-ln(p_c)` | 0 if `p_c = 1`, and grows without bound as `p_c → 0` |

Substitute and simplify: `L = -z_c + ln Σⱼ e^(zⱼ)`. Differentiate with respect
to `zᵢ`:

```
∂L/∂zᵢ = -[i = c] + e^(zᵢ)/Σⱼ e^(zⱼ)  =  pᵢ - yᵢ
```

where `yᵢ` is 1 for the true class and 0 otherwise (one-hot). After all
those exponentials and logs, **the gradient is just "predicted minus
actual"**.

**Example:** logits `z = (2.0, 1.0, 0.1)`, true class 0:

```
p = softmax(z) = (0.659, 0.242, 0.099)
L = -ln(0.659) = 0.417
∂L/∂z = p - y = (0.659-1, 0.242, 0.099) = (-0.341, 0.242, 0.099)
```

Push the correct logit up, push the wrong ones down, each in proportion to
how wrong it was. Every token of every LLM's pre-training uses this
gradient, trillions of times over.

### The Hessian: curvature in many dimensions

The matrix of all *second* partial derivatives:

```
H[i][j] = ∂²f / ∂xᵢ ∂xⱼ
```

It's symmetric (the order of differentiation doesn't matter for smooth `f`,
which is Schwarz's theorem), so by the spectral theorem (Ch. 11) its eigenvalues
are real. They give the **second-derivative test in `n` dimensions**. At a
point where `∇f = 0`:

| Eigenvalues of `H` | Point is | Shape |
|---|---|---|
| all `> 0` | local minimum | a bowl |
| all `< 0` | local maximum | a dome |
| mixed signs | **saddle point** | a horse saddle or mountain pass |

**Examples:**

```
f(x,y) = x² + xy + y²   ->   H = [2 1]   eigenvalues 1, 3  (Ch. 11!)  ->  minimum at (0,0)
                                 [1 2]
f(x,y) = x² - y²        ->   H = [2  0]  eigenvalues 2, -2            ->  saddle at (0,0)
                                 [0 -2]
```

**Why this matters for deep learning:** in a model with a billion
parameters, a critical point is a minimum only if *all* billion eigenvalues
are positive. If signs were coin flips, that would be absurdly unlikely, so
most points where the gradient vanishes are **saddles, not bad local
minima** (Dauphin et al., 2014). Plain gradient descent slows to a crawl near
saddles, and momentum (Ch. 21) helps carry it through. The Hessian's largest
eigenvalue also sets the maximum stable learning rate, `η < 2/λ_max`, which
Ch. 21 derives.

Putting it together, the **multivariable Taylor expansion** (Ch. 12, in `n`
dimensions) is the model behind every optimizer:

```
f(x + Δ) ≈ f(x)  +  ∇f(x)ᵀ·Δ  +  ½ · Δᵀ·H·Δ
           value     slope         curvature
```

Gradient descent uses the first two terms. Newton's method and its
approximations (L-BFGS, K-FAC, Shampoo) use all three.

### Networking in the wild: splitting traffic by marginal delay

Two backends: a fast one (`μ₁ = 100` req/s) and a slow one (`μ₂ = 50` req/s).
Total traffic `λ = 90` req/s. Send `x` to the fast one and `90 - x` to the
slow one. By Little's Law and the M/M/1 formula (Ch. 25), the average number
of requests in the system is:

```
N(x) = x/(μ₁ - x) + (λ - x)/(μ₂ - (λ - x))         average delay = N / λ
```

Set the derivative to zero. `d/dx [x/(μ - x)] = μ/(μ - x)²`, so the optimum
balances **marginal delays**, not loads:

```
μ₁/(μ₁ - x)² = μ₂/(μ₂ - λ + x)²    ->   x* ≈ 64.9 req/s to the fast server
```

| Policy | Split (fast / slow) | Avg delay |
|---|---|---|
| Round-robin (equal) | 45 / 45 | **109 ms**, since the slow server is at 90% load |
| Proportional to capacity | 60 / 30 | 33.3 ms |
| Gradient-optimal | 64.9 / 25.1 | **31.7 ms** |

Proportional-to-capacity looks "obviously fair", but the derivative says to
overload the fast server slightly, because its queue grows more slowly. This
is Gallager's 1977 **minimum-delay routing**: at the optimum, every path in use
has equal marginal delay `∂D/∂xᵢ`. It's the same "equalize the derivatives"
condition as the fleet-sizing example in Ch. 12, and the same idea behind
least-loaded and "power of two choices" load balancing in Envoy and NGINX.

```python
from scipy.optimize import minimize_scalar
mu1, mu2, lam = 100, 50, 90
N = lambda x: x/(mu1-x) + (lam-x)/(mu2-(lam-x))
best = minimize_scalar(N, bounds=(40.01, 89.99), method="bounded")
print(best.x, 1000 * N(best.x) / lam)    # ~64.85 req/s, ~31.7 ms
```

### Real-life engineering ties

- **AI:** this chapter *is* the prerequisite for understanding
  backpropagation and automatic differentiation (autograd/`.backward()` in
  PyTorch) rather than treating it as a magic black box.
- **SWE:** understanding sensitivity analysis — "how much does output
  change if this input changes slightly" — for tuning any system with
  multiple interacting parameters.
- **Network/Security:** less direct, but the same "how sensitive is the
  output to a small input change" question underlies stability analysis of
  control loops (Ch. 29) like TCP congestion control and rate limiters.

**More real-life examples:**
- Options pricing in finance defines "the Greeks" (Delta, Gamma, Vega) as
  partial derivatives of an option's price with respect to the underlying
  stock price, volatility, and time — the same partial-derivative
  machinery as this chapter, applied to markets instead of models.
- PyTorch's and TensorFlow's **autograd** engines build a computation
  graph of every tensor operation and apply the chain rule automatically,
  backward through the graph, to compute every gradient needed for
  training — this chapter is literally what `.backward()` does under the
  hood (Ch. 21 walks through it by hand).
- Distributed-systems performance engineers ask chain-rule-flavored
  sensitivity questions constantly: "if the auth service gets 10ms
  slower, how much slower is the end-to-end checkout request" — the
  answer depends on how many chained calls sit between the two, exactly
  like nested functions in the chain rule.
- Chemical and biological process control systems use partial derivatives
  of a target variable (e.g., reactor temperature) with respect to
  multiple inputs (flow rate, pressure) to decide which knob to turn — the
  same gradient idea used to train a neural net (Ch. 21), applied to a
  physical process.

### Common mistake

Forgetting to multiply by the *inner* derivative — computing only `f'(g(x))`
and stopping, instead of `f'(g(x)) · g'(x)`. This single dropped term is
the most common manual backpropagation error.

---

## 14. Probability Distributions

**Level:** Intermediate
**You'll use this for:** modeling request arrivals, error rates, time between failures, and measurement noise.

### The key distributions, and when each shows up

| Distribution | Models | Real example |
|---|---|---|
| **Bernoulli** | one yes/no trial | did this single request fail? |
| **Binomial** | count of successes in `n` independent trials | how many of 1000 requests failed? |
| **Poisson** | count of events in a fixed time window, at a known average rate | requests arriving per second at an API |
| **Normal (Gaussian)** | sum of many small independent effects | measurement noise, human heights, aggregated latencies |
| **Exponential** | time between events in a Poisson process | time until the next server failure |

### Binomial

```
P(k successes in n trials) = C(n,k) · p^k · (1-p)^(n-k)
```
```python
from scipy.stats import binom
# 1000 requests, each independently has 0.1% error rate — P(exactly 3 errors)?
binom.pmf(3, n=1000, p=0.001)
```

**A small binomial computed fully by hand** — 5 coin flips (`p=0.5`),
what's `P(exactly 2 heads)`?

```
C(5,2) = 5! / (2!·3!) = 120 / (2×6) = 10     <- how many ways to choose which 2 of 5 flips are heads (Ch. 15)
p^k = 0.5² = 0.25                             <- probability those 2 specific flips are heads
(1-p)^(n-k) = 0.5³ = 0.125                    <- probability the other 3 specific flips are tails

P(exactly 2 heads) = 10 × 0.25 × 0.125 = 0.3125
```

Every binomial probability is this same three-part product: "how many
arrangements achieve this outcome" × "probability of the success pattern"
× "probability of the failure pattern."

### Poisson — the request-arrival workhorse

Models "how many events happen in a fixed window" when events happen
independently at a constant average rate `λ`.

```
P(k events) = (λ^k · e^-λ) / k!
```
```python
from scipy.stats import poisson
# API averages 100 requests/sec. What's P(more than 130 arrive in one second)?
1 - poisson.cdf(130, mu=100)
```

**A small Poisson computed by hand** — a server averages `λ=2` crashes per
week, what's `P(exactly 3 crashes this week)`?

```
λ^k = 2³ = 8
e^-λ = e^-2 ≈ 0.1353         <- a fixed constant once λ is fixed, from a table or calculator
k! = 3! = 3×2×1 = 6

P(3 crashes) = (8 × 0.1353) / 6 = 1.0824 / 6 ≈ 0.1804
```

Same shape every time: raise `λ` to the `k`, multiply by the constant
`e^-λ`, divide by `k!` to account for the fact that the `k` events could
have landed in any order within the window.

**Why this matters directly for capacity planning:** if you provision for
*exactly* the average load, you will be under capacity roughly half the
time, because real traffic follows a Poisson-like distribution with
variance, not a flat line. This is the mathematical justification for
"provision for p99 load, not average load" (Ch. 25 builds directly on this
with queueing theory).

### Normal distribution and the Central Limit Theorem

The **Central Limit Theorem** says: the sum (or average) of many independent
random variables tends toward a normal distribution, *regardless of the
shape of the original distribution*. This is why so many real-world
measurements (and so many statistical tests, Ch. 22) assume normality —
it's not an arbitrary choice, it's a consequence of aggregating many small
independent effects.

### Exponential distribution — time between failures

If failures follow a Poisson process (constant hazard rate), the *time
between* failures follows an exponential distribution. This is the basis
of MTBF (Mean Time Between Failures) reasoning in reliability engineering.

```python
from scipy.stats import expon
# Mean time between failures = 720 hours. P(next failure within 100 hours)?
expon.cdf(100, scale=720)
```

### PMF, PDF, CDF: three ways to describe a distribution

```
PMF (discrete):    p(x) = P(X = x)                       Σₓ p(x) = 1
PDF (continuous):  f(x),   P(a ≤ X ≤ b) = ∫ₐᵇ f(x) dx     ∫ f(x) dx = 1
CDF (both):        F(x) = P(X ≤ x)                       F'(x) = f(x)
```

| Symbol | Say it as | Meaning |
|---|---|---|
| `p(x)` | "probability mass at x" | a real probability, used for countable outcomes (errors, packets, retries) |
| `f(x)` | "probability density at x" | **not** a probability. It's probability *per unit of x*, and can exceed 1 |
| `F(x)` | "cumulative distribution" | the fraction of outcomes at or below `x`, always rising from 0 to 1 |
| `F'(x) = f(x)` | | the Fundamental Theorem of Calculus (Ch. 12): the PDF is the CDF's derivative |

The CDF is the one engineers already use: **a percentile is the inverse
CDF**. "p99 latency = 250 ms" means `F(250 ms) = 0.99`.

### The normal density, every symbol explained

```
f(x) = 1/(σ·√(2π)) · e^( -(x - μ)² / (2σ²) )
```

| Piece | What it does |
|---|---|
| `μ` (mu) | the mean: where the peak sits |
| `σ` (sigma) | the standard deviation: how wide the bell is |
| `(x - μ)²` | squared distance from the centre. Symmetric, so the curve is too |
| `/(2σ²)` | measures that distance in units of `σ` |
| `e^(-…)` | turns distance into a weight: 1 at the centre, falling *very* fast in the tails |
| `1/(σ√(2π))` | the normalizer that makes the total area exactly 1, using the Gaussian integral `√π` (Ch. 12) |

**Example:** API latency `~ N(μ = 100 ms, σ = 20 ms)`. What fraction of
requests exceed 140 ms?

```
z = (140 - 100) / 20 = 2       ->   P(Z > 2) = 1 - Φ(2) ≈ 0.0228   (2.3%)
```

`Φ` (capital phi) is the standard normal CDF. You look it up in a table or
call `scipy.stats.norm.cdf(2)`. The `e^(-x²)` shape decides how quickly
tails die. At 4σ the probability is 0.003%, and at 6σ it's one in a billion.
Real latency is *not* normal (see heavy tails below), and that difference is
where outages come from.

### Every distribution's mean and variance at a glance

| Distribution | Parameters | `E[X]` | `Var(X)` | Engineering example |
|---|---|---|---|---|
| Bernoulli | `p` | `p` | `p(1-p)` | one request fails or not |
| Binomial | `n, p` | `np` | `np(1-p)` | errors in `n` requests |
| Geometric | `p` | `1/p` | `(1-p)/p²` | attempts until first success (Ch. 8) |
| Poisson | `λ` | `λ` | `λ` | arrivals per second |
| Uniform | `a, b` | `(a+b)/2` | `(b-a)²/12` | random backoff jitter in `[a, b]` |
| Exponential | `λ` | `1/λ` | `1/λ²` | time between arrivals or failures |
| Normal | `μ, σ²` | `μ` | `σ²` | aggregated measurement noise |

Notice the **Poisson's mean equals its variance**. That gives you a free
diagnostic. If your "requests per second" metric has a variance much larger
than its mean (*overdispersion*), arrivals aren't independent. That points to
bursty clients, retry storms, or cron jobs firing on the minute.

### The Poisson limit theorem: where Poisson comes from

> **Theorem.** If `n → ∞` and `p → 0` with `n·p = λ` held fixed, then
> `Binomial(n, p) → Poisson(λ)`.

Many independent chances, each tiny, give Poisson. That's why it models
requests (millions of users, each rarely clicking), disk failures (many disks,
each rarely failing) and radioactive decay alike.

**Example:** 1000 requests at a 0.1% error rate, so `λ = 1000·0.001 = 1`.

```
Binomial: P(exactly 3) = C(1000,3)·0.001³·0.999⁹⁹⁷ = 0.06128
Poisson:  P(exactly 3) = 1³·e⁻¹/3!                 = 0.06131
```

The two agree to three decimal places, and the Poisson version can be worked
out by hand.

### Memorylessness: the exponential distribution's strange property

```
P(T > s + t | T > s) = P(T > t)
```

**Proof in one line**, using `P(T > t) = e^(-λt)` from the CDF:

```
P(T > s+t | T > s) = P(T > s+t) / P(T > s) = e^(-λ(s+t)) / e^(-λs) = e^(-λt)
```

If failures are exponential, a server that has run for 1,000 hours is exactly
as likely to fail in the next hour as a brand-new one. That's why
time-based replacement of exponentially-failing parts is pointless. It's
also why "it's been quiet for an hour, so a request must be due" is
the gambler's fallacy again. Real hardware isn't fully memoryless. Disks
follow a "bathtub curve" with high infant mortality, a flat middle and a
wear-out phase, and Backblaze's data shows exactly that shape.

### The Law of Large Numbers and the Central Limit Theorem, stated properly

Let `X₁, …, Xₙ` be independent draws from *any* distribution with mean `μ` and
standard deviation `σ`, and let `X̄ₙ = (X₁ + … + Xₙ)/n` be their average.

> **Law of Large Numbers.** `X̄ₙ → μ` as `n → ∞`.
> Averages settle on the true mean.

> **Central Limit Theorem.** For large `n`,
> ```
> (X̄ₙ - μ) / (σ / √n)   ≈   N(0, 1)
> ```
> Averages are approximately normal, centred on `μ`, with spread `σ/√n`,
> whatever the shape of the original distribution.

| Symbol | Meaning |
|---|---|
| `X̄ₙ` | the sample mean of `n` observations |
| `σ/√n` | the **standard error**: how much the sample mean itself jitters |
| `√n` | the price of precision. **4× the data halves the error** |
| `N(0, 1)` | the standard normal: mean 0, standard deviation 1 |

**Example:** individual request latencies are exponential with mean 50 ms,
so `σ = 50 ms` as well. That's strongly skewed and nothing like a bell. Now
average them in batches of 100:

```
standard error = 50 / √100 = 5 ms
95% of batch averages fall in  50 ± 1.96·5  =  [40.2, 59.8] ms
```

A simulation of 100,000 such batches gives a mean of 49.99, a standard
deviation of 5.01, and 95.0% inside that interval. The skew has disappeared,
as the theorem promises.

```python
import numpy as np
rng = np.random.default_rng(0)
batch_means = rng.exponential(scale=50, size=(100_000, 100)).mean(axis=1)
print(batch_means.mean(), batch_means.std())          # ~50.0, ~5.0
print(np.mean(np.abs(batch_means - 50) < 1.96 * 5))   # ~0.95
```

The CLT is why A/B tests (Ch. 22) can use normal-based z-tests on
conversion rates, and why mini-batch gradients in ML are noisy estimates of
the true gradient whose noise shrinks like `1/√batch_size`. Doubling the batch
cuts gradient noise by only about 30%, which is why "just use bigger batches"
has diminishing returns.

### Concentration inequalities: guarantees without assuming normality

The CLT says "approximately, for large `n`". Sometimes you need a
**guarantee** that holds for any distribution. These three inequalities give
one:

```
Markov:      P(X ≥ a)          ≤ E[X] / a                 (X ≥ 0)
Chebyshev:   P(|X - μ| ≥ kσ)   ≤ 1 / k²
Hoeffding:   P(|X̄ₙ - μ| ≥ ε)   ≤ 2·e^(-2nε²)               (each Xᵢ in [0, 1])
```

| Symbol | Meaning |
|---|---|
| `a` | a threshold you're worried about |
| `k` | how many standard deviations out |
| `ε` | "epsilon", the error tolerance you'll accept on an average |
| `n` | number of independent samples |

**Markov example:** average queue length is 4. Without knowing anything else,
`P(queue ≥ 40) ≤ 4/40 = 10%`. It's a crude bound, but it holds for any
non-negative quantity.

**Chebyshev example:** for *any* distribution, at most `1/9 ≈ 11%` of
values lie beyond 3σ. The normal distribution puts only 0.3% out there, so a
"3-sigma alert" on non-normal data can fire about 40× more often than you
designed for.

**Hoeffding, in AI: how big should an eval set be?** You want the model's
measured accuracy within `ε = 2` points of the truth, with 95% confidence
(failure probability `δ = 0.05`). Set `2e^(-2nε²) = δ` and solve:

```
n ≥ ln(2/δ) / (2ε²) = ln(40) / (2 · 0.02²) = 3.689 / 0.0008 ≈ 4,612 examples
```

Run it the other way and a 500-question benchmark only pins accuracy to
`ε = √(ln 40 / 1000) ≈ ±6` points. A leaderboard gap of 3 points on such a
benchmark is noise. Hoeffding is also the backbone of PAC learning theory,
the formal answer to "how much data do I need to learn?"

### Heavy tails and the "tail at scale"

Many systems quantities follow a **Pareto (power-law)** distribution, where
the tail decays polynomially instead of exponentially:

```
P(X > x) = (x_m / x)^α        for x ≥ x_m
```

| Symbol | Meaning |
|---|---|
| `x_m` | the minimum value (the scale) |
| `α` | tail index. Smaller is heavier. `α ≤ 2` means infinite variance, `α ≤ 1` means infinite mean |

With `α ≈ 1.16` you get the "80/20 rule": 20% of the flows carry 80% of the
bytes. Internet flow sizes, file sizes and request costs all look like this.
In networking these are the "elephant" flows (few, huge) and "mice" flows (many,
tiny), and datacenter switches schedule them differently.

**Tail latency compounds under fan-out.** A search query fans out to 100
leaf servers and must wait for the slowest. If each leaf is slow (above its
p99) independently 1% of the time:

```
P(query is slow) = 1 - P(all 100 fast) = 1 - 0.99¹⁰⁰ ≈ 63%
```

**One leaf's p99 becomes the whole service's median.** This is the central
result of Google's "The Tail at Scale" (Dean & Barroso, 2013). It's why
hedged requests (send a backup after the p95 has passed) and tied requests
exist, and why p99.9 matters more than p50 in microservice architectures.

### AI in the wild: sampling text is sampling a distribution

At every step an LLM outputs logits `z`, turns them into a categorical
distribution with softmax (Ch. 13), and **draws a random sample** from it.
**Temperature** `T` reshapes that distribution before sampling:

```
pᵢ = e^(zᵢ / T) / Σⱼ e^(zⱼ / T)
```

Logits `z = (2.0, 1.0, 0.1)` for three candidate tokens:

| Temperature | p(token 1) | p(token 2) | p(token 3) | Behaviour |
|---|---|---|---|---|
| `T = 0.5` | 0.864 | 0.117 | 0.019 | sharp, near-deterministic |
| `T = 1.0` | 0.659 | 0.242 | 0.099 | the model's raw beliefs |
| `T = 2.0` | 0.502 | 0.304 | 0.194 | flat, more "creative" and more errors |

As `T → 0` this becomes `argmax` (greedy decoding). As `T → ∞` it becomes
uniform. **Top-p (nucleus) sampling** truncates to the smallest set of
tokens whose probabilities sum to `p` (say 0.9), which is just the CDF from
earlier in this chapter applied to the sorted probabilities.

### Real-life engineering ties

- **SWE/SRE:** error-budget math, capacity planning against percentile —
  not average — load (Ch. 9, 25).
- **AI:** assumed noise models in regression (Gaussian noise), Poisson
  regression for count data (e.g., predicting clicks/purchases).
- **Network:** packet arrival modeling (Poisson), queueing theory inputs
  (Ch. 25 uses this directly — M/M/1 queues assume Poisson arrivals and
  exponential service times).
- **Security:** modeling time-between-attacks/scans for anomaly baselining;
  Poisson-based alerting ("we normally see ~5 failed logins/hour — 200 in
  the last hour is statistically extreme," tying back to Ch. 9's 3-sigma idea).

**More real-life examples:**
- Telephone exchanges — and modern call-center staffing tools — use the
  Erlang C formula, built on the same Poisson-arrival assumption as this
  chapter, to compute how many agents are needed to hit a target
  wait-time SLA; the unit "Erlang" is named for the mathematician who
  solved this exact problem in 1917.
- Backblaze's publicly published hard-drive failure statistics report
  **Annualized Failure Rate**, used with an exponential (constant hazard
  rate) model to estimate the probability a specific drive fails within
  the next year — the same math as the MTBF example above, with real,
  public data.
- Insurance actuaries model claim frequency with Poisson distributions
  and claim severity with heavy-tailed distributions (log-normal or
  Pareto) — nearly identical in structure to the risk model in Ch. 26.
- A/B testing platforms model each user's conversion as a Bernoulli trial
  and the aggregate conversion count as Binomial — the statistical
  foundation underneath every "statistically significant" banner those
  tools show.

### Common mistake

Applying Normal-distribution intuition (symmetric, thin tails) to data
that's actually **heavy-tailed** (network latency, financial returns,
security incident severity) — heavy-tailed data has "black swan" events far
more often than a Normal model predicts, badly underestimating tail risk.

---

## 15. Combinatorics & Counting

**Level:** Intermediate
**You'll use this for:** password/keyspace sizing, hash collision estimation, and any "how many possible X" question.

### Permutations vs. combinations

- **Permutations** (order matters): arranging `n` items in `n!` ways;
  choosing `k` ordered items from `n`: `P(n,k) = n!/(n-k)!`
- **Combinations** (order doesn't matter): choosing `k` unordered items
  from `n`: `C(n,k) = n!/(k!(n-k)!)`

```python
from math import factorial, comb, perm
perm(10, 3)   # 720  -- ordered arrangements of 3 from 10
comb(10, 3)   # 120  -- unordered groups of 3 from 10
```

### Where `P(n,k)` and `C(n,k)` actually come from

**Permutations, built from the fundamental counting principle.** To
arrange 3 items chosen (in order) from 10: the 1st pick has 10 options,
the 2nd pick has 9 remaining, the 3rd has 8 remaining — multiply:

```
P(10,3) = 10 × 9 × 8 = 720

Written as factorials, this is "10! but stop after 3 terms," which is
exactly what dividing by the leftover tail (7!) achieves:
  10! / 7!  =  (10×9×8×7×6×5×4×3×2×1) / (7×6×5×4×3×2×1)  =  10×9×8  =  720
  and 7 = 10 - 3, so this is P(n,k) = n! / (n-k)!
```

**Combinations = permutations, with duplicate orderings divided back out.**
`C(10,3)` counts the same picks as `P(10,3)`, but treats `{A,B,C}` and
`{B,A,C}` etc. as the *same* group. Each group of 3 can be internally
reordered `3! = 6` ways, so every distinct group got counted 6 times over
in `P(10,3)`:

```
C(10,3) = P(10,3) / 3! = 720 / 6 = 120
        = n! / (n-k)! / k! = n! / (k!(n-k)!)     <- matches the formula above
```

That's the whole derivation: combinations are permutations with the
"which order" duplicates divided away.

### Worked example: password keyspace

A password of length `L` using an alphabet of size `A` has `A^L` possible
values (this is a counting problem, not a probability problem — it's the
*size of the sample space*).

```python
def keyspace(alphabet_size, length):
    return alphabet_size ** length

lower_only = keyspace(26, 8)          # 26^8 ≈ 2.09 * 10^11
mixed_case_digits_symbols = keyspace(95, 12)  # 95^12 ≈ 5.4 * 10^23
```
This single formula is the mathematical basis of every "password strength"
estimate (fully connected to entropy and crack-time in Ch. 26).

### The pigeonhole principle

If you put `n+1` items into `n` boxes, at least one box holds 2+ items.
Trivial-sounding, but it's the *proof* that hash collisions are
mathematically guaranteed once you hash more items than the output space
size — a hash function mapping an unbounded input space to a fixed-size
(say, 256-bit) output space **must** have collisions somewhere; the only
question is how hard they are to find (Ch. 26's birthday paradox
sharpens this into a concrete number).

### Real-life engineering ties

- **SWE:** estimating combinatorial test-case explosion, counting valid
  states in a state machine.
- **AI:** hyperparameter search space sizing (grid search over `k`
  parameters, each with `m` values, is `m^k` combinations — the
  combinatorial explosion that justifies random/Bayesian search instead).
- **Network:** counting possible subnet allocations, VLAN ID space
  (12-bit VLAN tags = `2^12 = 4096` possible VLANs, a real capacity limit).
- **Security:** keyspace-size calculations underlie every brute-force
  time estimate (Ch. 26); counting attack paths through a system
  (attack-graph enumeration) is a permutations/combinations problem.

**More real-life examples:**
- Lottery odds (e.g., choosing 6 numbers from 49) are a pure `C(n,k)`
  combinations calculation — a concrete, intuitive way to build number
  sense for the same formula used in password-keyspace reasoning.
- QA teams facing combinatorial explosion in configuration testing (e.g.,
  5 browsers x 4 OSes x 3 screen sizes = thousands of combinations) use
  **pairwise (combinatorial) testing**, which mathematically guarantees
  covering all *pairs* of parameter values with far fewer test cases than
  the full combinatorial product.
- API key and session token generators size their random identifier
  space (e.g., 128-bit UUIDs) using the same pigeonhole/collision
  reasoning as password entropy, making accidental collisions between two
  independently generated tokens astronomically unlikely.
- NIST's 2017 password guidance (SP 800-63B) explicitly moved away from
  complexity rules (mandatory symbols/digits) toward length-based entropy
  reasoning, citing exactly the "complexity rules reduce real entropy"
  effect covered in Ch. 23's common mistake.

### Common mistake

Confusing "order matters" (permutation) with "order doesn't matter"
(combination) — e.g., a 4-digit PIN is a permutation problem (order
matters: 1234 ≠ 4321, and repeats are allowed: `10^4`), while "choose 4
servers out of 10 to take offline for maintenance" is a combination
problem.

---

## 16. Bayes' Theorem

**Level:** Intermediate
**You'll use this for:** spam filters, medical/security test interpretation, and understanding why "99% accurate" detectors can still be mostly wrong.

### The formula

```
P(A|B) = [ P(B|A) · P(A) ] / P(B)
```
In words: "the probability of A given that B happened" is related to, but
**not equal to**, "the probability of B given that A happened." Confusing
these two is called the **prosecutor's fallacy**, and it's one of the most
consequential, most common reasoning errors in engineering, medicine, and law.

### Applying the formula by hand, with natural frequencies (no algebra needed)

The cleanest way to compute Bayes' theorem by hand is to skip the fraction
and just count people. Say 1% of a population of 1,000 has a disease, and
a test is 90% accurate both ways:

```
1,000 people total
  -> 10 actually sick   (1% of 1,000)
  -> 990 actually healthy

Of the 10 sick people, the test catches 90%:
  -> 9 correctly test positive (true positives)
  -> 1 sick person tests negative (missed)

Of the 990 healthy people, the test correctly clears 90%:
  -> 891 correctly test negative
  -> 99 healthy people wrongly test positive (false positives)

Total positive tests = 9 (true) + 99 (false) = 108

P(actually sick | tested positive) = 9 / 108 ≈ 0.083   -> only ~8.3%!
```

That fraction — true positives divided by *all* positives — is Bayes'
theorem with the algebra hidden: you never had to touch `P(B|A)·P(A)/P(B)`
directly, just count outcomes in a hypothetical group of real people. This
"natural frequency" method is the fastest way to sanity-check any
Bayes-flavored claim without a calculator.

### The base rate fallacy — the single most important lesson in this chapter

Suppose an intrusion detection system is 99% accurate (both true-positive
and true-negative rate), and genuinely malicious events are rare: 1 in
10,000 events. Out of 1,000,000 events/day, how many alerts are *actually*
malicious?

```python
p_malicious = 1 / 10_000
p_benign = 1 - p_malicious
accuracy = 0.99

total_events = 1_000_000
actual_malicious = total_events * p_malicious        # 100
actual_benign    = total_events * p_benign            # 999,900

true_positives  = actual_malicious * accuracy         # 99
false_positives = actual_benign * (1 - accuracy)      # 9,999

precision = true_positives / (true_positives + false_positives)
print(precision)   # ≈ 0.0098 -> under 1% of alerts are real!
```

**This is not a contrived example — it is the mathematical reason "alert
fatigue" is a structural, unavoidable problem for any detector applied to a
rare-event population**, no matter how "accurate" the detector sounds.
When the base rate of the thing you're looking for is very low, even a
low false-positive *rate* produces a false-positive *majority* in absolute
terms. This is a direct, load-bearing consequence of Bayes' theorem and
must inform how any Sr. security/ML engineer reads a confusion matrix
(Ch. 22 and 26 build directly on this with precision/recall/ROC curves).

### Naive Bayes classifier (spam filtering)

Bayes' theorem, applied to classification: given the words in an email,
what's `P(spam | words)`? "Naive" because it assumes words are
conditionally independent given the class (usually false, but works well
in practice).

```
P(spam | word1, word2, ...) ∝ P(spam) · P(word1|spam) · P(word2|spam) · ...
```

Multiplying hundreds of small probabilities underflows to 0.0 in floating
point (Ch. 1), so real implementations add **log**-probabilities instead
(Ch. 5): `log P(spam) + Σ log P(wordᵢ | spam)`. The `∝` ("proportional to")
means we dropped the denominator `P(words)`. It's the same for both classes,
so it can't change which one wins.

### Every symbol in Bayes' theorem has a name

```
              likelihood   prior
              ┌───┴───┐   ┌─┴─┐
P(H | E)  =   P(E | H)  · P(H)
└──┬───┘      ─────────────────
posterior          P(E)
                   └─┬─┘
                  evidence
```

| Symbol | Name | Meaning in an intrusion-detection example |
|---|---|---|
| `H` | hypothesis | "this login is an attacker" |
| `E` | evidence | "the login came from a new country at 3 a.m." |
| `P(H)` | **prior** | how common attackers are *before* looking at this login (the base rate) |
| `P(E \| H)` | **likelihood** | how often attackers' logins look like this |
| `P(E)` | **evidence** (normalizer) | how often *any* login looks like this, attacker or not |
| `P(H \| E)` | **posterior** | the updated belief, which is what you actually want |

The formula is "**posterior ∝ likelihood × prior**". Every learning system
that updates beliefs from data, from spam filters to Kalman filters in GPS to
Bayesian optimization of hyperparameters, runs this multiplication.

### Bayesian updating in practice: the Beta distribution

How do you estimate a rate (a conversion rate, a link's loss rate, a model's
accuracy) *and* how uncertain you are about it? Use a **Beta**
distribution as the belief about the unknown rate `θ`:

```
prior:      θ ~ Beta(α, β)
observe:    k successes in n trials
posterior:  θ ~ Beta(α + k, β + n - k)          mean = (α + k) / (α + β + n)
```

| Symbol | Meaning |
|---|---|
| `θ` | the unknown true rate, between 0 and 1 |
| `α`, `β` | "pseudo-counts" of prior successes and failures. `Beta(1,1)` is uniform: "no idea" |
| `k`, `n - k` | the observed successes and failures. You literally just add them on |

This is called a **conjugate prior**: Beta prior × binomial likelihood =
Beta posterior. A Bayesian update is two additions.

**Example:** a new feature variant, flat prior `Beta(1, 1)`. On day 1 you see
13 conversions in 100 visitors, so the posterior is `Beta(14, 88)`, with mean
`14/102 ≈ 13.7%` and a 95% credible interval of about **[7.8%, 21.0%]**. The
interval is wide because 100 visitors isn't much. After 10,000 visitors it
would narrow to about ±0.7 points.

**AI in the wild: Thompson sampling.** To choose among several variants (ads,
recommendation slots, LLM prompts), draw one random sample from each
variant's Beta posterior and show the one with the highest draw. Uncertain
variants sometimes produce high draws, so they get explored. Proven variants
usually win, so they get exploited. This "multi-armed bandit" strategy wastes
far less traffic than a fixed 50/50 A/B test, and it runs in production
recommendation and ad systems.

**Networking in the wild:** the same update tracks a link's delivery
probability from probe results (each probe is a Bernoulli trial), with a
built-in uncertainty that shrinks as probes accumulate. That's more honest
than an EWMA (Ch. 9) when probes are sparse.

```python
from scipy.stats import beta
post = beta(1 + 13, 1 + 87)            # flat prior + 13 of 100
print(post.mean(), post.interval(0.95))  # 0.137, (0.078, 0.210)

import numpy as np                       # Thompson sampling: pick the arm with the highest draw
rng = np.random.default_rng(0)
arms = {"A": (1 + 120, 1 + 880), "B": (1 + 13, 1 + 87)}   # (α, β) per variant
draws = {k: rng.beta(a, b) for k, (a, b) in arms.items()}
print(max(draws, key=draws.get), draws)
```

### Real-life engineering ties

- **SWE:** interpreting any classifier's output correctly instead of
  treating "confidence score" as "probability of being right."
- **AI:** Naive Bayes classifiers, Bayesian A/B testing, Bayesian
  hyperparameter optimization; Bayes is also the formal foundation of
  updating a model's belief as new data arrives (Bayesian inference,
  Ch. 22).
- **Network:** anomaly detection systems that flag "unusual" traffic must
  be evaluated with base-rate awareness — most traffic is normal, so even
  a "good" anomaly score threshold floods you with false positives.
- **Security:** **this is arguably the single highest-leverage math
  concept for a security engineer** — every SIEM tuning decision, every
  "should we lower this detection threshold" conversation, is implicitly a
  Bayes' theorem question about precision vs. recall tradeoffs.

**More real-life examples:**
- The classic medical teaching example — a mammogram with 90% sensitivity
  screening for a cancer with 1% prevalence — produces a *majority* of
  positive results being false positives, for exactly the same base-rate
  reason as the IDS example above; doctors are famously shown to get this
  calculation wrong as often as patients.
- Early spam filters (SpamAssassin, and Paul Graham's influential 2002
  "A Plan for Spam" essay) popularized Naive Bayes spam classification,
  directly kicking off widespread industry adoption of Bayesian filtering.
- Payment fraud systems must explicitly balance **false declines**
  (blocking a legitimate purchase) against **missed fraud**, and because
  legitimate transactions vastly outnumber fraudulent ones, a fraud
  model's naive "accuracy" is nearly meaningless — precision/recall framed
  around the true (low) fraud base rate is what actually matters.
- Recommendation systems handle new users with no history (the "cold
  start" problem) by falling back to a **prior** — a population-average
  preference distribution — and updating toward personalized
  recommendations as interaction data arrives, a live, continuous
  application of Bayesian updating.

### Common mistake

Reading "99% accurate" and concluding "99% of alerts are real." Accuracy
(or even true-positive rate) says nothing about **precision** without
knowing the base rate — always ask "how rare is the thing being detected"
before trusting a detector's headline number.

---

## 17. Graph Theory

**Level:** Intermediate
**You'll use this for:** network topology, routing, dependency graphs, social graphs, and attack-path modeling.

### The idea

A graph is a set of **nodes** (vertices) connected by **edges**. Graphs can
be **directed** (edges have a direction, like "A follows B") or
**undirected** (edges are symmetric, like "A is physically cabled to B"),
and **weighted** (edges have a cost, like latency or bandwidth) or
unweighted.

### Representations

```python
# Adjacency list — efficient for sparse graphs (most real networks)
graph = {
    "A": ["B", "C"],
    "B": ["A", "D"],
    "C": ["A", "D"],
    "D": ["B", "C"],
}

# Adjacency matrix — efficient for dense graphs, enables linear-algebra tricks (Ch. 11)
import numpy as np
nodes = ["A", "B", "C", "D"]
M = np.array([
    [0,1,1,0],
    [1,0,0,1],
    [1,0,0,1],
    [0,1,1,0],
])
```

### Traversal: BFS vs. DFS

- **BFS (breadth-first search):** explores level by level — finds the
  *shortest path* (in unweighted graphs) and is used for "closest nodes
  first" problems.
- **DFS (depth-first search):** explores as deep as possible before
  backtracking — used for cycle detection, topological sort, and
  exhaustive path enumeration.

```python
from collections import deque

def bfs_shortest_path(graph, start, goal):
    queue = deque([[start]])
    visited = {start}
    while queue:
        path = queue.popleft()
        node = path[-1]
        if node == goal:
            return path
        for neighbor in graph.get(node, []):
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(path + [neighbor])
    return None

print(bfs_shortest_path(graph, "A", "D"))   # ['A', 'B', 'D'] or ['A', 'C', 'D']
```

### Dijkstra's algorithm — shortest path with weighted edges

Used by **every link-state routing protocol** (OSPF, IS-IS) to compute the
shortest path tree from a router to every other router in the network.

```python
import heapq

def dijkstra(graph, start):
    # graph: {node: [(neighbor, weight), ...]}
    distances = {node: float('inf') for node in graph}
    distances[start] = 0
    pq = [(0, start)]
    while pq:
        dist, node = heapq.heappop(pq)
        if dist > distances[node]:
            continue
        for neighbor, weight in graph[node]:
            new_dist = dist + weight
            if new_dist < distances[neighbor]:
                distances[neighbor] = new_dist
                heapq.heappush(pq, (new_dist, neighbor))
    return distances

weighted_graph = {
    "A": [("B", 1), ("C", 4)],
    "B": [("A", 1), ("C", 2), ("D", 5)],
    "C": [("A", 4), ("B", 2), ("D", 1)],
    "D": [("B", 5), ("C", 1)],
}
print(dijkstra(weighted_graph, "A"))   # {'A': 0, 'B': 1, 'C': 3, 'D': 4}
```

**Tracing Dijkstra by hand** on that same graph, starting from `A`. Keep a
running "best known distance" per node, always expanding the closest
unvisited node next:

```
Start: distances = {A:0, B:∞, C:∞, D:∞}

Visit A (dist 0), its edges: B (1), C (4)
  -> B: 0+1=1 < ∞, update B to 1
  -> C: 0+4=4 < ∞, update C to 4
  distances = {A:0, B:1, C:4, D:∞}

Visit B (dist 1, the smallest unvisited), its edges: A (1), C (2), D (5)
  -> A: already visited, skip
  -> C: 1+2=3 < 4, update C to 3   <- found a shorter path to C via B!
  -> D: 1+5=6 < ∞, update D to 6
  distances = {A:0, B:1, C:3, D:6}

Visit C (dist 3, now the smallest unvisited), its edges: A (4), B (2), D (1)
  -> D: 3+1=4 < 6, update D to 4   <- found a shorter path to D via C!
  distances = {A:0, B:1, C:3, D:4}

Visit D (dist 4) — no shorter paths found through it. Done.

Final: A=0, B=1, C=3, D=4   ✓ matches the code output
```

The key move is that C's distance got revised *twice* (4, then 3) —
Dijkstra never commits to a distance as final until that node is actually
visited, which is exactly why it always finds the true shortest path
instead of just the first path found.

**Bellman-Ford** is the alternative used when edge weights can be
*negative* (Dijkstra breaks with negative weights) — this is why
distance-vector protocols like the classic RIP use a Bellman-Ford variant,
while OSPF (Dijkstra-based) requires non-negative link costs.

### Real-life engineering ties

- **SWE:** dependency graphs (build systems, package managers — a
  dependency cycle is a graph cycle-detection problem), org charts, state
  machines.
- **AI:** knowledge graphs, graph neural networks (GNNs), attention as a
  fully-connected weighted graph between tokens.
- **Network:** literally how routing works — network topology *is* a
  graph, OSPF/IS-IS *is* Dijkstra, BGP path selection is graph traversal
  with policy-weighted edges (Ch. 25).
- **Security:** **attack graphs** (modeling every path an attacker could
  take from initial foothold to target) are directed graphs, and finding
  the shortest/cheapest attack path is literally a Dijkstra/BFS problem;
  lateral-movement analysis in incident response is graph traversal over
  an access/trust graph.

**More real-life examples:**
- BGP, the routing protocol that holds the internet together, is a
  graph-based path-vector protocol — and BGP hijacking incidents (like
  the widely documented 2008 Pakistan Telecom incident that briefly took
  YouTube offline globally) happen because the graph has no built-in
  trust verification for advertised paths.
- Package managers (npm, pip, Cargo) and build systems (Bazel, Make)
  represent dependencies as a **Directed Acyclic Graph (DAG)** — a
  dependency cycle is a graph-cycle-detection bug, and build order is a
  topological sort of that DAG.
- Git's commit history is itself a DAG; commands like `git log --graph`,
  merge-base computation, and `git bisect` are all graph algorithms
  running directly on your repository.
- Service mesh tools (Istio, Linkerd) and distributed tracing platforms
  model microservice call patterns as a graph specifically to compute
  blast radius, critical-path latency, and single-point-of-failure
  analysis — all direct applications of this chapter.

### Common mistake

Using Dijkstra on a graph with negative edge weights (it will produce
silently wrong answers, not an error) — check for negative weights before
choosing an algorithm.

---

## 18. Complexity & Asymptotic Analysis

**Level:** Intermediate
**You'll use this for:** choosing the right algorithm/data structure at scale, and estimating brute-force attack feasibility.

### Big-O notation

Big-O describes how an algorithm's resource usage (time or memory) grows as
input size `n` grows — the *worst-case upper bound*, ignoring constant
factors, because constants stop mattering as `n` gets large.

| Complexity | Name | Example |
|---|---|---|
| `O(1)` | constant | hash table lookup |
| `O(log n)` | logarithmic | binary search (Ch. 5) |
| `O(n)` | linear | scanning a list once |
| `O(n log n)` | linearithmic | efficient sorting (merge sort, quicksort avg) |
| `O(n²)` | quadratic | nested loops, naive sorting |
| `O(2^n)` | exponential | brute-force subset enumeration |
| `O(n!)` | factorial | brute-force traveling salesman |

```python
import time

def linear_search(lst, target):        # O(n)
    return target in lst

def binary_search(sorted_lst, target): # O(log n) — requires sorted input
    lo, hi = 0, len(sorted_lst) - 1
    while lo <= hi:
        mid = (lo + hi) // 2
        if sorted_lst[mid] == target: return mid
        elif sorted_lst[mid] < target: lo = mid + 1
        else: hi = mid - 1
    return -1
```

### How to derive Big-O by hand: count the operations

**Linear search** on a list of `n` items: in the worst case (target is
last, or missing entirely), the loop runs once per element — `n`
comparisons. Drop the word "comparisons," keep the shape: `O(n)`.

**Binary search** on a sorted list of `n` items: each pass through the
`while` loop discards *half* the remaining candidates. Starting from `n`
items, after 1 comparison you have `n/2` left, after 2 you have `n/4`,
after `k` comparisons you have `n/2^k` left. The loop stops once 1 item
remains:

```
n / 2^k = 1   =>   2^k = n   =>   k = log2(n)     (Ch. 5's log rule, applied)
```

So for `n = 1,000,000`, binary search needs at most `log2(1,000,000) ≈ 20`
comparisons — count them by hand for a small case to see the pattern:
`n=16 -> 8 -> 4 -> 2 -> 1` is 4 halvings, and `log2(16) = 4`, confirming
the formula.

**The general method** for any piece of code: count how many times the
innermost operation runs as a function of `n`. A single loop over `n`
items is `O(n)`. A loop *inside* another loop, both over `n`, multiplies:
`O(n) × O(n) = O(n²)`. A loop that *halves* its remaining work each pass
(like binary search) is `O(log n)`. Reading nested loops and halving
patterns this way is literally how you derive a Big-O class without
memorizing a table of answers.

### Why the growth rate matters more than the constant

At `n = 10`, an `O(n²)` algorithm might beat an `O(n log n)` one due to
constant factors/overhead. At `n = 10,000,000`, `O(n²)` is *catastrophically*
worse — `n²` is 10 trillion operations vs. `n log n`'s ~230 million. This is
why "premature optimization" advice applies to constants, not to
asymptotic class — **algorithm choice matters far more at scale than
micro-optimization**.

### Worked example: brute-force time estimate

If a keyspace has `2^128` possible keys (Ch. 15) and an attacker can try
`10^12` keys/second (a fast, dedicated cracking rig), how long to brute
force?

```python
keyspace = 2**128
attempts_per_sec = 10**12
seconds = keyspace / attempts_per_sec
years = seconds / (60*60*24*365)
print(f"{years:.2e} years")   # ~1.08e19 years — vastly longer than the age of the universe
```
This single Big-O-adjacent calculation (exponential keyspace growth vs.
linear attacker speed growth) is *why* doubling a key length doesn't
double the difficulty — it squares it, which is the entire mathematical
argument for why AES-256 and RSA-2048+ are considered safe against
brute force for the foreseeable future (Ch. 24, 26).

### Real-life engineering ties

- **SWE:** choosing `O(1)` hash-map lookups over `O(n)` list scans at
  scale; recognizing an accidental `O(n²)` (e.g., a `list.remove()` or
  string concatenation inside a loop) before it becomes a production
  incident.
- **AI:** attention mechanism complexity is `O(n²)` in sequence length —
  the entire reason "long context" is expensive and the motivation behind
  efficient-attention research (sparse/linear attention).
- **Network:** routing table lookup complexity, packet processing
  pipelines designed for `O(1)`/`O(log n)` per-packet cost at line rate.
- **Security:** brute-force feasibility calculations (above), and
  recognizing algorithmic-complexity attacks (an attacker crafting input
  that forces your `O(n²)` worst-case path — e.g., hash-flooding attacks
  that force hash-table collisions deliberately, CWE-407).

**More real-life examples:**
- Cloudflare's July 2019 global outage was caused by a single regular
  expression with catastrophic backtracking — effectively an accidental
  `O(2^n)` worst case — that pegged CPU to 100% across their edge
  network; this failure mode has a name, **ReDoS (Regular Expression
  Denial of Service)**, and is a listed CWE (CWE-1333).
- Python's and Java's default sort (Timsort) was engineered to be
  `O(n log n)` worst-case while detecting and exploiting already-sorted
  runs in real data to approach `O(n)` best-case — a deliberate
  complexity-class engineering decision, not an accident.
- Database query planners (PostgreSQL's `EXPLAIN ANALYZE`, MySQL's
  optimizer) choose between an index scan (`O(log n)`) and a full table
  scan (`O(n)`) based on estimated row counts — reading a query plan is a
  direct, daily application of Big-O reasoning.
- Hash-flooding DoS attacks (CWE-407) exploit predictable hash functions
  to force many keys into the same bucket, degrading a hash table's
  expected `O(1)` lookup into `O(n)` — the reason modern language runtimes
  (Python, Ruby, PHP) randomize their hash seed at process start.

### Common mistake

Reporting Big-O without stating *which* case (best/average/worst).
Quicksort is `O(n log n)` average but `O(n²)` worst-case on adversarially
chosen (or already-sorted) input — a distinction that matters a lot if an
attacker controls your input.

---

## 19. Number Theory Basics

**Level:** Intermediate
**You'll use this for:** the mathematical foundation of every public-key cryptosystem (Ch. 24).

### Prime numbers

A prime has exactly two divisors: 1 and itself. Primes are the "atoms" of
multiplication — every integer > 1 has a unique prime factorization
(the Fundamental Theorem of Arithmetic). **Factoring a number back into its
primes is believed to be computationally hard for large numbers** — this
single asymmetry (easy to multiply two primes, hard to factor the product
back apart) is the entire basis of RSA (Ch. 24).

### GCD and the Euclidean algorithm

The greatest common divisor of two numbers, computed efficiently:

```python
def gcd(a, b):
    while b:
        a, b = b, a % b
    return a

gcd(48, 18)   # 6
```

**Tracing it by hand** — repeatedly replace the larger number with the
remainder of dividing it by the smaller, until the remainder hits 0:

```
gcd(48, 18):
  48 = 2×18 + 12    ->  gcd(48,18) = gcd(18,12)
  18 = 1×12 + 6     ->  gcd(18,12) = gcd(12,6)
  12 = 2×6  + 0     ->  remainder is 0, STOP

The last non-zero remainder is the GCD: 6.
```

**Why this works:** any number that divides both `48` and `18` must also
divide their difference (and therefore their remainder) — so
`gcd(48,18) = gcd(18,12) = gcd(12,6) = 6`. Each step shrinks the numbers
fast (roughly by half every two steps), which is exactly why this runs in
`O(log(min(a,b)))` instead of checking every possible common divisor.

This 2,300-year-old algorithm (Euclid, ~300 BCE) runs in `O(log(min(a,b)))`
time — remarkably fast — and is a core primitive inside RSA key generation
and modular inverse computation.

### Modular inverse and Fermat's Little Theorem

The modular inverse of `a` mod `n` is a number `x` such that `a·x ≡ 1
(mod n)` — "division" in modular arithmetic. **Fermat's Little Theorem**:
if `p` is prime and `a` is not a multiple of `p`, then `a^(p-1) ≡ 1
(mod p)`. This theorem is the basis for:
1. Fast primality testing (Fermat/Miller-Rabin primality tests).
2. Computing modular inverses efficiently.
3. Part of the mathematical machinery behind RSA's correctness proof.

```python
# Modular inverse via Python's built-in pow() with a negative exponent (3-arg pow)
a, n = 3, 11
inverse = pow(a, -1, n)
print(inverse)                 # 4, because 3*4 = 12 ≡ 1 (mod 11)
print((a * inverse) % n)       # 1
```

**Finding that inverse by hand** — for small numbers, just try
multiplying `a` by each candidate `x` from 1 upward until you hit `≡ 1
(mod n)`:

```
Find x such that 3x ≡ 1 (mod 11):
  3×1 = 3    mod 11 = 3    no
  3×2 = 6    mod 11 = 6    no
  3×3 = 9    mod 11 = 9    no
  3×4 = 12   mod 11 = 1    yes!  -> inverse of 3 mod 11 is 4
```

This brute-force search works fine for small `n`, but for cryptographic
sizes (hundreds of digits) it's infeasible — real implementations use the
**extended Euclidean algorithm**, which runs the GCD trace above
*backwards*, substituting each remainder equation back into the previous
one to express `1` as `3·x + 11·y`, and reads `x mod 11` off as the
inverse — same `O(log n)` speed as ordinary GCD, just tracked with extra
bookkeeping at each step.

### Euler's totient function

`φ(n)` counts how many integers from 1 to `n` are coprime to `n` (share no
common factor). For a prime `p`: `φ(p) = p - 1`. For `n = p·q` (product of
two distinct primes): `φ(n) = (p-1)(q-1)`. This exact formula is the
crux of RSA key generation (Ch. 24 walks through the full derivation).

### Real-life engineering ties

- **Security/Crypto:** this entire chapter is direct preparation for
  understanding *why* RSA, Diffie-Hellman, and DSA work — not just how to
  call the library functions (Ch. 24).
- **SWE:** hash function design borrows number-theoretic properties
  (good hash functions often use large primes to reduce clustering).
- **Network:** less direct, but TLS/SSH key exchange (used on every HTTPS
  connection) is number theory in production.

**More real-life examples:**
- Generating a GPG/PGP or TLS certificate key pair involves searching for
  large probable primes using fast primality tests (Miller-Rabin, built
  on Fermat's Little Theorem) — a process you can watch run
  (`openssl genrsa` pausing while it searches) every time a key is generated.
- NIST formally deprecated 1024-bit RSA keys around 2013-2015 specifically
  because advances in factoring algorithms and computing power eroded
  their safety margin — a dated, real-world consequence of the
  factoring-hardness assumption discussed in the common mistake below.
- Blockchain wallet addresses are derived from elliptic-curve public keys
  (Ch. 24), which rest on the same finite-field number theory as this
  chapter — losing the private key (a large number) means losing the
  funds forever, with no "reset password" option.
- Historical cryptanalysis of the Enigma machine at Bletchley Park
  exploited *structural* weaknesses (a letter could never encrypt to
  itself) rather than breaking the underlying math — an early lesson that
  a cipher's implementation details matter as much as its theoretical
  hardness.

### Common mistake

Believing "prime factorization is hard" is a proven mathematical fact — it
is **not proven**; it's a *widely believed, extensively tested*
computational hardness assumption. This distinction matters: it's why
"post-quantum" cryptography exists — Shor's algorithm, running on a
sufficiently large quantum computer, *would* break this assumption
efficiently, which is why NIST has standardized quantum-resistant
alternatives (ML-KEM, ML-DSA) as of 2024.

---

# Part 3 — Expert: Domain-Specific Deep Dives (Senior Level)

Parts 1–2 built a general-purpose toolkit that every engineer needs some of.
This Part is ten domain-specific deep dives instead — ML internals,
information theory, cryptography, quantitative security, control loops —
each aimed at a specific senior-level problem. Pick the chapters that match
your actual domain; nothing here assumes you've read the other nine.

## 20. Linear Algebra for AI/ML (SVD, PCA, Tensors)

**Level:** Expert
**You'll use this for:** dimensionality reduction, understanding embeddings, model compression (LoRA), and reading ML papers without glazing over.

### Tensors and broadcasting

A tensor is the generalization of scalar (0D) → vector (1D) → matrix (2D) →
`n`-D array. A batch of color images is a 4D tensor:
`[batch, height, width, channels]`. **Broadcasting** is the rule set that
lets numpy/PyTorch operate on tensors of different (but compatible) shapes
without explicit loops.

```python
import numpy as np
batch = np.random.rand(32, 128)      # 32 samples, 128 features each
bias = np.random.rand(128)           # one bias per feature
result = batch + bias                # broadcasts bias across all 32 rows
```
Misunderstanding broadcasting rules is one of the most common silent-bug
sources in ML code — shapes that are "compatible" by broadcasting rules
but not what you intended will run without error and produce wrong
numbers.

### Matrix rank and low-rank structure

The **rank** of a matrix is the number of linearly independent
rows/columns — intuitively, the "true dimensionality" of the information it
contains. A matrix can be `1000×1000` but have rank 10, meaning it's
"secretly" much simpler than its size suggests.

### Singular Value Decomposition (SVD)

Any matrix `A` (even non-square) can be decomposed:
```
A = U · Σ · Vᵀ
```
where `Σ` is a diagonal matrix of **singular values**, sorted largest to
smallest, representing how much "information" each component carries.
Keeping only the top `k` singular values gives the *best possible*
rank-`k` approximation of `A` — the mathematical foundation of
compression, denoising, and recommendation systems.

```python
A = np.random.rand(100, 50)
U, S, Vt = np.linalg.svd(A, full_matrices=False)
k = 10
A_approx = U[:, :k] @ np.diag(S[:k]) @ Vt[:k, :]
print("compression: kept", k, "of", len(S), "singular values")
```

### PCA — Principal Component Analysis, derived (not memorized)

PCA finds the directions (principal components) along which data varies
the most. Mechanically: **PCA = eigendecomposition of the covariance
matrix**, and its top components are the top singular vectors of the
(centered) data matrix — i.e., PCA is SVD applied to centered data.

**A 4-point PCA, by hand**, to see each step concretely. Points:
`(0,0), (2,2), (4,4), (6,6)` — all on a line, so PCA should find "the
diagonal direction" as the one that captures all the variance.

```
1. Center the data (subtract the mean):
   mean = (3, 3)
   centered points: (-3,-3), (-1,-1), (1,1), (3,3)

2. Build the covariance matrix. For 2D centered data with columns x, y:
     cov[0][0] = mean(x·x) = (9+1+1+9)/4 = 5
     cov[1][1] = mean(y·y) = (9+1+1+9)/4 = 5
     cov[0][1] = cov[1][0] = mean(x·y) = (9+1+1+9)/4 = 5
   Cov = [5 5]
         [5 5]

3. Eigendecompose Cov (same characteristic-equation method as Ch. 11):
   det(Cov - λI) = (5-λ)² - 25 = 0  =>  λ² - 10λ = 0  =>  λ = 0 or λ = 10

   λ=10's eigenvector is the direction (1,1) (normalized: 1/√2, 1/√2) —
   this is the "diagonal" principal component, carrying ALL the variance
   (the other eigenvalue, 0, means zero variance perpendicular to it,
   which matches the data lying exactly on a line).
```

That's the entire mechanism `np.linalg.eigh` runs internally, just at
whatever dimensionality your real dataset has — center, build the
covariance matrix, and find its eigenvectors.

```python
from numpy.linalg import eigh

def pca(X, k):
    X_centered = X - X.mean(axis=0)
    cov = np.cov(X_centered, rowvar=False)
    eigvals, eigvecs = eigh(cov)              # eigh: symmetric matrix, real eigenvalues
    order = np.argsort(eigvals)[::-1]
    top_k = eigvecs[:, order[:k]]
    return X_centered @ top_k                  # projected data

X = np.random.rand(200, 20)
reduced = pca(X, k=2)   # 20 dimensions -> 2, keeping max variance
```

### LoRA (Low-Rank Adaptation) — this math in production

Fine-tuning a large model normally updates a full weight matrix `W`
(potentially billions of parameters). **LoRA** observes that the *update*
`ΔW` needed for fine-tuning often has low intrinsic rank, so instead of
learning a full `ΔW`, you learn two small matrices `B` (n×r) and `A` (r×m)
where `r << n,m`, and use `ΔW = B·A`. This is a direct, production-scale
application of the low-rank idea from SVD above — it's why LoRA fine-tuning
uses a tiny fraction of the memory/compute of full fine-tuning.

**The arithmetic:** one attention projection in a 7B-class model is
`4096 × 4096`.

```
full ΔW:       4096 · 4096          = 16,777,216 trainable numbers
LoRA, r = 8:   8 · (4096 + 4096)    =     65,536 trainable numbers   (0.39%)
```

| Symbol | Meaning |
|---|---|
| `W` | the frozen pretrained weights (`n × m`) |
| `ΔW = B·A` | the learned update. `B` is `n × r` and starts at zero, `A` is `r × m` and starts random |
| `r` | the rank, typically 4–64. It caps how many independent "directions of change" the update can have |
| `α/r` | a scaling factor applied to `B·A` so that changing `r` doesn't change the learning rate |

Because `B` starts at zero, `ΔW = 0` at step 0 and fine-tuning begins
*exactly* at the pretrained model. That's a small detail of initialization
that makes LoRA stable.

### SVD, symbol by symbol, and by hand

```
A = U · Σ · Vᵀ          (A is m × n)
```

| Symbol | Shape | Meaning |
|---|---|---|
| `V` | `n × n` | orthonormal input directions (the "right singular vectors") |
| `Σ` | `m × n` | diagonal, with `σ₁ ≥ σ₂ ≥ … ≥ 0`, the stretch along each direction |
| `U` | `m × m` | orthonormal output directions (the "left singular vectors") |
| `σᵢ` | scalar | the `i`-th singular value, where `σᵢ² = ` eigenvalues of `AᵀA` |

Read right to left, **every matrix does three things**: rotate (`Vᵀ`),
stretch along the axes (`Σ`), rotate again (`U`). Any linear map, however
messy, is a rotation, a stretch and another rotation.

**By hand for `A = [[3, 0], [4, 5]]`:**

```
AᵀA = [3 4]·[3 0] = [25 20]
      [0 5] [4 5]   [20 25]

eigenvalues of AᵀA:  (25-λ)² - 400 = 0  ->  λ = 45 or 5
singular values:     σ₁ = √45 ≈ 6.708,  σ₂ = √5 ≈ 2.236
```

(Same `[[a,b],[b,a]]` pattern as Ch. 11, with eigenvalues `a ± b`.)

> **Theorem (Eckart–Young, 1936).** The best rank-`k` approximation of `A`,
> measured by the sum of squared errors (Frobenius norm), is the truncated SVD
> `Aₖ`, and its error is exactly the energy you dropped:
> ```
> ‖A - Aₖ‖²_F = σ²ₖ₊₁ + σ²ₖ₊₂ + ...
> ```

For our matrix, `‖A‖²_F = 9 + 0 + 16 + 25 = 50 = 45 + 5`. Keeping only
`σ₁` (rank 1) leaves an error of exactly `σ₂² = 5`, so **one number pair keeps
90% of the energy**. Real data is far more extreme. The singular values of
natural images, user–item ratings and LLM weight updates fall off a cliff,
which is why compression, recommendation and LoRA work at all.

### Attention, the equation behind every transformer

```
Attention(Q, K, V) = softmax( Q·Kᵀ / √d_k ) · V
```

| Symbol | Shape | Meaning |
|---|---|---|
| `Q` | `n × d_k` | **queries**: for each token, "what am I looking for?" |
| `K` | `n × d_k` | **keys**: for each token, "what do I contain?" |
| `V` | `n × d_v` | **values**: for each token, "what do I pass on if selected?" |
| `Q·Kᵀ` | `n × n` | every query dotted with every key: a table of match scores (Ch. 10) |
| `√d_k` | scalar | the scaling factor (explained below) |
| `softmax(…)` | `n × n` | each row becomes weights that sum to 1 (Ch. 13) |
| `… · V` | `n × d_v` | each token's output is a weighted average of all values |

Attention is a **soft dictionary lookup**. A Python dict returns the one
value whose key matches exactly. Attention returns a *blend* of all values,
weighted by how well each key matches.

**A tiny example by hand.** One query `q = (1, 0)`, three tokens with keys
`(1, 0)`, `(0, 1)`, `(1, 1)` and scalar values `10, 20, 90`, with `d_k = 2`:

```
scores  = q·kᵢ        = (1, 0, 1)
scaled  = / √2        = (0.707, 0, 0.707)
weights = softmax     = (0.401, 0.198, 0.401)
output  = 0.401·10 + 0.198·20 + 0.401·90  ≈  44.1
```

The two keys that match the query get twice the weight of the one that
doesn't.

**Why divide by `√d_k`? The variance argument (Ch. 8 and 14 earning their keep).**
If the entries of `q` and `k` are independent with mean 0 and variance 1,
then `q·k = Σᵢ qᵢkᵢ` is a sum of `d_k` terms each with variance 1, so
`Var(q·k) = d_k` and its standard deviation is `√d_k`. With `d_k = 64`, raw
scores have a standard deviation of about 8. Softmax of scores like
`(8, 0, -8)` is `(0.9997, 0.0003, 0.0000)`, which is nearly one-hot. In that
regime the softmax gradient `p - y` (Ch. 13) is almost zero and learning
stalls. Dividing by `√d_k = 8` restores a standard deviation of 1, and the
same scores become `(0.67, 0.24, 0.09)`. The `√d_k` that appears in every
transformer paper is a probability fact used to fix a calculus problem.

### High-dimensional geometry: why embeddings work

Intuition from 2D and 3D fails badly in 768 dimensions. Draw two random unit
vectors in `ℝᵈ`. Their cosine similarity has mean 0 and standard deviation
about `1/√d`:

```
d = 3:    std ≈ 0.58     random vectors often look "similar"
d = 768:  std ≈ 0.036    random vectors are almost exactly perpendicular
```

In high dimensions **almost every pair of random directions is nearly
orthogonal**, so a space can hold an exponential number of nearly
independent concepts. That's why a cosine similarity of 0.3 between two
embeddings is a meaningful signal and not noise. It's also why
approximate nearest-neighbour indexes (HNSW, IVF-PQ in FAISS and vector databases)
can prune most of the space quickly.

### Networking in the wild: PCA finds anomalies across a whole network

Stack a backbone network's link loads over time into a matrix: rows are 5-minute
intervals, columns are links. Normal traffic is driven by a few shared
forces (time of day, weekly cycle, a handful of big customers), so the
matrix is close to **low rank**, and a few principal components explain most of
the variance. Lakhina, Crovella and Diot (SIGCOMM 2004) split each time step into
the part explained by the top components (the "normal subspace") and the
**residual**. A DDoS, a routing change or an outage on a few links shows up
as a residual spike, even when no single link crosses its own alert
threshold. Network-wide anomaly detection is Eckart–Young applied to
SNMP counters.

### Real-life engineering ties

- **AI:** PCA for feature reduction/visualization, SVD for recommendation
  systems (matrix factorization: Users × Items ≈ low-rank approximation),
  LoRA/QLoRA for efficient fine-tuning, embedding compression.
- **Security:** PCA-based anomaly detection (project traffic/log features
  onto principal components; large residual = anomalous, doesn't fit the
  "normal" subspace).
- **Network:** less direct, but spectral graph analysis (eigenvectors of a
  network's Laplacian matrix) detects community structure/partitions.

**More real-life examples:**
- The 2006-2009 Netflix Prize competition was won using
  matrix-factorization techniques directly descended from SVD, applied to
  the giant sparse users x movies ratings matrix — arguably the event that
  popularized this math for recommendation systems industry-wide.
- Visualizing high-dimensional word or sentence embeddings (e.g., a
  768-dimensional BERT embedding space) is typically done by reducing to
  2-3 dimensions with PCA or the related t-SNE/UMAP techniques, so humans
  can look at a scatter plot.
- JPEG image compression uses the Discrete Cosine Transform, a close
  mathematical relative of SVD, to concentrate an image's information into
  a few significant coefficients and discard the rest — the same
  "keep the top-k components" idea as the SVD compression example above.
- Hugging Face's widely used PEFT library implements LoRA and QLoRA in
  production, letting engineers fine-tune multi-billion-parameter models
  on a single consumer GPU by training only the small low-rank matrices
  described in this chapter.

### Common mistake

Running PCA without **centering** (subtracting the mean) first — PCA is
mathematically defined on centered data; skipping this step gives a result
dominated by the mean offset, not the actual variance structure.

---

## 21. Optimization & Backpropagation

**Level:** Expert
**You'll use this for:** actually understanding how a neural network trains, instead of treating `.fit()`/`loss.backward()` as magic.

### Gradient descent, from the ground up

To minimize a loss function `L(w)`, repeatedly step in the direction
opposite the gradient (the direction of steepest ascent, Ch. 13):

```
w_new = w_old - η · ∇L(w_old)
```
`η` (eta) is the **learning rate** — too large and you overshoot/diverge,
too small and training crawls.

```python
def gradient_descent(grad_fn, w0, lr=0.1, steps=50):
    w = w0
    for _ in range(steps):
        w = w - lr * grad_fn(w)
    return w

# Minimize f(w) = (w-3)^2, whose derivative is 2(w-3)
w_min = gradient_descent(grad_fn=lambda w: 2*(w-3), w0=0.0)
print(w_min)   # converges to ~3.0
```

### Full worked backpropagation example (2 layers, by hand)

Network: `x → [w1] → h → relu → [w2] → y_pred`, loss `L = (y_pred - target)²`.

```python
# Forward pass
x, target = 2.0, 10.0
w1, w2 = 1.5, 2.0

h = w1 * x                 # h = 3.0
h_relu = max(0, h)         # h_relu = 3.0
y_pred = w2 * h_relu       # y_pred = 6.0
loss = (y_pred - target)**2   # loss = 16.0

# Backward pass (chain rule, Ch. 13, applied link by link)
dL_dy = 2 * (y_pred - target)          # dL/dy_pred = 2*(6-10) = -8.0
dy_dw2 = h_relu                          # dy_pred/dw2 = h_relu = 3.0
dL_dw2 = dL_dy * dy_dw2                  # -8.0 * 3.0 = -24.0

dy_dh_relu = w2                          # dy_pred/dh_relu = w2 = 2.0
dh_relu_dh = 1.0 if h > 0 else 0.0       # ReLU derivative = 1 (since h=3>0)
dL_dh = dL_dy * dy_dh_relu * dh_relu_dh  # -8.0 * 2.0 * 1.0 = -16.0

dh_dw1 = x                               # dh/dw1 = x = 2.0
dL_dw1 = dL_dh * dh_dw1                  # -16.0 * 2.0 = -32.0

print("dL/dw1 =", dL_dw1, " dL/dw2 =", dL_dw2)
# Update: w1_new = w1 - lr*dL_dw1, w2_new = w2 - lr*dL_dw2
```
**This hand-computed example is literally what `loss.backward()` in
PyTorch does automatically**, for every parameter, across every layer,
via a computation graph — reverse-mode automatic differentiation is
just this chain-rule bookkeeping, done exactly, at scale.

### Why deep networks need more than plain gradient descent

- **Vanishing/exploding gradients:** chaining many derivatives through the
  chain rule means multiplying many terms together — if each is < 1, the
  product shrinks toward 0 exponentially with depth (vanishing); if each
  is > 1, it explodes. This motivates activation function choice (ReLU
  over sigmoid) and architectural fixes (residual/skip connections,
  normalization layers).
- **Momentum:** accumulates a running average of past gradients to smooth
  out noisy steps and accelerate through consistent directions.
- **Adam:** combines momentum with a per-parameter adaptive learning rate
  based on the running variance of gradients — the default optimizer for
  most modern deep learning.

### Convexity, and why it matters

A function is **convex** if a straight line between any two points on it
never dips below the curve — convex functions have exactly one global
minimum, guaranteeing gradient descent finds *the* best answer. Most
real deep learning loss landscapes are **non-convex** (many local minima
and saddle points) — which is why training is empirically finicky and why
techniques like momentum, learning-rate schedules, and good initialization
matter so much in practice.

### Lagrange multipliers (constrained optimization)

To minimize `f(x)` subject to a constraint `g(x) = 0`, form the
Lagrangian `L(x, λ) = f(x) - λ·g(x)` and solve `∇L = 0`. This is the
mechanism behind Support Vector Machines (maximize margin subject to
correct-classification constraints) and appears throughout constrained
resource-allocation problems (e.g., "minimize latency subject to a fixed
compute budget").

### The learning-rate speed limit: `η < 2/λ_max`

Why does a learning rate that's too large diverge? Take the simplest loss,
a 1-D bowl with curvature `λ`:

```
L(w) = ½ · λ · w²       L'(w) = λ·w       (minimum at w = 0)

update:  w_new = w - η·λ·w = (1 - η·λ) · w
```

| Symbol | Meaning |
|---|---|
| `λ` | curvature, the second derivative `L''`. In `n` dimensions, an eigenvalue of the Hessian (Ch. 13) |
| `η` | learning rate |
| `1 - ηλ` | the factor the error is multiplied by **every step** |

After `t` steps, `w_t = (1 - ηλ)ᵗ · w₀`. That's a geometric sequence (Ch. 6),
so it converges if and only if `|1 - ηλ| < 1`, that is, **`0 < η < 2/λ`**.
With `λ = 10`:

| `η` | factor `1 - ηλ` | What the loss curve does |
|---|---|---|
| 0.05 | 0.5 | smooth, fast convergence |
| 0.15 | -0.5 | converges, overshooting back and forth |
| 0.19 | -0.9 | barely converges, violent oscillation |
| 0.25 | -1.5 | **diverges**: the error grows 1.5× per step until you see NaN |

In a real network every Hessian eigenvalue imposes its own limit, and the
**sharpest direction** (`λ_max`) sets the speed limit for all of them. But the
flattest direction (`λ_min`) then crawls at a factor of `1 - ηλ_min ≈ 1`.
The ratio `κ = λ_max/λ_min`, the **condition number**, measures how hard the
problem is. Plain gradient descent needs on the order of `κ` steps. Feature
standardization (Ch. 9 z-scores), BatchNorm/LayerNorm, momentum and Adam are
all, at heart, ways to reduce `κ` or work around it. Researchers have observed
that training tends to drive `λ_max` up until it hovers right at `2/η`
("the edge of stability", Cohen et al., 2021), so this formula is a live
dynamic in real training runs.

### Momentum and Adam, every symbol explained

**Momentum** (the "heavy ball"):

```
v ← β · v + ∇L(w)
w ← w - η · v
```

`v` is a velocity that accumulates past gradients, an EWMA from Ch. 9 with
`β ≈ 0.9`, so about 10 steps of memory. Along consistent directions,
steps add up and the effective learning rate is about `η/(1-β) = 10η`.
Across oscillating directions the opposite signs cancel. That accelerates
the slow, flat directions without speeding up the sharp ones, which is
exactly the fix for high `κ`.

**Adam** (Kingma & Ba, 2014), the default optimizer for transformers:

```
m ← β₁·m + (1-β₁)·g                 1st moment: EWMA of gradients
v ← β₂·v + (1-β₂)·g²                2nd moment: EWMA of squared gradients
m̂ = m / (1 - β₁ᵗ)                   bias correction
v̂ = v / (1 - β₂ᵗ)
w ← w - η · m̂ / (√v̂ + ε)
```

| Symbol | Typical value | Meaning |
|---|---|---|
| `g` | | the current gradient `∇L(w)` |
| `m` | starts at 0 | smoothed gradient: *which way* to go (momentum) |
| `v` | starts at 0 | smoothed squared gradient: *how big* gradients usually are for this parameter |
| `β₁`, `β₂` | 0.9, 0.999 | memory of about 10 and about 1,000 steps |
| `t` | 1, 2, 3, … | step count, used only in bias correction |
| `ε` | 1e-8 | stops division by zero |
| `m̂/√v̂` | ≈ ±1 | a **unit-less** step: a direction divided by its typical size |

**Why bias correction?** `m` and `v` start at 0, so early on they're
dragged toward 0. At step 1, `m = 0.1·g`, which is 10× too small. Dividing by
`1 - 0.9¹ = 0.1` fixes it exactly.

**Worked step 1, gradient `g = 0.5`:**

```
m = 0.1 · 0.5   = 0.05          m̂ = 0.05    / (1 - 0.9)   = 0.5
v = 0.001 · 0.25 = 0.00025      v̂ = 0.00025 / (1 - 0.999) = 0.25
step = η · 0.5 / √0.25 = η · 1.0
```

Try `g = 0.0005` instead and the step is *still* `η · 1.0`. Adam's step size
is roughly `η` whatever the gradient's scale. That's why one learning rate
works across embedding tables, attention weights and layer norms whose
gradients differ by orders of magnitude, and why `η` in Adam reads directly as
"the maximum change per parameter per step."

### Lagrange multipliers, worked: how should flows share a network?

A problem from networking that shows what a Lagrange multiplier *means*.
Two links, each with capacity 10 Mbit/s. Flow 0 crosses both links. Flow 1
uses only link A, and flow 2 uses only link B:

```
flow 0:  ──[ link A ]──[ link B ]──
flow 1:  ──[ link A ]──
flow 2:              ──[ link B ]──
```

Frank Kelly (1998) proposed choosing rates `xᵢ` to maximize total
"utility" `Σ log(xᵢ)`. The `log` gives diminishing returns, so starving any
flow is very costly. This is **proportional fairness**:

```
maximize    log x₀ + log x₁ + log x₂
subject to  x₀ + x₁ ≤ 10      (link A, multiplier λ_A)
            x₀ + x₂ ≤ 10      (link B, multiplier λ_B)
```

**The Lagrangian** folds each constraint into the objective, weighted by a
multiplier:

```
ℒ = log x₀ + log x₁ + log x₂ - λ_A·(x₀ + x₁ - 10) - λ_B·(x₀ + x₂ - 10)
```

| Symbol | Meaning |
|---|---|
| `ℒ` | the Lagrangian: objective minus "price × overuse" for each constraint |
| `λ_A`, `λ_B` | Lagrange multipliers, the **price** per unit of bandwidth on each link |
| `∂ℒ/∂xᵢ = 0` | at the optimum, no flow can gain by changing its rate |

Set each partial derivative to zero:

```
∂ℒ/∂x₁ = 1/x₁ - λ_A = 0              ->  x₁ = 1/λ_A
∂ℒ/∂x₂ = 1/x₂ - λ_B = 0              ->  x₂ = 1/λ_B
∂ℒ/∂x₀ = 1/x₀ - λ_A - λ_B = 0        ->  x₀ = 1/(λ_A + λ_B)
```

**Flow 0 pays the price of every link it crosses.** By symmetry
`λ_A = λ_B = λ`, both links are full, so `x₀ + x₁ = 1/(2λ) + 1/λ = 10`, giving
`λ = 0.15`:

```
x₀ = 1/0.30 = 3.33 Mbit/s       (the 2-hop flow)
x₁ = x₂ = 1/0.15 = 6.67 Mbit/s  (the 1-hop flows)
```

Compare max-min fairness, where everyone gets 5. Proportional fairness gives
the long flow less, because it uses twice the network resources. Kelly showed
that TCP-style congestion control, where each flow backs off when it sees loss or
delay (the "price") on its path, is a **distributed algorithm solving this
optimization**. The multipliers are the congestion signals, and the links
compute them without any central coordinator. This is also why long
multi-hop TCP flows get less throughput than short ones in practice. ML uses
the identical machinery: the SVM's support vectors are the training points
whose Lagrange multipliers are non-zero.

### Real-life engineering ties

- **AI:** this chapter is the actual mechanism of `model.fit()` — knowing
  it lets you debug training instability (exploding loss, NaN gradients,
  stuck training) as a first-principles problem instead of randomly
  changing hyperparameters.
- **SWE/Network:** gradient-based tuning shows up in auto-scaling
  controllers and adaptive rate limiters that continuously adjust a
  parameter to minimize an error signal (closely related to Ch. 29,
  control theory).
- **Security:** adversarial examples (inputs crafted to fool an ML model)
  are literally found via gradient ascent on the model's own loss function
  with respect to the *input* instead of the weights — understanding
  backprop is a prerequisite for understanding adversarial ML attacks
  (Ch. 28 connects this to game theory).

**More real-life examples:**
- The 2012 AlexNet result that kicked off the modern deep learning boom
  was, mathematically, nothing more than the gradient descent and
  backpropagation in this chapter run at then-unprecedented scale on
  GPUs — the math didn't change, the compute did.
- Modern large language model training runs use carefully engineered
  learning-rate *schedules* (a warmup period followed by cosine decay)
  specifically to avoid the too-high-divergence and too-low-slowness
  failure modes described in this chapter.
- ResNet (2015) introduced "skip connections" specifically to fix the
  vanishing-gradient problem in very deep networks — a direct, famous
  architectural response to the chain-rule multiplication problem
  explained above.
- Researchers have physically demonstrated adversarial attacks by placing
  small stickers on a stop sign that reliably fool a vision model into
  classifying it as a speed-limit sign (Eykholt et al., 2018) — a
  real-world instance of the gradient-ascent-on-the-input attack.

### Common mistake

Setting the learning rate too high "to train faster" — this is the single
most common cause of a loss curve that diverges or oscillates instead of
converging; always look at the loss curve shape, not just the final number.

---

## 22. Statistics for ML & Experimentation

**Level:** Expert
**You'll use this for:** trusting (or correctly distrusting) an A/B test, a model evaluation, or a paper's claimed result.

### Maximum Likelihood Estimation (MLE) vs. MAP

**MLE** picks the parameter values that make the observed data most
probable: `θ_MLE = argmax P(data | θ)`. Training a model by minimizing
cross-entropy loss **is** MLE under a specific probabilistic assumption
about the output distribution — it's not a separate technique, it's the
same math wearing a different name.

**MAP (Maximum A Posteriori)** adds a prior belief about `θ`:
`θ_MAP = argmax P(data|θ)·P(θ)`. Adding L2 regularization to a loss
function is mathematically equivalent to assuming a Gaussian prior on the
weights — regularization isn't an ad-hoc trick, it's Bayesian MAP
estimation in disguise.

### Hypothesis testing and p-values

A p-value answers: "if there were truly no effect (the null hypothesis),
how likely would data this extreme (or more extreme) be?" A small p-value
(conventionally < 0.05) is evidence against the null — **it is not** "the
probability the null hypothesis is true," a near-universal misreading.

```python
from scipy import stats

# A/B test: control conversion 100/1000, variant 130/1000
table = [
    [100, 900],   # control: converted, did not convert
    [130, 870],   # variant: converted, did not convert
]
chi2, p_value, dof, expected = stats.chi2_contingency(table)
print(p_value)   # if < 0.05, the observed gap is unlikely under "no difference"
```

For binary outcomes like conversion/no-conversion, use a proportion test
or chi-square test like this. A t-test is more natural for continuous outcomes
such as latency, revenue per user, or time spent.

**Where the t-statistic itself comes from**, computed by hand for two
small samples (`A = [10, 12, 9]`, `B = [15, 14, 16]`):

```
1. Mean of each group:
     mean_A = (10+12+9)/3 = 10.33
     mean_B = (15+14+16)/3 = 15.0

2. Standard deviation of each group (Ch. 9's method — deviations, square,
   average, sqrt):
     std_A ≈ 1.53      std_B ≈ 1.0

3. Standard error of the difference (roughly, for similar sample sizes):
     SE ≈ sqrt(std_A²/n_A + std_B²/n_B) = sqrt(1.53²/3 + 1.0²/3) ≈ 1.02

4. t-statistic = (difference in means) / (standard error):
     t = (15.0 - 10.33) / 1.02 ≈ 4.58
```

A `t` far from 0 (like 4.58) says "the gap between the two means is large
relative to how much the data naturally jitters within each group" — the
`p_value` is just that `t` translated into "how often would a gap this
large happen by pure chance," via the t-distribution's known shape (the
one calculation in this chapter genuinely worth leaving to a library,
since it requires integrating a probability density).

### Confidence intervals

A 95% confidence interval means: "if we repeated this experiment many
times, 95% of the intervals we'd construct this way would contain the true
value" — **not** "there's a 95% chance the true value is in this specific
interval" (a frequentist/Bayesian distinction that trips up even
experienced engineers).

### p-hacking and multiple comparisons — the trap that invalidates most naive A/B testing

If you test 20 independent metrics at `p < 0.05` significance, you expect
**~1 false positive by pure chance alone** — even with zero real effect.

```python
import numpy as np
from scipy import stats

np.random.seed(0)
false_positives = 0
for _ in range(20):
    # Two random samples with NO real difference
    a = np.random.normal(0, 1, 1000)
    b = np.random.normal(0, 1, 1000)
    _, p = stats.ttest_ind(a, b)
    if p < 0.05:
        false_positives += 1
print(false_positives)   # expect ~1, purely from testing 20 things at once
```
This is *the* mathematical justification for **Bonferroni correction**
(divide your significance threshold by the number of comparisons) and for
being deeply skeptical of any "we tested 30 metrics and found 2 significant
results" dashboard.

### Cross-entropy loss and KL divergence, the ML-training connection

```
Cross-entropy(p, q) = -Σ p(x)·log(q(x))
```
where `p` is the true label distribution and `q` is the model's predicted
distribution. Minimizing cross-entropy loss during classifier training is
mathematically minimizing the **KL divergence** between the true and
predicted distributions (Ch. 23 covers KL divergence directly) — this is
*why* cross-entropy, not raw accuracy, is the loss function, even though
accuracy is the metric you actually care about: it's differentiable and it
directly optimizes distributional closeness.

### MLE, derived with calculus: the coin and the bell

**Bernoulli MLE.** You observe `k` failures in `n` independent requests. What
failure rate `p` makes that most likely?

```
likelihood:       ℒ(p) = pᵏ · (1-p)ⁿ⁻ᵏ
log-likelihood:   ℓ(p) = k·ln p + (n-k)·ln(1-p)          (logs turn products into sums, Ch. 5)
derivative:       ℓ'(p) = k/p - (n-k)/(1-p) = 0
solve:            k(1-p) = (n-k)p   ->   p̂ = k/n
```

| Symbol | Meaning |
|---|---|
| `ℒ(p)` | likelihood: the probability of the data you saw, as a function of the unknown `p` |
| `ℓ(p)` | log-likelihood, which has the same maximum but is easier to differentiate |
| `p̂` | "p hat", the estimate. Hats mean "estimated from data" |

The "obvious" estimate, failures ÷ total, *is* the maximum-likelihood
estimate. Calculus confirms the intuition.

**Gaussian MLE gives you MSE.** Assume each target is the model's prediction
plus Gaussian noise: `yᵢ = f(xᵢ; θ) + noise`, with noise `~ N(0, σ²)`. The
log-likelihood is:

```
ℓ(θ) = Σᵢ ln[ 1/(σ√(2π)) · e^(-(yᵢ - f(xᵢ;θ))² / (2σ²)) ]
     = constant  -  (1/(2σ²)) · Σᵢ (yᵢ - f(xᵢ; θ))²
```

Maximizing `ℓ` is the same as **minimizing the sum of squared errors**. MSE
loss, least squares (Ch. 10) and linear regression are all MLE under a
Gaussian-noise assumption. Change the noise assumption and the loss changes
with it: Laplace noise gives absolute error (MAE), and Bernoulli outputs give
cross-entropy. **Choosing a loss function is choosing a probability model.**

### The two-proportion z-test, by hand

Here's the A/B test above (control 100/1000, variant 130/1000) done by hand.
Under the null hypothesis "no difference", both groups share one pooled rate:

```
p̂_pool = (100 + 130) / (1000 + 1000) = 0.115

SE = √( p̂_pool · (1 - p̂_pool) · (1/n₁ + 1/n₂) )
   = √( 0.115 · 0.885 · 0.002 ) = √0.0002036 ≈ 0.01427

z = (p̂₂ - p̂₁) / SE = (0.13 - 0.10) / 0.01427 ≈ 2.10

p-value = 2 · P(Z > 2.10) ≈ 0.035
```

| Symbol | Meaning |
|---|---|
| `p̂₁`, `p̂₂` | observed conversion rates in control and variant |
| `n₁`, `n₂` | sample sizes |
| `SE` | standard error of the difference, from the Bernoulli variance `p(1-p)` (Ch. 8) and the CLT (Ch. 14) |
| `z` | how many standard errors apart the two rates are |
| `2 · P(Z > z)` | two-sided: a difference this large in *either* direction |

`scipy`'s `chi2_contingency` in the code above prints about 0.042, not 0.035,
because it applies Yates' continuity correction by default. Pass
`correction=False` and it matches the hand calculation exactly. A 2×2
chi-square test *is* this z-test, squared (`χ² = z²`).

### Confidence intervals, the formula

```
p̂ ± z_(α/2) · √( p̂(1 - p̂) / n )
```

| Symbol | Meaning |
|---|---|
| `z_(α/2)` | the critical value: 1.96 for 95%, 2.576 for 99% |
| `√(p̂(1-p̂)/n)` | standard error of a proportion |

**Example:** variant at 130/1000:
`0.13 ± 1.96 · √(0.13 · 0.87 / 1000) = 0.13 ± 0.021 = [10.9%, 15.1%]`.
The `√n` in the denominator is the CLT's law of diminishing returns again:
to halve the interval's width you need **4×** the traffic.

### How many users does an A/B test need? Power analysis

Before you launch, decide the smallest effect worth detecting, then solve for `n`:

```
n per group = (z_(α/2) + z_β)² · [ p₁(1-p₁) + p₂(1-p₂) ] / δ²
```

| Symbol | Meaning | Typical value |
|---|---|---|
| `α` | false-positive rate you accept (Type I error) | 0.05, so `z_(α/2) = 1.96` |
| `1 - β` | **power**: the chance of detecting a real effect | 0.80, so `z_β = 0.84` |
| `p₁`, `p₂` | baseline rate and the rate you hope to reach | 10% → 11% |
| `δ` | `p₂ - p₁`, the minimum detectable effect | 0.01 |

```
n = (1.96 + 0.84)² · (0.10·0.90 + 0.11·0.89) / 0.01²
  = 7.85 · 0.1879 / 0.0001  ≈  14,750 users per group
```

Detecting a 1-point lift on a 10% baseline takes about 30,000 users in
total. Because of `δ²`, detecting **half** that lift takes **4×** as many.
This is why small companies can't run the same experiments as large ones,
and why "we ran it for a day and it looked good" usually means
"underpowered."

### Networking in the wild: the rule of three

You send 3,000 probe packets across a link and **none** are lost. Is the loss
rate zero? No, but you can bound it. The 95% upper bound `p` is the rate at
which seeing zero losses still has a 5% chance:

```
(1 - p)ⁿ = 0.05   ->   n·ln(1 - p) = ln 0.05   ->   n·(-p) ≈ -3.0   ->   p ≈ 3/n
```

(using `ln(1-p) ≈ -p` from Ch. 12's Taylor expansion, and `-ln 0.05 ≈ 3.0`)

**Zero failures in `n` trials means the true rate is below `3/n` with 95%
confidence.** For 3,000 clean probes, loss is below 0.1%. If your SLO is
0.01% loss, 3,000 probes can't verify it. You need 30,000. The same rule
answers "we ran 500 red-team prompts and the model never leaked the system
prompt". That only shows the leak rate is below `3/500 = 0.6%`, which is not
zero.

### Real-life engineering ties

- **AI:** correctly interpreting model evaluation metrics, understanding
  why regularization works probabilistically, designing experiments that
  won't be invalidated by multiple-comparisons issues.
- **SWE/Product:** every feature-flag A/B test at a company runs on this
  math — misreading a p-value here directly costs the business money
  (shipping "winning" features that were actually noise).
- **Security:** evaluating a new detection rule/model requires the same
  rigor — "we tested it and false positives dropped" needs a confidence
  interval, not just a single before/after number, especially against
  naturally noisy attack-volume data.

**More real-life examples:**
- The 1973 UC Berkeley graduate admissions data is a famous, documented
  real-world case of **Simpson's Paradox** — the university appeared to
  favor male applicants overall, but broken down by department, most
  departments individually favored female applicants; the aggregate trend
  reversed once properly segmented.
- The ongoing "replication crisis" in psychology and biomedical research
  (many high-profile 2010s replication attempts failing) is widely
  attributed in part to the multiple-comparisons and p-hacking problems
  demonstrated in this chapter's code example.
- Large tech companies with heavy experimentation cultures (Booking.com,
  Airbnb, Microsoft's ExP platform) run thousands of simultaneous A/B
  tests and invest in statistical infrastructure (sequential testing,
  false-discovery-rate control) to guard against inflated false-positive
  rates.
- Clinical drug trials conventionally require `p < 0.05` for a single
  pre-registered primary endpoint specifically *because* pre-registration
  prevents the multiple-comparisons/p-hacking trap — a regulatory rule
  that is a direct, high-stakes application of this chapter's math.

### Common mistake

**Peeking** — checking A/B test results repeatedly and stopping "as soon as
it's significant." This inflates the false-positive rate far above the
nominal 5%, because you're implicitly running many sequential tests
(similar to the multiple-comparisons problem above). Fix: pre-register a
sample size/stopping rule, or use sequential testing methods designed for
this (e.g., always-valid p-values).

---

## 23. Information Theory

**Level:** Expert
**You'll use this for:** password/key entropy, compression limits, ML loss functions, and network channel capacity.

### Shannon entropy — quantifying uncertainty/information

```
H(X) = -Σ p(x) · log2(p(x))
```
Entropy measures, in **bits**, the average "surprise" of a random
variable. A fair coin has `H = 1` bit (maximally uncertain). A coin that
always lands heads has `H = 0` bits (no uncertainty, no information
gained by observing it).

```python
import numpy as np

def entropy(probs):
    probs = np.array(probs)
    probs = probs[probs > 0]   # avoid log(0)
    return -np.sum(probs * np.log2(probs))

entropy([0.5, 0.5])          # 1.0 bit  -- fair coin
entropy([0.99, 0.01])        # 0.08 bits -- barely any uncertainty
entropy([0.25]*4)            # 2.0 bits  -- uniform over 4 outcomes
```

**Computing `entropy([0.25]*4)` by hand** — a fair 4-sided die, each face
equally likely:

```
H = -Σ p(x)·log2(p(x))

  = -(0.25·log2(0.25)) × 4 terms, since all 4 probabilities are equal

  log2(0.25) = log2(1/4) = -2      (2^-2 = 0.25)

  each term: -(0.25 × -2) = 0.5
  sum over 4 outcomes: 0.5 × 4 = 2.0 bits
```

Sanity check against Ch. 5: a uniform choice among 4 outcomes needs
exactly `log2(4) = 2` bits to identify which one occurred — entropy for a
*uniform* distribution always equals `log2(number of outcomes)`, which is
why the fair-coin case above gives exactly `log2(2) = 1` bit.

### Password entropy — the direct security application

If a password is drawn uniformly at random from a keyspace of size `N`,
its entropy is `log2(N)` bits. This is the *correct*, rigorous way to
quantify password strength — not "has a symbol and a number."

```python
import math
def password_entropy_bits(alphabet_size, length):
    return length * math.log2(alphabet_size)

print(password_entropy_bits(26, 8))    # lowercase, 8 chars: 37.6 bits
print(password_entropy_bits(95, 16))   # full ASCII printable, 16 chars: 105.2 bits
print(password_entropy_bits(2048, 6))  # 6-word Diceware passphrase (2048-word list): 66 bits
```
**Why a 6-word Diceware passphrase (66 bits) can beat an 8-character
"complex" password (~38–52 bits depending on charset assumptions):** entropy
is about the *size of the space an attacker must search*, not surface
complexity — length from a large-but-memorable word list often wins.

### Cross-entropy and KL divergence, formally

```
Cross-entropy: H(p, q) = -Σ p(x)·log2(q(x))
KL divergence: D_KL(p || q) = H(p, q) - H(p) = Σ p(x)·log2(p(x)/q(x))
```
KL divergence measures how different distribution `q` is from `p` — it's
**not symmetric** (`D_KL(p||q) ≠ D_KL(q||p)`), which matters: it means
"how surprised is a model expecting q, when reality is p" is a directional
question. This exact quantity is the loss function minimized during
classifier training (Ch. 22) and is used in diffusion models, variational
autoencoders (VAEs), and reinforcement learning (e.g., PPO's KL penalty
term that keeps policy updates from moving too far in one step).

### Mutual information — for feature selection

```
I(X;Y) = H(X) - H(X|Y)
```
"How many bits of uncertainty about `X` are removed once you know `Y`."
Used directly in ML feature selection: pick the features with the highest
mutual information with the target label — i.e., the features that are
*actually* informative, beyond simple correlation (mutual information
captures non-linear relationships that correlation coefficients miss).

### Shannon-Hartley theorem — the hard limit on channel capacity

```
C = B · log2(1 + S/N)
```
`C` = maximum error-free channel capacity (bits/sec), `B` = bandwidth
(Hz), `S/N` = signal-to-noise ratio. This is a **hard physical limit** —
no amount of clever engineering can transmit faster than this over a given
channel with a given noise level. It's why increasing Wi-Fi throughput
requires either more bandwidth (wider channels, Ch. 27) or better SNR
(closer to the router, less interference), and it directly explains why
"5G is faster" partly comes down to using wider channels and better
modulation to approach this Shannon limit more closely.

```python
import math
def shannon_capacity(bandwidth_hz, snr_db):
    snr_linear = 10 ** (snr_db / 10)
    return bandwidth_hz * math.log2(1 + snr_linear)

print(shannon_capacity(20_000_000, 30))  # 20 MHz channel, 30dB SNR -> bits/sec
```

### Real-life engineering ties

- **Security:** rigorous password/key entropy calculation (above); also
  the theoretical basis for evaluating whether a random number generator
  is producing enough real entropy (`/dev/random` blocking behavior on
  Linux exists precisely because of this concern).
- **AI:** cross-entropy loss (classification), KL divergence
  (VAEs/diffusion/RLHF), mutual information (feature selection,
  information bottleneck theory of deep learning).
- **Network:** Shannon-Hartley is the fundamental ceiling for every
  wireless and wired link capacity planning exercise (Ch. 25, 27).
- **SWE:** compression algorithms (gzip, Huffman coding) are literally
  bounded by the source's entropy — you cannot losslessly compress random
  (maximum-entropy) data, which is why compressing already-compressed or
  already-encrypted data does nothing.

**More real-life examples:**
- The Electronic Frontier Foundation publishes a curated 7,776-word
  Diceware wordlist specifically so people can generate the kind of
  high-entropy, memorable passphrases calculated in this chapter using
  physical dice rolls.
- Linux's `/dev/random` and `/dev/urandom` are built around entropy
  *estimation* — the kernel tracks how many bits of real-world
  unpredictability (keyboard timing, disk timing, interrupts) have
  accumulated, directly applying Shannon entropy to decide whether output
  is safe for cryptographic keys.
- "Zip bomb" files (a small compressed file that expands to petabytes)
  exploit extremely low-entropy, highly repetitive data — since
  compression is bounded by source entropy, and repetitive data has
  almost none, it compresses (and decompresses) at an extreme ratio,
  which attackers weaponize as a denial-of-service payload.
- Shannon's 1948 source coding theorem set the *theoretical* best-possible
  compression ratio for a given data source decades before formats like
  gzip, Zstandard, or Brotli were built to approach that limit.

### Common mistake

Treating password "complexity rules" (must have 1 uppercase, 1 digit, 1
symbol) as equivalent to high entropy — such rules often *reduce* the
effective keyspace (everyone puts the digit at the end, the symbol is
almost always `!`), producing lower real entropy than the naive `log2(N^L)`
calculation suggests. Real-world password entropy is lower than
theoretical entropy because human choices aren't uniformly random.

---

## 24. Cryptographic Mathematics

**Level:** Expert
**You'll use this for:** actually understanding *why* RSA, Diffie-Hellman, and elliptic-curve crypto are secure — not just calling `openssl`.

### Modular exponentiation (the core operation of public-key crypto)

Computing `a^b mod n` for huge `b` naively (multiply `a` by itself `b`
times) is infeasible for cryptographic-sized numbers. **Fast (modular)
exponentiation** ("square and multiply") does it in `O(log b)`
multiplications instead of `O(b)`.

```python
def fast_pow_mod(base, exp, mod):
    result = 1
    base = base % mod
    while exp > 0:
        if exp & 1:                      # if current bit is 1
            result = (result * base) % mod
        exp >>= 1                        # shift to next bit
        base = (base * base) % mod       # square the base
    return result

print(fast_pow_mod(7, 128, 13))          # matches pow(7, 128, 13)
print(pow(7, 128, 13))                    # Python's built-in does this internally
```

**Tracing "square and multiply" by hand** for `3^13 mod 7`. Write the
exponent `13` in binary (`1101`), then process it one bit at a time from
the *least significant* end, squaring the base every step and folding it
into the result only when the current bit is 1:

```
exp = 13 = 1101 (binary)     result starts at 1, base starts at 3

bit 1 (rightmost, value 1): result = 1×3 mod 7 = 3       base = 3² mod 7 = 2
bit 0            (value 0): result unchanged = 3          base = 2² mod 7 = 4
bit 1            (value 1): result = 3×4 mod 7 = 12 mod 7 = 5    base = 4² mod 7 = 2
bit 1 (leftmost, value 1):  result = 5×2 mod 7 = 10 mod 7 = 3    (last bit, base unused after)

Final result: 3
```

Only 4 squarings and 3 multiplications were needed for exponent 13 —
compare that to naively multiplying `3` by itself 13 times. For a
2048-bit RSA exponent, that's the difference between roughly 2,048
operations and `2^2048` operations — the entire reason this algorithm, not
the exponent's size, determines whether modern cryptography is
computationally feasible at all.

This is the actual algorithm running every time a server performs an RSA
or Diffie-Hellman operation — without it, cryptography using 2048+ bit
numbers would be computationally impossible.

### RSA, derived end to end

1. Pick two large primes `p`, `q`. Compute `n = p·q`.
2. Compute `φ(n) = (p-1)(q-1)` (Euler's totient, Ch. 19).
3. Pick public exponent `e` (commonly 65537) such that `gcd(e, φ(n)) = 1`.
4. Compute private exponent `d = e⁻¹ mod φ(n)` (modular inverse, Ch. 19).
5. **Public key:** `(n, e)`. **Private key:** `(n, d)`.
6. **Encrypt:** `ciphertext = message^e mod n`.
7. **Decrypt:** `message = ciphertext^d mod n`.

```python
# Toy RSA with small primes (NEVER do this with small primes in real life)
p, q = 61, 53
n = p * q                       # 3233
phi = (p - 1) * (q - 1)         # 3120
e = 17                          # public exponent, gcd(17, 3120) = 1
d = pow(e, -1, phi)             # private exponent (modular inverse)

message = 65
ciphertext = pow(message, e, n)
decrypted = pow(ciphertext, d, n)
print(ciphertext, decrypted)    # decrypted == 65 again
```
**Why this is secure with large primes:** anyone can see the public key
`(n, e)`, but recovering the private key `d` requires knowing `φ(n)`,
which requires factoring `n` back into `p` and `q` — and factoring a
2048-bit product of two large primes is (currently) computationally
infeasible. RSA's core hardness rests on this number-theoretic assumption
(Ch. 19), plus correct padding, parameter choices, and implementation.

### Diffie-Hellman key exchange — agreeing on a secret over a public channel

Two parties, Alice and Bob, agree on a public prime `p` and generator `g`.

```python
p = 23      # public prime (toy size — real DH uses 2048+ bit primes)
g = 5       # public generator

a = 6       # Alice's private key (secret, never transmitted)
b = 15      # Bob's private key (secret, never transmitted)

A = pow(g, a, p)   # Alice's public value, sent to Bob openly
B = pow(g, b, p)   # Bob's public value, sent to Alice openly

alice_shared_secret = pow(B, a, p)   # Bob's public value ^ Alice's private key
bob_shared_secret   = pow(A, b, p)   # Alice's public value ^ Bob's private key

print(alice_shared_secret == bob_shared_secret)   # True — both derive the same secret!
```
**Why an eavesdropper who sees `p`, `g`, `A`, `B` can't compute the shared
secret:** doing so requires solving the **discrete logarithm problem**
(recovering `a` from `g^a mod p`), which is believed computationally
infeasible for large primes — a different hardness assumption from RSA's
factoring problem, but the same flavor: an operation that's easy in one
direction (exponentiate) and (believed) hard to reverse (take the discrete
log).

### Elliptic curve cryptography (ECC) — same idea, smaller keys

ECC replaces modular exponentiation with **point addition on an elliptic
curve** (`y² = x³ + ax + b mod p`). "Multiplying" a point `G` by a scalar
`k` (i.e., adding `G` to itself `k` times using the curve's geometric
group law) is easy; recovering `k` from `k·G` and `G` (the **elliptic
curve discrete logarithm problem**) is believed *even harder* than the
classic discrete log problem relative to key size — which is why a
256-bit ECC key (`P-256`, `Curve25519`) offers security roughly comparable
to a 3072-bit RSA key, with much smaller keys and faster operations. This
is why modern TLS and SSH default to ECC (ECDHE, Ed25519) over classic
RSA/DH.

### The birthday attack — hash collision math

**The birthday paradox:** in a room of just 23 people, there's a >50%
chance two share a birthday — far fewer than the 366 you might intuitively
guess, because you're comparing *all pairs*, which grows quadratically.

```
P(collision) ≈ 1 - e^(-k²/2N)     for k items into N buckets
```
Applied to hash functions: to find *any* collision in an `n`-bit hash
function, an attacker needs only about `2^(n/2)` attempts (not `2^n`) —
this is why a hash function needs **double** the bit-length of its
intended collision-resistance security level. A 128-bit hash only offers
~64 bits of collision resistance, which is why SHA-256 (256-bit output,
~128-bit collision resistance) is the practical minimum for
collision-critical uses today, and *why* MD5 (128-bit, ~64-bit collision
resistance, further weakened by structural cryptanalytic attacks) is
considered broken.

```python
def birthday_bound(hash_bits):
    return 2 ** (hash_bits / 2)

print(f"{birthday_bound(128):.2e}")   # MD5-scale: ~1.8 * 10^19 attempts to find a collision
print(f"{birthday_bound(256):.2e}")   # SHA-256-scale: ~1.3 * 10^38 attempts
```

### Real-life engineering ties

- **Security:** this chapter is the mathematical grounding for evaluating
  cryptographic choices — key sizes, algorithm deprecation decisions
  (why MD5/SHA-1 are banned for security use), TLS cipher suite selection.
- **Network:** every TLS handshake performs (EC)DHE key exchange and RSA
  or ECDSA signature verification live, on every HTTPS connection —
  this math runs billions of times a second across the internet.
- **AI:** less direct, but cryptographic hashing underlies content
  deduplication in training data pipelines, and differential privacy
  (adding calibrated noise to protect individual data points) borrows
  concepts from this same hardness-assumption tradition.
- **SWE:** understanding why you should never write your own crypto —
  the security rests entirely on decades of cryptanalysis against
  *specific, standardized* constructions, not on the general idea being sound.

**More real-life examples:**
- In 2010, Sony's PlayStation 3 signing system was broken because its
  ECDSA implementation reused the same "random" nonce `k` for every
  signature instead of generating a fresh one each time — with two
  signatures sharing a nonce, simple algebra recovers the private key
  directly, a real demonstration of how fragile these assumptions are to
  implementation mistakes.
- A 2008 bug in Debian's OpenSSL package (CVE-2008-0166) accidentally
  removed almost all sources of entropy from key generation, meaning the
  "random" keys it generated came from a pool of only about 32,768
  possible values — every key generated on affected systems for nearly
  two years was trivially guessable, despite using "correct" RSA/DSA math.
- The 2017 ROCA vulnerability (CVE-2017-15361) found that a widely used
  Infineon smart-card/TPM chip generated RSA primes with a subtle
  structural weakness, making factoring dramatically faster — millions of
  real government ID cards and TPM-backed keys were affected.
- Modern TLS 1.3 mandates (EC)DHE for every handshake specifically to
  guarantee **forward secrecy** — even if a server's long-term private
  key is stolen later, past recorded traffic can't be decrypted, because
  each session's shared secret was never derived from that long-term key.

### Common mistake

Believing "hard to reverse" cryptographic assumptions are mathematically
*proven* — they are not (this echoes Ch. 19's point). RSA, DH, and ECC
security all rest on **unproven computational hardness assumptions** that
have simply resisted decades of attack. This is precisely why quantum
computing (Shor's algorithm) is such a big deal: it would break the
factoring and discrete-log assumptions efficiently, which is why NIST
finalized post-quantum standards (ML-KEM for key exchange, ML-DSA for
signatures) in 2024.

---

## 25. Network & Systems Mathematics

**Level:** Expert
**You'll use this for:** subnet design, capacity planning, latency budgets, and understanding TCP behavior under load.

### Subnetting / CIDR math

An IPv4 address is 32 bits. A `/24` CIDR block means the first 24 bits are
the fixed network portion, leaving `32 - 24 = 8` bits for host addresses:
`2^8 = 256` addresses (254 usable, since the first is the network address
and the last is the broadcast address).

```python
def subnet_info(cidr_bits):
    total_addresses = 2 ** (32 - cidr_bits)
    usable_hosts = total_addresses - 2   # minus network + broadcast
    return total_addresses, usable_hosts

print(subnet_info(24))   # (256, 254)
print(subnet_info(16))   # (65536, 65534)

# How many /24s fit inside a /16?
def subnets_within(outer_cidr, inner_cidr):
    return 2 ** (inner_cidr - outer_cidr)

print(subnets_within(16, 24))   # 256 — a /16 contains 256 /24 networks
```
Subnet masking itself is literal bitwise AND (Ch. 3): an IP is "in" a
subnet if `ip & mask == network_address & mask`.

```python
import ipaddress
net = ipaddress.ip_network("10.0.0.0/24")
print(net.num_addresses, list(net.hosts())[:3])
ip = ipaddress.ip_address("10.0.0.55")
print(ip in net)   # True
```

**Checking `10.0.0.55` against `10.0.0.0/24` bit by bit** — convert both
to binary, apply the mask (Ch. 1, 3), and compare:

```
IP:    10.0.0.55  =  00001010.00000000.00000000.00110111
Mask:  /24         =  11111111.11111111.11111111.00000000   <- 24 ones, then 8 zeros

IP AND Mask:          00001010.00000000.00000000.00000000  =  10.0.0.0
Network address:                                               10.0.0.0

They match  ->  10.0.0.55 IS inside the 10.0.0.0/24 network.
```

The AND mask zeroes out the last 8 bits (the host portion) on both sides —
if what's left over matches the network address, the IP belongs to that
subnet. The **broadcast address** is the same network bits with every
*host* bit forced to 1 instead of 0: `10.0.0.0` with the last 8 bits set
gives `10.0.0.255` — the top of the range, and the address the code's
`net.num_addresses` (256) and usable-host count (254) are built from.

### Bandwidth-delay product — why "just add more bandwidth" isn't the whole story

```
BDP = bandwidth × round-trip-time (RTT)
```
This tells you how many bytes can be "in flight" on the link at once — and
therefore the **minimum TCP window size** needed to fully utilize the
link. A 1 Gbps link with 100ms RTT has a BDP of `1,000,000,000 × 0.1 / 8 =
12.5 MB`. If your TCP window is smaller than this (the classic default
64KB window is 200× too small here), you cannot use the full bandwidth
**no matter how fast the link is** — this is the actual mathematical
reason "long fat networks" (high bandwidth, high latency, e.g.,
transcontinental links) need TCP window scaling enabled.

```python
def bandwidth_delay_product_bytes(bandwidth_bps, rtt_seconds):
    return (bandwidth_bps * rtt_seconds) / 8

bdp = bandwidth_delay_product_bytes(1_000_000_000, 0.1)
print(f"{bdp/1_000_000:.1f} MB needed in flight to saturate the link")
```

### Little's Law and queueing theory (M/M/1)

**Little's Law**, astonishingly general (applies to any stable queueing
system): `L = λ · W`
(average number of items in the system = arrival rate × average time each
item spends in the system). This applies to network packet queues, web
server request queues, and literally any waiting line.

The **M/M/1 queue** (Poisson arrivals, exponential service time, 1 server)
gives a closed-form for average wait time as a function of **utilization**
`ρ = λ/μ` (arrival rate ÷ service rate):

```
Average time in system: W = 1 / (μ - λ)     (only valid when λ < μ)
```

```python
def mm1_avg_wait(arrival_rate, service_rate):
    if arrival_rate >= service_rate:
        return float('inf')   # unstable queue: grows without bound
    return 1 / (service_rate - arrival_rate)

for util in [0.5, 0.7, 0.9, 0.95, 0.99]:
    mu = 100                       # server handles 100 req/sec
    lam = util * mu
    print(f"utilization={util:.0%}  avg wait={mm1_avg_wait(lam, mu)*1000:.1f}ms")
```
**This is the single most important, most underused piece of math in
capacity planning:** wait time doesn't grow linearly with utilization — it
grows **hyperbolically**, exploding as utilization approaches 100%. Going
from 90% to 95% utilization doesn't add 5% more latency, it can *double*
it. This is the rigorous justification for the common SRE rule of thumb
"never run production at >70-80% sustained utilization."

### TCP congestion control — AIMD (Additive Increase, Multiplicative Decrease)

TCP's congestion window `cwnd` grows **additively** (+1 MSS per RTT) when
things are going well, and shrinks **multiplicatively** (halved) on packet
loss:

```python
def aimd_simulation(rounds, loss_at_round):
    cwnd = 1
    history = []
    for r in range(rounds):
        if r == loss_at_round:
            cwnd = cwnd / 2            # multiplicative decrease
        else:
            cwnd += 1                  # additive increase
        history.append(cwnd)
    return history

print(aimd_simulation(20, loss_at_round=10))
```
This produces the classic loss-based TCP **"sawtooth"** throughput graph. The
asymmetry (slow linear growth, fast halving) is deliberately conservative
— a direct, real-world application of control theory (Ch. 29) tuned to
avoid congestion collapse across the shared, uncoordinated internet.

### Real-life engineering ties

- **Network:** subnetting/CIDR is daily-driver math for any network
  engineer; BDP explains "why is my transcontinental transfer slow despite
  huge bandwidth"; queueing theory is the rigorous basis for capacity
  planning and SLA design.
- **SWE/SRE:** Little's Law applies directly to any request queue, worker
  pool, or connection pool sizing decision.
- **Security:** queueing math explains and helps mitigate **SYN flood /
  slow-loris style DoS**: an attacker's goal is to push `ρ → 1` (or past
  it) deliberately, and understanding the hyperbolic wait-time curve
  explains why even modest additional load near saturation can cause
  disproportionate service degradation.

**More real-life examples:**
- Google's BBR congestion-control algorithm (deployed widely since 2016,
  including at YouTube) replaced classic loss-based AIMD/CUBIC with a
  model that directly estimates the bandwidth-delay product from this
  chapter, rather than waiting for packet loss as its congestion signal.
- CDN capacity planning (Netflix's Open Connect, Akamai, Cloudflare)
  applies Little's Law and queueing theory at enormous scale to decide how
  many edge servers and how much cache capacity a region needs to keep
  p99 latency low during peak load.
- Cloud auto-scaling policies (AWS Target Tracking, Kubernetes HPA) are
  configured with a target *utilization* — commonly 60-70%, not 100% — a
  direct, practiced application of the M/M/1 hyperbolic wait-time curve.
- Home routers and ISPs implementing **Smart Queue Management** (e.g.,
  `cake`, `fq_codel`) apply Active Queue Management theory to reduce
  "bufferbloat," trading a little throughput for a large drop in latency.

### Common mistake

Sizing infrastructure off **average** utilization instead of the queueing
math above — average utilization looks "fine" (e.g., 60%) while p99
latency is already terrible, because traffic isn't smooth; it's Poisson
(Ch. 14), so instantaneous utilization spikes well above the average
regularly.

---

## 26. Quantitative Security

**Level:** Expert
**You'll use this for:** turning "this feels risky" into a defensible, numeric argument.

### Password/key strength, connected end-to-end

Combining Ch. 15 (keyspace), Ch. 18 (brute-force time), and Ch. 23
(entropy) into one practical workflow:

```python
import math

def crack_time_estimate(entropy_bits, guesses_per_second):
    total_guesses = 2 ** entropy_bits
    avg_guesses = total_guesses / 2       # attacker expected to find it halfway through
    seconds = avg_guesses / guesses_per_second
    return seconds / (60*60*24*365)       # in years

# Offline attack against a fast, unsalted-hash-cracking rig (10^11 guesses/sec)
print(crack_time_estimate(40, 1e11), "years")   # weak: cracked in seconds
print(crack_time_estimate(80, 1e11), "years")   # strong: computationally infeasible

# Online attack, rate-limited to 10 guesses/sec (a login form with lockout)
print(crack_time_estimate(40, 10), "years")     # even "weak" 40-bit entropy holds up when rate-limited
```
**The takeaway that changes real security architecture decisions:** entropy
alone doesn't determine practical security — **entropy combined with
attacker guess rate** does. This is *why* rate-limiting a login endpoint is
often more impactful than forcing users into unmemorable passwords, and
why offline-crackable hash leaks (fast hash, no rate limit possible) are
categorically worse than online-only exposure.

### Detection system quality: precision, recall, ROC/AUC

Building directly on Bayes' theorem (Ch. 16):

```
Precision = TP / (TP + FP)     "of things I flagged, how many were real?"
Recall    = TP / (TP + FN)     "of real things, how many did I catch?"
```

**A tiny worked example, by hand.** Say a detector runs against 20 events:
8 real attacks, 12 benign. It flags 10 events, and of those 10, 6 turn out
to be real attacks (and it missed 2 real attacks entirely):

```
TP (flagged AND real)        = 6
FP (flagged but benign)      = 10 - 6 = 4
FN (real but NOT flagged)    = 8 - 6  = 2

Precision = TP / (TP+FP) = 6 / (6+4) = 6/10  = 0.60   "60% of alerts were real"
Recall    = TP / (TP+FN) = 6 / (6+2) = 6/8   = 0.75   "caught 75% of real attacks"
```

If the threshold were loosened to flag more events, more of those missed 2
attacks would likely get caught (recall rises) — but some of the 12
benign events would likely get swept in too (precision falls). That
push-pull, made concrete with real counts, is what the ROC curve below
plots across every possible threshold at once.

These trade off against each other — lowering your alert threshold
increases recall (catch more real attacks) but decreases precision (more
false alarms), and vice versa. An **ROC curve** plots this tradeoff across
every possible threshold; **AUC** (area under the curve) summarizes overall
detector quality independent of any one threshold choice.

```python
from sklearn.metrics import roc_curve, roc_auc_score

y_true = [0,0,0,0,1,0,0,1,0,1]                       # 1 = actual attack
y_scores = [0.1,0.4,0.2,0.3,0.8,0.05,0.6,0.9,0.15,0.7]  # detector's risk score
auc = roc_auc_score(y_true, y_scores)
fpr, tpr, thresholds = roc_curve(y_true, y_scores)
print("AUC:", auc)
```
Choosing the *operating threshold* on that curve is a business decision,
not a math one — but the math tells you exactly what tradeoff you're
accepting at each point.

### Risk = Likelihood × Impact (and why this needs distributions, not point estimates)

The classic security risk formula is often applied with single-number
guesses, which hides massive uncertainty. A more rigorous approach (used in
quantitative risk frameworks like FAIR) treats likelihood and impact as
**probability distributions** and runs a Monte Carlo simulation:

```python
import numpy as np
np.random.seed(0)

# Simulate 100,000 possible "years" of risk for a given threat scenario
n_sims = 100_000
# Likelihood: how many incidents per year (Poisson-distributed, Ch. 14)
incidents_per_year = np.random.poisson(lam=0.3, size=n_sims)
# Impact per incident: log-normally distributed (heavy-tailed cost, common for breach costs)
cost_per_incident = np.random.lognormal(mean=11, sigma=1.2, size=n_sims)  # ~$60k median

annual_loss = incidents_per_year * cost_per_incident
print("Median annual loss estimate: $", np.median(annual_loss))
print("95th percentile (tail risk): $", np.percentile(annual_loss, 95))
```
This Monte Carlo approach (repeatedly sampling from the input
distributions and tabulating outcomes) directly generalizes Ch. 14's
distributions and Ch. 9's percentile thinking into a decision-support
tool — instead of a single scary/reassuring number, you get a full
distribution of plausible outcomes, including the tail risk that a single
"expected value" calculation would hide.

### k-anonymity — quantifying anonymization strength

A dataset satisfies **k-anonymity** if every record is indistinguishable
from at least `k-1` other records on the "quasi-identifying" fields (zip
code, birth date, gender, etc.). This is a **combinatorics/counting**
question (Ch. 15): how many records share each combination of
quasi-identifiers?

```python
import pandas as pd
df = pd.DataFrame({
    "zip": ["02138","02138","02139","02138","02139"],
    "age": [29, 29, 45, 29, 45],
    "disease": ["flu","cold","flu","flu","cold"],
})
group_sizes = df.groupby(["zip","age"]).size()
k = group_sizes.min()
print("k-anonymity level:", k)   # smallest group size = the anonymity guarantee
```
A famous real-world result of this exact math: 87% of the US population is
uniquely identifiable from just `{zip code, birth date, gender}` alone
(Sweeney, 2000) — a striking demonstration that "anonymized" data with
insufficient k-anonymity often isn't anonymous at all.

### Real-life engineering ties

- **Security:** this entire chapter is core Sr.-security-engineer
  vocabulary — precision/recall/AUC for evaluating any detection system,
  entropy+rate-limiting reasoning for authentication design, Monte Carlo
  risk quantification for prioritizing a remediation backlog defensibly
  instead of by gut feeling.
- **AI:** the same precision/recall/ROC math applies to any binary
  classifier, security-related or not (fraud detection, content
  moderation, spam filtering).
- **SWE/Network:** capacity/incident risk modeling (Monte Carlo over
  failure-rate distributions) generalizes directly to reliability
  engineering (Ch. 14's exponential/Poisson distributions feeding in
  here).

**More real-life examples:**
- The "Have I Been Pwned" Pwned Passwords API (built by researcher Troy
  Hunt) lets you check whether a password appears in a breach corpus
  **without ever sending the full password or its full hash** — the
  client sends only the first 5 characters of a SHA-1 hash and receives
  back all matching suffixes, directly applying k-anonymity thinking in a
  widely deployed production tool.
- The FAIR (Factor Analysis of Information Risk) model is an
  industry-standard quantitative risk framework that formalizes "risk =
  probability distribution of loss frequency x probability distribution
  of loss magnitude" — a standardized, auditable version of the Monte
  Carlo simulation in this chapter, used for real enterprise risk
  reporting to boards and regulators.
- Credit scoring and insurance underwriting use the same precision/recall
  and cost-weighted tradeoff reasoning as fraud/intrusion detection: too
  permissive approves bad loans (false negatives), too strict rejects
  good customers (false positives).
- Bug bounty programs and CVSS (Common Vulnerability Scoring System)
  scores are attempts to formalize "likelihood x impact" into a
  repeatable number for prioritizing a remediation backlog.

### Common mistake

Optimizing a detection system purely for accuracy or purely for recall in
isolation. In a Bayes-driven rare-event context (Ch. 16), a 100%-recall
detector that flags everything is worthless (zero precision); the right
metric is almost always a threshold chosen on a precision/recall or
cost-weighted curve, not a single global metric.

---

## 27. Signal Processing Basics

**Level:** Expert
**You'll use this for:** understanding RF/wireless engineering, audio/video ML, and reading network jitter correctly.

### The Fourier transform — decomposing a signal into frequencies

Any signal (audio waveform, radio signal, even a time series of
request latencies) can be decomposed into a sum of sine waves of different
frequencies. The Fourier transform converts a signal from the **time
domain** (value vs. time) to the **frequency domain** (how much of each
frequency is present).

```python
import numpy as np

# A signal made of a 5Hz wave and a 50Hz wave mixed together
t = np.linspace(0, 1, 1000, endpoint=False)
signal = np.sin(2*np.pi*5*t) + 0.5*np.sin(2*np.pi*50*t)

fft_result = np.fft.fft(signal)
freqs = np.fft.fftfreq(len(t), d=t[1]-t[0])
magnitude = np.abs(fft_result)

# The two dominant frequencies show up as spikes at 5Hz and 50Hz
dominant_freqs = freqs[np.argsort(magnitude)[-4:]]
print(sorted(set(abs(f) for f in dominant_freqs if f > 0)))  # ~[5.0, 50.0]
```
This exact operation (FFT — Fast Fourier Transform, an `O(n log n)`
algorithm for computing this decomposition) runs inside Wi-Fi/cellular
modems (OFDM modulation splits data across many simultaneous frequency
sub-carriers), audio codecs, and image compression (JPEG uses a related
transform, the DCT).

### The Nyquist-Shannon sampling theorem

To perfectly reconstruct a signal containing frequencies up to `f_max`, you
must sample at a rate of at least `2 × f_max` (the **Nyquist rate**).
Sampling below this causes **aliasing** — high frequencies masquerade as
false low frequencies in your sampled data (the classic "wagon wheel
spinning backward" effect in old films is aliasing).

```
sample_rate >= 2 * highest_frequency_of_interest
```

**Computing an alias frequency by hand.** If you sample a 6Hz signal at
only 5 samples/second (below the required `2×6=12Hz` Nyquist rate), the
signal doesn't just look choppy — it appears as a *different, false*
frequency:

```
alias_frequency = | true_frequency - nearest_multiple_of(sample_rate) |
                 = | 6 - 5 |
                 = 1 Hz
```

Sampling a 6Hz wave at 5Hz makes it look exactly like a slow 1Hz wave in
the recorded data — indistinguishable from the real thing once sampled,
which is *why* aliasing is dangerous: nothing in the sampled data itself
flags that anything went wrong.

This is why audio CDs sample at 44.1kHz (comfortably above 2×20kHz, the
upper edge of human hearing), and it applies directly to *any* discrete
sampling of a continuous signal — including monitoring systems: if you
sample a metric (e.g., CPU usage) once a minute, you *cannot* detect
spikes/oscillations happening faster than once every 2 minutes — they'll
alias into misleading patterns in your dashboard.

### Real-life engineering ties

- **Network:** RF/wireless engineering (channel width, OFDM sub-carrier
  spacing) is direct Fourier-domain reasoning; understanding *why*
  wider Wi-Fi channels (Ch. 23's Shannon-Hartley) carry more data.
- **AI:** audio/speech models operate on spectrograms (a Fourier-domain
  representation of sound); positional encodings in transformers are
  literally built from sine/cosine functions of different frequencies —
  Fourier-style reasoning by design.
- **Security:** side-channel analysis (e.g., analyzing power-consumption
  or electromagnetic-emission traces to extract a cryptographic key) is
  fundamentally frequency-domain signal analysis; timing-attack detection
  also benefits from understanding sampling-rate limitations.
- **SWE/SRE:** understanding aliasing explains why your 1-minute-resolution
  monitoring can completely miss (or badly misrepresent) sub-minute
  latency spikes — a direct, practical consequence of the sampling theorem.

**More real-life examples:**
- Shazam's music-identification algorithm converts audio into a
  spectrogram (a Fourier-domain representation) and fingerprints the
  pattern of frequency peaks over time — matching a few seconds of noisy
  audio against a database of millions of songs.
- Wi-Fi (802.11a/g/n/ac/ax) and 4G/5G cellular both use OFDM
  (Orthogonal Frequency-Division Multiplexing), which splits a channel
  into dozens or hundreds of narrow sub-carriers computed via FFT — a
  massive-scale production use of the Fourier transform running in every
  phone and router.
- Speech-to-text systems (including OpenAI's Whisper) preprocess raw
  audio into a mel-spectrogram — a Fourier-transform-based, perceptually
  weighted frequency representation — before feeding it into the network.
- MRI machines acquire raw data directly in the frequency domain (called
  "k-space" in radiology) and apply an inverse Fourier transform to
  reconstruct the anatomical image radiologists actually view.

### Common mistake

Assuming higher monitoring-sample-rate is "just more data, no downside" —
it's the opposite problem people usually hit: under-sampling and believing
the resulting (aliased) chart accurately reflects reality.

---

## 28. Game Theory & Adversarial Thinking

**Level:** Expert
**You'll use this for:** reasoning formally about attacker/defender dynamics, adversarial ML, and any system with competing incentives.

### The idea

Game theory studies situations where multiple rational actors' outcomes
depend on each other's choices. A **security system is inherently a
game**: the defender chooses controls, the attacker chooses a strategy in
response, and both are trying to optimize against a thinking adversary —
not against random chance (which is why pure probability/statistics,
Ch. 8–9, is necessary but not sufficient for security reasoning).

### Nash equilibrium

A **Nash equilibrium** is a set of strategies where no player can improve
their outcome by unilaterally changing their own strategy, given what
everyone else is doing. Classic example — the **Prisoner's Dilemma**:

|  | Bob cooperates | Bob defects |
|---|---|---|
| **Alice cooperates** | (-1, -1) | (-3, 0) |
| **Alice defects** | (0, -3) | (-2, -2) |

**Verifying it's an equilibrium, by checking every possible unilateral
switch.** A cell is a Nash equilibrium only if *neither* player can do
better by changing just their own move, holding the other's move fixed:

```
Start at (defect, defect) = payoffs (-2, -2).

Would Alice switch to cooperate, if Bob keeps defecting?
  Alice's payoff would become -3 (worse than -2)  -> no, she won't switch.

Would Bob switch to cooperate, if Alice keeps defecting?
  Bob's payoff would become -3 (worse than -2)    -> no, he won't switch.

Neither player benefits from switching alone -> (defect, defect) IS a
Nash equilibrium.

Now check (cooperate, cooperate) = payoffs (-1, -1), the jointly better outcome:
  Would Alice switch to defect, if Bob keeps cooperating?
    Alice's payoff would become 0 (better than -1) -> YES, she'd switch.

Since at least one player benefits from switching, (cooperate, cooperate)
is NOT a Nash equilibrium — it's unstable, even though it's better for both.
```

Both defecting is the Nash equilibrium (neither benefits from unilaterally
switching), even though both cooperating would be better *for both of
them jointly* — the mathematical core of why purely self-interested
security investment decisions across organizations (e.g., ISPs not
filtering spoofed traffic because it costs them money but benefits
everyone else) tend to under-invest relative to the collective optimum.
This is literally why some security problems (BGP route hijacking, DDoS
amplification via spoofing) persist for decades despite known technical
fixes — it's a game-theory/incentive problem, not a technical one.

### Zero-sum vs. non-zero-sum framing of security

Classic security thinking often (wrongly) treats attacker/defender as
purely zero-sum (attacker's gain = defender's loss). In reality, many
security investments are **non-zero-sum**: raising the attacker's cost
(e.g., MFA) can deter an attacker entirely (redirecting them to an easier
target) without any actual "loss" transferred — both the specific
defender and the broader ecosystem can be better off.

### Adversarial ML as a minimax game

Adversarial example generation is literally: find the smallest input
perturbation `δ` that maximizes the model's loss:
```
δ* = argmax_δ  Loss(model(x + δ), true_label)   subject to |δ| < ε
```
This is a direct callback to gradient ascent from Ch. 21 (instead of
descending the loss w.r.t. weights, you ascend it w.r.t. the input),
framed as an adversary optimizing against your model. **Adversarial
training** (a defense) turns this into a full minimax game: the defender
trains to minimize loss *even under the worst-case attacker perturbation*:

```
min_θ  max_δ  Loss(model_θ(x + δ), true_label)
```
GAN (Generative Adversarial Network) training is the same minimax
structure applied generatively — a generator and discriminator playing
an adversarial game against each other until they reach (ideally) a Nash
equilibrium.

### Real-life engineering ties

- **Security:** threat modeling is applied game theory — reasoning about a
  *rational, adaptive* adversary is fundamentally different from reasoning
  about random failures (Ch. 14), and mixing up the two framings is a
  common root cause of security controls that look good on paper but fail
  against a motivated attacker who simply routes around them.
- **AI:** adversarial robustness, GAN training, multi-agent reinforcement
  learning, and RLHF's reward-model/policy dynamic are all game-theoretic.
- **Network/SWE:** auction-based systems (real-time bidding for ads, cloud
  spot-instance pricing) are direct game-theory applications; incentive
  design for decentralized systems (peer-to-peer networks, blockchain
  consensus mechanisms) is applied mechanism design (game theory in
  reverse: design the game so the *equilibrium* is the outcome you want).

**More real-life examples:**
- Google Ads and most real-time ad exchanges use second-price (Vickrey)
  auction mechanisms, chosen specifically because they make "bid your
  true value" the dominant strategy for advertisers, unlike a naive
  first-price auction.
- Bitcoin's proof-of-work design is deliberate mechanism design: a "51%
  attack" is technically possible but game-theoretically irrational for a
  large miner, because acquiring enough hash power to attack the network
  is more expensive than the profit from honestly mining.
- Coordinated vulnerability disclosure and bug bounty programs shift the
  game-theoretic payoff for a researcher who finds a bug: paying a
  bounty makes responsible disclosure competitive with selling the
  exploit on a gray/black market, changing the equilibrium outcome.
- **Braess's paradox** — where adding a new road (or network link) can
  *increase* overall congestion because individually rational routing no
  longer aligns with the system optimum — has been documented in real
  city traffic studies (Seoul reported reduced congestion after removing
  a highway in the 2000s).

### Common mistake

Modeling an attacker as a fixed, non-adaptive probability distribution
(e.g., "attacks happen with 1% probability per day" — pure Ch. 14
thinking) instead of a rational actor who will adapt strategy in response
to your defenses. This is *why* purely statistical anomaly detection
(Ch. 26) is necessary but insufficient — a game-theoretically savvy
attacker learns your detection thresholds and deliberately stays under
them.

---

## 29. Control Theory Basics (Feedback Loops, PID)

**Level:** Expert
**You'll use this for:** understanding autoscalers, rate limiters, and TCP congestion control as the same underlying pattern.

### The idea

A **control system** continuously measures an output, compares it to a
target ("setpoint"), computes an **error** (`error = setpoint - measured`),
and adjusts an input to drive that error toward zero. Your home
thermostat, a cloud autoscaler, and TCP congestion control are all the same
mathematical pattern.

### PID controller — the workhorse of feedback control

```
output = Kp·error + Ki·∫error dt + Kd·(d(error)/dt)
```
- **P (Proportional):** react proportionally to the *current* error — bigger
  error, bigger correction.
- **I (Integral):** react to *accumulated* past error — corrects persistent
  small biases (steady-state error) that pure proportional control leaves
  behind.
- **D (Derivative):** react to the *rate of change* of error — dampens
  oscillation/overshoot by anticipating where the error is heading (this
  literally reuses the derivative from Ch. 12).

```python
class PIDController:
    def __init__(self, kp, ki, kd, setpoint):
        self.kp, self.ki, self.kd = kp, ki, kd
        self.setpoint = setpoint
        self.integral = 0
        self.prev_error = 0

    def step(self, measured_value, dt=1.0):
        error = self.setpoint - measured_value
        self.integral += error * dt
        derivative = (error - self.prev_error) / dt
        self.prev_error = error
        return self.kp*error + self.ki*self.integral + self.kd*derivative

# Autoscaler: target 60% CPU utilization, adjust replica count
pid = PIDController(kp=0.5, ki=0.1, kd=0.05, setpoint=60)
current_cpu = 85
adjustment = pid.step(current_cpu)
print(adjustment)   # negative -> scale up (reduce per-replica load)
```

**Tracing two PID steps by hand**, `kp=0.5, ki=0.1, kd=0.05, setpoint=60`,
starting with `integral=0, prev_error=0`:

```
Step 1 — measured CPU = 85:
  error      = setpoint - measured = 60 - 85 = -25
  integral  += error × dt = 0 + (-25×1) = -25
  derivative = (error - prev_error) / dt = (-25 - 0) / 1 = -25
  output = kp·error + ki·integral + kd·derivative
         = 0.5×(-25) + 0.1×(-25) + 0.05×(-25)
         = -12.5   + -2.5        + -1.25
         = -16.25                                <- negative: scale up

Step 2 — CPU eased to 70 (prev_error is now -25):
  error      = 60 - 70 = -10
  integral  += -10  ->  integral = -25 + -10 = -35
  derivative = (-10 - (-25)) / 1 = 15              <- error is shrinking, so this is positive (damping)
  output = 0.5×(-10) + 0.1×(-35) + 0.05×15
         = -5        + -3.5      + 0.75
         = -7.75                                  <- smaller correction: still scaling up, but less aggressively
```

Notice the **integral** term keeps accumulating (`-25` then `-35`) even as
the error shrinks — that's the "remembers persistent past error" behavior
described below — while the **derivative** term flipped positive the
moment the error started improving, pulling the output back toward zero
before the P-term alone would.

### Why pure proportional control isn't enough (and why over-tuning causes oscillation)

- **P-only control** typically leaves a persistent **steady-state error**
  (it never quite reaches the setpoint, because the correction shrinks to
  zero as the error shrinks to zero — a vanishing feedback signal, similar
  in spirit to Ch. 21's vanishing gradient problem).
- **Too-aggressive gains** (`Kp` too high) cause **overshoot and
  oscillation** — the exact failure mode of an autoscaler that
  thrashes (scale up, overshoot, scale down, undershoot, repeat) instead
  of converging smoothly. This is mathematically the same instability
  pattern as a too-high learning rate in gradient descent (Ch. 21) —
  both are discrete-time feedback systems that can overshoot and diverge
  if the step size relative to the system's sensitivity is too large.

### TCP congestion control as control theory

Revisiting Ch. 25's AIMD from a control-theory lens: TCP's congestion
window is a **feedback-controlled variable**, where packet loss is the
"error signal" telling the sender it exceeded the network's capacity. The
conservative asymmetry (slow additive increase, fast multiplicative
decrease) is a deliberately *stable* control law, chosen specifically so
that thousands of independent, uncoordinated TCP senders across the global
internet converge toward fair bandwidth sharing without a central
controller — a remarkable real-world distributed-control-theory result.

### Real-life engineering ties

- **SWE/SRE/Network:** autoscalers, rate limiters, adaptive timeout/retry
  systems, and load-shedding controllers are all PID-flavored (even when
  implemented as simpler heuristics, they're approximating the same
  feedback-control idea, often badly, which is why they sometimes
  oscillate in production).
- **Network:** TCP congestion control (AIMD, Ch. 25), Wi-Fi's
  CSMA/CA backoff, and Active Queue Management (CoDel, RED) are all
  feedback control systems reacting to a congestion signal.
- **AI:** reinforcement learning control policies for physical/robotic
  systems build directly on control theory; some learning-rate schedulers
  are explicitly PID-inspired ("adjust the learning rate based on the
  rate of change of the loss").
- **Security:** adaptive rate-limiting against brute-force/DDoS attempts
  is a feedback controller reacting to a suspicious-traffic "error
  signal" — poorly tuned gains here cause exactly the oscillation problem
  above (over-aggressive throttling that also blocks legitimate users in
  bursts, then relaxes too far, repeating).

**More real-life examples:**
- Bitcoin's mining difficulty adjusts automatically roughly every two
  weeks specifically to keep the average block time near 10 minutes
  regardless of how much total mining power joins or leaves the
  network — a real, deployed feedback controller operating at global
  scale with no central operator.
- Cruise control (and adaptive cruise control / lane-keeping in modern
  cars) is a textbook PID controller: it measures actual speed (or lane
  position), compares it to the target, and adjusts throttle (or
  steering) — the same three-term structure as this chapter's code example.
- Data center cooling systems use classic industrial PID controllers to
  hold server-room temperature at a setpoint despite constantly varying
  heat load from compute racks — the literal thermostat analogy from this
  chapter, at industrial scale, and a meaningful lever in a data center's
  energy efficiency (PUE).
- Kubernetes' Horizontal Pod Autoscaler and most cloud autoscalers are, in
  effect, simplified proportional (P-only) controllers — exactly why
  they're prone to the oscillation ("flapping") failure mode described in
  this chapter, and why production setups add cooldown windows as a crude
  substitute for proper derivative damping.

### Common mistake

Tuning a feedback system (autoscaler, rate limiter) by trial and error in
production without understanding *why* it oscillates — the oscillation is
a textbook symptom of gains that are too aggressive relative to the
system's response lag, and the fix (add damping via a derivative term, or
simply reduce the proportional gain and increase the reaction window) is
well understood control theory, not guesswork.

---

# Part 4 — Putting It Together

Nothing new to learn from here on — this Part is the reference layer: a
formula cheat sheet, hands-on practice labs ordered by difficulty, a
role-based further-reading map back into the 29 chapters, and an honest
table of what this guide deliberately left out and where it actually lives.

## Formula cheat sheet

```
NUMBER SYSTEMS
  bits needed for n values         = ceil(log2(n))
  n-bit signed integer range       = [-2^(n-1), 2^(n-1) - 1]

PROBABILITY & STATISTICS
  P(A ∪ B)                          = P(A) + P(B) - P(A ∩ B)
  P(A|B)  [Bayes]                   = P(B|A)·P(A) / P(B)
  variance                          = (1/n)·Σ(xi - mean)²
  std dev                           = sqrt(variance)
  binomial P(k successes of n)      = C(n,k)·p^k·(1-p)^(n-k)
  poisson P(k events, rate λ)       = (λ^k · e^-λ) / k!

COMBINATORICS
  permutations                      = n! / (n-k)!
  combinations                      = n! / (k!(n-k)!)
  keyspace size                     = alphabet_size ^ length

INFORMATION THEORY
  entropy (bits)                    = -Σ p(x)·log2(p(x))
  password entropy                  = log2(keyspace_size) = length · log2(alphabet_size)
  KL divergence                     = Σ p(x)·log2(p(x)/q(x))
  Shannon channel capacity          = B · log2(1 + S/N)

CRYPTOGRAPHY
  RSA: n = p·q, φ(n) = (p-1)(q-1), d = e^-1 mod φ(n)
  Diffie-Hellman shared secret      = g^(a·b) mod p   (computed as (g^b)^a = (g^a)^b)
  birthday-bound collision attempts ≈ 2^(hash_bits / 2)
  brute-force expected time         = 2^(entropy_bits) / (2 · guesses_per_second)

NETWORKING
  subnet addresses                  = 2^(32 - cidr_bits)
  bandwidth-delay product           = bandwidth × RTT
  Little's Law                      = L (items in system) = λ (arrival rate) · W (time in system)
  M/M/1 average wait                = 1 / (service_rate - arrival_rate)

PROBABILITY DEEPER (Ch. 8, 14)
  conditional probability           P(A|B) = P(A ∩ B) / P(B)
  total probability                 P(A) = Σ P(A|Bᵢ)·P(Bᵢ)
  expectation                       E[X] = Σ x·P(X=x)   or   ∫ x·f(x) dx
  linearity (always holds)          E[X + Y] = E[X] + E[Y]
  variance shortcut                 Var(X) = E[X²] - (E[X])²
  normal density                    f(x) = 1/(σ√(2π)) · e^(-(x-μ)²/(2σ²))
  standard error (CLT)              σ / √n
  Hoeffding sample size             n ≥ ln(2/δ) / (2ε²)
  zero failures in n trials         95% upper bound ≈ 3/n
  Beta update                       Beta(α, β) + k of n  ->  Beta(α+k, β+n-k)

STATISTICS DEEPER (Ch. 9, 22)
  z-score                           z = (x - μ) / σ
  correlation                       r = Cov(X,Y) / (s_X · s_Y)
  EWMA                              S = (1-α)·S + α·x
  TCP RTO (RFC 6298)                RTO = SRTT + 4·RTTVAR
  CI for a proportion               p̂ ± 1.96·√(p̂(1-p̂)/n)
  A/B sample size per group         (z_(α/2) + z_β)² · [p₁(1-p₁) + p₂(1-p₂)] / δ²

LINEAR ALGEBRA (Ch. 10, 11, 20)
  dot product, geometric            A·B = ‖A‖‖B‖cos θ
  dense layer                       h = ReLU(W·x + b)
  least squares                     w = (XᵀX)⁻¹ Xᵀ y
  eigen-equation                    A·v = λ·v,  solve det(A - λI) = 0
  matrix powers                     Aᵏ = P·Dᵏ·P⁻¹
  PageRank                          PR = (1-d)/N + d·M·PR
  Markov steady state               π·P = π
  graph Laplacian                   L = D - A   (λ₂ = algebraic connectivity)
  SVD                               A = U·Σ·Vᵀ,  ‖A - Aₖ‖²_F = Σ_{i>k} σᵢ²
  attention                         softmax(Q·Kᵀ / √d_k)·V

CALCULUS / ML
  derivative                        f'(x) = lim(h→0) [f(x+h) - f(x)] / h
  sigmoid derivative                σ'(x) = σ(x)·(1 - σ(x))
  Taylor (2nd order)                f(x+Δ) ≈ f(x) + ∇fᵀΔ + ½ΔᵀHΔ
  Newton's method                   x_{n+1} = x_n - f(x_n)/f'(x_n)
  Fundamental Theorem               ∫ₐᵇ f(x) dx = F(b) - F(a)
  Gaussian integral                 ∫ e^(-x²) dx = √π
  gradient descent update           = w_new = w_old - learning_rate · ∇Loss(w_old)
  chain rule                        = dy/dx = f'(g(x)) · g'(x)
  linear layer backprop             ∂L/∂W = δ·xᵀ,   ∂L/∂x = Wᵀ·δ
  softmax + cross-entropy gradient  ∂L/∂z = p - y
  stable learning rate              η < 2 / λ_max(Hessian)
  Adam step                         w -= η · m̂ / (√v̂ + ε)
  Mathis TCP throughput             ≈ (MSS/RTT) · 1.22/√p

CONTROL THEORY
  PID output = Kp·error + Ki·∫error·dt + Kd·(d(error)/dt)
```

## Practice labs (hands-on, in order of difficulty)

1. **Overflow bug hunt:** write an 8-bit signed-integer simulator (Ch. 1)
   and find the exact input that wraps from 127 to -128.
2. **Percentile dashboard:** given a CSV of latency samples, compute and
   plot p50/p95/p99 vs. the mean, and explain the gap (Ch. 9).
3. **Toy PageRank:** build a 6-node link graph and compute the dominant
   eigenvector by hand and with `numpy.linalg.eig` (Ch. 11).
4. **Backprop by hand:** extend the 2-layer worked example in Ch. 21 to 3
   layers and verify your hand-derivative against PyTorch's `autograd`.
5. **Base-rate simulator:** build the IDS false-positive simulation from
   Ch. 16 and plot precision as the base rate varies from 1-in-10 to
   1-in-1,000,000.
6. **Toy RSA:** implement full key generation, encryption, and decryption
   with 200-bit primes (using a primality-testing library), and time how
   long factoring `n` takes vs. key size.
7. **M/M/1 capacity planner:** simulate a request queue at increasing
   utilization and reproduce the hyperbolic wait-time curve from Ch. 25.
8. **PID autoscaler:** simulate a service under bursty load and tune a PID
   controller to scale replica count without oscillating (Ch. 29).

## Further Reading

Organized by topic cluster, roughly matching Parts 1–3. Free/legally
available resources are marked **(free)**.

### Foundations (Ch. 1–9: number systems, logic, probability, statistics)

**Books**
- *Mathematics for Computer Science* — Lehman, Leighton, Meyer (MIT
  6.042). **(free)** The single best broad foundation text for this whole
  part; covers logic, number theory, graphs, probability in one book.
- *How to Bake Pi* — Eugenia Cheng. Builds real mathematical intuition
  from everyday analogies; great for Ch. 3–5.
- *Naked Statistics* — Charles Wheelan. The most readable plain-English
  introduction to Ch. 8–9's ideas, no formulas required to get the intuition.
- *The Art of Computer Programming, Vol. 1* — Donald Knuth. Denser, but
  the canonical reference for Ch. 1's bit-level fundamentals.

**Blogs / articles**
- BetterExplained (betterexplained.com) — intuitive, visual explanations
  of exponents, logs, and Euler's number; consistently the best "aha"
  source for Ch. 5–6.
- Math Is Fun (mathsisfun.com) — quick, clear reference for any formula
  in Ch. 1–9 when you need a refresher, not a lecture.

**Video**
- 3Blue1Brown (YouTube) — "Essence of..." series; short, visual,
  intuition-first. Start here for anything that feels abstract.
- Khan Academy — full, free, structured courses covering Ch. 1–9 in order
  with exercises.
- Professor Leonard (YouTube) — long-form, thorough lecture-style algebra
  and calculus, good if you want the full derivation, not just intuition.

### Core engineering math (Ch. 10–19: linear algebra, calculus, discrete math)

**Books**
- *Linear Algebra Done Right* — Sheldon Axler. Rigorous, conceptual
  treatment of Ch. 10–11 for readers who want the "why," not just the
  mechanics.
- *Introduction to Algorithms* (CLRS) — Cormen, Leiserson, Rivest, Stein.
  The standard reference for Ch. 17–18 (graphs, complexity).
- *Probability and Statistics for Computer Science* — David Forsyth.
  Written specifically for engineers, ties directly into Ch. 14–16.
- *Elementary Number Theory* — David Burton. Approachable, thorough
  treatment of Ch. 19, written as preparation for cryptography.

**Blogs / articles**
- Distill.pub (archived but still excellent) — beautifully illustrated,
  rigorous explanations bridging Ch. 10–13 and ML.
- Terence Tao's blog (terrytao.wordpress.com) — advanced but exceptionally
  clear; useful once the basics from Part 2 are solid.

**Video**
- MIT OCW 18.06 Linear Algebra — Gilbert Strang's full course, free,
  the standard reference lecture series for Ch. 10–11.
- 3Blue1Brown — "Essence of Linear Algebra" and "Essence of Calculus"
  series, directly aligned with Ch. 10–13.
- StatQuest with Josh Starmer (YouTube) — short, clear videos on Ch. 14–16
  (distributions, Bayes) with a friendly, repetition-based teaching style.

### AI/ML mathematics (Ch. 20–23)

**Books**
- *Mathematics for Machine Learning* — Deisenroth, Faisal, Ong. **(free
  PDF)** Purpose-built bridge from Part 2 math straight into ML; the most
  direct next step after this guide.
- *Deep Learning* — Goodfellow, Bengio, Courville. **(free online)** The
  standard reference for Ch. 20–21 in full mathematical depth.
- *Pattern Recognition and Machine Learning* — Christopher Bishop. Denser,
  more classically statistical treatment of Ch. 22.
- *Information Theory, Inference, and Learning Algorithms* — David MacKay.
  **(free)** The best single book connecting Ch. 23 to ML and coding theory.

**Blogs / articles**
- colah's blog (colah.github.io) — Christopher Olah's visual explanations
  of backpropagation and neural network internals; the best intuition
  bridge for Ch. 21.
- Jay Alammar's blog (jalammar.github.io) — "The Illustrated Transformer"
  and related posts; excellent for seeing Ch. 10's matrix math in a real
  production architecture.
- Sebastian Raschka's blog and newsletter — practical, math-grounded
  deep-learning explanations aimed at working engineers.

**Video**
- Andrej Karpathy — "Neural Networks: Zero to Hero" (YouTube series).
  Builds backpropagation and a GPT from scratch, in code, matching this
  guide's "derive it, don't memorize it" approach to Ch. 21.
- 3Blue1Brown — "Neural Networks" series, the visual companion to Ch. 20–21.
- StatQuest — ML-specific playlist covering gradient descent, cross-entropy,
  PCA, and regularization at an approachable pace.

### Cryptography & security mathematics (Ch. 16, 19, 23, 24, 26, 28)

**Books**
- *Serious Cryptography* — Jean-Philippe Aumasson. The best modern,
  practitioner-focused treatment of Ch. 24.
- *Handbook of Applied Cryptography* — Menezes, van Oorschot, Vanstone.
  **(free PDF, author-hosted)** The classic deep reference for the number
  theory behind Ch. 19 and 24.
- *Cryptography Engineering* — Ferguson, Schneier, Kohno. Practical,
  implementation-aware follow-up once the math in Ch. 24 is solid.
- *Thinking Strategically* — Avinash Dixit & Barry Nalebuff. The most
  accessible general introduction to Ch. 28's game theory.

**Blogs / articles**
- A Few Thoughts on Cryptographic Engineering — Matthew Green
  (blog.cryptographyengineering.com). Deep, readable dives into real
  cryptographic failures and design decisions.
- Troy Hunt's blog (troyhunt.com) — practical, real-incident-driven
  writing on password entropy, breaches, and the k-anonymity design used
  in Have I Been Pwned (Ch. 26).

**Video**
- Computerphile (YouTube) — short, clear explainers on RSA,
  Diffie-Hellman, and hashing, well matched to Ch. 19 and 24.
- Khan Academy's cryptography series — a free, structured path through
  the same modular-arithmetic-to-RSA progression as Ch. 19/24.

### Networking & systems mathematics (Ch. 17, 25, 27, 29)

**Books**
- *Computer Networking: A Top-Down Approach* — Kurose & Ross. The
  standard reference for Ch. 25 (congestion control, queueing) and Ch. 27.
- *The Art of Computer Systems Performance Analysis* — Raj Jain. The
  deep, rigorous reference for queueing theory and capacity planning
  math in Ch. 25.
- *TCP/IP Illustrated, Vol. 1* — W. Richard Stevens. Ties Ch. 25's AIMD
  math directly to real packet traces.
- *Feedback Control for Computer Systems* — Philipp K. Janert. A rare,
  excellent book applying Ch. 29's control theory specifically to software
  systems (autoscalers, rate limiters) rather than mechanical/electrical ones.

**Blogs / articles**
- High Scalability (highscalability.com) — real-world architecture
  write-ups that repeatedly touch Ch. 25's queueing and capacity math.
- The Cloudflare Blog (blog.cloudflare.com) — frequent, math-heavy
  postmortems and deep dives directly illustrating Ch. 18, 25, and 27
  concepts in production incidents.

**Video**
- MIT 6.033 (Computer System Engineering) lecture recordings — covers
  distributed systems and networking fundamentals underlying Ch. 17/25.
- Practical Networking (YouTube) — clear walkthroughs of subnetting/CIDR
  math (Ch. 25) for hands-on practice.

## Topics Not Covered (and where they'd fit)

This guide is deliberately scoped to the math that shows up *directly* and
*repeatedly* in day-to-day software/AI/network/security engineering work.
The following are real, legitimate branches of mathematics that a senior
engineer may eventually need, depending on specialization — they're left
out here because they're either less universally applicable, or deep
enough to warrant their own dedicated guide. Listed with why they matter
and who tends to need them:

| Topic | Why it matters | Who needs it most |
|---|---|---|
| **Differential equations (ODEs/PDEs)** | Simulation, epidemiological/worm-spread modeling beyond simple exponential growth, physics-based rendering | AI (simulation, physics-informed ML), robotics |
| **Complex numbers & complex analysis** | Proper treatment of the Fourier transform, AC circuit analysis, some quantum computing math | Signal processing, DSP/RF engineers |
| **Abstract algebra (groups, rings, fields)** | Deeper cryptography (lattice-based post-quantum schemes), error-correcting code theory | Cryptography engineers, post-quantum specialists |
| **Coding theory (Hamming, Reed-Solomon codes)** | RAID/storage redundancy, QR codes, forward error correction in networks and storage media | Storage/systems engineers, network engineers |
| **Numerical methods / numerical analysis** | Floating-point error accumulation in long computations, iterative solvers, root-finding at scale | ML infra/training-stability engineers, scientific computing |
| **Convex optimization (duality, KKT, interior-point methods)** | Deeper SVM theory, resource-allocation and network-flow optimization at scale | ML researchers, operations research |
| **Formal logic & proof systems (predicate logic, SAT/SMT)** | Formal verification, automated theorem proving, security protocol verification | Security engineers doing formal verification, compiler engineers |
| **Automata theory & formal languages** | Regex engine internals (and why ReDoS happens, Ch. 18), parser/compiler design | Compiler engineers, security engineers auditing parsers |
| **Measure theory & rigorous probability** | Formal foundations under Ch. 8–9/14/22, needed for continuous-time stochastic processes | ML researchers (theory-heavy roles), quants |
| **Stochastic calculus / Brownian motion** | Diffusion/score-based generative model theory, quantitative finance | AI researchers (generative models), fintech quants |
| **Time-series analysis (ARIMA, seasonality, wavelets)** | Forecasting, capacity planning over time, anomaly detection with seasonality | SRE/capacity planning, AI (forecasting models) |
| **Differential privacy math** | Formal privacy guarantees when training on or querying sensitive data | AI/ML engineers handling regulated data, privacy engineers |
| **Zero-knowledge proofs, homomorphic encryption, secure multi-party computation** | Privacy-preserving computation, blockchain, verifiable computation without revealing inputs | Security/cryptography engineers, blockchain engineers |
| **Lattice-based cryptography (deep dive)** | Full mathematical detail behind NIST's post-quantum standards (ML-KEM, ML-DSA), beyond the overview in Ch. 24 | Cryptography engineers preparing PQC migrations |
| **Quantum computing math (qubits, superposition, Shor's/Grover's algorithms)** | Understanding the actual mechanism of the quantum threat referenced in Ch. 19/24 | Security engineers tracking PQC timelines, quantum researchers |
| **Computational geometry** | Collision detection, GIS systems, mesh processing beyond Ch. 7's basic trigonometry | Graphics/game engineers, GIS engineers |
| **Robotics kinematics & linear dynamical systems** | Multi-joint motion planning, state-space control beyond Ch. 29's single-loop PID | Robotics/embedded control engineers |
| **Actuarial & financial mathematics (Black-Scholes, options pricing, time value of money)** | Risk-adjusted return modeling, fintech product engineering | Fintech/quant engineers |
| **Compiler & type theory math (lambda calculus, type systems)** | Programming language design, static analysis, formal type-safety guarantees | Compiler/PL engineers |
| **Distributed systems formal proofs (consensus correctness, CAP theorem formalism)** | Rigorously proving a consensus protocol (Raft, Paxos) or database is correct under failure | Distributed systems/database engineers |

If your role touches any of these regularly, treat this table as a map of
where to go next, not a gap in what you needed to know today.
